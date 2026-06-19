# Wave 03: Workout Diary

## Status
ready-for-dev

## User Approval
pending — awaiting user approval of detailed brief.

## Source Wave Summary
WAVE-03 from docs/prd-waves/waves/wave-03.md. Core workout tracking: daily workouts by date with exercises, sets, working weight snapshots, and cardio. Source status: user-approved (2026-06-18).

## Outcome After Implementation
- OUT-W03-001: Workouts can be logged by date via DailyLog upsert
- OUT-W03-002: Exercises added with snapshot working weight from Exercise library
- OUT-W03-003: Sets with weight, reps, optional RPE/RIR per exercise
- OUT-W03-004: Cardio can be added per day (inline to DailyLog)
- OUT-W03-005: Per-exercise comments stored and retrievable
- OUT-W03-006: Calendar date switching supported via DailyLog-by-date query
- OUT-W03-007: Data available for charts and AI via service layer

## Scope Included
- CAP-W03-001: DailyLog CRUD by date (upsert, unique per user per date)
- CAP-W03-002: WorkoutExercise with order
- CAP-W03-003: WorkoutSet with weight, reps, optional RPE, optional RIR
- CAP-W03-004: Working weight snapshot from Exercise (read at add time)
- CAP-W03-005: CardioEntry inline to DailyLog (type, duration, pulse, zone)
- CAP-W03-006: Comments per exercise (notes field on WorkoutExercise)
- CAP-W03-007: Calendar date switching (query DailyLog by date)

## Scope Excluded
- Exercise CRUD (WAVE-02)
- Charts visualization (WAVE-06)
- AI export (WAVE-07)
- Standalone CardioEntry (WAVE-04 — outside workout context)
- Body weight entries (WAVE-04)
- Workout templates (future scope per PAGE-002 deferral)

## Dependencies And Other-Wave Fit
- WAVE-01 (Foundation): prerequisite — provides PIN auth middleware, common GraphQL types (ValidationError, AuthError, NotFoundError), gqlgen+sqlc codegen infrastructure, config extension pattern, and Goose migration framework. WAVE-03 cannot start until WAVE-01 provides these contracts.
- WAVE-02 (Exercise Library): prerequisite — provides allExercises query for exercise selector, Exercise.workingWeight for snapshot, and exercises table FK target. WAVE-03 cannot start until WAVE-02 provides these.
- WAVE-04 (Cardio/Body): WAVE-03 creates DailyLog-linked CardioEntry (dailyLogId required). WAVE-04 handles standalone cardio and body tracking. No scope collision — ownership boundary documented.
- WAVE-06 (Charts): WAVE-03 provides data model (sets, weights, volumes) consumed by chart aggregation queries.
- WAVE-07/08 (AI Export/Review): WAVE-03 provides exercise comments and set data consumed via service layer.
- WAVE-09 (Backup): WAVE-03 tables are JSON-serializable for export compatibility.

## Frontend Pages Dependencies
- PAGE-002 (Workout Diary): primary frontend consumer. Consumes all WAVE-03 GraphQL operations listed below. Dependency context only; no frontend pages, UI, or UX work in this wave.
- Backend dependencies: dailyLogByDate query, upsertDailyLog mutation, addWorkoutExercise mutation, addWorkoutSet mutation, addCardioEntry mutation.
- WAVE-02 allExercises query is also consumed by PAGE-002 for exercise selector (not provided by WAVE-03).

