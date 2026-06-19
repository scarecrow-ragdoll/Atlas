// FILE: apps/api/internal/atlas/service/exercise.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement the transport-neutral ExerciseService for WAVE-02 Exercise Library with validation.
//   SCOPE: Create, Get, List, ListAll, Update, Archive, Restore for exercises; media management. Validates name (required, trimmed, non-empty) and working weight (> 0 when provided). Does not log personalNotes or media content.
//   DEPENDS: apps/api/internal/atlas/repository/postgres.ExerciseRepository, apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-02.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   ExerciseService - Interface for exercise business operations.
//   NewExerciseService - Creates a new ExerciseService.
//   Create - Validates and creates a new exercise.
//   GetByID - Gets an exercise by ID, includes media.
//   List - Paginated exercise list with cursor support.
//   ListAll - Unpaginated exercise list for WAVE-03.
//   Update - Validates and updates exercise fields.
//   Archive - Soft-deletes an exercise (sets is_active=false).
//   Restore - Restores an archived exercise (sets is_active=true).
//   CreateMedia - Validates and creates exercise media.
//   GetMediaByID - Gets exercise media by ID.
//   DeleteMedia - Deletes exercise media and returns metadata for file cleanup.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added exercise service for WAVE-02.
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
	ErrExerciseNotFound    = errors.New("exercise not found")
	ErrExerciseNameEmpty   = errors.New("exercise name is required")
	ErrWeightInvalid       = errors.New("working weight must be greater than 0")
)

type ExerciseService interface {
	Create(ctx context.Context, userID string, input models.CreateExerciseInput) (*models.Exercise, error)
	GetByID(ctx context.Context, userID string, id string) (*models.Exercise, error)
	List(ctx context.Context, userID string, first int32, after *string, includeInactive bool) (*models.ExerciseConnection, error)
	ListAll(ctx context.Context, userID string, includeInactive bool) ([]models.Exercise, error)
	Update(ctx context.Context, userID string, id string, input models.UpdateExerciseInput) (*models.Exercise, error)
	Archive(ctx context.Context, userID string, id string) (*models.Exercise, error)
	Restore(ctx context.Context, userID string, id string) (*models.Exercise, error)
	CreateMedia(ctx context.Context, userID string, exerciseID string, fileName string, filePath string, mimeType string, fileSize int64) (*models.ExerciseMedia, error)
	GetMediaByID(ctx context.Context, userID string, id string) (*models.ExerciseMedia, error)
	GetMediaRecordByID(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error)
	DeleteMedia(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error)
}

type exerciseService struct {
	repo atlasRepo.ExerciseRepository
}

func NewExerciseService(repo atlasRepo.ExerciseRepository) ExerciseService {
	return &exerciseService{repo: repo}
}

func (s *exerciseService) Create(ctx context.Context, userID string, input models.CreateExerciseInput) (*models.Exercise, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrExerciseNameEmpty
	}
	input.Name = name

	if input.WorkingWeight != nil && *input.WorkingWeight <= 0 {
		return nil, ErrWeightInvalid
	}

	record, err := s.repo.Create(ctx, userID, input)
	if err != nil {
		return nil, fmt.Errorf("exercise_service.Create: %w", err)
	}

	return exerciseFromRecord(record), nil
}

func (s *exerciseService) GetByID(ctx context.Context, userID string, id string) (*models.Exercise, error) {
	record, err := s.repo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("exercise_service.GetByID: %w", err)
	}
	if record == nil {
		return nil, ErrExerciseNotFound
	}

	ex := exerciseFromRecord(record)
	media, err := s.repo.ListMediaByExercise(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("exercise_service.GetByID: %w", err)
	}
	ex.Media = media
	return ex, nil
}

func (s *exerciseService) List(ctx context.Context, userID string, first int32, after *string, includeInactive bool) (*models.ExerciseConnection, error) {
	isActive := !includeInactive

	var records []models.ExerciseRecord
	var err error

	if first <= 0 {
		first = 20
	}

	if after != nil && *after != "" {
		records, err = s.repo.ListCursor(ctx, userID, isActive, *after, first+1)
	} else {
		records, err = s.repo.List(ctx, userID, isActive, first+1)
	}
	if err != nil {
		return nil, fmt.Errorf("exercise_service.List: %w", err)
	}

	totalCount, err := s.repo.Count(ctx, userID, isActive)
	if err != nil {
		return nil, fmt.Errorf("exercise_service.List: %w", err)
	}

	hasNextPage := len(records) > int(first)
	if hasNextPage {
		records = records[:first]
	}

	items := make([]models.Exercise, len(records))
	for i, rec := range records {
		items[i] = *exerciseFromRecord(&rec)
	}

	var endCursor *string
	if len(items) > 0 {
		cursor := items[len(items)-1].Name
		endCursor = &cursor
	}

	return &models.ExerciseConnection{
		Items:      items,
		TotalCount: totalCount,
		PageInfo: models.PageInfo{
			HasNextPage: hasNextPage,
			EndCursor:   endCursor,
		},
	}, nil
}

