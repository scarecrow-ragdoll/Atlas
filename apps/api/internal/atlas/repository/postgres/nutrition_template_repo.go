// FILE: apps/api/internal/atlas/repository/postgres/nutrition_template_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement NutritionTemplateRepository using sqlc-generated queries for the nutrition_template table.
//   SCOPE: Upsert, GetByID, GetByWeek, ListByRange, Update, Delete. All user-scoped. Uses pgtype.Date for weekStartDate. Not-found returns nil, nil.
//   DEPENDS: sqlc generated NutritionTemplate model, atlas/models for record types.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT

package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/repository/postgres/generated"
)

type NutritionTemplateRepository interface {
	Upsert(ctx context.Context, userID string, weekStartDate string, title, notes *string) (*models.NutritionTemplateRecord, error)
	GetByID(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error)
	GetByWeek(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error)
	ListByRange(ctx context.Context, userID string, startDate, endDate string) ([]models.NutritionTemplateRecord, error)
	Update(ctx context.Context, userID string, id string, title, notes *string) (*models.NutritionTemplateRecord, error)
	Delete(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error)
}

type nutritionTemplateRepository struct {
	q *generated.Queries
}

func NewNutritionTemplateRepository(pool *pgxpool.Pool) NutritionTemplateRepository {
	return &nutritionTemplateRepository{q: generated.New(pool)}
}

func (r *nutritionTemplateRepository) Upsert(ctx context.Context, userID string, weekStartDate string, title, notes *string) (*models.NutritionTemplateRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.Upsert: %w", err)
	}

	wd, err := parseDate(weekStartDate)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.Upsert: %w", err)
	}

	row, err := r.q.UpsertNutritionTemplate(ctx, generated.UpsertNutritionTemplateParams{
		UserID:        uid,
		WeekStartDate: wd,
		Title:         nullableText(title),
		Notes:         nullableText(notes),
	})
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.Upsert: %w", err)
	}

	return nutritionTemplateRecordFromRow(row), nil
}

func (r *nutritionTemplateRepository) GetByID(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error) {
	uid, tid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.GetByID: %w", err)
	}

	row, err := r.q.GetNutritionTemplateByID(ctx, generated.GetNutritionTemplateByIDParams{ID: tid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_template_repo.GetByID: %w", err)
	}

	return nutritionTemplateRecordFromRow(row), nil
}

func (r *nutritionTemplateRepository) GetByWeek(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.GetByWeek: %w", err)
	}

	wd, err := parseDate(weekStartDate)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.GetByWeek: %w", err)
	}

	row, err := r.q.GetNutritionTemplateByWeek(ctx, generated.GetNutritionTemplateByWeekParams{
		UserID:        uid,
		WeekStartDate: wd,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_template_repo.GetByWeek: %w", err)
	}

	return nutritionTemplateRecordFromRow(row), nil
}

func (r *nutritionTemplateRepository) ListByRange(ctx context.Context, userID string, startDate, endDate string) ([]models.NutritionTemplateRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.ListByRange: %w", err)
	}

	sd, err := parseDate(startDate)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.ListByRange: %w", err)
	}

	ed, err := parseDate(endDate)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.ListByRange: %w", err)
	}

	rows, err := r.q.ListNutritionTemplatesByRange(ctx, generated.ListNutritionTemplatesByRangeParams{
		UserID:          uid,
		WeekStartDate:   sd,
		WeekStartDate_2: ed,
	})
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.ListByRange: %w", err)
	}

	out := make([]models.NutritionTemplateRecord, len(rows))
	for i, row := range rows {
		out[i] = *nutritionTemplateRecordFromRow(row)
	}
	return out, nil
}

func (r *nutritionTemplateRepository) Update(ctx context.Context, userID string, id string, title, notes *string) (*models.NutritionTemplateRecord, error) {
	uid, tid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.Update: %w", err)
	}

	row, err := r.q.UpdateNutritionTemplate(ctx, generated.UpdateNutritionTemplateParams{
		ID:    tid,
		UserID: uid,
		Title: nullableText(title),
		Notes: nullableText(notes),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_template_repo.Update: %w", err)
	}

	return nutritionTemplateRecordFromRow(row), nil
}

func (r *nutritionTemplateRepository) Delete(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error) {
	uid, tid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_repo.Delete: %w", err)
	}

	row, err := r.q.DeleteNutritionTemplate(ctx, generated.DeleteNutritionTemplateParams{ID: tid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_template_repo.Delete: %w", err)
	}

	return nutritionTemplateRecordFromRow(row), nil
}

func parseDate(value string) (pgtype.Date, error) {
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return pgtype.Date{}, fmt.Errorf("parse date: %w", err)
	}
	return pgtype.Date{Time: t, Valid: true}, nil
}

func nutritionTemplateRecordFromRow(row generated.NutritionTemplate) *models.NutritionTemplateRecord {
	return &models.NutritionTemplateRecord{
		ID:            row.ID.String(),
		UserID:        row.UserID.String(),
		WeekStartDate: dateFromPGDate(row.WeekStartDate),
		Title:         textPtr(row.Title),
		Notes:         textPtr(row.Notes),
		CreatedAt:     formatTimestamp(row.CreatedAt),
		UpdatedAt:     formatTimestamp(row.UpdatedAt),
	}
}
