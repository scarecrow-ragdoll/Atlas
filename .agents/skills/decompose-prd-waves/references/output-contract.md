<!-- FILE: .agents/skills/decompose-prd-waves/references/output-contract.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Define the backend-only wave and per-page frontend output structure produced by decompose-prd-waves. -->
<!--   SCOPE: Covers required files, headings, status values, stable ids, reviewer verdicts, question ledgers, backend-only wave gates, frontend page gates, raw/verified PRD traceability, and shallow-only gates; excludes subagent prompt templates. -->
<!--   DEPENDS: .agents/skills/decompose-prd-waves/SKILL.md. -->
<!--   LINKS: M-PRD-WAVE-DECOMPOSER / V-M-PRD-WAVE-DECOMPOSER. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Required Structure - Lists files that must exist in docs/prd-waves. -->
<!--   Required Headings - Lists minimum headings for package and wave files. -->
<!--   Status Model - Defines package and wave statuses. -->
<!--   Backend Wave And Frontend Page Gate - Defines criteria for an approvable backend wave map plus per-page frontend handoff files. -->
<!--   Ledger Formats - Defines question and reviewer table shapes. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.2 - Replaced the single frontend sequence artifact with per-page frontend files sourced from raw and verified PRDs. -->
<!-- END_CHANGE_SUMMARY -->

# PRD Waves Output Contract

Create this structure:

```text
docs/prd-waves/
  index.md
  source-inventory.md
  scope-inventory.md
  wave-map.md
  frontend-pages/
    index.md
    page-001.md
    page-002.md
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

Create one `waves/wave-<nn>.md` file per top-level backend wave. All wave files remain shallow and backend-only. Do not create frontend waves. Frontend scope is represented only under `frontend-pages/`, with one `frontend-pages/page-<nnn>.md` file per page.

## Required File Purposes

- `index.md`: package status, source gate, shallow-only gate, wave count, and next action.
- `source-inventory.md`: raw product, verified product, technical, prior wave, source delta, and source gap inventory.
- `scope-inventory.md`: full PRD capability and concern inventory before grouping backend scope into waves and frontend scope into per-page frontend files.
- `wave-map.md`: complete top-level backend wave list, dependency order, risk classes, and downstream planner recommendations.
- `frontend-pages/index.md`: ordered frontend page index, source coverage from raw and verified PRDs, shared UX states, backend dependencies by page, deferrals, and frontend questions.
- `frontend-pages/page-<nnn>.md`: one page brief per frontend page, sourced from raw and verified PRDs, describing what is on the page, functional parts, empty states, loading/error states, backend dependencies, deferrals, questions, and traceability.
- `open-questions.md`: unresolved decomposition-blocking, owner-decision, deferred, and watchlist questions.
- `waves/index.md`: wave list, status, dependency order, and user approval state.
- `waves/wave-<nn>.md`: shallow wave summary.
- `appendix/reviewer-verdicts.md`: canonical scope and consistency reviewer verdict ledger.
- `appendix/traceability.md`: source-to-scope-to-wave mapping.
- `appendix/question-ledger.md`: canonical question ledger across scopes and waves.
- `appendix/decision-log.md`: source gates, scope cuts, deferrals, user approvals, and rejected assumptions.
- `appendix/run-history.md`: runs, source deltas, scope mapper cycles, consistency cycles, and approvals.

## Required Headings

### index.md

- `# PRD Waves`
- `## Status`
- `## Source Gate`
- `## Shallow Wave Gate`
- `## Wave Count`
- `## Source Set`
- `## Next Action`

### source-inventory.md

- `# Source Inventory`
- `## Raw Product Sources`
- `## Verified Product Sources`
- `## Technical Sources`
- `## Prior Wave Sources`
- `## Source Delta`
- `## Source Gaps`

### scope-inventory.md

- `# Scope Inventory`
- `## Capability Groups`
- `## User Journey Groups`
- `## Data Lifecycle Groups`
- `## Integration And Operations Groups`
- `## Client Experience Groups`
- `## Security Compliance Groups`
- `## Explicit Deferrals`

### wave-map.md

- `# Wave Map`
- `## Top-Level Wave List`
- `## Dependency Order`
- `## Coverage Matrix`
- `## More Than Eight Wave Check`
- `## Downstream Planning Recommendations`

