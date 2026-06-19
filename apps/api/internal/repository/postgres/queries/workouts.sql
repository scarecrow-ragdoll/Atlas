-- name: GetDailyLogByID :one
SELECT id, user_id, date, notes, version, created_at, updated_at
FROM daily_logs
WHERE user_id = sqlc.arg('user_id') AND id = sqlc.arg('id')
LIMIT 1;

-- name: GetDailyLogByDate :one
SELECT id, user_id, date, notes, version, created_at, updated_at
FROM daily_logs
WHERE user_id = sqlc.arg('user_id') AND date = sqlc.arg('date')
LIMIT 1;

-- name: CreateDailyLog :one
INSERT INTO daily_logs (user_id, date, notes)
VALUES (sqlc.arg('user_id'), sqlc.arg('date'), sqlc.narg('notes'))
ON CONFLICT (user_id, date) DO UPDATE SET updated_at = daily_logs.updated_at
RETURNING id, user_id, date, notes, version, created_at, updated_at;

-- name: LockDailyLogByID :one
SELECT id, user_id, date, notes, version, created_at, updated_at
FROM daily_logs
WHERE user_id = sqlc.arg('user_id') AND id = sqlc.arg('id')
FOR UPDATE;

-- name: LockDailyLogByDate :one
SELECT id, user_id, date, notes, version, created_at, updated_at
FROM daily_logs
WHERE user_id = sqlc.arg('user_id') AND date = sqlc.arg('date')
FOR UPDATE;

-- name: LockDailyLogByWorkoutExerciseID :one
SELECT dl.id, dl.user_id, dl.date, dl.notes, dl.version, dl.created_at, dl.updated_at
FROM daily_logs dl
JOIN workout_exercises we ON we.daily_log_id = dl.id
WHERE dl.user_id = sqlc.arg('user_id') AND we.id = sqlc.arg('workout_exercise_id')
FOR UPDATE OF dl;

-- name: LockDailyLogByWorkoutSetID :one
SELECT dl.id, dl.user_id, dl.date, dl.notes, dl.version, dl.created_at, dl.updated_at
FROM daily_logs dl
JOIN workout_exercises we ON we.daily_log_id = dl.id
JOIN workout_sets ws ON ws.workout_exercise_id = we.id
WHERE dl.user_id = sqlc.arg('user_id') AND ws.id = sqlc.arg('workout_set_id')
FOR UPDATE OF dl;

-- name: IncrementDailyLogVersion :one
UPDATE daily_logs
SET version = version + 1, updated_at = now()
WHERE user_id = sqlc.arg('user_id') AND id = sqlc.arg('id')
RETURNING id, user_id, date, notes, version, created_at, updated_at;

-- name: UpdateDailyLogNotes :one
UPDATE daily_logs
SET notes = sqlc.narg('notes'), updated_at = now()
WHERE user_id = sqlc.arg('user_id') AND id = sqlc.arg('id')
RETURNING id, user_id, date, notes, version, created_at, updated_at;

-- name: ListDailyLogSummaries :many
SELECT
    dl.id,
    dl.date,
    dl.version,
    COUNT(DISTINCT we.id)::int AS workout_exercise_count,
    COUNT(ws.id)::int AS workout_set_count,
    COALESCE(SUM(ws.weight * ws.reps), 0)::float8 AS total_volume,
    dl.updated_at
FROM daily_logs dl
LEFT JOIN workout_exercises we ON we.daily_log_id = dl.id
LEFT JOIN workout_sets ws ON ws.workout_exercise_id = we.id
WHERE dl.user_id = sqlc.arg('user_id')
  AND dl.date >= sqlc.arg('from_date')
  AND dl.date <= sqlc.arg('to_date')
GROUP BY dl.id
ORDER BY dl.date ASC;

-- name: ListWorkoutExercisesByDailyLog :many
SELECT id, user_id, daily_log_id, exercise_id, position, working_weight_snapshot, notes, created_at, updated_at
FROM workout_exercises
WHERE user_id = sqlc.arg('user_id') AND daily_log_id = sqlc.arg('daily_log_id')
ORDER BY position ASC;

