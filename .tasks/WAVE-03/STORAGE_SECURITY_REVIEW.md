<!-- FILE: .tasks/WAVE-03/STORAGE_SECURITY_REVIEW.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Record Atlas-qb2.3.2 backend storage and security review for WAVE-03 Workout Diary. -->
<!--   SCOPE: Reviews migrations, sqlc queries, repository transactions, user scoping, FK/cascade/restrict behavior, optimistic version locks, error handling, and logging privacy; excludes GraphQL/generated deep review, GRACE consistency review, and final readiness packet work owned by later QA beads. -->
<!--   DEPENDS: apps/api/internal/repository/postgres/migrations/00083_daily_logs.sql, apps/api/internal/repository/postgres/migrations/00084_workout_exercises.sql, apps/api/internal/repository/postgres/migrations/00085_workout_sets.sql, apps/api/internal/repository/postgres/queries/workouts.sql, apps/api/internal/atlas/repository/postgres/workout_repo.go, apps/api/internal/atlas/service/workout.go. -->
<!--   LINKS: M-API / V-M-API / WAVE-03 / Atlas-qb2.3.2. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Review Verdict - States storage and security readiness for the bead. -->
<!--   Static Review - Records findings for schema, queries, transactions, versioning, and privacy. -->
<!--   Verification Evidence - Records focused commands rerun for this review. -->
<!--   Findings And Handoff - Records severity-classified issues and next QA boundaries. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.0 - Added W03 backend storage and security review for Atlas-qb2.3.2. -->
<!-- END_CHANGE_SUMMARY -->

# W03 Storage And Security Review

## Review Verdict

Status: PASS for `Atlas-qb2.3.2`.

No Critical or Important storage/security findings were found.

Product/code blockers: none.

Follow-up blockers created by this review: none.

## Static Review

| Area | Evidence | Verdict |
| --- | --- | --- |
| `daily_logs` canonical container constraints | `00083_daily_logs.sql` creates `daily_logs` with `id`, `user_id`, `date`, nullable `notes`, `version`, timestamps, `UNIQUE(user_id, date)`, `CHECK(version >= 0)`, and user/date indexes. | PASS. |
| `workout_exercises` constraints and ownership shape | `00084_workout_exercises.sql` creates ordered exercise instances with `user_id`, `daily_log_id` `ON DELETE CASCADE`, `exercise_id` `ON DELETE RESTRICT`, positive position, optional positive working-weight snapshot, unique `(daily_log_id, position)`, and user/daily-log indexes. Duplicate `exercise_id` values per day remain allowed. | PASS. |
| `workout_sets` constraints | `00085_workout_sets.sql` creates ordered set rows with `workout_exercise_id ON DELETE CASCADE`, positive set number/weight/reps, optional bounded RPE/RIR, and unique `(workout_exercise_id, set_number)`. | PASS. |
| User-scoped query surface | `workouts.sql` scopes DailyLog reads/locks/summaries by `user_id`; workout exercise CRUD/reorder queries scope by `user_id` plus `daily_log_id` or row id; set mutations scope through a previously locked user-owned workout exercise/DailyLog. | PASS. |
| Transaction boundaries | `workout_repo.go` wraps aggregate mutations in `withTx`, begins pgx transactions, defers rollback, commits after mutation callback success, and performs child writes plus version increments inside the same transaction. | PASS. |
| Optimistic version lock path | Service methods validate `expectedVersion >= 0`, lock the relevant DailyLog through repository lock helpers, compare `locked.Version` before mutation, return conflict with current aggregate on mismatch, then increment and reload the aggregate after successful mutation. | PASS. |
| No stale save silently overwrites another tab | `requireVersion` returns `DailyLogConflictErr` when the locked version differs from `expectedVersion`; service tests cover stale notes and child mutation conflicts without repository mutation. | PASS. |
| FK/cascade/restrict behavior | Workout exercise deletion cascades sets and keeps DailyLog; deleting all exercises leaves DailyLog; exercise FK uses `ON DELETE RESTRICT`; set FK cascades through workout exercise. Repository tests cover cascade and empty DailyLog retention. | PASS. |
| Cross-user isolation | Add exercise rejects another user's Exercise Library row; add set rejects missing/wrong-user parent without mutation or version bump; update/delete set with wrong parent does not change the target set or version. | PASS. |
| Logging and privacy | `rg -n "log\\.|slog|zap|zerolog|fmt\\.Printf|fmt\\.Println" apps/api/internal/atlas/repository/postgres/workout_repo.go apps/api/internal/atlas/service/workout.go apps/api/internal/atlas/graph/resolver/workout.go` returned no matches. Error wrapping records operation names and wrapped errors, not note payload values. | PASS. |
| Reviewed non-blocker: `workout_sets` has no direct `user_id` column | User scope is inherited through `workout_sets -> workout_exercises -> daily_logs`. Repository set writes are reached only after locking a user-owned DailyLog by workout exercise/set id, and wrong-parent tests prove no mutation/version bump. This matches the approved design. | ACCEPTED. |
| Hard non-goals | No reviewed storage/API path introduces `cardio_entries`, cardio enums/placeholders, `body_weight`/`bodyWeight`, frontend, charts, or AI export. The dedicated no-scope proof remains `.tasks/WAVE-03/REGRESSION_NO_SCOPE_PROOF.md`. | PASS. |

