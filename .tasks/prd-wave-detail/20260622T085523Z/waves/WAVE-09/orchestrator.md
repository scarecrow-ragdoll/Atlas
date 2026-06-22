# Orchestration State: WAVE-09 (Backup Import/Export)

## Run Metadata
- **Run ID:** 20260622T085523Z
- **Wave ID:** WAVE-09
- **Source Wave:** docs/prd-waves/waves/wave-09.md
- **Orchestrator Mode:** Combined (planner + reviewer in single session — no sub-subagents spawned)
- **Start Time:** 2026-06-22T08:55:23Z
- **Status:** reports-generated

## Contexts Loaded
### Product Sources
- docs/product-verified/features/backup-and-restore.md
- docs/product-verified/acceptance-criteria.md (AC-093 through AC-102, AC-114-116, AC-124)
- docs/product-verified/domain-model.md
- docs/product-verified/functional-spec.md
- docs/product-verified/product-brief.md (performance targets)
- docs/product-verified/business-rules.md
- docs/product-verified/user-flows.md
- docs/product-verified/edge-cases.md
- docs/product-verified/actors-and-permissions.md

### Codebase Sources
- apps/api/internal/atlas/service/export_zip.go (reusable BuildZIP + ExportArchive)
- apps/api/internal/atlas/service/ai_export_service.go (service pattern, ZIP generation)
- apps/api/internal/atlas/service/ai_review_service.go (ListAllByUserID pattern for backup)
- apps/api/internal/handler/ai_export_handler.go (REST handler pattern)
- apps/api/internal/handler/atlas_media.go (file download pattern)
- apps/api/internal/atlas/graph/resolver/resolver.go (resolver struct)
- apps/api/cmd/server/main.go (wiring pattern)
- apps/api/internal/atlas/models/ai_export.go (model pattern)
- apps/api/atlas-gqlgen.yml (gqlgen bindings)
- apps/api/internal/appconfig/config.go (config pattern)
- All entity repositories and services via grep on ListAllByUserID patterns

### Migration Context
- Latest migration: 00093_ai_reviews.sql
- Next: 00094

## Planner Reports Written
All 6 planner reports completed in reports/planner/

## Reviewer Reports Written  
All 7 reviewer reports completed in reports/reviewer/

## Decisions Made
1. **No new backup_export / backup_import tables** — backup metadata is ephemeral (generated on demand); no DB persistence needed. All import/export records are transient.
2. **New BackupConfig in appconfig** — follows AiExportConfig pattern for base_path, max_export_size_bytes, max_media_size_bytes.
3. **Reuse ExportArchive + BuildZIP** with new BackupManifest and BackupData types (separate from AiExport types).
4. **All-or-nothing via PostgreSQL transaction** — wrap all RESTORE INSERTs in a single tx with rollback on any error.
5. **Schema version = current migration count** — implement as a migration 00094 entry that stamps schema_version, resolvable at runtime.
6. **No CSV files in backup** — not required by ACs; only manifest.json, data.json, media/.

## Open Questions
See question-ledger.md (DQ-W09-001 through DQ-W09-005). Wave is **questions-open** — 2 blocking questions need product owner decisions before dev can start.