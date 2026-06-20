<!-- FILE: .tasks/WAVE-03/COVERAGE_MATRIX.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Map WAVE-03 Workout Diary source-plan requirements to committed tests, verification commands, deterministic test-environment assumptions, and follow-up coverage beads. -->
<!--   SCOPE: Coverage evidence for DailyLog, strength workout exercises, workout sets, Date scalar, generated-artifact gates, Docker-backed Postgres assumptions, explicit non-goals, and known non-W03 blockers; excludes product code, test code, shared GRACE XML edits, Beads DB edits, frontend, cardio, body weight, charts, and AI export. -->
<!--   DEPENDS: docs/superpowers/plans/2026-06-19-wave-03-workout-diary.md, docs/superpowers/specs/2026-06-19-wave-03-workout-diary-design.md, docs/verification-plan.xml, .tasks/WAVE-03/HANDOFF.md. -->
<!--   LINKS: M-API / V-M-API / WAVE-03 / Atlas-qb2.2.1. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Scope Boundary - Restates the exact WAVE-03 backend scope and non-goals under test. -->
<!--   Deterministic Test Environment - Documents Postgres, Docker, fixture, safe-skip, and cleanup assumptions. -->
<!--   Coverage Matrix - Maps source-plan requirements to committed test names or explicit follow-up beads. -->
<!--   Verification Command Matrix - Maps TEST-W03 commands, generated gates, and closure evidence. -->
<!--   Follow-up Coverage Beads - Routes remaining gaps to Atlas-qb2.2.2 through Atlas-qb2.2.6 and Atlas-qb2.3.1. -->
<!--   Current Blockers / Explicit Non-Blockers - Separates W03 coverage gaps from baseline GRACE/tooling blockers. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.0 - Added WAVE-03 coverage matrix and deterministic test-environment handoff. -->
<!-- END_CHANGE_SUMMARY -->

# WAVE-03 Coverage Matrix

## Scope Boundary

WAVE-03 covers the Atlas strength workout diary backend only:

- DailyLog canonical aggregate: `daily_logs`.
- Strength workout children: `workout_exercises` and `workout_sets`.
- GraphQL `Date`, `dailyLog`, `dailyLogs`, DailyLog notes, workout exercise mutations, workout set mutations, typed `DailyLogResult` errors, and optimistic aggregate versioning.
- Generated sqlc and Atlas gqlgen artifacts are replacement-gated by codegen, tests, and build commands, not by handwritten behavior tests.

Explicit non-goals preserved by the source plan and this matrix:

- No `cardio_entries`.
- No cardio CRUD, cardio API/schema fields, cardio enums, `CardioType`, `HeartRateZone`, placeholder cardio arrays, or fake empty cardio fields.
- No `body_weight` column and no `bodyWeight` API/persistence.
- No legacy `WorkoutDay` implementation names, except when quoting legacy source text.
- No automatic working weight progression.
- No frontend, public web, web-admin, charts, e1RM chart endpoint, AI export payload, backup/import, or starter workout template implementation.

## Deterministic Test Environment

Repository integration and migration tests require the safe test Postgres target:

- Default DSN from `apps/api/internal/testinfra/safe_targets.go`: `postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable`.
- Override env: `API_TEST_DATABASE_DSN`.
- Safe-target guard: `RequireSafePostgresDSN` rejects empty DSNs, non-Postgres schemes, `monorepo_dev`, wrong DB names, missing ports, dev port `7501`, and ports other than `TEST_POSTGRES_PORT`.
- Coverage mode: `COVERAGE_GATE=1` disables safe skips for unavailable Postgres and should turn environment failures into hard failures.
- Docker test stack: `docker/docker-compose.test.yml` exposes Postgres on `${TEST_POSTGRES_PORT:-17501}` and Redis on `${TEST_REDIS_PORT:-17502}`.

Canonical Docker setup for WAVE-03 repository evidence:

