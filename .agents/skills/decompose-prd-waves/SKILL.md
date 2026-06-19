---
name: decompose-prd-waves
description: Use after verify-product-docs or verify-technical-docs when a PRD package needs backend implementation phasing and frontend page handoff files before detailed planning, GRACE planning, Beads, Jira, or developer handoff.
---

<!-- FILE: .agents/skills/decompose-prd-waves/SKILL.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Define the project-local Codex workflow for splitting raw and verified PRD sources into backend-only implementation waves plus per-page frontend handoff files. -->
<!--   SCOPE: Covers controller gates, multi-scope subagent mapping, mandatory review, backend-only shallow wave synthesis, per-page frontend files, references, scripts, and final response shape; excludes frontend waves, per-wave implementation detail, Beads/Jira mutation, and source PRD verification. -->
<!--   DEPENDS: docs/product, docs/product-verified, docs/technical-verified, .agents/skills/decompose-prd-waves/references/output-contract.md, .agents/skills/decompose-prd-waves/references/subagent-roles.md, scripts/scaffold_prd_waves.py, scripts/validate_prd_waves.py. -->
<!--   LINKS: M-PRD-WAVE-DECOMPOSER / V-M-PRD-WAVE-DECOMPOSER / DF-GRACE-CHANGE. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Overview - Establishes the post-verification backend wave and per-page frontend file loop. -->
<!--   Non-Negotiables - Defines source gates, backend-only waves, per-page frontend files, subagent/reviewer requirements, and blockers. -->
<!--   Workflow - Defines run setup, input inventory, scope orchestration, synthesis, validation, and pause points. -->
<!--   Backend Wave And Frontend Page Rules - Defines what backend waves and frontend page files may and must not contain. -->
<!--   Final Response - Defines the closeout evidence to report. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.2 - Replaced the single frontend sequence artifact with page files sourced from raw and verified PRDs. -->
<!-- END_CHANGE_SUMMARY -->

# Decompose PRD Waves

## Overview

Convert raw and verified PRD sources into a shallow, reviewer-approved handoff: backend implementation is grouped into top-level waves, while frontend is captured as one markdown file per page. The main session acts as controller: it dispatches scope-mapper subagents, requires independent review, synthesizes backend-only waves plus per-page frontend files, and stops before detailed planning or task creation.

Default raw product input is `docs/product`. Default verified product input is `docs/product-verified`. Default technical input is `docs/technical-verified` when present. Default output is `docs/prd-waves`. Orchestration reports live under `.tasks/prd-wave-decomposition/<run-id>/`.

## Non-Negotiables

- Start only when verified sources exist, and also read the raw PRD sources for frontend page extraction. Use `docs/product-verified` or an equivalent verified PRD package, plus `docs/product` or an equivalent raw PRD folder. Use `docs/technical-verified` when available; if it exists with blocking questions, record the technical limitations instead of pretending the wave map is dev-ready.
- Do not verify product docs or technical docs inside this skill. Route messy raw product docs to `$verify-product-docs`; route implementation-blocking technical gaps to `$verify-technical-docs`.
- Keep every backend wave shallow. A backend wave may include intent, outcome, included/excluded backend capability groups, dependency order, risk class, recommended next planning surface, and traceability. It must not include implementation tasks, module designs, migration plans, API payloads, acceptance criteria, exit criteria, Beads, Jira issues, code, or frontend page/screen/route/navigation scope.
- Decompose only backend implementation scope into waves. Frontend, UI, UX, navigation, mobile, and client-experience scope must not become waves or wave scope.
- Capture frontend only under `docs/prd-waves/frontend-pages/`: use `frontend-pages/index.md` for page order and shared frontend summary, and one `frontend-pages/page-<nnn>.md` per page. Each page file must describe what is on the page, functional parts, empty states, loading/error states, backend dependencies, deferrals, questions, and traceability to both raw PRD and verified PRD sources. It must not define visual design specs, implementation tasks, API payloads, acceptance criteria, test plans, or Jira/Beads work.
- Account for the whole PRD by routing material backend scope to backend waves, material frontend scope to `frontend-pages/page-<nnn>.md`, and explicit out-of-scope work to deferrals/questions.
- Use subagents. The main session must dispatch one scope-mapper orchestrator per decomposition scope. Each scope-mapper must spawn its own mapper worker and reviewer internally.
- Require review. A final wave map requires approved scope reviews plus a final consistency review. If nested subagents or reviewers are unavailable, mark the run blocked; do not fall back to direct controller synthesis.
- Preserve open questions. Record decomposition-blocking and owner-decision questions in scope ledgers, the aggregate ledger, and final `docs/prd-waves/open-questions.md`.
- Do not hide scope creep. If the shallow backend map needs more than 8 top-level backend waves for an MVP-sized product, stop with a re-scope question unless the user explicitly confirms a broader release scope.
- Do not mutate Beads, Jira, GRACE plans, downstream detailed backend-wave artifacts, or implementation files unless the user asks for a separate follow-up action. This skill may create and promote its own backend-only wave map under `docs/prd-waves/**`.

