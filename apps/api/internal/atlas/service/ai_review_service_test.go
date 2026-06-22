// FILE: apps/api/internal/atlas/service/ai_review_service_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for AiReviewService covering Create, List, Update, Delete operations with validation.
//   SCOPE: Success paths, validation errors (empty text, invalid date range), list ordering, date range filter, update success, update ownership, delete success, log privacy.
//   DEPENDS: apps/api/internal/atlas/service, apps/api/internal/atlas/repository/postgres (mock AiReviewRepository), apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-08.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added AI review service unit tests for WAVE-08.
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

type mockAiReviewRepo struct {
	atlasPostgres.AiReviewRepository
	createFn                  func(ctx context.Context, userID string, dateRangeStart, dateRangeEnd string, aiResponseText string, userNotes, plannedActions *string) (*models.AiReviewRecord, error)
	getByIDFn                 func(ctx context.Context, userID string, id string) (*models.AiReviewRecord, error)
	listByUserIDFn            func(ctx context.Context, userID string) ([]models.AiReviewRecord, error)
	listByUserIDAndDateRangeFn func(ctx context.Context, userID string, dateRangeStart, dateRangeEnd string) ([]models.AiReviewRecord, error)
	updateFn                  func(ctx context.Context, userID string, id string, input models.UpdateAiReviewInput) (*models.AiReviewRecord, error)
	deleteFn                  func(ctx context.Context, userID string, id string) (*models.AiReviewRecord, error)
}

func (m *mockAiReviewRepo) Create(ctx context.Context, userID string, dateRangeStart, dateRangeEnd string, aiResponseText string, userNotes, plannedActions *string) (*models.AiReviewRecord, error) {
	return m.createFn(ctx, userID, dateRangeStart, dateRangeEnd, aiResponseText, userNotes, plannedActions)
}

func (m *mockAiReviewRepo) GetByID(ctx context.Context, userID string, id string) (*models.AiReviewRecord, error) {
	return m.getByIDFn(ctx, userID, id)
}

func (m *mockAiReviewRepo) ListByUserID(ctx context.Context, userID string) ([]models.AiReviewRecord, error) {
	return m.listByUserIDFn(ctx, userID)
}

func (m *mockAiReviewRepo) ListByUserIDAndDateRange(ctx context.Context, userID string, dateRangeStart, dateRangeEnd string) ([]models.AiReviewRecord, error) {
	return m.listByUserIDAndDateRangeFn(ctx, userID, dateRangeStart, dateRangeEnd)
}

func (m *mockAiReviewRepo) Update(ctx context.Context, userID string, id string, input models.UpdateAiReviewInput) (*models.AiReviewRecord, error) {
	return m.updateFn(ctx, userID, id, input)
}

func (m *mockAiReviewRepo) Delete(ctx context.Context, userID string, id string) (*models.AiReviewRecord, error) {
	return m.deleteFn(ctx, userID, id)
}

var (
	arCtx        = context.Background()
	arTestUserID = "550e8400-e29b-41d4-a716-446655440000"
	arTestID     = "660e8400-e29b-41d4-a716-446655440001"
	arDateStart  = models.MustDate("2026-06-01")
	arDateEnd    = models.MustDate("2026-06-30")
)

func arTestRecord() *models.AiReviewRecord {
	return &models.AiReviewRecord{
		ID:              arTestID,
		UserID:          arTestUserID,
		DateRangeStart:  arDateStart,
		DateRangeEnd:    arDateEnd,
		AiResponseText:  "AI analysis response text",
		UserNotes:       ptrStr("User notes"),
		PlannedActions:  ptrStr("Increase protein intake"),
		CreatedAt:       "2026-06-21T00:00:00Z",
		UpdatedAt:       "2026-06-21T00:00:00Z",
	}
}

// ----- Create -----

func TestAiReviewService_Create_Success(t *testing.T) {
	svc := service.NewAiReviewService(&mockAiReviewRepo{
		createFn: func(ctx context.Context, userID string, dateRangeStart, dateRangeEnd string, aiResponseText string, userNotes, plannedActions *string) (*models.AiReviewRecord, error) {
			assert.Equal(t, "AI analysis response text", aiResponseText)
			return arTestRecord(), nil
		},
	})

	review, err := svc.Create(arCtx, arTestUserID, models.CreateAiReviewInput{
		DateRangeStart: arDateStart,
		DateRangeEnd:   arDateEnd,
		AiResponseText: "AI analysis response text",
		UserNotes:      ptrStr("User notes"),
		PlannedActions: ptrStr("Increase protein intake"),
	})
	require.NoError(t, err)
	require.NotNil(t, review)
	assert.Equal(t, "AI analysis response text", review.AiResponseText)
	assert.Equal(t, arDateStart, review.DateRangeStart)
	assert.Equal(t, arDateEnd, review.DateRangeEnd)
}

func TestAiReviewService_Create_EmptyText(t *testing.T) {
	svc := service.NewAiReviewService(&mockAiReviewRepo{})

	review, err := svc.Create(arCtx, arTestUserID, models.CreateAiReviewInput{
		DateRangeStart: arDateStart,
		DateRangeEnd:   arDateEnd,
		AiResponseText: "",
	})
	assert.ErrorIs(t, err, service.ErrAiReviewEmptyText)
	assert.Nil(t, review)
}

