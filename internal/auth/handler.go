package auth

import (
	"github.com/V2G-Minor-Fontys/server/internal/config"
	"github.com/V2G-Minor-Fontys/server/internal/httpx"
	"github.com/V2G-Minor-Fontys/server/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"time"
)

const (
	TimeDay = time.Hour * 24
)

type Handler struct {
	svc Service
}

func NewHandler(cfg *config.Jwt, db *pgxpool.Pool, queries *repository.Queries) *Handler {
	return &Handler{
		svc: NewService(cfg, db, queries),
	}
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) error {
	var req RegisterRequest
	if err := httpx.DecodeJSONBody(r, &req); err != nil {
		return err
	}

	res, err := h.svc.Register(r.Context(), req)
	if err != nil {
		return err
	}

	httpx.SetRefreshToken(w, res.RefreshToken)
	httpx.ResponseWithJSON(w, http.StatusCreated, mapResultToResponse(res))

	return nil
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) error {
	var req LoginRequest
	if err := httpx.DecodeJSONBody(r, &req); err != nil {
		return err
	}

	res, err := h.svc.Login(r.Context(), req)
	if err != nil {
		return err
	}

	httpx.SetRefreshToken(w, res.RefreshToken)
	httpx.ResponseWithJSON(w, http.StatusOK, mapResultToResponse(res))

	return nil
}

func (h *Handler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	rtCookie, err := r.Cookie("refresh-token")
	if err != nil {
		return httpx.NotFound(ctx, "Refresh token could not be extracted")
	}

	res, err := h.svc.RefreshToken(ctx, rtCookie.Value)
	if err != nil {
		return err
	}

	httpx.SetRefreshToken(w, res.RefreshToken)
	httpx.ResponseWithJSON(w, http.StatusOK, mapResultToResponse(res))

	return nil
}

func (h *Handler) RevokeTokenHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	rtCookie, err := r.Cookie("refresh-token")
	if err != nil {
		return httpx.NotFound(ctx, "Refresh token could not be extracted")
	}

	err = h.svc.RevokeToken(ctx, rtCookie.Value)
	if err != nil {
		return err
	}

	httpx.SetRefreshToken(w, nil)
	httpx.ResponseWithJSON(w, http.StatusNoContent, nil)

	return nil
}
