# Planner Report: Product-AC — WAVE-07 (AI Export and Prompt Builder)

**Run ID:** 20260621T170113Z  
**Planner role:** product-ac  
**Wave:** WAVE-07  
**Source wave file:** `docs/prd-waves/waves/wave-07.md` (user-approved 2026-06-18)  
**Attempt:** 1

---

## 1. Outcome Definition — What Backend Delivers After WAVE-07

After WAVE-07, the backend delivers a complete AI export and prompt builder system:

| Outcome ID | Summary | Backend Deliverable |
|---|---|---|
| OUT-W07-001 | Prompt builder with period selection | `POST /api/ai-export/generate` (or GraphQL mutation `generateAiExport`) accepts `dateRangeStart`, `dateRangeEnd`, section toggles, optional `userComment`, generates a formatted AI prompt text stored on the `AiExport` record |
| OUT-W07-002 | ZIP export with manifest.json, data.json, summary.md, CSV | Generation service writes a ZIP archive to the configured export directory; returns a download token or file path |
| OUT-W07-003 | Week flags support | API endpoint to list week flags for a date range (summarizes flags across multiple weeks); already exists for single-week via `weekFlags(weekStartDate)` — add range-based list or use existing query per week |
| OUT-W07-004 | One-time comment support | `AiExport.userComment` field stored on the export record, included in `data.json` and `summary.md` |
| OUT-W07-005 | Section toggles (photos optional) | Input accepts per-section boolean flags (`includePhotos`, `includeNutrition`, `includeCardio`, `includeMeasurements`, `includeBodyWeight`, `includeWorkouts`); photos default false per RULE-025 |

---

## 2. Included Backend Scope

| Capability | Source Document | What to Build |
|---|---|---|
| CAP-W07-001 Persistent AI context | domain-model.md §UserProfile | `UserProfile` model/table: `goal`, `height`, `birthDate`, `trainingExperience`, `currentTrainingSplit`, `preferredProgressionStyle`, `nutritionStrategy`, `persistentAiContext`. Create migration, models, repository, service, GraphQL schema (`UserProfile` type + `getUserProfile`/`updateUserProfile`). |
| CAP-W07-002 User goal storage | domain-model.md §UserProfile | Part of CAP-W07-001; `goal` and `persistentAiContext` are fields on the same `UserProfile` entity |
| CAP-W07-003 Week flags CRUD | domain-model.md §WeekFlag | **ALREADY EXISTS** in codebase (WAVE-04). No new backend work needed for basic CRUD. May need new query: `weekFlagsByDateRange(from, to)` or consolidate client-side. |
| CAP-W07-004 Prompt generation | functional-spec.md §17-18 | Prompt builder service — reads `UserProfile` context + selected week flags + one-time comment + period data, generates plain-text prompt engineered for manual copy-paste (RULE-027, AC-089) |
| CAP-W07-005 AI export ZIP creation | functional-spec.md §17.4 | ZIP assembly service — creates temp directory, writes `manifest.json`, `data.json`, `summary.md`, `workouts.csv`, `measurements.csv`, `nutrition.csv`, `cardio.csv`, optionally `photos/` directory, returns ZIP file |
| CAP-W07-006 manifest.json | functional-spec.md §17.5 | JSON with: `exportType`, `schemaVersion`, `appVersion`, `date`, `period` (start/end), `includedSections` array, `userComment` (if provided) |
| CAP-W07-007 data.json | functional-spec.md §17.6 | JSON with all selected entities for the period: workouts (with exercises, sets, comments, RPE/RIR), body weight entries, measurements, cardio entries, nutrition data, goal, additional context, week flags summary, one-time comment |
| CAP-W07-008 summary.md | functional-spec.md §17.7 | Markdown with: period description, user goal, workout stats (total sessions, volume, best sets), exercise trends, weight/measurement changes, nutrition summary (avg KJBJU), cardio summary, notable comments |
| CAP-W07-009 CSV files | functional-spec.md §17.8 | Four CSV files: `workouts.csv`, `measurements.csv`, `nutrition.csv`, `cardio.csv` |

### Existing code already providing WAVE-07 scope

