# PAGE-006: Progress Photos

## Status

user-approved

## Page Purpose

Manage progress photos within check-ins. View by date/check-in.

## What Is On This Page

- Photos grouped by check-in
- Photo viewer (full size)
- Photo labels and angles
- Delete photo

## Functional Parts

- Gallery/grid view
- Photo detail view
- Angle badges (front, side, back, custom)
- Label display

## Empty States

- No photos - shown in PAGE-005 check-in form

## Loading And Error States

- Loading photos - skeleton grid

## Backend Dependencies

- GET /api/progress-photos?checkin={id}
- DELETE /api/progress-photos/{id}

## Explicit Deferrals

- Photo taken within app - not supported (upload only)

## Open Questions

- None blocking

## Raw PRD Traceability

docs/product/prd.md Section 14

## Verified PRD Traceability

docs/product-verified/domain-model.md#ProgressPhoto