package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aldian78/go-react-ecommerce/backend/internal/repository"
	"github.com/aldian78/go-react-ecommerce/backend/internal/service"
	"github.com/aldian78/go-react-ecommerce/backend/internal/utils"
	protoApi "github.com/aldian78/go-react-ecommerce/proto/pb/api"
	"github.com/aldian78/go-react-ecommerce/proto/pb/basecommon"
	"github.com/aldian78/go-react-ecommerce/proto/pb/product"
	gtc "github.com/shengyanli1982/go-trycatch"
	"go-micro.dev/v4/logger"
	"runtime/debug"
	"strconv"
)

type ProductHandler struct {
	productService service.IProductService
}

func NewProductHandler(db *sql.DB) *ProductHandler {
	productRepository := repository.NewProductRepository(db)
	productSrv := service.NewProductService(productRepository)

	return &ProductHandler{productService: productSrv}
}

func (ph *ProductHandler) CreateProduct(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	priceFloat, err := strconv.ParseFloat(req.Params["price"], 64)
	request := &product.CreateProductRequest{
		Name:          req.Params["name"],
		Description:   req.Params["description"],
		Price:         priceFloat,
		ImageFileName: req.Params["imageFileName"],
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

	result, err := ph.productService.CreateProduct(ctx, request, customerParams)
	if err != nil {
		panic(err.Error())
	}

	res.Response = utils.ResSuccess(result)
	return nil
}

func (ph *ProductHandler) DetailProduct(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	request := &product.DetailProductRequest{
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

	result, err := ph.productService.DetailProduct(ctx, request)
	if err != nil {
		panic(err.Error())
	}

	res.Response = utils.ResSuccess(result)
	return nil
}

func (ph *ProductHandler) EditProduct(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	priceFloat, err := strconv.ParseFloat(req.Params["price"], 64)
	request := &product.EditProductRequest{
		Id:            req.Params["id"],
		Name:          req.Params["name"],
		Description:   req.Params["description"],
		Price:         priceFloat,
		ImageFileName: req.Params["imageFileName"],
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

	result, err := ph.productService.EditProduct(ctx, request, customerParams)
	if err != nil {
		panic(err.Error())
	}

	res.Response = utils.ResSuccess(result)
	return nil
}

func (ph *ProductHandler) DeleteProduct(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	request := &product.DeleteProductRequest{
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
	result, err := ph.productService.DeleteProduct(ctx, request, customerParams)
	if err != nil {
		panic(err.Error())
	}

	res.Response = utils.ResSuccess(result)
	return nil
}

func (ph *ProductHandler) ListProduct(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	mapParams := make(map[string]string)
	json.Unmarshal([]byte(req.Params["pagination"]), &mapParams)

	current, _ := strconv.Atoi(mapParams["currentPage"])
	itemPerPage, _ := strconv.Atoi(mapParams["itemPerPage"])

	request := &product.ListProductRequest{
		Pagination: &basecommon.PaginationRequest{
			CurrentPage: int32(current),
			ItemPerPage: int32(itemPerPage),
			Sort: &basecommon.PaginationSortRequest{
				Field:     "created_at",
				Direction: "DESC",
			},
		},
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

	result, err := ph.productService.ListProduct(ctx, request)
	if err != nil {
		panic(err.Error())
	}

	res.Response = utils.ResSuccess(result)
	return nil
}

func (ph *ProductHandler) ListProductAdmin(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	mapParams := make(map[string]string)
	json.Unmarshal([]byte(req.Params["pagination"]), &mapParams)

	current, _ := strconv.Atoi(mapParams["currentPage"])
	itemPerPage, _ := strconv.Atoi(mapParams["itemPerPage"])

	request := &product.ListProductAdminRequest{
		Pagination: &basecommon.PaginationRequest{
			CurrentPage: int32(current),
			ItemPerPage: int32(itemPerPage),
			Sort: &basecommon.PaginationSortRequest{
				Field:     "created_at",
				Direction: "DESC",
			},
		},
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
	result, err := ph.productService.ListProductAdmin(ctx, request, customerParams)
	if err != nil {
		panic(err.Error())
	}

	res.Response = utils.ResSuccess(result)
	return nil
}

func (ph *ProductHandler) HighlightProducts(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	request := &product.HighlightProductsRequest{}
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		panic(err.Error())
	}
	if validationErrors != nil {
		errorsMsg := utils.LoopValidationError(validationErrors)
		res.Response = utils.Error(400, errorsMsg)
		return nil
	}

	result, err := ph.productService.HighlightProducts(ctx, request)
	if err != nil {
		panic(err.Error())
	}

	res.Response = utils.ResSuccess(result)
	return nil
}
