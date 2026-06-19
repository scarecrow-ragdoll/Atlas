// FILE: apps/web-admin/src/shared/ui/layout/nav-projects.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Render the Reference navigation group adapted from sidebar-07 projects.
//   SCOPE: Owns reference link and disabled placeholder display; excludes real project actions.
//   DEPENDS: react-router, apps/web-admin/src/shared/ui/primitives/sidebar.tsx, apps/web-admin/src/shared/ui/layout/admin-shell-types.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   NavProjects - Reference navigation group adapted from sidebar-07 projects.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added disabled reference placeholders for future admin surfaces.
// END_CHANGE_SUMMARY

import { Link } from 'react-router';
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuBadge,
  SidebarMenuButton,
  SidebarMenuItem,
} from '../primitives/sidebar';
import type { AdminProjectItem } from './admin-shell-types';

type NavProjectsProps = {
  projects: AdminProjectItem[];
};

// START_CONTRACT: NavProjects
//   PURPOSE: Render reference navigation placeholders without owning app route metadata.
//   INPUTS: { projects: AdminProjectItem[] - app-owned reference navigation items }
//   OUTPUTS: { JSX.Element - sidebar reference group }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: NavProjects
export function NavProjects({ projects }: NavProjectsProps) {
  return (
    <SidebarGroup className="group-data-[collapsible=icon]:hidden">
      <SidebarGroupLabel>Reference</SidebarGroupLabel>
      <SidebarMenu>
        {projects.map((item) => {
          const Icon = item.icon;
          return (
            <SidebarMenuItem key={item.id}>
              <SidebarMenuButton aria-disabled={item.disabled} asChild={!item.disabled}>
                {item.disabled ? (
                  <span>
                    <Icon aria-hidden="true" />
                    <span>{item.name}</span>
                  </span>
                ) : (
                  <Link to={item.href}>
                    <Icon aria-hidden="true" />
                    <span>{item.name}</span>
                  </Link>
                )}
              </SidebarMenuButton>
              {item.disabled ? <SidebarMenuBadge>Coming soon</SidebarMenuBadge> : null}
            </SidebarMenuItem>
          );
        })}
      </SidebarMenu>
    </SidebarGroup>
  );
}
