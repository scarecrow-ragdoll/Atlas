// FILE: apps/api/internal/atlas/graph/resolver/user_profile.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement GraphQL resolvers for WAVE-07 UserProfile queries and mutations.
//   SCOPE: getUserProfile, updateUserProfile.
//   DEPENDS: apps/api/internal/atlas/service.UserProfileService, apps/api/internal/atlas/middleware, apps/api/internal/atlas/models, generated user_profile.resolvers.go.
//   LINKS: M-API / V-M-API / WAVE-07.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT

package resolver

import (
	"context"
	"errors"

	"monorepo-template/apps/api/internal/atlas/middleware"
	"monorepo-template/apps/api/internal/atlas/models"
	atlasService "monorepo-template/apps/api/internal/atlas/service"
)

func (r *Resolver) GetUserProfile(ctx context.Context) (*models.UserProfileResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.UserProfileResult{
			AuthErr: &models.UserProfileAuthErr{
				Message: "unauthorized",
				Code:    models.UserProfileErrorAuth,
			},
		}, nil
	}

	profile, err := r.UserProfileService.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, atlasService.ErrUserProfileNotFound) {
			return &models.UserProfileResult{}, nil
		}
		return nil, nil
	}

	return &models.UserProfileResult{Profile: profile}, nil
}

func (r *Resolver) UpdateUserProfile(ctx context.Context, input models.UserProfileInput) (*models.UserProfileResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.UserProfileResult{
			AuthErr: &models.UserProfileAuthErr{
				Message: "unauthorized",
				Code:    models.UserProfileErrorAuth,
			},
		}, nil
	}

	profile, err := r.UserProfileService.Update(ctx, userID, input)
	if err != nil {
		return nil, nil
	}

	return &models.UserProfileResult{Profile: profile}, nil
}