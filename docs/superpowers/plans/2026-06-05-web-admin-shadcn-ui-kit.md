# Web Admin Shadcn UI Kit Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add shadcn/ui to `apps/web-admin`, expose a visible `/ui-kit` reference page, migrate existing admin routes to the approved UI kit, and enforce page-level UI imports through docs and ESLint.

**Architecture:** Keep the UI kit local to `apps/web-admin/src/shared/ui`. shadcn-generated primitives live under `shared/ui/primitives`, admin-specific layout compositions live under `shared/ui/layout`, and pages import UI only from `@shared/ui`. The current GraphQL behavior remains unchanged while the route visuals move from global CSS classes to UI-kit components.

**Tech Stack:** Bun workspaces, Nx 20, Vite 5, React 19, React Router 7, React Query 5, TypeScript 5.5, Tailwind CSS 4, shadcn/ui, Radix primitives through shadcn, Vitest/jsdom, Testing Library, ESLint boundaries, GRACE XML.

---

<!-- FILE: docs/superpowers/plans/2026-06-05-web-admin-shadcn-ui-kit.md -->
<!-- VERSION: 1.0.1 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Provide the task-by-task implementation plan for adding shadcn/ui as the required web-admin UI kit. -->
<!--   SCOPE: Covers Tailwind/shadcn setup, shared UI primitives, admin layout compositions, /ui-kit route, current page migration, ESLint enforcement, GRACE docs, coverage, and final verification; excludes implementation performed by this document. -->
<!--   DEPENDS: docs/superpowers/specs/2026-06-05-web-admin-shadcn-ui-kit-design.md, apps/web-admin, docs/*.xml, tools/coverage/coverage.config.json. -->
<!--   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRACE-WORKFLOW / V-M-GRACE-WORKFLOW. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Header - Defines the implementation goal, architecture, and required execution sub-skill. -->
<!--   Source Spec - Anchors the approved and subagent-reviewed design. -->
<!--   File Structure - Defines files to create and modify before task decomposition. -->
<!--   Tasks - Provides TDD-oriented setup, implementation, verification, and commit steps. -->
<!--   Self-Review - Records spec coverage, placeholder scan, and type consistency checks. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.1 - Tightened review-loop blockers for import boundaries, verification gates, GRACE sync, and UI-kit completeness. -->
<!-- END_CHANGE_SUMMARY -->

## Source Spec

- Design: `docs/superpowers/specs/2026-06-05-web-admin-shadcn-ui-kit-design.md`
- Review-loop commits:
  - `07bcd8b docs(docs): design web-admin shadcn ui kit`
  - `75f1156 docs(docs): refine web-admin ui kit verification`
- Approved decisions:
  - UI-kit layer is local to `apps/web-admin/src/shared/ui`.
  - Pages import UI from `@shared/ui`, not from primitive implementation subpaths.
  - `/ui-kit` is a normal visible route linked from home.
  - Existing `/`, `/users`, and `/users/:id` pages are migrated in this wave.
  - Enforcement is docs/review plus ESLint import boundaries.
  - `bunx nx run web-admin:test-coverage` is required because the UI kit adds source files under `apps/web-admin/src`.
  - `apps/web-admin/vite.config.ts` is the source of truth for `web-admin:test-coverage`.
  - shadcn placement is explicit through `components.json` aliases or equivalent `shadcn add --path`.
  - Broad UI-kit coverage exclusions are not allowed.

## Scope Check

This is one implementation plan. It touches setup, UI primitives, pages, lint policy, and GRACE docs, but those surfaces are tightly coupled: the feature is not complete until the generated component placement, imported page surface, migrated routes, lint rule, and coverage/docs contracts agree.

## File Structure

### Create

- `apps/web-admin/components.json` - shadcn CLI config with aliases pointing at `src/shared/ui`.
- `apps/web-admin/src/shared/ui/lib/utils.ts` - `cn` class merge helper.
- `apps/web-admin/src/shared/ui/lib/utils.test.ts` - focused test for `cn`.
- `apps/web-admin/src/shared/ui/primitives/*.tsx` - shadcn-generated primitives listed in the spec.
- `apps/web-admin/src/shared/ui/primitives/ui-primitives.test.tsx` - export/render smoke coverage for primitive surface used by `/ui-kit`.
- `apps/web-admin/src/shared/ui/layout/admin-page-shell.tsx` - constrained admin page container.
- `apps/web-admin/src/shared/ui/layout/admin-page-header.tsx` - title/description/actions header.
- `apps/web-admin/src/shared/ui/layout/admin-toolbar.tsx` - responsive command toolbar.
- `apps/web-admin/src/shared/ui/layout/admin-section.tsx` - titled section wrapper around `Card`.
- `apps/web-admin/src/shared/ui/layout/admin-empty-state.tsx` - reusable empty/not-found panel.
- `apps/web-admin/src/shared/ui/layout/admin-layout.test.tsx` - layout composition behavior tests.
- `apps/web-admin/src/shared/ui/index.ts` - only public UI import entrypoint for pages.
- `apps/web-admin/src/pages/ui-kit-page.tsx` - visible component reference page.
- `apps/web-admin/src/pages/ui-kit-page.test.tsx` - route-independent `/ui-kit` page smoke tests.
- `.tasks/web-admin-shadcn-ui-kit/verification.md` - command evidence and decisions for this wave.

### Modify

- `apps/web-admin/package.json` - Tailwind/shadcn-related dependencies.
- `bun.lock` - workspace lockfile after dependency install.
- `apps/web-admin/vite.config.ts` - add `@tailwindcss/vite` plugin and update file-local contract.
- `apps/web-admin/src/styles.css` - replace manual mini design system with Tailwind/shadcn theme entrypoint.
- `apps/web-admin/src/App.tsx` - add `/ui-kit` route and update route map contract.
- `apps/web-admin/src/App.test.tsx` - cover `/ui-kit` route and home links.
- `apps/web-admin/src/pages/home.tsx` - migrate home route to UI-kit layout/cards.
- `apps/web-admin/src/pages/users-page.tsx` - migrate list/create visuals while preserving GraphQL behavior.
- `apps/web-admin/src/pages/users-page.test.tsx` - keep behavior assertions aligned with migrated UI.
- `apps/web-admin/src/pages/user-detail-page.tsx` - migrate detail states to UI-kit layout/cards.
- `apps/web-admin/src/pages/user-detail-page.test.tsx` - keep detail state assertions aligned with migrated UI.
- `apps/web-admin/e2e/users-flow.spec.ts` - update browser-flow locators for table markup and add `/ui-kit` route coverage.
- `apps/web-admin/.eslintrc.json` - add page-level `no-restricted-imports`.
- `AGENTS.md` - add admin UI-kit rule for future agents.
- `docs/requirements.xml` - record admin UI-kit constraint.
- `docs/technology.xml` - record Tailwind/shadcn dependencies and checks.
- `docs/development-plan.xml` - add UI-kit ownership to `M-WEB-ADMIN`.
- `docs/knowledge-graph.xml` - add `shared/ui` path and exported UI-kit annotation.
- `docs/verification-plan.xml` - add UI-kit files and lint/coverage assertions to `V-M-WEB-ADMIN`.
- `docs/operational-packets.xml` - update only when `rg -n "web-admin|M-WEB-ADMIN|V-M-WEB-ADMIN|UI" docs/operational-packets.xml` shows stale ownership text.

### Do Not Modify

- `apps/api/**` - GraphQL behavior is preserved.
- `apps/web/**` - public web app is not part of this UI-kit wave.
- `apps/web-admin/src/shared/api/generated/**` - GraphQL generated types are not part of this wave.
- `tools/coverage/coverage.config.json` - keep unchanged in the normal path. If a worker believes a generated shadcn file is impossible to cover, stop and treat it as a plan deviation requiring narrow `apps/web-admin/vite.config.ts` coverage exclusion, matching root coverage allowlist/replacement-gate updates, `docs/verification-plan.xml` updates, and `bun run test:coverage` evidence.

## Execution Discipline

- Start with Task 0 before editing implementation files.
- Do not create source-only intermediate commits in Tasks 1-8. Stage or leave files uncommitted, then commit implementation, docs, and evidence together in Task 9 after GRACE docs validate.
- If any target file already has unrelated local changes, do not broad-stage it. Use a clean worktree or pause for user direction before editing that file.
- JSON configs cannot carry file-local GRACE comments. `apps/web-admin/components.json` and `apps/web-admin/.eslintrc.json` are governed through Task 9 XML/docs entries and `.tasks/web-admin-shadcn-ui-kit/verification.md` evidence instead of inline markup.
- All non-JSON governed files created or meaningfully changed in this plan must include file-local GRACE markup and useful `START_CONTRACT` / `START_BLOCK_*` anchors for public components, helpers, and critical state branches.

## Task 0: Preflight And Evidence Scaffold

**Files:**

- Create: `.tasks/web-admin-shadcn-ui-kit/verification.md`

- [ ] **Step 1: Inspect target-file dirtiness**

Run:

```bash
git status --short -- apps/web-admin bun.lock AGENTS.md docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/knowledge-graph.xml docs/verification-plan.xml docs/operational-packets.xml .tasks/web-admin-shadcn-ui-kit
```

Expected:

- No target files have unrelated local changes.
- If a target file is already dirty, inspect it with `git diff -- <path>` before editing.
- If the dirty hunk is unrelated, switch to a clean worktree or pause for user direction. Do not stage broad shared files that contain unrelated hunks.

- [ ] **Step 2: Create verification evidence scaffold**

Create `.tasks/web-admin-shadcn-ui-kit/verification.md`:

```markdown
<!-- FILE: .tasks/web-admin-shadcn-ui-kit/verification.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Record verification evidence for the web-admin shadcn UI-kit wave. -->
<!--   SCOPE: Captures commands, outcomes, known retries, lint-boundary negative evidence, coverage decisions, and final gate evidence; excludes durable architecture contracts. -->
<!--   DEPENDS: docs/superpowers/plans/2026-06-05-web-admin-shadcn-ui-kit.md, apps/web-admin, docs/*.xml. -->
<!--   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Command Evidence - Lists each focused verification command and result. -->
<!--   Boundary Evidence - Records the expected failing ESLint fixture proving bypass imports are rejected. -->
<!--   Coverage Decision - Records that UI-kit source is covered by web-admin:test-coverage without broad exclusions. -->
<!--   JSON Governance - Records governance coverage for commentless machine-readable config files. -->
<!--   Final Status - States whether the wave is ready for handoff. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.0 - Added verification log for web-admin UI-kit implementation. -->
<!-- END_CHANGE_SUMMARY -->

# Web-admin shadcn UI-kit Verification

## Command Evidence

| Command                                                                                                                                                                | Result  | Notes                  |
| ---------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | ---------------------- |
| `bunx nx test web-admin`                                                                                                                                               | NOT RUN | Filled during Task 10. |
| `bunx nx run web-admin:test-coverage`                                                                                                                                  | NOT RUN | Filled during Task 10. |
| `bunx nx run web-admin:typecheck`                                                                                                                                      | NOT RUN | Filled during Task 10. |
| `bunx nx lint web-admin`                                                                                                                                               | NOT RUN | Filled during Task 10. |
| `bunx nx build web-admin`                                                                                                                                              | NOT RUN | Filled during Task 10. |
| `bunx nx run web-admin:e2e`                                                                                                                                            | NOT RUN | Filled during Task 10. |
| `xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml` | NOT RUN | Filled during Task 10. |
| `grace lint --path .`                                                                                                                                                  | NOT RUN | Filled during Task 10. |

## Boundary Evidence

| Command                                                                              | Expected Result                                              | Actual Result |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------ | ------------- |
| `cd apps/web-admin && bunx eslint src/pages/ui-boundary-lint.fixture.tsx --ext .tsx` | FAIL with `Admin pages must import UI from @shared/ui only.` | NOT RUN       |

## Coverage Decision

UI-kit source under `apps/web-admin/src/shared/ui` is covered by `web-admin:test-coverage` through `apps/web-admin/vite.config.ts`. No broad UI-kit coverage exclusion was added.

## JSON Governance

`apps/web-admin/components.json` and `apps/web-admin/.eslintrc.json` are machine-readable JSON files and cannot contain inline GRACE comments. Their ownership, constraints, and verification are recorded in `docs/technology.xml`, `docs/development-plan.xml`, `docs/knowledge-graph.xml`, `docs/verification-plan.xml`, and this evidence file.

## Final Status

NOT READY - implementation still in progress.
```

- [ ] **Step 3: Do not commit the scaffold alone**

Leave the evidence scaffold uncommitted until Task 9 commits implementation, docs, and evidence together.

## Task 1: Tailwind And shadcn Base Setup

**Files:**

- Create: `apps/web-admin/components.json`
- Create: `apps/web-admin/src/shared/ui/lib/utils.ts`
- Create: `apps/web-admin/src/shared/ui/lib/utils.test.ts`
- Modify: `apps/web-admin/package.json`
- Modify: `bun.lock`
- Modify: `apps/web-admin/vite.config.ts`
- Modify: `apps/web-admin/src/styles.css`

- [ ] **Step 1: Write the failing `cn` helper test**

Create `apps/web-admin/src/shared/ui/lib/utils.test.ts`:

```ts
// FILE: apps/web-admin/src/shared/ui/lib/utils.test.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify shared web-admin UI utility behavior.
//   SCOPE: Covers class composition and Tailwind conflict merging; excludes component rendering.
//   DEPENDS: apps/web-admin/src/shared/ui/lib/utils.ts, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   cn tests - Prove class values are merged and Tailwind conflicts resolve predictably.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added red coverage for the UI class merge helper.
// END_CHANGE_SUMMARY

import { describe, expect, it } from 'vitest';
import { cn } from './utils';

describe('cn', () => {
  it('merges conditional class values and resolves Tailwind conflicts', () => {
    expect(cn('px-2 text-sm', false && 'hidden', ['font-medium'], 'px-4')).toBe(
      'text-sm font-medium px-4',
    );
  });
});
```

- [ ] **Step 2: Run the test and verify it fails**

Run:

```bash
bunx nx test web-admin -- src/shared/ui/lib/utils.test.ts
```

Expected: FAIL with an import error for `./utils`.

- [ ] **Step 3: Install Tailwind and shadcn runtime dependencies**

Run from repo root:

```bash
cd apps/web-admin
bun add class-variance-authority@^0.7.1 clsx@^2.1.1 tailwind-merge@^3.6.0 lucide-react@^1.17.0 radix-ui@^1.4.3 tw-animate-css@^1.4.0
bun add -d tailwindcss@^4.3.0 @tailwindcss/vite@^4.3.0
cd ../..
```

Expected:

- `apps/web-admin/package.json` includes the new dependencies.
- `bun.lock` changes.
- No dependency is added to `apps/web` or root `package.json` unless Bun workspace hoisting updates the lockfile metadata.

- [ ] **Step 4: Create explicit shadcn config**

Create `apps/web-admin/components.json`:

```json
{
  "$schema": "https://ui.shadcn.com/schema.json",
  "style": "new-york",
  "rsc": false,
  "tsx": true,
  "tailwind": {
    "config": "",
    "css": "src/styles.css",
    "baseColor": "zinc",
    "cssVariables": true,
    "prefix": ""
  },
  "aliases": {
    "components": "@/shared/ui/primitives",
    "ui": "@/shared/ui/primitives",
    "utils": "@/shared/ui/lib/utils",
    "lib": "@/shared/ui/lib",
    "hooks": "@/shared/ui/hooks"
  },
  "iconLibrary": "lucide"
}
```

- [ ] **Step 5: Implement the `cn` helper**

Create `apps/web-admin/src/shared/ui/lib/utils.ts`:

```ts
// FILE: apps/web-admin/src/shared/ui/lib/utils.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide shared utility helpers for web-admin UI components.
//   SCOPE: Owns class name composition used by shadcn primitives and admin compositions; excludes component rendering.
//   DEPENDS: clsx, tailwind-merge.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   cn - Merge conditional class values and resolve Tailwind utility conflicts.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added shadcn-compatible class merge helper.
// END_CHANGE_SUMMARY

import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

// START_CONTRACT: cn
//   PURPOSE: Compose conditional class values and resolve Tailwind conflicts for UI-kit internals.
//   INPUTS: { inputs: ClassValue[] - conditional class values from shadcn primitives and layout compositions }
//   OUTPUTS: { string - merged class name string }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: cn
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
```

- [ ] **Step 6: Add Tailwind plugin to Vite**

Update `apps/web-admin/vite.config.ts`:

```ts
import tailwindcss from '@tailwindcss/vite';
import react from '@vitejs/plugin-react';
import { resolve } from 'path';
import { defineConfig } from 'vitest/config';

export default defineConfig({
  plugins: [react(), tailwindcss()],
  test: {
    environment: 'jsdom',
    globals: true,
    include: ['src/**/*.test.{ts,tsx}'],
    passWithNoTests: false,
    setupFiles: ['./vitest.setup.ts'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json-summary'],
      reportsDirectory: '../../dist/coverage/web-admin',
      include: ['src/**/*.{ts,tsx}'],
      exclude: ['src/**/*.test.{ts,tsx}', 'src/main.tsx', 'src/shared/api/generated/**'],
      thresholds: {
        statements: 100,
        branches: 100,
        functions: 100,
        lines: 100,
      },
    },
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
      '@app': resolve(__dirname, './src/app'),
      '@pages': resolve(__dirname, './src/pages'),
      '@widgets': resolve(__dirname, './src/widgets'),
      '@features': resolve(__dirname, './src/features'),
      '@entities': resolve(__dirname, './src/entities'),
      '@shared': resolve(__dirname, './src/shared'),
    },
  },
});
```

Update the `START_CHANGE_SUMMARY` in the same file to:

```ts
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Added Tailwind CSS Vite plugin for shadcn UI primitives.
// END_CHANGE_SUMMARY
```

- [ ] **Step 7: Replace global CSS with Tailwind/shadcn theme entrypoint**

Replace `apps/web-admin/src/styles.css` with:

```css
/* FILE: apps/web-admin/src/styles.css */
/* VERSION: 1.0.1 */
/* START_MODULE_CONTRACT */
/*   PURPOSE: Provide Tailwind and shadcn theme globals for the web-admin Vite SPA. */
/*   SCOPE: Imports Tailwind, defines theme tokens, and applies minimal base styling; excludes route-specific layout classes. */
/*   DEPENDS: @tailwindcss/vite, tw-animate-css, apps/web-admin/src/main.tsx. */
/*   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN. */
/*   ROLE: RUNTIME */
/*   MAP_MODE: SUMMARY */
/* END_MODULE_CONTRACT */
/* START_MODULE_MAP */
/*   global theme - Tailwind/shadcn tokens and base document styles. */
/* END_MODULE_MAP */
/* START_CHANGE_SUMMARY */
/*   LAST_CHANGE: 1.0.1 - Replaced manual admin CSS classes with Tailwind/shadcn theme globals. */
/* END_CHANGE_SUMMARY */

