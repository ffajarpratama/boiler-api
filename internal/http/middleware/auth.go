package middleware

import (
	"context"
	"net/http"

	"github.com/ffajarpratama/boiler-api/constant"
	"github.com/ffajarpratama/boiler-api/internal/http/response"
	custom_jwt "github.com/ffajarpratama/boiler-api/lib/jwt"
	"github.com/ffajarpratama/boiler-api/util"
	"github.com/golang-jwt/jwt/v5"
)

func Authorize(secret string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			token, err := util.GetTokenFromHeader(r)
			if err != nil {
				response.UnauthorizedError(w)
				return
			}

			resJwt, err := jwt.ParseWithClaims(token, &custom_jwt.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil {
				response.UnauthorizedError(w)
				return
			}

			customClaims, ok := resJwt.Claims.(*custom_jwt.CustomClaims)
			if !ok && !resJwt.Valid {
				response.UnauthorizedError(w)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, constant.UserIDKey, customClaims.UserID)
			h.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
