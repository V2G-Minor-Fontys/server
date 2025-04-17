package httpx

import (
	"context"
	"encoding/json"
	"fmt"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

const (
	UnauthorizedType = "https://datatracker.ietf.org/doc/html/rfc7235#section-3.1"
	NotFoundType     = "https://datatracker.ietf.org/doc/html/rfc7231#section-6.5.4"
	ConflictType     = "https://datatracker.ietf.org/doc/html/rfc7231#section-6.5.8"
	BadRequestType   = "https://datatracker.ietf.org/doc/html/rfc7231#section-6.5.1"
	InternalType     = "https://datatracker.ietf.org/doc/html/rfc7231#section-6.6.1"
)

type Problem struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
	err      error
}

func (p *Problem) Unwrap() error {
	return p.err
}

func (p *Problem) Error() string {
	if p.err != nil {
		return fmt.Sprintf("%s: %v", p.Detail, p.err)
	}

	return p.Detail
}

func ProblemResponseWithJSON(w http.ResponseWriter, problem *Problem) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(problem.Status)
	if err := json.NewEncoder(w).Encode(problem); err != nil {
		slog.Error("Could not encode problem response",
			slog.String("error", err.Error()),
			slog.String("problem.instance", problem.Instance),
			slog.String("problem.error", problem.err.Error()))
	}
}

func newProblem(ctx context.Context, status int, title, detail, typ string, err error) *Problem {
	return &Problem{
		Type:     typ,
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: chiMiddleware.GetReqID(ctx),
		err:      err,
	}
}

func Unauthorized(ctx context.Context, detail string) *Problem {
	return newProblem(
		ctx,
		http.StatusUnauthorized,
		"Unauthorized",
		detail,
		UnauthorizedType,
		nil,
	)
}

func NotFound(ctx context.Context, detail string) *Problem {
	return newProblem(
		ctx,
		http.StatusNotFound,
		"Resource Not Found",
		detail,
		NotFoundType,
		nil,
	)
}

func Conflict(ctx context.Context, detail string) *Problem {
	return newProblem(
		ctx,
		http.StatusConflict,
		"Conflict",
		detail,
		ConflictType,
		nil,
	)
}

func BadRequest(ctx context.Context, detail string) *Problem {
	return newProblem(
		ctx,
		http.StatusBadRequest,
		"Bad Request",
		detail,
		BadRequestType,
		nil,
	)
}

func InternalErr(ctx context.Context, detail string, err error) *Problem {
	return newProblem(
		ctx,
		http.StatusInternalServerError,
		"Internal Server Error",
		detail,
		InternalType,
		err,
	)
}
