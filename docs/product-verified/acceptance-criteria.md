# Acceptance Criteria

## Product-Level Criteria

| AC ID | Criterion | Source |
| --- | --- | --- |
| AC-001 | User can enable and disable PIN | §7.2 |
| AC-002 | User can create an exercise with name, muscle groups, working weight | §11.1, §11.2 |
| AC-003 | User can set and modify exercise working weight | §10.6, §11.2 |
| AC-004 | User can upload images and video to an exercise | §11.3 |
| AC-005 | User can open the current day in workout diary | §10.2 |
| AC-006 | User can select a past date via calendar | §10.2 |
| AC-007 | User can add an exercise to a workout day | §10.4 |
| AC-008 | User can add sets with weight and reps | §10.5 |
| AC-009 | User can optionally specify RPE per set | §10.5 |
| AC-010 | User can optionally specify RIR per set | §10.5 |
| AC-011 | User can add a comment to an exercise in a workout | §10.4 |
| AC-012 | User can add cardio with type and duration | §12.2 |
| AC-013 | User can optionally specify pulse and heart rate zone for cardio | §12.2, §12.4 |
| AC-014 | User can create a weekly body check-in | §13.2 |
| AC-015 | Check-in includes date, weight, optional body fat %, measurements, 2-4 photos, comment | §13.2 |
| AC-016 | User can enter body weight for any date | §13.5 |
| AC-017 | User can create a nutrition product with name and KJBJU per 100g | §15.2 |
| AC-018 | User can create a weekly nutrition template | §15.3 |
| AC-019 | User can override nutrition for a specific day | §15.5 |
| AC-020 | User can view an exercise progress chart | §16.2 |
| AC-021 | User can view body weight and measurement charts | §16.3 |
| AC-022 | User can view basic nutrition charts (weekly KJBJU averages) | §16.4 |
| AC-023 | User can generate an AI prompt | §18 |
| AC-024 | User can download an AI export ZIP for the last 4 weeks | §17 |
| AC-025 | User can save an AI review | §19 |
| AC-026 | User can download a full backup ZIP | §20.1 |
| AC-027 | User can import a backup into a clean instance | §20.4 |
| AC-028 | User can run verification command with passing tests and coverage | §24.4 |

## Feature-Level Criteria

### PIN Guard

| AC ID | Criterion | Source |
| --- | --- | --- |
| AC-029 | PIN is optional — user can use app without PIN | §7.2 |
| AC-030 | When PIN is enabled, all pages require PIN entry before access | §7.2 |
| AC-031 | PIN is stored as hash, not plaintext | §7.2 |
| AC-032 | User can change PIN after entering current PIN | §7.2 |
| AC-033 | User can disable PIN from settings | §7.2 |
| AC-034 | PIN session persists via cookie | §7.2 |

### Workout Diary

| AC ID | Criterion | Source |
| --- | --- | --- |
| AC-035 | Diary opens to today's date by default | §10.2 |
| AC-036 | User can navigate to any date via calendar | §10.2 |
| AC-037 | Selecting a date with an existing record opens that record | §10.2 |
| AC-038 | Selecting a date with no record shows an empty slot; first save creates the record | §10.2 |
| AC-039 | Adding an exercise auto-populates working weight from exercise library | §10.6 |
| AC-040 | User can add multiple sets per exercise | §10.5 |
| AC-041 | Working weight snapshot is stored per exercise in workout day | §10.6 |
| AC-042 | Exercise comment is included in AI export | §10.4 |

### Exercise Library

| AC ID | Criterion | Source |
| --- | --- | --- |
| AC-043 | User creates exercises manually (no starter catalog) | §11.1 |
| AC-044 | Exercise includes name, muscle groups, description, working weight | §11.2 |
| AC-045 | User can upload images and video after exercise creation | §11.3 |
| AC-046 | User can delete media from an exercise | §11.3 |
| AC-047 | Exercise can be marked as active or inactive | §11.2 |

### Cardio

| AC ID | Criterion | Source |
| --- | --- | --- |
| AC-048 | User can select cardio type from predefined list | §12.3 |
| AC-049 | Cardio duration is recorded in minutes | §12.2 |
| AC-050 | User can optionally record average pulse | §12.4 |
| AC-051 | User can optionally select heart rate zone (1-5 or unknown) | §12.4 |

### Body Tracking

| AC ID | Criterion | Source |
| --- | --- | --- |
| AC-052 | Weekly check-in records date, weight, optional body fat % | §13.2 |
| AC-053 | Check-in includes 10 measurement types | §13.3 |
| AC-054 | Paired measurements can record left, right, or single (common) value | §13.4 |
| AC-055 | Second value in paired measurement is not required | §13.4 |
| AC-056 | 2-4 photos can be attached to check-in | §13.2 |
| AC-057 | Standalone weight entry can be added for any date | §13.5 |

### Nutrition

| AC ID | Criterion | Source |
| --- | --- | --- |
| AC-058 | User creates products with KJBJU per 100g | §15.2 |
| AC-059 | Weekly template contains products with gram amounts | §15.3 |
| AC-060 | Template items can have optional meal label | §15.3 |
| AC-061 | Template auto-applies to all days of its week | §15.4 |
| AC-062 | Daily override can add, remove, or replace products | §15.5 |
| AC-063 | KJBJU calculated from template item values | §15.3 |
| AC-064 | KJBJU recalculated after daily override change | §15.5 |

### Charts

