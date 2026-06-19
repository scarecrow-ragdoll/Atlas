// FILE: apps/api/internal/atlas/graph/resolver/workout_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for WAVE-03 workout diary GraphQL resolvers.
//   SCOPE: Auth handling, WorkoutService delegation, DailyLogResult error mapping, and GraphQL omittable update input mapping; excludes repository, service validation internals, UI, analytics, and export work.
//   DEPENDS: apps/api/internal/atlas/graph/resolver, apps/api/internal/atlas/models, apps/api/internal/atlas/service.
//   LINKS: M-API / V-M-API / WAVE-03.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   mockWorkoutService - Captures resolver calls and returns configured service outcomes.
//   TestDailyLogResolver_UnauthorizedReturnsAuthError - Proves missing Atlas user context returns typed auth error.
//   TestDailyLogResolver_DelegatesAuthenticatedDailyLog - Proves authenticated dailyLog delegates with user and date.
//   TestUpdateDailyLogNotesResolver_MapsConflictError - Proves conflict service errors map to DailyLogResult conflictError without Go error.
//   TestAddWorkoutExerciseResolver_MapsValidationError - Proves validation service errors map to DailyLogResult validationError.
//   TestWorkoutSetResolvers_MapNotFoundError - Proves workout set mutation not-found errors map to DailyLogResult notFoundError.
//   TestWorkoutResolvers_DoNotLeakUnexpectedErrors - Proves unexpected service errors do not leak raw internal text.
//   TestUpdateWorkoutExerciseResolver_MapsExplicitNullNotes - Proves explicit null notes are forwarded as SetNotes with nil Notes.
//   TestUpdateWorkoutSetResolver_MapsExplicitNullNullableFields - Proves explicit null RPE, RIR, and notes set the corresponding service flags.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added WAVE-03 workout resolver TDD coverage.
// END_CHANGE_SUMMARY

package resolver_test

import (
	"context"
	"errors"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/atlas/graph/resolver"
	"monorepo-template/apps/api/internal/atlas/models"
)

type mockWorkoutService struct {
	getDailyLogFn             func(ctx context.Context, userID string, date models.Date) (*models.DailyLog, error)
	listDailyLogSummariesFn   func(ctx context.Context, userID string, from models.Date, to models.Date) ([]models.DailyLogSummary, error)
	updateDailyLogNotesFn     func(ctx context.Context, userID string, date models.Date, expectedVersion int32, notes *string) (*models.DailyLog, error)
	addWorkoutExerciseFn      func(ctx context.Context, userID string, date models.Date, expectedVersion int32, input models.AddWorkoutExerciseInput) (*models.DailyLog, error)
	updateWorkoutExerciseFn   func(ctx context.Context, userID string, id string, expectedVersion int32, input models.UpdateWorkoutExerciseInput) (*models.DailyLog, error)
	removeWorkoutExerciseFn   func(ctx context.Context, userID string, id string, expectedVersion int32) (*models.DailyLog, error)
	reorderWorkoutExercisesFn func(ctx context.Context, userID string, date models.Date, expectedVersion int32, orderedIDs []string) (*models.DailyLog, error)
	addWorkoutSetFn           func(ctx context.Context, userID string, workoutExerciseID string, expectedVersion int32, input models.AddWorkoutSetInput) (*models.DailyLog, error)
	updateWorkoutSetFn        func(ctx context.Context, userID string, id string, expectedVersion int32, input models.UpdateWorkoutSetInput) (*models.DailyLog, error)
	removeWorkoutSetFn        func(ctx context.Context, userID string, id string, expectedVersion int32) (*models.DailyLog, error)
	reorderWorkoutSetsFn      func(ctx context.Context, userID string, workoutExerciseID string, expectedVersion int32, orderedIDs []string) (*models.DailyLog, error)
}

func (m *mockWorkoutService) GetDailyLog(ctx context.Context, userID string, date models.Date) (*models.DailyLog, error) {
	return m.getDailyLogFn(ctx, userID, date)
}