-- name: CreateWorkoutExercise :one
INSERT INTO workout_exercises (
    user_id,
    daily_log_id,
    exercise_id,
    position,
    working_weight_snapshot,
    notes
)
VALUES (
    sqlc.arg('user_id'),
    sqlc.arg('daily_log_id'),
    sqlc.arg('exercise_id'),
    sqlc.arg('position'),
    sqlc.narg('working_weight_snapshot'),
    sqlc.narg('notes')
)
RETURNING id, user_id, daily_log_id, exercise_id, position, working_weight_snapshot, notes, created_at, updated_at;

-- name: UpdateWorkoutExercise :one
UPDATE workout_exercises
SET position = COALESCE(sqlc.narg('position'), position),
    notes = CASE WHEN sqlc.arg('set_notes')::boolean THEN sqlc.narg('notes') ELSE notes END,
    updated_at = now()
WHERE user_id = sqlc.arg('user_id') AND id = sqlc.arg('id')
RETURNING id, user_id, daily_log_id, exercise_id, position, working_weight_snapshot, notes, created_at, updated_at;

-- name: DeleteWorkoutExercise :one
DELETE FROM workout_exercises
WHERE user_id = sqlc.arg('user_id') AND id = sqlc.arg('id')
RETURNING id, user_id, daily_log_id, exercise_id, position, working_weight_snapshot, notes, created_at, updated_at;

-- name: TempShiftWorkoutExercisePositionsForInsert :many
UPDATE workout_exercises
SET position = position + 1000000,
    updated_at = now()
WHERE workout_exercises.user_id = sqlc.arg('user_id')
  AND workout_exercises.daily_log_id = sqlc.arg('daily_log_id')
  AND workout_exercises.position >= sqlc.arg('position')
RETURNING id, user_id, daily_log_id, exercise_id, position, working_weight_snapshot, notes, created_at, updated_at;

-- name: NormalizeWorkoutExercisePositionsForInsert :many
UPDATE workout_exercises
SET position = position - 999999,
    updated_at = now()
WHERE workout_exercises.user_id = sqlc.arg('user_id')
  AND workout_exercises.daily_log_id = sqlc.arg('daily_log_id')
  AND workout_exercises.position >= sqlc.arg('position')::integer + 1000000
RETURNING id, user_id, daily_log_id, exercise_id, position, working_weight_snapshot, notes, created_at, updated_at;

-- name: TempShiftWorkoutExercisePositionsAfterDelete :many
UPDATE workout_exercises
SET position = position + 1000000,
    updated_at = now()
WHERE workout_exercises.user_id = sqlc.arg('user_id')
  AND workout_exercises.daily_log_id = sqlc.arg('daily_log_id')
  AND workout_exercises.position > sqlc.arg('deleted_position')
RETURNING id, user_id, daily_log_id, exercise_id, position, working_weight_snapshot, notes, created_at, updated_at;

-- name: NormalizeWorkoutExercisePositionsAfterDelete :many
UPDATE workout_exercises
SET position = position - 1000001,
    updated_at = now()
WHERE workout_exercises.user_id = sqlc.arg('user_id')
  AND workout_exercises.daily_log_id = sqlc.arg('daily_log_id')
  AND workout_exercises.position > sqlc.arg('deleted_position')::integer + 1000000
RETURNING id, user_id, daily_log_id, exercise_id, position, working_weight_snapshot, notes, created_at, updated_at;

-- name: SetWorkoutExercisePosition :one
UPDATE workout_exercises
SET position = sqlc.arg('position'),
    updated_at = now()
WHERE user_id = sqlc.arg('user_id')
  AND daily_log_id = sqlc.arg('daily_log_id')
  AND id = sqlc.arg('id')
RETURNING id, user_id, daily_log_id, exercise_id, position, working_weight_snapshot, notes, created_at, updated_at;

