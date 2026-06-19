// FILE: apps/api/internal/atlas/graph/resolver/exercise_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for exercise GraphQL resolvers.
//   SCOPE: Exercises, Exercise, AllExercises queries; CreateExercise, UpdateExercise, ArchiveExercise, RestoreExercise mutations. Covers happy paths, validation errors, not-found errors, auth errors, and empty results.
//   DEPENDS: apps/api/internal/atlas/graph/resolver, apps/api/internal/atlas/service, apps/api/internal/atlas/models, apps/api/internal/atlas/middleware.
//   LINKS: M-API / V-M-API / WAVE-02.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT

package resolver_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasSvc "monorepo-template/apps/api/internal/atlas/service"

	"monorepo-template/apps/api/internal/atlas/graph/resolver"
)

type mockExerciseService struct {
	createFn                 func(ctx context.Context, userID string, input models.CreateExerciseInput) (*models.Exercise, error)
	getByIDFn                func(ctx context.Context, userID string, id string) (*models.Exercise, error)
	listFn                   func(ctx context.Context, userID string, first int32, after *string, includeInactive bool) (*models.ExerciseConnection, error)
	listAllFn                func(ctx context.Context, userID string, includeInactive bool) ([]models.Exercise, error)
	updateFn                 func(ctx context.Context, userID string, id string, input models.UpdateExerciseInput) (*models.Exercise, error)
	archiveFn                func(ctx context.Context, userID string, id string) (*models.Exercise, error)
	restoreFn                func(ctx context.Context, userID string, id string) (*models.Exercise, error)
	createMediaFn            func(ctx context.Context, userID string, exerciseID string, fileName string, filePath string, mimeType string, fileSize int64) (*models.ExerciseMedia, error)
	getMediaByIDFn           func(ctx context.Context, userID string, id string) (*models.ExerciseMedia, error)
	getMediaRecordByIDFn     func(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error)
	deleteMediaFn            func(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error)
}

func (m *mockExerciseService) Create(ctx context.Context, userID string, input models.CreateExerciseInput) (*models.Exercise, error) {
	return m.createFn(ctx, userID, input)
}

func (m *mockExerciseService) GetByID(ctx context.Context, userID string, id string) (*models.Exercise, error) {
	return m.getByIDFn(ctx, userID, id)
}

func (m *mockExerciseService) List(ctx context.Context, userID string, first int32, after *string, includeInactive bool) (*models.ExerciseConnection, error) {
	return m.listFn(ctx, userID, first, after, includeInactive)
}

func (m *mockExerciseService) ListAll(ctx context.Context, userID string, includeInactive bool) ([]models.Exercise, error) {
	return m.listAllFn(ctx, userID, includeInactive)
}

func (m *mockExerciseService) Update(ctx context.Context, userID string, id string, input models.UpdateExerciseInput) (*models.Exercise, error) {
	return m.updateFn(ctx, userID, id, input)
}

func (m *mockExerciseService) Archive(ctx context.Context, userID string, id string) (*models.Exercise, error) {
	return m.archiveFn(ctx, userID, id)
}

func (m *mockExerciseService) Restore(ctx context.Context, userID string, id string) (*models.Exercise, error) {
	return m.restoreFn(ctx, userID, id)
}

func (m *mockExerciseService) CreateMedia(ctx context.Context, userID string, exerciseID string, fileName string, filePath string, mimeType string, fileSize int64) (*models.ExerciseMedia, error) {
	return m.createMediaFn(ctx, userID, exerciseID, fileName, filePath, mimeType, fileSize)
}

func (m *mockExerciseService) GetMediaByID(ctx context.Context, userID string, id string) (*models.ExerciseMedia, error) {
	return m.getMediaByIDFn(ctx, userID, id)
}

func (m *mockExerciseService) GetMediaRecordByID(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error) {
	return m.getMediaRecordByIDFn(ctx, userID, id)
}

func (m *mockExerciseService) DeleteMedia(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error) {
	return m.deleteMediaFn(ctx, userID, id)
}

