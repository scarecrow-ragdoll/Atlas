# Integrations-Events Worker Attempt 1

## Sources Read
- docs/product-verified/functional-spec.md
- docs/product-verified/business-rules.md
- docs/product-verified/product-brief.md
- docs/product-verified/scope.md
- docs/product-verified/edge-cases.md
- docs/product-verified/open-questions.md
- docs/product-verified/features/ai-export.md
- docs/product-verified/features/backup-and-restore.md
- docs/product-verified/features/ai-prompt-builder.md
- docs/product-verified/features/index.md

## Source Delta Reviewed
None (initial run).

## Product Signals

### External Systems
- "No external API integrations in MVP" (§22 functional-spec.md)
- "AI analysis is manual copy-paste (user sends prompt + ZIP to ChatGPT)" (§22)
- "go-telegram/bot library is explicitly out of MVP scope (§23)"
- "Apple Health is explicitly out of MVP scope (§22)"
- "No automatic external API calls in MVP" (RULE-029)
- Target AI platforms (ChatGPT only or Claude, Gemini, local LLMs) unresolved (Q-SCOPE-006)

### Async Operations (User-Triggered)
- **AI Export**: Configurable date range, ZIP generation with CSVs, JSON, markdown, photos.
  - Performance: 4w no photos ≤5s, 4w with photos ≤20s, 12mo no photos ≤15s, 12mo with photos "best effort, show progress" (p95 product-brief.md)
  - UX rule: operations >2s must show loading state
- **Backup Export**: Full ZIP with manifest.json, data.json, media/
  - Performance: db-only ≤15s, with media "best effort (media-size dependent)"
- **Backup Import**: Upload ZIP → validate → dry-run → summary → confirm → restore
  - Performance: dry-run db-only ≤15s, import db-only ≤30s
  - All-or-nothing: "complete success or clean failure"

### Failure Handling Signals
- "Data not lost on restart" (§24.2)
- EDGE-021: mid-import failure — transaction rollback undefined
- EDGE-024: disk full during export generation — no quota handling
- EDGE-022: PostgreSQL connection lost during save
- EDGE-023: Redis unavailable for session store
- EDGE-025: Docker volume full — media save fails
- EDGE-028: Schema migration after backup import from older version

### Constraints
- Single self-hosted user — no multi-tenancy rate-limit concerns
- Redis in stack (scope.md), available for session store and optional job progress
- Docker + filesystem volume for media

## Technical Facts

1. **No external API integrations in MVP**: The product explicitly excludes API calls to OpenAI, Telegram, Apple Health, or any third-party service. All AI workflows are manual copy-paste. Zero webhooks, zero outbound API calls, zero inbound webhooks.

2. **No job queue specified**: No queuing system, event bus, or background worker architecture is defined. Redis is available but no job/queue usage is documented.

3. **Export/Backup generation is user-triggered synchronous or async-adjacent**: Performance targets distinguish synchronous-looking targets (≤5s, ≤15s) from "best effort, show progress" (large exports). This implies an undefined async path with progress reporting.

4. **Import has explicit rollback requirement**: "no silent partial import — complete success or clean failure" (RULE-009) and "Import must complete fully or fail" (business-rules.md). No transaction/rollback mechanism specified.

5. **Redis available but role undefined**: Scope.md lists Redis in stack, but only session store (PIN guard) is a known use case. No queue, rate-limit, or job-progress use documented.

6. **Rate limiting unspecified**: No rate limits on export, backup, or any operation. Single-user self-hosted context makes this low risk but not zero (e.g., user could accidentally trigger multiple large exports).

## Technical Gaps

### Missing Artifact: Async Job / Progress Contract
The product requires progress indication for large exports ("show progress" for 12mo with photos) but defines no async job mechanism, progress API, or polling contract. Without this, "show progress" cannot be implemented.

### Missing Artifact: Export/Backup Timeout Specification
No timeouts defined for any export or backup operation. A 12-month export with photos on a slow filesystem could hang indefinitely.

