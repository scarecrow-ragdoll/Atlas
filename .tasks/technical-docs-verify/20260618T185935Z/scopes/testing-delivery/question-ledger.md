# Testing-Delivery Question Ledger

## Open Questions

| ID | Scope | Severity | Parent | Question | Why It Matters | Needed Artifact Or Decision | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| TQ-TEST-001 | testing-delivery | dev-blocking | none | What is the test data builder/factory strategy for Atlas domain entities? | Every integration test needs domain seed data. Without a strategy, tests will be brittle and inconsistent. | Test factory library decision or pattern | worker-attempt-1 | open | TBD |
| TQ-TEST-002 | testing-delivery | dev-blocking | none | What is the e2e strategy for the weekly workflow (AC-121)? | Product mandates e2e for the full round-trip. Needs Playwright page objects, API helpers, and test data seeding. | E2E scenario definition and test plan | worker-attempt-1 | open | TBD |
| TQ-TEST-003 | testing-delivery | dev-blocking | none | What AI export schema format and versioning strategy is used? | Snapshot/schema tests need a versioned schema contract. PRD defines output files but not schema format. | AI export schema spec | worker-attempt-1 | open | TBD |
| TQ-TEST-004 | testing-delivery | dev-blocking | none | What is the backup manifest schema and versioning strategy? | Integration tests need to validate manifest.json structure and schema version compatibility. | Backup manifest schema spec | worker-attempt-1 | open | TBD |
| TQ-TEST-005 | testing-delivery | needs-owner-decision | none | What is the log redaction policy? | DEC-006 requires no sensitive data in logs. Needs explicit redaction rules and test coverage. | Log redaction specification | worker-attempt-1 | open | TBD |
| TQ-TEST-006 | testing-delivery | needs-owner-decision | none | What is the test isolation strategy for cross-domain integration tests? | Atlas entities reference each other. TRUNCATE-based isolation may break cross-domain scenarios. | Test isolation policy | worker-attempt-1 | open | TBD |
| TQ-TEST-007 | testing-delivery | needs-owner-decision | none | Should performance tests be part of the MVP release gate? | Performance targets exist but no tooling. Decision needed on p95 gates. | Performance test policy decision | worker-attempt-1 | open | TBD |
| TQ-TEST-008 | testing-delivery | deferred | none | What is the QA handoff format? | Release criteria exist but no QA checklist format. | QA handoff template | worker-attempt-1 | open | Deferred to implementation phase |
| TQ-TEST-009 | testing-delivery | watchlist | none | Are flaky test detection and retry policies needed? | E2E and integration tests may have flaky behavior. | Flaky test management decision | worker-attempt-1 | open | TBD |

## Answered Questions
None.

## Follow-Up Questions
None.

## Resolved Questions
None.

## Deferred Questions
| ID | Scope | Severity | Question | Rationale |
| --- | --- | --- | --- | --- |
| TQ-TEST-008 | testing-delivery | deferred | QA handoff format | Can be defined during implementation; release criteria in product-brief.md provide sufficient guidance for now |