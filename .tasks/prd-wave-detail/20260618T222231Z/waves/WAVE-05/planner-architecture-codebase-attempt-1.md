# WAVE-05 Architecture-Codebase Planner Attempt 1

## Sources Read
- apps/api/internal/atlas/repository/postgres/settings_repo.go
- apps/api/internal/atlas/service/settings_service.go
- apps/api/internal/atlas/models/settings.go
- apps/api/internal/atlas/graph/resolver/resolver.go
- apps/api/internal/atlas/graph/resolver/settings.go
- apps/api/internal/atlas/graph/schema/settings.graphql
- apps/api/internal/atlas/graph/schema/schema.graphql
- apps/api/cmd/server/main.go (Atlas wiring sections)
- apps/api/atlas-gqlgen.yml
- apps/api/sqlc.yaml
- apps/api/internal/repository/postgres/queries/atlas_settings.sql
- apps/api/internal/repository/postgres/migrations/00080_atlas_foundation.sql
- docs/prd-wave-details/waves/wave-04.md (pattern reference)
- docs/prd-wave-details/waves/wave-01.md (pattern reference)

## Selected Backend Wave Boundary
WAVE-05 adds 5 new entities (NutritionProduct, NutritionTemplate, NutritionTemplateItem, DailyNutritionOverride, DailyNutritionOverrideItem) plus a macro calculation service. All CRUD via GraphQL. No REST endpoints needed (no binary data).

## Neighboring Backend Wave Fit
- WAVE-04 migration sequence: last migration is 00080 (atlas_foundation). WAVE-04 likely adds 00081-00087. WAVE-05 should start at 00088.
- No domain overlap with WAVE-01/02/03/04. Nutrition tables are independent.
- WAVE-06 (Charts) and WAVE-07 (AI Export) consume nutrition data via service layer queries — no direct table dependency.

## Frontend Pages Context
PAGE-007 (Nutrition) depends on all WAVE-05 GraphQL queries/mutations plus macro calculation endpoint.

## Codebase Evidence

### Existing Module Structure
```
apps/api/internal/atlas/
  models/                    # Domain types
    settings.go              # SettingsRecord, Settings, SettingsInput, SettingsResult
    pin.go                   # Pin types
  repository/
    postgres/
      settings_repo.go       # SettingsRepository interface + impl (Pattern B — Atlas-style)
  service/
    settings_service.go      # SettingsService interface + impl
    pin_service.go           # PinService interface + impl
    bootstrap_service.go     # Default user/settings bootstrap
  middleware/
    pin_auth.go              # PIN guard middleware
    user_context.go          # Atlas user context
  graph/
    schema/
      schema.graphql         # Root Query/Mutation extensions
      settings.graphql       # Settings types
      pin.graphql            # PIN types
    resolver/
      resolver.go            # Resolver struct (DI container)
      settings.go            # Settings resolvers
      schema.resolvers.go    # Schema root resolvers
    generated/
      exec.go                # gqlgen generated
      models.go              # gqlgen generated models
```

### Pattern (Pattern B — Atlas-style)
1. **Migration**: goose file in `repository/postgres/migrations/`
2. **sqlc queries**: `.sql` file in `repository/postgres/queries/` with `-- name: Foo :one|:many|:exec`
3. **Repository**: interface + private struct in `atlas/repository/postgres/`, uses sqlc-generated code
4. **Service**: interface + private struct in `atlas/service/`, holds repo dependency
5. **Models**: Go types in `atlas/models/` for both DB records and public/result types
6. **GraphQL schema**: `.graphql` file in `atlas/graph/schema/`
7. **GraphQL resolver**: Go file in `atlas/graph/resolver/` on the shared `Resolver` struct
8. **gqlgen config**: `atlas-gqlgen.yml` — explicit model bindings per type
9. **sqlc config**: `sqlc.yaml` — single project with `queries/` glob
10. **Main wiring**: `cmd/server/main.go` — create repo, service, add to `atlasRes`

### Key Integration Points
- `atlas-gqlgen.yml`: new model bindings for all 5 entities and their result types
- `sqlc.yaml`: auto-discovers new `.sql` files in `queries/` — no config change needed
- `cmd/server/main.go`: add 5 repos, 6 services (5 CRUD + 1 macro calculation), add to `atlasRes`
- `atlas/graph/resolver/resolver.go`: add new service fields to `Resolver` struct

