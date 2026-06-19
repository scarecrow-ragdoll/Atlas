# Traceability

## Slice Map
| Slice ID | Source | Entity | Verification |
| --- | --- | --- | --- |
| SLICE-W03-001 | DB migration: daily_logs | domain-model.md DailyLog | TEST-W03-023 |
| SLICE-W03-002 | DB migration: workout_exercises | domain-model.md WorkoutExercise | TEST-W03-023 |
| SLICE-W03-003 | DB migration: workout_sets | domain-model.md WorkoutSet | TEST-W03-023 |
| SLICE-W03-004 | DB migration: cardio_entries | domain-model.md CardioEntry | TEST-W03-023 |
| SLICE-W03-005 | sqlc queries: daily_logs | data-contracts.md TDEC-020 | TEST-W03-001 |
| SLICE-W03-006 | sqlc queries: workout_exercises | data-contracts.md TDEC-020 | TEST-W03-002 |
| SLICE-W03-007 | sqlc queries: workout_sets | data-contracts.md TDEC-020 | TEST-W03-003 |
| SLICE-W03-008 | sqlc queries: cardio_entries | data-contracts.md TDEC-020 | TEST-W03-004 |
| SLICE-W03-009 | DailyLog repository | user_repo.go pattern | TEST-W03-001 |
| SLICE-W03-010 | WorkoutExercise repository | user_repo.go pattern | TEST-W03-002 |
| SLICE-W03-011 | WorkoutSet repository | user_repo.go pattern | TEST-W03-003 |
| SLICE-W03-012 | CardioEntry repository | user_repo.go pattern | TEST-W03-004 |
| SLICE-W03-013 | Workout service | admin_auth.go pattern | TEST-W03-005 through 012 |
| SLICE-W03-014 | GraphQL schema | admin_auth.graphql pattern | TEST-W03-025 |
| SLICE-W03-015 | GraphQL resolvers + wiring | main.go pattern | TEST-W03-013 through 022 |

## Acceptance Criteria Map
| AC ID | Product AC | Source | Verification |
| --- | --- | --- | --- |
| AC-W03-001 | AC-005, AC-006, AC-035, AC-036, AC-037 | functional-spec.md §10.2 | TEST-W03-013, TEST-W03-022 |
| AC-W03-002 | AC-038 | user-flows.md §26.2 | TEST-W03-014 |
| AC-W03-003 | AC-038 | user-flows.md §26.2, RULE-016 | TEST-W03-010 |
| AC-W03-004 | AC-038 | user-flows.md §26.3 | TEST-W03-010 |
| AC-W03-005 | AC-007 | functional-spec.md §10.4 | TEST-W03-015 |
| AC-W03-006 | AC-039, AC-041, AC-003 | functional-spec.md §10.6, RULE-017 | TEST-W03-006 |
| AC-W03-007 | AC-007 | functional-spec.md §10.4 | TEST-W03-005 |
| AC-W03-008 | AC-008 | functional-spec.md §10.5, RULE-004 | TEST-W03-016 |
| AC-W03-009 | AC-009 | functional-spec.md §10.5 | TEST-W03-008 |
| AC-W03-010 | AC-010 | functional-spec.md §10.5 | TEST-W03-008 |
| AC-W03-011 | AC-008 | functional-spec.md §10.5 | TEST-W03-016 |
| AC-W03-012 | AC-040 | functional-spec.md §10.5 | TEST-W03-005 |
| AC-W03-013 | AC-011, AC-042 | functional-spec.md §10.4 | TEST-W03-005 |
| AC-W03-014 | AC-040 | functional-spec.md §10.4 | TEST-W03-015 |
| AC-W03-015 | AC-012 | functional-spec.md §12.2 | TEST-W03-017 |
| AC-W03-016 | AC-013 | functional-spec.md §12.4 | TEST-W03-017 |
| AC-W03-017 | — | edge-cases.md EDGE-005 | TEST-W03-009 |
| AC-W03-018 | AC-008 | functional-spec.md §10.5 | TEST-W03-016 |
| AC-W03-019 | AC-008 | functional-spec.md §10.5 | TEST-W03-016 |
| AC-W03-020 | AC-012 | functional-spec.md §12.2 | TEST-W03-017 |
| AC-W03-021 | — | edge-cases.md EDGE-005 | TEST-W03-009 |
| AC-W03-022 | — | edge-cases.md EDGE-001, RULE-004 | TEST-W03-007, TEST-W03-019 |
| AC-W03-023 | — | edge-cases.md EDGE-001, RULE-004 | TEST-W03-007, TEST-W03-019 |
| AC-W03-024 | — | business-rules.md RULE-017 | TEST-W03-011, TEST-W03-021 |
| AC-W03-025 | AC-038 | functional-spec.md §10.2 | TEST-W03-022 |
| AC-W03-026 | — | domain-model.md invariant 6 | TEST-W03-005 |
| AC-W03-027 | — | domain-model.md invariant 5 | TEST-W03-005 |
| AC-W03-028 | — | data-contracts.md TDEC-020 | TEST-W03-007 |
| AC-W03-029 | AC-110 | auth-security-compliance.md TDEC-037 | TEST-W03-018 |
| AC-W03-030 | AC-110 | auth-security-compliance.md TDEC-037 | TEST-W03-018 |

