// FILE: apps/web-admin/src/shared/ui/layout/admin-shell-types.ts
// VERSION: 1.1.0
// START_MODULE_CONTRACT
//   PURPOSE: Define shared web-admin shell data contracts consumed by UI-kit shell compositions.
//   SCOPE: Owns navigation, breadcrumb, authenticated user, user actions, and team types; excludes app-owned metadata values and route matching implementation.
//   DEPENDS: react, lucide-react.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TYPES
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminIcon - Icon component contract for shell navigation.
//   AdminNavigationChild - Child navigation item contract.
//   AdminNavigationItem - Main navigation item contract.
//   AdminNavigationGroup - Sidebar navigation group contract.
//   AdminProjectItem - Reference navigation item contract.
//   AdminBreadcrumbItem - Header breadcrumb item contract.
//   AdminUser - Authenticated admin shell user contract.
//   AdminUserAction - Authenticated admin user-menu action contract.
//   AdminTeamItem - Team switcher placeholder contract.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.1.0 - Added authenticated user menu logout action state.
// END_CHANGE_SUMMARY

/* v8 ignore start -- Type-only shell contracts are erased at runtime and have no behavior to execute. */
import type { ComponentType, SVGProps } from 'react';

export type AdminIcon = ComponentType<SVGProps<SVGSVGElement>>;

export type AdminNavigationChild = {
  id: string;
  label: string;
  href: string;
  disabled?: boolean;
  isActive?: boolean;
};

export type AdminNavigationItem = {
  id: string;
  label: string;
  href: string;
  icon: AdminIcon;
  disabled?: boolean;
  isActive?: boolean;
  match?: (pathname: string) => boolean;
  children?: AdminNavigationChild[];
};

export type AdminNavigationGroup = {
  id: string;
  label: string;
  items: AdminNavigationItem[];
};

export type AdminProjectItem = {
  id: string;
  name: string;
  href: string;
  icon: AdminIcon;
  disabled?: boolean;
};

export type AdminBreadcrumbItem = {
  label: string;
  href?: string;
};

export type AdminUser = {
  name: string;
  email: string;
  initials: string;
  avatarUrl?: string;
};

export type AdminUserAction = {
  isLogoutPending?: boolean;
  onLogout?: () => void | Promise<void>;
};

export type AdminTeamItem = {
  id: string;
  name: string;
  plan: string;
  icon: AdminIcon;
};
/* v8 ignore stop */
