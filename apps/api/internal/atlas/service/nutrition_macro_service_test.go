// FILE: apps/api/internal/atlas/service/nutrition_macro_service_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for NutritionMacroService covering macro calculation for template week with override operations (add/subtract/replace).
//   SCOPE: Empty template (0 macros), template with items, override ADD/SUBTRACT/REPLACE, soft-deleted products contribute 0, no override for requested date.
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
	"monorepo-template/apps/api/internal/atlas/service"
)

type mockMacroRepos struct {
	tmplRepo         *mockNutritionTemplateRepo
	itemRepo         *mockNutritionTemplateItemRepo
	overrideRepo     *mockNutritionOverrideRepo
	overrideItemRepo *mockNutritionOverrideItemRepo
	productRepo      *mockNutritionProductRepo
}

func newMacroService(m *mockMacroRepos) service.NutritionMacroService {
	return service.NewNutritionMacroService(m.tmplRepo, m.itemRepo, m.overrideRepo, m.overrideItemRepo, m.productRepo, zap.NewNop())
}

var macroProductRecord = &models.NutritionProductRecord{
	ID: "880e8400-e29b-41d4-a716-446655440000", UserID: testUserID,
	Name: "Oats", CaloriesPer100g: 389, ProteinPer100g: 16.9, FatPer100g: 6.9, CarbsPer100g: 66.3,
	IsActive: true, CreatedAt: "2026-06-20T00:00:00Z", UpdatedAt: "2026-06-20T00:00:00Z",
}

var macroSoftDeletedProduct = &models.NutritionProductRecord{
	ID: "990e8400-e29b-41d4-a716-446655440000", UserID: testUserID,
	Name: "Old Product", CaloriesPer100g: 100, ProteinPer100g: 10, FatPer100g: 5, CarbsPer100g: 20,
	IsActive: false, CreatedAt: "2026-06-20T00:00:00Z", UpdatedAt: "2026-06-20T00:00:00Z",
}

func TestNutritionMacroService_EmptyWeek(t *testing.T) {
	m := &mockMacroRepos{
		tmplRepo: &mockNutritionTemplateRepo{
			getByWeekFn: func(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error) {
				return nil, nil
			},
		},
		itemRepo:         &mockNutritionTemplateItemRepo{},
		overrideRepo:     &mockNutritionOverrideRepo{},
		overrideItemRepo: &mockNutritionOverrideItemRepo{},
		productRepo:      &mockNutritionProductRepo{},
	}
	svc := newMacroService(m)

	macros, err := svc.Calculate(ctx, testUserID, "2026-06-15", "2026-06-16")
	require.NoError(t, err)
	require.NotNil(t, macros)
	assert.Equal(t, 0.0, macros.Calories)
	assert.Equal(t, 0.0, macros.Protein)
	assert.Equal(t, 0.0, macros.Fat)
	assert.Equal(t, 0.0, macros.Carbs)
}

func TestNutritionMacroService_TemplateWithItems(t *testing.T) {
	m := &mockMacroRepos{
		tmplRepo: &mockNutritionTemplateRepo{
			getByWeekFn: func(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error) {
				return tmplTestRecord, nil
			},
		},
		itemRepo: &mockNutritionTemplateItemRepo{
			listByTemplateFn: func(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error) {
				return []models.NutritionTemplateItemRecord{
					{ID: "i1", TemplateID: testID, ProductID: macroProductRecord.ID, AmountGrams: 100, CreatedAt: "", UpdatedAt: ""},
				}, nil
			},
		},
		overrideRepo: &mockNutritionOverrideRepo{},
		overrideItemRepo: &mockNutritionOverrideItemRepo{},
		productRepo: &mockNutritionProductRepo{
			getByIDIncludeInactiveFn: func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
				return macroProductRecord, nil
			},
		},
	}
	svc := newMacroService(m)

	macros, err := svc.Calculate(ctx, testUserID, "2026-06-15", "")
	require.NoError(t, err)
	require.NotNil(t, macros)
	assert.InDelta(t, 389, macros.Calories, 0.01)
	assert.InDelta(t, 16.9, macros.Protein, 0.01)
	assert.InDelta(t, 6.9, macros.Fat, 0.01)
	assert.InDelta(t, 66.3, macros.Carbs, 0.01)
}