// ---------------------------------------------------------------------------
// Exercises (paginated list)
// ---------------------------------------------------------------------------

func TestExercises_HappyPath_Defaults(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			listFn: func(ctx context.Context, userID string, first int32, after *string, includeInactive bool) (*models.ExerciseConnection, error) {
				assert.Equal(t, "test-uid", userID)
				assert.Equal(t, int32(20), first)
				assert.Nil(t, after)
				assert.False(t, includeInactive)
				return &models.ExerciseConnection{
					Items: []models.Exercise{
						{ID: "ex-1", Name: "Bench Press"},
						{ID: "ex-2", Name: "Squat"},
					},
					TotalCount: 2,
					PageInfo:   models.PageInfo{HasNextPage: false},
				}, nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.Exercises(ctx, nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Items, 2)
	assert.Equal(t, 2, result.TotalCount)
	assert.Equal(t, "Bench Press", result.Items[0].Name)
	assert.Equal(t, "Squat", result.Items[1].Name)
}

func TestExercises_HappyPath_WithPagination(t *testing.T) {
	after := "cursor-1"
	first := int(10)
	include := true

	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			listFn: func(ctx context.Context, userID string, first int32, after *string, includeInactive bool) (*models.ExerciseConnection, error) {
				assert.Equal(t, "test-uid", userID)
				assert.Equal(t, int32(10), first)
				require.NotNil(t, after)
				assert.Equal(t, "cursor-1", *after)
				assert.True(t, includeInactive)
				cursor := "cursor-1"
				return &models.ExerciseConnection{
					Items:      []models.Exercise{{ID: "ex-3", Name: "Deadlift"}},
					TotalCount: 1,
					PageInfo:   models.PageInfo{HasNextPage: false, EndCursor: &cursor},
				}, nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.Exercises(ctx, &first, &after, &include)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, "Deadlift", result.Items[0].Name)
}

func TestExercises_EmptyList(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			listFn: func(ctx context.Context, userID string, first int32, after *string, includeInactive bool) (*models.ExerciseConnection, error) {
				return &models.ExerciseConnection{
					Items:      []models.Exercise{},
					TotalCount: 0,
					PageInfo:   models.PageInfo{HasNextPage: false},
				}, nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.Exercises(ctx, nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Empty(t, result.Items)
	assert.Equal(t, 0, result.TotalCount)
}

func TestExercises_Unauthorized_ReturnsNil(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{},
	}

	result, err := r.Exercises(context.Background(), nil, nil, nil)
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestExercises_ServiceError_ReturnsNil(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			listFn: func(ctx context.Context, userID string, first int32, after *string, includeInactive bool) (*models.ExerciseConnection, error) {
				return nil, errors.New("db error")
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.Exercises(ctx, nil, nil, nil)
	require.NoError(t, err)
	assert.Nil(t, result)
}

// ---------------------------------------------------------------------------
// Exercise (single by ID)
// ---------------------------------------------------------------------------

func TestExercise_HappyPath(t *testing.T) {
	expected := &models.Exercise{ID: "ex-1", Name: "Bench Press"}
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			getByIDFn: func(ctx context.Context, userID string, id string) (*models.Exercise, error) {
				assert.Equal(t, "test-uid", userID)
				assert.Equal(t, "ex-1", id)
				return expected, nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.GetExercise(ctx, "ex-1")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expected, result.Exercise)
	assert.Nil(t, result.NotFoundErr)
	assert.Nil(t, result.AuthErr)
}

func TestExercise_NotFound(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			getByIDFn: func(ctx context.Context, userID string, id string) (*models.Exercise, error) {
				return nil, atlasSvc.ErrExerciseNotFound
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.GetExercise(ctx, "missing-id")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Exercise)
	assert.NotNil(t, result.NotFoundErr)
	assert.Equal(t, models.ExerciseErrorNotFound, result.NotFoundErr.Code)
}

func TestExercise_AuthError(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{},
	}

	result, err := r.GetExercise(context.Background(), "any-id")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Exercise)
	assert.NotNil(t, result.AuthErr)
	assert.Equal(t, models.ExerciseErrorAuth, result.AuthErr.Code)
}

func TestExercise_ServiceError_ReturnsNil(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			getByIDFn: func(ctx context.Context, userID string, id string) (*models.Exercise, error) {
				return nil, errors.New("unexpected")
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.GetExercise(ctx, "ex-1")
	require.NoError(t, err)
	assert.Nil(t, result)
}

// ---------------------------------------------------------------------------
// AllExercises (unpaginated)
// ---------------------------------------------------------------------------

func TestAllExercises_HappyPath(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			listAllFn: func(ctx context.Context, userID string, includeInactive bool) ([]models.Exercise, error) {
				assert.Equal(t, "test-uid", userID)
				assert.False(t, includeInactive)
				return []models.Exercise{
					{ID: "ex-1", Name: "Bench Press"},
					{ID: "ex-2", Name: "Squat"},
				}, nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.AllExercises(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "Bench Press", result[0].Name)
	assert.Equal(t, "Squat", result[1].Name)
}

func TestAllExercises_WithIncludeInactive(t *testing.T) {
	include := true
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			listAllFn: func(ctx context.Context, userID string, includeInactive bool) ([]models.Exercise, error) {
				assert.True(t, includeInactive)
				return []models.Exercise{
					{ID: "ex-1", Name: "Bench Press", IsActive: true},
					{ID: "ex-3", Name: "Inactive Ex", IsActive: false},
				}, nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.AllExercises(ctx, &include)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 2)
}

func TestAllExercises_EmptyList(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			listAllFn: func(ctx context.Context, userID string, includeInactive bool) ([]models.Exercise, error) {
				return []models.Exercise{}, nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.AllExercises(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Empty(t, result)
}

func TestAllExercises_Unauthorized_ReturnsNil(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{},
	}

	result, err := r.AllExercises(context.Background(), nil)
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestAllExercises_ServiceError_ReturnsNil(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			listAllFn: func(ctx context.Context, userID string, includeInactive bool) ([]models.Exercise, error) {
				return nil, errors.New("db error")
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.AllExercises(ctx, nil)
	require.NoError(t, err)
	assert.Nil(t, result)
}

// ---------------------------------------------------------------------------
// CreateExercise
// ---------------------------------------------------------------------------

func TestCreateExercise_HappyPath(t *testing.T) {
	input := models.CreateExerciseInput{
		Name:         "Bench Press",
		MuscleGroups: []string{"chest"},
		WorkingWeight: ptrFloat64(100.0),
	}
	expected := &models.Exercise{ID: "ex-new", Name: "Bench Press", WorkingWeight: ptrFloat64(100.0)}

	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			createFn: func(ctx context.Context, userID string, inp models.CreateExerciseInput) (*models.Exercise, error) {
				assert.Equal(t, "test-uid", userID)
				assert.Equal(t, input, inp)
				return expected, nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.CreateExercise(ctx, input)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expected, result.Exercise)
	assert.Nil(t, result.ValidationErr)
	assert.Nil(t, result.AuthErr)
}

func TestCreateExercise_ValidationError_EmptyName(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			createFn: func(ctx context.Context, userID string, input models.CreateExerciseInput) (*models.Exercise, error) {
				return nil, atlasSvc.ErrExerciseNameEmpty
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.CreateExercise(ctx, models.CreateExerciseInput{Name: ""})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Exercise)
	assert.NotNil(t, result.ValidationErr)
	assert.Equal(t, models.ExerciseErrorValidation, result.ValidationErr.Code)
}

func TestCreateExercise_ValidationError_InvalidWeight(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			createFn: func(ctx context.Context, userID string, input models.CreateExerciseInput) (*models.Exercise, error) {
				return nil, atlasSvc.ErrWeightInvalid
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.CreateExercise(ctx, models.CreateExerciseInput{Name: "Squat", WorkingWeight: ptrFloat64(-5)})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Exercise)
	assert.NotNil(t, result.ValidationErr)
	assert.Equal(t, models.ExerciseErrorValidation, result.ValidationErr.Code)
}

func TestCreateExercise_AuthError(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{},
	}

	result, err := r.CreateExercise(context.Background(), models.CreateExerciseInput{Name: "Test"})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Exercise)
	assert.NotNil(t, result.AuthErr)
	assert.Equal(t, models.ExerciseErrorAuth, result.AuthErr.Code)
}

func TestCreateExercise_ServiceError_ReturnsNil(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			createFn: func(ctx context.Context, userID string, input models.CreateExerciseInput) (*models.Exercise, error) {
				return nil, errors.New("unexpected error")
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.CreateExercise(ctx, models.CreateExerciseInput{Name: "Test"})
	require.NoError(t, err)
	assert.Nil(t, result)
}

// ---------------------------------------------------------------------------
// UpdateExercise
// ---------------------------------------------------------------------------

func TestUpdateExercise_HappyPath(t *testing.T) {
	newName := "Updated Bench"
	input := models.UpdateExerciseInput{Name: &newName}
	expected := &models.Exercise{ID: "ex-1", Name: "Updated Bench"}

	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			updateFn: func(ctx context.Context, userID string, id string, inp models.UpdateExerciseInput) (*models.Exercise, error) {
				assert.Equal(t, "test-uid", userID)
				assert.Equal(t, "ex-1", id)
				assert.Equal(t, input, inp)
				return expected, nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.UpdateExercise(ctx, "ex-1", input)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expected, result.Exercise)
	assert.Nil(t, result.ValidationErr)
	assert.Nil(t, result.NotFoundErr)
	assert.Nil(t, result.AuthErr)
}

func TestUpdateExercise_ValidationError_EmptyName(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			updateFn: func(ctx context.Context, userID string, id string, input models.UpdateExerciseInput) (*models.Exercise, error) {
				return nil, atlasSvc.ErrExerciseNameEmpty
			},
		},
	}

	ctx := userCtx("test-uid")
	emptyName := ""
	result, err := r.UpdateExercise(ctx, "ex-1", models.UpdateExerciseInput{Name: &emptyName})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Exercise)
	assert.NotNil(t, result.ValidationErr)
	assert.Equal(t, models.ExerciseErrorValidation, result.ValidationErr.Code)
}

func TestUpdateExercise_ValidationError_InvalidWeight(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			updateFn: func(ctx context.Context, userID string, id string, input models.UpdateExerciseInput) (*models.Exercise, error) {
				return nil, atlasSvc.ErrWeightInvalid
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.UpdateExercise(ctx, "ex-1", models.UpdateExerciseInput{WorkingWeight: ptrFloat64(-10)})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Exercise)
	assert.NotNil(t, result.ValidationErr)
	assert.Equal(t, models.ExerciseErrorValidation, result.ValidationErr.Code)
}

func TestUpdateExercise_NotFound(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			updateFn: func(ctx context.Context, userID string, id string, input models.UpdateExerciseInput) (*models.Exercise, error) {
				return nil, atlasSvc.ErrExerciseNotFound
			},
		},
	}

	ctx := userCtx("test-uid")
	newName := "Ghost"
	result, err := r.UpdateExercise(ctx, "missing", models.UpdateExerciseInput{Name: &newName})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Exercise)
	assert.NotNil(t, result.NotFoundErr)
	assert.Equal(t, models.ExerciseErrorNotFound, result.NotFoundErr.Code)
}

