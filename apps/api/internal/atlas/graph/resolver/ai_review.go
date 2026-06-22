// FILE: apps/api/internal/atlas/graph/resolver/ai_review.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement GraphQL resolvers for WAVE-08 AiReview queries and mutations.
//   SCOPE: getAiReview, listAiReviews, createAiReview, updateAiReview, deleteAiReview.
//   DEPENDS: apps/api/internal/atlas/service.AiReviewService, apps/api/internal/atlas/middleware, apps/api/internal/atlas/models, generated ai_review.resolvers.go.
//   LINKS: M-API / V-M-API / WAVE-08.
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

func (r *Resolver) GetAiReview(ctx context.Context, id string) (*models.AiReviewResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.AiReviewResult{
			AuthErr: &models.AiReviewAuthErr{
				Message: "unauthorized",
				Code:    models.AiReviewErrorAuth,
			},
		}, nil
	}

	review, err := r.AiReviewService.GetByID(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrAiReviewNotFound) {
			return &models.AiReviewResult{
				NotFoundErr: &models.AiReviewNotFoundErr{
					Message: "ai review not found",
					Code:    models.AiReviewErrorNotFound,
				},
			}, nil
		}
		return nil, nil
	}

	return &models.AiReviewResult{Review: review}, nil
}

func (r *Resolver) ListAiReviews(ctx context.Context, dateRangeStart, dateRangeEnd *models.Date) (*models.AiReviewsResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.AiReviewsResult{
			AuthErr: &models.AiReviewAuthErr{
				Message: "unauthorized",
				Code:    models.AiReviewErrorAuth,
			},
		}, nil
	}

	reviews, err := r.AiReviewService.ListByUserIDAndDateRange(ctx, userID, dateRangeStart, dateRangeEnd)
	if err != nil {
		return nil, nil
	}

	return &models.AiReviewsResult{Reviews: reviews}, nil
}

func (r *Resolver) CreateAiReview(ctx context.Context, input models.CreateAiReviewInput) (*models.AiReviewResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.AiReviewResult{
			AuthErr: &models.AiReviewAuthErr{
				Message: "unauthorized",
				Code:    models.AiReviewErrorAuth,
			},
		}, nil
	}

	review, err := r.AiReviewService.Create(ctx, userID, input)
	if err != nil {
		if errors.Is(err, atlasService.ErrAiReviewEmptyText) || errors.Is(err, atlasService.ErrAiReviewInvalidDateRange) {
			return &models.AiReviewResult{
				ValidationErr: &models.AiReviewValidationErr{
					Message: err.Error(),
					Code:    models.AiReviewErrorValidation,
				},
			}, nil
		}
		return nil, nil
	}

	return &models.AiReviewResult{Review: review}, nil
}

func (r *Resolver) UpdateAiReview(ctx context.Context, id string, input models.UpdateAiReviewInput) (*models.AiReviewResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.AiReviewResult{
			AuthErr: &models.AiReviewAuthErr{
				Message: "unauthorized",
				Code:    models.AiReviewErrorAuth,
			},
		}, nil
	}

	review, err := r.AiReviewService.Update(ctx, userID, id, input)
	if err != nil {
		if errors.Is(err, atlasService.ErrAiReviewNotFound) {
			return &models.AiReviewResult{
				NotFoundErr: &models.AiReviewNotFoundErr{
					Message: "ai review not found",
					Code:    models.AiReviewErrorNotFound,
				},
			}, nil
		}
		if errors.Is(err, atlasService.ErrAiReviewEmptyText) || errors.Is(err, atlasService.ErrAiReviewInvalidDateRange) {
			return &models.AiReviewResult{
				ValidationErr: &models.AiReviewValidationErr{
					Message: err.Error(),
					Code:    models.AiReviewErrorValidation,
				},
			}, nil
		}
		return nil, nil
	}

	return &models.AiReviewResult{Review: review}, nil
}

func (r *Resolver) DeleteAiReview(ctx context.Context, id string) (*models.AiReviewResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return &models.AiReviewResult{
			AuthErr: &models.AiReviewAuthErr{
				Message: "unauthorized",
				Code:    models.AiReviewErrorAuth,
			},
		}, nil
	}

	review, err := r.AiReviewService.Delete(ctx, userID, id)
	if err != nil {
		if errors.Is(err, atlasService.ErrAiReviewNotFound) {
			return &models.AiReviewResult{
				NotFoundErr: &models.AiReviewNotFoundErr{
					Message: "ai review not found",
					Code:    models.AiReviewErrorNotFound,
				},
			}, nil
		}
		return nil, nil
	}

	return &models.AiReviewResult{Review: review}, nil
}