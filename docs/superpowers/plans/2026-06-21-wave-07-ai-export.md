# WAVE-07: AI Export and Prompt Builder Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Generate AI-ready exports with structured data and prompts for ChatGPT analysis — UserProfile CRUD, prompt builder, ZIP export with manifest/data/summary/CSVs/photos.

**Architecture:** Two new entity domains (UserProfile, AiExport) across the full stack — migration → sqlc → model → repo → service → resolver/schema → wiring. AiExportDataProvider aggregates from 7 existing data sources. ZIP generation via Go stdlib with temp-file-atomic-rename.

**Tech Stack:** Go 1.25, pgx/v5, sqlc, gqlgen, archive/zip (stdlib), chi router (existing REST handlers)

**Spec:** `docs/superpowers/specs/2026-06-21-wave-07-ai-export-design.md`
**Full brief:** `docs/prd-wave-details/waves/wave-07.md`

---

## File Structure

### New Files (17)
| File | Purpose |
|------|---------|
| `apps/api/internal/repository/postgres/migrations/00091_user_profiles.sql` | UserProfile table |
| `apps/api/internal/repository/postgres/migrations/00092_ai_exports.sql` | AiExport table |
| `apps/api/internal/repository/postgres/queries/user_profiles.sql` | sqlc queries |
| `apps/api/internal/repository/postgres/queries/ai_exports.sql` | sqlc queries |
| `apps/api/internal/atlas/models/user_profile.go` | UserProfile types, result unions, errors |
| `apps/api/internal/atlas/models/ai_export.go` | AiExport types, result unions, errors |
| `apps/api/internal/atlas/models/zipper.go` | ZIP generation utility |
| `apps/api/internal/atlas/repository/postgres/user_profile_repo.go` | UserProfile repo adapter |
| `apps/api/internal/atlas/repository/postgres/ai_export_repo.go` | AiExport repo adapter |
| `apps/api/internal/atlas/service/user_profile_service.go` | UserProfile business logic |
| `apps/api/internal/atlas/service/ai_export_service.go` | AiExport business logic |
| `apps/api/internal/atlas/service/ai_export_data_provider.go` | Data aggregation interface+impl |
| `apps/api/internal/atlas/graph/schema/user_profile.graphql` | UserProfile GraphQL types |
| `apps/api/internal/atlas/graph/schema/ai_export.graphql` | AiExport GraphQL types |
| `apps/api/internal/atlas/graph/resolver/user_profile.go` | UserProfile resolver |
| `apps/api/internal/atlas/graph/resolver/ai_export.go` | AiExport resolver |
| `apps/api/internal/atlas/handler/ai_export_handler.go` | REST handler (generate, download) |

### Modified Files (5)
| File | Changes |
|------|---------|
| `apps/api/atlas-gqlgen.yml` | Model bindings for UserProfile, AiExport types |
| `apps/api/internal/atlas/graph/schema/schema.graphql` | UserProfile + AiExport queries/mutations |
| `apps/api/internal/atlas/graph/resolver/resolver.go` | Add UserProfileService, AiExportService fields |
| `apps/api/cmd/server/main.go` | Wire services, repos, routes |
| `apps/api/internal/atlas/service/bootstrap_service.go` | EnsureDefaultUser creates default UserProfile |

### Test Files (created alongside each component)
`apps/api/internal/atlas/service/user_profile_service_test.go`
`apps/api/internal/atlas/service/ai_export_service_test.go`
`apps/api/internal/atlas/handler/ai_export_handler_test.go`
`apps/api/internal/atlas/graph/resolver/user_profile_resolver_test.go`
`apps/api/internal/atlas/graph/resolver/ai_export_resolver_test.go`


### Task 1: UserProfile Migration and Config

**Files:**
- Create: `apps/api/internal/repository/postgres/migrations/00091_user_profiles.sql`

- [ ] **Step 1: Write migration**

