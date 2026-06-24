// FILE: apps/api/internal/atlas/service/daily_nutrition_legacy_resolver_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for deterministic legacy daily nutrition override resolution into factual food-log entries.
//   SCOPE: ADD, REPLACE, SUBTRACT, repeated same-product template rows, no-template ADD, and unresolved diagnostics for ambiguous or incomplete legacy rows.
//   DEPENDS: apps/api/internal/atlas/service, apps/api/internal/atlas/models, atlas postgres repository mock interfaces.
//   LINKS: M-API-NUTRITION / V-M-API-NUTRITION.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Task 6 RED coverage for legacy override resolver behavior.
// END_CHANGE_SUMMARY

package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/atlas/service"
)

const (
	legacyTemplateID = "100e8400-e29b-41d4-a716-446655440000"
	legacyOverrideID = "200e8400-e29b-41d4-a716-446655440000"
	legacyChickenID  = "300e8400-e29b-41d4-a716-446655440000"
	legacyRiceID     = "400e8400-e29b-41d4-a716-446655440000"
	legacyYogurtID   = "500e8400-e29b-41d4-a716-446655440000"
	legacyMissingID  = "600e8400-e29b-41d4-a716-446655440000"
)

func newLegacyResolver(
	templateRepo *mockNutritionTemplateRepo,
	templateItemRepo *mockNutritionTemplateItemRepo,
	overrideRepo *mockNutritionOverrideRepo,
	overrideItemRepo *mockNutritionOverrideItemRepo,
	productRepo *mockNutritionProductRepo,
) service.DailyNutritionLegacyResolver {
	return service.NewDailyNutritionLegacyResolver(templateRepo, templateItemRepo, overrideRepo, overrideItemRepo, productRepo)
}

func legacyTemplateRecord() *models.NutritionTemplateRecord {
	return &models.NutritionTemplateRecord{
		ID:            legacyTemplateID,
		UserID:        testUserID,
		WeekStartDate: models.MustDate("2026-06-22"),
		CreatedAt:     "2026-06-22T00:00:00Z",
		UpdatedAt:     "2026-06-22T00:00:00Z",
	}
}

func legacyOverrideRecord() *models.DailyNutritionOverrideRecord {
	return &models.DailyNutritionOverrideRecord{
		ID:        legacyOverrideID,
		UserID:    testUserID,
		Date:      models.MustDate("2026-06-24"),
		CreatedAt: "2026-06-24T00:00:00Z",
		UpdatedAt: "2026-06-24T00:00:00Z",
	}
}

func legacyProduct(id string) *models.NutritionProductRecord {
	products := map[string]*models.NutritionProductRecord{
		legacyChickenID: {
			ID: legacyChickenID, UserID: testUserID, Name: "Chicken",
			CaloriesPer100g: 100, ProteinPer100g: 20, FatPer100g: 2, CarbsPer100g: 0,
			IsActive: true, CreatedAt: "2026-06-20T00:00:00Z", UpdatedAt: "2026-06-20T00:00:00Z",
		},
		legacyRiceID: {
			ID: legacyRiceID, UserID: testUserID, Name: "Rice",
			CaloriesPer100g: 200, ProteinPer100g: 4, FatPer100g: 1, CarbsPer100g: 40,
			IsActive: true, CreatedAt: "2026-06-20T00:00:00Z", UpdatedAt: "2026-06-20T00:00:00Z",
		},
		legacyYogurtID: {
			ID: legacyYogurtID, UserID: testUserID, Name: "Yogurt",
			CaloriesPer100g: 80, ProteinPer100g: 8, FatPer100g: 2, CarbsPer100g: 6,
			IsActive: true, CreatedAt: "2026-06-20T00:00:00Z", UpdatedAt: "2026-06-20T00:00:00Z",
		},
	}
	return products[id]
}

func legacyProductRepo(t *testing.T) *mockNutritionProductRepo {
	t.Helper()
	return &mockNutritionProductRepo{
		getByIDIncludeInactiveFn: func(ctx context.Context, userID string, id string) (*models.NutritionProductRecord, error) {
			assert.Equal(t, testUserID, userID)
			return legacyProduct(id), nil
		},
	}
}

func legacyEntryByProduct(t *testing.T, entries []models.DailyNutritionEntry, productID string) models.DailyNutritionEntry {
	t.Helper()
	for _, entry := range entries {
		if entry.ProductID == productID {
			return entry
		}
	}
	t.Fatalf("entry for product %s not found", productID)
	return models.DailyNutritionEntry{}
}

