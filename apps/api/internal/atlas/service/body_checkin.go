// FILE: apps/api/internal/atlas/service/body_checkin.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement the transport-neutral BodyCheckInService and BodyMeasurementService for WAVE-04 body check-in and measurement operations with validation.
//   SCOPE: Check-in CRUD, measurement CRUD, validation (weight > 0, body fat % 0-100, measurement type enum, value > 0, side rules for paired types). Includes measurement side validation per DDEC-W04-003.
//   DEPENDS: apps/api/internal/atlas/repository/postgres.BodyCheckInRepository, BodyMeasurementRepository, ProgressPhotoRepository; apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-04.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   BodyCheckInService - Interface for body check-in business operations.
//   NewBodyCheckInService - Creates a new BodyCheckInService.
//   Create - Validates and creates a body check-in.
//   GetByID - Gets a body check-in by ID with nested measurements and photos.
//   ListByDateRange - Lists body check-ins by date range.
//   Update - Validates and updates a body check-in.
//   Delete - Deletes a check-in (cascade handled by DB FK).
//   BodyMeasurementService - Interface for body measurement business operations.
//   NewBodyMeasurementService - Creates a new BodyMeasurementService.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added body check-in and measurement service for WAVE-04.
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
	ErrCheckInWeightInvalid      = errors.New("check-in weight must be greater than 0")
	ErrCheckInBodyFatInvalid     = errors.New("body fat percentage must be between 0 and 100")
	ErrCheckInNotFound           = errors.New("body check-in not found")
	ErrMeasurementTypeInvalid    = errors.New("invalid measurement type")
	ErrMeasurementValueInvalid   = errors.New("measurement value must be greater than 0")
	ErrMeasurementSideInvalid    = errors.New("side is only allowed for paired measurement types (forearm, biceps, thigh, calf)")
	ErrMeasurementNotFound       = errors.New("body measurement not found")
)

type BodyCheckInService interface {
	Create(ctx context.Context, userID string, input models.CreateCheckInInput) (*models.BodyCheckIn, error)
	GetByID(ctx context.Context, userID string, id string) (*models.BodyCheckIn, error)
	ListByDateRange(ctx context.Context, userID string, fromDate models.Date, toDate models.Date) ([]models.BodyCheckIn, error)
	Update(ctx context.Context, userID string, id string, input models.UpdateCheckInInput) (*models.BodyCheckIn, error)
	Delete(ctx context.Context, userID string, id string) (*models.BodyCheckIn, error)
}

type bodyCheckInService struct {
	checkInRepo    atlasRepo.BodyCheckInRepository
	measurementRepo atlasRepo.BodyMeasurementRepository
	photoRepo      atlasRepo.ProgressPhotoRepository
}

func NewBodyCheckInService(checkInRepo atlasRepo.BodyCheckInRepository, measurementRepo atlasRepo.BodyMeasurementRepository, photoRepo atlasRepo.ProgressPhotoRepository) BodyCheckInService {
	return &bodyCheckInService{
		checkInRepo:     checkInRepo,
		measurementRepo: measurementRepo,
		photoRepo:       photoRepo,
	}
}

func (s *bodyCheckInService) Create(ctx context.Context, userID string, input models.CreateCheckInInput) (*models.BodyCheckIn, error) {
	if input.Weight != nil && *input.Weight <= 0 {
		return nil, ErrCheckInWeightInvalid
	}
	if input.BodyFatPercentage != nil {
		if *input.BodyFatPercentage <= 0 || *input.BodyFatPercentage > 100 {
			return nil, ErrCheckInBodyFatInvalid
		}
	}

	record, err := s.checkInRepo.Create(ctx, userID, input.Date, input.Weight, input.BodyFatPercentage, input.Notes)
	if err != nil {
		return nil, fmt.Errorf("body_checkin_service.Create: %w", err)
	}

	return models.BodyCheckInFromRecord(record), nil
}

