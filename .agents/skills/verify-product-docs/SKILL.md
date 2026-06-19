---
name: verify-product-docs
description: Use when product documentation under docs/product or a similar folder is incomplete, contradictory, noisy, or must become a verified product source of truth for decomposition, GRACE planning, PRD handoff, edge-case coverage, or acceptance criteria.
---

# Verify Product Docs

## Overview

Convert messy product notes into a source-backed `docs/product-verified` package. The main session acts as product-level controller: it starts one scope-orchestrator subagent per review scope, and each scope-orchestrator internally runs its own worker/reviewer loop until that scope is approved or blocked. The final package is synthesized only from approved scope outputs, approved derivations, recorded decisions, and explicit open questions.

## Non-Negotiables

- Treat the input folder as read-only. Default input is `./docs/product`; default output is `./docs/product-verified`.
- The main session must not dispatch scope workers or scope reviewers directly. It dispatches one scope-orchestrator per scope.
- Each scope-orchestrator must be able to spawn its own worker and reviewer subagents for its scope. If a required scope cannot complete for any reason, pause the run and follow the interruption recovery protocol.
- Do not use a controller fallback for nested subagents. If scope-orchestrators cannot spawn their scoped worker/reviewer internally, mark the scope blocked and stop before synthesis.
- Do not treat an interrupted, stalled, or missing-report scope as final failure while retry budget remains. Relaunch the same scope-orchestrator with the same `RUN_ID`, prior artifacts, and latest findings.
- Do not perform degraded controller synthesis. The main session must not synthesize final `docs/product-verified` from partial reports, direct analysis, interrupted scope output, or unreviewed worker output.
- Use evidence-constrained synthesis. Direct facts must cite source docs; derived roles, permissions, states, data fields, acceptance criteria, and edge cases must cite the source signal, derivation rationale, and confidence.
- Do not invent unsupported product behavior, API contracts, integrations, or implementation contracts to fill gaps. Missing API or integration contracts are source-gap blockers, not permission to fabricate endpoints, payloads, auth transport, retries, or error formats.
- Keep source evidence. Every material requirement, rule, actor, permission, data field, flow, edge case, or acceptance criterion must trace to a source file, approved scope finding, derivation rationale, explicit decision, or open question.
- Treat additional docs and answered questions as source deltas, not a separate workflow. Record what changed or was answered, pass the delta to every scope-orchestrator, and run the same approval pipeline.
- Preserve contradictions until synthesis. Do not silently pick a side in reviewer reports.
- Record missing information as open questions immediately. Open questions must appear in each relevant scope folder, in the aggregated `.tasks/product-docs-verify/<run-id>/question-ledger.md`, and in `docs/product-verified/open-questions.md` plus `appendix/question-ledger.md` after synthesis.
- Consolidate missing source classes. If an entire required artifact family is absent, create one blocker question for that source gap instead of many speculative questions derived from the missing artifact.
- Write orchestration reports under `.tasks/product-docs-verify/<run-id>/`. Keep final product docs only under `docs/product-verified`.
- Do not leave TODO/TBD placeholders in final verified docs.

## Workflow

1. Resolve paths and create a run id.
   - Use user-provided paths when present.
   - Otherwise use `SOURCE=docs/product`, `OUTPUT=docs/product-verified`, `RUN_ID=$(date -u +%Y%m%dT%H%M%SZ)`, `STAGING=.tasks/product-docs-verify/$RUN_ID/staging/product-verified`.
   - Use `SOURCE_DELTA=.tasks/product-docs-verify/$RUN_ID/source-delta.md` for any new documents, changed documents, removed documents, or answers to previous open questions.
