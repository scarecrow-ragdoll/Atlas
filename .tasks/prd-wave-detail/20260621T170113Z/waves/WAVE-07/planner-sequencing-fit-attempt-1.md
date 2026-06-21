# WAVE-07 Sequencing and Fit Report

**Run:** 20260621T170113Z  
**Wave:** WAVE-07 — AI Export and Prompt Builder  
**Role:** sequencing-fit planner  
**Attempt:** 1  

---

## 1. Prior Wave Compatibility

### 1.1 WAVE-01 (Foundation) — READY
WAVE-01 Settings table (AC-W01-001) stores `default_export_weeks`, `ai_goal`, `ai_height`, `ai_age`, `ai_experience`, `ai_split`, `ai_limits`, `ai_progression`, `ai_nutrition_strategy`. This directly serves WAVE-07's CAP-W07-001 (persistent AI context) and CAP-W07-002 (user goal storage). **Compatible — WAVE-07 reads Settings for AI context and goal.**

### 1.2 WAVE-01 / UserProfile gap
PAGE-009 requires `GET /api/user-profile` for goal context. WAVE-01 Settings contains AI context fields but the domain model separates UserProfile as a distinct entity. WAVE-01 does not define a UserProfile table (only Settings + pin_sessions). **Gap: WAVE-07 must either (a) add a UserProfile migration/endpoint, or (b) surface relevant settings fields through the existing Settings GraphQL resolvers. Recommend (b): query Settings via existing GraphQL `settings` query rather than adding a new REST endpoint.**

### 1.3 WAVE-02 (Exercise Library) — READY
WAVE-02 provides exercise metadata (name, muscleGroups, workingWeight) and media for export. WAVE-07 consumes via service layer (read-only). **Compatible — no collision.**

### 1.4 WAVE-04 (Cardio and Body Tracking) — questions-open (DQ-W04-001)
WAVE-04 provides CardioEntry, BodyWeightEntry, BodyCheckIn, BodyMeasurement, ProgressPhoto, and **WeekFlag** CRUD (CAP-W04-006). **CRITICAL: WAVE-07 source wave claims CAP-W07-003 "Week flags CRUD" — this directly collides with WAVE-04's existing scope.** WAVE-04 already specifies `createWeekFlag`, `deleteWeekFlag`, `weekFlags(by weekStartDate)` GraphQL mutations/queries, week_flag table, FlagType enum with 10 types, unique constraint per (week_start_date, flag_type). WAVE-07 must not recreate week flags CRUD; instead it must **import WAVE-04's WeekFlagService** for read operations in the export prompt builder.

### 1.5 WAVE-05 (Nutrition) — READY
WAVE-05 provides NutritionProduct, NutritionTemplate, NutritionTemplateItem, DailyNutritionOverride, DailyNutritionOverrideItem, and the NutritionMacroService for KJBJU calculation. WAVE-07 consumes via service layer (read-only). **Compatible — no collision.**

### 1.6 WAVE-06 (Charts) — READY
WAVE-06 is read-only chart queries. WAVE-07 shares the same underlying data but does not depend on WAVE-06's query services. **Compatible — no collision.**

---

## 2. Future Wave Compatibility

### 2.1 WAVE-08 (AI Review History) — depends on WAVE-07
WAVE-08 depends on WAVE-07. WAVE-08's AiReview entity links to a date range + AI response text. The AiExport table (WAVE-07) provides the export record and generated prompt. **Boundary is clean:** WAVE-07 creates the AiExport record (draft→generated lifecycle). WAVE-08 creates an independent AiReview record. WAVE-08 does not modify AiExport. The AiExport table serves as a source-of-truth for what was sent to the user; AiReview is a separate feedback/action store.

**Recommendation:** WAVE-07's AiExport table should store `generatedPrompt` (the prompt text sent to AI) so WAVE-08 can reference the export context. This is already in the domain model (AiExport.generatedPrompt). No additional contract needed.

