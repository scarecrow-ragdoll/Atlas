---
name: plan-backend-waves
description: Use after verify-technical-docs when docs/technical-verified is approved-to-dev and backend work must be decomposed into strict sequential, reviewer-approved, ready-for-dev implementation waves before Jira, Beads, GRACE execution, or developer handoff.
---

<!-- FILE: .agents/skills/plan-backend-waves/SKILL.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Define the project-local Codex workflow for turning approved technical readiness docs into sequential backend implementation waves. -->
<!--   SCOPE: Covers controller gates, wave-by-wave orchestration, reviewer approval, open-question handling, user approval, references, scripts, and final response shape; excludes actual backend implementation and Jira mutation. -->
<!--   DEPENDS: docs/technical-verified, docs/product-verified, .agents/skills/plan-backend-waves/references/output-contract.md, .agents/skills/plan-backend-waves/references/subagent-roles.md, scripts/scaffold_backend_waves.py, scripts/validate_backend_waves.py. -->
<!--   LINKS: M-BACKEND-WAVE-PLANNER / V-M-BACKEND-WAVE-PLANNER / DF-GRACE-CHANGE. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Overview - Establishes the post-technical-approval backend wave decomposition loop. -->
<!--   Non-Negotiables - Defines approved-to-dev, strict sequencing, open-question, reviewer, and user-approval gates. -->
<!--   Workflow - Defines run setup, technical approval validation, wave inventory, per-wave orchestration, synthesis, validation, and pause points. -->
<!--   Wave Approval Loop - Defines how a wave reaches ready-for-dev and how user approval unlocks the next wave. -->
<!--   Final Response - Defines the closeout or pause evidence to report. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.0 - Added backend wave planning workflow skill. -->
<!-- END_CHANGE_SUMMARY -->

# Plan Backend Waves

## Overview

Convert an `approved-to-dev` technical readiness package into a sequential backend wave plan. The main session is the controller: it validates the technical approval gate, dispatches one wave-orchestrator for exactly one wave, waits for reviewer-approved `ready-for-dev`, then stops for explicit user approval before any next wave is planned in detail.

Default inputs are `docs/technical-verified` and `docs/product-verified`. Default output is `docs/backend-waves`. Orchestration reports live under `.tasks/backend-wave-plan/<run-id>/`.

## Non-Negotiables

- Require explicit technical approval before wave planning. `docs/technical-verified/index.md` must contain status `approved-to-dev`, and its open-question ledgers must not contain open `dev-blocking` or `needs-owner-decision` questions.
- Do not bypass the technical gate with chat approval. If technical docs are not approved, stop and route the user back to `$verify-technical-docs`.
- Plan waves strictly sequentially. Only the current wave may be detailed to `ready-for-dev`; later waves may appear only as tentative names, outcomes, dependencies, and risks in the wave map.
- Do not plan wave N+1 until wave N has `ready-for-dev`, has no open wave-blocking questions, has all required reviewer approvals, has been shown to the user, and the user explicitly approves moving on.
- Do not let the user approve a wave with open `wave-blocking` or `needs-owner-decision` questions. Record the attempted approval, list the questions, and keep the wave status `questions-open` or `blocked`.
- Treat more than 8 backend implementation waves as an MVP scope smell. Stop with a re-scope question unless the user explicitly confirms this is not an MVP-sized backend or approves splitting the product into releases.
- Each wave must be ready for development by itself: detailed tasks, dependencies, backend modules, data/API/integration/auth/ops impacts, acceptance criteria, exit criteria, verification plan, rollback/compatibility notes, and source traceability.
- Use multiple reviewer perspectives for every wave: backend architecture, data/API contract, security/integration, testing/delivery, sequencing/MVP, and traceability/consistency.
- Record questions immediately in the wave-local ledger, aggregate ledger, and final `docs/backend-waves/open-questions.md`.
- Do not create or mutate Jira issues unless the user explicitly asks. Produce Jira-ready task text inside the wave artifact.

## Workflow

1. Resolve paths and create a run id.
   - Defaults: `TECHNICAL=docs/technical-verified`, `PRODUCT=docs/product-verified`, `OUTPUT=docs/backend-waves`, `RUN_ID=$(date -u +%Y%m%dT%H%M%SZ)`, `STAGING=.tasks/backend-wave-plan/$RUN_ID/staging/backend-waves`.
   - Use user-provided paths when present.
2. Load the local contracts.
   - Read `references/output-contract.md`.
   - Read `references/subagent-roles.md`.
3. Validate technical approval.
   - Read `$TECHNICAL/index.md`, `$TECHNICAL/open-questions.md`, and `$TECHNICAL/appendix/question-ledger.md` when present.
   - Stop if status is not `approved-to-dev`.
   - Stop if open `dev-blocking` or `needs-owner-decision` technical questions remain.
   - Write the gate result to `.tasks/backend-wave-plan/<run-id>/technical-approval-gate.md`.
