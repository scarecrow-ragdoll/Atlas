<!-- FILE: .agents/skills/detail-prd-wave/references/output-contract.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Define the detailed backend PRD wave output structure produced by detail-prd-wave. -->
<!--   SCOPE: Covers required files, headings, status values, stable ids, reviewer verdicts, question ledgers, backend-only and ready-for-dev gates; excludes subagent prompt templates and frontend planning. -->
<!--   DEPENDS: .agents/skills/detail-prd-wave/SKILL.md. -->
<!--   LINKS: M-PRD-WAVE-DETAILER / V-M-PRD-WAVE-DETAILER. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Required Structure - Lists files that must exist in docs/prd-wave-details. -->
<!--   Required Headings - Lists minimum headings for package and selected wave files. -->
<!--   Status Model - Defines package and wave statuses. -->
<!--   Ready-For-Dev Gate - Defines criteria for user-approvable detailed waves. -->
<!--   Ledger Formats - Defines question and reviewer table shapes. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.2 - Clarified detailed waves are backend-only and frontend-pages are dependency context only. -->
<!-- END_CHANGE_SUMMARY -->

# Detailed Backend PRD Wave Output Contract

Create this structure as selected waves are detailed:

```text
docs/prd-wave-details/
  index.md
  source-inventory.md
  wave-map-context.md
  codebase-fit.md
  open-questions.md
  waves/
    index.md
    wave-01.md
  appendix/
    reviewer-verdicts.md
    traceability.md
    question-ledger.md
    decision-log.md
    run-history.md
```

Only create or update the selected `waves/wave-<nn>.md` in a run. Later waves remain absent or shallow references until their own run.

## Required File Purposes

- `index.md`: package status, selected wave, source wave gate, current wave gate, source set, and next action.
- `source-inventory.md`: backend wave, frontend-pages context, product, technical, GRACE, codebase, source delta, and gap inventories.
- `wave-map-context.md`: selected backend wave boundary, prior-backend-wave fit, future-backend-wave fit, frontend dependency context, dependency order, and scope-collision checks.
- `codebase-fit.md`: read-only codebase touchpoints, current module contracts, likely graph deltas, and unsupported assumptions.
- `open-questions.md`: unresolved wave-blocking, owner-decision, deferred, and watchlist questions.
- `waves/index.md`: detailed wave list, status, dependency order, and user approval state.
- `waves/wave-<nn>.md`: one detailed selected wave brief.
- `appendix/reviewer-verdicts.md`: canonical reviewer verdict ledger across detailed waves.
- `appendix/traceability.md`: map slices, ACs, ECs, tests, questions, code touchpoints, and decisions to sources.
- `appendix/question-ledger.md`: canonical question ledger across detailed waves.
- `appendix/decision-log.md`: source wave gates, scope decisions, user approvals, and rejected assumptions.
- `appendix/run-history.md`: runs, selected wave, source deltas, planner cycles, review cycles, and approvals.

## Required Headings

### index.md

- `# Detailed Backend PRD Waves`
- `## Status`
- `## Selected Wave`
- `## Source Wave Gate`
- `## Current Wave Gate`
- `## Source Set`
- `## Next Action`

### source-inventory.md

- `# Source Inventory`
- `## PRD Wave Sources`
- `## Frontend Pages Source`
- `## Product Sources`
- `## Technical Sources`
- `## GRACE Sources`
- `## Codebase Sources`
- `## Source Delta`
- `## Source Gaps`

### wave-map-context.md

- `# Wave Map Context`
- `## Selected Backend Wave Boundary`
- `## Prior Backend Wave Fit`
- `## Future Backend Wave Fit`
- `## Frontend Pages Context`
- `## Dependency Order`
- `## Scope Collision Check`

### codebase-fit.md

- `# Codebase Fit`
- `## Relevant Modules`
- `## Relevant Files Read`
- `## Public Contracts`
- `## Generated Artifact Impact`
- `## Integration Points`
- `## Likely Graph Deltas`
- `## Unsupported Assumptions`

### open-questions.md

- `# Open Questions`
- `## Wave-Blocking`
- `## Needs Owner Decision`
- `## Deferred`
- `## Watchlist`
- `## Resolved This Run`

### waves/index.md

- `# Detailed Backend Waves`
- `## Wave List`
- `## Dependency Order`
- `## Approval State`

### waves/wave-<nn>.md

- `# Wave <nn>: <Name>`
- `## Status`
- `## User Approval`
- `## Source Wave Summary`
- `## Outcome After Implementation`
- `## Scope Included`
- `## Scope Excluded`
- `## Dependencies And Other-Wave Fit`
- `## Frontend Pages Dependencies`
- `## Codebase Fit And Touchpoints`
- `## Design Contracts`
- `## Data API Integration And Operations`
- `## Security Privacy And Compliance`
- `## Implementation Slices`
- `## Acceptance Criteria`
- `## Exit Criteria`
- `## Verification Obligations`
- `## Rollout Rollback And Compatibility`
- `## Handoff Packets`
- `## Reviewer Verdicts`
- `## Open Questions`
- `## Traceability`

