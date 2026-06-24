// FILE: apps/web-admin/src/pages/atlas/nutrition-api.test.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Prove the Atlas nutrition API adapter sends factual food-log GraphQL operations and normalizes result errors.
//   SCOPE: Covers adapter request documents, variables, daily log/product/template apply mapping, and typed error bubbling; excludes page rendering and backend resolver execution.
//   DEPENDS: vitest, graphql-request, apps/web-admin/src/pages/atlas/nutrition-api.ts.
//   LINKS: M-WEB-ADMIN / M-API-NUTRITION / V-M-WEB-ADMIN / V-M-API-NUTRITION.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   nutrition api adapter tests - Verify raw Atlas GraphQL calls, response mapping, and normalized validation/not-found/auth errors.
// END_MODULE_MAP

import { beforeEach, describe, expect, it, vi } from 'vitest';

const requestMock = vi.hoisted(() => vi.fn());
const graphQLClientMock = vi.hoisted(() =>
  vi.fn(() => ({
    request: requestMock,
  })),
);

vi.mock('graphql-request', () => ({
  GraphQLClient: graphQLClientMock,
}));

import {
  AtlasNutritionApiError,
  addAtlasDailyNutritionEntry,
  applyAtlasNutritionTemplateToWeek,
  createAtlasNutritionGraphQLClient,
  createAtlasNutritionTemplateItem,
  deleteAtlasDailyNutritionEntry,
  deleteAtlasNutritionTemplateItem,
  getAtlasGraphQLApiUrl,
  getAtlasDailyNutritionLog,
  listAtlasNutritionProducts,
  restoreAtlasNutritionProduct,
  updateAtlasDailyNutritionEntry,
  updateAtlasNutritionTemplateItem,
} from './nutrition-api';

const riceProduct = {
  id: 'product-1',
  userId: 'user-1',
  name: 'Rice',
  caloriesPer100g: 130,
  proteinPer100g: 2.7,
  fatPer100g: 0.3,
  carbsPer100g: 28,
  notes: 'white rice',
  isActive: true,
  createdAt: '2026-06-24T10:00:00Z',
  updatedAt: '2026-06-24T10:00:00Z',
};

const dailyLog = {
  id: 'log-1',
  userId: 'user-1',
  date: '2026-06-24',
  notes: 'normal day',
  totals: { calories: 325, protein: 6.75, fat: 0.75, carbs: 70 },
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
      notes: null,
      position: 1,
      macros: { calories: 325, protein: 6.75, fat: 0.75, carbs: 70 },
      createdAt: '2026-06-24T10:00:00Z',
      updatedAt: '2026-06-24T10:00:00Z',
    },
  ],
  createdAt: '2026-06-24T10:00:00Z',
  updatedAt: '2026-06-24T10:00:00Z',
};

const templateItem = {
  id: 'item-1',
  templateId: 'template-1',
  productId: 'product-1',
  amountGrams: 250,
  mealLabel: 'Lunch',
  notes: 'planned rice',
  createdAt: '2026-06-24T10:00:00Z',
  updatedAt: '2026-06-24T10:00:00Z',
};

function expectDocumentToContainFields(document: string, fields: string[]) {
  for (const field of fields) {
    expect(document).toContain(field);
  }
}

