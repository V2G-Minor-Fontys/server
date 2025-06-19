package router

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/V2G-Minor-Fontys/server/internal/auth"
	"github.com/V2G-Minor-Fontys/server/internal/charging_preferences"
	"github.com/V2G-Minor-Fontys/server/internal/config"
	"github.com/V2G-Minor-Fontys/server/internal/middleware"
	"github.com/V2G-Minor-Fontys/server/internal/repository"
	"github.com/V2G-Minor-Fontys/server/internal/system"
	"github.com/V2G-Minor-Fontys/server/internal/user"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	*chi.Mux
	cfg                 *config.Config
	httpServer          *http.Server
	auth                *auth.Handler
	user                *user.Handler
	chargingPreferences *charging_preferences.Handler
}

func NewServer(cfg *config.Config, pool *pgxpool.Pool, queries *repository.Queries) *Server {
	srv := &Server{
		cfg:                 cfg,
		auth:                auth.NewHandler(cfg.Jwt, pool, queries),
		user:                user.NewHandler(cfg.Jwt, pool, queries),
		chargingPreferences: charging_preferences.NewHandler(cfg.Jwt, pool, queries),
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

	r.Route("/api", func(r chi.Router) {
		r.Get("/healthz", middleware.ErrHandler(system.HealthHandler))
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", middleware.ErrHandler(s.auth.RegisterHandler))
			r.Post("/login", middleware.ErrHandler(s.auth.LoginHandler))

			r.Route("/token", func(r chi.Router) {
				r.Post("/refresh", middleware.ErrHandler(s.auth.RefreshTokenHandler))
				r.Delete("/revoke", middleware.ErrHandler(s.auth.RevokeTokenHandler))
			})
		})
		r.Route("/user", func(r chi.Router) {
			r.Get("/{id}", middleware.ErrHandler(s.user.GetUserHandler))
			r.Delete("/{id}", middleware.ErrHandler(s.user.DeleteUserHandler))
		})
		r.Route("/charging-preferences", func(r chi.Router) {
			r.Get("/user/{user-id}", middleware.ErrHandler(s.chargingPreferences.GetChargingPreferencesOfUserHandler))
			r.Get("/schema/{user-id}", middleware.ErrHandler(s.chargingPreferences.CreateChargingPreferencesSchemaHandler)) // move to mqtt component
			r.Post("/", middleware.ErrHandler(s.chargingPreferences.CreateChargingPreferenceHandler))
			r.Patch("/{id}", middleware.ErrHandler(s.chargingPreferences.GetChargingPreferencesOfUserHandler))
			r.Delete("/{id}", middleware.ErrHandler(s.chargingPreferences.GetChargingPreferencesOfUserHandler))
		})
	})

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
