package httpx

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

const (
	NotFoundType   = "https://datatracker.ietf.org/doc/html/rfc7231#section-6.5.4"
	BadRequestType = "https://datatracker.ietf.org/doc/html/rfc7231#section-6.5.1"
	InternalType   = "https://datatracker.ietf.org/doc/html/rfc7231#section-6.6.1"
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

func newProblem(status int, title, detail, instance, typ string, err error) *Problem {
	return &Problem{
		Type:     typ,
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: instance,
		err:      err,
	}
}

func NotFound(detail, instance string) *Problem {
	return newProblem(
		http.StatusNotFound,
		"Resource Not Found",
		detail,
		instance,
		NotFoundType,
		nil,
	)
}

func NotFoundErr(detail, instance string, err error) *Problem {
	return newProblem(
		http.StatusNotFound,
		"Resource Not Found",
		detail,
		instance,
		NotFoundType,
		err,
	)
}

func BadRequest(detail, instance string) *Problem {
	return newProblem(
		http.StatusBadRequest,
		"Bad Request",
		detail,
		instance,
		BadRequestType,
		nil,
	)
}

func BadRequestErr(detail, instance string, err error) *Problem {
	return newProblem(
		http.StatusBadRequest,
		"Bad Request",
		detail,
		instance,
		BadRequestType,
		err,
	)
}

func Internal(detail, instance string) *Problem {
	return newProblem(
		http.StatusInternalServerError,
		"Internal Server Error",
		detail,
		instance,
		InternalType,
		nil,
	)
}

func InternalErr(detail, instance string, err error) *Problem {
	return newProblem(
		http.StatusInternalServerError,
		"Internal Server Error",
		detail,
		instance,
		InternalType,
		err,
	)
}