2. Inventory source docs.
   - List files with `find "$SOURCE" -type f | sort`.
   - Treat this as the complete raw source inventory. Files unsupported by the scaffold script must still be recorded as excluded/noisy or source gaps, not silently ignored.
   - If a previous `docs/product-verified/source-inventory.md`, previous `.tasks/product-docs-verify/*/source-delta.md`, or previous `.tasks/product-docs-verify/*/question-ledger.md` exists, create `SOURCE_DELTA` and mark this run as a re-verification/update.
   - Keep the delta lightweight: `Added Sources`, `Changed Sources`, `Removed Sources`, `Answered Questions`, and `Notes`. Use file paths and question ids when available; do not create a parallel planning workflow.
   - If the user provides answers in chat instead of files, record those answers in `SOURCE_DELTA` as `answered-by-user` decisions with question ids when known. They may resolve questions, but they still need traceability in final docs.
   - If the folder is missing or empty, stop and report the blocker unless the user asked to create an empty product template.
3. Load the local contracts.
   - Read `references/output-contract.md`.
   - Read `references/subagent-roles.md`.
4. Scaffold the staging verified folder.
   - Run `python3 .agents/skills/verify-product-docs/scripts/scaffold_product_verified.py --source "$SOURCE" --output "$STAGING"`.
   - Re-run with `--force` only when intentionally regenerating the managed output files.
   - Validate staging shape with `python3 .agents/skills/verify-product-docs/scripts/validate_product_verified.py "$STAGING" --allow-placeholders`.
   - Do not write, touch, delete, or regenerate `docs/product-verified/**` during this step.
5. Dispatch scope-orchestrator subagents.
   - Use `references/subagent-roles.md` as the orchestration contract.
   - Start Phase 1 scope-orchestrators for the seven primary scopes: product scope, roles/permissions, actors/journeys, domain model, feature behavior, edge/risk, and acceptance criteria.
   - Give each scope-orchestrator `SOURCE`, `OUTPUT`, `STAGING`, `RUN_ID`, `SOURCE_DELTA` when present, its scope name/focus, source inventory, and output contract.
   - Give each scope-orchestrator write scope only for `.tasks/product-docs-verify/<run-id>/scopes/<scope>/`.
   - Instruct each scope-orchestrator to spawn its own worker and reviewer subagents internally.
   - Run Phase 1 scope-orchestrators in parallel when the tool supports it.
6. Each scope-orchestrator runs its own worker/reviewer loop.
   - Run the scope worker first.
   - After the worker report, run the scope reviewer.
   - If the reviewer does not approve, send reviewer findings back through the same scope-orchestrator and repeat.
   - If a worker or reviewer attempt is interrupted, stalls, or fails to write its required report/verdict, the scope-orchestrator must launch a replacement worker/reviewer attempt with the same context before declaring the scope blocked.
   - Stop only after the configured review budget or interruption retry budget is exhausted, then record unresolved issues as blocking scope questions.
   - Maintain `.tasks/product-docs-verify/<run-id>/scopes/<scope>/scope-status.md` and `.tasks/product-docs-verify/<run-id>/scopes/<scope>/question-ledger.md`.
   - Use scope-prefixed question ids such as `Q-ROLE-001`, `Q-ACTOR-001`, or `Q-DOMAIN-001` to prevent collisions across parallel ledgers.
7. Main session runs the consistency phase.
   - Aggregate Phase 1 statuses and question ledgers first.
   - If any required Phase 1 scope is not approved, stop before consistency and write a blocked run report.
   - Start the `consistency-reviewer` scope-orchestrator only after Phase 1 scope outputs exist.
   - Give it read access to source docs plus Phase 1 scope outputs, and write access only to `.tasks/product-docs-verify/<run-id>/scopes/consistency-reviewer/`.
   - Apply the same retry policy to consistency as to Phase 1 scopes.
   - If consistency returns `needs-revision`, relaunch `consistency-reviewer` with previous consistency attempt files, reviewer findings, Phase 1 outputs, and the aggregate question ledger until it is approved or the review budget is exhausted.
   - If consistency is interrupted, stalls, or produces no final approved/blocked/needs-revision status, relaunch `consistency-reviewer` with the same `RUN_ID`; do not mark final synthesis blocked until the interruption retry budget is exhausted.
