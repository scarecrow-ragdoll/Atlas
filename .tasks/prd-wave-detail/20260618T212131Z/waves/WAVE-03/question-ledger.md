# WAVE-03 Question Ledger

## Open Questions

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Q-WORKOUT-001 | WAVE-03 | operations | needs-owner-decision | None | Concurrent edit handling? Two browser tabs could lead to last-write-wins data loss. | Data integrity for workout diary entries. | Confirm strategy: optimistic locking (version field), last-write-wins, or deferred (acknowledge risk for MVP). | docs/prd-waves/open-questions.md | open | needs-owner-decision |
| DQ-W03-001 | WAVE-03 | operations | deferred | EDGE-016 | Should DailyLog updates use optimistic concurrency control (version field)? | Prevents last-write-wins data loss across browser tabs. | Deferred to post-MVP or owner decision. MVP uses last-write-wins. | planner-data-integration-ops-attempt-1.md | deferred | Deferred to post-MVP. MVP uses last-write-wins. |

## Resolved This Run

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W03-002 | WAVE-03 | data-ops | wave-blocking | WAVE-01 | What migration number should WAVE-03 start from? | Sequential migration order. | Start at 00082 after WAVE-02's 00080/00081 and any WAVE-01 migrations. | planner-architecture-codebase-attempt-1.md | resolved | Start at 00082. WAVE-02 uses 00080/00081. |
| DQ-W03-003 | WAVE-03 | data-ops | wave-blocking | WAVE-02 | Does allExercises query contract include exercise.workingWeight for snapshot? | Working weight snapshot requires reading current value from exercise. | Yes, allExercises returns workingWeight field. | planner-sequencing-fit-attempt-1.md | resolved | WAVE-02 allExercises returns workingWeight. WAVE-03 reads it on exercise add to workout. |
| DQ-W03-004 | WAVE-03 | product | needs-owner-decision | EDGE-004 | Backdating workout — is there any lower bound on how far back a date can be? | UI constraint vs database constraint. | No lower bound. Any valid date allowed. | planner-product-ac-attempt-1.md | resolved | No lower bound. Any date accepted. Year 0000 to 9999 validation deferred to DB DATE type. |
| DQ-W03-005 | WAVE-03 | product | needs-owner-decision | EDGE-005 | Empty workout day — can a DailyLog exist with no exercises, sets, or cardio? | Upsert behavior on first save. | DailyLog created on first save with at least one exercise or cardio entry. Empty day record is allowed but not required. | planner-product-ac-attempt-1.md | resolved | DailyLog created on first save that includes at least one WorkoutExercise or CardioEntry. Empty record (date only) not created. |
| DQ-W03-006 | WAVE-03 | data-ops | needs-owner-decision | EDGE-004 | Cascade delete behavior when a DailyLog is deleted? | Referential integrity with WorkoutExercise, WorkoutSet, CardioEntry. | Delete cascades: DailyLog -> WorkoutExercise -> WorkoutSet, and DailyLog -> CardioEntry. | planner-data-integration-ops-attempt-1.md | resolved | CASCADE delete from DailyLog to WorkoutExercise and CardioEntry. CASCADE delete from WorkoutExercise to WorkoutSet. |
| DQ-W03-007 | WAVE-03 | product | wave-blocking | AC-041 | Working weight snapshot: captured when exercise is ADDED to workout or when workout is SAVED? | Timing matters for snapshot accuracy. | Captured when exercise is added (POST / mutation), read current value from Exercise at that moment. | planner-product-ac-attempt-1.md | resolved | Snapshot captured at exercise-add time, not save time. Read Exercise.workingWeight at ADD moment. |
