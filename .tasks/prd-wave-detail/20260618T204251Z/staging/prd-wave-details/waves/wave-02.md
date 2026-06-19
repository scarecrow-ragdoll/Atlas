# Wave 02: Exercise Library

## Status
ready-for-dev

## User Approval
Not approved yet. Awaiting user review.

## Source Wave Summary
WAVE-02 from docs/prd-waves/waves/wave-02.md. Full CRUD for exercises with working weight and media management. Source status: user-approved (2026-06-18).

## Outcome After Implementation
- OUT-W02-001: Exercises can be created, listed, edited, deleted via GraphQL
- OUT-W02-002: Media (images/video) can be attached and retrieved per exercise
- OUT-W02-003: Working weight stored per exercise, snapshot-ready for WAVE-03
- OUT-W02-004: API ready for workout diary (allExercises query, exercise by ID)

## Scope Included
- CAP-W02-001: Exercise CRUD via GraphQL (create, read, list with pagination, update, soft delete)
- CAP-W02-002: ExerciseMedia upload, download, and deletion via REST
- CAP-W02-003: Working weight field (Float, validated > 0)
- CAP-W02-004: Muscle groups (string array), description, personal notes fields
- CAP-W02-005: isActive flag for soft delete (default true, exclude from default list)

## Scope Excluded
- Workout diary integration (WAVE-03 — uses allExercises query interface)
- AI export specifics (WAVE-07)
- Exercise builder catalog (future scope)
- Exercise sharing (future scope)
- Full-text search (deferred — basic ILIKE if needed)
- Frontend pages, UI, UX, routes, navigation, components

## Dependencies And Other-Wave Fit
- WAVE-01 (Foundation): prerequisite — provides PIN auth middleware, media REST scaffold, migration infrastructure, fitness GraphQL common types, codegen config, config extension pattern. WAVE-02 cannot start until WAVE-01 provides these contracts.
- WAVE-03 (Workout Diary): WAVE-02 provides allExercises query for exercise selector, exercise by ID for working weight snapshot
- WAVE-04/05: No direct dependency
- WAVE-06 (Charts): WAVE-02 provides exercise metadata (name, muscleGroups) for chart labels; WAVE-03 provides historical working weight data
- WAVE-07/08 (AI Export): WAVE-02 exercise + media data consumed via service layer
- WAVE-09 (Backup): WAVE-02 tables are designed for JSON-serializable export compatibility

## Frontend Pages Dependencies
- PAGE-003 (Exercise Library): primary frontend consumer — depends on all WAVE-02 backend endpoints listed below. Dependency context only; no frontend pages, UI, or UX work in this wave.
- PAGE-002 (Workout Diary): depends on allExercises GraphQL query for exercise selector in workout entry.
- Frontend backend dependencies from PAGE-003: GET/POST/PUT/DELETE /api/exercises (via GraphQL mutations/queries), POST/DELETE /api/exercise-media (REST for binary uploads).

## Codebase Fit And Touchpoints
- apps/api/internal/repository/postgres/migrations/00080_exercises.sql: new migration for exercises table
- apps/api/internal/repository/postgres/migrations/00081_exercise_media.sql: new migration for exercise_media table
- apps/api/internal/repository/postgres/queries/exercises.sql: sqlc query definitions
- apps/api/internal/repository/postgres/exercise_repo.go: repository adapter for exercise CRUD
- apps/api/internal/service/exercise.go: transport-neutral exercise service with validation
- apps/api/internal/handler/exercise_media.go: REST handler for media upload/download/delete
- apps/api/internal/graph/exercise.resolvers.go: GraphQL resolvers for exercise queries/mutations
- libs/graphql/schema/exercises.graphql: exercise GraphQL types and operations
- apps/api/cmd/server/main.go: wire ExerciseRepo, ExerciseService, ExerciseMediaHandler, resolver dependency, and PIN-protected route group
- apps/api/internal/appconfig/config.go: reads MediaConfig from WAVE-01 (BasePath, MaxUploadSize)
- apps/api/gqlgen.yml: auto-discovers new schema files via glob
- apps/api/sqlc.yaml: auto-discovers new queries via glob

