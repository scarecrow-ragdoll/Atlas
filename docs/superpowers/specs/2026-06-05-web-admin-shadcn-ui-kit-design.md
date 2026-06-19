<!-- FILE: docs/superpowers/specs/2026-06-05-web-admin-shadcn-ui-kit-design.md -->
<!-- VERSION: 1.0.1 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Capture the approved design for adding shadcn/ui to web-admin as the required admin UI kit. -->
<!--   SCOPE: Design-level architecture, UI-kit boundaries, demo page scope, existing page migration, enforcement, testing, and GRACE updates; excludes implementation code. -->
<!--   DEPENDS: AGENTS.md, docs/requirements.xml, docs/technology.xml, docs/development-plan.xml, docs/knowledge-graph.xml, docs/verification-plan.xml, apps/web-admin. -->
<!--   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRACE-WORKFLOW / V-M-GRACE-WORKFLOW. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Goal - Defines shadcn/ui as the required UI foundation for admin pages. -->
<!--   Current Context - Summarizes the existing Vite admin page and style surface. -->
<!--   Key Decisions - Captures the approved brainstorming choices. -->
<!--   Architecture - Defines shared/ui primitives, compositions, and public import boundaries. -->
<!--   UI Kit Reference Page - Defines the visible /ui-kit demo page scope. -->
<!--   Existing Page Migration - Defines how home, users, and detail routes move onto the UI kit. -->
<!--   Enforcement - Defines documentation and ESLint boundary rules. -->
<!--   Testing And Verification - Defines focused checks and GRACE validation. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.1 - Required web-admin coverage gate and explicit shadcn placement/coverage treatment after subagent review. -->
<!-- END_CHANGE_SUMMARY -->

# Web-admin shadcn UI kit design

**Status:** Approved
**Date:** 2026-06-05

## Goal

Add shadcn/ui to `apps/web-admin` and make it the required UI foundation for all admin pages.

The work should produce:

- a local `web-admin` UI-kit layer under `apps/web-admin/src/shared/ui`;
- a visible `/ui-kit` reference page linked from the admin home page;
- migrated existing admin pages (`/`, `/users`, `/users/:id`) that use the UI-kit layer;
- a documented rule that admin pages must be built from UI-kit components;
- ESLint import-boundary enforcement so page code cannot bypass `@shared/ui` for UI primitives.

This is a template-quality improvement, not a product-specific feature. The result should help future admin pages start from a consistent, tested component surface.

## Current Context

`apps/web-admin` is currently a Vite + React SPA with React Router, React Query, generated GraphQL types, and a small manual CSS file.

Existing route surface:

- `/` renders a simple home page with a link to users.
- `/users` renders a GraphQL-backed users list and create-user form.
- `/users/:id` renders user loading, error, not-found, and detail states.

Current styling is hand-authored in `apps/web-admin/src/styles.css` through classes such as `page-shell`, `toolbar`, `create-form`, `user-card`, `muted`, and `error-message`. That CSS acts like an informal design system, but it is not componentized or enforceable.

The repository already has the Vite alias required by shadcn (`@/*`) in `apps/web-admin/vite.config.ts` and `apps/web-admin/tsconfig.json`.

The implementation should follow the current shadcn Vite path: add Tailwind CSS integration, keep the `@/*` alias, run shadcn init for project configuration, and add concrete components with the shadcn CLI. The official references used for this design are:

- https://ui.shadcn.com/docs/installation/vite
- https://ui.shadcn.com/docs/cli
- https://ui.shadcn.com/docs/components-json

## Key Decisions

| Decision             | Choice                                      | Rationale                                                                     |
| -------------------- | ------------------------------------------- | ----------------------------------------------------------------------------- |
| UI location          | Local to `apps/web-admin/src/shared/ui`     | The rule is admin-specific; a shared workspace UI package is premature.       |
| Public page import   | `@shared/ui` only                           | Pages should consume the UI kit, not implementation details.                  |
| Component scope      | Broad showcase                              | The template should demonstrate enough primitives for future admin CRUD work. |
| `/ui-kit` visibility | Visible route linked from home              | The UI kit is part of developer experience for the template.                  |
| Existing pages       | Migrate `/`, `/users`, and `/users/:id` now | Avoids creating a new rule with old exceptions.                               |
| Enforcement          | Documentation plus ESLint import boundaries | Matches the requested strictness without brittle raw-HTML scanning.           |
| Raw HTML policing    | Not included                                | Import-level enforcement is stable; tag-level scanning would be fragile.      |

