// FILE: apps/api/internal/atlas/service/cardio_service_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for CardioService covering Create, GetByID, Update, Delete with validation and DailyLog auto-creation.
//   SCOPE: Success paths, validation errors (invalid type, zero duration, negative pulse, invalid zone), not-found, existing DailyLog reuse.
//   DEPENDS: apps/api/internal/atlas/service, apps/api/internal/atlas/repository/postgres (mock CardioEntryRepository, DailyLogRepository), apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added cardio service unit tests for WAVE-04.
// END_CHANGE_SUMMARY

package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasPostgres "monorepo-template/apps/api/internal/atlas/repository/postgres"
	"monorepo-template/apps/api/internal/atlas/service"
)

type mockCardioRepo struct {
	atlasPostgres.CardioEntryRepository
	createFn        func(ctx context.Context, userID string, dailyLogID string, cardioType string, durationMinutes int32, avgPulse *int32, heartRateZone *string, notes *string) (*models.CardioRecord, error)
	getByIDFn       func(ctx context.Context, userID string, id string) (*models.CardioRecord, error)
	listByDailyLogFn func(ctx context.Context, userID string, dailyLogID string) ([]models.CardioRecord, error)
	updateFn        func(ctx context.Context, userID string, id string, cardioType string, durationMinutes int32, avgPulse *int32, heartRateZone *string, notes *string) (*models.CardioRecord, error)
	deleteFn        func(ctx context.Context, userID string, id string) (*models.CardioRecord, error)
}

func (m *mockCardioRepo) Create(ctx context.Context, userID string, dailyLogID string, cardioType string, durationMinutes int32, avgPulse *int32, heartRateZone *string, notes *string) (*models.CardioRecord, error) {
	return m.createFn(ctx, userID, dailyLogID, cardioType, durationMinutes, avgPulse, heartRateZone, notes)
}

func (m *mockCardioRepo) GetByID(ctx context.Context, userID string, id string) (*models.CardioRecord, error) {
	return m.getByIDFn(ctx, userID, id)
}

func (m *mockCardioRepo) ListByDailyLog(ctx context.Context, userID string, dailyLogID string) ([]models.CardioRecord, error) {
	return m.listByDailyLogFn(ctx, userID, dailyLogID)
}

func (m *mockCardioRepo) Update(ctx context.Context, userID string, id string, cardioType string, durationMinutes int32, avgPulse *int32, heartRateZone *string, notes *string) (*models.CardioRecord, error) {
	return m.updateFn(ctx, userID, id, cardioType, durationMinutes, avgPulse, heartRateZone, notes)
}

func (m *mockCardioRepo) Delete(ctx context.Context, userID string, id string) (*models.CardioRecord, error) {
	return m.deleteFn(ctx, userID, id)
}

type mockCardioDailyLogRepo struct {
	atlasPostgres.DailyLogRepository
	getDailyLogByDateFn         func(ctx context.Context, userID string, date models.Date) (*models.DailyLogRecord, error)
	getOrCreateDailyLogByDateFn func(ctx context.Context, userID string, date models.Date) (*models.DailyLogRecord, error)
}

func (m *mockCardioDailyLogRepo) GetDailyLogByDate(ctx context.Context, userID string, date models.Date) (*models.DailyLogRecord, error) {
	return m.getDailyLogByDateFn(ctx, userID, date)
}

func (m *mockCardioDailyLogRepo) GetOrCreateDailyLogByDate(ctx context.Context, userID string, date models.Date) (*models.DailyLogRecord, error) {
	return m.getOrCreateDailyLogByDateFn(ctx, userID, date)
}

var (
	cardioTestDate = models.MustDate("2026-06-20")
)

func cardioTestDailyLog() *models.DailyLogRecord {
	return &models.DailyLogRecord{
		ID:        "11111111-1111-1111-1111-111111111111",
		UserID:    testUserID,
		Date:      cardioTestDate,
		CreatedAt: "2026-06-20T00:00:00Z",
		UpdatedAt: "2026-06-20T00:00:00Z",
	}
}

func cardioTestRecord(dailyLogID string) *models.CardioRecord {
	return &models.CardioRecord{
		ID:              testID,
		UserID:          testUserID,
		DailyLogID:      dailyLogID,
		CardioType:      "RUNNING",
		DurationMinutes: 30,
		AvgPulse:        ptrInt32(145),
		HeartRateZone:   ptrStr("ZONE_3"),
		Notes:           ptrStr("Morning run"),
		CreatedAt:       "2026-06-20T00:00:00Z",
		UpdatedAt:       "2026-06-20T00:00:00Z",
	}
}

