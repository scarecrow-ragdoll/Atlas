// FILE: apps/api/internal/atlas/graph/resolver/cardio.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement GraphQL resolvers for WAVE-04 CardioEntry queries and mutations.
//   SCOPE: cardioEntry(id), cardioEntries(date), createCardioEntry, updateCardioEntry, deleteCardioEntry.
//   DEPENDS: apps/api/internal/atlas/service.CardioService, apps/api/internal/atlas/middleware, apps/api/internal/atlas/models, generated cardio.resolvers.go.
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

func (r *Resolver) GetCardioEntry(ctx context.Context, id string) (*models.CardioEntryResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.CardioEntryResult{
			AuthErr: &models.CardioAuthErr{
				Message: "unauthorized",
				Code:    models.CardioErrorAuth,
			},
		}, nil
	}

	entry, err := r.CardioService.GetByID(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrCardioNotFound) {
			return &models.CardioEntryResult{
				NotFoundErr: &models.CardioNotFoundErr{
					Message: "cardio entry not found",
					Code:    models.CardioErrorNotFound,
				},
			}, nil
		}
		return nil, nil
	}

	return &models.CardioEntryResult{CardioEntry: entry}, nil
}

func (r *Resolver) GetCardioEntries(ctx context.Context, date models.Date) (*models.CardioEntriesResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.CardioEntriesResult{
			AuthErr: &models.CardioAuthErr{
				Message: "unauthorized",
				Code:    models.CardioErrorAuth,
			},
		}, nil
	}

	entries, err := r.CardioService.ListByDate(ctx, userID, date)
	if err != nil {
		return nil, nil
	}

	return &models.CardioEntriesResult{Entries: entries}, nil
}

func (r *Resolver) CreateCardioEntry(ctx context.Context, input models.CreateCardioInput) (*models.CardioEntryResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.CardioEntryResult{
			AuthErr: &models.CardioAuthErr{
				Message: "unauthorized",
				Code:    models.CardioErrorAuth,
			},
		}, nil
	}

	entry, err := r.CardioService.Create(ctx, userID, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrCardioInvalidType),
			errors.Is(err, atlasService.ErrCardioDurationInvalid),
			errors.Is(err, atlasService.ErrCardioPulseInvalid),
			errors.Is(err, atlasService.ErrCardioZoneInvalid):
			return &models.CardioEntryResult{
				ValidationErr: &models.CardioValidationErr{
					Message: err.Error(),
					Code:    models.CardioErrorValidation,
				},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.CardioEntryResult{CardioEntry: entry}, nil
}

func (r *Resolver) UpdateCardioEntry(ctx context.Context, id string, input models.UpdateCardioInput) (*models.CardioEntryResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.CardioEntryResult{
			AuthErr: &models.CardioAuthErr{
				Message: "unauthorized",
				Code:    models.CardioErrorAuth,
			},
		}, nil
	}

	entry, err := r.CardioService.Update(ctx, userID, id, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrCardioInvalidType),
			errors.Is(err, atlasService.ErrCardioDurationInvalid),
			errors.Is(err, atlasService.ErrCardioPulseInvalid),
			errors.Is(err, atlasService.ErrCardioZoneInvalid):
			return &models.CardioEntryResult{
				ValidationErr: &models.CardioValidationErr{
					Message: err.Error(),
					Code:    models.CardioErrorValidation,
				},
			}, nil
		case errors.Is(err, atlasService.ErrCardioNotFound):
			return &models.CardioEntryResult{
				NotFoundErr: &models.CardioNotFoundErr{
					Message: "cardio entry not found",
					Code:    models.CardioErrorNotFound,
				},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.CardioEntryResult{CardioEntry: entry}, nil
}

func (r *Resolver) DeleteCardioEntry(ctx context.Context, id string) (*models.CardioEntryResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.CardioEntryResult{
			AuthErr: &models.CardioAuthErr{
				Message: "unauthorized",
				Code:    models.CardioErrorAuth,
			},
		}, nil
	}

	entry, err := r.CardioService.Delete(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrCardioNotFound) {
			return &models.CardioEntryResult{
				NotFoundErr: &models.CardioNotFoundErr{
					Message: "cardio entry not found",
					Code:    models.CardioErrorNotFound,
				},
			}, nil
		}
		return nil, nil
	}

	return &models.CardioEntryResult{CardioEntry: entry}, nil
}
