// FILE: apps/api/internal/atlas/service/nutrition_template_service.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement NutritionTemplateService with upsert semantics (one per user per week), cascade item loading, and log markers.
//   SCOPE: Create (upsert), GetByID (with items), GetCurrent (by week), ListByRange, Update, Delete (cascade). Validation: weekStartDate required and valid date format. Log markers: [NutritionTemplate][create|get|list|update|delete|current].
//   DEPENDS: postgres.NutritionTemplateRepository, postgres.NutritionTemplateItemRepository, models.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT

package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/atlas/repository/postgres"
)

var (
	ErrTemplateWeekRequired = errors.New("weekStartDate is required")
	ErrTemplateNotFound     = errors.New("nutrition template not found")
	ErrTemplateItemNotFound = errors.New("nutrition template item not found")
)

type NutritionTemplateService interface {
	Create(ctx context.Context, userID string, input models.CreateTemplateInput) (*models.NutritionTemplate, error)
	GetByID(ctx context.Context, userID string, id string) (*models.NutritionTemplate, error)
	GetCurrent(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplate, error)
	ListByRange(ctx context.Context, userID string, startDate, endDate string) ([]models.NutritionTemplate, error)
	Update(ctx context.Context, userID string, id string, input models.UpdateTemplateInput) (*models.NutritionTemplate, error)
	Delete(ctx context.Context, userID string, id string) (*models.NutritionTemplate, error)
}

type nutritionTemplateService struct {
	repo     postgres.NutritionTemplateRepository
	itemRepo postgres.NutritionTemplateItemRepository
	logger   *zap.Logger
}

func NewNutritionTemplateService(repo postgres.NutritionTemplateRepository, itemRepo postgres.NutritionTemplateItemRepository, logger *zap.Logger) NutritionTemplateService {
	return &nutritionTemplateService{repo: repo, itemRepo: itemRepo, logger: logger}
}

func (s *nutritionTemplateService) loadItems(ctx context.Context, templateID string) []models.NutritionTemplateItem {
	records, err := s.itemRepo.ListByTemplate(ctx, templateID)
	if err != nil {
		return []models.NutritionTemplateItem{}
	}
	return models.NutritionTemplateItemsFromRecords(records)
}

func (s *nutritionTemplateService) Create(ctx context.Context, userID string, input models.CreateTemplateInput) (*models.NutritionTemplate, error) {
	s.logger.Info("[NutritionTemplate][create]")
	wd := strings.TrimSpace(input.WeekStartDate.String())
	if wd == "" {
		return nil, ErrTemplateWeekRequired
	}

	if _, err := time.Parse("2006-01-02", wd); err != nil {
		return nil, fmt.Errorf("%w: invalid date format", ErrTemplateWeekRequired)
	}

	record, err := s.repo.Upsert(ctx, userID, wd, input.Title, input.Notes)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_service.Create: %w", err)
	}

	return models.NutritionTemplateFromRecord(record, nil), nil
}

func (s *nutritionTemplateService) GetByID(ctx context.Context, userID string, id string) (*models.NutritionTemplate, error) {
	s.logger.Info("[NutritionTemplate][get]")
	record, err := s.repo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_service.GetByID: %w", err)
	}
	if record == nil {
		return nil, ErrTemplateNotFound
	}

	items := s.loadItems(ctx, record.ID)
	return models.NutritionTemplateFromRecord(record, items), nil
}

func (s *nutritionTemplateService) GetCurrent(ctx context.Context, userID string, weekStartDate string) (*models.NutritionTemplate, error) {
	s.logger.Info("[NutritionTemplate][current]")
	record, err := s.repo.GetByWeek(ctx, userID, weekStartDate)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_service.GetCurrent: %w", err)
	}
	if record == nil {
		return nil, nil
	}

	items := s.loadItems(ctx, record.ID)
	return models.NutritionTemplateFromRecord(record, items), nil
}

func (s *nutritionTemplateService) ListByRange(ctx context.Context, userID string, startDate, endDate string) ([]models.NutritionTemplate, error) {
	s.logger.Info("[NutritionTemplate][list]")
	records, err := s.repo.ListByRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_service.ListByRange: %w", err)
	}

	out := make([]models.NutritionTemplate, len(records))
	for i := range records {
		items := s.loadItems(ctx, records[i].ID)
		out[i] = *models.NutritionTemplateFromRecord(&records[i], items)
	}
	return out, nil
}

func (s *nutritionTemplateService) Update(ctx context.Context, userID string, id string, input models.UpdateTemplateInput) (*models.NutritionTemplate, error) {
	s.logger.Info("[NutritionTemplate][update]")
	existing, err := s.repo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_service.Update: %w", err)
	}
	if existing == nil {
		return nil, ErrTemplateNotFound
	}

	title := input.Title
	if title == nil {
		title = existing.Title
	}
	notes := input.Notes
	if notes == nil {
		notes = existing.Notes
	}

	record, err := s.repo.Update(ctx, userID, id, title, notes)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_service.Update: %w", err)
	}
	if record == nil {
		return nil, ErrTemplateNotFound
	}

	items := s.loadItems(ctx, record.ID)
	return models.NutritionTemplateFromRecord(record, items), nil
}

func (s *nutritionTemplateService) Delete(ctx context.Context, userID string, id string) (*models.NutritionTemplate, error) {
	s.logger.Info("[NutritionTemplate][delete]")
	record, err := s.repo.Delete(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_template_service.Delete: %w", err)
	}
	if record == nil {
		return nil, ErrTemplateNotFound
	}
	return models.NutritionTemplateFromRecord(record, nil), nil
}
