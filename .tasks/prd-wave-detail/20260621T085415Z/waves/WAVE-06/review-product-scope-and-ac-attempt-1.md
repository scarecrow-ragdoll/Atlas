# WAVE-06 Product-Scope-And-AC Review Attempt 1

## Verdict
needs-revision

## Sources Read
- planner-product-ac-attempt-1.md
- planner-architecture-codebase-attempt-1.md
- planner-data-integration-ops-attempt-1.md
- planner-security-compliance-attempt-1.md
- planner-testing-exit-attempt-1.md
- planner-sequencing-fit-attempt-1.md
- docs/prd-waves/waves/wave-06.md
- docs/product-verified/acceptance-criteria.md
- docs/product-verified/business-rules.md

## Coverage Check
All source ACs (AC-020-022, AC-065-073) are covered. Edge cases (EDGE-008, EDGE-026) addressed.

## Evidence Check
Source wave ACs mapped to wave-local ACs. However, AC-W06-001 bundles too many chart types into one AC — should split per chart data type.

## Codebase Fit Check
Planners correctly identified WAVE-03 gap for exercise charts. Product AC AC-W06-001 should be a stub/conditional.

## Other-Wave Fit Check
No scope collision with WAVE-04 or WAVE-05.

## Acceptance Criteria Check
Issues:
1. AC-W06-001 is a mega-AC covering 8 data fields from 5 chart types. Should split: AC-W06-001 (exercise progress data structure), separate ACs for each query type.
2. No AC for "exercise chart returns empty series" — EDGE-008 applies to exercise charts too, even if stubbed.
3. AC-W06-012 default period (12 weeks) — not sourced from PRD. PRD says "last 4 weeks" for AI export (AC-074). Chart default should be "last 12 weeks" but needs source evidence or explicit decision.
4. AC-W06-013 (from > to returns ValidationError) — not sourced from any PRD requirement. Acceptable design choice but should be flagged.
5. Need AC for measurement overlay ordering (by measurementType).

## Exit Criteria Check
EC-W06-001 through EC-W06-008 — adequate coverage.

## Verification Check
TEST list covers 20 tests. Missing: test for measurement overlay with 0 selected types (should return empty groups).

## Question Ledger Check
DQ-W06-001 (best set definition) correctly identifies open question. DQ-W06-003 (working weight source) correctly identifies open question. DQ-W06-002 (WAVE-03 dependency) correctly identified.

## Unsupported Or Invented Claims
1. Default 12-week period — needs decision source.
2. 52-week max range — needs decision source.
3. Measurement overlay ordering assumption.

## Required Revisions
1. Split AC-W06-001 into separate ACs per chart data type (exercise, body weight, measurement, nutrition).
2. Add AC for exercise chart empty series (even if stubbed).
3. Add default period to open questions or provide source evidence.
4. Add AC for measurement overlay with empty types list.
5. Add AC for measurement overlay ordering contract.
6. Clarify whether exercise progress ACs are "implemented" or "stubbed/conditional on WAVE-03".

## Approval Notes
Product scope is correct. Source wave boundaries respected. ACs need splitting for traceability but overall direction is sound. Awaiting revision.