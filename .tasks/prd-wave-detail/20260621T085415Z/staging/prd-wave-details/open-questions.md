# Open Questions

## Wave-Blocking
None resolved this run.

## Needs Owner Decision
None — all resolved this run.

## Deferred
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W06-005 | WAVE-06 | data-ops | deferred | — | Max date range to prevent expensive queries | Nutrition weekly averages iterate per week. 52-week max chosen. | — | user decision 2026-06-21 | open | 52-week max enforced via server constant. |
| DQ-W06-007 | WAVE-06 | security | deferred | — | Log detail level for chart queries | Balance between observability and privacy | — | reviewer-security-privacy-compliance-attempt-1 | open | Body weight/measurement values NOT logged. Dates and counts logged. |

## Watchlist
None

## Resolved This Run
| ID | Wave | Scope | Severity | Question | Resolution |
| --- | --- | --- | --- | --- | --- |
| DQ-W06-001 | WAVE-06 | product-ac | resolved | Best set definition | Highest e1RM per session |
| DQ-W06-002 | WAVE-06 | architecture | resolved | Exercise chart with no WAVE-03 | Stubs returning empty |
| DQ-W06-003 | WAVE-06 | product-ac | resolved | Working weight source | Per-session snapshot |
| DQ-W06-004 | WAVE-06 | data-ops | resolved | Default chart period | 4 weeks |
| DQ-W06-006 | WAVE-06 | data-ops | resolved | Exercise chart stubs or omitted | Stubs returning empty |