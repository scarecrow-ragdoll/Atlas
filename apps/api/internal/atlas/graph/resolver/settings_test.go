package resolver_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/atlas/middleware"
	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/atlas/service"
	atlasSvc "monorepo-template/apps/api/internal/atlas/service"

	"monorepo-template/apps/api/internal/atlas/graph/resolver"
)

type mockSettingsService struct {
	getFn    func(ctx context.Context, userID string) (*models.Settings, error)
	updateFn func(ctx context.Context, userID string, input models.SettingsInput) (*models.Settings, error)
}

func (m *mockSettingsService) Get(ctx context.Context, userID string) (*models.Settings, error) {
	return m.getFn(ctx, userID)
}

func (m *mockSettingsService) Update(ctx context.Context, userID string, input models.SettingsInput) (*models.Settings, error) {
	return m.updateFn(ctx, userID, input)
}

type mockPinService struct {
	enableFn     func(ctx context.Context, userID string, pin string) error
	disableFn    func(ctx context.Context, userID string, currentPin string) error
	changeFn     func(ctx context.Context, userID string, currentPin, newPin string) error
	verifyFn     func(ctx context.Context, userID string, pin string) (bool, error)
	isEnabledFn  func(ctx context.Context, userID string) (bool, error)
}

func (m *mockPinService) Enable(ctx context.Context, userID string, pin string) error {
	return m.enableFn(ctx, userID, pin)
}

func (m *mockPinService) Disable(ctx context.Context, userID string, currentPin string) error {
	return m.disableFn(ctx, userID, currentPin)
}

func (m *mockPinService) Change(ctx context.Context, userID string, currentPin, newPin string) error {
	return m.changeFn(ctx, userID, currentPin, newPin)
}

func (m *mockPinService) Verify(ctx context.Context, userID string, pin string) (bool, error) {
	return m.verifyFn(ctx, userID, pin)
}

func (m *mockPinService) IsEnabled(ctx context.Context, userID string) (bool, error) {
	return m.isEnabledFn(ctx, userID)
}

func userCtx(userID string) context.Context {
	return middleware.ContextWithAtlasUserID(context.Background(), userID)
}

