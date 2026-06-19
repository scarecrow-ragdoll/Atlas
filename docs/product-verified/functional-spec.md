# Functional Specification

## Capability Map

| Capability Area | Features |
| --- | --- |
| Access Control | Optional PIN guard with session management |
| Dashboard | Weekly summary with quick actions |
| Exercise Library | CRUD exercises, media upload, working weight |
| Workout Diary | Date-based workout entry, calendar navigation, sets, RPE/RIR |
| Cardio | Type, duration, pulse/zone logging |
| Body Tracking | Weight entries, weekly check-ins, measurements, photos |
| Nutrition | Product catalog, weekly template, daily overrides, KJBJU calculation |
| Charts | Training progress, body measurements, nutrition averages |
| AI Export | Prompt builder, persistent context, week flags, ZIP export |
| AI Review | Manual AI response storage, planned actions |
| Backup/Restore | Full ZIP export, dry-run import, full restore |
| Settings | PIN config, units, user profile |

## Feature Behavior

### PIN Guard (§7.2) — REQ-001
- Optional: enabled/disabled in settings
- When enabled: all pages require valid session
- PIN stored as hash only
- Changeable and removable
- Session via cookie

### Dashboard (§9) — REQ-002
- Shows current date
- Last recorded body weight
- Training days this week count
- Cardio sessions this week count
- Current goal
- Upcoming weekly check-in reminder
- Quick action buttons: add workout, add cardio, add weight, weekly check-in, generate AI report

### Exercise Library (§11) — REQ-003
- User creates exercises with name, muscle groups, description, notes, working weight
- User uploads images and video per exercise
- Media can be added or removed after exercise creation
- Exercise can be active or inactive
- Future built-in exercises (post-MVP) are user-editable like custom ones

### Workout Diary (§10) — REQ-004
- Default view: today's date
- Calendar navigation to any past date
- Backdating: data entry for past dates allowed
- One workout record per date
- Day contains: exercises with sets, cardio, body weight, notes
- Per exercise: linked to library, order, working weight snapshot, notes, sets
- Per set: number, weight, reps, optional RPE, optional RIR, optional comment
- Working weight auto-populated when adding exercise
- Progression display: current working weight, actual weights, best set, volume, e1RM, weekly trend
- Progression signals: stable top-end reps, weight increase, volume increase, stagnation, regression, negative comments

### Cardio (§12) — REQ-007
- Standalone entry by date or optionally attached to workout day
- Fields: date, type (walking/running/bike/elliptical/treadmill/other), duration (minutes), avg pulse, heart rate zone (1-5/unknown), notes

### Body Tracking (§13) — REQ-008, REQ-009
- Weekly check-in: date, weight (optional), body fat % (optional), 2-4 photos, 10 measurements, notes
- Standalone weight entry by any date
- 10 measurement types: neck, shoulders, forearms, biceps, chest, waist, abdomen, hips, thigh, calves
- Paired measurements (forearms, biceps, thigh, calves): left/right values, single value treated as common

### Nutrition (§15) — REQ-010, REQ-011
- Product catalog: name, calories/protein/fat/carbs per 100g, notes
- Weekly template: start date, products with grams, optional meal label, optional notes
- Template auto-applied to all days of the week
- Daily override: add/subtract/replace products for specific date
- KJBJU calculation for template and overridden days

### Charts (§16) — REQ-012
- Training charts per exercise: working weight, best set, e1RM, volume, total reps, working sets count
- Body charts: weight, body fat %, individual measurements, multi-measurement overlay
- Nutrition charts: weekly averages for calories, protein, fat, carbs
- Configurable period with filters

### AI Export (§17-18) — REQ-013, REQ-014
- Configurable date range (default: last 4 weeks)
- Selectable data sections: workouts, exercises, sets, working weights, comments, RPE/RIR, cardio, body weight, measurements, photos (opt-in), nutrition, goal, additional context
- Output: ZIP with manifest.json, data.json, summary.md, workouts.csv, measurements.csv, nutrition.csv, cardio.csv, photos/
- Prompt builder with persistent context (goal, height, age, training experience, split, limitations, progression style, nutrition strategy, persistent comment)
- One-time comment per export
- Week flags: poor sleep, high stress, illness, injury/pain, AAS/cycle, calorie deficit, surplus, maintenance, missed workouts, travel/disrupted routine

### AI Review (§19) — REQ-015
- Manual entry of AI response text
- Linked to date range
- User notes and planned actions
- Review history view

### Backup Import/Export (§20) — REQ-016
- Full backup ZIP: manifest.json, data.json, media/
- Optional CSV files for manual inspection
- Schema version in manifest
- Import: upload ZIP, validate manifest, validate schema version, dry-run validation, show summary, user confirms, full restore

## Validations

Field-level validation rules not specified in PRD. Known validation needs:
- Working weight: numeric, positive
- Sets: weight numeric >= 0, reps positive integer
- Duration: positive integer (minutes)
- Photo count per check-in: 2-4 recommended (requirement vs recommendation unclear)
- Product nutritional values: numeric per 100g
- Amount grams: positive numeric
- Backup ZIP: must contain manifest.json and data.json

## Notifications

- Dashboard shows "upcoming weekly check-in reminder" (§9)
- No proactive notifications (email, push) in MVP
- No reminders configured by user

## Integrations

- No external API integrations in MVP
- AI analysis is manual copy-paste (user sends prompt + ZIP to ChatGPT)
- Technology stack includes go-telegram/bot library but Telegram is explicitly out of MVP scope (§23)
- Apple Health is explicitly out of MVP scope (§22)