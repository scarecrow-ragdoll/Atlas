<!-- FILE: .agents/skills/decompose-prd-waves/references/subagent-roles.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Define controller, scope-mapper, mapper worker, reviewer, and consistency roles for backend-wave and frontend-page decomposition. -->
<!--   SCOPE: Covers orchestration ownership, backend-only wave rules, per-page frontend rules, reviewer perspectives, retry policy, report formats, question handling, and prompt templates. -->
<!--   DEPENDS: .agents/skills/decompose-prd-waves/SKILL.md, .agents/skills/decompose-prd-waves/references/output-contract.md. -->
<!--   LINKS: M-PRD-WAVE-DECOMPOSER / V-M-PRD-WAVE-DECOMPOSER. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Main Controller Contract - Defines what the main session owns. -->
<!--   Scope-Mapper Contract - Defines nested mapper/reviewer loop ownership for one scope. -->
<!--   Reviewer Perspectives - Lists required reviewers for final approval. -->
<!--   Prompt Templates - Provides copy-ready scope-mapper, mapper, reviewer, and consistency prompts. -->
<!--   Write Ownership Matrix - Defines file ownership boundaries. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.2 - Replaced frontend page-sequence inputs with per-page frontend inputs sourced from raw and verified PRDs. -->
<!-- END_CHANGE_SUMMARY -->

# PRD Wave Decomposition Subagent Roles

The main session is the controller. It validates source readiness, prepares source inventory and staging, starts one `scope-mapper` subagent per scope, aggregates approved scope outputs, synthesizes a candidate package with backend-only waves plus per-page frontend files under staging, starts final consistency review on that candidate, and promotes `docs/prd-waves/**` only after required approvals or blockers are known.

## Main Controller Contract

The main session must:

- Create `.tasks/prd-wave-decomposition/<run-id>/main-orchestration.md`.
- Create `.tasks/prd-wave-decomposition/<run-id>/source-gate.md`.
- Create `.tasks/prd-wave-decomposition/<run-id>/source-delta.md` when prior questions were answered or sources changed.
- Create staging skeletons and candidate wave packages only under `.tasks/prd-wave-decomposition/<run-id>/staging/prd-waves`.
- Dispatch one scope-mapper per primary scope.
- Give each scope-mapper write scope limited to `.tasks/prd-wave-decomposition/<run-id>/scopes/<scope>/`.
- Aggregate scope statuses and question ledgers.
- Keep `wave-map.md`, `waves/index.md`, and `waves/wave-<nn>.md` backend-only.
- Put frontend page, route, navigation, mobile, UI, UX, and client-experience scope only under `frontend-pages/`.
- Dispatch final `wave-map-consistency` after primary scope outputs and the candidate wave package exist.
- Pass the concrete candidate package path to final consistency.
- Promote final `docs/prd-waves/**` only from a candidate package that passed final consistency review.
- Stop after reporting a `waves-ready-awaiting-user-approval` map and ask for explicit user approval before downstream detailed planning.

The main session must not:

- Run this skill against raw, unverified product docs when `$verify-product-docs` is needed.
- Dispatch mapper workers or reviewers directly.
- Synthesize final docs from unreviewed mapper output or from an unreviewed candidate package.
- Turn frontend page or UX scope into implementation waves.
- Add implementation tasks, technical designs, acceptance criteria, test cases, Jira, Beads, or code.
- Treat chat approval as a substitute for missing source traceability.

## Scope-Mapper Contract

Each scope-mapper must:

- Create `.tasks/prd-wave-decomposition/<run-id>/scopes/<scope>/orchestrator.md`.
- Create `.tasks/prd-wave-decomposition/<run-id>/scopes/<scope>/scope-status.md`.
- Create `.tasks/prd-wave-decomposition/<run-id>/scopes/<scope>/question-ledger.md`.
- Spawn a mapper worker for the assigned scope.
- Spawn a scope reviewer after each mapper report.
- Repeat when the reviewer returns `needs-revision`.
- Relaunch interrupted, stalled, missing-report, or missing-verdict mapper/reviewer attempts while interruption retry budget remains.
- Convert unresolved issues into scope-local ledger questions immediately.
- Mark scope status as `approved`, `needs-revision`, `questions-open`, or `blocked`.

Reviewer approval requires:

- source-traceable capability and concern inventory for the assigned scope;
- proposed backend wave grouping at shallow level only when the assigned scope has backend implementation implications; otherwise an explicit `NO_BACKEND_WAVE_CONTRIBUTION` note;
- frontend page-file inputs when the scope touches client experience, journeys, navigation, mobile, or UI/UX;
- explicit exclusions, deferrals, dependencies, and open questions;
- no invented product or technical behavior;
- no hidden implementation detail;
- enough traceability for the main controller to build the backend wave map and per-page frontend files.

## Primary Scopes

- `product-capabilities`: product features, backend-owned business capabilities, roles, permissions, and explicit in/out scope.
- `user-journeys`: actor journeys, lifecycle states, handoffs, backend dependencies, and frontend page-file inputs.
- `data-lifecycle`: domain entities, data ownership, lifecycle, reporting, imports/exports, retention, and migration pressure at category level.
- `integrations-operations`: external systems, async work, notifications, observability, rollout, support, and operational readiness at category level.
- `client-experience`: frontend/admin/public/mobile pages, navigation, UX states, accessibility, and content flows as per-page frontend inputs only; never as wave proposals.
- `security-compliance`: auth, authorization, audit, privacy, compliance, abuse/rate-limit, and risk surfaces at category level.
- `delivery-sequencing`: dependencies, stage order, MVP cuts, prerequisite backend waves, parallel-safe groups, frontend page order dependencies, and more-than-8-backend-wave risk.

## Reviewer Perspectives

Final consistency must verify these perspectives:

- `product-scope-coverage`
- `technical-boundary-fit`
- `sequencing-dependencies`
- `backend-wave-boundary-quality`
- `traceability-consistency`

## Mapper Report Format

The mapper writes `.tasks/prd-wave-decomposition/<run-id>/scopes/<scope>/mapper-attempt-<n>.md`.

```text
# <Scope> Mapper Attempt <n>
## Sources Read
## Scope Inventory
## Candidate Backend Wave Contributions
## Frontend Page File Inputs
## Dependencies And Order Signals
## Explicit Exclusions And Deferrals
## Risk Signals
## Open Questions
## Traceability Candidates
## Shallow-Only Check
```

## Reviewer Report Format

The scope reviewer writes `.tasks/prd-wave-decomposition/<run-id>/scopes/<scope>/review-attempt-<n>.md`.

```text
# <Scope> Review Attempt <n>
## Verdict
approved | needs-revision | blocked
## Sources Read
## Coverage Check
## Evidence Check
## Shallow-Only Check
## Dependency Check
## Question Ledger Check
## Unsupported Or Invented Claims
## Required Revisions
## Approval Notes
```

## Consistency Report Format

Final consistency writes `.tasks/prd-wave-decomposition/<run-id>/scopes/wave-map-consistency/consistency-attempt-<n>.md`.

```text
# Wave Map Consistency Attempt <n>
## Verdict
approved | needs-revision | blocked
## Sources Read
## Candidate Package Reviewed
## Scope Coverage Check
## Wave Duplication Check
## Dependency Order Check
## Shallow-Only Check
## More Than Eight Wave Check
## Traceability Check
## Question Ledger Check
## Required Revisions
## Approval Notes
```

## Question Ledger Format

Each scope ledger and aggregate ledger use the output contract table:

```text
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
```

## Write Ownership Matrix

| Actor                | May write                                                                                                                                                                                             | Must not write                                                                  |
| -------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------- |
| Main controller      | `.tasks/prd-wave-decomposition/<run-id>/main-orchestration.md`, `source-gate.md`, `source-delta.md`, `recovery.md`, aggregate ledgers, staging skeleton, candidate package, final `docs/prd-waves/**` | Scope-local mapper/reviewer files except explicit controller intervention notes |
| Scope-mapper         | `.tasks/prd-wave-decomposition/<run-id>/scopes/<scope>/orchestrator.md`, `scope-status.md`, `question-ledger.md`                                                                                      | Other scopes, aggregate files, staging, final docs                              |
| Mapper worker        | `.tasks/prd-wave-decomposition/<run-id>/scopes/<scope>/mapper-attempt-<n>.md`                                                                                                                         | Reviewer files, other scopes, aggregate files, staging, final docs              |
| Scope reviewer       | `.tasks/prd-wave-decomposition/<run-id>/scopes/<scope>/review-attempt-<n>.md`                                                                                                                         | Mapper files, other scopes, aggregate files, staging, final docs                |
| Consistency reviewer | `.tasks/prd-wave-decomposition/<run-id>/scopes/wave-map-consistency/consistency-attempt-<n>.md`                                                                                                       | Primary scope files, staging, final docs                                        |

