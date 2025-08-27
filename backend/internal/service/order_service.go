package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aldian78/go-react-ecommerce/backend/internal/entity"
	baseutil "github.com/aldian78/go-react-ecommerce/common/utils"
	operatingsystem "os"
	"runtime/debug"
	"time"

	"github.com/aldian78/go-react-ecommerce/backend/internal/repository"
	"github.com/aldian78/go-react-ecommerce/proto/pb/order"
	"github.com/google/uuid"
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IOrderService interface {
	CreateOrder(ctx context.Context, request *order.CreateOrderRequest, params map[string]string) (*order.CreateOrderResponse, error)
	ListOrderAdmin(ctx context.Context, request *order.ListOrderAdminRequest, params map[string]string) (*order.ListOrderAdminResponse, error)
	ListOrder(ctx context.Context, request *order.ListOrderRequest, params map[string]string) (*order.ListOrderResponse, error)
	DetailOrder(ctx context.Context, request *order.DetailOrderRequest, params map[string]string) (*order.DetailOrderResponse, error)
	UpdateOrderStatus(ctx context.Context, request *order.UpdateOrderStatusRequest, params map[string]string) (*order.UpdateOrderStatusResponse, error)
}

type orderService struct {
	db                *sql.DB
	orderRepository   repository.IOrderRepository
	productRepository repository.IProductRepository
}

func NewOrderService(db *sql.DB, orderRepository repository.IOrderRepository, productRepository repository.IProductRepository) IOrderService {
	return &orderService{
		db:                db,
		orderRepository:   orderRepository,
		productRepository: productRepository,
	}
}

func (os *orderService) CreateOrder(ctx context.Context, request *order.CreateOrderRequest, params map[string]string) (*order.CreateOrderResponse, error) {
	tx, err := os.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := recover(); e != nil {
			if tx != nil {
				tx.Rollback()
			}

			debug.PrintStack()
			panic(e)
		}
	}()
	defer func() {
		if err != nil && tx != nil {
			tx.Rollback()
		}
	}()

	orderRepo := os.orderRepository.WithTransaction(tx)
	productRepo := os.productRepository.WithTransaction(tx)

	numbering, err := orderRepo.GetNumbering(ctx, "order")
	if err != nil {
		return nil, err
	}

	var productIds = make([]string, len(request.Products))
	for i := range request.Products {
		productIds[i] = request.Products[i].Id
	}

	products, err := productRepo.GetProductsByIds(ctx, productIds)
	if err != nil {
		return nil, err
	}

	productMap := make(map[string]*entity.Product)
	for i := range products {
		productMap[products[i].Id] = products[i]
	}

	var total float64 = 0
	for _, p := range request.Products {
		if productMap[p.Id] == nil {
			return &order.CreateOrderResponse{
				Base: baseutil.NotFoundResponse(fmt.Sprintf("Product %s not found", p.Id)),
			}, nil
		}
		total += productMap[p.Id].Price * float64(p.Quantity)
	}

	now := time.Now()
	expiredAt := now.Add(24 * time.Hour)
	orderEntity := entity.Order{
		Id:              uuid.NewString(),
		Number:          fmt.Sprintf("ORD-%d%08d", now.Year(), numbering.Number),
		UserId:          params["custId"],
		OrderStatusCode: entity.OrderStatusCodeUnpaid,
		UserFullName:    request.FullName,
		Address:         request.Address,
		PhoneNumber:     request.PhoneNumber,
		Notes:           &request.Notes,
		Total:           total,
		ExpiredAt:       &expiredAt,
		CreatedAt:       now,
		CreatedBy:       params["fullName"],
	}

	invoiceItems := make([]xendit.InvoiceItem, 0)
	for _, p := range request.Products {
		prod := productMap[p.Id]
		if prod != nil {
			invoiceItems = append(invoiceItems, xendit.InvoiceItem{
				Name:     prod.Name,
				Price:    prod.Price,
				Quantity: int(p.Quantity),
			})
		}
	}
	xenditInvoice, xenditErr := invoice.CreateWithContext(ctx, &invoice.CreateParams{
		ExternalID: orderEntity.Id,
		Amount:     total,
		Customer: xendit.InvoiceCustomer{
			GivenNames: params["fullName"],
		},
		Currency:           "IDR",
		SuccessRedirectURL: fmt.Sprintf("%s/checkout/%s/success", operatingsystem.Getenv("FRONTEND_BASE_URL"), orderEntity.Id),
		Items:              invoiceItems,
	})
	if xenditErr != nil {
		err = xenditErr
		return nil, err
	}

	orderEntity.XenditInvoiceId = &xenditInvoice.ID
	orderEntity.XenditInvoiceUrl = &xenditInvoice.InvoiceURL

	err = orderRepo.CreateOrder(ctx, &orderEntity)
	if err != nil {
		return nil, err
	}

	for _, p := range request.Products {
		var orderItem = entity.OrderItem{
			Id:                   uuid.NewString(),
			ProductId:            p.Id,
			ProductName:          productMap[p.Id].Name,
			ProductImageFileName: productMap[p.Id].ImageFileName,
			ProductPrice:         productMap[p.Id].Price,
			Quantity:             p.Quantity,
			OrderId:              orderEntity.Id,
			CreatedAt:            now,
			CreatedBy:            params["fullName"],
		}

		err = orderRepo.CreateOrderItem(ctx, &orderItem)
		if err != nil {
			return nil, err
		}
	}

	numbering.Number++
	err = orderRepo.UpdateNumbering(ctx, numbering)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &order.CreateOrderResponse{
		Base: baseutil.SuccessResponse("Create order success"),
		Id:   orderEntity.Id,
	}, nil
}

