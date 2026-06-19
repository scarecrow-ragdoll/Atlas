---
name: plan-to-beads
description: Use when the user asks to transform an implementation plan, architecture document, research plan, PRD, or verification plan into Beads milestone/epic/task graphs. Focus on minimizing semantic loss, preserving source-plan meaning, strict acceptance criteria, and post-implementation test/QA closure. Assumes Beads already exists in the project; do not explain Beads basics.
---

<!-- FILE: .agents/skills/plan-to-beads/SKILL.md -->
<!-- VERSION: 1.1.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Define the project-local Codex workflow for converting source plans into Beads milestone, epic, and task graphs without semantic loss. -->
<!--   SCOPE: Covers source anchoring, requirement mapping, mandatory milestone and child-epic structure, delivery contours, acceptance criteria, blocker handling, graph checks, and final reporting; excludes Beads basics and project-specific Beads convention design. -->
<!--   DEPENDS: Beads CLI and project Beads state, source implementation or verification plans, .agents/skills/plan-to-beads/agents/openai.yaml. -->
<!--   LINKS: M-PLAN-TO-BEADS-SKILL / V-M-PLAN-TO-BEADS-SKILL / DF-GRACE-CHANGE. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Semantic-Loss Principle - Keeps source plans as the source of truth for all generated Beads work. -->
<!--   Required Milestone And Three-Epic Shape - Defines the mandatory non-executable milestone and child execution epics. -->
<!--   Responsibility Boundary And Delivery Mode - Defines closure contours and blocker ownership. -->
<!--   Task Writing Rules - Preserves common objective, acceptance criteria, source pointers, and agent freedom. -->
<!--   Verification Checklist And Final Response - Defines graph validation and closeout reporting. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.1.0 - Updated the Beads graph shape to require a non-executable milestone with three child epics. -->
<!-- END_CHANGE_SUMMARY -->

# Plan To Beads

Convert a plan into an executable Beads graph without losing requirements.

This skill is not a Beads tutorial and not a replacement for project Beads conventions. Assume project Beads workflow, fields, and command syntax are already known.

The job is narrower: shape Beads tasks so agents follow the source plan, keep the larger goal in view, and do not drift while still having freedom to solve the work well.

## Semantic-Loss Principle

Any transformation of information loses some meaning. Treat plan-to-Beads conversion as a risky semantic migration, not a harmless summary.

The goal is not to compress the plan. The goal is to create execution handles that keep pointing back to the original plan and preserve enough context for each worker to understand why their slice exists.

Beads tasks, matrices, and coverage ledgers are navigation aids. They are not a second source of truth. If a task summary and the source plan disagree, the source plan wins unless the task explicitly records an approved plan change.

Rules:

- Do not replace the source plan with Beads wording.
- Keep the source plan path/spec id in every epic and task.
- Repeat the common objective in every task so an agent working on one slice sees the larger goal.
- Preserve exact names and hard constraints even when they feel repetitive.
- Use matrices to prevent omissions, not to rewrite the plan into a new plan.
- If the source plan or user intent is unclear from context, ask for the missing plan/source before creating the graph.
- If a conversion step would require interpretation beyond the source, mark it as a question or discovery task rather than silently inventing scope.

## Core Promise

The Beads graph must be strong enough that another agent can implement the whole plan from the tasks without reading the chat history.

Every unique requirement from the source plan must appear in at least one of:

- epic or task scope
- acceptance criteria
- design/notes
- explicit non-goal
- explicit blocker/follow-up task

If a requirement cannot be mapped, stop and add a task or note before claiming the graph is complete.

Prefer source anchoring over paraphrase. Each task should tell the worker what plan to open and what overall outcome they are contributing to.

## Agent Freedom And Guardrails

Tasks should guide the agent's thinking, not micromanage every edit.

Write tasks so they constrain:

- the source plan to follow
- the desired final invariant
- the hard requirements and non-goals
- the proof needed to close the task

Leave freedom in:

- exact implementation details
- local refactoring choices that preserve the plan
- test structure when the proof remains equivalent or stronger
- safe problem-solving steps for unexpected blockers

Good task wording says: "make this outcome true and prove it." It should not become a brittle mechanical checklist unless the source plan requires exact steps.

Each task should nudge the worker to:

