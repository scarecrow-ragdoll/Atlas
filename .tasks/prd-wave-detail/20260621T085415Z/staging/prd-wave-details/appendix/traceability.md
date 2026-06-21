# Traceability

## Slice Map
| Slice | Source |
| --- | --- |
| SLICE-W06-001 | CAP-W06-001 through CAP-W06-007 — chart models for all chart types |
| SLICE-W06-002 | CAP-W06-004 — body measurement range sqlc query |
| SLICE-W06-003 | CAP-W06-003, CAP-W06-004 — body chart service (weight + measurement) |
| SLICE-W06-004 | CAP-W06-005 — nutrition weekly average service |
| SLICE-W06-005 | CAP-W06-001 through CAP-W06-007 — charts.graphql schema |
| SLICE-W06-006 | CAP-W06-001 through CAP-W06-007 — chart resolvers |
| SLICE-W06-007 | Wiring — all service injections |
| SLICE-W06-008 | CAP-W06-002 — Epley e1RM helper |

## Acceptance Criteria Map
| AC | Source AC | Source Section |
| --- | --- | --- |
| AC-W06-001 | AC-020, AC-065, AC-066 | PRD §16.2 |
| AC-W06-002 | AC-067, RULE-012 | PRD §16.2, §10.7 |
| AC-W06-003 | EDGE-008 | PRD §16 |
| AC-W06-004 | AC-069 | PRD §16.3 |
| AC-W06-005 | EDGE-008 | PRD §16 |
| AC-W06-006 | AC-070 | PRD §16.3 |
| AC-W06-007 | AC-071 | PRD §16.3 |
| AC-W06-008 | EDGE-008 | PRD §16 |
| AC-W06-009 | Derived — empty types edge case | planner-product-ac |
| AC-W06-010 | AC-072 | PRD §16.4 |
| AC-W06-011 | RULE-015 | PRD §16.4 |
| AC-W06-012 | AC-073 | PRD §16.1 |
| AC-W06-013 | Derived — validation | planner-product-ac |
| AC-W06-014 | Derived — PIN auth | TDEC-037 |
| AC-W06-015 | EDGE-008 | PRD §16 |

## Exit Criteria Map
| EC | Source |
| --- | --- |
| EC-W06-001 | Coverage of all ACs |
| EC-W06-002 | WAVE-01 PIN auth pattern |
| EC-W06-003 | Codegen workflow |
| EC-W06-004 | EDGE-008 — empty series behavior |
| EC-W06-005 | RULE-015 — weekly average |
| EC-W06-006 | Lint workflow |
| EC-W06-007 | Privacy policy |
| EC-W06-008 | AC-W06-004 accuracy |
| EC-W06-009 | AC-W06-009 — empty types |
| EC-W06-010 | Conditional WAVE-03 ACs |

## Verification Obligation Map
| Test | AC/EC Coverage |
| --- | --- |
| TEST-W06-001 | AC-W06-001 |
| TEST-W06-002 | AC-W06-002, RULE-012 |
| TEST-W06-003 | AC-W06-003, EC-W06-004 |
| TEST-W06-004 | AC-W06-004, EC-W06-008 |
| TEST-W06-005 | AC-W06-005, EC-W06-004 |
| TEST-W06-006 | AC-W06-006, EC-W06-004 |
| TEST-W06-007 | AC-W06-007 |
| TEST-W06-008 | AC-W06-010, AC-W06-011, EC-W06-005 |
| TEST-W06-009 | AC-W06-015, EC-W06-004 |
| TEST-W06-010 | AC-W06-013 |
| TEST-W06-011 | AC-W06-012, DQ-W06-004 |
| TEST-W06-012 | AC-W06-014, EC-W06-002 |
| TEST-W06-013 | EC-W06-003 |
| TEST-W06-014 | EC-W06-006 |
| TEST-W06-015 | gqlgen schema validation |
| TEST-W06-016 | EC-W06-007 |
| TEST-W06-017 | DQ-W06-005, DDEC-W06-011 |
| TEST-W06-018 | AC-W06-006 measurement side |
| TEST-W06-019 | Body weight single point behavior |
| TEST-W06-020 | Nutrition partial week calculation |
| TEST-W06-021 | AC-W06-009 — empty types |
| TEST-W06-022 | AC-W06-007, DDEC-W06-007 |

## Code Touchpoint Map
- apps/api/internal/atlas/models/chart.go — new file
- apps/api/internal/atlas/service/body_chart_service.go — new file
- apps/api/internal/atlas/service/nutrition_weekly_avg_service.go — new file
- apps/api/internal/atlas/repository/postgres/queries/body_measurements_range.sql — new sqlc query
- apps/api/internal/atlas/repository/postgres/body_measurement_repo.go — add ListByUserTypeRange
- apps/api/internal/atlas/graph/schema/charts.graphql — new schema
- apps/api/internal/atlas/graph/resolver/charts.go — new resolver
- apps/api/internal/atlas/graph/resolver/resolver.go — add services
- apps/api/cmd/server/main.go — wire services
- apps/api/atlas-gqlgen.yml — add model bindings

## Question Map
| ID | Wave | Scope | Severity | Status | Resolution |
| --- | --- | --- | --- | --- | --- |
| Q-CHART-001 | 06 | operations | medium | resolved | Epley formula selected |
| DQ-W06-001 | 06 | product-ac | resolved | resolved | Highest e1RM per session |
| DQ-W06-002 | 06 | architecture | resolved | resolved | Stubs returning empty |
| DQ-W06-003 | 06 | product-ac | resolved | resolved | Per-session snapshot |
| DQ-W06-004 | 06 | data-ops | resolved | resolved | 4 weeks default |
| DQ-W06-005 | 06 | data-ops | deferred | open | 52-week max |
| DQ-W06-006 | 06 | data-ops | resolved | resolved | Stubs returning empty |

## Source Map
- docs/prd-waves/waves/wave-06.md — source wave boundary
- docs/product-verified/features/charts.md — chart feature spec
- docs/product-verified/acceptance-criteria.md — AC-020–AC-022, AC-065–AC-073
- docs/product-verified/business-rules.md — RULE-012 through RULE-015
- docs/product-verified/edge-cases.md — EDGE-008, EDGE-026
- docs/prd-waves/frontend-pages/page-008.md — PAGE-008 backend dependencies
- docs/prd-wave-details/waves/wave-04.md — WAVE-04 body tracking data contracts
- docs/prd-wave-details/waves/wave-05.md — WAVE-05 nutrition macro service contract
- apps/api/internal/atlas — codebase patterns