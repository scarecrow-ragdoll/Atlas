<!-- FILE: .tasks/prd-wave-detail/20260621T000001Z/waves/WAVE-08/planner-architecture-codebase-attempt-1.md -->
<!-- VERSION: 1.0.0 -->

# WAVE-08 Architecture-Codebase Planner Attempt 1

## Sources Read
- docs/prd-wave-details/waves/wave-07.md (template for patterns)
- apps/api/internal/atlas/models/ai_export.go (model pattern)
- apps/api/internal/atlas/models/week_flag.go (model pattern — Record/Result/ErrorCode)
- apps/api/internal/atlas/service/week_flag.go (service interface pattern)
- apps/api/internal/atlas/repository/postgres/week_flag_repo.go (repository pattern)
- apps/api/internal/atlas/graph/schema/week_flag.graphql (GraphQL schema pattern)
- apps/api/internal/atlas/graph/resolver/resolver.go (resolver container)
- apps/api/internal/repository/postgres/migrations/00092_ai_exports.sql (migration DDL pattern)
- apps/api/internal/atlas/models/settings.go (reference)

## Selected Backend Wave Boundary
WAVE-08: backend-only CRUD for AiReview entity. No file storage, no external integrations, no async processing.

## Neighboring Backend Wave Fit
Prior WAVE-07 creates AiExport + UserProfile. WAVE-08 creates independent AiReview with the same pattern triples (model/repo/service/resolver). No scope collision.

## Frontend Pages Context
No dedicated frontend page for AiReview. PAGE-009 may reference review history. Backend provides GraphQL queries only.

## Codebase Evidence

### Existing Patterns to Follow (WAVE-07/WeekFlag)
1. **Models**: AiReviewRecord (DB), AiReview (public JSON), CreateAiReviewInput, UpdateAiReviewInput, AiReviewResult, AiReviewsResult, AiReviewErrorCode, AiReviewFromRecord converter
2. **Repository**: Interface + private struct + New*Repository(pool) + *FromRow() helpers
3. **Service**: Interface + private struct + constructor + sentinel errors + FromRecord()
4. **Resolvers**: middleware.GetAtlasUserID + union result types
5. **GraphQL schema**: Types with nullable fields + input types + result types with inline errors
6. **Migrations**: goose +pgx pattern with +goose Up/+goose Down

### Relevant Existing Files
- apps/api/internal/repository/postgres/migrations/ — migration directory (next: 00093)
- apps/api/internal/repository/postgres/queries/ — sqlc queries directory
- apps/api/internal/atlas/models/ — model types directory
- apps/api/internal/atlas/repository/postgres/ — repository implementations
- apps/api/internal/atlas/service/ — service implementations
- apps/api/internal/atlas/graph/schema/ — GraphQL schema files
- apps/api/internal/atlas/graph/resolver/ — GraphQL resolvers
- apps/api/internal/atlas/graph/resolver/resolver.go — resolver container
- apps/api/internal/atlas/graph/schema/schema.graphql — main schema with type extensions
- apps/api/cmd/server/main.go — wiring
- apps/api/atlas-gqlgen.yml — gqlgen type bindings
- apps/api/sqlc.yaml — sqlc config

## Public Contracts

### New Exports
- AiReviewService (interface): Create, GetByID, ListByUserID, ListByUserIDAndDateRange, Update, Delete
- AiReviewRepository (interface): Create, GetByID, ListByUserID, ListByUserIDAndDateRange, Update, Delete
- AiReview model types: Record, Public, Inputs, Results, ErrorCodes

## Generated Artifact Impact
- sqlc: new ai_reviews.sql queries → generated Go code
- gqlgen: new AiReview type bindings in atlas-gqlgen.yml → generated resolvers

## Integration Points
- Resolver container: add AiReviewService field
- main.go: wire AiReviewRepository → AiReviewService → Resolver
- schema.graphql: add #extend type Query/Mutation for aiReviews

## Likely Graph Deltas
- M-API: add AiReview service, repository, resolvers
- V-M-API: add AiReview verification entries

## Implementation Slices

### SLICE-W08-001: AiReview Migration (00093_ai_reviews.sql)
CREATE TABLE ai_reviews (id UUID PK, user_id UUID FK, date_range_start DATE NOT NULL, date_range_end DATE NOT NULL, ai_response_text TEXT NOT NULL, user_notes TEXT, planned_actions TEXT, created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ). Down: DROP TABLE.

### SLICE-W08-002: AiReview Model (models/ai_review.go)
AiReviewRecord (DB row), AiReview (public JSON), CreateAiReviewInput, UpdateAiReviewInput, AiReviewResult, AiReviewsResult, error types matching AiExport pattern. All nullable fields use pointer types.

### SLICE-W08-003: AiReview SQLc Queries (queries/ai_reviews.sql)
CreateAiReview, GetAiReviewByID (user-scoped), ListAiReviewsByUserID, ListAiReviewsByUserIDAndDateRange, UpdateAiReview, DeleteAiReview.

### SLICE-W08-004: AiReview Repository (repository/postgres/ai_review_repo.go)
Create, GetByID, ListByUserID, ListByUserIDAndDateRange, Update, Delete. Interface + private struct pattern. Uses generated sqlc code + pgtype helpers.

### SLICE-W08-005: AiReview Service (service/ai_review_service.go)
Create(ctx, userID, input), GetByID(ctx, userID, id), List(ctx, userID), ListByDateRange(ctx, userID, start, end), Update(ctx, userID, id, input), Delete(ctx, userID, id). Interface + private struct. Sentinel errors: ErrAiReviewNotFound, ErrAiReviewInvalidDateRange. Validates dateRangeEnd >= dateRangeStart.

### SLICE-W08-006: AiReview GraphQL Schema + Resolver
schema/ai_review.graphql: AiReview type, CreateAiReviewInput, UpdateAiReviewInput, AiReviewResult, AiReviewsResult with inline errors, AiReviewErrorCode enum. Mutations: createAiReview(input), updateAiReview(id, input), deleteAiReview(id). Queries: aiReview(id), aiReviews(dateRangeStart, dateRangeEnd).
resolver/ai_review.go: CreateAiReview, UpdateAiReview, DeleteAiReview, AiReview, AiReviews. Auth guard through middleware.GetAtlasUserID.

### SLICE-W08-007: Main Wiring (main.go, resolver.go, gqlgen)
Wire AiReviewRepository → AiReviewService → Resolver. Add AiReviewService to Resolver struct. Add gqlgen bindings in atlas-gqlgen.yml. Add mutations/queries to schema.graphql.

## Unsupported Assumptions
- No assumption that ui content must be server-validated beyond required fields
- No encryption needed for ai_response_text (manual user entry, single-tenant)

## Questions Raised
Q-W08-ARC-001: Should planned_actions be TEXT or structured? Recommendation: TEXT for MVP.
Q-W08-ARC-002: Does AiReview need an index on (user_id, date_range_start, date_range_end)? Recommendation: yes for query performance on review history filter.

## Traceability Candidates
- SLICE-W08-001 → CAP-W08-001 (AiReview CRUD)
- SLICE-W08-002→006 → CAP-W08-001 through CAP-W08-005
- SLICE-W08-007 → main.go wiring pattern from WAVE-07