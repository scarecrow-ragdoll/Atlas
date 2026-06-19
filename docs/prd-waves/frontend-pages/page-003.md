# PAGE-003: Exercise Library

## Status

user-approved

## Page Purpose

CRUD exercises with working weight and optional media. Foundation for workout diary.

## What Is On This Page

- Exercise list with search
- Add exercise form
- Edit exercise form
- Media upload per exercise
- Working weight field

## Functional Parts

- Exercise list/table
- Search filter
- Create/edit modal/form
- Media upload component
- Delete with confirmation

## Empty States

- No exercises - "Create your first exercise"

## Loading And Error States

- Loading list - skeleton rows
- Media upload error - toast

## Backend Dependencies

- GET /api/exercises
- POST /api/exercises
- PUT /api/exercises/{id}
- DELETE /api/exercises/{id}
- POST /api/exercise-media
- DELETE /api/exercise-media/{id}

## Explicit Deferrals

- Pre-built exercise catalog - future scope
- Exercise sharing - future scope

## Open Questions

- None blocking

## Raw PRD Traceability

docs/product/prd.md Section 11

## Verified PRD Traceability

docs/product-verified/domain-model.md#Exercise