## Scope-Mapper Prompt Template

```text
Use the decompose-prd-waves scope-mapper role for <SCOPE>.

Raw product source folder: <RAW_PRODUCT>
Verified product source folder: <PRODUCT>
Technical source folder: <TECHNICAL>
Output folder: <OUTPUT>
Staging folder: <STAGING>
Run id: <RUN_ID>
Scope: <SCOPE>
Source delta file if present: <SOURCE_DELTA>
Output contract: <OUTPUT_CONTRACT>
Role contract: <THIS_FILE>
Available source files:
<SOURCE_FILE_LIST>

Rules:
- Spawn mapper and reviewer subagents internally for this scope only.
- Keep all output shallow and top-level.
- Write orchestration artifacts only under `.tasks/prd-wave-decomposition/<RUN_ID>/scopes/<SCOPE>/`.
- Maintain this scope's question-ledger.md from the first gap.
- Do not create implementation tasks, acceptance criteria, exit criteria, Beads, Jira, or code.
- If nested subagents are unavailable, mark this scope blocked; do not ask the main session to run mapper or reviewer roles.
```

## Mapper Prompt Template

```text
Use the decompose-prd-waves mapper role for <SCOPE>.

Raw product source folder: <RAW_PRODUCT>
Verified product source folder: <PRODUCT>
Technical source folder: <TECHNICAL>
Run id: <RUN_ID>
Attempt: <N>
Report path: .tasks/prd-wave-decomposition/<RUN_ID>/scopes/<SCOPE>/mapper-attempt-<N>.md

Scope focus:
<SCOPE_FOCUS>

Source delta if present:
<SOURCE_DELTA_SUMMARY>

Write the mapper report only. Do not edit docs/prd-waves.
Identify source-backed scope groups, likely backend wave contributions, frontend page-file inputs when applicable, dependencies, exclusions, deferrals, risks, and open questions.
Keep everything shallow. Do not turn frontend pages, routes, navigation, mobile, UI, UX, or client-experience scope into waves. Do not write implementation tasks, technical designs, acceptance criteria, tests, Beads, Jira, or code.
```

## Reviewer Prompt Template

```text
Use the decompose-prd-waves reviewer role for <SCOPE>.

Raw product source folder: <RAW_PRODUCT>
Verified product source folder: <PRODUCT>
Technical source folder: <TECHNICAL>
Run id: <RUN_ID>
Scope: <SCOPE>
Mapper report: .tasks/prd-wave-decomposition/<RUN_ID>/scopes/<SCOPE>/mapper-attempt-<N>.md
Review path: .tasks/prd-wave-decomposition/<RUN_ID>/scopes/<SCOPE>/review-attempt-<N>.md

Review for source coverage, backend-only wave compliance, per-page frontend handling, shallow-only compliance, dependency signals, question capture, and traceability.
Reject mapper reports that do not distinguish raw PRD evidence from verified PRD evidence for frontend page files.
Reject reports that invent behavior or drift into implementation detail.
Reject reports that leave open blockers outside the question ledger.
Return one verdict: approved, needs-revision, or blocked.
```

## Consistency Prompt Template

```text
Use the decompose-prd-waves consistency reviewer role.

Raw product source folder: <RAW_PRODUCT>
Verified product source folder: <PRODUCT>
Technical source folder: <TECHNICAL>
Run id: <RUN_ID>
Attempt: <N>
Approved scope reports:
<SCOPE_REPORT_LIST>
Aggregate question ledger: .tasks/prd-wave-decomposition/<RUN_ID>/question-ledger.md
Candidate package: .tasks/prd-wave-decomposition/<RUN_ID>/staging/prd-waves
Consistency report path: .tasks/prd-wave-decomposition/<RUN_ID>/scopes/wave-map-consistency/consistency-attempt-<N>.md

Review the proposed candidate package for complete PRD scope coverage across backend-only waves and per-page frontend files, duplicate or missing backend wave coverage, missing frontend pages/deferrals, missing raw PRD or verified PRD page traceability, dependency order, shallow-only compliance, more-than-eight-backend-wave re-scope handling, traceability, and open-question sync.
Reject maps that invent behavior, hide blockers, omit material PRD scope, or include implementation tasks, acceptance criteria, exit criteria, API payloads, schemas, migrations, test cases, Beads, Jira, or code-level detail.
Do not edit docs/prd-waves.
Return one verdict: approved, needs-revision, or blocked.
```
