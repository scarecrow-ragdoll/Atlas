// FILE: apps/web-admin/src/pages/user-detail-page.tsx
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Render the web-admin user detail route.
//   SCOPE: Loads one user by route id and displays loading, error, not-found, and detail states through UI-kit components; excludes list and mutation behavior.
//   DEPENDS: @tanstack/react-query, react-router, apps/web-admin/src/entities/user/api/user.graphql, apps/web-admin/src/shared/api/graphql-client.ts, generated GraphQL types, apps/web-admin/src/shared/ui.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - User detail route backed by the GetUser GraphQL document and UI-kit states.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Migrated user detail route visuals to the web-admin UI kit.
// END_CHANGE_SUMMARY

import { useQuery } from '@tanstack/react-query';
import getUserQueryDocument from '@entities/user/api/user.graphql?raw';
import { graphqlClient } from '@shared/api/graphql-client';
import type { GetUserQuery } from '@shared/api/generated/types';
import { Link, useParams } from 'react-router';
import {
  AdminEmptyState,
  AdminPageHeader,
  AdminPageShell,
  AdminSection,
  Alert,
  AlertDescription,
  AlertTitle,
  Badge,
  Button,
  Skeleton,
} from '@shared/ui';

// START_CONTRACT: UserDetailPage
//   PURPOSE: Render one user loaded by route id with UI-kit loading, error, not-found, and detail states.
//   INPUTS: none.
//   OUTPUTS: { JSX.Element - user detail route for the current route id }
//   SIDE_EFFECTS: Sends GetUser GraphQL query through React Query.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
// END_CONTRACT: UserDetailPage
export default function UserDetailPage() {
  const { id } = useParams<{ id: string }>();
  const userQuery = useQuery({
    enabled: Boolean(id),
    queryKey: ['admin-user', id],
    queryFn: () => graphqlClient.request<GetUserQuery>(getUserQueryDocument, { id }),
  });

  const user = userQuery.data?.user || null;

  // START_BLOCK_USER_DETAIL_STATES
  if (userQuery.isLoading) {
    return (
      <AdminPageShell>
        <div aria-label="Loading user" aria-live="polite" className="space-y-4" role="status">
          <Skeleton className="h-9 w-64" />
          <Skeleton className="h-48 w-full" />
        </div>
      </AdminPageShell>
    );
  }

  if (userQuery.isError) {
    return (
      <AdminPageShell>
        <Alert variant="destructive">
          <AlertTitle>Failed to load user.</AlertTitle>
          <AlertDescription>Refresh the page after the GraphQL API is available.</AlertDescription>
        </Alert>
        <Button asChild variant="outline">
          <Link to="/users">Back to users</Link>
        </Button>
      </AdminPageShell>
    );
  }

  if (!user) {
    return (
      <AdminPageShell>
        <AdminEmptyState
          title="User not found"
          description="The requested user does not exist."
          action={
            <Button asChild variant="outline">
              <Link to="/users">Back to users</Link>
            </Button>
          }
        />
      </AdminPageShell>
    );
  }

  return (
    <AdminPageShell>
      <AdminPageHeader
        title={user.name}
        description="Reference user loaded through the admin GraphQL API."
        actions={
          <Button asChild variant="outline">
            <Link to="/users">Back to users</Link>
          </Button>
        }
      />

      <AdminSection title="Profile" description="Stable user fields from GraphQL.">
        <dl className="grid gap-4 text-sm sm:grid-cols-2">
          <div className="space-y-1">
            <dt className="font-medium text-muted-foreground">Email</dt>
            <dd>{user.email}</dd>
          </div>
          <div className="space-y-1">
            <dt className="font-medium text-muted-foreground">ID</dt>
            <dd className="break-all">{user.id}</dd>
          </div>
          <div className="space-y-1">
            <dt className="font-medium text-muted-foreground">Created</dt>
            <dd>{new Date(user.createdAt).toLocaleString()}</dd>
          </div>
          <div className="space-y-1">
            <dt className="font-medium text-muted-foreground">Updated</dt>
            <dd>{new Date(user.updatedAt).toLocaleString()}</dd>
          </div>
        </dl>
        <div className="mt-4">
          <Badge variant="secondary">GraphQL user</Badge>
        </div>
      </AdminSection>
    </AdminPageShell>
  );
  // END_BLOCK_USER_DETAIL_STATES
}
