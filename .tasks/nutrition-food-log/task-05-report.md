<!-- FILE: .tasks/nutrition-food-log/task-05-report.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Record Task 5 GraphQL schema/resolver/codegen evidence for the nutrition food log rollout. -->
<!--   SCOPE: RED/GREEN commands, implementation summary, files changed, deviations, and known blockers for Task 5 only. -->
<!--   DEPENDS: docs/superpowers/plans/2026-06-24-nutrition-food-log.md, apps/api Atlas GraphQL files. -->
<!--   LINKS: M-API-NUTRITION / V-M-API-NUTRITION. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->

# Task 5 Report: GraphQL Schema, Resolvers, Codegen, And Wiring

## Status

Implemented.

## RED Evidence

Command:

```bash
cd apps/api && go test ./internal/atlas/graph/resolver -run "TestDailyNutritionGraphQL|TestNutritionGraphQL" -count=1
```

Expected failures were captured before implementation:

- unknown `AddDailyNutritionEntryInput`
- unknown `addDailyNutritionEntry`
- unknown `nutritionProductsAll`
- unknown `NutritionTemplateApplyMode`
- unknown `applyNutritionTemplateToWeek`
- existing `nutritionProducts` executable schema panic stub

## Implementation Summary

- Added factual daily nutrition GraphQL schema for daily logs, entry snapshots, entry mutations, notes mutation, and typed result errors.
- Added executable GraphQL tests that exercise gqlgen root query/mutation wiring instead of direct resolver calls only.
- Added nutrition product all-list and restore GraphQL operations.
- Added weekly nutrition template apply GraphQL operation with `SEED_EMPTY_DAYS` enum mapping.
- Wired nutrition root query/mutation stubs to service-backed resolver methods so existing nutrition fields no longer panic at runtime.
- Added daily nutrition resolver methods and error mapping.
- Added Atlas gqlgen model bindings for daily nutrition and template apply models.
- Wired daily nutrition log service and template apply service into `cmd/server`.
- Regenerated Atlas GraphQL executable schema and resolver helpers.
- Updated GRACE knowledge graph and verification plan for the new GraphQL coverage.

## Intentional API Shape Note

`updateDailyNutritionLogNotes` is exposed as:

```graphql
updateDailyNutritionLogNotes(id: ID!, input: UpdateDailyNutritionLogNotesInput!): DailyNutritionLogResult!
```

This differs from the early plan sketch that put `date` in the input. The current `DailyNutritionLogService.UpdateNotes` contract updates by daily log ID, so the GraphQL API follows the existing service boundary and avoids inventing a separate date-to-id mutation path in this task.

`UpdateDailyNutritionEntryInput.position` is required. The service treats entry updates as full replacements, so GraphQL does not silently default an omitted position to zero.

## GREEN Evidence

Commands:

```bash
bunx nx run api:codegen
bunx nx run api:codegen:atlas
cd apps/api && go test ./internal/atlas/graph/resolver -run "TestDailyNutritionGraphQL|TestNutritionGraphQL" -count=1
cd apps/api && go test ./cmd/server -run '^$' -count=1
bunx nx build api
xmllint --noout docs/knowledge-graph.xml docs/verification-plan.xml
```

Results:

- `api:codegen`: passed
- `api:codegen:atlas`: passed
- focused GraphQL resolver tests: passed
- `cmd/server` compile test: passed
- API build: passed
- XML validation for updated docs: passed

## Known Gaps

- No backend schema changes were made in this task.
- Daily nutrition legacy compatibility remains owned by Task 6.
- Frontend API/page wiring remains owned by Tasks 7-10.
