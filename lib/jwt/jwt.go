package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID string
	jwt.RegisteredClaims
}

const (
	ACCESS_TTL = time.Duration(1 * 24 * time.Hour)
)

func GenerateToken(claims *CustomClaims, secret string) (string, error) {
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ACCESS_TTL)),
		Issuer:    "elemes.id",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(secret))
}
