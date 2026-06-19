# Product Scope Reviewer Worker Attempt 2

## Sources Read

- docs/product/prd.md (1665 lines, complete)
- .tasks/product-docs-verify/20260618T185935Z/scopes/product-scope-reviewer/worker-attempt-1.md
- .tasks/product-docs-verify/20260618T185935Z/scopes/product-scope-reviewer/review-attempt-1.md

## Source Delta Reviewed

No source delta present. This is a revision incorporating reviewer findings from attempt 1.

## Changes from Attempt 1

- Removed contradiction #4 (working weight snapshot) — correctly described intended behavior, not a contradiction
- Added explicit handoff readiness assessment
- Added PIN-vs-registration/auth tension
- Added Q-SCOPE-006, Q-SCOPE-007, Q-SCOPE-008
- Separated "Contradictions" from "Ambiguities/Tensions"

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

## Contradictions (Strict)

1. **"No registration" vs PIN auth system**: Section 4 states "без регистрации" and section 7.1 reinforces no registration, no roles, no user management. However, section 7.2 introduces a full PIN authentication system with hashed storage, session cookies, and protected routes. This is effectively an auth system — the PRD describes it as "not registration" but the PIN gate controls data access the same way a password would.

2. **Cardio entity ambiguity**: Section 10.3 lists cardio as part of the workout day ("кардио за этот день"). The data model (section 25.8, CardioEntry) has an optional `workoutDayId`, making cardio a separate entity that can optionally belong to a workout day. These are two different models — cardio as a field of WorkoutDay vs cardio as an independent entity with optional relationship. This needs resolution.

3. **Telegram bot in tech stack**: `go-telegram/bot` is listed in section 5 (Technology Stack) but section 23 explicitly states Telegram bot is not in MVP. Including the library in the stack creates a build-time dependency that may not be needed for MVP.

## Ambiguities / Tensions

1. **Single-user vs future multi-user architecture**: Section 4 says "один пользователь" and "без multi-user режима в MVP", but also says "в будущем архитектура может быть расширена под multi-user". It is unclear whether the MVP should build with multi-user considerations (table prefixes, tenant isolation patterns, etc.) or remain purely single-user with a future migration planned.

2. **Working weight snapshot design**: Section 10.6 describes snapshot behavior. This is correctly designed (snapshot at time of workout, independent from exercise library changes) and not a contradiction. It is an implicit guarantee that past workout data remains stable when the exercise library is updated.

## Missing Source Artifacts

- **No success metrics or KPIs** defined anywhere. This is the most significant gap for a scope review.
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
6. **No minimum AI compatibility specification** — The export format mentions ChatGPT specifically but is expected to work with other AI platforms.

## Open Questions Raised

| ID | Question | Why It Matters |
|----|----------|----------------|
| Q-SCOPE-001 | What are the quantitative success metrics for the MVP? | Without success metrics, there is no way to validate the product is achieving its goals |
| Q-SCOPE-002 | Should the MVP architecture build for future multi-user or remain strictly single-user? | Affects data model, auth design, API structure, and database schema decisions |
| Q-SCOPE-003 | What is the target user's expected technical proficiency? | Self-hosted requires Docker/CLI knowledge — affects documentation and deployment UX requirements |
| Q-SCOPE-004 | What are the specific performance targets (page load, export time, chart render time)? | Section 24.3 uses vague language ("быстрым") with no measurable targets |
| Q-SCOPE-005 | Should cardio be a separate entity or always part of a workout day? | The data model shows cardio as separate with optional workoutDayId, but section 10.3 includes cardio in the workout day — ambiguous |
| Q-SCOPE-006 | What AI models/platforms must the export format support? Only ChatGPT, or Claude, Gemini, local LLMs? | The core value proposition depends on AI compatibility; format design affects all models differently |
| Q-SCOPE-007 | Is there a maximum photo/media storage limit? What happens when storage volume runs out? | Photos over years of training accumulate; no storage management policy exists |
| Q-SCOPE-008 | What data portability standard is required for "data belongs to user" commitment? | Export/backup formats need an interoperability guarantee beyond the current ZIP structure |

## Edge Cases Or Risks

1. **Scope creep risk**: The PRD extensively discusses future features (multi-user, Apple Health, Telegram bot, mobile app, coach mode, cloud backup). While section 28 explicitly excludes them from MVP, the volume of future scope discussion could lead to architectural over-engineering unless disciplined.

2. **Single point of failure**: Self-hosted single-user means all data is lost if the user's Docker volume or database is corrupted without a backup. The backup/restore feature mitigates this but depends on the user proactively using it.

3. **AI export dependency risk**: The core value proposition depends on the user having access to ChatGPT or another AI. If AI access changes (paywall, availability, API changes), the primary workflow is broken.

4. **Photo storage growth**: Hundreds of photos over years of training could lead to significant storage requirements. No guidance on storage limits, compression, or cleanup.

5. **Nutrition template complexity**: The weekly template + daily override model is elegant for regular eaters but may be too restrictive for users with highly variable daily diets.

## Handoff Readiness Assessment

**This PRD is NOT ready for development handoff in its current state.**

Required before handoff:
1. Success metrics must be defined (see Q-SCOPE-001).
2. The single-user vs multi-user architectural decision must be resolved (Q-SCOPE-002).
3. The cardio entity relationship must be clarified (Q-SCOPE-005).
4. Performance targets must be specified (Q-SCOPE-004).
5. AI platform compatibility must be defined (Q-SCOPE-006).

The PRD is strong in scope definition, user flows, data model draft, acceptance criteria, and non-functional requirements. These four items are the only blockers to handoff readiness from a scope perspective. Other missing artifacts (UX specs, deployment guide, glossary) are important but can be addressed during development.

## Recommended Decisions

1. **Define MVP success metrics before development**: At minimum: (a) user can complete the full weekly cycle without data loss, (b) AI export generates valid ZIP in under X seconds for 4 weeks of data, (c) user can enter a workout in under Y minutes.

2. **Decide single-user architectural boundary**: Document whether future multi-user compatibility is a design concern during MVP or explicitly deferred.

3. **Define performance targets**: At minimum: page load < 2s, AI export generation < 5s for 4 weeks of data, chart render < 1s.

4. **Accept the PIN system as de facto auth**: Rename "no registration" to "no user management" to avoid the contradiction with the PIN auth implementation.

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