## Generated Artifact Impact
- gqlgen: new schema files → new generated `exec.go` and `models.go` entries
- sqlc: new `.sql` files → new generated `.sql.go` in `generated/`

## Integration Points
- PIN auth middleware (`atlasMiddleware.AtlasPinGuard`) protects `/graphql/atlas` handler — no change needed for new resolvers
- `atlasUserID` from middleware used for all user-scoped operations
- Macro calculation is a stateless service function — no DB dependency beyond reading items

## Likely Graph Deltas
New modules to add to knowledge-graph.xml:
- Models: NutritionProduct, NutritionTemplate, NutritionTemplateItem, DailyNutritionOverride, DailyNutritionOverrideItem, NutiritionMacroResult
- Repository adapters: nutrition_product_repo.go, nutrition_template_repo.go, nutrition_template_item_repo.go, nutrition_override_repo.go, nutrition_override_item_repo.go
- Services: nutrition_product_service.go, nutrition_template_service.go, nutrition_override_service.go, nutrition_macro_service.go
- GraphQL schema: nutrition.graphql
- GraphQL resolvers: nutrition.resolvers.go
- sqlc queries: nutrition_products.sql, nutrition_templates.sql, nutrition_template_items.sql, nutrition_overrides.sql, nutrition_override_items.sql
- Migrations: 00081_nutrition_tables.sql (single migration for all 5 tables)

## Unsupported Assumptions
- Assume sqlc.yaml `queries: "internal/repository/postgres/queries"` glob picks up new `.sql` files — confirmed by config
- Assume gqlgen `schema: - internal/atlas/graph/schema/*.graphql` glob picks up new files — confirmed
- Assume WAVE-01 provides `atlasUserID` in context via `atlasMiddleware.GetAtlasUserID(ctx)` — confirmed by settings resolver pattern
- Template upsert (replace by week) requires custom INSERT ... ON CONFLICT logic in sqlc, not a simple CRUD query

## Proposed Implementation Slices

| Slice ID | Name | Description |
| --- | --- | --- |
| SLICE-W05-001 | DB migrations | Create goose migration 00081_nutrition_tables.sql for all 5 nutrition tables with indexes, FKs, cascades, and constraints |
| SLICE-W05-002 | sqlc queries | Define CRUD queries for all 5 entities: nutrition products (list all, by ID), templates (by week, current), template items (by templateId), overrides (by date, range), override items (by overrideId) |
| SLICE-W05-003 | Repository adapters | Implement 5 repo adapters with sqlc-generated code and error mapping |
| SLICE-W05-004 | Models | Define all 5 entity models (DB records, public models, inputs, result types, error types) |
| SLICE-W05-005 | Services layer | Implement 4 CRUD services (product, template, override, override-item) plus 1 macro calculation service |
| SLICE-W05-006 | GraphQL schema | Add nutrition.graphql with types, enums, inputs, queries, mutations, union results |
| SLICE-W05-007 | GraphQL resolvers | Implement nutrition resolvers with PIN auth guard and union error returns |
| SLICE-W05-008 | gqlgen config and wiring | Add model bindings to atlas-gqlgen.yml, wire repos/services in main.go, add to Resolver struct |

## Risks And Rollback
- Single migration (00081) for all 5 tables: easier to manage than 5 separate migrations. Rollback available.
- Template upsert (ON CONFLICT DO UPDATE) is an upsert pattern — test carefully to avoid data loss.
- No REST endpoints — fewer security surfaces.

## Questions Raised
- DQ-W05-004: Should template and override queries return resolved macro data (calculated server-side) or raw items (calculated frontend)? Recommended: server-side calculation for consistency.

## Traceability Candidates
- apps/api/internal/atlas → all implementation slices reference this module tree
- apps/api/atlas-gqlgen.yml → SLICE-W05-008
- apps/api/sqlc.yaml → SLICE-W05-002
- apps/api/cmd/server/main.go → SLICE-W05-008
- apps/api/internal/repository/postgres/migrations/ → SLICE-W05-001
- apps/api/internal/repository/postgres/queries/ → SLICE-W05-002
- apps/api/internal/repository/postgres/generated/ → SLICE-W05-002 (output)