```sql
CREATE TABLE user_profiles (
    id                           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                      UUID NOT NULL REFERENCES atlas_users(id) UNIQUE,
    goal                         TEXT,
    height                       REAL,
    birth_date                   DATE,
    training_experience          TEXT,
    current_training_split       TEXT,
    preferred_progression_style  TEXT,
    nutrition_strategy           TEXT,
    persistent_ai_context        TEXT,
    created_at                   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                   TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

- [ ] **Step 2: Run migration check**

Run: `INTEGRATION_TESTS=1 go test -run TestWave07Migration -v`
Expected: Migration applies without error

- [ ] **Step 3: Add AiExportConfig to appconfig**

Read `apps/api/internal/appconfig/config.go` to find existing config pattern. Add:
```go
type AiExportConfig struct {
    BasePath       string `yaml:"base_path" env:"AI_EXPORT_BASE_PATH" env-default:"./data/ai-export"`
    MaxSizeMB      int    `yaml:"max_size_mb" env:"AI_EXPORT_MAX_SIZE_MB" env-default:"100"`
    MaxPhotos      int    `yaml:"max_photos" env:"AI_EXPORT_MAX_PHOTOS" env-default:"20"`
    DefaultWeeks   int    `yaml:"default_weeks" env:"AI_EXPORT_DEFAULT_WEEKS" env-default:"4"`
    TTLDays        int    `yaml:"ttl_days" env:"AI_EXPORT_TTL_DAYS" env-default:"7"`
    MaxRangeDays   int    `yaml:"max_range_days" env:"AI_EXPORT_MAX_RANGE_DAYS" env-default:"365"`
}
```

- [ ] **Step 4: Wire AiExportConfig in main Config struct**

Add `AiExport AiExportConfig` field to root Config in `apps/api/internal/appconfig/config.go`

- [ ] **Step 5: Commit**

```bash
git add apps/api/internal/repository/postgres/migrations/00091_user_profiles.sql apps/api/internal/appconfig/config.go
git commit -m "feat(wave-07): add user_profiles migration and AiExport config"
```


### Task 2: AiExport Migration

**Files:**
- Create: `apps/api/internal/repository/postgres/migrations/00092_ai_exports.sql`

- [ ] **Step 1: Write migration**

```sql
CREATE TABLE ai_exports (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES atlas_users(id),
    date_range_start    DATE NOT NULL,
    date_range_end      DATE NOT NULL,
    include_photos      BOOLEAN NOT NULL DEFAULT false,
    include_nutrition   BOOLEAN NOT NULL DEFAULT true,
    include_cardio      BOOLEAN NOT NULL DEFAULT true,
    include_measurements BOOLEAN NOT NULL DEFAULT true,
    user_comment        TEXT,
    generated_prompt    TEXT NOT NULL,
    export_file_path    TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_ai_exports_user_created ON ai_exports(user_id, created_at DESC);
```

- [ ] **Step 2: Run migration check**

Run: `INTEGRATION_TESTS=1 go test -run TestWave07Migration -v`

- [ ] **Step 3: Commit**

```bash
git add apps/api/internal/repository/postgres/migrations/00092_ai_exports.sql
git commit -m "feat(wave-07): add ai_exports migration"
```


### Task 3: sqlc Queries (UserProfile + AiExport)

**Files:**
- Create: `apps/api/internal/repository/postgres/queries/user_profiles.sql`
- Create: `apps/api/internal/repository/postgres/queries/ai_exports.sql`

- [ ] **Step 1: Write user_profiles.sql queries**

```sql
-- name: GetUserProfileByUserID :one
SELECT id, user_id, goal, height, birth_date, training_experience,
       current_training_split, preferred_progression_style,
       nutrition_strategy, persistent_ai_context, created_at, updated_at
FROM user_profiles WHERE user_id = $1;

-- name: UpsertUserProfile :one
INSERT INTO user_profiles (user_id, goal, height, birth_date, training_experience,
    current_training_split, preferred_progression_style, nutrition_strategy, persistent_ai_context)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (user_id) DO UPDATE SET
    goal = COALESCE($2, user_profiles.goal),
    height = COALESCE($3, user_profiles.height),
    birth_date = COALESCE($4, user_profiles.birth_date),
    training_experience = COALESCE($5, user_profiles.training_experience),
    current_training_split = COALESCE($6, user_profiles.current_training_split),
    preferred_progression_style = COALESCE($7, user_profiles.preferred_progression_style),
    nutrition_strategy = COALESCE($8, user_profiles.nutrition_strategy),
    persistent_ai_context = COALESCE($9, user_profiles.persistent_ai_context),
    updated_at = now()
RETURNING id, user_id, goal, height, birth_date, training_experience,
          current_training_split, preferred_progression_style,
          nutrition_strategy, persistent_ai_context, created_at, updated_at;
```

- [ ] **Step 2: Write ai_exports.sql queries**

```sql
-- name: CreateAiExport :one
INSERT INTO ai_exports (user_id, date_range_start, date_range_end, include_photos,
    include_nutrition, include_cardio, include_measurements, user_comment, generated_prompt, export_file_path)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetAiExportByID :one
SELECT * FROM ai_exports WHERE id = $1;

-- name: GetAiExportByUserID :many
SELECT * FROM ai_exports WHERE user_id = $1 ORDER BY created_at DESC;

-- name: DeleteAiExport :exec
DELETE FROM ai_exports WHERE id = $1;

-- name: DeleteExpiredAiExports :exec
DELETE FROM ai_exports WHERE created_at < NOW() - ($1::int * INTERVAL '1 day');

-- name: UpdateAiExportFilePath :one
UPDATE ai_exports SET export_file_path = $2 WHERE id = $1 RETURNING *;
```

- [ ] **Step 3: Run sqlc codegen**

Run: `bunx nx run api:codegen`
Expected: Generated files appear under `apps/api/internal/repository/postgres/generated/`

- [ ] **Step 4: Commit**

```bash
git add apps/api/internal/repository/postgres/queries/user_profiles.sql apps/api/internal/repository/postgres/queries/ai_exports.sql
git commit -m "feat(wave-07): add sqlc queries for user_profiles and ai_exports"
```


### Task 4: UserProfile Model

**Files:**
- Create: `apps/api/internal/atlas/models/user_profile.go`

- [ ] **Step 1: Write model types**

Follow `models/settings.go` pattern:
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

- [ ] **Step 2: Write result/error types**

```go
type UserProfileResult struct {
    Success      *UserProfile
    NotFoundErr  *UserProfileNotFoundErr
    ValidationErr *UserProfileValidationErr
    AuthErr      *UserProfileAuthErr
}

type UserProfileNotFoundErr struct { Message string; Code UserProfileErrorCode }
type UserProfileValidationErr struct { Message string; Code UserProfileErrorCode }
type UserProfileAuthErr struct { Message string; Code UserProfileErrorCode }
type UserProfileErrorCode string
const (
    UserProfileErrorValidation UserProfileErrorCode = "VALIDATION_ERROR"
    UserProfileErrorNotFound   UserProfileErrorCode = "NOT_FOUND"
    UserProfileErrorAuth       UserProfileErrorCode = "AUTH_ERROR"
)
```

- [ ] **Step 3: Commit**

```bash
git add apps/api/internal/atlas/models/user_profile.go
git commit -m "feat(wave-07): add UserProfile model types"
```


### Task 5: AiExport Model

**Files:**
- Create: `apps/api/internal/atlas/models/ai_export.go`

- [ ] **Step 1: Write model types**

Follow same pattern as UserProfile model:
- AiExportRecord (DB), AiExport (public), AiExportInput
- AiExportResult union type with Success/ValidationErr/NotFoundErr/AuthErr
- AiExportErrorCode enum
- AiExportListResult with Success/Error variants

- [ ] **Step 2: Commit**

```bash
git add apps/api/internal/atlas/models/ai_export.go
git commit -m "feat(wave-07): add AiExport model types"
```


### Task 6: ZIP Generation Utility

**Files:**
- Create: `apps/api/internal/atlas/models/zipper.go`

- [ ] **Step 1: Write ZIP builder**

```go
type ExportZIPBuilder struct {
    buf        *bytes.Buffer
    zipWriter  *zip.Writer
    tmpPath    string
    finalPath  string
}

func NewExportZIPBuilder(finalPath string) *ExportZIPBuilder
func (b *ExportZIPBuilder) AddManifest(manifest ExportManifest) error
func (b *ExportZIPBuilder) AddDataJSON(data interface{}) error
func (b *ExportZIPBuilder) AddSummaryMD(summary string) error
func (b *ExportZIPBuilder) AddCSV(name string, records [][]string) error
func (b *ExportZIPBuilder) AddPhoto(path, name string) error
func (b *ExportZIPBuilder) Close() error  // temp-file-atomic-rename
```

Use temp file pattern:
1. Create temp file in same directory as finalPath
2. Write ZIP content
3. Close writer
4. Rename(tempPath, finalPath) — atomic

Include size check: if estimated content exceeds `MaxSizeMB`, return error before building.

- [ ] **Step 2: Commit**

```bash
git add apps/api/internal/atlas/models/zipper.go
git commit -m "feat(wave-07): add ZIP generation utility with temp-file-atomic-rename"
```


### Task 7: UserProfile Repository

**Files:**
- Create: `apps/api/internal/atlas/repository/postgres/user_profile_repo.go`

- [ ] **Step 1: Write repo**

```go
type UserProfileRepository interface {
    GetByUserID(ctx context.Context, userID string) (*models.UserProfileRecord, error)
    Upsert(ctx context.Context, userID string, input models.UserProfileInput) (*models.UserProfileRecord, error)
}

type userProfileRepository struct {
    q *generated.Queries
}

func NewUserProfileRepository(pgxPool pgxpool) UserProfileRepository
```

Read generated function names from `apps/api/internal/repository/postgres/generated/querier.go` after codegen.

- [ ] **Step 2: Write unit test with mock**

- [ ] **Step 3: Commit**

```bash
git add apps/api/internal/atlas/repository/postgres/user_profile_repo.go
git commit -m "feat(wave-07): add UserProfile repository"
```


### Task 8: AiExport Repository

**Files:**
- Create: `apps/api/internal/atlas/repository/postgres/ai_export_repo.go`

- [ ] **Step 1: Write repo**

```go
type AiExportRepository interface {
    Create(ctx context.Context, params CreateAiExportParams) (*models.AiExportRecord, error)
    GetByID(ctx context.Context, id string) (*models.AiExportRecord, error)
    GetByUserID(ctx context.Context, userID string) ([]*models.AiExportRecord, error)
    Delete(ctx context.Context, id string) error
    DeleteExpired(ctx context.Context, ttlDays int32) error
    UpdateFilePath(ctx context.Context, id, filePath string) (*models.AiExportRecord, error)
}
```

- [ ] **Step 2: Commit**

```bash
git add apps/api/internal/atlas/repository/postgres/ai_export_repo.go
git commit -m "feat(wave-07): add AiExport repository"
```


### Task 9: UserProfile Service

**Files:**
- Create: `apps/api/internal/atlas/service/user_profile_service.go`
- Create: `apps/api/internal/atlas/service/user_profile_service_test.go`

- [ ] **Step 1: Write service interface + implementation**

```go
type UserProfileService interface {
    Get(ctx context.Context, userID string) (*models.UserProfile, error)
    Update(ctx context.Context, userID string, input models.UserProfileInput) (*models.UserProfile, error)
}
```

Implementation:
- `Get`: query repo, convert Record → public, if not found return empty with all nil fields
- `Update`: call repo.Upsert, convert Record → public

- [ ] **Step 2: Write test**

```go
func TestUserProfileService_Get_ReturnsEmpty(t *testing.T)
func TestUserProfileService_Get_ReturnsProfile(t *testing.T)
func TestUserProfileService_Update_Success(t *testing.T)
func TestUserProfileService_Update_Partial(t *testing.T)
```

- [ ] **Step 3: Run tests**

Run: `go test ./internal/atlas/service/ -run TestUserProfileService -v`

- [ ] **Step 4: Commit**

```bash
git add apps/api/internal/atlas/service/user_profile_service.go apps/api/internal/atlas/service/user_profile_service_test.go
git commit -m "feat(wave-07): add UserProfile service"
```


### Task 10: AiExport Service

**Files:**
- Create: `apps/api/internal/atlas/service/ai_export_service.go`
- Create: `apps/api/internal/atlas/service/ai_export_service_test.go`

- [ ] **Step 1: Write service interface**

```go
type AiExportService interface {
    Generate(ctx context.Context, userID string, input GenerateAiExportInput) (*GenerateAiExportResult, error)
    GetByID(ctx context.Context, userID, exportID string) (*models.AiExport, error)
    Cleanup(ctx context.Context, ttlDays int) error
}

type GenerateAiExportInput struct {
    DateRangeStart    Date
    DateRangeEnd      Date
    IncludePhotos     bool
    IncludeNutrition  bool
    IncludeCardio     bool
    IncludeMeasurements bool
    UserComment       *string
}
```

- [ ] **Step 2: Implement Generate method**

1. Validate dates (end >= start, max 365 days)
2. Delete prior exports for this user (regeneration cleanup per DDEC-W07-010)
3. Build prompt from UserProfile context + sections + comment + week flags
4. Build data JSON via AiExportDataProvider
5. Build ZIP via ExportZIPBuilder
6. Save AiExport record to DB (with export file path)
7. Return result with export ID and generated prompt

- [ ] **Step 3: Write tests**

```go
func TestAiExportService_Generate_Basic(t *testing.T)
func TestAiExportService_Generate_WithPhotos(t *testing.T)
func TestAiExportService_Generate_EmptyPeriod(t *testing.T)
func TestAiExportService_Generate_InvalidDates(t *testing.T)
func TestAiExportService_Generate_DateRangeOverLimit(t *testing.T)
func TestAiExportService_GetByID_Success(t *testing.T)
func TestAiExportService_GetByID_Ownership(t *testing.T)
func TestAiExportService_Generate_PromptContainsContext(t *testing.T)
```

- [ ] **Step 4: Run tests**

Run: `go test ./internal/atlas/service/ -run TestAiExportService -v`

- [ ] **Step 5: Commit**

```bash
git add apps/api/internal/atlas/service/ai_export_service.go apps/api/internal/atlas/service/ai_export_service_test.go
git commit -m "feat(wave-07): add AiExport service"
```


### Task 11: AiExportDataProvider

**Files:**
- Create: `apps/api/internal/atlas/service/ai_export_data_provider.go`

- [ ] **Step 1: Write data provider interface**

```go
type AiExportDataProvider interface {
    GetWorkoutData(ctx context.Context, userID string, from, to Date) ([]WorkoutExport, error)
    GetCardioData(ctx context.Context, userID string, from, to Date) ([]CardioExport, error)
    GetBodyWeightData(ctx context.Context, userID string, from, to Date) ([]BodyWeightExport, error)
    GetMeasurementData(ctx context.Context, userID string, from, to Date) ([]MeasurementExport, error)
    GetNutritionData(ctx context.Context, userID string, from, to Date) ([]NutritionExport, error)
    GetWeekFlags(ctx context.Context, userID string, from, to Date) ([]WeekFlagExport, error)
    GetPhotos(ctx context.Context, userID string, from, to Date) ([]PhotoExport, error)
    GetUserProfile(ctx context.Context, userID string) (*UserProfileExport, error)
}
```

- [ ] **Step 2: Implement data provider**

Each method queries the corresponding existing service/repo:
- Workout: calls WAVE-03 repos (returns empty if not deployed)
- Cardio: calls CardioEntryRepo.ListByDateRange
- BodyWeight: calls BodyWeightService.ListByDateRange
- Measurements: calls BodyCheckInService with measurement joins
- Nutrition: calls NutritionMacroService + template/override services
- WeekFlags: calls WeekFlagService
- Photos: calls ProgressPhotoRepo (with count cap at MaxPhotos)
- UserProfile: calls UserProfileService.Get

- [ ] **Step 3: Commit**

```bash
git add apps/api/internal/atlas/service/ai_export_data_provider.go
git commit -m "feat(wave-07): add AiExportDataProvider for export data aggregation"
```


### Task 12: GraphQL Schemas (UserProfile + AiExport)

**Files:**
- Create: `apps/api/internal/atlas/graph/schema/user_profile.graphql`
- Create: `apps/api/internal/atlas/graph/schema/ai_export.graphql`
- Modify: `apps/api/internal/atlas/graph/schema/schema.graphql`

- [ ] **Step 1: Write user_profile.graphql**

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
    success: UserProfile
    notFoundErr: UserProfileNotFoundError
    validationErr: UserProfileValidationError
    authErr: UserProfileAuthError
}

type UserProfileNotFoundError { message: String!; code: UserProfileErrorCode! }
type UserProfileValidationError { message: String!; code: UserProfileErrorCode! }
type UserProfileAuthError { message: String!; code: UserProfileErrorCode! }
enum UserProfileErrorCode { VALIDATION_ERROR NOT_FOUND AUTH_ERROR INTERNAL_ERROR }
```

- [ ] **Step 2: Write ai_export.graphql**

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
}

