// FILE: apps/web-admin/src/pages/atlas/ai-export-builder-page.test.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify the Atlas AI export builder page uses the local REST adapter and renders real export states.
//   SCOPE: Covers section toggles, local privacy warning, loading/progress, error/retry, ready prompt/download states, and absence of reference/mock UI.
//   DEPENDS: apps/web-admin/src/pages/atlas/ai-export-builder-page.tsx, apps/web-admin/src/pages/atlas/ai-export-api.ts, apps/web-admin/src/app/i18n.tsx, @tanstack/react-query, @testing-library/react, vitest.
//   LINKS: M-WEB-ADMIN / M-API / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AI export builder page tests - Prove the real local export UI state machine and safe download behavior.
// END_MODULE_MAP

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterAll, afterEach, beforeAll, beforeEach, describe, expect, it, vi } from 'vitest';
import { I18nProvider, type Language } from '../../app/i18n';
import { generateAtlasAiExport } from './ai-export-api';
import AiExportBuilderPage from './ai-export-builder-page';

vi.mock('./ai-export-api', () => ({
  AtlasAiExportApiError: class MockAtlasAiExportApiError extends Error {
    readonly code: string;
    readonly type: string;

    constructor(message: string, code: string, type: string) {
      super(message);
      this.name = 'AtlasAiExportApiError';
      this.code = code;
      this.type = type;
    }
  },
  generateAtlasAiExport: vi.fn(),
  getAtlasAiExportDownloadUrl: (exportId: string) =>
    `http://localhost:8090/api/ai-export/download?exportId=${encodeURIComponent(exportId)}`,
}));

const generateMock = vi.mocked(generateAtlasAiExport);
const originalResizeObserver = globalThis.ResizeObserver;

const readyExport = {
  export: {
    id: 'export-1',
    dateRangeStart: '2026-06-01',
    dateRangeEnd: '2026-06-24',
    includePhotos: false,
    includeNutrition: true,
    includeCardio: true,
    includeMeasurements: true,
    userComment: null,
    generatedPrompt: 'Analyze nutrition adherence and cardio trend.',
    createdAt: '2026-06-24T12:00:00Z',
    downloadUrl: 'http://localhost:8090/api/ai-export/download?exportId=export-1',
  },
};

function createQueryClient() {
  return new QueryClient({
    defaultOptions: { mutations: { retry: false }, queries: { retry: false } },
  });
}

function renderAiExportPage(language: Language = 'en') {
  return render(
    <QueryClientProvider client={createQueryClient()}>
      <I18nProvider initialLanguage={language}>
        <AiExportBuilderPage initialStartDate="2026-06-01" initialEndDate="2026-06-24" />
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

class TestResizeObserver {
  observe() {
    return undefined;
  }

  unobserve() {
    return undefined;
  }

  disconnect() {
    return undefined;
  }
}

describe('AiExportBuilderPage', () => {
  beforeAll(() => {
    globalThis.ResizeObserver = TestResizeObserver as typeof ResizeObserver;
  });

  beforeEach(() => {
    vi.resetAllMocks();
  });

  afterEach(() => {
    cleanup();
  });

  afterAll(() => {
    if (originalResizeObserver) {
      globalThis.ResizeObserver = originalResizeObserver;
    } else {
      delete (globalThis as Partial<typeof globalThis>).ResizeObserver;
    }
  });

  it('generates a local AI export with nutrition enabled and shows the ready state', async () => {
    generateMock.mockResolvedValueOnce(readyExport);

    renderAiExportPage();

    expect(screen.getByRole('heading', { name: 'AI Export' })).toBeInTheDocument();
    expect(screen.getByRole('checkbox', { name: 'Nutrition' })).toBeChecked();
    fireEvent.click(screen.getByRole('button', { name: 'Generate export' }));

    await waitFor(() => {
      expect(generateMock).toHaveBeenCalledWith({
        dateRangeStart: '2026-06-01',
        dateRangeEnd: '2026-06-24',
        includePhotos: false,
        includeNutrition: true,
        includeCardio: true,
        includeMeasurements: true,
        userComment: null,
      });
    });
    expect(await screen.findByRole('heading', { name: 'Export ready' })).toBeInTheDocument();
    expect(screen.getByText('Analyze nutrition adherence and cardio trend.')).toBeInTheDocument();
  });

  it('renders a real local/internal privacy warning and no reference/mock page text', () => {
    renderAiExportPage();

    expect(
      screen.getByText(
        /This export is local and internal\. Atlas does not call external AI APIs\./i,
      ),
    ).toBeInTheDocument();
    expect(screen.getByText(/Photos are excluded unless selected\./i)).toBeInTheDocument();
    expect(screen.queryByText(/AtlasReferencePage/i)).not.toBeInTheDocument();
    expect(screen.queryByText(/mock/i)).not.toBeInTheDocument();
    expect(screen.queryByText(/reference/i)).not.toBeInTheDocument();
  });

  it('shows loading and progress while generating', async () => {
    const deferred = createDeferred<typeof readyExport>();
    generateMock.mockReturnValueOnce(deferred.promise);

    renderAiExportPage();
    fireEvent.click(screen.getByRole('button', { name: 'Generate export' }));

    expect(await screen.findByRole('button', { name: 'Generating export' })).toBeDisabled();
    expect(screen.getByText('Preparing local ZIP and prompt preview.')).toBeInTheDocument();

    await act(async () => {
      deferred.resolve(readyExport);
    });

    expect(await screen.findByRole('heading', { name: 'Export ready' })).toBeInTheDocument();
  });

  it('shows error and retries the local generate request', async () => {
    generateMock.mockRejectedValueOnce(new Error('export generation failed'));
    generateMock.mockResolvedValueOnce(readyExport);

    renderAiExportPage();
    fireEvent.click(screen.getByRole('button', { name: 'Generate export' }));

    expect(await screen.findByText('Export failed')).toBeInTheDocument();
    expect(screen.getByText('export generation failed')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'Retry export' }));

    expect(await screen.findByRole('heading', { name: 'Export ready' })).toBeInTheDocument();
    expect(generateMock).toHaveBeenCalledTimes(2);
  });

  it('uses a local download endpoint and never exposes server file paths', async () => {
    generateMock.mockResolvedValueOnce({
      export: {
        ...readyExport.export,
        id: 'export-2',
        downloadUrl: 'http://localhost:8090/api/ai-export/download?exportId=export-2',
      },
    });

    renderAiExportPage();
    fireEvent.click(screen.getByRole('button', { name: 'Generate export' }));

    const downloadLink = await screen.findByRole('link', { name: 'Download ZIP' });
    expect(downloadLink).toHaveAttribute(
      'href',
      'http://localhost:8090/api/ai-export/download?exportId=export-2',
    );
    expect(downloadLink.getAttribute('href')).not.toContain('/srv/atlas');
    expect(downloadLink.getAttribute('href')).not.toContain('exportFilePath');
  });
});
