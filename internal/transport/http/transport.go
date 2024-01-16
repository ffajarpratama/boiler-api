package http

import (
	"encoding/json"
	"net/http"

	"github.com/ffajarpratama/boiler-api/config"
	"github.com/ffajarpratama/boiler-api/internal/http/response"
	"github.com/ffajarpratama/boiler-api/internal/middleware"
	v1_http_handler "github.com/ffajarpratama/boiler-api/internal/transport/http/v1"
	"github.com/ffajarpratama/boiler-api/internal/usecase"
	"github.com/ffajarpratama/boiler-api/pkg/constant"
	"github.com/ffajarpratama/boiler-api/pkg/recover"
	"github.com/ffajarpratama/boiler-api/pkg/redis"
	custom_validator "github.com/ffajarpratama/boiler-api/pkg/validator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	go_validator "github.com/go-playground/validator/v10"
	en_trans "github.com/go-playground/validator/v10/translations/en"
	"github.com/sirupsen/logrus"
	swagger "github.com/swaggo/http-swagger/v2"
)

func NewHTTPHandler(cnf *config.Config, uc usecase.IFaceUsecase, redis redis.IFaceRedis) http.Handler {
	r := chi.NewRouter()
	validator := go_validator.New()
	enTrans := en.New()
	uni := ut.New(enTrans, enTrans)
	trans, _ := uni.GetTranslator("en")

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap:    logrus.FieldMap{},
		PrettyPrint: true,
	})

	en_trans.RegisterDefaultTranslations(validator, trans)
	v := custom_validator.New(validator, trans)
	m := middleware.Middleware{
		Redis:     redis,
		JWTConfig: cnf.JWTConfig,
	}

	r.Use(middleware.CustomTraceID)
	r.Use(middleware.Logger(logger))
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
			Errors: &response.ErrorResponse{
				Code:    constant.DefaultNotFoundError,
				Status:  http.StatusNotFound,
				Message: "please check url",
			},
		})
	})

	r.Use(recover.RecoverWrap)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	if cnf.App.Environment != "production" {
		r.Get("/swagger/*", swagger.WrapHandler)
	}

	r.Route("/api/v1", func(r chi.Router) {
		v1_http_handler.New(cnf, uc, v, m, r)
	})

	return r
}
