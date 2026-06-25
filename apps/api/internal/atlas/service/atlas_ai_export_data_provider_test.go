// FILE: apps/api/internal/atlas/service/atlas_ai_export_data_provider_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for the Atlas repo/service-backed AI export data provider.
//   SCOPE: Daily nutrition factual logs, weekly template planned entries with product snapshots, and unresolved legacy nutrition export payloads.
//   DEPENDS: apps/api/internal/atlas/service, apps/api/internal/atlas/models, atlas postgres repository mocks.
//   LINKS: M-API-NUTRITION / V-M-API-NUTRITION.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Task 11 RED coverage for service-backed nutrition AI export data provider.
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

type atlasAiExportDailyNutritionLogServiceMock struct {
	service.DailyNutritionLogService
	listByRangeFn func(ctx context.Context, userID string, from, to models.Date) ([]models.DailyNutritionLog, error)
}

func (m *atlasAiExportDailyNutritionLogServiceMock) ListByRange(ctx context.Context, userID string, from, to models.Date) ([]models.DailyNutritionLog, error) {
	return m.listByRangeFn(ctx, userID, from, to)
}

func TestAtlasAiExportDataProvider_ReturnsDailyTemplatesAndLegacy(t *testing.T) {
	from := models.MustDate("2026-06-24")
	to := models.MustDate("2026-06-24")
	chickenID := "a10e8400-e29b-41d4-a716-446655440000"
	riceID := "a20e8400-e29b-41d4-a716-446655440000"
	templateID := "a30e8400-e29b-41d4-a716-446655440000"
	overrideID := "a40e8400-e29b-41d4-a716-446655440000"

	dailyEntry := models.DailyNutritionEntry{
		ID:                      "entry-1",
		DailyLogID:              "log-1",
		ProductID:               chickenID,
		ProductNameSnapshot:     "Chicken Breast",
		AmountGrams:             150,
		CaloriesPer100gSnapshot: 165,
		ProteinPer100gSnapshot:  31,
		FatPer100gSnapshot:      3.6,
		CarbsPer100gSnapshot:    0,
		MealLabel:               ptrStr("Lunch"),
		Notes:                   ptrStr("grilled"),
	}
	dailyEntry.Macros = models.DailyNutritionEntryMacros(dailyEntry)

	provider := service.NewAtlasAiExportDataProvider(
		&atlasAiExportDailyNutritionLogServiceMock{
			listByRangeFn: func(ctx context.Context, userID string, fromDate, toDate models.Date) ([]models.DailyNutritionLog, error) {
				assert.Equal(t, testUserID, userID)
				assert.Equal(t, from, fromDate)
				assert.Equal(t, to, toDate)
				return []models.DailyNutritionLog{
					{
						ID:      "log-1",
						UserID:  testUserID,
						Date:    "2026-06-24",
						Notes:   ptrStr("training day"),
						Entries: []models.DailyNutritionEntry{dailyEntry},
						Totals:  models.DailyNutritionTotalsFromEntries([]models.DailyNutritionEntry{dailyEntry}),
					},
				}, nil
			},
		},
		&mockNutritionTemplateRepo{
			listByRangeFn: func(ctx context.Context, userID string, startDate, endDate string) ([]models.NutritionTemplateRecord, error) {
				assert.Equal(t, "2026-06-24", startDate)
				assert.Equal(t, "2026-06-24", endDate)
				return []models.NutritionTemplateRecord{
					{
						ID:            templateID,
						UserID:        testUserID,
						WeekStartDate: models.MustDate("2026-06-22"),
						Title:         ptrStr("Week A"),
						Notes:         ptrStr("planned high protein week"),
					},
				}, nil
			},
		},
		&mockNutritionTemplateItemRepo{
			listByTemplateFn: func(ctx context.Context, requestedTemplateID string) ([]models.NutritionTemplateItemRecord, error) {
				assert.Equal(t, templateID, requestedTemplateID)
				return []models.NutritionTemplateItemRecord{
					{
						ID:          "item-1",
						TemplateID:  templateID,
						ProductID:   riceID,
						AmountGrams: 200,
						MealLabel:   ptrStr("Dinner"),
						Notes:       ptrStr("post workout"),
					},
				}, nil
			},
		},
		&mockNutritionProductRepo{
			getByIDIncludeInactiveFn: func(ctx context.Context, userID string, productID string) (*models.NutritionProductRecord, error) {
				assert.Equal(t, testUserID, userID)
				assert.Equal(t, riceID, productID)
				return &models.NutritionProductRecord{
					ID:              riceID,
					UserID:          userID,
					Name:            "Rice",
					CaloriesPer100g: 130,
					ProteinPer100g:  2.7,
					FatPer100g:      0.3,
					CarbsPer100g:    28,
					IsActive:        true,
				}, nil
			},
		},
		&mockDailyNutritionLegacyResolver{
			resolveFn: func(ctx context.Context, userID string, date models.Date) (*models.DailyNutritionLegacyResolution, error) {
				assert.Equal(t, testUserID, userID)
				assert.Equal(t, from, date)
				return &models.DailyNutritionLegacyResolution{
					Status:           models.LegacyResolutionUnresolved,
					Date:             date.String(),
					WeekStartDate:    "2026-06-22",
					SourceOverrideID: overrideID,
					LegacyTotals: models.NutritionMacros{
						Calories: 500,
						Protein:  40,
						Fat:      10,
						Carbs:    60,
					},
					RawOperations: []models.DailyNutritionLegacyOperation{
						{ID: "raw-1", OverrideID: overrideID, ProductID: riceID, AmountGrams: 200, Operation: "replace"},
					},
					UnresolvedReasons: []string{"ambiguous legacy operation"},
				}, nil
			},
		},
	)

	daily, err := provider.GetDailyNutritionExport(ctx, testUserID, from, to)
	require.NoError(t, err)
	require.Len(t, daily, 1)
	dailyLog := requireMap(t, daily[0])
	assert.Equal(t, "2026-06-24", dailyLog["date"])
	dailyEntries := requireSlice(t, dailyLog["entries"])
	require.Len(t, dailyEntries, 1)
	assert.Equal(t, "Chicken Breast", requireMap(t, dailyEntries[0])["productNameSnapshot"])
	assert.Equal(t, 247.5, requireMap(t, dailyEntries[0])["entryCalories"])

	templates, err := provider.GetNutritionTemplateExport(ctx, testUserID, from, to)
	require.NoError(t, err)
	require.Len(t, templates, 1)
	template := requireMap(t, templates[0])
	assert.Equal(t, "2026-06-22", template["weekStartDate"])
	assert.Equal(t, "2026-06-28", template["weekEndDate"])
	plannedEntries := requireSlice(t, template["plannedEntries"])
	require.Len(t, plannedEntries, 1)
	assert.Equal(t, "Rice", requireMap(t, plannedEntries[0])["productNameSnapshot"])
	assert.Equal(t, 260.0, requireMap(t, plannedEntries[0])["entryCalories"])

	legacy, err := provider.GetLegacyNutritionExport(ctx, testUserID, from, to)
	require.NoError(t, err)
	require.Len(t, legacy, 1)
	legacyDay := requireMap(t, legacy[0])
	assert.Equal(t, "unresolved", legacyDay["legacyResolutionStatus"])
	assert.Equal(t, "replace", requireMap(t, requireSlice(t, legacyDay["rawOperations"])[0])["operation"])
}
