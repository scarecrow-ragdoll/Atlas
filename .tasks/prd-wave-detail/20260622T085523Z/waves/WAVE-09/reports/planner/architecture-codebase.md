# Planner Report: Architecture & Codebase Fit (WAVE-09)

## Scope
Analyze WAVE-09 against existing codebase: modules, code paths, contracts, generated artifacts, likely graph deltas, implementation slices.

## Existing Modules & Patterns

### WAVE-07 AI Export Pattern (Primary Reference)
- **Service:** AiExportService in ai_export_service.go — generates ZIP with manifest.json, data.json, summary.md, CSV files, photos/
- **ZIP Infrastructure:** ExportArchive struct + BuildZIP method in export_zip.go — reusable ZIP builder
- **REST Handlers:** AiExportHandler in ai_export_handler.go — GenerateExport (POST), DownloadExport (GET), GetUserProfile (GET)  
- **Models:** AiExportRecord (DB), AiExport (public), CreateAiExportInput in models/ai_export.go
- **Config:** AiExportConfig in appconfig.go — BasePath, MaxRangeDays, MaxExportSizeBytes

### Code Paths for WAVE-09
1. **Resolver struct** (resolver.go:8) — needs BackupService field
2. **Main.go wiring** — repo → service → resolver/handler → routes
3. **ExportArchive + BuildZIP** — reusable but needs new BackupManifest / BackupData types
4. **Migration 00094** — next migration number after 00093_ai_reviews.sql

### All Entity Services for Data Aggregation
Backup must read all entities. Current state of ListAllByUserID:

| Service | ListAllByUserID Status | Action Needed |
| --- | --- | --- |
| SettingsService | Not exposed | Need to verify existing query (only FindByUserID exists) |
| UserProfileService | Not exposed | Need to verify |
| ExerciseService | Has List(ctx, userID, ...) paginated | Need all-exercises variant |
| ExerciseMediaService | Via ExerciseService | Needs all-media for all exercises |
| CardioService | Not exposed | Need ListAllByUserID |
| BodyWeightService | Not exposed | Need ListAllByUserID |
| BodyCheckInService | Not exposed | Need ListAllByUserID |
| BodyMeasurementService | Not exposed | Need ListAllByUserID |
| NutritionProductService | Not exposed | Need ListAllByUserID |
| NutritionTemplateService | Not exposed | Need ListAllByUserID |
| DailyNutritionOverrideService | Not exposed | Need ListAllByUserID |
| WeekFlagService | Not exposed | Need ListAllByUserID |
| AiExportService | Has List(ctx, userID) | Reuse |
| AiReviewService | Has ListAllByUserID | Reuse (already implemented for WAVE-09) |

## Likely Graph Deltas
1. **M-API** → add BackupService to M-API module exports
2. **V-M-API** → add verification refs for backup service unit tests, REST handler tests, repository integration tests

## Implementation Slices

| Slice ID | Description | Files Created/Modified |
| --- | --- | --- |
| SLICE-W09-001 | Migration 00094_backup_export.sql + backup service scaffolding | migrations/00094_backup_schema_version.sql, service/backup_export.go |
| SLICE-W09-002 | Backup export service — data aggregation | service/backup_export.go (buildDataJSON, collect all entities) |
| SLICE-W09-003 | Backup export service — ZIP generation (reuse ExportArchive) | service/backup_export.go (buildBackupArchive) |
| SLICE-W09-004 | Backup REST handler — export/download endpoints | handler/backup_handler.go |
| SLICE-W09-005 | Backup import service — validation (dry-run, schema version, manifest) | service/backup_import.go |
| SLICE-W09-006 | Backup import service — restore (all-or-nothing transaction) | service/backup_import.go (processImport) |
| SLICE-W09-007 | Backup REST handler — import validate/confirm endpoints | handler/backup_handler.go |
| SLICE-W09-008 | GraphQL schema + resolvers for backup metadata | graph/schema/backup.graphql, graph/resolver/backup.go |
| SLICE-W09-009 | Main wiring — resolver.go, main.go, atlas-gqlgen.yml | resolver.go, main.go, atlas-gqlgen.yml |

## Unsupported Assumptions (from source docs)
1. **No new backup_export/backup_import DB tables needed** — backup is ephemeral (generated on demand, not persisted). Import validation state lives in memory between validate and confirm steps.
2. **Schema version = migration count** — implement as migration 00094 that stamps schema_version in a settings-like record or a dedicated `atlas_schema_version` table.
3. **No CSV files in backup** — only manifest.json, data.json, media/ per AC-093.