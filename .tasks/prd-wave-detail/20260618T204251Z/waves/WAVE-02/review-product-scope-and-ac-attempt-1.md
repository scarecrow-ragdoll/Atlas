# WAVE-02 product-scope-and-ac Review Attempt 1

## Verdict
needs-revision

## Sources Read
- planner-product-ac-attempt-1.md
- planner-data-integration-ops-attempt-1.md
- planner-architecture-codebase-attempt-1.md
- planner-security-compliance-attempt-1.md
- planner-testing-exit-attempt-1.md
- planner-sequencing-fit-attempt-1.md
- docs/prd-waves/waves/wave-02.md (source wave)
- docs/product-verified/acceptance-criteria.md
- docs/product-verified/functional-spec.md (Section 11)
- docs/product-verified/domain-model.md
- docs/prd-waves/frontend-pages/page-003.md

## Coverage Check
All 5 capability groups from the source wave are covered by the planners:
- CAP-W02-001 Exercise CRUD → covered (AC-W02-001 through AC-W02-005)
- CAP-W02-002 ExerciseMedia → covered (AC-W02-006 through AC-W02-008, AC-W02-013)
- CAP-W02-003 Working weight → covered (AC-W02-009)
- CAP-W02-004 Muscle groups/description/notes → covered (AC-W02-001, AC-W02-004)
- CAP-W02-005 isActive soft delete → covered (AC-W02-005, AC-W02-011, AC-W02-012)

All 10 product AC references mapped:
- AC-002, AC-003, AC-043, AC-044 → AC-W02-001, AC-W02-009
- AC-004, AC-045, AC-046 → AC-W02-006, AC-W02-007
- AC-047 → AC-W02-005, AC-W02-011

## Evidence Check
All AC claims trace to source documents (PRD, product-verified, domain-model). No unsourced claims.

## Codebase Fit Check
Exercise CRUD is new domain code—no conflict with existing admin auth or user management.

## Other-Wave Fit Check
WAVE-03 dependency on GET /api/exercises covered by AC-W02-021 (allExercises query). No scope collision with later waves.

## Acceptance Criteria Check
23 ACs proposed across all planners. Issues found:

### Required Revisions
1. **AC-W02-009 (working weight)** is merged into AC-W02-001 (exercise creation includes workingWeight). These should be separate: AC for store+retrieve working weight is valid but should explicitly test storage fidelity (stored value == retrieved value).
2. **AC-W02-013 (media download)** duplicates WAVE-01 AC-W01-010. WAVE-02 should test that exercise-media download works specifically (not generic media download). Rename to focus on exercise-media context.
3. **AC-W02-019 (no personalNotes in logs)** is valid but hard to test deterministically. Should be an EC, not an AC — EC-W02-013 covers this.
4. **Missing AC for exercise list pagination**: AC-W02-003 says "listed with pagination" but doesn't test cursor behavior or default page size.
5. **Missing AC for exercise name uniqueness**: EDGE-002 says no duplicate rule. Should have explicit AC that duplicate names are allowed (system behavior).
6. **Missing AC for restoring a soft-deleted exercise**: No reactivate mutation defined. Should be explicit that reactivation is not in scope (or add it if needed).

### Summary
23 ACs → 18 unique ACs needed after merge/deduplication. 5 missing or combined ACs identified.

## Exit Criteria Check
22 ECs across planners. EC-W02-004 (media scaffold extension) depends on WAVE-01 — acceptable as dependency. EC-W02-006 (existing auth unchanged) should reference WAVE-01 test suite.

## Verification Check
23 test obligations. Good coverage of unit, integration, lint, codegen. TEST-W02-019 (log sanitize) may be fragile — consider replacing with a more stable assertion pattern.

## Question Ledger Check
8 open questions: DQ-W02-001 through DQ-W02-008. DQ-W02-001 (delete physical file) and DQ-W02-004 (same topic flagged by two planners) should be merged into one question. DQ-W02-006 (signed URLs) is deferred-compatible. DQ-W02-005 (MIME detection) and DQ-W02-007 (mock vs integration) need owner decisions.

## Unsupported Or Invented Claims
None found. All proposed ACs trace to source docs.

## Required Revisions
1. Split AC-W02-009 into separate "working weight stored faithfully" AC
2. Repurpose AC-W02-013 to focus on exercise-media specific download
3. Move AC-W02-019 to EC (log sanitization is non-functional, not behavioral)
4. Add AC for list pagination cursor behavior
5. Add AC for duplicate exercise names being allowed
6. Add AC or explicit exclusion for reactivate exercise
7. Merge DQ-W02-001 and DQ-W02-004 (same topic)
8. Add AC for exercise field update persistence (update name → query returns new name)

## Approval Notes
Strong foundation. After addressing 8 revision items, this perspective will approve.