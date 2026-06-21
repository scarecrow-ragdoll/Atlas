# Source Inventory

## PRD Wave Sources
- /Users/vlad/Develop/Atlas/docs/prd-waves/index.md
- /Users/vlad/Develop/Atlas/docs/prd-waves/wave-map.md
- /Users/vlad/Develop/Atlas/docs/prd-waves/source-inventory.md
- /Users/vlad/Develop/Atlas/docs/prd-waves/scope-inventory.md
- /Users/vlad/Develop/Atlas/docs/prd-waves/open-questions.md
- /Users/vlad/Develop/Atlas/docs/prd-waves/waves/index.md
- /Users/vlad/Develop/Atlas/docs/prd-waves/waves/wave-07.md
- /Users/vlad/Develop/Atlas/docs/prd-waves/appendix/question-ledger.md
- /Users/vlad/Develop/Atlas/docs/prd-waves/appendix/decision-log.md
- /Users/vlad/Develop/Atlas/docs/prd-waves/appendix/reviewer-verdicts.md
- /Users/vlad/Develop/Atlas/docs/prd-waves/appendix/run-history.md
- /Users/vlad/Develop/Atlas/docs/prd-waves/appendix/traceability.md

## Frontend Pages Source
- /Users/vlad/Develop/Atlas/docs/prd-waves/frontend-pages/index.md
- /Users/vlad/Develop/Atlas/docs/prd-waves/frontend-pages/page-009.md

## Product Sources
- /Users/vlad/Develop/Atlas/docs/product-verified/index.md
- /Users/vlad/Develop/Atlas/docs/product-verified/product-brief.md
- /Users/vlad/Develop/Atlas/docs/product-verified/scope.md
- /Users/vlad/Develop/Atlas/docs/product-verified/actors-and-permissions.md
- /Users/vlad/Develop/Atlas/docs/product-verified/domain-model.md
- /Users/vlad/Develop/Atlas/docs/product-verified/functional-spec.md
- /Users/vlad/Develop/Atlas/docs/product-verified/user-flows.md
- /Users/vlad/Develop/Atlas/docs/product-verified/business-rules.md
- /Users/vlad/Develop/Atlas/docs/product-verified/edge-cases.md
- /Users/vlad/Develop/Atlas/docs/product-verified/acceptance-criteria.md
- /Users/vlad/Develop/Atlas/docs/product-verified/open-questions.md
- /Users/vlad/Develop/Atlas/docs/product-verified/features/ai-export.md (if exists)

## Technical Sources
- SOURCE_MISSING: /Users/vlad/Develop/Atlas/docs/technical-verified

## GRACE Sources
- /Users/vlad/Develop/Atlas/docs/development-plan.xml
- /Users/vlad/Develop/Atlas/docs/knowledge-graph.xml
- /Users/vlad/Develop/Atlas/docs/verification-plan.xml

## Codebase Sources
- apps/api/internal/atlas/models/ — existing model types
- apps/api/internal/atlas/service/ — existing service pattern
- apps/api/internal/atlas/repository/postgres/ — repository adapters
- apps/api/internal/atlas/graph/resolver/ — resolver container
- apps/api/internal/atlas/graph/schema/ — GraphQL schema files
- apps/api/cmd/server/main.go — wiring pattern
- apps/api/atlas-gqlgen.yml — gqlgen model binding config

## Source Delta
- WAVE-07: initial detailed run. 6 planner reports, 7 reviewer verdicts, 1 final fit review.
- WAVE-07 scope: CAP-W07-003 (week flags CRUD) removed — owned by WAVE-04
- WAVE-07 scope: UserProfile as separate table (not extending atlas_users or Settings)

## Source Gaps
- docs/technical-verified absent — no API contract, auth, or data contract docs from technical verification
- No existing UserProfile service/repo/model — needs full implementation in WAVE-07
- No existing AiExport service/repo/model — needs full implementation in WAVE-07