func (s *exerciseService) ListAll(ctx context.Context, userID string, includeInactive bool) ([]models.Exercise, error) {
	records, err := s.repo.ListAll(ctx, userID, includeInactive)
	if err != nil {
		return nil, fmt.Errorf("exercise_service.ListAll: %w", err)
	}

	items := make([]models.Exercise, len(records))
	for i, rec := range records {
		items[i] = *exerciseFromRecord(&rec)
	}
	return items, nil
}

func (s *exerciseService) Update(ctx context.Context, userID string, id string, input models.UpdateExerciseInput) (*models.Exercise, error) {
	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, ErrExerciseNameEmpty
		}
		input.Name = &name
	}

	if input.WorkingWeight != nil && *input.WorkingWeight <= 0 {
		return nil, ErrWeightInvalid
	}

	record, err := s.repo.Update(ctx, userID, id, input)
	if err != nil {
		return nil, fmt.Errorf("exercise_service.Update: %w", err)
	}
	if record == nil {
		return nil, ErrExerciseNotFound
	}

	return exerciseFromRecord(record), nil
}

func (s *exerciseService) Archive(ctx context.Context, userID string, id string) (*models.Exercise, error) {
	record, err := s.repo.Archive(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("exercise_service.Archive: %w", err)
	}
	if record == nil {
		return nil, ErrExerciseNotFound
	}

	return exerciseFromRecord(record), nil
}

func (s *exerciseService) Restore(ctx context.Context, userID string, id string) (*models.Exercise, error) {
	record, err := s.repo.Restore(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("exercise_service.Restore: %w", err)
	}
	if record == nil {
		return nil, ErrExerciseNotFound
	}

	return exerciseFromRecord(record), nil
}

func (s *exerciseService) CreateMedia(ctx context.Context, userID, exerciseID, fileName, filePath, mimeType string, fileSize int64) (*models.ExerciseMedia, error) {
	media, err := s.repo.CreateMedia(ctx, userID, exerciseID, fileName, filePath, mimeType, fileSize)
	if err != nil {
		return nil, fmt.Errorf("exercise_service.CreateMedia: %w", err)
	}
	return media, nil
}

func (s *exerciseService) GetMediaByID(ctx context.Context, userID string, id string) (*models.ExerciseMedia, error) {
	media, err := s.repo.GetMediaByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("exercise_service.GetMediaByID: %w", err)
	}
	if media == nil {
		return nil, ErrExerciseNotFound
	}
	return media, nil
}

func (s *exerciseService) GetMediaRecordByID(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error) {
	record, err := s.repo.GetMediaRecordByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("exercise_service.GetMediaRecordByID: %w", err)
	}
	if record == nil {
		return nil, ErrExerciseNotFound
	}
	return record, nil
}

// DeleteMedia deletes exercise media by ID and returns the record including file path for cleanup.
func (s *exerciseService) DeleteMedia(ctx context.Context, userID string, id string) (*models.ExerciseMediaRecord, error) {
	record, err := s.repo.DeleteMedia(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("exercise_service.DeleteMedia: %w", err)
	}
	if record == nil {
		return nil, ErrExerciseNotFound
	}
	return record, nil
}

func exerciseFromRecord(record *models.ExerciseRecord) *models.Exercise {
	return &models.Exercise{
		ID:            record.ID,
		UserID:        record.UserID,
		Name:          record.Name,
		MuscleGroups:  record.MuscleGroups,
		Description:   record.Description,
		PersonalNotes: record.PersonalNotes,
		WorkingWeight: record.WorkingWeight,
		IsActive:      record.IsActive,
		Media:         []models.ExerciseMedia{},
		CreatedAt:     record.CreatedAt,
		UpdatedAt:     record.UpdatedAt,
	}
}
