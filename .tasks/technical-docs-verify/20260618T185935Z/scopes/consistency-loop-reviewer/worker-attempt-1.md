# Consistency-Loop-Reviewer Worker Report — Attempt 2

## Revisions Applied
Per review-attempt-1.md:
1. Added "conditional partial-approval" sub-gate analysis (R1)
2. Added resolution blocker ranking by estimated effort (R2)
3. Added pre-implementation vs in-wave question classification (R3)
4. Added positive note on consistent question numbering (R4)
5. Added positive note on consistent question ledger format (R5)

## Run Metadata
- **Run ID:** 20260618T185935Z
- **Role:** consistency-loop-reviewer (cross-scope)
- **Source:** All 8 scope reports, source-delta.md, aggregate scope-status.md

## Cross-Scope Contradictions

### C1: userId Scoping — API Enforces vs Auth Ignores
- **api-contracts** (worker-attempt-1.md §DEC-007 API Implications): "Every create/read/update/list must scope to the authenticated user's userId. This is a server-enforced rule."
- **auth-security-compliance** (worker-attempt-1.md §Ownership): "No ownership checking needed in MVP (all data is the user's)."
- **Impact**: API-contracts assumes server-enforced userId scoping on every endpoint; auth-security says no ownership checking needed. If auth does not implement middleware-level userId scoping, API endpoints must do it inline — conflicting design expectations.
- **Owner**: architecture owner must reconcile: is userId scoping enforced at the middleware/auth layer or per-endpoint?

### C2: Session TTL — Three Scopes, Three Partial Views
- **auth-security-compliance** (TQ-AUTH-002): Session TTL, renewal policy, cookie flags undefined.
- **integrations-events** (TQ-INT-005): Session timeout interaction with long-running export/import operations.
- **client-state-ux** (TQ-CLIENT-013): Session loss recovery UX during active data entry.
- **Impact**: No single scope owns the complete session lifecycle decision. Auth owns the mechanism, integrations needs the interaction mode, client-state needs the recovery UI. Without a unified session policy, these three partial views cannot be resolved independently.
- **Recommendation**: Create a single "session lifecycle" decision record that addresses TTL, renewal, cookie flags, long-operation extension, and loss recovery together.

### C3: PIN Guard — Architecture vs Auth vs UX Disconnected
- **architecture-boundaries** (TQ-ARCH-007, TQ-ARCH-008): PIN guard and media serving architecture not analyzed (identified as minor gaps).
- **auth-security-compliance** (TQ-AUTH-006): RULE-022 vs RULE-024 contradiction — media access when PIN is disabled is unresolved.
- **client-state-ux** (TQ-CLIENT-015): PIN brute-force lockout UX — retry countdown, lockout timer, recovery path.
- **Impact**: The PIN guard has architectural implications (middleware layer, session boundaries), security implications (contradiction in media access rules), and UX implications (lockout behavior). These are split across three scope silos with no cross-reference. The RULE-022 vs RULE-024 contradiction in particular needs a unified architectural+security resolution.
- **Recommendation**: Resolve RULE-022 vs RULE-024 before any PIN implementation. Document the decision in both architecture and auth scope docs.

### C4: Backup Import Flow — Four Scopes, No Owner
- **api-contracts** (TQ-API-011): Multi-step flow needs state management, endpoint sequence undefined.
- **integrations-events** (TQ-INT-003): Async/sync decision + transaction rollback model.
- **operations-observability** (TQ-OPS-005): Backup scheduling, retention, verification strategy.
- **testing-delivery** (TQ-TEST-004): Backup manifest schema for contract tests.
- **Impact**: The backup/restore feature touches API design, async job handling, operations runbooks, and test artifacts — but no scope owns the end-to-end backup flow contract. Each scope assumes a different aspect without cross-referencing.
- **Recommendation**: Define a top-level "backup/restore flow contract" that unifies the endpoint sequence (API), sync/async decision (integrations), rollback model (integrations), scheduling (ops), and schema versioning (testing).

