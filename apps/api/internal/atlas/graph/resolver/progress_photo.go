// FILE: apps/api/internal/atlas/graph/resolver/progress_photo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement GraphQL resolver for WAVE-04 ProgressPhotos query.
//   SCOPE: progressPhotos(checkInID). Returns empty photos list when no service is available.
//   DEPENDS: apps/api/internal/atlas/middleware, apps/api/internal/atlas/models, generated progress_photo.resolvers.go.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT

package resolver

import (
	"context"

	"monorepo-template/apps/api/internal/atlas/middleware"
	"monorepo-template/apps/api/internal/atlas/models"
)

func (r *Resolver) GetProgressPhotos(ctx context.Context, checkInID string) (*models.ProgressPhotosResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.ProgressPhotosResult{
			AuthErr: &models.ProgressPhotoAuthErr{
				Message: "unauthorized",
				Code:    models.ProgressPhotoErrorAuth,
			},
		}, nil
	}

	_ = checkInID
	_ = userID

	return &models.ProgressPhotosResult{Photos: []models.ProgressPhoto{}}, nil
}
