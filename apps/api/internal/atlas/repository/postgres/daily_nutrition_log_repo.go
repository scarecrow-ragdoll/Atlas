// FILE: apps/api/internal/atlas/repository/postgres/daily_nutrition_log_repo.go
// VERSION: 1.0.3
// START_MODULE_CONTRACT
//   PURPOSE: Implement DailyNutritionLogRepository using sqlc-generated factual daily nutrition log and entry queries.
//   SCOPE: Get-or-create, get by ID, get by date, range list, notes update, entry add/list/update/delete, and atomic template seeding. All reads and mutations are user-scoped; ListEntries requires both userID and dailyLogID.
//   DEPENDS: sqlc generated daily_nutrition_logs queries, atlas/models for records and inputs.
//   LINKS: M-API-NUTRITION / V-M-API-NUTRITION.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   DailyNutritionLogRepository - Interface for factual daily nutrition logs and entries.
//   NewDailyNutritionLogRepository - Creates a sqlc-backed daily nutrition repository.
//   GetOrCreate/GetByID/GetByDate/ListByRange/UpdateNotes - Daily log operations.
//   AddEntry/ListEntries/UpdateEntry/DeleteEntry - User-scoped entry operations.
//   SeedEntriesIfEmpty - Transactional per-date seed helper for weekly template apply.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.3 - Added transactional seed-if-empty helper for weekly template apply.
// END_CHANGE_SUMMARY

package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/repository/postgres/generated"
)

type DailyNutritionLogRepository interface {
	GetOrCreate(ctx context.Context, userID string, date models.Date, notes *string) (*models.DailyNutritionLogRecord, error)
	GetByID(ctx context.Context, userID string, id string) (*models.DailyNutritionLogRecord, error)
	GetByDate(ctx context.Context, userID string, date models.Date) (*models.DailyNutritionLogRecord, error)
	ListByRange(ctx context.Context, userID string, from, to models.Date) ([]models.DailyNutritionLogRecord, error)
	UpdateNotes(ctx context.Context, userID string, id string, notes *string) (*models.DailyNutritionLogRecord, error)
	AddEntry(ctx context.Context, userID string, input models.CreateDailyNutritionEntryRecordInput) (*models.DailyNutritionEntryRecord, error)
	ListEntries(ctx context.Context, userID string, dailyLogID string) ([]models.DailyNutritionEntryRecord, error)
	UpdateEntry(ctx context.Context, userID string, id string, input models.UpdateDailyNutritionEntryInput) (*models.DailyNutritionEntryRecord, error)
	DeleteEntry(ctx context.Context, userID string, id string) (*models.DailyNutritionEntryRecord, error)
	SeedEntriesIfEmpty(ctx context.Context, userID string, date models.Date, items []models.DailyNutritionSeedEntryInput) (*models.DailyNutritionSeedResult, error)
}

type dailyNutritionLogRepository struct {
	q    *generated.Queries
	pool *pgxpool.Pool
}

func NewDailyNutritionLogRepository(pool *pgxpool.Pool) DailyNutritionLogRepository {
	return &dailyNutritionLogRepository{q: generated.New(pool), pool: pool}
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

func (r *dailyNutritionLogRepository) GetByID(ctx context.Context, userID string, id string) (*models.DailyNutritionLogRecord, error) {
	uid, lid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.GetByID: %w", err)
	}

	row := generated.DailyNutritionLog{}
	err = r.pool.QueryRow(ctx, `
SELECT id, user_id, date, notes, created_at, updated_at
FROM daily_nutrition_logs
WHERE id = $1 AND user_id = $2
LIMIT 1
`, lid, uid).Scan(&row.ID, &row.UserID, &row.Date, &row.Notes, &row.CreatedAt, &row.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("daily_nutrition_log_repo.GetByID: %w", err)
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
	uid, eid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.UpdateEntry: %w", err)
	}
	lid, err := uuidFromString(input.DailyLogID)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.UpdateEntry: %w", err)
	}

	row, err := r.q.UpdateDailyNutritionEntry(ctx, generated.UpdateDailyNutritionEntryParams{
		ID:          eid,
		UserID:      uid,
		DailyLogID:  lid,
		AmountGrams: float32(input.AmountGrams),
		MealLabel:   nullableText(input.MealLabel),
		Notes:       nullableText(input.Notes),
		Position:    input.Position,
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

// START_CONTRACT: SeedEntriesIfEmpty
//
//	PURPOSE: Atomically seed one user's date with planned nutrition entries only when the factual day is empty.
//	INPUTS: { userID: string - owner scope, date: models.Date - daily log date, items: []DailyNutritionSeedEntryInput - complete planned entries }
//	OUTPUTS: { *DailyNutritionSeedResult - created flag and committed/existing entry count }
//	SIDE_EFFECTS: Creates or reuses one daily log, locks it, inserts all entries, and commits only after the complete planned set is written.
//	LINKS: M-API-NUTRITION / V-M-API-NUTRITION.
//
// END_CONTRACT: SeedEntriesIfEmpty
func (r *dailyNutritionLogRepository) SeedEntriesIfEmpty(ctx context.Context, userID string, date models.Date, items []models.DailyNutritionSeedEntryInput) (*models.DailyNutritionSeedResult, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.SeedEntriesIfEmpty: %w", err)
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.SeedEntriesIfEmpty: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	qtx := r.q.WithTx(tx)
	logRow, err := qtx.CreateDailyNutritionLog(ctx, generated.CreateDailyNutritionLogParams{
		UserID: uid,
		Date:   modelsToPGDate(date),
		Notes:  nullableText(nil),
	})
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.SeedEntriesIfEmpty: %w", err)
	}

	var lockedID pgtype.UUID
	err = tx.QueryRow(ctx, `
SELECT id
FROM daily_nutrition_logs
WHERE id = $1 AND user_id = $2
FOR UPDATE
`, logRow.ID, uid).Scan(&lockedID)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.SeedEntriesIfEmpty: %w", err)
	}

	existing, err := qtx.ListDailyNutritionEntriesByLog(ctx, generated.ListDailyNutritionEntriesByLogParams{
		DailyLogID: logRow.ID,
		UserID:     uid,
	})
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.SeedEntriesIfEmpty: %w", err)
	}
	if len(existing) > 0 {
		if err := tx.Commit(ctx); err != nil {
			return nil, fmt.Errorf("daily_nutrition_log_repo.SeedEntriesIfEmpty: %w", err)
		}
		committed = true
		return &models.DailyNutritionSeedResult{Created: false, EntryCount: int32(len(existing))}, nil
	}

	for _, item := range items {
		productID, err := uuidFromString(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("daily_nutrition_log_repo.SeedEntriesIfEmpty: %w", err)
		}
		if _, err := qtx.CreateDailyNutritionEntry(ctx, generated.CreateDailyNutritionEntryParams{
			AmountGrams: float32(item.AmountGrams),
			MealLabel:   nullableText(item.MealLabel),
			Notes:       nullableText(item.Notes),
			Position:    item.Position,
			ProductID:   productID,
			DailyLogID:  logRow.ID,
			UserID:      uid,
		}); err != nil {
			return nil, fmt.Errorf("daily_nutrition_log_repo.SeedEntriesIfEmpty: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_repo.SeedEntriesIfEmpty: %w", err)
	}
	committed = true
	return &models.DailyNutritionSeedResult{Created: true, EntryCount: int32(len(items))}, nil
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
