// FILE: apps/api/internal/repository/postgres/workout_repo_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify WAVE-03 workout repository behavior against the goose-managed PostgreSQL test database.
//   SCOPE: DailyLog get/create/user isolation/versioning, workout exercise duplicates/snapshots/reordering/deletion, workout set constraints/reordering, and safe destructive setup.
//   DEPENDS: apps/api/internal/atlas/repository/postgres, apps/api/internal/atlas/models, apps/api/internal/repository/postgres, apps/api/internal/testinfra.
//   LINKS: M-API / V-M-API / WAVE-03.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   TestWorkoutRepo_* - Real database coverage for workout repository aggregate behavior.
//   workoutRepoTestSetup - Applies migrations, enforces safe test DSN, truncates WAVE-03 tables, and creates an Atlas user.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added WAVE-03 workout repository integration coverage.
// END_CHANGE_SUMMARY

package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	atlasModels "monorepo-template/apps/api/internal/atlas/models"
	atlasRepo "monorepo-template/apps/api/internal/atlas/repository/postgres"
	postgresrepo "monorepo-template/apps/api/internal/repository/postgres"
	"monorepo-template/apps/api/internal/testinfra"
)

func TestWorkoutRepo_GetOrCreateDailyLog_UniquePerUserDate(t *testing.T) {
	pool, repo, userID := workoutRepoTestSetup(t)
	date := atlasModels.MustDate("2026-06-19")

	first, err := repo.GetOrCreateDailyLogByDate(context.Background(), userID, date)
	require.NoError(t, err)
	require.NotNil(t, first)

	second, err := repo.GetOrCreateDailyLogByDate(context.Background(), userID, date)
	require.NoError(t, err)
	require.NotNil(t, second)

	assert.Equal(t, first.ID, second.ID)
	assert.Equal(t, int32(0), first.Version)
	assert.Equal(t, "2026-06-19", second.Date.String())

	var count int
	err = pool.QueryRow(context.Background(), `
		SELECT COUNT(*) FROM daily_logs WHERE user_id = $1::uuid AND date = $2::date
	`, userID, date.String()).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestWorkoutRepo_DailyLog_UserScopedIsolation(t *testing.T) {
	pool, repo, userID := workoutRepoTestSetup(t)
	otherUserID := ensureWorkoutAtlasUser(t, pool, "other-workout-user")
	date := atlasModels.MustDate("2026-06-20")

	mine, err := repo.GetOrCreateDailyLogByDate(context.Background(), userID, date)
	require.NoError(t, err)
	theirs, err := repo.GetOrCreateDailyLogByDate(context.Background(), otherUserID, date)
	require.NoError(t, err)

	assert.NotEqual(t, mine.ID, theirs.ID)

	foundMine, err := repo.GetDailyLogByDate(context.Background(), userID, date)
	require.NoError(t, err)
	require.NotNil(t, foundMine)
	assert.Equal(t, mine.ID, foundMine.ID)

	crossUserAggregate, err := repo.GetDailyLogAggregate(context.Background(), otherUserID, mine.ID)
	require.NoError(t, err)
	assert.Nil(t, crossUserAggregate)
}

func TestWorkoutRepo_AddWorkoutExercise_AllowsDuplicateExercise(t *testing.T) {
	pool, repo, userID := workoutRepoTestSetup(t)
	dailyLog := mustWorkoutDailyLog(t, repo, userID, "2026-06-21")
	exerciseID := seedWorkoutExerciseRecord(t, pool, userID, "Bench Press", nil)

	first, err := repo.AddWorkoutExercise(context.Background(), userID, dailyLog.ID, atlasRepo.AddWorkoutExerciseInput{
		ExerciseID: exerciseID,
		Position:   1,
	})
	require.NoError(t, err)
	second, err := repo.AddWorkoutExercise(context.Background(), userID, dailyLog.ID, atlasRepo.AddWorkoutExerciseInput{
		ExerciseID: exerciseID,
		Position:   2,
	})
	require.NoError(t, err)

	assert.NotEqual(t, first.ID, second.ID)
	assert.Equal(t, exerciseID, first.ExerciseID)
	assert.Equal(t, exerciseID, second.ExerciseID)

	aggregate, err := repo.GetDailyLogAggregate(context.Background(), userID, dailyLog.ID)
	require.NoError(t, err)
	require.NotNil(t, aggregate)
	require.Len(t, aggregate.WorkoutExercises, 2)
	assert.Equal(t, []int32{1, 2}, workoutExercisePositions(aggregate.WorkoutExercises))
}

func TestWorkoutRepo_AddWorkoutExercise_CapturesWorkingWeightSnapshot(t *testing.T) {
	pool, repo, userID := workoutRepoTestSetup(t)
	dailyLog := mustWorkoutDailyLog(t, repo, userID, "2026-06-22")
	snapshot := 72.5
	exerciseID := seedWorkoutExerciseRecord(t, pool, userID, "Back Squat", &snapshot)

	created, err := repo.AddWorkoutExercise(context.Background(), userID, dailyLog.ID, atlasRepo.AddWorkoutExerciseInput{
		ExerciseID:            exerciseID,
		Position:              1,
		WorkingWeightSnapshot: &snapshot,
		Notes:                 workoutStringPtr("felt strong"),
	})
	require.NoError(t, err)

	require.NotNil(t, created.WorkingWeightSnapshot)
	assert.InDelta(t, snapshot, *created.WorkingWeightSnapshot, 0.001)

	aggregate, err := repo.GetDailyLogAggregate(context.Background(), userID, dailyLog.ID)
	require.NoError(t, err)
	require.NotNil(t, aggregate)
	require.Len(t, aggregate.WorkoutExercises, 1)
	require.NotNil(t, aggregate.WorkoutExercises[0].WorkingWeightSnapshot)
	assert.InDelta(t, snapshot, *aggregate.WorkoutExercises[0].WorkingWeightSnapshot, 0.001)
}

func TestWorkoutRepo_ReorderWorkoutExercises_ReindexesContiguously(t *testing.T) {
	pool, repo, userID := workoutRepoTestSetup(t)
	dailyLog := mustWorkoutDailyLog(t, repo, userID, "2026-06-23")
	first := mustAddWorkoutExercise(t, repo, pool, userID, dailyLog.ID, "A", 1)
	second := mustAddWorkoutExercise(t, repo, pool, userID, dailyLog.ID, "B", 2)
	third := mustAddWorkoutExercise(t, repo, pool, userID, dailyLog.ID, "C", 3)

	err := repo.ReorderWorkoutExercises(context.Background(), userID, dailyLog.ID, []string{third.ID, first.ID, second.ID})
	require.NoError(t, err)

	aggregate, err := repo.GetDailyLogAggregate(context.Background(), userID, dailyLog.ID)
	require.NoError(t, err)
	require.NotNil(t, aggregate)
	require.Len(t, aggregate.WorkoutExercises, 3)
	assert.Equal(t, []string{third.ID, first.ID, second.ID}, workoutExerciseIDs(aggregate.WorkoutExercises))
	assert.Equal(t, []int32{1, 2, 3}, workoutExercisePositions(aggregate.WorkoutExercises))
}

func TestWorkoutRepo_DeleteWorkoutExercise_CascadesSetsAndKeepsDailyLog(t *testing.T) {
	pool, repo, userID := workoutRepoTestSetup(t)
	dailyLog := mustWorkoutDailyLog(t, repo, userID, "2026-06-24")
	deletedExercise := mustAddWorkoutExercise(t, repo, pool, userID, dailyLog.ID, "Delete Me", 1)
	remainingExercise := mustAddWorkoutExercise(t, repo, pool, userID, dailyLog.ID, "Keep Me", 2)
	_ = mustAddWorkoutSet(t, repo, userID, deletedExercise.ID, 1, 100, 5)
	_ = mustAddWorkoutSet(t, repo, userID, deletedExercise.ID, 2, 110, 3)

	deleted, err := repo.DeleteWorkoutExercise(context.Background(), userID, deletedExercise.ID)
	require.NoError(t, err)
	require.NotNil(t, deleted)

	var deletedSets int
	err = pool.QueryRow(context.Background(), `
		SELECT COUNT(*) FROM workout_sets WHERE workout_exercise_id = $1::uuid
	`, deletedExercise.ID).Scan(&deletedSets)
	require.NoError(t, err)
	assert.Equal(t, 0, deletedSets)

	aggregate, err := repo.GetDailyLogAggregate(context.Background(), userID, dailyLog.ID)
	require.NoError(t, err)
	require.NotNil(t, aggregate)
	assert.Equal(t, dailyLog.ID, aggregate.DailyLog.ID)
	require.Len(t, aggregate.WorkoutExercises, 1)
	assert.Equal(t, remainingExercise.ID, aggregate.WorkoutExercises[0].ID)
	assert.Equal(t, int32(1), aggregate.WorkoutExercises[0].Position)
}

func TestWorkoutRepo_AddWorkoutSet_ValidatesDBConstraints(t *testing.T) {
	pool, repo, userID := workoutRepoTestSetup(t)
	dailyLog := mustWorkoutDailyLog(t, repo, userID, "2026-06-25")
	exercise := mustAddWorkoutExercise(t, repo, pool, userID, dailyLog.ID, "Constraint Exercise", 1)
	_ = mustAddWorkoutSet(t, repo, userID, exercise.ID, 1, 100, 5)

	_, err := repo.AddWorkoutSet(context.Background(), userID, exercise.ID, atlasRepo.AddWorkoutSetInput{
		SetNumber: 1,
		Weight:    0,
		Reps:      5,
	})

	requireWorkoutCheckViolation(t, err, "chk_workout_sets_weight")

	aggregate, err := repo.GetDailyLogAggregate(context.Background(), userID, dailyLog.ID)
	require.NoError(t, err)
	require.NotNil(t, aggregate)
	require.Len(t, aggregate.WorkoutExercises, 1)
	require.Len(t, aggregate.WorkoutExercises[0].Sets, 1)
	assert.Equal(t, int32(1), aggregate.WorkoutExercises[0].Sets[0].SetNumber)
}

func TestWorkoutRepo_ReorderWorkoutSets_ReindexesContiguously(t *testing.T) {
	pool, repo, userID := workoutRepoTestSetup(t)
	dailyLog := mustWorkoutDailyLog(t, repo, userID, "2026-06-26")
	exercise := mustAddWorkoutExercise(t, repo, pool, userID, dailyLog.ID, "Rows", 1)
	first := mustAddWorkoutSet(t, repo, userID, exercise.ID, 1, 100, 5)
	second := mustAddWorkoutSet(t, repo, userID, exercise.ID, 2, 105, 4)
	third := mustAddWorkoutSet(t, repo, userID, exercise.ID, 3, 110, 3)

	err := repo.ReorderWorkoutSets(context.Background(), userID, exercise.ID, []string{third.ID, first.ID, second.ID})
	require.NoError(t, err)

	aggregate, err := repo.GetDailyLogAggregate(context.Background(), userID, dailyLog.ID)
	require.NoError(t, err)
	require.NotNil(t, aggregate)
	require.Len(t, aggregate.WorkoutExercises, 1)
	require.Len(t, aggregate.WorkoutExercises[0].Sets, 3)
	assert.Equal(t, []string{third.ID, first.ID, second.ID}, workoutSetIDs(aggregate.WorkoutExercises[0].Sets))
	assert.Equal(t, []int32{1, 2, 3}, workoutSetNumbers(aggregate.WorkoutExercises[0].Sets))
}

func TestWorkoutRepo_IncrementDailyLogVersion(t *testing.T) {
	pool, repo, userID := workoutRepoTestSetup(t)
	dailyLog := mustWorkoutDailyLog(t, repo, userID, "2026-06-27")

	first, err := repo.IncrementDailyLogVersion(context.Background(), userID, dailyLog.ID)
	require.NoError(t, err)
	require.NotNil(t, first)
	assert.Equal(t, int32(1), first.Version)

	second, err := repo.IncrementDailyLogVersion(context.Background(), userID, dailyLog.ID)
	require.NoError(t, err)
	require.NotNil(t, second)
	assert.Equal(t, int32(2), second.Version)

	otherUser := ensureWorkoutAtlasUser(t, pool, "wrong-user")
	wrongUser, err := repo.IncrementDailyLogVersion(context.Background(), otherUser, dailyLog.ID)
	require.NoError(t, err)
	assert.Nil(t, wrongUser)
}

func workoutRepoTestSetup(t *testing.T) (*pgxpool.Pool, atlasRepo.WorkoutRepository, string) {
	t.Helper()
	cfg := testinfra.PostgresConfig(t)
	if err := postgresrepo.RunMigrations(cfg.DSN(), zap.NewNop()); err != nil {
		if !testinfra.CoverageGateEnabled() {
			t.Skipf("postgres integration database is unavailable: %v", err)
		}
		require.NoError(t, err)
	}
	db, err := postgresrepo.New(cfg, zap.NewNop())
	if err != nil && !testinfra.CoverageGateEnabled() {
		t.Skipf("postgres integration database is unavailable: %v", err)
	}
	require.NoError(t, err)
	t.Cleanup(db.Close)

	truncateWorkoutRepoTables(t, db.Pool)
	userID := ensureWorkoutAtlasUser(t, db.Pool, "workout-user")
	return db.Pool, atlasRepo.NewWorkoutRepository(db.Pool), userID
}

func truncateWorkoutRepoTables(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `
		TRUNCATE workout_sets, workout_exercises, daily_logs, exercise_media, exercises, atlas_settings, atlas_users RESTART IDENTITY CASCADE
	`)
	require.NoError(t, err)
}

