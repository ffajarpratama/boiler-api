package user

import (
	"net/http"

	"github.com/ffajarpratama/boiler-api/config"
	"github.com/ffajarpratama/boiler-api/internal/http/middleware"
	"github.com/ffajarpratama/boiler-api/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	uc usecase.IFaceUsecase
}

func NewV1Handler(cnf *config.Config, uc usecase.IFaceUsecase) http.Handler {
	h := UserHandler{uc: uc}
	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)

		r.Route("/profile", func(r chi.Router) {
			r.Use(middleware.Authorize(cnf.JWT.Secret))
			r.Get("/", h.GetProfile)
		})
	})

	return r
}
