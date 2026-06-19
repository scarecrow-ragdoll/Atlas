<!-- FILE: .tasks/web-admin-login-route-guard/readiness.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Record the final no-push readiness packet for the web-admin login route guard milestone. -->
<!--   SCOPE: Summarizes implementation, coverage, QA audits, verification evidence, blockers, accepted risks, and delivery state; excludes source-level implementation detail already represented in code and tests. -->
<!--   DEPENDS: .tasks/web-admin-login-route-guard/verification.md, docs/superpowers/plans/2026-06-07-web-admin-login-route-guard.md, apps/web-admin, docs/*.xml. -->
<!--   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-COVERAGE-GATE / V-M-COVERAGE-GATE. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Scope Summary - States source plan, delivered scope, and non-goals preserved. -->
<!--   Commit Summary - Lists local commits produced for the milestone. -->
<!--   QA Audit Summary - Records traceability, security/privacy, UX/runtime, and GRACE/generated audit outcomes. -->
<!--   Verification Summary - Lists final command evidence. -->
<!--   Blocker Ledger - Records blockers, accepted risks, and follow-ups. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.0 - Added final no-push readiness packet for web-admin login route guard. -->
<!-- END_CHANGE_SUMMARY -->

# Web-admin Login Route Guard Readiness

## Scope Summary

Source plan: `docs/superpowers/plans/2026-06-07-web-admin-login-route-guard.md`.

Delivered scope:

- Split admin auth GraphQL operation documents into one operation per file and regenerated web-admin types.
- Added typed admin auth model helpers, GraphQL client helpers, `AuthProvider`, protected route layout, public `/login` page, safe return-to behavior, and sidebar logout.
- Converted web-admin browser e2e from direct cookie installation to real UI login/logout for protected browser flows.
- Synchronized `docs/requirements.xml`, `docs/development-plan.xml`, `docs/knowledge-graph.xml`, and `docs/verification-plan.xml`.
- Added coverage branch tests and a durable coverage matrix/evidence packet.

Non-goals preserved:

- No backend auth redesign.
- No public registration, forgot-password, admin-management UI, role UI, MFA, OAuth, or SSO.
- No public web behavior change beyond verification coverage.
- No product tenant or resource-owner authorization was added.
- No browser token, session id, raw cookie, or password storage was introduced.
- Credentialed GraphQL transport remains cookie-based through the existing shared client.

## Commit Summary

| Commit    | Summary                                                                          |
| --------- | -------------------------------------------------------------------------------- |
| `26c74bd` | Split admin auth operation documents and regenerated web-admin codegen contract. |
| `f11ac91` | Added admin auth model and GraphQL client helpers.                               |
| `3bd029e` | Added AuthProvider and protected admin route layout.                             |
| `50a35d5` | Added public admin login page.                                                   |
| `a9e791e` | Wired sidebar current admin identity and logout action.                          |
| `a8d146e` | Converted web-admin e2e to real browser login/logout.                            |
| `0b184f5` | Cleaned login-page lint warning.                                                 |
| `7f1fd0d` | Updated GRACE login guard contracts.                                             |
| `1d847bd` | Added branch coverage for auth guard behavior.                                   |
| `0f22ef6` | Recorded coverage evidence packet.                                               |
| `6791680` | Refined security evidence wording for test-only e2e credentials.                 |

## QA Audit Summary

| Audit area                  | Result | Evidence                                                                                                                                                                                                                                                                                                                                     |
| --------------------------- | ------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Source-plan traceability    | PASS   | `.tasks/web-admin-login-route-guard/verification.md` maps each source-plan obligation to code/tests/docs/evidence.                                                                                                                                                                                                                           |
| Semantic-loss preservation  | PASS   | Non-goals above remain true; no scope expansion was introduced in `apps/api`, `apps/web`, or product auth.                                                                                                                                                                                                                                   |
| Security and privacy        | PASS   | No client-side auth token/session storage; test-only e2e credentials are environment-backed defaults; GraphQL client has no bearer helper; raw `Set-Cookie` assertion output was redacted to a boolean cookie-name check; the unused browser cookie-install helper was removed; logout clears backend session and protected frontend caches. |
| Redirect safety             | PASS   | `resolveSafeReturnTo` rejects absolute, protocol-relative, backslash, malformed, relative, `/login`, `/login/`, and login-loop targets.                                                                                                                                                                                                      |
| UX/runtime                  | PASS   | `/login` renders outside the sidebar shell; protected routes show accessible loading state, avoid protected-content flash, preserve safe `from`, and render stable errors.                                                                                                                                                                   |
| Sidebar logout              | PASS   | Sidebar shows real admin name/email/initials; logout is discoverable, disables/relabels while pending, revokes backend session, clears protected caches, and returns to `/login`.                                                                                                                                                            |
| Browser proof               | PASS   | Playwright verifies unauthenticated redirect, UI login, logout, fresh-page post-logout redirect, users CRUD, duplicate email, desktop sidebar, and mobile sheet flows.                                                                                                                                                                       |
| GRACE/generated consistency | PASS   | Shared XML docs parse; `grace lint --path .` exits 0; generated auth operation types exist; no stale aggregate auth document refs remain in key GRACE docs.                                                                                                                                                                                  |

## Verification Summary

- `bunx nx run web-admin:codegen`: PASS.
- `bunx nx run web-admin:typecheck`: PASS.
- `bunx nx test web-admin`: PASS before coverage additions, 18 files and 74 tests.
- `bunx nx lint web-admin`: PASS.
- `bunx nx build web-admin`: PASS with existing Vite large-chunk warning only.
- `bunx nx run web-admin:e2e`: PASS on isolated `mt-login-guard`, 7 Playwright tests.
- `bun run --cwd apps/web-admin test -- src/entities/admin-auth/model.test.ts src/entities/admin-auth/provider.test.tsx src/App.test.tsx`: PASS, 3 files and 18 tests after coverage additions.
- `bunx nx run web-admin:test-coverage`: PASS, 18 files and 78 tests, 100 percent statements, branches, functions, and lines.
- `bun run test:coverage`: PASS on isolated `mt-login-coverage`, `[Coverage][gate] all thresholds passed`.
- `bun run verify:coverage`: PASS on isolated `mt-login-verify` after cleaning milestone compose scopes; covered lint, codegen, typecheck, build, coverage, public web e2e, web-admin e2e, XML, and GRACE lint.
- Post-review `bunx nx run web-admin:e2e`: PASS on isolated `mt-login-verify`, 7 Playwright tests after removing the browser cookie-install helper.
- Post-review `bunx nx lint web-admin`: PASS.
- Post-review `bunx nx run web-admin:typecheck`: PASS.
- `rg -n "adminAuth\\.graphql" docs/development-plan.xml docs/knowledge-graph.xml docs/verification-plan.xml`: PASS, no matches.
- `xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml`: PASS.
- `grace lint --path .`: PASS with 16 existing heuristic-export warnings in project skill scripts and the historical `.worktrees/mt-0xn-sqlc` copy.
- `git diff --check`: PASS.

## Blocker Ledger

Blockers: none.

Accepted risks:

- The existing Vite build emits a large-chunk warning; it is not caused by this feature and does not fail the build.
- `grace lint --path .` still reports 16 heuristic-export warnings in project skill scripts and a historical worktree copy; no new GRACE errors were introduced.
- During the first `verify:coverage` attempt Docker could not allocate a new compose network because local predefined address pools were exhausted. Milestone compose scopes were cleaned and the full gate passed on rerun.

Follow-ups:

- Future protected React Query keys must be added to `isProtectedAdminQueryKey` or moved under a shared protected key namespace when new protected admin data surfaces are introduced.
- Review finding `low`: future protected query-key namespace handling is accepted as the follow-up above because current protected query keys are `admin-users` and `admin-user`, both covered by logout cleanup.

Delivery state:

- No branch was pushed.
- No MR was created.
- Local milestone state is implementation-ready under the requested no-push contour.