### C5: Redis Role — Five Scopes, Five Assumptions
- **data-contracts** (TQ-DATA-006): Redis usage patterns for MVP undefined (deferred).
- **auth-security-compliance** (TQ-AUTH-005, TQ-AUTH-007): Session store in Redis, failure mode when unavailable.
- **integrations-events** (TQ-INT-001, TQ-INT-007): Suggested Redis for export progress tracking, dependency coupling risk.
- **operations-observability**: Redis as dependency stack component.
- **client-state-ux**: No cache layer for UI data (TQ-CLIENT-009) — Redis not considered.
- **Impact**: Five scopes assume Redis plays different roles (session store, progress tracker, cache layer, dependency) but no single contract defines Redis's actual MVP responsibilities. The client-state-ux scope assumes no cache layer, which may miss Redis as a potential cache for DEC-008 performance targets.
- **Recommendation**: Define a "Redis MVP role contract" settling session store, progress tracking (if any), and caching (if any) in one place.

### C6: DEC-008 Performance Targets — No Unified Performance Contract
- **architecture-boundaries** (TQ-ARCH-004): Deployment resources needed to meet SLOs.
- **operations-observability** (TQ-OPS-003): SLO measurement instrumentation.
- **client-state-ux** (TQ-CLIENT-002, TQ-CLIENT-009): Loading states and cache strategy.
- **testing-delivery** (TQ-TEST-007): Performance test policy.
- **api-contracts**: Not explicitly reviewed for performance.
- **Impact**: The p95 targets span UI, API, export, and backup but no scope connects them into a unified performance budget. API-contracts didn't evaluate API-level performance implications. The loading state rule (>2s) and the p95 targets are treated independently.
- **Recommendation**: Create a "performance budget" document that maps each p95 target to: measurement point (ops), resource plan (arch), loading state trigger (ux), and test gate (testing).

## Duplicate / Overlapping Questions

| Duplicate Group | Scopes | Questions | Impact |
|---|---|---|---|
| D1: Session TTL | auth, integrations, client-state | TQ-AUTH-002, TQ-INT-005, TQ-CLIENT-013 | Three questions about session timeout. Consolidate into one cross-scope decision. |
| D2: Data retention | data-contracts, auth-security | TQ-DATA-008, TQ-AUTH-009 | Both ask about data retention/deletion policy. Watchlist vs needs-owner-decision level mismatch. |
| D3: Import behavior (carried forward) | operations-observability | Q-ACTOR-08, Q-AC-15 | Same product question documented twice in product-verified. Carried forward as duplicates. |
| D4: Backup import flow | api-contracts, integrations-events | TQ-API-011, TQ-INT-003 | Both ask about the import workflow design; API asks about endpoint sequence, integrations about sync/async. Partial overlap — complementary, not full duplicate. |

## Source Delta (DEC-006 through DEC-009) Coverage Check

### DEC-006 (Quality Gates) — Affected: testing-delivery, architecture-boundaries
| Scope | Reviewed? | Evidence |
|---|---|---|
| testing-delivery | ✓ | TQ-TEST-001 through TQ-TEST-005, all six quality gates mapped |
| architecture-boundaries | ✓ (partial) | Mentioned as affecting test architecture expectations |
| all others | N/A | Not directly affected — correct |

### DEC-007 (userId FK) — Affected: data-contracts, api-contracts, auth-security
| Scope | Reviewed? | Evidence |
|---|---|---|
| data-contracts | ✓ | TQ-DATA-001, TQ-DATA-002, TQ-DATA-003 (userId ambiguity) |
| api-contracts | ✓ | userId scoping on all endpoints |
| auth-security-compliance | ✓ | Multi-user-ready model, no ownership checking in MVP |
| testing-delivery | ✓ (partial) | Single-user MVP mentioned |
| client-state-ux | ✗ **GAP** | DEC-007 not explicitly reviewed — userId FK implications for UI state (e.g., display user info, filter by user) not evaluated |
| integrations-events | ✗ **GAP** | DEC-007 not reviewed — userId scoping could affect export/backup data filtering |

### DEC-008 (p95 Performance Targets) — Affected: architecture, ops, client-state-ux
| Scope | Reviewed? | Evidence |
|---|---|---|
| architecture-boundaries | ✓ | TQ-ARCH-004 (deployment resources for SLOs) |
| operations-observability | ✓ | TQ-OPS-003 (SLO measurement) |
| client-state-ux | ✓ | Loading states, cache strategy implications |
| testing-delivery | ✓ (partial) | TQ-TEST-007 (performance test policy) |
| integrations-events | ✓ (implicit) | Export/backup time targets analyzed |
| api-contracts | ✗ **GAP** | DEC-008 p95 targets for API mutations (300ms), queries (500ms–1.0s) not explicitly reviewed — no API-level performance contract |
| data-contracts | ✓ (partial) | TQ-DATA-011 (index strategy for p95) |