## Architecture

`apps/web-admin/src/shared/ui` becomes the only approved UI-kit entrypoint for admin pages.

Target shape:

```text
apps/web-admin/src/shared/ui
  index.ts
  lib/
    utils.ts
  primitives/
    alert.tsx
    badge.tsx
    button.tsx
    card.tsx
    checkbox.tsx
    dialog.tsx
    dropdown-menu.tsx
    input.tsx
    label.tsx
    select.tsx
    separator.tsx
    skeleton.tsx
    switch.tsx
    table.tsx
    tabs.tsx
    textarea.tsx
    tooltip.tsx
  layout/
    admin-empty-state.tsx
    admin-page-header.tsx
    admin-page-shell.tsx
    admin-section.tsx
    admin-toolbar.tsx
```

The exact generated filenames can follow shadcn CLI output during implementation. The important boundary is public ownership:

- shadcn-generated primitives and their direct dependencies live inside `shared/ui`;
- admin page compositions also live inside `shared/ui`;
- `shared/ui/index.ts` re-exports the approved page-facing surface;
- pages import UI components from `@shared/ui`;
- pages do not import Radix packages, shadcn implementation files, class composition utilities, icon libraries, or UI implementation subpaths directly.

This keeps the page layer focused on business routes and data behavior. It also lets the UI-kit layer absorb future shadcn updates, theme changes, and admin layout decisions without rewriting every page.

The shadcn configuration must make this placement explicit. Use `components.json` aliases and paths, or an equivalent `shadcn add --path` workflow, so generated primitives land under `apps/web-admin/src/shared/ui/primitives` and utilities land under `apps/web-admin/src/shared/ui/lib` rather than the default `src/components/ui` and `src/lib` locations.

## UI Kit Reference Page

Add `/ui-kit` as a normal admin route in `apps/web-admin/src/App.tsx`.

The page should be visible from the home page and should demonstrate the approved component surface for future admin work. It should not be a marketing page or a long tutorial. It should be a practical reference page with realistic admin examples.

Sections:

1. Foundation
   - theme tokens, typography scale, spacing examples, border radius examples;
   - compact text that helps developers see the visual system without adding onboarding copy to production pages.

2. Actions
   - button variants;
   - disabled and loading examples;
   - icon-capable action patterns if icons are added to the UI kit.

3. Forms
   - label plus input;
   - textarea;
   - select;
   - checkbox;
   - switch;
   - validation/error alert example.

4. Feedback
   - alert variants;
   - badges;
   - skeleton loading example;
   - admin empty state composition.

5. Data
   - table with a typical admin CRUD row shape;
   - per-row action dropdown;
   - status badge usage.

6. Overlays And Navigation
   - dialog;
   - dropdown menu;
   - tabs;
   - tooltip;
   - separator.

7. Admin Compositions
   - `AdminPageShell`;
   - `AdminPageHeader`;
   - `AdminToolbar`;
   - `AdminSection`;
   - `AdminEmptyState`.

The page may use local static demo data. It must not require API calls.

## Existing Page Migration

### Home Route

`/` should become a compact admin home built from the UI kit:

- `AdminPageShell`;
- `AdminPageHeader`;
- cards or sections for `Users` and `UI Kit`;
- `Button` or card action links for navigation.

The route should remain simple and should not introduce product-specific dashboard data.

### Users Route

`/users` should preserve all existing GraphQL behavior:

- list query through `graphqlClient`;
- create-user mutation;
- successful mutation invalidates users query and clears the form;
- validation, auth, request, loading, empty, and data states remain user-visible.

The visual structure should move from raw CSS classes to UI-kit components:

- form controls use `Label`, `Input`, and `Button`;
- errors use `Alert`;
- empty state uses `AdminEmptyState`;
- loading state uses `Skeleton`;
- returned users render through `Table` with detail links and status/action patterns where useful;
- page layout uses `AdminPageShell`, `AdminPageHeader`, `AdminToolbar`, and `AdminSection`.

### User Detail Route

`/users/:id` should preserve loading, error, not-found, and detail behavior.

The visual structure should use:

- `AdminPageShell`;
- `AdminPageHeader`;
- `Button` or approved link action back to users;
- `Card` or `AdminSection` for detail groups;
- `Badge` for stable metadata where useful;
- `Skeleton`, `Alert`, and `AdminEmptyState` for non-happy states.

## Styling And Theming

`apps/web-admin/src/styles.css` should become the Tailwind/shadcn theme entrypoint plus minimal global body styles.

