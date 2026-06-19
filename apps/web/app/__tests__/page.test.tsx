// FILE: apps/web/app/__tests__/page.test.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify the public Next root layout and server page.
//   SCOPE: Covers metadata, layout wrapping, and server-fetched users passed to the client component; excludes client create behavior.
//   DEPENDS: apps/web/app/layout.tsx, apps/web/app/page.tsx, @testing-library/react, vitest.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   public page tests - Prove Next layout metadata and server page rendering.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Added coverage for server-side REST fetch failures rendering the users fallback.
// END_CHANGE_SUMMARY

import { cleanup, render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { Providers } from '../../src/app/providers';
import RootLayout, { metadata } from '../layout';
import Page, { dynamic } from '../page';

afterEach(() => {
  cleanup();
  vi.restoreAllMocks();
});

describe('public web page', () => {
  it('exports metadata and wraps children in the root layout', () => {
    expect(metadata.title).toBe('Monorepo Template');
    const layout = RootLayout({ children: <span>layout child</span> });
    expect(layout.type).toBe('html');
    expect(layout.props.lang).toBe('en');
    expect(layout.props.children.type).toBe('body');
  });

  it('uses dynamic rendering for runtime REST data', () => {
    expect(dynamic).toBe('force-dynamic');
  });

  it('renders server-fetched users into the public page', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(
        JSON.stringify({
          data: [
            {
              createdAt: '2026-05-24T00:00:00Z',
              email: 'one@example.com',
              id: 'u1',
              name: 'One',
              updatedAt: '2026-05-24T00:00:00Z',
            },
          ],
          meta: { totalCount: 1 },
        }),
        { status: 200 },
      ),
    );

    const ui = await Page();
    render(<Providers>{ui}</Providers>);

    expect(screen.getByText('REST Web')).toBeInTheDocument();
    expect(screen.getByText('One')).toBeInTheDocument();
  });

  it('renders a load failure fallback when the server REST fetch fails', async () => {
    vi.spyOn(globalThis, 'fetch').mockRejectedValueOnce(new Error('offline'));

    const ui = await Page();
    render(<Providers>{ui}</Providers>);

    expect(screen.getByText('REST Web')).toBeInTheDocument();
    expect(screen.getByText('Failed to load users.')).toBeInTheDocument();
    expect(screen.queryByText('No users yet.')).not.toBeInTheDocument();
  });
});
