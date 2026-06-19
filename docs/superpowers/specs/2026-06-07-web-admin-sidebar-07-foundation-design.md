<!-- FILE: docs/superpowers/specs/2026-06-07-web-admin-sidebar-07-foundation-design.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Capture the approved design for adapting shadcn sidebar-07 into the web-admin application shell foundation. -->
<!--   SCOPE: Design-level architecture, component boundaries, route metadata, shell states, verification, and GRACE handoff; excludes implementation code. -->
<!--   DEPENDS: AGENTS.md, docs/requirements.xml, docs/technology.xml, docs/development-plan.xml, docs/knowledge-graph.xml, docs/verification-plan.xml, docs/superpowers/specs/2026-06-05-web-admin-shadcn-ui-kit-design.md, apps/web-admin. -->
<!--   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRACE-WORKFLOW / V-M-GRACE-WORKFLOW. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Goal - Defines sidebar-07 as the target web-admin shell foundation. -->
<!--   Current Context - Summarizes the existing UI-kit and page-shell state. -->
<!--   Key Decisions - Captures approved brainstorming choices. -->
<!--   Architecture And Layout - Defines the adapted sidebar-07 shell structure. -->
<!--   Components And Boundaries - Defines primitive, shell, navigation, and import ownership. -->
<!--   Route Data Flow And States - Defines static route metadata, active navigation, breadcrumbs, and shell-owned state. -->
<!--   Testing And Verification - Defines focused web-admin checks, e2e coverage, and GRACE updates. -->
<!--   Implementation Handoff - Defines what the follow-up implementation plan must cover. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.1 - Addressed subagent review findings for route ownership, sidebar tokens, main landmarks, collapsed-state e2e, and internal UI coverage. -->
<!-- END_CHANGE_SUMMARY -->

# Web-admin sidebar-07 foundation design

**Status:** Approved by subagent review loop
**Date:** 2026-06-07

## Goal

Adapt the shadcn `sidebar-07` block into the `apps/web-admin` foundation so future admin pages can be added inside a ready application shell instead of rebuilding layout, navigation, breadcrumbs, theme actions, and sidebar behavior per page.

The work should produce:

- an adapted `sidebar-07` app shell with a sidebar that collapses to icons;
- template-native navigation content instead of shadcn demo content;
- missing shadcn primitives required by the block under the existing `@shared/ui` surface;
- route metadata for sidebar sections, active matching, and breadcrumbs;
- existing `/`, `/users`, `/users/:id`, and `/ui-kit` routes rendered inside the shell;
- focused tests and e2e checks proving the shell works without weakening existing GraphQL behavior.

Official source references:

- https://ui.shadcn.com/blocks#sidebar-07
- https://ui.shadcn.com/docs/components/sidebar
- https://ui.shadcn.com/r/styles/radix-rhea/sidebar-07.json

## Current Context

`apps/web-admin` is a Vite + React Router + React Query admin app with generated GraphQL client types. A local shadcn-based UI kit already exists under `apps/web-admin/src/shared/ui`, and admin pages under `apps/web-admin/src/pages/**` must import UI only from the bare `@shared/ui` barrel.

Current page-level compositions include:

- `AdminPageShell`;
- `AdminPageHeader`;
- `AdminToolbar`;
- `AdminSection`;
- `AdminEmptyState`;
- `ThemeToggle`;
- a visible `/ui-kit` reference route.

This is a good page kit, but it is not yet an application shell. Each route still owns its outer page container and navigation links. The new foundation should lift global chrome into one shared shell while keeping page data behavior local.

The official `sidebar-07` registry item contains these block files:

- `page.tsx`;
- `components/app-sidebar.tsx`;
- `components/nav-main.tsx`;
- `components/nav-projects.tsx`;
- `components/nav-user.tsx`;
- `components/team-switcher.tsx`.

Its registry dependencies are `sidebar`, `breadcrumb`, `separator`, `collapsible`, `dropdown-menu`, and `avatar`. The `sidebar` primitive itself depends on supporting primitives and helpers such as `sheet`, `tooltip`, `input`, `skeleton`, `button`, `separator`, and `use-mobile`.

## Key Decisions

