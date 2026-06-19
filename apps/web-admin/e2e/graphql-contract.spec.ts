// FILE: apps/web-admin/e2e/graphql-contract.spec.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify web-admin e2e API health and authenticated GraphQL user CRUD contracts.
//   SCOPE: Public health/ready endpoints and cookie-authenticated GraphQL user operations; excludes browser rendering and UI navigation.
//   DEPENDS: @playwright/test, apps/web-admin/e2e/helpers.ts, apps/api GraphQL schema.
//   LINKS: M-WEB-ADMIN / M-GRAPHQL-SCHEMA / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   health/readiness flow - Confirms public API status endpoints.
//   authenticated GraphQL CRUD flow - Confirms protected user operations through an admin session.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Switched protected GraphQL e2e calls to cookie-backed admin sessions.
// END_CHANGE_SUMMARY

import { expect, test } from '@playwright/test';
import {
  apiBaseURL,
  createUser,
  graphQL,
  uniqueEmail,
  withAuthenticatedGraphQLContext,
  type User,
} from './helpers';

test('health and readiness endpoints are available', async ({ request }) => {
  const health = await request.get(`${apiBaseURL}/healthz`);
  await expect(health).toBeOK();
  await expect(await health.json()).toEqual({ status: 'ok' });

  const ready = await request.get(`${apiBaseURL}/readyz`);
  await expect(ready).toBeOK();
  await expect(await ready.json()).toEqual({ status: 'ok' });
});

test('GraphQL user CRUD contract works against the real API', async () => {
  await withAuthenticatedGraphQLContext(async (context) => {
    const firstEmail = uniqueEmail('contract-primary');
    const secondEmail = uniqueEmail('contract-secondary');
    const created = await createUser(context, {
      email: firstEmail,
      name: 'Contract User',
    });
    const duplicateSeed = await createUser(context, {
      email: secondEmail,
      name: 'Duplicate Seed',
    });

    const list = await graphQL<{
      users: {
        edges: Array<{ node: User }>;
        totalCount: number;
        pageInfo: { hasNextPage: boolean; hasPreviousPage: boolean };
      };
    }>(
      context,
      `query ListUsers($first: Int) {
        users(pagination: { first: $first }) {
          edges {
            node {
              id
              email
              name
              createdAt
              updatedAt
            }
          }
          totalCount
          pageInfo {
            hasNextPage
            hasPreviousPage
          }
        }
      }`,
      { first: 50 },
    );
    expect(list.errors).toBeUndefined();
    expect(list.data?.users.edges.some(({ node }) => node.id === created.id)).toBe(true);

    const read = await graphQL<{ user: User | null }>(
      context,
      `query ReadUser($id: UUID!) {
        user(id: $id) {
          id
          email
          name
          createdAt
          updatedAt
        }
      }`,
      { id: created.id },
    );
    expect(read.data?.user).toMatchObject({ id: created.id, email: firstEmail });

    const updatedEmail = uniqueEmail('contract-updated');
    const update = await graphQL<{
      updateUser:
        | { __typename: 'UpdateUserSuccess'; user: User }
        | { __typename: 'ValidationError'; field: string; message: string }
        | { __typename: 'NotFoundError'; entityType: string; id: string }
        | { __typename: 'AuthError'; message: string };
    }>(
      context,
      `mutation UpdateUser($id: UUID!, $input: UpdateUserInput!) {
        updateUser(id: $id, input: $input) {
          __typename
          ... on UpdateUserSuccess {
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
          ... on NotFoundError {
            entityType
            id
          }
        }
      }`,
      { id: created.id, input: { name: 'Contract Updated', email: updatedEmail } },
    );
    expect(update.data?.updateUser).toMatchObject({
      __typename: 'UpdateUserSuccess',
      user: { id: created.id, email: updatedEmail, name: 'Contract Updated' },
    });

    const duplicate = await graphQL<{
      updateUser:
        | { __typename: 'UpdateUserSuccess'; user: User }
        | { __typename: 'ValidationError'; field: string; message: string }
        | { __typename: 'AuthError'; message: string };
    }>(
      context,
      `mutation DuplicateUpdate($id: UUID!, $input: UpdateUserInput!) {
        updateUser(id: $id, input: $input) {
          __typename
          ... on UpdateUserSuccess {
            user {
              id
            }
          }
          ... on ValidationError {
            field
            message
          }
        }
      }`,
      { id: created.id, input: { email: duplicateSeed.email } },
    );
    expect(duplicate.data?.updateUser).toEqual({
      __typename: 'ValidationError',
      field: 'email',
      message: 'already exists',
    });

    const deleted = await graphQL<{
      deleteUser:
        | { __typename: 'DeleteUserSuccess'; ok: boolean }
        | { __typename: 'AuthError'; message: string };
    }>(
      context,
      `mutation DeleteUser($id: UUID!) {
        deleteUser(id: $id) {
          __typename
          ... on DeleteUserSuccess {
            ok
          }
          ... on AuthError {
            message
          }
        }
      }`,
      { id: created.id },
    );
    expect(deleted.data?.deleteUser).toEqual({ __typename: 'DeleteUserSuccess', ok: true });

    const readDeleted = await graphQL<{ user: User | null }>(
      context,
      `query ReadDeleted($id: UUID!) {
        user(id: $id) {
          id
        }
      }`,
      { id: created.id },
    );
    expect(readDeleted.data?.user).toBeNull();

    const deleteMissing = await graphQL<{
      deleteUser:
        | { __typename: 'DeleteUserSuccess'; ok: boolean }
        | { __typename: 'AuthError'; message: string };
    }>(
      context,
      `mutation DeleteMissing($id: UUID!) {
        deleteUser(id: $id) {
          __typename
          ... on DeleteUserSuccess {
            ok
          }
          ... on AuthError {
            message
          }
        }
      }`,
      { id: '00000000-0000-0000-0000-000000000000' },
    );
    expect(deleteMissing.data?.deleteUser).toEqual({
      __typename: 'DeleteUserSuccess',
      ok: false,
    });
  });
});
