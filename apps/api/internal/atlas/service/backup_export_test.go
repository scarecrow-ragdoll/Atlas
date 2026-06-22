// FILE: apps/api/internal/atlas/service/backup_export_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for BackupExportService covering Generate (success, media, size limit, create record error) and GetDownloadPath (success, not found).
//   SCOPE: Success paths for Generate with and without media, size limit exceeded, backup repo create failure, download path success and not found.
//   DEPENDS: apps/api/internal/atlas/service, apps/api/internal/atlas/repository/postgres (mocked), apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-09.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added BackupExportService unit tests for WAVE-09.
// END_CHANGE_SUMMARY

package service_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasPostgres "monorepo-template/apps/api/internal/atlas/repository/postgres"
	"monorepo-template/apps/api/internal/atlas/service"
)

// ----- mocks for repos used by BackupExportService -----

type beMockBackupRepo struct {
	atlasPostgres.BackupArchiveRepository
	createFn         func(ctx context.Context, userID string, includeMedia bool, sizeBytes int64, entityCounts []byte) (*models.BackupArchiveRecord, error)
	getByIDFn        func(ctx context.Context, userID string, id string) (*models.BackupArchiveRecord, error)
	updateFilePathFn func(ctx context.Context, id string, filePath *string) (*models.BackupArchiveRecord, error)
}

func (m *beMockBackupRepo) Create(ctx context.Context, userID string, includeMedia bool, sizeBytes int64, entityCounts []byte) (*models.BackupArchiveRecord, error) {
	return m.createFn(ctx, userID, includeMedia, sizeBytes, entityCounts)
}

func (m *beMockBackupRepo) GetByID(ctx context.Context, userID string, id string) (*models.BackupArchiveRecord, error) {
	return m.getByIDFn(ctx, userID, id)
}

func (m *beMockBackupRepo) UpdateFilePath(ctx context.Context, id string, filePath *string) (*models.BackupArchiveRecord, error) {
	return m.updateFilePathFn(ctx, id, filePath)
}

type beMockSettingsRepo struct {
	atlasPostgres.SettingsRepository
	findByUserIDFn func(ctx context.Context, userID string) (*models.SettingsRecord, error)
}

func (m *beMockSettingsRepo) FindByUserID(ctx context.Context, userID string) (*models.SettingsRecord, error) {
	return m.findByUserIDFn(ctx, userID)
}

type beMockProfileRepo struct {
	atlasPostgres.UserProfileRepository
	findByUserIDFn func(ctx context.Context, userID string) (*models.UserProfileRecord, error)
}

func (m *beMockProfileRepo) FindByUserID(ctx context.Context, userID string) (*models.UserProfileRecord, error) {
	return m.findByUserIDFn(ctx, userID)
}

type beMockExerciseRepo struct {
	atlasPostgres.ExerciseRepository
	listAllFn   func(ctx context.Context, userID string, includeInactive bool) ([]models.ExerciseRecord, error)
	listMediaFn func(ctx context.Context, userID string, exerciseID string) ([]models.ExerciseMedia, error)
}

func (m *beMockExerciseRepo) ListAll(ctx context.Context, userID string, includeInactive bool) ([]models.ExerciseRecord, error) {
	return m.listAllFn(ctx, userID, includeInactive)
}

func (m *beMockExerciseRepo) ListMediaByExercise(ctx context.Context, userID string, exerciseID string) ([]models.ExerciseMedia, error) {
	return m.listMediaFn(ctx, userID, exerciseID)
}

type beMockAiExportRepo struct {
	atlasPostgres.AiExportRepository
	listByUserIDFn func(ctx context.Context, userID string) ([]models.AiExportRecord, error)
}

func (m *beMockAiExportRepo) ListByUserID(ctx context.Context, userID string) ([]models.AiExportRecord, error) {
	return m.listByUserIDFn(ctx, userID)
}

type beMockAiReviewSvc struct {
	service.AiReviewService
	listAllByUserIDFn func(ctx context.Context, userID string) ([]models.AiReview, error)
}

func (m *beMockAiReviewSvc) ListAllByUserID(ctx context.Context, userID string) ([]models.AiReview, error) {
	return m.listAllByUserIDFn(ctx, userID)
}

// ----- noop mocks for unused repos (embed nil interface, never called) -----

type beNoopCardioRepo struct{ atlasPostgres.CardioEntryRepository }
type beNoopBodyWeightRepo struct{ atlasPostgres.BodyWeightEntryRepository }
type beNoopBodyCheckInRepo struct{ atlasPostgres.BodyCheckInRepository }
type beNoopBodyMeasurementRepo struct{ atlasPostgres.BodyMeasurementRepository }
type beNoopProgressPhotoRepo struct{ atlasPostgres.ProgressPhotoRepository }
type beNoopNutritionProductRepo struct{ atlasPostgres.NutritionProductRepository }
type beNoopNutritionTemplateRepo struct{ atlasPostgres.NutritionTemplateRepository }
type beNoopNutritionTemplateItemRepo struct{ atlasPostgres.NutritionTemplateItemRepository }
type beNoopNutritionOverrideRepo struct{ atlasPostgres.DailyNutritionOverrideRepository }
type beNoopNutritionOverrideItemRepo struct{ atlasPostgres.DailyNutritionOverrideItemRepository }
type beNoopWeekFlagRepo struct{ atlasPostgres.WeekFlagRepository }
type beNoopDailyLogRepo struct{ atlasPostgres.DailyLogRepository }

