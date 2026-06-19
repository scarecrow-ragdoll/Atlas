# acceptance-criteria-reviewer — Worker Report (Attempt 1)

**Run ID**: 20260618T185935Z
**Source**: docs/product/prd.md
**Date**: 2026-06-18
**Scope**: Derive observable acceptance criteria from documented or strongly implied behavior (success, failure, negative, handoff). No new product behavior added.

---

## Methodology

Criteria are grouped by product section. Each entry includes:
- **ID**: `AC-<N>` for derived, `AC-E<M>` for the explicit list in Section 29
- **Source**: PRD section/line reference
- **Rationale**: Why this criterion is derivable
- **Confidence**: High (explicit), Medium (strongly implied), Low (weakly implied)

---

## 1. PIN Guard (Section 7.2, lines 169–184)

| ID | Criterion | Source | Rationale | Confidence |
|----|-----------|--------|-----------|------------|
| AC-01 | User can enable PIN code in settings | §7.2 L177 | Explicit | High |
| AC-02 | User can disable PIN code after it was enabled | §7.2 L182 | Explicit | High |
| AC-03 | User can change PIN code | §7.2 L181 | Explicit | High |
| AC-04 | When PIN is disabled, application opens without any authentication prompt | §7.2 L178 | Explicit | High |
| AC-05 | When PIN is enabled, user must enter correct PIN before accessing the application | §7.2 L179 | Explicit | High |
| AC-06 | PIN code is not stored in plain text | §7.2 L180 | Explicit | High |
| AC-07 | After correct PIN entry, session persists via cookie/session mechanism | §7.2 L183 | Explicit | High |
| AC-08 | Sensitive data (media, personal data) is not accessible without a valid PIN session | §7.2 L184, §24.1 L1073 | Explicit | High |

**Open questions (Q-AC-01)**: What happens when PIN is entered incorrectly? No lockout, rate limit, or retry limit is documented. What is the exact gate — full app redirect vs modal overlay vs specific API 401?

**Open questions (Q-AC-02)**: What is the session TTL? Is there a "remember me" or does it expire on browser close?

---

## 2. Dashboard (Section 9, lines 202–220)

| ID | Criterion | Source | Rationale | Confidence |
|----|-----------|--------|-----------|------------|
| AC-09 | Dashboard shows current date | §9 L208 | Explicit | High |
| AC-10 | Dashboard shows last recorded body weight | §9 L209 | Explicit | High |
| AC-11 | Dashboard shows number of training days for the current week | §9 L210 | Explicit | High |
| AC-12 | Dashboard shows number of cardio entries for the current week | §9 L211 | Explicit | High |
| AC-13 | Dashboard shows current user goal | §9 L212 | Explicit | High |
| AC-14 | Dashboard shows a reminder for upcoming weekly check-in | §9 L213 | Explicit | High |
| AC-15 | Dashboard provides quick action: add workout for today | §9 L216 | Explicit | High |
| AC-16 | Dashboard provides quick action: add cardio | §9 L217 | Explicit | High |
| AC-17 | Dashboard provides quick action: add body weight | §9 L218 | Explicit | High |
| AC-18 | Dashboard provides quick action: open weekly check-in | §9 L219 | Explicit | High |
| AC-19 | Dashboard provides quick action: generate AI report | §9 L220 | Explicit | High |

**Open questions (Q-AC-03)**: How is "last recorded body weight" determined — most recent entry regardless of date, or most recent entry within a time window?

**Open questions (Q-AC-04)**: What defines "current week" — ISO week (Mon-Sun) or Sunday-Saturday? What timezone is used?

**Open questions (Q-AC-05)**: What triggers the "weekly check-in reminder"? Is it date-based (e.g. 7 days since last check-in), configurable, or hardcoded?

---

## 3. Workout Diary — Date Handling (Section 10.2, lines 236–244)

