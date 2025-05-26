package router

import (
	"context"
	"errors"
	"fmt"
	"github.com/V2G-Minor-Fontys/server/internal/controller"
	"github.com/V2G-Minor-Fontys/server/internal/mqtt"
	"log/slog"
	"net/http"

	"github.com/V2G-Minor-Fontys/server/internal/auth"
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
	cfg         *config.Config
	httpServer  *http.Server
	auth        *auth.Handler
	user        *user.Handler
	controllers *controller.Handler
}

func NewServer(cfg *config.Config, mqttService *mqtt.Service, pool *pgxpool.Pool, queries *repository.Queries) *Server {
	srv := &Server{
		cfg:         cfg,
		auth:        auth.NewHandler(cfg.Jwt, pool, queries),
		user:        user.NewHandler(pool, queries),
		controllers: controller.NewHandler(mqttService, pool, queries),
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
		r.With(middleware.AuthVerifier(s.cfg.Jwt)).Route("/controllers", func(r chi.Router) {
			r.Post("/register", middleware.ErrHandler(s.controllers.RegisterControllerHandler))
			r.Get("/{cpuId}", middleware.ErrHandler(s.controllers.GetControllerByCpuIdHandler))
		})
		r.With(middleware.AuthVerifier(s.cfg.Jwt)).Route("/users/{userId}", func(r chi.Router) {
			r.Get("/", middleware.ErrHandler(s.user.GetUserHandler))
			r.Delete("/", middleware.ErrHandler(s.user.DeleteUserHandler))

			r.Route("/controllers", func(r chi.Router) {
				r.Get("/", middleware.ErrHandler(s.controllers.GetUserControllerHandler))
				r.Post("/", middleware.ErrHandler(s.controllers.PairUserToControllerHandler))

				r.Route("/{controllerId}", func(r chi.Router) {
					r.Put("/settings", middleware.ErrHandler(s.controllers.UpdateControllerSettingsHandler))

					r.Post("/actions", middleware.ErrHandler(s.controllers.ExecuteControllerActionHandler))
					r.Get("/history", middleware.ErrHandler(s.controllers.GetControllerTelemetryById))
				})

			})
		})
	})

	s.controllers.MountMqttMessageHandlers()

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
	s.controllers.ShutdownMQTT(2000)
	return s.httpServer.Shutdown(ctx)
}
