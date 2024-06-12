package custom_validator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ffajarpratama/boiler-api/constant"
	"github.com/ffajarpratama/boiler-api/lib/custom_error"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	go_validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type ValidatorError struct {
	Code    constant.InternalCode `json:"code"`
	Status  int                   `json:"status"`
	Message string                `json:"message"`
	Details []string              `json:"details"`
}

func (e ValidatorError) Error() string {
	return e.Message
}

func ValidateStruct(r *http.Request, data interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	err = json.Unmarshal(body, data)
	if err != nil {
		fmt.Println("[error-parse-body]", err.Error())
		err = custom_error.SetCustomError(&custom_error.ErrorContext{
			HTTPCode: http.StatusUnprocessableEntity,
			Message:  "please check your body request",
		})

		return err
	}

	validate := go_validator.New()
	eng := en.New()
	uni := ut.New(eng, eng)
	trans, _ := uni.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	err = validate.Struct(data)
	if err == nil {
		return nil
	}

	var message string
	var details = make([]string, 0)
	for _, field := range err.(go_validator.ValidationErrors) {
		message = field.Translate(trans)
		details = append(details, message)
	}

	err = ValidatorError{
		Code:    constant.DefaultBadRequestError,
		Status:  http.StatusBadRequest,
		Message: message,
		Details: details,
	}

	return err
}
