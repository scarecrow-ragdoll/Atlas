// FILE: apps/api/internal/atlas/service/exercise_service_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for ExerciseService covering Create, GetByID, List, ListAll, Update, Archive, Restore, media operations, validation, pagination, and media field initialization.
//   SCOPE: Success paths, validation errors (empty name, whitespace name, invalid weight), not-found errors, pagination with/without cursor, default page size, hasNextPage logic, empty list, ListAll filter, media CRUD, and media non-nil guarantee for all exercise return paths.
//   DEPENDS: apps/api/internal/atlas/service, apps/api/internal/atlas/repository/postgres (mock ExerciseRepository), apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-02.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added exercise service unit tests.
// END_CHANGE_SUMMARY

package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasPostgres "monorepo-template/apps/api/internal/atlas/repository/postgres"
	"monorepo-template/apps/api/internal/atlas/service"
)

var (
	ctx        = context.Background()
	testUserID = "550e8400-e29b-41d4-a716-446655440000"
	testID     = "660e8400-e29b-41d4-a716-446655440001"
)

func ptrFloat64(f float64) *float64 { return &f }

// mockExerciseRepo embeds ExerciseRepository and overrides selected methods.
type mockExerciseRepo struct {
	atlasPostgres.ExerciseRepository
	createFn             func(ctx context.Context, userID string, input models.CreateExerciseInput) (*models.ExerciseRecord, error)
	getByIDFn            func(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error)
	listFn               func(ctx context.Context, userID string, isActive bool, limit int32) ([]models.ExerciseRecord, error)
	listCursorFn         func(ctx context.Context, userID string, isActive bool, cursor string, limit int32) ([]models.ExerciseRecord, error)
	listAllFn            func(ctx context.Context, userID string, includeInactive bool) ([]models.ExerciseRecord, error)
	countFn              func(ctx context.Context, userID string, isActive bool) (int, error)
	updateFn             func(ctx context.Context, userID string, id string, input models.UpdateExerciseInput) (*models.ExerciseRecord, error)
	archiveFn            func(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error)
	restoreFn            func(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error)
	createMediaFn        func(ctx context.Context, userID string, exerciseID string, fileName string, filePath string, mimeType string, fileSize int64) (*models.ExerciseMedia, error)
	getMediaByIDFn       func(ctx context.Context, userID string, id string) (*models.ExerciseMedia, error)
	getMediaRecordByIDFn func(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error)
	listMediaByExerciseFn func(ctx context.Context, userID string, exerciseID string) ([]models.ExerciseMedia, error)
	deleteMediaFn        func(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error)
}

func (m *mockExerciseRepo) Create(ctx context.Context, userID string, input models.CreateExerciseInput) (*models.ExerciseRecord, error) {
	return m.createFn(ctx, userID, input)
}

func (m *mockExerciseRepo) GetByID(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error) {
	return m.getByIDFn(ctx, userID, id)
}

func (m *mockExerciseRepo) List(ctx context.Context, userID string, isActive bool, limit int32) ([]models.ExerciseRecord, error) {
	return m.listFn(ctx, userID, isActive, limit)
}

func (m *mockExerciseRepo) ListCursor(ctx context.Context, userID string, isActive bool, cursor string, limit int32) ([]models.ExerciseRecord, error) {
	return m.listCursorFn(ctx, userID, isActive, cursor, limit)
}

func (m *mockExerciseRepo) ListAll(ctx context.Context, userID string, includeInactive bool) ([]models.ExerciseRecord, error) {
	return m.listAllFn(ctx, userID, includeInactive)
}

func (m *mockExerciseRepo) Count(ctx context.Context, userID string, isActive bool) (int, error) {
	return m.countFn(ctx, userID, isActive)
}

func (m *mockExerciseRepo) Update(ctx context.Context, userID string, id string, input models.UpdateExerciseInput) (*models.ExerciseRecord, error) {
	return m.updateFn(ctx, userID, id, input)
}

func (m *mockExerciseRepo) Archive(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error) {
	return m.archiveFn(ctx, userID, id)
}

func (m *mockExerciseRepo) Restore(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error) {
	return m.restoreFn(ctx, userID, id)
}

