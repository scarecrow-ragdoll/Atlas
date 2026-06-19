# Testing-Delivery Worker Attempt 1

## Sources Read
- docs/product-verified/product-brief.md (quality gates §Quality Gates, DEC-006)
- docs/product-verified/acceptance-criteria.md (AC-028, AC-121-125)
- docs/product-verified/functional-spec.md (all capability areas)
- docs/product-verified/edge-cases.md (EDGE-001-031)
- docs/product-verified/open-questions.md
- docs/product-verified/source-inventory.md
- docs/product-verified/scope.md
- docs/verification-plan.xml (V-M-COVERAGE-GATE, V-M-API, V-M-WEB-ADMIN, V-M-WEB)
- docs/technology.xml (§Testing)
- tools/coverage/coverage.config.json
- Existing test files (glob patterns for *_test.go, *.spec.ts, *.test.*)

## Source Delta Reviewed
- DEC-006: Quality gates mandate verify:coverage, e2e weekly workflow, backup integration tests, PIN guard tests, AI export schema tests, log redaction
- DEC-007: Single-user MVP with multi-user-ready data model
- DEC-008: Performance targets defined with p95 thresholds
- DEC-009: DailyLog replaces WorkoutDay, cardio separate entity attached via dailyLogId

## Product Signals

### Quality Gates (DEC-006)
From product-brief.md §Quality Gates:
1. `bun run verify:coverage` passes
2. Critical user flows covered by automated tests
3. E2E tests cover the weekly workflow
4. Backup export/import covered by integration tests
5. PIN guard covered by tests
6. AI export schema covered by snapshot/schema tests
7. No sensitive data written to logs

### Handoff Criteria
AC-121: Round-trip: create exercise -> log workout -> add cardio -> weekly check-in -> AI export
AC-122: PIN enable -> close/open browser -> enter PIN -> access data
AC-123: Create nutrition template -> override one day -> verify different values
AC-124: Create backup -> reset app -> import -> verify restore
AC-125: Run test suite with passing results and coverage gate

### AC-028
User can run verification command with passing tests and coverage.

## Technical Facts

### Existing Test Infrastructure
- Go unit/integration tests for: handlers (health, users), middleware (admin_auth, admin_origin, cors), services (admin_auth, user_service), resolvers (schema, admin_auth), repos (user, admin, redis session, cache), appconfig, testinfra (safe_targets)
- Vitest tests for web-admin: pages (login, ui-kit), hooks (use-mobile), layout (admin-shell), config, graphql-client, admin-navigation, theme-contract
- Vitest tests for web: page, users-client, config, api/users/route, shared/api/users
- Playwright e2e: web-admin (graphql-contract, users-flow), web (rest-users-flow)
- Coverage config: 100% thresholds, 10 allowlist entries (generated code, bootstrap entrypoints)
- Coverage projects: go (api, bot, go-config, go-logger), typescript (web-admin, web, workspace-tools)

### Existing Testing Commands (from technology.xml)
- Module-level: bunx nx test {api,bot,go-config,go-logger,web-admin,web}, bunx nx run web-admin:test-coverage, bunx nx run web:test-coverage
- Wave-level: graphql:validate, api:codegen, web-admin:codegen, web-admin:typecheck, web:typecheck
- Phase-level: bun run lint, bun run test, bun run build, bunx nx run web-admin:e2e, bunx nx run web:e2e, bun run verify:coverage

### Existing Coverage Gaps
- No tests exist for Atlas domain features: exercise library, workout diary, cardio, body tracking, nutrition, charts, AI export, AI review, backup/restore
- These are entirely new feature domains that will need test coverage
- Current tests cover only admin-auth (login/logout/session/CRUD admin users), public REST users, health check, config, logger, and bot handlers

## Technical Gaps

### TGAP-TEST-001: No test fixtures or seed data for Atlas domain entities
The existing test suite has seed data for admin_users and regular users via postgres_test.go helpers. No fixtures exist for: exercises, workouts, sets, cardio entries, body weight entries, body check-ins, measurements, photos, nutrition products, nutrition templates, daily overrides, AI exports, AI reviews, backup ZIPs. Every Atlas feature test will need domain-specific test data builders.

### TGAP-TEST-002: No e2e scenarios for the weekly workflow (AC-121)
Product requires e2e coverage for the full round-trip: create exercise -> log workout -> add cardio -> weekly check-in -> AI export. No Playwright scenarios exist for this. The web-admin and web e2e suites only cover user management CRUD flows.

### TGAP-TEST-003: No integration tests for backup export/import (DEC-006)
Backup/restore is required to be covered by integration tests. The backup feature (AC-093-102) involves manifest.json, data.json, media/, schema validation, dry-run, and full restore. No integration test skeleton exists.

### TGAP-TEST-004: No PIN guard tests (DEC-006)
PIN guard has 9 acceptance criteria (AC-029-034, AC-109-110, AC-117) and 5 edge cases (EDGE-011-015). Sessions, hashing, cookie management, brute-force prevention, and log redaction all need tests. No PIN guard implementation exists yet.

### TGAP-TEST-005: No AI export schema snapshot tests (DEC-006)
The AI export produces structured ZIP output with manifest.json, data.json, summary.md, CSVs, and photos/. Product requires snapshot/schema tests to ensure the export contract does not regress. No export schema or snapshot test infrastructure exists.

### TGAP-TEST-006: No log redaction tests (DEC-006)
AC-117-120 require PIN, AI export content, photos, and sensitive comments to not be logged. No log audit tests exist. The current logger tests (go-logger) test middleware and context propagation but not redaction rules.

