# Wave 02: Exercise Library

## Status

user-approved

## User Approval

user-approved (2026-06-18)

## Purpose

Full CRUD for exercises with working weight and media management.

## Outcome After Wave

- OUT-W02-001 Exercises can be created, listed, edited, deleted
- OUT-W02-002 Media (images/video) can be attached
- OUT-W02-003 Working weight stored per exercise
- OUT-W02-004 API ready for workout diary

## Included Scope

- CAP-W02-001 Exercise CRUD (GraphQL mutations/queries)
- CAP-W02-002 ExerciseMedia upload and retrieval
- CAP-W02-003 Working weight field
- CAP-W02-004 Muscle groups, description, notes
- CAP-W02-005 isActive flag for soft delete

## Excluded Scope

- Workout diary integration (uses API)
- AI export specifics

## Dependencies

WAVE-01

## Surface Categories

backend, data, operations

## Risk Class

Low - Standard CRUD with file handling

## Recommended Next Planning

$detail-prd-wave for WAVE-02

## Open Questions

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |

## Traceability

- docs/product/prd.md Section 11
- docs/product-verified/domain-model.md#Exercise, #ExerciseMedia