input GenerateAiExportInput {
    dateRangeStart: Date!
    dateRangeEnd: Date!
    includePhotos: Boolean = false
    includeNutrition: Boolean = true
    includeCardio: Boolean = true
    includeMeasurements: Boolean = true
    userComment: String
}

type AiExportResult {
    success: AiExport
    notFoundErr: AiExportNotFoundError
    validationErr: AiExportValidationError
    authErr: AiExportAuthError
}

type AiExportNotFoundError { message: String!; code: AiExportErrorCode! }
type AiExportValidationError { message: String!; code: AiExportErrorCode! }
type AiExportAuthError { message: String!; code: AiExportErrorCode! }
enum AiExportErrorCode { VALIDATION_ERROR NOT_FOUND AUTH_ERROR INTERNAL_ERROR }

type GenerateAiExportResult {
    success: GenerateAiExportSuccess
    validationErr: AiExportValidationError
    authErr: AiExportAuthError
}

type GenerateAiExportSuccess {
    export: AiExport!
    generatedPrompt: String!
}

type AiExportResult {
    success: AiExport
    notFoundErr: AiExportNotFoundError
    validationErr: AiExportValidationError
    authErr: AiExportAuthError
}

type AiExportListResult {
    success: [AiExport!]
    authErr: AiExportAuthError
}

