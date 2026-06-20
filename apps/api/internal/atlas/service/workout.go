// FILE: apps/api/internal/atlas/service/workout.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement the transport-neutral WorkoutService for WAVE-03 DailyLog strength diary operations.
//   SCOPE: DailyLog reads, summaries, notes, workout exercise mutations, workout set mutations, ordering validation, snapshots, and optimistic version checks.
//   DEPENDS: apps/api/internal/atlas/models, apps/api/internal/atlas/repository/postgres.
//   LINKS: M-API / V-M-API / WAVE-03.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   WorkoutService - Interface for DailyLog aggregate operations.
//   NewWorkoutService - Creates a WorkoutService backed by workout and exercise repositories.
//   GetDailyLog - Reads a DailyLog aggregate without creating it.
//   ListDailyLogSummaries - Reads DailyLog summaries for a valid date span.
//   UpdateDailyLogNotes - Creates or updates notes after version validation.
//   AddWorkoutExercise - Adds an exercise instance with a working weight snapshot.
//   UpdateWorkoutExercise - Updates exercise notes or position after version validation.
//   RemoveWorkoutExercise - Removes an exercise while retaining the DailyLog row.
//   ReorderWorkoutExercises - Reorders exercises only when IDs exactly match the aggregate.
//   AddWorkoutSet - Adds a validated set to a workout exercise.
//   UpdateWorkoutSet - Updates set values or position after version validation.
//   RemoveWorkoutSet - Removes a set and keeps sibling numbering contiguous.
//   ReorderWorkoutSets - Reorders sets only when IDs exactly match the exercise.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.3 - Bounded append-position int32 conversion for lint-clean WAVE-03 service gates.
// END_CHANGE_SUMMARY

package service

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5"

	"monorepo-template/apps/api/internal/atlas/models"
	atlasRepo "monorepo-template/apps/api/internal/atlas/repository/postgres"
)

type WorkoutService interface {
	GetDailyLog(ctx context.Context, userID string, date models.Date) (*models.DailyLog, error)
	ListDailyLogSummaries(ctx context.Context, userID string, from models.Date, to models.Date) ([]models.DailyLogSummary, error)
	UpdateDailyLogNotes(ctx context.Context, userID string, date models.Date, expectedVersion int32, notes *string) (*models.DailyLog, error)
	AddWorkoutExercise(ctx context.Context, userID string, date models.Date, expectedVersion int32, input models.AddWorkoutExerciseInput) (*models.DailyLog, error)
	UpdateWorkoutExercise(ctx context.Context, userID string, id string, expectedVersion int32, input models.UpdateWorkoutExerciseInput) (*models.DailyLog, error)
	RemoveWorkoutExercise(ctx context.Context, userID string, id string, expectedVersion int32) (*models.DailyLog, error)
	ReorderWorkoutExercises(ctx context.Context, userID string, date models.Date, expectedVersion int32, orderedIDs []string) (*models.DailyLog, error)
	AddWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, expectedVersion int32, input models.AddWorkoutSetInput) (*models.DailyLog, error)
	UpdateWorkoutSet(ctx context.Context, userID string, id string, expectedVersion int32, input models.UpdateWorkoutSetInput) (*models.DailyLog, error)
	RemoveWorkoutSet(ctx context.Context, userID string, id string, expectedVersion int32) (*models.DailyLog, error)
	ReorderWorkoutSets(ctx context.Context, userID string, workoutExerciseID string, expectedVersion int32, orderedIDs []string) (*models.DailyLog, error)
}

type workoutService struct {
	workoutRepo  atlasRepo.WorkoutRepository
	exerciseRepo atlasRepo.ExerciseRepository
}

func NewWorkoutService(workoutRepo atlasRepo.WorkoutRepository, exerciseRepo atlasRepo.ExerciseRepository) WorkoutService {
	return &workoutService{
		workoutRepo:  workoutRepo,
		exerciseRepo: exerciseRepo,
	}
}

