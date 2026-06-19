# Dashboard

## Source Evidence

PRD §9.

## User Problem

See a quick summary of the current week's activity without navigating through multiple sections.

## Scope

In MVP. Minimal set of blocks, no customization.

## Behavior

- Current date displayed
- Last recorded body weight
- Training days count for current week
- Cardio sessions count for current week
- Current goal from user profile
- Upcoming weekly check-in reminder
- Quick actions: add workout, add cardio, add weight, weekly check-in, generate AI report

## Edge Cases

EDGE-010 (empty dashboard on first launch — no data exists).

## Acceptance Criteria

AC-103 through AC-108.

## Derived Requirements

None beyond source evidence.

## Dependencies

Workout diary, cardio, body weight, user profile.

## Open Questions

Q-FEAT-001: "Training days this week" — calendar week or trailing 7 days?
Q-FEAT-002: Weekly check-in reminder trigger.
Q-AC-03: "Last body weight" definition.