### frontend-pages/index.md

- `# Frontend Pages`
- `## Status`
- `## Scope Source`
- `## Page Order`
- `## Raw PRD Source Coverage`
- `## Verified PRD Source Coverage`
- `## Shared UX States`
- `## Backend Dependencies By Page`
- `## Explicit Frontend Deferrals`
- `## Open Questions`
- `## Traceability`

### frontend-pages/page-<nnn>.md

- `# PAGE-<nnn>: <Name>`
- `## Status`
- `## Page Purpose`
- `## What Is On This Page`
- `## Functional Parts`
- `## Empty States`
- `## Loading And Error States`
- `## Backend Dependencies`
- `## Explicit Deferrals`
- `## Open Questions`
- `## Raw PRD Traceability`
- `## Verified PRD Traceability`

### open-questions.md

- `# Open Questions`
- `## Decomposition Blocking`
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
- `## Purpose`
- `## Outcome After Wave`
- `## Included Scope`
- `## Excluded Scope`
- `## Dependencies`
- `## Surface Categories`
- `## Risk Class`
- `## Recommended Next Planning`
- `## Open Questions`
- `## Traceability`

### appendix/reviewer-verdicts.md

- `# Reviewer Verdicts`
- `## Scope Reviews`
- `## Consistency Review`
- `## Rejected Findings`

### appendix/traceability.md

- `# Traceability`
- `## Source To Scope Map`
- `## Scope To Wave Map`
- `## Wave To Source Map`
- `## Question Map`
- `## Decision Map`

### appendix/question-ledger.md

- `# Question Ledger`
- `## Open Questions`
- `## Answered Questions`
- `## Follow-Up Questions`
- `## Resolved Questions`
- `## Deferred Questions`

### appendix/decision-log.md

- `# Decision Log`
- `## Source Gate`
- `## Scope Decisions`
- `## Deferrals`
- `## User Wave Map Approvals`
- `## Rejected Assumptions`

### appendix/run-history.md

- `# Run History`
- `## Runs`
- `## Scope Mapper Cycles`
- `## Consistency Cycles`
- `## Source Delta History`
- `## Approval Gate History`

## Status Model

Use exactly one package status in `index.md`:

- `draft`: decomposition has started but the map is not review-approved.
- `questions-open`: decomposition-blocking or owner-decision questions remain.
- `blocked`: missing sources, unavailable subagents/reviewers, or exhausted budgets prevent completion.
- `waves-ready-awaiting-user-approval`: all required reviews approved and no open blockers remain; the user has not yet approved the map.
- `waves-approved`: the user explicitly approved the shallow map.
- `superseded`: package was replaced by a later run or re-scope.

Use exactly one wave status in each wave file:

- `draft`: wave grouping is being prepared.
- `needs-revision`: reviewers requested changes.
- `questions-open`: open decomposition-blocking or owner-decision questions remain for the wave.
- `blocked`: missing source, unavailable reviewers, or exhausted budgets prevent shallow approval.
- `top-level-ready`: reviewers approved the shallow wave and no open blockers remain.
- `user-approved`: the user explicitly approved the shallow wave map including this wave.
- `superseded`: wave was replaced by a later approved re-scope.

## Stable IDs

- Waves: `WAVE-01`, `WAVE-02`, ...
- Capability groups: `CAP-W01-001`, `CAP-W02-001`, ...
- Outcomes: `OUT-W01-001`, `OUT-W02-001`, ...
- Frontend pages: `PAGE-001`, `PAGE-002`, ...
- Handoff checkpoints: `HANDOFF-W01-001`, `HANDOFF-W02-001`, ...
- Questions: `PQ-W01-001`, `PQ-W02-001`, ...
- Decisions: `PDEC-W01-001`, `PDEC-W02-001`, ...

## Required Reviewer Perspectives

Final approval requires approved verdicts from:

- `product-scope-coverage`
- `technical-boundary-fit`
- `sequencing-dependencies`
- `backend-wave-boundary-quality`
- `traceability-consistency`

Final approval also requires approved scope-review rows for:

- `product-capabilities`
- `user-journeys`
- `data-lifecycle`
- `integrations-operations`
- `client-experience`
- `security-compliance`
- `delivery-sequencing`
- `wave-map-consistency`

