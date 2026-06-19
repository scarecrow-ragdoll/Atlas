# PAGE-001: Dashboard

## Status

user-approved

## Page Purpose

Weekly summary with quick actions for today's date. Entry point to weekly check-in cycle.

## What Is On This Page

- Current date
- Last body weight
- Workout days count this week
- Cardio count this week
- Current goal
- Next check-in reminder
- Quick actions: Add workout, Add cardio, Add weight, Open check-in, Generate AI report

## Functional Parts

- Date display
- Weight summary card
- Training summary card
- Cardio summary card
- Goal display
- Check-in reminder badge
- Action buttons grid

## Empty States

- No workouts this week - "Add your first workout"
- No measurements - "Complete your first check-in"

## Loading And Error States

- Loading summary data - skeleton cards
- Data error - error message with retry

## Backend Dependencies

- GET /api/workouts?week=current
- GET /api/body-weight/latest
- GET /api/user-profile
- GET /api/check-ins?upcoming=true

## Explicit Deferrals

- Visual design not specified

## Open Questions

- None blocking

## Raw PRD Traceability

docs/product/prd.md Section 9

## Verified PRD Traceability

docs/product-verified/functional-spec.md