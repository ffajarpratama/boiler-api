package mongo

import (
	"context"
	"fmt"
	"log"

	"github.com/ffajarpratama/boiler-api/config"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(cnf *config.Config) (*mongo.Database, error) {
	credential := options.Credential{
		Username:      cnf.Mongo.Username,
		Password:      cnf.Mongo.Password,
		AuthSource:    "admin",
		AuthMechanism: "SCRAM-SHA-256",
	}

	cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			log.Print(evt.Command)
		},
	}

	uri := fmt.Sprintf("mongodb://%s:%s", cnf.Mongo.Host, cnf.Mongo.Port)
	opts := options.Client().ApplyURI(uri).SetAuth(credential).SetMonitor(cmdMonitor)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	log.Println("[mongodb-connected]")

	return client.Database(cnf.Mongo.Database), nil
}
