// FILE: apps/web/app/api/users/proxy.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Share public Next REST route proxy helpers.
//   SCOPE: Builds API URLs, forwards methods/bodies/headers, and preserves upstream response status/body; excludes individual route exports.
//   DEPENDS: apps/web/src/shared/config.ts, Next route handler runtime.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   proxyUsersRequest - Forward a REST request to the Go API users endpoint.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added shared proxy helper for public Next REST route handlers.
// END_CHANGE_SUMMARY

import { resolveServerApiBaseUrl } from '../../../src/shared/config';

type ProxyOptions = {
  body?: string;
  id?: string;
  method: string;
};

function usersUrl(id?: string): string {
  const baseUrl = resolveServerApiBaseUrl();
  const suffix = id ? `/${encodeURIComponent(id)}` : '';
  return `${baseUrl}/api/users${suffix}`;
}

export async function proxyUsersRequest({ body, id, method }: ProxyOptions): Promise<Response> {
  const upstream = await fetch(usersUrl(id), {
    body,
    headers: body
      ? {
          Accept: 'application/json',
          'Content-Type': 'application/json',
        }
      : {
          Accept: 'application/json',
        },
    method,
  });

  if (upstream.status === 204) {
    return new Response(null, { status: 204 });
  }

  return new Response(await upstream.text(), {
    headers: {
      'Content-Type': upstream.headers.get('Content-Type') || 'application/json',
    },
    status: upstream.status,
  });
}
