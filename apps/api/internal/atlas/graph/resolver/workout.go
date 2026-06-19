// FILE: apps/api/internal/atlas/graph/resolver/workout.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement GraphQL resolvers for WAVE-03 DailyLog strength workout diary operations.
//   SCOPE: DailyLog reads, summaries, notes updates, workout exercise mutations, workout set mutations, auth handling, service error mapping, and GraphQL omittable update input mapping; excludes repository, service validation internals, UI, analytics, and export work.
//   DEPENDS: apps/api/internal/atlas/middleware, apps/api/internal/atlas/models, apps/api/internal/atlas/service.WorkoutService, generated workouts.resolvers.go.
//   LINKS: M-API / V-M-API / WAVE-03.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   GetDailyLog - Reads a DailyLog result envelope for the authenticated Atlas user.
//   DailyLogs - Lists DailyLog summaries for the authenticated Atlas user.
//   UpdateDailyLogNotes - Updates DailyLog notes through WorkoutService and maps typed result errors.
//   AddWorkoutExercise - Adds a strength exercise instance through WorkoutService.
//   UpdateWorkoutExercise - Maps GraphQL omittable fields into service update input and delegates.
//   RemoveWorkoutExercise - Removes a workout exercise through WorkoutService.
//   ReorderWorkoutExercises - Reorders workout exercises through WorkoutService.
//   AddWorkoutSet - Adds a workout set through WorkoutService.
//   UpdateWorkoutSet - Maps GraphQL omittable fields into service update input and delegates.
//   RemoveWorkoutSet - Removes a workout set through WorkoutService.
//   ReorderWorkoutSets - Reorders workout sets through WorkoutService.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added WAVE-03 workout diary resolver methods and typed DailyLogResult mapping.
// END_CHANGE_SUMMARY

package resolver

import (
	"context"
	"errors"

	"monorepo-template/apps/api/internal/atlas/middleware"
	"monorepo-template/apps/api/internal/atlas/models"
)

func (r *Resolver) GetDailyLog(ctx context.Context, date models.Date) (*models.DailyLogResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return unauthorizedDailyLogResult(), nil
	}

	log, err := r.WorkoutService.GetDailyLog(ctx, userID, date)
	return dailyLogResult(log, err), nil
}

func (r *Resolver) DailyLogs(ctx context.Context, from models.Date, to models.Date) ([]*models.DailyLogSummary, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return nil, nil
	}

	summaries, err := r.WorkoutService.ListDailyLogSummaries(ctx, userID, from, to)
	if err != nil {
		return nil, nil
	}

	out := make([]*models.DailyLogSummary, len(summaries))
	for i := range summaries {
		out[i] = &summaries[i]
	}
	return out, nil
}

func (r *Resolver) UpdateDailyLogNotes(ctx context.Context, date models.Date, expectedVersion int, notes *string) (*models.DailyLogResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return unauthorizedDailyLogResult(), nil
	}

	log, err := r.WorkoutService.UpdateDailyLogNotes(ctx, userID, date, int32(expectedVersion), notes)
	return dailyLogResult(log, err), nil
}

func (r *Resolver) AddWorkoutExercise(ctx context.Context, date models.Date, expectedVersion int, input models.AddWorkoutExerciseInput) (*models.DailyLogResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return unauthorizedDailyLogResult(), nil
	}

	log, err := r.WorkoutService.AddWorkoutExercise(ctx, userID, date, int32(expectedVersion), input)
	return dailyLogResult(log, err), nil
}

func (r *Resolver) UpdateWorkoutExercise(ctx context.Context, id string, expectedVersion int, input models.UpdateWorkoutExerciseGraphQLInput) (*models.DailyLogResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return unauthorizedDailyLogResult(), nil
	}

	serviceInput := models.UpdateWorkoutExerciseInput{}
	if input.Position.IsSet() {
		serviceInput.Position = input.Position.Value()
	}
	if input.Notes.IsSet() {
		serviceInput.SetNotes = true
		serviceInput.Notes = input.Notes.Value()
	}

	log, err := r.WorkoutService.UpdateWorkoutExercise(ctx, userID, id, int32(expectedVersion), serviceInput)
	return dailyLogResult(log, err), nil
}