func (s *workoutService) GetDailyLog(ctx context.Context, userID string, date models.Date) (*models.DailyLog, error) {
	record, err := s.workoutRepo.GetDailyLogByDate(ctx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("workout_service.GetDailyLog: %w", err)
	}
	if record == nil {
		return nil, nil
	}

	aggregate, err := s.workoutRepo.GetDailyLogAggregate(ctx, userID, record.ID)
	if err != nil {
		return nil, fmt.Errorf("workout_service.GetDailyLog: %w", err)
	}
	return s.dailyLogFromAggregate(ctx, aggregate)
}

func (s *workoutService) ListDailyLogSummaries(ctx context.Context, userID string, from models.Date, to models.Date) ([]models.DailyLogSummary, error) {
	if from.Time().After(to.Time()) {
		return nil, validationError("from date must be on or before to date")
	}

	records, err := s.workoutRepo.ListDailyLogSummaries(ctx, userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("workout_service.ListDailyLogSummaries: %w", err)
	}
	out := make([]models.DailyLogSummary, len(records))
	for i := range records {
		out[i] = dailyLogSummaryFromRepoRecord(records[i])
	}
	return out, nil
}

func (s *workoutService) UpdateDailyLogNotes(ctx context.Context, userID string, date models.Date, expectedVersion int32, notes *string) (*models.DailyLog, error) {
	if err := validateExpectedVersion(expectedVersion); err != nil {
		return nil, err
	}
	if err := s.ensureDailyLogForDateMutation(ctx, userID, date, expectedVersion, "workout_service.UpdateDailyLogNotes", nil); err != nil {
		return nil, err
	}

	var out *models.DailyLog
	err := s.workoutRepo.WithLockedDailyLogByDate(ctx, userID, date, func(ctx context.Context, tx atlasRepo.WorkoutTx, locked *atlasRepo.DailyLogRecord) error {
		if _, err := s.requireVersion(ctx, tx, userID, locked, expectedVersion); err != nil {
			return err
		}
		if _, err := tx.UpdateDailyLogNotes(ctx, userID, locked.ID, notes); err != nil {
			return fmt.Errorf("workout_service.UpdateDailyLogNotes: %w", err)
		}
		var err error
		out, err = s.incrementAndLoad(ctx, tx, userID, locked.ID)
		return err
	})
	if err != nil {
		return nil, mapLockError(err, "daily log not found")
	}
	return out, nil
}

func (s *workoutService) AddWorkoutExercise(ctx context.Context, userID string, date models.Date, expectedVersion int32, input models.AddWorkoutExerciseInput) (*models.DailyLog, error) {
	if err := validateExpectedVersion(expectedVersion); err != nil {
		return nil, err
	}
	if input.ExerciseID == "" {
		return nil, validationError("exercise id is required")
	}
	if input.Position != nil && *input.Position <= 0 {
		return nil, validationError("position must be greater than 0")
	}

	exercise, err := s.getExerciseRecord(ctx, userID, input.ExerciseID)
	if err != nil {
		return nil, err
	}
	if exercise == nil {
		return nil, notFoundError("exercise not found")
	}

	if err := s.ensureDailyLogForDateMutation(ctx, userID, date, expectedVersion, "workout_service.AddWorkoutExercise", func() error {
		if input.Position != nil && *input.Position > 1 {
			return validationError("position cannot exceed append position")
		}
		return nil
	}); err != nil {
		return nil, err
	}

	var out *models.DailyLog
	err = s.workoutRepo.WithLockedDailyLogByDate(ctx, userID, date, func(ctx context.Context, tx atlasRepo.WorkoutTx, locked *atlasRepo.DailyLogRecord) error {
		aggregate, err := s.requireVersion(ctx, tx, userID, locked, expectedVersion)
		if err != nil {
			return err
		}
		appendPosition := len(aggregate.WorkoutExercises) + 1
		position, err := boundedInt32(appendPosition, "position cannot exceed append position")
		if err != nil {
			return err
		}
		if input.Position != nil {
			position = *input.Position
		}
		if int64(position) > int64(appendPosition) {
			return validationError("position cannot exceed append position")
		}
		if _, err := tx.AddWorkoutExercise(ctx, userID, locked.ID, atlasRepo.AddWorkoutExerciseInput{
			ExerciseID:            input.ExerciseID,
			Position:              position,
			WorkingWeightSnapshot: exercise.WorkingWeight,
			Notes:                 input.Notes,
		}); err != nil {
			return fmt.Errorf("workout_service.AddWorkoutExercise: %w", err)
		}
		out, err = s.incrementAndLoad(ctx, tx, userID, locked.ID)
		return err
	})
	if err != nil {
		return nil, mapLockError(err, "daily log not found")
	}
	return out, nil
}