## Workflow

1. Resolve paths and create a run id.
   - Defaults: `RAW_PRODUCT=docs/product`, `PRODUCT=docs/product-verified`, `TECHNICAL=docs/technical-verified`, `OUTPUT=docs/prd-waves`, `RUN_ID=$(date -u +%Y%m%dT%H%M%SZ)`, `STAGING=.tasks/prd-wave-decomposition/$RUN_ID/staging/prd-waves`.
   - Use user-provided paths when present.
   - If continuing a prior run or answering prior questions, write `.tasks/prd-wave-decomposition/<run-id>/source-delta.md`.
2. Inventory inputs.
   - List files with `{ find "$RAW_PRODUCT" -type f; find "$PRODUCT" -type f; if [ -d "$TECHNICAL" ]; then find "$TECHNICAL" -type f; fi; } | sort`, skipping missing optional technical docs only when verified product docs have enough traceability.
   - Stop if the raw product source or verified product source is missing or empty.
   - Read prior `docs/prd-waves/**` and `.tasks/prd-wave-decomposition/*/question-ledger.md` when this is a continuation.
3. Load local contracts.
   - Read `references/output-contract.md`.
   - Read `references/subagent-roles.md`.
4. Scaffold staging.
   - Run `python3 .agents/skills/decompose-prd-waves/scripts/scaffold_prd_waves.py --raw-product-source "$RAW_PRODUCT" --product-source "$PRODUCT" --technical-source "$TECHNICAL" --output "$STAGING"`.
   - Validate staging with `python3 .agents/skills/decompose-prd-waves/scripts/validate_prd_waves.py "$STAGING" --allow-placeholders`.
   - Do not write final `docs/prd-waves/**` during staging.
5. Dispatch scope-mapper orchestrators.
   - Use `references/subagent-roles.md` as the role contract.
   - Start these scopes: product-capabilities, user-journeys, data-lifecycle, integrations-operations, client-experience, security-compliance, delivery-sequencing.
   - Tell every scope-mapper that wave contributions are backend-only. The `client-experience` mapper produces frontend page-file inputs, not frontend waves.
   - Give each scope-mapper `RAW_PRODUCT`, `PRODUCT`, `TECHNICAL`, `OUTPUT`, `STAGING`, `RUN_ID`, `SOURCE_DELTA` when present, source inventory, scope focus, output contract, and role contract.
   - Give each scope-mapper write scope only for `.tasks/prd-wave-decomposition/<run-id>/scopes/<scope>/`.
   - The main session must not dispatch mapper workers or scope reviewers directly.
6. Each scope-mapper runs its worker/reviewer loop.
   - The mapper worker writes `mapper-attempt-<n>.md`.
   - The scope reviewer writes `review-attempt-<n>.md`.
   - The scope-mapper repeats until reviewer approval, blocker, or budget exhaustion.
   - Maintain `scope-status.md` and `question-ledger.md` from the first gap.
7. Synthesize the candidate shallow wave package.
   - Aggregate primary scope statuses and question ledgers first.
   - Write aggregate `.tasks/prd-wave-decomposition/<run-id>/scope-status.md`.
   - Write aggregate `.tasks/prd-wave-decomposition/<run-id>/question-ledger.md`.
   - Synthesize a controller-owned candidate package under `.tasks/prd-wave-decomposition/<run-id>/staging/prd-waves/**` from approved scope reports, source deltas, decisions, and question ledgers.
   - Put backend implementation grouping only in `wave-map.md`, `waves/index.md`, and `waves/wave-<nn>.md`.
   - Put frontend page order and shared frontend summary only in `frontend-pages/index.md`.
   - Put page-specific frontend content only in `frontend-pages/page-<nnn>.md` files, one file per page, with raw PRD and verified PRD traceability.
   - Do not write or update final `docs/prd-waves/**` before the candidate package passes final consistency review.
