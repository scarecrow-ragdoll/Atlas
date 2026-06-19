// FILE: apps/web-admin/src/shared/ui/layout/admin-page-header.tsx
// VERSION: 1.1.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide the standard title, description, and action area for web-admin route content.
//   SCOPE: Owns page heading structure and optional page actions; excludes global shell navigation, theme controls, and data fetching.
//   DEPENDS: react, apps/web-admin/src/shared/ui/lib/utils.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminPageHeader - Accessible page header composition for route-specific headings and actions.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.1.0 - Moved global theme control ownership to the admin app shell.
// END_CHANGE_SUMMARY

import type { ReactNode } from 'react';
import { cn } from '../lib/utils';

type AdminPageHeaderProps = {
  title: string;
  description?: string;
  actions?: ReactNode;
  className?: string;
};

// START_CONTRACT: AdminPageHeader
//   PURPOSE: Render the standard admin page heading, description, and route-specific action slot.
//   INPUTS: { props: AdminPageHeaderProps - title, optional description, optional actions, and optional className }
//   OUTPUTS: { JSX.Element - accessible page header without global shell controls }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: AdminPageHeader
export function AdminPageHeader({ title, description, actions, className }: AdminPageHeaderProps) {
  return (
    <header
      className={cn('flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between', className)}
    >
      <div className="min-w-0 space-y-1">
        <h1 className="text-2xl font-semibold tracking-normal text-foreground">{title}</h1>
        {description ? (
          <p className="max-w-3xl text-sm text-muted-foreground">{description}</p>
        ) : null}
      </div>
      {actions ? <div className="flex shrink-0 flex-wrap gap-2">{actions}</div> : null}
    </header>
  );
}
