// FILE: apps/api/internal/atlas/graph/resolver/workout_api_contract_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: GraphQL executable-schema contract tests for WAVE-03 workout diary operations.
//   SCOPE: API-visible operation signatures, Date scalar binding, DailyLog result envelopes, DailyLog summary error propagation, and WAVE-03 no-scope schema proof; excludes repository and service internals.
//   DEPENDS: apps/api/internal/atlas/graph/generated, apps/api/internal/atlas/graph/resolver, apps/api/internal/atlas/models, gqlgen client.
//   LINKS: M-API / V-M-API / WAVE-03.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   TestWorkoutGraphQLSchema_OperationSignaturesAndNoWave04Placeholders - Proves WAVE-03 schema operation signatures and no cardio/body placeholders.
//   TestWorkoutGraphQLDailyLog_AuthNoCreateAndDateBinding - Proves executable schema Date binding, auth envelope, and no-create dailyLog query behavior.
//   TestWorkoutGraphQLDailyLogs_RangeSuccessMapping - Proves dailyLogs range query maps summary samples through generated schema.
//   TestWorkoutGraphQLDailyLogs_InvalidRangePropagatesValidationError - Proves dailyLogs exposes service validation errors instead of hiding them behind a null list.
//   TestWorkoutGraphQLMutations_ResultMappings - Proves generated mutation signatures map success, null, validation, conflict, and not-found envelopes.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added WAVE-03 GraphQL API contract coverage.
// END_CHANGE_SUMMARY

package resolver_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/atlas/graph/generated"
	"monorepo-template/apps/api/internal/atlas/graph/resolver"
	"monorepo-template/apps/api/internal/atlas/middleware"
	"monorepo-template/apps/api/internal/atlas/models"
)

func TestWorkoutGraphQLSchema_OperationSignaturesAndNoWave04Placeholders(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "schema", "workouts.graphql"))
	require.NoError(t, err)
	schema := string(raw)

	expectedSignatures := []string{
		"scalar Date",
		"dailyLog(date: Date!): DailyLogResult!",
		"dailyLogs(from: Date!, to: Date!): [DailyLogSummary!]!",
		"updateDailyLogNotes(date: Date!, expectedVersion: Int!, notes: String): DailyLogResult!",
		"addWorkoutExercise(date: Date!, expectedVersion: Int!, input: AddWorkoutExerciseInput!): DailyLogResult!",
		"updateWorkoutExercise(id: ID!, expectedVersion: Int!, input: UpdateWorkoutExerciseInput!): DailyLogResult!",
		"removeWorkoutExercise(id: ID!, expectedVersion: Int!): DailyLogResult!",
		"reorderWorkoutExercises(date: Date!, expectedVersion: Int!, orderedIds: [ID!]!): DailyLogResult!",
		"addWorkoutSet(workoutExerciseId: ID!, expectedVersion: Int!, input: AddWorkoutSetInput!): DailyLogResult!",
		"updateWorkoutSet(id: ID!, expectedVersion: Int!, input: UpdateWorkoutSetInput!): DailyLogResult!",
		"removeWorkoutSet(id: ID!, expectedVersion: Int!): DailyLogResult!",
		"reorderWorkoutSets(workoutExerciseId: ID!, expectedVersion: Int!, orderedIds: [ID!]!): DailyLogResult!",
	}
	for _, signature := range expectedSignatures {
		assert.Contains(t, schema, signature)
	}

	for _, forbidden := range []string{"cardio_entries", "CardioType", "HeartRateZone", "body_weight", "bodyWeight", "WorkoutDay"} {
		assert.NotContains(t, schema, forbidden)
	}
	assert.Equal(t, 0, strings.Count(schema, "cardio"), "schema must not expose fake cardio placeholders")
}

