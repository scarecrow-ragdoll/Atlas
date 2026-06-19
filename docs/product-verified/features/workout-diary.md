# Workout Diary

## Source Evidence

PRD §10, §26.2, §26.3.

## User Problem

Log workout data (exercises, sets, reps, weights) by date with minimal friction.

## Scope

In MVP. No workout templates, no plan scheduling, no split organization.

## Behavior

- Default view: today's date
- Calendar navigation to any past date
- Backdating allowed (create/edit data for past dates)
- One workout record per date
- Per date: list of exercises, cardio entries (same-day), body weight, general notes
- Per exercise: link to library exercise, display order, working weight snapshot, notes, sets
- Per set: number, weight, reps, optional RPE, optional RIR, optional comment
- Working weight auto-populated on exercise add
- Progression display: current working weight, actual weights, best set, volume, e1RM, weekly trend
- Progression signals: stable top-end reps, weight/volume growth, stagnation, regression, negative comments
- No automatic working weight change without user confirmation

## Derived Requirements

| Requirement | Source | Rationale | Confidence |
| --- | --- | --- | --- |
| Working weight snapshot stored per exercise in workout day | §10.6 | "в тренировочном дне сохраняется snapshot рабочего веса на момент выполнения" | High |
| Working weight auto-populated when adding exercise | §6.1, §10.6 | "рабочий вес упражнения должен подставляться автоматически" | High |
| Copy previous set value within exercise | §21 | "копирование значения прошлого подхода внутри упражнения" | High |

## Edge Cases

EDGE-001 (0 weight/reps), EDGE-004 (backdating before instance creation), EDGE-005 (empty workout day), EDGE-016 (concurrent tab access), EDGE-022 (DB connection lost during save).

## Acceptance Criteria

AC-005 through AC-011, AC-035 through AC-042.

## Dependencies

Exercise library must exist first. Calendar UI component.

## Open Questions

Q-FEAT-001: Dashboard "training days this week" — calendar week or trailing 7 days?
Q-AC-06: Working weight auto-populate UI behavior.
Q-AC-08: Best set definition (heaviest weight vs highest e1RM).