// FILE: apps/web-admin/src/app/admin-navigation.ts
// VERSION: 1.2.0
// START_MODULE_CONTRACT
//   PURPOSE: Own static web-admin shell navigation metadata and active route resolution.
//   SCOPE: Defines template-native sidebar groups, disabled placeholders, breadcrumbs, user/team placeholders, and route matching; excludes shared UI rendering.
//   DEPENDS: lucide-react, apps/web-admin/src/shared/ui/layout/admin-shell-types.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   adminNavigationGroups - Platform navigation groups for the sidebar shell.
//   adminReferenceItems - Disabled reference placeholders for future admin sections.
//   adminShellUser - Template-native user menu placeholder.
//   adminShellTeams - Template-native team/project switcher placeholder.
//   resolveAdminShellState - Derive active navigation and breadcrumbs from a pathname.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.2.0 - Pointed Atlas nutrition navigation to the factual daily food log route.
// END_CHANGE_SUMMARY

import {
  BookOpenIcon,
  DatabaseIcon,
  HomeIcon,
  UtensilsIcon,
  SettingsIcon,
  UsersIcon,
  WrenchIcon,
} from 'lucide-react';
import type {
  AdminBreadcrumbItem,
  AdminNavigationGroup,
  AdminNavigationItem,
  AdminProjectItem,
  AdminTeamItem,
  AdminUser,
} from '@shared/ui';

export const adminReferenceItems: AdminProjectItem[] = [
  {
    id: 'graphql-admin',
    name: 'GraphQL/Admin',
    href: '#graphql-admin',
    icon: DatabaseIcon,
    disabled: true,
  },
  {
    id: 'system-settings',
    name: 'System/Settings',
    href: '#system-settings',
    icon: SettingsIcon,
    disabled: true,
  },
];

export const adminNavigationGroups: AdminNavigationGroup[] = [
  {
    id: 'platform',
    label: 'Platform',
    items: [
      {
        id: 'overview',
        label: 'Overview',
        href: '/',
        icon: HomeIcon,
        match: (pathname) => pathname === '/',
      },
      {
        id: 'users',
        label: 'Users',
        href: '/users',
        icon: UsersIcon,
        match: (pathname) => pathname === '/users' || pathname.startsWith('/users/'),
        children: [
          { id: 'users-list', label: 'Directory', href: '/users' },
          { id: 'users-detail', label: 'User detail', href: '/users', disabled: true },
        ],
      },
      {
        id: 'ui-kit',
        label: 'UI Kit',
        href: '/ui-kit',
        icon: WrenchIcon,
      },
      {
        id: 'atlas-nutrition',
        label: 'Nutrition',
        href: '/atlas/nutrition',
        icon: UtensilsIcon,
        match: (pathname) =>
          pathname === '/atlas/nutrition' || pathname.startsWith('/atlas/nutrition/'),
        children: [
          {
            id: 'atlas-nutrition-daily-log',
            label: 'Daily Log',
            href: '/atlas/nutrition',
          },
          {
            id: 'atlas-nutrition-products',
            label: 'Product Library',
            href: '/atlas/nutrition/products',
          },
        ],
      },
    ],
  },
];

export const adminShellUser: AdminUser = {
  name: 'Developer',
  email: 'developer@example.local',
  initials: 'DV',
};

export const adminShellTeams: AdminTeamItem[] = [
  {
    id: 'monorepo-template',
    name: 'Monorepo Template',
    plan: 'Admin shell',
    icon: BookOpenIcon,
  },
];

function normalizePathname(pathname: string): string {
  const cleanPath = pathname.split('?')[0]?.split('#')[0] || '/';
  if (cleanPath.length > 1 && cleanPath.endsWith('/')) {
    return cleanPath.slice(0, -1);
  }
  return cleanPath;
}

function cloneItemWithActive(item: AdminNavigationItem, pathname: string): AdminNavigationItem {
  const isActive = item.match ? item.match(pathname) : item.href === pathname;
  return {
    ...item,
    isActive,
    children: item.children?.map((child) => ({
      ...child,
      isActive: !child.disabled && child.href === pathname,
    })),
  };
}

function buildBreadcrumbs(pathname: string): AdminBreadcrumbItem[] {
  if (pathname === '/users') {
    return [{ label: 'Users' }];
  }
  if (pathname.startsWith('/users/')) {
    return [{ label: 'Users', href: '/users' }, { label: 'User detail' }];
  }
  if (pathname === '/ui-kit') {
    return [{ label: 'UI Kit' }];
  }
  if (pathname === '/atlas/nutrition') {
    return [{ label: 'Nutrition' }];
  }
  if (pathname === '/atlas/nutrition/products') {
    return [{ label: 'Nutrition' }, { label: 'Product Library' }];
  }
  if (pathname.startsWith('/atlas/nutrition/')) {
    return [{ label: 'Nutrition' }];
  }
  return [{ label: 'Overview' }];
}

// START_CONTRACT: resolveAdminShellState
//   PURPOSE: Derive shell navigation display state from a browser pathname without importing app data into shared UI.
//   INPUTS: { pathname: string - browser location pathname, optionally with search or hash }
//   OUTPUTS: { navigation, referenceItems, breadcrumbs, user, teams - shell props for AdminAppShell }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: resolveAdminShellState
export function resolveAdminShellState(pathname: string) {
  const normalizedPathname = normalizePathname(pathname);
  return {
    navigation: adminNavigationGroups.map((group) => ({
      ...group,
      items: group.items.map((item) => cloneItemWithActive(item, normalizedPathname)),
    })),
    referenceItems: adminReferenceItems,
    breadcrumbs: buildBreadcrumbs(normalizedPathname),
    user: adminShellUser,
    teams: adminShellTeams,
  };
}