func TestNutritionMacroService_OverrideAdd(t *testing.T) {
	m := &mockMacroRepos{
		tmplRepo: &mockNutritionTemplateRepo{
			getByWeekFn: func(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error) {
				return tmplTestRecord, nil
			},
		},
		itemRepo: &mockNutritionTemplateItemRepo{
			listByTemplateFn: func(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error) {
				return []models.NutritionTemplateItemRecord{
					{ID: "i1", TemplateID: testID, ProductID: macroProductRecord.ID, AmountGrams: 100},
				}, nil
			},
		},
		overrideRepo: &mockNutritionOverrideRepo{
			getByDateFn: func(ctx context.Context, userID string, date string) (*models.DailyNutritionOverrideRecord, error) {
				return &models.DailyNutritionOverrideRecord{ID: "o1", UserID: userID, Date: models.MustDate("2026-06-16")}, nil
			},
		},
		overrideItemRepo: &mockNutritionOverrideItemRepo{
			listByOverrideFn: func(ctx context.Context, overrideID string) ([]models.DailyNutritionOverrideItemRecord, error) {
				return []models.DailyNutritionOverrideItemRecord{
					{ID: "oi1", OverrideID: "o1", ProductID: macroProductRecord.ID, AmountGrams: 50, Operation: "add"},
				}, nil
			},
		},
		productRepo: &mockNutritionProductRepo{
			getByIDIncludeInactiveFn: func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
				return macroProductRecord, nil
			},
		},
	}
	svc := newMacroService(m)

	macros, err := svc.Calculate(ctx, testUserID, "2026-06-15", "2026-06-16")
	require.NoError(t, err)
	// Template: 100g oats = 389 cal. Override ADD 50g oats = 194.5 cal. Total = 583.5
	assert.InDelta(t, 583.5, macros.Calories, 0.01)
}

func TestNutritionMacroService_OverrideSubtract(t *testing.T) {
	m := &mockMacroRepos{
		tmplRepo: &mockNutritionTemplateRepo{
			getByWeekFn: func(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error) {
				return tmplTestRecord, nil
			},
		},
		itemRepo: &mockNutritionTemplateItemRepo{
			listByTemplateFn: func(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error) {
				return []models.NutritionTemplateItemRecord{
					{ID: "i1", TemplateID: testID, ProductID: macroProductRecord.ID, AmountGrams: 100},
				}, nil
			},
		},
		overrideRepo: &mockNutritionOverrideRepo{
			getByDateFn: func(ctx context.Context, userID string, date string) (*models.DailyNutritionOverrideRecord, error) {
				return &models.DailyNutritionOverrideRecord{ID: "o1", UserID: userID, Date: models.MustDate("2026-06-16")}, nil
			},
		},
		overrideItemRepo: &mockNutritionOverrideItemRepo{
			listByOverrideFn: func(ctx context.Context, overrideID string) ([]models.DailyNutritionOverrideItemRecord, error) {
				return []models.DailyNutritionOverrideItemRecord{
					{ID: "oi1", OverrideID: "o1", ProductID: macroProductRecord.ID, AmountGrams: 30, Operation: "subtract"},
				}, nil
			},
		},
		productRepo: &mockNutritionProductRepo{
			getByIDIncludeInactiveFn: func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
				return macroProductRecord, nil
			},
		},
	}
	svc := newMacroService(m)

	macros, err := svc.Calculate(ctx, testUserID, "2026-06-15", "2026-06-16")
	require.NoError(t, err)
	// Template: 389 cal. Subtract 30g oats = 116.7 cal. Total = 272.3
	assert.InDelta(t, 272.3, macros.Calories, 0.01)
}

