// FILE: apps/web-admin/src/app/admin-layout.tsx
// VERSION: 1.1.0
// START_MODULE_CONTRACT
//   PURPOSE: Bridge app-owned route metadata into the shared web-admin shell.
//   SCOPE: Reads React Router location, resolves shell state, maps current admin into the shell user, and renders Outlet content; excludes page data behavior and shared UI internals.
//   DEPENDS: react-router, apps/web-admin/src/app/admin-navigation.ts, apps/web-admin/src/entities/admin-auth/model.ts, apps/web-admin/src/shared/ui.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminLayout - React Router layout route for the shared sidebar shell.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.1.0 - Mapped authenticated admin into the shell user contract.
// END_CHANGE_SUMMARY

import { Outlet, useLocation, useNavigate } from 'react-router';
import { adminToShellUser, type AdminPrincipal } from '@entities/admin-auth/model';
import { useAdminAuth } from '@entities/admin-auth/provider';
import { AdminAppShell } from '@shared/ui';
import { resolveAdminShellState } from './admin-navigation';

type AdminLayoutProps = {
  admin: AdminPrincipal;
};

// START_CONTRACT: AdminLayout
//   PURPOSE: Render current child route content inside the global web-admin shell.
//   INPUTS: { admin: AdminPrincipal - authenticated backend admin }
//   OUTPUTS: { JSX.Element - AdminAppShell with current route Outlet }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: AdminLayout
export function AdminLayout({ admin }: AdminLayoutProps) {
  const location = useLocation();
  const navigate = useNavigate();
  const { isLogoutPending, logout } = useAdminAuth();
  const shellState = resolveAdminShellState(location.pathname);

  async function handleLogout() {
    await logout();
    navigate('/login', { replace: true });
  }

  return (
    <AdminAppShell
      pathname={location.pathname}
      {...shellState}
      isLogoutPending={isLogoutPending}
      onLogout={handleLogout}
      user={adminToShellUser(admin)}
    >
      <Outlet />
    </AdminAppShell>
  );
}
