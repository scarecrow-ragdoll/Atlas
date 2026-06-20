// FILE: apps/api/internal/atlas/service/nutrition_override_service.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement DailyNutritionOverrideService with upsert semantics (one per user per date), override isolation, and item CRUD.
//   SCOPE: Create (upsert), GetByID (with items), GetByDate, ListByRange, Update, Delete (cascade), CreateItem, UpdateItem, DeleteItem. Validation: date required and valid format, amountGrams > 0 for items, operation must be add/subtract/replace. Log markers: [DailyNutritionOverride][create|update|delete|get|list], [DailyNutritionOverrideItem][create|update|delete].
//   DEPENDS: postgres.DailyNutritionOverrideRepository, postgres.DailyNutritionOverrideItemRepository, models.
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
	ErrOverrideDateRequired         = errors.New("date is required")
	ErrOverrideNotFound             = errors.New("daily nutrition override not found")
	ErrOverrideItemNotFound         = errors.New("daily nutrition override item not found")
	ErrOverrideItemAmountInvalid    = errors.New("amountGrams must be greater than 0")
	ErrOverrideItemOperationInvalid = errors.New("operation must be add, subtract, or replace")
)

type DailyNutritionOverrideService interface {
	Create(ctx context.Context, userID string, input models.CreateOverrideInput) (*models.DailyNutritionOverride, error)
	GetByID(ctx context.Context, userID string, id string) (*models.DailyNutritionOverride, error)
	GetByDate(ctx context.Context, userID string, date string) (*models.DailyNutritionOverride, error)
	ListByRange(ctx context.Context, userID string, startDate, endDate string) ([]models.DailyNutritionOverride, error)
	Update(ctx context.Context, userID string, id string, input models.UpdateOverrideInput) (*models.DailyNutritionOverride, error)
	Delete(ctx context.Context, userID string, id string) (*models.DailyNutritionOverride, error)
	CreateItem(ctx context.Context, userID string, input models.CreateOverrideItemInput) (*models.DailyNutritionOverrideItem, error)
	UpdateItem(ctx context.Context, userID string, itemID string, input models.UpdateOverrideItemInput) (*models.DailyNutritionOverrideItem, error)
	DeleteItem(ctx context.Context, userID string, itemID string) (*models.DailyNutritionOverrideItem, error)
}

type dailyNutritionOverrideService struct {
	repo     postgres.DailyNutritionOverrideRepository
	itemRepo postgres.DailyNutritionOverrideItemRepository
	logger   *zap.Logger
}

func NewNutritionOverrideService(repo postgres.DailyNutritionOverrideRepository, itemRepo postgres.DailyNutritionOverrideItemRepository, logger *zap.Logger) DailyNutritionOverrideService {
	return &dailyNutritionOverrideService{repo: repo, itemRepo: itemRepo, logger: logger}
}

func (s *dailyNutritionOverrideService) loadItems(ctx context.Context, overrideID string) []models.DailyNutritionOverrideItem {
	records, err := s.itemRepo.ListByOverride(ctx, overrideID)
	if err != nil {
		return []models.DailyNutritionOverrideItem{}
	}
	return models.DailyNutritionOverrideItemsFromRecords(records)
}

func (s *dailyNutritionOverrideService) Create(ctx context.Context, userID string, input models.CreateOverrideInput) (*models.DailyNutritionOverride, error) {
	s.logger.Info("[DailyNutritionOverride][create]")
	d := strings.TrimSpace(input.Date.String())
	if d == "" {
		return nil, ErrOverrideDateRequired
	}

	if _, err := time.Parse("2006-01-02", d); err != nil {
		return nil, fmt.Errorf("%w: invalid date", ErrOverrideDateRequired)
	}

	record, err := s.repo.Upsert(ctx, userID, d, input.Notes)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.Create: %w", err)
	}

	return models.DailyNutritionOverrideFromRecord(record, nil), nil
}

func (s *dailyNutritionOverrideService) GetByID(ctx context.Context, userID string, id string) (*models.DailyNutritionOverride, error) {
	s.logger.Info("[DailyNutritionOverride][get]")
	record, err := s.repo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.GetByID: %w", err)
	}
	if record == nil {
		return nil, ErrOverrideNotFound
	}

	items := s.loadItems(ctx, record.ID)
	return models.DailyNutritionOverrideFromRecord(record, items), nil
}

func (s *dailyNutritionOverrideService) GetByDate(ctx context.Context, userID string, date string) (*models.DailyNutritionOverride, error) {
	record, err := s.repo.GetByDate(ctx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.GetByDate: %w", err)
	}
	if record == nil {
		return nil, nil
	}

	items := s.loadItems(ctx, record.ID)
	return models.DailyNutritionOverrideFromRecord(record, items), nil
}

