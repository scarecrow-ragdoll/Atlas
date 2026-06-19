// FILE: apps/web-admin/src/shared/ui/layout/admin-toolbar.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide a responsive toolbar for web-admin filters and commands.
//   SCOPE: Owns horizontal command wrapping; excludes command behavior.
//   DEPENDS: react, apps/web-admin/src/shared/ui/lib/utils.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminToolbar - Responsive admin command toolbar.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added standard admin toolbar.
// END_CHANGE_SUMMARY

import type { ComponentPropsWithoutRef } from 'react';
import { cn } from '../lib/utils';

// START_CONTRACT: AdminToolbar
//   PURPOSE: Render the standard responsive toolbar for admin filters and commands.
//   INPUTS: { props: ComponentPropsWithoutRef<'div'> - native div props, optional className, and children }
//   OUTPUTS: { JSX.Element - responsive toolbar container }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: AdminToolbar
export function AdminToolbar({ className, children, ...props }: ComponentPropsWithoutRef<'div'>) {
  return (
    <div
      className={cn(
        'flex flex-col gap-3 rounded-lg border bg-card p-3 sm:flex-row sm:items-center sm:justify-between',
        className,
      )}
      {...props}
    >
      {children}
    </div>
  );
}
