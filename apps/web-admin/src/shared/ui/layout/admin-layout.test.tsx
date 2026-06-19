// FILE: apps/web-admin/src/shared/ui/layout/admin-layout.test.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify web-admin layout composition behavior.
//   SCOPE: Covers page shell, header actions, toolbar layout, sections, and empty states; excludes page data behavior.
//   DEPENDS: apps/web-admin/src/shared/ui/layout, @testing-library/react, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   admin layout tests - Prove shared admin compositions render accessible structure and optional actions.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Added coverage for optional header and section branches.
// END_CHANGE_SUMMARY

import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import {
  AdminEmptyState,
  AdminPageHeader,
  AdminPageShell,
  AdminSection,
  AdminToolbar,
  Button,
} from '@shared/ui';

describe('admin layout compositions', () => {
  it('renders a page shell with header actions and section content', () => {
    render(
      <AdminPageShell>
        <AdminPageHeader
          title="Users"
          description="Manage reference users."
          actions={<Button>New user</Button>}
        />
        <AdminToolbar>
          <Button variant="outline">Refresh</Button>
        </AdminToolbar>
        <AdminSection title="Directory" description="Current users in the system.">
          <p>One User</p>
        </AdminSection>
      </AdminPageShell>,
    );

    expect(screen.queryByRole('main')).not.toBeInTheDocument();
    expect(screen.getByTestId('admin-page-shell')).toHaveClass('max-w-6xl');
    expect(screen.getByRole('heading', { name: 'Users' })).toBeInTheDocument();
    expect(screen.getByText('Manage reference users.')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'New user' })).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: /Switch to .* theme/ })).not.toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Refresh' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Directory' })).toBeInTheDocument();
    expect(screen.getByText('One User')).toBeInTheDocument();
  });

  it('renders empty states with optional action', () => {
    render(
      <AdminEmptyState
        title="No users yet"
        description="Create the first reference user."
        action={<Button>Create user</Button>}
      />,
    );

    expect(screen.getByRole('heading', { name: 'No users yet' })).toBeInTheDocument();
    expect(screen.getByText('Create the first reference user.')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Create user' })).toBeInTheDocument();
  });

  it('renders header and sections without optional slots', () => {
    render(
      <AdminPageShell>
        <AdminPageHeader title="Settings" />
        <AdminSection title="General">
          <p>General settings</p>
        </AdminSection>
      </AdminPageShell>,
    );

    expect(screen.getByRole('heading', { name: 'Settings' })).toBeInTheDocument();
    expect(screen.queryByText('Manage reference users.')).not.toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'General' })).toBeInTheDocument();
    expect(screen.getByText('General settings')).toBeInTheDocument();
  });
});
