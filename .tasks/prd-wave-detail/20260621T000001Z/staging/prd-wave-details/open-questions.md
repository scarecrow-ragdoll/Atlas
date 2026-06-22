<!-- FILE: docs/prd-wave-details/open-questions.md -->
<!-- VERSION: 1.0.0 -->

# Open Questions

## Wave-Blocking
None.

## Needs Owner Decision
None.

## Deferred
None.

## Watchlist
None.

## Resolved This Run

| ID | Wave | Scope | Severity | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
|---|---|---|---|---|---|---|---|---|---|
| DQ-W08-001 | WAVE-08 | data-api-integration-ops | needs-owner-decision | Should planned_actions be a simple TEXT field (MVP) or a structured child table? | PRD says "planned actions storage" — structured enables queryable action tracking for future features; simple TEXT matches MVP constraints. All 6 planners recommend TEXT for MVP. | Confirm: simple TEXT for MVP, structured in post-MVP | planner-product-ac-attempt-1, planner-architecture-codebase-attempt-1 | resolved | Simple TEXT for MVP (user-approved 2026-06-21) |
| DQ-W08-002 | WAVE-08 | sequencing-fit | needs-owner-decision | Should WAVE-08 expose ListAllByUserID for WAVE-09 backup consumption? | WAVE-07 context states "WAVE-08 must provide service layer for WAVE-09 to include AiReview data in backups." Without this interface, WAVE-09 cannot include AiReview data. | Confirm: yes, expose ListAllByUserID | planner-sequencing-fit-attempt-1, wave-07.md "Future Wave Compatibility" | resolved | Yes, expose ListAllByUserID (user-approved 2026-06-21) |