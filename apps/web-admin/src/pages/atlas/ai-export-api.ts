// FILE: apps/web-admin/src/pages/atlas/ai-export-api.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide the Atlas AI export REST adapter for guarded local prompt and ZIP generation.
//   SCOPE: Owns POST /api/ai-export/generate requests, safe GET /api/ai-export/download URL construction, and typed REST error normalization; excludes page state, external AI APIs, and server file paths.
//   DEPENDS: apps/web-admin/src/shared/config/index.ts, apps/api/internal/handler/ai_export_handler.go.
//   LINKS: M-WEB-ADMIN / M-API / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   getAtlasRestApiUrl - Derives the local Atlas REST origin from appConfig GraphQL URLs.
//   getAtlasAiExportDownloadUrl - Builds a session-guarded download URL using exportId only.
//   generateAtlasAiExport - Sends the local generate request with browser credentials and maps safe response fields.
//   AtlasAiExportApiError - Typed frontend error for backend AI export error envelopes.
// END_MODULE_MAP

import { appConfig } from '@shared/config';

export interface GenerateAtlasAiExportInput {
  dateRangeStart: string;
  dateRangeEnd: string;
  includePhotos: boolean;
  includeNutrition: boolean;
  includeCardio: boolean;
  includeMeasurements: boolean;
  userComment: string | null;
}

export interface AtlasAiExport {
  id: string;
  dateRangeStart: string;
  dateRangeEnd: string;
  includePhotos: boolean;
  includeNutrition: boolean;
  includeCardio: boolean;
  includeMeasurements: boolean;
  userComment: string | null;
  generatedPrompt: string;
  createdAt: string;
  downloadUrl: string;
}

export interface GenerateAtlasAiExportResult {
  export: AtlasAiExport;
}

export type AtlasAiExportApiErrorType = 'validation' | 'not_found' | 'auth' | 'internal';

export class AtlasAiExportApiError extends Error {
  readonly code: string;
  readonly type: AtlasAiExportApiErrorType;

  constructor(message: string, code: string, type: AtlasAiExportApiErrorType) {
    super(message);
    this.name = 'AtlasAiExportApiError';
    this.code = code;
    this.type = type;
  }
}

interface BackendAiExport {
  id: string;
  dateRangeStart: string;
  dateRangeEnd: string;
  includePhotos: boolean;
  includeNutrition: boolean;
  includeCardio: boolean;
  includeMeasurements: boolean;
  userComment?: string | null;
  generatedPrompt: string;
  exportFilePath?: string | null;
  createdAt: string;
}

interface BackendAiExportError {
  code: string;
  message: string;
}

interface BackendGenerateResponse {
  export?: BackendAiExport | null;
  error?: BackendAiExportError | null;
}

function trimTrailingSlashes(value: string) {
  return value.replace(/\/+$/, '');
}

function stripGraphQLEndpoint(apiUrl: string) {
  const normalizedUrl = trimTrailingSlashes(apiUrl.trim());

  if (normalizedUrl.endsWith('/graphql/atlas')) {
    return normalizedUrl.slice(0, -'/graphql/atlas'.length);
  }

  if (normalizedUrl.endsWith('/graphql')) {
    return normalizedUrl.slice(0, -'/graphql'.length);
  }

  return normalizedUrl;
}

function joinRestPath(restApiUrl: string, path: string) {
  const normalizedRestApiUrl = trimTrailingSlashes(restApiUrl);
  return normalizedRestApiUrl ? `${normalizedRestApiUrl}${path}` : path;
}

export function getAtlasRestApiUrl(apiUrl = appConfig.apiUrl): string {
  return stripGraphQLEndpoint(apiUrl);
}

export function getAtlasAiExportDownloadUrl(
  exportId: string,
  apiUrl = getAtlasRestApiUrl(),
): string {
  return joinRestPath(
    getAtlasRestApiUrl(apiUrl),
    `/api/ai-export/download?exportId=${encodeURIComponent(exportId)}`,
  );
}

function errorTypeFromStatus(status: number): AtlasAiExportApiErrorType {
  if (status === 400) {
    return 'validation';
  }
  if (status === 401 || status === 403) {
    return 'auth';
  }
  if (status === 404) {
    return 'not_found';
  }
  return 'internal';
}

async function parseJson(response: Response): Promise<BackendGenerateResponse> {
  try {
    return (await response.json()) as BackendGenerateResponse;
  } catch {
    return {};
  }
}

function backendErrorToApiError(
  error: BackendAiExportError | null | undefined,
  status: number,
): AtlasAiExportApiError {
  return new AtlasAiExportApiError(
    error?.message || 'AI export request failed',
    error?.code || 'AI_EXPORT_REQUEST_FAILED',
    errorTypeFromStatus(status),
  );
}

function mapBackendExport(exportResult: BackendAiExport, restApiUrl: string): AtlasAiExport {
  return {
    id: exportResult.id,
    dateRangeStart: exportResult.dateRangeStart,
    dateRangeEnd: exportResult.dateRangeEnd,
    includePhotos: exportResult.includePhotos,
    includeNutrition: exportResult.includeNutrition,
    includeCardio: exportResult.includeCardio,
    includeMeasurements: exportResult.includeMeasurements,
    userComment: exportResult.userComment ?? null,
    generatedPrompt: exportResult.generatedPrompt,
    createdAt: exportResult.createdAt,
    downloadUrl: getAtlasAiExportDownloadUrl(exportResult.id, restApiUrl),
  };
}

// START_CONTRACT: generateAtlasAiExport
//   PURPOSE: Generate an AI-ready local export package and prompt through the guarded Atlas REST endpoint.
//   INPUTS: { input: GenerateAtlasAiExportInput - date range and section flags, apiUrl?: string - GraphQL or REST base override for tests }
//   OUTPUTS: { Promise<GenerateAtlasAiExportResult> - safe export fields plus local download URL }
//   SIDE_EFFECTS: Sends POST /api/ai-export/generate with credentials include.
//   LINKS: M-WEB-ADMIN / M-API / V-M-WEB-ADMIN.
// END_CONTRACT: generateAtlasAiExport
export async function generateAtlasAiExport(
  input: GenerateAtlasAiExportInput,
  apiUrl = getAtlasRestApiUrl(),
): Promise<GenerateAtlasAiExportResult> {
  const restApiUrl = getAtlasRestApiUrl(apiUrl);
  const response = await fetch(joinRestPath(restApiUrl, '/api/ai-export/generate'), {
    body: JSON.stringify(input),
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    method: 'POST',
  });
  const body = await parseJson(response);

  if (!response.ok || body.error) {
    throw backendErrorToApiError(body.error, response.status);
  }

  if (!body.export) {
    throw new AtlasAiExportApiError(
      'AI export response did not include an export.',
      'AI_EXPORT_EMPTY_RESPONSE',
      'internal',
    );
  }

  return { export: mapBackendExport(body.export, restApiUrl) };
}
