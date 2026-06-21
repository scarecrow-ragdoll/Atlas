# Review Report: Traceability & Consistency — WAVE-07

**Run ID:** 20260621T170113Z  
**Wave ID:** WAVE-07 (AI Export and Prompt Builder)  
**Reviewer role:** traceability-consistency  
**Attempt:** 1  
**Verdict: needs-revision**

---

## Findings

### F1: CRITICAL — `include_photos` default value contradiction (architecture-codebase vs all sources)

| Source | Default value |
|--------|--------------|
| Source wave domain-model.md invariant #10 | `false` |
| RULE-025, AC-077, AC-112 | `false` (opt-in) |
| product-ac planner AC-W07-005 | `false` |
| data-integration-ops migration §2.2 | `DEFAULT false` |
| security-compliance planner §3 | `false` |
| **architecture-codebase SLICE-W07-007 migration** | **`DEFAULT true`** |

The architecture-codebase planner's `ai_exports` migration defines `include_photos BOOLEAN NOT NULL DEFAULT true` (line 240). This directly violates RULE-025, AC-112, and the domain model invariant. This is a clear bug — the default must be `false`.

---

### F2: CRITICAL — Fundamental architecture contradiction: UserProfile vs WAVE-01 Settings

| Planner | Position |
|---------|----------|
| product-ac | Creates `UserProfile` as new entity (CAP-W07-001) |
| architecture-codebase | Creates full UserProfile: migration, model, repo, service, resolver, GraphQL schema (SLICES W07-001–W07-006) |
| data-integration-ops | Creates user_profiles migration + `GET /api/v1/user-profile` REST endpoint |
| testing-exit | Writes full test suite for UserProfile CRUD |
| security-compliance | No UserProfile creation discussion (assumes it exists) |
| **sequencing-fit** | **Flags this as a gap: WAVE-01 Settings already stores AI context fields (`ai_goal`, `ai_height`, `ai_age`, `ai_experience`, etc.). Recommends reusing Settings instead of creating UserProfile.** |

**DQ-W07-001** captures this question, but 4 of 6 planners already committed to creating UserProfile as a new entity. If WAVE-01 Settings is already deployed with AI context fields, this is data duplication. This must be resolved before implementation — either confirm WAVE-01 does NOT have these fields, or reconverge on Settings reuse.

---

### F3: HIGH — Missing section toggle fields across planners

| Toggle field | product-ac | architecture-codebase | data-integration-ops |
|---|---|---|---|
| `includePhotos` | ✅ | ✅ | ✅ |
| `includeNutrition` | ✅ | ✅ | ✅ |
| `includeCardio` | ✅ | ✅ | ✅ |
| `includeMeasurements` | ✅ | ✅ | ✅ |
| **`includeBodyWeight`** | ✅ (AC-W07-005 lists this) | ❌ missing | ❌ missing |
| **`includeWorkouts`** | ✅ (AC-W07-005 lists this) | ❌ missing | ❌ missing |

product-ac's AC-W07-005 lists 6 toggle fields: `includePhotos`, `includeNutrition`, `includeCardio`, `includeMeasurements`, `includeBodyWeight`, `includeWorkouts`. But architecture-codebase and data-integration-ops migrations only define 4 toggle booleans. `includeBodyWeight` and `includeWorkouts` are missing from the table schema.

---

### F4: HIGH — `display_name` field inconsistency in user_profiles

architecture-codebase's `user_profiles` migration does NOT include `display_name`. data-integration-ops's migration DOES include `display_name TEXT NOT NULL DEFAULT ''`. The product-ac planner does not mention `display_name` in the UserProfile field list. These are different table schemas for the same entity.

---

### F5: HIGH — Migration numbering conflict

| Planner | user_profiles migration | ai_exports migration |
|---------|------------------------|---------------------|
| architecture-codebase | 00093_user_profiles.sql | 00094_ai_exports.sql |
| data-integration-ops | 00091_user_profiles.sql | 00092_ai_exports.sql |

Neither planner cites the actual current migration sequence. The `wave-status.md` confirms all planners produced reports — this numbering discrepancy means the actual migration files would collide if both were followed.

---

### F6: MEDIUM — AC/EC/ID numbering and content mismatch across planners

All 6 planners propose different AC sets with overlapping ID ranges:

| Planner | AC count | AC ID pattern | Content overlap |
|---------|----------|---------------|-----------------|
| product-ac | 17 | AC-W07-001–017 | Reference set (from source docs) |
| architecture-codebase | 11 | AC-W07-001–011 | Different body than product-ac at same IDs |
| data-integration-ops | 24 | AC-W07-DIO-001–024 | Separate namespace (good pattern) |
| security-compliance | 10 | AC-W07-SEC-001–010 | Separate namespace (good pattern) |
| sequencing-fit | 16 | AC-W07-001–016 | Different from both product-ac and arch-codebase |
| testing-exit | 20 ECs | EC-W07-001–020 | EC-W07-XXX used by 3 planners with different content |

For example, AC-W07-001 means:
- product-ac: "User profile context stored and retrievable"
- architecture-codebase: "User creates an AI export prompt with date range..."
- sequencing-fit: "AiExport record created via POST /api/ai-export..."

