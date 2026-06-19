# FILE: docker/web.Dockerfile
# VERSION: 1.0.1
# START_MODULE_CONTRACT
#   PURPOSE: Build the deployable public Next web image for development and production runtime.
#   SCOPE: Owns apps/web container build stages only; excludes admin SPA deployment and API/Bot images.
#   DEPENDS: oven/bun, node:22-alpine, apps/web, apps/web-admin/package.json, libs, tools.
#   LINKS: M-CI-CD / M-WEB / V-M-CI-CD.
#   ROLE: CONFIG
#   MAP_MODE: SUMMARY
# END_MODULE_CONTRACT
# START_MODULE_MAP
#   dev - Runs the public Next web dev server on port 3000 for local Docker compose.
#   builder - Builds apps/web with Next standalone output.
#   prod - Runs the apps/web standalone server for Dokploy and CI images.
# END_MODULE_MAP
# START_CHANGE_SUMMARY
#   LAST_CHANGE: 1.0.1 - Preserve Next standalone root layout so Bun node_modules symlinks resolve at runtime.
# END_CHANGE_SUMMARY

# ---- Dev stage: deployable web image is the Next.js public web app ----
FROM oven/bun:1.3.5-alpine AS dev

WORKDIR /app
COPY package.json bun.lock tsconfig.base.json .eslintrc.json ./
COPY apps/web/package.json apps/web/
COPY apps/web-admin/package.json apps/web-admin/
COPY libs/ libs/
COPY tools/ tools/
RUN bun install --frozen-lockfile --ignore-scripts
COPY apps/web/ apps/web/
WORKDIR /app/apps/web
CMD ["bun", "exec", "next", "dev", "--hostname", "0.0.0.0", "--port", "3000"]

# ---- Build stage: deployable web image is the Next.js public web app ----
FROM oven/bun:1.3.5-alpine AS builder

WORKDIR /app
COPY package.json bun.lock tsconfig.base.json .eslintrc.json ./
COPY apps/web/package.json apps/web/
COPY apps/web-admin/package.json apps/web-admin/
COPY libs/ libs/
COPY tools/ tools/
RUN bun install --frozen-lockfile --ignore-scripts
COPY apps/web/ apps/web/
WORKDIR /app/apps/web
RUN bun run build

# ---- Production stage ----
FROM node:22-alpine AS prod

WORKDIR /app
COPY --from=builder /app/apps/web/.next/standalone ./
WORKDIR /app/apps/web
COPY --from=builder /app/apps/web/.next/static ./.next/static
COPY --from=builder /app/apps/web/public ./public
EXPOSE 3000
CMD ["node", "server.js"]
