// FILE: apps/web/src/shared/ui/primitives/input.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide the shadcn-compatible input primitive for the public web app.
//   SCOPE: Owns text input rendering and styling hooks; excludes form validation behavior.
//   DEPENDS: react, apps/web/src/shared/ui/lib/utils.ts.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Input - Public web input primitive.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added independent public web input primitive.
// END_CHANGE_SUMMARY

import * as React from 'react';

import { cn } from '../lib/utils';

function Input({ className, type, ...props }: React.ComponentProps<'input'>) {
  return <input data-slot="input" type={type} className={cn('web-input', className)} {...props} />;
}

export { Input };
