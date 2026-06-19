// FILE: apps/web-admin/src/entities/admin-auth/provider.test.tsx
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Verify the web-admin admin-auth provider boundary.
//   SCOPE: Covers current-admin query state, login refetch behavior, login error normalization, and logout cache clearing; excludes route navigation and page rendering.
//   DEPENDS: apps/web-admin/src/entities/admin-auth/provider.tsx, apps/web-admin/src/entities/admin-auth/client.ts, @tanstack/react-query, @testing-library/react, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   auth provider tests - Prove auth state and actions are available to route components.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Added completed-logout reset and hook misuse coverage.
// END_CHANGE_SUMMARY

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { AuthProvider, useAdminAuth } from './provider';
import type { AdminPrincipal } from './model';

const fetchCurrentAdminMock = vi.hoisted(() => vi.fn());
const loginAdminMock = vi.hoisted(() => vi.fn());
const logoutAdminMock = vi.hoisted(() => vi.fn());

vi.mock('./client', () => ({
  fetchCurrentAdmin: fetchCurrentAdminMock,
  loginAdmin: loginAdminMock,
  logoutAdmin: logoutAdminMock,
}));

const admin: AdminPrincipal = {
  id: 'admin-1',
  email: 'owner@example.test',
  name: 'Owner Admin',
  role: 'ADMIN',
  createdAt: '2026-06-07T00:00:00Z',
  updatedAt: '2026-06-07T00:00:00Z',
};

let latestAuth: ReturnType<typeof useAdminAuth> | null = null;

function Probe() {
  latestAuth = useAdminAuth();
  return (
    <div>
      <span>{latestAuth.isLoading ? 'loading' : 'ready'}</span>
      <span>{latestAuth.admin?.email ?? 'anonymous'}</span>
      <span>{latestAuth.hasCompletedLogout ? 'completed-logout' : 'not-completed-logout'}</span>
    </div>
  );
}

function OutsideProviderProbe() {
  useAdminAuth();
  return null;
}

function renderProvider({ staleTime = 0 }: { staleTime?: number } = {}) {
  latestAuth = null;
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, staleTime }, mutations: { retry: false } },
  });

  const result = render(
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <Probe />
      </AuthProvider>
    </QueryClientProvider>,
  );

  return { queryClient, ...result };
}

describe('AuthProvider', () => {
  beforeEach(() => {
    fetchCurrentAdminMock.mockReset();
    loginAdminMock.mockReset();
    logoutAdminMock.mockReset();
  });

  it('loads the current admin from the backend session', async () => {
    fetchCurrentAdminMock.mockResolvedValue(admin);

    renderProvider();

    expect(await screen.findByText('owner@example.test')).toBeInTheDocument();
    expect(screen.getByText('ready')).toBeInTheDocument();
  });

  it('logs in and refreshes the current admin even when cached anonymous auth is fresh', async () => {
    fetchCurrentAdminMock.mockResolvedValueOnce(null).mockResolvedValueOnce(admin);
    loginAdminMock.mockResolvedValue({ ok: true, admin });

    renderProvider({ staleTime: 60_000 });
    await screen.findByText('anonymous');

    await act(async () => {
      const result = await latestAuth?.login({
        email: 'owner@example.test',
        password: 'StrongPassword123!',
      });
      expect(result).toEqual({ ok: true, admin });
    });

    await waitFor(() => expect(fetchCurrentAdminMock).toHaveBeenCalledTimes(2));
    expect(await screen.findByText('owner@example.test')).toBeInTheDocument();
  });

  it('returns normalized login errors without refetching current admin', async () => {
    fetchCurrentAdminMock.mockResolvedValue(null);
    loginAdminMock.mockResolvedValue({ ok: false, error: { message: 'invalid credentials' } });

    renderProvider();
    await screen.findByText('anonymous');

    await act(async () => {
      const result = await latestAuth?.login({
        email: 'owner@example.test',
        password: 'bad-password',
      });
      expect(result).toEqual({ ok: false, error: { message: 'invalid credentials' } });
    });

    expect(fetchCurrentAdminMock).toHaveBeenCalledTimes(1);
  });

  it('logs out and clears current admin plus protected route caches', async () => {
    fetchCurrentAdminMock.mockResolvedValue(admin);
    logoutAdminMock.mockResolvedValue(true);

    const { queryClient } = renderProvider();
    queryClient.setQueryData(['admin-users'], { users: { edges: [{ node: { id: 'user-1' } }] } });
    queryClient.setQueryData(['admin-user', 'user-1'], { user: { id: 'user-1' } });
    await screen.findByText('owner@example.test');

    await act(async () => {
      await latestAuth?.logout();
    });

    expect(logoutAdminMock).toHaveBeenCalledTimes(1);
    expect(screen.getByText('completed-logout')).toBeInTheDocument();
    expect(queryClient.getQueryData(['admin-users'])).toBeUndefined();
    expect(queryClient.getQueryData(['admin-user', 'user-1'])).toBeUndefined();
    expect(await screen.findByText('anonymous')).toBeInTheDocument();

    await act(async () => {
      latestAuth?.clearCompletedLogout();
    });

    expect(screen.getByText('not-completed-logout')).toBeInTheDocument();
  });

  it('requires consumers to be rendered within AuthProvider', () => {
    const consoleError = vi.spyOn(console, 'error').mockImplementation(() => undefined);

    try {
      expect(() => render(<OutsideProviderProbe />)).toThrow(
        'useAdminAuth must be used within AuthProvider',
      );
    } finally {
      consoleError.mockRestore();
    }
  });
});
