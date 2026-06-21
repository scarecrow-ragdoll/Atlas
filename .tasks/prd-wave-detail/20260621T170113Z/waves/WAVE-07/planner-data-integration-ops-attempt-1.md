# Planner Report: WAVE-07 Data Integration & Operations

**Run**: 20260621T170113Z  
**Role**: data-integration-ops  
**Source**: docs/prd-waves/waves/wave-07.md + product-verified docs + code evidence  
**Attempt**: 1

---

## 1. REST API Endpoint Design

### 1.1 POST /api/v1/ai-export — Generate Export

**Route**: `POST /api/v1/ai-export`  
**Auth**: Atlas PIN-guarded group (existing middleware).  
**Content-Type**: `application/json`  

**Request Body**:
```json
{
  "dateRangeStart": "2026-01-01",
  "dateRangeEnd": "2026-01-28",
  "includePhotos": false,
  "includeNutrition": true,
  "includeCardio": true,
  "includeMeasurements": true,
  "userComment": "Optional user note"
}
```

**Fields**:
| Field | Type | Required | Default | Notes |
|---|---|---|---|---|
| dateRangeStart | string (date) | yes | — | Inclusive |
| dateRangeEnd | string (date) | yes | — | Inclusive |
| includePhotos | bool | no | false | Matches PRD §17.3 opt-in default |
| includeNutrition | bool | no | true | |
| includeCardio | bool | no | true | |
| includeMeasurements | bool | no | true | |
| userComment | string | no | null | One-time comment per export |

**Validation**:
- `dateRangeStart <= dateRangeEnd`
- Range must not exceed 365 days (hard limit to bound ZIP size). Open question: is 365 days appropriate, or should this match the frontend's max (e.g. 52 weeks)?
- `includePhotos` defaults to false — explicit in PRD domain model invariant #10.

**Success Response (201)**:
```json
{
  "export": {
    "id": "uuid",
    "dateRangeStart": "2026-01-01",
    "dateRangeEnd": "2026-01-28",
    "includePhotos": false,
    "includeNutrition": true,
    "includeCardio": true,
    "includeMeasurements": true,
    "userComment": "Optional user note",
    "generatedPrompt": "...",
    "exportFilePath": "exports/uuid/export.zip",
    "createdAt": "2026-06-21T17:00:00Z"
  }
}
```

**Error responses**:
| Status | Code | When |
|---|---|---|
| 400 | INVALID_DATE_RANGE | Start > end or range > 365 days |
| 400 | NO_DATA_IN_RANGE | Zero entities found for the period |
| 500 | EXPORT_GENERATION_FAILED | ZIP write fails, DB write fails |

**Logic flow (handler → service → repo)**:

1. Handler reads userID from middleware context, parses/validates request body.
2. Calls `AiExportService.Generate(ctx, userID, input)`.
3. Service creates an `AiExportRecord` in state `draft` (no exportFilePath).
4. Service queries all relevant data sources for the date range:
   - DailyLog + WorkoutExercise + WorkoutSet (WAVE-01, WAVE-02)
   - CardioEntry (WAVE-03)
   - BodyWeightEntry (WAVE-03)
   - BodyCheckIn + BodyMeasurement (WAVE-03)
   - ProgressPhoto file paths (WAVE-03) — only when includePhotos=true
   - Nutrition products/templates/overrides (WAVE-06)
   - WeekFlag (WAVE-04)
   - UserProfile for goal/context
5. Service composes `manifest.json`, `data.json`, `summary.md`, CSV files in memory.
6. If includePhotos=true, copies photo files into `photos/` subdirectory inside ZIP.
7. Writes ZIP to `{exportBasePath}/{exportID}/export.zip`.
8. Generates the prompt string using:
   - Persistent AI context from UserProfile
   - One-time userComment
   - Data summary (date range, entity counts, week flags)
   - Section toggle flags
9. Updates AiExport record to state `generated` (sets exportFilePath, generatedPrompt).
10. Returns the now-full AiExport model.

**Data source queries**: Service assembles all data. No new sqlc queries needed for reading — all entities already have read queries from prior waves. The service needs to call the existing repositories directly or compose through existing service interfaces.

### 1.2 GET /api/v1/ai-export/download — Download ZIP

