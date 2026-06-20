// FILE: apps/api/internal/atlas/service/week_flag.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement the transport-neutral WeekFlagService for WAVE-04 week flag tracking with validation.
//   SCOPE: Create (with unique per week per type enforcement), ListByWeekStart, Delete. Validates flag type enum.
//   DEPENDS: apps/api/internal/atlas/repository/postgres.WeekFlagRepository, apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   WeekFlagService - Interface for week flag business operations.
//   NewWeekFlagService - Creates a new WeekFlagService.
//   Create - Validates and creates a week flag.
//   ListByWeekStart - Lists week flags by week start date.
//   Delete - Deletes a week flag.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added week flag service for WAVE-04.
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
	ErrWeekFlagInvalidType = errors.New("invalid week flag type")
	ErrWeekFlagNotFound    = errors.New("week flag not found")
)

type WeekFlagService interface {
	Create(ctx context.Context, userID string, input models.CreateWeekFlagInput) (*models.WeekFlag, error)
	ListByWeekStart(ctx context.Context, userID string, weekStartDate models.Date) ([]models.WeekFlag, error)
	Delete(ctx context.Context, userID string, id string) (*models.WeekFlag, error)
}

type weekFlagService struct {
	repo atlasRepo.WeekFlagRepository
}

func NewWeekFlagService(repo atlasRepo.WeekFlagRepository) WeekFlagService {
	return &weekFlagService{repo: repo}
}

func (s *weekFlagService) Create(ctx context.Context, userID string, input models.CreateWeekFlagInput) (*models.WeekFlag, error) {
	ft := strings.TrimSpace(string(input.FlagType))
	if !models.IsValidWeekFlagType(ft) {
		return nil, ErrWeekFlagInvalidType
	}

	record, err := s.repo.Create(ctx, userID, input.WeekStartDate, ft, input.Notes)
	if err != nil {
		return nil, fmt.Errorf("week_flag_service.Create: %w", err)
	}

	return models.WeekFlagFromRecord(record), nil
}

func (s *weekFlagService) ListByWeekStart(ctx context.Context, userID string, weekStartDate models.Date) ([]models.WeekFlag, error) {
	records, err := s.repo.ListByWeekStart(ctx, userID, weekStartDate)
	if err != nil {
		return nil, fmt.Errorf("week_flag_service.ListByWeekStart: %w", err)
	}

	out := make([]models.WeekFlag, len(records))
	for i := range records {
		out[i] = *models.WeekFlagFromRecord(&records[i])
	}
	return out, nil
}

func (s *weekFlagService) Delete(ctx context.Context, userID string, id string) (*models.WeekFlag, error) {
	record, err := s.repo.Delete(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("week_flag_service.Delete: %w", err)
	}
	if record == nil {
		return nil, ErrWeekFlagNotFound
	}
	return models.WeekFlagFromRecord(record), nil
}