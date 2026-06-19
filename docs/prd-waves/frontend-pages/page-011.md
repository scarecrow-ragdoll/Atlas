# PAGE-011: Settings

## Status

user-approved

## Page Purpose

Application settings: PIN configuration, AI context, export preferences.

## What Is On This Page

- PIN enabled toggle
- PIN change form (requires current PIN)
- AI context form:
  - Goal
  - Height
  - Age (optional)
  - Training experience
  - Current split (optional)
  - Limits/injuries
  - Progression style
  - Nutrition strategy
- Default export weeks setting

## Functional Parts

- PIN toggle and management
- AI context form
- Preferences form

## Empty States

- No settings configured yet - defaults apply
- PIN not set

## Loading And Error States

- Saving settings - spinner
- PIN validation error - toast

## Backend Dependencies

- GET /api/settings
- PUT /api/settings
- POST /api/settings/pin-enable
- POST /api/settings/pin-disable
- POST /api/settings/pin-change

## Explicit Deferrals

- User profile editing - combined with settings

## Open Questions

- None blocking

## Raw PRD Traceability

docs/product/prd.md Sections 7.2, 18.2

## Verified PRD Traceability

docs/product-verified/actors-and-permissions.md