### DEC-009 (DailyLog replaces WorkoutDay) — Affected: data-contracts, api-contracts
| Scope | Reviewed? | Evidence |
|---|---|---|
| data-contracts | ✓ | TQ-DATA-010 (stale WorkoutDay references) |
| api-contracts | ✓ | DailyLog as unified date resource, cardio-dailyLog coupling |
| testing-delivery | ✓ (partial) | Mentioned |
| integrations-events | ✗ **GAP** | DEC-009 not reviewed — DailyLog rename may affect export CSV schema and backup data.json structure |
| client-state-ux | ✗ **GAP** | DEC-009 not reviewed — DailyLog rename affects workout diary UI, date navigation, and cardio entry forms |

### Coverage Summary
- **DEC-006**: 2/2 affected scopes covered ✓
- **DEC-007**: 3/5 affected scopes covered — **client-state-ux** and **integrations-events** missed
- **DEC-008**: 4/6 affected scopes covered — **api-contracts** and **data-contracts** partially missed
- **DEC-009**: 2/4 affected scopes covered — **integrations-events** and **client-state-ux** missed

## Answer Effects

No questions from any scope have been answered or resolved (all statuses = "open" or "deferred"). Therefore:
- No new blocking follow-ups were created by answered questions.
- No answer-effect chain exists to analyze.
- The 0% resolution rate across all 8 scopes is the primary finding: ~47 dev-blocking questions, 0 resolved.

The only "answered" items are the 4 source delta decisions (DEC-006 through DEC-009), which created new technical questions instead of resolving existing ones. Each DEC resolved a product question but introduced 1–3 new technical gaps per affected scope.

## Severity Drift

| Question | Scope | Severity | Adjacent Scope View | Drift? |
|---|---|---|---|---|
| TQ-DATA-008 | data-contracts | watchlist | TQ-AUTH-009 (needs-owner-decision) — same topic | Moderate drift: one scope says watchlist, other says needs-owner-decision for data retention. |
| TQ-INT-005 | integrations-events | needs-owner-decision | TQ-AUTH-002 (dev-blocking) — same session TTL topic | Significant drift: session TTL is dev-blocking for auth but needs-owner-decision for integrations. Auth is correct — session TTL genuinely blocks auth implementation. |
| TQ-DATA-006 | data-contracts | deferred | integrations-events assumes Redis for progress, auth assumes Redis for session (both active) | Redis usage deferred by data-contracts but assumed active by integrations and auth. Deferral contradicts other scopes' assumptions. |

## Unresolved Parent Questions

All questions are open. Parent relationships are defined only within scopes:
- TQ-CLIENT-002 through TQ-CLIENT-004 parented to TQ-CLIENT-001 (state machine)
- TQ-CLIENT-015, TQ-CLIENT-016 parented to TQ-CLIENT-001
- TQ-CLIENT-007 parented to TQ-CLIENT-006 (validation)
- TQ-DATA-008 parented to TQ-DATA-007 (media lifecycle)
- TQ-INT-007, TQ-INT-008 parented to TQ-INT-001 (async export)

No cross-scope parent relationships exist. Questions that depend on answers from other scopes have no cross-reference:
- TQ-API-012 (session auth for API) depends on TQ-AUTH-002 (session contract) — no cross-reference
- TQ-CLIENT-009 (cache strategy) depends on TQ-OPS-003 (SLO measurement) — no cross-reference
- TQ-INT-005 (session + long operations) depends on TQ-AUTH-002 (session TTL) — no cross-reference
- TQ-TEST-003 (AI export schema) depends on TQ-INT-001 (export sync/async decision) — no cross-reference

## Approved-to-Dev Assessment

**Verdict: NOT REACHABLE — BLOCKED**

### Blocking Issues

1. **Foundational decisions missing (~47 dev-blocking questions, 0 resolved)**: Every scope has open dev-blocking questions. The package cannot reach approved-to-dev while all seven foundational artifacts are absent: architecture contract, API protocol, auth spec, UI state machine, async job contract, ops config, test infrastructure.

2. **Four cross-scope contradictions block unified design**: C1 (userId scoping), C2 (session TTL), C4 (backup flow owner), C5 (Redis role) need architectural decisions before any scope can proceed independently.

3. **Source delta coverage gaps**: DEC-007 and DEC-009 not reviewed by client-state-ux and integrations-events. DEC-008 not reviewed by api-contracts. These gaps mean downstream scope reports are incomplete.