8. Run final consistency review.
   - Dispatch `wave-map-consistency` only after all required primary scope outputs and the candidate package exist.
   - Pass the concrete candidate package path to the consistency reviewer.
   - Consistency reviews the candidate package for backend-only wave coverage, per-page frontend coverage, raw and verified frontend traceability, duplicate backend waves, dependency order, traceability, shallow-only compliance, more-than-8-backend-wave risk, and open-question sync.
   - If consistency returns `needs-revision`, revise only the candidate package and rerun consistency before final output.
9. Promote the reviewed shallow wave package.
   - Write aggregate `.tasks/prd-wave-decomposition/<run-id>/scope-status.md`.
   - Write aggregate `.tasks/prd-wave-decomposition/<run-id>/question-ledger.md`.
   - Write or update `docs/prd-waves/**` only by promoting the reviewed candidate package plus approved revisions.
   - Set status to `questions-open`, `blocked`, `waves-ready-awaiting-user-approval`, `waves-approved`, or `superseded` using the output contract.
10. Validate and report.

- Run `python3 .agents/skills/decompose-prd-waves/scripts/validate_prd_waves.py "$OUTPUT"`.
- Run `git diff --check -- "$OUTPUT" .tasks/prd-wave-decomposition .agents/skills/decompose-prd-waves`.
  - If the map is `waves-ready-awaiting-user-approval`, report the backend wave overview and frontend page-file status, then ask for explicit user approval before downstream detailed planning.

## Backend Wave And Frontend Page Rules

A top-level wave is a backend wave. It may contain:

- wave id, name, purpose, and outcome after the wave;
- backend capability groups and source-backed backend PRD areas included;
- explicit exclusions and deferrals;
- predecessor and successor dependencies;
- affected backend, data, integration, security, and operations surfaces at category level;
- risk class and open questions;
- recommended next planning surface, such as `$plan-backend-waves`, `$grace-plan`, `$plan-to-beads`, or a backend domain-specific detailed planner;
- traceability to verified PRD and technical sources.

A backend wave must not contain:

- implementation tasks or task estimates;
- detailed module architecture;
- API request/response payloads;
- database schemas, migrations, jobs, or event contracts;
- acceptance criteria, exit criteria, test cases, or coverage gates;
- Jira tickets, Beads tasks, branch names, or code changes.
- frontend pages, screens, routes, navigation, UI states, UX flows, mobile screens, or client-experience scope.

`frontend-pages/index.md` may contain:

- ordered page ids such as `PAGE-001`, names, and source-backed purpose;
- sequential page order and page-to-page dependencies;
- raw PRD and verified PRD source coverage;
- shared UX states at category level only;
- backend dependency category per page;
- explicit frontend deferrals, frontend questions, and traceability.

Each `frontend-pages/page-<nnn>.md` file may contain:

- page id, page name, and source-backed page purpose;
- what is on the page, functional parts, empty states, and loading/error states;
- backend dependency categories for that page;
- explicit page-level deferrals and questions;
- traceability to raw PRD and verified PRD sources.

Frontend page files must not contain component architecture, visual design specs, copy decks, API payloads, test cases, acceptance criteria, implementation tasks, Jira/Beads tasks, or code.

If the mapper needs those details to make a backend wave or page file coherent, record a decomposition question or route the approved artifact to the appropriate downstream planning skill after user approval.

## Interruption And Recovery

On interruption:

- Update `.tasks/prd-wave-decomposition/<run-id>/main-orchestration.md` with the current phase, interrupted scope, and next recovery action.
- Aggregate available scope statuses and question ledgers.
- Create `.tasks/prd-wave-decomposition/<run-id>/recovery.md`.
- Do not mark partial wave maps `waves-ready-awaiting-user-approval`.

To recover:

1. Resume with the same `RUN_ID`.
2. Preserve approved scope folders unless sources, answers, or output contract changed.
3. Relaunch only missing, interrupted, blocked, or needs-revision scope-mappers.
4. Run final consistency only after required primary scopes have current approved reports.

Default budgets unless the user sets another value:

- `REVIEW_BUDGET=3` complete mapper/reviewer cycles per scope.
- `INTERRUPTION_RETRY_BUDGET=3` controller relaunches per interrupted, stalled, missing-report, or missing-verdict scope-mapper.

## Final Response

Report:

- Output folder path and status.
- Orchestration report folder path and scope statuses.
- Backend wave count, wave ids, wave names, dependency order, frontend page-file status, frontend page count, and open question count.
- Reviewer verdict summary and whether final consistency approved.
- Confirmation that backend waves are shallow/backend-only, frontend is captured only under `frontend-pages/`, each page file traces to raw and verified PRDs, and no implementation detail was produced.
- Recommended next planning surface for wave 1 and whether the user may approve the wave map now.
- Validation commands and results.
