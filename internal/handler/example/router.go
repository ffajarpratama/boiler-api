package example

import (
	"net/http"

	"github.com/ffajarpratama/boiler-api/internal/middleware"
	"github.com/ffajarpratama/boiler-api/internal/usecase"
	custom_validator "github.com/ffajarpratama/boiler-api/pkg/validator"
	"github.com/go-chi/chi/v5"
)

type exampleHandler struct {
	uc usecase.IFaceUsecase
	v  custom_validator.Validator
	m  middleware.Middleware
}

func NewExampleHandler(uc usecase.IFaceUsecase, v custom_validator.Validator, m middleware.Middleware, r chi.Router) http.Handler {
	h := exampleHandler{
		uc: uc,
		v:  v,
		m:  m,
	}

	r.Get("/foo", h.Foo)
	r.Get("/bar", h.Bar)
	r.Get("/baz", h.Baz)

	return r
}
