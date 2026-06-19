// FILE: apps/web-admin/src/entities/admin-auth/model.test.ts
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Verify web-admin auth model helpers.
//   SCOPE: Covers admin initials, shell-user mapping, and safe same-app return-to parsing; excludes GraphQL transport and route rendering.
//   DEPENDS: apps/web-admin/src/entities/admin-auth/model.ts, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   auth model tests - Prove initials, shell user mapping, and return-to safety rules.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Added malformed URL fallback coverage for safe return paths.
// END_CHANGE_SUMMARY

import { describe, expect, it, vi } from 'vitest';
import {
  adminToShellUser,
  getAdminInitials,
  resolveSafeReturnTo,
  type AdminPrincipal,
} from './model';

const admin: AdminPrincipal = {
  id: 'admin-1',
  email: 'owner@example.test',
  name: 'Owner Admin',
  role: 'ADMIN',
  createdAt: '2026-06-07T00:00:00Z',
  updatedAt: '2026-06-07T00:00:00Z',
};

describe('admin auth model helpers', () => {
  it('derives readable initials from admin names and emails', () => {
    expect(getAdminInitials('Owner Admin', 'owner@example.test')).toBe('OA');
    expect(getAdminInitials('Owner', 'owner@example.test')).toBe('O');
    expect(getAdminInitials('', 'support@example.test')).toBe('S');
    expect(getAdminInitials('', '@example.test')).toBe('A');
  });

  it('maps the current admin to the sidebar user contract', () => {
    expect(adminToShellUser(admin)).toEqual({
      name: 'Owner Admin',
      email: 'owner@example.test',
      initials: 'OA',
    });
  });

  it('accepts same-app return paths with search and hash', () => {
    expect(resolveSafeReturnTo('/users?status=active#row-2')).toBe('/users?status=active#row-2');
    expect(resolveSafeReturnTo('/ui-kit')).toBe('/ui-kit');
  });

  it('rejects unsafe or login-loop return paths', () => {
    expect(resolveSafeReturnTo(undefined)).toBe('/');
    expect(resolveSafeReturnTo('https://evil.example/users')).toBe('/');
    expect(resolveSafeReturnTo('//evil.example/users')).toBe('/');
    expect(resolveSafeReturnTo('/\\evil.example/users')).toBe('/');
    expect(resolveSafeReturnTo('users')).toBe('/');
    expect(resolveSafeReturnTo('/login')).toBe('/');
    expect(resolveSafeReturnTo('/login/')).toBe('/');
    expect(resolveSafeReturnTo('/login?from=/users')).toBe('/');
  });

  it('falls back when browser URL parsing rejects a same-app-looking path', () => {
    const OriginalURL = globalThis.URL;
    vi.stubGlobal(
      'URL',
      class ThrowingURL {
        constructor() {
          throw new TypeError('malformed URL');
        }
      },
    );

    try {
      expect(resolveSafeReturnTo('/users')).toBe('/');
    } finally {
      vi.stubGlobal('URL', OriginalURL);
    }
  });

  it('falls back when browser URL parsing produces an empty candidate path', () => {
    const OriginalURL = globalThis.URL;
    vi.stubGlobal(
      'URL',
      class EmptyPathURL {
        hash = '';
        origin = window.location.origin;
        pathname = '';
        search = '';
      },
    );

    try {
      expect(resolveSafeReturnTo('/empty')).toBe('/');
    } finally {
      vi.stubGlobal('URL', OriginalURL);
    }
  });
});
