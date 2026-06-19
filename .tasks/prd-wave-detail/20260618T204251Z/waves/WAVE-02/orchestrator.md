# WAVE-02 Orchestrator

## Run Metadata
- Run ID: 20260618T204251Z
- Wave ID: WAVE-02
- Wave Name: Exercise Library - CRUD exercises with media
- Source Wave: docs/prd-waves/waves/wave-02.md (user-approved)
- Prior Detailed Wave: WAVE-01 (docs/prd-wave-details/waves/wave-01.md, ready-for-dev)
- Staging Folder: .tasks/prd-wave-detail/20260618T204251Z/staging/prd-wave-details

## Planner Scopes
1. product-ac
2. architecture-codebase
3. data-integration-ops
4. security-compliance
5. testing-exit
6. sequencing-fit

## Reviewer Perspectives
1. product-scope-and-ac
2. architecture-codebase-fit
3. data-api-integration-ops
4. security-privacy-compliance
5. testing-exit-criteria
6. sequencing-other-wave-fit
7. traceability-consistency

## Sources Read
- docs/prd-waves/waves/wave-02.md (selected source wave)
- docs/prd-wave-details/waves/wave-01.md (prior detailed wave)
- docs/technical-verified/api-contracts.md
- docs/technical-verified/data-contracts.md
- docs/technical-verified/architecture-and-boundaries.md
- docs/technical-verified/implementation-slices.md
- docs/technical-verified/auth-security-compliance.md
- docs/technical-verified/operations-observability.md
- docs/technical-verified/testing-and-delivery.md
- docs/technical-verified/integrations-and-events.md
- docs/technical-verified/technical-brief.md
- docs/technical-verified/client-state-and-ux-contracts.md
- docs/product-verified/domain-model.md
- docs/product-verified/acceptance-criteria.md
- docs/product-verified/functional-spec.md
- docs/product-verified/edge-cases.md
- docs/product-verified/business-rules.md
- docs/prd-waves/frontend-pages/page-002.md
- docs/prd-waves/frontend-pages/page-003.md
- apps/api/cmd/server/main.go
- apps/api/internal/appconfig/config.go
- apps/api/internal/graph/resolver.go
- apps/api/internal/graph/schema.resolvers.go
- apps/api/internal/graph/admin_auth.resolvers.go
- apps/api/internal/graph/admin_auth_helpers.go
- apps/api/internal/service/admin_auth.go
- apps/api/internal/repository/postgres/user_repo.go
- apps/api/internal/repository/postgres/admin_repo.go
- apps/api/internal/repository/postgres/migrations/00001_init.sql
- apps/api/internal/repository/postgres/migrations/00079_admin_users.sql
- apps/api/internal/repository/postgres/queries/users.sql
- apps/api/internal/repository/postgres/queries/admin_users.sql
- apps/api/internal/middleware/admin_auth.go
- apps/api/internal/handler/health.go
- apps/api/internal/handler/users.go
- apps/api/gqlgen.yml
- apps/api/sqlc.yaml
- libs/graphql/schema/schema.graphql
- libs/graphql/schema/admin_auth.graphql

## Key Codebase Patterns
- **Hybrid API**: GraphQL for CRUD, REST for binary uploads
- **Service/Repository pattern**: service layer with interface contracts, repository adapters with sqlc-generated code
- **GraphQL**: extend type Query/Mutation pattern, union result types (Success | ValidationError | AuthError | NotFoundError)
- **sqlc**: queries in .sql files under queries/, generated output in generated/
- **Migrations**: goose format with +goose Up/Down markers
- **Resolver pattern**: per-schema-file resolver, dependency injection through Resolver struct
- **Admin auth guard**: requireAdmin() helper checks context for AdminPrincipal
- **Log markers**: [Domain][action][BLOCK_NAME] pattern
- **Error envelope**: { "error": { "code", "message", "field"? } }

## Max Attempts
Max planner attempts: 3
Max review cycles: 3
Current cycle: 1