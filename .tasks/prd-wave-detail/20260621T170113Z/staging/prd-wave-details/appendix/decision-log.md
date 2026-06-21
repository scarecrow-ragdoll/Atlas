# Decision Log

## Source Wave Gate
Source wave: docs/prd-waves/waves/wave-07.md — user-approved (2026-06-18)
Source wave gate status: passed

## User Wave Approvals
- 2026-06-18: Full 9-wave backend map and 11 frontend pages approved by user
- WAVE-07 specifically: user-approved on 2026-06-18

## Scope Decisions
| Decision | ID | Value | Source |
|---|---|---|---|
| UserProfile as separate entity | DDEC-W07-001 | Create separate user_profiles table (NOT extending atlas_users or settings) | domain-model.md, reviewer R2 resolution |
| Week flags CRUD removed | DDEC-W07-002 | CAP-W07-003 removed from WAVE-07 scope. WAVE-04 owns week flags. WAVE-07 reads via WAVE-04 service. | sequencing-fit review, planner agreement |
| REST endpoint design | DDEC-W07-003 | POST /api/ai-export/generate, GET /api/ai-export/download?exportId=, GET /api/user-profile | User design decision |
| include_photos default false | DDEC-W07-004 | DEFAULT false enforced at DB DDL + service layer | RULE-025, AC-077, AC-112 |
| Migration numbers | DDEC-W07-005 | 00091_user_profiles.sql, 00092_ai_exports.sql | Latest is 00090 |
| Storage path | DDEC-W07-006 | {ExportBasePath}/{userId}/{exportId}.zip | Security recommendation |
| Export lifecycle | DDEC-W07-007 | 7-day TTL + delete-on-regeneration | Security GAP-2 resolution |
| Temp-file-atomic-rename | DDEC-W07-008 | Write to temp, rename on success, clean up on failure | EDGE-024, security GAP-1 resolution |
| Max export size | DDEC-W07-009 | 100MB uncompressed hard limit | Security GAP-4 resolution |
| Photo in export | DDEC-W07-010 | Files in photos/ subfolder, named {checkInId}_{angle}.{ext} | planner consensus |
| ZIP format | DDEC-W07-011 | manifest.json (schemaVersion=1), data.json, summary.md, CSVs, photos/ | planner consensus |
| WAVE-03 stub pattern | DDEC-W07-012 | Empty arrays when WAVE-03 not deployed | WAVE-06 precedent |
| Sync generation | DDEC-W07-013 | Sync for MVP. Architecture supports future async. | User design decision |
| gqlgen bindings | DDEC-W07-014 | Add 16 type bindings to atlas-gqlgen.yml | architecture reviewer requirement |
| display_name | DDEC-W07-015 | NOT in user_profiles. Use atlas_users.display_name. | planner inconsistency resolution |

## Codebase Fit Decisions
- Follow WeekFlag model/repo/service/resolver pattern (not Settings pattern)
- AiExportDataProvider interface wraps all data source dependencies (prevents constructor explosion)
- GraphQL for CRUD mutations/queries; REST only for binary ZIP download
- ZIP generation in service package with in-memory struct export (testable, no filesystem deps)
- Bootstrap service extended to create default UserProfile (prevents 404 on first GET)

## Deferrals
- WeekFlagsByDateRange query — deferred to WAVE-04 or client-side per-week calls
- Build-time app version injection — omitted from manifest.json for MVP
- Max AiExport records per user cap — follow-up after MVP
- Async ZIP generation with status polling — follow-up after MVP
- Photo downscaling/resizing in export — follow-up optimization
- Concurrent generation lock — frontend debouncing sufficient for MVP
- WeekFlag service REST proxy for PAGE-009 — PAGE-009 uses WAVE-04 GraphQL directly

## Rejected Assumptions
- UserProfile can reuse WAVE-01 Settings: REJECTED. Separate entity per domain-model.md.
- include_photos DEFAULT true: REJECTED. Domain model invariant #10 requires false.
- CAP-W07-003 in WAVE-07: REJECTED. WAVE-04 owns week flags CRUD.
- display_name in user_profiles: REJECTED. Duplicates atlas_users.display_name.
- Migration 00093/00094: REJECTED. Correct numbers are 00091/00092.