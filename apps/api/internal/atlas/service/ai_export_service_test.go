// FILE: apps/api/internal/atlas/service/ai_export_service_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for BuildPrompt covering user context, persistent AI context, one-time comment, week flags, empty date range, no-data-in-period, all profile fields, and nil profile.
//   SCOPE: Pure function tests for prompt generation logic. Does not cover Generate/GetByID/List/Delete service methods (those require repo mocks and are covered by integration tests).
//   DEPENDS: apps/api/internal/atlas/service (BuildPrompt, UserProfileExport, SectionToggles).
//   LINKS: M-API / V-M-API / WAVE-07.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added BuildPrompt pure function tests for WAVE-07.
// END_CHANGE_SUMMARY

package service_test

import (
	"archive/zip"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasPostgres "monorepo-template/apps/api/internal/atlas/repository/postgres"
	"monorepo-template/apps/api/internal/atlas/service"
)

func TestBuildPrompt_IncludesUserContext(t *testing.T) {
	goal := "Build muscle"
	profile := &service.UserProfileExport{Goal: &goal}
	prompt := service.BuildPrompt(profile, "2026-01-01", "2026-01-28", service.SectionToggles{}, nil, nil, "Some data")
	assert.Contains(t, prompt, "Build muscle")
	assert.Contains(t, prompt, "2026-01-01")
	assert.Contains(t, prompt, "2026-01-28")
	assert.Contains(t, prompt, "Analysis Requests")
}

func TestBuildPrompt_WithPersistentContext(t *testing.T) {
	ctx := "Focus on progressive overload and recovery"
	profile := &service.UserProfileExport{PersistentAiContext: &ctx}
	prompt := service.BuildPrompt(profile, "2026-01-01", "2026-01-28", service.SectionToggles{}, nil, nil, "")
	assert.Contains(t, prompt, "progressive overload")
}

func TestBuildPrompt_WithOneTimeComment(t *testing.T) {
	comment := "Trying a new deload protocol this month"
	profile := &service.UserProfileExport{}
	prompt := service.BuildPrompt(profile, "2026-01-01", "2026-01-28", service.SectionToggles{}, &comment, nil, "")
	assert.Contains(t, prompt, "deload protocol")
}

func TestBuildPrompt_WithWeekFlags(t *testing.T) {
	profile := &service.UserProfileExport{}
	flags := []string{"POOR_SLEEP", "HIGH_STRESS"}
	prompt := service.BuildPrompt(profile, "2026-01-01", "2026-01-28", service.SectionToggles{}, nil, flags, "")
	assert.Contains(t, prompt, "POOR_SLEEP")
	assert.Contains(t, prompt, "HIGH_STRESS")
}

func TestBuildPrompt_EmptyDateRange(t *testing.T) {
	profile := &service.UserProfileExport{}
	prompt := service.BuildPrompt(profile, "", "", service.SectionToggles{}, nil, nil, "")
	assert.NotEmpty(t, prompt)
	assert.Contains(t, prompt, "Analysis Requests")
}

func TestBuildPrompt_NoDataInPeriod(t *testing.T) {
	profile := &service.UserProfileExport{}
	prompt := service.BuildPrompt(profile, "2026-01-01", "2026-01-28", service.SectionToggles{}, nil, nil, "")
	assert.NotEmpty(t, prompt)
	assert.Contains(t, prompt, "Analysis Requests")
}

func TestBuildPrompt_WithAllProfileFields(t *testing.T) {
	goal := "Lose weight"
	height := 175.0
	birthDate := "1990-06-15"
	exp := "Advanced"
	split := "Upper/Lower"
	prog := "Linear"
	nutri := "Low Carb"
	aiCtx := "Focus on Zone 2 cardio"

	profile := &service.UserProfileExport{
		Goal:                      &goal,
		Height:                    &height,
		BirthDate:                 &birthDate,
		TrainingExperience:        &exp,
		CurrentTrainingSplit:      &split,
		PreferredProgressionStyle: &prog,
		NutritionStrategy:         &nutri,
		PersistentAiContext:       &aiCtx,
	}
	prompt := service.BuildPrompt(profile, "2026-01-01", "2026-01-28", service.SectionToggles{}, nil, nil, "")
	assert.Contains(t, prompt, "Lose weight")
	assert.Contains(t, prompt, "175.0")
	assert.Contains(t, prompt, "Advanced")
	assert.Contains(t, prompt, "Upper/Lower")
	assert.Contains(t, prompt, "Low Carb")
	assert.Contains(t, prompt, "Zone 2 cardio")
}

