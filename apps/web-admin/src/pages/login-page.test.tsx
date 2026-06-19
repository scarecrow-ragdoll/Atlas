// FILE: apps/web-admin/src/pages/login-page.test.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify the web-admin login page behavior.
//   SCOPE: Covers form rendering, successful safe redirects, unsafe return fallback, already-authenticated redirect, and visible auth/validation/network errors; excludes backend session creation.
//   DEPENDS: apps/web-admin/src/pages/login-page.tsx, apps/web-admin/src/entities/admin-auth/provider.tsx, react-router, @testing-library/react, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   login page tests - Prove login UX and redirect behavior.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added login page coverage for route guard behavior.
// END_CHANGE_SUMMARY

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { BrowserRouter, Route, Routes } from 'react-router';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { AuthProvider } from '@entities/admin-auth/provider';
import LoginPage from './login-page';

const fetchCurrentAdminMock = vi.hoisted(() => vi.fn());
const loginAdminMock = vi.hoisted(() => vi.fn());
const logoutAdminMock = vi.hoisted(() => vi.fn());

vi.mock('@entities/admin-auth/client', () => ({
  fetchCurrentAdmin: fetchCurrentAdminMock,
  loginAdmin: loginAdminMock,
  logoutAdmin: logoutAdminMock,
}));

const admin = {
  id: 'admin-1',
  email: 'owner@example.test',
  name: 'Owner Admin',
  role: 'ADMIN',
  createdAt: '2026-06-07T00:00:00Z',
  updatedAt: '2026-06-07T00:00:00Z',
};

function renderLogin(path = '/login') {
  window.history.pushState({}, '', path);
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, staleTime: 60_000 }, mutations: { retry: false } },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <AuthProvider>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/" element={<h1>Overview</h1>} />
            <Route path="/users" element={<h1>Users</h1>} />
          </Routes>
        </AuthProvider>
      </BrowserRouter>
    </QueryClientProvider>,
  );
}

async function submitLogin() {
  fireEvent.change(screen.getByLabelText('Email'), {
    target: { value: 'owner@example.test' },
  });
  fireEvent.change(screen.getByLabelText('Password'), {
    target: { value: 'StrongPassword123!' },
  });
  fireEvent.click(screen.getByRole('button', { name: 'Sign in' }));
}

describe('LoginPage', () => {
  beforeEach(() => {
    fetchCurrentAdminMock.mockReset();
    loginAdminMock.mockReset();
    logoutAdminMock.mockReset();
    fetchCurrentAdminMock.mockResolvedValue(null);
  });

  it('renders the login form', async () => {
    renderLogin();

    expect(await screen.findByRole('heading', { name: 'Admin sign in' })).toBeInTheDocument();
    expect(screen.getByLabelText('Email')).toBeInTheDocument();
    expect(screen.getByLabelText('Password')).toHaveAttribute('type', 'password');
    expect(screen.getByRole('button', { name: 'Sign in' })).toBeInTheDocument();
  });

  it('redirects to a safe return path after login', async () => {
    fetchCurrentAdminMock.mockResolvedValueOnce(null).mockResolvedValueOnce(admin);
    loginAdminMock.mockResolvedValue({ ok: true, admin });

    renderLogin('/login?from=%2Fusers%3Fstatus%3Dactive%23directory');
    await submitLogin();

    await waitFor(() => expect(window.location.pathname).toBe('/users'));
    expect(window.location.search).toBe('?status=active');
    expect(window.location.hash).toBe('#directory');
  });

  it('falls back to overview for unsafe return paths', async () => {
    fetchCurrentAdminMock.mockResolvedValueOnce(null).mockResolvedValueOnce(admin);
    loginAdminMock.mockResolvedValue({ ok: true, admin });

    renderLogin('/login?from=https%3A%2F%2Fevil.example%2Fusers');
    await submitLogin();

    await waitFor(() => expect(window.location.pathname).toBe('/'));
  });

  it('redirects an already-authenticated admin away from login', async () => {
    fetchCurrentAdminMock.mockResolvedValue(admin);

    renderLogin('/login?from=%2Fusers');

    await waitFor(() => expect(window.location.pathname).toBe('/users'));
  });

  it('renders auth and validation errors', async () => {
    loginAdminMock.mockResolvedValueOnce({
      ok: false,
      error: { message: 'invalid credentials' },
    });

    renderLogin();
    await submitLogin();

    expect(await screen.findByText('invalid credentials')).toBeInTheDocument();

    loginAdminMock.mockResolvedValueOnce({
      ok: false,
      error: { field: 'email', message: 'invalid email' },
    });
    await submitLogin();

    expect(await screen.findByText('email: invalid email')).toBeInTheDocument();
  });

  it('renders a stable fallback for network failures', async () => {
    loginAdminMock.mockRejectedValue(new Error('network failed'));

    renderLogin();
    await submitLogin();

    expect(
      await screen.findByText('Unable to sign in. Try again after the API is available.'),
    ).toBeInTheDocument();
  });
});
