// FILE: apps/web-admin/src/shared/ui/layout/team-switcher.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Render the template-native shell project switcher adapted from sidebar-07.
//   SCOPE: Owns placeholder team display and local active selection; excludes real team persistence or API behavior.
//   DEPENDS: react, lucide-react, apps/web-admin/src/shared/ui/primitives/dropdown-menu.tsx, apps/web-admin/src/shared/ui/primitives/sidebar.tsx, apps/web-admin/src/shared/ui/layout/admin-shell-types.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   TeamSwitcher - Template-native project switcher placeholder.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added template-native sidebar project switcher.
// END_CHANGE_SUMMARY

import * as React from 'react';
import { ChevronsUpDownIcon } from 'lucide-react';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from '../primitives/dropdown-menu';
import { SidebarMenu, SidebarMenuButton, SidebarMenuItem, useSidebar } from '../primitives/sidebar';
import type { AdminTeamItem } from './admin-shell-types';

type TeamSwitcherProps = {
  teams: AdminTeamItem[];
};

// START_CONTRACT: TeamSwitcher
//   PURPOSE: Render the active template workspace placeholder and available placeholder teams.
//   INPUTS: { teams: AdminTeamItem[] - app-provided static team placeholders }
//   OUTPUTS: { JSX.Element | null - sidebar team switcher menu }
//   SIDE_EFFECTS: Updates local active-team state after menu selection.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: TeamSwitcher
export function TeamSwitcher({ teams }: TeamSwitcherProps) {
  const { isMobile } = useSidebar();
  const [activeTeam, setActiveTeam] = React.useState(teams[0]);

  if (!activeTeam) {
    return null;
  }

  const ActiveIcon = activeTeam.icon;

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
              size="lg"
            >
              <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                <ActiveIcon aria-hidden="true" className="size-4" />
              </div>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-medium">{activeTeam.name}</span>
                <span className="truncate text-xs">{activeTeam.plan}</span>
              </div>
              <ChevronsUpDownIcon aria-hidden="true" className="ml-auto size-4" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            align="start"
            className="w-fit"
            side={isMobile ? 'bottom' : 'right'}
            sideOffset={4}
          >
            <DropdownMenuLabel className="text-xs text-muted-foreground">
              Workspace
            </DropdownMenuLabel>
            {teams.map((team) => {
              const Icon = team.icon;
              return (
                <DropdownMenuItem
                  className="gap-2 p-2"
                  key={team.id}
                  onClick={() => setActiveTeam(team)}
                >
                  <div className="flex size-6 items-center justify-center rounded-md border">
                    <Icon aria-hidden="true" className="size-3.5" />
                  </div>
                  {team.name}
                </DropdownMenuItem>
              );
            })}
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
