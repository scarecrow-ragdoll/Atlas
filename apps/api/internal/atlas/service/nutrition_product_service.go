// FILE: apps/api/internal/atlas/service/nutrition_product_service.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement NutritionProductService with validation for product name and macro values, soft-delete, and log markers.
//   SCOPE: Create, GetByID (includes inactive), ListActive, Update, Delete (soft-delete). Validation: name required (<=255 chars), macros >= 0. Log markers: [NutritionProduct][create|get|list|update|delete].
//   DEPENDS: postgres.NutritionProductRepository, models.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT

package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/atlas/repository/postgres"
)

var (
	ErrProductNameRequired  = errors.New("product name is required")
	ErrProductMacroNegative = errors.New("nutritional values must be >= 0")
	ErrProductNotFound      = errors.New("nutrition product not found")
	ErrProductNameTooLong   = errors.New("product name must not exceed 255 characters")
)

type NutritionProductService interface {
	Create(ctx context.Context, userID string, input models.CreateProductInput) (*models.NutritionProduct, error)
	GetByID(ctx context.Context, userID string, id string) (*models.NutritionProduct, error)
	ListActive(ctx context.Context, userID string) ([]models.NutritionProduct, error)
	Update(ctx context.Context, userID string, id string, input models.UpdateProductInput) (*models.NutritionProduct, error)
	Delete(ctx context.Context, userID string, id string) (*models.NutritionProduct, error)
}

type nutritionProductService struct {
	repo   postgres.NutritionProductRepository
	logger *zap.Logger
}

func NewNutritionProductService(repo postgres.NutritionProductRepository, logger *zap.Logger) NutritionProductService {
	return &nutritionProductService{repo: repo, logger: logger}
}

func (s *nutritionProductService) Create(ctx context.Context, userID string, input models.CreateProductInput) (*models.NutritionProduct, error) {
	s.logger.Info("[NutritionProduct][create]")
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrProductNameRequired
	}
	if len(name) > 255 {
		return nil, ErrProductNameTooLong
	}
	if input.CaloriesPer100g < 0 || input.ProteinPer100g < 0 || input.FatPer100g < 0 || input.CarbsPer100g < 0 {
		return nil, ErrProductMacroNegative
	}

	record, err := s.repo.Create(ctx, userID, name, input.CaloriesPer100g, input.ProteinPer100g, input.FatPer100g, input.CarbsPer100g, input.Notes)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_service.Create: %w", err)
	}

	return models.NutritionProductFromRecord(record), nil
}

func (s *nutritionProductService) GetByID(ctx context.Context, userID string, id string) (*models.NutritionProduct, error) {
	s.logger.Info("[NutritionProduct][get]")
	record, err := s.repo.GetByIDIncludeInactive(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_service.GetByID: %w", err)
	}
	if record == nil {
		return nil, ErrProductNotFound
	}
	return models.NutritionProductFromRecord(record), nil
}

func (s *nutritionProductService) ListActive(ctx context.Context, userID string) ([]models.NutritionProduct, error) {
	records, err := s.repo.ListActive(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_service.ListActive: %w", err)
	}

	out := make([]models.NutritionProduct, len(records))
	for i := range records {
		out[i] = *models.NutritionProductFromRecord(&records[i])
	}
	return out, nil
}

func (s *nutritionProductService) Update(ctx context.Context, userID string, id string, input models.UpdateProductInput) (*models.NutritionProduct, error) {
	s.logger.Info("[NutritionProduct][update]")
	existing, err := s.repo.GetByIDIncludeInactive(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_service.Update: %w", err)
	}
	if existing == nil {
		return nil, ErrProductNotFound
	}

	name := existing.Name
	if input.Name != nil {
		name = strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, ErrProductNameRequired
		}
		if len(name) > 255 {
			return nil, ErrProductNameTooLong
		}
	}

	calories := existing.CaloriesPer100g
	if input.CaloriesPer100g != nil {
		calories = *input.CaloriesPer100g
	}
	protein := existing.ProteinPer100g
	if input.ProteinPer100g != nil {
		protein = *input.ProteinPer100g
	}
	fat := existing.FatPer100g
	if input.FatPer100g != nil {
		fat = *input.FatPer100g
	}
	carbs := existing.CarbsPer100g
	if input.CarbsPer100g != nil {
		carbs = *input.CarbsPer100g
	}

	if calories < 0 || protein < 0 || fat < 0 || carbs < 0 {
		return nil, ErrProductMacroNegative
	}

	notes := input.Notes
	if notes == nil {
		notes = existing.Notes
	}

	record, err := s.repo.Update(ctx, userID, id, name, calories, protein, fat, carbs, notes)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_service.Update: %w", err)
	}
	if record == nil {
		return nil, ErrProductNotFound
	}

	return models.NutritionProductFromRecord(record), nil
}

func (s *nutritionProductService) Delete(ctx context.Context, userID string, id string) (*models.NutritionProduct, error) {
	s.logger.Info("[NutritionProduct][delete]")
	record, err := s.repo.SoftDelete(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_product_service.Delete: %w", err)
	}
	if record == nil {
		return nil, ErrProductNotFound
	}
	return models.NutritionProductFromRecord(record), nil
}
