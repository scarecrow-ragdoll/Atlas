// FILE: apps/web/e2e/rest-users-flow.spec.ts
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Verify the public Next REST browser user flow.
//   SCOPE: Covers creating a user through the UI, list refresh, selecting the user, and selected-detail rendering; excludes admin GraphQL behavior.
//   DEPENDS: @playwright/test, apps/web/e2e/playwright.config.ts, apps/api REST users endpoints.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   public REST web flow - Browser-level create/list/detail test for the public web app.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Scoped selected-user email assertion to match the migrated Next UI.
// END_CHANGE_SUMMARY

import { expect, test } from '@playwright/test';
import { uniqueEmail } from './helpers';

test('public REST web creates, lists, and opens a user detail', async ({ page }) => {
  const email = uniqueEmail('rest-web');
  const name = `REST Web User ${Date.now()}`;

  await page.goto('/');
  await page.getByPlaceholder('Name').fill(name);
  await page.getByPlaceholder('Email').fill(email);
  await page.getByPlaceholder('Password').fill('secret123');
  await page.getByRole('button', { name: 'Create' }).click();

  await expect(page.getByText(name)).toBeVisible();
  await expect(page.getByLabel('Users').getByText(email)).toBeVisible();
  await page.getByRole('button', { name }).click();
  await expect(page.getByLabel('Selected user').getByText(email)).toBeVisible();
});
