# WAVE-06 Product-AC Planner Attempt 1

## Sources Read
- docs/prd-waves/waves/wave-06.md
- docs/product-verified/features/charts.md
- docs/product-verified/acceptance-criteria.md (AC-020-022, AC-065-073)
- docs/product-verified/business-rules.md (RULE-012, RULE-013, RULE-014, RULE-015)
- docs/product-verified/edge-cases.md (EDGE-008, EDGE-026)
- docs/prd-waves/frontend-pages/page-008.md
- docs/prd-wave-details/waves/wave-05.md
- docs/prd-wave-details/waves/wave-04.md

## Selected Backend Wave Boundary
WAVE-06 provides ONLY backend data queries for chart consumption. No mutations, no new storage, no frontend rendering. This is a pure query/aggregation wave.

## Neighboring Backend Wave Fit
- WAVE-03 (Workout Diary): NOT implemented. Required for exercise chart queries that need set-level data. WAVE-06 cannot implement exercise progress queries until WAVE-03 creates the daily_log_exercises and workout_sets tables.
- WAVE-04 (Cardio/Body Tracking): Fully implemented. BodyWeightEntry and BodyMeasurement repos/services exist. WAVE-06 can query body weight and measurements directly.
- WAVE-05 (Nutrition): Fully implemented. NutritionMacroService exists. WAVE-06 can reuse macro calculation for weekly averages.

## Frontend Pages Context
- PAGE-008 (Charts): Needs exercise selector (exercise list query already exists), date range filter, chart type selector. Backend provides raw time-series data — frontend shapes it for specific chart types.

## Codebase Evidence
- No workout_set table or model exists. Only daily_logs table exists with id, user_id, date, notes.
- ExerciseService provides ListAll (for exercise selector dropdown).
- BodyWeightService provides ListByDateRange (date filter).
- BodyCheckInService and BodyMeasurementService exist with measurement queries.
- NutritionMacroService provides per-week calculation — WAVE-06 needs a weekly-average variant.

## Proposed Details

### Outcome
- Exercise progress queries (weight, 1RM, volume, reps, working sets)
- Body weight trend queries (date range)
- Body measurement trend queries (single + overlay)
- Nutrition weekly macro average queries
- Period filtering on all queries

### Key Design Decisions
- Q-CHART-001: RESOLVED — Epley formula: weight × (1 + reps / 30)
- Best set definition: OPEN (DQ-W06-001)
- Working weight source: OPEN (DQ-W06-003)

## Acceptance Criteria Contributions

| AC ID | Description |
|---|---|
| AC-W06-001 | Exercise progress query returns per-session time-series data for a given exercise over a date range, including: session date, working weight, best set weight, best set reps, e1RM, volume, total reps, working sets count |
| AC-W06-002 | e1RM calculated per-set using Epley formula: weight × (1 + reps / 30). Per-session e1RM = best (highest) e1RM value among all sets in that session |
| AC-W06-003 | Exercise chart queries return empty series (zero items, no error) when no data exists for the selected period/exercise |
| AC-W06-004 | Body weight trend query returns time-series data (date, weight) over a given date range, ordered by date ascending |
| AC-W06-005 | Body weight trend query returns empty series when no data exists for the selected period |
| AC-W06-006 | Body measurement trend query returns time-series data (date, measurementType, side, value) for a given measurement type over a date range |
| AC-W06-007 | Body measurement overlay query returns time-series data for multiple measurement types in a single response |
| AC-W06-008 | Body measurement queries return empty series when no data exists for the selected period |
| AC-W06-009 | Nutrition weekly macro averages query returns time-series data (weekStartDate, calories, protein, fat, carbs) over a given date range |
| AC-W06-010 | Nutrition weekly averages use RULE-015: sum of daily values / 7 |
| AC-W06-011 | All chart queries accept optional from/to date parameters for period filtering |
| AC-W06-012 | Period filter defaults to last 12 weeks when not specified |
| AC-W06-013 | Period filter with from > to date returns ValidationError |
| AC-W06-014 | All chart queries return AuthError when PIN session header is missing or invalid |
| AC-W06-015 | Nutrition weekly averages query returns empty series when no data exists for the selected period |

## Exit Criteria Contributions
| EC ID | Description |
|---|---|
| EC-W06-001 | AC-W06-001 through AC-W06-015 pass via TEST-W06-001 through TEST-W06-020 |
| EC-W06-002 | All chart queries protected by WAVE-01 PIN auth middleware |
| EC-W06-003 | gqlgen codegen produces valid Go code for WAVE-06 schema without drift |
| EC-W06-004 | Exercise chart query returns correct e1RM using Epley formula |
| EC-W06-005 | Empty series returned for no-data periods (no errors) |
| EC-W06-006 | Nutrition weekly average calculation matches RULE-015 |
| EC-W06-007 | Lint passes for all changed packages |
| EC-W06-008 | Body weight trend query returns accurate date-ascending data |

## Verification Contributions
| Test ID | Description | Type |
|---|---|---|
| TEST-W06-001 | Exercise progress query returns correct time-series structure | integration |
| TEST-W06-002 | e1RM calculation matches Epley formula | unit |
| TEST-W06-003 | Exercise chart empty series for no data period | integration |
| TEST-W06-004 | Body weight trend query returns correct data | integration |
| TEST-W06-005 | Body weight trend query returns empty series for no data | integration |
| TEST-W06-006 | Body measurement trend query returns correct data | integration |
| TEST-W06-007 | Body measurement overlay query returns multiple types | integration |
| TEST-W06-008 | Nutrition weekly average query returns correct RULE-015 values | integration |
| TEST-W06-009 | Nutrition weekly average query returns empty series for no data | integration |
| TEST-W06-010 | Period filter validation (from > to returns error) | unit |
| TEST-W06-011 | Period filter default (12 weeks when not specified) | integration |
| TEST-W06-012 | Auth check — all chart queries return AuthError without PIN session | integration |

## Risks And Rollback
- WAVE-03 not implemented — exercise chart queries depend on workout sets data model. Mitigation: WAVE-06 provides all other chart queries (body, nutrition) that work independently. Exercise chart queries are documented as requiring WAVE-03.
- No session snapshot of working weight per exercise — RULE-017 defines static workingWeight on exercise. WAVE-03 should record workingWeight per daily_log_exercise.

## Questions Raised
- DQ-W06-001: Best set — heaviest weight or highest e1RM?
- DQ-W06-002: Exercise chart queries need WAVE-03 sets table — does WAVE-06 create its own query infrastructure or wait?
- DQ-W06-003: Working weight source — static exercise.workingWeight or per-session snapshot?

## Traceability Candidates
- docs/prd-waves/waves/wave-06.md (source boundary)
- docs/product-verified/acceptance-criteria.md (AC-020-022, AC-065-073)
- docs/product-verified/business-rules.md (RULE-012 Epley, RULE-013 volume, RULE-014 best set, RULE-015 weekly avg)
- docs/product-verified/edge-cases.md (EDGE-008, EDGE-026)
- docs/prd-waves/frontend-pages/page-008.md (PAGE-008 backend deps)