This makes it impossible for a developer to know which AC-W07-001 to implement. All planners using the `AC-W07-XXX` namespace must converge on a single set.

---

### F7: MEDIUM — Site toggle `user_comment` not included in architecture-codebase AiExport model

architecture-codebase `CreateAiExportInput` (SLICE-W07-008) does not include `userComment`/`UserComment` — wait, it does include `UserComment *string`. Let me recheck... Yes, line 303: `UserComment *string \`json:"userComment"\``. Actually this is fine. Strike this finding. But the GraphQL schema (SLICE-W07-013) for `CreateAiExportInput` also includes `userComment: String`. OK, consistent.

---

### F7 (corrected): MEDIUM — Validation overflow in testing-exit planner

testing-exit EC-W07-002 defines: "UserProfile goal, height, birthDate validation rejects empty goal, invalid height (<=0 or >300), and invalid birthDate (future date)". These validation rules appear nowhere in the product-verified docs, source wave, or product-ac planner. Testing planner is inventing constraints not traced to requirements.

---

### F8: MEDIUM — REST endpoint path inconsistency

| Planner | POST path | GET path |
|---------|-----------|----------|
| product-ac | `/api/ai-export/generate` | — |
| architecture-codebase | GraphQL mutation | `GET /api/v1/ai-exports/{id}/download` |
| data-integration-ops | `POST /api/v1/ai-export` | `GET /api/v1/ai-export/download?exportId={uuid}` |
| security-compliance | `POST /api/ai-export` | `GET /api/ai-export/download` |

Three different URL patterns (`/api/ai-export`, `/api/v1/ai-export`, `/api/v1/ai-exports/{id}`). The architecture-codebase planner uses plural resource path with path-parameter ID; data-integration-ops uses singular resource with query parameter.

---

### F9: LOW — Date range max bound inconsistency

| Planner | Max range |
|---------|-----------|
| architecture-codebase | 52 weeks (via Settings setting, SLICE-W07-011) |
| data-integration-ops | 365 days (configurable, `max_range_days`) |
| product-ac | not specified |

These are different but reconcilable (52 weeks ≈ 364 days). Minor.

---

### F10: LOW — question-ledger.md is empty

The `question-ledger.md` is a template with "None yet — initial planner dispatch in progress." Despite 6 planners raising ~20+ questions (Q-W07-001–006 from product-ac, Q-W07-001–004 from architecture-codebase, Q-W07-DIO-001–007 from data-integration-ops, Q-W07-SEC-001–005 from security-compliance, Q-TC-W07-001–006 from testing-exit, DQ-W07-001–005 from sequencing-fit), the ledgers is not populated. This is a process gap, not a content issue.

---

## Summary of Required Revisions

| # | Severity | Finding | Action |
|---|----------|---------|--------|
| F1 | CRITICAL | `include_photos DEFAULT true` in arch migration | Change to `DEFAULT false` |
| F2 | CRITICAL | UserProfile vs WAVE-01 Settings architecture question | Resolve DQ-W07-001 before slicing |
| F3 | HIGH | Missing `includeBodyWeight`/`includeWorkouts` toggles | Add to migration and all ACs |
| F4 | HIGH | `display_name` field in user_profiles | Converge on one schema |
| F5 | HIGH | Migration numbering 00091 vs 00093 | Audit current migrations, pick correct next |
| F6 | MEDIUM | AC-W07-XXX ID collision across 3 planners | Unify into one authoritative set per orchestrator |
| F7 | MEDIUM | Test-EC-W07-002 invents validation rules | Remove unless sourced from product docs |
| F8 | MEDIUM | 3 different REST URL patterns | Converge on one pattern |
| F9 | LOW | Date range max bound differs | Align (52 weeks ≈ 364 days) |
| F10 | LOW | question-ledger.md empty | Populate from all planner questions |

---

## Traceability Assessment

- **Product-verified docs → AC mapping**: Strong in product-ac, security-compliance, data-integration-ops. Weak in architecture-codebase (ACs not directly sourced to specific product ACs).
- **Source wave → outcomes mapping**: Consistent. All 5 OUT-W07 entries and 9 CAP entries are addressed by at least one planner.
- **Codebase evidence**: Good. architecture-codebase and sequencing-fit reference existing code patterns and migration sequences.
- **Prior wave dependencies**: Well-documented by sequencing-fit planner. WAVE-03 stub pattern, WAVE-04 week flag read-only, WAVE-01 dependency clearly identified.
- **Question capture**: Fragmented. Questions spread across 5 planners with different ID schemes (Q-W07-XXX, Q-W07-DIO-XXX, Q-W07-SEC-XXX, Q-TC-W07-XXX, DQ-W07-XXX). No consolidation into question-ledger.md.

---

## Verdict

**needs-revision** — The core wave scope is correctly bounded and the implementation is well-understood, but F1 (include_photos default bug) and F2 (UserProfile/Settings architecture contradiction) are critical blockers for safe implementation. F3–F6 and F8 require coordination across planners before a developer can write coherent code. After resolving F1 (trivial fix) and F2 (requires orchestrator decision on DQ-W07-001), and converging AC IDs (F6), the wave is ready for implementation.