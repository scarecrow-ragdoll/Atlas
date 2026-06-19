// FILE: apps/web-admin/src/shared/ui/primitives/skeleton.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide the shadcn skeleton primitive for web-admin UI loading states.
//   SCOPE: Owns skeleton placeholder rendering; excludes page-specific loading decisions.
//   DEPENDS: react, apps/web-admin/src/shared/ui/lib/utils.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Skeleton - shadcn skeleton primitive.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added shadcn skeleton primitive under the web-admin UI kit.
// END_CHANGE_SUMMARY

import { cn } from '@/shared/ui/lib/utils';

function Skeleton({ className, ...props }: React.ComponentProps<'div'>) {
  return (
    <div
      data-slot="skeleton"
      className={cn('animate-pulse rounded-md bg-accent', className)}
      {...props}
    />
  );
}

export { Skeleton };
