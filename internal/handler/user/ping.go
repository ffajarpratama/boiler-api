package user

import (
	"net/http"

	"github.com/ffajarpratama/boiler-api/internal/http/response"
)

func (h *userHandler) Ping(w http.ResponseWriter, r *http.Request) {
	res, err := h.uc.Ping(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, map[string]interface{}{"message": res})
}
