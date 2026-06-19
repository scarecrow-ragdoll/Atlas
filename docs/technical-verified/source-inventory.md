# Source Inventory

## Included Sources

- docs/product-verified/index.md
- docs/product-verified/product-brief.md
- docs/product-verified/scope.md
- docs/product-verified/actors-and-permissions.md
- docs/product-verified/domain-model.md
- docs/product-verified/functional-spec.md
- docs/product-verified/user-flows.md
- docs/product-verified/business-rules.md
- docs/product-verified/edge-cases.md
- docs/product-verified/acceptance-criteria.md
- docs/product-verified/open-questions.md
- docs/product-verified/source-inventory.md
- docs/product-verified/appendix/decision-log.md
- docs/product-verified/features/*.md (13 feature files)

## Source Delta

Technical run 20260618T185935Z includes source delta with 4 product owner decisions:
- DEC-006: Success metrics and quality gates (testing-delivery)
- DEC-007: Multi-user-ready data model with userId (data-contracts, api-contracts, auth-security)
- DEC-008: Performance p95 targets (architecture, operations, client-state)
- DEC-009: DailyLog replaces WorkoutDay (data-contracts, api-contracts)

All decisions reviewed by affected scopes.

## Answered Questions

4 resolved product questions carried forward from product-verified run (Q-SCOPE-001..005).

No technical questions were answered in this run — first pass identified ~58 new technical questions.

## Excluded Or Noisy Sources

None.

## Coverage Gaps

- No API contract or GraphQL schema exists
- No system/component architecture diagrams
- No deployment environment specification
- No test data fixtures or seed strategies
- No UI state machine definitions
- No monitoring/logging framework specification