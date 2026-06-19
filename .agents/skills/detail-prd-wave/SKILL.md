---
name: detail-prd-wave
description: Use after decompose-prd-waves when one selected backend wave from docs/prd-waves, such as WAVE-01, must be detailed into a code-aware, reviewer-approved, ready-for-dev backend wave brief with acceptance criteria, exit criteria, backend implementation slices, verification obligations, open-question ledgers, and explicit fit checks against the codebase, neighboring backend waves, and the separate frontend-pages dependency context.
---

<!-- FILE: .agents/skills/detail-prd-wave/SKILL.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Define the project-local Codex workflow for detailing exactly one approved backend PRD wave into a ready-for-dev backend wave brief. -->
<!--   SCOPE: Covers source gates, selected backend-wave isolation, codebase and neighboring-backend-wave fit, frontend-pages dependency context, multi-subagent planning and review, output contracts, scripts, validation, and final response shape; excludes implementation, Beads/Jira mutation, frontend planning, and detailing future waves. -->
<!--   DEPENDS: docs/prd-waves, docs/product-verified, docs/technical-verified, docs/development-plan.xml, docs/knowledge-graph.xml, docs/verification-plan.xml, .agents/skills/detail-prd-wave/references/output-contract.md, .agents/skills/detail-prd-wave/references/subagent-roles.md, scripts/scaffold_detail_prd_wave.py, scripts/validate_detail_prd_wave.py. -->
<!--   LINKS: M-PRD-WAVE-DETAILER / V-M-PRD-WAVE-DETAILER / DF-GRACE-CHANGE. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Overview - Establishes the selected-wave deep planning loop. -->
<!--   Non-Negotiables - Defines one-wave focus, approved source gates, code-aware fit, reviewer, and blocker rules. -->
<!--   Workflow - Defines run setup, inventory, orchestration, synthesis, validation, and approval pause points. -->
<!--   Detailed Wave Rules - Defines what a ready-for-dev wave brief must and must not contain. -->
<!--   Final Response - Defines closeout evidence to report. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.2 - Clarified selected waves are backend-only and frontend stays under frontend-pages/. -->
<!-- END_CHANGE_SUMMARY -->

# Detail Backend PRD Wave

## Overview

Turn exactly one approved backend wave from `docs/prd-waves` into a detailed, code-aware, reviewer-approved ready-for-dev backend brief. The main session is the controller: it selects one backend `WAVE-<nn>`, inventories source docs, neighboring backend waves, the separate `frontend-pages/**` dependency context, GRACE contracts, and relevant code, dispatches a wave-orchestrator, requires multiple independent reviewer perspectives, then writes a detailed package only for the selected backend wave.

Default inputs are `docs/prd-waves`, `docs/product-verified`, and `docs/technical-verified` when present. Default output is `docs/prd-wave-details`. Orchestration reports live under `.tasks/prd-wave-detail/<run-id>/`.

## Non-Negotiables

- Work on exactly one backend wave per run. Require a selected `WAVE_ID` from the user or from the earliest backend wave that is not already detailed; never detail wave N+1 while wave N is still open unless the user explicitly skips or supersedes wave N.
- Start from an approved shallow backend source wave. `docs/prd-waves` must exist, the selected backend wave must exist, and the selected wave must be `top-level-ready` or `user-approved` with no open decomposition-blocking or owner-decision questions affecting it. If the backend wave map is still awaiting approval, ask for explicit user approval before deep planning.
- Do not re-decompose the PRD. If the selected wave boundary is incoherent, record a blocker and route back to `$decompose-prd-waves`.
- Keep frontend separate. Read `docs/prd-waves/frontend-pages/**` only to understand backend dependencies, page consumers, and sequencing constraints. Do not detail pages, routes, screens, navigation, UX states, components, frontend tests, or frontend implementation.
- Read surrounding context. The detailed backend wave must check neighboring backend waves, prior detailed backend waves, `docs/development-plan.xml`, `docs/knowledge-graph.xml`, `docs/verification-plan.xml`, and relevant source files so implementation slices fit the existing codebase and do not steal scope from other backend waves.
- Use subagents. The main session must dispatch exactly one wave-orchestrator for the selected wave. The wave-orchestrator must spawn specialist planners and required reviewers internally.
- Require review. A wave may be `ready-for-dev` only when every required reviewer perspective approves and final fit review approves. If nested subagents or reviewers are unavailable, mark the wave blocked.
- Preserve open questions. Record wave-blocking, owner-decision, deferred, and watchlist questions in wave-local and aggregate ledgers from the first gap.
- Do not implement code, mutate Beads/Jira, or edit source modules. This skill may read code and produce developer-ready backend task text, acceptance criteria, exit criteria, verification obligations, and handoff packets only.
- Do not hide scope creep. If the selected backend wave cannot be independently implemented without pulling material scope from future backend waves or frontend page planning, stop with a sequencing question or route back to the wave map.

