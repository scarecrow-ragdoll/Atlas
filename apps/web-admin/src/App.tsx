// FILE: apps/web-admin/src/App.tsx
// VERSION: 1.4.0
// START_MODULE_CONTRACT
//   PURPOSE: Own the web-admin Vite route table and auth-guarded app shell layout route.
//   SCOPE: Maps /login publicly, protects every other admin route through CurrentAdmin, and redirects unknown routes; excludes page internals.
//   DEPENDS: react-router, apps/web-admin/src/app/protected-admin-layout.tsx, apps/web-admin/src/entities/admin-auth/provider.tsx, apps/web-admin/src/pages/login-page.tsx, apps/web-admin/src/pages/home.tsx, apps/web-admin/src/pages/users-page.tsx, apps/web-admin/src/pages/user-detail-page.tsx, apps/web-admin/src/pages/ui-kit-page.tsx, apps/web-admin/src/pages/atlas/nutrition-overview-page.tsx, apps/web-admin/src/pages/atlas/product-library-page.tsx.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - BrowserRouter-backed route table for public login and protected admin pages.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.4.0 - Added factual Atlas nutrition daily log route and legacy override redirect.
// END_CHANGE_SUMMARY

import { BrowserRouter, Navigate, Route, Routes } from 'react-router';
import { AuthProvider } from '@entities/admin-auth/provider';
import { I18nProvider } from './app/i18n';
import { ProtectedAdminLayout } from './app/protected-admin-layout';
import DailyNutritionOverridePage from './pages/atlas/daily-nutrition-override-page';
import NutritionOverviewPage from './pages/atlas/nutrition-overview-page';
import ProductLibraryPage from './pages/atlas/product-library-page';
import HomePage from './pages/home';
import LoginPage from './pages/login-page';
import UiKitPage from './pages/ui-kit-page';
import UserDetailPage from './pages/user-detail-page';
import UsersPage from './pages/users-page';

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <I18nProvider>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route element={<ProtectedAdminLayout />}>
              <Route path="/" element={<HomePage />} />
              <Route path="/ui-kit" element={<UiKitPage />} />
              <Route path="/users" element={<UsersPage />} />
              <Route path="/users/:id" element={<UserDetailPage />} />
              <Route path="/atlas/nutrition" element={<NutritionOverviewPage />} />
              <Route
                path="/atlas/nutrition/overrides/new"
                element={<DailyNutritionOverridePage />}
              />
              <Route path="/atlas/nutrition/products" element={<ProductLibraryPage />} />
            </Route>
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </I18nProvider>
      </AuthProvider>
    </BrowserRouter>
  );
}
