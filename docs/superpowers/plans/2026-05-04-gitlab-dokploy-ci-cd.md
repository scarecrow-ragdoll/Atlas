# GitLab Dokploy CI/CD Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a GitLab-first CI/CD baseline where merge requests run fast checks, `develop` deploys automatically to Dokploy dev, `main` validates release candidates, and SemVer tags from `main` produce manually approved Dokploy production releases.

**Architecture:** GitLab CI remains the orchestrator for validation, full gates, Docker builds, registry publishing, release metadata, and deployment approval. Dokploy is treated as the runtime deployment target for Docker Compose stacks and receives CI-built image references through a `compose.update` then `compose.deploy` flow. Small TypeScript CI helpers hold testable release, image, and Dokploy logic so `.gitlab-ci.yml` stays readable.

**Tech Stack:** GitLab CI, GitLab Container Registry, Docker Buildx, Dokploy Compose API, Bun, Nx, Go 1.25, Next.js, Playwright, Vitest, GRACE XML.

---

## Source Spec

- Design: `docs/superpowers/specs/2026-05-04-gitlab-dokploy-ci-cd-design.md`
- Approved behavior:
  - MR: fast affected checks only.
  - `develop`: full gate, build/push images, auto deploy to Dokploy dev.
  - `main`: full gate as release-candidate validation, no production deploy job.
  - `vX.Y.Z`: assert tag is reachable from `origin/main`, full gate, build/push release images, GitLab Release, manual `deploy:prod`.
  - Dokploy deploy helper updates compose env or compose file before triggering compose deploy.

## File Structure

- Modify `.gitlab-ci.yml`: replace the stale MR-only pipeline with GitLab-first CI/CD stages and rules.
- Create `tools/ci/package.json`: workspace package marker for CI tooling.
- Create `tools/ci/project.json`: Nx target for CI helper tests.
- Create `tools/ci/src/core.ts`: pure helpers for release tags, affected base selection, image refs, metadata, and redaction.
- Create `tools/ci/src/core.test.ts`: unit tests for the pure CI helpers.
- Create `tools/ci/src/dokploy.ts`: Dokploy Compose API client with injectable fetch.
- Create `tools/ci/src/dokploy.test.ts`: unit tests for update/deploy request order and redaction.
- Create `tools/ci/src/cli.ts`: Bun-executed CLI entrypoint used by GitLab jobs.
- Create `tools/ci/src/cli.test.ts`: unit tests for CLI command dispatch without real Git, Docker, or Dokploy calls.
- Modify `tools/vitest.config.ts`: include `tools/ci/src/**/*.ts` in tool coverage.
- Create `deploy/dokploy/docker-compose.template.yml`: environment-driven Dokploy Compose stack template for `api`, `web`, and `bot`.
- Create `docs/infrastructure/ci-cd.md`: operator-facing CI/CD and Dokploy setup guide.
- Modify `README.md`: link to CI/CD docs and release flow.
- Modify `docs/requirements.xml`: add CI/CD release/deploy use case, constraints, risks, and open questions if still relevant.
- Modify `docs/technology.xml`: add GitLab CI, Dokploy, Docker Buildx, release-cli, and CI helper tooling.
- Modify `docs/development-plan.xml`: add `M-CI-CD` module and `DF-CI-CD-RELEASE` flow.
- Modify `docs/knowledge-graph.xml`: add `M-CI-CD` graph node and cross-links to workspace, coverage, API, web, and bot modules.
- Modify `docs/verification-plan.xml`: add `V-M-CI-CD` checks and phase gate coverage.
- Modify `docs/operational-packets.xml`: add CI/CD checkpoint fields only if implementation introduces CI-specific worker packets.

## Task 1: CI Helper Project and Pure Helper Tests

**Files:**

- Create: `tools/ci/package.json`
- Create: `tools/ci/project.json`
- Create: `tools/ci/src/core.ts`
- Create: `tools/ci/src/core.test.ts`
- Modify: `tools/vitest.config.ts`

- [ ] **Step 1: Write the failing tests for release tags, affected base, image refs, metadata, and redaction**

Create `tools/ci/src/core.test.ts` with this content:

```ts
import {
  buildImageRefs,
  createImageMetadata,
  redactValue,
  releaseTagPattern,
  requireEnv,
  resolveAffectedBase,
  services,
} from './core';

describe('core CI helpers', () => {
  it('accepts strict SemVer release tags only', () => {
    expect(releaseTagPattern.test('v1.2.3')).toBe(true);
    expect(releaseTagPattern.test('1.2.3')).toBe(false);
    expect(releaseTagPattern.test('v1.2')).toBe(false);
    expect(releaseTagPattern.test('v1.2.3-rc.1')).toBe(false);
  });

  it('prefers merge request diff base for affected ranges', () => {
    expect(
      resolveAffectedBase({
        CI_MERGE_REQUEST_DIFF_BASE_SHA: 'abc123',
        CI_DEFAULT_BRANCH: 'main',
      }),
    ).toBe('abc123');
  });

  it('falls back to origin default branch when no merge request base exists', () => {
    expect(resolveAffectedBase({ CI_DEFAULT_BRANCH: 'main' })).toBe('origin/main');
  });

  it('builds image refs for every deployable service', () => {
    expect(buildImageRefs('registry.example.com/group/app', 'v1.2.3')).toEqual({
      api: 'registry.example.com/group/app/api:v1.2.3',
      web: 'registry.example.com/group/app/web:v1.2.3',
      bot: 'registry.example.com/group/app/bot:v1.2.3',
    });
  });

  it('creates metadata entries with digest and pipeline evidence', () => {
    const metadata = createImageMetadata({
      registryImage: 'registry.example.com/group/app',
      imageTag: 'v1.2.3',
      commitSha: 'abc123',
      pipelineId: '42',
      releaseTag: 'v1.2.3',
      digests: {
        api: 'sha256:api',
        web: 'sha256:web',
        bot: 'sha256:bot',
      },
    });

    expect(metadata.services).toHaveLength(services.length);
    expect(metadata.services[0]).toMatchObject({
      service: 'api',
      dockerfile: 'docker/api.Dockerfile',
      image: 'registry.example.com/group/app/api:v1.2.3',
      digest: 'sha256:api',
      commitSha: 'abc123',
      pipelineId: '42',
      releaseTag: 'v1.2.3',
    });
  });

  it('fails fast when required environment values are empty', () => {
    expect(() => requireEnv({}, 'DOKPLOY_PROD_API_KEY')).toThrow(
      'Missing required CI variable: DOKPLOY_PROD_API_KEY',
    );
  });

  it('redacts secret-looking values before logging', () => {
    expect(redactValue('DOKPLOY_PROD_API_KEY', 'secret')).toBe('[redacted]');
    expect(redactValue('API_IMAGE', 'registry.example.com/app/api:v1.2.3')).toBe(
      'registry.example.com/app/api:v1.2.3',
    );
  });
});
```

- [ ] **Step 2: Run the failing CI helper test**

Run:

```bash
bunx vitest run --config tools/vitest.config.ts tools/ci/src/core.test.ts
```

Expected: FAIL with an import error for `./core`.