## Design Contracts
- Soft delete: isActive=false preserves referential integrity for WAVE-03 workout history (DDEC-W02-001)
- Physical file deletion: ExerciseMedia delete removes file from disk; soft-fail (log error, return 204) on deletion failure (DDEC-W02-002)
- Duplicate names: allowed, no unique constraint (DDEC-W02-003, EDGE-002)
- allExercises: GraphQL-only query returning active exercises ordered by name ASC (DDEC-W02-004)
- MIME detection: http.DetectContentType() server-side (512 bytes) as primary validation (DQ-W02-005 resolved)
- PIN auth: WAVE-01 middleware guards all WAVE-02 GraphQL and REST endpoints
- File storage: <WAVE-01 BasePath>/exercise/<exercise_id>/<uuid>.<ext>
- Working weight: REAL type, validated > 0 when provided

## Data API Integration And Operations
- PostgreSQL: goose migrations 00080_exercises.sql and 00081_exercise_media.sql
- Indexes: idx_exercises_is_active (is_active), idx_exercises_name (name), idx_exercise_media_exercise (exercise_id)
- FK constraint: exercise_media.exercise_id → exercises(id) ON DELETE NO ACTION (prevents accidental cascade)
- File validation: server-side MIME detection (JPEG/PNG/WEBP/MP4/MOV/WEBM), per-type size limits (25MB images, 250MB video, 300MB max single upload)
- Memory-safe upload: r.ParseMultipartForm(maxBytes) with 300MB limit per TDEC-008
- Log markers: [Exercise][create|update|delete|get|list], [ExerciseMedia][upload|download|delete]
- Error codes (REST): FILE_TOO_LARGE, INVALID_FILE_TYPE, NOT_FOUND, INTERNAL_ERROR, UNAUTHORIZED
- Error format: { "error": { "code": "ERROR_CODE", "message": "Human readable" } } per TDEC-027
- CORS: uses existing publicCORS config from main.go with PIN auth middleware

## Security Privacy And Compliance
- All endpoints protected by WAVE-01 PIN auth middleware (GraphQL mutations/queries + REST exercise-media endpoints)
- When PIN is disabled, endpoints accessible without auth (consistent with TDEC-037)
- Server-side MIME detection prevents Content-Type spoofing
- UUID-based parameters prevent path traversal (file path resolved from DB, not user input)
- Uploaded file names sanitized: UUID-based storage, no user-provided path segments
- No sensitive content (personalNotes, file content) logged — log markers record exercise_id, action, success/failure only
- Exercise and ExerciseMedia operations scoped to default user per MVP constraint

## Implementation Slices

| Slice ID | Name | Description |
| --- | --- | --- |
| SLICE-W02-001 | DB migrations | Create goose migrations 00080_exercises.sql and 00081_exercise_media.sql with indexes and FK |
| SLICE-W02-002 | sqlc queries | Define exercise CRUD, list with pagination, allExercises, media CRUD, and soft delete queries |
| SLICE-W02-003 | Exercise repository | Implement ExerciseRepo adapter with sqlc-generated queries and error mapping |
| SLICE-W02-004 | Exercise service | Transport-neutral service with validation (name required, weight > 0, duplicate names allowed) |
| SLICE-W02-005 | GraphQL schema | Add exercises.graphql with Exercise, ExerciseMedia types, CRUD mutations, queries, and union results |
| SLICE-W02-006 | GraphQL resolvers | Implement exercise resolvers with PIN auth guard and WAVE-01 common error types |
| SLICE-W02-007 | ExerciseMedia REST handler | Upload (multipart with validation), download, and delete endpoints |
| SLICE-W02-008 | Main wiring | Wire ExerciseRepo, ExerciseService, ExerciseMediaHandler, and resolver DI; register PIN-protected route group |

## Acceptance Criteria

