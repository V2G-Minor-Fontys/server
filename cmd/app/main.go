package main

import (
	"context"
	"github.com/V2G-Minor-Fontys/server/internal/config"
	"github.com/V2G-Minor-Fontys/server/internal/repository"
	"github.com/V2G-Minor-Fontys/server/internal/router"
	"github.com/V2G-Minor-Fontys/server/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	logger.Init(cfg.Server.Environment, cfg.Logger)
	ctx, serverStopCtx := context.WithCancel(context.Background())

	conn, err := pgxpool.New(ctx, cfg.Database.Uri)
	if err != nil {
		panic(err)
	}

	repo := repository.New(conn)
	srv := router.NewServer(cfg, conn, repo)
	if err = srv.MountHandlers(); err != nil {
		panic(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		slog.InfoContext(ctx, "Shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.ErrorContext(ctx, "Error during shutdown", "error", err)
		}
		serverStopCtx()
	}()

	if err = srv.ListenAndServe(); err != nil {
		slog.ErrorContext(ctx, "Server error", "error", err)
	}

	<-ctx.Done()
	slog.InfoContext(ctx, "Server stopped")
}
