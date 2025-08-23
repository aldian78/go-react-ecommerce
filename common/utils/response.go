package utils

import (
	"github.com/aldian78/go-react-ecommerce/proto/pb/basecommon"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func SuccessResponse(message string) *basecommon.BaseResponse {
	return &basecommon.BaseResponse{
		StatusCode: 200,
		Message:    message,
	}
}

func BadRequestResponse(message string) *basecommon.BaseResponse {
	return &basecommon.BaseResponse{
		StatusCode: 400,
		Message:    message,
		IsError:    true,
	}
}

func NotFoundResponse(message string) *basecommon.BaseResponse {
	return &basecommon.BaseResponse{
		StatusCode: 404,
		Message:    message,
		IsError:    true,
	}
}

func UnauthenticatedResponse() error {
	return status.Error(codes.Unauthenticated, "Unauthenticated")
}

func ValidationErrorResponse(validationErrors []*basecommon.ValidationError) *basecommon.BaseResponse {
	return &basecommon.BaseResponse{
		StatusCode:       400,
		Message:          "Validation error",
		IsError:          true,
		ValidationErrors: validationErrors,
	}
}