- open the source plan before editing
- inspect the current code/runtime path for its slice
- follow existing project patterns before inventing new abstractions
- preserve current behavior that the plan says must not regress
- record gaps as findings/follow-ups instead of silently widening scope

This keeps agents autonomous while keeping them inside the plan.

## Preserve Intent And Failure Mode

When a plan requirement exists because of a risk, preserve that risk in the task wording.

Do not only say what to build. Say what must not happen.

Examples:

- Prefer: "Implement safe relation replacement so a conclusion never temporarily loses the old stored document relation if StoreApp/session creation fails."
- Avoid: "Implement relation replacement."

Guidelines:

- For migration tasks, include the preserved current behavior or no-regression invariant.
- For reliability tasks, include the failure mode being prevented.
- For test tasks, include the regression or edge case the test must catch.
- For QA tasks, include the drift, hidden bug, or contract mismatch being hunted.

This is a low-cost way to carry the original reasoning into each executable slice.

## Required Milestone And Three-Epic Shape

For a serious implementation plan, always create or normalize exactly one top-level non-executable `milestone` plus exactly these three execution epics unless the user explicitly asks for a different structure.

The milestone is mandatory in the supported Beads baseline. Create it directly. If milestone creation fails, treat that as a broken local Beads environment and report the tooling mismatch; do not silently fake the same meaning with an epic, task, label, or note.

The required hierarchy is:

1. **Milestone**
   - Represents the whole plan, release, phase, or work package.
   - Contains no direct implementation work.
   - Preserves the source plan pointer, common objective, scope boundary, non-goals, and final readiness definition.
   - Must not have executable implementation tasks directly under it.

2. **Implementation Epic**
   - Implements the plan strictly from start to finish.
   - Covers code, migrations, contracts, docs/artifacts, cleanup, and focused implementation verification.
   - Must not stop at a partial migration, compatibility layer, or “happy path only” state.
   - Must be attached to the milestone with `--parent <milestone-id>`.

3. **Full Test Coverage Epic**
   - Runs after the implementation epic.
   - Means practically complete coverage of the target area, not a vague “add tests” task.
   - Includes test environment/data preparation, deterministic fixtures, E2E/integration/browser tests, edge cases, retries, failure paths, download/regression proof, and CI/local command wiring.
   - Interactive browser/MCP checks are exploratory only; they do not replace committed test files and runnable commands unless the user explicitly allows that.
   - Must be attached to the milestone with `--parent <milestone-id>`.

4. **Pre-MR QA Epic**
   - Runs after the test coverage epic.
   - Independently audits implementation, tests, docs, API contracts, project artifacts, and runtime behavior before rollout/MR.
   - Must include project-appropriate gates such as typecheck, build, lint/static checks, focused tests, E2E commands, and documentation/artifact consistency.
   - Must end with a pushed branch plus MR/readiness packet when the user requested MR delivery, or an explicit blocker list when MR readiness cannot honestly be reached.
   - Must be attached to the milestone with `--parent <milestone-id>`.

Do not write milestone acceptance criteria as if the milestone were executable work. Keep `bd epic status` and `bd epic close-eligible` semantics centered on the child epics, not on the milestone.

Use this creation pattern:

```bash
milestone_id=$(bd create --type=milestone --title "<source-plan work package>" --silent)
impl_id=$(bd create --type=epic --parent "$milestone_id" --title "[epic] <plan>: implementation" --silent)
coverage_id=$(bd create --type=epic --parent "$milestone_id" --title "[epic] <plan>: full test coverage" --silent)
qa_id=$(bd create --type=epic --parent "$milestone_id" --title "[epic] <plan>: pre-MR QA / readiness" --silent)

bd dep "$impl_id" --blocks "$coverage_id"
bd dep "$coverage_id" --blocks "$qa_id"

bd create --type=task --parent "$impl_id" --title "<implementation slice>"
bd create --type=task --parent "$coverage_id" --title "<coverage slice>"
bd create --type=task --parent "$qa_id" --title "<QA slice>"
```

Tasks belong under the relevant epic, never directly under the milestone.

Do not mark work complete just because tasks were created.

## Responsibility Boundary And Delivery Mode

Before creating QA/readiness tasks, record the closure contour the user expects:

- `implementation-ready`: code and focused implementation proof only; no MR claim.
- `MR-ready`: code, tests, business-logic verification, branch push, MR creation/update, and handoff risks are complete.
- `staging-verified`: MR-ready plus deployed staging/test-contour runtime proof.
- `production-verified`: staging-verified plus production rollout proof.

Default to `MR-ready` when the user asks to make an MR, prepare for merge, or verify quality before stand rollout. Do not make Pre-MR QA epics wait for external deployment, target database mutation, unavailable CI runners, or operator-only runtime actions unless the user explicitly sets the contour to `staging-verified`, `production-verified`, or `rollout-owned`.

Use a delivery-mode field in milestone/epic notes:

- `no-push`: local graph or local implementation only.
- `branch-pushed`: branch delivery is in scope, MR is not.
- `MR-created`: branch push and MR creation/update are in scope.
- `rollout-owned`: the agent also owns deployment and target runtime verification.

Classify unresolved work by ownership:

- `code/logic blocker`: must block implementation, coverage, and QA closure.
- `verification blocker`: blocks closure only when it prevents the selected closure contour.
- `external infra blocker`: record in the MR/readiness packet and handoff notes; do not block `MR-ready` closure when code quality and local/CI-accessible checks are complete.
- `post-MR operator handoff`: create or link a follow-up task/runbook; do not block Pre-MR QA unless rollout is explicitly agent-owned.

## Input Handling

Read the source material first:

- architecture/research document
- implementation plan
- acceptance criteria
- open questions and decisions
- existing Beads state, if any, to avoid duplicates
- project governance or artifact rules when relevant

If no source document/path is available and the chat context is not enough to reconstruct the plan safely, ask the user for the source. Do not generate a confident Beads graph from a blurry memory of the plan.

Preserve exact domain terms:

- endpoint names
- use-case names
- field names
- statuses
- department names
- file paths
- “must not” rules
- non-goals
- failure modes

Do not paraphrase away a hard constraint. A phrase like “must not keep legacy fallback” belongs in acceptance criteria, not only in a loose description.

## Requirement Extraction

Before creating tasks, build an internal requirement ledger:

- business operations and user flows
- server/API contracts
- storage/database/migration changes
- frontend/browser behavior
- security and authorization rules
- idempotency, retry, and partial-failure rules
- compatibility and cleanup rules
- verification commands and evidence expectations
- docs/artifacts that must be updated
- explicit out-of-scope items
- unknowns/blockers

For each important item, map it across the three-stage graph:

```text
source requirement -> implementation task -> test coverage task -> QA/release-readiness task
```

This ledger may stay internal for small plans. For large or risky plans, create a coverage/traceability artifact or a dedicated first task that owns it.

The ledger is a cross-check, not a new authority. Keep it short and use it to detect lost requirements before creating tasks.

Then group by executable ownership, not by document headings alone.

Good task boundaries:

- one coherent runtime slice
- one primary ownership area
- enough acceptance detail for a worker to finish without guessing
- not so small that a worker only edits one trivial field
- not so broad that failures become invisible

## Task Content Contract

Every epic and task should contain:

- the common objective in plain language
- link/path to the source plan
- scope
- hard requirements
- non-goals or exclusions when relevant
- acceptance criteria with proof expectations

Task text should usually include a short instruction like: "Start by re-reading the source plan section and inspecting the current implementation for this slice." This prevents execution from drifting away from the real code and source plan.

Prefer outcome/invariant phrasing over long procedural lists. A task may include suggested order, but its closure should depend on the final state and proof, not on following a guessed micro-plan.

For implementation tasks, acceptance must mention behavior and verification, not just files changed.

For test tasks, acceptance must mention test files, commands, fixtures/data, and edge cases.

For QA tasks, acceptance must mention findings handling: every issue is either fixed, converted into a follow-up Bead, explicitly accepted, or recorded as a blocker.

Do not make a task depend on memory of the conversation. If a worker needs context, put it in task text or point to the exact plan section/file.

Do not duplicate all Beads formatting rules inside task text. Use task text for plan context, boundaries, invariants, and proof.

When tasks mention blockers, phrase them as problems to solve, not as automatic stop signs. A worker should understand that they are expected to investigate and resolve ordinary blockers before escalating.