func (m *mockWorkoutService) ListDailyLogSummaries(ctx context.Context, userID string, from models.Date, to models.Date) ([]models.DailyLogSummary, error) {
	return m.listDailyLogSummariesFn(ctx, userID, from, to)
}

func (m *mockWorkoutService) UpdateDailyLogNotes(ctx context.Context, userID string, date models.Date, expectedVersion int32, notes *string) (*models.DailyLog, error) {
	return m.updateDailyLogNotesFn(ctx, userID, date, expectedVersion, notes)
}

func (m *mockWorkoutService) AddWorkoutExercise(ctx context.Context, userID string, date models.Date, expectedVersion int32, input models.AddWorkoutExerciseInput) (*models.DailyLog, error) {
	return m.addWorkoutExerciseFn(ctx, userID, date, expectedVersion, input)
}

func (m *mockWorkoutService) UpdateWorkoutExercise(ctx context.Context, userID string, id string, expectedVersion int32, input models.UpdateWorkoutExerciseInput) (*models.DailyLog, error) {
	return m.updateWorkoutExerciseFn(ctx, userID, id, expectedVersion, input)
}

func (m *mockWorkoutService) RemoveWorkoutExercise(ctx context.Context, userID string, id string, expectedVersion int32) (*models.DailyLog, error) {
	return m.removeWorkoutExerciseFn(ctx, userID, id, expectedVersion)
}

func (m *mockWorkoutService) ReorderWorkoutExercises(ctx context.Context, userID string, date models.Date, expectedVersion int32, orderedIDs []string) (*models.DailyLog, error) {
	return m.reorderWorkoutExercisesFn(ctx, userID, date, expectedVersion, orderedIDs)
}

func (m *mockWorkoutService) AddWorkoutSet(ctx context.Context, userID string, workoutExerciseID string, expectedVersion int32, input models.AddWorkoutSetInput) (*models.DailyLog, error) {
	return m.addWorkoutSetFn(ctx, userID, workoutExerciseID, expectedVersion, input)
}

func (m *mockWorkoutService) UpdateWorkoutSet(ctx context.Context, userID string, id string, expectedVersion int32, input models.UpdateWorkoutSetInput) (*models.DailyLog, error) {
	return m.updateWorkoutSetFn(ctx, userID, id, expectedVersion, input)
}

func (m *mockWorkoutService) RemoveWorkoutSet(ctx context.Context, userID string, id string, expectedVersion int32) (*models.DailyLog, error) {
	return m.removeWorkoutSetFn(ctx, userID, id, expectedVersion)
}

func (m *mockWorkoutService) ReorderWorkoutSets(ctx context.Context, userID string, workoutExerciseID string, expectedVersion int32, orderedIDs []string) (*models.DailyLog, error) {
	return m.reorderWorkoutSetsFn(ctx, userID, workoutExerciseID, expectedVersion, orderedIDs)
}

func TestDailyLogResolver_UnauthorizedReturnsAuthError(t *testing.T) {
	r := &resolver.Resolver{
		WorkoutService: &mockWorkoutService{},
	}

	result, err := r.GetDailyLog(context.Background(), models.MustDate("2026-06-19"))

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.DailyLog)
	require.NotNil(t, result.AuthErr)
	assert.Equal(t, "unauthorized", result.AuthErr.Message)
	assert.Equal(t, models.DailyLogErrorAuth, result.AuthErr.Code)
}

func TestDailyLogResolver_DelegatesAuthenticatedDailyLog(t *testing.T) {
	date := models.MustDate("2026-06-19")
	expected := &models.DailyLog{ID: "log-1", UserID: "user-1", Date: date, Version: 3}
	r := &resolver.Resolver{
		WorkoutService: &mockWorkoutService{
			getDailyLogFn: func(ctx context.Context, userID string, gotDate models.Date) (*models.DailyLog, error) {
				assert.Equal(t, "user-1", userID)
				assert.Equal(t, date.String(), gotDate.String())
				return expected, nil
			},
		},
	}

	result, err := r.GetDailyLog(userCtx("user-1"), date)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expected, result.DailyLog)
	assert.Nil(t, result.AuthErr)
	assert.Nil(t, result.ValidationErr)
	assert.Nil(t, result.NotFoundErr)
	assert.Nil(t, result.ConflictErr)
}

