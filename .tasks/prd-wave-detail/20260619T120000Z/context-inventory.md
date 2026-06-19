# Context Inventory

## Run ID
20260619T120000Z

## PRD Wave Sources
- docs/prd-waves/index.md
- docs/prd-waves/wave-map.md
- docs/prd-waves/open-questions.md
- docs/prd-waves/waves/index.md
- docs/prd-waves/waves/wave-04.md (Cardio and Body Tracking)
- docs/prd-waves/source-inventory.md
- docs/prd-waves/scope-inventory.md

## Frontend Pages Source
- docs/prd-waves/frontend-pages/index.md
- docs/prd-waves/frontend-pages/page-004.md (Cardio — primary consumer of cardio CRUD)
- docs/prd-waves/frontend-pages/page-005.md (Body Measurements — primary consumer of check-in/measurements/weight)
- docs/prd-waves/frontend-pages/page-006.md (Progress Photos — photo management within check-ins)
- docs/prd-waves/frontend-pages/page-001.md (Dashboard — cardio count, last weight)

## Product Sources
- docs/product-verified/index.md
- docs/product-verified/product-brief.md
- docs/product-verified/domain-model.md (CardioEntry §1.17, BodyWeightEntry §1.18, BodyCheckIn §1.19, BodyMeasurement §1.20, ProgressPhoto §1.21, WeekFlag §1.27)
- docs/product-verified/functional-spec.md (Cardio §12, Body Tracking §13 — REQ-007, REQ-008, REQ-009)
- docs/product-verified/acceptance-criteria.md (AC-012–AC-016 for cardio/body; AC-048–AC-057)
- docs/product-verified/edge-cases.md (EDGE-006, EDGE-007)
- docs/product-verified/business-rules.md (RULE-005 photo count)
- docs/product-verified/user-flows.md (Add Cardio §26.4, Weekly Check-In §26.5, Add Body Weight Standalone §26.6)
- docs/product-verified/actors-and-permissions.md
- docs/product-verified/scope.md
- docs/product-verified/open-questions.md

## Technical Sources
- docs/technical-verified (not present — SOURCE_MISSING per source-inventory.md)

## GRACE Sources
- docs/development-plan.xml (M-API module, M-PRD-WAVE-DETAILER)
- docs/knowledge-graph.xml (M-API, M-PRD-WAVE-DETAILER, M-GRACE-WORKFLOW)
- docs/verification-plan.xml (V-M-API for API tests)

## Prior Detailed Waves
- docs/prd-wave-details/index.md
- docs/prd-wave-details/wave-map-context.md
- docs/prd-wave-details/codebase-fit.md
- docs/prd-wave-details/source-inventory.md
- docs/prd-wave-details/open-questions.md
- docs/prd-wave-details/waves/wave-01.md (Foundation — ready-for-dev)
- docs/prd-wave-details/waves/wave-02.md (Exercise Library — user-approved)
- docs/prd-wave-details/appendix/question-ledger.md
- docs/prd-wave-details/appendix/reviewer-verdicts.md
- docs/prd-wave-details/appendix/traceability.md
- docs/prd-wave-details/appendix/decision-log.md
- docs/prd-wave-details/appendix/run-history.md

## Codebase Patterns (from WAVE-02 detail)
- apps/api/cmd/server/main.go — wiring pattern for repos, services, handlers, resolvers
- apps/api/internal/appconfig/config.go — config struct extension pattern
- apps/api/internal/middleware/admin_auth.go — auth middleware pattern
- apps/api/internal/service/admin_auth.go — service layer pattern
- apps/api/internal/repository/postgres/user_repo.go — repository adapter with sqlc
- apps/api/internal/repository/redis/admin_session_store.go — HMAC key derivation pattern
- apps/api/gqlgen.yml — schema glob for auto-discovery
- apps/api/sqlc.yaml — query glob for auto-discovery
- libs/graphql/schema/schema.graphql — extend type Query/Mutation pattern

## WAVE-04 Codebase Touchpoints (planned)
- apps/api/internal/repository/postgres/migrations/: new migrations for cardio_entry, body_weight_entry, body_check_in, body_measurement, progress_photo, week_flag tables
- apps/api/internal/repository/postgres/queries/: new sqlc query definitions
- apps/api/internal/repository/postgres/: new repository adapters
- apps/api/internal/service/: new transport-neutral services
- apps/api/internal/handler/: new REST handlers for media (progress photos use WAVE-01 media scaffold pattern)
- apps/api/internal/graph/: new GraphQL resolvers for cardio/body CRUD
- libs/graphql/schema/: new schema files for cardio, body, check-in types
- apps/api/cmd/server/main.go: wire new services/resolvers

## Source Gaps
- No cardio/body-related GraphQL schema, sqlc queries, repository, service, or handler code exists yet
- WAVE-01 PIN auth middleware not yet implemented (blocking dependency — same as WAVE-02)
- WAVE-01 media REST scaffold pattern not yet implemented (blocking for progress photos — reusable from WAVE-01/WAVE-02 pattern)
- No Server-Side MIME detection utility exists