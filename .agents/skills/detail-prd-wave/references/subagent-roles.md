<!-- FILE: .agents/skills/detail-prd-wave/references/subagent-roles.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Define controller, selected-backend-wave orchestrator, specialist planner, reviewer, and final fit roles for detailed backend PRD wave planning. -->
<!--   SCOPE: Covers orchestration ownership, one-backend-wave focus, codebase and neighboring-backend-wave fit, frontend-pages dependency context, reviewer perspectives, retry policy, report formats, question handling, and prompt templates. -->
<!--   DEPENDS: .agents/skills/detail-prd-wave/SKILL.md, .agents/skills/detail-prd-wave/references/output-contract.md. -->
<!--   LINKS: M-PRD-WAVE-DETAILER / V-M-PRD-WAVE-DETAILER. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Main Controller Contract - Defines what the main session owns. -->
<!--   Wave-Orchestrator Contract - Defines nested planner and reviewer loop ownership for one selected backend wave. -->
<!--   Planner Scopes - Lists specialist planner responsibilities. -->
<!--   Reviewer Perspectives - Lists required reviewers for ready-for-dev. -->
<!--   Prompt Templates - Provides copy-ready wave, planner, reviewer, and final fit prompts. -->
<!--   Write Ownership Matrix - Defines file ownership boundaries. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.1 - Clarified selected waves are backend-only and frontend planning is out of scope. -->
<!-- END_CHANGE_SUMMARY -->

# Detailed PRD Wave Orchestrators And Subagent Roles

The main session is the controller. It validates the selected backend source wave, inventories surrounding backend waves, the separate frontend-pages dependency context, GRACE, and codebase context, then starts exactly one `wave-orchestrator` subagent for that backend wave. The wave-orchestrator owns its backend wave, spawns specialist planners and reviewers internally, loops until approval or blocker, and writes only wave-local orchestration artifacts.

## Main Controller Contract

The main session must:

- Create `.tasks/prd-wave-detail/<run-id>/main-orchestration.md`.
- Create `.tasks/prd-wave-detail/<run-id>/source-wave-gate.md`.
- Create `.tasks/prd-wave-detail/<run-id>/context-inventory.md`.
- Create `.tasks/prd-wave-detail/<run-id>/source-delta.md` when prior questions were answered or source docs changed.
- Create staging skeletons only under `.tasks/prd-wave-detail/<run-id>/staging/prd-wave-details`.
- Dispatch exactly one wave-orchestrator for the selected wave.
- Give the wave-orchestrator write scope limited to `.tasks/prd-wave-detail/<run-id>/waves/<wave-id>/`.
- Aggregate selected-wave status and question ledgers.
- Dispatch final fit review only after the candidate selected-wave package exists.
- Promote final `docs/prd-wave-details/**` only for the selected wave after approvals or blockers are known.
- Stop after reporting a `ready-for-dev` wave and ask for explicit user approval.
- If nested subagents or reviewers are unavailable, write a `wave-blocking` open question row such as `DQ-W<nn>-001`, set the selected wave and package status to `blocked`, and do not synthesize unreviewed ready-for-dev output.

The main session must not:

- Detail more than one wave in a run.
- Plan frontend pages, routes, screens, navigation, UX states, components, frontend tests, or frontend implementation tasks.
- Dispatch planner workers or reviewers directly.
- Mark a wave `ready-for-dev` while open wave-blocking or owner-decision questions remain.
- Edit implementation code, generated artifacts, Beads, or Jira.
- Hide missing source or codebase evidence in assumptions.

## Wave-Orchestrator Contract

Each wave-orchestrator must:

- Create `.tasks/prd-wave-detail/<run-id>/waves/<wave-id>/orchestrator.md`.
- Create `.tasks/prd-wave-detail/<run-id>/waves/<wave-id>/wave-status.md`.
- Create `.tasks/prd-wave-detail/<run-id>/waves/<wave-id>/question-ledger.md`.
- Spawn all specialist planners for the selected wave.
- Spawn all required reviewer perspectives after planner reports exist.
- Repeat when any reviewer returns `needs-revision`.
- Relaunch interrupted, stalled, missing-report, or missing-verdict planner/reviewer attempts while interruption retry budget remains.
- Convert unresolved issues into wave-local ledger questions immediately.
- Mark wave status as `ready-for-dev`, `needs-revision`, `questions-open`, or `blocked`.
- When marking `blocked` or `questions-open`, write at least one matching open row to `question-ledger.md` with severity `wave-blocking` or `needs-owner-decision`.
- Use stable wave-prefixed ids.

Reviewer approval requires:

