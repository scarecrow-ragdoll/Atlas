# Scope Orchestrators And Subagent Roles

The main session is the product-level controller. It starts one `scope-orchestrator` subagent per review scope. Each scope-orchestrator owns exactly one scope: it spawns that scope's worker, spawns that scope's reviewer, loops until reviewer approval or budget exhaustion, and writes only that scope's orchestration artifacts.

## Quick Map

- `Main Controller Contract`: what the main session owns, including staging and final synthesis.
- `Scope-Orchestrator Contract`: what each scope-orchestrator must run internally.
- `Worker/Reviewer Report Format`: required files for scoped worker/reviewer loops.
- `Retry And Relaunch Policy`: what to do when an orchestrator, worker, reviewer, or consistency attempt is interrupted or returns needs-revision.
- `Source Delta And Re-Verification`: how new docs and answered questions are folded into the normal pipeline.
- `Source-Gap Consolidation`: how missing artifact classes become blocker questions.
- `Evidence-Constrained Derivation`: when roles, fields, states, criteria, and edge cases may be derived.
- `Write Ownership Matrix`: exact write boundaries for main, scope-orchestrators, workers, and reviewers.
- `Dispatch Phases` and `Scopes`: phase order and scope focus.
- `Scoped Worker/Reviewer Prompt Template`: copy-ready prompts for the scope-orchestrator.

## Main Controller Contract

The main session must:

- Create `.tasks/product-docs-verify/<run-id>/main-orchestration.md`.
- Create `.tasks/product-docs-verify/<run-id>/source-delta.md` when new docs, changed docs, removed docs, previous verified output, previous question ledgers, or answered questions are present.
- Create staging skeletons only under `.tasks/product-docs-verify/<run-id>/staging/product-verified`.
- Dispatch Phase 1 scope-orchestrators for primary scopes in parallel when possible.
- Dispatch the `consistency-reviewer` scope-orchestrator only after Phase 1 scope outputs are available.
- Give every scope-orchestrator write scope limited to `.tasks/product-docs-verify/<run-id>/scopes/<scope>/`.
- Never dispatch scope workers or scope reviewers directly.
- Aggregate `.tasks/product-docs-verify/<run-id>/scope-status.md` from scope status files.
- Aggregate `.tasks/product-docs-verify/<run-id>/question-ledger.md` from scope question ledgers.
- Pass `source-delta.md` to every scope-orchestrator when it exists.
- Synthesize final docs only from source evidence, approved scope outputs, approved derivations, unresolved questions, assumptions, and explicit decisions.
- Stop before synthesis if any required scope-orchestrator is missing, interrupted, incomplete, or finishes without reviewer approval.
- Relaunch interrupted, stalled, missing-report, or missing-verdict scope-orchestrators while interruption retry budget remains.
- Never replace a failed scope-orchestrator with direct main-session analysis.
- Never fall back to main-session worker/reviewer execution when nested subagents are unavailable.
- Never invent unsupported requirements or implementation contracts to fill source gaps. Derived product requirements are allowed only when they trace to source behavior and include derivation rationale.

## Scope-Orchestrator Contract

Each scope-orchestrator must:

- Create `.tasks/product-docs-verify/<run-id>/scopes/<scope>/orchestrator.md`.
- Create `.tasks/product-docs-verify/<run-id>/scopes/<scope>/scope-status.md`.
- Create `.tasks/product-docs-verify/<run-id>/scopes/<scope>/question-ledger.md`.
- Spawn a worker subagent for its scope.
- After each worker report, spawn a reviewer subagent for the same scope.
- If the reviewer returns `needs-revision`, send the review findings back to the worker through the same scope-orchestrator and repeat.
- If a worker or reviewer subagent is interrupted, stalls, or fails to write the required report/verdict, spawn a replacement worker/reviewer attempt with the same context and record the interrupted attempt.
- Stop after 3 worker/reviewer attempts unless the user set a different budget.
- Mark the scope as `approved`, `blocked`, or `deferred` in `scope-status.md`.
- Convert unresolved scope issues into blocking or non-blocking questions.
- Use scope-prefixed question ids, for example `Q-ROLE-001`, `Q-ACTOR-001`, or `Q-DOMAIN-001`.
- Report interrupted execution as `blocked` or `interrupted`; do not convert it into partial approval.
- Consolidate missing source artifacts into one source-gap question per artifact class.

