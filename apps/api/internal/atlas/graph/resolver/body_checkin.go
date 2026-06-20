// FILE: apps/api/internal/atlas/graph/resolver/body_checkin.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement GraphQL resolvers for WAVE-04 BodyCheckIn and BodyMeasurement queries and mutations.
//   SCOPE: bodyCheckIn(id), bodyCheckIns(from, to), bodyMeasurements(checkInID), createBodyCheckIn, updateBodyCheckIn, deleteBodyCheckIn, createBodyMeasurement, updateBodyMeasurement, deleteBodyMeasurement.
//   DEPENDS: apps/api/internal/atlas/service.BodyCheckInService, apps/api/internal/atlas/service.BodyMeasurementService, apps/api/internal/atlas/middleware, apps/api/internal/atlas/models, generated body_tracking.resolvers.go.
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

func (r *Resolver) GetBodyCheckIn(ctx context.Context, id string) (*models.BodyCheckInResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.BodyCheckInResult{
			AuthErr: &models.BodyAuthErr{
				Message: "unauthorized",
				Code:    models.BodyErrorAuth,
			},
		}, nil
	}

	checkIn, err := r.BodyCheckInService.GetByID(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrCheckInNotFound) {
			return &models.BodyCheckInResult{
				NotFoundErr: &models.BodyNotFoundErr{
					Message: "body check-in not found",
					Code:    models.BodyErrorNotFound,
				},
			}, nil
		}
		return nil, nil
	}

	return &models.BodyCheckInResult{CheckIn: checkIn}, nil
}

func (r *Resolver) GetBodyCheckIns(ctx context.Context, from models.Date, to models.Date) (*models.BodyCheckInsResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.BodyCheckInsResult{
			AuthErr: &models.BodyAuthErr{
				Message: "unauthorized",
				Code:    models.BodyErrorAuth,
			},
		}, nil
	}

	checkIns, err := r.BodyCheckInService.ListByDateRange(ctx, userID, from, to)
	if err != nil {
		return nil, nil
	}

	return &models.BodyCheckInsResult{CheckIns: checkIns}, nil
}

func (r *Resolver) GetBodyMeasurements(ctx context.Context, checkInID string) ([]*models.BodyMeasurement, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return nil, nil
	}

	measurements, err := r.BodyMeasurementService.ListByCheckIn(ctx, userID, checkInID)
	if err != nil {
		return nil, nil
	}

	items := make([]*models.BodyMeasurement, len(measurements))
	for i := range measurements {
		items[i] = &measurements[i]
	}
	return items, nil
}

func (r *Resolver) CreateBodyCheckIn(ctx context.Context, input models.CreateCheckInInput) (*models.BodyCheckInResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.BodyCheckInResult{
			AuthErr: &models.BodyAuthErr{
				Message: "unauthorized",
				Code:    models.BodyErrorAuth,
			},
		}, nil
	}

	checkIn, err := r.BodyCheckInService.Create(ctx, userID, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrCheckInWeightInvalid),
			errors.Is(err, atlasService.ErrCheckInBodyFatInvalid):
			return &models.BodyCheckInResult{
				ValidationErr: &models.BodyValidationErr{
					Message: err.Error(),
					Code:    models.BodyErrorValidation,
				},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.BodyCheckInResult{CheckIn: checkIn}, nil
}

func (r *Resolver) UpdateBodyCheckIn(ctx context.Context, id string, input models.UpdateCheckInInput) (*models.BodyCheckInResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.BodyCheckInResult{
			AuthErr: &models.BodyAuthErr{
				Message: "unauthorized",
				Code:    models.BodyErrorAuth,
			},
		}, nil
	}

	checkIn, err := r.BodyCheckInService.Update(ctx, userID, id, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrCheckInWeightInvalid),
			errors.Is(err, atlasService.ErrCheckInBodyFatInvalid):
			return &models.BodyCheckInResult{
				ValidationErr: &models.BodyValidationErr{
					Message: err.Error(),
					Code:    models.BodyErrorValidation,
				},
			}, nil
		case errors.Is(err, atlasService.ErrCheckInNotFound):
			return &models.BodyCheckInResult{
				NotFoundErr: &models.BodyNotFoundErr{
					Message: "body check-in not found",
					Code:    models.BodyErrorNotFound,
				},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.BodyCheckInResult{CheckIn: checkIn}, nil
}

