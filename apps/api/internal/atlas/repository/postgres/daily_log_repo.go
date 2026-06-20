// FILE: apps/api/internal/atlas/repository/postgres/daily_log_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement DailyLogRepository for WAVE-04 cardio FK support using sqlc-generated queries.
//   SCOPE: GetOrCreateDailyLogByDate, GetDailyLogByDate. Minimal — WAVE-03 adds full aggregate queries when deployed.
//   DEPENDS: apps/api/internal/repository/postgres/generated, apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   DailyLogRepository - Interface for DailyLog data access.
//   NewDailyLogRepository - Creates a new DailyLogRepository.
//   GetOrCreateDailyLogByDate - Gets existing or creates new DailyLog for user+date.
//   GetDailyLogByDate - Gets DailyLog by user+date.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added DailyLog repository for WAVE-04.
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

type DailyLogRepository interface {
	GetDailyLogByDate(ctx context.Context, userID string, date models.Date) (*models.DailyLogRecord, error)
	GetOrCreateDailyLogByDate(ctx context.Context, userID string, date models.Date) (*models.DailyLogRecord, error)
}

type dailyLogRepository struct {
	q *generated.Queries
}

func NewDailyLogRepository(pool *pgxpool.Pool) DailyLogRepository {
	return &dailyLogRepository{q: generated.New(pool)}
}

func (r *dailyLogRepository) GetDailyLogByDate(ctx context.Context, userID string, date models.Date) (*models.DailyLogRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("daily_log_repo.GetDailyLogByDate: %w", err)
	}

	row, err := r.q.GetDailyLogByDate(ctx, generated.GetDailyLogByDateParams{
		UserID: uid,
		Date:   pgtype.Date{Time: date.Time(), Valid: !date.Time().IsZero()},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("daily_log_repo.GetDailyLogByDate: %w", err)
	}

	return dailyLogRecordFromRow(row), nil
}

func (r *dailyLogRepository) GetOrCreateDailyLogByDate(ctx context.Context, userID string, date models.Date) (*models.DailyLogRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("daily_log_repo.GetOrCreateDailyLogByDate: %w", err)
	}

	row, err := r.q.CreateDailyLog(ctx, generated.CreateDailyLogParams{
		UserID: uid,
		Date:   pgtype.Date{Time: date.Time(), Valid: !date.Time().IsZero()},
		Notes:  pgtype.Text{},
	})
	if err != nil {
		return nil, fmt.Errorf("daily_log_repo.GetOrCreateDailyLogByDate: %w", err)
	}

	return dailyLogRecordFromRow(row), nil
}

func dailyLogRecordFromRow(row generated.DailyLog) *models.DailyLogRecord {
	return &models.DailyLogRecord{
		ID:        row.ID.String(),
		UserID:    row.UserID.String(),
		Date:      dateFromPGDate(row.Date),
		Notes:     textPtr(row.Notes),
		CreatedAt: formatTimestamp(row.CreatedAt),
		UpdatedAt: formatTimestamp(row.UpdatedAt),
	}
}

func dateFromPGDate(d pgtype.Date) models.Date {
	if !d.Valid {
		return models.Date{}
	}
	return models.MustDate(d.Time.Format("2006-01-02"))
}