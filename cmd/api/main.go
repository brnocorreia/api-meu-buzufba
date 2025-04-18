package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/brnocorreia/api-meu-buzufba/internal/config"
	"github.com/brnocorreia/api-meu-buzufba/internal/infra/database/pg"
	"github.com/brnocorreia/api-meu-buzufba/internal/infra/database/redis"
	"github.com/brnocorreia/api-meu-buzufba/internal/infra/http/middleware"
	"github.com/brnocorreia/api-meu-buzufba/internal/infra/mail"
	"github.com/brnocorreia/api-meu-buzufba/internal/infra/server"
	"github.com/brnocorreia/api-meu-buzufba/internal/modules/auth"
	"github.com/brnocorreia/api-meu-buzufba/internal/modules/session"
	"github.com/brnocorreia/api-meu-buzufba/internal/modules/stop"
	"github.com/brnocorreia/api-meu-buzufba/internal/modules/user"
	"github.com/brnocorreia/api-meu-buzufba/pkg/cache"
	"github.com/brnocorreia/api-meu-buzufba/pkg/logging"
	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
)

func main() {
	ctx := context.Background()
	cfg := config.GetConfig()

	r := chi.NewRouter()
	middleware.Apply(r)

	redisConn, err := redis.NewConnection(ctx, cfg)
	if err != nil {
		logging.Error("failed to connect to redis", err, zap.String("journey", "main"))
		panic(err)
	}
	defer redisConn.Close()

	cache, err := cache.New(ctx, redisConn.DB())
	if err != nil {
		logging.Error("failed to connect to cache", err, zap.String("journey", "main"))
		panic(err)
	}
	defer cache.Close()

	pgConn, err := pg.NewConnection(ctx, cfg.PostgresDSN)
	if err != nil {
		logging.Error("failed to connect to database", err, zap.String("journey", "main"))
		panic(err)
	}
	// Migrating database
	err = pgConn.Migrate()
	if err != nil {
		logging.Error("failed to migrate database", err, zap.String("journey", "main"))
		panic(err)
	}
	defer pgConn.Close()

	// Repositories
	userRepo := user.NewRepo(pgConn.DB())
	sessionRepo := session.NewRepo(pgConn.DB())
	stopRepo := stop.NewRepo(pgConn.DB())

	// Services
	mailService := mail.New(ctx, mail.Config{
		MaxRetries: 3,
		APIKey:     cfg.ResendKey,
		RetryDelay: time.Second * 2,
		Timeout:    time.Second * 5,
	})
	// userService := user.NewService(log, userRepo)
	sessionService := session.NewService(session.ServiceConfig{
		SessionRepo: sessionRepo,
		UserRepo:    userRepo,
		Cache:       cache,
		SecretKey:   cfg.JWTSecretKey,
	})
	authService := auth.NewService(auth.ServiceConfig{
		UserRepo:       userRepo,
		SessionService: sessionService,
		SessionRepo:    sessionRepo,
		Mailer:         mailService,
		Cache:          cache,
		SecretKey:      cfg.JWTSecretKey,
	})
	stopService := stop.NewService(stop.ServiceConfig{
		StopRepo: stopRepo,
	})

	// Handlers
	session.NewHandler(sessionService, cfg.JWTSecretKey).Register(r)
	auth.NewHandler(authService, cfg.JWTSecretKey).Register(r)
	stop.NewHandler(stopService).Register(r)

	srv := server.New(server.Config{
		Port:         cfg.Port,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
		Router:       r,
	})

	shutdoewnErr := srv.GracefulShutdown(ctx, time.Second*30)

	err = srv.Start()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logging.Error("failed to start server", err, zap.String("journey", "main"))
		os.Exit(1)
	}

	err = <-shutdoewnErr
	if err != nil {
		logging.Error("failed to shutdown server", err, zap.String("journey", "main"))
		os.Exit(1)
	}

	logging.Info("server shutdown gracefully", zap.String("journey", "main"))
}
