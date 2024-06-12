package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ffajarpratama/boiler-api/config"
	"github.com/ffajarpratama/boiler-api/constant"
	"github.com/ffajarpratama/boiler-api/internal/http/handler/example"
	"github.com/ffajarpratama/boiler-api/internal/http/handler/user"
	"github.com/ffajarpratama/boiler-api/internal/http/middleware"
	"github.com/ffajarpratama/boiler-api/internal/http/response"
	"github.com/ffajarpratama/boiler-api/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func NewHTTPRouter(cnf *config.Config, uc usecase.IFaceUsecase) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recover)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response.JsonResponse{
			Error: &response.ErrorResponse{
				Code:    constant.DefaultNotFoundError,
				Status:  http.StatusNotFound,
				Message: "please check url",
			},
		})
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response.JsonResponse{
			Error: &response.ErrorResponse{
				Code:    constant.DefaultMethodNotAllowedError,
				Status:  http.StatusMethodNotAllowed,
				Message: "method not allowed",
			},
		})
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response.JsonResponse{
			Success: true,
			Data: map[string]interface{}{
				"app-name": "elemes go-boilerplate",
			},
		})
	})

	r.Route("/api", func(api chi.Router) {
		api.Route("/v1", func(v1 chi.Router) {
			v1.Mount("/user", user.NewV1Handler(cnf, uc))
			v1.Mount("/example", example.NewV1Handler(uc))
		})
	})

	return r
}
