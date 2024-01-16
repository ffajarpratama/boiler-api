package usecase

import (
	"github.com/ffajarpratama/boiler-api/config"
	"github.com/ffajarpratama/boiler-api/internal/repository"
	"github.com/ffajarpratama/boiler-api/pkg/aws"
	"github.com/ffajarpratama/boiler-api/pkg/google"
	"github.com/ffajarpratama/boiler-api/pkg/redis"
	"gorm.io/gorm"
)

type Usecase struct {
	Cnf   *config.Config
	Repo  repository.IFaceRepository
	DB    *gorm.DB
	Redis redis.IFaceRedis
	AWS   aws.IFaceAWS
	FCM   google.IFaceFCM
}

func New(params *Usecase) IFaceUsecase {
	return &Usecase{
		Cnf:   params.Cnf,
		Repo:  params.Repo,
		DB:    params.DB,
		Redis: params.Redis,
		AWS:   params.AWS,
		FCM:   params.FCM,
	}
}
