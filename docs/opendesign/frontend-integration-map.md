<!-- FILE: docs/opendesign/frontend-integration-map.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Track Atlas OpenDesign frontend artifact integration status against the current web-admin routes. -->
<!--   SCOPE: Records implemented route/component/data-wiring status for the nutrition food-log workstream and related AI export page; excludes backend schema design and untracked HTML artifact storage. -->
<!--   DEPENDS: docs/opendesign/atlas-frontend-design-brief.md, apps/web-admin/src/pages/atlas, apps/api/internal/atlas. -->
<!--   LINKS: M-WEB-ADMIN / M-API-NUTRITION / M-API-AI-EXPORT / V-M-WEB-ADMIN. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Integration Table - Current route/component/data-wiring status for implemented OpenDesign Atlas nutrition and AI export pages. -->
<!--   Known Gaps - Missing branch-local HTML artifacts and deferred non-nutrition pages. -->
<!-- END_MODULE_MAP -->

# Atlas Frontend Integration Map

This branch does not track the generated `docs/opendesign/artifacts/*.html` files that earlier frontend integration prompts referenced. The source design guidance currently present in the branch is `docs/opendesign/atlas-frontend-design-brief.md`. Status below maps the implemented project-native React/TSX surfaces to that design brief and the nutrition food-log plan evidence.

| Source area                         | Target route                | Target component                                                    | Status   | Data wiring notes                                                                                                                                                                                                                                  | Known gaps                                                                                                                         |
| ----------------------------------- | --------------------------- | ------------------------------------------------------------------- | -------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| Nutrition Overview / Daily Food Log | `/atlas/nutrition`          | `apps/web-admin/src/pages/atlas/nutrition-overview-page.tsx`        | verified | API-backed via `getAtlasDailyNutritionLog`, `listAtlasNutritionProducts`, `addAtlasDailyNutritionEntry`, `updateAtlasDailyNutritionEntry`, and `deleteAtlasDailyNutritionEntry`; entries show product snapshots, grams, and returned daily totals. | Historical `/atlas/nutrition/overrides/new` route redirects to the factual daily log instead of rendering target-only override UI. |
| Product Library                     | `/atlas/nutrition/products` | `apps/web-admin/src/pages/atlas/product-library-page.tsx`           | verified | API-backed via `listAtlasNutritionProducts({ includeArchived: true })`, create/update/archive/restore product mutations, active/archived filtering, and EN/RU labels.                                                                              | No public food database, barcode scanner, recipe builder, or external nutrition API by design.                                     |
| Weekly Nutrition Template           | `/atlas/nutrition/template` | `apps/web-admin/src/pages/atlas/weekly-nutrition-template-page.tsx` | verified | API-backed via current-template load, create/update template, create/update/delete template items, active product list, and `applyAtlasNutritionTemplateToWeek(..., SEED_EMPTY_DAYS)`.                                                             | Existing item product changes reconcile as create-new plus delete-old because update-item API does not change `productId`.         |
| AI Export Builder                   | `/atlas/ai-export`          | `apps/web-admin/src/pages/atlas/ai-export-builder-page.tsx`         | verified | API-backed through guarded local REST adapter `apps/web-admin/src/pages/atlas/ai-export-api.ts`; `POST /api/ai-export/generate` uses browser credentials and download links use user-scoped `exportId` only.                                       | No external AI API calls; user still manually sends prompt/ZIP to the AI tool.                                                     |

## Backend Export Payload Notes

- Daily nutrition export includes product name snapshots, grams consumed, per-100g macro snapshots, entry macros, daily totals, and notes/meal labels where present.
- Weekly nutrition template export includes planned products and grams for the selected week-start range.
- Legacy daily nutrition override diagnostics remain a separate export block for migration/review context.
- Provider errors now abort export generation instead of silently producing incomplete nutrition data.

## Known Gaps

- Branch-local OpenDesign HTML artifacts are absent, so this map uses the checked-in design brief plus implemented route/component evidence.
- Non-nutrition Atlas pages are outside the nutrition food-log workstream unless listed above.
- `grace lint --path .` is locally blocked in this shell because the `grace` binary is not installed.
