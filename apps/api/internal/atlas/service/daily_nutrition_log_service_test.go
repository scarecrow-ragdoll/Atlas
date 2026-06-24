// FILE: apps/api/internal/atlas/service/daily_nutrition_log_service_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for DailyNutritionLogService factual food-log behavior.
//   SCOPE: Get-or-create daily logs, entry create/update/delete validation, product snapshot totals, inactive/missing product rejection, and user-scoped delete semantics.
//   DEPENDS: apps/api/internal/atlas/service, apps/api/internal/atlas/repository/postgres (mock), apps/api/internal/atlas/models.
//   LINKS: M-API-NUTRITION / V-M-API-NUTRITION.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Task 2 RED tests for factual daily nutrition logs.
// END_CHANGE_SUMMARY

package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasPostgres "monorepo-template/apps/api/internal/atlas/repository/postgres"
	"monorepo-template/apps/api/internal/atlas/service"
)

type mockDailyNutritionLogRepo struct {
	atlasPostgres.DailyNutritionLogRepository
	getOrCreateFn func(ctx context.Context, userID string, date models.Date, notes *string) (*models.DailyNutritionLogRecord, error)
	getByDateFn   func(ctx context.Context, userID string, date models.Date) (*models.DailyNutritionLogRecord, error)
	listByRangeFn func(ctx context.Context, userID string, from, to models.Date) ([]models.DailyNutritionLogRecord, error)
	updateNotesFn func(ctx context.Context, userID string, id string, notes *string) (*models.DailyNutritionLogRecord, error)
	addEntryFn    func(ctx context.Context, userID string, input models.CreateDailyNutritionEntryRecordInput) (*models.DailyNutritionEntryRecord, error)
	listEntriesFn func(ctx context.Context, userID string, dailyLogID string) ([]models.DailyNutritionEntryRecord, error)
	updateEntryFn func(ctx context.Context, userID string, id string, input models.UpdateDailyNutritionEntryInput) (*models.DailyNutritionEntryRecord, error)
	deleteEntryFn func(ctx context.Context, userID string, id string) (*models.DailyNutritionEntryRecord, error)
}

func (m *mockDailyNutritionLogRepo) GetOrCreate(ctx context.Context, userID string, date models.Date, notes *string) (*models.DailyNutritionLogRecord, error) {
	return m.getOrCreateFn(ctx, userID, date, notes)
}
func (m *mockDailyNutritionLogRepo) GetByDate(ctx context.Context, userID string, date models.Date) (*models.DailyNutritionLogRecord, error) {
	return m.getByDateFn(ctx, userID, date)
}
func (m *mockDailyNutritionLogRepo) ListByRange(ctx context.Context, userID string, from, to models.Date) ([]models.DailyNutritionLogRecord, error) {
	return m.listByRangeFn(ctx, userID, from, to)
}
func (m *mockDailyNutritionLogRepo) UpdateNotes(ctx context.Context, userID string, id string, notes *string) (*models.DailyNutritionLogRecord, error) {
	return m.updateNotesFn(ctx, userID, id, notes)
}
func (m *mockDailyNutritionLogRepo) AddEntry(ctx context.Context, userID string, input models.CreateDailyNutritionEntryRecordInput) (*models.DailyNutritionEntryRecord, error) {
	return m.addEntryFn(ctx, userID, input)
}
func (m *mockDailyNutritionLogRepo) ListEntries(ctx context.Context, userID string, dailyLogID string) ([]models.DailyNutritionEntryRecord, error) {
	return m.listEntriesFn(ctx, userID, dailyLogID)
}
func (m *mockDailyNutritionLogRepo) UpdateEntry(ctx context.Context, userID string, id string, input models.UpdateDailyNutritionEntryInput) (*models.DailyNutritionEntryRecord, error) {
	return m.updateEntryFn(ctx, userID, id, input)
}
func (m *mockDailyNutritionLogRepo) DeleteEntry(ctx context.Context, userID string, id string) (*models.DailyNutritionEntryRecord, error) {
	return m.deleteEntryFn(ctx, userID, id)
}

type mockDailyNutritionProductService struct {
	service.NutritionProductService
	getByIDFn func(ctx context.Context, userID string, id string) (*models.NutritionProduct, error)
}

func (m *mockDailyNutritionProductService) GetByID(ctx context.Context, userID string, id string) (*models.NutritionProduct, error) {
	return m.getByIDFn(ctx, userID, id)
}

var (
	dailyNutritionDate      = models.MustDate("2026-06-24")
	dailyNutritionLogID     = "111e8400-e29b-41d4-a716-446655440000"
	dailyNutritionEntryID   = "222e8400-e29b-41d4-a716-446655440000"
	dailyNutritionProductID = "333e8400-e29b-41d4-a716-446655440000"
)

func newDailyNutritionService(repo *mockDailyNutritionLogRepo, productSvc *mockDailyNutritionProductService) service.DailyNutritionLogService {
	return service.NewDailyNutritionLogService(repo, productSvc, zap.NewNop())
}