func TestUpdateDailyLogNotesResolver_MapsConflictError(t *testing.T) {
	current := &models.DailyLog{ID: "log-current", Version: 7}
	r := &resolver.Resolver{
		WorkoutService: &mockWorkoutService{
			updateDailyLogNotesFn: func(ctx context.Context, userID string, date models.Date, expectedVersion int32, notes *string) (*models.DailyLog, error) {
				assert.Equal(t, "user-1", userID)
				assert.Equal(t, int32(2), expectedVersion)
				return nil, &models.DailyLogConflictErr{
					Message:         "daily log version conflict",
					Code:            models.DailyLogErrorConflict,
					CurrentVersion:  7,
					CurrentDailyLog: current,
				}
			},
		},
	}

	result, err := r.UpdateDailyLogNotes(userCtx("user-1"), models.MustDate("2026-06-19"), 2, stringPtr("notes"))

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.DailyLog)
	require.NotNil(t, result.ConflictErr)
	assert.Equal(t, models.DailyLogErrorConflict, result.ConflictErr.Code)
	assert.Equal(t, int32(7), result.ConflictErr.CurrentVersion)
	assert.Equal(t, current, result.ConflictErr.CurrentDailyLog)
}

func TestAddWorkoutExerciseResolver_MapsValidationError(t *testing.T) {
	r := &resolver.Resolver{
		WorkoutService: &mockWorkoutService{
			addWorkoutExerciseFn: func(ctx context.Context, userID string, date models.Date, expectedVersion int32, input models.AddWorkoutExerciseInput) (*models.DailyLog, error) {
				assert.Equal(t, "user-1", userID)
				return nil, &models.DailyLogValidationErr{
					Message: "position must be greater than 0",
					Code:    models.DailyLogErrorValidation,
				}
			},
		},
	}

	result, err := r.AddWorkoutExercise(userCtx("user-1"), models.MustDate("2026-06-19"), 1, models.AddWorkoutExerciseInput{ExerciseID: "exercise-1"})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.DailyLog)
	require.NotNil(t, result.ValidationErr)
	assert.Equal(t, models.DailyLogErrorValidation, result.ValidationErr.Code)
	assert.Contains(t, result.ValidationErr.Message, "position")
}

func TestWorkoutSetResolvers_MapNotFoundError(t *testing.T) {
	t.Run("add set", func(t *testing.T) {
		r := &resolver.Resolver{
			WorkoutService: &mockWorkoutService{
				addWorkoutSetFn: func(ctx context.Context, userID string, workoutExerciseID string, expectedVersion int32, input models.AddWorkoutSetInput) (*models.DailyLog, error) {
					assert.Equal(t, "workout-exercise-1", workoutExerciseID)
					return nil, &models.DailyLogNotFoundErr{Message: "workout exercise not found", Code: models.DailyLogErrorNotFound}
				},
			},
		}

		result, err := r.AddWorkoutSet(userCtx("user-1"), "workout-exercise-1", 1, models.AddWorkoutSetInput{Weight: 100, Reps: 5})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.NotFoundErr)
		assert.Equal(t, models.DailyLogErrorNotFound, result.NotFoundErr.Code)
		assert.Nil(t, result.DailyLog)
	})

	t.Run("update set", func(t *testing.T) {
		r := &resolver.Resolver{
			WorkoutService: &mockWorkoutService{
				updateWorkoutSetFn: func(ctx context.Context, userID string, id string, expectedVersion int32, input models.UpdateWorkoutSetInput) (*models.DailyLog, error) {
					assert.Equal(t, "set-1", id)
					return nil, &models.DailyLogNotFoundErr{Message: "workout set not found", Code: models.DailyLogErrorNotFound}
				},
			},
		}

		result, err := r.UpdateWorkoutSet(userCtx("user-1"), "set-1", 1, models.UpdateWorkoutSetGraphQLInput{})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.NotFoundErr)
		assert.Equal(t, models.DailyLogErrorNotFound, result.NotFoundErr.Code)
		assert.Nil(t, result.DailyLog)
	})

	t.Run("remove set", func(t *testing.T) {
		r := &resolver.Resolver{
			WorkoutService: &mockWorkoutService{
				removeWorkoutSetFn: func(ctx context.Context, userID string, id string, expectedVersion int32) (*models.DailyLog, error) {
					assert.Equal(t, "set-1", id)
					return nil, &models.DailyLogNotFoundErr{Message: "workout set not found", Code: models.DailyLogErrorNotFound}
				},
			},
		}

		result, err := r.RemoveWorkoutSet(userCtx("user-1"), "set-1", 1)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.NotFoundErr)
		assert.Equal(t, models.DailyLogErrorNotFound, result.NotFoundErr.Code)
		assert.Nil(t, result.DailyLog)
	})
}