type DeleteAiExportResult {
    success: Boolean
    notFoundErr: AiExportNotFoundError
    authErr: AiExportAuthError
}
```

- [ ] **Step 3: Add queries/mutations to schema.graphql**

```graphql
# WAVE-07 UserProfile
userProfile: UserProfileResult!
updateUserProfile(input: UserProfileInput!): UserProfileResult!

# WAVE-07 AiExport
generateAiExport(input: GenerateAiExportInput!): GenerateAiExportResult!
aiExport(id: ID!): AiExportResult!
aiExports: AiExportListResult!
deleteAiExport(id: ID!): DeleteAiExportResult!
```

- [ ] **Step 4: Commit**

```bash
git add apps/api/internal/atlas/graph/schema/user_profile.graphql apps/api/internal/atlas/graph/schema/ai_export.graphql apps/api/internal/atlas/graph/schema/schema.graphql
git commit -m "feat(wave-07): add GraphQL schemas for UserProfile and AiExport"
```


### Task 13: GraphQL Resolvers (UserProfile + AiExport)

**Files:**
- Create: `apps/api/internal/atlas/graph/resolver/user_profile.go`
- Create: `apps/api/internal/atlas/graph/resolver/ai_export.go`

- [ ] **Step 1: Write UserProfile resolver**

```go
func (r *Resolver) UserProfile(ctx context.Context) (*models.UserProfileResult, error) { ... }
func (r *Resolver) UpdateUserProfile(ctx context.Context, input models.UserProfileInput) (*models.UserProfileResult, error) { ... }
```

Pattern: `middleware.GetAtlasUserID(ctx)`, empty → return auth error. Call service, map to result union.

- [ ] **Step 2: Write AiExport resolver**

```go
func (r *Resolver) GenerateAiExport(ctx context.Context, input models.GenerateAiExportInput) (*models.GenerateAiExportResult, error) { ... }
```

- [ ] **Step 3: Run gqlgen codegen**

Run: `bunx nx run api:codegen` then `bunx nx run graphql:validate`
Fix any drift before proceeding.

- [ ] **Step 4: Commit**

```bash
git add apps/api/internal/atlas/graph/resolver/user_profile.go apps/api/internal/atlas/graph/resolver/ai_export.go
git commit -m "feat(wave-07): add GraphQL resolvers for UserProfile and AiExport"
```


### Task 14: REST Handlers

**Files:**
- Create: `apps/api/internal/atlas/handler/ai_export_handler.go`
- Create: `apps/api/internal/atlas/handler/user_profile_handler.go`
- Create: `apps/api/internal/atlas/handler/ai_export_handler_test.go`

- [ ] **Step 1: Write download handler**

Follow `ProgressPhotoHandler.Download` pattern:
```go
type AiExportHandler struct {
    aiExportService service.AiExportService
    exportBasePath  string
}

