# Aggregate Question Ledger — Run 20260618T185935Z

## Blocking Questions

| ID | Scope | Severity | Question | Why It Matters | Source | Status |
| --- | --- | --- | --- | --- | --- | --- |
| Q-SCOPE-001 | product-scope | blocking | What are the quantitative success metrics for the MVP? | Without success metrics, there is no way to validate the product is achieving its goals | prd.md (no success metrics found) | resolved |
| Q-SCOPE-002 | product-scope | blocking | Should the MVP architecture build for future multi-user or remain strictly single-user? | Affects data model, auth design, API structure, and database schema decisions | prd.md sections 4, 28 | resolved |
| Q-SCOPE-004 | product-scope | blocking | What are the specific performance targets (page load, export time, chart render time)? | Section 24.3 uses vague language with no measurable targets | prd.md section 24.3 | resolved |
| Q-SCOPE-005 | product-scope | blocking | Should cardio be a separate entity or always part of a workout day? | Data model shows separate with optional workoutDayId, but section 10.3 includes cardio in the workout day | prd.md sections 10.3, 25.8 | resolved |

## Non-Blocking Questions

### Product Scope

| ID | Scope | Severity | Question | Source | Status |
| --- | --- | --- | --- | --- | --- |
| Q-SCOPE-003 | product-scope | non-blocking | What is the target user's expected technical proficiency? | prd.md section 1 | open |
| Q-SCOPE-006 | product-scope | non-blocking | What AI models/platforms must the export format support? | prd.md sections 1, 17 | open |
| Q-SCOPE-007 | product-scope | non-blocking | Is there a maximum photo/media storage limit? | prd.md sections 14, 24.2 | open |
| Q-SCOPE-008 | product-scope | non-blocking | What data portability standard is required? | prd.md section 6.3 | open |

### Roles and Permissions

| ID | Scope | Severity | Question | Source | Status |
| --- | --- | --- | --- | --- | --- |
| Q-ROLE-001 | roles-permissions | non-blocking | PIN session lifetime and renewal policy | worker-attempt-1 | open |
| Q-ROLE-002 | roles-permissions | non-blocking | Logout mechanism when PIN is enabled | worker-attempt-1 | open |
| Q-ROLE-003 | roles-permissions | non-blocking | Access control when PIN is disabled | worker-attempt-1 | open |
| Q-ROLE-004 | roles-permissions | deferred | Resource-level visibility control | worker-attempt-1 | deferred |
| Q-ROLE-005 | roles-permissions | non-blocking | Deployer setup flow / configuration | worker-attempt-1 | open |

### Actor Journey

| ID | Scope | Severity | Question | Source | Status |
| --- | --- | --- | --- | --- | --- |
| Q-ACTOR-01 | actor-journey | non-blocking | Add exercise to workout when exercise doesn't exist | worker-attempt-2 | open |
| Q-ACTOR-02 | actor-journey | non-blocking | Add set with invalid values (0 reps, 0 weight) | worker-attempt-2 | open |
| Q-ACTOR-03 | actor-journey | non-blocking | Weekly check-in without photos — 2-4 required or recommended? | worker-attempt-2 | open |
| Q-ACTOR-04 | actor-journey | non-blocking | Removing last exercise from workout day — delete day? | worker-attempt-2 | open |
| Q-ACTOR-05 | actor-journey | non-blocking | Same exercise twice in one workout day | worker-attempt-2 | open |
| Q-ACTOR-06 | actor-journey | non-blocking | Reset nutrition overrides to template? | worker-attempt-2 | open |
| Q-ACTOR-07 | actor-journey | non-blocking | Empty data results UI (AI export, charts) | worker-attempt-2 | open |
| Q-ACTOR-08 | actor-journey | non-blocking | Backup import when data already exists | worker-attempt-2 | open |
| Q-ACTOR-09 | actor-journey | non-blocking | Wrong PIN handling, lockout? (cross-scope: roles/permissions) | worker-attempt-2 | open |
| Q-ACTOR-10 | actor-journey | non-blocking | Dashboard empty state on first launch | worker-attempt-2 | open |
| Q-ACTOR-11 | actor-journey | non-blocking | Nutrition template expiry — what happens week over week? | worker-attempt-2 | open |
| Q-ACTOR-12 | actor-journey | non-blocking | First-run empty-state convention for all sections | worker-attempt-2 | open |
| Q-ACTOR-13 | actor-journey | non-blocking | PIN lost/forgotten recovery (cross-scope: roles/permissions) | worker-attempt-2 | open |
| Q-ACTOR-14 | actor-journey | non-blocking | ZIP export generation failure (disk, timeout) | worker-attempt-2 | open |
| Q-ACTOR-15 | actor-journey | non-blocking | Backup import fails mid-way — partial data? | worker-attempt-2 | open |
| Q-ACTOR-16 | actor-journey | non-blocking | Media upload failure | worker-attempt-2 | open |
| Q-ACTOR-17 | actor-journey | non-blocking | Mid-session resilience (connection lost, session expiry) (cross-scope: roles/permissions) | worker-attempt-2 | open |
| Q-ACTOR-18 | actor-journey | non-blocking | Navigate away during workout entry — unsaved data? | worker-attempt-2 | open |
| Q-ACTOR-19 | actor-journey | non-blocking | Invalid backup ZIP uploaded | worker-attempt-2 | open |
| Q-ACTOR-20 | actor-journey | non-blocking | Concurrent access (two tabs) | worker-attempt-2 | open |
| Q-ACTOR-21 | actor-journey | non-blocking | Nutrition template after goal change | worker-attempt-2 | open |
| Q-ACTOR-22 | actor-journey | non-blocking | App setup on first launch | worker-attempt-2 | open |
| Q-ACTOR-23 | actor-journey | non-blocking | Goal/context setup before first AI export | worker-attempt-2 | open |
| Q-ACTOR-24 | actor-journey | non-blocking | PIN setup during first run | worker-attempt-2 | open |
| Q-ACTOR-25 | actor-journey | non-blocking | Exercise search/filter for large library | worker-attempt-2 | open |
| Q-ACTOR-26 | actor-journey | non-blocking | Cardio entry without pulse/zone | worker-attempt-2 | open |
| Q-ACTOR-27 | actor-journey | non-blocking | Multiple workout days navigation (list vs calendar) | worker-attempt-2 | open |
| Q-ACTOR-28 | actor-journey | non-blocking | Body weight duplicate on same date | worker-attempt-2 | open |