func (s *bodyCheckInService) GetByID(ctx context.Context, userID string, id string) (*models.BodyCheckIn, error) {
	record, err := s.checkInRepo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("body_checkin_service.GetByID: %w", err)
	}
	if record == nil {
		return nil, ErrCheckInNotFound
	}

	checkIn := models.BodyCheckInFromRecord(record)

	measurements, err := s.measurementRepo.ListByCheckIn(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("body_checkin_service.GetByID: %w", err)
	}
	checkIn.Measurements = make([]models.BodyMeasurement, len(measurements))
	for i := range measurements {
		checkIn.Measurements[i] = *models.BodyMeasurementFromRecord(&measurements[i])
	}

	photos, err := s.photoRepo.ListByCheckIn(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("body_checkin_service.GetByID: %w", err)
	}
	checkIn.ProgressPhotos = make([]models.ProgressPhoto, len(photos))
	for i := range photos {
		checkIn.ProgressPhotos[i] = *models.ProgressPhotoFromRecord(&photos[i])
	}

	return checkIn, nil
}

func (s *bodyCheckInService) ListByDateRange(ctx context.Context, userID string, fromDate models.Date, toDate models.Date) ([]models.BodyCheckIn, error) {
	if fromDate.Time().After(toDate.Time()) {
		return nil, fmt.Errorf("from date must be on or before to date")
	}

	records, err := s.checkInRepo.ListByDateRange(ctx, userID, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("body_checkin_service.ListByDateRange: %w", err)
	}

	out := make([]models.BodyCheckIn, len(records))
	for i := range records {
		ci := models.BodyCheckInFromRecord(&records[i])
		if ci != nil {
			out[i] = *ci
		}
	}
	return out, nil
}

func (s *bodyCheckInService) Update(ctx context.Context, userID string, id string, input models.UpdateCheckInInput) (*models.BodyCheckIn, error) {
	existing, err := s.checkInRepo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("body_checkin_service.Update: %w", err)
	}
	if existing == nil {
		return nil, ErrCheckInNotFound
	}

	weight := existing.Weight
	if input.Weight != nil {
		if *input.Weight <= 0 {
			return nil, ErrCheckInWeightInvalid
		}
		weight = input.Weight
	}

	bfp := existing.BodyFatPercentage
	if input.BodyFatPercentage != nil {
		if *input.BodyFatPercentage <= 0 || *input.BodyFatPercentage > 100 {
			return nil, ErrCheckInBodyFatInvalid
		}
		bfp = input.BodyFatPercentage
	}

	notes := input.Notes
	if notes == nil {
		notes = existing.Notes
	}

	record, err := s.checkInRepo.Update(ctx, userID, id, weight, bfp, notes)
	if err != nil {
		return nil, fmt.Errorf("body_checkin_service.Update: %w", err)
	}
	if record == nil {
		return nil, ErrCheckInNotFound
	}

	return models.BodyCheckInFromRecord(record), nil
}

func (s *bodyCheckInService) Delete(ctx context.Context, userID string, id string) (*models.BodyCheckIn, error) {
	record, err := s.checkInRepo.Delete(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("body_checkin_service.Delete: %w", err)
	}
	if record == nil {
		return nil, ErrCheckInNotFound
	}
	return models.BodyCheckInFromRecord(record), nil
}

type BodyMeasurementService interface {
	Create(ctx context.Context, userID string, input models.CreateMeasurementInput) (*models.BodyMeasurement, error)
	GetByID(ctx context.Context, userID string, id string) (*models.BodyMeasurement, error)
	ListByCheckIn(ctx context.Context, userID string, checkInID string) ([]models.BodyMeasurement, error)
	Update(ctx context.Context, userID string, id string, input models.UpdateMeasurementInput) (*models.BodyMeasurement, error)
	Delete(ctx context.Context, userID string, id string) (*models.BodyMeasurement, error)
}

type bodyMeasurementService struct {
	measurementRepo atlasRepo.BodyMeasurementRepository
	checkInRepo     atlasRepo.BodyCheckInRepository
}

func NewBodyMeasurementService(measurementRepo atlasRepo.BodyMeasurementRepository, checkInRepo atlasRepo.BodyCheckInRepository) BodyMeasurementService {
	return &bodyMeasurementService{
		measurementRepo: measurementRepo,
		checkInRepo:     checkInRepo,
	}
}

