// FILE: apps/api/internal/atlas/repository/postgres/cardio_entry_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement CardioEntryRepository for WAVE-04 cardio tracking using sqlc-generated queries.
//   SCOPE: CRUD operations for cardio entries, list by daily log ID, all user-scoped.
//   DEPENDS: apps/api/internal/repository/postgres/generated, apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   CardioEntryRepository - Interface for cardio entry data access.
//   NewCardioEntryRepository - Creates a new CardioEntryRepository.
//   Create - Creates a cardio entry.
//   GetByID - Gets a cardio entry by ID (user-scoped).
//   ListByDailyLog - Lists cardio entries by daily log ID.
//   Update - Updates a cardio entry.
//   Delete - Deletes a cardio entry.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added cardio entry repository for WAVE-04.
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

type CardioEntryRepository interface {
	Create(ctx context.Context, userID string, dailyLogID string, cardioType string, durationMinutes int32, avgPulse *int32, heartRateZone *string, notes *string) (*models.CardioRecord, error)
	GetByID(ctx context.Context, userID string, id string) (*models.CardioRecord, error)
	ListByDailyLog(ctx context.Context, userID string, dailyLogID string) ([]models.CardioRecord, error)
	Update(ctx context.Context, userID string, id string, cardioType string, durationMinutes int32, avgPulse *int32, heartRateZone *string, notes *string) (*models.CardioRecord, error)
	Delete(ctx context.Context, userID string, id string) (*models.CardioRecord, error)
}

type cardioEntryRepository struct {
	q *generated.Queries
}

func NewCardioEntryRepository(pool *pgxpool.Pool) CardioEntryRepository {
	return &cardioEntryRepository{q: generated.New(pool)}
}

func (r *cardioEntryRepository) Create(ctx context.Context, userID string, dailyLogID string, cardioType string, durationMinutes int32, avgPulse *int32, heartRateZone *string, notes *string) (*models.CardioRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("cardio_repo.Create: %w", err)
	}
	dlid, err := uuidFromString(dailyLogID)
	if err != nil {
		return nil, fmt.Errorf("cardio_repo.Create: %w", err)
	}

	row, err := r.q.CreateCardioEntry(ctx, generated.CreateCardioEntryParams{
		UserID:          uid,
		DailyLogID:      dlid,
		CardioType:      cardioType,
		DurationMinutes: durationMinutes,
		AvgPulse:        nullableInt4(avgPulse),
		HeartRateZone:   nullableText(heartRateZone),
		Notes:           nullableText(notes),
	})
	if err != nil {
		return nil, fmt.Errorf("cardio_repo.Create: %w", err)
	}

	return cardioRecordFromRow(row), nil
}

func (r *cardioEntryRepository) GetByID(ctx context.Context, userID string, id string) (*models.CardioRecord, error) {
	uid, eid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("cardio_repo.GetByID: %w", err)
	}

	row, err := r.q.GetCardioEntryByID(ctx, generated.GetCardioEntryByIDParams{ID: eid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("cardio_repo.GetByID: %w", err)
	}

	return cardioRecordFromRow(row), nil
}

func (r *cardioEntryRepository) ListByDailyLog(ctx context.Context, userID string, dailyLogID string) ([]models.CardioRecord, error) {
	uid, dlid, err := parseTwoUUIDs(userID, dailyLogID)
	if err != nil {
		return nil, fmt.Errorf("cardio_repo.ListByDailyLog: %w", err)
	}

	rows, err := r.q.ListCardioEntriesByDailyLog(ctx, generated.ListCardioEntriesByDailyLogParams{
		DailyLogID: dlid,
		UserID:     uid,
	})
	if err != nil {
		return nil, fmt.Errorf("cardio_repo.ListByDailyLog: %w", err)
	}

	out := make([]models.CardioRecord, len(rows))
	for i, row := range rows {
		out[i] = *cardioRecordFromRow(row)
	}
	return out, nil
}

func (r *cardioEntryRepository) Update(ctx context.Context, userID string, id string, cardioType string, durationMinutes int32, avgPulse *int32, heartRateZone *string, notes *string) (*models.CardioRecord, error) {
	uid, eid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("cardio_repo.Update: %w", err)
	}

	row, err := r.q.UpdateCardioEntry(ctx, generated.UpdateCardioEntryParams{
		ID:              eid,
		UserID:          uid,
		CardioType:      cardioType,
		DurationMinutes: durationMinutes,
		AvgPulse:        nullableInt4(avgPulse),
		HeartRateZone:   nullableText(heartRateZone),
		Notes:           nullableText(notes),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("cardio_repo.Update: %w", err)
	}

	return cardioRecordFromRow(row), nil
}

func (r *cardioEntryRepository) Delete(ctx context.Context, userID string, id string) (*models.CardioRecord, error) {
	uid, eid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("cardio_repo.Delete: %w", err)
	}

	row, err := r.q.DeleteCardioEntry(ctx, generated.DeleteCardioEntryParams{ID: eid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("cardio_repo.Delete: %w", err)
	}

	return cardioRecordFromRow(row), nil
}

func cardioRecordFromRow(row generated.CardioEntry) *models.CardioRecord {
	return &models.CardioRecord{
		ID:              row.ID.String(),
		UserID:          row.UserID.String(),
		DailyLogID:      row.DailyLogID.String(),
		CardioType:      row.CardioType,
		DurationMinutes: row.DurationMinutes,
		AvgPulse:        int4Ptr(row.AvgPulse),
		HeartRateZone:   textPtr(row.HeartRateZone),
		Notes:           textPtr(row.Notes),
		CreatedAt:       formatTimestamp(row.CreatedAt),
		UpdatedAt:       formatTimestamp(row.UpdatedAt),
	}
}