func (s *workoutService) UpdateWorkoutExercise(ctx context.Context, userID string, id string, expectedVersion int32, input models.UpdateWorkoutExerciseInput) (*models.DailyLog, error) {
	if err := validateExpectedVersion(expectedVersion); err != nil {
		return nil, err
	}
	if err := validateUpdateWorkoutExerciseInput(input); err != nil {
		return nil, err
	}

	var out *models.DailyLog
	err := s.workoutRepo.WithLockedDailyLogByWorkoutExerciseID(ctx, userID, id, func(ctx context.Context, tx atlasRepo.WorkoutTx, locked *atlasRepo.DailyLogRecord) error {
		aggregate, err := s.requireVersion(ctx, tx, userID, locked, expectedVersion)
		if err != nil {
			return err
		}
		current := findWorkoutExerciseRecord(aggregate, id)
		if current == nil {
			return notFoundError("workout exercise not found")
		}
		if input.Position != nil && int64(*input.Position) > int64(len(aggregate.WorkoutExercises)) {
			return validationError("position must refer to an existing workout exercise slot")
		}
		if isWorkoutExercisePositionOnlyUpdate(input) && current.Position == *input.Position {
			return validationError("update workout exercise requires at least one meaningful change")
		}
		updated, err := tx.UpdateWorkoutExercise(ctx, userID, id, atlasRepo.UpdateWorkoutExerciseInput{
			Position: input.Position,
			SetNotes: input.SetNotes || input.Notes != nil,
			Notes:    input.Notes,
		})
		if err != nil {
			return fmt.Errorf("workout_service.UpdateWorkoutExercise: %w", err)
		}
		if updated == nil {
			return notFoundError("workout exercise not found")
		}
		out, err = s.incrementAndLoad(ctx, tx, userID, locked.ID)
		return err
	})
	if err != nil {
		return nil, mapLockError(err, "workout exercise not found")
	}
	return out, nil
}

func (s *workoutService) RemoveWorkoutExercise(ctx context.Context, userID string, id string, expectedVersion int32) (*models.DailyLog, error) {
	if err := validateExpectedVersion(expectedVersion); err != nil {
		return nil, err
	}

	var out *models.DailyLog
	err := s.workoutRepo.WithLockedDailyLogByWorkoutExerciseID(ctx, userID, id, func(ctx context.Context, tx atlasRepo.WorkoutTx, locked *atlasRepo.DailyLogRecord) error {
		aggregate, err := s.requireVersion(ctx, tx, userID, locked, expectedVersion)
		if err != nil {
			return err
		}
		remainingIDs := workoutExerciseIDsExcept(aggregate.WorkoutExercises, id)
		deleted, err := tx.DeleteWorkoutExercise(ctx, userID, id)
		if err != nil {
			return fmt.Errorf("workout_service.RemoveWorkoutExercise: %w", err)
		}
		if deleted == nil {
			return notFoundError("workout exercise not found")
		}
		if err := tx.ReorderWorkoutExercises(ctx, userID, locked.ID, remainingIDs); err != nil {
			return fmt.Errorf("workout_service.RemoveWorkoutExercise: %w", err)
		}
		out, err = s.incrementAndLoad(ctx, tx, userID, locked.ID)
		return err
	})
	if err != nil {
		return nil, mapLockError(err, "workout exercise not found")
	}
	return out, nil
}