## Codebase Fit And Touchpoints
- apps/api/internal/repository/postgres/migrations/00082_daily_logs.sql: new migration for daily_logs table
- apps/api/internal/repository/postgres/migrations/00083_workout_exercises.sql: new migration for workout_exercises table
- apps/api/internal/repository/postgres/migrations/00084_workout_sets.sql: new migration for workout_sets table
- apps/api/internal/repository/postgres/migrations/00085_cardio_entries.sql: new migration for cardio_entries table
- apps/api/internal/repository/postgres/queries/daily_logs.sql: sqlc query definitions
- apps/api/internal/repository/postgres/queries/workout_exercises.sql: sqlc query definitions
- apps/api/internal/repository/postgres/queries/workout_sets.sql: sqlc query definitions
- apps/api/internal/repository/postgres/queries/cardio_entries.sql: sqlc query definitions
- apps/api/internal/repository/postgres/daily_log_repo.go: repository adapter for daily logs
- apps/api/internal/repository/postgres/workout_exercise_repo.go: repository adapter for workout exercises
- apps/api/internal/repository/postgres/workout_set_repo.go: repository adapter for workout sets
- apps/api/internal/repository/postgres/cardio_entry_repo.go: repository adapter for cardio entries
- apps/api/internal/service/workout.go: transport-neutral workout service
- libs/graphql/schema/workout.graphql: GraphQL types and operations
- apps/api/internal/graph/workout.resolvers.go: GraphQL resolvers
- apps/api/cmd/server/main.go: wire repos, service, resolvers; register PIN-protected route group
- apps/api/gqlgen.yml: auto-discovers new schema files via glob (no change needed)
- apps/api/sqlc.yaml: auto-discovers new queries via glob (no change needed)

## Design Contracts
- DailyLog upsert: created on first save with content (exercises or cardio). Not created for empty dates (DDEC-W03-001).
- Working weight snapshot: captured at exercise-add time from Exercise.workingWeight. Not updated retroactively (RULE-017).
- Exercise order: WorkoutExercise.display_order determines display order within DailyLog (DDEC-W03-002).
- Set order: WorkoutSet.set_number determines order within WorkoutExercise, starts at 1, sequential (DDEC-W03-003).
- Cascade delete: DailyLog -> WorkoutExercise -> WorkoutSet and DailyLog -> CardioEntry. FK ON DELETE CASCADE (DDEC-W03-004).
- Exercise FK: workout_exercises.exercise_id -> exercises(id) ON DELETE NO ACTION. Preserves workout history on exercise soft-delete (DDEC-W03-005).
- Concurrent edit: last-write-wins for MVP. No optimistic locking (DDEC-W03-006, DQ-W03-001 deferred).
- Backdating: no lower bound on date. Any valid DATE accepted (DDEC-W03-007).
- PIN auth: WAVE-01 middleware guards all WAVE-03 GraphQL endpoints. AuthError returned when session missing or invalid.
- Error format: { "error": { "code": "ERROR_CODE", "message": "Human readable" } } per TDEC-027.
- Error codes: VALIDATION_ERROR, NOT_FOUND, AUTH_ERROR, INTERNAL_ERROR, DUPLICATE_DATE.
- Set weight: REAL >= 0 (0 allowed for bodyweight exercises).
- Set reps: INT > 0 (positive integer).
- RPE: optional, REAL 1.0-10.0, step 0.5.
- RIR: optional, INT 0-5.

## Data API Integration And Operations

### PostgreSQL Tables

**daily_logs:**
Columns: id (UUID PK), user_id (UUID FK -> default_user), date (DATE), notes (TEXT nullable), body_weight (REAL nullable), created_at (TIMESTAMPTZ), updated_at (TIMESTAMPTZ).
Constraints: UNIQUE(user_id, date).
Indexes: idx_daily_logs_user_date (user_id, date).

**workout_exercises:**
Columns: id (UUID PK), user_id (UUID FK -> default_user), daily_log_id (UUID FK -> daily_logs, ON DELETE CASCADE), exercise_id (UUID FK -> exercises, ON DELETE NO ACTION), display_order (INT), working_weight_snapshot (REAL nullable), notes (TEXT nullable), created_at (TIMESTAMPTZ), updated_at (TIMESTAMPTZ).
Indexes: idx_workout_exercises_daily_log (daily_log_id), idx_workout_exercises_exercise (exercise_id).

