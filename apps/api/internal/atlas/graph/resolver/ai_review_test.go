// FILE: apps/api/internal/atlas/graph/resolver/ai_review_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for WAVE-08 AiReview GraphQL resolvers.
//   SCOPE: createAiReview, aiReviews (list), deleteAiReview mutations/queries. Covers happy paths, validation errors, not-found errors, auth errors.
//   DEPENDS: apps/api/internal/atlas/graph/resolver, apps/api/internal/atlas/service, apps/api/internal/atlas/models, apps/api/internal/atlas/middleware.
//   LINKS: M-API / V-M-API / WAVE-08.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added AI review resolver unit tests for WAVE-08.
// END_CHANGE_SUMMARY

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

type mockAiReviewService struct {
	createFn                  func(ctx context.Context, userID string, input models.CreateAiReviewInput) (*models.AiReview, error)
	getByIDFn                 func(ctx context.Context, userID string, id string) (*models.AiReview, error)
	listByUserIDFn            func(ctx context.Context, userID string) ([]models.AiReview, error)
	listByUserIDAndDateRangeFn func(ctx context.Context, userID string, start, end *models.Date) ([]models.AiReview, error)
	updateFn                  func(ctx context.Context, userID string, id string, input models.UpdateAiReviewInput) (*models.AiReview, error)
	deleteFn                  func(ctx context.Context, userID string, id string) (*models.AiReview, error)
	listAllByUserIDFn         func(ctx context.Context, userID string) ([]models.AiReview, error)
}

func (m *mockAiReviewService) Create(ctx context.Context, userID string, input models.CreateAiReviewInput) (*models.AiReview, error) {
	return m.createFn(ctx, userID, input)
}

func (m *mockAiReviewService) GetByID(ctx context.Context, userID string, id string) (*models.AiReview, error) {
	return m.getByIDFn(ctx, userID, id)
}

func (m *mockAiReviewService) ListByUserID(ctx context.Context, userID string) ([]models.AiReview, error) {
	return m.listByUserIDFn(ctx, userID)
}

func (m *mockAiReviewService) ListByUserIDAndDateRange(ctx context.Context, userID string, start, end *models.Date) ([]models.AiReview, error) {
	return m.listByUserIDAndDateRangeFn(ctx, userID, start, end)
}

func (m *mockAiReviewService) Update(ctx context.Context, userID string, id string, input models.UpdateAiReviewInput) (*models.AiReview, error) {
	return m.updateFn(ctx, userID, id, input)
}

func (m *mockAiReviewService) Delete(ctx context.Context, userID string, id string) (*models.AiReview, error) {
	return m.deleteFn(ctx, userID, id)
}

func (m *mockAiReviewService) ListAllByUserID(ctx context.Context, userID string) ([]models.AiReview, error) {
	return m.listAllByUserIDFn(ctx, userID)
}

var (
	arReviewID    = "660e8400-e29b-41d4-a716-446655440001"
	arDateStart   = models.MustDate("2026-06-01")
	arDateEnd     = models.MustDate("2026-06-30")
	arTestReviews = []models.AiReview{
		{ID: arReviewID, DateRangeStart: arDateStart, DateRangeEnd: arDateEnd, AiResponseText: "Test analysis"},
	}
)

// ----- CreateAiReview -----

