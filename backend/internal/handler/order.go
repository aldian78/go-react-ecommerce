package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aldian78/go-react-ecommerce/backend/internal/constants"
	"github.com/aldian78/go-react-ecommerce/backend/internal/dto"
	"github.com/aldian78/go-react-ecommerce/backend/internal/model"
	"github.com/aldian78/go-react-ecommerce/backend/internal/repository"
	"github.com/aldian78/go-react-ecommerce/backend/internal/service"
	"github.com/aldian78/go-react-ecommerce/backend/internal/utils"
	protoApi "github.com/aldian78/go-react-ecommerce/proto/pb/api"
	"github.com/aldian78/go-react-ecommerce/proto/pb/basecommon"
	"github.com/aldian78/go-react-ecommerce/proto/pb/order"
	gtc "github.com/shengyanli1982/go-trycatch"
	"go-micro.dev/v4/logger"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"
)

type OrderHandler struct {
	webhookService service.IWebhookService
	orderService   service.IOrderService
}

func NewOrderHandler(db *sql.DB) *OrderHandler {
	orderRepo := repository.NewOrderRepository(db)
	productRepo := repository.NewProductRepository(db)

	webhookService := service.NewWebhookService(orderRepo)
	orderSvc := service.NewOrderService(db, orderRepo, productRepo)
	return &OrderHandler{
		webhookService: webhookService,
		orderService:   orderSvc,
	}
}

func (oh *OrderHandler) CreateOrder(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
	defer gtc.New().
		Try(func() error {
			fmt.Println("Execute... ", req.Headers["Request-Func"])
			return nil
		}).
		Catch(func(err error) {
			// Tangkap error atau panic
			logger.Errorf("[%s] exception: %v\nStack:\n%s",
				req.Headers["Request-ID"], err, string(debug.Stack()))

			res.Response = utils.InternalServerError()
		}).
		Finally(func() {
			// Logging response
			logger.Infof("[%s] response: %v",
				req.Headers["Request-ID"], string(res.Response))
		}).
		Do()

	jsonReq, _ := json.Marshal(req)
	logger.Infof("[%s] request: %v", req.Headers["Request-ID"], string(jsonReq))

	var products []*order.CreateOrderRequestProductItem
	_ = json.Unmarshal([]byte(req.Params["products"]), &products)

	request := &order.CreateOrderRequest{
		Address:     req.Params["address"],
		FullName:    req.Params["fullName"],
		Notes:       req.Params["notes"],
		PhoneNumber: req.Params["phoneNumber"],
		Products:    products,
	}

	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		panic(err.Error())
	}
	if validationErrors != nil {
		errorsMsg := utils.LoopValidationError(validationErrors)
		res.Response = utils.Error(400, errorsMsg)
		return nil
	}

	//param from jwt
	customerParams := make(map[string]string)
	customerParams["email"] = req.Params["email"]
	customerParams["role"] = req.Params["role"]
	customerParams["fullName"] = req.Params["fullName"]
	customerParams["custId"] = req.Params["customerId"]

	result, err := oh.orderService.CreateOrder(ctx, request, customerParams)
	if err != nil {
		panic(err.Error())
	}

	if result.Base.StatusCode != 200 {
		res.Response = utils.Error(result.Base.StatusCode, result.Base.Message)
		return nil
	}

	response := &model.MResponse{
		Id:      result.Id,
		Message: result.Base.Message,
	}

	res.Response = utils.ResSuccess(response)
	return nil
}

func (oh *OrderHandler) ListOrder(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
	defer gtc.New().
		Try(func() error {
			fmt.Println("Execute... ", req.Headers["Request-Func"])
			return nil
		}).
		Catch(func(err error) {
			// Tangkap error atau panic
			logger.Errorf("[%s] exception: %v\nStack:\n%s",
				req.Headers["Request-ID"], err, string(debug.Stack()))

			res.Response = utils.InternalServerError()
		}).
		Finally(func() {
			// Logging response
			logger.Infof("[%s] response: %v",
				req.Headers["Request-ID"], string(res.Response))
		}).
		Do()

	jsonReq, _ := json.Marshal(req)
	logger.Infof("[%s] request: %v", req.Headers["Request-ID"], string(jsonReq))

	var pagingRes basecommon.PaginationRequest
	_ = json.Unmarshal([]byte(req.Params["pagination"]), &pagingRes)

	request := &order.ListOrderRequest{
		Pagination: &pagingRes,
	}

	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		panic(err.Error())
	}
	if validationErrors != nil {
		errorsMsg := utils.LoopValidationError(validationErrors)
		res.Response = utils.Error(400, errorsMsg)
		return nil
	}

	//param from jwt
	customerParams := make(map[string]string)
	customerParams["email"] = req.Params["email"]
	customerParams["role"] = req.Params["role"]
	customerParams["fullName"] = req.Params["fullName"]
	customerParams["custId"] = req.Params["customerId"]

	result, err := oh.orderService.ListOrder(ctx, request, customerParams)
	if err != nil {
		panic(err.Error())
	}

	if result.Base.StatusCode != 200 {
		res.Response = utils.Error(result.Base.StatusCode, result.Base.Message)
		return nil
	}

	res.Response = utils.ResSuccess(result)
	return nil
}

