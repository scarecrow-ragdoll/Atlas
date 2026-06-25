// FILE: apps/web-admin/src/pages/atlas/nutrition-api.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide the Atlas nutrition page adapter for factual food logs, products, weekly templates, and typed GraphQL result errors.
//   SCOPE: Owns raw Atlas GraphQL documents, local TypeScript contracts, /graphql/atlas client construction, and result-error normalization; excludes generated admin GraphQL types, page state, and legacy override editing APIs.
//   DEPENDS: graphql-request, apps/web-admin/src/shared/config/index.ts, apps/api/internal/atlas/graph/schema/nutrition.graphql, apps/api/internal/atlas/graph/schema/daily_nutrition.graphql.
//   LINKS: M-WEB-ADMIN / M-API-NUTRITION / V-M-WEB-ADMIN / V-M-API-NUTRITION.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   createAtlasNutritionGraphQLClient - Constructs a credentialed GraphQL client for Atlas guarded endpoints.
//   getAtlasDailyNutritionLog - Loads a factual daily nutrition log by date.
//   addAtlasDailyNutritionEntry - Adds a product snapshot entry and returns the updated daily log.
//   updateAtlasDailyNutritionEntry - Replaces mutable entry fields, including required position, and returns the updated daily log.
//   deleteAtlasDailyNutritionEntry - Deletes a factual entry and returns the updated daily log.
//   listAtlasNutritionProducts - Lists active products by default or all products when archived products are requested.
//   createAtlasNutritionProduct/updateAtlasNutritionProduct/archiveAtlasNutritionProduct/restoreAtlasNutritionProduct - Manage product catalog records.
//   listAtlasNutritionTemplates/getAtlasNutritionTemplate/getAtlasNutritionTemplateCurrent/createAtlasNutritionTemplate/updateAtlasNutritionTemplate/deleteAtlasNutritionTemplate - Manage weekly template headers.
//   createAtlasNutritionTemplateItem/updateAtlasNutritionTemplateItem/deleteAtlasNutritionTemplateItem - Manage weekly template product rows.
//   applyAtlasNutritionTemplateToWeek - Seeds empty factual days from a weekly template.
//   AtlasNutritionApiError - Typed API error thrown from GraphQL result error envelopes.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Mapped empty current-template success responses to null for missing weekly templates.
// END_CHANGE_SUMMARY

import { GraphQLClient } from 'graphql-request';
import { appConfig } from '@shared/config';

export interface AtlasNutritionMacros {
  calories: number;
  protein: number;
  fat: number;
  carbs: number;
}

