// FILE: apps/web-admin/src/shared/ui/layout/admin-shell-header.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Render the global header inside the web-admin sidebar shell.
//   SCOPE: Owns sidebar trigger, breadcrumb display, separator, and global theme placement; excludes page-specific actions and data fetching.
//   DEPENDS: react-router, apps/web-admin/src/shared/ui/primitives/breadcrumb.tsx, apps/web-admin/src/shared/ui/primitives/separator.tsx, apps/web-admin/src/shared/ui/primitives/sidebar.tsx, apps/web-admin/src/shared/ui/layout/theme-toggle.tsx, apps/web-admin/src/shared/ui/layout/admin-shell-types.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminShellHeader - Global header for sidebar trigger, breadcrumbs, and theme action.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added adapted sidebar-07 shell header.
// END_CHANGE_SUMMARY

import { Fragment } from 'react';
import { Link } from 'react-router';
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '../primitives/breadcrumb';
import { Separator } from '../primitives/separator';
import { SidebarTrigger } from '../primitives/sidebar';
import type { AdminBreadcrumbItem } from './admin-shell-types';
import { ThemeToggle } from './theme-toggle';

type AdminShellHeaderProps = {
  breadcrumbs: AdminBreadcrumbItem[];
};

// START_CONTRACT: AdminShellHeader
//   PURPOSE: Render global admin shell navigation controls and breadcrumb context.
//   INPUTS: { breadcrumbs: AdminBreadcrumbItem[] - current app-layer breadcrumb trail }
//   OUTPUTS: { JSX.Element - shell header }
//   SIDE_EFFECTS: ThemeToggle may update persisted theme when clicked.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: AdminShellHeader
export function AdminShellHeader({ breadcrumbs }: AdminShellHeaderProps) {
  return (
    <header className="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12">
      <div className="flex min-w-0 flex-1 items-center gap-2 px-4">
        <SidebarTrigger aria-label="Toggle sidebar" className="-ml-1" />
        <Separator orientation="vertical" className="mr-2 data-[orientation=vertical]:h-4" />
        <Breadcrumb>
          <BreadcrumbList>
            {breadcrumbs.map((breadcrumb, index) => {
              const isLast = index === breadcrumbs.length - 1;
              return (
                <Fragment key={`${breadcrumb.label}-${index}`}>
                  <BreadcrumbItem
                    className={index === 0 && !isLast ? 'hidden md:block' : undefined}
                  >
                    {breadcrumb.href && !isLast ? (
                      <BreadcrumbLink asChild>
                        <Link aria-label={`${breadcrumb.label} breadcrumb`} to={breadcrumb.href}>
                          {breadcrumb.label}
                        </Link>
                      </BreadcrumbLink>
                    ) : (
                      <BreadcrumbPage>{breadcrumb.label}</BreadcrumbPage>
                    )}
                  </BreadcrumbItem>
                  {!isLast ? <BreadcrumbSeparator className="hidden md:block" /> : null}
                </Fragment>
              );
            })}
          </BreadcrumbList>
        </Breadcrumb>
      </div>
      <div className="flex items-center gap-2 px-4">
        <ThemeToggle />
      </div>
    </header>
  );
}
