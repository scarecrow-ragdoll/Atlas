// FILE: apps/web/vitest.config.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Configure Vitest for the public Next web app.
//   SCOPE: Defines jsdom setup, app/src test discovery, aliases, and 100 percent coverage settings; excludes Playwright e2e orchestration.
//   DEPENDS: vitest, apps/web/app, apps/web/src.
//   LINKS: M-WEB / V-M-WEB / V-M-COVERAGE-GATE.
//   ROLE: CONFIG
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - Vitest config for Next app, route proxy, shared REST client, and coverage tests.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Next-aware Vitest configuration for public web.
// END_CHANGE_SUMMARY

import { resolve } from 'path';
import { defineConfig } from 'vitest/config';

export default defineConfig({
  esbuild: {
    jsx: 'automatic',
    jsxImportSource: 'react',
  },
  test: {
    coverage: {
      exclude: ['app/**/*.test.{ts,tsx}', 'src/**/*.test.{ts,tsx}', 'next-env.d.ts'],
      include: ['app/**/*.{ts,tsx}', 'src/**/*.{ts,tsx}'],
      provider: 'v8',
      reporter: ['text', 'json-summary'],
      reportsDirectory: '../../dist/coverage/web',
      thresholds: {
        branches: 100,
        functions: 100,
        lines: 100,
        statements: 100,
      },
    },
    environment: 'jsdom',
    include: ['app/**/*.test.{ts,tsx}', 'src/**/*.test.{ts,tsx}'],
    setupFiles: ['./vitest.setup.ts'],
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
      '@app': resolve(__dirname, './src/app'),
      '@shared': resolve(__dirname, './src/shared'),
    },
  },
});