```bash
TEST_COMPOSE_PROJECT=atlas-w03-test \
TEST_POSTGRES_CONTAINER_NAME=atlas-w03-test-postgres \
TEST_REDIS_CONTAINER_NAME=atlas-w03-test-redis \
TEST_POSTGRES_VOLUME=atlas-w03-test-pg-data \
docker compose -f docker/docker-compose.test.yml up -d --wait postgres redis
```

Canonical repository run after the stack is healthy:

```bash
cd apps/api && \
API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable \
go test ./internal/repository/postgres -run 'TestWorkoutRepo|TestDailyLog' -count=1 -v
```

Coverage-gate variant for closure work:

```bash
cd apps/api && \
COVERAGE_GATE=1 \
API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable \
go test ./internal/repository/postgres -run 'TestWorkoutMigrations|TestWorkoutRepo|TestDailyLog|TestNew_ConnectsAndPings' -count=1 -v
```

Safe cleanup when the temporary stack is no longer needed:

```bash
TEST_COMPOSE_PROJECT=atlas-w03-test \
TEST_POSTGRES_CONTAINER_NAME=atlas-w03-test-postgres \
TEST_REDIS_CONTAINER_NAME=atlas-w03-test-redis \
TEST_POSTGRES_VOLUME=atlas-w03-test-pg-data \
docker compose -f docker/docker-compose.test.yml down -v
```

Fixture/helper assumptions:

- `workoutMigrationTestPool` applies goose migrations and skips only when Postgres is unavailable outside `COVERAGE_GATE=1`.
- `workoutRepoTestSetup` applies migrations, creates a real `postgresrepo.New` pool, truncates WAVE-03 and related Atlas tables, and creates an Atlas user.
- `truncateWorkoutRepoTables` uses `TRUNCATE workout_sets, workout_exercises, daily_logs, exercise_media, exercises, atlas_settings, atlas_users RESTART IDENTITY CASCADE`; only safe test DSNs may reach this destructive helper.
- `fakeWorkoutRepo` and `fakeExerciseRepo` in service tests isolate service behavior from database I/O while preserving aggregate state/version semantics.
- Resolver tests use `mockWorkoutService`; they prove resolver mapping/delegation, not repository or service internals.

Safe remedy ladder before escalating DB blockers:

1. Confirm Docker is running and the WAVE-03 compose stack is healthy with the command above.
2. Confirm `API_TEST_DATABASE_DSN` points at `monorepo_test` on `17501` or matches `TEST_POSTGRES_PORT`.
3. Re-run the focused repository command with `-v`; a skip is acceptable only outside coverage closure and only with the exact unavailable-service reason recorded.
4. Re-run with `COVERAGE_GATE=1` during coverage closure; skips are not acceptable there.
5. If migrations fail, capture the migration error, `docker compose ps`, and relevant Postgres logs before escalation.

## Coverage Matrix

Legend: `COVERED` means committed tests directly prove the behavior. `PARTIAL` means tests cover part of the behavior or a lower layer. `UNCOVERED` means a later coverage bead must add direct proof. `GATE` means covered by command/static generated-artifact replacement gates. `NON-GOAL` means the behavior must remain absent.