| ID | Criterion | Source | Rationale | Confidence |
|----|-----------|--------|-----------|------------|
| AC-20 | Workout diary opens with today's date by default | §10.2 L240, §6.1 L116 | Explicit | High |
| AC-21 | User can select any past or future date via calendar date picker | §10.2 L241, §6.1 L117 | Explicit | High |
| AC-22 | User can add and edit workouts for past dates (backdating) | §10.2 L242, §6.1 L118 | Explicit | High |
| AC-23 | If a workout exists for a selected date, the existing record is opened | §10.2 L243 | Explicit | High |
| AC-24 | If no workout exists for a selected date, a new record is created only upon first data save | §10.2 L244 | Explicit | High |

---

## 4. Workout Diary — Exercise & Sets (Sections 10.3–10.6, lines 246–306)

| ID | Criterion | Source | Rationale | Confidence |
|----|-----------|--------|-----------|------------|
| AC-25 | Workout day stores date, exercise list, sets per exercise, cardio, day comment, optional body weight | §10.3 L248–255 | Explicit | High |
| AC-26 | Each workout exercise stores reference to exercise library, display order, working weight snapshot, user comment, set list | §10.4 L259–265 | Explicit | High |
| AC-27 | Exercise comment is included in AI export | §10.4 L267 | Explicit | High |
| AC-28 | Each set stores set number, weight, reps, optional RPE, optional RIR, optional comment | §10.5 L279–286 | Explicit | High |
| AC-29 | RPE and RIR are optional; user can log workouts without them using only weight and reps | §10.5 L288–290 | Explicit | High |
| AC-30 | Working weight is stored per exercise in the exercise library | §10.6 L300 | Explicit | High |
| AC-31 | When adding an exercise to a workout day, working weight is auto-populated | §10.6 L301, §6.1 L120 | Explicit | High |
| AC-32 | User can override the auto-populated weight for individual sets | §10.6 L302 | Implicit from "фактический вес подходов хранится отдельно" | Medium |
| AC-33 | Working weight snapshot is preserved at time of workout (not live-linked to library) | §10.6 L304 | Explicit | High |
| AC-34 | AI export includes both working weight and actual set weights | §10.6 L305 | Explicit | High |
| AC-35 | Workout day body weight is optional | §10.3 L255 | Explicit | High |
| AC-36 | Exercise comment field is optional | §10.4 L264 | Implicit from "опционально" pattern and example comments | Medium |

**Open questions (Q-AC-06)**: What is the exact UI for "auto-populate working weight" — is it a pre-filled field the user can change before adding sets? Is it applied per-set or per-exercise as a default?

---

## 5. Progression Tracking (Section 10.7, lines 308–331)

| ID | Criterion | Source | Rationale | Confidence |
|----|-----------|--------|-----------|------------|
| AC-37 | System shows current working weight for an exercise | §10.7 L314 | Explicit | High |
| AC-38 | System shows actual weights for a selected period | §10.7 L315 | Explicit | High |
| AC-39 | System shows best set for an exercise | §10.7 L316 | Explicit | High |
| AC-40 | System shows volume for an exercise | §10.7 L317 | Explicit | High |
| AC-41 | System shows estimated 1RM for an exercise | §10.7 L318 | Explicit | High |
| AC-42 | System shows weekly progress trend for an exercise | §10.7 L319 | Explicit | High |
| AC-43 | System does NOT automatically change working weight without user confirmation | §10.7 L330 | Explicit | High |

**Open questions (Q-AC-07)**: Which formula is used for estimated 1RM (Epley, Brzycki, Lombardi)? Not specified.

**Open questions (Q-AC-08)**: How is "best set" determined — highest weight, highest volume (weight × reps), or highest e1RM?

**Open questions (Q-AC-09)**: What signals are surfaced to the user for progression signals (stable upper rep range, weight increasing, weight plateau, etc.)? The PRD says system "должно отмечать" but is vague — is this a visual indicator on the exercise view, a separate report, or an inline note?

---

## 6. Exercise Library (Section 11, lines 342–387)