| Component | Location | Status |
|---|---|---|
| WeekFlag model + enum + validation | `apps/api/internal/atlas/models/week_flag.go` | ✅ Complete |
| WeekFlag repository (CRUD) | `apps/api/internal/atlas/repository/postgres/week_flag_repo.go` | ✅ Complete |
| WeekFlag service (create/list/delete) | `apps/api/internal/atlas/service/week_flag.go` | ✅ Complete |
| WeekFlag GraphQL resolvers | `apps/api/internal/atlas/graph/resolver/week_flag.go` + `.resolvers.go` | ✅ Complete |
| WeekFlag GraphQL schema | `apps/api/internal/atlas/graph/schema/week_flag.graphql` | ✅ Complete |
| WeekFlag DB migration | `apps/api/internal/repository/postgres/migrations/00089_week_flags.sql` | ✅ Complete |
| Settings with `defaultAiExportWeeks` | `apps/api/internal/atlas/models/settings.go`, repo, service, resolver, schema | ✅ Complete |
| `default_ai_export_weeks` in DB | `apps/api/internal/repository/postgres/migrations/00080_atlas_foundation.sql` | ✅ Complete |

### New backend scope needed

| Component | Priority | Notes |
|---|---|---|
| `user_profile` DB table + migration | Required | New table: `id UUID PK`, `user_id UUID FK`, `goal TEXT`, `height NUMERIC`, `birth_date DATE`, `training_experience TEXT`, `current_training_split TEXT`, `preferred_progression_style TEXT`, `nutrition_strategy TEXT`, `persistent_ai_context TEXT`, `created_at`, `updated_at`. One row per user, enforced by UNIQUE(user_id) or single-row insert at bootstrap. |
| `UserProfile` model types | Required | `UserProfileRecord` (DB), `UserProfile` (public), `UpdateUserProfileInput`, `UserProfileResult` with error union. Pattern: same as `settings.go`. |
| `UserProfile` repository | Required | `Get(ctx, userID)`, `Upsert(ctx, userID, input)`. Returns `UserProfileRecord`. |
| `UserProfile` service | Required | `Get(ctx, userID)`, `Update(ctx, userID, input)`. Pattern: same as `settings_service.go`. |
| `ai_exports` DB table + migration | Required | New table: `id UUID PK`, `user_id UUID FK`, `date_range_start DATE NOT NULL`, `date_range_end DATE NOT NULL`, `include_photos BOOLEAN NOT NULL DEFAULT FALSE`, `include_nutrition BOOLEAN NOT NULL DEFAULT TRUE`, `include_cardio BOOLEAN NOT NULL DEFAULT TRUE`, `include_measurements BOOLEAN NOT NULL DEFAULT TRUE`, `user_comment TEXT`, `generated_prompt TEXT NOT NULL`, `export_file_path TEXT`, `created_at TIMESTAMPTZ NOT NULL DEFAULT now()`. |
| `AiExport` model types | Required | `AiExportRecord` (DB), `AiExport` (public), `GenerateAiExportInput`, `AiExportResult` with error union. |
| `AiExport` repository | Required | `Create(ctx, userID, params)`, `GetByID(ctx, userID, id)`, `List(ctx, userID, limit, offset)`. |
| Prompt generator service | Required | Takes `UserProfile` + date range + section toggles + `userComment` + aggregated data; produces formatted prompt text per AC-089. |
| ZIP assembly service | Required | Takes `AiExport` record + aggregated data + photos (if opted-in); writes files to temp dir, creates ZIP. |
| Data aggregation queries | Required | Fetch workouts + sets + exercises + body weight + measurements + cardio + nutrition for date range. These queries span multiple existing tables. |
| `summary.md` generator | Required | Takes aggregated data + prompt context; produces human-readable Markdown. |
| CSV generators | Required | Four CSV files from aggregated data. |
| GraphQL schema additions | Required | Types: `UserProfile`, `UserProfileResult`, `UpdateUserProfileInput`, `UserProfileError`, `AiExport`, `AiExportResult`, `GenerateAiExportInput`, `AiExportError`. Queries: `userProfile`, `aiExport(id)`, `aiExports`. Mutations: `updateUserProfile`, `generateAiExport`, `downloadAiExport(id)`. |
| GraphQL resolvers | Required | Resolvers for all new schema types and mutations. |
| `weekFlagsByDateRange` query (optional) | Low | For frontend to fetch flags across the export period in one call. Currently only `weekFlags(weekStartDate)` exists — client can call per week. |

