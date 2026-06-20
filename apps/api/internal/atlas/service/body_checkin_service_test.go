// FILE: apps/api/internal/atlas/service/body_checkin_service_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for BodyCheckInService and BodyMeasurementService covering Create, GetByID, Update, Delete operations with validation.
//   SCOPE: Check-in and measurement success paths, validation errors (zero weight, invalid body fat, invalid type, zero value, side on unpaired type), not-found, nested measurements and photos in GetByID.
//   DEPENDS: apps/api/internal/atlas/service, apps/api/internal/atlas/repository/postgres (mock BodyCheckInRepository, BodyMeasurementRepository, ProgressPhotoRepository), apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added body check-in and measurement service unit tests for WAVE-04.
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

type mockCheckInRepo struct {
	atlasPostgres.BodyCheckInRepository
	createFn          func(ctx context.Context, userID string, date models.Date, weight *float64, bodyFatPercentage *float64, notes *string) (*models.BodyCheckInRecord, error)
	getByIDFn         func(ctx context.Context, userID string, id string) (*models.BodyCheckInRecord, error)
	listByDateRangeFn func(ctx context.Context, userID string, fromDate models.Date, toDate models.Date) ([]models.BodyCheckInRecord, error)
	updateFn          func(ctx context.Context, userID string, id string, weight *float64, bodyFatPercentage *float64, notes *string) (*models.BodyCheckInRecord, error)
	deleteFn          func(ctx context.Context, userID string, id string) (*models.BodyCheckInRecord, error)
}

func (m *mockCheckInRepo) Create(ctx context.Context, userID string, date models.Date, weight *float64, bodyFatPercentage *float64, notes *string) (*models.BodyCheckInRecord, error) {
	return m.createFn(ctx, userID, date, weight, bodyFatPercentage, notes)
}

func (m *mockCheckInRepo) GetByID(ctx context.Context, userID string, id string) (*models.BodyCheckInRecord, error) {
	return m.getByIDFn(ctx, userID, id)
}

func (m *mockCheckInRepo) ListByDateRange(ctx context.Context, userID string, fromDate models.Date, toDate models.Date) ([]models.BodyCheckInRecord, error) {
	return m.listByDateRangeFn(ctx, userID, fromDate, toDate)
}

func (m *mockCheckInRepo) Update(ctx context.Context, userID string, id string, weight *float64, bodyFatPercentage *float64, notes *string) (*models.BodyCheckInRecord, error) {
	return m.updateFn(ctx, userID, id, weight, bodyFatPercentage, notes)
}

func (m *mockCheckInRepo) Delete(ctx context.Context, userID string, id string) (*models.BodyCheckInRecord, error) {
	return m.deleteFn(ctx, userID, id)
}

type mockMeasurementRepo struct {
	atlasPostgres.BodyMeasurementRepository
	createFn        func(ctx context.Context, checkInID string, measurementType string, side *string, value float64) (*models.BodyMeasurementRecord, error)
	getByIDFn       func(ctx context.Context, userID string, id string) (*models.BodyMeasurementRecord, error)
	listByCheckInFn func(ctx context.Context, userID string, checkInID string) ([]models.BodyMeasurementRecord, error)
	updateFn        func(ctx context.Context, userID string, id string, measurementType string, side *string, value float64) (*models.BodyMeasurementRecord, error)
	deleteFn        func(ctx context.Context, userID string, id string) (*models.BodyMeasurementRecord, error)
}

func (m *mockMeasurementRepo) Create(ctx context.Context, checkInID string, measurementType string, side *string, value float64) (*models.BodyMeasurementRecord, error) {
	return m.createFn(ctx, checkInID, measurementType, side, value)
}

func (m *mockMeasurementRepo) GetByID(ctx context.Context, userID string, id string) (*models.BodyMeasurementRecord, error) {
	return m.getByIDFn(ctx, userID, id)
}

func (m *mockMeasurementRepo) ListByCheckIn(ctx context.Context, userID string, checkInID string) ([]models.BodyMeasurementRecord, error) {
	return m.listByCheckInFn(ctx, userID, checkInID)
}

func (m *mockMeasurementRepo) Update(ctx context.Context, userID string, id string, measurementType string, side *string, value float64) (*models.BodyMeasurementRecord, error) {
	return m.updateFn(ctx, userID, id, measurementType, side, value)
}

func (m *mockMeasurementRepo) Delete(ctx context.Context, userID string, id string) (*models.BodyMeasurementRecord, error) {
	return m.deleteFn(ctx, userID, id)
}

