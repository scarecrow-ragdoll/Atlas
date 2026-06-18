<!-- FILE: AGENTS.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Define the repository operating contract for AI agents working under GRACE. -->
<!--   SCOPE: Prime-before-work, GRACE skill routing, docs-first changes, file-local semantic markup, verification, subagent, and quality-gate rules. -->
<!--   DEPENDS: docs/requirements.xml, docs/technology.xml, docs/development-plan.xml, docs/knowledge-graph.xml, docs/verification-plan.xml, docs/operational-packets.xml. -->
<!--   LINKS: M-GRACE-WORKFLOW / V-M-GRACE-WORKFLOW / DF-GRACE-CHANGE. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Prime Before Work - Defines the GRACE artifacts agents must load before changing behavior. -->
<!--   Primary Workflow - Establishes docs/*.xml as the durable shared contract. -->
<!--   File-Local GRACE Contracts - Defines required MODULE_CONTRACT, MODULE_MAP, semantic anchor, and CHANGE_SUMMARY markup for governed file changes. -->
<!--   Web-admin UI Kit Rule - Requires admin pages to import UI only from the approved @shared/ui surface. -->
<!--   Public Web UI Kit Rule - Requires public web pages to use the local apps/web UI kit and stay independent from web-admin UI files. -->
<!--   Testing Policy - Defines behavior-first test expectations and evidence requirements. -->
<!--   Coverage Phases - Separates active development from explicit 100 percent coverage work. -->
<!--   Active Development Verification - Keeps iteration checks focused before final gates. -->
<!--   Subagents - Defines controller and worker ownership boundaries. -->
<!--   Installed Codex Skills - Lists project-local non-GRACE skill entry points. -->
<!--   Quality Gates - Lists required Nx, Bun, Go, XML, and GRACE validation commands. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.6 - Split active development testing from explicit 100 percent coverage phases. -->
<!-- END_CHANGE_SUMMARY -->

# Agent Operating Guide (GRACE)

## Prime Before Work

1. Read `docs/requirements.xml` for product intent, actors, constraints, risks, and open questions.
2. Read `docs/technology.xml` for the approved stack, commands, testing policy, and autonomy policy.
3. Read `docs/development-plan.xml` and `docs/knowledge-graph.xml` for module boundaries, dependencies, write scopes, and implementation order.
4. Read `docs/verification-plan.xml` before changing behavior or tests.
5. Use `docs/operational-packets.xml` for worker packets, graph deltas, verification deltas, failure handoffs, and checkpoint reports.

## Primary Workflow

GRACE is the durable engineering contract for this repository.

- `docs/*.xml` is the shared source of truth for requirements, architecture, graph, execution, and verification.
- `docs/product/` remains product background and glossary material.
- `.tasks/` is for operational reports and review outputs.
- `.protocols/` is for long-running execution state.
- The previous legacy docs workflow is removed. Do not recreate old durable-docs folders or old command aliases.

## Installed GRACE Surface

GRACE skills come from `https://github.com/osovv/grace-marketplace`.

- Codex project skills: `.agents/skills/grace-*/SKILL.md`
- Claude Code project skills: `.claude/skills/grace-*/SKILL.md`
- `.codex/` is only for Codex configuration if needed, not project skills.

Use these entry points:

- `$grace-status` for project health and next safe action.
- `$grace-ask` for grounded questions about the project.
- `$grace-plan` for module, flow, contract, and execution planning.
- `$grace-verification` for test, trace, and evidence design.
- `$grace-execute` for sequential implementation.
- `$grace-multiagent-execute` for controller-managed parallel waves.
- `$grace-refresh` after code changes or suspected drift.
- `$grace-reviewer` for integrity review.
- `$grace-cli` for `grace lint`, `grace status`, `grace module find`, `grace module show`, and `grace file show`.

## Installed Codex Skills

- `$decompose-prd-waves` lives in `.agents/skills/decompose-prd-waves` and splits raw plus verified PRD or technical packages into reviewer-approved backend-only implementation waves plus one source-backed markdown file per frontend page without detailing each wave.
- `$detail-prd-wave` lives in `.agents/skills/detail-prd-wave` and turns exactly one approved backend wave into a code-aware, reviewer-approved ready-for-dev backend wave brief with AC, EC, verification obligations, other-backend-wave fit checks, and read-only `frontend-pages/**` dependency context.
- `$plan-to-beads` lives in `.agents/skills/plan-to-beads` and converts source plans into one Beads milestone with implementation, full-test-coverage, and pre-MR-QA child epics plus source-anchored tasks.

## Docs First

After a meaningful unit of work:

1. Update the relevant `docs/*.xml` contracts and navigation.
2. Update source, tests, and semantic markup consistently.
3. Run the focused verification commands from `docs/verification-plan.xml`.
4. Run broader quality gates when shared surfaces changed.
5. Commit only after the GRACE artifacts and verification evidence are in sync.

## File-Local GRACE Contracts

New or meaningfully edited governed files must carry file-local GRACE markup. This applies to source, tests, tooling, config, scripts, Docker, CI, and durable docs outside `docs/*.xml` when the file owns behavior, workflow, verification, or a public contract. Shared GRACE XML artifacts are governed by their XML structure and references directly.

Do not mass-backfill every historical file only because it lacks markup. When an existing unmarked file becomes part of a meaningful change, add the contract then.

Use the file's native comment syntax:

```ts
// FILE: path/to/file.ext
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: One sentence describing what this file owns.
//   SCOPE: Included responsibilities and explicit exclusions.
//   DEPENDS: Local modules, generated artifacts, external services, or none.
//   LINKS: GRACE module and verification refs, for example M-API / V-M-API.
//   ROLE: RUNTIME | TEST | BARREL | CONFIG | TYPES | SCRIPT | DOC
//   MAP_MODE: EXPORTS | LOCALS | SUMMARY | NONE
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   exportedSymbol - one-line contract summary
// END_MODULE_MAP
```

For Markdown and XML, use HTML comments. For YAML, Dockerfiles, shell, and Python, use `#`. For Go and TypeScript, use `//`.

Function, component, and critical block anchors are required when they improve navigation or verification:

```ts
// START_CONTRACT: functionName
//   PURPOSE: Observable behavior.
//   INPUTS: { paramName: Type - description }
//   OUTPUTS: { ReturnType - description }
//   SIDE_EFFECTS: External state changes or none.
//   LINKS: Related module or verification refs.
// END_CONTRACT: functionName

// START_BLOCK_VALIDATE_INPUT
// ... code ...
// END_BLOCK_VALIDATE_INPUT
```

Use `START_CHANGE_SUMMARY` for bug fixes, contract changes, migrations, or risky behavior changes:

```ts
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - What changed and why.
// END_CHANGE_SUMMARY
```

Rules for modifications:

1. Read the existing `MODULE_CONTRACT` before editing a governed file.
2. Create or update `MODULE_CONTRACT` before changing behavior in a new or previously unmarked governed file.
3. Keep `MODULE_MAP` aligned with exports, public entry points, test helpers, or config surfaces according to `MAP_MODE`.
4. Add or update `START_CONTRACT` and `START_BLOCK_*` anchors for non-trivial public functions, components, handlers, migrations, critical branches, and log-marker paths.
5. After adding, moving, or removing public modules or dependencies, update `docs/knowledge-graph.xml`.
6. After changing scenarios, tests, commands, coverage policy, or log markers, update `docs/verification-plan.xml`.
7. Never remove semantic markup anchors unless the structure is intentionally replaced with better anchors.

## Web-admin UI Kit Rule

All admin pages under `apps/web-admin/src/pages/**` must build UI from the approved `@shared/ui` surface.

- Store shadcn-generated primitives under `apps/web-admin/src/shared/ui/primitives/**`.
- Store admin page compositions under `apps/web-admin/src/shared/ui/layout/**`.
- Import UI in pages from `@shared/ui` only.
- Do not import Radix primitives, shadcn implementation subpaths, UI-kit aliases such as `@/shared/ui` or `@shared/ui/*`, relative `shared/ui` paths, class composition helpers, or icon libraries directly from page files.
- Prove the rule with `bunx nx lint web-admin` after page or UI-kit changes.

## Public Web UI Kit Rule

Public pages and client components under `apps/web/**` must use the local public web UI-kit surface when shared primitives are needed.

- Store public web shadcn-compatible primitives under `apps/web/src/shared/ui/primitives/**`.
- Store public web UI helpers under `apps/web/src/shared/ui/lib/**`.
- Import public web primitives from `@shared/ui`.
- Do not import from `apps/web-admin/src/shared/ui/**`, `apps/web-admin` aliases, or web-admin generated shadcn implementation files into `apps/web`.
- Keep public web `components.json` separate from `apps/web-admin/components.json`.
- Prove the rule with `bunx nx test web` and `bunx nx run web:typecheck` after public web page or UI-kit changes.

## Testing Policy

Tests are part of the implementation contract, not a follow-up task. Every behavior change must update the nearest useful tests first or in the same change set, then prove the affected surface with focused checks before broader gates.

- Prefer behavior and contract tests over implementation-detail tests. Test what the user, API client, or module consumer can observe: success paths, validation errors, not-found cases, duplicate/conflict cases, auth/config failures, retries, and integration boundaries.
- Keep unit tests fast and local. Mock outside dependencies at clear module boundaries, but do not mock the code path being verified just to make coverage pass.
- Use integration tests for real persistence, cache, GraphQL/API, or cross-package behavior. When Docker-backed services are required, make unavailable-service behavior explicit and fail under the coverage gate when required services are expected.
- Use e2e tests for real vertical slices that cross UI, API, GraphQL, and database boundaries. Use deterministic setup data, unique identifiers, accessible UI locators, and explicit assertions on the final user-visible state.
- Frontend tests must cover page states, forms, loading/empty/success/error paths, API client behavior, and config parsing. Prefer roles, labels, placeholders, and visible text over brittle DOM structure.
- Backend tests must cover services, repositories, handlers/resolvers, validation mapping, error translation, and external resource lifecycle. Repository tests should prove both happy paths and database-level edge cases.
- Tooling code is production code. Nx executors, codegen config, coverage scripts, and project-local utilities need tests when changed.
- Coverage applies to handwritten behavior. Generated files and bootstrap entrypoints may be excluded only through `tools/coverage/coverage.config.json`, and every exclusion must have a replacement gate such as codegen, build, typecheck, or e2e startup coverage.
- Do not add broad coverage exclusions, snapshot-only assertions, or empty tests. If coverage is hard to reach, improve the seam with a small testable boundary instead of weakening the gate.
- Record meaningful verification evidence in `.tasks/` when closing a task, epic, handoff, or risky change.

## Coverage Phases

Coverage work has two explicit modes:

- During active development, agents write the minimal useful tests needed to prove the new behavior, regression fix, or integration contract. Do not chase 100 percent coverage, do not add tests only to satisfy uncovered implementation details, and do not treat a coverage percentage as a completion blocker.
- During an explicit coverage phase, the 100 percent coverage contract applies to non-allowlisted handwritten Go and TypeScript behavior. Enter this mode only when the user explicitly asks for 100 percent coverage, when closing a coverage epic, when changing coverage policy or allowlists, or when preparing release/template handoff where the full gate is part of the delivery contract.
- If the user asks for normal feature development, bug fixing, or refactoring without mentioning 100 percent coverage, coverage gaps may be recorded as follow-up risk but should not expand the task beyond the smallest behavior-focused test set.

## Active Development Verification

During active development, keep verification tight and local. Do not run heavy repo-wide gates while iterating unless the user explicitly asks for them or the work is at final closeout.

- Run only focused checks for the changed files, packages, or Nx projects: exact unit test files, nearest package tests, affected typecheck/lint targets, codegen for touched schemas, and small contract checks tied to the current diff.
- Avoid heavy gates during the implementation loop, including full `bun run test`, full `bun run build`, `bun run test:coverage`, `bun run test:e2e`, `bun run verify:coverage`, and broad run-many commands that exercise unrelated modules.
- Save heavy e2e, full build, full test, and release/handoff gates for the end, after code, tests, generated artifacts, and GRACE docs are synchronized. Save 100 percent coverage gates for explicit coverage phases.
- If a heavier command is genuinely needed before final closeout because the changed surface is shared or risky, state the reason first and keep the run to the narrowest command that can prove the risk.

## Subagents

- Main session is the controller.
- Workers operate from compact execution packets and explicit write scopes.
- Workers write detailed reports to `.tasks/TASK-*/`.
- Controller owns shared artifact updates in `docs/development-plan.xml`, `docs/knowledge-graph.xml`, and `docs/verification-plan.xml`.
- Independent module work may run in parallel only when write scopes and dependency order allow it.

## Quality Gates

Use the real Nx/Bun/Go commands from this repository:

- `bun run lint`
- `bun run test`
- `bun run codegen`
- `bunx nx run web:typecheck`
- `bun run build`
- `bun run test:coverage` during explicit coverage phases or after coverage-policy/allowlist changes
- `bun run test:e2e` when UI/API flow coverage is required
- `bun run verify:coverage` before release/template handoff, coverage epic closeout, explicit 100 percent coverage requests, or coverage-policy changes

For GRACE integrity:

- `xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml`
- `grace lint --path .`
- `grace lint --profile autonomous --path .` before long autonomous execution

<!-- BEGIN BEADS INTEGRATION v:1 profile:minimal hash:7510c1e2 -->
## Beads Issue Tracker

This project uses **bd (beads)** for issue tracking. Run `bd prime` to see full workflow context and commands.

### Quick Reference

```bash
bd ready              # Find available work
bd show <id>          # View issue details
bd update <id> --claim  # Claim work
bd close <id>         # Complete work
```

### Rules

- Use `bd` for ALL task tracking — do NOT use TodoWrite, TaskCreate, or markdown TODO lists
- Run `bd prime` for detailed command reference and session close protocol
- Use `bd remember` for persistent knowledge — do NOT use MEMORY.md files

**Architecture in one line:** issues live in a local Dolt DB; sync uses `refs/dolt/data` on your git remote; `.beads/issues.jsonl` is a passive export. See https://github.com/gastownhall/beads/blob/main/docs/SYNC_CONCEPTS.md for details and anti-patterns.

## Session Completion

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** - Clear stashes, prune remote branches
6. **Verify** - All changes committed AND pushed
7. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds
<!-- END BEADS INTEGRATION -->
