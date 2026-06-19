# Domain-Model-Reviewer Worker Attempt 1

## Sources Read
- docs/product/prd.md (1665 lines)

## Source Delta Reviewed
No source delta present. First run.

## Confirmed Facts

### Entities (20 total from sec 25 Data Model Draft)

| # | Entity | Source | Key Attributes | Relationships |
|---|--------|--------|----------------|---------------|
| 1 | Settings | sec 25.1 | id, pinEnabled, pinHash?, defaultAiExportWeeks, units, createdAt, updatedAt | Singleton |
| 2 | UserProfile | sec 25.2 | id, displayName, goal, height, birthDate?, trainingExperience?, currentTrainingSplit?, preferredProgressionStyle?, nutritionStrategy?, persistentAiContext, createdAt, updatedAt | Singleton |
| 3 | Exercise | sec 25.3 | id, name, muscleGroups, description, personalNotes, workingWeight, isActive, createdAt, updatedAt | 1:N → ExerciseMedia |
| 4 | ExerciseMedia | sec 25.4 | id, exerciseId, mediaType, filePath, originalFileName, mimeType, sizeBytes, createdAt | N:1 → Exercise |
| 5 | WorkoutDay | sec 25.5 | id, date, bodyWeight?, notes, createdAt, updatedAt | 1:N → WorkoutExercise |
| 6 | WorkoutExercise | sec 25.6 | id, workoutDayId, exerciseId, order, workingWeightSnapshot, notes, createdAt, updatedAt | N:1 → WorkoutDay, N:1 → Exercise, 1:N → WorkoutSet |
| 7 | WorkoutSet | sec 25.7 | id, workoutExerciseId, setNumber, weight, reps, rpe?, rir?, notes?, createdAt, updatedAt | N:1 → WorkoutExercise |
| 8 | CardioEntry | sec 25.8 | id, workoutDayId?, date, cardioType, durationMinutes, avgPulse?, heartRateZone?, notes, createdAt, updatedAt | N:1 → WorkoutDay (optional) |
| 9 | BodyWeightEntry | sec 25.9 | id, date, weight, source, notes?, createdAt, updatedAt | Independent of check-in |
| 10 | BodyCheckIn | sec 25.10 | id, date, weight?, bodyFatPercentage?, notes, createdAt, updatedAt | 1:N → BodyMeasurement, 1:N → ProgressPhoto |
| 11 | BodyMeasurement | sec 25.11 | id, checkInId, measurementType, side?, value, createdAt, updatedAt | N:1 → BodyCheckIn |
| 12 | ProgressPhoto | sec 25.12 | id, checkInId, filePath, originalFileName, mimeType, sizeBytes, angle?, label?, notes?, createdAt, updatedAt | N:1 → BodyCheckIn |
| 13 | NutritionProduct | sec 25.13 | id, name, caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g, notes?, createdAt, updatedAt | 1:N → NutritionTemplateItem, 1:N → DailyNutritionOverrideItem |
| 14 | NutritionTemplate | sec 25.14 | id, weekStartDate, title, notes?, createdAt, updatedAt | 1:N → NutritionTemplateItem |
| 15 | NutritionTemplateItem | sec 25.15 | id, templateId, productId, amountGrams, mealLabel?, notes?, createdAt, updatedAt | N:1 → NutritionTemplate, N:1 → NutritionProduct |
| 16 | DailyNutritionOverride | sec 25.16 | id, date, notes?, createdAt, updatedAt | 1:N → DailyNutritionOverrideItem |
| 17 | DailyNutritionOverrideItem | sec 25.17 | id, overrideId, productId, amountGrams, operation, mealLabel?, notes?, createdAt, updatedAt | N:1 → DailyNutritionOverride, N:1 → NutritionProduct |
| 18 | WeekFlag | sec 25.18 | id, weekStartDate, flagType, notes?, createdAt, updatedAt | Independent, scoped by week |
| 19 | AiExport | sec 25.19 | id, dateRangeStart, dateRangeEnd, includePhotos, includeNutrition, includeCardio, includeMeasurements, userComment, generatedPrompt, exportFilePath?, createdAt | Snapshot, no live relationships |
| 20 | AiReview | sec 25.20 | id, dateRangeStart, dateRangeEnd, aiResponseText, userNotes?, plannedActions?, createdAt, updatedAt | Standalone, linked by period |

