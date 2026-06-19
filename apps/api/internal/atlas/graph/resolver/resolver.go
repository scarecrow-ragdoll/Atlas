// FILE: apps/api/internal/atlas/graph/resolver/resolver.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Define the Atlas GraphQL root resolver dependency container.
//   SCOPE: Service dependency injection for Atlas resolver methods; excludes resolver behavior implementation.
//   DEPENDS: apps/api/internal/atlas/service.
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Resolver.SettingsService - Settings resolver service dependency.
//   Resolver.PinService - PIN resolver service dependency.
//   Resolver.ExerciseService - Exercise library resolver service dependency.
//   Resolver.WorkoutService - DailyLog strength workout diary resolver service dependency.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added WAVE-03 WorkoutService injection to the Atlas root resolver.
// END_CHANGE_SUMMARY

package resolver

import (
	"monorepo-template/apps/api/internal/atlas/service"
)

type Resolver struct {
	SettingsService service.SettingsService
	PinService      service.PinService
	ExerciseService service.ExerciseService
	WorkoutService  service.WorkoutService
}
