package recover

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
)

const (
	MAX_BYTE = 2084
)

func RecoverWrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, MAX_BYTE)
				n := runtime.Stack(buf, false)
				buf = buf[:n]

				fmt.Printf("[ERROR::RECOVER] \n%v\n\n", err)
				fmt.Printf("[StackTrace::] \n%s\n", buf)

				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)

				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"error": "Something went wrong, please try again later",
				})
			}
		}()

		h.ServeHTTP(w, r)
	})
}

func GetPanicErrorMsg(err interface{}) string {
	v := reflect.ValueOf(err)

	//nolint // no need all case
	switch v.Kind() {
	case reflect.String:
		return err.(string)
	default:
		return err.(error).Error()
	}
}
