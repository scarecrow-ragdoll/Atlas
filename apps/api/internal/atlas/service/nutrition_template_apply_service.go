// FILE: apps/api/internal/atlas/service/nutrition_template_apply_service.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Apply weekly nutrition templates into factual daily nutrition logs through legacy-safe seed modes.
//   SCOPE: seed_empty_days weekly apply orchestration, template ownership checks, template item loading, active product validation, legacy daily override skip checks, and per-day atomic repository seeding.
//   DEPENDS: postgres nutrition template, template item, product, daily log, daily override repositories; atlas/models; zap.
//   LINKS: M-API-NUTRITION / V-M-API-NUTRITION.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   NutritionTemplateApplyService - Interface for applying one weekly template to one week.
//   NewNutritionTemplateApplyService - Creates the dedicated apply service without expanding template CRUD.
//   ApplyToWeek - Applies supported weekly-template modes and returns per-date outcomes.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Moved legacy resolver implementation to daily_nutrition_legacy_resolver.go for Task 6 resolution metadata.
//   LAST_CHANGE: 1.0.0 - Added Task 4 seed_empty_days template apply service and legacy resolver.
// END_CHANGE_SUMMARY

package service

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/atlas/repository/postgres"
)

const legacyNutritionSeedSkipReason = "legacy nutrition exists; migrate or review before seeding"

var (
	ErrNutritionTemplateApplyModeUnsupported = errors.New("nutrition template apply mode is unsupported")
)

type NutritionTemplateApplyService interface {
	ApplyToWeek(ctx context.Context, userID string, templateID string, mode models.NutritionTemplateApplyMode) (*models.NutritionTemplateApplyResult, error)
}

type nutritionTemplateApplyService struct {
	templateRepo     postgres.NutritionTemplateRepository
	templateItemRepo postgres.NutritionTemplateItemRepository
	productRepo      postgres.NutritionProductRepository
	dailyRepo        postgres.DailyNutritionLogRepository
	legacyResolver   DailyNutritionLegacyResolver
	logger           *zap.Logger
}

func NewNutritionTemplateApplyService(
	templateRepo postgres.NutritionTemplateRepository,
	templateItemRepo postgres.NutritionTemplateItemRepository,
	productRepo postgres.NutritionProductRepository,
	dailyRepo postgres.DailyNutritionLogRepository,
	legacyResolver DailyNutritionLegacyResolver,
	logger *zap.Logger,
) NutritionTemplateApplyService {
	if logger == nil {
		logger = zap.NewNop()
	}
	if legacyResolver == nil {
		legacyResolver = &dailyNutritionLegacyResolver{}
	}
	return &nutritionTemplateApplyService{
		templateRepo:     templateRepo,
		templateItemRepo: templateItemRepo,
		productRepo:      productRepo,
		dailyRepo:        dailyRepo,
		legacyResolver:   legacyResolver,
		logger:           logger,
	}
}

// START_CONTRACT: ApplyToWeek
//
//	PURPOSE: Seed one weekly template into empty factual daily logs while preserving existing factual or legacy days.
//	INPUTS: { userID: string - owner scope, templateID: string - weekly template id, mode: NutritionTemplateApplyMode - currently seed_empty_days only }
//	OUTPUTS: { *NutritionTemplateApplyResult - seven per-date outcomes for successful orchestration }
//	SIDE_EFFECTS: Creates factual daily logs and entries only through the atomic daily repository seed helper.
//	LINKS: M-API-NUTRITION / V-M-API-NUTRITION.
//
// END_CONTRACT: ApplyToWeek
func (s *nutritionTemplateApplyService) ApplyToWeek(ctx context.Context, userID string, templateID string, mode models.NutritionTemplateApplyMode) (*models.NutritionTemplateApplyResult, error) {
	s.logger.Info("[NutritionTemplateApply][apply]")
	if mode != models.ApplyModeSeedEmptyDays {
		return nil, ErrNutritionTemplateApplyModeUnsupported
	}

	template, err := s.templateRepo.GetByID(ctx, userID, templateID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_apply_service.ApplyToWeek: %w", err)
	}
	if template == nil {
		return nil, ErrTemplateNotFound
	}

	items, err := s.templateItemRepo.ListByTemplate(ctx, template.ID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_apply_service.ApplyToWeek: %w", err)
	}

	seedItems, productConflictReason, err := s.buildSeedItems(ctx, userID, items)
	if err != nil {
		return nil, err
	}

	weekStart := template.WeekStartDate.Time()
	result := &models.NutritionTemplateApplyResult{
		WeekStartDate: template.WeekStartDate.String(),
		WeekEndDate:   weekStart.AddDate(0, 0, 6).Format("2006-01-02"),
		Mode:          mode,
		Dates:         make([]models.NutritionTemplateApplyDateResult, 0, 7),
	}

	for offset := 0; offset < 7; offset++ {
		date := models.MustDate(weekStart.AddDate(0, 0, offset).Format("2006-01-02"))
		hasLegacy, err := s.legacyResolver.HasLegacyNutrition(ctx, userID, date)
		if err != nil {
			return nil, fmt.Errorf("nutrition_template_apply_service.ApplyToWeek: %w", err)
		}
		if hasLegacy {
			result.Dates = append(result.Dates, models.NutritionTemplateApplyDateResult{
				Date:   date.String(),
				Status: models.ApplyDateSkipped,
				Reason: stringPtr(legacyNutritionSeedSkipReason),
			})
			continue
		}
		if productConflictReason != nil {
			result.Dates = append(result.Dates, models.NutritionTemplateApplyDateResult{
				Date:   date.String(),
				Status: models.ApplyDateConflict,
				Reason: productConflictReason,
			})
			continue
		}

		seeded, err := s.dailyRepo.SeedEntriesIfEmpty(ctx, userID, date, seedItems)
		if err != nil {
			return nil, fmt.Errorf("nutrition_template_apply_service.ApplyToWeek: %w", err)
		}
		if seeded.Created {
			result.Dates = append(result.Dates, models.NutritionTemplateApplyDateResult{
				Date:       date.String(),
				Status:     models.ApplyDateCreated,
				EntryCount: seeded.EntryCount,
			})
			continue
		}
		result.Dates = append(result.Dates, models.NutritionTemplateApplyDateResult{
			Date:       date.String(),
			Status:     models.ApplyDateSkipped,
			EntryCount: seeded.EntryCount,
			Reason:     stringPtr("day has entries"),
		})
	}

	return result, nil
}

func (s *nutritionTemplateApplyService) buildSeedItems(ctx context.Context, userID string, items []models.NutritionTemplateItemRecord) ([]models.DailyNutritionSeedEntryInput, *string, error) {
	seedItems := make([]models.DailyNutritionSeedEntryInput, 0, len(items))
	for i, item := range items {
		product, err := s.productRepo.GetByID(ctx, userID, item.ProductID)
		if err != nil {
			return nil, nil, fmt.Errorf("nutrition_template_apply_service.buildSeedItems: %w", err)
		}
		if product == nil || !product.IsActive {
			return nil, stringPtr("template product missing or inactive"), nil
		}
		seedItems = append(seedItems, models.DailyNutritionSeedEntryInput{
			ProductID:   item.ProductID,
			AmountGrams: item.AmountGrams,
			MealLabel:   item.MealLabel,
			Notes:       item.Notes,
			Position:    int32(i),
		})
	}
	return seedItems, nil, nil
}

func stringPtr(value string) *string {
	return &value
}
