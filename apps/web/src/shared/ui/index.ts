// FILE: apps/web/src/shared/ui/index.ts
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Expose the independent public web UI-kit surface.
//   SCOPE: Re-exports apps/web shadcn-compatible primitives and public theme toggle; excludes web-admin UI files and implementation-only utility subpaths.
//   DEPENDS: apps/web/src/shared/ui/primitives, apps/web/src/shared/ui/theme-toggle.tsx.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: BARREL
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Alert - Public web alert primitive export.
//   AlertDescription - Public web alert body primitive export.
//   AlertTitle - Public web alert title primitive export.
//   Badge - Public web badge primitive export.
//   Button - Public web button primitive export.
//   Card - Public web card primitive export.
//   CardContent - Public web card content primitive export.
//   CardDescription - Public web card description primitive export.
//   CardHeader - Public web card header primitive export.
//   CardTitle - Public web card title primitive export.
//   Input - Public web input primitive export.
//   Label - Public web label primitive export.
//   ThemeToggle - Public web theme toggle export.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Exported the public web theme toggle through the local UI-kit barrel.
// END_CHANGE_SUMMARY

export { Alert, AlertDescription, AlertTitle } from './primitives/alert';
export { Badge } from './primitives/badge';
export { Button } from './primitives/button';
export { Card, CardContent, CardDescription, CardHeader, CardTitle } from './primitives/card';
export { Input } from './primitives/input';
export { Label } from './primitives/label';
export { ThemeToggle } from './theme-toggle';