func TestWorkoutGraphQLDailyLog_AuthNoCreateAndDateBinding(t *testing.T) {
	query := `query DailyLogByDate($date: Date!) {
		dailyLog(date: $date) {
			dailyLog { id date version }
			authError { code message }
			validationError { code message }
		}
	}`

	t.Run("missing user context returns auth envelope", func(t *testing.T) {
		c := workoutGraphQLClient(&mockWorkoutService{
			getDailyLogFn: func(ctx context.Context, userID string, date models.Date) (*models.DailyLog, error) {
				t.Fatalf("service must not be called without Atlas user context")
				return nil, nil
			},
		})

		var response struct {
			DailyLog struct {
				DailyLog  *struct{ ID string } `json:"dailyLog"`
				AuthError *struct {
					Code    models.DailyLogErrorCode `json:"code"`
					Message string                   `json:"message"`
				} `json:"authError"`
				ValidationError *struct {
					Code models.DailyLogErrorCode `json:"code"`
				} `json:"validationError"`
			} `json:"dailyLog"`
		}
		err := c.Post(query, &response, client.Var("date", "2026-06-19"))

		require.NoError(t, err)
		assert.Nil(t, response.DailyLog.DailyLog)
		require.NotNil(t, response.DailyLog.AuthError)
		assert.Equal(t, models.DailyLogErrorAuth, response.DailyLog.AuthError.Code)
		assert.Equal(t, "unauthorized", response.DailyLog.AuthError.Message)
	})

	t.Run("absent day returns nil dailyLog without creating", func(t *testing.T) {
		called := false
		c := workoutGraphQLClient(&mockWorkoutService{
			getDailyLogFn: func(ctx context.Context, userID string, date models.Date) (*models.DailyLog, error) {
				called = true
				assert.Equal(t, "user-1", userID)
				assert.Equal(t, "2026-06-19", date.String())
				return nil, nil
			},
		})

		var response struct {
			DailyLog struct {
				DailyLog  *struct{ ID string } `json:"dailyLog"`
				AuthError *struct {
					Code models.DailyLogErrorCode `json:"code"`
				} `json:"authError"`
				ValidationError *struct {
					Code models.DailyLogErrorCode `json:"code"`
				} `json:"validationError"`
			} `json:"dailyLog"`
		}
		err := c.Post(query, &response, withAtlasUser("user-1"), client.Var("date", "2026-06-19"))

		require.NoError(t, err)
		assert.True(t, called)
		assert.Nil(t, response.DailyLog.DailyLog)
		assert.Nil(t, response.DailyLog.AuthError)
	})

	t.Run("timestamp date variable fails before resolver", func(t *testing.T) {
		called := false
		c := workoutGraphQLClient(&mockWorkoutService{
			getDailyLogFn: func(ctx context.Context, userID string, date models.Date) (*models.DailyLog, error) {
				called = true
				return nil, nil
			},
		})

		var response struct {
			DailyLog any `json:"dailyLog"`
		}
		err := c.Post(query, &response, withAtlasUser("user-1"), client.Var("date", "2026-06-19T10:00:00Z"))

		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid date")
		assert.False(t, called)
	})
}

