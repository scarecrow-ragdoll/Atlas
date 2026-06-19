package resolver

import (
	"monorepo-template/apps/api/internal/atlas/service"
)

type Resolver struct {
	SettingsService service.SettingsService
	PinService      service.PinService
	ExerciseService service.ExerciseService
}