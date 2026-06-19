// FILE: apps/api/internal/repository/postgres/workout_migration_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify WAVE-03 workout diary migrations create the DailyLog aggregate schema contract.
//   SCOPE: daily_logs, workout_exercises, workout_sets columns, constraints, indexes, FK behavior, and excluded cardio/body-weight schema tokens.
//   DEPENDS: apps/api/internal/repository/postgres, apps/api/internal/testinfra.
//   LINKS: M-API / V-M-API / WAVE-03 / TEST-W03-001.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   TestWorkoutMigrations_FilesExistWithGraceMarkup - Proves WAVE-03 migration source files are present and governed.
//   TestWorkoutMigrations_DailyLogSchema - Proves the canonical DailyLog container schema.
//   TestWorkoutMigrations_WorkoutExerciseSchema - Proves exercise instances are ordered, user-scoped, and duplicate exercise IDs remain allowed.
//   TestWorkoutMigrations_WorkoutSetSchema - Proves ordered strength set schema and validation constraints.
//   TestWorkoutMigrations_WorkoutSetValidationRejectsInvalidBounds - Proves DB CHECK constraints reject invalid strength set values.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Tightened WAVE-03 migration verification for review feedback.
// END_CHANGE_SUMMARY

package postgres_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	postgresrepo "monorepo-template/apps/api/internal/repository/postgres"
	"monorepo-template/apps/api/internal/testinfra"
)

func TestWorkoutMigrations_FilesExistWithGraceMarkup(t *testing.T) {
	expected := map[string][]string{
		"migrations/00083_daily_logs.sql": {
			"FILE: apps/api/internal/repository/postgres/migrations/00083_daily_logs.sql",
			"START_MODULE_CONTRACT",
			"daily_logs - Canonical daily aggregate container",
			"-- +goose Up",
			"CREATE TABLE daily_logs",
			"CONSTRAINT uq_daily_logs_user_date UNIQUE (user_id, date)",
			"CONSTRAINT chk_daily_logs_version CHECK (version >= 0)",
			"CREATE INDEX idx_daily_logs_user_date ON daily_logs (user_id, date)",
			"CREATE INDEX idx_daily_logs_user_date_desc ON daily_logs (user_id, date DESC)",
			"-- +goose Down",
		},
		"migrations/00084_workout_exercises.sql": {
			"FILE: apps/api/internal/repository/postgres/migrations/00084_workout_exercises.sql",
			"START_MODULE_CONTRACT",
			"workout_exercises - Ordered strength exercise instances",
			"-- +goose Up",
			"CREATE TABLE workout_exercises",
			"daily_log_id            UUID NOT NULL REFERENCES daily_logs(id) ON DELETE CASCADE",
			"exercise_id             UUID NOT NULL REFERENCES exercises(id) ON DELETE RESTRICT",
			"CONSTRAINT chk_workout_exercises_position CHECK (position > 0)",
			"CONSTRAINT chk_workout_exercises_working_weight_snapshot CHECK (working_weight_snapshot IS NULL OR working_weight_snapshot > 0)",
			"CONSTRAINT uq_workout_exercises_daily_log_position UNIQUE (daily_log_id, position)",
			"-- +goose Down",
		},
		"migrations/00085_workout_sets.sql": {
			"FILE: apps/api/internal/repository/postgres/migrations/00085_workout_sets.sql",
			"START_MODULE_CONTRACT",
			"workout_sets - Ordered strength set rows",
			"-- +goose Up",
			"CREATE TABLE workout_sets",
			"workout_exercise_id UUID NOT NULL REFERENCES workout_exercises(id) ON DELETE CASCADE",
			"CONSTRAINT chk_workout_sets_set_number CHECK (set_number > 0)",
			"CONSTRAINT chk_workout_sets_weight CHECK (weight > 0)",
			"CONSTRAINT chk_workout_sets_reps CHECK (reps > 0)",
			"CONSTRAINT chk_workout_sets_rpe CHECK (rpe IS NULL OR (rpe >= 1 AND rpe <= 10))",
			"CONSTRAINT chk_workout_sets_rir CHECK (rir IS NULL OR (rir >= 0 AND rir <= 10))",
			"CONSTRAINT uq_workout_sets_exercise_set_number UNIQUE (workout_exercise_id, set_number)",
			"-- +goose Down",
		},
	}
	prohibited := []string{
		"later WAVE-04",
		"cardio_entries",
		"CardioType",
		"HeartRateZone",
		"body_weight",
		"bodyWeight",
	}

	for path, snippets := range expected {
		t.Run(path, func(t *testing.T) {
			content, err := os.ReadFile(path)
			require.NoError(t, err)
			raw := string(content)
			for _, snippet := range snippets {
				assert.Contains(t, raw, snippet)
			}
			for _, token := range prohibited {
				assert.NotContains(t, raw, token)
			}
		})
	}
}