**workout_sets:**
Columns: id (UUID PK), workout_exercise_id (UUID FK -> workout_exercises, ON DELETE CASCADE), set_number (INT), weight (REAL), reps (INT), rpe (REAL nullable), rir (INT nullable), notes (TEXT nullable), created_at (TIMESTAMPTZ), updated_at (TIMESTAMPTZ).
Indexes: idx_workout_sets_exercise (workout_exercise_id).

**cardio_entries:**
Columns: id (UUID PK), user_id (UUID FK -> default_user), daily_log_id (UUID FK -> daily_logs, ON DELETE CASCADE), cardio_type (VARCHAR 32), duration_minutes (INT), avg_pulse (INT nullable), heart_rate_zone (INT nullable), notes (TEXT nullable), created_at (TIMESTAMPTZ), updated_at (TIMESTAMPTZ).
Indexes: idx_cardio_entries_daily_log (daily_log_id).

### GraphQL Operations (via single /graphql endpoint)

Queries:
- dailyLogByDate(date: Date!): DailyLogResult!
- dailyLogsByDateRange(startDate: Date!, endDate: Date!): [DailyLog!]!

Mutations:
- upsertDailyLog(input: UpsertDailyLogInput!): UpsertDailyLogResult!
- addWorkoutExercise(dailyLogId: UUID!, exerciseId: UUID!): AddWorkoutExerciseResult!
- updateWorkoutExercise(id: UUID!, input: UpdateWorkoutExerciseInput!): UpdateWorkoutExerciseResult!
- removeWorkoutExercise(id: UUID!): RemoveWorkoutExerciseResult!
- addWorkoutSet(workoutExerciseId: UUID!, input: AddWorkoutSetInput!): AddWorkoutSetResult!
- updateWorkoutSet(id: UUID!, input: UpdateWorkoutSetInput!): UpdateWorkoutSetResult!
- removeWorkoutSet(id: UUID!): RemoveWorkoutSetResult!
- addCardioEntry(dailyLogId: UUID!, input: AddCardioEntryInput!): AddCardioEntryResult!
- updateCardioEntry(id: UUID!, input: UpdateCardioEntryInput!): UpdateCardioEntryResult!
- removeCardioEntry(id: UUID!): RemoveCardioEntryResult!
- deleteDailyLog(id: UUID!): DeleteDailyLogResult!

### Observability
- Log markers: [DailyLog][create|update|delete|get], [WorkoutExercise][create|update|delete], [WorkoutSet][create|update|delete], [CardioEntry][create|update|delete]
- No sensitive content logged: notes, comments not included in log messages
- Error codes: VALIDATION_ERROR, NOT_FOUND, AUTH_ERROR, INTERNAL_ERROR, DUPLICATE_DATE

### Rollout and Rollback
- Rollout: merge PR, CI builds and runs tests, deploy via Dokploy compose update. Migrations 00082-00085 run at startup.
- Rollback: revert PR, CI builds previous image, Dokploy compose update rolls back. Down migrations available for 00082-00085.
- Compatibility: all new GraphQL operations are additive. No existing API changes. WAVE-01/WAVE-02 endpoints unchanged.
- WAVE-04 compatibility: cardio_entries table exists. WAVE-03 only creates DailyLog-linked entries. WAVE-04 will coordinate for standalone entries.

## Security Privacy And Compliance
- All endpoints protected by WAVE-01 PIN auth middleware (GraphQL queries and mutations)
- When PIN is disabled, endpoints accessible without auth (consistent with TDEC-037)
- Single-user scope: all data owned by default user; service layer scopes queries by user_id
- No sensitive content logged: notes, RPE/RIR values, cardio zones not logged
- Audit events: DailyLog CRUD, WorkoutExercise CRUD, WorkoutSet CRUD, CardioEntry CRUD
- Input validation enforced server-side: weight >= 0, reps > 0, RPE 1.0-10.0, RIR 0-5, date valid
- String field length limits: notes max 1000 characters
- No file uploads or binary data in this wave
- FK constraints enforce referential integrity for all relationships