| ID | Criterion | Source | Rationale | Confidence |
|----|-----------|--------|-----------|------------|
| AC-44 | User can create an exercise with name, muscle groups, description, personal notes, working weight | §11.2 L355–360 | Explicit | High |
| AC-45 | User can upload images and videos to an exercise | §11.2 L361–362, §11.3 L366 | Explicit | High |
| AC-46 | Media can be added after exercise creation | §11.3 L371 | Explicit | High |
| AC-47 | Media can be deleted | §11.3 L372 | Explicit | High |
| AC-48 | Media is included in full backup export | §11.3 L373 | Explicit | High |
| AC-49 | Media is NOT automatically included in AI export | §11.3 L374 | Explicit | High |
| AC-50 | Exercise has active/inactive status | §11.2 L363 | Explicit | High |
| AC-51 | User does NOT receive a pre-populated starter catalog of exercises in MVP | §11.1 L349 | Explicit | High |

---

## 7. Cardio (Section 12, lines 389–437)

| ID | Criterion | Source | Rationale | Confidence |
|----|-----------|--------|-----------|------------|
| AC-52 | Cardio entry includes date, type, duration (minutes), average pulse, heart rate zone, comment | §12.2 L399–404 | Explicit | High |
| AC-53 | User can select cardio type from the predefined list (walking, running, stationary bike, elliptical, treadmill, other) | §12.3 L408–417 | Explicit | High |
| AC-54 | Cardio entry supports average pulse and heart rate zone | §12.4 L421–426 | Explicit | High |
| AC-55 | Heart rate zone is an enum: Zone 1–5 + unknown | §12.4 L428–436 | Explicit | High |
| AC-56 | If user does not know the zone, they can leave the zone field empty or specify only pulse | §12.4 L437 | Explicit | High |
| AC-57 | Cardio entry pulse field is optional | §12.4 L437 | Implicit from "оставить поле пустым" | Medium |

**Open questions (Q-AC-10)**: Is a cardio entry always linked to a workout day (date-based), or can it exist as a standalone record? The data model shows optional `workoutDayId` but the PRD §10.3 says cardio is part of training day. Clarify: standalone vs day-attached.

---

## 8. Body Measurements & Progress Photos (Sections 13–14, lines 439–537)

| ID | Criterion | Source | Rationale | Confidence |
|----|-----------|--------|-----------|------------|
| AC-58 | User can create a weekly body check-in with date, weight, optional body fat %, 2–4 photos, measurements, and comment | §13.2 L451–458 | Explicit | High |
| AC-59 | User can enter individual body weight entries on any date independent of check-in | §13.5 L496–498 | Explicit | High |
| AC-60 | Body weight entries appear in weight charts and AI export for the selected period | §13.5 L501–502 | Explicit | High |
| AC-61 | Measurements include: neck, shoulders, forearms, biceps, chest, waist, abdomen, hips, thigh, calves | §13.3 L463–474 | Explicit | High |
| AC-62 | For paired measurements (forearms, biceps, thighs, calves), user can specify left and right values | §13.4 L477–490 | Explicit | High |
| AC-63 | If user specifies only one value for a paired measurement, it is treated as the general value | §13.4 L488 | Explicit | High |
| AC-64 | If user specifies left and right, both are stored | §13.4 L489 | Explicit | High |
| AC-65 | The second value for a paired measurement is NOT required if left-only or right-only is provided | §13.4 L490 | Explicit | High |
| AC-66 | Progress photos are linked to a weekly check-in | §14.1 L508 | Explicit | High |
| AC-67 | Each progress photo stores date, check-in reference, file, optional label, optional angle, optional comment | §14.1 L510–517 | Explicit | High |
| AC-68 | Photo angle can be: front, side, back, custom | §14.2 L522–527 | Explicit | High |
| AC-69 | User can view photos at any time | §14.3 L534 | Explicit | High |
| AC-70 | Photos are NOT publicly accessible without PIN session / authorization | §14.3 L535, §24.1 L1073 | Explicit | High |
| AC-71 | Photos are included in full backup export | §14.3 L536 | Explicit | High |
| AC-72 | Photos are included in AI export only when user explicitly opts in | §14.3 L537, §17.3 L712 | Explicit | High |

---

## 9. Nutrition (Section 15, lines 539–618)

