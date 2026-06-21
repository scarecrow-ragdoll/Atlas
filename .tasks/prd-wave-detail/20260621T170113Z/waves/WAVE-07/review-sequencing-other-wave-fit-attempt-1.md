# Sequencing and Other-Wave Fit Review — WAVE-07

**Run:** 20260621T170113Z  
**Wave:** WAVE-07 — AI Export and Prompt Builder  
**Role:** sequencing-other-wave-fit  
**Attempt:** 1  
**Reviewer verdict:** **needs-revision**

---

## 1. Week Flags CRUD Collision (CAP-W07-003 vs CAP-W04-006) — RESOLVED

The architecture-codebase planner correctly respects the boundary:
- No week flag create/update/delete mutations or schema in WAVE-07
- `WeekFlagService` injected as a read-only dependency via `AiExportDataProvider.GetWeekFlags`
- AC-W07-005 only reads week flags for prompt and data.json inclusion
- EC-W07-009 explicitly prohibits week flag write operations

**Verdict: Compliant.** The source wave claim CAP-W07-003 must be struck from the wave scope doc before implementation, but both planners agree on read-only treatment. No revision needed.

---

## 2. UserProfile vs Settings — SCOPE COLLISION (needs revision)

The sequencing-fit planner raised DQ-W07-001 with a clear recommendation: **reuse WAVE-01 Settings** for AI context and goal data (option b from the report). The setting fields `ai_goal`, `ai_height`, `ai_age`, `ai_experience`, `ai_split`, `ai_limits`, `ai_progression`, `ai_nutrition_strategy` already cover what WAVE-07 needs for prompt generation.

**The architecture-codebase planner violates this recommendation.** It creates a full `user_profiles` table (SLICE-W07-001 through SLICE-W07-006) with 7 nullable profile fields that overlap significantly with WAVE-01 Settings fields:

| UserProfile field | WAVE-01 Settings equivalent |
|---|---|
| `goal` | `ai_goal` |
| `height` | `ai_height` |
| `birth_date` | `ai_age` (age derived from birth_date) |
| `training_experience` | `ai_experience` |
| `current_training_split` | `ai_split` |
| `preferred_progression_style` | `ai_progression` |
| `nutrition_strategy` | `ai_nutrition_strategy` |
| `persistent_ai_context` | new (not in Settings) |

**Problems with this approach:**
1. **Data duplication**: Two sources of truth for goal/context. Will diverge.
2. **Scope creep**: 5 implementation slices for a table that duplicates existing WAVE-01 functionality. WAVE-01 already delivers "persistent AI context" and "user goal storage" via Settings.
3. **No sync mechanism**: The planner does not specify how UserProfile and Settings stay in sync. If a user updates their goal in Settings, does UserProfile also update? If not, which one is authoritative?
4. **PAGE-009 surface mismatch**: PAGE-009 needs `GET /api/user-profile (goal context)`. Rather than a new table + REST endpoint, this can be a simple Settings query.

**The only genuinely new field is `persistent_ai_context`** — a freeform text field for custom AI instructions. This should be added to WAVE-01 Settings (settings table already supports extensible fields), not used to justify a separate table.

**Required revision:** Rebase WAVE-07 to read AI context and goal from WAVE-01 Settings. Remove SLICE-W07-001 through SLICE-W07-006 (migration, model, sqlc, repo, service, resolver). Add `persistent_ai_context` field to either the WAVE-01 Settings table or the WAVE-07 AiExport model as per-session context (not global). PAGE-009's `GET /api/user-profile` should be a GraphQL query on Settings with AI context fields + a lightweight wrapper, not a new entity.

---

## 3. WAVE-03 Workout Stub Pattern — RESOLVED

The stub pattern follows WAVE-06 precedent (DDEC-W06-010) correctly:
- `AiExportDataProvider.GetWorkoutSummary` returns empty slices when WAVE-03 tables don't exist
- AC-W07-009 explicitly tests this: empty workout section for no WAVE-03 data
- No WAVE-03 table creation in WAVE-07 (EC-W07-010)
- Architecture-codebase planner documents conditional behavior in risk table

**Verdict: Compliant.** No revision needed.

---

## 4. WAVE-08 Foundation — RESOLVED

WAVE-07 correctly creates the right foundation for WAVE-08:
- `AiExport.generatedPrompt` stores the full prompt text (available for WAVE-08 reference)
- Clean boundary: AiExport = export record, AiReview = separate feedback record
- WAVE-07 does not create AiReview entities or reference them

**Verdict: Compliant.** No revision needed.

---

## 5. WAVE-09 ZIP Pattern Sharing — RESOLVED

Both planners acknowledge the ZIP pattern similarity but correctly avoid shared utility extraction:
- WAVE-07: per-period AI export ZIP in `exports/{userID}/{exportID}.zip`
- WAVE-09: full data backup with version manifest in a different structure
- No scope collision — different purpose, different ZIP contents

The architecture-codebase planner uses single-file ZIP paths (`exports/{userID}/{exportID}.zip`) rather than the subfolder pattern recommended by the sequencing-fit planner. This is acceptable: the ZIP files have unique UUID names, so no collision with WAVE-09.

**Verdict: Compliant.** Pattern similarity noted, no shared code extraction needed at this stage. Can be revisited during WAVE-09 planning.

---

## 6. PAGE-009 REST vs GraphQL for Week Flags — ACCEPTABLE WITH DOCUMENTATION GAP

The sequencing-fit planner identified DQ-W07-002: PAGE-009 specifies `GET /api/week-flags` (REST), but WAVE-04 exposes `weekFlags(by weekStartDate)` (GraphQL). Recommendation: PAGE-009 uses GraphQL directly.

The architecture-codebase planner's `AiExportDataProvider` reads week flags server-side for prompt inclusion — this is separate from the frontend browsing interaction. For PAGE-009's "Week flags selector" UI element, the frontend must query WAVE-04's GraphQL `weekFlags(weekStartDate:)` query directly.

**Documentation gap:** The planner does not explicitly state that PAGE-009 must use WAVE-04's GraphQL for week flag browsing. This should be documented to prevent a future REST proxy from being added to WAVE-07.

**Required revision:** Add a note to the wave brief or frontend dependency context that week flag selection on PAGE-009 is served by WAVE-04's existing GraphQL `weekFlags` query — no REST endpoint needed from WAVE-07.

---

## Summary of Required Revisions

| # | Issue | Severity | Required Change |
|---|-------|----------|-----------------|
| R1 | UserProfile table duplicates WAVE-01 Settings | **Critical** | Remove SLICE-W07-001–006. Read AI context and goal from WAVE-01 Settings. Add `persistent_ai_context` to Settings or handle per-session. |
| R2 | Week flag REST vs GraphQL documentation gap | Low | Document that PAGE-009 uses WAVE-04 GraphQL `weekFlags` for browsing — no REST proxy needed. |

**R1 is blocking.** If the UserProfile table is created, it introduces a permanent data-duplication problem with WAVE-01 Settings that will cause bugs, drift, and confusion. The architecture-codebase planner must be revised to remove all UserProfile slices.

## Verdict

**needs-revision** — The UserProfile scope collision with WAVE-01 Settings (R1) must be resolved before this wave can proceed. All other sequencing and other-wave-fit concerns are properly addressed.