package app

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/ffajarpratama/boiler-api/config"
	"github.com/ffajarpratama/boiler-api/internal/repository"
	http_transport "github.com/ffajarpratama/boiler-api/internal/transport/http"
	"github.com/ffajarpratama/boiler-api/internal/usecase"
	"github.com/ffajarpratama/boiler-api/pkg/aws"
	"github.com/ffajarpratama/boiler-api/pkg/google"
	custom_http "github.com/ffajarpratama/boiler-api/pkg/http"
	"github.com/ffajarpratama/boiler-api/pkg/mongodb"
	"github.com/ffajarpratama/boiler-api/pkg/postgres"
	"github.com/ffajarpratama/boiler-api/pkg/redis"
	"github.com/ffajarpratama/boiler-api/pkg/util"
)

func Exec() (err error) {
	cnf := config.New()
	rand.New(rand.NewSource(time.Now().UnixNano()))
	util.SetTimeZone("UTC")

	dbClient, err := postgres.NewDBClient(cnf)
	if err != nil {
		return err
	}

	mongoDb, err := mongodb.NewMongoDBClient(cnf)
	if err != nil {
		return err
	}

	redisClient, err := redis.NewRedisClient(cnf)
	if err != nil {
		return err
	}

	aws, err := aws.NewAWSClient(cnf.AWSConfig)
	if err != nil {
		return err
	}

	fcm, err := google.NewFCMClient(cnf.Firebase)
	if err != nil {
		return err
	}

	repo := repository.New(dbClient, mongoDb)
	uc := usecase.New(&usecase.Usecase{
		Cnf:   cnf,
		Repo:  repo,
		DB:    dbClient,
		Redis: redisClient,
		AWS:   aws,
		FCM:   fcm,
	})

	addr := flag.String("http", fmt.Sprintf(":%d", cnf.App.Port), "HTTP listen address")
	handler := http_transport.NewHTTPHandler(cnf, uc, redisClient)

	err = custom_http.NewHTTPServer(*addr, handler, cnf)
	if err != nil {
		return err
	}

	return
}
