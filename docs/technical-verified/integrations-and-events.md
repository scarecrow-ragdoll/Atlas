# Integrations And Events

## External Systems

No external API integrations in MVP. AI analysis is manual copy-paste. All integration points are internal:
- Browser ↔ Go API (undecided protocol)
- Go API ↔ PostgreSQL
- Go API ↔ Redis (sessions)
- Go API ↔ Filesystem (media, exports, backups)

## Events And Jobs

**Decision (TDEC-009):** Long operations continue after being started. Each operation receives a `jobId`. Progress available via `GET /api/v1/jobs/{id}`. If session expires before download, user must unlock with PIN again. Completed download URLs not public — must validate current session.

No event system or background job infrastructure defined. Three user-triggered operations require async/background execution:
- AI export ZIP generation (TQ-INT-001)
- Full backup ZIP generation
- Full backup import

Sync vs async contract for these operations is undefined (TQ-INT-001). Progress reporting mechanism absent (TQ-INT-002).

## Sync And Retry Rules

- No transaction specification for "no silent partial import" (RULE-009) — TQ-INT-003
- No retry/backoff for long-running operations
- No dead letter or reconciliation for failed jobs

## Rate Limits And Failure Handling

**Decision (TDEC-008):**
- Progress photo image: 25 MB per file
- Exercise image: 25 MB per file
- Exercise video: 250 MB per file
- Single upload request: 300 MB
- AI export without photos: no practical limit
- AI export with photos: 500 MB soft limit
- Full backup with media: best effort, disk-space limited
- Allowed formats: JPEG/PNG/WEBP (images), MP4/MOV/WEBM (video)
- Exceeding limit returns clear validation error

Low risk for single-user. Rate limiting deferred (TQ-INT-006).

## Integration Questions

| ID | Question | Severity | Status |
| --- | --- | --- | --- |
| TQ-INT-001 | Async vs sync contract for AI export, backup, import | dev-blocking | **resolved** (TDEC-041) |
| TQ-INT-002 | Progress reporting mechanism for long operations | dev-blocking | **resolved** (TDEC-042) |
| TQ-INT-003 | "No silent partial import" — transaction specification missing | dev-blocking | **resolved** (TDEC-043) |
| TQ-INT-004 | File size limits for export/import/media upload | needs-owner | **resolved** (TDEC-008) |
| TQ-INT-005 | Session continuity for browser during long operations | needs-owner | **resolved** (TDEC-009) |
| TQ-INT-006 | Rate limiting (low risk, single-user) | deferred | deferred |
| TQ-INT-007 | Redis-based session failure model | watchlist | open |
| TQ-INT-008 | ZIP memory allocation strategy for large exports | watchlist | open |