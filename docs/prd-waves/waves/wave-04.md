# Wave 04: Cardio and Body Tracking

## Status

user-approved

## User Approval

user-approved (2026-06-18)

## Purpose

Cardio tracking and body measurements with weekly check-ins and progress photos.

## Outcome After Wave

- OUT-W04-001 Cardio entries with type, duration, pulse, zone
- OUT-W04-002 Body weight entries (standalone or in check-in)
- OUT-W04-003 Weekly body check-ins
- OUT-W04-004 Body measurements (paired left/right)
- OUT-W04-005 Progress photos with angles

## Included Scope

- CAP-W04-001 CardioEntry CRUD
- CAP-W04-002 BodyWeightEntry CRUD
- CAP-W04-003 BodyCheckIn CRUD
- CAP-W04-004 BodyMeasurement with side (left/right)
- CAP-W04-005 ProgressPhoto with angle (front/side/back/custom)
- CAP-W04-006 WeekFlag tracking

## Excluded Scope

- Photo taken in app (upload only)
- Chart visualization

## Dependencies

WAVE-01

## Surface Categories

backend, data, operations

## Risk Class

Low - Standard CRUD with photo handling

## Recommended Next Planning

$detail-prd-wave for WAVE-04

## Open Questions

| ID | Wave | Scope | Severity | Parent | Question | Why It Matters | Source Or Report | Status | Resolution |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |

## Traceability

- docs/product/prd.md Sections 12, 13
- docs/product-verified/domain-model.md#CardioEntry, #BodyWeightEntry, #BodyCheckIn, #BodyMeasurement, #ProgressPhoto