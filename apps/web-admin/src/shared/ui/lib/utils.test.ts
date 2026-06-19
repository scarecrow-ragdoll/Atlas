// FILE: apps/web-admin/src/shared/ui/lib/utils.test.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify shared web-admin UI utility behavior.
//   SCOPE: Covers class composition and Tailwind conflict merging; excludes component rendering.
//   DEPENDS: apps/web-admin/src/shared/ui/lib/utils.ts, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   cn tests - Prove class values are merged and Tailwind conflicts resolve predictably.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added red coverage for the UI class merge helper.
// END_CHANGE_SUMMARY

import { describe, expect, it } from 'vitest';
import { cn } from './utils';

describe('cn', () => {
  it('merges conditional class values and resolves Tailwind conflicts', () => {
    expect(cn('px-2 text-sm', false && 'hidden', ['font-medium'], 'px-4')).toBe(
      'text-sm font-medium px-4',
    );
  });
});
