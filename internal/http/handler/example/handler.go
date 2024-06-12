package example

import (
	"net/http"

	"github.com/ffajarpratama/boiler-api/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type ExampleHandler struct {
	uc usecase.IFaceUsecase
}

func NewV1Handler(uc usecase.IFaceUsecase) http.Handler {
	h := ExampleHandler{uc: uc}
	r := chi.NewRouter()

	r.Route("/ping", func(r chi.Router) {
		r.Get("/", h.ExamplePing)
	})

	return r
}
