# Decision Log

## Source Wave Gate
source-wave-gate: passed for WAVE-06 (2026-06-21). Source: docs/prd-waves/waves/wave-06.md.
source-wave-gate: passed for WAVE-07 (2026-06-21). Source: docs/prd-waves/waves/wave-07.md.

## User Wave Approvals
- WAVE-06 source wave: user-approved (2026-06-18) via wave-map.md
- WAVE-07 source wave: user-approved (2026-06-18) via wave-map.md
- Detailed WAVE-06: awaiting user approval
- Detailed WAVE-07: awaiting owner decisions then user approval

## Scope Decisions

### WAVE-07
- UserProfile as separate table (not extending atlas_users or Settings) — follows domain separation (DDEC-W07-001)
- AiExport as separate table with generatedPrompt, exportFilePath (DDEC-W07-002)
- include_photos DEFAULT false — enforced server-side per RULE-025 (DDEC-W07-003)
- Migration numbers: 00091 (user_profiles), 00092 (ai_exports) — additive to existing 00090 max (DDEC-W07-004)
- REST endpoints: POST /api/ai-export/generate, GET /api/ai-export/download?exportId=, GET /api/user-profile (DDEC-W07-005)
- ZIP generation: sync with temp-file-atomic-rename pattern (DDEC-W07-006)
- ZIP storage path: {ExportBasePath}/{userId}/{exportId}.zip (DDEC-W07-007)
- WAVE-03 stub pattern: empty arrays in export when no workout data available (DDEC-W07-008)
- Week flags: WAVE-07 reads only — WAVE-04 owns CRUD (DDEC-W07-009)
- Export lifecycle: 7-day TTL + delete-on-regeneration (DDEC-W07-010)
- Max export size: 100MB hard limit (DDEC-W07-011)
- Photo in ZIP: files in photos/ subfolder (DDEC-W07-012)
- ZIP format: manifest.json (schemaVersion=1), data.json, summary.md, CSVs (DDEC-W07-013)
- CAP-W07-003 removed from WAVE-07 scope — week flags owned by WAVE-04 (DDEC-W07-014)
- Prompt returned in generate response body for frontend display (DDEC-W07-015)

### WAVE-06
- Read-only wave: no mutations, no storage changes (DDEC-W06-001)
- Epley formula for e1RM (DDEC-W06-002)
- Measurement range via check_in JOIN (DDEC-W06-003)
- Nutrition weekly average via iteration per RULE-015 (DDEC-W06-004)
- Empty series for no-data periods (DDEC-W06-005)
- Default chart period: 4 weeks (DDEC-W06-006)
- Measurement overlay alphabetically ordered (DDEC-W06-007)
- Best set = highest e1RM per session (DDEC-W06-008)
- Working weight per session from WorkoutExercise.workingWeightSnapshot (DDEC-W06-009)
- Exercise chart stubs returning empty series until WAVE-03 (DDEC-W06-010)
- 52-week max date range cap (DDEC-W06-011)

## Codebase Fit Decisions
### WAVE-07
- All new files in apps/api/internal/atlas/ — consistent with existing module structure
- sqlc auto-discovery via glob — add user_profiles.sql and ai_exports.sql
- atlas-gqlgen.yml needs new model bindings for UserProfile, AiExport types
- REST download handler follows ProgressPhotoHandler.Download pattern
- AiExportDataProvider interface provides clean seam for aggregating export data
- display_name NOT duplicated in user_profiles — use atlas_users.display_name

### WAVE-06
- All new files in apps/api/internal/atlas/ — consistent with existing module structure
- sqlc auto-discovery via glob — no config changes needed
- atlas-gqlgen.yml needs new model bindings for chart types

## Deferrals
- WAVE-07 DQ-W07-003: Max AiExport records per user — follow-up after MVP
- WAVE-07 DQ-W07-006: WeekFlagsByDateRange query — client calls per week for MVP
- Exercise chart queries: stubs returning empty series until WAVE-03 deployment
- 52-week max range: enforcement via server constant

## Rejected Assumptions
### WAVE-07
- UserProfile extends atlas_users table: REJECTED — new table follows domain separation
- UserProfile extends Settings: REJECTED — Settings is app config, UserProfile is user data
- CAP-W07-003 (week flags CRUD) in WAVE-07: REJECTED — owned by WAVE-04
- include_photos DEFAULT true: REJECTED — must be false per RULE-025