func (m *mockExerciseRepo) CreateMedia(ctx context.Context, userID string, exerciseID string, fileName string, filePath string, mimeType string, fileSize int64) (*models.ExerciseMedia, error) {
	return m.createMediaFn(ctx, userID, exerciseID, fileName, filePath, mimeType, fileSize)
}

func (m *mockExerciseRepo) GetMediaByID(ctx context.Context, userID string, id string) (*models.ExerciseMedia, error) {
	return m.getMediaByIDFn(ctx, userID, id)
}

func (m *mockExerciseRepo) GetMediaRecordByID(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error) {
	return m.getMediaRecordByIDFn(ctx, userID, id)
}

func (m *mockExerciseRepo) ListMediaByExercise(ctx context.Context, userID string, exerciseID string) ([]models.ExerciseMedia, error) {
	return m.listMediaByExerciseFn(ctx, userID, exerciseID)
}

func (m *mockExerciseRepo) DeleteMedia(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error) {
	return m.deleteMediaFn(ctx, userID, id)
}

func baseRecord() *models.ExerciseRecord {
	return &models.ExerciseRecord{
		ID:            testID,
		UserID:        testUserID,
		Name:          "Bench Press",
		MuscleGroups:  []string{"chest", "triceps"},
		Description:   ptrStr("Barbell bench press"),
		PersonalNotes: ptrStr("Focus on form"),
		WorkingWeight: ptrFloat64(80.5),
		IsActive:      true,
		CreatedAt:     "2025-01-01T00:00:00Z",
		UpdatedAt:     "2025-01-01T00:00:00Z",
	}
}

func baseMedia() *models.ExerciseMedia {
	return &models.ExerciseMedia{
		ID:         "770e8400-e29b-41d4-a716-446655440000",
		UserID:     testUserID,
		ExerciseID: testID,
		FileName:   "bench.mp4",
		MimeType:   "video/mp4",
		FileSize:   1024000,
		CreatedAt:  "2025-01-01T00:00:00Z",
	}
}

// ----- Create -----

func TestExerciseService_Create_Success(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{
		createFn: func(ctx context.Context, userID string, input models.CreateExerciseInput) (*models.ExerciseRecord, error) {
			assert.Equal(t, "Bench Press", input.Name)
			return baseRecord(), nil
		},
	})

	ex, err := svc.Create(ctx, testUserID, models.CreateExerciseInput{
		Name:          "  Bench Press  ",
		MuscleGroups:  []string{"chest", "triceps"},
		Description:   ptrStr("Barbell bench press"),
		PersonalNotes: ptrStr("Focus on form"),
		WorkingWeight: ptrFloat64(80.5),
	})
	require.NoError(t, err)
	require.NotNil(t, ex)
	assert.Equal(t, "Bench Press", ex.Name)
	assert.NotNil(t, ex.Media)
	assert.Empty(t, ex.Media)
}

func TestExerciseService_Create_EmptyName(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{})

	ex, err := svc.Create(ctx, testUserID, models.CreateExerciseInput{
		Name: "",
	})
	assert.ErrorIs(t, err, service.ErrExerciseNameEmpty)
	assert.Nil(t, ex)
}

func TestExerciseService_Create_WhitespaceName(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{})

	ex, err := svc.Create(ctx, testUserID, models.CreateExerciseInput{
		Name: "   \t\n  ",
	})
	assert.ErrorIs(t, err, service.ErrExerciseNameEmpty)
	assert.Nil(t, ex)
}

func TestExerciseService_Create_InvalidWeight(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{})

	t.Run("zero weight", func(t *testing.T) {
		ex, err := svc.Create(ctx, testUserID, models.CreateExerciseInput{
			Name:          "Bench Press",
			WorkingWeight: ptrFloat64(0),
		})
		assert.ErrorIs(t, err, service.ErrWeightInvalid)
		assert.Nil(t, ex)
	})

	t.Run("negative weight", func(t *testing.T) {
		ex, err := svc.Create(ctx, testUserID, models.CreateExerciseInput{
			Name:          "Bench Press",
			WorkingWeight: ptrFloat64(-10),
		})
		assert.ErrorIs(t, err, service.ErrWeightInvalid)
		assert.Nil(t, ex)
	})
}

