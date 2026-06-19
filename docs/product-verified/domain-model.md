# Domain Model

## Entities

20 entities identified from PRD §25 data model draft.

| Entity | Description | Key Identifier |
| --- | --- | --- |
| Settings | Application configuration | id (userId FK) |
| UserProfile | User personal data, goal, AI context | id (userId FK) |
| DefaultUser | Default system user created at bootstrap | id |
| Exercise | Exercise definition in library | id (userId FK) |
| ExerciseMedia | Images/video attached to exercise | id (userId FK) |
| DailyLog | Daily aggregate activity record for a date | id (userId + date unique) |
| WorkoutExercise | Exercise instance within a daily log | id (userId FK) |
| WorkoutSet | Individual set within a workout exercise | id |
| CardioEntry | Cardio session record attached to daily log | id (dailyLogId required) |
| BodyWeightEntry | Standalone body weight record | id (userId FK) |
| BodyCheckIn | Weekly check-in record | id (userId FK) |
| BodyMeasurement | Measurement within a check-in | id |
| ProgressPhoto | Photo within a check-in | id |
| NutritionProduct | Food item in product catalog | id (userId FK) |
| NutritionTemplate | Weekly meal plan template | id (userId FK) |
| NutritionTemplateItem | Product entry in template | id |
| DailyNutritionOverride | Daily nutrition deviation from template | id (userId FK) |
| DailyNutritionOverrideItem | Product change within override | id |
| WeekFlag | Weekly context flag for AI | id (userId FK) |
| AiExport | Generated AI export record | id (userId FK) |
| AiReview | Saved AI analysis/review record | id (userId FK) |

## Attributes

### Settings
- id, pinEnabled, pinHash (optional), defaultAiExportWeeks, units, createdAt, updatedAt

### UserProfile
- id, displayName, goal, height, birthDate (optional), trainingExperience (optional), currentTrainingSplit (optional), preferredProgressionStyle (optional), nutritionStrategy (optional), persistentAiContext, createdAt, updatedAt

### Exercise
- id, name, muscleGroups, description, personalNotes, workingWeight, isActive, createdAt, updatedAt

### ExerciseMedia
- id, exerciseId, mediaType, filePath, originalFileName, mimeType, sizeBytes, createdAt

### DefaultUser
- id, displayName, createdAt, updatedAt

### DailyLog (replaces WorkoutDay)
- id, userId, date (unique per user), notes (optional), bodyWeight (optional), createdAt, updatedAt

### WorkoutExercise
- id, userId, dailyLogId, exerciseId, order, workingWeightSnapshot, notes, createdAt, updatedAt

### WorkoutSet
- id, workoutExerciseId, setNumber, weight, reps, rpe (optional), rir (optional), notes (optional), createdAt, updatedAt

### CardioEntry
- id, userId, dailyLogId, cardioType, durationMinutes, avgPulse (optional), heartRateZone (optional), notes (optional), createdAt, updatedAt

### BodyWeightEntry
- id, date, weight, source (enum undefined), notes (optional), createdAt, updatedAt

### BodyCheckIn
- id, date, weight (optional), bodyFatPercentage (optional), notes, createdAt, updatedAt

### BodyMeasurement
- id, checkInId, measurementType (enum undefined), side (optional, left/right/null), value, createdAt, updatedAt

### ProgressPhoto
- id, checkInId, filePath, originalFileName, mimeType, sizeBytes, angle (optional), label (optional), notes (optional), createdAt, updatedAt

### NutritionProduct
- id, name, caloriesPer100g, proteinPer100g, fatPer100g, carbsPer100g, notes (optional), createdAt, updatedAt

### NutritionTemplate
- id, weekStartDate, title, notes (optional), createdAt, updatedAt

### NutritionTemplateItem
- id, templateId, productId, amountGrams, mealLabel (optional), notes (optional), createdAt, updatedAt

### DailyNutritionOverride
- id, date, notes (optional), createdAt, updatedAt

### DailyNutritionOverrideItem
- id, overrideId, productId, amountGrams, operation (add/subtract/replace), mealLabel (optional), notes (optional), createdAt, updatedAt

### WeekFlag
- id, weekStartDate, flagType (enum undefined), notes (optional), createdAt, updatedAt

### AiExport
- id, dateRangeStart, dateRangeEnd, includePhotos, includeNutrition, includeCardio, includeMeasurements, userComment, generatedPrompt, exportFilePath (optional), createdAt

### AiReview
- id, dateRangeStart, dateRangeEnd, aiResponseText, userNotes (optional), plannedActions (optional), createdAt, updatedAt

## Relationships

- Exercise 1:N ExerciseMedia
- DailyLog 1:N WorkoutExercise (via dailyLogId)
- WorkoutExercise 1:N WorkoutSet
- WorkoutExercise N:1 Exercise (via exerciseId)
- DailyLog 0:N CardioEntry (via dailyLogId, required)
- BodyCheckIn 1:N BodyMeasurement
- BodyCheckIn 1:N ProgressPhoto
- NutritionTemplate 1:N NutritionTemplateItem
- NutritionTemplateItem N:1 NutritionProduct
- DailyNutritionOverride 1:N DailyNutritionOverrideItem
- DailyNutritionOverrideItem N:1 NutritionProduct
- WorkoutDay 0:1 BodyWeightEntry (bodyWeight field on day, also standalone entries by date)

## Lifecycle States

| Entity | States | Source Evidence |
| --- | --- | --- |
| Exercise | active / inactive | PRD §11.2: `isActive` field |
| AiExport | draft (no exportFilePath) / generated (exportFilePath present) | PRD §25.19: exportFilePath optional |
| AiReview | created with response text | PRD §19.2: simple create/read |
| NutritionTemplate | active for its week / superseded by new template | Derived: single template at a time per §15.3 |
| DailyNutritionOverride | exists only when template is overridden | PRD §15.5: override affects single date |

No formal state machines documented. WorkoutDay has implicit states: empty (no exercises saved) / active (has exercises).

## Invariants

1. **Unique DailyLog per user per date**: One daily activity record per user per date (DEC-009, resolved Q-SCOPE-005)
2. **Cardio must belong to a DailyLog**: dailyLogId is required on CardioEntry; system auto-creates DailyLog when needed (DEC-009)
3. **Working weight snapshot**: WorkingWeightSnapshot in WorkoutExercise captures value at execution time, independent of Exercise.workingWeight changes (PRD §10.6)
3. **Single nutrition template**: At most one active template at a time (PRD §15.3)
4. **Check-in photo range**: 2-4 photos per check-in (PRD §13.2)
5. **Set order**: WorkoutSet.setNumber determines order within WorkoutExercise (PRD §10.5)
6. **Exercise order**: WorkoutExercise.order determines display order within WorkoutDay (PRD §10.4)
7. **Body weight optionality**: BodyCheckIn.weight is optional (PRD §13.2)
8. **Override operation constraint**: DailyNutritionOverrideItem.operation is add/subtract/replace (PRD §25.17)
9. **BodyMeasurement side**: Paired measurements (forearm, biceps, thigh, calf) may have left/right/nulled side (PRD §13.4)
10. **AiExport includePhotos defaults to false**: Photos excluded from AI export by default (PRD §17.3)