func dailyNutritionLogRecord() *models.DailyNutritionLogRecord {
	return &models.DailyNutritionLogRecord{
		ID:        dailyNutritionLogID,
		UserID:    testUserID,
		Date:      dailyNutritionDate,
		Notes:     ptrStr("training day"),
		CreatedAt: "2026-06-24T00:00:00Z",
		UpdatedAt: "2026-06-24T00:00:00Z",
	}
}

func activeDailyNutritionProduct() *models.NutritionProduct {
	return &models.NutritionProduct{
		ID:              dailyNutritionProductID,
		UserID:          testUserID,
		Name:            "Chicken Breast",
		CaloriesPer100g: 165,
		ProteinPer100g:  31,
		FatPer100g:      3.6,
		CarbsPer100g:    0,
		IsActive:        true,
		CreatedAt:       "2026-06-24T00:00:00Z",
		UpdatedAt:       "2026-06-24T00:00:00Z",
	}
}

func dailyNutritionEntryRecord(amount float64) *models.DailyNutritionEntryRecord {
	return &models.DailyNutritionEntryRecord{
		ID:                      dailyNutritionEntryID,
		DailyLogID:              dailyNutritionLogID,
		ProductID:               dailyNutritionProductID,
		ProductNameSnapshot:     "Chicken Breast",
		CaloriesPer100gSnapshot: 165,
		ProteinPer100gSnapshot:  31,
		FatPer100gSnapshot:      3.6,
		CarbsPer100gSnapshot:    0,
		AmountGrams:             amount,
		MealLabel:               ptrStr("Lunch"),
		Notes:                   ptrStr("grilled"),
		Position:                1,
		CreatedAt:               "2026-06-24T01:00:00Z",
		UpdatedAt:               "2026-06-24T01:00:00Z",
	}
}

func TestDailyNutritionLogService_GetByDateCreatesEmptyLogWithZeroTotals(t *testing.T) {
	svc := newDailyNutritionService(&mockDailyNutritionLogRepo{
		getOrCreateFn: func(ctx context.Context, userID string, date models.Date, notes *string) (*models.DailyNutritionLogRecord, error) {
			assert.Equal(t, testUserID, userID)
			assert.Equal(t, "2026-06-24", date.String())
			return dailyNutritionLogRecord(), nil
		},
		listEntriesFn: func(ctx context.Context, userID string, dailyLogID string) ([]models.DailyNutritionEntryRecord, error) {
			assert.Equal(t, testUserID, userID)
			assert.Equal(t, dailyNutritionLogID, dailyLogID)
			return []models.DailyNutritionEntryRecord{}, nil
		},
	}, &mockDailyNutritionProductService{})

	log, err := svc.GetByDate(ctx, testUserID, dailyNutritionDate)
	require.NoError(t, err)
	require.NotNil(t, log)
	assert.Empty(t, log.Entries)
	assert.Equal(t, models.NutritionMacros{}, log.Totals)
}

func TestDailyNutritionLogService_AddEntrySnapshotsProductAndCalculatesTotals(t *testing.T) {
	svc := newDailyNutritionService(&mockDailyNutritionLogRepo{
		getOrCreateFn: func(ctx context.Context, userID string, date models.Date, notes *string) (*models.DailyNutritionLogRecord, error) {
			return dailyNutritionLogRecord(), nil
		},
		addEntryFn: func(ctx context.Context, userID string, input models.CreateDailyNutritionEntryRecordInput) (*models.DailyNutritionEntryRecord, error) {
			assert.Equal(t, dailyNutritionLogID, input.DailyLogID)
			assert.Equal(t, dailyNutritionProductID, input.ProductID)
			assert.Equal(t, 50.0, input.AmountGrams)
			assert.Equal(t, int32(1), input.Position)
			return dailyNutritionEntryRecord(50), nil
		},
		listEntriesFn: func(ctx context.Context, userID string, dailyLogID string) ([]models.DailyNutritionEntryRecord, error) {
			assert.Equal(t, testUserID, userID)
			assert.Equal(t, dailyNutritionLogID, dailyLogID)
			return []models.DailyNutritionEntryRecord{*dailyNutritionEntryRecord(50)}, nil
		},
	}, &mockDailyNutritionProductService{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.NutritionProduct, error) {
			assert.Equal(t, testUserID, userID)
			assert.Equal(t, dailyNutritionProductID, id)
			return activeDailyNutritionProduct(), nil
		},
	})

	log, err := svc.AddEntry(ctx, testUserID, models.AddDailyNutritionEntryInput{
		Date:        dailyNutritionDate,
		ProductID:   dailyNutritionProductID,
		AmountGrams: 50,
		MealLabel:   ptrStr("Lunch"),
		Notes:       ptrStr("grilled"),
		Position:    1,
	})
	require.NoError(t, err)
	require.NotNil(t, log)
	require.Len(t, log.Entries, 1)
	assert.Equal(t, "Chicken Breast", log.Entries[0].ProductNameSnapshot)
	assert.Equal(t, 82.5, log.Totals.Calories)
	assert.Equal(t, 15.5, log.Totals.Protein)
	assert.Equal(t, 1.8, log.Totals.Fat)
	assert.Equal(t, 0.0, log.Totals.Carbs)
}

