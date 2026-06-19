# Question Ledger

## Open Questions

| ID | Scope | Severity | Parent | Question | Why It Matters | Needed Artifact Or Decision | Source | Status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| TQ-ARCH-001 | architecture | dev-blocking | none | No system context diagram | Cannot reason about system boundaries | Architecture diagram | worker-attempt-1 | resolved-by-decision |
| TQ-ARCH-002 | architecture | dev-blocking | none | No component architecture | No service/component boundaries | Component architecture document | worker-attempt-1 | resolved-by-decision |
| TQ-ARCH-003 | architecture | dev-blocking | none | Default user bootstrap mechanism | DEC-007 requires default user at install | Bootstrap specification | worker-attempt-1 | resolved-by-decision |
| TQ-ARCH-004 | architecture | dev-blocking | none | No deployment topology | Docker Compose structure undefined | Deployment diagram | worker-attempt-1 | resolved-by-decision |
| TQ-ARCH-005 | architecture | dev-blocking | none | No API surface boundaries | Public web vs admin boundaries undefined | API surface map | worker-attempt-1 | resolved-by-decision |
| TQ-ARCH-006 | architecture | dev-blocking | none | Go role undefined | GraphQL resolver or REST API provider? | Go service spec | worker-attempt-1 | resolved-by-decision |
| TQ-DATA-001 | data | dev-blocking | none | userId FK missing from 6 child entity descriptions | Breaks DEC-007 multi-user readiness | Data model alignment | worker-attempt-2 | resolved-by-decision |
| TQ-DATA-002 | data | dev-blocking | none | No index strategy for p95 query targets | Performance targets at risk | Index design doc | worker-attempt-2 | resolved-by-decision |
| TQ-DATA-003 | data | dev-blocking | none | 3 enum types undefined (source, heartRateZone, cardioType) | Data integrity risk | Enum definitions | worker-attempt-2 | resolved-by-decision |
| TQ-DATA-004 | data | dev-blocking | none | 3 enum types undefined (measurementType, flagType, mediaType) | Data integrity risk | Enum definitions | worker-attempt-2 | resolved-by-decision |
| TQ-DATA-005 | data | dev-blocking | none | No migration strategy | Schema changes risky | Migration strategy doc | worker-attempt-2 | resolved-by-decision |
| TQ-DATA-010 | data | dev-blocking | none | Import collision handling undefined | Data loss risk on re-import | Import collision spec | worker-attempt-2 | resolved-by-decision |
| TQ-DATA-006 | data | needs-owner | none | Stale WorkoutDay references after DEC-009 | Inconsistency risk | Owner decision | worker-attempt-2 | resolved-by-decision |
| TQ-DATA-007 | data | needs-owner | none | BodyWeightEntry.source enum undefined | Data integrity | Owner decision | worker-attempt-2 | resolved-by-decision |
| TQ-DATA-008 | data | deferred | none | Media lifecycle policy | Storage growth unbounded | Owner decision | worker-attempt-2 | deferred |
| TQ-DATA-009 | data | watchlist | none | Fixture/seed data format | Will be resolved by TQ-TEST-001 | - | worker-attempt-2 | resolved-by-decision |
| TQ-API-001 | api | dev-blocking | none | Primary API protocol undecided (GraphQL vs REST) | Blocks all endpoint work | Protocol decision | worker-attempt-1 | resolved-by-decision |
| TQ-API-002 | api | dev-blocking | TQ-API-001 | No endpoint catalog | All client work blocked | Endpoint/operation catalog | worker-attempt-1 | resolved-by-decision |
| TQ-API-003 | api | dev-blocking | TQ-API-001 | No error/validation response format | Client error handling blocked | Error contract | worker-attempt-1 | resolved-by-decision |
| TQ-API-007 | api | dev-blocking | TQ-API-001 | No file upload contract | Media upload blocked | Upload API spec | worker-attempt-1 | resolved-by-decision |
| TQ-API-008 | api | dev-blocking | TQ-API-001 | No session auth contract for API | Auth middleware blocked | Session contract | worker-attempt-1 | resolved-by-decision |
| TQ-API-010 | api | dev-blocking | TQ-API-001 | Chart data aggregation query undefined | Charts blocked | Chart query contract | worker-attempt-1 | resolved-by-decision |
| TQ-API-011 | api | dev-blocking | TQ-API-001 | Backup import flow API contract undefined | Import blocked | Import API contract | worker-attempt-1 | resolved-by-decision |
| TQ-API-012 | api | dev-blocking | TQ-API-001 | AI export operation contract undefined | Export blocked | Export operation spec | worker-attempt-1 | resolved-by-decision |
| TQ-API-004 | api | needs-owner | TQ-API-001 | Pagination contract | List queries undefined | Pagination decision | worker-attempt-1 | resolved-by-decision |
| TQ-API-006 | api | needs-owner | TQ-API-001 | Idempotency for mutations | Data integrity risk | Idempotency decision | worker-attempt-1 | resolved-by-decision |
| TQ-API-009 | api | needs-owner | TQ-API-001 | API versioning strategy | Future compatibility | Versioning decision | worker-attempt-1 | resolved-by-decision |
| TQ-AUTH-001 | auth | dev-blocking | none | PIN hash algorithm unspecified | Security risk | Hash algorithm decision | worker-attempt-1 | resolved-by-decision |
| TQ-AUTH-002 | auth | dev-blocking | none | PIN session TTL, cookie flags, renewal undefined | Session management blocked | Session spec | worker-attempt-1 | resolved-by-decision |
| TQ-AUTH-003 | auth | dev-blocking | none | No brute-force protection for PIN | Security risk | Brute-force spec | worker-attempt-1 | resolved-by-decision |
| TQ-AUTH-005 | auth | dev-blocking | TQ-AUTH-002 | Session token generation mechanism undefined | Session implementation blocked | Token spec | worker-attempt-1 | resolved-by-decision |
| TQ-AUTH-006 | auth | dev-blocking | none | Media access contradiction (RULE-022 vs RULE-024) | Access control undefined | Auth resolution | worker-attempt-1 | resolved-by-decision |
| TQ-AUTH-007 | auth | dev-blocking | none | Redis session store failure mode undefined | Session reliability | Redis failure mode | worker-attempt-1 | resolved-by-decision |
| TQ-AUTH-010 | auth | dev-blocking | TQ-ARCH-003 | DefaultUser bootstrap security | User creation during install | Bootstrap decision | worker-attempt-1 | resolved-by-decision |
| TQ-AUTH-011 | auth | dev-blocking | none | PIN: global vs per-app | Access scope undefined | PIN scope decision | worker-attempt-1 | resolved-by-decision |
| TQ-AUTH-004 | auth | needs-owner | none | Audit trail for sensitive operations | Compliance risk | Audit decision | worker-attempt-1 | resolved-by-decision |
| TQ-AUTH-009 | auth | needs-owner | none | Data retention/deletion policy | Compliance risk | Data policy | worker-attempt-1 | resolved-by-decision |
| TQ-INT-001 | integrations | dev-blocking | none | Sync/async contract for export/backup/import | Long ops undefined | Async decision | worker-attempt-2 | resolved-by-decision |
| TQ-INT-002 | integrations | dev-blocking | TQ-INT-001 | Progress reporting for long operations | UX blocked for export | Progress spec | worker-attempt-2 | resolved-by-decision |
| TQ-INT-003 | integrations | dev-blocking | none | Import transaction spec for RULE-009 | Data integrity on import | Transaction spec | worker-attempt-2 | resolved-by-decision |
| TQ-INT-004 | integrations | needs-owner | none | File size limits | UX/resource limits | Size limits decision | worker-attempt-2 | resolved-by-decision |
| TQ-INT-005 | integrations | needs-owner | TQ-INT-001 | Session continuity during long ops | User may close tab | Continuity decision | worker-attempt-2 | resolved-by-decision |
| TQ-CLIENT-001 | client-ux | dev-blocking | none | No UI state machine for any page | All pages blocked | State machine spec | worker-attempt-1 | resolved-by-decision |
| TQ-CLIENT-002 | client-ux | dev-blocking | TQ-CLIENT-001 | No page loading state contract | UX blocked | Loading state spec | worker-attempt-1 | resolved-by-decision |
| TQ-CLIENT-003 | client-ux | dev-blocking | TQ-CLIENT-001 | No empty state contract | UX blocked | Empty state spec | worker-attempt-1 | resolved-by-decision |
| TQ-CLIENT-004 | client-ux | dev-blocking | TQ-CLIENT-001 | No error state contract | Error UX blocked | Error state spec | worker-attempt-1 | resolved-by-decision |
| TQ-CLIENT-005 | client-ux | dev-blocking | TQ-CLIENT-001 | No offline state contract | Offline UX blocked | Offline state spec | worker-attempt-1 | resolved-by-decision |
| TQ-CLIENT-006 | client-ux | dev-blocking | none | No form validation contract | All forms blocked | Validation contract | worker-attempt-1 | resolved-by-decision |
| TQ-CLIENT-007 | client-ux | dev-blocking | none | No CRUD feedback pattern | Save UX undefined | Feedback pattern | worker-attempt-1 | resolved-by-decision |
| TQ-CLIENT-008 | client-ux | dev-blocking | none | No client error display pattern | Errors undefined | Error display spec | worker-attempt-1 | resolved-by-decision |
| TQ-CLIENT-009 | client-ux | dev-blocking | none | No cache/data freshness contract | Performance targets at risk | Cache strategy | worker-attempt-1 | resolved-by-decision |
| TQ-CLIENT-010 | client-ux | dev-blocking | none | No navigation state management | Unsaved data loss risk | Navigation state spec | worker-attempt-1 | resolved-by-decision |
| TQ-CLIENT-011 | client-ux | dev-blocking | none | No accessibility standard | WCAG/ARIA undefined | Accessibility standard | worker-attempt-1 | resolved-by-decision |
| TQ-CLIENT-012 | client-ux | dev-blocking | none | No localization strategy | i18n undefined | Localization strategy | worker-attempt-1 | resolved-by-decision |
| TQ-OPS-001 | operations | dev-blocking | none | No environment/config topology | Deployment blocked | Config spec | worker-attempt-1 | resolved-by-decision |
| TQ-OPS-002 | operations | dev-blocking | none | No logging framework or health check | Observability blocked | Logging spec | worker-attempt-1 | resolved-by-decision |
| TQ-OPS-003 | operations | dev-blocking | none | No SLO/SLI/metrics for DEC-008 p95 targets | Performance unmeasurable | Metrics spec | worker-attempt-1 | resolved-by-decision |
| TQ-OPS-005 | operations | dev-blocking | none | No backup scheduling/retention/verification | Operability blocked | Backup ops spec | worker-attempt-1 | resolved-by-decision |
| TQ-OPS-006 | operations | needs-owner | none | No resource estimates | Capacity planning impossible | Resource estimates | worker-attempt-2 | resolved-by-decision |
| TQ-TEST-001 | testing | dev-blocking | none | No test data factory/fixture strategy | Tests blocked | Fixture strategy | worker-attempt-1 | resolved-by-decision |
| TQ-TEST-002 | testing | dev-blocking | none | Weekly workflow e2e strategy undefined | Quality gate blocked | E2E plan | worker-attempt-1 | resolved-by-decision |
| TQ-TEST-003 | testing | dev-blocking | none | AI export schema snapshot test strategy undefined | Quality gate blocked | Snapshot strategy | worker-attempt-1 | resolved-by-decision |
| TQ-TEST-004 | testing | dev-blocking | none | Backup manifest schema test strategy undefined | Quality gate blocked | Schema test strategy | worker-attempt-1 | resolved-by-decision |
| TQ-TEST-005 | testing | needs-owner | none | Log redaction test policy | Quality gate blocked | Log test policy | worker-attempt-1 | resolved-by-decision |
| TQ-TEST-006 | testing | needs-owner | none | Test isolation strategy | Test reliability | Isolation decision | worker-attempt-1 | resolved-by-decision |
| TQ-TEST-007 | testing | needs-owner | none | Performance test policy for DEC-008 | Quality gate blocked | Perf test policy | worker-attempt-1 | resolved-by-decision |

