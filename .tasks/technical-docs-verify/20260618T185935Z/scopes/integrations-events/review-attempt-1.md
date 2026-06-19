# Integrations-Events Review Attempt 1

## Verdict
needs-revision

## Sources Read
- All sources listed in worker-attempt-1.md
- docs/product-verified/features/ai-export.md
- docs/product-verified/features/backup-and-restore.md

## Coverage Check
- External systems: covered — none in MVP, correctly identified.
- Sync ownership: covered — all 3 user-triggered operations identified.
- Async jobs: covered — export/backup/import identified as async-adjacent operations.
- Events: covered — none defined, correctly called out.
- Webhooks: covered — none, explicitly excluded.
- Queues: covered — not needed, correctly assessed.
- Retries/backoff/dead letters: covered — not applicable for MVP.
- Reconciliation: covered — not applicable.
- Rate limits: covered — deferred as low risk, reasonable.
- Failure handling: covered — disk full, session timeout, mid-import, DB connection loss identified.

## Evidence Check
All factual claims trace to product-verified sources. Performance targets cite product-brief.md §AI Export and §Backup. Edge cases cite edge-cases.md by ID. Business rules cite business-rules.md by RULE-ID.

No claim is unsupported.

## No-Invention Check
- No endpoints, schemas, event payloads, auth rules, infra topology, SLOs, migrations, or test gates invented. ✓
- Suggested Decisions section is clearly labeled as decisions, not facts. ✓
- Questions are framed as gaps, not assumptions. ✓

## Source-Gap Consolidation Check
Gaps are consolidated into 6 distinct questions covering:
1. Async job/progress contract (TQ-INT-001)
2. Backup sync vs async + timeout (TQ-INT-002)
3. Import transaction model + rollback (TQ-INT-003)
4. Disk space pre-check (TQ-INT-004)
5. Session continuity (TQ-INT-005)
6. Rate limits (TQ-INT-006)

No redundant splitting. Each addresses a distinct missing artifact class.

## Question Ledger Check
- IDs use TQ-INT-* prefix ✓
- Severity levels match allowed values (dev-blocking, needs-owner-decision, deferred) ✓
- Status field present with all values "open" (correct for initial run) ✓
- "Why It Matters" and "Needed Artifact Or Decision" columns populated ✓

## Answer Effect Check
No answers to review (initial run).

## Missing Or Unsupported Claims

1. **Redis dependency for progress tracking not addressed as a gap**: Suggested Decision 1 proposes using Redis for progress polling during AI Export, but EDGE-023 (Redis unavailable) is only cited as a session-store concern in the worker. If Redis is used for progress, Redis unavailability becomes a blocking dependency for the export feature itself. This should be captured explicitly — either as a follow-up question or as a risk expansion.

2. **No discussion of ZIP streaming vs in-memory buffering for large exports**: For 12-month exports with photos, in-memory ZIP generation can exhaust RAM. This is an implementation detail, but the performance target ("show progress") and the risk of OOM make it relevant to the async/scalability discussion. At minimum, the risk should be flagged.

## Required Revisions

1. Expand the "Redis dependency for progress" concern — either add a note in the risk section about Redis unavailability blocking export progress, or raise a watchlist-level question.
2. Add a note about ZIP generation memory pressure for large exports in the risks section.

## Approval Notes
Will approve after revisions above are incorporated. Remaining analysis is thorough, evidence-backed, and does not invent implementation contracts.