## Blocker Handling Posture

A blocker is a real problem in the work, not automatically a request for user help.

Encode this posture in implementation, test, and QA tasks:

- First understand the blocker: reproduce it, isolate the layer, inspect logs/state/config, and compare against the source plan.
- Try safe and correct remedies before escalating: configuration alignment, deterministic setup/reset, focused fixes, small harness improvements, clearer fixtures, targeted retries, or splitting an oversized step into a precise follow-up.
- Use engineering judgment and ingenuity, but do not use hacks that weaken the plan, bypass security, hide failures, or reduce required coverage.
- Treat resource constraints as design constraints: use available CPU/memory/storage consciously, avoid saturating the machine by default, and prefer bounded commands over all-core/all-container fanout unless the project expects it.
- Only escalate to the user when the blocker cannot be resolved safely after serious attempts, requires a product decision, requires unavailable credentials/external access, would need destructive/risky action, or has no honest path forward.
- When escalating, include what was tried, what evidence was found, why the safe paths are exhausted, and the smallest decision or access needed.

The task graph should not encourage early handoff. It should encourage autonomous problem solving with honest blocker reporting only when the blocker is truly hard.

## Implementation Epic Pattern

Create tasks that cover the whole plan in implementation order.

Typical slices:

- foundation/storage/database
- backend API and domain use cases
- cross-service clients/contracts
- boundary renames or module normalization
- business use-case migrations
- frontend/client behavior
- admin/secondary surfaces
- legacy removal and project artifacts
- final implementation verification

Implementation tasks must include cleanup and artifact updates when those are part of the plan. Do not leave docs/GRACE/verification artifacts for an unspecified future pass.

## Full Test Coverage Epic Pattern

This epic is not optional when the user asks for “100%” or “fully covered” behavior.

Use a clear sequence:

1. Coverage matrix.
2. Test environment and deterministic fixture/data setup.
3. Service/backend E2E suites.
4. Application/server E2E suites.
5. Domain/business-flow E2E suites.
6. Browser/client E2E suites from committed test files.
7. Cross-service edge cases, retries, idempotency, partial failures.
8. Regression/no-legacy/CI gate wiring.

“100%” means every in-scope operation and meaningful edge case is covered or explicitly blocked. It does not mean inventing fake line-coverage claims.

When files are uploaded in tests, require deterministic test documents in known fixture folders or generated by test setup. The test should prove the actual file path and file bytes flow through the intended runtime path.

Test environment and infrastructure rules:

- Prefer a reproducible test environment over ad hoc local state.
- Prepare database/storage fixtures deliberately and reset them deterministically.
- Use machine resources intentionally, but do not saturate the developer machine by default. Avoid all-core/all-container fanout unless the project explicitly allows it.
- If infrastructure is broken, first reason through and try safe local fixes: env alignment, ports, containers, seeds, permissions, storage buckets, migrations, and test data.
- If the test environment cannot be prepared after serious safe attempts, or only risky/unclear changes remain, stop, explain the hard blocker, and ask the user how to proceed.
- Do not hide infrastructure failures by weakening E2E requirements.
- Coverage tasks may add test harnesses and fixtures. They must not quietly change product architecture to make tests pass; architecture-impacting fixes need explicit implementation/follow-up tasks.

## Pre-MR QA Epic Pattern

This epic challenges the result after implementation and tests.

Typical tasks:

- traceability audit against source plan and coverage matrix
- backend/storage/security review
- server/domain/business review
- frontend/UX/runtime-boundary review
- docs/project-artifact/API-contract consistency review
- regression and exploratory failure sweep in the prepared test environment
- final MR/readiness packet

QA must include project gates:

- typecheck or equivalent
- build or equivalent
- lint/static checks when used by the project
- focused unit/integration/E2E commands
- docs/artifact consistency checks

QA is not allowed to hide unresolved issues inside the final summary. File or create follow-up tasks for findings.

The final QA task should produce a release/MR readiness packet with:

- what was implemented
- source plan used
- commands run and results
- typecheck/build/lint status or project equivalents
- E2E/runtime evidence
- no-regression proof for critical flows
- docs/artifacts updated
- remaining blockers or accepted risks
- explicit note that nothing was pushed unless the user requested it

## Dependency Semantics

