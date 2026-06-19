<!-- FILE: docs/infrastructure/gitlab-runners.md -->
<!-- VERSION: 1.0.0 -->
<!-- START_MODULE_CONTRACT -->
<!--   PURPOSE: Define the department GitLab Runner pool contract for repositories derived from this template. -->
<!--   SCOPE: Runner classes, tags, concurrency, Docker-in-Docker risk acceptance, host sizing, and GitLab group settings; excludes installation tokens and project-specific secrets. -->
<!--   DEPENDS: .gitlab-ci.yml, docs/infrastructure/ci-cd.md, docs/technology.xml, docs/development-plan.xml, docs/knowledge-graph.xml, docs/verification-plan.xml. -->
<!--   LINKS: M-CI-CD / V-M-CI-CD / DF-CI-CD-RELEASE. -->
<!--   ROLE: DOC -->
<!--   MAP_MODE: SUMMARY -->
<!-- END_MODULE_CONTRACT -->
<!-- START_MODULE_MAP -->
<!--   Runner Pool Goal - Explains why downstream repositories use a shared department runner pool. -->
<!--   Practical Risk Decision - Records the accepted Docker-in-Docker strategy and required guardrails. -->
<!--   Runner Classes - Defines tags, limits, privilege mode, and intended job ownership. -->
<!--   Template CI Contract - Maps monorepo-template jobs to runner tags. -->
<!--   Host Baseline - Records the initial sizing policy for the first runner host. -->
<!--   Operations Checklist - Lists GitLab group settings, cleanup, cache, and scaling rules. -->
<!-- END_MODULE_MAP -->
<!-- START_CHANGE_SUMMARY -->
<!--   LAST_CHANGE: 1.0.0 - Added the department runner pool operating contract for template-derived projects. -->
<!-- END_CHANGE_SUMMARY -->

# GitLab Runner Pool

## Runner Pool Goal

Repositories derived from this template should use a shared department runner pool, registered at the top-level GitLab group that owns the department projects. This makes the runner setup reusable for all downstream repositories while keeping access scoped to the department group rather than the whole GitLab instance.

The pool is part of the template contract:

- `.gitlab-ci.yml` defines which jobs exist and which tags they request.
- This document defines which runner classes must satisfy those tags.
- `docs/infrastructure/ci-cd.md` defines release, image, and Dokploy deployment behavior.

GitLab group runners are the right default because GitLab makes them available to all projects in a group and its subgroups. Project runners should be reserved for special one-off workloads, customer-specific isolation, or repositories with unusual compliance needs.

## Practical Risk Decision

The practical default is to use Docker executor runners and keep Docker-in-Docker only on separate protected runner classes. This is not zero-risk. Privileged Docker-in-Docker disables important container isolation controls and can expose the host to container breakout if untrusted jobs run there.

For this department template, the accepted risk is reasonable when all of these guardrails hold:

- Merge request jobs run on an unprivileged runner.
- Docker-in-Docker jobs run only on protected branches and protected release tags.
- Full-gate and image-build runners have `limit = 1`.
- Protected variables are exposed only to protected refs.
- Runner tags are required in every template job.
- The runner host is not used for application runtime, databases, Dokploy, or unrelated services.
- No project enables untrusted external fork pipelines on the protected Docker-in-Docker runner classes.

If any of these guardrails are not true, treat the risk as high and move image builds to a stronger isolation model such as rootless BuildKit, Kaniko-like builds, autoscaled disposable runners, or Kubernetes runners.

## Runner Classes

Start with these group runners:

| Runner                 | Tags                                      | Protected | Privileged | Limit | Purpose                                                                                                                              |
| ---------------------- | ----------------------------------------- | --------- | ---------- | ----- | ------------------------------------------------------------------------------------------------------------------------------------ |
| `dept-template-fast`   | `dept`, `template`, `fast`, `docker`      | No        | No         | 3     | Merge request lint, typecheck, affected tests, affected builds, XML and GRACE validation.                                            |
| `dept-template-full`   | `dept`, `template`, `full`, `dind`        | Yes       | Yes        | 1     | Full `bun run verify:coverage`, Playwright e2e, and Docker Compose based test infrastructure on `develop`, `main`, and release tags. |
| `dept-template-build`  | `dept`, `template`, `image-build`, `dind` | Yes       | Yes        | 1     | Docker Buildx image builds and pushes to GitLab Container Registry.                                                                  |
| `dept-template-deploy` | `dept`, `template`, `deploy`              | Yes       | No         | 1     | Dokploy API calls and release deployment jobs that receive protected deployment variables.                                           |

