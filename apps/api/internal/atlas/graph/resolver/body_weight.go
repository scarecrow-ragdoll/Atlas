// FILE: apps/api/internal/atlas/graph/resolver/body_weight.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement GraphQL resolvers for WAVE-04 BodyWeightEntry queries and mutations.
//   SCOPE: bodyWeightEntry(id), bodyWeightEntries(from, to), latestBodyWeight, createBodyWeightEntry, updateBodyWeightEntry, deleteBodyWeightEntry.
//   DEPENDS: apps/api/internal/atlas/service.BodyWeightService, apps/api/internal/atlas/middleware, apps/api/internal/atlas/models, generated body_tracking.resolvers.go.
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

func (r *Resolver) GetBodyWeightEntry(ctx context.Context, id string) (*models.BodyWeightResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.BodyWeightResult{
			AuthErr: &models.BodyAuthErr{
				Message: "unauthorized",
				Code:    models.BodyErrorAuth,
			},
		}, nil
	}

	entry, err := r.BodyWeightService.GetByID(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrBodyWeightNotFound) {
			return &models.BodyWeightResult{
				NotFoundErr: &models.BodyNotFoundErr{
					Message: "body weight entry not found",
					Code:    models.BodyErrorNotFound,
				},
			}, nil
		}
		return nil, nil
	}

	return &models.BodyWeightResult{Entry: entry}, nil
}

func (r *Resolver) GetBodyWeightEntries(ctx context.Context, from models.Date, to models.Date) (*models.BodyWeightEntriesResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.BodyWeightEntriesResult{
			AuthErr: &models.BodyAuthErr{
				Message: "unauthorized",
				Code:    models.BodyErrorAuth,
			},
		}, nil
	}

	entries, err := r.BodyWeightService.ListByDateRange(ctx, userID, from, to)
	if err != nil {
		return nil, nil
	}

	return &models.BodyWeightEntriesResult{Entries: entries}, nil
}

func (r *Resolver) GetLatestBodyWeight(ctx context.Context) (*models.BodyWeightResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.BodyWeightResult{
			AuthErr: &models.BodyAuthErr{
				Message: "unauthorized",
				Code:    models.BodyErrorAuth,
			},
		}, nil
	}

	entry, err := r.BodyWeightService.Latest(ctx, userID)
	if err != nil {
		return nil, nil
	}
	if entry == nil {
		return &models.BodyWeightResult{}, nil
	}

	return &models.BodyWeightResult{Entry: entry}, nil
}

func (r *Resolver) CreateBodyWeightEntry(ctx context.Context, input models.CreateBodyWeightInput) (*models.BodyWeightResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.BodyWeightResult{
			AuthErr: &models.BodyAuthErr{
				Message: "unauthorized",
				Code:    models.BodyErrorAuth,
			},
		}, nil
	}

	entry, err := r.BodyWeightService.Create(ctx, userID, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrBodyWeightInvalid),
			errors.Is(err, atlasService.ErrBodyWeightInvalidSource):
			return &models.BodyWeightResult{
				ValidationErr: &models.BodyValidationErr{
					Message: err.Error(),
					Code:    models.BodyErrorValidation,
				},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.BodyWeightResult{Entry: entry}, nil
}

func (r *Resolver) UpdateBodyWeightEntry(ctx context.Context, id string, input models.UpdateBodyWeightInput) (*models.BodyWeightResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.BodyWeightResult{
			AuthErr: &models.BodyAuthErr{
				Message: "unauthorized",
				Code:    models.BodyErrorAuth,
			},
		}, nil
	}

	entry, err := r.BodyWeightService.Update(ctx, userID, id, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrBodyWeightInvalid),
			errors.Is(err, atlasService.ErrBodyWeightInvalidSource):
			return &models.BodyWeightResult{
				ValidationErr: &models.BodyValidationErr{
					Message: err.Error(),
					Code:    models.BodyErrorValidation,
				},
			}, nil
		case errors.Is(err, atlasService.ErrBodyWeightNotFound):
			return &models.BodyWeightResult{
				NotFoundErr: &models.BodyNotFoundErr{
					Message: "body weight entry not found",
					Code:    models.BodyErrorNotFound,
				},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.BodyWeightResult{Entry: entry}, nil
}

func (r *Resolver) DeleteBodyWeightEntry(ctx context.Context, id string) (*models.BodyWeightResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.BodyWeightResult{
			AuthErr: &models.BodyAuthErr{
				Message: "unauthorized",
				Code:    models.BodyErrorAuth,
			},
		}, nil
	}

	entry, err := r.BodyWeightService.Delete(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrBodyWeightNotFound) {
			return &models.BodyWeightResult{
				NotFoundErr: &models.BodyNotFoundErr{
					Message: "body weight entry not found",
					Code:    models.BodyErrorNotFound,
				},
			}, nil
		}
		return nil, nil
	}

	return &models.BodyWeightResult{Entry: entry}, nil
}
