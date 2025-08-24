package utils

import (
	"fmt"
	"github.com/aldian78/go-react-ecommerce/backend/internal/model"
	protoApi "github.com/aldian78/go-react-ecommerce/proto/pb/api"
	gtc "github.com/shengyanli1982/go-trycatch"
	"go-micro.dev/v4/logger"
	"runtime/debug"
)

func LoggingReqRes(req *protoApi.APIREQ) (res *protoApi.APIRES) {
	gtc.New().
		Try(func() error {
			fmt.Println("Execute... ", req.Headers["Request-Func"])
			return nil
		}).
		Catch(func(err error) {
			// Tangkap error atau panic
			logger.Errorf("[%s] exception: %v\nStack:\n%s",
				req.Headers["Request-ID"], err, string(debug.Stack()))

			res.Response = InternalServerError()
		}).
		Finally(func() {
			// Logging response
			logger.Infof("[%s] response: %v",
				req.Headers["Request-ID"], string(res.Response))
		}).
		Do()

	return res
}

func SetParamSession(params map[string]string) *model.ParseParamJWT2 {
	//param from jwt
	paramJwt := model.ParseParamJWT2{
		CustomerId: params["customerId"],
		Name:       params["name"],
		FullName:   params["fullName"],
		Email:      params["email"],
		Role:       params["role"],
	}

	return &paramJwt
}