func (s *workoutService) ReorderWorkoutExercises(ctx context.Context, userID string, date models.Date, expectedVersion int32, orderedIDs []string) (*models.DailyLog, error) {
	if err := validateExpectedVersion(expectedVersion); err != nil {
		return nil, err
	}

	var out *models.DailyLog
	err := s.workoutRepo.WithLockedDailyLogByDate(ctx, userID, date, func(ctx context.Context, tx atlasRepo.WorkoutTx, locked *atlasRepo.DailyLogRecord) error {
		aggregate, err := s.requireVersion(ctx, tx, userID, locked, expectedVersion)
		if err != nil {
			return err
		}
		if err := validateExactIDs(workoutExerciseIDs(aggregate.WorkoutExercises), orderedIDs, "workout exercise ids"); err != nil {
			return err
		}
		if err := tx.ReorderWorkoutExercises(ctx, userID, locked.ID, orderedIDs); err != nil {
			return fmt.Errorf("workout_service.ReorderWorkoutExercises: %w", err)
		}
		out, err = s.incrementAndLoad(ctx, tx, userID, locked.ID)
		return err
	})
	if err != nil {
		return nil, mapLockError(err, "daily log not found")
	}
	return out, nil
}

func (s *workoutService) AddWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, expectedVersion int32, input models.AddWorkoutSetInput) (*models.DailyLog, error) {
	if err := validateExpectedVersion(expectedVersion); err != nil {
		return nil, err
	}
	if err := validateAddWorkoutSetInput(input); err != nil {
		return nil, err
	}

	var out *models.DailyLog
	err := s.workoutRepo.WithLockedDailyLogByWorkoutExerciseID(ctx, userID, workoutExerciseID, func(ctx context.Context, tx atlasRepo.WorkoutTx, locked *atlasRepo.DailyLogRecord) error {
		aggregate, err := s.requireVersion(ctx, tx, userID, locked, expectedVersion)
		if err != nil {
			return err
		}
		exercise := findWorkoutExerciseRecord(aggregate, workoutExerciseID)
		if exercise == nil {
			return notFoundError("workout exercise not found")
		}
		appendSetNumber := len(exercise.Sets) + 1
		setNumber, err := boundedInt32(appendSetNumber, "set number cannot exceed append position")
		if err != nil {
			return err
		}
		if input.SetNumber != nil {
			setNumber = *input.SetNumber
		}
		if int64(setNumber) > int64(appendSetNumber) {
			return validationError("set number cannot exceed append position")
		}
		if _, err := tx.AddWorkoutSet(ctx, userID, workoutExerciseID, atlasRepo.AddWorkoutSetInput{
			SetNumber: setNumber,
			Weight:    input.Weight,
			Reps:      input.Reps,
			RPE:       input.RPE,
			RIR:       input.RIR,
			Notes:     input.Notes,
		}); err != nil {
			return fmt.Errorf("workout_service.AddWorkoutSet: %w", err)
		}
		out, err = s.incrementAndLoad(ctx, tx, userID, locked.ID)
		return err
	})
	if err != nil {
		return nil, mapLockError(err, "workout exercise not found")
	}
	return out, nil
}

func (s *workoutService) UpdateWorkoutSet(ctx context.Context, userID string, id string, expectedVersion int32, input models.UpdateWorkoutSetInput) (*models.DailyLog, error) {
	if err := validateExpectedVersion(expectedVersion); err != nil {
		return nil, err
	}
	if err := validateUpdateWorkoutSetInput(input); err != nil {
		return nil, err
	}

	var out *models.DailyLog
	err := s.workoutRepo.WithLockedDailyLogByWorkoutSetID(ctx, userID, id, func(ctx context.Context, tx atlasRepo.WorkoutTx, locked *atlasRepo.DailyLogRecord) error {
		aggregate, err := s.requireVersion(ctx, tx, userID, locked, expectedVersion)
		if err != nil {
			return err
		}
		exercise := findWorkoutExerciseRecordBySetID(aggregate, id)
		if exercise == nil {
			return notFoundError("workout set not found")
		}
		currentSet := findWorkoutSetRecord(exercise.Sets, id)
		if currentSet == nil {
			return notFoundError("workout set not found")
		}
		if input.SetNumber != nil && int64(*input.SetNumber) > int64(len(exercise.Sets)) {
			return validationError("set number must refer to an existing set slot")
		}
		if isWorkoutSetNumberOnlyUpdate(input) && currentSet.SetNumber == *input.SetNumber {
			return validationError("update workout set requires at least one meaningful change")
		}
		updated, err := tx.UpdateWorkoutSet(ctx, userID, exercise.ID, id, atlasRepo.UpdateWorkoutSetInput{
			SetNumber: input.SetNumber,
			Weight:    input.Weight,
			Reps:      input.Reps,
			SetRPE:    input.SetRPE || input.RPE != nil,
			RPE:       input.RPE,
			SetRIR:    input.SetRIR || input.RIR != nil,
			RIR:       input.RIR,
			SetNotes:  input.SetNotes || input.Notes != nil,
			Notes:     input.Notes,
		})
		if err != nil {
			return fmt.Errorf("workout_service.UpdateWorkoutSet: %w", err)
		}
		if updated == nil {
			return notFoundError("workout set not found")
		}
		out, err = s.incrementAndLoad(ctx, tx, userID, locked.ID)
		return err
	})
	if err != nil {
		return nil, mapLockError(err, "workout set not found")
	}
	return out, nil
}

