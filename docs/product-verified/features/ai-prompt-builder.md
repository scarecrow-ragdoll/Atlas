# AI Prompt Builder

## Source Evidence

PRD §18.

## User Problem

Generate a ready-to-send AI prompt that includes all relevant context, goal, and week-specific flags for meaningful analysis.

## Scope

In MVP. Manual copy-paste, no API integration.

## Behavior

- Persistent AI context stored and reused: goal, height, optional age, training experience, split, limitations, progression style, nutrition strategy, persistent comment
- User can update context at any time
- One-time comment per export
- Week flags: poor sleep, high stress, illness, injury/pain, AAS/cycle, calorie deficit, surplus, maintenance, missed workouts, travel/disrupted routine
- Generated prompt asks AI for: progress analysis, working weight dynamics, set vs actual weight comparison, volume evaluation, RPE/RIR consideration, cardio context, weight/measurement correlation, nutrition review, next-week recommendations
- Prompt designed for ChatGPT (compatible with other AI models)

## Acceptance Criteria

AC-023, AC-084 through AC-089.

## Derived Requirements

None beyond source evidence.

## Edge Cases

No feature-specific edge cases beyond general edge cases documented in edge-cases.md.

## Open Questions

Q-SCOPE-006: Target AI platforms.
Q-ACTOR-23: Goal/context setup before first AI export.

## Dependencies

User profile (goal, context), AI export ZIP generation.