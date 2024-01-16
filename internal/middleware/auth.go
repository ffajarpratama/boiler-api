package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/ffajarpratama/boiler-api/config"
	"github.com/ffajarpratama/boiler-api/internal/http/response"
	"github.com/ffajarpratama/boiler-api/pkg/constant"
	custom_jwt "github.com/ffajarpratama/boiler-api/pkg/jwt"
	"github.com/ffajarpratama/boiler-api/pkg/redis"
	"github.com/golang-jwt/jwt/v5"
)

type Middleware struct {
	Redis redis.IFaceRedis
	config.JWTConfig
}

func (m Middleware) AuthenticateUser() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			token, err := GetTokenFromHeader(r)
			if err != nil {
				response.UnauthorizedError(w)
				return
			}

			resJwt, err := jwt.ParseWithClaims(token, &custom_jwt.CustomUserClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(m.JWTConfig.User), nil
			})

			if err != nil {
				response.UnauthorizedError(w)
				return
			}

			customClaims, ok := resJwt.Claims.(*custom_jwt.CustomUserClaims)
			if !ok && !resJwt.Valid {
				response.UnauthorizedError(w)
				return
			}

			// _, err = m.Redis.Get(fmt.Sprintf("auth:%s", customClaims.SessionID))
			// if err != nil {
			// 	response.UnauthorizedError(w)
			// 	return
			// }

			ctx := r.Context()
			ctx = context.WithValue(ctx, constant.UserIDKey, customClaims.ID)
			ctx = context.WithValue(ctx, constant.RoleKey, customClaims.Role)
			h.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func (m Middleware) AuthenticateAdmin() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			token, err := GetTokenFromHeader(r)
			if err != nil {
				response.UnauthorizedError(w)
				return
			}

			resJwt, err := jwt.ParseWithClaims(token, &custom_jwt.CustomUserClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(m.JWTConfig.Admin), nil
			})

			if err != nil {
				response.UnauthorizedError(w)
				return
			}

			customClaims, ok := resJwt.Claims.(*custom_jwt.CustomUserClaims)
			if !ok && !resJwt.Valid {
				response.UnauthorizedError(w)
				return
			}

			// _, err = m.Redis.Get(fmt.Sprintf("auth:%s", customClaims.SessionID))
			// if err != nil {
			// 	response.UnauthorizedError(w)
			// 	return
			// }

			ctx := r.Context()
			ctx = context.WithValue(ctx, constant.UserIDKey, customClaims.ID)
			ctx = context.WithValue(ctx, constant.RoleKey, customClaims.Role)
			h.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func (m Middleware) OptionalAuth(cnf *config.Config) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			token, err := GetTokenFromHeader(r)
			if err != nil {
				ctx := r.Context()
				ctx = context.WithValue(ctx, constant.UserIDKey, custom_jwt.CustomUserClaims{ID: ""})
				h.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			resJwt, err := jwt.ParseWithClaims(token, &custom_jwt.CustomUserClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(cnf.JWTConfig.Admin), nil
			})

			if err != nil {
				resJwt, err = jwt.ParseWithClaims(token, &custom_jwt.CustomUserClaims{}, func(t *jwt.Token) (interface{}, error) {
					return []byte(cnf.JWTConfig.User), nil
				})

				if err != nil {
					response.UnauthorizedError(w)
					return
				}
			}

			customClaims, ok := resJwt.Claims.(*custom_jwt.CustomUserClaims)
			if !ok && !resJwt.Valid {
				response.UnauthorizedError(w)
				return
			}

			// _, err = m.Redis.Get(fmt.Sprintf("auth:%s", customClaims.SessionID))
			// if err != nil {
			// 	response.UnauthorizedError(w)
			// 	return
			// }

			ctx := r.Context()
			ctx = context.WithValue(ctx, constant.UserIDKey, customClaims.ID)
			ctx = context.WithValue(ctx, constant.RoleKey, customClaims.Role)
			h.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

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

func (m Middleware) GetUserIDFromCtx(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if userID, ok := ctx.Value(constant.UserIDKey).(string); ok {
		return userID
	}

	return ""
}

func (m Middleware) GetRole(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if roleKey, ok := ctx.Value(constant.RoleKey).(string); ok {
		return roleKey
	}

	return ""
}
