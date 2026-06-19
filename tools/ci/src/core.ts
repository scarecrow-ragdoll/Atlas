// FILE: tools/ci/src/core.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide deterministic CI helper primitives for affected ranges, image refs, deploy envs, release metadata, and log redaction.
//   SCOPE: Pure data transformations and validation only; excludes process execution, filesystem writes, and Dokploy HTTP calls.
//   DEPENDS: none.
//   LINKS: M-CI-CD / V-M-CI-CD.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   releaseTagPattern - Strict SemVer tag pattern accepted by release deploys.
//   services - Deployable service registry and Dockerfile ownership.
//   ServiceName - Deployable service name union derived from services.
//   EnvMap - CI environment map accepted by helper functions.
//   ImageMetadataInput - Input contract for release image metadata generation.
//   resolveAffectedBase - Selects the Git base ref for affected Nx checks.
//   requireEnv - Fails fast when required CI variables are missing.
//   assertReleaseTag - Validates release tags before production deployment.
//   buildImageRefs - Builds per-service registry image refs.
//   renderDokployImageEnv - Renders image-only Dokploy env values.
//   renderDokployDeployEnv - Renders image values plus public web runtime API base URL.
//   createImageMetadata - Creates release image metadata for CI artifacts.
//   redactValue - Redacts secret-looking values for logs.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Dokploy deploy env rendering with WEB_API_BASE_URL for public Next web.
// END_CHANGE_SUMMARY

export const releaseTagPattern = /^v[0-9]+\.[0-9]+\.[0-9]+$/;

export const services = [
  { name: 'api', dockerfile: 'docker/api.Dockerfile' },
  { name: 'web', dockerfile: 'docker/web.Dockerfile' },
  { name: 'bot', dockerfile: 'docker/bot.Dockerfile' },
] as const;

export type ServiceName = (typeof services)[number]['name'];

export type EnvMap = Record<string, string | undefined>;

export type ImageMetadataInput = {
  registryImage: string;
  imageTag: string;
  commitSha: string;
  pipelineId: string;
  releaseTag?: string;
  digests: Record<ServiceName, string>;
};

export function resolveAffectedBase(env: EnvMap): string {
  if (env.CI_MERGE_REQUEST_DIFF_BASE_SHA) {
    return env.CI_MERGE_REQUEST_DIFF_BASE_SHA;
  }

  return `origin/${env.CI_DEFAULT_BRANCH || 'main'}`;
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

// START_CONTRACT: renderDokployDeployEnv
//   PURPOSE: Render the full Dokploy env update set for image refs and public web runtime REST proxy config.
//   INPUTS: { registryImage: string - GitLab registry namespace; imageTag: string - image tag to deploy; webApiBaseUrl: string - server-side REST API base URL for Next web }
//   OUTPUTS: { Record<string, string> - Dokploy env keys to append or replace }
//   SIDE_EFFECTS: none.
//   LINKS: M-CI-CD / V-M-CI-CD.
// END_CONTRACT: renderDokployDeployEnv
export function renderDokployDeployEnv(
  registryImage: string,
  imageTag: string,
  webApiBaseUrl: string,
): Record<string, string> {
  return {
    ...renderDokployImageEnv(registryImage, imageTag),
    WEB_API_BASE_URL: webApiBaseUrl,
  };
}

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
