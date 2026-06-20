// FILE: apps/api/internal/atlas/repository/postgres/body_check_in_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement BodyCheckInRepository for WAVE-04 body check-in tracking using sqlc-generated queries.
//   SCOPE: CRUD operations for body check-ins, list by date range, all user-scoped.
//   DEPENDS: apps/api/internal/repository/postgres/generated, apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   BodyCheckInRepository - Interface for body check-in data access.
//   NewBodyCheckInRepository - Creates a new BodyCheckInRepository.
//   Create - Creates a body check-in.
//   GetByID - Gets a body check-in by ID (user-scoped).
//   ListByDateRange - Lists body check-ins by date range.
//   Update - Updates a body check-in.
//   Delete - Deletes a body check-in.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added body check-in repository for WAVE-04.
// END_CHANGE_SUMMARY

package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/repository/postgres/generated"
)

type BodyCheckInRepository interface {
	Create(ctx context.Context, userID string, date models.Date, weight *float64, bodyFatPercentage *float64, notes *string) (*models.BodyCheckInRecord, error)
	GetByID(ctx context.Context, userID string, id string) (*models.BodyCheckInRecord, error)
	ListByDateRange(ctx context.Context, userID string, fromDate models.Date, toDate models.Date) ([]models.BodyCheckInRecord, error)
	Update(ctx context.Context, userID string, id string, weight *float64, bodyFatPercentage *float64, notes *string) (*models.BodyCheckInRecord, error)
	Delete(ctx context.Context, userID string, id string) (*models.BodyCheckInRecord, error)
}

type bodyCheckInRepository struct {
	q *generated.Queries
}

func NewBodyCheckInRepository(pool *pgxpool.Pool) BodyCheckInRepository {
	return &bodyCheckInRepository{q: generated.New(pool)}
}

func (r *bodyCheckInRepository) Create(ctx context.Context, userID string, date models.Date, weight *float64, bodyFatPercentage *float64, notes *string) (*models.BodyCheckInRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("body_checkin_repo.Create: %w", err)
	}

	row, err := r.q.CreateBodyCheckIn(ctx, generated.CreateBodyCheckInParams{
		UserID:            uid,
		Date:              modelsToPGDate(date),
		Weight:            nullableFloat4(weight),
		BodyFatPercentage: nullableFloat4(bodyFatPercentage),
		Notes:             nullableText(notes),
	})
	if err != nil {
		return nil, fmt.Errorf("body_checkin_repo.Create: %w", err)
	}

	return bodyCheckInRecordFromRow(row), nil
}

func (r *bodyCheckInRepository) GetByID(ctx context.Context, userID string, id string) (*models.BodyCheckInRecord, error) {
	uid, cid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("body_checkin_repo.GetByID: %w", err)
	}

	row, err := r.q.GetBodyCheckInByID(ctx, generated.GetBodyCheckInByIDParams{ID: cid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("body_checkin_repo.GetByID: %w", err)
	}

	return bodyCheckInRecordFromRow(row), nil
}

func (r *bodyCheckInRepository) ListByDateRange(ctx context.Context, userID string, fromDate models.Date, toDate models.Date) ([]models.BodyCheckInRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("body_checkin_repo.ListByDateRange: %w", err)
	}

	rows, err := r.q.ListBodyCheckInsByDateRange(ctx, generated.ListBodyCheckInsByDateRangeParams{
		UserID: uid,
		Date:   modelsToPGDate(fromDate),
		Date_2: modelsToPGDate(toDate),
	})
	if err != nil {
		return nil, fmt.Errorf("body_checkin_repo.ListByDateRange: %w", err)
	}

	out := make([]models.BodyCheckInRecord, len(rows))
	for i, row := range rows {
		out[i] = *bodyCheckInRecordFromRow(row)
	}
	return out, nil
}

func (r *bodyCheckInRepository) Update(ctx context.Context, userID string, id string, weight *float64, bodyFatPercentage *float64, notes *string) (*models.BodyCheckInRecord, error) {
	uid, cid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("body_checkin_repo.Update: %w", err)
	}

	row, err := r.q.UpdateBodyCheckIn(ctx, generated.UpdateBodyCheckInParams{
		ID:                cid,
		UserID:            uid,
		Weight:            nullableFloat4(weight),
		BodyFatPercentage: nullableFloat4(bodyFatPercentage),
		Notes:             nullableText(notes),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("body_checkin_repo.Update: %w", err)
	}

	return bodyCheckInRecordFromRow(row), nil
}

func (r *bodyCheckInRepository) Delete(ctx context.Context, userID string, id string) (*models.BodyCheckInRecord, error) {
	uid, cid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("body_checkin_repo.Delete: %w", err)
	}

	row, err := r.q.DeleteBodyCheckIn(ctx, generated.DeleteBodyCheckInParams{ID: cid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("body_checkin_repo.Delete: %w", err)
	}

	return bodyCheckInRecordFromRow(row), nil
}

func bodyCheckInRecordFromRow(row generated.BodyCheckIn) *models.BodyCheckInRecord {
	return &models.BodyCheckInRecord{
		ID:                row.ID.String(),
		UserID:            row.UserID.String(),
		Date:              dateFromPGDate(row.Date),
		Weight:            float4Ptr(row.Weight),
		BodyFatPercentage: float4Ptr(row.BodyFatPercentage),
		Notes:             textPtr(row.Notes),
		CreatedAt:         formatTimestamp(row.CreatedAt),
		UpdatedAt:         formatTimestamp(row.UpdatedAt),
	}
}