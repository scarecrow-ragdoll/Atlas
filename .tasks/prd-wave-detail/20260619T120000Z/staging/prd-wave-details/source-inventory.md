# Source Inventory
## PRD Wave Sources
- docs/prd-waves/index.md
- docs/prd-waves/wave-map.md
- docs/prd-waves/open-questions.md
- docs/prd-waves/waves/index.md
- docs/prd-waves/waves/wave-01.md (Foundation)
- docs/prd-waves/waves/wave-02.md (Exercise Library)
- docs/prd-waves/waves/wave-04.md (Cardio and Body Tracking)
## Frontend Pages Source
- docs/prd-waves/frontend-pages/index.md
- docs/prd-waves/frontend-pages/page-001.md (Dashboard — depends on latestBodyWeight)
- docs/prd-waves/frontend-pages/page-004.md (Cardio — primary frontend consumer for cardio CRUD)
- docs/prd-waves/frontend-pages/page-005.md (Body Measurements — primary frontend consumer for check-in/measurements/weight)
- docs/prd-waves/frontend-pages/page-006.md (Progress Photos — photo management within check-ins)
## Product Sources
- docs/product-verified/index.md
- docs/product-verified/product-brief.md
- docs/product-verified/functional-spec.md (Cardio §12 REQ-007, Body Tracking §13 REQ-008/REQ-009)
- docs/product-verified/domain-model.md (CardioEntry, BodyWeightEntry, BodyCheckIn, BodyMeasurement, ProgressPhoto, WeekFlag)
- docs/product-verified/acceptance-criteria.md (AC-012–AC-016, AC-048–AC-057)
- docs/product-verified/edge-cases.md (EDGE-006, EDGE-007)
- docs/product-verified/business-rules.md (RULE-005)
- docs/product-verified/user-flows.md (§26.4 Add Cardio, §26.5 Weekly Check-In, §26.6 Add Body Weight)
- docs/product-verified/actors-and-permissions.md
- docs/product-verified/scope.md
## Technical Sources
- docs/technical-verified/api-contracts.md (hybrid GraphQL/REST, TDEC-027)
- docs/technical-verified/auth-security-compliance.md (PIN auth, TDEC-037)
- docs/technical-verified/operations-observability.md (log markers, error format)
- docs/technical-verified/data-contracts.md (domain entities)
## GRACE Sources
- docs/development-plan.xml
- docs/knowledge-graph.xml
- docs/verification-plan.xml
## Codebase Sources
- apps/api/cmd/server/main.go (API wiring pattern, route groups)
- apps/api/internal/appconfig/config.go (Config extension pattern)
- apps/api/internal/middleware/admin_auth.go (Auth middleware pattern)
- apps/api/internal/service/exercise.go (Service layer pattern)
- apps/api/internal/repository/postgres/user_repo.go (Repository adapter pattern)
- apps/api/internal/handler/exercise_media.go (REST handler pattern)
- apps/api/internal/graph/exercise.resolvers.go (GraphQL resolver pattern)
- apps/api/gqlgen.yml (Codegen config)
- apps/api/sqlc.yaml (SQLC config)
- libs/graphql/schema/exercises.graphql (Schema extension pattern)
- libs/graphql/schema/schema.graphql (Schema extension with extend type pattern)
## Source Delta
- Added WAVE-04 (Cardio and Body Tracking) detail to the existing WAVE-01/WAVE-02 details
- WAVE-01 is ready-for-dev awaiting user approval — no source delta between runs
- WAVE-02 is user-approved — no source delta between runs
## Source Gaps
- No cardio/body-related GraphQL schema, sqlc queries, repository, service, or handler code exists yet
- WAVE-01 PIN auth middleware not yet implemented (blocking dependency)
- WAVE-01 media REST scaffold not yet implemented (blocking dependency for progress photos)
- No server-side MIME detection utility exists
- DailyLog table migration strategy pending WAVE-03 coordination (DQ-W04-001)