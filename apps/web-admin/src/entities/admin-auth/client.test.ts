// FILE: apps/web-admin/src/entities/admin-auth/client.test.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify web-admin admin-auth GraphQL client helpers.
//   SCOPE: Covers current-admin, login, logout, and union error normalization through the shared credentialed GraphQL client; excludes React context and route rendering.
//   DEPENDS: apps/web-admin/src/entities/admin-auth/client.ts, apps/web-admin/src/shared/api/graphql-client.ts, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   auth client tests - Prove generated auth operations are requested and normalized.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added auth GraphQL client helper coverage.
// END_CHANGE_SUMMARY

import { beforeEach, describe, expect, it, vi } from 'vitest';
import { graphqlClient } from '@shared/api/graphql-client';
import { fetchCurrentAdmin, loginAdmin, logoutAdmin } from './client';

vi.mock('@shared/api/graphql-client', () => ({
  graphqlClient: {
    request: vi.fn(),
  },
}));

const requestMock = vi.mocked(graphqlClient.request);

describe('admin auth GraphQL client helpers', () => {
  beforeEach(() => {
    requestMock.mockReset();
  });

  it('fetches the current admin through the CurrentAdmin operation', async () => {
    requestMock.mockResolvedValue({
      me: {
        id: 'admin-1',
        email: 'owner@example.test',
        name: 'Owner Admin',
        role: 'ADMIN',
        createdAt: '2026-06-07T00:00:00Z',
        updatedAt: '2026-06-07T00:00:00Z',
      },
    });

    await expect(fetchCurrentAdmin()).resolves.toMatchObject({
      id: 'admin-1',
      email: 'owner@example.test',
    });
    expect(requestMock).toHaveBeenCalledWith(expect.stringContaining('query CurrentAdmin'));
  });

  it('returns the admin on successful login', async () => {
    requestMock.mockResolvedValue({
      loginAdmin: {
        __typename: 'LoginAdminSuccess',
        admin: {
          id: 'admin-1',
          email: 'owner@example.test',
          name: 'Owner Admin',
          role: 'ADMIN',
          createdAt: '2026-06-07T00:00:00Z',
          updatedAt: '2026-06-07T00:00:00Z',
        },
      },
    });

    await expect(
      loginAdmin({ email: 'owner@example.test', password: 'StrongPassword123!' }),
    ).resolves.toMatchObject({ ok: true, admin: { email: 'owner@example.test' } });
    expect(requestMock).toHaveBeenCalledWith(expect.stringContaining('mutation LoginAdmin'), {
      input: { email: 'owner@example.test', password: 'StrongPassword123!' },
    });
  });

  it('normalizes auth and validation login failures', async () => {
    requestMock.mockResolvedValueOnce({
      loginAdmin: { __typename: 'AuthError', message: 'invalid credentials' },
    });
    await expect(
      loginAdmin({ email: 'owner@example.test', password: 'bad-password' }),
    ).resolves.toEqual({ ok: false, error: { message: 'invalid credentials' } });

    requestMock.mockResolvedValueOnce({
      loginAdmin: {
        __typename: 'ValidationError',
        field: 'email',
        message: 'invalid email',
      },
    });
    await expect(loginAdmin({ email: 'broken', password: 'StrongPassword123!' })).resolves.toEqual({
      ok: false,
      error: { field: 'email', message: 'invalid email' },
    });
  });

  it('logs out through the LogoutAdmin operation', async () => {
    requestMock.mockResolvedValue({ logoutAdmin: { __typename: 'LogoutAdminSuccess', ok: true } });

    await expect(logoutAdmin()).resolves.toBe(true);
    expect(requestMock).toHaveBeenCalledWith(expect.stringContaining('mutation LogoutAdmin'));
  });
});
