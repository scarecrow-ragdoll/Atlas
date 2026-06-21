// FILE: apps/api/internal/atlas/service/user_profile_service_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for UserProfileService covering Get and Update operations.
//   SCOPE: Success paths, not-found on Get, create-when-not-exists on Update.
//   DEPENDS: apps/api/internal/atlas/service, apps/api/internal/atlas/repository/postgres (mock UserProfileRepository), apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-07.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added user profile service unit tests for WAVE-07.
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

type mockUserProfileRepo struct {
	atlasPostgres.UserProfileRepository
	findByUserIDFn func(ctx context.Context, userID string) (*models.UserProfileRecord, error)
	upsertFn       func(ctx context.Context, userID string, input models.UserProfileInput) (*models.UserProfileRecord, error)
	createFn       func(ctx context.Context, userID string, input models.UserProfileInput) (*models.UserProfileRecord, error)
}

func (m *mockUserProfileRepo) FindByUserID(ctx context.Context, userID string) (*models.UserProfileRecord, error) {
	return m.findByUserIDFn(ctx, userID)
}

func (m *mockUserProfileRepo) Upsert(ctx context.Context, userID string, input models.UserProfileInput) (*models.UserProfileRecord, error) {
	return m.upsertFn(ctx, userID, input)
}

func (m *mockUserProfileRepo) Create(ctx context.Context, userID string, input models.UserProfileInput) (*models.UserProfileRecord, error) {
	return m.createFn(ctx, userID, input)
}

func testUserProfileRecord() *models.UserProfileRecord {
	goal := "Build muscle"
	height := 180.5
	birthDate := models.MustDate("1990-01-15")
	exp := "Intermediate"
	split := "Push Pull Legs"
	progression := "Double Progression"
	nutrition := "High Protein"
	aiContext := "Focus on progressive overload"
	return &models.UserProfileRecord{
		ID:                        testID,
		UserID:                    testUserID,
		Goal:                      &goal,
		Height:                    &height,
		BirthDate:                 &birthDate,
		TrainingExperience:        &exp,
		CurrentTrainingSplit:      &split,
		PreferredProgressionStyle: &progression,
		NutritionStrategy:         &nutrition,
		PersistentAiContext:       &aiContext,
		CreatedAt:                 "2026-06-21T00:00:00Z",
		UpdatedAt:                 "2026-06-21T00:00:00Z",
	}
}

// ----- Get -----

func TestUserProfileService_Get_Success(t *testing.T) {
	svc := service.NewUserProfileService(&mockUserProfileRepo{
		findByUserIDFn: func(ctx context.Context, userID string) (*models.UserProfileRecord, error) {
			return testUserProfileRecord(), nil
		},
	})

	profile, err := svc.Get(ctx, testUserID)
	require.NoError(t, err)
	require.NotNil(t, profile)
	assert.Equal(t, testUserID, profile.UserID)
	assert.Equal(t, "Build muscle", *profile.Goal)
	assert.Equal(t, 180.5, *profile.Height)
}

func TestUserProfileService_Get_NotFound(t *testing.T) {
	svc := service.NewUserProfileService(&mockUserProfileRepo{
		findByUserIDFn: func(ctx context.Context, userID string) (*models.UserProfileRecord, error) {
			return nil, nil
		},
	})

	profile, err := svc.Get(ctx, testUserID)
	assert.ErrorIs(t, err, service.ErrUserProfileNotFound)
	assert.Nil(t, profile)
}

// ----- Update -----

func TestUserProfileService_Update_Success(t *testing.T) {
	svc := service.NewUserProfileService(&mockUserProfileRepo{
		upsertFn: func(ctx context.Context, userID string, input models.UserProfileInput) (*models.UserProfileRecord, error) {
			assert.Equal(t, "Build muscle", *input.Goal)
			return testUserProfileRecord(), nil
		},
	})

	profile, err := svc.Update(ctx, testUserID, models.UserProfileInput{
		Goal: ptrStr("Build muscle"),
	})
	require.NoError(t, err)
	require.NotNil(t, profile)
	assert.Equal(t, testUserID, profile.UserID)
	assert.Equal(t, "Build muscle", *profile.Goal)
}

func TestUserProfileService_Update_CreatesWhenNotExists(t *testing.T) {
	svc := service.NewUserProfileService(&mockUserProfileRepo{
		upsertFn: func(ctx context.Context, userID string, input models.UserProfileInput) (*models.UserProfileRecord, error) {
			return testUserProfileRecord(), nil
		},
	})

	profile, err := svc.Update(ctx, testUserID, models.UserProfileInput{
		Goal: ptrStr("Build muscle"),
	})
	require.NoError(t, err)
	require.NotNil(t, profile)
	assert.Equal(t, testUserID, profile.UserID)
	assert.Equal(t, "Build muscle", *profile.Goal)
}
