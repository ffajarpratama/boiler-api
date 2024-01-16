package user

import (
	"net/http"

	"github.com/ffajarpratama/boiler-api/internal/middleware"
	"github.com/ffajarpratama/boiler-api/internal/usecase"
	custom_validator "github.com/ffajarpratama/boiler-api/pkg/validator"
	"github.com/go-chi/chi/v5"
)

type userHandler struct {
	uc usecase.IFaceUsecase
	v  custom_validator.Validator
	m  middleware.Middleware
}

func NewUserHandler(uc usecase.IFaceUsecase, v custom_validator.Validator, m middleware.Middleware, r chi.Router) http.Handler {
	h := userHandler{
		uc: uc,
		v:  v,
		m:  m,
	}

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Route("/profile", func(r chi.Router) {
			r.Use(m.AuthenticateUser())
			r.Get("/", h.GetProfile)
		})
	})

	r.Route("/ping", func(r chi.Router) {
		r.Get("/", h.Ping)
	})

	return r
}
