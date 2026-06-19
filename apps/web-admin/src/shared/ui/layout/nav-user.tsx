// FILE: apps/web-admin/src/shared/ui/layout/nav-user.tsx
// VERSION: 1.1.0
// START_MODULE_CONTRACT
//   PURPOSE: Render the authenticated admin user menu adapted from sidebar-07.
//   SCOPE: Owns user display and logout menu item rendering; excludes authentication state ownership and user profile loading.
//   DEPENDS: lucide-react, apps/web-admin/src/shared/ui/primitives/avatar.tsx, apps/web-admin/src/shared/ui/primitives/dropdown-menu.tsx, apps/web-admin/src/shared/ui/primitives/sidebar.tsx, apps/web-admin/src/shared/ui/layout/admin-shell-types.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   NavUser - Template-native authenticated admin user menu.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.1.0 - Replaced placeholder menu actions with real logout action rendering.
// END_CHANGE_SUMMARY

import { ChevronsUpDownIcon, LogOutIcon } from 'lucide-react';
import { Avatar, AvatarFallback, AvatarImage } from '../primitives/avatar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '../primitives/dropdown-menu';
import { SidebarMenu, SidebarMenuButton, SidebarMenuItem, useSidebar } from '../primitives/sidebar';
import type { AdminUser, AdminUserAction } from './admin-shell-types';

type NavUserProps = AdminUserAction & {
  user: AdminUser;
};

// START_CONTRACT: NavUser
//   PURPOSE: Render the authenticated admin user menu and logout command.
//   INPUTS: { user: AdminUser - current admin display data, onLogout: optional logout command, isLogoutPending: pending state }
//   OUTPUTS: { JSX.Element - sidebar user menu }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: NavUser
export function NavUser({ isLogoutPending, onLogout, user }: NavUserProps) {
  const { isMobile } = useSidebar();

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
              size="lg"
            >
              <Avatar className="h-8 w-8 rounded-lg">
                {user.avatarUrl ? <AvatarImage alt={user.name} src={user.avatarUrl} /> : null}
                <AvatarFallback className="rounded-lg">{user.initials}</AvatarFallback>
              </Avatar>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-medium">{user.name}</span>
                <span className="truncate text-xs">{user.email}</span>
              </div>
              <ChevronsUpDownIcon aria-hidden="true" className="ml-auto size-4" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            align="end"
            className="w-fit"
            side={isMobile ? 'bottom' : 'right'}
            sideOffset={4}
          >
            <DropdownMenuLabel className="p-0 font-normal">
              <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                <Avatar className="h-8 w-8 rounded-lg">
                  {user.avatarUrl ? <AvatarImage alt={user.name} src={user.avatarUrl} /> : null}
                  <AvatarFallback className="rounded-lg">{user.initials}</AvatarFallback>
                </Avatar>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-medium">{user.name}</span>
                  <span className="truncate text-xs">{user.email}</span>
                </div>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <DropdownMenuItem
                disabled={!onLogout || isLogoutPending}
                onClick={() => void onLogout?.()}
              >
                <LogOutIcon aria-hidden="true" />
                {isLogoutPending ? 'Logging out...' : 'Logout'}
              </DropdownMenuItem>
            </DropdownMenuGroup>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
