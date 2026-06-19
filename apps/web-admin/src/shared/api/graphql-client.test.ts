// FILE: apps/web-admin/src/shared/api/graphql-client.test.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Prove the web-admin GraphQL transport sends cookie-backed admin sessions.
//   SCOPE: Covers GraphQLClient construction and absence of stale bearer-token helpers; excludes generated operation documents and page-level request behavior.
//   DEPENDS: vitest, graphql-request, apps/web-admin/src/shared/api/graphql-client.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   graphql client tests - Verifies cookie credentials and no bearer mutation helper.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Switched transport expectations from bearer headers to cookie credentials.
// END_CHANGE_SUMMARY

import { describe, expect, it } from 'vitest';
import * as GraphqlClientModule from './graphql-client';

const { createGraphQLClient } = GraphqlClientModule;

describe('graphql client', () => {
  it('creates a GraphQLClient for the supplied URL with cookie credentials', () => {
    const client = createGraphQLClient('http://example.test/graphql') as unknown as {
      requestConfig: { credentials?: RequestCredentials; headers?: Record<string, string> };
    };

    expect(client).toBeDefined();
    expect(client.requestConfig.credentials).toBe('include');
    expect(client.requestConfig.headers?.Authorization).toBeUndefined();
  });

  it('does not expose a bearer-token mutation helper', () => {
    expect('setAuthToken' in GraphqlClientModule).toBe(false);
  });
});