The old route-level manual classes should be removed or retired when the pages are migrated. The implementation should not keep a second ad hoc design system in global CSS.

The UI should stay suitable for an admin template:

- quiet, utilitarian, and scan-friendly;
- no marketing hero treatment;
- no decorative gradient/orb background;
- responsive controls that do not overflow on small screens;
- cards only for real grouped items or panels, not nested decorative page sections.

## Enforcement

The rule is:

> Admin pages under `apps/web-admin/src/pages/**` must build UI from `@shared/ui`. They may not import UI primitives, Radix packages, shadcn implementation subpaths, class composition helpers, or icon libraries directly.

Implementation should enforce this in two ways.

### Documentation And Review

Update the durable repository contract during implementation:

- `AGENTS.md` should include the admin UI-kit rule.
- `docs/development-plan.xml` and `docs/knowledge-graph.xml` should reflect the UI-kit layer under `M-WEB-ADMIN`.
- `docs/verification-plan.xml` should list the lint/typecheck/test/build checks that prove the rule and migrated pages.
- `docs/technology.xml` should list Tailwind/shadcn dependencies when they are added.

New or meaningfully edited governed files should carry file-local GRACE markup according to `AGENTS.md`.

### ESLint Import Boundaries

Add a `web-admin` ESLint override for page files, likely in `apps/web-admin/.eslintrc.json`.

The override should prohibit direct UI dependency imports from pages, including:

- `@radix-ui/*` or equivalent Radix primitive packages;
- direct shadcn primitive subpaths such as `@shared/ui/primitives/*`;
- implementation-only UI utility paths such as `@shared/ui/lib/*`;
- `class-variance-authority`, `clsx`, and `tailwind-merge`;
- icon libraries such as `lucide-react`, unless a specific exported icon wrapper is added to `@shared/ui`;
- future raw component paths that bypass `@shared/ui`.

The override should allow normal route/data dependencies:

- `react`;
- `react-router`;
- `@tanstack/react-query`;
- GraphQL documents and generated types;
- `@shared/api/*`;
- `@shared/ui`.

Tag-level scanning for raw HTML controls is out of scope for this design. The repository should use review and UI-kit examples to keep page markup consistent, while ESLint catches architectural bypasses.

## Testing And Verification

Focused implementation checks:

- `bunx nx lint web-admin`;
- `bunx nx test web-admin`;
- `bunx nx run web-admin:test-coverage`;
- `bunx nx run web-admin:typecheck`;
- `bunx nx build web-admin`.

Expected test updates:

- route smoke tests cover the new `/ui-kit` route and home link;
- existing users route tests continue proving loading, empty, returned user, create success, validation/auth error, request error, and pending mutation states;
- detail page tests continue proving fetched user, not-found, and load failure states;
- UI-kit composition tests cover any `shared/ui/layout` component with conditional behavior.

Coverage treatment:

- exported admin layout compositions are handwritten behavior and should be covered by focused tests when they contain branching, conditional rendering, or accessibility behavior;
- shadcn-generated primitives should remain inside the normal `web-admin:test-coverage` contract when practical;
- if implementation treats some generated primitive files as generated code rather than handwritten behavior, it must add narrow allowlist entries to `tools/coverage/coverage.config.json` and matching replacement gates in `docs/verification-plan.xml`;
- broad UI-kit coverage exclusions are not allowed.

GRACE integrity checks:

- `xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml`;
- `grace lint --path .`.

Broader root gates are not required during iteration. A final root build or root coverage gate may be run if the implementation changes shared coverage policy, generated artifacts, or cross-project command surfaces, but the `web-admin:test-coverage` gate is required for this wave because the UI-kit adds source files under `apps/web-admin/src`.

## Risks And Stop Conditions

Stop and replan if:

- shadcn CLI output cannot be cleanly directed under `apps/web-admin/src/shared/ui`;
- Tailwind integration breaks Vite build or Vitest CSS handling;
- ESLint boundaries block legitimate page data dependencies instead of only UI bypasses;
- the migration changes GraphQL behavior rather than only page composition and visual state rendering;
- the UI-kit starts becoming a workspace-wide UI package without a separate design decision.

## Approval

Approved direction:

- local `apps/web-admin` shadcn/ui kit;
- broad visible `/ui-kit` showcase;
- immediate migration of existing admin routes;
- hard rule through docs/review and ESLint import boundaries;
- no tag-level raw HTML scanner in this wave.
