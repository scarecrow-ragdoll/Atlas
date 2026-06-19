# Actor-Journey Review — Worker Report (Attempt 2)

**Run ID:** 20260618T185935Z
**Attempt:** 2
**Source:** docs/product/prd.md
**Scope focus:** Personas, user journeys, happy paths, alternative paths, empty states, recovery paths
**Reviewer:** actor-journey-reviewer
**Previous revision:** Attempt 1 received "needs-revision" — empty states compressed, cross-scope tags added, duplicates consolidated, positive findings section added.

---

## 0. Positive Finding: Happy-Path Map Consistency

The PRD defines **12 well-structured happy-path scenarios** in §26 (26.1–26.12). These align cleanly with the 25 acceptance criteria in §29:

| Scenario (§26) | Covered ACs (§29) |
|---|---|
| 26.1 Add exercise | AC 2, 3 |
| 26.2 Log today's workout | AC 5, 7, 8, 9, 10 |
| 26.3 Log backdated workout | AC 6, 7, 8, 9, 10 |
| 26.4 Add cardio | AC 11 |
| 26.5 Weekly check-in | AC 12, 13 |
| 26.6 Add weight separately | AC 14 |
| 26.7 Create nutrition template | AC 15, 16 |
| 26.8 Override daily nutrition | AC 17 |
| 26.9 Generate AI report | AC 21, 22 |
| 26.10 Save AI review | AC 23 |
| 26.11 Full backup export | AC 24 |
| 26.12 Restore from backup | AC 25 |

AC 1 (PIN enable/disable) has no dedicated scenario but is described in §7.2. AC 26 (test/coverage gate) is a build-time requirement.

The scenario list also respects the **Out of Scope boundary** in §28 — no workout templates, barcode scanner, Apple Health, etc. are introduced in the happy paths.

This gives the development team a reliable "happy-path map" with clear traceability to acceptance criteria. **This is a strong foundation and should be preserved as-is.**

---

## 1. Personas Coverage

| Persona | Defined? | Assessment |
|---|---|---|
| **Single user (athlete/tracker)** | Partially | Implicit throughout but no persona depth (goals, habits, tech comfort, frequency, motivation). §7.1 states "one user" but does not describe who they are. |
| **Self-hoster / deployer** | Not as a persona | §4 mentions self-hosted, single-user deployment. The deployer's journey (install, first setup, PIN config) is absent as a distinct flow. |
| **AI reviewer (ChatGPT)** | Not a persona | §6.2, §17 describe AI as consumer of exported data (prompt + file → text). Interaction model is outlined but not formalized as an actor. |

**Finding:** Only one implicit persona exists. Deployer and AI-as-consumer are present in requirements but missing as explicit actors with journeys.

---

## 2. Happy Paths (Well-Defined, §26)

12 scenarios fully specified — see §0 above for consistency mapping with §29 and §28.

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
| AI report or chart with no data in selected period | Empty export / empty chart | **Missing** | Q-ACTOR-07 |
| Backup import when data already exists in instance | Overwrite, merge, or reject? | **Missing** | Q-ACTOR-08 |
| User opens app with PIN enabled but no valid session | PIN prompt appears | **Covered** (§7.2) | — |
| User enters wrong PIN | Failed attempt feedback, lockout? | **Missing** | Q-ACTOR-09 ★ |

★ = cross-scope: roles/permissions

---

## 4. Empty States

No section in the PRD describes what the UI should show when each section has no data. After revision (compressed from per-section enumeration to a first-run convention plus behavior-critical states):

| Section | Empty State | Status | Question ID |
|---|---|---|---|
| **Behavior-critical empty states** | | | |
| Dashboard | First launch, zero data (no weight, no workouts, no goal) | **Missing** | Q-ACTOR-10 |
| Nutrition template expires | Template week (e.g., week 1) passes, no new template for week 2. What does the daily view show? | **Missing** | Q-ACTOR-11 |
| **First-run empty-state convention** | Consolidated — applies to: exercise library, workout diary (first visit), body measurements, progress photos, nutrition products, AI export history, AI review history, backup history, charts section | **Missing** | Q-ACTOR-12 |

