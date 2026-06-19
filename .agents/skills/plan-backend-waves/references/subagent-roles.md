<!-- FILE: .agents/skills/plan-backend-waves/references/subagent-roles.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Define controller, wave-orchestrator, planner, reviewer, and consistency roles for backend wave planning. -->
<!--   SCOPE: Covers orchestration ownership, strict sequential wave rules, reviewer perspectives, retry policy, report formats, question handling, and prompt templates. -->
<!--   DEPENDS: .agents/skills/plan-backend-waves/SKILL.md, .agents/skills/plan-backend-waves/references/output-contract.md. -->
<!--   LINKS: M-BACKEND-WAVE-PLANNER / V-M-BACKEND-WAVE-PLANNER. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Main Controller Contract - Defines what the main session owns. -->
<!--   Wave-Orchestrator Contract - Defines nested planner/reviewer loop ownership for one wave. -->
<!--   Reviewer Perspectives - Lists required reviewers for each wave. -->
<!--   Prompt Templates - Provides copy-ready wave, planner, and reviewer prompts. -->
<!--   Write Ownership Matrix - Defines file ownership boundaries. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.0 - Added backend wave orchestration role contract. -->
<!-- END_CHANGE_SUMMARY -->

# Backend Wave Orchestrators And Subagent Roles

The main session is the controller. It validates the technical approval gate, prepares source inventory and a shallow wave map, then starts exactly one `wave-orchestrator` subagent for the current wave. The wave-orchestrator owns its wave, spawns one planner worker and all required reviewers internally, loops until approval or blocker, and writes only wave-local orchestration artifacts.

## Main Controller Contract

The main session must:

- Create `.tasks/backend-wave-plan/<run-id>/main-orchestration.md`.
- Create `.tasks/backend-wave-plan/<run-id>/technical-approval-gate.md`.
- Create `.tasks/backend-wave-plan/<run-id>/source-delta.md` when prior questions were answered or source docs changed.
- Create staging skeletons only under `.tasks/backend-wave-plan/<run-id>/staging/backend-waves`.
- Build a shallow backend wave map from approved technical docs.
- Stop for re-scope when the tentative backend plan exceeds 8 waves unless the user explicitly confirms the broader release scope.
- Dispatch exactly one wave-orchestrator for the current wave.
- Give the wave-orchestrator write scope limited to `.tasks/backend-wave-plan/<run-id>/waves/<wave-id>/`.
- Aggregate current wave status and question ledgers.
- Write final `docs/backend-waves/**` for the current wave only after wave reviewer approvals or blockers are known.
- Stop after reporting a `ready-for-dev` wave and ask for explicit user approval.
- Plan the next wave only after the previous wave has explicit user approval recorded.

The main session must not:

- Start this skill from `questions-open` or `blocked` technical docs.
- Plan all waves in implementation-level detail in one run.
- Dispatch planner workers or wave reviewers directly.
- Mark a wave `ready-for-dev` while open wave-blocking or owner-decision questions remain.
- Treat Jira-ready text as permission to mutate Jira.
- Hide missing information in assumptions.

## Wave-Orchestrator Contract

Each wave-orchestrator must:

- Create `.tasks/backend-wave-plan/<run-id>/waves/<wave-id>/orchestrator.md`.
- Create `.tasks/backend-wave-plan/<run-id>/waves/<wave-id>/wave-status.md`.
- Create `.tasks/backend-wave-plan/<run-id>/waves/<wave-id>/question-ledger.md`.
- Spawn a planner worker for the current wave.
- Spawn all required reviewer perspectives after each planner report.
- Repeat when any reviewer returns `needs-revision`.
- Relaunch interrupted, stalled, missing-report, or missing-verdict planner/reviewer attempts while interruption retry budget remains.
- Convert unresolved issues into wave-local ledger questions immediately.
- Mark wave status as `ready-for-dev`, `needs-revision`, `questions-open`, or `blocked`.
- Use stable wave-prefixed ids.

Reviewer approval requires:

- all claims trace to `docs/technical-verified`, `docs/product-verified`, source deltas, prior approved waves, or explicit decisions;
- the wave is independently implementable and does not depend on unplanned later work except named feature flags or compatibility scaffolding;
- acceptance criteria and exit criteria are executable by developers and QA;
- test obligations match the backend risk surface;
- no open blocker is hidden in assumptions;
- the wave can be shown to the user as ready for approval.

## Reviewer Perspectives

### backend-architecture

Review service boundaries, module ownership, layering, dependency direction, transaction boundaries, concurrency, compatibility, and whether the wave is independently shippable.

### data-api-contract

Review persistence, migrations, seed data, fixtures, API contracts, request/response shapes, errors, idempotency, pagination, jobs, events, and compatibility.

### security-integration

Review auth, authorization, tenant/owner scope, audit, secrets, rate limits, external integrations, retries, reconciliation, and failure handling.

### testing-delivery

Review unit, contract, integration, e2e, fixture, migration, observability, rollout, pre-MR, and release-gate coverage.

### sequencing-mvp

Review wave order, dependencies, MVP fit, more-than-8-wave smell, scope cuts, deferrals, and whether the outcome after this wave is meaningful.

### traceability-consistency

Review source traceability, duplicate/contradictory requirements, criteria-to-task mapping, open question sync, and consistency with approved technical docs.

## Planner Report Format

The planner writes `.tasks/backend-wave-plan/<run-id>/waves/<wave-id>/planner-attempt-<n>.md`.

