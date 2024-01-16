package google

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/ffajarpratama/boiler-api/config"
	"google.golang.org/api/option"
)

type IFaceFCM interface {
	SendMessage(message messaging.MulticastMessage) error
}

type FCM struct {
	client *messaging.Client
}

// SendMessage implements IFaceFCM.
func (f *FCM) SendMessage(message messaging.MulticastMessage) error {
	_, err := f.client.SendMulticast(context.Background(), &message)
	if err != nil {
		return err
	}

	return nil
}

func NewFCMClient(cnf config.Firebase) (IFaceFCM, error) {
	opts := option.WithCredentialsFile(cnf.ServiceAccountKeyPath)
	app, err := firebase.NewApp(context.Background(), nil, opts)
	if err != nil {
		log.Fatal("[firebase-app-error] \n", err.Error())
		return nil, err
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatal("[fcm-client-error] \n", err.Error())
		return nil, err
	}

	return &FCM{client: client}, nil
}
