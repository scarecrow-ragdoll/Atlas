# Feature-Behavior Worker Attempt 1

## Sources Read
- docs/product/prd.md (full document, 1665 lines)

## Source Delta Reviewed
No source delta present.

## Confirmed Facts

### Dashboard (Section 9)
- Shows weekly summary with date, last weight, training days count, cardio count, current goal, nearest check-in reminder.
- Quick actions: add today's workout, add cardio, add weight, open weekly check-in, generate AI report.

### Workout Diary (Section 10)
- One workout record per date (daily activity model, section 10.1).
- Default open to current day; calendar date picker; backfill supported (section 10.2).
- Existing record opens; new record created on first save (section 10.2).
- Fields per workout day: date, exercise list, sets per exercise, cardio, day comment, optional body weight (section 10.3).
- Per exercise: reference to exercise catalog, order, working weight snapshot, user comment, set list (section 10.4).
- Per set: number, weight, reps, optional RPE, optional RIR, optional comment (section 10.5).
- RPE/RIR fully optional (section 10.5).
- Working weight: stored in exercise catalog, auto-substituted on add, actual set weights stored separately, snapshot captured at workout time (section 10.6).
- Progression tracking: current working weight, actual weights for period, best set, volume, e1RM, weekly dynamics (section 10.7).
- Progression signals: stable top-end reps, weight increase, volume increase, stagnation, regression, frequent negative comments (section 10.7).
- No automatic working weight change without user confirmation (section 10.7).
- Out of MVP: workout templates, workout planning, repeat last workout, training programs, split training, future calendar (section 10.8).

### Exercise Library (Section 11)
- User creates exercises; no starter catalog in MVP (section 11.1).
- Fields: name, muscle groups, description, personal notes, working weight, images, video, active status (section 11.2).
- Media: image/video upload per exercise, addable after creation, deletable, included in full backup, not required in AI export (section 11.3).
- Future starter catalog: standard exercises behave identically to user exercises; user can edit, add notes, change working weight (section 11.4).

### Cardio (Section 12)
- Fields: date, type, duration (minutes), pulse, heart rate zone, comment (section 12.2).
- Basic types enum: walking, running, exercise bike, elliptical, treadmill, other (section 12.3).
- Heart rate zone enum: Zone 1-5 + unknown (section 12.4).
- User can enter only pulse or leave empty if zone unknown (section 12.4).

### Body Measurements (Section 13)
- Weekly check-in: date, weight, optional body fat %, 2-4 photos, body measurements, comment (section 13.2).
- Measurement list: neck, shoulders, forearms, biceps, chest, waist, abdomen, hips, thigh, calves (section 13.3).
- Paired measurements (forearms, biceps, thigh, calves): left/right side, single value treated as common, second value not required (section 13.4).
- Body weight: can be entered standalone on any date, part of check-in or separate, included in charts and AI export (section 13.5).

### Progress Photos (Section 14)
- Attached to weekly check-in.
- Fields: date, check-in link, file, optional label, optional angle, optional comment (section 14.1).
- Angles: front, side, back, custom (section 14.2).
- Storage: not publicly accessible without PIN session, included in full backup, included in AI export only if user opts in (section 14.3).

### Nutrition (Section 15)
- Weekly template model: user defines typical day, each day uses template by default (section 15.1).
- Product catalog: user-created, fields: name, calories/100g, protein/100g, fat/100g, carbs/100g, optional comment (section 15.2).
- Template: one per week, list of products with grams, optional meal label, optional comment (section 15.3).
- App calculates: calories, protein, fat, carbs (section 15.3).
- Daily override: add, remove, change product quantity, comment; affects only selected date (section 15.5).
- Out of MVP: recipes, prepared meals, barcode scanner, water, fiber, salt, sugar, alcohol tracking, auto food recognition, public food database (section 15.6).

### Charts (Section 16)
- User selects period, filters data (section 16.1).
- Exercise charts: working weight, best set, e1RM, total volume, total reps, working sets count; user selects exercise + period (section 16.2).
- Body charts: body weight, body fat %, individual measurements, multiple measurements on one chart (section 16.3).
- Nutrition charts: weekly averages for calories, protein, fat, carbs (section 16.4).

### AI Export (Section 17)
- Configurable period, default 4 weeks, arbitrary start/end dates (section 17.2).
- User selects included sections: workouts, exercises, sets, working weights, comments, RPE/RIR, cardio, body weight, measurements, photos, nutrition, goal, additional context (section 17.3).
- Photos excluded by default (section 17.3).
- Format: ZIP archive with manifest.json, data.json, summary.md, workouts.csv, measurements.csv, nutrition.csv, cardio.csv, photos/ (section 17.4).
- manifest.json: export type, schema version, app version, export date, period, included sections, file list, photo presence (section 17.5).
- data.json: profile, goals, exercises, workouts, workoutExercises, sets, cardio, bodyWeightEntries, bodyCheckIns, measurements, nutritionProducts, nutritionTemplates, dailyNutritionOverrides, userComments, computedSummary (section 17.6).
- summary.md: period, goal, training stats, exercise dynamics, weight changes, measurement changes, nutrition summary, cardio, user comments (section 17.7).
- CSV files: workouts.csv, measurements.csv, nutrition.csv, cardio.csv (section 17.8).