export interface AtlasNutritionProduct {
  id: string;
  userId: string;
  name: string;
  caloriesPer100g: number;
  proteinPer100g: number;
  fatPer100g: number;
  carbsPer100g: number;
  notes: string | null;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface AtlasNutritionTemplateItem {
  id: string;
  templateId: string;
  productId: string;
  amountGrams: number;
  mealLabel: string | null;
  notes: string | null;
  createdAt: string;
  updatedAt: string;
}

export interface AtlasNutritionTemplate {
  id: string;
  userId: string;
  weekStartDate: string;
  title: string | null;
  notes: string | null;
  items: AtlasNutritionTemplateItem[];
  createdAt: string;
  updatedAt: string;
}

export interface AtlasDailyNutritionEntry {
  id: string;
  dailyLogId: string;
  productId: string;
  productNameSnapshot: string;
  caloriesPer100gSnapshot: number;
  proteinPer100gSnapshot: number;
  fatPer100gSnapshot: number;
  carbsPer100gSnapshot: number;
  amountGrams: number;
  mealLabel: string | null;
  notes: string | null;
  position: number;
  macros: AtlasNutritionMacros;
  createdAt: string;
  updatedAt: string;
}

export interface AtlasDailyNutritionLog {
  id: string;
  userId: string;
  date: string;
  notes: string | null;
  entries: AtlasDailyNutritionEntry[];
  totals: AtlasNutritionMacros;
  createdAt: string;
  updatedAt: string;
}

export interface AtlasNutritionTemplateApplyDateResult {
  date: string;
  status: string;
  entryCount: number;
  reason: string | null;
}

export interface AtlasNutritionTemplateApplyResult {
  weekStartDate: string | null;
  weekEndDate: string | null;
  mode: 'SEED_EMPTY_DAYS' | null;
  dates: AtlasNutritionTemplateApplyDateResult[];
}

export interface AddAtlasDailyNutritionEntryInput {
  date: string;
  productId: string;
  amountGrams: number;
  mealLabel?: string | null;
  notes?: string | null;
}

export interface UpdateAtlasDailyNutritionEntryInput {
  dailyLogId: string;
  amountGrams: number;
  mealLabel?: string | null;
  notes?: string | null;
  position: number;
}

export interface CreateAtlasNutritionProductInput {
  name: string;
  caloriesPer100g: number;
  proteinPer100g: number;
  fatPer100g: number;
  carbsPer100g: number;
  notes?: string | null;
}

export interface UpdateAtlasNutritionProductInput {
  name?: string | null;
  caloriesPer100g?: number | null;
  proteinPer100g?: number | null;
  fatPer100g?: number | null;
  carbsPer100g?: number | null;
  notes?: string | null;
}

export interface ListAtlasNutritionTemplatesOptions {
  startDate: string;
  endDate: string;
}

export interface CreateAtlasNutritionTemplateInput {
  weekStartDate: string;
  title?: string | null;
  notes?: string | null;
}

export interface UpdateAtlasNutritionTemplateInput {
  title?: string | null;
  notes?: string | null;
}

export interface CreateAtlasNutritionTemplateItemInput {
  templateId: string;
  productId: string;
  amountGrams: number;
  mealLabel?: string | null;
  notes?: string | null;
}

export interface UpdateAtlasNutritionTemplateItemInput {
  amountGrams?: number | null;
  mealLabel?: string | null;
  notes?: string | null;
}

export type AtlasNutritionApiErrorType = 'validation' | 'not_found' | 'auth' | 'internal';

export class AtlasNutritionApiError extends Error {
  readonly code: string;
  readonly type: AtlasNutritionApiErrorType;