type mockPhotoRepo struct {
	atlasPostgres.ProgressPhotoRepository
	createFn         func(ctx context.Context, checkInID string, filePath string, originalFileName string, mimeType string, sizeBytes int64, angle *string, label *string, notes *string) (*models.ProgressPhotoRecord, error)
	getByIDFn        func(ctx context.Context, userID string, id string) (*models.ProgressPhotoRecord, error)
	listByCheckInFn  func(ctx context.Context, userID string, checkInID string) ([]models.ProgressPhotoRecord, error)
	deleteFn         func(ctx context.Context, userID string, id string) (*models.ProgressPhotoRecord, error)
	countByCheckInFn func(ctx context.Context, checkInID string) (int64, error)
}

func (m *mockPhotoRepo) Create(ctx context.Context, checkInID string, filePath string, originalFileName string, mimeType string, sizeBytes int64, angle *string, label *string, notes *string) (*models.ProgressPhotoRecord, error) {
	return m.createFn(ctx, checkInID, filePath, originalFileName, mimeType, sizeBytes, angle, label, notes)
}

func (m *mockPhotoRepo) GetByID(ctx context.Context, userID string, id string) (*models.ProgressPhotoRecord, error) {
	return m.getByIDFn(ctx, userID, id)
}

func (m *mockPhotoRepo) ListByCheckIn(ctx context.Context, userID string, checkInID string) ([]models.ProgressPhotoRecord, error) {
	return m.listByCheckInFn(ctx, userID, checkInID)
}

func (m *mockPhotoRepo) Delete(ctx context.Context, userID string, id string) (*models.ProgressPhotoRecord, error) {
	return m.deleteFn(ctx, userID, id)
}

func (m *mockPhotoRepo) CountByCheckIn(ctx context.Context, checkInID string) (int64, error) {
	return m.countByCheckInFn(ctx, checkInID)
}

var (
	bciTestDate = models.MustDate("2026-06-20")
)

func bciTestRecord() *models.BodyCheckInRecord {
	return &models.BodyCheckInRecord{
		ID:                testID,
		UserID:            testUserID,
		Date:              bciTestDate,
		Weight:            ptrFloat64(75.5),
		BodyFatPercentage: ptrFloat64(15.0),
		Notes:             ptrStr("Morning check-in"),
		CreatedAt:         "2026-06-20T00:00:00Z",
		UpdatedAt:         "2026-06-20T00:00:00Z",
	}
}

func bciTestMeasurementRecord() *models.BodyMeasurementRecord {
	return &models.BodyMeasurementRecord{
		ID:              "measurement-1",
		CheckInID:       testID,
		MeasurementType: "CHEST",
		Side:            nil,
		Value:           100.0,
		CreatedAt:       "2026-06-20T00:00:00Z",
		UpdatedAt:       "2026-06-20T00:00:00Z",
	}
}

func bciTestPhotoRecord() *models.ProgressPhotoRecord {
	return &models.ProgressPhotoRecord{
		ID:               "photo-1",
		CheckInID:        testID,
		FilePath:         "/uploads/photo.jpg",
		OriginalFileName: "photo.jpg",
		MimeType:         "image/jpeg",
		SizeBytes:        204800,
		Angle:            ptrStr("FRONT"),
		Label:            nil,
		Notes:            nil,
		CreatedAt:        "2026-06-20T00:00:00Z",
		UpdatedAt:        "2026-06-20T00:00:00Z",
	}
}

// ----- BodyCheckInService Create -----

func TestBodyCheckInService_Create_WeightZero(t *testing.T) {
	svc := service.NewBodyCheckInService(&mockCheckInRepo{}, &mockMeasurementRepo{}, &mockPhotoRepo{})

	entry, err := svc.Create(ctx, testUserID, models.CreateCheckInInput{
		Date:   bciTestDate,
		Weight: ptrFloat64(0),
	})
	assert.ErrorIs(t, err, service.ErrCheckInWeightInvalid)
	assert.Nil(t, entry)
}

func TestBodyCheckInService_Create_BodyFatInvalid(t *testing.T) {
	svc := service.NewBodyCheckInService(&mockCheckInRepo{}, &mockMeasurementRepo{}, &mockPhotoRepo{})

	t.Run("zero body fat", func(t *testing.T) {
		entry, err := svc.Create(ctx, testUserID, models.CreateCheckInInput{
			Date:              bciTestDate,
			BodyFatPercentage: ptrFloat64(0),
		})
		assert.ErrorIs(t, err, service.ErrCheckInBodyFatInvalid)
		assert.Nil(t, entry)
	})

	t.Run("over 100 body fat", func(t *testing.T) {
		entry, err := svc.Create(ctx, testUserID, models.CreateCheckInInput{
			Date:              bciTestDate,
			BodyFatPercentage: ptrFloat64(101),
		})
		assert.ErrorIs(t, err, service.ErrCheckInBodyFatInvalid)
		assert.Nil(t, entry)
	})
}

