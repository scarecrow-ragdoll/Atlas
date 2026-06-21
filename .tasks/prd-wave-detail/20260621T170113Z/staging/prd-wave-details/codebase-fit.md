# Codebase Fit

## Relevant Modules
- **M-API**: Shared models (Date, UUID, auth middleware), graphql, sqlc, handler patterns
- **M-API-SETTINGS**: Settings model, repo, service, resolver, schema (defaultAiExportWeeks referenced)
- **M-API-WEEK-FLAG**: WeekFlag model, repo, service, resolver, schema (read-only dependency via AiExportDataProvider)
- **M-API-BODY**: CardioEntry, BodyWeightEntry, BodyCheckIn, BodyMeasurement, ProgressPhoto (read-only)
- **M-API-NUTRITION**: NutritionProduct, NutritionTemplate, NutritionMacroService (read-only)
- **M-API-USER-PROFILE**: NEW — UserProfile model, repo, service, resolver, schema
- **M-API-AI-EXPORT**: NEW — AiExport model, repo, service, resolver, schema, ZIP generator, download handler

## Relevant Files Read
- apps/api/internal/atlas/models/week_flag.go — model pattern (record, public, input, result types)
- apps/api/internal/atlas/models/settings.go — model pattern with DefaultAiExportWeeks
- apps/api/internal/atlas/service/week_flag.go — service pattern (interface, sentinel errors, FromRecord)
- apps/api/internal/atlas/service/settings_service.go — service pattern
- apps/api/internal/atlas/service/nutrition_macro_service.go — aggregate service pattern
- apps/api/internal/atlas/repository/postgres/week_flag_repo.go — repo pattern (interface, struct, New*, *FromRow)
- apps/api/internal/atlas/repository/postgres/settings_repo.go — repo pattern with nullable helpers
- apps/api/internal/atlas/repository/postgres/cardio_entry_repo.go — repo pattern
- apps/api/internal/atlas/graph/schema/week_flag.graphql — GraphQL schema pattern (result types, inline errors, error code enum)
- apps/api/internal/atlas/graph/schema/settings.graphql — GraphQL schema pattern
- apps/api/internal/atlas/graph/resolver/week_flag.go — resolver pattern (auth middleare, error mapping, result types)
- apps/api/internal/atlas/graph/resolver/resolver.go — resolver container with service field injection
- apps/api/internal/atlas/handler/progress_photo_handler.go — binary download handler pattern
- apps/api/internal/repository/postgres/migrations/00080_atlas_foundation.sql — atlas_users, settings tables
- apps/api/internal/repository/postgres/migrations/00089_week_flags.sql — week_flags table
- apps/api/internal/repository/postgres/migrations/00090_nutrition_tables.sql — latest migration
- apps/api/cmd/server/main.go — main wiring pattern
- apps/api/atlas-gqlgen.yml — gqlgen config (needs new bindings)
- apps/api/internal/appconfig/config.go — config struct pattern

## Public Contracts
- **WeekFlagService**: weekFlags(weekStartDate) query — WAVE-07 consumes via WeekFlagService for prompt/data.json
- **Settings**: Settings query with defaultAiExportWeeks — WAVE-07 reads for date range defaults
- **atlas_users**: User identity — WAVE-07 uses user_id FK for both tables

## Generated Artifact Impact
- sqlc: New query files (user_profiles.sql, ai_exports.sql) require codegen regeneration
- gqlgen: 16 new type bindings added to atlas-gqlgen.yml + 2 new .graphql schema files require codegen
- Both generated outputs excluded from coverage allowlist with replacement gates (codegen + build + integration tests)

## Integration Points
- **PIN auth middleware**: All 3 new endpoints under existing PIN-guarded route group
- **WAVE-01 bootstrap**: Extend EnsureDefaultUser to create default UserProfile row
- **WAVE-03 workout data**: AiExportDataProvider.GetWorkoutSummary returns empty when tables absent
- **WAVE-04 data**: AiExportDataProvider reads CardioEntry, BodyWeightEntry, BodyCheckIn, BodyMeasurement, ProgressPhoto, WeekFlag
- **WAVE-05 nutrition**: AiExportDataProvider reads NutritionMacroService for KJBJU averages
- **Media storage**: Photo files referenced via existing media storage (ProgressPhoto.filePath) — copied, not moved

## Likely Graph Deltas
```
M-API-USER-PROFILE:
  depends: M-API (shared models, date, auth)
  files: models/user_profile.go, service/user_profile_service.go,
         repository/postgres/user_profile_repo.go, graph/schema/user_profile.graphql,
         graph/resolver/user_profile.go, repository/postgres/queries/user_profiles.sql,
         repository/postgres/migrations/00091_user_profiles.sql

M-API-AI-EXPORT:
  depends: M-API, M-API-USER-PROFILE, M-API-WEEK-FLAG, M-API-BODY, M-API-NUTRITION
  files: models/ai_export.go, service/ai_export_service.go, service/export_zip.go,
         repository/postgres/ai_export_repo.go, graph/schema/ai_export.graphql,
         graph/resolver/ai_export.go, handler/ai_export_handler.go,
         handler/user_profile_handler.go, repository/postgres/queries/ai_exports.sql,
         repository/postgres/migrations/00092_ai_exports.sql
```

## Unsupported Assumptions
- No technical-verified docs exist — API contracts derived from product-verified docs + codebase patterns
- No prior temp-file-atomic-rename pattern in codebase — new utility needed
- No build-time app version injection — omitted from manifest.json for MVP
- No existing weekFlagsByDateRange query — client calls per week for MVP
- No existing concurrent generation guard — acceptable for MVP with single user