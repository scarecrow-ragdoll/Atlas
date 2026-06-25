// FILE: apps/api/internal/atlas/service/default_data_provider.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Test/stub implementation of AiExportDataProvider that returns empty data for all sources.
//   SCOPE: Provides empty responses for workout, cardio, body weight, check-in, measurement, week flag, detailed nutrition, and photo data; runtime must use AtlasAiExportDataProvider instead.
//   DEPENDS: apps/api/internal/atlas/models, apps/api/internal/atlas/service.
//   LINKS: M-API / V-M-API / WAVE-07 / DDEC-W07-012.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   NewDefaultAiExportDataProvider - Creates a stub AiExportDataProvider returning empty data.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Kept default provider as test/stub-only and updated detailed nutrition export methods.
//   LAST_CHANGE: 1.0.0 - Added default stub data provider for WAVE-07 (DDEC-W07-012).
// END_CHANGE_SUMMARY

package service

import (
	"context"
	"monorepo-template/apps/api/internal/atlas/models"
)

func NewDefaultAiExportDataProvider() AiExportDataProvider {
	return &defaultDataProvider{}
}

type defaultDataProvider struct{}

func (p *defaultDataProvider) GetWorkoutSummary(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return []any{}, nil
}

func (p *defaultDataProvider) GetCardioEntries(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return []any{}, nil
}

func (p *defaultDataProvider) GetBodyWeightEntries(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return []any{}, nil
}

func (p *defaultDataProvider) GetBodyCheckIns(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return []any{}, nil
}

func (p *defaultDataProvider) GetBodyMeasurements(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return []any{}, nil
}

func (p *defaultDataProvider) GetWeekFlags(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return []any{}, nil
}

func (p *defaultDataProvider) GetDailyNutritionExport(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return []any{}, nil
}

func (p *defaultDataProvider) GetNutritionTemplateExport(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return []any{}, nil
}

func (p *defaultDataProvider) GetLegacyNutritionExport(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return []any{}, nil
}

func (p *defaultDataProvider) GetProgressPhotos(ctx context.Context, userID string, from, to models.Date) ([]ExportPhoto, error) {
	return []ExportPhoto{}, nil
}
