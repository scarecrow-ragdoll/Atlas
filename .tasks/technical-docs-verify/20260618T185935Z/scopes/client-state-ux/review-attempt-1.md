# client-state-ux Review Attempt 1

## Verdict
**needs-revision**

## Sources Read
- worker-attempt-1.md
- docs/product-verified/functional-spec.md
- docs/product-verified/user-flows.md
- docs/product-verified/edge-cases.md
- docs/product-verified/product-brief.md
- docs/product-verified/scope.md
- docs/product-verified/domain-model.md
- docs/product-verified/features/pin-guard.md
- docs/product-verified/features/dashboard.md
- docs/product-verified/features/ai-export.md
- docs/product-verified/features/backup-and-restore.md

## Coverage Check
- UI state machines: covered (TQ-CLIENT-001)
- Loading states: covered (TQ-CLIENT-002)
- Empty states: covered (TQ-CLIENT-003)
- Error states: covered (TQ-CLIENT-004)
- Offline states: covered (TQ-CLIENT-005)
- Form validation: covered (TQ-CLIENT-006, TQ-CLIENT-007)
- Optimistic updates: covered (TQ-CLIENT-008)
- Realtime updates: covered (TQ-CLIENT-010)
- Cache invalidation: covered (TQ-CLIENT-009)
- Accessibility: covered (TQ-CLIENT-011)
- Localization: covered (TQ-CLIENT-012)
- Session/recovery: covered (TQ-CLIENT-013)
- SPA routing: covered (TQ-CLIENT-014)

Missing coverage:
- PIN guard specific UX states (retry countdown, lockout timer, brute-force delay display) — partially addressed in suggested decisions but not as a formal question
- Empty state shape for PIN lockout screen
- Calendar navigation UX (month-to-month loading, year-to-year)
- Nutrition template auto-apply UX (mid-week creation — show week preview or apply forward only)

## Evidence Check
- Every technical claim traces to a source document, edge case, or product question. Pass.
- Source delta DEC-008 is correctly referenced in loading state gap and performance targets.
- Product questions Q-ACTOR-12, Q-FEAT-009, Q-FEAT-005, Q-ROLE-001, EDGE-* are correctly traced.

## No-Invention Check
- The worker does not invent API endpoints, schemas, or implementation details for gaps.
- Suggested decisions are clearly labeled as suggestions, not requirements.
- Missing artifact classes are correctly consolidated. PASS.

## Source-Gap Consolidation Check
- Missing artifacts are consolidated into 15 classes. This is appropriate but could be slightly tighter:
  - Gap 1 (state machine) + Gap 3 (loading) + Gap 4 (error) + Gap 14 (export/import UX) could be one "unified UX state contract" gap.
  - SPA routing (Gap 12) partially overlaps with state machine behavior.
- OK to keep separate for clarity given different scope owner responsibilities. PASS (minor suggestion).

## Question Ledger Check
- 14 questions, 11 dev-blocking, 3 needs-owner-decision. All use TQ-CLIENT-* prefix.
- Parent links correct (TQ-CLIENT-002, 003, 004 → TQ-CLIENT-001; TQ-CLIENT-007 → TQ-CLIENT-006).
- Statuses: all open. Correct.
- Needs revision:
  - **TQ-CLIENT-013 severity is dev-blocking, not needs-owner-decision.** Session loss during data entry can corrupt or lose user data. This is an implementation blocker.
  - **Missing question: PIN brute-force lockout UX** — the worker suggests a lockout policy in "Suggested Decisions" but does not formalize a question. If the PIN guard locks the user out, the UI must display lockout state, countdown, and recovery path. Add as TQ-CLIENT-015.
  - **Missing question: Calendar navigation loading UX** — Workout Diary uses calendar navigation to any past date. Loading behavior for month-to-month pagination, year switching, and date selection is unspecified. Add as TQ-CLIENT-016.
  - **Missing question: Nutrition template mid-week creation UX** — EDGE-017: template created mid-week. The UX for week-preview and forward-only vs retroactive application needs design. Add as TQ-CLIENT-017.
  - **TQ-CLIENT-010 (realtime) could be deferred** — single user, self-hosted, non-realtime is acceptable. Lower to deferred severity, not needs-owner-decision. Or merge into cache question TQ-CLIENT-009.

## Answer Effect Check
No answers to check. Source delta DEC-008 was reviewed. PASS.

## Missing Or Unsupported Claims
None. All claims are supported by source evidence.

## Required Revisions

1. **Change TQ-CLIENT-013 severity from `needs-owner-decision` to `dev-blocking`** — Session loss UX is an implementation blocker.
2. **Add TQ-CLIENT-015: PIN brute-force lockout UX** — Lockout state, timer display, retry countdown display are unspecified.
3. **Add TQ-CLIENT-016: Calendar navigation UX for workout diary** — Month loading, year switching, date selection feedback not specified.
4. **Add TQ-CLIENT-017: Nutrition template mid-week creation UX** — Week preview and application boundary (forward vs retroactive) not specified.
5. **Consider lowering TQ-CLIENT-010 to `deferred`** — Single-user self-hosted app can accept manual tab refresh. If not deferred, adjust severity to `needs-owner-decision` is acceptable.
6. **Minor: Add PIN-guard-specific states to worker text description** — idle->pending->success/error/locked states for the PIN entry screen are not called out in the gap description, even though they follow the general state machine pattern.

## Approval Notes
Worker is well-researched and correctly identifies all major client-state-ux gaps. Three ledger questions are missing, one severity is wrong. After revisions, this scope should be ready for approval.