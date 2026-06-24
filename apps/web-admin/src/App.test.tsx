// FILE: apps/web-admin/src/App.test.tsx
// VERSION: 1.3.0
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
//   LAST_CHANGE: 1.3.0 - Added Atlas weekly nutrition template route coverage.
// END_CHANGE_SUMMARY

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterAll, afterEach, beforeAll, describe, expect, it, vi } from 'vitest';
import App from './App';

const requestMock = vi.hoisted(() => vi.fn());
const getAtlasDailyNutritionLogMock = vi.hoisted(() => vi.fn());
const getAtlasNutritionTemplateCurrentMock = vi.hoisted(() => vi.fn());
const listAtlasNutritionProductsMock = vi.hoisted(() => vi.fn());
const addAtlasDailyNutritionEntryMock = vi.hoisted(() => vi.fn());
const applyAtlasNutritionTemplateToWeekMock = vi.hoisted(() => vi.fn());
const createAtlasNutritionTemplateMock = vi.hoisted(() => vi.fn());
const createAtlasNutritionTemplateItemMock = vi.hoisted(() => vi.fn());
const deleteAtlasNutritionTemplateItemMock = vi.hoisted(() => vi.fn());
const updateAtlasDailyNutritionEntryMock = vi.hoisted(() => vi.fn());
const updateAtlasNutritionTemplateMock = vi.hoisted(() => vi.fn());
const updateAtlasNutritionTemplateItemMock = vi.hoisted(() => vi.fn());
const deleteAtlasDailyNutritionEntryMock = vi.hoisted(() => vi.fn());
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