| ID | Criterion | Source | Rationale | Confidence |
|----|-----------|--------|-----------|------------|
| AC-73 | User can create nutrition products with name, calories/100g, protein/100g, fat/100g, carbs/100g, optional notes | §15.2 L555–562 | Explicit | High |
| AC-74 | User can create one weekly nutrition template | §15.3 L566 | Explicit | High |
| AC-75 | Template includes product list with grams per day, optional meal label, optional comment | §15.3 L568–575 | Explicit | High |
| AC-76 | System calculates calories, protein, fat, carbs from template products | §15.3 L577–582 | Explicit | High |
| AC-77 | By default, every day of the week uses the weekly template | §15.4 L586 | Explicit | High |
| AC-78 | User can override nutrition for a specific day (daily override) | §15.5 L596–601 | Explicit | High |
| AC-79 | Daily override supports: add product, remove product, change product amount, add comment | §15.5 L597–601 | Explicit | High |
| AC-80 | Daily override affects ONLY the selected date | §15.5 L603 | Explicit | High |
| AC-81 | When viewing a specific day, user sees template-calculated nutrition with overrides applied | §26.8 L1430 | Explicit | High |
| AC-82 | System recalculates macros when override is applied | §26.8 L1432 | Explicit | High |

**Open questions (Q-AC-11)**: Can user have multiple weekly templates? §15.3 says "один шаблон" — is replacement the only option (creating a new one replaces the old)?

**Open questions (Q-AC-12)**: What happens when a product used in a template/override is deleted? Is the reference orphaned, or is deletion blocked?

---

## 10. Charts (Section 16, lines 620–671)

| ID | Criterion | Source | Rationale | Confidence |
|----|-----------|--------|-----------|------------|
| AC-83 | User can view charts for a selected time period | §16.1 L624, L628 | Explicit | High |
| AC-84 | Charts can be filtered | §16.1 L629 | Explicit | High |
| AC-85 | Charts use data from workouts, measurements, and nutrition | §16.1 L631 | Explicit | High |
| AC-86 | For a selected exercise, charts show: working weight, best set, e1RM, total volume, total reps, number of working sets | §16.2 L637–642 | Explicit | High |
| AC-87 | User selects specific exercise and period for exercise charts | §16.2 L644 | Explicit | High |
| AC-88 | Body charts show: body weight, body fat %, each individual measurement, multiple selected measurements on one chart | §16.3 L649–653 | Explicit | High |
| AC-89 | Nutrition charts show: average calories per week, average protein per week, average fat per week, average carbs per week | §16.4 L668–671 | Explicit | High |

**Open questions (Q-AC-13)**: What does "фильтровать" mean? By date range only, or by exercise/category, or both?

---

## 11. AI Export (Section 17, lines 673–792)

| ID | Criterion | Source | Rationale | Confidence |
|----|-----------|--------|-----------|------------|
| AC-90 | User can configure export period; default is last 4 weeks | §17.2 L690–691 | Explicit | High |
| AC-91 | User can select custom start and end dates | §17.2 L692 | Explicit | High |
| AC-92 | User can choose which data sections to include in export | §17.3 L696–710 | Explicit | High |
| AC-93 | Photos are excluded from export by default | §17.3 L712 | Explicit | High |
| AC-94 | Export is delivered as a ZIP archive | §17.4 L716 | Explicit | High |
| AC-95 | ZIP filename follows pattern `atlas-ai-export-YYYY-MM-DD.zip` | §17.4 L720 | Explicit | High |
| AC-96 | ZIP contains: manifest.json, data.json, summary.md, workouts.csv, measurements.csv, nutrition.csv, cardio.csv, and optionally photos/ directory | §17.4 L721–730 | Explicit | High |
| AC-97 | manifest.json includes: export type, schema version, app version, export date, selected period, included sections, file list, photo presence info | §17.5 L736–743 | Explicit | High |
| AC-98 | data.json includes: profile, goals, exercises, workouts, workoutExercises, sets, cardio, bodyWeightEntries, bodyCheckIns, measurements, nutritionProducts, nutritionTemplates, dailyNutritionOverrides, userComments, computedSummary | §17.6 L751–765 | Explicit | High |
| AC-99 | summary.md includes: report period, current goal, training stats, exercise dynamics, weight changes, measurement changes, nutrition summary, cardio, user comments | §17.7 L771–781 | Explicit | High |
| AC-100 | CSV files exist for: workouts, measurements, nutrition, cardio | §17.8 L787–792 | Explicit | High |