| AC ID | Criterion | Source |
| --- | --- | --- |
| AC-065 | Exercise chart shows working weight over selected period | §16.2 |
| AC-066 | Exercise chart shows best set per session | §16.2 |
| AC-067 | Exercise chart shows e1RM over time | §16.2 |
| AC-068 | Exercise chart shows volume per session | §16.2 |
| AC-069 | Body weight chart shows weight over period | §16.3 |
| AC-070 | Measurement chart shows individual measurement over period | §16.3 |
| AC-071 | User can overlay multiple measurements on one chart | §16.3 |
| AC-072 | Nutrition chart shows weekly average KJBJU | §16.4 |
| AC-073 | User can select chart period | §16.1 |

### AI Export

| AC ID | Criterion | Source |
| --- | --- | --- |
| AC-074 | Default date range is last 4 weeks | §17.2 |
| AC-075 | User can select custom date range | §17.2 |
| AC-076 | User can toggle which data sections to include | §17.3 |
| AC-077 | Photos excluded from export by default | §17.3 |
| AC-078 | Generated ZIP contains manifest.json, data.json, summary.md | §17.4 |
| AC-079 | Generated ZIP contains CSV files (workouts, measurements, nutrition, cardio) | §17.8 |
| AC-080 | Generated ZIP contains photos/ directory when photos included | §17.4 |
| AC-081 | manifest.json includes export type, schema version, app version, date, period, sections | §17.5 |
| AC-082 | data.json includes all selected data sections | §17.6 |
| AC-083 | summary.md includes period, goal, workout stats, exercise trends, weight/measurement changes, nutrition summary, cardio, comments | §17.7 |

### AI Prompt Builder

| AC ID | Criterion | Source |
| --- | --- | --- |
| AC-084 | Persistent AI context is stored and reused across exports | §18.2 |
| AC-085 | Persistent context includes goal, height, optional age, training experience, split, limitations, progression style, nutrition strategy, persistent comment | §18.2 |
| AC-086 | User can update context at any time | §18.2 |
| AC-087 | User can add one-time comment for a single export | §18.3 |
| AC-088 | User can select week flags for the export period | §18.4 |
| AC-089 | Generated prompt asks AI to analyze progress, compare actual vs working weights, evaluate volume, consider RPE/RIR and cardio, compare training vs body changes, give next-week recommendations | §18.5 |

### AI Review

| AC ID | Criterion | Source |
| --- | --- | --- |
| AC-090 | User can paste AI response text | §19.2 |
| AC-091 | User can link review to a date range | §19.2 |
| AC-092 | User can add notes and planned actions | §19.2 |

### Backup

| AC ID | Criterion | Source |
| --- | --- | --- |
| AC-093 | Full backup ZIP contains manifest.json, data.json, media/ | §20.1 |
| AC-094 | manifest.json includes type, schema version, app version, date, sections | §20.2 |
| AC-095 | data.json includes all entities (settings, profile, exercises, workouts, cardio, body, nutrition, AI) | §20.3 |
| AC-096 | User can include or exclude media from backup | §20.1 |
| AC-097 | Import validates manifest.json structure | §20.4 |
| AC-098 | Import validates schema version | §20.5 |
| AC-099 | Import runs dry-run validation before actual restore | §20.4 |
| AC-100 | Import shows summary before user confirmation | §20.4 |
| AC-101 | Import restores data and media fully, or fails without partial import | §20.4 |
| AC-102 | Import displays clear error messages on validation failure | §20.4 |

### Dashboard

| AC ID | Criterion | Source |
| --- | --- | --- |
| AC-103 | Dashboard shows current date | §9 |
| AC-104 | Dashboard shows last recorded body weight | §9 |
| AC-105 | Dashboard shows training days count for current week | §9 |
| AC-106 | Dashboard shows cardio sessions count for current week | §9 |
| AC-107 | Dashboard shows current goal from user profile | §9 |
| AC-108 | Dashboard shows quick action buttons (add workout, add cardio, add weight, check-in, AI report) | §9 |

## Negative Criteria

| AC ID | Criterion | Source |
| --- | --- | --- |
| AC-109 | PIN change fails without correct current PIN | §7.2 |
| AC-110 | PIN-protected pages inaccessible without valid session | §7.2 |
| AC-111 | Media files return 401/403 without valid session | §14.3, §24.1 |
| AC-112 | Photos not included in AI export unless explicitly opted in | §17.3 |
| AC-113 | Daily override does not affect other dates | §15.5 |
| AC-114 | Backup import fails cleanly with error when manifest is missing/invalid | §20.4 |
| AC-115 | Backup import fails cleanly when schema version is incompatible | §20.5 |
| AC-116 | No silent partial import — incomplete restore fails and reports error | §20.4 |
| AC-117 | PIN not logged in application logs | §24.1 |
| AC-118 | AI export content not logged | §24.1 |
| AC-119 | Photos not logged | §24.1 |
| AC-120 | Sensitive comments not logged | §24.1 |

## Handoff Criteria

| AC ID | Criterion | Source |
| --- | --- | --- |
| AC-121 | User can complete round-trip: create exercise → log workout → add cardio → create weekly check-in → generate AI export | §29 |
| AC-122 | User can enable PIN → close browser → reopen → enter PIN → access data | §29 |
| AC-123 | User can create nutrition template → override one day → verify different values | §29 |
| AC-124 | User can create backup → reset app → import backup → verify data restored | §29 |
| AC-125 | User can run test suite with passing results and coverage gate | §24.4, §29 |