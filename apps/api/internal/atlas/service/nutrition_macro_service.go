// FILE: apps/api/internal/atlas/service/nutrition_macro_service.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement NutritionMacroService for KJBJU calculation per RULE-010/RULE-011. Calculates per-day macros from template items with override operations (add/subtract/replace).
//   SCOPE: Calculate - takes userID, weekStartDate, optional date. Returns NutritionMacros (calories/protein/fat/carbs). Handles: empty template (0 macros), soft-deleted products (0 contribution), override ADD/SUBTRACT/REPLACE operations. Log markers: [NutritionMacros][calculate].
//   DEPENDS: All 5 nutrition repository interfaces.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT

package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/atlas/repository/postgres"
)

type NutritionMacroService interface {
	Calculate(ctx context.Context, userID string, weekStartDate string, date string) (*models.NutritionMacros, error)
}

type nutritionMacroService struct {
	tmplRepo         postgres.NutritionTemplateRepository
	itemRepo         postgres.NutritionTemplateItemRepository
	overrideRepo     postgres.DailyNutritionOverrideRepository
	overrideItemRepo postgres.DailyNutritionOverrideItemRepository
	productRepo      postgres.NutritionProductRepository
	logger           *zap.Logger
}

func NewNutritionMacroService(
	tmplRepo postgres.NutritionTemplateRepository,
	itemRepo postgres.NutritionTemplateItemRepository,
	overrideRepo postgres.DailyNutritionOverrideRepository,
	overrideItemRepo postgres.DailyNutritionOverrideItemRepository,
	productRepo postgres.NutritionProductRepository,
	logger *zap.Logger,
) NutritionMacroService {
	return &nutritionMacroService{
		tmplRepo:         tmplRepo,
		itemRepo:         itemRepo,
		overrideRepo:     overrideRepo,
		overrideItemRepo: overrideItemRepo,
		productRepo:      productRepo,
		logger:           logger,
	}
}

func (s *nutritionMacroService) Calculate(ctx context.Context, userID string, weekStartDate string, date string) (*models.NutritionMacros, error) {
	s.logger.Info("[NutritionMacros][calculate]")
	tmpl, err := s.tmplRepo.GetByWeek(ctx, userID, weekStartDate)
	if err != nil {
		return nil, fmt.Errorf("nutrition_macro_service.Calculate: %w", err)
	}
	if tmpl == nil {
		return &models.NutritionMacros{}, nil
	}

	items, err := s.itemRepo.ListByTemplate(ctx, tmpl.ID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_macro_service.Calculate: %w", err)
	}

	result := &models.NutritionMacros{}

	type templateMacro struct {
		calories, protein, fat, carbs float64
	}
	tmplMacrosByProduct := make(map[string]templateMacro)

	for _, item := range items {
		product, err := s.productRepo.GetByIDIncludeInactive(ctx, userID, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("nutrition_macro_service.Calculate: %w", err)
		}
		if product == nil || !product.IsActive {
			continue
		}

		factor := item.AmountGrams / 100.0
		m := templateMacro{
			calories: product.CaloriesPer100g * factor,
			protein:  product.ProteinPer100g * factor,
			fat:      product.FatPer100g * factor,
			carbs:    product.CarbsPer100g * factor,
		}
		result.Calories += m.calories
		result.Protein += m.protein
		result.Fat += m.fat
		result.Carbs += m.carbs

		prev := tmplMacrosByProduct[item.ProductID]
		prev.calories += m.calories
		prev.protein += m.protein
		prev.fat += m.fat
		prev.carbs += m.carbs
		tmplMacrosByProduct[item.ProductID] = prev
	}

	if date == "" {
		return result, nil
	}

	override, err := s.overrideRepo.GetByDate(ctx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("nutrition_macro_service.Calculate: %w", err)
	}
	if override == nil {
		return result, nil
	}

	overrideItems, err := s.overrideItemRepo.ListByOverride(ctx, override.ID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_macro_service.Calculate: %w", err)
	}

	for _, oi := range overrideItems {
		product, err := s.productRepo.GetByIDIncludeInactive(ctx, userID, oi.ProductID)
		if err != nil {
			return nil, fmt.Errorf("nutrition_macro_service.Calculate: %w", err)
		}
		if product == nil || !product.IsActive {
			continue
		}

		factor := oi.AmountGrams / 100.0
		overrideCal := product.CaloriesPer100g * factor
		overrideProtein := product.ProteinPer100g * factor
		overrideFat := product.FatPer100g * factor
		overrideCarbs := product.CarbsPer100g * factor

		switch oi.Operation {
		case string(models.OperationAdd):
			result.Calories += overrideCal
			result.Protein += overrideProtein
			result.Fat += overrideFat
			result.Carbs += overrideCarbs
		case string(models.OperationSubtract):
			result.Calories -= overrideCal
			result.Protein -= overrideProtein
			result.Fat -= overrideFat
			result.Carbs -= overrideCarbs
		case string(models.OperationReplace):
			tmpl, exists := tmplMacrosByProduct[oi.ProductID]
			if exists {
				result.Calories = result.Calories - tmpl.calories + overrideCal
				result.Protein = result.Protein - tmpl.protein + overrideProtein
				result.Fat = result.Fat - tmpl.fat + overrideFat
				result.Carbs = result.Carbs - tmpl.carbs + overrideCarbs
			} else {
				result.Calories += overrideCal
				result.Protein += overrideProtein
				result.Fat += overrideFat
				result.Carbs += overrideCarbs
			}
		}
	}

	return result, nil
}