## Exit Criteria Map
| EC ID | Source | Verification |
| --- | --- | --- |
| EC-W03-001 | All ACs | TEST-W03-001 through TEST-W03-022 |
| EC-W03-002 | gqlgen.yml | TEST-W03-026 |
| EC-W03-003 | sqlc.yaml | TEST-W03-026 |
| EC-W03-004 | TDEC-024 | TEST-W03-023 |
| EC-W03-005 | user-flows.md §26.2 | TEST-W03-010 |
| EC-W03-006 | functional-spec.md §10.6, RULE-017 | TEST-W03-006 |
| EC-W03-007 | domain-model.md invariant 6 | TEST-W03-005 |
| EC-W03-008 | domain-model.md invariant 5 | TEST-W03-005 |
| EC-W03-009 | functional-spec.md §12 | TEST-W03-017 |
| EC-W03-010 | TDEC-005 | TEST-W03-009 |
| EC-W03-011 | auth-security-compliance.md TDEC-037 | TEST-W03-018 |
| EC-W03-012 | auth-security-compliance.md TDEC-004 | TEST-W03-020 |
| EC-W03-013 | edge-cases.md EDGE-001 | TEST-W03-007, TEST-W03-019 |
| EC-W03-014 | functional-spec.md §10.2 | TEST-W03-022 |
| EC-W03-015 | WAVE-01 regression | TEST-W03-027 |
| EC-W03-016 | WAVE-02 regression | TEST-W03-028 |
| EC-W03-017 | quality gate | TEST-W03-024 |
| EC-W03-018 | quality gate | typecheck |

## Verification Obligation Map
| Test ID | Source | Verification |
| --- | --- | --- |
| TEST-W03-001 | repos/daily_log_repo.go | AC-W03-001, EC-W03-001 |
| TEST-W03-002 | repos/workout_exercise_repo.go | AC-W03-005, EC-W03-001 |
| TEST-W03-003 | repos/workout_set_repo.go | AC-W03-008, EC-W03-001 |
| TEST-W03-004 | repos/cardio_entry_repo.go | AC-W03-015, EC-W03-001 |
| TEST-W03-005 | service/workout.go domain-model invariants | AC-W03-007, AC-W03-012, AC-W03-013, AC-W03-026, AC-W03-027, EC-W03-007, EC-W03-008 |
| TEST-W03-006 | service/workout.go RULE-017 | AC-W03-006, EC-W03-006 |
| TEST-W03-007 | service/workout.go EDGE-001, RULE-004 | AC-W03-022, AC-W03-023, EC-W03-013 |
| TEST-W03-008 | service/workout.go RPE/RIR bounds | AC-W03-009, AC-W03-010 |
| TEST-W03-009 | service/workout.go TDEC-005 | AC-W03-017, AC-W03-021, EC-W03-010 |
| TEST-W03-010 | service/workout.go upsert semantics | AC-W03-003, AC-W03-004, EC-W03-005 |
| TEST-W03-011 | service/workout.go FK constraints | AC-W03-024 |
| TEST-W03-012 | service/workout.go unique constraint | AC-W03-002 |
| TEST-W03-013 | resolvers/dailyLogByDate | AC-W03-001 |
| TEST-W03-014 | resolvers/upsertDailyLog | AC-W03-002 |
| TEST-W03-015 | resolvers/addWorkoutExercise | AC-W03-005, AC-W03-014 |
| TEST-W03-016 | resolvers/set CRUD | AC-W03-008, AC-W03-011, AC-W03-018, AC-W03-019 |
| TEST-W03-017 | resolvers/cardio CRUD | AC-W03-015, AC-W03-016, AC-W03-020 |
| TEST-W03-018 | resolvers/auth guard TDEC-037 | AC-W03-029, AC-W03-030, EC-W03-011 |
| TEST-W03-019 | resolvers/validation errors | AC-W03-022, AC-W03-023, AC-W03-028, EC-W03-013 |
| TEST-W03-020 | log sanitization TDEC-004 | EC-W03-012 |
| TEST-W03-021 | FK exercise_id constraint | AC-W03-024 |
| TEST-W03-022 | empty date query | AC-W03-025, EC-W03-014 |
| TEST-W03-023 | migrations smoke test TDEC-024 | EC-W03-004 |
| TEST-W03-024 | lint | EC-W03-017 |
| TEST-W03-025 | schema validation | EC-W03-002 |
| TEST-W03-026 | codegen drift | EC-W03-002, EC-W03-003 |
| TEST-W03-027 | WAVE-01 regression | EC-W03-015 |
| TEST-W03-028 | WAVE-02 regression | EC-W03-016 |

