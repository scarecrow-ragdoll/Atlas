# Codebase Fit

## Relevant Modules
- M-API (GoSharedHTTPAPI): apps/api — existing Go HTTP server with chi router, gqlgen GraphQL, goose migrations, sqlc queries. WAVE-01 extends this with fitness-domain schema, resolvers, services, migrations.
- M-GO-CONFIG: libs/go/config — shared config loading. WAVE-01 adds fitness-app config fields (PIN settings, session TTL, media paths).
- M-GO-LOGGER: libs/go/logger — shared zap logger. WAVE-01 uses existing logger, no changes.
- M-GRAPHQL-SCHEMA: libs/graphql/schema — admin GraphQL SDL. WAVE-01 adds fitness-domain schema types alongside existing admin schema.
- M-WORKSPACE: package.json, nx.json — WAVE-01 extends with codegen targets for fitness-domain gqlgen/sqlc.
- M-WEB-ADMIN / M-WEB: existing frontends — not touched by WAVE-01.

## Relevant Files Read
- apps/api/cmd/server/main.go: wires config, logger, storage, middleware, routes, GraphQL. WAVE-01 adds fitness-domain wiring.
- apps/api/internal/appconfig/config.go: config struct. WAVE-01 adds Session, Pin, Media config sections.
- apps/api/internal/handler/health.go: existing health endpoint. No change.
- apps/api/internal/middleware/admin_auth.go: cookie-based admin session middleware. WAVE-01 adds separate PIN auth middleware.
- apps/api/internal/service/admin_auth.go: admin auth service. WAVE-01 adds fitness-domain PIN service.
- apps/api/internal/repository/redis/admin_session_store.go: HMAC-keyed Redis sessions. WAVE-01 reuses or extends for PIN sessions.
- apps/api/gqlgen.yml: gqlgen config. WAVE-01 adds fitness-domain schema paths.
- apps/api/sqlc.yaml: sqlc config. WAVE-01 adds fitness-domain query paths.
- libs/graphql/schema/schema.graphql: root schema. WAVE-01 adds fitness-domain type extensions.
- docker-compose.yml: root Docker Compose. WAVE-01 adds fitness service configuration.
- package.json: root scripts. WAVE-01 adds fitness codegen commands.

## Public Contracts
- M-API contract: serves HTTP on configurable port, health at /health, GraphQL at /graphql, admin cookie session protected. WAVE-01 preserves all existing contracts and adds fitness-domain operations under the same /graphql endpoint with PIN session protection.
- M-GRAPHQL-SCHEMA contract: admin auth types and user CRUD. WAVE-01 extends with fitness-domain types without modifying existing admin operations.

## Generated Artifact Impact
- gqlgen: adding fitness-domain .graphql files will regenerate Go resolver interfaces and models. Existing admin resolvers unaffected.
- sqlc: adding fitness-domain SQL queries will regenerate Go query code. Existing admin queries unaffected.
- web-admin codegen: admin GraphQL client regenerated when schema changes. WAVE-01 adds fitness-domain types but admin client only needs admin operations — generated types may include new types but admin pages don't use them.
- No frontend codegen impact expected for existing pages.

## Integration Points
- The admin auth middleware (cookie + HMAC Redis session) and PIN auth middleware (PIN hash + Redis session) are separate middleware chains. The GraphQL endpoint will use a directive or middleware that routes requests to PIN-authenticated resolvers for fitness-domain operations vs. cookie-authenticated resolvers for admin operations.
- Docker Compose extensions: add media volume, worker service placeholder for async jobs.

## Likely Graph Deltas
- M-API: fitness-domain service, repository, resolver, middleware packages added alongside existing admin packages.
- M-GRAPHQL-SCHEMA: fitness-domain schema files added to libs/graphql/schema/.
- knowledge-graph.xml: add fitness domain entities, new service/repository modules.
- verification-plan.xml: add fitness-domain verification flows (VF-PIN-AUTH, VF-SETTINGS, VF-FOUNDATION).

## Unsupported Assumptions
- The PIN session lifespan and refresh strategy are not specified in verified docs. Recommend 7-day session with configurable TTL per TDEC-AUTH-001.
- Rate limiting for PIN attempts is deferred (Q-PIN-001). First implementation: no rate limit, documented as risk.
- The API protocol decision (TDEC-001) confirmed hybrid model: GraphQL for CRUD, REST for binary uploads. WAVE-01 sets up both paths.
- Session storage: Redis already exists for admin sessions. PIN sessions reuse the same Redis instance.