| Source-plan requirement | Current evidence | Status | Follow-up |
| --- | --- | --- | --- |
| `daily_logs` has `id`, `user_id`, `date`, `notes`, `version`, `created_at`, `updated_at`, `UNIQUE(user_id,date)`, version check, FK to `atlas_users`, date indexes. | `TestWorkoutMigrations_FilesExistWithGraceMarkup`, `TestWorkoutMigrations_DailyLogSchema`; `TestWorkoutRepo_GetOrCreateDailyLog_UniquePerUserDate`. | COVERED | `Atlas-qb2.2.2` should retain this as repository baseline. |
| WAVE-03 adds only `daily_logs`, `workout_exercises`, `workout_sets` for diary persistence; no `cardio_entries`, no `body_weight`, no `bodyWeight`. | `TestWorkoutMigrations_FilesExistWithGraceMarkup` checks prohibited migration tokens; `TestWorkoutMigrations_DailyLogSchema` checks no `body_weight`, no `bodyWeight`, and no `cardio_entries`; static grep target is `rg -n "cardio_entries|CardioType|HeartRateZone|body_weight|bodyWeight" apps/api/internal/atlas apps/api/internal/repository/postgres/migrations apps/api/internal/repository/postgres/queries`. | COVERED / NON-GOAL | `Atlas-qb2.2.5` owns cross-wave no-scope proof. |
| GraphQL/API exposes no cardio placeholders/enums/API fields and no body-weight fields. | `apps/api/internal/atlas/graph/schema/workouts.graphql` defines only strength diary fields; no-scope grep over `apps/api/internal/atlas` currently returns no product-code hits for `cardio_entries`, `CardioType`, `HeartRateZone`, `body_weight`, `bodyWeight`. | PARTIAL | `Atlas-qb2.2.4` should add schema/API assertions; `Atlas-qb2.2.5` should run final static no-scope proof. |
| Legacy `WorkoutDay` implementation names are not introduced. | Static proof command: `rg -n "WorkoutDay" apps/api/internal/atlas apps/api/internal/repository/postgres/migrations apps/api/internal/repository/postgres/queries`; current source search found no product-code hits. | GATE / NON-GOAL | `Atlas-qb2.2.5` should rerun this as part of final no-scope proof. |
| Broader non-goals stay absent: frontend/public web/web-admin routes or UI, charts/e1RM endpoint, AI export payload assembly, backup/import behavior, starter workout templates, and automatic working weight progression. | Final proof should combine narrow checks: `git diff --name-only -- apps/web apps/web-admin` for no frontend/admin changes; `rg -n "AIExport|e1RM|ChartEndpoint|WorkoutChart|StarterWorkoutTemplate|starter workout|workout template|working weight progression|automatic working weight|backup/import|BackupImport" apps/api/internal/atlas/graph/schema/workouts.graphql apps/api/internal/atlas/graph/resolver/workout.go apps/api/internal/atlas/models/workout.go apps/api/internal/atlas/models/workout_graphql.go apps/api/internal/atlas/service/workout.go apps/api/internal/atlas/repository/postgres/workout_repo.go apps/api/internal/repository/postgres/migrations/00083_daily_logs.sql apps/api/internal/repository/postgres/migrations/00084_workout_exercises.sql apps/api/internal/repository/postgres/migrations/00085_workout_sets.sql apps/api/internal/repository/postgres/queries/workouts.sql` for backend/API non-goal terms. Avoid generic `restore`, generic `backup`, broad `apps/api/internal/atlas`, and broad `apps/web` or `apps/web-admin` content greps because baseline text can produce false blockers. | GATE / NON-GOAL | `Atlas-qb2.2.5` owns the final broad no-scope proof. |
| `GetOrCreateDailyLogByDate` creates first DailyLog and reuses the existing `(user_id,date)` row. | `TestWorkoutRepo_GetOrCreateDailyLog_UniquePerUserDate`; service create path in `TestWorkoutService_UpdateNotes_CreatesDailyLogAtExpectedVersionZero`; absent non-zero version does not create in `TestWorkoutService_UpdateNotes_AbsentDateWithNonZeroExpectedVersionDoesNotCreateDailyLog` and `TestWorkoutService_AddExercise_AbsentDateWithNonZeroExpectedVersionDoesNotCreateDailyLog`. | COVERED | `Atlas-qb2.2.2` can add lock-specific coverage if needed. |
| DailyLog repository access is user-scoped and blocks cross-user aggregate reads. | `TestWorkoutRepo_DailyLog_UserScopedIsolation` creates two users on the same date, proves separate DailyLog IDs, and verifies another user cannot load the first user's aggregate. | COVERED for repository isolation. | Service/API cross-user coverage can remain in `Atlas-qb2.2.5` regression/no-scope checks if broader auth isolation proof is requested. |
| `dailyLog(date:)` query returns existing data but must not create a row when absent. | Resolver delegation is covered by `TestDailyLogResolver_DelegatesAuthenticatedDailyLog`; service implementation reads via `GetDailyLogByDate`. No committed test directly proves absent-date query leaves the DB/fake repository unchanged. | UNCOVERED | `Atlas-qb2.2.3` for service no-create test; `Atlas-qb2.2.4` for resolver/API contract no-create path. |
| `dailyLogs(from,to)` summary query and invalid date range behavior. | Service has `ListDailyLogSummaries` implementation; no committed test directly covers summary success or `from > to` validation. | UNCOVERED | `Atlas-qb2.2.3` service behavior; `Atlas-qb2.2.4` resolver/API signatures. |
| DailyLog notes create/update increments version and stale version returns typed conflict with current aggregate. | `TestWorkoutService_UpdateNotes_CreatesDailyLogAtExpectedVersionZero`; `TestWorkoutService_UpdateNotes_RejectsStaleVersion`; resolver conflict mapping in `TestUpdateDailyLogNotesResolver_MapsConflictError`. | COVERED for create/stale conflict; PARTIAL for existing-row update. | `Atlas-qb2.2.3` should add existing-row successful update if closure requires it. |
| DailyLog notes explicit null/clear behavior. | Code path accepts `notes *string`; no committed service/resolver test directly calls `UpdateDailyLogNotes(..., nil)` and asserts cleared notes. | UNCOVERED | `Atlas-qb2.2.3` service null clear; `Atlas-qb2.2.4` resolver/API null argument proof. |
| Service rule `expectedVersion >= 0` rejects negative versions before mutation. | Implementation has `validateExpectedVersion(expectedVersion)` returning validation error for `< 0`; no committed test directly passes a negative `expectedVersion` through service or GraphQL resolver paths. | UNCOVERED | `Atlas-qb2.2.3` for service validation/no-mutation proof; `Atlas-qb2.2.4` for API/resolver contract proof that negative GraphQL `expectedVersion` maps to typed validation. |
| Workout exercise add: validates existing exercise, creates DailyLog at expected version zero, appends/positions, captures working weight snapshot. | `TestWorkoutService_AddExercise_RequiresExistingExercise`; `TestWorkoutService_AddExercise_CapturesWorkingWeightSnapshot`; `TestWorkoutRepo_AddWorkoutExercise_CapturesWorkingWeightSnapshot`; absent-date guards in service tests. | COVERED for required exercise/snapshot/create guards; PARTIAL for insert-at-position success. | `Atlas-qb2.2.3` should add insert-at-position success if needed. |
| Duplicate `exercise_id` instances are allowed in one DailyLog. | `TestWorkoutRepo_AddWorkoutExercise_AllowsDuplicateExercise`; `TestWorkoutService_AddExercise_AllowsDuplicateExerciseID`; schema test confirms no unique constraint on `(daily_log_id, exercise_id)`. | COVERED | None. |
| Workout exercise update: position move/reindex, no-op validation, exercise-level notes update and null/clear, snapshot stays historical. | Repo move/reindex: `TestWorkoutRepo_UpdateWorkoutExercise_MoveReindexesContiguously`; service no-op validation: `TestWorkoutService_UpdateExercise_RejectsEmptyInputWithoutVersionChange`, `TestWorkoutService_UpdateExercise_RejectsSamePositionOnlyWithoutVersionChange`; resolver explicit null mapping: `TestUpdateWorkoutExerciseResolver_MapsExplicitNullNotes`. No direct service/repository assertion for notes persistence/clear or snapshot immutability after source exercise update. | PARTIAL / UNCOVERED for notes persistence and snapshot immutability after exercise update. | `Atlas-qb2.2.2` repository notes/snapshot persistence; `Atlas-qb2.2.3` service notes/null and snapshot immutability. |
| Workout exercise remove: cascades sets, reindexes remaining exercises, keeps empty DailyLog. | `TestWorkoutRepo_DeleteWorkoutExercise_CascadesSetsAndKeepsDailyLog`; `TestWorkoutRepo_DeleteLastWorkoutExercise_KeepsEmptyDailyLog`; `TestWorkoutService_RemoveExercise_KeepsEmptyDailyLog`. | COVERED | None. |
| Workout exercise reorder: success reindexes contiguously; missing/duplicate/foreign/extra IDs rejected without version change. | Success/reindex: `TestWorkoutRepo_ReorderWorkoutExercises_ReindexesContiguously`; validation failures: `TestWorkoutService_ReorderExercises_RejectsMissingDuplicateOrForeignIDs`. | COVERED for repo success and service invalid IDs; PARTIAL for service success/version increment. | `Atlas-qb2.2.3` can add service reorder success/version proof. |
| Wrong-user exercise cannot be attached. | `TestWorkoutRepo_AddWorkoutExercise_RejectsOtherUsersExercise` proves no row and no version change. | COVERED | `Atlas-qb2.2.5` should keep WAVE-02/user-scope regression evidence. |
| Workout set add: validates `weight > 0`, `reps > 0`, optional `rpe` in `1..10`, optional `rir` in `0..10`; appends or inserts with contiguous set numbers. | Service validation: `TestWorkoutService_AddSet_ValidatesWeightRepsRpeRir`; DB schema/invalid bounds: `TestWorkoutMigrations_WorkoutSetSchema`, `TestWorkoutMigrations_WorkoutSetValidationRejectsInvalidBounds`; repository constraint smoke: `TestWorkoutRepo_AddWorkoutSet_ValidatesDBConstraints`. | COVERED for validation and append default; PARTIAL for insert-at-set-number success. | `Atlas-qb2.2.2` repository insert success; `Atlas-qb2.2.3` service insert success if required. |
| `addWorkoutSet` parent workout exercise must belong to current user; missing or foreign parent returns not-found without mutation/version bump. | Resolver not-found envelope is sampled by `TestWorkoutSetResolvers_MapNotFoundError` for add set; repository/service tests do not directly prove missing-parent or wrong-user parent ownership for `AddWorkoutSet`. | PARTIAL / UNCOVERED for repository and service ownership behavior. | `Atlas-qb2.2.2` repository missing/wrong-user parent proof; `Atlas-qb2.2.3` service parent not-found/no-mutation proof; `Atlas-qb2.2.4` API not-found contract proof. |
| Workout set update: weight/reps/rpe/rir/notes and optional set-number move/reindex; explicit null clears nullable fields. | Move/reindex: `TestWorkoutRepo_UpdateWorkoutSet_MoveReindexesContiguously`, `TestWorkoutService_UpdateSet_ReindexesWhenSetNumberChanges`; no-op validation: `TestWorkoutService_UpdateSet_RejectsEmptyInputWithoutVersionChange`, `TestWorkoutService_UpdateSet_RejectsSameSetNumberOnlyWithoutVersionChange`; resolver explicit null mapping: `TestUpdateWorkoutSetResolver_MapsExplicitNullNullableFields`; wrong parent: `TestWorkoutRepo_UpdateWorkoutSet_WrongParentDoesNotChangeSetOrVersion`. No direct service/repo assertions for updating weight/reps/rpe/rir/notes values or clearing nullable fields through persistence. | PARTIAL / UNCOVERED for value updates and notes/null persistence. | `Atlas-qb2.2.2` repository update values/null persistence; `Atlas-qb2.2.3` service update validation/value/null paths; `Atlas-qb2.2.4` GraphQL API null paths. |
| Workout set remove: deletes only the target set, wrong parent/foreign IDs do not delete or bump version, remaining sets reindex. | Wrong parent: `TestWorkoutRepo_DeleteWorkoutSet_WrongParentDoesNotDeleteOrVersion`; resolver not-found mapping for remove set in `TestWorkoutSetResolvers_MapNotFoundError`. No direct repository/service success test for remove-one-set reindex. | PARTIAL / UNCOVERED for success remove/reindex. | `Atlas-qb2.2.2` repository success remove/reindex; `Atlas-qb2.2.3` service success remove/reindex. |
| Workout set reorder: success reindexes contiguously; missing/duplicate/foreign/extra IDs rejected. | Success: `TestWorkoutRepo_ReorderWorkoutSets_ReindexesContiguously`; validation failures: `TestWorkoutService_ReorderSets_RejectsMissingDuplicateOrForeignIDs`. | COVERED for repo success and service invalid IDs; PARTIAL for service success/version increment. | `Atlas-qb2.2.3` can add service reorder success/version proof. |
| Optimistic concurrency: stale DailyLog version returns conflict, no last-write-wins, child mutations validate `expectedVersion` before mutation, version increments are reusable for future child domains. | DailyLog stale conflict: `TestWorkoutService_UpdateNotes_RejectsStaleVersion`; absent non-zero conflict tests; reusable increment: `TestWorkoutRepo_IncrementDailyLogVersion`; multiple success tests assert incremented versions. No direct stale-version tests for each child mutation. | PARTIAL / UNCOVERED for child stale conflicts. | `Atlas-qb2.2.3` owns child conflict/no-last-write-wins coverage. |
| Resolver auth, delegation, typed domain error mapping, unexpected internal error behavior, explicit null mapping. | Auth: `TestDailyLogResolver_UnauthorizedReturnsAuthError`; delegation: `TestDailyLogResolver_DelegatesAuthenticatedDailyLog`; conflict: `TestUpdateDailyLogNotesResolver_MapsConflictError`; validation: `TestAddWorkoutExerciseResolver_MapsValidationError`; not found: `TestWorkoutSetResolvers_MapNotFoundError`; internal: `TestWorkoutResolvers_DoNotLeakUnexpectedErrors`; explicit null: `TestUpdateWorkoutExerciseResolver_MapsExplicitNullNotes`, `TestUpdateWorkoutSetResolver_MapsExplicitNullNullableFields`. | COVERED for sampled mappings; PARTIAL for every GraphQL operation signature/success path. | `Atlas-qb2.2.4` owns resolver/API contract expansion. |
| Date scalar accepts strict `YYYY-MM-DD`, rejects timestamps, marshals quoted dates, and marshals zero value as null. | `TestDate_UnmarshalStrictYYYYMMDD`; `TestDate_RejectsTimestamp`; `TestDate_MarshalGQLWritesQuotedDate`; `TestDate_MarshalGQLZeroValueWritesNull`. | COVERED | `Atlas-qb2.2.4` should include generated binding/schema compile proof. |
| Generated sqlc and gqlgen artifacts stay current. | `TEST-W03-009`: `bunx nx run api:codegen && bunx nx run api:codegen:atlas`; handoff records both passed with `--skip-nx-cache`. | GATE | `Atlas-qb2.2.6` reruns final closure commands. |
| API package and build gates prove generated artifacts are usable. | `TEST-W03-012`: `bunx nx test api && bunx nx build api`; handoff records both passed with `--skip-nx-cache`. | GATE | `Atlas-qb2.2.6` reruns final closure commands. |
| Docker-backed repository tests prove committed migrations and sqlc queries against real Postgres. | `TEST-W03-011`; handoff says Docker-backed repository integration passed, migrations applied through `00085`, and 17 `TestWorkoutRepo_*` tests ran with no DB skip. | COVERED for handoff run; closure rerun required. | `Atlas-qb2.2.6` reruns and records exact output. |

