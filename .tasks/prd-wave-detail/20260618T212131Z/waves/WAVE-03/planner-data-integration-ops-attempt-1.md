# WAVE-03 data-integration-ops Planner Attempt 1

## Sources Read
- docs/technical-verified/data-contracts.md
- docs/technical-verified/api-contracts.md
- docs/technical-verified/operations-observability.md
- docs/product-verified/domain-model.md
- docs/product-verified/edge-cases.md
- docs/prd-waves/frontend-pages/page-002.md
- apps/api/internal/repository/postgres/migrations/00079_admin_users.sql
- apps/api/sqlc.yaml

## Selected Backend Wave Boundary
WAVE-03 manages the complete data lifecycle for DailyLog, WorkoutExercise, WorkoutSet, and CardioEntry entities via GraphQL CRUD operations backed by PostgreSQL. No REST endpoints needed (binary uploads not in scope).

## Neighboring Backend Wave Fit
- WAVE-01: Provides Redis-backed PIN sessions, config for DB connections, migration infrastructure.
- WAVE-02: Provides exercises table FK target for workout_exercises.exercise_id. allExercises query for working weight snapshot.
- WAVE-04: CardioEntry entity shared. WAVE-03 creates only DailyLog-linked entries.

## Frontend Pages Context
- PAGE-002: consumes GraphQL operations for daily log by date, exercise/set/cardio CRUD. Backend dependency only.

## Proposed Details

### Data Lifecycle

#### DailyLog
- **Create**: UPSERT operation — GET to check existence, then INSERT or UPDATE
- **Read**: query by userId + date (composite unique). Returns full nested tree: exercises -> sets, cardio
- **Update**: upsert replaces the entire day's exercise/set/cardio content per save operation
- **Delete**: removes daily_log row, cascades to workout_exercises -> workout_sets and cardio_entries
- **Unique constraint**: (user_id, date) prevents duplicate entries per day
- **bodyWeight**: optional REAL field on DailyLog (not a separate table in WAVE-03)

#### WorkoutExercise
- **Create**: INSERT with daily_log_id FK, exercise_id FK, order (auto-assigned), working_weight_snapshot from Exercise.workingWeight
- **Read**: list by daily_log_id ORDER BY order ASC, eager-loaded with WorkoutSets
- **Update**: order (reordering), notes, working_weight_snapshot (not updated retroactively per RULE-017)
- **Delete**: CASCADE to workout_sets

#### WorkoutSet
- **Create**: INSERT with workout_exercise_id FK, auto-incrementing setNumber (max existing + 1), weight, reps, optional rpe/rir/notes
- **Read**: list by workout_exercise_id ORDER BY set_number ASC
- **Update**: weight, reps, rpe, rir, notes
- **Delete**: remove individual set; remaining sets keep their set_number (no renumbering)

#### CardioEntry
- **Create**: INSERT with daily_log_id FK (required), cardioType, durationMinutes, optional avgPulse/heartRateZone/notes
- **Read**: list by daily_log_id
- **Update**: all fields except daily_log_id
- **Delete**: remove individual entry

### API Endpoints (GraphQL via single /graphql endpoint)

All operations are GraphQL mutations/queries. No REST endpoints in WAVE-03.

**Queries:**
- `dailyLogByDate(date: Date!): DailyLogResult!` — returns DailyLog with full nested tree or empty
- `dailyLogsByDateRange(startDate: Date!, endDate: Date!): [DailyLog!]!` — list for date range (calendar navigation)

**Mutations:**
- `upsertDailyLog(input: UpsertDailyLogInput!): UpsertDailyLogResult!` — save entire day's data
- `addWorkoutExercise(dailyLogId: UUID!, exerciseId: UUID!): AddWorkoutExerciseResult!` — add exercise with snapshot
- `updateWorkoutExercise(id: UUID!, input: UpdateWorkoutExerciseInput!): UpdateWorkoutExerciseResult!` — update order/notes
- `removeWorkoutExercise(id: UUID!): RemoveWorkoutExerciseResult!` — remove exercise and its sets
- `addWorkoutSet(workoutExerciseId: UUID!, input: AddWorkoutSetInput!): AddWorkoutSetResult!` — add set
- `updateWorkoutSet(id: UUID!, input: UpdateWorkoutSetInput!): UpdateWorkoutSetResult!` — update set
- `removeWorkoutSet(id: UUID!): RemoveWorkoutSetResult!` — remove set
- `addCardioEntry(dailyLogId: UUID!, input: AddCardioEntryInput!): AddCardioEntryResult!` — add cardio
- `updateCardioEntry(id: UUID!, input: UpdateCardioEntryInput!): UpdateCardioEntryResult!` — update cardio
- `removeCardioEntry(id: UUID!): RemoveCardioEntryResult!` — remove cardio
- `deleteDailyLog(id: UUID!): DeleteDailyLogResult!` — delete entire day