func (r *Resolver) RemoveWorkoutExercise(ctx context.Context, id string, expectedVersion int) (*models.DailyLogResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return unauthorizedDailyLogResult(), nil
	}

	log, err := r.WorkoutService.RemoveWorkoutExercise(ctx, userID, id, int32(expectedVersion))
	return dailyLogResult(log, err), nil
}

func (r *Resolver) ReorderWorkoutExercises(ctx context.Context, date models.Date, expectedVersion int, orderedIDs []string) (*models.DailyLogResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return unauthorizedDailyLogResult(), nil
	}

	log, err := r.WorkoutService.ReorderWorkoutExercises(ctx, userID, date, int32(expectedVersion), orderedIDs)
	return dailyLogResult(log, err), nil
}

func (r *Resolver) AddWorkoutSet(ctx context.Context, workoutExerciseID string, expectedVersion int, input models.AddWorkoutSetInput) (*models.DailyLogResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return unauthorizedDailyLogResult(), nil
	}

	log, err := r.WorkoutService.AddWorkoutSet(ctx, userID, workoutExerciseID, int32(expectedVersion), input)
	return dailyLogResult(log, err), nil
}

func (r *Resolver) UpdateWorkoutSet(ctx context.Context, id string, expectedVersion int, input models.UpdateWorkoutSetGraphQLInput) (*models.DailyLogResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return unauthorizedDailyLogResult(), nil
	}

	serviceInput := models.UpdateWorkoutSetInput{}
	if input.SetNumber.IsSet() {
		serviceInput.SetNumber = input.SetNumber.Value()
	}
	if input.Weight.IsSet() {
		serviceInput.Weight = input.Weight.Value()
	}
	if input.Reps.IsSet() {
		serviceInput.Reps = input.Reps.Value()
	}
	if input.RPE.IsSet() {
		serviceInput.SetRPE = true
		serviceInput.RPE = input.RPE.Value()
	}
	if input.RIR.IsSet() {
		serviceInput.SetRIR = true
		serviceInput.RIR = input.RIR.Value()
	}
	if input.Notes.IsSet() {
		serviceInput.SetNotes = true
		serviceInput.Notes = input.Notes.Value()
	}

	log, err := r.WorkoutService.UpdateWorkoutSet(ctx, userID, id, int32(expectedVersion), serviceInput)
	return dailyLogResult(log, err), nil
}

func (r *Resolver) RemoveWorkoutSet(ctx context.Context, id string, expectedVersion int) (*models.DailyLogResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return unauthorizedDailyLogResult(), nil
	}

	log, err := r.WorkoutService.RemoveWorkoutSet(ctx, userID, id, int32(expectedVersion))
	return dailyLogResult(log, err), nil
}

func (r *Resolver) ReorderWorkoutSets(ctx context.Context, workoutExerciseID string, expectedVersion int, orderedIDs []string) (*models.DailyLogResult, error) {
	userID := middleware.GetAtlasUserID(ctx)
	if userID == "" {
		return unauthorizedDailyLogResult(), nil
	}

	log, err := r.WorkoutService.ReorderWorkoutSets(ctx, userID, workoutExerciseID, int32(expectedVersion), orderedIDs)
	return dailyLogResult(log, err), nil
}

func unauthorizedDailyLogResult() *models.DailyLogResult {
	return &models.DailyLogResult{
		AuthErr: &models.DailyLogAuthErr{
			Message: "unauthorized",
			Code:    models.DailyLogErrorAuth,
		},
	}
}

func dailyLogResult(log *models.DailyLog, err error) *models.DailyLogResult {
	if err == nil {
		return &models.DailyLogResult{DailyLog: log}
	}

	var validationErr *models.DailyLogValidationErr
	if errors.As(err, &validationErr) {
		return &models.DailyLogResult{ValidationErr: validationErr}
	}

	var notFoundErr *models.DailyLogNotFoundErr
	if errors.As(err, &notFoundErr) {
		return &models.DailyLogResult{NotFoundErr: notFoundErr}
	}

	var conflictErr *models.DailyLogConflictErr
	if errors.As(err, &conflictErr) {
		return &models.DailyLogResult{ConflictErr: conflictErr}
	}

	var authErr *models.DailyLogAuthErr
	if errors.As(err, &authErr) {
		return &models.DailyLogResult{AuthErr: authErr}
	}

	return nil
}