---

## 12. AI Prompt Builder (Section 18, lines 794–873)

| ID | Criterion | Source | Rationale | Confidence |
|----|-----------|--------|-----------|------------|
| AC-101 | System stores persistent AI context: current goal, height, optional age, training experience, optional training split, injuries/limitations, preferred progression style, current nutrition strategy, additional persistent comment | §18.2 L804–814 | Explicit | High |
| AC-102 | User can change goal and context at any time | §18.2 L817 | Explicit | High |
| AC-103 | If goal was entered once, it is remembered for subsequent exports | §18.2 L818 | Explicit | High |
| AC-104 | During AI export generation, user can add a one-time comment | §18.3 L823 | Explicit | High |
| AC-105 | One-time comment appears ONLY in the current AI export | §18.3 L836 | Explicit | High |
| AC-106 | System supports week flags: poor sleep, high stress, illness, injury/pain, AAS/cycle, calorie deficit, calorie surplus, maintenance, missed workouts, travel/disrupted schedule | §18.4 L841–852 | Explicit | High |
| AC-107 | Week flags are included in AI prompt and data.json | §18.4 L854 | Explicit | High |
| AC-108 | Generated prompt asks AI to: analyze exercise progress, evaluate working weight dynamics, compare actual vs working weights, evaluate volume, consider RPE/RIR and comments, consider cardio, correlate training with weight/measurement changes, consider nutrition, give weekly recommendations, suggest working weight changes, flag exercises for increase/repeat/deload, give actionable plan | §18.5 L858–873 | Explicit | High |

**Open questions (Q-AC-14)**: Are week flags per-user per-week or per-export? Is there a single set of flags per week that applies to all exports in that week?

---

## 13. AI Review History (Section 19, lines 875–898)

| ID | Criterion | Source | Rationale | Confidence |
|----|-----------|--------|-----------|------------|
| AC-109 | User can create an AI review record | §19.2 L892 | Explicit | High |
| AC-110 | User can paste AI response text into the review | §19.2 L893 | Explicit | High |
| AC-111 | Review is linked to a date period | §19.2 L894 | Explicit | High |
| AC-112 | User can add notes to the review | §19.2 L895 | Explicit | High |
| AC-113 | User can mark planned actions from the review | §19.2 L896 | Explicit | High |
| AC-114 | No automatic OpenAI/ChatGPT API integration in MVP | §19.2 L898 | Explicit | High |

---

## 14. Import / Export — Full Backup (Section 20, lines 900–1000)