func (oh *OrderHandler) ListOrderAdmin(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
	defer gtc.New().
		Try(func() error {
			fmt.Println("Execute... ", req.Headers["Request-Func"])
			return nil
		}).
		Catch(func(err error) {
			// Tangkap error atau panic
			logger.Errorf("[%s] exception: %v\nStack:\n%s",
				req.Headers["Request-ID"], err, string(debug.Stack()))

			res.Response = utils.InternalServerError()
		}).
		Finally(func() {
			// Logging response
			logger.Infof("[%s] response: %v",
				req.Headers["Request-ID"], string(res.Response))
		}).
		Do()

	jsonReq, _ := json.Marshal(req)
	logger.Infof("[%s] request: %v", req.Headers["Request-ID"], string(jsonReq))

	var pagingRes basecommon.PaginationRequest
	_ = json.Unmarshal([]byte(req.Params["pagination"]), &pagingRes)

	request := &order.ListOrderAdminRequest{
		Pagination: &pagingRes,
	}

	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		panic(err.Error())
	}
	if validationErrors != nil {
		errorsMsg := utils.LoopValidationError(validationErrors)
		res.Response = utils.Error(400, errorsMsg)
		return nil
	}

	//param from jwt
	customerParams := make(map[string]string)
	customerParams["email"] = req.Params["email"]
	customerParams["role"] = req.Params["role"]
	customerParams["fullName"] = req.Params["fullName"]
	customerParams["custId"] = req.Params["customerId"]

	result, err := oh.orderService.ListOrderAdmin(ctx, request, customerParams)
	if err != nil {
		panic(err.Error())
	}

	res.Response = utils.ResSuccess(result)
	return nil
}

func (oh *OrderHandler) DetailOrder(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
	defer gtc.New().
		Try(func() error {
			fmt.Println("Execute... ", req.Headers["Request-Func"])
			return nil
		}).
		Catch(func(err error) {
			// Tangkap error atau panic
			logger.Errorf("[%s] exception: %v\nStack:\n%s",
				req.Headers["Request-ID"], err, string(debug.Stack()))

			res.Response = utils.InternalServerError()
		}).
		Finally(func() {
			// Logging response
			logger.Infof("[%s] response: %v",
				req.Headers["Request-ID"], string(res.Response))
		}).
		Do()

	jsonReq, _ := json.Marshal(req)
	logger.Infof("[%s] request: %v", req.Headers["Request-ID"], string(jsonReq))

	request := &order.DetailOrderRequest{
		Id: req.Params["id"],
	}

	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		panic(err.Error())
	}
	if validationErrors != nil {
		errorsMsg := utils.LoopValidationError(validationErrors)
		res.Response = utils.Error(400, errorsMsg)
		return nil
	}

	//param from jwt
	customerParams := make(map[string]string)
	customerParams["email"] = req.Params["email"]
	customerParams["role"] = req.Params["role"]
	customerParams["fullName"] = req.Params["fullName"]
	customerParams["custId"] = req.Params["customerId"]

	result, err := oh.orderService.DetailOrder(ctx, request, customerParams)
	if err != nil {
		panic(err.Error())
	}

	res.Response = utils.ResSuccess(result)
	return nil
}