## Workflow

1. Resolve paths and create a run id.
   - Defaults: `PRD_WAVES=docs/prd-waves`, `PRODUCT=docs/product-verified`, `TECHNICAL=docs/technical-verified`, `OUTPUT=docs/prd-wave-details`, `RUN_ID=$(date -u +%Y%m%dT%H%M%SZ)`, `STAGING=.tasks/prd-wave-detail/$RUN_ID/staging/prd-wave-details`.
   - Require `WAVE_ID` such as `WAVE-01`; normalize `wave-01`, `1`, and `01` to `WAVE-01`.
   - Use user-provided paths when present.
2. Load local contracts.
   - Read `references/output-contract.md`.
   - Read `references/subagent-roles.md`.
3. Validate the source wave gate.
   - Read `docs/prd-waves/index.md`, `wave-map.md`, `frontend-pages/index.md`, relevant `frontend-pages/page-<nnn>.md` files, `open-questions.md`, `waves/index.md`, `waves/wave-<nn>.md`, and `appendix/question-ledger.md`.
   - Stop if the selected backend wave is missing, not shallow-approved, or has open decomposition-blocking or owner-decision questions.
   - Write `.tasks/prd-wave-detail/<run-id>/source-wave-gate.md`.
   - When promoting final docs, mirror the gate in `docs/prd-wave-details/index.md` `## Source Wave Gate` with `source-wave-gate: passed` for ready waves or `source-wave-gate: blocked` plus a matching open question row for blocked waves.
4. Inventory context.
   - List verified docs with `{ find "$PRD_WAVES" -type f; find "$PRODUCT" -type f; if [ -d "$TECHNICAL" ]; then find "$TECHNICAL" -type f; fi; } | sort`.
   - Read prior `docs/prd-wave-details/**` and `.tasks/prd-wave-detail/*/question-ledger.md` when resuming or when earlier backend waves may affect this one.
   - Read `docs/development-plan.xml`, `docs/knowledge-graph.xml`, and `docs/verification-plan.xml` for module and verification references.
   - Inspect relevant code paths named by the source wave, technical docs, GRACE graph, or repository search. Keep this read-only.
   - Write `.tasks/prd-wave-detail/<run-id>/context-inventory.md`.
5. Scaffold staging.
   - Run `python3 .agents/skills/detail-prd-wave/scripts/scaffold_detail_prd_wave.py --wave-id "$WAVE_ID" --prd-waves "$PRD_WAVES" --product-source "$PRODUCT" --technical-source "$TECHNICAL" --output "$STAGING"`.
   - If GRACE docs live outside the defaults, pass `--development-plan`, `--knowledge-graph`, and `--verification-plan`.
   - Validate staging with `python3 .agents/skills/detail-prd-wave/scripts/validate_detail_prd_wave.py "$STAGING" --wave-id "$WAVE_ID" --allow-placeholders`.
   - Do not write final `docs/prd-wave-details/**` during staging.
6. Dispatch the selected wave-orchestrator.
   - Use `references/subagent-roles.md` as the role contract.
   - Give it write scope only for `.tasks/prd-wave-detail/<run-id>/waves/<wave-id>/`.
   - The wave-orchestrator must spawn specialist planners and reviewers internally.
   - The main session must not dispatch planner workers or reviewers directly.
7. Run the planner/reviewer loop.
   - Planner specialists write their reports under `.tasks/prd-wave-detail/<run-id>/waves/<wave-id>/planner-<scope>-attempt-<n>.md`.
   - Reviewers write `.tasks/prd-wave-detail/<run-id>/waves/<wave-id>/review-<perspective>-attempt-<n>.md`.
   - The wave-orchestrator repeats until every required reviewer approves, a blocker is recorded, or budgets are exhausted.
   - Maintain wave-local `wave-status.md` and `question-ledger.md`.