@import 'tailwindcss';
@import 'tw-animate-css';

@custom-variant dark (&:is(.dark *));

@theme inline {
  --color-background: var(--background);
  --color-foreground: var(--foreground);
  --color-card: var(--card);
  --color-card-foreground: var(--card-foreground);
  --color-popover: var(--popover);
  --color-popover-foreground: var(--popover-foreground);
  --color-primary: var(--primary);
  --color-primary-foreground: var(--primary-foreground);
  --color-secondary: var(--secondary);
  --color-secondary-foreground: var(--secondary-foreground);
  --color-muted: var(--muted);
  --color-muted-foreground: var(--muted-foreground);
  --color-accent: var(--accent);
  --color-accent-foreground: var(--accent-foreground);
  --color-destructive: var(--destructive);
  --color-border: var(--border);
  --color-input: var(--input);
  --color-ring: var(--ring);
  --radius-sm: calc(var(--radius) - 4px);
  --radius-md: calc(var(--radius) - 2px);
  --radius-lg: var(--radius);
}

:root {
  --radius: 0.5rem;
  --background: oklch(0.985 0.004 106);
  --foreground: oklch(0.22 0.018 245);
  --card: oklch(1 0 0);
  --card-foreground: oklch(0.22 0.018 245);
  --popover: oklch(1 0 0);
  --popover-foreground: oklch(0.22 0.018 245);
  --primary: oklch(0.39 0.087 184);
  --primary-foreground: oklch(0.985 0.004 106);
  --secondary: oklch(0.94 0.012 246);
  --secondary-foreground: oklch(0.24 0.025 245);
  --muted: oklch(0.94 0.012 246);
  --muted-foreground: oklch(0.48 0.026 245);
  --accent: oklch(0.92 0.028 156);
  --accent-foreground: oklch(0.24 0.025 245);
  --destructive: oklch(0.58 0.19 28);
  --border: oklch(0.88 0.014 245);
  --input: oklch(0.88 0.014 245);
  --ring: oklch(0.52 0.095 184);
}