- all claims trace to `docs/prd-waves`, verified docs, GRACE XML, source deltas, code paths, prior detailed waves, or explicit decisions;
- implementation slices are scoped to the selected backend wave and independently useful;
- acceptance criteria and exit criteria are testable by developers and QA;
- verification obligations match the changed backend code, data, API, integration, security, and operations surface;
- neighboring backend waves are not duplicated or silently changed;
- frontend-pages references are dependency context only;
- no open blocker is hidden in assumptions.

## Planner Scopes

- `product-ac`: outcome, included/excluded backend scope, acceptance criteria, product edge cases, and non-frontend acceptance boundaries.
- `architecture-codebase`: existing modules, code paths, contracts, generated artifacts, likely graph deltas, and implementation slices.
- `data-integration-ops`: data lifecycle, API/events/jobs, external integrations, observability, rollout, rollback, and operations.
- `security-compliance`: auth, authorization, privacy, audit, compliance, abuse, rate limits, and secrets.
- `testing-exit`: exit criteria, backend verification obligations, focused commands, integration or contract needs, fixtures, traces, and evidence.
- `sequencing-fit`: prior detailed backend wave compatibility, future backend wave boundaries, frontend-pages dependency context, dependencies, deferrals, and independent value.

## Reviewer Perspectives

- `product-scope-and-ac`
- `architecture-codebase-fit`
- `data-api-integration-ops`
- `security-privacy-compliance`
- `testing-exit-criteria`
- `sequencing-other-wave-fit`
- `traceability-consistency`
- `final-wave-fit-review`

## Planner Report Format

Each planner writes `.tasks/prd-wave-detail/<run-id>/waves/<wave-id>/planner-<scope>-attempt-<n>.md`.

```text
# <Wave ID> <Scope> Planner Attempt <n>
## Sources Read
## Selected Backend Wave Boundary
## Neighboring Backend Wave Fit
## Frontend Pages Context
## Codebase Evidence
## Proposed Details
## Acceptance Criteria Contributions
## Exit Criteria Contributions
## Verification Contributions
## Risks And Rollback
## Questions Raised
## Traceability Candidates
```

## Reviewer Report Format

Each reviewer writes `.tasks/prd-wave-detail/<run-id>/waves/<wave-id>/review-<perspective>-attempt-<n>.md`.

```text
# <Wave ID> <Perspective> Review Attempt <n>
## Verdict
approved | needs-revision | blocked
## Sources Read
## Coverage Check
## Evidence Check
## Codebase Fit Check
## Other-Wave Fit Check
## Acceptance Criteria Check
## Exit Criteria Check
## Verification Check
## Question Ledger Check
## Unsupported Or Invented Claims
## Required Revisions
## Approval Notes
```

## Final Fit Report Format

Final fit writes `.tasks/prd-wave-detail/<run-id>/waves/<wave-id>/final-wave-fit-review-attempt-<n>.md`.

```text
# <Wave ID> Final Wave Fit Review Attempt <n>
## Verdict
approved | needs-revision | blocked
## Sources Read
## Candidate Package Reviewed
## One-Wave Focus Check
## Source Wave Gate Check
## Codebase Fit Check
## Neighboring Wave Fit Check
## AC EC Verification Check
## Reviewer Verdict Check
## Question Ledger Check
## Required Revisions
## Approval Notes
```

## Question Ledger Format

Each wave ledger and aggregate ledger use the output contract table:

```text
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
```

## Write Ownership Matrix

| Actor              | May write                                                                                                                                                                                                       | Must not write                                                                    |
| ------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------- |
| Main controller    | `.tasks/prd-wave-detail/<run-id>/main-orchestration.md`, `source-wave-gate.md`, `context-inventory.md`, `source-delta.md`, `recovery.md`, aggregate ledgers, staging skeleton, final `docs/prd-wave-details/**` | Wave-local planner/reviewer files except explicit controller intervention notes   |
| Wave-orchestrator  | `.tasks/prd-wave-detail/<run-id>/waves/<wave-id>/orchestrator.md`, `wave-status.md`, `question-ledger.md`                                                                                                       | Other waves, aggregate files, staging, final docs                                 |
| Planner specialist | `.tasks/prd-wave-detail/<run-id>/waves/<wave-id>/planner-<scope>-attempt-<n>.md`                                                                                                                                | Reviewer files, other waves, aggregate files, staging, final docs                 |
| Reviewer           | `.tasks/prd-wave-detail/<run-id>/waves/<wave-id>/review-<perspective>-attempt-<n>.md`                                                                                                                           | Planner files, other reviewers, other waves, aggregate files, staging, final docs |
| Final fit reviewer | `.tasks/prd-wave-detail/<run-id>/waves/<wave-id>/final-wave-fit-review-attempt-<n>.md`                                                                                                                          | Planner files, reviewer files, staging, final docs                                |

