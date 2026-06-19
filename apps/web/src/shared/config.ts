// FILE: apps/web/src/shared/config.ts
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Provide the public web REST API base configuration.
//   SCOPE: Resolves same-origin browser REST calls and server/runtime WEB_API_BASE_URL for Next route handlers; excludes API request/response parsing.
//   DEPENDS: Next.js browser/server runtime, process.env.WEB_API_BASE_URL.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: CONFIG
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   appConfig - REST API base URL for browser and server/runtime execution.
//   resolveServerApiBaseUrl - Runtime-only REST API base URL resolver for Next route handlers and server components.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Moved public web config from Vite env to Next same-origin plus WEB_API_BASE_URL runtime env.
// END_CHANGE_SUMMARY

function normalizeApiBaseUrl(value: string | undefined): string {
  const baseUrl = value?.trim() || 'http://localhost:8090';
  return baseUrl.replace(/\/+$/, '');
}

export function resolveServerApiBaseUrl(): string {
  return normalizeApiBaseUrl(process.env.WEB_API_BASE_URL);
}

function resolveApiBaseUrl(): string {
  if (typeof window !== 'undefined') {
    return '';
  }

  return resolveServerApiBaseUrl();
}

export const appConfig = {
  apiBaseUrl: resolveApiBaseUrl(),
};
