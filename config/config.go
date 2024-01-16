package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App
	JWTConfig
	Postgres
	MongoDB
	Redis
	SMTP
	AWSConfig
	Google
	Firebase
	ImgixURL string
}

type App struct {
	Name        string
	Environment string
	Host        int
	Port        int
	URL         string
}

type JWTConfig struct {
	Admin string
	User  string
}

type Postgres struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	SSLMode  string
	URI      string
}

type MongoDB struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
	URI      string
}

type Redis struct {
	Host     string
	Port     string
	Address  string
	Password string
	Database int
}

type SMTP struct {
	User      string
	Pass      string
	Host      string
	Port      int
	EmailFrom string
}

type AWSConfig struct {
	AccessKeyID     string
	AccessKeySecret string
	Bucket          string
	Region          string
	URL             string
}

type Google struct {
	ClientID string
}

type Firebase struct {
	AndroidPackageName    string
	IosBundleID           string
	ServiceAccountKeyPath string
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
			Name:        v.GetString("APP_NAME"),
			Environment: v.GetString("APP_ENV"),
			Host:        v.GetInt("APP_HOST"),
			Port:        v.GetInt("APP_PORT"),
			URL:         v.GetString("APP_URL"),
		},
		JWTConfig: JWTConfig{
			Admin: v.GetString("JWT_SECRET_ADMIN"),
			User:  v.GetString("JWT_SECRET_USER"),
		},
		Postgres: Postgres{
			Host:     v.GetString("POSTGRES_HOST"),
			Port:     v.GetString("POSTGRES_PORT"),
			Database: v.GetString("POSTGRES_DATABASE"),
			User:     v.GetString("POSTGRES_USER"),
			Password: v.GetString("POSTGRES_PASS"),
			SSLMode:  v.GetString("POSTGRES_SSL_MODE"),
			URI:      v.GetString("POSTGRES_URI"),
		},
		MongoDB: MongoDB{
			Host:     v.GetString("MONGODB_HOST"),
			Port:     v.GetString("MONGODB_PORT"),
			DBName:   v.GetString("MONGODB_DATABASE"),
			User:     v.GetString("MONGODB_USER"),
			Password: v.GetString("MONGODB_PASSWORD"),
			URI:      v.GetString("MONGODB_URI"),
		},
		Redis: Redis{
			Host:     v.GetString("REDIS_HOST"),
			Port:     v.GetString("REDIS_PORT"),
			Address:  v.GetString("REDIS_ADDRESS"),
			Password: v.GetString("REDIS_PASSWORD"),
			Database: v.GetInt("REDIS_DATABASE"),
		},
		SMTP: SMTP{
			User:      v.GetString("SMTP_USER"),
			Pass:      v.GetString("SMTP_PASSWORD"),
			Host:      v.GetString("SMTP_HOST"),
			Port:      v.GetInt("SMTP_PORT"),
			EmailFrom: v.GetString("EMAIL_FROM"),
		},
		AWSConfig: AWSConfig{
			AccessKeyID:     v.GetString("AWS_ACCESS_KEY_ID"),
			AccessKeySecret: v.GetString("AWS_SECRET_ACCESS_KEY"),
			Bucket:          v.GetString("AWS_BUCKET"),
			Region:          v.GetString("AWS_REGION"),
			URL:             v.GetString("AWS_URL"),
		},
		Google: Google{
			ClientID: v.GetString("GOOGLE_CLIENT_ID"),
		},
		Firebase: Firebase{
			AndroidPackageName:    v.GetString("ANDROID_PACKAGE_NAME"),
			IosBundleID:           v.GetString("IOS_BUNDLE_ID"),
			ServiceAccountKeyPath: v.GetString("SERVICE_ACCOUNT_KEY_PATH"),
		},
		ImgixURL: v.GetString("IMGIX_URL"),
	}
}
