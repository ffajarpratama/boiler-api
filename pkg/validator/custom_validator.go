package custom_validator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ffajarpratama/boiler-api/pkg/constant"
	custom_error "github.com/ffajarpratama/boiler-api/pkg/error"
	uni_trans "github.com/go-playground/universal-translator"
	go_validator "github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *go_validator.Validate
	trans    uni_trans.Translator
}

type ValidatorError struct {
	Code    int    `json:"code"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (o ValidatorError) Error() string {
	return "validate.request"
}

func New(v *go_validator.Validate, trans uni_trans.Translator) Validator {
	return Validator{
		validate: v,
		trans:    trans,
	}
}

func (v *Validator) ValidateStruct(r *http.Request, values interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	defer r.Body.Close()
	err = json.Unmarshal(body, values)
	if err != nil {
		fmt.Println("[error-parse-body]", err.Error())
		err = custom_error.SetCustomError(&custom_error.ErrorContext{
			Code:     constant.DefaultBadRequestError,
			HTTPCode: http.StatusUnprocessableEntity,
			Message:  "please check your body request",
		})

		return err
	}

	err = v.validate.Struct(values)
	if err == nil {
		return nil
	}

	var message string
	for _, field := range err.(go_validator.ValidationErrors) {
		message = field.Translate(v.trans)
	}

	validatorError := ValidatorError{
		Code:    constant.DefaultBadRequestError,
		Status:  http.StatusBadRequest,
		Message: message,
	}

	return validatorError
}
