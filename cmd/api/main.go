package main

import (
	"context"

	"github.com/brnocorreia/api-meu-buzufba/internal/infra/database/pg"
	"github.com/brnocorreia/api-meu-buzufba/internal/infra/http/middleware"
	"github.com/brnocorreia/api-meu-buzufba/pkg/cache"
	"github.com/brnocorreia/api-meu-buzufba/pkg/config"
	"github.com/brnocorreia/api-meu-buzufba/pkg/logging"
	"github.com/go-chi/chi/v5"
	cmid "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	ctx := context.Background()
	cfg := config.GetConfig()

	log := logging.New(logging.LogParams{
		AppName:                  cfg.AppName,
		DebugLevel:               cfg.DebugMode,
		AddAttributesFromContext: nil,
		LogToFile:                false,
	})

	r := chi.NewRouter()
	r.Use(
		cmid.Logger,
		middleware.WithIP,
		middleware.WithRateLimit,
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
			MaxAge:           300,
		}),
	)

	cache, err := cache.New(ctx, &cache.Config{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
	})
	if err != nil {
		log.Criticalw(ctx, "failed to connect to cache", logging.Err(err))
		panic(err)
	}
	defer cache.Close()

	con, err := pg.NewConnection(cfg.PostgresDSN)
	if err != nil {
		log.Criticalw(ctx, "failed to connect database", logging.Err(err))
		panic(err)
	}
	defer con.Close()
}
