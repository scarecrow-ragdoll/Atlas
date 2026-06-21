// FILE: apps/api/internal/atlas/service/ai_export_data_provider.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Define the AiExportDataProvider interface for WAVE-07 AI export data source dependencies.
//   SCOPE: Data source abstraction for workout, cardio, body weight, check-in, measurement, week flag, nutrition, and photo retrieval.
//   DEPENDS: apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-07.
//   ROLE: TYPES
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AiExportDataProvider - Interface for data source dependencies in AI export generation.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added AI export data provider interface for WAVE-07.
// END_CHANGE_SUMMARY

package service

import (
	"context"
	"monorepo-template/apps/api/internal/atlas/models"
)

// AiExportDataProvider defines the data source dependencies for AI export generation.
type AiExportDataProvider interface {
	GetWorkoutSummary(ctx context.Context, userID string, from, to models.Date) ([]any, error)
	GetCardioEntries(ctx context.Context, userID string, from, to models.Date) ([]any, error)
	GetBodyWeightEntries(ctx context.Context, userID string, from, to models.Date) ([]any, error)
	GetBodyCheckIns(ctx context.Context, userID string, from, to models.Date) ([]any, error)
	GetBodyMeasurements(ctx context.Context, userID string, from, to models.Date) ([]any, error)
	GetWeekFlags(ctx context.Context, userID string, from, to models.Date) ([]any, error)
	GetNutritionMacros(ctx context.Context, userID string, from, to models.Date) ([]any, error)
	GetProgressPhotos(ctx context.Context, userID string, from, to models.Date) ([]ExportPhoto, error)
}