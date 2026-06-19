# WAVE-02 product-ac Planner Attempt 1

## Sources Read
- docs/prd-waves/waves/wave-02.md (selected source wave)
- docs/product-verified/acceptance-criteria.md (AC-002, AC-003, AC-004, AC-043, AC-044, AC-045, AC-046, AC-047)
- docs/product-verified/functional-spec.md (Section 11 - Exercise Library)
- docs/product-verified/domain-model.md (Exercise, ExerciseMedia entities)
- docs/product-verified/business-rules.md (RULE-003)
- docs/product-verified/edge-cases.md (EDGE-002, EDGE-018, EDGE-020)
- docs/product-verified/user-flows.md
- docs/prd-waves/frontend-pages/page-003.md (Exercise Library - primary consumer)
- docs/prd-waves/frontend-pages/page-002.md (Workout Diary - depends on GET /api/exercises)

## Selected Backend Wave Boundary
WAVE-02 implements full CRUD for exercises (GraphQL mutations/queries) with media association (REST upload/download through WAVE-01 scaffold). Working weight, muscle groups, description, notes, and isActive (soft delete) are included. This is backend-only — no frontend pages, UI, or UX.

## Neighboring Backend Wave Fit
- WAVE-01 (Foundation): provides db migration infrastructure, PIN auth middleware, media REST scaffold (POST/GET /api/v1/media), fitness GraphQL schema extension pattern. WAVE-02 builds on all of these.
- WAVE-03 (Workout Diary): depends on WAVE-02's GET /api/exercises for exercise selection in workout diary. No scope collision — WAVE-02 owns exercise CRUD, WAVE-03 owns workout logging.
- No scope overlap with WAVE-04 (Cardio), WAVE-05 (Nutrition), WAVE-06 (Charts), WAVE-07/08 (AI), WAVE-09 (Backup).

## Frontend Pages Context
- PAGE-003 (Exercise Library): reads exercises (GET), creates (POST), edits (PUT), deletes (DELETE), manages media (POST/DELETE exercise-media). Backend provides all these operations.
- PAGE-002 (Workout Diary): depends on GET /api/exercises for exercise selection list. Only read dependency.
- No frontend planning, UI, or UX in this wave.

## Codebase Evidence
- No existing exercise-related SQL, GraphQL schema, repository, service, handler, or resolver code exists.
- WAVE-01 will establish: media REST scaffold (POST /api/v1/media/upload, GET /api/v1/media/{id}), fitness GraphQL schema pattern (extend type Query/Mutation), sqlc query pattern, migration pattern, PIN auth guard middleware.
- Existing patterns: user_repo.go (sqlc adapter), admin_auth.go (service/repository pattern), users.go (handler pattern), resolver pattern in graph/.

## Proposed Details

### Exercise CRUD (GraphQL)
- CreateExercise: name (required), muscleGroups (string array), description (optional), personalNotes (optional), workingWeight (float, optional), isActive (default true)
- UpdateExercise: all optional fields, id required
- DeleteExercise: soft-delete via isActive = false (not hard delete due to WAVE-03 referential integrity)
- GetExercise: single exercise by id
- ListExercises: paginated list, filterable by isActive, search by name
- ListAllExercises: simple list without pagination for WAVE-03 exercise selector

### ExerciseMedia Association (REST)
- POST /api/v1/exercise-media: multipart upload, associates media with exerciseId
- DELETE /api/v1/exercise-media/{id}: removes media association and optionally physical file
- GET /api/v1/exercise-media/{id}: returns media file (reuses WAVE-01 GET /api/v1/media/{id} pattern if media stored the same way)
- List media for exercise: GraphQL query on exercise.media

### Working Weight
- Stored on Exercise entity as workingWeight (float)
- WAVE-03 captures snapshot at workout-log time (RULE-017)
- No history tracking in WAVE-02 — working weight is current value only

### Soft Delete
- isActive flag on Exercise
- DELETE mutation sets isActive = false
- List queries respect isActive filter (default: only active, or include inactive)
- Referential integrity: exercise rows preserved for WAVE-03 workout history

### Validation
- Name: required, non-empty, trimmed
- Working weight: nullable float, must be > 0 if provided
- Muscle groups: array of strings, each non-empty
- Exercise name uniqueness: not strictly required per EDGE-002 (no duplicate rule in PRD), but backend should handle gracefully

## Acceptance Criteria Contributions

