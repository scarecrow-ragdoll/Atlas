// FILE: apps/web/app/__tests__/users-client.test.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify the public Next users client component.
//   SCOPE: Covers initial server data rendering, REST create/refetch, duplicate-email errors, and selected-user detail behavior; excludes route proxy internals.
//   DEPENDS: apps/web/app/users-client.tsx, @tanstack/react-query, @testing-library/react, vitest.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   UsersClient tests - Prove interactive public REST users behavior in the Next client component.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.2 - Added saved-theme and toggle-back coverage for the public theme switch.
// END_CHANGE_SUMMARY

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import type { User } from '../../src/shared/api/users';
import UsersClient from '../users-client';

function renderUsersClient(
  initialUsers: User[] = [],
  totalCount = initialUsers.length,
  options: { refetchOnMount?: boolean | 'always'; staleTime?: number } = {},
) {
  const client = new QueryClient({
    defaultOptions: {
      mutations: { retry: false },
      queries: {
        refetchOnMount: options.refetchOnMount ?? false,
        retry: false,
        staleTime: options.staleTime ?? 60_000,
      },
    },
  });
  return render(
    <QueryClientProvider client={client}>
      <UsersClient initialTotalCount={totalCount} initialUsers={initialUsers} />
    </QueryClientProvider>,
  );
}

afterEach(() => {
  cleanup();
  document.documentElement.classList.remove('dark');
  window.localStorage.clear();
  vi.restoreAllMocks();
});

describe('UsersClient', () => {
  it('renders initial users from the server page', () => {
    renderUsersClient([
      {
        createdAt: '2026-05-24T00:00:00Z',
        email: 'one@example.com',
        id: 'u1',
        name: 'One',
        updatedAt: '2026-05-24T00:00:00Z',
      },
    ]);

    expect(screen.getByText('One')).toBeInTheDocument();
    expect(screen.getByText('one@example.com')).toBeInTheDocument();
    expect(screen.getByText('1 users')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Create' })).toHaveAttribute('data-slot', 'button');
    expect(screen.getByLabelText('Name')).toHaveAttribute('data-slot', 'input');
  });

  it('toggles and persists the public web theme', () => {
    renderUsersClient();

    const toggle = screen.getByRole('button', { name: 'Switch to dark theme' });
    fireEvent.click(toggle);

    expect(document.documentElement).toHaveClass('dark');
    expect(window.localStorage.getItem('web-theme')).toBe('dark');
    const lightToggle = screen.getByRole('button', { name: 'Switch to light theme' });
    expect(lightToggle).toBeInTheDocument();

    fireEvent.click(lightToggle);

    expect(document.documentElement).not.toHaveClass('dark');
    expect(window.localStorage.getItem('web-theme')).toBe('light');
    expect(screen.getByRole('button', { name: 'Switch to dark theme' })).toBeInTheDocument();
  });

  it('restores a saved dark public web theme on mount', () => {
    window.localStorage.setItem('web-theme', 'dark');

    renderUsersClient();

    expect(document.documentElement).toHaveClass('dark');
    expect(screen.getByRole('button', { name: 'Switch to light theme' })).toBeInTheDocument();
  });

  it('shows load errors from list refetches', async () => {
    vi.spyOn(globalThis, 'fetch').mockRejectedValueOnce(new Error('offline'));

    renderUsersClient([], 0, { refetchOnMount: 'always', staleTime: 0 });

    expect(await screen.findByText('Failed to load users.')).toBeInTheDocument();
  });

  it('creates a user and refreshes the list through REST', async () => {
    vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(
        new Response(
          JSON.stringify({
            data: {
              createdAt: '2026-05-24T00:00:00Z',
              email: 'new@example.com',
              id: 'u1',
              name: 'New User',
              updatedAt: '2026-05-24T00:00:00Z',
            },
          }),
          { status: 201 },
        ),
      )
      .mockResolvedValueOnce(
        new Response(
          JSON.stringify({
            data: [
              {
                createdAt: '2026-05-24T00:00:00Z',
                email: 'new@example.com',
                id: 'u1',
                name: 'New User',
                updatedAt: '2026-05-24T00:00:00Z',
              },
            ],
            meta: { totalCount: 1 },
          }),
          { status: 200 },
        ),
      );

    renderUsersClient();
    fireEvent.change(screen.getByPlaceholderText('Name'), { target: { value: 'New User' } });
    fireEvent.change(screen.getByPlaceholderText('Email'), {
      target: { value: 'new@example.com' },
    });
    fireEvent.change(screen.getByPlaceholderText('Password'), { target: { value: 'secret123' } });
    fireEvent.click(screen.getByRole('button', { name: 'Create' }));

    expect(await screen.findByText('New User')).toBeInTheDocument();
    await waitFor(() => expect(screen.getByText('new@example.com')).toBeInTheDocument());
  });

  it('shows duplicate email errors from REST', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce(
      new Response(
        JSON.stringify({
          error: { code: 'DUPLICATE_EMAIL', field: 'email', message: 'email already exists' },
        }),
        { status: 409 },
      ),
    );

    renderUsersClient();
    fireEvent.change(screen.getByPlaceholderText('Name'), { target: { value: 'Taken' } });
    fireEvent.change(screen.getByPlaceholderText('Email'), {
      target: { value: 'taken@example.com' },
    });
    fireEvent.change(screen.getByPlaceholderText('Password'), { target: { value: 'secret123' } });
    fireEvent.click(screen.getByRole('button', { name: 'Create' }));

    expect(await screen.findByText('email: email already exists')).toBeInTheDocument();
  });

  it('shows non-field REST errors', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce(
      new Response(JSON.stringify({ error: { code: 'AUTH_ERROR', message: 'not authorized' } }), {
        status: 401,
      }),
    );

    renderUsersClient();
    fireEvent.change(screen.getByPlaceholderText('Name'), { target: { value: 'Blocked' } });
    fireEvent.change(screen.getByPlaceholderText('Email'), {
      target: { value: 'blocked@example.com' },
    });
    fireEvent.change(screen.getByPlaceholderText('Password'), { target: { value: 'secret123' } });
    fireEvent.click(screen.getByRole('button', { name: 'Create' }));

    expect(await screen.findByText('not authorized')).toBeInTheDocument();
  });

  it('shows generic create failures', async () => {
    vi.spyOn(globalThis, 'fetch').mockRejectedValueOnce(new Error('network failed'));

    renderUsersClient();
    fireEvent.change(screen.getByPlaceholderText('Name'), { target: { value: 'Offline' } });
    fireEvent.change(screen.getByPlaceholderText('Email'), {
      target: { value: 'offline@example.com' },
    });
    fireEvent.change(screen.getByPlaceholderText('Password'), { target: { value: 'secret123' } });
    fireEvent.click(screen.getByRole('button', { name: 'Create' }));

    expect(await screen.findByText('network failed')).toBeInTheDocument();
  });

  it('opens a selected user detail panel', () => {
    renderUsersClient([
      {
        createdAt: '2026-05-24T00:00:00Z',
        email: 'one@example.com',
        id: 'u1',
        name: 'One',
        updatedAt: '2026-05-24T00:00:00Z',
      },
    ]);

    fireEvent.click(screen.getByRole('button', { name: 'One' }));

    expect(
      within(screen.getByLabelText('Selected user')).getByText('one@example.com'),
    ).toBeInTheDocument();
  });
});