Each scope-orchestrator must not:

- Treat a worker report as true without reviewer approval.
- Drop missing information because it is inconvenient.
- Present derived roles, permissions, lifecycle states, data fields, acceptance criteria, or edge cases as direct facts without source evidence and derivation rationale.
- Invent API details, integration details, implementation contracts, or unrelated product behavior beyond source evidence.
- Split one missing source artifact into many speculative endpoint, payload, state, or acceptance questions.
- Ask the main session to manually run each scoped subagent.

## Scope-Orchestrator Prompt Template

Use this template from the main session:

```text
Use the verify-product-docs scope-orchestrator role: <SCOPE_NAME>.

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

Source delta if present:
<SOURCE_DELTA_SUMMARY>

Rules:
- Spawn worker and reviewer subagents internally for this scope only.
- Keep raw source docs read-only.
- Write orchestration artifacts only under `.tasks/product-docs-verify/<RUN_ID>/scopes/<SCOPE_NAME>/`.
- Do not edit staging or final verified product docs.
- Maintain this scope's question-ledger.md from the first missing-info finding.
- Do not claim this scope is approved until its reviewer approves it.
- If nested subagents are unavailable, mark this scope blocked and explain why; do not ask the main session to run worker/reviewer for you.
```

## Retry And Relaunch Policy

Default budgets unless the user sets another value:

- `REVIEW_BUDGET=3` complete worker/reviewer cycles per scope.
- `INTERRUPTION_RETRY_BUDGET=3` controller relaunches per scope-orchestrator for interrupted, stalled, missing-report, or missing-verdict runs.

Definitions:

- `needs-revision`: a reviewer completed and requested changes. Relaunch the same scope-orchestrator with reviewer findings while review budget remains.
- `interrupted`: execution stopped before the scope-orchestrator wrote a final approved/blocked/deferred `scope-status.md`.
- `stalled`: execution produced no new required artifact or final status after the controller's wait window.
- `missing-report`: worker attempt path is absent or incomplete.
- `missing-verdict`: reviewer attempt path is absent, incomplete, or lacks `approved`, `needs-revision`, or `blocked`.

Controller relaunch rules:

1. Preserve approved scopes.
2. For any non-approved required scope, inspect `scope-status.md`, worker attempts, reviewer attempts, and question ledger.
3. If status is `needs-revision`, relaunch that scope-orchestrator with prior attempts and reviewer findings.
4. If status is `interrupted`, `stalled`, missing, or lacks a final reviewer verdict, relaunch that scope-orchestrator with the same `RUN_ID` and prior artifacts.
5. Do not mark the product run blocked until the relevant review or interruption retry budget is exhausted.
6. Do not synthesize final docs from incomplete consistency or scope outputs.

Consistency-specific rules:

- `consistency-reviewer` follows the same retry policy as every Phase 1 scope.
- If consistency attempt 1 returns `needs-revision`, attempt 2 must receive attempt 1 findings and all Phase 1 approved reports.
- If any consistency attempt is interrupted or stalls without a final verdict, the main controller must relaunch `consistency-reviewer` while interruption retry budget remains.
- Final synthesis is allowed only after `consistency-reviewer` reaches `approved`.

## Source Delta And Re-Verification

New product docs and answers to previous open questions do not create a separate workflow. They create a normal new run with a source delta marker.

The main controller writes `.tasks/product-docs-verify/<run-id>/source-delta.md` when any of these are true:

- new source files were added;
- existing source files changed;
- source files were removed;
- previous `docs/product-verified/**` exists and the user is re-verifying;
- previous `.tasks/product-docs-verify/*/question-ledger.md` exists and some questions were answered;
- the user gave answers in the current session.

Minimum source-delta structure:

```text
# Source Delta
## Previous Baseline
## Added Sources
## Changed Sources
## Removed Sources
## Answered Questions
## Notes
```

Rules:

- Scope-orchestrators must read source delta as context before spawning workers.
- Workers must report which delta entries affect their scope.
- Answered questions must keep their original ids when known and move to `answered-by-source`, `answered-by-user`, or `resolved-by-decision`.
- A source delta can resolve blockers, create new contradictions, or supersede prior verified behavior.
- Prior approved reports from older runs are context only; current run approval still requires current scope approval and approved consistency.
- If a delta supersedes previous verified behavior, final synthesis must preserve the old behavior as superseded in traceability or decision log.

## Worker Report Format

Each worker writes to `.tasks/product-docs-verify/<run-id>/scopes/<scope>/worker-attempt-<n>.md`.

```text
# <Scope> Worker Attempt <n>
## Sources Read
## Source Delta Reviewed
## Confirmed Facts
## Contradictions
## Missing Source Artifacts
## Derived Requirements
## Missing Information
## Open Questions Raised
## Edge Cases Or Risks
## Recommended Decisions
## Traceability Candidates
```

## Reviewer Report Format

Each reviewer writes to `.tasks/product-docs-verify/<run-id>/scopes/<scope>/review-attempt-<n>.md`.

```text
# <Scope> Review Attempt <n>
## Verdict
approved | needs-revision | blocked
## Sources Read
## Coverage Check
## Evidence Check
## Invention Check
## Derivation Check
## Source-Gap Consolidation Check
## Missing Or Unsupported Claims
## Contradictions Not Preserved
## Open Questions That Must Be Recorded
## Required Revisions
## Approval Notes
```

Reviewer approval requires:

- Claims cite source files or explicit assumptions/questions.
- Source delta entries that affect the scope are addressed or explicitly marked not applicable.
- Derived roles, permissions, states, data fields, acceptance criteria, and edge cases cite source signals, derivation rationale, and confidence.
- Missing information is written as open questions.
- Missing source artifacts are consolidated instead of expanded into speculative detail questions.
- Contradictions are named, not smoothed over.
- Edge cases and acceptance criteria are observable enough for decomposition.
- The report separates product facts from recommended decisions.

## Question Ledger Format

Each scope-orchestrator must append every missing-info item to `.tasks/product-docs-verify/<run-id>/scopes/<scope>/question-ledger.md` as soon as it appears. The main session later aggregates all scope ledgers into `.tasks/product-docs-verify/<run-id>/question-ledger.md`.

```text
| ID | Scope | Severity | Question | Why It Matters | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Q-ROLE-001 | roles-permissions | blocking | Who can approve projects? | Affects permissions and acceptance criteria. | roles-permissions worker attempt 1 | open | TBD |
```

Statuses: `open`, `answered-by-source`, `answered-by-user`, `resolved-by-decision`, `deferred`.

Any `open` blocking question must appear in `docs/product-verified/open-questions.md`.

## Source-Gap Consolidation

When a whole required artifact class is absent, create one blocker question for that missing source. Do not generate detailed questions that would only make sense after that source exists.

Examples:

- Missing API contract: one `Q-API-001` asking which API contract/source will be used; mark technical continuation blocked for endpoint, payload, error, retry, and auth mapping work.
- Missing authorization policy: one `Q-AUTH-001` asking for the role/ownership policy; mark permission-sensitive flows blocked.
- Missing integration specification: one `Q-INT-001` asking for the integration contract/source; mark integration behavior and failure handling blocked.
- Missing data retention or compliance policy: one `Q-COMP-001` asking for the policy source; mark retention, deletion, audit, and export criteria blocked.

Question ledger entries for source gaps must include:

- missing artifact class;
- impacted scopes;
- why this blocks product decomposition or technical continuation;
- the exact artifact, owner decision, or source needed to unblock.

## Evidence-Constrained Derivation

Some product outputs must be produced even when raw docs are incomplete. These are derived requirements, not unsupported inventions.

Allowed derivations:

- Roles and permissions from described actors, actions, ownership, approvals, visibility, responsibility boundaries, and denied/allowed flows.
- Lifecycle states from documented status words, transitions, queues, approvals, cancellations, retries, archival, or completion conditions.
- Data fields from named entities, forms, workflows, statuses, calculations, validations, reports, imports, exports, and acceptance needs.
- Acceptance criteria from documented or strongly implied behavior; criteria must cover the written behavior and must not add new product behavior.
- Edge cases from documented operations plus standard boundary and failure classes around those operations.

Every derived item must include:

- source reference;
- derivation rationale;
- confidence: `high`, `medium`, or `low`;
- linked open question when confidence is low or the item affects money, identity, authorization, compliance, irreversible transitions, or external contracts.

Do not derive API endpoints, request/response schemas, auth transport, retry policies, rate limits, or error formats when no API/source contract exists. Treat absent API, auth, integration, or compliance artifacts as consolidated source-gap blockers.

## Write Ownership Matrix

| Actor                          | May write                                                                                                                                                                                                                                        | Must not write                                                                                            |
| ------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------- |
| Main controller                | `.tasks/product-docs-verify/<run-id>/main-orchestration.md`, `recovery.md`, aggregate `scope-status.md`, aggregate `question-ledger.md`, `.tasks/product-docs-verify/<run-id>/staging/product-verified/**`, and final `docs/product-verified/**` | Scope-local worker/reviewer attempt files except when explicitly recording a controller intervention      |
| Phase 1 scope-orchestrator     | `.tasks/product-docs-verify/<run-id>/scopes/<scope>/orchestrator.md`, `scope-status.md`, `question-ledger.md`                                                                                                                                    | Any other scope folder, staging folder, aggregate files, `docs/product-verified/**`                       |
| Scope worker                   | `.tasks/product-docs-verify/<run-id>/scopes/<scope>/worker-attempt-<n>.md`                                                                                                                                                                       | Reviewer files, other attempts, other scopes, staging folder, aggregate files, `docs/product-verified/**` |
| Scope reviewer                 | `.tasks/product-docs-verify/<run-id>/scopes/<scope>/review-attempt-<n>.md`                                                                                                                                                                       | Worker files, other scopes, staging folder, aggregate files, `docs/product-verified/**`                   |
| Consistency scope-orchestrator | `.tasks/product-docs-verify/<run-id>/scopes/consistency-reviewer/**`                                                                                                                                                                             | Phase 1 scope folders, staging folder, aggregate files, `docs/product-verified/**`                        |

## Approval Gate

Required before writing `docs/product-verified/**`:

- all Phase 1 scope-orchestrators reached `approved`;
- the consistency scope-orchestrator reached `approved`;
- aggregate `scope-status.md` contains no required scope with any status other than `approved`;
- aggregate `question-ledger.md` contains all blocking questions from scope ledgers;
- no scope was replaced by direct main-session analysis.
- no nested subagent failure was bypassed by controller-run worker/reviewer execution.

If the gate fails:

- write a blocked run report in `.tasks/product-docs-verify/<run-id>/main-orchestration.md`;
- write/update aggregate `scope-status.md` and `question-ledger.md`;
- write/update `.tasks/product-docs-verify/<run-id>/recovery.md` with the next scope-level recovery action;
- stop without writing `docs/product-verified/**`;
- do not call partial output a source of truth.

## Recovery Protocol

Use the same `RUN_ID` when resuming an interrupted run. The main controller should:

1. Read aggregate `scope-status.md`, aggregate `question-ledger.md`, and scope-local `scope-status.md` files.
2. Preserve `approved` scope folders unchanged.
3. Relaunch scope-orchestrators only for missing, interrupted, blocked, or `needs-revision` scopes.
4. Provide each relaunched scope-orchestrator with its previous worker attempts, review attempts, scope status, and question ledger.
5. Re-aggregate statuses and question ledgers after resumed scopes finish.
6. Run `consistency-reviewer` only after all Phase 1 required scopes are approved.
7. Write final `docs/product-verified/**` only after the approval gate passes.

If source docs or the output contract changed since the interrupted run began, record the delta in `recovery.md` and rerun affected approved scopes too.

## Dispatch Phases

Phase 1 can run in parallel:

- `product-scope-reviewer`
- `roles-permissions-reviewer`
- `actor-journey-reviewer`
- `domain-model-reviewer`
- `feature-behavior-reviewer`
- `edge-case-risk-reviewer`
- `acceptance-criteria-reviewer`

Phase 2 runs after Phase 1 aggregation:

- `consistency-reviewer`

Phase 3 is controller-owned synthesis:

- aggregate scope statuses and question ledgers
- write final `docs/product-verified/**`
- validate output

## Scopes

### 1. product-scope-reviewer

