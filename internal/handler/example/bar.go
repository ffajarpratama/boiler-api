package example

import (
	"net/http"

	"github.com/ffajarpratama/boiler-api/internal/http/response"
)

func (h *exampleHandler) Bar(w http.ResponseWriter, r *http.Request) {
	res := map[string]interface{}{
		"message": "coming from bar...",
	}

	response.OK(w, res)
}
