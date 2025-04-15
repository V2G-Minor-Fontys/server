package middleware

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		rw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		requestId := middleware.GetReqID(r.Context())

		slog.InfoContext(ctx, "request started",
			slog.String("request.id", requestId),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path))

		start := time.Now()
		next.ServeHTTP(rw, r)
		stop := time.Since(start)

		slog.InfoContext(ctx, "request completed",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rw.Status()),
			slog.String("status.text", http.StatusText(rw.Status())),
			slog.Duration("duration", stop),
			slog.String("request.id", requestId))

	})
}
