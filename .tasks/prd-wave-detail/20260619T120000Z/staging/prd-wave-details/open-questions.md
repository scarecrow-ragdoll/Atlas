# Open Questions
## Wave-Blocking
None.
## Needs Owner Decision
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W04-001 | WAVE-04 | operations | needs-owner-decision | WAVE-03 | Should WAVE-04 include daily_log table creation or require WAVE-03 to be deployed first? | Affects deployment ordering and migration strategy | WAVE-04 should create daily_log table migration as prerequisite | planner-sequencing-fit-attempt-1.md | open | DailyLog auto-creation logic required. WAVE-04 should include a daily_log migration if WAVE-03 not yet deployed. |
## Deferred
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W04-004 | WAVE-04 | security | deferred | TDEC-008 | Should progress photo URLs be time-limited (signed URLs)? | Session-gated access sufficient for MVP self-hosted | Signed URLs add complexity for single-user MVP | planner-security-compliance-attempt-1.md | deferred | Deferred post-MVP. |
## Watchlist
None.
## Resolved This Run
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W04-002 | WAVE-04 | operations | deferred | EDGE-006 | Is 2-4 photos per check-in a hard requirement or recommendation? | Affects validation logic | Soft recommendation (warn, don't block) with hard limit of 10 | planner-product-ac-attempt-1.md | resolved | Soft guidance. DDEC-W04-002. |
| DQ-W04-003 | WAVE-04 | data-ops | resolved | EDGE-007 | Should body measurement value 0 or negative be rejected? | Data integrity | Reject 0 and negative values. Validated in service layer. | planner-product-ac-attempt-1.md | resolved | AC-W04-028: value must be > 0. |