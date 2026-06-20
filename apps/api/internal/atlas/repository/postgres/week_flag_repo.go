// FILE: apps/api/internal/atlas/repository/postgres/week_flag_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement WeekFlagRepository for WAVE-04 week flag tracking using sqlc-generated queries.
//   SCOPE: CRUD operations for week flags, list by week start date, all user-scoped.
//   DEPENDS: apps/api/internal/repository/postgres/generated, apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   WeekFlagRepository - Interface for week flag data access.
//   NewWeekFlagRepository - Creates a new WeekFlagRepository.
//   Create - Creates a week flag.
//   GetByID - Gets a week flag by ID (user-scoped).
//   ListByWeekStart - Lists week flags by week start date.
//   Delete - Deletes a week flag.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added week flag repository for WAVE-04.
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

type WeekFlagRepository interface {
	Create(ctx context.Context, userID string, weekStartDate models.Date, flagType string, notes *string) (*models.WeekFlagRecord, error)
	GetByID(ctx context.Context, userID string, id string) (*models.WeekFlagRecord, error)
	ListByWeekStart(ctx context.Context, userID string, weekStartDate models.Date) ([]models.WeekFlagRecord, error)
	Delete(ctx context.Context, userID string, id string) (*models.WeekFlagRecord, error)
}

type weekFlagRepository struct {
	q *generated.Queries
}

func NewWeekFlagRepository(pool *pgxpool.Pool) WeekFlagRepository {
	return &weekFlagRepository{q: generated.New(pool)}
}

func (r *weekFlagRepository) Create(ctx context.Context, userID string, weekStartDate models.Date, flagType string, notes *string) (*models.WeekFlagRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("week_flag_repo.Create: %w", err)
	}

	row, err := r.q.CreateWeekFlag(ctx, generated.CreateWeekFlagParams{
		UserID:        uid,
		WeekStartDate: modelsToPGDate(weekStartDate),
		FlagType:      flagType,
		Notes:         nullableText(notes),
	})
	if err != nil {
		return nil, fmt.Errorf("week_flag_repo.Create: %w", err)
	}

	return weekFlagRecordFromRow(row), nil
}

func (r *weekFlagRepository) GetByID(ctx context.Context, userID string, id string) (*models.WeekFlagRecord, error) {
	uid, fid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("week_flag_repo.GetByID: %w", err)
	}

	row, err := r.q.GetWeekFlagByID(ctx, generated.GetWeekFlagByIDParams{ID: fid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("week_flag_repo.GetByID: %w", err)
	}

	return weekFlagRecordFromRow(row), nil
}

func (r *weekFlagRepository) ListByWeekStart(ctx context.Context, userID string, weekStartDate models.Date) ([]models.WeekFlagRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("week_flag_repo.ListByWeekStart: %w", err)
	}

	rows, err := r.q.ListWeekFlagsByWeekStart(ctx, generated.ListWeekFlagsByWeekStartParams{
		WeekStartDate: modelsToPGDate(weekStartDate),
		UserID:        uid,
	})
	if err != nil {
		return nil, fmt.Errorf("week_flag_repo.ListByWeekStart: %w", err)
	}

	out := make([]models.WeekFlagRecord, len(rows))
	for i, row := range rows {
		out[i] = *weekFlagRecordFromRow(row)
	}
	return out, nil
}

func (r *weekFlagRepository) Delete(ctx context.Context, userID string, id string) (*models.WeekFlagRecord, error) {
	uid, fid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("week_flag_repo.Delete: %w", err)
	}

	row, err := r.q.DeleteWeekFlag(ctx, generated.DeleteWeekFlagParams{ID: fid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("week_flag_repo.Delete: %w", err)
	}

	return weekFlagRecordFromRow(row), nil
}

func weekFlagRecordFromRow(row generated.WeekFlag) *models.WeekFlagRecord {
	return &models.WeekFlagRecord{
		ID:            row.ID.String(),
		UserID:        row.UserID.String(),
		WeekStartDate: dateFromPGDate(row.WeekStartDate),
		FlagType:      row.FlagType,
		Notes:         textPtr(row.Notes),
		CreatedAt:     formatTimestamp(row.CreatedAt),
		UpdatedAt:     formatTimestamp(row.UpdatedAt),
	}
}