### appendix/reviewer-verdicts.md

- `# Reviewer Verdicts`
- `## Current Wave`
- `## Historical Waves`
- `## Final Fit Reviews`
- `## Rejected Findings`

### appendix/traceability.md

- `# Traceability`
- `## Slice Map`
- `## Acceptance Criteria Map`
- `## Exit Criteria Map`
- `## Verification Obligation Map`
- `## Code Touchpoint Map`
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
- `## Source Wave Gate`
- `## User Wave Approvals`
- `## Scope Decisions`
- `## Codebase Fit Decisions`
- `## Deferrals`
- `## Rejected Assumptions`

### appendix/run-history.md

- `# Run History`
- `## Runs`
- `## Selected Wave History`
- `## Planner Cycles`
- `## Review Cycles`
- `## Source Delta History`
- `## Approval Gate History`

## Status Model

Use exactly one package status in `index.md`:

- `draft`: selected backend wave detailing has started but no current wave is ready for dev.
- `questions-open`: wave-blocking or owner-decision questions remain.
- `blocked`: source, codebase-fit, reviewer, validation, or budget blockers prevent progress.
- `wave-ready-awaiting-user-approval`: selected backend wave is `ready-for-dev` and awaits explicit user approval.
- `wave-approved`: selected backend wave is user-approved.
- `superseded`: package was replaced by a later run or re-scope.

Use exactly one wave status in each detailed wave file:

- `draft`: wave details are being prepared.
- `needs-revision`: reviewers requested changes.
- `questions-open`: open wave-blocking or owner-decision questions remain.
- `blocked`: missing source, codebase-fit, unavailable reviewers, or exhausted budgets prevent ready-for-dev.
- `ready-for-dev`: reviewers approved and no open blockers remain, but user approval is still pending.
- `user-approved`: the user explicitly approved this detailed wave.
- `superseded`: wave was replaced by a later approved re-scope.

## Stable IDs

- Waves: `WAVE-01`, `WAVE-02`, ...
- Backend implementation slices: `SLICE-W01-001`, `SLICE-W01-002`, ...
- Acceptance criteria: `AC-W01-001`, `AC-W01-002`, ...
- Exit criteria: `EC-W01-001`, `EC-W01-002`, ...
- Verification obligations: `TEST-W01-001`, `TEST-W01-002`, ...
- Handoff packets: `HANDOFF-W01-001`, `HANDOFF-W01-002`, ...
- Questions: `DQ-W01-001`, `DQ-W01-002`, ...
- Decisions: `DDEC-W01-001`, `DDEC-W01-002`, ...

## Required Reviewer Perspectives

Each ready-for-dev or user-approved wave must have approved verdicts from:

- `product-scope-and-ac`
- `architecture-codebase-fit`
- `data-api-integration-ops`
- `security-privacy-compliance`
- `testing-exit-criteria`
- `sequencing-other-wave-fit`
- `traceability-consistency`
- `final-wave-fit-review`

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

Open `wave-blocking` or `needs-owner-decision` rows block `ready-for-dev`, block user approval, and block downstream Beads, Jira, GRACE execution, or implementation.

## Source Wave Gate Format

`index.md` `## Source Wave Gate` must include one of these markers:

- `source-wave-gate: passed`
- `source-wave-gate: blocked`

Ready-for-dev or user-approved packages require `source-wave-gate: passed` plus the selected source wave id and source wave path. Blocked or questions-open packages require `source-wave-gate: blocked` or an equivalent blocker statement plus a matching open question row in `open-questions.md`, `appendix/question-ledger.md`, or the selected wave file.

## Ready-For-Dev Gate

A wave may be `ready-for-dev` only when:

- source wave gate names the selected source wave and confirms no open source-wave blockers;
- all required reviewer perspectives have `approved`;
- final-wave-fit review names the reviewed candidate package path and approves it;
- the selected source wave is a backend wave and any frontend-pages references are dependency context only;
- the wave file contains at least one `SLICE-W<nn>-...` implementation slice;
- the wave file contains at least one `AC-W<nn>-...` acceptance criterion;
- the wave file contains at least one `EC-W<nn>-...` exit criterion;
- the wave file contains at least one `TEST-W<nn>-...` verification obligation;
- codebase-fit evidence names relevant modules or explicitly records why none are needed;
- other-wave fit evidence records prior and future backend wave compatibility;
- frontend-pages evidence records only backend dependencies or explicitly says no frontend dependency;
- aggregate and wave-local question ledgers have no open `wave-blocking` or `needs-owner-decision` rows for the selected wave;
- source evidence and traceability point to `docs/prd-waves`, verified docs, GRACE XML, source deltas, code paths, reviewer reports, or explicit decisions.

## User Approval Gate

`user-approved` requires:

- the wave was already `ready-for-dev`;
- the user approved the wave after receiving the overview;
- `## User Approval` includes an `approved-by-user` entry with date or conversation reference;
- `appendix/decision-log.md` records the approval.