The deploy runner is cheap but useful because it keeps Dokploy API keys away from privileged Docker-in-Docker jobs. If the first rollout must stay smaller, deploy jobs may temporarily run on `dept-template-full`, but the preferred target is a separate unprivileged protected deploy runner.

## Template CI Contract

The template pipeline should map jobs to runner classes explicitly:

```yaml
.toolchain:
  tags: [dept, template, fast]

verify:full:
  tags: [dept, template, full]

.docker_build:
  tags: [dept, template, image-build]

deploy:dev:
  tags: [dept, template, deploy]

deploy:prod:
  tags: [dept, template, deploy]
```

The current pipeline already disables the Nx daemon and skips Nx cache for deterministic CI graph resolution:

```yaml
variables:
  NX_DAEMON: 'false'
  NX_SKIP_NX_CACHE: '1'
```

Keep that behavior in downstream repositories unless a future remote cache policy is designed and verified separately.

## Toolchain Image

Do not install the full toolchain in every job. Build and publish a department CI image instead:

```text
registry.gitlab.com/<department>/ci-images/template-toolchain:node22-bun-go1.25-playwright-grace
```

The image should include:

- Node.js 22
- Bun
- Go 1.25
- Docker CLI and Docker Compose plugin
- Playwright browser dependencies matching the template Playwright version
- `xmllint`
- GRACE CLI
- Git and CA certificates

The image should not include project secrets, deployment credentials, or a Docker daemon. Docker daemon access comes only from the dedicated Docker-in-Docker service used by the protected full and image-build runners.

Use immutable or deliberately versioned image tags. Avoid `latest` for the default template image because silent toolchain drift creates noisy generated diffs and hard-to-reproduce CI failures.

## Host Baseline

The first runner host baseline is:

| Resource | Baseline                           |
| -------- | ---------------------------------- |
| CPU      | 8 vCPU                             |
| RAM      | 15 GiB                             |
| Disk     | 100 GiB root disk                  |
| OS       | Ubuntu LTS                         |
| Runtime  | Docker executor, cgroup v2 capable |

Use this initial capacity:

```text
global concurrent = 4
dept-template-fast limit = 3
dept-template-full limit = 1
dept-template-build limit = 1
dept-template-deploy limit = 1
```

The global concurrency limit is the real host safety valve. Keep it at `4` for the first rollout. The per-runner limits describe what each class is allowed to do, but the global limit prevents the host from running every class at full limit at the same time.

If queues become painful, add another runner host or increase disk/RAM before raising concurrency. Increasing concurrency first usually turns CI into a memory, Docker layer, or disk pressure problem.

## Runner Configuration Defaults

Use separate runner entries rather than one catch-all runner. A representative `config.toml` shape:

```toml
concurrent = 4
check_interval = 0

[[runners]]
  name = "dept-template-fast-01"
  executor = "docker"
  limit = 3
  builds_dir = "/srv/gitlab-runner/builds"
  cache_dir = "/srv/gitlab-runner/cache"
  [runners.docker]
    image = "registry.gitlab.com/<department>/ci-images/template-toolchain:node22-bun-go1.25-playwright-grace"
    privileged = false
    cpus = "2"
    memory = "3g"
    pull_policy = "always"

[[runners]]
  name = "dept-template-full-01"
  executor = "docker"
  limit = 1
  builds_dir = "/srv/gitlab-runner/builds"
  cache_dir = "/srv/gitlab-runner/cache"
  [runners.docker]
    image = "registry.gitlab.com/<department>/ci-images/template-toolchain:node22-bun-go1.25-playwright-grace"
    privileged = true
    cpus = "4"
    memory = "8g"
    pull_policy = "always"
    allowed_services = ["docker:*dind"]

[[runners]]
  name = "dept-template-build-01"
  executor = "docker"
  limit = 1
  builds_dir = "/srv/gitlab-runner/builds"
  cache_dir = "/srv/gitlab-runner/cache"
  [runners.docker]
    image = "docker:27"
    privileged = true
    cpus = "4"
    memory = "8g"
    pull_policy = "always"
    allowed_services = ["docker:*dind"]

[[runners]]
  name = "dept-template-deploy-01"
  executor = "docker"
  limit = 1
  builds_dir = "/srv/gitlab-runner/builds"
  cache_dir = "/srv/gitlab-runner/cache"
  [runners.docker]
    image = "registry.gitlab.com/<department>/ci-images/template-toolchain:node22-bun-go1.25-playwright-grace"
    privileged = false
    cpus = "1"
    memory = "1g"
    pull_policy = "always"
```

Register tags and protected-runner flags in GitLab during runner creation. Do not commit registration or authentication tokens to the repository.

## GitLab Group Settings

At the department group level:

- Create the runners as group runners using runner authentication tokens.
- Require jobs to use tags.
- Protect `develop`, `main`, and release tags matching `v*`.
- Mark `dept-template-full`, `dept-template-build`, and `dept-template-deploy` as protected runners.
- Disable instance runners for template-derived projects unless the project explicitly needs fallback capacity.
- Keep Dokploy, registry, and production variables protected and masked.
- Do not expose protected variables to merge request pipelines from untrusted forks.

At each downstream project:

- Keep the template runner tags unless the project has a documented exception.
- Keep `NX_DAEMON=false` in CI.
- Keep full gate and image build jobs off merge request pipelines.
- Document any additional runner tag in the project infrastructure docs.

## Cache and Disk Policy

The first host has only 100 GiB of disk. Docker layers, Playwright artifacts, coverage output, and build caches will fill it faster than source checkouts.

Use these defaults:

- Keep GitLab cache project-scoped, not global across all repositories.
- Do not share `node_modules` or `.bun` caches through host paths between projects.
- Prefer GitLab cache keys that include the lockfile hash.
- Keep artifacts in GitLab, not in long-lived host directories.
- Run scheduled Docker cleanup for stopped containers, old build cache, and unused images.
- Alert before disk reaches 80 percent.

Example cleanup policy:

```bash
docker system prune --all --force --filter "until=168h"
docker buildx prune --force --filter "until=168h"
```

Do not run cleanup while pipelines are active unless the host is already in emergency disk pressure.

## Operating Checklist

Before enabling the pool for the department:

- Password SSH login is disabled on the runner host.
- Root SSH login is key-only or disabled in favor of a sudo user.
- The host firewall allows SSH and outbound GitLab, registry, package mirror, and Dokploy API traffic.
- Docker and GitLab Runner are installed from a documented package source.
- `/srv/gitlab-runner/builds` and `/srv/gitlab-runner/cache` exist.
- Runner host monitoring covers CPU load, memory, disk, Docker daemon health, and failed jobs.
- A simple "runner smoke" pipeline passes on each runner class.
- The template `.gitlab-ci.yml` uses explicit tags for fast, full, image-build, and deploy jobs.

## Scaling Rule

Scale by adding runner hosts with the same classes and a numbered suffix:

```text
dept-template-fast-02
dept-template-full-02
dept-template-build-02
dept-template-deploy-02
```

Do not turn the first host into an oversized mixed workload machine. CI runner reliability comes from predictable isolation, boring limits, and repeatable images more than from squeezing every CPU core.

## References

- GitLab group runners: https://docs.gitlab.com/ci/runners/runners_scope/
- GitLab Runner security: https://docs.gitlab.com/runner/security/
- GitLab Docker executor: https://docs.gitlab.com/runner/executors/docker/
- GitLab Runner advanced configuration: https://docs.gitlab.com/runner/configuration/advanced-configuration/
