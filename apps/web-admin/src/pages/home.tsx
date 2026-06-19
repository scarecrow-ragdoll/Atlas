// FILE: apps/web-admin/src/pages/home.tsx
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Render the web-admin home route.
//   SCOPE: Shows admin entry cards for users and UI-kit reference routes; excludes data fetching and mutation behavior.
//   DEPENDS: react-router, apps/web-admin/src/shared/ui.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - Admin home route with users and UI-kit navigation cards.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Migrated home route to the web-admin UI kit and added UI-kit navigation.
// END_CHANGE_SUMMARY

import { Link } from 'react-router';
import {
  AdminPageHeader,
  AdminPageShell,
  Button,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@shared/ui';

// START_CONTRACT: HomePage
//   PURPOSE: Render admin route entry cards for users and the UI-kit reference page.
//   INPUTS: none.
//   OUTPUTS: { JSX.Element - home route with admin navigation cards }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: HomePage
export default function HomePage() {
  return (
    <AdminPageShell>
      <AdminPageHeader
        title="Monorepo Template Admin"
        description="GraphQL admin client and UI reference for new admin pages."
      />

      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Users</CardTitle>
            <CardDescription>Manage the reference GraphQL users flow.</CardDescription>
          </CardHeader>
          <CardContent>
            <Button asChild>
              <Link to="/users">Open users</Link>
            </Button>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>UI Kit</CardTitle>
            <CardDescription>Review the approved components for admin pages.</CardDescription>
          </CardHeader>
          <CardContent>
            <Button asChild variant="outline">
              <Link to="/ui-kit">Open UI kit</Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    </AdminPageShell>
  );
}