describe('Atlas nutrition API adapter', () => {
  beforeEach(() => {
    requestMock.mockReset();
    graphQLClientMock.mockClear();
    vi.unstubAllEnvs();
  });

  it('derives the Atlas GraphQL endpoint from admin config and explicit env overrides', () => {
    vi.stubEnv('VITE_ATLAS_GRAPHQL_API_URL', '');

    expect(getAtlasGraphQLApiUrl('http://localhost:8090/graphql')).toBe(
      'http://localhost:8090/graphql/atlas',
    );
    expect(getAtlasGraphQLApiUrl('http://localhost:8090/graphql///')).toBe(
      'http://localhost:8090/graphql/atlas',
    );
    expect(getAtlasGraphQLApiUrl('http://localhost:8090/api')).toBe(
      'http://localhost:8090/api/graphql/atlas',
    );

    vi.stubEnv('VITE_ATLAS_GRAPHQL_API_URL', ' https://atlas.test/graphql/atlas/// ');

    expect(getAtlasGraphQLApiUrl('http://localhost:8090/graphql')).toBe(
      'https://atlas.test/graphql/atlas',
    );
  });

  it('creates a credentialed Atlas GraphQL client', () => {
    createAtlasNutritionGraphQLClient('https://atlas.test/graphql/atlas');

    expect(graphQLClientMock).toHaveBeenCalledWith('https://atlas.test/graphql/atlas', {
      credentials: 'include',
    });
  });

  it('loads a factual daily nutrition log with entry snapshots and totals', async () => {
    requestMock.mockResolvedValueOnce({
      dailyNutritionLog: {
        dailyNutritionLog: dailyLog,
      },
    });

    const result = await getAtlasDailyNutritionLog('2026-06-24');

    expect(requestMock).toHaveBeenCalledWith(expect.stringContaining('dailyNutritionLog'), {
      date: '2026-06-24',
    });
    expectDocumentToContainFields(requestMock.mock.calls[0][0], [
      'productNameSnapshot',
      'caloriesPer100gSnapshot',
      'proteinPer100gSnapshot',
      'fatPer100gSnapshot',
      'carbsPer100gSnapshot',
      'amountGrams',
      'macros',
      'totals',
    ]);
    expect(result.entries[0].productNameSnapshot).toBe('Rice');
    expect(result.entries[0].macros.calories).toBe(325);
    expect(result.totals.calories).toBe(325);
  });

  it('adds a factual daily nutrition entry by date/product/grams and returns the updated daily log', async () => {
    const input = {
      date: '2026-06-24',
      productId: 'product-1',
      amountGrams: 250,
      mealLabel: 'Lunch',
      notes: 'post training',
    };
    requestMock.mockResolvedValueOnce({
      addDailyNutritionEntry: {
        dailyNutritionLog: dailyLog,
      },
    });

    const result = await addAtlasDailyNutritionEntry(input);

    expect(requestMock).toHaveBeenCalledWith(expect.stringContaining('addDailyNutritionEntry'), {
      input,
    });
    expect(result.id).toBe('log-1');
    expect(result.entries[0].amountGrams).toBe(250);
  });

  it('updates a factual daily nutrition entry with required position and returns the updated daily log', async () => {
    const input = {
      dailyLogId: 'log-1',
      amountGrams: 300,
      mealLabel: 'Dinner',
      notes: null,
      position: 2,
    };
    requestMock.mockResolvedValueOnce({
      updateDailyNutritionEntry: {
        dailyNutritionLog: {
          ...dailyLog,
          entries: [{ ...dailyLog.entries[0], amountGrams: 300, position: 2 }],
        },
      },
    });

    const result = await updateAtlasDailyNutritionEntry('entry-1', input);

    expect(requestMock).toHaveBeenCalledWith(expect.stringContaining('updateDailyNutritionEntry'), {
      id: 'entry-1',
      input,
    });
    expect(result.entries[0].amountGrams).toBe(300);
    expect(result.entries[0].position).toBe(2);
  });

  it('deletes a factual daily nutrition entry and returns the updated daily log', async () => {
    requestMock.mockResolvedValueOnce({
      deleteDailyNutritionEntry: {
        dailyNutritionLog: { ...dailyLog, entries: [] },
      },
    });

    const result = await deleteAtlasDailyNutritionEntry('entry-1');

    expect(requestMock).toHaveBeenCalledWith(expect.stringContaining('deleteDailyNutritionEntry'), {
      id: 'entry-1',
    });
    expect(result.entries).toEqual([]);
  });

  it('lists active nutrition products by default and all products when archived products are requested', async () => {
    requestMock
      .mockResolvedValueOnce({
        nutritionProducts: {
          products: [riceProduct],
        },
      })
      .mockResolvedValueOnce({
        nutritionProductsAll: {
          products: [riceProduct, { ...riceProduct, id: 'product-2', isActive: false }],
        },
      });

    const activeProducts = await listAtlasNutritionProducts();
    const allProducts = await listAtlasNutritionProducts({ includeArchived: true });

    expect(requestMock).toHaveBeenNthCalledWith(
      1,
      expect.stringContaining('nutritionProducts'),
      {},
    );
    expectDocumentToContainFields(requestMock.mock.calls[0][0], [
      'caloriesPer100g',
      'proteinPer100g',
      'fatPer100g',
      'carbsPer100g',
      'isActive',
    ]);
    expect(requestMock.mock.calls[0][0]).not.toContain('nutritionProductsAll');
    expect(requestMock).toHaveBeenNthCalledWith(
      2,
      expect.stringContaining('nutritionProductsAll'),
      {},
    );
    expect(activeProducts).toHaveLength(1);
    expect(allProducts).toHaveLength(2);
    expect(allProducts[1].isActive).toBe(false);
  });

  it('restores an archived nutrition product through the restore mutation', async () => {
    requestMock.mockResolvedValueOnce({
      restoreNutritionProduct: {
        nutritionProduct: { ...riceProduct, id: 'product-2', isActive: true },
      },
    });

    const result = await restoreAtlasNutritionProduct('product-2');

    expect(requestMock).toHaveBeenCalledWith(expect.stringContaining('restoreNutritionProduct'), {
      id: 'product-2',
    });
    expect(result.id).toBe('product-2');
    expect(result.isActive).toBe(true);
  });

  it('applies a nutrition template to the week and maps per-date statuses', async () => {
    requestMock.mockResolvedValueOnce({
      applyNutritionTemplateToWeek: {
        weekStartDate: '2026-06-22',
        weekEndDate: '2026-06-28',
        mode: 'SEED_EMPTY_DAYS',
        dates: [
          { date: '2026-06-22', status: 'created', entryCount: 2, reason: null },
          { date: '2026-06-23', status: 'skipped', entryCount: 1, reason: 'day has entries' },
        ],
      },
    });

    const result = await applyAtlasNutritionTemplateToWeek('template-1', 'SEED_EMPTY_DAYS');

    expect(requestMock).toHaveBeenCalledWith(
      expect.stringContaining('applyNutritionTemplateToWeek'),
      {
        templateId: 'template-1',
        mode: 'SEED_EMPTY_DAYS',
      },
    );
    expectDocumentToContainFields(requestMock.mock.calls[0][0], [
      'weekStartDate',
      'weekEndDate',
      'mode',
      'date',
      'status',
      'entryCount',
      'reason',
    ]);
    expect(result.mode).toBe('SEED_EMPTY_DAYS');
    expect(result.dates).toEqual([
      { date: '2026-06-22', status: 'created', entryCount: 2, reason: null },
      { date: '2026-06-23', status: 'skipped', entryCount: 1, reason: 'day has entries' },
    ]);
  });

  it('creates, updates, and deletes weekly template product rows', async () => {
    const createInput = {
      templateId: 'template-1',
      productId: 'product-1',
      amountGrams: 250,
      mealLabel: 'Lunch',
      notes: 'planned rice',
    };
    const updateInput = {
      amountGrams: 300,
      mealLabel: 'Dinner',
      notes: null,
    };
    requestMock
      .mockResolvedValueOnce({
        createNutritionTemplateItem: {
          nutritionTemplateItem: templateItem,
        },
      })
      .mockResolvedValueOnce({
        updateNutritionTemplateItem: {
          nutritionTemplateItem: { ...templateItem, ...updateInput },
        },
      })
      .mockResolvedValueOnce({
        deleteNutritionTemplateItem: {
          nutritionTemplateItem: { ...templateItem, ...updateInput },
        },
      });

    const created = await createAtlasNutritionTemplateItem(createInput);
    const updated = await updateAtlasNutritionTemplateItem('item-1', updateInput);
    const deleted = await deleteAtlasNutritionTemplateItem('item-1');

    expect(requestMock).toHaveBeenNthCalledWith(
      1,
      expect.stringContaining('createNutritionTemplateItem'),
      { input: createInput },
    );
    expectDocumentToContainFields(requestMock.mock.calls[0][0], [
      'templateId',
      'productId',
      'amountGrams',
      'mealLabel',
      'notes',
    ]);
    expect(requestMock).toHaveBeenNthCalledWith(
      2,
      expect.stringContaining('updateNutritionTemplateItem'),
      { id: 'item-1', input: updateInput },
    );
    expectDocumentToContainFields(requestMock.mock.calls[1][0], [
      'templateId',
      'productId',
      'amountGrams',
      'mealLabel',
      'notes',
    ]);
    expect(requestMock).toHaveBeenNthCalledWith(
      3,
      expect.stringContaining('deleteNutritionTemplateItem'),
      { id: 'item-1' },
    );
    expectDocumentToContainFields(requestMock.mock.calls[2][0], [
      'templateId',
      'productId',
      'amountGrams',
      'mealLabel',
      'notes',
    ]);
    expect(created.productId).toBe('product-1');
    expect(updated.amountGrams).toBe(300);
    expect(deleted.id).toBe('item-1');
  });

  it.each([
    ['validationError', 'validation', 'VALIDATION_ERROR'],
    ['notFoundError', 'not_found', 'NOT_FOUND'],
    ['authError', 'auth', 'AUTH_ERROR'],
  ] as const)(
    'normalizes %s result errors into AtlasNutritionApiError',
    async (field, type, code) => {
      requestMock.mockResolvedValueOnce({
        dailyNutritionLog: {
          dailyNutritionLog: null,
          [field]: {
            message: `${code} from backend`,
            code,
          },
        },
      });

      const result = getAtlasDailyNutritionLog('2026-06-24');

      await expect(result).rejects.toMatchObject({
        name: 'AtlasNutritionApiError',
        type,
        code,
        message: `${code} from backend`,
      });

      await expect(result).rejects.toBeInstanceOf(AtlasNutritionApiError);
    },
  );
});