### 2.2 WAVE-09 (Backup Import/Export) — similar ZIP pattern
WAVE-09 generates a full backup ZIP with `manifest.json` (data version), `data.json` (all entities), and `media/` folder. WAVE-07 generates an AI export ZIP with `manifest.json` (export metadata), `data.json` (period-scoped entities), `summary.md`, and CSV files. **Similar pattern, different purpose:** WAVE-07's manifest contains export-specific metadata (date range, selection toggles, user comment); WAVE-09's manifest contains data version for migration compatibility. WAVE-07's ZIP can be treated as a WAVE-09 export source (the ZIP file lives on disk, linkable from AiExport.exportFilePath). WAVE-09 can optionally back up AiExport records. **No scope collision — WAVE-07 creates per-period AI exports; WAVE-09 creates full data backups.**

**Recommendation:** Design WAVE-07's ZIP structure as a self-contained subfolder pattern (e.g., `ai-exports/<uuid>/manifest.json`, `data.json`, `summary.md`) to avoid filename collision with WAVE-09's root-level ZIP structure.

---

## 3. WAVE-03 Dependency — Workout Data Availability

### 3.1 Current state
WAVE-03 (Workout Diary) has **no detailed wave doc** and is **not implemented**. No `daily_log`, `workout_exercise`, `workout_set` tables exist. The domain model and source wave are approved but no detailed AC, EC, implementation slices, or code exist.

### 3.2 Impact on WAVE-07
WAVE-07's export includes workout data: exercises performed, sets with weight/reps/RPE/RIR, working weight snapshots. Without WAVE-03, this data is empty.

### 3.3 Stub pattern (same as WAVE-06)
WAVE-06 faced the same problem: exercise chart queries are stubbed returning empty series when WAVE-03 not deployed (DDEC-W06-010, DQ-W06-006). **Apply same pattern to WAVE-07:** the export `data.json` must include a `workout` section that is:
- An empty array `[]` when no WAVE-03 data exists for the period
- A populated array with relevant workout records when WAVE-03 data exists

The export query service must:
1. Attempt to query WAVE-03 tables via repository/service layer
2. If tables don't exist or return zero rows, return empty array
3. Include all workout data when available (no conditional compilation — runtime detection via SQL/table existence check or empty result handling)

### 3.4 Recommendation for detailed planning
Document in SLICEs that workout export queries are **conditional**: the service calls WAVE-03 repositories, which return empty results when the tables have no data. The summary.md must indicate "No workout data for this period" when empty. No need for table-existence checks — standard query returning zero rows achieves the same behavior.

---

## 4. Frontend Dependency Check — PAGE-009

| PAGE-009 Dependency | WAVE-07 Delivers? | Source |
|---|---|---|
| POST /api/ai-export (generate) | ✅ Must be created — generates prompt + ZIP | Core WAVE-07 scope |
| GET /api/ai-export/download | ✅ Must be created — downloads generated ZIP | Core WAVE-07 scope |
| GET /api/user-profile (goal context) | ⚠️ Gap — WAVE-01 Settings covers AI context but no UserProfile table/endpoint. Recommend reusing Settings GraphQL query or adding a lightweight REST endpoint. | See §1.2 |
| GET /api/week-flags | ❌ Scope collision — WAVE-04 owns WeekFlag CRUD via GraphQL. PAGE-009 specifies REST; WAVE-04 exposes GraphQL. WAVE-07 must either (a) add a REST proxy for week flags, or (b) have PAGE-009 use GraphQL. Recommend (b): PAGE-009 uses existing GraphQL `weekFlags(weekStartDate:)` query from WAVE-04. No new REST endpoint needed. | See §1.4 |

**Additional frontend requirement not in PAGE-009 but in WAVE-07 source wave:** Prompt display/copy. This is generated server-side (stored in AiExport.generatedPrompt, returned in POST /api/ai-export response). WAVE-07 must return the prompt text in the generate response body so the frontend can display it without downloading the ZIP.

---

## 5. Scope Collision Analysis