func TestSettingsResolver_Settings_HappyPath(t *testing.T) {
	r := &resolver.Resolver{
		SettingsService: &mockSettingsService{
			getFn: func(ctx context.Context, userID string) (*models.Settings, error) {
				assert.Equal(t, "test-uid", userID)
				return &models.Settings{PinEnabled: true, Units: "metric", DefaultAiExportWeeks: 4}, nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.Settings(ctx)
	require.NoError(t, err)
	require.NotNil(t, result.Settings)
	assert.True(t, result.Settings.PinEnabled)
	assert.Equal(t, "metric", result.Settings.Units)
	assert.Nil(t, result.Error)
}

func TestSettingsResolver_Settings_Unauthorized(t *testing.T) {
	r := &resolver.Resolver{
		SettingsService: &mockSettingsService{},
	}

	result, err := r.Settings(context.Background())
	require.NoError(t, err)
	assert.Nil(t, result.Settings)
	assert.NotNil(t, result.Error)
	assert.Equal(t, models.SettingsErrorUnauthorized, result.Error.Code)
}

func TestSettingsResolver_UpdateSettings_HappyPath(t *testing.T) {
	units := "imperial"
	weeks := 8
	r := &resolver.Resolver{
		SettingsService: &mockSettingsService{
			updateFn: func(ctx context.Context, userID string, input models.SettingsInput) (*models.Settings, error) {
				assert.Equal(t, "test-uid", userID)
				return &models.Settings{PinEnabled: false, Units: *input.Units, DefaultAiExportWeeks: *input.DefaultAiExportWeeks}, nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.UpdateSettings(ctx, models.SettingsInput{Units: &units, DefaultAiExportWeeks: &weeks})
	require.NoError(t, err)
	require.NotNil(t, result.Settings)
	assert.Equal(t, "imperial", result.Settings.Units)
	assert.Equal(t, 8, result.Settings.DefaultAiExportWeeks)
}

func TestSettingsResolver_UpdateSettings_ValidationError(t *testing.T) {
	units := "fahrenheit"
	r := &resolver.Resolver{
		SettingsService: &mockSettingsService{
			updateFn: func(ctx context.Context, userID string, input models.SettingsInput) (*models.Settings, error) {
				return nil, atlasSvc.ErrInvalidUnits
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.UpdateSettings(ctx, models.SettingsInput{Units: &units})
	require.NoError(t, err)
	assert.Nil(t, result.Settings)
	assert.NotNil(t, result.Error)
	assert.Equal(t, models.SettingsErrorValidation, result.Error.Code)
}

func TestSettingsResolver_UpdateSettings_InternalError(t *testing.T) {
	r := &resolver.Resolver{
		SettingsService: &mockSettingsService{
			updateFn: func(ctx context.Context, userID string, input models.SettingsInput) (*models.Settings, error) {
				return nil, errors.New("db failure")
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.UpdateSettings(ctx, models.SettingsInput{})
	require.NoError(t, err)
	assert.Nil(t, result.Settings)
	assert.NotNil(t, result.Error)
	assert.Equal(t, models.SettingsErrorValidation, result.Error.Code)
}

func TestSettingsResolver_EnablePin_HappyPath(t *testing.T) {
	r := &resolver.Resolver{
		PinService: &mockPinService{
			enableFn: func(ctx context.Context, userID string, pin string) error {
				assert.Equal(t, "test-uid", userID)
				assert.Equal(t, "1234", pin)
				return nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.EnablePin(ctx, models.PinEnableInput{Pin: "1234"})
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Nil(t, result.Error)
}

func TestSettingsResolver_EnablePin_AlreadyEnabled(t *testing.T) {
	r := &resolver.Resolver{
		PinService: &mockPinService{
			enableFn: func(ctx context.Context, userID string, pin string) error {
				return service.ErrPinAlreadyEnabled
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.EnablePin(ctx, models.PinEnableInput{Pin: "1234"})
	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, models.PinErrorAlreadyEnabled, result.Error.Code)
}

func TestSettingsResolver_DisablePin_HappyPath(t *testing.T) {
	r := &resolver.Resolver{
		PinService: &mockPinService{
			disableFn: func(ctx context.Context, userID string, currentPin string) error {
				return nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.DisablePin(ctx, models.PinDisableInput{CurrentPin: "1234"})
	require.NoError(t, err)
	assert.True(t, result.Success)
}

func TestSettingsResolver_DisablePin_WrongPin(t *testing.T) {
	r := &resolver.Resolver{
		PinService: &mockPinService{
			disableFn: func(ctx context.Context, userID string, currentPin string) error {
				return service.ErrPinWrongPin
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.DisablePin(ctx, models.PinDisableInput{CurrentPin: "9999"})
	require.NoError(t, err)
	assert.Equal(t, models.PinErrorWrongPin, result.Error.Code)
}

func TestSettingsResolver_DisablePin_AlreadyDisabled(t *testing.T) {
	r := &resolver.Resolver{
		PinService: &mockPinService{
			disableFn: func(ctx context.Context, userID string, currentPin string) error {
				return service.ErrPinAlreadyDisabled
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.DisablePin(ctx, models.PinDisableInput{CurrentPin: "1234"})
	require.NoError(t, err)
	assert.Equal(t, models.PinErrorAlreadyDisabled, result.Error.Code)
}

func TestSettingsResolver_ChangePin_HappyPath(t *testing.T) {
	r := &resolver.Resolver{
		PinService: &mockPinService{
			changeFn: func(ctx context.Context, userID string, currentPin, newPin string) error {
				return nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.ChangePin(ctx, models.PinChangeInput{CurrentPin: "1234", NewPin: "5678"})
	require.NoError(t, err)
	assert.True(t, result.Success)
}

func TestSettingsResolver_ChangePin_WrongPin(t *testing.T) {
	r := &resolver.Resolver{
		PinService: &mockPinService{
			changeFn: func(ctx context.Context, userID string, currentPin, newPin string) error {
				return service.ErrPinWrongPin
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.ChangePin(ctx, models.PinChangeInput{CurrentPin: "9999", NewPin: "5678"})
	require.NoError(t, err)
	assert.Equal(t, models.PinErrorWrongPin, result.Error.Code)
}

func TestSettingsResolver_Unauthorized_WhenNoUser(t *testing.T) {
	r := &resolver.Resolver{PinService: &mockPinService{}}
	ctx := context.Background()

	result, err := r.EnablePin(ctx, models.PinEnableInput{Pin: "1234"})
	require.NoError(t, err)
	assert.Equal(t, models.PinErrorSessionExpired, result.Error.Code)
}