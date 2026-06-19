// FILE: apps/web-admin/src/app/config.ts
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Expose app-layer access to shared web-admin configuration.
//   SCOPE: Re-exports shared configuration for app imports; excludes config resolution logic.
//   DEPENDS: apps/web-admin/src/shared/config/index.ts.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: BARREL
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   appConfig - Re-exported shared configuration object.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Documented Vite env config re-export.
// END_CHANGE_SUMMARY

export { appConfig } from '@shared/config';
