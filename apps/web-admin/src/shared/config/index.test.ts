// FILE: apps/web-admin/src/shared/config/index.test.ts
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Verify the web-admin browser configuration contract.
//   SCOPE: Covers Vite environment defaults, overrides, URL normalization, and app-level re-export identity; excludes runtime network calls.
//   DEPENDS: apps/web-admin/src/shared/config/index.ts, apps/web-admin/src/app/config.ts, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   web-admin appConfig tests - Prove the Vite env contract and re-export behavior.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Switched expectations from Next public env to Vite browser env.
// END_CHANGE_SUMMARY

import { afterEach, describe, expect, it, vi } from 'vitest';

describe('web-admin appConfig', () => {
  afterEach(() => {
    vi.unstubAllEnvs();
    vi.resetModules();
  });

  it('defaults to the local GraphQL endpoint and app name', async () => {
    vi.stubEnv('VITE_GRAPHQL_API_URL', '');
    vi.stubEnv('VITE_APP_NAME', '');

    const { appConfig } = await import('./index');

    expect(appConfig).toEqual({
      apiUrl: 'http://localhost:8090/graphql',
      appName: 'MonorepoApp',
    });
  });

  it('uses Vite environment overrides and trims trailing slashes', async () => {
    vi.stubEnv('VITE_GRAPHQL_API_URL', ' https://api.test/graphql/// ');
    vi.stubEnv('VITE_APP_NAME', 'TemplateAdmin');

    const shared = await import('./index');
    const app = await import('../../app/config');

    expect(shared.appConfig.apiUrl).toBe('https://api.test/graphql');
    expect(shared.appConfig.appName).toBe('TemplateAdmin');
    expect(app.appConfig).toBe(shared.appConfig);
  });
});
