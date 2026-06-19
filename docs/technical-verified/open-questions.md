# Open Questions

## Dev-Blocking

None — all dev-blocking questions resolved by owner decisions TDEC-014..059.

## Needs Owner Decision

None — all needs-owner questions resolved by TDEC-001..013.

## Deferred

## Deferred

| ID | Question | Owner |
| --- | --- | --- |
| TQ-DATA-008 | Media lifecycle policy | product owner |
| TQ-API-005 | Filtering contract | product owner |
| TQ-OPS-004 | Alerting/runbooks | product owner |
| TQ-INT-006 | Rate limiting | product owner |
| TQ-AUTH-008 | Backup identity validation | product owner |

## Watchlist

| ID | Question | Rationale |
| --- | --- | --- |
| TQ-INT-007 | Redis-based session failure model | Single-user, low risk |
| TQ-INT-008 | ZIP memory allocation for large exports | Needs monitoring post-launch |
| TQ-DATA-009 | Fixture/seed data format | Will be resolved by TQ-TEST-001 |
| TQ-API-013 | Health check endpoint | Nice-to-have for Docker |
| TQ-CLIENT-011 | Accessibility standard | Deferred to post-MVP |
| TQ-CLIENT-012 | Localization strategy | Deferred to post-MVP |

## Resolved This Run

All questions resolved. 59 TDEC entries in decision-log.md covering: architecture (6), data (6), API (7 + 3 from previous batch), auth (8 + 2), integrations (3 + 2), client UX (12), operations (4 + 1), testing (4 + 3).

## Dev-Blocking Remaining

None. All ~80 technical questions resolved.

| ID | Question | Resolution |
| --- | --- | --- |
| TQ-API-001 | GraphQL vs REST | Hybrid: GraphQL primary API, REST for binary/long ops. See api-contracts.md |
| TQ-API-004 | Pagination contract | Cursor-based, default 50, max 200 (TDEC-001) |
| TQ-API-006 | Idempotency for mutations | Required for long ops, optional clientMutationId for CRUD (TDEC-002) |
| TQ-API-009 | API versioning strategy | REST /api/v1, GraphQL evolutionary, export schema 1.0.0 (TDEC-003) |
| TQ-AUTH-004 | Audit trail | Yes, minimal audit for PIN/media/export/backup events (TDEC-004) |
| TQ-AUTH-009 | Data retention/deletion policy | No auto-deletion; entity-specific rules (TDEC-005) |
| TQ-DATA-006 | Stale WorkoutDay references | Resolved: canonical name is DailyLog everywhere (TDEC-006) |
| TQ-DATA-007 | BodyWeightEntry.source enum | Defined: MANUAL, DAILY_LOG, CHECK_IN, BACKUP_IMPORT, EXTERNAL_IMPORT (TDEC-007) |
| TQ-INT-004 | File size limits | Image 25MB, video 250MB, upload 300MB, AI export 500MB (TDEC-008) |
| TQ-INT-005 | Session continuity | JobId + progress polling + session-validated downloads (TDEC-009) |
| TQ-OPS-006 | Resource estimates | Min 1vCPU/1GB/10GB, recommended 2vCPU/4GB/50GB (TDEC-010) |
| TQ-TEST-005 | Log redaction test policy | Must be tested with specific no-log list (TDEC-011) |
| TQ-TEST-006 | Test isolation strategy | Isolated DB/Redis/media per type (TDEC-012) |
| TQ-TEST-007 | Performance test policy | Critical paths in CI smoke + release verification (TDEC-013) |

## Dev-Blocking Remaining

~47 dev-blocking questions remain. Key additions from resolved parent:
- TQ-API-001 (protocol) resolved → TQ-API-002 (endpoints), TQ-API-003 (errors), TQ-API-007 (upload), TQ-API-008 (auth), TQ-API-010 (charts), TQ-API-011 (backup), TQ-API-012 (AI export) still open but unblocked by protocol decision
- TQ-API-004, TQ-API-006, TQ-API-009 moved from needs-owner to resolved
- All 13 needs-owner questions resolved