package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimid "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type middlewareFn func(http.Handler) http.Handler

func setup() []middlewareFn {
	return []middlewareFn{
		chimid.Logger,
		withIP,
		withRateLimit,
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
			MaxAge:           300,
		}),
	}
}
func Apply(r *chi.Mux) {
	for _, midleware := range setup() {
		r.Use(midleware)
	}
}
