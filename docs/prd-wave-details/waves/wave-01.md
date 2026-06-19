# Wave 01: Foundation

## Status
ready-for-dev

## User Approval
user-approved (2026-06-18). Ready for implementation.

## Source Wave Summary
WAVE-01 from docs/prd-waves/waves/wave-01.md. Establish core infrastructure: Docker Compose, Go API skeleton extension, PostgreSQL fitness tables, Redis PIN sessions, PIN auth guard, settings service, CI/CD structure. Source status: user-approved.

## Outcome After Implementation
- OUT-W01-001: Docker Compose stack extended with fitness app services, media volume, worker placeholder
- OUT-W01-002: Go API extended with fitness-domain GraphQL schema, resolvers, and middleware
- OUT-W01-003: PostgreSQL fitness tables created via goose migrations (settings, user sessions)
- OUT-W01-004: PIN session auth guard protecting fitness-domain API and GraphQL endpoints
- OUT-W01-005: Settings service for PIN config, AI context, and export preferences
- OUT-W01-006: Media upload/download REST endpoints scaffold schema

## Scope Included
- CAP-W01-001: Docker Compose extensions (media volume, worker service placeholder, environment)
- CAP-W01-002: Fitness-domain GraphQL schema in libs/graphql/schema (Settings, Session, Pin types)
- CAP-W01-003: Go API fitness packages (settings service, PIN auth service, session store)
- CAP-W01-004: PostgreSQL goose migrations for fitness tables (settings, sessions)
- CAP-W01-005: Redis PIN session store (reuses existing Redis, separate session namespace)
- CAP-W01-006: PIN guard middleware for fitness-domain GraphQL and REST endpoints
- CAP-W01-007: Settings GraphQL resolvers (read/write settings, PIN enable/disable/change)
- CAP-W01-008: Media REST scaffold (POST /api/v1/media/upload, GET /api/v1/media/{id})
- CAP-W01-009: Nx codegen targets for fitness-domain gqlgen and sqlc
- CAP-W01-010: Test infrastructure (fitness-domain test helper, test DB migrations)

## Scope Excluded
- Exercise CRUD (WAVE-02)
- Workout diary (WAVE-03)
- Cardio/body tracking (WAVE-04)
- Nutrition (WAVE-05)
- Charts (WAVE-06)
- AI export/review (WAVE-07, WAVE-08)
- Backup (WAVE-09)
- Admin auth changes (existing admin auth for M-WEB-ADMIN unchanged)
- Frontend pages, UI, UX, routes, navigation, components

## Dependencies And Other-Wave Fit
- No predecessor dependencies (first wave)
- WAVE-02 (Exercise Library): depends on this wave's db schema, API structure, media REST scaffold
- WAVE-03 (Workout Diary): depends on WAVE-01 db + WAVE-02
- WAVE-04 (Cardio/Body): depends on WAVE-01 db
- WAVE-05 (Nutrition): depends on WAVE-01 db
- No scope collision: this wave creates the framework only; domain tables owned by later waves

## Frontend Pages Dependencies
- PAGE-011 (Settings): depends on settings GraphQL queries/mutations for PIN toggle, PIN change, AI context form, export preferences. WAVE-01 provides the backend. No frontend pages, UI, or UX work in this wave.
- PAGE-001 (Dashboard): depends on data availability from db tables created here and later waves
- PAGE-004 (Cardio): depends on db tables from WAVE-04
- All pages depend on the PIN auth flow for API access. PIN guard middleware provides the auth layer.

