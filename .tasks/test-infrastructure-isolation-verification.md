# Test Infrastructure Isolation Verification

Date: 2026-05-03

## Scope

The coverage and e2e gates now use isolated PostgreSQL and Redis test infrastructure. Destructive repository cleanup is guarded before it can run against a development database.

## Commands

- `rg -n "monorepo_dev|7501|7502|docker-compose.dev.yml" apps/api/config apps/api/internal apps/web/e2e tools/coverage docker README.md docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml`: PASS; hits were limited to local development docs/config, isolated test port references, and guard or regression-test rejection text.
- `docker compose -f docker/docker-compose.test.yml config`: PASS
- `cd apps/web && bun run e2e:preflight`: PASS
- `cd apps/api && COVERAGE_GATE=1 API_TEST_DATABASE_DSN=postgres://app:secret@localhost:7501/monorepo_dev?sslmode=disable go test ./internal/repository/postgres -run TestUserRepo_CreateGetListUpdateDelete -count=1`: FAIL as expected with `unsafe postgres test DSN`
- `cd apps/api && COVERAGE_GATE=1 API_TEST_DATABASE_DSN=postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable TEST_REDIS_PORT=17502 go test ./internal/repository/postgres ./internal/repository/redis ./internal/testinfra`: PASS
- `bun run test:coverage`: PASS with `[Coverage][gate] all thresholds passed`
- `bun run test:e2e`: PASS
- `bun run verify:coverage`: PASS
- `xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml`: PASS
- `grace lint --path .`: PASS
- `git diff --check`: PASS

## Evidence Summary

Gate runners and repository integration tests target `monorepo_test` on the test PostgreSQL port. The explicit unsafe-DSN regression check fails before cleanup when pointed at `monorepo_dev` on the development port.