func (oh *OrderHandler) UpdateOrderStatus(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
	defer gtc.New().
		Try(func() error {
			fmt.Println("Execute... ", req.Headers["Request-Func"])
			return nil
		}).
		Catch(func(err error) {
			// Tangkap error atau panic
			logger.Errorf("[%s] exception: %v\nStack:\n%s",
				req.Headers["Request-ID"], err, string(debug.Stack()))

			res.Response = utils.InternalServerError()
		}).
		Finally(func() {
			// Logging response
			logger.Infof("[%s] response: %v",
				req.Headers["Request-ID"], string(res.Response))
		}).
		Do()

	jsonReq, _ := json.Marshal(req)
	logger.Infof("[%s] request: %v", req.Headers["Request-ID"], string(jsonReq))

	request := &order.UpdateOrderStatusRequest{
		OrderId:       req.Params["id"],
		NewStatusCode: req.Params["statusOrder"],
	}

	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		panic(err.Error())
	}
	if validationErrors != nil {
		errorsMsg := utils.LoopValidationError(validationErrors)
		res.Response = utils.Error(400, errorsMsg)
		return nil
	}

	//param from jwt
	customerParams := make(map[string]string)
	customerParams["email"] = req.Params["email"]
	customerParams["role"] = req.Params["role"]
	customerParams["fullName"] = req.Params["fullName"]
	customerParams["custId"] = req.Params["customerId"]

	result, err := oh.orderService.UpdateOrderStatus(ctx, request, customerParams)
	if err != nil {
		panic(err.Error())
	}

	response := &model.MResponse{
		Message: result.Base.Message,
	}

	res.Response = utils.ResSuccess(response)
	return nil
}

func (wh *OrderHandler) WebhookXendit(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
	defer gtc.New().
		Try(func() error {
			fmt.Println("Execute... ", req.Headers["Request-Func"])
			return nil
		}).
		Catch(func(err error) {
			// Tangkap error atau panic
			logger.Errorf("[%s] exception: %v\nStack:\n%s",
				req.Headers["Request-ID"], err, string(debug.Stack()))

			res.Response = utils.InternalServerError()
		}).
		Finally(func() {
			// Logging response
			logger.Infof("[%s] response: %v",
				req.Headers["Request-ID"], string(res.Response))
		}).
		Do()

	jsonReq, _ := json.Marshal(req)
	logger.Infof("[%s] request: %v", req.Headers["Request-ID"], string(jsonReq))

	var isHigh bool
	if req.Params["is_high"] == "true" {
		isHigh = true
	} else {
		isHigh = false
	}

	receiveAmount, _ := strconv.Atoi(req.Params["adjusted_received_amount"])
	amount, _ := strconv.Atoi(req.Params["amount"])
	fPaidAmount, _ := strconv.Atoi(req.Params["fees_paid_amount"])
	paidAmount, _ := strconv.Atoi(req.Params["paid_amount"])

	createdTime, err := time.Parse(time.RFC3339Nano, req.Params["created"])
	if err != nil {
		// Handle the error (e.g., the string didn't match the layout)
		panic(err)
	}

	paidAtTime, err := time.Parse(time.RFC3339Nano, req.Params["paid_at"])
	if err != nil {
		// Handle the error (e.g., the string didn't match the layout)
		panic(err)
	}

	updatedTime, err := time.Parse(time.RFC3339Nano, req.Params["updated"])
	if err != nil {
		// Handle the error (e.g., the string didn't match the layout)
		panic(err)
	}

	request := &dto.XenditInvoiceRequest{
		ExternalID:             req.Params["external_id"],
		AdjustedReceivedAmount: receiveAmount,
		Amount:                 amount,
		BankCode:               req.Params["bank_code"],
		Created:                createdTime,
		Currency:               req.Params["currency"],
		Description:            req.Params["description"],
		FeesPaidAmount:         fPaidAmount,
		ID:                     req.Params["id"],
		IsHigh:                 isHigh,
		MerchantName:           req.Params["merchant_name"],
		PaidAmount:             paidAmount,
		PaidAt:                 paidAtTime,
		PayerEmail:             req.Params["payer_email"],
		PaymentChannel:         req.Params["payment_channel"],
		PaymentDestination:     req.Params["payment_destination"],
		PaymentMethod:          req.Params["payment_method"],
		Status:                 req.Params["status"],
		Updated:                updatedTime,
		UserID:                 req.Params["user_id"],
	}

	err = wh.webhookService.ReceiveInvoice(ctx, request)
	if err != nil {
		logger.Infof("error Receive Invoice : %s ", err.Error())
		res.Response = utils.Error(http.StatusInternalServerError, "Internal Server Error")
		return nil
	}

	res.Response = utils.Error(http.StatusOK, constants.ResSuccess)
	return nil
}
