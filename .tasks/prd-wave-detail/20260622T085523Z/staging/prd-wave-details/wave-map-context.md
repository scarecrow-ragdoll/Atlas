# Wave Map Context
## Selected Backend Wave Boundary
WAVE-09 is the final backend wave. It implements full backup import/export with ZIP generation, dry-run validation, all-or-nothing restore, and migration framework. It depends on ALL prior waves for data.
## Prior Backend Wave Fit
- WAVE-07: Primary pattern reference (AiExportService + ExportArchive + REST handler). AiExportService.List exists but WAVE-09 needs all-entity aggregation (not just date-range-scoped). ExportArchive.BuildZIP reusable with new types. AiExportConfig pattern in appconfig.
- WAVE-08: AiReviewService.ListAllByUserID already implemented for backup. No changes needed.
- WAVE-01 through WAVE-06: All entity services need ListAllByUserID methods added or verified.
## Future Backend Wave Fit
No future backend waves defined. WAVE-09 is the terminal wave. Excluded scope (cloud/incremental backup) is future work but not planned.
## Frontend Pages Context
PAGE-010 (Import/Export) requires 4 REST endpoints:
1. POST /api/backup/export {includeMedia} → {downloadId, size, timestamp}
2. GET /api/backup/download?downloadId=X → ZIP file
3. POST /api/backup/import/validate (multipart) → {validationId, summary}
4. POST /api/backup/import/confirm {validationId} → {status, entityCounts}

The frontend page specification says "POST /api/backup/import (dry-run then confirm)" — this is a simplification. Backend provides two separate endpoints. The frontend implementation must chain them.
## Dependency Order
```
SLICE-W09-001 (migration + config) → SLICE-W09-002 (data aggregation) → SLICE-W09-003 (ZIP gen) → SLICE-W09-004 (export handler)
SLICE-W09-005 (import validation) → SLICE-W09-006 (import restore) → SLICE-W09-007 (import handler)
SLICE-W09-008 (GraphQL) → SLICE-W09-009 (main wiring)
```
## Scope Collision Check
No scope collision with any prior or future wave. WAVE-09 introduces entirely new service, handler, and endpoints. The only shared infrastructure is ExportArchive (which is generic and designed for reuse).