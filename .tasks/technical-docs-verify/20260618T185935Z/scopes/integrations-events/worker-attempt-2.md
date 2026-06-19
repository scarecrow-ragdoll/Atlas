# Integrations-Events Worker Attempt 2

## Sources Read
(same as attempt 1)

## Source Delta Reviewed
None.

## Previous Reviewer Findings Addressed
1. Redis dependency for progress tracking — expanded in Risks and added follow-up question.
2. ZIP generation memory pressure — added to Risks.

## Product Signals
(same as attempt 1)

## Technical Facts
(same as attempt 1, with additions:)

6. **Redis dependency for progress is a coupling risk**: If Suggested Decision 1 (Redis-backed progress polling) is adopted, then Redis unavailability (EDGE-023) blocks not just session management but the AI Export progress reporting feature itself.
7. **ZIP generation memory profile**: Large exports with photos may generate ZIP archives exceeding available RAM if buffered in memory. Streaming ZIP generation or temp-file-based approaches mitigate this but add complexity.

## Technical Gaps
(same as attempt 1, plus:)

### Missing Artifact: Progress-Reporting Dependency Model
If progress reporting uses Redis (via Suggested Decision 1), the dependency model must specify behavior when Redis is unavailable: degrade to no progress, block the operation, or use an alternative progress mechanism.

### Missing Artifact: Export Memory Budget
No specification of whether ZIP generation uses in-memory buffering, streaming, or temp-file staging. For 12-month exports with photos, this is a correctness and reliability concern.

## Questions Raised

| ID | Question | Severity | Parent | Why It Matters |
| --- | --- | --- | --- | --- |
| TQ-INT-001 | Is AI Export ZIP generation synchronous or asynchronous, and what progress reporting mechanism is used for large exports? | dev-blocking | none | Performance targets specify "show progress" for 12mo with photos — no progress API defined. |
| TQ-INT-002 | Is Backup ZIP generation synchronous or asynchronous? | dev-blocking | none | Backup with media is "best effort" with no progress or timeout contract. |
| TQ-INT-003 | Is Backup Import synchronous or asynchronous, and what happens on mid-import failure beyond "no silent partial import"? | dev-blocking | none | EDGE-021: transaction rollback undefined. Recovery model absent. |
| TQ-INT-004 | What are the file size limits and disk space pre-checks for AI Export and Backup ZIP generation? | needs-owner-decision | none | EDGE-024: disk full during export has no pre-check or graceful handling. |
| TQ-INT-005 | What is the session timeout for PIN guard, and how does it interact with long-running export/import operations? | needs-owner-decision | none | Operation could exceed session lifetime; inconsistent behavior risk. |
| TQ-INT-006 | What rate limits or abuse prevention apply to export/backup operations? | deferred | none | Low risk for single-user but no guardrail defined. |
| TQ-INT-007 | If progress reporting uses Redis (e.g., polling endpoint), what is the behavior when Redis is unavailable during export/import? | watchlist | TQ-INT-001 | EDGE-023 (Redis unavailable) affects progress reporting if Redis is the progress store. Operation should not silently hang. |
| TQ-INT-008 | What is the ZIP generation memory strategy (in-memory, streaming, or temp-file) for large exports with photos? | watchlist | TQ-INT-001 | 12-month exports with photos can exceed available RAM if fully buffered in memory. |

## Answer Effects
None.

## Risks

1. **Async contract missing blocks "show progress" in UI** — unchanged.
2. **Import transaction model undefined risks data corruption** — unchanged.
3. **Disk full failure is an operational blind spot** — unchanged.
4. **Session timeout + long operation coupling** — unchanged.
5. **NEW: Redis unavailability blocks progress reporting**: If progress is coupled to Redis, an unavailable Redis (EDGE-023) silently blocks progress visibility. Consider fallback: synchronous fallback, static progress estimate, or at least an error message.
6. **NEW: In-memory ZIP buffering can OOM for large exports**: A 12-month export including all photos may generate a multi-GB ZIP. In-memory generation will crash the process. Streaming or temp-file staging required.

## Suggested Decisions

1. DEC: AI Export should use synchronous generation with a progress-polling endpoint for exports exceeding ~5s estimated time. Use Redis to store progress state.
2. DEC: If Redis is used for progress tracking and Redis is unavailable, fall back to synchronous generation without progress updates rather than blocking the user.
3. DEC: Backup Export should be synchronous with a timeout (e.g., 120s for db-only, best-effort with progress for media).
4. DEC: Backup Import should use database transactions with a full-rollback on any failure.
5. DEC: Disk space should be checked before export/import starts; fail early with a user-facing error message.
6. DEC: Export/Backup operations should extend or bypass the PIN session TTL for their duration.
7. DEC: Use streaming ZIP generation (Write-ZIP-to-temp-file or pipe) rather than in-memory buffering for AI Export and Backup.
8. DEC: Allow one export/backup operation at a time; queue or block concurrent triggers.

## Traceability Candidates

- TQ-INT-001 → product-brief.md §AI Export
- TQ-INT-002 → product-brief.md §Backup
- TQ-INT-003 → edge-cases.md EDGE-021, business-rules.md RULE-009
- TQ-INT-004 → edge-cases.md EDGE-024, EDGE-025
- TQ-INT-005 → open-questions.md Q-ROLE-001
- TQ-INT-006 → (implied)
- TQ-INT-007 → edge-cases.md EDGE-023, DEC-[Redis fallback]
- TQ-INT-008 → product-brief.md §AI Export (12mo with photos target)