func TestBuildPrompt_NilProfile(t *testing.T) {
	prompt := service.BuildPrompt(nil, "2026-01-01", "2026-01-28", service.SectionToggles{}, nil, nil, "")
	assert.NotEmpty(t, prompt)
	assert.Contains(t, prompt, "Analysis Requests")
}

type aiExportRepoMock struct {
	atlasPostgres.AiExportRepository
	createFn         func(ctx context.Context, userID, dateRangeStart, dateRangeEnd string, includePhotos, includeNutrition, includeCardio, includeMeasurements bool, userComment *string, generatedPrompt string) (*models.AiExportRecord, error)
	updateFilePathFn func(ctx context.Context, id string, filePath *string) (*models.AiExportRecord, error)
}

func (m *aiExportRepoMock) Create(ctx context.Context, userID, dateRangeStart, dateRangeEnd string, includePhotos, includeNutrition, includeCardio, includeMeasurements bool, userComment *string, generatedPrompt string) (*models.AiExportRecord, error) {
	return m.createFn(ctx, userID, dateRangeStart, dateRangeEnd, includePhotos, includeNutrition, includeCardio, includeMeasurements, userComment, generatedPrompt)
}

func (m *aiExportRepoMock) UpdateFilePath(ctx context.Context, id string, filePath *string) (*models.AiExportRecord, error) {
	return m.updateFilePathFn(ctx, id, filePath)
}

type aiExportProfileRepoMock struct {
	atlasPostgres.UserProfileRepository
	findByUserIDFn func(ctx context.Context, userID string) (*models.UserProfileRecord, error)
}

func (m *aiExportProfileRepoMock) FindByUserID(ctx context.Context, userID string) (*models.UserProfileRecord, error) {
	return m.findByUserIDFn(ctx, userID)
}

type aiExportProviderMock struct {
	workouts          []any
	cardio            []any
	bodyWeight        []any
	bodyCheckIns      []any
	measurements      []any
	weekFlags         []any
	dailyNutrition    []any
	templateNutrition []any
	legacyNutrition   []any
	photos            []service.ExportPhoto
}

func (m *aiExportProviderMock) GetWorkoutSummary(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return m.workouts, nil
}

func (m *aiExportProviderMock) GetCardioEntries(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return m.cardio, nil
}

func (m *aiExportProviderMock) GetBodyWeightEntries(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return m.bodyWeight, nil
}

func (m *aiExportProviderMock) GetBodyCheckIns(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return m.bodyCheckIns, nil
}

func (m *aiExportProviderMock) GetBodyMeasurements(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return m.measurements, nil
}

func (m *aiExportProviderMock) GetWeekFlags(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return m.weekFlags, nil
}

func (m *aiExportProviderMock) GetDailyNutritionExport(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return m.dailyNutrition, nil
}

func (m *aiExportProviderMock) GetNutritionTemplateExport(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return m.templateNutrition, nil
}

func (m *aiExportProviderMock) GetLegacyNutritionExport(ctx context.Context, userID string, from, to models.Date) ([]any, error) {
	return m.legacyNutrition, nil
}

func (m *aiExportProviderMock) GetProgressPhotos(ctx context.Context, userID string, from, to models.Date) ([]service.ExportPhoto, error) {
	return m.photos, nil
}

func TestAiExportService_IncludesDailyNutritionEntriesWithoutExternalCall(t *testing.T) {
	provider := &aiExportProviderMock{
		dailyNutrition: []any{
			map[string]any{
				"date": "2026-06-24",
				"totals": map[string]any{
					"calories": 345.0,
					"protein":  32.5,
					"fat":      7.5,
					"carbs":    30.0,
				},
				"entries": []any{
					map[string]any{
						"productId":               "product-chicken",
						"productNameSnapshot":     "Chicken Breast",
						"amountGrams":             150.0,
						"caloriesPer100gSnapshot": 165.0,
						"proteinPer100gSnapshot":  31.0,
						"fatPer100gSnapshot":      3.6,
						"carbsPer100gSnapshot":    0.0,
						"entryCalories":           247.5,
						"entryProtein":            46.5,
						"entryFat":                5.4,
						"entryCarbs":              0.0,
						"mealLabel":               "Lunch",
						"notes":                   "grilled",
					},
				},
			},
		},
	}
	zipPath := generateAiExportForTest(t, provider)

	data := readDataJSONFromZipPath(t, zipPath)
	nutrition := requireMap(t, data["nutrition"])
	dailyLogs := requireSlice(t, nutrition["dailyLogs"])
	require.Len(t, dailyLogs, 1)
	day := requireMap(t, dailyLogs[0])
	assert.Equal(t, "2026-06-24", day["date"])
	entries := requireSlice(t, day["entries"])
	require.Len(t, entries, 1)
	entry := requireMap(t, entries[0])
	assert.Equal(t, "product-chicken", entry["productId"])
	assert.Equal(t, "Chicken Breast", entry["productNameSnapshot"])
	assert.Equal(t, 150.0, entry["amountGrams"])
	assert.Equal(t, 165.0, entry["caloriesPer100gSnapshot"])
	assert.Equal(t, 31.0, entry["proteinPer100gSnapshot"])
	assert.Equal(t, 3.6, entry["fatPer100gSnapshot"])
	assert.Equal(t, 0.0, entry["carbsPer100gSnapshot"])
	assert.Equal(t, 247.5, entry["entryCalories"])
	assert.Equal(t, 46.5, entry["entryProtein"])
	assert.Equal(t, 5.4, entry["entryFat"])
	assert.Equal(t, 0.0, entry["entryCarbs"])
	assert.Equal(t, "Lunch", entry["mealLabel"])
	assert.Equal(t, "grilled", entry["notes"])
}