| AC ID | Description |
| --- | --- |
| AC-W02-001 | Exercise can be created with name (required), muscleGroups, description, personalNotes, workingWeight, isActive via GraphQL mutation. Created exercise returned with generated id and timestamps. |
| AC-W02-002 | Exercise can be read by ID via GraphQL query. Returns full exercise with all fields. |
| AC-W02-003 | Exercises can be listed with cursor pagination via GraphQL query. Default page size applies. Response includes page items and total count. |
| AC-W02-004 | Exercise list pagination cursor works correctly: cursor at end returns no results; cursor in middle returns remaining items. |
| AC-W02-005 | Exercise can be updated (name, muscleGroups, description, personalNotes, workingWeight) via GraphQL mutation. Updated exercise returned. |
| AC-W02-006 | Exercise working weight is stored faithfully: created value equals retrieved value; updated value equals retrieved value. |
| AC-W02-007 | Exercise can be soft-deleted (isActive set to false) via GraphQL mutation. Mutation returns success indicator. |
| AC-W02-008 | Soft-deleted exercise (isActive=false) is excluded from default exercise list. |
| AC-W02-009 | Soft-deleted exercise (isActive=false) can be queried by ID — returns the exercise regardless of isActive status. |
| AC-W02-010 | Soft-deleted exercise can be included in list with includeInactive=true parameter. |
| AC-W02-011 | Reactivation of soft-deleted exercise is not in WAVE-02 scope. Exercise must be re-created or handled via direct DB update. |
| AC-W02-012 | Exercise name is required and validated: mutation returns ValidationError when name is empty or whitespace-only. |
| AC-W02-013 | Duplicate exercise names are allowed (no uniqueness constraint). Creating two exercises with the same name succeeds. |
| AC-W02-014 | Working weight, if provided, must be > 0. Mutation returns ValidationError when weight is <= 0. |
| AC-W02-015 | Media file can be uploaded and associated with an exercise via REST POST /api/v1/exercise-media. Request body: multipart with exerciseId and file. Response: ExerciseMedia JSON. |
| AC-W02-016 | Exercise media file can be downloaded via GET /api/v1/exercise-media/{id}. Returns file with correct content type. |
| AC-W02-017 | Exercise media association can be removed via REST DELETE /api/v1/exercise-media/{id}. Physical file is deleted from disk. Returns 204 No Content. |
| AC-W02-018 | Exercise's media list is returned in GraphQL exercise query (media field on Exercise type). |
| AC-W02-019 | WAVE-02 provides allExercises(includeInactive: Boolean = false): [Exercise!]! GraphQL query returning all active exercises without pagination for WAVE-03 exercise selector. |
| AC-W02-020 | Exercise GraphQL mutations return AuthError when PIN session header is missing or invalid. |
| AC-W02-021 | Exercise media REST endpoints return 401 when PIN session header is missing or invalid. |
| AC-W02-022 | File upload rejects files with disallowed MIME types (only JPEG/PNG/WEBP/MP4/MOV/WEBM allowed). Returns validation error with supported types listed. |
| AC-W02-023 | File upload rejects files exceeding size limits (25MB for images, 250MB for video). Returns validation error. |
| AC-W02-024 | Uploaded file name is sanitized: path separators removed, UUID-based storage name used. No path traversal possible. |

## Exit Criteria

| EC ID | Description |
| --- | --- |
| EC-W02-001 | AC-W02-001 through AC-W02-024 pass via TEST-W02-001 through TEST-W02-022 |
| EC-W02-002 | gqlgen codegen produces valid Go code for Exercise schema without drift |
| EC-W02-003 | sqlc codegen produces valid Go code for exercise queries without drift |
| EC-W02-004 | WAVE-01 media REST scaffold extended for exercise-media association. Existing media endpoints unchanged. |
| EC-W02-005 | WAVE-01 PIN auth guard protects all WAVE-02 GraphQL and REST endpoints. Existing admin auth unchanged. |
| EC-W02-006 | WAVE-01 admin auth and health test suite still passes after WAVE-02 changes |
| EC-W02-007 | Lint passes for all changed packages |
| EC-W02-008 | Migrations 00080 (exercises) and 00081 (exercise_media) apply and roll back in sequence without errors |
| EC-W02-009 | File size and type validation enforced for exercise media uploads per TDEC-008 |
| EC-W02-010 | Exercise and ExerciseMedia operations appear in audit log markers |
| EC-W02-011 | No sensitive content (personalNotes, media file content) appears in application logs |
| EC-W02-012 | Exercise round-trip integration test passes: create exercise -> upload media -> verify media -> delete media -> soft-delete -> verify inactive |
| EC-W02-013 | allExercises query works for WAVE-03 exercise selector dependency |

