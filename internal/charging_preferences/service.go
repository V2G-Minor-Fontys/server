package charging_preferences

import (
	"github.com/V2G-Minor-Fontys/server/internal/config"
	"github.com/V2G-Minor-Fontys/server/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	cfg     *config.Jwt
	db      *pgxpool.Pool
	queries *repository.Queries
}

func NewService(cfg *config.Jwt, db *pgxpool.Pool, queries *repository.Queries) *Service {
	return &Service{cfg: cfg, db: db, queries: queries}
}
