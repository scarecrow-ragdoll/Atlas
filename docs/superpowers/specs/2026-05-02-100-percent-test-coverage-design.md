# 100 Percent Test Coverage Design

## Goal

Make this repository a template-grade baseline with strict, repeatable test infrastructure. The repository should enforce 100 percent coverage for handwritten production code and include mandatory end-to-end scenarios that prove the generated project works as a real full-stack system.

The desired outcome is not a decorative metric. The coverage system must make downstream repositories inherit a reliable testing contract: every meaningful handwritten behavior is tested, generated artifacts are validated by generation/build gates, and full-stack user flows are exercised through the browser and API.

## Current Context

The repository already has the core testing stack:

- Nx and Bun scripts for root `lint`, `test`, `build`, and `codegen`.
- Go unit tests in API, bot, and shared Go libraries.
- Vitest and Testing Library in `apps/web`.
- A configured Playwright target for `web:e2e`.
- GRACE artifacts in `docs/*.xml` that define verification expectations.

The current gaps are:

- Coverage thresholds are not enforced consistently at the root.
- Go projects write local `coverage.out` files, but there is no 100 percent enforcement layer.
- Vitest config has a 70 percent threshold, and the project test target does not explicitly run coverage.
- Playwright is configured, but there are no actual e2e test files.
- Generated files and entrypoints are not governed by a clear coverage exclusion policy.
- API GraphQL update/delete paths and repository update/delete methods are currently not implemented, yet the schema exposes them.
- First-party workspace tooling under `tools/nx-go` and `tools/codegen` is handwritten template behavior, but it does not yet have explicit test or coverage gates.

## Coverage Contract

The repository will define a single coverage contract:

1. Handwritten production code must have 100 percent coverage.
   - TypeScript/Vitest: statements, branches, functions, and lines must all be 100 percent.
   - Go: statement coverage must be 100 percent using the standard Go coverage tool.
   - This includes runtime app/library code and first-party workspace tooling that changes template behavior.

2. Generated code is excluded from line coverage.
   - Generated code must still be protected by codegen, typecheck, schema validation, and build gates.
   - Generated files must be listed explicitly in the coverage allowlist.

3. Entry points and bootstrap glue may be excluded from line coverage only when covered by startup, build, or e2e gates.
   - Examples: `apps/api/cmd/server/main.go`, `apps/bot/cmd/bot/main.go`.
   - These exclusions must be justified in the allowlist and in `docs/verification-plan.xml`.

4. No file can be silently excluded.
   - Every exclusion must have an owner-facing reason.
   - The coverage enforcement tool must fail if a new uncovered handwritten file is not tested or explicitly allowlisted.

5. End-to-end coverage is scenario based, not percentage based.
   - E2E scenarios must prove the template runs through real browser, API, GraphQL, and persistence paths.

## Allowed Coverage Exclusions

Initial exclusions should be limited to:

- `apps/api/internal/graph/generated.go`
- `apps/api/internal/graph/model/models_gen.go`
- `apps/web/src/shared/api/generated/**`
- `apps/api/cmd/server/main.go`, only if covered by API startup or full-stack e2e gates.
- `apps/bot/cmd/bot/main.go`, only if covered by build/startup smoke gates.
- Framework or tool artifacts such as `.next`, `node_modules`, `coverage`, `dist`, and TypeScript build info files.
- Config-only files may be excluded from line coverage only when a dedicated lint, typecheck, build, codegen, or command-contract gate validates them. These exclusions must be explicit file entries, not broad globs.

This list is intentionally small. If implementation discovers another candidate, the implementation plan must require an explicit decision rather than adding a broad glob.

## Test Infrastructure

Add a dedicated coverage enforcement surface instead of scattering policy across project files:

- `tools/coverage/`
  - Runs coverage commands.
  - Reads and applies the allowlist.
  - Parses Go and Vitest coverage outputs.
  - Fails below 100 percent for handwritten code.
  - Writes normalized artifacts under `dist/coverage`.

- Root scripts:
  - `test:coverage`: run unit, contract, and integration coverage gates.
  - `test:e2e`: run mandatory e2e scenario matrix.
  - `verify:coverage`: run the full template gate.

- Nx targets remain module focused:
  - `api:test`
  - `bot:test`
  - `go-config:test`
  - `go-logger:test`
  - `web:test`
  - `web:e2e`
  - `nx-go:test`
  - `codegen:validate`

