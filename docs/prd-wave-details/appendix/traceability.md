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
| SLICE-W05-001 | DB migrations | docs/prd-waves/waves/wave-05.md, docs/product-verified/domain-model.md |
| SLICE-W05-002 | sqlc queries | docs/product-verified/domain-model.md, apps/api/internal/atlas (existing patterns) |
| SLICE-W05-003 | Repository adapters | apps/api/internal/atlas/repository/postgres/settings_repo.go (pattern) |
| SLICE-W05-004 | Models | apps/api/internal/atlas/models/settings.go (pattern) |
| SLICE-W05-005 | Services layer | apps/api/internal/atlas/service/settings_service.go (pattern), docs/product-verified/business-rules.md |
| SLICE-W05-006 | GraphQL schema | apps/api/internal/atlas/graph/schema/settings.graphql (pattern) |
| SLICE-W05-007 | GraphQL resolvers | apps/api/internal/atlas/graph/resolver/settings.go (pattern) |
| SLICE-W05-008 | gqlgen config and wiring | apps/api/atlas-gqlgen.yml, apps/api/cmd/server/main.go |
## Acceptance Criteria Map
| AC Range | Source |
| --- | --- |
| AC-W04-001–010 | docs/product-verified/functional-spec.md §12 (Cardio REQ-007), AC-012, AC-013, AC-048–051 |
| AC-W04-011–018 | docs/product-verified/functional-spec.md §13 (Body Tracking REQ-008/REQ-009), AC-016, AC-057 |
| AC-W04-019–031 | docs/product-verified/functional-spec.md §13 (Body Tracking), AC-014, AC-015, AC-052–056 |
| AC-W04-032–038 | docs/product-verified/functional-spec.md §14 (Progress Photos), AC-056 |
| AC-W04-039–042 | docs/product-verified/domain-model.md WeekFlag entity |
| AC-W04-043–044 | docs/technical-verified/auth-security-compliance.md TDEC-037 |
| AC-W05-001–006 | docs/product-verified/functional-spec.md §15 (Nutrition REQ-010), AC-017, AC-058, AC-059 |
| AC-W05-007–017 | docs/product-verified/functional-spec.md §15 (Nutrition REQ-010), AC-018, AC-060 |
| AC-W05-018–027 | docs/product-verified/functional-spec.md §15 (Nutrition REQ-011), AC-019, AC-061, AC-062, AC-113 |
| AC-W05-028–033 | docs/product-verified/business-rules.md RULE-010, RULE-011, docs/product-verified/edge-cases.md EDGE-003, EDGE-009 |
| AC-W05-034–036 | docs/technical-verified/auth-security-compliance.md TDEC-037, soft-delete design decisions |
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
| EC-W05-001 | All ACs passing |
| EC-W05-002–003 | docs/technical-verified/implementation-slices.md codegen requirements |
| EC-W05-004–005 | docs/technical-verified/auth-security-compliance.md |
| EC-W05-006 | Migration patterns from WAVE-01 |
| EC-W05-007–009 | docs/product-verified/edge-cases.md EDGE-003, docs/product-verified/business-rules.md RULE-006 |
| EC-W05-010 | docs/product-verified/user-flows.md §26.7, §26.8 |
| EC-W05-011 | docs/verification-plan.xml V-M-API |
| EC-W05-012 | docs/technical-verified/operations-observability.md |
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
| TEST-W05-001–003 | docs/product-verified/functional-spec.md §15, AC-W05-001–006 |
| TEST-W05-004–007 | docs/product-verified/functional-spec.md §15, AC-W05-007–017 |
| TEST-W05-008–009 | docs/product-verified/domain-model.md NutritionTemplateItem, AC-W05-014–017 |
| TEST-W05-010–014 | docs/product-verified/functional-spec.md §15, AC-W05-018–027 |
| TEST-W05-015–018 | docs/product-verified/business-rules.md RULE-010, RULE-011, AC-W05-029–033 |
| TEST-W05-019–020 | docs/product-verified/user-flows.md §26.7, §26.8, AC-W05-028 |
| TEST-W05-021 | docs/technical-verified/auth-security-compliance.md TDEC-037 |
| TEST-W05-022 | docs/prd-wave-details/waves/wave-04.md migration smoke test pattern |
| TEST-W05-023 | docs/technical-verified/implementation-slices.md codegen requirements |
| TEST-W05-024 | docs/technical-verified/operations-observability.md |
| TEST-W05-025–026 | docs/verification-plan.xml V-M-API |
| TEST-W05-027–029 | docs/product-verified/domain-model.md Nutrition invariants |
| TEST-W05-030 | docs/verification-plan.xml admin auth regression |
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
| apps/api/cmd/server/main.go | WAVE-04/05 | Wiring |
| apps/api/internal/atlas/graph/schema/nutrition.graphql | WAVE-05 | GraphQL types |
| apps/api/internal/atlas/models/nutrition.go | WAVE-05 | Domain models |
| apps/api/internal/atlas/service/nutrition_*_service.go | WAVE-05 | Services |
| apps/api/internal/atlas/repository/postgres/nutrition_*_repo.go | WAVE-05 | Repository adapters |
| apps/api/internal/atlas/graph/resolver/nutrition.go | WAVE-05 | GraphQL resolvers |
| apps/api/internal/atlas/graph/resolver/resolver.go | WAVE-05 | Resolver container |
| apps/api/atlas-gqlgen.yml | WAVE-05 | gqlgen config |
| apps/api/internal/repository/postgres/migrations/00081_nutrition_tables.sql | WAVE-05 | Migration |
| apps/api/internal/repository/postgres/queries/nutrition_*.sql | WAVE-05 | sqlc query defs |
## Question Map
| Question ID | Source |
| --- | --- |
| DQ-W04-001 | planner-sequencing-fit-attempt-1.md — WAVE-03 DailyLog dependency |
| DQ-W04-002 | planner-product-ac-attempt-1.md — EDGE-006 photo count |
| DQ-W04-003 | planner-product-ac-attempt-1.md — EDGE-007 measurement validation |
| DQ-W04-004 | planner-security-compliance-attempt-1.md — TDEC-008 signed URLs |
| DQ-W04-005 | planner-data-integration-ops-attempt-1.md — WAVE-01 MediaConfig path |
| DQ-W05-001 | planner-product-ac-attempt-2.md — EDGE-019 soft-delete |
| DQ-W05-002 | planner-product-ac-attempt-1.md — RULE-020 single template |
| DQ-W05-003 | planner-product-ac-attempt-1.md — mealLabel type |
| DQ-W05-004 | planner-architecture-codebase-attempt-1.md — macro calc location |
| DQ-W05-005 | planner-data-integration-ops-attempt-1.md — macro query design |
| DQ-W05-007 | planner-security-compliance-attempt-1.md — soft-delete recovery |
| DQ-W05-008 | planner-testing-exit-attempt-1.md — macro test scope |
| DQ-W05-009 | planner-sequencing-fit-attempt-1.md — migration number collision |
## Source Map
| Artifact | Source |
| --- | --- |
| WAVE-04 detailed brief | docs/prd-waves/waves/wave-04.md, docs/product-verified/*, docs/technical-verified/* |
| WAVE-05 detailed brief | docs/prd-waves/waves/wave-05.md, docs/product-verified/*, docs/technical-verified/* |
| Planner reports | docs/prd-waves/waves/wave-05.md, docs/product-verified, docs/technical-verified, docs/prd-wave-details/waves/* |
| Reviewer reports | Planner reports, source evidence, codebase patterns |
| Design decisions | Planner reports, product-verified edge-cases.md, business-rules.md |