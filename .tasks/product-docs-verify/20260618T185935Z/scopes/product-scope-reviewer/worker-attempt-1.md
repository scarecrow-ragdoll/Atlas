# Product Scope Reviewer Worker Attempt 1

## Sources Read

- docs/product/prd.md (1665 lines, complete)

## Source Delta Reviewed

No source delta present. First run.

## Confirmed Facts

### Product Identity
- Name: Atlas
- Type: Self-hosted web application
- Format: Single user, no registration, no multi-user in MVP
- GitHub description provided

### Core Purpose
Atlas is a personal fitness diary that tracks workouts, nutrition, body measurements, progress photos, and prepares structured AI-export for analysis. The core weekly cycle: track during week → body check-in → 4-week data aggregation → AI prompt + file generation → user sends to ChatGPT/other AI → AI gives recommendations.

### MVP Sections (11 total)
1. Dashboard
2. Workout diary
3. Exercise library
4. Cardio
5. Body measurements
6. Progress photos
7. Nutrition
8. Charts
9. AI export / AI prompt builder
10. Import/Export
11. Settings

### MVP Features (explicit list in section 27)
Single-user mode, optional PIN, exercise library with media, workout diary by date, backdating workouts, sets with weight/reps, optional RPE/RIR, exercise comments, working weight, cardio, body check-ins, individual weight entries, body measurements, progress photos, product database, weekly nutrition template, daily nutrition overrides, workout charts, body charts, basic nutrition charts, AI prompt builder, AI export ZIP, AI review history, full backup export/import, tests for key scenarios.

### Explicit Non-Goals / Out of Scope (section 28)
Registration, multi-user, roles, SaaS mode, public pages, workout templates, training planning, quick repeat of past workout, prebuilt exercise database, recipes/meals, barcode scanner, Apple Health integration, cloud backup, Telegram bot, OpenAI API integration, automatic plan generation, mobile app.

### Tech Stack
Bun, Node, Go, Docker Compose, Nx, Next.js, Vite, React, gqlgen, PostgreSQL, Redis, chi router, GraphQL, Tailwind, shadcn, vitest, Playwright. 100% coverage gate via `bun run verify:coverage`.

### Acceptance Criteria (section 29)
26 explicit criteria listed (AC-1 through AC-26). Cover PIN, exercises, workout diary, cardio, body check-in, weight entries, nutrition, charts, AI export, AI review, backup/restore, test gate.

### Development Epics (section 30)
9 epics: Foundation, Exercise Library, Workout Diary, Cardio and Body Tracking, Nutrition, Charts, AI Export and Prompt Builder, AI Review History, Backup Import/Export.

### Data Model Draft (section 25)
20 entities defined with fields: Settings, UserProfile, Exercise, ExerciseMedia, WorkoutDay, WorkoutExercise, WorkoutSet, CardioEntry, BodyWeightEntry, BodyCheckIn, BodyMeasurement, ProgressPhoto, NutritionProduct, NutritionTemplate, NutritionTemplateItem, DailyNutritionOverride, DailyNutritionOverrideItem, WeekFlag, AiExport, AiReview.

### Key User Flows (section 26)
12 flows: Add exercise, log today's workout, backdate workout, add cardio, weekly check-in, add weight separately, create nutrition template, change single-day nutrition, generate AI report, save AI review, full backup, restore from backup.

## Contradictions

