// FILE: apps/web-admin/src/shared/ui/layout/admin-page-shell.tsx
// VERSION: 1.1.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide the standard content container for web-admin routes rendered inside the app shell.
//   SCOPE: Owns responsive content spacing and max width; excludes main landmark ownership, route-specific headers, and data behavior.
//   DEPENDS: react, apps/web-admin/src/shared/ui/lib/utils.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminPageShell - Standard web-admin content container inside SidebarInset.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.1.0 - Refactored from main landmark to content container for sidebar app shell nesting.
// END_CHANGE_SUMMARY

import type { ComponentPropsWithoutRef } from 'react';
import { cn } from '../lib/utils';

// START_CONTRACT: AdminPageShell
//   PURPOSE: Render the standard responsive content container for admin pages inside the global app shell.
//   INPUTS: { props: ComponentPropsWithoutRef<'div'> - native div props, optional className, and children }
//   OUTPUTS: { JSX.Element - content container without a main landmark }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: AdminPageShell
export function AdminPageShell({ className, children, ...props }: ComponentPropsWithoutRef<'div'>) {
  return (
    <div
      data-testid="admin-page-shell"
      className={cn(
        'mx-auto flex w-full max-w-6xl flex-col gap-6 px-4 py-6 sm:px-6 lg:px-8',
        className,
      )}
      {...props}
    >
      {children}
    </div>
  );
}