- [ ] **Step 3: Add the CI helper package and Nx project**

Create `tools/ci/package.json`:

```json
{
  "name": "@monorepo-template/ci-tools",
  "private": true,
  "type": "module"
}
```

Create `tools/ci/project.json`:

```json
{
  "name": "ci-tools",
  "$schema": "../../node_modules/nx/schemas/project-schema.json",
  "sourceRoot": "tools/ci/src",
  "projectType": "library",
  "targets": {
    "test": {
      "executor": "nx:run-commands",
      "options": {
        "command": "bunx vitest run --config tools/vitest.config.ts --coverage tools/ci/src"
      }
    }
  }
}
```

- [ ] **Step 4: Implement the pure helper module**

Create `tools/ci/src/core.ts`:

```ts
export const releaseTagPattern = /^v[0-9]+\.[0-9]+\.[0-9]+$/;

export const services = [
  { name: 'api', dockerfile: 'docker/api.Dockerfile' },
  { name: 'web', dockerfile: 'docker/web.Dockerfile' },
  { name: 'bot', dockerfile: 'docker/bot.Dockerfile' },
] as const;

export type ServiceName = (typeof services)[number]['name'];

export type EnvMap = Record<string, string | undefined>;

export function resolveAffectedBase(env: EnvMap): string {
  return env.CI_MERGE_REQUEST_DIFF_BASE_SHA || `origin/${env.CI_DEFAULT_BRANCH || 'main'}`;
}

export function requireEnv(env: EnvMap, key: string): string {
  const value = env[key];
  if (!value) {
    throw new Error(`Missing required CI variable: ${key}`);
  }
  return value;
}

export function assertReleaseTag(tag: string): string {
  if (!releaseTagPattern.test(tag)) {
    throw new Error(`Release tag must match vX.Y.Z: ${tag}`);
  }
  return tag;
}

export function buildImageRefs(
  registryImage: string,
  imageTag: string,
): Record<ServiceName, string> {
  return Object.fromEntries(
    services.map((service) => [service.name, `${registryImage}/${service.name}:${imageTag}`]),
  ) as Record<ServiceName, string>;
}

export function renderDokployImageEnv(
  registryImage: string,
  imageTag: string,
): Record<string, string> {
  const refs = buildImageRefs(registryImage, imageTag);
  return {
    IMAGE_TAG: imageTag,
    API_IMAGE: refs.api,
    WEB_IMAGE: refs.web,
    BOT_IMAGE: refs.bot,
  };
}

export type ImageMetadataInput = {
  registryImage: string;
  imageTag: string;
  commitSha: string;
  pipelineId: string;
  releaseTag?: string;
  digests: Record<ServiceName, string>;
};

export function createImageMetadata(input: ImageMetadataInput) {
  return {
    imageTag: input.imageTag,
    commitSha: input.commitSha,
    pipelineId: input.pipelineId,
    releaseTag: input.releaseTag || null,
    services: services.map((service) => ({
      service: service.name,
      dockerfile: service.dockerfile,
      image: `${input.registryImage}/${service.name}:${input.imageTag}`,
      digest: input.digests[service.name],
      commitSha: input.commitSha,
      pipelineId: input.pipelineId,
      releaseTag: input.releaseTag || null,
    })),
  };
}

export function redactValue(key: string, value: string): string {
  if (/TOKEN|PASSWORD|SECRET|KEY|WEBHOOK/i.test(key)) {
    return '[redacted]';
  }
  return value;
}
```

- [ ] **Step 5: Add CI helper sources to tools coverage**

Modify `tools/vitest.config.ts` coverage include:

```ts
include: ['tools/nx-go/src/**/*.ts', 'tools/codegen/**/*.ts', 'tools/ci/src/**/*.ts'],
```

- [ ] **Step 6: Run the CI helper tests**

Run:

```bash
bunx vitest run --config tools/vitest.config.ts tools/ci/src/core.test.ts
```

Expected: PASS.

- [ ] **Step 7: Run the full tools coverage target**

Run:

```bash
bunx nx run ci-tools:test
```

Expected: PASS with 100 percent coverage for the new CI helper files.

- [ ] **Step 8: Commit Task 1**

```bash
git add tools/ci/package.json tools/ci/project.json tools/ci/src/core.ts tools/ci/src/core.test.ts tools/vitest.config.ts
git commit -m "feat(ci): add tested ci helper core"
```

## Task 2: Dokploy API Client and CLI Tests

**Files:**

- Create: `tools/ci/src/dokploy.ts`
- Create: `tools/ci/src/dokploy.test.ts`
- Create: `tools/ci/src/cli.ts`
- Create: `tools/ci/src/cli.test.ts`
- Modify: `tools/ci/src/core.ts`
- Modify: `tools/ci/src/core.test.ts`

- [ ] **Step 1: Write failing Dokploy client tests**

Create `tools/ci/src/dokploy.test.ts`:

```ts
import { deployDokployCompose, type FetchLike } from './dokploy';

function okResponse(body: unknown = { ok: true }) {
  return Promise.resolve({
    ok: true,
    status: 200,
    text: () => Promise.resolve(JSON.stringify(body)),
  });
}

describe('deployDokployCompose', () => {
  it('updates compose env before triggering deploy', async () => {
    const calls: Array<{
      url: string;
      method: string;
      headers: Record<string, string>;
      body: unknown;
    }> = [];
    const fetchImpl: FetchLike = async (url, init) => {
      calls.push({
        url: String(url),
        method: init?.method ?? 'GET',
        headers: init?.headers ?? {},
        body: init?.body ? JSON.parse(String(init.body)) : null,
      });
      if (String(url).includes('/api/compose.one')) {
        return okResponse({ env: 'POSTGRES_HOST=db\nIMAGE_TAG=old\nAUTH_JWT_SECRET=keep-me' });
      }
      return okResponse();
    };

    await deployDokployCompose({
      baseUrl: 'https://dokploy.example.com',
      apiKey: 'secret',
      composeId: 'cmp_123',
      imageEnv: {
        IMAGE_TAG: 'v1.2.3',
        API_IMAGE: 'registry/app/api:v1.2.3',
        WEB_IMAGE: 'registry/app/web:v1.2.3',
        BOT_IMAGE: 'registry/app/bot:v1.2.3',
      },
      fetchImpl,
    });

    expect(calls).toHaveLength(3);
    expect(calls[0].url).toBe('https://dokploy.example.com/api/compose.one?composeId=cmp_123');
    expect(calls[0].method).toBe('GET');
    expect(calls[0].headers['x-api-key']).toBe('secret');
    expect(calls[1].url).toBe('https://dokploy.example.com/api/compose.update');
    expect(calls[1].body).toEqual({
      composeId: 'cmp_123',
      env: 'POSTGRES_HOST=db\nIMAGE_TAG=v1.2.3\nAUTH_JWT_SECRET=keep-me\nAPI_IMAGE=registry/app/api:v1.2.3\nWEB_IMAGE=registry/app/web:v1.2.3\nBOT_IMAGE=registry/app/bot:v1.2.3',
    });
    expect(calls[2].url).toBe('https://dokploy.example.com/api/compose.deploy');
    expect(calls[2].body).toEqual({ composeId: 'cmp_123' });
  });

  it('fails when Dokploy update returns an error', async () => {
    const fetchImpl: FetchLike = async (url) => {
      if (String(url).includes('/api/compose.one')) {
        return okResponse({ env: '' });
      }
      return Promise.resolve({
        ok: false,
        status: 500,
        text: () => Promise.resolve('boom'),
      });
    };

    await expect(
      deployDokployCompose({
        baseUrl: 'https://dokploy.example.com',
        apiKey: 'secret',
        composeId: 'cmp_123',
        imageEnv: { IMAGE_TAG: 'v1.2.3' },
        fetchImpl,
      }),
    ).rejects.toThrow('Dokploy compose.update failed with status 500');
  });
});
```