// ----- helpers -----

var beCtx = context.Background()
var beUserID = "550e8400-e29b-41d4-a716-446655440000"
var beBackupID = "660e8400-e29b-41d4-a716-446655440001"

func beNewService(
	backupRepo *beMockBackupRepo,
	settingsRepo *beMockSettingsRepo,
	profileRepo *beMockProfileRepo,
	exerciseRepo *beMockExerciseRepo,
	aiExportRepo *beMockAiExportRepo,
	aiReviewSvc *beMockAiReviewSvc,
) service.BackupExportService {
	return service.NewBackupExportService(
		backupRepo,
		settingsRepo,
		profileRepo,
		exerciseRepo,
		&beNoopCardioRepo{},
		&beNoopBodyWeightRepo{},
		&beNoopBodyCheckInRepo{},
		&beNoopBodyMeasurementRepo{},
		&beNoopProgressPhotoRepo{},
		&beNoopNutritionProductRepo{},
		&beNoopNutritionTemplateRepo{},
		&beNoopNutritionTemplateItemRepo{},
		&beNoopNutritionOverrideRepo{},
		&beNoopNutritionOverrideItemRepo{},
		&beNoopWeekFlagRepo{},
		aiExportRepo,
		aiReviewSvc,
		&beNoopDailyLogRepo{},
		zap.NewNop(),
	)
}

func beDefaultMocks() (
	*beMockBackupRepo,
	*beMockSettingsRepo,
	*beMockProfileRepo,
	*beMockExerciseRepo,
	*beMockAiExportRepo,
	*beMockAiReviewSvc,
) {
	backupRepo := &beMockBackupRepo{
		createFn: func(ctx context.Context, userID string, includeMedia bool, sizeBytes int64, entityCounts []byte) (*models.BackupArchiveRecord, error) {
			return &models.BackupArchiveRecord{
				ID:           beBackupID,
				UserID:       userID,
				IncludeMedia: includeMedia,
				SizeBytes:    sizeBytes,
				EntityCounts: string(entityCounts),
				ArchivePath:  nil,
				CreatedAt:    "2026-06-22T00:00:00Z",
				UpdatedAt:    "2026-06-22T00:00:00Z",
			}, nil
		},
		updateFilePathFn: func(ctx context.Context, id string, filePath *string) (*models.BackupArchiveRecord, error) {
			return &models.BackupArchiveRecord{ID: id, ArchivePath: filePath}, nil
		},
	}
	settingsRepo := &beMockSettingsRepo{
		findByUserIDFn: func(ctx context.Context, userID string) (*models.SettingsRecord, error) {
			return &models.SettingsRecord{ID: "s1", UserID: userID, Units: "metric"}, nil
		},
	}
	profileRepo := &beMockProfileRepo{
		findByUserIDFn: func(ctx context.Context, userID string) (*models.UserProfileRecord, error) {
			return &models.UserProfileRecord{ID: "p1", UserID: userID}, nil
		},
	}
	exerciseRepo := &beMockExerciseRepo{
		listAllFn: func(ctx context.Context, userID string, includeInactive bool) ([]models.ExerciseRecord, error) {
			return nil, nil
		},
	}
	aiExportRepo := &beMockAiExportRepo{
		listByUserIDFn: func(ctx context.Context, userID string) ([]models.AiExportRecord, error) {
			return nil, nil
		},
	}
	aiReviewSvc := &beMockAiReviewSvc{
		listAllByUserIDFn: func(ctx context.Context, userID string) ([]models.AiReview, error) {
			return nil, nil
		},
	}
	return backupRepo, settingsRepo, profileRepo, exerciseRepo, aiExportRepo, aiReviewSvc
}

// ----- Generate tests -----

func TestBackupExportService_Generate_Success(t *testing.T) {
	backupRepo, settingsRepo, profileRepo, exerciseRepo, aiExportRepo, aiReviewSvc := beDefaultMocks()
	svc := beNewService(backupRepo, settingsRepo, profileRepo, exerciseRepo, aiExportRepo, aiReviewSvc)
	exportDir := t.TempDir()

	result, err := svc.Generate(beCtx, beUserID, false, 10*1024*1024, exportDir)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, beBackupID, result.DownloadID)
	assert.Greater(t, result.SizeBytes, int64(0))
	assert.NotEmpty(t, result.Timestamp)
}

