# Review Report: WAVE-07 (AI Export and Prompt Builder)

**Role**: data-api-integration-ops
**Review ID**: review-data-api-integration-ops-attempt-1
**Run**: 20260621T170113Z | **Attempt**: 1
**Planners reviewed**: data-integration-ops (attempt-1), architecture-codebase (attempt-1)

---

## Verdict: **needs-revision**

Both planners contain strong individual work. The data-ops planner excels at ops detail (log markers, cleanup lifecycle, ZIP format spec, disk safeguards). The architecture-codebase planner correctly handles codebase patterns (GraphQL mutations + REST binary endpoints, `AiExportDataProvider` interface, two-step prompt/export flow). However, they diverge on several critical design points that must be reconciled before implementation. No blockers — all issues are reconcilable.

---

## Findings Requiring Revision

### F1 — Include Photos Default Mismatch (MUST FIX)

| Source | Value |
|---|---|
| Domain model invariant #10 | `DEFAULT false` |
| Data-ops planner migration | `DEFAULT false` — correct |
| Architecture-codebase migration | `DEFAULT true` — violates invariant #10 |

**Action**: Align the architecture-codebase migration to `DEFAULT false`. This is the single most critical issue — implementing `DEFAULT true` silently ships a domain invariant violation.

### F2 — Route Convention Mismatch (MUST FIX)

| Aspect | Data-ops (singular, query param) | Architecture-codebase (plural, path param) | Codebase precedent |
|---|---|---|---|
| Generate endpoint | `POST /api/v1/ai-export` | GraphQL mutation `createAiExportPrompt` + `generateAiExport` | GraphQL for mutations |
| Download endpoint | `GET /api/v1/ai-export/download?exportId={uuid}` | `GET /api/v1/ai-exports/{id}/download` | All REST binary routes use path params: `/api/v1/media/{id}`, `/api/v1/progress-photos/{id}` |
| Resource naming | singular `ai-export` | plural `ai-exports` | All existing routes are plural: `media`, `progress-photos`, `exercises` |

**Action**: Adopt the architecture-codebase approach: GraphQL mutations for prompt/metadata, REST `GET /api/v1/ai-exports/{id}/download` for binary download with path parameter. Use plural `ai-exports` consistently.

### F3 — ZIP Storage Path Lacks User-ID Scope (MUST FIX)

| Source | Path pattern |
|---|---|
| Data-ops | `{export_base_path}/{export_uuid}/export.zip` |
| Architecture-codebase | `cfg.Media.BasePath/exports/{userID}/{exportID}.zip` |
| Media handler precedent | `{basePath}/{userID}/{uuid}.{ext}` (indirect — media paths use UUID only) |

The data-ops path lacks a userID directory. While UUID collision risk is negligible, omitting userID from the path makes manual debugging, disk usage attribution, and bulk cleanup harder. The architecture-codebase pattern is preferable.

**Action**: Adopt `{export_base_path}/{userID}/{exportID}.zip`.

### F4 — `display_name` in UserProfile Migration (SHOULD FIX)

Data-ops migration includes `display_name TEXT NOT NULL DEFAULT ''` on `user_profiles`. The domain model lists `displayName` on `UserProfile` (line 37), but `atlas_users` already stores `display_name`. The architecture-codebase correctly omits it.

Duplicating display_name creates a consistency risk: two sources of truth for the same field. The prompt builder does not need display_name in the AI context.

**Action**: Remove `display_name` from the `user_profiles` migration. If display_name access is needed in the UserProfile API, derive it from `atlas_users` via a JOIN.

### F5 — Migration Numbering Must Be Coordinated Across Waves

| Source | Migration numbers |
|---|---|
| Latest existing | `00090_nutrition_tables.sql` |
| Data-ops | 00091, 00092 |
| Architecture-codebase | 00093, 00094 |

If WAVE-05/06 consume 00091 and 00092, WAVE-07 uses 00093/00094 — correct. If those are unused, 00091/00092 is correct. This must be resolved during implementation based on the actual latest migration.

**Action**: Use the actual next available migration numbers at implementation time.

### F6 — Two-Step vs One-Step Export Flow (SHOULD FIX)

Architecture-codebase splits into `GeneratePrompt` (creates record + prompt) and `GenerateExport` (asynchronously builds ZIP). Data-ops has a single synchronous POST that does everything.

The functional spec (§17-18) describes a "Prompt builder" separate from ZIP generation — users should preview the prompt before committing to a potentially expensive ZIP build. The two-step approach better matches product intent and provides a better UX for large exports with photos.

**Action**: Adopt the two-step flow (architecture-codebase approach). The data-ops handler design should be revised to match: `POST` creates prompt + record, a separate endpoint (GraphQL mutation) triggers ZIP generation.

### F7 — Cleanup Task Missing from Architecture-Codebase Plan

Data-ops has a complete cleanup design (TTL config, `ListStaleExports` sqlc query, periodic task, delete order safety). Architecture-codebase mentions cleanup only as a risk mitigation note. Cleanup is critical for ops — stale ZIP files on disk are a liability.

**Action**: Merge the data-ops cleanup design (section 2.5, 3, 5) into the consolidated plan.

### F8 — Merge Log Markers (SHOULD FIX)

Data-ops provides excellent detailed block markers (e.g., `BLOCK_EXPORT_DATA_QUERY`, `BLOCK_EXPORT_ZIP_BUILD`, `BLOCK_EXPORT_ZIP_WRITE`) with 14 defined markers. Architecture-codebase has simpler markers. For ops observability and profiling, the data-ops marker set is far superior.

**Action**: Adopt the data-ops log markers set in full.

---

## Design Points Confirmed Correct

| Point | Assessment |
|---|---|
| `/api/v1/` prefix | Both planners use it — matches codebase. Confirmed. |
| PIN-guarded middleware | Both planners use `AtlasPinGuard` — matches codebase pattern. Confirmed. |
| AI prompt is generated server-side, no external API call | Both planners agree — consistent with functional spec §22 (manual copy-paste to ChatGPT). Confirmed. |
| `AiExportDataProvider` interface pattern | Architecture-codebase approach is clean and preferred over injecting 10+ repos. Confirmed. |
| ZIP format: manifest.json, data.json, summary.md, CSVs, photos/ | Both planners agree — matches functional spec §17-18. Confirmed. |
| Photo files in `photos/` subdirectory (not base64) | Both planners agree — correct choice for usability. Confirmed. |
| Manifest schema version `1.0.0` | Both planners agree — fine. |
| Date range hard limit (365 days) | Data-ops proposes this with open question Q-W07-DIO-01. Acceptable for MVP with product alignment note. |
| `includePhotos` defaults to false in handler | Both planners agree — matches domain model invariant #10 from handler layer. Confirmed. |
| No idempotency key for concurrent generates | Both planners accept this risk with frontend debouncing. Acceptable for MVP. |
| Export ZIP is hard-deleted (record + file) on cleanup | Data-ops section 2.5 makes this explicit. Confirmed. |
| Bootstrap creates default UserProfile | Both planners agree. Confirmed. |

---

## Consolidated Recommendation

The implementation should:
1. **Adopt architecture-codebase** as the structural template (GraphQL + REST split, two-step flow, `AiExportDataProvider`, codebase-consistent patterns)
2. **Overlay data-ops** material onto it: the ZIP format spec, log markers, cleanup lifecycle, config structure, validation rules, and AC/EC sets
3. **Fix the discrepancies** listed above (F1–F8) during consolidation

The planners collectively produce a complete spec — neither is wrong in isolation, but they differ on enough detail that implementation from either alone would produce a misaligned result.