1. **Tech stack vs MVP scope conflict**: The PRD lists `go-telegram/bot` in the Technology Stack (section 5) but states explicitly that Telegram bot is not in MVP (section 23) and is future scope only. This is not a contradiction per se (it's listed as future scope), but it creates noise in the stack specification for MVP developers who may wonder why a dependency is included.

2. **"Один пользователь" vs future multi-user architecture**: Section 4 says "один пользователь" and "без multi-user режима в MVP", but also says "в будущем архитектура может быть расширена под multi-user". This creates a tension — should the MVP architecture be built with multi-user considerations now, or purely single-user? Unclear.

3. **Workout model vs real-world training patterns**: Section 10.1 states all exercises for a date are one workout record. But the example (morning cardio + evening strength) conflates cardio and strength into one record, while cardio has its own separate section/entity (CardioEntry with optional workoutDayId). This creates ambiguity about whether cardio lives inside the workout day or is a separate parallel concept.

4. **Working weight snapshot timing**: Section 10.6 says working weight is stored in exercise library and snapshot is saved in WorkoutExercise at time of workout. Section 10.4 says "рабочий вес на момент выполнения". However, there's no clarity on what happens if the user changes the working weight in the exercise library after a workout — does the snapshot remain correct? This is implicit from the snapshot design but not explicitly stated.

## Missing Source Artifacts

- **No success metrics** defined anywhere. This is the most significant gap for a scope review.
- **No UX/design specifications** — wireframes, mockups, or design system requirements.
- **No deployment guide or infrastructure requirements** beyond Docker Compose.
- **No data retention policy** — how long data is kept, what happens to media files.
- **No privacy policy or compliance artifacts** despite the product storing sensitive health/fitness data.
- **No glossary/terminology document** — key terms (RPE, RIR, working weight, etc.) are explained inline but not in a centralized glossary.
- **No error handling specification** — beyond implicit PIN guard and import validation.

## Derived Requirements

1. **Success metrics must be defined before handoff** — Source: absence of any success indicators in entire PRD. Rationale: Without success metrics, the product team cannot determine when the MVP is successful, only when it's technically complete. Confidence: high.

2. **Single-user data model should not pre-allocate multi-user infrastructure** — Source: sections 4 and 28 both state single-user MVP. Rationale: Multi-user mentions in the PRD create ambiguity about architectural investment trade-offs. A clear decision is needed. Confidence: medium.

3. **Target user technical proficiency assumption** — Source: self-hosted requirement with Docker Compose. Rationale: Users must have Docker/self-hosting knowledge. This is an implicit requirement for the target audience. Confidence: high.

## Missing Information

1. **No success metrics or KPIs** — The most significant gap. No quantitative or qualitative measures of product success.
2. **No definition of target user's technical expertise** — Self-hosted requires Docker knowledge; is the user expected to be technical?
3. **No performance SLAs** — "быстрым на персональном объёме данных" is vague. No specific page load targets, query response times, or export generation time expectations.
4. **No design system or UI guidelines** — Charts are described in prose but no design specifications.
5. **No backup frequency recommendations** — Not specified how often users should back up.
6. **No definition of "быстрым" (fast) in UX** — Section 24.3 says interface should remain "fast" but without quantitative targets.

## Open Questions Raised

| ID | Question | Why It Matters |
|----|----------|----------------|
| Q-SCOPE-001 | What are the quantitative success metrics for the MVP? | Without success metrics, there is no way to validate the product is achieving its goals |
| Q-SCOPE-002 | Should the MVP architecture build for future multi-user or remain strictly single-user? | Affects data model, auth design, API structure, and database schema decisions |
| Q-SCOPE-003 | What is the target user's expected technical proficiency? | Self-hosted requires Docker/CLI knowledge — affects documentation and deployment UX requirements |
| Q-SCOPE-004 | What are the specific performance targets (page load, export time, chart render time)? | Section 24.3 uses vague language ("быстрым") with no measurable targets |
| Q-SCOPE-005 | Should cardio be a separate entity or always part of a workout day? | The data model shows cardio as separate with optional workoutDayId, but section 10.3 includes cardio in the workout day — ambiguous |

## Edge Cases Or Risks

1. **Scope creep risk**: The PRD extensively discusses future features (multi-user, Apple Health, Telegram bot, mobile app, coach mode, cloud backup). While section 28 explicitly excludes them from MVP, the volume of future scope discussion could lead to architectural over-engineering unless disciplined.

2. **Single point of failure**: Self-hosted single-user means all data is lost if the user's Docker volume or database is corrupted without a backup. The backup/restore feature mitigates this but depends on the user proactively using it.

3. **AI export dependency risk**: The core value proposition depends on the user having access to ChatGPT or another AI. If AI access changes (paywall, availability, API changes), the primary workflow is broken.

4. **Photo storage growth**: Hundreds of photos over years of training could lead to significant storage requirements. No guidance on storage limits, compression, or cleanup.

5. **Nutrition template complexity**: The weekly template + daily override model is elegant for regular eaters but may be too restrictive for users with highly variable daily diets.

## Recommended Decisions

1. **Define MVP success metrics before development**: At minimum: (a) user can complete the full weekly cycle without data loss, (b) AI export generates valid ZIP in under X seconds for 4 weeks of data, (c) user can enter a workout in under Y minutes.

2. **Decide single-user architectural boundary**: Document whether future multi-user compatibility is a design concern during MVP or explicitly deferred.

3. **Define performance targets**: At minimum: page load < 2s, AI export generation < 5s for 4 weeks of data, chart render < 1s.

## Traceability Candidates

| Claim | Source Reference |
|-------|-----------------|
| Product name, type, format | prd.md sections 1, 3, 4 |
| MVP sections | prd.md section 8 |
| MVP features | prd.md section 27 |
| Non-goals | prd.md section 28 |
| Acceptance criteria | prd.md section 29 |
| Development epics | prd.md section 30 |
| Data model | prd.md section 25 |
| User flows | prd.md section 26 |
| Tech stack | prd.md section 5 |
| Non-functional requirements | prd.md section 24 |
| Value proposition | prd.md sections 1, 2, 6 |