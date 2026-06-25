// FILE: apps/web-admin/src/pages/atlas/weekly-nutrition-template-page.test.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify the Atlas weekly nutrition template editor behavior.
//   SCOPE: Covers template/product loading, missing-template creation, editable item save reconciliation, seed-empty-days apply results, validation, API errors, empty/loading states, and no body-weight UI; excludes GraphQL transport internals.
//   DEPENDS: apps/web-admin/src/pages/atlas/weekly-nutrition-template-page.tsx, apps/web-admin/src/pages/atlas/nutrition-api.ts, apps/web-admin/src/app/i18n.tsx, @tanstack/react-query, react-router, @testing-library/react, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-API-NUTRITION.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   WeeklyNutritionTemplatePage tests - Prove weekly plan editing persists templates separately from factual day seeding.
// END_MODULE_MAP

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { I18nProvider, type Language } from '../../app/i18n';
import {
  applyAtlasNutritionTemplateToWeek,
  AtlasNutritionApiError,
  createAtlasNutritionTemplate,
  createAtlasNutritionTemplateItem,
  deleteAtlasNutritionTemplateItem,
  getAtlasNutritionTemplateCurrent,
  listAtlasNutritionProducts,
  updateAtlasNutritionTemplate,
  updateAtlasNutritionTemplateItem,
  type AtlasNutritionProduct,
  type AtlasNutritionTemplate,
  type AtlasNutritionTemplateApplyResult,
  type AtlasNutritionTemplateItem,
} from './nutrition-api';
import WeeklyNutritionTemplatePage from './weekly-nutrition-template-page';

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
    applyAtlasNutritionTemplateToWeek: vi.fn(),
    AtlasNutritionApiError: MockAtlasNutritionApiError,
    createAtlasNutritionTemplate: vi.fn(),
    createAtlasNutritionTemplateItem: vi.fn(),
    deleteAtlasNutritionTemplateItem: vi.fn(),
    getAtlasNutritionTemplateCurrent: vi.fn(),
    listAtlasNutritionProducts: vi.fn(),
    updateAtlasNutritionTemplate: vi.fn(),
    updateAtlasNutritionTemplateItem: vi.fn(),
  };
});

const applyTemplateMock = vi.mocked(applyAtlasNutritionTemplateToWeek);
const createTemplateMock = vi.mocked(createAtlasNutritionTemplate);
const createTemplateItemMock = vi.mocked(createAtlasNutritionTemplateItem);
const deleteTemplateItemMock = vi.mocked(deleteAtlasNutritionTemplateItem);
const getCurrentTemplateMock = vi.mocked(getAtlasNutritionTemplateCurrent);
const listProductsMock = vi.mocked(listAtlasNutritionProducts);
const updateTemplateMock = vi.mocked(updateAtlasNutritionTemplate);
const updateTemplateItemMock = vi.mocked(updateAtlasNutritionTemplateItem);

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

function makeTemplateItem(
  overrides: Partial<AtlasNutritionTemplateItem> = {},
): AtlasNutritionTemplateItem {
  return {
    id: overrides.id ?? 'item-1',
    templateId: overrides.templateId ?? 'template-1',
    productId: overrides.productId ?? 'product-1',
    amountGrams: overrides.amountGrams ?? 100,
    mealLabel: overrides.mealLabel ?? 'Lunch',
    notes: overrides.notes ?? 'Steamed',
    createdAt: '2026-06-24T00:00:00Z',
    updatedAt: '2026-06-24T00:00:00Z',
    ...overrides,
  };
}

function makeTemplate(overrides: Partial<AtlasNutritionTemplate> = {}): AtlasNutritionTemplate {
  return {
    id: 'template-1',
    userId: 'user-1',
    weekStartDate: '2026-06-22',
    title: 'Base week',
    notes: 'Steady plan',
    items: [],
    createdAt: '2026-06-24T00:00:00Z',
    updatedAt: '2026-06-24T00:00:00Z',
    ...overrides,
  };
}

function makeApplyResult(
  overrides: Partial<AtlasNutritionTemplateApplyResult> = {},
): AtlasNutritionTemplateApplyResult {
  return {
    weekStartDate: '2026-06-22',
    weekEndDate: '2026-06-28',
    mode: 'SEED_EMPTY_DAYS',
    dates: [
      { date: '2026-06-22', status: 'CREATED', entryCount: 2, reason: null },
      { date: '2026-06-23', status: 'SKIPPED', entryCount: 0, reason: 'already has entries' },
    ],
    ...overrides,
  };
}

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });
}

