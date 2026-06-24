// FILE: apps/web-admin/src/app/admin-navigation.test.ts
// VERSION: 1.1.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify web-admin app-owned navigation metadata and shell state resolution.
//   SCOPE: Covers route groups, active matching, breadcrumbs, disabled placeholders, and placeholder shell data; excludes rendering.
//   DEPENDS: apps/web-admin/src/app/admin-navigation.ts, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   admin navigation tests - Prove metadata is deterministic and app-owned for the shell.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.1.0 - Added Atlas nutrition daily-log navigation coverage.
// END_CHANGE_SUMMARY

import { describe, expect, it } from 'vitest';
import {
  adminReferenceItems,
  adminShellTeams,
  adminShellUser,
  resolveAdminShellState,
} from './admin-navigation';

describe('admin navigation metadata', () => {
  it('marks overview active and returns overview breadcrumbs for root', () => {
    const state = resolveAdminShellState('/');
    const emptyState = resolveAdminShellState('');

    expect(state.navigation[0].items.find((item) => item.id === 'overview')?.isActive).toBe(true);
    expect(emptyState.navigation[0].items.find((item) => item.id === 'overview')?.isActive).toBe(
      true,
    );
    expect(state.breadcrumbs).toEqual([{ label: 'Overview' }]);
    expect(emptyState.breadcrumbs).toEqual([{ label: 'Overview' }]);
  });

  it('marks users active for list and detail routes', () => {
    const listState = resolveAdminShellState('/users');
    const trailingSlashState = resolveAdminShellState('/users/');
    const detailState = resolveAdminShellState('/users/123');

    expect(listState.navigation[0].items.find((item) => item.id === 'users')?.isActive).toBe(true);
    expect(
      trailingSlashState.navigation[0].items.find((item) => item.id === 'users')?.isActive,
    ).toBe(true);
    expect(listState.breadcrumbs).toEqual([{ label: 'Users' }]);
    expect(trailingSlashState.breadcrumbs).toEqual([{ label: 'Users' }]);
    expect(detailState.navigation[0].items.find((item) => item.id === 'users')?.isActive).toBe(
      true,
    );
    expect(detailState.breadcrumbs).toEqual([
      { label: 'Users', href: '/users' },
      { label: 'User detail' },
    ]);
  });

  it('marks UI Kit active and keeps reference items disabled', () => {
    const state = resolveAdminShellState('/ui-kit');

    expect(state.navigation[0].items.find((item) => item.id === 'ui-kit')?.isActive).toBe(true);
    expect(state.breadcrumbs).toEqual([{ label: 'UI Kit' }]);
    expect(adminReferenceItems).toEqual([
      expect.objectContaining({ id: 'graphql-admin', disabled: true, name: 'GraphQL/Admin' }),
      expect.objectContaining({ id: 'system-settings', disabled: true, name: 'System/Settings' }),
    ]);
  });

  it('marks Atlas nutrition routes active and returns nutrition breadcrumbs', () => {
    const dailyState = resolveAdminShellState('/atlas/nutrition');
    const state = resolveAdminShellState('/atlas/nutrition/products');
    const dailyNutritionItem = dailyState.navigation[0].items.find(
      (item) => item.id === 'atlas-nutrition',
    );
    const nutritionItem = state.navigation[0].items.find((item) => item.id === 'atlas-nutrition');

    expect(dailyNutritionItem?.isActive).toBe(true);
    expect(dailyNutritionItem?.href).toBe('/atlas/nutrition');
    expect(
      dailyNutritionItem?.children?.find((child) => child.id === 'atlas-nutrition-daily-log'),
    ).toEqual(
      expect.objectContaining({
        href: '/atlas/nutrition',
        isActive: true,
        label: 'Daily Log',
      }),
    );
    expect(dailyState.breadcrumbs).toEqual([{ label: 'Nutrition' }]);
    expect(nutritionItem?.isActive).toBe(true);
    expect(
      nutritionItem?.children?.find((child) => child.id === 'atlas-nutrition-products'),
    ).toEqual(
      expect.objectContaining({
        href: '/atlas/nutrition/products',
        isActive: true,
        label: 'Product Library',
      }),
    );
    expect(state.breadcrumbs).toEqual([{ label: 'Nutrition' }, { label: 'Product Library' }]);
  });

  it('provides template-native user and team placeholders', () => {
    expect(adminShellUser).toEqual({
      name: 'Developer',
      email: 'developer@example.local',
      initials: 'DV',
    });
    expect(adminShellTeams[0]).toMatchObject({
      name: 'Monorepo Template',
      plan: 'Admin shell',
    });
  });
});
