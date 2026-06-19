// FILE: apps/web/src/shared/api/users.test.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify the public web REST users client contract.
//   SCOPE: Covers same-origin browser URLs, request methods, request bodies, response mapping, and REST error mapping; excludes Next route proxy forwarding.
//   DEPENDS: apps/web/src/shared/api/users.ts, apps/web/src/shared/config.ts, vitest.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   REST users client tests - Prove public web REST client requests and errors.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Updated tests for same-origin browser REST proxy calls.
// END_CHANGE_SUMMARY

import { afterEach, describe, expect, it, vi } from 'vitest';
import { createUser, deleteUser, getUser, listUsers, updateUser } from './users';

afterEach(() => vi.restoreAllMocks());

describe('REST users client', () => {
  it('loads users from the REST endpoint', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(JSON.stringify({ data: [], meta: { totalCount: 0 } }), { status: 200 }),
    );

    await expect(listUsers()).resolves.toEqual({ users: [], totalCount: 0 });
    expect(fetch).toHaveBeenCalledWith('/api/users', expect.objectContaining({ method: 'GET' }));
  });

  it('throws API errors with code, message, and field', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(
        JSON.stringify({
          error: { code: 'DUPLICATE_EMAIL', message: 'email already exists', field: 'email' },
        }),
        { status: 409 },
      ),
    );

    await expect(
      createUser({ email: 'taken@example.com', name: 'Taken', password: 'secret123' }),
    ).rejects.toMatchObject({ code: 'DUPLICATE_EMAIL', field: 'email' });
  });

  it('throws fallback HTTP errors when the response has no error envelope', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(JSON.stringify({ data: null }), { status: 500 }),
    );

    await expect(
      createUser({ email: 'broken@example.com', name: 'Broken', password: 'secret123' }),
    ).rejects.toMatchObject({
      code: 'HTTP_ERROR',
      message: 'Request failed with status 500',
      status: 500,
    });
  });

  it('creates users and returns the created user', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(
        JSON.stringify({
          data: {
            id: 'u1',
            email: 'new@example.com',
            name: 'New',
            createdAt: '2026-05-24T00:00:00Z',
            updatedAt: '2026-05-24T00:00:00Z',
          },
        }),
        { status: 201 },
      ),
    );

    await expect(
      createUser({ email: 'new@example.com', name: 'New', password: 'secret123' }),
    ).resolves.toMatchObject({ id: 'u1', email: 'new@example.com' });
    expect(fetch).toHaveBeenCalledWith(
      '/api/users',
      expect.objectContaining({
        body: JSON.stringify({ email: 'new@example.com', name: 'New', password: 'secret123' }),
        method: 'POST',
      }),
    );
  });

  it('loads one user by encoded id', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(
        JSON.stringify({
          data: {
            id: 'user/1',
            email: 'one@example.com',
            name: 'One',
            createdAt: '2026-05-24T00:00:00Z',
            updatedAt: '2026-05-24T00:00:00Z',
          },
        }),
        { status: 200 },
      ),
    );

    await expect(getUser('user/1')).resolves.toMatchObject({ id: 'user/1' });
    expect(fetch).toHaveBeenCalledWith(
      '/api/users/user%2F1',
      expect.objectContaining({ method: 'GET' }),
    );
  });

  it('updates users by encoded id', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(
        JSON.stringify({
          data: {
            id: 'u1',
            email: 'updated@example.com',
            name: 'Updated',
            createdAt: '2026-05-24T00:00:00Z',
            updatedAt: '2026-05-24T00:00:01Z',
          },
        }),
        { status: 200 },
      ),
    );

    await expect(
      updateUser('u1', { email: 'updated@example.com', name: 'Updated' }),
    ).resolves.toMatchObject({
      email: 'updated@example.com',
      name: 'Updated',
    });
    expect(fetch).toHaveBeenCalledWith(
      '/api/users/u1',
      expect.objectContaining({
        body: JSON.stringify({ email: 'updated@example.com', name: 'Updated' }),
        method: 'PATCH',
      }),
    );
  });

  it('deletes users and accepts no-content responses', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(null, { status: 204 }));

    await expect(deleteUser('u1')).resolves.toBeUndefined();
    expect(fetch).toHaveBeenCalledWith(
      '/api/users/u1',
      expect.objectContaining({ method: 'DELETE' }),
    );
  });
});
