<!-- FILE: .tasks/prd-wave-detail/20260621T000001Z/waves/WAVE-08/planner-data-integration-ops-attempt-1.md -->
<!-- VERSION: 1.0.0 -->

# WAVE-08 Data-Integration-Ops Planner Attempt 1

## Sources Read
- docs/product-verified/domain-model.md (AiReview entity, lifecycle)
- docs/technical-verified/data-contracts.md
- docs/technical-verified/architecture-and-boundaries.md
- docs/prd-wave-details/waves/wave-07.md (API patterns, operations model)
- docs/prd-waves/waves/wave-08.md (source wave)

## Selected Backend Wave Boundary
WAVE-08 is pure data persistence — no external integrations, no files, no async processing.

## Neighboring Backend Wave Fit
- WAVE-07: REST endpoints exist for download; WAVE-08 has no REST endpoints (GraphQL-only)
- WAVE-09: WAVE-08 service must expose ListAllByUserID for backup inclusion

## Frontend Pages Context
No frontend page specific to AiReview. PAGE-009 may reference review history. Backend provides GraphQL queries/mutations only.

## Codebase Evidence
AiReview entity: id, userId, dateRangeStart, dateRangeEnd, aiResponseText, userNotes, plannedActions, createdAt, updatedAt. All text fields. Lifecycle: created → read → updated → deleted. No state machine.

## Proposed Details

### Data Lifecycle
- Created: user pastes AI response text, links to date range, optionally adds notes/planned actions
- Read: user views review history, filtered by date range
- Updated: user modifies notes, planned actions, or date range (text content is manual user entry)
- Deleted: user removes review entry
- No TTL, no cleanup needed (user-managed)

### API Design
GraphQL-only (no REST endpoints):
- Mutation `createAiReview(input: CreateAiReviewInput!): AiReviewResult!`
- Mutation `updateAiReview(id: ID!, input: UpdateAiReviewInput!): AiReviewResult!`
- Mutation `deleteAiReview(id: ID!): AiReviewResult!`
- Query `aiReview(id: ID!): AiReviewResult!`
- Query `aiReviews(dateRangeStart: Date, dateRangeEnd: Date): AiReviewsResult!`

### Input Types
```
input CreateAiReviewInput {
  dateRangeStart: Date!
  dateRangeEnd: Date!
  aiResponseText: String!
  userNotes: String
  plannedActions: String
}

input UpdateAiReviewInput {
  dateRangeStart: Date
  dateRangeEnd: Date
  aiResponseText: String
  userNotes: String
  plannedActions: String
}
```

### Validation
- aiResponseText: required, non-empty
- dateRangeStart, dateRangeEnd: required on create, must be valid dates
- dateRangeEnd >= dateRangeStart validation (service layer)
- User-scoped: all queries filtered by userId from session context
- No max text length validation (PostgreSQL TEXT handles up to 1GB)

### Rollout
1. Create migration 00093_ai_reviews.sql
2. Add sqlc queries, run bun run codegen
3. Create AiReview model types
4. Create AiReview repository
5. Create AiReview service
6. Create GraphQL schema + resolvers
7. Wire in main.go, resolver.go, atlas-gqlgen.yml

### Rollback
1. goose down 00093
2. Remove new source files
3. Revert main.go, resolver.go, atlas-gqlgen.yml, schema.graphql
4. bun run codegen to purge generated code

### Compatibility
- Additive changes only — no existing tables modified
- No impact on existing API endpoints or resolvers
- No impact on existing migrations

### Observability
- Log markers: [AiReview][create], [AiReview][update], [AiReview][delete], [AiReview][list]
- No ai_response_text, user_notes, or planned_actions content in logs (following AC-118 pattern)
- Only metadata (ID, date range) logged

## Questions Raised
Q-W08-DIO-001: GraphQL-only or include REST endpoints? WAVE-07 has REST for file download. AiReview has no files. Recommendation: GraphQL-only.
Q-W08-DIO-002: Should Max reviews per user be bounded? Recommendation: No for MVP (unbounded, user-managed).

## Traceability Candidates
- GraphQL-only approach → follows WAVE-07 pattern for CRUD, no need for REST since no file download
- Validation patterns → WAVE-07 service validation patterns
- Rollout/rollback → WAVE-07 deployment patterns