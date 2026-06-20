// FILE: apps/api/internal/atlas/repository/postgres/nutrition_override_item_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement DailyNutritionOverrideItemRepository using sqlc-generated queries for the daily_nutrition_override_item table.
//   SCOPE: Create, GetByID, ListByOverride, Update, Delete. Override-scoped (not user-scoped for item CRUD). Not-found returns nil, nil.
//   DEPENDS: sqlc generated DailyNutritionOverrideItem model, atlas/models for record types.
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

type DailyNutritionOverrideItemRepository interface {
	Create(ctx context.Context, overrideID string, productID string, amountGrams float64, operation string, mealLabel, notes *string) (*models.DailyNutritionOverrideItemRecord, error)
	GetByID(ctx context.Context, id string) (*models.DailyNutritionOverrideItemRecord, error)
	ListByOverride(ctx context.Context, overrideID string) ([]models.DailyNutritionOverrideItemRecord, error)
	Update(ctx context.Context, id string, amountGrams float64, operation string, mealLabel, notes *string) (*models.DailyNutritionOverrideItemRecord, error)
	Delete(ctx context.Context, id string) (*models.DailyNutritionOverrideItemRecord, error)
}

type dailyNutritionOverrideItemRepository struct {
	q *generated.Queries
}

func NewDailyNutritionOverrideItemRepository(pool *pgxpool.Pool) DailyNutritionOverrideItemRepository {
	return &dailyNutritionOverrideItemRepository{q: generated.New(pool)}
}

func (r *dailyNutritionOverrideItemRepository) Create(ctx context.Context, overrideID string, productID string, amountGrams float64, operation string, mealLabel, notes *string) (*models.DailyNutritionOverrideItemRecord, error) {
	oid, err := uuidFromString(overrideID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_item_repo.Create: %w", err)
	}
	pid, err := uuidFromString(productID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_item_repo.Create: %w", err)
	}

	row, err := r.q.CreateDailyNutritionOverrideItem(ctx, generated.CreateDailyNutritionOverrideItemParams{
		OverrideID:  oid,
		ProductID:   pid,
		AmountGrams: float32(amountGrams),
		Operation:   operation,
		MealLabel:   nullableText(mealLabel),
		Notes:       nullableText(notes),
	})
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_item_repo.Create: %w", err)
	}

	return dailyNutritionOverrideItemRecordFromRow(row), nil
}

func (r *dailyNutritionOverrideItemRepository) GetByID(ctx context.Context, id string) (*models.DailyNutritionOverrideItemRecord, error) {
	iid, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_item_repo.GetByID: %w", err)
	}

	row, err := r.q.GetDailyNutritionOverrideItemByID(ctx, iid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_override_item_repo.GetByID: %w", err)
	}

	return dailyNutritionOverrideItemRecordFromRow(row), nil
}

func (r *dailyNutritionOverrideItemRepository) ListByOverride(ctx context.Context, overrideID string) ([]models.DailyNutritionOverrideItemRecord, error) {
	oid, err := uuidFromString(overrideID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_item_repo.ListByOverride: %w", err)
	}

	rows, err := r.q.ListDailyNutritionOverrideItemsByOverride(ctx, oid)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_item_repo.ListByOverride: %w", err)
	}

	out := make([]models.DailyNutritionOverrideItemRecord, len(rows))
	for i, row := range rows {
		out[i] = *dailyNutritionOverrideItemRecordFromRow(row)
	}
	return out, nil
}

func (r *dailyNutritionOverrideItemRepository) Update(ctx context.Context, id string, amountGrams float64, operation string, mealLabel, notes *string) (*models.DailyNutritionOverrideItemRecord, error) {
	iid, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_item_repo.Update: %w", err)
	}

	row, err := r.q.UpdateDailyNutritionOverrideItem(ctx, generated.UpdateDailyNutritionOverrideItemParams{
		ID:          iid,
		AmountGrams: float32(amountGrams),
		Operation:   operation,
		MealLabel:   nullableText(mealLabel),
		Notes:       nullableText(notes),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_override_item_repo.Update: %w", err)
	}

	return dailyNutritionOverrideItemRecordFromRow(row), nil
}

func (r *dailyNutritionOverrideItemRepository) Delete(ctx context.Context, id string) (*models.DailyNutritionOverrideItemRecord, error) {
	iid, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_item_repo.Delete: %w", err)
	}

	row, err := r.q.DeleteDailyNutritionOverrideItem(ctx, iid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_override_item_repo.Delete: %w", err)
	}

	return dailyNutritionOverrideItemRecordFromRow(row), nil
}

func dailyNutritionOverrideItemRecordFromRow(row generated.DailyNutritionOverrideItem) *models.DailyNutritionOverrideItemRecord {
	return &models.DailyNutritionOverrideItemRecord{
		ID:          row.ID.String(),
		OverrideID:  row.OverrideID.String(),
		ProductID:   row.ProductID.String(),
		AmountGrams: float64(row.AmountGrams),
		Operation:   row.Operation,
		MealLabel:   textPtr(row.MealLabel),
		Notes:       textPtr(row.Notes),
		CreatedAt:   formatTimestamp(row.CreatedAt),
		UpdatedAt:   formatTimestamp(row.UpdatedAt),
	}
}