### TGAP-TEST-007: No performance/load test plan
Performance targets are defined (product-brief.md §Performance Targets) with p95 thresholds for UI, API, AI export, and backup operations. No performance test scenarios, benchmark harness, or CI gates exist for these targets.

### TGAP-TEST-008: No contract tests for backup/import ZIP structure
Backup manifest schema versioning, data.json entity completeness, and media/ directory structure need contract tests. AC-093-102 define the expected ZIP structure with specific validation rules.

### TGAP-TEST-009: No test isolation strategy for full feature set
The existing postgres_test.go uses schema-level isolation with TRUNCATE between test groups. The full Atlas feature set (exercises, workouts, cardio, body, nutrition, AI) will share entities across domains. A test isolation strategy is needed for cross-domain tests.

### TGAP-TEST-010: No QA handoff checklist or release criteria
Product defines success metrics and quality gates but no explicit QA handoff document, pre-release checklist, or smoke test suite for release validation.

## Missing Source Artifacts

| Missing Artifact | Why | Needed For |
| --- | --- | --- |
| Test data builder/factory library | Every Atlas domain test needs seed data | Deterministic, maintainable tests |
| AI export schema definition | Schema needs to be versioned and validated | Snapshot/schema tests (DEC-006) |
| Backup manifest schema definition | Schema versioning for compatibility checks | Contract tests, migration tests |
| Log redaction specification | Which fields must be redacted and how | Log audit tests (DEC-006) |
| Test isolation policy | Strategy for cross-domain entity cleanup | All integration tests |
| Performance test harness | Benchmark tooling for p95 targets | Performance gates |
| QA handoff checklist | Pre-release validation steps | Release process |

## Questions Raised

| ID | Severity | Question | Why It Matters | Needed Artifact | Status |
| --- | --- | --- | --- | --- | --- |
| TQ-TEST-001 | dev-blocking | What is the test data builder/factory strategy for Atlas domain entities? | Every integration test needs domain seed data (exercises, workouts, sets, cardio, body, nutrition, AI). Without a strategy, tests will be brittle and inconsistent. | Test factory library decision or pattern | open |
| TQ-TEST-002 | dev-blocking | What is the e2e strategy for the weekly workflow (AC-121)? | Product mandates e2e for the full round-trip. Needs Playwright page objects, API helpers, and test data seeding for the 5-step flow. | E2E scenario definition and test plan | open |
| TQ-TEST-003 | dev-blocking | What AI export schema format and versioning strategy is used? | Snapshot/schema tests (DEC-006) need a versioned schema contract to compare against. The PRD defines output files but not the schema format. | AI export schema spec | open |
| TQ-TEST-004 | dev-blocking | What is the backup manifest schema and versioning strategy? | Integration tests (DEC-006) and contract tests need to validate manifest.json structure and schema version compatibility (AC-094, AC-098). | Backup manifest schema spec | open |
| TQ-TEST-005 | needs-owner-decision | What is the log redaction policy? | DEC-006 requires no sensitive data (PIN, AI content, photos, comments) in logs. Needs explicit redaction rules and test coverage. | Log redaction specification | open |
| TQ-TEST-006 | needs-owner-decision | What is the test isolation strategy for cross-domain integration tests? | Atlas entities reference each other (workout->exercise, check-in->photos, template->products). TRUNCATE-based isolation may break cross-domain test scenarios. | Test isolation policy | open |
| TQ-TEST-007 | needs-owner-decision | Should performance tests be part of the MVP release gate? | Performance targets exist but no tooling. Decision needed on whether p95 tests block release or are deferred. | Performance test policy decision | open |
| TQ-TEST-008 | deferred | What is the QA handoff format? | Release criteria exist but no QA checklist format. Can be defined during implementation. | QA handoff template | deferred |
| TQ-TEST-009 | watchlist | Are flaky test detection and retry policies needed? | E2E tests and integration tests may have flaky behavior. No current flaky detection or retry policy. | Flaky test management decision | watchlist |

## Answer Effects
No answered questions from previous runs affect this scope directly. Q-SCOPE-001 (DEC-006 quality gates) is the primary source delta driving testing requirements.

## Risks
1. **Schedule risk**: 7 DEC-006 quality gates requiring new test infrastructure (fixtures, e2e, integration, snapshots, log audit) before MVP is product-ready
2. **Coverage threshold risk**: 100% coverage targets for Go and TypeScript mean every new Atlas feature must have near-complete test coverage; allowlist management needed
3. **E2E flakiness risk**: Weekly workflow e2e spans multiple pages and data mutations; test isolation and cleanup must be robust
4. **Schema evolution risk**: AI export and backup manifest schemas will evolve; snapshot tests must handle schema versioning
5. **Integration test dependency risk**: Backup/restore integration tests need Docker Compose services (PostgreSQL, Redis, filesystem) — same infra as current API tests, but with larger setup

## Suggested Decisions
1. Adopt a builder/factory pattern for test data (Go: functional options or builder structs; TypeScript: test data factories)
2. Use structured Playwright page objects for the weekly workflow e2e
3. Define AI export and backup manifest schemas as JSON Schema files in a `libs/schemas/` package
4. Add `logAudit` test helpers that capture log output and assert redacted fields
5. Defer performance tests to post-MVP with a watchlist item
6. Use transaction-based rollback instead of TRUNCATE for integration tests when feasible

## Traceability Candidates
- DEC-006 -> TQ-TEST-001, TQ-TEST-002, TQ-TEST-003, TQ-TEST-004, TQ-TEST-005
- AC-121 -> TQ-TEST-002
- AC-124 -> TQ-TEST-004
- AC-125 -> AC-028
- AC-117-120 -> TQ-TEST-005
- EDGE-022-023 -> test isolation strategy
- product-brief.md §Performance Targets -> TQ-TEST-007