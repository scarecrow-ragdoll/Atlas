// FILE: apps/api/internal/atlas/service/week_flag_service_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for WeekFlagService covering Create, ListByWeekStart, and Delete operations with validation.
//   SCOPE: Success paths, validation error (invalid flag type), not-found on delete, list by week start date.
//   DEPENDS: apps/api/internal/atlas/service, apps/api/internal/atlas/repository/postgres (mock WeekFlagRepository), apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added week flag service unit tests for WAVE-04.
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

type mockWeekFlagRepo struct {
	atlasPostgres.WeekFlagRepository
	createFn          func(ctx context.Context, userID string, weekStartDate models.Date, flagType string, notes *string) (*models.WeekFlagRecord, error)
	getByIDFn         func(ctx context.Context, userID string, id string) (*models.WeekFlagRecord, error)
	listByWeekStartFn func(ctx context.Context, userID string, weekStartDate models.Date) ([]models.WeekFlagRecord, error)
	deleteFn          func(ctx context.Context, userID string, id string) (*models.WeekFlagRecord, error)
}

func (m *mockWeekFlagRepo) Create(ctx context.Context, userID string, weekStartDate models.Date, flagType string, notes *string) (*models.WeekFlagRecord, error) {
	return m.createFn(ctx, userID, weekStartDate, flagType, notes)
}

func (m *mockWeekFlagRepo) GetByID(ctx context.Context, userID string, id string) (*models.WeekFlagRecord, error) {
	return m.getByIDFn(ctx, userID, id)
}

func (m *mockWeekFlagRepo) ListByWeekStart(ctx context.Context, userID string, weekStartDate models.Date) ([]models.WeekFlagRecord, error) {
	return m.listByWeekStartFn(ctx, userID, weekStartDate)
}

func (m *mockWeekFlagRepo) Delete(ctx context.Context, userID string, id string) (*models.WeekFlagRecord, error) {
	return m.deleteFn(ctx, userID, id)
}

var (
	wfTestWeekStart = models.MustDate("2026-06-15")
)

func wfTestRecord() *models.WeekFlagRecord {
	return &models.WeekFlagRecord{
		ID:            testID,
		UserID:        testUserID,
		WeekStartDate: wfTestWeekStart,
		FlagType:      "POOR_SLEEP",
		Notes:         ptrStr("Had trouble sleeping"),
		CreatedAt:     "2026-06-15T00:00:00Z",
		UpdatedAt:     "2026-06-15T00:00:00Z",
	}
}

// ----- Create -----

func TestWeekFlagService_Create_InvalidType(t *testing.T) {
	svc := service.NewWeekFlagService(&mockWeekFlagRepo{})

	flag, err := svc.Create(ctx, testUserID, models.CreateWeekFlagInput{
		WeekStartDate: wfTestWeekStart,
		FlagType:      "INVALID_FLAG",
	})
	assert.ErrorIs(t, err, service.ErrWeekFlagInvalidType)
	assert.Nil(t, flag)
}

func TestWeekFlagService_Create_Success(t *testing.T) {
	svc := service.NewWeekFlagService(&mockWeekFlagRepo{
		createFn: func(ctx context.Context, userID string, weekStartDate models.Date, flagType string, notes *string) (*models.WeekFlagRecord, error) {
			assert.Equal(t, "POOR_SLEEP", flagType)
			return wfTestRecord(), nil
		},
	})

	flag, err := svc.Create(ctx, testUserID, models.CreateWeekFlagInput{
		WeekStartDate: wfTestWeekStart,
		FlagType:      models.WeekFlagTypePoorSleep,
		Notes:         ptrStr("Had trouble sleeping"),
	})
	require.NoError(t, err)
	require.NotNil(t, flag)
	assert.Equal(t, models.WeekFlagTypePoorSleep, flag.FlagType)
}

// ----- ListByWeekStart -----

func TestWeekFlagService_ListByWeekStart_Success(t *testing.T) {
	svc := service.NewWeekFlagService(&mockWeekFlagRepo{
		listByWeekStartFn: func(ctx context.Context, userID string, weekStartDate models.Date) ([]models.WeekFlagRecord, error) {
			return []models.WeekFlagRecord{*wfTestRecord()}, nil
		},
	})

	flags, err := svc.ListByWeekStart(ctx, testUserID, wfTestWeekStart)
	require.NoError(t, err)
	require.Len(t, flags, 1)
	assert.Equal(t, models.WeekFlagTypePoorSleep, flags[0].FlagType)
}

// ----- Delete -----

func TestWeekFlagService_Delete_Success(t *testing.T) {
	svc := service.NewWeekFlagService(&mockWeekFlagRepo{
		deleteFn: func(ctx context.Context, userID string, id string) (*models.WeekFlagRecord, error) {
			return wfTestRecord(), nil
		},
	})

	flag, err := svc.Delete(ctx, testUserID, testID)
	require.NoError(t, err)
	require.NotNil(t, flag)
	assert.Equal(t, testID, flag.ID)
}

func TestWeekFlagService_Delete_NotFound(t *testing.T) {
	svc := service.NewWeekFlagService(&mockWeekFlagRepo{
		deleteFn: func(ctx context.Context, userID string, id string) (*models.WeekFlagRecord, error) {
			return nil, nil
		},
	})

	flag, err := svc.Delete(ctx, testUserID, testID)
	assert.ErrorIs(t, err, service.ErrWeekFlagNotFound)
	assert.Nil(t, flag)
}