func TestBodyCheckInService_Create_Success(t *testing.T) {
	svc := service.NewBodyCheckInService(&mockCheckInRepo{
		createFn: func(ctx context.Context, userID string, date models.Date, weight *float64, bodyFatPercentage *float64, notes *string) (*models.BodyCheckInRecord, error) {
			return bciTestRecord(), nil
		},
	}, &mockMeasurementRepo{}, &mockPhotoRepo{})

	entry, err := svc.Create(ctx, testUserID, models.CreateCheckInInput{
		Date:              bciTestDate,
		Weight:            ptrFloat64(75.5),
		BodyFatPercentage: ptrFloat64(15.0),
		Notes:             ptrStr("Morning check-in"),
	})
	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.Equal(t, 75.5, *entry.Weight)
	assert.Equal(t, 15.0, *entry.BodyFatPercentage)
}

// ----- BodyCheckInService GetByID -----

func TestBodyCheckInService_GetByID_Success(t *testing.T) {
	svc := service.NewBodyCheckInService(&mockCheckInRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.BodyCheckInRecord, error) {
			return bciTestRecord(), nil
		},
	}, &mockMeasurementRepo{
		listByCheckInFn: func(ctx context.Context, userID string, checkInID string) ([]models.BodyMeasurementRecord, error) {
			return []models.BodyMeasurementRecord{*bciTestMeasurementRecord()}, nil
		},
	}, &mockPhotoRepo{
		listByCheckInFn: func(ctx context.Context, userID string, checkInID string) ([]models.ProgressPhotoRecord, error) {
			return []models.ProgressPhotoRecord{*bciTestPhotoRecord()}, nil
		},
	})

	entry, err := svc.GetByID(ctx, testUserID, testID)
	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.Equal(t, 1, len(entry.Measurements))
	assert.Equal(t, 1, len(entry.ProgressPhotos))
	assert.Equal(t, models.MeasurementType("CHEST"), entry.Measurements[0].MeasurementType)
	assert.Equal(t, "FRONT", string(*entry.ProgressPhotos[0].Angle))
}

func TestBodyCheckInService_GetByID_NotFound(t *testing.T) {
	svc := service.NewBodyCheckInService(&mockCheckInRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.BodyCheckInRecord, error) {
			return nil, nil
		},
	}, &mockMeasurementRepo{}, &mockPhotoRepo{})

	entry, err := svc.GetByID(ctx, testUserID, testID)
	assert.ErrorIs(t, err, service.ErrCheckInNotFound)
	assert.Nil(t, entry)
}

// ----- BodyCheckInService Update -----

func TestBodyCheckInService_Update_BodyFatInvalid(t *testing.T) {
	svc := service.NewBodyCheckInService(&mockCheckInRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.BodyCheckInRecord, error) {
			return bciTestRecord(), nil
		},
	}, &mockMeasurementRepo{}, &mockPhotoRepo{})

	entry, err := svc.Update(ctx, testUserID, testID, models.UpdateCheckInInput{
		BodyFatPercentage: ptrFloat64(101),
	})
	assert.ErrorIs(t, err, service.ErrCheckInBodyFatInvalid)
	assert.Nil(t, entry)
}

func TestBodyCheckInService_Update_Success(t *testing.T) {
	svc := service.NewBodyCheckInService(&mockCheckInRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.BodyCheckInRecord, error) {
			return bciTestRecord(), nil
		},
		updateFn: func(ctx context.Context, userID string, id string, weight *float64, bodyFatPercentage *float64, notes *string) (*models.BodyCheckInRecord, error) {
			assert.Equal(t, 80.0, *weight)
			return &models.BodyCheckInRecord{
				ID:                id,
				UserID:            userID,
				Date:              bciTestDate,
				Weight:            weight,
				BodyFatPercentage: bodyFatPercentage,
				Notes:             notes,
				CreatedAt:         "2026-06-20T00:00:00Z",
				UpdatedAt:         "2026-06-20T12:00:00Z",
			}, nil
		},
	}, &mockMeasurementRepo{}, &mockPhotoRepo{})

	entry, err := svc.Update(ctx, testUserID, testID, models.UpdateCheckInInput{
		Weight: ptrFloat64(80.0),
	})
	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.Equal(t, 80.0, *entry.Weight)
}