## Verification Obligations

| Test ID | Description | Type | Command |
| --- | --- | --- | --- |
| TEST-W02-001 | Exercise repository CRUD unit tests | unit | bunx nx run api:test -- --run '(?i)exercise_repo' |
| TEST-W02-002 | Exercise service validation tests (name required, weight > 0) | unit | bunx nx run api:test -- --run '(?i)exercise_service' |
| TEST-W02-003 | Exercise GraphQL resolver integration tests (union results) | integration | bunx nx run api:test -- --run '(?i)exercise_resolver' |
| TEST-W02-004 | ExerciseMedia REST handler integration tests | integration | bunx nx run api:test -- --run '(?i)exercise_media_handler' |
| TEST-W02-005 | Migration smoke test (00080 + 00081 up + down) | integration | bunx nx run api:test -- --run '(?i)migration' |
| TEST-W02-006 | Codegen drift check (gqlgen + sqlc) | codegen | bunx nx run api:codegen && bunx nx run graphql:codegen |
| TEST-W02-007 | allExercises query for WAVE-03 dependency | integration | bunx nx run api:test -- --run '(?i)exercise_list_all' |
| TEST-W02-008 | Soft delete referential integrity | integration | bunx nx run api:test -- --run '(?i)exercise_soft_delete' |
| TEST-W02-009 | isActive filter: default list excludes inactive, includeInactive=true includes | integration | bunx nx run api:test -- --run '(?i)exercise_active_filter' |
| TEST-W02-010 | Exercise list pagination (cursor behavior, edge cases) | integration | bunx nx run api:test -- --run '(?i)exercise_pagination' |
| TEST-W02-011 | Duplicate exercise names allowed | integration | bunx nx run api:test -- --run '(?i)exercise_duplicate_name' |
| TEST-W02-012 | Exercise field update persistence | integration | bunx nx run api:test -- --run '(?i)exercise_update' |
| TEST-W02-013 | Exercise GraphQL returns AuthError without valid PIN session | integration | bunx nx run api:test -- --run '(?i)exercise_auth' |
| TEST-W02-014 | ExerciseMedia upload returns 401 without valid PIN session | integration | bunx nx run api:test -- --run '(?i)exercise_media_auth' |
| TEST-W02-015 | File type validation rejects unauthorized MIME types | unit | bunx nx run api:test -- --run '(?i)exercise_media_filetype' |
| TEST-W02-016 | File size validation rejects oversized uploads | unit | bunx nx run api:test -- --run '(?i)exercise_media_filesize' |
| TEST-W02-017 | Path traversal prevention in upload handler | unit | bunx nx run api:test -- --run '(?i)exercise_media_path_traversal' |
| TEST-W02-018 | Log privacy: personalNotes not appearing in logs | unit | bunx nx run api:test -- --run '(?i)exercise_log_sanitize' |
| TEST-W02-019 | Go lint for API package | lint | bunx nx run api:lint |
| TEST-W02-020 | GraphQL schema validate | codegen | bunx nx run graphql:validate |
| TEST-W02-021 | Exercise round-trip integration test (full lifecycle) | integration | bunx nx run api:test -- --run '(?i)exercise_roundtrip' |
| TEST-W02-022 | WAVE-01 admin auth regression test | unit | bunx nx run api:test -- --run '(?i)admin_auth' |

