# WAVE-03 testing-exit Planner Attempt 1

## Sources Read
- docs/technical-verified/testing-and-delivery.md
- docs/prd-wave-details/waves/wave-02.md
- docs/prd-wave-details/waves/wave-01.md
- docs/product-verified/acceptance-criteria.md
- docs/product-verified/edge-cases.md
- docs/product-verified/business-rules.md

## Selected Backend Wave Boundary
WAVE-03 adds four new database tables, sqlc queries, repository adapters, a workout service, GraphQL schema, and GraphQL resolvers. Testing must cover all layers: repository, service, resolver, and integration.

## Neighboring Backend Wave Fit
- WAVE-01 tests: admin auth regression test (TEST-W01-008) must still pass after WAVE-03 changes.
- WAVE-02 tests: exercise tests must still pass; allExercises query still works.
- WAVE-03 test helpers will create exercises via WAVE-02 allExercises or exerciseById for snapshot testing.

## Frontend Pages Context
- No frontend testing in this wave. Backend tests only.

## Proposed Details

### Exit Criteria

| EC ID | Description |
| --- | --- |
| EC-W03-001 | All acceptance criteria (AC-W03-001 through AC-W03-030) pass via focused tests. |
| EC-W03-002 | gqlgen codegen produces valid Go code for workout schema without drift. |
| EC-W03-003 | sqlc codegen produces valid Go code for workout queries without drift. |
| EC-W03-004 | Migrations 00082-00085 apply and roll back in sequence without errors. |
| EC-W03-005 | DailyLog upsert round-trip: create with exercises/sets/cardio, read back, update, delete. |
| EC-W03-006 | Working weight snapshot correctly read from Exercise.workingWeight at add time. |
| EC-W03-007 | WorkoutExercise order returned correctly ordered in DailyLog query. |
| EC-W03-008 | WorkoutSet setNumber returned correctly ordered in WorkoutExercise query. |
| EC-W03-009 | CardioEntry CRUD within DailyLog works correctly. |
| EC-W03-010 | Cascade delete: deleting DailyLog removes all nested WorkoutExercises, WorkoutSets, CardioEntries. |
| EC-W03-011 | All WAVE-03 GraphQL operations return AuthError without valid PIN session. |
| EC-W03-012 | No sensitive content (notes, comments) appears in application logs. |
| EC-W03-013 | Input validation enforced: negative weight rejected, zero reps rejected, invalid dates rejected. |
| EC-W03-014 | DailyLog query by date returns empty result (not error) for nonexistent dates. |
| EC-W03-015 | WAVE-01 admin auth and health test suites still pass after WAVE-03 changes. |
| EC-W03-016 | WAVE-02 exercise tests still pass after WAVE-03 changes. |
| EC-W03-017 | Lint passes for all changed packages. |
| EC-W03-018 | Typecheck passes for Go API. |

### Verification Obligations

