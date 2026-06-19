// FILE: apps/web/app/api/users/[id]/route.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Expose same-origin public REST detail, update, and delete routes in Next.
//   SCOPE: Proxies GET, PATCH, and DELETE `/api/users/:id` requests to the Go API using runtime WEB_API_BASE_URL; excludes root list/create operations.
//   DEPENDS: apps/web/app/api/users/proxy.ts.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   GET - Proxy one-user read requests.
//   PATCH - Proxy one-user update requests.
//   DELETE - Proxy one-user delete requests.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added public Next REST users id-scoped proxy route.
// END_CHANGE_SUMMARY

import { proxyUsersRequest } from '../proxy';

type RouteContext = {
  params: Promise<{ id: string }>;
};

async function userId(context: RouteContext): Promise<string> {
  return (await context.params).id;
}

export async function GET(_request: Request, context: RouteContext) {
  return proxyUsersRequest({ id: await userId(context), method: 'GET' });
}

export async function PATCH(request: Request, context: RouteContext) {
  return proxyUsersRequest({
    body: await request.text(),
    id: await userId(context),
    method: 'PATCH',
  });
}

export async function DELETE(_request: Request, context: RouteContext) {
  return proxyUsersRequest({ id: await userId(context), method: 'DELETE' });
}
