# Wave 07: AI Export and Prompt Builder

## Status
ready-for-dev

## User Approval
Source wave: user-approved (2026-06-18). All questions resolved and design decisions user-approved (2026-06-21).

## Source Wave Summary
Generate AI-ready exports with structured data and prompts for ChatGPT analysis. Backend delivers prompt builder with period selection, ZIP export with manifest.json/data.json/summary.md/CSV files, one-time comment support, section toggles (photos optional), and UserProfile for persistent AI context and goal storage. CAP-W07-003 (week flags CRUD) removed from scope — WAVE-04 owns it; WAVE-07 reads week flags via WAVE-04 service.

## Outcome After Implementation
- OUT-W07-001 Prompt builder with period selection — POST /api/ai-export/generate accepts dateRangeStart, dateRangeEnd, section toggles, optional userComment; returns generated prompt string and export record
- OUT-W07-002 ZIP export with manifest.json, data.json, summary.md, CSV — Generation service writes ZIP to {ExportBasePath}/{userId}/{exportId}.zip, returns export ID and download URL
- OUT-W07-003 Week flags support — WAVE-07 reads week flags from WAVE-04 WeekFlagService; included in prompt and data.json. Week flag CRUD is WAVE-04 scope, not WAVE-07.
- OUT-W07-004 One-time comment support — userComment stored on AiExport record, included in data.json, summary.md, and manifest.json
- OUT-W07-005 Section toggles (photos optional) — Input fields: includePhotos (default false), includeNutrition (default true), includeCardio (default true), includeMeasurements (default true). Photos opt-in enforced server-side per RULE-025.

## Scope Included
- CAP-W07-001 Persistent AI context — UserProfile table with goal, height, birthDate, trainingExperience, currentTrainingSplit, preferredProgressionStyle, nutritionStrategy, persistentAiContext
- CAP-W07-002 User goal storage — goal field on UserProfile
- CAP-W07-004 Prompt generation — prompt builder service reads UserProfile context + week flags + userComment + period data; produces plain-text prompt
- CAP-W07-005 AI export ZIP creation — ZIP assembly service writes manifest.json, data.json, summary.md, CSVs, optional photos/
- CAP-W07-006 manifest.json with export metadata — schemaVersion=1, exportTimestamp, dateRange, includedSections
- CAP-W07-007 data.json with all entities for period — workouts (with exercises, sets, RPE/RIR), cardio, body weight, measurements, nutrition, week flags, user profile
- CAP-W07-008 summary.md with human-readable overview — period summary, goal, workout stats, weight/measurement deltas, nutrition summary, cardio, comments
- CAP-W07-009 CSV files for compatibility — workouts.csv, measurements.csv, nutrition.csv, cardio.csv

## Scope Excluded
- CAP-W07-003 Week flags CRUD — REMOVED from WAVE-07. WAVE-04 owns WeekFlag CRUD. WAVE-07 reads week flags via WAVE-04 WeekFlagService.
- Direct ChatGPT/OpenAI API call — explicitly excluded per RULE-029
- AI review history (AiReview) — belongs to WAVE-08
- Backup/restore ZIP import — belongs to WAVE-09
- Frontend pages, routes, UX states — PAGE-009 delivers frontend
- Frontend ZIP download UX — frontend calls GET /api/ai-export/download?exportId=
- Photo resizing/compression for export — not specified; include originals in MVP
- Watermarking or metadata stripping on exported photos — not specified
- Async ZIP generation — sync for MVP; architecture supports future async

## Dependencies And Other-Wave Fit

### Prior Wave Compatibility
- **WAVE-01 (Foundation)** — Required. PIN auth middleware for all endpoints. Settings for defaultAiExportWeeks. atlas_users for user identity.
- **WAVE-02 (Exercise Library)** — Exercise metadata consumed via service layer (read-only). Compatible.
- **WAVE-03 (Workout Diary)** — Workout data consumed if available. Empty arrays when WAVE-03 not deployed (stub pattern from WAVE-06 DDEC-W06-010). Compatible.
- **WAVE-04 (Cardio and Body Tracking)** — CardioEntry, BodyWeightEntry, BodyCheckIn, BodyMeasurement, ProgressPhoto, WeekFlag consumed via service layer (read-only). Compatible. Week flag browsing on PAGE-009 uses WAVE-04 GraphQL weekFlags query directly.
- **WAVE-05 (Nutrition)** — NutritionProduct, NutritionTemplate, NutritionTemplateItem, DailyNutritionOverride consumed via service layer (read-only). Compatible.
- **WAVE-06 (Charts)** — Shares same underlying data. Compatible.

### Future Wave Compatibility
- **WAVE-08 (AI Review History)** — Depends on WAVE-07. WAVE-07 creates AiExport record with generatedPrompt. WAVE-08 creates independent AiReview record. Clean boundary.
- **WAVE-09 (Backup Import/Export)** — Similar ZIP pattern, different purpose. WAVE-07: per-period AI export ZIP at {ExportBasePath}/{userId}/{exportId}.zip. WAVE-09: full data backup with version manifest. No scope collision.

### Independent Deliverability
- Can be implemented without WAVE-03 (workout data returns empty arrays)
- Can be implemented without WAVE-04 (cardio/body/flags data returns empty)
- Can be implemented without WAVE-05 (nutrition data returns empty)
- **Cannot be implemented without WAVE-01** — hard dependency on PIN auth and user identity

