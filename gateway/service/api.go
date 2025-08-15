package service

import (
	"context"
	"gateway/internal/constants"
	"gateway/internal/entity"
	"gateway/internal/grpcmiddleware"
	"gateway/internal/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/logger"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	protoApi "gateway/protoc/proto-api"
	"net/http"
)

// GoClient is go gateway client for call gRPC services

func GetConfig() map[string]string {
	return make(map[string]string)
}

// Init API services
var GoClient *grpcmiddleware.Client

func Init(baseURL string, gatewayTimeout int, corsOrigins []string, corsMehthods []string, corsHeaders []string, routesRepository repository.IRoutesRepository) http.Handler {

	//GoClient = cl.Get()

	GoClient = grpcmiddleware.NewClient() // inisialisasi GoClient di awal

	r := gin.New()
	config := cors.DefaultConfig()
	config.AllowOrigins = corsOrigins
	config.AllowMethods = corsMehthods
	config.AllowHeaders = corsHeaders
	r.Use(grpcmiddleware.Recovery, cors.New(config))

	r.NoMethod(func(c *gin.Context) {
		c.String(http.StatusMethodNotAllowed, "Method not Allowed")
	})
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Invalid URL")
	})

	var g *gin.RouterGroup
	if baseURL == "" {
		g = r.Group("/")
	} else {
		g = r.Group(baseURL)
	}

	//logger := middleware.Logger()

	apis, err := routesRepository.GetAllEndpoint(context.TODO())
	if err != nil {
		logger.Infof("apisss  ", err.Error())
		panic(err)
	}

	//db.MRateLimit = db.GetAllMRateLimitBlock()

	conf := GetConfig()
	logger.Info("REDIS :" + conf["REDIS_URL"])
	// Start the Datadog tracer
	//if v, ok := GetConfig()["DATA_DOG"]; ok && v == "Y" {
	//	tracer.Start(
	//		tracer.WithAgentAddr(gateway.GetConfig()["URL_DATA_DOG"]),
	//	)
	//	tracer.Start(tracer.WithEnv("development"), tracer.WithService("gateway"))
	//	defer tracer.Stop()
	//}

	grpcmiddleware.InitAuthorizationMiddleware()
	for _, api := range apis {

		handlers := make([]gin.HandlerFunc, 0)
		handlers = append(handlers, grpcmiddleware.Logger(api))
		logger.Infof("Loading api.Middleware %s", api.Middleware)
		handlers = append(handlers, grpcmiddleware.LoadMiddlewareFromFactory(api.Middleware)...)

		handlers = append(handlers, globalAPI(api))
		//handlers = append(handlers, audit(proto-api))
		//if v, ok := GetConfig()["DATA_DOG"]; ok && v == "Y" {
		//	handlers = append(handlers, gintrace.Middleware("gateway"))
		//}

		if api.Method == "GET" {
			g.GET(api.ContextPath, handlers...)
		} else if api.Method == "POST" {
			g.POST(api.ContextPath, handlers...)
		} else if api.Method == "PUT" {
			g.PUT(api.ContextPath, handlers...)
		} else if api.Method == "DELETE" {
			g.DELETE(api.ContextPath, handlers...)
		} else {
			logger.Warn("invalid or not supported API method: " + api.Method + ", for path: " + api.ContextPath)
		}
	}

	return r
}

func globalAPI(api *entity.ApiGateway) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Writer.Header()
		span, _ := tracer.StartSpanFromContext(c.Request.Context(), api.EndpointName+":"+api.EndpointFunc)
		defer span.Finish()

		if api.ContentType != "" {
			//c.Header(constant.ContentType, proto-api.ContentType)
			header["Content-Type"] = []string{api.ContentType}
		}

		if api.EndpointName != "" && api.EndpointFunc != "" {

			logger.Infof("call to %v:%v", api.EndpointName, api.EndpointFunc)

			//apiReq, _ := c.Get("req")
			//reqApi, ok := apiReq.(*protoApi.APIREQ)

			// Bind JSON dari request ke DTO
			//var reqx rest.APIREQDTO
			//if err := c.ShouldBindJSON(&reqx); err != nil {
			//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			//	return
			//}

			// Inisialisasi map agar tidak nil
			//if reqx.Headers == nil {
			//	reqx.Headers = make(map[string]string)
			//}
			//if reqx.Params == nil {
			//	reqx.Params = make(map[string]string)
			//}
			//
			//// Konversi ke protobuf
			//protoReq := &protoApi.APIREQ{
			//	TxType:  reqx.TxType,
			//	Headers: reqx.Headers,
			//	Params:  reqx.Params,
			//}

			//logger.Infof("Creating request with: name=%s func=%s txType=%s headers=%v params=%v",
			//	api.EndpointName, api.EndpointFunc, protoReq.TxType, protoReq.Headers, protoReq.Params)

			// Cek GoClient
			if GoClient == nil {
				logger.Error("GoClient is nil")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "GoClient not initialized"})
				return
			}

			//// Cek protoReq
			//if protoReq == nil {
			//	logger.Error("protoReq is nil before creating request")
			//	c.JSON(http.StatusInternalServerError, gin.H{"error": "protoReq is nil"})
			//	return
			//}

			apiReq, _ := c.Get(constants.APIReqKey)
			reqApi := apiReq.(*protoApi.APIREQ)

			// Panggil services
			req := GoClient.NewRequest(api.EndpointName, api.EndpointFunc, reqApi)

			apiRes := new(protoApi.APIRES)
			err := GoClient.Call(context.TODO(), req, apiRes)
			if err != nil {
				logger.Errorf("call failed: %v", err)
				c.AbortWithStatus(http.StatusServiceUnavailable)
				return
			}

			c.Set("proto-api.res", apiRes)

			if apiRes.Headers != nil {
				for k, v := range apiRes.Headers {
					header[k] = []string{v}
				}
			}
			if apiRes.Response != nil {
				//translator.Translate(c.GetHeader(constant.HeaderAcceptLanguage), apiRes)
				//
				//encryptor.EncryptResponse(c, apiRes, reqApi.Headers[constant.HeaderRequestID])

				if apiRes.HttpStatus != 0 {
					c.Status(int(apiRes.HttpStatus))
				}
				_, err := c.Writer.Write(apiRes.Response)
				if err != nil {
					logger.Errorf("writting response: %v", err)
				}
			} else {
				c.AbortWithStatus(http.StatusNoContent)
			}
		}
	}
}
