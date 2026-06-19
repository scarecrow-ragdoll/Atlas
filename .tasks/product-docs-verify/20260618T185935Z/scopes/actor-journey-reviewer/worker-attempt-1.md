# Actor-Journey Review — Worker Report (Attempt 1)

**Run ID:** 20260618T185935Z  
**Source:** docs/product/prd.md  
**Scope focus:** Personas, user journeys, happy paths, alternative paths, empty states, recovery paths  
**Reviewer:** actor-journey-reviewer

---

## 1. Personas Coverage

| Persona | Defined? | Assessment |
|---|---|---|
| **Single user (athlete/tracker)** | Partially | Sections 7.1 and implicit throughout. The user is described as someone who works out, tracks nutrition, and reviews progress. No persona depth: goals, habits, tech comfort, frequency of use, or motivation profile. |
| **Self-hoster / deployer** | Not as a persona | The app is self-hosted and single-user (§4), but the deployer's journey (install, first setup, PIN config) is not described as a distinct flow. |
| **AI reviewer (ChatGPT)** | Not a persona | The AI is a consumer of exported data (§6.2, §17). Its interaction model (consumes prompt + data, returns text) is outlined but not formalized as an actor. |

**Finding:** Only one implicit persona exists. The deployer and AI-as-consumer are present in requirements but missing as explicit actors with journeys.

---

## 2. Happy Paths (Well-Defined)

The following scenarios in §26 are well-structured happy paths with clear step sequences:

- Add exercise (§26.1)
- Log today's workout (§26.2)
- Log backdated workout (§26.3)
- Add cardio (§26.4)
- Weekly check-in (§26.5)
- Add weight separately (§26.6)
- Create weekly nutrition template (§26.7)
- Override daily nutrition (§26.8)
- Generate AI report (§26.9)
- Save AI review (§26.10)
- Full backup export (§26.11)
- Restore from backup (§26.12)

These 12 journeys cover the core MVP flows listed in §27 and §29.

---

## 3. Alternative Paths (Gaps)

| Scenario | Alternative Path | Status | Question ID |
|---|---|---|---|
| Open workout diary for a date with existing data | Record opens for editing (§26.3, step 3) | **Covered** | — |
| Add exercise to workout when exercise doesn't exist in library | Must exit workout entry, go to exercise library, create exercise, return | **Missing** | Q-ACTOR-01 |
| Add a set with invalid values (0 reps, 0 weight) | Validation error | **Missing** | Q-ACTOR-02 |
| Weekly check-in without photos | User doesn't have 2-4 photos ready | **Missing** | Q-ACTOR-03 |
| Remove the last exercise from a workout day | Day becomes empty | **Missing** | Q-ACTOR-04 |
| Add the same exercise twice in one workout day | Duplicate entry or two entries? | **Missing** | Q-ACTOR-05 |
| Nutrition override reverting to template | User wants to "remove all overrides for today" | **Missing** | Q-ACTOR-06 |
| AI report generation with no data in selected period | Empty export | **Missing** | Q-ACTOR-07 |
| Backup import when data already exists in instance | Overwrite, merge, or reject? | **Missing** | Q-ACTOR-08 |
| User opens app with PIN enabled but no valid session | PIN prompt appears | **Covered** (§7.2) | — |
| User enters wrong PIN | Failed attempt feedback, lockout? | **Missing** | Q-ACTOR-09 |

---

## 4. Empty States (Major Gap)

No section in the PRD describes what the UI should show when each section has no data. This affects:

| Section | Empty State | Status | Question ID |
|---|---|---|---|
| Dashboard | First launch, zero data (no weight, no workouts, no goal) | **Missing** | Q-ACTOR-10 |
| Exercise library | No exercises created yet | **Missing** | Q-ACTOR-11 |
| Workout diary (first visit) | No days with data | **Missing** | Q-ACTOR-12 |
| Body measurements | No check-ins yet | **Missing** | Q-ACTOR-13 |
| Progress photos | No photos uploaded | **Missing** | Q-ACTOR-14 |
| Nutrition products | No products created | **Missing** | Q-ACTOR-15 |
| Nutrition template | No template for current week | **Missing** | Q-ACTOR-16 |
| AI export history | No exports yet | **Missing** | Q-ACTOR-17 |
| AI review history | No reviews saved | **Missing** | Q-ACTOR-18 |
| Charts (any section) | No data for selected period | **Missing** | Q-ACTOR-19 |
| Backup export | No backups yet | **Missing** | Q-ACTOR-20 |

---

## 5. Recovery / Error Paths (Major Gap)

| Scenario | Error/Recovery | Status | Question ID |
|---|---|---|---|
| PIN lost/forgotten | User cannot access app. No recovery mechanism (no email, no account). | **Missing** | Q-ACTOR-21 |
| ZIP export generation fails | Disk full, permission error | **Missing** | Q-ACTOR-22 |
| Backup import validation passes but actual restore fails | Partial data restored? Transaction rollback? | **Missing** | Q-ACTOR-23 |
| Media upload fails | File too large, wrong format, network error | **Missing** | Q-ACTOR-24 |
| Database connection lost mid-entry | Data not saved, user retries | **Missing** | Q-ACTOR-25 |
| User navigates away during workout entry | Unsaved data lost? Autosave? | **Missing** | Q-ACTOR-26 |
| Invalid backup ZIP uploaded | Clear error message | **Missing** | Q-ACTOR-27 |
| Session expires during long data entry | Data loss risk | **Missing** | Q-ACTOR-28 |
| Concurrent access (same day opened in two tabs) | Last write wins? Conflict? | **Missing** | Q-ACTOR-29 |
| User creates nutrition template then changes goal | Template no longer fits, but no recalculation is triggered | **Missing** | Q-ACTOR-30 |

---

## 6. First-Run / Onboarding Journey (Gap)

No first-run experience is described:

| Aspect | Finding | Question ID |
|---|---|---|
| App setup on first launch | User opens app → no data → no guidance | Q-ACTOR-31 |
| Goal/context setup | User must open Settings → UserProfile manually. No prompt to set goal before first AI export. | Q-ACTOR-32 |
| PIN setup during first run | PIN is optional but no nudge shown at first launch if user wants security. | Q-ACTOR-33 |

---

## 7. Edge Cases (Miscellaneous)

| Finding | Question ID |
|---|---|
| Exercise search/filter for large exercise library (>50 exercises) is undefined | Q-ACTOR-34 |
| Cardio entry without pulse/zone data is allowed but behavior not specified (e.g., "unknown" zone default) | Q-ACTOR-35 |
| Multiple workout days: can user view a list of "all days with workouts" or only navigate by calendar? | Q-ACTOR-36 |
| Nutrition template expires: what happens when template week (e.g., week 1) passes and no new template for week 2 is created? | Q-ACTOR-37 |
| Body weight entry duplicates: what if user enters weight twice on same date? Overwrite, reject, keep latest? | Q-ACTOR-38 |
| Exercise working weight changes after it is used in past workouts: should past snapshots remain or update? | **Covered** (§10.6, step 4) — snapshot is stored per workout day. |

---

## 8. Summary

| Metric | Value |
|---|---|
| Happy paths fully specified | 12 |
| Alternative / edge paths specified | 1 of ~15 |
| Empty states specified | 0 of ~11 |
| Recovery / error paths specified | 0 of ~10 |
| Personas with explicit journeys | 0 (1 implicit) |
| Open questions raised | 38 (Q-ACTOR-01 to Q-ACTOR-38) |

**Overall assessment:** The PRD defines clear happy-path scenarios for the core workflows, which is a strong foundation. However, empty states, recovery paths, error handling, first-run experience, and most alternative paths are absent. These gaps will surface as UX ambiguity during development and testing. The PIN flow (§7.2) is the only area with basic alternative and recovery coverage.