# Source Inventory

## PRD Wave Sources
- docs/prd-waves/index.md
- docs/prd-waves/wave-map.md
- docs/prd-waves/open-questions.md
- docs/prd-waves/waves/index.md
- docs/prd-waves/waves/wave-01.md

## Frontend Pages Source
- docs/prd-waves/frontend-pages/index.md
- docs/prd-waves/frontend-pages/page-011.md

## Product Sources
- docs/product-verified/index.md
- docs/product-verified/product-brief.md
- docs/product-verified/functional-spec.md
- docs/product-verified/domain-model.md
- docs/product-verified/actors-and-permissions.md
- docs/product-verified/business-rules.md
- docs/product-verified/edge-cases.md
- docs/product-verified/scope.md
- docs/product-verified/user-flows.md
- docs/product-verified/acceptance-criteria.md

## Technical Sources
- docs/technical-verified/index.md
- docs/technical-verified/technical-brief.md
- docs/technical-verified/architecture-and-boundaries.md
- docs/technical-verified/api-contracts.md
- docs/technical-verified/data-contracts.md
- docs/technical-verified/auth-security-compliance.md
- docs/technical-verified/integrations-and-events.md
- docs/technical-verified/operations-observability.md
- docs/technical-verified/implementation-slices.md
- docs/technical-verified/testing-and-delivery.md
- docs/technical-verified/client-state-and-ux-contracts.md

## GRACE Sources
- docs/development-plan.xml
- docs/knowledge-graph.xml
- docs/verification-plan.xml

## Codebase Sources
- apps/api/cmd/server/main.go
- apps/api/internal/appconfig/config.go
- apps/api/internal/handler/health.go
- apps/api/internal/middleware/admin_auth.go
- apps/api/internal/service/admin_auth.go
- apps/api/internal/repository/redis/admin_session_store.go
- apps/api/internal/service/user_service.go
- apps/api/internal/repository/postgres/user_repo.go
- apps/api/gqlgen.yml
- apps/api/sqlc.yaml
- libs/go/config/config.go
- libs/go/logger/logger.go
- libs/graphql/schema/schema.graphql
- libs/graphql/schema/admin_auth.graphql
- docker-compose.yml
- package.json

## Source Delta
- N/A (initial detailed wave run)

## Source Gaps
- No fitness-domain GraphQL schema (admin auth schema only)
- No fitness-domain database migrations (goose seed migrations exist but no product entities)
- No PIN auth service (admin session-based auth exists for admin context)
- No settings service for fitness app context
- No Docker Compose fitness-service overrides