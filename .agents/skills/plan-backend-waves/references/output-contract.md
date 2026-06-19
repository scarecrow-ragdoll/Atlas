<!-- FILE: .agents/skills/plan-backend-waves/references/output-contract.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Define the backend wave planning output structure produced by plan-backend-waves. -->
<!--   SCOPE: Covers required files, headings, status values, stable ids, reviewer verdicts, question ledgers, and wave/user approval gates; excludes subagent prompt templates. -->
<!--   DEPENDS: .agents/skills/plan-backend-waves/SKILL.md. -->
<!--   LINKS: M-BACKEND-WAVE-PLANNER / V-M-BACKEND-WAVE-PLANNER. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Required Structure - Lists files that must exist in docs/backend-waves. -->
<!--   Required Headings - Lists minimum headings for package and wave files. -->
<!--   Status Model - Defines package and wave statuses. -->
<!--   Ready-For-Dev Gate - Defines criteria for a wave to be approvable by the user. -->
<!--   Ledger Formats - Defines question and reviewer table shapes. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.0 - Added backend wave output contract. -->
<!-- END_CHANGE_SUMMARY -->

# Backend Waves Output Contract

Create this structure as the package grows:

```text
docs/backend-waves/
  index.md
  source-inventory.md
  wave-map.md
  open-questions.md
  waves/
    index.md
    wave-01.md
    wave-02.md
  appendix/
    reviewer-verdicts.md
    traceability.md
    question-ledger.md
    decision-log.md
    run-history.md
```

`waves/wave-02.md` and later files are illustrative. Create them only after the previous wave has user approval.

## Required File Purposes

- `index.md`: package status, technical approval gate, current wave gate, source set, and next action.
- `source-inventory.md`: technical and product sources, prior wave artifacts, source deltas, answered questions, and coverage gaps.
- `wave-map.md`: complete but shallow backend wave map with tentative future waves.
- `open-questions.md`: unresolved wave-blocking, owner-decision, deferred, and watchlist questions.
- `waves/index.md`: wave list, status, dependency order, and user approval state.
- `waves/wave-<nn>.md`: one detailed current or approved wave.
- `appendix/reviewer-verdicts.md`: canonical reviewer verdict ledger across all waves.
- `appendix/traceability.md`: map wave tasks, criteria, tests, and questions to technical/product sources and decisions.
- `appendix/question-ledger.md`: canonical question ledger across waves.
- `appendix/decision-log.md`: technical approval gate decisions, user approvals, scope deferrals, and rejected assumptions.
- `appendix/run-history.md`: runs, current wave, source deltas, reviewer cycles, and approval history.

## Required Headings

### index.md

- `# Backend Waves`
- `## Status`
- `## Technical Approval Gate`
- `## Current Wave Gate`
- `## Source Set`
- `## Next Action`

### source-inventory.md

- `# Source Inventory`
- `## Technical Sources`
- `## Product Sources`
- `## Prior Wave Sources`
- `## Source Delta`
- `## Coverage Gaps`

### wave-map.md

- `# Wave Map`
- `## Backend Scope Inventory`
- `## Tentative Wave Count`
- `## Sequential Wave Map`
- `## MVP Scope Check`
- `## Dependency Notes`

### open-questions.md

- `# Open Questions`
- `## Wave-Blocking`
- `## Needs Owner Decision`
- `## Deferred`
- `## Watchlist`
- `## Resolved This Run`

### waves/index.md

- `# Waves`
- `## Wave List`
- `## Dependency Order`
- `## Approval State`

### waves/wave-<nn>.md

- `# Wave <nn>: <Name>`
- `## Status`
- `## User Approval`
- `## Outcome After Implementation`
- `## Source Evidence`
- `## Scope Included`
- `## Scope Excluded`
- `## Dependencies`
- `## Backend Design`
- `## Data And Migration Work`
- `## API Jobs And Events`
- `## Auth Security And Compliance`
- `## Operations Observability`
- `## Implementation Tasks`
- `## Acceptance Criteria`
- `## Exit Criteria`
- `## Verification Plan`
- `## Rollback And Compatibility`
- `## Jira Ready Tasks`
- `## Reviewer Verdicts`
- `## Open Questions`
- `## Traceability`

### appendix/reviewer-verdicts.md

- `# Reviewer Verdicts`
- `## Current Wave`
- `## Historical Waves`
- `## Rejected Findings`

### appendix/traceability.md

- `# Traceability`
- `## Wave Task Map`
- `## Acceptance Criteria Map`
- `## Exit Criteria Map`
- `## Test Obligation Map`
- `## Question Map`
- `## Source Map`

