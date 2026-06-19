<!-- FILE: .tasks/WAVE-03/HANDOFF.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Record WAVE-03 implementation handoff evidence for the Atlas DailyLog and strength workout diary backend. -->
<!--   SCOPE: Implemented scope, verification commands, generated artifact status, WAVE-04 compatibility, and known risks for WAVE-03; excludes MR creation and WAVE-04 cardio delivery. -->
<!--   DEPENDS: docs/superpowers/plans/2026-06-19-wave-03-workout-diary.md, docs/development-plan.xml, docs/knowledge-graph.xml, docs/verification-plan.xml, Beads Atlas-qb2. -->
<!--   LINKS: M-API / V-M-API / WAVE-03 / Atlas-qb2.1.10. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Scope Implemented - Summarizes WAVE-03 backend behavior and explicit non-goals. -->
<!--   Verification Evidence - Lists focused command evidence from implementation closeout. -->
<!--   Generated Artifact Status - Documents sqlc/gqlgen gate expectations and ignored admin GraphQL generated outputs. -->
<!--   WAVE-04 Compatibility - Records DailyLog container contracts reserved for cardio attachment later. -->
<!--   Known Risks - Captures follow-up risks for coverage and final QA beads. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.0 - Added WAVE-03 handoff evidence after focused integration verification. -->
<!-- END_CHANGE_SUMMARY -->

# WAVE-03 Handoff

## Scope Implemented

WAVE-03 implements the Atlas DailyLog and strength workout diary backend:

- `daily_logs` canonical daily container with `id`, `user_id`, `date`, `version`, `created_at`, `updated_at`, nullable notes, and `UNIQUE(user_id, date)`.
- `workout_exercises` and `workout_sets` persistence with user scoping, ordering, FK behavior, snapshots, and repository transaction helpers.
- DailyLog get-or-create behavior for mutations and reusable aggregate version increments.
- Date scalar, workout domain models, repository, service validation, optimistic concurrency, typed errors, GraphQL schema, generated Atlas gqlgen artifacts, handwritten resolvers, and runtime wiring.

Explicit non-goals preserved:

- No `cardio_entries` table.
- No cardio CRUD, schema, resolvers, validation, placeholder fields, `CardioType`, or `HeartRateZone`.
- No `body_weight` or `bodyWeight` persistence/API.
- No frontend, charts, or AI export behavior.

## Verification Evidence

Focused Task 9 evidence:

- `bunx nx run api:codegen --skip-nx-cache` passed after deleting ignored admin gqlgen outputs.
- `bunx nx run api:codegen:atlas --skip-nx-cache` passed.
- `bunx nx test api --skip-nx-cache` passed.
- `bunx nx build api --skip-nx-cache` passed.
- `cd apps/api && go test ./internal/atlas/models -run TestDate -count=1` passed.
- `cd apps/api && go test ./internal/atlas/service -run 'TestWorkoutService' -count=1` passed.
- `cd apps/api && go test ./internal/atlas/graph/resolver -run 'Test.*DailyLog|Test.*Workout' -count=1` passed.
- Docker-backed repository integration passed with:

```bash
TEST_COMPOSE_PROJECT=atlas-w03-test \
TEST_POSTGRES_CONTAINER_NAME=atlas-w03-test-postgres \
TEST_REDIS_CONTAINER_NAME=atlas-w03-test-redis \
TEST_POSTGRES_VOLUME=atlas-w03-test-pg-data \
docker compose -f docker/docker-compose.test.yml up -d --wait postgres redis

cd apps/api && \
API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable \
go test ./internal/repository/postgres -run 'TestWorkoutRepo|TestDailyLog' -count=1 -v
```

That run applied migrations through `00085` and executed 17 `TestWorkoutRepo_*` tests with no DB skip.

Docs/GRACE verification for this handoff:

- `xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml` passed.
- `grace` was not installed on PATH, so the handoff used `bunx @osovv/grace-cli lint --path .`.
- Initial GRACE lint found one in-scope WAVE-03 issue: `apps/api/internal/atlas/service/workout_service_test.go` missing `START_MODULE_MAP`.
- After adding that module map, `bunx @osovv/grace-cli lint --path .` still reports 32 baseline issues: 24 errors and 8 warnings. Remaining errors are pre-existing Atlas WAVE-01/WAVE-02 markup gaps such as missing module maps in Atlas health, PIN, settings, exercise, Redis, and generated sqlc files; none are introduced by the WAVE-03 handoff patch.

