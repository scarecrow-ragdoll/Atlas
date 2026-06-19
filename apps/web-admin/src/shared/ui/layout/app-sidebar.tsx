// FILE: apps/web-admin/src/shared/ui/layout/app-sidebar.tsx
// VERSION: 1.1.0
// START_MODULE_CONTRACT
//   PURPOSE: Render the adapted sidebar-07 sidebar for web-admin.
//   SCOPE: Owns sidebar header, content, footer, rail, user action pass-through, and template-native nav composition; excludes route matching and page data behavior.
//   DEPENDS: apps/web-admin/src/shared/ui/primitives/sidebar.tsx, apps/web-admin/src/shared/ui/layout/nav-main.tsx, apps/web-admin/src/shared/ui/layout/nav-projects.tsx, apps/web-admin/src/shared/ui/layout/nav-user.tsx, apps/web-admin/src/shared/ui/layout/team-switcher.tsx, apps/web-admin/src/shared/ui/layout/admin-shell-types.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AppSidebar - Template-native sidebar-07 composition.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.1.0 - Passed authenticated admin logout action state into the user menu.
// END_CHANGE_SUMMARY

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from '../primitives/sidebar';
import type {
  AdminNavigationGroup,
  AdminProjectItem,
  AdminTeamItem,
  AdminUser,
  AdminUserAction,
} from './admin-shell-types';
import { NavMain } from './nav-main';
import { NavProjects } from './nav-projects';
import { NavUser } from './nav-user';
import { TeamSwitcher } from './team-switcher';

type AppSidebarProps = AdminUserAction & {
  navigation: AdminNavigationGroup[];
  referenceItems: AdminProjectItem[];
  teams: AdminTeamItem[];
  user: AdminUser;
};

// START_CONTRACT: AppSidebar
//   PURPOSE: Compose the sidebar header, navigation, reference placeholders, user menu, and rail.
//   INPUTS: { navigation, referenceItems, teams, user - app-owned shell display data }
//   OUTPUTS: { JSX.Element - sidebar-07 shell sidebar }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: AppSidebar
export function AppSidebar({
  isLogoutPending,
  navigation,
  onLogout,
  referenceItems,
  teams,
  user,
}: AppSidebarProps) {
  return (
    <Sidebar aria-label="Admin navigation" collapsible="icon" role="navigation">
      <SidebarHeader>
        <TeamSwitcher teams={teams} />
      </SidebarHeader>
      <SidebarContent>
        <NavMain groups={navigation} />
        <NavProjects projects={referenceItems} />
      </SidebarContent>
      <SidebarFooter>
        <NavUser isLogoutPending={isLogoutPending} onLogout={onLogout} user={user} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