Focus on product intent, target segments, value proposition, scope, non-goals, success metrics, and handoff readiness. Flag vague goals, conflicting audiences, and scope creep.

### 2. roles-permissions-reviewer

Focus on extracting and deriving actors, roles, permissions, ownership rules, approval authority, visibility rules, responsibility boundaries, and the permissions matrix. Every derived role or permission must cite source behavior, derivation rationale, and confidence.

### 3. actor-journey-reviewer

Focus on personas, user journeys, happy paths, alternative paths, empty states, and recovery paths. Refer permission questions to the roles/permissions scope instead of owning the permissions matrix.

### 4. domain-model-reviewer

Focus on entities, attributes, identifiers, statuses, lifecycle transitions, relationships, data dictionary terms, invariants, and glossary consistency. Derive data fields only from source-supported entities, forms, workflows, statuses, calculations, validations, reports, imports, exports, and acceptance needs. Flag naming collisions and impossible states.

### 5. feature-behavior-reviewer

Focus on features, workflows, UI/API-visible behavior, validations, notifications, integrations, and business process details. Flag unsupported feature claims and missing behavior.

### 6. edge-case-risk-reviewer

Focus on negative cases, boundary values, duplicates, missing data, concurrency, permissions denial, external dependency failure, retries, rate limits, auditability, privacy, data retention, and migration concerns that follow from documented operations. Do not add unrelated business scenarios.

### 7. acceptance-criteria-reviewer

Focus on deriving observable acceptance criteria from documented or strongly implied behavior, including success, failure, negative, and handoff criteria. Criteria must cover what is written or source-supported and must not add new product behavior.

### 8. consistency-reviewer

Run after the other role reports are available. Focus on cross-report contradictions, duplicate concepts, unresolved terms, decision quality, and whether the final `docs/product-verified` package can become a source of truth.

## Scoped Worker Prompt Template

Each scope-orchestrator uses this template for its worker:

```text
Use the verify-product-docs scoped worker role: <SCOPE_NAME>.

Input folder: <SOURCE>
Run id: <RUN_ID>
Attempt: <N>
Report path: .tasks/product-docs-verify/<RUN_ID>/scopes/<SCOPE_NAME>/worker-attempt-<N>.md

Role focus:
<SCOPE_FOCUS>

Available source files:
<SOURCE_FILE_LIST>

Previous reviewer findings if any:
<REVIEWER_FINDINGS>

Previous interrupted/stalled attempt artifacts if any:
<INTERRUPTED_ATTEMPT_CONTEXT>

Source delta if present:
<SOURCE_DELTA_SUMMARY>

Write the worker report only. Do not edit docs/product-verified.
Do not edit staging skeletons.
Address every source-delta entry that affects this scope. Mark unrelated delta entries as out of scope.
Record missing information as open questions.
Derive roles, permissions, states, data fields, acceptance criteria, and edge cases when source behavior supports them; include source reference, derivation rationale, and confidence.
Do not invent unrelated behavior, API details, integration contracts, or implementation contracts. Consolidate missing artifact classes into source-gap blocker questions.
```

## Scoped Reviewer Prompt Template

Each scope-orchestrator uses this template after each worker attempt:

```text
Use the verify-product-docs scoped reviewer role: <SCOPE_NAME>.

Input folder: <SOURCE>
Run id: <RUN_ID>
Worker report: .tasks/product-docs-verify/<RUN_ID>/scopes/<SCOPE_NAME>/worker-attempt-<N>.md
Review path: .tasks/product-docs-verify/<RUN_ID>/scopes/<SCOPE_NAME>/review-attempt-<N>.md

Review for evidence, missing information capture, contradictions, edge cases, acceptance readiness, and traceability.
Approve derived roles, permissions, states, data fields, acceptance criteria, and edge cases only when they trace to source behavior and include derivation rationale.
Reject reports that invent unrelated behavior, API details, integration contracts, or implementation contracts beyond source evidence.
Reject reports that turn one missing source artifact into many speculative questions.
If the worker report is missing or incomplete, return `blocked` with the missing-report reason so the scope-orchestrator can relaunch the worker within interruption retry budget.

Return one verdict: approved, needs-revision, or blocked.
If not approved, list exact required revisions.
Do not edit docs/product-verified.
```
