# WAVE-02 product-ac Planner Attempt 2

## Cycle 1 Reviewer Feedback Addressed
All revision items from 7 reviewers incorporated. Key changes from attempt 1:

## Resolved AC List (Deduplicated, Sequential)

| AC ID | Description | Source Trace |
| --- | --- | --- |
| AC-W02-001 | Exercise can be created with name (required), muscleGroups, description, personalNotes, workingWeight, isActive via GraphQL mutation. Created exercise returned with generated id and timestamps. | OUT-W02-001, CAP-W02-001, AC-002, AC-043, AC-044 |
| AC-W02-002 | Exercise can be read by ID via GraphQL query. Returns full exercise with all fields. | OUT-W02-001, CAP-W02-001 |
| AC-W02-003 | Exercises can be listed with cursor pagination via GraphQL query. Default page size applies. Response includes page items and total count. | OUT-W02-001, CAP-W02-001 |
| AC-W02-004 | Exercise list pagination cursor works correctly: cursor at end returns no results; cursor in middle returns remaining items. | OUT-W02-001, CAP-W02-001 |
| AC-W02-005 | Exercise can be updated (name, muscleGroups, description, personalNotes, workingWeight) via GraphQL mutation. Updated exercise returned. | OUT-W02-001, CAP-W02-001, AC-044 |
| AC-W02-006 | Exercise working weight is stored faithfully: created value equals retrieved value; updated value equals retrieved value. | OUT-W02-003, CAP-W02-003, AC-003, RULE-003 |
| AC-W02-007 | Exercise can be soft-deleted (isActive set to false) via GraphQL mutation. Mutation returns success indicator. | OUT-W02-001, CAP-W02-005, AC-047 |
| AC-W02-008 | Soft-deleted exercise (isActive=false) is excluded from default exercise list. | OUT-W02-001, CAP-W02-005, AC-047 |
| AC-W02-009 | Soft-deleted exercise (isActive=false) can be queried by ID (exercise(id: UUID!)) — returns the exercise regardless of isActive status. | OUT-W02-004, CAP-W02-005 |
| AC-W02-010 | Soft-deleted exercise can be included in list with includeInactive=true parameter. | OUT-W02-001, CAP-W02-005 |
| AC-W02-011 | Reactivation of soft-deleted exercise is not in WAVE-02 scope. Exercise must be re-created or handled via direct DB update. | CAP-W02-005 |
| AC-W02-012 | Exercise name is required and validated: mutation returns ValidationError when name is empty or whitespace-only. | CAP-W02-001, EDGE-002 |
| AC-W02-013 | Duplicate exercise names are allowed (no uniqueness constraint). Creating two exercises with the same name succeeds. | CAP-W02-001, EDGE-002 |
| AC-W02-014 | Working weight, if provided, must be > 0. Mutation returns ValidationError when weight is <= 0. | CAP-W02-003, RULE-003 |
| AC-W02-015 | Media file can be uploaded and associated with an exercise via REST POST /api/v1/exercise-media. Request body: multipart with exerciseId and file. Response: ExerciseMedia JSON. | OUT-W02-002, CAP-W02-002, AC-004, AC-045 |
| AC-W02-016 | Exercise media file can be downloaded via GET /api/v1/exercise-media/{id} (reuses WAVE-01 media delivery pattern). Returns file with correct content type. | OUT-W02-002, CAP-W02-002 |
| AC-W02-017 | Exercise media association can be removed via REST DELETE /api/v1/exercise-media/{id}. Physical file is deleted from disk. Returns 204 No Content. | OUT-W02-002, CAP-W02-002, AC-046, TDEC-005 |
| AC-W02-018 | Exercise's media list is returned in GraphQL exercise query (exercises field on Exercise type). | CAP-W02-002, AC-045 |
| AC-W02-019 | WAVE-02 provides `allExercises(includeInactive: Boolean = false): [Exercise!]!` GraphQL query returning all active exercises without pagination for WAVE-03 exercise selector. | OUT-W02-004, PAGE-002, PAGE-003 |
| AC-W02-020 | Exercise GraphQL mutations return AuthError when PIN session header is missing or invalid. | RULE-023, WAVE-01 PIN auth |
| AC-W02-021 | Exercise media REST endpoints return 401 when PIN session header is missing or invalid. | RULE-024, WAVE-01 PIN auth |
| AC-W02-022 | File upload rejects files with disallowed MIME types (only JPEG/PNG/WEBP/MP4/MOV/WEBM allowed). Returns validation error with supported types listed. | TDEC-008, AC-045 |
| AC-W02-023 | File upload rejects files exceeding size limits (25MB for images, 250MB for video). Returns validation error. | TDEC-008, RULE-024 |
| AC-W02-024 | Uploaded file name is sanitized: path separators removed, UUID-based storage name used. No path traversal possible. | Security, EDGE-014 |

## Exit Criteria Contributions (Deduplicated)

