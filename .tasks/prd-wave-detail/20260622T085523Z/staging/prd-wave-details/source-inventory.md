<!-- FILE: docs/prd-wave-details/source-inventory.md -->
<!-- VERSION: 1.0.1 -->

# Source Inventory

## PRD Wave Sources
- docs/prd-waves/index.md
- docs/prd-waves/wave-map.md
- docs/prd-waves/waves/wave-09.md
- docs/prd-waves/frontend-pages/index.md
- docs/prd-waves/frontend-pages/page-010.md

## Frontend Pages Source
- docs/prd-waves/frontend-pages/index.md
- docs/prd-waves/frontend-pages/page-010.md

## Product Sources
- docs/product-verified/domain-model.md
- docs/product-verified/functional-spec.md §20 (REQ-016)
- docs/product-verified/acceptance-criteria.md (AC-093-102, AC-114-116, AC-124)
- docs/product-verified/features/backup-and-restore.md
- docs/product-verified/user-flows.md §26.11-§26.12
- docs/product-verified/edge-cases.md (EDGE-010, EDGE-021, EDGE-028)
- docs/product-verified/business-rules.md (RULE-007, RULE-008, RULE-028)
- docs/product-verified/product-brief.md (performance targets)
- docs/product-verified/actors-and-permissions.md

## Technical Sources
- (none — docs/technical-verified not required for shallow waves)

## GRACE Sources
- docs/development-plan.xml
- docs/knowledge-graph.xml
- docs/verification-plan.xml

## Codebase Sources
- apps/api/internal/atlas/service/ai_export_service.go (ZIP pattern)
- apps/api/internal/atlas/service/export_zip.go (BuildZIP reusable)
- apps/api/internal/atlas/handler/ai_export_handler.go (REST handler pattern)
- apps/api/internal/atlas/handler/atlas_media.go (file download pattern)
- apps/api/internal/atlas/service/ai_review_service.go (ListAllByUserID pattern)
- apps/api/internal/atlas/models/ai_export.go (model pattern)
- apps/api/internal/atlas/graph/resolver/resolver.go (resolver struct)
- apps/api/cmd/server/main.go (wiring pattern)
- apps/api/atlas-gqlgen.yml (gqlgen bindings)
- apps/api/internal/appconfig/config.go (config pattern)
- apps/api/internal/repository/postgres/migrations/ (next: 00094)

## Source Delta
WAVE-09 is the ninth wave. All prior waves (WAVE-01 through WAVE-08) are completed. Source wave is user-approved (2026-06-18).

## Source Gaps
- Q-ACTOR-08, Q-AC-15: Import behavior when data already exists — recorded as DQ-W09-001
- Q-AC-16: CSV files mandatory or optional — recorded as DQ-W09-002
- Q-EDGE-11: Schema version migration strategy — recorded as DQ-W09-005