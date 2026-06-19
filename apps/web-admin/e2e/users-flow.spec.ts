// FILE: apps/web-admin/e2e/users-flow.spec.ts
// VERSION: 1.1.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify admin users browser flows through the Vite UI, real login/logout, and real GraphQL API.
//   SCOPE: Covers protected-route login/logout, create/list/detail navigation, duplicate-email validation, and shell navigation in the admin app; excludes public REST web flows and API unit contracts.
//   DEPENDS: @playwright/test, apps/web-admin/e2e/helpers.ts, apps/web-admin/e2e/playwright.config.ts, apps/api GraphQL users API.
//   LINKS: M-WEB-ADMIN / M-GRAPHQL-SCHEMA / V-M-WEB-ADMIN
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   protected auth flow - Verifies unauthenticated redirects, real browser login, logout, and post-logout protection.
//   create/list/detail flow - Creates a user through the UI, verifies table refresh, and opens the detail route.
//   duplicate-email flow - Seeds a duplicate through GraphQL and verifies the admin UI validation error.
//   desktop shell flow - Verifies sidebar navigation to /ui-kit and the collapsed icon rail.
//   mobile shell flow - Verifies mobile sidebar sheet navigation and route usability.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.1.0 - Replaced cookie installation with real browser login/logout coverage.
// END_CHANGE_SUMMARY

import { expect, test } from '@playwright/test';
import {
  createUser,
  loginThroughUi,
  uniqueEmail,
  withAuthenticatedGraphQLContext,
} from './helpers';

test('protected routes require login and logout revokes browser access', async ({
  context,
  page,
}) => {
  await page.goto('/users');
  await expect(page).toHaveURL(/\/login\?from=%2Fusers$/);
  await expect(page.getByRole('heading', { name: 'Admin sign in' })).toBeVisible();

  await loginThroughUi(page, '/users');
  await expect(page.getByRole('heading', { level: 1, name: 'Users' })).toBeVisible();

  await page.getByRole('button', { name: /E2E Admin/ }).click();
  await page.getByRole('menuitem', { name: 'Logout' }).click();
  await expect(page).toHaveURL(/\/login$/);

  const freshPage = await context.newPage();
  await freshPage.goto('/users');
  await expect(freshPage).toHaveURL(/\/login\?from=%2Fusers$/);
  await expect(freshPage.getByRole('heading', { name: 'Admin sign in' })).toBeVisible();
});

test('users page creates, lists, and opens a user detail page', async ({ page }) => {
  const email = uniqueEmail('browser-create');

  await loginThroughUi(page, '/users');
  await page.getByPlaceholder('Name').fill('Browser User');
  await page.getByPlaceholder('Email').fill(email);
  await page.getByPlaceholder('Password').fill('Password123!');
  await page.getByRole('button', { name: 'Create' }).click();

  const createdRow = page.getByRole('row').filter({ hasText: email });
  const userLink = createdRow.getByRole('link', { name: 'Browser User' });
  await expect(userLink).toBeVisible();
  await expect(createdRow.getByRole('cell', { name: email })).toBeVisible();

  await userLink.click();
  await expect(page.getByRole('heading', { name: 'Browser User' })).toBeVisible();
  await expect(page.getByText(email)).toBeVisible();
  await expect(page.getByTestId('admin-page-shell').getByText('ID', { exact: true })).toBeVisible();
});

test('users page shows duplicate email validation errors', async ({ page }) => {
  const email = uniqueEmail('browser-duplicate');
  await withAuthenticatedGraphQLContext(async (context) => {
    await createUser(context, { email, name: 'Duplicate Browser Seed' });
  });

  await loginThroughUi(page, '/users');
  await page.getByPlaceholder('Name').fill('Duplicate Browser User');
  await page.getByPlaceholder('Email').fill(email);
  await page.getByPlaceholder('Password').fill('Password123!');
  await page.getByRole('button', { name: 'Create' }).click();

  await expect(page.getByText('email: already exists')).toBeVisible();
});

test('shell navigation reaches UI kit and collapses to icon rail', async ({ page }) => {
  await page.setViewportSize({ width: 1280, height: 900 });
  await loginThroughUi(page, '/users');

  const navigation = page.getByRole('navigation', { name: 'Admin navigation' });
  await expect(navigation).toBeVisible();
  await navigation.getByRole('link', { name: 'UI Kit' }).click();

  await expect(page).toHaveURL(/\/ui-kit$/);
  await expect(page.getByRole('heading', { name: 'UI Kit' })).toBeVisible();
  await expect(page.getByRole('heading', { name: 'Shell Foundation' })).toBeVisible();
  await expect(navigation.getByRole('link', { name: 'UI Kit' })).toHaveAttribute(
    'data-active',
    'true',
  );

  const shell = page.locator('[data-slot="sidebar-wrapper"]').first();
  await expect(shell).toHaveAttribute('data-state', 'expanded');

  await page.locator('[data-slot="sidebar-trigger"]').first().click();
  await expect(shell).toHaveAttribute('data-state', 'collapsed');
  await expect
    .poll(async () => {
      return page
        .locator('[data-slot="sidebar-container"]')
        .first()
        .evaluate((element) => Math.round(element.getBoundingClientRect().width));
    })
    .toBeLessThanOrEqual(80);
});

test('mobile shell opens sidebar sheet navigation and returns to usable content', async ({
  page,
}) => {
  await page.setViewportSize({ width: 390, height: 844 });
  await loginThroughUi(page, '/users');

  await page.getByRole('button', { name: 'Toggle sidebar' }).click();
  const mobileSidebar = page.getByRole('dialog', { name: 'Sidebar' });
  await expect(mobileSidebar).toBeVisible();

  await mobileSidebar.getByRole('link', { name: 'UI Kit' }).click();

  await expect(page).toHaveURL(/\/ui-kit$/);
  await expect(mobileSidebar).toBeHidden();
  await expect(page.getByRole('heading', { name: 'UI Kit' })).toBeVisible();
  await expect(page.getByRole('button', { name: 'Toggle sidebar' })).toBeVisible();
});
