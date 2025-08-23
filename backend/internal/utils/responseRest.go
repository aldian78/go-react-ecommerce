package utils

import (
	"encoding/json"
	"fmt"
	cons "github.com/aldian78/go-react-ecommerce/backend/internal/constants"
	model "github.com/aldian78/go-react-ecommerce/common/rest"
	"github.com/aldian78/go-react-ecommerce/proto/pb/basecommon"
	"reflect"
	"strconv"
	"strings"
)

// ResSuccess .
func ResSuccess(resData ...interface{}) (result []byte) {
	r := model.RestResult{}
	r.Code = cons.Success
	r.Message = cons.ResSuccess
	for _, val := range resData {
		if !reflect.ValueOf(val).IsNil() {
			r.Data = append(r.Data, val)
		} else {
		}
	}
	result, _ = json.Marshal(r)
	return
}

func ResSuccessWithData(responseCode int64, message string, resData ...interface{}) (result []byte) {
	r := model.RestResult{}
	rc := strconv.Itoa(int(responseCode))
	r.Code = rc
	r.Message = message
	for _, val := range resData {
		if !reflect.ValueOf(val).IsNil() {
			r.Data = append(r.Data, val)
		}
	}
	result, _ = json.Marshal(r)
	return
}

// ResSuccessArray
func ResSuccessList(resData ...interface{}) model.RestResultSingle {
	r := model.RestResultSingle{}
	r.Code = cons.Success
	r.Message = cons.ResSuccess
	for _, val := range resData {
		if !reflect.ValueOf(val).IsNil() {
			r.Data = val
		}
	}
	return r
}

// ResError .
func Error(responseCode int64, message string) (result []byte) {
	r := model.RestResult{}
	rc := strconv.Itoa(int(responseCode))
	r.Code = rc
	r.Message = message
	result, _ = json.Marshal(r)
	return
}

// ResErrorList
func ErrorWithDataList(responseCode string, message string, resData ...interface{}) (result []byte) {
	r := model.RestResultSingle{}
	r.Code = responseCode
	r.Message = message
	for _, val := range resData {
		if !reflect.ValueOf(val).IsNil() {
			r.Data = val
		}
	}
	result, _ = json.Marshal(r)
	return
}

func ErrorWithData(responseCode string, message string, resData ...interface{}) (result []byte) {
	r := model.RestResult{}
	r.Code = responseCode
	r.Message = message
	for _, val := range resData {
		if !reflect.ValueOf(val).IsNil() {
			r.Data = append(r.Data, val)
		}
	}
	result, _ = json.Marshal(r)
	return
}

func InternalServerError() (result []byte) {
	r := model.RestResult{}
	r.Code = cons.InternalErr
	r.Message = cons.ResInternalError
	result, _ = json.Marshal(r)
	return
}

func LoopValidationError(validationErrors []*basecommon.ValidationError) string {
	var messages []string
	for _, ve := range validationErrors {
		messages = append(messages, fmt.Sprintf("%s: %s", ve.Field, ve.Message))
	}
	finalMessage := strings.Join(messages, ", ")

	return finalMessage
}
