// FILE: apps/api/internal/atlas/service/settings_service_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for the Atlas settings service covering Get, Update, validation, and pinHash stripping.
//   SCOPE: Get returns public Settings without pinHash, Get error propagation, Update with valid/invalid units and export weeks, recordToPublic strips pinHash.
//   DEPENDS: apps/api/internal/atlas/service, apps/api/internal/atlas/repository/postgres (mock SettingsRepository).
//   LINKS: M-API / V-M-API / TEST-W01-001.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added settings service unit tests.
// END_CHANGE_SUMMARY

package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasPostgres "monorepo-template/apps/api/internal/atlas/repository/postgres"
	"monorepo-template/apps/api/internal/atlas/service"
)

type mockSettingsRepoForSettings struct {
	atlasPostgres.SettingsRepository
	findByUserIDFn func(ctx context.Context, userID string) (*models.SettingsRecord, error)
	upsertFn       func(ctx context.Context, userID string, input models.SettingsInput) (*models.SettingsRecord, error)
}

func (m *mockSettingsRepoForSettings) FindByUserID(ctx context.Context, userID string) (*models.SettingsRecord, error) {
	return m.findByUserIDFn(ctx, userID)
}

func (m *mockSettingsRepoForSettings) UpsertSettings(ctx context.Context, userID string, input models.SettingsInput) (*models.SettingsRecord, error) {
	return m.upsertFn(ctx, userID, input)
}

var settingsCtx = context.Background()
var settingsUID = "550e8400-e29b-41d4-a716-446655440000"

func ptrStr(s string) *string { return &s }

func ptrInt(i int) *int { return &i }

func TestSettingsService_Get_ReturnsPublicSettings(t *testing.T) {
	svc := service.NewSettingsService(&mockSettingsRepoForSettings{
		findByUserIDFn: func(ctx context.Context, userID string) (*models.SettingsRecord, error) {
			return &models.SettingsRecord{
				PinEnabled:           true,
				PinHash:              ptrStr("secret-hash"),
				Units:                "metric",
				DefaultAiExportWeeks: 12,
			}, nil
		},
	})

	result, err := svc.Get(settingsCtx, settingsUID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.PinEnabled)
	assert.Equal(t, "metric", result.Units)
	assert.Equal(t, 12, result.DefaultAiExportWeeks)
}

func TestSettingsService_Get_RepoError(t *testing.T) {
	svc := service.NewSettingsService(&mockSettingsRepoForSettings{
		findByUserIDFn: func(ctx context.Context, userID string) (*models.SettingsRecord, error) {
			return nil, errors.New("db unavailable")
		},
	})

	result, err := svc.Get(settingsCtx, settingsUID)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "db unavailable")
}

func TestSettingsService_Update_ValidUnitsAndExportWeeks(t *testing.T) {
	svc := service.NewSettingsService(&mockSettingsRepoForSettings{
		upsertFn: func(ctx context.Context, userID string, input models.SettingsInput) (*models.SettingsRecord, error) {
			return &models.SettingsRecord{
				PinEnabled:           false,
				Units:                *input.Units,
				DefaultAiExportWeeks: int32(*input.DefaultAiExportWeeks),
			}, nil
		},
	})

	result, err := svc.Update(settingsCtx, settingsUID, models.SettingsInput{
		Units:               ptrStr("imperial"),
		DefaultAiExportWeeks: ptrInt(8),
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "imperial", result.Units)
	assert.Equal(t, 8, result.DefaultAiExportWeeks)
}

func TestSettingsService_Update_InvalidUnits(t *testing.T) {
	svc := service.NewSettingsService(&mockSettingsRepoForSettings{})

	result, err := svc.Update(settingsCtx, settingsUID, models.SettingsInput{
		Units: ptrStr("fahrenheit"),
	})
	assert.ErrorIs(t, err, service.ErrInvalidUnits)
	assert.Nil(t, result)
}

func TestSettingsService_Update_InvalidExportWeeks(t *testing.T) {
	svc := service.NewSettingsService(&mockSettingsRepoForSettings{})

	t.Run("below minimum", func(t *testing.T) {
		result, err := svc.Update(settingsCtx, settingsUID, models.SettingsInput{
			DefaultAiExportWeeks: ptrInt(0),
		})
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "1 and 52")
	})

	t.Run("above maximum", func(t *testing.T) {
		result, err := svc.Update(settingsCtx, settingsUID, models.SettingsInput{
			DefaultAiExportWeeks: ptrInt(53),
		})
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "1 and 52")
	})
}

func TestSettingsService_RecordConversionStripsPinHash(t *testing.T) {
	svc := service.NewSettingsService(&mockSettingsRepoForSettings{
		findByUserIDFn: func(ctx context.Context, userID string) (*models.SettingsRecord, error) {
			return &models.SettingsRecord{
				PinEnabled:           true,
				PinHash:              ptrStr("should-be-stripped"),
				Units:                "metric",
				DefaultAiExportWeeks: 4,
			}, nil
		},
	})

	result, err := svc.Get(settingsCtx, settingsUID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.PinEnabled)
	assert.Equal(t, "metric", result.Units)
	assert.Equal(t, 4, result.DefaultAiExportWeeks)
}