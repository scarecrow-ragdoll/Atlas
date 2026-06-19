# WAVE-02 sequencing-fit Planner Attempt 1

## Sources Read
- docs/prd-waves/wave-map.md
- docs/prd-waves/waves/wave-02.md
- docs/prd-wave-details/waves/wave-01.md (prior detailed wave)
- docs/prd-waves/frontend-pages/page-002.md (Workout Diary)
- docs/prd-waves/frontend-pages/page-003.md (Exercise Library)
- docs/technical-verified/implementation-slices.md (Slice 1)
- docs/technical-verified/architecture-and-boundaries.md

## Selected Backend Wave Boundary
WAVE-02 implements Exercise CRUD and ExerciseMedia management exclusively. It depends on WAVE-01 (infrastructure) and provides dependencies to WAVE-03. No overlap with WAVE-04 through WAVE-09.

## Prior Backend Wave Fit (WAVE-01)
- WAVE-01 provides: migration infrastructure, PIN auth middleware, media REST scaffold (POST/GET /api/v1/media), fitness GraphQL schema extension pattern, sqlc/gqlgen codegen config for fitness domain, go API fitness packages structure, config extension pattern.
- WAVE-02 consumes all of these directly.
- **Compatibility check**: WAVE-01's media REST scaffold must support ExerciseMedia association. WAVE-01 creates POST /api/v1/media/upload and GET /api/v1/media/{id}. WAVE-02 needs POST /api/v1/exercise-media and DELETE /api/v1/exercise-media/{id}. These are additive — WAVE-01 scaffold is extended, not modified.
- **WAVE-01 dependency exposure**: if WAVE-01's media handler stores files without an exercise association mechanism, WAVE-02 must either (a) extend the existing handler or (b) create a separate handler that links to exercise. Recommendation: separate handler to maintain clean boundaries.
- **No scope collision**: WAVE-01 does not touch exercise domain.

## Future Backend Wave Fit (WAVE-03)
- WAVE-03 (Workout Diary) depends on GET /api/exercises for exercise selection in workout creation flow.
- WAVE-02 must provide a simple exercise list query (`allExercises`) that returns all active exercises without pagination for WAVE-03's exercise selector dropdown.
- WAVE-03 reads Exercise.workingWeight at log time and snapshots it into WorkoutExercise.workingWeightSnapshot (RULE-017). This is WAVE-03 responsibility, not WAVE-02.
- WAVE-03 also depends on exercise rows persisting after soft delete (for historical workout data). WAVE-02's soft-delete (isActive=false) preserves rows.
- **Collision check**: WAVE-02 does not duplicate any WAVE-03 functionality. Exercise CRUD (WAVE-02) and Workout logging (WAVE-03) are cleanly separated.

## Future Backend Wave Fit (WAVE-04+)
- WAVE-04 (Cardio/Body): no dependency on exercises.
- WAVE-05 (Nutrition): no dependency on exercises.
- WAVE-06 (Charts): may depend on exercise data for training charts (working weight over time, volume, e1RM). WAVE-02 provides the exercise entity; WAVE-06 reads it. No scope collision.
- WAVE-07/08 (AI Export/Review): may include exercise data in exports. WAVE-02 provides the data source.
- WAVE-09 (Backup): includes exercise and exercise_media tables in full backup. WAVE-02 data model must be compatible with WAVE-09 backup format (data.json with exercise entities).

## Frontend Pages Context
- PAGE-003 (Exercise Library): primary consumer. Backend dependencies: exercise CRUD (GraphQL) + media management (REST).
- PAGE-002 (Workout Diary): reads exercise list for exercise selector in workout creation.
- PAGE-006 (Charts): future consumer of exercise data for training charts.

## Dependency Order
```
WAVE-01 (Foundation) → WAVE-02 (Exercise Library) → WAVE-03 (Workout Diary)
                                                → WAVE-06 (Charts — reads)
                                                → WAVE-07/08 (AI — reads)
                                                → WAVE-09 (Backup — reads)
```

