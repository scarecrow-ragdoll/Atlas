// FILE: tools/workspace/dev-config.test.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify local development command and port contracts for the Nx workspace.
//   SCOPE: Covers root dev project selection and app-level frontend dev ports; excludes live server startup.
//   DEPENDS: package.json, apps/web-admin/package.json, apps/web/package.json.
//   LINKS: M-WORKSPACE / V-M-WORKSPACE / VF-LOCAL-DEV.
//   ROLE: TEST
//   MAP_MODE: LOCALS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   packageJson - Reads JSON package manifests relative to the workspace root.
//   projectJson - Reads Nx project manifests relative to the workspace root.
//   scriptPort - Extracts a numeric --port value from a package script.
//   dev config tests - Prove bun run dev starts the local web stack on distinct ports.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Added API local serve admin-auth default coverage.
// END_CHANGE_SUMMARY

import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import { describe, expect, it } from 'vitest';

type PackageJson = {
  scripts?: Record<string, string>;
};

type ProjectJson = {
  targets?: Record<string, { options?: { command?: string } }>;
};

const workspaceRoot = resolve(__dirname, '../..');

function packageJson(path: string): PackageJson {
  return JSON.parse(readFileSync(resolve(workspaceRoot, path), 'utf8')) as PackageJson;
}

function projectJson(path: string): ProjectJson {
  return JSON.parse(readFileSync(resolve(workspaceRoot, path), 'utf8')) as ProjectJson;
}

function scriptPort(script: string | undefined): number {
  const match = script?.match(/--port\s+(\d+)/);

  if (!match) {
    throw new Error(`script does not declare an explicit --port: ${script ?? '<missing>'}`);
  }

  return Number(match[1]);
}

describe('workspace local dev configuration', () => {
  it('runs the local web stack without requiring the bot service', () => {
    const root = packageJson('package.json');

    expect(root.scripts?.dev).toBe(
      'bunx nx run-many --target=serve --projects=api,web-admin,web --parallel=3',
    );
  });

  it('provides local API admin auth defaults for the serve target', () => {
    const api = projectJson('apps/api/project.json');
    const command = api.targets?.serve?.options?.command || '';

    expect(command).toContain('ADMIN_INITIAL_EMAIL="${ADMIN_INITIAL_EMAIL:-admin@example.com}"');
    expect(command).toContain(
      'ADMIN_INITIAL_PASSWORD="${ADMIN_INITIAL_PASSWORD:-ChangeMeAdmin123!}"',
    );
    expect(command).toContain('ADMIN_INITIAL_NAME="${ADMIN_INITIAL_NAME:-Template Admin}"');
    expect(command).toContain(
      'ADMIN_SESSION_KEY_SECRET="${ADMIN_SESSION_KEY_SECRET:-dev-session-key-secret}"',
    );
    expect(command).toContain('cd apps/api &&');
    expect(command).toContain('air -c air.toml');
  });

  it('assigns stable distinct ports to web-admin and web dev servers', () => {
    const webAdmin = packageJson('apps/web-admin/package.json');
    const web = packageJson('apps/web/package.json');

    expect(scriptPort(webAdmin.scripts?.dev)).toBe(3100);
    expect(scriptPort(webAdmin.scripts?.preview)).toBe(3100);
    expect(scriptPort(web.scripts?.dev)).toBe(3101);
    expect(scriptPort(web.scripts?.start)).toBe(3101);
    expect(scriptPort(webAdmin.scripts?.dev)).not.toBe(scriptPort(web.scripts?.dev));
  });
});