function renderWeeklyTemplatePage(
  language: Language = 'en',
  queryClient = createTestQueryClient(),
) {
  const result = render(
    <QueryClientProvider client={queryClient}>
      <I18nProvider initialLanguage={language}>
        <MemoryRouter>
          <WeeklyNutritionTemplatePage initialWeekStartDate="2026-06-22" />
        </MemoryRouter>
      </I18nProvider>
    </QueryClientProvider>,
  );

  return { ...result, queryClient };
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

async function chooseProductForEntry(entryNumber: number, productName: RegExp) {
  if (!HTMLElement.prototype.scrollIntoView) {
    HTMLElement.prototype.scrollIntoView = vi.fn();
  }
  fireEvent.keyDown(await screen.findByLabelText(`Product for entry ${entryNumber}`), {
    key: 'ArrowDown',
  });
  fireEvent.click(await screen.findByRole('option', { name: productName }));
}

describe('WeeklyNutritionTemplatePage', () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  afterEach(() => {
    cleanup();
  });

  it('saves template items without applying them to factual days', async () => {
    const rice = makeProduct({ id: 'product-1', name: 'Rice' });
    const oats = makeProduct({ id: 'product-2', name: 'Oats', caloriesPer100g: 389 });
    const milk = makeProduct({ id: 'product-3', name: 'Milk', caloriesPer100g: 64 });
    const riceItem = makeTemplateItem({ id: 'item-1', productId: 'product-1', amountGrams: 100 });
    const oatsItem = makeTemplateItem({ id: 'item-2', productId: 'product-2', amountGrams: 80 });
    getCurrentTemplateMock.mockResolvedValue(makeTemplate({ items: [riceItem, oatsItem] }));
    listProductsMock.mockResolvedValue([rice, oats, milk]);
    updateTemplateMock.mockResolvedValue(
      makeTemplate({ title: 'Training week', notes: 'More carbs' }),
    );
    updateTemplateItemMock.mockResolvedValue(
      makeTemplateItem({ id: 'item-1', productId: 'product-1', amountGrams: 200 }),
    );
    deleteTemplateItemMock.mockResolvedValue(oatsItem);
    createTemplateItemMock.mockResolvedValue(
      makeTemplateItem({ id: 'item-3', productId: 'product-3', amountGrams: 300 }),
    );

    renderWeeklyTemplatePage();

    expect(await screen.findByRole('heading', { name: 'Weekly Plan' })).toBeInTheDocument();
    fireEvent.change(await screen.findByLabelText('Template title'), {
      target: { value: 'Training week' },
    });
    fireEvent.change(screen.getByLabelText('Template notes'), { target: { value: 'More carbs' } });
    fireEvent.change(screen.getByLabelText('Grams for entry 1'), { target: { value: '200' } });
    fireEvent.click(screen.getByRole('button', { name: 'Delete Oats' }));
    fireEvent.click(screen.getByRole('button', { name: 'Add planned entry' }));
    await chooseProductForEntry(2, /Milk/);
    fireEvent.change(screen.getByLabelText('Grams for entry 2'), { target: { value: '300' } });
    fireEvent.change(screen.getByLabelText('Meal label for entry 2'), {
      target: { value: 'Snack' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'Save Template' }));

    expect(await screen.findByText('Template saved')).toBeInTheDocument();
    expect(updateTemplateMock).toHaveBeenCalledWith('template-1', {
      title: 'Training week',
      notes: 'More carbs',
    });
    expect(updateTemplateItemMock).toHaveBeenCalledWith('item-1', {
      amountGrams: 200,
      mealLabel: 'Lunch',
      notes: 'Steamed',
    });
    expect(deleteTemplateItemMock).toHaveBeenCalledWith('item-2');
    expect(createTemplateItemMock).toHaveBeenCalledWith({
      templateId: 'template-1',
      productId: 'product-3',
      amountGrams: 300,
      mealLabel: 'Snack',
      notes: null,
    });
    expect(applyTemplateMock).not.toHaveBeenCalled();
  });

  it('sends empty strings when clearing existing template and item text fields', async () => {
    const riceItem = makeTemplateItem({
      id: 'item-1',
      productId: 'product-1',
      amountGrams: 100,
      mealLabel: 'Lunch',
      notes: 'Steamed',
    });
    getCurrentTemplateMock.mockResolvedValue(
      makeTemplate({ title: 'Base week', notes: 'Steady plan', items: [riceItem] }),
    );
    listProductsMock.mockResolvedValue([makeProduct({ id: 'product-1', name: 'Rice' })]);
    updateTemplateMock.mockResolvedValue(makeTemplate({ title: '', notes: '' }));
    updateTemplateItemMock.mockResolvedValue(
      makeTemplateItem({ id: 'item-1', mealLabel: '', notes: '' }),
    );

    renderWeeklyTemplatePage();

    fireEvent.change(await screen.findByLabelText('Template title'), { target: { value: '' } });
    fireEvent.change(screen.getByLabelText('Template notes'), { target: { value: '' } });
    fireEvent.change(screen.getByLabelText('Meal label for entry 1'), { target: { value: '' } });
    fireEvent.change(screen.getByLabelText('Notes for entry 1'), { target: { value: '' } });
    fireEvent.click(screen.getByRole('button', { name: 'Save Template' }));

    expect(await screen.findByText('Template saved')).toBeInTheDocument();
    expect(updateTemplateMock).toHaveBeenCalledWith('template-1', {
      title: '',
      notes: '',
    });
    expect(updateTemplateItemMock).toHaveBeenCalledWith('item-1', {
      amountGrams: 100,
      mealLabel: '',
      notes: '',
    });
  });

  it('keeps unsaved local edits when the current week query data refetches', async () => {
    getCurrentTemplateMock.mockResolvedValue(makeTemplate({ title: 'Server title' }));
    listProductsMock.mockResolvedValue([makeProduct()]);

    const { queryClient } = renderWeeklyTemplatePage();

    fireEvent.change(await screen.findByLabelText('Template title'), {
      target: { value: 'Unsaved local title' },
    });

    await act(async () => {
      queryClient.setQueryData(
        ['atlas-weekly-nutrition-template', '2026-06-22'],
        makeTemplate({ title: 'Refetched title' }),
      );
    });

    await waitFor(() =>
      expect(screen.getByLabelText('Template title')).toHaveValue('Unsaved local title'),
    );
    expect(screen.queryByDisplayValue('Refetched title')).not.toBeInTheDocument();
  });

  it('disables Apply to Week while a save is pending', async () => {
    const saveDeferred = createDeferred<AtlasNutritionTemplate>();
    getCurrentTemplateMock.mockResolvedValue(makeTemplate());
    listProductsMock.mockResolvedValue([makeProduct()]);
    updateTemplateMock.mockReturnValue(saveDeferred.promise);

    renderWeeklyTemplatePage();

    expect(await screen.findByRole('button', { name: 'Apply to Week' })).toBeEnabled();
    fireEvent.click(screen.getByRole('button', { name: 'Save Template' }));

    await waitFor(() =>
      expect(screen.getByRole('button', { name: 'Apply to Week' })).toBeDisabled(),
    );
    expect(applyTemplateMock).not.toHaveBeenCalled();

    await act(async () => {
      saveDeferred.resolve(makeTemplate());
    });
    expect(await screen.findByText('Template saved')).toBeInTheDocument();
  });

  it('disables Apply to Week when an existing template has unsaved local edits', async () => {
    getCurrentTemplateMock.mockResolvedValue(makeTemplate({ title: 'Saved plan' }));
    listProductsMock.mockResolvedValue([makeProduct()]);

    renderWeeklyTemplatePage();

    expect(await screen.findByRole('button', { name: 'Apply to Week' })).toBeEnabled();
    fireEvent.change(screen.getByLabelText('Template title'), {
      target: { value: 'Unsaved plan' },
    });

    await waitFor(() =>
      expect(screen.getByRole('button', { name: 'Apply to Week' })).toBeDisabled(),
    );
    fireEvent.click(screen.getByRole('button', { name: 'Apply to Week' }));
    expect(applyTemplateMock).not.toHaveBeenCalled();
  });

  it('applies seed_empty_days and reports created and skipped dates', async () => {
    getCurrentTemplateMock.mockResolvedValue(makeTemplate());
    listProductsMock.mockResolvedValue([makeProduct()]);
    applyTemplateMock.mockResolvedValue(makeApplyResult());

    renderWeeklyTemplatePage();

    expect(await screen.findByRole('button', { name: 'Apply to Week' })).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Apply to Week' }));

    expect(await screen.findByText('1 created')).toBeInTheDocument();
    expect(screen.getByText('1 skipped')).toBeInTheDocument();
    expect(screen.getByText('2026-06-22')).toBeInTheDocument();
    expect(screen.getByText('2026-06-23')).toBeInTheDocument();
    expect(screen.getByText('already has entries')).toBeInTheDocument();
    expect(screen.getByText('2 entries')).toBeInTheDocument();
    expect(applyTemplateMock).toHaveBeenCalledWith('template-1', 'SEED_EMPTY_DAYS');
  });

  it('recreates an existing item when its product changes', async () => {
    const rice = makeProduct({ id: 'product-1', name: 'Rice' });
    const milk = makeProduct({ id: 'product-3', name: 'Milk', caloriesPer100g: 64 });
    const riceItem = makeTemplateItem({ id: 'item-1', productId: 'product-1', amountGrams: 100 });
    getCurrentTemplateMock.mockResolvedValue(makeTemplate({ items: [riceItem] }));
    listProductsMock.mockResolvedValue([rice, milk]);
    updateTemplateMock.mockResolvedValue(makeTemplate());
    updateTemplateItemMock.mockResolvedValue(riceItem);
    deleteTemplateItemMock.mockResolvedValue(riceItem);
    createTemplateItemMock.mockResolvedValue(
      makeTemplateItem({ id: 'item-2', productId: 'product-3', amountGrams: 150 }),
    );

    renderWeeklyTemplatePage();

    await chooseProductForEntry(1, /Milk/);
    fireEvent.change(screen.getByLabelText('Grams for entry 1'), { target: { value: '150' } });
    fireEvent.click(screen.getByRole('button', { name: 'Save Template' }));

    expect(await screen.findByText('Template saved')).toBeInTheDocument();
    expect(createTemplateItemMock).toHaveBeenCalledWith({
      templateId: 'template-1',
      productId: 'product-3',
      amountGrams: 150,
      mealLabel: 'Lunch',
      notes: 'Steamed',
    });
    expect(deleteTemplateItemMock).toHaveBeenCalledWith('item-1');
    expect(updateTemplateItemMock).not.toHaveBeenCalled();
  });

  it('renders the loading state while template and product data load', () => {
    const templateDeferred = createDeferred<AtlasNutritionTemplate>();
    const productsDeferred = createDeferred<AtlasNutritionProduct[]>();
    getCurrentTemplateMock.mockReturnValue(templateDeferred.promise);
    listProductsMock.mockReturnValue(productsDeferred.promise);

    renderWeeklyTemplatePage();

    expect(screen.getByRole('status', { name: 'Loading weekly plan' })).toBeInTheDocument();
  });

  it('renders an empty template state without reference-only blocks', async () => {
    getCurrentTemplateMock.mockResolvedValue(null);
    listProductsMock.mockResolvedValue([]);

    renderWeeklyTemplatePage();

    expect(await screen.findByRole('heading', { name: 'No weekly plan yet' })).toBeInTheDocument();
    expect(screen.queryByText(/represented states/i)).not.toBeInTheDocument();
    expect(screen.queryByText(/mock/i)).not.toBeInTheDocument();
  });

  it('renders load errors and retries when template data fails', async () => {
    getCurrentTemplateMock.mockRejectedValueOnce(new Error('Template API unavailable'));
    listProductsMock.mockResolvedValue([]);

    renderWeeklyTemplatePage();

    expect(await screen.findByText('Failed to load weekly plan')).toBeInTheDocument();

    getCurrentTemplateMock.mockResolvedValueOnce(makeTemplate());
    fireEvent.click(screen.getByRole('button', { name: 'Retry' }));

    expect(await screen.findByDisplayValue('Base week')).toBeInTheDocument();
    expect(getCurrentTemplateMock).toHaveBeenCalledTimes(2);
  });

  it('renders load errors and retries when product data fails', async () => {
    getCurrentTemplateMock.mockResolvedValue(makeTemplate());
    listProductsMock.mockRejectedValueOnce(new Error('Products unavailable'));

    renderWeeklyTemplatePage();

    expect(await screen.findByText('Failed to load products')).toBeInTheDocument();

    listProductsMock.mockResolvedValueOnce([makeProduct({ name: 'Buckwheat' })]);
    fireEvent.click(screen.getByRole('button', { name: 'Retry' }));

    fireEvent.click(await screen.findByRole('button', { name: 'Add planned entry' }));
    await chooseProductForEntry(1, /Buckwheat/);
    expect(screen.getAllByText(/Buckwheat/).length).toBeGreaterThan(0);
    expect(listProductsMock).toHaveBeenCalledTimes(2);
  });

  it('validates missing product and non-positive grams', async () => {
    getCurrentTemplateMock.mockResolvedValue(makeTemplate());
    listProductsMock.mockResolvedValue([makeProduct({ name: 'Rice' })]);

    renderWeeklyTemplatePage();

    fireEvent.click(await screen.findByRole('button', { name: 'Add planned entry' }));
    fireEvent.click(screen.getByRole('button', { name: 'Save Template' }));

    expect(await screen.findByText('Choose a product')).toBeInTheDocument();
    expect(updateTemplateMock).not.toHaveBeenCalled();

    await chooseProductForEntry(1, /Rice/);
    fireEvent.change(screen.getByLabelText('Grams for entry 1'), { target: { value: '0' } });
    fireEvent.click(screen.getByRole('button', { name: 'Save Template' }));

    expect(await screen.findByText('Grams must be greater than 0')).toBeInTheDocument();
    expect(updateTemplateMock).not.toHaveBeenCalled();
  });

  it('bubbles API validation and not-found errors to the editor', async () => {
    getCurrentTemplateMock.mockResolvedValue(makeTemplate());
    listProductsMock.mockResolvedValue([makeProduct()]);
    updateTemplateMock.mockRejectedValue(
      new AtlasNutritionApiError('Title is too long', 'VALIDATION_ERROR', 'validation'),
    );
    applyTemplateMock.mockRejectedValue(
      new AtlasNutritionApiError('Template missing', 'NOT_FOUND', 'not_found'),
    );

    renderWeeklyTemplatePage();

    expect(await screen.findByRole('button', { name: 'Save Template' })).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Save Template' }));

    expect(await screen.findByText('Title is too long')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'Apply to Week' }));

    expect(await screen.findByText('Template missing')).toBeInTheDocument();
  });

  it('creates a missing template and renders success after Save Template', async () => {
    getCurrentTemplateMock.mockResolvedValue(null);
    listProductsMock.mockResolvedValue([]);
    createTemplateMock.mockResolvedValue(makeTemplate({ title: 'New week', notes: null }));

    renderWeeklyTemplatePage();

    fireEvent.change(await screen.findByLabelText('Template title'), {
      target: { value: 'New week' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'Save Template' }));

    expect(await screen.findByText('Template saved')).toBeInTheDocument();
    expect(createTemplateMock).toHaveBeenCalledWith({
      weekStartDate: '2026-06-22',
      title: 'New week',
      notes: null,
    });
  });

  it('does not render a body weight field', async () => {
    getCurrentTemplateMock.mockResolvedValue(makeTemplate());
    listProductsMock.mockResolvedValue([makeProduct()]);

    renderWeeklyTemplatePage();

    expect(await screen.findByRole('heading', { name: 'Weekly Plan' })).toBeInTheDocument();
    expect(screen.queryByLabelText(/body weight/i)).not.toBeInTheDocument();
    expect(screen.queryByText(/body weight/i)).not.toBeInTheDocument();
  });

  it('renders Russian weekly plan labels when Russian is selected', async () => {
    getCurrentTemplateMock.mockResolvedValue(makeTemplate());
    listProductsMock.mockResolvedValue([makeProduct()]);

    renderWeeklyTemplatePage('ru');

    expect(await screen.findByRole('heading', { name: 'Недельный план' })).toBeInTheDocument();
    expect(await screen.findByRole('button', { name: 'Сохранить шаблон' })).toBeInTheDocument();
  });

  it('localizes unknown product fallback text and entry aria labels in Russian', async () => {
    getCurrentTemplateMock.mockResolvedValue(
      makeTemplate({
        items: [
          makeTemplateItem({
            id: 'item-1',
            productId: 'missing-product',
            mealLabel: null,
            notes: null,
          }),
        ],
      }),
    );
    listProductsMock.mockResolvedValue([]);

    renderWeeklyTemplatePage('ru');

    expect(await screen.findAllByText('Неизвестный продукт (missing-product)')).not.toHaveLength(0);
    expect(screen.getByLabelText('Продукт для записи 1')).toBeInTheDocument();
    expect(screen.getByLabelText('Граммы для записи 1')).toBeInTheDocument();
  });
});
