// FILE: apps/web-admin/src/shared/ui/layout/admin-app-shell.tsx
// VERSION: 1.1.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide the global sidebar-07 application shell for web-admin routes.
//   SCOPE: Owns SidebarProvider, adapted sidebar, SidebarInset main landmark, shell header, authenticated user action pass-through, and content slot; excludes app-owned route metadata and page data behavior.
//   DEPENDS: react, apps/web-admin/src/shared/ui/primitives/sidebar.tsx, apps/web-admin/src/shared/ui/primitives/tooltip.tsx, apps/web-admin/src/shared/ui/layout/app-sidebar.tsx, apps/web-admin/src/shared/ui/layout/admin-shell-header.tsx, apps/web-admin/src/shared/ui/layout/admin-shell-types.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminAppShell - Global sidebar app shell for all web-admin routes.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.1.0 - Passed authenticated admin logout action state into the sidebar.
// END_CHANGE_SUMMARY

import type { ReactNode } from 'react';
import { SidebarInset, SidebarProvider } from '../primitives/sidebar';
import { TooltipProvider } from '../primitives/tooltip';
import { AdminShellHeader } from './admin-shell-header';
import { AppSidebar } from './app-sidebar';
import type {
  AdminBreadcrumbItem,
  AdminNavigationGroup,
  AdminProjectItem,
  AdminTeamItem,
  AdminUser,
  AdminUserAction,
} from './admin-shell-types';

type AdminAppShellProps = AdminUserAction & {
  breadcrumbs: AdminBreadcrumbItem[];
  children: ReactNode;
  navigation: AdminNavigationGroup[];
  pathname: string;
  referenceItems: AdminProjectItem[];
  teams: AdminTeamItem[];
  user: AdminUser;
};

// START_CONTRACT: AdminAppShell
//   PURPOSE: Render all web-admin route content inside the adapted sidebar-07 shell.
//   INPUTS: { props: AdminAppShellProps - app-owned nav data, breadcrumbs, pathname, placeholders, and children }
//   OUTPUTS: { JSX.Element - sidebar provider with app sidebar, header, and main content inset }
//   SIDE_EFFECTS: SidebarProvider may persist sidebar open state in document.cookie.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: AdminAppShell
export function AdminAppShell({
  breadcrumbs,
  children,
  isLogoutPending,
  navigation,
  onLogout,
  referenceItems,
  teams,
  user,
}: AdminAppShellProps) {
  return (
    <SidebarProvider>
      <TooltipProvider>
        <AppSidebar
          isLogoutPending={isLogoutPending}
          navigation={navigation}
          onLogout={onLogout}
          referenceItems={referenceItems}
          teams={teams}
          user={user}
        />
        <SidebarInset>
          <AdminShellHeader breadcrumbs={breadcrumbs} />
          {children}
        </SidebarInset>
      </TooltipProvider>
    </SidebarProvider>
  );
}
