// FILE: apps/api/internal/atlas/service/settings_service.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement the SettingsService for the Atlas fitness tracker.
//   SCOPE: Get and Update operations; Get returns public Settings (no pinHash), Update preserves pinHash.
//   DEPENDS: apps/api/internal/atlas/repository/postgres.SettingsRepository, apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Atlas settings service for WAVE-01.
// END_CHANGE_SUMMARY

package service

import (
	"context"
	"errors"
	"fmt"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasRepo "monorepo-template/apps/api/internal/atlas/repository/postgres"
)

var (
	ErrSettingsNotFound = errors.New("settings not found")
	ErrInvalidUnits     = errors.New("units must be 'metric' or 'imperial'")
)

type SettingsService interface {
	Get(ctx context.Context, userID string) (*models.Settings, error)
	Update(ctx context.Context, userID string, input models.SettingsInput) (*models.Settings, error)
}

type settingsService struct {
	repo atlasRepo.SettingsRepository
}

func NewSettingsService(repo atlasRepo.SettingsRepository) SettingsService {
	return &settingsService{repo: repo}
}

func (s *settingsService) Get(ctx context.Context, userID string) (*models.Settings, error) {
	record, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("settings_service.Get: %w", err)
	}
	return recordToPublic(record), nil
}

func (s *settingsService) Update(ctx context.Context, userID string, input models.SettingsInput) (*models.Settings, error) {
	if input.Units != nil {
		u := *input.Units
		if u != "metric" && u != "imperial" {
			return nil, ErrInvalidUnits
		}
	}
	if input.DefaultAiExportWeeks != nil {
		w := *input.DefaultAiExportWeeks
		if w < 1 || w > 52 {
			return nil, fmt.Errorf("default_ai_export_weeks must be between 1 and 52")
		}
	}

	record, err := s.repo.UpsertSettings(ctx, userID, input)
	if err != nil {
		return nil, fmt.Errorf("settings_service.Update: %w", err)
	}
	return recordToPublic(record), nil
}

func recordToPublic(record *models.SettingsRecord) *models.Settings {
	return &models.Settings{
		PinEnabled:           record.PinEnabled,
		Units:                record.Units,
		DefaultAiExportWeeks: int(record.DefaultAiExportWeeks),
	}
}