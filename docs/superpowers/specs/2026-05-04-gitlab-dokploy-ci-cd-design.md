# GitLab Dokploy CI/CD Design

Date: 2026-05-04
Status: Approved for implementation planning

## Summary

This template needs a universal CI/CD baseline that downstream repositories can inherit with minimal project-specific changes. The approved direction is GitLab-first CI and Dokploy-first CD.

GitLab CI owns validation, coverage, Docker image builds, registry publishing, release metadata, and environment approval gates. Dokploy owns runtime deployment of Docker Compose stacks. Dokploy must not rebuild production artifacts during deployment; it should pull immutable images that GitLab CI has already built and pushed.

The default environment model is:

- Merge requests run fast affected checks.
- `develop` runs the full gate, builds images, pushes them to GitLab Container Registry, and deploys automatically to Dokploy `dev`.
- `main` runs the full gate as release-candidate validation, but does not deploy production.
- SemVer tags like `v1.2.3`, created from commits reachable from `origin/main`, run the release pipeline and expose a manual `deploy:prod` job.

## Goals

- Provide a production-grade `.gitlab-ci.yml` that works as the default for new repositories derived from this template.
- Keep CI universal for the reference stack: Go API, Next.js web app, Telegram bot, GraphQL codegen, Docker images, coverage, e2e, and GRACE integrity.
- Make production deployment traceable to a human-readable release tag, not an arbitrary branch SHA.
- Use GitLab Container Registry as the artifact source of truth.
- Deploy with Dokploy Docker Compose stacks using images produced by CI.
- Keep the downstream customization surface small and explicit.

## Non-Goals

- Do not add GitHub Actions or another CI provider in the baseline.
- Do not make Dokploy build from Git for production.
- Do not require a `stage` environment in the template.
- Do not encode project-specific domains, secrets, database credentials, or Dokploy instance details in the repository.
- Do not replace the existing Bun, Nx, Go, Docker, coverage, e2e, or GRACE command surfaces.

## Current Problem

The existing `.gitlab-ci.yml` predates the current GRACE contract:

- It uses Go `1.23` while the repository declares Go `1.25`.
- It does not run the full `bun run verify:coverage` gate before deploy-capable branches.
- It treats e2e as manual and allowed to fail.
- It uses GitLab service containers rather than the repository's isolated `docker/docker-compose.test.yml` contract for full gates.
- It lacks release-tag production semantics.
- It does not publish coverage and Playwright artifacts as first-class evidence.
- It does not model Dokploy deployment.

The new design should replace this with a coherent template pipeline rather than patch the old jobs one by one.

## Recommended Approach

Use one GitLab-first pipeline with stages:

1. `prepare`
2. `validate`
3. `test`
4. `coverage`
5. `build`
6. `release`
7. `deploy`

The pipeline stays readable in `.gitlab-ci.yml`. Helper scripts may be introduced only for repeated safety logic such as asserting that a tag is on `origin/main`, producing image metadata, or calling Dokploy. The CI provider abstraction layer is intentionally out of scope because this repository has chosen GitLab-first.

Rejected alternatives:

- Provider-neutral CI core plus thin GitLab adapter. This is more portable, but adds indirection before a second provider exists.
- Dokploy auto-build from Git. This is simpler to wire, but weakens artifact provenance and makes the deployment host responsible for production builds.
- GitLab plus GitHub Actions in the template. This increases maintenance and drift risk for little immediate value.

## Pipeline Rules

### Merge Requests

Merge requests optimize for fast feedback:

- Install dependencies with `bun install --frozen-lockfile`.
- Resolve the Nx graph with daemon disabled.
- Run schema and codegen drift checks where relevant.
- Run lint, typecheck, affected tests, and affected builds.
- Use the merge request diff base SHA for Nx affected ranges when available; fall back to `origin/main` only for branch pipelines that do not have merge request metadata.
- Do not build deploy images unless an affected build requires Docker validation.
- Do not deploy.

The MR pipeline should use `NX_DAEMON=false`. Previous validation in this repository showed Nx daemon instability after resets, so CI should prefer deterministic graph resolution.

### `develop`

