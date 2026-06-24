// FILE: apps/web-admin/src/pages/atlas/daily-nutrition-override-page.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Retire the legacy Daily Nutrition Override route in favor of the factual daily food log.
//   SCOPE: Redirects old override URLs to /atlas/nutrition; excludes rendering legacy target/override forms or mock reference UI.
//   DEPENDS: react-router.
//   LINKS: M-WEB-ADMIN / M-API-NUTRITION / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - Legacy override route redirect to the factual Nutrition route.
// END_MODULE_MAP

import { Navigate } from 'react-router';

// START_CONTRACT: DailyNutritionOverridePage
//   PURPOSE: Redirect legacy override URLs to the factual daily food log route.
//   INPUTS: none.
//   OUTPUTS: { JSX.Element - react-router redirect to /atlas/nutrition }
//   SIDE_EFFECTS: Updates browser navigation through react-router.
//   LINKS: M-WEB-ADMIN / M-API-NUTRITION / V-M-WEB-ADMIN.
// END_CONTRACT: DailyNutritionOverridePage
export default function DailyNutritionOverridePage() {
  return <Navigate replace to="/atlas/nutrition" />;
}
