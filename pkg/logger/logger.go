package logger

import (
	"github.com/V2G-Minor-Fontys/server/internal/config"
	"log/slog"
	"os"
	"strings"
)

func Init(environment string, cfg *config.Logger) {
	var handler slog.Handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: cfg.AddSource,
		Level:     mapLogLevel(cfg.Level),
	})

	if strings.ToLower(environment) == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: cfg.AddSource,
			Level:     mapLogLevel(cfg.Level),
		})
	}

	slog.SetDefault(slog.New(handler))
}

func mapLogLevel(levelStr string) slog.Level {
	switch levelStr = strings.ToLower(levelStr); levelStr {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		slog.Info("Unsupported log level. Supported values: DEBUG, INFO, WARN, ERROR", "level", levelStr)
		return slog.LevelInfo
	}
}
