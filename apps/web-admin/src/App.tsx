// FILE: apps/web-admin/src/App.tsx
// VERSION: 1.2.0
// START_MODULE_CONTRACT
//   PURPOSE: Own the web-admin Vite route table and auth-guarded app shell layout route.
//   SCOPE: Maps /login publicly, protects every other admin route through CurrentAdmin, and redirects unknown routes; excludes page internals.
//   DEPENDS: react-router, apps/web-admin/src/app/protected-admin-layout.tsx, apps/web-admin/src/entities/admin-auth/provider.tsx, apps/web-admin/src/pages/login-page.tsx, apps/web-admin/src/pages/home.tsx, apps/web-admin/src/pages/users-page.tsx, apps/web-admin/src/pages/user-detail-page.tsx, apps/web-admin/src/pages/ui-kit-page.tsx.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - BrowserRouter-backed route table for public login and protected admin pages.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.2.0 - Added AuthProvider and protected route layout around non-login admin routes.
// END_CHANGE_SUMMARY

import { BrowserRouter, Navigate, Route, Routes } from 'react-router';
import { AuthProvider } from '@entities/admin-auth/provider';
import { ProtectedAdminLayout } from './app/protected-admin-layout';
import HomePage from './pages/home';
import LoginPage from './pages/login-page';
import UiKitPage from './pages/ui-kit-page';
import UserDetailPage from './pages/user-detail-page';
import UsersPage from './pages/users-page';

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route element={<ProtectedAdminLayout />}>
            <Route path="/" element={<HomePage />} />
            <Route path="/ui-kit" element={<UiKitPage />} />
            <Route path="/users" element={<UsersPage />} />
            <Route path="/users/:id" element={<UserDetailPage />} />
          </Route>
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}
