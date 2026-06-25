// FILE: apps/api/internal/atlas/service/export_zip_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for AI export ZIP archive shape, manifest/data files, CSV headers, photo inclusion, and user profile serialization.
//   SCOPE: In-memory ZIP generation only; excludes AiExportService repository writes and runtime data-provider assembly.
//   DEPENDS: apps/api/internal/atlas/service export ZIP helpers.
//   LINKS: M-API / V-M-API / M-API-NUTRITION / V-M-API-NUTRITION.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Added exact detailed nutrition CSV header coverage for Task 11 nutrition AI export payloads.
// END_CHANGE_SUMMARY

package service_test

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/atlas/service"
)

func TestAiExportZIP_ValidArchive(t *testing.T) {
	archive := service.NewDefaultExportArchive("2026-01-01", "2026-01-28", service.ExportProfile{})
	data, err := archive.BuildZIP()
	require.NoError(t, err)
	require.NotEmpty(t, data)

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	require.NoError(t, err)
	assert.Len(t, zr.File, 7) // manifest.json, data.json, summary.md, 4 CSV files (workouts, measurements, nutrition, cardio)
}

func TestAiExportZIP_ManifestStructure(t *testing.T) {
	archive := service.NewDefaultExportArchive("2026-01-01", "2026-01-28", service.ExportProfile{})
	archive.Manifest.IncludedSections.Workouts = true
	archive.Manifest.IncludedSections.Nutrition = true

	data, err := archive.BuildZIP()
	require.NoError(t, err)

	content := readFileFromZIP(t, data, "manifest.json")
	var manifest service.Manifest
	err = json.Unmarshal([]byte(content), &manifest)
	require.NoError(t, err)

	assert.Equal(t, 1, manifest.SchemaVersion)
	assert.Equal(t, "2026-01-01", manifest.DateRangeStart)
	assert.Equal(t, "2026-01-28", manifest.DateRangeEnd)
	assert.NotEmpty(t, manifest.ExportTimestamp)
	assert.True(t, manifest.IncludedSections.Workouts)
	assert.True(t, manifest.IncludedSections.Nutrition)
	assert.False(t, manifest.IncludedSections.Photos)
}

func TestAiExportZIP_DataJSONStructure(t *testing.T) {
	archive := service.NewDefaultExportArchive("2026-01-01", "2026-01-28", service.ExportProfile{})
	data, err := archive.BuildZIP()
	require.NoError(t, err)

	content := readFileFromZIP(t, data, "data.json")
	var dataJSON map[string]any
	err = json.Unmarshal([]byte(content), &dataJSON)
	require.NoError(t, err)

	assert.Contains(t, dataJSON, "workouts")
	assert.Contains(t, dataJSON, "cardio")
	assert.Contains(t, dataJSON, "bodyWeightEntries")
	assert.Contains(t, dataJSON, "measurements")
	assert.Contains(t, dataJSON, "nutrition")
	assert.Contains(t, dataJSON, "weekFlags")
	assert.Contains(t, dataJSON, "userProfile")
}

func TestAiExportZIP_SummaryMDContent(t *testing.T) {
	archive := service.NewDefaultExportArchive("2026-01-01", "2026-01-28", service.ExportProfile{})
	archive.SummaryMD = "# AI Export Summary\nPeriod: 2026-01-01 to 2026-01-28"

	data, err := archive.BuildZIP()
	require.NoError(t, err)

	content := readFileFromZIP(t, data, "summary.md")
	assert.Contains(t, content, "AI Export Summary")
	assert.Contains(t, content, "2026-01-01")
	assert.Contains(t, content, "2026-01-28")
}

func TestAiExportZIP_CSVFilesExist(t *testing.T) {
	archive := service.NewDefaultExportArchive("2026-01-01", "2026-01-28", service.ExportProfile{})
	data, err := archive.BuildZIP()
	require.NoError(t, err)

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	require.NoError(t, err)

	names := make([]string, len(zr.File))
	for i, f := range zr.File {
		names[i] = f.Name
	}

	assert.Contains(t, names, "workouts.csv")
	assert.Contains(t, names, "measurements.csv")
	assert.Contains(t, names, "nutrition.csv")
	assert.Contains(t, names, "cardio.csv")
}

func TestAiExportZIP_CSVHeadersAndRows(t *testing.T) {
	archive := service.NewDefaultExportArchive("2026-01-01", "2026-01-28", service.ExportProfile{})
	archive.WorkoutsCSV = service.CSVData{
		Headers: []string{"date", "exercise_name", "set_number", "weight", "reps", "rpe", "rir", "set_notes", "exercise_notes", "day_notes"},
		Rows: [][]string{
			{"2026-01-02", "Bench Press", "1", "80", "10", "8", "2", "", "", "Push Day"},
		},
	}

	data, err := archive.BuildZIP()
	require.NoError(t, err)

	content := readFileFromZIP(t, data, "workouts.csv")
	assert.Contains(t, content, "date,exercise_name")
	assert.Contains(t, content, "Bench Press")
	assert.Contains(t, content, "Push Day")
}