### Database Schema

**daily_logs table:**
```sql
CREATE TABLE daily_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES default_user(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    notes TEXT,
    body_weight REAL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, date)
);
CREATE INDEX idx_daily_logs_user_date ON daily_logs (user_id, date);
```

**workout_exercises table:**
```sql
CREATE TABLE workout_exercises (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES default_user(id) ON DELETE CASCADE,
    daily_log_id UUID NOT NULL REFERENCES daily_logs(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE NO ACTION,
    display_order INT NOT NULL,
    working_weight_snapshot REAL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_workout_exercises_daily_log ON workout_exercises (daily_log_id);
CREATE INDEX idx_workout_exercises_exercise ON workout_exercises (exercise_id);
```

**workout_sets table:**
```sql
CREATE TABLE workout_sets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workout_exercise_id UUID NOT NULL REFERENCES workout_exercises(id) ON DELETE CASCADE,
    set_number INT NOT NULL,
    weight REAL NOT NULL,
    reps INT NOT NULL,
    rpe REAL,
    rir INT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_workout_sets_exercise ON workout_sets (workout_exercise_id);
```

**cardio_entries table:**
```sql
CREATE TABLE cardio_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES default_user(id) ON DELETE CASCADE,
    daily_log_id UUID NOT NULL REFERENCES daily_logs(id) ON DELETE CASCADE,
    cardio_type VARCHAR(32) NOT NULL,
    duration_minutes INT NOT NULL,
    avg_pulse INT,
    heart_rate_zone INT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_cardio_entries_daily_log ON cardio_entries (daily_log_id);
```

### Observability
- Log markers per TDEC-053 pattern: [DailyLog][create|update|delete|get], [WorkoutExercise][create|update|delete], [WorkoutSet][create|update|delete], [CardioEntry][create|update|delete]
- Error format per TDEC-027: { "error": { "code": "ERROR_CODE", "message": "..." } }
- Error codes: VALIDATION_ERROR, NOT_FOUND, AUTH_ERROR, INTERNAL_ERROR, DUPLICATE_DATE
- No sensitive content logged: workout notes, cardio notes, set notes not included in log messages

### Rollout and Rollback
- Rollout: merge PR -> CI builds -> Dokploy compose update. Migrations 00082-00085 run at startup.
- Rollback: revert PR -> CI rebuilds previous image -> Dokploy compose update. Down migrations remove tables.
- Compatibility: all WAVE-03 operations are additive. Existing WAVE-01 health/admin endpoints and WAVE-02 exercise endpoints unchanged.
- WAVE-04 compatibility: cardio_entries table exists but WAVE-03 only creates daily_log_id-linked entries. WAVE-04 will extend if needed.

## Acceptance Criteria Contributions
All 30 ACs supported by the data/API contracts described above.

## Exit Criteria Contributions
- EC-W03-004: Migrations apply and roll back cleanly
- EC-W03-005: DailyLog upsert round-trip with full nested data
- EC-W03-006: CardioEntry CRUD within DailyLog

## Verification Contributions
- TEST-W03-001 through TEST-W03-012: repository tests
- TEST-W03-013 through TEST-W03-022: resolver integration tests

## Risks And Rollback
- Concurrent edit risk: Q-WORKOUT-001, deferred. Last-write-wins for MVP.
- exercise_id FK with NO ACTION: prevents deletion of exercises referenced in workout history. Compatible with WAVE-02 soft delete.
- Large daily logs: single-day data bounded by practical exercise count (typically < 20 exercises, < 50 sets).

## Questions Raised
- DQ-W03-006: Cascade delete behavior confirmed CASCADE throughout.
- DQ-W03-001: Concurrent edit handling deferred.

## Traceability Candidates
- docs/technical-verified/data-contracts.md: TDEC-020, TDEC-021 (index strategy), TDEC-022 (enum types)
- docs/technical-verified/api-contracts.md: TDEC-001 (hybrid model), TDEC-026 (endpoint catalog), TDEC-027 (error format)
- docs/technical-verified/operations-observability.md: log markers, error format
