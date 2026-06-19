// FILE: tools/ci/src/dokploy.test.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verifies Dokploy compose env migration and deployment HTTP sequencing.
//   SCOPE: Covers injected and global fetch behavior; excludes live Dokploy network calls.
//   DEPENDS: vitest, tools/ci/src/dokploy.ts.
//   LINKS: M-CI-CD / V-M-CI-CD.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   deployDokployCompose - Guards compose env merge, update, and deploy ordering.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added coverage for WEB_API_BASE_URL compose env append and update behavior.
// END_CHANGE_SUMMARY

import { afterEach, describe, expect, it, vi } from 'vitest';
import { deployDokployCompose, type FetchLike } from './dokploy';

function okResponse(body: unknown = { ok: true }) {
  return Promise.resolve({
    ok: true,
    status: 200,
    text: () => Promise.resolve(JSON.stringify(body)),
  });
}

function errorResponse(status: number, body: string) {
  return Promise.resolve({
    ok: false,
    status,
    text: () => Promise.resolve(body),
  });
}

describe('deployDokployCompose', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('updates compose env before triggering deploy', async () => {
    const calls: Array<{
      url: string;
      method: string;
      headers: Record<string, string>;
      body: unknown;
    }> = [];
    const fetchImpl: FetchLike = async (url, init) => {
      calls.push({
        url,
        method: init.method,
        headers: init.headers,
        body: init.body ? JSON.parse(init.body) : null,
      });
      if (url.includes('/api/compose.one')) {
        return okResponse({
          env: 'POSTGRES_HOST=db\n# runtime-only\nIMAGE_TAG=old\nAUTH_JWT_SECRET=keep-me\n\nAPI_IMAGE=old-api',
        });
      }
      return okResponse();
    };

    await deployDokployCompose({
      baseUrl: 'https://dokploy.example.com/',
      apiKey: 'secret',
      composeId: 'cmp_123',
      imageEnv: {
        IMAGE_TAG: 'v1.2.3',
        API_IMAGE: 'registry/app/api:v1.2.3',
        WEB_IMAGE: 'registry/app/web:v1.2.3',
        BOT_IMAGE: 'registry/app/bot:v1.2.3',
        WEB_API_BASE_URL: 'https://api.example.com',
      },
      fetchImpl,
    });

    expect(calls).toHaveLength(3);
    expect(calls[0].url).toBe('https://dokploy.example.com/api/compose.one?composeId=cmp_123');
    expect(calls[0].method).toBe('GET');
    expect(calls[0].headers['x-api-key']).toBe('secret');
    expect(calls[0].body).toBeNull();
    expect(calls[1].url).toBe('https://dokploy.example.com/api/compose.update');
    expect(calls[1].body).toEqual({
      composeId: 'cmp_123',
      env: [
        'POSTGRES_HOST=db',
        '# runtime-only',
        'IMAGE_TAG=v1.2.3',
        'AUTH_JWT_SECRET=keep-me',
        'API_IMAGE=registry/app/api:v1.2.3',
        'WEB_IMAGE=registry/app/web:v1.2.3',
        'BOT_IMAGE=registry/app/bot:v1.2.3',
        'WEB_API_BASE_URL=https://api.example.com',
      ].join('\n'),
    });
    expect(calls[2].url).toBe('https://dokploy.example.com/api/compose.deploy');
    expect(calls[2].body).toEqual({ composeId: 'cmp_123' });
  });

  it('updates existing WEB_API_BASE_URL in compose env', async () => {
    const calls: Array<{ url: string; body: unknown }> = [];
    const fetchImpl: FetchLike = async (url, init) => {
      calls.push({
        url,
        body: init.body ? JSON.parse(init.body) : null,
      });
      if (url.includes('/api/compose.one')) {
        return okResponse({
          env: 'IMAGE_TAG=old\nWEB_API_BASE_URL=https://old-api.example.com\nWEB_IMAGE=old-web',
        });
      }
      return okResponse();
    };

    await deployDokployCompose({
      baseUrl: 'https://dokploy.example.com',
      apiKey: 'secret',
      composeId: 'cmp_123',
      imageEnv: {
        IMAGE_TAG: 'v1.2.3',
        WEB_IMAGE: 'registry/app/web:v1.2.3',
        WEB_API_BASE_URL: 'https://api.example.com',
      },
      fetchImpl,
    });

    expect(calls[1].body).toEqual({
      composeId: 'cmp_123',
      env: [
        'IMAGE_TAG=v1.2.3',
        'WEB_API_BASE_URL=https://api.example.com',
        'WEB_IMAGE=registry/app/web:v1.2.3',
      ].join('\n'),
    });
  });

  it('uses global fetch when no fetch dependency is injected', async () => {
    const fetchImpl = vi.fn(async (url: string) => {
      if (url.includes('/api/compose.one')) {
        return okResponse({});
      }
      return okResponse();
    });
    vi.stubGlobal('fetch', fetchImpl);

    await deployDokployCompose({
      baseUrl: 'https://dokploy.example.com',
      apiKey: 'secret',
      composeId: 'cmp with space',
      imageEnv: {
        IMAGE_TAG: 'v1.2.3',
      },
    });

    expect(fetchImpl.mock.calls[0][0]).toBe(
      'https://dokploy.example.com/api/compose.one?composeId=cmp%20with%20space',
    );
  });

  it('fails when Dokploy update returns an error and does not deploy', async () => {
    const calls: string[] = [];
    const fetchImpl: FetchLike = async (url) => {
      if (url.includes('/api/compose.one')) {
        calls.push('compose.one');
        return okResponse({ env: null });
      }
      if (url.includes('/api/compose.update')) {
        calls.push('compose.update');
        return errorResponse(500, 'boom');
      }
      calls.push('compose.deploy');
      return okResponse();
    };

    await expect(
      deployDokployCompose({
        baseUrl: 'https://dokploy.example.com',
        apiKey: 'secret',
        composeId: 'cmp_123',
        imageEnv: { IMAGE_TAG: 'v1.2.3' },
        fetchImpl,
      }),
    ).rejects.toThrow('Dokploy compose.update failed with status 500: boom');
    expect(calls).toEqual(['compose.one', 'compose.update']);
  });

  it('fails when Dokploy deploy returns an error after a successful update', async () => {
    const fetchImpl: FetchLike = async (url) => {
      if (url.includes('/api/compose.one')) {
        return okResponse({ env: '' });
      }
      if (url.includes('/api/compose.deploy')) {
        return errorResponse(503, 'unavailable');
      }
      return okResponse();
    };

    await expect(
      deployDokployCompose({
        baseUrl: 'https://dokploy.example.com',
        apiKey: 'secret',
        composeId: 'cmp_123',
        imageEnv: { IMAGE_TAG: 'v1.2.3' },
        fetchImpl,
      }),
    ).rejects.toThrow('Dokploy compose.deploy failed with status 503: unavailable');
  });
});