## Verification Command Matrix

| ID / command | Purpose | Required closure evidence | Current note |
| --- | --- | --- | --- |
| `TEST-W03-001`: `cd apps/api && go test ./internal/repository/postgres -run 'TestWorkoutMigrations|TestNew_ConnectsAndPings' -count=1` | Migration schema and baseline Postgres connectivity. | PASS, or exact DB setup blocker. Under coverage closure, run with safe Docker DSN and prefer `COVERAGE_GATE=1`. | Migration tests cover W03 schema/non-goal table absence. |
| `TEST-W03-007`: `cd apps/api && go test ./internal/atlas/graph/resolver -run 'Test.*DailyLog|Test.*Workout' -count=1` | Resolver auth/delegation/error/null mapping. | PASS with test names visible if run with `-v`. | Existing resolver coverage is sampled; `Atlas-qb2.2.4` expands API contract proof. |
| `TEST-W03-009`: `bunx nx run api:codegen && bunx nx run api:codegen:atlas` | sqlc and Atlas gqlgen replacement gate. | Both commands exit 0 and produce no unexpected generated drift. | Handoff records passed with `--skip-nx-cache`. |
| `TEST-W03-010`: `cd apps/api && go test ./internal/atlas/models -run TestDate -count=1 && go test ./internal/atlas/service -run 'TestWorkoutService' -count=1 && go test ./internal/atlas/graph/resolver -run 'Test.*DailyLog|Test.*Workout' -count=1` | Date, service, resolver focused package tests. | PASS; if failures move downstream, record exact failing test and assign to `.2.3` or `.2.4`. | Current matrix identifies direct service/resolver coverage gaps. |
| `TEST-W03-011`: Docker up plus `API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable go test ./internal/repository/postgres -run 'TestWorkoutRepo|TestDailyLog' -count=1 -v` | Real Postgres repository integration. | Docker stack health, no DB skip, migrations through `00085`, PASS output. | Handoff records passed and no skip. |
| `TEST-W03-012`: `bunx nx test api && bunx nx build api` | API target and build closure. | Both exit 0 after generated artifacts are refreshed. | Handoff records passed with `--skip-nx-cache`. |
| Artifact check: `rg -n "UNCOVERED|Atlas-qb2.2|TEST-W03|cardio_entries|body_weight|docker compose" .tasks/WAVE-03/COVERAGE_MATRIX.md` | Proves matrix contains gap markers, follow-up beads, verification IDs, non-goal tokens, and Docker commands. | Command exits 0 with relevant lines. | Required for Atlas-qb2.2.1 only. |
| Artifact diff: `git add -N .tasks/WAVE-03/COVERAGE_MATRIX.md && git diff -- .tasks/WAVE-03/COVERAGE_MATRIX.md` | Proves write scope and reviewability of this initially untracked artifact. | Diff shows only this new file. | Required for Atlas-qb2.2.1 only. |