func (os *orderService) ListOrderAdmin(ctx context.Context, request *order.ListOrderAdminRequest, params map[string]string) (*order.ListOrderAdminResponse, error) {

	if params["role"] != entity.UserRoleAdmin {
		return nil, baseutil.UnauthenticatedResponse()
	}

	orders, metadata, err := os.orderRepository.GetListOrderAdminPagination(ctx, request.Pagination)
	if err != nil {
		return nil, err
	}

	items := make([]*order.ListOrderAdminResponseItem, 0)
	for _, o := range orders {

		products := make([]*order.ListOrderAdminResponseItemProduct, 0)
		for _, oi := range o.Items {
			products = append(products, &order.ListOrderAdminResponseItemProduct{
				Id:       oi.ProductId,
				Name:     oi.ProductName,
				Price:    oi.ProductPrice,
				Quantity: oi.Quantity,
			})
		}

		orderStatusCode := o.OrderStatusCode
		if o.OrderStatusCode == entity.OrderStatusCodeUnpaid && time.Now().After(*o.ExpiredAt) {
			orderStatusCode = entity.OrderStatusCodeExpired
		}

		items = append(items, &order.ListOrderAdminResponseItem{
			Id:         o.Id,
			Number:     o.Number,
			Customer:   o.UserFullName,
			StatusCode: orderStatusCode,
			Total:      o.Total,
			CreatedAt:  timestamppb.New(o.CreatedAt),
			Products:   products,
		})
	}

	return &order.ListOrderAdminResponse{
		Base:       baseutil.SuccessResponse("Get list order admin success"),
		Pagination: metadata,
		Items:      items,
	}, nil
}

func (os *orderService) ListOrder(ctx context.Context, request *order.ListOrderRequest, params map[string]string) (*order.ListOrderResponse, error) {
	orders, metadata, err := os.orderRepository.GetListOrderPagination(ctx, request.Pagination, params["custId"])
	if err != nil {
		return nil, err
	}

	items := make([]*order.ListOrderResponseItem, 0)
	for _, o := range orders {

		products := make([]*order.ListOrderResponseItemProduct, 0)
		for _, oi := range o.Items {
			products = append(products, &order.ListOrderResponseItemProduct{
				Id:       oi.ProductId,
				Name:     oi.ProductName,
				Price:    oi.ProductPrice,
				Quantity: oi.Quantity,
			})
		}

		orderStatusCode := o.OrderStatusCode
		if o.OrderStatusCode == entity.OrderStatusCodeUnpaid && time.Now().After(*o.ExpiredAt) {
			orderStatusCode = entity.OrderStatusCodeExpired
		}

		xenditInvoiceUrl := ""
		if o.XenditInvoiceUrl != nil {
			xenditInvoiceUrl = *o.XenditInvoiceUrl
		}
		items = append(items, &order.ListOrderResponseItem{
			Id:               o.Id,
			Number:           o.Number,
			Customer:         o.UserFullName,
			StatusCode:       orderStatusCode,
			Total:            o.Total,
			CreatedAt:        timestamppb.New(o.CreatedAt),
			Products:         products,
			XenditInvoiceUrl: xenditInvoiceUrl,
		})
	}

	return &order.ListOrderResponse{
		Base:       baseutil.SuccessResponse("Get list order success"),
		Pagination: metadata,
		Items:      items,
	}, nil
}