func TestDailyNutritionLegacyResolver_ResolvesAddReplaceSubtract(t *testing.T) {
	resolver := newLegacyResolver(
		&mockNutritionTemplateRepo{
			getByWeekFn: func(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error) {
				assert.Equal(t, "2026-06-22", weekStartDate)
				return legacyTemplateRecord(), nil
			},
		},
		&mockNutritionTemplateItemRepo{
			listByTemplateFn: func(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error) {
				assert.Equal(t, legacyTemplateID, templateID)
				return []models.NutritionTemplateItemRecord{
					{ID: "base-chicken", TemplateID: legacyTemplateID, ProductID: legacyChickenID, AmountGrams: 100, MealLabel: ptrStr("Lunch")},
					{ID: "base-rice", TemplateID: legacyTemplateID, ProductID: legacyRiceID, AmountGrams: 200, MealLabel: ptrStr("Lunch")},
				}, nil
			},
		},
		&mockNutritionOverrideRepo{
			getByDateFn: func(ctx context.Context, userID string, date string) (*models.DailyNutritionOverrideRecord, error) {
				assert.Equal(t, "2026-06-24", date)
				return legacyOverrideRecord(), nil
			},
		},
		&mockNutritionOverrideItemRepo{
			listByOverrideFn: func(ctx context.Context, overrideID string) ([]models.DailyNutritionOverrideItemRecord, error) {
				assert.Equal(t, legacyOverrideID, overrideID)
				return []models.DailyNutritionOverrideItemRecord{
					{ID: "replace-chicken", OverrideID: legacyOverrideID, ProductID: legacyChickenID, AmountGrams: 150, Operation: string(models.OperationReplace), MealLabel: ptrStr("Dinner")},
					{ID: "subtract-rice", OverrideID: legacyOverrideID, ProductID: legacyRiceID, AmountGrams: 50, Operation: string(models.OperationSubtract)},
					{ID: "add-yogurt", OverrideID: legacyOverrideID, ProductID: legacyYogurtID, AmountGrams: 50, Operation: string(models.OperationAdd), MealLabel: ptrStr("Snack")},
				}, nil
			},
		},
		legacyProductRepo(t),
	)

	day, err := resolver.Resolve(ctx, testUserID, models.MustDate("2026-06-24"))
	require.NoError(t, err)
	require.NotNil(t, day)
	assert.Equal(t, models.LegacyResolutionResolved, day.Status)
	require.Len(t, day.ResolvedEntries, 3)
	assert.Equal(t, 150.0, legacyEntryByProduct(t, day.ResolvedEntries, legacyChickenID).AmountGrams)
	assert.Equal(t, 150.0, legacyEntryByProduct(t, day.ResolvedEntries, legacyRiceID).AmountGrams)
	assert.Equal(t, 50.0, legacyEntryByProduct(t, day.ResolvedEntries, legacyYogurtID).AmountGrams)
	assert.InDelta(t, 490, day.Totals.Calories, 0.001)
	assert.InDelta(t, day.LegacyTotals.Calories, day.Totals.Calories, 0.001)
	assert.Len(t, day.RawOperations, 3)
}

func TestDailyNutritionLegacyResolver_ResolvesRepeatedSameProductTemplateRows(t *testing.T) {
	resolver := newLegacyResolver(
		&mockNutritionTemplateRepo{
			getByWeekFn: func(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error) {
				return legacyTemplateRecord(), nil
			},
		},
		&mockNutritionTemplateItemRepo{
			listByTemplateFn: func(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error) {
				return []models.NutritionTemplateItemRecord{
					{ID: "base-rice-a", TemplateID: legacyTemplateID, ProductID: legacyRiceID, AmountGrams: 100, MealLabel: ptrStr("Lunch")},
					{ID: "base-rice-b", TemplateID: legacyTemplateID, ProductID: legacyRiceID, AmountGrams: 75, MealLabel: ptrStr("Dinner")},
				}, nil
			},
		},
		&mockNutritionOverrideRepo{
			getByDateFn: func(ctx context.Context, userID string, date string) (*models.DailyNutritionOverrideRecord, error) {
				return legacyOverrideRecord(), nil
			},
		},
		&mockNutritionOverrideItemRepo{
			listByOverrideFn: func(ctx context.Context, overrideID string) ([]models.DailyNutritionOverrideItemRecord, error) {
				return []models.DailyNutritionOverrideItemRecord{
					{ID: "subtract-rice", OverrideID: legacyOverrideID, ProductID: legacyRiceID, AmountGrams: 25, Operation: string(models.OperationSubtract)},
				}, nil
			},
		},
		legacyProductRepo(t),
	)

	day, err := resolver.Resolve(ctx, testUserID, models.MustDate("2026-06-24"))
	require.NoError(t, err)
	require.NotNil(t, day)
	assert.Equal(t, models.LegacyResolutionResolved, day.Status)
	require.Len(t, day.ResolvedEntries, 1)
	assert.Equal(t, legacyRiceID, day.ResolvedEntries[0].ProductID)
	assert.Equal(t, 150.0, day.ResolvedEntries[0].AmountGrams)
	assert.InDelta(t, 300, day.Totals.Calories, 0.001)
}

