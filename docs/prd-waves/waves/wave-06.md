# Wave 06: Charts

## Status

user-approved

## User Approval

user-approved (2026-06-18)

## Purpose

Progress visualization for workouts, body measurements, and nutrition.

## Outcome After Wave

- OUT-W06-001 Exercise charts (weight, 1RM, volume, reps)
- OUT-W06-002 Body charts (weight, measurements)
- OUT-W06-003 Nutrition charts (weekly macro averages)
- OUT-W06-004 Period filtering
- OUT-W06-005 Ready for AI export data aggregation

## Included Scope

- CAP-W06-001 Exercise progress queries
- CAP-W06-002 1RM calculation (needs formula decision)
- CAP-W06-003 Body weight trend queries
- CAP-W06-004 Body measurement trend queries
- CAP-W06-005 Nutrition macro summary queries
- CAP-W06-006 Period filtering
- CAP-W06-007 Chart data queries

## Excluded Scope

- Photo charts
- Advanced analytics

## Dependencies

WAVE-01, WAVE-03, WAVE-04, WAVE-05

## Surface Categories

backend, data, operations

## Risk Class

Medium - 1RM formula needs technical decision

## Recommended Next Planning

$detail-prd-wave for WAVE-06

## Open Questions

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Q-CHART-001 | 06 | operations | Medium | None | Exact 1RM formula? | Chart accuracy | docs/product/prd.md Section 16.2 | open | needs-owner-decision |

## Traceability

- docs/product/prd.md Section 16
- docs/product-verified/domain-model.md#WorkoutSet