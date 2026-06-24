// FILE: apps/api/internal/atlas/repository/postgres/daily_nutrition_log_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement DailyNutritionLogRepository using sqlc-generated factual daily nutrition log and entry queries.
//   SCOPE: Get-or-create, get by date, range list, notes update, entry add/list/update/delete. All reads and mutations are user-scoped; ListEntries requires both userID and dailyLogID.
//   DEPENDS: sqlc generated daily_nutrition_logs queries, atlas/models for records and inputs.
//   LINKS: M-API-NUTRITION / V-M-API-NUTRITION.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   DailyNutritionLogRepository - Interface for factual daily nutrition logs and entries.
//   NewDailyNutritionLogRepository - Creates a sqlc-backed daily nutrition repository.
//   GetOrCreate/GetByDate/ListByRange/UpdateNotes - Daily log operations.
//   AddEntry/ListEntries/UpdateEntry/DeleteEntry - User-scoped entry operations.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Task 2 repository adapter preserving user-scoped ListEntries.
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

type DailyNutritionLogRepository interface {
	GetOrCreate(ctx context.Context, userID string, date models.Date, notes *string) (*models.DailyNutritionLogRecord, error)
	GetByDate(ctx context.Context, userID string, date models.Date) (*models.DailyNutritionLogRecord, error)
	ListByRange(ctx context.Context, userID string, from, to models.Date) ([]models.DailyNutritionLogRecord, error)
	UpdateNotes(ctx context.Context, userID string, id string, notes *string) (*models.DailyNutritionLogRecord, error)
	AddEntry(ctx context.Context, userID string, input models.CreateDailyNutritionEntryRecordInput) (*models.DailyNutritionEntryRecord, error)
	ListEntries(ctx context.Context, userID string, dailyLogID string) ([]models.DailyNutritionEntryRecord, error)
	UpdateEntry(ctx context.Context, userID string, id string, input models.UpdateDailyNutritionEntryInput) (*models.DailyNutritionEntryRecord, error)
	DeleteEntry(ctx context.Context, userID string, id string) (*models.DailyNutritionEntryRecord, error)
}

type dailyNutritionLogRepository struct {
	q *generated.Queries
}

func NewDailyNutritionLogRepository(pool *pgxpool.Pool) DailyNutritionLogRepository {
	return &dailyNutritionLogRepository{q: generated.New(pool)}
}

func (r *dailyNutritionLogRepository) GetOrCreate(ctx context.Context, userID string, date models.Date, notes *string) (*models.DailyNutritionLogRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.GetOrCreate: %w", err)
	}

	row, err := r.q.CreateDailyNutritionLog(ctx, generated.CreateDailyNutritionLogParams{
		UserID: uid,
		Date:   modelsToPGDate(date),
		Notes:  nullableText(notes),
	})
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.GetOrCreate: %w", err)
	}

	return dailyNutritionLogRecordFromRow(row), nil
}

func (r *dailyNutritionLogRepository) GetByDate(ctx context.Context, userID string, date models.Date) (*models.DailyNutritionLogRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.GetByDate: %w", err)
	}

	row, err := r.q.GetDailyNutritionLogByDate(ctx, generated.GetDailyNutritionLogByDateParams{
		UserID: uid,
		Date:   modelsToPGDate(date),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("daily_nutrition_log_repo.GetByDate: %w", err)
	}

	return dailyNutritionLogRecordFromRow(row), nil
}

func (r *dailyNutritionLogRepository) ListByRange(ctx context.Context, userID string, from, to models.Date) ([]models.DailyNutritionLogRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.ListByRange: %w", err)
	}

	rows, err := r.q.ListDailyNutritionLogsByRange(ctx, generated.ListDailyNutritionLogsByRangeParams{
		UserID:    uid,
		StartDate: modelsToPGDate(from),
		EndDate:   modelsToPGDate(to),
	})
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.ListByRange: %w", err)
	}

	out := make([]models.DailyNutritionLogRecord, len(rows))
	for i, row := range rows {
		out[i] = *dailyNutritionLogRecordFromRow(row)
	}
	return out, nil
}

func (r *dailyNutritionLogRepository) UpdateNotes(ctx context.Context, userID string, id string, notes *string) (*models.DailyNutritionLogRecord, error) {
	uid, lid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.UpdateNotes: %w", err)
	}

	row, err := r.q.UpdateDailyNutritionLogNotes(ctx, generated.UpdateDailyNutritionLogNotesParams{
		ID:     lid,
		UserID: uid,
		Notes:  nullableText(notes),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("daily_nutrition_log_repo.UpdateNotes: %w", err)
	}

	return dailyNutritionLogRecordFromRow(row), nil
}