func TestDailyNutritionLegacyResolver_PreservesSameProductAddBeforeReplace(t *testing.T) {
	resolver := newLegacyResolver(
		&mockNutritionTemplateRepo{
			getByWeekFn: func(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error) {
				return legacyTemplateRecord(), nil
			},
		},
		&mockNutritionTemplateItemRepo{
			listByTemplateFn: func(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error) {
				return []models.NutritionTemplateItemRecord{
					{ID: "base-chicken", TemplateID: legacyTemplateID, ProductID: legacyChickenID, AmountGrams: 100, MealLabel: ptrStr("Lunch")},
				}, nil
			},
		},
		&mockNutritionOverrideRepo{
			getByDateFn: func(ctx context.Context, userID string, date string) (*models.DailyNutritionOverrideRecord, error) {
				return legacyOverrideRecord(), nil
			},
		},
		&mockNutritionOverrideItemRepo{
			listByOverrideFn: func(ctx context.Context, overrideID string) ([]models.DailyNutritionOverrideItemRecord, error) {
				return []models.DailyNutritionOverrideItemRecord{
					{ID: "add-chicken", OverrideID: legacyOverrideID, ProductID: legacyChickenID, AmountGrams: 50, Operation: string(models.OperationAdd), MealLabel: ptrStr("Snack")},
					{ID: "replace-chicken", OverrideID: legacyOverrideID, ProductID: legacyChickenID, AmountGrams: 150, Operation: string(models.OperationReplace), MealLabel: ptrStr("Dinner")},
				}, nil
			},
		},
		legacyProductRepo(t),
	)

	day, err := resolver.Resolve(ctx, testUserID, models.MustDate("2026-06-24"))
	require.NoError(t, err)
	require.NotNil(t, day)
	assert.Equal(t, models.LegacyResolutionResolved, day.Status)
	require.Len(t, day.ResolvedEntries, 2)
	assert.InDelta(t, 200, day.Totals.Calories, 0.001)
	assert.InDelta(t, day.LegacyTotals.Calories, day.Totals.Calories, 0.001)

	var totalChickenGrams float64
	for _, entry := range day.ResolvedEntries {
		if entry.ProductID == legacyChickenID {
			totalChickenGrams += entry.AmountGrams
		}
	}
	assert.InDelta(t, 200, totalChickenGrams, 0.001)
}

func TestDailyNutritionLegacyResolver_ResolvesAddWithoutTemplate(t *testing.T) {
	resolver := newLegacyResolver(
		&mockNutritionTemplateRepo{
			getByWeekFn: func(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error) {
				return nil, nil
			},
		},
		&mockNutritionTemplateItemRepo{},
		&mockNutritionOverrideRepo{
			getByDateFn: func(ctx context.Context, userID string, date string) (*models.DailyNutritionOverrideRecord, error) {
				return legacyOverrideRecord(), nil
			},
		},
		&mockNutritionOverrideItemRepo{
			listByOverrideFn: func(ctx context.Context, overrideID string) ([]models.DailyNutritionOverrideItemRecord, error) {
				return []models.DailyNutritionOverrideItemRecord{
					{ID: "add-yogurt", OverrideID: legacyOverrideID, ProductID: legacyYogurtID, AmountGrams: 100, Operation: string(models.OperationAdd)},
				}, nil
			},
		},
		legacyProductRepo(t),
	)

	day, err := resolver.Resolve(ctx, testUserID, models.MustDate("2026-06-24"))
	require.NoError(t, err)
	require.NotNil(t, day)
	assert.Equal(t, models.LegacyResolutionResolved, day.Status)
	require.Len(t, day.ResolvedEntries, 1)
	assert.Equal(t, legacyYogurtID, day.ResolvedEntries[0].ProductID)
	assert.InDelta(t, 80, day.Totals.Calories, 0.001)
}

