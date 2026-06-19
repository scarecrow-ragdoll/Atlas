# Wave 02: Exercise Library

## Status

user-approved

## User Approval

user-approved (2026-06-19) — corrections applied: user_id FK on all entities, archive/restore instead of delete, media via WAVE-01 scaffold routes, GraphQL-only CRUD, exercise query returns union type, paths aligned with WAVE-01 Approach A, NOT NULL exercise_id, working_weight CHECK constraint.

## Purpose

Full CRUD for exercises with working weight and media management.

## Outcome After Wave

- OUT-W02-001 Exercises can be created, listed, edited, archived/restored
- OUT-W02-002 Media (images/video) can be attached via WAVE-01 media scaffold
- OUT-W02-003 Working weight stored per exercise
- OUT-W02-004 API ready for workout diary

## Included Scope

- CAP-W02-001 Exercise CRUD (GraphQL mutations/queries)
- CAP-W02-002 ExerciseMedia upload and retrieval via WAVE-01 media scaffold
- CAP-W02-003 Working weight field
- CAP-W02-004 Muscle groups, description, notes
- CAP-W02-005 isActive flag for soft archive

## Excluded Scope

- Workout diary integration (uses API)
- AI export specifics
- New REST media namespace (uses WAVE-01 scaffold)
- Hard delete of exercises (archive only)

## Dependencies

WAVE-01 (Foundation) — PIN auth middleware, media scaffold routes, GraphQL infra, migration infra

## Surface Categories

backend, data, operations

## Risk Class

Low - Standard CRUD with file handling

## Recommended Next Planning

$detail-prd-wave for WAVE-02 (completed)

## Open Questions

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |

## Traceability

- docs/product/prd.md Section 11
- docs/product-verified/domain-model.md#Exercise, #ExerciseMedia
- docs/prd-wave-details/waves/wave-02.md: detailed wave brief with corrections and final consistency fixes applied
- Future reference: AI exports use exerciseId as primary identity
- WAVE-01 Approach A paths: atlas/ subdirectory for service, repo, graph/schema, graph/resolver; atlas-gqlgen.yml for codegen