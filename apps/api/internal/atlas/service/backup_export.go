// FILE: apps/api/internal/atlas/service/backup_export.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement BackupExportService for WAVE-09 full user data backup export with ZIP archive generation.
//   SCOPE: Generate backup archives containing all user data, write to disk with atomic temp-rename, and provide download path resolution.
//   DEPENDS: apps/api/internal/atlas/repository/postgres, apps/api/internal/atlas/models, libs/go/logger.
//   LINKS: M-API / V-M-API / WAVE-09.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   BackupExportService - Interface for backup export operations.
//   NewBackupExportService - Creates a new BackupExportService.
//   Generate - Generates a full user data backup archive.
//   GetDownloadPath - Resolves the download path for a completed backup.
// END_MODULE_MAP

package service

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasRepo "monorepo-template/apps/api/internal/atlas/repository/postgres"
	"monorepo-template/libs/go/logger"
)

var (
	ErrBackupExportNotFound    = errors.New("backup export not found")
	ErrBackupExportSizeLimit   = errors.New("backup export size exceeds maximum allowed")
	ErrBackupExportFailed      = errors.New("backup export generation failed")
)

type BackupExportService interface {
	Generate(ctx context.Context, userID string, includeMedia bool, maxExportSize int64, exportBasePath string) (*models.BackupExportResult, error)
	GetDownloadPath(ctx context.Context, userID string, downloadID string) (string, error)
}

type backupExportService struct {
	backupRepo            atlasRepo.BackupArchiveRepository
	settingsRepo          atlasRepo.SettingsRepository
	profileRepo           atlasRepo.UserProfileRepository
	exerciseRepo          atlasRepo.ExerciseRepository
	cardioRepo            atlasRepo.CardioEntryRepository
	bodyWeightRepo        atlasRepo.BodyWeightEntryRepository
	checkInRepo           atlasRepo.BodyCheckInRepository
	measurementRepo       atlasRepo.BodyMeasurementRepository
	progressPhotoRepo     atlasRepo.ProgressPhotoRepository
	nutritionProductRepo  atlasRepo.NutritionProductRepository
	nutritionTemplateRepo atlasRepo.NutritionTemplateRepository
	nutritionTemplateItemRepo atlasRepo.NutritionTemplateItemRepository
	nutritionOverrideRepo atlasRepo.DailyNutritionOverrideRepository
	nutritionOverrideItemRepo atlasRepo.DailyNutritionOverrideItemRepository
	weekFlagRepo          atlasRepo.WeekFlagRepository
	aiExportRepo          atlasRepo.AiExportRepository
	aiReviewService       AiReviewService
	dailyLogRepo          atlasRepo.DailyLogRepository
	logger                *zap.Logger
}

func NewBackupExportService(
	backupRepo atlasRepo.BackupArchiveRepository,
	settingsRepo atlasRepo.SettingsRepository,
	profileRepo atlasRepo.UserProfileRepository,
	exerciseRepo atlasRepo.ExerciseRepository,
	cardioRepo atlasRepo.CardioEntryRepository,
	bodyWeightRepo atlasRepo.BodyWeightEntryRepository,
	checkInRepo atlasRepo.BodyCheckInRepository,
	measurementRepo atlasRepo.BodyMeasurementRepository,
	progressPhotoRepo atlasRepo.ProgressPhotoRepository,
	nutritionProductRepo atlasRepo.NutritionProductRepository,
	nutritionTemplateRepo atlasRepo.NutritionTemplateRepository,
	nutritionTemplateItemRepo atlasRepo.NutritionTemplateItemRepository,
	nutritionOverrideRepo atlasRepo.DailyNutritionOverrideRepository,
	nutritionOverrideItemRepo atlasRepo.DailyNutritionOverrideItemRepository,
	weekFlagRepo atlasRepo.WeekFlagRepository,
	aiExportRepo atlasRepo.AiExportRepository,
	aiReviewService AiReviewService,
	dailyLogRepo atlasRepo.DailyLogRepository,
	logger *zap.Logger,
) BackupExportService {
	return &backupExportService{
		backupRepo:               backupRepo,
		settingsRepo:             settingsRepo,
		profileRepo:              profileRepo,
		exerciseRepo:             exerciseRepo,
		cardioRepo:               cardioRepo,
		bodyWeightRepo:           bodyWeightRepo,
		checkInRepo:              checkInRepo,
		measurementRepo:          measurementRepo,
		progressPhotoRepo:        progressPhotoRepo,
		nutritionProductRepo:     nutritionProductRepo,
		nutritionTemplateRepo:    nutritionTemplateRepo,
		nutritionTemplateItemRepo: nutritionTemplateItemRepo,
		nutritionOverrideRepo:    nutritionOverrideRepo,
		nutritionOverrideItemRepo: nutritionOverrideItemRepo,
		weekFlagRepo:             weekFlagRepo,
		aiExportRepo:             aiExportRepo,
		aiReviewService:          aiReviewService,
		dailyLogRepo:             dailyLogRepo,
		logger:                   logger,
	}
}

