<!-- FILE: docs/prd-wave-details/waves/wave-08.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Define the detailed WAVE-08 (AI Review History) backend wave brief for developer handoff. -->
<!--   SCOPE: Covers all backend implementation slices, ACs, ECs, verification obligations, codebase fit, data/API/security design, open questions, and traceability for AiReview CRUD. -->
<!--   DEPENDS: docs/prd-waves/waves/wave-08.md, docs/product-verified/domain-model.md, docs/product-verified/functional-spec.md §19, docs/product-verified/acceptance-criteria.md, docs/prd-wave-details/waves/wave-07.md. -->
<!--   LINKS: M-API / V-M-API / WAVE-08. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->

# Wave 08: AI Review History

## Status
ready-for-dev

## User Approval
Source wave: user-approved (2026-06-18). Detailed wave: ready-for-dev — all questions resolved and user-approved (2026-06-21).

## Source Wave Summary
Store AI responses with period linkage and planned actions. Simple CRUD for AiReview entity. Backend-only wave — no AI call, no OpenAI integration, no file storage. Risk: Low.

## Outcome After Implementation
- Save AI response text (manual paste)
- Link review to a date range (dateRangeStart, dateRangeEnd)
- Add user notes (optional)
- Track planned actions (optional TEXT, MVP)
- Queryable review history with date range filtering

## Scope Included
- AiReview CRUD (Create, Read, List, Update, Delete) via GraphQL
- Period linkage with date range validation (end >= start)
- User notes (optional TEXT field)
- Planned actions storage (optional TEXT field, MVP)
- Review history queries (list by user, filter by date range, ordered by createdAt DESC)
- Service layer for WAVE-09 backup consumption (ListAllByUserID)

## Scope Excluded
- Automatic AI call (explicitly excluded)
- OpenAI integration (explicitly excluded)
- Frontend pages, UI components, routes
- File storage for review data
- Structured planned actions table (deferred to post-MVP)

## Dependencies And Other-Wave Fit

### Prior Wave Compatibility
- **WAVE-01 (Foundation)** — Required. PIN auth middleware for all GraphQL operations. atlas_users for user identity.
- **WAVE-02 through WAVE-06** — Compatible. No dependency.
- **WAVE-07 (AI Export)** — Compatible. Clean boundary. WAVE-07 creates AiExport, WAVE-08 creates independent AiReview. No shared tables or services.

### Future Wave Compatibility
- **WAVE-09 (Backup Import/Export)** — WAVE-08 exposes AiReviewService.ListAllByUserID(ctx, userID) for WAVE-09 to include AiReview data in backup data.json. Read-only interface.

### Independent Deliverability
- Cannot be implemented without WAVE-01 (hard dependency on PIN auth and user identity)
- Can be implemented without WAVE-02 through WAVE-07 (no code dependencies)

## Frontend Pages Dependencies
No dedicated AiReview frontend page exists. PAGE-009 (AI Export) does not list AiReview as a backend dependency. Backend provides GraphQL queries as dependency context only:

- `Query aiReview(id: ID!): AiReviewResult!`
- `Query aiReviews(dateRangeStart: Date, dateRangeEnd: Date): AiReviewsResult!`

## Codebase Fit And Touchpoints

### New Files Required
- `apps/api/internal/repository/postgres/migrations/00093_ai_reviews.sql` — AiReview table
- `apps/api/internal/repository/postgres/queries/ai_reviews.sql` — sqlc queries
- `apps/api/internal/atlas/models/ai_review.go` — AiReview model types
- `apps/api/internal/atlas/repository/postgres/ai_review_repo.go` — AiReview repository
- `apps/api/internal/atlas/service/ai_review_service.go` — AiReview service
- `apps/api/internal/atlas/graph/schema/ai_review.graphql` — AiReview GraphQL schema
- `apps/api/internal/atlas/graph/resolver/ai_review.go` — AiReview resolvers

### Existing Files to Modify
- `apps/api/internal/atlas/graph/resolver/resolver.go` — add AiReviewService field
- `apps/api/internal/atlas/graph/schema/schema.graphql` — add type extensions for Query and Mutation
- `apps/api/cmd/server/main.go` — wire new repository, service, resolvers
- `apps/api/atlas-gqlgen.yml` — add bindings for all new types

