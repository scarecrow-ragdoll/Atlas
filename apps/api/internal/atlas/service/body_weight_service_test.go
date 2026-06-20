// FILE: apps/api/internal/atlas/service/body_weight_service_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for BodyWeightService covering Create, GetByID, Update, Delete, and Latest operations with validation.
//   SCOPE: Success paths, validation errors (zero weight, invalid source), not-found, Latest returning nil when empty, Latest returning entry.
//   DEPENDS: apps/api/internal/atlas/service, apps/api/internal/atlas/repository/postgres (mock BodyWeightEntryRepository), apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added body weight service unit tests for WAVE-04.
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

type mockBodyWeightRepo struct {
	atlasPostgres.BodyWeightEntryRepository
	createFn          func(ctx context.Context, userID string, date models.Date, weight float64, source string, notes *string) (*models.BodyWeightRecord, error)
	getByIDFn         func(ctx context.Context, userID string, id string) (*models.BodyWeightRecord, error)
	listByDateRangeFn func(ctx context.Context, userID string, fromDate models.Date, toDate models.Date) ([]models.BodyWeightRecord, error)
	latestFn          func(ctx context.Context, userID string) (*models.BodyWeightRecord, error)
	updateFn          func(ctx context.Context, userID string, id string, weight *float64, source *string, notes *string) (*models.BodyWeightRecord, error)
	deleteFn          func(ctx context.Context, userID string, id string) (*models.BodyWeightRecord, error)
}

func (m *mockBodyWeightRepo) Create(ctx context.Context, userID string, date models.Date, weight float64, source string, notes *string) (*models.BodyWeightRecord, error) {
	return m.createFn(ctx, userID, date, weight, source, notes)
}

func (m *mockBodyWeightRepo) GetByID(ctx context.Context, userID string, id string) (*models.BodyWeightRecord, error) {
	return m.getByIDFn(ctx, userID, id)
}

func (m *mockBodyWeightRepo) ListByDateRange(ctx context.Context, userID string, fromDate models.Date, toDate models.Date) ([]models.BodyWeightRecord, error) {
	return m.listByDateRangeFn(ctx, userID, fromDate, toDate)
}

func (m *mockBodyWeightRepo) Latest(ctx context.Context, userID string) (*models.BodyWeightRecord, error) {
	return m.latestFn(ctx, userID)
}

func (m *mockBodyWeightRepo) Update(ctx context.Context, userID string, id string, weight *float64, source *string, notes *string) (*models.BodyWeightRecord, error) {
	return m.updateFn(ctx, userID, id, weight, source, notes)
}

func (m *mockBodyWeightRepo) Delete(ctx context.Context, userID string, id string) (*models.BodyWeightRecord, error) {
	return m.deleteFn(ctx, userID, id)
}

var (
	bwTestDate = models.MustDate("2026-06-20")
)

func bwTestRecord() *models.BodyWeightRecord {
	return &models.BodyWeightRecord{
		ID:        testID,
		UserID:    testUserID,
		Date:      bwTestDate,
		Weight:    75.5,
		Source:    "SCALE",
		Notes:     ptrStr("Morning weight"),
		CreatedAt: "2026-06-20T00:00:00Z",
		UpdatedAt: "2026-06-20T00:00:00Z",
	}
}

// ----- Create -----

func TestBodyWeightService_Create_WeightZero(t *testing.T) {
	svc := service.NewBodyWeightService(&mockBodyWeightRepo{})

	entry, err := svc.Create(ctx, testUserID, models.CreateBodyWeightInput{
		Date:   bwTestDate,
		Weight: 0,
		Source: models.BodyWeightSourceScale,
	})
	assert.ErrorIs(t, err, service.ErrBodyWeightInvalid)
	assert.Nil(t, entry)
}

func TestBodyWeightService_Create_InvalidSource(t *testing.T) {
	svc := service.NewBodyWeightService(&mockBodyWeightRepo{})

	entry, err := svc.Create(ctx, testUserID, models.CreateBodyWeightInput{
		Date:   bwTestDate,
		Weight: 75.5,
		Source: "INVALID_SOURCE",
	})
	assert.ErrorIs(t, err, service.ErrBodyWeightInvalidSource)
	assert.Nil(t, entry)
}

