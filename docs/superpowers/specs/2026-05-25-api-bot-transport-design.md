# API Bot Transport Design

**Status:** Approved
**Date:** 2026-05-25

## Goal

Add a template-friendly communication baseline between the Go API and Telegram bot without turning the monorepo template into a product-specific platform.

The baseline demonstrates two reusable integration patterns:

1. Synchronous service-to-service calls through an internal gRPC contract.
2. Asynchronous background work through an Asynq task queue backed by Redis.

Both examples use a neutral demo domain rather than the existing users CRUD domain. Downstream projects can keep, replace, or delete the demo transport slice without untangling it from product business logic.

## Current Context

The repository currently has separate `api` and `bot` Go applications. The API owns public HTTP surfaces for health, readiness, REST users, GraphQL users, PostgreSQL, and Redis. The bot owns Telegram polling, command handlers, logging middleware, recovery middleware, and config.

GRACE currently models `M-API` and `M-BOT` as separate entry points. `M-BOT` depends only on shared config and logger libraries, while `M-WEB` is the public REST consumer of `M-API`. This change introduces an explicit internal transport slice between `M-BOT` and `M-API` while keeping users CRUD independent.

## Key Decisions

| Decision               | Choice                                        | Rationale                                                                                                                 |
| ---------------------- | --------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| Sync transport         | gRPC                                          | Shows a typed internal service-to-service contract for downstream products.                                               |
| Async transport        | Asynq over Redis                              | Gives a real task queue with workers, retries, queue selection, and concurrency instead of low-level Redis Streams usage. |
| Demo domain            | Neutral demo service/tasks                    | Keeps the template reusable and avoids coupling Telegram to users CRUD.                                                   |
| Worker ownership       | Separate `apps/worker` app                    | Background jobs are not inherently bot-owned, and separating polling from workers keeps process responsibilities clear.   |
| Auth for internal gRPC | Shared token interceptor                      | Internal-only networking is not enough for a production-oriented template; the token is simple and removable.             |
| Asynq exposure         | Thin queue adapter and task contract packages | Asynq is useful but still a `v0.x` API, so direct dependency spread should be limited.                                    |

Asynq is intentionally chosen as a task queue, not as a general message bus. It is a good default for background jobs such as notifications, email, sync, exports, media processing, and webhook delivery. If a downstream product needs event streaming, cross-language pub/sub, or long-lived event replay, it should introduce a dedicated broker or stream module rather than stretching this demo queue.

## Architecture

The transport slice adds four public concepts:

1. Demo gRPC contract.
2. API gRPC server adapter.
3. Bot gRPC client adapter.
4. Asynq task queue contract and worker.

The API runs its existing HTTP server and a new internal gRPC listener in the same process. The bot keeps its Telegram polling loop and calls the gRPC client from a demo command. A new worker process consumes queue tasks from Redis through Asynq.

The sync and async demos are related but not hidden inside each other:

- `DemoService.GetGreeting` is the synchronous request/response example used by the bot `/demo` command.
- `DemoService.EnqueueEcho` is the explicit asynchronous producer example. It enqueues one `demo.echo` task and returns the Asynq task ID and queue name.

`GetGreeting` must not enqueue background work implicitly. This keeps the Telegram command deterministic and makes the queue trigger obvious in tests and downstream copies.

The default runtime shape becomes:

```text
Telegram user
  -> apps/bot /demo command
  -> bot internal DemoClient interface
  -> gRPC client
  -> apps/api internal gRPC listener
  -> API demo service
  -> gRPC response
  -> Telegram reply

DemoService.EnqueueEcho
  -> queue adapter
  -> Asynq client
  -> Redis
  -> apps/worker Asynq server
  -> demo task handler
```

## GRACE Module Mapping

The implementation plan should add or update these GRACE module IDs and verification refs:

| Module ID            | Verification ref       | Path                                                                           | Role                                                                                        |
| -------------------- | ---------------------- | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------- |
| `M-DEMOAPI`          | `V-M-DEMOAPI`          | `libs/go/demoapi`                                                              | Raw proto, generated Go gRPC code, and generated contract package for the neutral demo API. |
| `M-API-GRPC`         | `V-M-API-GRPC`         | `apps/api/internal/grpc`                                                       | API-side internal gRPC listener, auth/logging interceptors, and demo service registration.  |
| `M-API-DEMO-SERVICE` | `V-M-API-DEMO-SERVICE` | `apps/api/internal/service/demo`                                               | Transport-neutral demo service used by gRPC handlers.                                       |
| `M-BOT-API-CLIENT`   | `V-M-BOT-API-CLIENT`   | `apps/bot/internal/apiclient`                                                  | Bot-side gRPC client adapter and error mapping.                                             |
| `M-GO-TASKS`         | `V-M-GO-TASKS`         | `libs/go/tasks`                                                                | Asynq task type constants, payloads, factory helpers, and parse helpers.                    |
| `M-GO-QUEUE`         | `V-M-GO-QUEUE`         | `libs/go/queue`                                                                | Thin Asynq client/server configuration adapter.                                             |
| `M-WORKER`           | `V-M-WORKER`           | `apps/worker`                                                                  | Background worker app that processes Asynq tasks.                                           |
| `M-CI-CD`            | `V-M-CI-CD`            | `.gitlab-ci.yml`, `tools/ci`, `deploy/dokploy`, `docs/infrastructure/ci-cd.md` | Existing deployment module extended with worker image and release metadata.                 |

Existing `M-API` remains the public HTTP API entry point and gains cross-links to `M-API-GRPC`, `M-API-DEMO-SERVICE`, `M-DEMOAPI`, `M-GO-QUEUE`, and `M-GO-TASKS`. Existing `M-BOT` remains the Telegram polling entry point and gains a cross-link to `M-BOT-API-CLIENT`. Existing `M-COVERAGE-GATE` gains cross-links to all new modules.

## Modules

### `libs/go/demoapi`

Owns the proto contract and generated Go package for the neutral demo API. This must be a Go module with:

- `libs/go/demoapi/go.mod`
- `libs/go/demoapi/proto/demo/v1/demo.proto`
- `libs/go/demoapi/gen/demo/v1/*.pb.go`
- `libs/go/demoapi/gen/demo/v1/*_grpc.pb.go`

The module path should be `monorepo-template/libs/go/demoapi`. Add `./libs/go/demoapi` to `go.work`, and import generated code from `monorepo-template/libs/go/demoapi/gen/demo/v1`.

The contract should stay tiny:

- `DemoService.GetGreeting`
- `DemoService.EnqueueEcho`
- request fields such as `name`
- enqueue request fields such as `message`
- response fields such as `message`, `request_id`, `task_id`, and `queue`

Generated Go code is committed. Codegen drift is checked by `bunx nx run go-demoapi:codegen`, and root `bun run codegen` must include `go-demoapi` along with the existing API and web-admin codegen targets.

### `apps/api/internal/grpc`

Owns the API-side gRPC server setup and handlers. It should:

- register generated demo service server implementations;
- apply internal auth and request logging interceptors;
- map domain/service errors to stable gRPC status codes;
- avoid placing business logic directly in generated handlers.

### `apps/api/internal/service/demo`

Owns neutral demo behavior. `GetGreeting` returns a greeting without side effects. `EnqueueEcho` validates a message and enqueues one `demo.echo` task through an interface. This keeps gRPC handlers thin and mirrors the existing transport-neutral `UserService` pattern.

### `apps/bot/internal/apiclient`

Owns the bot-side gRPC client adapter. Bot handlers depend on a local interface, not generated client types, so command tests can use fakes without opening sockets.

The client should set per-call timeouts and convert transport failures into simple typed errors such as `ErrUnavailable`.

### `libs/go/tasks`

Owns Asynq task type constants, payload structs, factory functions, and parse helpers.

Example surface:

- `TypeDemoEcho`
- `DemoEchoPayload`
- `NewDemoEchoTask(payload)`
- `ParseDemoEchoTask(task)`

Task payloads are JSON by default because they are easy to inspect and stable enough for template examples. Payload structs should be versionable through additive fields.