- [ ] **Step 2: Run the failing Dokploy test**

Run:

```bash
bunx vitest run --config tools/vitest.config.ts tools/ci/src/dokploy.test.ts
```

Expected: FAIL with an import error for `./dokploy`.

- [ ] **Step 3: Implement the Dokploy API client**

Create `tools/ci/src/dokploy.ts`:

```ts
export type FetchResponse = {
  ok: boolean;
  status: number;
  text(): Promise<string>;
};

export type FetchLike = (
  url: string,
  init: {
    method: string;
    headers: Record<string, string>;
    body?: string;
  },
) => Promise<FetchResponse>;

export type DeployDokployComposeInput = {
  baseUrl: string;
  apiKey: string;
  composeId: string;
  imageEnv: Record<string, string>;
  fetchImpl?: FetchLike;
};

function normalizeBaseUrl(baseUrl: string): string {
  return baseUrl.replace(/\/+$/, '');
}

function mergeEnv(existingEnv: string, updates: Record<string, string>): string {
  const seen = new Set<string>();
  const mergedLines = existingEnv
    .split('\n')
    .filter((line) => line.trim() !== '')
    .map((line) => {
      const match = line.match(/^([A-Za-z_][A-Za-z0-9_]*)=(.*)$/);
      if (!match) {
        return line;
      }
      const key = match[1];
      if (!(key in updates)) {
        return line;
      }
      seen.add(key);
      return `${key}=${updates[key]}`;
    });

  for (const [key, value] of Object.entries(updates)) {
    if (!seen.has(key)) {
      mergedLines.push(`${key}=${value}`);
    }
  }

  return mergedLines.join('\n');
}

async function callDokploy(
  path: string,
  input: DeployDokployComposeInput,
  options: { method: string; body?: Record<string, unknown> },
): Promise<string> {
  const fetchImpl = input.fetchImpl || fetch;
  const response = await fetchImpl(`${normalizeBaseUrl(input.baseUrl)}/api/${path}`, {
    method: options.method,
    headers: {
      'Content-Type': 'application/json',
      'x-api-key': input.apiKey,
    },
    body: options.body ? JSON.stringify(options.body) : undefined,
  });

  const text = await response.text();
  if (!response.ok) {
    throw new Error(`Dokploy ${path} failed with status ${response.status}: ${text}`);
  }
  return text;
}

async function getComposeEnv(input: DeployDokployComposeInput): Promise<string> {
  const text = await callDokploy(
    `compose.one?composeId=${encodeURIComponent(input.composeId)}`,
    input,
    {
      method: 'GET',
    },
  );
  const data = JSON.parse(text) as { env?: string | null };
  return data.env ?? '';
}

async function updateComposeEnv(input: DeployDokployComposeInput, env: string): Promise<void> {
  await callDokploy('compose.update', input, {
    method: 'POST',
    body: {
      composeId: input.composeId,
      env,
    },
  });
}

async function deployCompose(input: DeployDokployComposeInput): Promise<void> {
  await callDokploy('compose.deploy', input, {
    method: 'POST',
    body: {
      composeId: input.composeId,
    },
  });
}

export async function deployDokployCompose(input: DeployDokployComposeInput): Promise<void> {
  const currentEnv = await getComposeEnv(input);
  await updateComposeEnv(input, mergeEnv(currentEnv, input.imageEnv));
  await deployCompose(input);
}
```

- [ ] **Step 4: Write failing CLI tests**

Create `tools/ci/src/cli.test.ts`:

```ts
import { runCli } from './cli';

describe('runCli', () => {
  it('prints affected base', async () => {
    const output: string[] = [];
    const code = await runCli(['affected-base'], {
      env: { CI_MERGE_REQUEST_DIFF_BASE_SHA: 'abc123' },
      write: (line) => output.push(line),
    });

    expect(code).toBe(0);
    expect(output).toEqual(['abc123']);
  });

  it('rejects non-SemVer release tags', async () => {
    const errors: string[] = [];
    const code = await runCli(['assert-release-tag'], {
      env: { CI_COMMIT_TAG: 'latest' },
      writeError: (line) => errors.push(line),
    });

    expect(code).toBe(1);
    expect(errors[0]).toContain('Release tag must match vX.Y.Z');
  });

  it('deploys Dokploy compose with release images', async () => {
    const calls: string[] = [];
    const code = await runCli(['deploy-dokploy', 'prod'], {
      env: {
        CI_REGISTRY_IMAGE: 'registry.example.com/group/app',
        CI_COMMIT_TAG: 'v1.2.3',
        DOKPLOY_PROD_URL: 'https://dokploy.example.com',
        DOKPLOY_PROD_API_KEY: 'secret',
        DOKPLOY_PROD_COMPOSE_ID: 'cmp_123',
      },
      deployDokploy: async () => {
        calls.push('deploy');
      },
    });

    expect(code).toBe(0);
    expect(calls).toEqual(['deploy']);
  });

  it('fails release tag assertion when git ancestry check fails', async () => {
    const errors: string[] = [];
    const code = await runCli(['assert-release-tag'], {
      env: { CI_COMMIT_TAG: 'v1.2.3' },
      runGit: () => 1,
      writeError: (line) => errors.push(line),
    });

    expect(code).toBe(1);
    expect(errors[0]).toContain('Failed to fetch origin/main and tags');
  });
});
```

- [ ] **Step 5: Implement the CLI entrypoint**

Create `tools/ci/src/cli.ts`:

