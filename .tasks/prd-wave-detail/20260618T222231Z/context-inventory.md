# Context Inventory

## Run ID

20260618T222231Z

## Selected Wave

WAVE-05: Nutrition

## PRD Wave Sources

- docs/prd-waves/index.md
- docs/prd-waves/wave-map.md
- docs/prd-waves/open-questions.md
- docs/prd-waves/source-inventory.md
- docs/prd-waves/scope-inventory.md
- docs/prd-waves/waves/index.md
- docs/prd-waves/waves/wave-04.md (Cardio and Body — parallelizable, used for pattern reference)
- docs/prd-waves/waves/wave-05.md (Nutrition — selected wave)
- docs/prd-waves/frontend-pages/index.md
- docs/prd-waves/frontend-pages/page-007.md (Nutrition page — backend dependency context only)
- docs/prd-waves/appendix/decision-log.md
- docs/prd-waves/appendix/question-ledger.md

## Product-verified Sources

- docs/product-verified/index.md
- docs/product-verified/product-brief.md
- docs/product-verified/functional-spec.md (Nutrition §15 REQ-010/REQ-011)
- docs/product-verified/domain-model.md (NutritionProduct, NutritionTemplate, NutritionTemplateItem, DailyNutritionOverride, DailyNutritionOverrideItem)
- docs/product-verified/acceptance-criteria.md (AC-017–AC-019, AC-058–AC-064, AC-113)
- docs/product-verified/edge-cases.md (EDGE-003, EDGE-009, EDGE-017, EDGE-019)
- docs/product-verified/business-rules.md (RULE-006, RULE-010, RULE-011, RULE-018, RULE-019, RULE-020)
- docs/product-verified/user-flows.md (§26.7 Create Nutrition Template, §26.8 Override Daily Nutrition)
- docs/product-verified/actors-and-permissions.md
- docs/product-verified/scope.md
- docs/product-verified/features/nutrition.md

## Technical Sources

- docs/technical-verified/api-contracts.md (hybrid GraphQL/REST, TDEC-027 error format)
- docs/technical-verified/auth-security-compliance.md (PIN auth, TDEC-037)
- docs/technical-verified/data-contracts.md (domain entities, userId FKs)
- docs/technical-verified/operations-observability.md (log markers, error format)
- docs/technical-verified/implementation-slices.md
- docs/technical-verified/testing-and-delivery.md

## Prior Detailed Waves

- docs/prd-wave-details/index.md (package status: questions-open for WAVE-04)
- docs/prd-wave-details/waves/wave-01.md (Foundation — ready-for-dev, pattern for detailed wave)
- docs/prd-wave-details/waves/wave-02.md (Exercise Library — user-approved, pattern for detailed wave)
- docs/prd-wave-details/waves/wave-04.md (Cardio and Body — questions-open, pattern for detailed wave)
- docs/prd-wave-details/wave-map-context.md
- docs/prd-wave-details/codebase-fit.md
- docs/prd-wave-details/source-inventory.md
- docs/prd-wave-details/open-questions.md
- docs/prd-wave-details/appendix/reviewer-verdicts.md
- docs/prd-wave-details/appendix/question-ledger.md
- docs/prd-wave-details/appendix/decision-log.md
- docs/prd-wave-details/appendix/traceability.md
- docs/prd-wave-details/appendix/run-history.md

## GRACE Sources

- docs/development-plan.xml
- docs/knowledge-graph.xml
- docs/verification-plan.xml

## Codebase Sources (Read-only)

See codebase exploration in task report. Key patterns:

- Migration pattern: apps/api/internal/repository/postgres/migrations/00080_atlas_foundation.sql (last migration: 00080)
- sqlc query pattern: apps/api/internal/repository/postgres/queries/atlas_settings.sql
- Repository adapter pattern: apps/api/internal/atlas/repository/postgres/settings_repo.go
- Service pattern: apps/api/internal/atlas/service/settings_service.go
- Models pattern: apps/api/internal/atlas/models/settings.go
- GraphQL schema pattern: apps/api/internal/atlas/graph/schema/settings.graphql
- Resolver pattern: apps/api/internal/atlas/graph/resolver/settings.go
- Resolver container: apps/api/internal/atlas/graph/resolver/resolver.go
- Main wiring: apps/api/cmd/server/main.go (atlasRes section at lines 161-164)
- gqlgen config: apps/api/atlas-gqlgen.yml
- sqlc config: apps/api/sqlc.yaml

## Source Delta

WAVE-04 was detailed in prior runs. No source changes affect WAVE-05 since wave map approval. This is the first detail run for WAVE-05.

## Source Gaps

- No nutrition-related sqlc queries, repository adapters, services, GraphQL schema, or resolvers exist yet
- WAVE-01 PIN auth middleware and Atlas GraphQL endpoint not yet implemented (blocking dependency)
- Nutrition macro calculation engine does not exist yet
- No test infrastructure for nutrition domain exists yet