func (s *workoutService) RemoveWorkoutSet(ctx context.Context, userID string, id string, expectedVersion int32) (*models.DailyLog, error) {
	if err := validateExpectedVersion(expectedVersion); err != nil {
		return nil, err
	}

	var out *models.DailyLog
	err := s.workoutRepo.WithLockedDailyLogByWorkoutSetID(ctx, userID, id, func(ctx context.Context, tx atlasRepo.WorkoutTx, locked *atlasRepo.DailyLogRecord) error {
		aggregate, err := s.requireVersion(ctx, tx, userID, locked, expectedVersion)
		if err != nil {
			return err
		}
		exercise := findWorkoutExerciseRecordBySetID(aggregate, id)
		if exercise == nil {
			return notFoundError("workout set not found")
		}
		remainingIDs := workoutSetIDsExcept(exercise.Sets, id)
		deleted, err := tx.DeleteWorkoutSet(ctx, userID, exercise.ID, id)
		if err != nil {
			return fmt.Errorf("workout_service.RemoveWorkoutSet: %w", err)
		}
		if deleted == nil {
			return notFoundError("workout set not found")
		}
		if err := tx.ReorderWorkoutSets(ctx, userID, exercise.ID, remainingIDs); err != nil {
			return fmt.Errorf("workout_service.RemoveWorkoutSet: %w", err)
		}
		out, err = s.incrementAndLoad(ctx, tx, userID, locked.ID)
		return err
	})
	if err != nil {
		return nil, mapLockError(err, "workout set not found")
	}
	return out, nil
}

func (s *workoutService) ReorderWorkoutSets(ctx context.Context, userID string, workoutExerciseID string, expectedVersion int32, orderedIDs []string) (*models.DailyLog, error) {
	if err := validateExpectedVersion(expectedVersion); err != nil {
		return nil, err
	}

	var out *models.DailyLog
	err := s.workoutRepo.WithLockedDailyLogByWorkoutExerciseID(ctx, userID, workoutExerciseID, func(ctx context.Context, tx atlasRepo.WorkoutTx, locked *atlasRepo.DailyLogRecord) error {
		aggregate, err := s.requireVersion(ctx, tx, userID, locked, expectedVersion)
		if err != nil {
			return err
		}
		exercise := findWorkoutExerciseRecord(aggregate, workoutExerciseID)
		if exercise == nil {
			return notFoundError("workout exercise not found")
		}
		if err := validateExactIDs(workoutSetIDs(exercise.Sets), orderedIDs, "workout set ids"); err != nil {
			return err
		}
		if err := tx.ReorderWorkoutSets(ctx, userID, workoutExerciseID, orderedIDs); err != nil {
			return fmt.Errorf("workout_service.ReorderWorkoutSets: %w", err)
		}
		out, err = s.incrementAndLoad(ctx, tx, userID, locked.ID)
		return err
	})
	if err != nil {
		return nil, mapLockError(err, "workout exercise not found")
	}
	return out, nil
}

/*
START_CONTRACT: requireVersion
  PURPOSE: Load the locked aggregate and reject stale optimistic versions before mutation.
  INPUTS: { tx: atlasRepo.WorkoutTx - active locked transaction, locked: DailyLogRecord - locked root, expectedVersion: int32 - client version }
  OUTPUTS: { DailyLogAggregate - current aggregate when the expected version matches }
  SIDE_EFFECTS: none.
  LINKS: V-M-API / WAVE-03.
END_CONTRACT: requireVersion
*/

