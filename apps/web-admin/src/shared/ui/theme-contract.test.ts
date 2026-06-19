// FILE: apps/web-admin/src/shared/ui/theme-contract.test.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify the web-admin shadcn and Tailwind theme contract.
//   SCOPE: Covers components.json, Tailwind v4 CSS setup, and semantic token mappings; excludes visual pixel assertions.
//   DEPENDS: node:fs, node:path, node:url, vitest, apps/web-admin/components.json, apps/web-admin/src/styles.css.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   web-admin theme contract tests - Prove shadcn config and CSS tokens remain aligned.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added web-admin token contract coverage.
// END_CHANGE_SUMMARY

import { readFileSync } from 'node:fs';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { describe, expect, it } from 'vitest';

const currentDir = dirname(fileURLToPath(import.meta.url));
const appDir = resolve(currentDir, '../../..');
const componentsJson = JSON.parse(readFileSync(resolve(appDir, 'components.json'), 'utf8')) as {
  iconLibrary: string;
  rsc: boolean;
  style: string;
  tailwind: {
    baseColor: string;
    config: string;
    css: string;
    cssVariables: boolean;
  };
};
const stylesCss = readFileSync(resolve(appDir, 'src/styles.css'), 'utf8');

describe('web-admin shadcn theme contract', () => {
  it('keeps the CLI config aligned with the requested style and token mode', () => {
    expect(componentsJson.style).toBe('radix-rhea');
    expect(componentsJson.rsc).toBe(false);
    expect(componentsJson.iconLibrary).toBe('lucide');
    expect(componentsJson.tailwind).toMatchObject({
      baseColor: 'zinc',
      config: '',
      css: 'src/styles.css',
      cssVariables: true,
    });
  });

  it('exposes shadcn semantic tokens through Tailwind v4 theme variables', () => {
    expect(stylesCss).toContain("@import 'tailwindcss';");
    expect(stylesCss).toContain("@import 'tw-animate-css';");
    expect(stylesCss).toContain('@custom-variant dark');
    expect(stylesCss).toContain('@theme inline');

    for (const token of [
      'background',
      'foreground',
      'card',
      'card-foreground',
      'primary',
      'primary-foreground',
      'secondary',
      'muted',
      'muted-foreground',
      'accent',
      'accent-foreground',
      'destructive',
      'border',
      'input',
      'ring',
    ]) {
      expect(stylesCss).toContain(`--color-${token}: var(--${token});`);
    }
  });

  it('exposes sidebar semantic tokens for shadcn sidebar primitives', () => {
    for (const token of [
      'sidebar',
      'sidebar-foreground',
      'sidebar-primary',
      'sidebar-primary-foreground',
      'sidebar-accent',
      'sidebar-accent-foreground',
      'sidebar-border',
      'sidebar-ring',
    ]) {
      expect(stylesCss).toContain(`--color-${token}: var(--${token});`);
      expect(stylesCss).toContain(`--${token}:`);
    }

    expect(stylesCss).toContain('.dark {');
  });
});
