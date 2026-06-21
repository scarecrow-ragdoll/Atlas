# Wave Map Context

## Selected Backend Wave Boundary
WAVE-07 (AI Export and Prompt Builder): 
- UserProfile CRUD (persistent AI context, goal storage)
- AiExport generate + download (prompt, ZIP, CSVs, photos)
- Prompt builder with section toggles and week flag integration
- 15 implementation slices across 2 entity chains + ZIP generation + wiring
- All endpoints PIN-guarded

## Prior Backend Wave Fit

### WAVE-01 (Foundation) — HARD DEPENDENCY
- PIN auth middleware required for all 3 endpoints
- atlas_users table for user identity
- Settings provides defaultAiExportWeeks config
- Bootstrap service extended to create default UserProfile
- **Cannot deliver WAVE-07 without WAVE-01**

### WAVE-02 (Exercise Library) — READY, READ-ONLY
- Exercise metadata for name, muscleGroups in export data.json
- No scope collision. WAVE-07 consumes via service layer.

### WAVE-03 (Workout Diary) — STUB PATTERN
- Workout data (daily_log, workout_exercise, workout_set) consumed if available
- Empty arrays when not deployed (matching WAVE-06 DDEC-W06-010 stub pattern)
- No WAVE-03 table creation in WAVE-07. EC-W07-020 enforces this.

### WAVE-04 (Cardio and Body Tracking) — READY, READ-ONLY
- CardioEntry, BodyWeightEntry, BodyCheckIn, BodyMeasurement, ProgressPhoto consumed via service layer
- WeekFlag CRUD is WAVE-04 scope — WAVE-07 reads via WeekFlagService
- PAGE-009 week flag browsing uses WAVE-04 GraphQL query (no REST proxy)
- Migration 00089 already exists for week_flags

### WAVE-05 (Nutrition) — READY, READ-ONLY
- NutritionProduct, NutritionTemplate, DailyNutritionOverride consumed via service layer
- NutritionMacroService for KJBJU averages
- Empty arrays when not deployed

### WAVE-06 (Charts) — READY, NO DIRECT DEPENDENCY
- Shares same underlying data. No dependency on WAVE-06 query services.
- Stub pattern precedent (DDEC-W06-010) applied to WAVE-07 workout data.

## Future Backend Wave Fit

### WAVE-08 (AI Review History) — CLEAN BOUNDARY
- WAVE-07 creates AiExport record with generatedPrompt
- WAVE-08 creates independent AiReview record (date range + AI response text)
- WAVE-08 does not modify AiExport. No scope collision.

### WAVE-09 (Backup Import/Export) — SIMILAR PATTERN, DIFFERENT PURPOSE
- WAVE-07: per-period AI export ZIP at {ExportBasePath}/{userId}/{exportId}.zip
- WAVE-09: full data backup with version manifest
- No shared code extraction needed. Pattern similarity noted.

## Frontend Pages Context
- **PAGE-009 (AI Export)**: Three backend endpoints — POST /api/ai-export/generate, GET /api/ai-export/download?exportId=, GET /api/user-profile
- Week flag browsing uses WAVE-04 GraphQL weekFlags query (not a WAVE-07 endpoint)
- Prompt display/copy enabled by returning generatedPrompt in generate response body

## Dependency Order
- SLICE-W07-001 (UserProfile migration) → SLICE-W07-002 (model) → SLICE-W07-003 (sqlc) → SLICE-W07-004 (repo) → SLICE-W07-005 (service) → SLICE-W07-006 (resolver+schema)
- SLICE-W07-007 (AiExport migration) → SLICE-W07-008 (model) → SLICE-W07-009 (sqlc) → SLICE-W07-010 (repo) → SLICE-W07-011 (service) + SLICE-W07-012 (ZIP utils) → SLICE-W07-013 (resolver+schema) → SLICE-W07-014 (download handler)
- SLICE-W07-015 (main wiring) depends on all above
- UserProfile chain and AiExport chain are independent until SLICE-015

## Scope Collision Check
- CAP-W07-003 (week flags CRUD): REMOVED — WAVE-04 owns this
- CAP-W07-001/002 (UserProfile): New entity. No overlap with Settings (separate concerns per domain-model.md)
- AiExport vs AiReview: Clean boundary with WAVE-08
- Export ZIP vs Backup ZIP: Different purpose, different structure