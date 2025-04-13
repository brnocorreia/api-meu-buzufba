package middleware

import (
	"net/http"

	"github.com/brnocorreia/api-meu-buzufba/pkg/metric"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Config struct {
	Metrics *metric.Metric
}
type middlewareFn func(http.Handler) http.Handler

func setup(cfg Config) []middlewareFn {
	return []middlewareFn{
		withIP,
		withRateLimit,
		withMetrics(cfg.Metrics),
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
			MaxAge:           300,
		}),
	}
}
func Apply(r *chi.Mux, cfg Config) {
	for _, midleware := range setup(cfg) {
		r.Use(midleware)
	}
}
