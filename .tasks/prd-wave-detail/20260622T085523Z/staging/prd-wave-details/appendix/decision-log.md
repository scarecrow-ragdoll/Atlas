<!-- FILE: docs/prd-wave-details/appendix/decision-log.md -->
<!-- VERSION: 1.0.1 -->

# Decision Log

## Source Wave Gate
- Source wave: docs/prd-waves/waves/wave-09.md
- Source wave status: user-approved (2026-06-18)
- Source wave gate result: passed
- Gate check date: 2026-06-22

## User Wave Approvals
- WAVE-09 source wave: user-approved (2026-06-18)
- WAVE-09 detailed wave: questions-open (awaiting blocking-question resolution)

## Scope Decisions

| ID | Decision | Source | Rationale |
|----|---------|--------|-----------|
| DDEC-W09-001 | No new backup DB tables | WAVE-07 pattern | Backup metadata is ephemeral; replicate AI Export model with file storage only |
| DDEC-W09-002 | Reuse ExportArchive + BuildZIP | WAVE-07 export_zip.go | Generic ZIP infrastructure — add BackupManifest/BackupData types |
| DDEC-W09-003 | All-or-nothing via single PostgreSQL transaction | AC-101, AC-116 | 14+ entity INSERTs in single tx; rollback on any failure |
| DDEC-W09-004 | Same-version-only import for MVP (DQ-W09-005) | Recommended resolution | Schema version = migration count; reject version mismatch |
| DDEC-W09-005 | Exclude CSV files from backup | AC-093 | No AC requires CSV; functional-spec says "optional" |
| DDEC-W09-006 | REST endpoints for file operations | WAVE-07 AI Export pattern | GraphQL unsuitable for binary upload/download |
| DDEC-W09-007 | Two-phase import (validate → confirm) | WAVE-07 pattern, PAGE-010 spec | POST /api/backup/import/validate + POST /api/backup/import/confirm |

## Codebase Fit Decisions
- Follow WAVE-07 AI Export pattern (service + handler + config + wiring)
- Reuse export_zip.go BuildZIP infrastructure for backup ZIP
- All entity services need ListAllByUserID methods added
- No sqlc changes needed (backup reads via existing queries + service calls)
- Migration number 00094 (next after 00093_ai_reviews.sql)

## Deferrals
- Cloud backup: deferred to future scope (explicitly excluded in source wave)
- Incremental backup: deferred to future scope (explicitly excluded in source wave)
- CSV files in backup: excluded (no AC requirement)
- Migration runner for schema version upgrades: deferred to post-MVP (DQ-W09-005)
- Import with existing-data merge strategy: deferred pending product owner decision (DQ-W09-001)

## Rejected Assumptions
- No assumption that backup requires new DB tables (uses file system + existing entity services)
- No assumption about CSV file generation (excluded)
- No assumption about cloud storage integration (excluded in source wave)
- No assumption about multi-user data merging (single-user instance)