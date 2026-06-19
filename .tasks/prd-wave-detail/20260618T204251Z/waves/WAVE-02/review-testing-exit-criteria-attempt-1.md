# WAVE-02 testing-exit-criteria Review Attempt 1

## Verdict
needs-revision

## Sources Read
- planner-testing-exit-attempt-1.md
- planner-product-ac-attempt-1.md
- planner-architecture-codebase-attempt-1.md
- planner-data-integration-ops-attempt-1.md
- planner-security-compliance-attempt-1.md
- planner-sequencing-fit-attempt-1.md
- docs/technical-verified/testing-and-delivery.md
- docs/prd-wave-details/waves/wave-01.md (verification obligations)
- docs/product-verified/acceptance-criteria.md
- apps/api/internal/graph/schema_resolvers_test.go

## Coverage Check
23 test obligations proposed across all planners. Coverage spans: unit (repo, service, file validation, pagination), integration (resolver, handler, migration, soft delete, auth), lint, codegen validation.

## Evidence Check
Test patterns match existing codebase (testify/require, httptest, test DB setup). WAVE-01 test infrastructure provides reusable helpers.

## Codebase Fit Check
Existing test patterns: internal test files, integration test files, resolver test files. All proposed tests follow established conventions.

### Issues Found

1. **No explicit WAVE-01 dependency test**: Since WAVE-02 depends on WAVE-01's PIN auth middleware, there should be a test that verifies WAVE-01's PIN auth is correctly plumbed. This could be a smoke test that calls an exercise endpoint without PIN and expects 401/AuthError. TEST-W02-014 covers this partially but doesn't explicitly test the WAVE-01 middleware chain integration.

2. **E2E test gap**: While the active development policy excludes heavy gates, the exit criteria should include at least one exercise round-trip e2e test (create exercise → upload media → verify media → delete media → soft-delete exercise → verify inactive). Add TEST-W02-024 for this round-trip test.

3. **Test fixture strategy**: The planner mentions "test fixtures create deterministic exercises" but doesn't specify the format. Should use WAVE-01's test helper pattern. Add a note about fixture style (factory functions vs raw SQL vs JSON).

4. **EC-W02-001 is too vague**: "All acceptance criteria passing" is not independently verifiable. Should reference the specific test patterns that prove each AC. Either break into per-AC ECs or reference the TEST-W02-* list.

5. **EC-W02-003 (sqlc codegen) depends on sqlc generation**: This should be paired with a `make generate` or equivalent script execution to verify no drift. The verification command in TEST-W02-006 covers this.

6. **Missing migration test**: TEST-W02-005 (migration smoke test) is mentioned but all planners must ensure the migration test covers both the exercises and exercise_media tables together (00080 and 00081 in sequence).

7. **No test for exercise pagination edge cases**: TEST-W02-020 (pagination) should explicitly test: (a) default page size, (b) cursor with no more results, (c) cursor in middle of dataset. Add as sub-cases.

8. **EC-W02-022 is a negative EC**: "No WAVE-03 functionality" is hard to prove in tests. Should be a code review gate, not an EC. Move to code review checklist.

## Other-Wave Fit Check
WAVE-03 will need exercise fixtures. The planners correctly identify the need for deterministic exercise IDs that WAVE-03 tests can depend on.

## Acceptance Criteria Check
All 23 ACs have corresponding test obligations. AC-to-test mapping is good:
- AC-W02-001 through AC-W02-005 → TEST-W02-001, TEST-W02-002, TEST-W02-003
- AC-W02-006 through AC-W02-008 → TEST-W02-004
- AC-W02-011 through AC-W02-012 → TEST-W02-013
- AC-W02-014 through AC-W02-020 → TEST-W02-014 through TEST-W02-019

## Exit Criteria Check
22 ECs proposed. EC-W02-001 (all AC passing) and EC-W02-022 (no WAVE-03 scope) need revision. See issues above.

## Verification Check
TEST-W02-001 through TEST-W02-023 provide solid coverage. Add TEST-W02-024 (round-trip e2e) as suggested.

## Question Ledger Check
DQ-W02-007 (mock vs integration auth) affects test strategy. Resolving this (prefer integration through full middleware chain) would simplify test design.

## Unsupported Or Invented Claims
No unsupported claims. Test plan is grounded in known patterns.

## Required Revisions
1. **Add TEST-W02-024**: Exercise round-trip integration test (create → upload media → verify → delete media → soft delete → verify).
2. **Strengthen EC-W02-001**: Reference specific TEST-W02-* patterns instead of vague "all ACs passing."
3. **Add EC for migration order**: Verify migrations 00080 and 00081 apply in sequence without errors.
4. **Document fixture strategy**: Specify factory function pattern for test exercise creation.
5. **Remove EC-W02-022**: Move to code review checklist (not independently testable).
6. **Expand TEST-W02-020**: Add cursor edge cases.

## Approval Notes
Good test coverage with clear revision items. After adjustments, will approve.