```ts
import { spawnSync } from 'node:child_process';
import { mkdirSync, writeFileSync } from 'node:fs';
import { dirname } from 'node:path';
import {
  assertReleaseTag,
  createImageMetadata,
  releaseTagPattern,
  renderDokployImageEnv,
  requireEnv,
  resolveAffectedBase,
  type EnvMap,
  type ServiceName,
} from './core';
import { deployDokployCompose, type DeployDokployComposeInput } from './dokploy';

export type CliDeps = {
  env?: EnvMap;
  write?: (line: string) => void;
  writeError?: (line: string) => void;
  runGit?: (args: string[]) => number;
  deployDokploy?: (input: DeployDokployComposeInput) => Promise<void>;
};

function envForTarget(target: 'dev' | 'prod', env: EnvMap) {
  const prefix = target === 'prod' ? 'DOKPLOY_PROD' : 'DOKPLOY_DEV';
  return {
    baseUrl: requireEnv(env, `${prefix}_URL`),
    apiKey: requireEnv(env, `${prefix}_API_KEY`),
    composeId: requireEnv(env, `${prefix}_COMPOSE_ID`),
  };
}

function assertTagReachableFromMain(tag: string, runGit: (args: string[]) => number): void {
  if (runGit(['fetch', 'origin', 'main:refs/remotes/origin/main', '--tags']) !== 0) {
    throw new Error('Failed to fetch origin/main and tags');
  }
  if (runGit(['merge-base', '--is-ancestor', `${tag}^{commit}`, 'origin/main']) !== 0) {
    throw new Error(`Release tag ${tag} is not reachable from origin/main`);
  }
}

export async function runCli(args: string[], deps: CliDeps = {}): Promise<number> {
  const env = deps.env || process.env;
  const write = deps.write || console.log;
  const writeError = deps.writeError || console.error;
  const runGit =
    deps.runGit ||
    ((gitArgs: string[]) => {
      const result = spawnSync('git', gitArgs, { stdio: 'inherit' });
      return typeof result.status === 'number' ? result.status : 1;
    });
  const deployDokploy = deps.deployDokploy || deployDokployCompose;

  try {
    const command = args[0];
    if (command === 'affected-base') {
      write(resolveAffectedBase(env));
      return 0;
    }

    if (command === 'assert-release-tag') {
      const tag = assertReleaseTag(requireEnv(env, 'CI_COMMIT_TAG'));
      assertTagReachableFromMain(tag, runGit);
      return 0;
    }

    if (command === 'write-image-metadata') {
      const registryImage = requireEnv(env, 'CI_REGISTRY_IMAGE');
      const imageTag = requireEnv(env, 'IMAGE_TAG');
      const outputPath = env.IMAGE_METADATA_PATH || 'dist/ci/image-metadata.json';
      const metadata = createImageMetadata({
        registryImage,
        imageTag,
        commitSha: requireEnv(env, 'CI_COMMIT_SHA'),
        pipelineId: requireEnv(env, 'CI_PIPELINE_ID'),
        releaseTag: releaseTagPattern.test(imageTag) ? imageTag : undefined,
        digests: {
          api: requireEnv(env, 'API_IMAGE_DIGEST'),
          web: requireEnv(env, 'WEB_IMAGE_DIGEST'),
          bot: requireEnv(env, 'BOT_IMAGE_DIGEST'),
        } satisfies Record<ServiceName, string>,
      });
      mkdirSync(dirname(outputPath), { recursive: true });
      writeFileSync(outputPath, `${JSON.stringify(metadata, null, 2)}\n`);
      write(outputPath);
      return 0;
    }

    if (command === 'deploy-dokploy') {
      const target = args[1] === 'prod' ? 'prod' : 'dev';
      const imageTag =
        target === 'prod'
          ? assertReleaseTag(requireEnv(env, 'CI_COMMIT_TAG'))
          : requireEnv(env, 'IMAGE_TAG');
      const dokployEnv = envForTarget(target, env);
      await deployDokploy({
        ...dokployEnv,
        imageEnv: renderDokployImageEnv(requireEnv(env, 'CI_REGISTRY_IMAGE'), imageTag),
      });
      return 0;
    }

    throw new Error(`Unknown CI helper command: ${command || '<empty>'}`);
  } catch (error) {
    writeError(error instanceof Error ? error.message : String(error));
    return 1;
  }
}

if (import.meta.main) {
  const code = await runCli(process.argv.slice(2));
  process.exit(code);
}
```

- [ ] **Step 6: Run focused CI helper tests**

Run:

```bash
bunx vitest run --config tools/vitest.config.ts tools/ci/src
```

Expected: PASS.

- [ ] **Step 7: Run full tools coverage**

Run:

```bash
bunx nx run ci-tools:test
```

Expected: PASS with 100 percent coverage.

- [ ] **Step 8: Commit Task 2**

```bash
git add tools/ci/src/dokploy.ts tools/ci/src/dokploy.test.ts tools/ci/src/cli.ts tools/ci/src/cli.test.ts tools/ci/src/core.ts tools/ci/src/core.test.ts
git commit -m "feat(ci): add dokploy deploy helpers"
```

## Task 3: Dokploy Compose Template and Operator Docs

**Files:**

- Create: `deploy/dokploy/docker-compose.template.yml`
- Create: `docs/infrastructure/ci-cd.md`
- Modify: `README.md`

- [ ] **Step 1: Create the Dokploy compose template**

Create `deploy/dokploy/docker-compose.template.yml`:

```yaml
services:
  api:
    image: ${API_IMAGE}
    pull_policy: always
    restart: unless-stopped
    environment:
      SERVER_PORT: '8080'
      POSTGRES_HOST: ${POSTGRES_HOST}
      POSTGRES_PORT: ${POSTGRES_PORT}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_SSLMODE: ${POSTGRES_SSLMODE:-disable}
      REDIS_HOST: ${REDIS_HOST}
      REDIS_PORT: ${REDIS_PORT}
      REDIS_PASSWORD: ${REDIS_PASSWORD:-}
      AUTH_JWT_SECRET: ${AUTH_JWT_SECRET}
      SERVER_ENV: ${SERVER_ENV:-production}
    expose:
      - '8080'

  web:
    image: ${WEB_IMAGE}
    pull_policy: always
    restart: unless-stopped
    environment:
      PORT: '3000'
      NEXT_PUBLIC_API_URL: ${NEXT_PUBLIC_API_URL}
    depends_on:
      - api
    expose:
      - '3000'

  bot:
    image: ${BOT_IMAGE}
    pull_policy: always
    restart: unless-stopped
    environment:
      BOT_TOKEN: ${BOT_TOKEN}
    depends_on:
      - api
```

- [ ] **Step 2: Add CI/CD operator documentation**

Create `docs/infrastructure/ci-cd.md`:

````markdown
# CI/CD and Dokploy

## Pipeline Model

- Merge requests run fast affected checks.
- `develop` runs `bun run verify:coverage`, builds `api`, `web`, and `bot` images, pushes them to GitLab Container Registry, and deploys automatically to Dokploy dev.
- `main` runs `bun run verify:coverage` as release-candidate validation and does not deploy production.
- `vX.Y.Z` tags that are reachable from `origin/main` run the release pipeline, publish immutable release images, create a GitLab Release, and expose manual `deploy:prod`.

## Required GitLab Variables

Set these as environment-scoped variables:

- `DOKPLOY_DEV_URL`
- `DOKPLOY_DEV_API_KEY`
- `DOKPLOY_DEV_COMPOSE_ID`
- `DOKPLOY_PROD_URL`
- `DOKPLOY_PROD_API_KEY`
- `DOKPLOY_PROD_COMPOSE_ID`

Production variables must be protected and masked.

## Dokploy Compose Stack

Use `deploy/dokploy/docker-compose.template.yml` as the base compose file in Dokploy. The GitLab deploy job updates the image variables before deployment:

- `IMAGE_TAG`
- `API_IMAGE`
- `WEB_IMAGE`
- `BOT_IMAGE`

Runtime secrets, external database credentials, domains, TLS settings, and registry pull credentials are configured in Dokploy or protected GitLab variables, not in git. The template compose file expects these runtime values to exist in Dokploy:

- `POSTGRES_HOST`
- `POSTGRES_PORT`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_DB`
- `POSTGRES_SSLMODE`
- `REDIS_HOST`
- `REDIS_PORT`
- `REDIS_PASSWORD`
- `AUTH_JWT_SECRET`
- `NEXT_PUBLIC_API_URL`
- `BOT_TOKEN`

## Production Release Flow

1. Merge into `main`.
2. Wait for the `main` release-candidate pipeline to pass.
3. Create a SemVer tag from the `main` commit:

   ```bash
   git tag v1.2.3
   git push origin v1.2.3
   ```
````

4. Wait for the release pipeline to pass.
5. Approve `deploy:prod` in GitLab.

Production deployment from a branch pipeline is intentionally unavailable.

````

- [ ] **Step 3: Link CI/CD docs from README**

Add this section to `README.md` after "Coverage and E2E Gate":

```markdown
## CI/CD and Dokploy

The template uses GitLab CI and Dokploy for deployment. Merge requests run fast affected checks, `develop` deploys automatically to dev after the full gate, `main` validates release candidates, and production deploys only from SemVer tags such as `v1.2.3` after manual approval.

See [docs/infrastructure/ci-cd.md](docs/infrastructure/ci-cd.md) for required GitLab variables, Dokploy compose setup, image tagging, and the production release flow.
````

- [ ] **Step 4: Run documentation checks**

Run:

```bash
rg -n "DOKPLOY_|deploy:prod|v1\\.2\\.3|origin/main" README.md docs/infrastructure/ci-cd.md deploy/dokploy/docker-compose.template.yml
```

Expected: output includes the required variables, `deploy:prod`, `v1.2.3`, and `origin/main`.

- [ ] **Step 5: Commit Task 3**

```bash
git add deploy/dokploy/docker-compose.template.yml docs/infrastructure/ci-cd.md README.md
git commit -m "docs(ci): document dokploy deployment contract"
```

## Task 4: Replace the GitLab CI Pipeline

**Files:**

- Modify: `.gitlab-ci.yml`

- [ ] **Step 1: Replace `.gitlab-ci.yml` with the new pipeline skeleton**

Replace `.gitlab-ci.yml` with this content:

```yaml
workflow:
  rules:
    - if: '$CI_MERGE_REQUEST_IID'
    - if: '$CI_COMMIT_BRANCH == "develop"'
    - if: '$CI_COMMIT_BRANCH == "main"'
    - if: '$CI_COMMIT_TAG =~ /^v\d+\.\d+\.\d+$/'
    - when: never

stages:
  - prepare
  - validate
  - test
  - coverage
  - build
  - release
  - deploy

variables:
  GIT_DEPTH: '0'
  NX_DAEMON: 'false'
  NX_SKIP_NX_CACHE: '1'
  DOCKER_DRIVER: overlay2
  DOCKER_TLS_CERTDIR: ''
  DOCKER_HOST: tcp://docker:2375
  CI_TOOLCHAIN_IMAGE: mcr.microsoft.com/playwright:v1.48.0-noble

.toolchain:
  image: $CI_TOOLCHAIN_IMAGE
  before_script:
    - apt-get update
    - apt-get install -y --no-install-recommends curl git ca-certificates libxml2-utils docker.io docker-compose-v2
    - curl -fsSL https://bun.sh/install | bash
    - export BUN_INSTALL="$HOME/.bun"
    - export PATH="$BUN_INSTALL/bin:/usr/local/go/bin:$PATH"
    - curl -fsSL https://go.dev/dl/go1.25.0.linux-amd64.tar.gz -o /tmp/go.tar.gz
    - rm -rf /usr/local/go
    - tar -C /usr/local -xzf /tmp/go.tar.gz
    - bun install --frozen-lockfile
    - bun add -g @osovv/grace-cli

.docker_build:
  image: docker:27
  services:
    - name: docker:27-dind
      alias: docker
  before_script:
    - docker login -u "$CI_REGISTRY_USER" -p "$CI_REGISTRY_PASSWORD" "$CI_REGISTRY"
    - mkdir -p dist/ci
    - docker buildx create --use --name template-builder || docker buildx use template-builder

prepare:toolchain:
  extends: .toolchain
  stage: prepare
  script:
    - export BUN_INSTALL="$HOME/.bun"
    - export PATH="$BUN_INSTALL/bin:/usr/local/go/bin:$PATH"
    - node --version
    - bun --version
    - go version
    - docker --version
    - docker compose version
    - xmllint --version
    - grace --version

validate:mr:
  extends: .toolchain
  stage: validate
  rules:
    - if: '$CI_MERGE_REQUEST_IID'
  script:
    - export BUN_INSTALL="$HOME/.bun"
    - export PATH="$BUN_INSTALL/bin:/usr/local/go/bin:$PATH"
    - NX_BASE="$(bun tools/ci/src/cli.ts affected-base)"
    - bunx nx affected --target=lint --base="$NX_BASE" --parallel=1
    - bunx nx affected --target=typecheck --base="$NX_BASE"
    - bunx nx affected --target=test --base="$NX_BASE"
    - bunx nx affected --target=build --base="$NX_BASE"

validate:grace:
  extends: .toolchain
  stage: validate
  rules:
    - if: '$CI_MERGE_REQUEST_IID'
    - if: '$CI_COMMIT_BRANCH == "develop"'
    - if: '$CI_COMMIT_BRANCH == "main"'
    - if: '$CI_COMMIT_TAG =~ /^v\d+\.\d+\.\d+$/'
  script:
    - export BUN_INSTALL="$HOME/.bun"
    - export PATH="$BUN_INSTALL/bin:/usr/local/go/bin:$PATH"
    - xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
    - grace lint --path .

verify:full:
  extends: .toolchain
  stage: coverage
  services:
    - name: docker:27-dind
      alias: docker
  rules:
    - if: '$CI_COMMIT_BRANCH == "develop"'
    - if: '$CI_COMMIT_BRANCH == "main"'
    - if: '$CI_COMMIT_TAG =~ /^v\d+\.\d+\.\d+$/'
  script:
    - export BUN_INSTALL="$HOME/.bun"
    - export PATH="$BUN_INSTALL/bin:/usr/local/go/bin:$PATH"
    - bun run verify:coverage
  artifacts:
    when: always
    expire_in: 14 days
    paths:
      - dist/coverage
      - dist/test-results/web-e2e
      - dist/playwright-report/web

assert:release-tag:
  extends: .toolchain
  stage: validate
  rules:
    - if: '$CI_COMMIT_TAG =~ /^v\d+\.\d+\.\d+$/'
  script:
    - export BUN_INSTALL="$HOME/.bun"
    - export PATH="$BUN_INSTALL/bin:/usr/local/go/bin:$PATH"
    - bun tools/ci/src/cli.ts assert-release-tag

build:images:dev:
  extends: .docker_build
  stage: build
  needs:
    - verify:full
  rules:
    - if: '$CI_COMMIT_BRANCH == "develop"'
  variables:
    IMAGE_TAG: develop-$CI_COMMIT_SHORT_SHA
  script:
    - docker buildx build --pull --target prod -f docker/api.Dockerfile -t "$CI_REGISTRY_IMAGE/api:$IMAGE_TAG" -t "$CI_REGISTRY_IMAGE/api:dev-latest" --metadata-file dist/ci/api-image.json --push .
    - docker buildx build --pull --target prod -f docker/web.Dockerfile -t "$CI_REGISTRY_IMAGE/web:$IMAGE_TAG" -t "$CI_REGISTRY_IMAGE/web:dev-latest" --metadata-file dist/ci/web-image.json --push .
    - docker buildx build --pull --target prod -f docker/bot.Dockerfile -t "$CI_REGISTRY_IMAGE/bot:$IMAGE_TAG" -t "$CI_REGISTRY_IMAGE/bot:dev-latest" --metadata-file dist/ci/bot-image.json --push .
  artifacts:
    expire_in: 14 days
    paths:
      - dist/ci

build:images:release:
  extends: .docker_build
  stage: build
  needs:
    - verify:full
    - assert:release-tag
  rules:
    - if: '$CI_COMMIT_TAG =~ /^v\d+\.\d+\.\d+$/'
  variables:
    IMAGE_TAG: $CI_COMMIT_TAG
  script:
    - apk add --no-cache jq curl bash
    - |
      for service in api web bot; do
        image="$CI_REGISTRY_IMAGE/$service:$IMAGE_TAG"
        if docker manifest inspect "$image" >/dev/null 2>&1; then
          echo "Release image already exists and will not be overwritten: $image" >&2
          exit 1
        fi
      done
    - docker buildx build --pull --target prod -f docker/api.Dockerfile -t "$CI_REGISTRY_IMAGE/api:$IMAGE_TAG" --metadata-file dist/ci/api-image.json --push .
    - docker buildx build --pull --target prod -f docker/web.Dockerfile -t "$CI_REGISTRY_IMAGE/web:$IMAGE_TAG" --metadata-file dist/ci/web-image.json --push .
    - docker buildx build --pull --target prod -f docker/bot.Dockerfile -t "$CI_REGISTRY_IMAGE/bot:$IMAGE_TAG" --metadata-file dist/ci/bot-image.json --push .
    - export API_IMAGE_DIGEST="$(jq -r '.["containerimage.digest"]' dist/ci/api-image.json)"
    - export WEB_IMAGE_DIGEST="$(jq -r '.["containerimage.digest"]' dist/ci/web-image.json)"
    - export BOT_IMAGE_DIGEST="$(jq -r '.["containerimage.digest"]' dist/ci/bot-image.json)"
    - curl -fsSL https://bun.sh/install | bash
    - export BUN_INSTALL="$HOME/.bun"
    - export PATH="$BUN_INSTALL/bin:$PATH"
    - bun tools/ci/src/cli.ts write-image-metadata
  artifacts:
    expire_in: 90 days
    paths:
      - dist/ci

release:gitlab:
  image: registry.gitlab.com/gitlab-org/release-cli:latest
  stage: release
  needs:
    - build:images:release
  rules:
    - if: '$CI_COMMIT_TAG =~ /^v\d+\.\d+\.\d+$/'
  script:
    - echo "Creating GitLab Release for $CI_COMMIT_TAG"
  release:
    tag_name: '$CI_COMMIT_TAG'
    name: '$CI_COMMIT_TAG'
    description: 'Release $CI_COMMIT_TAG from $CI_COMMIT_SHA. Image metadata is attached as a pipeline artifact.'

deploy:dev:
  extends: .toolchain
  stage: deploy
  needs:
    - build:images:dev
  rules:
    - if: '$CI_COMMIT_BRANCH == "develop"'
  variables:
    IMAGE_TAG: develop-$CI_COMMIT_SHORT_SHA
  script:
    - export BUN_INSTALL="$HOME/.bun"
    - export PATH="$BUN_INSTALL/bin:/usr/local/go/bin:$PATH"
    - bun tools/ci/src/cli.ts deploy-dokploy dev
  environment:
    name: dev

deploy:prod:
  extends: .toolchain
  stage: deploy
  needs:
    - release:gitlab
  rules:
    - if: '$CI_COMMIT_TAG =~ /^v\d+\.\d+\.\d+$/'
      when: manual
  script:
    - export BUN_INSTALL="$HOME/.bun"
    - export PATH="$BUN_INSTALL/bin:/usr/local/go/bin:$PATH"
    - bun tools/ci/src/cli.ts deploy-dokploy prod
  environment:
    name: production
```

