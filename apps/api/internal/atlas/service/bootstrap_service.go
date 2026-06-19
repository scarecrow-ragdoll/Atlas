// FILE: apps/api/internal/atlas/service/bootstrap_service.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement the BootstrapService that ensures the default user and default settings exist on startup.
//   SCOPE: EnsureDefaultUser creates the first atlas_users row if empty; EnsureDefaultSettings creates atlas_settings row with defaults.
//   DEPENDS: apps/api/internal/repository/postgres/generated (sqlc queries for atlas_users and atlas_settings).
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Atlas bootstrap service for WAVE-01.
// END_CHANGE_SUMMARY

package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"monorepo-template/apps/api/internal/repository/postgres/generated"
)

var (
	ErrBootstrapFailed = errors.New("atlas bootstrap failed")
)

type BootstrapService interface {
	EnsureDefaultUser(ctx context.Context) (string, error)
	EnsureDefaultSettings(ctx context.Context, userID string) error
}

type bootstrapService struct {
	pool *pgxpool.Pool
	q    *generated.Queries
}

func NewBootstrapService(pool *pgxpool.Pool) BootstrapService {
	return &bootstrapService{
		pool: pool,
		q:    generated.New(pool),
	}
}

func (s *bootstrapService) EnsureDefaultUser(ctx context.Context) (string, error) {
	existing, err := s.q.GetAtlasDefaultUser(ctx)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("bootstrap_service.EnsureDefaultUser: query: %w", err)
	}

	if err == nil && existing.ID.Valid {
		return existing.ID.String(), nil
	}

	userID, err := s.q.InsertAtlasDefaultUser(ctx, "Default User")
	if err != nil {
		return "", fmt.Errorf("bootstrap_service.EnsureDefaultUser: insert: %w", err)
	}

	return userID.String(), nil
}

func (s *bootstrapService) EnsureDefaultSettings(ctx context.Context, userID string) error {
	uid, err := uuidFromString(userID)
	if err != nil {
		return fmt.Errorf("bootstrap_service.EnsureDefaultSettings: %w", err)
	}

	if err := s.q.CreateAtlasSettings(ctx, uid); err != nil {
		return fmt.Errorf("bootstrap_service.EnsureDefaultSettings: %w", err)
	}

	return nil
}

func uuidFromString(value string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(value); err != nil {
		return pgtype.UUID{}, err
	}
	return uuid, nil
}