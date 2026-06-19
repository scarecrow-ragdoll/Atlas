# Open Questions
## Wave-Blocking
No wave-blocking questions open.
## Needs Owner Decision
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W02-002 | WAVE-02 | product | needs-owner-decision | AC-043 | Are exercise names unique per user or can duplicates exist? | Affects validation logic and UI behavior | Duplicates allowed (no constraint) — consistent with EDGE-002. | planner-product-ac-attempt-2.md | answered | Tentative: duplicates allowed per EDGE-002. Awaiting user confirmation. |
## Deferred
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W02-003 | WAVE-02 | data-ops | wave-blocking | WAVE-01 | What exact file storage path pattern does WAVE-01 MediaConfig provide for exercise media? | Drives migration and handler design | Use WAVE-01 BasePath/<exercise_id>/<uuid>.<ext>. Confirm after WAVE-01 implementation. | planner-data-integration-ops-attempt-2.md | deferred | WAVE-01 coordination item. |
| DQ-W02-006 | WAVE-02 | security | deferred | EDGE-014 | Should exercise media URLs be time-limited (signed URLs)? | Signed URLs add complexity for single-user MVP | Session-gated access sufficient for MVP self-hosted deployment. | planner-security-compliance-attempt-2.md | deferred | Deferred post-MVP. |
## Watchlist
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W02-008 | WAVE-02 | sequencing | watchlist | WAVE-03 | Does allExercises need filtering beyond isActive for WAVE-03 exercise selector? | Might be needed if library grows large | Deferred — current scope is unfiltered active list ordered by name. | planner-sequencing-fit-attempt-2.md | deferred | Watchlist. Not needed for MVP. |
## Resolved This Run
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W02-001 | WAVE-02 | data-ops | wave-blocking | EDGE-020 | Should deleting ExerciseMedia delete the physical media file from disk? | Orphaned files accumulate | Yes, per TDEC-005. | planner-data-integration-ops-attempt-2.md | resolved | Physical file deleted. On failure, log error and return 204. |
| DQ-W02-005 | WAVE-02 | security | needs-owner-decision | TDEC-008 | Server-side MIME detection or trust Content-Type header? | Content-Type spoofable | Use http.DetectContentType() server-side. | planner-security-compliance-attempt-2.md | resolved | http.DetectContentType() as primary check. |
| DQ-W02-007 | WAVE-02 | testing | needs-owner-decision | WAVE-01 | Mocked PIN auth or full middleware chain? | Test complexity vs realism | Prefer full middleware chain integration tests. | planner-testing-exit-attempt-2.md | resolved | Full middleware chain with WAVE-01 PIN test helpers. |