8. Synthesize the selected wave package.
   - Aggregate planner reports, reviewer verdicts, source-wave evidence, codebase fit evidence, other-wave fit evidence, decisions, and questions.
   - Write or update only selected-wave files under staging.
   - Do not add detailed wave files for later waves.
9. Run final fit review.
   - Dispatch `final-wave-fit-review` only after the selected-wave candidate package exists.
   - Pass the concrete candidate package path to the reviewer.
   - If the reviewer returns `needs-revision`, revise only the selected-wave candidate and rerun final review.
10. Promote and validate.

- Promote the reviewed candidate into `docs/prd-wave-details/**` only after final fit review returns `approved` or after blockers are fully recorded.
- Set package status to `questions-open`, `blocked`, `wave-ready-awaiting-user-approval`, `wave-approved`, or `superseded` using the output contract.
- Run `python3 .agents/skills/detail-prd-wave/scripts/validate_detail_prd_wave.py "$OUTPUT" --wave-id "$WAVE_ID"`.
- Run `git diff --check -- "$OUTPUT" .tasks/prd-wave-detail .agents/skills/detail-prd-wave`.
- If `$OUTPUT` is outside the repository during a fixture or dry run, use `git diff --check --no-index /dev/null "$OUTPUT"` or another no-index whitespace check instead of treating the path as a repo diff path.
- If the selected wave is `ready-for-dev`, report the overview and ask for explicit user approval before downstream Beads, Jira, GRACE execution, or implementation.

## Detailed Wave Rules

A selected backend wave brief must include:

- selected wave id, source wave path, status, and user approval state;
- outcome after implementation and explicitly excluded scope;
- compatibility notes for prior, current, and future backend waves;
- read-only dependency notes from `frontend-pages/**`, without frontend planning;
- codebase fit: relevant modules, files, public contracts, generated artifacts, integration points, and likely GRACE graph deltas;
- backend implementation slices with stable ids;
- acceptance criteria with stable `AC-W<nn>-...` ids;
- exit criteria with stable `EC-W<nn>-...` ids;
- verification obligations with stable `TEST-W<nn>-...` ids and commands or required evidence;
- security, privacy, data, API, integration, operations, rollout, rollback, and compatibility notes where relevant;
- open questions and traceability to source docs, wave map entries, code references, and reviewer reports.

A selected backend wave brief must not include:

- code changes or generated artifacts;
- final implementation decisions invented without source, code, explicit user decision, or reviewer support;
- detailed plans for later backend waves;
- frontend pages, screens, routes, navigation, UX states, component architecture, visual design, copy decks, frontend tests, or frontend implementation tasks;
- Jira or Beads mutations;
- broad repository refactors outside the selected wave.

## Interruption And Recovery

On interruption:

- Update `.tasks/prd-wave-detail/<run-id>/main-orchestration.md` with current phase, selected wave, interrupted actor, and next recovery action.
- Aggregate available wave status and question ledgers.
- Create `.tasks/prd-wave-detail/<run-id>/recovery.md`.
- Do not mark a partial wave `ready-for-dev`.

To recover:

1. Resume with the same `RUN_ID` and `WAVE_ID`.
2. Read `main-orchestration.md`, `source-wave-gate.md`, aggregate `question-ledger.md`, and the selected wave folder.
3. Preserve approved reviewer verdicts unless sources, answers, codebase context, or wave scope changed.
4. Relaunch only the selected wave-orchestrator for missing, interrupted, blocked, or needs-revision work.
5. Continue to require user approval before downstream execution or detailing another wave.

Default budgets unless the user sets another value:

- `REVIEW_BUDGET=3` complete planner/reviewer cycles for the selected wave.
- `INTERRUPTION_RETRY_BUDGET=3` controller relaunches per interrupted, stalled, missing-report, or missing-verdict actor.

## Final Response

Report:

- Selected wave id, source wave gate result, output folder, and orchestration folder.
- Wave status, reviewer verdict summary, final fit verdict, AC count, EC count, verification obligation count, and open question count.
- Codebase and neighboring-backend-wave fit summary.
- Whether the user may approve this wave now.
- Validation commands and results.
