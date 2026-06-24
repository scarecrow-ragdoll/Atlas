// FILE: apps/api/internal/atlas/service/daily_nutrition_log_service.go
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Implement factual DailyNutritionLogService with product snapshot entry CRUD and snapshot-based aggregate totals.
//   SCOPE: GetByDate get-or-create, range listing, notes update, entry add/update/delete validation, product existence/active checks for new entries, and user-scoped aggregate reloads.
//   DEPENDS: postgres.DailyNutritionLogRepository, NutritionProductService, apps/api/internal/atlas/models, zap.
//   LINKS: M-API-NUTRITION / V-M-API-NUTRITION.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   DailyNutritionLogService - Interface for factual daily nutrition log operations.
//   NewDailyNutritionLogService - Creates a service with repository, product service, and logger dependencies.
//   GetByDate/ListByRange/UpdateNotes - Daily log aggregate operations.
//   AddEntry/UpdateEntry/DeleteEntry - Entry mutations returning refreshed snapshot totals where possible.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Reload parent daily log metadata for entry update/delete aggregate responses.
// END_CHANGE_SUMMARY

package service

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/atlas/repository/postgres"
)

var (
	ErrDailyNutritionAmountInvalid   = errors.New("amountGrams must be greater than 0")
	ErrDailyNutritionProductNotFound = errors.New("nutrition product not found for daily entry")
	ErrDailyNutritionProductInactive = errors.New("nutrition product is inactive")
	ErrDailyNutritionLogNotFound     = errors.New("daily nutrition log not found")
	ErrDailyNutritionEntryNotFound   = errors.New("daily nutrition entry not found")
)

type DailyNutritionLogService interface {
	GetByDate(ctx context.Context, userID string, date models.Date) (*models.DailyNutritionLog, error)
	ListByRange(ctx context.Context, userID string, from, to models.Date) ([]models.DailyNutritionLog, error)
	UpdateNotes(ctx context.Context, userID string, id string, input models.UpdateDailyNutritionLogNotesInput) (*models.DailyNutritionLog, error)
	AddEntry(ctx context.Context, userID string, input models.AddDailyNutritionEntryInput) (*models.DailyNutritionLog, error)
	UpdateEntry(ctx context.Context, userID string, id string, input models.UpdateDailyNutritionEntryInput) (*models.DailyNutritionLog, error)
	DeleteEntry(ctx context.Context, userID string, id string) (*models.DailyNutritionLog, error)
}

type dailyNutritionLogService struct {
	repo           postgres.DailyNutritionLogRepository
	productService NutritionProductService
	logger         *zap.Logger
}

func NewDailyNutritionLogService(repo postgres.DailyNutritionLogRepository, productService NutritionProductService, logger *zap.Logger) DailyNutritionLogService {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &dailyNutritionLogService{repo: repo, productService: productService, logger: logger}
}

func (s *dailyNutritionLogService) GetByDate(ctx context.Context, userID string, date models.Date) (*models.DailyNutritionLog, error) {
	s.logger.Info("[DailyNutritionLog][get]")
	record, err := s.repo.GetOrCreate(ctx, userID, date, nil)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_service.GetByDate: %w", err)
	}
	return s.loadAggregate(ctx, userID, record)
}

func (s *dailyNutritionLogService) ListByRange(ctx context.Context, userID string, from, to models.Date) ([]models.DailyNutritionLog, error) {
	s.logger.Info("[DailyNutritionLog][list]")
	records, err := s.repo.ListByRange(ctx, userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_service.ListByRange: %w", err)
	}

	out := make([]models.DailyNutritionLog, len(records))
	for i := range records {
		log, err := s.loadAggregate(ctx, userID, &records[i])
		if err != nil {
			return nil, err
		}
		out[i] = *log
	}
	return out, nil
}