func TestBodyWeightService_Create_Success(t *testing.T) {
	svc := service.NewBodyWeightService(&mockBodyWeightRepo{
		createFn: func(ctx context.Context, userID string, date models.Date, weight float64, source string, notes *string) (*models.BodyWeightRecord, error) {
			assert.Equal(t, 75.5, weight)
			assert.Equal(t, "SCALE", source)
			return bwTestRecord(), nil
		},
	})

	entry, err := svc.Create(ctx, testUserID, models.CreateBodyWeightInput{
		Date:   bwTestDate,
		Weight: 75.5,
		Source: models.BodyWeightSourceScale,
		Notes:  ptrStr("Morning weight"),
	})
	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.Equal(t, 75.5, entry.Weight)
	assert.Equal(t, models.BodyWeightSourceScale, entry.Source)
}

// ----- GetByID -----

func TestBodyWeightService_GetByID_Success(t *testing.T) {
	svc := service.NewBodyWeightService(&mockBodyWeightRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.BodyWeightRecord, error) {
			return bwTestRecord(), nil
		},
	})

	entry, err := svc.GetByID(ctx, testUserID, testID)
	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.Equal(t, 75.5, entry.Weight)
}

func TestBodyWeightService_GetByID_NotFound(t *testing.T) {
	svc := service.NewBodyWeightService(&mockBodyWeightRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.BodyWeightRecord, error) {
			return nil, nil
		},
	})

	entry, err := svc.GetByID(ctx, testUserID, testID)
	assert.ErrorIs(t, err, service.ErrBodyWeightNotFound)
	assert.Nil(t, entry)
}

// ----- Update -----

func TestBodyWeightService_Update_WeightZero(t *testing.T) {
	svc := service.NewBodyWeightService(&mockBodyWeightRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.BodyWeightRecord, error) {
			return bwTestRecord(), nil
		},
	})

	entry, err := svc.Update(ctx, testUserID, testID, models.UpdateBodyWeightInput{
		Weight: ptrFloat64(0),
	})
	assert.ErrorIs(t, err, service.ErrBodyWeightInvalid)
	assert.Nil(t, entry)
}

func TestBodyWeightService_Update_Success(t *testing.T) {
	svc := service.NewBodyWeightService(&mockBodyWeightRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.BodyWeightRecord, error) {
			return bwTestRecord(), nil
		},
		updateFn: func(ctx context.Context, userID string, id string, weight *float64, source *string, notes *string) (*models.BodyWeightRecord, error) {
			assert.Equal(t, 80.0, *weight)
			return &models.BodyWeightRecord{
				ID:        id,
				UserID:    userID,
				Date:      bwTestDate,
				Weight:    *weight,
				Source:    *source,
				Notes:     notes,
				CreatedAt: "2026-06-20T00:00:00Z",
				UpdatedAt: "2026-06-20T12:00:00Z",
			}, nil
		},
	})

	entry, err := svc.Update(ctx, testUserID, testID, models.UpdateBodyWeightInput{
		Weight: ptrFloat64(80.0),
		Source: sourcePtr(models.BodyWeightSourceManual),
	})
	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.Equal(t, 80.0, entry.Weight)
}

// ----- Delete -----

func TestBodyWeightService_Delete_Success(t *testing.T) {
	svc := service.NewBodyWeightService(&mockBodyWeightRepo{
		deleteFn: func(ctx context.Context, userID string, id string) (*models.BodyWeightRecord, error) {
			return bwTestRecord(), nil
		},
	})

	entry, err := svc.Delete(ctx, testUserID, testID)
	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.Equal(t, testID, entry.ID)
}

// ----- Latest -----

func TestBodyWeightService_Latest_Success(t *testing.T) {
	svc := service.NewBodyWeightService(&mockBodyWeightRepo{
		latestFn: func(ctx context.Context, userID string) (*models.BodyWeightRecord, error) {
			return nil, nil
		},
	})

	entry, err := svc.Latest(ctx, testUserID)
	require.NoError(t, err)
	assert.Nil(t, entry)
}

func TestBodyWeightService_Latest_ReturnsEntry(t *testing.T) {
	svc := service.NewBodyWeightService(&mockBodyWeightRepo{
		latestFn: func(ctx context.Context, userID string) (*models.BodyWeightRecord, error) {
			return bwTestRecord(), nil
		},
	})

	entry, err := svc.Latest(ctx, testUserID)
	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.Equal(t, 75.5, entry.Weight)
}

// ----- Helpers -----

func sourcePtr(s models.BodyWeightSource) *models.BodyWeightSource {
	return &s
}