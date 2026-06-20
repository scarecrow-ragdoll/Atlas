// FILE: apps/api/internal/atlas/service/body_weight.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement the transport-neutral BodyWeightService for WAVE-04 body weight tracking with validation.
//   SCOPE: Create, GetByID, ListByDateRange, Latest, Update, Delete. Validates weight > 0, source enum.
//   DEPENDS: apps/api/internal/atlas/repository/postgres.BodyWeightEntryRepository, apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   BodyWeightService - Interface for body weight business operations.
//   NewBodyWeightService - Creates a new BodyWeightService.
//   Create - Validates and creates a body weight entry.
//   GetByID - Gets a body weight entry by ID.
//   ListByDateRange - Lists body weight entries by date range.
//   Latest - Gets the latest body weight entry (returns nil if none).
//   Update - Validates and updates a body weight entry.
//   Delete - Deletes a body weight entry.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added body weight service for WAVE-04.
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
	ErrBodyWeightInvalid = errors.New("weight must be greater than 0")
	ErrBodyWeightInvalidSource = errors.New("invalid body weight source")
	ErrBodyWeightNotFound = errors.New("body weight entry not found")
)

type BodyWeightService interface {
	Create(ctx context.Context, userID string, input models.CreateBodyWeightInput) (*models.BodyWeightEntry, error)
	GetByID(ctx context.Context, userID string, id string) (*models.BodyWeightEntry, error)
	ListByDateRange(ctx context.Context, userID string, fromDate models.Date, toDate models.Date) ([]models.BodyWeightEntry, error)
	Latest(ctx context.Context, userID string) (*models.BodyWeightEntry, error)
	Update(ctx context.Context, userID string, id string, input models.UpdateBodyWeightInput) (*models.BodyWeightEntry, error)
	Delete(ctx context.Context, userID string, id string) (*models.BodyWeightEntry, error)
}

type bodyWeightService struct {
	repo atlasRepo.BodyWeightEntryRepository
}

func NewBodyWeightService(repo atlasRepo.BodyWeightEntryRepository) BodyWeightService {
	return &bodyWeightService{repo: repo}
}

func (s *bodyWeightService) Create(ctx context.Context, userID string, input models.CreateBodyWeightInput) (*models.BodyWeightEntry, error) {
	if input.Weight <= 0 {
		return nil, ErrBodyWeightInvalid
	}
	src := strings.TrimSpace(string(input.Source))
	if !models.IsValidBodyWeightSource(src) {
		return nil, ErrBodyWeightInvalidSource
	}

	record, err := s.repo.Create(ctx, userID, input.Date, input.Weight, src, input.Notes)
	if err != nil {
		return nil, fmt.Errorf("body_weight_service.Create: %w", err)
	}

	return models.BodyWeightEntryFromRecord(record), nil
}

func (s *bodyWeightService) GetByID(ctx context.Context, userID string, id string) (*models.BodyWeightEntry, error) {
	record, err := s.repo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("body_weight_service.GetByID: %w", err)
	}
	if record == nil {
		return nil, ErrBodyWeightNotFound
	}
	return models.BodyWeightEntryFromRecord(record), nil
}

func (s *bodyWeightService) ListByDateRange(ctx context.Context, userID string, fromDate models.Date, toDate models.Date) ([]models.BodyWeightEntry, error) {
	if fromDate.Time().After(toDate.Time()) {
		return nil, fmt.Errorf("from date must be on or before to date")
	}

	records, err := s.repo.ListByDateRange(ctx, userID, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("body_weight_service.ListByDateRange: %w", err)
	}

	out := make([]models.BodyWeightEntry, len(records))
	for i := range records {
		out[i] = *models.BodyWeightEntryFromRecord(&records[i])
	}
	return out, nil
}

func (s *bodyWeightService) Latest(ctx context.Context, userID string) (*models.BodyWeightEntry, error) {
	record, err := s.repo.Latest(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("body_weight_service.Latest: %w", err)
	}
	if record == nil {
		return nil, nil
	}
	return models.BodyWeightEntryFromRecord(record), nil
}

func (s *bodyWeightService) Update(ctx context.Context, userID string, id string, input models.UpdateBodyWeightInput) (*models.BodyWeightEntry, error) {
	existing, err := s.repo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("body_weight_service.Update: %w", err)
	}
	if existing == nil {
		return nil, ErrBodyWeightNotFound
	}

	weight := existing.Weight
	if input.Weight != nil {
		if *input.Weight <= 0 {
			return nil, ErrBodyWeightInvalid
		}
		weight = *input.Weight
	}

	source := existing.Source
	if input.Source != nil {
		src := strings.TrimSpace(string(*input.Source))
		if !models.IsValidBodyWeightSource(src) {
			return nil, ErrBodyWeightInvalidSource
		}
		source = src
	}

	notes := input.Notes
	if notes == nil {
		notes = existing.Notes
	}

	record, err := s.repo.Update(ctx, userID, id, &weight, &source, notes)
	if err != nil {
		return nil, fmt.Errorf("body_weight_service.Update: %w", err)
	}
	if record == nil {
		return nil, ErrBodyWeightNotFound
	}

	return models.BodyWeightEntryFromRecord(record), nil
}

func (s *bodyWeightService) Delete(ctx context.Context, userID string, id string) (*models.BodyWeightEntry, error) {
	record, err := s.repo.Delete(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("body_weight_service.Delete: %w", err)
	}
	if record == nil {
		return nil, ErrBodyWeightNotFound
	}
	return models.BodyWeightEntryFromRecord(record), nil
}