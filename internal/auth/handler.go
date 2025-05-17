package auth

import (
	"encoding/json"
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
	ctx := r.Context()
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httpx.BadRequest(ctx, "Could not parse JSON body")
	}

	res, err := h.svc.Register(ctx, req)
	if err != nil {
		return err
	}

	httpx.SetRefreshToken(w, res.RefreshToken)
	httpx.ResponseWithJSON(w, http.StatusCreated, res.ToAuthenticationResponse())

	return nil
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httpx.BadRequest(ctx, "Invalid JSON body")
	}

	res, err := h.svc.Login(ctx, req)
	if err != nil {
		return err
	}

	httpx.SetRefreshToken(w, res.RefreshToken)
	httpx.ResponseWithJSON(w, http.StatusOK, res.ToAuthenticationResponse())

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
	httpx.ResponseWithJSON(w, http.StatusOK, res.ToAuthenticationResponse())

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