Generated-artifact replacement gates:

- sqlc generated files under `apps/api/internal/repository/postgres/generated/**`: `bunx nx run api:codegen`, `bunx nx build api`, and Docker-backed repository tests.
- Atlas gqlgen generated files under `apps/api/internal/atlas/graph/generated/**` and generated resolver shims: `bunx nx run api:codegen:atlas`, resolver tests, `bunx nx test api`, and `bunx nx build api`.
- Admin GraphQL generated outputs under `apps/api/internal/graph/generated.go` and `apps/api/internal/graph/model/models_gen.go` are intentionally ignored; clean-checkout API reproducibility depends on `bunx nx run api:codegen` before `bunx nx test api` or `bunx nx build api`.

## Follow-up Coverage Beads

| Bead | Scope from Beads metadata | Matrix assignments |
| --- | --- | --- |
| `Atlas-qb2.2.2` | W03 coverage: repository integration edge cases. | Add/retain repository tests for absent `GetDailyLogByDate`, lock/version helpers, exercise notes persistence/clear, historical snapshot immutability, `AddWorkoutSet` missing/wrong-user parent ownership, set update values/null persistence, set remove success/reindex, insert-at-position/set-number success, and wrong-parent/foreign-ID persistence invariants. |
| `Atlas-qb2.2.3` | W03 coverage: service behavior and conflict paths. | Add service tests for `dailyLog` no-create, `dailyLogs` date range, existing notes update, notes null clear, negative `expectedVersion >= 0` validation/no mutation, child stale `expectedVersion` conflicts, `AddWorkoutSet` parent workout exercise not-found/no mutation, exercise notes/null, set update values/null, set remove success, reorder success/version, and no last-write-wins paths. |
| `Atlas-qb2.2.4` | W03 coverage: GraphQL resolver and API contract paths. | Add resolver/API/schema tests for all operation signatures, success mapping samples, `dailyLog` no-create behavior, `dailyLogs` range behavior, updateDailyLogNotes null mapping, negative `expectedVersion` validation mapping, `AddWorkoutSet` parent workout exercise not-found mapping, Date scalar generated binding, and no cardio/body placeholders in schema. |
| `Atlas-qb2.2.5` | W03 coverage: cross-wave regression and no-scope proof. | Run WAVE-01/WAVE-02 focused regression checks and final no-scope proof for all WAVE-03 boundaries: no `cardio_entries`, `CardioType`, `HeartRateZone`, `body_weight`, `bodyWeight`, legacy `WorkoutDay` implementation names, frontend/public web/web-admin routes or UI, charts/e1RM endpoint, AI export payload assembly, backup/import behavior, starter workout templates, or automatic working weight progression. This matrix intentionally extends `.2.5` static proof to all W03 source-plan non-goals, even though Beads metadata emphasizes cardio/body-weight. Use targeted path checks for `apps/web` and `apps/web-admin`, and targeted domain terms for API/backend paths. |
| `Atlas-qb2.2.6` | W03 coverage: coverage epic closure commands. | Rerun `api:codegen`, `api:codegen:atlas`, WAVE-03 package tests, Docker repository tests, `bunx nx test api`, and `bunx nx build api`; record exact output and no-skip DB evidence. |
| `Atlas-qb2.3.1` | W03 QA: source traceability audit. | Use this matrix plus final command evidence to audit source-plan traceability after coverage beads close. |

