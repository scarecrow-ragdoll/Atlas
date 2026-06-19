// FILE: apps/web-admin/e2e/helpers.ts
// VERSION: 1.1.1
// START_MODULE_CONTRACT
//   PURPOSE: Provide Playwright helpers for authenticated web-admin GraphQL e2e flows.
//   SCOPE: Unique test data, API-bound GraphQL requests, API admin login, browser UI login, and user seed helpers; excludes scenario assertions, browser cookie shortcuts, and Playwright server startup.
//   DEPENDS: @playwright/test, apps/api GraphQL admin auth and users schema.
//   LINKS: M-WEB-ADMIN / M-GRAPHQL-SCHEMA / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   apiBaseURL - API base URL used by API-bound Playwright request contexts.
//   adminEmail - Test admin email supplied by Playwright environment.
//   adminPassword - Test admin password supplied by Playwright environment.
//   graphqlURL - Full GraphQL endpoint URL for direct API assertions.
//   uniqueEmail - Builds process/time-scoped test email addresses.
//   withGraphQLContext - Runs API-bound GraphQL calls without installing an admin session.
//   graphQL - Posts one GraphQL document to the API and returns the decoded payload.
//   loginAdmin - Authenticates the seeded e2e admin through API-bound GraphQL.
//   loginThroughUi - Logs in through the real browser login page.
//   withAuthenticatedGraphQLContext - Runs API-bound GraphQL calls with an admin session.
//   User - Shared e2e user shape returned by GraphQL operations.
//   createUser - Creates one user through the authenticated GraphQL API.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.1.1 - Removed browser cookie-install helper and redacted cookie assertion output.
// END_CHANGE_SUMMARY

import {
  expect,
  request as playwrightRequest,
  type APIRequestContext,
  type Page,
} from '@playwright/test';

const apiPort = process.env.E2E_API_PORT ?? '18080';

export const apiBaseURL = process.env.E2E_API_URL ?? `http://localhost:${apiPort}`;
export const graphqlURL = `${apiBaseURL}/graphql`;
const webPort = process.env.E2E_WEB_PORT ?? '13000';
const webOrigin = process.env.E2E_WEB_URL ?? `http://localhost:${webPort}`;
export const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'e2e-admin@example.test';
export const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'StrongPassword123!';

interface GraphQLPayload<TData> {
  data?: TData;
  errors?: Array<{ message: string }>;
}

export function uniqueEmail(prefix: string) {
  const normalized = prefix.replace(/[^a-z0-9-]/gi, '-').toLowerCase();
  return `${normalized}-${process.pid}-${Date.now()}@example.test`;
}

export async function withGraphQLContext<T>(fn: (context: APIRequestContext) => Promise<T>) {
  const context = await playwrightRequest.newContext({ baseURL: apiBaseURL });
  try {
    return await fn(context);
  } finally {
    await context.dispose();
  }
}

export async function graphQL<TData>(
  context: APIRequestContext,
  query: string,
  variables: Record<string, unknown> = {},
) {
  const response = await context.post('/graphql', {
    data: { query, variables },
    headers: { 'Content-Type': 'application/json' },
  });
  expect(response.ok()).toBeTruthy();
  return (await response.json()) as GraphQLPayload<TData>;
}

export async function loginAdmin(apiContext: APIRequestContext) {
  const response = await apiContext.post('/graphql', {
    data: {
      query: `mutation LoginAdmin($input: LoginAdminInput!) {
        loginAdmin(input: $input) {
          __typename
          ... on LoginAdminSuccess {
            admin {
              id
              email
              role
            }
          }
          ... on AuthError {
            message
          }
        }
      }`,
      variables: {
        input: {
          email: adminEmail,
          password: adminPassword,
        },
      },
    },
    headers: { 'Content-Type': 'application/json', Origin: webOrigin },
  });

  expect(response.ok()).toBeTruthy();
  const payload = (await response.json()) as GraphQLPayload<{
    loginAdmin:
      | { __typename: 'LoginAdminSuccess'; admin: { id: string; email: string; role: string } }
      | { __typename: 'AuthError'; message: string };
  }>;
  expect(payload.errors).toBeUndefined();
  expect(payload.data?.loginAdmin.__typename).toBe('LoginAdminSuccess');
  const hasAdminSessionCookie = (response.headers()['set-cookie'] ?? '').includes(
    'web_admin_session=',
  );
  expect(hasAdminSessionCookie).toBe(true);
}

// START_CONTRACT: loginThroughUi
//   PURPOSE: Authenticate the browser through the real public login route and backend cookie session.
//   INPUTS: { page: Page - Playwright page, from: string - same-app return path }
//   OUTPUTS: none.
//   SIDE_EFFECTS: Navigates the browser, submits credentials, and stores the httpOnly admin cookie in the browser context.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: loginThroughUi
export async function loginThroughUi(page: Page, from = '/users') {
  await page.goto(`/login?from=${encodeURIComponent(from)}`);
  await page.getByLabel('Email').fill(adminEmail);
  await page.getByLabel('Password').fill(adminPassword);
  await page.getByRole('button', { name: 'Sign in' }).click();
  await expect(page).toHaveURL((url) => `${url.pathname}${url.search}${url.hash}` === from);
}

export async function withAuthenticatedGraphQLContext<T>(
  fn: (context: APIRequestContext) => Promise<T>,
) {
  const context = await playwrightRequest.newContext({
    baseURL: apiBaseURL,
    extraHTTPHeaders: { Origin: webOrigin },
  });
  try {
    await loginAdmin(context);
    return await fn(context);
  } finally {
    await context.dispose();
  }
}

export interface User {
  id: string;
  email: string;
  name: string;
  createdAt: string;
  updatedAt: string;
}

export async function createUser(
  context: APIRequestContext,
  input: { email: string; name?: string; password?: string },
) {
  const result = await graphQL<{
    createUser:
      | { __typename: 'CreateUserSuccess'; user: User }
      | { __typename: 'ValidationError'; field: string; message: string };
  }>(
    context,
    `mutation CreateE2EUser($input: CreateUserInput!) {
      createUser(input: $input) {
        __typename
        ... on CreateUserSuccess {
          user {
            id
            email
            name
            createdAt
            updatedAt
          }
        }
        ... on ValidationError {
          field
          message
        }
      }
    }`,
    {
      input: {
        email: input.email,
        name: input.name ?? 'E2E User',
        password: input.password ?? 'Password123!',
      },
    },
  );

  expect(result.errors).toBeUndefined();
  expect(result.data?.createUser.__typename).toBe('CreateUserSuccess');
  const created = result.data?.createUser;
  if (!created || created.__typename !== 'CreateUserSuccess') {
    throw new Error('createUser did not return CreateUserSuccess');
  }
  return created.user;
}