## Wave-Orchestrator Prompt Template

```text
Use the detail-prd-wave wave-orchestrator role for <WAVE_ID>.

PRD waves folder: <PRD_WAVES>
Product source folder: <PRODUCT>
Technical source folder: <TECHNICAL>
Output folder: <OUTPUT>
Staging folder: <STAGING>
Run id: <RUN_ID>
Wave id: <WAVE_ID>
Selected source wave: <SOURCE_WAVE_FILE>
Context inventory: .tasks/prd-wave-detail/<RUN_ID>/context-inventory.md
Source delta file if present: <SOURCE_DELTA>
Output contract: <OUTPUT_CONTRACT>
Role contract: <THIS_FILE>

Rules:
- Spawn specialist planner and reviewer subagents internally for this selected wave only.
- Read relevant code and neighboring waves before claiming codebase or sequencing fit.
- Write orchestration artifacts only under `.tasks/prd-wave-detail/<RUN_ID>/waves/<WAVE_ID>/`.
- Maintain this wave's question-ledger.md from the first gap.
- Use stable wave-prefixed ids: SLICE, AC, EC, TEST, HANDOFF, DQ, and DDEC.
- Do not detail later waves and do not edit implementation code.
- If nested subagents are unavailable, mark this wave blocked; do not ask the main session to run planner/reviewers.
```

## Planner Prompt Template

```text
Use the detail-prd-wave planner role: <SCOPE>.

Run id: <RUN_ID>
Wave id: <WAVE_ID>
Attempt: <N>
Report path: .tasks/prd-wave-detail/<RUN_ID>/waves/<WAVE_ID>/planner-<SCOPE>-attempt-<N>.md

Selected backend source wave: <SOURCE_WAVE_FILE>
Neighboring backend wave context: <NEIGHBORING_WAVE_SUMMARY>
Frontend pages context: <FRONTEND_PAGES_SUMMARY>
Relevant GRACE refs: <GRACE_REF_SUMMARY>
Relevant code refs: <CODE_REF_SUMMARY>
Previous reviewer findings if any: <REVIEWER_FINDINGS>

Write the planner report only. Do not edit docs/prd-wave-details.
Make only the selected backend wave ready for developers. Record source-backed backend slices, ACs, ECs, verification obligations, risks, rollback, questions, and traceability for your scope.
Do not plan frontend pages, routes, screens, navigation, UX states, components, frontend tests, or frontend implementation tasks.
Do not invent behavior beyond source docs, code evidence, prior approved waves, or explicit decisions.
```

## Reviewer Prompt Template

```text
Use the detail-prd-wave reviewer role: <PERSPECTIVE>.

Run id: <RUN_ID>
Wave id: <WAVE_ID>
Planner reports:
<PLANNER_REPORT_LIST>
Candidate selected-wave package if present: <CANDIDATE_PACKAGE>
Review path: .tasks/prd-wave-detail/<RUN_ID>/waves/<WAVE_ID>/review-<PERSPECTIVE>-attempt-<N>.md

Review for your perspective, source evidence, codebase fit, other-backend-wave fit, frontend-pages boundary compliance, AC quality, EC quality, verification obligations, question capture, and traceability.
Reject reports that invent unsupported details, hide blockers, skip codebase context, detail later backend waves, or plan frontend work.
Return one verdict: approved, needs-revision, or blocked.
```

## Final Fit Prompt Template

```text
Use the detail-prd-wave final-wave-fit-review role.

Run id: <RUN_ID>
Wave id: <WAVE_ID>
Candidate package: .tasks/prd-wave-detail/<RUN_ID>/staging/prd-wave-details
Reviewer verdict ledger: <REVIEWER_VERDICTS>
Question ledger: <QUESTION_LEDGER>
Review path: .tasks/prd-wave-detail/<RUN_ID>/waves/<WAVE_ID>/final-wave-fit-review-attempt-<N>.md

Review the candidate package for one-backend-wave focus, source-wave gate integrity, codebase fit, other-backend-wave fit, frontend-pages boundary compliance, AC/EC completeness, verification obligations, reviewer approvals, traceability, and open-question sync.
Reject packages that detail later backend waves, plan frontend work, invent unsupported contracts, miss codebase fit evidence, or mark ready-for-dev with open blockers.
Return one verdict: approved, needs-revision, or blocked.
```