**Route**: `GET /api/v1/ai-export/download?exportId={uuid}`  
**Auth**: Atlas PIN-guarded group  

**Query params**: `exportId` (required, UUID)

**Success (200)**: Streams the ZIP file with:
```
Content-Type: application/zip
Content-Disposition: attachment; filename="ai-export-{dateRangeStart}-{dateRangeEnd}.zip"
```

**Error responses**:
| Status | Code | When |
|---|---|---|
| 400 | MISSING_EXPORT_ID | exportId param absent |
| 404 | EXPORT_NOT_FOUND | No AiExport with that ID or not owned by user |
| 404 | EXPORT_FILE_NOT_FOUND | Record exists but file missing from disk |
| 500 | INTERNAL_ERROR | File stat/read fails |

**Logic**:
1. Resolve AiExport record from DB by ID + userID.
2. Verify `exportFilePath` is non-empty (state = generated).
3. Open file from disk, stat it, stream with `Content-Disposition: attachment`.

### 1.3 GET /api/v1/user-profile — Read User Profile

**Route**: `GET /api/v1/user-profile`  
**Auth**: Atlas PIN-guarded group  

**Success (200)**:
```json
{
  "userProfile": {
    "id": "uuid",
    "displayName": "User",
    "goal": "Build muscle",
    "height": 180.0,
    "birthDate": "1990-01-01",
    "trainingExperience": "intermediate",
    "currentTrainingSplit": "Push/Pull/Legs",
    "preferredProgressionStyle": "double progression",
    "nutritionStrategy": "high protein",
    "persistentAiContext": "I respond well to high volume..."
  }
}
```

**Error responses**:
| Status | Code | When |
|---|---|---|
| 404 | PROFILE_NOT_FOUND | No user profile record exists for this user |
| 500 | INTERNAL_ERROR | DB failure |

**Notes**: All fields except `id` and `displayName` are optional per domain model. This is a new entity — see section 2.2 for the migration.

---

## 2. Data Lifecycle

### 2.1 AiExport Lifecycle

```
[draft] → (ZIP generation succeeds) → [generated]
[draft] → (ZIP generation fails)     → [draft] (orphan record, file never written)
[generated] → (cleanup TTL expires)  → [deleted]
```

**States per domain-model.md §AiExport**:
- **draft**: `exportFilePath` is NULL. AiExport record exists but no ZIP on disk. Generated when POST is accepted but before ZIP write completes. If generation fails, record stays in draft — treat as orphan.
- **generated**: `exportFilePath` is set. Record + ZIP are complete and downloadable.

**Cleanup**: A scheduled or on-demand cleanup deletes AiExport records and their ZIP files older than a configurable TTL. Default: 7 days. See section 3 for config.

### 2.2 New Database Migrations Required

#### Migration 00091_user_profiles.sql

Adds `user_profiles` table matching the domain model `UserProfile` entity:

```sql
CREATE TABLE user_profiles (
    id                        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                   UUID NOT NULL REFERENCES atlas_users(id),
    display_name              TEXT NOT NULL DEFAULT '',
    goal                      TEXT,
    height                    REAL,
    birth_date                DATE,
    training_experience       TEXT,
    current_training_split    TEXT,
    preferred_progression_style TEXT,
    nutrition_strategy        TEXT,
    persistent_ai_context     TEXT,
    created_at                TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id)
);

CREATE INDEX idx_user_profiles_user_id ON user_profiles (user_id);
```

#### Migration 00092_ai_exports.sql

```sql
CREATE TABLE ai_exports (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           UUID NOT NULL REFERENCES atlas_users(id),
    date_range_start  DATE NOT NULL,
    date_range_end    DATE NOT NULL,
    include_photos    BOOLEAN NOT NULL DEFAULT false,
    include_nutrition BOOLEAN NOT NULL DEFAULT true,
    include_cardio    BOOLEAN NOT NULL DEFAULT true,
    include_measurements BOOLEAN NOT NULL DEFAULT true,
    user_comment      TEXT,
    generated_prompt  TEXT,
    export_file_path  TEXT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_ai_exports_user_id ON ai_exports (user_id);
CREATE INDEX idx_ai_exports_created_at ON ai_exports (created_at);
```

### 2.3 New sqlc Queries

