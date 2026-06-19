// FILE: apps/web-admin/src/shared/config/index.ts
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Provide the web-admin browser configuration values.
//   SCOPE: Resolves Vite browser environment values for the GraphQL API URL and app name; excludes server-only configuration.
//   DEPENDS: Vite import.meta.env.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: CONFIG
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   appConfig - Normalized web-admin app name and GraphQL endpoint.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Moved web-admin config from Next public env to Vite browser env.
// END_CHANGE_SUMMARY

function normalizeUrl(value: string | undefined, fallback: string): string {
  const url = value?.trim() || fallback;
  return url.replace(/\/+$/, '');
}

export const appConfig = {
  apiUrl: normalizeUrl(import.meta.env.VITE_GRAPHQL_API_URL, 'http://localhost:8090/graphql'),
  appName: import.meta.env.VITE_APP_NAME?.trim() || 'MonorepoApp',
} as const;
