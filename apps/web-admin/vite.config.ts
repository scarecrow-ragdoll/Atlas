// FILE: apps/web-admin/vite.config.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Configure Vite and Vitest for the web-admin SPA.
//   SCOPE: Defines React plugin wiring, module aliases, jsdom tests, and 100 percent coverage settings; excludes Playwright e2e server orchestration.
//   DEPENDS: @tailwindcss/vite, @vitejs/plugin-react, vitest, apps/web-admin/src.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / V-M-COVERAGE-GATE.
//   ROLE: CONFIG
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - Vite/Vitest config for app build, test, aliases, and coverage.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Added Tailwind CSS Vite plugin for shadcn UI primitives.
// END_CHANGE_SUMMARY

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
