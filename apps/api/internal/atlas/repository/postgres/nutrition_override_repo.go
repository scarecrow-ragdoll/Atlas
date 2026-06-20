// FILE: apps/api/internal/atlas/repository/postgres/nutrition_override_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement DailyNutritionOverrideRepository using sqlc-generated queries for the daily_nutrition_override table.
//   SCOPE: Upsert, GetByID, GetByDate, ListByRange, Update, Delete. All user-scoped. Uses pgtype.Date for date column. Not-found returns nil, nil.
//   DEPENDS: sqlc generated DailyNutritionOverride model, atlas/models for record types.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT

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

type DailyNutritionOverrideRepository interface {
	Upsert(ctx context.Context, userID string, date string, notes *string) (*models.DailyNutritionOverrideRecord, error)
	GetByID(ctx context.Context, userID string, id string) (*models.DailyNutritionOverrideRecord, error)
	GetByDate(ctx context.Context, userID string, date string) (*models.DailyNutritionOverrideRecord, error)
	ListByRange(ctx context.Context, userID string, startDate, endDate string) ([]models.DailyNutritionOverrideRecord, error)
	Update(ctx context.Context, userID string, id string, notes *string) (*models.DailyNutritionOverrideRecord, error)
	Delete(ctx context.Context, userID string, id string) (*models.DailyNutritionOverrideRecord, error)
}

type dailyNutritionOverrideRepository struct {
	q *generated.Queries
}

func NewDailyNutritionOverrideRepository(pool *pgxpool.Pool) DailyNutritionOverrideRepository {
	return &dailyNutritionOverrideRepository{q: generated.New(pool)}
}

func (r *dailyNutritionOverrideRepository) Upsert(ctx context.Context, userID string, date string, notes *string) (*models.DailyNutritionOverrideRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.Upsert: %w", err)
	}

	d, err := parseDate(date)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.Upsert: %w", err)
	}

	row, err := r.q.UpsertDailyNutritionOverride(ctx, generated.UpsertDailyNutritionOverrideParams{
		UserID: uid,
		Date:   d,
		Notes:  nullableText(notes),
	})
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.Upsert: %w", err)
	}

	return dailyNutritionOverrideRecordFromRow(row), nil
}

func (r *dailyNutritionOverrideRepository) GetByID(ctx context.Context, userID string, id string) (*models.DailyNutritionOverrideRecord, error) {
	uid, oid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.GetByID: %w", err)
	}

	row, err := r.q.GetDailyNutritionOverrideByID(ctx, generated.GetDailyNutritionOverrideByIDParams{ID: oid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_override_repo.GetByID: %w", err)
	}

	return dailyNutritionOverrideRecordFromRow(row), nil
}

func (r *dailyNutritionOverrideRepository) GetByDate(ctx context.Context, userID string, date string) (*models.DailyNutritionOverrideRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.GetByDate: %w", err)
	}

	d, err := parseDate(date)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.GetByDate: %w", err)
	}

	row, err := r.q.GetDailyNutritionOverrideByDate(ctx, generated.GetDailyNutritionOverrideByDateParams{
		UserID: uid,
		Date:   d,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_override_repo.GetByDate: %w", err)
	}

	return dailyNutritionOverrideRecordFromRow(row), nil
}

func (r *dailyNutritionOverrideRepository) ListByRange(ctx context.Context, userID string, startDate, endDate string) ([]models.DailyNutritionOverrideRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.ListByRange: %w", err)
	}

	sd, err := parseDate(startDate)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.ListByRange: %w", err)
	}

	ed, err := parseDate(endDate)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.ListByRange: %w", err)
	}

	rows, err := r.q.ListDailyNutritionOverridesByRange(ctx, generated.ListDailyNutritionOverridesByRangeParams{
		UserID: uid,
		Date:   sd,
		Date_2: ed,
	})
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.ListByRange: %w", err)
	}

	out := make([]models.DailyNutritionOverrideRecord, len(rows))
	for i, row := range rows {
		out[i] = *dailyNutritionOverrideRecordFromRow(row)
	}
	return out, nil
}

func (r *dailyNutritionOverrideRepository) Update(ctx context.Context, userID string, id string, notes *string) (*models.DailyNutritionOverrideRecord, error) {
	uid, oid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.Update: %w", err)
	}

	row, err := r.q.UpdateDailyNutritionOverride(ctx, generated.UpdateDailyNutritionOverrideParams{
		ID:    oid,
		UserID: uid,
		Notes: nullableText(notes),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_override_repo.Update: %w", err)
	}

	return dailyNutritionOverrideRecordFromRow(row), nil
}

func (r *dailyNutritionOverrideRepository) Delete(ctx context.Context, userID string, id string) (*models.DailyNutritionOverrideRecord, error) {
	uid, oid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_repo.Delete: %w", err)
	}

	row, err := r.q.DeleteDailyNutritionOverride(ctx, generated.DeleteDailyNutritionOverrideParams{ID: oid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_override_repo.Delete: %w", err)
	}

	return dailyNutritionOverrideRecordFromRow(row), nil
}

func dailyNutritionOverrideRecordFromRow(row generated.DailyNutritionOverride) *models.DailyNutritionOverrideRecord {
	return &models.DailyNutritionOverrideRecord{
		ID:        row.ID.String(),
		UserID:    row.UserID.String(),
		Date:      dateFromPGDate(row.Date),
		Notes:     textPtr(row.Notes),
		CreatedAt: formatTimestamp(row.CreatedAt),
		UpdatedAt: formatTimestamp(row.UpdatedAt),
	}
}
