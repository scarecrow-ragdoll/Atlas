# Product Scope Reviewer Question Ledger

Run ID: 20260618T185935Z
Scope: product-scope-reviewer

| ID | Scope | Severity | Question | Why It Matters | Source | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Q-SCOPE-001 | product-scope | blocking | What are the quantitative success metrics for the MVP? | Without success metrics, there is no way to validate the product is achieving its goals | prd.md (entire document — no success metrics found) | open | TBD |
| Q-SCOPE-002 | product-scope | blocking | Should the MVP architecture build for future multi-user or remain strictly single-user? | Affects data model, auth design, API structure, and database schema decisions | prd.md sections 4, 28 | open | TBD |
| Q-SCOPE-003 | product-scope | non-blocking | What is the target user's expected technical proficiency? | Self-hosted requires Docker/CLI knowledge — affects documentation and deployment UX requirements | prd.md section 1 (self-hosted requirement) | open | TBD |
| Q-SCOPE-004 | product-scope | blocking | What are the specific performance targets (page load, export time, chart render time)? | Section 24.3 uses vague language ("быстрым") with no measurable targets | prd.md section 24.3 | open | TBD |
| Q-SCOPE-005 | product-scope | blocking | Should cardio be a separate entity or always part of a workout day? | The data model shows cardio as separate with optional workoutDayId, but section 10.3 includes cardio in the workout day — ambiguous | prd.md sections 10.3, 25.8 | open | TBD |
| Q-SCOPE-006 | product-scope | non-blocking | What AI models/platforms must the export format support? Only ChatGPT, or Claude, Gemini, local LLMs? | The core value proposition depends on AI compatibility; format design affects all models differently | prd.md sections 1, 17 | open | TBD |
| Q-SCOPE-007 | product-scope | non-blocking | Is there a maximum photo/media storage limit? What happens when storage volume runs out? | Photos over years of training accumulate; no storage management policy exists | prd.md sections 14, 24.2 | open | TBD |
| Q-SCOPE-008 | product-scope | non-blocking | What data portability standard is required for "data belongs to user" commitment? | Export/backup formats need an interoperability guarantee beyond the current ZIP structure | prd.md section 6.3 | open | TBD |