# Context Inventory

## PRD Wave Sources
- docs/prd-waves/index.md
- docs/prd-waves/wave-map.md
- docs/prd-waves/open-questions.md
- docs/prd-waves/waves/index.md
- docs/prd-waves/waves/wave-02.md

## Frontend Pages Source
- docs/prd-waves/frontend-pages/index.md
- docs/prd-waves/frontend-pages/page-002.md (Workout Diary - depends on exercises)
- docs/prd-waves/frontend-pages/page-003.md (Exercise Library - primary frontend consumer)

## Product Sources
- docs/product-verified/product-brief.md
- docs/product-verified/functional-spec.md (Section: Exercise Library §11 — REQ-003)
- docs/product-verified/domain-model.md (Exercise, ExerciseMedia entities)
- docs/product-verified/actors-and-permissions.md
- docs/product-verified/business-rules.md
- docs/product-verified/edge-cases.md
- docs/product-verified/scope.md
- docs/product-verified/user-flows.md
- docs/product-verified/acceptance-criteria.md (AC-002, AC-003, AC-004, AC-043, AC-044, AC-045, AC-046, AC-047)

## Technical Sources
- docs/technical-verified/technical-brief.md
- docs/technical-verified/architecture-and-boundaries.md
- docs/technical-verified/api-contracts.md
- docs/technical-verified/data-contracts.md (Exercise, ExerciseMedia entities)
- docs/technical-verified/auth-security-compliance.md
- docs/technical-verified/integrations-and-events.md
- docs/technical-verified/operations-observability.md
- docs/technical-verified/implementation-slices.md (Slice 1: Exercise Library)
- docs/technical-verified/testing-and-delivery.md
- docs/technical-verified/client-state-and-ux-contracts.md

## GRACE Sources
- docs/development-plan.xml
- docs/knowledge-graph.xml
- docs/verification-plan.xml

## Prior Detailed Waves
- docs/prd-wave-details/waves/wave-01.md (Foundation - ready-for-dev)

## Codebase Sources
- apps/api/cmd/server/main.go (API wiring pattern)
- apps/api/internal/appconfig/config.go (Config extension pattern)
- apps/api/internal/middleware/admin_auth.go (Auth middleware pattern)
- apps/api/internal/service/admin_auth.go (Service pattern)
- apps/api/internal/repository/postgres/user_repo.go (Repository pattern)
- apps/api/internal/repository/redis/admin_session_store.go (Redis store pattern)
- apps/api/internal/handler/health.go (Handler pattern)
- apps/api/gqlgen.yml (Codegen config)
- apps/api/sqlc.yaml (SQLC config)
- libs/graphql/schema/schema.graphql (Schema extension pattern)
- libs/graphql/schema/admin_auth.graphql (Auth schema extension)

## Key Observations
- API uses hybrid GraphQL (CRUD) + REST (binary uploads) pattern per TDEC-001
- WAVE-01 provides: media REST scaffold (POST/GET /api/v1/media), fitness GraphQL schema extension pattern, PIN auth middleware, db migration infrastructure
- Exercise entity: id, name, muscleGroups, description, personalNotes, workingWeight, isActive, createdAt, updatedAt
- ExerciseMedia entity: id, exerciseId, mediaType, filePath, originalFileName, mimeType, sizeBytes, createdAt
- Existing media REST scaffold from WAVE-01 handles upload/download; WAVE-02 adds exercise-media association
- PAGE-003 (Exercise Library) lists backend dependencies: GET/POST/PUT/DELETE /api/exercises, POST/DELETE /api/exercise-media
- PAGE-002 (Workout Diary) depends on GET /api/exercises for exercise library integration
- No existing exercise-related SQL queries, GraphQL schema, repository, service, or resolver code