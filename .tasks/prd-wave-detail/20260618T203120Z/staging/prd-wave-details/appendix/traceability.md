# Traceability

## Slice Map
| Slice ID | Source | Wave |
| --- | --- | --- |
| SLICE-W01-001 | docs/technical-verified/implementation-slices.md (Slice 0) | WAVE-01 |
| SLICE-W01-002 | docs/technical-verified/api-contracts.md (TDEC-001) | WAVE-01 |
| SLICE-W01-003 | docs/technical-verified/auth-security-compliance.md | WAVE-01 |
| SLICE-W01-004 | docs/technical-verified/architecture-and-boundaries.md | WAVE-01 |
| SLICE-W01-005 | docs/product-verified/functional-spec.md | WAVE-01 |
| SLICE-W01-006 | docs/product-verified/domain-model.md | WAVE-01 |
| SLICE-W01-007 | docs/technical-verified/api-contracts.md (REST endpoints) | WAVE-01 |
| SLICE-W01-008 | apps/api/internal/appconfig/config.go pattern | WAVE-01 |
| SLICE-W01-009 | apps/api/gqlgen.yml pattern | WAVE-01 |
| SLICE-W01-010 | docker-compose.yml | WAVE-01 |
| SLICE-W01-011 | apps/api/internal/testinfra pattern | WAVE-01 |

## Acceptance Criteria Map
| AC ID | Source | Verified By |
| --- | --- | --- |
| AC-W01-001 | docs/product-verified/domain-model.md#Settings | TEST-W01-001, TEST-W01-007 |
| AC-W01-002 | docs/product-verified/functional-spec.md (PIN section) | TEST-W01-002 |
| AC-W01-003 | docs/product-verified/functional-spec.md (PIN change) | TEST-W01-002 |
| AC-W01-004 | docs/product-verified/functional-spec.md (PIN disable) | TEST-W01-002 |
| AC-W01-005 | docs/technical-verified/auth-security-compliance.md (session) | TEST-W01-003 |
| AC-W01-006 | docs/technical-verified/architecture-and-boundaries.md (auth boundary) | TEST-W01-004 |
| AC-W01-007 | docs/product-verified/functional-spec.md (settings read) | TEST-W01-005 |
| AC-W01-008 | docs/product-verified/functional-spec.md (settings update) | TEST-W01-005 |
| AC-W01-009 | docs/technical-verified/api-contracts.md (media upload) | TEST-W01-006 |
| AC-W01-010 | docs/technical-verified/api-contracts.md (media download) | TEST-W01-006 |
| AC-W01-011 | docs/technical-verified/api-contracts.md (not found) | TEST-W01-006 |
| AC-W01-012 | docs/development-plan.xml (M-API contract) | TEST-W01-008 |
| AC-W01-013 | apps/api/internal/handler/health.go | TEST-W01-009 |
| AC-W01-014 | docs/technical-verified/architecture-and-boundaries.md | TEST-W01-004 |

## Exit Criteria Map
| EC ID | Validated By |
| --- | --- |
| EC-W01-001 | TEST-W01-001 through TEST-W01-009 |
| EC-W01-002 | TEST-W01-011, TEST-W01-012 |
| EC-W01-003 | TEST-W01-008 |
| EC-W01-004 | TEST-W01-009 |
| EC-W01-005 | docker compose up smoke test |
| EC-W01-006 | Config validation in TEST-W01-002 |
| EC-W01-007 | TEST-W01-010 |
| EC-W01-008 | bunx nx run api:typecheck |
| EC-W01-009 | No frontend files in diff |

## Verification Obligation Map
| TEST ID | Type | Coverage |
| --- | --- | --- |
| TEST-W01-001 | unit | AC-W01-001 |
| TEST-W01-002 | unit | AC-W01-002 through AC-W01-004 |
| TEST-W01-003 | unit | AC-W01-005 |
| TEST-W01-004 | unit | AC-W01-006, AC-W01-014 |
| TEST-W01-005 | integration | AC-W01-007, AC-W01-008 |
| TEST-W01-006 | integration | AC-W01-009 through AC-W01-011 |
| TEST-W01-007 | integration | EC-W01-001 |
| TEST-W01-008 | unit | AC-W01-012, EC-W01-003 |
| TEST-W01-009 | unit | AC-W01-013, EC-W01-004 |
| TEST-W01-010 | lint | EC-W01-007 |
| TEST-W01-011 | codegen | EC-W01-002 |
| TEST-W01-012 | codegen | EC-W01-002 |

## Code Touchpoint Map
| Code Path | WAVE-01 Impact |
| --- | --- |
| apps/api/cmd/server/main.go | Add fitness-domain wiring |
| apps/api/internal/appconfig/config.go | Add Session/Pin/Media config |
| apps/api/internal/middleware/ | Add pin_auth.go |
| apps/api/internal/service/ | Add pin_service.go, settings_service.go |
| apps/api/internal/repository/ | Add pin_session_store.go, settings_sqlc |
| apps/api/internal/handler/ | Add media_handler.go |
| apps/api/internal/graph/ | Add settings_resolver.go |
| libs/graphql/schema/ | Add fitness.graphql, settings.graphql |
| apps/api/gqlgen.yml | Add fitness schema paths |
| apps/api/sqlc.yaml | Add fitness query paths |
| apps/api/migrations/ | Add fitness migration files |
| docker-compose.yml | Add media volume, worker service |

## Question Map
| Question ID | Source | Wave |
| --- | --- | --- |
| DQ-W01-001 | docs/prd-waves/open-questions.md (Q-PIN-001) | WAVE-01 |
| DQ-W01-002 | docs/technical-verified/auth-security-compliance.md | WAVE-01 |

## Source Map
| File | Role |
| --- | --- |
| docs/prd-waves/waves/wave-01.md | Source wave boundary |
| docs/product-verified/functional-spec.md | Product behavior |
| docs/product-verified/domain-model.md | Entity contracts |
| docs/product-verified/actors-and-permissions.md | User roles |
| docs/technical-verified/api-contracts.md | API protocol decisions |
| docs/technical-verified/architecture-and-boundaries.md | System architecture |
| docs/technical-verified/auth-security-compliance.md | Auth and security |
| docs/technical-verified/implementation-slices.md | Slice mapping |
| docs/development-plan.xml | Module contracts |
| docs/knowledge-graph.xml | Module graph |
| apps/api/internal/ | Codebase patterns |