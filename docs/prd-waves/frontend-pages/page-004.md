# PAGE-004: Cardio

## Status

user-approved

## Page Purpose

Add/edit cardio entries with type, duration, pulse, heart rate zone.

## What Is On This Page

- Cardio list by date
- Add cardio form
- Cardio type selector (walk, run, bike, elliptical, treadmill, other)
- Duration input (minutes)
- Pulse input
- Heart rate zone selector (Zone 1-5, unknown)

## Functional Parts

- List view with edit/delete
- Create/edit form
- Type dropdown
- Zone radio buttons
- Date picker

## Empty States

- No cardio entries - "Add your first cardio session"

## Loading And Error States

- Loading list - skeleton
- Save error - toast

## Backend Dependencies

- GET /api/cardio
- POST /api/cardio
- PUT /api/cardio/{id}
- DELETE /api/cardio/{id}

## Explicit Deferrals

- None

## Open Questions

- Q-CARDIO-001: Can cardio be linked to workout day or standalone? Both supported per PRD.

## Raw PRD Traceability

docs/product/prd.md Section 12

## Verified PRD Traceability

docs/product-verified/domain-model.md#CardioEntry