# Planner Report: Architecture & Codebase — WAVE-07 (AI Export and Prompt Builder)

**Run**: `20260621T170113Z` | **Wave**: WAVE-07 | **Role**: architecture-codebase | **Attempt**: 1

---

## 1. Source-Backed Implementation Slices

### SLICE-W07-001: UserProfile Migration

**Files**: `.../migrations/00093_user_profiles.sql`

New table `user_profiles` (not `atlas_user_profiles` — using singular-convention like `week_flags`, not `atlas_` prefix):

```sql
CREATE TABLE user_profiles (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                     UUID NOT NULL REFERENCES atlas_users(id) UNIQUE,
    goal                        TEXT,
    height                      REAL,
    birth_date                  DATE,
    training_experience         TEXT,
    current_training_split      TEXT,
    preferred_progression_style TEXT,
    nutrition_strategy          TEXT,
    persistent_ai_context       TEXT,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

- `UNIQUE(user_id)` — one profile per user, created on bootstrap alongside settings.
- All fields nullable except `id`, `user_id`, `created_at`, `updated_at`.
- Why not extend `atlas_users`? The existing table only has `id, display_name`. Adding 9 optional fields to the user identity table conflates identity with profile. A separate table keeps migration small and follows the domain-separation pattern already used (settings are a separate table from users).

**Backfill**: Update `atlasBootstrapService.EnsureDefaultUser` to also create an empty user_profile row (or return nil gracefully when queried). This keeps the profile optional — users who never open the export feature don't need a row.

**Down**: `DROP TABLE IF EXISTS user_profiles;`

### SLICE-W07-002: UserProfile Model

**File**: `apps/api/internal/atlas/models/user_profile.go`

Follow the exact pattern from `models/settings.go`:

```go
type UserProfileRecord struct {
    ID                       string
    UserID                   string
    Goal                     *string
    Height                   *float32
    BirthDate                *Date
    TrainingExperience       *string
    CurrentTrainingSplit     *string
    PreferredProgressionStyle *string
    NutritionStrategy        *string
    PersistentAiContext      *string
    CreatedAt                string
    UpdatedAt                string
}

type UserProfile struct {
    ID                       string   `json:"id"`
    UserID                   string   `json:"userId"`
    Goal                     *string  `json:"goal"`
    Height                   *float32 `json:"height"`
    BirthDate                *Date    `json:"birthDate"`
    TrainingExperience       *string  `json:"trainingExperience"`
    CurrentTrainingSplit     *string  `json:"currentTrainingSplit"`
    PreferredProgressionStyle *string `json:"preferredProgressionStyle"`
    NutritionStrategy        *string  `json:"nutritionStrategy"`
    PersistentAiContext      *string  `json:"persistentAiContext"`
    CreatedAt                string   `json:"createdAt"`
    UpdatedAt                string   `json:"updatedAt"`
}

type UserProfileInput struct {
    Goal                     *string  `json:"goal"`
    Height                   *float32 `json:"height"`
    BirthDate                *Date    `json:"birthDate"`
    TrainingExperience       *string  `json:"trainingExperience"`
    CurrentTrainingSplit     *string  `json:"currentTrainingSplit"`
    PreferredProgressionStyle *string `json:"preferredProgressionStyle"`
    NutritionStrategy        *string  `json:"nutritionStrategy"`
    PersistentAiContext      *string  `json:"persistentAiContext"`
}
```

Result/error types follow the WeekFlag/Cardio pattern: `UserProfileResult`, `UserProfileNotFoundErr`, `UserProfileValidationErr`, `UserProfileAuthErr`, `UserProfileErrorCode`.

### SLICE-W07-003: UserProfile SQLc Queries

**File**: `apps/api/internal/repository/postgres/queries/user_profiles.sql`

```sql
-- name: GetUserProfileByUserID :one
SELECT id, user_id, goal, height, birth_date, training_experience, current_training_split,
       preferred_progression_style, nutrition_strategy, persistent_ai_context,
       created_at, updated_at
FROM user_profiles
WHERE user_id = $1
LIMIT 1;

-- name: UpsertUserProfile :one
INSERT INTO user_profiles (user_id, goal, height, birth_date, training_experience,
                           current_training_split, preferred_progression_style,
                           nutrition_strategy, persistent_ai_context)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (user_id)
