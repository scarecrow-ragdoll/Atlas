// FILE: apps/web-admin/src/pages/atlas/product-library-page.test.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify the Atlas nutrition product library page behavior.
//   SCOPE: Covers API-backed product loading, active/archive filtering, create/edit/archive/restore flows, validation, and page states; excludes GraphQL transport internals.
//   DEPENDS: apps/web-admin/src/pages/atlas/product-library-page.tsx, apps/web-admin/src/pages/atlas/nutrition-api.ts, apps/web-admin/src/app/i18n.tsx, @tanstack/react-query, react-router, @testing-library/react, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-API-NUTRITION.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   ProductLibraryPage tests - Prove the product library uses the Atlas nutrition API instead of mock/reference content.
// END_MODULE_MAP

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { I18nProvider, type Language } from '../../app/i18n';
import {
  archiveAtlasNutritionProduct,
  AtlasNutritionApiError,
  createAtlasNutritionProduct,
  listAtlasNutritionProducts,
  restoreAtlasNutritionProduct,
  updateAtlasNutritionProduct,
  type AtlasNutritionProduct,
} from './nutrition-api';
import ProductLibraryPage from './product-library-page';

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
    AtlasNutritionApiError: MockAtlasNutritionApiError,
    archiveAtlasNutritionProduct: vi.fn(),
    createAtlasNutritionProduct: vi.fn(),
    listAtlasNutritionProducts: vi.fn(),
    restoreAtlasNutritionProduct: vi.fn(),
    updateAtlasNutritionProduct: vi.fn(),
  };
});

const listProductsMock = vi.mocked(listAtlasNutritionProducts);
const createProductMock = vi.mocked(createAtlasNutritionProduct);
const updateProductMock = vi.mocked(updateAtlasNutritionProduct);
const archiveProductMock = vi.mocked(archiveAtlasNutritionProduct);
const restoreProductMock = vi.mocked(restoreAtlasNutritionProduct);

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

function renderProductLibraryPage(language: Language = 'en') {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <I18nProvider initialLanguage={language}>
        <MemoryRouter>
          <ProductLibraryPage />
        </MemoryRouter>
      </I18nProvider>
    </QueryClientProvider>,
  );
}

function getRowByProductName(name: string) {
  return screen.getByText(name).closest('tr') as HTMLTableRowElement;
}

async function fillProductForm(values: {
  name?: string;
  calories?: string;
  protein?: string;
  fat?: string;
  carbs?: string;
  notes?: string;
}) {
  if (values.name !== undefined) {
    fireEvent.change(screen.getByLabelText('Product name'), { target: { value: values.name } });
  }
  if (values.calories !== undefined) {
    fireEvent.change(screen.getByLabelText('Calories per 100g'), {
      target: { value: values.calories },
    });
  }
  if (values.protein !== undefined) {
    fireEvent.change(screen.getByLabelText('Protein per 100g'), {
      target: { value: values.protein },
    });
  }
  if (values.fat !== undefined) {
    fireEvent.change(screen.getByLabelText('Fat per 100g'), { target: { value: values.fat } });
  }
  if (values.carbs !== undefined) {
    fireEvent.change(screen.getByLabelText('Carbs per 100g'), { target: { value: values.carbs } });
  }
  if (values.notes !== undefined) {
    fireEvent.change(screen.getByLabelText('Notes'), { target: { value: values.notes } });
  }
}

