# Review Report: product-scope-and-ac — WAVE-07 Attempt 1

**Run ID:** 20260621T170113Z
**Wave:** WAVE-07 (AI Export and Prompt Builder)
**Reviewer role:** product-scope-and-ac
**Attempt:** 1
**Date:** 2026-06-21

---

## Verdict: needs-revision

The AC coverage is strong overall — all 9 source wave capabilities are addressed, outcomes map correctly, and the excluded scope is clean. Two items need revision before approval.

---

## Required Revisions

### R1 — Missing AC: prompt text returned in generate response body

The product-ac ACs do not include an AC that the `generateAiExport` mutation (or equivalent POST) returns the generated prompt text in the response body. This is required for frontend display/copy:

- PAGE-009 lists "Prompt display/copy" as a page element.
- The sequencing-fit planner explicitly calls this out at §4: "WAVE-07 must return the prompt text in the generate response body so the frontend can display it without downloading the ZIP."

The sequencing-fit planner's draft AC-W07-001 covers this ("Returns export ID, generated prompt text, download URL"), but the product-ac planner's AC-W07-001 only covers UserProfile retrieval. Add an AC (or amend AC-W07-006/AC-W07-013) stating the generate endpoint returns `generatedPrompt` in its response payload.

### R2 — UserProfile scope vs WAVE-01 Settings: conflicting approaches, needs resolution

The product-ac planner defines a full UserProfile CRUD as required new scope: new `user_profile` table migration, model, repository, service, and GraphQL schema with queries and mutations.

The sequencing-fit planner explicitly recommends the opposite at §1.2 and §5.2: reuse WAVE-01's existing Settings table (which stores `ai_goal`, `ai_height`, `ai_age`, `ai_experience`, `ai_split`, `ai_limits`, `ai_progression`, `ai_nutrition_strategy`) rather than adding a new UserProfile entity.

These approaches are contradictory and the product-ac planner does not acknowledge the conflict. Two possible resolutions:

- **Option A (product-ac approach):** Build UserProfile as a new entity. Requires justification in the plan — e.g., UserProfile has fields not in Settings (`birthDate`, `trainingExperience`, `currentTrainingSplit`, `persistentAiContext`), and the domain model separates them. If this is the path, update the plan to explicitly note that this is additive scope beyond what Settings provides and that WAVE-07's prompt service reads from UserProfile, not Settings.

- **Option B (sequencing-fit approach):** Do not create a UserProfile table. Read AI context directly from WAVE-01's Settings GraphQL query. Remove all UserProfile-related migrations, models, repositories, services, and GraphQL additions from WAVE-07 scope. Adjust AC-W07-001 accordingly.

**Action:** Resolve DQ-W07-001 (already raised by sequencing-fit planner) before finalizing. Either revise the plan to align with one consistent approach, or document why UserProfile is justified new scope.

---

## AC Coverage Verdict

| Check | Result |
|---|---|
| All 9 CAP-W07-XXX covered? | ✅ Yes — each maps to at least one AC |
| All 5 OUT-W07-XXX covered? | ✅ Yes |
| ACs testable and specific? | ✅ Yes — clear input/output definitions |
| No frontend scope leaked? | ✅ Clean — frontend explicitly excluded |
| No scope stolen from other waves? | ✅ Week flags correctly marked as read-only; no write collision |
| RULE-021 (4-week default) covered? | ✅ AC-W07-002 |
| RULE-025 (photos opt-in) covered? | ✅ AC-W07-005, AC-W07-016 |
| RULE-026 (on-demand) covered? | ✅ AC-W07-002/003 (user initiates with params) |
| RULE-027 (manual copy-paste) covered? | ✅ AC-W07-013 (AI-readable prompt) |
| RULE-029 (no external API) covered? | ✅ In excluded scope §3 |
| EDGE-008 (empty period) covered? | ✅ AC-W07-014 |
| EDGE-024 (disk full) covered? | ✅ AC-W07-015 |
| AC-118 (no content in logs) covered? | ✅ AC-W07-017 |

### Missing but low-severity

- EDGE-018 (deleted exercise in historical data): No explicit AC. Acceptable — covered implicitly by data query implementation (left-join with fallback name).
- EDGE-031 (timezone handling): No explicit AC. Acceptable — all dates are date-only; addressed in product-ac planner §5.
- WAVE-03 stub (empty workout): No explicit AC about workout section specifically. AC-W07-014 (empty period) covers the general case. Low severity.

---

## Traceability Check

| Source | Mapped in report? | Notes |
|---|---|---|
| Source wave CAPs | ✅ | §2 table |
| Functional spec §17-18 | ✅ | §8 traceability |
| Business rules RULE-021/025/026/027 | ✅ | §8 traceability |
| PAGE-009 deps | ✅ | Addressed in both planners |
| WAVE-04 collision | ✅ | Correctly identified |

---

## Summary

The product-ac planner produces well-structured, testable ACs with good edge-case coverage. The two revision items are actionable:

1. Add an AC that the generate response returns the prompt text (R1).
2. Resolve the UserProfile-vs-Settings conflict — either justify the new entity or remove it (R2).

Neither issue blocks the wave, but both must be resolved before the wave is ready for implementation planning.