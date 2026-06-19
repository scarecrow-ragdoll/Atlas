// FILE: apps/web-admin/src/app/providers.test.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify the web-admin provider boundary.
//   SCOPE: Covers provider rendering around children; excludes React Query internals.
//   DEPENDS: apps/web-admin/src/app/providers.tsx, @testing-library/react, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Providers tests - Prove app context provider renders child routes.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added coverage for the Vite admin provider boundary.
// END_CHANGE_SUMMARY

import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { Providers } from './providers';

describe('Providers', () => {
  it('renders children inside app providers', () => {
    render(
      <Providers>
        <span>admin child</span>
      </Providers>,
    );

    expect(screen.getByText('admin child')).toBeInTheDocument();
  });
});