Q-ACTOR-12 represents a single design decision: what is the convention for all first-visit/empty sections? (e.g., "no data yet — here's how to start" vs. blank slate vs. guided setup)

Retained behavior-critical states (Q-ACTOR-10, Q-ACTOR-11) affect core UX decisions (dashboard as landing page, nutrition daily view).

---

## 5. Recovery / Error Paths (Major Gap)

| Scenario | Error/Recovery | Status | Question ID |
|---|---|---|---|
| PIN lost/forgotten | User cannot access app. No recovery mechanism (no email, no account). | **Missing** | Q-ACTOR-13 ★ |
| ZIP export generation fails | Disk full, permission error | **Missing** | Q-ACTOR-14 |
| Backup import validation passes but actual restore fails | Partial data restored? Transaction rollback? | **Missing** | Q-ACTOR-15 |
| Media upload fails | File too large, wrong format, network error | **Missing** | Q-ACTOR-16 |
| Mid-session resilience | Database connection lost or session expires during long data entry. Data loss risk, retry behavior. | **Missing** | Q-ACTOR-17 ★ |
| User navigates away during workout entry | Unsaved data lost? Autosave? | **Missing** | Q-ACTOR-18 |
| Invalid backup ZIP uploaded | Clear error message | **Missing** | Q-ACTOR-19 |
| Concurrent access (same day opened in two tabs) | Last write wins? Conflict? | **Missing** | Q-ACTOR-20 |
| User creates nutrition template then changes goal | Template no longer fits, but no recalculation is triggered | **Missing** | Q-ACTOR-21 |

★ = cross-scope: roles/permissions (Q-ACTOR-13: PIN recovery is auth policy; Q-ACTOR-17: session expiry touches auth)

**Deduplication note:** Q-ACTOR-07 (empty data results) consolidates the previously separate Q-ACTOR-07 (empty AI export) and Q-ACTOR-19 (charts with no data). Q-ACTOR-17 (mid-session resilience) consolidates the previously separate Q-ACTOR-25 (DB connection lost) and Q-ACTOR-28 (session expiry).

---

## 6. First-Run / Onboarding Journey (Gap)

No first-run experience is described:

| Aspect | Finding | Question ID |
|---|---|---|
| App setup on first launch | User opens app → no data → no guidance | Q-ACTOR-22 |
| Goal/context setup | User must open Settings → UserProfile manually. No prompt to set goal before first AI export. | Q-ACTOR-23 |
| PIN setup during first run | PIN is optional but no nudge shown at first launch if user wants security. | Q-ACTOR-24 |

---

## 7. Edge Cases (Miscellaneous)

| Finding | Question ID |
|---|---|
| Exercise search/filter for large exercise library (>50 exercises) is undefined | Q-ACTOR-25 |
| Cardio entry without pulse/zone data is allowed but behavior not specified (e.g., "unknown" zone default) | Q-ACTOR-26 |
| Multiple workout days: can user view a list of "all days with workouts" or only navigate by calendar? | Q-ACTOR-27 |
| Body weight entry duplicates: what if user enters weight twice on same date? Overwrite, reject, keep latest? | Q-ACTOR-28 |
| Exercise working weight changes after it is used in past workouts: should past snapshots remain or update? | **Covered** (§10.6, step 4) — snapshot is stored per workout day. |

---

## 8. Summary

| Metric | Value |
|---|---|
| Happy paths fully specified | 12 (traceable to §29 ACs, respecting §28 boundaries) |
| Alternative / edge paths specified | 1 of ~15 |
| Empty states specified | 0 (1 convention + 2 behavior-critical = 3 questions) |
| Recovery / error paths specified | 0 of ~9 |
| Personas with explicit journeys | 0 (1 implicit) |
| Open questions raised | 28 (Q-ACTOR-01 to Q-ACTOR-28) — reduced from 38 via compression and deduplication |
| Cross-scope references (roles/permissions) | Q-ACTOR-09, Q-ACTOR-13, Q-ACTOR-17 |

