# WAVE-06 Testing-Exit Planner Attempt 2

## Revisions Applied
1. Added test for measurement overlay with empty types list (returns empty groups)
2. Documented that exercise chart tests are conditional on WAVE-03
3. Added exit criterion EC-W06-009 for measurement overlay empty types

## Proposed Details (Revisions Only)

### Revised Exit Criteria
| EC ID | Description |
|---|---|
| EC-W06-009 | Measurement overlay with empty measurementTypes list returns empty groups array (no error) |

### Revised Verification Obligations
| Test ID | Description | Type | Command |
|---|---|---|---|
| TEST-W06-021 | Measurement overlay with empty types list returns empty groups | integration | bunx nx run api:test -- --run '(?i)measurement_overlay_empty_types' |
| TEST-W06-022 | Measurement overlay returns groups ordered alphabetically by measurement type | integration | bunx nx run api:test -- --run '(?i)measurement_overlay_ordering' |

### Notes on Conditional Tests
- TEST-W06-001 through TEST-W06-003 (exercise chart, e1RM) are conditional on WAVE-03
- If WAVE-03 is not implemented when WAVE-06 is developed, these tests are omitted
- If WAVE-03 is implemented, these tests become active with a pre-implementation gate check
