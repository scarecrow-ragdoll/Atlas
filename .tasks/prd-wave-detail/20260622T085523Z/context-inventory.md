# Context Inventory: WAVE-09

## PRD Wave Sources
| File | Status |
|------|--------|
| docs/prd-waves/waves/wave-09.md | user-approved source wave |
| docs/prd-waves/wave-map.md | user-approved wave map |
| docs/prd-waves/index.md | wave package status |
| docs/prd-waves/scope-inventory.md | capability groups |
| docs/prd-waves/source-inventory.md | source list |
| docs/prd-waves/open-questions.md | no WAVE-09 questions |

## Frontend Pages Source
| File | Purpose |
|------|---------|
| docs/prd-waves/frontend-pages/index.md | PAGE-010 Import/Export page listing |
| docs/prd-waves/frontend-pages/page-010.md | PAGE-010 backend dependencies |

PAGE-010 backend dependencies:
- POST /api/backup/export
- POST /api/backup/import (dry-run then confirm)

## Product Sources
| File | Content |
|------|---------|
| docs/product-verified/features/backup-and-restore.md | Feature spec: full backup ZIP, import flow |
| docs/product-verified/acceptance-criteria.md | AC-093 through AC-102, AC-114-116, AC-124 |
| docs/product-verified/domain-model.md | All 20 entities for backup |
| docs/product-verified/functional-spec.md §20 | REQ-016: backup/import spec |
| docs/product-verified/product-brief.md | Performance targets, success metrics |
| docs/product-verified/business-rules.md | RULE-007: ZIP schema validation |
| docs/product-verified/user-flows.md | Export flow (§26.11), import flow (§26.12) |
| docs/product-verified/edge-cases.md | EDGE-010 (invalid ZIP), EDGE-021 (partial restore), EDGE-028 (schema migration) |
| docs/product-verified/actors-and-permissions.md | Export/import permissions |

## Technical Sources
None. docs/technical-verified not created (not required for shallow waves).

## GRACE Sources
| File | Purpose |
|------|---------|
| docs/development-plan.xml | Module definitions, M-API patterns |
| docs/knowledge-graph.xml | Module graph, cross-references |
| docs/verification-plan.xml | Verification references |

## Codebase Sources
Existing WAVE-07 (AI Export) pattern is the closest reference for ZIP generation, file storage, and data aggregation.

### Relevant Existing Modules
| Module | Purpose |
|--------|---------|
| AiExport service | ZIP creation, file storage, export metadata CRUD via GraphQL + REST |
| AiReview service | WAVE-08 exposes ListAllByUserID for backup consumption |
| All other entity services | Must provide ListAllByUserID for backup data aggregation |
| export_zip.go | Reusable ZIP build infrastructure (Manifest, ExportArchive, BuildZIP) |
| atlas_media.go | File storage pattern (directory per entity, temp files, atomic rename) |
| settings/goal/profile services | Foundation entities for backup |

### Key Code Locations
| Path | Purpose |
|------|---------|
| apps/api/internal/atlas/service/ai_export_service.go | Service pattern, ZIP handling |
| apps/api/internal/atlas/service/export_zip.go | Reusable ExportArchive + BuildZIP |
| apps/api/internal/atlas/models/ai_export.go | Record + Public + Input + Result patterns |
| apps/api/internal/atlas/graph/resolver/resolver.go | Resolver struct (must add BackupService) |
| apps/api/cmd/server/main.go | Wiring pattern |
| apps/api/internal/handler/atlas_media.go | File download pattern |
| apps/api/internal/handler/ai_export_handler.go | REST handler pattern for ZIP download |
| apps/api/internal/atlas/repository/postgres/queries/ | sqlc query directory |
| apps/api/internal/repository/postgres/migrations/ | Migration directory (next: 00094) |
| apps/api/atlas-gqlgen.yml | gqlgen bindings for atlas modules |

### All Entity Services (WAVE-01 through WAVE-08)
Backup must read ALL entities. Each service needs a ListAllByUserID or equivalent reader:
- SettingsService
- UserProfileService
- ExerciseService + ExerciseMediaService
- CardioService
- BodyWeightService / BodyCheckInService / BodyMeasurementService
- NutritionProductService / NutritionTemplateService / NutritionOverrideService
- WeekFlagService
- AiExportService
- AiReviewService

## Source Delta
No source delta — WAVE-09 is the ninth wave and all prior waves are completed.

## Source Gaps
1. Q-ACTOR-08: Import behavior when data already exists (merge/replace/error)
2. Q-AC-15: Import with existing data behavior
3. Q-AC-16: CSV files — mandatory or optional in backup
4. Q-EDGE-11: Migration strategy for schema version changes
5. No docs/technical-verified exists — technical decisions must be derived from product docs and codebase patterns

## Prior Detailed Waves Check
| Wave | Status | WAVE-09 Dependency |
|------|--------|-------------------|
| WAVE-01 | user-approved | Foundation, PIN auth, atlas_users |
| WAVE-02 | user-approved | Exercises + media for backup |
| WAVE-03 | user-approved | Workout diary (daily logs, workout exercises, sets) |
| WAVE-04 | user-approved | Cardio, body check-ins, measurements |
| WAVE-05 | user-approved | Nutrition products, templates, overrides |
| WAVE-06 | user-approved | Charts (no direct backup dependency) |
| WAVE-07 | ready-for-dev | AI export records + UserProfile |
| WAVE-08 | user-approved | AiReview records (ListAllByUserID exposed) |