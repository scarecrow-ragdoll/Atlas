<!-- FILE: .tasks/swap-web-admin-vite-web-next/readiness.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Summarize final no-push readiness for the web-admin Vite GraphQL and public web Next REST swap milestone. -->
<!--   SCOPE: Records completed scope, command evidence, QA audit outcomes, blockers, risks, follow-ups, and no-push/no-MR status; excludes implementation details already captured in source files. -->
<!--   DEPENDS: docs/superpowers/plans/2026-06-05-swap-web-admin-vite-web-next.md, .tasks/swap-web-admin-vite-web-next/verification.md, docs/*.xml, apps/web-admin, apps/web, docker, deploy/dokploy, tools/ci. -->
<!--   LINKS: M-WEB-ADMIN / M-WEB / M-COVERAGE-GATE / M-CI-CD / V-M-WEB-ADMIN / V-M-WEB / V-M-COVERAGE-GATE / V-M-CI-CD. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Source And Scope - Identifies the source plan and completed ownership swap. -->
<!--   Command Evidence - Summarizes verification gates and where detailed evidence lives. -->
<!--   QA Review - Records independent subagent audit outcomes and the fixed finding. -->
<!--   Risks And Follow-Ups - Lists blockers, accepted non-blocking warnings, and follow-up Bead IDs. -->
<!--   Delivery Status - States local/no-push/no-MR status. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.0 - Created final no-push readiness packet for mt-f82.3.4. -->
<!-- END_CHANGE_SUMMARY -->

# Web Swap No-Push Readiness Packet

Source plan: `docs/superpowers/plans/2026-06-05-swap-web-admin-vite-web-next.md`

Detailed verification ledger: `.tasks/swap-web-admin-vite-web-next/verification.md`

## Source And Scope

Completed scope:

- `apps/web-admin` remains Nx project `web-admin` and is now a Vite + React Router admin SPA using GraphQL documents from `apps/web-admin/src/entities/user/api/*.graphql`.
- `apps/web` remains Nx project `web` and is now a Next App Router public REST app.
- Public browser code uses same-origin `/api/users`; Next server/runtime proxy code uses `WEB_API_BASE_URL`.
- Existing deployable `web` image is retargeted to public Next web. Vite admin deployment remains explicitly out of scope for this no-push milestone.
- Coverage preflight, coverage allowlists, e2e configs, Docker/Dokploy/local compose, CI helper env rendering, GRACE XML, and file-local GRACE markup were synchronized.

Primary implementation commits:

- Implementation: `6a9e6db`, `97f54f8`, `e693661`, `e000a7c`, `23bcd26`, `f707b8f`, `a8742e6`, `6b611c9`.
- Full coverage evidence: `9b32a52`, `bbdf56a`, `541ba60`, `58e4cec`, `9fd4bdb`.
- Pre-MR QA evidence: `502a9eb`, `f58e636`, `3ca02cb`.

## Command Evidence

All final gates passed. Detailed command rows and output summaries are in `verification.md`.

| Gate                                        | Result                                                                                                                   |
| ------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------ |
| `bun install --frozen-lockfile`             | PASS                                                                                                                     |
| `bun run codegen` plus generated diff check | PASS                                                                                                                     |
| `bunx nx run codegen:validate`              | PASS - 49 tool tests at 100 percent plus generated drift check                                                           |
| `node tools/coverage/preflight.mjs`         | PASS                                                                                                                     |
| `bun run test:coverage`                     | PASS - Go, admin Vite, public Next, tools, and coverage thresholds                                                       |
| `bun run verify:coverage`                   | PASS - lint, codegen, typecheck, build, coverage, e2e, XML, and GRACE                                                    |
| `bun run build`                             | PASS                                                                                                                     |
| `bunx nx run web-admin:e2e`                 | PASS - 4 Chromium tests after GraphQL markup fix                                                                         |
| `bunx nx run web:e2e`                       | PASS - 1 Chromium REST test via `verify:coverage` and focused coverage evidence                                          |
| Docker build/run/curl smoke                 | PASS - `monorepo-template-web-swap:mt-f82-2-5`, `REST Web`, `Failed to load users`, container cleaned up                 |
| `xmllint --noout docs/*.xml`                | PASS                                                                                                                     |
| `grace lint --path .`                       | PASS - 0 errors, 16 known heuristic warnings                                                                             |
| Task 8 file-local GRACE scan                | PASS - 49 touched governed files                                                                                         |
| Forbidden ownership scans                   | PASS - no public GraphQL ownership, stale public Vite env, stale `NEXT_PUBLIC_*` API env, or stale admin image ownership |
| Local and Dokploy compose render checks     | PASS - `WEB_API_BASE_URL` present; no rendered `web-admin` or `NEXT_PUBLIC` strings                                      |

## QA Review

Independent read-only audits:

- Volta: PASS for source-plan traceability and review-loop fix preservation.
- Hume: initial FAIL for missing file-local GRACE headers in three GraphQL docs, admin users e2e spec, and `.gitlab-ci.yml`; fixed in `f58e636`, then verified with exact Task 8 scan, codegen drift, admin e2e, XML, and GRACE lint.
- Boole: PASS for deployment/runtime env and forbidden ownership boundaries.

## Risks And Follow-Ups

Blockers: none.

Follow-up Bead IDs: none created.

Accepted non-blocking notes:

- `grace lint --path .` still reports 16 known `analysis.heuristic-export-surface` warnings in project-local PRD/wave skill Python scripts and duplicated `.worktrees/mt-0xn-sqlc` copies. These are unrelated to the web swap and remain warnings, not errors.
- `docker compose -f docker/docker-compose.yml config` reports Docker Compose's obsolete top-level `version` warning. The rendered service/env contract is correct.
- Partial local Dokploy template rendering can warn about unrelated missing deployment env vars when only the web-swap variables are supplied for inspection. The rendered `WEB_API_BASE_URL` contract was verified.

## Delivery Status

No branch was pushed and no MR was created as part of this milestone.

This packet is the local readiness handoff seed if an MR or rollout is requested later.