func TestWorkoutGraphQLDailyLogs_RangeSuccessMapping(t *testing.T) {
	c := workoutGraphQLClient(&mockWorkoutService{
		listDailyLogSummariesFn: func(ctx context.Context, userID string, from models.Date, to models.Date) ([]models.DailyLogSummary, error) {
			assert.Equal(t, "user-1", userID)
			assert.Equal(t, "2026-06-18", from.String())
			assert.Equal(t, "2026-06-20", to.String())
			return []models.DailyLogSummary{{
				ID:                   "log-1",
				Date:                 models.MustDate("2026-06-19"),
				Version:              4,
				WorkoutExerciseCount: 2,
				WorkoutSetCount:      5,
				TotalVolume:          1234.5,
				UpdatedAt:            "2026-06-20T10:00:00Z",
			}}, nil
		},
	})

	var response struct {
		DailyLogs []struct {
			ID                   string  `json:"id"`
			Date                 string  `json:"date"`
			Version              int     `json:"version"`
			WorkoutExerciseCount int     `json:"workoutExerciseCount"`
			WorkoutSetCount      int     `json:"workoutSetCount"`
			TotalVolume          float64 `json:"totalVolume"`
		} `json:"dailyLogs"`
	}
	err := c.Post(
		`query DailyLogsRange($from: Date!, $to: Date!) {
			dailyLogs(from: $from, to: $to) {
				id
				date
				version
				workoutExerciseCount
				workoutSetCount
				totalVolume
			}
		}`,
		&response,
		withAtlasUser("user-1"),
		client.Var("from", "2026-06-18"),
		client.Var("to", "2026-06-20"),
	)

	require.NoError(t, err)
	require.Len(t, response.DailyLogs, 1)
	assert.Equal(t, "log-1", response.DailyLogs[0].ID)
	assert.Equal(t, "2026-06-19", response.DailyLogs[0].Date)
	assert.Equal(t, 4, response.DailyLogs[0].Version)
	assert.Equal(t, 2, response.DailyLogs[0].WorkoutExerciseCount)
	assert.Equal(t, 5, response.DailyLogs[0].WorkoutSetCount)
	assert.Equal(t, 1234.5, response.DailyLogs[0].TotalVolume)
}

func TestWorkoutGraphQLDailyLogs_InvalidRangePropagatesValidationError(t *testing.T) {
	c := workoutGraphQLClient(&mockWorkoutService{
		listDailyLogSummariesFn: func(ctx context.Context, userID string, from models.Date, to models.Date) ([]models.DailyLogSummary, error) {
			assert.Equal(t, "user-1", userID)
			assert.Equal(t, "2026-06-20", from.String())
			assert.Equal(t, "2026-06-19", to.String())
			return nil, &models.DailyLogValidationErr{
				Message: "from date must be on or before to date",
				Code:    models.DailyLogErrorValidation,
			}
		},
	})

	var response struct {
		DailyLogs []struct {
			ID string `json:"id"`
		} `json:"dailyLogs"`
	}
	err := c.Post(
		`query DailyLogsRange($from: Date!, $to: Date!) {
			dailyLogs(from: $from, to: $to) { id }
		}`,
		&response,
		withAtlasUser("user-1"),
		client.Var("from", "2026-06-20"),
		client.Var("to", "2026-06-19"),
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "from date must be on or before to date")
	assert.Nil(t, response.DailyLogs)
}