**Overall assessment:** Strong happy-path foundation with clear §29 traceability. Empty states, recovery paths, error handling, first-run experience, and most alternative paths are absent. Without addressing these gaps, development will encounter UX ambiguity, especially on first launch, error conditions, and data recovery scenarios. The PIN flow (§7.2) is the only area with basic alternative and recovery coverage.

### Old-to-New Question ID Mapping

| Old ID | New ID | Action |
|---|---|---|
| Q-ACTOR-01 | Q-ACTOR-01 | Preserved |
| Q-ACTOR-02 | Q-ACTOR-02 | Preserved |
| Q-ACTOR-03 | Q-ACTOR-03 | Preserved |
| Q-ACTOR-04 | Q-ACTOR-04 | Preserved |
| Q-ACTOR-05 | Q-ACTOR-05 | Preserved |
| Q-ACTOR-06 | Q-ACTOR-06 | Preserved |
| Q-ACTOR-07 | Q-ACTOR-07 | Consolidated with Q-ACTOR-19 (old) |
| Q-ACTOR-08 | Q-ACTOR-08 | Preserved |
| Q-ACTOR-09 | Q-ACTOR-09 | Preserved, tagged cross-scope |
| Q-ACTOR-10 | Q-ACTOR-10 | Preserved |
| Q-ACTOR-11 | → Q-ACTOR-12 | Folded into first-run convention |
| Q-ACTOR-12 | → Q-ACTOR-12 | Folded into first-run convention |
| Q-ACTOR-13 | → Q-ACTOR-12 | Folded into first-run convention |
| Q-ACTOR-14 | → Q-ACTOR-12 | Folded into first-run convention |
| Q-ACTOR-15 | → Q-ACTOR-12 | Folded into first-run convention |
| Q-ACTOR-16 | → Q-ACTOR-12 | Folded into first-run convention |
| Q-ACTOR-17 | → Q-ACTOR-12 | Folded into first-run convention |
| Q-ACTOR-18 | → Q-ACTOR-12 | Folded into first-run convention |
| Q-ACTOR-19 | → Q-ACTOR-07 | Consolidated with Q-ACTOR-07 (new) |
| Q-ACTOR-20 | → Q-ACTOR-12 | Folded into first-run convention |
| Q-ACTOR-21 | Q-ACTOR-13 | Renumbered, tagged cross-scope |
| Q-ACTOR-22 | Q-ACTOR-14 | Renumbered |
| Q-ACTOR-23 | Q-ACTOR-15 | Renumbered |
| Q-ACTOR-24 | Q-ACTOR-16 | Renumbered |
| Q-ACTOR-25 | → Q-ACTOR-17 | Consolidated with Q-ACTOR-28 (old) |
| Q-ACTOR-26 | Q-ACTOR-18 | Renumbered |
| Q-ACTOR-27 | Q-ACTOR-19 | Renumbered |
| Q-ACTOR-28 | → Q-ACTOR-17 | Consolidated with Q-ACTOR-25 (old), tagged cross-scope |
| Q-ACTOR-29 | Q-ACTOR-20 | Renumbered |
| Q-ACTOR-30 | Q-ACTOR-21 | Renumbered |
| Q-ACTOR-31 | Q-ACTOR-22 | Renumbered |
| Q-ACTOR-32 | Q-ACTOR-23 | Renumbered |
| Q-ACTOR-33 | Q-ACTOR-24 | Renumbered |
| Q-ACTOR-34 | Q-ACTOR-25 | Renumbered |
| Q-ACTOR-35 | Q-ACTOR-26 | Renumbered |
| Q-ACTOR-36 | Q-ACTOR-27 | Renumbered |
| Q-ACTOR-37 | Q-ACTOR-11 | Renumbered (moved to empty states) |
| Q-ACTOR-38 | Q-ACTOR-28 | Renumbered |