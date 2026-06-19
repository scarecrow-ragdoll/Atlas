# PAGE-008: Charts

## Status

user-approved

## Page Purpose

Visualize progress: workout weights, body measurements, nutrition macros.

## What Is On This Page

- Exercise selector dropdown
- Date range picker
- Chart type selector
- Multiple chart views:
  - Working weight trend
  - Best set progression
  - Estimated 1RM
  - Volume trend
  - Body weight trend
  - Body measurements
  - Nutrition macros (weekly average)

## Functional Parts

- Chart container
- Period filter
- Entity selectors
- Multiple chart types (line, bar)

## Empty States

- No data for period - "No data to display"

## Loading And Error States

- Loading chart data - spinner
- Data error - error message

## Backend Dependencies

- GET /api/workouts?period=... for exercise charts
- GET /api/body-weight?period=...
- GET /api/measurements?period=...
- GET /api/nutrition-summary?period=...

## Explicit Deferrals

- Advanced chart types - future scope

## Open Questions

- Q-CHART-001: Exact 1RM formula needed

## Raw PRD Traceability

docs/product/prd.md Section 16

## Verified PRD Traceability

docs/product-verified/functional-spec.md#Charts