### Domain Model

| ID | Scope | Severity | Question | Source | Status |
| --- | --- | --- | --- | --- | --- |
| Q-DOMAIN-001 | domain-model | non-blocking | Valid values for BodyWeightEntry.source | prd.md sec 25.9 | open |
| Q-DOMAIN-002 | domain-model | non-blocking | Enum values for heartRateZone | prd.md sec 12.4, 25.8 | open |
| Q-DOMAIN-003 | domain-model | non-blocking | Enum values for cardioType | prd.md sec 12.3, 25.8 | open |
| Q-DOMAIN-004 | domain-model | non-blocking | Enum values for measurementType | prd.md sec 13.3, 25.11 | open |
| Q-DOMAIN-005 | domain-model | non-blocking | Enum values for side | prd.md sec 13.4, 25.11 | open |
| Q-DOMAIN-006 | domain-model | non-blocking | Enum values for flagType | prd.md sec 18.4, 25.18 | open |
| Q-DOMAIN-007 | domain-model | non-blocking | Enum values for mediaType | prd.md sec 11.3, 25.4 | open |
| Q-DOMAIN-008 | domain-model | non-blocking | Enum values for mealLabel | prd.md sec 15.3, 25.15, 25.17 | open |
| Q-DOMAIN-009 | domain-model | non-blocking | Values for Settings.units | prd.md sec 25.1 | open |
| Q-DOMAIN-010 | domain-model | non-blocking | AiExport include flags — model incomplete vs feature description | prd.md sec 17.3, 25.19 | open |

### Feature Behavior

| ID | Scope | Severity | Question | Source | Status |
| --- | --- | --- | --- | --- | --- |
| Q-FEAT-001 | feature-behavior | non-blocking | Dashboard "training days" count — calendar week or trailing 7 days? | prd.md section 9 | open |
| Q-FEAT-002 | feature-behavior | non-blocking | Weekly check-in reminder trigger | prd.md section 9 | open |
| Q-FEAT-003 | feature-behavior | non-blocking | Cardio standalone vs workout day — contradiction | prd.md sections 10.1, 10.3, 25.8 | open |
| Q-FEAT-004 | feature-behavior | non-blocking | Muscle groups representation | prd.md section 11.2 | open |
| Q-FEAT-005 | feature-behavior | non-blocking | Exercise media file size limits and formats | prd.md section 11.3 | open |
| Q-FEAT-006 | feature-behavior | non-blocking | Custom cardio types beyond basic enum | prd.md section 12.3 | open |
| Q-FEAT-007 | feature-behavior | non-blocking | Standalone body measurements without check-in | prd.md section 13 | open |
| Q-FEAT-008 | feature-behavior | non-blocking | Nutrition template week-over-week lifecycle | prd.md sections 15.3, 15.4 | open |
| Q-FEAT-009 | feature-behavior | non-blocking | Charting library and interactivity | prd.md section 16 | open |
| Q-FEAT-010 | feature-behavior | non-blocking | Max AI export ZIP size, large export handling | prd.md section 17 | open |
| Q-FEAT-011 | feature-behavior | non-blocking | CSV encoding and escaping rules | prd.md section 17.8 | open |
| Q-FEAT-012 | feature-behavior | non-blocking | PIN session TTL and refresh | prd.md section 7.2 | open |
| Q-FEAT-013 | feature-behavior | non-blocking | Nutrition template applied mid-week (retroactive or forward?) | prd.md section 15.4 | open |
| Q-FEAT-014 | feature-behavior | non-blocking | Delete exercise with historical data | worker attempt 1 | open |
| Q-FEAT-015 | feature-behavior | non-blocking | Backup import max file size | prd.md section 20.4 | open |
| Q-FEAT-016 | feature-behavior | non-blocking | Valid values for BodyWeightEntry.source | prd.md section 25.9 | open |

