// FILE: apps/api/internal/atlas/service/nutrition_override_service_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for DailyNutritionOverrideService and override item CRUD covering Create (upsert), GetByID, GetByDate, ListByRange, Update, Delete, CreateItem, UpdateItem, DeleteItem.
//   SCOPE: Success paths, validation errors, not-found, override isolation pattern, operation enum validation.
//   DEPENDS: apps/api/internal/atlas/service, apps/api/internal/atlas/repository/postgres (mock), apps/api/internal/atlas/models.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT

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

type mockNutritionOverrideRepo struct {
	atlasPostgres.DailyNutritionOverrideRepository
	upsertFn      func(ctx context.Context, userID string, date string, notes *string) (*models.DailyNutritionOverrideRecord, error)
	getByIDFn     func(ctx context.Context, userID string, id string) (*models.DailyNutritionOverrideRecord, error)
	getByDateFn   func(ctx context.Context, userID string, date string) (*models.DailyNutritionOverrideRecord, error)
	listByRangeFn func(ctx context.Context, userID string, startDate, endDate string) ([]models.DailyNutritionOverrideRecord, error)
	updateFn      func(ctx context.Context, userID string, id string, notes *string) (*models.DailyNutritionOverrideRecord, error)
	deleteFn      func(ctx context.Context, userID string, id string) (*models.DailyNutritionOverrideRecord, error)
}

func (m *mockNutritionOverrideRepo) Upsert(ctx context.Context, userID string, date string, notes *string) (*models.DailyNutritionOverrideRecord, error) {
	return m.upsertFn(ctx, userID, date, notes)
}
func (m *mockNutritionOverrideRepo) GetByID(ctx context.Context, userID string, id string) (*models.DailyNutritionOverrideRecord, error) {
	return m.getByIDFn(ctx, userID, id)
}
func (m *mockNutritionOverrideRepo) GetByDate(ctx context.Context, userID string, date string) (*models.DailyNutritionOverrideRecord, error) {
	return m.getByDateFn(ctx, userID, date)
}
func (m *mockNutritionOverrideRepo) ListByRange(ctx context.Context, userID string, startDate, endDate string) ([]models.DailyNutritionOverrideRecord, error) {
	return m.listByRangeFn(ctx, userID, startDate, endDate)
}
func (m *mockNutritionOverrideRepo) Update(ctx context.Context, userID string, id string, notes *string) (*models.DailyNutritionOverrideRecord, error) {
	return m.updateFn(ctx, userID, id, notes)
}
func (m *mockNutritionOverrideRepo) Delete(ctx context.Context, userID string, id string) (*models.DailyNutritionOverrideRecord, error) {
	return m.deleteFn(ctx, userID, id)
}

type mockNutritionOverrideItemRepo struct {
	atlasPostgres.DailyNutritionOverrideItemRepository
	createFn         func(ctx context.Context, overrideID string, productID string, amountGrams float64, operation string, mealLabel, notes *string) (*models.DailyNutritionOverrideItemRecord, error)
	getByIDFn        func(ctx context.Context, id string) (*models.DailyNutritionOverrideItemRecord, error)
	listByOverrideFn func(ctx context.Context, overrideID string) ([]models.DailyNutritionOverrideItemRecord, error)
	updateFn         func(ctx context.Context, id string, amountGrams float64, operation string, mealLabel, notes *string) (*models.DailyNutritionOverrideItemRecord, error)
	deleteFn         func(ctx context.Context, id string) (*models.DailyNutritionOverrideItemRecord, error)
}

func (m *mockNutritionOverrideItemRepo) Create(ctx context.Context, overrideID string, productID string, amountGrams float64, operation string, mealLabel, notes *string) (*models.DailyNutritionOverrideItemRecord, error) {
	return m.createFn(ctx, overrideID, productID, amountGrams, operation, mealLabel, notes)
}
func (m *mockNutritionOverrideItemRepo) GetByID(ctx context.Context, id string) (*models.DailyNutritionOverrideItemRecord, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockNutritionOverrideItemRepo) ListByOverride(ctx context.Context, overrideID string) ([]models.DailyNutritionOverrideItemRecord, error) {
	return m.listByOverrideFn(ctx, overrideID)
}
func (m *mockNutritionOverrideItemRepo) Update(ctx context.Context, id string, amountGrams float64, operation string, mealLabel, notes *string) (*models.DailyNutritionOverrideItemRecord, error) {
	return m.updateFn(ctx, id, amountGrams, operation, mealLabel, notes)
}
func (m *mockNutritionOverrideItemRepo) Delete(ctx context.Context, id string) (*models.DailyNutritionOverrideItemRecord, error) {
	return m.deleteFn(ctx, id)
}

