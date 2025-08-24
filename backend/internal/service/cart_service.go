package service

import (
	"context"
	"fmt"
	"github.com/aldian78/go-react-ecommerce/backend/internal/entity"
	"github.com/aldian78/go-react-ecommerce/backend/internal/model"
	baseutil "github.com/aldian78/go-react-ecommerce/common/utils"
	"go-micro.dev/v4/logger"
	"os"
	"time"

	"github.com/aldian78/go-react-ecommerce/backend/internal/repository"
	"github.com/aldian78/go-react-ecommerce/proto/pb/cart"
	"github.com/google/uuid"
)

type ICartService interface {
	AddProductToCart(ctx context.Context, request *cart.AddProductToCartRequest, addParams *model.ParseParamJWT2) (*cart.AddProductToCartResponse, error)
	ListCart(ctx context.Context, request *cart.ListCartRequest, addParams *model.ParseParamJWT2) (*cart.ListCartResponse, error)
	DeleteCart(ctx context.Context, request *cart.DeleteCartRequest, addParams *model.ParseParamJWT2) (*cart.DeleteCartResponse, error)
	UpdateCartQuantity(ctx context.Context, request *cart.UpdateCartQuantityRequest, addParams *model.ParseParamJWT2) (*cart.UpdateCartQuantityResponse, error)
}

type cartService struct {
	productRepository repository.IProductRepository
	cartRepository    repository.ICartRepository
}

func (cs *cartService) AddProductToCart(ctx context.Context, request *cart.AddProductToCartRequest, addParams *model.ParseParamJWT2) (*cart.AddProductToCartResponse, error) {
	custId := addParams.CustomerId
	custName := addParams.Name

	productEntity, err := cs.productRepository.GetProductById(ctx, request.ProductId)
	if err != nil {
		logger.Infof("Get Product By id : %v", err.Error())
		return &cart.AddProductToCartResponse{
			Base: baseutil.BadRequestResponse("Product id not valid"),
		}, nil
	}
	if productEntity == nil {
		return &cart.AddProductToCartResponse{
			Base: baseutil.NotFoundResponse("Product not found"),
		}, nil
	}

	cartEntity, err := cs.cartRepository.GetCartByProductAndUserId(ctx, request.ProductId, custId)
	if err != nil {
		logger.Infof("Get cart failed: %v", err.Error())
		return &cart.AddProductToCartResponse{
			Base: baseutil.BadRequestResponse("Add product to cart failed"),
		}, nil
	}

	if cartEntity != nil {
		now := time.Now()
		cartEntity.Quantity += 1
		cartEntity.UpdatedAt = &now
		cartEntity.UpdatedBy = &custId

		err = cs.cartRepository.UpdateCart(ctx, cartEntity)
		if err != nil {
			logger.Infof("Update cart failed: %v", err.Error())
			return &cart.AddProductToCartResponse{
				Base: baseutil.BadRequestResponse("Update cart failed"),
				Id:   cartEntity.Id,
			}, nil
		}

		return &cart.AddProductToCartResponse{
			Base: baseutil.SuccessResponse("Add product to cart success"),
			Id:   cartEntity.Id,
		}, nil
	}

	newCartEntity := entity.UserCart{
		Id:        uuid.NewString(),
		UserId:    custId,
		ProductId: request.ProductId,
		Quantity:  1,
		CreatedAt: time.Now(),
		CreatedBy: custName,
	}

	err = cs.cartRepository.CreateNewCart(ctx, &newCartEntity)
	if err != nil {
		logger.Infof("Create new cart failed: %v", err.Error())
		return &cart.AddProductToCartResponse{
			Base: baseutil.BadRequestResponse("Create Cart failed"),
			Id:   newCartEntity.Id,
		}, nil
	}

	return &cart.AddProductToCartResponse{
		Base: baseutil.SuccessResponse("Add product to cart success"),
		Id:   newCartEntity.Id,
	}, nil
}