func TestBackupExportService_Generate_WithMedia(t *testing.T) {
	backupRepo, settingsRepo, profileRepo, exerciseRepo, aiExportRepo, aiReviewSvc := beDefaultMocks()
	svc := beNewService(backupRepo, settingsRepo, profileRepo, exerciseRepo, aiExportRepo, aiReviewSvc)
	exportDir := t.TempDir()

	result, err := svc.Generate(beCtx, beUserID, true, 10*1024*1024, exportDir)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, beBackupID, result.DownloadID)
	assert.Greater(t, result.SizeBytes, int64(0))
	assert.NotEmpty(t, result.Timestamp)
}

func TestBackupExportService_Generate_SizeLimitExceeded(t *testing.T) {
	backupRepo, settingsRepo, profileRepo, exerciseRepo, aiExportRepo, aiReviewSvc := beDefaultMocks()
	svc := beNewService(backupRepo, settingsRepo, profileRepo, exerciseRepo, aiExportRepo, aiReviewSvc)
	exportDir := t.TempDir()

	result, err := svc.Generate(beCtx, beUserID, false, 0, exportDir)

	assert.ErrorIs(t, err, service.ErrBackupExportSizeLimit)
	assert.Nil(t, result)
}

func TestBackupExportService_Generate_CreateRecordError(t *testing.T) {
	backupRepo := &beMockBackupRepo{
		createFn: func(ctx context.Context, userID string, includeMedia bool, sizeBytes int64, entityCounts []byte) (*models.BackupArchiveRecord, error) {
			return nil, errors.New("database error")
		},
	}
	settingsRepo := &beMockSettingsRepo{
		findByUserIDFn: func(ctx context.Context, userID string) (*models.SettingsRecord, error) {
			return nil, nil
		},
	}
	profileRepo := &beMockProfileRepo{
		findByUserIDFn: func(ctx context.Context, userID string) (*models.UserProfileRecord, error) {
			return nil, nil
		},
	}
	exerciseRepo := &beMockExerciseRepo{
		listAllFn: func(ctx context.Context, userID string, includeInactive bool) ([]models.ExerciseRecord, error) {
			return nil, nil
		},
	}
	aiExportRepo := &beMockAiExportRepo{
		listByUserIDFn: func(ctx context.Context, userID string) ([]models.AiExportRecord, error) {
			return nil, nil
		},
	}
	aiReviewSvc := &beMockAiReviewSvc{
		listAllByUserIDFn: func(ctx context.Context, userID string) ([]models.AiReview, error) {
			return nil, nil
		},
	}
	svc := beNewService(backupRepo, settingsRepo, profileRepo, exerciseRepo, aiExportRepo, aiReviewSvc)
	exportDir := t.TempDir()

	result, err := svc.Generate(beCtx, beUserID, false, 10*1024*1024, exportDir)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "create record")
}

// ----- GetDownloadPath tests -----

func TestBackupExportService_GetDownloadPath_Success(t *testing.T) {
	expectedPath := filepath.Join("/tmp", "backups", "atlas-backup-1234.zip")
	backupRepo := &beMockBackupRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.BackupArchiveRecord, error) {
			return &models.BackupArchiveRecord{
				ID:          id,
				UserID:      userID,
				ArchivePath: &expectedPath,
			}, nil
		},
	}
	svc := beNewService(backupRepo,
		&beMockSettingsRepo{},
		&beMockProfileRepo{},
		&beMockExerciseRepo{},
		&beMockAiExportRepo{},
		&beMockAiReviewSvc{},
	)

	path, err := svc.GetDownloadPath(beCtx, beUserID, beBackupID)

	require.NoError(t, err)
	assert.Equal(t, expectedPath, path)
}

func TestBackupExportService_GetDownloadPath_NotFound(t *testing.T) {
	t.Run("nil record", func(t *testing.T) {
		backupRepo := &beMockBackupRepo{
			getByIDFn: func(ctx context.Context, userID string, id string) (*models.BackupArchiveRecord, error) {
				return nil, nil
			},
		}
		svc := beNewService(backupRepo,
			&beMockSettingsRepo{},
			&beMockProfileRepo{},
			&beMockExerciseRepo{},
			&beMockAiExportRepo{},
			&beMockAiReviewSvc{},
		)

		path, err := svc.GetDownloadPath(beCtx, beUserID, beBackupID)

		assert.ErrorIs(t, err, service.ErrBackupExportNotFound)
		assert.Empty(t, path)
	})

	t.Run("nil archive path", func(t *testing.T) {
		backupRepo := &beMockBackupRepo{
			getByIDFn: func(ctx context.Context, userID string, id string) (*models.BackupArchiveRecord, error) {
				return &models.BackupArchiveRecord{
					ID:          id,
					UserID:      userID,
					ArchivePath: nil,
				}, nil
			},
		}
		svc := beNewService(backupRepo,
			&beMockSettingsRepo{},
			&beMockProfileRepo{},
			&beMockExerciseRepo{},
			&beMockAiExportRepo{},
			&beMockAiReviewSvc{},
		)

		path, err := svc.GetDownloadPath(beCtx, beUserID, beBackupID)

		assert.ErrorIs(t, err, service.ErrBackupExportNotFound)
		assert.Empty(t, path)
	})
}
