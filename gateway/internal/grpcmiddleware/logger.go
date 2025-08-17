package grpcmiddleware

import (
	"context"
	"encoding/json"
	"github.com/aldian78/go-react-ecommerce/gateway/internal/constants"
	"github.com/aldian78/go-react-ecommerce/gateway/internal/entity"
	"github.com/aldian78/go-react-ecommerce/gateway/internal/model"
	protoApi "github.com/aldian78/go-react-ecommerce/proto/pb/api"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/logger"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var (
	loggedContentType = []string{
		"application/json",
		"application/xml",
		"text/plain",
	}
)

// Logger .
func Logger(api *entity.ApiGateway) gin.HandlerFunc {

	//Generator ID unik
	node, _ := snowflake.NewNode(1)

	return func(c *gin.Context) {

		/* get headers */
		headers := make(map[string]string)

		for k := range c.Request.Header {
			headers[k] = c.Request.Header.Get(k)
		}

		params := make(map[string]string)

		/* get path parameters */
		for _, v := range c.Params {
			params[v.Key] = v.Value
		}

		if c.Request.Method == "GET" {

			/* get query parameters */
			query := c.Request.URL.Query()
			for k := range query {
				params[k] = c.Query(k)
			}

		} else {

			logger.Info("Content-Type: " + c.ContentType())

			/* get body parameters by content type */
			if strings.Contains(c.ContentType(), "application/x-www-form-urlencoded") ||
				strings.Contains(c.ContentType(), "multipart/form-data") {

				c.Request.ParseForm()
				for k := range c.Request.Form {
					params[k] = c.PostForm(k)
				}

			} else if strings.Contains(c.ContentType(), "application/json") {

				body := make(map[string]interface{})

				decoder := json.NewDecoder(c.Request.Body)
				err := decoder.Decode(&body)
				if err != nil {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}

				for k, v := range body {
					val := reflect.ValueOf(v)
					switch val.Kind() {
					case reflect.Map:
						bjson, _ := json.Marshal(val.Interface())
						params[k] = string(bjson)
					case reflect.Slice:
						bjson, _ := json.Marshal(val.Interface())
						params[k] = string(bjson)
					case reflect.Float64:
						params[k] = strconv.FormatFloat(val.Float(), 'f', -1, 64)
					case reflect.Bool:
						params[k] = strconv.FormatBool(val.Bool())
					default:
						params[k] = val.String()
					}
				}
			}
		}

		ID := node.Generate()
		IP := c.ClientIP()
		if net.IP(c.ClientIP()).To4() != nil {
			IP = net.IP(c.ClientIP()).To4().String()
		}

		apiReq := protoApi.APIREQ{
			TxType:  api.TxType,
			Headers: headers,
			Params:  params,
		}

		apiReq.Headers[constants.HeaderRequestID] = ID.String()
		apiReq.Headers[constants.HeaderRequestIP] = IP
		apiReq.Headers[constants.HeaderRequestFunc] = api.EndpointFunc
		apiReq.Headers[constants.HeaderRequestDesc] = api.Description

		// decrypt request for api-version 2
		//result, errs := encryptor.DecryptRequest(c, &apiReq)
		//if errs != nil {
		//	c.AbortWithStatus(http.StatusBadRequest)
		//	return
		//}

		logParam, _ := json.Marshal(&apiReq)

		logger.Infof("[%s] - [%s] request: %v", IP, ID.String(), string(logParam))

		c.Set(constants.APIReqKey, &apiReq)

		//blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		var blw *model.ResponseBuffer
		blw = model.NewResponseBuffer(c.Writer)
		c.Writer = blw

		c.Next()

		//if !c.IsAborted() {
		//	cType := c.Writer.Header().Get(constants.ContentType)
		//if isLoggedContentType(cType) {
		//	if c.GetHeader(constants.HeaderApiVersion) != encryptor.ApiVersion {
		//		logger.Infof("[%s] - [%s] response: %v", IP, ID.String(), blw.Body.String())
		//	}
		//} else {
		//	logger.Infof("[%s] - [%s] response: not printable(%v)", IP, ID.String(), cType)
		//}
		//} else {
		//insert to logging for error aborted
		body := blw.Body.String()
		if body == "" {
			body = `{"code":"UE","message":"Undefined Error"}`
		}

		//translator.TranslatePlugin(result.Headers[constant.HeaderAcceptLanguage], blw)
		logger.Infof("[%s] - [%s] response : %v", IP, ID.String(), body)

		//_, err := createAndSendRequestLogging(context.Background(), api, &apiReq, body)
		//if err != nil {
		//	logger.Errorf("error createAndSendRequestLogging : %v", err.Error())
		//}

		//}
		blw.Flush()
	}
}

func createAndSendRequestLogging(ctx context.Context, api *entity.ApiGateway, request *protoApi.APIREQ, body string) (*protoApi.APIRES, error) {
	GoClient := NewClient()

	reqJsonParam, err := json.Marshal(request.Params)
	if err != nil {
		logger.Errorf("error marshal : %v", err.Error())
		return nil, err
	}

	reqJsonHeader, err := json.Marshal(request.Headers)
	if err != nil {
		logger.Errorf("error marshal : %v", err.Error())
		return nil, err
	}

	reqAudit := protoApi.APIREQ{
		TxType:  api.TxType,
		Headers: request.Headers,
		Params: map[string]string{
			"requestData":       string(reqJsonParam),
			"responseData":      body,
			"requestDataHeader": string(reqJsonHeader),
		},
	}

	req := GoClient.NewRequest(
		"logging-service",
		"LoggingService.WriteLog",
		&reqAudit,
	)

	response := new(protoApi.APIRES)
	if err := GoClient.Call(ctx, req, response); err != nil {
		logger.Errorf("failed to call go client : %v", err.Error())
		return nil, err
	}

	return response, nil
}

func isLoggedContentType(contentType string) bool {
	for _, v := range loggedContentType {
		if strings.Contains(contentType, v) {
			return true
		}
	}
	return false
}