## Reviewer Verdict Table

Use this shape in `appendix/reviewer-verdicts.md`:

```text
| Scope | Perspective | Attempt | Verdict | Reviewer Report | Required Revisions | Notes |
| --- | --- | --- | --- | --- | --- | --- |
```

Allowed final verdicts: `approved`, `needs-revision`, `blocked`. Staging scaffolds may use `pending-review` only while validation is run with `--allow-placeholders`.

Final `Reviewer Report` values must point to scope-local `.tasks/prd-wave-decomposition/<run-id>/scopes/<scope>/review-attempt-<n>.md` files for primary scopes and `.tasks/prd-wave-decomposition/<run-id>/scopes/wave-map-consistency/consistency-attempt-<n>.md` for consistency. The consistency report must name the reviewed candidate package path `.tasks/prd-wave-decomposition/<run-id>/staging/prd-waves`, and that candidate package must exist in run evidence.

## Question Ledger Format

Use this table shape in every scope-local and aggregate ledger:

```text
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
```

Allowed severities: `decomposition-blocking`, `needs-owner-decision`, `deferred`, `watchlist`.

Allowed statuses: `open`, `answered`, `resolved`, `deferred`, `superseded`.

Open `decomposition-blocking` or `needs-owner-decision` rows block `top-level-ready`, block user approval, and block downstream detailed planning.

## Backend Wave And Frontend Page Gate

The package may be `waves-ready-awaiting-user-approval` only when:

- raw product source inventory exists and is non-empty;
- verified product source inventory exists and is non-empty;
- every required primary scope has an approved scope review;
- consistency review approved a concrete candidate package under `.tasks/prd-wave-decomposition/<run-id>/staging/prd-waves`;
- no package file contains implementation tasks, acceptance criteria, exit criteria, API payloads, schemas, component architecture, component designs, migrations, test cases, Beads, Jira, or code changes;
- maps with more than 8 top-level backend waves include an explicit `broader-release-scope-approved` or `non-mvp-release-scope-approved` decision in `appendix/decision-log.md`;
- every wave file has `top-level-ready` status;
- every wave file contains at least one `CAP-W<nn>-...` id and one `OUT-W<nn>-...` id;
- every wave file and `wave-map.md` contain only backend, data, integration, operations, security, and delivery-sequencing scope;
- frontend page, screen, route, navigation, mobile, UI, UX, user-interface, client-facing, and client-experience scope appears only under `frontend-pages/`, not in `waves/**` or `wave-map.md`;
- final packages do not contain `frontend-page-sequence.md`;
- `frontend-pages/index.md` has `top-level-ready` or `user-approved` status and contains ordered `PAGE-...` entries in `## Page Order`, matching `frontend-pages/page-<nnn>.md` files, raw PRD source coverage, verified PRD source coverage, and package-level traceability; or it contains an explicit `FRONTEND_NONE_CONFIRMED` / `FRONTEND_DEFERRED` marker in status or frontend deferrals when no frontend pages can be produced;
- every `frontend-pages/page-<nnn>.md` file has `top-level-ready` or `user-approved` status, includes its canonical `PAGE-<nnn>` id, appears in `frontend-pages/index.md` `## Page Order`, and has non-placeholder sections for page purpose, what is on the page, functional parts, empty states, loading/error states, backend dependencies, raw PRD traceability, and verified PRD traceability;
- all material backend PRD scope groups map to exactly one included backend wave or explicit deferral;
- all material frontend PRD scope maps to exactly one ordered page entry or explicit frontend deferral;
- the aggregate question ledger has no open `decomposition-blocking` or `needs-owner-decision` rows;
- backend source evidence and traceability point to verified product docs, technical docs, source deltas, scope reports, or explicit decisions;
- frontend page evidence and traceability point to both raw product docs and verified product docs, with technical docs, source deltas, scope reports, or explicit decisions used only as supplemental context.

`waves-approved` additionally requires explicit user approval recorded in every wave's `## User Approval` section and in `appendix/decision-log.md`. Use one of these literal markers in the approval evidence: `approved-by-user`, `user-approved`, or `waves-approved-by-user`. Each approved wave must also have status `user-approved`.
