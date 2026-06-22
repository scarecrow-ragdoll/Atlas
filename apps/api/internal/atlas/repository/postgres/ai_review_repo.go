// FILE: apps/api/internal/atlas/repository/postgres/ai_review_repo.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement AiReviewRepository for WAVE-08 AI review CRUD using sqlc-generated queries.
//   SCOPE: CRUD operations for AI reviews, list by user ID, list by user ID and date range, all user-scoped.
//   DEPENDS: apps/api/internal/repository/postgres/generated, apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-08.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AiReviewRepository - Interface for AI review data access.
//   NewAiReviewRepository - Creates a new AiReviewRepository.
//   Create - Creates an AI review record.
//   GetByID - Gets an AI review by ID (user-scoped).
//   ListByUserID - Lists AI reviews by user ID.
//   ListByUserIDAndDateRange - Lists AI reviews by user ID and date range.
//   Update - Updates an AI review (partial merge).
//   Delete - Deletes an AI review.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added AI review repository for WAVE-08.
// END_CHANGE_SUMMARY

package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/repository/postgres/generated"
)

type AiReviewRepository interface {
	Create(ctx context.Context, userID string, dateRangeStart, dateRangeEnd string, aiResponseText string, userNotes, plannedActions *string) (*models.AiReviewRecord, error)
	GetByID(ctx context.Context, userID string, id string) (*models.AiReviewRecord, error)
	ListByUserID(ctx context.Context, userID string) ([]models.AiReviewRecord, error)
	ListByUserIDAndDateRange(ctx context.Context, userID string, dateRangeStart, dateRangeEnd string) ([]models.AiReviewRecord, error)
	Update(ctx context.Context, userID string, id string, input models.UpdateAiReviewInput) (*models.AiReviewRecord, error)
	Delete(ctx context.Context, userID string, id string) (*models.AiReviewRecord, error)
}

type aiReviewRepository struct {
	q *generated.Queries
}

func NewAiReviewRepository(pool *pgxpool.Pool) AiReviewRepository {
	return &aiReviewRepository{q: generated.New(pool)}
}

func (r *aiReviewRepository) Create(ctx context.Context, userID string, dateRangeStart, dateRangeEnd string, aiResponseText string, userNotes, plannedActions *string) (*models.AiReviewRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("ai_review_repo.Create: %w", err)
	}

	startDate, err := models.ParseDate(dateRangeStart)
	if err != nil {
		return nil, fmt.Errorf("ai_review_repo.Create: parse dateRangeStart: %w", err)
	}

	endDate, err := models.ParseDate(dateRangeEnd)
	if err != nil {
		return nil, fmt.Errorf("ai_review_repo.Create: parse dateRangeEnd: %w", err)
	}

	row, err := r.q.CreateAiReview(ctx, generated.CreateAiReviewParams{
		UserID:         uid,
		DateRangeStart: modelsToPGDate(startDate),
		DateRangeEnd:   modelsToPGDate(endDate),
		AiResponseText: aiResponseText,
		UserNotes:      nullableText(userNotes),
		PlannedActions: nullableText(plannedActions),
	})
	if err != nil {
		return nil, fmt.Errorf("ai_review_repo.Create: %w", err)
	}

	return aiReviewRecordFromRow(row), nil
}

func (r *aiReviewRepository) GetByID(ctx context.Context, userID string, id string) (*models.AiReviewRecord, error) {
	uid, rid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("ai_review_repo.GetByID: %w", err)
	}

	row, err := r.q.GetAiReviewByID(ctx, generated.GetAiReviewByIDParams{ID: rid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("ai_review_repo.GetByID: %w", err)
	}

	return aiReviewRecordFromRow(row), nil
}

func (r *aiReviewRepository) ListByUserID(ctx context.Context, userID string) ([]models.AiReviewRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("ai_review_repo.ListByUserID: %w", err)
	}

	rows, err := r.q.ListAiReviewsByUserID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("ai_review_repo.ListByUserID: %w", err)
	}

	out := make([]models.AiReviewRecord, len(rows))
	for i, row := range rows {
		out[i] = *aiReviewRecordFromRow(row)
	}
	return out, nil
}

