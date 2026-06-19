// FILE: apps/web/vitest.setup.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Register public web Vitest DOM matchers.
//   SCOPE: Loads jest-dom matchers for Vitest assertions; excludes test fixtures and mocks.
//   DEPENDS: @testing-library/jest-dom/vitest.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: TEST
//   MAP_MODE: NONE
// END_MODULE_CONTRACT
// START_MODULE_MAP
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Vitest DOM matcher setup for Next web tests.
// END_CHANGE_SUMMARY

import '@testing-library/jest-dom/vitest';
