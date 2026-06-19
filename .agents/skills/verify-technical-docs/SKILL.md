---
name: verify-technical-docs
description: Use after verify-product-docs when docs/product-verified or another verified PRD package must be turned into a technical readiness package, after answered product questions have been incorporated into a verified product package and need another technical gap loop, or when Codex must decide questions-open versus approved-to-dev before GRACE planning, Beads decomposition, implementation planning, or developer handoff.
---

<!-- FILE: .agents/skills/verify-technical-docs/SKILL.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Define the project-local Codex workflow for technical readiness verification after product docs are verified. -->
<!--   SCOPE: Covers controller workflow, question-loop rules, approval gates, source-delta handling, references, scripts, and final response shape; excludes product-doc synthesis and implementation planning. -->
<!--   DEPENDS: .agents/skills/verify-technical-docs/references/output-contract.md, .agents/skills/verify-technical-docs/references/subagent-roles.md, scripts/scaffold_technical_verified.py, scripts/validate_technical_verified.py, docs/product-verified. -->
<!--   LINKS: M-TECHNICAL-DOCS-VERIFIER / V-M-TECHNICAL-DOCS-VERIFIER / DF-GRACE-CHANGE. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Overview - Establishes the technical verification loop that follows product verification. -->
<!--   Non-Negotiables - Defines read-only source handling, no-invention policy, question ledgers, and approved-to-dev gate. -->
<!--   Workflow - Defines run setup, source delta, staging, scope orchestration, synthesis, and validation. -->
<!--   Question Loop - Defines how answered questions are folded into later runs and when the loop may close. -->
<!--   Final Response - Defines the closeout evidence to report. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.0 - Added the technical docs verification workflow skill. -->
<!-- END_CHANGE_SUMMARY -->

# Verify Technical Docs

## Overview

Convert a verified product package into a technical readiness package that names every implementation-blocking technical gap before development starts. The main session acts as technical controller: it starts one scope-orchestrator subagent per technical scope, each scope runs worker/reviewer loops internally, and the final status is `approved-to-dev` only when answered questions do not create new blocking questions.

Default input is `docs/product-verified`. Default output is `docs/technical-verified`. Orchestration reports live under `.tasks/technical-docs-verify/<run-id>/`.

## Non-Negotiables

- Treat `docs/product-verified` and prior product reports as read-only inputs. Record answers or changed sources as source deltas.
- Do not invent implementation contracts. Missing API, data, auth, integration, infrastructure, migration, observability, or testing artifacts are technical source gaps, not permission to fabricate endpoints, schemas, jobs, retries, permissions, environments, or SLOs.
- Write every technical hole as a question immediately in the scope ledger, aggregate ledger, and final `docs/technical-verified/open-questions.md`.
- Consolidate absent artifact classes. If the API contract is missing, ask one API-contract blocker instead of separate endpoint, payload, auth, retry, pagination, and error-shape questions.
- The final package may be written while questions remain, but its status must be `questions-open` or `blocked`. Never call it dev-ready unless the approval gate passes.
- `approved-to-dev` requires all required scopes and consistency approved, zero open `dev-blocking` or `needs-owner-decision` questions, no unreviewed source delta, and no new blocking question spawned by answers in the current loop.
- If answered questions create new blocking questions, keep the original answers traceable, add follow-up question ids linked to the parent question, and run another loop.
- Do not start GRACE planning, Beads decomposition, or implementation planning from this skill unless the user explicitly asks after an `approved-to-dev` result.

## Workflow

1. Resolve paths and create a run id.
   - Use user-provided paths when present.
   - Otherwise use `SOURCE=docs/product-verified`, `OUTPUT=docs/technical-verified`, `RUN_ID=$(date -u +%Y%m%dT%H%M%SZ)`, and `STAGING=.tasks/technical-docs-verify/$RUN_ID/staging/technical-verified`.
   - Use `SOURCE_DELTA=.tasks/technical-docs-verify/$RUN_ID/source-delta.md` when prior technical output, prior question ledgers, changed product docs, or answered questions exist.
   - Write `SOURCE_DELTA` with the minimum structure below; keep it lightweight, but do not omit answered question ids, answer sources, affected scopes, or expected second-order effects.
2. Inventory inputs.
   - List files with `find "$SOURCE" -type f | sort`.
   - If previous `docs/technical-verified/appendix/question-ledger.md` or `.tasks/technical-docs-verify/*/question-ledger.md` exists, treat the run as a loop continuation.
   - If the user answered questions in chat, record them in `SOURCE_DELTA` as `answered-by-user` with original question ids when known.
   - Stop if `SOURCE` is missing or empty unless the user explicitly asks for an empty technical template.
3. Load the local contracts.
   - Read `references/output-contract.md`.
   - Read `references/subagent-roles.md`.
4. Scaffold the staging technical package.
   - Run `python3 .agents/skills/verify-technical-docs/scripts/scaffold_technical_verified.py --source "$SOURCE" --output "$STAGING"`.
   - Re-run with `--force` only when intentionally regenerating managed staging files.
   - Validate staging shape with `python3 .agents/skills/verify-technical-docs/scripts/validate_technical_verified.py "$STAGING" --allow-placeholders`.
   - Do not write final `docs/technical-verified/**` during staging.
