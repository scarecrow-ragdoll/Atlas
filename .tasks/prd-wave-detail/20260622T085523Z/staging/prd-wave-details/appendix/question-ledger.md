<!-- FILE: docs/prd-wave-details/appendix/question-ledger.md -->
<!-- VERSION: 1.0.1 -->

# Question Ledger

## Open Questions

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
|----|------|-------|----------|--------|----------|---------------|--------------|------------------|--------|-----------|
| DQ-W09-001 | WAVE-09 | product-ac | BLOCKING | Q-ACTOR-08, Q-AC-15 | Import behavior when data already exists (merge/replace/error) | Restore transaction design depends on this | Decision: merge, replace-silently, or reject-with-error | planner/product-ac, product-verified/features/backup-and-restore.md | open | No resolution — needs product owner decision |
| DQ-W09-005 | WAVE-09 | architecture-codebase | BLOCKING | Q-EDGE-11 | Migration strategy for schema version differences (same-version-only vs migration runner) | Schema comparison logic during import validation | Decision: same-version-only for MVP or implement migration runner | planner/architecture-codebase, product-verified/edge-cases.md | open | Recommended: same-version-only for MVP |
| DQ-W09-002 | WAVE-09 | product-ac | WATCHLIST | Q-AC-16 | CSV files in backup — mandatory or optional? | Affects ZIP content generation | Confirm: exclude CSV files from backup | planner/product-ac, product-verified/functional-spec.md §20 | open | Recommended: exclude CSV (no AC requires it) |
| DQ-W09-003 | WAVE-09 | data-integration-ops | WATCHLIST | — | Import ZIP upload size limit | MaxBytesReader config parameter | Max size in MB for import uploads | planner/data-integration-ops | open | Recommended: 500MB |
| DQ-W09-004 | WAVE-09 | security-privacy-compliance | WATCHLIST | — | Should backup/import operations be logged differently than AI export? | Privacy compliance (AC-117-120) | Confirm: log event metadata only | planner/security-compliance | open | Recommended: log event metadata only, not data content |

## Answered Questions
None.

## Follow-Up Questions
None.

## Resolved Questions

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
|----|------|-------|----------|--------|----------|---------------|--------------|------------------|--------|-----------|
| DQ-W08-001 | WAVE-08 | data-api-integration-ops | needs-owner-decision | None | Should planned_actions be a simple TEXT field (MVP) or a structured child table? | PRD says "planned actions storage" — structured enables queryable action tracking; simple TEXT matches MVP constraints | Confirm: simple TEXT for MVP, structured in post-MVP | planner-product-ac-attempt-1, planner-architecture-codebase-attempt-1 | resolved | Simple TEXT for MVP (user-approved 2026-06-21) |
| DQ-W08-002 | WAVE-08 | sequencing-fit | needs-owner-decision | None | Should WAVE-08 expose ListAllByUserID for WAVE-09 backup consumption? | WAVE-07 context states "WAVE-08 must provide service layer for WAVE-09 to include AiReview data in backups" | Confirm: yes, expose ListAllByUserID | planner-sequencing-fit-attempt-1, wave-07.md | resolved | Yes, expose ListAllByUserID (user-approved 2026-06-21) |

## Deferred Questions
None.