func TestUpdateExercise_AuthError(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{},
	}

	result, err := r.UpdateExercise(context.Background(), "ex-1", models.UpdateExerciseInput{})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Exercise)
	assert.NotNil(t, result.AuthErr)
	assert.Equal(t, models.ExerciseErrorAuth, result.AuthErr.Code)
}

func TestUpdateExercise_ServiceError_ReturnsNil(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			updateFn: func(ctx context.Context, userID string, id string, input models.UpdateExerciseInput) (*models.Exercise, error) {
				return nil, errors.New("unexpected")
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.UpdateExercise(ctx, "ex-1", models.UpdateExerciseInput{})
	require.NoError(t, err)
	assert.Nil(t, result)
}

// ---------------------------------------------------------------------------
// ArchiveExercise
// ---------------------------------------------------------------------------

func TestArchiveExercise_HappyPath(t *testing.T) {
	expected := &models.Exercise{ID: "ex-1", Name: "Bench Press", IsActive: false}

	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			archiveFn: func(ctx context.Context, userID string, id string) (*models.Exercise, error) {
				assert.Equal(t, "test-uid", userID)
				assert.Equal(t, "ex-1", id)
				return expected, nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.ArchiveExercise(ctx, "ex-1")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expected, result.Exercise)
	assert.Nil(t, result.NotFoundErr)
	assert.Nil(t, result.AuthErr)
}

func TestArchiveExercise_NotFound(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			archiveFn: func(ctx context.Context, userID string, id string) (*models.Exercise, error) {
				return nil, atlasSvc.ErrExerciseNotFound
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.ArchiveExercise(ctx, "missing-id")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Exercise)
	assert.NotNil(t, result.NotFoundErr)
	assert.Equal(t, models.ExerciseErrorNotFound, result.NotFoundErr.Code)
}

func TestArchiveExercise_AuthError(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{},
	}

	result, err := r.ArchiveExercise(context.Background(), "ex-1")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Exercise)
	assert.NotNil(t, result.AuthErr)
	assert.Equal(t, models.ExerciseErrorAuth, result.AuthErr.Code)
}