// ----- BodyCheckInService Delete -----

func TestBodyCheckInService_Delete_Success(t *testing.T) {
	svc := service.NewBodyCheckInService(&mockCheckInRepo{
		deleteFn: func(ctx context.Context, userID string, id string) (*models.BodyCheckInRecord, error) {
			return bciTestRecord(), nil
		},
	}, &mockMeasurementRepo{}, &mockPhotoRepo{})

	entry, err := svc.Delete(ctx, testUserID, testID)
	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.Equal(t, testID, entry.ID)
}

// ----- BodyMeasurementService Create -----

func TestBodyMeasurementService_Create_InvalidType(t *testing.T) {
	svc := service.NewBodyMeasurementService(&mockMeasurementRepo{}, &mockCheckInRepo{})

	entry, err := svc.Create(ctx, testUserID, models.CreateMeasurementInput{
		MeasurementType: "INVALID_TYPE",
		Value:           100,
	})
	assert.ErrorIs(t, err, service.ErrMeasurementTypeInvalid)
	assert.Nil(t, entry)
}

func TestBodyMeasurementService_Create_ValueZero(t *testing.T) {
	svc := service.NewBodyMeasurementService(&mockMeasurementRepo{}, &mockCheckInRepo{})

	entry, err := svc.Create(ctx, testUserID, models.CreateMeasurementInput{
		MeasurementType: models.MeasurementTypeChest,
		Value:           0,
	})
	assert.ErrorIs(t, err, service.ErrMeasurementValueInvalid)
	assert.Nil(t, entry)
}

func TestBodyMeasurementService_Create_SideOnUnpairedType(t *testing.T) {
	svc := service.NewBodyMeasurementService(&mockMeasurementRepo{}, &mockCheckInRepo{})

	entry, err := svc.Create(ctx, testUserID, models.CreateMeasurementInput{
		MeasurementType: models.MeasurementTypeChest,
		Value:           100,
		Side:            sideModelPtr(models.MeasurementSideLeft),
	})
	assert.ErrorIs(t, err, service.ErrMeasurementSideInvalid)
	assert.Nil(t, entry)
}

func TestBodyMeasurementService_Create_Success(t *testing.T) {
	svc := service.NewBodyMeasurementService(&mockMeasurementRepo{
		createFn: func(ctx context.Context, checkInID string, measurementType string, side *string, value float64) (*models.BodyMeasurementRecord, error) {
			assert.Equal(t, "CHEST", measurementType)
			assert.Nil(t, side)
			return bciTestMeasurementRecord(), nil
		},
	}, &mockCheckInRepo{})

	entry, err := svc.Create(ctx, testUserID, models.CreateMeasurementInput{
		MeasurementType: models.MeasurementTypeChest,
		Value:           100,
	})
	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.Equal(t, models.MeasurementType("CHEST"), entry.MeasurementType)
	assert.Equal(t, 100.0, entry.Value)
}

// ----- BodyMeasurementService Update -----

func TestBodyMeasurementService_Update_Success(t *testing.T) {
	svc := service.NewBodyMeasurementService(&mockMeasurementRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.BodyMeasurementRecord, error) {
			return bciTestMeasurementRecord(), nil
		},
		updateFn: func(ctx context.Context, userID string, id string, measurementType string, side *string, value float64) (*models.BodyMeasurementRecord, error) {
			assert.Equal(t, 105.0, value)
			return &models.BodyMeasurementRecord{
				ID:              id,
				CheckInID:       testID,
				MeasurementType: measurementType,
				Side:            side,
				Value:           value,
				CreatedAt:       "2026-06-20T00:00:00Z",
				UpdatedAt:       "2026-06-20T12:00:00Z",
			}, nil
		},
	}, &mockCheckInRepo{})

	entry, err := svc.Update(ctx, testUserID, "measurement-1", models.UpdateMeasurementInput{
		Value: ptrFloat64(105),
	})
	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.Equal(t, 105.0, entry.Value)
}

// ----- BodyMeasurementService Delete -----

func TestBodyMeasurementService_Delete_Success(t *testing.T) {
	svc := service.NewBodyMeasurementService(&mockMeasurementRepo{
		deleteFn: func(ctx context.Context, userID string, id string) (*models.BodyMeasurementRecord, error) {
			return bciTestMeasurementRecord(), nil
		},
	}, &mockCheckInRepo{})

	entry, err := svc.Delete(ctx, testUserID, "measurement-1")
	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.Equal(t, "measurement-1", entry.ID)
}

// ----- Helpers -----

func sideModelPtr(s models.MeasurementSide) *models.MeasurementSide {
	return &s
}