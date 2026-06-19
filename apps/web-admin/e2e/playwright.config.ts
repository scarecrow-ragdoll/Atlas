// FILE: apps/web-admin/e2e/playwright.config.ts
// VERSION: 1.1.0
// START_MODULE_CONTRACT
//   PURPOSE: Configure Playwright e2e for the Vite web-admin GraphQL app.
//   SCOPE: Starts test dependencies, the Go API, and the Vite admin dev server with GraphQL env; excludes test scenario assertions.
//   DEPENDS: @playwright/test, apps/web-admin/e2e/preflight.mjs, apps/api, apps/web-admin.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: CONFIG
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - Playwright config for admin GraphQL e2e browser and API setup.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.1.0 - Aligned API bootstrap credentials with browser e2e login credentials.
// END_CHANGE_SUMMARY

import { defineConfig, devices } from '@playwright/test';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = dirname(fileURLToPath(import.meta.url));
const repoRoot = resolve(__dirname, '../../..');
const webRoot = resolve(__dirname, '..');
const apiPort = process.env.E2E_API_PORT ?? '18080';
const apiBaseURL = process.env.E2E_API_URL ?? `http://localhost:${apiPort}`;
const webPort = process.env.E2E_WEB_PORT ?? '13000';
const webBaseURL = process.env.E2E_WEB_URL ?? `http://localhost:${webPort}`;
const testPostgresHost = process.env.TEST_POSTGRES_HOST ?? 'localhost';
const testPostgresPort = process.env.TEST_POSTGRES_PORT ?? '17501';
const testPostgresUser = process.env.TEST_POSTGRES_USER ?? 'app';
const testPostgresPassword = process.env.TEST_POSTGRES_PASSWORD ?? 'secret';
const testPostgresDB = process.env.TEST_POSTGRES_DB ?? 'monorepo_test';
const testRedisHost = process.env.TEST_REDIS_HOST ?? 'localhost';
const testRedisPort = process.env.TEST_REDIS_PORT ?? '17502';
const testRedisPassword = process.env.TEST_REDIS_PASSWORD ?? '';
const e2eAdminEmail = process.env.E2E_ADMIN_EMAIL ?? 'e2e-admin@example.test';
const e2eAdminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'StrongPassword123!';
const e2eAdminName = process.env.E2E_ADMIN_NAME ?? 'E2E Admin';

export default defineConfig({
  testDir: __dirname,
  outputDir: resolve(repoRoot, 'dist/test-results/web-admin-e2e'),
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: [
    ['list'],
    [
      'html',
      { outputFolder: resolve(repoRoot, 'dist/playwright-report/web-admin'), open: 'never' },
    ],
  ],
  use: {
    baseURL: webBaseURL,
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  webServer: [
    {
      command: 'bun run e2e:preflight && cd ../api && go run ./cmd/server',
      url: `${apiBaseURL}/readyz`,
      reuseExistingServer: false,
      timeout: 120_000,
      cwd: webRoot,
      env: {
        SERVER_PORT: apiPort,
        SERVER_CORS_ORIGINS: webBaseURL,
        POSTGRES_HOST: testPostgresHost,
        POSTGRES_PORT: testPostgresPort,
        POSTGRES_USER: testPostgresUser,
        POSTGRES_PASSWORD: testPostgresPassword,
        POSTGRES_DB: testPostgresDB,
        POSTGRES_SSLMODE: 'disable',
        REDIS_HOST: testRedisHost,
        REDIS_PORT: testRedisPort,
        REDIS_PASSWORD: testRedisPassword,
        ADMIN_INITIAL_EMAIL: e2eAdminEmail,
        ADMIN_INITIAL_PASSWORD: e2eAdminPassword,
        ADMIN_INITIAL_NAME: e2eAdminName,
        ADMIN_ORIGINS: webBaseURL,
        ADMIN_SESSION_COOKIE_NAME: 'web_admin_session',
        ADMIN_SESSION_TTL: '168h',
        ADMIN_SESSION_COOKIE_SECURE: 'false',
        ADMIN_SESSION_SAME_SITE: 'Lax',
        ADMIN_SESSION_KEY_SECRET: 'e2e-session-key-secret',
      },
    },
    {
      command: `bun run dev -- --host 127.0.0.1 --port ${webPort}`,
      url: webBaseURL,
      reuseExistingServer: false,
      timeout: 120_000,
      cwd: webRoot,
      env: {
        VITE_GRAPHQL_API_URL: `${apiBaseURL}/graphql`,
      },
    },
  ],
});
