# Traceability
## Slice Map
| Slice ID | Name | Source |
| --- | --- | --- |
| SLICE-W04-001 | DB migrations | docs/prd-waves/waves/wave-04.md, docs/product-verified/domain-model.md |
| SLICE-W04-002 | sqlc queries | docs/product-verified/domain-model.md, docs/technical-verified/data-contracts.md |
| SLICE-W04-003 | Repository adapters | apps/api/internal/repository/postgres/user_repo.go (pattern), docs/prd-wave-details/waves/wave-02.md |
| SLICE-W04-004 | Services layer | apps/api/internal/service/exercise.go (pattern), docs/product-verified/business-rules.md |
| SLICE-W04-005 | GraphQL schema | libs/graphql/schema/exercises.graphql (pattern), docs/product-verified/domain-model.md |
| SLICE-W04-006 | GraphQL resolvers | apps/api/internal/graph/exercise.resolvers.go (pattern) |
| SLICE-W04-007 | ProgressPhoto REST handler | apps/api/internal/handler/exercise_media.go (pattern) |
| SLICE-W04-008 | Main wiring | apps/api/cmd/server/main.go (pattern) |
## Acceptance Criteria Map
| AC Range | Source |
| --- | --- |
| AC-W04-001–010 | docs/product-verified/functional-spec.md §12 (Cardio REQ-007), AC-012, AC-013, AC-048–051 |
| AC-W04-011–018 | docs/product-verified/functional-spec.md §13 (Body Tracking REQ-008/REQ-009), AC-016, AC-057 |
| AC-W04-019–031 | docs/product-verified/functional-spec.md §13 (Body Tracking), AC-014, AC-015, AC-052–056 |
| AC-W04-032–038 | docs/product-verified/functional-spec.md §14 (Progress Photos), AC-056 |
| AC-W04-039–042 | docs/product-verified/domain-model.md WeekFlag entity |
| AC-W04-043–044 | docs/technical-verified/auth-security-compliance.md TDEC-037 |
## Exit Criteria Map
| EC Range | Source |
| --- | --- |
| EC-W04-001 | All ACs passing |
| EC-W04-002–003 | docs/technical-verified/implementation-slices.md codegen requirements |
| EC-W04-004–005 | docs/technical-verified/auth-security-compliance.md |
| EC-W04-006 | Migration patterns from WAVE-01/WAVE-02 |
| EC-W04-007–009 | docs/product-verified/edge-cases.md EDGE-006, EDGE-007 |
| EC-W04-010 | docs/product-verified/user-flows.md §26.4, §26.5, §26.6 |
| EC-W04-011 | docs/verification-plan.xml V-M-API |
| EC-W04-012 | docs/technical-verified/operations-observability.md |
| EC-W04-013–014 | docs/product-verified/domain-model.md CardioEntry, BodyCheckIn invariants |
## Verification Obligation Map
| Test Range | Source |
| --- | --- |
| TEST-W04-001–004 | docs/product-verified/functional-spec.md §12, AC-W04-001–010 |
| TEST-W04-005–008 | docs/product-verified/functional-spec.md §13, AC-W04-011–018 |
| TEST-W04-009–014 | docs/product-verified/functional-spec.md §13, AC-W04-019–031 |
| TEST-W04-015–019 | docs/product-verified/functional-spec.md §14, AC-W04-032–038 |
| TEST-W04-020–022 | docs/product-verified/domain-model.md WeekFlag, AC-W04-039–042 |
| TEST-W04-023–024 | docs/technical-verified/auth-security-compliance.md TDEC-037 |
| TEST-W04-025 | docs/prd-wave-details/waves/wave-02.md migration smoke test pattern |
| TEST-W04-026–029 | docs/verification-plan.xml V-M-API, V-M-GRAPHQL-SCHEMA |
| TEST-W04-030 | docs/product-verified/user-flows.md §26.4, §26.5, §26.6 |
## Code Touchpoint Map
| File | Wave | Purpose |
| --- | --- | --- |
| apps/api/internal/repository/postgres/migrations/00082–00087.sql | WAVE-04 | New entity tables |
| apps/api/internal/repository/postgres/queries/*.sql | WAVE-04 | sqlc query defs |
| apps/api/internal/repository/postgres/*_repo.go | WAVE-04 | Repository adapters |
| apps/api/internal/service/*.go | WAVE-04 | Transport-neutral services |
| apps/api/internal/graph/*.resolvers.go | WAVE-04 | GraphQL resolvers |
| apps/api/internal/handler/progress_photo_handler.go | WAVE-04 | REST handler |
| libs/graphql/schema/*.graphql | WAVE-04 | GraphQL types |
| apps/api/cmd/server/main.go | WAVE-04 | Wiring |
## Question Map
| Question ID | Source |
| --- | --- |
| DQ-W04-001 | planner-sequencing-fit-attempt-1.md — WAVE-03 DailyLog dependency |
| DQ-W04-002 | planner-product-ac-attempt-1.md — EDGE-006 photo count |
| DQ-W04-003 | planner-product-ac-attempt-1.md — EDGE-007 measurement validation |
| DQ-W04-004 | planner-security-compliance-attempt-1.md — TDEC-008 signed URLs |
| DQ-W04-005 | planner-data-integration-ops-attempt-1.md — WAVE-01 MediaConfig path |
## Source Map
| Artifact | Source |
| --- | --- |
| WAVE-04 detailed brief | docs/prd-waves/waves/wave-04.md, docs/product-verified/*, docs/technical-verified/* |
| Planner reports | docs/prd-waves/waves/wave-04.md, docs/product-verified, docs/prd-wave-details/waves/wave-01.md, wave-02.md |
| Reviewer reports | Planner reports, source evidence, codebase patterns |
| Design decisions | Planner reports, product-verified edge-cases.md, business-rules.md |