## Implementation Slices

| Slice ID | Name | Description |
| --- | --- | --- |
| SLICE-W03-001 | DB migration: daily_logs | Create goose migration 00082_daily_logs.sql. Columns: id, user_id, date (UNIQUE), notes, body_weight, timestamps. Index on (user_id, date). |
| SLICE-W03-002 | DB migration: workout_exercises | Create goose migration 00083_workout_exercises.sql. Columns: id, user_id, daily_log_id FK CASCADE, exercise_id FK NO ACTION, display_order, working_weight_snapshot, notes, timestamps. Indexes on daily_log_id and exercise_id. |
| SLICE-W03-003 | DB migration: workout_sets | Create goose migration 00084_workout_sets.sql. Columns: id, workout_exercise_id FK CASCADE, set_number, weight, reps, rpe, rir, notes, timestamps. Index on workout_exercise_id. |
| SLICE-W03-004 | DB migration: cardio_entries | Create goose migration 00085_cardio_entries.sql. Columns: id, user_id, daily_log_id FK CASCADE, cardio_type, duration_minutes, avg_pulse, heart_rate_zone, notes, timestamps. Index on daily_log_id. |
| SLICE-W03-005 | sqlc queries: daily_logs | CRUD queries: get by user_id+date, upsert, delete by id, list by date range. |
| SLICE-W03-006 | sqlc queries: workout_exercises | CRUD queries: list by daily_log_id ordered by display_order, create, update order/notes, delete (cascading to sets). |
| SLICE-W03-007 | sqlc queries: workout_sets | CRUD queries: list by workout_exercise_id ordered by set_number, create (auto set_number = max+1), update, delete. |
| SLICE-W03-008 | sqlc queries: cardio_entries | CRUD queries: list by daily_log_id, create, update, delete. |
| SLICE-W03-009 | DailyLog repository | Repository adapter for daily_logs with GetByDate(userID, date), Upsert, Delete. |
| SLICE-W03-010 | WorkoutExercise repository | Repository adapter for workout_exercises with ListByDailyLog, Create, Update, Delete. |
| SLICE-W03-011 | WorkoutSet repository | Repository adapter for workout_sets with ListByWorkoutExercise, Create, Update, Delete. |
| SLICE-W03-012 | CardioEntry repository | Repository adapter for cardio_entries with ListByDailyLog, Create, Update, Delete. |
| SLICE-W03-013 | Workout service | Transport-neutral service: DailyLogByDate, UpsertDailyLog, AddExercise (snapshot weight), RemoveExercise, UpdateExercise, AddSet, UpdateSet, RemoveSet, AddCardio, RemoveCardio, DeleteDailyLog. Validates inputs. Calls WAVE-02 for snapshot. |
| SLICE-W03-014 | GraphQL schema | Add workout.graphql with DailyLog, WorkoutExercise, WorkoutSet, CardioEntry types, inputs, and union result types. Extend root Query and Mutation. |
| SLICE-W03-015 | GraphQL resolvers + wiring | Implement resolvers for all operations. Wire repos and service in main.go. Register PIN-protected fitness GraphQL endpoint. |

## Acceptance Criteria

