package resolver

import (
	"context"
	"errors"

	"monorepo-template/apps/api/internal/atlas/middleware"
	"monorepo-template/apps/api/internal/atlas/models"
	atlasService "monorepo-template/apps/api/internal/atlas/service"
)

func (r *Resolver) Settings(ctx context.Context) (*models.SettingsResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.SettingsResult{
			Error: &models.SettingsError{
				Message: "unauthorized",
				Code:    models.SettingsErrorUnauthorized,
			},
		}, nil
	}

	settings, err := r.SettingsService.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, atlasService.ErrSettingsNotFound) {
			return &models.SettingsResult{
				Settings: &models.Settings{
					PinEnabled:           false,
					Units:                "metric",
					DefaultAiExportWeeks: 4,
				},
			}, nil
		}
		return &models.SettingsResult{
			Error: &models.SettingsError{
				Message: "internal error",
				Code:    models.SettingsErrorInternal,
			},
		}, nil
	}

	return &models.SettingsResult{Settings: settings}, nil
}

func (r *Resolver) UpdateSettings(ctx context.Context, input models.SettingsInput) (*models.SettingsResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.SettingsResult{
			Error: &models.SettingsError{
				Message: "unauthorized",
				Code:    models.SettingsErrorUnauthorized,
			},
		}, nil
	}

	settings, err := r.SettingsService.Update(ctx, userID, input)
	if err != nil {
		return &models.SettingsResult{
			Error: &models.SettingsError{
				Message: err.Error(),
				Code:    models.SettingsErrorValidation,
			},
		}, nil
	}

	return &models.SettingsResult{Settings: settings}, nil
}

func (r *Resolver) EnablePin(ctx context.Context, input models.PinEnableInput) (*models.PinOperationResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.PinOperationResult{
			Error: &models.PinError{
				Message: "unauthorized",
				Code:    models.PinErrorSessionExpired,
			},
		}, nil
	}

	if err := r.PinService.Enable(ctx, userID, input.Pin); err != nil {
		switch {
		case errors.Is(err, atlasService.ErrPinTooShort):
			return &models.PinOperationResult{
				Error: &models.PinError{Message: err.Error(), Code: models.PinErrorTooShort},
			}, nil
		case errors.Is(err, atlasService.ErrPinTooLong):
			return &models.PinOperationResult{
				Error: &models.PinError{Message: err.Error(), Code: models.PinErrorTooLong},
			}, nil
		case errors.Is(err, atlasService.ErrPinAlreadyEnabled):
			return &models.PinOperationResult{
				Error: &models.PinError{Message: err.Error(), Code: models.PinErrorAlreadyEnabled},
			}, nil
		default:
			return &models.PinOperationResult{
				Error: &models.PinError{Message: "internal error", Code: models.PinErrorInternal},
			}, nil
		}
	}

	return &models.PinOperationResult{Success: true}, nil
}

func (r *Resolver) DisablePin(ctx context.Context, input models.PinDisableInput) (*models.PinOperationResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.PinOperationResult{
			Error: &models.PinError{
				Message: "unauthorized",
				Code:    models.PinErrorSessionExpired,
			},
		}, nil
	}

	if err := r.PinService.Disable(ctx, userID, input.CurrentPin); err != nil {
		switch {
		case errors.Is(err, atlasService.ErrPinWrongPin):
			return &models.PinOperationResult{
				Error: &models.PinError{Message: err.Error(), Code: models.PinErrorWrongPin},
			}, nil
		case errors.Is(err, atlasService.ErrPinAlreadyDisabled):
			return &models.PinOperationResult{
				Error: &models.PinError{Message: err.Error(), Code: models.PinErrorAlreadyDisabled},
			}, nil
		default:
			return &models.PinOperationResult{
				Error: &models.PinError{Message: "internal error", Code: models.PinErrorInternal},
			}, nil
		}
	}

	return &models.PinOperationResult{Success: true}, nil
}

func (r *Resolver) ChangePin(ctx context.Context, input models.PinChangeInput) (*models.PinOperationResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.PinOperationResult{
			Error: &models.PinError{
				Message: "unauthorized",
				Code:    models.PinErrorSessionExpired,
			},
		}, nil
	}

	if err := r.PinService.Change(ctx, userID, input.CurrentPin, input.NewPin); err != nil {
		switch {
		case errors.Is(err, atlasService.ErrPinWrongPin):
			return &models.PinOperationResult{
				Error: &models.PinError{Message: err.Error(), Code: models.PinErrorWrongPin},
			}, nil
		case errors.Is(err, atlasService.ErrPinTooShort):
			return &models.PinOperationResult{
				Error: &models.PinError{Message: err.Error(), Code: models.PinErrorTooShort},
			}, nil
		case errors.Is(err, atlasService.ErrPinTooLong):
			return &models.PinOperationResult{
				Error: &models.PinError{Message: err.Error(), Code: models.PinErrorTooLong},
			}, nil
		default:
			return &models.PinOperationResult{
				Error: &models.PinError{Message: "internal error", Code: models.PinErrorInternal},
			}, nil
		}
	}

	return &models.PinOperationResult{Success: true}, nil
}