### `libs/go/queue`

Owns a small adapter around Asynq client/server configuration. This package should hide Asynq setup from API and worker entrypoints where practical.

The adapter should cover:

- Redis connection options from shared config;
- default queue names and weights;
- enqueue options such as max retry, timeout, and queue;
- graceful close of clients.

### `apps/worker`

Owns the background task worker process. It starts an Asynq server, registers demo task handlers, logs startup/shutdown, and exits on fatal Redis or handler registration errors.

This is a separate Go application rather than a bot command because background jobs are a platform concern, not a Telegram concern.

The worker should have its own Go module, Nx project, and Dockerfile:

- `apps/worker/go.mod`
- `apps/worker/project.json`
- `docker/worker.Dockerfile`

Add `./apps/worker` to `go.work`.

## Configuration

### API

Add API config for internal gRPC:

- `grpc.enabled`
- `grpc.port`
- `grpc.auth_token`
- `grpc.read_timeout` or equivalent server defaults if used by the selected library

Environment bindings:

- `GRPC_ENABLED` -> `grpc.enabled`
- `GRPC_PORT` -> `grpc.port`
- `INTERNAL_GRPC_TOKEN` -> `grpc.auth_token`

When `grpc.enabled=true`, `grpc.port` and `grpc.auth_token` are required.

Add queue config through the existing Redis config plus small queue settings:

- `queue.enabled`
- `queue.default_queue`
- enqueue retry/timeout defaults if not hardcoded in task factories

Environment bindings:

- `QUEUE_ENABLED` -> `queue.enabled`
- `QUEUE_DEFAULT_QUEUE` -> `queue.default_queue`

### Bot

Add bot config for internal API access:

- `internal_api.enabled`
- `internal_api.grpc_address`
- `internal_api.auth_token`
- `internal_api.timeout`

Environment bindings:

- `INTERNAL_API_ENABLED` -> `internal_api.enabled`
- `INTERNAL_API_GRPC_ADDRESS` -> `internal_api.grpc_address`
- `INTERNAL_GRPC_TOKEN` -> `internal_api.auth_token`
- `INTERNAL_API_TIMEOUT` -> `internal_api.timeout`

When `internal_api.enabled=true`, `internal_api.grpc_address`, `internal_api.auth_token`, and `internal_api.timeout` are required.

The Docker and Dokploy defaults should point to the API service name, for example `api:9090`. Local config can default to `localhost:9090`.

### Worker

Add worker config for Redis and queue processing:

- `redis.host`
- `redis.port`
- `redis.password`
- `redis.db`
- `queue.concurrency`
- `queue.queues`
- `log.level`
- `log.format`

Environment bindings:

- `WORKER_CONCURRENCY` -> `queue.concurrency`

Queue weights stay YAML-only in v1. Do not add a `WORKER_QUEUES` env parser unless a downstream product needs runtime queue-weight overrides.

The worker should use the same shared config and logger libraries as API and bot.

## Runtime Flows

### Sync gRPC Flow

1. A Telegram user sends `/demo`.
2. The bot handler reads request context and calls `DemoClient.GetGreeting`.
3. The gRPC client sends metadata containing the internal auth token.
4. The API gRPC listener validates the token and invokes the demo service.
5. The demo service returns a greeting.
6. The bot sends the greeting to Telegram.

If the gRPC call times out or returns unavailable, the bot sends a short service-unavailable response and logs the transport error without crashing the polling loop.

### Async Queue Flow

1. A client calls `DemoService.EnqueueEcho` over internal gRPC.
2. The API queue adapter enqueues the task through Asynq.
3. Redis stores the task.
4. The API returns the Asynq task ID and queue name to the caller.
5. `apps/worker` pulls the task through Asynq and runs the demo handler.
6. The handler validates payload and logs the processed demo event.
7. Handler errors are returned to Asynq for retry unless the payload is invalid and should skip retry.

The async example is intentionally independent of Telegram. Downstream projects can use it for bot notifications later, but the template should not imply that every queued task belongs to `apps/bot`.