### Identifiers
- All entities use UUID-style `id` as primary key (implied by draft model convention)
- Foreign keys use parent entity name + `Id` suffix (exerciseId, workoutDayId, etc.)

### Statuses & Lifecycle States
- **Exercise.isActive** (boolean): active ↔ inactive transition. No explicit workflow for deactivation. Source: sec 25.3, 11.2.
- **AiExport.exportFilePath**: presence indicates generation completed; absence implies not yet generated or failed. No explicit status enum. Source: sec 25.19.
- No other entities have explicit status/lifecycle fields. WorkoutDay, CardioEntry, NutritionTemplate, etc. are created and updated without formal state transitions.

### Invariants
1. **Unique date per WorkoutDay**: "если за дату уже есть запись, открывается существующая запись" (sec 10.2). Implies unique constraint on WorkoutDay.date.
2. **One NutritionTemplate per week**: weekStartDate appears unique (sec 15.3). Implicit unique constraint.
3. **WorkingWeightSnapshot captures Exercise.workingWeight at creation**: sec 10.6 mandates snapshot, stored in WorkoutExercise.workingWeightSnapshot.
4. **BodyWeightEntry.source has defined values**: field exists (sec 25.9) but values undefined.
5. **DailyNutritionOverrideItem.operation in {add, subtract, replace}**: sec 25.17 explicitly defines three values.
6. **BodyMeasurement.side in {left, right, null}**: "если пользователь указал только одно значение, оно считается общим" (sec 13.4).
7. **ProgressPhoto count 2-4 per BodyCheckIn**: "фотографии 2-4 штуки" (sec 13.2, 26.5 step 7).
8. **No orphan media**: ExerciseMedia.exerciseId, BodyMeasurement.checkInId, ProgressPhoto.checkInId are required FKs (sec 25.4, 25.11, 25.12).
9. **CardioEntry.workoutDayId is optional**: cardio can exist as standalone entry or attached to a workout day (sec 12, 25.8).

## Contradictions
- **AiExport include flags inconsistency**: Sec 25.19 model has `includePhotos`, `includeNutrition`, `includeCardio`, `includeMeasurements`. Sec 17.3 lists 14 toggles (workouts, exercises, sets, working weights, comments, RPE/RIR, cardio, body weight, measurements, photos, nutrition, user goal, additional context). The model covers only 4 of 14. Possible derivation: model uses broad flags while UI uses granular toggles, but this is not documented.
- **BodyCheckIn notes vs WorkoutDay notes vs workout-level comments**: Sec 10.3 includes "общий комментарий к дню" (WorkoutDay.notes). Sec 25.10 BodyCheckIn has "notes". Both exist independently, but PRD doesn't clarify if a daily comment entered during a workout carries over to check-in context.

## Missing Source Artifacts
- No enum definitions for cardioType, heartRateZone, measurementType, mediaType, flagType, mealLabel, side, units, BodyWeightEntry.source
- No formal state machine or workflow model for any entity lifecycle
- No database-level constraint specifications (unique, not-null, cascade)
- No entity relationship diagram

## Derived Requirements

### Derived data fields

| Field | Source | Rationale | Confidence |
|-------|--------|-----------|------------|
| BodyWeightEntry.source enum: check-in, manual, import | sec 25.9 has source field, sec 13.5 distinguishes check-in weight from standalone weight | Field exists; values must cover described scenarios | medium |
| CardioEntry.heartRateZone as enum {zone1..zone5, unknown} | sec 12.4 explicitly lists zones as enum | High - direct source text | high |
| CardioEntry.cardioType as enum {walking, running, bike, elliptical, treadmill, other} | sec 12.3 lists 6 types | high | high |
| BodyMeasurement.measurementType as enum {neck, shoulders, forearms, biceps, chest, waist, abdomen, hips, thigh, calf} | sec 13.3 lists 10 measurements | high | high |
| BodyMeasurement.side as enum {left, right} or null | sec 13.4 describes left/right for paired measurements | high | high |
| ExerciseMedia.mediaType as enum {image, video} | sec 11.3 mentions images and video | medium | high |
| WeekFlag.flagType as enum with ~10 values from sec 18.4 | sec 18.4 lists bad-sleep, high-stress, sickness, injury/pain, AAS-cycle, calorie-deficit, calorie-surplus, maintenance, missed-workouts, travel | high | high |

