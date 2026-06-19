// FILE: apps/web/src/shared/ui/primitives/label.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide the shadcn-compatible label primitive for the public web app.
//   SCOPE: Owns accessible label rendering; excludes form validation behavior.
//   DEPENDS: react, radix-ui Label, apps/web/src/shared/ui/lib/utils.ts.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Label - Public web label primitive.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added independent public web label primitive.
// END_CHANGE_SUMMARY

import * as React from 'react';
import { Label as LabelPrimitive } from 'radix-ui';

import { cn } from '../lib/utils';

function Label({ className, ...props }: React.ComponentProps<typeof LabelPrimitive.Root>) {
  return (
    <LabelPrimitive.Root data-slot="label" className={cn('web-label', className)} {...props} />
  );
}

export { Label };