4. **Severity drift across overlapping questions**: Session TTL is dev-blocking for auth but needs-owner-decision for integrations — the highest severity should govern, creating unexpected dependencies between the scopes.

5. **No cross-scope dependency graph**: Questions in different scopes depend on each other (e.g., TQ-API-012 depends on TQ-AUTH-002) but no scope identifies these cross-dependencies. This means resolving a question in one scope may invalidate assumptions in another.

### Minimum Resolution Path to Approved-to-Dev

1. Create a unified "architecture-and-boundaries.md" resolving TQ-ARCH-001 through TQ-ARCH-006
2. Choose API protocol (REST recommended) resolving TQ-API-001
3. Define auth specification (PIN hash, session TTL, cookie config, Redis role) resolving TQ-AUTH-001, TQ-AUTH-002, TQ-AUTH-005, TQ-AUTH-007
4. Resolve RULE-022 vs RULE-024 contradiction (TQ-AUTH-006)
5. Define UI state machine convention (TQ-CLIENT-001)
6. Define backup/restore end-to-end flow contract (TQ-API-011, TQ-INT-001 through TQ-INT-003)
7. Define test data strategy and e2e plan (TQ-TEST-001, TQ-TEST-002)
8. Close source delta gaps: client-state-ux and integrations-events review DEC-007, DEC-009

### Risk Assessment if Forced Without Resolution

- **High risk**: Implementation starts with 47 unresolved questions → constant backtracking, contradictory code, wasted work.
- **Medium risk**: Self-hosted single-user app limits damage from missing auth, ops, and performance contracts.
- **Low risk**: Dataset volume expectations (5yr) minimize immediate data-model scaling issues.

## Conditional Partial-Approval Sub-Gate Analysis

Even though the full package cannot reach approved-to-dev, some scopes can proceed while others remain blocked:

### Scopes That Can Receive Qualified Approval (with explicit dependency list)
| Scope | Precondition | Can Proceed When | Remaining Blockers After Precondition |
|---|---|---|---|
| data-contracts | TQ-API-001 (protocol) + TQ-ARCH-002 (component arch) resolved | Endpoint catalog + schema format known | ~4 dev-blocking (userId ambiguity, enums, migration tool, stale refs) |
| testing-delivery | TQ-API-001 (protocol) + TQ-API-003 (schemas) resolved | API shapes known → test fixtures designable | ~3 dev-blocking (e2e strategy, export schema, backup schema) |

### Scopes That Cannot Receive Any Approval Until Foundational Decisions
| Scope | Minimum Foundation Needed | Reason |
|---|---|---|
| api-contracts | TQ-ARCH-002 (component arch) + TQ-API-001 (protocol) | No API design possible without protocol choice and system boundaries |
| auth-security-compliance | TQ-ARCH-002 + TQ-ARCH-005 (service boundaries) | Auth middleware placement depends on component architecture |
| client-state-ux | TQ-ARCH-002 + TQ-API-001 + TQ-API-012 (session auth) | UI state machines depend on API surface and auth flow |
| integrations-events | TQ-ARCH-005 (background jobs) + TQ-AUTH-002 (session TTL) | Async job model depends on service boundaries and session contract |
| operations-observability | TQ-ARCH-004 (deployment topology) | Ops config depends on deployment architecture |
| architecture-boundaries | None (already self-contained) | Architecture gaps can be resolved independently — highest priority scope |

## Resolution Blocker Ranking by Estimated Effort

| Rank | Question(s) | Type | Est. Effort | Owner | Notes |
|---|---|---|---|---|---|
| 1 | TQ-API-001 (protocol choice) | Decision | 30 min | Tech lead | REST vs GraphQL — clear tradeoffs, small decision |
| 2 | TQ-AUTH-002 (session TTL + cookie) | Decision | 30 min | Tech lead / Product | Standard values apply; low risk if wrong |
| 3 | TQ-ARCH-001 + TQ-ARCH-002 (system context + component arch) | Design artifact | 2–4 hours | Architect | Single diagram + component table. Highest leverage. |
| 4 | TQ-CLIENT-001 (UI state machine) | Design artifact | 2–4 hours | UX / Frontend lead | Convention + page examples. Can parallelize with arch. |
| 5 | TQ-AUTH-001 + TQ-AUTH-005 + TQ-AUTH-006 (PIN + session + media access) | Decision + Design | 2–3 hours | Tech lead + Product | Includes resolving RULE-022 vs RULE-024 contradiction |
| 6 | TQ-OPS-001 + TQ-OPS-003 (env topology + SLO measurement) | Design artifact | 2–4 hours | Ops / Full-stack | Docker Compose spec + instrumentation plan |
| 7 | TQ-INT-001 + TQ-INT-003 (async job + import transaction model) | Design artifact | 2–4 hours | Tech lead | Architecturally significant — affects API, UX, ops |
| 8 | TQ-TEST-001 + TQ-TEST-002 (test data + e2e strategy) | Decision + Plan | 1–2 hours | QA / Full-stack | Pattern decision, not implementation |
| 9 | TQ-API-002 + TQ-API-003 (endpoint catalog + schemas) | Design artifact | 4–8 hours | Tech lead | Large surface but mechanical: endpoint list from entities |
| 10 | TQ-DATA-001 through TQ-DATA-005 (data model contradictions) | Decision + Fix | 2–4 hours | Tech lead | Schema cleanup + enum definitions + migration tool choice |

