// FILE: apps/api/internal/atlas/service/ai_export_service.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement AiExportService for local/internal AI export generation, retrieval, and deletion.
//   SCOPE: Generate AI export archives with data summaries, detailed nutrition payloads, private ZIP bundles, list user exports, delete exports.
//   DEPENDS: apps/api/internal/atlas/repository/postgres, apps/api/internal/atlas/models, apps/api/internal/atlas/service/export_zip.go, libs/go/logger.
//   LINKS: M-API / V-M-API / WAVE-07.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AiExportService - Interface for AI export operations.
//   NewAiExportService - Creates a new AiExportService.
//   Generate - Generates an AI export with archive and prompt.
//   GetByID - Gets an AI export by ID (user-scoped).
//   List - Lists AI exports for a user.
//   Delete - Deletes an AI export.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Added detailed nutrition daily/template/legacy export payloads and private archive permissions.
//   LAST_CHANGE: 1.0.0 - Added AI export service for WAVE-07.
// END_CHANGE_SUMMARY

package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasRepo "monorepo-template/apps/api/internal/atlas/repository/postgres"
	"monorepo-template/libs/go/logger"
)

var (
	ErrAiExportInvalidDateRange   = errors.New("date range end must be after start")
	ErrAiExportDateRangeOverLimit = errors.New("date range exceeds maximum of 365 days")
	ErrAiExportNotFound           = errors.New("ai export not found")
	ErrAiExportSizeLimit          = errors.New("export data size exceeds maximum allowed")
	ErrAiExportGenerationFailed   = errors.New("export generation failed")
)

type AiExportService interface {
	Generate(ctx context.Context, userID string, input models.CreateAiExportInput, maxRangeDays int, maxExportSize int64, exportBasePath string) (*models.AiExport, string, error)
	GetByID(ctx context.Context, userID string, id string) (*models.AiExport, error)
	List(ctx context.Context, userID string) ([]models.AiExport, error)
	Delete(ctx context.Context, userID string, id string) (*models.AiExport, error)
}

type aiExportService struct {
	exportRepo   atlasRepo.AiExportRepository
	profileRepo  atlasRepo.UserProfileRepository
	dataProvider AiExportDataProvider
	logger       *zap.Logger
}

func NewAiExportService(exportRepo atlasRepo.AiExportRepository, profileRepo atlasRepo.UserProfileRepository, dataProvider AiExportDataProvider, logger *zap.Logger) AiExportService {
	return &aiExportService{
		exportRepo:   exportRepo,
		profileRepo:  profileRepo,
		dataProvider: dataProvider,
		logger:       logger,
	}
}

