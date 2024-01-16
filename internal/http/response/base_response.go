package response

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/ffajarpratama/boiler-api/pkg/constant"
	custom_error "github.com/ffajarpratama/boiler-api/pkg/error"
	custom_validator "github.com/ffajarpratama/boiler-api/pkg/validator"
)

const (
	CONTENT_TYPE_HEADER  = "Content-Type"
	CONTENT_DESC_HEADER  = "Content-Description"
	CONTENT_DISPO_HEADER = "Content-Disposition"

	CONTENT_TYPE_JSON         = "application/json"
	CONTENT_TYPE_PDF          = "application/pdf"
	CONTENT_TYPE_OCTET_STREAM = "aplication/octet-stream"

	CONTENT_DESC_FILE_TRANSFER = "File Transfer"
)

type JsonResponse struct {
	Success bool           `json:"success"`
	Paging  *PagingJSON    `json:"paging"`
	Data    interface{}    `json:"data"`
	Errors  *ErrorResponse `json:"errors"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type PagingJSON struct {
	Page      int   `json:"page"`
	PerPage   int   `json:"per_page"`
	Count     int64 `json:"count"`
	PageCount int   `json:"page_count"`
	Next      bool  `json:"next"`
	Prev      bool  `json:"prev"`
}

func OK(w http.ResponseWriter, data interface{}) {
	w.Header().Set(CONTENT_TYPE_HEADER, CONTENT_TYPE_JSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JsonResponse{
		Data:    data,
		Success: true,
	})
}

func Paging(w http.ResponseWriter, list interface{}, page, perPage int, cnt int64) {
	var paging *PagingJSON
	total := calculateTotalPage(cnt, perPage)
	if page > 0 {
		paging = &PagingJSON{
			Page:      page,
			PerPage:   perPage,
			Count:     cnt,
			PageCount: total,
			Next:      hasNext(page, total),
			Prev:      hasPrev(page),
		}
	}

	w.Header().Set(CONTENT_TYPE_HEADER, CONTENT_TYPE_JSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JsonResponse{
		Success: true,
		Paging:  paging,
		Data:    list,
		Errors:  nil,
	})
}

func Error(w http.ResponseWriter, err error) {
	v, isValidationErr := err.(custom_validator.ValidatorError)
	if isValidationErr {
		w.Header().Set(CONTENT_TYPE_HEADER, CONTENT_TYPE_JSON)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JsonResponse{
			Errors: &ErrorResponse{
				Code:    v.Code,
				Status:  v.Status,
				Message: v.Message,
			},
		})

		return
	}

	e, isCustomErr := err.(*custom_error.CustomErrors)
	if !isCustomErr {
		if err != nil && !errors.Is(err, context.Canceled) {
			// bugsnag.Notify(err)
			fmt.Println(err.Error(), "[unhandled-error]")
		}

		w.Header().Set(CONTENT_TYPE_HEADER, CONTENT_TYPE_JSON)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JsonResponse{
			Errors: &ErrorResponse{
				Code:    constant.DefaultUnhandledError,
				Status:  http.StatusInternalServerError,
				Message: custom_error.GetErrorMessageByErrorCode(constant.DefaultUnhandledError),
			},
		})

		return
	}

	httpCode := http.StatusInternalServerError
	internalCode := constant.DefaultUnhandledError

	if e.ErrorContext != nil && e.ErrorContext.HTTPCode > 0 {
		httpCode = e.ErrorContext.HTTPCode
		internalCode = e.ErrorContext.Code
	}

	if e.ErrorContext != nil && e.ErrorContext.Code == 0 {
		e.ErrorContext.Code = constant.DefaultUnhandledError
		internalCode = http.StatusInternalServerError
	}

	msg := e.ErrorContext.Message
	if msg == "" {
		msg = custom_error.GetErrorMessageByErrorCode(e.ErrorContext.Code)
	}

	if httpCode == http.StatusInternalServerError {
		msg = custom_error.GetErrorMessageByErrorCode(constant.DefaultUnhandledError)
	}

	w.Header().Set(CONTENT_TYPE_HEADER, CONTENT_TYPE_JSON)
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(JsonResponse{
		Errors: &ErrorResponse{
			Code:    internalCode,
			Status:  httpCode,
			Message: msg,
		},
	})
}

func UnauthorizedError(w http.ResponseWriter) {
	w.Header().Set(CONTENT_TYPE_HEADER, CONTENT_TYPE_JSON)
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(JsonResponse{
		Data:    nil,
		Success: false,
		Errors: &ErrorResponse{
			Code:    constant.DefaultUnauthorizedError,
			Status:  http.StatusUnauthorized,
			Message: "you are not authorized to access this api",
		},
	})
}

func BinaryExcel(w http.ResponseWriter, filename string, b *bytes.Buffer) {
	filename = fmt.Sprintf("%s.xlsx", filename)
	w.Header().Set(CONTENT_DESC_HEADER, CONTENT_DESC_FILE_TRANSFER)
	w.Header().Set(CONTENT_DISPO_HEADER, "attachment; filename="+filename)
	w.Header().Set(CONTENT_TYPE_HEADER, CONTENT_TYPE_OCTET_STREAM)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(b.Bytes())
	if err != nil {
		panic(err)
	}
}

func BinaryPdf(w http.ResponseWriter, filename string, b *bytes.Buffer) {
	filename = fmt.Sprintf("%s.pdf", filename)
	w.Header().Set(CONTENT_DESC_HEADER, CONTENT_DESC_FILE_TRANSFER)
	w.Header().Set(CONTENT_DISPO_HEADER, "attachment; filename="+filename)
	w.Header().Set(CONTENT_TYPE_HEADER, CONTENT_TYPE_PDF)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(b.Bytes())
	if err != nil {
		panic(err)
	}
}

func BinaryCsv(w http.ResponseWriter, filename string, b *bytes.Buffer) {
	filename = fmt.Sprintf("%s.csv", filename)
	w.Header().Set(CONTENT_DESC_HEADER, CONTENT_DESC_FILE_TRANSFER)
	w.Header().Set(CONTENT_DISPO_HEADER, "attachment; filename="+filename)
	w.Header().Set(CONTENT_TYPE_HEADER, CONTENT_TYPE_OCTET_STREAM)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(b.Bytes())
	if err != nil {
		panic(err)
	}
}

func hasNext(currentPage, totalPages int) bool {
	return currentPage < totalPages
}

func hasPrev(currentPage int) bool {
	return currentPage > 1
}

func calculateTotalPage(cnt int64, limit int) (total int) {
	return int(math.Ceil(float64(cnt) / float64(limit)))
}