## Rollout Rollback And Compatibility
- Rollout: merge PR, CI builds and runs tests, deploy via Dokploy compose update. New services start alongside existing WAVE-01 infrastructure.
- Rollback: revert PR, CI builds previous image, Dokploy compose update rolls back. Existing media scaffold and admin auth unaffected.
- Compatibility: all new operations are additive. No existing API changes. WAVE-01 endpoints (health, admin GraphQL, media REST scaffold, users REST) unchanged.
- Migration: goose migrations 00080 and 00081 run at startup. Down migrations available for rollback.
- WAVE-03 compatibility: allExercises query interface is stable. WAVE-02 can be deployed before WAVE-03 without breaking anything.

## Handoff Packets
- HANDOFF-W02-001: This wave brief document
- HANDOFF-W02-002: Planner reports (6 scopes, 2 cycles)
- HANDOFF-W02-003: Reviewer evidence (7 perspectives, 2 cycles)
- HANDOFF-W02-004: Final fit review evidence

## Reviewer Verdicts

| Wave | Perspective | Attempt | Verdict | Reviewer Report | Required Revisions | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| WAVE-02 | product-scope-and-ac | 1 | needs-revision | review-product-scope-and-ac-attempt-1.md | AC deduplication, exercise lifecycle boundary, media lifecycle edge cases | Revised in cycle 2 |
| WAVE-02 | product-scope-and-ac | 2 | approved | review-product-scope-and-ac-attempt-2.md | none | 24 ACs cover all scope, edge cases documented |
| WAVE-02 | architecture-codebase-fit | 1 | needs-revision | review-architecture-codebase-fit-attempt-1.md | WAVE-01 dependency contract explicit, codegen auto-discovery, resolver DI | Revised in cycle 2 |
| WAVE-02 | architecture-codebase-fit | 2 | approved | review-architecture-codebase-fit-attempt-2.md | none | Codebase touchpoints well-documented |
| WAVE-02 | data-api-integration-ops | 1 | needs-revision | review-data-api-integration-ops-attempt-1.md | pg_trgm removed, ON DELETE CASCADE changed to NO ACTION, GET endpoint added | Revised in cycle 2 |
| WAVE-02 | data-api-integration-ops | 2 | approved | review-data-api-integration-ops-attempt-2.md | none | Data/API/ops coverage adequate |
| WAVE-02 | security-privacy-compliance | 1 | needs-revision | review-security-privacy-compliance-attempt-1.md | MIME detection, PIN-disabled access, CORS, log privacy | Revised in cycle 2 |
| WAVE-02 | security-privacy-compliance | 2 | approved | review-security-privacy-compliance-attempt-2.md | none | Server-side MIME, file validation, log privacy covered |
| WAVE-02 | testing-exit-criteria | 1 | needs-revision | review-testing-exit-criteria-attempt-1.md | Round-trip test added, EC strength, fixture strategy | Revised in cycle 2 |
| WAVE-02 | testing-exit-criteria | 2 | approved | review-testing-exit-criteria-attempt-2.md | none | 22 test obligations cover all AC and EC |
| WAVE-02 | sequencing-other-wave-fit | 1 | needs-revision | review-sequencing-other-wave-fit-attempt-1.md | WAVE-01 block, allExercises interface, WAVE-06 data flow correction | Revised in cycle 2 |
| WAVE-02 | sequencing-other-wave-fit | 2 | approved | review-sequencing-other-wave-fit-attempt-2.md | none | Dependency order correct, no collision |
| WAVE-02 | traceability-consistency | 1 | needs-revision | review-traceability-consistency-attempt-1.md | Stable IDs, source traces, ledger consistency | Revised in cycle 2 |
| WAVE-02 | traceability-consistency | 2 | approved | review-traceability-consistency-attempt-2.md | none | Source traceability documented |
| WAVE-02 | final-wave-fit-review | 1 | approved | final-wave-fit-review-attempt-1.md | none | Package is ready-for-dev |