func (s *dailyNutritionLogService) UpdateNotes(ctx context.Context, userID string, id string, input models.UpdateDailyNutritionLogNotesInput) (*models.DailyNutritionLog, error) {
	s.logger.Info("[DailyNutritionLog][update]")
	record, err := s.repo.UpdateNotes(ctx, userID, id, input.Notes)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_service.UpdateNotes: %w", err)
	}
	if record == nil {
		return nil, ErrDailyNutritionLogNotFound
	}
	return s.loadAggregate(ctx, userID, record)
}

func (s *dailyNutritionLogService) AddEntry(ctx context.Context, userID string, input models.AddDailyNutritionEntryInput) (*models.DailyNutritionLog, error) {
	s.logger.Info("[DailyNutritionEntry][create]")
	if input.AmountGrams <= 0 {
		return nil, ErrDailyNutritionAmountInvalid
	}

	product, err := s.productService.GetByID(ctx, userID, input.ProductID)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			return nil, ErrDailyNutritionProductNotFound
		}
		return nil, fmt.Errorf("daily_nutrition_log_service.AddEntry: %w", err)
	}
	if product == nil {
		return nil, ErrDailyNutritionProductNotFound
	}
	if !product.IsActive {
		return nil, ErrDailyNutritionProductInactive
	}

	logRecord, err := s.repo.GetOrCreate(ctx, userID, input.Date, nil)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_service.AddEntry: %w", err)
	}
	if logRecord == nil {
		return nil, ErrDailyNutritionLogNotFound
	}

	created, err := s.repo.AddEntry(ctx, userID, models.CreateDailyNutritionEntryRecordInput{
		DailyLogID:  logRecord.ID,
		ProductID:   input.ProductID,
		AmountGrams: input.AmountGrams,
		MealLabel:   input.MealLabel,
		Notes:       input.Notes,
		Position:    input.Position,
	})
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_service.AddEntry: %w", err)
	}
	if created == nil {
		return nil, ErrDailyNutritionEntryNotFound
	}

	return s.loadAggregate(ctx, userID, logRecord)
}

func (s *dailyNutritionLogService) UpdateEntry(ctx context.Context, userID string, id string, input models.UpdateDailyNutritionEntryInput) (*models.DailyNutritionLog, error) {
	s.logger.Info("[DailyNutritionEntry][update]")
	if input.AmountGrams == nil || *input.AmountGrams <= 0 {
		return nil, ErrDailyNutritionAmountInvalid
	}

	record, err := s.repo.UpdateEntry(ctx, userID, id, input)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_service.UpdateEntry: %w", err)
	}
	if record == nil {
		return nil, ErrDailyNutritionEntryNotFound
	}

	return s.loadEntryAggregate(ctx, userID, record)
}

func (s *dailyNutritionLogService) DeleteEntry(ctx context.Context, userID string, id string) (*models.DailyNutritionLog, error) {
	s.logger.Info("[DailyNutritionEntry][delete]")
	record, err := s.repo.DeleteEntry(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_service.DeleteEntry: %w", err)
	}
	if record == nil {
		return nil, ErrDailyNutritionEntryNotFound
	}

	return s.loadEntryAggregate(ctx, userID, record)
}

func (s *dailyNutritionLogService) loadAggregate(ctx context.Context, userID string, record *models.DailyNutritionLogRecord) (*models.DailyNutritionLog, error) {
	if record == nil {
		return nil, ErrDailyNutritionLogNotFound
	}
	entries, err := s.repo.ListEntries(ctx, userID, record.ID)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_service.loadAggregate: %w", err)
	}
	return models.DailyNutritionLogFromRecord(record, models.DailyNutritionEntriesFromRecords(entries)), nil
}

func (s *dailyNutritionLogService) loadEntryAggregate(ctx context.Context, userID string, entry *models.DailyNutritionEntryRecord) (*models.DailyNutritionLog, error) {
	parent, err := s.repo.GetByID(ctx, userID, entry.DailyLogID)
	if err != nil {
		return nil, fmt.Errorf("daily_nutrition_log_service.loadEntryAggregate: %w", err)
	}
	if parent == nil {
		return nil, ErrDailyNutritionLogNotFound
	}
	return s.loadAggregate(ctx, userID, parent)
}
