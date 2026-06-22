package resolver

import (
	"monorepo-template/apps/api/internal/appconfig"
	"monorepo-template/apps/api/internal/atlas/service"
)

type Resolver struct {
	SettingsService               service.SettingsService
	PinService                    service.PinService
	ExerciseService               service.ExerciseService
	CardioService                 service.CardioService
	BodyWeightService             service.BodyWeightService
	BodyCheckInService            service.BodyCheckInService
	BodyMeasurementService        service.BodyMeasurementService
	WeekFlagService               service.WeekFlagService
	NutritionProductService       service.NutritionProductService
	NutritionTemplateService      service.NutritionTemplateService
	NutritionTemplateItemService  service.NutritionTemplateItemService
	DailyNutritionOverrideService service.DailyNutritionOverrideService
	NutritionMacroService         service.NutritionMacroService
	BodyChartService              service.BodyChartService
	NutritionWeeklyAvgService     service.NutritionWeeklyAvgService
	UserProfileService            service.UserProfileService
	AiExportService               service.AiExportService
	AiReviewService               service.AiReviewService
	BackupExportService           service.BackupExportService
	BackupImportService           service.BackupImportService
	AiExportConfig                appconfig.AiExportConfig
}