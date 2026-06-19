# Consistency Reviewer — Question Ledger

**Run ID:** 20260618T185935Z
**Scope:** consistency-reviewer

| ID | Scope | Severity | Question | Why It Matters | Source Or Report | Status |
| --- | --- | --- | --- | --- | --- | --- |
| Q-CONS-001 | consistency-reviewer | non-blocking | Should the 4 blocking questions (Q-SCOPE-001, Q-SCOPE-002, Q-SCOPE-004, Q-SCOPE-005) be resolved before synthesis proceeds, or can synthesis proceed with them as acknowledged open items? | Affects handoff gate definition and development planning. product-scope-reviewer says "NOT ready for handoff" but all other scopes approved. | worker-attempt-1.md §7 | open |
| Q-CONS-002 | consistency-reviewer | non-blocking | What timezone handling strategy should be used for all date-based features (workout day, dashboard week, check-in date, nutrition week start, AI export period)? | Every date feature in the PRD assumes implicit single timezone. Self-hosted users may be in any timezone. Server clock vs user clock vs browser clock affects data consistency. | Consistency pass — no Phase 1 scope raised this as a standalone question (noted as minor gap by edge-case-risk reviewer, suggested by AC reviewer as Q-AC-19) | open |