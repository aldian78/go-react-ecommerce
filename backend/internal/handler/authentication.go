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
	"github.com/aldian78/go-react-ecommerce/proto/pb/authentication"
	gocache "github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	gtc "github.com/shengyanli1982/go-trycatch"
	"go-micro.dev/v4/logger"
	"runtime/debug"
)

type AuthenticationHandler struct {
	authService service.IAuthenticationService
}

func NewAuthenticationHandler(db *sql.DB, rdb *redis.Client, cacheService *gocache.Cache) *AuthenticationHandler {
	authRepository := repository.NewAuthenticationRepository(db)
	authSrv := service.NewAuthenticationService(authRepository, rdb, cacheService)

	return &AuthenticationHandler{authService: authSrv}
}

func (sh *AuthenticationHandler) RegisterHandler(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	protoReq := &authentication.RegisterRequest{
		FullName:             req.Params["fullName"],
		Email:                req.Params["email"],
		Password:             req.Params["password"],
		PasswordConfirmation: req.Params["passwordConfirmation"],
	}

	validationErrors, err := utils.CheckValidation(protoReq)
	if err != nil {
		return err
	}
	if validationErrors != nil {
		res.Response = utils.Error(400, "Authentication Failed")
		return nil
	}

	// Process register
	result, err := sh.authService.RegisterService(ctx, protoReq)
	if err != nil {
		return err
	}

	jsonxx, _ := json.Marshal(res)
	logger.Infof("check res register : %s", string(jsonxx))

	if result.Base.StatusCode != 200 {
		res.Response = utils.Error(result.Base.StatusCode, result.Base.Message)
		return nil
	}

	res.Response = utils.ResSuccessWithData(result.Base.StatusCode, result.Base.Message)
	return nil
}

func (sh *AuthenticationHandler) LoginHandler(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	protoReq := &authentication.LoginRequest{
		Email:    req.Params["email"],
		Password: req.Params["password"],
	}

	// Validasi params
	validationErrors, err := utils.CheckValidation(protoReq)
	if err != nil {
		return err
	}
	if validationErrors != nil {
		res.Response = utils.Error(400, "Authentication Failed")
		return nil
	}

	// Panggil services
	result, err := sh.authService.LoginService(ctx, protoReq)
	if err != nil {
		return err
	}

	mapResult := make(map[string]interface{})
	mapResult["acces_token"] = result.AccessToken
	res.Response = utils.ResSuccess(mapResult)
	return nil
}

func (sh *AuthenticationHandler) LogoutHandler(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	var protoReq *authentication.LogoutRequest
	validationErrors, err := utils.CheckValidation(protoReq)
	if err != nil {
		return err
	}
	if validationErrors != nil {
		res.Response = utils.Error(400, "Authentication Failed")
		return nil
	}

	result, err := sh.authService.LogoutService(ctx, req)
	if err != nil {
		return err
	}

	res.Response = utils.ResSuccess(result.Base)
	return nil
}

func (sh *AuthenticationHandler) ChangePasswordHandler(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	protoReq := &authentication.ChangePasswordRequest{
		OldPassword:             req.Params["oldPassword"],
		NewPassword:             req.Params["newPassword"],
		NewPasswordConfirmation: req.Params["newPasswordConfirmation"],
	}

	validationErrors, err := utils.CheckValidation(protoReq)
	if err != nil {
		return err
	}
	if validationErrors != nil {
		res.Response = utils.Error(400, "Authentication Failed")
		return nil
	}

	result, err := sh.authService.ChangePasswordService(ctx, protoReq)
	if err != nil {
		return err
	}

	res.Response = utils.ResSuccess(result.Base)
	return nil
}

func (sh *AuthenticationHandler) GetProfileHandler(ctx context.Context, req *protoApi.APIREQ, res *protoApi.APIRES) error {
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

	var protoReq *authentication.GetProfileRequest

	result, err := sh.authService.GetProfileService(ctx, protoReq)
	if err != nil {
		return err
	}

	res.Response = utils.ResSuccess(result.Base)
	return nil
}
