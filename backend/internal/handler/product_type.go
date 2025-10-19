package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aldian78/go-react-ecommerce/backend/internal/model"
	"github.com/aldian78/go-react-ecommerce/backend/internal/repository"
	"github.com/aldian78/go-react-ecommerce/backend/internal/service"
	"github.com/aldian78/go-react-ecommerce/backend/internal/utils"
	protoApi "github.com/aldian78/go-react-ecommerce/proto/pb/api"
	"github.com/aldian78/go-react-ecommerce/proto/pb/basecommon"
	ProductType "github.com/aldian78/go-react-ecommerce/proto/pb/product_type"
	gtc "github.com/shengyanli1982/go-trycatch"
	"go-micro.dev/v4/logger"
	"runtime/debug"
	"strconv"
)

type ProductTypeHandler struct {
	productTypeService service.IProductTypeService
}

func NewProductTypeHandler(db *sql.DB) *ProductTypeHandler {
	productTypeRepository := repository.NewProductTypeRepository(db)
	productTypeSrv := service.NewProductTypeService(productTypeRepository)

	return &ProductTypeHandler{productTypeService: productTypeSrv}
}

func (h *ProductTypeHandler) CreateProductType(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	request := &ProductType.CreateProductTypeRequest{
		Type: req.Params["type"],
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

	result, err := h.productTypeService.CreateProductType(ctx, request, customerParams)
	if err != nil {
		panic(err.Error())
	}

	if result.Base.StatusCode != 200 {
		res.Response = utils.Error(result.Base.StatusCode, result.Base.Message)
		return nil
	}

	response := &model.MProductRes{
		Id:      result.Id,
		Message: result.Base.Message,
	}

	res.Response = utils.ResSuccess(response)
	return nil
}

func (h *ProductTypeHandler) EditProductType(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	request := &ProductType.EditProductTypeRequest{
		Id:   req.Params["id"],
		Type: req.Params["type"],
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

	result, err := h.productTypeService.EditProductType(ctx, request, customerParams)
	if err != nil {
		panic(err.Error())
	}

	if result.Base.StatusCode != 200 {
		res.Response = utils.Error(result.Base.StatusCode, result.Base.Message)
		return nil
	}

	response := &model.MProductRes{
		Id:      result.Id,
		Message: result.Base.Message,
	}

	res.Response = utils.ResSuccess(response)
	return nil
}
func (h *ProductTypeHandler) DeleteProductType(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	request := &ProductType.DeleteProductTypeRequest{
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
	result, err := h.productTypeService.DeleteProductType(ctx, request, customerParams)
	if err != nil {
		panic(err.Error())
	}

	if result.Base.StatusCode != 200 {
		res.Response = utils.Error(result.Base.StatusCode, result.Base.Message)
		return nil
	}

	response := &model.MProductRes{
		Message: result.Base.Message,
	}

	res.Response = utils.ResSuccess(response)
	return nil
}
func (h *ProductTypeHandler) ListProductType(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	request := &ProductType.ListProductTypeRequest{
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

	result, err := h.productTypeService.ListProductType(ctx, request)
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