8. Main session aggregates all scope outputs.
   - Write `.tasks/product-docs-verify/<run-id>/scope-status.md` from all scope status files.
   - Write `.tasks/product-docs-verify/<run-id>/question-ledger.md` from all scope question ledgers.
   - Treat any unapproved required scope as a product handoff blocker.
9. Synthesize `docs/product-verified`.
   - Only synthesize after all required Phase 1 scopes and the consistency scope are approved.
   - If the approval gate is not met, stop after updating `.tasks/product-docs-verify/<run-id>/main-orchestration.md`, `scope-status.md`, and `question-ledger.md`; do not write or update final product docs.
   - Read approved scope outputs, reviewer verdicts, aggregated question ledger entries, and explicit decisions.
   - Apply `SOURCE_DELTA` by resolving answered questions as `answered-by-source` or `answered-by-user`, updating affected decisions, and keeping superseded requirements traceable rather than silently deleting them.
   - Use `STAGING` only as a structure reference; do not copy placeholder prose into final docs.
   - Resolve contradictions with this precedence: explicit current source > repeated cross-source agreement > explicit approved decision > open question.
   - Record every approved derivation in `appendix/derivation-log.md` and link it from `appendix/traceability.md`.
   - Mark assumptions in `appendix/decision-log.md`; mark every unresolved blocker in `open-questions.md` and `appendix/question-ledger.md`.
   - If a contradiction changes a core user flow, data ownership rule, money/security behavior, or acceptance criteria, ask the user unless they explicitly requested autonomous resolution.
10. Write the final package using the output contract.

- Cover product intent, scope, actors, permissions, domain model, features, user flows, business rules, edge cases, acceptance criteria, open questions, traceability, and decisions.
- Create one `features/<feature-id>.md` file per material feature when the product has multiple feature areas.

11. Main session validates and fixes only integration issues.

- Run `python3 .agents/skills/verify-product-docs/scripts/validate_product_verified.py "$OUTPUT"`.
- Run `git diff --check -- "$OUTPUT" .tasks/product-docs-verify`.
- Do not overrule scope-orchestrator decisions silently. If validation exposes product ambiguity, send it back through the relevant scope-orchestrator or record a blocking open question.

## Interruption And Recovery

If orchestration is interrupted for any reason, pause the run instead of switching to direct synthesis.

On interruption:

- Update `.tasks/product-docs-verify/<run-id>/main-orchestration.md` with the last completed phase, interrupted scope, and next recovery action.
- Update aggregate `scope-status.md` and `question-ledger.md` from available scope-local artifacts.
- Create or update `.tasks/product-docs-verify/<run-id>/recovery.md`.
- Do not write or update `docs/product-verified/**`.
- Keep any staged skeleton under `.tasks/product-docs-verify/<run-id>/staging/product-verified` as non-final run state.
- Do not label partial output as `usable`, `approved`, `verified`, `source of truth`, or `handoff ready`.

To recover:

1. Resume with the same `RUN_ID`.
2. Read `.tasks/product-docs-verify/<run-id>/scope-status.md`, `question-ledger.md`, and every `scopes/<scope>/scope-status.md`.
3. Preserve approved scopes; do not rerun them unless the source docs or output contract changed.
4. Relaunch scope-orchestrators only for interrupted, missing, blocked, or needs-revision scopes.
5. Pass each relaunched scope-orchestrator its previous attempt files, latest reviewer findings, and scope question ledger.
6. Run the consistency scope only after all required Phase 1 scopes are approved.
7. Synthesize `docs/product-verified/**` only after the approval gate passes.

Default budgets unless the user sets another value:

- `REVIEW_BUDGET=3` complete worker/reviewer cycles per scope.
- `INTERRUPTION_RETRY_BUDGET=3` controller relaunches per scope-orchestrator for interrupted, stalled, missing-report, or missing-verdict runs.
- Interrupted/stalled attempts without a required report or verdict consume interruption retry budget, not review budget.

For the `consistency-reviewer` specifically:

- `needs-revision` means relaunch consistency with the prior reviewer findings; it is not a final blocker while review budget remains.
- `interrupted`, `stalled`, missing `scope-status.md`, or missing final reviewer verdict means relaunch consistency with the same `RUN_ID`; it is not a final blocker while interruption retry budget remains.
- only `approved` opens final synthesis; `blocked` or exhausted budgets stop the run and write recovery instructions.

If source docs changed materially after interruption, record the source delta in `recovery.md` and rerun affected scopes before synthesis.

## Additional Docs And Answered Questions

Do not switch to a separate update workflow when new product docs or answers arrive. Start a normal run with a new `RUN_ID`, then mark the run as a re-verification/update by writing `.tasks/product-docs-verify/<run-id>/source-delta.md`.

Use this minimum structure:

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

- Treat files newly added under `SOURCE` as normal source evidence.
- Treat answers to previous open questions as explicit source evidence only when they are in a source file; otherwise record them as `answered-by-user` decisions in `SOURCE_DELTA`.
- Pass `SOURCE_DELTA` to all scope-orchestrators, including consistency.
- Run the same Phase 1, consistency, synthesis, and validation pipeline.
- Previous approved scope reports are context, not current approval. A re-verification/update run still needs current scope approvals and approved consistency before writing `docs/product-verified/**`.
- If an answered question changes existing verified behavior, preserve the old behavior in `appendix/decision-log.md` or traceability as superseded; do not silently overwrite it.

## Synthesis Rules

- Prefer crisp product requirements over narrative summaries.
- Separate "in scope", "out of scope", "assumption", and "open question".
- Treat missing information as a first-class output, not a weakness to hide. If the docs do not specify behavior, record the exact question, impacted scope, decision needed, and whether it blocks decomposition.
- Separate direct facts, derived requirements, assumptions, and open questions.
- Assumptions are never verified behavior by themselves. An assumption may only appear in `appendix/decision-log.md`, `scope.md`, traceability, or an open question until a source-backed derivation or explicit approved decision promotes it.
- Derived requirements are allowed when they are constrained by the source docs. For each derived role, permission, lifecycle state, data field, acceptance criterion, or edge case, record the source signal, derivation rationale, and confidence.
- Derive roles and permissions from described actors, actions, ownership, approvals, visibility, responsibility boundaries, and denied/allowed flows. If authorization policy is absent or high-risk, record a consolidated auth/source-gap blocker.
- Derive data fields from named entities, forms, workflows, statuses, calculations, validations, reports, imports, exports, and acceptance needs. If a field is needed only because an absent API/database contract is missing, record the missing source artifact instead.
- Generate acceptance criteria from documented or strongly implied behavior. Criteria must cover the behavior already present in the docs and must not add new product behavior.
- Generate edge cases from documented operations plus standard boundary and failure classes around those operations. Do not add unrelated business scenarios.
- Do not promote a speculative assumption into verified product behavior. If a behavior is not supported by source signal or approved derivation, omit it or record it as an open question or explicit decision.
- When missing information belongs to one absent source artifact, create one consolidated blocker. Example: if API documentation is absent, ask `Which API contract/source should be used for implementation?` and mark technical continuation blocked; do not create separate questions for every endpoint, payload, error, retry, and auth detail.
- Source-gap blockers should name the missing artifact class, impacted scopes, why it blocks decomposition or technical continuation, and what artifact or decision is needed to unblock.
- Write acceptance criteria as observable outcomes. Use Given/When/Then when it improves precision.
- Include negative and boundary cases: missing data, invalid input, duplicate actions, permission denial, not found, empty state, concurrent edits, retries, external service failure, offline state, rate limits, data retention, auditability, migrations, and backwards compatibility where relevant.
- Do not import implementation details unless product behavior depends on them.
- Do not weaken conflicts into vague prose. Name the conflicting sources and the chosen decision.

## Final Response

Report:

- Output folder path.
- Orchestration report folder path, scope-orchestrator statuses, and which worker/reviewer loops approved.
- Source delta path and whether the run included new docs, changed docs, removed docs, or answered questions.
- Validation commands and results.
- Any blocking open questions left in `docs/product-verified/open-questions.md`.