**queries/ai_exports.sql**:
```sql
-- name: CreateAiExport :one
INSERT INTO ai_exports (user_id, date_range_start, date_range_end, include_photos,
    include_nutrition, include_cardio, include_measurements, user_comment)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, user_id, date_range_start, date_range_end, include_photos,
    include_nutrition, include_cardio, include_measurements, user_comment,
    generated_prompt, export_file_path, created_at;

-- name: GetAiExportByID :one
SELECT id, user_id, date_range_start, date_range_end, include_photos,
    include_nutrition, include_cardio, include_measurements, user_comment,
    generated_prompt, export_file_path, created_at
FROM ai_exports
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: UpdateAiExportGenerated :one
UPDATE ai_exports
SET generated_prompt = $3, export_file_path = $4
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, date_range_start, date_range_end, include_photos,
    include_nutrition, include_cardio, include_measurements, user_comment,
    generated_prompt, export_file_path, created_at;

-- name: ListAiExportsByUser :many
SELECT id, user_id, date_range_start, date_range_end, include_photos,
    include_nutrition, include_cardio, include_measurements, user_comment,
    generated_prompt, export_file_path, created_at
FROM ai_exports
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: DeleteAiExport :one
DELETE FROM ai_exports
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, date_range_start, date_range_end, include_photos,
    include_nutrition, include_cardio, include_measurements, user_comment,
    generated_prompt, export_file_path, created_at;

-- name: ListStaleAiExports :many
SELECT id, user_id, date_range_start, date_range_end, include_photos,
    include_nutrition, include_cardio, include_measurements, user_comment,
    generated_prompt, export_file_path, created_at
FROM ai_exports
WHERE created_at < $1
  AND export_file_path IS NOT NULL
ORDER BY created_at ASC;
```

**queries/user_profiles.sql**:
```sql
-- name: GetUserProfileByUserID :one
SELECT id, user_id, display_name, goal, height, birth_date,
    training_experience, current_training_split, preferred_progression_style,
    nutrition_strategy, persistent_ai_context, created_at, updated_at
FROM user_profiles
WHERE user_id = $1
LIMIT 1;

-- name: UpsertUserProfile :one
INSERT INTO user_profiles (user_id, display_name, goal, height, birth_date,
    training_experience, current_training_split, preferred_progression_style,
    nutrition_strategy, persistent_ai_context)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (user_id)
DO UPDATE SET
    display_name = COALESCE($2, user_profiles.display_name),
    goal = COALESCE($3, user_profiles.goal),
    height = COALESCE($4, user_profiles.height),
    birth_date = COALESCE($5, user_profiles.birth_date),
    training_experience = COALESCE($6, user_profiles.training_experience),
    current_training_split = COALESCE($7, user_profiles.current_training_split),
    preferred_progression_style = COALESCE($8, user_profiles.preferred_progression_style),
    nutrition_strategy = COALESCE($9, user_profiles.nutrition_strategy),
    persistent_ai_context = COALESCE($10, user_profiles.persistent_ai_context),
    updated_at = now()
RETURNING id, user_id, display_name, goal, height, birth_date,
    training_experience, current_training_split, preferred_progression_style,
    nutrition_strategy, persistent_ai_context, created_at, updated_at;
```

### 2.4 Bootstrap: Default UserProfile Record

When the bootstrap service ensures default user (WAVE-01 pattern), it must also upsert a default `user_profiles` record so `GET /api/v1/user-profile` never 404s for the default user. The default display_name should match the `atlas_users.display_name`.

### 2.5 ZIP File Lifecycle

- **Path pattern**: `{export_base_path}/{export_uuid}/export.zip`
- `export_base_path` defaults to `./data/exports` (new config entry, see section 3).
- **TTL**: 7 days after creation. Cleanup is handled by a periodic task — see section 3.
- **Cleanup**: Delete the directory `{export_base_path}/{export_uuid}/` and set `export_file_path = NULL` on the AiExport record (transition to draft, or delete the record — choose: hard-delete for simplicity since re-download of expired exports should not be possible). **Decision**: hard-delete both file and record.

### 2.6 Rollback