func (h *AiExportHandler) Generate(w http.ResponseWriter, r *http.Request) { ... }
func (h *AiExportHandler) Download(w http.ResponseWriter, r *http.Request) { ... }
```

Generate logic:
1. Extract userID from PIN middleware
2. Parse JSON body: dateRangeStart, dateRangeEnd, section toggles, userComment
3. Call service.Generate
4. Return JSON with exportId and generatedPrompt

Download logic:
1. Extract userID from PIN middleware
2. Read `exportId` from query param
3. Call service.GetByID to verify ownership
4. Open file from `{exportBasePath}/{userID}/{exportId}.zip`
5. Set Content-Type, Content-Disposition headers
6. Stream file to response

- [ ] **Step 2: Write user profile handler**

```go
type UserProfileHandler struct {
    userProfileService service.UserProfileService
}

func (h *UserProfileHandler) Get(w http.ResponseWriter, r *http.Request) { ... }
```

Get logic:
1. Extract userID from PIN middleware
2. Call service.Get
3. Return JSON with user profile (or empty profile with nil fields)

- [ ] **Step 3: Write handler tests**

```go
func TestAiExportHandler_Generate_Success(t *testing.T)
func TestAiExportHandler_Generate_InvalidInput(t *testing.T)
func TestAiExportHandler_Download_Success(t *testing.T)
func TestAiExportHandler_Download_MissingAuth(t *testing.T)
func TestAiExportHandler_Download_NotFound(t *testing.T)
func TestAiExportHandler_Download_OwnershipMismatch(t *testing.T)
func TestUserProfileHandler_Get_Success(t *testing.T)
func TestUserProfileHandler_Get_MissingAuth(t *testing.T)
```

- [ ] **Step 4: Commit**

```bash
git add apps/api/internal/atlas/handler/ai_export_handler.go apps/api/internal/atlas/handler/user_profile_handler.go apps/api/internal/atlas/handler/ai_export_handler_test.go
git commit -m "feat(wave-07): add REST handlers for AiExport generate/download and UserProfile get"
```


### Task 15: Bootstrap Service Update

**Files:**
- Modify: `apps/api/internal/atlas/service/bootstrap_service.go`

- [ ] **Step 1: Extend EnsureDefaultUser**

After creating atlas_users row and settings row, add:
```go
// Create default user profile with nil fields
_, err := s.userProfileRepo.Upsert(ctx, userID, models.UserProfileInput{})
```

Add `userProfileRepo` field to bootstrap service constructor.

- [ ] **Step 2: Commit**

```bash
git add apps/api/internal/atlas/service/bootstrap_service.go
git commit -m "feat(wave-07): extend bootstrap to create default UserProfile"
```


### Task 16: Main Wiring

**Files:**
- Modify: `apps/api/internal/atlas/graph/resolver/resolver.go`
- Modify: `apps/api/atlas-gqlgen.yml`
- Modify: `apps/api/cmd/server/main.go`
- Modify: `apps/api/config.yml`

- [ ] **Step 1: Add AiExport config to config.yml**

```yaml
ai_export:
  base_path: ./data/ai-export
  max_size_mb: 100
  max_photos: 20
  default_weeks: 4
  ttl_days: 7
  max_range_days: 365
