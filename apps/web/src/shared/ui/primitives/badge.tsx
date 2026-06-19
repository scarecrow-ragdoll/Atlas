// FILE: apps/web/src/shared/ui/primitives/badge.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide the shadcn-compatible badge primitive for the public web app.
//   SCOPE: Owns compact status badge rendering; excludes page-specific status semantics.
//   DEPENDS: react, apps/web/src/shared/ui/lib/utils.ts.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Badge - Public web badge primitive.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added independent public web badge primitive.
// END_CHANGE_SUMMARY

import * as React from 'react';

import { cn } from '../lib/utils';

function Badge({ className, ...props }: React.ComponentProps<'span'>) {
  return <span data-slot="badge" className={cn('web-badge', className)} {...props} />;
}

export { Badge };
