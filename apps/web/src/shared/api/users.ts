// FILE: apps/web/src/shared/api/users.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide the public web REST users client.
//   SCOPE: Maps REST envelopes, sends CRUD requests, encodes user ids, and normalizes API errors; excludes Next route proxy forwarding and UI state.
//   DEPENDS: apps/web/src/shared/config.ts, global fetch.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   User - Public REST user shape.
//   CreateUserInput - Public REST create-user request shape.
//   UpdateUserInput - Public REST update-user request shape.
//   ApiError - Structured REST client error.
//   listUsers - Load users from /api/users.
//   createUser - Create one user through /api/users.
//   getUser - Load one user through an encoded /api/users/:id path.
//   updateUser - Update one user through an encoded /api/users/:id path.
//   deleteUser - Delete one user through an encoded /api/users/:id path.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added GRACE contract while moving browser calls behind the same-origin Next proxy.
// END_CHANGE_SUMMARY

import { appConfig } from '../config';

export type User = {
  id: string;
  email: string;
  name: string;
  createdAt: string;
  updatedAt: string;
};

export type CreateUserInput = {
  email: string;
  name: string;
  password: string;
};

export type UpdateUserInput = {
  email?: string;
  name?: string;
};

type ApiEnvelope<T> = { data: T };
type ApiListEnvelope<T> = { data: T[]; meta: { totalCount: number } };
type ApiErrorEnvelope = { error: { code: string; message: string; field?: string } };

export class ApiError extends Error {
  readonly code: string;
  readonly field?: string;
  readonly status: number;

  constructor(error: ApiErrorEnvelope['error'], status: number) {
    super(error.message);
    this.name = 'ApiError';
    this.code = error.code;
    this.field = error.field;
    this.status = status;
  }
}

async function request<T>(path: string, init: RequestInit): Promise<T> {
  const response = await fetch(`${appConfig.apiBaseUrl}${path}`, {
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      ...init.headers,
    },
    ...init,
  });

  if (response.status === 204) {
    return undefined as T;
  }

  const payload = (await response.json()) as ApiEnvelope<T> | ApiListEnvelope<T> | ApiErrorEnvelope;

  if (!response.ok || 'error' in payload) {
    const error =
      'error' in payload
        ? payload.error
        : { code: 'HTTP_ERROR', message: `Request failed with status ${response.status}` };
    throw new ApiError(error, response.status);
  }

  return payload as T;
}

export async function listUsers(): Promise<{ users: User[]; totalCount: number }> {
  const envelope = await request<ApiListEnvelope<User>>('/api/users', { method: 'GET' });
  return { users: envelope.data, totalCount: envelope.meta.totalCount };
}

export async function createUser(input: CreateUserInput): Promise<User> {
  const envelope = await request<ApiEnvelope<User>>('/api/users', {
    method: 'POST',
    body: JSON.stringify(input),
  });
  return envelope.data;
}

export async function getUser(id: string): Promise<User> {
  const envelope = await request<ApiEnvelope<User>>(`/api/users/${encodeURIComponent(id)}`, {
    method: 'GET',
  });
  return envelope.data;
}

export async function updateUser(id: string, input: UpdateUserInput): Promise<User> {
  const envelope = await request<ApiEnvelope<User>>(`/api/users/${encodeURIComponent(id)}`, {
    method: 'PATCH',
    body: JSON.stringify(input),
  });
  return envelope.data;
}

export async function deleteUser(id: string): Promise<void> {
  await request<void>(`/api/users/${encodeURIComponent(id)}`, { method: 'DELETE' });
}
