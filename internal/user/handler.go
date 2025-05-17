package user

import (
	"encoding/json"
	"net/http"

	"github.com/V2G-Minor-Fontys/server/internal/config"
	"github.com/V2G-Minor-Fontys/server/internal/httpx"
	"github.com/V2G-Minor-Fontys/server/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	svc *Service
}

func NewHandler(cfg *config.Jwt, db *pgxpool.Pool, queries *repository.Queries) *Handler {
	return &Handler{
		svc: NewService(cfg, db, queries),
	}
}

func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	idParam := chi.URLParam(r, "userId")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return httpx.BadRequest(ctx, "Invalid id")
	}

	user, err := h.svc.queries.GetUserById(r.Context(), id)
	if err != nil {
		return httpx.BadRequest(ctx, "Could not find user")
	}

	return json.NewEncoder(w).Encode(user)
}

func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	idParam := chi.URLParam(r, "userId")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return httpx.BadRequest(ctx, "Invalid id")
	}

	err = h.svc.queries.DeleteIdentityById(r.Context(), id)
	if err != nil {
		return httpx.BadRequest(ctx, "Could not find user")
	}

	return nil
}