### AI Prompt Builder (Section 18)
- Persistent context: goal, height, optional age, training experience, optional training split, limitations/injuries, preferred progression style, nutrition strategy, additional persistent comment (section 18.2).
- User can update goal/context anytime; goal remembered across exports (section 18.2).
- One-time comment per export (section 18.3).
- Week flags: poor sleep, high stress, illness, injury/pain, AAS/cycle context, calorie deficit, calorie surplus, maintenance, missed workouts, travel/disrupted routine (section 18.4).
- Flags included in prompt and data.json (section 18.4).
- Prompt instructs AI to: analyze exercise progress, evaluate working weight dynamics, compare actual vs working weights, evaluate volume, consider RPE/RIR and comments, consider cardio, cross-reference training with weight/measurement trends, consider nutrition, give next-week recommendations, suggest working weight changes, flag exercises for increase/repeat/overload risk, provide action plan (section 18.5).

### AI Review History (Section 19)
- Save AI response text, link to period, add comment, mark planned actions (section 19.2).
- No automatic OpenAI/ChatGPT integration in MVP (section 19.2).

### Import/Export (Section 20)
- Full backup ZIP: manifest.json, data.json, media/ (section 20.1).
- Optional CSV: exercises.csv, workouts.csv, sets.csv, body_measurements.csv, nutrition.csv (section 20.1).
- manifest.json includes: export type (full_backup), schemaVersion, appVersion, exportedAt, includedSections, mediaIncluded, file list, optional checksums (section 20.2).
- data.json includes all entities (section 20.3).
- Import: ZIP upload, manifest check, schema version check, data.json structure check, dry-run validation, summary display, clear errors, no silent partial import, media restore, entity relationship restore (section 20.4).
- Schema versioning for backward compatibility (section 20.5).

### PIN Guard (Section 7.2)
- Optional, disabled by default.
- If enabled: PIN entry on app open, not stored in plain text, changeable, disableable.
- Session via cookie after PIN entry; sensitive data not accessible without valid session.

### Privacy & Non-Functional (Section 24)
- No PIN logging, no AI export logging, no photo logging, no sensitive comment logging.
- Media not served without auth/PIN session.
- Full backup and AI export generated on user request only.
- Data not lost on container restart; media in volume.

## Contradictions

1. **Workout-day model vs. CardioEntry independence (sections 10.1, 10.3, 12.1, 25.8).** Section 10.1 states all activity on a date belongs to one workout record, and section 10.3 lists cardio as a field within the workout day. However, the data model draft (25.8) defines CardioEntry with its own `date` field and an optional `workoutDayId`, implying cardio can exist standalone. The PRD does not clarify whether cardio can be entered completely independently of a workout day (e.g., a cardio-only day that has no workout record).

2. **BodyWeightEntry source field (section 25.9).** The data model includes a `source` field on BodyWeightEntry, but the product description (section 13.5) does not define what source values exist. Possible derivations: check-in, standalone entry, import.

3. **Nutrition template week start (sections 15.3, 15.4).** The template has a `weekStartDate` field, but the PRD does not specify whether the template auto-advances to the next week, whether the user must manually create a new template each week, or what happens when the week ends and no new template exists.

## Missing Source Artifacts

1. **No API contract or GraphQL schema.** The PRD references gqlgen in the stack but does not define any API operations, mutations, queries, or resolver behavior. This affects every feature's backend contract.

2. **No UI wireframes, mockups, or page layouts.** Dashboard blocks, workout entry forms, chart views, check-in screens — none are specified beyond textual descriptions.

3. **No validation rules.** No field-level validation for required/optional fields, value ranges, uniqueness, format, or cross-field constraints. For example: is `reps` required to be positive? Is `weight` required to be positive? Is `durationMinutes` required to be > 0? Can `muscleGroups` be free text or an enum?

## Derived Requirements

### Derived Fields
| Field | Source Signal | Derivation Rationale | Confidence |
|---|---|---|---|
| BodyWeightEntry.source enum | Data model (25.9) includes `source` field | Must document valid source values from described flows: check-in, standalone, import | Medium |
| NutritionTemplate.weekStartDate behavior | Section 15.3, 25.14 | Template has weekStartDate; implicit lifecycle: active during that week, what happens after? | Medium |
| DailyNutritionOverrideItem.operation enum | Section 25.17 defines add/subtract/replace | Derived from described override behavior (15.5): add product, remove product, change quantity | High |

### Derived Behaviors
| Behavior | Source Signal | Derivation Rationale | Confidence |
|---|---|---|---|
| Chart period defaults to last 4 weeks | Section 16.1, consistent with AI export default (17.2) | Cross-reference: period filter behavior implied for all charts | High |
| AI export ZIP must be downloadable via browser | Section 26.9 step 9: user downloads file | Standard web download behavior derived from "скачивает файл" | High |
| Backup import requires confirmation after dry-run | Section 26.12 steps 4-6: dry-run then user confirms | Clear flow described | High |

## Missing Information

