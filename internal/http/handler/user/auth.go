package user

import (
	"net/http"

	"github.com/ffajarpratama/boiler-api/internal/http/request"
	"github.com/ffajarpratama/boiler-api/internal/http/response"
	"github.com/ffajarpratama/boiler-api/lib/custom_validator"
	"github.com/ffajarpratama/boiler-api/util"
	"github.com/google/uuid"
)

func (h UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req request.Register
	err := custom_validator.ValidateStruct(r, &req)
	if err != nil {
		response.Error(w, err)
		return
	}

	res, err := h.uc.Register(r.Context(), &req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, res)
}

func (h UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req request.Login
	err := custom_validator.ValidateStruct(r, &req)
	if err != nil {
		response.Error(w, err)
		return
	}

	res, err := h.uc.Login(r.Context(), &req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, res)
}

func (h UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := uuid.Parse(util.GetUserIDFromContext(ctx))
	res, err := h.uc.GetProfile(ctx, userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, res)
}
