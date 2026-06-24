// FILE: apps/web-admin/src/app/i18n.test.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify Atlas web-admin i18n state remains safe when browser storage is unavailable.
//   SCOPE: Covers default language fallback and render safety for storage read/write failures; excludes page-specific translations.
//   DEPENDS: apps/web-admin/src/app/i18n.tsx, @testing-library/react, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   i18n storage tests - Prove I18nProvider defaults to English and does not crash when localStorage throws.
// END_MODULE_MAP

import { render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { I18nProvider, useI18n } from './i18n';

function Probe() {
  const { language, t } = useI18n();
  return (
    <div>
      <span>{language}</span>
      <span>{t('nutrition.productLibrary')}</span>
    </div>
  );
}

describe('I18nProvider', () => {
  const originalLocalStorage = window.localStorage;

  afterEach(() => {
    Object.defineProperty(window, 'localStorage', {
      configurable: true,
      value: originalLocalStorage,
    });
    vi.restoreAllMocks();
  });

  it('defaults to English when localStorage access throws', () => {
    Object.defineProperty(window, 'localStorage', {
      configurable: true,
      get() {
        throw new Error('storage unavailable');
      },
    });

    render(
      <I18nProvider>
        <Probe />
      </I18nProvider>,
    );

    expect(screen.getByText('en')).toBeInTheDocument();
    expect(screen.getByText('Product Library')).toBeInTheDocument();
  });
});
