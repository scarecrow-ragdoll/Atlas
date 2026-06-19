// FILE: apps/web/src/shared/ui/lib/utils.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide class name composition helpers for the public web UI kit.
//   SCOPE: Owns conditional class merging for apps/web shadcn-compatible primitives; excludes component rendering.
//   DEPENDS: clsx, tailwind-merge.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   cn - Merge conditional class values and Tailwind-style conflicts.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added public web UI-kit class merge helper.
// END_CHANGE_SUMMARY

import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

// START_CONTRACT: cn
//   PURPOSE: Compose conditional class values and resolve utility conflicts for public web UI primitives.
//   INPUTS: { inputs: ClassValue[] - conditional class values from UI primitives and route compositions }
//   OUTPUTS: { string - merged class name string }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB / V-M-WEB.
// END_CONTRACT: cn
export function cn(...inputs: ClassValue[]): string {
  return twMerge(clsx(inputs));
}
