<!-- FILE: .agents/skills/verify-technical-docs/references/subagent-roles.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Define controller, scope-orchestrator, worker, reviewer, and consistency roles for the technical docs verification loop. -->
<!--   SCOPE: Covers orchestration ownership, retry policy, scope list, prompts, report formats, question-loop closure, and write boundaries. -->
<!--   DEPENDS: .agents/skills/verify-technical-docs/SKILL.md, .agents/skills/verify-technical-docs/references/output-contract.md. -->
<!--   LINKS: M-TECHNICAL-DOCS-VERIFIER / V-M-TECHNICAL-DOCS-VERIFIER. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Main Controller Contract - Defines what the main session owns. -->
<!--   Scope-Orchestrator Contract - Defines nested worker/reviewer loop ownership. -->
<!--   Scopes - Lists technical review scopes and their focus. -->
<!--   Question Loop Closure - Defines how answered questions can or cannot approve the package. -->
<!--   Prompt Templates - Provides copy-ready scope, worker, and reviewer prompts. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.0 - Added technical scope orchestration contract. -->
<!-- END_CHANGE_SUMMARY -->

# Technical Scope Orchestrators And Subagent Roles

The main session is the technical controller. It starts one `scope-orchestrator` subagent per technical scope. Each scope-orchestrator owns exactly one scope, spawns that scope's worker and reviewer internally, loops until reviewer approval or budget exhaustion, and writes only that scope's orchestration artifacts.

## Main Controller Contract

The main session must:

- Create `.tasks/technical-docs-verify/<run-id>/main-orchestration.md`.
- Create `.tasks/technical-docs-verify/<run-id>/source-delta.md` when product docs changed, prior technical output exists, prior ledgers exist, or questions were answered.
- Create staging skeletons only under `.tasks/technical-docs-verify/<run-id>/staging/technical-verified`.
- Dispatch Phase 1 scope-orchestrators for all required scopes.
- Dispatch `consistency-loop-reviewer` only after Phase 1 scope outputs exist.
- Give every scope-orchestrator write scope limited to `.tasks/technical-docs-verify/<run-id>/scopes/<scope>/`.
- Never dispatch scope workers or scope reviewers directly.
- Aggregate scope status and question ledgers.
- Write final `docs/technical-verified/**` with an accurate status: `questions-open`, `blocked`, or `approved-to-dev`.
- Stop `approved-to-dev` if any scope lacks reviewer approval, if nested worker/reviewer execution was bypassed, or if answer deltas created unresolved blockers.

The main session must not:

- Invent endpoints, schemas, event payloads, auth rules, infra topology, SLOs, migrations, or test gates not present in sources or explicit decisions.
- Hide unresolved technical gaps in assumptions.
- Start implementation planning from a `questions-open` package unless the user explicitly accepts the risk.

## Scope-Orchestrator Contract

Each scope-orchestrator must:

- Create `.tasks/technical-docs-verify/<run-id>/scopes/<scope>/orchestrator.md`.
- Create `.tasks/technical-docs-verify/<run-id>/scopes/<scope>/scope-status.md`.
- Create `.tasks/technical-docs-verify/<run-id>/scopes/<scope>/question-ledger.md`.
- Spawn a worker subagent for its scope.
- Spawn a reviewer subagent after each worker report.
- Repeat when the reviewer returns `needs-revision`.
- Relaunch interrupted, stalled, missing-report, or missing-verdict worker/reviewer attempts while interruption retry budget remains.
- Convert unresolved scope issues into ledger questions immediately.
- Mark scope status as `approved`, `needs-revision`, `blocked`, or `interrupted`.
- Do not defer a required scope. Deferral is allowed only for individual ledger questions with owner, rationale, and no implementation-blocking impact.
- Use stable scope-prefixed ids.

Reviewer approval requires:

- all technical claims trace to product-verified docs, prior verified technical docs, source delta, explicit owner decisions, or scope reports;
- all missing artifact classes are consolidated;
- no product behavior or implementation contract is invented;
- question severities and statuses match the output contract;
- answer deltas are reviewed for second-order effects;
- the scope report can be safely synthesized into `docs/technical-verified`.

## Scopes

### architecture-boundaries

Focus on system context, component boundaries, ownership, tenancy, deployment boundary, service boundaries, build-vs-buy boundaries, and architecture decisions implied by product behavior.

### data-contracts

Focus on entities, identifiers, relationships, persistence, storage engines, migrations, seed data, fixtures, retention, privacy, imports, exports, and data lineage.

### api-contracts

Focus on public and internal API surfaces, GraphQL/REST/RPC choices, request/response schemas, error formats, validation mapping, pagination, filtering, idempotency, compatibility, versioning, and client/server contract gaps.

### auth-security-compliance

Focus on identity, authentication, authorization, ownership, tenant scoping, audit, rate limiting, abuse prevention, secrets, privacy, compliance, legal retention, and irreversible action controls.

### integrations-events

Focus on external systems, sync ownership, async jobs, events, webhooks, queues, retries, backoff, dead letters, reconciliation, rate limits, and failure handling.

### client-state-ux

Focus on UI state machines, loading/empty/error/offline states, form validation, optimistic updates, realtime updates, cache invalidation, accessibility, localization, and user-facing technical edge cases.

### operations-observability

Focus on environments, config, deployment, rollout, feature flags, monitoring, logs, metrics, traces, alerts, runbooks, backup/restore, SLOs, capacity, and operational ownership.

### testing-delivery

Focus on test strategy, contract tests, integration tests, e2e flows, fixtures, seed data, test isolation, coverage gates, release criteria, QA handoff, and pre-MR checks.

### consistency-loop-reviewer

Run after Phase 1. Focus on contradictions, duplicate questions, severity drift, unresolved parent questions, answer effects, and whether the package can be `approved-to-dev`.

## Question Loop Closure

An answered question closes only when:

1. The original question id is preserved.
2. The answer source is recorded as source file, user answer, or explicit decision.
3. Every affected scope reviewed the answer.
4. No affected scope raises an open `dev-blocking` or `needs-owner-decision` follow-up.
5. Consistency confirms the answer did not contradict other technical decisions or product evidence.

If an answer creates a follow-up blocker, the package status is `questions-open`. Link the follow-up question to the parent id with the `Parent` column.

## Worker Report Format

Each worker writes to `.tasks/technical-docs-verify/<run-id>/scopes/<scope>/worker-attempt-<n>.md`.

```text
# <Scope> Worker Attempt <n>
## Sources Read
## Source Delta Reviewed
## Product Signals
## Technical Facts
## Technical Gaps
## Missing Source Artifacts
## Questions Raised
## Answer Effects
## Risks
## Suggested Decisions
## Traceability Candidates
```

## Reviewer Report Format

Each reviewer writes to `.tasks/technical-docs-verify/<run-id>/scopes/<scope>/review-attempt-<n>.md`.

```text
# <Scope> Review Attempt <n>
## Verdict
approved | needs-revision | blocked
## Sources Read
## Coverage Check
## Evidence Check
## No-Invention Check
## Source-Gap Consolidation Check
## Question Ledger Check
## Answer Effect Check
## Missing Or Unsupported Claims
## Required Revisions
## Approval Notes
```

## Question Ledger Format

Each scope ledger and the aggregate ledger use the output contract table:

```text
| ID | Scope | Severity | Parent | Question | Why It Matters | Needed Artifact Or Decision | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
```

## Write Ownership Matrix

