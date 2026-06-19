// FILE: apps/web/src/shared/config.test.ts
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Verify the public web REST configuration contract.
//   SCOPE: Covers browser same-origin REST base and server/runtime WEB_API_BASE_URL resolution; excludes route proxy behavior.
//   DEPENDS: apps/web/src/shared/config.ts, vitest.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   web config tests - Prove browser and server REST base URL behavior.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Switched expectations from Vite public env to Next runtime REST config.
// END_CHANGE_SUMMARY

import { afterEach, describe, expect, it, vi } from 'vitest';

afterEach(() => {
  vi.unstubAllGlobals();
  vi.unstubAllEnvs();
  vi.resetModules();
});

describe('web config', () => {
  it('defaults browser REST requests to the same-origin Next proxy', async () => {
    vi.stubEnv('WEB_API_BASE_URL', 'https://api.example.test');

    const { appConfig } = await import('./config');

    expect(appConfig.apiBaseUrl).toBe('');
  });

  it('uses the runtime server API base URL when no browser window exists', async () => {
    vi.stubGlobal('window', undefined);
    vi.stubEnv('WEB_API_BASE_URL', 'https://api.example.test');

    const { appConfig } = await import('./config');

    expect(appConfig.apiBaseUrl).toBe('https://api.example.test');
  });

  it('trims whitespace and trailing slashes from the runtime server API base URL', async () => {
    vi.stubGlobal('window', undefined);
    vi.stubEnv('WEB_API_BASE_URL', ' https://api.example.test/// ');

    const { appConfig } = await import('./config');

    expect(appConfig.apiBaseUrl).toBe('https://api.example.test');
  });

  it('defaults server REST requests to the local API base URL', async () => {
    vi.stubGlobal('window', undefined);
    vi.stubEnv('WEB_API_BASE_URL', '');

    const { appConfig } = await import('./config');

    expect(appConfig.apiBaseUrl).toBe('http://localhost:8090');
  });
});
