// FILE: apps/api/internal/atlas/service/nutrition_template_item_service.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement NutritionTemplateItemService with validation for amountGrams > 0 and template ownership check.
//   SCOPE: Create, Update, Delete. Create validates template exists and belongs to user. Update validates amountGrams > 0. Log markers: [NutritionTemplateItem][create|update|delete].
//   DEPENDS: postgres.NutritionTemplateItemRepository, postgres.NutritionTemplateRepository, models.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT

package service

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/atlas/repository/postgres"
)

var (
	ErrTemplateItemAmountInvalid = errors.New("amountGrams must be greater than 0")
)

type NutritionTemplateItemService interface {
	Create(ctx context.Context, userID string, input models.CreateTemplateItemInput) (*models.NutritionTemplateItem, error)
	Update(ctx context.Context, userID string, id string, input models.UpdateTemplateItemInput) (*models.NutritionTemplateItem, error)
	Delete(ctx context.Context, userID string, id string) (*models.NutritionTemplateItem, error)
}

type nutritionTemplateItemService struct {
	itemRepo postgres.NutritionTemplateItemRepository
	tmplRepo postgres.NutritionTemplateRepository
	logger   *zap.Logger
}

func NewNutritionTemplateItemService(itemRepo postgres.NutritionTemplateItemRepository, tmplRepo postgres.NutritionTemplateRepository, logger *zap.Logger) NutritionTemplateItemService {
	return &nutritionTemplateItemService{itemRepo: itemRepo, tmplRepo: tmplRepo, logger: logger}
}

func (s *nutritionTemplateItemService) Create(ctx context.Context, userID string, input models.CreateTemplateItemInput) (*models.NutritionTemplateItem, error) {
	s.logger.Info("[NutritionTemplateItem][create]")
	if input.AmountGrams <= 0 {
		return nil, ErrTemplateItemAmountInvalid
	}

	tmpl, err := s.tmplRepo.GetByID(ctx, userID, input.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_service.Create: %w", err)
	}
	if tmpl == nil {
		return nil, ErrTemplateNotFound
	}

	record, err := s.itemRepo.Create(ctx, input.TemplateID, input.ProductID, input.AmountGrams, input.MealLabel, input.Notes)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_service.Create: %w", err)
	}

	return models.NutritionTemplateItemFromRecord(record), nil
}

func (s *nutritionTemplateItemService) Update(ctx context.Context, userID string, id string, input models.UpdateTemplateItemInput) (*models.NutritionTemplateItem, error) {
	s.logger.Info("[NutritionTemplateItem][update]")
	existing, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_service.Update: %w", err)
	}
	if existing == nil {
		return nil, ErrTemplateItemNotFound
	}

	amount := existing.AmountGrams
	if input.AmountGrams != nil {
		if *input.AmountGrams <= 0 {
			return nil, ErrTemplateItemAmountInvalid
		}
		amount = *input.AmountGrams
	}

	mealLabel := input.MealLabel
	if mealLabel == nil {
		mealLabel = existing.MealLabel
	}

	notes := input.Notes
	if notes == nil {
		notes = existing.Notes
	}

	record, err := s.itemRepo.Update(ctx, id, amount, mealLabel, notes)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_service.Update: %w", err)
	}
	if record == nil {
		return nil, ErrTemplateItemNotFound
	}

	return models.NutritionTemplateItemFromRecord(record), nil
}

func (s *nutritionTemplateItemService) Delete(ctx context.Context, userID string, id string) (*models.NutritionTemplateItem, error) {
	s.logger.Info("[NutritionTemplateItem][delete]")
	record, err := s.itemRepo.Delete(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_item_service.Delete: %w", err)
	}
	if record == nil {
		return nil, ErrTemplateItemNotFound
	}
	return models.NutritionTemplateItemFromRecord(record), nil
}