func (s *workoutService) requireVersion(ctx context.Context, tx atlasRepo.WorkoutTx, userID string, locked *atlasRepo.DailyLogRecord, expectedVersion int32) (*atlasRepo.DailyLogAggregate, error) {
	if locked == nil {
		return nil, notFoundError("daily log not found")
	}
	aggregate, err := tx.GetDailyLogAggregate(ctx, userID, locked.ID)
	if err != nil {
		return nil, fmt.Errorf("workout_service.requireVersion: %w", err)
	}
	if locked.Version != expectedVersion {
		current, err := s.dailyLogFromAggregate(ctx, aggregate)
		if err != nil {
			return nil, err
		}
		return nil, conflictError(expectedVersion, locked.Version, current)
	}
	return aggregate, nil
}

func (s *workoutService) incrementAndLoad(ctx context.Context, tx atlasRepo.WorkoutTx, userID string, dailyLogID string) (*models.DailyLog, error) {
	if _, err := tx.IncrementDailyLogVersion(ctx, userID, dailyLogID); err != nil {
		return nil, fmt.Errorf("workout_service.incrementAndLoad: %w", err)
	}
	aggregate, err := tx.GetDailyLogAggregate(ctx, userID, dailyLogID)
	if err != nil {
		return nil, fmt.Errorf("workout_service.incrementAndLoad: %w", err)
	}
	return s.dailyLogFromAggregate(ctx, aggregate)
}

func (s *workoutService) ensureDailyLogForDateMutation(ctx context.Context, userID string, date models.Date, expectedVersion int32, operation string, validateAbsent func() error) error {
	record, err := s.workoutRepo.GetDailyLogByDate(ctx, userID, date)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	if record != nil {
		return nil
	}
	if expectedVersion != 0 {
		return conflictError(expectedVersion, 0, nil)
	}
	if validateAbsent != nil {
		if err := validateAbsent(); err != nil {
			return err
		}
	}
	if _, err := s.workoutRepo.GetOrCreateDailyLogByDate(ctx, userID, date); err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	return nil
}

func (s *workoutService) dailyLogFromAggregate(ctx context.Context, aggregate *atlasRepo.DailyLogAggregate) (*models.DailyLog, error) {
	if aggregate == nil {
		return nil, nil
	}
	out := &models.DailyLog{
		ID:               aggregate.DailyLog.ID,
		UserID:           aggregate.DailyLog.UserID,
		Date:             aggregate.DailyLog.Date,
		Notes:            aggregate.DailyLog.Notes,
		Version:          aggregate.DailyLog.Version,
		WorkoutExercises: make([]models.WorkoutExercise, len(aggregate.WorkoutExercises)),
		CreatedAt:        aggregate.DailyLog.CreatedAt,
		UpdatedAt:        aggregate.DailyLog.UpdatedAt,
	}
	for i := range aggregate.WorkoutExercises {
		exercise, err := s.workoutExerciseFromRepoRecord(ctx, aggregate.DailyLog.UserID, aggregate.WorkoutExercises[i])
		if err != nil {
			return nil, err
		}
		out.WorkoutExercises[i] = exercise
	}
	return out, nil
}

func (s *workoutService) workoutExerciseFromRepoRecord(ctx context.Context, userID string, record atlasRepo.WorkoutExerciseRecord) (models.WorkoutExercise, error) {
	out := models.WorkoutExercise{
		ID:                    record.ID,
		UserID:                record.UserID,
		DailyLogID:            record.DailyLogID,
		ExerciseID:            record.ExerciseID,
		Position:              record.Position,
		WorkingWeightSnapshot: record.WorkingWeightSnapshot,
		Notes:                 record.Notes,
		Sets:                  make([]models.WorkoutSet, len(record.Sets)),
		CreatedAt:             record.CreatedAt,
		UpdatedAt:             record.UpdatedAt,
	}
	exercise, err := s.getExerciseRecord(ctx, userID, record.ExerciseID)
	if err != nil {
		return models.WorkoutExercise{}, err
	}
	if exercise != nil {
		out.Exercise = exerciseFromRecord(exercise)
	}
	for i := range record.Sets {
		out.Sets[i] = workoutSetFromRepoRecord(record.Sets[i])
	}
	return out, nil
}

