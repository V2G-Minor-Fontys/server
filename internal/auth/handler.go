package auth

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/V2G-Minor-Fontys/server/internal/config"
	"github.com/V2G-Minor-Fontys/server/internal/httpx"
	"github.com/V2G-Minor-Fontys/server/internal/repository"
	"github.com/V2G-Minor-Fontys/server/pkg/crypto"
	"github.com/V2G-Minor-Fontys/server/pkg/jwt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"net/http"
	"time"
)

const (
	TimeDay = time.Hour * 24
)

type Handler struct {
	cfg     *config.Jwt
	db      *pgxpool.Pool
	queries *repository.Queries
}

func NewHandler(cfg *config.Jwt, db *pgxpool.Pool, queries *repository.Queries) *Handler {
	return &Handler{
		cfg:     cfg,
		db:      db,
		queries: queries,
	}
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tx, err := h.db.Begin(ctx)
	if err != nil {
		return httpx.InternalErr("Transaction could not be created", r.RequestURI, err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.ErrorContext(ctx, "Rollback failed", "err", err)
		}
	}()

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httpx.BadRequest("Could not parse JSON body", r.RequestURI)
	}

	params, err := req.ToRegisterParams()
	if err != nil {
		return httpx.BadRequest("Invalid register request", r.RequestURI)
	}

	qtx := h.queries.WithTx(tx)
	if err := qtx.Register(ctx, *params); err != nil {
		return httpx.BadRequest("Bad Request Error", r.RequestURI)
	}

	rtBytes, err := jwt.GenerateRefreshToken()
	if err != nil {
		return httpx.InternalErr("Could not generate refresh token", r.RequestURI, err)
	}

	at, err := jwt.GenerateAccessToken(params.ID.String(), h.cfg)
	if err != nil {
		return httpx.InternalErr("Could not generate access token", r.RequestURI, err)
	}

	rt, err := qtx.CreateRefreshToken(ctx, repository.CreateRefreshTokenParams{
		Token:      rtBytes,
		IdentityID: params.ID,
		ExpiresAt:  time.Now().UTC().Add(TimeDay * 30),
	})
	if err != nil {
		return httpx.BadRequest("Bad Request Error", r.RequestURI)
	}

	if err := qtx.CreateUser(ctx, repository.CreateUserParams{
		ID:       params.ID,
		Username: params.Username,
	}); err != nil {
		return httpx.BadRequest("Bad Request Error", r.RequestURI)
	}

	if err := tx.Commit(ctx); err != nil {
		return httpx.InternalErr("Transaction could not be committed", r.RequestURI, err)
	}

	httpx.SetRefreshToken(w, &rt)
	httpx.ResponseWithJSON(w, http.StatusCreated, AuthenticationResponse{
		ID:          params.ID,
		AccessToken: at,
	})

	return nil
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httpx.BadRequest("Invalid JSON body", r.RequestURI)
	}

	identity, err := h.queries.GetIdentityByUsername(ctx, req.Username)
	if err != nil {
		return httpx.BadRequest("Invalid Identity Request", r.RequestURI)
	}

	if match := crypto.CheckPasswordHash(req.Password, identity.PasswordHash); !match {
		return httpx.BadRequest("Invalid username or password.", r.RequestURI)
	}

	rt, err := h.queries.GetRefreshTokenByIdentityId(ctx, identity.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return httpx.InternalErr("Could not get refresh token", r.RequestURI, err)
		}

		rtBytes, err := jwt.GenerateRefreshToken()
		if err != nil {
			return httpx.InternalErr("Could not generate refresh token", r.RequestURI, err)
		}

		rt, err = h.queries.CreateRefreshToken(ctx, repository.CreateRefreshTokenParams{
			Token:      rtBytes,
			IdentityID: identity.ID,
			ExpiresAt:  time.Now().UTC().Add(TimeDay * 30),
		})
	}

	at, err := jwt.GenerateAccessToken(identity.ID.String(), h.cfg)
	if err != nil {
		return httpx.InternalErr("Could not generate access token", r.RequestURI, err)
	}

	httpx.SetRefreshToken(w, &rt)
	httpx.ResponseWithJSON(w, http.StatusCreated, AuthenticationResponse{
		ID:          identity.ID,
		AccessToken: at,
	})

	return nil
}

func (h *Handler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	rtCookie, err := r.Cookie("refresh-token")
	if err != nil {
		return httpx.NotFound("Refresh token could not be extracted", r.RequestURI)
	}

	rtBytes, err := hex.DecodeString(rtCookie.Value)
	if err != nil {
		return httpx.BadRequest("Refresh token could not be decoded", r.RequestURI)
	}

	rt, err := h.queries.GetRefreshToken(ctx, rtBytes)
	if err != nil {
		return httpx.NotFound("Refresh token could not be found", r.RequestURI)
	}

	identity, err := h.queries.GetIdentityById(ctx, rt.IdentityID)
	if err != nil {
		return httpx.BadRequest("Identity attached to this refresh token is not longer exists.", r.RequestURI)
	}

	at, err := jwt.GenerateAccessToken(identity.ID.String(), h.cfg)
	if err != nil {
		return httpx.InternalErr("Could not generate access token", r.RequestURI, err)
	}

	httpx.SetRefreshToken(w, &rt)
	httpx.ResponseWithJSON(w, http.StatusCreated, AuthenticationResponse{
		ID:          identity.ID,
		AccessToken: at,
	})

	return nil
}

func (h *Handler) RevokeTokenHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	rtCookie, err := r.Cookie("refresh-token")
	if err != nil {
		return httpx.NotFound("Refresh token could not be extracted", r.RequestURI)
	}

	rtBytes, err := hex.DecodeString(rtCookie.Value)
	if err != nil {
		return httpx.BadRequest("Refresh token could not be decoded", r.RequestURI)
	}

	rows, err := h.queries.DeleteRefreshToken(ctx, rtBytes)
	if err != nil {
		return httpx.InternalErr("Unexpected Error Occurred.", r.RequestURI, err)
	}

	if rows == 0 {
		return httpx.NotFound("Refresh token could not be found", r.RequestURI)
	}

	httpx.SetRefreshToken(w, nil)
	httpx.ResponseWithJSON(w, http.StatusNoContent, nil)

	return nil
}
