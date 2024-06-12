package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Repository struct {
	BaseRepository
	db    *gorm.DB
	mongo *mongo.Database
}

func New(db *gorm.DB, mongo *mongo.Database) IFaceRepository {
	return &Repository{
		db:    db,
		mongo: mongo,
	}
}