Dependencies should express when work can safely start:

- Test coverage starts after implementation is complete enough to test.
- QA starts after the E2E/test coverage epic is complete.
- Within an epic, sequence tasks when later tasks require earlier code, data, or evidence.
- Keep independent tasks parallelizable only when their write scopes and proof scopes are truly independent.

After creating or updating the graph, verify:

- exactly one top-level milestone represents the source plan, unless an equivalent milestone already existed
- all children are under the expected epic
- the three required epics are attached to the milestone
- dependencies have no cycles
- ready tasks make sense
- representative tasks show the source plan link and common objective
- no duplicate stale Beads contradict the new graph

## No Information Loss Checklist

Before final answer, check:

- Every plan phase has a task.
- Every “must” and “must not” appears in acceptance or non-goals.
- Every edge case appears in implementation, test, or QA work.
- Every artifact/doc update appears in a task.
- Every out-of-scope item is explicitly named.
- Every known blocker is represented.
- Blocker wording tells agents to solve safe/ordinary blockers first and escalate only true hard blockers with evidence.
- Every task has the source plan path/spec id or an equivalent explicit source pointer.
- Every task repeats the common objective or enough parent context to prevent isolated-slice confusion.
- Every important requirement maps to implementation, test coverage, and QA/readiness or has an explicit reason why not.
- Final test and QA epics cannot become ready before implementation closure.
- The milestone has no direct executable implementation tasks.
- No task asks for push/deploy unless the user explicitly requested it.

If a source plan is too vague to decompose honestly, create a preliminary discovery/coverage-matrix task instead of inventing missing requirements.

## Beads Graph Definition Of Done

Before reporting success, verify the graph itself:

- the top-level milestone exists and is non-executable
- all expected epics exist
- all expected epics are children of the milestone
- every child is under the intended epic
- dependencies have no cycles
- ready tasks match the intended first executable tasks
- representative tasks from each epic show the source plan, common objective, scope, and acceptance
- no stale duplicate task contradicts the new graph
- local Beads export/sync warnings are reported separately from task creation success

Also check that the graph has not become more bureaucratic than useful. If a task exists only to satisfy process and does not protect meaning, sequencing, proof, or execution clarity, merge or remove it.

## Anti-Patterns

Avoid:

- one vague “QA” task
- one vague “add tests” task
- putting implementation tasks directly under the milestone
- using an epic as a fake milestone when the tooling baseline is mismatched
- adding conditional milestone fallback phrasing when the agreed Beads baseline requires milestone
- a Beads graph that can close while known blockers remain
- dropping non-goals because they are not implementation work
- turning “100% coverage” into fake line-coverage claims
- hiding infrastructure problems by reducing test scope
- leaving typecheck/build/lint to an unspecified future step
- replacing source-plan wording with a simplified task summary when exact contract names matter
- turning a plan into a matrix, then treating the matrix as more authoritative than the plan
- over-constraining implementation details when the plan only requires an outcome and proof

## Final Response

Return:

- milestone ID and title
- epic IDs and titles
- child task IDs grouped by epic
- dependency summary
- what proof was run on the Beads graph
- any warnings, such as duplicate existing tasks or local Beads export issues
- short prompt-goal for a new session that names the milestone ID

Keep the response concise but concrete. The user should be able to start from the first ready task immediately.

The final prompt-goal must be ready to paste into a new session. Use this short shape and replace placeholders with real IDs:

```text
Выполни полностью Beads milestone <milestone-id>: закрыть на 100% все его epic-и и все дочерние задачи. Начни с `bd show <milestone-id> --children` / `bd children <milestone-id> --pretty`, затем последовательно выполняй child epic-и milestone: implementation, full test coverage, pre-MR QA/readiness. Распределяй работу агентам по независимым задачам, но сохраняй единый контроль качества и зависимости.

При блокерах сначала логически разбери причину, проверь код/контракты/тесты/окружение и постарайся исправить блокер сам; к пользователю выноси только настоящий продуктовый, доступовый или рискованный blocker с доказательствами. Выполняй milestone без костылей, без legacy/fallback как финального состояния, без ослабления требований и без вреда продукту: следуй source plan, существующим архитектурным решениям, тестам, GRACE/документации и обязательной верификации.
```
