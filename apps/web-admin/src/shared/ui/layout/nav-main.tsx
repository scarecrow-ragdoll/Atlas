// FILE: apps/web-admin/src/shared/ui/layout/nav-main.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Render the primary Platform navigation group for the web-admin sidebar.
//   SCOPE: Owns main nav item, child item, active, disabled, and tooltip display; excludes route metadata ownership.
//   DEPENDS: react-router, apps/web-admin/src/shared/ui/primitives/collapsible.tsx, apps/web-admin/src/shared/ui/primitives/sidebar.tsx, apps/web-admin/src/shared/ui/layout/admin-shell-types.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   NavMain - Platform navigation group adapted from sidebar-07.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added template-native main sidebar navigation.
// END_CHANGE_SUMMARY

import { ChevronRightIcon } from 'lucide-react';
import { Link } from 'react-router';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '../primitives/collapsible';
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
  useSidebar,
} from '../primitives/sidebar';
import type { AdminNavigationGroup } from './admin-shell-types';

type NavMainProps = {
  groups: AdminNavigationGroup[];
};

// START_CONTRACT: NavMain
//   PURPOSE: Render app-owned navigation groups through reusable sidebar primitives.
//   INPUTS: { groups: AdminNavigationGroup[] - app-owned navigation groups with active state }
//   OUTPUTS: { JSX.Element - sidebar main navigation groups }
//   SIDE_EFFECTS: Closes the mobile sidebar sheet after successful navigation clicks.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: NavMain
export function NavMain({ groups }: NavMainProps) {
  const { isMobile, setOpenMobile } = useSidebar();
  const closeMobileSidebar = () => {
    if (isMobile) {
      setOpenMobile(false);
    }
  };

  return (
    <>
      {groups.map((group) => (
        <SidebarGroup key={group.id}>
          <SidebarGroupLabel>{group.label}</SidebarGroupLabel>
          <SidebarMenu>
            {group.items.map((item) => {
              const Icon = item.icon;
              const hasChildren = Boolean(item.children?.length);
              return (
                <Collapsible
                  asChild
                  className="group/collapsible"
                  defaultOpen={item.isActive && hasChildren}
                  key={item.id}
                >
                  <SidebarMenuItem>
                    <CollapsibleTrigger asChild disabled={item.disabled}>
                      <SidebarMenuButton
                        aria-disabled={item.disabled}
                        asChild={!item.disabled}
                        isActive={item.isActive}
                        tooltip={item.label}
                      >
                        {item.disabled ? (
                          <span>
                            <Icon aria-hidden="true" />
                            <span>{item.label}</span>
                          </span>
                        ) : (
                          <Link onClick={closeMobileSidebar} to={item.href}>
                            <Icon aria-hidden="true" />
                            <span>{item.label}</span>
                            {hasChildren ? (
                              <ChevronRightIcon
                                aria-hidden="true"
                                className="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90"
                              />
                            ) : null}
                          </Link>
                        )}
                      </SidebarMenuButton>
                    </CollapsibleTrigger>
                    {hasChildren ? (
                      <CollapsibleContent>
                        <SidebarMenuSub>
                          {item.children?.map((child) => (
                            <SidebarMenuSubItem key={child.id}>
                              <SidebarMenuSubButton
                                aria-disabled={child.disabled}
                                asChild={!child.disabled}
                                isActive={child.isActive}
                              >
                                {child.disabled ? (
                                  <span>{child.label}</span>
                                ) : (
                                  <Link onClick={closeMobileSidebar} to={child.href}>
                                    {child.label}
                                  </Link>
                                )}
                              </SidebarMenuSubButton>
                            </SidebarMenuSubItem>
                          ))}
                        </SidebarMenuSub>
                      </CollapsibleContent>
                    ) : null}
                  </SidebarMenuItem>
                </Collapsible>
              );
            })}
          </SidebarMenu>
        </SidebarGroup>
      ))}
    </>
  );
}