The `develop` branch is the source for the dev environment:

- Run `bun run verify:coverage`.
- Publish coverage and Playwright artifacts.
- Build and push `api`, `web`, and `bot` images.
- Tag images with a dev trace tag such as `develop-$CI_COMMIT_SHORT_SHA`.
- Optionally also update a mutable `dev-latest` tag.
- Automatically trigger Dokploy `dev` compose deployment.

### `main`

The `main` branch is release-candidate validation:

- Run `bun run verify:coverage`.
- Publish coverage and Playwright artifacts.
- Optionally build and push SHA-tagged candidate images for cache/provenance.
- Do not deploy production.
- Do not expose a production manual deploy job from a plain branch pipeline.

### Release Tags

Production release pipelines are created only by SemVer tags that match:

```text
^v[0-9]+\.[0-9]+\.[0-9]+$
```

The release pipeline must:

- Fetch `origin/main`.
- Assert the tag commit is reachable from `origin/main`.
- Run `bun run verify:coverage` by default.
- Build and push release images:
  - `api:vX.Y.Z`
  - `web:vX.Y.Z`
  - `bot:vX.Y.Z`
- Write an image metadata artifact mapping each service to its tag and digest.
- Create a GitLab Release for the tag.
- Expose `deploy:prod` as a manual job.

The default template should rerun the full gate on tag pipelines. A downstream repository may later optimize this to verify an already-passed `main` pipeline for the same SHA, but the template should favor self-contained release proof.

## Docker Images

The baseline builds one image per deployable service:

- `api` from `docker/api.Dockerfile`
- `web` from `docker/web.Dockerfile`
- `bot` from `docker/bot.Dockerfile`

GitLab Container Registry is the source of truth:

```text
$CI_REGISTRY_IMAGE/api:<tag>
$CI_REGISTRY_IMAGE/web:<tag>
$CI_REGISTRY_IMAGE/bot:<tag>
```

Release tags are immutable by policy. Jobs should avoid overwriting an existing `vX.Y.Z` image unless the release is explicitly rebuilt through a controlled maintenance procedure outside the normal pipeline.

The image metadata artifact should record:

- service name
- Dockerfile path
- image tag
- image digest
- commit SHA
- pipeline ID
- release tag when present

## Dokploy Deployment Contract

Dokploy deployment is modeled as one Docker Compose stack per environment:

- `dev` stack for `develop`
- `prod` stack for SemVer release tags

The repository should provide a template Compose file for Dokploy that references environment variables rather than hard-coded image names:

```yaml
services:
  api:
    image: ${API_IMAGE}
    pull_policy: always
  web:
    image: ${WEB_IMAGE}
    pull_policy: always
  bot:
    image: ${BOT_IMAGE}
    pull_policy: always
```

GitLab deploy jobs compute the image variables for the target environment:

```text
IMAGE_TAG=v1.2.3
API_IMAGE=$CI_REGISTRY_IMAGE/api:$IMAGE_TAG
WEB_IMAGE=$CI_REGISTRY_IMAGE/web:$IMAGE_TAG
BOT_IMAGE=$CI_REGISTRY_IMAGE/bot:$IMAGE_TAG
```

Real domains, secrets, database URLs, bot tokens, TLS settings, and Dokploy IDs remain outside the template and are configured as protected GitLab variables or Dokploy environment values.

For immutable release deployments, the CI job must update the Dokploy compose environment or compose file before triggering deployment. Dokploy's compose deploy endpoint starts a deployment for the current compose configuration; it does not itself carry a new `IMAGE_TAG` payload. The default CI helper should therefore perform:

1. `compose.update` with the non-secret image variables or rendered compose file for the target release.
2. `compose.deploy` for the same compose ID.

The deploy job should trigger Dokploy through one of two supported mechanisms:

- A protected Dokploy webhook URL for the target compose stack.
- The Dokploy API, using a protected API key and compose identifier.

The implementation should choose the API path as the default because it can update the compose environment or file and then deploy the compose stack with clearer status handling. Webhook support can remain a documented fallback only when the Dokploy stack already resolves the intended image tag without a CI-side update.

## Environment Policy

