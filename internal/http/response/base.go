package response

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/ffajarpratama/boiler-api/constant"
	"github.com/ffajarpratama/boiler-api/lib/custom_error"
	"github.com/ffajarpratama/boiler-api/lib/custom_validator"
)

type JsonResponse struct {
	Success bool            `json:"success"`
	Paging  *PagingResponse `json:"paging"`
	Data    interface{}     `json:"data"`
	Error   *ErrorResponse  `json:"error"`
}

type PagingResponse struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Count int64 `json:"count"`
	Total int   `json:"total"`
	Next  bool  `json:"next"`
	Prev  bool  `json:"prev"`
}

type ErrorResponse struct {
	Code    constant.InternalCode `json:"code"`
	Status  int                   `json:"status"`
	Message string                `json:"message"`
	Details []string              `json:"details"`
}

func OK(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JsonResponse{
		Data:    data,
		Success: true,
	})
}

func Paging(w http.ResponseWriter, data interface{}, page int, limit int, cnt int64) {
	var paging *PagingResponse

	total := calculateTotalPage(cnt, limit)
	if page > 0 {
		paging = &PagingResponse{
			Page:  page,
			Limit: limit,
			Count: cnt,
			Total: total,
			Next:  hasNext(page, total),
			Prev:  hasPrev(page),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JsonResponse{
		Success: true,
		Paging:  paging,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, err error) {
	v, isValidationErr := err.(custom_validator.ValidatorError)
	if isValidationErr {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JsonResponse{
			Error: &ErrorResponse{
				Code:    v.Code,
				Status:  v.Status,
				Message: v.Message,
				Details: v.Details,
			},
		})

		return
	}

	e, isCustomErr := err.(*custom_error.CustomError)
	if !isCustomErr {
		if err != nil && !errors.Is(err, context.Canceled) {
			fmt.Println(err.Error(), "[unhandled-error]")

			// if config.GlobalConfig.App.Environment != "local" {
			// 	// todo: notify bugsnag
			// }
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JsonResponse{
			Error: &ErrorResponse{
				Code:    constant.DefaultUnhandledError,
				Status:  http.StatusInternalServerError,
				Message: constant.HTTPStatusText(http.StatusInternalServerError),
			},
		})

		return
	}

	code := constant.DefaultUnhandledError
	httpCode := http.StatusInternalServerError
	msg := constant.HTTPStatusText(httpCode)

	if e.ErrorContext != nil && e.ErrorContext.HTTPCode > 0 {
		code = constant.InternalCodeMap[e.ErrorContext.HTTPCode]
		httpCode = e.ErrorContext.HTTPCode
		msg = constant.HTTPStatusText(httpCode)

		if e.ErrorContext.Message != "" {
			msg = e.ErrorContext.Message
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(JsonResponse{
		Error: &ErrorResponse{
			Code:    code,
			Status:  httpCode,
			Message: msg,
		},
	})
}

func UnauthorizedError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(JsonResponse{
		Data:    nil,
		Success: false,
		Error: &ErrorResponse{
			Code:    constant.DefaultUnauthorizedError,
			Status:  http.StatusUnauthorized,
			Message: constant.HTTPStatusText(http.StatusUnauthorized),
		},
	})
}

func BinaryExcel(w http.ResponseWriter, filename string, b *bytes.Buffer) {
	filename = fmt.Sprintf("%s.xlsx", filename)
	w.Header().Set("Content-Description", "File Transfer")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(b.Bytes())
	if err != nil {
		panic(err)
	}
}

func BinaryPdf(w http.ResponseWriter, filename string, b *bytes.Buffer) {
	filename = fmt.Sprintf("%s.pdf", filename)
	w.Header().Set("Content-Description", "File Transfer")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(b.Bytes())
	if err != nil {
		panic(err)
	}
}

func BinaryCsv(w http.ResponseWriter, filename string, b *bytes.Buffer) {
	filename = fmt.Sprintf("%s.csv", filename)
	w.Header().Set("Content-Description", "File Transfer")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")
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
