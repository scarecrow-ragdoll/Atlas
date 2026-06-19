# WAVE-04 Review: Product Scope and AC

## Review Cycle
1

## Planner Reports Reviewed
- planner-product-ac-attempt-1.md
- planner-architecture-codebase-attempt-1.md
- planner-data-integration-ops-attempt-1.md
- planner-security-compliance-attempt-1.md
- planner-testing-exit-attempt-1.md
- planner-sequencing-fit-attempt-1.md

## Verdict
approved

## Findings

### Scope Coverage
The 44 proposed ACs (AC-W04-001 through AC-W04-044) comprehensively cover:
- CardioEntry CRUD: AC-W04-001 through AC-W04-010 (type, duration, pulse, zone, dailyLog linkage)
- BodyWeightEntry CRUD: AC-W04-011 through AC-W04-018 (date, weight > 0, source enum, latest query)
- BodyCheckIn CRUD: AC-W04-019 through AC-W04-025 (weight, bodyFat%, cascade delete)
- BodyMeasurement CRUD: AC-W04-026 through AC-W04-031 (type enum, side validation for paired types)
- ProgressPhoto CRUD: AC-W04-032 through AC-W04-038 (upload, MIME validation, delete cascade)
- WeekFlag CRUD: AC-W04-039 through AC-W04-042 (flagType enum, week listing)
- Auth: AC-W04-043 through AC-W04-044 (PIN auth for GraphQL and REST)

### Edge Cases Covered
- EDGE-006 (photo count): addressed as soft guidance, no hard block. Recommend adding this to AC.
- EDGE-007 (measurement value 0/negative): covered by AC-W04-028 (value > 0)

### Product AC Traceability
All product-level ACs (AC-012 through AC-016, AC-048 through AC-057) are mapped to WAVE-04 ACs.

### Missing Considerations
1. **BodyWeightEntry uniqueness per date**: Not addressed in ACs. Allow multiple entries per date (source varies: scale vs manual). Add AC for this.
2. **Soft delete or hard delete?**: All deletions appear to be hard delete. Verify this is consistent with WAVE-02 (soft delete for exercises). **Recommend:** Hard delete is acceptable for cardio entries, body weight entries, measurements, photos, week flags — no referential integrity concerns.

### Required Revisions
- None. All scope covered. Proceed to synthesis.

## Notes
- The 2-4 photo count requirement (RULE-005) is ambiguous per EDGE-006. The planner recommends soft guidance. This should be documented as a design decision.
- All ACs use stable W04 prefix as required.