```text
# <Wave ID> Planner Attempt <n>
## Sources Read
## Technical Approval Gate
## Wave Goal
## Outcome After Implementation
## Source Evidence
## Scope Included
## Scope Excluded
## Dependencies
## Backend Design
## Data And Migration Work
## API Jobs And Events
## Auth Security And Compliance
## Operations Observability
## Implementation Tasks
## Acceptance Criteria
## Exit Criteria
## Verification Plan
## Rollback And Compatibility
## Jira Ready Tasks
## Questions Raised
## Traceability Candidates
```

## Reviewer Report Format

Each reviewer writes `.tasks/backend-wave-plan/<run-id>/waves/<wave-id>/review-<perspective>-attempt-<n>.md`.

```text
# <Wave ID> <Perspective> Review Attempt <n>
## Verdict
approved | needs-revision | blocked
## Sources Read
## Coverage Check
## Evidence Check
## Ready-For-Dev Check
## Acceptance Criteria Check
## Exit Criteria Check
## Verification Check
## Question Ledger Check
## Unsupported Or Invented Claims
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

| Actor             | May write                                                                                                                                                                                     | Must not write                                                                    |
| ----------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------- |
| Main controller   | `.tasks/backend-wave-plan/<run-id>/main-orchestration.md`, `technical-approval-gate.md`, `source-delta.md`, `recovery.md`, aggregate ledgers, staging skeleton, final `docs/backend-waves/**` | Wave-local planner/reviewer files except explicit controller intervention notes   |
| Wave-orchestrator | `.tasks/backend-wave-plan/<run-id>/waves/<wave-id>/orchestrator.md`, `wave-status.md`, `question-ledger.md`                                                                                   | Other waves, aggregate files, staging, final docs                                 |
| Wave planner      | `.tasks/backend-wave-plan/<run-id>/waves/<wave-id>/planner-attempt-<n>.md`                                                                                                                    | Reviewer files, other waves, aggregate files, staging, final docs                 |
| Wave reviewer     | `.tasks/backend-wave-plan/<run-id>/waves/<wave-id>/review-<perspective>-attempt-<n>.md`                                                                                                       | Planner files, other reviewers, other waves, aggregate files, staging, final docs |

## Wave-Orchestrator Prompt Template

```text
Use the plan-backend-waves wave-orchestrator role for <WAVE_ID>.

Technical source folder: <TECHNICAL>
Product source folder: <PRODUCT>
Output folder: <OUTPUT>
Staging folder: <STAGING>
Run id: <RUN_ID>
Wave id: <WAVE_ID>
Wave map: <TENTATIVE_WAVE_MAP>
Source delta file if present: <SOURCE_DELTA>
Output contract: <OUTPUT_CONTRACT>
Role contract: <THIS_FILE>
Available source files:
<SOURCE_FILE_LIST>

Rules:
- Spawn planner and reviewer subagents internally for this wave only.
- Do not plan later waves in implementation-level detail.
- Write orchestration artifacts only under `.tasks/backend-wave-plan/<RUN_ID>/waves/<WAVE_ID>/`.
- Maintain this wave's question-ledger.md from the first gap.
- Use stable wave-prefixed ids.
- Do not claim ready-for-dev until every required reviewer approves and no open wave-blocking or owner-decision questions remain.
- If nested subagents are unavailable, mark this wave blocked; do not ask the main session to run planner/reviewers.
```

## Wave Planner Prompt Template

```text
Use the plan-backend-waves planner role for <WAVE_ID>.

Technical source folder: <TECHNICAL>
Product source folder: <PRODUCT>
Run id: <RUN_ID>
Attempt: <N>
Report path: .tasks/backend-wave-plan/<RUN_ID>/waves/<WAVE_ID>/planner-attempt-<N>.md

Wave goal:
<WAVE_GOAL>

Tentative wave map:
<TENTATIVE_WAVE_MAP>

Previous reviewer findings if any:
<REVIEWER_FINDINGS>

Source delta if present:
<SOURCE_DELTA_SUMMARY>

Write the planner report only. Do not edit docs/backend-waves.
Make the wave ready for developers: implementation tasks, ACs, exit criteria, tests, migrations, APIs/jobs/events, auth/security, ops, rollback, Jira-ready task text, and traceability.
Record every unresolved implementation or owner decision as a question.
Do not invent backend contracts beyond approved technical docs, source deltas, prior approved waves, or explicit decisions.
```

## Wave Reviewer Prompt Template

```text
Use the plan-backend-waves reviewer role: <PERSPECTIVE>.

Technical source folder: <TECHNICAL>
Product source folder: <PRODUCT>
Run id: <RUN_ID>
Wave id: <WAVE_ID>
Planner report: .tasks/backend-wave-plan/<RUN_ID>/waves/<WAVE_ID>/planner-attempt-<N>.md
Review path: .tasks/backend-wave-plan/<RUN_ID>/waves/<WAVE_ID>/review-<PERSPECTIVE>-attempt-<N>.md

Review for your perspective, source evidence, ready-for-dev detail, question capture, acceptance criteria, exit criteria, verification, and traceability.
Reject reports that invent backend details beyond the approved technical package, source deltas, prior approved waves, or explicit decisions.
Reject reports that leave open blockers outside the question ledger.
If not approved, list exact required revisions.
Do not edit docs/backend-waves.

Return one verdict: approved, needs-revision, or blocked.
```
