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
| DQ-W05-007 | WAVE-05 | security | deferred | — | Should soft-deleted products be recoverable via API or admin DB only? | Data recovery for accidental deletion | Admin-only recovery via direct DB. No restore API in MVP. | planner-security-compliance-attempt-1.md | open | Admin-only DB for MVP. |
| DQ-W05-009 | WAVE-05 | sequencing | deferred | WAVE-04 | What migration number should WAVE-05 use? | Avoiding migration number collision with WAVE-04 | Check current migration state at implementation time. Start at 00081 or next available. | planner-sequencing-fit-attempt-1.md | open | Coordinate with WAVE-04 implementation. Use next available number. |
## Watchlist
None.
## Resolved This Run
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W04-002 | WAVE-04 | operations | deferred | EDGE-006 | Is 2-4 photos per check-in a hard requirement or recommendation? | Affects validation logic | Soft recommendation (warn, don't block) with hard limit of 10 | planner-product-ac-attempt-1.md | resolved | Soft guidance. DDEC-W04-002. |
| DQ-W04-003 | WAVE-04 | data-ops | resolved | EDGE-007 | Should body measurement value 0 or negative be rejected? | Data integrity | Reject 0 and negative values. Validated in service layer. | planner-product-ac-attempt-1.md | resolved | AC-W04-028: value must be > 0. |
| DQ-W05-001 | WAVE-05 | data | resolved | EDGE-019 | Should NutritionProduct use soft-delete (isActive flag) or hard-delete with FK block? | Referential integrity for historical template/override data | Soft-delete with isActive flag. Products remain in DB but excluded from default queries. | planner-product-ac-attempt-2.md | resolved | Soft-delete (isActive flag). DDEC-W05-003. |
| DQ-W05-002 | WAVE-05 | product-ac | resolved | RULE-020 | What is exact "single template at a time" semantic — per-week or per-user? | Drives upsert behavior | Per-week: template for week X replaces previous template for that week. | planner-product-ac-attempt-1.md | resolved | Per-week upsert. DDEC-W05-001. |
| DQ-W05-003 | WAVE-05 | product-ac | resolved | — | mealLabel — free text or enum? | Flexibility vs validation | Free-text string. | planner-product-ac-attempt-1.md | resolved | Free-text. DDEC-W05-004. |
| DQ-W05-004 | WAVE-05 | architecture | resolved | — | Macro calculation server-side or client-side? | Consistency | Server-side in Go service. | planner-architecture-codebase-attempt-1.md | resolved | Server-side. DDEC-W05-002. |
| DQ-W05-005 | WAVE-05 | data-ops | resolved | — | Macro query: separate or inline? | API surface design | Separate query. | planner-data-integration-ops-attempt-1.md | resolved | nutritionMacros query. |
| DQ-W05-006 | WAVE-05 | data-ops | resolved | — | NutritionProduct deletion: soft or hard? | Same as DQ-W05-001 | Soft-delete. | planner-data-integration-ops-attempt-1.md | resolved | Soft-delete. Merged to DQ-W05-001. |
| DQ-W05-008 | WAVE-05 | testing | resolved | — | Macro tests: unit or integration? | Test scope | Both: unit (calculation) + integration (round-trip). | planner-testing-exit-attempt-1.md | resolved | Both types used. |