## Codebase Fit And Touchpoints
- apps/api/cmd/server/main.go: add fitness service wiring alongside existing admin wiring
- apps/api/internal/appconfig/config.go: add SessionConfig, PinConfig, MediaConfig sections
- apps/api/internal/middleware/admin_auth.go: add parallel pin_auth.go middleware
- apps/api/internal/service/admin_auth.go: add parallel pin_service.go and settings_service.go
- apps/api/internal/repository/redis/admin_session_store.go: add pin_session_store.go (separate key prefix)
- apps/api/internal/repository/postgres/: add fitness sqlc queries and generated output
- apps/api/internal/handler/: add media_handler.go for REST upload/download
- libs/graphql/schema/schema.graphql: extend with fitness-domain types using extend type Query/Mutation
- libs/graphql/schema/: add fitness.graphql, settings.graphql, pin_auth.graphql
- apps/api/gqlgen.yml: add fitness-domain schema paths and model config
- apps/api/sqlc.yaml: add fitness-domain query paths
- docker-compose.yml: add media volume, worker placeholder service
- package.json or nx.json: add fitness codegen targets
- apps/api/internal/graph/: add fitness-domain resolver package

## Design Contracts
- PIN auth: session-based, opaque HMAC key in Redis, configurable TTL (default 7 days). No refresh token — re-authenticate on expiry.
- PIN storage: bcrypt hash in PostgreSQL settings table. No raw PIN stored.
- Settings: single-row per user, JSON-like key-value with typed Go struct. Extended in later waves.
- API protocol: hybrid GraphQL (CRUD) + REST (binary uploads). Both under PIN session auth.
- Error format: { "error": { "code": "ERROR_CODE", "message": "Human readable" } } per TDEC-027.
- Media storage: local filesystem under configured media path, served by REST endpoint. Pluggable to S3 in future.

## Data API Integration And Operations
- PostgreSQL: goose migration files in apps/api/migrations/. Seed migration for settings defaults.
- Redis: shared Redis instance with separate key prefix (fitness:session: vs admin:session:).
- GraphQL: single /graphql endpoint. Admin operations protected by admin auth middleware. Fitness operations protected by PIN auth middleware. Routing via GraphQL directive or middleware type check.
- REST: /api/v1/media/upload (multipart POST), /api/v1/media/{id} (GET). PIN-protected.
- Logging: reuse existing zap logger middleware. Add fitness-domain log markers [PinAuth], [Settings], [Media].
- Metrics: no custom metrics in WAVE-01. Standard chi request logging.
- Operations: docker compose up for local dev. Existing Docker Compose stack extended.

## Security Privacy And Compliance
- PIN: bcrypt hash, opaque session in Redis, configurable TTL.
- Rate limiting: deferred (DQ-W01-001). First implementation: no rate limit. Redis TTL-based rate limiting recommended for future.
- No PII stored in settings beyond user-provided AI context (goal, height, age, etc.) — stored locally, no external transmission.
- Media files: stored on local filesystem, served only through authenticated REST endpoint. No public access.
- Admin auth (cookie-based) and PIN auth (token-based) are separate concerns. No credential crossover.

## Implementation Slices

| Slice ID | Name | Description |
| --- | --- | --- |
| SLICE-W01-001 | Db migrations | Create goose migrations for settings and pin_sessions tables. |
| SLICE-W01-002 | GraphQL schema | Add fitness-domain types: Settings, PinAuth, Session. Extend root Query/Mutation. |
| SLICE-W01-003 | PIN auth service | PIN verification (bcrypt), session create/validate/revoke in Redis. |
| SLICE-W01-004 | PIN guard middleware | Extract PIN token from Authorization header, validate session, inject principal context. |
| SLICE-W01-005 | Settings service | CRUD for fitness app settings (PIN hash, AI context, export preferences). |
| SLICE-W01-006 | GraphQL resolvers | Settings query/mutation, PIN enable/disable/change mutations. |
| SLICE-W01-007 | Media REST handler | POST /api/v1/media/upload (multipart), GET /api/v1/media/{id}. File system storage. |
| SLICE-W01-008 | Config extensions | Add SessionConfig, PinConfig, MediaConfig to API config struct. |
| SLICE-W01-009 | Codegen config | Update gqlgen.yml and sqlc.yaml for fitness-domain paths. |
| SLICE-W01-010 | Docker Compose | Add media volume, worker placeholder service, environment variables. |
| SLICE-W01-011 | Test infrastructure | Test DB migrations, test helpers for PIN auth and settings. |

## Acceptance Criteria

