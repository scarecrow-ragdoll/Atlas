// FILE: apps/api/internal/atlas/service/nutrition_product_service_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for NutritionProductService covering Create, GetByID, ListActive, Update, Delete with validation.
//   SCOPE: Success paths, validation errors (name required, name too long, negative macros), not-found, soft-delete behavior.
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

type mockNutritionProductRepo struct {
	atlasPostgres.NutritionProductRepository
	createFn             func(ctx context.Context, userID string, name string, caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g float64, notes *string) (*models.NutritionProductRecord, error)
	getByIDFn            func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error)
	listActiveFn         func(ctx context.Context, userID string) ([]models.NutritionProductRecord, error)
	updateFn             func(ctx context.Context, userID string, id string, name string, caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g float64, notes *string) (*models.NutritionProductRecord, error)
	softDeleteFn         func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error)
	getByIDIncludeInactiveFn func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error)
}

func (m *mockNutritionProductRepo) Create(ctx context.Context, userID string, name string, caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g float64, notes *string) (*models.NutritionProductRecord, error) {
	return m.createFn(ctx, userID, name, caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g, notes)
}
func (m *mockNutritionProductRepo) GetByID(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
	return m.getByIDFn(ctx, userID, id)
}
func (m *mockNutritionProductRepo) GetByIDIncludeInactive(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
	return m.getByIDIncludeInactiveFn(ctx, userID, id)
}
func (m *mockNutritionProductRepo) ListActive(ctx context.Context, userID string) ([]models.NutritionProductRecord, error) {
	return m.listActiveFn(ctx, userID)
}
func (m *mockNutritionProductRepo) Update(ctx context.Context, userID string, id string, name string, caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g float64, notes *string) (*models.NutritionProductRecord, error) {
	return m.updateFn(ctx, userID, id, name, caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g, notes)
}
func (m *mockNutritionProductRepo) SoftDelete(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
	return m.softDeleteFn(ctx, userID, id)
}

var productTestRecord = &models.NutritionProductRecord{
	ID:              testID,
	UserID:          testUserID,
	Name:            "Chicken Breast",
	CaloriesPer100g: 165,
	ProteinPer100g:  31,
	FatPer100g:      3.6,
	CarbsPer100g:    0,
	Notes:           ptrStr("Boneless skinless"),
	IsActive:        true,
	CreatedAt:       "2026-06-20T00:00:00Z",
	UpdatedAt:       "2026-06-20T00:00:00Z",
}

var productSoftDeletedRecord = &models.NutritionProductRecord{
	ID:              testID,
	UserID:          testUserID,
	Name:            "Chicken Breast",
	CaloriesPer100g: 165,
	ProteinPer100g:  31,
	FatPer100g:      3.6,
	CarbsPer100g:    0,
	Notes:           ptrStr("Boneless skinless"),
	IsActive:        false,
	CreatedAt:       "2026-06-20T00:00:00Z",
	UpdatedAt:       "2026-06-20T00:00:00Z",
}

func newProductService(repo *mockNutritionProductRepo) service.NutritionProductService {
	return service.NewNutritionProductService(repo, zap.NewNop())
}

// ----- Create -----

func TestNutritionProductService_Create_Success(t *testing.T) {
	svc := newProductService(&mockNutritionProductRepo{
		createFn: func(ctx context.Context, userID string, name string, c, p, f, carbs float64, notes *string) (*models.NutritionProductRecord, error) {
			assert.Equal(t, "Chicken Breast", name)
			assert.Equal(t, 165.0, c)
			return productTestRecord, nil
		},
	})

	product, err := svc.Create(ctx, testUserID, models.CreateProductInput{
		Name:            "Chicken Breast",
		CaloriesPer100g: 165,
		ProteinPer100g:  31,
		FatPer100g:      3.6,
		CarbsPer100g:    0,
		Notes:           ptrStr("Boneless skinless"),
	})
	require.NoError(t, err)
	require.NotNil(t, product)
	assert.Equal(t, "Chicken Breast", product.Name)
}

func TestNutritionProductService_Create_NameEmpty(t *testing.T) {
	svc := newProductService(&mockNutritionProductRepo{})
	product, err := svc.Create(ctx, testUserID, models.CreateProductInput{
		Name:            "  ",
		CaloriesPer100g: 100,
		ProteinPer100g:  10,
		FatPer100g:      5,
		CarbsPer100g:    10,
	})
	assert.ErrorIs(t, err, service.ErrProductNameRequired)
	assert.Nil(t, product)
}

func TestNutritionProductService_Create_NameTooLong(t *testing.T) {
	svc := newProductService(&mockNutritionProductRepo{})
	name := make([]byte, 256)
	for i := range name {
		name[i] = 'a'
	}
	product, err := svc.Create(ctx, testUserID, models.CreateProductInput{
		Name:            string(name),
		CaloriesPer100g: 100,
		ProteinPer100g:  10,
		FatPer100g:      5,
		CarbsPer100g:    10,
	})
	assert.ErrorIs(t, err, service.ErrProductNameTooLong)
	assert.Nil(t, product)
}

