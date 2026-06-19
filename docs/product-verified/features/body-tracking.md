# Body Tracking

## Source Evidence

PRD §13, §26.5, §26.6.

## User Problem

Track body weight and measurements over time for progress analysis.

## Scope

In MVP. Weekly check-ins plus standalone weight entries.

## Behavior

- Weekly check-in: date, weight (optional), body fat % (optional), 10 measurements, 2-4 photos, comment
- Standalone weight entry: date, weight, comment (optional)
- 10 measurements: neck, shoulders, forearms, biceps, chest, waist, abdomen, hips, thigh, calves
- Paired measurements: left/right values or single common value

## Acceptance Criteria

AC-014 through AC-016, AC-052 through AC-057.

## Derived Requirements

None beyond source evidence.

## Edge Cases

EDGE-006: Check-in with 0-1 photos.
EDGE-007: Body measurement value 0 or negative.

## Dependencies

Progress photos.

## Open Questions

Q-ACTOR-03: Check-in without 2-4 photos — required or recommended?