### Missing Artifact: Failure Recovery Model
Mid-import failure (EDGE-021) requires "clean failure" but no transaction model, partial-cleanup mechanism, or rollback approach is defined. Database-level transaction vs application-level compensation is unspecified.

### Missing Artifact: Disk Space Pre-Check
EDGE-024 and EDGE-025 identify disk-full failure paths with no pre-check or graceful error handling. No size estimation or capacity validation before starting export/import.

### Missing Artifact: Session Continuity for Long Operations
PIN session TTL is unspecified (Q-ROLE-001 / EDGE-012). Long export/import operations could exceed session lifetime, causing inconsistent behavior (operation continues server-side but session expires, or operation fails mid-way).

### Missing Artifact: Rate Limit / Abuse Prevention
No guardrails against repeated export/backup triggers. Low risk for single-user, but no cooldown, max-frequency, or concurrent-operation limit defined.

## Questions Raised

| ID | Question | Severity | Why It Matters |
| --- | --- | --- | --- |
| TQ-INT-001 | Is AI Export ZIP generation synchronous or asynchronous, and what progress reporting mechanism is used for large exports? | dev-blocking | Performance targets specify "show progress" for 12mo with photos — opposite of sync. No progress API defined. |
| TQ-INT-002 | Is Backup ZIP generation synchronous or asynchronous? | dev-blocking | Backup with media is "best effort" with no progress or timeout contract. |
| TQ-INT-003 | Is Backup Import synchronous or asynchronous, and what happens on mid-import failure beyond "no silent partial import"? | dev-blocking | EDGE-021: transaction rollback undefined. Recovery model absent. |
| TQ-INT-004 | What are the file size limits and disk space pre-checks for AI Export and Backup ZIP generation? | needs-owner-decision | EDGE-024: disk full during export has no pre-check or graceful handling. |
| TQ-INT-005 | What is the session timeout for PIN guard, and how does it interact with long-running export/import operations? | needs-owner-decision | Operation could exceed session lifetime; inconsistent behavior risk. |
| TQ-INT-006 | What rate limits or abuse prevention apply to export/backup operations? | deferred | Low risk for single-user but no guardrail defined. |

## Answer Effects
None — initial run.

## Risks

1. **Async contract missing blocks use of "show progress" in UI**: Without a progress API, the performance target for large exports is unimplementable.
2. **Import transaction model undefined risks data corruption**: "No silent partial import" is a strong requirement; without transaction specification, implementation may not meet it.
3. **Disk full failure is an operational blind spot**: Mid-operation disk full causes undefined behavior.
4. **Session timeout + long operation coupling**: If PIN session expires during import, the user may see a redirect to PIN screen mid-import while the server continues processing.

## Suggested Decisions

1. DEC: AI Export should use synchronous generation with a progress-polling endpoint for exports exceeding ~5s estimated time. Use Redis to store progress state.
2. DEC: Backup Export should be synchronous with a timeout (e.g., 120s for db-only, best-effort with progress for media).
3. DEC: Backup Import should use database transactions with a full-rollback on any failure. Wrap the entire import in a single transaction or use a compensating cleanup step.
4. DEC: Disk space should be checked before export/import starts; fail early with a user-facing error message.
5. DEC: Export/Backup operations should extend or bypass the PIN session TTL for their duration.
6. DEC: Allow one export/backup operation at a time; queue or block concurrent triggers.

## Traceability Candidates

- TQ-INT-001 → product-brief.md §AI Export (p95 targets), functional-spec.md §AI Export
- TQ-INT-002 → product-brief.md §Backup (p95 targets)
- TQ-INT-003 → edge-cases.md EDGE-021, business-rules.md RULE-009
- TQ-INT-004 → edge-cases.md EDGE-024, EDGE-025
- TQ-INT-005 → open-questions.md Q-ROLE-001, edge-cases.md EDGE-012
- TQ-INT-006 → (implied by single-user operation safety)