- [ ] **Step 2: Verify release build job package assumptions before merging**

The `build:images:release` job uses the `docker:27` image and installs `jq`, `curl`, `bash`, and Bun before writing metadata. It must also fail before build when a release image tag already exists. Verify these commands are ordered exactly as shown in Step 1:

```bash
rg -n "apk add --no-cache jq curl bash|docker manifest inspect|Release image already exists|export API_IMAGE_DIGEST|bun tools/ci/src/cli.ts write-image-metadata" .gitlab-ci.yml
```

Expected: output shows package installation, immutable tag check before the first `docker buildx build`, digest export before `write-image-metadata`, and the release image overwrite error text.

- [ ] **Step 3: Add a focused CI YAML parse check**

Run:

```bash
ruby -e "require 'yaml'; YAML.load_file('.gitlab-ci.yml'); puts 'gitlab ci yaml parses'"
```

Expected: PASS and output includes `gitlab ci yaml parses`. If Ruby is not installed locally, run:

```bash
bunx prettier --check .gitlab-ci.yml
```

Expected: PASS with `.gitlab-ci.yml` reported as formatted.

- [ ] **Step 4: Run local helper commands used by CI**

Run:

```bash
bun tools/ci/src/cli.ts affected-base
```

Expected: output is `origin/main` in a local shell without MR environment variables.

Run:

```bash
CI_COMMIT_TAG=latest bun tools/ci/src/cli.ts assert-release-tag
```

Expected: FAIL with `Release tag must match vX.Y.Z: latest`.

- [ ] **Step 5: Commit Task 4**

```bash
git add .gitlab-ci.yml
git commit -m "feat(ci): add gitlab dokploy pipeline"
```

## Task 5: GRACE Contract Updates

**Files:**

- Modify: `docs/requirements.xml`
- Modify: `docs/technology.xml`
- Modify: `docs/development-plan.xml`
- Modify: `docs/knowledge-graph.xml`
- Modify: `docs/verification-plan.xml`
- Modify: `docs/operational-packets.xml` only if CI-specific packets are added

- [ ] **Step 1: Add CI/CD use case to requirements**

In `docs/requirements.xml`, add a new high-priority use case after `UC-005`:

```xml
    <UC-006>
      <Actor>Operator</Actor>
      <Action>Runs GitLab CI/CD to validate, release, and deploy template services through Dokploy.</Action>
      <Goal>Provide a reusable deployment baseline for downstream repositories with dev auto deploys and tag-gated production releases.</Goal>
      <Preconditions>GitLab CI variables for registry and Dokploy are configured, Docker-in-Docker is available, and the release tag is reachable from `origin/main`.</Preconditions>
      <AcceptanceCriteria>Merge requests run affected checks, `develop` deploys to Dokploy dev after `bun run verify:coverage`, `main` performs release-candidate validation without production deploy, and SemVer tags expose manual `deploy:prod` after release images are pushed.</AcceptanceCriteria>
      <Priority>high</Priority>
      <RelatedFlows>DF-CI-CD-RELEASE</RelatedFlows>
    </UC-006>
```