var overrideTestRecord = &models.DailyNutritionOverrideRecord{
	ID: testID, UserID: testUserID,
	Date:      models.MustDate("2026-06-16"),
	Notes:     ptrStr("Cheat day"),
	CreatedAt: "2026-06-20T00:00:00Z", UpdatedAt: "2026-06-20T00:00:00Z",
}

var overrideItemTestRecord = &models.DailyNutritionOverrideItemRecord{
	ID: "770e8400-e29b-41d4-a716-446655440000", OverrideID: testID,
	ProductID: "880e8400-e29b-41d4-a716-446655440000",
	AmountGrams: 100, Operation: "add", MealLabel: ptrStr("Snack"),
	CreatedAt: "2026-06-20T00:00:00Z", UpdatedAt: "2026-06-20T00:00:00Z",
}

func newOverrideService(overrideRepo *mockNutritionOverrideRepo, itemRepo *mockNutritionOverrideItemRepo) service.DailyNutritionOverrideService {
	return service.NewNutritionOverrideService(overrideRepo, itemRepo, zap.NewNop())
}

// ----- Override CRUD -----

func TestOverrideService_Create_Success(t *testing.T) {
	svc := newOverrideService(&mockNutritionOverrideRepo{
		upsertFn: func(ctx context.Context, userID string, date string, notes *string) (*models.DailyNutritionOverrideRecord, error) {
			return overrideTestRecord, nil
		},
	}, &mockNutritionOverrideItemRepo{})

	override, err := svc.Create(ctx, testUserID, models.CreateOverrideInput{
		Date:  models.MustDate("2026-06-16"),
		Notes: ptrStr("Cheat day"),
	})
	require.NoError(t, err)
	require.NotNil(t, override)
	assert.Equal(t, "2026-06-16", override.Date)
}

func TestOverrideService_Create_DateEmpty(t *testing.T) {
	svc := newOverrideService(&mockNutritionOverrideRepo{}, &mockNutritionOverrideItemRepo{})
	override, err := svc.Create(ctx, testUserID, models.CreateOverrideInput{Date: models.Date{}})
	assert.ErrorIs(t, err, service.ErrOverrideDateRequired)
	assert.Nil(t, override)
}

func TestOverrideService_GetByID_Success(t *testing.T) {
	svc := newOverrideService(&mockNutritionOverrideRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.DailyNutritionOverrideRecord, error) {
			return overrideTestRecord, nil
		},
	}, &mockNutritionOverrideItemRepo{
		listByOverrideFn: func(ctx context.Context, overrideID string) ([]models.DailyNutritionOverrideItemRecord, error) {
			return []models.DailyNutritionOverrideItemRecord{*overrideItemTestRecord}, nil
		},
	})

	override, err := svc.GetByID(ctx, testUserID, testID)
	require.NoError(t, err)
	require.NotNil(t, override)
	assert.Len(t, override.Items, 1)
}

func TestOverrideService_GetByID_NotFound(t *testing.T) {
	svc := newOverrideService(&mockNutritionOverrideRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.DailyNutritionOverrideRecord, error) {
			return nil, nil
		},
	}, &mockNutritionOverrideItemRepo{})

	override, err := svc.GetByID(ctx, testUserID, testID)
	assert.ErrorIs(t, err, service.ErrOverrideNotFound)
	assert.Nil(t, override)
}

func TestOverrideService_GetByDate_Success(t *testing.T) {
	svc := newOverrideService(&mockNutritionOverrideRepo{
		getByDateFn: func(ctx context.Context, userID string, date string) (*models.DailyNutritionOverrideRecord, error) {
			return overrideTestRecord, nil
		},
	}, &mockNutritionOverrideItemRepo{
		listByOverrideFn: func(ctx context.Context, overrideID string) ([]models.DailyNutritionOverrideItemRecord, error) {
			return []models.DailyNutritionOverrideItemRecord{}, nil
		},
	})

	override, err := svc.GetByDate(ctx, testUserID, "2026-06-16")
	require.NoError(t, err)
	require.NotNil(t, override)
}