func (s *workoutService) getExerciseRecord(ctx context.Context, userID string, exerciseID string) (*models.ExerciseRecord, error) {
	if s.exerciseRepo == nil {
		return nil, nil
	}
	record, err := s.exerciseRepo.GetByID(ctx, userID, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("workout_service.getExerciseRecord: %w", err)
	}
	return record, nil
}

func dailyLogSummaryFromRepoRecord(record atlasRepo.DailyLogSummaryRecord) models.DailyLogSummary {
	return models.DailyLogSummary{
		ID:                   record.ID,
		Date:                 record.Date,
		Version:              record.Version,
		WorkoutExerciseCount: record.WorkoutExerciseCount,
		WorkoutSetCount:      record.WorkoutSetCount,
		TotalVolume:          record.TotalVolume,
		UpdatedAt:            record.UpdatedAt,
	}
}

func workoutSetFromRepoRecord(record atlasRepo.WorkoutSetRecord) models.WorkoutSet {
	return models.WorkoutSet{
		ID:                record.ID,
		WorkoutExerciseID: record.WorkoutExerciseID,
		SetNumber:         record.SetNumber,
		Weight:            record.Weight,
		Reps:              record.Reps,
		RPE:               record.RPE,
		RIR:               record.RIR,
		Notes:             record.Notes,
		CreatedAt:         record.CreatedAt,
		UpdatedAt:         record.UpdatedAt,
	}
}

func validateExpectedVersion(expectedVersion int32) error {
	if expectedVersion < 0 {
		return validationError("expected version must be greater than or equal to 0")
	}
	return nil
}

func boundedInt32(value int, overflowMessage string) (int32, error) {
	if value > math.MaxInt32 {
		return 0, validationError("%s", overflowMessage)
	}
	return int32(value), nil //nolint:gosec // bounded above by math.MaxInt32 before conversion.
}

func validateAddWorkoutSetInput(input models.AddWorkoutSetInput) error {
	if input.SetNumber != nil && *input.SetNumber <= 0 {
		return validationError("set number must be greater than 0")
	}
	if input.Weight <= 0 {
		return validationError("weight must be greater than 0")
	}
	if input.Reps <= 0 {
		return validationError("reps must be greater than 0")
	}
	return validateOptionalSetFields(input.RPE, input.RIR)
}

func validateUpdateWorkoutExerciseInput(input models.UpdateWorkoutExerciseInput) error {
	if input.Position != nil && *input.Position <= 0 {
		return validationError("position must be greater than 0")
	}
	if !hasWorkoutExerciseUpdate(input) {
		return validationError("update workout exercise requires at least one meaningful change")
	}
	return nil
}

func validateUpdateWorkoutSetInput(input models.UpdateWorkoutSetInput) error {
	if input.SetNumber != nil && *input.SetNumber <= 0 {
		return validationError("set number must be greater than 0")
	}
	if input.Weight != nil && *input.Weight <= 0 {
		return validationError("weight must be greater than 0")
	}
	if input.Reps != nil && *input.Reps <= 0 {
		return validationError("reps must be greater than 0")
	}
	if !hasWorkoutSetUpdate(input) {
		return validationError("update workout set requires at least one meaningful change")
	}
	return validateOptionalSetFields(input.RPE, input.RIR)
}

func validateOptionalSetFields(rpe *float64, rir *int32) error {
	if rpe != nil && (*rpe < 1 || *rpe > 10) {
		return validationError("rpe must be between 1 and 10")
	}
	if rir != nil && (*rir < 0 || *rir > 10) {
		return validationError("rir must be between 0 and 10")
	}
	return nil
}

func validateExactIDs(current []string, ordered []string, label string) error {
	if len(current) != len(ordered) {
		return validationError("%s must exactly match current ids", label)
	}
	expected := make(map[string]struct{}, len(current))
	for _, id := range current {
		expected[id] = struct{}{}
	}
	seen := make(map[string]struct{}, len(ordered))
	for _, id := range ordered {
		if _, duplicate := seen[id]; duplicate {
			return validationError("%s contain duplicate id %q", label, id)
		}
		seen[id] = struct{}{}
		if _, ok := expected[id]; !ok {
			return validationError("%s contain foreign id %q", label, id)
		}
	}
	for _, id := range current {
		if _, ok := seen[id]; !ok {
			return validationError("%s are missing id %q", label, id)
		}
	}
	return nil
}