4. Inventory inputs.
   - List files with `find "$TECHNICAL" "$PRODUCT" -type f | sort`, ignoring missing optional product source only when the technical package contains enough product traceability.
   - Read prior `docs/backend-waves/index.md`, wave files, and `.tasks/backend-wave-plan/*/question-ledger.md` if this is a continuation.
   - If the user answered prior wave questions in chat, record them in `.tasks/backend-wave-plan/<run-id>/source-delta.md` with original question ids.
5. Scaffold staging.
   - Run `python3 .agents/skills/plan-backend-waves/scripts/scaffold_backend_waves.py --technical-source "$TECHNICAL" --product-source "$PRODUCT" --output "$STAGING"`.
   - Validate staging with `python3 .agents/skills/plan-backend-waves/scripts/validate_backend_waves.py "$STAGING" --allow-placeholders`.
   - Do not write final `docs/backend-waves/**` during staging.
6. Build or update the tentative backend wave map.
   - Use approved technical docs plus product traceability to identify the complete backend scope.
   - Keep later waves intentionally shallow: name, outcome, dependency, rough scope, and risk only.
   - If the tentative map needs more than 8 waves, stop before current-wave planning and ask the re-scope question.
7. Select the current wave.
   - Start with `wave-01` unless resuming.
   - On resume, continue the earliest wave that is not `user-approved`.
   - Never skip a wave unless the user explicitly marks it out of scope and the traceability impact is recorded.
8. Dispatch one wave-orchestrator.
   - Use `references/subagent-roles.md` as the role contract.
   - Give it write scope only for `.tasks/backend-wave-plan/<run-id>/waves/<wave-id>/`.
   - The wave-orchestrator must spawn its own planner worker and required reviewers.
   - The main session must not dispatch wave workers or reviewers directly.
9. Run the wave planner/reviewer loop.
   - The planner writes `planner-attempt-<n>.md`.
   - Each reviewer writes `review-<perspective>-attempt-<n>.md`.
   - The wave-orchestrator repeats until every required reviewer approves, a blocker is recorded, or budgets are exhausted.
   - Maintain wave-local `wave-status.md` and `question-ledger.md` from the first gap.
10. Synthesize the current wave.

- After the wave loop terminates with either reviewer approval or recorded blockers, write or update `docs/backend-waves/**` for the wave map, current wave file, ledgers, traceability, decisions, and reviewer verdicts.
- Set the current wave to `ready-for-dev` only when it has no open wave blockers and all required reviewers approved.
- Set it to `questions-open` or `blocked` when questions, missing artifacts, unavailable reviewers, or exhausted budgets remain; persist the questions and do not ask for approval.

11. Validate and report.

- Run `python3 .agents/skills/plan-backend-waves/scripts/validate_backend_waves.py "$OUTPUT"`.
- Run `git diff --check -- "$OUTPUT" .tasks/backend-wave-plan .agents/skills/plan-backend-waves`.
- If the current wave is `ready-for-dev`, report its overview and ask for explicit user approval to move to the next wave.
- If it has open questions, report the questions and do not ask for approval yet.

## Wave Approval Loop

A wave is `ready-for-dev` only when all of these are true:

- technical approval gate passed for the source package;
- wave scope is source-traceable and does not invent backend behavior;
- all dependencies and exclusions are explicit;
- detailed backend tasks are implementation-ready;
- acceptance criteria and exit criteria are concrete and testable;
- unit, contract, integration, e2e, migration, fixture, and observability checks are assigned where relevant;
- no open `wave-blocking` or `needs-owner-decision` questions remain;
- all required reviewer perspectives approved;
- the final wave file passes validation.

User approval is separate from `ready-for-dev`:

- The user must explicitly approve the current wave after seeing the overview.
- Record approval in the wave file under `## User Approval` and in `appendix/decision-log.md`.
- Only then may the next wave be planned in detail.
- If user approval arrives with unanswered blocking questions, refuse the approval and keep the current wave open.

## Interruption And Recovery

On interruption:

- Update `.tasks/backend-wave-plan/<run-id>/main-orchestration.md` with the current phase, wave id, interrupted actor, and next recovery action.
- Aggregate available wave statuses and question ledgers.
- Create `.tasks/backend-wave-plan/<run-id>/recovery.md`.
- Do not mark a partial wave `ready-for-dev`.

To recover:

1. Resume with the same `RUN_ID`.
2. Read `main-orchestration.md`, aggregate `wave-status.md`, aggregate `question-ledger.md`, and the current wave folder.
3. Preserve approved reviewer verdicts unless sources, answers, or wave scope changed.
4. Relaunch only the current wave-orchestrator for missing, interrupted, blocked, or needs-revision work.
5. Continue to require user approval before the next wave.

Default budgets unless the user sets another value:

- `REVIEW_BUDGET=3` complete planner/reviewer cycles per wave.
- `INTERRUPTION_RETRY_BUDGET=3` controller relaunches per interrupted, stalled, missing-report, or missing-verdict wave-orchestrator.

## Final Response

Report:

- Technical approval gate result.
- Output folder path and orchestration report folder path.
- Current wave id, status, reviewer verdict summary, acceptance criteria count, exit criteria count, and open question count.
- What will be ready after the wave if implemented.
- Whether the user may approve this wave now.
- Validation commands and results.