func (s *aiExportService) Generate(ctx context.Context, userID string, input models.CreateAiExportInput, maxRangeDays int, maxExportSize int64, exportBasePath string) (*models.AiExport, string, error) {
	log := logger.FromContext(ctx)
	if log == nil {
		log = s.logger
	}

	log.Info("[AiExport][generate][BLOCK_EXPORT_START] generating export")

	if !input.DateRangeEnd.Time().IsZero() && !input.DateRangeStart.Time().IsZero() {
		if input.DateRangeEnd.Time().Before(input.DateRangeStart.Time()) {
			log.Warn("[AiExport][generate][BLOCK_EXPORT_FAILURE] invalid date range")
			return nil, "", ErrAiExportInvalidDateRange
		}
		days := int(input.DateRangeEnd.Time().Sub(input.DateRangeStart.Time()).Hours() / 24)
		if days > maxRangeDays {
			log.Warn("[AiExport][generate][BLOCK_EXPORT_FAILURE] date range over limit")
			return nil, "", ErrAiExportDateRangeOverLimit
		}
	}

	includePhotos := false
	includeNutrition := true
	includeCardio := true
	includeMeasurements := true
	if input.IncludePhotos != nil {
		includePhotos = *input.IncludePhotos
	}
	if input.IncludeNutrition != nil {
		includeNutrition = *input.IncludeNutrition
	}
	if input.IncludeCardio != nil {
		includeCardio = *input.IncludeCardio
	}
	if input.IncludeMeasurements != nil {
		includeMeasurements = *input.IncludeMeasurements
	}

	profileRecord, _ := s.profileRepo.FindByUserID(ctx, userID)
	profileExport := &UserProfileExport{}
	if profileRecord != nil {
		profileExport = profileRecordToExport(profileRecord)
	}

	toggles := SectionToggles{
		Photos:       includePhotos,
		Nutrition:    includeNutrition,
		Cardio:       includeCardio,
		Measurements: includeMeasurements,
	}

	dataSummary := s.buildDataSummary(ctx, userID, input.DateRangeStart, input.DateRangeEnd, toggles)

	log.Info("[AiExport][generate][BLOCK_EXPORT_PROMPT_GENERATE] building prompt")

	dateRangeStart := input.DateRangeStart.String()
	dateRangeEnd := input.DateRangeEnd.String()

	weekFlags := []string{}
	weekFlagData, _ := s.dataProvider.GetWeekFlags(ctx, userID, input.DateRangeStart, input.DateRangeEnd)
	for _, f := range weekFlagData {
		if m, ok := f.(map[string]any); ok {
			if ft, ok := m["flagType"]; ok {
				weekFlags = append(weekFlags, fmt.Sprintf("%v", ft))
			}
		}
	}

	prompt := BuildPrompt(profileExport, dateRangeStart, dateRangeEnd, toggles, input.UserComment, weekFlags, dataSummary)

	log.Info("[AiExport][generate][BLOCK_EXPORT_DATA_QUERY] querying data sources")

	dateRangeStartStr := input.DateRangeStart.String()
	dateRangeEndStr := input.DateRangeEnd.String()
	aiExportRecord, err := s.exportRepo.Create(ctx, userID, dateRangeStartStr, dateRangeEndStr, includePhotos, includeNutrition, includeCardio, includeMeasurements, input.UserComment, prompt)
	if err != nil {
		log.Error("[AiExport][generate][BLOCK_EXPORT_FAILURE] failed to create export record", zap.Error(err))
		return nil, "", fmt.Errorf("ai_export_service.Generate: create record: %w", err)
	}

	log.Info("[AiExport][generate][BLOCK_EXPORT_ZIP_BUILD] building ZIP", zap.String("export_id", aiExportRecord.ID))

	archive, err := s.buildArchive(ctx, userID, dateRangeStart, dateRangeEnd, toggles, profileExport, prompt)
	if err != nil {
		log.Error("[AiExport][generate][BLOCK_EXPORT_FAILURE] failed to build archive", zap.Error(err))
		return nil, "", fmt.Errorf("ai_export_service.Generate: build archive: %w", err)
	}

	zipData, err := archive.BuildZIP()
	if err != nil {
		log.Error("[AiExport][generate][BLOCK_EXPORT_FAILURE] failed to build ZIP", zap.Error(err))
		return nil, "", fmt.Errorf("ai_export_service.Generate: build zip: %w", err)
	}

	if int64(len(zipData)) > maxExportSize {
		log.Warn("[AiExport][generate][BLOCK_EXPORT_FAILURE] export size exceeds limit",
			zap.Int("size", len(zipData)),
			zap.Int64("max", maxExportSize),
		)
		return nil, "", ErrAiExportSizeLimit
	}

	log.Info("[AiExport][generate][BLOCK_EXPORT_ZIP_WRITE] writing ZIP to disk")

	exportDir := filepath.Join(exportBasePath, userID)
	if err := os.MkdirAll(exportDir, 0700); err != nil {
		log.Error("[AiExport][generate][BLOCK_EXPORT_FAILURE] cannot create export dir", zap.Error(err))
		return nil, "", fmt.Errorf("ai_export_service.Generate: mkdir: %w", err)
	}
	if err := os.Chmod(exportDir, 0700); err != nil {
		log.Error("[AiExport][generate][BLOCK_EXPORT_FAILURE] cannot protect export dir", zap.Error(err))
		return nil, "", fmt.Errorf("ai_export_service.Generate: chmod dir: %w", err)
	}

	tmpName := fmt.Sprintf(".tmp-%x.zip", newRandomSuffix())
	tmpPath := filepath.Join(exportDir, tmpName)
	finalName := fmt.Sprintf("%s.zip", aiExportRecord.ID)
	finalPath := filepath.Join(exportDir, finalName)

	if err := os.WriteFile(tmpPath, zipData, 0600); err != nil {
		log.Error("[AiExport][generate][BLOCK_EXPORT_FAILURE] failed to write temp file", zap.Error(err))
		os.Remove(tmpPath)
		return nil, "", fmt.Errorf("ai_export_service.Generate: write temp: %w", err)
	}
	if err := os.Chmod(tmpPath, 0600); err != nil {
		log.Error("[AiExport][generate][BLOCK_EXPORT_FAILURE] failed to protect temp file", zap.Error(err))
		os.Remove(tmpPath)
		return nil, "", fmt.Errorf("ai_export_service.Generate: chmod temp: %w", err)
	}

	if err := os.Rename(tmpPath, finalPath); err != nil {
		log.Error("[AiExport][generate][BLOCK_EXPORT_FAILURE] failed to rename temp file", zap.Error(err))
		os.Remove(tmpPath)
		return nil, "", fmt.Errorf("ai_export_service.Generate: rename: %w", err)
	}
	if err := os.Chmod(finalPath, 0600); err != nil {
		log.Error("[AiExport][generate][BLOCK_EXPORT_FAILURE] failed to protect final file", zap.Error(err))
		return nil, "", fmt.Errorf("ai_export_service.Generate: chmod final: %w", err)
	}

	log.Info("[AiExport][generate][BLOCK_EXPORT_DB_SAVE] saving export record")

	_, err = s.exportRepo.UpdateFilePath(ctx, aiExportRecord.ID, &finalPath)
	if err != nil {
		log.Error("[AiExport][generate][BLOCK_EXPORT_FAILURE] failed to update file path", zap.Error(err))
		return nil, "", fmt.Errorf("ai_export_service.Generate: update path: %w", err)
	}

	log.Info("[AiExport][generate][BLOCK_EXPORT_SUCCESS] export complete",
		zap.String("export_id", aiExportRecord.ID),
	)

	export := models.AiExportFromRecord(aiExportRecord)
	return export, prompt, nil
}