func (s *dailyNutritionOverrideService) ListByRange(ctx context.Context, userID string, startDate, endDate string) ([]models.DailyNutritionOverride, error) {
	s.logger.Info("[DailyNutritionOverride][list]")
	records, err := s.repo.ListByRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.ListByRange: %w", err)
	}

	out := make([]models.DailyNutritionOverride, len(records))
	for i := range records {
		items := s.loadItems(ctx, records[i].ID)
		out[i] = *models.DailyNutritionOverrideFromRecord(&records[i], items)
	}
	return out, nil
}

func (s *dailyNutritionOverrideService) Update(ctx context.Context, userID string, id string, input models.UpdateOverrideInput) (*models.DailyNutritionOverride, error) {
	s.logger.Info("[DailyNutritionOverride][update]")
	existing, err := s.repo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.Update: %w", err)
	}
	if existing == nil {
		return nil, ErrOverrideNotFound
	}

	notes := input.Notes
	if notes == nil {
		notes = existing.Notes
	}

	record, err := s.repo.Update(ctx, userID, id, notes)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.Update: %w", err)
	}
	if record == nil {
		return nil, ErrOverrideNotFound
	}

	items := s.loadItems(ctx, record.ID)
	return models.DailyNutritionOverrideFromRecord(record, items), nil
}

func (s *dailyNutritionOverrideService) Delete(ctx context.Context, userID string, id string) (*models.DailyNutritionOverride, error) {
	s.logger.Info("[DailyNutritionOverride][delete]")
	record, err := s.repo.Delete(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.Delete: %w", err)
	}
	if record == nil {
		return nil, ErrOverrideNotFound
	}
	return models.DailyNutritionOverrideFromRecord(record, nil), nil
}

func (s *dailyNutritionOverrideService) CreateItem(ctx context.Context, userID string, input models.CreateOverrideItemInput) (*models.DailyNutritionOverrideItem, error) {
	s.logger.Info("[DailyNutritionOverrideItem][create]")
	if input.AmountGrams <= 0 {
		return nil, ErrOverrideItemAmountInvalid
	}
	if input.Operation != models.OperationAdd && input.Operation != models.OperationSubtract && input.Operation != models.OperationReplace {
		return nil, ErrOverrideItemOperationInvalid
	}

	override, err := s.repo.GetByID(ctx, userID, input.OverrideID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.CreateItem: %w", err)
	}
	if override == nil {
		return nil, ErrOverrideNotFound
	}

	record, err := s.itemRepo.Create(ctx, input.OverrideID, input.ProductID, input.AmountGrams, string(input.Operation), input.MealLabel, input.Notes)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.CreateItem: %w", err)
	}

	return models.DailyNutritionOverrideItemFromRecord(record), nil
}

func (s *dailyNutritionOverrideService) UpdateItem(ctx context.Context, userID string, itemID string, input models.UpdateOverrideItemInput) (*models.DailyNutritionOverrideItem, error) {
	s.logger.Info("[DailyNutritionOverrideItem][update]")
	existing, err := s.itemRepo.GetByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.UpdateItem: %w", err)
	}
	if existing == nil {
		return nil, ErrOverrideItemNotFound
	}

	amount := existing.AmountGrams
	if input.AmountGrams != nil {
		if *input.AmountGrams <= 0 {
			return nil, ErrOverrideItemAmountInvalid
		}
		amount = *input.AmountGrams
	}

	op := existing.Operation
	if input.Operation != nil {
		if *input.Operation != models.OperationAdd && *input.Operation != models.OperationSubtract && *input.Operation != models.OperationReplace {
			return nil, ErrOverrideItemOperationInvalid
		}
		op = string(*input.Operation)
	}

	mealLabel := input.MealLabel
	if mealLabel == nil {
		mealLabel = existing.MealLabel
	}

	notes := input.Notes
	if notes == nil {
		notes = existing.Notes
	}

	record, err := s.itemRepo.Update(ctx, itemID, amount, op, mealLabel, notes)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.UpdateItem: %w", err)
	}
	if record == nil {
		return nil, ErrOverrideItemNotFound
	}

	return models.DailyNutritionOverrideItemFromRecord(record), nil
}

func (s *dailyNutritionOverrideService) DeleteItem(ctx context.Context, userID string, itemID string) (*models.DailyNutritionOverrideItem, error) {
	s.logger.Info("[DailyNutritionOverrideItem][delete]")
	record, err := s.itemRepo.Delete(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("nutrition_override_service.DeleteItem: %w", err)
	}
	if record == nil {
		return nil, ErrOverrideItemNotFound
	}
	return models.DailyNutritionOverrideItemFromRecord(record), nil
}
