package usecase

import (
	"context"

	"github.com/ffajarpratama/boiler-api/internal/http/request"
	"github.com/ffajarpratama/boiler-api/internal/model"
	"github.com/google/uuid"
)

type IFaceUsecase interface {
	// auth
	Register(ctx context.Context, req *request.Register) (*model.User, error)
	Login(ctx context.Context, req *request.Login) (*model.User, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*model.User, error)
}
