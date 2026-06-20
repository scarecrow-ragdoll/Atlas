// FILE: apps/api/internal/atlas/service/cardio.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement the transport-neutral CardioService for WAVE-04 cardio entry operations with validation and DailyLog auto-creation.
//   SCOPE: Create, GetByID, ListByDate, Update, Delete for cardio entries. Validates cardio type enum, duration > 0, pulse positive, zone valid. Auto-creates DailyLog when date has no log.
//   DEPENDS: apps/api/internal/atlas/repository/postgres.CardioEntryRepository, apps/api/internal/atlas/repository/postgres.DailyLogRepository, apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   CardioService - Interface for cardio business operations.
//   NewCardioService - Creates a new CardioService.
//   Create - Validates and creates a cardio entry, auto-creating DailyLog if needed.
//   GetByID - Gets a cardio entry by ID.
//   ListByDate - Lists cardio entries for a user+date via DailyLog.
//   Update - Validates and updates a cardio entry.
//   Delete - Deletes a cardio entry.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added cardio service for WAVE-04.
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
	ErrCardioInvalidType     = errors.New("invalid cardio type")
	ErrCardioDurationInvalid = errors.New("duration must be greater than 0")
	ErrCardioPulseInvalid    = errors.New("avg pulse must be greater than 0")
	ErrCardioZoneInvalid     = errors.New("invalid heart rate zone")
	ErrCardioNotFound        = errors.New("cardio entry not found")
)

type CardioService interface {
	Create(ctx context.Context, userID string, input models.CreateCardioInput) (*models.CardioEntry, error)
	GetByID(ctx context.Context, userID string, id string) (*models.CardioEntry, error)
	ListByDate(ctx context.Context, userID string, date models.Date) ([]models.CardioEntry, error)
	Update(ctx context.Context, userID string, id string, input models.UpdateCardioInput) (*models.CardioEntry, error)
	Delete(ctx context.Context, userID string, id string) (*models.CardioEntry, error)
}

type cardioService struct {
	cardioRepo atlasRepo.CardioEntryRepository
	dailyLogRepo atlasRepo.DailyLogRepository
}

func NewCardioService(cardioRepo atlasRepo.CardioEntryRepository, dailyLogRepo atlasRepo.DailyLogRepository) CardioService {
	return &cardioService{
		cardioRepo:   cardioRepo,
		dailyLogRepo: dailyLogRepo,
	}
}

func (s *cardioService) Create(ctx context.Context, userID string, input models.CreateCardioInput) (*models.CardioEntry, error) {
	ct := strings.TrimSpace(string(input.CardioType))
	if !models.IsValidCardioType(ct) {
		return nil, ErrCardioInvalidType
	}
	if input.DurationMinutes <= 0 {
		return nil, ErrCardioDurationInvalid
	}
	if input.AvgPulse != nil && *input.AvgPulse <= 0 {
		return nil, ErrCardioPulseInvalid
	}
	if input.HeartRateZone != nil {
		z := strings.TrimSpace(string(*input.HeartRateZone))
		if !models.IsValidHeartRateZone(z) {
			return nil, ErrCardioZoneInvalid
		}
	}

	dailyLog, err := s.dailyLogRepo.GetDailyLogByDate(ctx, userID, input.Date)
	if err != nil {
		return nil, fmt.Errorf("cardio_service.Create: %w", err)
	}
	if dailyLog == nil {
		dailyLog, err = s.dailyLogRepo.GetOrCreateDailyLogByDate(ctx, userID, input.Date)
		if err != nil {
			return nil, fmt.Errorf("cardio_service.Create: %w", err)
		}
	}

	zone := (*string)(nil)
	if input.HeartRateZone != nil {
		z := string(*input.HeartRateZone)
		zone = &z
	}

	record, err := s.cardioRepo.Create(ctx, userID, dailyLog.ID, ct, input.DurationMinutes, input.AvgPulse, zone, input.Notes)
	if err != nil {
		return nil, fmt.Errorf("cardio_service.Create: %w", err)
	}

	return models.CardioEntryFromRecord(record), nil
}

func (s *cardioService) GetByID(ctx context.Context, userID string, id string) (*models.CardioEntry, error) {
	record, err := s.cardioRepo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("cardio_service.GetByID: %w", err)
	}
	if record == nil {
		return nil, ErrCardioNotFound
	}
	return models.CardioEntryFromRecord(record), nil
}

func (s *cardioService) ListByDate(ctx context.Context, userID string, date models.Date) ([]models.CardioEntry, error) {
	dailyLog, err := s.dailyLogRepo.GetDailyLogByDate(ctx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("cardio_service.ListByDate: %w", err)
	}
	if dailyLog == nil {
		return []models.CardioEntry{}, nil
	}

	records, err := s.cardioRepo.ListByDailyLog(ctx, userID, dailyLog.ID)
	if err != nil {
		return nil, fmt.Errorf("cardio_service.ListByDate: %w", err)
	}

	out := make([]models.CardioEntry, len(records))
	for i := range records {
		out[i] = *models.CardioEntryFromRecord(&records[i])
	}
	return out, nil
}

func (s *cardioService) Update(ctx context.Context, userID string, id string, input models.UpdateCardioInput) (*models.CardioEntry, error) {
	existing, err := s.cardioRepo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("cardio_service.Update: %w", err)
	}
	if existing == nil {
		return nil, ErrCardioNotFound
	}

	ct := existing.CardioType
	if input.CardioType != nil {
		ct = strings.TrimSpace(string(*input.CardioType))
		if !models.IsValidCardioType(ct) {
			return nil, ErrCardioInvalidType
		}
	}

	duration := existing.DurationMinutes
	if input.DurationMinutes != nil {
		if *input.DurationMinutes <= 0 {
			return nil, ErrCardioDurationInvalid
		}
		duration = *input.DurationMinutes
	}

	avgPulse := existing.AvgPulse
	if input.AvgPulse != nil {
		if *input.AvgPulse <= 0 {
			return nil, ErrCardioPulseInvalid
		}
		avgPulse = input.AvgPulse
	}

	zone := existing.HeartRateZone
	if input.HeartRateZone != nil {
		z := string(*input.HeartRateZone)
		if !models.IsValidHeartRateZone(z) {
			return nil, ErrCardioZoneInvalid
		}
		zone = &z
	}

	notes := input.Notes
	if notes == nil {
		notes = existing.Notes
	}

	var zoneStr *string
	if zone != nil {
		z := string(*zone)
		zoneStr = &z
	}

	record, err := s.cardioRepo.Update(ctx, userID, id, ct, duration, avgPulse, zoneStr, notes)
	if err != nil {
		return nil, fmt.Errorf("cardio_service.Update: %w", err)
	}
	if record == nil {
		return nil, ErrCardioNotFound
	}

	return models.CardioEntryFromRecord(record), nil
}

func (s *cardioService) Delete(ctx context.Context, userID string, id string) (*models.CardioEntry, error) {
	record, err := s.cardioRepo.Delete(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("cardio_service.Delete: %w", err)
	}
	if record == nil {
		return nil, ErrCardioNotFound
	}
	return models.CardioEntryFromRecord(record), nil
}