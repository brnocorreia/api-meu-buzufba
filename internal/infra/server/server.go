package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
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
	slog.Info("server started", "port", s.config.Port)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
