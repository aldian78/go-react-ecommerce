package utils

import (
	"gateway/protoc/proto-common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func SuccessResponse(message string) *proto_common.BaseResponse {
	return &proto_common.BaseResponse{
		StatusCode: 200,
		Message:    message,
	}
}

func BadRequestResponse(message string) *proto_common.BaseResponse {
	return &proto_common.BaseResponse{
		StatusCode: 400,
		Message:    message,
		IsError:    true,
	}
}

func NotFoundResponse(message string) *proto_common.BaseResponse {
	return &proto_common.BaseResponse{
		StatusCode: 404,
		Message:    message,
		IsError:    true,
	}
}

func UnauthenticatedResponse() error {
	return status.Error(codes.Unauthenticated, "Unauthenticated")
}

func ValidationErrorResponse(validationErrors []*proto_common.ValidationError) *proto_common.BaseResponse {
	return &proto_common.BaseResponse{
		StatusCode:       400,
		Message:          "Validation error",
		IsError:          true,
		ValidationErrors: validationErrors,
	}
}
