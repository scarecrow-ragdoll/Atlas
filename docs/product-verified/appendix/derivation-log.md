# Derivation Log

## Derived Roles And Permissions

| ID | Derived Item | Source Signal | Rationale | Confidence |
| --- | --- | --- | --- | --- |
 | D-ROLE-001 | User has all CRUD permissions on exercises | PRD §11, §26.1 | User creates, edits, and deletes exercises and media | High |
| D-ROLE-002 | User has all permissions on workout data | PRD §10, §26.2-§26.3 | User creates, edits, and saves workout days, exercises, sets | High |
| D-ROLE-003 | User has all permissions on cardio data | PRD §12, §26.4 | User adds, edits cardio entries | High |
| D-ROLE-004 | User has all permissions on body tracking data | PRD §13, §26.5-§26.6 | User creates check-ins and weight entries | High |
| D-ROLE-005 | User has all permissions on nutrition data | PRD §15, §26.7-§26.8 | User creates products, templates, overrides | High |
| D-ROLE-006 | User can generate AI exports | PRD §17, §26.9 | Export is on-demand user action | High |
| D-ROLE-007 | User can save AI reviews | PRD §19, §26.10 | Review is manual user entry | High |
| D-ROLE-008 | User can export and import backups | PRD §20, §26.11-§26.12 | Backup/restore is user-initiated | High |
| D-ROLE-009 | User can configure settings (PIN, units, profile) | PRD §7.2, §25.1-§25.2 | Settings section accessible to user | High |

## Derived Data Fields

| ID | Derived Field | Source Signal | Rationale | Confidence |
| --- | --- | --- | --- | --- |
| D-FIELD-001 | CardioType enum values | PRD §12.3 | Listed types: walking, running, bike, elliptical, treadmill, other | High |
| D-FIELD-002 | HeartRateZone enum values | PRD §12.4 | Zone 1-5 plus unknown | High |
| D-FIELD-003 | MeasurementType enum values | PRD §13.3 | 10 listed measurement types | High |
| D-FIELD-004 | Side enum for paired measurements | PRD §13.4 | Left, right, or null (common) | High |
| D-FIELD-005 | FlagType enum values | PRD §18.4 | 10 listed flag types | High |
| D-FIELD-006 | MediaType enum values | PRD §11.3 | Image and video | Medium |
| D-FIELD-007 | DailyNutritionOverrideItem operation values | PRD §25.17 | Add, subtract, replace | High |
| D-FIELD-008 | Photo angle enum values | PRD §14.2 | Front, side, back, custom | High |

## Derived States

| ID | Derived State | Source Signal | Rationale | Confidence |
| --- | --- | --- | --- | --- |
| D-STATE-001 | Exercise: active / inactive | PRD §11.2 | isActive boolean field | High |
| D-STATE-002 | AiExport: draft / generated | PRD §25.19 | exportFilePath optional — presence indicates generation | Medium |
| D-STATE-003 | NutritionTemplate: active / superseded | PRD §15.3 | Single template at a time | Medium |
| D-STATE-004 | WorkoutDay: empty / has-exercises | PRD §10.2 | Created on first save | Medium |

## Derived Acceptance Criteria

See docs/product-verified/acceptance-criteria.md for 125 derived AC items. All ACs trace to PRD source sections as noted in traceability.md.

## Derived Edge Cases

See docs/product-verified/edge-cases.md for 31 derived edge cases. All edge cases are derived from documented operations plus standard boundary and failure classes.

## Low-Confidence Derivations

| ID | Item | Reason | Linked Question |
| --- | --- | --- | --- |
| D-LOW-001 | Deployer role with setup flow | No source evidence for setup behavior | Q-ROLE-005 |
| D-LOW-002 | e1RM formula | Formula unspecified in PRD | Q-AC-07 |
| D-LOW-003 | Best set definition (heaviest vs highest e1RM) | PRD mentions best approach but not definition | Q-AC-08 |