func TestWorkoutMigrations_DailyLogSchema(t *testing.T) {
	pool := workoutMigrationTestPool(t)

	requireTable(t, pool, "daily_logs")
	requireColumn(t, pool, "daily_logs", "id", "uuid", "NO", "gen_random_uuid")
	requireColumn(t, pool, "daily_logs", "user_id", "uuid", "NO", "")
	requireColumn(t, pool, "daily_logs", "date", "date", "NO", "")
	requireColumn(t, pool, "daily_logs", "notes", "text", "YES", "")
	requireColumn(t, pool, "daily_logs", "version", "integer", "NO", "0")
	requireColumn(t, pool, "daily_logs", "created_at", "timestamp with time zone", "NO", "now()")
	requireColumn(t, pool, "daily_logs", "updated_at", "timestamp with time zone", "NO", "now()")
	requireNoColumn(t, pool, "daily_logs", "body_weight")
	requireNoColumn(t, pool, "daily_logs", "bodyWeight")
	requireNoTable(t, pool, "cardio_entries")

	requireConstraint(t, pool, "daily_logs", "uq_daily_logs_user_date", "UNIQUE")
	requireUniqueConstraintColumns(t, pool, "daily_logs", "uq_daily_logs_user_date", []string{"user_id", "date"})
	requireConstraint(t, pool, "daily_logs", "chk_daily_logs_version", "CHECK")
	requireCheckConstraintDefinition(t, pool, "daily_logs", "chk_daily_logs_version", []string{"version >= 0"})
	requireForeignKeyDeleteRule(t, pool, "daily_logs", "user_id", "atlas_users", "NO ACTION")
	requireIndex(t, pool, "idx_daily_logs_user_date")
	requireIndex(t, pool, "idx_daily_logs_user_date_desc")
	requireIndexDefinitionContains(t, pool, "idx_daily_logs_user_date_desc", "date DESC")
}

func TestWorkoutMigrations_WorkoutExerciseSchema(t *testing.T) {
	pool := workoutMigrationTestPool(t)

	requireTable(t, pool, "workout_exercises")
	requireColumn(t, pool, "workout_exercises", "id", "uuid", "NO", "gen_random_uuid")
	requireColumn(t, pool, "workout_exercises", "user_id", "uuid", "NO", "")
	requireColumn(t, pool, "workout_exercises", "daily_log_id", "uuid", "NO", "")
	requireColumn(t, pool, "workout_exercises", "exercise_id", "uuid", "NO", "")
	requireColumn(t, pool, "workout_exercises", "position", "integer", "NO", "")
	requireColumn(t, pool, "workout_exercises", "working_weight_snapshot", "real", "YES", "")
	requireColumn(t, pool, "workout_exercises", "notes", "text", "YES", "")
	requireColumn(t, pool, "workout_exercises", "created_at", "timestamp with time zone", "NO", "now()")
	requireColumn(t, pool, "workout_exercises", "updated_at", "timestamp with time zone", "NO", "now()")
	requireNoColumn(t, pool, "workout_exercises", "body_weight")

	requireConstraint(t, pool, "workout_exercises", "chk_workout_exercises_position", "CHECK")
	requireCheckConstraintDefinition(t, pool, "workout_exercises", "chk_workout_exercises_position", []string{"position > 0"})
	requireConstraint(t, pool, "workout_exercises", "chk_workout_exercises_working_weight_snapshot", "CHECK")
	requireConstraint(t, pool, "workout_exercises", "uq_workout_exercises_daily_log_position", "UNIQUE")
	requireUniqueConstraintColumns(t, pool, "workout_exercises", "uq_workout_exercises_daily_log_position", []string{"daily_log_id", "position"})
	requireNoUniqueConstraintOnColumns(t, pool, "workout_exercises", []string{"daily_log_id", "exercise_id"})
	requireForeignKeyDeleteRule(t, pool, "workout_exercises", "daily_log_id", "daily_logs", "CASCADE")
	requireForeignKeyDeleteRule(t, pool, "workout_exercises", "exercise_id", "exercises", "RESTRICT")
	requireIndex(t, pool, "idx_workout_exercises_user_daily_log")
	requireIndex(t, pool, "idx_workout_exercises_exercise")
}