No data migration — additive only. Reverting WAVE-07 means:
- Drop `ai_exports` and `user_profiles` tables (migration down).
- Remove config entries from `config.yml`.
- Remove REST route registrations from `main.go`.
- ZIP file cleanup is manual (delete `./data/exports/`).

---

## 3. Configuration Additions

### config.yml
```yaml
ai_export:
  base_path: ./data/exports
  max_range_days: 365
  max_photos_in_export: 20
```

### appconfig/config.go Additions
```go
type AiExportConfig struct {
    BasePath       string `mapstructure:"base_path"`
    MaxRangeDays   int    `mapstructure:"max_range_days"`
    MaxPhotosCount int    `mapstructure:"max_photos_in_export"`
}
```

Add to `Config` struct:
```go
AiExport AiExportConfig `mapstructure:"ai_export"`
```

---

## 4. ZIP Format Specification

### 4.1 Top-Level Structure
```
export.zip
├── manifest.json
├── data.json
├── summary.md
├── workouts.csv
├── measurements.csv
├── nutrition.csv
├── cardio.csv
└── photos/
    ├── photo1.jpg
    └── photo2.png
```

### 4.2 manifest.json Schema
```json
{
  "schemaVersion": "1.0.0",
  "exportTimestamp": "2026-06-21T17:00:00Z",
  "dateRangeStart": "2026-01-01",
  "dateRangeEnd": "2026-01-28",
  "sections": {
    "workouts": true,
    "cardio": true,
    "bodyWeight": true,
    "measurements": true,
    "nutrition": true,
    "photos": false
  }
}
```

### 4.3 data.json Schema

Contains all entities for the selected date range as a flat JSON object with named arrays:

```json
{
  "workouts": [
    {
      "date": "2026-01-01",
      "notes": "Leg day",
      "bodyWeight": 80.5,
      "exercises": [
        {
          "exerciseName": "Squat",
          "exerciseId": "uuid",
          "order": 1,
          "workingWeightSnapshot": 100.0,
          "notes": "Felt strong",
          "sets": [
            {
              "setNumber": 1,
              "weight": 100.0,
              "reps": 8,
              "rpe": 8,
              "rir": null,
              "notes": null
            }
          ]
        }
      ]
    }
  ],
  "cardio": [ /* CardioEntry fields */ ],
  "bodyWeightEntries": [ /* BodyWeightEntry fields */ ],
  "measurements": [ /* BodyMeasurement with check-in date */ ],
  "nutrition": {
    "products": [ /* NutritionProduct fields */ ],
    "template": { /* NutritionTemplate with items */ },
    "overrides": [ /* DailyNutritionOverride with items */ ]
  },
  "weekFlags": [ /* WeekFlag fields */ ],
  "userProfile": {
    "goal": "...",
    "height": 180.0,
    "trainingExperience": "intermediate",
    "currentTrainingSplit": "Push/Pull/Legs",
    "persistentAiContext": "..."
  }
}
```

### 4.4 summary.md Format

```markdown
# AI Export Summary
**Period**: 2026-01-01 → 2026-01-28
**Generated**: 2026-06-21T17:00:00Z

## Overview
- Workout days: 16
- Exercises performed: 8
- Total sets: 120
- Cardio sessions: 4
- Weight entries: 3
- Measurements: 5
- Nutrition tracked: 7 days
- Week flags: 2

## Goal
Build muscle

## Week Flags
- Week of 2026-01-06: HIGH_STRESS, POOR_SLEEP
- Week of 2026-01-13: TRAVEL

## Comment
Optional user note about this period.
```

### 4.5 CSV Column Layouts

**workouts.csv**:
```
date,exercise_name,set_number,weight,reps,rpe,rir,set_notes,exercise_notes,day_notes
2026-01-01,Squat,1,100.0,8,8,,Felt strong,Leg day
2026-01-01,Squat,2,100.0,7,9,,,Leg day
```

**cardio.csv**:
```
date,type,duration_minutes,avg_pulse,heart_rate_zone,notes
2026-01-02,running,30,145,4,
```

**measurements.csv**:
```
check_in_date,measurement_type,side,value,notes
2026-01-06,waist,,82.0,
2026-01-06,biceps,left,36.5,
2026-01-06,biceps,right,37.0,
```

