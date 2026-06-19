'use client';

// FILE: apps/web/app/users-client.tsx
// VERSION: 1.0.2
// START_MODULE_CONTRACT
//   PURPOSE: Provide the interactive public REST users experience inside the Next app.
//   SCOPE: Renders initial users, creates users, refetches the list, displays REST errors, and shows selected-user details; excludes server data fetching and route proxy internals.
//   DEPENDS: @tanstack/react-query, react, apps/web/src/shared/api/users.ts, apps/web/src/shared/ui.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - Client component for public REST users list, create form, and detail panel.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.3 - Added the persisted public web theme toggle.
// END_CHANGE_SUMMARY

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { type FormEvent, useState } from 'react';
import {
  Alert,
  AlertDescription,
  Badge,
  Button,
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  Input,
  Label,
  ThemeToggle,
} from '@shared/ui';
import { ApiError, createUser, listUsers, type User } from '../src/shared/api/users';

type FormState = {
  name: string;
  email: string;
  password: string;
};

type UsersClientProps = {
  initialLoadError?: boolean;
  initialUsers: User[];
  initialTotalCount: number;
};

const initialFormState: FormState = {
  email: '',
  name: '',
  password: '',
};

function formatCreateError(error: Error): string {
  if (error instanceof ApiError && error.field) {
    return `${error.field}: ${error.message}`;
  }

  return error.message;
}

export default function UsersClient({
  initialLoadError = false,
  initialTotalCount,
  initialUsers,
}: UsersClientProps) {
  const queryClient = useQueryClient();
  const [form, setForm] = useState<FormState>(initialFormState);
  const [selectedUserId, setSelectedUserId] = useState<string | null>(null);
  const [createError, setCreateError] = useState<string | null>(null);

  const usersQuery = useQuery({
    initialData: { totalCount: initialTotalCount, users: initialUsers },
    queryFn: listUsers,
    queryKey: ['users'],
  });

  const createUserMutation = useMutation({
    mutationFn: createUser,
    onError: (error: Error) => setCreateError(formatCreateError(error)),
    onMutate: () => setCreateError(null),
    onSuccess: async () => {
      setForm(initialFormState);
      await queryClient.invalidateQueries({ queryKey: ['users'] });
    },
  });

  const users = usersQuery.data.users;
  const hasLoadError = initialLoadError || usersQuery.isError;
  const selectedUser = users.find((user) => user.id === selectedUserId) ?? null;

  function updateField(field: keyof FormState, value: string) {
    setForm((current) => ({ ...current, [field]: value }));
  }

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    createUserMutation.mutate(form);
  }

  return (
    <main className="app-shell">
      <div className="app-actions">
        <ThemeToggle />
      </div>

      <Card className="users-panel">
        <CardHeader className="panel-heading">
          <CardTitle>REST Web</CardTitle>
          <Badge>{usersQuery.data.totalCount} users</Badge>
        </CardHeader>

        <CardContent>
          <form className="create-form" onSubmit={handleSubmit}>
            <div className="field-control">
              <Label className="sr-only" htmlFor="user-name">
                Name
              </Label>
              <Input
                id="user-name"
                onChange={(event) => updateField('name', event.target.value)}
                placeholder="Name"
                value={form.name}
              />
            </div>
            <div className="field-control">
              <Label className="sr-only" htmlFor="user-email">
                Email
              </Label>
              <Input
                id="user-email"
                onChange={(event) => updateField('email', event.target.value)}
                placeholder="Email"
                type="email"
                value={form.email}
              />
            </div>
            <div className="field-control">
              <Label className="sr-only" htmlFor="user-password">
                Password
              </Label>
              <Input
                id="user-password"
                onChange={(event) => updateField('password', event.target.value)}
                placeholder="Password"
                type="password"
                value={form.password}
              />
            </div>
            <Button disabled={createUserMutation.isPending} type="submit">
              Create
            </Button>
          </form>

          {createError ? (
            <Alert className="flow-message" variant="destructive">
              <AlertDescription>{createError}</AlertDescription>
            </Alert>
          ) : null}
          {hasLoadError ? (
            <Alert className="flow-message" variant="destructive">
              <AlertDescription>Failed to load users.</AlertDescription>
            </Alert>
          ) : null}

          {!usersQuery.isLoading && !hasLoadError && users.length === 0 ? (
            <p className="empty-state">No users yet.</p>
          ) : null}

          <div aria-label="Users" className="user-list">
            {users.map((user) => (
              <div className="user-row" key={user.id}>
                <Button
                  className="user-row-button"
                  onClick={() => setSelectedUserId(user.id)}
                  type="button"
                  variant="ghost"
                >
                  {user.name}
                </Button>
                <span>{user.email}</span>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {selectedUser ? (
        <Card aria-label="Selected user" className="detail-panel">
          <CardHeader>
            <CardTitle>{selectedUser.name}</CardTitle>
          </CardHeader>
          <CardContent>
            <p>{selectedUser.email}</p>
          </CardContent>
        </Card>
      ) : null}
    </main>
  );
}
