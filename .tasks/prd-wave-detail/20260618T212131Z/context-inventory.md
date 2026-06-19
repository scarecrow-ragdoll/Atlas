# Context Inventory

## Selected Wave
WAVE-03: Workout Diary

## PRD Wave Sources
- docs/prd-waves/index.md — waves-approved, source gate: docs/product/prd.md + docs/product-verified/
- docs/prd-waves/wave-map.md — WAVE-03: Workout Diary - Daily workouts with sets
- docs/prd-waves/waves/index.md — WAVE-03 status: user-approved
- docs/prd-waves/waves/wave-03.md — user-approved, purpose: "Workout day by date with exercises and sets"
- docs/prd-waves/open-questions.md — Q-WORKOUT-001: concurrent edit handling
- docs/prd-waves/source-inventory.md — source gap: Q-WORKOUT-001

## Frontend Pages Source
- docs/prd-waves/frontend-pages/index.md — PAGE-002: Workout Diary
- docs/prd-waves/frontend-pages/page-002.md — Daily workout entry by date, exercise cards, sets table, cardio section
- docs/prd-waves/frontend-pages/page-004.md — Cardio entries (WAVE-04 dependency)

## Product Sources
- docs/product-verified/index.md
- docs/product-verified/functional-spec.md
- docs/product-verified/domain-model.md (DailyLog, WorkoutExercise, WorkoutSet, CardioEntry entities)
- docs/product-verified/acceptance-criteria.md
- docs/product-verified/user-flows.md (Log Workout Today)
- docs/product-verified/business-rules.md
- docs/product-verified/edge-cases.md

## Technical Sources
- docs/technical-verified/index.md
- docs/technical-verified/architecture-and-boundaries.md
- docs/technical-verified/api-contracts.md (hybrid GraphQL/REST, TDEC-001, TDEC-026, TDEC-027)
- docs/technical-verified/data-contracts.md (DailyLog, WorkoutExercise, WorkoutSet, CardioEntry — TDEC-020, TDEC-021, TDEC-022)
- docs/technical-verified/auth-security-compliance.md (PIN auth per TDEC-033..040, TDEC-037)
- docs/technical-verified/implementation-slices.md (Slice 2: DailyLog + Cardio + WorkoutExercise + WorkoutSet)
- docs/technical-verified/operations-observability.md
- docs/technical-verified/testing-and-delivery.md

## GRACE Sources
- docs/development-plan.xml
- docs/knowledge-graph.xml
- docs/verification-plan.xml

## Prior Detailed Waves
- docs/prd-wave-details/waves/wave-01.md (Foundation)
- docs/prd-wave-details/waves/wave-02.md (Exercise Library)
- docs/prd-wave-details/index.md — status: wave-approved (WAVE-02)
- docs/prd-wave-details/wave-map-context.md
- docs/prd-wave-details/codebase-fit.md
- docs/prd-wave-details/source-inventory.md
- docs/prd-wave-details/open-questions.md

## Codebase Sources
- apps/api/cmd/server/main.go — API wiring pattern
- apps/api/internal/appconfig/config.go — Config extension pattern
- apps/api/internal/middleware/admin_auth.go — Auth middleware pattern
- apps/api/internal/service/admin_auth.go — Service layer pattern
- apps/api/internal/repository/postgres/user_repo.go — Repository adapter pattern
- apps/api/internal/repository/redis/admin_session_store.go — Redis store pattern
- apps/api/gqlgen.yml — Codegen config
- apps/api/sqlc.yaml — SQLC config
- libs/graphql/schema/schema.graphql — Schema extension pattern

## Source Delta
- Initial detailed planning for WAVE-03
- WAVE-01: ready-for-dev awaiting user approval
- WAVE-02: user-approved
- No source delta from previous runs for WAVE-03

## Source Gaps
- Concurrent edit handling strategy undefined (Q-WORKOUT-001)
- No Workout Diary-related GraphQL schema, sqlc queries, repository, service, or handler code exists yet
- WAVE-01 PIN auth middleware not yet implemented (blocking dependency)
- WAVE-02 allExercises query not yet implemented (blocking dependency)