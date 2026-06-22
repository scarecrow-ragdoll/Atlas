<!-- FILE: .tasks/prd-wave-detail/20260621T000001Z/waves/WAVE-08/question-ledger.md -->
<!-- VERSION: 1.0.0 -->

# WAVE-08 Question Ledger

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
|---|---|---|---|---|---|---|---|---|---|---|
| DQ-W08-001 | WAVE-08 | data-api-integration-ops | needs-owner-decision | None | Should planned_actions be a simple TEXT field or a separate child table? | PRD says "planned actions storage" — structured approach enables queryable action tracking; simple text matches MVP constraints | Simple TEXT for MVP (matching domain model and userNotes pattern) | planner-product-ac-attempt-1 Q-W08-PAC-001, planner-architecture-codebase-attempt-1 Q-W08-ARC-001 | open | Pending: recommended "TEXT for MVP, structured in post-MVP" per all planners |
| DQ-W08-002 | WAVE-08 | sequencing-fit | needs-owner-decision | None | Does WAVE-09 (Backup) need AiReview data included in data.json? If yes, AiReviewService must expose ListAllByUserID. | WAVE-07 context says "WAVE-08 must provide service layer for WAVE-09 to include AiReview data in backups" | Confirm yes: WAVE-09 consumes AiReview via AiReviewService.ListAllByUserID | planner-sequencing-fit-attempt-1 Q-W08-SEQ-001, wave-07.md "Future Wave Compatibility" | open | Pending: recommended "Yes, expose ListAllByUserID for WAVE-09" per sequencing-fit planner |