| ID | Criterion | Source | Rationale | Confidence |
|----|-----------|--------|-----------|------------|
| AC-115 | User can generate a full backup export of all data | §20.1 L904 | Explicit | High |
| AC-116 | Backup export is a ZIP archive | §20.1 L913 | Explicit | High |
| AC-117 | ZIP filename follows pattern `atlas-backup-YYYY-MM-DD.zip` | §20.1 L917 | Explicit | High |
| AC-118 | Backup ZIP contains manifest.json, data.json, media/ directory | §20.1 L918–922 | Explicit | High |
| AC-119 | manifest.json includes: export type (`full_backup`), schemaVersion, appVersion, exportedAt, includedSections, mediaIncluded, file list, optional checksums | §20.2 L938–946 | Explicit | High |
| AC-120 | data.json includes all entities: settings, profile, goals, exercises, exercise media, workout days, workout exercises, sets, cardio, body weight, check-ins, measurements, photos metadata, nutrition products, nutrition templates, daily overrides, AI prompt settings, AI review history | §20.3 L952–969 | Explicit | High |
| AC-121 | User can import a backup ZIP to restore data in a new instance | §20.4 L973 | Explicit | High |
| AC-122 | Import supports ZIP file upload | §20.4 L977 | Explicit | High |
| AC-123 | Import validates manifest.json before restoring | §20.4 L978 | Explicit | High |
| AC-124 | Import validates schema version | §20.4 L979 | Explicit | High |
| AC-125 | Import validates data.json structure | §20.4 L980 | Explicit | High |
| AC-126 | Import performs dry-run validation before actual restore | §20.4 L981 | Explicit | High |
| AC-127 | Import displays summary to user before confirming | §20.4 L982 | Explicit | High |
| AC-128 | Import shows clear error messages on validation failure | §20.4 L983 | Explicit | High |
| AC-129 | Import prohibits silent partial import (all-or-nothing) | §20.4 L984 | Explicit | High |
| AC-130 | Import restores media files | §20.4 L985 | Explicit | High |
| AC-131 | Import restores entity relationships | §20.4 L986 | Explicit | High |
| AC-132 | Export format includes schema version | §20.5 L990 | Explicit | High |

**Open questions (Q-AC-15)**: What happens on import when data already exists in the target instance? Merge, replace, or error? Not specified.

**Open questions (Q-AC-16)**: CSV files in backup — are they mandatory or optional? §20.1 says "опционально можно добавить" — this is optional, not required.

---

## 15. Routine Optimization (Section 21, lines 1002–1025)

| ID | Criterion | Source | Rationale | Confidence |
|----|-----------|--------|-----------|------------|
| AC-133 | Fields are pre-filled where possible | §6.1 L119, §21 L1011 | Explicit | Medium |
| AC-134 | User can copy previous set's weight/reps within an exercise | §21 L1013 | Explicit | High |
| AC-135 | System calculates tonnage (volume) | §21 L1014 | Explicit | High |
| AC-136 | System calculates estimated 1RM | §21 L1015 | Explicit | High |
| AC-137 | System generates AI prompt text | §21 L1020 | Explicit | High |
| AC-138 | System generates ZIP export | §21 L1021 | Explicit | High |

**Open questions (Q-AC-17)**: "Копирование значения прошлого подхода" — is this a one-tap duplication of weight/reps from the previous set row, or auto-fill as user types? Not enough detail.

---

## 16. Non-Functional Requirements (Section 24, lines 1061–1119)

| ID | Criterion | Source | Rationale | Confidence |
|----|-----------|--------|-----------|------------|
| AC-139 | PIN code is not written to logs | §24.1 L1069 | Explicit | High |
| AC-140 | AI export content is not written to logs | §24.1 L1070 | Explicit | High |
| AC-141 | Photos are not written to logs | §24.1 L1071 | Explicit | High |
| AC-142 | Sensitive comments are not written to logs | §24.1 L1072 | Explicit | High |
| AC-143 | Media files are not served without valid PIN session | §24.1 L1073 | Explicit | High |
| AC-144 | Full backup and AI export are generated only on user request | §24.1 L1074 | Explicit | High |
| AC-145 | Data is not lost on container restart | §24.2 L1080 | Explicit | High |
| AC-146 | Media files are stored in a Docker volume | §24.2 L1081 | Explicit | High |
| AC-147 | Backup export includes both data and media | §24.2 L1082 | Explicit | High |
| AC-148 | Import validates archive integrity | §24.2 L1083 | Explicit | High |
| AC-149 | Interface remains responsive at personal data volume (years of workouts, thousands of sets, hundreds of photos/products) | §24.3 L1089–1096 | Explicit | Medium |

**Open questions (Q-AC-18)**: What defines "чувствительные комментарии"? Is this a labeled field (e.g. `isSensitive`) or does it refer to exercise/check-in comment fields in general?

---

## 17. Explicit Acceptance Criteria (Section 29, lines 1531–1560) — Coverage Mapping

These 26 criteria are explicitly listed. Map to derived criteria:

| # | Explicit AC | Derived Coverage | Status |
|---|-------------|------------------|--------|
| 1 | Enable/disable PIN | AC-01, AC-02 | Covered |
| 2 | Create exercise | AC-44 | Covered |
| 3 | Set working weight | AC-44 | Covered |
| 4 | Upload media to exercise | AC-45, AC-46 | Covered |
| 5 | Open current day in diary | AC-20 | Covered |
| 6 | Select past date via calendar | AC-21 | Covered |
| 7 | Add exercise to workout day | AC-26 | Covered |
| 8 | Add sets with weight and reps | AC-28 | Covered |
| 9 | Optional RPE/RIR | AC-29 | Covered |
| 10 | Add exercise comment | AC-36 | Covered |
| 11 | Add cardio | AC-52, AC-53 | Covered |
| 12 | Create weekly body check-in | AC-58 | Covered |
| 13 | Add weight, body fat %, measurements, photos | AC-58, AC-61, AC-66 | Covered |
| 14 | Add separate weight entry | AC-59 | Covered |
| 15 | Create product | AC-73 | Covered |
| 16 | Create weekly nutrition template | AC-74 | Covered |
| 17 | Change nutrition for a specific day | AC-78, AC-79 | Covered |
| 18 | View exercise progress chart | AC-86 | Covered |
| 19 | View weight/measurement chart | AC-88 | Covered |
| 20 | View basic KBJU chart | AC-89 | Covered |
| 21 | Generate AI prompt | AC-108, AC-137 | Covered |
| 22 | Download AI export ZIP (last 4 weeks) | AC-90, AC-94 | Covered |
| 23 | Save AI review | AC-109 | Covered |
| 24 | Download full backup | AC-115 | Covered |
| 25 | Import backup into clean instance | AC-121 | Covered |
| 26 | Verification command passes tests and coverage gate | AC-(not derived — infra/test NFR) | Not derived as it's CI/test infra |

---

## 18. Explicit "Not in MVP" Constraints (Sections 10.8, 15.6, 22.1, 23, 28)

These are not criteria but constraints. Key ones for AC boundary:

- Workout templates not in MVP (§10.8 L336)
- Workout planning not in MVP (§10.8 L337)
- Quick repeat of last workout not in MVP (§10.8 L338, §21 L1025)
- Training programs not in MVP (§10.8 L339)
- Split training not in MVP (§10.8 L340)
- Future workout calendar not in MVP (§10.8 L341)
- Recipes and meals not in MVP (§15.6 L609–610)
- Barcode scanner not in MVP (§15.6 L611)
- Water, fiber, salt, sugar, alcohol tracking not in MVP (§15.6 L612–616)
- Food recognition not in MVP (§15.6 L617)
- Public food database not in MVP (§15.6 L618)
- Apple Health not in MVP (§22.1 L1031)
- Telegram bot not in MVP (§23 L1051)
- Registration, multi-user, roles, SaaS, public pages not in MVP (§28 L1513–1517)
- OpenAI API integration not in MVP (§28 L1527, §19.2 L898)

---

## 19. Summary

| Metric | Count |
|--------|-------|
| Total derived criteria (AC-*) | 139 |
| Explicit criteria (AC-E*) | 26 (mapped coverage) |
| Open questions (Q-AC-*) | 18 |
| High confidence criteria | 134 |
| Medium confidence criteria | 5 |
| Low confidence criteria | 0 |

### Confidence Assessment

The PRD is well-structured with explicit requirements in most sections. Sections 7–20 and 24–29 provide enough detail for high-confidence AC derivation. The explicit acceptance criteria list (§29) covers the main user journeys but misses edge cases and failure modes (e.g., what happens when PIN is wrong, what happens on empty data states, what happens when media upload fails).

### Key Gaps
1. PIN entry failure behavior (retry limit, lockout, error display) — not specified
2. Session TTL for PIN-authenticated sessions
3. Estimated 1RM formula selection
4. Best set / progression signal definition
5. Nutrition template lifecycle when products are deleted
6. Import behavior when data already exists (merge vs replace vs reject)
7. What constitutes "sensitive comments" for logging exclusion