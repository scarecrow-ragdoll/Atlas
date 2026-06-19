# Traceability

## Technical Requirement Map

No formal TREQ-xxx IDs defined. Technical gaps are tracked as questions (TQ-xxx-xxx). See question-ledger for full mapping.

## Question Map

| Technical Question | Product Source | Scope |
| --- | --- | --- |
| TQ-ARCH-001..006 | docs/product-verified/product-brief.md, scope.md, domain-model.md | architecture-boundaries |
| TQ-DATA-001..010 | docs/product-verified/domain-model.md, scope.md (DEC-007, DEC-009) | data-contracts |
| TQ-API-001..013 | docs/product-verified/functional-spec.md, domain-model.md | api-contracts |
| TQ-AUTH-001..012 | docs/product-verified/actors-and-permissions.md, edge-cases.md, business-rules.md | auth-security-compliance |
| TQ-INT-001..008 | docs/product-verified/functional-spec.md, features/ai-export.md, features/backup-and-restore.md | integrations-events |
| TQ-CLIENT-001..012 | docs/product-verified/user-flows.md, functional-spec.md, features/*.md | client-state-ux |
| TQ-OPS-001..006 | docs/product-verified/product-brief.md (DEC-008), scope.md, features/backup-and-restore.md | operations-observability |
| TQ-TEST-001..007 | docs/product-verified/product-brief.md (DEC-006), acceptance-criteria.md | testing-delivery |

## Decision Map

| Technical Decision | Source | Scope Impact |
| --- | --- | --- |
| DEC-006 | product-verified DEC-006 | testing-delivery |
| DEC-007 | product-verified DEC-007 | data-contracts, api-contracts, auth-security |
| DEC-008 | product-verified DEC-008 | architecture, operations, client-state |
| DEC-009 | product-verified DEC-009 | data-contracts, api-contracts, integrations |

## Source Map

All technical questions trace to docs/product-verified/ as the primary source. No external technical sources were used.

## Slice Map

Slices defined in implementation-slices.md. All blocked by foundational technical decisions.