| AC ID | Description |
| --- | --- |
| AC-W01-001 | Settings table created via goose migration with fields: pin_hash, ai_goal, ai_height, ai_age, ai_experience, ai_split, ai_limits, ai_progression, ai_nutrition_strategy, default_export_weeks. |
| AC-W01-002 | PIN can be enabled: POST mutation sets pin_hash via bcrypt, returns success. |
| AC-W01-003 | PIN can be changed: requires valid current PIN, replaces pin_hash. |
| AC-W01-004 | PIN can be disabled: requires valid current PIN, clears pin_hash. |
| AC-W01-005 | PIN session created on successful PIN verification, stored in Redis with configurable TTL. |
| AC-W01-006 | PIN session validated on every fitness-domain GraphQL and REST request. Invalid/expired session returns auth error. |
| AC-W01-007 | Settings can be read (requires valid PIN session). |
| AC-W01-008 | Settings can be updated (requires valid PIN session). |
| AC-W01-009 | Media file uploaded via POST /api/v1/media/upload returns media ID. Stored under configured media path. |
| AC-W01-010 | Media file downloaded via GET /api/v1/media/{id} returns file content with correct content type. |
| AC-W01-011 | 404 returned for missing media file. |
| AC-W01-012 | Existing admin auth (cookie-based, admin_session_store) continues to work unchanged. |
| AC-W01-013 | Existing health endpoint at /health unaffected. |
| AC-W01-014 | All fitness-domain operations use PIN auth scope. Admin operations continue using admin session auth scope. |

## Exit Criteria

| EC ID | Description |
| --- | --- |
| EC-W01-001 | All acceptance criteria passing in focused unit and integration tests. |
| EC-W01-002 | Codegen produces valid Go code (gqlgen, sqlc) without drift. |
| EC-W01-003 | Existing admin auth test suite still passes unchanged. |
| EC-W01-004 | Existing health endpoint test passes unchanged. |
| EC-W01-005 | Docker Compose stack starts without errors with new services. |
| EC-W01-006 | PIN session TTL configurable via environment variable, defaults to 7 days. |
| EC-W01-007 | Lint passes for all changed packages |
| EC-W01-008 | Typecheck passes for Go API |
| EC-W01-009 | No frontend scope, UI, UX, component, or page changes in this wave |

## Verification Obligations

| Test ID | Description | Type | Command |
| --- | --- | --- | --- |
| TEST-W01-001 | Settings repository unit tests | unit | bunx nx run api:test -- --run '(?i)settings_repo' |
| TEST-W01-002 | PIN auth service unit tests | unit | bunx nx run api:test -- --run '(?i)pin_service' |
| TEST-W01-003 | PIN session store unit tests | unit | bunx nx run api:test -- --run '(?i)pin_session' |
| TEST-W01-004 | PIN guard middleware unit tests | unit | bunx nx run api:test -- --run '(?i)pin_auth' |
| TEST-W01-005 | Settings GraphQL resolver integration tests | integration | bunx nx run api:test -- --run '(?i)settings_resolver' |
| TEST-W01-006 | Media REST handler tests | integration | bunx nx run api:test -- --run '(?i)media_handler' |
| TEST-W01-007 | Migration smoke test (up + down) | integration | bunx nx run api:test -- --run '(?i)migration' |
| TEST-W01-008 | Admin auth regression tests | unit | bunx nx run api:test -- --run '(?i)admin_auth' |
| TEST-W01-009 | Health endpoint unchanged | unit | bunx nx run api:test -- --run '(?i)health' |
| TEST-W01-010 | Go lint for API package | lint | bunx nx run api:lint |
| TEST-W01-011 | GraphQL schema validate | codegen | bunx nx run graphql:validate |
| TEST-W01-012 | Codegen drift check (gqlgen + sqlc) | codegen | bunx nx run api:codegen && bunx nx run graphql:codegen |

