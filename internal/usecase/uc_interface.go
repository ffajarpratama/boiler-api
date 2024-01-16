package usecase

import (
	"context"

	"github.com/ffajarpratama/boiler-api/internal/http/request"
	"github.com/ffajarpratama/boiler-api/internal/model"
	"github.com/google/uuid"
)

type IFaceUsecase interface {
	// example
	Ping(ctx context.Context) (string, error)

	// auth
	RegisterUser(ctx context.Context, req *request.ReqRegisterUser) error
	Login(ctx context.Context, req *request.ReqLoginUser) (*model.User, error)
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*model.User, error)
}