Coverage artifacts should be written to `dist/coverage/...` rather than module roots. This avoids stale `coverage.out` files beside source code and makes CI artifact collection straightforward.

## Go Coverage Strategy

Use the standard Go coverage tooling for Go code.

Required behavior:

- Each Go project emits a coverage profile under `dist/coverage/go/<project>/coverage.out`.
- The enforcement tool reads the profile and requires 100 percent statement coverage after applying the allowlist.
- Table-driven tests must cover success and failure paths for branch-heavy code.
- Integration-style tests are allowed where repository or HTTP behavior is more valuable than narrow mocks.

Go does not provide branch and function percentages through the standard coverage tool in the same way Vitest does. The contract for Go is therefore 100 percent statement coverage plus scenario requirements in `docs/verification-plan.xml`.

## TypeScript Coverage Strategy

Use Vitest coverage for web unit and component tests.

Required behavior:

- `apps/web` test coverage thresholds must be set to 100 for statements, branches, functions, and lines.
- Web generated GraphQL types are excluded from coverage and covered by codegen/typecheck/build.
- Component tests should cover loading, success, empty, and error states where those states exist.
- GraphQL client behavior should be tested through stable request/mocking boundaries rather than ad hoc global monkey-patching.

## Workspace Tooling Coverage Strategy

First-party tooling is part of the template contract because downstream repositories inherit it.

Required behavior:

- `tools/nx-go/src/**` must have unit tests or command-contract tests for executor argument construction, working-directory handling, option handling, and failure propagation.
- `tools/codegen/codegen.ts` must be validated by the actual codegen command and by a focused assertion that the schema, document glob, and output paths match the repository layout.
- Tooling tests must run from root verification and produce coverage artifacts under `dist/coverage/tools/...` when the tested surface is executable TypeScript.
- Config-only tooling files can be excluded from line coverage only if their command-contract validation is listed in `docs/verification-plan.xml`.

## E2E Matrix

E2E tests must prove the template works as a full-stack starting point.

Minimum required scenarios:

1. Infrastructure startup
   - Start `docker compose -f docker/docker-compose.dev.yml up -d`.
   - Confirm PostgreSQL and Redis are reachable.

2. API startup
   - Start the API.
   - Confirm migrations run.
   - Confirm `/healthz` and `/readyz` return successful responses.

3. GraphQL contract
   - Exercise `users`.
   - Exercise `user(id)`.
   - Exercise `createUser`.
   - Exercise `updateUser`.
   - Exercise `deleteUser`.

4. Validation and error behavior
   - Duplicate email or invalid input returns a stable typed error.
   - Failure paths do not persist partial state.

5. Web happy path
   - Open `/users`.
   - See the user list.
   - Create a user.
   - See the created user in the list.
   - Open `/users/[id]` and see the detail page.

6. Web failure path
   - Simulate or trigger a GraphQL failure.
   - Confirm the UI shows a controlled error state.

7. Template smoke
   - A clean checkout can run the full verification gate after install and local infrastructure startup.

Test isolation should prefer cleanup-before-each using a test helper or direct database cleanup. This is simpler for downstream repositories than adding multiple Docker profiles at the template stage.

## API Contract Decision

The GraphQL schema exposes create, read, update, and delete operations. The reference vertical slice should therefore implement and test CRUD completely instead of leaving update/delete as stubs.

Implementation should:

- Implement repository `Update` and `Delete`.
- Implement GraphQL `updateUser` and `deleteUser` resolvers.
- Add unit and integration tests for success, not-found, validation, and duplicate paths.
- Cover the corresponding GraphQL operations in e2e tests.

If a future template owner wants a smaller reference slice, the correct alternative is to remove update/delete from the schema and docs. Leaving public schema operations as panics or `not implemented` is not compatible with the 100 percent coverage goal.

## Verification Commands

The final verification surface should include:

- `bun run lint`
- `bun run test`
- `bun run build`
- `bun run codegen`
- `bun run test:coverage`
- `bun run test:e2e`
- `bun run verify:coverage`
- `bunx nx run web:typecheck`
- `xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml`
- `grace lint --path .`

For deterministic Nx validation in this workspace, final runs should use:

```bash
bunx nx reset
NX_SKIP_NX_CACHE=1 NX_DAEMON=false bunx nx run-many --target=lint --all --parallel=1
NX_SKIP_NX_CACHE=1 NX_DAEMON=false bunx nx run-many --target=test --all
NX_SKIP_NX_CACHE=1 NX_DAEMON=false bunx nx run-many --target=build --all
```

`verify:coverage` must be the single user-facing root command for the full gate. It should delegate to smaller scripts where that keeps output readable, but the command order and required checks must be fixed by the implementation plan rather than left to individual contributors.

## GRACE Artifact Updates

The implementation must keep GRACE in sync.

Required updates:

- `docs/requirements.xml`
  - Add the 100 percent handwritten coverage requirement.
  - Add e2e scenario coverage as a template acceptance criterion.

- `docs/technology.xml`
  - Add coverage enforcement tooling.
  - Document the Go and TypeScript coverage semantics.

- `docs/development-plan.xml`
  - Add a test infrastructure module or extend `M-WORKSPACE`.
  - Add first-party tooling coverage to the workspace orchestration contract.
  - Add implementation order for coverage enforcement and e2e closure.

- `docs/knowledge-graph.xml`
  - Add graph entries or annotations for the coverage enforcement surface.

- `docs/verification-plan.xml`
  - Add the coverage contract.
  - Add the e2e matrix.
  - Add explicit allowlist policy.
  - Add final verification commands and evidence expectations.

- `docs/operational-packets.xml`
  - Add worker packet guidance for coverage closure work.

## Implementation Slices

The implementation plan should split work into these slices:

1. Baseline inventory
   - Capture current coverage numbers.
   - Identify uncovered handwritten files.
   - Confirm generated and entrypoint exclusion candidates.

2. Coverage enforcement foundation
   - Add `tools/coverage`.
   - Add root scripts.
   - Route artifacts to `dist/coverage`.
   - Fail below 100 percent.

3. Workspace tooling coverage closure
   - Add coverage or command-contract gates for `tools/nx-go` and `tools/codegen`.
   - Ensure tooling checks run from the root coverage gate.

4. Go coverage closure
   - Complete API CRUD behavior.
   - Add repository, resolver, service, middleware, config, logger, and bot tests as needed.
   - Keep generated gqlgen files excluded but validated.

5. Web coverage closure
   - Add Vitest tests for pages, providers, config, GraphQL client behavior, and error states.
   - Raise thresholds to 100.

6. E2E closure
   - Add Playwright tests for the full `/users` flow and API/GraphQL behavior.
   - Add deterministic cleanup.
   - Ensure e2e scenarios are part of the final root gate.

7. GRACE and docs sync
   - Update all required `docs/*.xml` contracts.
   - Keep README command examples aligned.

8. Final verification
   - Run deterministic Nx gates.
   - Run coverage and e2e gates.
   - Run XML and GRACE lint.

## Error Handling and Failure Reporting

Coverage and e2e failures should produce actionable evidence:

- Which module failed.
- Which coverage metric failed.
- Which files or functions are uncovered.
- Whether the uncovered file is allowlisted.
- Which e2e scenario failed.
- First divergent command or observable behavior.

Failure messages should distinguish:

- Code behavior failure.
- Coverage policy failure.
- Missing test fixture or local infrastructure failure.
- Generated artifact drift.
- GRACE artifact drift.

This distinction matters because the repository is a template. A downstream maintainer must be able to tell whether they need to fix product code, add tests, regenerate artifacts, or start local infrastructure.

## Non-Goals

- Do not test generated code line by line.
- Do not add broad coverage exclusion globs.
- Do not rely on e2e tests to inflate unit coverage metrics.
- Do not introduce a second task orchestrator beside Nx and Bun.
- Do not make a CI-only system that cannot run locally.

## Acceptance Criteria

The design is complete when implementation can prove:

- Handwritten production code coverage is enforced at 100 percent.
- First-party workspace tooling is covered or explicitly command-contract validated.
- Generated and entrypoint exclusions are explicit and justified.
- E2E scenarios cover the required full-stack matrix.
- API CRUD schema operations are implemented or removed from the public contract.
- `verify:coverage` is the single root command for the full coverage gate.
- Coverage artifacts are stored under `dist/coverage`.
- GRACE artifacts describe the same coverage contract that the tooling enforces.