func TestExportZip_NutritionCSVHeadersExact(t *testing.T) {
	archive := service.NewDefaultExportArchive("2026-01-01", "2026-01-28", service.ExportProfile{})

	data, err := archive.BuildZIP()
	require.NoError(t, err)

	content := readFileFromZIP(t, data, "nutrition.csv")
	reader := csv.NewReader(strings.NewReader(content))
	rows, err := reader.ReadAll()
	require.NoError(t, err)
	require.NotEmpty(t, rows)
	assert.Equal(t, []string{
		"date",
		"product_id",
		"product_name_snapshot",
		"amount_grams",
		"calories_per_100g_snapshot",
		"protein_per_100g_snapshot",
		"fat_per_100g_snapshot",
		"carbs_per_100g_snapshot",
		"entry_calories",
		"entry_protein",
		"entry_fat",
		"entry_carbs",
		"meal_label",
		"notes",
	}, rows[0])
}

func TestAiExportZIP_PhotosIncluded(t *testing.T) {
	archive := service.NewDefaultExportArchive("2026-01-01", "2026-01-28", service.ExportProfile{})
	archive.Photos = []service.ExportPhoto{
		{CheckInID: "checkin-1", Angle: "FRONT", Extension: "jpg", Data: []byte("fake-image-data")},
	}

	data, err := archive.BuildZIP()
	require.NoError(t, err)

	content := readFileFromZIP(t, data, "photos/checkin-1_FRONT.jpg")
	assert.Equal(t, "fake-image-data", content)
}

func TestAiExportZIP_PhotosExcludedByDefault(t *testing.T) {
	archive := service.NewDefaultExportArchive("2026-01-01", "2026-01-28", service.ExportProfile{})
	data, err := archive.BuildZIP()
	require.NoError(t, err)

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	require.NoError(t, err)

	for _, f := range zr.File {
		assert.NotContains(t, f.Name, "photos/", "photos/ should not exist in default export")
	}
}

func TestAiExportZIP_NoPhotosDirWhenOptedOut(t *testing.T) {
	archive := service.NewDefaultExportArchive("2026-01-01", "2026-01-28", service.ExportProfile{})
	data, err := archive.BuildZIP()
	require.NoError(t, err)

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	require.NoError(t, err)

	for _, f := range zr.File {
		assert.False(t, strings.HasPrefix(f.Name, "photos/"), "photos/ directory should not exist")
	}
}

func TestAiExportZIP_WorkoutData(t *testing.T) {
	archive := service.NewDefaultExportArchive("2026-01-01", "2026-01-28", service.ExportProfile{})
	archive.Data.Workouts = []any{
		map[string]any{"date": "2026-01-02", "exercise": "Bench Press", "sets": 3},
	}

	data, err := archive.BuildZIP()
	require.NoError(t, err)

	content := readFileFromZIP(t, data, "data.json")
	var dataJSON map[string]any
	err = json.Unmarshal([]byte(content), &dataJSON)
	require.NoError(t, err)

	workouts := dataJSON["workouts"].([]any)
	assert.Len(t, workouts, 1)
}

func TestAiExportZIP_UserProfileInManifest(t *testing.T) {
	goal := "Lose weight"
	height := 175.0
	archive := service.NewDefaultExportArchive("2026-01-01", "2026-01-28", service.ExportProfile{
		Goal:   &goal,
		Height: &height,
	})

	data, err := archive.BuildZIP()
	require.NoError(t, err)

	content := readFileFromZIP(t, data, "data.json")
	var dataJSON map[string]any
	err = json.Unmarshal([]byte(content), &dataJSON)
	require.NoError(t, err)

	profile := dataJSON["userProfile"].(map[string]any)
	assert.Equal(t, "Lose weight", profile["goal"])
	assert.Equal(t, 175.0, profile["height"])
}

func readFileFromZIP(t *testing.T, zipData []byte, fileName string) string {
	t.Helper()
	zr, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	require.NoError(t, err)

	for _, f := range zr.File {
		if f.Name == fileName {
			rc, err := f.Open()
			require.NoError(t, err)
			defer rc.Close()

			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(rc)
			require.NoError(t, err)
			return buf.String()
		}
	}
	t.Fatalf("file %s not found in zip", fileName)
	return ""
}
