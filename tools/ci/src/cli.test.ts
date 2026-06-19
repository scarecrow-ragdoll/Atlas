// FILE: tools/ci/src/cli.test.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verifies CI CLI command behavior, dependency injection, and deployment env requirements.
//   SCOPE: Covers command routing and injected side effects; excludes real Git, filesystem, and Dokploy network execution.
//   DEPENDS: vitest, node:fs, node:os, node:path, tools/ci/src/cli.ts.
//   LINKS: M-CI-CD / V-M-CI-CD.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   runCli - Guards affected-base, release tag, metadata, and Dokploy deploy command behavior.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added coverage for WEB_API_BASE_URL deploy requirements.
// END_CHANGE_SUMMARY

import { afterEach, describe, expect, it, vi } from 'vitest';
import { mkdtempSync, readFileSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';

const spawnSyncMock = vi.hoisted(() => vi.fn());

vi.mock('node:child_process', () => ({
  spawnSync: spawnSyncMock,
}));

import { runCli } from './cli';

describe('runCli', () => {
  afterEach(() => {
    spawnSyncMock.mockReset();
    vi.restoreAllMocks();
    vi.unstubAllGlobals();
  });

  it('prints affected base from injected environment', async () => {
    const output: string[] = [];
    const code = await runCli(['affected-base'], {
      env: { CI_MERGE_REQUEST_DIFF_BASE_SHA: 'abc123' },
      write: (line) => output.push(line),
    });

    expect(code).toBe(0);
    expect(output).toEqual(['abc123']);
  });

  it('prints affected base through default process and console dependencies', async () => {
    const log = vi.spyOn(console, 'log').mockImplementation(() => undefined);

    const code = await runCli(['affected-base']);

    expect(code).toBe(0);
    expect(log).toHaveBeenCalledWith(expect.stringMatching(/^origin\//));
  });

  it('rejects non-SemVer release tags before running git', async () => {
    const errors: string[] = [];
    const gitCalls: string[][] = [];
    const code = await runCli(['assert-release-tag'], {
      env: { CI_COMMIT_TAG: 'latest' },
      runGit: (args) => {
        gitCalls.push(args);
        return 0;
      },
      writeError: (line) => errors.push(line),
    });

    expect(code).toBe(1);
    expect(errors[0]).toContain('Release tag must match vX.Y.Z');
    expect(gitCalls).toEqual([]);
  });

  it('asserts SemVer release tags are reachable from origin main', async () => {
    const gitCalls: string[][] = [];
    const code = await runCli(['assert-release-tag'], {
      env: { CI_COMMIT_TAG: 'v1.2.3' },
      runGit: (args) => {
        gitCalls.push(args);
        return 0;
      },
    });

    expect(code).toBe(0);
    expect(gitCalls).toEqual([
      ['fetch', 'origin', 'main:refs/remotes/origin/main', '--tags'],
      ['merge-base', '--is-ancestor', 'v1.2.3^{commit}', 'origin/main'],
    ]);
  });

  it('fails release tag assertion when git fetch fails', async () => {
    const errors: string[] = [];
    const code = await runCli(['assert-release-tag'], {
      env: { CI_COMMIT_TAG: 'v1.2.3' },
      runGit: () => 1,
      writeError: (line) => errors.push(line),
    });

    expect(code).toBe(1);
    expect(errors[0]).toContain('Failed to fetch origin/main and tags');
  });

  it('fails release tag assertion when ancestry check fails', async () => {
    const errors: string[] = [];
    let calls = 0;
    const code = await runCli(['assert-release-tag'], {
      env: { CI_COMMIT_TAG: 'v1.2.3' },
      runGit: () => {
        calls += 1;
        return calls === 1 ? 0 : 1;
      },
      writeError: (line) => errors.push(line),
    });

    expect(code).toBe(1);
    expect(errors[0]).toContain('Release tag v1.2.3 is not reachable from origin/main');
  });

  it('runs git through the default runner with spawnSync', async () => {
    spawnSyncMock.mockReturnValue({ status: 0 });

    const code = await runCli(['assert-release-tag'], {
      env: { CI_COMMIT_TAG: 'v1.2.3' },
    });

    expect(code).toBe(0);
    expect(spawnSyncMock).toHaveBeenNthCalledWith(
      1,
      'git',
      ['fetch', 'origin', 'main:refs/remotes/origin/main', '--tags'],
      { stdio: 'inherit' },
    );
    expect(spawnSyncMock).toHaveBeenNthCalledWith(
      2,
      'git',
      ['merge-base', '--is-ancestor', 'v1.2.3^{commit}', 'origin/main'],
      { stdio: 'inherit' },
    );
  });

  it('treats missing default git process status as a failure', async () => {
    const errors: string[] = [];
    spawnSyncMock.mockReturnValue({});

    const code = await runCli(['assert-release-tag'], {
      env: { CI_COMMIT_TAG: 'v1.2.3' },
      writeError: (line) => errors.push(line),
    });

    expect(code).toBe(1);
    expect(errors[0]).toContain('Failed to fetch origin/main and tags');
  });

  it('writes image metadata with injected filesystem dependencies', async () => {
    const createdDirs: string[] = [];
    const files = new Map<string, string>();
    const output: string[] = [];

    const code = await runCli(['write-image-metadata'], {
      env: {
        CI_REGISTRY_IMAGE: 'registry.example.com/group/app',
        IMAGE_TAG: 'v1.2.3',
        IMAGE_METADATA_PATH: 'dist/custom/metadata.json',
        CI_COMMIT_SHA: 'abc123',
        CI_PIPELINE_ID: '42',
        API_IMAGE_DIGEST: 'sha256:api',
        WEB_IMAGE_DIGEST: 'sha256:web',
        BOT_IMAGE_DIGEST: 'sha256:bot',
      },
      makeDir: (path) => createdDirs.push(path),
      writeFile: (path, contents) => files.set(path, contents),
      write: (line) => output.push(line),
    });

    expect(code).toBe(0);
    expect(createdDirs).toEqual(['dist/custom']);
    expect(output).toEqual(['dist/custom/metadata.json']);
    expect(JSON.parse(files.get('dist/custom/metadata.json') ?? '{}')).toMatchObject({
      imageTag: 'v1.2.3',
      releaseTag: 'v1.2.3',
      services: [
        {
          service: 'api',
          image: 'registry.example.com/group/app/api:v1.2.3',
          digest: 'sha256:api',
        },
        {
          service: 'web',
          image: 'registry.example.com/group/app/web:v1.2.3',
          digest: 'sha256:web',
        },
        {
          service: 'bot',
          image: 'registry.example.com/group/app/bot:v1.2.3',
          digest: 'sha256:bot',
        },
      ],
    });
  });

  it('writes default-path metadata for non-release images', async () => {
    const files = new Map<string, string>();
    const code = await runCli(['write-image-metadata'], {
      env: {
        CI_REGISTRY_IMAGE: 'registry.example.com/group/app',
        IMAGE_TAG: 'develop-abc123',
        CI_COMMIT_SHA: 'abc123',
        CI_PIPELINE_ID: '42',
        API_IMAGE_DIGEST: 'sha256:api',
        WEB_IMAGE_DIGEST: 'sha256:web',
        BOT_IMAGE_DIGEST: 'sha256:bot',
      },
      makeDir: () => undefined,
      writeFile: (path, contents) => files.set(path, contents),
      write: () => undefined,
    });

    expect(code).toBe(0);
    expect(JSON.parse(files.get('dist/ci/image-metadata.json') ?? '{}').releaseTag).toBeNull();
  });

  it('writes image metadata through the default filesystem dependencies', async () => {
    const tempDir = mkdtempSync(join(tmpdir(), 'mt-ci-metadata-'));
    const outputPath = join(tempDir, 'nested', 'metadata.json');

    try {
      const code = await runCli(['write-image-metadata'], {
        env: {
          CI_REGISTRY_IMAGE: 'registry.example.com/group/app',
          IMAGE_TAG: 'develop-abc123',
          IMAGE_METADATA_PATH: outputPath,
          CI_COMMIT_SHA: 'abc123',
          CI_PIPELINE_ID: '42',
          API_IMAGE_DIGEST: 'sha256:api',
          WEB_IMAGE_DIGEST: 'sha256:web',
          BOT_IMAGE_DIGEST: 'sha256:bot',
        },
        write: () => undefined,
      });

      expect(code).toBe(0);
      expect(JSON.parse(readFileSync(outputPath, 'utf8'))).toMatchObject({
        imageTag: 'develop-abc123',
        releaseTag: null,
      });
    } finally {
      rmSync(tempDir, { recursive: true, force: true });
    }
  });

  it('writes metadata through the default filesystem dependencies', async () => {
    const root = mkdtempSync(join(tmpdir(), 'ci-tools-'));
    const outputPath = join(root, 'nested', 'metadata.json');

    try {
      const code = await runCli(['write-image-metadata'], {
        env: {
          CI_REGISTRY_IMAGE: 'registry.example.com/group/app',
          IMAGE_TAG: 'v1.2.3',
          IMAGE_METADATA_PATH: outputPath,
          CI_COMMIT_SHA: 'abc123',
          CI_PIPELINE_ID: '42',
          API_IMAGE_DIGEST: 'sha256:api',
          WEB_IMAGE_DIGEST: 'sha256:web',
          BOT_IMAGE_DIGEST: 'sha256:bot',
        },
        write: () => undefined,
      });

      expect(code).toBe(0);
      expect(JSON.parse(readFileSync(outputPath, 'utf8')).imageTag).toBe('v1.2.3');
    } finally {
      rmSync(root, { recursive: true, force: true });
    }
  });

  it('deploys Dokploy dev compose with injected deployment dependency', async () => {
    const deployments: unknown[] = [];
    const code = await runCli(['deploy-dokploy', 'dev'], {
      env: {
        CI_REGISTRY_IMAGE: 'registry.example.com/group/app',
        IMAGE_TAG: 'develop-abc123',
        WEB_API_BASE_URL: 'https://api.dev.example.com',
        DOKPLOY_DEV_URL: 'https://dokploy-dev.example.com',
        DOKPLOY_DEV_API_KEY: 'secret',
        DOKPLOY_DEV_COMPOSE_ID: 'cmp_dev',
      },
      deployDokploy: async (input) => {
        deployments.push(input);
      },
    });

    expect(code).toBe(0);
    expect(deployments).toEqual([
      {
        baseUrl: 'https://dokploy-dev.example.com',
        apiKey: 'secret',
        composeId: 'cmp_dev',
        imageEnv: {
          IMAGE_TAG: 'develop-abc123',
          API_IMAGE: 'registry.example.com/group/app/api:develop-abc123',
          WEB_IMAGE: 'registry.example.com/group/app/web:develop-abc123',
          BOT_IMAGE: 'registry.example.com/group/app/bot:develop-abc123',
          WEB_API_BASE_URL: 'https://api.dev.example.com',
        },
      },
    ]);
  });

  it('requires WEB_API_BASE_URL for Dokploy deploy even when legacy public env exists', async () => {
    const errors: string[] = [];
    const deployments: unknown[] = [];
    const code = await runCli(['deploy-dokploy', 'dev'], {
      env: {
        CI_REGISTRY_IMAGE: 'registry.example.com/group/app',
        IMAGE_TAG: 'develop-abc123',
        NEXT_PUBLIC_API_URL: 'https://api.dev.example.com/graphql',
        DOKPLOY_DEV_URL: 'https://dokploy-dev.example.com',
        DOKPLOY_DEV_API_KEY: 'secret',
        DOKPLOY_DEV_COMPOSE_ID: 'cmp_dev',
      },
      deployDokploy: async (input) => {
        deployments.push(input);
      },
      writeError: (line) => errors.push(line),
    });

    expect(code).toBe(1);
    expect(errors).toEqual(['Missing required CI variable: WEB_API_BASE_URL']);
    expect(deployments).toEqual([]);
  });

  it('deploys Dokploy prod compose with SemVer release images', async () => {
    const deployments: unknown[] = [];
    const code = await runCli(['deploy-dokploy', 'prod'], {
      env: {
        CI_REGISTRY_IMAGE: 'registry.example.com/group/app',
        CI_COMMIT_TAG: 'v1.2.3',
        WEB_API_BASE_URL: 'https://api.example.com',
        DOKPLOY_PROD_URL: 'https://dokploy.example.com',
        DOKPLOY_PROD_API_KEY: 'secret',
        DOKPLOY_PROD_COMPOSE_ID: 'cmp_123',
      },
      deployDokploy: async (input) => {
        deployments.push(input);
      },
    });

    expect(code).toBe(0);
    expect(deployments).toEqual([
      {
        baseUrl: 'https://dokploy.example.com',
        apiKey: 'secret',
        composeId: 'cmp_123',
        imageEnv: {
          IMAGE_TAG: 'v1.2.3',
          API_IMAGE: 'registry.example.com/group/app/api:v1.2.3',
          WEB_IMAGE: 'registry.example.com/group/app/web:v1.2.3',
          BOT_IMAGE: 'registry.example.com/group/app/bot:v1.2.3',
          WEB_API_BASE_URL: 'https://api.example.com',
        },
      },
    ]);
  });

  it('uses the default Dokploy deployer with stubbed fetch', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn(async (url: string) => {
        if (url.includes('/api/compose.one')) {
          return {
            ok: true,
            status: 200,
            text: () => Promise.resolve(JSON.stringify({ env: '' })),
          };
        }
        return {
          ok: true,
          status: 200,
          text: () => Promise.resolve(JSON.stringify({ ok: true })),
        };
      }),
    );

    const code = await runCli(['deploy-dokploy', 'dev'], {
      env: {
        CI_REGISTRY_IMAGE: 'registry.example.com/group/app',
        IMAGE_TAG: 'develop-abc123',
        WEB_API_BASE_URL: 'https://api.dev.example.com',
        DOKPLOY_DEV_URL: 'https://dokploy.example.com',
        DOKPLOY_DEV_API_KEY: 'secret',
        DOKPLOY_DEV_COMPOSE_ID: 'cmp_123',
      },
    });

    expect(code).toBe(0);
    expect(fetch).toHaveBeenCalledTimes(3);
  });

  it('prints unknown command errors through default error writer', async () => {
    const error = vi.spyOn(console, 'error').mockImplementation(() => undefined);

    const code = await runCli([]);

    expect(code).toBe(1);
    expect(error).toHaveBeenCalledWith('Unknown CI helper command: <empty>');
  });

  it('prints non-Error failures from injected dependencies', async () => {
    const errors: string[] = [];
    const code = await runCli(['deploy-dokploy', 'dev'], {
      env: {
        CI_REGISTRY_IMAGE: 'registry.example.com/group/app',
        IMAGE_TAG: 'develop-abc123',
        WEB_API_BASE_URL: 'https://api.dev.example.com',
        DOKPLOY_DEV_URL: 'https://dokploy.example.com',
        DOKPLOY_DEV_API_KEY: 'secret',
        DOKPLOY_DEV_COMPOSE_ID: 'cmp_123',
      },
      deployDokploy: async () => {
        throw 'string failure';
      },
      writeError: (line) => errors.push(line),
    });

    expect(code).toBe(1);
    expect(errors).toEqual(['string failure']);
  });
});