func (r *Resolver) DeleteBodyCheckIn(ctx context.Context, id string) (*models.BodyCheckInResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.BodyCheckInResult{
			AuthErr: &models.BodyAuthErr{
				Message: "unauthorized",
				Code:    models.BodyErrorAuth,
			},
		}, nil
	}

	checkIn, err := r.BodyCheckInService.Delete(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrCheckInNotFound) {
			return &models.BodyCheckInResult{
				NotFoundErr: &models.BodyNotFoundErr{
					Message: "body check-in not found",
					Code:    models.BodyErrorNotFound,
				},
			}, nil
		}
		return nil, nil
	}

	return &models.BodyCheckInResult{CheckIn: checkIn}, nil
}

func (r *Resolver) CreateBodyMeasurement(ctx context.Context, checkInID string, input models.CreateMeasurementInput) (*models.BodyMeasurementResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.BodyMeasurementResult{
			AuthErr: &models.BodyAuthErr{
				Message: "unauthorized",
				Code:    models.BodyErrorAuth,
			},
		}, nil
	}

	measurement, err := r.BodyMeasurementService.Create(ctx, userID, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrMeasurementTypeInvalid),
			errors.Is(err, atlasService.ErrMeasurementValueInvalid),
			errors.Is(err, atlasService.ErrMeasurementSideInvalid):
			return &models.BodyMeasurementResult{
				ValidationErr: &models.BodyValidationErr{
					Message: err.Error(),
					Code:    models.BodyErrorValidation,
				},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.BodyMeasurementResult{Measurement: measurement}, nil
}

func (r *Resolver) UpdateBodyMeasurement(ctx context.Context, id string, input models.UpdateMeasurementInput) (*models.BodyMeasurementResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.BodyMeasurementResult{
			AuthErr: &models.BodyAuthErr{
				Message: "unauthorized",
				Code:    models.BodyErrorAuth,
			},
		}, nil
	}

	measurement, err := r.BodyMeasurementService.Update(ctx, userID, id, input)
	if err != nil {
		switch {
		case errors.Is(err, atlasService.ErrMeasurementTypeInvalid),
			errors.Is(err, atlasService.ErrMeasurementValueInvalid),
			errors.Is(err, atlasService.ErrMeasurementSideInvalid):
			return &models.BodyMeasurementResult{
				ValidationErr: &models.BodyValidationErr{
					Message: err.Error(),
					Code:    models.BodyErrorValidation,
				},
			}, nil
		case errors.Is(err, atlasService.ErrMeasurementNotFound):
			return &models.BodyMeasurementResult{
				NotFoundErr: &models.BodyNotFoundErr{
					Message: "body measurement not found",
					Code:    models.BodyErrorNotFound,
				},
			}, nil
		default:
			return nil, nil
		}
	}

	return &models.BodyMeasurementResult{Measurement: measurement}, nil
}

func (r *Resolver) DeleteBodyMeasurement(ctx context.Context, id string) (*models.BodyMeasurementResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.BodyMeasurementResult{
			AuthErr: &models.BodyAuthErr{
				Message: "unauthorized",
				Code:    models.BodyErrorAuth,
			},
		}, nil
	}

	measurement, err := r.BodyMeasurementService.Delete(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrMeasurementNotFound) {
			return &models.BodyMeasurementResult{
				NotFoundErr: &models.BodyNotFoundErr{
					Message: "body measurement not found",
					Code:    models.BodyErrorNotFound,
				},
			}, nil
		}
		return nil, nil
	}

	return &models.BodyMeasurementResult{Measurement: measurement}, nil
}
