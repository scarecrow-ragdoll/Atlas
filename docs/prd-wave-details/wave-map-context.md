# Wave Map Context

## Selected Backend Wave Boundary
WAVE-07: AI Export and Prompt Builder. Generate AI-ready exports with structured data and prompts for ChatGPT analysis. Includes persistent AI context (UserProfile), prompt builder with period selection, ZIP export with manifest.json/data.json/summary.md/CSVs, week flags support, one-time comment support, and section toggles (photos optional). No direct ChatGPT API call.

## Prior Backend Wave Fit
- WAVE-01 (Foundation): prerequisite — provides PIN auth middleware, Atlas GraphQL endpoint, Settings with defaultAiExportWeeks
- WAVE-02 (Exercise Library): provides exercises table — WAVE-07 reads exercise metadata for export data.json
- WAVE-03 (Workout Diary): provides workout_sets and daily_log_exercises tables. WAVE-03 NOT implemented yet — workout export data returns empty arrays
- WAVE-04 (Cardio and Body Tracking): provides body_weight_entries, body_check_ins, body_measurements, progress_photos, cardio_entries tables. WAVE-07 reads body/cardio/photo data for export. WAVE-04 also owns WeekFlag CRUD (CAP-W04-006) — WAVE-07 reads week flags via WAVE-04 service
- WAVE-05 (Nutrition): provides nutrition tables and NutritionMacroService. WAVE-07 reads nutrition data for export
- WAVE-06 (Charts): read-only chart queries. WAVE-07 shares same underlying data — no direct dependency

No pattern or contract conflicts with prior detailed waves after removing CAP-W07-003 (week flags CRUD owned by WAVE-04).

## Future Backend Wave Fit
- WAVE-08 (AI Review History): depends on WAVE-07 — AiExport table provides export record and generatedPrompt for review context
- WAVE-09 (Backup Import/Export): shares ZIP serialization patterns but different purpose (full backup vs per-period AI export). WAVE-07 ZIP uses ai-exports/<uuid>/ subfolder pattern to avoid collision

No scope collision. WAVE-07 creates foundations for WAVE-08 (AiExport table) and shares ZIP patterns with WAVE-09.

## Frontend Pages Context
- PAGE-009 (AI Export): primary consumer. Backend provides:
  - POST /api/ai-export/generate — generate prompt and ZIP
  - GET /api/ai-export/download?exportId= — download generated ZIP
  - GET /api/user-profile — read user goal and AI context
  - GET /api/week-flags — via WAVE-04 GraphQL (no REST endpoint needed)
- PAGE-001 (Dashboard): quick action button to generate AI report
Dependency context only; no frontend pages, UI, or UX work in this wave.

## Dependency Order
WAVE-01 → WAVE-02 → WAVE-03 (partial) → WAVE-04 → WAVE-05 → WAVE-06 → WAVE-07 → WAVE-08 → WAVE-09

## Scope Collision Check
- Week flags CRUD: REMOVED from WAVE-07 (CAP-W07-003). WAVE-04 owns WeekFlag CRUD. WAVE-07 reads via WAVE-04 service
- UserProfile: new table (00091) — no collision with any prior wave (atlas_users only has display_name)
- AiExport: new table (00092) — no collision. WAVE-08 builds on this but does not modify
- ZIP generation: no collision with WAVE-09 — different purpose, different path prefix
- Migration numbers: 00091 and 00092 are additive — no migration number collision
- All new operations are additive — no existing tables modified