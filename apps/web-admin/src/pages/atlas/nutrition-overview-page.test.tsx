// FILE: apps/web-admin/src/pages/atlas/nutrition-overview-page.test.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify the Atlas Nutrition overview behaves as a factual daily food log.
//   SCOPE: Covers daily-log loading, date switching, add/edit/delete entry flows, validation, errors, empty/loading states, and EN/RU labels; excludes GraphQL transport internals.
//   DEPENDS: apps/web-admin/src/pages/atlas/nutrition-overview-page.tsx, apps/web-admin/src/pages/atlas/nutrition-api.ts, apps/web-admin/src/app/i18n.tsx, @tanstack/react-query, react-router, @testing-library/react, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-API-NUTRITION.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   NutritionOverviewPage tests - Prove the nutrition route uses factual daily product entries instead of mock target/override UI.
// END_MODULE_MAP

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { I18nProvider, type Language } from '../../app/i18n';
import {
  addAtlasDailyNutritionEntry,
  AtlasNutritionApiError,
  deleteAtlasDailyNutritionEntry,
  getAtlasDailyNutritionLog,
  listAtlasNutritionProducts,
  updateAtlasDailyNutritionEntry,
  type AtlasDailyNutritionEntry,
  type AtlasDailyNutritionLog,
  type AtlasNutritionMacros,
  type AtlasNutritionProduct,
} from './nutrition-api';
import NutritionOverviewPage from './nutrition-overview-page';

vi.mock('./nutrition-api', () => {
  class MockAtlasNutritionApiError extends Error {
    readonly code: string;
    readonly type: string;

    constructor(message: string, code: string, type: string) {
      super(message);
      this.name = 'AtlasNutritionApiError';
      this.code = code;
      this.type = type;
    }
  }

  return {
    addAtlasDailyNutritionEntry: vi.fn(),
    AtlasNutritionApiError: MockAtlasNutritionApiError,
    deleteAtlasDailyNutritionEntry: vi.fn(),
    getAtlasDailyNutritionLog: vi.fn(),
    listAtlasNutritionProducts: vi.fn(),
    updateAtlasDailyNutritionEntry: vi.fn(),
  };
});

const getDailyLogMock = vi.mocked(getAtlasDailyNutritionLog);
const addEntryMock = vi.mocked(addAtlasDailyNutritionEntry);
const updateEntryMock = vi.mocked(updateAtlasDailyNutritionEntry);
const deleteEntryMock = vi.mocked(deleteAtlasDailyNutritionEntry);
const listProductsMock = vi.mocked(listAtlasNutritionProducts);

const zeroMacros: AtlasNutritionMacros = {
  calories: 0,
  protein: 0,
  fat: 0,
  carbs: 0,
};

function makeProduct(overrides: Partial<AtlasNutritionProduct> = {}): AtlasNutritionProduct {
  const id = overrides.id ?? 'product-1';
  return {
    id,
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
    ...overrides,
  };
}

function makeEntry(overrides: Partial<AtlasDailyNutritionEntry> = {}): AtlasDailyNutritionEntry {
  const amountGrams = overrides.amountGrams ?? 100;
  const caloriesPer100gSnapshot = overrides.caloriesPer100gSnapshot ?? 130;
  const proteinPer100gSnapshot = overrides.proteinPer100gSnapshot ?? 2.7;
  const fatPer100gSnapshot = overrides.fatPer100gSnapshot ?? 0.3;
  const carbsPer100gSnapshot = overrides.carbsPer100gSnapshot ?? 28;

  return {
    id: 'entry-1',
    dailyLogId: 'log-1',
    productId: 'product-1',
    productNameSnapshot: 'Rice',
    caloriesPer100gSnapshot,
    proteinPer100gSnapshot,
    fatPer100gSnapshot,
    carbsPer100gSnapshot,
    amountGrams,
    mealLabel: 'Lunch',
    notes: 'Steamed',
    position: 0,
    macros: {
      calories: (caloriesPer100gSnapshot * amountGrams) / 100,
      protein: (proteinPer100gSnapshot * amountGrams) / 100,
      fat: (fatPer100gSnapshot * amountGrams) / 100,
      carbs: (carbsPer100gSnapshot * amountGrams) / 100,
    },
    createdAt: '2026-06-24T00:00:00Z',
    updatedAt: '2026-06-24T00:00:00Z',
    ...overrides,
  };
}