## Answered Questions

| ID | Scope | Answer | Resolved By |
| --- | --- | --- | --- |
| TQ-API-004 | api | Cursor-based pagination, default 50, max 200, stable ordering | Owner decision (TDEC-001) |
| TQ-API-006 | api | Idempotency for long-running ops; optional clientMutationId for CRUD | Owner decision (TDEC-002) |
| TQ-API-009 | api | REST /api/v1; GraphQL evolutionary; export schema 1.0.0 | Owner decision (TDEC-003) |
| TQ-AUTH-004 | auth | Minimal audit trail with specific events and no-log list | Owner decision (TDEC-004) |
| TQ-AUTH-009 | auth | No auto-deletion in MVP; specific deletion rules per entity | Owner decision (TDEC-005) |
| TQ-DATA-006 | data | Remove all stale WorkoutDay references; use DailyLog canonical | Owner decision (TDEC-006) |
| TQ-DATA-007 | data | BodyWeightEntry.source enum: MANUAL, DAILY_LOG, CHECK_IN, BACKUP_IMPORT, EXTERNAL_IMPORT; APPLE_HEALTH_IMPORT future | Owner decision (TDEC-007) |
| TQ-INT-004 | integrations | Image 25MB, video 250MB, single upload 300MB, AI export 500MB soft limit; allowed formats | Owner decision (TDEC-008) |
| TQ-INT-005 | integrations | JobId-based continuation; session-validated downloads | Owner decision (TDEC-009) |
| TQ-OPS-006 | operations | Min 1vCPU/1GB/10GB, recommended 2vCPU/4GB/50GB | Owner decision (TDEC-010) |
| TQ-TEST-005 | testing | Log redaction must be tested; specific no-log list | Owner decision (TDEC-011) |
| TQ-TEST-006 | testing | Isolated DB/Redis/media per test type; Docker Compose e2e | Owner decision (TDEC-012) |
| TQ-TEST-007 | testing | Performance tests for critical paths; CI smoke + release verification | Owner decision (TDEC-013) |