func (s *aiExportService) GetByID(ctx context.Context, userID string, id string) (*models.AiExport, error) {
	record, err := s.exportRepo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("ai_export_service.GetByID: %w", err)
	}
	if record == nil {
		return nil, ErrAiExportNotFound
	}
	return models.AiExportFromRecord(record), nil
}

func (s *aiExportService) List(ctx context.Context, userID string) ([]models.AiExport, error) {
	records, err := s.exportRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ai_export_service.List: %w", err)
	}
	out := make([]models.AiExport, len(records))
	for i := range records {
		out[i] = *models.AiExportFromRecord(&records[i])
	}
	return out, nil
}

func (s *aiExportService) Delete(ctx context.Context, userID string, id string) (*models.AiExport, error) {
	record, err := s.exportRepo.Delete(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("ai_export_service.Delete: %w", err)
	}
	if record == nil {
		return nil, ErrAiExportNotFound
	}
	return models.AiExportFromRecord(record), nil
}

func (s *aiExportService) buildDataSummary(ctx context.Context, userID string, from, to models.Date, toggles SectionToggles) string {
	var parts []string

	if toggles.Nutrition {
		dailyLogs, _ := s.dataProvider.GetDailyNutritionExport(ctx, userID, from, to)
		templates, _ := s.dataProvider.GetNutritionTemplateExport(ctx, userID, from, to)
		legacy, _ := s.dataProvider.GetLegacyNutritionExport(ctx, userID, from, to)
		if len(dailyLogs) > 0 || len(templates) > 0 || len(legacy) > 0 {
			parts = append(parts, fmt.Sprintf("- Nutrition data: %d daily logs, %d weekly templates, %d unresolved legacy days", len(dailyLogs), len(templates), len(legacy)))
		} else {
			parts = append(parts, "- Nutrition data: No entries")
		}
	}

	if toggles.Cardio {
		cardio, _ := s.dataProvider.GetCardioEntries(ctx, userID, from, to)
		if len(cardio) > 0 {
			parts = append(parts, fmt.Sprintf("- Cardio sessions: %d entries", len(cardio)))
		} else {
			parts = append(parts, "- Cardio sessions: None")
		}
	}

	if len(parts) == 0 {
		return "No workout data recorded for this period.\n"
	}

	return strings.Join(parts, "\n") + "\n"
}

