<!-- FILE: docs/prd-wave-details/codebase-fit.md -->
<!-- VERSION: 1.0.0 -->

# Codebase Fit

## Relevant Modules
- M-API: contains all Go backend code — new AiReview model, repository, service, resolvers
- No new config, Docker, or deployment changes needed

## Relevant Files Read
- apps/api/internal/atlas/models/week_flag.go, ai_export.go (model patterns)
- apps/api/internal/atlas/service/week_flag.go (service interface pattern)
- apps/api/internal/atlas/repository/postgres/week_flag_repo.go (repository pattern)
- apps/api/internal/atlas/graph/schema/week_flag.graphql (GraphQL schema pattern)
- apps/api/internal/atlas/graph/resolver/resolver.go (resolver container)
- apps/api/internal/atlas/graph/resolver/week_flag.go (resolver implementation pattern)
- apps/api/internal/repository/postgres/migrations/00092_ai_exports.sql (migration DDL pattern)
- apps/api/cmd/server/main.go (wiring pattern)

## Public Contracts

### New Interface Exports
- AiReviewService: Create(ctx, userID, input), GetByID(ctx, userID, id), ListByUserID(ctx, userID), ListByUserIDAndDateRange(ctx, userID, start, end), Update(ctx, userID, id, input), Delete(ctx, userID, id)
- AiReviewRepository: Create(ctx, userID, input), GetByID(ctx, userID, id), ListByUserID(ctx, userID), ListByUserIDAndDateRange(ctx, userID, start, end), Update(ctx, userID, id, input), Delete(ctx, userID, id)
- AiReview model types: Record, Public, Input types, Result types, ErrorCodes

## Generated Artifact Impact
- sqlc: new ai_reviews.sql query file → generated Go code
- gqlgen: new AiReview type bindings in atlas-gqlgen.yml → generated resolvers
- Both require bun run codegen after creation

## Integration Points
- resolver.go: add `AiReviewService service.AiReviewService` field
- main.go: wire `aiReviewRepo := postgres.NewAiReviewRepository(pool)`, `aiReviewSvc := service.NewAiReviewService(aiReviewRepo)`, add to Resolver
- schema.graphql: add extend type Query and extend type Mutation blocks
- atlas-gqlgen.yml: add AiReview, AiReviewResult, AiReviewsResult, CreateAiReviewInput, UpdateAiReviewInput, AiReviewErrorCode bindings

## Likely Graph Deltas
- M-API: add AiReview service, repository, resolvers sub-modules
- V-M-API: add V-M-API-AI-REVIEW verification entry

## Unsupported Assumptions
- No assumption that external AI integration will exist (explicitly excluded)
- No assumption about multi-user access patterns (single-tenant)
- No assumption about data retention/cleanup policies (user-managed)