func TestArchiveExercise_ServiceError_ReturnsNil(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			archiveFn: func(ctx context.Context, userID string, id string) (*models.Exercise, error) {
				return nil, errors.New("unexpected")
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.ArchiveExercise(ctx, "ex-1")
	require.NoError(t, err)
	assert.Nil(t, result)
}

// ---------------------------------------------------------------------------
// RestoreExercise
// ---------------------------------------------------------------------------

func TestRestoreExercise_HappyPath(t *testing.T) {
	expected := &models.Exercise{ID: "ex-1", Name: "Bench Press", IsActive: true}

	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			restoreFn: func(ctx context.Context, userID string, id string) (*models.Exercise, error) {
				assert.Equal(t, "test-uid", userID)
				assert.Equal(t, "ex-1", id)
				return expected, nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.RestoreExercise(ctx, "ex-1")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expected, result.Exercise)
	assert.Nil(t, result.NotFoundErr)
	assert.Nil(t, result.AuthErr)
}

func TestRestoreExercise_NotFound(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			restoreFn: func(ctx context.Context, userID string, id string) (*models.Exercise, error) {
				return nil, atlasSvc.ErrExerciseNotFound
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.RestoreExercise(ctx, "missing-id")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Exercise)
	assert.NotNil(t, result.NotFoundErr)
	assert.Equal(t, models.ExerciseErrorNotFound, result.NotFoundErr.Code)
}

func TestRestoreExercise_AuthError(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{},
	}

	result, err := r.RestoreExercise(context.Background(), "ex-1")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Exercise)
	assert.NotNil(t, result.AuthErr)
	assert.Equal(t, models.ExerciseErrorAuth, result.AuthErr.Code)
}

func TestRestoreExercise_ServiceError_ReturnsNil(t *testing.T) {
	r := &resolver.Resolver{
		ExerciseService: &mockExerciseService{
			restoreFn: func(ctx context.Context, userID string, id string) (*models.Exercise, error) {
				return nil, errors.New("unexpected")
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.RestoreExercise(ctx, "ex-1")
	require.NoError(t, err)
	assert.Nil(t, result)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func ptrFloat64(v float64) *float64 {
	return &v
}