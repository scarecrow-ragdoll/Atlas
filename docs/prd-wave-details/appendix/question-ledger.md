# Question Ledger
## Open Questions
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W04-001 | WAVE-04 | operations | needs-owner-decision | WAVE-03 | Should WAVE-04 include daily_log table creation or require WAVE-03 to be deployed first? | Affects deployment ordering and migration strategy | WAVE-04 should create daily_log table migration as prerequisite | planner-sequencing-fit-attempt-1.md | open | DailyLog auto-creation logic required. WAVE-04 should include a daily_log migration if WAVE-03 not yet deployed. |
| DQ-W05-007 | WAVE-05 | security | deferred | — | Should soft-deleted products be recoverable via API or admin DB only? | Data recovery for accidental deletion | Admin-only DB recovery. | planner-security-compliance-attempt-1.md | open | Admin-only for MVP. |
| DQ-W05-009 | WAVE-05 | sequencing | deferred | WAVE-04 | What migration number should WAVE-05 use? | Avoiding collision with WAVE-04 | Check current state at implementation time. | planner-sequencing-fit-attempt-1.md | open | Coordinate with WAVE-04. Use next available. |
## Answered Questions
None.
## Follow-Up Questions
None.
## Resolved Questions
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W04-002 | WAVE-04 | operations | deferred | EDGE-006 | Is 2-4 photos per check-in a hard requirement or recommendation? | Affects validation logic | Soft recommendation (warn, don't block) with hard limit of 10 | planner-product-ac-attempt-1.md | resolved | Soft guidance. DDEC-W04-002. |
| DQ-W04-003 | WAVE-04 | data-ops | resolved | EDGE-007 | Should body measurement value 0 or negative be rejected? | Data integrity | Reject 0 and negative values. Validated in service layer. | planner-product-ac-attempt-1.md | resolved | AC-W04-028: value must be > 0. |
| DQ-W05-001 | WAVE-05 | data | resolved | EDGE-019 | Should NutritionProduct use soft-delete (isActive flag) or hard-delete? | Referential integrity for historical data | Soft-delete with isActive flag. | planner-product-ac-attempt-2.md | resolved | Soft-delete. DDEC-W05-003. |
| DQ-W05-002 | WAVE-05 | product-ac | resolved | RULE-020 | Single template per-week or per-user? | Drives upsert behavior | Per-week upsert. | planner-product-ac-attempt-1.md | resolved | Per-week upsert. DDEC-W05-001. |
| DQ-W05-003 | WAVE-05 | product-ac | resolved | — | mealLabel: free text or enum? | Flexibility | Free-text string. | planner-product-ac-attempt-1.md | resolved | Free-text. DDEC-W05-004. |
| DQ-W05-004 | WAVE-05 | architecture | resolved | — | Macro calc server-side or client-side? | Consistency | Server-side. | planner-architecture-codebase-attempt-1.md | resolved | Server-side. DDEC-W05-002. |
| DQ-W05-005 | WAVE-05 | data-ops | resolved | — | Macro query separate or inline? | API surface | Separate query. | planner-data-integration-ops-attempt-1.md | resolved | nutritionMacros query. |
| DQ-W05-008 | WAVE-05 | testing | resolved | — | Macro tests: unit or integration? | Test scope | Both. | planner-testing-exit-attempt-1.md | resolved | Both types used. |
## Deferred Questions
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W04-004 | WAVE-04 | security | deferred | TDEC-008 | Should progress photo URLs be time-limited (signed URLs)? | Session-gated access sufficient for MVP self-hosted | Signed URLs add complexity for single-user MVP | planner-security-compliance-attempt-1.md | deferred | Deferred post-MVP. |
| DQ-W04-005 | WAVE-04 | data-ops | deferred | WAVE-01 | What exact file storage path pattern does WAVE-01 MediaConfig provide for progress photos? | Drives migration and handler design | Use WAVE-01 BasePath/progress-photos/<checkin_id>/<uuid>.<ext> | planner-data-integration-ops-attempt-1.md | deferred | Confirmed after WAVE-01 implementation. WAVE-04 assumes composable BasePath. |