## Error Handling

### gRPC

- Missing or invalid auth token returns `codes.Unauthenticated`.
- Invalid demo input returns `codes.InvalidArgument`.
- Missing demo resource, if added later, returns `codes.NotFound`.
- Unexpected failures return `codes.Internal` and are logged server-side.
- Bot client timeout maps to `ErrUnavailable` and a user-facing fallback message.

### Asynq

- Enqueue failures are returned to the caller and logged by the API.
- Invalid task payloads return an error wrapping `asynq.SkipRetry`.
- Transient handler failures return a normal error so Asynq can retry.
- Worker shutdown should let in-flight tasks finish according to Asynq/server defaults where possible.

## Security

- The gRPC listener is internal-only in Docker and Dokploy: use `expose`, not public host `ports`.
- The API gRPC server and bot client use a shared `INTERNAL_GRPC_TOKEN`.
- The token must never be logged.
- gRPC auth errors must not echo token values.
- Asynq payloads must not include secrets, raw auth headers, credentials, or Telegram tokens.
- Any future public exposure of gRPC requires a separate security review and likely mTLS or gateway-level auth.

## Docker And Deployment

Local Docker should include:

- API HTTP port as today.
- API gRPC service port exposed only inside the compose network.
- Bot env vars for `INTERNAL_API_GRPC_ADDRESS=api:9090`.
- Worker service using the worker image or app target.
- Redis dependency for API queue producer and worker.

Dokploy should include:

- `api` exposing HTTP and internal gRPC ports inside the Dokploy network.
- `bot` depending on `api` and receiving the internal gRPC address/token.
- `worker` depending on Redis and using the same image-build/release metadata pattern as other services.
- `WORKER_IMAGE` as a first-class compose variable beside `API_IMAGE`, `WEB_IMAGE`, and `BOT_IMAGE`.
- CI helper metadata and release artifacts that include worker image ref, digest, commit SHA, and pipeline ID.

The existing public web and web-admin API URLs should remain unchanged.

Update `.gitlab-ci.yml`, `tools/ci`, `deploy/dokploy/docker-compose.template.yml`, and `docs/infrastructure/ci-cd.md` in the implementation so worker build/push/deploy behavior is explicit.

## Codegen And Tooling

The implementation uses `libs/go/demoapi` as the Go proto/codegen module. Its Nx project name is `go-demoapi`.

Required targets:

- `bunx nx run go-demoapi:codegen` regenerates `libs/go/demoapi/gen/demo/v1/*.pb.go` and `libs/go/demoapi/gen/demo/v1/*_grpc.pb.go`.
- `bunx nx test go-demoapi` runs demo API contract tests.
- `bunx nx build go-demoapi` verifies generated code compiles.
- `bun run codegen` runs `api`, `web-admin`, and `go-demoapi` codegen.

The `go-demoapi:codegen` target command is `cd libs/go/demoapi && go generate ./...`. The module owns the `go:generate` directive and any generator tool pinning required to produce the committed files.

Coverage allowlist entries:

- `libs/go/demoapi/gen/demo/v1/demo.pb.go` reason: protobuf generated Go message code; replacement gate: `bunx nx run go-demoapi:codegen && bunx nx build go-demoapi`.
- `libs/go/demoapi/gen/demo/v1/demo_grpc.pb.go` reason: protobuf generated Go gRPC transport code; replacement gate: `bunx nx run go-demoapi:codegen && bunx nx build go-demoapi`.

The current Go coverage filter in `tools/coverage/run.mjs` matches exact `.go` allowlist paths by suffix/includes checks, not glob syntax. Implementation must either use the exact generated file entries above or update the coverage tool with tested glob support before using wildcard allowlist paths.

Add these Go projects to `tools/coverage/coverage.config.json`:

- `go-demoapi`: `cwd=libs/go/demoapi`, `profile=dist/coverage/go/go-demoapi/coverage.out`, `packages=./...`
- `go-tasks`: `cwd=libs/go/tasks`, `profile=dist/coverage/go/go-tasks/coverage.out`, `packages=./...`
- `go-queue`: `cwd=libs/go/queue`, `profile=dist/coverage/go/go-queue/coverage.out`, `packages=./...`
- `worker`: `cwd=apps/worker`, `profile=dist/coverage/go/worker/coverage.out`, `packages=./...`