func TestAiExportService_IncludesWeeklyNutritionPlanContext(t *testing.T) {
	provider := &aiExportProviderMock{
		templateNutrition: []any{
			map[string]any{
				"weekStartDate": "2026-06-22",
				"weekEndDate":   "2026-06-28",
				"notes":         "planned high protein week",
				"plannedEntries": []any{
					map[string]any{
						"productId":               "product-rice",
						"productNameSnapshot":     "Rice",
						"amountGrams":             200.0,
						"mealLabel":               "Dinner",
						"notes":                   "post workout",
						"caloriesPer100gSnapshot": 130.0,
						"proteinPer100gSnapshot":  2.7,
						"fatPer100gSnapshot":      0.3,
						"carbsPer100gSnapshot":    28.0,
						"entryCalories":           260.0,
						"entryProtein":            5.4,
						"entryFat":                0.6,
						"entryCarbs":              56.0,
					},
				},
			},
		},
	}
	zipPath := generateAiExportForTest(t, provider)

	data := readDataJSONFromZipPath(t, zipPath)
	nutrition := requireMap(t, data["nutrition"])
	templates := requireSlice(t, nutrition["templates"])
	require.Len(t, templates, 1)
	template := requireMap(t, templates[0])
	assert.Equal(t, "2026-06-22", template["weekStartDate"])
	assert.Equal(t, "2026-06-28", template["weekEndDate"])
	assert.Equal(t, "planned high protein week", template["notes"])
	plannedEntries := requireSlice(t, template["plannedEntries"])
	require.Len(t, plannedEntries, 1)
	entry := requireMap(t, plannedEntries[0])
	assert.Equal(t, "product-rice", entry["productId"])
	assert.Equal(t, "Rice", entry["productNameSnapshot"])
	assert.Equal(t, 200.0, entry["amountGrams"])
	assert.Equal(t, "Dinner", entry["mealLabel"])
	assert.Equal(t, "post workout", entry["notes"])
	assert.Equal(t, 130.0, entry["caloriesPer100gSnapshot"])
	assert.Equal(t, 260.0, entry["entryCalories"])
}

func TestAiExportService_IncludesUnresolvedLegacyNutritionBlock(t *testing.T) {
	provider := &aiExportProviderMock{
		legacyNutrition: []any{
			map[string]any{
				"date":                   "2026-06-24",
				"legacyResolutionStatus": "unresolved",
				"legacyTotals": map[string]any{
					"calories": 500.0,
					"protein":  40.0,
					"fat":      10.0,
					"carbs":    60.0,
				},
				"rawOperations": []any{
					map[string]any{
						"productId":   "product-rice",
						"amountGrams": 100.0,
						"operation":   "add",
						"mealLabel":   "Lunch",
						"notes":       "legacy add",
					},
					map[string]any{"productId": "product-rice", "amountGrams": 50.0, "operation": "subtract"},
					map[string]any{"productId": "product-rice", "amountGrams": 200.0, "operation": "replace"},
				},
			},
		},
	}
	zipPath := generateAiExportForTest(t, provider)

	data := readDataJSONFromZipPath(t, zipPath)
	nutrition := requireMap(t, data["nutrition"])
	legacy := requireSlice(t, nutrition["legacy"])
	require.Len(t, legacy, 1)
	day := requireMap(t, legacy[0])
	assert.Equal(t, "unresolved", day["legacyResolutionStatus"])
	totals := requireMap(t, day["legacyTotals"])
	assert.Equal(t, 500.0, totals["calories"])
	raw := requireSlice(t, day["rawOperations"])
	require.Len(t, raw, 3)
	assert.Equal(t, "add", requireMap(t, raw[0])["operation"])
	assert.Equal(t, "subtract", requireMap(t, raw[1])["operation"])
	assert.Equal(t, "replace", requireMap(t, raw[2])["operation"])
}

