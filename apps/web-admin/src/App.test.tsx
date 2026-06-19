// FILE: apps/web-admin/src/App.test.tsx
// VERSION: 1.1.1
// START_MODULE_CONTRACT
//   PURPOSE: Verify the web-admin Vite route table and auth guard at the app boundary.
//   SCOPE: Covers protected home, users, UI-kit, login redirect, and no protected-content flash behavior through browser history; excludes page-level create/detail edge cases.
//   DEPENDS: apps/web-admin/src/App.tsx, @tanstack/react-query, @testing-library/react, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   web-admin routes tests - Prove the Vite admin app exposes home, users, and UI-kit routes.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.1.1 - Added shell logout navigation coverage for protected routes.
// END_CHANGE_SUMMARY

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterAll, afterEach, beforeAll, describe, expect, it, vi } from 'vitest';
import App from './App';

const requestMock = vi.hoisted(() => vi.fn());
const originalMatchMedia = window.matchMedia;
const currentAdmin = {
  id: 'admin-1',
  email: 'owner@example.test',
  name: 'Owner Admin',
  role: 'ADMIN',
  createdAt: '2026-06-07T00:00:00Z',
  updatedAt: '2026-06-07T00:00:00Z',
};

type MockGraphQLRequest = string | { document?: string; operationName?: string };

function getGraphQLDocument(request: MockGraphQLRequest) {
  return typeof request === 'string' ? request : (request.document ?? '');
}

function getGraphQLOperationName(request: MockGraphQLRequest) {
  return typeof request === 'string' ? null : (request.operationName ?? null);
}

function mockGraphQLResponse(request: MockGraphQLRequest) {
  const document = getGraphQLDocument(request);
  const operationName = getGraphQLOperationName(request);

  if (operationName === 'CurrentAdmin' || document.includes('query CurrentAdmin')) {
    return Promise.resolve({ me: currentAdmin });
  }
  if (document.includes('query GetUsers')) {
    return Promise.resolve({
      users: { edges: [], pageInfo: { hasNextPage: false, endCursor: null }, totalCount: 0 },
    });
  }
  if (document.includes('mutation LogoutAdmin')) {
    return Promise.resolve({ logoutAdmin: { __typename: 'LogoutAdminSuccess', ok: true } });
  }
  return Promise.reject(new Error(`Unhandled GraphQL document: ${document.slice(0, 80)}`));
}

vi.mock('@shared/api/graphql-client', () => ({
  graphqlClient: {
    request: requestMock,
  },
}));

beforeAll(() => {
  window.matchMedia = vi.fn().mockImplementation((query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    addListener: vi.fn(),
    removeListener: vi.fn(),
    dispatchEvent: vi.fn(),
  }));
  requestMock.mockImplementation(mockGraphQLResponse);
});

function renderApp(path = '/') {
  window.history.pushState({}, '', path);
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <App />
    </QueryClientProvider>,
  );
}

afterEach(() => {
  cleanup();
  requestMock.mockClear();
  requestMock.mockImplementation(mockGraphQLResponse);
  document.documentElement.classList.remove('dark');
  document.cookie = 'sidebar_state=; path=/; max-age=0';
  window.localStorage.clear();
});

afterAll(() => {
  if (originalMatchMedia) {
    window.matchMedia = originalMatchMedia;
  } else {
    delete (window as Partial<Window>).matchMedia;
  }
});