func TestExerciseService_Create_DuplicateNamesAllowed(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{
		createFn: func(ctx context.Context, userID string, input models.CreateExerciseInput) (*models.ExerciseRecord, error) {
			return &models.ExerciseRecord{
				ID:       testID,
				UserID:   userID,
				Name:     input.Name,
				IsActive: true,
			}, nil
		},
	})

	ex1, err := svc.Create(ctx, testUserID, models.CreateExerciseInput{Name: "Squat"})
	require.NoError(t, err)
	require.NotNil(t, ex1)
	assert.Equal(t, "Squat", ex1.Name)

	ex2, err := svc.Create(ctx, testUserID, models.CreateExerciseInput{Name: "Squat"})
	require.NoError(t, err)
	require.NotNil(t, ex2)
	assert.Equal(t, "Squat", ex2.Name)
}

// ----- GetByID -----

func TestExerciseService_GetByID_Success(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error) {
			return baseRecord(), nil
		},
		listMediaByExerciseFn: func(ctx context.Context, userID string, exerciseID string) ([]models.ExerciseMedia, error) {
			return []models.ExerciseMedia{*baseMedia()}, nil
		},
	})

	ex, err := svc.GetByID(ctx, testUserID, testID)
	require.NoError(t, err)
	require.NotNil(t, ex)
	assert.Equal(t, "Bench Press", ex.Name)
	require.Len(t, ex.Media, 1)
	assert.Equal(t, "bench.mp4", ex.Media[0].FileName)
}

func TestExerciseService_GetByID_NotFound(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error) {
			return nil, nil
		},
	})

	ex, err := svc.GetByID(ctx, testUserID, testID)
	assert.ErrorIs(t, err, service.ErrExerciseNotFound)
	assert.Nil(t, ex)
}

func TestExerciseService_GetByID_MediaAlwaysInitialized(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error) {
			return baseRecord(), nil
		},
		listMediaByExerciseFn: func(ctx context.Context, userID string, exerciseID string) ([]models.ExerciseMedia, error) {
			return []models.ExerciseMedia{}, nil
		},
	})

	ex, err := svc.GetByID(ctx, testUserID, testID)
	require.NoError(t, err)
	require.NotNil(t, ex)
	assert.NotNil(t, ex.Media)
	assert.Empty(t, ex.Media)
}

// ----- List -----

func TestExerciseService_List_Success(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{
		listFn: func(ctx context.Context, userID string, isActive bool, limit int32) ([]models.ExerciseRecord, error) {
			assert.True(t, isActive)
			assert.Equal(t, int32(21), limit)
			return []models.ExerciseRecord{*baseRecord()}, nil
		},
		countFn: func(ctx context.Context, userID string, isActive bool) (int, error) {
			return 1, nil
		},
	})

	conn, err := svc.List(ctx, testUserID, 20, nil, false)
	require.NoError(t, err)
	require.NotNil(t, conn)
	assert.Len(t, conn.Items, 1)
	assert.Equal(t, 1, conn.TotalCount)
	assert.False(t, conn.PageInfo.HasNextPage)
}

func TestExerciseService_List_DefaultPageSize(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{
		listFn: func(ctx context.Context, userID string, isActive bool, limit int32) ([]models.ExerciseRecord, error) {
			assert.Equal(t, int32(21), limit)
			return []models.ExerciseRecord{}, nil
		},
		countFn: func(ctx context.Context, userID string, isActive bool) (int, error) {
			return 0, nil
		},
	})

	conn, err := svc.List(ctx, testUserID, 0, nil, false)
	require.NoError(t, err)
	assert.Empty(t, conn.Items)
	assert.Equal(t, 0, conn.TotalCount)
	assert.False(t, conn.PageInfo.HasNextPage)
}

