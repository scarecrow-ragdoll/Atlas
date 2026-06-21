# WAVE-06 Product-AC Planner Attempt 2

## Revisions Applied
1. Split AC-W06-001 into separate ACs per chart data type
2. Added AC for exercise chart empty series (conditional/stubbed)
3. Moved default period to open question (DQ-W06-004)
4. Added AC for measurement overlay with empty types list
5. Added AC for measurement overlay ordering
6. Clarified exercise progress ACs as "conditional on WAVE-03"

## Proposed Details (Revisions Only)

### Outcome
Same as Attempt 1. Exercise progress queries conditional on WAVE-03.

### Revised Acceptance Criteria

| AC ID | Description | Status |
|---|---|---|
| AC-W06-001 | Exercise progress query for a given exercise over a date range returns per-session time-series data including: session date, working weight, best set weight, best set reps, e1RM, volume, total reps, working sets count | conditional (WAVE-03) |
| AC-W06-002 | e1RM calculated per-set using Epley formula: weight × (1 + reps / 30). Per-session e1RM = best (highest) e1RM value among all sets in that session | conditional (WAVE-03) |
| AC-W06-003 | Exercise chart query returns empty series (zero items, no error) when no data exists for the selected period/exercise | conditional (WAVE-03) |
| AC-W06-004 | Body weight trend query returns time-series data (date, weight) over a given date range, ordered by date ascending | implementable |
| AC-W06-005 | Body weight trend query returns empty series when no data exists for the selected period | implementable |
| AC-W06-006 | Body measurement trend query returns time-series data (date, measurementType, side, value) for a given measurement type over a date range | implementable |
| AC-W06-007 | Body measurement overlay query returns time-series data for multiple measurement types in a single response, grouped by measurementType, ordered alphabetically by type name | implementable |
| AC-W06-008 | Body measurement queries return empty series when no data exists for the selected period | implementable |
| AC-W06-009 | Body measurement overlay query with empty measurementTypes list returns empty groups array (no error) | implementable |
| AC-W06-010 | Nutrition weekly macro averages query returns time-series data (weekStartDate, calories, protein, fat, carbs) over a given date range | implementable |
| AC-W06-011 | Nutrition weekly averages use RULE-015: sum of daily values / 7 | implementable |
| AC-W06-012 | All chart queries accept optional from/to date parameters for period filtering | implementable |
| AC-W06-013 | Period filter with from > to date returns ValidationError | implementable |
| AC-W06-014 | All chart queries return AuthError when PIN session header is missing or invalid | implementable |
| AC-W06-015 | Nutrition weekly averages query returns empty series when no data exists for the selected period | implementable |

### Revised Exit Criteria
| EC ID | Description |
|---|---|
| EC-W06-001 | AC-W06-001 through AC-W06-015 pass via TEST-W06-001 through TEST-W06-022 |
| EC-W06-002 | All chart queries protected by WAVE-01 PIN auth middleware |
| EC-W06-003 | gqlgen codegen produces valid Go code for WAVE-06 schema without drift |
| EC-W06-004 | Empty series returned for no-data periods (no errors) — verified for body, measurement, and nutrition queries |
| EC-W06-005 | Nutrition weekly average calculation matches RULE-015 |
| EC-W06-006 | Lint passes for all changed packages |
| EC-W06-007 | No sensitive data (body weight values, measurement values) in application logs |
| EC-W06-008 | Body weight trend query returns accurate date-ascending data |
| EC-W06-009 | Measurement overlay returns correct multi-type structure with alphabetical ordering |

### Revised Verification Obligations
Identical to attempt 1 with additions:
- TEST-W06-021: Measurement overlay with empty types list returns empty groups array
- TEST-W06-022: Measurement overlay returns groups ordered alphabetically by type

## Questions Raised (Updated)
(All existing DQ-W06 questions remain. DQ-W06-004 added for default period.)
