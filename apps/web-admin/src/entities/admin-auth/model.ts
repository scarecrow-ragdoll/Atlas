// FILE: apps/web-admin/src/entities/admin-auth/model.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Define frontend admin auth model helpers for the web-admin app.
//   SCOPE: Owns current-admin principal typing, initials derivation, sidebar user mapping, and safe same-app return-to parsing; excludes GraphQL transport and React context.
//   DEPENDS: apps/web-admin/src/shared/api/generated/types.ts, apps/web-admin/src/shared/ui.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
//   ROLE: TYPES
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminPrincipal - Frontend current-admin principal shape.
//   LoginCredentials - Login form input shape.
//   AuthMutationError - Normalized login mutation error.
//   getAdminInitials - Derive sidebar initials from current admin data.
//   adminToShellUser - Map current admin to the shared shell user contract.
//   resolveSafeReturnTo - Accept only safe same-app return paths.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added auth model helpers for login route guard behavior.
// END_CHANGE_SUMMARY

import type { CurrentAdminQuery } from '@shared/api/generated/types';
import type { AdminUser as ShellAdminUser } from '@shared/ui';

export type AdminPrincipal = NonNullable<CurrentAdminQuery['me']>;

export type LoginCredentials = {
  email: string;
  password: string;
};

export type AuthMutationError = {
  message: string;
  field?: string;
};

// START_CONTRACT: getAdminInitials
//   PURPOSE: Derive short stable initials for sidebar avatar display.
//   INPUTS: { name: string - admin display name, email: string - admin email fallback }
//   OUTPUTS: { string - one or two uppercase initials }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: getAdminInitials
export function getAdminInitials(name: string, email: string): string {
  const nameParts = name.trim().split(/\s+/).filter(Boolean);
  const source = nameParts.length > 0 ? nameParts : [email.split('@')[0] || 'A'];

  return source
    .slice(0, 2)
    .map((part) => part.charAt(0).toUpperCase())
    .join('');
}

// START_CONTRACT: adminToShellUser
//   PURPOSE: Map backend current-admin data into the shared sidebar user contract.
//   INPUTS: { admin: AdminPrincipal - current backend admin }
//   OUTPUTS: { ShellAdminUser - sidebar display user }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: adminToShellUser
export function adminToShellUser(admin: AdminPrincipal): ShellAdminUser {
  return {
    name: admin.name,
    email: admin.email,
    initials: getAdminInitials(admin.name, admin.email),
  };
}

// START_CONTRACT: resolveSafeReturnTo
//   PURPOSE: Accept safe same-app return paths and reject external or login-loop redirects.
//   INPUTS: { rawValue: string | null | undefined - untrusted from query parameter }
//   OUTPUTS: { string - safe app-relative path }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: resolveSafeReturnTo
export function resolveSafeReturnTo(rawValue: string | null | undefined): string {
  const fallback = '/';
  const value = rawValue?.trim();

  if (!value || !value.startsWith('/') || value.startsWith('//') || value.includes('\\')) {
    return fallback;
  }

  try {
    const parsed = new URL(value, window.location.origin);
    const candidate = `${parsed.pathname}${parsed.search}${parsed.hash}`;

    if (
      parsed.origin !== window.location.origin ||
      parsed.pathname === '/login' ||
      parsed.pathname.startsWith('/login/')
    ) {
      return fallback;
    }

    return candidate || fallback;
  } catch {
    return fallback;
  }
}