## Frontend Pages Dependencies
PAGE-009 (AI Export) requires:
- `POST /api/ai-export/generate` — generates prompt + ZIP, returns generatedPrompt in response body
- `GET /api/ai-export/download?exportId=` — streams generated ZIP
- `GET /api/user-profile` — returns UserProfile fields (goal, persistentAiContext)
- Week flag browsing — uses WAVE-04 GraphQL `weekFlags(weekStartDate:)` query directly (no REST proxy needed)

Prompt display/copy: Backend returns generatedPrompt in POST /api/ai-export/generate response body so frontend can display without downloading ZIP.

## Codebase Fit And Touchpoints

### New Files Required
- `apps/api/internal/repository/postgres/migrations/00091_user_profiles.sql` — UserProfile table
- `apps/api/internal/repository/postgres/migrations/00092_ai_exports.sql` — AiExport table
- `apps/api/internal/repository/postgres/queries/user_profiles.sql` — sqlc queries for user_profiles
- `apps/api/internal/repository/postgres/queries/ai_exports.sql` — sqlc queries for ai_exports
- `apps/api/internal/atlas/models/user_profile.go` — UserProfile model types
- `apps/api/internal/atlas/models/ai_export.go` — AiExport model types
- `apps/api/internal/atlas/repository/postgres/user_profile_repo.go` — UserProfileRepository
- `apps/api/internal/atlas/repository/postgres/ai_export_repo.go` — AiExportRepository
- `apps/api/internal/atlas/service/user_profile_service.go` — UserProfileService
- `apps/api/internal/atlas/service/ai_export_service.go` — AiExportService (prompt generation + data aggregation)
- `apps/api/internal/atlas/service/export_zip.go` — ZIP generation utilities
- `apps/api/internal/atlas/graph/schema/user_profile.graphql` — UserProfile GraphQL schema
- `apps/api/internal/atlas/graph/schema/ai_export.graphql` — AiExport GraphQL schema
- `apps/api/internal/atlas/graph/resolver/user_profile.go` — UserProfile resolvers
- `apps/api/internal/atlas/graph/resolver/ai_export.go` — AiExport resolvers
- `apps/api/internal/handler/ai_export_handler.go` — REST download handler
- `apps/api/internal/handler/user_profile_handler.go` — REST user-profile handler

### Existing Files to Modify
- `apps/api/internal/atlas/graph/resolver/resolver.go` — add UserProfileService, AiExportService fields
- `apps/api/internal/atlas/graph/schema/schema.graphql` — add type extensions
- `apps/api/cmd/server/main.go` — wire new repositories, services, resolvers, routes
- `apps/api/atlas-gqlgen.yml` — add bindings for all new types
- `apps/api/internal/atlas/service/atlas_bootstrap_service.go` — create default UserProfile during bootstrap
- `apps/api/internal/appconfig/config.go` — add AiExportConfig
- `apps/api/config.yml` — add ai_export config section

### Existing Code Already Providing WAVE-07 Scope
- WeekFlag model + enum + validation — complete (WAVE-04)
- WeekFlag repository (CRUD) — complete (WAVE-04)
- WeekFlag service (create/list/delete) — complete (WAVE-04)
- WeekFlag GraphQL resolvers — complete (WAVE-04)
- WeekFlag GraphQL schema — complete (WAVE-04)
- WeekFlag DB migration 00089 — complete (WAVE-04)
- Settings with defaultAiExportWeeks — complete (WAVE-01)

### Patterns to Follow
- Models: UserProfileRecord/UserProfile/UserProfileInput triple matching week_flag.go pattern
- Repository: Interface + private struct + New*Repository(pool) + *FromRow() matching week_flag_repo.go
- Service: Interface + private struct + constructor + sentinel errors + FromRecord() matching week_flag.go
- Resolvers: middleware.GetAtlasUserID + union result types matching resolver/week_flag.go
- GraphQL schema: Types with nullable fields + input types + result types with inline errors matching week_flag.graphql
- REST handler: Binary download via http.ServeContent matching progress_photo_handler.go
- AiExportDataProvider: Interface wrapping all data source dependencies (keeps AiExportService constructor manageable)

## Design Contracts

### DDEC-W07-001: UserProfile as Separate Entity
Create separate `user_profiles` table (NOT extending `atlas_users` or `settings`). UserProfile entity: goal, height, birthDate, trainingExperience, currentTrainingSplit, preferredProgressionStyle, nutritionStrategy, persistentAiContext. Domain-model.md defines UserProfile separately from Settings. Settings exists for app configuration (PIN, units, defaultAiExportWeeks); UserProfile exists for user-specific data.

### DDEC-W07-002: Week Flags Read-Only
CAP-W07-003 removed from WAVE-07 scope. WAVE-04 owns WeekFlag CRUD. WAVE-07 reads week flags via WAVE-04 WeekFlagService. No week flag create/update/delete in WAVE-07. PAGE-009 uses WAVE-04 GraphQL weekFlags query for browsing.

### DDEC-W07-003: REST Endpoint Design
POST /api/ai-export/generate — generates prompt + ZIP, returns generatedPrompt in response body
GET /api/ai-export/download?exportId= — streams ZIP file
GET /api/user-profile — returns UserProfile fields
All endpoints use /api/ prefix (no /api/v1/). GraphQL used for AiExport query/mutation CRUD operations.

### DDEC-W07-004: include_photos Default False
include_photos defaults to false. Enforced server-side per RULE-025 and domain model invariant #10. AiExport migration DDL: `include_photos BOOLEAN NOT NULL DEFAULT false`.

### DDEC-W07-005: Migration Numbers
00091_user_profiles.sql, 00092_ai_exports.sql. Based on actual latest migration 00090_nutrition_tables.sql.

