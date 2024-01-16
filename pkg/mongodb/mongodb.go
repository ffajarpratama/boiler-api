package mongodb

import (
	"context"
	"log"

	"github.com/ffajarpratama/boiler-api/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDBClient(conf *config.Config) (*mongo.Database, error) {
	credential := options.Credential{
		Username: conf.MongoDB.User,
		Password: conf.MongoDB.Password,
	}

	opts := options.Client().ApplyURI(conf.MongoDB.URI).SetAuth(credential)
	conn, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal("[mongodb-connection-error] \n", err.Error())
		return nil, err
	}

	if err = conn.Ping(context.TODO(), nil); err != nil {
		log.Fatal("[mongodb-connection-error] \n", err.Error())
		return nil, err
	}

	return conn.Database(conf.MongoDB.DBName), nil
}
