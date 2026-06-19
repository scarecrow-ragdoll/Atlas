# Actor-Journey Reviewer — Open Question Ledger

**Run ID:** 20260618T185935Z
**Source:** docs/product/prd.md
**Scope:** actor-journey-reviewer (personas, user journeys, happy paths, alternative paths, empty states, recovery paths)

---

| ID | Category | Question | PRD Gaps | Cross-Scope |
|---|---|---|---|---|
| Q-ACTOR-01 | Alternative path | Add exercise to workout when exercise doesn't exist in exercise library — must exit, create, return? Or inline creation? | §26.2, §26.3 describe adding exercises but not the "missing exercise" path | — |
| Q-ACTOR-02 | Alternative path | Add a set with invalid values (0 reps, 0 weight) — validation error behavior? | §10.5 defines set fields but not validation | — |
| Q-ACTOR-03 | Alternative path | Weekly check-in without photos — can user proceed with 0-1 photos? §13.2 says "2-4 photos" — is this required or recommended? | §13.2, §26.5 | — |
| Q-ACTOR-04 | Alternative path | Remove the last exercise from a workout day — day becomes empty, should it be deleted? | No mention of empty workout day deletion | — |
| Q-ACTOR-05 | Alternative path | Add the same exercise twice in one workout day — duplicate entry allowed or one entry with combined sets? | No duplicate exercise guidance in §10 | — |
| Q-ACTOR-06 | Alternative path | Nutrition override reverting to template — user wants to remove all overrides for today. Is there a "reset to template" action? | §15.5 defines adding overrides but not removing them | — |
| Q-ACTOR-07 | Alternative path | Empty data results — what does the UI show when a query returns zero results? (AI export with no data in period, chart with no data for selected period/exercise) | §17 (AI export), §16 (charts) describe behavior with data but not empty queries | — |
| Q-ACTOR-08 | Alternative path | Backup import when data already exists in instance — overwrite, merge, or reject? | §20.4 describes import to "new instance" but not re-import | — |
| Q-ACTOR-09 | Alternative path | User enters wrong PIN — failed attempt feedback, attempt limit, lockout? | §7.2 describes PIN enable/disable but not wrong-PIN handling | roles/permissions |
| Q-ACTOR-10 | Empty state | Dashboard empty state — what does first launch look like with zero data (no weight, no workouts, no goal)? | §9 describes dashboard blocks assuming data exists | — |
| Q-ACTOR-11 | Empty state | Nutrition template expiry — what happens when template week (e.g., week 1) passes and no new template for week 2 is created? Does daily view show blank? Continue last template? | §15.3-§15.4 describe template lifecycle but not expiry/absence | — |
| Q-ACTOR-12 | Empty state | First-run empty-state convention — what is the consistent UI pattern for all sections with no data yet? (exercise library, workout diary, body measurements, progress photos, nutrition products, AI exports, AI reviews, backup history, charts section) | No empty-state convention described anywhere in PRD | — |
| Q-ACTOR-13 | Recovery | PIN lost/forgotten — user cannot access app. No email, no account, no recovery mechanism. What happens? | §7.2 omits PIN recovery entirely | roles/permissions |
| Q-ACTOR-14 | Recovery | ZIP export generation fails — disk full, permission error, timeout. User feedback? | §20.1, §17.4 describe ZIP format but not generation failure | — |
| Q-ACTOR-15 | Recovery | Backup import validation passes but actual restore fails mid-way — partial data? Transaction rollback? | §20.4 step 4 (dry-run) and step 7 (restore) but no failure handling between | — |
| Q-ACTOR-16 | Recovery | Media upload fails — file too large, wrong format, network error. User feedback and retry? | §11.3 describes that upload is possible but not upload failure | — |
| Q-ACTOR-17 | Recovery | Mid-session resilience — database connection lost or session expires during long data entry. Data loss risk, retry, autosave? | No connection/resilience guidance for mid-session failures | roles/permissions (session expiry) |
| Q-ACTOR-18 | Recovery | User navigates away during workout entry — unsaved data lost? Autosave? Confirm-before-navigate? | No save-or-abandon guidance for in-progress entries | — |
| Q-ACTOR-19 | Recovery | Invalid backup ZIP uploaded — clear error message and graceful rejection? | §20.4 step 2 (check manifest) and step 4 (dry-run) but format validation error behavior unspecified | — |
| Q-ACTOR-20 | Recovery | Concurrent access — same day opened in two browser tabs — last write wins? Conflict detection? | §4 states single-user, but single-user can open two tabs | — |
| Q-ACTOR-21 | Recovery | User creates nutrition template then changes goal — template no longer fits, but no recalculation is triggered. Should recalculation be prompted? | §15.3-§15.4, §18.2 (goal change) — no interaction defined | — |
| Q-ACTOR-22 | First-run | App setup on first launch — user opens app → no data → no guidance or onboarding flow | No first-run experience defined | — |
| Q-ACTOR-23 | First-run | Goal/context setup — user must manually open Settings → UserProfile. No prompt to set goal before first AI export. | §18.2 defines profile fields but not setup prompting | — |
| Q-ACTOR-24 | First-run | PIN setup during first run — PIN is optional but no nudge at first launch. Should first launch suggest PIN? | §7.2 defines PIN as optional but not first-run suggestion | — |
| Q-ACTOR-25 | Edge case | Exercise search/filter for large exercise library (>50 exercises) — search/filter UI undefined | §11 defines exercise creation but not browsing/search for large libraries | — |
| Q-ACTOR-26 | Edge case | Cardio entry without pulse/zone — allowed per §12.4 but what default zone label? "Unknown"? | §12.4 says "unknown" zone exists but behavior with no entry unspecified | — |
| Q-ACTOR-27 | Edge case | Multiple workout days navigation — can user view a list of "all days with workouts" or only navigate by calendar? | §10.2 describes date-based access but not list view | — |
| Q-ACTOR-28 | Edge case | Body weight duplicate on same date — if user enters weight twice on same date, what happens? Overwrite last, reject, keep latest? | §13.5 allows weight by date but no duplicate handling | — |

---

## Cross-Scope References

| Question | Scope | Reason |
|---|---|---|
| Q-ACTOR-09 | roles/permissions | PIN wrong-entry policy and lockout is an auth/permission concern |
| Q-ACTOR-13 | roles/permissions | PIN recovery/forgotten mechanism is an auth/permission concern |
| Q-ACTOR-17 | roles/permissions | Session expiry policy is an auth/permission concern |

## Migration from Attempt 1

Open questions reduced from 38 to 28:
- 10 empty-state questions compressed into 1 convention question (Q-ACTOR-12)
- 2 duplicate pairs consolidated (Q-ACTOR-07, Q-ACTOR-17)
- 10 questions renumbered