func TestExerciseService_List_HasNextPage(t *testing.T) {
	records := make([]models.ExerciseRecord, 21)
	for i := range records {
		records[i] = models.ExerciseRecord{
			ID:   testID,
			Name: "Exercise",
		}
	}

	svc := service.NewExerciseService(&mockExerciseRepo{
		listFn: func(ctx context.Context, userID string, isActive bool, limit int32) ([]models.ExerciseRecord, error) {
			return records, nil
		},
		countFn: func(ctx context.Context, userID string, isActive bool) (int, error) {
			return 25, nil
		},
	})

	conn, err := svc.List(ctx, testUserID, 20, nil, false)
	require.NoError(t, err)
	assert.Len(t, conn.Items, 20)
	assert.True(t, conn.PageInfo.HasNextPage)
	assert.Equal(t, 25, conn.TotalCount)
	require.NotNil(t, conn.PageInfo.EndCursor)
}

func TestExerciseService_List_WithCursor(t *testing.T) {
	after := "Bench Press"
	records := []models.ExerciseRecord{
		{
			ID:   testID,
			Name: "Deadlift",
		},
	}

	svc := service.NewExerciseService(&mockExerciseRepo{
		listCursorFn: func(ctx context.Context, userID string, isActive bool, cursor string, limit int32) ([]models.ExerciseRecord, error) {
			assert.Equal(t, "Bench Press", cursor)
			assert.Equal(t, int32(21), limit)
			return records, nil
		},
		countFn: func(ctx context.Context, userID string, isActive bool) (int, error) {
			return 2, nil
		},
	})

	conn, err := svc.List(ctx, testUserID, 20, &after, false)
	require.NoError(t, err)
	require.Len(t, conn.Items, 1)
	assert.Equal(t, "Deadlift", conn.Items[0].Name)
	assert.False(t, conn.PageInfo.HasNextPage)
}

func TestExerciseService_List_Empty(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{
		listFn: func(ctx context.Context, userID string, isActive bool, limit int32) ([]models.ExerciseRecord, error) {
			return []models.ExerciseRecord{}, nil
		},
		countFn: func(ctx context.Context, userID string, isActive bool) (int, error) {
			return 0, nil
		},
	})

	conn, err := svc.List(ctx, testUserID, 20, nil, false)
	require.NoError(t, err)
	assert.Empty(t, conn.Items)
	assert.Equal(t, 0, conn.TotalCount)
	assert.False(t, conn.PageInfo.HasNextPage)
	assert.Nil(t, conn.PageInfo.EndCursor)
}

func TestExerciseService_List_IncludeInactive(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{
		listFn: func(ctx context.Context, userID string, isActive bool, limit int32) ([]models.ExerciseRecord, error) {
			assert.False(t, isActive)
			return []models.ExerciseRecord{*baseRecord()}, nil
		},
		countFn: func(ctx context.Context, userID string, isActive bool) (int, error) {
			return 1, nil
		},
	})

	conn, err := svc.List(ctx, testUserID, 20, nil, true)
	require.NoError(t, err)
	require.Len(t, conn.Items, 1)
}

// ----- ListAll -----

func TestExerciseService_ListAll_IncludeInactiveTrue(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{
		listAllFn: func(ctx context.Context, userID string, includeInactive bool) ([]models.ExerciseRecord, error) {
			assert.True(t, includeInactive)
			return []models.ExerciseRecord{
				{ID: "1", Name: "Active", IsActive: true},
				{ID: "2", Name: "Archived", IsActive: false},
			}, nil
		},
	})

	items, err := svc.ListAll(ctx, testUserID, true)
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.NotNil(t, items[0].Media)
	assert.Empty(t, items[0].Media)
}

func TestExerciseService_ListAll_IncludeInactiveFalse(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{
		listAllFn: func(ctx context.Context, userID string, includeInactive bool) ([]models.ExerciseRecord, error) {
			assert.False(t, includeInactive)
			return []models.ExerciseRecord{
				{ID: "1", Name: "Active", IsActive: true},
			}, nil
		},
	})

	items, err := svc.ListAll(ctx, testUserID, false)
	require.NoError(t, err)
	require.Len(t, items, 1)
}

// ----- Update -----

func TestExerciseService_Update_Success(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{
		updateFn: func(ctx context.Context, userID string, id string, input models.UpdateExerciseInput) (*models.ExerciseRecord, error) {
			assert.Equal(t, "Squat", *input.Name)
			return &models.ExerciseRecord{
				ID:     id,
				UserID: userID,
				Name:   *input.Name,
			}, nil
		},
	})

	ex, err := svc.Update(ctx, testUserID, testID, models.UpdateExerciseInput{
		Name: ptrStr("  Squat  "),
	})
	require.NoError(t, err)
	require.NotNil(t, ex)
	assert.Equal(t, "Squat", ex.Name)
	assert.NotNil(t, ex.Media)
	assert.Empty(t, ex.Media)
}

