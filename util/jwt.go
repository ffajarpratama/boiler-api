package util

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/ffajarpratama/boiler-api/constant"
	custom_jwt "github.com/ffajarpratama/boiler-api/lib/jwt"
	"github.com/golang-jwt/jwt/v5"
)

func GetTokenFromHeader(r *http.Request) (token string, err error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		err = errors.New("token is empty")
		return "", err
	}

	lenToken := 2
	s := strings.Split(authHeader, " ")
	if len(s) != lenToken {
		err = errors.New("token is invalid")
		return "", err
	}

	token = s[1]
	return token, nil
}

func ParseWithoutVerified(token string) *custom_jwt.CustomClaims {
	res, _, err := new(jwt.Parser).ParseUnverified(token, &custom_jwt.CustomClaims{})
	if err != nil {
		return nil
	}

	claims, ok := res.Claims.(*custom_jwt.CustomClaims)
	if ok && claims.UserID != "" {
		return claims
	}

	return nil
}

func GetUserIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if userID, ok := ctx.Value(constant.UserIDKey).(string); ok {
		return userID
	}

	return ""
}
