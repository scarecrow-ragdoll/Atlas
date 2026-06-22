# Planner Report: Sequencing & Fit (WAVE-09)

## Scope
Prior detailed wave compatibility, future wave boundaries, frontend dependency context.

## Prior Wave Dependencies
WAVE-09 depends on ALL prior waves (WAVE-01 through WAVE-08) for entity data:

| Wave | Entities Used | Status |
| --- | --- | --- |
| WAVE-01 | Foundation (atlas_users, settings) | user-approved |
| WAVE-02 | ExerciseService, ExerciseMediaService | user-approved |
| WAVE-03 | DailyLog, WorkoutExercise, WorkoutSet entities | user-approved |
| WAVE-04 | CardioEntry, BodyWeightEntry, BodyCheckIn, BodyMeasurement, ProgressPhoto | user-approved |
| WAVE-05 | NutritionProduct, NutritionTemplate, NutritionTemplateItem, DailyNutritionOverride, DailyNutritionOverrideItem | user-approved |
| WAVE-06 | No direct backup dependency | user-approved |
| WAVE-07 | UserProfile, AiExport (ListByUserID for backup) | ready-for-dev |
| WAVE-08 | AiReview (ListAllByUserID already exposed for backup) | user-approved |

## Prior Detailed Wave Compatibility

### WAVE-07 Readyness Check
- AiExportService exists with List(ctx, userID) — needs ListAllByUserID variant for backup
- ExportArchive + BuildZIP exists — backup will reuse this infrastructure
- AiExportConfig exists in appconfig — BackupConfig should follow same pattern
- Rest handler pattern (GenerateExport + DownloadExport) — backup handler follows same API design

### WAVE-08 Readyness Check
- AiReviewService already has ListAllByUserID (implemented explicitly for WAVE-09 backup consumption)
- Test mocks already include ListAllByUserID

## Future Wave Boundaries
- No future waves defined after WAVE-09 — this is the final wave
- Future: cloud backup, incremental backup — explicitly excluded from WAVE-09 scope
- Migration framework: current MVP implementation accepts same-version only. Future migration runner can be added without breaking existing import flow.

## Frontend Pages Context

### PAGE-010: Import/Export
**Backend Dependencies:**
- POST /api/backup/export — generate ZIP, return downloadId
- GET /api/backup/download — download generated ZIP
- POST /api/backup/import/validate — upload ZIP, dry-run, return summary
- POST /api/backup/import/confirm — execute restore

**Note:** The frontend page spec says "POST /api/backup/import (dry-run then confirm)" — this is a simplification. The backend will expose two separate endpoints: `/import/validate` and `/import/confirm`.

**API contract for frontend:**
1. User clicks "Export" → POST /api/backup/export {includeMedia: bool} → {downloadId, size, timestamp}
2. User clicks "Download" → GET /api/backup/download?downloadId=xxx → ZIP file download
3. User selects file → POST /api/backup/import/validate (multipart) → {validationId, summary}
4. User reviews summary → clicks "Confirm" → POST /api/backup/import/confirm {validationId} → {status, entityCounts}

## Dependency Order
```
SLICE-W09-001 (migration) → SLICE-W09-002 (data aggregation) → SLICE-W09-003 (ZIP generation) → SLICE-W09-004 (export handler) → SLICE-W09-009 (wiring)
SLICE-W09-005 (import validation) → SLICE-W09-006 (import restore) → SLICE-W09-007 (import handler) → SLICE-W09-009 (wiring)
SLICE-W09-008 (GraphQL) → SLICE-W09-009 (wiring)
```

## Scope Collision Check
- No collision with prior waves — WAVE-09 introduces new service, handler, and endpoints
- No collision with future waves — WAVE-09 is the final backend wave