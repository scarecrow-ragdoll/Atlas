# Charts

## Source Evidence

PRD §16.

## User Problem

Visualize progress trends in training, body measurements, and nutrition over configurable periods.

## Scope

In MVP. Simple charts with period selection and data filtering.

## Behavior

- Training charts per exercise: working weight, best set, e1RM, volume, total reps, working sets count
- Body charts: weight, body fat %, individual measurements, multiple measurements overlay
- Nutrition charts: weekly averages for calories, protein, fat, carbs
- Configurable period with filters

## Acceptance Criteria

AC-020 through AC-022, AC-065 through AC-073.

## Derived Requirements

None beyond source evidence.

## Edge Cases

EDGE-008: Chart with no data for selected period/exercise.

## Dependencies

All data entities.

## Open Questions

Q-FEAT-009: Charting library and interactivity requirements.