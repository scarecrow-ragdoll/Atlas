# Cardio

## Source Evidence

PRD §12, §26.4.

## User Problem

Log cardio sessions for recovery and load context in AI analysis.

## Scope

In MVP. Simple entries with type, duration, pulse/zone.

## Behavior

- Cardio entry by date (standalone or day-attached — ambiguity per Q-SCOPE-005)
- Field: date, cardio type (enum: walking/running/bike/elliptical/treadmill/other), duration (minutes), avg pulse (optional), heart rate zone (1-5/unknown, optional), notes (optional)

## Edge Cases

EDGE-006 (cardio without pulse/zone — allowed).

## Acceptance Criteria

AC-012, AC-013, AC-050, AC-051.

## Derived Requirements

None beyond source evidence.

## Dependencies

Date navigation.

## Open Questions

Q-SCOPE-005: Cardio — standalone vs workout day entity.
Q-FEAT-006: Custom cardio types beyond basic enum.