func TestDailyNutritionLogService_RejectsMissingProduct(t *testing.T) {
	svc := newDailyNutritionService(&mockDailyNutritionLogRepo{}, &mockDailyNutritionProductService{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.NutritionProduct, error) {
			return nil, service.ErrProductNotFound
		},
	})

	log, err := svc.AddEntry(ctx, testUserID, models.AddDailyNutritionEntryInput{
		Date:        dailyNutritionDate,
		ProductID:   dailyNutritionProductID,
		AmountGrams: 50,
	})
	assert.ErrorIs(t, err, service.ErrDailyNutritionProductNotFound)
	assert.Nil(t, log)
}

func TestDailyNutritionLogService_RejectsNonPositiveAmount(t *testing.T) {
	svc := newDailyNutritionService(&mockDailyNutritionLogRepo{}, &mockDailyNutritionProductService{})

	log, err := svc.AddEntry(ctx, testUserID, models.AddDailyNutritionEntryInput{
		Date:        dailyNutritionDate,
		ProductID:   dailyNutritionProductID,
		AmountGrams: 0,
	})
	assert.ErrorIs(t, err, service.ErrDailyNutritionAmountInvalid)
	assert.Nil(t, log)
}

func TestDailyNutritionLogService_RejectsInactiveProductForNewEntry(t *testing.T) {
	inactive := activeDailyNutritionProduct()
	inactive.IsActive = false
	svc := newDailyNutritionService(&mockDailyNutritionLogRepo{}, &mockDailyNutritionProductService{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.NutritionProduct, error) {
			return inactive, nil
		},
	})

	log, err := svc.AddEntry(ctx, testUserID, models.AddDailyNutritionEntryInput{
		Date:        dailyNutritionDate,
		ProductID:   dailyNutritionProductID,
		AmountGrams: 50,
	})
	assert.ErrorIs(t, err, service.ErrDailyNutritionProductInactive)
	assert.Nil(t, log)
}

func TestDailyNutritionLogService_ProductEditDoesNotChangeExistingEntryTotals(t *testing.T) {
	svc := newDailyNutritionService(&mockDailyNutritionLogRepo{
		getOrCreateFn: func(ctx context.Context, userID string, date models.Date, notes *string) (*models.DailyNutritionLogRecord, error) {
			return dailyNutritionLogRecord(), nil
		},
		listEntriesFn: func(ctx context.Context, userID string, dailyLogID string) ([]models.DailyNutritionEntryRecord, error) {
			return []models.DailyNutritionEntryRecord{*dailyNutritionEntryRecord(100)}, nil
		},
	}, &mockDailyNutritionProductService{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.NutritionProduct, error) {
			t.Fatal("GetByDate must not recalculate totals from current product values")
			return nil, nil
		},
	})

	log, err := svc.GetByDate(ctx, testUserID, dailyNutritionDate)
	require.NoError(t, err)
	require.NotNil(t, log)
	assert.Equal(t, 165.0, log.Totals.Calories)
	assert.Equal(t, 31.0, log.Totals.Protein)
	assert.Equal(t, 3.6, log.Totals.Fat)
	assert.Equal(t, 0.0, log.Totals.Carbs)
}

func TestDailyNutritionLogService_DeleteEntryRequiresOwnedParentLog(t *testing.T) {
	svc := newDailyNutritionService(&mockDailyNutritionLogRepo{
		deleteEntryFn: func(ctx context.Context, userID string, id string) (*models.DailyNutritionEntryRecord, error) {
			assert.Equal(t, testUserID, userID)
			assert.Equal(t, dailyNutritionEntryID, id)
			return nil, nil
		},
	}, &mockDailyNutritionProductService{})

	log, err := svc.DeleteEntry(ctx, testUserID, dailyNutritionEntryID)
	assert.ErrorIs(t, err, service.ErrDailyNutritionEntryNotFound)
	assert.Nil(t, log)
}

func TestDailyNutritionLogService_UpdateEntryRejectsNonPositiveAmount(t *testing.T) {
	svc := newDailyNutritionService(&mockDailyNutritionLogRepo{}, &mockDailyNutritionProductService{})

	log, err := svc.UpdateEntry(ctx, testUserID, dailyNutritionEntryID, models.UpdateDailyNutritionEntryInput{
		DailyLogID:  dailyNutritionLogID,
		AmountGrams: ptrFloat64(0),
	})
	assert.ErrorIs(t, err, service.ErrDailyNutritionAmountInvalid)
	assert.Nil(t, log)
}
