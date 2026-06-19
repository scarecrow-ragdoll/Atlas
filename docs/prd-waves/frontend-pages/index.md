# Frontend Pages

## Status

user-approved

## Scope Source

docs/product/prd.md Sections 9-20

## Page Order

| Page ID | Name | Purpose |
| --- | --- | --- |
| PAGE-001 | Dashboard | Weekly summary with quick actions |
| PAGE-002 | Workout Diary | Daily workout entry by date |
| PAGE-003 | Exercise Library | CRUD exercises with media |
| PAGE-004 | Cardio | Add/edit cardio entries |
| PAGE-005 | Body Measurements | Weekly check-in and weight entries |
| PAGE-006 | Progress Photos | Photo management within check-in |
| PAGE-007 | Nutrition | Products, template, daily overrides |
| PAGE-008 | Charts | Progress visualization |
| PAGE-009 | AI Export | Prompt builder and ZIP download |
| PAGE-010 | Import/Export | Full backup operations |
| PAGE-011 | Settings | PIN, AI context, preferences |

## Raw PRD Source Coverage

- Sections 9-20 cover all 11 pages
- Each page purpose traced to PRD sections

## Verified PRD Source Coverage

- Product capabilities mapped to pages

## Shared UX States

- Loading: Fetching data from API
- Error: API/data errors
- Empty: No data yet (first use)
- Auth: PIN entry screen (overlay or route)

## Backend Dependencies By Page

- PAGE-001: workouts, body-weight, user-profile, exercises
- PAGE-002: workouts, exercises, cardio
- PAGE-003: exercises, exercise-media
- PAGE-004: cardio
- PAGE-005: body-check-ins, measurements, photos
- PAGE-006: photos (embedded in check-in)
- PAGE-007: nutrition-products, nutrition-templates, daily-overrides
- PAGE-008: all entities for charting
- PAGE-009: all entities for export
- PAGE-010: all entities for backup
- PAGE-011: settings, user-profile

## Explicit Frontend Deferrals

- No visual design specs
- No mobile screens

## Open Questions

- Q-PAGE-001: How should PIN auth be handled? Modal or separate route?

## Traceability

Each page traces to docs/product/prd.md sections.