Add constraints:

```xml
    <constraint-8>Production deployment is allowed only from SemVer tags reachable from `origin/main`.</constraint-8>
    <constraint-9>Dokploy deployment must use GitLab-built images and must update the target compose image variables before deploy.</constraint-9>
```

Add risks:

```xml
    <risk-7>CI runner images can drift from the required Bun, Go, Docker, Playwright, xmllint, and GRACE toolchain if bootstrap checks are weakened.</risk-7>
    <risk-8>Dokploy can redeploy stale images if CI triggers compose deploy without first updating the target compose image variables.</risk-8>
```

- [ ] **Step 2: Add technology entries**

In `docs/technology.xml`, add tooling entries:

```xml
    <tool name="ci-provider" value="GitLab CI" version="GitLab SaaS or compatible self-managed GitLab" />
    <tool name="deployment-target" value="Dokploy Docker Compose" version="current Dokploy API" />
    <tool name="image-builder" value="Docker Buildx" version="Docker 27 compatible" />
    <tool name="release-tool" value="GitLab release-cli" version="latest compatible" />
```

Add phase-level commands:

```xml
      <command>bun tools/ci/src/cli.ts affected-base</command>
      <command>bun tools/ci/src/cli.ts assert-release-tag</command>
      <command>bun tools/ci/src/cli.ts deploy-dokploy dev</command>
      <command>bun tools/ci/src/cli.ts deploy-dokploy prod</command>
```

- [ ] **Step 3: Add `M-CI-CD` module to development plan**

In `docs/development-plan.xml`, add a module after `M-COVERAGE-GATE`:

```xml
    <M-CI-CD NAME="GitLabDokployCICD" TYPE="UTILITY" LAYER="0" ORDER="1.6" STATUS="implemented">
      <contract>
        <purpose>Validate merge requests, run release gates, build service images, create GitLab releases, and trigger Dokploy Docker Compose deployments.</purpose>
        <inputs>
          <param name="git-refs" type="merge requests, develop, main, and SemVer tags" />
          <param name="ci-variables" type="GitLab registry and Dokploy environment variables" />
          <param name="dockerfiles" type="docker/api.Dockerfile, docker/web.Dockerfile, docker/bot.Dockerfile" />
        </inputs>
        <outputs>
          <param name="pipeline-evidence" type="coverage, e2e, image metadata, release, and deployment artifacts" />
          <param name="dokploy-deployments" type="dev and production Dokploy compose deployments" />
        </outputs>
        <errors>
          <error code="CI_TOOLCHAIN_UNAVAILABLE" />
          <error code="RELEASE_TAG_NOT_ON_MAIN" />
          <error code="IMAGE_BUILD_FAILED" />
          <error code="DOKPLOY_UPDATE_FAILED" />
          <error code="DOKPLOY_DEPLOY_FAILED" />
        </errors>
      </contract>
      <interface>
        <export-gitlab-pipeline PURPOSE="Expose GitLab CI stages and jobs." />
        <export-ci-helpers PURPOSE="Expose tested helper commands for affected base, release tag validation, image metadata, and Dokploy deployment." />
        <export-dokploy-template PURPOSE="Expose environment-driven Dokploy compose stack template." />
      </interface>
      <depends>M-WORKSPACE, M-COVERAGE-GATE, M-API, M-WEB, M-BOT</depends>
      <target>
        <source>.gitlab-ci.yml</source>
        <source>tools/ci</source>
        <source>deploy/dokploy</source>
        <source>docs/infrastructure/ci-cd.md</source>
        <tests>bunx nx run ci-tools:test and CI helper command smoke tests</tests>
      </target>
      <observability>
        <log-prefix>[CI-CD]</log-prefix>
        <critical-block>BLOCK_RELEASE_AND_DEPLOY</critical-block>
      </observability>
      <verification-ref>V-M-CI-CD</verification-ref>
    </M-CI-CD>
```

Add flow:

```xml
    <DF-CI-CD-RELEASE NAME="GitLabDokployReleaseFlow" TRIGGER="Merge request, develop push, main push, or SemVer tag">
      <step-1>Merge requests run deterministic affected checks using the merge request diff base SHA.</step-1>
      <step-2>`develop` and `main` run `bun run verify:coverage` and publish coverage plus Playwright evidence.</step-2>
      <step-3>`develop` builds and pushes dev images, then updates and deploys the Dokploy dev compose stack.</step-3>
      <step-4>SemVer tags reachable from `origin/main` build immutable release images and create GitLab Release metadata.</step-4>
      <step-5>Manual `deploy:prod` updates the Dokploy production compose image variables and triggers compose deploy.</step-5>
      <evidence>GitLab pipeline status, image metadata artifact, GitLab Release, and Dokploy update/deploy response summary.</evidence>
    </DF-CI-CD-RELEASE>
```

- [ ] **Step 4: Add knowledge graph entry**

In `docs/knowledge-graph.xml`, add:

```xml
    <M-CI-CD NAME="GitLabDokployCICD" TYPE="UTILITY" STATUS="implemented">
      <purpose>GitLab CI/CD pipeline, tested CI helpers, image publishing, release metadata, and Dokploy compose deployment contract.</purpose>
      <path>.gitlab-ci.yml</path>
      <path>tools/ci</path>
      <path>deploy/dokploy</path>
      <path>docs/infrastructure/ci-cd.md</path>
      <depends>M-WORKSPACE, M-COVERAGE-GATE, M-API, M-WEB, M-BOT</depends>
      <verification-ref>V-M-CI-CD</verification-ref>
      <annotations>
        <export-pipeline PURPOSE="Run GitLab CI workflow for MR, develop, main, and release tags." />
        <export-releaseGate PURPOSE="Assert production tags are SemVer and reachable from origin/main." />
        <export-imageMetadata PURPOSE="Record service image refs, digests, commit SHA, and pipeline ID." />
        <export-dokployDeploy PURPOSE="Update Dokploy compose image variables and trigger compose deployment." />
      </annotations>
    </M-CI-CD>
```

Add cross-links:

```xml
    <CrossLink from="M-CI-CD" to="M-WORKSPACE" relation="runs Bun and Nx scripts" />
    <CrossLink from="M-CI-CD" to="M-COVERAGE-GATE" relation="requires full coverage and e2e gate before deploy-capable branches and release tags" />
    <CrossLink from="M-CI-CD" to="M-API" relation="builds and deploys API image" />
    <CrossLink from="M-CI-CD" to="M-WEB" relation="builds and deploys web image" />
    <CrossLink from="M-CI-CD" to="M-BOT" relation="builds and deploys bot image" />
```

- [ ] **Step 5: Add verification plan entry**

In `docs/verification-plan.xml`, add `V-M-CI-CD`:

```xml
    <V-M-CI-CD MODULE="M-CI-CD" PRIORITY="high">
      <test-files>
        <file>.gitlab-ci.yml</file>
        <file>tools/ci/src/core.test.ts</file>
        <file>tools/ci/src/dokploy.test.ts</file>
        <file>tools/ci/src/cli.test.ts</file>
        <file>deploy/dokploy/docker-compose.template.yml</file>
        <file>docs/infrastructure/ci-cd.md</file>
      </test-files>
      <module-checks>
        <check-1>bunx nx run ci-tools:test</check-1>
        <check-2>bun tools/ci/src/cli.ts affected-base</check-2>
        <check-3>CI_COMMIT_TAG=latest bun tools/ci/src/cli.ts assert-release-tag</check-3>
        <check-4>xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml</check-4>
        <check-5>grace lint --path .</check-5>
      </module-checks>
      <scenarios>
        <scenario-1 kind="success">Merge request pipelines choose the merge request diff base for affected checks.</scenario-1>
        <scenario-2 kind="success">`develop` pipelines run the full gate and deploy Dokploy dev after image push.</scenario-2>
        <scenario-3 kind="success">SemVer tag pipelines validate tag ancestry, publish release images, create release metadata, and expose manual production deploy.</scenario-3>
        <scenario-4 kind="failure">Production deploy is unavailable from branch pipelines.</scenario-4>
        <scenario-5 kind="failure">Dokploy deploy fails when compose update or compose deploy fails.</scenario-5>
      </scenarios>
      <required-trace-assertions>
        <assertion-1>`main` must not expose a production deploy job.</assertion-1>
        <assertion-2>Production deploy requires a SemVer tag reachable from `origin/main`.</assertion-2>
        <assertion-3>Dokploy compose update must happen before compose deploy when immutable image tags are used.</assertion-3>
      </required-trace-assertions>
      <wave-follow-up>Run helper tests and CI YAML parse check after pipeline or helper changes.</wave-follow-up>
      <phase-follow-up>Run full `bun run verify:coverage` before release handoff.</phase-follow-up>
    </V-M-CI-CD>
```

- [ ] **Step 6: Run XML and GRACE validation**

Run:

```bash
xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
grace lint --path .
```

Expected: both commands pass.

- [ ] **Step 7: Commit Task 5**

```bash
git add docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/knowledge-graph.xml docs/verification-plan.xml docs/operational-packets.xml
git commit -m "docs(ci): add gitlab dokploy grace contracts"
```

## Task 6: Focused Verification and Release Safety Review

**Files:**

- Review all files changed by Tasks 1-5.
- Write report: `.tasks/ci-cd-review.md`

- [ ] **Step 1: Run focused helper tests**

Run:

```bash
bunx nx run ci-tools:test
```

Expected: PASS with 100 percent tool coverage.

- [ ] **Step 2: Run root CI helper smoke tests**

Run:

```bash
bun tools/ci/src/cli.ts affected-base
CI_COMMIT_TAG=latest bun tools/ci/src/cli.ts assert-release-tag
```

Expected:

- First command prints `origin/main`.
- Second command exits non-zero and prints `Release tag must match vX.Y.Z: latest`.

- [ ] **Step 3: Run GRACE integrity checks**

Run:

```bash
xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
grace lint --path .
```

Expected: both pass.

- [ ] **Step 4: Run broad template gate**

Run:

```bash
bun run verify:coverage
```

Expected: PASS. Artifacts are written under:

- `dist/coverage`
- `dist/test-results/web-e2e`
- `dist/playwright-report/web`

- [ ] **Step 5: Review pipeline safety rules**

Run:

```bash
rg -n "deploy:prod|CI_COMMIT_BRANCH == \"main\"|CI_COMMIT_TAG|compose.update|compose.deploy|NX_DAEMON|CI_MERGE_REQUEST_DIFF_BASE_SHA" .gitlab-ci.yml tools/ci docs README.md deploy/dokploy
```

Expected:

- `deploy:prod` appears only under a SemVer tag rule.
- `main` appears only in validation or ancestry checks, not production deploy rules.
- `compose.update` appears before `compose.deploy` in helper code or tests.
- `NX_DAEMON=false` appears in CI variables.
- `CI_MERGE_REQUEST_DIFF_BASE_SHA` appears in affected base helper or tests.

- [ ] **Step 6: Write verification report**

Create `.tasks/ci-cd-review.md`:

```markdown
# CI/CD Review

## Scope

- GitLab CI/CD pipeline
- CI helper tooling
- Dokploy compose deployment contract
- GRACE CI/CD module and verification updates

## Commands

- `bunx nx run ci-tools:test`
- `bun tools/ci/src/cli.ts affected-base`
- `CI_COMMIT_TAG=latest bun tools/ci/src/cli.ts assert-release-tag`
- `xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml`
- `grace lint --path .`
- `bun run verify:coverage`
- `rg -n "deploy:prod|CI_COMMIT_BRANCH == \"main\"|CI_COMMIT_TAG|compose.update|compose.deploy|NX_DAEMON|CI_MERGE_REQUEST_DIFF_BASE_SHA" .gitlab-ci.yml tools/ci docs README.md deploy/dokploy`

## Findings

- Production deploy is tag-only.
- `main` has no production deploy job.
- Dokploy compose update precedes deploy.
- Full gate evidence is retained for deploy-capable pipelines.

## Residual Risks

- First real GitLab run must verify runner package availability for Docker Buildx, Docker Compose, Go 1.25, Bun, Playwright, xmllint, and GRACE CLI.
- Real Dokploy API variables and compose IDs are environment-specific and must be configured outside git.
```

- [ ] **Step 7: Commit Task 6**

```bash
git add .tasks/ci-cd-review.md
git commit -m "docs(ci): record ci cd verification evidence"
```

## Task 7: Final Integration Check

**Files:**

- No new source files expected.
- Review final git history and working tree.

- [ ] **Step 1: Check final working tree**

Run:

```bash
git status --short --branch
```

Expected: clean working tree on the implementation branch.

- [ ] **Step 2: Show final commit stack**

Run:

```bash
git log --oneline --max-count=8
```

Expected: commits include:

- `feat(ci): add tested ci helper core`
- `feat(ci): add dokploy deploy helpers`
- `docs(ci): document dokploy deployment contract`
- `feat(ci): add gitlab dokploy pipeline`
- `docs(ci): add gitlab dokploy grace contracts`
- `docs(ci): record ci cd verification evidence`

- [ ] **Step 3: Prepare handoff summary**

Include this in the final implementation handoff:

```text
Implemented GitLab-first, Dokploy-first CI/CD.

Key behavior:
- MR fast affected checks.
- develop full gate, image publish, automatic Dokploy dev deploy.
- main full gate only, no production deploy.
- vX.Y.Z tags from origin/main build release images, create GitLab Release, and expose manual production deploy.
- Dokploy production deploy updates compose image variables before compose deploy.

Verification:
- [list exact commands and pass/fail status]

Residual operational setup:
- Configure GitLab protected variables for Dokploy.
- Configure Dokploy registry pull credentials.
- Run first GitLab pipeline to validate runner toolchain package availability.
```