func (r *dailyNutritionLogRepository) AddEntry(ctx context.Context, userID string, input models.CreateDailyNutritionEntryRecordInput) (*models.DailyNutritionEntryRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.AddEntry: %w", err)
	}
	dailyLogID, productID, err := parseTwoUUIDs(input.DailyLogID, input.ProductID)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.AddEntry: %w", err)
	}

	row, err := r.q.CreateDailyNutritionEntry(ctx, generated.CreateDailyNutritionEntryParams{
		UserID:      uid,
		DailyLogID:  dailyLogID,
		ProductID:   productID,
		AmountGrams: float32(input.AmountGrams),
		MealLabel:   nullableText(input.MealLabel),
		Notes:       nullableText(input.Notes),
		Position:    input.Position,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("daily_nutrition_log_repo.AddEntry: %w", err)
	}

	return dailyNutritionEntryRecordFromRow(row), nil
}

func (r *dailyNutritionLogRepository) ListEntries(ctx context.Context, userID string, dailyLogID string) ([]models.DailyNutritionEntryRecord, error) {
	uid, lid, err := parseTwoUUIDs(userID, dailyLogID)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.ListEntries: %w", err)
	}

	rows, err := r.q.ListDailyNutritionEntriesByLog(ctx, generated.ListDailyNutritionEntriesByLogParams{
		UserID:     uid,
		DailyLogID: lid,
	})
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.ListEntries: %w", err)
	}

	out := make([]models.DailyNutritionEntryRecord, len(rows))
	for i, row := range rows {
		out[i] = *dailyNutritionEntryRecordFromRow(row)
	}
	return out, nil
}

func (r *dailyNutritionLogRepository) UpdateEntry(ctx context.Context, userID string, id string, input models.UpdateDailyNutritionEntryInput) (*models.DailyNutritionEntryRecord, error) {
	if input.AmountGrams == nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.UpdateEntry: amountGrams is required")
	}
	uid, eid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.UpdateEntry: %w", err)
	}
	lid, err := uuidFromString(input.DailyLogID)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.UpdateEntry: %w", err)
	}

	position := int32(0)
	if input.Position != nil {
		position = *input.Position
	}

	row, err := r.q.UpdateDailyNutritionEntry(ctx, generated.UpdateDailyNutritionEntryParams{
		ID:          eid,
		UserID:      uid,
		DailyLogID:  lid,
		AmountGrams: float32(*input.AmountGrams),
		MealLabel:   nullableText(input.MealLabel),
		Notes:       nullableText(input.Notes),
		Position:    position,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("daily_nutrition_log_repo.UpdateEntry: %w", err)
	}

	return dailyNutritionEntryRecordFromRow(row), nil
}

func (r *dailyNutritionLogRepository) DeleteEntry(ctx context.Context, userID string, id string) (*models.DailyNutritionEntryRecord, error) {
	uid, eid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.DeleteEntry: %w", err)
	}

	row, err := r.q.DeleteDailyNutritionEntry(ctx, generated.DeleteDailyNutritionEntryParams{ID: eid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("daily_nutrition_log_repo.DeleteEntry: %w", err)
	}

	return dailyNutritionEntryRecordFromRow(row), nil
}

func dailyNutritionLogRecordFromRow(row generated.DailyNutritionLog) *models.DailyNutritionLogRecord {
	return &models.DailyNutritionLogRecord{
		ID:        row.ID.String(),
		UserID:    row.UserID.String(),
		Date:      dateFromPGDate(row.Date),
		Notes:     textPtr(row.Notes),
		CreatedAt: formatTimestamp(row.CreatedAt),
		UpdatedAt: formatTimestamp(row.UpdatedAt),
	}
}

func dailyNutritionEntryRecordFromRow(row generated.DailyNutritionEntry) *models.DailyNutritionEntryRecord {
	return &models.DailyNutritionEntryRecord{
		ID:                      row.ID.String(),
		DailyLogID:              row.DailyLogID.String(),
		ProductID:               row.ProductID.String(),
		ProductNameSnapshot:     row.ProductNameSnapshot,
		CaloriesPer100gSnapshot: float64(row.CaloriesPer100gSnapshot),
		ProteinPer100gSnapshot:  float64(row.ProteinPer100gSnapshot),
		FatPer100gSnapshot:      float64(row.FatPer100gSnapshot),
		CarbsPer100gSnapshot:    float64(row.CarbsPer100gSnapshot),
		AmountGrams:             float64(row.AmountGrams),
		MealLabel:               textPtr(row.MealLabel),
		Notes:                   textPtr(row.Notes),
		Position:                row.Position,
		CreatedAt:               formatTimestamp(row.CreatedAt),
		UpdatedAt:               formatTimestamp(row.UpdatedAt),
	}
}
