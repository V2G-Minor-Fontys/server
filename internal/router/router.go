package router

import (
	"context"
	"errors"
	"fmt"
	"github.com/V2G-Minor-Fontys/server/internal/config"
	"github.com/V2G-Minor-Fontys/server/internal/handler/system"
	"github.com/V2G-Minor-Fontys/server/internal/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

type Server struct {
	*chi.Mux
	cfg        *config.Config
	httpServer *http.Server
}

func NewServer(cfg *config.Config) *Server {
	srv := &Server{
		cfg: cfg,
	}
	srv.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: srv,
	}

	return srv
}

func (s *Server) MountHandlers() error {
	r := chi.NewRouter()
	r.Use(
		chiMiddleware.RequestID,
		chiMiddleware.Recoverer,
		middleware.Logger,
	)

	r.Get("/api/healthz", middleware.ErrHandler(system.HealthHandler))

	s.Mux = r
	return nil
}

func (s *Server) ListenAndServe() error {
	go func() {
		slog.Info(fmt.Sprintf("HTTP server is listening on %d", s.cfg.Server.Port))
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
