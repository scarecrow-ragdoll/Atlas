# WAVE-05 Product-Scope-and-AC Review Attempt 1

## Verdict
needs-revision

## Sources Read
- planner-product-ac-attempt-1.md
- docs/prd-waves/waves/wave-05.md
- docs/product-verified/acceptance-criteria.md
- docs/product-verified/edge-cases.md
- docs/product-verified/business-rules.md
- docs/product-verified/domain-model.md

## Coverage Check
All product-level ACs mapped: AC-017–AC-019, AC-058–AC-064, AC-113. All edge cases covered: EDGE-003, EDGE-009, EDGE-017, EDGE-019. All business rules: RULE-006, RULE-010, RULE-011, RULE-018, RULE-019, RULE-020. Good.

## Evidence Check
Claims trace to source docs. No unsupported inventions.

## Codebase Fit Check
Not applicable for this perspective.

## Other-Wave Fit Check
Not applicable for this perspective.

## Acceptance Criteria Check
34 ACs mapped. Good coverage. Some gaps:
1. AC-W05-006 says "hard deletion blocked when referenced by active template/override items" but only soft-delete (isActive flag) is mentioned. Need consistency across reports — is delete soft or hard-blocked? The data planner uses soft-delete.
2. AC-W05-029 mentions "for each day" but the query interface is per-week or per-day. The wording should match the actual query shape.
3. AC-W05-032 mentions "warning marker" — what form does this take in the API response? Need clarification.

## Exit Criteria Check
12 ECs cover all ACs. Adequate.

## Verification Check
29 test IDs cover all major code paths. Adequate.

## Question Ledger Check
3 questions raised (DQ-W05-001, DQ-W05-002, DQ-W05-003). All are reasonable.

## Unsupported Or Invented Claims
- The `mealLabel` as free-text is an assumption (not specified in PRD). This is fine — documented as decision in DQ-W05-003.
- Soft-delete (isActive) for NutritionProduct is an assumption. PRD doesn't specify deletion behavior. Documented in DQ-W05-001.

## Required Revisions
1. Make the soft-delete vs hard-blocked decision consistent across ACs. Either (a) soft-delete with isActive flag (products remain in DB, excluded from queries by default) or (b) hard-delete blocked when referenced. The data planner says soft-delete. If soft-delete: AC-W05-006 should say "soft-deleted (isActive=false)" not "hard deletion blocked."
2. Add AC for: soft-deleted product is NOT visible in active product list but IS visible in template/override items (historical reference).
3. Add a note about the "warning marker" for soft-deleted products in KJBJU — suggest returning 0 for that item's macros with no special error, as the product simply contributes 0.

## Approval Notes
Good foundation. Minor consistency fixes needed. Approve after revisions.