No circular dependencies. WAVE-02 is a necessary prerequisite for WAVE-03, WAVE-06, WAVE-07/08, and WAVE-09 exercise data.

## Scope Collision Check
| Neighboring Wave | Scope Collision | Notes |
| --- | --- | --- |
| WAVE-01 (Foundation) | None | WAVE-01 is infrastructure only; WAVE-02 is first domain wave |
| WAVE-03 (Workout Diary) | None | WAVE-02 owns exercise CRUD; WAVE-03 owns workout logging. Exercise selector is a read-only dependency. |
| WAVE-04 (Cardio/Body) | None | Independent domain |
| WAVE-05 (Nutrition) | None | Independent domain |
| WAVE-06 (Charts) | None | WAVE-02 provides data; WAVE-06 reads it |
| WAVE-07/08 (AI) | None | WAVE-02 provides data; AI waves read it |
| WAVE-09 (Backup) | None | WAVE-02 provides data tables; backup wave reads them |

## Independent Value
WAVE-02 delivers standalone value after WAVE-01: users can manage their exercise library even before WAVE-03 (Workout Diary) is built. Exercise CRUD is independently useful and testable.

## Design Contracts For Sequencing

### WAVE-02 → WAVE-03 interface
- WAVE-02 provides GraphQL query `allExercises(includeInactive: Boolean = false): [Exercise!]!` returning all exercises (no pagination) for WAVE-03 exercise selector.
- WAVE-02 exercises remain queryable by ID even after soft-delete (for historical workout references in WAVE-03).
- WAVE-02 exercise.workingWeight is the current value. WAVE-03 snapshots it at log time into WorkoutExercise.workingWeightSnapshot.

### WAVE-02 → WAVE-06 interface
- WAVE-02 provides exercise data (working weight, isActive) for training chart queries.
- Chart aggregation queries are WAVE-06 responsibility.

### WAVE-02 → WAVE-09 interface
- Exercise and exercise_media tables are included in full backup data.json.
- WAVE-09 reads directly from tables or through service layer.

## Acceptance Criteria Contributions

| AC ID | Description |
| --- | --- |
| AC-W02-021 | WAVE-02 provides `allExercises` query returning all active exercises (no pagination) for WAVE-03 selector |
| AC-W02-022 | Exercise can be queried by ID after soft-delete (isActive=false) for WAVE-03 historical reference |

## Exit Criteria Contributions

| EC ID | Description |
| --- | --- |
| EC-W02-020 | `allExercises` query returns all active exercises without pagination |
| EC-W02-021 | Exercise by ID query returns exercise even when isActive=false |
| EC-W02-022 | No WAVE-03 or later wave functionality implemented in WAVE-02 |

## Risks And Rollback
- Risk: if WAVE-01 media REST scaffold is not yet delivered, WAVE-02 media upload is blocked. Mitigation: WAVE-01 is `ready-for-dev`; implement WAVE-01 first.
- Risk: if WAVE-03 needs additional exercise fields (e.g., defaultRestTime, exerciseCategory) that are not in WAVE-02, a schema change is needed. Mitigation: all fields currently known are in domain model; future fields can be added as optional migration.
- Rollback: WAVE-02 code and migrations are fully reversible. No other wave is blocked by WAVE-02 rollback since WAVE-03 (which depends on this) would not exist yet.

## Questions Raised

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Source Or Report | Status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DQ-W02-008 | WAVE-02 | sequencing-fit | needs-owner-decision | WAVE-03 | Does `allExercises` need to support additional filtering beyond isActive (e.g., by muscle group) for WAVE-03? | Affects WAVE-02 query API surface; may need filtering if exercise library grows large (future concern) | sequencing-fit planner | open |

## Traceability Candidates
- Dependency order → docs/prd-waves/wave-map.md
- WAVE-03 dependency → docs/prd-waves/frontend-pages/page-002.md (GET /api/exercises)
- Independent value → docs/technical-verified/implementation-slices.md (Slice 1)
- No scope collision → docs/prd-waves/waves/wave-02.md (excluded scope: workout diary)