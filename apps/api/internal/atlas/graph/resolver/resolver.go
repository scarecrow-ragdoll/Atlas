package resolver

import (
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
}