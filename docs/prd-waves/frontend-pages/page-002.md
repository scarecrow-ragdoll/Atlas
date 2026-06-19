# PAGE-002: Workout Diary

## Status

user-approved

## Page Purpose

Daily workout entry by date. Add exercises, sets with weights/reps/RPE/RIR.

## What Is On This Page

- Date selector (default today)
- Exercise list for selected date
- Add exercise button
- Each exercise with sets table
- Add set button per exercise
- Comments per exercise
- Cardio section for the day
- Save button

## Functional Parts

- Calendar date picker
- Exercise cards with reorder
- Set rows with weight/reps/RPE/RIR
- Auto-working weight from Exercise Library
- Snapshot working weight on save

## Empty States

- No exercises for date - "Add your first exercise"

## Loading And Error States

- Loading day data - skeleton
- Save error - toast with retry

## Backend Dependencies

- GET /api/workouts?date=YYYY-MM-DD
- POST /api/workouts (create day)
- POST /api/workout-exercises
- POST /api/sets
- GET /api/exercises

## Explicit Deferrals

- Workout templates - future scope

## Open Questions

- Q-WORKOUT-001: How to handle concurrent edits? (Medium/risk)

## Raw PRD Traceability

docs/product/prd.md Sections 10.1-10.8

## Verified PRD Traceability

docs/product-verified/user-flows.md