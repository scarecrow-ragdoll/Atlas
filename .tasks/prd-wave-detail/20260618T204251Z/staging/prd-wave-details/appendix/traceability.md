# Traceability

## Slice Map
| Slice ID | Source Outcome | Source Capability | Source Doc |
| --- | --- | --- | --- |
| SLICE-W02-001 | OUT-W02-001 | CAP-W02-001 | docs/prd-waves/waves/wave-02.md |
| SLICE-W02-002 | OUT-W02-001, OUT-W02-002 | CAP-W02-001, CAP-W02-002 | docs/prd-waves/waves/wave-02.md |
| SLICE-W02-003 | OUT-W02-001 | CAP-W02-001 | docs/prd-waves/waves/wave-02.md |
| SLICE-W02-004 | OUT-W02-001, OUT-W02-003 | CAP-W02-001, CAP-W02-003 | docs/prd-waves/waves/wave-02.md |
| SLICE-W02-005 | OUT-W02-001, OUT-W02-002, OUT-W02-003 | CAP-W02-001, CAP-W02-002, CAP-W02-003, CAP-W02-004 | docs/prd-waves/waves/wave-02.md |
| SLICE-W02-006 | OUT-W02-001, OUT-W02-004 | CAP-W02-001, CAP-W02-005 | docs/prd-waves/waves/wave-02.md |
| SLICE-W02-007 | OUT-W02-002 | CAP-W02-002 | docs/prd-waves/waves/wave-02.md |
| SLICE-W02-008 | OUT-W02-001, OUT-W02-004 | CAP-W02-001 | docs/prd-waves/waves/wave-02.md |

## Acceptance Criteria Map
| AC ID | Source Requirement | Source Doc | Planner |
| --- | --- | --- | --- |
| AC-W02-001 | CAP-W02-001, AC-002, AC-043, AC-044 | docs/product-verified/acceptance-criteria.md | product-ac |
| AC-W02-002 | CAP-W02-001, OUT-W02-001 | docs/prd-waves/waves/wave-02.md | product-ac |
| AC-W02-003 | CAP-W02-001, OUT-W02-001 | docs/prd-waves/waves/wave-02.md | product-ac |
| AC-W02-004 | CAP-W02-001, OUT-W02-001 | docs/prd-waves/waves/wave-02.md | product-ac |
| AC-W02-005 | CAP-W02-001, AC-044 | docs/product-verified/acceptance-criteria.md | product-ac |
| AC-W02-006 | CAP-W02-003, AC-003, OUT-W02-003 | docs/product-verified/acceptance-criteria.md | product-ac |
| AC-W02-007 | CAP-W02-005, AC-047 | docs/product-verified/acceptance-criteria.md | product-ac |
| AC-W02-008 | CAP-W02-005, AC-047 | docs/product-verified/acceptance-criteria.md | product-ac |
| AC-W02-009 | OUT-W02-004, CAP-W02-005 | docs/prd-waves/waves/wave-02.md | product-ac |
| AC-W02-010 | CAP-W02-005 | docs/prd-waves/waves/wave-02.md | product-ac |
| AC-W02-011 | CAP-W02-005 | docs/prd-waves/waves/wave-02.md | product-ac |
| AC-W02-012 | CAP-W02-001, EDGE-002 | docs/product-verified/edge-cases.md | product-ac |
| AC-W02-013 | EDGE-002 | docs/product-verified/edge-cases.md | product-ac |
| AC-W02-014 | CAP-W02-003 | docs/prd-waves/waves/wave-02.md | product-ac |
| AC-W02-015 | CAP-W02-002, AC-004, AC-045 | docs/product-verified/acceptance-criteria.md | product-ac |
| AC-W02-016 | CAP-W02-002, OUT-W02-002 | docs/prd-waves/waves/wave-02.md | product-ac |
| AC-W02-017 | CAP-W02-002, AC-046, TDEC-005 | docs/technical-verified/data-contracts.md | product-ac |
| AC-W02-018 | CAP-W02-002, AC-045 | docs/product-verified/acceptance-criteria.md | product-ac |
| AC-W02-019 | OUT-W02-004 | docs/prd-waves/waves/wave-02.md | product-ac |
| AC-W02-020 | TDEC-029, WAVE-01 PIN auth | docs/technical-verified/api-contracts.md | security-compliance |
| AC-W02-021 | TDEC-029, WAVE-01 PIN auth | docs/technical-verified/api-contracts.md | security-compliance |
| AC-W02-022 | TDEC-008 | docs/technical-verified/auth-security-compliance.md | security-compliance |
| AC-W02-023 | TDEC-008 | docs/technical-verified/auth-security-compliance.md | security-compliance |
| AC-W02-024 | EDGE-014 | docs/product-verified/edge-cases.md | security-compliance |