func TestWorkoutMigrations_WorkoutSetSchema(t *testing.T) {
	pool := workoutMigrationTestPool(t)

	requireTable(t, pool, "workout_sets")
	requireColumn(t, pool, "workout_sets", "id", "uuid", "NO", "gen_random_uuid")
	requireColumn(t, pool, "workout_sets", "workout_exercise_id", "uuid", "NO", "")
	requireColumn(t, pool, "workout_sets", "set_number", "integer", "NO", "")
	requireColumn(t, pool, "workout_sets", "weight", "real", "NO", "")
	requireColumn(t, pool, "workout_sets", "reps", "integer", "NO", "")
	requireColumn(t, pool, "workout_sets", "rpe", "real", "YES", "")
	requireColumn(t, pool, "workout_sets", "rir", "integer", "YES", "")
	requireColumn(t, pool, "workout_sets", "notes", "text", "YES", "")
	requireColumn(t, pool, "workout_sets", "created_at", "timestamp with time zone", "NO", "now()")
	requireColumn(t, pool, "workout_sets", "updated_at", "timestamp with time zone", "NO", "now()")

	requireConstraint(t, pool, "workout_sets", "chk_workout_sets_set_number", "CHECK")
	requireCheckConstraintDefinition(t, pool, "workout_sets", "chk_workout_sets_set_number", []string{"set_number > 0"})
	requireConstraint(t, pool, "workout_sets", "chk_workout_sets_weight", "CHECK")
	requireConstraint(t, pool, "workout_sets", "chk_workout_sets_reps", "CHECK")
	requireCheckConstraintDefinition(t, pool, "workout_sets", "chk_workout_sets_reps", []string{"reps > 0"})
	requireConstraint(t, pool, "workout_sets", "chk_workout_sets_rpe", "CHECK")
	requireConstraint(t, pool, "workout_sets", "chk_workout_sets_rir", "CHECK")
	requireCheckConstraintDefinition(t, pool, "workout_sets", "chk_workout_sets_rir", []string{"rir IS NULL", "rir >= 0", "rir <= 10"})
	requireConstraint(t, pool, "workout_sets", "uq_workout_sets_exercise_set_number", "UNIQUE")
	requireUniqueConstraintColumns(t, pool, "workout_sets", "uq_workout_sets_exercise_set_number", []string{"workout_exercise_id", "set_number"})
	requireForeignKeyDeleteRule(t, pool, "workout_sets", "workout_exercise_id", "workout_exercises", "CASCADE")
	requireIndex(t, pool, "idx_workout_sets_workout_exercise")
}

func TestWorkoutMigrations_WorkoutSetValidationRejectsInvalidBounds(t *testing.T) {
	pool := workoutMigrationTestPool(t)
	workoutExerciseID := seedWorkoutExercise(t, pool)

	_, err := pool.Exec(context.Background(), `
		INSERT INTO workout_sets (workout_exercise_id, set_number, weight, reps, rpe, rir)
		VALUES ($1::uuid, 1, 1, 1, 1, 0)
	`, workoutExerciseID)
	require.NoError(t, err)

	_, err = pool.Exec(context.Background(), `
		INSERT INTO workout_sets (workout_exercise_id, set_number, weight, reps, rpe, rir)
		VALUES ($1::uuid, 2, 1, 1, 10, 10)
	`, workoutExerciseID)
	require.NoError(t, err)

	cases := []struct {
		name       string
		setNumber  int
		weight     float64
		reps       int
		rpe        any
		rir        any
		constraint string
	}{
		{
			name:       "zero weight",
			setNumber:  3,
			weight:     0,
			reps:       1,
			rpe:        nil,
			rir:        nil,
			constraint: "chk_workout_sets_weight",
		},
		{
			name:       "negative weight",
			setNumber:  4,
			weight:     -1,
			reps:       1,
			rpe:        nil,
			rir:        nil,
			constraint: "chk_workout_sets_weight",
		},
		{
			name:       "zero rpe",
			setNumber:  5,
			weight:     1,
			reps:       1,
			rpe:        0,
			rir:        nil,
			constraint: "chk_workout_sets_rpe",
		},
		{
			name:       "rpe above ten",
			setNumber:  6,
			weight:     1,
			reps:       1,
			rpe:        11,
			rir:        nil,
			constraint: "chk_workout_sets_rpe",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := pool.Exec(context.Background(), `
				INSERT INTO workout_sets (workout_exercise_id, set_number, weight, reps, rpe, rir)
				VALUES ($1::uuid, $2, $3, $4, $5, $6)
			`, workoutExerciseID, tc.setNumber, tc.weight, tc.reps, tc.rpe, tc.rir)
			requireCheckViolation(t, err, tc.constraint)
		})
	}
}

func workoutMigrationTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := testinfra.PostgresDSN()
	testinfra.RequireSafePostgresDSN(t, dsn)
	if err := postgresrepo.RunMigrations(dsn, zap.NewNop()); err != nil {
		if !testinfra.CoverageGateEnabled() {
			t.Skipf("postgres integration database is unavailable: %v", err)
		}
		require.NoError(t, err)
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	t.Cleanup(pool.Close)
	return pool
}

func seedWorkoutExercise(t *testing.T, pool *pgxpool.Pool) string {
	t.Helper()
	var workoutExerciseID string
	err := pool.QueryRow(context.Background(), `
		WITH user_row AS (
			INSERT INTO atlas_users DEFAULT VALUES
			RETURNING id
		),
		daily_log_row AS (
			INSERT INTO daily_logs (user_id, date)
			SELECT id, CURRENT_DATE FROM user_row
			RETURNING id, user_id
		),
		exercise_row AS (
			INSERT INTO exercises (user_id, name)
			SELECT id, 'WAVE-03 migration constraint test exercise' FROM user_row
			RETURNING id
		)
		INSERT INTO workout_exercises (user_id, daily_log_id, exercise_id, position)
		SELECT daily_log_row.user_id, daily_log_row.id, exercise_row.id, 1
		FROM daily_log_row, exercise_row
		RETURNING id::text
	`).Scan(&workoutExerciseID)
	require.NoError(t, err)
	return workoutExerciseID
}

func requireTable(t *testing.T, pool *pgxpool.Pool, table string) {
	t.Helper()
	var exists bool
	err := pool.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = $1
		)
	`, table).Scan(&exists)
	require.NoError(t, err)
	require.True(t, exists, "expected table %s to exist", table)
}

func requireColumn(t *testing.T, pool *pgxpool.Pool, table, column, dataType, nullable, defaultContains string) {
	t.Helper()
	var actualType, actualNullable string
	var columnDefault *string
	err := pool.QueryRow(context.Background(), `
		SELECT data_type, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = $1 AND column_name = $2
	`, table, column).Scan(&actualType, &actualNullable, &columnDefault)
	require.NoError(t, err, "expected column %s.%s", table, column)
	assert.Equal(t, dataType, actualType, "data type for %s.%s", table, column)
	assert.Equal(t, nullable, actualNullable, "nullable for %s.%s", table, column)
	if defaultContains != "" {
		if assert.NotNil(t, columnDefault, "default for %s.%s", table, column) {
			assert.Contains(t, *columnDefault, defaultContains, "default for %s.%s", table, column)
		}
	}
}

func requireNoColumn(t *testing.T, pool *pgxpool.Pool, table, column string) {
	t.Helper()
	var exists bool
	err := pool.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.columns
			WHERE table_schema = 'public' AND table_name = $1 AND column_name = $2
		)
	`, table, column).Scan(&exists)
	require.NoError(t, err)
	assert.False(t, exists, "did not expect column %s.%s", table, column)
}

func requireNoTable(t *testing.T, pool *pgxpool.Pool, table string) {
	t.Helper()
	var exists bool
	err := pool.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = $1
		)
	`, table).Scan(&exists)
	require.NoError(t, err)
	assert.False(t, exists, "did not expect table %s to exist", table)
}

func requireConstraint(t *testing.T, pool *pgxpool.Pool, table, constraintName, constraintType string) {
	t.Helper()
	var actualType string
	err := pool.QueryRow(context.Background(), `
		SELECT constraint_type
		FROM information_schema.table_constraints
		WHERE table_schema = 'public' AND table_name = $1 AND constraint_name = $2
	`, table, constraintName).Scan(&actualType)
	require.NoError(t, err, "expected constraint %s on %s", constraintName, table)
	assert.Equal(t, constraintType, actualType)
}

func requireUniqueConstraintColumns(t *testing.T, pool *pgxpool.Pool, table, constraintName string, columns []string) {
	t.Helper()
	var constrainedColumns string
	err := pool.QueryRow(context.Background(), `
		SELECT string_agg(kcu.column_name, ',' ORDER BY kcu.ordinal_position) AS constrained_columns
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu
			ON kcu.constraint_catalog = tc.constraint_catalog
			AND kcu.constraint_schema = tc.constraint_schema
			AND kcu.constraint_name = tc.constraint_name
		WHERE tc.table_schema = 'public'
			AND tc.table_name = $1
			AND tc.constraint_name = $2
			AND tc.constraint_type = 'UNIQUE'
		GROUP BY tc.constraint_name
	`, table, constraintName).Scan(&constrainedColumns)
	require.NoError(t, err, "expected unique constraint %s on %s", constraintName, table)
	assert.Equal(t, strings.Join(columns, ","), constrainedColumns)
}

func requireCheckConstraintDefinition(t *testing.T, pool *pgxpool.Pool, table, constraintName string, snippets []string) {
	t.Helper()
	var definition string
	err := pool.QueryRow(context.Background(), `
		SELECT pg_get_constraintdef(c.oid)
		FROM pg_constraint c
		JOIN pg_class t ON t.oid = c.conrelid
		JOIN pg_namespace n ON n.oid = t.relnamespace
		WHERE n.nspname = 'public'
			AND t.relname = $1
			AND c.conname = $2
	`, table, constraintName).Scan(&definition)
	require.NoError(t, err, "expected check constraint %s on %s", constraintName, table)
	lowerDefinition := strings.ToLower(strings.ReplaceAll(definition, `"`, ""))
	for _, snippet := range snippets {
		assert.Contains(t, lowerDefinition, strings.ToLower(snippet), "constraint %s definition", constraintName)
	}
}

