// FILE: apps/web-admin/src/pages/atlas/ai-export-api.test.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Prove the Atlas AI export REST adapter sends guarded local requests and maps backend responses safely.
//   SCOPE: Covers generate request credentials, response mapping, local download URL construction, and typed backend error mapping; excludes page rendering and backend handler execution.
//   DEPENDS: vitest, apps/web-admin/src/pages/atlas/ai-export-api.ts.
//   LINKS: M-WEB-ADMIN / M-API / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   ai export api adapter tests - Verify local REST endpoint usage, credentials include, safe download URLs, and error normalization.
// END_MODULE_MAP

import { beforeEach, describe, expect, it, vi } from 'vitest';
import {
  AtlasAiExportApiError,
  generateAtlasAiExport,
  getAtlasAiExportDownloadUrl,
  getAtlasRestApiUrl,
} from './ai-export-api';

const fetchMock = vi.fn();

function jsonResponse(body: unknown, init: ResponseInit = {}) {
  return new Response(JSON.stringify(body), {
    headers: { 'Content-Type': 'application/json' },
    ...init,
  });
}

describe('Atlas AI export REST adapter', () => {
  beforeEach(() => {
    fetchMock.mockReset();
    vi.stubGlobal('fetch', fetchMock);
  });

  it('derives local REST endpoint roots from appConfig GraphQL URLs', () => {
    expect(getAtlasRestApiUrl('http://localhost:8090/graphql')).toBe('http://localhost:8090');
    expect(getAtlasRestApiUrl('http://localhost:8090/graphql/atlas')).toBe('http://localhost:8090');
    expect(getAtlasRestApiUrl('http://localhost:8090/api')).toBe('http://localhost:8090/api');
  });

  it('generates a local AI export with credentials and maps backend success safely', async () => {
    fetchMock.mockResolvedValueOnce(
      jsonResponse({
        export: {
          id: 'export-1',
          dateRangeStart: '2026-06-01',
          dateRangeEnd: '2026-06-24',
          includePhotos: false,
          includeNutrition: true,
          includeCardio: true,
          includeMeasurements: true,
          userComment: 'Prefer concise coaching.',
          generatedPrompt: 'Analyze the nutrition and cardio trend.',
          exportFilePath: '/srv/atlas/private/user-1/export-1.zip',
          createdAt: '2026-06-24T12:00:00Z',
        },
      }),
    );

    const result = await generateAtlasAiExport(
      {
        dateRangeStart: '2026-06-01',
        dateRangeEnd: '2026-06-24',
        includePhotos: false,
        includeNutrition: true,
        includeCardio: true,
        includeMeasurements: true,
        userComment: 'Prefer concise coaching.',
      },
      'http://localhost:8090/graphql',
    );

    expect(fetchMock).toHaveBeenCalledWith('http://localhost:8090/api/ai-export/generate', {
      body: JSON.stringify({
        dateRangeStart: '2026-06-01',
        dateRangeEnd: '2026-06-24',
        includePhotos: false,
        includeNutrition: true,
        includeCardio: true,
        includeMeasurements: true,
        userComment: 'Prefer concise coaching.',
      }),
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      method: 'POST',
    });
    expect(result.export.id).toBe('export-1');
    expect(result.export.generatedPrompt).toBe('Analyze the nutrition and cardio trend.');
    expect(result.export.downloadUrl).toBe(
      'http://localhost:8090/api/ai-export/download?exportId=export-1',
    );
    expect(JSON.stringify(result)).not.toContain('/srv/atlas/private');
  });

  it('builds encoded local download URLs without server file paths', () => {
    const url = getAtlasAiExportDownloadUrl('export/with spaces', 'http://localhost:8090/graphql');

    expect(url).toBe(
      'http://localhost:8090/api/ai-export/download?exportId=export%2Fwith%20spaces',
    );
    expect(url).not.toContain('.zip');
    expect(url).not.toContain('/srv/');
  });

  it('maps backend error envelopes into typed frontend errors', async () => {
    fetchMock.mockResolvedValueOnce(
      jsonResponse(
        {
          error: {
            code: 'INVALID_DATE_RANGE',
            message: 'invalid dateRangeStart',
          },
        },
        { status: 400 },
      ),
    );

    let caughtError: unknown;
    try {
      await generateAtlasAiExport(
        {
          dateRangeStart: 'bad-date',
          dateRangeEnd: '2026-06-24',
          includePhotos: false,
          includeNutrition: true,
          includeCardio: false,
          includeMeasurements: false,
          userComment: null,
        },
        'http://localhost:8090/graphql',
      );
    } catch (error) {
      caughtError = error;
    }

    expect(caughtError).toBeInstanceOf(AtlasAiExportApiError);
    expect(caughtError).toMatchObject({
      code: 'INVALID_DATE_RANGE',
      message: 'invalid dateRangeStart',
      type: 'validation',
    });
  });
});
