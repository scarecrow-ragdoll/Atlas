# PAGE-009: AI Export

## Status

user-approved

## Page Purpose

Generate AI prompt and export files for AI analysis. Period selection, section toggles, ZIP download.

## What Is On This Page

- Date range picker (default 4 weeks)
- Section checkboxes (workouts, exercises, sets, RPE/RIR, cardio, weight, measurements, photos, nutrition)
- One-time comment box
- User goal display
- Week flags selector
- Generate button
- Prompt display/copy
- ZIP download

## Functional Parts

- Period selector
- Section selector
- Prompt preview
- ZIP generation and download

## Empty States

- No data in period - warning message

## Loading And Error States

- Generating export - spinner
- Error - toast with details

## Backend Dependencies

- POST /api/ai-export (generate)
- GET /api/ai-export/download
- GET /api/user-profile (goal context)
- GET /api/week-flags

## Explicit Deferrals

- Direct OpenAI API call - future scope

## Open Questions

- None blocking

## Raw PRD Traceability

docs/product/prd.md Sections 17, 18

## Verified PRD Traceability

docs/product-verified/functional-spec.md#AI Export