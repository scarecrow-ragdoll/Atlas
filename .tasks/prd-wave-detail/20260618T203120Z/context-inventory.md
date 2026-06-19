# Context Inventory: WAVE-01

## PRD Wave Sources
- docs/prd-waves/index.md (status: waves-approved)
- docs/prd-waves/wave-map.md (9 backend waves)
- docs/prd-waves/waves/wave-01.md (Foundation source wave)
- docs/prd-waves/frontend-pages/index.md (11 frontend pages)
- docs/prd-waves/frontend-pages/page-011.md (Settings - PIN, AI context, export preferences)

## Product Sources
- docs/product-verified/functional-spec.md
- docs/product-verified/domain-model.md
- docs/product-verified/actors-and-permissions.md
- docs/product-verified/business-rules.md
- docs/product-verified/edge-cases.md

## Technical Sources
- docs/technical-verified/implementation-slices.md (Slice 0 = Foundation)
- docs/technical-verified/api-contracts.md (hybrid GraphQL/REST, TDEC decisions resolved)
- docs/technical-verified/architecture-and-boundaries.md
- docs/technical-verified/data-contracts.md
- docs/technical-verified/auth-security-compliance.md
- docs/technical-verified/integrations-and-events.md
- docs/technical-verified/operations-observability.md
- docs/technical-verified/testing-and-delivery.md

## GRACE Sources
- docs/development-plan.xml (M-API, M-GO-CONFIG, M-GO-LOGGER, M-GRAPHQL-SCHEMA)
- docs/knowledge-graph.xml (full module graph)
- docs/verification-plan.xml (VF-LOCAL-DEV, VF-USER-GRAPHQL, VF-WEB-ADMIN-AUTH)

## Codebase Sources (read-only)
- apps/api/cmd/server/main.go
- apps/api/internal/appconfig/config.go
- apps/api/internal/handler/health.go
- apps/api/internal/service/
- apps/api/internal/repository/
- apps/api/internal/middleware/
- apps/api/internal/graph/
- apps/api/gqlgen.yml
- apps/api/sqlc.yaml
- libs/go/config/
- libs/go/logger/
- libs/graphql/schema/
- docker-compose.yml (root)
- package.json

## Source Delta
- N/A (first detailed wave run)

## Source Gaps
- No fitness-domain GraphQL schema (schema defined for admin only)
- No fitness-domain database migrations
- No PIN auth service (admin auth exists)
- No settings service for fitness app
- No Docker Compose fitness-service overrides