---

## 3. Excluded Backend Scope

| Feature | Reason |
|---|---|
| Direct ChatGPT/OpenAI API call | Explicitly excluded per RULE-029, WAVE-07 wave doc §Excluded Scope |
| AI review history (AiReview) | Belongs to WAVE-08 |
| Backup/restore ZIP import | Belongs to WAVE-09 |
| Frontend pages, routes, screens, UX states, components, frontend tests | Out of scope per planner instructions |
| Frontend ZIP download handler | Frontend owns the download UX; backend provides file at known path |
| Photo resizing/compression for export | Not specified in source docs |
| Watermarking or metadata stripping on exported photos | Not specified in source docs |
| CSV handling for missing optional data (RPE, RIR, notes) | Include empty columns — non-blocking implementation detail |

---

## 4. Acceptance Criteria Contributions

### AC-W07-001 — User profile context stored and retrievable
Source: AC-084, AC-085, AC-086  
Backend: `GET userProfile` returns all persistent context fields. `updateUserProfile` accepts partial updates.

### AC-W07-002 — Date range defaults to last 4 weeks
Source: AC-074, RULE-021  
Backend: If `dateRangeStart` and `dateRangeEnd` are not provided, default to `(now() - 28 days, now())`. Use `Settings.defaultAiExportWeeks` as base multiplier (4 weeks default).

### AC-W07-003 — Custom date range accepted
Source: AC-075  
Backend: Accept `dateRangeStart` and `dateRangeEnd` as optional input fields; validate end >= start, no future dates (or cap to today).

### AC-W07-004 — Section toggles respected
Source: AC-076, AC-077, RULE-025  
Backend: `includePhotos` defaults to `false`. Only included sections appear in `data.json`, `summary.md`, and CSV files. `manifest.json.includedSections` reflects actual inclusions.

### AC-W07-005 — Photos excluded by default
Source: AC-077, AC-112, RULE-025  
Backend: Input field `includePhotos` has default value `false`. Even when `true`, only photos from `ProgressPhoto` records in the date range are included.

### AC-W07-006 — ZIP contains manifest.json
Source: AC-078, AC-081  
Backend: `manifest.json` includes `exportType: "ai-export"`, `schemaVersion: "1.0"`, `appVersion` (from build config), `createdAt`, `period: {start, end}`, `includedSections: [...]`.

### AC-W07-007 — ZIP contains data.json
Source: AC-078, AC-082  
Backend: `data.json` includes all selected sections: `workouts` (array with nested exercises, sets, comments, RPE/RIR), `bodyWeight`, `measurements`, `cardio`, `nutrition`, `goal`, `additionalContext` (profile fields), `weekFlags` (array), `userComment`.

### AC-W07-008 — ZIP contains summary.md
Source: AC-078, AC-083  
Backend: Markdown file with period, user goal, workout stats (total sessions, volume lifted, PRs), exercise trends, weight/measurement deltas, nutrition averages, cardio summary, notable week flags.

### AC-W07-009 — ZIP contains CSV files
Source: AC-079  
Backend: `workouts.csv`, `measurements.csv`, `nutrition.csv`, `cardio.csv`. Each with headers appropriate to the data shape.

### AC-W07-010 — ZIP includes photos/ directory when photos included
Source: AC-080  
Backend: When `includePhotos === true`, copies referenced photo files into `photos/{uuid}.{ext}` within the ZIP. Files are read from the configured media storage directory.

### AC-W07-011 — One-time comment stored and included
Source: AC-087  
Backend: `GenerateAiExportInput.userComment` (optional string). Stored on `AiExportRecord.userComment`. Included in `data.json`, `summary.md`, and `manifest.json`.

