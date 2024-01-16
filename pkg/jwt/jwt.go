package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	JWT_TTL = 3600 * 24 * 1 // 1d in seconds
)

type CustomAdminClaims struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type CustomUserClaims struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(claims *CustomUserClaims, secret string) (token string, err error) {
	claims.Role = "public"
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(JWT_TTL * time.Second)),
		Issuer:    "example.com",
	}

	jwtClaims := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return jwtClaims.SignedString([]byte(secret))
}

func GenerateTokenAdmin(claims *CustomAdminClaims, secret string) (token string, err error) {
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(JWT_TTL * time.Second)),
		Issuer:    "example.com",
	}

	jwtClaims := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return jwtClaims.SignedString([]byte(secret))
}
