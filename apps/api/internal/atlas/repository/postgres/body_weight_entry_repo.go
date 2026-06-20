// FILE: apps/api/internal/atlas/repository/postgres/body_weight_entry_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement BodyWeightEntryRepository for WAVE-04 body weight tracking using sqlc-generated queries.
//   SCOPE: CRUD operations for body weight entries, list by date range, latest entry, all user-scoped.
//   DEPENDS: apps/api/internal/repository/postgres/generated, apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   BodyWeightEntryRepository - Interface for body weight entry data access.
//   NewBodyWeightEntryRepository - Creates a new BodyWeightEntryRepository.
//   Create - Creates a body weight entry.
//   GetByID - Gets a body weight entry by ID (user-scoped).
//   ListByDateRange - Lists body weight entries by date range.
//   Latest - Gets the latest body weight entry.
//   Update - Updates a body weight entry.
//   Delete - Deletes a body weight entry.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added body weight entry repository for WAVE-04.
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

type BodyWeightEntryRepository interface {
	Create(ctx context.Context, userID string, date models.Date, weight float64, source string, notes *string) (*models.BodyWeightRecord, error)
	GetByID(ctx context.Context, userID string, id string) (*models.BodyWeightRecord, error)
	ListByDateRange(ctx context.Context, userID string, fromDate models.Date, toDate models.Date) ([]models.BodyWeightRecord, error)
	Latest(ctx context.Context, userID string) (*models.BodyWeightRecord, error)
	Update(ctx context.Context, userID string, id string, weight *float64, source *string, notes *string) (*models.BodyWeightRecord, error)
	Delete(ctx context.Context, userID string, id string) (*models.BodyWeightRecord, error)
}

type bodyWeightEntryRepository struct {
	q *generated.Queries
}

func NewBodyWeightEntryRepository(pool *pgxpool.Pool) BodyWeightEntryRepository {
	return &bodyWeightEntryRepository{q: generated.New(pool)}
}

func (r *bodyWeightEntryRepository) Create(ctx context.Context, userID string, date models.Date, weight float64, source string, notes *string) (*models.BodyWeightRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("body_weight_repo.Create: %w", err)
	}

	row, err := r.q.CreateBodyWeightEntry(ctx, generated.CreateBodyWeightEntryParams{
		UserID: uid,
		Date:   modelsToPGDate(date),
		Weight: float32(weight),
		Source: source,
		Notes:  nullableText(notes),
	})
	if err != nil {
		return nil, fmt.Errorf("body_weight_repo.Create: %w", err)
	}

	return bodyWeightRecordFromRow(row), nil
}

func (r *bodyWeightEntryRepository) GetByID(ctx context.Context, userID string, id string) (*models.BodyWeightRecord, error) {
	uid, eid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("body_weight_repo.GetByID: %w", err)
	}

	row, err := r.q.GetBodyWeightEntryByID(ctx, generated.GetBodyWeightEntryByIDParams{ID: eid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("body_weight_repo.GetByID: %w", err)
	}

	return bodyWeightRecordFromRow(row), nil
}

func (r *bodyWeightEntryRepository) ListByDateRange(ctx context.Context, userID string, fromDate models.Date, toDate models.Date) ([]models.BodyWeightRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("body_weight_repo.ListByDateRange: %w", err)
	}

	rows, err := r.q.ListBodyWeightEntriesByDateRange(ctx, generated.ListBodyWeightEntriesByDateRangeParams{
		UserID: uid,
		Date:   modelsToPGDate(fromDate),
		Date_2: modelsToPGDate(toDate),
	})
	if err != nil {
		return nil, fmt.Errorf("body_weight_repo.ListByDateRange: %w", err)
	}

	out := make([]models.BodyWeightRecord, len(rows))
	for i, row := range rows {
		out[i] = *bodyWeightRecordFromRow(row)
	}
	return out, nil
}

func (r *bodyWeightEntryRepository) Latest(ctx context.Context, userID string) (*models.BodyWeightRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("body_weight_repo.Latest: %w", err)
	}

	row, err := r.q.LatestBodyWeightEntry(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("body_weight_repo.Latest: %w", err)
	}

	return bodyWeightRecordFromRow(row), nil
}

func (r *bodyWeightEntryRepository) Update(ctx context.Context, userID string, id string, weight *float64, source *string, notes *string) (*models.BodyWeightRecord, error) {
	uid, eid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("body_weight_repo.Update: %w", err)
	}

	w := float32(0)
	if weight != nil {
		w = float32(*weight)
	}
	s := ""
	if source != nil {
		s = *source
	}

	row, err := r.q.UpdateBodyWeightEntry(ctx, generated.UpdateBodyWeightEntryParams{
		ID:     eid,
		UserID: uid,
		Weight: w,
		Source: s,
		Notes:  nullableText(notes),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("body_weight_repo.Update: %w", err)
	}

	return bodyWeightRecordFromRow(row), nil
}

func (r *bodyWeightEntryRepository) Delete(ctx context.Context, userID string, id string) (*models.BodyWeightRecord, error) {
	uid, eid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("body_weight_repo.Delete: %w", err)
	}

	row, err := r.q.DeleteBodyWeightEntry(ctx, generated.DeleteBodyWeightEntryParams{ID: eid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("body_weight_repo.Delete: %w", err)
	}

	return bodyWeightRecordFromRow(row), nil
}

func bodyWeightRecordFromRow(row generated.BodyWeightEntry) *models.BodyWeightRecord {
	return &models.BodyWeightRecord{
		ID:        row.ID.String(),
		UserID:    row.UserID.String(),
		Date:      dateFromPGDate(row.Date),
		Weight:    float64(row.Weight),
		Source:    row.Source,
		Notes:     textPtr(row.Notes),
		CreatedAt: formatTimestamp(row.CreatedAt),
		UpdatedAt: formatTimestamp(row.UpdatedAt),
	}
}

func modelsToPGDate(date models.Date) pgtype.Date {
	return pgtype.Date{Time: date.Time(), Valid: !date.Time().IsZero()}
}