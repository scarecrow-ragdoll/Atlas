# Source Inventory

## PRD Wave Sources
- docs/prd-waves/index.md
- docs/prd-waves/wave-map.md
- docs/prd-waves/open-questions.md
- docs/prd-waves/waves/index.md
- docs/prd-waves/waves/wave-01.md (Foundation)
- docs/prd-waves/waves/wave-02.md

## Frontend Pages Source
- docs/prd-waves/frontend-pages/index.md
- docs/prd-waves/frontend-pages/page-002.md (Workout Diary — depends on WAVE-02 allExercises)
- docs/prd-waves/frontend-pages/page-003.md (Exercise Library — primary frontend consumer)
- docs/prd-waves/frontend-pages/page-011.md (Settings — PIN auth dependency from WAVE-01)

## Product Sources
- docs/product-verified/index.md
- docs/product-verified/product-brief.md
- docs/product-verified/functional-spec.md (Exercise Library §11 — REQ-003)
- docs/product-verified/domain-model.md (Exercise, ExerciseMedia entities)
- docs/product-verified/actors-and-permissions.md
- docs/product-verified/business-rules.md
- docs/product-verified/edge-cases.md
- docs/product-verified/scope.md
- docs/product-verified/user-flows.md
- docs/product-verified/acceptance-criteria.md (AC-002, AC-003, AC-004, AC-043, AC-044, AC-045, AC-046, AC-047)

## Technical Sources
- docs/technical-verified/index.md
- docs/technical-verified/technical-brief.md
- docs/technical-verified/architecture-and-boundaries.md
- docs/technical-verified/api-contracts.md (hybrid GraphQL/REST, TDEC-001, TDEC-027, TDEC-008)
- docs/technical-verified/data-contracts.md (Exercise, ExerciseMedia, TDEC-020, TDEC-022, TDEC-023)
- docs/technical-verified/auth-security-compliance.md (PIN auth, TDEC-037)
- docs/technical-verified/integrations-and-events.md
- docs/technical-verified/operations-observability.md
- docs/technical-verified/implementation-slices.md (Slice 1: Exercise Library)
- docs/technical-verified/testing-and-delivery.md
- docs/technical-verified/client-state-and-ux-contracts.md

## GRACE Sources
- docs/development-plan.xml
- docs/knowledge-graph.xml
- docs/verification-plan.xml

## Codebase Sources
- apps/api/cmd/server/main.go (API wiring pattern, route groups)
- apps/api/internal/appconfig/config.go (Config extension pattern)
- apps/api/internal/middleware/admin_auth.go (Auth middleware pattern)
- apps/api/internal/service/admin_auth.go (Service layer pattern)
- apps/api/internal/repository/postgres/user_repo.go (Repository adapter pattern)
- apps/api/internal/repository/redis/admin_session_store.go (Redis store pattern)
- apps/api/internal/handler/health.go (REST handler pattern)
- apps/api/gqlgen.yml (Codegen config)
- apps/api/sqlc.yaml (SQLC config)
- libs/graphql/schema/schema.graphql (Schema extension with extend type pattern)
- libs/graphql/schema/admin_auth.graphql (Auth schema extension pattern)
- apps/api/internal/repository/postgres/migrations/00079_admin_users.sql (Migration pattern)

## Source Delta
- Added WAVE-02 (Exercise Library) detail to the existing WAVE-01 detail
- WAVE-01 is ready-for-dev awaiting user approval — no source delta between runs

## Source Gaps
- No exercise-related GraphQL schema, sqlc queries, repository, service, or handler code exists yet
- WAVE-01 PIN auth middleware not yet implemented (blocking dependency)
- WAVE-01 media REST scaffold not yet implemented (blocking dependency)
- No server-side MIME detection utility exists