// FILE: apps/api/internal/atlas/service/user_profile_service.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement the transport-neutral UserProfileService for WAVE-07 AI Export user profile management.
//   SCOPE: Get (find by user ID), Update (upsert). Maps repository records to public models.
//   DEPENDS: apps/api/internal/atlas/repository/postgres.UserProfileRepository, apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-07.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   UserProfileService - Interface for user profile business operations.
//   NewUserProfileService - Creates a new UserProfileService.
//   Get - Gets a user profile by user ID.
//   Update - Creates or updates a user profile.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added user profile service for WAVE-07.
// END_CHANGE_SUMMARY

package service

import (
	"context"
	"errors"
	"fmt"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasRepo "monorepo-template/apps/api/internal/atlas/repository/postgres"
)

var (
	ErrUserProfileNotFound = errors.New("user profile not found")
)

type UserProfileService interface {
	Get(ctx context.Context, userID string) (*models.UserProfile, error)
	Update(ctx context.Context, userID string, input models.UserProfileInput) (*models.UserProfile, error)
}

type userProfileService struct {
	repo atlasRepo.UserProfileRepository
}

func NewUserProfileService(repo atlasRepo.UserProfileRepository) UserProfileService {
	return &userProfileService{repo: repo}
}

func (s *userProfileService) Get(ctx context.Context, userID string) (*models.UserProfile, error) {
	record, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user_profile_service.Get: %w", err)
	}
	if record == nil {
		return nil, ErrUserProfileNotFound
	}
	return models.UserProfileFromRecord(record), nil
}

func (s *userProfileService) Update(ctx context.Context, userID string, input models.UserProfileInput) (*models.UserProfile, error) {
	record, err := s.repo.Upsert(ctx, userID, input)
	if err != nil {
		return nil, fmt.Errorf("user_profile_service.Update: %w", err)
	}
	return models.UserProfileFromRecord(record), nil
}