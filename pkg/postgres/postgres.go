package postgres

import (
	"log"

	"github.com/ffajarpratama/boiler-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDBClient(cnf *config.Config) (*gorm.DB, error) {
	logLevel := logger.Error
	switch cnf.App.Environment {
	case "production":
		logLevel = logger.Error
	case "development", "staging":
		logLevel = logger.Warn
	default:
		logLevel = logger.Info
	}

	gormConfig := gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	conn, err := gorm.Open(postgres.Open(cnf.Postgres.URI), &gormConfig)
	if err != nil {
		log.Fatal("[postgres-connection-error] \n", err.Error())
		return nil, err
	}

	return conn, nil
}
