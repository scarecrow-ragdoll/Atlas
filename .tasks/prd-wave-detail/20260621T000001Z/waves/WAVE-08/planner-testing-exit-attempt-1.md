<!-- FILE: .tasks/prd-wave-detail/20260621T000001Z/waves/WAVE-08/planner-testing-exit-attempt-1.md -->
<!-- VERSION: 1.0.0 -->

# WAVE-08 Testing-Exit Planner Attempt 1

## Sources Read
- docs/prd-wave-details/waves/wave-07.md (test pattern, test IDs)
- docs/technical-verified/testing-and-delivery.md
- docs/verification-plan.xml (verification references)
- docs/product-verified/acceptance-criteria.md (AC-025, AC-090-092)

## Selected Backend Wave Boundary
Backend CRUD for AiReview. GraphQL-only, no REST endpoints, no file handling.

## Neighboring Backend Wave Fit
WAVE-07 test patterns: Mock repos via testify, INTEGRATION_TESTS=1 guard for repo tests.

## Frontend Pages Context
No frontend test obligations.

## Codebase Evidence
- Existing test files: ai_export_service_test.go, user_profile_service_test.go, week_flag_service_test.go, wave07_migration_test.go
- Pattern: testify, mock repos, httptest for handlers
- INTEGRATION_TESTS=1 guard for repo/integration tests

## Proposed Details

### Exit Criteria

**EC-W08-001 — Migration 00093 applies cleanly**
Migration creates ai_reviews table with correct columns (id, user_id, date_range_start, date_range_end, ai_response_text, user_notes, planned_actions, created_at, updated_at). Down works.

**EC-W08-002 — All sqlc queries compile and regenerate**
bun run codegen succeeds for ai_reviews.sql.

**EC-W08-003 — All Go code compiles**
bun run build succeeds. No type errors, no missing imports.

**EC-W08-004 — GraphQL schema passes gqlgen generation**
gqlgen type bindings in atlas-gqlgen.yml are complete. bun run codegen succeeds.

**EC-W08-005 — AiReview service unit tests pass**
TEST-W08-001 through TEST-W08-009 all pass.

**EC-W08-006 — AiReview resolver unit tests pass**
TEST-W08-010 through TEST-W08-012 all pass.

**EC-W08-007 — Repository integration tests pass**
TEST-W08-013 through TEST-W08-014 all pass (with INTEGRATION_TESTS=1).

**EC-W08-008 — Codegen drift check passes**
bunx nx run api:codegen && bunx nx build api succeeds.

**EC-W08-009 — All endpoints return 401 without valid PIN session**
All GraphQL mutations/queries protected.

**EC-W08-010 — Lint passes**
bun run lint succeeds for all changed packages.

### Verification Obligations

#### AiReview Service Tests
| ID | Description | Command |
|---|---|---|
| TEST-W08-001 | AiReview create succeeds with valid fields | go test -run TestAiReviewService_Create_Success |
| TEST-W08-002 | AiReview create rejects empty aiResponseText | go test -run TestAiReviewService_Create_EmptyText |
| TEST-W08-003 | AiReview create rejects invalid date range (end < start) | go test -run TestAiReviewService_Create_InvalidDateRange |
| TEST-W08-004 | AiReview list returns reviews ordered by createdAt DESC | go test -run TestAiReviewService_List_Ordered |
| TEST-W08-005 | AiReview list filters by date range | go test -run TestAiReviewService_List_DateRangeFilter |
| TEST-W08-006 | AiReview update succeeds | go test -run TestAiReviewService_Update_Success |
| TEST-W08-007 | AiReview update returns not found for wrong user | go test -run TestAiReviewService_Update_Ownership |
| TEST-W08-008 | AiReview delete succeeds | go test -run TestAiReviewService_Delete_Success |
| TEST-W08-009 | AiReview log privacy — no content in logs | go test -run TestAiReviewService_Logs_NoContent |

#### AiReview Resolver Tests
| ID | Description | Command |
|---|---|---|
| TEST-W08-010 | AiReview create resolver | go test ./internal/atlas/graph/resolver/ -run TestAiReviewResolver_Create |
| TEST-W08-011 | AiReview list resolver | go test ./internal/atlas/graph/resolver/ -run TestAiReviewResolver_List |
| TEST-W08-012 | AiReview delete resolver | go test ./internal/atlas/graph/resolver/ -run TestAiReviewResolver_Delete |

#### Repository Integration Tests (requires INTEGRATION_TESTS=1)
| ID | Description | Command |
|---|---|---|
| TEST-W08-013 | AiReview repository operations | INTEGRATION_TESTS=1 go test -run TestAiReviewRepo |
| TEST-W08-014 | Migration 00093 applies cleanly | INTEGRATION_TESTS=1 go test -run TestWave08Migration |

## Questions Raised
Q-W08-TST-001: Should we include an updateAiReview mutation resolver test? Yes, included in TEST-W08-006/007 (service level covers all update logic).

## Traceability Candidates
- EC-W08-001 through EC-W08-010 → standard exit criteria matching WAVE-07 pattern
- TEST-W08-001→009 → AC-W08-001 through AC-W08-008
- TEST-W08-010→012 → GraphQL resolver coverage
- TEST-W08-013→014 → integration coverage