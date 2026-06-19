# Wave 02: Exercise Library

## Status
user-approved

## User Approval
approved-by-user: 2026-06-19 — user corrected spec with 9 items: (1) user_id added to both tables; (2) indexes prefixed with user_id; (3) media uses WAVE-01 scaffold routes (POST/GET/DELETE /api/v1/media/* with purpose=EXERCISE_MEDIA + exerciseId), no new /api/v1/exercise-media namespace; (4) Exercise CRUD is GraphQL only — no REST; (5) deleteExercise replaced with archiveExercise/restoreExercise; (6) archiving exercise does NOT cascade delete media; (7) exercise(id: ID!) returns ExerciseResult union; (8) duplicate names allowed (unchanged); (9) AI/future references use exerciseId as primary identity, name as display. Final consistency fixes applied: (A) codebase touchpoints aligned with WAVE-01 Approach A — atlas-specific paths for service, repo, graph/schema, graph/resolver, atlas-gqlgen.yml; (B) exercise_media.exercise_id NOT NULL; (C) working_weight CHECK (IS NULL OR > 0) at DB level. Scope: Exercise CRUD via GraphQL, muscle groups, working weight, description/notes, active/inactive behavior, media upload/download/delete via WAVE-01 media scaffold routes with purpose/entity metadata, PIN-auth protection through WAVE-01 Atlas PIN middleware, multi-user-ready (user_id FK on all entities). Boundary: no DailyLog/workout, sets, cardio, body tracking, nutrition, charts, AI export, backup/import. Dependency: WAVE-01 provides Atlas PIN middleware, media REST scaffold (501 — extended by WAVE-02), GraphQL foundation/common types, default user context, and settings/session infrastructure.

## Source Wave Summary
WAVE-02 from docs/prd-waves/waves/wave-02.md. Full CRUD for exercises with working weight and media management. Source status: user-approved (2026-06-18).

## Outcome After Implementation
- OUT-W02-001: Exercises can be created, listed, edited, archived/restored via GraphQL
- OUT-W02-002: Media (images/video) can be attached, retrieved, and deleted per exercise via WAVE-01 media scaffold routes
- OUT-W02-003: Working weight stored per exercise, snapshot-ready for WAVE-03
- OUT-W02-004: API ready for workout diary (allExercises query, exercise by ID)

## Scope Included
- CAP-W02-001: Exercise CRUD via GraphQL (create, read, list with pagination, update, soft archive/restore)
- CAP-W02-002: ExerciseMedia upload, download, and deletion via WAVE-01 media scaffold routes with purpose/entity metadata
- CAP-W02-003: Working weight field (Float, validated > 0)
- CAP-W02-004: Muscle groups (string array), description, personal notes fields
- CAP-W02-005: isActive flag for soft archive (default true, exclude from default list)

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
- Frontend backend dependencies from PAGE-003: GET/POST/PUT/DELETE /api/exercises (via GraphQL mutations/queries), POST/DELETE /api/v1/media/upload + DELETE /api/v1/media/{id} (via WAVE-01 REST media scaffold).

## Codebase Fit And Touchpoints
- apps/api/internal/repository/postgres/migrations/00081_exercises.sql: new migration for exercises table (with user_id FK, working_weight CHECK)
- apps/api/internal/repository/postgres/migrations/00082_exercise_media.sql: new migration for exercise_media table (with user_id FK, exercise_id NOT NULL)
- apps/api/internal/repository/postgres/queries/exercises.sql: sqlc query definitions
- apps/api/internal/atlas/repository/postgres/exercise_repo.go: repository adapter for exercise CRUD (imports from internal/repository/postgres/generated)
- apps/api/internal/atlas/service/exercise.go: transport-neutral exercise service with validation
- apps/api/internal/handler/atlas_media.go: EXTENDED by WAVE-02 — implement upload/download/delete logic with purpose/entity routing (EXERCISE_MEDIA)
- apps/api/internal/atlas/graph/resolver/exercise.resolvers.go: GraphQL resolvers for exercise queries/mutations
- apps/api/internal/atlas/graph/schema/exercises.graphql: exercise GraphQL types and operations
- apps/api/cmd/server/main.go: wire ExerciseRepo, ExerciseService, resolver DI; media routes already registered by WAVE-01 scaffold, no new route group needed
- apps/api/internal/appconfig/config.go: reads MediaConfig from WAVE-01 (BasePath, MaxUploadSize) — implement if not present
- apps/api/atlas-gqlgen.yml: auto-discovers new schema files from internal/atlas/graph/schema via glob
- apps/api/sqlc.yaml: auto-discovers new queries from internal/repository/postgres/queries via glob

## Design Contracts
- Soft archive: isActive=false preserves referential integrity for WAVE-03 workout history. Archiving exercise does NOT cascade delete exercise_media (DDEC-W02-001)
- Physical file deletion: ExerciseMedia delete (via WAVE-01 media scaffold) removes file from disk; soft-fail (log error, return 204) on deletion failure (DDEC-W02-002)
- Duplicate names: allowed, no unique constraint (DDEC-W02-003, EDGE-002)
- allExercises: GraphQL-only query returning active exercises ordered by name ASC (DDEC-W02-004)
- MIME detection: http.DetectContentType() server-side (512 bytes) as primary validation (DQ-W02-005 resolved)
- PIN auth: WAVE-01 middleware guards all WAVE-02 GraphQL and REST endpoints
- File storage: <WAVE-01 BasePath>/exercise/<exercise_id>/<uuid>.<ext>
- Working weight: REAL type, validated > 0 when provided. DB-level CHECK (working_weight IS NULL OR working_weight > 0) enforces invariant.
- exercise_id in exercise_media: UUID NOT NULL with FK ON DELETE NO ACTION — every media record is always attached to an exercise
- user_id FK: all exercises and exercise_media records are scoped to default Atlas user (multi-user-ready)
- exercise(id: ID!) returns ExerciseResult union (Exercise | ValidationError | NotFoundError | AuthError) — never raw Exercise type
- archiveExercise(id: ID!): ArchiveResult union — replaces deleteExercise. restoreExercise optional.
- AI/future export references must use exerciseId as primary identity, name as display field only.

## Data API Integration And Operations
- PostgreSQL: goose migrations 00081_exercises.sql and 00082_exercise_media.sql
- Indexes: idx_exercises_user_active (user_id, is_active), idx_exercises_user_name (user_id, name), idx_exercises_user_created_at (user_id, created_at, optional), idx_exercise_media_user_exercise (user_id, exercise_id)
- FK constraints: exercise_media.exercise_id → exercises(id) ON DELETE NO ACTION, exercise_id NOT NULL; exercises.user_id → atlas_users(id); exercise_media.user_id → atlas_users(id)
- Working weight: REAL type, CHECK (working_weight IS NULL OR working_weight > 0) at DB level
- File validation: server-side MIME detection (JPEG/PNG/WEBP/MP4/MOV/WEBM), per-type size limits (25MB images, 250MB video, 300MB max single upload)
- Memory-safe upload: r.ParseMultipartForm(maxBytes) with 300MB limit per TDEC-008
- Media routes: WAVE-01 scaffold routes POST /api/v1/media/upload, GET /api/v1/media/{id}, DELETE /api/v1/media/{id}. WAVE-02 extends handlers with purpose/entity routing. Upload includes purpose=EXERCISE_MEDIA + exerciseId metadata.
- Log markers: [Exercise][create|update|archive|restore|get|list], [Media][upload|download|delete]
- Error codes (REST): FILE_TOO_LARGE, INVALID_FILE_TYPE, NOT_FOUND, INTERNAL_ERROR, UNAUTHORIZED
- Error format: { "error": { "code": "ERROR_CODE", "message": "Human readable" } } per TDEC-027
- CORS: uses existing publicCORS config from main.go with PIN auth middleware

## Security Privacy And Compliance
- All GraphQL endpoints protected by WAVE-01 PIN auth middleware. Media endpoints protected by WAVE-01 middleware chain (AtlasUserContext + AtlasPinGuard).
- When PIN is disabled, endpoints accessible without auth (consistent with TDEC-037)
- Server-side MIME detection prevents Content-Type spoofing
- UUID-based parameters prevent path traversal (file path resolved from DB, not user input)
- Uploaded file names sanitized: UUID-based storage, no user-provided path segments
- No sensitive content (personalNotes, file content) logged — log markers record action, success/failure only
- Exercise and ExerciseMedia operations scoped to default user via user_id FK (multi-user-ready)

## Implementation Slices

| Slice ID | Name | Description |
| --- | --- | --- |
| SLICE-W02-001 | DB migrations | Create goose migrations 00081_exercises.sql and 00082_exercise_media.sql with user_id FK, NOT NULL exercise_id, working_weight CHECK, indexed for multi-user |
| SLICE-W02-002 | sqlc queries | Define exercise CRUD, list with pagination, allExercises, media CRUD, and soft archive queries |
| SLICE-W02-003 | Exercise repository | Implement ExerciseRepo adapter in atlas/repository/postgres with sqlc-generated queries and error mapping |
| SLICE-W02-004 | Exercise service | Transport-neutral service in atlas/service with validation (name required, weight > 0, duplicate names allowed) |
| SLICE-W02-005 | GraphQL schema | Add exercises.graphql under atlas/graph/schema with Exercise, ExerciseMedia types, CRUD mutations (archive/restore instead of delete), queries (exercise returns union), and union results |
| SLICE-W02-006 | GraphQL resolvers | Implement exercise resolvers in atlas/graph/resolver with PIN auth guard and WAVE-01 common error types |
| SLICE-W02-007 | WAVE-01 media scaffold extension | Extend handler/atlas_media.go: implement upload with purpose/entity routing (EXERCISE_MEDIA), download, delete with file validation and storage |
| SLICE-W02-008 | Main wiring | Wire ExerciseRepo, ExerciseService, resolver DI in main.go; media routes already registered by WAVE-01 scaffold |

## Acceptance Criteria

| AC ID | Description |
| --- | --- |
| AC-W02-001 | Exercise can be created with name (required), muscleGroups, description, personalNotes, workingWeight, isActive via GraphQL mutation. Created exercise returned with generated id and timestamps. |
| AC-W02-002 | Exercise can be read by ID via GraphQL query. Returns full exercise with all fields. |
| AC-W02-003 | Exercises can be listed with cursor pagination via GraphQL query. Default page size applies. Response includes page items and total count. |
| AC-W02-004 | Exercise list pagination cursor works correctly: cursor at end returns no results; cursor in middle returns remaining items. |
| AC-W02-005 | Exercise can be updated (name, muscleGroups, description, personalNotes, workingWeight) via GraphQL mutation. Updated exercise returned. |
| AC-W02-006 | Exercise working weight is stored faithfully: created value equals retrieved value; updated value equals retrieved value. |
| AC-W02-007 | Exercise can be archived (isActive set to false) via GraphQL archiveExercise mutation. Archived exercise returned. |
| AC-W02-008 | Archived exercise (isActive=false) is excluded from default exercise list. |
| AC-W02-009 | Archived exercise (isActive=false) can be queried by ID via exercise(id: ID!): ExerciseResult! — returns the exercise regardless of isActive status via successful result union member. |
| AC-W02-010 | Archived exercise can be included in list with includeInactive=true parameter. |
| AC-W02-011 | Archived exercise can be restored (isActive set to true) via GraphQL restoreExercise mutation. Restored exercise returned. |
| AC-W02-012 | Exercise name is required and validated: mutation returns ValidationError when name is empty or whitespace-only. |
| AC-W02-013 | Duplicate exercise names are allowed (no uniqueness constraint). Creating two exercises with the same name succeeds. |
| AC-W02-014 | Working weight, if provided, must be > 0. Mutation returns ValidationError when weight is <= 0. |
| AC-W02-015 | Media file can be uploaded and associated with an exercise via REST POST /api/v1/media/upload with purpose=EXERCISE_MEDIA and exerciseId. Response: ExerciseMedia JSON. |
| AC-W02-016 | Exercise media file can be downloaded via GET /api/v1/media/{id}. Returns file with correct content type. |
| AC-W02-017 | Exercise media association can be removed via REST DELETE /api/v1/media/{id}. Physical file is deleted from disk. Returns 204 No Content. Archiving the exercise does NOT automatically delete its media records or files. |
| AC-W02-018 | Exercise's media list is returned in GraphQL exercise query (media field on Exercise type). |
| AC-W02-019 | WAVE-02 provides allExercises(includeInactive: Boolean = false): [Exercise!]! GraphQL query returning all active exercises without pagination for WAVE-03 exercise selector. |
| AC-W02-020 | Exercise GraphQL mutations return AuthError when PIN session header is missing or invalid. |
| AC-W02-021 | Media REST endpoints return 401 when PIN session header is missing or invalid. |
| AC-W02-022 | File upload rejects files with disallowed MIME types (only JPEG/PNG/WEBP/MP4/MOV/WEBM allowed). Returns validation error with supported types listed. |
| AC-W02-023 | File upload rejects files exceeding size limits (25MB for images, 250MB for video). Returns validation error. |
| AC-W02-024 | Uploaded file name is sanitized: path separators removed, UUID-based storage name used. No path traversal possible. |
| AC-W02-025 | All exercises and exercise_media records are scoped to the requesting user via user_id FK. Operations only return/modify records for the authenticated user's context. |
| AC-W02-026 | Archiving an exercise (isActive=false) does NOT cascade delete associated exercise_media records or physical files. Media remains queryable. |
| AC-W02-027 | exercise(id: ID!): ExerciseResult! returns NotFoundError when exercise does not exist or belongs to a different user. |

## Exit Criteria

| EC ID | Description |
| --- | --- |
| EC-W02-001 | AC-W02-001 through AC-W02-027 pass via TEST-W02-001 through TEST-W02-022 |
| EC-W02-002 | atlas-gqlgen codegen produces valid Go code for Exercise schema without drift |
| EC-W02-003 | sqlc codegen produces valid Go code for exercise queries without drift |
| EC-W02-004 | WAVE-01 media REST scaffold extended for exercise-media upload/download/delete with purpose/entity routing. Existing media scaffold routes unchanged. |
| EC-W02-005 | WAVE-01 PIN auth guard protects all WAVE-02 GraphQL endpoints. Media endpoints protected by WAVE-01 middleware chain (unchanged). |
| EC-W02-006 | WAVE-01 admin auth and health test suite still passes after WAVE-02 changes |
| EC-W02-007 | Lint passes for all changed packages |
| EC-W02-008 | Migrations 00081 (exercises) and 00082 (exercise_media) apply and roll back in sequence without errors |
| EC-W02-009 | File size and type validation enforced for exercise media uploads per TDEC-008 |
| EC-W02-010 | Exercise CRUD operations appear in log markers |
| EC-W02-011 | No sensitive content (personalNotes, media file content) appears in application logs |
| EC-W02-012 | Exercise round-trip integration test passes: create exercise -> upload media -> verify media -> delete media -> archive exercise -> restore exercise -> verify final state |
| EC-W02-013 | allExercises query works for WAVE-03 exercise selector dependency |

## Verification Obligations

| Test ID | Description | Type | Command |
| --- | --- | --- | --- |
| TEST-W02-001 | Exercise repository CRUD unit tests | unit | bunx nx run api:test -- --run '(?i)exercise_repo' |
| TEST-W02-002 | Exercise service validation tests (name required, weight > 0) | unit | bunx nx run api:test -- --run '(?i)exercise_service' |
| TEST-W02-003 | Exercise GraphQL resolver integration tests (union results, archive/restore) | integration | bunx nx run api:test -- --run '(?i)exercise_resolver' |
| TEST-W02-004 | ExerciseMedia integration via WAVE-01 media scaffold (upload/download/delete with purpose routing) | integration | bunx nx run api:test -- --run '(?i)atlas_media' |
| TEST-W02-005 | Migration smoke test (00080 + 00081 up + down) | integration | bunx nx run api:test -- --run '(?i)migration' |
| TEST-W02-006 | Codegen drift check (sqlc + atlas gqlgen) | codegen | bunx nx run api:codegen && bunx nx run api:codegen:atlas |
| TEST-W02-007 | allExercises query for WAVE-03 dependency | integration | bunx nx run api:test -- --run '(?i)exercise_list_all' |
| TEST-W02-008 | Soft archive referential integrity (archive does NOT cascade media delete) | integration | bunx nx run api:test -- --run '(?i)exercise_archive' |
| TEST-W02-009 | isActive filter: default list excludes inactive, includeInactive=true includes | integration | bunx nx run api:test -- --run '(?i)exercise_active_filter' |
| TEST-W02-010 | Exercise list pagination (cursor behavior, edge cases) | integration | bunx nx run api:test -- --run '(?i)exercise_pagination' |
| TEST-W02-011 | Duplicate exercise names allowed | integration | bunx nx run api:test -- --run '(?i)exercise_duplicate_name' |
| TEST-W02-012 | Exercise field update persistence | integration | bunx nx run api:test -- --run '(?i)exercise_update' |
| TEST-W02-013 | Exercise GraphQL returns AuthError without valid PIN session | integration | bunx nx run api:test -- --run '(?i)exercise_auth' |
| TEST-W02-014 | Exercise exercise(id: ID!) returns NotFoundError for missing or wrong-user exercise | integration | bunx nx run api:test -- --run '(?i)exercise_not_found' |
| TEST-W02-015 | File type validation rejects unauthorized MIME types | unit | bunx nx run api:test -- --run '(?i)media_filetype' |
| TEST-W02-016 | File size validation rejects oversized uploads | unit | bunx nx run api:test -- --run '(?i)media_filesize' |
| TEST-W02-017 | Path traversal prevention in upload handler | unit | bunx nx run api:test -- --run '(?i)media_path_traversal' |
| TEST-W02-018 | Log privacy: personalNotes not appearing in logs | unit | bunx nx run api:test -- --run '(?i)exercise_log_sanitize' |
| TEST-W02-019 | Go lint for API package | lint | bunx nx run api:lint |
| TEST-W02-020 | GraphQL schema validate | codegen | bunx nx run graphql:validate |
| TEST-W02-021 | Exercise round-trip integration test (create, upload media, verify, delete media, archive, restore, verify final state) | integration | bunx nx run api:test -- --run '(?i)exercise_roundtrip' |
| TEST-W02-022 | WAVE-01 admin auth regression test | unit | bunx nx run api:test -- --run '(?i)admin_auth' |

## Rollout Rollback And Compatibility
- Rollout: merge PR, CI builds and runs tests, deploy via Dokploy compose update. New services start alongside existing WAVE-01 infrastructure.
- Rollback: revert PR, CI builds previous image, Dokploy compose update rolls back. Existing media scaffold and admin auth unaffected.
- Compatibility: all new operations are additive. No existing API changes. WAVE-01 endpoints (health, admin GraphQL, media REST scaffold, users REST) unchanged.
- Migration: goose migrations 00081 and 00082 run at startup. Down migrations available for rollback.
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
| DQ-W02-002 | WAVE-02 | product | needs-owner-decision | AC-043 | Are exercise names unique per user or can duplicates exist? | Affects validation logic and UI behavior | Duplicates allowed (no constraint) — consistent with EDGE-002. | planner-product-ac-attempt-2.md | resolved | User approved: duplicates allowed. No UNIQUE(user_id, name). Name required, trimmed, no empty/whitespace-only. Duplicate normalized names allowed. Non-blocking UI warning allowed but save must succeed. Identity always by exercise.id. AI export/workout references include both exerciseId and name. |
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
- docs/prd-wave-details/waves/wave-01.md: WAVE-01 dependency contracts, media scaffold routes
- apps/api/internal: existing codebase patterns for service/repository/middleware/handler structure
- docs/prd-waves/frontend-pages/page-002.md: workout diary backend dependencies
- docs/prd-waves/frontend-pages/page-003.md: exercise library backend dependencies
- Future reference: AI/WAVE-07/WAVE-08 export references must use exerciseId as primary identity and name as display field. Identity is always exercise.id, never exercise.name.