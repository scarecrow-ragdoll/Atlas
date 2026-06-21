# Decision Log

## Source Wave Gate
source-wave-gate: passed for WAVE-06 (2026-06-21). Source: docs/prd-waves/waves/wave-06.md.

## User Wave Approvals
- WAVE-06 source wave: user-approved (2026-06-18) via wave-map.md
- Detailed WAVE-06: awaiting user approval

## Scope Decisions
- Read-only wave: no mutations, no storage changes (DDEC-W06-001)
- Epley formula for e1RM (DDEC-W06-002)
- Measurement range via check_in JOIN (DDEC-W06-003)
- Nutrition weekly average via iteration per RULE-015 (DDEC-W06-004)
- Empty series for no-data periods (DDEC-W06-005)
- Default chart period: 4 weeks (DDEC-W06-006)
- Measurement overlay alphabetically ordered (DDEC-W06-007)
- Best set = highest e1RM per session (DDEC-W06-008)
- Working weight per session from WorkoutExercise.workingWeightSnapshot (DDEC-W06-009)
- Exercise chart stubs returning empty series until WAVE-03 (DDEC-W06-010)
- 52-week max date range cap (DDEC-W06-011)

## Codebase Fit Decisions
- All new files in apps/api/internal/atlas/ — consistent with existing module structure
- sqlc auto-discovery via glob — no config changes needed
- atlas-gqlgen.yml needs new model bindings for chart types

## Deferrals
- Exercise chart queries: stubs returning empty series until WAVE-03 deployment
- 52-week max range: enforcement via server constant

## Rejected Assumptions
None.