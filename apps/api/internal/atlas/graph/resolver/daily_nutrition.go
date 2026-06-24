// FILE: apps/api/internal/atlas/graph/resolver/daily_nutrition.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement factual daily nutrition log and weekly template apply GraphQL adapter methods.
//   SCOPE: PIN-auth user extraction, DailyNutritionLogService query/mutations, NutritionTemplateApplyService mutation, and nutrition error-result mapping; excludes repository/service business rules.
//   DEPENDS: atlas middleware user context, atlas nutrition models, atlas nutrition services.
//   LINKS: M-API-NUTRITION / V-M-API-NUTRITION.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   GetDailyNutritionLog - Query adapter for get-or-create factual daily log by date.
//   UpdateDailyNutritionLogNotes - Mutation adapter for notes update by daily log id.
//   AddDailyNutritionEntry/UpdateDailyNutritionEntry/DeleteDailyNutritionEntry - Entry mutation adapters returning refreshed daily aggregates.
//   ApplyNutritionTemplateToWeek - Mutation adapter for seed_empty_days weekly template application.
//   atlasGraphQLDatePtr/atlasGraphQLTimePtr - Shared generated-resolver scalar parsing helpers.
// END_MODULE_MAP

package resolver

import (
	"context"
	"errors"
	"time"

	"monorepo-template/apps/api/internal/atlas/middleware"
	"monorepo-template/apps/api/internal/atlas/models"
	atlasService "monorepo-template/apps/api/internal/atlas/service"
)

func (r *Resolver) GetDailyNutritionLog(ctx context.Context, date models.Date) (*models.DailyNutritionLogResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return dailyNutritionAuthResult(), nil
	}

	log, err := r.DailyNutritionLogService.GetByDate(ctx, userID, date)
	if err != nil {
		return dailyNutritionResultFromError(err), nil
	}
	return &models.DailyNutritionLogResult{DailyNutritionLog: log}, nil
}

func (r *Resolver) UpdateDailyNutritionLogNotes(ctx context.Context, id string, input models.UpdateDailyNutritionLogNotesInput) (*models.DailyNutritionLogResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return dailyNutritionAuthResult(), nil
	}

	log, err := r.DailyNutritionLogService.UpdateNotes(ctx, userID, id, input)
	if err != nil {
		return dailyNutritionResultFromError(err), nil
	}
	return &models.DailyNutritionLogResult{DailyNutritionLog: log}, nil
}

func (r *Resolver) AddDailyNutritionEntry(ctx context.Context, input models.AddDailyNutritionEntryInput) (*models.DailyNutritionLogResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return dailyNutritionAuthResult(), nil
	}

	log, err := r.DailyNutritionLogService.AddEntry(ctx, userID, input)
	if err != nil {
		return dailyNutritionResultFromError(err), nil
	}
	return &models.DailyNutritionLogResult{DailyNutritionLog: log}, nil
}

func (r *Resolver) UpdateDailyNutritionEntry(ctx context.Context, id string, input models.UpdateDailyNutritionEntryInput) (*models.DailyNutritionLogResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return dailyNutritionAuthResult(), nil
	}

	log, err := r.DailyNutritionLogService.UpdateEntry(ctx, userID, id, input)
	if err != nil {
		return dailyNutritionResultFromError(err), nil
	}
	return &models.DailyNutritionLogResult{DailyNutritionLog: log}, nil
}

func (r *Resolver) DeleteDailyNutritionEntry(ctx context.Context, id string) (*models.DailyNutritionLogResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return dailyNutritionAuthResult(), nil
	}

	log, err := r.DailyNutritionLogService.DeleteEntry(ctx, userID, id)
	if err != nil {
		return dailyNutritionResultFromError(err), nil
	}
	return &models.DailyNutritionLogResult{DailyNutritionLog: log}, nil
}

func (r *Resolver) ApplyNutritionTemplateToWeek(ctx context.Context, templateID string, mode models.NutritionTemplateApplyMode) (*models.NutritionTemplateApplyResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.NutritionTemplateApplyResult{
			AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
		}, nil
	}

	result, err := r.NutritionTemplateApplyService.ApplyToWeek(ctx, userID, templateID, mode)
	if err != nil {
		return nutritionTemplateApplyResultFromError(err), nil
	}
	if result == nil {
		return &models.NutritionTemplateApplyResult{}, nil
	}
	return result, nil
}

func dailyNutritionAuthResult() *models.DailyNutritionLogResult {
	return &models.DailyNutritionLogResult{
		AuthErr: &models.NutritionAuthErr{Message: "unauthorized", Code: models.NutritionErrorAuth},
	}
}

func dailyNutritionResultFromError(err error) *models.DailyNutritionLogResult {
	switch {
	case errors.Is(err, atlasService.ErrDailyNutritionAmountInvalid),
		errors.Is(err, atlasService.ErrDailyNutritionProductInactive):
		return &models.DailyNutritionLogResult{
			ValidationErr: &models.NutritionValidationErr{Message: err.Error(), Code: models.NutritionErrorValidation},
		}
	case errors.Is(err, atlasService.ErrDailyNutritionProductNotFound),
		errors.Is(err, atlasService.ErrDailyNutritionLogNotFound),
		errors.Is(err, atlasService.ErrDailyNutritionEntryNotFound),
		errors.Is(err, atlasService.ErrProductNotFound):
		return &models.DailyNutritionLogResult{
			NotFoundErr: &models.NutritionNotFoundErr{Message: err.Error(), Code: models.NutritionErrorNotFound},
		}
	default:
		return &models.DailyNutritionLogResult{
			ValidationErr: &models.NutritionValidationErr{Message: err.Error(), Code: models.NutritionErrorInternal},
		}
	}
}

func nutritionTemplateApplyResultFromError(err error) *models.NutritionTemplateApplyResult {
	switch {
	case errors.Is(err, atlasService.ErrNutritionTemplateApplyModeUnsupported):
		return &models.NutritionTemplateApplyResult{
			ValidationErr: &models.NutritionValidationErr{Message: err.Error(), Code: models.NutritionErrorValidation},
		}
	case errors.Is(err, atlasService.ErrTemplateNotFound):
		return &models.NutritionTemplateApplyResult{
			NotFoundErr: &models.NutritionNotFoundErr{Message: "template not found", Code: models.NutritionErrorNotFound},
		}
	default:
		return &models.NutritionTemplateApplyResult{
			ValidationErr: &models.NutritionValidationErr{Message: err.Error(), Code: models.NutritionErrorInternal},
		}
	}
}

func atlasGraphQLDatePtr(value string) (*models.Date, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := models.ParseDate(value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func atlasGraphQLTimePtr(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339} {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return &parsed, nil
		}
	}
	parsed, err := time.Parse("2006-01-02 15:04:05.999999-07", value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