  constructor(message: string, code: string, type: AtlasNutritionApiErrorType) {
    super(message);
    this.name = 'AtlasNutritionApiError';
    this.code = code;
    this.type = type;
  }
}

interface AtlasNutritionResultError {
  message: string;
  code: string;
}

interface AtlasNutritionResultErrors {
  validationError?: AtlasNutritionResultError | null;
  notFoundError?: AtlasNutritionResultError | null;
  authError?: AtlasNutritionResultError | null;
}

type AtlasNutritionResult<TValueKey extends string, TValue> = AtlasNutritionResultErrors &
  Record<TValueKey, TValue | null | undefined>;

interface AtlasNutritionProductsResult extends AtlasNutritionResultErrors {
  products: AtlasNutritionProduct[];
}

interface AtlasNutritionTemplatesResult extends AtlasNutritionResultErrors {
  templates: AtlasNutritionTemplate[];
}

type GraphQLVariables = Record<string, unknown>;

const NUTRITION_ERROR_FIELDS = `
  validationError { message code }
  notFoundError { message code }
  authError { message code }
`;

const NUTRITION_LIST_ERROR_FIELDS = `
  validationError { message code }
  authError { message code }
`;

const MACROS_FIELDS = `
  calories
  protein
  fat
  carbs
`;

const PRODUCT_FIELDS = `
  id
  userId
  name
  caloriesPer100g
  proteinPer100g
  fatPer100g
  carbsPer100g
  notes
  isActive
  createdAt
  updatedAt
`;

const TEMPLATE_ITEM_FIELDS = `
  id
  templateId
  productId
  amountGrams
  mealLabel
  notes
  createdAt
  updatedAt
`;

const TEMPLATE_FIELDS = `
  id
  userId
  weekStartDate
  title
  notes
  items {
    ${TEMPLATE_ITEM_FIELDS}
  }
  createdAt
  updatedAt
`;

const DAILY_NUTRITION_LOG_FIELDS = `
  id
  userId
  date
  notes
  totals {
    ${MACROS_FIELDS}
  }
  entries {
    id
    dailyLogId
    productId
    productNameSnapshot
    caloriesPer100gSnapshot
    proteinPer100gSnapshot
    fatPer100gSnapshot
    carbsPer100gSnapshot
    amountGrams
    mealLabel
    notes
    position
    macros {
      ${MACROS_FIELDS}
    }
    createdAt
    updatedAt
  }
  createdAt
  updatedAt
`;

const DAILY_NUTRITION_LOG_QUERY = `
  query AtlasDailyNutritionLog($date: Date!) {
    dailyNutritionLog(date: $date) {
      dailyNutritionLog {
        ${DAILY_NUTRITION_LOG_FIELDS}
      }
      ${NUTRITION_ERROR_FIELDS}
    }
  }
`;

const ADD_DAILY_NUTRITION_ENTRY_MUTATION = `
  mutation AtlasAddDailyNutritionEntry($input: AddDailyNutritionEntryInput!) {
    addDailyNutritionEntry(input: $input) {
      dailyNutritionLog {
        ${DAILY_NUTRITION_LOG_FIELDS}
      }
      ${NUTRITION_ERROR_FIELDS}
    }
  }
`;

const UPDATE_DAILY_NUTRITION_ENTRY_MUTATION = `
  mutation AtlasUpdateDailyNutritionEntry($id: ID!, $input: UpdateDailyNutritionEntryInput!) {
    updateDailyNutritionEntry(id: $id, input: $input) {
      dailyNutritionLog {
        ${DAILY_NUTRITION_LOG_FIELDS}
      }
      ${NUTRITION_ERROR_FIELDS}
    }
  }
`;

const DELETE_DAILY_NUTRITION_ENTRY_MUTATION = `
  mutation AtlasDeleteDailyNutritionEntry($id: ID!) {
    deleteDailyNutritionEntry(id: $id) {
      dailyNutritionLog {
        ${DAILY_NUTRITION_LOG_FIELDS}
      }
      ${NUTRITION_ERROR_FIELDS}
    }
  }
`;

const NUTRITION_PRODUCTS_QUERY = `
  query AtlasNutritionProducts {
    nutritionProducts {
      products {
        ${PRODUCT_FIELDS}
      }
      ${NUTRITION_LIST_ERROR_FIELDS}
    }
  }
`;

const NUTRITION_PRODUCTS_ALL_QUERY = `
  query AtlasNutritionProductsAll {
    nutritionProductsAll {
      products {
        ${PRODUCT_FIELDS}
      }
      ${NUTRITION_LIST_ERROR_FIELDS}
    }
  }
`;

const CREATE_NUTRITION_PRODUCT_MUTATION = `
  mutation AtlasCreateNutritionProduct($input: CreateProductInput!) {
    createNutritionProduct(input: $input) {
      nutritionProduct {
        ${PRODUCT_FIELDS}
      }
      ${NUTRITION_ERROR_FIELDS}
    }
  }
`;

const UPDATE_NUTRITION_PRODUCT_MUTATION = `
  mutation AtlasUpdateNutritionProduct($id: ID!, $input: UpdateProductInput!) {
    updateNutritionProduct(id: $id, input: $input) {
      nutritionProduct {
        ${PRODUCT_FIELDS}
      }
      ${NUTRITION_ERROR_FIELDS}
    }
  }
`;

const ARCHIVE_NUTRITION_PRODUCT_MUTATION = `
  mutation AtlasArchiveNutritionProduct($id: ID!) {
    deleteNutritionProduct(id: $id) {
      nutritionProduct {
        ${PRODUCT_FIELDS}
      }
      ${NUTRITION_ERROR_FIELDS}
    }
  }
`;

const RESTORE_NUTRITION_PRODUCT_MUTATION = `
  mutation AtlasRestoreNutritionProduct($id: ID!) {
    restoreNutritionProduct(id: $id) {
      nutritionProduct {
        ${PRODUCT_FIELDS}
      }
      ${NUTRITION_ERROR_FIELDS}
    }
  }
`;

const NUTRITION_TEMPLATES_QUERY = `
  query AtlasNutritionTemplates($startDate: Date!, $endDate: Date!) {
    nutritionTemplates(startDate: $startDate, endDate: $endDate) {
      templates {
        ${TEMPLATE_FIELDS}
      }
      ${NUTRITION_LIST_ERROR_FIELDS}
    }
  }
`;

const NUTRITION_TEMPLATE_QUERY = `
  query AtlasNutritionTemplate($id: ID!) {
    nutritionTemplate(id: $id) {
      nutritionTemplate {
        ${TEMPLATE_FIELDS}
      }
      ${NUTRITION_ERROR_FIELDS}
    }
  }
`;

const NUTRITION_TEMPLATE_CURRENT_QUERY = `
  query AtlasNutritionTemplateCurrent($weekStartDate: Date!) {
    nutritionTemplateCurrent(weekStartDate: $weekStartDate) {
      nutritionTemplate {
        ${TEMPLATE_FIELDS}
      }
      ${NUTRITION_ERROR_FIELDS}
    }
  }
`;

const CREATE_NUTRITION_TEMPLATE_MUTATION = `
  mutation AtlasCreateNutritionTemplate($input: CreateTemplateInput!) {
    createNutritionTemplate(input: $input) {
      nutritionTemplate {
        ${TEMPLATE_FIELDS}
      }
      ${NUTRITION_ERROR_FIELDS}
    }
  }
`;

const UPDATE_NUTRITION_TEMPLATE_MUTATION = `
  mutation AtlasUpdateNutritionTemplate($id: ID!, $input: UpdateTemplateInput!) {
    updateNutritionTemplate(id: $id, input: $input) {
      nutritionTemplate {
        ${TEMPLATE_FIELDS}
      }
      ${NUTRITION_ERROR_FIELDS}
    }
  }
`;

const DELETE_NUTRITION_TEMPLATE_MUTATION = `
  mutation AtlasDeleteNutritionTemplate($id: ID!) {
    deleteNutritionTemplate(id: $id) {
      nutritionTemplate {
        ${TEMPLATE_FIELDS}
      }
      ${NUTRITION_ERROR_FIELDS}
    }
  }
`;

const CREATE_NUTRITION_TEMPLATE_ITEM_MUTATION = `
  mutation AtlasCreateNutritionTemplateItem($input: CreateTemplateItemInput!) {
    createNutritionTemplateItem(input: $input) {
      nutritionTemplateItem {
        ${TEMPLATE_ITEM_FIELDS}
      }
      ${NUTRITION_ERROR_FIELDS}
    }
  }
`;

const UPDATE_NUTRITION_TEMPLATE_ITEM_MUTATION = `
  mutation AtlasUpdateNutritionTemplateItem($id: ID!, $input: UpdateTemplateItemInput!) {
    updateNutritionTemplateItem(id: $id, input: $input) {
      nutritionTemplateItem {
        ${TEMPLATE_ITEM_FIELDS}
      }
      ${NUTRITION_ERROR_FIELDS}
    }
  }
`;

const DELETE_NUTRITION_TEMPLATE_ITEM_MUTATION = `
  mutation AtlasDeleteNutritionTemplateItem($id: ID!) {
    deleteNutritionTemplateItem(id: $id) {
      nutritionTemplateItem {
        ${TEMPLATE_ITEM_FIELDS}
      }
      ${NUTRITION_ERROR_FIELDS}
    }
  }
`;

const APPLY_NUTRITION_TEMPLATE_TO_WEEK_MUTATION = `
  mutation AtlasApplyNutritionTemplateToWeek($templateId: ID!, $mode: NutritionTemplateApplyMode!) {
    applyNutritionTemplateToWeek(templateId: $templateId, mode: $mode) {
      weekStartDate
      weekEndDate
      mode
      dates {
        date
        status
        entryCount
        reason
      }
      ${NUTRITION_ERROR_FIELDS}
    }
  }
`;

export function getAtlasGraphQLApiUrl(apiUrl = appConfig.apiUrl): string {
  const explicitUrl = import.meta.env.VITE_ATLAS_GRAPHQL_API_URL?.trim();

  if (explicitUrl) {
    return explicitUrl.replace(/\/+$/, '');
  }

  const normalizedAdminUrl = apiUrl.replace(/\/+$/, '');

  if (normalizedAdminUrl.endsWith('/graphql')) {
    return normalizedAdminUrl.replace(/\/graphql$/, '/graphql/atlas');
  }

  return `${normalizedAdminUrl}/graphql/atlas`;
}

export function createAtlasNutritionGraphQLClient(apiUrl = getAtlasGraphQLApiUrl()) {
  return new GraphQLClient(apiUrl, { credentials: 'include' });
}

const atlasNutritionGraphQLClient = createAtlasNutritionGraphQLClient();

async function requestAtlasNutrition<TResult>(
  document: string,
  variables: object = {},
): Promise<TResult> {
  return atlasNutritionGraphQLClient.request<TResult>(document, variables as GraphQLVariables);
}

function resultErrorToApiError(result: AtlasNutritionResultErrors | null | undefined) {
  if (result?.validationError) {
    return new AtlasNutritionApiError(
      result.validationError.message,
      result.validationError.code,
      'validation',
    );
  }

  if (result?.notFoundError) {
    return new AtlasNutritionApiError(
      result.notFoundError.message,
      result.notFoundError.code,
      'not_found',
    );
  }

  if (result?.authError) {
    return new AtlasNutritionApiError(result.authError.message, result.authError.code, 'auth');
  }

  return null;
}

function unwrapResult<TValue, TValueKey extends string>(
  result: AtlasNutritionResult<TValueKey, TValue> | null | undefined,
  valueKey: TValueKey,
  operationName: string,
): TValue {
  const apiError = resultErrorToApiError(result);

  if (apiError) {
    throw apiError;
  }

  const value = result?.[valueKey];

  if (value == null) {
    throw new AtlasNutritionApiError(
      `${operationName} returned no data and no typed error.`,
      'INTERNAL_ERROR',
      'internal',
    );
  }

  return value;
}

function unwrapNullableResult<TValue, TValueKey extends string>(
  result: Partial<AtlasNutritionResult<TValueKey, TValue>> | null | undefined,
  valueKey: TValueKey,
): TValue | null {
  const apiError = resultErrorToApiError(result);

  if (apiError) {
    throw apiError;
  }

  return result?.[valueKey] ?? null;
}

function unwrapListResult<TValue, TValueKey extends string>(
  result: (AtlasNutritionResultErrors & Record<TValueKey, TValue[]>) | null | undefined,
  valueKey: TValueKey,
  operationName: string,
): TValue[] {
  const apiError = resultErrorToApiError(result);

  if (apiError) {
    throw apiError;
  }

  const value = result?.[valueKey];

  if (!Array.isArray(value)) {
    throw new AtlasNutritionApiError(
      `${operationName} returned no list data and no typed error.`,
      'INTERNAL_ERROR',
      'internal',
    );
  }

  return value;
}

function unwrapApplyResult(
  result: (AtlasNutritionTemplateApplyResult & AtlasNutritionResultErrors) | null | undefined,
): AtlasNutritionTemplateApplyResult {
  const apiError = resultErrorToApiError(result);

  if (apiError) {
    throw apiError;
  }

  if (!result?.dates) {
    throw new AtlasNutritionApiError(
      'applyNutritionTemplateToWeek returned no date results and no typed error.',
      'INTERNAL_ERROR',
      'internal',
    );
  }

  return {
    weekStartDate: result.weekStartDate,
    weekEndDate: result.weekEndDate,
    mode: result.mode,
    dates: result.dates,
  };
}

export async function getAtlasDailyNutritionLog(date: string): Promise<AtlasDailyNutritionLog> {
  const response = await requestAtlasNutrition<{
    dailyNutritionLog: AtlasNutritionResult<'dailyNutritionLog', AtlasDailyNutritionLog>;
  }>(DAILY_NUTRITION_LOG_QUERY, { date });

  return unwrapResult(response.dailyNutritionLog, 'dailyNutritionLog', 'dailyNutritionLog');
}

export async function addAtlasDailyNutritionEntry(
  input: AddAtlasDailyNutritionEntryInput,
): Promise<AtlasDailyNutritionLog> {
  const response = await requestAtlasNutrition<{
    addDailyNutritionEntry: AtlasNutritionResult<'dailyNutritionLog', AtlasDailyNutritionLog>;
  }>(ADD_DAILY_NUTRITION_ENTRY_MUTATION, { input });

  return unwrapResult(
    response.addDailyNutritionEntry,
    'dailyNutritionLog',
    'addDailyNutritionEntry',
  );
}

export async function updateAtlasDailyNutritionEntry(
  id: string,
  input: UpdateAtlasDailyNutritionEntryInput,
): Promise<AtlasDailyNutritionLog> {
  const response = await requestAtlasNutrition<{
    updateDailyNutritionEntry: AtlasNutritionResult<'dailyNutritionLog', AtlasDailyNutritionLog>;
  }>(UPDATE_DAILY_NUTRITION_ENTRY_MUTATION, { id, input });

  return unwrapResult(
    response.updateDailyNutritionEntry,
    'dailyNutritionLog',
    'updateDailyNutritionEntry',
  );
}

export async function deleteAtlasDailyNutritionEntry(id: string): Promise<AtlasDailyNutritionLog> {
  const response = await requestAtlasNutrition<{
    deleteDailyNutritionEntry: AtlasNutritionResult<'dailyNutritionLog', AtlasDailyNutritionLog>;
  }>(DELETE_DAILY_NUTRITION_ENTRY_MUTATION, { id });

  return unwrapResult(
    response.deleteDailyNutritionEntry,
    'dailyNutritionLog',
    'deleteDailyNutritionEntry',
  );
}

export async function listAtlasNutritionProducts(options?: {
  includeArchived?: boolean;
}): Promise<AtlasNutritionProduct[]> {
  const includeArchived = options?.includeArchived === true;
  const response = await requestAtlasNutrition<{
    nutritionProducts?: AtlasNutritionProductsResult;
    nutritionProductsAll?: AtlasNutritionProductsResult;
  }>(includeArchived ? NUTRITION_PRODUCTS_ALL_QUERY : NUTRITION_PRODUCTS_QUERY);

  return unwrapListResult(
    includeArchived ? response.nutritionProductsAll : response.nutritionProducts,
    'products',
    includeArchived ? 'nutritionProductsAll' : 'nutritionProducts',
  );
}

export async function createAtlasNutritionProduct(
  input: CreateAtlasNutritionProductInput,
): Promise<AtlasNutritionProduct> {
  const response = await requestAtlasNutrition<{
    createNutritionProduct: AtlasNutritionResult<'nutritionProduct', AtlasNutritionProduct>;
  }>(CREATE_NUTRITION_PRODUCT_MUTATION, { input });

  return unwrapResult(
    response.createNutritionProduct,
    'nutritionProduct',
    'createNutritionProduct',
  );
}

export async function updateAtlasNutritionProduct(
  id: string,
  input: UpdateAtlasNutritionProductInput,
): Promise<AtlasNutritionProduct> {
  const response = await requestAtlasNutrition<{
    updateNutritionProduct: AtlasNutritionResult<'nutritionProduct', AtlasNutritionProduct>;
  }>(UPDATE_NUTRITION_PRODUCT_MUTATION, { id, input });

  return unwrapResult(
    response.updateNutritionProduct,
    'nutritionProduct',
    'updateNutritionProduct',
  );
}

export async function archiveAtlasNutritionProduct(id: string): Promise<AtlasNutritionProduct> {
  const response = await requestAtlasNutrition<{
    deleteNutritionProduct: AtlasNutritionResult<'nutritionProduct', AtlasNutritionProduct>;
  }>(ARCHIVE_NUTRITION_PRODUCT_MUTATION, { id });

  return unwrapResult(
    response.deleteNutritionProduct,
    'nutritionProduct',
    'deleteNutritionProduct',
  );
}

export async function restoreAtlasNutritionProduct(id: string): Promise<AtlasNutritionProduct> {
  const response = await requestAtlasNutrition<{
    restoreNutritionProduct: AtlasNutritionResult<'nutritionProduct', AtlasNutritionProduct>;
  }>(RESTORE_NUTRITION_PRODUCT_MUTATION, { id });

  return unwrapResult(
    response.restoreNutritionProduct,
    'nutritionProduct',
    'restoreNutritionProduct',
  );
}

export async function listAtlasNutritionTemplates(
  options: ListAtlasNutritionTemplatesOptions,
): Promise<AtlasNutritionTemplate[]> {
  const response = await requestAtlasNutrition<{
    nutritionTemplates: AtlasNutritionTemplatesResult;
  }>(NUTRITION_TEMPLATES_QUERY, options);

  return unwrapListResult(response.nutritionTemplates, 'templates', 'nutritionTemplates');
}

export async function getAtlasNutritionTemplate(id: string): Promise<AtlasNutritionTemplate> {
  const response = await requestAtlasNutrition<{
    nutritionTemplate: AtlasNutritionResult<'nutritionTemplate', AtlasNutritionTemplate>;
  }>(NUTRITION_TEMPLATE_QUERY, { id });

  return unwrapResult(response.nutritionTemplate, 'nutritionTemplate', 'nutritionTemplate');
}

export async function getAtlasNutritionTemplateCurrent(
  weekStartDate: string,
): Promise<AtlasNutritionTemplate | null> {
  const response = await requestAtlasNutrition<{
    nutritionTemplateCurrent: Partial<
      AtlasNutritionResult<'nutritionTemplate', AtlasNutritionTemplate>
    >;
  }>(NUTRITION_TEMPLATE_CURRENT_QUERY, { weekStartDate });

  return unwrapNullableResult(response.nutritionTemplateCurrent, 'nutritionTemplate');
}

export async function createAtlasNutritionTemplate(
  input: CreateAtlasNutritionTemplateInput,
): Promise<AtlasNutritionTemplate> {
  const response = await requestAtlasNutrition<{
    createNutritionTemplate: AtlasNutritionResult<'nutritionTemplate', AtlasNutritionTemplate>;
  }>(CREATE_NUTRITION_TEMPLATE_MUTATION, { input });

  return unwrapResult(
    response.createNutritionTemplate,
    'nutritionTemplate',
    'createNutritionTemplate',
  );
}

export async function updateAtlasNutritionTemplate(
  id: string,
  input: UpdateAtlasNutritionTemplateInput,
): Promise<AtlasNutritionTemplate> {
  const response = await requestAtlasNutrition<{
    updateNutritionTemplate: AtlasNutritionResult<'nutritionTemplate', AtlasNutritionTemplate>;
  }>(UPDATE_NUTRITION_TEMPLATE_MUTATION, { id, input });

  return unwrapResult(
    response.updateNutritionTemplate,
    'nutritionTemplate',
    'updateNutritionTemplate',
  );
}

export async function deleteAtlasNutritionTemplate(id: string): Promise<AtlasNutritionTemplate> {
  const response = await requestAtlasNutrition<{
    deleteNutritionTemplate: AtlasNutritionResult<'nutritionTemplate', AtlasNutritionTemplate>;
  }>(DELETE_NUTRITION_TEMPLATE_MUTATION, { id });

  return unwrapResult(
    response.deleteNutritionTemplate,
    'nutritionTemplate',
    'deleteNutritionTemplate',
  );
}

export async function createAtlasNutritionTemplateItem(
  input: CreateAtlasNutritionTemplateItemInput,
): Promise<AtlasNutritionTemplateItem> {
  const response = await requestAtlasNutrition<{
    createNutritionTemplateItem: AtlasNutritionResult<
      'nutritionTemplateItem',
      AtlasNutritionTemplateItem
    >;
  }>(CREATE_NUTRITION_TEMPLATE_ITEM_MUTATION, { input });

  return unwrapResult(
    response.createNutritionTemplateItem,
    'nutritionTemplateItem',
    'createNutritionTemplateItem',
  );
}

export async function updateAtlasNutritionTemplateItem(
  id: string,
  input: UpdateAtlasNutritionTemplateItemInput,
): Promise<AtlasNutritionTemplateItem> {
  const response = await requestAtlasNutrition<{
    updateNutritionTemplateItem: AtlasNutritionResult<
      'nutritionTemplateItem',
      AtlasNutritionTemplateItem
    >;
  }>(UPDATE_NUTRITION_TEMPLATE_ITEM_MUTATION, { id, input });

  return unwrapResult(
    response.updateNutritionTemplateItem,
    'nutritionTemplateItem',
    'updateNutritionTemplateItem',
  );
}

export async function deleteAtlasNutritionTemplateItem(
  id: string,
): Promise<AtlasNutritionTemplateItem> {
  const response = await requestAtlasNutrition<{
    deleteNutritionTemplateItem: AtlasNutritionResult<
      'nutritionTemplateItem',
      AtlasNutritionTemplateItem
    >;
  }>(DELETE_NUTRITION_TEMPLATE_ITEM_MUTATION, { id });

  return unwrapResult(
    response.deleteNutritionTemplateItem,
    'nutritionTemplateItem',
    'deleteNutritionTemplateItem',
  );
}

export async function applyAtlasNutritionTemplateToWeek(
  templateId: string,
  mode: 'SEED_EMPTY_DAYS',
): Promise<AtlasNutritionTemplateApplyResult> {
  const response = await requestAtlasNutrition<{
    applyNutritionTemplateToWeek: AtlasNutritionTemplateApplyResult & AtlasNutritionResultErrors;
  }>(APPLY_NUTRITION_TEMPLATE_TO_WEEK_MUTATION, { templateId, mode });

  return unwrapApplyResult(response.applyNutritionTemplateToWeek);
}
