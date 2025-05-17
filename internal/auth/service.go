package auth

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/V2G-Minor-Fontys/server/internal/config"
	"github.com/V2G-Minor-Fontys/server/internal/httpx"
	"github.com/V2G-Minor-Fontys/server/internal/repository"
	"github.com/V2G-Minor-Fontys/server/pkg/crypto"
	"github.com/V2G-Minor-Fontys/server/pkg/jwt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthenticationResult, error)
	Login(ctx context.Context, req LoginRequest) (*AuthenticationResult, error)
	RefreshToken(ctx context.Context, token string) (*AuthenticationResult, error)
	RevokeToken(ctx context.Context, token string) error
}

type ServiceImpl struct {
	cfg     *config.Jwt
	db      *pgxpool.Pool
	queries *repository.Queries
}

func NewService(cfg *config.Jwt, db *pgxpool.Pool, queries *repository.Queries) *ServiceImpl {
	return &ServiceImpl{cfg: cfg, db: db, queries: queries}
}

func (s *ServiceImpl) Register(ctx context.Context, req RegisterRequest) (*AuthenticationResult, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, httpx.InternalErr(ctx, "Failed to begin transaction", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.ErrorContext(ctx, "Rollback failed", "err", err)
		}
	}()

	params, err := mapRegisterRequestToParams(crypto.HashPassword, &req)
	if err != nil {
		return nil, httpx.BadRequest(ctx, "Invalid registration data provided")
	}

	qtx := s.queries.WithTx(tx)
	if err := qtx.Register(ctx, *params); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, httpx.Conflict(ctx, "Username is already in use")
		}

		return nil, httpx.InternalErr(ctx, "Failed to register identity", err)
	}

	rtBytes, err := jwt.GenerateRefreshToken()
	if err != nil {
		return nil, httpx.InternalErr(ctx, "Could not generate refresh token", err)
	}

	at, err := jwt.GenerateAccessToken(params.ID.String(), s.cfg)
	if err != nil {
		return nil, httpx.InternalErr(ctx, "Could not generate access token", err)
	}

	rt, err := qtx.CreateRefreshToken(ctx, repository.CreateRefreshTokenParams{
		Token:      rtBytes,
		IdentityID: params.ID,
		ExpiresAt:  time.Now().UTC().Add(TimeDay * 30),
	})
	if err != nil {
		return nil, httpx.InternalErr(ctx, "Could not store refresh token", err)
	}

	if err := qtx.CreateUser(ctx, repository.CreateUserParams{
		ID:       params.ID,
		Username: params.Username,
	}); err != nil {
		return nil, httpx.InternalErr(ctx, "Could not create user profile", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, httpx.InternalErr(ctx, "Failed to commit transaction", err)
	}

	return &AuthenticationResult{
		ID:           params.ID,
		AccessToken:  at,
		RefreshToken: mapDatabaseRefreshTokenToToken(&rt),
	}, nil
}

func (s *ServiceImpl) Login(ctx context.Context, req LoginRequest) (*AuthenticationResult, error) {
	identity, err := s.queries.GetIdentityByUsername(ctx, req.Username)
	if err != nil {
		return nil, httpx.BadRequest(ctx, "Invalid username or password")
	}

	if match := crypto.CheckPasswordHash(req.Password, identity.PasswordHash); !match {
		return nil, httpx.BadRequest(ctx, "Invalid username or password")
	}

	rt, err := s.queries.GetRefreshTokenByIdentityId(ctx, identity.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, httpx.InternalErr(ctx, "Refresh token could not be retrieved", err)
		}

		rtBytes, err := jwt.GenerateRefreshToken()
		if err != nil {
			return nil, httpx.InternalErr(ctx, "Could not generate refresh token", err)
		}

		if rt, err = s.queries.CreateRefreshToken(ctx, repository.CreateRefreshTokenParams{
			Token:      rtBytes,
			IdentityID: identity.ID,
			ExpiresAt:  time.Now().UTC().Add(TimeDay * 30),
		}); err != nil {
			return nil, httpx.InternalErr(ctx, "Could not store refresh token", err)
		}
	}

	at, err := jwt.GenerateAccessToken(identity.ID.String(), s.cfg)
	if err != nil {
		return nil, httpx.InternalErr(ctx, "Could not generate access token", err)
	}

	return &AuthenticationResult{
		ID:           identity.ID,
		AccessToken:  at,
		RefreshToken: mapDatabaseRefreshTokenToToken(&rt),
	}, nil
}

func (s *ServiceImpl) RefreshToken(ctx context.Context, token string) (*AuthenticationResult, error) {
	rtBytes, err := hex.DecodeString(token)
	if err != nil {
		return nil, httpx.BadRequest(ctx, "Refresh token could not be decoded")
	}

	rt, err := s.queries.GetRefreshToken(ctx, rtBytes)
	if err != nil {
		return nil, httpx.NotFound(ctx, "Refresh token could not be found")
	}

	identity, err := s.queries.GetIdentityById(ctx, rt.IdentityID)
	if err != nil {
		return nil, httpx.BadRequest(ctx, "Identity attached to this refresh token no longer exists")
	}

	at, err := jwt.GenerateAccessToken(identity.ID.String(), s.cfg)
	if err != nil {
		return nil, httpx.InternalErr(ctx, "Could not generate access token", err)
	}

	return &AuthenticationResult{
		ID:           identity.ID,
		AccessToken:  at,
		RefreshToken: mapDatabaseRefreshTokenToToken(&rt),
	}, nil
}

func (s *ServiceImpl) RevokeToken(ctx context.Context, token string) error {
	rtBytes, err := hex.DecodeString(token)
	if err != nil {
		return httpx.BadRequest(ctx, "Refresh token could not be decoded")
	}

	rows, err := s.queries.DeleteRefreshToken(ctx, rtBytes)
	if err != nil {
		return httpx.InternalErr(ctx, "Unexpected error occurred during refresh token deletion. Please try again", err)
	}

	if rows == 0 {
		return httpx.NotFound(ctx, "Refresh token could not be found")
	}

	return nil
}