func (s *aiExportService) buildArchive(ctx context.Context, userID string, dateRangeStart, dateRangeEnd string, toggles SectionToggles, profile *UserProfileExport, prompt string) (*ExportArchive, error) {
	from, _ := models.ParseDate(dateRangeStart)
	to, _ := models.ParseDate(dateRangeEnd)

	profileData := ExportProfile{
		Goal:                      profile.Goal,
		Height:                    profile.Height,
		BirthDate:                 profile.BirthDate,
		TrainingExperience:        profile.TrainingExperience,
		CurrentTrainingSplit:      profile.CurrentTrainingSplit,
		PreferredProgressionStyle: profile.PreferredProgressionStyle,
		NutritionStrategy:         profile.NutritionStrategy,
		PersistentAiContext:       profile.PersistentAiContext,
	}

	archive := NewDefaultExportArchive(dateRangeStart, dateRangeEnd, profileData)

	if toggles.Nutrition {
		dailyLogs, _ := s.dataProvider.GetDailyNutritionExport(ctx, userID, from, to)
		templates, _ := s.dataProvider.GetNutritionTemplateExport(ctx, userID, from, to)
		legacy, _ := s.dataProvider.GetLegacyNutritionExport(ctx, userID, from, to)
		archive.Data.Nutrition.DailyLogs = dailyLogs
		archive.Data.Nutrition.Templates = templates
		archive.Data.Nutrition.Legacy = legacy
		archive.NutritionCSV.Rows = nutritionCSVRowsFromExport(dailyLogs, templates)
		archive.Manifest.IncludedSections.Nutrition = true
	}

	if toggles.Cardio {
		cardio, _ := s.dataProvider.GetCardioEntries(ctx, userID, from, to)
		archive.Data.Cardio = cardio
		archive.Manifest.IncludedSections.Cardio = true
	}

	workouts, _ := s.dataProvider.GetWorkoutSummary(ctx, userID, from, to)
	archive.Data.Workouts = workouts
	if len(workouts) > 0 {
		archive.Manifest.IncludedSections.Workouts = true
	}

	bodyWeight, _ := s.dataProvider.GetBodyWeightEntries(ctx, userID, from, to)
	archive.Data.BodyWeightEntries = bodyWeight
	if len(bodyWeight) > 0 {
		archive.Manifest.IncludedSections.BodyWeight = true
	}

	measurements, _ := s.dataProvider.GetBodyMeasurements(ctx, userID, from, to)
	archive.Data.Measurements = measurements
	if len(measurements) > 0 {
		archive.Manifest.IncludedSections.Measurements = true
	}

	weekFlags, _ := s.dataProvider.GetWeekFlags(ctx, userID, from, to)
	archive.Data.WeekFlags = weekFlags

	if toggles.Photos {
		photos, _ := s.dataProvider.GetProgressPhotos(ctx, userID, from, to)
		archive.Photos = photos
		if len(photos) > 0 {
			archive.Manifest.IncludedSections.Photos = true
		}
	}

	var summaryParts []string
	summaryParts = append(summaryParts, fmt.Sprintf("# AI Export Summary\n"))
	summaryParts = append(summaryParts, fmt.Sprintf("Period: %s to %s\n\n", dateRangeStart, dateRangeEnd))
	if profile.Goal != nil && *profile.Goal != "" {
		summaryParts = append(summaryParts, fmt.Sprintf("Goal: %s\n\n", *profile.Goal))
	}
	summaryParts = append(summaryParts, "## Data\n")
	summaryParts = append(summaryParts, s.buildDataSummary(ctx, userID, from, to, toggles))
	archive.SummaryMD = strings.Join(summaryParts, "")

	rawWeekFlags, _ := json.Marshal(weekFlags)
	var parsedWeekFlags []any
	json.Unmarshal(rawWeekFlags, &parsedWeekFlags)
	archive.Data.WeekFlags = parsedWeekFlags

	return archive, nil
}

