// FILE: apps/web/app/page.tsx
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Render the public Next users page.
//   SCOPE: Fetches initial REST users data on the server and passes it to the interactive client component; excludes mutation behavior and route proxy internals.
//   DEPENDS: apps/web/app/users-client.tsx, apps/web/src/shared/api/users.ts.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   dynamic - Opts the REST-backed root page out of build-time prerendering.
//   default - Server-first public users route.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.2 - Render the users fallback when the server REST fetch fails at runtime.
// END_CHANGE_SUMMARY

import UsersClient from './users-client';
import { listUsers } from '../src/shared/api/users';

export const dynamic = 'force-dynamic';

export default async function Page() {
  try {
    const { users, totalCount } = await listUsers();

    return <UsersClient initialTotalCount={totalCount} initialUsers={users} />;
  } catch {
    return <UsersClient initialLoadError initialTotalCount={0} initialUsers={[]} />;
  }
}