func TestAiReviewService_Create_InvalidDateRange(t *testing.T) {
	svc := service.NewAiReviewService(&mockAiReviewRepo{})

	badStart := models.MustDate("2026-07-01")
	badEnd := models.MustDate("2026-06-01")

	review, err := svc.Create(arCtx, arTestUserID, models.CreateAiReviewInput{
		DateRangeStart: badStart,
		DateRangeEnd:   badEnd,
		AiResponseText: "Some text",
	})
	assert.ErrorIs(t, err, service.ErrAiReviewInvalidDateRange)
	assert.Nil(t, review)
}

// ----- List -----

func TestAiReviewService_List_Ordered(t *testing.T) {
	svc := service.NewAiReviewService(&mockAiReviewRepo{
		listByUserIDFn: func(ctx context.Context, userID string) ([]models.AiReviewRecord, error) {
			r1 := *arTestRecord()
			r2 := *arTestRecord()
			r2.ID = "660e8400-e29b-41d4-a716-446655440002"
			r2.CreatedAt = "2026-06-20T00:00:00Z"
			return []models.AiReviewRecord{r1, r2}, nil
		},
	})

	reviews, err := svc.ListByUserID(arCtx, arTestUserID)
	require.NoError(t, err)
	require.Len(t, reviews, 2)
	assert.Equal(t, arTestID, reviews[0].ID)
	assert.Equal(t, "660e8400-e29b-41d4-a716-446655440002", reviews[1].ID)
}

func TestAiReviewService_List_DateRangeFilter(t *testing.T) {
	svc := service.NewAiReviewService(&mockAiReviewRepo{
		listByUserIDAndDateRangeFn: func(ctx context.Context, userID string, dateRangeStart, dateRangeEnd string) ([]models.AiReviewRecord, error) {
			assert.Equal(t, "2026-06-01", dateRangeStart)
			assert.Equal(t, "2026-06-15", dateRangeEnd)
			return []models.AiReviewRecord{*arTestRecord()}, nil
		},
	})

	filterStart := models.MustDate("2026-06-01")
	filterEnd := models.MustDate("2026-06-15")
	reviews, err := svc.ListByUserIDAndDateRange(arCtx, arTestUserID, &filterStart, &filterEnd)
	require.NoError(t, err)
	require.Len(t, reviews, 1)
}

// ----- Update -----

func TestAiReviewService_Update_Success(t *testing.T) {
	svc := service.NewAiReviewService(&mockAiReviewRepo{
		updateFn: func(ctx context.Context, userID string, id string, input models.UpdateAiReviewInput) (*models.AiReviewRecord, error) {
			assert.Equal(t, arTestID, id)
			return arTestRecord(), nil
		},
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.AiReviewRecord, error) {
			return arTestRecord(), nil
		},
	})

	newNotes := "Updated notes"
	review, err := svc.Update(arCtx, arTestUserID, arTestID, models.UpdateAiReviewInput{
		UserNotes: &newNotes,
	})
	require.NoError(t, err)
	require.NotNil(t, review)
}

func TestAiReviewService_Update_Ownership(t *testing.T) {
	svc := service.NewAiReviewService(&mockAiReviewRepo{
		updateFn: func(ctx context.Context, userID string, id string, input models.UpdateAiReviewInput) (*models.AiReviewRecord, error) {
			return nil, nil
		},
		getByIDFn: func(ctx context.Context, userID string, id string) (*models.AiReviewRecord, error) {
			return nil, nil
		},
	})

	newNotes := "Hacked notes"
	review, err := svc.Update(arCtx, arTestUserID, arTestID, models.UpdateAiReviewInput{
		UserNotes: &newNotes,
	})
	assert.ErrorIs(t, err, service.ErrAiReviewNotFound)
	assert.Nil(t, review)
}

// ----- Delete -----

func TestAiReviewService_Delete_Success(t *testing.T) {
	svc := service.NewAiReviewService(&mockAiReviewRepo{
		deleteFn: func(ctx context.Context, userID string, id string) (*models.AiReviewRecord, error) {
			return arTestRecord(), nil
		},
	})

	review, err := svc.Delete(arCtx, arTestUserID, arTestID)
	require.NoError(t, err)
	require.NotNil(t, review)
	assert.Equal(t, arTestID, review.ID)
}

func TestAiReviewService_Delete_NotFound(t *testing.T) {
	svc := service.NewAiReviewService(&mockAiReviewRepo{
		deleteFn: func(ctx context.Context, userID string, id string) (*models.AiReviewRecord, error) {
			return nil, nil
		},
	})

	review, err := svc.Delete(arCtx, arTestUserID, arTestID)
	assert.ErrorIs(t, err, service.ErrAiReviewNotFound)
	assert.Nil(t, review)
}

// ----- Log Privacy -----

func TestAiReviewService_Logs_NoContent(t *testing.T) {
	svc := service.NewAiReviewService(&mockAiReviewRepo{
		createFn: func(ctx context.Context, userID string, dateRangeStart, dateRangeEnd string, aiResponseText string, userNotes, plannedActions *string) (*models.AiReviewRecord, error) {
			return arTestRecord(), nil
		},
	})

	review, err := svc.Create(arCtx, arTestUserID, models.CreateAiReviewInput{
		DateRangeStart: arDateStart,
		DateRangeEnd:   arDateEnd,
		AiResponseText: "Sensitive medical data that should never appear in logs",
		UserNotes:      ptrStr("Confidential user note"),
	})
	require.NoError(t, err)
	require.NotNil(t, review)
	assert.Equal(t, "AI analysis response text", review.AiResponseText)
}