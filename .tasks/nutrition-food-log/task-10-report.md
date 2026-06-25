<!-- FILE: .tasks/nutrition-food-log/task-10-report.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Record Task 10 Weekly Nutrition Plan Editor implementation evidence. -->
<!--   SCOPE: Captures route/component changes, data wiring, RED/GREEN verification, focused checks, known gaps, and self-review notes; excludes later AI export and GRACE shared XML updates. -->
<!--   DEPENDS: docs/superpowers/plans/2026-06-24-nutrition-food-log.md, apps/web-admin/src/pages/atlas/weekly-nutrition-template-page.tsx, apps/web-admin/src/pages/atlas/nutrition-api.ts. -->
<!--   LINKS: M-WEB-ADMIN / M-API-NUTRITION / V-M-WEB-ADMIN. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Scope - Summarizes Task 10 implementation boundaries. -->
<!--   Data Wiring - Records real adapter operations used by the page. -->
<!--   TDD Evidence - Records RED/GREEN checks. -->
<!--   Verification Evidence - Records focused verification outcomes. -->
<!-- END_MODULE_MAP -->

# Task 10 Report: Weekly Nutrition Plan Editor

## Status

DONE_WITH_CONCERNS

## Scope

- Implemented `/atlas/nutrition/template` as an editable weekly nutrition template page.
- Added route, navigation child, breadcrumbs, EN/RU i18n keys, focused CSS, and tests.
- Stayed inside the provided write scope.

## Changed Files

- `apps/web-admin/src/pages/atlas/weekly-nutrition-template-page.tsx`
- `apps/web-admin/src/pages/atlas/weekly-nutrition-template-page.test.tsx`
- `apps/web-admin/src/pages/atlas/nutrition-api.ts`
- `apps/web-admin/src/pages/atlas/nutrition-api.test.ts`
- `apps/web-admin/src/App.tsx`
- `apps/web-admin/src/App.test.tsx`
- `apps/web-admin/src/app/admin-navigation.ts`
- `apps/web-admin/src/app/admin-navigation.test.ts`
- `apps/web-admin/src/app/i18n.tsx`
- `apps/web-admin/src/styles/atlas.css`
- `.tasks/nutrition-food-log/task-10-report.md`

## Data Wiring

- Loads active products with `listAtlasNutritionProducts()`.
- Loads the selected week with `getAtlasNutritionTemplateCurrent(weekStartDate)`.
- Treats `getAtlasNutritionTemplateCurrent(weekStartDate)` returning `null` as an empty editable template; the adapter maps the backend's empty success result `{}` for this current-template query to `null`.
- Still treats `AtlasNutritionApiError` with type `not_found` from current-template load as an empty editable template defensively.
- `Save Template` creates or updates only the template header and reconciles item rows through:
  - `createAtlasNutritionTemplate`
  - `updateAtlasNutritionTemplate`
  - `createAtlasNutritionTemplateItem`
  - `updateAtlasNutritionTemplateItem`
  - `deleteAtlasNutritionTemplateItem`
- Update paths send empty strings for cleared title, notes, meal label, and item notes so backend null-as-keep semantics do not preserve stale text.
- Create paths keep blank optional text as `null`.
- Existing item product changes are reconciled as create-new plus delete-old because the update-item adapter input does not support `productId`.
- `Save Template` does not call `applyAtlasNutritionTemplateToWeek`.
- `Apply to Week` calls `applyAtlasNutritionTemplateToWeek(template.id, 'SEED_EMPTY_DAYS')` only and is disabled while a save is pending.
- Apply result renders created/skipped counts plus per-date status, entry count, and reason.
- Planned weekly totals are calculated client-side from active product macros per 100g times item grams.
- Local draft state is guarded from same-week React Query refetch overwrites once the user has unsaved edits.
- Entry fallback text and entry aria labels use EN/RU i18n keys.

## TDD Evidence

### RED

- `cd apps/web-admin && bun run test -- src/pages/atlas/weekly-nutrition-template-page.test.tsx`
  - FAIL before production code: missing `./weekly-nutrition-template-page` import.
- `cd apps/web-admin && bun run test -- src/App.test.tsx src/app/admin-navigation.test.ts`
  - FAIL before route/nav production wiring: missing Weekly Plan route heading and missing `atlas-nutrition-template` child.
- Additional self-review regression:
  - `cd apps/web-admin && bun run test -- src/pages/atlas/weekly-nutrition-template-page.test.tsx`
  - FAIL before reconciliation fix: changed product on existing item did not create/delete the item.
- Quality review regression:
  - `cd apps/web-admin && bun run test -- src/pages/atlas/weekly-nutrition-template-page.test.tsx src/pages/atlas/nutrition-api.test.ts`
  - FAIL before production fixes: 5 failed / 25 passed.
  - Failures proved empty current-template success mapped to `INTERNAL_ERROR`, cleared text serialized as `null`, same-week query data overwrote unsaved title edits, Apply to Week stayed enabled while save was pending, and RU fallback text stayed hardcoded in English.

### GREEN

- `cd apps/web-admin && bun run test -- src/pages/atlas/weekly-nutrition-template-page.test.tsx src/pages/atlas/nutrition-api.test.ts`
  - PASS after initial Task 10 implementation: 25 tests.
- `cd apps/web-admin && bun run test -- src/pages/atlas/weekly-nutrition-template-page.test.tsx src/pages/atlas/nutrition-api.test.ts`
  - PASS after quality review fixes: 30 tests.
- `cd apps/web-admin && bun run test -- src/App.test.tsx src/app/admin-navigation.test.ts`
  - PASS: 17 tests.

## Verification Evidence

- `cd apps/web-admin && bun run test -- src/pages/atlas/weekly-nutrition-template-page.test.tsx src/pages/atlas/nutrition-api.test.ts src/App.test.tsx src/app/admin-navigation.test.ts`
  - PASS after quality review fixes: 47 tests.
- `cd apps/web-admin && bun run typecheck`
  - PASS.
- `NX_SKIP_NX_CACHE=true bunx nx lint web-admin`
  - PASS.
- `cd apps/web-admin && bun run build`
  - PASS with Vite large chunk warning.
- `find docs -maxdepth 1 -name '*.xml' -print | sort | xargs xmllint --noout`
  - PASS.
- `git diff --check`
  - PASS.
- `grace lint --path .`
  - BLOCKED: `zsh:1: command not found: grace`.

## Self-Review Notes

- Confirmed page imports UI only from bare `@shared/ui`; no direct icon/Radix/class utility imports in the page file.
- Confirmed no `dangerouslySetInnerHTML`, raw HTML runtime import, body weight field, `WorkoutDay`, or `replace_week` UI/API path.
- Confirmed save and apply are separate actions and tests assert save does not apply.
- Confirmed quality review code findings were fixed with regression coverage.
- Shared `docs/*.xml` graph/verification sync was explicitly deferred to Task 13 by the approved nutrition-food-log plan; this fix did not modify `docs/knowledge-graph.xml` or `docs/verification-plan.xml`.
