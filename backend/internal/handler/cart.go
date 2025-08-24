package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aldian78/go-react-ecommerce/backend/internal/model"
	"github.com/aldian78/go-react-ecommerce/backend/internal/repository"
	"github.com/aldian78/go-react-ecommerce/proto/pb/api"
	gtc "github.com/shengyanli1982/go-trycatch"
	"go-micro.dev/v4/logger"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/aldian78/go-react-ecommerce/backend/internal/service"
	"github.com/aldian78/go-react-ecommerce/backend/internal/utils"
	"github.com/aldian78/go-react-ecommerce/proto/pb/cart"
)

type CartHandler struct {
	cartService service.ICartService
}

func NewCartHandler(db *sql.DB) *CartHandler {
	productRepository := repository.NewProductRepository(db)
	cartRepository := repository.NewCartRepository(db)
	cartSrv := service.NewCartService(productRepository, cartRepository)

	return &CartHandler{cartService: cartSrv}
}

func (ch *CartHandler) AddProductToCart(ctx context.Context, req *api.APIREQ, res *api.APIRES) error {
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

	request := &cart.AddProductToCartRequest{
		ProductId: req.Params["productId"],
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

	addParams := utils.SetParamSession(req.Params)
	result, err := ch.cartService.AddProductToCart(ctx, request, addParams)
	if err != nil {
		panic(err.Error())
	}

	if result.Base.StatusCode != 200 {
		res.Response = utils.Error(result.Base.StatusCode, result.Base.Message)
		return nil
	}

	response := &model.HandlerResponseCart{
		Id:      result.Id,
		Message: result.Base.Message,
	}
	res.Response = utils.ResSuccess(response)
	return nil
}

func (ch *CartHandler) ListCart(ctx context.Context, req *api.APIREQ, res *api.APIRES) error {
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

	request := &cart.ListCartRequest{}

	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		panic(err.Error())
	}
	if validationErrors != nil {
		errorsMsg := utils.LoopValidationError(validationErrors)
		res.Response = utils.Error(400, errorsMsg)
		return nil
	}

	addParams := utils.SetParamSession(req.Params)
	result, err := ch.cartService.ListCart(ctx, request, addParams)
	if err != nil {
		panic(err.Error())
	}

	if result.Base.StatusCode != 200 {
		res.Response = utils.Error(result.Base.StatusCode, result.Base.Message)
		return nil
	}

	var items []*model.ListCart = make([]*model.ListCart, 0)
	for _, cartEntity := range result.Items {
		item := model.ListCart{
			CartId:         cartEntity.CartId,
			ProductId:      cartEntity.ProductId,
			ProductName:    cartEntity.ProductName,
			ProductImagUrl: fmt.Sprintf("%s/product/%s", os.Getenv("STORAGE_SERVICE_URL"), cartEntity.ProductImageUrl),
			Price:          cartEntity.ProductPrice,
			Quantity:       cartEntity.Quantity,
		}

		items = append(items, &item)
	}

	response := &model.HandlerResponseCart{
		Message: result.Base.Message,
		Items:   items,
	}

	res.Response = utils.ResSuccess(response)
	return nil
}

func (ch *CartHandler) DeleteCart(ctx context.Context, req *api.APIREQ, res *api.APIRES) error {
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

	request := &cart.DeleteCartRequest{
		CartId: req.Params["cartId"],
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

	addParams := utils.SetParamSession(req.Params)
	result, err := ch.cartService.DeleteCart(ctx, request, addParams)
	if err != nil {
		panic(err.Error())
	}

	if result.Base.StatusCode != 200 {
		res.Response = utils.Error(result.Base.StatusCode, result.Base.Message)
		return nil
	}

	response := &model.HandlerResponseCart{
		Message: result.Base.Message,
	}

	res.Response = utils.ResSuccess(response)
	return nil
}

func (ch *CartHandler) UpdateCartQuantity(ctx context.Context, req *api.APIREQ, res *api.APIRES) error {
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

	qty, _ := strconv.Atoi(req.Params["newQty"])
	request := &cart.UpdateCartQuantityRequest{
		CartId:      req.Params["cartId"],
		NewQuantity: int64(qty),
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

	addParams := utils.SetParamSession(req.Params)
	result, err := ch.cartService.UpdateCartQuantity(ctx, request, addParams)
	if err != nil {
		panic(err.Error())
	}

	if result.Base.StatusCode != 200 {
		res.Response = utils.Error(result.Base.StatusCode, result.Base.Message)
		return nil
	}

	response := &model.HandlerResponseCart{
		Message: result.Base.Message,
	}

	res.Response = utils.ResSuccess(response)
	return nil
}
