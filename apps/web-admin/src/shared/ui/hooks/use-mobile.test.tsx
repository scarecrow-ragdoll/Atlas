// FILE: apps/web-admin/src/shared/ui/hooks/use-mobile.test.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify the web-admin sidebar mobile breakpoint hook.
//   SCOPE: Covers initial desktop/mobile detection and media-query change cleanup; excludes sidebar rendering.
//   DEPENDS: apps/web-admin/src/shared/ui/hooks/use-mobile.ts, @testing-library/react, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   useIsMobile tests - Prove desktop and mobile breakpoint behavior for sidebar state.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added focused coverage for the sidebar mobile hook.
// END_CHANGE_SUMMARY

import { act, renderHook, waitFor } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { useIsMobile } from './use-mobile';

type MediaListener = () => void;

const originalMatchMedia = window.matchMedia;
const originalInnerWidth = window.innerWidth;

function setViewportWidth(width: number) {
  Object.defineProperty(window, 'innerWidth', {
    configurable: true,
    value: width,
  });
}

function installMatchMedia() {
  const listeners = new Set<MediaListener>();
  window.matchMedia = vi.fn().mockImplementation((query: string) => ({
    matches: query.includes('767') ? window.innerWidth < 768 : false,
    media: query,
    onchange: null,
    addEventListener: (_event: string, listener: MediaListener) => listeners.add(listener),
    removeEventListener: (_event: string, listener: MediaListener) => listeners.delete(listener),
    addListener: vi.fn(),
    removeListener: vi.fn(),
    dispatchEvent: vi.fn(),
  }));
  return listeners;
}

afterEach(() => {
  window.matchMedia = originalMatchMedia;
  setViewportWidth(originalInnerWidth);
});

describe('useIsMobile', () => {
  it('returns false for desktop widths', async () => {
    setViewportWidth(1024);
    installMatchMedia();

    const { result } = renderHook(() => useIsMobile());

    await waitFor(() => expect(result.current).toBe(false));
  });

  it('returns true for mobile widths and responds to media changes', async () => {
    setViewportWidth(500);
    const listeners = installMatchMedia();

    const { result } = renderHook(() => useIsMobile());

    await waitFor(() => expect(result.current).toBe(true));

    act(() => {
      setViewportWidth(900);
      listeners.forEach((listener) => listener());
    });

    await waitFor(() => expect(result.current).toBe(false));
  });
});
