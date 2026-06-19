// FILE: apps/web/src/shared/ui/theme-contract.test.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify the public web shadcn and Tailwind theme contract.
//   SCOPE: Covers components.json, Tailwind v4 CSS setup, semantic token mappings, and avoidance of legacy app-local color variables; excludes visual pixel assertions.
//   DEPENDS: node:fs, node:path, node:url, vitest, apps/web/components.json, apps/web/app/globals.css.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   public web theme contract tests - Prove shadcn config and CSS tokens remain aligned.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added public web token contract coverage.
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
const globalsCss = readFileSync(resolve(appDir, 'app/globals.css'), 'utf8');

describe('public web shadcn theme contract', () => {
  it('keeps the CLI config aligned with the requested style and token mode', () => {
    expect(componentsJson.style).toBe('radix-rhea');
    expect(componentsJson.rsc).toBe(true);
    expect(componentsJson.iconLibrary).toBe('lucide');
    expect(componentsJson.tailwind).toMatchObject({
      baseColor: 'zinc',
      config: '',
      css: 'app/globals.css',
      cssVariables: true,
    });
  });

  it('exposes shadcn semantic tokens through Tailwind v4 theme variables', () => {
    expect(globalsCss).toContain("@import 'tailwindcss';");
    expect(globalsCss).toContain("@import 'tw-animate-css';");
    expect(globalsCss).toContain('@custom-variant dark');
    expect(globalsCss).toContain('@theme inline');

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
      expect(globalsCss).toContain(`--color-${token}: var(--${token});`);
    }
  });

  it('uses semantic tokens instead of legacy public-web color variables in component CSS', () => {
    expect(globalsCss).not.toMatch(/--(?:accent-dark|danger|panel|panel-soft|shadow):/);
    expect(globalsCss).not.toMatch(/var\(--(?:accent-dark|danger|panel|panel-soft|shadow)\)/);
  });
});
