# CI/CD Review

## Scope

- GitLab CI/CD pipeline
- CI helper tooling
- Dokploy compose deployment contract
- GRACE CI/CD module and verification updates

## Commands

- `bunx vitest run --config tools/vitest.config.ts tools/ci/src`: PASS, 3 files and 30 tests.
- `bunx nx run ci-tools:test`: PASS, CI helper coverage 100 percent statements, branches, functions, and lines.
- `bunx vitest run --config tools/vitest.config.ts --coverage`: PASS, workspace tooling coverage 100 percent statements, branches, functions, and lines.
- `bunx prettier --check .gitlab-ci.yml`: PASS.
- `bun tools/ci/src/cli.ts affected-base`: PASS, printed `origin/main`.
- `env CI_COMMIT_TAG=latest bun tools/ci/src/cli.ts assert-release-tag`: EXPECTED FAIL, printed `Release tag must match vX.Y.Z: latest`.
- `xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml`: PASS.
- `grace lint --path .`: PASS, no GRACE integrity issues.
- `rg -n "apk add --no-cache jq curl bash|docker manifest inspect|Release image already exists|export API_IMAGE_DIGEST|bun tools/ci/src/cli.ts write-image-metadata|deploy:prod|CI_COMMIT_BRANCH == \"main\"|CI_COMMIT_TAG|compose.update|compose.deploy|NX_DAEMON|CI_MERGE_REQUEST_DIFF_BASE_SHA" .gitlab-ci.yml tools/ci docs README.md deploy/dokploy`: PASS.
- `bun run verify:coverage`: first run reached e2e and failed because unrelated `/home/nolood/general/auto/.worktrees/auto-eylq/web` Next server occupied `localhost:13000`.
- `env E2E_WEB_PORT=13001 E2E_API_PORT=18081 bun run verify:coverage`: PASS. Lint, codegen, web typecheck, build, coverage, e2e, XML, and GRACE gates completed. Playwright e2e passed 4/4.

## Findings

- Production deploy is tag-only and manual.
- `main` has no production deploy job.
- Dokploy compose update precedes compose deploy.
- Full gate evidence is retained for deploy-capable pipelines.
- CI helper tests cover affected-base resolution, strict release tags, image refs, image metadata, env requirements, redaction, Dokploy env merge, update-before-deploy order, and CLI command dispatch.

## Residual Risks

- First real GitLab run must verify runner package availability for Docker Buildx, Docker Compose, Go 1.25, Bun, Playwright, xmllint, and GRACE CLI.
- Real Dokploy API variables and compose IDs are environment-specific and must be configured outside git.
- Local default e2e web port `13000` can conflict with unrelated dev servers; the Playwright config supports `E2E_WEB_PORT` and `E2E_API_PORT` overrides and the successful verification used `13001` and `18081`.
