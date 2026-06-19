# acceptance-criteria-reviewer — Review (Attempt 1)

**Run ID**: 20260618T185935Z
**Worker report**: `.tasks/product-docs-verify/20260618T185935Z/scopes/acceptance-criteria-reviewer/worker-attempt-1.md`
**Reviewer scope**: Verify the worker did not add new product behavior, did not miss observable criteria, and correctly identified gaps.

---

## Review Criteria

1. **Fidelity** — All derived criteria trace to PRD text or strong implication.
2. **Completeness** — No observable behavior from the PRD is left un-captured.
3. **Boundary** — No new product behavior was invented.
4. **Open questions** — Missing specifications are honestly flagged.
5. **Structure** — Criteria are observable (pass/fail), not design prescriptions.

---

## Findings

### 1. Fidelity — All criteria trace to source

Each AC entry includes a section/line reference. Spot-checked:

- AC-01 through AC-08 (PIN): All reference §7.2 directly. ✅
- AC-30 (working weight auto-populate): §10.6 L301 + §6.1 L120. ✅
- AC-72 (photo AI opt-in): §14.3 L537 + §17.3 L712. ✅
- AC-108 (prompt asks): §18.5 L858–873 — all 13 items listed individually. ✅

No invented criteria found.

### 2. Completeness — Minor gaps

The worker covers 139 derived criteria. Reviewing the PRD for missed observable behavior:

**Missed — Section 10.2 L243 "если за дату уже есть запись, открывается существующая запись"**: Implicitly covered in AC-23 ✅.

**Missed — Section 10.8 L330 "Автоматическое изменение рабочего веса без подтверждения пользователя не требуется"**: Covered in AC-43 ✅.

**Missed — Section 24.1 L1069–1074 privacy rules**: Covered in AC-139 through AC-144 ✅.

**Potential miss — Section 20.4 L984 "запрет silent partial import"**: Covered in AC-129 ✅.

After full review, I found **one gap**:

- **Section 6.1 L121**: "пользователь не должен каждый день вручную вбивать питание, если оно соответствует недельному шаблону" — this is an explicit usability requirement ("user should not have to manually enter nutrition daily if it matches template"). The worker covers template auto-application (AC-77, AC-81) but does not explicitly state that the UI does NOT require daily nutrition re-entry. This is an observable AC: *When a weekly template exists, the user does not need to enter nutrition for any day that matches the template*. Consider adding as AC-149b or noting in the routine optimization section.

**Minor**: Section 21 L1011 "автоподстановку рабочего веса" is covered in AC-31 but could be split: auto-populate on add (AC-31) vs pre-fill on calendar day view open. Currently merged.

### 3. Boundary — No new product behavior

All criteria are traceable to the PRD. No criteria require features not described in the document. The worker correctly excluded future scope items (Apple Health, Telegram bot, etc.) and correctly flagged them as "not in MVP" constraints.

### 4. Open Questions — Appropriately flagged

The 18 open questions (Q-AC-01 through Q-AC-18) are genuine specification gaps. None are answerable from the PRD alone. Particularly valuable:

- Q-AC-01 (PIN failure behavior) — critical for implementation
- Q-AC-07 (e1RM formula) — affects AC-41 and AC-136
- Q-AC-08 (best set definition) — affects AC-39
- Q-AC-15 (import with existing data) — silent data loss risk if unaddressed

**Suggestion**: Add Q-AC-19: *What is the timezone handling for date-based features (dashboard week, training day, check-in)?* This is implicit in AC-09 through AC-14 and AC-20 but never specified.

### 5. Structure

Criteria are formatted as observable pass/fail statements. They are grouped by product section for traceability. Confidence levels are appropriately assigned. No criteria are phrased as UI/implementation details.

---

## Verdict

**approved** — The worker report is faithful to the source, does not invent behavior, covers the vast majority of observable criteria, and honestly identifies gaps. The one minor gap (explicit "no daily nutrition re-entry" criterion) is partially covered by AC-77/AC-81 and does not warrant revision.

### Minor enhancement suggestion (non-blocking)
Add derived criterion: "When a weekly nutrition template exists, the system does not require the user to enter nutrition for any day that falls within the template's week." (Source: §6.1 L121)

### Open questions to carry forward
- Q-AC-01: PIN failure behavior (retries, lockout, error display)
- Q-AC-02: PIN session TTL
- Q-AC-03: "Last body weight" definition
- Q-AC-04: Week start day and timezone
- Q-AC-05: Weekly check-in reminder trigger
- Q-AC-06: Working weight auto-populate UI behavior
- Q-AC-07: e1RM formula
- Q-AC-08: Best set definition
- Q-AC-09: Progression signal surfacing
- Q-AC-10: Cardio — standalone vs day-attached
- Q-AC-11: Single nutrition template — replacement semantics
- Q-AC-12: Nutrition product deletion behavior
- Q-AC-13: Chart filter scope
- Q-AC-14: Week flags — per-week vs per-export
- Q-AC-15: Import with existing data (merge/replace/error)
- Q-AC-16: Backup CSV files — mandatory vs optional
- Q-AC-17: "Copy previous set" — one-tap vs auto-fill
- Q-AC-18: "Sensitive comments" definition
- Q-AC-19: (suggested) Timezone handling for date features