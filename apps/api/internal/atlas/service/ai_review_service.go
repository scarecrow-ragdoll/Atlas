// FILE: apps/api/internal/atlas/service/ai_review_service.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement the transport-neutral AiReviewService for WAVE-08 AI review CRUD with validation.
//   SCOPE: Create (with date range validation), GetByID, ListByUserID, ListByUserIDAndDateRange, Update, Delete, and ListAllByUserID for WAVE-09 backup consumption. Validates aiResponseText non-empty and dateRangeEnd >= dateRangeStart.
//   DEPENDS: apps/api/internal/atlas/repository/postgres.AiReviewRepository, apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-08.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AiReviewService - Interface for AI review business operations.
//   NewAiReviewService - Creates a new AiReviewService.
//   Create - Validates and creates an AI review.
//   GetByID - Gets an AI review by ID (user-scoped).
//   ListByUserID - Lists AI reviews for a user.
//   ListByUserIDAndDateRange - Lists AI reviews filtered by date range.
//   Update - Updates an AI review with partial merge.
//   Delete - Deletes an AI review.
//   ListAllByUserID - Lists all reviews for a user (for WAVE-09 backup consumption).
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added AI review service for WAVE-08.
// END_CHANGE_SUMMARY

package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasRepo "monorepo-template/apps/api/internal/atlas/repository/postgres"
)

var (
	ErrAiReviewNotFound          = errors.New("ai review not found")
	ErrAiReviewInvalidDateRange  = errors.New("date range end must be after or equal to date range start")
	ErrAiReviewEmptyText         = errors.New("ai response text must not be empty")
)

type AiReviewService interface {
	Create(ctx context.Context, userID string, input models.CreateAiReviewInput) (*models.AiReview, error)
	GetByID(ctx context.Context, userID string, id string) (*models.AiReview, error)
	ListByUserID(ctx context.Context, userID string) ([]models.AiReview, error)
	ListByUserIDAndDateRange(ctx context.Context, userID string, start, end *models.Date) ([]models.AiReview, error)
	Update(ctx context.Context, userID string, id string, input models.UpdateAiReviewInput) (*models.AiReview, error)
	Delete(ctx context.Context, userID string, id string) (*models.AiReview, error)
	ListAllByUserID(ctx context.Context, userID string) ([]models.AiReview, error)
}

type aiReviewService struct {
	repo atlasRepo.AiReviewRepository
}

func NewAiReviewService(repo atlasRepo.AiReviewRepository) AiReviewService {
	return &aiReviewService{repo: repo}
}

func (s *aiReviewService) Create(ctx context.Context, userID string, input models.CreateAiReviewInput) (*models.AiReview, error) {
	text := strings.TrimSpace(input.AiResponseText)
	if text == "" {
		return nil, ErrAiReviewEmptyText
	}

	if input.DateRangeEnd.Time().Before(input.DateRangeStart.Time()) {
		return nil, ErrAiReviewInvalidDateRange
	}

	record, err := s.repo.Create(ctx, userID, input.DateRangeStart.String(), input.DateRangeEnd.String(), text, input.UserNotes, input.PlannedActions)
	if err != nil {
		return nil, fmt.Errorf("ai_review_service.Create: %w", err)
	}

	return models.AiReviewFromRecord(record), nil
}

func (s *aiReviewService) GetByID(ctx context.Context, userID string, id string) (*models.AiReview, error) {
	record, err := s.repo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("ai_review_service.GetByID: %w", err)
	}
	if record == nil {
		return nil, ErrAiReviewNotFound
	}
	return models.AiReviewFromRecord(record), nil
}

func (s *aiReviewService) ListByUserID(ctx context.Context, userID string) ([]models.AiReview, error) {
	records, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ai_review_service.ListByUserID: %w", err)
	}

	out := make([]models.AiReview, len(records))
	for i := range records {
		out[i] = *models.AiReviewFromRecord(&records[i])
	}
	return out, nil
}

func (s *aiReviewService) ListByUserIDAndDateRange(ctx context.Context, userID string, start, end *models.Date) ([]models.AiReview, error) {
	var records []models.AiReviewRecord
	var err error

	if start != nil && end != nil {
		records, err = s.repo.ListByUserIDAndDateRange(ctx, userID, start.String(), end.String())
	} else {
		records, err = s.repo.ListByUserID(ctx, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("ai_review_service.ListByUserIDAndDateRange: %w", err)
	}

	out := make([]models.AiReview, len(records))
	for i := range records {
		out[i] = *models.AiReviewFromRecord(&records[i])
	}
	return out, nil
}

func (s *aiReviewService) Update(ctx context.Context, userID string, id string, input models.UpdateAiReviewInput) (*models.AiReview, error) {
	if input.AiResponseText != nil {
		text := strings.TrimSpace(*input.AiResponseText)
		if text == "" {
			return nil, ErrAiReviewEmptyText
		}
		input.AiResponseText = &text
	}

	if input.DateRangeStart != nil && input.DateRangeEnd != nil {
		if input.DateRangeEnd.Time().Before(input.DateRangeStart.Time()) {
			return nil, ErrAiReviewInvalidDateRange
		}
	}

	record, err := s.repo.Update(ctx, userID, id, input)
	if err != nil {
		return nil, fmt.Errorf("ai_review_service.Update: %w", err)
	}
	if record == nil {
		return nil, ErrAiReviewNotFound
	}
	return models.AiReviewFromRecord(record), nil
}

func (s *aiReviewService) Delete(ctx context.Context, userID string, id string) (*models.AiReview, error) {
	record, err := s.repo.Delete(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("ai_review_service.Delete: %w", err)
	}
	if record == nil {
		return nil, ErrAiReviewNotFound
	}
	return models.AiReviewFromRecord(record), nil
}

func (s *aiReviewService) ListAllByUserID(ctx context.Context, userID string) ([]models.AiReview, error) {
	return s.ListByUserID(ctx, userID)
}