func TestAiExportService_WritesPrivateNutritionZipPermissions(t *testing.T) {
	provider := &aiExportProviderMock{
		dailyNutrition: []any{
			map[string]any{"date": "2026-06-24", "entries": []any{}},
		},
	}
	zipPath := generateAiExportForTest(t, provider)

	exportDirInfo, err := os.Stat(filepath.Dir(zipPath))
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0700), exportDirInfo.Mode().Perm())

	zipInfo, err := os.Stat(zipPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), zipInfo.Mode().Perm())
}

func generateAiExportForTest(t *testing.T, provider *aiExportProviderMock) string {
	t.Helper()

	exportID := "770e8400-e29b-41d4-a716-446655440000"
	exportBasePath := t.TempDir()
	exportRepo := &aiExportRepoMock{
		createFn: func(ctx context.Context, userID, dateRangeStart, dateRangeEnd string, includePhotos, includeNutrition, includeCardio, includeMeasurements bool, userComment *string, generatedPrompt string) (*models.AiExportRecord, error) {
			return &models.AiExportRecord{
				ID:                  exportID,
				UserID:              userID,
				DateRangeStart:      models.MustDate(dateRangeStart),
				DateRangeEnd:        models.MustDate(dateRangeEnd),
				IncludePhotos:       includePhotos,
				IncludeNutrition:    includeNutrition,
				IncludeCardio:       includeCardio,
				IncludeMeasurements: includeMeasurements,
				UserComment:         userComment,
				GeneratedPrompt:     generatedPrompt,
				CreatedAt:           "2026-06-24T00:00:00Z",
				UpdatedAt:           "2026-06-24T00:00:00Z",
			}, nil
		},
		updateFilePathFn: func(ctx context.Context, id string, filePath *string) (*models.AiExportRecord, error) {
			return &models.AiExportRecord{
				ID:                  id,
				UserID:              testUserID,
				DateRangeStart:      models.MustDate("2026-06-24"),
				DateRangeEnd:        models.MustDate("2026-06-30"),
				IncludeNutrition:    true,
				IncludeCardio:       true,
				IncludeMeasurements: true,
				GeneratedPrompt:     "prompt",
				ExportFilePath:      filePath,
				CreatedAt:           "2026-06-24T00:00:00Z",
				UpdatedAt:           "2026-06-24T00:00:00Z",
			}, nil
		},
	}
	profileRepo := &aiExportProfileRepoMock{
		findByUserIDFn: func(ctx context.Context, userID string) (*models.UserProfileRecord, error) {
			return nil, nil
		},
	}
	svc := service.NewAiExportService(exportRepo, profileRepo, provider, zap.NewNop())

	_, _, err := svc.Generate(ctx, testUserID, models.CreateAiExportInput{
		DateRangeStart:      models.MustDate("2026-06-24"),
		DateRangeEnd:        models.MustDate("2026-06-30"),
		IncludeNutrition:    ptrBool(true),
		IncludeCardio:       ptrBool(false),
		IncludeMeasurements: ptrBool(false),
		IncludePhotos:       ptrBool(false),
	}, 365, 10*1024*1024, exportBasePath)
	require.NoError(t, err)

	return filepath.Join(exportBasePath, testUserID, exportID+".zip")
}

func readDataJSONFromZipPath(t *testing.T, zipPath string) map[string]any {
	t.Helper()

	zr, err := zip.OpenReader(zipPath)
	require.NoError(t, err)
	defer zr.Close()

	for _, file := range zr.File {
		if file.Name != "data.json" {
			continue
		}
		rc, err := file.Open()
		require.NoError(t, err)
		defer rc.Close()
		var data map[string]any
		require.NoError(t, json.NewDecoder(rc).Decode(&data))
		return data
	}

	t.Fatalf("data.json not found in %s", zipPath)
	return nil
}

func requireMap(t *testing.T, value any) map[string]any {
	t.Helper()
	result, ok := value.(map[string]any)
	require.Truef(t, ok, "expected map[string]any, got %T", value)
	return result
}

func requireSlice(t *testing.T, value any) []any {
	t.Helper()
	result, ok := value.([]any)
	require.Truef(t, ok, "expected []any, got %T", value)
	return result
}

func ptrBool(v bool) *bool {
	return &v
}
