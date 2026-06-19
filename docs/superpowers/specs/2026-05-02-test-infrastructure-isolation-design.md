# Test Infrastructure Isolation Design

Date: 2026-05-02
Status: Approved for implementation planning

## Summary

The coverage and e2e gates must use dedicated test infrastructure instead of the normal local development database and Redis. Today, Playwright e2e starts PostgreSQL and Redis through `docker/docker-compose.dev.yml`, the API defaults point at `monorepo_dev` on port `7501`, and repository integration tests can truncate `users` when their fallback DSN targets the dev database.

This design changes the test contract so `test:coverage`, `test:e2e`, and `verify:coverage` use only an isolated test Docker Compose stack, an isolated `monorepo_test` database, a test Redis instance, and explicit environment overrides. A destructive test must fail before touching data if it is not pointed at the test database.

## Current Problem

The current gate is isolated by API and web ports, but not by persistence:

- `apps/web/e2e/preflight.mjs` starts services from `docker/docker-compose.dev.yml`.
- `docker/docker-compose.dev.yml` defaults to `monorepo_dev`, `mt-postgres`, port `7501`, and volume `pg-data`.
- `apps/web/e2e/playwright.config.ts` uses isolated API and web ports, but does not isolate database or Redis targets.
- `apps/api/config/config.yml` defaults the API to `localhost:7501/monorepo_dev`.
- `apps/api/internal/repository/postgres/user_repo_test.go` accepts `API_TEST_DATABASE_DSN`, but falls back to `monorepo_dev`.
- The repository integration test runs `TRUNCATE users RESTART IDENTITY CASCADE`, so a fallback to the dev database can destroy local development data.

## Goals

- Keep local development compose and test compose separate.
- Make `test:coverage`, `test:e2e`, and `verify:coverage` target only test infrastructure.
- Prevent destructive integration cleanup from running against `monorepo_dev`.
- Keep the existing 100 percent coverage policy intact.
- Preserve normal local `nx serve api` behavior through the existing development config.
- Record the changed contract in GRACE XML docs, README, and verification evidence during implementation.

## Non-Goals

- Do not replace the normal local development database or ports.
- Do not require dropping the test volume on every run.
- Do not lower coverage thresholds or add broad coverage exclusions.
- Do not redesign the user GraphQL domain or Playwright scenario matrix beyond the infrastructure target.

## Recommended Approach

Use a dedicated `docker/docker-compose.test.yml` stack as the only infrastructure source for coverage and e2e gates.

The rejected lighter approach was to keep one dev compose and create a separate database inside it. That would reduce file count, but the container, volume, and service lifecycle would still be shared with development. The rejected stricter approach was to make the test stack fully ephemeral on every run. That maximizes repeatability, but adds unnecessary local cost for this template.

The approved baseline is a dedicated test compose with separate names, ports, volumes, and database. Data cleanup remains in test helpers after safety guards prove the DSN targets `monorepo_test`.

## Architecture

`M-COVERAGE-GATE` changes from "dev compose plus isolated app ports" to "dedicated test infrastructure plus isolated app ports".

The development stack remains:

- Compose file: `docker/docker-compose.dev.yml`
- Database: `monorepo_dev`
- PostgreSQL container: `mt-postgres`
- Redis container: `mt-redis`
- Default ports: `7501` and `7502`
- Volume: `pg-data`

The test stack becomes:

- Compose file: `docker/docker-compose.test.yml`
- Database: `monorepo_test`
- PostgreSQL container: `mt-test-postgres`
- Redis container: `mt-test-redis`
- Default PostgreSQL port: `17501`
- Default Redis port: `17502`
- PostgreSQL volume: `pg-test-data`

The normal API config can stay development-oriented. Gate behavior must be driven by environment overrides from the coverage and e2e runners.

## Components

### Test Docker Compose

`docker/docker-compose.test.yml` defines PostgreSQL and Redis services with test-specific container names, ports, volumes, and `POSTGRES_DB=monorepo_test`.

### E2E Preflight

`apps/web/e2e/preflight.mjs` becomes the single Playwright infrastructure bootstrap path. It starts only the test compose services, waits for health, and exposes enough log context to prove the test stack was selected without printing secrets.

### Playwright Config

`apps/web/e2e/playwright.config.ts` starts:

- API on `18080`, with `SERVER_PORT`, `SERVER_CORS_ORIGINS`, `POSTGRES_*`, and `REDIS_*` pointed at the test stack.
- Web on `13000`, with `NEXT_PUBLIC_API_URL` pointed at the isolated API.