vi.mock('./pages/atlas/nutrition-api', () => ({
  addAtlasDailyNutritionEntry: addAtlasDailyNutritionEntryMock,
  applyAtlasNutritionTemplateToWeek: applyAtlasNutritionTemplateToWeekMock,
  archiveAtlasNutritionProduct: vi.fn(),
  AtlasNutritionApiError: class MockAtlasNutritionApiError extends Error {
    readonly code: string;
    readonly type: string;

    constructor(message: string, code: string, type: string) {
      super(message);
      this.name = 'AtlasNutritionApiError';
      this.code = code;
      this.type = type;
    }
  },
  createAtlasNutritionProduct: vi.fn(),
  createAtlasNutritionTemplate: createAtlasNutritionTemplateMock,
  createAtlasNutritionTemplateItem: createAtlasNutritionTemplateItemMock,
  deleteAtlasDailyNutritionEntry: deleteAtlasDailyNutritionEntryMock,
  deleteAtlasNutritionTemplateItem: deleteAtlasNutritionTemplateItemMock,
  getAtlasDailyNutritionLog: getAtlasDailyNutritionLogMock,
  getAtlasNutritionTemplateCurrent: getAtlasNutritionTemplateCurrentMock,
  listAtlasNutritionProducts: listAtlasNutritionProductsMock,
  restoreAtlasNutritionProduct: vi.fn(),
  updateAtlasDailyNutritionEntry: updateAtlasDailyNutritionEntryMock,
  updateAtlasNutritionTemplate: updateAtlasNutritionTemplateMock,
  updateAtlasNutritionTemplateItem: updateAtlasNutritionTemplateItemMock,
  updateAtlasNutritionProduct: vi.fn(),
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
  getAtlasDailyNutritionLogMock.mockReset();
  getAtlasNutritionTemplateCurrentMock.mockReset();
  listAtlasNutritionProductsMock.mockReset();
  addAtlasDailyNutritionEntryMock.mockReset();
  applyAtlasNutritionTemplateToWeekMock.mockReset();
  createAtlasNutritionTemplateMock.mockReset();
  createAtlasNutritionTemplateItemMock.mockReset();
  deleteAtlasNutritionTemplateItemMock.mockReset();
  updateAtlasDailyNutritionEntryMock.mockReset();
  updateAtlasNutritionTemplateMock.mockReset();
  updateAtlasNutritionTemplateItemMock.mockReset();
  deleteAtlasDailyNutritionEntryMock.mockReset();
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

  it('renders the Atlas nutrition products route through the browser router', async () => {
    listAtlasNutritionProductsMock.mockResolvedValue([
      {
        id: 'product-1',
        userId: 'user-1',
        name: 'Rice',
        caloriesPer100g: 130,
        proteinPer100g: 2.7,
        fatPer100g: 0.3,
        carbsPer100g: 28,
        notes: 'Base carbs',
        isActive: true,
        createdAt: '2026-06-24T00:00:00Z',
        updatedAt: '2026-06-24T00:00:00Z',
      },
    ]);

    renderApp('/atlas/nutrition/products');

    expect(await screen.findByRole('heading', { name: 'Product Library' })).toBeInTheDocument();
    expect(await screen.findByText('Rice')).toBeInTheDocument();
    expect(screen.getAllByRole('link', { name: 'Nutrition' })[0]).toHaveAttribute(
      'href',
      '/atlas/nutrition',
    );
    expect(listAtlasNutritionProductsMock).toHaveBeenCalledWith({ includeArchived: true });
  });

  it('renders the Atlas nutrition daily log route through the browser router', async () => {
    getAtlasDailyNutritionLogMock.mockResolvedValue({
      id: 'log-1',
      userId: 'user-1',
      date: '2026-06-24',
      notes: null,
      totals: { calories: 325, protein: 6.8, fat: 0.8, carbs: 70 },
      entries: [
        {
          id: 'entry-1',
          dailyLogId: 'log-1',
          productId: 'product-1',
          productNameSnapshot: 'Rice',
          caloriesPer100gSnapshot: 130,
          proteinPer100gSnapshot: 2.7,
          fatPer100gSnapshot: 0.3,
          carbsPer100gSnapshot: 28,
          amountGrams: 250,
          mealLabel: 'Lunch',
          notes: 'Steamed',
          position: 0,
          macros: { calories: 325, protein: 6.8, fat: 0.8, carbs: 70 },
          createdAt: '2026-06-24T00:00:00Z',
          updatedAt: '2026-06-24T00:00:00Z',
        },
      ],
      createdAt: '2026-06-24T00:00:00Z',
      updatedAt: '2026-06-24T00:00:00Z',
    });
    listAtlasNutritionProductsMock.mockResolvedValue([
      {
        id: 'product-1',
        userId: 'user-1',
        name: 'Rice',
        caloriesPer100g: 130,
        proteinPer100g: 2.7,
        fatPer100g: 0.3,
        carbsPer100g: 28,
        notes: 'Base carbs',
        isActive: true,
        createdAt: '2026-06-24T00:00:00Z',
        updatedAt: '2026-06-24T00:00:00Z',
      },
    ]);

    renderApp('/atlas/nutrition');

    expect(await screen.findByRole('heading', { name: 'Nutrition' })).toBeInTheDocument();
    expect(await screen.findByText('Rice')).toBeInTheDocument();
    expect(screen.getByText('325 kcal')).toBeInTheDocument();
    expect(getAtlasDailyNutritionLogMock).toHaveBeenCalled();
    expect(listAtlasNutritionProductsMock).toHaveBeenCalledWith();
  });

  it('renders the Atlas weekly nutrition template route through the browser router', async () => {
    getAtlasNutritionTemplateCurrentMock.mockResolvedValue({
      id: 'template-1',
      userId: 'user-1',
      weekStartDate: '2026-06-22',
      title: 'Base week',
      notes: null,
      items: [
        {
          id: 'item-1',
          templateId: 'template-1',
          productId: 'product-1',
          amountGrams: 250,
          mealLabel: 'Lunch',
          notes: null,
          createdAt: '2026-06-24T00:00:00Z',
          updatedAt: '2026-06-24T00:00:00Z',
        },
      ],
      createdAt: '2026-06-24T00:00:00Z',
      updatedAt: '2026-06-24T00:00:00Z',
    });
    listAtlasNutritionProductsMock.mockResolvedValue([
      {
        id: 'product-1',
        userId: 'user-1',
        name: 'Rice',
        caloriesPer100g: 130,
        proteinPer100g: 2.7,
        fatPer100g: 0.3,
        carbsPer100g: 28,
        notes: 'Base carbs',
        isActive: true,
        createdAt: '2026-06-24T00:00:00Z',
        updatedAt: '2026-06-24T00:00:00Z',
      },
    ]);

    renderApp('/atlas/nutrition/template');

    expect(await screen.findByRole('heading', { name: 'Weekly Plan' })).toBeInTheDocument();
    expect(await screen.findByText('Rice')).toBeInTheDocument();
    expect(screen.getAllByRole('link', { name: 'Nutrition' })[0]).toHaveAttribute(
      'href',
      '/atlas/nutrition',
    );
    expect(getAtlasNutritionTemplateCurrentMock).toHaveBeenCalled();
    expect(listAtlasNutritionProductsMock).toHaveBeenCalledWith();
  });

  it('redirects the legacy daily override route to the factual nutrition route', async () => {
    getAtlasDailyNutritionLogMock.mockResolvedValue({
      id: 'log-1',
      userId: 'user-1',
      date: '2026-06-24',
      notes: null,
      totals: { calories: 0, protein: 0, fat: 0, carbs: 0 },
      entries: [],
      createdAt: '2026-06-24T00:00:00Z',
      updatedAt: '2026-06-24T00:00:00Z',
    });
    listAtlasNutritionProductsMock.mockResolvedValue([]);

    renderApp('/atlas/nutrition/overrides/new');

    expect(await screen.findByRole('heading', { name: 'Nutrition' })).toBeInTheDocument();
    expect(window.location.pathname).toBe('/atlas/nutrition');
    expect(screen.queryByText(/Daily Override/i)).not.toBeInTheDocument();
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
