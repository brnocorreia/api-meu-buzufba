package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/brnocorreia/api-meu-buzufba/pkg/logging"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Config struct {
	Port         string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Router       *chi.Mux
}

type Server struct {
	server *http.Server
	config Config
}

func New(c Config) *Server {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", c.Port),
		IdleTimeout:  c.IdleTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		Handler:      c.Router,
	}

	return &Server{
		server: srv,
		config: c,
	}
}

func (s *Server) Start() error {
	logging.Info("server started",
		zap.String("journey", "server"),
		zap.String("port", s.config.Port))
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	logging.Info("shutting down server",
		zap.String("journey", "server"))
	return s.server.Shutdown(ctx)
}
