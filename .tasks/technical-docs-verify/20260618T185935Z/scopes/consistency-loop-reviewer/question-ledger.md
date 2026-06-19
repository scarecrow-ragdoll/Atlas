# Consistency-Loop-Reviewer Question Ledger

These questions are cross-scope findings raised by the consistency loop — not new technical questions, but structural issues in how the 8 scope reports relate.

## Cross-Scope Contradictions

| ID | Scopes Involved | Contradiction | Severity | Resolution |
|---|---|---|---|---|
| CL-CON-001 | api-contracts, auth-security-compliance | userId scoping: API says server-enforced everywhere; auth says no ownership checking needed in MVP | blocking | Architecture owner must reconcile: middleware-level vs endpoint-level enforcement |
| CL-CON-002 | auth-security-compliance, integrations-events, client-state-ux | Session TTL split across 3 scopes with no unified decision | blocking | Create single session lifecycle decision record covering TTL, renewal, cookies, long-operation extension, loss recovery |
| CL-CON-003 | architecture-boundaries, auth-security-compliance, client-state-ux | PIN guard: arch missing middleware analysis, auth has rule contradiction (RULE-022 vs RULE-024), UX has lockout design | blocking | Resolve RULE-022 vs RULE-024 before any PIN implementation |
| CL-CON-004 | api-contracts, integrations-events, operations-observability, testing-delivery | Backup import flow owned by 4 scopes with no end-to-end contract | blocking | Define top-level backup/restore flow contract unifying endpoint sequence, sync/async, rollback, scheduling, schema versioning |
| CL-CON-005 | data-contracts, auth-security-compliance, integrations-events, operations-observability, client-state-ux | Redis role: 5 scopes assume different roles (session, progress, cache, dependency) with no unified contract | blocking | Define Redis MVP role contract settling session, progress, and caching responsibilities |
| CL-CON-006 | architecture-boundaries, operations-observability, client-state-ux, testing-delivery, api-contracts | DEC-008 performance targets span 5 scopes but no unified performance budget | blocking | Create performance budget mapping each p95 target to measurement point, resource plan, loading state, and test gate |

## Source Delta Coverage Gaps

| ID | Gap | Severity | Resolution |
|---|---|---|---|
| CL-GAP-001 | DEC-007 (userId FK) not reviewed by client-state-ux and integrations-events | high | Request those scopes to review userId FK effects on UI state and export data |
| CL-GAP-002 | DEC-008 (p95 performance) not explicitly reviewed by api-contracts | high | Request api-contracts to review API-level performance implications |
| CL-GAP-003 | DEC-009 (DailyLog rename) not reviewed by integrations-events and client-state-ux | high | Request those scopes to review effects on export CSV schema and workout diary UI |

## Severity Drift

| ID | Question | Scope A | Severity A | Scope B | Severity B | Recommended |
|---|---|---|---|---|---|---|
| CL-DRIFT-001 | Session TTL | auth-security | dev-blocking | integrations-events | needs-owner-decision | Promote TQ-INT-005 to dev-blocking |
| CL-DRIFT-002 | Data retention | data-contracts | watchlist | auth-security | needs-owner-decision | Align both to needs-owner-decision |
| CL-DRIFT-003 | Redis usage | data-contracts | deferred | integrations-events, auth-security | assumed active | Resolve deferral vs active assumption |

## Question Status Summary

| Status | Count | Notes |
|---|---|---|
| blocking | 9 | 6 contradictions + 3 source delta gaps |
| high | 3 | Source delta coverage gaps |
| resolved | 0 | All still open |
| **Total** | **12** |