# WAVE-02 testing-exit Planner Attempt 1

## Sources Read
- docs/technical-verified/testing-and-delivery.md
- docs/technical-verified/implementation-slices.md
- docs/prd-wave-details/waves/wave-01.md (verification obligations)
- apps/api/internal/repository/postgres/user_repo_test.go
- apps/api/internal/graph/schema_resolvers_test.go
- apps/api/internal/handler/users_test.go
- apps/api/internal/handler/users_internal_test.go

## Selected Backend Wave Boundary
WAVE-02 test coverage: Exercise CRUD (GraphQL service + resolver + repository), ExerciseMedia upload/delete (REST handler), migration tests, auth integration, file validation.

## Neighboring Backend Wave Fit
- WAVE-01 test infrastructure: provides test DB migrations, test helpers for PIN auth, test patterns. WAVE-02 adds exercise-related tests.
- WAVE-03 tests will depend on exercise fixtures created here.

## Frontend Pages Context
No frontend tests in WAVE-02. PAGE-003 and PAGE-002 frontend tests are separate scope.

## Codebase Evidence
- Existing test patterns: internal tests (handler/users_internal_test.go), integration tests (handler/users_test.go), resolver tests (graph/schema_resolvers_test.go), repo tests (postgres/user_repo_test.go, user_repo_unit_test.go)
- Go test framework: testing + testify/require
- Existing test helpers: postgres_test.go (DB setup), handler patterns with httptest

## Exit Criteria Contributions

| EC ID | Description |
| --- | --- |
| EC-W02-001 | All acceptance criteria passing in focused unit and integration tests |
| EC-W02-002 | gqlgen codegen produces valid Go code for Exercise schema (no drift) |
| EC-W02-003 | sqlc codegen produces valid Go code for exercise queries (no drift) |
| EC-W02-007 | Lint passes for all changed packages |
| EC-W02-015 | Exercise service tests cover all CRUD operations (create, read, update, soft delete) |
| EC-W02-016 | ExerciseMedia handler tests cover upload success, file type rejection, file size rejection, delete |
| EC-W02-017 | Exercise resolver tests verify union result types (Success, ValidationError, AuthError, NotFoundError) |
| EC-W02-018 | Exercise list tests verify pagination and isActive filtering |
| EC-W02-019 | Test fixtures create deterministic exercises with known IDs for WAVE-03 consumption |

## Verification Obligations

| Test ID | Description | Type | Command | Coverage |
| --- | --- | --- | --- | --- |
| TEST-W02-001 | Exercise repository CRUD unit tests (create, get, list, update, soft delete) | unit | bunx nx run api:test -- --run '(?i)exercise_repo' | ExerciseRepo |
| TEST-W02-002 | Exercise service unit tests (validation, soft delete logic, list filtering) | unit | bunx nx run api:test -- --run '(?i)exercise_service' | ExerciseService |
| TEST-W02-003 | Exercise GraphQL resolver integration tests (union results, auth errors) | integration | bunx nx run api:test -- --run '(?i)exercise_resolver' | Exercise resolver |
| TEST-W02-004 | ExerciseMedia REST handler integration tests (upload, download, delete) | integration | bunx nx run api:test -- --run '(?i)exercise_media_handler' | ExerciseMedia handler |
| TEST-W02-005 | Migration smoke test (00080 + 00081 up + down) | integration | bunx nx run api:test -- --run '(?i)migration' | Migrations |
| TEST-W02-006 | Codegen drift check (gqlgen + sqlc) | codegen | bunx nx run api:codegen && bunx nx run graphql:codegen | Codegen |
| TEST-W02-007 | Exercise list all (simple list for WAVE-03 dependency) | integration | bunx nx run api:test -- --run '(?i)exercise_list_all' | ListAll query |
| TEST-W02-008 | Soft delete referential integrity (exercise accessible by ID after soft delete) | integration | bunx nx run api:test -- --run '(?i)exercise_soft_delete' | Soft delete |
| TEST-W02-013 | Soft-deleted exercise excluded from default list, included with includeInactive=true | integration | bunx nx run api:test -- --run '(?i)exercise_active_filter' | Active filter |
| TEST-W02-014 | Exercise GraphQL returns AuthError without valid PIN session | integration | bunx nx run api:test -- --run '(?i)exercise_auth' | PIN auth |
| TEST-W02-015 | ExerciseMedia upload returns 401 without valid PIN session | integration | bunx nx run api:test -- --run '(?i)exercise_media_auth' | PIN auth media |
| TEST-W02-016 | File type rejection for unauthorized MIME types | unit | bunx nx run api:test -- --run '(?i)exercise_media_filetype' | File validation |
| TEST-W02-017 | File size rejection for oversized uploads | unit | bunx nx run api:test -- --run '(?i)exercise_media_filesize' | Size validation |
| TEST-W02-018 | Path traversal prevention in upload handler | unit | bunx nx run api:test -- --run '(?i)exercise_media_path_traversal' | Security |
| TEST-W02-019 | Sensitive content not appearing in log output | unit | bunx nx run api:test -- --run '(?i)exercise_log_sanitize' | Log privacy |
| TEST-W02-020 | Exercise list pagination (cursor-based, default page size) | integration | bunx nx run api:test -- --run '(?i)exercise_pagination' | Pagination |
| TEST-W02-021 | Exercise name validation (required, non-empty) | unit | bunx nx run api:test -- --run '(?i)exercise_validation' | Validation |
| TEST-W02-022 | Go lint for API package | lint | bunx nx run api:lint | Lint |
| TEST-W02-023 | GraphQL schema validate | codegen | bunx nx run graphql:validate | Schema |

## Test Fixtures
- Exercise fixtures: create exercises with known name, muscleGroups, workingWeight values
- ExerciseMedia fixtures: create media records associated with known exercises, mock file system
- Follow WAVE-01 test helper pattern for DB and file system setup

## Risks And Rollback
- Risk: tests that use real file uploads need temporary directories and cleanup. Mitigation: use t.TempDir() for test file storage.
- Risk: tests depending on WAVE-01 PIN auth middleware may be fragile. Mitigation: mock PIN auth middleware or test through the middleware chain with a valid session.

## Questions Raised

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Source Or Report | Status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W02-007 | WAVE-02 | testing-exit | needs-owner-decision | WAVE-01 test patterns | Should exercise tests use the WAVE-01 PIN test helper (mocked session) or integration test through full middleware chain? | Affects test complexity and CI run time. Mocked is faster, full chain is more realistic. | testing-exit planner | open |

## Traceability Candidates
- All test IDs → docs/technical-verified/testing-and-delivery.md
- Migration tests → docs/technical-verified/data-contracts.md
- Auth tests → docs/technical-verified/auth-security-compliance.md
- File validation tests → docs/technical-verified/integrations-and-events.md (TDEC-008)