func (cs *cartService) ListCart(ctx context.Context, request *cart.ListCartRequest, addParams *model.ParseParamJWT2) (*cart.ListCartResponse, error) {
	custId := addParams.CustomerId

	carts, err := cs.cartRepository.GetListCart(ctx, custId)
	if err != nil {
		logger.Infof("Get List cart failed: %v", err.Error())
		return &cart.ListCartResponse{
			Base: baseutil.BadRequestResponse("Get List Cart failed"),
		}, nil
	}

	var items []*cart.ListCartResponseItem = make([]*cart.ListCartResponseItem, 0)
	for _, cartEntity := range carts {
		item := cart.ListCartResponseItem{
			CartId:          cartEntity.Id,
			ProductId:       cartEntity.Product.Id,
			ProductName:     cartEntity.Product.Name,
			ProductImageUrl: fmt.Sprintf("%s/product/%s", os.Getenv("STORAGE_SERVICE_URL"), cartEntity.Product.ImageFileName),
			ProductPrice:    cartEntity.Product.Price,
			Quantity:        int64(cartEntity.Quantity),
		}

		items = append(items, &item)
	}

	return &cart.ListCartResponse{
		Base:  baseutil.SuccessResponse("Get list cart success"),
		Items: items,
	}, nil
}

func (cs *cartService) DeleteCart(ctx context.Context, request *cart.DeleteCartRequest, addParams *model.ParseParamJWT2) (*cart.DeleteCartResponse, error) {
	custId := addParams.CustomerId

	cartEntity, err := cs.cartRepository.GetCartById(ctx, request.CartId)
	if err != nil {
		return nil, err
	}
	if cartEntity == nil {
		return &cart.DeleteCartResponse{
			Base: baseutil.NotFoundResponse("Cart not found"),
		}, nil
	}

	if cartEntity.UserId != custId {
		return &cart.DeleteCartResponse{
			Base: baseutil.BadRequestResponse("Cart user is is not matched"),
		}, nil
	}

	err = cs.cartRepository.DeleteCart(ctx, request.CartId)
	if err != nil {
		logger.Infof("Delete cart failed: %v", err.Error())
		return &cart.DeleteCartResponse{
			Base: baseutil.BadRequestResponse("Delete cart failed"),
		}, nil
	}

	return &cart.DeleteCartResponse{
		Base: baseutil.SuccessResponse("Delete cart success"),
	}, nil
}

func (cs *cartService) UpdateCartQuantity(ctx context.Context, request *cart.UpdateCartQuantityRequest, addParams *model.ParseParamJWT2) (*cart.UpdateCartQuantityResponse, error) {
	custId := addParams.CustomerId
	custName := addParams.Name

	cartEntity, err := cs.cartRepository.GetCartById(ctx, request.CartId)
	if err != nil {
		return nil, err
	}
	if cartEntity == nil {
		return &cart.UpdateCartQuantityResponse{
			Base: baseutil.NotFoundResponse("Cart not found"),
		}, nil
	}

	if cartEntity.UserId != custId {
		return &cart.UpdateCartQuantityResponse{
			Base: baseutil.BadRequestResponse("Cart user id is not matched"),
		}, nil
	}

	if request.NewQuantity == 0 {
		err = cs.cartRepository.DeleteCart(ctx, request.CartId)
		if err != nil {
			logger.Infof("Delete cart qty failed: %v", err.Error())
			return &cart.UpdateCartQuantityResponse{
				Base: baseutil.BadRequestResponse("Delete cart qty failed"),
			}, nil
		}

		return &cart.UpdateCartQuantityResponse{
			Base: baseutil.SuccessResponse("Update cart quantity success"),
		}, nil
	}
	now := time.Now()
	cartEntity.Quantity = int(request.NewQuantity)
	cartEntity.UpdatedAt = &now
	cartEntity.UpdatedBy = &custName

	err = cs.cartRepository.UpdateCart(ctx, cartEntity)
	if err != nil {
		logger.Infof("Update cart qty failed: %v", err.Error())
		return &cart.UpdateCartQuantityResponse{
			Base: baseutil.BadRequestResponse("Update cart qty failed"),
		}, nil
	}

	return &cart.UpdateCartQuantityResponse{
		Base: baseutil.SuccessResponse("Update cart quantity success"),
	}, nil
}

func NewCartService(productRepository repository.IProductRepository, cartRepository repository.ICartRepository) ICartService {
	return &cartService{
		productRepository: productRepository,
		cartRepository:    cartRepository,
	}
}
