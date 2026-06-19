# WAVE-02 testing-exit Planner Attempt 2

## Cycle 1 Reviewer Feedback Addressed

### 1. Round-Trip Integration Test (TEST-W02-024)
Added: Exercise round-trip integration test:
- Create exercise with name, muscleGroups, workingWeight → verify response
- Upload media file to exercise → verify media metadata
- List exercise media → verify media appears
- Download media → verify file content
- Delete media → verify 204, verify media removed from list
- Soft-delete exercise → verify isActive=false
- Query exercise by ID → verify it's still returned
- Default list → verify exercise NOT in list
- List with includeInactive=true → verify exercise IS in list

### 2. EC-W02-001 Strengthened
Changed from "All acceptance criteria passing" to explicit mapping:
```text
| EC-W02-001 | AC-W02-001 through AC-W02-024 pass via TEST-W02-001 through TEST-W02-024 | bunx nx run api:test -- --run exercise |
```

### 3. EC for Migration Order Added
```text
| EC-W02-008 | Migrations 00080 (exercises) and 00081 (exercise_media) apply and roll back in sequence | TEST-W02-005 |
```

### 4. Test Fixture Strategy
Following WAVE-01's test helper pattern:
```go
// Helper function for creating test exercises
func createTestExercise(t *testing.T, repo *ExerciseRepo, name string, opts ...func(*service.Exercise)) *service.Exercise {
    t.Helper()
    // ... creates exercise with unique name
}

// Helper function for creating test media
func createTestExerciseMedia(t *testing.T, repo *ExerciseMediaRepo, exerciseID string) *service.ExerciseMedia {
    t.Helper()
    // ... creates media record with mock file path
}
```

Fixture pattern: factory functions returning typed objects, with sensible defaults and option overrides. Same pattern as WAVE-01 test helpers.

### 5. EC-W02-022 Removed (per testing-exit and sequencing reviewers)
"No WAVE-03 functionality" is a code review gate, not a testable EC. Removed.

### 6. TEST-W02-020 Expanded (Pagination Edge Cases)
Pagination test cases:
- Default page size (no first/after params) — returns default N items
- first=1 returns 1 item
- Cursor from first page applied to second page — returns remaining
- Cursor at end of dataset — returns empty
- Cursor with no items — returns empty

### 7. TEST-W02-019 (Log Sanitize) Moved to EC
AC-W02-019 (from attempt 1) moved to EC-W02-011 (log sanitization). Test retained as TEST-W02-019.

### 8. Auth Test Strategy (DQ-W02-007 Resolution)
Decision: Prefer integration through full middleware chain for auth tests. This means:
- Create a valid PIN session (via test helper from WAVE-01)
- Call exercise endpoints with/without the session
- Verify AuthError for missing session vs normal response for valid session

This is more realistic and end-to-end than mocking the middleware.

## Updated Deduplicated Test Obligations

