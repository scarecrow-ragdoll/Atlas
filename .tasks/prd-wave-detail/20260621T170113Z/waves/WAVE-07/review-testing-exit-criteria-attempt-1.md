# Review: Testing-Exit-Criteria — WAVE-07

**Run ID:** 20260621T170113Z  
**Reviewer role:** testing-exit-criteria  
**Attempt:** 1  
**Verdict: needs-revision**

---

## What passes

1. **Source outcome coverage** — All 5 product outcomes (OUT-W07-001 through OUT-W07-005) are mapped to exit criteria and verification obligations.
2. **UserProfile CRUD** — TEST-W07-001 through TEST-W07-006 cover create, update, get-defaults, get-existing, and validation (empty goal, invalid height, future birthDate).
3. **AiExport prompt generation** — TEST-W07-008 through TEST-W07-014 cover all section toggles, persistent context, one-time comment, week flags, empty date range, and no-data period.
4. **ZIP generation** — TEST-W07-015 through TEST-W07-024 cover archive structure, manifest.json, data.json, summary.md, CSV headers/rows, photos included/excluded, photos/ dir absent when opted out. Very thorough (10 tests).
5. **Download & auth guard** — TEST-W07-030 through TEST-W07-035 cover export generation, download, missing auth (both endpoints), download-not-found, user-profile get, and user-profile no-session.
6. **Log privacy** — TEST-W07-007 (no goal in log), TEST-W07-029 (no export content logged).
7. **Codegen drift** — TEST-W07-039 (sqlc) and TEST-W07-040 (gqlgen).
8. **Migration** — TEST-W07-038 (additive migration against test PostgreSQL).
9. **Pattern consistency** — Mock repo embedding, testify, `httptest` + chi for handlers, `INTEGRATION_TESTS=1` guard, `package service_test` / `package handler_test` — all match existing codebase conventions.
10. **Edge cases** — EC-W07-011 / TEST-W07-013 (empty date range), EC-W07-020 / TEST-W07-025-027 (cleanup, disk full).

## What needs revision

1. **Missing: date range validation for `generateAiExport`** — Product-AC AC-W07-003 states: *"validate end >= start, no future dates (or cap to today)."* There is no EC or TEST that verifies rejection of `dateRangeEnd < dateRangeStart` or `dateRangeStart > today`. Add either:
   - A new EC (e.g., EC-W07-021) with a targeted test, or
   - A TEST entry under §2.2 (e.g., TEST-W07-XXX) validating input rejection.

   The existing EC-W07-002 covers only UserProfile field validation, not AiExport date range validation.

2. **Missing: download endpoint ownership validation** — Security report AC-W07-SEC-006 recommends the download endpoint verify `AiExport.userId` matches session user, returning 404 on mismatch. TEST-W07-033 tests "download not found" (non-existent ID) but not "download of another user's export." In single-user MVP this is a no-op, but the security report explicitly recommends it for future-proofing. Add a TEST entry under §2.5 for this scenario.

3. **Q-TC-W07-001 (sync vs async) affects test design** — The testing planner assumes synchronous generation (EC-W07-012: "returns 200 with export ID and triggers ZIP generation on disk"). If the decision is async, EC-W07-012 and TEST-W07-030 must change entirely (202 Accepted, polling endpoint). Mark this as a dependency: resolve Q-TC-W07-001 before writing handler tests.

4. **Missing: max export size enforcement** — Security report AC-W07-SEC-009 proposes a 100MB uncompressed limit. While not a confirmed AC (it is a recommended addition), the question ledger should cross-reference this. If accepted, add EC and TEST.

## Non-blocking observations

- **WeekFlagsByDateRange query** — If the optional `weekFlagsByDateRange(from, to)` query is built, a separate TEST entry should be added under §2.1 or a new subsection.
- **Security §3 concurrent generation limit** — No test required for MVP, but worth a note in the risk section (§3.5) that this is unenforced.
- **Coverage allowlist** — The exclusion of generated sqlc/gqlgen files in §3.4 is correct and consistent with prior waves.

---

## Required actions to reach approval

1. Add EC-W07-021 (or equivalent) for date range validation rejection on export generation.
2. Add TEST entry under §2.5 for download endpoint ownership mismatch (wrong user returns 404).
3. Cross-reference Q-TC-W07-001 as a dependency that may invalidate handler test design.
4. After fixes, re-run this review.