### appendix/question-ledger.md

- `# Question Ledger`
- `## Open Questions`
- `## Answered Questions`
- `## Follow-Up Questions`
- `## Resolved Questions`
- `## Deferred Questions`

### appendix/decision-log.md

- `# Decision Log`
- `## Technical Approval Gate`
- `## User Wave Approvals`
- `## Scope Decisions`
- `## Deferrals`
- `## Rejected Assumptions`

### appendix/run-history.md

- `# Run History`
- `## Runs`
- `## Wave Planning Cycles`
- `## Source Delta History`
- `## Approval Gate History`

## Status Model

Use exactly one package status in `index.md`:

- `draft`: wave planning has started but no current wave is ready for dev.
- `blocked`: technical approval, source, reviewer, validation, or question blockers prevent progress.
- `wave-ready-awaiting-user-approval`: current wave is `ready-for-dev` and awaits explicit user approval.
- `wave-approved-planning-next`: latest wave is user-approved and the next wave may be planned.
- `waves-verified`: all required backend waves are planned, reviewer-approved, user-approved, and traceable.
- `superseded`: package was replaced by a later run or re-scope.

Use exactly one wave status in each wave file:

- `draft`: wave details are being prepared.
- `needs-revision`: reviewers requested changes.
- `questions-open`: open wave-blocking or owner-decision questions remain.
- `blocked`: missing source, unavailable reviewers, or exhausted budgets prevent ready-for-dev.
- `ready-for-dev`: reviewers approved and no open blockers remain, but user approval is still pending.
- `user-approved`: the user explicitly approved this wave after reviewing the ready-for-dev overview.
- `superseded`: wave was replaced by a later approved re-scope.

## Stable IDs

- Waves: `WAVE-01`, `WAVE-02`, ...
- Backend tasks: `BTASK-W01-001`, `BTASK-W01-002`, ...
- Acceptance criteria: `AC-W01-001`, `AC-W01-002`, ...
- Exit criteria: `EXIT-W01-001`, `EXIT-W01-002`, ...
- Test obligations: `BTEST-W01-001`, `BTEST-W01-002`, ...
- Questions: `BQ-W01-001`, `BQ-W01-002`, ...
- Decisions: `BDEC-W01-001`, `BDEC-W01-002`, ...

## Required Reviewer Perspectives

Each ready-for-dev or user-approved wave must have approved verdicts from:

- `backend-architecture`
- `data-api-contract`
- `security-integration`
- `testing-delivery`
- `sequencing-mvp`
- `traceability-consistency`

## Reviewer Verdict Table

Use this shape in each wave and in `appendix/reviewer-verdicts.md`:

```text
| Wave | Perspective | Attempt | Verdict | Reviewer Report | Required Revisions | Notes |
| --- | --- | --- | --- | --- | --- | --- |
```

Allowed final verdicts: `approved`, `needs-revision`, `blocked`. Staging scaffolds may use `pending-review` only while validation is run with `--allow-placeholders`.

## Question Ledger Format

Use this table shape in every wave-local and aggregate ledger:

```text
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
```

Allowed severities: `wave-blocking`, `needs-owner-decision`, `deferred`, `watchlist`.

Allowed statuses: `open`, `answered`, `resolved`, `deferred`, `superseded`.

Open `wave-blocking` or `needs-owner-decision` rows block `ready-for-dev`, block user approval, and block planning the next wave.

## Ready-For-Dev Gate

A wave may be `ready-for-dev` only when:

- `index.md` `## Technical Approval Gate` records the approved technical source path, run or decision reference, and literal status `approved-to-dev`;
- all required reviewer perspectives have `approved`;
- the wave file contains at least one `AC-W<nn>-...` acceptance criterion;
- the wave file contains at least one `EXIT-W<nn>-...` exit criterion;
- the wave file contains implementation tasks with stable `BTASK-W<nn>-...` ids;
- all test obligations needed by the changed backend surface are explicit;
- the aggregate question ledger has no open `wave-blocking` or `needs-owner-decision` rows for the wave;
- source evidence and traceability point to `docs/technical-verified`, `docs/product-verified`, source deltas, or explicit decisions.

## User Approval Gate

`user-approved` requires:

- the wave was already `ready-for-dev`;
- the user approved the wave after receiving the overview;
- `## User Approval` includes an `approved-by-user` entry with date or conversation reference;
- `appendix/decision-log.md` records the approval as a decision;
- no open wave-blocking or owner-decision questions exist.
