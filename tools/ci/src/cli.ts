// FILE: tools/ci/src/cli.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide the CI command-line entrypoint for affected checks, release validation, image metadata, and Dokploy deploys.
//   SCOPE: Owns command routing and injected side effects; excludes pure helper implementations and Dokploy HTTP details.
//   DEPENDS: node:child_process, node:fs, node:path, tools/ci/src/core.ts, tools/ci/src/dokploy.ts.
//   LINKS: M-CI-CD / V-M-CI-CD.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   CliDeps - Injectable dependencies for deterministic CLI tests and default runtime wiring.
//   runCli - Executes supported CI helper commands with injectable process, filesystem, and Dokploy dependencies.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Requires WEB_API_BASE_URL for Dokploy deploys and injects it with image refs.
// END_CHANGE_SUMMARY

import { spawnSync } from 'node:child_process';
import { mkdirSync, writeFileSync } from 'node:fs';
import { dirname } from 'node:path';
import {
  assertReleaseTag,
  createImageMetadata,
  releaseTagPattern,
  renderDokployDeployEnv,
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
  makeDir?: (path: string) => void;
  writeFile?: (path: string, contents: string) => void;
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

/* v8 ignore start -- real git is exercised by CI smoke commands; unit tests inject it. */
function defaultGitRunner(gitArgs: string[]): number {
  const result = spawnSync('git', gitArgs, { stdio: 'inherit' });

  return typeof result.status === 'number' ? result.status : 1;
}
/* v8 ignore stop */

function gitRunnerFromDeps(deps: CliDeps): (args: string[]) => number {
  if (deps.runGit) {
    return deps.runGit;
  }

  return defaultGitRunner;
}

function assertTagReachableFromMain(tag: string, runGit: (args: string[]) => number): void {
  if (runGit(['fetch', 'origin', 'main:refs/remotes/origin/main', '--tags']) !== 0) {
    throw new Error('Failed to fetch origin/main and tags');
  }

  if (runGit(['merge-base', '--is-ancestor', `${tag}^{commit}`, 'origin/main']) !== 0) {
    throw new Error(`Release tag ${tag} is not reachable from origin/main`);
  }
}

function toErrorMessage(error: unknown): string {
  return error instanceof Error ? error.message : String(error);
}

export async function runCli(args: string[], deps: CliDeps = {}): Promise<number> {
  const env = deps.env || process.env;
  const write = deps.write || console.log;
  const writeError = deps.writeError || console.error;

  try {
    const command = args[0];
    if (command === 'affected-base') {
      write(resolveAffectedBase(env));
      return 0;
    }

    if (command === 'assert-release-tag') {
      const tag = assertReleaseTag(requireEnv(env, 'CI_COMMIT_TAG'));
      assertTagReachableFromMain(tag, gitRunnerFromDeps(deps));
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
      const makeDir = deps.makeDir || ((path: string) => mkdirSync(path, { recursive: true }));
      const writeFile =
        deps.writeFile || ((path: string, contents: string) => writeFileSync(path, contents));

      makeDir(dirname(outputPath));
      writeFile(outputPath, `${JSON.stringify(metadata, null, 2)}\n`);
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
      const deployDokploy = deps.deployDokploy || deployDokployCompose;

      await deployDokploy({
        ...dokployEnv,
        imageEnv: renderDokployDeployEnv(
          requireEnv(env, 'CI_REGISTRY_IMAGE'),
          imageTag,
          requireEnv(env, 'WEB_API_BASE_URL'),
        ),
      });
      return 0;
    }

    throw new Error(`Unknown CI helper command: ${command || '<empty>'}`);
  } catch (error) {
    writeError(toErrorMessage(error));
    return 1;
  }
}

/* v8 ignore next 4 -- Bun entrypoint; runCli is covered through injected dependencies. */
if (import.meta.main) {
  const code = await runCli(process.argv.slice(2));
  process.exit(code);
}
