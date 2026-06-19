// FILE: apps/web/app/api/users/route.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Expose same-origin public REST list and create routes in Next.
//   SCOPE: Proxies GET and POST `/api/users` requests to the Go API using runtime WEB_API_BASE_URL; excludes id-scoped operations.
//   DEPENDS: apps/web/app/api/users/proxy.ts.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   GET - Proxy public users list requests.
//   POST - Proxy public users create requests.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added public Next REST users root proxy route.
// END_CHANGE_SUMMARY

import { proxyUsersRequest } from './proxy';

export async function GET() {
  return proxyUsersRequest({ method: 'GET' });
}

export async function POST(request: Request) {
  return proxyUsersRequest({ body: await request.text(), method: 'POST' });
}
