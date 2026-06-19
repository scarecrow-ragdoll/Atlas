<!-- FILE: README.md -->
<!-- VERSION: 1.0.1 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Explain the template setup, local development commands, verification gates, and deployment overview. -->
<!--   SCOPE: Documents developer-facing command usage and stack orientation; excludes exhaustive GRACE artifact details. -->
<!--   DEPENDS: package.json, apps/*/project.json, docker/docker-compose.dev.yml, docs/verification-plan.xml. -->
<!--   LINKS: M-WORKSPACE / M-API / M-WEB-ADMIN / M-WEB / VF-LOCAL-DEV. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Quick Start - Shows initial dependency installation. -->
<!--   Local Development - Shows local infrastructure and dev server ports. -->
<!--   Nx Commands - Lists focused project commands. -->
<!--   Coverage and E2E Gate - Summarizes verification commands and artifacts. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.1 - Aligned local dev docs with Vite admin, Next public web, and separated ports. -->
<!-- END_CHANGE_SUMMARY -->

# Monorepo Template

Full-stack monorepo: Go API with admin GraphQL and public REST, Vite admin UI, Next.js public UI, PostgreSQL, and Redis.

## Prerequisites

| Tool                    | Version | Install                                         |
| ----------------------- | ------- | ----------------------------------------------- |
| Bun                     | 1.1+    | [bun.sh](https://bun.sh/)                       |
| Go                      | 1.25+   | [go.dev](https://go.dev/dl/)                    |
| Docker & Docker Compose | latest  | [docker.com](https://www.docker.com/)           |
| air                     | latest  | `go install github.com/air-verse/air@latest`    |
| golangci-lint           | latest  | [golangci-lint.run](https://golangci-lint.run/) |

## Quick Start

```bash
git clone <repo-url>
cd monorepo-template
bun install
```

## Local Development

```bash
# 1. Start infra (PostgreSQL :7501, Redis :7502)
docker compose -f docker/docker-compose.dev.yml up -d

# 2. Start API, admin frontend, and public frontend
bun run dev          # API :8090, web-admin :3100, web :3101
```

If you want to run each service separately:

```bash
# API with hot-reload
nx serve api          # Go API on :8090, auto-migrates DB on startup

# Admin frontend
nx serve web-admin    # Vite admin UI on :3100, uses /graphql

# Public frontend
nx serve web          # Next.js public UI on :3101, uses /api/users
```

> **Tip:** Add `alias nx="bunx nx"` to your `.zshrc` if `nx` is not on PATH.

Open the running frontend on its configured port to verify the admin GraphQL or public REST user flow.

## Full Stack via Docker

```bash
docker compose -f docker/docker-compose.yml up --build
```

The current Docker web image runs the public Next.js `web` app. Vite admin deployment is intentionally separate future work.

## Nx Commands

```bash
nx serve api              # Go API dev server (air hot-reload)
nx serve web-admin        # Vite admin dev server on :3100
nx serve web              # Next.js public dev server on :3101
nx test api               # Go tests with coverage
nx test web-admin         # Admin Vitest
nx test web               # Vitest
nx lint api               # golangci-lint
nx lint web-admin         # Admin ESLint
nx lint web               # ESLint
nx run web-admin:typecheck # tsc --noEmit for admin
nx run web:typecheck      # tsc --noEmit
nx run api:codegen        # gqlgen
nx run web-admin:codegen  # graphql-codegen for admin GraphQL
nx run web-admin:test-coverage # admin Vitest coverage threshold
nx run web:test-coverage  # web Vitest coverage threshold
nx run web-admin:e2e      # admin GraphQL Playwright e2e
nx run web:e2e            # public REST Playwright e2e
nx affected --target=test # run tests for affected projects
```

## Coverage and E2E Gate

The handoff gate is:

```bash
bun run verify:coverage
```

It runs lint, codegen, web typecheck, build, the 100 percent coverage gate, Playwright e2e, XML validation, and `grace lint --path .`.

Focused commands:

```bash
bun run test:coverage     # Go + web + tools coverage, all thresholds at 100%
bun run test:e2e          # Playwright e2e for web-admin and web, sequentially
```

`test:e2e` starts PostgreSQL and Redis from `docker/docker-compose.test.yml`, then starts the API on `18080` and runs the web-admin and web Playwright projects sequentially unless their `E2E_*` overrides are set. Local development still uses `docker/docker-compose.dev.yml` with `monorepo_dev`; coverage and e2e gates use the isolated `monorepo_test` stack on default ports `17501` and `17502`.

Artifacts:

- `dist/coverage/go/*/coverage.out`
- `dist/coverage/web-admin/coverage-summary.json`
- `dist/coverage/web/coverage-summary.json`
- `dist/coverage/tools/coverage-summary.json`
- `dist/test-results/web-admin-e2e`
- `dist/test-results/web-e2e`
- `dist/playwright-report/web-admin`
- `dist/playwright-report/web`

Generated files and bootstrap entrypoints are excluded only through `tools/coverage/coverage.config.json`, where every allowlist entry has a replacement gate such as codegen, build, typecheck, or Playwright startup coverage.

## CI/CD and Dokploy

The template uses GitLab CI and Dokploy for deployment. Merge requests run fast affected checks, `develop` deploys automatically to dev after the full gate, `main` validates release candidates, and production deploys only from SemVer tags such as `v1.2.3` after manual approval.

See [docs/infrastructure/ci-cd.md](docs/infrastructure/ci-cd.md) for required GitLab variables, Dokploy compose setup, image tagging, and the production release flow. See [docs/infrastructure/gitlab-runners.md](docs/infrastructure/gitlab-runners.md) for the shared department runner pool contract.

## Agent Tooling Setup

```bash
# Install agent-browser CLI (browser automation for E2E/UI testing)
npm install -g agent-browser
agent-browser install

# Optional GRACE CLI companion for the installed project GRACE skills
bun add -g @osovv/grace-cli
grace lint --path .
```

## Development Flow

Every feature goes through the GRACE contract pipeline:

```bash
$grace-status        # 1. Check project health and next safe action
$grace-ask           # 2. Ground questions in docs/*.xml and code
$grace-plan          # 3. Update module contracts, flows, and execution order
$grace-verification  # 4. Update tests, traces, and evidence expectations
$grace-execute       # 5. Implement sequentially, or use $grace-multiagent-execute for safe waves
$grace-refresh       # 6. Re-sync GRACE artifacts after code changes
```

GRACE skills are installed project-locally from `https://github.com/osovv/grace-marketplace` under `.agents/skills/grace-*` and `.claude/skills/grace-*`. See `AGENTS.md` for the full agent operating guide.