func ensureWorkoutAtlasUser(t *testing.T, pool *pgxpool.Pool, displayName string) string {
	t.Helper()
	var id string
	err := pool.QueryRow(context.Background(), `
		INSERT INTO atlas_users (display_name) VALUES ($1) RETURNING id::text
	`, displayName).Scan(&id)
	require.NoError(t, err)
	return id
}

func seedWorkoutExerciseRecord(t *testing.T, pool *pgxpool.Pool, userID string, name string, workingWeight *float64) string {
	t.Helper()
	var id string
	err := pool.QueryRow(context.Background(), `
		INSERT INTO exercises (user_id, name, muscle_groups, working_weight)
		VALUES ($1::uuid, $2, ARRAY['strength'], $3)
		RETURNING id::text
	`, userID, name, workingWeight).Scan(&id)
	require.NoError(t, err)
	return id
}

func mustWorkoutDailyLog(t *testing.T, repo atlasRepo.WorkoutRepository, userID string, rawDate string) *atlasRepo.DailyLogRecord {
	t.Helper()
	dailyLog, err := repo.GetOrCreateDailyLogByDate(context.Background(), userID, atlasModels.MustDate(rawDate))
	require.NoError(t, err)
	require.NotNil(t, dailyLog)
	return dailyLog
}

