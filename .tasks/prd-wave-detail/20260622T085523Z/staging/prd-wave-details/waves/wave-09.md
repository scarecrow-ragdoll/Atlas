# Wave 09: Backup Import/Export
## Status
questions-open
## User Approval
user-approved (2026-06-18)
## Source Wave Summary
Full backup and restore capability with versioning and media support. Export generates a ZIP with manifest.json, data.json, and optional media/. Import validates manifest structure → checks schema version → dry-runs → shows summary → user confirms → all-or-nothing restore.
## Outcome After Implementation
- OUT-W09-001: Full backup ZIP export via POST /api/backup/export + GET /api/backup/download
- OUT-W09-002: Import with dry-run validation via POST /api/backup/import/validate + POST /api/backup/import/confirm
- OUT-W09-003: Schema version check during import (manifest version vs runtime version)
- OUT-W09-004: Media inclusion toggle (exercise media + progress photos optionally included in backup)
- OUT-W09-005: All-or-nothing transaction — import wraps all entity restores in single DB tx with full rollback on any failure
## Scope Included
- Backup ZIP generation with manifest.json, data.json, media/
- manifest.json with schema version, app version, export date, included sections
- data.json with ALL entities: settings, profile, exercises, workouts, cardio, body, nutrition, AI exports, AI reviews
- Media inclusion toggle
- Import dry-run validation (ZIP parse, manifest structure, schema version, entity counts)
- Import summary display (entity counts per type, media count, warnings)
- All-or-nothing restore in a single DB transaction
- Migration framework stub (schema version comparison, same-version-only for MVP)
## Scope Excluded
- Cloud backup — future scope
- Incremental backup — future scope
- CSV files in backup — not required by any AC, recommended for exclusion
- Automatic/scheduled backups — user-invoked only per RULE-028
## Dependencies And Other-Wave Fit
WAVE-09 depends on ALL prior waves (WAVE-01 through WAVE-08). All entity services are consumed during export data aggregation. Each service needs a ListAllByUserID method — AiReviewService already has it (WAVE-08), others need addition. No future backend waves defined.
## Frontend Pages Dependencies
PAGE-010 (Import/Export) requires 4 endpoints:
- POST /api/backup/export {includeMedia: bool} → {downloadId, size, timestamp}
- GET /api/backup/download?downloadId=X → ZIP stream
- POST /api/backup/import/validate (multipart) → {validationId, summary}
- POST /api/backup/import/confirm {validationId} → {status, entityCounts, mediaCount}
## Codebase Fit And Touchpoints
- WAVE-07 AI Export pattern is the primary reference: AiExportService, ExportArchive+BuildZIP, AiExportHandler, AiExportConfig
- ExportArchive+BuiZIP reusable with new BackupManifest/BackupData types
- Resolver struct (resolver.go) needs BackupService field
- Main.go wiring follows: repo → service → handler → routes
- Migration 00094 (next after 00093_ai_reviews.sql)
- 12 entity services need ListAllByUserID methods added
## Design Contracts
- **BackupExportService**: Generate(ctx, userID, includeMedia) → {downloadId string, sizeBytes int64, timestamp string, error}
- **BackupExportService**: GetDownloadPath(ctx, userID, downloadId) → {filePath string, error}
- **BackupImportService**: Validate(ctx, userID, zipData []byte) → {validationId string, summary ImportSummary, error}
- **BackupImportService**: Confirm(ctx, userID, validationId string) → {status string, entityCounts map[string]int, error}
- **ImportSummary**: {schemaVersion, appVersion, entityCounts map[string]int, mediaCount int, warnings []string}
- **BackupManifest**: {type, schemaVersion, appVersion, exportDate, sections []string}
- **BackupData**: {settings, profile, exercises, exerciseMedia, dailyLogs, workoutExercises, workoutSets, cardio, bodyWeight, bodyCheckIns, measurements, progressPhotos, nutritionProducts, nutritionTemplates, templateItems, nutritionOverrides, overrideItems, weekFlags, aiExports, aiReviews}
- **Config**: BackupConfig {BasePath, MaxExportSizeBytes, MaxImportSizeBytes}
## Data Api Integration And Operations
- Export: collect all entities → build BackupArchive → BuildZIP → temp file → atomic rename
- Import (validate): parse ZIP → validate manifest → validate schema → dry-run entity counts → return validationId
- Import (confirm): begin tx → DELETE existing (if replace) → INSERT all in order → restore media → COMMIT/ROLLBACK
- Log markers: [Backup][export][BLOCK_EXPORT_START], [Backup][export][BLOCK_EXPORT_SUCCESS], [Backup][import][BLOCK_IMPORT_COMMIT], [Backup][import][BLOCK_IMPORT_ROLLBACK]
- Privacy: log metadata only (operation, status, size), never content
- Performance: db-only export <= 15s p95, with media best-effort; db-only import <= 30s p95
## Security Privacy And Compliance
- All endpoints behind AtlasPinGuard (same as existing Atlas guarded routes)
- Single-user instance — no additional authorization
- No PIN/content/photo/comment logging (AC-117-120)
- Upload size limited via MaxBytesReader (DoS prevention)
- Validation token one-time-use with 15-minute TTL
- Backup generated on user request only (RULE-028)
## Implementation Slices
| Slice | Description | Key Files |
| --- | --- | --- |
| SLICE-W09-001 | Migration 00094 + BackupConfig | migrations/00094_backup_schema_version.sql, appconfig/config.go |
| SLICE-W09-002 | Backup export data aggregation | service/backup_export.go |
| SLICE-W09-003 | Backup export ZIP generation | service/backup_export.go (BackupArchive + BackupManifest + BackupData) |
| SLICE-W09-004 | Export REST handler | handler/backup_handler.go |
| SLICE-W09-005 | Import validation service | service/backup_import.go |
| SLICE-W09-006 | Import restore service | service/backup_import.go |
| SLICE-W09-007 | Import REST handler | handler/backup_handler.go |
| SLICE-W09-008 | GraphQL backup metadata | graph/schema/backup.graphql, graph/resolver/backup.go |
| SLICE-W09-009 | Main wiring | resolver.go, main.go, atlas-gqlgen.yml, appconfig |
## Acceptance Criteria
| AC ID | Source AC | Criterion |
| --- | --- | --- |
| AC-W09-001 | AC-093, AC-026 | Full backup ZIP contains manifest.json, data.json, media/ |
| AC-W09-002 | AC-094 | manifest.json includes type, schema version, app version, date, sections |
| AC-W09-003 | AC-095 | data.json includes all entities (settings, profile, exercises, workouts, cardio, body, nutrition, AI) |
| AC-W09-004 | AC-096 | User can include or exclude media from backup |
| AC-W09-005 | AC-097 | Import validates manifest.json structure |
| AC-W09-006 | AC-098, AC-115 | Import validates schema version |
| AC-W09-007 | AC-099 | Import runs dry-run validation before actual restore |
| AC-W09-008 | AC-100 | Import shows summary before user confirmation |
| AC-W09-009 | AC-101, AC-116 | Import restores data and media fully, or fails without partial import |
| AC-W09-010 | AC-102, AC-114 | Import displays clear error messages on validation failure |
| AC-W09-011 | AC-124 | User can create backup → reset app → import backup → verify data restored |
| AC-W09-012 | AC-026, RULE-028 | Backup export generated only on user request |
| AC-W09-013 | RULE-008 | Dry-run validation before import |
## Exit Criteria
| EC ID | Criterion | Verification |
| --- | --- | --- |
| EC-W09-001 | Migration 00094 applies/rolls back cleanly | bun run migrate |
| EC-W09-002 | Go build succeeds | bun run build |
| EC-W09-003 | Backup service unit tests pass | go test ./.../service/ -run Backup |
| EC-W09-004 | Backup handler unit tests pass | go test ./.../handler/ -run Backup |
| EC-W09-005 | Repository integration tests pass | go test ./.../postgres/ -run Backup |
| EC-W09-006 | All ACs traced in test files | grep AC-W09 test files |
| EC-W09-007 | Transaction rollback proven by test | Import failure test verifies no side effects |
| EC-W09-008 | ZIP content verified (manifest, data, media/) | BuildZIP output test |
## Verification Obligations
| TEST ID | Type | Scope |
| --- | --- | --- |
| TEST-W09-001 | Unit | BackupService.Generate — all entity services called, ZIP built |
| TEST-W09-002 | Unit | BackupService.Generate — media toggle respected |
| TEST-W09-003 | Unit | BackupService.Generate — error propagation |
| TEST-W09-004 | Unit | ImportService.Validate — manifest structure |
| TEST-W09-005 | Unit | ImportService.Validate — schema version check |
| TEST-W09-006 | Unit | ImportService.Validate — invalid ZIP rejection |
| TEST-W09-007 | Unit | ImportService.Confirm — entity insertion in correct order |
| TEST-W09-008 | Unit | ImportService.Confirm — transaction rollback |
| TEST-W09-009 | Unit | ImportService.Confirm — media file restore |
| TEST-W09-010 | Integration | Export → import round-trip, data matches |
| TEST-W09-011 | Integration | Media toggle correct |
| TEST-W09-012 | Unit | BackupArchive manifest + data format |
| TEST-W09-013 | Unit | Handler.GenerateExport — valid request response |
| TEST-W09-014 | Unit | Handler.DownloadExport — file stream |
| TEST-W09-015 | Unit | Handler.ImportValidate — summary returned |
| TEST-W09-016 | Unit | Handler.ImportConfirm — restore executed |
| TEST-W09-017 | Unit | Handler error cases |
| TEST-W09-018 | Unit | Privacy: backup logs don't contain entity content |
| TEST-W09-019 | Benchmark | db-only export < 15s p95 |
## Rollout Rollback And Compatibility
- Migration 00094 — standard apply/rollback via existing framework
- BackupConfig — new config section, hot-reloadable
- No impact on existing routes or behavior
- Import accepts same-version only (DQ-W09-005 pending)
## Handoff Packets
Wave output: `.tasks/prd-wave-detail/20260622T085523Z/waves/WAVE-09/`
## Reviewer Verdicts
| Wave | Perspective | Attempt | Verdict | Reviewer Report | Required Revisions | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| WAVE-09 | product-scope-and-ac | 1 | approved-with-questions | reports/reviewer/product-scope-and-ac.md | None | DQ-W09-001 blocking |
| WAVE-09 | architecture-codebase-fit | 1 | approved-with-revisions | reports/reviewer/architecture-codebase-fit.md | Add ListAllByUserID scope note | 12 services need new methods |
| WAVE-09 | data-api-integration-ops | 1 | approved | reports/reviewer/data-api-integration-ops.md | None | Use Redis for import state |
| WAVE-09 | security-privacy-compliance | 1 | approved-with-notes | reports/reviewer/security-privacy-compliance.md | None | Log metadata only |
| WAVE-09 | testing-exit-criteria | 1 | approved-with-notes | reports/reviewer/testing-exit-criteria.md | Add 2 tests | Add privacy + perf tests |
| WAVE-09 | sequencing-other-wave-fit | 1 | approved | reports/reviewer/sequencing-other-wave-fit.md | None | |
| WAVE-09 | traceability-consistency | 1 | approved | reports/reviewer/traceability-consistency.md | None | |
| WAVE-09 | final-wave-fit-review | 1 | pending-review | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |
## Open Questions
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W09-001 | WAVE-09 | product-ac | BLOCKING | Q-ACTOR-08, Q-AC-15 | Import behavior when data already exists (merge/replace/error) | Restore transaction design depends on this | Decision: merge, replace-silently, or reject-with-error | planner/product-ac | open | No resolution |
| DQ-W09-005 | WAVE-09 | architecture-codebase | BLOCKING | Q-EDGE-11 | Migration strategy for schema version differences | Schema comparison logic depends on this | Decision: same-version-only or implement migration runner | planner/architecture-codebase | open | Recommended: same-version-only for MVP |
| DQ-W09-002 | WAVE-09 | product-ac | WATCHLIST | Q-AC-16 | CSV files in backup — mandatory or optional? | Affects ZIP content | Confirm: exclude CSV | planner/product-ac | open | Recommended: exclude |
| DQ-W09-003 | WAVE-09 | data-integration-ops | WATCHLIST | — | Import ZIP size limit | MaxBytesReader config parameter | Size limit in MB | planner/data-integration-ops | open | Recommended: 500MB |
| DQ-W09-004 | WAVE-09 | security-privacy-compliance | WATCHLIST | — | Backup/import logging policy | Privacy compliance | Confirm log metadata only | planner/security-compliance | open | Recommended: log event metadata only |
## Traceability
- docs/product-verified/acceptance-criteria.md AC-093 through AC-102, AC-114-116, AC-124
- docs/product-verified/features/backup-and-restore.md
- docs/product-verified/functional-spec.md §20 (REQ-016)
- docs/product-verified/user-flows.md §26.11 (export), §26.12 (import)
- docs/product-verified/edge-cases.md EDGE-010, EDGE-021, EDGE-028
- docs/product-verified/business-rules.md RULE-007, RULE-008, RULE-009, RULE-028
- docs/product-verified/product-brief.md (performance targets)
- docs/prd-waves/waves/wave-09.md (source wave)
- docs/prd-waves/frontend-pages/page-010.md