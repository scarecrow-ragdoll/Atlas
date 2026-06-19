// FILE: apps/web/src/shared/ui/lib/utils.test.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify public web UI utility helpers.
//   SCOPE: Covers class composition behavior used by apps/web UI primitives; excludes component rendering.
//   DEPENDS: apps/web/src/shared/ui/lib/utils.ts, vitest.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   cn tests - Prove conditional classes and conflict resolution.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added public web UI utility coverage.
// END_CHANGE_SUMMARY

import { describe, expect, it } from 'vitest';
import { cn } from './utils';

describe('cn', () => {
  it('merges conditional class values and resolves utility conflicts', () => {
    expect(cn('px-2 text-sm', false && 'hidden', ['font-medium'], 'px-4')).toBe(
      'text-sm font-medium px-4',
    );
  });
});
