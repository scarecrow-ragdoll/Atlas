// FILE: apps/api/internal/atlas/repository/postgres/nutrition_product_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement NutritionProductRepository using sqlc-generated queries for the nutrition_product table.
//   SCOPE: Create, GetByID, ListActive, ListAll, Update, SoftDelete, Restore, GetByIDIncludeInactive. All user-scoped. Not-found returns nil, nil.
//   DEPENDS: sqlc generated NutritionProduct model and query functions, atlas/models for record types.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Exposed all-product listing and restore using generated sqlc queries.
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

type NutritionProductRepository interface {
	Create(ctx context.Context, userID string, name string, caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g float64, notes *string) (*models.NutritionProductRecord, error)
	GetByID(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error)
	ListActive(ctx context.Context, userID string) ([]models.NutritionProductRecord, error)
	ListAll(ctx context.Context, userID string) ([]models.NutritionProductRecord, error)
	Update(ctx context.Context, userID string, id string, name string, caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g float64, notes *string) (*models.NutritionProductRecord, error)
	SoftDelete(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error)
	Restore(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error)
	GetByIDIncludeInactive(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error)
}

type nutritionProductRepository struct {
	q *generated.Queries
}

func NewNutritionProductRepository(pool *pgxpool.Pool) NutritionProductRepository {
	return &nutritionProductRepository{q: generated.New(pool)}
}

func (r *nutritionProductRepository) Create(ctx context.Context, userID string, name string, caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g float64, notes *string) (*models.NutritionProductRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_repo.Create: %w", err)
	}

	row, err := r.q.CreateNutritionProduct(ctx, generated.CreateNutritionProductParams{
		UserID:          uid,
		Name:            name,
		CaloriesPer100g: float32(caloriesPer100g),
		ProteinPer100g:  float32(proteinPer100g),
		FatPer100g:      float32(fatPer100g),
		CarbsPer100g:    float32(carbsPer100g),
		Notes:           nullableText(notes),
	})
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_repo.Create: %w", err)
	}

	return nutritionProductRecordFromRow(row), nil
}

func (r *nutritionProductRepository) GetByID(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
	uid, pid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_repo.GetByID: %w", err)
	}

	row, err := r.q.GetNutritionProductByID(ctx, generated.GetNutritionProductByIDParams{ID: pid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_product_repo.GetByID: %w", err)
	}

	return nutritionProductRecordFromRow(row), nil
}

func (r *nutritionProductRepository) ListActive(ctx context.Context, userID string) ([]models.NutritionProductRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_repo.ListActive: %w", err)
	}

	rows, err := r.q.ListActiveNutritionProducts(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_repo.ListActive: %w", err)
	}

	out := make([]models.NutritionProductRecord, len(rows))
	for i, row := range rows {
		out[i] = *nutritionProductRecordFromRow(row)
	}
	return out, nil
}

func (r *nutritionProductRepository) ListAll(ctx context.Context, userID string) ([]models.NutritionProductRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_repo.ListAll: %w", err)
	}

	rows, err := r.q.ListNutritionProductsAll(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_repo.ListAll: %w", err)
	}

	out := make([]models.NutritionProductRecord, len(rows))
	for i, row := range rows {
		out[i] = *nutritionProductRecordFromRow(row)
	}
	return out, nil
}

func (r *nutritionProductRepository) Update(ctx context.Context, userID string, id string, name string, caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g float64, notes *string) (*models.NutritionProductRecord, error) {
	uid, pid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_repo.Update: %w", err)
	}

	row, err := r.q.UpdateNutritionProduct(ctx, generated.UpdateNutritionProductParams{
		ID:              pid,
		UserID:          uid,
		Name:            name,
		CaloriesPer100g: float32(caloriesPer100g),
		ProteinPer100g:  float32(proteinPer100g),
		FatPer100g:      float32(fatPer100g),
		CarbsPer100g:    float32(carbsPer100g),
		Notes:           nullableText(notes),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_product_repo.Update: %w", err)
	}

	return nutritionProductRecordFromRow(row), nil
}

func (r *nutritionProductRepository) SoftDelete(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
	uid, pid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_repo.SoftDelete: %w", err)
	}

	row, err := r.q.SoftDeleteNutritionProduct(ctx, generated.SoftDeleteNutritionProductParams{ID: pid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_product_repo.SoftDelete: %w", err)
	}

	return nutritionProductRecordFromRow(row), nil
}

func (r *nutritionProductRepository) Restore(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
	uid, pid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_repo.Restore: %w", err)
	}

	row, err := r.q.RestoreNutritionProduct(ctx, generated.RestoreNutritionProductParams{ID: pid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_product_repo.Restore: %w", err)
	}

	return nutritionProductRecordFromRow(row), nil
}

func (r *nutritionProductRepository) GetByIDIncludeInactive(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
	uid, pid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_repo.GetByIDIncludeInactive: %w", err)
	}

	row, err := r.q.GetNutritionProductByIDIncludeInactive(ctx, generated.GetNutritionProductByIDIncludeInactiveParams{ID: pid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("nutrition_product_repo.GetByIDIncludeInactive: %w", err)
	}

	return nutritionProductRecordFromRow(row), nil
}

func nutritionProductRecordFromRow(row generated.NutritionProduct) *models.NutritionProductRecord {
	return &models.NutritionProductRecord{
		ID:              row.ID.String(),
		UserID:          row.UserID.String(),
		Name:            row.Name,
		CaloriesPer100g: float64(row.CaloriesPer100g),
		ProteinPer100g:  float64(row.ProteinPer100g),
		FatPer100g:      float64(row.FatPer100g),
		CarbsPer100g:    float64(row.CarbsPer100g),
		Notes:           textPtr(row.Notes),
		IsActive:        row.IsActive,
		CreatedAt:       formatTimestamp(row.CreatedAt),
		UpdatedAt:       formatTimestamp(row.UpdatedAt),
	}
}