**nutrition.csv**:
```
date,product_name,amount_grams,calories,protein,fat,carbs,meal_label,operation
2026-01-06,Whey Protein,30,115,25,1,3,breakfast,
2026-01-06,Chicken Breast,200,330,62,7,0,lunch,add
```

### 4.6 Photos

When `includePhotos=true`, photos are included as files in `photos/` subdirectory. The files are copied from their original storage paths (from `ProgressPhoto.filePath`). Photos are renamed to a flat naming scheme: `{checkInId}_{angle}.{ext}` to avoid collisions and preserve context.

---

## 5. Log Markers

| Marker | Location | When |
|---|---|---|
| `[AiExport][generate][BLOCK_EXPORT_START]` | Handler/service | POST received, validation passed |
| `[AiExport][generate][BLOCK_EXPORT_DATA_QUERY]` | Service | Querying data sources |
| `[AiExport][generate][BLOCK_EXPORT_ZIP_BUILD]` | Service | Building ZIP in-memory |
| `[AiExport][generate][BLOCK_EXPORT_ZIP_WRITE]` | Service | Writing ZIP to disk |
| `[AiExport][generate][BLOCK_EXPORT_PROMPT_GENERATE]` | Service | Building prompt string |
| `[AiExport][generate][BLOCK_EXPORT_DB_SAVE]` | Service | Saving AiExport record |
| `[AiExport][generate][BLOCK_EXPORT_SUCCESS]` | Service | Export complete |
| `[AiExport][generate][BLOCK_EXPORT_NO_DATA]` | Service | Zero entities in range |
| `[AiExport][generate][BLOCK_EXPORT_FAILURE]` | Service | Any error during generation |
| `[AiExport][download]` | Handler | ZIP downloaded |
| `[AiExport][download][BLOCK_EXPORT_NOT_FOUND]` | Handler | ID not found |
| `[AiExport][download][BLOCK_EXPORT_FILE_MISSING]` | Handler | Record OK but file absent |
| `[AiExport][cleanup]` | Cleanup task | Running cleanup |
| `[AiExport][cleanup][BLOCK_EXPORT_DELETED]` | Cleanup task | File + DB removed |
| `[AiExport][cleanup][BLOCK_EXPORT_FILE_DELETE_FAILED]` | Cleanup task | File delete error |

---

## 6. Rollout / Rollback

### 6.1 Rollout Order
1. Add `AiExportConfig` to `appconfig` and `config.yml`.
2. Create `00091_user_profiles.sql` and `00092_ai_exports.sql` migrations.
3. Add sqlc queries for `ai_exports` and `user_profiles`.
4. Run `sqlc generate`.
5. Create `AiExportRecord` and `UserProfileRecord` (DB) models + `AiExport` and `UserProfile` (public) models in `apps/api/internal/atlas/models/`.
6. Create `AiExportRepository` and `UserProfileRepository` in `apps/api/internal/atlas/repository/postgres/`.
7. Update bootstrap service to upsert default user_profile.
8. Implement ZIP generation utility (`apps/api/internal/atlas/export/zipper.go` or similar).
9. Implement `AiExportService` in `apps/api/internal/atlas/service/`:
   - `Generate(ctx, userID, input)` — orchestrates data gathering, ZIP build, prompt generation, DB save.
   - `GetByID(ctx, userID, exportID)` — for download.
   - `ListStaleExports(ctx, cutoff)` — for cleanup.
   - `DeleteExport(ctx, userID, exportID)` — manual delete (optional).
10. Implement `UserProfileService` with `GetProfile(ctx, userID)`.
11. Create REST handlers for the three endpoints.
12. Register routes in `main.go` (under the existing atlas PIN-guarded group).
13. Create the exports directory at startup or on first write.

### 6.2 Rollback
- Run `goose down` for the two new migrations (or one combined migration step).
- Remove route registrations from `main.go`.
- Remove new handler files, service files, repo files.
- Remove config entries.
- Clean up exported ZIP files: `rm -rf ./data/exports/`.

---

## 7. Risks & Concerns

### EDGE-024: Disk Space
**Mitigation**: Limit ZIP lifetime (7 days). Log and monitor disk usage close to the exports directory. The `max_range_days: 365` bound limits the maximum ZIP size. If photos are included, the per-export max photo count (`max_photos_in_export: 20`) prevents unreasonably large archives.