## Verification Evidence

Focused Docker-backed storage/security command:

```bash
TEST_COMPOSE_PROJECT=atlas-w03-test TEST_POSTGRES_CONTAINER_NAME=atlas-w03-test-postgres TEST_REDIS_CONTAINER_NAME=atlas-w03-test-redis TEST_POSTGRES_VOLUME=atlas-w03-test-pg-data docker compose -f docker/docker-compose.test.yml up -d --wait postgres redis && cd apps/api && COVERAGE_GATE=1 API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable go test ./internal/repository/postgres -run 'TestWorkoutMigrations|TestWorkoutRepo_(GetOrCreateDailyLog_UniquePerUserDate|DailyLog_UserScopedIsolation|AddWorkoutExercise_RejectsOtherUsersExercise|AddWorkoutSet_MissingOrWrongUserParentDoesNotMutateOrVersion|UpdateWorkoutSet_WrongParentDoesNotChangeSetOrVersion|DeleteWorkoutSet_WrongParentDoesNotDeleteOrVersion|IncrementDailyLogVersion|LockedDailyLogHelpersExposeOwnedAggregateAndVersion|DeleteWorkoutExercise_CascadesSetsAndKeepsDailyLog|DeleteLastWorkoutExercise_KeepsEmptyDailyLog)' -count=1 -v
```

Result: PASS.

Observed evidence:

- Docker reported `atlas-w03-test-postgres` and `atlas-w03-test-redis` healthy.
- Goose reported current migration version `85`.
- Migration schema/constraint tests passed.
- User isolation, wrong-user parent, wrong-parent update/delete, cascade, empty DailyLog retention, version increment, and locked aggregate helper tests passed.
- No DB skip or unavailable-service output occurred under `COVERAGE_GATE=1`.

Focused service conflict/no-last-write-wins command:

```bash
cd apps/api && go test ./internal/atlas/service -run 'TestWorkoutService_(UpdateNotes_RejectsStaleVersion|RejectsNegativeExpectedVersionBeforeRepositoryMutation|ChildMutationsRejectStaleVersionWithCurrentAggregateAndNoMutation|AddSet_MissingParentReturnsNotFoundWithoutMutation)' -count=1 -v
```

Result: PASS.

Observed evidence:

- `TestWorkoutService_UpdateNotes_RejectsStaleVersion` passed.
- `TestWorkoutService_RejectsNegativeExpectedVersionBeforeRepositoryMutation` passed for notes, exercise, set, and reorder mutations.
- `TestWorkoutService_ChildMutationsRejectStaleVersionWithCurrentAggregateAndNoMutation` passed for add/update/remove/reorder exercise and add/update/remove/reorder set paths.
- `TestWorkoutService_AddSet_MissingParentReturnsNotFoundWithoutMutation` passed.

Focused logging/privacy command:

```bash
rg -n "log\\.|slog|zap|zerolog|fmt\\.Printf|fmt\\.Println" apps/api/internal/atlas/repository/postgres/workout_repo.go apps/api/internal/atlas/service/workout.go apps/api/internal/atlas/graph/resolver/workout.go
```

Result: PASS/no matches, exit 1 from `rg`.

## Findings And Handoff

Critical findings: none.

Important findings: none.

Minor findings: none.

Independent review: PASS. A read-only backend storage/security reviewer found no Critical, Important, or Minor findings and confirmed the reviewed non-blocker that `workout_sets` inherits user scope through the locked DailyLog/workout exercise path.

No Bead was created by this review because no storage/security issue requires follow-up.

This review does not replace `Atlas-qb2.3.3` GraphQL/API/generated artifact review, `Atlas-qb2.3.4` GRACE consistency review, or `Atlas-qb2.3.5` final gates/readiness packet.
