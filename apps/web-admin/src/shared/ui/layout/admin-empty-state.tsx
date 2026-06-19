// FILE: apps/web-admin/src/shared/ui/layout/admin-empty-state.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide a reusable empty/not-found state for web-admin routes.
//   SCOPE: Owns empty-state copy layout and optional action slot; excludes data fetching and route decisions.
//   DEPENDS: react, apps/web-admin/src/shared/ui/primitives/card.tsx.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminEmptyState - Reusable admin empty-state panel.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added standard admin empty state.
// END_CHANGE_SUMMARY

import type { ReactNode } from 'react';
import { Card, CardContent } from '../primitives/card';

type AdminEmptyStateProps = {
  title: string;
  description: string;
  action?: ReactNode;
};

// START_CONTRACT: AdminEmptyState
//   PURPOSE: Render a reusable empty or not-found state for admin routes.
//   INPUTS: { props: AdminEmptyStateProps - title, description, and optional action }
//   OUTPUTS: { JSX.Element - centered empty-state card }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: AdminEmptyState
export function AdminEmptyState({ title, description, action }: AdminEmptyStateProps) {
  return (
    <Card>
      <CardContent className="flex min-h-44 flex-col items-center justify-center gap-3 text-center">
        <div className="space-y-1">
          <h2 className="text-lg font-medium">{title}</h2>
          <p className="text-sm text-muted-foreground">{description}</p>
        </div>
        {action ? <div className="pt-1">{action}</div> : null}
      </CardContent>
    </Card>
  );
}
