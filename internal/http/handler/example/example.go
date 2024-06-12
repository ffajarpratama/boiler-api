package example

import (
	"net/http"

	"github.com/ffajarpratama/boiler-api/internal/http/response"
)

func (h ExampleHandler) ExamplePing(w http.ResponseWriter, r *http.Request) {
	res := map[string]interface{}{
		"message": "PONG",
	}

	response.OK(w, res)
}