### 5.1 COLLISION: CAP-W07-003 (Week flags CRUD) vs CAP-W04-006
**Direct collision.** WAVE-04 detailed wave already specifies complete WeekFlag CRUD: table, sqlc queries, repo, service, GraphQL schema, resolvers, 3 ACs (AC-W04-039 to AC-W04-042), 3 tests (TEST-W04-020 to TEST-W04-022). WAVE-07 source wave claims "Week flags CRUD" as included scope.

**Resolution:** Strike CAP-W07-003 from WAVE-07. WAVE-07 imports WAVE-04's WeekFlagService for read queries. No write operations needed for AI export. The prompt builder needs to read week flags by week start date for display and inclusion in the prompt context — this is a read-only dependency.

### 5.2 PARTIAL OVERLAP: CAP-W07-001 (Persistent AI context) + CAP-W07-002 (User goal storage) vs WAVE-01 Settings
WAVE-01 Settings table stores AI context fields (`ai_goal`, `ai_height`, `ai_age`, `ai_experience`, `ai_split`, `ai_limits`, `ai_progression`, `ai_nutrition_strategy`). WAVE-07's "persistent AI context" and "user goal storage" are **served by reading Settings**. If WAVE-01 is implemented as specified, WAVE-07 does not need to create or store any user context data — it reads from Settings.

**Resolution:** WAVE-07 reads Settings for AI context and goal. If WAVE-01 Settings is not yet deployed, WAVE-07 needs a migration for a minimal settings/UserProfile table. Include a fallback: if Settings query returns empty for a field, use a sensible default or omit it from the prompt.

### 5.3 BOUNDARY OK: WAVE-07 vs WAVE-08 (AiExport vs AiReview)
- AiExport: id, dateRangeStart, dateRangeEnd, includePhotos, includeNutrition, includeCardio, includeMeasurements, userComment, generatedPrompt, exportFilePath
- AiReview: id, dateRangeStart, dateRangeEnd, aiResponseText, userNotes, plannedActions

Clean separation. WAVE-07 creates export records; WAVE-08 creates review records linked to a date range (not directly to an export). **No collision.**

### 5.4 BOUNDARY OK: WAVE-07 vs WAVE-09 (AI Export ZIP vs Backup ZIP)
Different ZIP structures, different purposes. WAVE-07: per-period AI-focused export. WAVE-09: full data backup with versioning. The shared pattern (manifest.json + data.json) is an implementation similarity, not a scope collision. WAVE-09 can reference WAVE-07's pattern when designing its own ZIP structure. **No collision; positive learning opportunity.**

---

## 6. Independent Deliverability

### 6.1 Can WAVE-07 be implemented without WAVE-03?
**YES.** Apply the same stub pattern as WAVE-06: workout data section in export returns empty when no WorkoutExercise/WorkoutSet data exists. The export query service calls WAVE-03 repositories which return zero rows. The summary.md includes "No workout data recorded for this period."

### 6.2 Can WAVE-07 be implemented without WAVE-04?
**YES, with limitations.** Without WAVE-04:
- Cardio data: empty in export
- Body weight/measurements: empty in export
- Week flags: empty (no flag context in prompt)
- Progress photos: unavailable (user must deselect photos in section toggles)

The export still generates data.json with available sections and marks empty sections as `[]`.

### 6.3 Can WAVE-07 be implemented without WAVE-05?
**YES.** Nutrition data section returns empty array. Summary.md indicates nutrition data unavailable.

### 6.4 Can WAVE-07 be implemented without WAVE-01?
**NO — hard dependency.** WAVE-01 provides:
- PIN auth middleware (required for all API endpoints)
- Settings with AI context and goal fields (CAP-W07-001, CAP-W07-002)
- Media storage scaffold (for storing generated ZIP files on disk)
- GraphQL infrastructure

Without WAVE-01, WAVE-07 has no auth layer, no AI context, no user identity.