func TestWorkoutResolvers_DoNotLeakUnexpectedErrors(t *testing.T) {
	r := &resolver.Resolver{
		WorkoutService: &mockWorkoutService{
			updateDailyLogNotesFn: func(ctx context.Context, userID string, date models.Date, expectedVersion int32, notes *string) (*models.DailyLog, error) {
				return nil, errors.New("pq: duplicate key raw internal detail")
			},
		},
	}

	result, err := r.UpdateDailyLogNotes(userCtx("user-1"), models.MustDate("2026-06-19"), 1, stringPtr("notes"))

	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestUpdateWorkoutExerciseResolver_MapsExplicitNullNotes(t *testing.T) {
	r := &resolver.Resolver{
		WorkoutService: &mockWorkoutService{
			updateWorkoutExerciseFn: func(ctx context.Context, userID string, id string, expectedVersion int32, input models.UpdateWorkoutExerciseInput) (*models.DailyLog, error) {
				assert.Equal(t, "user-1", userID)
				assert.Equal(t, "workout-exercise-1", id)
				assert.Equal(t, int32(4), expectedVersion)
				assert.True(t, input.SetNotes)
				assert.Nil(t, input.Notes)
				assert.Nil(t, input.Position)
				return &models.DailyLog{ID: "log-1"}, nil
			},
		},
	}

	input := models.UpdateWorkoutExerciseGraphQLInput{
		Notes: graphql.OmittableOf[*string](nil),
	}
	result, err := r.UpdateWorkoutExercise(userCtx("user-1"), "workout-exercise-1", 4, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.DailyLog)
	assert.Equal(t, "log-1", result.DailyLog.ID)
}

func TestUpdateWorkoutSetResolver_MapsExplicitNullNullableFields(t *testing.T) {
	r := &resolver.Resolver{
		WorkoutService: &mockWorkoutService{
			updateWorkoutSetFn: func(ctx context.Context, userID string, id string, expectedVersion int32, input models.UpdateWorkoutSetInput) (*models.DailyLog, error) {
				assert.Equal(t, "user-1", userID)
				assert.Equal(t, "set-1", id)
				assert.Equal(t, int32(5), expectedVersion)
				assert.True(t, input.SetRPE)
				assert.Nil(t, input.RPE)
				assert.True(t, input.SetRIR)
				assert.Nil(t, input.RIR)
				assert.True(t, input.SetNotes)
				assert.Nil(t, input.Notes)
				assert.Nil(t, input.SetNumber)
				assert.Nil(t, input.Weight)
				assert.Nil(t, input.Reps)
				return &models.DailyLog{ID: "log-1"}, nil
			},
		},
	}

	input := models.UpdateWorkoutSetGraphQLInput{
		RPE:   graphql.OmittableOf[*float64](nil),
		RIR:   graphql.OmittableOf[*int32](nil),
		Notes: graphql.OmittableOf[*string](nil),
	}
	result, err := r.UpdateWorkoutSet(userCtx("user-1"), "set-1", 5, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.DailyLog)
	assert.Equal(t, "log-1", result.DailyLog.ID)
}

func stringPtr(value string) *string {
	return &value
}