The config should avoid running preflight twice through both global setup and webServer command.

### Coverage Runner

`tools/coverage/preflight.mjs` validates that the test compose and required scripts/configs exist.

`tools/coverage/run.mjs` starts the test infrastructure before Go coverage and passes:

- `COVERAGE_GATE=1`
- `API_TEST_DATABASE_DSN` pointing to `monorepo_test`
- test Redis host and port overrides where repository tests need Redis

### Repository Integration Tests

PostgreSQL and Redis integration tests become env-first and gate-strict. Under `COVERAGE_GATE=1`, missing or unsafe test infrastructure is a failure, not a skip.

Destructive cleanup such as `TRUNCATE users RESTART IDENTITY CASCADE` may run only after a safety guard verifies the DSN is a test DSN.

## Data Flow

### `bun run test:e2e`

1. Nx runs `web:e2e`.
2. Playwright starts e2e preflight.
3. Preflight starts PostgreSQL and Redis from `docker/docker-compose.test.yml`.
4. The API starts on `18080` with environment overrides targeting `localhost:<test-postgres-port>/monorepo_test` and the test Redis port.
5. The API runs migrations against `monorepo_test`.
6. Playwright verifies health, readiness, GraphQL CRUD, and browser user flows.
7. Any data cleanup touches only the test database.

### `bun run test:coverage`

1. `tools/coverage/preflight.mjs` validates the coverage and isolated infrastructure contract.
2. `tools/coverage/run.mjs` starts test infrastructure once.
3. Go coverage runs with `COVERAGE_GATE=1`, `API_TEST_DATABASE_DSN=postgres://.../monorepo_test?...`, and test Redis overrides.
4. Web and tools coverage run as they do today.
5. Coverage thresholds remain 100 percent.

### `bun run verify:coverage`

`verify:coverage` continues to run lint, codegen, web typecheck, build, coverage, e2e, XML validation, and `grace lint --path .`. Its coverage and e2e phases inherit the isolated test infrastructure contract.

## Failure Rules

- If `COVERAGE_GATE=1` and `API_TEST_DATABASE_DSN` is empty, the gate fails.
- If `COVERAGE_GATE=1` and the DSN points at `monorepo_dev`, the gate fails.
- If `COVERAGE_GATE=1` and the DSN uses the dev PostgreSQL port `7501`, the gate fails.
- If the e2e API environment points at `monorepo_dev`, preflight or config validation fails before destructive tests run.
- If repository cleanup is about to truncate data and the DSN does not identify `monorepo_test`, the test fails before cleanup.
- If test PostgreSQL or Redis is unavailable during a gate run, the gate fails instead of skipping.
- Ordinary local runs of individual integration tests may still skip when optional services are unavailable, but their fallback must be a test target or an explicit opt-in environment value, not `monorepo_dev`.

## Testing Strategy

Implementation must add focused coverage for unsafe target rejection:

- Empty DSN under `COVERAGE_GATE=1`.
- DSN targeting `monorepo_dev`.
- DSN targeting dev port `7501`.
- Malformed DSN.
- Cleanup guard rejects unsafe targets before `TRUNCATE`.

Focused verification after implementation:

- `bun run test:coverage`
- `bun run test:e2e`

Final verification after source and docs are synchronized:

- `bun run verify:coverage`
- `xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml`
- `grace lint --path .`

## Documentation And GRACE Updates

Implementation must update:

- `docs/requirements.xml` so the e2e/coverage contract names dedicated test infrastructure instead of dev compose.
- `docs/technology.xml` so the e2e policy points to `docker/docker-compose.test.yml`.
- `docs/development-plan.xml` so `M-COVERAGE-GATE` includes isolated test infrastructure and unsafe-target guards.
- `docs/knowledge-graph.xml` so CoverageGate annotations name test compose and `monorepo_test`.
- `docs/verification-plan.xml` so failure scenarios explicitly assert that destructive gate tests never target the dev database.
- `README.md` so developers know local development uses dev compose while gates use test compose.

## Acceptance Criteria

- `test:coverage`, `test:e2e`, and `verify:coverage` do not use `docker/docker-compose.dev.yml` for test infrastructure.
- Gate-launched API and repository tests target `monorepo_test`, not `monorepo_dev`.
- Test PostgreSQL and Redis use separate container names, ports, and volumes from the development stack.
- Destructive repository cleanup cannot run before the safe test DSN guard passes.
- The existing local development flow remains on `docker/docker-compose.dev.yml` and `monorepo_dev`.
- GRACE XML docs, README, source, tests, and verification evidence are synchronized before implementation closeout.