func requireForeignKeyDeleteRule(t *testing.T, pool *pgxpool.Pool, table, column, referencedTable, deleteRule string) {
	t.Helper()
	var actualRule string
	err := pool.QueryRow(context.Background(), `
		SELECT rc.delete_rule
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu
			ON kcu.constraint_catalog = tc.constraint_catalog
			AND kcu.constraint_schema = tc.constraint_schema
			AND kcu.constraint_name = tc.constraint_name
		JOIN information_schema.referential_constraints rc
			ON rc.constraint_catalog = tc.constraint_catalog
			AND rc.constraint_schema = tc.constraint_schema
			AND rc.constraint_name = tc.constraint_name
		JOIN information_schema.constraint_column_usage ccu
			ON ccu.constraint_catalog = rc.unique_constraint_catalog
			AND ccu.constraint_schema = rc.unique_constraint_schema
			AND ccu.constraint_name = rc.unique_constraint_name
		WHERE tc.table_schema = 'public'
			AND tc.table_name = $1
			AND tc.constraint_type = 'FOREIGN KEY'
			AND kcu.column_name = $2
			AND ccu.table_name = $3
	`, table, column, referencedTable).Scan(&actualRule)
	require.NoError(t, err, "expected FK from %s.%s to %s", table, column, referencedTable)
	assert.Equal(t, deleteRule, actualRule)
}

func requireCheckViolation(t *testing.T, err error, constraintName string) {
	t.Helper()
	var pgErr *pgconn.PgError
	if assert.ErrorAs(t, err, &pgErr) {
		assert.Equal(t, "23514", pgErr.Code)
		assert.Equal(t, constraintName, pgErr.ConstraintName)
	}
}

func requireIndex(t *testing.T, pool *pgxpool.Pool, indexName string) {
	t.Helper()
	var exists bool
	err := pool.QueryRow(context.Background(), `SELECT to_regclass($1) IS NOT NULL`, "public."+indexName).Scan(&exists)
	require.NoError(t, err)
	require.True(t, exists, "expected index %s", indexName)
}

func requireIndexDefinitionContains(t *testing.T, pool *pgxpool.Pool, indexName, snippet string) {
	t.Helper()
	var definition string
	err := pool.QueryRow(context.Background(), `SELECT pg_get_indexdef(to_regclass($1))`, "public."+indexName).Scan(&definition)
	require.NoError(t, err, "expected index %s", indexName)
	assert.Contains(t, definition, snippet)
}

func requireNoUniqueConstraintOnColumns(t *testing.T, pool *pgxpool.Pool, table string, columns []string) {
	t.Helper()
	rows, err := pool.Query(context.Background(), `
		SELECT tc.constraint_name, string_agg(kcu.column_name, ',' ORDER BY kcu.ordinal_position) AS constrained_columns
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu
			ON kcu.constraint_catalog = tc.constraint_catalog
			AND kcu.constraint_schema = tc.constraint_schema
			AND kcu.constraint_name = tc.constraint_name
		WHERE tc.table_schema = 'public'
			AND tc.table_name = $1
			AND tc.constraint_type = 'UNIQUE'
		GROUP BY tc.constraint_name
	`, table)
	require.NoError(t, err)
	defer rows.Close()

	unwanted := strings.Join(columns, ",")
	for rows.Next() {
		var name, constrainedColumns string
		require.NoError(t, rows.Scan(&name, &constrainedColumns))
		assert.NotEqual(t, unwanted, constrainedColumns, "unexpected unique constraint %s on %s(%s)", name, table, unwanted)
	}
	require.NoError(t, rows.Err())
}