func (s *bodyMeasurementService) Create(ctx context.Context, userID string, input models.CreateMeasurementInput) (*models.BodyMeasurement, error) {
	mt := strings.TrimSpace(string(input.MeasurementType))
	if !models.IsValidMeasurementType(mt) {
		return nil, ErrMeasurementTypeInvalid
	}
	if input.Value <= 0 {
		return nil, ErrMeasurementValueInvalid
	}

	if input.Side != nil && *input.Side != models.MeasurementSideNone {
		sideStr := strings.TrimSpace(string(*input.Side))
		if !models.IsValidMeasurementSide(sideStr) {
			return nil, ErrMeasurementSideInvalid
		}
		if !models.IsPairedMeasurementType(models.MeasurementType(mt)) {
			return nil, ErrMeasurementSideInvalid
		}
	}

	record, err := s.measurementRepo.Create(ctx, "", mt, sidePtr(input.Side), input.Value)
	if err != nil {
		return nil, fmt.Errorf("body_measurement_service.Create: %w", err)
	}

	return models.BodyMeasurementFromRecord(record), nil
}

func (s *bodyMeasurementService) GetByID(ctx context.Context, userID string, id string) (*models.BodyMeasurement, error) {
	record, err := s.measurementRepo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("body_measurement_service.GetByID: %w", err)
	}
	if record == nil {
		return nil, ErrMeasurementNotFound
	}
	return models.BodyMeasurementFromRecord(record), nil
}

func (s *bodyMeasurementService) ListByCheckIn(ctx context.Context, userID string, checkInID string) ([]models.BodyMeasurement, error) {
	records, err := s.measurementRepo.ListByCheckIn(ctx, userID, checkInID)
	if err != nil {
		return nil, fmt.Errorf("body_measurement_service.ListByCheckIn: %w", err)
	}
	out := make([]models.BodyMeasurement, len(records))
	for i := range records {
		out[i] = *models.BodyMeasurementFromRecord(&records[i])
	}
	return out, nil
}

func (s *bodyMeasurementService) Update(ctx context.Context, userID string, id string, input models.UpdateMeasurementInput) (*models.BodyMeasurement, error) {
	existing, err := s.measurementRepo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("body_measurement_service.Update: %w", err)
	}
	if existing == nil {
		return nil, ErrMeasurementNotFound
	}

	mt := existing.MeasurementType
	if input.MeasurementType != nil {
		mt = strings.TrimSpace(string(*input.MeasurementType))
		if !models.IsValidMeasurementType(mt) {
			return nil, ErrMeasurementTypeInvalid
		}
	}

	value := existing.Value
	if input.Value != nil {
		if *input.Value <= 0 {
			return nil, ErrMeasurementValueInvalid
		}
		value = *input.Value
	}

	var side *string
	if input.Side != nil {
		if *input.Side != models.MeasurementSideNone {
			sideStr := strings.TrimSpace(string(*input.Side))
			if !models.IsValidMeasurementSide(sideStr) {
				return nil, ErrMeasurementSideInvalid
			}
			if !models.IsPairedMeasurementType(models.MeasurementType(mt)) {
				return nil, ErrMeasurementSideInvalid
			}
			side = &sideStr
		}
	} else {
		side = existing.Side
	}

	record, err := s.measurementRepo.Update(ctx, userID, id, mt, side, value)
	if err != nil {
		return nil, fmt.Errorf("body_measurement_service.Update: %w", err)
	}
	if record == nil {
		return nil, ErrMeasurementNotFound
	}

	return models.BodyMeasurementFromRecord(record), nil
}

func (s *bodyMeasurementService) Delete(ctx context.Context, userID string, id string) (*models.BodyMeasurement, error) {
	record, err := s.measurementRepo.Delete(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("body_measurement_service.Delete: %w", err)
	}
	if record == nil {
		return nil, ErrMeasurementNotFound
	}
	return models.BodyMeasurementFromRecord(record), nil
}

func sidePtr(side *models.MeasurementSide) *string {
	if side == nil {
		return nil
	}
	s := string(*side)
	return &s
}