func findWorkoutExerciseRecord(aggregate *atlasRepo.DailyLogAggregate, id string) *atlasRepo.WorkoutExerciseRecord {
	if aggregate == nil {
		return nil
	}
	for i := range aggregate.WorkoutExercises {
		if aggregate.WorkoutExercises[i].ID == id {
			return &aggregate.WorkoutExercises[i]
		}
	}
	return nil
}

func findWorkoutExerciseRecordBySetID(aggregate *atlasRepo.DailyLogAggregate, setID string) *atlasRepo.WorkoutExerciseRecord {
	if aggregate == nil {
		return nil
	}
	for i := range aggregate.WorkoutExercises {
		for j := range aggregate.WorkoutExercises[i].Sets {
			if aggregate.WorkoutExercises[i].Sets[j].ID == setID {
				return &aggregate.WorkoutExercises[i]
			}
		}
	}
	return nil
}

func findWorkoutSetRecord(records []atlasRepo.WorkoutSetRecord, id string) *atlasRepo.WorkoutSetRecord {
	for i := range records {
		if records[i].ID == id {
			return &records[i]
		}
	}
	return nil
}

func hasWorkoutExerciseUpdate(input models.UpdateWorkoutExerciseInput) bool {
	return input.Position != nil || input.SetNotes || input.Notes != nil
}

func isWorkoutExercisePositionOnlyUpdate(input models.UpdateWorkoutExerciseInput) bool {
	return input.Position != nil && !input.SetNotes && input.Notes == nil
}

func hasWorkoutSetUpdate(input models.UpdateWorkoutSetInput) bool {
	return input.SetNumber != nil || input.Weight != nil || input.Reps != nil || input.SetRPE || input.RPE != nil || input.SetRIR || input.RIR != nil || input.SetNotes || input.Notes != nil
}

func isWorkoutSetNumberOnlyUpdate(input models.UpdateWorkoutSetInput) bool {
	return input.SetNumber != nil && input.Weight == nil && input.Reps == nil && !input.SetRPE && input.RPE == nil && !input.SetRIR && input.RIR == nil && !input.SetNotes && input.Notes == nil
}

func workoutExerciseIDs(records []atlasRepo.WorkoutExerciseRecord) []string {
	out := make([]string, len(records))
	for i := range records {
		out[i] = records[i].ID
	}
	return out
}

func workoutExerciseIDsExcept(records []atlasRepo.WorkoutExerciseRecord, excludedID string) []string {
	out := make([]string, 0, len(records))
	for i := range records {
		if records[i].ID != excludedID {
			out = append(out, records[i].ID)
		}
	}
	return out
}

func workoutSetIDs(records []atlasRepo.WorkoutSetRecord) []string {
	out := make([]string, len(records))
	for i := range records {
		out[i] = records[i].ID
	}
	return out
}

func workoutSetIDsExcept(records []atlasRepo.WorkoutSetRecord, excludedID string) []string {
	out := make([]string, 0, len(records))
	for i := range records {
		if records[i].ID != excludedID {
			out = append(out, records[i].ID)
		}
	}
	return out
}

func validationError(format string, args ...any) *models.DailyLogValidationErr {
	return &models.DailyLogValidationErr{
		Message: fmt.Sprintf(format, args...),
		Code:    models.DailyLogErrorValidation,
	}
}

func notFoundError(message string) *models.DailyLogNotFoundErr {
	return &models.DailyLogNotFoundErr{
		Message: message,
		Code:    models.DailyLogErrorNotFound,
	}
}

func conflictError(expectedVersion int32, currentVersion int32, current *models.DailyLog) *models.DailyLogConflictErr {
	return &models.DailyLogConflictErr{
		Message:         fmt.Sprintf("daily log version conflict: expected %d, current %d", expectedVersion, currentVersion),
		Code:            models.DailyLogErrorConflict,
		CurrentVersion:  currentVersion,
		CurrentDailyLog: current,
	}
}

func mapLockError(err error, message string) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return notFoundError(message)
	}
	return err
}