## Current Blockers / Explicit Non-Blockers

Current W03 coverage blockers:

- `UNCOVERED`: direct absent-date `dailyLog(date:)` no-create tests.
- `UNCOVERED`: `dailyLogs(from,to)` summary/date-range behavior tests.
- `UNCOVERED`: DailyLog notes explicit null/clear tests.
- `UNCOVERED`: negative `expectedVersion >= 0` service validation and API/resolver mapping tests.
- `UNCOVERED`: child mutation stale-version conflict tests across exercise/set operations.
- `UNCOVERED`: `addWorkoutSet` parent workout exercise missing/wrong-user ownership tests below the sampled resolver not-found envelope.
- `UNCOVERED`: exercise notes persistence/null clear and historical snapshot immutability after Exercise Library source changes.
- `UNCOVERED`: set update values/null persistence and set remove success/reindex tests.
- `UNCOVERED`: full GraphQL operation signature/success/no-placeholder assertions.

Explicit non-blockers for W03 coverage matrix creation:

- The W03 implementation handoff already records focused codegen, test, build, and Docker repository evidence; this artifact does not rerun heavy gates.
- `.tasks/WAVE-03/HANDOFF.md` records baseline GRACE lint issues after the in-scope W03 markup issue was fixed: `bunx @osovv/grace-cli lint --path .` still reported 32 baseline issues, 24 errors and 8 warnings, in WAVE-01/WAVE-02/generated/skills surfaces. Those are not W03 coverage matrix blockers.
- The temporary Docker stack `atlas-w03-test` may remain running from handoff; it is an environment cleanup item, not a product/code blocker, unless it prevents deterministic reruns.
- This task intentionally does not edit product code, tests, `docs/*.xml`, generated files, Beads DB, git config, cardio/body-weight/frontend scope, or commits.
