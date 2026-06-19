// FILE: apps/web-admin/src/shared/ui/primitives/separator.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide the shadcn separator primitive for web-admin UI compositions.
//   SCOPE: Owns separator primitive rendering; excludes page-specific grouping semantics.
//   DEPENDS: react, radix-ui Separator, apps/web-admin/src/shared/ui/lib/utils.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Separator - shadcn separator primitive.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added shadcn separator primitive under the web-admin UI kit.
// END_CHANGE_SUMMARY

import * as React from 'react';
import { Separator as SeparatorPrimitive } from 'radix-ui';

import { cn } from '@/shared/ui/lib/utils';

function Separator({
  className,
  orientation = 'horizontal',
  decorative = true,
  ...props
}: React.ComponentProps<typeof SeparatorPrimitive.Root>) {
  return (
    <SeparatorPrimitive.Root
      data-slot="separator"
      decorative={decorative}
      orientation={orientation}
      className={cn(
        'shrink-0 bg-border data-[orientation=horizontal]:h-px data-[orientation=horizontal]:w-full data-[orientation=vertical]:h-full data-[orientation=vertical]:w-px',
        className,
      )}
      {...props}
    />
  );
}

export { Separator };