func TestDailyNutritionLegacyResolver_MarksUnresolvedSubtractOverBase(t *testing.T) {
	resolver := newLegacyResolver(
		&mockNutritionTemplateRepo{
			getByWeekFn: func(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error) {
				return legacyTemplateRecord(), nil
			},
		},
		&mockNutritionTemplateItemRepo{
			listByTemplateFn: func(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error) {
				return []models.NutritionTemplateItemRecord{
					{ID: "base-rice", TemplateID: legacyTemplateID, ProductID: legacyRiceID, AmountGrams: 50},
				}, nil
			},
		},
		&mockNutritionOverrideRepo{
			getByDateFn: func(ctx context.Context, userID string, date string) (*models.DailyNutritionOverrideRecord, error) {
				return legacyOverrideRecord(), nil
			},
		},
		&mockNutritionOverrideItemRepo{
			listByOverrideFn: func(ctx context.Context, overrideID string) ([]models.DailyNutritionOverrideItemRecord, error) {
				return []models.DailyNutritionOverrideItemRecord{
					{ID: "subtract-rice", OverrideID: legacyOverrideID, ProductID: legacyRiceID, AmountGrams: 75, Operation: string(models.OperationSubtract)},
				}, nil
			},
		},
		legacyProductRepo(t),
	)

	day, err := resolver.Resolve(ctx, testUserID, models.MustDate("2026-06-24"))
	require.NoError(t, err)
	require.NotNil(t, day)
	assert.Equal(t, models.LegacyResolutionUnresolved, day.Status)
	assert.NotEmpty(t, day.RawOperations)
	assert.Contains(t, day.UnresolvedReasons, "subtract exceeds base amount")
}

func TestDailyNutritionLegacyResolver_MarksUnresolvedMissingProduct(t *testing.T) {
	resolver := newLegacyResolver(
		&mockNutritionTemplateRepo{
			getByWeekFn: func(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error) {
				return nil, nil
			},
		},
		&mockNutritionTemplateItemRepo{},
		&mockNutritionOverrideRepo{
			getByDateFn: func(ctx context.Context, userID string, date string) (*models.DailyNutritionOverrideRecord, error) {
				return legacyOverrideRecord(), nil
			},
		},
		&mockNutritionOverrideItemRepo{
			listByOverrideFn: func(ctx context.Context, overrideID string) ([]models.DailyNutritionOverrideItemRecord, error) {
				return []models.DailyNutritionOverrideItemRecord{
					{ID: "add-missing", OverrideID: legacyOverrideID, ProductID: legacyMissingID, AmountGrams: 100, Operation: string(models.OperationAdd)},
				}, nil
			},
		},
		legacyProductRepo(t),
	)

	day, err := resolver.Resolve(ctx, testUserID, models.MustDate("2026-06-24"))
	require.NoError(t, err)
	require.NotNil(t, day)
	assert.Equal(t, models.LegacyResolutionUnresolved, day.Status)
	assert.Contains(t, day.UnresolvedReasons, "missing product context")
}

func TestDailyNutritionLegacyResolver_MarksUnresolvedMultipleConflictingReplacements(t *testing.T) {
	resolver := newLegacyResolver(
		&mockNutritionTemplateRepo{
			getByWeekFn: func(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplateRecord, error) {
				return legacyTemplateRecord(), nil
			},
		},
		&mockNutritionTemplateItemRepo{
			listByTemplateFn: func(ctx context.Context, templateID string) ([]models.NutritionTemplateItemRecord, error) {
				return []models.NutritionTemplateItemRecord{
					{ID: "base-chicken", TemplateID: legacyTemplateID, ProductID: legacyChickenID, AmountGrams: 100},
				}, nil
			},
		},
		&mockNutritionOverrideRepo{
			getByDateFn: func(ctx context.Context, userID string, date string) (*models.DailyNutritionOverrideRecord, error) {
				return legacyOverrideRecord(), nil
			},
		},
		&mockNutritionOverrideItemRepo{
			listByOverrideFn: func(ctx context.Context, overrideID string) ([]models.DailyNutritionOverrideItemRecord, error) {
				return []models.DailyNutritionOverrideItemRecord{
					{ID: "replace-chicken-a", OverrideID: legacyOverrideID, ProductID: legacyChickenID, AmountGrams: 120, Operation: string(models.OperationReplace)},
					{ID: "replace-chicken-b", OverrideID: legacyOverrideID, ProductID: legacyChickenID, AmountGrams: 140, Operation: string(models.OperationReplace)},
				}, nil
			},
		},
		legacyProductRepo(t),
	)

	day, err := resolver.Resolve(ctx, testUserID, models.MustDate("2026-06-24"))
	require.NoError(t, err)
	require.NotNil(t, day)
	assert.Equal(t, models.LegacyResolutionUnresolved, day.Status)
	assert.Contains(t, day.UnresolvedReasons, "multiple conflicting replacements")
	assert.Len(t, day.RawOperations, 2)
}
