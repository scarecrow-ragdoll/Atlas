<!-- FILE: docs/prd-wave-details/appendix/traceability.md -->
<!-- VERSION: 1.0.0 -->

# Traceability

## Slice Map

| Slice ID | Source Wave Capability | Description |
|----------|----------------------|-------------|
| SLICE-W08-001 | CAP-W08-001 (AiReview CRUD) | Migration 00093_ai_reviews.sql |
| SLICE-W08-002 | CAP-W08-001 (AiReview CRUD) | AiReview model types |
| SLICE-W08-003 | CAP-W08-001 (AiReview CRUD) | AiReview sqlc queries |
| SLICE-W08-004 | CAP-W08-001 (AiReview CRUD) | AiReview repository |
| SLICE-W08-005 | CAP-W08-002 through CAP-W08-005 | AiReview service with period linkage, notes, planned actions, history |
| SLICE-W08-006 | CAP-W08-001, CAP-W08-005 | GraphQL schema + resolvers |
| SLICE-W08-007 | CAP-W08-001 | Main wiring (resolver.go, main.go, gqlgen) |

## Acceptance Criteria Map

| AC ID | Source | Description |
|-------|--------|-------------|
| AC-W08-001 | AC-025, §19 | createAiReview mutation with required fields |
| AC-W08-002 | AC-090, §19.2 | Paste AI response text |
| AC-W08-003 | AC-091, §19.2 | Link to date range with validation |
| AC-W08-004 | AC-092, §19.2 | Add notes and planned actions |
| AC-W08-005 | W08-005, §19 | View review history (createdAt DESC) |
| AC-W08-006 | Functional spec §19 | Filter reviews by date range |
| AC-W08-007 | CAP-W08-001 | Update review entry |
| AC-W08-008 | CAP-W08-001 | Delete review entry |

## Exit Criteria Map

| EC ID | Validation Type | Source |
|-------|----------------|--------|
| EC-W08-001 | Migration | SLICE-W08-001 |
| EC-W08-002 | Codegen (sqlc) | SLICE-W08-003 |
| EC-W08-003 | Build | All slices |
| EC-W08-004 | Codegen (gqlgen) | SLICE-W08-006 |
| EC-W08-005 | Service tests | SLICE-W08-005 |
| EC-W08-006 | Resolver tests | SLICE-W08-006 |
| EC-W08-007 | Repository integration tests | SLICE-W08-004 |
| EC-W08-008 | Codegen drift | All slices |
| EC-W08-009 | Auth protection | All slices |
| EC-W08-010 | Lint | All slices |

## Verification Obligation Map

| TEST ID | Layer | Description |
|---------|-------|-------------|
| TEST-W08-001 | Service | Create success |
| TEST-W08-002 | Service | Empty text rejection |
| TEST-W08-003 | Service | Invalid date range |
| TEST-W08-004 | Service | List ordered by createdAt DESC |
| TEST-W08-005 | Service | Date range filter |
| TEST-W08-006 | Service | Update success |
| TEST-W08-007 | Service | Update ownership check |
| TEST-W08-008 | Service | Delete success |
| TEST-W08-009 | Service | Log privacy (AC-118) |
| TEST-W08-010 | Resolver | Create resolver |
| TEST-W08-011 | Resolver | List resolver |
| TEST-W08-012 | Resolver | Delete resolver |
| TEST-W08-013 | Repository | Repository operations (INTEGRATION_TESTS=1) |
| TEST-W08-014 | Migration | Migration 00093 applies cleanly (INTEGRATION_TESTS=1) |

## Code Touchpoint Map

| File | Slice | Operation |
|------|-------|-----------|
| internal/repository/postgres/migrations/00093_ai_reviews.sql | SLICE-W08-001 | Create table DDL |
| internal/atlas/models/ai_review.go | SLICE-W08-002 | New file |
| internal/repository/postgres/queries/ai_reviews.sql | SLICE-W08-003 | New file |
| internal/repository/postgres/ai_review_repo.go | SLICE-W08-004 | New file |
| internal/atlas/service/ai_review_service.go | SLICE-W08-005 | New file |
| internal/atlas/graph/schema/ai_review.graphql | SLICE-W08-006 | New file |
| internal/atlas/graph/resolver/ai_review.go | SLICE-W08-006 | New file |
| internal/atlas/graph/resolver/resolver.go | SLICE-W08-007 | Add AiReviewService field |
| cmd/server/main.go | SLICE-W08-007 | Wire repo→service→resolver |
| atlas-gqlgen.yml | SLICE-W08-007 | Add type bindings |
| internal/atlas/graph/schema/schema.graphql | SLICE-W08-007 | Add extend type Query/Mutation |

## Question Map

| Question ID | Source Report | Status |
|-------------|-------------|--------|
| DQ-W08-001 | planner-product-ac-attempt-1, planner-architecture-codebase-attempt-1 | open (needs-owner-decision) |
| DQ-W08-002 | planner-sequencing-fit-attempt-1, wave-07.md | open (needs-owner-decision) |

## Source Map

| Source Document | Relevant Artifacts |
|----------------|-------------------|
| docs/prd-waves/waves/wave-08.md | All slices, AC-W08-001 through AC-W08-008 |
| docs/product-verified/functional-spec.md §19 | AC-W08-001 through AC-W08-006 |
| docs/product-verified/domain-model.md#AiReview | SLICE-W08-001, SLICE-W08-002 |
| docs/product-verified/acceptance-criteria.md | AC-W08-001 through AC-W08-004 |
| docs/prd-wave-details/waves/wave-07.md | All slice patterns |
| docs/knowledge-graph.xml | Codebase fit module references |
| docs/verification-plan.xml | Verification entry patterns |