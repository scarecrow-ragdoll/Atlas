// FILE: apps/api/internal/atlas/graph/resolver/week_flag.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement GraphQL resolvers for WAVE-04 WeekFlag queries and mutations.
//   SCOPE: weekFlags(weekStartDate), createWeekFlag, deleteWeekFlag.
//   DEPENDS: apps/api/internal/atlas/service.WeekFlagService, apps/api/internal/atlas/middleware, apps/api/internal/atlas/models, generated week_flag.resolvers.go.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT

package resolver

import (
	"context"
	"errors"

	"monorepo-template/apps/api/internal/atlas/middleware"
	"monorepo-template/apps/api/internal/atlas/models"
	atlasService "monorepo-template/apps/api/internal/atlas/service"
)

func (r *Resolver) GetWeekFlags(ctx context.Context, weekStartDate models.Date) (*models.WeekFlagsResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.WeekFlagsResult{
			AuthErr: &models.WeekFlagAuthErr{
				Message: "unauthorized",
				Code:    models.WeekFlagErrorAuth,
			},
		}, nil
	}

	flags, err := r.WeekFlagService.ListByWeekStart(ctx, userID, weekStartDate)
	if err != nil {
		return nil, nil
	}

	return &models.WeekFlagsResult{Flags: flags}, nil
}

func (r *Resolver) CreateWeekFlag(ctx context.Context, input models.CreateWeekFlagInput) (*models.WeekFlagResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.WeekFlagResult{
			AuthErr: &models.WeekFlagAuthErr{
				Message: "unauthorized",
				Code:    models.WeekFlagErrorAuth,
			},
		}, nil
	}

	flag, err := r.WeekFlagService.Create(ctx, userID, input)
	if err != nil {
		if errors.Is(err, atlasService.ErrWeekFlagInvalidType) {
			return &models.WeekFlagResult{
				ValidationErr: &models.WeekFlagValidationErr{
					Message: err.Error(),
					Code:    models.WeekFlagErrorValidation,
				},
			}, nil
		}
		return nil, nil
	}

	return &models.WeekFlagResult{WeekFlag: flag}, nil
}

func (r *Resolver) DeleteWeekFlag(ctx context.Context, id string) (*models.WeekFlagResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.WeekFlagResult{
			AuthErr: &models.WeekFlagAuthErr{
				Message: "unauthorized",
				Code:    models.WeekFlagErrorAuth,
			},
		}, nil
	}

	flag, err := r.WeekFlagService.Delete(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrWeekFlagNotFound) {
			return &models.WeekFlagResult{
				NotFoundErr: &models.WeekFlagNotFoundErr{
					Message: "week flag not found",
					Code:    models.WeekFlagErrorNotFound,
				},
			}, nil
		}
		return nil, nil
	}

	return &models.WeekFlagResult{WeekFlag: flag}, nil
}
