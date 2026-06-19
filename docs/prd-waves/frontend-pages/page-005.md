# PAGE-005: Body Measurements

## Status

user-approved

## Page Purpose

Weekly check-in with weight, body fat %, body measurements, and photos; plus standalone weight entries.

## What Is On This Page

- Check-ins list/history
- Create check-in button
- Weight input
- Body fat % optional
- Body measurements form (neck, shoulders, bicep, chest, waist, abdomen, hips, thigh, calf)
- Left/right side toggles for paired measurements
- Photo upload (2-4 photos)
- Comments

## Functional Parts

- Check-in history list
- Measurement form with paired fields
- Photo upload (multiple, angles)
- Save check-in

## Empty States

- No check-ins - "Complete your first weekly check-in"

## Loading And Error States

- Loading check-in data - skeleton
- Photo upload error - toast

## Backend Dependencies

- GET /api/body-check-ins
- POST /api/body-check-ins
- POST /api/measurements
- POST /api/progress-photos
- GET /api/body-weight

## Explicit Deferrals

- None

## Open Questions

- None blocking

## Raw PRD Traceability

docs/product/prd.md Section 13

## Verified PRD Traceability

docs/product-verified/domain-model.md#BodyCheckIn, #BodyMeasurement, #ProgressPhoto