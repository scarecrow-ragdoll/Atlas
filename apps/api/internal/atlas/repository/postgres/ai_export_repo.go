// FILE: apps/api/internal/atlas/repository/postgres/ai_export_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement AiExportRepository for WAVE-07 AI export management using sqlc-generated queries.
//   SCOPE: CRUD operations for AI exports, list by user ID, list stale exports, update file path.
//   DEPENDS: apps/api/internal/repository/postgres/generated, apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-07.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AiExportRepository - Interface for AI export data access.
//   NewAiExportRepository - Creates a new AiExportRepository.
//   Create - Creates an AI export record.
//   GetByID - Gets an AI export by ID (user-scoped).
//   ListByUserID - Lists AI exports by user ID.
//   UpdateFilePath - Updates the export file path.
//   Delete - Deletes an AI export.
//   ListStale - Lists stale AI exports older than the given interval.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added AI export repository for WAVE-07.
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

type AiExportRepository interface {
	Create(ctx context.Context, userID, dateRangeStart, dateRangeEnd string, includePhotos, includeNutrition, includeCardio, includeMeasurements bool, userComment *string, generatedPrompt string) (*models.AiExportRecord, error)
	GetByID(ctx context.Context, userID string, id string) (*models.AiExportRecord, error)
	ListByUserID(ctx context.Context, userID string) ([]models.AiExportRecord, error)
	UpdateFilePath(ctx context.Context, id string, filePath *string) (*models.AiExportRecord, error)
	Delete(ctx context.Context, userID string, id string) (*models.AiExportRecord, error)
	ListStale(ctx context.Context, interval string) ([]models.AiExportRecord, error)
}

type aiExportRepository struct {
	q *generated.Queries
}

func NewAiExportRepository(pool *pgxpool.Pool) AiExportRepository {
	return &aiExportRepository{q: generated.New(pool)}
}

func (r *aiExportRepository) Create(ctx context.Context, userID, dateRangeStart, dateRangeEnd string, includePhotos, includeNutrition, includeCardio, includeMeasurements bool, userComment *string, generatedPrompt string) (*models.AiExportRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("ai_export_repo.Create: %w", err)
	}

	startDate, err := models.ParseDate(dateRangeStart)
	if err != nil {
		return nil, fmt.Errorf("ai_export_repo.Create: parse dateRangeStart: %w", err)
	}

	endDate, err := models.ParseDate(dateRangeEnd)
	if err != nil {
		return nil, fmt.Errorf("ai_export_repo.Create: parse dateRangeEnd: %w", err)
	}

	row, err := r.q.CreateAiExport(ctx, generated.CreateAiExportParams{
		UserID:              uid,
		DateRangeStart:      modelsToPGDate(startDate),
		DateRangeEnd:        modelsToPGDate(endDate),
		IncludePhotos:       includePhotos,
		IncludeNutrition:    includeNutrition,
		IncludeCardio:       includeCardio,
		IncludeMeasurements: includeMeasurements,
		UserComment:         nullableText(userComment),
		GeneratedPrompt:     generatedPrompt,
		ExportFilePath:      pgtype.Text{},
	})
	if err != nil {
		return nil, fmt.Errorf("ai_export_repo.Create: %w", err)
	}

	return aiExportRecordFromRow(row), nil
}

func (r *aiExportRepository) GetByID(ctx context.Context, userID string, id string) (*models.AiExportRecord, error) {
	uid, eid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("ai_export_repo.GetByID: %w", err)
	}

	row, err := r.q.GetAiExportByID(ctx, generated.GetAiExportByIDParams{ID: eid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("ai_export_repo.GetByID: %w", err)
	}

	return aiExportRecordFromRow(row), nil
}

func (r *aiExportRepository) ListByUserID(ctx context.Context, userID string) ([]models.AiExportRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("ai_export_repo.ListByUserID: %w", err)
	}

	rows, err := r.q.ListAiExportsByUserID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("ai_export_repo.ListByUserID: %w", err)
	}

	out := make([]models.AiExportRecord, len(rows))
	for i, row := range rows {
		out[i] = *aiExportRecordFromRow(row)
	}
	return out, nil
}

func (r *aiExportRepository) UpdateFilePath(ctx context.Context, id string, filePath *string) (*models.AiExportRecord, error) {
	eid, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("ai_export_repo.UpdateFilePath: %w", err)
	}

	row, err := r.q.UpdateAiExportFilePath(ctx, generated.UpdateAiExportFilePathParams{
		ID:             eid,
		ExportFilePath: nullableText(filePath),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("ai_export_repo.UpdateFilePath: %w", err)
	}

	return aiExportRecordFromRow(row), nil
}

func (r *aiExportRepository) Delete(ctx context.Context, userID string, id string) (*models.AiExportRecord, error) {
	uid, eid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("ai_export_repo.Delete: %w", err)
	}

	row, err := r.q.DeleteAiExport(ctx, generated.DeleteAiExportParams{ID: eid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("ai_export_repo.Delete: %w", err)
	}

	return aiExportRecordFromRow(row), nil
}

func (r *aiExportRepository) ListStale(ctx context.Context, interval string) ([]models.AiExportRecord, error) {
	var intervalPg pgtype.Interval
	if err := intervalPg.Scan(interval); err != nil {
		return nil, fmt.Errorf("ai_export_repo.ListStale: parse interval: %w", err)
	}

	rows, err := r.q.ListStaleAiExports(ctx, intervalPg)
	if err != nil {
		return nil, fmt.Errorf("ai_export_repo.ListStale: %w", err)
	}

	out := make([]models.AiExportRecord, len(rows))
	for i, row := range rows {
		out[i] = *aiExportRecordFromRow(row)
	}
	return out, nil
}

func aiExportRecordFromRow(row generated.AiExport) *models.AiExportRecord {
	return &models.AiExportRecord{
		ID:                  row.ID.String(),
		UserID:              row.UserID.String(),
		DateRangeStart:      pgDateToModelsDate(row.DateRangeStart),
		DateRangeEnd:        pgDateToModelsDate(row.DateRangeEnd),
		IncludePhotos:       row.IncludePhotos,
		IncludeNutrition:    row.IncludeNutrition,
		IncludeCardio:       row.IncludeCardio,
		IncludeMeasurements: row.IncludeMeasurements,
		UserComment:         textPtr(row.UserComment),
		GeneratedPrompt:     row.GeneratedPrompt,
		ExportFilePath:      textPtr(row.ExportFilePath),
		CreatedAt:           formatTimestamp(row.CreatedAt),
		UpdatedAt:           formatTimestamp(row.UpdatedAt),
	}
}

func pgDateToModelsDate(d pgtype.Date) models.Date {
	if !d.Valid {
		return models.Date{}
	}
	return models.MustDate(d.Time.Format("2006-01-02"))
}