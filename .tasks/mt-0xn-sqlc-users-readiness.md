<!-- FILE: .tasks/mt-0xn-sqlc-users-readiness.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Record coverage, QA, and readiness evidence for the mt-0xn sqlc users reference milestone. -->
<!--   SCOPE: Source-plan traceability, command evidence, generated-code gates, accepted risks, and no-push readiness; excludes implementation instructions. -->
<!--   DEPENDS: docs/superpowers/plans/2026-06-05-sqlc-pgx-goose-users-reference.md, docs/verification-plan.xml, commit c2e745b. -->
<!--   LINKS: mt-0xn / M-API / V-M-API. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Coverage Matrix - Maps source-plan requirements to committed implementation and verification evidence. -->
<!--   Verification Evidence - Records focused and broad command outcomes used for closeout. -->
<!--   QA Readiness - Records boundary review, residual risks, and no-push status. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.0 - Added local readiness packet for the sqlc users reference milestone. -->
<!-- END_CHANGE_SUMMARY -->

# mt-0xn sqlc Users Reference Readiness

Date: 2026-06-05
Branch: `mt-0xn-sqlc-users`
Commit: `d4e356b feat(api): adopt sqlc users reference`
Source plan: `docs/superpowers/plans/2026-06-05-sqlc-pgx-goose-users-reference.md`
Result: PASS for local no-push readiness

## Coverage Matrix

| Source-plan area                        | Evidence                                                                                                                                                                                                                     |
| --------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| sqlc tool/config/query/generated output | `apps/api/sqlc.yaml`, `apps/api/tools.go`, `apps/api/project.json`, `apps/api/internal/repository/postgres/queries/users.sql`, and generated package under `apps/api/internal/repository/postgres/generated`.                |
| pgx/v5 + goose alignment                | sqlc config reads existing goose migrations; no schema migration was introduced; `github.com/jackc/pgx/v5` is on `v5.7.5`.                                                                                                   |
| UserRepo adapter boundary               | `apps/api/internal/repository/postgres/user_repo.go` keeps `UserRepo` as the service-facing adapter and hides generated sqlc types behind the postgres repository package.                                                   |
| Unit edge cases                         | `user_repo_unit_test.go` covers generated-query mapping, UUID parsing, cursor decoding, limit overflow, trimming, nullable updates, duplicate email, not-found, query, delete, and malformed timestamp paths.                |
| Real DB integration                     | `user_repo_test.go` covers CRUD, pagination, duplicate email, name-only/email-only/empty updates, missing update, and missing delete against `monorepo_test` with `COVERAGE_GATE=1`.                                         |
| REST/service regression coverage        | `users_test.go`, `users_internal_test.go`, and `user_service_test.go` cover handler success/errors, response encoding edge cases, and service duplicate/error helpers needed to keep API coverage at 100%.                   |
| Generated drift prevention              | `tools/codegen/project.json` runs tool coverage, `bun run codegen`, and `git diff --exit-code` over sqlc/API/web-admin generated outputs. Negative drift proof failed on an intentional query rename, then restored cleanly. |
| Coverage replacement gates              | `tools/coverage/coverage.config.json` and `docs/verification-plan.xml` allowlist only exact sqlc generated files with API codegen, build, and non-skipping repository integration replacement gates.                         |
| GRACE sync                              | `docs/requirements.xml`, `docs/technology.xml`, `docs/development-plan.xml`, `docs/knowledge-graph.xml`, and `docs/verification-plan.xml` reflect sqlc, generated-code drift, coverage policy, and repository proof.         |
| No-goals                                | No service/REST/GraphQL generated model leakage, no manual generated-file markup, no goose schema migration, no branch push, and no MR creation.                                                                             |
| Independent review fixes                | Review findings were fixed before final readiness: query-file GRACE markup no longer copies into generated sqlc output, and `user_service_test.go` now has file-local GRACE metadata.                                        |

## Verification Evidence

| Command                                                                                                                                                                                                                                        | Result                                                                                            |
| ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| `bunx nx run api:codegen`                                                                                                                                                                                                                      | PASS                                                                                              |
| `cd apps/api && go test ./internal/repository/postgres -run TestUserRepo -count=1`                                                                                                                                                             | PASS                                                                                              |
| `bunx nx test api`                                                                                                                                                                                                                             | PASS                                                                                              |
| `bunx nx build api`                                                                                                                                                                                                                            | PASS                                                                                              |
| `bunx nx run codegen:validate`                                                                                                                                                                                                                 | PASS; Nx reported the task as flaky because of the intentional earlier negative drift proof.      |
| `TEST_POSTGRES_PORT=17511 TEST_REDIS_PORT=17512 API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17511/monorepo_test?sslmode=disable cd apps/api && COVERAGE_GATE=1 go test ./internal/repository/postgres -run TestUserRepo -count=1 -v` | PASS with output containing `PASS` and not `SKIP`.                                                |
| `xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml`                                                                         | PASS                                                                                              |
| `grace lint --path .`                                                                                                                                                                                                                          | PASS with 8 pre-existing heuristic export-surface warnings on project-local skill Python scripts. |
| `git diff --cached --check` and `git diff --check` before commit                                                                                                                                                                               | PASS                                                                                              |
| `TEST_POSTGRES_PORT=17521 TEST_REDIS_PORT=17522 API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17521/monorepo_test?sslmode=disable bun run test:coverage`                                                                               | PASS; printed `[Coverage][gate] all thresholds passed`.                                           |
| `TEST_POSTGRES_PORT=17521 TEST_REDIS_PORT=17522 API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17521/monorepo_test?sslmode=disable bun run verify:coverage`                                                                             | PASS; lint, codegen, typecheck, build, coverage, e2e, XML, and GRACE gates completed.             |
| Pre-commit hook on `d4e356b`                                                                                                                                                                                                                   | PASS; lint-staged, `go-lint`, `go-test`, and commitlint passed.                                   |

## QA Readiness

- Backend boundary: PASS. `UserRepo` remains the adapter and service-facing contract; generated sqlc package is internal to postgres repository code.
- Data boundary: PASS. The existing goose schema is reused; there is no schema migration or rollback requirement.
- Tooling boundary: PASS. Codegen validation proves generated drift detection and restores clean generated output.
- Coverage boundary: PASS. Generated sqlc files are excluded only through exact allowlist entries with replacement gates; handwritten API/service/repository/handler paths reach 100% in the full gate.
- GRACE boundary: PASS. XML and file-local contracts align with the committed sqlc implementation; remaining GRACE warnings are pre-existing project-local skill heuristic warnings.
- Review boundary: PASS. Independent review findings were fixed and reverified with focused Go tests, `codegen:validate`, `verify:coverage`, and `grace lint`.
- Tracker boundary: PASS locally using `bd --sandbox` for closeout after the ordinary `bd close` path hit a Dolt auto-sync hang. Local close state is still updated.

## Follow-Up Ledger

No blocking follow-up Beads are required for the selected no-push readiness contour.

Accepted residual notes:

- Nx marks `codegen:validate` as flaky because the rollout intentionally ran a negative drift proof before restoring the generated output.
- Docker ports `17501/17502` were occupied by an existing test stack during verification, so the final broad gates used safe alternate ports `17521/17522` with the same `monorepo_test` DSN policy.
- No branch was pushed and no MR was created.