func (s *backupExportService) Generate(ctx context.Context, userID string, includeMedia bool, maxExportSize int64, exportBasePath string) (*models.BackupExportResult, error) {
	log := logger.FromContext(ctx)
	if log == nil {
		log = s.logger
	}

	log.Info("[Backup][export][BLOCK_EXPORT_START] generating backup export",
		zap.String("user_id", userID),
		zap.Bool("include_media", includeMedia),
	)

	data, entityCounts, err := s.collectAllData(ctx, userID, includeMedia)
	if err != nil {
		log.Error("[Backup][export][BLOCK_EXPORT_FAILURE] failed to collect data", zap.Error(err))
		return nil, fmt.Errorf("backup_export_service.Generate: collect data: %w", err)
	}

	manifest := models.NewBackupManifest(includeMedia, entityCounts)
	backupArchive := models.BackupDataArchive{
		Manifest: manifest,
		Data:     *data,
	}

	manifestBytes, err := json.MarshalIndent(backupArchive.Manifest, "", "  ")
	if err != nil {
		log.Error("[Backup][export][BLOCK_EXPORT_FAILURE] failed to marshal manifest", zap.Error(err))
		return nil, fmt.Errorf("backup_export_service.Generate: marshal manifest: %w", err)
	}

	dataBytes, err := json.MarshalIndent(backupArchive.Data, "", "  ")
	if err != nil {
		log.Error("[Backup][export][BLOCK_EXPORT_FAILURE] failed to marshal data", zap.Error(err))
		return nil, fmt.Errorf("backup_export_service.Generate: marshal data: %w", err)
	}

	timestamp := time.Now().UTC().Unix()
	zipFilename := fmt.Sprintf("atlas-backup-%d.zip", timestamp)

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	mf, err := zw.Create("manifest.json")
	if err != nil {
		return nil, fmt.Errorf("backup_export_service.Generate: create manifest.json: %w", err)
	}
	if _, err := mf.Write(manifestBytes); err != nil {
		return nil, fmt.Errorf("backup_export_service.Generate: write manifest.json: %w", err)
	}

	df, err := zw.Create("data.json")
	if err != nil {
		return nil, fmt.Errorf("backup_export_service.Generate: create data.json: %w", err)
	}
	if _, err := df.Write(dataBytes); err != nil {
		return nil, fmt.Errorf("backup_export_service.Generate: write data.json: %w", err)
	}

	if includeMedia && data.ProgressPhotos != nil {
		for _, photo := range data.ProgressPhotos {
			if photo.FilePath == "" {
				continue
			}
			photoName := fmt.Sprintf("media/%s", photo.OriginalFileName)
			pf, err := zw.Create(photoName)
			if err != nil {
				log.Warn("[Backup][export][BLOCK_EXPORT_MEDIA_SKIP] failed to create photo entry in zip",
					zap.String("photo_id", photo.ID),
					zap.Error(err),
				)
				continue
			}
			photoData, err := os.ReadFile(photo.FilePath)
			if err != nil {
				log.Warn("[Backup][export][BLOCK_EXPORT_MEDIA_SKIP] failed to read photo file",
					zap.String("photo_id", photo.ID),
					zap.String("path", photo.FilePath),
					zap.Error(err),
				)
				continue
			}
			if _, err := pf.Write(photoData); err != nil {
				log.Warn("[Backup][export][BLOCK_EXPORT_MEDIA_SKIP] failed to write photo data to zip",
					zap.String("photo_id", photo.ID),
					zap.Error(err),
				)
				continue
			}
		}
	}

	if err := zw.Close(); err != nil {
		log.Error("[Backup][export][BLOCK_EXPORT_FAILURE] failed to close zip writer", zap.Error(err))
		return nil, fmt.Errorf("backup_export_service.Generate: close zip: %w", err)
	}

	zipData := buf.Bytes()

	if int64(len(zipData)) > maxExportSize {
		log.Warn("[Backup][export][BLOCK_EXPORT_FAILURE] backup size exceeds limit",
			zap.Int("size", len(zipData)),
			zap.Int64("max", maxExportSize),
		)
		return nil, ErrBackupExportSizeLimit
	}

	countsJSON, _ := json.Marshal(entityCounts)

	record, err := s.backupRepo.Create(ctx, userID, includeMedia, int64(len(zipData)), countsJSON)
	if err != nil {
		log.Error("[Backup][export][BLOCK_EXPORT_FAILURE] failed to create backup record", zap.Error(err))
		return nil, fmt.Errorf("backup_export_service.Generate: create record: %w", err)
	}

	exportDir := filepath.Join(exportBasePath, userID)
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		log.Error("[Backup][export][BLOCK_EXPORT_FAILURE] cannot create export dir", zap.Error(err))
		return nil, fmt.Errorf("backup_export_service.Generate: mkdir: %w", err)
	}

	tmpName := fmt.Sprintf(".tmp-%x.zip", newRandomSuffix())
	tmpPath := filepath.Join(exportDir, tmpName)
	finalName := zipFilename
	finalPath := filepath.Join(exportDir, finalName)

	if err := os.WriteFile(tmpPath, zipData, 0644); err != nil {
		log.Error("[Backup][export][BLOCK_EXPORT_FAILURE] failed to write temp file", zap.Error(err))
		os.Remove(tmpPath)
		return nil, fmt.Errorf("backup_export_service.Generate: write temp: %w", err)
	}

	if err := os.Rename(tmpPath, finalPath); err != nil {
		log.Error("[Backup][export][BLOCK_EXPORT_FAILURE] failed to rename temp file", zap.Error(err))
		os.Remove(tmpPath)
		return nil, fmt.Errorf("backup_export_service.Generate: rename: %w", err)
	}

	_, err = s.backupRepo.UpdateFilePath(ctx, record.ID, &finalPath)
	if err != nil {
		log.Error("[Backup][export][BLOCK_EXPORT_FAILURE] failed to update backup file path", zap.Error(err))
		return nil, fmt.Errorf("backup_export_service.Generate: update path: %w", err)
	}

	log.Info("[Backup][export][BLOCK_EXPORT_SUCCESS] backup export complete",
		zap.String("backup_id", record.ID),
		zap.Int("size_bytes", len(zipData)),
	)

	return &models.BackupExportResult{
		DownloadID: record.ID,
		SizeBytes:  int64(len(zipData)),
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func (s *backupExportService) GetDownloadPath(ctx context.Context, userID string, downloadID string) (string, error) {
	record, err := s.backupRepo.GetByID(ctx, userID, downloadID)
	if err != nil {
		return "", fmt.Errorf("backup_export_service.GetDownloadPath: %w", err)
	}
	if record == nil {
		return "", ErrBackupExportNotFound
	}
	if record.ArchivePath == nil {
		return "", ErrBackupExportNotFound
	}
	return *record.ArchivePath, nil
}

func (s *backupExportService) collectAllData(ctx context.Context, userID string, includeMedia bool) (*models.BackupData, map[string]int, error) {
	log := logger.FromContext(ctx)
	if log == nil {
		log = s.logger
	}

	data := &models.BackupData{}
	entityCounts := make(map[string]int)

	settings, err := s.settingsRepo.FindByUserID(ctx, userID)
	if err != nil {
		log.Warn("[Backup][export][BLOCK_COLLECT_SKIP] failed to fetch settings", zap.Error(err))
	}
	if settings != nil {
		data.Settings = settings
		entityCounts["settings"] = 1
	}

	profile, err := s.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		log.Warn("[Backup][export][BLOCK_COLLECT_SKIP] failed to fetch user profile", zap.Error(err))
	}
	if profile != nil {
		data.UserProfile = profile
		entityCounts["userProfile"] = 1
	}

	exercises, err := s.exerciseRepo.ListAll(ctx, userID, true)
	if err != nil {
		log.Warn("[Backup][export][BLOCK_COLLECT_SKIP] failed to fetch exercises", zap.Error(err))
	}
	if len(exercises) > 0 {
		data.Exercises = exercises
		entityCounts["exercises"] = len(exercises)

		for _, ex := range exercises {
			media, mediaErr := s.exerciseRepo.ListMediaByExercise(ctx, userID, ex.ID)
			if mediaErr != nil {
				log.Warn("[Backup][export][BLOCK_COLLECT_SKIP] failed to fetch exercise media",
					zap.String("exercise_id", ex.ID),
					zap.Error(mediaErr),
				)
				continue
			}
			for _, m := range media {
				data.ExerciseMedia = append(data.ExerciseMedia, models.ExerciseMediaRecord{
					ID:         m.ID,
					UserID:     m.UserID,
					ExerciseID: m.ExerciseID,
					FileName:   m.FileName,
					MimeType:   m.MimeType,
					FileSize:   m.FileSize,
					CreatedAt:  m.CreatedAt,
				})
			}
		}
		entityCounts["exerciseMedia"] = len(data.ExerciseMedia)
	}

	entityCounts["dailyLogs"] = 0
	entityCounts["workoutExercises"] = 0
	entityCounts["workoutSets"] = 0
	entityCounts["cardioEntries"] = 0
	entityCounts["bodyWeightEntries"] = 0
	entityCounts["bodyCheckIns"] = 0
	entityCounts["bodyMeasurements"] = 0

	entityCounts["nutritionProducts"] = 0
	entityCounts["nutritionTemplates"] = 0
	entityCounts["nutritionTemplateItems"] = 0
	entityCounts["nutritionOverrides"] = 0
	entityCounts["nutritionOverrideItems"] = 0
	entityCounts["weekFlags"] = 0

	aiExports, err := s.aiExportRepo.ListByUserID(ctx, userID)
	if err != nil {
		log.Warn("[Backup][export][BLOCK_COLLECT_SKIP] failed to fetch AI exports", zap.Error(err))
	}
	if len(aiExports) > 0 {
		data.AiExports = aiExports
		entityCounts["aiExports"] = len(aiExports)
	}

	aiReviews, err := s.aiReviewService.ListAllByUserID(ctx, userID)
	if err != nil {
		log.Warn("[Backup][export][BLOCK_COLLECT_SKIP] failed to fetch AI reviews", zap.Error(err))
	}
	if len(aiReviews) > 0 {
		out := make([]models.AiReviewRecord, len(aiReviews))
		for i, r := range aiReviews {
			out[i] = models.AiReviewRecord{
				ID:              r.ID,
				UserID:          r.UserID,
				DateRangeStart:  r.DateRangeStart,
				DateRangeEnd:    r.DateRangeEnd,
				AiResponseText:  r.AiResponseText,
				UserNotes:       r.UserNotes,
				PlannedActions:  r.PlannedActions,
				CreatedAt:       r.CreatedAt,
				UpdatedAt:       r.UpdatedAt,
			}
		}
		data.AiReviews = out
		entityCounts["aiReviews"] = len(out)
	}

	log.Info("[Backup][export][BLOCK_COLLECT_COMPLETE] data collection complete",
		zap.Any("entity_counts", entityCounts),
	)

	return data, entityCounts, nil
}