func (r *aiReviewRepository) ListByUserIDAndDateRange(ctx context.Context, userID string, dateRangeStart, dateRangeEnd string) ([]models.AiReviewRecord, error) {
	uid, err := uuidFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("ai_review_repo.ListByUserIDAndDateRange: %w", err)
	}

	startDate, err := models.ParseDate(dateRangeStart)
	if err != nil {
		return nil, fmt.Errorf("ai_review_repo.ListByUserIDAndDateRange: parse dateRangeStart: %w", err)
	}

	endDate, err := models.ParseDate(dateRangeEnd)
	if err != nil {
		return nil, fmt.Errorf("ai_review_repo.ListByUserIDAndDateRange: parse dateRangeEnd: %w", err)
	}

	rows, err := r.q.ListAiReviewsByUserIDAndDateRange(ctx, generated.ListAiReviewsByUserIDAndDateRangeParams{
		UserID:         uid,
		DateRangeStart: modelsToPGDate(startDate),
		DateRangeEnd:   modelsToPGDate(endDate),
	})
	if err != nil {
		return nil, fmt.Errorf("ai_review_repo.ListByUserIDAndDateRange: %w", err)
	}

	out := make([]models.AiReviewRecord, len(rows))
	for i, row := range rows {
		out[i] = *aiReviewRecordFromRow(row)
	}
	return out, nil
}

func (r *aiReviewRepository) Update(ctx context.Context, userID string, id string, input models.UpdateAiReviewInput) (*models.AiReviewRecord, error) {
	existing, err := r.GetByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("ai_review_repo.Update: get existing: %w", err)
	}
	if existing == nil {
		return nil, nil
	}

	dateRangeStart := existing.DateRangeStart
	if input.DateRangeStart != nil {
		dateRangeStart = *input.DateRangeStart
	}

	dateRangeEnd := existing.DateRangeEnd
	if input.DateRangeEnd != nil {
		dateRangeEnd = *input.DateRangeEnd
	}

	aiResponseText := existing.AiResponseText
	if input.AiResponseText != nil {
		aiResponseText = *input.AiResponseText
	}

	userNotes := existing.UserNotes
	if input.UserNotes != nil {
		userNotes = input.UserNotes
	}

	plannedActions := existing.PlannedActions
	if input.PlannedActions != nil {
		plannedActions = input.PlannedActions
	}

	uid, rid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("ai_review_repo.Update: %w", err)
	}

	row, err := r.q.UpdateAiReview(ctx, generated.UpdateAiReviewParams{
		ID:             rid,
		UserID:         uid,
		DateRangeStart: modelsToPGDate(dateRangeStart),
		DateRangeEnd:   modelsToPGDate(dateRangeEnd),
		AiResponseText: aiResponseText,
		UserNotes:      nullableText(userNotes),
		PlannedActions: nullableText(plannedActions),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("ai_review_repo.Update: %w", err)
	}

	return aiReviewRecordFromRow(row), nil
}

func (r *aiReviewRepository) Delete(ctx context.Context, userID string, id string) (*models.AiReviewRecord, error) {
	uid, rid, err := parseTwoUUIDs(userID, id)
	if err != nil {
		return nil, fmt.Errorf("ai_review_repo.Delete: %w", err)
	}

	row, err := r.q.DeleteAiReview(ctx, generated.DeleteAiReviewParams{ID: rid, UserID: uid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("ai_review_repo.Delete: %w", err)
	}

	return aiReviewRecordFromRow(row), nil
}

func aiReviewRecordFromRow(row generated.AiReview) *models.AiReviewRecord {
	return &models.AiReviewRecord{
		ID:              row.ID.String(),
		UserID:          row.UserID.String(),
		DateRangeStart:  pgDateToModelsDate(row.DateRangeStart),
		DateRangeEnd:    pgDateToModelsDate(row.DateRangeEnd),
		AiResponseText:  row.AiResponseText,
		UserNotes:       textPtr(row.UserNotes),
		PlannedActions:  textPtr(row.PlannedActions),
		CreatedAt:       formatTimestamp(row.CreatedAt),
		UpdatedAt:       formatTimestamp(row.UpdatedAt),
	}
}