### Patterns to Follow
- **Models**: AiReviewRecord/AiReview/CreateAiReviewInput/UpdateAiReviewInput triple matching ai_export.go
- **Repository**: Interface + private struct + New*Repository(pool) + *FromRow() matching week_flag_repo.go
- **Service**: Interface + private struct + constructor + sentinel errors + FromRecord() matching week_flag.go
- **Resolvers**: middleware.GetAtlasUserID + union result types matching resolver/week_flag.go
- **GraphQL schema**: Types with nullable fields + input types + result types with inline errors matching week_flag.graphql

## Design Contracts

### DDEC-W08-001: GraphQL for AiReview CRUD
All AiReview operations use GraphQL (no REST endpoints). Existing gqlgen setup supports CRUD patterns. No file download or binary operations needed — REST is unnecessary.

### DDEC-W08-002: Migration Number 00093
Latest migration is 00092_ai_exports.sql. Next sequential migration is 00093_ai_reviews.sql.

### DDEC-W08-003: planned_actions as TEXT
planned_actions stored as optional TEXT column. MVP simplicity — no separate table. Structured planned actions deferred to post-MVP.

### DDEC-W08-004: GraphQL-Only API
No REST endpoints for AiReview. All operations through gqlgen mutations/queries. WAVE-07 uses REST only for ZIP file download — not applicable to WAVE-08.

### DDEC-W08-005: User-Scoped Queries
All AiReview data access filtered by userId from session context (middleware.GetAtlasUserID). No admin/read-all operations in MVP.

### DDEC-W08-006: Date Range Index
Composite index on (user_id, date_range_start, date_range_end) recommended for review history filter queries. Not required for MVP correctness but good practice.

## Data API Integration And Operations

### GraphQL Operations

**Mutation createAiReview(input: CreateAiReviewInput!): AiReviewResult!**
- Input: dateRangeStart (Date!), dateRangeEnd (Date!), aiResponseText (String!), userNotes (String), plannedActions (String)
- Validation: aiResponseText non-empty, dateRangeEnd >= dateRangeStart
- Success: AiReview with id, timestamps
- Error: VALIDATION_ERROR

**Mutation updateAiReview(id: ID!, input: UpdateAiReviewInput!): AiReviewResult!**
- Input: dateRangeStart (Date), dateRangeEnd (Date), aiResponseText (String), userNotes (String), plannedActions (String)
- Partial update: only provided fields updated
- Success: Updated AiReview
- Error: NOT_FOUND, VALIDATION_ERROR, AUTH_ERROR

**Mutation deleteAiReview(id: ID!): AiReviewResult!**
- Success: Deleted AiReview
- Error: NOT_FOUND, AUTH_ERROR

**Query aiReview(id: ID!): AiReviewResult!**
- User-scoped: filters by session user
- Error: NOT_FOUND, AUTH_ERROR

**Query aiReviews(dateRangeStart: Date, dateRangeEnd: Date): AiReviewsResult!**
- All user's reviews when no filters; filtered by date range when provided
- Ordered by createdAt DESC
- Error: AUTH_ERROR

### Data Lifecycle
- Created: user pastes AI response text, links to date range, optionally adds notes/planned actions
- Read: user views individual review or review history
- Updated: user modifies text, date range, notes, or planned actions
- Deleted: user removes review entry
- No TTL, no cleanup: user-managed lifecycle

### Validation Rules
- aiResponseText: required, non-empty (service layer)
- dateRangeStart, dateRangeEnd: required on create (GraphQL input non-null)
- dateRangeEnd >= dateRangeStart (service layer validation)
- All text fields: no max length validation in MVP (PostgreSQL TEXT up to 1GB)

### Rollout Order
1. Create migration 00093_ai_reviews.sql
2. Add sqlc queries in ai_reviews.sql, run bun run codegen
3. Create AiReview model types (ai_review.go)
4. Create AiReview repository (ai_review_repo.go)
5. Create AiReview service (ai_review_service.go)
6. Create GraphQL schema (ai_review.graphql)
7. Create GraphQL resolvers (ai_review.go)
8. Wire everything: resolver.go, main.go, atlas-gqlgen.yml, schema.graphql
9. Run migrations, deploy, verify with tests

### Rollback
1. goose down 00093
2. Remove all new source files
3. Revert main.go, resolver.go, atlas-gqlgen.yml, schema.graphql
4. bun run codegen to purge generated code