### DDEC-W07-006: Storage Path
{ExportBasePath}/{userId}/{exportId}.zip. User-scoped path for future multi-user isolation. ExportBasePath defaults to ./data/exports.

### DDEC-W07-007: Export Lifecycle
7-day TTL + delete-on-regeneration. Cleanup task deletes stale exports (file + DB record). On regeneration, old export for same user is deleted first. Hard-delete of both file and record.

### DDEC-W07-008: Temp-File-Atomic-Rename
ZIP generation writes to temp file first (.tmp-{uuid}.zip), then atomically renames to final path on success. On failure, temp file cleaned up. AiExport.exportFilePath only set after rename succeeds.

### DDEC-W07-009: Max Export Size
100MB uncompressed hard limit. Generation rejected with error if estimated data size exceeds limit. Configurable via max_export_size_bytes.

### DDEC-W07-010: Photo in Export
Files in photos/ subfolder within ZIP. Named as {checkInId}_{angle}.{ext}. Copied from media storage (not moved).

### DDEC-W07-011: ZIP Format
manifest.json (schemaVersion=1, exportTimestamp, dateRange, includedSections), data.json (all entities), summary.md (human-readable), workouts.csv, measurements.csv, nutrition.csv, cardio.csv, photos/ (when opted in).

### DDEC-W07-012: WAVE-03 Stub Pattern
Empty arrays for workout data when WAVE-03 not deployed (same stub pattern as WAVE-06 DDEC-W06-010). AiExportDataProvider.GetWorkoutSummary returns empty slices.

### DDEC-W07-013: Sync Generation
MVP uses sync generation. Architecture supports future async if needed for large exports.

### DDEC-W07-014: gqlgen Config
Add bindings for all new types: UserProfile, UserProfileInput, UserProfileResult, UserProfileValidationError, UserProfileNotFoundError, UserProfileAuthError, UserProfileErrorCode, AiExport, CreateAiExportInput, AiExportResult, AiExportsResult, AiExportValidationError, AiExportNotFoundError, AiExportAuthError, AiExportErrorCode.

### DDEC-W07-015: display_name
NOT in user_profiles. Use atlas_users.display_name. If needed in UserProfile API response, derive via query join.

## Data API Integration And Operations

### REST Endpoints

**POST /api/ai-export/generate** — Generate Export (prompt + ZIP)
- Auth: PIN-guarded middleware
- Content-Type: application/json
- Request: { dateRangeStart, dateRangeEnd, includePhotos (default false), includeNutrition (default true), includeCardio (default true), includeMeasurements (default true), userComment (optional) }
- Validation: dateRangeStart <= dateRangeEnd, range not exceeding 365 days
- Success 200: { export: { id, dateRangeStart, dateRangeEnd, includePhotos, includeNutrition, includeCardio, includeMeasurements, userComment, generatedPrompt, exportFilePath, createdAt } }
- Error 400: INVALID_DATE_RANGE, DATA_SIZE_EXCEEDS_LIMIT
- Error 500: EXPORT_GENERATION_FAILED

**GET /api/ai-export/download?exportId=** — Download ZIP
- Auth: PIN-guarded middleware
- Query: exportId (UUID, required)
- Success 200: application/zip with Content-Disposition: attachment; filename="ai-export-{dateRangeStart}-{dateRangeEnd}-{shortUuid}.zip"
- Error 400: MISSING_EXPORT_ID
- Error 404: EXPORT_NOT_FOUND, EXPORT_FILE_NOT_FOUND (record OK but file absent)
- Ownership check: validate AiExport.userId matches session user

**GET /api/user-profile** — Read User Profile
- Auth: PIN-guarded middleware
- Success 200: { id, goal, height, birthDate, trainingExperience, currentTrainingSplit, preferredProgressionStyle, nutritionStrategy, persistentAiContext }
- Error 404: PROFILE_NOT_FOUND (resolved by bootstrap default)

### ZIP Format Specification
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
    ├── {checkInId}_{angle}.{ext}
    └── ...
```

**manifest.json**: { schemaVersion: 1, exportTimestamp, dateRangeStart, dateRangeEnd, sections: { workouts, cardio, bodyWeight, measurements, nutrition, photos } }

**data.json**: { workouts[], cardio[], bodyWeightEntries[], measurements[], nutrition{products, template, overrides}, weekFlags[], userProfile{goal, height, trainingExperience, ...} }

**summary.md**: Markdown with period, goal, workout stats, weight/measurement changes, nutrition summary, cardio summary, week flags, user comment.

**CSV Columns**: workouts.csv (date, exercise_name, set_number, weight, reps, rpe, rir, set_notes, exercise_notes, day_notes), cardio.csv (date, type, duration_minutes, avg_pulse, heart_rate_zone, notes), measurements.csv (check_in_date, measurement_type, side, value, notes), nutrition.csv (date, product_name, amount_grams, calories, protein, fat, carbs, meal_label, operation).

### AiExport Lifecycle
draft (exportFilePath NULL) -> generated (exportFilePath set) -> deleted (7-day TTL or regeneration)
On generation failure: record stays in draft, temp file cleaned up. No partial export saved.

### Configuration Additions
```yaml
ai_export:
  base_path: ./data/exports
  max_range_days: 365
  max_photos_in_export: 20
  max_export_size_bytes: 104857600