### Performance: Large Date Range with Photos
**Mitigation**: Build ZIP in-memory using `archive/zip` buffer, then write to disk in one pass. For very large exports (many photos), consider streaming directly to a file without holding the full ZIP in memory. The `max_photos_in_export` limit caps the worst-case size.

### Photo File Copy in ZIP
Photos are stored in the media directory (`./data/media`). The export service must copy them into the ZIP — not move them. This adds I/O but preserves the originals. If photos are large, this is the dominant cost in ZIP generation.

### Cleanup Safety
The cleanup task must never delete AiExport records without first confirming the file deletion succeeded. If file deletion fails, log the error and skip that record to avoid dangling files. Re-run on the next cycle.

### No Concurrent Generation Guard
If the user hits "Generate" twice in rapid succession, two ZIP files will be generated for the same range. This is acceptable (each gets its own UUID). Consider frontend debouncing as the primary guard, not a backend idempotency key.

---

## 8. Open Questions

| ID | Question | Why It Matters | Resolution |
|---|---|---|---|
| Q-W07-DIO-01 | Is 365 days the correct max range, or should it match the frontend's 52-week default ceiling? | Hard limit on ZIP size and query scope. | Needs product alignment. Default to 365 days for now. |
| Q-W07-DIO-02 | Should the cleanup TTL be configurable at runtime or only via config file? | Operational flexibility vs. simplicity. | Config file is sufficient for MVP. |
| Q-W07-DIO-03 | What is the max ZIP size we should tolerate before streaming to disk? | In-memory ZIP buffer may OOM for large exports with 20 photos at 25MB each = ~500MB. | Set a threshold: if estimated size > 100MB, stream to temp file; else build in-memory. |
| Q-W07-DIO-04 | Should photos be downscaled/resized when included in the export ZIP? | Full-resolution check-in photos may be 25MB each. 20 photos = 500MB. | For MVP, include originals. Downscaling is a future optimization. |
| Q-W07-DIO-05 | Does `GET /api/user-profile` need to support updating the profile in this wave, or is it read-only? | Frontend page-009 says "User goal display" — only read is needed. Separate create/update endpoint can be added in a later wave. |
| Q-W07-DIO-06 | Should the bootstrap service create a default UserProfile record with display_name populated from atlas_users? | Prevents 404 on first GET. | Yes — add to bootstrap. |
| Q-W07-DIO-07 | What is the exact prompt template format? | The prompt needs to be a useful ChatGPT input. The format should include goal, context, data summary, flags, comment, and instructions. | MVP: compose a structured plain-text prompt with sections. Future: make prompt template configurable. |

---

## 9. Traceability

| Artifact | Source |
|---|---|
| AiExport entity | domain-model.md:90-92, §AiExport |
| UserProfile entity | domain-model.md:36-37, §UserProfile |
| AiExport lifecycle states | domain-model.md:116, §AiExport |
| Section toggles | functional-spec.md:82-83, §AI Export |
| ZIP output format | functional-spec.md:83, page-009.md:14 |
| Include photos default false | domain-model.md:135, invariant #10 |
| Prompt builder | functional-spec.md:84-86 |
| Week flags | functional-spec.md:86, WAVE-04 |
| REST endpoints | page-009.md:39-43 |
| ZIP storage path pattern | Derived from media handler `basePath` pattern (atlas_media.go:153) |

---

## 10. New Files Required

| File | Purpose |
|---|---|
| `apps/api/internal/repository/postgres/migrations/00091_user_profiles.sql` | UserProfile table |
| `apps/api/internal/repository/postgres/migrations/00092_ai_exports.sql` | AiExport table |
| `apps/api/internal/repository/postgres/queries/user_profiles.sql` | sqlc queries for user_profiles |
| `apps/api/internal/repository/postgres/queries/ai_exports.sql` | sqlc queries for ai_exports |
| `apps/api/internal/atlas/models/user_profile.go` | UserProfile domain models |
| `apps/api/internal/atlas/models/ai_export.go` | AiExport domain models |
| `apps/api/internal/atlas/repository/postgres/user_profile_repo.go` | UserProfileRepository |
| `apps/api/internal/atlas/repository/postgres/ai_export_repo.go` | AiExportRepository |
| `apps/api/internal/atlas/service/user_profile_service.go` | UserProfileService |
| `apps/api/internal/atlas/service/ai_export_service.go` | AiExportService |
| `apps/api/internal/handler/ai_export_handler.go` | REST handler for ai-export endpoints |
| `apps/api/internal/handler/user_profile_handler.go` | REST handler for user-profile endpoint |
| `apps/api/internal/atlas/export/zipper.go` | ZIP building + prompt generation utility |

