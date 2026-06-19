# WAVE-02: Exercise Library Design

## Status
ready-for-dev — approved by user 2026-06-19, all 7 reviewer perspectives approved (2 cycles).

## Corrections Applied (2026-06-19)
1. Added user_id to exercises and exercise_media tables
2. Updated indexes to include user_id
3. Use WAVE-01 media scaffold routes (no new /api/v1/exercise-media namespace)
4. Exercise CRUD is GraphQL only — no REST
5. Replaced deleteExercise with archiveExercise (soft archive)
6. Clarified media delete: archive exercise does NOT cascade delete media
7. exercise(id: ID!) returns ExerciseResult union (not Exercise!)
8. Duplicate names allowed (unchanged — already correct)
9. AI/future references use exerciseId as primary identity, name as display

## Final Consistency Fixes (2026-06-19)
A. Codebase paths aligned with WAVE-01 Approach A (Atlas-specific directories)
B. exercise_media.exercise_id UUID NOT NULL
C. working_weight CHECK (IS NULL OR > 0) at DB level

## Data Model

### exercises
| Column | Type | Constraints |
|--------|------|-------------|
| id | UUID | PK, default gen_random_uuid() |
| user_id | UUID | NOT NULL, FK → atlas_users(id) |
| name | TEXT | NOT NULL |
| muscle_groups | TEXT[] | |
| description | TEXT | |
| personal_notes | TEXT | |
| working_weight | REAL | nullable, CHECK (working_weight IS NULL OR working_weight > 0) |
| is_active | BOOLEAN | DEFAULT true |
| created_at | TIMESTAMPTZ | DEFAULT now() |
| updated_at | TIMESTAMPTZ | DEFAULT now() |

Indexes:
- `idx_exercises_user_active (user_id, is_active)`
- `idx_exercises_user_name (user_id, name)`
- `idx_exercises_user_created_at (user_id, created_at)` (optional)

No unique constraint on name — duplicates allowed.

### exercise_media
| Column | Type | Constraints |
|--------|------|-------------|
| id | UUID | PK, default gen_random_uuid() |
| user_id | UUID | NOT NULL, FK → atlas_users(id) |
| exercise_id | UUID | NOT NULL, FK → exercises(id), ON DELETE NO ACTION |
| file_name | TEXT | NOT NULL |
| file_path | TEXT | NOT NULL |
| mime_type | TEXT | NOT NULL |
| file_size | BIGINT | NOT NULL |
| created_at | TIMESTAMPTZ | DEFAULT now() |

Index: `idx_exercise_media_user_exercise (user_id, exercise_id)`.

FK `ON DELETE NO ACTION` — archiving an exercise does NOT cascade delete media.

## API Surface

### GraphQL (on /graphql/atlas, PIN-protected)

CRUD is GraphQL only. No REST for exercise operations.

**Types:**
- `Exercise` — id, userId, name, muscleGroups, description, personalNotes, workingWeight, isActive, media ([ExerciseMedia]), createdAt, updatedAt
- `ExerciseMedia` — id, userId, exerciseId, fileName, mimeType, fileSize, createdAt
- `ExerciseConnection` — items ([Exercise]), totalCount, pageInfo (hasNextPage, endCursor)
- `CreateExerciseInput` — name (required), muscleGroups, description, personalNotes, workingWeight
- `UpdateExerciseInput` — name, muscleGroups, description, personalNotes, workingWeight
- `ExerciseResult` — union: Exercise | ValidationError | NotFoundError | AuthError
- `ArchiveResult` — union: Exercise | NotFoundError | AuthError

**Queries:**
- `exercises(first: Int = 20, after: String, includeInactive: Boolean = false): ExerciseConnection!`
- `exercise(id: ID!): ExerciseResult!`
- `allExercises(includeInactive: Boolean = false): [Exercise!]!` (unpaginated, for WAVE-03)

**Mutations:**
- `createExercise(input: CreateExerciseInput!): ExerciseResult!`
- `updateExercise(id: ID!, input: UpdateExerciseInput!): ExerciseResult!`
- `archiveExercise(id: ID!): ArchiveResult!`
- `restoreExercise(id: ID!): ArchiveResult!` (optional)

### Media (REST via WAVE-01 scaffold, PIN-protected)

Exercise media uses the WAVE-01 media scaffold routes with a purpose/entity metadata pattern:

- `POST /api/v1/media/upload` — multipart with purpose=EXERCISE_MEDIA + exerciseId
- `GET /api/v1/media/{id}` — file download, correct Content-Type
- `DELETE /api/v1/media/{id}` — delete DB record + file from disk, 204

File validation: server-side http.DetectContentType() — JPEG/PNG/WEBP/MP4/MOV/WEBM only. Image ≤ 25MB, video ≤ 250MB, max single upload 300MB.

Media delete removes DB record + physical file — same as WAVE-02's original media handler logic, but routed through the WAVE-01 scaffold paths.

Archiving an exercise does NOT cascade delete its media. Media records remain in DB with is_active exercise reference.

## Architecture (WAVE-01 Approach A paths)

```
GraphQL resolvers:    apps/api/internal/atlas/graph/resolver/exercise.resolvers.go
    → ExerciseService:  apps/api/internal/atlas/service/exercise.go
        → ExerciseRepo:   apps/api/internal/atlas/repository/postgres/exercise_repo.go
            → sqlc:          apps/api/internal/repository/postgres/queries/exercises.sql
                → PostgreSQL via goose migrations 00081_exercises, 00082_exercise_media

Media via WAVE-01 scaffold: apps/api/internal/handler/atlas_media.go (extended)
    atlas-gqlgen config: apps/api/atlas-gqlgen.yml (auto-discovers internal/atlas/graph/schema/*.graphql)
```

## Implementation Slices

1. **SLICE-W02-001** — DB migrations (00081_exercises.sql, 00082_exercise_media.sql)
2. **SLICE-W02-002** — sqlc queries (exercises.sql in internal/repository/postgres/queries)
3. **SLICE-W02-003** — Exercise repository (internal/atlas/repository/postgres/exercise_repo.go)
4. **SLICE-W02-004** — Exercise service (internal/atlas/service/exercise.go)
5. **SLICE-W02-005** — GraphQL schema (internal/atlas/graph/schema/exercises.graphql)
6. **SLICE-W02-006** — GraphQL resolvers (internal/atlas/graph/resolver/exercise.resolvers.go)
7. **SLICE-W02-007** — Extend WAVE-01 media scaffold (handler/atlas_media.go) with purpose/entity routing
8. **SLICE-W02-008** — Main wiring (main.go)

## Testing

22 verification obligations: 6 unit, 13 integration, 2 codegen, 1 lint.
Full round-trip integration test, PIN auth guards, file validation, log privacy.
Codegen uses `bunx nx run api:codegen:atlas` for Atlas gqlgen.

## Future Reference Note

Future Workout/AI export references must use exerciseId as primary identity and name as display field. Identity is always exercise.id, never exercise.name.

## Source Documents
- docs/prd-wave-details/waves/wave-02.md — full detailed wave brief
- docs/prd-waves/waves/wave-02.md — shallow wave definition
- docs/prd-waves/wave-map.md — wave sequence context
- docs/prd-waves/frontend-pages/page-003.md — frontend dependency context