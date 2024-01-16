package repository

import (
	"context"

	"github.com/ffajarpratama/boiler-api/internal/model"
	"gorm.io/gorm"
)

// CreateUser implements IFaceRepository.
func (r *Repository) CreateUser(ctx context.Context, data *model.User, db *gorm.DB) error {
	return r.BaseRepository.Create(db.WithContext(ctx), data)
}

// FindOneUser implements IFaceRepository.
func (r *Repository) FindOneUser(ctx context.Context, query ...interface{}) (*model.User, error) {
	var res *model.User

	if err := r.BaseRepository.FindOne(r.db.WithContext(ctx).Where(query[0], query[1:]...), &res); err != nil {
		return nil, err
	}

	return res, nil
}
