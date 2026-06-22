# Codebase Fit
## Relevant Modules
- **M-API** (GoSharedHTTPAPI, layer 1) — all WAVE-09 code lives here
- Apps that interact: M-WEB-ADMIN (PAGE-010 frontend consumes backup endpoints)
- No Go module changes outside M-API

## Relevant Files Read
| File | Purpose | Result |
| --- | --- | --- |
| ai_export_service.go | Service pattern for ZIP generation | Reference pattern for BackupService |
| export_zip.go | ExportArchive + BuildZIP | Reusable, needs new BackupManifest/BackupData |
| ai_export_handler.go | REST handler pattern | Reference for BackupHandler |
| atlas_media.go | File download pattern | Reference for BackupHandler.Download |
| resolver.go | Resolver struct | Needs BackupService field |
| main.go | Wiring pattern | Follows repo → service → handler/resolver → routes |
| ai_review_service.go | ListAllByUserID implementation | Already done for backup |
| atlas-gqlgen.yml | gqlgen bindings | Needs Backup types |
| appconfig.go | Config | Needs BackupConfig section |

## Public Contracts
- BackupExportService interface (Generate, GetDownloadPath)
- BackupImportService interface (Validate, Confirm)
- BackupManifest, BackupData structs (new types in models/ or service/)
- BackupConfig in appconfig

## Generated Artifact Impact
- **atlas-gqlgen.yml** — needs BackupResult, BackupSummary, BackupErrorCode types
- **GraphQL resolvers** — new backup.graphql schema file
- **No sqlc changes** — backup uses existing entity queries, no new backup tables

## Integration Points
- All 14+ entity services consumed during export (each needs ListAllByUserID)
- All 14+ entity repositories consumed during import INSERT
- Media storage directory read during export and write during import
- File system for temp ZIP files during export generation

## Likely Graph Deltas
- M-API: add BackupService, BackupImportService to module exports
- V-M-API: add verification refs for backup service/handler/integration tests
- Add M-BACKUP sub-module reference in knowledge-graph.xml

## Unsupported Assumptions
- No new backup_export/backup_import DB tables needed (backup is ephemeral)
- Schema version = migration count (implemented via migration 00094)
- No CSV files in backup (not required by ACs)