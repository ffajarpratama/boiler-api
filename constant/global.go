package constant

import "net/http"

type InternalCode int
type ContextKey string

const (
	DefaultUnhandledError InternalCode = 1000 + iota
	DefaultNotFoundError
	DefaultMethodNotAllowedError
	DefaultBadRequestError
	DefaultUnauthorizedError
	DefaultDuplicateDataError
)

const (
	UserIDKey ContextKey = "user_id"
)

var InternalCodeMap = map[int]InternalCode{
	http.StatusInternalServerError: DefaultUnhandledError,
	http.StatusNotFound:            DefaultNotFoundError,
	http.StatusBadRequest:          DefaultBadRequestError,
	http.StatusUnauthorized:        DefaultUnauthorizedError,
	http.StatusConflict:            DefaultDuplicateDataError,
}

func HTTPStatusText(code int) string {
	switch code {
	case http.StatusInternalServerError:
		return "something went wrong with our side, please try again"
	case http.StatusNotFound:
		return "data not found"
	case http.StatusUnauthorized:
		return "you are not authorized to access this api"
	case http.StatusConflict:
		return "duplicate data error"
	case http.StatusUnprocessableEntity:
		return "please check your body request"
	case http.StatusBadRequest:
		return "request doesn't pass validation"
	default:
		return ""
	}
}