func (os *orderService) DetailOrder(ctx context.Context, request *order.DetailOrderRequest, params map[string]string) (*order.DetailOrderResponse, error) {

	orderEntity, err := os.orderRepository.GetOrderById(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if params["role"] != entity.UserRoleAdmin && params["custId"] != orderEntity.UserId {
		return &order.DetailOrderResponse{
			Base: baseutil.BadRequestResponse("User id is not matched"),
		}, nil
	}

	notes := ""
	if orderEntity.Notes != nil {
		notes = *orderEntity.Notes
	}
	xenditInvoiceUrl := ""
	if orderEntity.XenditInvoiceUrl != nil {
		xenditInvoiceUrl = *orderEntity.XenditInvoiceUrl
	}

	orderStatusCode := orderEntity.OrderStatusCode
	if orderEntity.OrderStatusCode == entity.OrderStatusCodeUnpaid && time.Now().After(*orderEntity.ExpiredAt) {
		orderStatusCode = entity.OrderStatusCodeExpired
	}

	items := make([]*order.DetailOrderResponseItem, 0)
	for _, oi := range orderEntity.Items {
		items = append(items, &order.DetailOrderResponseItem{
			Id:       oi.ProductId,
			Name:     oi.ProductName,
			Price:    oi.ProductPrice,
			Quantity: oi.Quantity,
		})
	}
	return &order.DetailOrderResponse{
		Base:             baseutil.SuccessResponse("Get order detail success"),
		Id:               orderEntity.Id,
		Number:           orderEntity.Number,
		UserFullName:     orderEntity.UserFullName,
		Address:          orderEntity.Address,
		PhoneNumber:      orderEntity.PhoneNumber,
		Notes:            notes,
		OrderStatusCode:  orderStatusCode,
		CreatedAt:        timestamppb.New(orderEntity.CreatedAt),
		XenditInvoiceUrl: xenditInvoiceUrl,
		Items:            items,
		Total:            orderEntity.Total,
		ExpiredAt:        timestamppb.New(*orderEntity.ExpiredAt),
	}, nil
}

func (os *orderService) UpdateOrderStatus(ctx context.Context, request *order.UpdateOrderStatusRequest, params map[string]string) (*order.UpdateOrderStatusResponse, error) {
	role := params["role"]
	custId := params["custId"]

	orderEntity, err := os.orderRepository.GetOrderById(ctx, request.OrderId)
	if err != nil {
		return nil, err
	}
	if orderEntity == nil {
		return &order.UpdateOrderStatusResponse{
			Base: baseutil.NotFoundResponse("Order not found"),
		}, nil
	}

	if role != entity.UserRoleAdmin && orderEntity.UserId != custId {
		return &order.UpdateOrderStatusResponse{
			Base: baseutil.BadRequestResponse("User id is not matched"),
		}, nil
	}

	if request.NewStatusCode == entity.OrderStatusCodePaid {
		if role != entity.UserRoleAdmin || orderEntity.OrderStatusCode != entity.OrderStatusCodeUnpaid {
			return &order.UpdateOrderStatusResponse{
				Base: baseutil.BadRequestResponse("Update status is not allowed"),
			}, nil
		}
	} else if request.NewStatusCode == entity.OrderStatusCodeCanceled {
		if orderEntity.OrderStatusCode != entity.OrderStatusCodeUnpaid {
			return &order.UpdateOrderStatusResponse{
				Base: baseutil.BadRequestResponse("Update status is not allowed"),
			}, nil
		}
	} else if request.NewStatusCode == entity.OrderStatusCodeShipped {
		if role != entity.UserRoleAdmin || orderEntity.OrderStatusCode != entity.OrderStatusCodePaid {
			return &order.UpdateOrderStatusResponse{
				Base: baseutil.BadRequestResponse("Update status is not allowed"),
			}, nil
		}

	} else if request.NewStatusCode == entity.OrderStatusCodeDone {
		if orderEntity.OrderStatusCode != entity.OrderStatusCodeShipped {
			return &order.UpdateOrderStatusResponse{
				Base: baseutil.BadRequestResponse("Update status is not allowed"),
			}, nil
		}
	} else {
		return &order.UpdateOrderStatusResponse{
			Base: baseutil.BadRequestResponse("Invalid new status code"),
		}, nil
	}

	now := time.Now()
	orderEntity.OrderStatusCode = request.NewStatusCode
	orderEntity.UpdatedAt = &now
	orderEntity.UpdatedBy = &custId

	err = os.orderRepository.UpdateOrder(ctx, orderEntity)
	if err != nil {
		return nil, err
	}

	return &order.UpdateOrderStatusResponse{
		Base: baseutil.SuccessResponse("Update order status success"),
	}, nil
}
