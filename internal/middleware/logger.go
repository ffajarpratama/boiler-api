package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	custom_jwt "github.com/ffajarpratama/boiler-api/pkg/jwt"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type userLog struct {
	ID   string
	Role string
}

type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (r customResponseWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *customResponseWriter) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func CustomTraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get(chi_middleware.RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = context.WithValue(ctx, chi_middleware.RequestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Logger(log *logrus.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			reqID := getTraceID(r.Context())
			ww := &customResponseWriter{body: &bytes.Buffer{}, ResponseWriter: w}
			t1 := time.Now()

			bodyJSON, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("[error-io-reader] \n%v\n", err)
			}

			var payload map[string]interface{}
			if len(bodyJSON) > 0 {
				err = json.Unmarshal(bodyJSON, &payload)
				if err != nil {
					log.Printf("[error-unmarshal] \n%v\n", err)
				}
			}

			token, _ := GetTokenFromHeader(r)
			userClaims := parseWithoutVerified(token)
			defer func() {
				fields := logrus.Fields{
					"http_code": ww.statusCode,
					"duration":  fmt.Sprintf("%dms", int(time.Since(t1).Milliseconds())),
					"action": map[string]interface{}{
						"method":  r.Method,
						"request": payload,
						"path":    r.RequestURI,
					},
				}

				if len(reqID) > 0 {
					fields["trace_id"] = reqID
				}

				if token != "" && userClaims != nil {
					fields["user"] = map[string]interface{}{
						"id":   userClaims.ID,
						"role": userClaims.Role,
					}
				}

				err = r.Body.Close()
				if err != nil {
					log.Printf("[error-body-close] \n%v\n", err)
				}

				if !isExcludeRouter(r.RequestURI) {
					log.WithFields(fields).Info()
				}
			}()

			// create new body for handlers
			newBody := io.NopCloser(bytes.NewBuffer(bodyJSON))
			r.Body = newBody

			h.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}

func isExcludeRouter(path string) bool {
	if strings.Contains(path, "swagger") {
		return true
	}

	if strings.Contains(path, "upload") {
		return true
	}

	if strings.Contains(path, "import") {
		return true
	}

	excluded := make(map[string]bool)
	excluded["news-and-promotion/upload"] = true

	if _, ok := excluded[path]; ok {
		return true
	}

	return false
}

func parseWithoutVerified(tokenString string) *userLog {
	resJwt, _, err := new(jwt.Parser).ParseUnverified(tokenString, &custom_jwt.CustomUserClaims{})
	if err != nil {
		return nil
	}

	if customClaims, ok := resJwt.Claims.(*custom_jwt.CustomUserClaims); ok {
		return &userLog{
			ID:   customClaims.ID,
			Role: customClaims.Role,
		}
	}

	return nil
}

// getTraceID returns a request ID from the given context if one is present.
// Returns the empty string if a request ID cannot be found.
func getTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if reqID, ok := ctx.Value(chi_middleware.RequestIDKey).(string); ok {
		return reqID
	}

	return ""
}