func TestOverrideService_GetByDate_Nil(t *testing.T) {
	svc := newOverrideService(&mockNutritionOverrideRepo{
		getByDateFn: func(ctx context.Context, userID string, date string) (*models.DailyNutritionOverrideRecord, error) {
			return nil, nil
		},
	}, &mockNutritionOverrideItemRepo{})

	override, err := svc.GetByDate(ctx, testUserID, "2026-06-16")
	require.NoError(t, err)
	assert.Nil(t, override)
}

func TestOverrideService_ListByRange_Success(t *testing.T) {
	svc := newOverrideService(&mockNutritionOverrideRepo{
		listByRangeFn: func(ctx context.Context, userID string, startDate, endDate string) ([]models.DailyNutritionOverrideRecord, error) {
			return []models.DailyNutritionOverrideRecord{*overrideTestRecord}, nil
		},
	}, &mockNutritionOverrideItemRepo{
		listByOverrideFn: func(ctx context.Context, overrideID string) ([]models.DailyNutritionOverrideItemRecord, error) {
			return []models.DailyNutritionOverrideItemRecord{}, nil
		},
	})

	overrides, err := svc.ListByRange(ctx, testUserID, "2026-06-01", "2026-06-30")
	require.NoError(t, err)
	assert.Len(t, overrides, 1)
}

func TestOverrideService_Update_Success(t *testing.T) {
	svc := newOverrideService(&mockNutritionOverrideRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.DailyNutritionOverrideRecord, error) {
			return overrideTestRecord, nil
		},
		updateFn: func(ctx context.Context, userID string, id string, notes *string) (*models.DailyNutritionOverrideRecord, error) {
			return &models.DailyNutritionOverrideRecord{
				ID: id, UserID: userID, Date: models.MustDate("2026-06-16"),
				Notes: notes, CreatedAt: "2026-06-20T00:00:00Z", UpdatedAt: "2026-06-20T12:00:00Z",
			}, nil
		},
	}, &mockNutritionOverrideItemRepo{
		listByOverrideFn: func(ctx context.Context, overrideID string) ([]models.DailyNutritionOverrideItemRecord, error) {
			return []models.DailyNutritionOverrideItemRecord{}, nil
		},
	})

	override, err := svc.Update(ctx, testUserID, testID, models.UpdateOverrideInput{
		Notes: ptrStr("Updated notes"),
	})
	require.NoError(t, err)
	require.NotNil(t, override)
	assert.Equal(t, "Updated notes", *override.Notes)
}

func TestOverrideService_Delete_Success(t *testing.T) {
	svc := newOverrideService(&mockNutritionOverrideRepo{
		deleteFn: func(ctx context.Context, userID string, id string) (*models.DailyNutritionOverrideRecord, error) {
			return overrideTestRecord, nil
		},
	}, &mockNutritionOverrideItemRepo{})

	override, err := svc.Delete(ctx, testUserID, testID)
	require.NoError(t, err)
	require.NotNil(t, override)
	assert.Equal(t, testID, override.ID)
}

func TestOverrideService_Delete_NotFound(t *testing.T) {
	svc := newOverrideService(&mockNutritionOverrideRepo{
		deleteFn: func(ctx context.Context, userID string, id string) (*models.DailyNutritionOverrideRecord, error) {
			return nil, nil
		},
	}, &mockNutritionOverrideItemRepo{})

	override, err := svc.Delete(ctx, testUserID, testID)
	assert.ErrorIs(t, err, service.ErrOverrideNotFound)
	assert.Nil(t, override)
}

// ----- Override Item CRUD -----

func TestOverrideService_CreateItem_Success(t *testing.T) {
	svc := newOverrideService(&mockNutritionOverrideRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.DailyNutritionOverrideRecord, error) {
			return overrideTestRecord, nil
		},
	}, &mockNutritionOverrideItemRepo{
		createFn: func(ctx context.Context, overrideID string, productID string, amountGrams float64, operation string, mealLabel, notes *string) (*models.DailyNutritionOverrideItemRecord, error) {
			return overrideItemTestRecord, nil
		},
	})

	item, err := svc.CreateItem(ctx, testUserID, models.CreateOverrideItemInput{
		OverrideID: testID, ProductID: "880e8400-e29b-41d4-a716-446655440000",
		AmountGrams: 100, Operation: models.OperationAdd,
	})
	require.NoError(t, err)
	require.NotNil(t, item)
}

