package controller

import (
	"context"
	"database/sql"
	"errors"
	"github.com/V2G-Minor-Fontys/server/internal/httpx"
	"github.com/V2G-Minor-Fontys/server/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

type Service interface {
	RegisterController(ctx context.Context, req *RegisterControllerRequest) (*Controller, error)
	PairUserToController(ctx context.Context, req *PairUserToControllerRequest) error
	GetControllerByUserId(ctx context.Context, id uuid.UUID) (*Controller, error)
	GetControllerByCpuId(ctx context.Context, cpuId string) (*Controller, error)
	GetControllerTelemetryById(ctx context.Context, id uuid.UUID) ([]*Telemetry, error)
	AddControllerTelemetry(ctx context.Context, req *AddControllerTelemetryRequest) error
	UpdateControllerSettings(ctx context.Context, req *UpdateControllerSettingsRequest) error
}

type ServiceImpl struct {
	db      *pgxpool.Pool
	queries *repository.Queries
}

func NewService(db *pgxpool.Pool, queries *repository.Queries) *ServiceImpl {
	return &ServiceImpl{db: db, queries: queries}
}

func (s *ServiceImpl) RegisterController(ctx context.Context, req *RegisterControllerRequest) (*Controller, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, httpx.InternalErr(ctx, "Failed to begin transaction", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.ErrorContext(ctx, "Rollback failed", "err", err)
		}
	}()

	qtx := s.queries.WithTx(tx)
	controllerParams := repository.AddControllerParams{
		ID:              uuid.New(),
		CpuID:           req.CpuID,
		FirmwareVersion: req.FirmwareVersion,
	}

	if err := qtx.AddController(ctx, controllerParams); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, httpx.Conflict(ctx, "Controller is already registered")
		}
	}

	controllerSettings := repository.AddControllerSettingsParams{
		ID:            controllerParams.ID,
		AutoStart:     true,
		HeartbeatRate: 5,
	}
	if err := qtx.AddControllerSettings(ctx, controllerSettings); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, httpx.Conflict(ctx, "Controller is already registered")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, httpx.InternalErr(ctx, "Failed to commit transaction", err)
	}

	return &Controller{
		ID:              controllerParams.ID,
		CpuID:           controllerParams.CpuID,
		FirmwareVersion: controllerParams.FirmwareVersion,
		Settings: &Settings{
			ID:            controllerParams.ID,
			AutoStart:     controllerSettings.AutoStart,
			HeartbeatRate: controllerSettings.HeartbeatRate,
			UpdatedAt:     time.Now().UTC(),
		},
	}, nil
}

func (s *ServiceImpl) PairUserToController(ctx context.Context, req *PairUserToControllerRequest) error {
	if err := s.queries.PairUserToController(ctx, repository.PairUserToControllerParams{
		CpuID:  req.CpuID,
		UserID: repository.GuidToPgUUID(req.UserId),
	}); err != nil {
		return httpx.InternalErr(ctx, "Could not pair the user to controller", err)
	}

	return nil
}

func (s *ServiceImpl) GetControllerByCpuId(ctx context.Context, cpuId string) (*Controller, error) {
	controllerRow, err := s.queries.GetControllerByCpuId(ctx, cpuId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, httpx.NotFound(ctx, "Controller not found")
	}

	return &Controller{
		ID:              controllerRow.ID,
		CpuID:           controllerRow.CpuID,
		FirmwareVersion: controllerRow.FirmwareVersion,
		Settings:        mapDatabaseSettingsToSettings(&controllerRow.ControllerSetting),
	}, nil
}

func (s *ServiceImpl) GetControllerTelemetryById(ctx context.Context, id uuid.UUID) ([]*Telemetry, error) {
	controllerTelemetries, err := s.queries.GetControllerTelemetryByControllerId(ctx, repository.GuidToPgUUID(id))
	if err != nil {
		return nil, httpx.InternalErr(ctx, "Unable to get controller telemetry", err)
	}

	if len(controllerTelemetries) == 0 {
		return nil, httpx.NotFound(ctx, "No telemetry found for controller")
	}

	return mapDatabaseTelemetrySliceToTelemetrySlice(controllerTelemetries), nil
}

func (s *ServiceImpl) GetControllerByUserId(ctx context.Context, id uuid.UUID) (*Controller, error) {
	controllerRow, err := s.queries.GetPairedControllerByUserId(ctx, repository.GuidToPgUUID(id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, httpx.NotFound(ctx, "Controller not found")
	}

	return &Controller{
		ID:              controllerRow.ID,
		CpuID:           controllerRow.CpuID,
		FirmwareVersion: controllerRow.FirmwareVersion,
		Settings:        mapDatabaseSettingsToSettings(&controllerRow.ControllerSetting),
	}, nil
}

func (s *ServiceImpl) AddControllerTelemetry(ctx context.Context, req *AddControllerTelemetryRequest) error {
	if err := s.queries.AddControllerTelemetry(ctx, repository.AddControllerTelemetryParams{
		ID:            uuid.New(),
		ControllerID:  repository.GuidToPgUUID(req.ControllerID),
		OutputPower:   req.OutputPower,
		Soc:           req.Soc,
		EvDischarging: req.EvDischarging,
	}); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return httpx.Conflict(ctx, "Controller is not registered")
		}
	}

	return nil
}

func (s *ServiceImpl) UpdateControllerSettings(ctx context.Context, req *UpdateControllerSettingsRequest) error {
	if err := s.queries.UpdateControllerSettings(ctx, repository.UpdateControllerSettingsParams{
		ID:            req.ID,
		AutoStart:     req.AutoStart,
		HeartbeatRate: req.Heartbeat,
	}); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return httpx.Conflict(ctx, "Controller settings is already registered")
		}
	}

	return nil
}