## Pre-Implementation vs In-Wave Question Classification

### Must-Resolve Before Any Implementation Begins
1. **TQ-API-001** — Protocol choice blocks every other scope
2. **TQ-ARCH-002** — Component architecture blocks all code organization
3. **TQ-AUTH-002** — Session TTL + cookie config blocks auth middleware
4. **TQ-AUTH-006** — RULE-022 vs RULE-024 contradiction blocks media access design
5. **TQ-ARCH-005** — Monolith vs modular + background job decision blocks service layer

### Must-Resolve Before Scope-Specific Implementation
6. **TQ-API-002, TQ-API-003** — Endpoint catalog + schemas before any API code
7. **TQ-CLIENT-001** — UI state machine before any page rendering
8. **TQ-AUTH-001, TQ-AUTH-005** — PIN hash + session token before PIN implementation
9. **TQ-OPS-001** — Environment/config topology before Docker setup
10. **TQ-INT-001, TQ-INT-003** — Async job contract + transaction model before export/backup code

### Can Be Resolved During Early Implementation Waves
11. **TQ-DATA-001 through TQ-DATA-005** — Data model cleanup can run in parallel with arch work
12. **TQ-DATA-011** — Index strategy can wait until schema is generated and tested
13. **TQ-CLIENT-006** — Form validation can be defined per-page during implementation
14. **TQ-CLIENT-011** — Accessibility can be added iteratively during frontend development
15. **TQ-CLIENT-012** — Localization can be deferred if single-locale MVP is accepted
16. **TQ-TEST-001, TQ-TEST-002** — Test patterns should be decided early but fixtures built during implementation
17. **TQ-TEST-005** — Log redaction policy can be refined during implementation
18. **TQ-OPS-004** — Alerting/runbooks (already deferred to post-MVP)
19. **TQ-CLIENT-010** — Realtime/subscriptions (already deferred)

## Format Consistency Check

### Question Numbering Convention
All 8 scopes use the required TQ-{SCOPE}-* prefix:
- TQ-ARCH-* (architecture-boundaries)
- TQ-DATA-* (data-contracts)
- TQ-API-* (api-contracts)
- TQ-AUTH-* (auth-security-compliance)
- TQ-INT-* (integrations-events)
- TQ-CLIENT-* (client-state-ux)
- TQ-OPS-* (operations-observability)
- TQ-TEST-* (testing-delivery)

No numbering collisions, no duplicates, no missing prefixes. ✓

### Question Ledger Format Consistency
All 8 scope ledgers consistently use the required columns:
Scope, Severity, Parent, Question, Why It Matters, Needed Artifact Or Decision, Source Or Report, Status, Resolution

All ledgers use the allowed severity values (dev-blocking, needs-owner-decision, deferred, watchlist). All statuses are "open" or "deferred" with deferral rationale. ✓

## Recommended Next Actions for Controller

1. **Synchronize session TTL severity**: Promote TQ-INT-005 to dev-blocking (matches TQ-AUTH-002).
2. **Request client-state-ux and integrations-events** to review DEC-007 and DEC-009 effects (gap closure).
3. **Request api-contracts** to review DEC-008 performance targets for API layer.
4. **Resolve severity drift**: Align TQ-DATA-008 with TQ-AUTH-009 (both needs-owner-decision for data retention).
5. **Create cross-scope dependency map** for the 6 identified blocking decision chains.
6. **Present blocking summary** to product owner for the 7 foundational artifact decisions.