// ----- Create -----

func TestCardioService_Create_InvalidType(t *testing.T) {
	svc := service.NewCardioService(&mockCardioRepo{}, &mockCardioDailyLogRepo{})

	entry, err := svc.Create(ctx, testUserID, models.CreateCardioInput{
		Date:            cardioTestDate,
		CardioType:      "INVALID",
		DurationMinutes: 30,
	})
	assert.ErrorIs(t, err, service.ErrCardioInvalidType)
	assert.Nil(t, entry)
}

func TestCardioService_Create_DurationZero(t *testing.T) {
	svc := service.NewCardioService(&mockCardioRepo{}, &mockCardioDailyLogRepo{})

	entry, err := svc.Create(ctx, testUserID, models.CreateCardioInput{
		Date:            cardioTestDate,
		CardioType:      models.CardioTypeRunning,
		DurationMinutes: 0,
	})
	assert.ErrorIs(t, err, service.ErrCardioDurationInvalid)
	assert.Nil(t, entry)
}

func TestCardioService_Create_NegativePulse(t *testing.T) {
	svc := service.NewCardioService(&mockCardioRepo{}, &mockCardioDailyLogRepo{})

	entry, err := svc.Create(ctx, testUserID, models.CreateCardioInput{
		Date:            cardioTestDate,
		CardioType:      models.CardioTypeRunning,
		DurationMinutes: 30,
		AvgPulse:        ptrInt32(-5),
	})
	assert.ErrorIs(t, err, service.ErrCardioPulseInvalid)
	assert.Nil(t, entry)
}

func TestCardioService_Create_InvalidZone(t *testing.T) {
	svc := service.NewCardioService(&mockCardioRepo{}, &mockCardioDailyLogRepo{})

	entry, err := svc.Create(ctx, testUserID, models.CreateCardioInput{
		Date:            cardioTestDate,
		CardioType:      models.CardioTypeRunning,
		DurationMinutes: 30,
		HeartRateZone:   zonePtr("INVALID_ZONE"),
	})
	assert.ErrorIs(t, err, service.ErrCardioZoneInvalid)
	assert.Nil(t, entry)
}

func TestCardioService_Create_Success(t *testing.T) {
	dailyLog := cardioTestDailyLog()
	svc := service.NewCardioService(&mockCardioRepo{
		createFn: func(ctx context.Context, userID string, dailyLogID string, cardioType string, durationMinutes int32, avgPulse *int32, heartRateZone *string, notes *string) (*models.CardioRecord, error) {
			assert.Equal(t, "RUNNING", cardioType)
			assert.Equal(t, int32(30), durationMinutes)
			return cardioTestRecord(dailyLogID), nil
		},
	}, &mockCardioDailyLogRepo{
		getDailyLogByDateFn: func(ctx context.Context, userID string, date models.Date) (*models.DailyLogRecord, error) {
			return nil, nil
		},
		getOrCreateDailyLogByDateFn: func(ctx context.Context, userID string, date models.Date) (*models.DailyLogRecord, error) {
			return dailyLog, nil
		},
	})

	entry, err := svc.Create(ctx, testUserID, models.CreateCardioInput{
		Date:            cardioTestDate,
		CardioType:      models.CardioTypeRunning,
		DurationMinutes: 30,
		AvgPulse:        ptrInt32(145),
		HeartRateZone:   zonePtr("ZONE_3"),
		Notes:           ptrStr("Morning run"),
	})
	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.Equal(t, models.CardioTypeRunning, entry.CardioType)
	assert.Equal(t, int32(30), entry.DurationMinutes)
}

func TestCardioService_Create_ExistingDailyLog(t *testing.T) {
	dailyLog := cardioTestDailyLog()
	svc := service.NewCardioService(&mockCardioRepo{
		createFn: func(ctx context.Context, userID string, dailyLogID string, cardioType string, durationMinutes int32, avgPulse *int32, heartRateZone *string, notes *string) (*models.CardioRecord, error) {
			assert.Equal(t, dailyLog.ID, dailyLogID)
			return cardioTestRecord(dailyLogID), nil
		},
	}, &mockCardioDailyLogRepo{
		getDailyLogByDateFn: func(ctx context.Context, userID string, date models.Date) (*models.DailyLogRecord, error) {
			return dailyLog, nil
		},
	})

	entry, err := svc.Create(ctx, testUserID, models.CreateCardioInput{
		Date:            cardioTestDate,
		CardioType:      models.CardioTypeRunning,
		DurationMinutes: 30,
	})
	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.Equal(t, dailyLog.ID, entry.DailyLogID)
}