DO UPDATE SET goal = COALESCE($2, user_profiles.goal),
              height = COALESCE($3, user_profiles.height),
              birth_date = COALESCE($4, user_profiles.birth_date),
              training_experience = COALESCE($5, user_profiles.training_experience),
              current_training_split = COALESCE($6, user_profiles.current_training_split),
              preferred_progression_style = COALESCE($7, user_profiles.preferred_progression_style),
              nutrition_strategy = COALESCE($8, user_profiles.nutrition_strategy),
              persistent_ai_context = COALESCE($9, user_profiles.persistent_ai_context),
              updated_at = now()
RETURNING id, user_id, goal, height, birth_date, training_experience, current_training_split,
          preferred_progression_style, nutrition_strategy, persistent_ai_context,
          created_at, updated_at;

-- name: CreateUserProfile :one
INSERT INTO user_profiles (user_id, goal, height, birth_date, training_experience,
                           current_training_split, preferred_progression_style,
                           nutrition_strategy, persistent_ai_context)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, user_id, goal, height, birth_date, training_experience, current_training_split,
          preferred_progression_style, nutrition_strategy, persistent_ai_context,
          created_at, updated_at;
```

After writing, run `bun run codegen` to regenerate sqlc.

### SLICE-W07-004: UserProfile Repository

**File**: `apps/api/internal/atlas/repository/postgres/user_profile_repo.go`

Interface:
```go
type UserProfileRepository interface {
    FindByUserID(ctx context.Context, userID string) (*models.UserProfileRecord, error)
    Upsert(ctx context.Context, userID string, input models.UserProfileInput) (*models.UserProfileRecord, error)
}
```

Implementation uses `*generated.Queries`, pgtype UUID conversion, nullable helper functions from `settings_repo.go`.

### SLICE-W07-005: UserProfile Service

**File**: `apps/api/internal/atlas/service/user_profile_service.go`

Interface:
```go
type UserProfileService interface {
    Get(ctx context.Context, userID string) (*models.UserProfile, error)
    Update(ctx context.Context, userID string, input models.UserProfileInput) (*models.UserProfile, error)
}
```

- `Get` returns `nil, nil` when no profile exists (graceful for users not yet configured).
- `Update` delegates to `Upsert` — no partial-update logic needed since every field is optional in the input (COALESCE handles it in SQL).

### SLICE-W07-006: UserProfile Resolver + GraphQL Schema

**Files**:
- `apps/api/internal/atlas/graph/schema/user_profile.graphql`
- `apps/api/internal/atlas/graph/resolver/user_profile.go`

Schema:
```graphql
type UserProfile {
  id: ID!
  userId: ID!
  goal: String
  height: Float
  birthDate: Date
  trainingExperience: String
  currentTrainingSplit: String
  preferredProgressionStyle: String
  nutritionStrategy: String
  persistentAiContext: String
  createdAt: Time!
  updatedAt: Time!
}

input UserProfileInput {
  goal: String
  height: Float
  birthDate: Date
  trainingExperience: String
  currentTrainingSplit: String
  preferredProgressionStyle: String
  nutritionStrategy: String
  persistentAiContext: String
}

type UserProfileResult {
  userProfile: UserProfile
  validationError: UserProfileValidationError
  notFoundError: UserProfileNotFoundError
  authError: UserProfileAuthError
}

type UserProfileValidationError {
  message: String!
  code: UserProfileErrorCode!
}

type UserProfileNotFoundError {
  message: String!
  code: UserProfileErrorCode!
}

type UserProfileAuthError {
  message: String!
  code: UserProfileErrorCode!
}