describe('ProductLibraryPage', () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  afterEach(() => {
    cleanup();
  });

  it('loads products from the API and archives/restores a product', async () => {
    listProductsMock.mockResolvedValue([
      makeProduct({ id: 'p1', name: 'Rice', isActive: true }),
      makeProduct({ id: 'p2', name: 'Milk', isActive: false, notes: 'Dairy' }),
    ]);
    archiveProductMock.mockResolvedValue(makeProduct({ id: 'p1', name: 'Rice', isActive: false }));
    restoreProductMock.mockResolvedValue(makeProduct({ id: 'p2', name: 'Milk', isActive: true }));

    renderProductLibraryPage();

    expect(await screen.findByRole('heading', { name: 'Product Library' })).toBeInTheDocument();
    expect(listProductsMock).toHaveBeenCalledWith({ includeArchived: true });
    expect(await screen.findByText('Rice')).toBeInTheDocument();
    expect(screen.queryByText('Milk')).not.toBeInTheDocument();
    expect(screen.queryByText('Represented states')).not.toBeInTheDocument();

    fireEvent.click(within(getRowByProductName('Rice')).getByRole('button', { name: 'Archive' }));

    expect(await screen.findByText('Product archived')).toBeInTheDocument();
    expect(archiveProductMock).toHaveBeenCalledWith('p1');

    fireEvent.click(screen.getByRole('tab', { name: 'Archived' }));
    expect(screen.getByText('Milk')).toBeInTheDocument();

    fireEvent.click(within(getRowByProductName('Milk')).getByRole('button', { name: 'Restore' }));

    expect(await screen.findByText('Product restored')).toBeInTheDocument();
    expect(restoreProductMock).toHaveBeenCalledWith('p2');
  });

  it('renders the loading state while products are pending', () => {
    listProductsMock.mockReturnValue(new Promise(() => undefined));

    renderProductLibraryPage();

    expect(screen.getByRole('status', { name: 'Loading products' })).toBeInTheDocument();
  });

  it('renders an empty state when no products exist', async () => {
    listProductsMock.mockResolvedValue([]);

    renderProductLibraryPage();

    expect(await screen.findByRole('heading', { name: 'No products yet' })).toBeInTheDocument();
    expect(
      screen.getByText('Create your first product to use it in food logs.'),
    ).toBeInTheDocument();
  });

  it('renders load errors and retries the product list request', async () => {
    listProductsMock.mockRejectedValueOnce(new Error('API unavailable'));

    renderProductLibraryPage();

    expect(await screen.findByText('Failed to load products')).toBeInTheDocument();

    listProductsMock.mockResolvedValueOnce([makeProduct({ id: 'p1', name: 'Buckwheat' })]);
    fireEvent.click(screen.getByRole('button', { name: 'Retry' }));

    expect(await screen.findByText('Buckwheat')).toBeInTheDocument();
    expect(listProductsMock).toHaveBeenCalledTimes(2);
  });

  it('validates missing product names and negative macro values', async () => {
    listProductsMock.mockResolvedValue([]);

    renderProductLibraryPage();
    await screen.findByRole('heading', { name: 'No products yet' });

    fireEvent.click(screen.getByRole('button', { name: 'Save product' }));

    expect(await screen.findByText('Product name is required')).toBeInTheDocument();
    expect(createProductMock).not.toHaveBeenCalled();

    await fillProductForm({
      name: 'Invalid oats',
      calories: '-1',
      protein: '10',
      fat: '5',
      carbs: '40',
    });
    fireEvent.click(screen.getByRole('button', { name: 'Save product' }));

    expect(await screen.findByText('Macro values must be zero or greater')).toBeInTheDocument();
    expect(createProductMock).not.toHaveBeenCalled();
  });

  it('shows success states after create and edit', async () => {
    listProductsMock.mockResolvedValue([makeProduct({ id: 'p1', name: 'Rice' })]);
    createProductMock.mockResolvedValue(
      makeProduct({ id: 'p3', name: 'Eggs', proteinPer100g: 13 }),
    );
    updateProductMock.mockResolvedValue(makeProduct({ id: 'p1', name: 'Brown rice' }));

    renderProductLibraryPage();
    expect(await screen.findByText('Rice')).toBeInTheDocument();

    await fillProductForm({
      name: 'Eggs',
      calories: '155',
      protein: '13',
      fat: '11',
      carbs: '1',
      notes: 'Breakfast protein',
    });
    fireEvent.click(screen.getByRole('button', { name: 'Save product' }));

    expect(await screen.findByText('Product created')).toBeInTheDocument();
    expect(createProductMock).toHaveBeenCalledWith({
      name: 'Eggs',
      caloriesPer100g: 155,
      proteinPer100g: 13,
      fatPer100g: 11,
      carbsPer100g: 1,
      notes: 'Breakfast protein',
    });

    fireEvent.click(within(getRowByProductName('Rice')).getByRole('button', { name: 'Edit' }));
    await fillProductForm({ name: 'Brown rice' });
    fireEvent.click(screen.getByRole('button', { name: 'Update product' }));

    expect(await screen.findByText('Product updated')).toBeInTheDocument();
    expect(updateProductMock).toHaveBeenCalledWith('p1', {
      name: 'Brown rice',
      caloriesPer100g: 130,
      proteinPer100g: 2.7,
      fatPer100g: 0.3,
      carbsPer100g: 28,
      notes: 'Base carbs',
    });
  });

  it('shows API validation errors in the form', async () => {
    listProductsMock.mockResolvedValue([]);
    createProductMock.mockRejectedValue(
      new AtlasNutritionApiError('Product name already exists', 'VALIDATION_ERROR', 'validation'),
    );

    renderProductLibraryPage();
    await screen.findByRole('heading', { name: 'No products yet' });

    await fillProductForm({
      name: 'Duplicate',
      calories: '100',
      protein: '1',
      fat: '1',
      carbs: '20',
    });
    fireEvent.click(screen.getByRole('button', { name: 'Save product' }));

    expect(await screen.findByText('Product name already exists')).toBeInTheDocument();
  });

  it('renders Russian product library labels when Russian is selected', async () => {
    listProductsMock.mockResolvedValue([makeProduct({ id: 'p1', name: 'Rice' })]);

    renderProductLibraryPage('ru');

    expect(
      await screen.findByRole('heading', { name: 'Библиотека продуктов' }),
    ).toBeInTheDocument();
    expect(await screen.findByText('Rice')).toBeInTheDocument();
    expect(screen.getByRole('columnheader', { name: 'Действия' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Архивировать' })).toBeInTheDocument();
  });
});