| AC ID | Description |
| --- | --- |
| AC-W02-001 | Exercise can be created with name, muscleGroups, description, personalNotes, workingWeight, isActive via GraphQL mutation |
| AC-W02-002 | Exercise can be read by ID via GraphQL query |
| AC-W02-003 | Exercises can be listed with pagination via GraphQL query |
| AC-W02-004 | Exercise can be updated (all optional fields) via GraphQL mutation |
| AC-W02-005 | Exercise can be soft-deleted (isActive set to false) via GraphQL mutation |
| AC-W02-006 | Media file can be uploaded and associated with an exercise via REST POST /api/v1/exercise-media |
| AC-W02-007 | Media association can be removed via REST DELETE /api/v1/exercise-media/{id} |
| AC-W02-008 | Exercise's media list is returned in GraphQL exercise query |
| AC-W02-009 | Working weight is stored and retrievable on exercise |
| AC-W02-010 | Exercise name is required and validated |
| AC-W02-011 | Exercise with isActive=false is excluded from default list queries |
| AC-W02-012 | Exercise with isActive=false can be queried with explicit includeInactive filter |
| AC-W02-013 | Media file download works through WAVE-01 REST scaffold |

## Exit Criteria Contributions

| EC ID | Description |
| --- | --- |
| EC-W02-001 | All acceptance criteria passing in focused unit and integration tests |
| EC-W02-002 | gqlgen codegen produces valid Go code for Exercise CRUD schema |
| EC-W02-003 | sqlc codegen produces valid Go code for exercise queries |
| EC-W02-004 | WAVE-01 media REST scaffold handles exercise-media association extension |
| EC-W02-005 | WAVE-01 PIN auth guard protects all WAVE-02 GraphQL and REST endpoints |
| EC-W02-006 | Existing admin auth and health endpoints still pass unchanged |
| EC-W02-007 | Lint passes for all changed packages |
| EC-W02-008 | No frontend scope changes in this wave |

## Verification Contributions

| Test ID | Description | Type | Command |
| --- | --- | --- | --- |
| TEST-W02-001 | Exercise repository unit tests | unit | bunx nx run api:test -- --run '(?i)exercise_repo' |
| TEST-W02-002 | Exercise service unit tests | unit | bunx nx run api:test -- --run '(?i)exercise_service' |
| TEST-W02-003 | Exercise GraphQL resolver integration tests | integration | bunx nx run api:test -- --run '(?i)exercise_resolver' |
| TEST-W02-004 | ExerciseMedia REST handler tests | integration | bunx nx run api:test -- --run '(?i)exercise_media_handler' |
| TEST-W02-005 | Migration smoke test (up + down) | integration | bunx nx run api:test -- --run '(?i)migration' |
| TEST-W02-006 | Codegen drift check (gqlgen + sqlc) | codegen | bunx nx run api:codegen && bunx nx run graphql:codegen |
| TEST-W02-007 | Exercise list without pagination for WAVE-03 dependency | integration | bunx nx run api:test -- --run '(?i)exercise_list_all' |
| TEST-W02-008 | Soft delete referential integrity check | integration | bunx nx run api:test -- --run '(?i)exercise_soft_delete' |

## Risks And Rollback
- Risk: media files orphaned when exercise-media association deleted but physical file still on disk. Mitigation: WAVE-02 deletes physical file when the last media association is removed, or defers file cleanup to WAVE-01 media lifecycle story. Recommendation: delete physical file when ExerciseMedia record is deleted (consistent with TDEC-005).
- Risk: exercise with isActive=false still referenced by WAVE-03 workout data. Mitigation: soft delete preserves rows; WAVE-03 queries exercises by ID directly, not by active status. No referential integrity break.
- Rollback: revert migration 00080_exercises.sql and 00081_exercise_media.sql, revert code. Existing exercises remain in database; WAVE-02 code removed.

## Questions Raised
| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Source Or Report | Status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W02-001 | WAVE-02 | product-ac | wave-blocking | EDGE-020 | Should deleting ExerciseMedia delete the physical media file from disk? | Orphaned files accumulate; contradicts TDEC-005 media deletion principle | product-ac planner | open |
| DQ-W02-002 | WAVE-02 | product-ac | needs-owner-decision | AC-043 | Are exercise names unique per user or can duplicates exist? | PRD §11 has no duplicate rule (EDGE-002). Affects validation logic and UI behavior. | product-ac planner | open |

## Traceability Candidates
- CAP-W02-001 Exercise CRUD → AC-W02-001 through AC-W02-005, AC-W02-010 through AC-W02-012
- CAP-W02-002 ExerciseMedia → AC-W02-006 through AC-W02-008, AC-W02-013
- CAP-W02-003 Working weight → AC-W02-009
- CAP-W02-004 Muscle groups/description/notes → AC-W02-001, AC-W02-004
- CAP-W02-005 isActive soft delete → AC-W02-005, AC-W02-011, AC-W02-012
- Product AC-002, AC-003, AC-043, AC-044 mapped to AC-W02-001
- Product AC-004, AC-045, AC-046 mapped to AC-W02-006, AC-W02-007
- Product AC-047 mapped to AC-W02-005, AC-W02-011