func TestNutritionMacroService_OverrideReplace(t *testing.T) {
	m := &mockMacroRepos{
		tmplRepo: &mockNutritionTemplateRepo{
			getByWeekFn: func(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error) {
				return tmplTestRecord, nil
			},
		},
		itemRepo: &mockNutritionTemplateItemRepo{
			listByTemplateFn: func(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error) {
				return []models.NutritionTemplateItemRecord{
					{ID: "i1", TemplateID: testID, ProductID: macroProductRecord.ID, AmountGrams: 100},
				}, nil
			},
		},
		overrideRepo: &mockNutritionOverrideRepo{
			getByDateFn: func(ctx context.Context, userID string, date string) (*models.DailyNutritionOverrideRecord, error) {
				return &models.DailyNutritionOverrideRecord{ID: "o1", UserID: userID, Date: models.MustDate("2026-06-16")}, nil
			},
		},
		overrideItemRepo: &mockNutritionOverrideItemRepo{
			listByOverrideFn: func(ctx context.Context, overrideID string) ([]models.DailyNutritionOverrideItemRecord, error) {
				return []models.DailyNutritionOverrideItemRecord{
					{ID: "oi1", OverrideID: "o1", ProductID: macroProductRecord.ID, AmountGrams: 80, Operation: "replace"},
				}, nil
			},
		},
		productRepo: &mockNutritionProductRepo{
			getByIDIncludeInactiveFn: func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
				return macroProductRecord, nil
			},
		},
	}
	svc := newMacroService(m)

	macros, err := svc.Calculate(ctx, testUserID, "2026-06-15", "2026-06-16")
	require.NoError(t, err)
	// Replace 100g oats (389 cal) with 80g oats (311.2 cal)
	assert.InDelta(t, 311.2, macros.Calories, 0.01)
}

func TestNutritionMacroService_SoftDeletedProduct(t *testing.T) {
	m := &mockMacroRepos{
		tmplRepo: &mockNutritionTemplateRepo{
			getByWeekFn: func(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error) {
				return tmplTestRecord, nil
			},
		},
		itemRepo: &mockNutritionTemplateItemRepo{
			listByTemplateFn: func(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error) {
				return []models.NutritionTemplateItemRecord{
					{ID: "i1", TemplateID: testID, ProductID: macroSoftDeletedProduct.ID, AmountGrams: 100},
				}, nil
			},
		},
		overrideRepo: &mockNutritionOverrideRepo{},
		overrideItemRepo: &mockNutritionOverrideItemRepo{},
		productRepo: &mockNutritionProductRepo{
			getByIDIncludeInactiveFn: func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
				return macroSoftDeletedProduct, nil
			},
		},
	}
	svc := newMacroService(m)

	macros, err := svc.Calculate(ctx, testUserID, "2026-06-15", "")
	require.NoError(t, err)
	// Soft-deleted product contributes 0
	assert.Equal(t, 0.0, macros.Calories)
	assert.Equal(t, 0.0, macros.Protein)
	assert.Equal(t, 0.0, macros.Fat)
	assert.Equal(t, 0.0, macros.Carbs)
}

func TestNutritionMacroService_NoOverride(t *testing.T) {
	m := &mockMacroRepos{
		tmplRepo: &mockNutritionTemplateRepo{
			getByWeekFn: func(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error) {
				return tmplTestRecord, nil
			},
		},
		itemRepo: &mockNutritionTemplateItemRepo{
			listByTemplateFn: func(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error) {
				return []models.NutritionTemplateItemRecord{
					{ID: "i1", TemplateID: testID, ProductID: macroProductRecord.ID, AmountGrams: 100},
				}, nil
			},
		},
		overrideRepo: &mockNutritionOverrideRepo{
			getByDateFn: func(ctx context.Context, userID string, date string) (*models.DailyNutritionOverrideRecord, error) {
				return nil, nil
			},
		},
		overrideItemRepo: &mockNutritionOverrideItemRepo{},
		productRepo: &mockNutritionProductRepo{
			getByIDIncludeInactiveFn: func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
				return macroProductRecord, nil
			},
		},
	}
	svc := newMacroService(m)

	macros, err := svc.Calculate(ctx, testUserID, "2026-06-15", "2026-06-16")
	require.NoError(t, err)
	// No override found for date, should return template values
	assert.InDelta(t, 389, macros.Calories, 0.01)
}
