# Wave 03: Workout Diary

## Status

user-approved

## User Approval

user-approved (2026-06-18)

## Purpose

Workout day by date with exercises and sets. Core workout tracking functionality.

## Outcome After Wave

- OUT-W03-001 Workouts can be logged by date
- OUT-W03-002 Exercises added with snapshot working weight
- OUT-W03-003 Sets with weight, reps, optional RPE/RIR
- OUT-W03-004 Cardio can be added per day
- OUT-W03-005 Data available for charts/AI

## Included Scope

- CAP-W03-001 WorkoutDay CRUD by date
- CAP-W03-002 WorkoutExercise with order
- CAP-W03-003 WorkoutSet with weight/reps/RPE/RIR
- CAP-W03-004 Working weight snapshot from Exercise
- CAP-W03-005 Cardio inline or linked
- CAP-W03-006 Comments per exercise
- CAP-W03-007 Calendar date switching

## Excluded Scope

- Exercise CRUD (from WAVE-02)
- Charts visualization
- AI export

## Dependencies

WAVE-01, WAVE-02

## Surface Categories

backend, data, operations

## Risk Class

Medium - Date handling, working weight snapshots

## Recommended Next Planning

$detail-prd-wave for WAVE-03

## Open Questions

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Q-WORKOUT-001 | 03 | operations | Medium | None | Concurrent edit handling? | Data integrity | docs/product/prd.md Section 24.3 | open | needs-owner-decision |

## Traceability

- docs/product/prd.md Section 10
- docs/product-verified/user-flows.md#Log Workout Today