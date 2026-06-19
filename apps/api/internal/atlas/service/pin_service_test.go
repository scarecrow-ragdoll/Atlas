// FILE: apps/api/internal/atlas/service/pin_service_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for the PIN service covering enable, disable, change, verify, and validation edge cases.
//   SCOPE: PIN enable/disable/change with correct and incorrect PINs, already-enabled/disabled states, PIN validation (length, digits), verify success/failure, IsEnabled, Argon2id hash verification.
//   DEPENDS: apps/api/internal/atlas/service, apps/api/internal/atlas/repository/postgres (mock SettingsRepository).
//   LINKS: M-API / V-M-API / TEST-W01-002.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added PIN service unit tests for WAVE-01.
// END_CHANGE_SUMMARY

package service_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasPostgres "monorepo-template/apps/api/internal/atlas/repository/postgres"
	atlasRedis "monorepo-template/apps/api/internal/atlas/repository/redis"
	"monorepo-template/apps/api/internal/atlas/service"
)

type mockSettingsRepo struct {
	atlasPostgres.SettingsRepository
	findByUserIDFn func(ctx context.Context, userID string) (*models.SettingsRecord, error)
	updatePinFn    func(ctx context.Context, userID string, pinEnabled bool, pinHash *string) error
}

func (m *mockSettingsRepo) FindByUserID(ctx context.Context, userID string) (*models.SettingsRecord, error) {
	return m.findByUserIDFn(ctx, userID)
}

func (m *mockSettingsRepo) UpdatePinState(ctx context.Context, userID string, pinEnabled bool, pinHash *string) error {
	return m.updatePinFn(ctx, userID, pinEnabled, pinHash)
}

type mockSessionStore struct {
	atlasRedis.PinSessionStore
	revokeAllFn func(ctx context.Context, userID string) error
}

func (m *mockSessionStore) RevokeAllByUser(ctx context.Context, userID string) error {
	return m.revokeAllFn(ctx, userID)
}

func ptr(s string) *string {
	return &s
}

var (
	defaultCtx   = context.Background()
	defaultUID   = "550e8400-e29b-41d4-a716-446655440000"
	argon2Params = service.Argon2Params{
		Memory:      64,
		Iterations:  1,
		Parallelism: 1,
		KeyLength:   32,
	}
)

func TestPinService_Enable_Success(t *testing.T) {
	var storedHash *string
	svc := service.NewPinService(
		&mockSettingsRepo{
			findByUserIDFn: func(ctx context.Context, userID string) (*models.SettingsRecord, error) {
				return &models.SettingsRecord{PinEnabled: false}, nil
			},
			updatePinFn: func(ctx context.Context, userID string, pinEnabled bool, pinHash *string) error {
				storedHash = pinHash
				return nil
			},
		},
		&mockSessionStore{},
		argon2Params, 4, 20,
	)

	err := svc.Enable(defaultCtx, defaultUID, "1234")
	require.NoError(t, err)
	assert.NotNil(t, storedHash)
	assert.True(t, strings.Contains(*storedHash, "$"))
}

func TestPinService_Enable_AlreadyEnabled(t *testing.T) {
	svc := service.NewPinService(
		&mockSettingsRepo{
			findByUserIDFn: func(ctx context.Context, userID string) (*models.SettingsRecord, error) {
				return &models.SettingsRecord{PinEnabled: true, PinHash: ptr("existing")}, nil
			},
		},
		&mockSessionStore{},
		argon2Params, 4, 20,
	)

	err := svc.Enable(defaultCtx, defaultUID, "1234")
	assert.ErrorIs(t, err, service.ErrPinAlreadyEnabled)
}

func TestPinService_Enable_TooShort(t *testing.T) {
	svc := service.NewPinService(
		&mockSettingsRepo{},
		&mockSessionStore{},
		argon2Params, 4, 20,
	)

	err := svc.Enable(defaultCtx, defaultUID, "12")
	assert.ErrorIs(t, err, service.ErrPinTooShort)
}

func TestPinService_Enable_TooLong(t *testing.T) {
	svc := service.NewPinService(
		&mockSettingsRepo{},
		&mockSessionStore{},
		argon2Params, 4, 20,
	)

	err := svc.Enable(defaultCtx, defaultUID, "123456789012345678901")
	assert.ErrorIs(t, err, service.ErrPinTooLong)
}

func TestPinService_Enable_NonDigits(t *testing.T) {
	svc := service.NewPinService(
		&mockSettingsRepo{},
		&mockSessionStore{},
		argon2Params, 4, 20,
	)

	err := svc.Enable(defaultCtx, defaultUID, "abcd")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "digits")
}

func TestPinService_Disable_Success(t *testing.T) {
	pin := "1234"
	svc := service.NewPinService(
		&mockSettingsRepo{
			findByUserIDFn: func(ctx context.Context, userID string) (*models.SettingsRecord, error) {
				return &models.SettingsRecord{PinEnabled: true, PinHash: ptr(hashPinForTest(pin))}, nil
			},
			updatePinFn: func(ctx context.Context, userID string, pinEnabled bool, pinHash *string) error {
				assert.False(t, pinEnabled)
				assert.Nil(t, pinHash)
				return nil
			},
		},
		&mockSessionStore{
			revokeAllFn: func(ctx context.Context, userID string) error {
				return nil
			},
		},
		argon2Params, 4, 20,
	)

	err := svc.Disable(defaultCtx, defaultUID, pin)
	assert.NoError(t, err)
}

