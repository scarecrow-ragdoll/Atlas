<!-- FILE: docs/infrastructure/ci-cd.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Document CI/CD, image, release, and Dokploy runtime variable contracts. -->
<!--   SCOPE: Covers operator-facing pipeline and Dokploy configuration guidance; excludes executable CI helper implementation. -->
<!--   DEPENDS: .gitlab-ci.yml, deploy/dokploy/docker-compose.template.yml, tools/ci. -->
<!--   LINKS: M-CI-CD / V-M-CI-CD. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Pipeline Model - Documents branch, tag, image, and deploy behavior. -->
<!--   Required GitLab Variables - Lists environment-scoped Dokploy and web runtime variables. -->
<!--   Dokploy Compose Stack - Describes image env mutation and runtime placeholders. -->
<!--   Production Release Flow - Defines the production deploy sequence. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.1 - Replaced bearer placeholder env docs with web-admin bootstrap and session env. -->
<!-- END_CHANGE_SUMMARY -->

# CI/CD and Dokploy

Runner capacity and security are part of the CI/CD contract. See [GitLab Runner Pool](gitlab-runners.md) for the department group runner tags, Docker-in-Docker policy, concurrency limits, and host operations checklist.

## Pipeline Model

- Merge requests run fast affected checks.
- `develop` runs `bun run verify:coverage`, builds `api`, the public Next `web` image, and `bot` images, pushes them to GitLab Container Registry, and deploys automatically to Dokploy dev.
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
- `WEB_API_BASE_URL`

Production variables must be protected and masked.

## Dokploy Compose Stack

Use `deploy/dokploy/docker-compose.template.yml` as the base compose file in Dokploy. The GitLab deploy job updates the image variables before deployment:

- `IMAGE_TAG`
- `API_IMAGE`
- `WEB_IMAGE` for the public Next `web` Docker image.
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
- `ADMIN_INITIAL_EMAIL`, required while `admin_users` is empty
- `ADMIN_INITIAL_PASSWORD`, required while `admin_users` is empty
- `ADMIN_INITIAL_NAME`, required while `admin_users` is empty
- `ADMIN_ORIGINS`
- `ADMIN_SESSION_COOKIE_NAME`
- `ADMIN_SESSION_TTL`
- `ADMIN_SESSION_COOKIE_SECURE`
- `ADMIN_SESSION_SAME_SITE`
- `ADMIN_SESSION_KEY_SECRET`, always required for API startup
- `WEB_API_BASE_URL` for the public Next server-side REST proxy, such as `https://api.example.com`
- `BOT_TOKEN`

The compose template sets the service-local runtime defaults `SERVER_PORT=8080`, `SERVER_ENV=production`, and `PORT=3000`; override `SERVER_ENV` in Dokploy only when the target environment needs a different value.

For a fresh local or Dokploy API start with an empty `admin_users` table, provide `ADMIN_INITIAL_EMAIL`, `ADMIN_INITIAL_PASSWORD`, `ADMIN_INITIAL_NAME`, and `ADMIN_SESSION_KEY_SECRET`. After the first admin exists, `ADMIN_INITIAL_*` may be unset, but `ADMIN_SESSION_KEY_SECRET` remains required.

## Production Release Flow

1. Merge into `main`.
2. Wait for the `main` release-candidate pipeline to pass.
3. Create a SemVer tag from `origin/main`:

   ```bash
   git fetch origin main --tags
   git tag v1.2.3 origin/main
   git push origin v1.2.3
   ```

4. Confirm the tag is reachable from `origin/main`.
5. Wait for the release pipeline to pass.
6. Approve `deploy:prod` in GitLab.

Production deployment from a branch pipeline is intentionally unavailable.