| Actor              | May write                                                                                                                                                                                                          | Must not write                                                                  |
| ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------- |
| Main controller    | `.tasks/technical-docs-verify/<run-id>/main-orchestration.md`, `source-delta.md`, `recovery.md`, aggregate `scope-status.md`, aggregate `question-ledger.md`, staging skeleton, final `docs/technical-verified/**` | Scope-local worker/reviewer files except explicit controller intervention notes |
| Scope-orchestrator | `.tasks/technical-docs-verify/<run-id>/scopes/<scope>/orchestrator.md`, `scope-status.md`, `question-ledger.md`                                                                                                    | Other scopes, aggregate files, staging, final docs                              |
| Scope worker       | `.tasks/technical-docs-verify/<run-id>/scopes/<scope>/worker-attempt-<n>.md`                                                                                                                                       | Reviewer files, other scopes, aggregate files, staging, final docs              |
| Scope reviewer     | `.tasks/technical-docs-verify/<run-id>/scopes/<scope>/review-attempt-<n>.md`                                                                                                                                       | Worker files, other scopes, aggregate files, staging, final docs                |
| Consistency scope  | `.tasks/technical-docs-verify/<run-id>/scopes/consistency-loop-reviewer/**`                                                                                                                                        | Phase 1 scope folders, staging, final docs                                      |

## Scope-Orchestrator Prompt Template

```text
Use the verify-technical-docs scope-orchestrator role: <SCOPE_NAME>.

Input folder: <SOURCE>
Output folder: <OUTPUT>
Staging folder: <STAGING>
Source delta file if present: <SOURCE_DELTA>
Run id: <RUN_ID>
Scope: <SCOPE_NAME>
Scope focus: <SCOPE_FOCUS>
Output contract: <OUTPUT_CONTRACT>
Role contract: <THIS_FILE>
Available source files:
<SOURCE_FILE_LIST>

Rules:
- Spawn worker and reviewer subagents internally for this scope only.
- Keep product-verified inputs read-only.
- Write orchestration artifacts only under `.tasks/technical-docs-verify/<RUN_ID>/scopes/<SCOPE_NAME>/`.
- Maintain this scope's question-ledger.md from the first gap.
- Preserve original question ids when answers are present.
- Link follow-up questions to parent ids.
- Do not claim this scope is approved until its reviewer approves it.
- If nested subagents are unavailable, mark this scope blocked; do not ask the main session to run worker/reviewer.
```

## Scoped Worker Prompt Template

```text
Use the verify-technical-docs scoped worker role: <SCOPE_NAME>.

Input folder: <SOURCE>
Run id: <RUN_ID>
Attempt: <N>
Report path: .tasks/technical-docs-verify/<RUN_ID>/scopes/<SCOPE_NAME>/worker-attempt-<N>.md

Role focus:
<SCOPE_FOCUS>

Available source files:
<SOURCE_FILE_LIST>

Previous reviewer findings if any:
<REVIEWER_FINDINGS>

Source delta if present:
<SOURCE_DELTA_SUMMARY>

Write the worker report only. Do not edit docs/technical-verified.
Record every missing implementation-critical artifact as a ledger question.
Consolidate missing source classes.
When answers are present, analyze whether they resolve, supersede, or create follow-up questions.
Do not invent endpoints, schemas, event payloads, auth rules, infra topology, SLOs, migrations, or test gates.
```

## Scoped Reviewer Prompt Template

```text
Use the verify-technical-docs scoped reviewer role: <SCOPE_NAME>.

Input folder: <SOURCE>
Run id: <RUN_ID>
Worker report: .tasks/technical-docs-verify/<RUN_ID>/scopes/<SCOPE_NAME>/worker-attempt-<N>.md
Review path: .tasks/technical-docs-verify/<RUN_ID>/scopes/<SCOPE_NAME>/review-attempt-<N>.md

Review for evidence, missing information capture, consolidated source gaps, answer effects, and traceability.
Reject reports that invent implementation details beyond source evidence or explicit decisions.
Reject reports that split one absent artifact into many speculative questions.
If not approved, list exact required revisions.
Do not edit docs/technical-verified.

Return one verdict: approved, needs-revision, or blocked.
```