| Decision           | Choice                                       | Rationale                                                                                                     |
| ------------------ | -------------------------------------------- | ------------------------------------------------------------------------------------------------------------- |
| Target block       | Adapt full `sidebar-07`                      | The user specifically requested the collapsible-icon block, not a generic sidebar.                            |
| Content strategy   | Template-native labels and routes            | Avoids shipping demo `Acme` content while preserving block behavior.                                          |
| Import boundary    | Pages still import UI only from `@shared/ui` | Preserves the existing web-admin UI-kit contract and ESLint boundary.                                         |
| Shell scope        | Global chrome only                           | Sidebar, header, breadcrumbs, theme, and user/team placeholders are shell-owned; data states stay page-owned. |
| Route metadata     | Static app-owned config                      | Keeps active nav and breadcrumbs deterministic and testable.                                                  |
| Placeholder routes | Disabled nav items by default                | Shows future foundation capacity without creating fake pages.                                                 |
| Dashboard data     | Not included                                 | The shell is foundation work, not a product dashboard feature.                                                |

## Architecture And Layout

The target app shape follows `sidebar-07`:

```text
SidebarProvider
  AppSidebar
  SidebarInset
    AdminShellHeader
    route content
```

`AdminAppShell` owns the `SidebarProvider`, adapted `AppSidebar`, `SidebarInset`, shell header, and content slot. Use a React Router layout route with `<Outlet />` so all current and future admin routes render inside the same shell.

The sidebar uses `collapsible="icon"` so desktop users can collapse it to an icon rail. Mobile behavior should follow the shadcn sidebar primitive, including the sheet/offcanvas behavior provided by the primitive dependency chain.

`AdminShellHeader` follows the block header pattern:

- `SidebarTrigger`;
- vertical `Separator`;
- `Breadcrumb`;
- right-side global controls such as `ThemeToggle` and a compact user/menu action.

Existing page content should sit below the shell header inside a constrained content area. `SidebarInset` should be the only `main` landmark in the shell. Refactor `AdminPageShell` from a `<main>` wrapper into a non-landmark content container, or replace its route usage with an equivalent content wrapper, so existing pages do not create nested `main` landmarks after they move inside `SidebarInset`.

`AdminPageHeader` remains useful for route-specific titles and actions, but it should no longer carry global shell controls. Move `ThemeToggle` out of `AdminPageHeader` and into the shell/header or user menu, then update layout tests to prove `AdminPageHeader` does not render the global theme control by default.

## Components And Boundaries

Add missing primitives under `apps/web-admin/src/shared/ui/primitives`:

- `sidebar`;
- `breadcrumb`;
- `collapsible`;
- `avatar`;
- `sheet`.

Keep existing primitives:

- `button`;
- `dropdown-menu`;
- `separator`;
- `tooltip`;
- `input`;
- `skeleton`;
- current form, table, card, dialog, tabs, alert, and feedback primitives.

Add implementation-only helper:

- `apps/web-admin/src/shared/ui/hooks/use-mobile.ts` for the shadcn sidebar primitive.

Update theme tokens in `apps/web-admin/src/styles.css` for the sidebar primitive:

- add `@theme inline` mappings for `--color-sidebar`, `--color-sidebar-foreground`, `--color-sidebar-primary`, `--color-sidebar-primary-foreground`, `--color-sidebar-accent`, `--color-sidebar-accent-foreground`, `--color-sidebar-border`, and `--color-sidebar-ring`;
- add matching light and dark `--sidebar*` values aligned with the current zinc base and indigo accent direction;
- update `apps/web-admin/src/shared/ui/theme-contract.test.ts` so token drift is caught by focused tests.

Add shell compositions under `apps/web-admin/src/shared/ui/layout`:

- `AdminAppShell`;
- `AdminShellHeader`;
- `AppSidebar`;
- `NavMain`;
- `NavProjects`;
- `NavUser`;
- `TeamSwitcher`.

These compositions may import shadcn primitives, `lucide-react`, and internal UI helpers because they are part of the UI kit. Page files must not import those implementation details directly.

Because `shared` modules cannot import from `app`, route metadata must stay app-owned and flow down into the shell through props. `AdminAppShell` may accept navigation items, user/team placeholder data, current pathname, and breadcrumb data as props, or it may accept navigation items plus a pathname and compute active display from those props. It must not import `apps/web-admin/src/app/admin-navigation.ts`.

The public `apps/web-admin/src/shared/ui/index.ts` barrel should export the approved page-facing shell and primitive surface. If a sidebar subcomponent is only needed internally by `AdminAppShell`, it can remain unexported from the public barrel.

## Route Data Flow And States

Create a static navigation config at `apps/web-admin/src/app/admin-navigation.ts` as the single source of truth for sidebar groups and shell breadcrumbs.

Each route entry should include:

- stable id;
- label;
- path;
- section or group;
- icon key or icon component owned by the UI layer;
- breadcrumb labels;
- optional children;
- optional disabled placeholder state.

Suggested adapted navigation:

```text
Platform
  Overview -> /
  Users -> /users
  UI Kit -> /ui-kit

Reference
  GraphQL/Admin -> disabled placeholder
  System/Settings -> disabled placeholder
```

`App.tsx` or an app-layer layout component reads the current `useLocation()` pathname, imports `admin-navigation.ts`, computes or passes the shell navigation inputs, and renders `AdminAppShell`. This preserves the dependency direction:

```text
app -> pages
app -> shared/ui
shared/ui -> shared/ui only
```

The app layer and shell together must provide:

- active nav item;
- active parent item for detail routes such as `/users/:id`;
- breadcrumb trail.

For dynamic detail routes, the shell should use static breadcrumbs such as `Users / User detail`. It should not fetch user data for breadcrumb names. The page can still render the loaded user name in `AdminPageHeader`.

The shell owns only layout and navigation state:

- expanded or collapsed sidebar state;
- mobile sidebar state via the shadcn primitive;
- active navigation;
- breadcrumbs;
- global theme control;
- user/team placeholder menus.

Pages remain responsible for:

- GraphQL queries and mutations;
- loading, error, empty, not-found, and success states;
- page-specific actions and forms.

Unknown routes can keep the current redirect-to-home behavior unless implementation discovers a better existing route fallback pattern. Placeholder items should not create fake data routes unless the implementation plan explicitly chooses a small placeholder route for browser verification.

## Testing And Verification

Focused web-admin tests should come before broad gates.

Add or update unit/component tests:

- `admin-navigation.test.ts` for route metadata, active matching, breadcrumbs, and disabled placeholders;
- `admin-shell.test.tsx` for shell rendering, active item, sidebar trigger, content slot, breadcrumbs, and global controls;
- `admin-layout.test.tsx` after moving `ThemeToggle` out of `AdminPageHeader`;
- `ui-primitives.test.tsx` for new primitive exports: sidebar, breadcrumb, collapsible, avatar, and sheet if public;
- focused tests for all new handwritten UI-kit source, including internal sidebar provider behavior, trigger behavior, `use-mobile`, and mobile sheet/offcanvas branches where they are part of the imported source surface;
- `theme-contract.test.ts` for sidebar color token mappings and light/dark `--sidebar*` variables;
- existing `App.test.tsx`, `users-page.test.tsx`, `user-detail-page.test.tsx`, and `ui-kit-page.test.tsx` where links, headings, or layout ownership change.

Do not add broad UI-kit coverage exclusions for sidebar, sheet, or `use-mobile`. If a generated shadcn branch is genuinely impossible to cover, the implementation plan must stop and define a narrow coverage-policy change with replacement gates before proceeding.

Update e2e coverage:

- `/` loads with the sidebar shell;
- `/users` still supports the current GraphQL browser flow;
- `/ui-kit` is reachable from the sidebar;
- the sidebar trigger is visible and clickable;
- clicking the trigger on desktop collapses the sidebar to the icon rail;
- clicking the trigger again expands the sidebar;
- the active navigation affordance remains visible in collapsed state;
- collapsed navigation keeps accessible labels or tooltips for icon-only items;
- route content remains visible and is not obscured after collapse and expand.

Focused commands:

```bash
bunx nx test web-admin
bunx nx run web-admin:test-coverage
bunx nx run web-admin:typecheck
bunx nx lint web-admin
bunx nx build web-admin
```

Run this near closeout because the shell changes real browser route chrome:

```bash
bunx nx run web-admin:e2e
```

Update shared GRACE artifacts during implementation:

- `docs/development-plan.xml` should describe the `sidebar-07` app shell under `M-WEB-ADMIN`;
- `docs/knowledge-graph.xml` should list new shell and primitive ownership paths;
- `docs/verification-plan.xml` should include shell tests and e2e assertions under `V-M-WEB-ADMIN`;
- new or meaningfully edited governed files should include file-local GRACE markup.

## Implementation Handoff

The follow-up implementation plan should cover these work packets:

1. Add missing shadcn primitives and helper dependencies in the approved `shared/ui` locations.
2. Add sidebar theme tokens and token contract tests.
3. Add adapted `sidebar-07` shell compositions with template-native navigation passed from app-owned metadata.
4. Add static navigation metadata and active/breadcrumb helpers in the app layer.
5. Move existing routes into the shell, make `SidebarInset` the only main landmark, and remove duplicated global controls from page headers.
6. Expand `/ui-kit` to demonstrate the shell and new primitives.
7. Update focused tests, e2e, and GRACE artifacts.
8. Run focused verification and record evidence.

Implementation must not copy the block into default `src/components` paths, must not leave demo navigation content as product content, and must not weaken the existing `@shared/ui` page import boundary.
