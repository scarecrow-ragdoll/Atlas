// FILE: apps/api/internal/atlas/service/nutrition_template_item_service_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for NutritionTemplateItemService ownership and mutation behavior.
//   SCOPE: Product ownership checks on create, parent template ownership checks, user-scoped update/delete, and successful owned create/update/delete paths.
//   DEPENDS: apps/api/internal/atlas/service, apps/api/internal/atlas/repository/postgres (mocks), apps/api/internal/atlas/models.
//   LINKS: M-API-NUTRITION / V-M-API-NUTRITION.
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
	"monorepo-template/apps/api/internal/atlas/service"
)

const templateItemProductID = "880e8400-e29b-41d4-a716-446655440000"

var templateItemTestRecord = &models.NutritionTemplateItemRecord{
	ID:          "770e8400-e29b-41d4-a716-446655440000",
	TemplateID:  testID,
	ProductID:   templateItemProductID,
	AmountGrams: 150,
	MealLabel:   ptrStr("Lunch"),
	Notes:       ptrStr("Grilled"),
	CreatedAt:   "2026-06-20T00:00:00Z",
	UpdatedAt:   "2026-06-20T00:00:00Z",
}

var templateItemProductRecord = &models.NutritionProductRecord{
	ID:              templateItemProductID,
	UserID:          testUserID,
	Name:            "Oats",
	CaloriesPer100g: 389,
	ProteinPer100g:  16.9,
	FatPer100g:      6.9,
	CarbsPer100g:    66.3,
	IsActive:        true,
	CreatedAt:       "2026-06-20T00:00:00Z",
	UpdatedAt:       "2026-06-20T00:00:00Z",
}

func newTemplateItemService(
	itemRepo *mockNutritionTemplateItemRepo,
	tmplRepo *mockNutritionTemplateRepo,
	productRepo *mockNutritionProductRepo,
) service.NutritionTemplateItemService {
	return service.NewNutritionTemplateItemService(itemRepo, tmplRepo, productRepo, zap.NewNop())
}

func TestNutritionTemplateItemService_Create_RejectsProductNotOwnedByUser(t *testing.T) {
	svc := newTemplateItemService(&mockNutritionTemplateItemRepo{
		createFn: func(ctx context.Context, templateID string, productID string, amountGrams float64, mealLabel, notes *string) (*models.NutritionTemplateItemRecord, error) {
			t.Fatal("Create should not persist item when product ownership check fails")
			return nil, nil
		},
	}, &mockNutritionTemplateRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error) {
			assert.Equal(t, testUserID, userID)
			assert.Equal(t, testID, id)
			return tmplTestRecord, nil
		},
	}, &mockNutritionProductRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
			assert.Equal(t, testUserID, userID)
			assert.Equal(t, templateItemProductID, id)
			return nil, nil
		},
	})

	item, err := svc.Create(ctx, testUserID, models.CreateTemplateItemInput{
		TemplateID:  testID,
		ProductID:   templateItemProductID,
		AmountGrams: 150,
	})

	assert.ErrorIs(t, err, service.ErrProductNotFound)
	assert.Nil(t, item)
}

func TestNutritionTemplateItemService_Create_VerifiesParentTemplateBelongsToUser(t *testing.T) {
	svc := newTemplateItemService(&mockNutritionTemplateItemRepo{
		createFn: func(ctx context.Context, templateID string, productID string, amountGrams float64, mealLabel, notes *string) (*models.NutritionTemplateItemRecord, error) {
			t.Fatal("Create should not persist item when template ownership check fails")
			return nil, nil
		},
	}, &mockNutritionTemplateRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error) {
			assert.Equal(t, testUserID, userID)
			return nil, nil
		},
	}, &mockNutritionProductRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
			t.Fatal("Create should check template ownership before product ownership")
			return nil, nil
		},
	})

	item, err := svc.Create(ctx, testUserID, models.CreateTemplateItemInput{
		TemplateID:  testID,
		ProductID:   templateItemProductID,
		AmountGrams: 150,
	})

	assert.ErrorIs(t, err, service.ErrTemplateNotFound)
	assert.Nil(t, item)
}

func TestNutritionTemplateItemService_Create_SuccessWithOwnedTemplateAndProduct(t *testing.T) {
	svc := newTemplateItemService(&mockNutritionTemplateItemRepo{
		createFn: func(ctx context.Context, templateID string, productID string, amountGrams float64, mealLabel, notes *string) (*models.NutritionTemplateItemRecord, error) {
			assert.Equal(t, testID, templateID)
			assert.Equal(t, templateItemProductID, productID)
			assert.Equal(t, 150.0, amountGrams)
			return templateItemTestRecord, nil
		},
	}, &mockNutritionTemplateRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error) {
			assert.Equal(t, testUserID, userID)
			return tmplTestRecord, nil
		},
	}, &mockNutritionProductRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
			assert.Equal(t, testUserID, userID)
			return templateItemProductRecord, nil
		},
	})

	item, err := svc.Create(ctx, testUserID, models.CreateTemplateItemInput{
		TemplateID:  testID,
		ProductID:   templateItemProductID,
		AmountGrams: 150,
		MealLabel:   ptrStr("Lunch"),
	})

	require.NoError(t, err)
	require.NotNil(t, item)
	assert.Equal(t, templateItemTestRecord.ID, item.ID)
}