function makeDailyLog(overrides: Partial<AtlasDailyNutritionLog> = {}): AtlasDailyNutritionLog {
  return {
    id: 'log-1',
    userId: 'user-1',
    date: '2026-06-24',
    notes: null,
    entries: [],
    totals: zeroMacros,
    createdAt: '2026-06-24T00:00:00Z',
    updatedAt: '2026-06-24T00:00:00Z',
    ...overrides,
  };
}

function renderNutritionPage(language: Language = 'en') {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <I18nProvider initialLanguage={language}>
        <MemoryRouter>
          <NutritionOverviewPage initialDate="2026-06-24" />
        </MemoryRouter>
      </I18nProvider>
    </QueryClientProvider>,
  );
}

function createDeferred<T>() {
  let resolve!: (value: T) => void;
  let reject!: (reason?: unknown) => void;
  const promise = new Promise<T>((promiseResolve, promiseReject) => {
    resolve = promiseResolve;
    reject = promiseReject;
  });

  return { promise, reject, resolve };
}

function getEntryRow(productName: string) {
  return screen.getByText(productName).closest('tr') as HTMLTableRowElement;
}

async function chooseProduct(productName: RegExp) {
  if (!HTMLElement.prototype.scrollIntoView) {
    HTMLElement.prototype.scrollIntoView = vi.fn();
  }
  fireEvent.keyDown(await screen.findByLabelText('Product'), { key: 'ArrowDown' });
  fireEvent.click(await screen.findByRole('option', { name: productName }));
}

