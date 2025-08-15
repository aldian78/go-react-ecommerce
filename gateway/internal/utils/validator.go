package utils

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"

	"buf.build/go/protovalidate"
	"gateway/protoc/proto-common"
	"google.golang.org/protobuf/proto"
)

func CheckValidation(req proto.Message) ([]*proto_common.ValidationError, error) {
	if err := protovalidate.Validate(req); err != nil {
		var validationError *protovalidate.ValidationError
		if errors.As(err, &validationError) {
			var validationErrorResponse []*proto_common.ValidationError = make([]*proto_common.ValidationError, 0)
			for _, violation := range validationError.Violations {
				validationErrorResponse = append(validationErrorResponse, &proto_common.ValidationError{
					Field:   *violation.Proto.Field.Elements[0].FieldName,
					Message: *violation.Proto.Message,
				})
			}

			return validationErrorResponse, nil
		}

		return nil, err
	}

	return nil, nil
}

var validate = validator.New()

// CheckValidationRest memvalidasi struct dan mengembalikan error validasi (jika ada)
func CheckValidationRest(obj interface{}) map[string]string {
	if err := validate.Struct(obj); err != nil {
		errs := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errs[e.Field()] = fmt.Sprintf("Field '%s' failed on the '%s' tag", e.Field(), e.Tag())
		}
		return errs
	}
	return nil
}
