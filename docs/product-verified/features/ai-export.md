# AI Export

## Source Evidence

PRD §17, §26.9.

## User Problem

Export structured training, body, and nutrition data for AI-powered weekly analysis and recommendations.

## Scope

In MVP. Manual copy-paste to external AI (no OpenAI API integration).

## Behavior

- Configurable date range (default: last 4 weeks)
- Selectable data sections: workouts, exercises, sets, working weights, comments, RPE/RIR, cardio, body weight, measurements, photos (opt-in), nutrition, goal, additional context
- Output: ZIP archive with manifest.json, data.json, summary.md, CSVs, photos/
- Photos excluded by default
- Generated only on user request

## Derived Requirements

| Requirement | Source | Rationale | Confidence |
| --- | --- | --- | --- |
| ZIP filename includes date | §17.4 | "atlas-ai-export-YYYY-MM-DD.zip" | High |
| manifest.json includes schema version, app version | §17.5 | Required for future compatibility | High |

## Edge Cases

EDGE-008 (empty data in period), EDGE-024 (disk full during generation), EDGE-010 (max export size).

## Acceptance Criteria

AC-074 through AC-083.

## Dependencies

All data entities (workouts, cardio, body, nutrition, profile) must be implemented.

## Open Questions

Q-FEAT-010: Max AI export ZIP size and large export handling.
Q-FEAT-011: CSV encoding and escaping rules.
Q-AC-07: e1RM formula.
Q-SCOPE-006: Target AI platforms (ChatGPT only or Claude, Gemini, local LLMs).