func profileRecordToExport(r *models.UserProfileRecord) *UserProfileExport {
	var height *float64
	if r.Height != nil {
		height = r.Height
	}
	var birthDate *string
	if r.BirthDate != nil {
		s := r.BirthDate.String()
		birthDate = &s
	}
	return &UserProfileExport{
		Goal:                      r.Goal,
		Height:                    height,
		BirthDate:                 birthDate,
		TrainingExperience:        r.TrainingExperience,
		CurrentTrainingSplit:      r.CurrentTrainingSplit,
		PreferredProgressionStyle: r.PreferredProgressionStyle,
		NutritionStrategy:         r.NutritionStrategy,
		PersistentAiContext:       r.PersistentAiContext,
	}
}

func newRandomSuffix() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func nutritionCSVRowsFromExport(dailyLogs []any, templates []any) [][]string {
	rows := [][]string{}
	for _, logAny := range dailyLogs {
		logMap := exportMap(logAny)
		date := exportString(logMap["date"])
		for _, entryAny := range exportSlice(logMap["entries"]) {
			rows = append(rows, nutritionEntryCSVRow(date, exportMap(entryAny)))
		}
	}
	for _, templateAny := range templates {
		templateMap := exportMap(templateAny)
		date := exportString(templateMap["weekStartDate"])
		for _, entryAny := range exportSlice(templateMap["plannedEntries"]) {
			rows = append(rows, nutritionEntryCSVRow(date, exportMap(entryAny)))
		}
	}
	return rows
}

func nutritionEntryCSVRow(date string, entry map[string]any) []string {
	return []string{
		date,
		exportString(entry["productId"]),
		exportString(entry["productNameSnapshot"]),
		exportString(entry["amountGrams"]),
		exportString(entry["caloriesPer100gSnapshot"]),
		exportString(entry["proteinPer100gSnapshot"]),
		exportString(entry["fatPer100gSnapshot"]),
		exportString(entry["carbsPer100gSnapshot"]),
		exportString(entry["entryCalories"]),
		exportString(entry["entryProtein"]),
		exportString(entry["entryFat"]),
		exportString(entry["entryCarbs"]),
		exportString(entry["mealLabel"]),
		exportString(entry["notes"]),
	}
}

func exportMap(value any) map[string]any {
	if value == nil {
		return map[string]any{}
	}
	if m, ok := value.(map[string]any); ok {
		return m
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return map[string]any{}
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		return map[string]any{}
	}
	return out
}

func exportSlice(value any) []any {
	if value == nil {
		return []any{}
	}
	if items, ok := value.([]any); ok {
		return items
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return []any{}
	}
	var out []any
	if err := json.Unmarshal(raw, &out); err != nil {
		return []any{}
	}
	return out
}

func exportString(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case *string:
		if v == nil {
			return ""
		}
		return *v
	default:
		return fmt.Sprintf("%v", v)
	}
}
