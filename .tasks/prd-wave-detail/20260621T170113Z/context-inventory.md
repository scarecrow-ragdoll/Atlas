# Context Inventory

## PRD Wave Sources
- docs/prd-waves/index.md — waves-approved
- docs/prd-waves/wave-map.md — waves-approved-by-user (2026-06-18)
- docs/prd-waves/source-inventory.md
- docs/prd-waves/scope-inventory.md
- docs/prd-waves/open-questions.md
- docs/prd-waves/waves/index.md
- docs/prd-waves/waves/wave-07.md — user-approved

## Frontend Pages Source (dependency context only)
- docs/prd-waves/frontend-pages/index.md
- docs/prd-waves/frontend-pages/page-009.md — AI Export page

## Product Sources
- docs/product-verified/index.md
- docs/product-verified/product-brief.md
- docs/product-verified/scope.md
- docs/product-verified/actors-and-permissions.md
- docs/product-verified/domain-model.md
- docs/product-verified/functional-spec.md (§17-18 AI Export)
- docs/product-verified/user-flows.md
- docs/product-verified/business-rules.md (RULE-021, RULE-025, RULE-026, RULE-027)
- docs/product-verified/edge-cases.md (EDGE-008, EDGE-024)
- docs/product-verified/acceptance-criteria.md (AC-023, AC-024)
- docs/product-verified/open-questions.md

## Technical Sources
- SOURCE_MISSING: docs/technical-verified (not created)

## Prior Detailed Waves
- WAVE-01 (Foundation): ready-for-dev (awaiting user approval)
- WAVE-02 (Exercise Library): user-approved
- WAVE-04 (Cardio and Body Tracking): questions-open (DQ-W04-001)
- WAVE-05 (Nutrition): ready-for-dev (awaiting user approval)
- WAVE-06 (Charts): ready-for-dev (awaiting user approval)

## GRACE Sources
- docs/development-plan.xml — M-PRD-WAVE-DETAILER at order 5.175
- docs/knowledge-graph.xml — module boundaries for Atlas services
- docs/verification-plan.xml — test patterns

## Codebase Sources (read-only)
- apps/api/internal/atlas/models/ — model types
- apps/api/internal/atlas/models/week_flag.go — WeekFlag model (exists)
- apps/api/internal/atlas/models/settings.go — Settings with DefaultAiExportWeeks
- apps/api/internal/atlas/service/week_flag.go — WeekFlag service (exists)
- apps/api/internal/atlas/service/settings_service.go — Settings service (exists)
- apps/api/internal/atlas/repository/postgres/week_flag_repo.go — WeekFlag repo (exists)
- apps/api/internal/atlas/repository/postgres/settings_repo.go — Settings repo (exists)
- apps/api/internal/atlas/graph/schema/week_flag.graphql — WeekFlag schema (exists)
- apps/api/internal/atlas/graph/schema/settings.graphql — Settings schema with defaultAiExportWeeks
- apps/api/internal/atlas/graph/resolver/week_flag.go — WeekFlag resolver (exists)
- apps/api/internal/atlas/graph/resolver/settings.go — Settings resolver (exists)
- apps/api/internal/repository/postgres/migrations/00089_week_flags.sql — WeekFlag table
- apps/api/internal/repository/postgres/migrations/00080_atlas_foundation.sql — atlas_users, atlas_settings
- apps/api/internal/atlas/service/ (27 files) — service patterns
- apps/api/internal/atlas/repository/postgres/ (15 files) — repository patterns
- apps/api/internal/atlas/graph/resolver/resolver.go — resolver container wiring
- apps/api/internal/atlas/graph/schema/schema.graphql — root schema

## Key Codebase Findings
- WeekFlag: FULLY IMPLEMENTED (model, service, test, repo, resolver, schema, gqlgen-generated resolver, migration 00089)
- UserProfile: NOT IMPLEMENTED as separate service/repo — only `atlas_users` table with id, display_name
  - UserProfile domain model has: goal, height, birthDate, trainingExperience, etc.
  - Need new migration + model + service + repo + resolver + schema
- AiExport: NOT IMPLEMENTED — only DefaultAiExportWeeks in Settings
  - Need: model, service, repo, resolver, schema, migration
  - Need ZIP generation logic (manifest.json, data.json, summary.md, CSVs, photos)
- Settings model exists with DefaultAiExportWeeks field

## Source Delta
- WAVE-07 is new — no prior detailed wave to compare

## Source Gaps
- docs/technical-verified absent — no API contract, auth, or data contract docs from technical verification
- UserProfile service/repo/model not yet implemented in codebase
- AiExport service/repo/model not yet implemented in codebase
- ZIP export generation not yet implemented
