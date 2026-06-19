// FILE: apps/web-admin/src/pages/ui-kit-page.test.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify the web-admin UI kit reference page.
//   SCOPE: Covers static component showcase sections and local-only rendering; excludes visual pixel assertions and API behavior.
//   DEPENDS: apps/web-admin/src/pages/ui-kit-page.tsx, @testing-library/react, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   UiKitPage tests - Prove the reference page demonstrates approved UI-kit areas without API calls.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Added shell foundation showcase coverage.
// END_CHANGE_SUMMARY

import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { afterAll, afterEach, beforeAll, describe, expect, it, vi } from 'vitest';
import UiKitPage from './ui-kit-page';

const originalMatchMedia = window.matchMedia;
const originalInnerWidth = window.innerWidth;

type UiKitMediaListener = () => void;

function installUiKitMatchMedia() {
  const listeners = new Set<UiKitMediaListener>();
  window.matchMedia = vi.fn().mockImplementation((query: string) => ({
    matches: query.includes('767') ? window.innerWidth < 768 : false,
    media: query,
    onchange: null,
    addEventListener: (_event: string, listener: UiKitMediaListener) => listeners.add(listener),
    removeEventListener: (_event: string, listener: UiKitMediaListener) =>
      listeners.delete(listener),
    addListener: vi.fn(),
    removeListener: vi.fn(),
    dispatchEvent: vi.fn(),
  }));
}

beforeAll(() => {
  Object.defineProperty(window, 'innerWidth', {
    configurable: true,
    value: 1024,
  });
  installUiKitMatchMedia();
});

afterEach(() => {
  document.cookie = 'sidebar_state=; path=/; max-age=0';
  installUiKitMatchMedia();
});

afterAll(() => {
  if (originalMatchMedia) {
    window.matchMedia = originalMatchMedia;
  } else {
    delete (window as Partial<Window>).matchMedia;
  }

  Object.defineProperty(window, 'innerWidth', {
    configurable: true,
    value: originalInnerWidth,
  });
});

describe('UiKitPage', () => {
  it('renders the broad UI-kit showcase sections from local data', () => {
    render(
      <MemoryRouter>
        <UiKitPage />
      </MemoryRouter>,
    );

    expect(screen.getByRole('heading', { name: 'UI Kit' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Actions' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Forms' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Feedback' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Data' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Overlays And Navigation' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Shell Foundation' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Admin Compositions' })).toBeInTheDocument();
    expect(screen.getByText('Typography scale')).toBeInTheDocument();
    expect(screen.getByText('Spacing examples')).toBeInTheDocument();
    expect(screen.getByText('Radius examples')).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'SidebarProvider' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Breadcrumb' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Avatar' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Collapsible' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Sheet' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Creating...' })).toHaveAttribute(
      'aria-busy',
      'true',
    );
    expect(screen.getByText('AdminToolbar')).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'No filtered users' })).toBeInTheDocument();
    expect(screen.getByText('ada@example.com')).toBeInTheDocument();
  });
});