### Dev

- Source: `develop`
- Gate: full `bun run verify:coverage`
- Image tag: `develop-$CI_COMMIT_SHORT_SHA`, optionally `dev-latest`
- Deploy: automatic after successful image push
- Dokploy variables: environment-scoped dev variables, protected when `develop` is a protected branch

### Prod

- Source: SemVer tag from `origin/main`
- Gate: full `bun run verify:coverage`
- Image tag: `vX.Y.Z`
- Release: GitLab Release created before deployment approval
- Deploy: manual `deploy:prod`
- GitLab environment: protected production environment
- Dokploy variables: protected and masked CI variables

## Required CI Variables

The template should document these variables and keep real values out of git:

- `DOKPLOY_DEV_URL` or `DOKPLOY_DEV_WEBHOOK_URL`
- `DOKPLOY_DEV_API_KEY`
- `DOKPLOY_DEV_COMPOSE_ID`
- `DOKPLOY_PROD_URL` or `DOKPLOY_PROD_WEBHOOK_URL`
- `DOKPLOY_PROD_API_KEY`
- `DOKPLOY_PROD_COMPOSE_ID`
- Any registry credentials needed by Dokploy to pull from GitLab Container Registry

GitLab's built-in `$CI_REGISTRY`, `$CI_REGISTRY_USER`, `$CI_REGISTRY_PASSWORD`, and `$CI_REGISTRY_IMAGE` should be used for CI-side push authentication.

## Artifacts and Evidence

Full gate pipelines must retain:

- `dist/coverage`
- `dist/test-results/web-e2e`
- `dist/playwright-report/web`
- XML and GRACE validation logs or summaries

Release pipelines must additionally retain:

- image metadata artifact
- GitLab Release link
- Dokploy deployment trigger response summary

Artifacts should be available for merge/release review. Retention duration can be tuned per downstream repository, but the template should set a reasonable default instead of dropping evidence immediately.

## Failure Rules

- MR pipelines must fail on lint, schema, typecheck, affected test, or affected build failures.
- `develop`, `main`, and release tag pipelines must fail if `bun run verify:coverage` fails.
- Release tag pipelines must fail if the tag is not reachable from `origin/main`.
- Production deploy jobs must not exist for branch pipelines.
- Production deploy jobs must not run automatically.
- Deploy jobs must fail if required Dokploy variables are missing.
- Deploy jobs must fail if they cannot update the target Dokploy compose environment or compose file to the intended image tag.
- Deploy jobs must not print Dokploy API keys, webhook URLs, registry passwords, application secrets, or high-risk payloads.
- A failed Dokploy deploy trigger is a failed CI job.
- Existing release tags must not be silently overwritten.

## Downstream Customization Surface

New repositories derived from this template should usually change only:

- Dokploy instance URLs, API keys, compose IDs, and environment variables.
- Domain names and runtime secrets in Dokploy.
- Service list if the project removes `bot` or adds workers.
- Optional stage environment jobs if the project wants `stage`.
- Retention duration for artifacts.
- Release versioning policy if the project needs prerelease tags such as `v1.2.3-rc.1`.

Everything else should work from the template baseline.

## Documentation and GRACE Updates for Implementation

The implementation plan should update:

- `.gitlab-ci.yml`
- Dokploy compose template or deployment documentation
- CI helper scripts if introduced
- `README.md`
- `docs/requirements.xml`
- `docs/technology.xml`
- `docs/development-plan.xml`
- `docs/knowledge-graph.xml`
- `docs/verification-plan.xml`
- `docs/operational-packets.xml` if CI worker packets need a new checkpoint shape

The GRACE module model should add or update a CI/CD module rather than treating deployment as an incidental workspace script.

## References

- Dokploy Auto Deploy: https://docs.dokploy.com/docs/core/auto-deploy
- Dokploy GitLab integration: https://docs.dokploy.com/docs/core/gitlab
- Dokploy Docker Compose: https://docs.dokploy.com/docs/core/docker-compose
- Dokploy Compose API: https://docs.dokploy.com/docs/api/compose
- Dokploy Registry: https://docs.dokploy.com/docs/core/registry