### Edge Cases and Risk

| ID | Scope | Severity | Question | Source | Status |
| --- | --- | --- | --- | --- | --- |
| Q-EDGE-01 | edge-case-risk | non-blocking | PIN enabled but pinHash missing | worker-attempt-1 | open |
| Q-EDGE-02 | edge-case-risk | non-blocking | PIN brute-force protection policy | worker-attempt-1 | open |
| Q-EDGE-03 | edge-case-risk | non-blocking | Session TTL for PIN-protected sessions | worker-attempt-1 | open |
| Q-EDGE-04 | edge-case-risk | non-blocking | Workout save interrupted (network/DB) | worker-attempt-1 | open |
| Q-EDGE-05 | edge-case-risk | non-blocking | Duplicate exercise names | worker-attempt-1 | open |
| Q-EDGE-06 | edge-case-risk | non-blocking | Concurrent edit (two tabs) | worker-attempt-1 | open |
| Q-EDGE-07 | edge-case-risk | non-blocking | Import with existing data | worker-attempt-1 | open |
| Q-EDGE-08 | edge-case-risk | non-blocking | EXIF metadata on exported photos | worker-attempt-1 | open |
| Q-EDGE-09 | edge-case-risk | non-blocking | Max file size for media uploads | worker-attempt-1 | open |
| Q-EDGE-10 | edge-case-risk | non-blocking | Backup export timeout for large datasets | worker-attempt-1 | open |
| Q-EDGE-11 | edge-case-risk | non-blocking | Media cleanup on exercise deletion | worker-attempt-1 | open |
| Q-EDGE-12 | edge-case-risk | non-blocking | Migration strategy for schema version changes | worker-attempt-1 | open |

### Acceptance Criteria

| ID | Scope | Severity | Question | Source | Status |
| --- | --- | --- | --- | --- | --- |
| Q-AC-01 | acceptance-criteria | non-blocking | PIN failure behavior (retries, lockout, error display) | worker-attempt-1 | open |
| Q-AC-02 | acceptance-criteria | non-blocking | PIN session TTL | worker-attempt-1 | open |
| Q-AC-03 | acceptance-criteria | non-blocking | "Last body weight" definition on dashboard | worker-attempt-1 | open |
| Q-AC-04 | acceptance-criteria | non-blocking | Week start day and timezone | worker-attempt-1 | open |
| Q-AC-05 | acceptance-criteria | non-blocking | Weekly check-in reminder trigger | worker-attempt-1 | open |
| Q-AC-06 | acceptance-criteria | non-blocking | Working weight auto-populate UI behavior | worker-attempt-1 | open |
| Q-AC-07 | acceptance-criteria | non-blocking | e1RM formula | worker-attempt-1 | open |
| Q-AC-08 | acceptance-criteria | non-blocking | Best set definition | worker-attempt-1 | open |
| Q-AC-09 | acceptance-criteria | non-blocking | Progression signal surfacing in UI | worker-attempt-1 | open |
| Q-AC-10 | acceptance-criteria | non-blocking | Cardio standalone vs day-attached | worker-attempt-1 | open |
| Q-AC-11 | acceptance-criteria | non-blocking | Single nutrition template replacement semantics | worker-attempt-1 | open |
| Q-AC-12 | acceptance-criteria | non-blocking | Nutrition product deletion behavior | worker-attempt-1 | open |
| Q-AC-13 | acceptance-criteria | non-blocking | Chart filter scope | worker-attempt-1 | open |
| Q-AC-14 | acceptance-criteria | non-blocking | Week flags per-week vs per-export | worker-attempt-1 | open |
| Q-AC-15 | acceptance-criteria | non-blocking | Import with existing data (merge/replace/error) | worker-attempt-1 | open |
| Q-AC-16 | acceptance-criteria | non-blocking | Backup CSV files mandatory vs optional | worker-attempt-1 | open |
| Q-AC-17 | acceptance-criteria | non-blocking | "Copy previous set" one-tap vs auto-fill | worker-attempt-1 | open |
| Q-AC-18 | acceptance-criteria | non-blocking | "Sensitive comments" definition | worker-attempt-1 | open |