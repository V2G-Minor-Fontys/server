package httpx

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

func ParseUUIDParam(r *http.Request, param string) (uuid.UUID, error) {
	val := chi.URLParam(r, param)
	id, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil, BadRequest(r.Context(), fmt.Sprintf("Invalid UUID: %s", param))
	}
	return id, nil
}
