// FILE: apps/web-admin/src/app/protected-admin-layout.tsx
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Protect web-admin routes behind the backend admin session.
//   SCOPE: Gates non-login routes on CurrentAdmin state, preserves safe return-to paths, and renders AdminLayout only after auth is known; excludes login page form behavior.
//   DEPENDS: react-router, apps/web-admin/src/app/admin-layout.tsx, apps/web-admin/src/entities/admin-auth/provider.tsx, apps/web-admin/src/shared/ui.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   ProtectedAdminLayout - Auth gate for all non-login web-admin routes.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Redirected completed logout transitions to plain /login.
// END_CHANGE_SUMMARY

import { Navigate, useLocation } from 'react-router';
import { useAdminAuth } from '@entities/admin-auth/provider';
import { AdminPageShell, Skeleton } from '@shared/ui';
import { AdminLayout } from './admin-layout';

function buildReturnTo(location: ReturnType<typeof useLocation>) {
  return `${location.pathname}${location.search}${location.hash}`;
}

// START_CONTRACT: ProtectedAdminLayout
//   PURPOSE: Render protected admin routes only after CurrentAdmin confirms an authenticated admin.
//   INPUTS: none.
//   OUTPUTS: { JSX.Element - loading state, login redirect, or AdminLayout }
//   SIDE_EFFECTS: Navigates unauthenticated users to /login with a safe return path.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: ProtectedAdminLayout
export function ProtectedAdminLayout() {
  const location = useLocation();
  const { admin, hasCompletedLogout, isLoading } = useAdminAuth();

  if (isLoading) {
    return (
      <AdminPageShell>
        <div
          aria-label="Loading admin session"
          aria-live="polite"
          className="space-y-4"
          role="status"
        >
          <Skeleton className="h-9 w-64" />
          <Skeleton className="h-48 w-full" />
        </div>
      </AdminPageShell>
    );
  }

  if (!admin) {
    const target = hasCompletedLogout
      ? '/login'
      : `/login?from=${encodeURIComponent(buildReturnTo(location))}`;
    return <Navigate replace to={target} />;
  }

  return <AdminLayout admin={admin} />;
}
