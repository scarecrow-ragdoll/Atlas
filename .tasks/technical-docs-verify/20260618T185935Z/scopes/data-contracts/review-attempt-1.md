# Data-Contracts Review Attempt 1

## Verdict
needs-revision

## Sources Read
Same 6 sources as worker-attempt-1.md.

## Coverage Check
All product-verified source files consumed. Source delta (DEC-007, DEC-009) fully analyzed. Per-entity userId compliance table present. DEC-009 stale-name gaps identified.

## Evidence Check
All technical claims in the worker report trace to specific lines in source documents. No undocumented claims found.

## No-Invention Check
Worker suggests specific enum values and migration tools in "Suggested Decisions" — acceptable as opinion, not asserted as facts. No implementation contract invented. PASS.

## Source-Gap Consolidation Check
GENERAL PASS but three gaps are missing from the question ledger:
1. **Seed data / fixtures**: Listed in "Missing Source Artifacts" but no TQ-DATA question captures seed data format, default user fixture structure, or test data strategy.
2. **Index strategy**: Listed as missing but no question captures what indexes are needed for the performance p95 targets (product-brief.md §Performance Targets).
3. **Schema version format**: Listed as missing but no question captures the manifest.json schema version scheme needed for backup import validation.

These should be consolidated into at most two additional questions (seed data could join the migration strategy question; index strategy is a separate concern).

## Question Ledger Check
- TQ-DATA-001 through TQ-DATA-009 are properly structured with scope, severity, parent, and status columns.
- IDs follow TQ-DATA-* convention.
- Severities correctly assigned (dev-blocking for schema-blocking items, needs-owner-decision for design ambiguity, deferred/watchlist for non-blocking).
- Question TQ-DATA-001 (child entity userId depth) properly linked to DEC-007.
- **Missing question**: Stale WorkoutDay references (gap T3) have no question ID. Worker report says "no follow-up blockers discovered" but stale references are schema blockers — need a TQ-DATA-010.
- **Missing question**: Seed data strategy (gap T2/T5 overlap) — add or fold into TQ-DATA-005.
- **Missing question**: Index strategy for p95 performance targets — needs separate question.

## Answer Effect Check
No prior technical answers exist. DEC-007/DEC-009 source deltas analyzed correctly. The ambiguity on "all entities" depth (aggregate-root vs truly all) is correctly flagged as needs-owner-decision.

## Missing Or Unsupported Claims
1. Worker claims "All 20 entities listed" — verified true.
2. Worker claims "No prior docs/technical-verified/" — verified true.
3. Worker states DEC-007 trace to 3 questions — correct.
4. Worker does NOT include a question for stale WorkoutDay references — needs to add.

## Required Revisions
1. Add question TQ-DATA-010 for stale WorkoutDay references (dev-blocking — affects schema naming and code).
2. Add or extend TQ-DATA-005 to cover seed data / fixture format.
3. Add question TQ-DATA-011 for index strategy needed to meet p95 performance targets (watchlist — can be deferred but should be documented).
4. Merge TQ-DATA-007 and TQ-DATA-008 scope: both deal with data lifecycle cleanup. Consider renaming TQ-DATA-007 to "Media and data lifecycle policy" to widen coverage.

## Approval Notes
Not approved this round. Three new questions needed. One consolidation suggested. After these revisions, the report is structurally sound and sources are fully consumed.