-- name: ListWorkoutSetsByExerciseIDs :many
SELECT id, workout_exercise_id, set_number, weight, reps, rpe, rir, notes, created_at, updated_at
FROM workout_sets
WHERE workout_exercise_id = ANY(sqlc.arg('workout_exercise_ids')::uuid[])
ORDER BY workout_exercise_id ASC, set_number ASC;

-- name: CreateWorkoutSet :one
INSERT INTO workout_sets (
    workout_exercise_id,
    set_number,
    weight,
    reps,
    rpe,
    rir,
    notes
)
VALUES (
    sqlc.arg('workout_exercise_id'),
    sqlc.arg('set_number'),
    sqlc.arg('weight'),
    sqlc.arg('reps'),
    sqlc.narg('rpe'),
    sqlc.narg('rir'),
    sqlc.narg('notes')
)
RETURNING id, workout_exercise_id, set_number, weight, reps, rpe, rir, notes, created_at, updated_at;

-- name: UpdateWorkoutSet :one
UPDATE workout_sets
SET set_number = COALESCE(sqlc.narg('set_number'), set_number),
    weight = COALESCE(sqlc.narg('weight'), weight),
    reps = COALESCE(sqlc.narg('reps'), reps),
    rpe = CASE WHEN sqlc.arg('set_rpe')::boolean THEN sqlc.narg('rpe') ELSE rpe END,
    rir = CASE WHEN sqlc.arg('set_rir')::boolean THEN sqlc.narg('rir') ELSE rir END,
    notes = CASE WHEN sqlc.arg('set_notes')::boolean THEN sqlc.narg('notes') ELSE notes END,
    updated_at = now()
WHERE workout_exercise_id = sqlc.arg('workout_exercise_id') AND id = sqlc.arg('id')
RETURNING id, workout_exercise_id, set_number, weight, reps, rpe, rir, notes, created_at, updated_at;

-- name: DeleteWorkoutSet :one
DELETE FROM workout_sets
WHERE workout_exercise_id = sqlc.arg('workout_exercise_id') AND id = sqlc.arg('id')
RETURNING id, workout_exercise_id, set_number, weight, reps, rpe, rir, notes, created_at, updated_at;

-- name: TempShiftWorkoutSetNumbersForInsert :many
UPDATE workout_sets
SET set_number = set_number + 1000000,
    updated_at = now()
WHERE workout_sets.workout_exercise_id = sqlc.arg('workout_exercise_id')
  AND workout_sets.set_number >= sqlc.arg('set_number')
RETURNING id, workout_exercise_id, set_number, weight, reps, rpe, rir, notes, created_at, updated_at;

-- name: NormalizeWorkoutSetNumbersForInsert :many
UPDATE workout_sets
SET set_number = set_number - 999999,
    updated_at = now()
WHERE workout_sets.workout_exercise_id = sqlc.arg('workout_exercise_id')
  AND workout_sets.set_number >= sqlc.arg('set_number')::integer + 1000000
RETURNING id, workout_exercise_id, set_number, weight, reps, rpe, rir, notes, created_at, updated_at;

-- name: TempShiftWorkoutSetNumbersAfterDelete :many
UPDATE workout_sets
SET set_number = set_number + 1000000,
    updated_at = now()
WHERE workout_sets.workout_exercise_id = sqlc.arg('workout_exercise_id')
  AND workout_sets.set_number > sqlc.arg('deleted_set_number')
RETURNING id, workout_exercise_id, set_number, weight, reps, rpe, rir, notes, created_at, updated_at;

-- name: NormalizeWorkoutSetNumbersAfterDelete :many
UPDATE workout_sets
SET set_number = set_number - 1000001,
    updated_at = now()
WHERE workout_sets.workout_exercise_id = sqlc.arg('workout_exercise_id')
  AND workout_sets.set_number > sqlc.arg('deleted_set_number')::integer + 1000000
RETURNING id, workout_exercise_id, set_number, weight, reps, rpe, rir, notes, created_at, updated_at;

