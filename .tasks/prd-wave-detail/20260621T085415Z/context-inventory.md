# Context Inventory

## Source Documents Read
- docs/prd-waves/waves/wave-06.md (source wave)
- docs/prd-waves/wave-map.md
- docs/prd-waves/index.md
- docs/prd-waves/source-inventory.md
- docs/prd-waves/scope-inventory.md
- docs/prd-waves/open-questions.md
- docs/prd-waves/frontend-pages/index.md
- docs/prd-waves/frontend-pages/page-008.md
- docs/prd-waves/waves/index.md
- docs/prd-waves/appendix/question-ledger.md
- docs/prd-waves/appendix/traceability.md
- docs/product-verified/features/charts.md
- docs/product-verified/functional-spec.md
- docs/product-verified/domain-model.md
- docs/product-verified/business-rules.md
- docs/product-verified/edge-cases.md
- docs/product-verified/acceptance-criteria.md
- docs/product-verified/actors-and-permissions.md

## Prior Detailed Waves Read
- docs/prd-wave-details/waves/wave-05.md (most recent, full pattern reference)
- docs/prd-wave-details/waves/wave-04.md
- docs/prd-wave-details/codebase-fit.md
- docs/prd-wave-details/index.md
- docs/prd-wave-details/waves/index.md

## GRACE Sources Read
- docs/development-plan.xml
- docs/knowledge-graph.xml
- docs/verification-plan.xml

## Codebase Context (Read-Only)
- apps/api/internal/atlas — existing module structure
- apps/api/internal/atlas/service/exercise.go — service pattern
- apps/api/internal/atlas/service/nutrition_macro_service.go — calculation service pattern
- apps/api/internal/atlas/repository/postgres/ — repository pattern
- apps/api/internal/atlas/graph/resolver/ — resolver pattern
- apps/api/internal/atlas/graph/schema/ — schema files
- apps/api/cmd/server/main.go — wiring

## Key Context
- WAVE-06 is a pure backend data wave: aggregation/query-only for chart data
- No mutations — only queries returning time-series data
- Depends on WAVE-03 (workout/sets), WAVE-04 (body weight/measurements), WAVE-05 (nutrition macros)
- Epley formula selected for e1RM calculation: weight × (1 + reps / 30)
- Frontend renders charts — this wave provides the data queries
- No media, no uploads, no new storage