package utils

import (
	"encoding/json"
	cons "gateway/internal/constants"
	model "gateway/internal/entity/rest"
	"reflect"
)

// ResSuccess .
func ResSuccess(resData ...interface{}) model.RestResult {
	r := model.RestResult{}
	r.Code = cons.Success
	r.Message = cons.ResSuccess
	for _, val := range resData {
		if !reflect.ValueOf(val).IsNil() {
			r.Data = append(r.Data, val)
		} else {
		}
	}
	return r
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
func Error(responseCode string, message string) model.RestResult {
	r := model.RestResult{}
	r.Code = responseCode
	r.Message = message
	return r
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
