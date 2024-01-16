package v1

import (
	"github.com/ffajarpratama/boiler-api/config"
	"github.com/ffajarpratama/boiler-api/internal/handler/example"
	"github.com/ffajarpratama/boiler-api/internal/handler/user"
	"github.com/ffajarpratama/boiler-api/internal/middleware"
	"github.com/ffajarpratama/boiler-api/internal/usecase"
	custom_validator "github.com/ffajarpratama/boiler-api/pkg/validator"
	"github.com/go-chi/chi/v5"
)

func New(cnf *config.Config, uc usecase.IFaceUsecase, v custom_validator.Validator, m middleware.Middleware, r chi.Router) {
	r.Mount("/example", example.NewExampleHandler(uc, v, m, r)) // api/v1/example/*
	r.Mount("/user", user.NewUserHandler(uc, v, m, r))          // api/v1/user/*
}