Exact remaining `bunx @osovv/grace-cli lint --path .` output excerpt:

```text
GRACE Lint Report
=================
Root: /Users/vlad/Develop/Atlas/.worktrees/wave-03-workout-diary
Profile: standard
Code files checked: 290
Governed files checked: 205
XML files checked: 5
Issues: 32 (errors: 24, warnings: 8)

Errors:
- [markup.missing-module-map] apps/api/internal/handler/atlas_health.go Governed files with ROLE RUNTIME and MAP_MODE EXPORTS must include a paired MODULE_MAP section.
- [markup.missing-module-map] apps/api/internal/handler/atlas_pin_auth.go Governed files with ROLE RUNTIME and MAP_MODE EXPORTS must include a paired MODULE_MAP section.
- [markup.missing-module-map] apps/api/internal/handler/atlas_media_test.go Governed files with ROLE TEST and MAP_MODE SUMMARY must include a paired MODULE_MAP section.
- [markup.missing-change-summary] apps/api/internal/handler/atlas_health_test.go Governed files must include a paired CHANGE_SUMMARY section.
- [markup.missing-module-map] apps/api/internal/repository/postgres/generated/exercises.sql.go Governed files with ROLE CONFIG and MAP_MODE SUMMARY must include a paired MODULE_MAP section.
- [markup.missing-module-map] apps/api/internal/repository/postgres/generated/querier.go Governed files with ROLE CONFIG and MAP_MODE SUMMARY must include a paired MODULE_MAP section.
- [markup.missing-module-map] apps/api/internal/repository/postgres/generated/atlas_settings.sql.go Governed files with ROLE CONFIG and MAP_MODE SUMMARY must include a paired MODULE_MAP section.
- [markup.missing-module-map] apps/api/internal/repository/postgres/exercise_repo_test.go Governed files with ROLE TEST and MAP_MODE SUMMARY must include a paired MODULE_MAP section.
- [markup.missing-module-map] apps/api/internal/repository/postgres/queries/exercises.sql Governed files with ROLE CONFIG and MAP_MODE SUMMARY must include a paired MODULE_MAP section.
- [markup.missing-module-map] apps/api/internal/repository/postgres/queries/atlas_settings.sql Governed files with ROLE CONFIG and MAP_MODE SUMMARY must include a paired MODULE_MAP section.
- [markup.missing-module-map] apps/api/internal/atlas/middleware/pin_guard.go Governed files with ROLE RUNTIME and MAP_MODE EXPORTS must include a paired MODULE_MAP section.
- [markup.missing-change-summary] apps/api/internal/atlas/middleware/pin_guard_test.go Governed files must include a paired CHANGE_SUMMARY section.
- [markup.missing-module-map] apps/api/internal/atlas/repository/redis/pin_attempt_store.go Governed files with ROLE RUNTIME and MAP_MODE EXPORTS must include a paired MODULE_MAP section.
- [markup.missing-module-map] apps/api/internal/atlas/repository/redis/pin_session_store.go Governed files with ROLE RUNTIME and MAP_MODE EXPORTS must include a paired MODULE_MAP section.
- [markup.missing-module-map] apps/api/internal/atlas/repository/postgres/settings_repo.go Governed files with ROLE RUNTIME and MAP_MODE EXPORTS must include a paired MODULE_MAP section.
- [markup.missing-module-map] apps/api/internal/atlas/graph/resolver/exercise.go Governed files with ROLE RUNTIME and MAP_MODE EXPORTS must include a paired MODULE_MAP section.
- [markup.missing-change-summary] apps/api/internal/atlas/graph/resolver/exercise_test.go Governed files must include a paired CHANGE_SUMMARY section.
- [markup.missing-module-map] apps/api/internal/atlas/graph/resolver/exercise_test.go Governed files with ROLE TEST and MAP_MODE SUMMARY must include a paired MODULE_MAP section.
- [markup.missing-module-map] apps/api/internal/atlas/service/settings_service.go Governed files with ROLE RUNTIME and MAP_MODE EXPORTS must include a paired MODULE_MAP section.
- [markup.missing-module-map] apps/api/internal/atlas/service/exercise_service_test.go Governed files with ROLE TEST and MAP_MODE SUMMARY must include a paired MODULE_MAP section.
- [markup.missing-module-map] apps/api/internal/atlas/service/settings_service_test.go Governed files with ROLE TEST and MAP_MODE SUMMARY must include a paired MODULE_MAP section.
- [markup.missing-module-map] apps/api/internal/atlas/service/pin_service.go Governed files with ROLE RUNTIME and MAP_MODE EXPORTS must include a paired MODULE_MAP section.
- [markup.missing-module-map] apps/api/internal/atlas/service/pin_service_test.go Governed files with ROLE TEST and MAP_MODE SUMMARY must include a paired MODULE_MAP section.
- [markup.missing-module-map] apps/api/internal/atlas/service/bootstrap_service.go Governed files with ROLE RUNTIME and MAP_MODE EXPORTS must include a paired MODULE_MAP section.

Warnings:
- [analysis.heuristic-export-surface] .agents/skills/verify-technical-docs/scripts/validate_technical_verified.py The python adapter inferred exports heuristically for this file. Exact MODULE_MAP parity may require explicit file ROLE/MAP_MODE or stronger language-specific export declarations.
- [analysis.heuristic-export-surface] .agents/skills/verify-technical-docs/scripts/scaffold_technical_verified.py The python adapter inferred exports heuristically for this file. Exact MODULE_MAP parity may require explicit file ROLE/MAP_MODE or stronger language-specific export declarations.
- [analysis.heuristic-export-surface] .agents/skills/detail-prd-wave/scripts/validate_detail_prd_wave.py The python adapter inferred exports heuristically for this file. Exact MODULE_MAP parity may require explicit file ROLE/MAP_MODE or stronger language-specific export declarations.
- [analysis.heuristic-export-surface] .agents/skills/detail-prd-wave/scripts/scaffold_detail_prd_wave.py The python adapter inferred exports heuristically for this file. Exact MODULE_MAP parity may require explicit file ROLE/MAP_MODE or stronger language-specific export declarations.
- [analysis.heuristic-export-surface] .agents/skills/plan-backend-waves/scripts/scaffold_backend_waves.py The python adapter inferred exports heuristically for this file. Exact MODULE_MAP parity may require explicit file ROLE/MAP_MODE or stronger language-specific export declarations.
- [analysis.heuristic-export-surface] .agents/skills/plan-backend-waves/scripts/validate_backend_waves.py The python adapter inferred exports heuristically for this file. Exact MODULE_MAP parity may require explicit file ROLE/MAP_MODE or stronger language-specific export declarations.
- [analysis.heuristic-export-surface] .agents/skills/decompose-prd-waves/scripts/validate_prd_waves.py The python adapter inferred exports heuristically for this file. Exact MODULE_MAP parity may require explicit file ROLE/MAP_MODE or stronger language-specific export declarations.
- [analysis.heuristic-export-surface] .agents/skills/decompose-prd-waves/scripts/scaffold_prd_waves.py The python adapter inferred exports heuristically for this file. Exact MODULE_MAP parity may require explicit file ROLE/MAP_MODE or stronger language-specific export declarations.
```

## Generated Artifact Status

Atlas generated artifacts are tracked and refreshed by `bunx nx run api:codegen:atlas`.

Admin GraphQL generated outputs under `apps/api/internal/graph/generated.go` and `apps/api/internal/graph/model/models_gen.go` are intentionally ignored by `.gitignore`; clean-checkout API build/test reproducibility depends on running `bunx nx run api:codegen` before `bunx nx test api` or `bunx nx build api`.

## WAVE-04 Compatibility

DailyLog is ready to serve as the shared daily container for WAVE-04 cardio attachment:

- DailyLog identity is stable by `(user_id, date)`.
- DailyLog aggregate `version` is reusable by child mutations.
- Repository/service behavior supports create-or-reuse daily container semantics.
- Future cardio mutations should increment DailyLog version when cardio entries change.

## Known Risks

- Full coverage work is intentionally deferred to `Atlas-qb2.2`. If coverage gates flag generated Atlas gqlgen files, handle that in the coverage epic with the repository coverage policy and replacement gates.
- Temporary Docker stack `atlas-w03-test` is running for downstream final checks and should be stopped after WAVE-03 closeout if no longer needed.
- No MR has been created in this handoff; the branch remains the WAVE-03 implementation branch.
