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
	"github.com/aldian78/go-react-ecommerce/proto/pb/category"
	gtc "github.com/shengyanli1982/go-trycatch"
	"go-micro.dev/v4/logger"
	"runtime/debug"
	"strconv"
)

type CategoryHandler struct {
	categoryService service.ICategoryService
}

func NewCategoryHandler(db *sql.DB) *CategoryHandler {
	categoryRepository := repository.NewCategoryRepository(db)
	categorySrv := service.NewCategoryService(categoryRepository)

	return &CategoryHandler{categoryService: categorySrv}
}

func (ch *CategoryHandler) CreateCategory(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	request := &category.CreateCategoryRequest{
		Name: req.Params["name"],
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

	result, err := ch.categoryService.CreateCategory(ctx, request, customerParams)
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

func (ch *CategoryHandler) EditCategory(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	request := &category.EditCategoryRequest{
		Id:   req.Params["id"],
		Name: req.Params["name"],
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

	result, err := ch.categoryService.EditCategory(ctx, request, customerParams)
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
func (ch *CategoryHandler) DeleteCategory(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	request := &category.DeleteCategoryRequest{
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
	result, err := ch.categoryService.DeleteCategory(ctx, request, customerParams)
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
func (ch *CategoryHandler) ListCategory(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	request := &category.ListCategoryRequest{
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

	result, err := ch.categoryService.ListCategory(ctx, request)
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
