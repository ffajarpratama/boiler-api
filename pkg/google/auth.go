package google

import (
	"context"

	"github.com/ffajarpratama/boiler-api/config"
	"google.golang.org/api/idtoken"
)

type UserData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func VerifyIdToken(cnf *config.Config, code string) (UserData, error) {
	payload, err := idtoken.Validate(context.Background(), code, cnf.Google.ClientID)
	if err != nil {
		return UserData{}, err
	}

	user := UserData{
		Email: payload.Claims["email"].(string),
		Name:  payload.Claims["name"].(string),
	}

	return user, nil
}