### AC-W07-012 — Week flags included in export
Source: AC-088  
Backend: For the export date range, aggregate all `WeekFlag` records whose `weekStartDate` falls within the range. Include in `data.json.weekFlags` and reference in `summary.md`.

### AC-W07-013 — Generated prompt asks specific analysis questions
Source: AC-089  
Backend: Generated prompt text must instruct the AI to: analyze progress, compare actual vs working weights, evaluate volume trends, consider RPE/RIR and cardio data, compare training vs body changes, and give next-week recommendations.

### AC-W07-014 — Empty date range handled gracefully
Source: EDGE-008  
Backend: When no data exists in the selected period, generate export with empty arrays/empty summary. The prompt still includes user context. The response must not error — an export with empty data is valid.

### AC-W07-015 — ZIP generation failure (disk full) handled gracefully
Source: EDGE-024  
Backend: If ZIP generation fails (I/O error, disk full), the `generateAiExport` mutation returns an error result; no `AiExport` record is created with `exportFilePath`. The temp directory is cleaned up on failure.

### AC-W07-016 — Photos not included in export unless explicitly opted in
Source: AC-112  
Backend: Enforced at input validation level. `includePhotos` must be explicitly `true` to include photos.

### AC-W07-017 — AI export content not logged
Source: AC-118  
Backend: The generated prompt and export data must not appear in application logs. Only non-content metadata (export ID, date range, file size) may be logged.

---

## 5. Product Edge Cases — AI Export Perspective

| Edge Case ID | Description | Backend Handling |
|---|---|---|
| EDGE-008 | AI export date range with no data in period | Export is still generated with empty arrays in `data.json`, a minimal `summary.md` stating "No data in selected period", and a valid `manifest.json`. Prompt still includes user context. |
| EDGE-024 | Disk full during export ZIP generation | Catch I/O errors during file writes. Clean up temp directory. Return error result. No partial export record saved. |
| EDGE-025 | Docker volume full — photo copy fails | Same as EDGE-024; if `includePhotos` is true and photo copy fails, the entire export fails. Do not generate ZIP with missing photos. |
| EDGE-006 | Check-in with 0-1 photos when photos included | Include whatever photos exist (0, 1, 2-4). No validation on photo count during export. |
| EDGE-018 | Exercise deleted but referenced in historical data | Historical data (workout sets) reference exercises by ID. If exercise is deleted, include `{id, name: "Deleted Exercise"}` fallback. The data query should left-join or include exercise data at time of set recording. |
| EDGE-031 | Timezone handling | All dates are stored as dates (no timezone). Export range is date-based, not time-based. Ensure date comparisons are date-only. |

---

## 6. Questions Raised for Product Scope

| ID | Question | Context | Suggested Resolution |
|---|---|---|---|
| Q-W07-001 | Should `UserProfile` auto-create when the app boots (like `DefaultUser`), or require the user to explicitly fill in profile? | The domain model shows UserProfile as a separate entity from `atlas_users`, but no migration exists. The profile may be empty at first use. | Auto-create a default `UserProfile` row during bootstrap (like `atlas_users`), with null/default fields. The prompt builder can still work with missing fields by omitting them from the prompt. |
| Q-W07-002 | What is the `schemaVersion` for manifest.json? Should it match the app version? | `manifest.json` needs a `schemaVersion` field. No schema version convention exists yet. | Use `"1.0"` for MVP. Document that schemaVersion is independent of appVersion; increment when export JSON structure changes. |
| Q-W07-003 | Should `AiExport` records be retained after download? Can user re-download? | `AiExport.exportFilePath` is optional — once downloaded, should the ZIP be kept on server? | Keep the ZIP file on disk after generation. Allow `downloadAiExport(id)` to stream it again. Clean up old exports via a configurable retention policy (e.g., 30 days, or never). |
| Q-W07-004 | Is there a limit on how many `AiExport` records a user can create? | No source document mentions export limits. Unbounded exports could fill disk. | Add a configurable max exports limit (default: unlimited for MVP). Document as follow-up risk. |
| Q-W07-005 | Should photos in ZIP be renamed/sequenced or keep original filenames? | Photos are stored on filesystem with UUID-based paths. Direct copy would lose context. | Use `photos/{checkin_date}_{angle}_{original_filename}` or `photos/{uuid}.{ext}`. Recommend recording the mapping in `data.json.photos`. |
| Q-W07-006 | Is `appVersion` available from the Go build? | `manifest.json` needs an app version string. No existing version injection in the codebase. | Add `-ldflags` to inject `main.appVersion` at build time and read it from the export service. |

