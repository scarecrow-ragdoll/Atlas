// FILE: apps/api/internal/atlas/graph/resolver/user_profile_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Tests for WAVE-07 UserProfile and AiExport GraphQL resolvers.
//   SCOPE: user profile get/update, ai export get/list/delete.
//   DEPENDS: mock services for UserProfileService and AiExportService.
//   LINKS: M-API / V-M-API / WAVE-07.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT

package resolver_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasSvc "monorepo-template/apps/api/internal/atlas/service"

	"monorepo-template/apps/api/internal/atlas/graph/resolver"
)

type mockUserProfileSvc struct {
	atlasSvc.UserProfileService
	getFn    func(ctx context.Context, userID string) (*models.UserProfile, error)
	updateFn func(ctx context.Context, userID string, input models.UserProfileInput) (*models.UserProfile, error)
}

func (m *mockUserProfileSvc) Get(ctx context.Context, userID string) (*models.UserProfile, error) {
	return m.getFn(ctx, userID)
}

func (m *mockUserProfileSvc) Update(ctx context.Context, userID string, input models.UserProfileInput) (*models.UserProfile, error) {
	return m.updateFn(ctx, userID, input)
}

type mockAiExportSvc struct {
	atlasSvc.AiExportService
	generateFn func(ctx context.Context, userID string, input models.CreateAiExportInput, maxRangeDays int, maxExportSize int64, exportBasePath string) (*models.AiExport, string, error)
	getByIDFn  func(ctx context.Context, userID string, id string) (*models.AiExport, error)
	listFn     func(ctx context.Context, userID string) ([]models.AiExport, error)
	deleteFn   func(ctx context.Context, userID string, id string) (*models.AiExport, error)
}

func (m *mockAiExportSvc) Generate(ctx context.Context, userID string, input models.CreateAiExportInput, maxRangeDays int, maxExportSize int64, exportBasePath string) (*models.AiExport, string, error) {
	return m.generateFn(ctx, userID, input, maxRangeDays, maxExportSize, exportBasePath)
}

func (m *mockAiExportSvc) GetByID(ctx context.Context, userID string, id string) (*models.AiExport, error) {
	return m.getByIDFn(ctx, userID, id)
}

func (m *mockAiExportSvc) List(ctx context.Context, userID string) ([]models.AiExport, error) {
	return m.listFn(ctx, userID)
}

func (m *mockAiExportSvc) Delete(ctx context.Context, userID string, id string) (*models.AiExport, error) {
	return m.deleteFn(ctx, userID, id)
}

func TestUserProfileResolver_Get_HappyPath(t *testing.T) {
	goal := "Build muscle"
	r := &resolver.Resolver{
		UserProfileService: &mockUserProfileSvc{
			getFn: func(ctx context.Context, userID string) (*models.UserProfile, error) {
				assert.Equal(t, "test-uid", userID)
				return &models.UserProfile{
					ID:   "profile-1",
					Goal: &goal,
				}, nil
			},
		},
	}

	result, err := r.GetUserProfile(userCtx("test-uid"))
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Profile)
	assert.Equal(t, "profile-1", result.Profile.ID)
	assert.Equal(t, "Build muscle", *result.Profile.Goal)
}

func TestUserProfileResolver_Get_NotFound(t *testing.T) {
	r := &resolver.Resolver{
		UserProfileService: &mockUserProfileSvc{
			getFn: func(ctx context.Context, userID string) (*models.UserProfile, error) {
				return nil, atlasSvc.ErrUserProfileNotFound
			},
		},
	}

	result, err := r.GetUserProfile(userCtx("test-uid"))
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Profile)
}

func TestUserProfileResolver_Get_Unauthorized(t *testing.T) {
	r := &resolver.Resolver{
		UserProfileService: &mockUserProfileSvc{
			getFn: func(ctx context.Context, userID string) (*models.UserProfile, error) {
				t.Error("should not be called")
				return nil, nil
			},
		},
	}

	result, err := r.GetUserProfile(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, result.AuthErr)
	assert.Equal(t, models.UserProfileErrorAuth, result.AuthErr.Code)
}

func TestUserProfileResolver_Update_HappyPath(t *testing.T) {
	goal := "Lose weight"
	r := &resolver.Resolver{
		UserProfileService: &mockUserProfileSvc{
			updateFn: func(ctx context.Context, userID string, input models.UserProfileInput) (*models.UserProfile, error) {
				assert.Equal(t, "test-uid", userID)
				assert.Equal(t, "Lose weight", *input.Goal)
				return &models.UserProfile{
					ID:   "profile-1",
					Goal: &goal,
				}, nil
			},
		},
	}

	input := models.UserProfileInput{Goal: &goal}
	result, err := r.UpdateUserProfile(userCtx("test-uid"), input)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Profile)
	assert.Equal(t, "Lose weight", *result.Profile.Goal)
}

func TestAiExportResolver_GetByID_HappyPath(t *testing.T) {
	r := &resolver.Resolver{
		AiExportService: &mockAiExportSvc{
			getByIDFn: func(ctx context.Context, userID string, id string) (*models.AiExport, error) {
				assert.Equal(t, "test-uid", userID)
				assert.Equal(t, "export-1", id)
				return &models.AiExport{
					ID:             "export-1",
					GeneratedPrompt: "test prompt",
				}, nil
			},
		},
	}

	result, err := r.GetAiExport(userCtx("test-uid"), "export-1")
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Export)
	assert.Equal(t, "export-1", result.Export.ID)
	assert.Equal(t, "test prompt", result.Export.GeneratedPrompt)
}

func TestAiExportResolver_List_HappyPath(t *testing.T) {
	r := &resolver.Resolver{
		AiExportService: &mockAiExportSvc{
			listFn: func(ctx context.Context, userID string) ([]models.AiExport, error) {
				assert.Equal(t, "test-uid", userID)
				return []models.AiExport{
					{ID: "export-1", GeneratedPrompt: "prompt 1"},
					{ID: "export-2", GeneratedPrompt: "prompt 2"},
				}, nil
			},
		},
	}

	result, err := r.ListAiExports(userCtx("test-uid"))
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Exports, 2)
}

func TestAiExportResolver_Delete_HappyPath(t *testing.T) {
	r := &resolver.Resolver{
		AiExportService: &mockAiExportSvc{
			deleteFn: func(ctx context.Context, userID string, id string) (*models.AiExport, error) {
				assert.Equal(t, "test-uid", userID)
				assert.Equal(t, "export-1", id)
				return &models.AiExport{ID: "export-1"}, nil
			},
		},
	}

	result, err := r.DeleteAiExport(userCtx("test-uid"), "export-1")
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Export)
	assert.Equal(t, "export-1", result.Export.ID)
}