## Exit Criteria Map
| EC ID | Source | Verifies |
| --- | --- | --- |
| EC-W02-001 | AC-W02-001 through AC-W02-024 | All AC coverage |
| EC-W02-002 | gqlgen config | Codegen validity |
| EC-W02-003 | sqlc config | Codegen validity |
| EC-W02-004 | WAVE-01 media scaffold | Backward compatibility |
| EC-W02-005 | WAVE-01 PIN auth | Auth guard coverage |
| EC-W02-006 | WAVE-01 test suite | Regression safety |
| EC-W02-007 | Lint config | Code quality |
| EC-W02-008 | 00080 + 00081 migrations | Migration integrity |
| EC-W02-009 | TDEC-008 | File validation |
| EC-W02-010 | Operations observability | Audit markers |
| EC-W02-011 | Privacy policy | Log privacy |
| EC-W02-012 | Round-trip integration | End-to-end correctness |
| EC-W02-013 | WAVE-03 dependency | Interface stability |

## Verification Obligation Map
| Test ID | AC/EC Coverage | Type |
| --- | --- | --- |
| TEST-W02-001 | AC-W02-001 through AC-W02-007 | unit |
| TEST-W02-002 | AC-W02-012, AC-W02-014 | unit |
| TEST-W02-003 | AC-W02-001 through AC-W02-014, AC-W02-018, AC-W02-019, AC-W02-020 | integration |
| TEST-W02-004 | AC-W02-015 through AC-W02-017 | integration |
| TEST-W02-005 | EC-W02-008 | integration |
| TEST-W02-006 | EC-W02-002, EC-W02-003 | codegen |
| TEST-W02-007 | AC-W02-019, EC-W02-013 | integration |
| TEST-W02-008 | AC-W02-009 | integration |
| TEST-W02-009 | AC-W02-008, AC-W02-010 | integration |
| TEST-W02-010 | AC-W02-003, AC-W02-004 | integration |
| TEST-W02-011 | AC-W02-013 | integration |
| TEST-W02-012 | AC-W02-005, AC-W02-006 | integration |
| TEST-W02-013 | AC-W02-020, EC-W02-005 | integration |
| TEST-W02-014 | AC-W02-021, EC-W02-005 | integration |
| TEST-W02-015 | AC-W02-022 | unit |
| TEST-W02-016 | AC-W02-023, EC-W02-009 | unit |
| TEST-W02-017 | AC-W02-024 | unit |
| TEST-W02-018 | EC-W02-011 | unit |
| TEST-W02-019 | EC-W02-007 | lint |
| TEST-W02-020 | EC-W02-002 | codegen |
| TEST-W02-021 | EC-W02-012 | integration |
| TEST-W02-022 | EC-W02-006 | unit |

## Code Touchpoint Map
| File | Purpose | Slice |
| --- | --- | --- |
| apps/api/internal/repository/postgres/migrations/00080_exercises.sql | Exercise table migration | SLICE-W02-001 |
| apps/api/internal/repository/postgres/migrations/00081_exercise_media.sql | ExerciseMedia table migration | SLICE-W02-001 |
| apps/api/internal/repository/postgres/queries/exercises.sql | sqlc query definitions | SLICE-W02-002 |
| apps/api/internal/repository/postgres/exercise_repo.go | Repository adapter | SLICE-W02-003 |
| apps/api/internal/service/exercise.go | Transport-neutral service | SLICE-W02-004 |
| libs/graphql/schema/exercises.graphql | GraphQL schema | SLICE-W02-005 |
| apps/api/internal/graph/exercise.resolvers.go | GraphQL resolvers | SLICE-W02-006 |
| apps/api/internal/handler/exercise_media.go | REST media handler | SLICE-W02-007 |
| apps/api/cmd/server/main.go | Wiring and route registration | SLICE-W02-008 |

## Question Map
| ID | Source | Status |
| --- | --- | --- |
| DQ-W02-001 | planner-data-integration-ops-attempt-2.md | resolved |
| DQ-W02-002 | planner-product-ac-attempt-2.md | answered |
| DQ-W02-003 | planner-data-integration-ops-attempt-2.md | deferred |
| DQ-W02-005 | planner-security-compliance-attempt-2.md | resolved |
| DQ-W02-006 | planner-security-compliance-attempt-2.md | deferred |
| DQ-W02-007 | planner-testing-exit-attempt-2.md | resolved |
| DQ-W02-008 | planner-sequencing-fit-attempt-2.md | deferred |

## Source Map
| Artifact | Primary Source |
| --- | --- |
| WAVE-02 boundary | docs/prd-waves/waves/wave-02.md |
| Exercise entity | docs/product-verified/domain-model.md |
| ExerciseMedia entity | docs/product-verified/domain-model.md |
| Hybrid API protocol | docs/technical-verified/api-contracts.md (TDEC-001) |
| Error format | docs/technical-verified/api-contracts.md (TDEC-027) |
| File validation | docs/technical-verified/auth-security-compliance.md (TDEC-008) |
| PIN auth contract | docs/technical-verified/api-contracts.md (TDEC-029) |
| Implementation ordering | docs/technical-verified/implementation-slices.md (Slice 1) |
| Frontend dependencies | docs/prd-waves/frontend-pages/page-003.md, page-002.md |
| WAVE-01 contracts | docs/prd-wave-details/waves/wave-01.md |