---

## 7. Verification Obligations

| Test Scope | What to Verify | Level |
|---|---|---|
| UserProfile repository | Upsert creates/updates record; `Get` returns profile; empty state | Unit |
| UserProfile service | Update validates optional fields; returns profile after update | Unit |
| UserProfile GraphQL resolver | Auth check, error mapping for not-found, internal error | Unit |
| AiExport model | Input validation (date range end >= start, no future dates) | Unit |
| AiExport service — generate | Creates DB record with correct fields; generates prompt string; rejects invalid inputs | Unit |
| AiExport service — download | Returns file path for existing export; returns error for missing/non-generated export | Unit |
| Prompt generator | Output contains user context fields; contains section-specific analysis instructions per AC-089; handles missing/null fields gracefully | Unit |
| ZIP assembly — structure | ZIP contains correct files based on section toggles; photos only when opted in | Unit |
| ZIP assembly — manifest.json | Correct schema, sections, metadata | Unit |
| ZIP assembly — data.json | All selected sections present; empty arrays for empty periods | Unit |
| ZIP assembly — summary.md | Non-empty content; references all sections | Unit |
| ZIP assembly — CSVs | Valid CSV with headers; empty files for empty data | Unit |
| EDGE-008 | Export with empty date range still succeeds with valid structure | Unit |
| EDGE-024 | Disk full during ZIP returns error, no partial record | Unit (mock filesystem) |
| AC-118 | No prompt/export content in log output | Integration (log capture) |

---

## 8. Traceability

| Source Document | Reference | Mapped To |
|---|---|---|
| `docs/product-verified/functional-spec.md` §17-18 | AI Export | WAVE-07 all ACs |
| `docs/product-verified/domain-model.md` §UserProfile | UserProfile entity | CAP-W07-001, CAP-W07-002 |
| `docs/product-verified/domain-model.md` §AiExport | AiExport entity | CAP-W07-005 through CAP-W07-009 |
| `docs/product-verified/domain-model.md` §WeekFlag | WeekFlag entity | CAP-W07-003 (already exists) |
| `docs/product-verified/business-rules.md` RULE-021 | Default period 4 weeks | AC-W07-002 |
| `docs/product-verified/business-rules.md` RULE-025 | Photos opt-in | AC-W07-005, AC-W07-016 |
| `docs/product-verified/business-rules.md` RULE-026 | On-demand generation | AC-W07-003 |
| `docs/product-verified/business-rules.md` RULE-027 | Manual copy-paste | AC-W07-013 |
| `docs/product-verified/business-rules.md` RULE-029 | No external API calls | Excluded scope |
| `docs/product-verified/edge-cases.md` EDGE-008 | Empty data period | AC-W07-014 |
| `docs/product-verified/edge-cases.md` EDGE-024 | Disk full during export | AC-W07-015 |
| `docs/product-verified/acceptance-criteria.md` AC-023 | Generate AI prompt | AC-W07-013 |
| `docs/product-verified/acceptance-criteria.md` AC-024 | Download export ZIP | AC-W07-006 through AC-W07-010 |
| `docs/product-verified/acceptance-criteria.md` AC-074-083 | Export criteria | AC-W07-001 through AC-W07-013 |
| `docs/product-verified/acceptance-criteria.md` AC-084-089 | Prompt builder criteria | AC-W07-001, AC-W07-011, AC-W07-012, AC-W07-013 |
| `docs/product-verified/acceptance-criteria.md` AC-112 | Photos excluded by default | AC-W07-016 |
| `docs/product-verified/acceptance-criteria.md` AC-118 | Export content not logged | AC-W07-017 |
| `docs/prd-waves/frontend-pages/page-009.md` | PAGE-009 backend deps | POST generate, GET download, GET userProfile, GET weekFlags |