func TestPinService_Disable_WrongPin(t *testing.T) {
	svc := service.NewPinService(
		&mockSettingsRepo{
			findByUserIDFn: func(ctx context.Context, userID string) (*models.SettingsRecord, error) {
				return &models.SettingsRecord{PinEnabled: true, PinHash: ptr(hashPinForTest("1234"))}, nil
			},
		},
		&mockSessionStore{},
		argon2Params, 4, 20,
	)

	err := svc.Disable(defaultCtx, defaultUID, "9999")
	assert.ErrorIs(t, err, service.ErrPinWrongPin)
}

func TestPinService_Disable_AlreadyDisabled(t *testing.T) {
	svc := service.NewPinService(
		&mockSettingsRepo{
			findByUserIDFn: func(ctx context.Context, userID string) (*models.SettingsRecord, error) {
				return &models.SettingsRecord{PinEnabled: false, PinHash: nil}, nil
			},
		},
		&mockSessionStore{},
		argon2Params, 4, 20,
	)

	err := svc.Disable(defaultCtx, defaultUID, "1234")
	assert.ErrorIs(t, err, service.ErrPinAlreadyDisabled)
}

func TestPinService_Change_Success(t *testing.T) {
	oldPin := "1234"
	var storedHash *string
	svc := service.NewPinService(
		&mockSettingsRepo{
			findByUserIDFn: func(ctx context.Context, userID string) (*models.SettingsRecord, error) {
				return &models.SettingsRecord{PinEnabled: true, PinHash: ptr(hashPinForTest(oldPin))}, nil
			},
			updatePinFn: func(ctx context.Context, userID string, pinEnabled bool, pinHash *string) error {
				storedHash = pinHash
				return nil
			},
		},
		&mockSessionStore{
			revokeAllFn: func(ctx context.Context, userID string) error {
				return nil
			},
		},
		argon2Params, 4, 20,
	)

	err := svc.Change(defaultCtx, defaultUID, oldPin, "5678")
	require.NoError(t, err)
	assert.NotNil(t, storedHash)
}

func TestPinService_Change_WrongCurrentPin(t *testing.T) {
	svc := service.NewPinService(
		&mockSettingsRepo{
			findByUserIDFn: func(ctx context.Context, userID string) (*models.SettingsRecord, error) {
				return &models.SettingsRecord{PinEnabled: true, PinHash: ptr(hashPinForTest("1234"))}, nil
			},
		},
		&mockSessionStore{},
		argon2Params, 4, 20,
	)

	err := svc.Change(defaultCtx, defaultUID, "9999", "5678")
	assert.ErrorIs(t, err, service.ErrPinWrongPin)
}

func TestPinService_Verify_Correct(t *testing.T) {
	svc := service.NewPinService(
		&mockSettingsRepo{
			findByUserIDFn: func(ctx context.Context, userID string) (*models.SettingsRecord, error) {
				return &models.SettingsRecord{PinHash: ptr(hashPinForTest("1234"))}, nil
			},
		},
		&mockSessionStore{},
		argon2Params, 4, 20,
	)

	valid, err := svc.Verify(defaultCtx, defaultUID, "1234")
	require.NoError(t, err)
	assert.True(t, valid)
}

func TestPinService_Verify_Incorrect(t *testing.T) {
	svc := service.NewPinService(
		&mockSettingsRepo{
			findByUserIDFn: func(ctx context.Context, userID string) (*models.SettingsRecord, error) {
				return &models.SettingsRecord{PinHash: ptr(hashPinForTest("1234"))}, nil
			},
		},
		&mockSessionStore{},
		argon2Params, 4, 20,
	)

	valid, err := svc.Verify(defaultCtx, defaultUID, "9999")
	require.NoError(t, err)
	assert.False(t, valid)
}

func TestPinService_Verify_NoHash(t *testing.T) {
	svc := service.NewPinService(
		&mockSettingsRepo{
			findByUserIDFn: func(ctx context.Context, userID string) (*models.SettingsRecord, error) {
				return &models.SettingsRecord{PinHash: nil}, nil
			},
		},
		&mockSessionStore{},
		argon2Params, 4, 20,
	)

	valid, err := svc.Verify(defaultCtx, defaultUID, "1234")
	require.NoError(t, err)
	assert.False(t, valid)
}

func TestPinService_IsEnabled(t *testing.T) {
	svc := service.NewPinService(
		&mockSettingsRepo{
			findByUserIDFn: func(ctx context.Context, userID string) (*models.SettingsRecord, error) {
				return &models.SettingsRecord{PinEnabled: true}, nil
			},
		},
		&mockSessionStore{},
		argon2Params, 4, 20,
	)

	enabled, err := svc.IsEnabled(defaultCtx, defaultUID)
	require.NoError(t, err)
	assert.True(t, enabled)
}

func TestPinService_IsEnabled_False(t *testing.T) {
	svc := service.NewPinService(
		&mockSettingsRepo{
			findByUserIDFn: func(ctx context.Context, userID string) (*models.SettingsRecord, error) {
				return &models.SettingsRecord{PinEnabled: false}, nil
			},
		},
		&mockSessionStore{},
		argon2Params, 4, 20,
	)

	enabled, err := svc.IsEnabled(defaultCtx, defaultUID)
	require.NoError(t, err)
	assert.False(t, enabled)
}

func hashPinForTest(pin string) string {
	return service.HashPinForTest(pin)
}