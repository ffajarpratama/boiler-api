package config

import (
	"log"

	"github.com/spf13/viper"
)

var GlobalConfig Config

type Config struct {
	App      App
	JWT      JWT
	Postgres Postgres
	Mongo    Mongo
}

type App struct {
	Environment string
	Port        int
	URL         string
}

type JWT struct {
	Secret string
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
	URI      string
}

type Mongo struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func New() *Config {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("[config-file-fail-load] \n", err.Error())
	}

	v := viper.GetViper()
	viper.AutomaticEnv()

	return &Config{
		App: App{
			Environment: v.GetString("APP_ENV"),
			Port:        v.GetInt("APP_PORT"),
			URL:         v.GetString("APP_URL"),
		},
		JWT: JWT{
			Secret: v.GetString("JWT_SECRET"),
		},
		Postgres: Postgres{
			Host:     v.GetString("PG_HOST"),
			Port:     v.GetString("PG_PORT"),
			User:     v.GetString("PG_USER"),
			Password: v.GetString("PG_PASSWORD"),
			Database: v.GetString("PG_DATABASE"),
			SSLMode:  v.GetString("PG_SSL_MODE"),
			URI:      v.GetString("PG_URI"),
		},
		Mongo: Mongo{
			Host:     v.GetString("MONGO_HOST"),
			Port:     v.GetString("MONGO_PORT"),
			Username: v.GetString("MONGO_USER"),
			Password: v.GetString("MONGO_PASSWORD"),
			Database: v.GetString("MONGO_DATABASE"),
		},
	}
}