enum UserProfileErrorCode {
  VALIDATION_ERROR
  NOT_FOUND
  AUTH_ERROR
  INTERNAL_ERROR
}
```

Resolver methods: `GetUserProfile` (Query), `UpdateUserProfile` (Mutation). Pattern = auth check from middleware, delegate to service, map errors to result types.

### SLICE-W07-007: AiExport Migration

**File**: `.../migrations/00094_ai_exports.sql`

```sql
CREATE TABLE ai_exports (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           UUID NOT NULL REFERENCES atlas_users(id),
    date_range_start  DATE NOT NULL,
    date_range_end    DATE NOT NULL,
    include_photos    BOOLEAN NOT NULL DEFAULT true,
    include_nutrition BOOLEAN NOT NULL DEFAULT true,
    include_cardio    BOOLEAN NOT NULL DEFAULT true,
    include_measurements BOOLEAN NOT NULL DEFAULT true,
    user_comment      TEXT,
    generated_prompt  TEXT NOT NULL,
    export_file_path  TEXT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

- `export_file_path` is NULL until ZIP generation completes.
- No `UNIQUE` constraints except PK — users can have multiple exports.

**Down**: `DROP TABLE IF EXISTS ai_exports;`

### SLICE-W07-008: AiExport Model

**File**: `apps/api/internal/atlas/models/ai_export.go`

Record + Public + Input types:

```go
type AiExportRecord struct {
    ID                  string
    UserID              string
    DateRangeStart      Date
    DateRangeEnd        Date
    IncludePhotos       bool
    IncludeNutrition    bool
    IncludeCardio       bool
    IncludeMeasurements bool
    UserComment         *string
    GeneratedPrompt     string
    ExportFilePath      *string
    CreatedAt           string
    UpdatedAt           string
}

type AiExport struct {
    ID                  string  `json:"id"`
    UserID              string  `json:"userId"`
    DateRangeStart      Date    `json:"dateRangeStart"`
    DateRangeEnd        Date    `json:"dateRangeEnd"`
    IncludePhotos       bool    `json:"includePhotos"`
    IncludeNutrition    bool    `json:"includeNutrition"`
    IncludeCardio       bool    `json:"includeCardio"`
    IncludeMeasurements bool    `json:"includeMeasurements"`
    UserComment         *string `json:"userComment"`
    GeneratedPrompt     string  `json:"generatedPrompt"`
    ExportFilePath      *string `json:"exportFilePath"`
    CreatedAt           string  `json:"createdAt"`
    UpdatedAt           string  `json:"updatedAt"`
}

type CreateAiExportInput struct {
    DateRangeStart      Date    `json:"dateRangeStart"`
    DateRangeEnd        Date    `json:"dateRangeEnd"`
    IncludePhotos       bool    `json:"includePhotos"`
    IncludeNutrition    bool    `json:"includeNutrition"`
    IncludeCardio       bool    `json:"includeCardio"`
    IncludeMeasurements bool    `json:"includeMeasurements"`
    UserComment         *string `json:"userComment"`
}
```

Result/error types follow the standard pattern.

### SLICE-W07-009: AiExport SQLc Queries

**File**: `apps/api/internal/repository/postgres/queries/ai_exports.sql`

```sql
-- name: CreateAiExport :one
INSERT INTO ai_exports (user_id, date_range_start, date_range_end, include_photos,
                        include_nutrition, include_cardio, include_measurements,
                        user_comment, generated_prompt)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, user_id, date_range_start, date_range_end, include_photos,
          include_nutrition, include_cardio, include_measurements,
          user_comment, generated_prompt, export_file_path,
          created_at, updated_at;

-- name: GetAiExportByID :one
SELECT id, user_id, date_range_start, date_range_end, include_photos,
       include_nutrition, include_cardio, include_measurements,
       user_comment, generated_prompt, export_file_path,
       created_at, updated_at
FROM ai_exports
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: ListAiExportsByUserID :many
SELECT id, user_id, date_range_start, date_range_end, include_photos,
       include_nutrition, include_cardio, include_measurements,
       user_comment, generated_prompt, export_file_path,
       created_at, updated_at
FROM ai_exports
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateAiExportFilePath :one
UPDATE ai_exports
SET export_file_path = $3, updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, date_range_start, date_range_end, include_photos,
          include_nutrition, include_cardio, include_measurements,
          user_comment, generated_prompt, export_file_path,
          created_at, updated_at;

-- name: DeleteAiExport :one
DELETE FROM ai_exports
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, date_range_start, date_range_end, include_photos,
          include_nutrition, include_cardio, include_measurements,
          user_comment, generated_prompt, export_file_path,
          created_at, updated_at;
```

### SLICE-W07-010: AiExport Repository

**File**: `apps/api/internal/atlas/repository/postgres/ai_export_repo.go`

Interface:
```go
type AiExportRepository interface {
    Create(ctx context.Context, userID string, input models.CreateAiExportInput, generatedPrompt string) (*models.AiExportRecord, error)
    GetByID(ctx context.Context, userID string, id string) (*models.AiExportRecord, error)
    ListByUserID(ctx context.Context, userID string) ([]models.AiExportRecord, error)
    UpdateFilePath(ctx context.Context, userID string, id string, filePath string) (*models.AiExportRecord, error)
    Delete(ctx context.Context, userID string, id string) (*models.AiExportRecord, error)
}
```

### SLICE-W07-011: AiExport Service (Prompt Generation + Data Aggregation)

**File**: `apps/api/internal/atlas/service/ai_export_service.go`

This is the most complex slice. The service orchestrates:

#### Prompt Generation (no ZIP yet)
```go
type AiExportService interface {
    GeneratePrompt(ctx context.Context, userID string, input models.CreateAiExportInput) (*models.AiExport, error)
    GenerateExport(ctx context.Context, userID string, aiExportID string) (*models.AiExport, error)
    GetByID(ctx context.Context, userID string, id string) (*models.AiExport, error)
    List(ctx context.Context, userID string) ([]models.AiExport, error)
    Delete(ctx context.Context, userID string, id string) (*models.AiExport, error)
}
```

#### Dependencies (constructor injection)
- `AiExportRepository`
- `WeekFlagService` — for week flags in the period
- `UserProfileService` — for persistent AI context + goal
- **Read-only query services from prior waves** (these already exist as repository interfaces or services):
  - `BodyWeightEntryRepository` (for body_weight_entries by date range)
  - `BodyCheckInRepository` (for body_check_ins by date range)
  - `BodyMeasurementRepository` (for measurements by check-in IDs)
  - `ProgressPhotoRepository` (for photos by check-in IDs)
  - `CardioEntryRepository` (for cardio by date range)
  - `NutritionMacroService` (for macro aggregates per day/week — already aggregates from templates + overrides)
  - DailyLog/workout queries — WAVE-03 entities (need read-only access)

Actually, for data.json aggregation, the service needs raw data across all domains. Rather than injecting 10+ repos into the AiExport service, define a lightweight **export-data query interface** that the AiExport service depends on:

```go
type AiExportDataProvider interface {
    // WAVE-03 (if available)
    GetWorkoutSummary(ctx context.Context, userID string, from, to models.Date) ([]WorkoutSummaryItem, error)
    // WAVE-04
    GetBodyWeightEntries(ctx context.Context, userID string, from, to models.Date) ([]models.BodyWeightEntry, error)
    GetBodyCheckIns(ctx context.Context, userID string, from, to models.Date) ([]models.BodyCheckIn, error)
    GetCardioEntries(ctx context.Context, userID string, from, to models.Date) ([]models.CardioEntry, error)
    GetWeekFlags(ctx context.Context, userID string, weekStarts []models.Date) ([]models.WeekFlag, error)
    // WAVE-05
    GetNutritionMacros(ctx context.Context, userID string, from, to models.Date) ([]models.NutritionWeeklyAverage, error)
}
```

This interface can be implemented by a single `dataProvider` struct in the same package that wraps the individual repos/services. This keeps the AiExport service constructor manageable.

#### Prompt Building Logic
```
You are a fitness and nutrition AI coach analyzing training data.
[PERSISTENT_AI_CONTEXT from user_profile]
[USER_GOAL from user_profile]
Date range: {dateRangeStart} to {dateRangeEnd}

## Training Summary
{workout summaries aggregated from daily_log + workout_exercise + workout_set}

## Body Weight Trends
{body weight entries}

## Body Check-in Data
{check-ins + measurements}

## Cardio Activity
{cardio entries}

## Week Flags
{week flags in period}

## Nutrition Summary
{macro averages}

## User's Comment
{userComment or "None provided"}

Analyze the above data and provide actionable insights.
```

#### Export Flow (GeneratePrompt)
1. Validate `dateRangeStart <= dateRangeEnd`, range not excessive (max 52 weeks matches Setting's max).
2. Fetch user profile for `persistentAiContext` and `goal` (okay if nil).
3. Build prompt template string.
4. Create `AiExportRecord` with `generatedPrompt` and `export_file_path = NULL`.
5. Return the record (with prompt) to caller.

#### Export Flow (GenerateExport = ZIP creation)
1. Fetch `AiExportRecord` by ID.
2. If `export_file_path` already non-nil, return existing (idempotent).
3. Fetch all data from `AiExportDataProvider` for the date range.
4. Build `data.json`, `summary.md`, `manifest.json`, and CSV files in memory.
5. Package into ZIP, write to `cfg.Media.BasePath/exports/{userID}/{exportID}.zip`.
6. Update `export_file_path` in DB.
7. Return updated record.

### SLICE-W07-012: ZIP Generation Utilities

**File**: `apps/api/internal/atlas/service/export_zip.go`

Package-level helpers (not a separate package — keep in `service`):

```go
type ExportArchive struct {
    Manifest    ExportManifest
    DataJSON    map[string]any   // all entities keyed by type
    SummaryMD   string
    CSVs        map[string]string // filename -> CSV content
}
```

No filesystem dependencies in the archive builder — it just builds in-memory structures. The service layer handles the actual ZIP writing and path management.

### SLICE-W07-013: AiExport Resolver + GraphQL Schema

**Files**:
- `apps/api/internal/atlas/graph/schema/ai_export.graphql`
- `apps/api/internal/atlas/graph/resolver/ai_export.go`

Schema:
```graphql
type AiExport {
  id: ID!
  userId: ID!
  dateRangeStart: Date!
  dateRangeEnd: Date!
  includePhotos: Boolean!
  includeNutrition: Boolean!
  includeCardio: Boolean!
  includeMeasurements: Boolean!
  userComment: String
  generatedPrompt: String!
  exportFilePath: String
  createdAt: Time!
  updatedAt: Time!
}

input CreateAiExportInput {
  dateRangeStart: Date!
  dateRangeEnd: Date!
  includePhotos: Boolean!
  includeNutrition: Boolean!
  includeCardio: Boolean!
  includeMeasurements: Boolean!
  userComment: String
}

type AiExportResult {
  aiExport: AiExport
  validationError: AiExportValidationError
  notFoundError: AiExportNotFoundError
  authError: AiExportAuthError
}

type AiExportsResult {
  aiExports: [AiExport!]!
  validationError: AiExportValidationError
  authError: AiExportAuthError
}
```

Mutations:
- `createAiExportPrompt(input: CreateAiExportInput!): AiExportResult!` — generates prompt only
- `generateAiExport(id: ID!): AiExportResult!` — generates ZIP
- `deleteAiExport(id: ID!): AiExportResult!`

Queries:
- `aiExport(id: ID!): AiExportResult!`
- `aiExports: AiExportsResult!`

### SLICE-W07-014: Download Endpoint (REST, not GraphQL)

**File**: `apps/api/internal/handler/ai_export_handler.go`

Since ZIP download is binary, use a REST handler (not GraphQL):

```go
GET /api/v1/ai-exports/{id}/download
```

Pattern: mirror `ProgressPhotoHandler.Download` — guard with `AtlasPinGuard`, read `export_file_path` from service, serve the file with `http.ServeFile` or `http.ServeContent`.

Wire in `main.go` alongside `progressPhotos` routes.

### SLICE-W07-015: Main Wiring

**File**: `apps/api/cmd/server/main.go`

Add to the existing wiring block:

```go
atlasUserProfileRepo := atlasPostgres.NewUserProfileRepository(db.Pool)
atlasUserProfileService := atlasService.NewUserProfileService(atlasUserProfileRepo)

atlasAiExportRepo := atlasPostgres.NewAiExportRepository(db.Pool)
atlasAiExportDataProvider := atlasService.NewAiExportDataProvider(
    atlasBodyWeightRepo,    // or bodyWeightService
    atlasCheckInRepo,
    atlasMeasurementRepo,
    atlasPhotoRepo,
    atlasCardioRepo,
    atlasWeekFlagService,
    atlasNutritionWeeklyAvgService,
    // plus WAVE-03 workout data when available
)
atlasAiExportService := atlasService.NewAiExportService(
    atlasAiExportRepo,
    atlasUserProfileService,
    atlasAiExportDataProvider,
    l,
)
```

Add to `atlasResolver.Resolver`:
- `UserProfileService`
- `AiExportService`

Wire the download endpoint:
```go
atlas.Get("/api/v1/ai-exports/{id}/download", aiExportHandler.Download)
```

---

## 2. Acceptance Criteria (AC)

AC-W07-001: User creates an AI export prompt with date range, section toggles, and optional comment. Response contains a generated prompt string with all configured context sections.

AC-W07-002: Generating a ZIP export (via `generateAiExport`) produces a valid ZIP containing `manifest.json`, `data.json`, `summary.md`, and toggled CSV sections.

AC-W07-003: Excluding photos/measurements/nutrition/cardio individually removes the corresponding section from prompt and ZIP content.

AC-W07-004: Persistent AI context from `user_profile` appears in the prompt when set; absent context fields are gracefully omitted.

AC-W07-005: Week flags in the selected date range appear in the prompt and `data.json`.

AC-W07-006: Downloading a completed export serves the ZIP file with correct MIME type and Content-Disposition.

AC-W07-007: User can update their profile fields; all fields accept null/omission for partial updates.

AC-W07-008: Requesting a prompt with `dateRangeStart > dateRangeEnd` returns a validation error.

AC-W07-009: Requesting a prompt with a range > 52 weeks returns a validation error.

AC-W07-010: Authenticated user sees only their own exports and profile. Requests without valid session return auth error.

AC-W07-011: Deleting an export removes the database row and the associated file from disk.

---

## 3. Exit Criteria (EC)

EC-W07-001: `user_profiles` migration file exists and passes sqlc generation.
EC-W07-002: `ai_exports` migration file exists and passes sqlc generation.
EC-W07-003: All sqlc queries compile and regenerate without errors.
EC-W07-004: All Go code compiles (`bun run build` succeeds).
EC-W07-005: GraphQL schema passes gqlgen generation without errors.
EC-W07-006: `generateAiExport(id)` produces a valid ZIP file on disk at the configured path.
EC-W07-007: Prompt contents reflect all toggles correctly (each section inclusion/exclusion tested).
EC-W07-008: At minimum, service unit tests exist for:
  - Prompt building with all sections included vs excluded
  - Profile partial update with nil/zero fields
  - Date range validation
  - Data aggregation with empty result sets (no data for period)
EC-W07-009: All lint checks pass (`bun run lint`).

---

## 4. Verification Obligations

### Verification against verification-plan.xml
- Add `V-M-API-USER-PROFILE` and `V-M-API-AI-EXPORT` entries to `docs/verification-plan.xml`.
- Extend the existing `V-M-API` entry's scope to cover export-download endpoint integration.
- Add trace markers: `[AiExport][generatePrompt]`, `[AiExport][generateExport]`, `[AiExport][download]`, `[UserProfile][update]`.

### Tests required
| Layer | File | What to test |
|-------|------|-------------|
| Service | `user_profile_service_test.go` | Get (exists/nil), Update (full/partial), validation |
| Service | `ai_export_service_test.go` | Prompt generation structure, section toggles, date validation, range limit |
| Service | `export_zip_test.go` | ZIP structure (manifest.json keys, data.json shape, summary.md non-empty), CSV format |
| Resolver | `ai_export_test.go` (if e2e) | Full GraphQL mutation flow with auth context |
| Integration | handler | Download endpoint serves correct file, 404 for missing export, unauth 401 |

### Log markers
```
[AiExport][generatePrompt] user_id={} date_range={}-{} toggles={}
[AiExport][generateExport] user_id={} export_id={}
[AiExport][download] user_id={} export_id={}
[UserProfile][update] user_id={} fields={}
```

### Traceability
```
docs/product/prd.md Sections 17, 18
docs/product-verified/domain-model.md#AiExport
docs/requirements.xml // after update
```

---

## 5. Risks

| Risk | Severity | Mitigation |
|------|----------|-----------|
| ZIP generation for large date ranges (many photos) is slow/expensive | Medium | Stream archive entries; set a max-photos limit or warn on > N photos; consider async generation with status polling |
| WAVE-03 (workout_exercise + workout_set) tables may not exist yet | High | Make workout data optional in prompt; AiExportDataProvider returns empty slices if tables don't exist or are nil |
| Export ZIP file fills up disk | Low | Use configured `cfg.Media.BasePath/exports/` subdirectory; add periodic cleanup for stale exports |
| Persistent AI context grows unbounded | Low | `persistent_ai_context` is TEXT (unlimited); add a character limit in the service (e.g. 5000 chars) with truncation warning |
| sqlc generates types for nullable fields (float32, Date) | Medium | Test sqlc generation for `REAL` → `*float32` and `DATE` → `pgtype.Date`; use `nullableText`/`nullableReal` helpers from `settings_repo.go` pattern |

---

## 6. Rollback Plan

1. **Reverse migrations**: `goose down` twice (00093 and 00094).
2. **Remove source files**: all files in slices W07-001 through W07-014.
3. **Revert main.go**: remove wiring, resolver fields, and route.
4. **Revert schema.graphql**: remove `userProfile`/`updateUserProfile`/`aiExport`/`aiExports`/`createAiExportPrompt`/`generateAiExport`/`deleteAiExport` queries/mutations.
5. **Regenerate**: `bun run codegen` (both sqlc and gqlgen) to purge generated code.
6. **Clean up exports dir**: `rm -rf <media-base>/exports/`.

---

## 7. Open Questions

| ID | Question | Why It Matters | Asked Of | Status |
|----|----------|---------------|----------|--------|
| Q-W07-001 | WAVE-03 (workout_exercise, workout_set) — do these tables exist yet? | AiExport data.json needs workout data; if absent, data.json has empty training section | WAVE-03 owner | Open |
| Q-W07-002 | Should ZIP generation be synchronous (wait in GraphQL) or async (create job, poll status)? | UX: large exports could timeout GraphQL | Product | Open |
| Q-W07-003 | What is the max allowed ZIP size? Should we limit date range dynamically based on photo count? | Disk usage and response time | Product | Open |
| Q-W07-004 | How does the frontend trigger download? Direct link or stream through API? | Influences if we need auth token on download URL (already guarded by PIN session cookie) | Frontend | Open |

---

## 8. Graph Deltas (for docs/knowledge-graph.xml)

```
M-API-USER-PROFILE → module for UserProfile CRUD
  depends: M-API (shared models, date, auth)
  files: models/user_profile.go, service/user_profile_service.go,
         repository/postgres/user_profile_repo.go, graph/schema/user_profile.graphql,
         graph/resolver/user_profile.go, repository/postgres/queries/user_profiles.sql

M-API-AI-EXPORT → module for AiExport prompt and ZIP generation
  depends: M-API, M-API-USER-PROFILE, M-API-WEEK-FLAG, M-API-BODY (WAVE-04),
           M-API-CARDIO (WAVE-04), M-API-NUTRITION (WAVE-05)
  files: models/ai_export.go, service/ai_export_service.go, service/export_zip.go,
         repository/postgres/ai_export_repo.go, graph/schema/ai_export.graphql,
         graph/resolver/ai_export.go, handler/ai_export_handler.go,
         repository/postgres/queries/ai_exports.sql
```

---

## 9. Dependency Order

```
SLICE-W07-001 (migration profile) → SLICE-W07-002 (model) → SLICE-W07-003 (sqlc) → SLICE-W07-004 (repo) → SLICE-W07-005 (service) → SLICE-W07-006 (resolver+schema)
SLICE-W07-007 (migration export) → SLICE-W07-008 (model) → SLICE-W07-009 (sqlc) → SLICE-W07-010 (repo) → SLICE-W07-011 (service) + SLICE-W07-012 (zip utils) → SLICE-W07-013 (resolver+schema) → SLICE-W07-014 (download handler)
SLICE-W07-015 (main wiring) — depends on all above
```

The user profile and ai export chains are **independent** and can be built in parallel until slice 015.

---

## 10. Reviewer Checklist

- [ ] Migration files use correct sequential numbering (00093, 00094)
- [ ] All nullable columns use `COALESCE` pattern in upsert queries
- [ ] sqlc generates without errors for REAL, DATE, UUID, TEXT nullable types
- [ ] gqlgen config has explicit bindings for all new types
- [ ] AiExportDataProvider doesn't leak into resolver layer
- [ ] Download endpoint is guarded by AtlasPinGuard (same as all other guarded routes)
- [ ] Prompt generation is separable from ZIP generation (user can get prompt without building ZIP)
- [ ] No secrets or API keys in prompt generation (no ChatGPT/OpenAI integration per excluded scope)