## Follow-Up Questions

None resolved — answered questions reduced the open count without spawning new dev-blocking questions. Child questions parented to TQ-API-001 (endpoints, error format, upload contract, session auth, chart queries, backup import, AI export) remain dev-blocking because they need endpoint-specific specs, but the parent protocol decision is resolved.

## Resolved Questions

| ID | Scope | Resolution |
| --- | --- | --- |
| TQ-API-004 | api | Resolved by TDEC-001 |
| TQ-API-006 | api | Resolved by TDEC-002 |
| TQ-API-009 | api | Resolved by TDEC-003 |
| TQ-AUTH-004 | auth | Resolved by TDEC-004 |
| TQ-AUTH-009 | auth | Resolved by TDEC-005 |
| TQ-DATA-006 | data | Resolved by TDEC-006 |
| TQ-DATA-007 | data | Resolved by TDEC-007 |
| TQ-INT-004 | integrations | Resolved by TDEC-008 |
| TQ-INT-005 | integrations | Resolved by TDEC-009 |
| TQ-OPS-006 | operations | Resolved by TDEC-010 |
| TQ-TEST-005 | testing | Resolved by TDEC-011 |
| TQ-TEST-006 | testing | Resolved by TDEC-012 |
| TQ-TEST-007 | testing | Resolved by TDEC-013 |

## Deferred Questions

See open-questions.md §Deferred for list of deferred questions.