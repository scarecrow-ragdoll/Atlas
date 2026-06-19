// FILE: apps/web/app/api/users/route.test.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify the public Next REST route proxy handlers.
//   SCOPE: Covers root list/create and id-scoped get/update/delete proxying, body forwarding, encoded ids, and runtime WEB_API_BASE_URL defaults; excludes Go API behavior.
//   DEPENDS: apps/web/app/api/users/route.ts, apps/web/app/api/users/[id]/route.ts, vitest.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   route proxy tests - Prove Next handlers forward REST requests to the configured API base.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added red coverage for public Next REST route proxy behavior.
// END_CHANGE_SUMMARY

import { afterEach, describe, expect, it, vi } from 'vitest';
import { DELETE, GET as GET_BY_ID, PATCH } from './[id]/route';
import { GET, POST } from './route';

afterEach(() => {
  vi.restoreAllMocks();
  vi.unstubAllEnvs();
  vi.resetModules();
});

function jsonResponse(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), { status });
}

describe('public users route proxy', () => {
  it('proxies list requests through the default local API base URL', async () => {
    vi.stubEnv('WEB_API_BASE_URL', '');
    const fetchMock = vi
      .spyOn(globalThis, 'fetch')
      .mockResolvedValue(jsonResponse({ data: [], meta: { totalCount: 0 } }));

    const response = await GET();

    expect(fetchMock).toHaveBeenCalledWith(
      'http://localhost:8090/api/users',
      expect.objectContaining({ method: 'GET' }),
    );
    await expect(response.json()).resolves.toEqual({ data: [], meta: { totalCount: 0 } });
  });

  it('proxies create requests with the request body', async () => {
    vi.stubEnv('WEB_API_BASE_URL', 'https://api.example.test///');
    const body = { email: 'new@example.com', name: 'New', password: 'secret123' };
    const fetchMock = vi
      .spyOn(globalThis, 'fetch')
      .mockResolvedValue(jsonResponse({ data: { id: 'u1', ...body } }, 201));

    const response = await POST(
      new Request('http://next.test/api/users', {
        body: JSON.stringify(body),
        method: 'POST',
      }),
    );

    expect(fetchMock).toHaveBeenCalledWith(
      'https://api.example.test/api/users',
      expect.objectContaining({ body: JSON.stringify(body), method: 'POST' }),
    );
    expect(response.status).toBe(201);
  });

  it('defaults upstream response content type when it is absent', async () => {
    vi.stubEnv('WEB_API_BASE_URL', 'https://api.example.test');
    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      headers: new Headers(),
      status: 200,
      text: async () => JSON.stringify({ data: [], meta: { totalCount: 0 } }),
    } as Response);

    const response = await GET();

    expect(response.headers.get('Content-Type')).toBe('application/json');
  });

  it('proxies id-scoped read, update, and delete requests with encoded ids', async () => {
    vi.stubEnv('WEB_API_BASE_URL', 'https://api.example.test');
    const fetchMock = vi
      .spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(jsonResponse({ data: { id: 'user/1' } }))
      .mockResolvedValueOnce(jsonResponse({ data: { id: 'user/1', name: 'Updated' } }))
      .mockResolvedValueOnce(new Response(null, { status: 204 }));
    const params = Promise.resolve({ id: 'user/1' });

    await GET_BY_ID(new Request('http://next.test/api/users/user%2F1'), { params });
    await PATCH(
      new Request('http://next.test/api/users/user%2F1', {
        body: JSON.stringify({ name: 'Updated' }),
        method: 'PATCH',
      }),
      { params },
    );
    const deleteResponse = await DELETE(new Request('http://next.test/api/users/user%2F1'), {
      params,
    });

    expect(fetchMock).toHaveBeenNthCalledWith(
      1,
      'https://api.example.test/api/users/user%2F1',
      expect.objectContaining({ method: 'GET' }),
    );
    expect(fetchMock).toHaveBeenNthCalledWith(
      2,
      'https://api.example.test/api/users/user%2F1',
      expect.objectContaining({ body: JSON.stringify({ name: 'Updated' }), method: 'PATCH' }),
    );
    expect(fetchMock).toHaveBeenNthCalledWith(
      3,
      'https://api.example.test/api/users/user%2F1',
      expect.objectContaining({ method: 'DELETE' }),
    );
    expect(deleteResponse.status).toBe(204);
  });
});
