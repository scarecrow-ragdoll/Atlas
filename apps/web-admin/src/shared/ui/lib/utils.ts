// FILE: apps/web-admin/src/shared/ui/lib/utils.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide shared utility helpers for web-admin UI components.
//   SCOPE: Owns class name composition used by shadcn primitives and admin compositions; excludes component rendering.
//   DEPENDS: clsx, tailwind-merge.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   cn - Merge conditional class values and resolve Tailwind utility conflicts.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added shadcn-compatible class merge helper.
// END_CHANGE_SUMMARY

import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

// START_CONTRACT: cn
//   PURPOSE: Compose conditional class values and resolve Tailwind conflicts for UI-kit internals.
//   INPUTS: { inputs: ClassValue[] - conditional class values from shadcn primitives and layout compositions }
//   OUTPUTS: { string - merged class name string }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: cn
export function cn(...inputs: ClassValue[]): string {
  return twMerge(clsx(inputs));
}
