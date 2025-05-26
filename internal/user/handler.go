package user

import (
	"net/http"

	"github.com/V2G-Minor-Fontys/server/internal/httpx"
	"github.com/V2G-Minor-Fontys/server/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	svc Service
}

func NewHandler(db *pgxpool.Pool, queries *repository.Queries) *Handler {
	return &Handler{
		svc: NewService(db, queries),
	}
}

func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) error {
	userId, err := httpx.ParseUUIDParam(r, "userId")
	if err != nil {
		return err
	}

	user, err := h.svc.GetUserById(r.Context(), userId)
	if err != nil {
		return err
	}

	httpx.ResponseWithJSON(w, http.StatusOK, mapUserToResponse(user))
	return nil
}

func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) error {
	userId, err := httpx.ParseUUIDParam(r, "userId")
	if err != nil {
		return err
	}

	if err = h.svc.DeleteUserById(r.Context(), userId); err != nil {
		return err
	}

	httpx.ResponseWithJSON(w, http.StatusNoContent, nil)
	return nil
}