// ----- GetByID -----

func TestCardioService_GetByID_Success(t *testing.T) {
	svc := service.NewCardioService(&mockCardioRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.CardioRecord, error) {
			return cardioTestRecord("daily-log-id"), nil
		},
	}, &mockCardioDailyLogRepo{})

	entry, err := svc.GetByID(ctx, testUserID, testID)
	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.Equal(t, models.CardioTypeRunning, entry.CardioType)
}

func TestCardioService_GetByID_NotFound(t *testing.T) {
	svc := service.NewCardioService(&mockCardioRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.CardioRecord, error) {
			return nil, nil
		},
	}, &mockCardioDailyLogRepo{})

	entry, err := svc.GetByID(ctx, testUserID, testID)
	assert.ErrorIs(t, err, service.ErrCardioNotFound)
	assert.Nil(t, entry)
}

// ----- Update -----

func TestCardioService_Update_InvalidType(t *testing.T) {
	svc := service.NewCardioService(&mockCardioRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.CardioRecord, error) {
			return cardioTestRecord("daily-log-id"), nil
		},
	}, &mockCardioDailyLogRepo{})

	entry, err := svc.Update(ctx, testUserID, testID, models.UpdateCardioInput{
		CardioType: typePtr("INVALID"),
	})
	assert.ErrorIs(t, err, service.ErrCardioInvalidType)
	assert.Nil(t, entry)
}

func TestCardioService_Update_Success(t *testing.T) {
	svc := service.NewCardioService(&mockCardioRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.CardioRecord, error) {
			return cardioTestRecord("daily-log-id"), nil
		},
		updateFn: func(ctx context.Context, userID string, id string, cardioType string, durationMinutes int32, avgPulse *int32, heartRateZone *string, notes *string) (*models.CardioRecord, error) {
			assert.Equal(t, "WALKING", cardioType)
			assert.Equal(t, int32(45), durationMinutes)
			return &models.CardioRecord{
				ID:              id,
				UserID:          userID,
				DailyLogID:      "daily-log-id",
				CardioType:      "WALKING",
				DurationMinutes: 45,
				AvgPulse:        ptrInt32(120),
				HeartRateZone:   ptrStr("ZONE_2"),
				Notes:           ptrStr("Evening walk"),
				CreatedAt:       "2026-06-20T00:00:00Z",
				UpdatedAt:       "2026-06-20T12:00:00Z",
			}, nil
		},
	}, &mockCardioDailyLogRepo{})

	entry, err := svc.Update(ctx, testUserID, testID, models.UpdateCardioInput{
		CardioType:      typePtr("WALKING"),
		DurationMinutes: ptrInt32(45),
		AvgPulse:        ptrInt32(120),
		HeartRateZone:   zonePtr("ZONE_2"),
		Notes:           ptrStr("Evening walk"),
	})
	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.Equal(t, models.CardioTypeWalking, entry.CardioType)
	assert.Equal(t, int32(45), entry.DurationMinutes)
}

// ----- Delete -----

func TestCardioService_Delete_Success(t *testing.T) {
	svc := service.NewCardioService(&mockCardioRepo{
		deleteFn: func(ctx context.Context, userID string, id string) (*models.CardioRecord, error) {
			return cardioTestRecord("daily-log-id"), nil
		},
	}, &mockCardioDailyLogRepo{})

	entry, err := svc.Delete(ctx, testUserID, testID)
	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.Equal(t, testID, entry.ID)
}

func TestCardioService_Delete_NotFound(t *testing.T) {
	svc := service.NewCardioService(&mockCardioRepo{
		deleteFn: func(ctx context.Context, userID string, id string) (*models.CardioRecord, error) {
			return nil, nil
		},
	}, &mockCardioDailyLogRepo{})

	entry, err := svc.Delete(ctx, testUserID, testID)
	assert.ErrorIs(t, err, service.ErrCardioNotFound)
	assert.Nil(t, entry)
}

// ----- Helpers -----

func ptrInt32(i int32) *int32 { return &i }

func zonePtr(z string) *models.HeartRateZone {
	v := models.HeartRateZone(z)
	return &v
}

func typePtr(tp string) *models.CardioType {
	v := models.CardioType(tp)
	return &v
}