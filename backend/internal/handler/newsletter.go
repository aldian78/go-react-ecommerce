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
	"github.com/aldian78/go-react-ecommerce/proto/pb/newsletter"
	gtc "github.com/shengyanli1982/go-trycatch"
	"go-micro.dev/v4/logger"
	"runtime/debug"
)

type NewsletterHandler struct {
	newsletterService service.INewsletterService
}

func NewNewsletterHandler(db *sql.DB) *NewsletterHandler {
	newsletterRepo := repository.NewNewsletterRepository(db)
	newsletterSvc := service.NewNewsletterService(newsletterRepo)
	return &NewsletterHandler{
		newsletterService: newsletterSvc,
	}
}

func (nh *NewsletterHandler) SubscribeNewsletter(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	request := &newsletter.SubcribeNewsletterRequest{
		Email:    req.Params["email"],
		FullName: req.Params["fullName"],
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

	result, err := nh.newsletterService.SubscribeNewsletter(ctx, request)
	if err != nil {
		panic(err.Error())
	}

	if result.Base.StatusCode != 200 {
		res.Response = utils.Error(result.Base.StatusCode, result.Base.Message)
		return nil
	}

	response := &model.MResponse{
		Message: result.Base.Message,
	}

	res.Response = utils.ResSuccess(response)
	return nil
}