func TestAiReviewResolver_Create_Success(t *testing.T) {
	r := &resolver.Resolver{
		AiReviewService: &mockAiReviewService{
			createFn: func(ctx context.Context, userID string, input models.CreateAiReviewInput) (*models.AiReview, error) {
				assert.Equal(t, "test-uid", userID)
				return &models.AiReview{ID: arReviewID, DateRangeStart: input.DateRangeStart, DateRangeEnd: input.DateRangeEnd, AiResponseText: input.AiResponseText}, nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.CreateAiReview(ctx, models.CreateAiReviewInput{
		DateRangeStart: arDateStart,
		DateRangeEnd:   arDateEnd,
		AiResponseText: "Test analysis",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Review)
	assert.Equal(t, arReviewID, result.Review.ID)
	assert.Nil(t, result.ValidationErr)
	assert.Nil(t, result.AuthErr)
}

func TestAiReviewResolver_Create_ValidationError(t *testing.T) {
	r := &resolver.Resolver{
		AiReviewService: &mockAiReviewService{
			createFn: func(ctx context.Context, userID string, input models.CreateAiReviewInput) (*models.AiReview, error) {
				return nil, atlasSvc.ErrAiReviewEmptyText
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.CreateAiReview(ctx, models.CreateAiReviewInput{DateRangeStart: arDateStart, DateRangeEnd: arDateEnd, AiResponseText: ""})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Review)
	assert.NotNil(t, result.ValidationErr)
	assert.Equal(t, models.AiReviewErrorValidation, result.ValidationErr.Code)
}

func TestAiReviewResolver_Create_AuthError(t *testing.T) {
	r := &resolver.Resolver{
		AiReviewService: &mockAiReviewService{},
	}

	result, err := r.CreateAiReview(context.Background(), models.CreateAiReviewInput{DateRangeStart: arDateStart, DateRangeEnd: arDateEnd, AiResponseText: "Test"})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Review)
	assert.NotNil(t, result.AuthErr)
	assert.Equal(t, models.AiReviewErrorAuth, result.AuthErr.Code)
}

func TestAiReviewResolver_Create_ServiceError_ReturnsNil(t *testing.T) {
	r := &resolver.Resolver{
		AiReviewService: &mockAiReviewService{
			createFn: func(ctx context.Context, userID string, input models.CreateAiReviewInput) (*models.AiReview, error) {
				return nil, errors.New("unexpected")
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.CreateAiReview(ctx, models.CreateAiReviewInput{DateRangeStart: arDateStart, DateRangeEnd: arDateEnd, AiResponseText: "Test"})
	require.NoError(t, err)
	assert.Nil(t, result)
}

// ----- ListAiReviews -----

func TestAiReviewResolver_List_Success(t *testing.T) {
	r := &resolver.Resolver{
		AiReviewService: &mockAiReviewService{
			listByUserIDAndDateRangeFn: func(ctx context.Context, userID string, start, end *models.Date) ([]models.AiReview, error) {
				assert.Equal(t, "test-uid", userID)
				return arTestReviews, nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.ListAiReviews(ctx, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Reviews, 1)
	assert.Nil(t, result.AuthErr)
}

func TestAiReviewResolver_List_Empty(t *testing.T) {
	r := &resolver.Resolver{
		AiReviewService: &mockAiReviewService{
			listByUserIDAndDateRangeFn: func(ctx context.Context, userID string, start, end *models.Date) ([]models.AiReview, error) {
				return []models.AiReview{}, nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.ListAiReviews(ctx, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Empty(t, result.Reviews)
}

func TestAiReviewResolver_List_AuthError(t *testing.T) {
	r := &resolver.Resolver{
		AiReviewService: &mockAiReviewService{},
	}

	result, err := r.ListAiReviews(context.Background(), nil, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, result.AuthErr)
	assert.Equal(t, models.AiReviewErrorAuth, result.AuthErr.Code)
}

func TestAiReviewResolver_List_ServiceError_ReturnsNil(t *testing.T) {
	r := &resolver.Resolver{
		AiReviewService: &mockAiReviewService{
			listByUserIDAndDateRangeFn: func(ctx context.Context, userID string, start, end *models.Date) ([]models.AiReview, error) {
				return nil, errors.New("db error")
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.ListAiReviews(ctx, nil, nil)
	require.NoError(t, err)
	assert.Nil(t, result)
}

// ----- GetAiReview -----

func TestAiReviewResolver_Get_Success(t *testing.T) {
	r := &resolver.Resolver{
		AiReviewService: &mockAiReviewService{
			getByIDFn: func(ctx context.Context, userID string, id string) (*models.AiReview, error) {
				assert.Equal(t, "test-uid", userID)
				assert.Equal(t, arReviewID, id)
				return &arTestReviews[0], nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.GetAiReview(ctx, arReviewID)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Review)
	assert.Equal(t, arReviewID, result.Review.ID)
}

func TestAiReviewResolver_Get_NotFound(t *testing.T) {
	r := &resolver.Resolver{
		AiReviewService: &mockAiReviewService{
			getByIDFn: func(ctx context.Context, userID string, id string) (*models.AiReview, error) {
				return nil, atlasSvc.ErrAiReviewNotFound
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.GetAiReview(ctx, "missing-id")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Review)
	assert.NotNil(t, result.NotFoundErr)
}

func TestAiReviewResolver_Get_AuthError(t *testing.T) {
	r := &resolver.Resolver{
		AiReviewService: &mockAiReviewService{},
	}

	result, err := r.GetAiReview(context.Background(), arReviewID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, result.AuthErr)
}

// ----- UpdateAiReview -----

func TestAiReviewResolver_Update_Success(t *testing.T) {
	newNotes := "Updated notes"
	r := &resolver.Resolver{
		AiReviewService: &mockAiReviewService{
			updateFn: func(ctx context.Context, userID string, id string, input models.UpdateAiReviewInput) (*models.AiReview, error) {
				assert.Equal(t, "test-uid", userID)
				assert.Equal(t, arReviewID, id)
				return &models.AiReview{ID: arReviewID, UserNotes: &newNotes}, nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.UpdateAiReview(ctx, arReviewID, models.UpdateAiReviewInput{UserNotes: &newNotes})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Review)
	assert.Equal(t, "Updated notes", *result.Review.UserNotes)
}

func TestAiReviewResolver_Update_NotFound(t *testing.T) {
	r := &resolver.Resolver{
		AiReviewService: &mockAiReviewService{
			updateFn: func(ctx context.Context, userID string, id string, input models.UpdateAiReviewInput) (*models.AiReview, error) {
				return nil, atlasSvc.ErrAiReviewNotFound
			},
		},
	}

	ctx := userCtx("test-uid")
	newNotes := "Ghost"
	result, err := r.UpdateAiReview(ctx, "missing", models.UpdateAiReviewInput{UserNotes: &newNotes})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Review)
	assert.NotNil(t, result.NotFoundErr)
}

func TestAiReviewResolver_Update_AuthError(t *testing.T) {
	r := &resolver.Resolver{
		AiReviewService: &mockAiReviewService{},
	}

	result, err := r.UpdateAiReview(context.Background(), arReviewID, models.UpdateAiReviewInput{})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, result.AuthErr)
}

// ----- DeleteAiReview -----

func TestAiReviewResolver_Delete_Success(t *testing.T) {
	r := &resolver.Resolver{
		AiReviewService: &mockAiReviewService{
			deleteFn: func(ctx context.Context, userID string, id string) (*models.AiReview, error) {
				assert.Equal(t, "test-uid", userID)
				assert.Equal(t, arReviewID, id)
				return &arTestReviews[0], nil
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.DeleteAiReview(ctx, arReviewID)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Review)
	assert.Equal(t, arReviewID, result.Review.ID)
}

func TestAiReviewResolver_Delete_NotFound(t *testing.T) {
	r := &resolver.Resolver{
		AiReviewService: &mockAiReviewService{
			deleteFn: func(ctx context.Context, userID string, id string) (*models.AiReview, error) {
				return nil, atlasSvc.ErrAiReviewNotFound
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.DeleteAiReview(ctx, "missing-id")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Review)
	assert.NotNil(t, result.NotFoundErr)
}

func TestAiReviewResolver_Delete_AuthError(t *testing.T) {
	r := &resolver.Resolver{
		AiReviewService: &mockAiReviewService{},
	}

	result, err := r.DeleteAiReview(context.Background(), arReviewID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, result.AuthErr)
}

func TestAiReviewResolver_Delete_ServiceError_ReturnsNil(t *testing.T) {
	r := &resolver.Resolver{
		AiReviewService: &mockAiReviewService{
			deleteFn: func(ctx context.Context, userID string, id string) (*models.AiReview, error) {
				return nil, errors.New("unexpected")
			},
		},
	}

	ctx := userCtx("test-uid")
	result, err := r.DeleteAiReview(ctx, arReviewID)
	require.NoError(t, err)
	assert.Nil(t, result)
}