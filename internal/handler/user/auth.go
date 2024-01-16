package user

import (
	"net/http"

	"github.com/ffajarpratama/boiler-api/internal/http/request"
	"github.com/ffajarpratama/boiler-api/internal/http/response"
	"github.com/google/uuid"
)

func (h *userHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req request.ReqRegisterUser

	err := h.v.ValidateStruct(r, &req)
	if err != nil {
		response.Error(w, err)
		return
	}

	err = h.uc.RegisterUser(r.Context(), &req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, nil)
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req request.ReqLoginUser

	err := h.v.ValidateStruct(r, &req)
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

func (h *userHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := uuid.Parse(h.m.GetUserIDFromCtx(ctx))
	res, err := h.uc.GetUserProfile(ctx, userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, res)
}