If worker bootstrap code needs an allowlist entry, use an exact path such as `apps/worker/cmd/worker/main.go` with replacement gate `bunx nx build worker && bunx nx run worker:e2e`.

CI must fail if generated code drifts after proto changes.

The implementation should pin the Asynq module version explicitly. Asynq remains below v1, so direct dependency usage should stay inside `libs/go/queue`, `libs/go/tasks`, and `apps/worker` unless a downstream product intentionally expands it.

## Testing

### Focused Tests

- API demo service tests: success, validation, enqueue success, enqueue failure.
- API gRPC handler tests: success, invalid input, unauthenticated request, internal failure mapping.
- Bot command tests: success with fake `DemoClient`, unavailable fallback, invalid response fallback if relevant.
- Bot gRPC client tests: timeout config, metadata token injection, error mapping.
- Task tests: task factory produces expected type/payload/options; parse helper accepts valid payload and rejects invalid payload with skip-retry semantics.
- Worker handler tests: valid task succeeds; invalid payload skips retry; transient dependency error retries.
- Queue adapter tests: Redis option mapping and enqueue option defaults.

Exact focused gates:

- `bunx nx run go-demoapi:codegen`
- `bunx nx test go-demoapi`
- `bunx nx build go-demoapi`
- `bunx nx test api`
- `bunx nx test bot`
- `bunx nx test go-tasks`
- `bunx nx test go-queue`
- `bunx nx test worker`
- `bunx nx build worker`
- `bun run codegen`

### Integration Tests

Add a Docker-backed integration smoke only after focused tests pass:

1. Start Redis.
2. Start API with gRPC enabled.
3. Start worker.
4. Run a small gRPC client or test command that calls `DemoService.GetGreeting`.
5. Assert the gRPC response succeeds.
6. Call `DemoService.EnqueueEcho`.
7. Assert the response includes a task ID and queue name.
8. Assert one `demo.echo` task is processed by `apps/worker`.

This smoke should be exposed as `bunx nx run worker:e2e` and included in `bun run verify:coverage`. It can stay out of the fastest module-level tests.

### GRACE And Coverage

Update:

- `docs/requirements.xml` with the new sync/async transport use case and boilerplate constraint.
- `docs/technology.xml` with gRPC/protobuf and Asynq dependencies/tooling.
- `docs/development-plan.xml` with modules for proto contracts, internal gRPC, queue adapter, and worker.
- `docs/knowledge-graph.xml` with the new module dependencies and cross-links.
- `docs/verification-plan.xml` with proto/codegen, gRPC, queue, worker, and integration checks.
- `docs/infrastructure/ci-cd.md` with worker image and Dokploy deployment behavior.
- `.gitlab-ci.yml`, `tools/ci`, and `deploy/dokploy/docker-compose.template.yml` with worker build, release metadata, and deploy variables.

Coverage allowlist entries are allowed only for generated proto/gRPC code and bootstrap entrypoints, each with replacement gates.

## Non-Goals

- Do not replace existing REST or GraphQL user flows.
- Do not make users CRUD depend on gRPC or Asynq.
- Do not put product-specific notification logic into the template.
- Do not expose gRPC publicly.
- Do not add Redis Streams for this slice.
- Do not add RabbitMQ, NATS, Kafka, or another broker in v1.
- Do not make the bot process also run the queue worker.

## Open Implementation Choices

The implementation plan must decide only lower-level mechanics that do not change the architecture:

1. The exact generator versions used by the `go generate ./...` directive in `libs/go/demoapi`.
2. Whether the worker integration smoke starts API and worker as subprocesses or uses in-process test servers around the same production wiring.

The module paths, task trigger, Dockerfile ownership, codegen target names, and CI/deploy surface are fixed by this spec.

## Sources

- Asynq README, `hibiken/asynq`: https://github.com/hibiken/asynq