| AC ID | Description |
| --- | --- |
| AC-W03-001 | DailyLog can be queried by date (YYYY-MM-DD). Returns existing log or empty result. |
| AC-W03-002 | DailyLog is created/updated via upsert mutation. Returns the full daily log with exercises, sets, and cardio. |
| AC-W03-003 | DailyLog upsert creates a new record when no record exists for the date and exercises/cardio are provided. |
| AC-W03-004 | DailyLog upsert updates an existing record when one already exists for the date. |
| AC-W03-005 | WorkoutExercise can be added to a DailyLog with exerciseId, order, and workingWeightSnapshot. |
| AC-W03-006 | WorkoutExercise workingWeightSnapshot is auto-populated from Exercise.workingWeight at add time. |
| AC-W03-007 | Multiple WorkoutExercises can be added to a single DailyLog, each with unique order values. |
| AC-W03-008 | WorkoutSet can be added to a WorkoutExercise with weight, reps, setNumber. |
| AC-W03-009 | WorkoutSet can optionally include RPE (1-10 scale). |
| AC-W03-010 | WorkoutSet can optionally include RIR (0-5 scale). |
| AC-W03-011 | WorkoutSet can optionally include notes. |
| AC-W03-012 | Multiple WorkoutSets can be added to a single WorkoutExercise with sequential setNumber. |
| AC-W03-013 | WorkoutExercise.notes (comment) can be stored and retrieved in queries. |
| AC-W03-014 | WorkoutExercise order can be updated (reordering). |
| AC-W03-015 | CardioEntry can be created within a DailyLog with cardioType, durationMinutes. |
| AC-W03-016 | CardioEntry optionally includes avgPulse and heartRateZone. |
| AC-W03-017 | WorkoutExercise can be removed from a DailyLog (cascading delete of its WorkoutSets). |
| AC-W03-018 | WorkoutSet can be updated (weight, reps, RPE, RIR, notes). |
| AC-W03-019 | WorkoutSet can be removed from a WorkoutExercise. |
| AC-W03-020 | CardioEntry can be removed from a DailyLog. |
| AC-W03-021 | DailyLog can be deleted entirely (cascades to WorkoutExercises, WorkoutSets, CardioEntries). |
| AC-W03-022 | Set weight must be >= 0 (0 allowed for warmup/bodyweight). Mutation returns ValidationError when negative. |
| AC-W03-023 | Set reps must be > 0 (positive integer). Mutation returns ValidationError when zero or negative. |
| AC-W03-024 | WorkoutExercise.exerciseId must reference an existing active Exercise. Mutation returns NotFoundError when invalid. |
| AC-W03-025 | DailyLog query by date returns empty result (not error) when no record exists for that date. |
| AC-W03-026 | DailyLog returns all WorkoutExercises ordered by WorkoutExercise.order ASC. |
| AC-W03-027 | WorkoutExercise returns all WorkoutSets ordered by WorkoutSet.setNumber ASC. |
| AC-W03-028 | DailyLog date field is validated as YYYY-MM-DD format. Mutation returns ValidationError for invalid dates. |
| AC-W03-029 | DailyLog, WorkoutExercise, WorkoutSet, CardioEntry mutations return AuthError when PIN session header is missing or invalid. |
| AC-W03-030 | DailyLog query by date returns AuthError when PIN session header is missing or invalid. |

## Exit Criteria

| EC ID | Description |
| --- | --- |
| EC-W03-001 | AC-W03-001 through AC-W03-030 pass via focused tests. |
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

## Verification Obligations

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

## Rollout Rollback And Compatibility
- Rollout: merge PR, CI builds and runs tests, deploy via Dokploy compose update. New services start alongside existing WAVE-01/WAVE-02 infrastructure.
- Rollback: revert PR, CI builds previous image, Dokploy compose update rolls back. Existing WAVE-01 and WAVE-02 endpoints unchanged.
- Compatibility: all new GraphQL operations are additive. No existing API changes. WAVE-01 health/admin endpoints, WAVE-02 exercise endpoints unchanged.
- Migration: goose migrations 00082-00085 run at startup. Down migrations available for rollback.
- WAVE-04 compatibility: cardio_entries table created. WAVE-03 creates only DailyLog-linked entries. WAVE-04 coordinates for standalone entries.