func TestNutritionTemplateItemService_Create_RejectsArchivedProduct(t *testing.T) {
	archivedProduct := *templateItemProductRecord
	archivedProduct.IsActive = false
	svc := newTemplateItemService(&mockNutritionTemplateItemRepo{
		createFn: func(ctx context.Context, templateID string, productID string, amountGrams float64, mealLabel, notes *string) (*models.NutritionTemplateItemRecord, error) {
			t.Fatal("Create should not persist item when product is archived")
			return nil, nil
		},
	}, &mockNutritionTemplateRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.NutritionTemplateRecord, error) {
			return tmplTestRecord, nil
		},
	}, &mockNutritionProductRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
			return &archivedProduct, nil
		},
	})

	item, err := svc.Create(ctx, testUserID, models.CreateTemplateItemInput{
		TemplateID:  testID,
		ProductID:   templateItemProductID,
		AmountGrams: 150,
	})

	assert.ErrorIs(t, err, service.ErrProductNotFound)
	assert.Nil(t, item)
}

func TestNutritionTemplateItemService_Update_RejectsItemOwnedByAnotherUser(t *testing.T) {
	svc := newTemplateItemService(&mockNutritionTemplateItemRepo{
		getByIDForUserFn: func(ctx context.Context, userID string, id string) (*models.NutritionTemplateItemRecord, error) {
			assert.Equal(t, testUserID, userID)
			assert.Equal(t, templateItemTestRecord.ID, id)
			return nil, nil
		},
		updateFn: func(ctx context.Context, id string, amountGrams float64, mealLabel, notes *string) (*models.NutritionTemplateItemRecord, error) {
			t.Fatal("Update should not persist cross-user item changes")
			return nil, nil
		},
	}, &mockNutritionTemplateRepo{}, &mockNutritionProductRepo{})

	item, err := svc.Update(ctx, testUserID, templateItemTestRecord.ID, models.UpdateTemplateItemInput{
		AmountGrams: ptrFloat64(175),
	})

	assert.ErrorIs(t, err, service.ErrTemplateItemNotFound)
	assert.Nil(t, item)
}

func TestNutritionTemplateItemService_Update_SuccessForOwnedItem(t *testing.T) {
	svc := newTemplateItemService(&mockNutritionTemplateItemRepo{
		getByIDForUserFn: func(ctx context.Context, userID string, id string) (*models.NutritionTemplateItemRecord, error) {
			assert.Equal(t, testUserID, userID)
			return templateItemTestRecord, nil
		},
		updateFn: func(ctx context.Context, id string, amountGrams float64, mealLabel, notes *string) (*models.NutritionTemplateItemRecord, error) {
			assert.Equal(t, templateItemTestRecord.ID, id)
			assert.Equal(t, 175.0, amountGrams)
			updated := *templateItemTestRecord
			updated.AmountGrams = amountGrams
			updated.MealLabel = mealLabel
			return &updated, nil
		},
	}, &mockNutritionTemplateRepo{}, &mockNutritionProductRepo{})

	item, err := svc.Update(ctx, testUserID, templateItemTestRecord.ID, models.UpdateTemplateItemInput{
		AmountGrams: ptrFloat64(175),
		MealLabel:   ptrStr("Dinner"),
	})

	require.NoError(t, err)
	require.NotNil(t, item)
	assert.Equal(t, 175.0, item.AmountGrams)
	assert.Equal(t, "Dinner", *item.MealLabel)
}

func TestNutritionTemplateItemService_Delete_RejectsItemOwnedByAnotherUser(t *testing.T) {
	svc := newTemplateItemService(&mockNutritionTemplateItemRepo{
		getByIDForUserFn: func(ctx context.Context, userID string, id string) (*models.NutritionTemplateItemRecord, error) {
			assert.Equal(t, testUserID, userID)
			assert.Equal(t, templateItemTestRecord.ID, id)
			return nil, nil
		},
		deleteFn: func(ctx context.Context, id string) (*models.NutritionTemplateItemRecord, error) {
			t.Fatal("Delete should not remove cross-user items")
			return nil, nil
		},
	}, &mockNutritionTemplateRepo{}, &mockNutritionProductRepo{})

	item, err := svc.Delete(ctx, testUserID, templateItemTestRecord.ID)

	assert.ErrorIs(t, err, service.ErrTemplateItemNotFound)
	assert.Nil(t, item)
}

func TestNutritionTemplateItemService_Delete_SuccessForOwnedItem(t *testing.T) {
	svc := newTemplateItemService(&mockNutritionTemplateItemRepo{
		getByIDForUserFn: func(ctx context.Context, userID string, id string) (*models.NutritionTemplateItemRecord, error) {
			assert.Equal(t, testUserID, userID)
			return templateItemTestRecord, nil
		},
		deleteFn: func(ctx context.Context, id string) (*models.NutritionTemplateItemRecord, error) {
			assert.Equal(t, templateItemTestRecord.ID, id)
			return templateItemTestRecord, nil
		},
	}, &mockNutritionTemplateRepo{}, &mockNutritionProductRepo{})

	item, err := svc.Delete(ctx, testUserID, templateItemTestRecord.ID)

	require.NoError(t, err)
	require.NotNil(t, item)
	assert.Equal(t, templateItemTestRecord.ID, item.ID)
}