func TestOverrideService_CreateItem_InvalidAmount(t *testing.T) {
	svc := newOverrideService(&mockNutritionOverrideRepo{}, &mockNutritionOverrideItemRepo{})
	item, err := svc.CreateItem(ctx, testUserID, models.CreateOverrideItemInput{
		OverrideID: testID, ProductID: "p1",
		AmountGrams: 0, Operation: models.OperationAdd,
	})
	assert.ErrorIs(t, err, service.ErrOverrideItemAmountInvalid)
	assert.Nil(t, item)
}

func TestOverrideService_CreateItem_InvalidOperation(t *testing.T) {
	svc := newOverrideService(&mockNutritionOverrideRepo{}, &mockNutritionOverrideItemRepo{})
	item, err := svc.CreateItem(ctx, testUserID, models.CreateOverrideItemInput{
		OverrideID: testID, ProductID: "p1",
		AmountGrams: 100, Operation: "invalid",
	})
	assert.ErrorIs(t, err, service.ErrOverrideItemOperationInvalid)
	assert.Nil(t, item)
}

func TestOverrideService_CreateItem_OverrideNotFound(t *testing.T) {
	svc := newOverrideService(&mockNutritionOverrideRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.DailyNutritionOverrideRecord, error) {
			return nil, nil
		},
	}, &mockNutritionOverrideItemRepo{})

	item, err := svc.CreateItem(ctx, testUserID, models.CreateOverrideItemInput{
		OverrideID: testID, ProductID: "p1",
		AmountGrams: 100, Operation: models.OperationAdd,
	})
	assert.ErrorIs(t, err, service.ErrOverrideNotFound)
	assert.Nil(t, item)
}

func TestOverrideService_UpdateItem_Success(t *testing.T) {
	svc := newOverrideService(&mockNutritionOverrideRepo{}, &mockNutritionOverrideItemRepo{
		getByIDFn: func(ctx context.Context, id string) (*models.DailyNutritionOverrideItemRecord, error) {
			return overrideItemTestRecord, nil
		},
		updateFn: func(ctx context.Context, id string, amountGrams float64, operation string, mealLabel, notes *string) (*models.DailyNutritionOverrideItemRecord, error) {
			return &models.DailyNutritionOverrideItemRecord{
				ID: id, OverrideID: testID, ProductID: "p1",
				AmountGrams: amountGrams, Operation: operation,
				CreatedAt: "2026-06-20T00:00:00Z", UpdatedAt: "2026-06-20T12:00:00Z",
			}, nil
		},
	})

	item, err := svc.UpdateItem(ctx, testUserID, "item1", models.UpdateOverrideItemInput{
		AmountGrams: ptrFloat64(200),
	})
	require.NoError(t, err)
	require.NotNil(t, item)
}

func TestOverrideService_UpdateItem_InvalidAmount(t *testing.T) {
	svc := newOverrideService(&mockNutritionOverrideRepo{}, &mockNutritionOverrideItemRepo{
		getByIDFn: func(ctx context.Context, id string) (*models.DailyNutritionOverrideItemRecord, error) {
			return overrideItemTestRecord, nil
		},
	})

	item, err := svc.UpdateItem(ctx, testUserID, "item1", models.UpdateOverrideItemInput{
		AmountGrams: ptrFloat64(0),
	})
	assert.ErrorIs(t, err, service.ErrOverrideItemAmountInvalid)
	assert.Nil(t, item)
}

func TestOverrideService_DeleteItem_Success(t *testing.T) {
	svc := newOverrideService(&mockNutritionOverrideRepo{}, &mockNutritionOverrideItemRepo{
		deleteFn: func(ctx context.Context, id string) (*models.DailyNutritionOverrideItemRecord, error) {
			return overrideItemTestRecord, nil
		},
	})

	item, err := svc.DeleteItem(ctx, testUserID, "item1")
	require.NoError(t, err)
	require.NotNil(t, item)
}

func TestOverrideService_DeleteItem_NotFound(t *testing.T) {
	svc := newOverrideService(&mockNutritionOverrideRepo{}, &mockNutritionOverrideItemRepo{
		deleteFn: func(ctx context.Context, id string) (*models.DailyNutritionOverrideItemRecord, error) {
			return nil, nil
		},
	})

	item, err := svc.DeleteItem(ctx, testUserID, "item1")
	assert.ErrorIs(t, err, service.ErrOverrideItemNotFound)
	assert.Nil(t, item)
}