describe('web-admin routes', () => {
  it('renders the home route with users and UI-kit links', async () => {
    renderApp('/');

    expect(
      await screen.findByRole('heading', { name: 'Monorepo Template Admin' }),
    ).toBeInTheDocument();
    expect(screen.getByRole('navigation', { name: 'Admin navigation' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Users' })).toHaveAttribute('href', '/users');
    expect(screen.getByRole('link', { name: 'UI Kit' })).toHaveAttribute('href', '/ui-kit');
    expect(screen.getByRole('button', { name: 'Toggle sidebar' })).toBeInTheDocument();
    expect(screen.getAllByText('Overview').length).toBeGreaterThan(0);
    expect(screen.getByRole('link', { name: 'Open users' })).toHaveAttribute('href', '/users');
  });

  it('renders the users route through the browser router', async () => {
    renderApp('/users');

    expect(await screen.findByText('No users yet. Create one above.')).toBeInTheDocument();
  });

  it('renders the UI-kit route through the browser router', async () => {
    renderApp('/ui-kit');

    expect(await screen.findByRole('heading', { name: 'UI Kit' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Actions' })).toBeInTheDocument();
  });

  it('toggles and persists the admin theme from the app shell', async () => {
    renderApp('/');

    const toggle = await screen.findByRole('button', { name: 'Switch to dark theme' });
    fireEvent.click(toggle);

    expect(document.documentElement).toHaveClass('dark');
    expect(window.localStorage.getItem('web-admin-theme')).toBe('dark');
    const lightToggle = screen.getByRole('button', { name: 'Switch to light theme' });
    expect(lightToggle).toBeInTheDocument();

    fireEvent.click(lightToggle);

    expect(document.documentElement).not.toHaveClass('dark');
    expect(window.localStorage.getItem('web-admin-theme')).toBe('light');
    expect(screen.getByRole('button', { name: 'Switch to dark theme' })).toBeInTheDocument();
  });

  it('restores a saved dark theme when the app shell mounts', async () => {
    window.localStorage.setItem('web-admin-theme', 'dark');

    renderApp('/');

    expect(
      await screen.findByRole('button', { name: 'Switch to light theme' }),
    ).toBeInTheDocument();
    expect(document.documentElement).toHaveClass('dark');
  });

  it('redirects protected routes to login when no current admin exists', async () => {
    requestMock.mockImplementation((request: MockGraphQLRequest) => {
      const document = getGraphQLDocument(request);
      const operationName = getGraphQLOperationName(request);
      if (operationName === 'CurrentAdmin' || document.includes('query CurrentAdmin')) {
        return Promise.resolve({ me: null });
      }
      return mockGraphQLResponse(request);
    });

    renderApp('/users?status=active#directory');

    expect(await screen.findByRole('heading', { name: 'Admin sign in' })).toBeInTheDocument();
    expect(window.location.pathname).toBe('/login');
    expect(window.location.search).toBe('?from=%2Fusers%3Fstatus%3Dactive%23directory');
    expect(screen.queryByText('No users yet. Create one above.')).not.toBeInTheDocument();
  });

  it('does not render protected content when the current-admin check fails', async () => {
    requestMock.mockImplementation((request: MockGraphQLRequest) => {
      const document = getGraphQLDocument(request);
      const operationName = getGraphQLOperationName(request);
      if (operationName === 'CurrentAdmin' || document.includes('query CurrentAdmin')) {
        return Promise.reject(new Error('API unavailable'));
      }
      return mockGraphQLResponse(request);
    });

    renderApp('/users');

    expect(await screen.findByRole('heading', { name: 'Admin sign in' })).toBeInTheDocument();
    expect(screen.queryByText('No users yet. Create one above.')).not.toBeInTheDocument();
  });

  it('logs out from the app shell and returns to the public login route', async () => {
    renderApp('/users');

    expect(await screen.findByText('No users yet. Create one above.')).toBeInTheDocument();
    const userButton = screen.getByRole('button', { name: /Owner Admin owner@example\.test/i });
    fireEvent.pointerDown(userButton, { button: 0, ctrlKey: false, pointerType: 'mouse' });
    fireEvent.keyDown(userButton, { key: 'Enter' });
    fireEvent.click(await screen.findByRole('menuitem', { name: 'Logout' }));

    expect(await screen.findByRole('heading', { name: 'Admin sign in' })).toBeInTheDocument();
    expect(requestMock).toHaveBeenCalledWith(expect.stringContaining('mutation LogoutAdmin'));
    expect(window.location.pathname).toBe('/login');
    expect(window.location.search).toBe('');
    expect(screen.queryByText('No users yet. Create one above.')).not.toBeInTheDocument();
  });
});