### Feature-Level Gaps
1. **Dashboard (Section 9):** No specification of how "количество тренировочных дней за неделю" and "количество кардио за неделю" are counted (current calendar week? trailing 7 days?). No definition of "ближайшее напоминание о weekly check-in" — what triggers this? No specification of what data refreshes the dashboard (real-time? on page load?).

2. **Workout Diary (Section 10):** No search/filter on exercises within a workout. No specification of what happens if user adds the same exercise twice in one day. How is exercise order managed (drag-and-drop? numeric input?).

3. **Exercise Library (Section 11):** No search, sort, or filter for exercise list. No specification of how muscle groups are represented (enum? free text? multi-select?). How many images/videos per exercise? File size limits? Accepted formats?

4. **Cardio (Section 12):** Can user add custom cardio types beyond the listed enum? Are cardio-only days allowed (workout without any strength exercises)? No validation on duration (max/min).

5. **Body Measurements (Section 13):** Can measurements be entered standalone without a full check-in? The check-in is the primary scenario, but the data model allows measurements per checkInId.

6. **Progress Photos (Section 14):** File size limits for photos? Accepted image formats? Video in progress photos? What happens when storage is full?

7. **Nutrition (Section 15):** What happens to the nutrition template when the week ends? Does it auto-apply the same template to the next week until changed? Can user have multiple templates? What if no template exists for the current week — does the user see zero nutrition data?

8. **Charts (Section 16):** What charting library is used? What interactivity is required (zoom, tooltip, legend toggle, data point selection, export to image)? Are charts responsive?

9. **AI Export (Section 17):** What is the maximum size for ZIP export (especially with photos)? What encoding for CSV files? What CSV quoting/escaping rules? What happens if export takes a long time — background job or synchronous? Where is the ZIP file stored — in-memory, temp directory, or database? Is there cleanup after download?

10. **AI Prompt Builder (Section 18):** Is the prompt template editable by the user? Is the prompt localized or always in one language?

11. **AI Review History (Section 19):** Can user view past AI reviews? Is there a list view? Search?

12. **Import/Export (Section 20):** What is the maximum upload file size for import? How are checksums computed (algorithm)? What happens on partial import failure — rollback? How long are temporary files kept?

13. **PIN Guard (Section 7.2):** Session TTL? Refresh behavior? Cookie name/domain? What happens if cookie expires mid-use? Logout mechanism?

14. **Settings (Section 25.1):** What are the `units` options (metric only? imperial?)? Is `defaultAiExportWeeks` configurable beyond the default of 4?

### Notifications
No notification system is described beyond the dashboard "reminder about weekly check-in." The PRD does not specify:
- How the reminder is delivered (in-app only? email? push?).
- Whether there are other notifications (e.g., "you haven't logged training in 3 days").
- Whether the check-in reminder is configurable (every N days? specific day of week?).

### Integration Behavior
- No integration with external services in MVP (Apple Health and Telegram are explicitly out of scope).
- No OpenAI/ChatGPT API integration (manual copy by user).
- No cloud backup integration.
- The "AI prompt" is generated and given to the user, who manually sends it to ChatGPT — this is not an API integration.

## Open Questions Raised

See question-ledger.md for full list with IDs Q-FEAT-001 through Q-FEAT-016.

## Edge Cases or Risks

1. **Race condition on workout day creation:** If workout day is created on first save (section 10.2), what happens if the tab is closed before save? User sees no record, reopens, creates again — but first save on re-open should create a new record.
2. **Exercise deletion with history:** What happens if user deletes an exercise that has workout history? Cascade delete? Soft delete? Block deletion?
3. **Working weight change retroactive effect:** If user changes working weight on an exercise, do historical workout snapshots remain as-is? Section 10.6 says snapshot is captured at workout time, implying no retroactive update.
4. **Nutrition template mid-week creation:** If user creates a template on Wednesday, does it apply to Mon-Wed retroactively or only Thu-Sun forward?
5. **Large photo volumes:** Hundreds of photos over years. ZIP export with photos could be very large. Browser download of multi-GB files may fail.
6. **Import from different schema version:** Section 20.5 mentions schema versioning but no migration/upgrade path if import version does not match app version.
7. **Concurrent cardiorespiratory and strength training on same date:** Section 10.3 includes cardio in workout day, but cardio might be logged independently. If both exist, is it duplicated in export?
8. **Empty state flows:** No specification of what the user sees when they have zero workouts, zero exercises, zero nutrition data, etc.

## Recommended Decisions

1. Document the cardio relationship: clarify whether cardio is always part of a workout day or can be standalone.
2. Document the BodyWeightEntry.source enum values.
3. Document nutrition template lifecycle across weeks.
4. Document exercise catalog search/filter behavior.
5. Document charting library decision.
6. Document session TTL and refresh behavior for PIN guard.
7. Document whether the user can add custom cardio types.
8. Document whether measurements can be standalone or only within check-ins.

## Traceability Candidates
- Every feature from Section 8 maps to a traceability entry.
- Every workflow from Section 26 maps to a user flow.
- Acceptance criteria from Section 29 map to AC-001 through AC-026.