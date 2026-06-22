# Context Inventory

## Selected Backend Wave
WAVE-08: AI Review History

## Source Docs
- docs/prd-waves/waves/wave-08.md (source wave)
- docs/prd-waves/index.md (waves-approved)
- docs/prd-waves/wave-map.md (WAVE-08: AI Review History)
- docs/prd-waves/frontend-pages/index.md (PAGE-009 AI Export as consumer)
- docs/prd-waves/frontend-pages/page-009.md (backend deps: AiReview backend)

## Verified Docs
- docs/product-verified/domain-model.md (AiReview entity: id, dateRangeStart, dateRangeEnd, aiResponseText, userNotes, plannedActions, createdAt, updatedAt)
- docs/product-verified/functional-spec.md §19 (AI Review: manual entry, linked to date range, notes, planned actions, review history)
- docs/product-verified/acceptance-criteria.md (AC-025, AC-090, AC-091, AC-092)
- docs/product-verified/features/ai-review-history.md (manual entry, paste AI text, link to date range, notes, planned actions, viewable history)
- docs/product-verified/business-rules.md
- docs/product-verified/edge-cases.md

## Technical Docs
- docs/technical-verified/index.md
- docs/technical-verified/data-contracts.md
- docs/technical-verified/architecture-and-boundaries.md
- docs/technical-verified/implementation-slices.md
- docs/technical-verified/testing-and-delivery.md

## GRACE Docs
- docs/development-plan.xml
- docs/knowledge-graph.xml
- docs/verification-plan.xml

## Prior Detailed Waves
- WAVE-01 (docs/prd-wave-details/waves/wave-01.md): Foundation — PIN auth, Settings, migrations
- WAVE-02 (docs/prd-wave-details/waves/wave-02.md): Exercise Library
- WAVE-04 (docs/prd-wave-details/waves/wave-04.md): Cardio and Body Tracking
- WAVE-05 (docs/prd-wave-details/waves/wave-05.md): Nutrition
- WAVE-06 (docs/prd-wave-details/waves/wave-06.md): Charts
- WAVE-07 (docs/prd-wave-details/waves/wave-07.md): AI Export and Prompt Builder (direct neighbor — AiExport entity, UserProfile, export infrastructure)

## Codebase Sources
- apps/api/internal/repository/postgres/migrations/ (latest: 00092_ai_exports.sql)
- apps/api/internal/atlas/models/ (existing model patterns)
- apps/api/internal/atlas/repository/postgres/ (existing repo patterns)
- apps/api/internal/atlas/service/ (existing service patterns)
- apps/api/internal/atlas/graph/ (existing GraphQL resolver/schema patterns)
- apps/api/internal/handler/ (existing REST handler patterns)
- apps/api/cmd/server/main.go (wiring)
- apps/api/atlas-gqlgen.yml (gqlgen config)
- apps/api/sqlc.yaml (sqlc config)

## Neighboring Backend Waves
- WAVE-07 (prior): WAVE-08 depends on WAVE-07. WAVE-07 creates AiExport record. WAVE-08 creates independent AiReview record. Clean boundary per wave-07.md.
- WAVE-09 (future): Backup — WAVE-09 includes AiReview data in backup. WAVE-08 must provide service layer for WAVE-09 consumption.

## Frontend Pages Context
- PAGE-009 (AI Export): may reference AiReview history. No dedicated frontend page for AiReview in the page list — likely exposed through PAGE-009 or a sub-section.
