package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/V2G-Minor-Fontys/server/internal/httpx"
	"github.com/V2G-Minor-Fontys/server/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	GetUserById(ctx context.Context, id uuid.UUID) (*User, error)
	DeleteUserById(ctx context.Context, id uuid.UUID) error
}

type ServiceImpl struct {
	db      *pgxpool.Pool
	queries *repository.Queries
}

func NewService(db *pgxpool.Pool, queries *repository.Queries) Service {
	return &ServiceImpl{db: db, queries: queries}
}

func (s *ServiceImpl) GetUserById(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := s.queries.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, httpx.NotFound(ctx, "User Not Found")
		}

		return nil, httpx.InternalErr(ctx, "Unexpected Error Occurred", err)
	}

	return mapDatabaseUserToUser(&user), nil
}

func (s *ServiceImpl) DeleteUserById(ctx context.Context, id uuid.UUID) error {
	rows, err := s.queries.DeleteIdentityById(ctx, id)
	if err != nil {
		return httpx.InternalErr(ctx, "Unexpected Error Occurred", err)
	}

	if rows < 1 {
		return httpx.NotFound(ctx, "User Not Found")
	}

	return nil
}
