# Decision Log

## Source Wave Gate
- WAVE-03 source gate: passed (2026-06-18)
- Source wave: docs/prd-waves/waves/wave-03.md
- Gate check: Q-WORKOUT-001 is not decomposition-blocking or owner-decision affecting wave boundary

## User Wave Approvals
- WAVE-01: user-approved (2026-06-18), ready-for-dev
- WAVE-02: user-approved (2026-06-18)
- WAVE-03: pending user approval of detailed brief

## Scope Decisions
| Decision ID | Description | Source | Date |
| --- | --- | --- | --- |
| DDEC-W03-001 | DailyLog created on first content save (exercises or cardio). Empty dates have no record. | planner-product-ac-attempt-1.md, DQ-W03-005 | 2026-06-18 |
| DDEC-W03-002 | WorkoutExercise.display_order determines display order within DailyLog. | planner-architecture-codebase-attempt-1.md | 2026-06-18 |
| DDEC-W03-003 | WorkoutSet.set_number determines order within WorkoutExercise, starts at 1, sequential. | planner-architecture-codebase-attempt-1.md | 2026-06-18 |
| DDEC-W03-004 | Cascade delete: DailyLog -> WorkoutExercise -> WorkoutSet and DailyLog -> CardioEntry. | planner-data-integration-ops-attempt-1.md, DQ-W03-006 | 2026-06-18 |
| DDEC-W03-005 | Exercise FK: workout_exercises.exercise_id -> exercises(id) ON DELETE NO ACTION. Preserves workout history. | planner-architecture-codebase-attempt-1.md | 2026-06-18 |
| DDEC-W03-006 | Concurrent edit: last-write-wins for MVP. No optimistic locking. | planner-data-integration-ops-attempt-1.md, DQ-W03-001 | 2026-06-18 |
| DDEC-W03-007 | Backdating: no lower bound on date. Any valid DATE accepted. | planner-product-ac-attempt-1.md, DQ-W03-004 | 2026-06-18 |
| DDEC-W03-008 | Working weight snapshot captured at exercise-add time, not save time. | planner-product-ac-attempt-1.md, DQ-W03-007 | 2026-06-18 |
| DDEC-W03-009 | CardioEntry requires dailyLogId FK. No standalone cardio in WAVE-03. | planner-sequencing-fit-attempt-1.md | 2026-06-18 |

## Codebase Fit Decisions
| Decision ID | Description | Source | Date |
| --- | --- | --- | --- |
| DDEC-W03-010 | Migration numbering starts at 00082 (after WAVE-02 00080/00081). | planner-architecture-codebase-attempt-1.md | 2026-06-18 |
| DDEC-W03-011 | All WAVE-03 operations via GraphQL. No REST endpoints. | planner-data-integration-ops-attempt-1.md | 2026-06-18 |
| DDEC-W03-012 | Auto-discovery via gqlgen/sqlc globs — no config changes needed for new files. | planner-architecture-codebase-attempt-1.md | 2026-06-18 |

## Deferrals
| Deferral ID | Description | Target | Reason | Date |
| --- | --- | --- | --- | --- |
| DDEC-W03-013 | Optimistic concurrency control for DailyLog updates | Post-MVP | MVP uses last-write-wins; concurrent edit acknowledged risk | 2026-06-18 |
| DDEC-W03-014 | Workout templates | Future scope | Explicitly deferred in PAGE-002 | 2026-06-18 |
| DDEC-W03-015 | Bulk exercise reorder (drag and drop) | Post-MVP | Per-exercise order update sufficient for MVP | 2026-06-18 |

## Rejected Assumptions
| Assumption | Reason for Rejection | Source | Date |
| --- | --- | --- | --- |
| CardioEntry should be planned in WAVE-04 only | WAVE-03 domain model shows CardioEntry with required dailyLogId — inline cardio is WAVE-03 scope. WAVE-04 handles standalone variant. | docs/product-verified/domain-model.md | 2026-06-18 |
| DailyLog uses "WorkoutDay" name | Domain model confirms DailyLog rename from WorkoutDay per DEC-009. | docs/product-verified/domain-model.md | 2026-06-18 |
