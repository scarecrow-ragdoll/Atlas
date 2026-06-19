// FILE: apps/web-admin/src/pages/user-detail-page.test.tsx
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Verify the Vite web-admin user detail route behavior.
//   SCOPE: Covers fetched user rendering, UI-kit navigation fields, not-found state, and load failure state; excludes list and create behavior.
//   DEPENDS: apps/web-admin/src/pages/user-detail-page.tsx, @tanstack/react-query, react-router, @testing-library/react, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   UserDetailPage tests - Prove user detail route states through mocked GraphQL requests.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Added red coverage for UI-kit detail navigation and not-found copy.
// END_CHANGE_SUMMARY

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router';
import { afterEach, describe, expect, it, vi } from 'vitest';
import UserDetailPage from './user-detail-page';

const requestMock = vi.hoisted(() => vi.fn());

vi.mock('@shared/api/graphql-client', () => ({
  graphqlClient: {
    request: requestMock,
  },
}));

function renderDetail(path = '/users/user-1') {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[path]}>
        <Routes>
          <Route path="/users/:id" element={<UserDetailPage />} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe('UserDetailPage', () => {
  afterEach(() => {
    vi.resetAllMocks();
  });

  it('renders a fetched user', async () => {
    requestMock.mockResolvedValue({
      user: {
        id: 'user-1',
        email: 'one@example.com',
        name: 'One User',
        createdAt: '2026-05-02T00:00:00Z',
        updatedAt: '2026-05-02T00:00:00Z',
      },
    });

    renderDetail();

    expect(screen.getByRole('status', { name: 'Loading user' })).toBeInTheDocument();
    expect(await screen.findByRole('heading', { name: 'One User' })).toBeInTheDocument();
    expect(requestMock).toHaveBeenCalledWith(expect.any(String), { id: 'user-1' });
    expect(screen.getByRole('link', { name: 'Back to users' })).toHaveAttribute('href', '/users');
    expect(screen.getByText('Email')).toBeInTheDocument();
    expect(screen.getByText('one@example.com')).toBeInTheDocument();
    expect(screen.getByText('user-1')).toBeInTheDocument();
  });

  it('renders not found when GraphQL returns null', async () => {
    requestMock.mockResolvedValue({ user: null });

    renderDetail('/users/missing');

    expect(await screen.findByRole('heading', { name: 'User not found' })).toBeInTheDocument();
    expect(screen.getByText('The requested user does not exist.')).toBeInTheDocument();
  });

  it('renders load failures', async () => {
    requestMock.mockRejectedValue(new Error('network failed'));

    renderDetail();

    expect(await screen.findByText('Failed to load user.')).toBeInTheDocument();
  });
});