### 6.5 Core WAVE-07 functionality that does NOT depend on upstream data
- AiExport record creation and lifecycle (draft → generated)
- Prompt generation using AI context from Settings + period selection
- ZIP creation (manifest.json, data.json with available data, summary.md)
- ZIP file management (storage, download, cleanup)
- Section toggle processing (which data types to include)
- User comment handling

All of this works even if all upstream waves return empty data.

---

## 7. Recommended ACs and ECs (Draft)

### Acceptance Criteria

| AC ID | Description | Source |
|---|---|---|
| AC-W07-001 | AiExport record created via POST /api/ai-export with dateRangeStart, dateRangeEnd, section toggles, optional userComment. Returns export ID, generated prompt text, download URL. | PAGE-009, AiExport domain model |
| AC-W07-002 | AiExport created without photos section returns includePhotos=false. Photos excluded by default per domain invariant #10. | Domain model invariant #10 |
| AC-W07-003 | AiExport generated prompt includes user goal and AI context from Settings. | CAP-W07-001, CAP-W07-002 |
| AC-W07-004 | AiExport generated prompt includes week flags for the selected period (reads from WAVE-04 WeekFlagService). | CAP-W07-003 (read-only reframe) |
| AC-W07-005 | AiExport ZIP downloaded via GET /api/ai-export/download returns ZIP file with manifest.json, data.json, summary.md. | OUT-W07-002 |
| AC-W07-006 | manifest.json contains export version, date range, user profile summary, section toggle states, generation timestamp. | CAP-W07-006 |
| AC-W07-007 | data.json contains structured data for all selected sections within the date range. Empty sections return empty arrays. | CAP-W07-007, DDEC-W07-001 |
| AC-W07-008 | summary.md contains human-readable overview with period summary, section-by-section highlights. | CAP-W07-008 |
| AC-W07-009 | Workout section in data.json is empty when no WAVE-03 data exists for the period. | WAVE-03 not implemented |
| AC-W07-010 | Section toggles are honored: deselected sections (e.g., photos, nutrition, cardio) excluded from ZIP. | OUT-W07-005 |
| AC-W07-011 | User comment included in manifest.json and summary.md when provided. | OUT-W07-004 |
| AC-W07-012 | AiExport record persisted with draft→generated lifecycle. exportFilePath populated after ZIP generation. | Domain model lifecycle |
| AC-W07-013 | POST /api/ai-export returns ValidationError for invalid date range (from > to) or missing required fields. | Standard error handling |
| AC-W07-014 | All WAVE-07 endpoints return 401/403 without valid PIN session. | WAVE-01 PIN auth guard |
| AC-W07-015 | ZIP file stored under configured media path, served only through authenticated GET endpoint. | Consistent with WAVE-01 media pattern |
| AC-W07-016 | Duplicate AiExport request for same date range returns new export (no dedup constraint — each export is independent). | Equivalent to WAVE-02 duplicate name rule |

### Exit Criteria

| EC ID | Description |
|---|---|
| EC-W07-001 | AC-W07-001 through AC-W07-016 pass via focused tests |
| EC-W07-002 | ZIP structure validated: manifest.json, data.json, summary.md present and well-formed |
| EC-W07-003 | Prompt generation produces valid AI-ready text with goal context, period, flags, and structured data references |
| EC-W07-004 | Empty data sections (no WAVE-03, empty period) produce valid empty arrays, not errors |
| EC-W07-005 | All endpoints protected by WAVE-01 PIN auth middleware |
| EC-W07-006 | ZIP stored on disk, deletable via export record deletion |
| EC-W07-007 | Lint passes for all changed packages |
| EC-W07-008 | Codegen produces valid Go code (gqlgen, sqlc) without drift |
| EC-W07-009 | No WAVE-04 week flag write operations in WAVE-07 (read-only via WeekFlagService) |
| EC-W07-010 | No WAVE-03 table creation in WAVE-07 (workout data queried when available) |

---

## 8. Risks

