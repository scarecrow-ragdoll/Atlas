// FILE: apps/web-admin/src/entities/admin-auth/provider.tsx
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Provide frontend admin-auth state and actions for web-admin routes.
//   SCOPE: Owns React Query current-admin state, login refetch, logout cache clearing, and context access; excludes route navigation and page-specific UI.
//   DEPENDS: react, @tanstack/react-query, apps/web-admin/src/entities/admin-auth/client.ts, apps/web-admin/src/entities/admin-auth/model.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   adminAuthQueryKey - Stable React Query key for current admin.
//   isProtectedAdminQueryKey - Identifies protected admin route data that must be dropped on logout.
//   AuthProvider - Navigation-free auth state provider.
//   useAdminAuth - Context hook for route and shell components.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Added one-shot completed logout state for plain /login redirects.
// END_CHANGE_SUMMARY

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { createContext, useContext, useState, type ReactNode } from 'react';
import {
  fetchCurrentAdmin,
  loginAdmin as requestLoginAdmin,
  logoutAdmin as requestLogoutAdmin,
  type LoginAdminResult,
} from './client';
import type { AdminPrincipal, LoginCredentials } from './model';

export const adminAuthQueryKey = ['admin-auth', 'current-admin'] as const;

export function isProtectedAdminQueryKey(queryKey: readonly unknown[]) {
  return queryKey[0] === 'admin-users' || queryKey[0] === 'admin-user';
}

type AuthContextValue = {
  admin: AdminPrincipal | null;
  clearCompletedLogout: () => void;
  hasCompletedLogout: boolean;
  isAuthenticated: boolean;
  isLoading: boolean;
  isLogoutPending: boolean;
  login: (credentials: LoginCredentials) => Promise<LoginAdminResult>;
  logout: () => Promise<void>;
  refetchCurrentAdmin: () => Promise<AdminPrincipal | null>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

// START_CONTRACT: AuthProvider
//   PURPOSE: Expose current admin auth state and login/logout actions without owning route navigation.
//   INPUTS: { children: ReactNode - route tree content }
//   OUTPUTS: { JSX.Element - auth context provider }
//   SIDE_EFFECTS: Sends CurrentAdmin/LoginAdmin/LogoutAdmin GraphQL requests through React Query actions.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: AuthProvider
export function AuthProvider({ children }: { children: ReactNode }) {
  const queryClient = useQueryClient();
  const [hasCompletedLogout, setHasCompletedLogout] = useState(false);
  const currentAdminQuery = useQuery({
    queryKey: adminAuthQueryKey,
    queryFn: fetchCurrentAdmin,
  });
  const logoutMutation = useMutation({ mutationFn: requestLogoutAdmin });

  async function refetchCurrentAdmin() {
    const result = await currentAdminQuery.refetch();
    return result.data ?? null;
  }

  async function login(credentials: LoginCredentials) {
    const result = await requestLoginAdmin(credentials);
    if (result.ok) {
      setHasCompletedLogout(false);
      await refetchCurrentAdmin();
    }
    return result;
  }

  async function logout() {
    await logoutMutation.mutateAsync();
    await queryClient.cancelQueries();
    setHasCompletedLogout(true);
    queryClient.setQueryData(adminAuthQueryKey, null);
    queryClient.removeQueries({
      predicate: (query) => isProtectedAdminQueryKey(query.queryKey),
    });
  }

  function clearCompletedLogout() {
    setHasCompletedLogout(false);
  }

  return (
    <AuthContext.Provider
      value={{
        admin: currentAdminQuery.data ?? null,
        clearCompletedLogout,
        hasCompletedLogout,
        isAuthenticated: Boolean(currentAdminQuery.data),
        isLoading: currentAdminQuery.isPending,
        isLogoutPending: logoutMutation.isPending,
        login,
        logout,
        refetchCurrentAdmin,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

// START_CONTRACT: useAdminAuth
//   PURPOSE: Read web-admin auth context from route, page, and shell components.
//   INPUTS: none.
//   OUTPUTS: { AuthContextValue - current admin state and actions }
//   SIDE_EFFECTS: Throws if used outside AuthProvider.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: useAdminAuth
export function useAdminAuth(): AuthContextValue {
  const value = useContext(AuthContext);
  if (!value) {
    throw new Error('useAdminAuth must be used within AuthProvider');
  }
  return value;
}