## Open Questions

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W02-001 | WAVE-02 | data-ops | wave-blocking | EDGE-020 | Should deleting ExerciseMedia delete the physical media file from disk? | Orphaned files accumulate | Yes, per TDEC-005. Add failure handling (log error if file deletion fails). | planner-data-integration-ops-attempt-2.md | resolved | Physical file deleted. On failure, log error and return 204 to client. |
| DQ-W02-002 | WAVE-02 | product | needs-owner-decision | AC-043 | Are exercise names unique per user or can duplicates exist? | Affects validation logic and UI behavior | Duplicates allowed (no constraint) — consistent with EDGE-002. | planner-product-ac-attempt-2.md | answered | Tentative: duplicates allowed per EDGE-002. Awaiting user confirmation. |
| DQ-W02-003 | WAVE-02 | data-ops | wave-blocking | WAVE-01 | What exact file storage path pattern does WAVE-01 MediaConfig provide for exercise media? | Drives migration and handler design | Use WAVE-01 BasePath/<exercise_id>/<uuid>.<ext>. Confirm after WAVE-01 implementation. | planner-data-integration-ops-attempt-2.md | deferred | WAVE-01 coordination item. WAVE-02 assumes composable BasePath. |
| DQ-W02-005 | WAVE-02 | security | needs-owner-decision | TDEC-008 | Should WAVE-02 use server-side MIME detection (file magic bytes) or trust Content-Type header? | Content-Type can be spoofed; magic bytes are more secure | Use http.DetectContentType() cross-check against provided Content-Type. | planner-security-compliance-attempt-2.md | resolved | Decision: http.DetectContentType() server-side as primary check. |
| DQ-W02-006 | WAVE-02 | security | deferred | EDGE-014 | Should exercise media URLs be time-limited (signed URLs) or always accessible with valid session? | Signed URLs add complexity for single-user MVP | Session-gated access sufficient for MVP self-hosted deployment. | planner-security-compliance-attempt-2.md | deferred | Deferred post-MVP. |
| DQ-W02-007 | WAVE-02 | testing | needs-owner-decision | WAVE-01 | Should exercise tests use mocked PIN auth or integration through full middleware chain? | Test complexity vs realism | Prefer integration through full middleware chain with WAVE-01 PIN test helpers. | planner-testing-exit-attempt-2.md | resolved | Decision: full middleware chain integration tests. |
| DQ-W02-008 | WAVE-02 | sequencing | watchlist | WAVE-03 | Does allExercises need filtering beyond isActive for WAVE-03 exercise selector? | Might be needed if library grows large | Deferred — current scope is unfiltered active list ordered by name. | planner-sequencing-fit-attempt-2.md | deferred | Watchlist. Not needed for MVP. |

## Traceability
- docs/prd-waves/waves/wave-02.md: source wave boundary, outcomes, capability groups
- docs/product-verified/functional-spec.md: Exercise Library §11 — REQ-003
- docs/product-verified/domain-model.md: Exercise, ExerciseMedia entities
- docs/product-verified/acceptance-criteria.md: AC-002, AC-003, AC-004, AC-043, AC-044, AC-045, AC-046, AC-047
- docs/technical-verified/api-contracts.md: hybrid GraphQL/REST protocol, TDEC-001
- docs/technical-verified/data-contracts.md: exercise entity contracts, TDEC-020, TDEC-022, TDEC-023
- docs/technical-verified/implementation-slices.md: Slice 1 Exercise Library mapping
- docs/technical-verified/auth-security-compliance.md: PIN auth, TDEC-037, TDEC-008
- docs/technical-verified/operations-observability.md: log markers, error format
- docs/development-plan.xml: M-API, M-PRD-WAVE-DETAILER module contracts
- docs/knowledge-graph.xml: existing module boundaries
- docs/prd-wave-details/waves/wave-01.md: WAVE-01 dependency contracts
- apps/api/internal: existing codebase patterns for service/repository/middleware/handler structure
- docs/prd-waves/frontend-pages/page-002.md: workout diary backend dependencies
- docs/prd-waves/frontend-pages/page-003.md: exercise library backend dependencies