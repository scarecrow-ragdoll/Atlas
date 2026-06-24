// FILE: apps/web-admin/vitest.setup.ts
// VERSION: 1.1.0
// START_MODULE_CONTRACT
//   PURPOSE: Configure shared web-admin Vitest browser test globals.
//   SCOPE: Installs jest-dom matchers and browser API fallbacks missing from the Bun/jsdom test runtime; excludes per-test mocks.
//   DEPENDS: @testing-library/jest-dom, vitest jsdom environment.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   localStorage fallback - Provides an in-memory Storage implementation when Bun/jsdom does not expose one.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.1.0 - Added a localStorage fallback for app-shell tests under Bun/jsdom.
// END_CHANGE_SUMMARY

import '@testing-library/jest-dom';

class InMemoryStorage implements Storage {
  private readonly values = new Map<string, string>();

  get length() {
    return this.values.size;
  }

  clear() {
    this.values.clear();
  }

  getItem(key: string) {
    return this.values.get(key) ?? null;
  }

  key(index: number) {
    return Array.from(this.values.keys())[index] ?? null;
  }

  removeItem(key: string) {
    this.values.delete(key);
  }

  setItem(key: string, value: string) {
    this.values.set(key, value);
  }
}

function canUseLocalStorage() {
  try {
    return typeof window !== 'undefined' && Boolean(window.localStorage);
  } catch {
    return false;
  }
}

if (typeof window !== 'undefined' && !canUseLocalStorage()) {
  const localStorageFallback = new InMemoryStorage();

  Object.defineProperty(window, 'localStorage', {
    configurable: true,
    value: localStorageFallback,
  });

  Object.defineProperty(globalThis, 'localStorage', {
    configurable: true,
    value: localStorageFallback,
  });
}
