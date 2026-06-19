// FILE: apps/web-admin/src/pages/users-page.tsx
// VERSION: 1.0.2
// START_MODULE_CONTRACT
//   PURPOSE: Render the web-admin users list and create-user route.
//   SCOPE: Loads users, displays list states, submits create-user GraphQL mutations, and links to details through UI-kit components; excludes detail rendering and GraphQL transport construction.
//   DEPENDS: @tanstack/react-query, react, react-router, apps/web-admin/src/entities/user/api/*.graphql, apps/web-admin/src/shared/api/graphql-client.ts, generated GraphQL types, apps/web-admin/src/shared/ui.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - Users list and create-form route backed by codegen-visible GraphQL documents and UI-kit components.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.2 - Removed redundant global navigation links after sidebar shell migration.
// END_CHANGE_SUMMARY

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import createUserMutationDocument from '@entities/user/api/createUser.graphql?raw';
import getUsersQueryDocument from '@entities/user/api/users.graphql?raw';
import { graphqlClient } from '@shared/api/graphql-client';
import type { CreateUserMutation, GetUsersQuery } from '@shared/api/generated/types';
import { type FormEvent, useState } from 'react';
import { Link } from 'react-router';
import {
  AdminEmptyState,
  AdminPageHeader,
  AdminPageShell,
  AdminSection,
  AdminToolbar,
  Alert,
  AlertDescription,
  AlertTitle,
  Button,
  Input,
  Label,
  Skeleton,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@shared/ui';

type FormState = {
  name: string;
  email: string;
  password: string;
};

const initialFormState: FormState = { name: '', email: '', password: '' };

// START_CONTRACT: errorMessageFromUnknown
//   PURPOSE: Convert unknown mutation errors into user-visible fallback copy.
//   INPUTS: { error: unknown - mutation error from React Query }
//   OUTPUTS: { string - safe error message }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: errorMessageFromUnknown
function errorMessageFromUnknown(error: unknown): string {
  return error instanceof Error ? error.message : 'Request failed';
}

// START_CONTRACT: UsersPage
//   PURPOSE: Render the users list and create-user GraphQL mutation flow through UI-kit components.
//   INPUTS: none.
//   OUTPUTS: { JSX.Element - users route with loading, error, empty, list, and create states }
//   SIDE_EFFECTS: Sends createUser mutation and invalidates admin-users query on successful user creation.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
// END_CONTRACT: UsersPage
export default function UsersPage() {
  const queryClient = useQueryClient();
  const [form, setForm] = useState<FormState>(initialFormState);
  const [error, setError] = useState<string | null>(null);

  const usersQuery = useQuery({
    queryKey: ['admin-users'],
    queryFn: () => graphqlClient.request<GetUsersQuery>(getUsersQueryDocument, { first: 20 }),
  });

  const mutation = useMutation({
    mutationFn: (input: FormState) =>
      graphqlClient.request<CreateUserMutation>(createUserMutationDocument, { input }),
    onError: (mutationError) => setError(errorMessageFromUnknown(mutationError)),
    onSuccess: async (response) => {
      const result = response.createUser;

      // START_BLOCK_CREATE_USER_RESULT
      if ('user' in result) {
        setForm(initialFormState);
        setError(null);
        await queryClient.invalidateQueries({ queryKey: ['admin-users'] });
        return;
      }

      if ('field' in result) {
        setError(`${result.field}: ${result.message}`);
        return;
      }

      setError(result.message);
      // END_BLOCK_CREATE_USER_RESULT
    },
  });

  // START_CONTRACT: updateField
  //   PURPOSE: Update one create-user form field without mutating the other fields.
  //   INPUTS: { field: keyof FormState - field to change, value: string - next field value }
  //   OUTPUTS: none.
  //   SIDE_EFFECTS: Updates local React state.
  //   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
  // END_CONTRACT: updateField
  function updateField(field: keyof FormState, value: string) {
    setForm((current) => ({ ...current, [field]: value }));
  }

  // START_CONTRACT: handleSubmit
  //   PURPOSE: Submit the current create-user form through the GraphQL mutation.
  //   INPUTS: { event: FormEvent<HTMLFormElement> - form submit event }
  //   OUTPUTS: none.
  //   SIDE_EFFECTS: Prevents default form submission and starts the createUser mutation.
  //   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN / M-GRAPHQL-SCHEMA.
  // END_CONTRACT: handleSubmit
  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    mutation.mutate(form);
  }

  const users = usersQuery.data?.users.edges || [];

  return (
    <AdminPageShell>
      <AdminPageHeader
        title="Users"
        description="Create and inspect reference users through the admin GraphQL API."
      />

      <AdminToolbar>
        <p className="text-sm text-muted-foreground">Showing the latest 20 users.</p>
      </AdminToolbar>

      <AdminSection title="Create user" description="Submit a GraphQL createUser mutation.">
        <form className="grid gap-4 md:grid-cols-[1fr_1fr_1fr_auto]" onSubmit={handleSubmit}>
          <div className="space-y-2">
            <Label htmlFor="user-name">Name</Label>
            <Input
              id="user-name"
              onChange={(event) => updateField('name', event.target.value)}
              placeholder="Name"
              value={form.name}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="user-email">Email</Label>
            <Input
              id="user-email"
              onChange={(event) => updateField('email', event.target.value)}
              placeholder="Email"
              type="email"
              value={form.email}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="user-password">Password</Label>
            <Input
              id="user-password"
              onChange={(event) => updateField('password', event.target.value)}
              placeholder="Password"
              type="password"
              value={form.password}
            />
          </div>
          <div className="flex items-end">
            <Button disabled={mutation.isPending} type="submit">
              {mutation.isPending ? 'Creating...' : 'Create'}
            </Button>
          </div>
        </form>
      </AdminSection>

      {error ? (
        <Alert variant="destructive">
          <AlertTitle>Create user failed</AlertTitle>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      ) : null}

      {usersQuery.isError ? (
        <Alert variant="destructive">
          <AlertTitle>Failed to load users.</AlertTitle>
          <AlertDescription>Refresh the page after the GraphQL API is available.</AlertDescription>
        </Alert>
      ) : null}

      {/* START_BLOCK_USERS_LIST_STATES */}
      <AdminSection
        title="Directory"
        description={
          usersQuery.data ? `Total: ${usersQuery.data.users.totalCount}` : 'Loading users.'
        }
      >
        {usersQuery.isLoading ? (
          <div className="space-y-2">
            <Skeleton className="h-10 w-full" />
            <Skeleton className="h-10 w-full" />
          </div>
        ) : null}

        {usersQuery.data && users.length === 0 ? (
          <AdminEmptyState title="No users yet" description="No users yet. Create one above." />
        ) : null}

        {users.length > 0 ? (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Email</TableHead>
                <TableHead className="text-right">Details</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {users.map(({ node }) => (
                <TableRow key={node.id}>
                  <TableCell>
                    <Button asChild className="h-auto p-0" variant="link">
                      <Link to={`/users/${node.id}`}>{node.name}</Link>
                    </Button>
                  </TableCell>
                  <TableCell>{node.email}</TableCell>
                  <TableCell className="text-right">
                    <Button asChild size="sm" variant="outline">
                      <Link to={`/users/${node.id}`}>Open</Link>
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        ) : null}
      </AdminSection>
      {/* END_BLOCK_USERS_LIST_STATES */}
    </AdminPageShell>
  );
}