### Observability
- Log markers: [AiReview][create], [AiReview][update], [AiReview][delete], [AiReview][list]
- No ai_response_text, user_notes, or planned_actions content in application logs
- Only metadata logged: review ID, date range, operation type
- Following AC-118 privacy pattern from WAVE-07

## Security Privacy And Compliance
- All GraphQL operations under PIN-guard middleware (consistent with WAVE-01 auth pattern)
- User-scoped queries: all data access filtered by userId from session context
- Ownership validation on GetByID: verify AiReview.userId matches session user
- Log privacy: no review content in logs (AC-118 pattern)
- Data stored in PostgreSQL (encrypted at rest via volume encryption)
- No external data transmission (manual user entry only)
- No file storage, no media upload risks

## Implementation Slices

### SLICE-W08-001: AiReview Migration (00093_ai_reviews.sql)
```sql
-- +goose Up
CREATE TABLE ai_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES atlas_users(id) ON DELETE CASCADE,
    date_range_start DATE NOT NULL,
    date_range_end DATE NOT NULL,
    ai_response_text TEXT NOT NULL,
    user_notes TEXT,
    planned_actions TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_ai_reviews_user_id_date_range ON ai_reviews(user_id, date_range_start, date_range_end);
-- +goose Down
DROP TABLE IF EXISTS ai_reviews;
```

### SLICE-W08-002: AiReview Model (models/ai_review.go)
AiReviewRecord (DB row), AiReview (public JSON), CreateAiReviewInput, UpdateAiReviewInput, AiReviewResult, AiReviewsResult, error types matching AiExport pattern. All nullable fields use pointer types.

### SLICE-W08-003: AiReview SQLc Queries (queries/ai_reviews.sql)
CreateAiReview, GetAiReviewByID (user-scoped), ListAiReviewsByUserID, ListAiReviewsByUserIDAndDateRange, UpdateAiReview, DeleteAiReview.

### SLICE-W08-004: AiReview Repository (repository/postgres/ai_review_repo.go)
Create, GetByID, ListByUserID, ListByUserIDAndDateRange, Update, Delete. Interface + private struct pattern. Uses generated sqlc code + pgtype helpers.

### SLICE-W08-005: AiReview Service (service/ai_review_service.go)
Service interface: Create(ctx, userID, input), GetByID(ctx, userID, id), ListByUserID(ctx, userID), ListByUserIDAndDateRange(ctx, userID, start, end), Update(ctx, userID, id, input), Delete(ctx, userID, id). ListAllByUserID(ctx, userID) for WAVE-09 backup consumption.
Sentinel errors: ErrAiReviewNotFound, ErrAiReviewInvalidDateRange.
Validation: dateRangeEnd >= dateRangeStart.

### SLICE-W08-006: AiReview GraphQL Schema + Resolver
schema/ai_review.graphql: AiReview type, CreateAiReviewInput, UpdateAiReviewInput, AiReviewResult, AiReviewsResult with inline errors, AiReviewErrorCode enum.
Mutations: createAiReview(input), updateAiReview(id, input), deleteAiReview(id).
Queries: aiReview(id), aiReviews(dateRangeStart, dateRangeEnd).
resolver/ai_review.go: CreateAiReview, UpdateAiReview, DeleteAiReview, AiReview, AiReviews. Auth guard through middleware.GetAtlasUserID.

### SLICE-W08-007: Main Wiring (main.go, resolver.go, gqlgen)
Wire AiReviewRepository → AiReviewService → Resolver. Add AiReviewService to Resolver struct. Add gqlgen bindings in atlas-gqlgen.yml. Add mutations/queries to schema.graphql.

## Acceptance Criteria

| ID | Criterion | Source |
|----|-----------|--------|
| AC-W08-001 | User can create an AI review entry. createAiReview mutation accepts aiResponseText (required), dateRangeStart (required), dateRangeEnd (required), userNotes (optional), plannedActions (optional). Returns created review. | AC-025, §19 |
| AC-W08-002 | User can paste AI response text. aiResponseText accepts arbitrary text. No content validation. | AC-090, §19.2 |
| AC-W08-003 | User can link review to a date range. dateRangeStart and dateRangeEnd are required DATE fields. Validation rejects end < start. | AC-091, §19.2 |
| AC-W08-004 | User can add notes and planned actions. userNotes (optional TEXT) and plannedActions (optional TEXT) accepted on create and update. | AC-092, §19.2 |
| AC-W08-005 | User can view review history. aiReviews query returns all reviews for the authenticated user, ordered by createdAt DESC. | W08-005, §19 |
| AC-W08-006 | User can filter reviews by date range. aiReviews accepts optional dateRangeStart and dateRangeEnd filter parameters. | §19 "review history view" |
| AC-W08-007 | User can update a review entry. updateAiReview mutation accepts partial updates for all fields. | CAP-W08-001 |
| AC-W08-008 | User can delete a review entry. deleteAiReview mutation deletes by id, user-scoped. | CAP-W08-001 |

