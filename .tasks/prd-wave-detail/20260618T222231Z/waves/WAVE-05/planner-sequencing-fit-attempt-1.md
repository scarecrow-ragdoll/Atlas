# WAVE-05 Sequencing-Fit Planner Attempt 1

## Sources Read
- docs/prd-waves/wave-map.md
- docs/prd-waves/waves/wave-05.md
- docs/prd-waves/frontend-pages/page-007.md
- docs/prd-wave-details/waves/wave-01.md
- docs/prd-wave-details/waves/wave-02.md
- docs/prd-wave-details/waves/wave-04.md
- docs/development-plan.xml
- docs/knowledge-graph.xml

## Selected Backend Wave Boundary
WAVE-05 is the Nutrition domain. Independent domain with its own tables, services, resolvers, and GraphQL schema file. No cross-domain data references.

## Neighboring Backend Wave Fit

### WAVE-01 (Foundation) — PREREQUISITE
- Provides PIN auth middleware (`atlasMiddleware.AtlasPinGuard`)
- Provides Atlas user context middleware (`atlasMiddleware.AtlasUserContext`)
- Provides `/graphql/atlas` PIN-protected endpoint where WAVE-05 resolvers mount
- Provides gqlgen config (`atlas-gqlgen.yml`) that WAVE-05 extends with new model bindings
- Provides sqlc config that auto-discovers new query files
- Provides `atlas_users` table for user_id FK
- Provides `atlas_bootstrap_service` for ensuring default user exists
- **WAVE-05 cannot start until WAVE-01 provides these contracts**

### WAVE-02 (Exercise Library) — NO DEPENDENCY
- Independent domain. WAVE-05 does not reference exercises or media.

### WAVE-03 (Workout Diary) — NO DEPENDENCY
- Independent domain. WAVE-05 does not reference workouts, DailyLog, or sets.
- Note: WAVE-03 creates `daily_log` table. WAVE-05 does NOT reference `daily_log` — nutrition uses its own date schema.
- **Fully parallelizable with WAVE-03**

### WAVE-04 (Cardio and Body Tracking) — NO DIRECT DEPENDENCY
- Independent domain. WAVE-05 does not reference cardio, body check-ins, measurements, or photos.
- **Fully parallelizable with WAVE-04** (confirmed by wave-map.md)

### WAVE-06 (Charts) — DOWNSTREAM CONSUMER
- WAVE-06 (Charts) consumes WAVE-05 data for weekly KJBJU average charts (AC-072)
- WAVE-06 depends on WAVE-05 providing a macro query or data export endpoint
- No contract change needed in WAVE-05 — the `nutritionMacros` query provides all data WAVE-06 needs

### WAVE-07 (AI Export) — DOWNSTREAM CONSUMER
- AI export includes nutrition data (template and override items) in export ZIP
- WAVE-05 must provide a service-layer method to export nutrition data grouped by week
- This is a read-only service contract, not a schema change. WAVE-07 calls the existing service
- Deferred: WAVE-05 exports a JSON-serializable representation of templates and overrides

### WAVE-08 (AI Review) — NO DEPENDENCY
- AI Review stores only text responses. No nutrition data dependency.

### WAVE-09 (Backup) — DOWNSTREAM CONSUMER
- Backup includes nutrition tables in export
- WAVE-05 tables are JSON-serializable (no binary data, no file paths)
- WAVE-09 can serialize all WAVE-05 entities via service layer or direct DB queries

## Frontend Pages Context

### PAGE-007 (Nutrition) — PRIMARY CONSUMER
- Products list/CRUD: depends on `nutritionProducts` and `createNutritionProduct`/`updateNutritionProduct`/`deleteNutritionProduct`
- Weekly template editor: depends on `nutritionTemplateCurrent` (by weekStartDate) and `createNutritionTemplate`/`updateNutritionTemplate`/`deleteNutritionTemplate` plus nested item mutations
- Daily override editor: depends on `dailyNutritionOverrideByDate` and `createDailyNutritionOverride`/`updateDailyNutritionOverride`/`deleteDailyNutritionOverride` plus nested item mutations
- Macro summary: depends on `nutritionMacros(weekStartDate, date)` query returning all 4 macros
- Empty states: no products → empty list, no template → no active template message
- Backend dependency context only — no frontend work in this wave
- PAGE-007 uses REST-style URL structure (e.g., `/api/nutrition-products`) in its backend dependencies but expects these to be GraphQL operations, not REST. No change needed — the frontend will use GraphQL queries.

## Dependency Order
WAVE-01 → WAVE-05 (WAVE-01 must ship first). WAVE-05 ↔ WAVE-02/03/04 (fully parallelizable). WAVE-05 → WAVE-06/07/09 (WAVE-05 provides data contracts).

## Scope Collision Check
- **No collision with WAVE-02**: Exercise Library does not touch nutrition
- **No collision with WAVE-03**: Workout Diary does not touch nutrition
- **No collision with WAVE-04**: Cardio and Body Tracking do not touch nutrition
- **Migration number check**: If WAVE-04 uses 00081-00087 (6 tables for 6 entities), WAVE-05 must use 00088+. If WAVE-04 hasn't claimed migration numbers, WAVE-05 should coordinate or use a designated block. **Recommendation**: WAVE-05 migration = 00081 if no other wave has claimed it. Coordinate during implementation.

## Risks And Rollback
- Migration number collision with WAVE-04 is the only coordination risk. Both waves can parallelize on source, but migration numbers must not conflict.
- No other sequencing blockers identified.

## Questions Raised
- DQ-W05-009: What migration number should WAVE-05 use? Depends on WAVE-04's migration range. WAVE-05 should check current migration state at implementation time.

## Traceability Candidates
- docs/prd-waves/wave-map.md → dependency order, parallelization note
- docs/prd-wave-details/waves/wave-01.md → WAVE-01 contracts
- docs/prd-wave-details/waves/wave-04.md → parallelization confirmation
- docs/prd-waves/frontend-pages/page-007.md → frontend dependency context