| Test ID | Description | Type | Command | AC/EC Coverage |
| --- | --- | --- | --- | --- |
| TEST-W02-001 | Exercise repository CRUD unit tests | unit | bunx nx run api:test -- --run '(?i)exercise_repo' | AC-W02-001 through AC-W02-007 |
| TEST-W02-002 | Exercise service validation tests (name required, weight > 0) | unit | bunx nx run api:test -- --run '(?i)exercise_service' | AC-W02-012, AC-W02-014 |
| TEST-W02-003 | Exercise GraphQL resolver integration tests (union results) | integration | bunx nx run api:test -- --run '(?i)exercise_resolver' | AC-W02-001 through AC-W02-014, AC-W02-018, AC-W02-019, AC-W02-020 |
| TEST-W02-004 | ExerciseMedia REST handler integration tests | integration | bunx nx run api:test -- --run '(?i)exercise_media_handler' | AC-W02-015 through AC-W02-017 |
| TEST-W02-005 | Migration smoke test (00080 + 00081 up + down) | integration | bunx nx run api:test -- --run '(?i)migration' | EC-W02-008 |
| TEST-W02-006 | Codegen drift check (gqlgen + sqlc) | codegen | bunx nx run api:codegen && bunx nx run graphql:codegen | EC-W02-002, EC-W02-003 |
| TEST-W02-007 | allExercises query for WAVE-03 dependency | integration | bunx nx run api:test -- --run '(?i)exercise_list_all' | AC-W02-019, EC-W02-013 |
| TEST-W02-008 | Soft delete referential integrity (exercise accessible by ID after soft delete) | integration | bunx nx run api:test -- --run '(?i)exercise_soft_delete' | AC-W02-009 |
| TEST-W02-009 | isActive filter: default list excludes inactive, includeInactive=true includes | integration | bunx nx run api:test -- --run '(?i)exercise_active_filter' | AC-W02-008, AC-W02-010 |
| TEST-W02-010 | Exercise list pagination (cursor behavior, default size, edge cases) | integration | bunx nx run api:test -- --run '(?i)exercise_pagination' | AC-W02-003, AC-W02-004 |
| TEST-W02-011 | Duplicate exercise names allowed (no constraint violation) | integration | bunx nx run api:test -- --run '(?i)exercise_duplicate_name' | AC-W02-013 |
| TEST-W02-012 | Exercise field update persistence (update name → query returns new name) | integration | bunx nx run api:test -- --run '(?i)exercise_update' | AC-W02-005, AC-W02-006 |
| TEST-W02-013 | Exercise GraphQL returns AuthError without valid PIN session | integration | bunx nx run api:test -- --run '(?i)exercise_auth' | AC-W02-020, EC-W02-005 |
| TEST-W02-014 | ExerciseMedia upload returns 401 without valid PIN session | integration | bunx nx run api:test -- --run '(?i)exercise_media_auth' | AC-W02-021, EC-W02-005 |
| TEST-W02-015 | File type validation rejects unauthorized MIME types | unit | bunx nx run api:test -- --run '(?i)exercise_media_filetype' | AC-W02-022 |
| TEST-W02-016 | File size validation rejects oversized uploads | unit | bunx nx run api:test -- --run '(?i)exercise_media_filesize' | AC-W02-023, EC-W02-009 |
| TEST-W02-017 | Path traversal prevention in upload handler | unit | bunx nx run api:test -- --run '(?i)exercise_media_path_traversal' | AC-W02-024 |
| TEST-W02-018 | Log privacy: personalNotes not appearing in logs | unit | bunx nx run api:test -- --run '(?i)exercise_log_sanitize' | EC-W02-011 |
| TEST-W02-019 | Go lint for API package | lint | bunx nx run api:lint | EC-W02-007 |
| TEST-W02-020 | GraphQL schema validate | codegen | bunx nx run graphql:validate | EC-W02-002 |
| TEST-W02-021 | Exercise round-trip integration test (full lifecycle) | integration | bunx nx run api:test -- --run '(?i)exercise_roundtrip' | EC-W02-012 |
| TEST-W02-022 | WAVE-01 admin auth regression test | unit | bunx nx run api:test -- --run '(?i)admin_auth' | EC-W02-006 |

## Updated Exit Criteria

| EC ID | Description | Proof |
| --- | --- | --- |
| EC-W02-001 | AC-W02-001 through AC-W02-024 pass via TEST-W02-001 through TEST-W02-022 | bunx nx run api:test -- --run exercise |
| EC-W02-002 | gqlgen codegen valid | bunx nx run graphql:validate; bunx nx run api:codegen |
| EC-W02-003 | sqlc codegen valid | bunx nx run api:codegen |
| EC-W02-004 | WAVE-01 media scaffold extended for exercise-media | TEST-W02-004 |
| EC-W02-005 | WAVE-01 PIN auth protects all WAVE-02 endpoints | TEST-W02-013, TEST-W02-014 |
| EC-W02-006 | WAVE-01 admin auth/health unchanged | TEST-W02-022 |
| EC-W02-007 | Lint passes | bunx nx run api:lint |
| EC-W02-008 | Migrations apply/rollback in sequence | TEST-W02-005 |
| EC-W02-009 | File validation per TDEC-008 | TEST-W02-015, TEST-W02-016 |
| EC-W02-010 | Audit log markers present | TEST-W02-018 (log markers) |
| EC-W02-011 | No sensitive content in logs | TEST-W02-018 |
| EC-W02-012 | Exercise round-trip integration passes | TEST-W02-021 |
| EC-W02-013 | allExercises query works for WAVE-03 | TEST-W02-007 |