func mustAddWorkoutExercise(t *testing.T, repo atlasRepo.WorkoutRepository, pool *pgxpool.Pool, userID string, dailyLogID string, name string, position int32) *atlasRepo.WorkoutExerciseRecord {
	t.Helper()
	exerciseID := seedWorkoutExerciseRecord(t, pool, userID, name, nil)
	created, err := repo.AddWorkoutExercise(context.Background(), userID, dailyLogID, atlasRepo.AddWorkoutExerciseInput{
		ExerciseID: exerciseID,
		Position:   position,
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	return created
}

func mustAddWorkoutSet(t *testing.T, repo atlasRepo.WorkoutRepository, userID string, workoutExerciseID string, setNumber int32, weight float64, reps int32) *atlasRepo.WorkoutSetRecord {
	t.Helper()
	created, err := repo.AddWorkoutSet(context.Background(), userID, workoutExerciseID, atlasRepo.AddWorkoutSetInput{
		SetNumber: setNumber,
		Weight:    weight,
		Reps:      reps,
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	return created
}

func workoutStringPtr(value string) *string {
	return &value
}

func workoutExerciseIDs(rows []atlasRepo.WorkoutExerciseRecord) []string {
	out := make([]string, len(rows))
	for i, row := range rows {
		out[i] = row.ID
	}
	return out
}

func workoutExercisePositions(rows []atlasRepo.WorkoutExerciseRecord) []int32 {
	out := make([]int32, len(rows))
	for i, row := range rows {
		out[i] = row.Position
	}
	return out
}

func workoutSetIDs(rows []atlasRepo.WorkoutSetRecord) []string {
	out := make([]string, len(rows))
	for i, row := range rows {
		out[i] = row.ID
	}
	return out
}

func workoutSetNumbers(rows []atlasRepo.WorkoutSetRecord) []int32 {
	out := make([]int32, len(rows))
	for i, row := range rows {
		out[i] = row.SetNumber
	}
	return out
}

func requireWorkoutCheckViolation(t *testing.T, err error, constraintName string) {
	t.Helper()
	var pgErr *pgconn.PgError
	if assert.True(t, errors.As(err, &pgErr), "expected PostgreSQL error, got %v", err) {
		assert.Equal(t, "23514", pgErr.Code)
		assert.Equal(t, constraintName, pgErr.ConstraintName)
	}
}
