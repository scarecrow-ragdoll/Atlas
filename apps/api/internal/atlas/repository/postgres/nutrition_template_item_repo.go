// FILE: apps/api/internal/atlas/repository/postgres/nutrition_template_item_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement NutritionTemplateItemRepository using sqlc-generated queries for the nutrition_template_item table.
//   SCOPE: Create, GetByID, ListByTemplate, Update, Delete. Template-scoped (not user-scoped for item CRUD). Not-found returns nil, nil.
//   DEPENDS: sqlc generated NutritionTemplateItem model, atlas/models for record types.
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

type NutritionTemplateItemRepository interface {
	Create(ctx context.Context, templateID string, productID string, amountGrams float64, mealLabel, notes *string) (*models.NutritionTemplateItemRecord, error)
	GetByID(ctx context.Context, id string) (*models.NutritionTemplateItemRecord, error)
	ListByTemplate(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error)
	Update(ctx context.Context, id string, amountGrams float64, mealLabel, notes *string) (*models.NutritionTemplateItemRecord, error)
	Delete(ctx context.Context, id string) (*models.NutritionTemplateItemRecord, error)
}

type nutritionTemplateItemRepository struct {
	q *generated.Queries
}

func NewNutritionTemplateItemRepository(pool *pgxpool.Pool) NutritionTemplateItemRepository {
	return &nutritionTemplateItemRepository{q: generated.New(pool)}
}

func (r *nutritionTemplateItemRepository) Create(ctx context.Context, templateID string, productID string, amountGrams float64, mealLabel, notes *string) (*models.NutritionTemplateItemRecord, error) {
	tid, err := uuidFromString(templateID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_repo.Create: %w", err)
	}
	pid, err := uuidFromString(productID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_repo.Create: %w", err)
	}

	row, err := r.q.CreateNutritionTemplateItem(ctx, generated.CreateNutritionTemplateItemParams{
		TemplateID:  tid,
		ProductID:   pid,
		AmountGrams: float32(amountGrams),
		MealLabel:   nullableText(mealLabel),
		Notes:       nullableText(notes),
	})
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_repo.Create: %w", err)
	}

	return nutritionTemplateItemRecordFromRow(row), nil
}

func (r *nutritionTemplateItemRepository) GetByID(ctx context.Context, id string) (*models.NutritionTemplateItemRecord, error) {
	iid, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_repo.GetByID: %w", err)
	}

	row, err := r.q.GetNutritionTemplateItemByID(ctx, iid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_template_item_repo.GetByID: %w", err)
	}

	return nutritionTemplateItemRecordFromRow(row), nil
}

func (r *nutritionTemplateItemRepository) ListByTemplate(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error) {
	tid, err := uuidFromString(templateID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_repo.ListByTemplate: %w", err)
	}

	rows, err := r.q.ListNutritionTemplateItemsByTemplate(ctx, tid)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_repo.ListByTemplate: %w", err)
	}

	out := make([]models.NutritionTemplateItemRecord, len(rows))
	for i, row := range rows {
		out[i] = *nutritionTemplateItemRecordFromRow(row)
	}
	return out, nil
}

func (r *nutritionTemplateItemRepository) Update(ctx context.Context, id string, amountGrams float64, mealLabel, notes *string) (*models.NutritionTemplateItemRecord, error) {
	iid, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_repo.Update: %w", err)
	}

	row, err := r.q.UpdateNutritionTemplateItem(ctx, generated.UpdateNutritionTemplateItemParams{
		ID:          iid,
		AmountGrams: float32(amountGrams),
		MealLabel:   nullableText(mealLabel),
		Notes:       nullableText(notes),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_template_item_repo.Update: %w", err)
	}

	return nutritionTemplateItemRecordFromRow(row), nil
}

func (r *nutritionTemplateItemRepository) Delete(ctx context.Context, id string) (*models.NutritionTemplateItemRecord, error) {
	iid, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_repo.Delete: %w", err)
	}

	row, err := r.q.DeleteNutritionTemplateItem(ctx, iid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_template_item_repo.Delete: %w", err)
	}

	return nutritionTemplateItemRecordFromRow(row), nil
}

func nutritionTemplateItemRecordFromRow(row generated.NutritionTemplateItem) *models.NutritionTemplateItemRecord {
	return &models.NutritionTemplateItemRecord{
		ID:          row.ID.String(),
		TemplateID:  row.TemplateID.String(),
		ProductID:   row.ProductID.String(),
		AmountGrams: float64(row.AmountGrams),
		MealLabel:   textPtr(row.MealLabel),
		Notes:       textPtr(row.Notes),
		CreatedAt:   formatTimestamp(row.CreatedAt),
		UpdatedAt:   formatTimestamp(row.UpdatedAt),
	}
}