```

### Log Markers
[AiExport][generate][BLOCK_EXPORT_START] — POST received, validation passed
[AiExport][generate][BLOCK_EXPORT_DATA_QUERY] — Querying data sources
[AiExport][generate][BLOCK_EXPORT_ZIP_BUILD] — Building ZIP in-memory
[AiExport][generate][BLOCK_EXPORT_ZIP_WRITE] — Writing ZIP to disk
[AiExport][generate][BLOCK_EXPORT_PROMPT_GENERATE] — Building prompt string
[AiExport][generate][BLOCK_EXPORT_DB_SAVE] — Saving AiExport record
[AiExport][generate][BLOCK_EXPORT_SUCCESS] — Export complete
[AiExport][generate][BLOCK_EXPORT_NO_DATA] — Zero entities in range
[AiExport][generate][BLOCK_EXPORT_FAILURE] — Any error during generation
[AiExport][download] — ZIP downloaded
[AiExport][download][BLOCK_EXPORT_NOT_FOUND] — ID not found
[AiExport][download][BLOCK_EXPORT_FILE_MISSING] — Record OK but file absent
[AiExport][cleanup] — Running cleanup
[AiExport][cleanup][BLOCK_EXPORT_DELETED] — File + DB removed
[AiExport][cleanup][BLOCK_EXPORT_FILE_DELETE_FAILED] — File delete error
[UserProfile][update] — Profile updated

No prompt content, user comments, body values, photo paths, or week flag notes in logs.

## Security Privacy And Compliance
- All three endpoints under PIN-guarded middleware (consistent with WAVE-01 auth pattern)
- Ownership validation on download endpoint — 404 on mismatch
- includePhotos defaults to false enforced at DB + service layer (defense-in-depth)
- Temp-file-atomic-rename for ZIP generation (EDGE-024)
- UUID-based export filenames to prevent enumeration
- No export content in application logs (AC-118)
- 100MB export size hard limit prevents disk exhaustion
- User-scoped storage paths for future multi-user isolation
- No external data transmission (RULE-027 — manual copy-paste)
- Each generation creates fresh data snapshot — no stale data leakage
- Cleanup: file deletion before record deletion (prevents dangling files)

## Implementation Slices

### SLICE-W07-001: UserProfile Migration (00091_user_profiles.sql)
CREATE TABLE user_profiles (id UUID PK, user_id UUID FK UNIQUE, goal TEXT, height REAL, birth_date DATE, training_experience TEXT, current_training_split TEXT, preferred_progression_style TEXT, nutrition_strategy TEXT, persistent_ai_context TEXT, created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ). All fields nullable except id, user_id, created_at, updated_at. Down: DROP TABLE.

### SLICE-W07-002: UserProfile Model (models/user_profile.go)
UserProfileRecord (DB row), UserProfile (public JSON), UserProfileInput (mutation input). Result/error types matching WeekFlag pattern. All nullable fields use pointer types.

### SLICE-W07-003: UserProfile SQLc Queries (queries/user_profiles.sql)
GetUserProfileByUserID (SELECT by user_id), UpsertUserProfile (INSERT ON CONFLICT with COALESCE for partial updates), CreateUserProfile (INSERT).

### SLICE-W07-004: UserProfile Repository (repository/postgres/user_profile_repo.go)
FindByUserID(ctx, userID), Upsert(ctx, userID, input). Uses generated sqlc code, pgtype helpers. Interface + private struct pattern.

### SLICE-W07-005: UserProfile Service (service/user_profile_service.go)
Get(ctx, userID) returns nil,nil when no profile. Update(ctx, userID, input) delegates to Upsert. Validates optional fields.

### SLICE-W07-006: UserProfile Resolver + GraphQL Schema
schema/user_profile.graphql: UserProfile type, UserProfileInput, UserProfileResult with inline error types, UserProfileErrorCode enum. resolver/user_profile.go: GetUserProfile (Query), UpdateUserProfile (Mutation). Auth guard through middleware.

### SLICE-W07-007: AiExport Migration (00092_ai_exports.sql)
CREATE TABLE ai_exports (id UUID PK, user_id UUID FK, date_range_start DATE NOT NULL, date_range_end DATE NOT NULL, include_photos BOOLEAN NOT NULL DEFAULT false, include_nutrition BOOLEAN NOT NULL DEFAULT true, include_cardio BOOLEAN NOT NULL DEFAULT true, include_measurements BOOLEAN NOT NULL DEFAULT true, user_comment TEXT, generated_prompt TEXT NOT NULL, export_file_path TEXT, created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ). Down: DROP TABLE.

### SLICE-W07-008: AiExport Model (models/ai_export.go)
AiExportRecord (DB row), AiExport (public JSON), CreateAiExportInput (mutation input). Result/error types matching WeekFlag pattern.

### SLICE-W07-009: AiExport SQLc Queries (queries/ai_exports.sql)
CreateAiExport, GetAiExportByID (user-scoped), ListAiExportsByUserID, UpdateAiExportFilePath, DeleteAiExport, ListStaleAiExports.

### SLICE-W07-010: AiExport Repository (repository/postgres/ai_export_repo.go)
Create, GetByID, ListByUserID, UpdateFilePath, Delete. Interface + private struct pattern.

### SLICE-W07-011: AiExport Service + Data Provider (service/ai_export_service.go)
AiExportService: GeneratePrompt (creates record + prompt), GenerateExport (builds ZIP), GetByID, List, Delete.
AiExportDataProvider interface: GetWorkoutSummary, GetBodyWeightEntries, GetBodyCheckIns, GetCardioEntries, GetWeekFlags, GetNutritionMacros.
Prompt building: structured plain-text with user context, period, section data, week flags, user comment.
Export flow: validate -> fetch profile -> build prompt -> create record with NULL exportFilePath -> fetch data -> build ZIP -> atomic rename -> update exportFilePath.

### SLICE-W07-012: ZIP Generation Utilities (service/export_zip.go)
ExportArchive struct: manifest, dataJSON, summaryMD, CSVs, photos. In-memory archive builder with no filesystem deps. Service handles actual file I/O. Temp-file-atomic-rename pattern for safe write.

### SLICE-W07-013: AiExport Resolver + GraphQL Schema
schema/ai_export.graphql: AiExport type, CreateAiExportInput, AiExportResult, AiExportsResult with inline errors, AiExportErrorCode enum. Mutations: createAiExportPrompt(input), generateAiExport(id), deleteAiExport(id). Queries: aiExport(id), aiExports.

### SLICE-W07-014: Download REST Handler (handler/ai_export_handler.go)
GET /api/ai-export/download?exportId= — PIN guard, resolve AiExport, verify ownership, stream ZIP via http.ServeContent. Follows progress_photo_handler.go pattern.

### SLICE-W07-015: Main Wiring (main.go, resolver.go, config.go, gqlgen)
Wire UserProfileService, AiExportService, AiExportDataProvider into resolver container. Register download route under PIN group. Add AiExportConfig to config.go. Add gqlgen bindings. Add bootstrap default UserProfile creation.

## Acceptance Criteria

### AC-W07-001 — User profile context stored and retrievable
Source: AC-084, AC-085, AC-086. GET /api/user-profile returns all persistent context fields. updateUserProfile mutation accepts partial updates with COALESCE.

### AC-W07-002 — Date range defaults to last 4 weeks
Source: RULE-021. If date range not provided, default to (now() - 28 days, now()). Uses Settings.defaultAiExportWeeks as base multiplier.

### AC-W07-003 — Custom date range accepted
Source: AC-075. Accept dateRangeStart and dateRangeEnd as optional input fields. Validate end >= start, no future dates.

### AC-W07-004 — Section toggles respected
Source: AC-076, AC-077, RULE-025. Only included sections appear in data.json, summary.md, CSV files. manifest.json.includedSections reflects actual inclusions.

### AC-W07-005 — Photos excluded by default
Source: AC-077, AC-112, RULE-025. includePhotos defaults to false. Even when true, only ProgressPhoto records in date range are included. Enforced at service layer.

### AC-W07-006 — ZIP contains valid manifest.json
Source: AC-078, AC-081. manifest.json includes schemaVersion=1, exportTimestamp, dateRange, includedSections.

### AC-W07-007 — ZIP contains valid data.json
Source: AC-078, AC-082. data.json includes all selected sections with correct data shapes.

### AC-W07-008 — ZIP contains valid summary.md
Source: AC-078, AC-083. Markdown with period, goal, workout stats, trends, nutrition summary, cardio summary, week flags.

### AC-W07-009 — ZIP contains CSV files
Source: AC-079. workouts.csv, measurements.csv, nutrition.csv, cardio.csv with correct column headers.

### AC-W07-010 — ZIP includes photos/ directory when photos included
Source: AC-080. Files in photos/{checkInId}_{angle}.{ext} when includePhotos=true.

### AC-W07-011 — One-time comment stored and included
Source: AC-087. userComment stored on AiExport, included in data.json, summary.md, manifest.json.

### AC-W07-012 — Week flags included in export
Source: AC-088. WeekFlag records in data.json.weekFlags and referenced in summary.md.

### AC-W07-013 — Generated prompt asks specific analysis questions
Source: AC-089. Prompt instructs AI to analyze progress, compare weights, evaluate volume trends, consider RPE/RIR and cardio, compare training vs body changes, give recommendations.

### AC-W07-014 — Empty date range handled gracefully
Source: EDGE-008. Export with empty arrays/empty summary. No error. Prompt still includes user context.

### AC-W07-015 — ZIP generation failure handled gracefully
Source: EDGE-024. Temp-file-atomic-rename. On failure: temp file cleaned up, no AiExport record with exportFilePath saved. Error returned.

### AC-W07-016 — Photos not included unless explicitly opted in
Source: AC-112. includePhotos must be explicitly true. Enforced at DB DDL (DEFAULT false) and service layer.

### AC-W07-017 — AI export content not logged
Source: AC-118. No prompt text, user comments, body values, photo paths in logs. Only metadata (export ID, dates, toggles) logged.

### AC-W07-018 — Generated prompt returned in generate response body
Source: PAGE-009 prompt display/copy. POST /api/ai-export/generate returns generatedPrompt in response body.

### AC-W07-019 — Date range validation rejects invalid ranges
Source: standard validation. dateRangeEnd < dateRangeStart returns error. dateRangeStart > today returns error. Range > 365 days returns error.

### AC-W07-020 — Authenticated user sees only their own data
Source: ownership principle. AiExport queries scoped by user_id from session context. Download endpoint validates ownership (returns 404 on mismatch).

### AC-W07-021 — Export file stored at correct path
Source: DDEC-W07-006. ZIP stored at {ExportBasePath}/{userId}/{exportId}.zip.

### AC-W07-022 — Export size limit enforced
Source: DDEC-W07-009. Generation rejected with clear error if estimated data size > 100MB.

### AC-W07-023 — Stale exports cleaned up
Source: DDEC-W07-007. Cleanup task deletes AiExport records and files older than 7 days. Regeneration deletes prior export.

### AC-W07-024 — UserProfile bootstrap creates default record
Source: DDEC-W07-001. Bootstrap service creates default UserProfile row (all nullable fields as NULL) when atlas_users is created.

### AC-W07-025 — Workout data returns empty arrays without WAVE-03
Source: DDEC-W07-012. data.json.workouts is [] when no WAVE-03 tables or no data. Summary.md states "No workout data recorded for this period."

## Exit Criteria

### EC-W07-001 — UserProfile migration 00091 applies cleanly
Migration creates user_profiles table with correct columns. Down works.

### EC-W07-002 — AiExport migration 00092 applies cleanly
Migration creates ai_exports table with correct columns, include_photos DEFAULT false. Down works.

### EC-W07-003 — All sqlc queries compile and regenerate without errors
bun run codegen succeeds for both user_profiles.sql and ai_exports.sql.

### EC-W07-004 — All Go code compiles
bun run build succeeds. No type errors, no missing imports.

### EC-W07-005 — GraphQL schema passes gqlgen generation
gqlgen type bindings in atlas-gqlgen.yml are complete. bun run codegen succeeds for GraphQL.

### EC-W07-006 — UserProfile service unit tests pass
TEST-W07-001 through TEST-W07-007 all pass.

### EC-W07-007 — AiExport prompt service unit tests pass
TEST-W07-008 through TEST-W07-014 all pass.

### EC-W07-008 — ZIP generation service unit tests pass
TEST-W07-015 through TEST-W07-024 all pass.

### EC-W07-009 — AiExport lifecycle tests pass
TEST-W07-025 through TEST-W07-029 all pass.

### EC-W07-010 — REST handler integration tests pass
TEST-W07-030 through TEST-W07-035 all pass.

### EC-W07-011 — Repository integration tests pass
TEST-W07-036 through TEST-W07-038 all pass (with INTEGRATION_TESTS=1).

### EC-W07-012 — Codegen drift check passes
bunx nx run api:codegen && bunx nx build api succeeds. No generated file drift.

### EC-W07-013 — All endpoints return 401 without valid PIN session
POST /api/ai-export/generate, GET /api/ai-export/download, GET /api/user-profile all protected.

### EC-W07-014 — Download endpoint validates ownership
Export owned by another user returns 404 (not 403 or 200).

### EC-W07-015 — ZIP structure validated
manifest.json, data.json, summary.md, CSVs present. manifest.json has correct schemaVersion and sections.

### EC-W07-016 — Empty date range produces valid export
No crash. Empty arrays in data.json. Summary.md states no data. Prompt still includes user context.

### EC-W07-017 — Photo opt-in enforced
Default export has no photos/ directory. With includePhotos=true, photos/ directory present with correct file names.

### EC-W07-018 — Log privacy verified
No prompt content, user comment, body values, photo paths in logs. Only metadata (IDs, dates, toggle booleans).

### EC-W07-019 — Lint passes
bun run lint succeeds for all changed packages.

### EC-W07-020 — No week flag write operations in WAVE-07
Code review confirms WAVE-07 only reads week flags via WAVE-04 service.

## Verification Obligations

### UserProfile Service Tests
| ID | Description | Command |
|---|---|---|
| TEST-W07-001 | UserProfile create succeeds with valid fields | go test -run TestUserProfileService_Create_Success |
| TEST-W07-002 | UserProfile create with empty goal | go test -run TestUserProfileService_Create_EmptyGoal |
| TEST-W07-003 | UserProfile create with invalid height | go test -run TestUserProfileService_Create_InvalidHeight |
| TEST-W07-004 | UserProfile update succeeds | go test -run TestUserProfileService_Update_Success |
| TEST-W07-005 | UserProfile get returns defaults when no profile | go test -run TestUserProfileService_Get_ReturnsDefaults |
| TEST-W07-006 | UserProfile get returns existing profile | go test -run TestUserProfileService_Get_ReturnsExisting |
| TEST-W07-007 | UserProfile log privacy — no goal in log | go test -run TestUserProfileService_Logs_NoGoalInLog |

### AiExport Prompt Service Tests
| ID | Description | Command |
|---|---|---|
| TEST-W07-008 | Prompt generation with all sections | go test -run TestAiExportPrompt_AllSections |
| TEST-W07-009 | Prompt generation with only workouts | go test -run TestAiExportPrompt_OnlyWorkouts |
| TEST-W07-010 | Prompt with persistent AI context | go test -run TestAiExportPrompt_WithPersistentContext |
| TEST-W07-011 | Prompt with one-time comment | go test -run TestAiExportPrompt_WithOneTimeComment |
| TEST-W07-012 | Prompt with week flags | go test -run TestAiExportPrompt_WithWeekFlags |
| TEST-W07-013 | Prompt with empty date range | go test -run TestAiExportPrompt_EmptyDateRange |
| TEST-W07-014 | Prompt with no data in period | go test -run TestAiExportPrompt_NoDataInPeriod |
| TEST-W07-043 | Validates date range (end before start returns error) | go test -run TestAiExportPrompt_InvalidDateRange |
| TEST-W07-044 | Validates max date range (over 365 days returns error) | go test -run TestAiExportPrompt_DateRangeOverLimit |

### ZIP Generation Tests
| ID | Description | Command |
|---|---|---|
| TEST-W07-015 | ZIP archive is valid | go test -run TestAiExportZIP_ValidArchive |
| TEST-W07-016 | manifest.json structure correct | go test -run TestAiExportZIP_ManifestStructure |
| TEST-W07-017 | data.json structure correct | go test -run TestAiExportZIP_DataJSONStructure |
| TEST-W07-018 | summary.md content present | go test -run TestAiExportZIP_SummaryMDContent |
| TEST-W07-019 | CSV files exist | go test -run TestAiExportZIP_CSVFilesExist |
| TEST-W07-020 | CSV headers and rows correct | go test -run TestAiExportZIP_CSVHeadersAndRows |
| TEST-W07-021 | Photos included when opted in | go test -run TestAiExportZIP_PhotosIncluded |
| TEST-W07-022 | Photos excluded by default | go test -run TestAiExportZIP_PhotosExcludedByDefault |
| TEST-W07-023 | No photos/ dir when opted out | go test -run TestAiExportZIP_NoPhotosDirWhenOptedOut |
| TEST-W07-024 | Workout data present in ZIP | go test -run TestAiExportZIP_WorkoutData |

### Lifecycle and Cleanup Tests
| ID | Description | Command |
|---|---|---|
| TEST-W07-025 | Temp files cleaned up on success | go test -run TestAiExportCleanup_RemovesTempFiles |
| TEST-W07-026 | Orphaned exports cleaned up | go test -run TestAiExportCleanup_OrphanedExports |
| TEST-W07-027 | Disk full returns error, no partial file | go test -run TestAiExportCleanup_DiskFull |
| TEST-W07-028 | Large date range completes | go test -run TestAiExportCleanup_LargeDateRange |
| TEST-W07-045 | Max export size rejection (over 100MB returns error) | go test -run TestAiExportCleanup_MaxSizeLimit |
| TEST-W07-029 | Log privacy — no export content in logs | go test -run TestAiExportCleanup_Logs_NoExportContent |

### REST Handler Integration Tests
| ID | Description | Command |
|---|---|---|
| TEST-W07-030 | Generate export handler | go test -run TestAiExportHandler_GenerateExport |
| TEST-W07-031 | Download export handler | go test -run TestAiExportHandler_DownloadExport |
| TEST-W07-032 | Generate export returns 401 without auth | go test -run TestAiExportHandler_GenerateExport_MissingAuth |
| TEST-W07-033 | Download returns 404 for missing export | go test -run TestAiExportHandler_DownloadExport_NotFound |
| TEST-W07-034 | Download returns 404 for wrong user's export | go test -run TestAiExportHandler_DownloadExport_OwnershipMismatch |
| TEST-W07-035 | Get user profile handler | go test -run TestUserProfileHandler_Get |
| TEST-W07-036 | User profile returns 401 without auth | go test -run TestUserProfileHandler_Get_NoSession |

### Repository Integration Tests (requires INTEGRATION_TESTS=1)
| ID | Description | Command |
|---|---|---|
| TEST-W07-037 | UserProfile repository operations | INTEGRATION_TESTS=1 go test -run TestUserProfileRepo |
| TEST-W07-038 | AiExport repository operations | INTEGRATION_TESTS=1 go test -run TestAiExportRepo |
| TEST-W07-039 | Migrations apply cleanly | INTEGRATION_TESTS=1 go test -run TestWave07Migration |

### Codegen Drift Checks
| ID | Description | Command |
|---|---|---|
| TEST-W07-040 | sqlc codegen passes | bunx nx run api:codegen && bunx nx build api |
| TEST-W07-041 | gqlgen codegen passes | bunx nx run graphql:validate && bunx nx run api:codegen |

### GraphQL Resolver Tests
| ID | Description | Command |
|---|---|---|
| TEST-W07-042 | UserProfile resolver | go test ./internal/atlas/graph/resolver/ -run TestUserProfileResolver |

## Rollout Rollback And Compatibility

### Rollout Order
1. Add AiExportConfig to appconfig/config.go and config.yml
2. Create 00091_user_profiles.sql and 00092_ai_exports.sql migrations
3. Add sqlc queries, run bun run codegen
4. Create model types for UserProfile and AiExport
5. Create repository implementations
6. Update bootstrap service to create default UserProfile
7. Create UserProfile service
8. Create AiExport service with AiExportDataProvider
9. Create ZIP generation utility with temp-file-atomic-rename
10. Create GraphQL schemas and resolvers
11. Create REST handler for download
12. Create REST handler for user-profile
13. Wire everything in main.go, resolver.go, atlas-gqlgen.yml
14. Create exports directory on startup
15. Run migrations, deploy, verify with integration tests

### Rollback
1. goose down twice (00092 then 00091)
2. Remove all new source files
3. Revert main.go wiring (services, repos, resolvers, routes)
4. Revert atlas-gqlgen.yml (remove bindings)
5. Revert schema.graphql (remove type extensions)
6. bun run codegen to purge generated code
7. Clean up exports directory: rm -rf {ExportBasePath}/
8. Remove config entries from config.yml

### Compatibility
- Additive changes only — no existing tables modified
- WAVE-03 dependency: empty arrays when not deployed (no breaking changes)
- WAVE-04 dependency: read-only week flag consumption (no breaking changes)
- WAVE-05 dependency: empty nutrition arrays when not deployed
- No existing migration sequence affected (00091 and 00092 are new)

## Handoff Packets

### Developer Handoff
Dependencies: All 6 prior waves must be deployed (WAVE-01 mandatory, others optional with empty stubs).
Implementation order: Follow SLICE dependency graph — UserProfile chain (SLICES 001-006) and AiExport chain (SLICES 007-014) are independent until SLICE-015.
Key patterns: Week_flag.go for models/repo/service/resolver patterns. Settings.go for config patterns. Progress_photo_handler.go for download handler pattern.
Generators: sqlc (user_profiles.sql, ai_exports.sql), gqlgen (atlas-gqlgen.yml bindings).
Testing: Mock repos via testify, INTEGRATION_TESTS=1 guard for repo tests, httptest for handler tests.

### Reviewer Handoff
7 reviewer verdicts in needs-revision state. All critical findings (RF-001 through RF-003) resolved in design decisions. 6 open questions (DQ-W07-001 through DQ-W07-006) documented in question-ledger.md. Ready for final fit review after revision consolidation.

## Reviewer Verdicts

| Wave | Perspective | Attempt | Verdict | Reviewer Report | Required Revisions | Notes |
|---|---|---|---|---|---|---|
| WAVE-07 | product-scope-and-ac | 1 | needs-revision | review-product-scope-and-ac-attempt-1.md | R1: Add AC for prompt in response body. R2: Resolve UserProfile/Settings conflict. | AC coverage strong. |
| WAVE-07 | architecture-codebase-fit | 1 | needs-revision | review-architecture-codebase-fit-attempt-1.md | Migration numbers, gqlgen config, display_name inconsistency. | Pattern fit approved. 3 issues. |
| WAVE-07 | data-api-integration-ops | 1 | needs-revision | review-data-api-integration-ops-attempt-1.md | F1-F8: photo default, route design, storage path, cleanup, log markers. | 8 discrepancies between planners. |
| WAVE-07 | security-privacy-compliance | 1 | needs-revision | review-security-privacy-compliance-attempt-1.md | GAP1-GAP4: temp-file, lifecycle, storage path, max size. | 10 ACs confirmed. 4 gaps. |
| WAVE-07 | testing-exit-criteria | 1 | needs-revision | review-testing-exit-criteria-attempt-1.md | Date validation tests, ownership test, sync/async dependency. | 10 items pass. 4 need revision. |
| WAVE-07 | sequencing-other-wave-fit | 1 | needs-revision | review-sequencing-other-wave-fit-attempt-1.md | R1: UserProfile duplicates Settings. R2: Week flag REST documentation. | R1 blocking. |
| WAVE-07 | traceability-consistency | 1 | needs-revision | review-traceability-consistency-attempt-1.md | F1-F10: photo default, UserProfile, AC namespace, migration numbers, URL patterns. | 10 issues, 2 critical. |

Design decisions resolve all critical and high-severity findings. Ready for final fit review.

## Open Questions

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
|---|---|---|---|---|---|---|---|---|---|---|
| DQ-W07-001 | WAVE-07 | architecture-codebase-fit | Medium | None | Schema version format for manifest.json — integer (1) or semver (1.0.0)? | Schema evolution compatibility for downstream consumers | Adopt integer schemaVersion = 1 for first version | planner-architecture-codebase-attempt-1 Q-W07-002 | resolved | MVP: integer 1 (user-approved 2026-06-21) |
| DQ-W07-002 | WAVE-07 | data-api-integration-ops | Medium | None | How to inject app version into manifest.json? | No existing -ldflags or build version injection in codebase | Add -ldflags for main.appVersion or omit from manifest for MVP | planner-architecture-codebase-attempt-1 Q-W07-006 | resolved | Omit appVersion for MVP (user-approved 2026-06-21) |
| DQ-W07-003 | WAVE-07 | security-privacy-compliance | Low | None | Max AiExport records per user — unbounded or capped? | Disk usage unbounded; cleanup TTL handles old records but count not limited | Configurable max_records_per_user (default: 50) | planner-product-ac-attempt-1 Q-W07-004 | resolved | Follow-up: add after MVP (user-approved 2026-06-21) |
| DQ-W07-004 | WAVE-07 | data-api-integration-ops | Medium | None | ZIP streaming threshold — when to switch from in-memory to streaming? | Large exports with many photos could OOM if built fully in memory | Set threshold: if estimated size > 100MB, stream to temp file | planner-data-integration-ops-attempt-1 Q-W07-DIO-03 | resolved | Threshold: 100MB (user-approved 2026-06-21) |
| DQ-W07-005 | WAVE-07 | data-api-integration-ops | Medium | None | Photo naming convention in export ZIP — UUID-based or descriptive? | Descriptive names (checkInId_angle) preserve context; UUID avoids collisions | Use {checkInId}_{angle}.{ext} | planner-product-ac-attempt-1 Q-W07-005 | resolved | Descriptive naming: {checkInId}_{angle}.{ext} (user-approved 2026-06-21) |
| DQ-W07-006 | WAVE-07 | architecture-codebase-fit | Low | None | Build WeekFlagsByDateRange query or let client call per week? | Reduces N+1 queries for the frontend; WAVE-04 only has single-week query | Build lightweight weekFlagsByDateRange query in WAVE-04 or handle in WAVE-07 | planner-sequencing-fit-attempt-1 §1.1 | resolved | Defer: client calls per week for MVP (user-approved 2026-06-21) |

## Traceability
- docs/product/prd.md Sections 17, 18
- docs/product-verified/domain-model.md#AiExport, #UserProfile
- docs/product-verified/functional-spec.md §17-18
- docs/product-verified/acceptance-criteria.md AC-074-089, AC-112, AC-118
- docs/product-verified/business-rules.md RULE-021, RULE-025, RULE-026, RULE-027, RULE-029
- docs/product-verified/edge-cases.md EDGE-008, EDGE-024
- docs/prd-waves/waves/wave-07.md (source wave)
- docs/prd-waves/frontend-pages/page-009.md
- docs/prd-wave-details/waves/wave-01.md (PIN auth, Settings)
- docs/prd-wave-details/waves/wave-04.md (week flags, body, cardio)
- docs/prd-wave-details/waves/wave-05.md (nutrition)
- docs/prd-wave-details/waves/wave-06.md (stub pattern DDEC-W06-010)
- docs/requirements.xml (after update)
- docs/development-plan.xml (M-API-USER-PROFILE, M-API-AI-EXPORT)
- docs/knowledge-graph.xml (after update)
- docs/verification-plan.xml (V-M-API-USER-PROFILE, V-M-API-AI-EXPORT)
