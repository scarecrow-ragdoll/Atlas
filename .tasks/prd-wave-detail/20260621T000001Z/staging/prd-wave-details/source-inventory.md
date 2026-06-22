<!-- FILE: docs/prd-wave-details/source-inventory.md -->
<!-- VERSION: 1.0.0 -->

# Source Inventory

## PRD Wave Sources
- docs/prd-waves/index.md
- docs/prd-waves/wave-map.md
- docs/prd-waves/waves/wave-08.md
- docs/prd-waves/frontend-pages/index.md
- docs/prd-waves/frontend-pages/page-009.md

## Frontend Pages Source
- docs/prd-waves/frontend-pages/index.md
- docs/prd-waves/frontend-pages/page-009.md

## Product Sources
- docs/product-verified/domain-model.md
- docs/product-verified/functional-spec.md
- docs/product-verified/acceptance-criteria.md
- docs/product-verified/features/ai-review-history.md
- docs/product-verified/edge-cases.md
- docs/product-verified/business-rules.md

## Technical Sources
- docs/technical-verified/data-contracts.md
- docs/technical-verified/architecture-and-boundaries.md
- docs/technical-verified/testing-and-delivery.md

## GRACE Sources
- docs/development-plan.xml
- docs/knowledge-graph.xml
- docs/verification-plan.xml

## Codebase Sources
- apps/api/internal/repository/postgres/migrations/00092_ai_exports.sql
- apps/api/internal/repository/postgres/queries/
- apps/api/internal/atlas/models/week_flag.go
- apps/api/internal/atlas/models/ai_export.go
- apps/api/internal/atlas/repository/postgres/week_flag_repo.go
- apps/api/internal/atlas/service/week_flag.go
- apps/api/internal/atlas/graph/schema/week_flag.graphql
- apps/api/internal/atlas/graph/resolver/week_flag.go
- apps/api/internal/atlas/graph/resolver/resolver.go
- apps/api/cmd/server/main.go
- apps/api/atlas-gqlgen.yml

## Source Delta
No source deltas since WAVE-08 source wave was approved (2026-06-18). All prior questions resolved.

## Source Gaps
- planned_actions storage format (TEXT vs structured) not specified in source docs — recorded as DQ-W08-001
- WAVE-09 backup contract for AiReview not explicitly defined in source docs — recorded as DQ-W08-002