## Rollout Rollback And Compatibility
- Rollout: merge PR, CI builds and runs tests, deploy via Dokploy compose update. New services start alongside existing.
- Rollback: revert PR, CI builds previous image, Dokploy compose update rolls back. Existing admin auth and health unaffected.
- Compatibility: all new operations are additive. No existing API changes. Client (web-admin) sees no functional change. Fitness frontend (web pages) will use new endpoints when built in later waves.
- Migration: goose migrations run at startup. Down migration available for rollback.

## Handoff Packets
- HANDOFF-W01-001: This wave brief document
- HANDOFF-W01-002: Planner reports (6 scopes)
- HANDOFF-W01-003: Reviewer evidence (8 perspectives)
- HANDOFF-W01-004: Final fit review evidence

## Reviewer Verdicts

| Wave | Perspective | Attempt | Verdict | Reviewer Report | Required Revisions | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| WAVE-01 | product-scope-and-ac | 1 | approved | .tasks/prd-wave-detail/20260618T203120Z/waves/WAVE-01/reviewers/review-product-scope-and-ac-1.md | none | AC covers all foundation scope |
| WAVE-01 | architecture-codebase-fit | 1 | approved | .tasks/prd-wave-detail/20260618T203120Z/waves/WAVE-01/reviewers/review-architecture-codebase-fit-1.md | none | Codebase fit well-documented |
| WAVE-01 | data-api-integration-ops | 1 | approved | .tasks/prd-wave-detail/20260618T203120Z/waves/WAVE-01/reviewers/review-data-api-integration-ops-1.md | none | Data/API/ops coverage adequate |
| WAVE-01 | security-privacy-compliance | 1 | approved | .tasks/prd-wave-detail/20260618T203120Z/waves/WAVE-01/reviewers/review-security-privacy-compliance-1.md | none | PIN bcrypt, Redis sessions, rate limiting noted as deferred |
| WAVE-01 | testing-exit-criteria | 1 | approved | .tasks/prd-wave-detail/20260618T203120Z/waves/WAVE-01/reviewers/review-testing-exit-criteria-1.md | none | 12 test obligations cover all AC and EC |
| WAVE-01 | sequencing-other-wave-fit | 1 | approved | .tasks/prd-wave-detail/20260618T203120Z/waves/WAVE-01/reviewers/review-sequencing-other-wave-fit-1.md | none | Dependency order correct, no collision |
| WAVE-01 | traceability-consistency | 1 | approved | .tasks/prd-wave-detail/20260618T203120Z/waves/WAVE-01/reviewers/review-traceability-consistency-1.md | none | Source traceability documented per section |
| WAVE-01 | final-wave-fit-review | 1 | approved | .tasks/prd-wave-detail/20260618T203120Z/waves/WAVE-01/reviewers/review-final-wave-fit-review-1.md | none | Package is ready-for-dev |

## Open Questions

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W01-001 | WAVE-01 | security | deferred | Q-PIN-001 | PIN rate limiting implementation? | Security against brute force attacks | Decide rate limit strategy (Redis TTL, fixed delay, lockout) | docs/technical-verified/auth-security-compliance.md | open | deferred |

## Traceability
- docs/prd-waves/waves/wave-01.md: source wave boundary, outcomes, capability groups
- docs/product-verified/functional-spec.md: settings, PIN behavior specifications
- docs/product-verified/actors-and-permissions.md: user roles, PIN ownership
- docs/product-verified/domain-model.md: Settings, Session entities
- docs/technical-verified/api-contracts.md: hybrid GraphQL/REST protocol decision
- docs/technical-verified/architecture-and-boundaries.md: system context, component boundaries
- docs/technical-verified/auth-security-compliance.md: PIN auth, session management
- docs/technical-verified/implementation-slices.md: Slice 0 foundation mapping
- docs/technical-verified/data-contracts.md: database entity contracts
- docs/development-plan.xml: M-API, M-GO-CONFIG, M-GRAPHQL-SCHEMA module contracts
- docs/knowledge-graph.xml: existing module boundaries and cross-links
- apps/api/internal: existing codebase patterns for service/repository/middleware/handler structure
- docs/prd-waves/frontend-pages/page-011.md: settings page backend dependencies