| EC ID | Description | Evidence/Command |
| --- | --- | --- |
| EC-W02-001 | All AC-W02-* acceptance criteria pass in focused unit and integration tests per TEST-W02-* mapping. | bunx nx run api:test -- --run exercise |
| EC-W02-002 | gqlgen codegen produces valid Go code for Exercise schema without drift. | bunx nx run api:codegen && bunx nx run graphql:codegen |
| EC-W02-003 | sqlc codegen produces valid Go code for exercise queries without drift. | bunx nx run api:codegen |
| EC-W02-004 | WAVE-01 media REST scaffold extended for exercise-media association. Existing media endpoints unchanged. | bunx nx run api:test -- --run media |
| EC-W02-005 | WAVE-01 PIN auth guard protects all WAVE-02 GraphQL and REST endpoints. Existing admin auth unchanged. | TEST-W02-020, TEST-W02-021 |
| EC-W02-006 | WAVE-01 test suite (admin auth, health) still passes after WAVE-02 changes. | bunx nx run api:test |
| EC-W02-007 | Lint passes for all changed packages. | bunx nx run api:lint |
| EC-W02-008 | Migrations 00080 (exercises) and 00081 (exercise_media) apply and roll back in sequence without errors. | TEST-W02-005 |
| EC-W02-009 | File size and type validation enforced for exercise media uploads per TDEC-008. | TEST-W02-022, TEST-W02-023 |
| EC-W02-010 | Exercise and ExerciseMedia operations appear in audit log markers ([Exercise][*], [ExerciseMedia][*]) per TDEC-004 pattern. | TEST-W02-019 (moved from AC) |
| EC-W02-011 | No sensitive content (personalNotes, media file content) appears in application logs. | TEST-W02-019 |
| EC-W02-012 | Exercise round-trip integration test passes: create exercise → upload media → verify media → delete media → soft-delete → verify inactive. | TEST-W02-024 |
| EC-W02-013 | Exercise list without pagination (allExercises) works for WAVE-03 exercise selector dependency. | TEST-W02-007 |

## Decision Log Entries

| ID | Decision | Rationale | Source |
| --- | --- | --- | --- |
| DDEC-W02-001 | Soft delete (isActive=false) instead of hard delete | Preserves referential integrity for WAVE-03 workout history | CAP-W02-005, EDGE-018 |
| DDEC-W02-002 | ExerciseMedia file deletion removes physical file from disk | Consistent with TDEC-005 media deletion principle | TDEC-005, DQ-W02-001 resolved |
| DDEC-W02-003 | Duplicate exercise names allowed | PRD §11 does not specify uniqueness; EDGE-002 | EDGE-002 |
| DDEC-W02-004 | allExercises is a GraphQL query (not REST) | Consistent with hybrid pattern; frontend uses GraphQL | PAGE-002, TDEC-001 |

## Questions Raised (Deduplicated)

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W02-001 | WAVE-02 | data-ops | wave-blocking | EDGE-020 | Should deleting ExerciseMedia delete the physical media file from disk? (Merged from DQ-W02-004) | Orphaned files accumulate | Yes, per TDEC-005. Add failure handling (log error if file deletion fails) | resolved |
| DQ-W02-002 | WAVE-02 | product | needs-owner-decision | AC-043 | Are exercise names unique per user or can duplicates exist? | Affects validation logic | Duplicates allowed (no constraint) — consistent with EDGE-002 | open |
| DQ-W02-003 | WAVE-02 | data-ops | wave-blocking | WAVE-01 | What exact file storage path pattern does WAVE-01 MediaConfig provide for exercise media? | Drives migration and handler design | Use WAVE-01 BasePath/<exercise_id>/<uuid>.<ext> | open |
| DQ-W02-005 | WAVE-02 | security | needs-owner-decision | TDEC-008 | Should WAVE-02 use server-side MIME detection (file magic bytes) or trust Content-Type header? | Content-Type spoofable; magic bytes more secure | Use http.DetectContentType() cross-check | open |
| DQ-W02-006 | WAVE-02 | security | deferred | EDGE-014 | Should exercise media URLs be time-limited (signed URLs)? | Signed URLs add complexity for single-user | Deferred — session-gated access sufficient for MVP | deferred |
| DQ-W02-007 | WAVE-02 | testing | needs-owner-decision | WAVE-01 | Should exercise tests use mocked PIN auth or integration through full middleware chain? | Test complexity vs realism | Prefer integration through full middleware chain | open |
| DQ-W02-008 | WAVE-02 | sequencing | watchlist | WAVE-03 | Does allExercises need filtering beyond isActive? | Might be needed if library grows large | Deferred — current scope is unfiltered active list | watchlist |

## Consolidated Traceability: Source Outcomes → ACs
- OUT-W02-001 (Exercises CRUD) → AC-W02-001 through AC-W02-014
- OUT-W02-002 (Media attached) → AC-W02-015 through AC-W02-018, AC-W02-022 through AC-W02-024
- OUT-W02-003 (Working weight) → AC-W02-006, AC-W02-014
- OUT-W02-004 (API ready for diary) → AC-W02-009, AC-W02-019