### Derived relationships
- **CardioEntry can exist without WorkoutDay**: workdayId optional (sec 25.8). Confirmed by sec 26.4 (cardio can be standalone).
- **BodyWeightEntry is independent from BodyCheckIn**: separate entity, can be logged on any date (sec 25.9, 13.5).
- **WorkoutDay.bodyWeight is a convenience field** that duplicates BodyWeightEntry for the same date. No explicit dedup rule documented.

## Missing Information

### Enums not formalized in data model (listed in text only)
1. heartRateZone (sec 12.4) - enum values in text, no field in CardioEntry model sec 25.8
2. cardioType (sec 12.3) - values in text, string field in model
3. measurementType (sec 13.3) - values in text, string field in model (sec 25.11)
4. side (sec 13.4) - implied left/right/null
5. mediaType (sec 11.3) - image/video implied
6. flagType (sec 18.4) - ~10 values in text
7. mealLabel (sec 15.3) - never defined
8. units (sec 25.1) - never defined
9. BodyWeightEntry.source (sec 25.9) - field exists, no values
10. AiExport include flags incomplete vs sec 17.3 feature description

### Missing lifecycle states
- AiExport has no status field (generating/ready/failed). Only tracked via exportFilePath presence.
- Exercise deactivation workflow unspecified (isActive = false implies what? Hide from selector? Keep historical data?).

### Missing constraints
- WorkoutDay.date unique constraint not explicit in model
- NutritionTemplate.weekStartDate unique constraint not explicit
- No cascade delete rules: what happens to WorkoutSet when WorkoutExercise is deleted?
- No deletion or archival rules for any entity

## Open Questions Raised
10 open questions written to question-ledger.md (Q-DOMAIN-001 through Q-DOMAIN-010).

## Edge Cases Or Risks
1. **Empty state**: Exercise without ExerciseMedia is valid. WorkoutDay without WorkoutExercise? PRD says "если за дату записи нет, создаётся новая запись при первом сохранении данных" (sec 10.2) — implies empty WorkoutDay should not be persisted.
2. **Duplicate body weight entries**: Same date can have both BodyWeightEntry and WorkoutDay.bodyWeight. Which takes priority for graphs/AI export? Not documented.
3. **CardioEntry without workoutDayId but on same date as a workout**: PRD considers all activity for a date as one "тренировочная запись" (sec 10.1), but the data model allows standalone CardioEntry.
4. **Template week overlap**: If user creates multiple NutritionTemplates with overlapping weeks, behavior is undefined.
5. **Photo deletion from BodyCheckIn**: If a ProgressPhoto is deleted, the BodyCheckIn continues to exist. Valid.
6. **Exercise deletion impacts historical WorkoutExercise records**: WorkoutExercise references Exercise by exerciseId. If Exercise is deleted, snapshots in historical workouts become orphaned.
7. **AiExport.data.json consistency with data model**: sec 17.6 and 20.3 list entities that must appear in export, but export schema doesn't explicitly map 1:1 with internal entities.

## Recommended Decisions
1. Formalize all text-listed enums as defined types in the data model.
2. Add AiExport.status field (enum: generating, ready, failed) to track generation lifecycle.
3. Add unique constraint on WorkoutDay.date and NutritionTemplate.weekStartDate.
4. Define cascade rules: ExerciseMedia cascade delete with Exercise; WorkoutSet cascade delete with WorkoutExercise; BodyMeasurement/ProgressPhoto cascade delete with BodyCheckIn; NutritionTemplateItem cascade delete with NutritionTemplate; DailyNutritionOverrideItem cascade delete with DailyNutritionOverride.
5. Add Exercise.isActive semantics: inactive exercises hidden from add-exercise selector but preserved in historical workouts and AI export.
6. Expand AiExport include flags to match sec 17.3 toggle list, or document aggregation rule if toggles map to the 4 model-level booleans.
7. Clarify BodyWeightEntry source enum values: at minimum "check-in" and "manual".

## Traceability Candidates
- Entity list → sec 25 Data Model Draft
- Relationships → foreign key fields in sec 25.x
- Invariants → sec 10.2, 10.6, 13.4, 15.3, 17.6, 25.17
- Enum values → sec 12.3, 12.4, 13.3, 13.4, 18.4, 25.17
- User scenarios → sec 26.x for entity relationships in action
- AI export structure → sec 17.4-17.8, 20.1-20.3 for export-focused entity lists