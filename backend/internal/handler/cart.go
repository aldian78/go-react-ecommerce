package handler

import (
	"context"
	baseutil "github.com/aldian78/go-react-ecommerce/common/utils"

	"github.com/aldian78/go-react-ecommerce/backend/internal/service"
	"github.com/aldian78/go-react-ecommerce/backend/internal/utils"
	"github.com/aldian78/go-react-ecommerce/proto/pb/cart"
)

type cartHandler struct {
	cart.UnimplementedCartServiceServer

	cartService service.ICartService
}

func (ch *cartHandler) AddProductToCart(ctx context.Context, request *cart.AddProductToCartRequest) (*cart.AddProductToCartResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &cart.AddProductToCartResponse{
			Base: baseutil.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ch.cartService.AddProductToCart(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ch *cartHandler) ListCart(ctx context.Context, request *cart.ListCartRequest) (*cart.ListCartResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &cart.ListCartResponse{
			Base: baseutil.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ch.cartService.ListCart(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ch *cartHandler) DeleteCart(ctx context.Context, request *cart.DeleteCartRequest) (*cart.DeleteCartResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &cart.DeleteCartResponse{
			Base: baseutil.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ch.cartService.DeleteCart(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ch *cartHandler) UpdateCartQuantity(ctx context.Context, request *cart.UpdateCartQuantityRequest) (*cart.UpdateCartQuantityResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &cart.UpdateCartQuantityResponse{
			Base: baseutil.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ch.cartService.UpdateCartQuantity(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewCartHandler(cartService service.ICartService) *cartHandler {
	return &cartHandler{
		cartService: cartService,
	}
}