@layer base {
  * {
    @apply border-border outline-ring/50;
  }

  body {
    @apply bg-background text-foreground;
    margin: 0;
    font-family:
      Inter,
      ui-sans-serif,
      system-ui,
      -apple-system,
      BlinkMacSystemFont,
      'Segoe UI',
      sans-serif;
  }

  button,
  input,
  textarea,
  select {
    font: inherit;
  }
}
```

- [ ] **Step 8: Run setup checks**

Run:

```bash
bunx nx test web-admin -- src/shared/ui/lib/utils.test.ts
bunx nx run web-admin:typecheck
bunx nx build web-admin
```

Expected:

- The utility test passes.
- Typecheck passes.
- Build passes and Tailwind CSS compiles through the Vite plugin.

If Nx reports project graph instability, rerun the same command with `NX_DAEMON=false`, for example:

```bash
NX_DAEMON=false bunx nx build web-admin
```

- [ ] **Step 9: Stage setup changes**

Run:

```bash
git add apps/web-admin/package.json bun.lock apps/web-admin/components.json apps/web-admin/vite.config.ts apps/web-admin/src/styles.css apps/web-admin/src/shared/ui/lib/utils.ts apps/web-admin/src/shared/ui/lib/utils.test.ts
```

Expected: setup files are staged or ready for Task 9. Do not commit until GRACE docs and evidence are synchronized.

## Task 2: Add shadcn Primitives And Public UI Barrel

**Files:**

- Create: `apps/web-admin/src/shared/ui/primitives/*.tsx`
- Create: `apps/web-admin/src/shared/ui/primitives/ui-primitives.test.tsx`
- Create: `apps/web-admin/src/shared/ui/index.ts`

- [ ] **Step 1: Write the failing primitive barrel test**

Create `apps/web-admin/src/shared/ui/primitives/ui-primitives.test.tsx`:

```tsx
// FILE: apps/web-admin/src/shared/ui/primitives/ui-primitives.test.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify the public web-admin UI primitive export surface.
//   SCOPE: Covers rendering and importability for primitives used by admin pages and /ui-kit; excludes visual pixel assertions.
//   DEPENDS: apps/web-admin/src/shared/ui/index.ts, @testing-library/react, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   primitive exports test - Prove primitives are available through @shared/ui and render basic accessible output.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added red coverage for UI primitive exports.
// END_CHANGE_SUMMARY

import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import {
  Alert,
  AlertDescription,
  AlertTitle,
  Badge,
  Button,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Checkbox,
  Input,
  Label,
  Separator,
  Skeleton,
  Switch,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
  Textarea,
} from '@shared/ui';

describe('web-admin UI primitive exports', () => {
  it('renders the primitive set through the public @shared/ui barrel', () => {
    render(
      <div>
        <Alert>
          <AlertTitle>Saved</AlertTitle>
          <AlertDescription>Changes are ready.</AlertDescription>
        </Alert>
        <Badge>Active</Badge>
        <Button>Save</Button>
        <Card>
          <CardHeader>
            <CardTitle>Card title</CardTitle>
            <CardDescription>Card description</CardDescription>
          </CardHeader>
          <CardContent>Card body</CardContent>
        </Card>
        <Label htmlFor="name">Name</Label>
        <Input id="name" defaultValue="Ada" />
        <Textarea aria-label="Notes" defaultValue="Reference notes" />
        <Checkbox aria-label="Enabled" defaultChecked />
        <Switch aria-label="Published" defaultChecked />
        <Separator />
        <Skeleton data-testid="loading-row" />
        <Tabs defaultValue="overview">
          <TabsList>
            <TabsTrigger value="overview">Overview</TabsTrigger>
          </TabsList>
          <TabsContent value="overview">Overview content</TabsContent>
        </Tabs>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Email</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow>
              <TableCell>ada@example.com</TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </div>,
    );

    expect(screen.getByRole('button', { name: 'Save' })).toBeInTheDocument();
    expect(screen.getByText('Active')).toBeInTheDocument();
    expect(screen.getByLabelText('Name')).toHaveValue('Ada');
    expect(screen.getByText('Overview content')).toBeInTheDocument();
    expect(screen.getByText('ada@example.com')).toBeInTheDocument();
  });
});
```

- [ ] **Step 2: Run the test and verify it fails**

Run:

```bash
bunx nx test web-admin -- src/shared/ui/primitives/ui-primitives.test.tsx
```

Expected: FAIL with missing `@shared/ui` exports.

- [ ] **Step 3: Generate shadcn primitives into `shared/ui/primitives`**

Run:

```bash
cd apps/web-admin
bunx shadcn@latest add button input label card alert badge table skeleton dialog dropdown-menu tabs tooltip separator select textarea switch checkbox -y
cd ../..
```

Expected generated paths:

- `apps/web-admin/src/shared/ui/primitives/button.tsx`
- `apps/web-admin/src/shared/ui/primitives/input.tsx`
- `apps/web-admin/src/shared/ui/primitives/label.tsx`
- `apps/web-admin/src/shared/ui/primitives/card.tsx`
- `apps/web-admin/src/shared/ui/primitives/alert.tsx`
- `apps/web-admin/src/shared/ui/primitives/badge.tsx`
- `apps/web-admin/src/shared/ui/primitives/table.tsx`
- `apps/web-admin/src/shared/ui/primitives/skeleton.tsx`
- `apps/web-admin/src/shared/ui/primitives/dialog.tsx`
- `apps/web-admin/src/shared/ui/primitives/dropdown-menu.tsx`
- `apps/web-admin/src/shared/ui/primitives/tabs.tsx`
- `apps/web-admin/src/shared/ui/primitives/tooltip.tsx`
- `apps/web-admin/src/shared/ui/primitives/separator.tsx`
- `apps/web-admin/src/shared/ui/primitives/select.tsx`
- `apps/web-admin/src/shared/ui/primitives/textarea.tsx`
- `apps/web-admin/src/shared/ui/primitives/switch.tsx`
- `apps/web-admin/src/shared/ui/primitives/checkbox.tsx`

If the CLI writes any file under `apps/web-admin/src/components`, move it into `apps/web-admin/src/shared/ui/primitives`, fix imports to use `@/shared/ui/lib/utils`, and update `apps/web-admin/components.json` before continuing.

- [ ] **Step 4: Add GRACE headers to generated primitive files**

Prepend this header pattern to every generated primitive file, with the `FILE` path and `MODULE_MAP` export names adjusted to the actual file:

```tsx
// FILE: apps/web-admin/src/shared/ui/primitives/button.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide the shadcn button primitive for web-admin UI compositions.
//   SCOPE: Owns button variants and primitive rendering; excludes page-specific behavior.
//   DEPENDS: react, class-variance-authority, @radix-ui/react-slot or radix-ui Slot equivalent, apps/web-admin/src/shared/ui/lib/utils.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Button - shadcn button primitive.
//   buttonVariants - Variant class generator for approved button styles.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added shadcn button primitive under the web-admin UI kit.
// END_CHANGE_SUMMARY
```

Use each file's real exports in the `MODULE_MAP`; for example `Card`, `CardHeader`, `CardTitle`, `CardDescription`, `CardContent`, and `CardFooter` for `card.tsx`.

- [ ] **Step 5: Create the public UI barrel**

Create `apps/web-admin/src/shared/ui/index.ts`:

```ts
// FILE: apps/web-admin/src/shared/ui/index.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Expose the approved public UI-kit surface for web-admin pages.
//   SCOPE: Re-exports shadcn primitives and admin layout compositions; excludes implementation-only utility subpaths from page imports.
//   DEPENDS: apps/web-admin/src/shared/ui/primitives, apps/web-admin/src/shared/ui/layout.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: BARREL
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   primitives - Approved shadcn component exports for admin pages.
//   layout - Approved admin page composition component exports.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added curated public web-admin UI-kit barrel without implementation helpers.
// END_CHANGE_SUMMARY

export { Alert, AlertDescription, AlertTitle } from './primitives/alert';
export { Badge } from './primitives/badge';
export { Button } from './primitives/button';
export {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from './primitives/card';
export { Checkbox } from './primitives/checkbox';
export {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogOverlay,
  DialogPortal,
  DialogTitle,
  DialogTrigger,
} from './primitives/dialog';
export {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from './primitives/dropdown-menu';
export { Input } from './primitives/input';
export { Label } from './primitives/label';
export {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectScrollDownButton,
  SelectScrollUpButton,
  SelectSeparator,
  SelectTrigger,
  SelectValue,
} from './primitives/select';
export { Separator } from './primitives/separator';
export { Skeleton } from './primitives/skeleton';
export { Switch } from './primitives/switch';
export {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableFooter,
  TableHead,
  TableHeader,
  TableRow,
} from './primitives/table';
export { Tabs, TabsContent, TabsList, TabsTrigger } from './primitives/tabs';
export { Textarea } from './primitives/textarea';
export { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from './primitives/tooltip';
```

- [ ] **Step 6: Run primitive checks**

Run:

```bash
bunx nx test web-admin -- src/shared/ui/primitives/ui-primitives.test.tsx src/shared/ui/lib/utils.test.ts
bunx nx run web-admin:typecheck
```

Expected: PASS.

- [ ] **Step 7: Stage primitives**

Run:

```bash
git add apps/web-admin/components.json apps/web-admin/package.json bun.lock apps/web-admin/src/shared/ui
```

Expected: files are staged or ready for Task 9. Do not commit until GRACE docs and evidence are synchronized.

## Task 3: Add Admin Layout Compositions

**Files:**

- Create: `apps/web-admin/src/shared/ui/layout/admin-page-shell.tsx`
- Create: `apps/web-admin/src/shared/ui/layout/admin-page-header.tsx`
- Create: `apps/web-admin/src/shared/ui/layout/admin-toolbar.tsx`
- Create: `apps/web-admin/src/shared/ui/layout/admin-section.tsx`
- Create: `apps/web-admin/src/shared/ui/layout/admin-empty-state.tsx`
- Create: `apps/web-admin/src/shared/ui/layout/admin-layout.test.tsx`
- Modify: `apps/web-admin/src/shared/ui/index.ts`

- [ ] **Step 1: Write the failing layout composition tests**

Create `apps/web-admin/src/shared/ui/layout/admin-layout.test.tsx`:

```tsx
// FILE: apps/web-admin/src/shared/ui/layout/admin-layout.test.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify web-admin layout composition behavior.
//   SCOPE: Covers page shell, header actions, toolbar layout, sections, and empty states; excludes page data behavior.
//   DEPENDS: apps/web-admin/src/shared/ui/layout, @testing-library/react, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   admin layout tests - Prove shared admin compositions render accessible structure and optional actions.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added red coverage for admin UI layout compositions.
// END_CHANGE_SUMMARY

import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import {
  AdminEmptyState,
  AdminPageHeader,
  AdminPageShell,
  AdminSection,
  AdminToolbar,
  Button,
} from '@shared/ui';

describe('admin layout compositions', () => {
  it('renders a page shell with header actions and section content', () => {
    render(
      <AdminPageShell>
        <AdminPageHeader
          title="Users"
          description="Manage reference users."
          actions={<Button>New user</Button>}
        />
        <AdminToolbar>
          <Button variant="outline">Refresh</Button>
        </AdminToolbar>
        <AdminSection title="Directory" description="Current users in the system.">
          <p>One User</p>
        </AdminSection>
      </AdminPageShell>,
    );

    expect(screen.getByRole('main')).toHaveClass('min-h-screen');
    expect(screen.getByRole('heading', { name: 'Users' })).toBeInTheDocument();
    expect(screen.getByText('Manage reference users.')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'New user' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Refresh' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Directory' })).toBeInTheDocument();
    expect(screen.getByText('One User')).toBeInTheDocument();
  });

  it('renders empty states with optional action', () => {
    render(
      <AdminEmptyState
        title="No users yet"
        description="Create the first reference user."
        action={<Button>Create user</Button>}
      />,
    );

    expect(screen.getByRole('heading', { name: 'No users yet' })).toBeInTheDocument();
    expect(screen.getByText('Create the first reference user.')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Create user' })).toBeInTheDocument();
  });
});
```

- [ ] **Step 2: Run the test and verify it fails**

Run:

```bash
bunx nx test web-admin -- src/shared/ui/layout/admin-layout.test.tsx
```

Expected: FAIL with missing layout exports.

- [ ] **Step 3: Create layout components**

Create `apps/web-admin/src/shared/ui/layout/admin-page-shell.tsx`:

```tsx
// FILE: apps/web-admin/src/shared/ui/layout/admin-page-shell.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide the standard page container for web-admin routes.
//   SCOPE: Owns responsive main-page spacing and max width; excludes route-specific headers and data behavior.
//   DEPENDS: react, apps/web-admin/src/shared/ui/lib/utils.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminPageShell - Standard web-admin main container.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added standard admin page shell.
// END_CHANGE_SUMMARY

import type { ComponentPropsWithoutRef } from 'react';
import { cn } from '../lib/utils';

// START_CONTRACT: AdminPageShell
//   PURPOSE: Render the standard responsive main container for admin pages.
//   INPUTS: { props: ComponentPropsWithoutRef<'main'> - native main props, optional className, and children }
//   OUTPUTS: { JSX.Element - main landmark with admin page spacing }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: AdminPageShell
export function AdminPageShell({
  className,
  children,
  ...props
}: ComponentPropsWithoutRef<'main'>) {
  return (
    <main
      className={cn(
        'mx-auto flex min-h-screen w-full max-w-6xl flex-col gap-6 px-4 py-8 sm:px-6 lg:px-8',
        className,
      )}
      {...props}
    >
      {children}
    </main>
  );
}
```

Create `apps/web-admin/src/shared/ui/layout/admin-page-header.tsx`:

```tsx
// FILE: apps/web-admin/src/shared/ui/layout/admin-page-header.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide the standard title, description, and action area for web-admin routes.
//   SCOPE: Owns page heading structure and optional actions; excludes navigation and data fetching.
//   DEPENDS: react, apps/web-admin/src/shared/ui/lib/utils.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminPageHeader - Accessible page header composition.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added standard admin page header.
// END_CHANGE_SUMMARY

import type { ReactNode } from 'react';
import { cn } from '../lib/utils';

type AdminPageHeaderProps = {
  title: string;
  description?: string;
  actions?: ReactNode;
  className?: string;
};

// START_CONTRACT: AdminPageHeader
//   PURPOSE: Render the standard admin page heading, description, and action slot.
//   INPUTS: { props: AdminPageHeaderProps - title, optional description, optional actions, and optional className }
//   OUTPUTS: { JSX.Element - accessible page header }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: AdminPageHeader
export function AdminPageHeader({ title, description, actions, className }: AdminPageHeaderProps) {
  return (
    <header
      className={cn('flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between', className)}
    >
      <div className="min-w-0 space-y-1">
        <h1 className="text-2xl font-semibold tracking-normal text-foreground">{title}</h1>
        {description ? (
          <p className="max-w-3xl text-sm text-muted-foreground">{description}</p>
        ) : null}
      </div>
      {actions ? <div className="flex shrink-0 flex-wrap gap-2">{actions}</div> : null}
    </header>
  );
}
```

Create `apps/web-admin/src/shared/ui/layout/admin-toolbar.tsx`:

```tsx
// FILE: apps/web-admin/src/shared/ui/layout/admin-toolbar.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide a responsive toolbar for web-admin filters and commands.
//   SCOPE: Owns horizontal command wrapping; excludes command behavior.
//   DEPENDS: react, apps/web-admin/src/shared/ui/lib/utils.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminToolbar - Responsive admin command toolbar.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added standard admin toolbar.
// END_CHANGE_SUMMARY

import type { ComponentPropsWithoutRef } from 'react';
import { cn } from '../lib/utils';

// START_CONTRACT: AdminToolbar
//   PURPOSE: Render the standard responsive toolbar for admin filters and commands.
//   INPUTS: { props: ComponentPropsWithoutRef<'div'> - native div props, optional className, and children }
//   OUTPUTS: { JSX.Element - responsive toolbar container }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: AdminToolbar
export function AdminToolbar({ className, children, ...props }: ComponentPropsWithoutRef<'div'>) {
  return (
    <div
      className={cn(
        'flex flex-col gap-3 rounded-lg border bg-card p-3 sm:flex-row sm:items-center sm:justify-between',
        className,
      )}
      {...props}
    >
      {children}
    </div>
  );
}
```

Create `apps/web-admin/src/shared/ui/layout/admin-section.tsx`:

```tsx
// FILE: apps/web-admin/src/shared/ui/layout/admin-section.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide a titled card-backed content section for web-admin routes.
//   SCOPE: Owns section heading and card framing; excludes page-level layout and business data behavior.
//   DEPENDS: react, apps/web-admin/src/shared/ui/primitives/card.tsx.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminSection - Card-backed section with optional description.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added standard admin content section.
// END_CHANGE_SUMMARY

import type { ReactNode } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../primitives/card';

type AdminSectionProps = {
  title: string;
  description?: string;
  children: ReactNode;
};

// START_CONTRACT: AdminSection
//   PURPOSE: Render a titled admin content section using the approved card primitive.
//   INPUTS: { props: AdminSectionProps - title, optional description, and section children }
//   OUTPUTS: { JSX.Element - card-backed content section }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: AdminSection
export function AdminSection({ title, description, children }: AdminSectionProps) {
  return (
    <section>
      <Card>
        <CardHeader>
          <CardTitle>{title}</CardTitle>
          {description ? <CardDescription>{description}</CardDescription> : null}
        </CardHeader>
        <CardContent>{children}</CardContent>
      </Card>
    </section>
  );
}
```

Create `apps/web-admin/src/shared/ui/layout/admin-empty-state.tsx`:

```tsx
// FILE: apps/web-admin/src/shared/ui/layout/admin-empty-state.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide a reusable empty/not-found state for web-admin routes.
//   SCOPE: Owns empty-state copy layout and optional action slot; excludes data fetching and route decisions.
//   DEPENDS: react, apps/web-admin/src/shared/ui/primitives/card.tsx.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminEmptyState - Reusable admin empty-state panel.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added standard admin empty state.
// END_CHANGE_SUMMARY

import type { ReactNode } from 'react';
import { Card, CardContent } from '../primitives/card';

type AdminEmptyStateProps = {
  title: string;
  description: string;
  action?: ReactNode;
};

// START_CONTRACT: AdminEmptyState
//   PURPOSE: Render a reusable empty or not-found state for admin routes.
//   INPUTS: { props: AdminEmptyStateProps - title, description, and optional action }
//   OUTPUTS: { JSX.Element - centered empty-state card }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: AdminEmptyState
export function AdminEmptyState({ title, description, action }: AdminEmptyStateProps) {
  return (
    <Card>
      <CardContent className="flex min-h-44 flex-col items-center justify-center gap-3 text-center">
        <div className="space-y-1">
          <h2 className="text-lg font-medium">{title}</h2>
          <p className="text-sm text-muted-foreground">{description}</p>
        </div>
        {action ? <div className="pt-1">{action}</div> : null}
      </CardContent>
    </Card>
  );
}
```

- [ ] **Step 4: Export layout compositions**

Append to `apps/web-admin/src/shared/ui/index.ts`:

```ts
export * from './layout/admin-empty-state';
export * from './layout/admin-page-header';
export * from './layout/admin-page-shell';
export * from './layout/admin-section';
export * from './layout/admin-toolbar';
```

- [ ] **Step 5: Run layout checks**

Run:

```bash
bunx nx test web-admin -- src/shared/ui/layout/admin-layout.test.tsx src/shared/ui/primitives/ui-primitives.test.tsx
bunx nx run web-admin:typecheck
```

Expected: PASS.

- [ ] **Step 6: Stage layout compositions**

Run:

```bash
git add apps/web-admin/src/shared/ui
```

Expected: layout files are staged or ready for Task 9. Do not commit until GRACE docs and evidence are synchronized.

## Task 4: Add Visible `/ui-kit` Reference Route

**Files:**

- Create: `apps/web-admin/src/pages/ui-kit-page.tsx`
- Create: `apps/web-admin/src/pages/ui-kit-page.test.tsx`
- Modify: `apps/web-admin/src/App.tsx`
- Modify: `apps/web-admin/src/App.test.tsx`
- Modify: `apps/web-admin/e2e/users-flow.spec.ts`

- [ ] **Step 1: Write the failing UI kit page test**

Create `apps/web-admin/src/pages/ui-kit-page.test.tsx`:

```tsx
// FILE: apps/web-admin/src/pages/ui-kit-page.test.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify the web-admin UI kit reference page.
//   SCOPE: Covers static component showcase sections and local-only rendering; excludes visual pixel assertions and API behavior.
//   DEPENDS: apps/web-admin/src/pages/ui-kit-page.tsx, @testing-library/react, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   UiKitPage tests - Prove the reference page demonstrates approved UI-kit areas without API calls.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added red coverage for the UI-kit reference page.
// END_CHANGE_SUMMARY

import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { describe, expect, it } from 'vitest';
import UiKitPage from './ui-kit-page';

describe('UiKitPage', () => {
  it('renders the broad UI-kit showcase sections from local data', () => {
    render(
      <MemoryRouter>
        <UiKitPage />
      </MemoryRouter>,
    );

    expect(screen.getByRole('heading', { name: 'UI Kit' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Actions' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Forms' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Feedback' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Data' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Overlays And Navigation' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Admin Compositions' })).toBeInTheDocument();
    expect(screen.getByText('Typography scale')).toBeInTheDocument();
    expect(screen.getByText('Spacing examples')).toBeInTheDocument();
    expect(screen.getByText('Radius examples')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Creating...' })).toHaveAttribute(
      'aria-busy',
      'true',
    );
    expect(screen.getByText('AdminToolbar')).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'No filtered users' })).toBeInTheDocument();
    expect(screen.getByText('ada@example.com')).toBeInTheDocument();
  });
});
```

- [ ] **Step 2: Run the test and verify it fails**

Run:

```bash
bunx nx test web-admin -- src/pages/ui-kit-page.test.tsx
```

Expected: FAIL with missing `./ui-kit-page`.

- [ ] **Step 3: Implement the UI kit page**

Create `apps/web-admin/src/pages/ui-kit-page.tsx`:

```tsx
// FILE: apps/web-admin/src/pages/ui-kit-page.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Render the web-admin UI-kit reference route.
//   SCOPE: Demonstrates approved primitives and admin compositions using local static data; excludes API calls and product-specific workflows.
//   DEPENDS: react-router, apps/web-admin/src/shared/ui.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - Visible admin UI-kit showcase route.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added broad UI-kit reference page.
// END_CHANGE_SUMMARY

import { Link } from 'react-router';
import {
  AdminEmptyState,
  AdminPageHeader,
  AdminPageShell,
  AdminSection,
  AdminToolbar,
  Alert,
  AlertDescription,
  AlertTitle,
  Badge,
  Button,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Checkbox,
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  Input,
  Label,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Separator,
  Skeleton,
  Switch,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
  Textarea,
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@shared/ui';

const demoRows = [
  { id: 'usr_ada', name: 'Ada Lovelace', email: 'ada@example.com', status: 'Active' },
  { id: 'usr_grace', name: 'Grace Hopper', email: 'grace@example.com', status: 'Pending' },
];

// START_CONTRACT: UiKitPage
//   PURPOSE: Render a local-only reference of approved admin UI components and compositions.
//   INPUTS: none.
//   OUTPUTS: { JSX.Element - visible /ui-kit route using static demo data }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: UiKitPage
export default function UiKitPage() {
  return (
    <TooltipProvider>
      <AdminPageShell>
        <AdminPageHeader
          title="UI Kit"
          description="Approved web-admin components and page compositions."
          actions={
            <Button asChild variant="outline">
              <Link to="/">Home</Link>
            </Button>
          }
        />

        <AdminSection title="Foundation" description="Theme tokens used by admin screens.">
          <div className="grid gap-4 lg:grid-cols-2">
            <div className="space-y-3">
              <h3 className="text-base font-semibold">Typography scale</h3>
              <div className="space-y-1">
                <p className="text-2xl font-semibold">Page title</p>
                <p className="text-sm text-muted-foreground">Muted helper text</p>
                <p className="text-xs uppercase tracking-normal text-muted-foreground">
                  Section label
                </p>
              </div>
            </div>
            <div className="space-y-3">
              <h3 className="text-base font-semibold">Spacing examples</h3>
              <div className="flex items-center gap-2">
                <span className="h-4 w-4 rounded-sm bg-primary" />
                <span className="h-6 w-6 rounded-sm bg-primary" />
                <span className="h-8 w-8 rounded-sm bg-primary" />
              </div>
              <h3 className="text-base font-semibold">Radius examples</h3>
              <div className="flex items-center gap-2">
                <span className="h-10 w-16 rounded-sm border bg-card" />
                <span className="h-10 w-16 rounded-md border bg-card" />
                <span className="h-10 w-16 rounded-lg border bg-card" />
              </div>
            </div>
            <div className="grid gap-3 sm:grid-cols-3 lg:col-span-2">
              {[
                'bg-background text-foreground',
                'bg-primary text-primary-foreground',
                'bg-muted text-muted-foreground',
              ].map((className) => (
                <div className={`rounded-md border p-4 text-sm ${className}`} key={className}>
                  {className}
                </div>
              ))}
            </div>
          </div>
        </AdminSection>

        <AdminSection title="Actions" description="Approved command surfaces.">
          <div className="flex flex-wrap gap-2">
            <Button>Primary</Button>
            <Button variant="secondary">Secondary</Button>
            <Button variant="outline">Outline</Button>
            <Button variant="destructive">Destructive</Button>
            <Button disabled>Disabled</Button>
            <Button aria-busy="true" disabled>
              Creating...
            </Button>
          </div>
        </AdminSection>

        <AdminSection title="Forms" description="Form primitives for admin CRUD flows.">
          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="kit-name">Name</Label>
              <Input id="kit-name" defaultValue="Ada Lovelace" />
            </div>
            <div className="space-y-2">
              <Label>Status</Label>
              <Select defaultValue="active">
                <SelectTrigger aria-label="Status">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="active">Active</SelectItem>
                  <SelectItem value="pending">Pending</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2 md:col-span-2">
              <Label htmlFor="kit-notes">Notes</Label>
              <Textarea id="kit-notes" defaultValue="Reference form copy." />
            </div>
            <label className="flex items-center gap-2 text-sm">
              <Checkbox defaultChecked />
              Receive status updates
            </label>
            <label className="flex items-center gap-2 text-sm">
              <Switch defaultChecked />
              Published
            </label>
          </div>
        </AdminSection>

        <AdminSection title="Feedback" description="Status, loading, and empty states.">
          <div className="grid gap-4 lg:grid-cols-2">
            <Alert>
              <AlertTitle>Validation message</AlertTitle>
              <AlertDescription>Email is already used by another user.</AlertDescription>
            </Alert>
            <Card>
              <CardHeader>
                <CardTitle>Loading skeleton</CardTitle>
                <CardDescription>Use for pending list or detail content.</CardDescription>
              </CardHeader>
              <CardContent className="space-y-2">
                <Skeleton className="h-4 w-2/3" />
                <Skeleton className="h-4 w-1/2" />
              </CardContent>
            </Card>
            <AdminEmptyState
              title="No records"
              description="Create a record to populate this table."
            />
            <div className="flex items-center gap-2">
              <Badge>Active</Badge>
              <Badge variant="secondary">Draft</Badge>
              <Badge variant="destructive">Blocked</Badge>
            </div>
          </div>
        </AdminSection>

        <AdminSection title="Data" description="Table pattern for admin collections.">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Email</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {demoRows.map((row) => (
                <TableRow key={row.id}>
                  <TableCell>{row.name}</TableCell>
                  <TableCell>{row.email}</TableCell>
                  <TableCell>
                    <Badge variant={row.status === 'Active' ? 'default' : 'secondary'}>
                      {row.status}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-right">
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button size="sm" variant="outline">
                          Actions
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem>Open</DropdownMenuItem>
                        <DropdownMenuItem>Archive</DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </AdminSection>

        <AdminSection
          title="Overlays And Navigation"
          description="Use overlays for bounded secondary tasks."
        >
          <Tabs defaultValue="dialog">
            <TabsList>
              <TabsTrigger value="dialog">Dialog</TabsTrigger>
              <TabsTrigger value="tooltip">Tooltip</TabsTrigger>
            </TabsList>
            <TabsContent value="dialog" className="pt-4">
              <Dialog>
                <DialogTrigger asChild>
                  <Button variant="outline">Open dialog</Button>
                </DialogTrigger>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Confirm action</DialogTitle>
                    <DialogDescription>This is the approved modal composition.</DialogDescription>
                  </DialogHeader>
                </DialogContent>
              </Dialog>
            </TabsContent>
            <TabsContent value="tooltip" className="pt-4">
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button variant="outline">Hover target</Button>
                </TooltipTrigger>
                <TooltipContent>Tooltip content</TooltipContent>
              </Tooltip>
            </TabsContent>
          </Tabs>
          <Separator className="my-4" />
          <p className="text-sm text-muted-foreground">Separators divide related admin panels.</p>
        </AdminSection>

        <AdminSection title="Admin Compositions" description="Recommended route-building blocks.">
          <div className="grid gap-3 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>AdminPageShell</CardTitle>
                <CardDescription>Wraps every admin route.</CardDescription>
              </CardHeader>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>AdminPageHeader</CardTitle>
                <CardDescription>Standardizes titles, descriptions, and actions.</CardDescription>
              </CardHeader>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>AdminToolbar</CardTitle>
                <CardDescription>Groups filters and route commands.</CardDescription>
              </CardHeader>
              <CardContent>
                <AdminToolbar>
                  <Input aria-label="Filter users" placeholder="Filter users" />
                  <Button variant="outline">Refresh</Button>
                </AdminToolbar>
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>AdminSection</CardTitle>
                <CardDescription>Frames repeated route panels.</CardDescription>
              </CardHeader>
            </Card>
            <AdminEmptyState
              title="No filtered users"
              description="Clear filters to show the full directory."
              action={<Button variant="outline">Clear filters</Button>}
            />
          </div>
        </AdminSection>
      </AdminPageShell>
    </TooltipProvider>
  );
}
```

- [ ] **Step 4: Add `/ui-kit` to the route table**

Update `apps/web-admin/src/App.tsx` imports and routes:

```tsx
import UiKitPage from './pages/ui-kit-page';

// inside <Routes>
<Route path="/ui-kit" element={<UiKitPage />} />;
```

Update `MODULE_MAP` in the same file to mention home, users, user detail, and UI-kit routes.

- [ ] **Step 5: Update route smoke tests**

Add this test to `apps/web-admin/src/App.test.tsx`:

```tsx
it('renders the UI-kit route through the browser router', () => {
  renderApp('/ui-kit');

  expect(screen.getByRole('heading', { name: 'UI Kit' })).toBeInTheDocument();
  expect(screen.getByRole('heading', { name: 'Actions' })).toBeInTheDocument();
});
```

- [ ] **Step 6: Add Playwright coverage for `/ui-kit` navigation**

Update `apps/web-admin/e2e/users-flow.spec.ts` file-local map and change summary:

```ts
// START_MODULE_MAP
//   create/list/detail flow - Creates a user through the UI, verifies table refresh, and opens the detail route.
//   duplicate-email flow - Seeds a duplicate through GraphQL and verifies the admin UI validation error.
//   ui-kit route flow - Verifies home navigation to the visible /ui-kit component reference page.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Added /ui-kit browser coverage and aligned users flow locators with UI-kit table markup.
// END_CHANGE_SUMMARY
```

Append this Playwright test:

```ts
test('home links to the UI kit reference page', async ({ page }) => {
  await page.goto('/');
  await page.getByRole('link', { name: 'Open UI kit' }).click();

  await expect(page).toHaveURL(/\/ui-kit$/);
  await expect(page.getByRole('heading', { name: 'UI Kit' })).toBeVisible();
  await expect(page.getByRole('heading', { name: 'Admin Compositions' })).toBeVisible();
  await expect(page.getByText('AdminToolbar')).toBeVisible();
  await expect(page.getByRole('button', { name: 'Creating...' })).toHaveAttribute(
    'aria-busy',
    'true',
  );
});
```

Do not run Playwright during active iteration unless needed for debugging; Task 10 runs the required `bunx nx run web-admin:e2e` gate.

- [ ] **Step 7: Run UI-kit route checks**

Run:

```bash
bunx nx test web-admin -- src/pages/ui-kit-page.test.tsx src/App.test.tsx
bunx nx run web-admin:typecheck
```

Expected: PASS.

- [ ] **Step 8: Stage `/ui-kit` changes**

Run:

```bash
git add apps/web-admin/src/App.tsx apps/web-admin/src/App.test.tsx apps/web-admin/src/pages/ui-kit-page.tsx apps/web-admin/src/pages/ui-kit-page.test.tsx apps/web-admin/e2e/users-flow.spec.ts
```

Expected: route files are staged or ready for Task 9. Do not commit until GRACE docs and evidence are synchronized.

## Task 5: Migrate Home Route To UI Kit

**Files:**

- Modify: `apps/web-admin/src/pages/home.tsx`
- Modify: `apps/web-admin/src/App.test.tsx`

- [ ] **Step 1: Update route smoke expectations first**

Replace the home test in `apps/web-admin/src/App.test.tsx` with:

```tsx
it('renders the home route with users and UI-kit links', () => {
  renderApp('/');

  expect(screen.getByRole('heading', { name: 'Monorepo Template Admin' })).toBeInTheDocument();
  expect(screen.getByRole('link', { name: 'Open users' })).toHaveAttribute('href', '/users');
  expect(screen.getByRole('link', { name: 'Open UI kit' })).toHaveAttribute('href', '/ui-kit');
});
```

- [ ] **Step 2: Run the test and verify it fails**

Run:

```bash
bunx nx test web-admin -- src/App.test.tsx
```

Expected: FAIL because `Open UI kit` is not on home yet.

- [ ] **Step 3: Replace home route implementation**

Replace `apps/web-admin/src/pages/home.tsx` with:

```tsx
// FILE: apps/web-admin/src/pages/home.tsx
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Render the web-admin home route.
//   SCOPE: Shows admin entry cards for users and UI-kit reference routes; excludes data fetching and mutation behavior.
//   DEPENDS: react-router, apps/web-admin/src/shared/ui.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - Admin home route with users and UI-kit navigation cards.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Migrated home route to the web-admin UI kit and added UI-kit navigation.
// END_CHANGE_SUMMARY

import { Link } from 'react-router';
import {
  AdminPageHeader,
  AdminPageShell,
  Button,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@shared/ui';

// START_CONTRACT: HomePage
//   PURPOSE: Render admin route entry cards for users and the UI-kit reference page.
//   INPUTS: none.
//   OUTPUTS: { JSX.Element - home route with admin navigation cards }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: HomePage
export default function HomePage() {
  return (
    <AdminPageShell>
      <AdminPageHeader
        title="Monorepo Template Admin"
        description="GraphQL admin client and UI reference for new admin pages."
      />

      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Users</CardTitle>
            <CardDescription>Manage the reference GraphQL users flow.</CardDescription>
          </CardHeader>
          <CardContent>
            <Button asChild>
              <Link to="/users">Open users</Link>
            </Button>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>UI Kit</CardTitle>
            <CardDescription>Review the approved components for admin pages.</CardDescription>
          </CardHeader>
          <CardContent>
            <Button asChild variant="outline">
              <Link to="/ui-kit">Open UI kit</Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    </AdminPageShell>
  );
}
```

- [ ] **Step 4: Run home checks**

Run:

```bash
bunx nx test web-admin -- src/App.test.tsx
bunx nx run web-admin:typecheck
```

Expected: PASS.

- [ ] **Step 5: Stage home migration**

Run:

```bash
git add apps/web-admin/src/pages/home.tsx apps/web-admin/src/App.test.tsx
```

Expected: home files are staged or ready for Task 9. Do not commit until GRACE docs and evidence are synchronized.

## Task 6: Migrate Users Route To UI Kit

**Files:**

- Modify: `apps/web-admin/src/pages/users-page.tsx`
- Modify: `apps/web-admin/src/pages/users-page.test.tsx`
- Modify: `apps/web-admin/e2e/users-flow.spec.ts`

- [ ] **Step 1: Update users tests for table semantics and preserved behavior**

In `apps/web-admin/src/pages/users-page.test.tsx`, keep the existing tests and update assertions that look for the list to also prove table semantics:

```tsx
expect(await screen.findByRole('link', { name: 'One User' })).toHaveAttribute(
  'href',
  '/users/user-1',
);
expect(screen.getByRole('columnheader', { name: 'Name' })).toBeInTheDocument();
expect(screen.getByRole('cell', { name: 'one@example.com' })).toBeInTheDocument();
expect(screen.getByText('Total: 1')).toBeInTheDocument();
expect(screen.getByText('Showing the latest 20 users.')).toBeInTheDocument();
```

Add this assertion to the empty-state test:

```tsx
expect(screen.getByRole('heading', { name: 'No users yet' })).toBeInTheDocument();
```

- [ ] **Step 2: Run users tests and verify they fail against old markup**

Run:

```bash
bunx nx test web-admin -- src/pages/users-page.test.tsx
```

Expected: FAIL because the current route uses a list, not UI-kit table/empty-state headings.

- [ ] **Step 3: Replace users route implementation**

Replace `apps/web-admin/src/pages/users-page.tsx` with:

```tsx
// FILE: apps/web-admin/src/pages/users-page.tsx
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Render the web-admin users list and create-user route.
//   SCOPE: Loads users, displays list states, submits create-user GraphQL mutations, and links to details through UI-kit components; excludes detail rendering and GraphQL transport construction.
//   DEPENDS: @tanstack/react-query, react, react-router, apps/web-admin/src/entities/user/api/*.graphql, apps/web-admin/src/shared/api/graphql-client.ts, generated GraphQL types, apps/web-admin/src/shared/ui.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - Users list and create-form route backed by codegen-visible GraphQL documents and UI-kit components.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Migrated users route visuals to the web-admin UI kit without changing GraphQL behavior.
// END_CHANGE_SUMMARY

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import createUserMutationDocument from '@entities/user/api/createUser.graphql?raw';
import getUsersQueryDocument from '@entities/user/api/users.graphql?raw';
import { graphqlClient } from '@shared/api/graphql-client';
import type { CreateUserMutation, GetUsersQuery } from '@shared/api/generated/types';
import { type FormEvent, useState } from 'react';
import { Link } from 'react-router';
import {
  AdminEmptyState,
  AdminPageHeader,
  AdminPageShell,
  AdminSection,
  AdminToolbar,
  Alert,
  AlertDescription,
  AlertTitle,
  Button,
  Input,
  Label,
  Skeleton,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@shared/ui';

type FormState = {
  name: string;
  email: string;
  password: string;
};

const initialFormState: FormState = { name: '', email: '', password: '' };

// START_CONTRACT: errorMessageFromUnknown
//   PURPOSE: Convert unknown mutation errors into user-visible fallback copy.
//   INPUTS: { error: unknown - mutation error from React Query }
//   OUTPUTS: { string - safe error message }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: errorMessageFromUnknown
function errorMessageFromUnknown(error: unknown): string {
  return error instanceof Error ? error.message : 'Request failed';
}

// START_CONTRACT: UsersPage
//   PURPOSE: Render the users list and create-user GraphQL mutation flow through UI-kit components.
//   INPUTS: none.
//   OUTPUTS: { JSX.Element - users route with loading, error, empty, list, and create states }
//   SIDE_EFFECTS: Sends createUser mutation and invalidates admin-users query on successful user creation.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
// END_CONTRACT: UsersPage
export default function UsersPage() {
  const queryClient = useQueryClient();
  const [form, setForm] = useState<FormState>(initialFormState);
  const [error, setError] = useState<string | null>(null);

  const usersQuery = useQuery({
    queryKey: ['admin-users'],
    queryFn: () => graphqlClient.request<GetUsersQuery>(getUsersQueryDocument, { first: 20 }),
  });

  const mutation = useMutation({
    mutationFn: (input: FormState) =>
      graphqlClient.request<CreateUserMutation>(createUserMutationDocument, { input }),
    onError: (mutationError) => setError(errorMessageFromUnknown(mutationError)),
    onSuccess: async (response) => {
      const result = response.createUser;

      // START_BLOCK_CREATE_USER_RESULT
      if ('user' in result) {
        setForm(initialFormState);
        setError(null);
        await queryClient.invalidateQueries({ queryKey: ['admin-users'] });
        return;
      }

      if ('field' in result) {
        setError(`${result.field}: ${result.message}`);
        return;
      }

      setError(result.message);
      // END_BLOCK_CREATE_USER_RESULT
    },
  });

  // START_CONTRACT: updateField
  //   PURPOSE: Update one create-user form field without mutating the other fields.
  //   INPUTS: { field: keyof FormState - field to change, value: string - next field value }
  //   OUTPUTS: none.
  //   SIDE_EFFECTS: Updates local React state.
  //   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
  // END_CONTRACT: updateField
  function updateField(field: keyof FormState, value: string) {
    setForm((current) => ({ ...current, [field]: value }));
  }

  // START_CONTRACT: handleSubmit
  //   PURPOSE: Submit the current create-user form through the GraphQL mutation.
  //   INPUTS: { event: FormEvent<HTMLFormElement> - form submit event }
  //   OUTPUTS: none.
  //   SIDE_EFFECTS: Prevents default form submission and starts the createUser mutation.
  //   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
  // END_CONTRACT: handleSubmit
  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    mutation.mutate(form);
  }

  const users = usersQuery.data?.users.edges || [];

  return (
    <AdminPageShell>
      <AdminPageHeader
        title="Users"
        description="Create and inspect reference users through the admin GraphQL API."
        actions={
          <Button asChild variant="outline">
            <Link to="/">Home</Link>
          </Button>
        }
      />

      <AdminToolbar>
        <p className="text-sm text-muted-foreground">Showing the latest 20 users.</p>
        <Button asChild variant="outline">
          <Link to="/ui-kit">UI kit</Link>
        </Button>
      </AdminToolbar>

      <AdminSection title="Create user" description="Submit a GraphQL createUser mutation.">
        <form className="grid gap-4 md:grid-cols-[1fr_1fr_1fr_auto]" onSubmit={handleSubmit}>
          <div className="space-y-2">
            <Label htmlFor="user-name">Name</Label>
            <Input
              id="user-name"
              onChange={(event) => updateField('name', event.target.value)}
              placeholder="Name"
              value={form.name}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="user-email">Email</Label>
            <Input
              id="user-email"
              onChange={(event) => updateField('email', event.target.value)}
              placeholder="Email"
              type="email"
              value={form.email}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="user-password">Password</Label>
            <Input
              id="user-password"
              onChange={(event) => updateField('password', event.target.value)}
              placeholder="Password"
              type="password"
              value={form.password}
            />
          </div>
          <div className="flex items-end">
            <Button disabled={mutation.isPending} type="submit">
              {mutation.isPending ? 'Creating...' : 'Create'}
            </Button>
          </div>
        </form>
      </AdminSection>

      {error ? (
        <Alert variant="destructive">
          <AlertTitle>Request failed</AlertTitle>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      ) : null}

      {usersQuery.isError ? (
        <Alert variant="destructive">
          <AlertTitle>Failed to load users.</AlertTitle>
          <AlertDescription>Refresh the page after the GraphQL API is available.</AlertDescription>
        </Alert>
      ) : null}

      {/* START_BLOCK_USERS_LIST_STATES */}
      <AdminSection
        title="Directory"
        description={
          usersQuery.data ? `Total: ${usersQuery.data.users.totalCount}` : 'Loading users.'
        }
      >
        {usersQuery.isLoading ? (
          <div className="space-y-2">
            <Skeleton className="h-10 w-full" />
            <Skeleton className="h-10 w-full" />
          </div>
        ) : null}

        {usersQuery.data && users.length === 0 ? (
          <AdminEmptyState title="No users yet" description="No users yet. Create one above." />
        ) : null}

        {users.length > 0 ? (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Email</TableHead>
                <TableHead className="text-right">Details</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {users.map(({ node }) => (
                <TableRow key={node.id}>
                  <TableCell>
                    <Button asChild className="h-auto p-0" variant="link">
                      <Link to={`/users/${node.id}`}>{node.name}</Link>
                    </Button>
                  </TableCell>
                  <TableCell>{node.email}</TableCell>
                  <TableCell className="text-right">
                    <Button asChild size="sm" variant="outline">
                      <Link to={`/users/${node.id}`}>Open</Link>
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        ) : null}
      </AdminSection>
      {/* END_BLOCK_USERS_LIST_STATES */}
    </AdminPageShell>
  );
}
```

- [ ] **Step 4: Update users Playwright flow for table markup**

In `apps/web-admin/e2e/users-flow.spec.ts`, replace the list-item locator in the create/list/detail test:

```ts
const createdUser = page.getByRole('listitem').filter({ hasText: email });
const userLink = createdUser.getByRole('link', { name: 'Browser User' });
await expect(userLink).toBeVisible();
await expect(createdUser.getByText(email)).toBeVisible();
```

with:

```ts
const createdRow = page.getByRole('row').filter({ hasText: email });
const userLink = createdRow.getByRole('link', { name: 'Browser User' });
await expect(userLink).toBeVisible();
await expect(createdRow.getByRole('cell', { name: email })).toBeVisible();
```

Do not run Playwright during active iteration unless needed for debugging; Task 10 runs the required `bunx nx run web-admin:e2e` gate.

- [ ] **Step 5: Run users checks**

Run:

```bash
bunx nx test web-admin -- src/pages/users-page.test.tsx
bunx nx run web-admin:typecheck
```

Expected: PASS.

- [ ] **Step 6: Stage users migration**

Run:

```bash
git add apps/web-admin/src/pages/users-page.tsx apps/web-admin/src/pages/users-page.test.tsx apps/web-admin/e2e/users-flow.spec.ts
```

Expected: users files are staged or ready for Task 9. Do not commit until GRACE docs and evidence are synchronized.

## Task 7: Migrate User Detail Route To UI Kit

**Files:**

- Modify: `apps/web-admin/src/pages/user-detail-page.tsx`
- Modify: `apps/web-admin/src/pages/user-detail-page.test.tsx`

- [ ] **Step 1: Update detail tests for UI-kit states**

Add these assertions to `apps/web-admin/src/pages/user-detail-page.test.tsx`:

```tsx
expect(screen.getByRole('link', { name: 'Back to users' })).toHaveAttribute('href', '/users');
expect(screen.getByText('Email')).toBeInTheDocument();
```

In the not-found test, add:

```tsx
expect(screen.getByText('The requested user does not exist.')).toBeInTheDocument();
```

- [ ] **Step 2: Run detail tests and verify they fail against old markup**

Run:

```bash
bunx nx test web-admin -- src/pages/user-detail-page.test.tsx
```

Expected: FAIL because the old not-found copy does not include the UI-kit empty-state description.

- [ ] **Step 3: Replace detail route implementation**

Replace `apps/web-admin/src/pages/user-detail-page.tsx` with:

```tsx
// FILE: apps/web-admin/src/pages/user-detail-page.tsx
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Render the web-admin user detail route.
//   SCOPE: Loads one user by route id and displays loading, error, not-found, and detail states through UI-kit components; excludes list and mutation behavior.
//   DEPENDS: @tanstack/react-query, react-router, apps/web-admin/src/entities/user/api/user.graphql, apps/web-admin/src/shared/api/graphql-client.ts, generated GraphQL types, apps/web-admin/src/shared/ui.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - User detail route backed by the GetUser GraphQL document and UI-kit states.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Migrated user detail route visuals to the web-admin UI kit.
// END_CHANGE_SUMMARY

import { useQuery } from '@tanstack/react-query';
import getUserQueryDocument from '@entities/user/api/user.graphql?raw';
import { graphqlClient } from '@shared/api/graphql-client';
import type { GetUserQuery } from '@shared/api/generated/types';
import { Link, useParams } from 'react-router';
import {
  AdminEmptyState,
  AdminPageHeader,
  AdminPageShell,
  AdminSection,
  Alert,
  AlertDescription,
  AlertTitle,
  Badge,
  Button,
  Skeleton,
} from '@shared/ui';

// START_CONTRACT: UserDetailPage
//   PURPOSE: Render one user loaded by route id with UI-kit loading, error, not-found, and detail states.
//   INPUTS: none.
//   OUTPUTS: { JSX.Element - user detail route for the current route id }
//   SIDE_EFFECTS: Sends GetUser GraphQL query through React Query.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
// END_CONTRACT: UserDetailPage
export default function UserDetailPage() {
  const { id } = useParams<{ id: string }>();
  const userQuery = useQuery({
    enabled: Boolean(id),
    queryKey: ['admin-user', id],
    queryFn: () => graphqlClient.request<GetUserQuery>(getUserQueryDocument, { id }),
  });

  const user = userQuery.data?.user || null;

  // START_BLOCK_USER_DETAIL_STATES
  if (userQuery.isLoading) {
    return (
      <AdminPageShell>
        <Skeleton className="h-9 w-64" />
        <Skeleton className="h-48 w-full" />
      </AdminPageShell>
    );
  }

  if (userQuery.isError) {
    return (
      <AdminPageShell>
        <Alert variant="destructive">
          <AlertTitle>Failed to load user.</AlertTitle>
          <AlertDescription>Refresh the page after the GraphQL API is available.</AlertDescription>
        </Alert>
        <Button asChild variant="outline">
          <Link to="/users">Back to users</Link>
        </Button>
      </AdminPageShell>
    );
  }

  if (!user) {
    return (
      <AdminPageShell>
        <AdminEmptyState
          title="User not found"
          description="The requested user does not exist."
          action={
            <Button asChild variant="outline">
              <Link to="/users">Back to users</Link>
            </Button>
          }
        />
      </AdminPageShell>
    );
  }

  return (
    <AdminPageShell>
      <AdminPageHeader
        title={user.name}
        description="Reference user loaded through the admin GraphQL API."
        actions={
          <Button asChild variant="outline">
            <Link to="/users">Back to users</Link>
          </Button>
        }
      />

      <AdminSection title="Profile" description="Stable user fields from GraphQL.">
        <dl className="grid gap-4 text-sm sm:grid-cols-2">
          <div className="space-y-1">
            <dt className="font-medium text-muted-foreground">Email</dt>
            <dd>{user.email}</dd>
          </div>
          <div className="space-y-1">
            <dt className="font-medium text-muted-foreground">ID</dt>
            <dd className="break-all">{user.id}</dd>
          </div>
          <div className="space-y-1">
            <dt className="font-medium text-muted-foreground">Created</dt>
            <dd>{new Date(user.createdAt).toLocaleString()}</dd>
          </div>
          <div className="space-y-1">
            <dt className="font-medium text-muted-foreground">Updated</dt>
            <dd>{new Date(user.updatedAt).toLocaleString()}</dd>
          </div>
        </dl>
        <div className="mt-4">
          <Badge variant="secondary">GraphQL user</Badge>
        </div>
      </AdminSection>
    </AdminPageShell>
  );
  // END_BLOCK_USER_DETAIL_STATES
}
```

- [ ] **Step 4: Run detail checks**

Run:

```bash
bunx nx test web-admin -- src/pages/user-detail-page.test.tsx
bunx nx run web-admin:typecheck
```

Expected: PASS.

- [ ] **Step 5: Stage detail migration**

Run:

```bash
git add apps/web-admin/src/pages/user-detail-page.tsx apps/web-admin/src/pages/user-detail-page.test.tsx
```

Expected: detail files are staged or ready for Task 9. Do not commit until GRACE docs and evidence are synchronized.

## Task 8: Enforce Page Imports Through ESLint

**Files:**

- Modify: `apps/web-admin/.eslintrc.json`
- Temporary local-only file during this task: `apps/web-admin/src/pages/ui-boundary-lint.fixture.tsx`

- [ ] **Step 1: Create a temporary bad page import fixture**

Create `apps/web-admin/src/pages/ui-boundary-lint.fixture.tsx`:

```tsx
import { Button } from '@shared/ui/primitives/button';
import { Card } from '@/shared/ui';
import { AdminToolbar } from '@/shared/ui/layout/admin-toolbar';
import { cn } from '../shared/ui/lib/utils';
import { Search } from 'lucide-react/icons/search';

export default function UiBoundaryLintFixture() {
  return (
    <AdminToolbar className={cn('gap-2')}>
      <Card>
        <Button>Bad import</Button>
      </Card>
      <Search aria-hidden="true" />
    </AdminToolbar>
  );
}
```

- [ ] **Step 2: Run lint and verify the bad import is not blocked yet**

Run:

```bash
cd apps/web-admin
bunx eslint src/pages/ui-boundary-lint.fixture.tsx --ext .tsx
cd ../..
```

Expected before adding the rule: PASS or a non-boundary lint result. If it fails because the primitive file is missing, return to Task 2 before continuing.

- [ ] **Step 3: Add page-level restricted imports**

Update `apps/web-admin/.eslintrc.json` to include an `overrides` array while preserving existing `extends`, `ignorePatterns`, `plugins`, `settings`, and `boundaries/element-types` rules:

```json
{
  "extends": ["../../.eslintrc.json"],
  "ignorePatterns": ["src/shared/api/generated/**", "next-env.d.ts"],
  "plugins": ["boundaries"],
  "settings": {
    "boundaries/elements": [
      { "type": "app", "pattern": "src/app/*" },
      { "type": "pages", "pattern": "src/pages/*" },
      { "type": "widgets", "pattern": "src/widgets/*" },
      { "type": "features", "pattern": "src/features/*" },
      { "type": "entities", "pattern": "src/entities/*" },
      { "type": "shared", "pattern": "src/shared/*" }
    ],
    "boundaries/ignore": ["**/*.test.*", "**/__tests__/**"]
  },
  "rules": {
    "boundaries/element-types": [
      "error",
      {
        "default": "disallow",
        "rules": [
          { "from": "app", "allow": ["pages", "widgets", "features", "entities", "shared"] },
          { "from": "pages", "allow": ["widgets", "features", "entities", "shared"] },
          { "from": "widgets", "allow": ["features", "entities", "shared"] },
          { "from": "features", "allow": ["entities", "shared"] },
          { "from": "entities", "allow": ["shared"] },
          { "from": "shared", "allow": ["shared"] }
        ]
      }
    ]
  },
  "overrides": [
    {
      "files": ["src/pages/**/*.{ts,tsx}"],
      "rules": {
        "no-restricted-imports": [
          "error",
          {
            "paths": [
              {
                "name": "class-variance-authority",
                "message": "Admin pages must use @shared/ui instead of UI implementation helpers."
              },
              {
                "name": "clsx",
                "message": "Admin pages must use @shared/ui instead of composing design-system classes directly."
              },
              {
                "name": "tailwind-merge",
                "message": "Admin pages must use @shared/ui instead of composing design-system classes directly."
              },
              {
                "name": "lucide-react",
                "message": "Export approved icon wrappers from @shared/ui before using icons on admin pages."
              },
              {
                "name": "radix-ui",
                "message": "Admin pages must use @shared/ui instead of Radix primitives directly."
              },
              {
                "name": "@/shared/ui",
                "message": "Admin pages must import UI from @shared/ui only."
              },
              {
                "name": "@shared/ui/primitives",
                "message": "Admin pages must import UI from @shared/ui only."
              },
              {
                "name": "@shared/ui/layout",
                "message": "Admin pages must import UI from @shared/ui only."
              },
              {
                "name": "@shared/ui/lib",
                "message": "Admin pages must import UI from @shared/ui only."
              },
              {
                "name": "@/shared/ui/primitives",
                "message": "Admin pages must import UI from @shared/ui only."
              },
              {
                "name": "@/shared/ui/layout",
                "message": "Admin pages must import UI from @shared/ui only."
              },
              {
                "name": "@/shared/ui/lib",
                "message": "Admin pages must import UI from @shared/ui only."
              },
              {
                "name": "../shared/ui",
                "message": "Admin pages must import UI from @shared/ui only."
              },
              {
                "name": "../../shared/ui",
                "message": "Admin pages must import UI from @shared/ui only."
              },
              {
                "name": "../../../shared/ui",
                "message": "Admin pages must import UI from @shared/ui only."
              }
            ],
            "patterns": [
              {
                "group": [
                  "@radix-ui/*",
                  "lucide-react/*",
                  "radix-ui/*",
                  "@shared/ui/*",
                  "@/shared/ui/*",
                  "../shared/ui/*",
                  "../shared/ui/**",
                  "../../shared/ui/*",
                  "../../shared/ui/**",
                  "../../../shared/ui/*",
                  "../../../shared/ui/**",
                  "**/shared/ui",
                  "**/shared/ui/**"
                ],
                "message": "Admin pages must import UI from @shared/ui only."
              }
            ]
          }
        ]
      }
    }
  ]
}
```

- [ ] **Step 4: Run lint and verify the temporary fixture fails**

Run:

```bash
cd apps/web-admin
bunx eslint src/pages/ui-boundary-lint.fixture.tsx --ext .tsx
cd ../..
```

Expected: FAIL with `Admin pages must import UI from @shared/ui only.`

Record the exact failure output in `.tasks/web-admin-shadcn-ui-kit/verification.md` under `Boundary Evidence`.

- [ ] **Step 5: Remove the temporary fixture and run real lint**

Run:

```bash
rm apps/web-admin/src/pages/ui-boundary-lint.fixture.tsx
bunx nx lint web-admin
```

Expected: PASS. This proves migrated pages use the public `@shared/ui` entrypoint.

- [ ] **Step 6: Stage enforcement**

Run:

```bash
git add apps/web-admin/.eslintrc.json
```

Expected: ESLint config is staged or ready for Task 9. Do not commit until GRACE docs and evidence are synchronized.

## Task 9: Update GRACE Docs And Verification Evidence

**Files:**

- Modify: `AGENTS.md`
- Modify: `docs/requirements.xml`
- Modify: `docs/technology.xml`
- Modify: `docs/development-plan.xml`
- Modify: `docs/knowledge-graph.xml`
- Modify: `docs/verification-plan.xml`
- Modify: `docs/operational-packets.xml` when the grep command in this task finds stale references
- Create or modify: `.tasks/web-admin-shadcn-ui-kit/verification.md`

- [ ] **Step 1: Update verification evidence file**

Open `.tasks/web-admin-shadcn-ui-kit/verification.md` from Task 0 and confirm it still contains:

- `bunx nx run web-admin:e2e` in `Command Evidence`;
- the expected failing `ui-boundary-lint.fixture.tsx` row in `Boundary Evidence`;
- the `JSON Governance` section for `components.json` and `.eslintrc.json`;
- the coverage decision naming `apps/web-admin/vite.config.ts` as the `web-admin:test-coverage` source of truth.

Add any Task 8 boundary failure output and operational packet decisions that have already been collected.

- [ ] **Step 2: Add admin UI-kit rule to `AGENTS.md`**

Add this section near the testing and file-local contract guidance:

```markdown
## Web-admin UI Kit Rule

All admin pages under `apps/web-admin/src/pages/**` must build UI from the approved `@shared/ui` surface.

- Store shadcn-generated primitives under `apps/web-admin/src/shared/ui/primitives/**`.
- Store admin page compositions under `apps/web-admin/src/shared/ui/layout/**`.
- Import UI in pages from `@shared/ui` only.
- Do not import Radix primitives, shadcn implementation subpaths, UI-kit aliases such as `@/shared/ui` or `@shared/ui/*`, relative `shared/ui` paths, class composition helpers, or icon libraries directly from page files.
- Prove the rule with `bunx nx lint web-admin` after page or UI-kit changes.
```

- [ ] **Step 3: Update `docs/requirements.xml`**

Add a constraint under `<Constraints>`:

```xml
<constraint-18>web-admin admin pages under `apps/web-admin/src/pages/**` must use the approved bare `@shared/ui` UI-kit surface for UI primitives and admin page compositions; direct page imports from `@shared/ui/*`, `@/shared/ui`, `@/shared/ui/*`, relative `shared/ui` paths, Radix packages, shadcn implementation subpaths, class composition helpers, and icon libraries are disallowed by ESLint.</constraint-18>
```

Add a risk under `<Risks>`:

```xml
<risk-16>Admin UI consistency can drift if future pages bypass `@shared/ui`; web-admin lint must enforce the page import boundary.</risk-16>
```

- [ ] **Step 4: Update `docs/technology.xml`**

Add dependency entries under `<Dependencies>`:

```xml
<dep name="tailwindcss" version="^4.3.0" purpose="Utility CSS and shadcn theme token pipeline for web-admin" />
<dep name="@tailwindcss/vite" version="^4.3.0" purpose="Tailwind CSS integration for the Vite web-admin build" />
<dep name="radix-ui" version="^1.4.3" purpose="Accessible primitive implementation dependency for shadcn web-admin components" />
<dep name="class-variance-authority" version="^0.7.1" purpose="Variant class generation for shadcn primitives" />
<dep name="tailwind-merge" version="^3.6.0" purpose="Tailwind class conflict resolution for shadcn UI helpers" />
<dep name="clsx" version="^2.1.1" purpose="Conditional class composition for shadcn UI helpers" />
<dep name="lucide-react" version="^1.17.0" purpose="Icon dependency available only through approved UI-kit exports on admin pages" />
<dep name="tw-animate-css" version="^1.4.0" purpose="Animation utilities imported by the web-admin shadcn theme CSS" />
```

Add a module-level command under web-admin checks:

```xml
<command>bunx nx lint web-admin</command>
```

- [ ] **Step 5: Update `docs/development-plan.xml`**

Inside `M-WEB-ADMIN`, update the purpose to mention the UI kit:

```xml
<purpose>Provide the Vite React Router web-admin SPA, generated GraphQL client, and approved shadcn-based UI-kit surface for the admin user flow.</purpose>
```

Add outputs:

```xml
<param name="ui-kit" type="@shared/ui primitives, admin layout compositions, and /ui-kit reference page" />
```

Add interface exports:

```xml
<export-ui-kit PURPOSE="Expose approved shadcn primitives and admin layout compositions through @shared/ui." />
<export-ui-kit-reference PURPOSE="Expose /ui-kit as the visible component reference route for future admin pages." />
```

Add target sources:

```xml
<source>apps/web-admin/src/shared/ui</source>
<source>apps/web-admin/src/pages/ui-kit-page.tsx</source>
```

- [ ] **Step 6: Update `docs/knowledge-graph.xml`**

Inside `M-WEB-ADMIN`, add paths:

```xml
<path>apps/web-admin/src/shared/ui</path>
<path>apps/web-admin/src/pages/ui-kit-page.tsx</path>
```

Add annotations:

```xml
<export-UiKit PURPOSE="Approved @shared/ui primitive and layout surface for admin pages." />
<export-UiKitReferenceRoute PURPOSE="Visible /ui-kit route demonstrating approved admin components." />
<constraint-PageUiImports PURPOSE="Admin pages import UI from @shared/ui only, enforced by web-admin ESLint." />
```

- [ ] **Step 7: Update `docs/verification-plan.xml`**

In `V-M-WEB-ADMIN`, add test files:

```xml
<file>apps/web-admin/components.json</file>
<file>apps/web-admin/src/styles.css</file>
<file>apps/web-admin/src/shared/ui</file>
<file>apps/web-admin/src/pages/ui-kit-page.tsx</file>
<file>apps/web-admin/src/pages/ui-kit-page.test.tsx</file>
<file>apps/web-admin/e2e/users-flow.spec.ts</file>
<file>apps/web-admin/.eslintrc.json</file>
```

Add required trace assertions:

```xml
<assertion-2>Admin pages under `apps/web-admin/src/pages/**` must import UI through bare `@shared/ui`, not `@shared/ui/*`, `@/shared/ui`, `@/shared/ui/*`, relative `shared/ui` paths, Radix packages, class-composition helpers, or icon libraries.</assertion-2>
<assertion-3>`bunx nx run web-admin:test-coverage` must cover the UI-kit source surface without broad UI-kit coverage exclusions.</assertion-3>
<assertion-4>`bunx nx run web-admin:e2e` must cover the changed admin route flow including the visible `/ui-kit` route or record a blocker instead of READY.</assertion-4>
```

Update scenario 1:

```xml
<scenario-1 kind="success">Vite web-admin app renders admin user flow and the /ui-kit reference route with generated GraphQL types, React Router routes, and approved @shared/ui components.</scenario-1>
```

- [ ] **Step 8: Check `docs/operational-packets.xml` for stale text**

Run:

```bash
rg -n "web-admin|M-WEB-ADMIN|V-M-WEB-ADMIN|UI kit|ui-kit|shared/ui" docs/operational-packets.xml
```

Expected:

- If output mentions web-admin implementation or verification surfaces, update the relevant packet text to include `@shared/ui`, `/ui-kit`, `bunx nx lint web-admin`, `bunx nx run web-admin:test-coverage`, and `bunx nx run web-admin:e2e`.
- If output is empty or only generic, leave `docs/operational-packets.xml` unchanged and record `No operational packet update needed` in `.tasks/web-admin-shadcn-ui-kit/verification.md`.

- [ ] **Step 9: Validate docs**

Run:

```bash
xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
grace lint --path .
```

Expected: PASS. If `grace lint --path .` reports pre-existing unrelated issues, record the exact output and do not mix them with this wave.

- [ ] **Step 10: Commit implementation, docs, and evidence**

Run:

```bash
git status --short -- apps/web-admin bun.lock AGENTS.md docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/knowledge-graph.xml docs/verification-plan.xml docs/operational-packets.xml .tasks/web-admin-shadcn-ui-kit
git add apps/web-admin/package.json bun.lock apps/web-admin/components.json apps/web-admin/vite.config.ts apps/web-admin/src/styles.css apps/web-admin/src/App.tsx apps/web-admin/src/App.test.tsx apps/web-admin/src/shared/ui apps/web-admin/src/pages/home.tsx apps/web-admin/src/pages/ui-kit-page.tsx apps/web-admin/src/pages/ui-kit-page.test.tsx apps/web-admin/src/pages/users-page.tsx apps/web-admin/src/pages/users-page.test.tsx apps/web-admin/src/pages/user-detail-page.tsx apps/web-admin/src/pages/user-detail-page.test.tsx apps/web-admin/e2e/users-flow.spec.ts apps/web-admin/.eslintrc.json AGENTS.md docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/knowledge-graph.xml docs/verification-plan.xml docs/operational-packets.xml .tasks/web-admin-shadcn-ui-kit/verification.md
git commit -m "feat(web): add web-admin shadcn ui kit"
```

Expected: the first commit includes implementation, docs, and evidence together. If `git status --short -- ...` shows unrelated hunks in shared files such as `AGENTS.md`, do not run the broad `git add`; move to a clean worktree or pause for user direction.

## Task 10: Final Focused Verification And Handoff

**Files:**

- Modify: `.tasks/web-admin-shadcn-ui-kit/verification.md`

- [ ] **Step 1: Run focused web-admin tests**

Run:

```bash
bunx nx test web-admin
```

Expected: PASS. Record the result in `.tasks/web-admin-shadcn-ui-kit/verification.md`.

- [ ] **Step 2: Run required web-admin coverage**

Run:

```bash
bunx nx run web-admin:test-coverage
```

Expected: PASS with 100 percent statements, branches, functions, and lines. Record the result in `.tasks/web-admin-shadcn-ui-kit/verification.md`.

- [ ] **Step 3: Run typecheck, lint, and build**

Run:

```bash
bunx nx run web-admin:typecheck
bunx nx lint web-admin
bunx nx build web-admin
```

Expected: all PASS. Record results in `.tasks/web-admin-shadcn-ui-kit/verification.md`.

- [ ] **Step 4: Run web-admin e2e**

Run:

```bash
bunx nx run web-admin:e2e
```

Expected: PASS, including the changed admin route flow and visible `/ui-kit` route. Record the result in `.tasks/web-admin-shadcn-ui-kit/verification.md`. If the e2e environment is unavailable, record the exact blocker and do not mark the wave READY.

- [ ] **Step 5: Run conditional root coverage policy gate**

Run this only if `apps/web-admin/vite.config.ts` coverage excludes, `tools/coverage/coverage.config.json`, or coverage policy docs changed to exclude any UI-kit source:

```bash
bun run test:coverage
```

Expected: PASS with exact allowlist/replacement-gate evidence recorded in `.tasks/web-admin-shadcn-ui-kit/verification.md`. If no coverage exclusion or coverage policy file changed, record `No coverage-policy gate needed; web-admin:test-coverage covers UI-kit source`.

- [ ] **Step 6: Run GRACE checks**

Run:

```bash
xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
grace lint --path .
```

Expected: both PASS. Record results in `.tasks/web-admin-shadcn-ui-kit/verification.md`.

- [ ] **Step 7: Inspect page import boundary evidence**

Run:

```bash
rg -n "from ['\"](@shared/ui/|@/shared/ui(/|['\"])|(\\.\\./)+shared/ui(/|['\"])|lucide-react(/|['\"])|radix-ui(/|['\"])|@radix-ui/|class-variance-authority|clsx|tailwind-merge)" apps/web-admin/src/pages
```

Expected: no matches. Record `No direct page UI bypass imports found` in `.tasks/web-admin-shadcn-ui-kit/verification.md`.

- [ ] **Step 8: Update final status**

Change `.tasks/web-admin-shadcn-ui-kit/verification.md` final status to:

```markdown
READY - focused web-admin tests, coverage, typecheck, lint, build, e2e, import-boundary scan, XML validation, and GRACE lint passed.
```

- [ ] **Step 9: Commit verification evidence**

Run:

```bash
git add .tasks/web-admin-shadcn-ui-kit/verification.md
git commit -m "docs(docs): record web-admin ui kit verification"
```

- [ ] **Step 10: Report handoff**

Report:

- changed UI-kit entrypoint: `apps/web-admin/src/shared/ui/index.ts`;
- demo route: `/ui-kit`;
- migrated routes: `/`, `/users`, `/users/:id`;
- enforcement: `apps/web-admin/.eslintrc.json` plus `AGENTS.md`;
- verification commands and PASS results from `.tasks/web-admin-shadcn-ui-kit/verification.md`;
- any pre-existing dirty worktree entries left untouched.

## Self-Review

### Spec Coverage

- Local `apps/web-admin/src/shared/ui` UI-kit layer: Tasks 1-3.
- Visible `/ui-kit` page linked from home: Tasks 4-5.
- Broad showcase: Task 4.
- Existing route migration: Tasks 5-7.
- Browser coverage for migrated users flow and `/ui-kit`: Tasks 4, 6, and 10.
- Docs plus ESLint enforcement: Tasks 8-9.
- Required `web-admin:test-coverage`: Tasks 9-10.
- shadcn placement under `shared/ui`: Tasks 1-2.
- No broad UI-kit coverage exclusions: Tasks 9-10.
- GRACE docs and file-local markup: comment-capable governed files include headers and useful anchors; commentless JSON configs are governed through Task 9 XML/docs entries and `.tasks/web-admin-shadcn-ui-kit/verification.md`.

### Placeholder Scan

The plan contains no placeholder tokens, no deferred implementation language, no cross-task shorthand, and no generic test-writing steps. Steps that change code include concrete code blocks or concrete CLI commands with expected files.

### Type Consistency

The plan consistently uses:

- `AdminPageShell`, `AdminPageHeader`, `AdminToolbar`, `AdminSection`, `AdminEmptyState`;
- `UiKitPage`;
- public UI imports from `@shared/ui`;
- current existing GraphQL types `CreateUserMutation`, `GetUsersQuery`, and `GetUserQuery`;
- current Nx commands `bunx nx test web-admin`, `bunx nx run web-admin:test-coverage`, `bunx nx run web-admin:typecheck`, `bunx nx lint web-admin`, `bunx nx build web-admin`, and `bunx nx run web-admin:e2e`.
