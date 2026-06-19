// FILE: apps/web-admin/src/pages/users-page.test.tsx
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Verify the Vite web-admin users route behavior.
//   SCOPE: Covers list loading, empty state, table semantics, returned users, load errors, create success, and create validation/auth errors; excludes GraphQL transport internals.
//   DEPENDS: apps/web-admin/src/pages/users-page.tsx, @tanstack/react-query, react-router, @testing-library/react, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   UsersPage tests - Prove users list and create-form behavior through mocked GraphQL requests.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.3 - Aligned users route assertions with sidebar-owned global navigation.
// END_CHANGE_SUMMARY

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import UsersPage from './users-page';

const requestMock = vi.hoisted(() => vi.fn());

vi.mock('@shared/api/graphql-client', () => ({
  graphqlClient: {
    request: requestMock,
  },
}));

function usersResponse(edges: Array<{ cursor: string; node: unknown }>, totalCount = edges.length) {
  return {
    users: {
      edges,
      pageInfo: { hasNextPage: false, endCursor: null },
      totalCount,
    },
  };
}

function renderUsersPage() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>
        <UsersPage />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

async function fillAndSubmit() {
  fireEvent.change(screen.getByPlaceholderText('Name'), { target: { value: 'Created User' } });
  fireEvent.change(screen.getByPlaceholderText('Email'), {
    target: { value: 'created@example.com' },
  });
  fireEvent.change(screen.getByPlaceholderText('Password'), {
    target: { value: 'Password123!' },
  });
  fireEvent.click(screen.getByRole('button', { name: 'Create' }));
}

describe('UsersPage', () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  it('renders loading, empty, and total count states', async () => {
    requestMock.mockResolvedValue(usersResponse([]));

    renderUsersPage();

    expect(await screen.findByText('No users yet. Create one above.')).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'No users yet' })).toBeInTheDocument();
    expect(screen.getByText('Total: 0')).toBeInTheDocument();
    expect(screen.getByText('Showing the latest 20 users.')).toBeInTheDocument();
  });

  it('renders returned users as detail links', async () => {
    requestMock.mockResolvedValue(
      usersResponse([
        {
          cursor: 'cursor-1',
          node: {
            id: 'user-1',
            email: 'one@example.com',
            name: 'One User',
            createdAt: '2026-05-02T00:00:00Z',
          },
        },
      ]),
    );

    renderUsersPage();

    expect(await screen.findByRole('link', { name: 'One User' })).toHaveAttribute(
      'href',
      '/users/user-1',
    );
    expect(screen.getByRole('columnheader', { name: 'Name' })).toBeInTheDocument();
    expect(screen.getByRole('cell', { name: 'one@example.com' })).toBeInTheDocument();
    expect(screen.getByText('Total: 1')).toBeInTheDocument();
    expect(screen.getByText('Showing the latest 20 users.')).toBeInTheDocument();
  });

  it('shows load errors', async () => {
    requestMock.mockRejectedValue(new Error('network failed'));

    renderUsersPage();

    expect(await screen.findByText('Failed to load users.')).toBeInTheDocument();
  });

  it('creates a user, invalidates the list, and clears the form', async () => {
    requestMock
      .mockResolvedValueOnce(usersResponse([]))
      .mockResolvedValueOnce({
        createUser: { __typename: 'CreateUserSuccess', user: { id: 'u2' } },
      })
      .mockResolvedValueOnce(usersResponse([]));

    renderUsersPage();
    await screen.findByText('No users yet. Create one above.');
    await fillAndSubmit();

    await waitFor(() => expect(requestMock).toHaveBeenCalledTimes(3));
    expect(screen.getByPlaceholderText('Name')).toHaveValue('');
    expect(screen.getByPlaceholderText('Email')).toHaveValue('');
    expect(screen.getByPlaceholderText('Password')).toHaveValue('');
  });

  it('shows pending create state while the mutation is in flight', async () => {
    let resolveCreate: (value: unknown) => void = () => undefined;
    const createPromise = new Promise((resolve) => {
      resolveCreate = resolve;
    });
    requestMock
      .mockResolvedValueOnce(usersResponse([]))
      .mockImplementationOnce(() => createPromise)
      .mockResolvedValueOnce(usersResponse([]));

    renderUsersPage();
    await screen.findByText('No users yet. Create one above.');
    await fillAndSubmit();

    expect(await screen.findByRole('button', { name: 'Creating...' })).toBeDisabled();
    resolveCreate({ createUser: { __typename: 'CreateUserSuccess', user: { id: 'u2' } } });
    await waitFor(() => expect(screen.getByRole('button', { name: 'Create' })).toBeEnabled());
  });

  it('shows validation and auth errors from createUser', async () => {
    requestMock.mockResolvedValueOnce(usersResponse([])).mockResolvedValueOnce({
      createUser: { __typename: 'ValidationError', field: 'email', message: 'already exists' },
    });

    renderUsersPage();
    await screen.findByText('No users yet. Create one above.');
    await fillAndSubmit();

    expect(await screen.findByText('email: already exists')).toBeInTheDocument();
  });

  it('shows non-field createUser result errors', async () => {
    requestMock.mockResolvedValueOnce(usersResponse([])).mockResolvedValueOnce({
      createUser: { __typename: 'AuthError', message: 'not authorized' },
    });

    renderUsersPage();
    await screen.findByText('No users yet. Create one above.');
    await fillAndSubmit();

    expect(await screen.findByText('not authorized')).toBeInTheDocument();
  });

  it('shows mutation error messages and fallback unknown errors', async () => {
    requestMock
      .mockResolvedValueOnce(usersResponse([]))
      .mockRejectedValueOnce(new Error('network failed'));

    renderUsersPage();
    await screen.findByText('No users yet. Create one above.');
    await fillAndSubmit();

    expect(await screen.findByText('network failed')).toBeInTheDocument();

    cleanup();
    requestMock.mockResolvedValueOnce(usersResponse([])).mockRejectedValueOnce('boom');
    renderUsersPage();
    await screen.findByText('No users yet. Create one above.');
    await fillAndSubmit();

    expect(await screen.findByText('Request failed')).toBeInTheDocument();
  });
});