-- name: SetWorkoutSetNumber :one
UPDATE workout_sets
SET set_number = sqlc.arg('set_number'),
    updated_at = now()
WHERE workout_exercise_id = sqlc.arg('workout_exercise_id') AND id = sqlc.arg('id')
RETURNING id, workout_exercise_id, set_number, weight, reps, rpe, rir, notes, created_at, updated_at;

-- FILE: apps/api/internal/repository/postgres/queries/workouts.sql
-- VERSION: 1.0.0
-- START_MODULE_CONTRACT
--   PURPOSE: Define sqlc queries for DailyLog, workout_exercises, and workout_sets.
--   SCOPE: Aggregate reads, date range summaries, row locks, version increments, CRUD, and reorder helpers for WAVE-03.
--   DEPENDS: daily_logs, workout_exercises, workout_sets, exercises.
--   LINKS: M-API / V-M-API / WAVE-03.
--   ROLE: CONFIG
--   MAP_MODE: SUMMARY
-- END_MODULE_CONTRACT
-- START_MODULE_MAP
--   GetDailyLogByID - Fetches one DailyLog row by user and id.
--   GetDailyLogByDate - Fetches one DailyLog row by user and date.
--   CreateDailyLog - Inserts or returns the user-date DailyLog row.
--   LockDailyLogByID - Locks one DailyLog aggregate row by id.
--   LockDailyLogByDate - Locks one DailyLog aggregate row by date.
--   LockDailyLogByWorkoutExerciseID - Locks the DailyLog that owns a workout exercise.
--   LockDailyLogByWorkoutSetID - Locks the DailyLog that owns a workout set.
--   IncrementDailyLogVersion - Bumps the optimistic aggregate version.
--   UpdateDailyLogNotes - Saves DailyLog notes and returns the persisted row.
--   ListDailyLogSummaries - Lists date-range DailyLog summary metrics.
--   ListWorkoutExercisesByDailyLog - Lists workout exercises for one DailyLog in position order.
--   CreateWorkoutExercise - Inserts a workout exercise and returns the persisted row.
--   UpdateWorkoutExercise - Updates workout exercise notes or position and returns the persisted row.
--   DeleteWorkoutExercise - Deletes a workout exercise and returns the deleted row.
--   TempShiftWorkoutExercisePositionsForInsert - Moves insert-affected exercise positions to a temporary range.
--   NormalizeWorkoutExercisePositionsForInsert - Maps temporary insert positions back to final position + 1 values.
--   TempShiftWorkoutExercisePositionsAfterDelete - Moves delete-affected exercise positions to a temporary range.
--   NormalizeWorkoutExercisePositionsAfterDelete - Maps temporary delete positions back to final position - 1 values.
--   SetWorkoutExercisePosition - Assigns a final exercise position during reorder.
--   ListWorkoutSetsByExerciseIDs - Lists sets for workout exercises in set-number order.
--   CreateWorkoutSet - Inserts a workout set and returns the persisted row.
--   UpdateWorkoutSet - Updates workout set fields and returns the persisted row.
--   DeleteWorkoutSet - Deletes a workout set and returns the deleted row.
--   TempShiftWorkoutSetNumbersForInsert - Moves insert-affected set numbers to a temporary range.
--   NormalizeWorkoutSetNumbersForInsert - Maps temporary insert set numbers back to final set_number + 1 values.
--   TempShiftWorkoutSetNumbersAfterDelete - Moves delete-affected set numbers to a temporary range.
--   NormalizeWorkoutSetNumbersAfterDelete - Maps temporary delete set numbers back to final set_number - 1 values.
--   SetWorkoutSetNumber - Assigns a final set number during reorder.
-- END_MODULE_MAP
-- START_CHANGE_SUMMARY
--   LAST_CHANGE: 1.0.1 - Split reordering helpers into sequential temp and normalize queries, and made nullable update fields explicit.
-- END_CHANGE_SUMMARY