| Test ID | Description | Type | Command |
| --- | --- | --- | --- |
| TEST-W03-001 | DailyLog repository: create, get by date, update, delete | unit | bunx nx run api:test -- --run '(?i)daily_log_repo' |
| TEST-W03-002 | WorkoutExercise repository: create, list by daily log, update order, delete | unit | bunx nx run api:test -- --run '(?i)workout_exercise_repo' |
| TEST-W03-003 | WorkoutSet repository: create, list by exercise, update, delete | unit | bunx nx run api:test -- --run '(?i)workout_set_repo' |
| TEST-W03-004 | CardioEntry repository: create, list by daily log, update, delete | unit | bunx nx run api:test -- --run '(?i)cardio_entry_repo' |
| TEST-W03-005 | Workout service: DailyLogByDate returns full nested tree with ordered exercises and sets | unit | bunx nx run api:test -- --run '(?i)workout_service' |
| TEST-W03-006 | Workout service: working weight snapshot from Exercise.workingWeight | unit | bunx nx run api:test -- --run '(?i)workout_snapshot' |
| TEST-W03-007 | Workout service: input validation (weight >= 0, reps > 0, valid date) | unit | bunx nx run api:test -- --run '(?i)workout_validation' |
| TEST-W03-008 | Workout service: RPE bounds (1.0-10.0, step 0.5) and RIR bounds (0-5) | unit | bunx nx run api:test -- --run '(?i)workout_rpe_rir' |
| TEST-W03-009 | Workout service: cascade delete from DailyLog to exercises/sets/cardio | integration | bunx nx run api:test -- --run '(?i)workout_cascade' |
| TEST-W03-010 | Workout service: DailyLog upsert creates new vs updates existing | integration | bunx nx run api:test -- --run '(?i)workout_upsert' |
| TEST-W03-011 | Workout service: FK constraint violation with invalid exercise ID | integration | bunx nx run api:test -- --run '(?i)workout_fk_constraint' |
| TEST-W03-012 | Workout service: duplicate date violation (unique constraint) | integration | bunx nx run api:test -- --run '(?i)workout_duplicate_date' |
| TEST-W03-013 | Workout GraphQL resolvers: dailyLogByDate query | integration | bunx nx run api:test -- --run '(?i)workout_resolver_query' |
| TEST-W03-014 | Workout GraphQL resolvers: upsertDailyLog mutation | integration | bunx nx run api:test -- --run '(?i)workout_resolver_upsert' |
| TEST-W03-015 | Workout GraphQL resolvers: addWorkoutExercise with snapshot | integration | bunx nx run api:test -- --run '(?i)workout_resolver_add_exercise' |
| TEST-W03-016 | Workout GraphQL resolvers: addWorkoutSet/updateWorkoutSet/removeWorkoutSet | integration | bunx nx run api:test -- --run '(?i)workout_resolver_sets' |
| TEST-W03-017 | Workout GraphQL resolvers: CardioEntry CRUD | integration | bunx nx run api:test -- --run '(?i)workout_resolver_cardio' |
| TEST-W03-018 | Workout GraphQL operations return AuthError without valid PIN session | integration | bunx nx run api:test -- --run '(?i)workout_auth' |
| TEST-W03-019 | Input validation: negative weight, zero reps, invalid date return ValidationError | integration | bunx nx run api:test -- --run '(?i)workout_validation_errors' |
| TEST-W03-020 | Log privacy: notes not appearing in application logs | unit | bunx nx run api:test -- --run '(?i)workout_log_sanitize' |
| TEST-W03-021 | FK constraint: invalid exercise_id returns error | integration | bunx nx run api:test -- --run '(?i)workout_fk_exercise' |
| TEST-W03-022 | DailyLog query by date returns empty (not error) for nonexistent date | integration | bunx nx run api:test -- --run '(?i)workout_empty_date' |
| TEST-W03-023 | Migration smoke test (00082-00085 up + down) | integration | bunx nx run api:test -- --run '(?i)workout_migration' |
| TEST-W03-024 | Go lint for API package | lint | bunx nx run api:lint |
| TEST-W03-025 | GraphQL schema validate | codegen | bunx nx run graphql:validate |
| TEST-W03-026 | Codegen drift check (gqlgen + sqlc) | codegen | bunx nx run api:codegen && bunx nx run graphql:codegen |
| TEST-W03-027 | WAVE-01 admin auth regression tests | unit | bunx nx run api:test -- --run '(?i)admin_auth' |
| TEST-W03-028 | WAVE-02 exercise regression tests | unit | bunx nx run api:test -- --run '(?i)exercise' |

### Fixture Strategy
- Use shared test DB helper from WAVE-01 (test DB migrations)
- Create exercise fixtures for working weight snapshot tests
- Create daily log fixtures with known exercises and sets for resolver tests
- Use WAVE-02 allExercises or exerciseById in service tests (mock or real)
- PIN auth: use WAVE-01 PIN test helpers for auth-gated resolver tests

## Acceptance Criteria Contributions
All 30 ACs covered by verification obligations.

## Exit Criteria Contributions
All 18 ECs defined above with matching verification.

## Verification Contributions
28 verification obligations covering all layers: unit (repos, service), integration (resolvers), lint, codegen, regression.

## Risks And Rollback
- Risk: WAVE-01 PIN test helpers may not exist yet. Mitigation: define test helpers in WAVE-01; WAVE-03 tests assume they exist.
- Risk: WAVE-02 exercise fixtures needed for snapshot tests. Mitigation: create minimal exercise data in test setup.
- Rollback: tests are additive. Old test suite unaffected.

## Questions Raised
- None new.

## Traceability Candidates
- docs/technical-verified/testing-and-delivery.md: TDEC-056 (test data factory), TDEC-057 (weekly workflow e2e)
- docs/prd-wave-details/waves/wave-02.md: TEST-W02 patterns
- docs/prd-wave-details/waves/wave-01.md: TEST-W01 patterns