func TestNutritionProductService_Create_NegativeMacro(t *testing.T) {
	svc := newProductService(&mockNutritionProductRepo{})
	product, err := svc.Create(ctx, testUserID, models.CreateProductInput{
		Name:            "Bad",
		CaloriesPer100g: 100,
		ProteinPer100g:  -1,
		FatPer100g:      5,
		CarbsPer100g:    10,
	})
	assert.ErrorIs(t, err, service.ErrProductMacroNegative)
	assert.Nil(t, product)
}

// ----- GetByID -----

func TestNutritionProductService_GetByID_Success(t *testing.T) {
	svc := newProductService(&mockNutritionProductRepo{
		getByIDIncludeInactiveFn: func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
			return productTestRecord, nil
		},
	})

	product, err := svc.GetByID(ctx, testUserID, testID)
	require.NoError(t, err)
	require.NotNil(t, product)
	assert.Equal(t, "Chicken Breast", product.Name)
	assert.True(t, product.IsActive)
}

func TestNutritionProductService_GetByID_SoftDeleted(t *testing.T) {
	svc := newProductService(&mockNutritionProductRepo{
		getByIDIncludeInactiveFn: func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
			return productSoftDeletedRecord, nil
		},
	})

	product, err := svc.GetByID(ctx, testUserID, testID)
	require.NoError(t, err)
	require.NotNil(t, product)
	assert.False(t, product.IsActive)
}

func TestNutritionProductService_GetByID_NotFound(t *testing.T) {
	svc := newProductService(&mockNutritionProductRepo{
		getByIDIncludeInactiveFn: func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
			return nil, nil
		},
	})

	product, err := svc.GetByID(ctx, testUserID, testID)
	assert.ErrorIs(t, err, service.ErrProductNotFound)
	assert.Nil(t, product)
}

// ----- ListActive -----

func TestNutritionProductService_ListActive_Success(t *testing.T) {
	svc := newProductService(&mockNutritionProductRepo{
		listActiveFn: func(ctx context.Context, userID string) ([]models.NutritionProductRecord, error) {
			return []models.NutritionProductRecord{*productTestRecord}, nil
		},
	})

	products, err := svc.ListActive(ctx, testUserID)
	require.NoError(t, err)
	assert.Len(t, products, 1)
}

func TestNutritionProductService_ListActive_Empty(t *testing.T) {
	svc := newProductService(&mockNutritionProductRepo{
		listActiveFn: func(ctx context.Context, userID string) ([]models.NutritionProductRecord, error) {
			return []models.NutritionProductRecord{}, nil
		},
	})

	products, err := svc.ListActive(ctx, testUserID)
	require.NoError(t, err)
	assert.Len(t, products, 0)
}

// ----- Update -----

func TestNutritionProductService_Update_Success(t *testing.T) {
	svc := newProductService(&mockNutritionProductRepo{
		getByIDIncludeInactiveFn: func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
			return productTestRecord, nil
		},
		updateFn: func(ctx context.Context, userID string, id string, name string, c, p, f, carbs float64, notes *string) (*models.NutritionProductRecord, error) {
			assert.Equal(t, "Updated Chicken", name)
			return &models.NutritionProductRecord{
				ID: id, UserID: userID, Name: name,
				CaloriesPer100g: c, ProteinPer100g: p, FatPer100g: f, CarbsPer100g: carbs,
				Notes: notes, IsActive: true,
				CreatedAt: "2026-06-20T00:00:00Z", UpdatedAt: "2026-06-20T12:00:00Z",
			}, nil
		},
	})

	product, err := svc.Update(ctx, testUserID, testID, models.UpdateProductInput{
		Name:            ptrStr("Updated Chicken"),
		CaloriesPer100g: ptrFloat64(200),
	})
	require.NoError(t, err)
	require.NotNil(t, product)
	assert.Equal(t, "Updated Chicken", product.Name)
}

func TestNutritionProductService_Update_NotFound(t *testing.T) {
	svc := newProductService(&mockNutritionProductRepo{
		getByIDIncludeInactiveFn: func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
			return nil, nil
		},
	})

	product, err := svc.Update(ctx, testUserID, testID, models.UpdateProductInput{Name: ptrStr("X")})
	assert.ErrorIs(t, err, service.ErrProductNotFound)
	assert.Nil(t, product)
}

func TestNutritionProductService_Update_NegativeMacro(t *testing.T) {
	svc := newProductService(&mockNutritionProductRepo{
		getByIDIncludeInactiveFn: func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
			return productTestRecord, nil
		},
	})

	product, err := svc.Update(ctx, testUserID, testID, models.UpdateProductInput{
		CaloriesPer100g: ptrFloat64(-10),
	})
	assert.ErrorIs(t, err, service.ErrProductMacroNegative)
	assert.Nil(t, product)
}

// ----- Delete (soft-delete) -----

func TestNutritionProductService_Delete_Success(t *testing.T) {
	svc := newProductService(&mockNutritionProductRepo{
		softDeleteFn: func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
			return productSoftDeletedRecord, nil
		},
	})

	product, err := svc.Delete(ctx, testUserID, testID)
	require.NoError(t, err)
	require.NotNil(t, product)
	assert.False(t, product.IsActive)
}

func TestNutritionProductService_Delete_NotFound(t *testing.T) {
	svc := newProductService(&mockNutritionProductRepo{
		softDeleteFn: func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
			return nil, nil
		},
	})

	product, err := svc.Delete(ctx, testUserID, testID)
	assert.ErrorIs(t, err, service.ErrProductNotFound)
	assert.Nil(t, product)
}