## Exit Criteria

| ID | Criterion | Validation |
|----|-----------|------------|
| EC-W08-001 | Migration 00093 applies cleanly | Migration creates ai_reviews table with correct columns. Down works. |
| EC-W08-002 | All sqlc queries compile and regenerate | bun run codegen succeeds for ai_reviews.sql |
| EC-W08-003 | All Go code compiles | bun run build succeeds. No type errors. |
| EC-W08-004 | GraphQL schema passes gqlgen generation | gqlgen type bindings complete. bun run codegen succeeds. |
| EC-W08-005 | AiReview service unit tests pass | TEST-W08-001 through TEST-W08-009 all pass |
| EC-W08-006 | AiReview resolver unit tests pass | TEST-W08-010 through TEST-W08-012 all pass |
| EC-W08-007 | Repository integration tests pass | TEST-W08-013 through TEST-W08-014 all pass (INTEGRATION_TESTS=1) |
| EC-W08-008 | Codegen drift check passes | bunx nx run api:codegen && bunx nx build api succeeds |
| EC-W08-009 | All endpoints return 401 without valid PIN session | All GraphQL mutations/queries protected |
| EC-W08-010 | Lint passes | bun run lint succeeds for all changed packages |

## Verification Obligations

### AiReview Service Tests
| ID | Description | Command |
|----|-------------|---------|
| TEST-W08-001 | AiReview create succeeds with valid fields | go test -run TestAiReviewService_Create_Success |
| TEST-W08-002 | AiReview create rejects empty aiResponseText | go test -run TestAiReviewService_Create_EmptyText |
| TEST-W08-003 | AiReview create rejects invalid date range (end < start) | go test -run TestAiReviewService_Create_InvalidDateRange |
| TEST-W08-004 | AiReview list returns reviews ordered by createdAt DESC | go test -run TestAiReviewService_List_Ordered |
| TEST-W08-005 | AiReview list filters by date range | go test -run TestAiReviewService_List_DateRangeFilter |
| TEST-W08-006 | AiReview update succeeds | go test -run TestAiReviewService_Update_Success |
| TEST-W08-007 | AiReview update returns not found for wrong user | go test -run TestAiReviewService_Update_Ownership |
| TEST-W08-008 | AiReview delete succeeds | go test -run TestAiReviewService_Delete_Success |
| TEST-W08-009 | AiReview log privacy — no content in logs | go test -run TestAiReviewService_Logs_NoContent |

### AiReview Resolver Tests
| ID | Description | Command |
|----|-------------|---------|
| TEST-W08-010 | AiReview create resolver | go test ./internal/atlas/graph/resolver/ -run TestAiReviewResolver_Create |
| TEST-W08-011 | AiReview list resolver | go test ./internal/atlas/graph/resolver/ -run TestAiReviewResolver_List |
| TEST-W08-012 | AiReview delete resolver | go test ./internal/atlas/graph/resolver/ -run TestAiReviewResolver_Delete |

### Repository Integration Tests (requires INTEGRATION_TESTS=1)
| ID | Description | Command |
|----|-------------|---------|
| TEST-W08-013 | AiReview repository operations | INTEGRATION_TESTS=1 go test -run TestAiReviewRepo |
| TEST-W08-014 | Migration 00093 applies cleanly | INTEGRATION_TESTS=1 go test -run TestWave08Migration |

## Rollout Rollback And Compatibility

### Rollout Order
1. Create migration 00093_ai_reviews.sql
2. Add sqlc queries, run bun run codegen
3. Create AiReview model types
4. Create AiReview repository
5. Create AiReview service
6. Create GraphQL schema + resolvers
7. Wire everything in main.go, resolver.go, atlas-gqlgen.yml, schema.graphql
8. Run migrations, deploy, verify with tests

