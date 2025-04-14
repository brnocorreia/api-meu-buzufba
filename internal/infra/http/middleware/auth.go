package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/brnocorreia/api-meu-buzufba/internal/infra/http/token"
	"github.com/brnocorreia/api-meu-buzufba/pkg/fault"
	"github.com/brnocorreia/api-meu-buzufba/pkg/logging"
	"go.uber.org/zap"
)

const (
	authJourney = "auth middleware"
)

type AuthKey struct{}

type middleware struct {
	secretKey string
}

func NewWithAuth(secretKey string) *middleware {
	return &middleware{
		secretKey: secretKey,
	}
}

func (m *middleware) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")

		if len(accessToken) == 0 {
			logging.Info("access token not provided",
				zap.String("journey", authJourney),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path))
			fault.NewHTTPError(w, fault.NewUnauthorized("access token not provided"))
			return
		}

		claims, err := token.Verify(m.secretKey, accessToken)
		if err != nil {
			if strings.Contains(err.Error(), "token has expired") {
				logging.Error("token has expired", err,
					zap.String("journey", authJourney),
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path))
				fault.NewHTTPError(w, fault.NewUnauthorized("token has expired"))
				return
			}
			logging.Error("invalid access token", err,
				zap.String("journey", authJourney),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path))
			fault.NewHTTPError(w, fault.NewUnauthorized("invalid access token"))
			return
		}

		ctx := context.WithValue(r.Context(), AuthKey{}, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