func TestExerciseService_Update_EmptyName(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{})

	ex, err := svc.Update(ctx, testUserID, testID, models.UpdateExerciseInput{
		Name: ptrStr(""),
	})
	assert.ErrorIs(t, err, service.ErrExerciseNameEmpty)
	assert.Nil(t, ex)
}

func TestExerciseService_Update_WhitespaceName(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{})

	ex, err := svc.Update(ctx, testUserID, testID, models.UpdateExerciseInput{
		Name: ptrStr("   "),
	})
	assert.ErrorIs(t, err, service.ErrExerciseNameEmpty)
	assert.Nil(t, ex)
}

func TestExerciseService_Update_InvalidWeight(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{})

	ex, err := svc.Update(ctx, testUserID, testID, models.UpdateExerciseInput{
		WorkingWeight: ptrFloat64(0),
	})
	assert.ErrorIs(t, err, service.ErrWeightInvalid)
	assert.Nil(t, ex)
}

func TestExerciseService_Update_NotFound(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{
		updateFn: func(ctx context.Context, userID string, id string, input models.UpdateExerciseInput) (*models.ExerciseRecord, error) {
			return nil, nil
		},
	})

	ex, err := svc.Update(ctx, testUserID, testID, models.UpdateExerciseInput{
		Name: ptrStr("New Name"),
	})
	assert.ErrorIs(t, err, service.ErrExerciseNotFound)
	assert.Nil(t, ex)
}

// ----- Archive -----

func TestExerciseService_Archive_Success(t *testing.T) {
	rec := baseRecord()
	rec.IsActive = false

	svc := service.NewExerciseService(&mockExerciseRepo{
		archiveFn: func(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error) {
			return rec, nil
		},
	})

	ex, err := svc.Archive(ctx, testUserID, testID)
	require.NoError(t, err)
	require.NotNil(t, ex)
	assert.False(t, ex.IsActive)
	assert.NotNil(t, ex.Media)
	assert.Empty(t, ex.Media)
}

func TestExerciseService_Archive_NotFound(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{
		archiveFn: func(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error) {
			return nil, nil
		},
	})

	ex, err := svc.Archive(ctx, testUserID, testID)
	assert.ErrorIs(t, err, service.ErrExerciseNotFound)
	assert.Nil(t, ex)
}

// ----- Restore -----

func TestExerciseService_Restore_Success(t *testing.T) {
	rec := baseRecord()
	rec.IsActive = true

	svc := service.NewExerciseService(&mockExerciseRepo{
		restoreFn: func(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error) {
			return rec, nil
		},
	})

	ex, err := svc.Restore(ctx, testUserID, testID)
	require.NoError(t, err)
	require.NotNil(t, ex)
	assert.True(t, ex.IsActive)
	assert.NotNil(t, ex.Media)
	assert.Empty(t, ex.Media)
}

func TestExerciseService_Restore_NotFound(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{
		restoreFn: func(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error) {
			return nil, nil
		},
	})

	ex, err := svc.Restore(ctx, testUserID, testID)
	assert.ErrorIs(t, err, service.ErrExerciseNotFound)
	assert.Nil(t, ex)
}

// ----- CreateMedia -----

func TestExerciseService_CreateMedia_Success(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{
		createMediaFn: func(ctx context.Context, userID string, exerciseID string, fileName string, filePath string, mimeType string, fileSize int64) (*models.ExerciseMedia, error) {
			return baseMedia(), nil
		},
	})

	media, err := svc.CreateMedia(ctx, testUserID, testID, "bench.mp4", "/uploads/bench.mp4", "video/mp4", 1024000)
	require.NoError(t, err)
	require.NotNil(t, media)
	assert.Equal(t, "bench.mp4", media.FileName)
}

// ----- GetMediaByID -----

