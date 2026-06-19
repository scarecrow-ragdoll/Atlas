// FILE: tools/coverage/preflight.mjs
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Validate that coverage gate prerequisites exist before expensive coverage runs.
//   SCOPE: Checks required repository files and root scripts for coverage, codegen, e2e, and verification gates; excludes running tests or interpreting coverage summaries.
//   DEPENDS: package.json, tools/coverage/coverage.config.json, apps/web-admin, apps/web, docker, docs.
//   LINKS: M-COVERAGE-GATE / V-M-COVERAGE-GATE.
//   ROLE: SCRIPT
//   MAP_MODE: LOCALS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   requiredFiles - Repository contract files that must exist for the coverage gate.
//   requiredScripts - Root scripts required by coverage and release handoff.
//   fail - Records a preflight failure without stopping subsequent checks.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Updated required files for Vite admin and Next public web swap.
// END_CHANGE_SUMMARY

import fs from 'node:fs';
import path from 'node:path';

const root = process.cwd();
const requiredFiles = [
  'tools/coverage/coverage.config.json',
  'package.json',
  'apps/web-admin/package.json',
  'apps/web-admin/project.json',
  'apps/web-admin/vite.config.ts',
  'apps/web-admin/tsconfig.json',
  'apps/web-admin/index.html',
  'apps/web-admin/src/App.tsx',
  'apps/web-admin/src/main.tsx',
  'apps/web-admin/src/entities/user/api/users.graphql',
  'apps/web-admin/src/entities/user/api/createUser.graphql',
  'apps/web-admin/src/entities/user/api/user.graphql',
  'apps/web-admin/src/shared/api/graphql-client.ts',
  'apps/web-admin/src/shared/config/index.ts',
  'apps/web-admin/e2e/playwright.config.ts',
  'apps/web-admin/e2e/preflight.mjs',
  'apps/web/package.json',
  'apps/web/project.json',
  'apps/web/next.config.js',
  'apps/web/tsconfig.json',
  'apps/web/vitest.config.ts',
  'apps/web/vitest.setup.ts',
  'apps/web/app/layout.tsx',
  'apps/web/app/page.tsx',
  'apps/web/app/users-client.tsx',
  'apps/web/app/api/users/route.ts',
  'apps/web/app/api/users/[id]/route.ts',
  'apps/web/src/shared/api/users.ts',
  'apps/web/src/shared/config.ts',
  'apps/web/e2e/playwright.config.ts',
  'docker/docker-compose.yml',
  'docker/docker-compose.test.yml',
  'docs/verification-plan.xml',
  '.tasks/swap-web-admin-vite-web-next/verification.md',
];

const requiredScripts = ['test:coverage', 'test:e2e', 'verify:coverage'];

function fail(message) {
  console.error(`[Coverage][preflight] ${message}`);
  process.exitCode = 1;
}

for (const file of requiredFiles) {
  if (!fs.existsSync(path.join(root, file))) {
    fail(`Missing required file: ${file}`);
  }
}

const pkg = JSON.parse(fs.readFileSync(path.join(root, 'package.json'), 'utf8'));
for (const script of requiredScripts) {
  if (!pkg.scripts?.[script]) {
    fail(`Missing package.json script: ${script}`);
  }
}

if (process.exitCode) {
  process.exit(process.exitCode);
}

console.log('[Coverage][preflight] ok');