### Rollback
1. goose down 00093
2. Remove all new source files
3. Revert main.go/resolver.go wiring
4. Revert atlas-gqlgen.yml (remove bindings)
5. Revert schema.graphql (remove type extensions)
6. bun run codegen to purge generated code

### Compatibility
- Additive changes only — no existing tables modified
- No impact on existing API endpoints or resolvers
- No existing services or data affected
- AiReview index creation is safe — no table lock concern on single-user DB

## Handoff Packets

### Developer Handoff
Dependencies: WAVE-01 mandatory (PIN auth, atlas_users). WAVE-02 through WAVE-07 optional.
Implementation order: Follow SLICE dependency graph — SLICES 001-006 sequential, SLICE-007 (wiring) last.
Key patterns: ai_export.go and week_flag.go for models/repo/service/resolver patterns.
Generators: sqlc (ai_reviews.sql), gqlgen (atlas-gqlgen.yml bindings).
Testing: Mock repos via testify, INTEGRATION_TESTS=1 guard for repo tests.

### Reviewer Handoff
All 7 reviewer perspectives approved on attempt 1. 2 open questions (DQ-W08-001, DQ-W08-002) resolved and user-approved (2026-06-21). Ready for user approval.

## Reviewer Verdicts

| WAVE-08 | product-scope-and-ac | 1 | approved | review-product-scope-and-ac-attempt-1.md | none | All outcomes covered, 8 ACs traceable |
| WAVE-08 | architecture-codebase-fit | 1 | approved | review-architecture-codebase-fit-attempt-1.md | none | 7 slices follow WAVE-07 patterns |
| WAVE-08 | data-api-integration-ops | 1 | approved | review-data-api-integration-ops-attempt-1.md | none | GraphQL-only correct, lifecycle complete |
| WAVE-08 | security-privacy-compliance | 1 | approved | review-security-privacy-compliance-attempt-1.md | none | PIN guard, user-scoping, log privacy covered |
| WAVE-08 | testing-exit-criteria | 1 | approved | review-testing-exit-criteria-attempt-1.md | none | 12 tests across service/resolver/integration |
| WAVE-08 | sequencing-other-wave-fit | 1 | approved | review-sequencing-other-wave-fit-attempt-1.md | none | Clean boundaries, WAVE-09 contract defined |
| WAVE-08 | traceability-consistency | 1 | approved | review-traceability-consistency-attempt-1.md | none | Full source traceability, consistent IDs |
| WAVE-08 | final-wave-fit-review | 1 | approved | final-wave-fit-review-attempt-1.md | none | Package complete, one-wave focus, ready for user approval |

## Open Questions

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Needed Answer | Source Or Report | Status | Resolution |
|----|------|-------|----------|--------|----------|---------------|--------------|------------------|--------|-----------|
| DQ-W08-001 | WAVE-08 | data-api-integration-ops | needs-owner-decision | None | Should planned_actions be a simple TEXT field (MVP) or a structured child table? | PRD says "planned actions storage" — structured enables queryable action tacking; simple TEXT matches MVP constraints | Confirm: simple TEXT for MVP, structured in post-MVP | planner-product-ac-attempt-1 | resolved | Simple TEXT for MVP (user-approved 2026-06-21) |
| DQ-W08-002 | WAVE-08 | sequencing-fit | needs-owner-decision | None | Should WAVE-08 expose ListAllByUserID for WAVE-09 backup consumption? | WAVE-07 context states "WAVE-08 must provide service layer for WAVE-09 to include AiReview data in backups" | Confirm: yes, expose ListAllByUserID | planner-sequencing-fit-attempt-1, wave-07.md | resolved | Yes, expose ListAllByUserID (user-approved 2026-06-21) |

## Traceability
- docs/prd-waves/waves/wave-08.md (source wave)
- docs/product-verified/functional-spec.md §19 (REQ-015)
- docs/product-verified/domain-model.md#AiReview
- docs/product-verified/acceptance-criteria.md (AC-025, AC-090-092)
- docs/product-verified/features/ai-review-history.md
- docs/prd-wave-details/waves/wave-01.md (PIN auth)
- docs/prd-wave-details/waves/wave-07.md (template patterns, boundary)
- docs/prd-wave-details/appendix/traceability.md