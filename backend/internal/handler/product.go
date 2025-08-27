package handler

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aldian78/go-react-ecommerce/backend/internal/model"
	"github.com/aldian78/go-react-ecommerce/backend/internal/repository"
	"github.com/aldian78/go-react-ecommerce/backend/internal/service"
	"github.com/aldian78/go-react-ecommerce/backend/internal/utils"
	protoApi "github.com/aldian78/go-react-ecommerce/proto/pb/api"
	"github.com/aldian78/go-react-ecommerce/proto/pb/basecommon"
	"github.com/aldian78/go-react-ecommerce/proto/pb/product"
	"github.com/gofiber/fiber/v2"
	gtc "github.com/shengyanli1982/go-trycatch"
	"go-micro.dev/v4/logger"
	"image"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
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

	if result.Base.StatusCode != 200 {
		res.Response = utils.Error(result.Base.StatusCode, result.Base.Message)
		return nil
	}

	response := &model.MGetProductRes{
		Id:          result.Id,
		Message:     result.Base.Message,
		ProductName: result.Name,
		Description: result.Description,
		Price:       result.Price,
		ImageUrl:    result.ImageUrl,
	}

	res.Response = utils.ResSuccess(response)
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

	if result.Base.StatusCode != 200 {
		res.Response = utils.Error(result.Base.StatusCode, result.Base.Message)
		return nil
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

	if result.Base.StatusCode != 200 {
		res.Response = utils.Error(result.Base.StatusCode, result.Base.Message)
		return nil
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

	if result.Base.StatusCode != 200 {
		res.Response = utils.Error(result.Base.StatusCode, result.Base.Message)
		return nil
	}

	res.Response = utils.ResSuccess(result)
	return nil
}

func (ph *ProductHandler) GetFileName(_ context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	fileNameParam := req.Params["filename"]
	filePath := path.Join("storage", "product", fileNameParam)
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			res.Response = utils.Error(http.StatusNotFound, "Not Found")
			return nil
		}

		logger.Infof(err.Error())
		res.Response = utils.Error(http.StatusInternalServerError, "Internal Server Error")
		return nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		logger.Infof(err.Error())
		res.Response = utils.Error(http.StatusInternalServerError, "Internal Server Error")
		return nil
	}

	ext := path.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)

	res.Headers["Content-Type"] = mimeType
	res.Response = utils.ResSuccess(file)
	return nil
}

func (ph *ProductHandler) UploadProductImage(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	//file, err := c.FormFile("image")
	//if err != nil {
	//	return c.Status(http.StatusBadRequest).JSON(fiber.Map{
	//		"success": false,
	//		"message": "image data not found",
	//	})
	//}

	// validasi gambar
	// validasi extension

	ext := strings.ToLower(filepath.Ext(req.Params["filename"]))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		res.Response = utils.Error(http.StatusBadRequest, "Image extension is not allowed (jpg, jpeg, png, webp")
		return nil
	}

	// validasi content type
	contentType := req.Headers["Content-Type"]
	allowedContentType := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	if !allowedContentType[contentType] {
		res.Response = utils.Error(http.StatusBadRequest, "Content type is not allowed")
		return nil
	}

	timestamp := time.Now().UnixNano()
	fileName := fmt.Sprintf("product_%d%s", timestamp, filepath.Ext(req.Params["filename"]))
	uploadPath := "./storage/product/" + fileName

	var saveImgRes bool
	if strings.Contains(req.Params["filename"], "base64") {
		sImages := strings.Split(req.Params["filename"], "base64,")

		resBool, err := ph.SaveImages(map[string]string{uploadPath: sImages[1]}, uploadPath)
		if err != nil {
			res.Response = utils.Error(http.StatusInternalServerError, "Internal server error")
			return nil
		}
		saveImgRes = resBool
	}

	result := &fiber.Map{
		"success":   saveImgRes,
		"message":   "Upload success",
		"file_name": fileName,
	}
	res.Response = utils.ResSuccess(result)
	return nil

}

func (c *ProductHandler) SaveImages(files map[string]string, location string) (bool, error) {
	logger.Info("SaveImages")
	for k, v := range files {
		dec, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			logger.Infof("error decode image: %v", err.Error())
			logger.Info("err 1")

			return false, err
		}

		f, err := os.Create(location + k)
		if err != nil {
			logger.Infof("error open location image: %v", err.Error())
			logger.Info("err 2")

			return false, err
		}
		defer func() {
			_ = f.Close()
		}()

		if _, err := f.Write(dec); err != nil {
			logger.Infof("error write image: %v", err.Error())
			logger.Info("err 3")

			return false, err
		}
		if err := f.Sync(); err != nil {
			logger.Infof("error Syncr image: %v", err.Error())
			logger.Info("err 4")

			return false, err
		}

		_, _, err = image.Decode(f)
		if err != nil {
			if f.Name() == "" {
				logger.Info(f.Name() + " seems corrupt")
				return false, err
			}
		}
	}
	return true, nil
}