func TestWorkoutGraphQLMutations_ResultMappings(t *testing.T) {
	t.Run("updateDailyLogNotes maps explicit null to successful DailyLog", func(t *testing.T) {
		c := workoutGraphQLClient(&mockWorkoutService{
			updateDailyLogNotesFn: func(ctx context.Context, userID string, date models.Date, expectedVersion int32, notes *string) (*models.DailyLog, error) {
				assert.Equal(t, "user-1", userID)
				assert.Equal(t, "2026-06-19", date.String())
				assert.Equal(t, int32(2), expectedVersion)
				assert.Nil(t, notes)
				return &models.DailyLog{ID: "log-1", UserID: userID, Date: date, Version: 3}, nil
			},
		})

		var response struct {
			UpdateDailyLogNotes struct {
				DailyLog *struct {
					ID      string  `json:"id"`
					Date    string  `json:"date"`
					Version int     `json:"version"`
					Notes   *string `json:"notes"`
				} `json:"dailyLog"`
				ValidationError *struct {
					Code models.DailyLogErrorCode `json:"code"`
				} `json:"validationError"`
			} `json:"updateDailyLogNotes"`
		}
		err := c.Post(
			`mutation UpdateNotes($date: Date!, $expectedVersion: Int!, $notes: String) {
				updateDailyLogNotes(date: $date, expectedVersion: $expectedVersion, notes: $notes) {
					dailyLog { id date version notes }
					validationError { code message }
				}
			}`,
			&response,
			withAtlasUser("user-1"),
			client.Var("date", "2026-06-19"),
			client.Var("expectedVersion", 2),
			client.Var("notes", nil),
		)

		require.NoError(t, err)
		require.NotNil(t, response.UpdateDailyLogNotes.DailyLog)
		assert.Equal(t, "log-1", response.UpdateDailyLogNotes.DailyLog.ID)
		assert.Equal(t, "2026-06-19", response.UpdateDailyLogNotes.DailyLog.Date)
		assert.Equal(t, 3, response.UpdateDailyLogNotes.DailyLog.Version)
		assert.Nil(t, response.UpdateDailyLogNotes.DailyLog.Notes)
		assert.Nil(t, response.UpdateDailyLogNotes.ValidationError)
	})

	t.Run("negative expectedVersion maps validation envelope", func(t *testing.T) {
		c := workoutGraphQLClient(&mockWorkoutService{
			addWorkoutExerciseFn: func(ctx context.Context, userID string, date models.Date, expectedVersion int32, input models.AddWorkoutExerciseInput) (*models.DailyLog, error) {
				assert.Equal(t, "user-1", userID)
				assert.Equal(t, int32(-1), expectedVersion)
				assert.Equal(t, "exercise-1", input.ExerciseID)
				return nil, &models.DailyLogValidationErr{
					Message: "expectedVersion must be greater than or equal to 0",
					Code:    models.DailyLogErrorValidation,
				}
			},
		})

		var response struct {
			AddWorkoutExercise struct {
				DailyLog        *struct{ ID string } `json:"dailyLog"`
				ValidationError *struct {
					Code    models.DailyLogErrorCode `json:"code"`
					Message string                   `json:"message"`
				} `json:"validationError"`
			} `json:"addWorkoutExercise"`
		}
		err := c.Post(
			`mutation AddExercise($date: Date!, $expectedVersion: Int!, $input: AddWorkoutExerciseInput!) {
				addWorkoutExercise(date: $date, expectedVersion: $expectedVersion, input: $input) {
					dailyLog { id }
					validationError { code message }
				}
			}`,
			&response,
			withAtlasUser("user-1"),
			client.Var("date", "2026-06-19"),
			client.Var("expectedVersion", -1),
			client.Var("input", map[string]any{"exerciseId": "exercise-1"}),
		)

		require.NoError(t, err)
		assert.Nil(t, response.AddWorkoutExercise.DailyLog)
		require.NotNil(t, response.AddWorkoutExercise.ValidationError)
		assert.Equal(t, models.DailyLogErrorValidation, response.AddWorkoutExercise.ValidationError.Code)
		assert.Contains(t, response.AddWorkoutExercise.ValidationError.Message, "expectedVersion")
	})

	t.Run("stale notes update maps conflict envelope", func(t *testing.T) {
		c := workoutGraphQLClient(&mockWorkoutService{
			updateDailyLogNotesFn: func(ctx context.Context, userID string, date models.Date, expectedVersion int32, notes *string) (*models.DailyLog, error) {
				return nil, &models.DailyLogConflictErr{
					Message:        "daily log version conflict",
					Code:           models.DailyLogErrorConflict,
					CurrentVersion: 7,
					CurrentDailyLog: &models.DailyLog{
						ID:      "log-current",
						UserID:  userID,
						Date:    date,
						Version: 7,
					},
				}
			},
		})

		var response struct {
			UpdateDailyLogNotes struct {
				ConflictError *struct {
					Code            models.DailyLogErrorCode `json:"code"`
					CurrentVersion  int                      `json:"currentVersion"`
					CurrentDailyLog *struct {
						ID      string `json:"id"`
						Version int    `json:"version"`
					} `json:"currentDailyLog"`
				} `json:"conflictError"`
			} `json:"updateDailyLogNotes"`
		}
		err := c.Post(
			`mutation UpdateNotes($date: Date!, $expectedVersion: Int!, $notes: String) {
				updateDailyLogNotes(date: $date, expectedVersion: $expectedVersion, notes: $notes) {
					conflictError { code currentVersion currentDailyLog { id version } }
				}
			}`,
			&response,
			withAtlasUser("user-1"),
			client.Var("date", "2026-06-19"),
			client.Var("expectedVersion", 2),
			client.Var("notes", "new notes"),
		)

		require.NoError(t, err)
		require.NotNil(t, response.UpdateDailyLogNotes.ConflictError)
		assert.Equal(t, models.DailyLogErrorConflict, response.UpdateDailyLogNotes.ConflictError.Code)
		assert.Equal(t, 7, response.UpdateDailyLogNotes.ConflictError.CurrentVersion)
		require.NotNil(t, response.UpdateDailyLogNotes.ConflictError.CurrentDailyLog)
		assert.Equal(t, "log-current", response.UpdateDailyLogNotes.ConflictError.CurrentDailyLog.ID)
		assert.Equal(t, 7, response.UpdateDailyLogNotes.ConflictError.CurrentDailyLog.Version)
	})

	t.Run("addWorkoutSet missing parent maps not-found envelope", func(t *testing.T) {
		c := workoutGraphQLClient(&mockWorkoutService{
			addWorkoutSetFn: func(ctx context.Context, userID string, workoutExerciseID string, expectedVersion int32, input models.AddWorkoutSetInput) (*models.DailyLog, error) {
				assert.Equal(t, "user-1", userID)
				assert.Equal(t, "workout-exercise-missing", workoutExerciseID)
				assert.Equal(t, int32(1), expectedVersion)
				assert.Equal(t, 100.0, input.Weight)
				assert.Equal(t, int32(5), input.Reps)
				return nil, &models.DailyLogNotFoundErr{
					Message: "workout exercise not found",
					Code:    models.DailyLogErrorNotFound,
				}
			},
		})

		var response struct {
			AddWorkoutSet struct {
				DailyLog      *struct{ ID string } `json:"dailyLog"`
				NotFoundError *struct {
					Code    models.DailyLogErrorCode `json:"code"`
					Message string                   `json:"message"`
				} `json:"notFoundError"`
			} `json:"addWorkoutSet"`
		}
		err := c.Post(
			`mutation AddSet($workoutExerciseId: ID!, $expectedVersion: Int!, $input: AddWorkoutSetInput!) {
				addWorkoutSet(workoutExerciseId: $workoutExerciseId, expectedVersion: $expectedVersion, input: $input) {
					dailyLog { id }
					notFoundError { code message }
				}
			}`,
			&response,
			withAtlasUser("user-1"),
			client.Var("workoutExerciseId", "workout-exercise-missing"),
			client.Var("expectedVersion", 1),
			client.Var("input", map[string]any{"weight": 100.0, "reps": 5}),
		)

		require.NoError(t, err)
		assert.Nil(t, response.AddWorkoutSet.DailyLog)
		require.NotNil(t, response.AddWorkoutSet.NotFoundError)
		assert.Equal(t, models.DailyLogErrorNotFound, response.AddWorkoutSet.NotFoundError.Code)
		assert.Contains(t, response.AddWorkoutSet.NotFoundError.Message, "workout exercise")
	})
}

func workoutGraphQLClient(workoutService *mockWorkoutService) *client.Client {
	r := &resolver.Resolver{WorkoutService: workoutService}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: r}))
	return client.New(srv)
}

func withAtlasUser(userID string) client.Option {
	return func(req *client.Request) {
		ctx := middleware.ContextWithAtlasUserID(req.HTTP.Context(), userID)
		req.HTTP = req.HTTP.WithContext(ctx)
	}
}