## Handoff Packets
- HANDOFF-W03-001: This wave brief document
- HANDOFF-W03-002: Planner reports (6 scopes, 1 cycle)
- HANDOFF-W03-003: Reviewer evidence (7 perspectives, 1 cycle)
- HANDOFF-W03-004: Final fit review evidence

## Reviewer Verdicts

| Wave | Perspective | Attempt | Verdict | Reviewer Report | Required Revisions | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| WAVE-03 | product-scope-and-ac | 1 | approved | review-product-scope-and-ac-attempt-1.md | none | 30 ACs cover all scope, edge cases documented |
| WAVE-03 | architecture-codebase-fit | 1 | approved | review-architecture-codebase-fit-attempt-1.md | none | 15 slices follow existing patterns |
| WAVE-03 | data-api-integration-ops | 1 | approved | review-data-api-integration-ops-attempt-1.md | none | Data/API/ops coverage adequate |
| WAVE-03 | security-privacy-compliance | 1 | approved | review-security-privacy-compliance-attempt-1.md | none | PIN auth, input validation, log privacy covered |
| WAVE-03 | testing-exit-criteria | 1 | approved | review-testing-exit-criteria-attempt-1.md | none | 28 test obligations cover all AC and EC |
| WAVE-03 | sequencing-other-wave-fit | 1 | approved | review-sequencing-other-wave-fit-attempt-1.md | none | Dependency order correct, no collision |
| WAVE-03 | traceability-consistency | 1 | approved | review-traceability-consistency-attempt-1.md | none | Source traceability documented |
| WAVE-03 | final-wave-fit-review | 1 | approved | final-wave-fit-review-attempt-1.md | none | Package is ready-for-dev |

## Open Questions

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Q-WORKOUT-001 | WAVE-03 | operations | needs-owner-decision | None | Concurrent edit handling? Two browser tabs could lead to last-write-wins data loss. | Data integrity for workout diary entries. | Confirm strategy: optimistic locking (version field), last-write-wins, or deferred (acknowledge risk for MVP). | docs/prd-waves/open-questions.md, planner-data-integration-ops.md | resolved | last-write-wins for MVP. No optimistic locking. Tracked as DQ-W03-001 (deferred post-MVP). |
| DQ-W03-001 | WAVE-03 | operations | deferred | EDGE-016 | Should DailyLog updates use optimistic concurrency control (version field)? | Prevents last-write-wins data loss across browser tabs. | Deferred to post-MVP or owner decision. MVP uses last-write-wins. | planner-data-integration-ops-attempt-1.md | deferred | Deferred to post-MVP. MVP uses last-write-wins. |

## Traceability
- docs/prd-waves/waves/wave-03.md: source wave boundary, outcomes, capability groups
- docs/product-verified/domain-model.md: DailyLog, WorkoutExercise, WorkoutSet, CardioEntry entities
- docs/product-verified/acceptance-criteria.md: AC-005 through AC-011, AC-035 through AC-042
- docs/product-verified/user-flows.md: Enter Workout For Today, Enter Workout Backdated
- docs/technical-verified/api-contracts.md: hybrid GraphQL/REST protocol, TDEC-001
- docs/technical-verified/data-contracts.md: entity contracts, TDEC-020, TDEC-021, TDEC-022
- docs/technical-verified/implementation-slices.md: Slice 2 DailyLog mapping
- docs/technical-verified/auth-security-compliance.md: PIN auth, TDEC-037
- docs/technical-verified/operations-observability.md: log markers, error format
- docs/development-plan.xml: M-API, M-PRD-WAVE-DETAILER module contracts
- docs/knowledge-graph.xml: existing module boundaries
- docs/prd-wave-details/waves/wave-01.md: WAVE-01 dependency contracts
- docs/prd-wave-details/waves/wave-02.md: WAVE-02 allExercises contract
- docs/prd-waves/frontend-pages/page-002.md: workout diary backend dependencies
- apps/api/internal: codebase patterns for service/repository/middleware/handler structure
