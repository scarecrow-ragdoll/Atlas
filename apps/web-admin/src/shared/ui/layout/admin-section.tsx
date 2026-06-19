// FILE: apps/web-admin/src/shared/ui/layout/admin-section.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide a titled card-backed content section for web-admin routes.
//   SCOPE: Owns section heading and card framing; excludes page-level layout and business data behavior.
//   DEPENDS: react, apps/web-admin/src/shared/ui/primitives/card.tsx.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AdminSection - Card-backed section with optional description.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added standard admin content section.
// END_CHANGE_SUMMARY

import type { ReactNode } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../primitives/card';

type AdminSectionProps = {
  title: string;
  description?: string;
  children: ReactNode;
};

// START_CONTRACT: AdminSection
//   PURPOSE: Render a titled admin content section using the approved card primitive.
//   INPUTS: { props: AdminSectionProps - title, optional description, and section children }
//   OUTPUTS: { JSX.Element - card-backed content section }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: AdminSection
export function AdminSection({ title, description, children }: AdminSectionProps) {
  return (
    <section>
      <Card>
        <CardHeader>
          <CardTitle>
            <h2>{title}</h2>
          </CardTitle>
          {description ? <CardDescription>{description}</CardDescription> : null}
        </CardHeader>
        <CardContent>{children}</CardContent>
      </Card>
    </section>
  );
}
