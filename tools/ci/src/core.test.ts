// FILE: tools/ci/src/core.test.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verifies pure CI helper contracts for image refs, deploy env rendering, release tags, metadata, and redaction.
//   SCOPE: Covers deterministic core helpers only; excludes CLI side effects and Dokploy HTTP calls.
//   DEPENDS: vitest, tools/ci/src/core.ts.
//   LINKS: M-CI-CD / V-M-CI-CD.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   core CI helpers - Guards image metadata, Dokploy deploy env, affected-base, release-tag, and redaction behavior.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added coverage for Dokploy deploy env rendering with WEB_API_BASE_URL.
// END_CHANGE_SUMMARY

import { describe, expect, it } from 'vitest';
import {
  assertReleaseTag,
  buildImageRefs,
  createImageMetadata,
  redactValue,
  releaseTagPattern,
  renderDokployDeployEnv,
  renderDokployImageEnv,
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
    expect(assertReleaseTag('v10.20.30')).toBe('v10.20.30');
    expect(() => assertReleaseTag('latest')).toThrow('Release tag must match vX.Y.Z: latest');
  });

  it('defines every deployable service with its Dockerfile', () => {
    expect(services).toEqual([
      { name: 'api', dockerfile: 'docker/api.Dockerfile' },
      { name: 'web', dockerfile: 'docker/web.Dockerfile' },
      { name: 'bot', dockerfile: 'docker/bot.Dockerfile' },
    ]);
  });

  it('prefers merge request diff base for affected ranges', () => {
    expect(
      resolveAffectedBase({
        CI_MERGE_REQUEST_DIFF_BASE_SHA: 'abc123',
        CI_DEFAULT_BRANCH: 'develop',
      }),
    ).toBe('abc123');
  });

  it('falls back to origin default branch or origin main when no merge request base exists', () => {
    expect(resolveAffectedBase({ CI_DEFAULT_BRANCH: 'develop' })).toBe('origin/develop');
    expect(resolveAffectedBase({})).toBe('origin/main');
  });

  it('fails fast when required environment values are empty', () => {
    expect(
      requireEnv({ CI_REGISTRY_IMAGE: 'registry.example.com/group/app' }, 'CI_REGISTRY_IMAGE'),
    ).toBe('registry.example.com/group/app');
    expect(() => requireEnv({}, 'DOKPLOY_PROD_API_KEY')).toThrow(
      'Missing required CI variable: DOKPLOY_PROD_API_KEY',
    );
    expect(() => requireEnv({ DOKPLOY_PROD_API_KEY: '' }, 'DOKPLOY_PROD_API_KEY')).toThrow(
      'Missing required CI variable: DOKPLOY_PROD_API_KEY',
    );
  });

  it('builds image refs and Dokploy image environment for every deployable service', () => {
    expect(buildImageRefs('registry.example.com/group/app', 'v1.2.3')).toEqual({
      api: 'registry.example.com/group/app/api:v1.2.3',
      web: 'registry.example.com/group/app/web:v1.2.3',
      bot: 'registry.example.com/group/app/bot:v1.2.3',
    });
    expect(renderDokployImageEnv('registry.example.com/group/app', 'v1.2.3')).toEqual({
      IMAGE_TAG: 'v1.2.3',
      API_IMAGE: 'registry.example.com/group/app/api:v1.2.3',
      WEB_IMAGE: 'registry.example.com/group/app/web:v1.2.3',
      BOT_IMAGE: 'registry.example.com/group/app/bot:v1.2.3',
    });
  });

  it('renders Dokploy deploy environment with public web runtime API base URL', () => {
    expect(
      renderDokployDeployEnv('registry.example.com/group/app', 'v1.2.3', 'https://api.example.com'),
    ).toEqual({
      IMAGE_TAG: 'v1.2.3',
      API_IMAGE: 'registry.example.com/group/app/api:v1.2.3',
      WEB_IMAGE: 'registry.example.com/group/app/web:v1.2.3',
      BOT_IMAGE: 'registry.example.com/group/app/bot:v1.2.3',
      WEB_API_BASE_URL: 'https://api.example.com',
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

    expect(metadata).toMatchObject({
      imageTag: 'v1.2.3',
      commitSha: 'abc123',
      pipelineId: '42',
      releaseTag: 'v1.2.3',
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

  it('uses null release metadata for non-release image tags', () => {
    expect(
      createImageMetadata({
        registryImage: 'registry.example.com/group/app',
        imageTag: 'develop-abc123',
        commitSha: 'abc123',
        pipelineId: '42',
        digests: {
          api: 'sha256:api',
          web: 'sha256:web',
          bot: 'sha256:bot',
        },
      }).releaseTag,
    ).toBeNull();
  });

  it('redacts secret-looking values before logging', () => {
    expect(redactValue('DOKPLOY_PROD_API_KEY', 'secret')).toBe('[redacted]');
    expect(redactValue('BOT_TOKEN', 'secret')).toBe('[redacted]');
    expect(redactValue('API_IMAGE', 'registry.example.com/app/api:v1.2.3')).toBe(
      'registry.example.com/app/api:v1.2.3',
    );
  });
});
