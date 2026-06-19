# User Flows

## Primary Flows

### Add Exercise (§26.1)
1. User opens exercise library
2. Clicks "Add exercise"
3. Enters name, muscle groups, working weight, description, notes
4. Optionally uploads images/video
5. Saves exercise

### Enter Workout For Today (§26.2)
1. User opens workout diary
2. Today's date is shown by default
3. User adds exercise from library (working weight auto-populated)
4. User adds sets with weight and reps
5. Optionally adds RPE/RIR
6. Optionally adds exercise comment
7. Saves the day

### Enter Workout Backdated (§26.3)
1. User opens diary
2. Selects a past date via calendar
3. If record exists, it opens; if not, a new record is created on first save
4. User adds exercises, sets, cardio, comments
5. Saves

### Add Cardio (§26.4)
1. User opens current day or selects date
2. Adds cardio entry
3. Selects type, duration, pulse/zone
4. Optionally adds comment
5. Saves

### Weekly Check-In (§26.5)
1. User opens body measurements section
2. Creates check-in
3. Enters date, weight, optional body fat %
4. Enters measurements (10 types)
5. Adds 2-4 photos
6. Adds comment
7. Saves

### Add Body Weight Standalone (§26.6)
1. User opens weight section or dashboard
2. Adds weight on selected date
3. Saves

### Create Nutrition Template (§26.7)
1. User opens nutrition section
2. Creates products in catalog (if not existing)
3. Creates weekly template
4. Adds products with gram amounts
5. System calculates KJBJU
6. User adjusts to meet goals
7. Saves template (auto-applied to all week days)

### Override Daily Nutrition (§26.8)
1. User selects a date
2. Sees nutrition calculated from template
3. Adds, removes, or modifies products
4. System recalculates KJBJU for the day
5. Changes apply only to selected date

### Generate AI Report (§26.9)
1. User opens AI Export
2. Selects date range (default: 4 weeks)
3. Selects data sections to include
4. Adds one-time comment
5. Reviews saved goal/context
6. System generates prompt + ZIP export
7. User downloads file
8. User sends prompt + file to ChatGPT

### Save AI Review (§26.10)
1. User receives AI response
2. Opens AI Review section
3. Pastes response text
4. Links to date range
5. Adds notes and planned actions
6. Saves

### Full Backup Export (§26.11)
1. User opens Import/Export
2. Clicks "Export all data"
3. Optionally includes media
4. System creates ZIP backup
5. User saves file locally

### Restore From Backup (§26.12)
1. User deploys new Atlas instance
2. Opens Import/Export
3. Uploads ZIP backup
4. System runs dry-run validation
5. Shows import summary
6. User confirms import
7. System restores data and media

## Alternative Flows

| Scenario | Alternative Path |
| --- | --- |
| Exercise not in library | User must exit workout, create exercise, return to workout |
| Invalid set values (0 reps, 0 weight) | Validation error behavior unspecified |
| Check-in with < 2 photos | Photo requirement vs recommendation ambiguous (§13.2) |
| Last exercise removed from day | Day becomes empty — deletion behavior unspecified |
| Same exercise twice in one day | Duplicate handling unspecified |
| Nutrition override reversion | No "reset to template" action defined |

## Failure And Recovery Flows

| Failure Mode | Current Coverage |
| --- | --- |
| PIN forgotten | No recovery mechanism specified (Q-ACTOR-13) |
| ZIP export fails (disk/timeout) | User feedback unspecified |
| Import fails mid-restore | Partial data risk — no rollback specified |
| Media upload fails | No retry or error feedback specified |
| Session lost during data entry | Data loss risk — no autosave specified |
| Invalid backup uploaded | Validation exists but error format unspecified |
| Concurrent tab access | Last-write-wins or conflict detection undefined |

## Empty States

The PRD defines no empty state behavior. All sections (dashboard, exercise library, workout diary, body measurements, progress photos, nutrition, AI exports, AI reviews, charts) have no specified first-run or empty-data state. A consistent first-run convention is needed (Q-ACTOR-12).

Recommended convention: First launch shows dashboard with "Getting started" guidance — create your first exercise, log your first workout.