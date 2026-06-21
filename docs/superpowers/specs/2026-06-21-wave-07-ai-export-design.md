# WAVE-07: AI Export and Prompt Builder — Design Spec

**Date:** 2026-06-21
**Status:** approved-by-user
**Source:** `docs/prd-waves/waves/wave-07.md`, `docs/product-verified/`
**Detailed wave brief:** `docs/prd-wave-details/waves/wave-07.md` (619 lines)

---

## 1. Purpose

Generate AI-ready exports with structured data and prompts for ChatGPT analysis. User selects a date range (default 4 weeks), toggles data sections (workouts, cardio, body, nutrition, photos opt-in), adds optional one-time comment, and receives a ZIP with prompt text + structured data in multiple formats.

**Not in scope:** Direct ChatGPT API call, OpenAI integration, frontend rendering.

---

## 2. Architecture

Two new entity domains across the full stack, following existing Atlas patterns:

### 2.1 UserProfile (Migration 00091)
Separate from `atlas_users` and `settings`. Stores AI context: goal, height, birthDate, trainingExperience, currentTrainingSplit, preferredProgressionStyle, nutritionStrategy, persistentAiContext. Auto-created on first access via bootstrap extension.

### 2.2 AiExport (Migration 00092)
Stores export records: dateRangeStart/End, section toggle flags, generatedPrompt, exportFilePath (optional until ZIP built).

### 2.3 Data Flow
```
POST /api/ai-export/generate
  → validate auth (PIN)
  → validate input (dates, sections)
  → build prompt from context + flags
  → aggregate data via AiExportDataProvider
  → build ZIP (temp file, atomic rename)
  → save AiExport record
  → return { exportId, generatedPrompt }

GET /api/ai-export/download?exportId=
  → validate auth
  → validate ownership
  → stream ZIP file

GET /api/user-profile
  → validate auth
  → return profile with AI context fields

**Primary frontend path** for profile data: `GET /api/user-profile` REST endpoint (used by PAGE-009). GraphQL resolver also available for admin/extension use.
```

### 2.4 AiExportDataProvider Interface
Reads from existing services/repos:
- WAVE-02: ExerciseService (exercise metadata)
- WAVE-03: WorkoutExercise/WorkoutSet repos (empty arrays if not deployed)
- WAVE-04: WeekFlagService, BodyWeightEntryRepo, BodyCheckInRepo, CardioEntryRepo, ProgressPhotoRepo
- WAVE-05: NutritionProductRepo, NutritionMacroService

---

## 3. ZIP Format

```
export.zip
├── manifest.json      { schemaVersion: 1, exportTime, dateRange, sections }
├── data.json          { workouts, cardio, body, nutrition, weekFlags, userProfile }
├── summary.md         Human-readable overview
├── workouts.csv
├── measurements.csv
├── nutrition.csv
├── cardio.csv
└── photos/            {checkInId}_{angle}.{ext} files (opt-in only, max 20 per export)
```

**Config:** `ExportBasePath` (configurable), `ExportMaxRecordsPerUser` (default 50, post-MVP), `ExportMaxPhotos` (20), `ExportMaxSizeMB` (100), `ExportDefaultWeeks` (4), `ExportTTLDays` (7).

---

## 4. Design Decisions

| ID | Decision | Rationale |
|---|---|---|
| DDEC-W07-001 | Separate user_profiles table | Domain separation from atlas_users and settings |
| DDEC-W07-002 | AiExport service with DataProvider interface | Clean seam for aggregating 7+ data sources |
| DDEC-W07-003 | include_photos DEFAULT false | RULE-025 enforced server-side |
| DDEC-W07-004 | Migration 00091/00092 | Additive, no existing table changes |
| DDEC-W07-005 | REST endpoints for export, GraphQL for user-profile | Download handler follows ProgressPhoto pattern |
| DDEC-W07-006 | temp-file-atomic-rename for ZIP | Prevents partial file on error (EDGE-024) |
| DDEC-W07-007 | {base}/{userId}/{exportId}.zip storage | User-scoped, UUID filenames prevent enumeration |
| DDEC-W07-008 | Empty arrays when WAVE-03 not deployed | Same stub pattern as WAVE-06 |
| DDEC-W07-009 | WAVE-07 reads week flags via WAVE-04 service | No duplicate CRUD |
| DDEC-W07-010 | 7-day TTL + delete-on-regeneration | Cleanup policy for export files |
| DDEC-W07-011 | 100MB hard limit | Prevents OOM during ZIP generation |
| DDEC-W07-012 | Photos as files in photos/ subfolder | Standard ZIP practice, avoids base64 bloat |
| DDEC-W07-013 | manifest.json schemaVersion = 1 | Integer version, simpler than semver for MVP |
| DDEC-W07-014 | CAP-W07-003 removed (week flags) | Owned by WAVE-04 |
| DDEC-W07-015 | Prompt in generate response body | Frontend displays prompt without downloading ZIP |

---

## 5. Implementation Slices (15 total)

1. **SLICE-W07-001–006**: UserProfile chain — migration → sqlc → model → repo → service → resolver+schema
2. **SLICE-W07-007–012**: AiExport chain — migration → sqlc → model → repo → service → resolver+schema
3. **SLICE-W07-013**: ZIP generation utility
4. **SLICE-W07-014**: AiExportDataProvider interface
5. **SLICE-W07-015**: Main wiring (resolver, routes, gqlgen config)

---

## 6. Security & Privacy

- All endpoints require PIN auth (WAVE-01 middleware)
- Photos opt-in enforced server-side (RULE-025)
- ZIP stored at user-scoped path with UUID filename
- No export content in logs — only metadata (IDs, dates, toggle booleans)
- 100MB limit prevents resource exhaustion
- Ownership validation on download

---

## 7. Verification

45 verification obligations across:
- UserProfile CRUD unit tests
- AiExport prompt generation tests (all section toggle combinations, empty periods)
- ZIP structure tests (manifest, data.json, summary, CSVs, photos)
- REST handler integration tests (auth, ownership, file serving)
- Lifecycle tests (temp cleanup, disk full, large range)
- Codegen drift checks (sqlc, gqlgen)

---

## 8. Open Questions (all resolved)

All 6 questions resolved by user on 2026-06-21: integer schemaVersion, omit appVersion for MVP, 100MB streaming threshold, descriptive photo naming, defer week flags query, defer max records per user.

---

## Traceability

- `docs/prd-wave-details/waves/wave-07.md` — full detailed wave brief
- `docs/prd-waves/waves/wave-07.md` — source shallow wave
- `docs/product-verified/domain-model.md` — UserProfile, AiExport entities
- `docs/product-verified/functional-spec.md §17-18` — AI Export specification
- `docs/product-verified/business-rules.md` — RULE-021, RULE-025, RULE-026, RULE-027
- `docs/product-verified/edge-cases.md` — EDGE-008, EDGE-024
- `docs/prd-waves/frontend-pages/page-009.md` — frontend backend dependencies