## Code Touchpoint Map
| File | Slice | Purpose |
| --- | --- | --- |
| apps/api/internal/repository/postgres/migrations/00082_daily_logs.sql | SLICE-W03-001 | daily_logs table |
| apps/api/internal/repository/postgres/migrations/00083_workout_exercises.sql | SLICE-W03-002 | workout_exercises table |
| apps/api/internal/repository/postgres/migrations/00084_workout_sets.sql | SLICE-W03-003 | workout_sets table |
| apps/api/internal/repository/postgres/migrations/00085_cardio_entries.sql | SLICE-W03-004 | cardio_entries table |
| apps/api/internal/repository/postgres/queries/daily_logs.sql | SLICE-W03-005 | sqlc queries |
| apps/api/internal/repository/postgres/queries/workout_exercises.sql | SLICE-W03-006 | sqlc queries |
| apps/api/internal/repository/postgres/queries/workout_sets.sql | SLICE-W03-007 | sqlc queries |
| apps/api/internal/repository/postgres/queries/cardio_entries.sql | SLICE-W03-008 | sqlc queries |
| apps/api/internal/repository/postgres/daily_log_repo.go | SLICE-W03-009 | repository adapter |
| apps/api/internal/repository/postgres/workout_exercise_repo.go | SLICE-W03-010 | repository adapter |
| apps/api/internal/repository/postgres/workout_set_repo.go | SLICE-W03-011 | repository adapter |
| apps/api/internal/repository/postgres/cardio_entry_repo.go | SLICE-W03-012 | repository adapter |
| apps/api/internal/service/workout.go | SLICE-W03-013 | service layer |
| libs/graphql/schema/workout.graphql | SLICE-W03-014 | GraphQL schema |
| apps/api/internal/graph/workout.resolvers.go | SLICE-W03-015 | GraphQL resolvers |
| apps/api/cmd/server/main.go | SLICE-W03-015 | wiring |

## Question Map
| Question ID | Scope | Severity | Status | Source |
| --- | --- | --- | --- | --- |
| Q-WORKOUT-001 | operations | needs-owner-decision | open | docs/prd-waves/open-questions.md |
| DQ-W03-001 | operations | deferred | deferred | planner-data-integration-ops-attempt-1.md |
| DQ-W03-002 | data-ops | wave-blocking | resolved | planner-architecture-codebase-attempt-1.md |
| DQ-W03-003 | data-ops | wave-blocking | resolved | planner-sequencing-fit-attempt-1.md |
| DQ-W03-004 | product | needs-owner-decision | resolved | planner-product-ac-attempt-1.md |
| DQ-W03-005 | product | needs-owner-decision | resolved | planner-product-ac-attempt-1.md |
| DQ-W03-006 | data-ops | needs-owner-decision | resolved | planner-data-integration-ops-attempt-1.md |
| DQ-W03-007 | product | wave-blocking | resolved | planner-product-ac-attempt-1.md |

## Source Map
| Source | Usage |
| --- | --- |
| docs/prd-waves/waves/wave-03.md | Primary source — outcomes, scope, dependencies, risk class |
| docs/prd-waves/frontend-pages/page-002.md | Backend dependency context for GraphQL operations |
| docs/product-verified/domain-model.md | Entity definitions, attributes, relationships, invariants |
| docs/product-verified/functional-spec.md | Feature behavior, validations (Workout Diary §10) |
| docs/product-verified/acceptance-criteria.md | Product-level AC-005 through AC-011, AC-035 through AC-042 |
| docs/product-verified/user-flows.md | Workout entry user flows (§26.2, §26.3) |
| docs/product-verified/business-rules.md | RULE-004, RULE-016, RULE-017 |
| docs/product-verified/edge-cases.md | EDGE-001, EDGE-004, EDGE-005, EDGE-016 |
| docs/technical-verified/api-contracts.md | Hybrid GraphQL/REST, error format (TDEC-001, TDEC-027) |
| docs/technical-verified/data-contracts.md | Entity contracts, indexes (TDEC-020, TDEC-021) |
| docs/technical-verified/auth-security-compliance.md | PIN auth, audit, log privacy (TDEC-037, TDEC-004) |
| docs/technical-verified/operations-observability.md | Log markers, error format |
| docs/technical-verified/implementation-slices.md | Slice 2 mapping |
| docs/technical-verified/testing-and-delivery.md | Test strategy, fixture strategy (TDEC-056) |
| docs/prd-wave-details/waves/wave-01.md | WAVE-01 dependency contracts |
| docs/prd-wave-details/waves/wave-02.md | WAVE-02 allExercises contract |