describe('NutritionOverviewPage', () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  afterEach(() => {
    cleanup();
  });

  it('renders daily food entries and recalculates totals after adding food', async () => {
    listProductsMock.mockResolvedValue([makeProduct({ id: 'product-1', name: 'Rice' })]);
    getDailyLogMock.mockResolvedValue(makeDailyLog());
    addEntryMock.mockResolvedValue(
      makeDailyLog({
        totals: { calories: 325, protein: 6.75, fat: 0.75, carbs: 70 },
        entries: [makeEntry({ productNameSnapshot: 'Rice', amountGrams: 250 })],
      }),
    );

    renderNutritionPage();

    expect(await screen.findByRole('heading', { name: 'Nutrition' })).toBeInTheDocument();
    expect(screen.getByText('Wednesday, June 24, 2026')).toBeInTheDocument();

    await chooseProduct(/Rice/);
    fireEvent.change(screen.getByLabelText('Grams'), { target: { value: '250' } });
    fireEvent.click(screen.getByRole('button', { name: 'Add food' }));

    expect(await screen.findByText('Food entry added')).toBeInTheDocument();
    expect(screen.getByText('Rice')).toBeInTheDocument();
    expect(screen.getByText('325 kcal')).toBeInTheDocument();
    expect(addEntryMock).toHaveBeenCalledWith({
      date: '2026-06-24',
      productId: 'product-1',
      amountGrams: 250,
      mealLabel: null,
      notes: null,
    });
  });

  it('changes the requested date from the date switcher', async () => {
    listProductsMock.mockResolvedValue([]);
    getDailyLogMock
      .mockResolvedValueOnce(makeDailyLog({ date: '2026-06-24' }))
      .mockResolvedValueOnce(makeDailyLog({ date: '2026-06-25' }));

    renderNutritionPage();

    expect(await screen.findByText('Wednesday, June 24, 2026')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'Next day' }));

    await waitFor(() => expect(getDailyLogMock).toHaveBeenLastCalledWith('2026-06-25'));
    expect(await screen.findByText('Thursday, June 25, 2026')).toBeInTheDocument();
  });

  it('edits grams and updates the entry contribution', async () => {
    const originalEntry = makeEntry({ amountGrams: 100, macros: { ...zeroMacros, calories: 130 } });
    const updatedEntry = makeEntry({
      amountGrams: 200,
      macros: { calories: 260, protein: 5.4, fat: 0.6, carbs: 56 },
    });
    listProductsMock.mockResolvedValue([makeProduct()]);
    getDailyLogMock.mockResolvedValue(
      makeDailyLog({ totals: originalEntry.macros, entries: [originalEntry] }),
    );
    updateEntryMock.mockResolvedValue(
      makeDailyLog({ totals: updatedEntry.macros, entries: [updatedEntry] }),
    );

    renderNutritionPage();

    expect(await screen.findByText('Rice')).toBeInTheDocument();
    const row = within(getEntryRow('Rice'));
    fireEvent.click(row.getByRole('button', { name: 'Edit Rice' }));
    fireEvent.change(screen.getByLabelText('Grams for Rice'), { target: { value: '200' } });
    fireEvent.click(screen.getByRole('button', { name: 'Save Rice' }));

    expect(await screen.findByText('Food entry updated')).toBeInTheDocument();
    expect(screen.getByText('260 kcal')).toBeInTheDocument();
    expect(updateEntryMock).toHaveBeenCalledWith('entry-1', {
      dailyLogId: 'log-1',
      amountGrams: 200,
      mealLabel: 'Lunch',
      notes: 'Steamed',
      position: 0,
    });
  });

  it('deletes entries and renders the empty day state', async () => {
    const entry = makeEntry({ amountGrams: 100 });
    listProductsMock.mockResolvedValue([makeProduct()]);
    getDailyLogMock.mockResolvedValue(makeDailyLog({ totals: entry.macros, entries: [entry] }));
    deleteEntryMock.mockResolvedValue(makeDailyLog());

    renderNutritionPage();

    const row = within(await screen.findByText('Rice').then(() => getEntryRow('Rice')));
    fireEvent.click(row.getByRole('button', { name: 'Delete Rice' }));

    expect(await screen.findByText('Food entry deleted')).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'No food entries yet' })).toBeInTheDocument();
    expect(deleteEntryMock).toHaveBeenCalledWith('entry-1');
  });

  it('renders the loading state while the day or product list is pending', () => {
    getDailyLogMock.mockReturnValue(new Promise(() => undefined));
    listProductsMock.mockReturnValue(new Promise(() => undefined));

    renderNutritionPage();

    expect(screen.getByRole('status', { name: 'Loading nutrition data' })).toBeInTheDocument();
  });

  it('renders an empty day with zero totals and no body weight field', async () => {
    listProductsMock.mockResolvedValue([]);
    getDailyLogMock.mockResolvedValue(makeDailyLog());

    renderNutritionPage();

    expect(await screen.findByRole('heading', { name: 'No food entries yet' })).toBeInTheDocument();
    expect(screen.getByText('0 kcal')).toBeInTheDocument();
    expect(screen.queryByLabelText(/body weight/i)).not.toBeInTheDocument();
    expect(screen.queryByText(/body weight/i)).not.toBeInTheDocument();
  });

  it('renders load errors and retries daily log loading', async () => {
    listProductsMock.mockResolvedValue([]);
    getDailyLogMock.mockRejectedValueOnce(new Error('API unavailable'));

    renderNutritionPage();

    expect(await screen.findByText('Failed to load nutrition data')).toBeInTheDocument();

    getDailyLogMock.mockResolvedValueOnce(makeDailyLog());
    fireEvent.click(screen.getByRole('button', { name: 'Retry' }));

    expect(await screen.findByRole('heading', { name: 'No food entries yet' })).toBeInTheDocument();
    expect(getDailyLogMock).toHaveBeenCalledTimes(2);
  });

  it('does not render the food form while products fail to load', async () => {
    getDailyLogMock.mockResolvedValue(makeDailyLog());
    listProductsMock.mockRejectedValueOnce(new Error('Products unavailable'));

    renderNutritionPage();

    expect(await screen.findByText('Failed to load nutrition data')).toBeInTheDocument();
    expect(screen.getByText('Products unavailable')).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Add food' })).not.toBeInTheDocument();
    expect(screen.queryByRole('heading', { name: 'Food entries' })).not.toBeInTheDocument();
  });

  it('validates missing product and non-positive grams before adding food', async () => {
    listProductsMock.mockResolvedValue([makeProduct()]);
    getDailyLogMock.mockResolvedValue(makeDailyLog());

    renderNutritionPage();
    await screen.findByRole('button', { name: 'Add food' });

    fireEvent.click(screen.getByRole('button', { name: 'Add food' }));

    expect(await screen.findByText('Choose a product')).toBeInTheDocument();
    expect(addEntryMock).not.toHaveBeenCalled();

    await chooseProduct(/Rice/);
    fireEvent.change(screen.getByLabelText('Grams'), { target: { value: '0' } });
    fireEvent.click(screen.getByRole('button', { name: 'Add food' }));

    expect(await screen.findByText('Grams must be greater than 0')).toBeInTheDocument();
    expect(addEntryMock).not.toHaveBeenCalled();
  });

  it('does not show stale mutation success after switching dates before add resolves', async () => {
    const addEntryDeferred = createDeferred<AtlasDailyNutritionLog>();
    listProductsMock.mockResolvedValue([makeProduct({ id: 'product-1', name: 'Rice' })]);
    getDailyLogMock
      .mockResolvedValueOnce(makeDailyLog({ date: '2026-06-24' }))
      .mockResolvedValueOnce(makeDailyLog({ date: '2026-06-25' }));
    addEntryMock.mockReturnValue(addEntryDeferred.promise);

    renderNutritionPage();

    await chooseProduct(/Rice/);
    fireEvent.change(screen.getByLabelText('Grams'), { target: { value: '250' } });
    fireEvent.click(screen.getByRole('button', { name: 'Add food' }));
    await waitFor(() => expect(addEntryMock).toHaveBeenCalledTimes(1));

    fireEvent.click(screen.getByRole('button', { name: 'Next day' }));
    expect(await screen.findByText('Thursday, June 25, 2026')).toBeInTheDocument();

    await act(async () => {
      addEntryDeferred.resolve(
        makeDailyLog({
          date: '2026-06-24',
          totals: { calories: 325, protein: 6.75, fat: 0.75, carbs: 70 },
          entries: [makeEntry({ productNameSnapshot: 'Rice', amountGrams: 250 })],
        }),
      );
      await addEntryDeferred.promise;
    });

    expect(screen.queryByText('Food entry added')).not.toBeInTheDocument();
  });

  it('shows API validation and not-found errors in the relevant action area', async () => {
    const entry = makeEntry();
    listProductsMock.mockResolvedValue([makeProduct()]);
    getDailyLogMock.mockResolvedValue(makeDailyLog({ totals: entry.macros, entries: [entry] }));
    addEntryMock.mockRejectedValue(
      new AtlasNutritionApiError('Product is archived', 'VALIDATION_ERROR', 'validation'),
    );
    deleteEntryMock.mockRejectedValue(
      new AtlasNutritionApiError('Entry not found', 'NOT_FOUND', 'not_found'),
    );

    renderNutritionPage();
    expect(await screen.findByText('Rice')).toBeInTheDocument();

    await chooseProduct(/Rice/);
    fireEvent.change(screen.getByLabelText('Grams'), { target: { value: '100' } });
    fireEvent.click(screen.getByRole('button', { name: 'Add food' }));

    expect(await screen.findByText('Product is archived')).toBeInTheDocument();

    fireEvent.click(within(getEntryRow('Rice')).getByRole('button', { name: 'Delete Rice' }));

    expect(await screen.findByText('Entry not found')).toBeInTheDocument();
  });

  it('renders Russian nutrition labels when Russian is selected', async () => {
    listProductsMock.mockResolvedValue([]);
    getDailyLogMock.mockResolvedValue(makeDailyLog());

    renderNutritionPage('ru');

    expect(await screen.findByRole('heading', { name: 'Питание' })).toBeInTheDocument();
    expect(await screen.findByRole('button', { name: 'Добавить продукт' })).toBeInTheDocument();
  });
});