---

## 11. Acceptance Criteria (Data & Integration Scope)

- AC-W07-DIO-01: `POST /api/v1/ai-export` returns 201 with export ID when given a valid date range with data
- AC-W07-DIO-02: `POST /api/v1/ai-export` returns 400 when date range exceeds 365 days
- AC-W07-DIO-03: `POST /api/v1/ai-export` returns 400 when start > end
- AC-W07-DIO-04: `POST /api/v1/ai-export` returns 400 when no entities exist in the date range
- AC-W07-DIO-05: Generated ZIP contains manifest.json with valid schema version, export timestamp, date range, sections
- AC-W07-DIO-06: Generated ZIP contains data.json with workouts, exercises, sets, cardio, body weight, measurements, nutrition, week flags, user profile
- AC-W07-DIO-07: Generated ZIP contains summary.md with human-readable overview
- AC-W07-DIO-08: Generated ZIP contains workouts.csv, cardio.csv, measurements.csv, nutrition.csv with correct column headers
- AC-W07-DIO-09: When includePhotos=false, ZIP does NOT contain photos/ directory
- AC-W07-DIO-10: When includePhotos=true, ZIP contains photos/ with copies of ProgressPhoto files
- AC-W07-DIO-11: When includeNutrition=false, data.json omits nutrition section and nutrition.csv is absent
- AC-W07-DIO-12: When includeCardio=false, data.json omits cardio section and cardio.csv is absent
- AC-W07-DIO-13: When includeMeasurements=false, data.json omits measurements section and measurements.csv is absent
- AC-W07-DIO-14: `GET /api/v1/ai-export/download?exportId={id}` returns ZIP with correct Content-Type and Content-Disposition
- AC-W07-DIO-15: `GET /api/v1/ai-export/download?exportId={id}` returns 404 for non-existent export
- AC-W07-DIO-16: `GET /api/v1/ai-export/download?exportId={id}` returns 404 for export owned by different user (user-scoped isolation)
- AC-W07-DIO-17: `GET /api/v1/user-profile` returns user profile with displayName, goal, persistentAiContext
- AC-W07-DIO-18: `GET /api/v1/user-profile` returns profile after bootstrap without a separate create step
- AC-W07-DIO-19: AiExport record is saved to DB with generated_prompt and export_file_path after successful ZIP generation
- AC-W07-DIO-20: AiExport record has `export_file_path = NULL` if ZIP generation fails
- AC-W07-DIO-21: generated_prompt contains userComment when provided
- AC-W07-DIO-22: generated_prompt contains persistentAiContext from user profile
- AC-W07-DIO-23: generated_prompt lists week flags for the period
- AC-W07-DIO-24: Export ZIP is stored at `{export_base_path}/{export_uuid}/export.zip`

## 12. Exit Criteria (Data & Integration Scope)

- EC-W07-DIO-01: All ACs passing with integration tests that write real ZIP to temp directory and verify contents
- EC-W07-DIO-02: `sqlc generate` succeeds for the two new query files
- EC-W07-DIO-03: Goose migration `00091` and `00092` apply and roll back cleanly
- EC-W07-DIO-04: All log markers present in handler and service code
- EC-W07-DIO-05: Existing tests pass (`bun run test` at the Go API package level)
- EC-W07-DIO-06: Lint passes (`bun run lint`)

## 13. Verification Obligations

| Verification | Type | When |
|---|---|---|
| ZIP content integrity tests | Unit | After zipper.go |
| Handler integration tests with temp storage | Integration | After ai_export_handler.go |
| Repository tests with sqlc-generated code | Unit | After repo files |
| Service tests with mock repo | Unit | After service files |
| Config parse + validation tests | Unit | After appconfig changes |
| Migration apply/down test | Script | After migration files |