func TestExerciseService_GetMediaByID_Success(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{
		getMediaByIDFn: func(ctx context.Context, userID string, id string) (*models.ExerciseMedia, error) {
			return baseMedia(), nil
		},
	})

	media, err := svc.GetMediaByID(ctx, testUserID, "770e8400-e29b-41d4-a716-446655440000")
	require.NoError(t, err)
	require.NotNil(t, media)
	assert.Equal(t, "bench.mp4", media.FileName)
}

func TestExerciseService_GetMediaByID_NotFound(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{
		getMediaByIDFn: func(ctx context.Context, userID string, id string) (*models.ExerciseMedia, error) {
			return nil, nil
		},
	})

	media, err := svc.GetMediaByID(ctx, testUserID, "nonexistent")
	assert.ErrorIs(t, err, service.ErrExerciseNotFound)
	assert.Nil(t, media)
}

// ----- DeleteMedia -----

func TestExerciseService_DeleteMedia_Success(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{
		deleteMediaFn: func(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error) {
			return &models.ExerciseMediaRecord{
				ID:       "770e8400-e29b-41d4-a716-446655440000",
				UserID:   userID,
				FileName: "bench.mp4",
				FilePath: "/uploads/bench.mp4",
			}, nil
		},
	})

	rec, err := svc.DeleteMedia(ctx, testUserID, "770e8400-e29b-41d4-a716-446655440000")
	require.NoError(t, err)
	require.NotNil(t, rec)
	assert.Equal(t, "bench.mp4", rec.FileName)
	assert.Equal(t, "/uploads/bench.mp4", rec.FilePath)
}

func TestExerciseService_DeleteMedia_NotFound(t *testing.T) {
	svc := service.NewExerciseService(&mockExerciseRepo{
		deleteMediaFn: func(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error) {
			return nil, nil
		},
	})

	rec, err := svc.DeleteMedia(ctx, testUserID, "nonexistent")
	assert.ErrorIs(t, err, service.ErrExerciseNotFound)
	assert.Nil(t, rec)
}

// ----- Media field always initialized -----

func TestExerciseService_MediaField_AlwaysInitialized(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		svc := service.NewExerciseService(&mockExerciseRepo{
			createFn: func(ctx context.Context, userID string, input models.CreateExerciseInput) (*models.ExerciseRecord, error) {
				return baseRecord(), nil
			},
		})
		ex, err := svc.Create(ctx, testUserID, models.CreateExerciseInput{Name: "Test"})
		require.NoError(t, err)
		require.NotNil(t, ex.Media)
	})

	t.Run("GetByID with no media", func(t *testing.T) {
		svc := service.NewExerciseService(&mockExerciseRepo{
			getByIDFn: func(ctx context.Context, userID string, id string) (*models.ExerciseRecord, error) {
				return baseRecord(), nil
			},
			listMediaByExerciseFn: func(ctx context.Context, userID string, exerciseID string) ([]models.ExerciseMedia, error) {
				return []models.ExerciseMedia{}, nil
			},
		})
		ex, err := svc.GetByID(ctx, testUserID, testID)
		require.NoError(t, err)
		require.NotNil(t, ex.Media)
	})

	t.Run("Update", func(t *testing.T) {
		svc := service.NewExerciseService(&mockExerciseRepo{
			updateFn: func(ctx context.Context, userID string, id string, input models.UpdateExerciseInput) (*models.ExerciseRecord, error) {
				return baseRecord(), nil
			},
		})
		ex, err := svc.Update(ctx, testUserID, testID, models.UpdateExerciseInput{Name: ptrStr("Updated")})
		require.NoError(t, err)
		require.NotNil(t, ex.Media)
	})

	t.Run("List", func(t *testing.T) {
		svc := service.NewExerciseService(&mockExerciseRepo{
			listFn: func(ctx context.Context, userID string, isActive bool, limit int32) ([]models.ExerciseRecord, error) {
				return []models.ExerciseRecord{*baseRecord()}, nil
			},
			countFn: func(ctx context.Context, userID string, isActive bool) (int, error) {
				return 1, nil
			},
		})
		conn, err := svc.List(ctx, testUserID, 20, nil, false)
		require.NoError(t, err)
		require.Len(t, conn.Items, 1)
		require.NotNil(t, conn.Items[0].Media)
	})
}