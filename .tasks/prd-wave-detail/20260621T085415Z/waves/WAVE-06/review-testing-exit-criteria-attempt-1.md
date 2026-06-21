# WAVE-06 Testing-Exit-Criteria Review Attempt 1

## Verdict
needs-revision

## Sources Read
- planner-testing-exit-attempt-1.md
- planner-product-ac-attempt-1.md
- planner-architecture-codebase-attempt-1.md
- planner-data-integration-ops-attempt-1.md
- planner-security-compliance-attempt-1.md
- docs/prd-wave-details/waves/wave-05.md (TEST section pattern)

## Coverage Check
20 tests proposed. Coverage: e1RM unit, body weight trend (happy + empty), measurement trend (happy + empty), nutrition weekly avg (happy + empty), date validation, auth, codegen, lint, log privacy, max range, side filter, single point, partial week.

## Evidence Check
Test patterns match WAVE-04/WAVE-05 conventions. Commands use `bunx nx run api:test -- --run '(?i)pattern'` consistently.

## Codebase Fit Check
Integration tests require DB-backed queries. Measurement range test depends on new sqlc query — test is valid.

## Other-Wave Fit Check
No test collision with WAVE-04 or WAVE-05 tests. Separate test regex patterns prevent overlap.

## Acceptance Criteria Check
All 15 ACs have corresponding tests. However:
- AC-W06-001 (exercise progress) has no test due to WAVE-03 dependency — should be documented as conditional.
- AC-W06-002 (e1RM formula) covered by TEST-W06-001.
- AC-W06-014 (auth) covered by TEST-W06-011.

## Exit Criteria Check
EC-W06-004 (empty series), EC-W06-005 (RULE-015), EC-W06-007 (lint), EC-W06-008 (body weight ordering) all have matching tests.

## Verification Check
All verification obligations are either unit or integration tests. No e2e tests needed for read-only queries. Codegen drift and lint included.

## Question Ledger Check
DQ-W06-008 (side inclusion in measurement overlay) appropriately raised.

## Unsupported Or Invented Claims
None directly. Some tests assume specific implementation details (e.g., default 12-week period) that need decisions first.

## Required Revisions
1. Add test for measurement overlay with empty types list (edge case).
2. Document that exercise chart tests are conditional on WAVE-03.
3. Add exit criterion EC for measurement overlay returning empty groups when types list empty.

## Approval Notes
Test strategy is sound. Coverage adequate for read-only query surface. Missing exercise chart tests are honest about dependency. Awaiting revision for edge cases.