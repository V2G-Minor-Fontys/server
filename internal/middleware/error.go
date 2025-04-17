package middleware

import (
	"errors"
	"github.com/V2G-Minor-Fontys/server/internal/httpx"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

type ErrHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func ErrHandler(h ErrHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var errResponse *httpx.Problem
		rw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		err := h(rw, r)
		if err != nil {
			if !errors.As(err, &errResponse) {
				errResponse = httpx.InternalErr(r.Context(), "An unexpected error occurred on the server while processing your request. Please try again later.", err)
			}

			if err = errResponse.Unwrap(); err != nil {
				slog.Error("Error occurred while unwrapping response",
					slog.String("error", err.Error()),
					slog.String("request.id", middleware.GetReqID(r.Context())),
					slog.String("detail", errResponse.Detail),
					slog.String("instance", errResponse.Instance),
				)
			}

			httpx.ProblemResponseWithJSON(w, errResponse)
		}
	}
}
