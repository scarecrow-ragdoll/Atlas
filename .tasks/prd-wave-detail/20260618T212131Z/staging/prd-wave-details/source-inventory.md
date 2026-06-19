# Source Inventory

## PRD Wave Sources
- docs/prd-waves/index.md
- docs/prd-waves/wave-map.md
- docs/prd-waves/open-questions.md
- docs/prd-waves/waves/index.md
- docs/prd-waves/waves/wave-01.md
- docs/prd-waves/waves/wave-02.md
- docs/prd-waves/waves/wave-03.md
- docs/prd-waves/waves/wave-04.md
- docs/prd-waves/waves/wave-05.md

## Frontend Pages Source
- docs/prd-waves/frontend-pages/index.md
- docs/prd-waves/frontend-pages/page-002.md

## Product Sources
- docs/product-verified/index.md
- docs/product-verified/domain-model.md (DailyLog, WorkoutExercise, WorkoutSet, CardioEntry entities)
- docs/product-verified/functional-spec.md (Workout Diary §10 — REQ-004)
- docs/product-verified/acceptance-criteria.md (AC-005 through AC-011, AC-035 through AC-042)
- docs/product-verified/user-flows.md (Enter Workout For Today, Enter Workout Backdated)
- docs/product-verified/business-rules.md (RULE-004, RULE-016, RULE-017)
- docs/product-verified/edge-cases.md (EDGE-001, EDGE-004, EDGE-005, EDGE-016)

## Technical Sources
- docs/technical-verified/index.md
- docs/technical-verified/architecture-and-boundaries.md
- docs/technical-verified/api-contracts.md (hybrid GraphQL/REST, TDEC-001, TDEC-026, TDEC-027)
- docs/technical-verified/data-contracts.md (DailyLog, WorkoutExercise, WorkoutSet, CardioEntry — TDEC-020, TDEC-021)
- docs/technical-verified/auth-security-compliance.md (PIN auth, TDEC-037)
- docs/technical-verified/implementation-slices.md (Slice 2: DailyLog mapping)
- docs/technical-verified/operations-observability.md (log markers, error format)
- docs/technical-verified/testing-and-delivery.md (TDEC-056, TDEC-057)

## GRACE Sources
- docs/development-plan.xml
- docs/knowledge-graph.xml
- docs/verification-plan.xml

## Codebase Sources
- apps/api/cmd/server/main.go — API wiring pattern, route groups
- apps/api/internal/appconfig/config.go — Config extension pattern
- apps/api/internal/middleware/admin_auth.go — Auth middleware pattern
- apps/api/internal/service/admin_auth.go — Service layer pattern
- apps/api/internal/repository/postgres/user_repo.go — Repository adapter with sqlc pattern
- apps/api/internal/repository/redis/admin_session_store.go — Redis store pattern
- apps/api/gqlgen.yml — Schema glob pattern for auto-discovery
- apps/api/sqlc.yaml — Query glob pattern for auto-discovery
- libs/graphql/schema/schema.graphql — Extend type Query/Mutation pattern
- libs/graphql/schema/admin_auth.graphql — Union result type pattern
- libs/graphql/schema/common.graphql — Common types (ValidationError, AuthError, NotFoundError)

## Source Delta
- WAVE-01: ready-for-dev awaiting user approval
- WAVE-02: user-approved
- WAVE-03: detailed planning complete, ready-for-dev

## Source Gaps
- Concurrent edit handling strategy undefined (Q-WORKOUT-001)
- No Workout Diary-related GraphQL schema, sqlc queries, repository, service, or handler code exists yet
- WAVE-01 PIN auth middleware not yet implemented (blocking dependency)
- WAVE-01 common GraphQL types not yet implemented (blocking dependency)
- WAVE-02 allExercises query not yet implemented (blocking dependency)
