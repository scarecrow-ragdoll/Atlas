// FILE: apps/web/src/shared/ui/primitives/ui-primitives.test.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify the public web UI-kit primitive exports.
//   SCOPE: Covers the independent apps/web @shared/ui barrel and basic accessible primitive rendering; excludes page composition behavior.
//   DEPENDS: apps/web/src/shared/ui, @testing-library/react, vitest.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   UI primitive tests - Prove public web shadcn-compatible primitives are exported independently from web-admin.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added red coverage for the public web UI-kit surface.
// END_CHANGE_SUMMARY

import { cleanup, render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it } from 'vitest';
import {
  Alert,
  AlertDescription,
  AlertTitle,
  Badge,
  Button,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Input,
  Label,
} from '@shared/ui';

afterEach(() => {
  cleanup();
});

describe('public web UI primitives', () => {
  it('exports shadcn-compatible primitives through the local @shared/ui barrel', () => {
    render(
      <Card>
        <CardHeader>
          <Badge>Public web</Badge>
          <CardTitle>REST users</CardTitle>
          <CardDescription>Independent UI kit</CardDescription>
        </CardHeader>
        <CardContent>
          <Alert>
            <AlertTitle>Ready</AlertTitle>
            <AlertDescription>Uses apps/web UI primitives.</AlertDescription>
          </Alert>
          <Label htmlFor="public-email">Email</Label>
          <Input id="public-email" placeholder="Email" />
          <Button type="button">Create</Button>
        </CardContent>
      </Card>,
    );

    expect(screen.getByText('REST users')).toBeInTheDocument();
    expect(screen.getByText('Independent UI kit')).toBeInTheDocument();
    expect(screen.getByText('Uses apps/web UI primitives.')).toBeInTheDocument();
    expect(screen.getByLabelText('Email')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Create' })).toHaveAttribute('data-slot', 'button');
  });
});
