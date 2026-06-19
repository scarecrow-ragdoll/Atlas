# Open Questions

## Wave-Blocking
No wave-blocking questions open for WAVE-03.

## Needs Owner Decision
No unresolved needs-owner-decision questions for WAVE-03.

## Deferred
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W03-001 | WAVE-03 | operations | deferred | EDGE-016 | Should DailyLog updates use optimistic concurrency control (version field)? | Prevents last-write-wins data loss across browser tabs. | Deferred to post-MVP or owner decision. MVP uses last-write-wins. | planner-data-integration-ops-attempt-1.md | deferred | Deferred to post-MVP. MVP uses last-write-wins. |

## Watchlist
None.

## Resolved This Run
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Q-WORKOUT-001 | WAVE-03 | operations | needs-owner-decision | None | Concurrent edit handling? Two browser tabs could lead to last-write-wins data loss. | Data integrity for workout diary entries. | Confirm strategy: optimistic locking (version field), last-write-wins, or deferred (acknowledge risk for MVP). | docs/prd-waves/open-questions.md | resolved | last-write-wins for MVP. No optimistic locking. Tracked as DQ-W03-001 (deferred post-MVP). |
| DQ-W03-002 | WAVE-03 | data-ops | wave-blocking | WAVE-01 | What migration number should WAVE-03 start from? | Sequential migration order. | Start at 00082 after WAVE-02's 00080/00081. | planner-architecture-codebase-attempt-1.md | resolved | Start at 00082. |
| DQ-W03-003 | WAVE-03 | data-ops | wave-blocking | WAVE-02 | Does allExercises query contract include exercise.workingWeight for snapshot? | Working weight snapshot requires reading current value. | Yes, allExercises returns workingWeight field. | planner-sequencing-fit-attempt-1.md | resolved | WAVE-02 allExercises returns workingWeight. |
| DQ-W03-004 | WAVE-03 | product | needs-owner-decision | EDGE-004 | Backdating lower bound? | UI vs DB constraint. | No lower bound. Any valid date accepted. | planner-product-ac-attempt-1.md | resolved | No lower bound. |
| DQ-W03-005 | WAVE-03 | product | needs-owner-decision | EDGE-005 | Empty workout day creation? | Upsert behavior. | Created on first content save. | planner-product-ac-attempt-1.md | resolved | DailyLog created on first content save. |
| DQ-W03-006 | WAVE-03 | data-ops | needs-owner-decision | EDGE-004 | Cascade delete behavior? | Referential integrity. | CASCADE all FK relationships. | planner-data-integration-ops-attempt-1.md | resolved | CASCADE delete throughout. |
| DQ-W03-007 | WAVE-03 | product | wave-blocking | AC-041 | Working weight snapshot timing? | When is snapshot captured? | Captured at exercise-add time. | planner-product-ac-attempt-1.md | resolved | Snapshot captured at exercise-add time. |
