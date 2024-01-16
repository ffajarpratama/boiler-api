package repository

import (
	"context"

	"github.com/ffajarpratama/boiler-api/internal/model"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../mock/mock_repo.go -package=mock_repo -source=repo_interface.go
type IFaceRepository interface {
	// example
	Ping(ctx context.Context) (string, error)

	// user
	CreateUser(ctx context.Context, data *model.User, db *gorm.DB) error
	FindOneUser(ctx context.Context, query ...interface{}) (*model.User, error)
}
