// FILE: apps/web-admin/src/shared/api/graphql-client.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Construct the web-admin GraphQL transport for cookie-backed admin sessions.
//   SCOPE: Owns GraphQLClient configuration for browser requests; excludes operation documents, generated types, and page request orchestration.
//   DEPENDS: graphql-request, apps/web-admin/src/shared/config/index.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   createGraphQLClient - Returns a GraphQLClient configured to include admin session cookies.
//   graphqlClient - Shared web-admin GraphQLClient instance.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Replaced bearer-token header mutation with cookie credential transport.
// END_CHANGE_SUMMARY

import { GraphQLClient } from 'graphql-request';
import { appConfig } from '@shared/config';

export function createGraphQLClient(apiUrl = appConfig.apiUrl) {
  return new GraphQLClient(apiUrl, { credentials: 'include' });
}

export const graphqlClient = createGraphQLClient();