| Risk | Severity | Mitigation |
|---|---|---|
| WAVE-03 not implemented — workout data empty | Medium | Stub pattern (empty arrays). Document in summary.md. Add AC-W07-009. |
| WAVE-04 not implemented — no week flags, cardio, body data | Medium | Empty section arrays. Section toggles still work. |
| WAVE-01 Settings not yet deployed — no AI context | High | Cannot deliver without WAVE-01. Add fallback: inline defaults for prompt generation, log warning. |
| CAP-W07-003 collision with WAVE-04 | High | Strike from WAVE-07. Import WAVE-04 WeekFlagService for reads. Add EC-W07-009 to enforce. |
| PAGE-009 specifies REST endpoints; WAVE-04 exposes GraphQL | Medium | WAVE-07 can add thin REST handlers that delegate to existing GraphQL queries, or PAGE-009 uses GraphQL directly. |
| Large date ranges generate large ZIPs | Low | Date range defaults from WAVE-01 Settings.defaultAiExportWeeks (default 4). Max 52 weeks per WAVE-06 cap precedent. |
| Export file cleanup/orphaned ZIPs | Low | AiExport deletion should clean up ZIP file from disk. Consistent with WAVE-04 progress photo deletion pattern. |

---

## 9. Open Questions

| ID | Scope | Question | Why It Matters | Source |
|---|---|---|---|---|
| DQ-W07-001 | product-scope-and-ac | Should WAVE-07 add a UserProfile table or reuse WAVE-01 Settings for goal/AI context? | PAGE-009 needs GET /api/user-profile. If Settings already stores AI context, no new table needed. | PAGE-009, WAVE-01 Settings AC-W01-001 |
| DQ-W07-002 | data-ops | PAGE-009 specifies GET /api/week-flags (REST). WAVE-04 exposes weekFlags (GraphQL). Should WAVE-07 add REST proxy or PAGE-009 use GraphQL? | WAVE-07 must not recreate week flags CRUD. API surface consistency. | PAGE-009, WAVE-04 week flags |
| DQ-W07-003 | product-ac | What prompt format should the AI export use? (Markdown, structured JSON, system message format?) | AI readiness of generated prompt. Affects CAP-W07-004 implementation. | CAP-W07-004 |
| DQ-W07-004 | data-ops | Should WAVE-07 create AiExport table migration or reuse existing migration sequence? | Migration numbering. WAVE-04 uses 00082-00087/00091, WAVE-05 needs next number. | Migration coordination |
| DQ-W07-005 | architecture | REST vs GraphQL for export endpoints? | Frontend PAGE-009 specifies REST. WAVE-07 could use REST for file operations, GraphQL for create. | WAVE-01 hybrid pattern precedent (GraphQL CRUD + REST binary) |

---

## 10. Traceability

- `docs/prd-waves/waves/wave-07.md`: source wave boundary (CAP-W07-001 through CAP-W07-009)
- `docs/product-verified/domain-model.md`: AiExport entity, AiReview entity, Settings attributes, UserProfile attributes, WeekFlag attributes
- `docs/prd-waves/frontend-pages/page-009.md`: POST /api/ai-export, GET /api/ai-export/download, GET /api/user-profile, GET /api/week-flags
- `docs/prd-wave-details/waves/wave-01.md`: Settings with AI context fields, PIN auth, media scaffold
- `docs/prd-wave-details/waves/wave-02.md`: exercise metadata and media for export consumption
- `docs/prd-wave-details/waves/wave-04.md`: WeekFlag CRUD (CAP-W04-006), cardio, body data
- `docs/prd-wave-details/waves/wave-05.md`: NutritionProduct, NutritionTemplate, NutritionMacroService
- `docs/prd-wave-details/waves/wave-06.md`: WAVE-03 stub pattern (DDEC-W06-010), empty series precedent
- `docs/prd-waves/waves/wave-03.md`: source wave (no detailed doc)
- `docs/prd-waves/waves/wave-08.md`: AiReview depends on WAVE-07
- `docs/prd-waves/waves/wave-09.md`: Backup ZIP pattern (learning opportunity)