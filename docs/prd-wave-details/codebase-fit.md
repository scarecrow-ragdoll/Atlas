# Codebase Fit

## Relevant Modules
- apps/api: Go HTTP API — target for all WAVE-07 code
- apps/api/internal/atlas: Atlas fitness module — all new files added here
- apps/api/internal/atlas/models: new UserProfile and AiExport model files
- apps/api/internal/atlas/service: new UserProfileService, AiExportService, AiExportDataProvider
- apps/api/internal/atlas/repository/postgres: new user_profile and ai_export repos + sqlc queries
- apps/api/internal/atlas/graph/schema: new UserProfile and AiExport graphql schemas
- apps/api/internal/atlas/graph/resolver: new UserProfile and AiExport resolvers
- apps/api/cmd/server/main.go: REST handler registration for download and user-profile

No new Nx packages or top-level modules.

## Relevant Files Read
- apps/api/internal/atlas/service/week_flag.go — existing service pattern
- apps/api/internal/atlas/service/settings_service.go — existing settings service with DefaultAiExportWeeks
- apps/api/internal/atlas/service/body_weight.go — generic service pattern
- apps/api/internal/atlas/models/week_flag.go — model pattern (records, types, result unions)
- apps/api/internal/atlas/models/settings.go — settings model with DefaultAiExportWeeks
- apps/api/internal/atlas/repository/postgres/week_flag_repo.go — repo pattern
- apps/api/internal/atlas/repository/postgres/settings_repo.go — repo upsert pattern
- apps/api/internal/atlas/graph/resolver/week_flag.go — resolver pattern
- apps/api/internal/atlas/graph/resolver/resolver.go — resolver container wiring
- apps/api/internal/atlas/graph/schema/week_flag.graphql — graphql schema pattern
- apps/api/internal/atlas/graph/schema/settings.graphql — settings schema
- apps/api/internal/atlas/graph/schema/schema.graphql — root schema
- apps/api/internal/repository/postgres/migrations/00089_week_flags.sql — migration pattern
- apps/api/internal/repository/postgres/migrations/00080_atlas_foundation.sql — atlas_users and atlas_settings tables
- apps/api/cmd/server/main.go — service/repo/route wiring pattern
- apps/api/atlas-gqlgen.yml — gqlgen config with model bindings

## Public Contracts
- UserProfile CRUD (via GraphQL or REST):
  - GetUserProfile: returns user profile with AI context fields
  - UpdateUserProfile: updates goal, height, birthDate, trainingExperience, etc.
  - UserProfile auto-created on first access (bootstrap extension)
- AiExport:
  - POST /api/ai-export/generate: accepts date range, section toggles, optional one-time comment. Returns generatedPrompt. Creates ZIP on filesystem.
  - GET /api/ai-export/download?exportId=: streams ZIP file with auth check
  - GET /api/user-profile: returns goal context for prompt builder
- All operations require PIN auth session when PIN is enabled
- Error format per existing patterns (union result types for GraphQL, JSON error responses for REST)

## Generated Artifact Impact
- gqlgen (atlas-gqlgen.yml): needs new model bindings for UserProfile, AiExport types, inputs, result unions
- sqlc: auto-discovers user_profiles.sql and ai_exports.sql via glob
- No existing generated artifacts affected — all additions are additive

## Integration Points
- PIN auth middleware from WAVE-01: guards all WAVE-07 endpoints
- WAVE-04 WeekFlagService: read week flags for AI context
- WAVE-04 BodyWeightEntryRepo: read body weight data for export
- WAVE-04 BodyCheckInRepo: read check-in/measurement/photo data for export
- WAVE-04 CardioEntryRepo: read cardio data for export
- WAVE-02 ExerciseService: read exercise metadata for export
- WAVE-05 NutritionProductRepo + macro services: read nutrition data for export
- WAVE-03 WorkoutExercise/WorkoutSet repos: read workout data for export (empty arrays when not deployed)
- File system: ZIP storage path (configurable ExportBasePath)
- bootstrap service: extend EnsureDefaultUser to create default UserProfile

## Likely Graph Deltas
- M-API (atlas module) gains: UserProfileService, AiExportService, AiExportDataProvider, ZIP generation utility
- apps/api/internal/atlas/models gains: user_profile.go, ai_export.go
- apps/api/internal/atlas/service gains: user_profile_service.go, ai_export_service.go, ai_export_data_provider.go
- apps/api/internal/atlas/repository/postgres/queries gains: user_profiles.sql, ai_exports.sql
- apps/api/internal/atlas/repository/postgres gains: user_profile_repo.go, ai_export_repo.go
- apps/api/internal/atlas/graph/schema gains: user_profile.graphql, ai_export.graphql
- apps/api/internal/atlas/graph/resolver gains: user_profile.go, ai_export.go
- apps/api/internal/repository/postgres/migrations gains: 00091_user_profiles.sql, 00092_ai_exports.sql
- apps/api/cmd/server/main.go gains: 2-3 service wiring additions, 2 REST route registrations
- apps/api/atlas-gqlgen.yml gains: model bindings for UserProfile, AiExport types

## Unsupported Assumptions
- WAVE-03 workout_sets/daily_log_exercises tables may not exist — export data returns empty arrays for workout section (same stub pattern as DDEC-W06-010)
- WAVE-04 WeekFlag service is fully operational — verified in codebase audit
- Settings.DefaultAiExportWeeks exists — verified in codebase
- ZIP generation uses Go stdlib archive/zip — no external dependencies needed
- File system path for export storage is configurable via AiExportConfig
- sqlc and gqlgen auto-discovery work via glob patterns — verified in prior waves
- No migration number collision — migrations are 00091 and 00092