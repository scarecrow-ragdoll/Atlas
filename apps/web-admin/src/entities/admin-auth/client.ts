// FILE: apps/web-admin/src/entities/admin-auth/client.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide typed web-admin GraphQL client helpers for admin authentication.
//   SCOPE: Owns current-admin, login, and logout requests plus auth union normalization; excludes React Query context and route navigation.
//   DEPENDS: apps/web-admin/src/entities/admin-auth/api/*.graphql, apps/web-admin/src/shared/api/graphql-client.ts, apps/web-admin/src/shared/api/generated/types.ts, apps/web-admin/src/entities/admin-auth/model.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   fetchCurrentAdmin - Request the current backend admin principal.
//   loginAdmin - Request backend login and normalize success/error result unions.
//   logoutAdmin - Request backend logout and return whether logout succeeded.
//   LoginAdminResult - Frontend normalized login result.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added admin-auth GraphQL client helpers for login route guard behavior.
// END_CHANGE_SUMMARY

import { graphqlClient } from '@shared/api/graphql-client';
import type {
  CurrentAdminQuery,
  LoginAdminMutation,
  LoginAdminMutationVariables,
  LogoutAdminMutation,
} from '@shared/api/generated/types';
import currentAdminQueryDocument from './api/currentAdmin.graphql?raw';
import loginAdminMutationDocument from './api/loginAdmin.graphql?raw';
import logoutAdminMutationDocument from './api/logoutAdmin.graphql?raw';
import type { AdminPrincipal, AuthMutationError, LoginCredentials } from './model';

export type LoginAdminResult =
  | { ok: true; admin: AdminPrincipal }
  | { ok: false; error: AuthMutationError };

// START_CONTRACT: fetchCurrentAdmin
//   PURPOSE: Load the current admin principal from the backend session cookie.
//   INPUTS: none.
//   OUTPUTS: { Promise<AdminPrincipal | null> - current admin or null when unauthenticated }
//   SIDE_EFFECTS: Sends a credentialed GraphQL request through graphqlClient.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
// END_CONTRACT: fetchCurrentAdmin
export async function fetchCurrentAdmin(): Promise<AdminPrincipal | null> {
  const response = await graphqlClient.request<CurrentAdminQuery>(currentAdminQueryDocument);
  return response.me ?? null;
}

// START_CONTRACT: loginAdmin
//   PURPOSE: Authenticate an admin and normalize backend auth result unions for UI consumption.
//   INPUTS: { credentials: LoginCredentials - email and password from the login form }
//   OUTPUTS: { Promise<LoginAdminResult> - normalized success or user-visible error }
//   SIDE_EFFECTS: Sends a credentialed GraphQL request; backend may set the httpOnly session cookie.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
// END_CONTRACT: loginAdmin
export async function loginAdmin(credentials: LoginCredentials): Promise<LoginAdminResult> {
  const response = await graphqlClient.request<LoginAdminMutation, LoginAdminMutationVariables>(
    loginAdminMutationDocument,
    { input: credentials },
  );

  const result = response.loginAdmin;

  if (result.__typename === 'LoginAdminSuccess') {
    return { ok: true, admin: result.admin };
  }

  if (result.__typename === 'ValidationError') {
    return { ok: false, error: { field: result.field, message: result.message } };
  }

  return { ok: false, error: { message: result.message } };
}

// START_CONTRACT: logoutAdmin
//   PURPOSE: Revoke the current admin session through the backend auth contract.
//   INPUTS: none.
//   OUTPUTS: { Promise<boolean> - true when backend reports logout success }
//   SIDE_EFFECTS: Sends a credentialed GraphQL request; backend clears the httpOnly session cookie.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
// END_CONTRACT: logoutAdmin
export async function logoutAdmin(): Promise<boolean> {
  const response = await graphqlClient.request<LogoutAdminMutation>(logoutAdminMutationDocument);
  return response.logoutAdmin.ok;
}