```

```go
UserProfileService service.UserProfileService
AiExportService    service.AiExportService
```

- [ ] **Step 2: Add model bindings to atlas-gqlgen.yml**

Add bindings for:
- UserProfile, UserProfileInput, UserProfileResult, UserProfileNotFoundError, UserProfileValidationError, UserProfileAuthError, UserProfileErrorCode
- AiExport, GenerateAiExportInput, GenerateAiExportResult, GenerateAiExportSuccess, AiExportResult, AiExportListResult, DeleteAiExportResult, AiExportNotFoundError, AiExportValidationError, AiExportAuthError, AiExportErrorCode

Each follows the pattern:
```yaml
UserProfile:
    model: monorepo-template/apps/api/internal/atlas/models.UserProfile
```

- [ ] **Step 3: Wire services and routes in main.go**

1. Create repo instances (similar to existing postgres repos)
2. Create service instances
3. Create AiExportHandler with config
4. Add to Resolver struct
5. Register REST routes:
   ```go
   atlas.Post("/api/ai-export/generate", aiExportHandler.Generate)
   atlas.Get("/api/ai-export/download", aiExportHandler.Download)
   atlas.Get("/api/user-profile", userProfileHandler.Get)
   ```
6. Create export base directory on startup: `os.MkdirAll(cfg.AiExport.BasePath, 0755)`

- [ ] **Step 4: Run codegen and build**

```bash
bunx nx run api:codegen && bunx nx build api
```

- [ ] **Step 5: Lint**

```bash
bunx nx run api:lint
```

- [ ] **Step 6: Commit**

```bash
git add apps/api/cmd/server/main.go apps/api/internal/atlas/graph/resolver/resolver.go apps/api/atlas-gqlgen.yml
git commit -m "feat(wave-07): wire UserProfile and AiExport services, resolvers, and routes"
```


### Task 17: Integration Tests

**Files:**
- Create/modify integration tests

- [ ] **Step 1: Write repository integration tests**

```
INTEGRATION_TESTS=1 go test -run TestUserProfileRepo -v
INTEGRATION_TESTS=1 go test -run TestAiExportRepo -v
INTEGRATION_TESTS=1 go test -run TestWave07Migration -v
```

- [ ] **Step 2: Run full test suite for WAVE-07**

```bash
go test ./internal/atlas/service/ -run "TestUserProfileService|TestAiExportService" -v
go test ./internal/atlas/handler/ -run TestAiExportHandler -v
INTEGRATION_TESTS=1 go test ./internal/atlas/repository/postgres/ -run "TestUserProfileRepo|TestAiExportRepo|TestWave07Migration" -v
```

- [ ] **Step 3: Run codegen drift check**

```bash
bunx nx run api:codegen && bunx nx run graphql:validate
```

- [ ] **Step 4: Lint all changed packages**

```bash
bunx nx run api:lint
```

- [ ] **Step 5: Commit**

```bash
git add apps/api/internal/atlas/repository/postgres/user_profile_repo_test.go apps/api/internal/atlas/repository/postgres/ai_export_repo_test.go
git commit -m "test(wave-07): add integration tests for WAVE-07"
```