5. Dispatch Phase 1 scope-orchestrators.
   - Use `references/subagent-roles.md` as the orchestration contract.
   - Start these scopes: architecture-boundaries, data-contracts, api-contracts, auth-security-compliance, integrations-events, client-state-ux, operations-observability, testing-delivery.
   - Give each scope-orchestrator `SOURCE`, `OUTPUT`, `STAGING`, `RUN_ID`, `SOURCE_DELTA` when present, source inventory, scope focus, output contract, and role contract.
   - Give each scope-orchestrator write scope only for `.tasks/technical-docs-verify/<run-id>/scopes/<scope>/`.
   - The main session must not dispatch scope workers or scope reviewers directly.
6. Each scope-orchestrator runs its worker/reviewer loop.
   - The scope worker writes `worker-attempt-<n>.md`.
   - The scope reviewer writes `review-attempt-<n>.md`.
   - The scope-orchestrator repeats until reviewer approval or budget exhaustion.
   - Maintain `scope-status.md` and `question-ledger.md` from the first missing-info finding.
7. Run consistency and loop-closure review.
   - Aggregate Phase 1 statuses and question ledgers.
   - Start `consistency-loop-reviewer` only after Phase 1 scope outputs exist.
   - This scope checks cross-scope contradictions, duplicate questions, unanswered parent questions, and whether answers created new blockers.
8. Aggregate and synthesize.
   - Write aggregate `.tasks/technical-docs-verify/<run-id>/scope-status.md`.
   - Write aggregate `.tasks/technical-docs-verify/<run-id>/question-ledger.md`.
   - Write or update `docs/technical-verified/**` from approved scope reports, reviewer verdicts, aggregate question ledger, source delta, decisions, and traceability.
   - Set status to `questions-open`, `blocked`, or `approved-to-dev` using the approval gate. Do not hide open questions to achieve approval.
9. Validate and fix integration issues only.
   - Run `python3 .agents/skills/verify-technical-docs/scripts/validate_technical_verified.py "$OUTPUT"`.
   - Run `git diff --check -- "$OUTPUT" .tasks/technical-docs-verify .agents/skills/verify-technical-docs`.
   - If validation exposes technical ambiguity, send it back through the relevant scope-orchestrator or record a blocking question.

## Question Loop

Use one normal workflow for initial runs and reruns. New product docs, changed `docs/product-verified`, technical answers, and owner decisions are all source deltas.

Minimum `source-delta.md` structure:

```text
# Source Delta
## Previous Baseline
## Product Verified Changes
## Added Sources
## Changed Sources
## Removed Sources
## Answered Technical Questions
| Question ID | Answer Source | Answer Summary | Affected Scopes | Expected Effect |
| --- | --- | --- | --- | --- |
## Answered Product Questions
| Product Question ID | Answer Source | Technical Impact | Affected Scopes |
| --- | --- | --- | --- |
## Notes
```

- Preserve original question ids when answers arrive.
- Mark answered questions as `answered-by-source`, `answered-by-user`, or `resolved-by-decision`.
- Add follow-up ids with a `Parent` field when an answer creates more questions.
- Treat follow-up `dev-blocking` or `needs-owner-decision` questions as proof the loop is not closed.
- Permit `approved-to-dev` only after consistency confirms that all source deltas and answer deltas were reviewed and produced no new blocking questions.
- Keep non-blocking `deferred` items only when they have an explicit owner, a deferral rationale, and no implementation-blocking impact.

## Output Rules

The final package must follow `references/output-contract.md`. It should describe:

- technical status and handoff readiness;
- architecture boundaries and unknowns;
- data model, storage, migration, and retention gaps;
- API, integration, async/event, and external contract gaps;
- auth, security, compliance, and auditability gaps;
- client state, UX technical states, offline/loading/error states, and accessibility blockers;
- operations, observability, deployment, environment, and SLO questions;
- testing, e2e, fixture, seed data, and coverage strategy blockers;
- implementation slices only when supported by the verified product package or explicit decisions.

## Interruption And Recovery

If orchestration is interrupted, do not synthesize an approved package.

- Update `.tasks/technical-docs-verify/<run-id>/main-orchestration.md` with the interrupted scope and next action.
- Aggregate any available scope statuses and question ledgers.
- Create `.tasks/technical-docs-verify/<run-id>/recovery.md`.
- Preserve approved scope folders.
- Resume with the same `RUN_ID`; relaunch only missing, interrupted, blocked, or needs-revision scopes.
- Run consistency only after required Phase 1 scopes have current approved reports.

Default budgets unless the user sets another value:

- `REVIEW_BUDGET=3` complete worker/reviewer cycles per scope.
- `INTERRUPTION_RETRY_BUDGET=3` controller relaunches per interrupted, stalled, missing-report, or missing-verdict scope-orchestrator.

## Final Response

Report:

- Output folder path and status: `questions-open`, `blocked`, or `approved-to-dev`.
- Orchestration report folder path and scope statuses.
- Source delta path and whether this was an initial run or question-loop rerun.
- Count of open `dev-blocking`, `needs-owner-decision`, `deferred`, and resolved questions.
- Whether answered questions created any follow-up blockers.
- Validation commands and results.
