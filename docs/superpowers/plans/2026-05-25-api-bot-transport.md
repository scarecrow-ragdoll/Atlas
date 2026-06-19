# API Bot Transport Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a reusable template baseline where `apps/bot` can call `apps/api` through internal gRPC and `apps/api` can enqueue neutral demo jobs for `apps/worker` through Asynq over Redis.

**Architecture:** The sync path is a tiny generated `DemoService` proto contract in `libs/go/demoapi`, an API-side gRPC adapter in `apps/api/internal/grpc`, and a bot-side client in `apps/bot/internal/apiclient`. The async path is a task contract package in `libs/go/tasks`, a thin Asynq adapter in `libs/go/queue`, and a separate `apps/worker` process that consumes `demo.echo` jobs from Redis. Existing REST and GraphQL users flows remain independent.

**Tech Stack:** Go 1.25, gRPC, protobuf, Buf code generation, Asynq v0.26.0, Redis, Nx, Bun, Docker Compose, GitLab CI, Dokploy, GRACE XML.

---

## Source Spec

- Design: `docs/superpowers/specs/2026-05-25-api-bot-transport-design.md`
- Status: Approved
- Fixed choices:
  - Sync transport: gRPC.
  - Async transport: Asynq over Redis.
  - Demo domain only: no users CRUD coupling.
  - Worker ownership: separate `apps/worker`.
  - gRPC auth: shared internal token through metadata.
  - No Redis Streams, RabbitMQ, NATS, Kafka, or public gRPC exposure.

## Execution Notes

- Start execution from an isolated worktree using `superpowers:using-git-worktrees`.
- The current working tree may contain user-owned edits in:
  - `apps/api/internal/middleware/auth.go`
  - `docs/development-plan.xml`
  - `docs/requirements.xml`
  - `docs/verification-plan.xml`
- Do not revert those edits. When execution reaches GRACE docs, merge the XML additions with the current file contents.
- Use TDD inside each task: add the nearest failing tests first, run the focused gate, implement, rerun the gate, then commit.
- Codegen output under `libs/go/demoapi/gen/demo/v1` is committed after `go-demoapi:codegen`.

## Implementation Choices Locked By This Plan

- Proto generator command: `cd libs/go/demoapi && go generate ./...`.
- `go generate` uses `go run github.com/bufbuild/buf/cmd/buf@v1.69.0 generate`.
- Buf remote plugins:
  - `buf.build/protocolbuffers/go:v1.36.11`
  - `buf.build/grpc/go:v1.6.2`
- Runtime modules:
  - `google.golang.org/grpc v1.81.1`
  - `google.golang.org/protobuf v1.36.11`
  - `github.com/hibiken/asynq v0.26.0`
- Integration smoke starts Docker-backed Postgres and Redis plus real API and worker subprocesses, then calls gRPC as a black-box client. This avoids illegal Go `internal` package imports while proving Redis, gRPC, queue, worker handler, and process wiring behavior.

## File Structure

- Create `libs/go/demoapi/go.mod`: Go module for the generated demo gRPC contract.
- Create `libs/go/demoapi/generate.go`: pinned `go:generate` entrypoint.
- Create `libs/go/demoapi/buf.yaml`: Buf module config.
- Create `libs/go/demoapi/buf.gen.yaml`: pinned Go and gRPC plugin config.
- Create `libs/go/demoapi/proto/demo/v1/demo.proto`: source proto for `DemoService`.
- Create `libs/go/demoapi/metadata.go`: shared internal gRPC metadata key.
- Create `libs/go/demoapi/contract_test.go`: generated contract and metadata smoke tests.
- Create `libs/go/demoapi/project.json`: Nx targets for codegen, test, build, lint.
- Generate `libs/go/demoapi/gen/demo/v1/demo.pb.go`: committed generated protobuf messages.
- Generate `libs/go/demoapi/gen/demo/v1/demo_grpc.pb.go`: committed generated gRPC bindings.
- Modify `go.work`: add `./libs/go/demoapi`, `./libs/go/tasks`, `./libs/go/queue`, and `./apps/worker`.
- Modify `package.json`: include `go-demoapi` in root `codegen`; include `worker:e2e` in release verification.
- Modify `libs/go/config/config.go`: bind explicit environment aliases for env names that do not match config paths.
- Create `libs/go/config/env_alias_test.go`: alias binding tests.
- Create `libs/go/config/testdata/env_alias.yml`: alias binding fixture.
- Create `libs/go/tasks/go.mod`: Go module for queue task contracts.
- Create `libs/go/tasks/tasks.go`: `demo.echo` task type, payload, factory, parse, and skip-retry helpers.
- Create `libs/go/tasks/tasks_test.go`: task contract tests.
- Create `libs/go/tasks/project.json`: Nx targets for test and lint.
- Create `libs/go/queue/go.mod`: Go module for the Asynq adapter.
- Create `libs/go/queue/queue.go`: Redis option mapping, enqueue defaults, server factory, demo enqueue helper.
- Create `libs/go/queue/queue_test.go`: Redis mapping and enqueue tests using `miniredis`.
- Create `libs/go/queue/project.json`: Nx targets for test and lint.
- Modify `apps/api/go.mod`: add `demoapi`, `go-queue`, `go-tasks`, `grpc`, and protobuf dependencies.
- Modify `apps/api/internal/appconfig/config.go`: add internal gRPC and queue config blocks.
- Modify `apps/api/config/config.yml`: add disabled-by-default gRPC and queue defaults.
- Create `apps/api/internal/appconfig/config_test.go`: config validation tests for enabled gRPC and queue settings.
- Create `apps/api/internal/service/demo/service.go`: transport-neutral demo service.
- Create `apps/api/internal/service/demo/service_test.go`: demo service success, validation, enqueue success, enqueue failure tests.
- Create `apps/api/internal/grpc/auth.go`: unary token auth interceptor.
- Create `apps/api/internal/grpc/auth_test.go`: unauthenticated and accepted token tests.
- Create `apps/api/internal/grpc/demo_server.go`: generated gRPC service implementation.
- Create `apps/api/internal/grpc/demo_server_test.go`: gRPC status mapping tests.
- Create `apps/api/internal/grpc/server.go`: gRPC server construction and listener runner.
- Create `apps/api/internal/grpc/server_test.go`: registration and shutdown tests.
- Modify `apps/api/cmd/server/main.go`: create queue client, demo service, and conditional gRPC listener.
- Modify `apps/bot/go.mod`: add `demoapi`, `grpc`, and protobuf dependencies.
- Modify `apps/bot/internal/appconfig/config.go`: add internal API client config.
- Modify `apps/bot/config/config.yml`: add disabled-by-default internal API defaults.
- Create `apps/bot/internal/appconfig/config_test.go`: bot internal API validation tests.
- Create `apps/bot/internal/apiclient/demo.go`: bot-side gRPC client, timeout, metadata token, error mapping.
- Create `apps/bot/internal/apiclient/demo_test.go`: client timeout, metadata, and unavailable mapping tests.
- Create `apps/bot/internal/handler/demo.go`: `/demo` command handler.
- Modify `apps/bot/internal/handler/help.go`: include `/demo`.
- Modify `apps/bot/internal/handler/handler_test.go`: add `/demo` success and fallback tests.
- Modify `apps/bot/cmd/bot/main.go`: initialize client and register `/demo`.
- Create `apps/worker/go.mod`: worker Go module.
- Create `apps/worker/config/config.yml`: Redis, queue, and log defaults.
- Create `apps/worker/internal/appconfig/config.go`: worker config types.
- Create `apps/worker/internal/appconfig/config_test.go`: worker config validation tests.
- Create `apps/worker/internal/handler/demo_echo.go`: Asynq handler for `demo.echo`.
- Create `apps/worker/internal/handler/demo_echo_test.go`: valid payload, skip retry, and transient retry tests.
- Create `apps/worker/cmd/worker/main.go`: worker process entrypoint.
- Create `apps/worker/project.json`: Nx targets for build, serve, test, lint, and e2e.
- Create `apps/worker/air.toml`: local worker reload config.
- Create `apps/worker/internal/integration/transport_smoke_test.go`: Redis-backed gRPC plus worker smoke.
- Create `docker/worker.Dockerfile`: worker image build.
- Modify `docker/api.Dockerfile`: copy new Go modules needed by API builds.
- Modify `docker/bot.Dockerfile`: copy `libs/go/demoapi` needed by bot builds.
- Modify `docker/docker-compose.dev.yml`: add API gRPC env, bot internal API env, and worker service.
- Modify `docker/docker-compose.yml`: expose API gRPC inside the compose network and add worker.
- Modify `docker/docker-compose.test.yml`: ensure Redis is available for `worker:e2e`.
- Create `tools/e2e/worker-smoke.sh`: deterministic Redis-backed worker smoke runner.
- Modify `tools/coverage/coverage.config.json`: add new Go projects and exact generated-code allowlist entries.
- Modify `tools/coverage/preflight.mjs`: require worker project, worker Dockerfile, and worker smoke runner.
- Create `tools/codegen/check-drift.sh`: run codegen and fail when committed generated files drift.
- Modify `tools/ci/src/core.ts`: add `worker` to deployable services and release metadata.
- Modify `tools/ci/src/core.test.ts`: assert worker image refs, digest metadata, and Dokploy env rendering.
- Modify `tools/ci/src/cli.ts`: include worker in Docker build/push commands through the shared service list.
- Modify `tools/ci/src/cli.test.ts`: assert worker build/push dispatch.
- Modify `.gitlab-ci.yml`: build, push, and release the worker image.
- Modify `deploy/dokploy/docker-compose.template.yml`: add `WORKER_IMAGE` and worker service.
- Modify `docs/infrastructure/ci-cd.md`: document worker image, internal gRPC, and release metadata.
- Modify `docs/requirements.xml`: add the boilerplate transport use case and constraints.
- Modify `docs/technology.xml`: add gRPC, protobuf, Buf, and Asynq tooling.
- Modify `docs/development-plan.xml`: add `M-DEMOAPI`, `M-API-GRPC`, `M-API-DEMO-SERVICE`, `M-BOT-API-CLIENT`, `M-GO-TASKS`, `M-GO-QUEUE`, and `M-WORKER`.
- Modify `docs/knowledge-graph.xml`: add new module nodes and edges.
- Modify `docs/verification-plan.xml`: add `V-M-*` verification refs for all new modules and the integration smoke.
- Modify `docs/operational-packets.xml`: add a concrete note for API bot transport packet write scopes.

## Task 0: Config Environment Alias Support

**Files:**

- Modify: `libs/go/config/config.go`
- Create: `libs/go/config/env_alias_test.go`
- Create: `libs/go/config/testdata/env_alias.yml`

- [ ] **Step 1: Add the failing alias binding test**

Create `libs/go/config/testdata/env_alias.yml`:

```yaml
service:
  token: ''
```

Create `libs/go/config/env_alias_test.go`:

```go
package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"monorepo-template/libs/go/config"
)

type envAliasConfig struct {
	Service struct {
		Token string `mapstructure:"token" validate:"required"`
	} `mapstructure:"service"`
}

func TestLoadBindsExplicitEnvAliases(t *testing.T) {
	t.Setenv("INTERNAL_GRPC_TOKEN", "secret-token")

	cfg, err := config.Load[envAliasConfig](config.Options{
		ConfigPath: "testdata/env_alias.yml",
		EnvAliases: map[string]string{
			"service.token": "INTERNAL_GRPC_TOKEN",
		},
	})

	require.NoError(t, err)
	require.Equal(t, "secret-token", cfg.Service.Token)
}
```

- [ ] **Step 2: Run the failing config test**

Run: `cd libs/go/config && go test -run TestLoadBindsExplicitEnvAliases ./...`

Expected: FAIL because `config.Options` has no `EnvAliases` field.

- [ ] **Step 3: Add alias binding to the shared config loader**

Modify `libs/go/config/config.go`:

```go
type Options struct {
	// ConfigPath is the path to the YAML config file (required).
	ConfigPath string

	// EnvFile is the path to a .env file (optional, "" = skip).
	EnvFile string

	// EnvPrefix is an optional prefix for env var lookup (e.g. "API" -> API_POSTGRES_HOST).
	EnvPrefix string

	// EnvAliases binds config keys to env var names when the env var cannot match the key automatically.
	EnvAliases map[string]string
}
```

Add this block after `v.SetEnvKeyReplacer(...)` and before `v.AutomaticEnv()`:

```go
for key, envName := range opts.EnvAliases {
	if err := v.BindEnv(key, envName); err != nil {
		return zero, fmt.Errorf("%w: bind env alias %s: %v", ErrConfigLoad, key, err)
	}
}
```

Keep `v.AutomaticEnv()` after the alias loop:

```go
v.AutomaticEnv()
```

- [ ] **Step 4: Verify and commit**

Run: `bunx nx test go-config`

Expected: PASS including `TestLoadBindsExplicitEnvAliases`.

```bash
git add libs/go/config/config.go libs/go/config/env_alias_test.go libs/go/config/testdata/env_alias.yml
git commit -m "feat: support config env aliases"
```

## Task 1: Demo gRPC Contract Module

**Files:**

- Create: `libs/go/demoapi/go.mod`
- Create: `libs/go/demoapi/generate.go`
- Create: `libs/go/demoapi/buf.yaml`
- Create: `libs/go/demoapi/buf.gen.yaml`
- Create: `libs/go/demoapi/proto/demo/v1/demo.proto`
- Create: `libs/go/demoapi/metadata.go`
- Create: `libs/go/demoapi/contract_test.go`
- Create: `libs/go/demoapi/project.json`
- Generate: `libs/go/demoapi/gen/demo/v1/demo.pb.go`
- Generate: `libs/go/demoapi/gen/demo/v1/demo_grpc.pb.go`
- Modify: `go.work`
- Modify: `package.json`

- [ ] **Step 1: Add the failing generated-contract test**

Create `libs/go/demoapi/contract_test.go`:

```go
package demoapi_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	"monorepo-template/libs/go/demoapi"
	demov1 "monorepo-template/libs/go/demoapi/gen/demo/v1"
)

func TestGeneratedGreetingResponseRoundTrips(t *testing.T) {
	original := &demov1.GetGreetingResponse{
		Message:   "Hello, Template!",
		RequestId: "request-1",
	}

	data, err := proto.Marshal(original)
	require.NoError(t, err)

	var decoded demov1.GetGreetingResponse
	require.NoError(t, proto.Unmarshal(data, &decoded))
	require.Equal(t, original.Message, decoded.Message)
	require.Equal(t, original.RequestId, decoded.RequestId)
}

func TestGeneratedEnqueueEchoResponseFields(t *testing.T) {
	response := &demov1.EnqueueEchoResponse{
		TaskId:    "task-1",
		Queue:     "default",
		RequestId: "request-2",
	}

	require.Equal(t, "task-1", response.GetTaskId())
	require.Equal(t, "default", response.GetQueue())
	require.Equal(t, "request-2", response.GetRequestId())
}

func TestInternalTokenMetadataKeyIsStable(t *testing.T) {
	require.Equal(t, "x-internal-token", demoapi.InternalTokenMetadataKey)
}

func TestAppendInternalTokenAddsOutgoingMetadata(t *testing.T) {
	ctx := demoapi.AppendInternalToken(context.Background(), "secret")
	md, ok := metadata.FromOutgoingContext(ctx)

	require.True(t, ok)
	require.Equal(t, []string{"secret"}, md.Get(demoapi.InternalTokenMetadataKey))
}
```

- [ ] **Step 2: Add the demoapi module and generator config**

Create `libs/go/demoapi/go.mod`:

```go
module monorepo-template/libs/go/demoapi

go 1.25.0

require (
	github.com/stretchr/testify v1.11.1
	google.golang.org/grpc v1.81.1
	google.golang.org/protobuf v1.36.11
)
```

Create `libs/go/demoapi/generate.go`:

```go
package demoapi

//go:generate go run github.com/bufbuild/buf/cmd/buf@v1.69.0 generate
```

Create `libs/go/demoapi/metadata.go`:

```go
package demoapi

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const InternalTokenMetadataKey = "x-internal-token"

func AppendInternalToken(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, InternalTokenMetadataKey, token)
}
```

Create `libs/go/demoapi/buf.yaml`:

```yaml
version: v2
lint:
  use:
    - STANDARD
breaking:
  use:
    - FILE
```

Create `libs/go/demoapi/buf.gen.yaml`:

```yaml
version: v2
plugins:
  - remote: buf.build/protocolbuffers/go:v1.36.11
    out: .
    opt:
      - paths=import
      - module=monorepo-template/libs/go/demoapi
  - remote: buf.build/grpc/go:v1.6.2
    out: .
    opt:
      - paths=import
      - module=monorepo-template/libs/go/demoapi
```

Create `libs/go/demoapi/proto/demo/v1/demo.proto`:

```proto
syntax = "proto3";

package demo.v1;

option go_package = "monorepo-template/libs/go/demoapi/gen/demo/v1;demov1";

service DemoService {
  rpc GetGreeting(GetGreetingRequest) returns (GetGreetingResponse);
  rpc EnqueueEcho(EnqueueEchoRequest) returns (EnqueueEchoResponse);
}

message GetGreetingRequest {
  string name = 1;
}

message GetGreetingResponse {
  string message = 1;
  string request_id = 2;
}

message EnqueueEchoRequest {
  string message = 1;
}

message EnqueueEchoResponse {
  string task_id = 1;
  string queue = 2;
  string request_id = 3;
}
```

- [ ] **Step 3: Add the Nx project**

Create `libs/go/demoapi/project.json`:

```json
{
  "name": "go-demoapi",
  "$schema": "../../../node_modules/nx/schemas/project-schema.json",
  "sourceRoot": "libs/go/demoapi",
  "projectType": "library",
  "tags": ["scope:shared", "lang:go"],
  "targets": {
    "codegen": {
      "executor": "nx:run-commands",
      "options": {
        "command": "cd libs/go/demoapi && go generate ./..."
      }
    },
    "build": {
      "executor": "nx:run-commands",
      "options": {
        "command": "cd libs/go/demoapi && go build ./..."
      },
      "dependsOn": ["codegen"]
    },
    "test": {
      "executor": "nx:run-commands",
      "options": {
        "command": "cd libs/go/demoapi && mkdir -p ../../../dist/coverage/go/go-demoapi && go test -v -coverprofile=../../../dist/coverage/go/go-demoapi/coverage.out ./..."
      }
    },
    "lint": {
      "executor": "nx:run-commands",
      "options": {
        "command": "cd libs/go/demoapi && golangci-lint run"
      }
    }
  }
}
```

- [ ] **Step 4: Register demoapi in workspace tooling**

Modify `go.work` so the `use` block is:

```go
use (
	./apps/api
	./apps/bot
	./libs/go/config
	./libs/go/demoapi
	./libs/go/logger
)
```

Modify `package.json` root script:

```json
"codegen": "bunx nx run-many --target=codegen --projects=api,web-admin,go-demoapi"
```

- [ ] **Step 5: Run the failing test before generated files exist**

Run: `bunx nx test go-demoapi`

Expected: FAIL with an import error for `monorepo-template/libs/go/demoapi/gen/demo/v1`.

- [ ] **Step 6: Generate and commit the generated gRPC code**

Run: `bunx nx run go-demoapi:codegen`

Expected:

```text
NX   Successfully ran target codegen for project go-demoapi
```

Confirm generated files:

```bash
test -f libs/go/demoapi/gen/demo/v1/demo.pb.go
test -f libs/go/demoapi/gen/demo/v1/demo_grpc.pb.go
```

- [ ] **Step 7: Verify the module**

Run: `bunx nx test go-demoapi`

Expected: PASS with `TestGeneratedGreetingResponseRoundTrips` and `TestGeneratedEnqueueEchoResponseFields`.

Run: `bunx nx build go-demoapi`

Expected: PASS.

- [ ] **Step 8: Commit**

```bash
git add go.work package.json libs/go/demoapi
git commit -m "feat: add demo gRPC contract"
```

## Task 2: Asynq Task Contract Library

**Files:**

- Create: `libs/go/tasks/go.mod`
- Create: `libs/go/tasks/tasks.go`
- Create: `libs/go/tasks/tasks_test.go`
- Create: `libs/go/tasks/project.json`
- Modify: `go.work`

- [ ] **Step 1: Add failing task contract tests**

Create `libs/go/tasks/tasks_test.go`:

```go
package tasks_test

import (
	"errors"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"

	"monorepo-template/libs/go/tasks"
)

func TestNewDemoEchoTaskProducesExpectedTypeAndPayload(t *testing.T) {
	task, err := tasks.NewDemoEchoTask(tasks.DemoEchoPayload{
		Message:   " hello ",
		RequestID: "request-1",
	})
	require.NoError(t, err)
	require.Equal(t, tasks.TypeDemoEcho, task.Type())

	payload, err := tasks.ParseDemoEchoTask(task)
	require.NoError(t, err)
	require.Equal(t, "hello", payload.Message)
	require.Equal(t, "request-1", payload.RequestID)
}

func TestNewDemoEchoTaskRejectsEmptyMessage(t *testing.T) {
	_, err := tasks.NewDemoEchoTask(tasks.DemoEchoPayload{Message: "  "})
	require.ErrorIs(t, err, tasks.ErrInvalidPayload)
}

func TestParseDemoEchoTaskRejectsWrongType(t *testing.T) {
	task := asynq.NewTask("other.type", []byte(`{"message":"hello"}`))

	_, err := tasks.ParseDemoEchoTask(task)
	require.ErrorIs(t, err, tasks.ErrInvalidPayload)
}

func TestParseDemoEchoTaskRejectsBadJSON(t *testing.T) {
	task := asynq.NewTask(tasks.TypeDemoEcho, []byte(`{`))

	_, err := tasks.ParseDemoEchoTask(task)
	require.ErrorIs(t, err, tasks.ErrInvalidPayload)
}

func TestSkipRetryInvalidPayloadWrapsAsynqSentinel(t *testing.T) {
	err := tasks.SkipRetryInvalidPayload(tasks.ErrInvalidPayload)

	require.ErrorIs(t, err, asynq.SkipRetry)
	require.True(t, tasks.IsInvalidPayload(err))
	require.True(t, errors.Is(err, tasks.ErrInvalidPayload))
}
```

- [ ] **Step 2: Add the task module**

Create `libs/go/tasks/go.mod`:

```go
module monorepo-template/libs/go/tasks

go 1.25.0

require (
	github.com/hibiken/asynq v0.26.0
	github.com/stretchr/testify v1.11.1
)
```

Create `libs/go/tasks/tasks.go`:

```go
package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hibiken/asynq"
)

const TypeDemoEcho = "demo.echo"

var ErrInvalidPayload = errors.New("invalid task payload")

type DemoEchoPayload struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

func NewDemoEchoTask(payload DemoEchoPayload) (*asynq.Task, error) {
	payload.Message = strings.TrimSpace(payload.Message)
	if payload.Message == "" {
		return nil, fmt.Errorf("%w: message is required", ErrInvalidPayload)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal demo echo payload: %w", err)
	}

	return asynq.NewTask(TypeDemoEcho, body), nil
}

func ParseDemoEchoTask(task *asynq.Task) (DemoEchoPayload, error) {
	if task.Type() != TypeDemoEcho {
		return DemoEchoPayload{}, fmt.Errorf("%w: unexpected task type %q", ErrInvalidPayload, task.Type())
	}

	var payload DemoEchoPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return DemoEchoPayload{}, fmt.Errorf("%w: decode payload: %v", ErrInvalidPayload, err)
	}

	payload.Message = strings.TrimSpace(payload.Message)
	if payload.Message == "" {
		return DemoEchoPayload{}, fmt.Errorf("%w: message is required", ErrInvalidPayload)
	}

	return payload, nil
}

func IsInvalidPayload(err error) bool {
	return errors.Is(err, ErrInvalidPayload)
}

func SkipRetryInvalidPayload(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w: %w", asynq.SkipRetry, err)
}
```

- [ ] **Step 3: Add the Nx project**

Create `libs/go/tasks/project.json`:

```json
{
  "name": "go-tasks",
  "$schema": "../../../node_modules/nx/schemas/project-schema.json",
  "sourceRoot": "libs/go/tasks",
  "projectType": "library",
  "tags": ["scope:shared", "lang:go"],
  "targets": {
    "test": {
      "executor": "nx:run-commands",
      "options": {
        "command": "cd libs/go/tasks && mkdir -p ../../../dist/coverage/go/go-tasks && go test -v -coverprofile=../../../dist/coverage/go/go-tasks/coverage.out ./..."
      }
    },
    "lint": {
      "executor": "nx:run-commands",
      "options": {
        "command": "cd libs/go/tasks && golangci-lint run"
      }
    }
  }
}
```

- [ ] **Step 4: Register go-tasks in `go.work`**

Modify `go.work`:

```go
use (
	./apps/api
	./apps/bot
	./libs/go/config
	./libs/go/demoapi
	./libs/go/logger
	./libs/go/tasks
)
```

- [ ] **Step 5: Verify and commit**

Run: `bunx nx test go-tasks`

Expected: PASS with all five task contract tests.

```bash
git add go.work libs/go/tasks
git commit -m "feat: add demo queue task contract"
```

## Task 3: Thin Asynq Queue Adapter

**Files:**

- Create: `libs/go/queue/go.mod`
- Create: `libs/go/queue/queue.go`
- Create: `libs/go/queue/queue_test.go`
- Create: `libs/go/queue/project.json`
- Modify: `go.work`

- [ ] **Step 1: Add failing queue adapter tests**

Create `libs/go/queue/queue_test.go`:

```go
package queue_test

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/require"

	"monorepo-template/libs/go/config"
	"monorepo-template/libs/go/queue"
	"monorepo-template/libs/go/tasks"
)

func TestRedisClientOptMapsSharedRedisConfig(t *testing.T) {
	opt := queue.RedisClientOpt(config.RedisConfig{
		Host:     "redis",
		Port:     6379,
		Password: "secret",
		DB:       2,
	})

	require.Equal(t, "redis:6379", opt.Addr)
	require.Equal(t, "secret", opt.Password)
	require.Equal(t, 2, opt.DB)
}

func TestNormalizeConfigDefaultsQueueConcurrencyRetryAndTimeout(t *testing.T) {
	cfg := queue.NormalizeConfig(queue.Config{})

	require.Equal(t, "default", cfg.DefaultQueue)
	require.Equal(t, 5, cfg.Concurrency)
	require.Equal(t, map[string]int{"default": 1}, cfg.Queues)
	require.Equal(t, 3, cfg.MaxRetry)
	require.Equal(t, 30*time.Second, cfg.EnqueueTimeout)
}

func TestClientEnqueueDemoEchoUsesDefaultQueue(t *testing.T) {
	redis := miniredis.RunT(t)
	host, portText, err := net.SplitHostPort(redis.Addr())
	require.NoError(t, err)
	port, err := strconv.Atoi(portText)
	require.NoError(t, err)

	client := queue.NewClient(queue.Config{
		Redis:        config.RedisConfig{Host: host, Port: port},
		DefaultQueue: "critical",
	})
	defer client.Close()

	info, err := client.EnqueueDemoEcho(context.Background(), tasks.DemoEchoPayload{Message: "hello"})
	require.NoError(t, err)
	require.NotEmpty(t, info.ID)
	require.Equal(t, "critical", info.Queue)
}
```

- [ ] **Step 2: Add the queue module**

Create `libs/go/queue/go.mod`:

```go
module monorepo-template/libs/go/queue

go 1.25.0

require (
	github.com/alicebob/miniredis/v2 v2.38.0
	github.com/hibiken/asynq v0.26.0
	github.com/stretchr/testify v1.11.1
	monorepo-template/libs/go/config v0.0.0
	monorepo-template/libs/go/tasks v0.0.0
)

replace (
	monorepo-template/libs/go/config => ../config
	monorepo-template/libs/go/tasks => ../tasks
)
```

Create `libs/go/queue/queue.go`:

```go
package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"

	"monorepo-template/libs/go/config"
	"monorepo-template/libs/go/tasks"
)

type Config struct {
	Redis          config.RedisConfig
	DefaultQueue   string
	Queues         map[string]int
	Concurrency    int
	MaxRetry       int
	EnqueueTimeout time.Duration
}

type EnqueueResult struct {
	ID    string
	Queue string
}

type Client struct {
	client *asynq.Client
	cfg    Config
}

func NormalizeConfig(cfg Config) Config {
	if cfg.DefaultQueue == "" {
		cfg.DefaultQueue = "default"
	}
	if cfg.Concurrency == 0 {
		cfg.Concurrency = 5
	}
	if len(cfg.Queues) == 0 {
		cfg.Queues = map[string]int{cfg.DefaultQueue: 1}
	}
	if cfg.MaxRetry == 0 {
		cfg.MaxRetry = 3
	}
	if cfg.EnqueueTimeout == 0 {
		cfg.EnqueueTimeout = 30 * time.Second
	}
	return cfg
}

func RedisClientOpt(cfg config.RedisConfig) asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	}
}

func NewClient(cfg Config) *Client {
	normalized := NormalizeConfig(cfg)
	return &Client{
		client: asynq.NewClient(RedisClientOpt(normalized.Redis)),
		cfg:    normalized,
	}
}

func (c *Client) Enqueue(ctx context.Context, task *asynq.Task) (*asynq.TaskInfo, error) {
	return c.client.EnqueueContext(
		ctx,
		task,
		asynq.Queue(c.cfg.DefaultQueue),
		asynq.MaxRetry(c.cfg.MaxRetry),
		asynq.Timeout(c.cfg.EnqueueTimeout),
	)
}

func (c *Client) EnqueueDemoEcho(ctx context.Context, payload tasks.DemoEchoPayload) (EnqueueResult, error) {
	task, err := tasks.NewDemoEchoTask(payload)
	if err != nil {
		return EnqueueResult{}, err
	}

	info, err := c.Enqueue(ctx, task)
	if err != nil {
		return EnqueueResult{}, err
	}

	return EnqueueResult{ID: info.ID, Queue: info.Queue}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func NewServer(cfg Config) *asynq.Server {
	normalized := NormalizeConfig(cfg)
	return asynq.NewServer(RedisClientOpt(normalized.Redis), asynq.Config{
		Concurrency: normalized.Concurrency,
		Queues:      normalized.Queues,
	})
}
```

- [ ] **Step 3: Add the Nx project and workspace entry**

Create `libs/go/queue/project.json`:

```json
{
  "name": "go-queue",
  "$schema": "../../../node_modules/nx/schemas/project-schema.json",
  "sourceRoot": "libs/go/queue",
  "projectType": "library",
  "tags": ["scope:shared", "lang:go"],
  "targets": {
    "test": {
      "executor": "nx:run-commands",
      "options": {
        "command": "cd libs/go/queue && mkdir -p ../../../dist/coverage/go/go-queue && go test -v -coverprofile=../../../dist/coverage/go/go-queue/coverage.out ./..."
      }
    },
    "lint": {
      "executor": "nx:run-commands",
      "options": {
        "command": "cd libs/go/queue && golangci-lint run"
      }
    }
  }
}
```

Modify `go.work`:

```go
use (
	./apps/api
	./apps/bot
	./libs/go/config
	./libs/go/demoapi
	./libs/go/logger
	./libs/go/queue
	./libs/go/tasks
)
```

- [ ] **Step 4: Verify and commit**

Run: `bunx nx test go-queue`

Expected: PASS with all three queue adapter tests.

```bash
git add go.work libs/go/queue
git commit -m "feat: add asynq queue adapter"
```

## Task 4: API Config And Transport-Neutral Demo Service

**Files:**

- Modify: `apps/api/go.mod`
- Modify: `apps/api/internal/appconfig/config.go`
- Modify: `apps/api/config/config.yml`
- Create: `apps/api/internal/appconfig/config_test.go`
- Create: `apps/api/internal/service/demo/service.go`
- Create: `apps/api/internal/service/demo/service_test.go`

- [ ] **Step 1: Add API config tests**

Create `apps/api/internal/appconfig/config_test.go`:

```go
package appconfig_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/appconfig"
	"monorepo-template/libs/go/config"
)

func validConfig() appconfig.Config {
	return appconfig.Config{
		Server: config.ServerConfig{Port: 8090},
		Log:    config.LogConfig{Level: "info", Format: "json"},
		Postgres: config.PostgresConfig{
			Host: "localhost",
			Port: 5432,
			User: "app",
			DB:   "app",
		},
		Redis:      config.RedisConfig{Host: "localhost", Port: 6379},
		Auth:       appconfig.AuthConfig{JWTSecret: "secret"},
		Pagination: appconfig.PaginationConfig{DefaultPageSize: 20, MaxPageSize: 100},
		GRPC:       appconfig.GRPCConfig{Enabled: false, Port: 9090, ReadTimeout: time.Second},
		Queue:      appconfig.QueueConfig{Enabled: false, DefaultQueue: "default"},
	}
}

func TestConfigAllowsDisabledGRPCWithoutToken(t *testing.T) {
	require.NoError(t, config.Validate(validConfig()))
}

func TestConfigRequiresGRPCPortAndTokenWhenEnabled(t *testing.T) {
	cfg := validConfig()
	cfg.GRPC.Enabled = true
	cfg.GRPC.Port = 0
	cfg.GRPC.AuthToken = ""

	err := config.Validate(cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Port")
	require.Contains(t, err.Error(), "AuthToken")
}

func TestConfigAcceptsEnabledGRPCAndQueue(t *testing.T) {
	cfg := validConfig()
	cfg.GRPC.Enabled = true
	cfg.GRPC.AuthToken = "internal"
	cfg.Queue.Enabled = true
	cfg.Queue.DefaultQueue = "critical"

	require.NoError(t, config.Validate(cfg))
}
```

- [ ] **Step 2: Extend API config types**

Modify `apps/api/internal/appconfig/config.go`:

```go
package appconfig

import (
	"time"

	"monorepo-template/libs/go/config"
)

type AuthConfig struct {
	JWTSecret string `mapstructure:"jwt_secret" validate:"required"`
}

type PaginationConfig struct {
	DefaultPageSize int `mapstructure:"default_page_size" validate:"gt=0"`
	MaxPageSize     int `mapstructure:"max_page_size"     validate:"gt=0"`
}

type GRPCConfig struct {
	Enabled     bool          `mapstructure:"enabled"`
	Port        int           `mapstructure:"port"         validate:"required_if=Enabled true,omitempty,gt=0"`
	AuthToken   string        `mapstructure:"auth_token"   validate:"required_if=Enabled true"`
	ReadTimeout time.Duration `mapstructure:"read_timeout" validate:"required_if=Enabled true"`
}

type QueueConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	DefaultQueue string `mapstructure:"default_queue" validate:"required_if=Enabled true"`
}

type Config struct {
	Server     config.ServerConfig   `mapstructure:"server"`
	Log        config.LogConfig      `mapstructure:"log"`
	Postgres   config.PostgresConfig `mapstructure:"postgres"`
	Redis      config.RedisConfig    `mapstructure:"redis"`
	Auth       AuthConfig            `mapstructure:"auth"`
	Pagination PaginationConfig      `mapstructure:"pagination"`
	GRPC       GRPCConfig            `mapstructure:"grpc"`
	Queue      QueueConfig           `mapstructure:"queue"`
}
```

Modify `apps/api/config/config.yml`:

```yaml
grpc:
  enabled: false
  port: 9090
  auth_token: ''
  read_timeout: 10s

queue:
  enabled: false
  default_queue: default
```

- [ ] **Step 3: Add demo service tests**

Create `apps/api/internal/service/demo/service_test.go`:

```go
package demo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/service/demo"
	"monorepo-template/libs/go/queue"
	"monorepo-template/libs/go/tasks"
)

type fakeEnqueuer struct {
	result  queue.EnqueueResult
	err     error
	payload tasks.DemoEchoPayload
	called  bool
}

func (f *fakeEnqueuer) EnqueueDemoEcho(ctx context.Context, payload tasks.DemoEchoPayload) (queue.EnqueueResult, error) {
	f.called = true
	f.payload = payload
	return f.result, f.err
}

func TestGetGreetingReturnsGreetingWithoutQueueSideEffect(t *testing.T) {
	enqueuer := &fakeEnqueuer{}
	service := demo.NewService(enqueuer)

	response, err := service.GetGreeting(context.Background(), demo.GetGreetingInput{Name: "Template"})

	require.NoError(t, err)
	require.Equal(t, "Hello, Template!", response.Message)
	require.NotEmpty(t, response.RequestID)
	require.False(t, enqueuer.called)
}

func TestGetGreetingUsesNeutralDefaultName(t *testing.T) {
	service := demo.NewService(&fakeEnqueuer{})

	response, err := service.GetGreeting(context.Background(), demo.GetGreetingInput{Name: "  "})

	require.NoError(t, err)
	require.Equal(t, "Hello, there!", response.Message)
}

func TestEnqueueEchoRejectsEmptyMessage(t *testing.T) {
	service := demo.NewService(&fakeEnqueuer{})

	_, err := service.EnqueueEcho(context.Background(), demo.EnqueueEchoInput{Message: " "})

	require.ErrorIs(t, err, demo.ErrInvalidInput)
}

func TestEnqueueEchoReturnsTaskIDAndQueue(t *testing.T) {
	enqueuer := &fakeEnqueuer{result: queue.EnqueueResult{ID: "task-1", Queue: "critical"}}
	service := demo.NewService(enqueuer)

	result, err := service.EnqueueEcho(context.Background(), demo.EnqueueEchoInput{Message: " hello "})

	require.NoError(t, err)
	require.Equal(t, "task-1", result.TaskID)
	require.Equal(t, "critical", result.Queue)
	require.NotEmpty(t, result.RequestID)
	require.Equal(t, "hello", enqueuer.payload.Message)
	require.Equal(t, result.RequestID, enqueuer.payload.RequestID)
}

func TestEnqueueEchoReturnsProducerError(t *testing.T) {
	enqueuer := &fakeEnqueuer{err: errors.New("redis down")}
	service := demo.NewService(enqueuer)

	_, err := service.EnqueueEcho(context.Background(), demo.EnqueueEchoInput{Message: "hello"})

	require.ErrorContains(t, err, "redis down")
}
```

- [ ] **Step 4: Add demo service implementation**

Create `apps/api/internal/service/demo/service.go`:

```go
package demo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"monorepo-template/libs/go/queue"
	"monorepo-template/libs/go/tasks"
)

var ErrInvalidInput = errors.New("invalid demo input")

type EchoEnqueuer interface {
	EnqueueDemoEcho(ctx context.Context, payload tasks.DemoEchoPayload) (queue.EnqueueResult, error)
}

type Service struct {
	enqueuer EchoEnqueuer
}

type GetGreetingInput struct {
	Name string
}

type GetGreetingResult struct {
	Message   string
	RequestID string
}

type EnqueueEchoInput struct {
	Message string
}

type EnqueueEchoResult struct {
	TaskID    string
	Queue     string
	RequestID string
}

func NewService(enqueuer EchoEnqueuer) *Service {
	return &Service{enqueuer: enqueuer}
}

func (s *Service) GetGreeting(ctx context.Context, input GetGreetingInput) (GetGreetingResult, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		name = "there"
	}

	return GetGreetingResult{
		Message:   fmt.Sprintf("Hello, %s!", name),
		RequestID: uuid.NewString(),
	}, nil
}

func (s *Service) EnqueueEcho(ctx context.Context, input EnqueueEchoInput) (EnqueueEchoResult, error) {
	message := strings.TrimSpace(input.Message)
	if message == "" {
		return EnqueueEchoResult{}, fmt.Errorf("%w: message is required", ErrInvalidInput)
	}
	if s.enqueuer == nil {
		return EnqueueEchoResult{}, errors.New("demo echo enqueuer is not configured")
	}

	requestID := uuid.NewString()
	result, err := s.enqueuer.EnqueueDemoEcho(ctx, tasks.DemoEchoPayload{
		Message:   message,
		RequestID: requestID,
	})
	if err != nil {
		return EnqueueEchoResult{}, fmt.Errorf("enqueue demo echo: %w", err)
	}

	return EnqueueEchoResult{
		TaskID:    result.ID,
		Queue:     result.Queue,
		RequestID: requestID,
	}, nil
}
```

Update `apps/api/go.mod` direct requirements:

```go
github.com/google/uuid v1.6.0
google.golang.org/grpc v1.81.1
google.golang.org/protobuf v1.36.11
monorepo-template/libs/go/demoapi v0.0.0
monorepo-template/libs/go/queue v0.0.0
monorepo-template/libs/go/tasks v0.0.0
```

Update `apps/api/go.mod` replace block:

```go
replace (
	monorepo-template/libs/go/config => ../../libs/go/config
	monorepo-template/libs/go/demoapi => ../../libs/go/demoapi
	monorepo-template/libs/go/logger => ../../libs/go/logger
	monorepo-template/libs/go/queue => ../../libs/go/queue
	monorepo-template/libs/go/tasks => ../../libs/go/tasks
)
```

- [ ] **Step 5: Verify and commit**

Run: `bunx nx test api`

Expected: PASS including API config tests and demo service tests.

```bash
git add apps/api/go.mod apps/api/go.sum apps/api/config/config.yml apps/api/internal/appconfig apps/api/internal/service/demo
git commit -m "feat: add api demo service"
```

## Task 5: API gRPC Server Adapter

**Files:**

- Create: `apps/api/internal/grpc/auth.go`
- Create: `apps/api/internal/grpc/auth_test.go`
- Create: `apps/api/internal/grpc/logging.go`
- Create: `apps/api/internal/grpc/logging_test.go`
- Create: `apps/api/internal/grpc/demo_server.go`
- Create: `apps/api/internal/grpc/demo_server_test.go`
- Create: `apps/api/internal/grpc/server.go`
- Create: `apps/api/internal/grpc/server_test.go`
- Modify: `apps/api/cmd/server/main.go`

- [ ] **Step 1: Add auth interceptor tests**

Create `apps/api/internal/grpc/auth_test.go`:

```go
package grpc_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	apigrpc "monorepo-template/apps/api/internal/grpc"
)

func TestAuthUnaryInterceptorRejectsMissingToken(t *testing.T) {
	interceptor := apigrpc.AuthUnaryInterceptor("secret")

	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/demo.v1.DemoService/GetGreeting"}, func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	})

	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestAuthUnaryInterceptorRejectsInvalidToken(t *testing.T) {
	interceptor := apigrpc.AuthUnaryInterceptor("secret")
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(apigrpc.InternalTokenMetadataKey, "wrong"))

	_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/demo.v1.DemoService/GetGreeting"}, func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	})

	require.Equal(t, codes.Unauthenticated, status.Code(err))
	require.NotContains(t, err.Error(), "wrong")
}

func TestAuthUnaryInterceptorAcceptsValidToken(t *testing.T) {
	interceptor := apigrpc.AuthUnaryInterceptor("secret")
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(apigrpc.InternalTokenMetadataKey, "secret"))

	got, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/demo.v1.DemoService/GetGreeting"}, func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	})

	require.NoError(t, err)
	require.Equal(t, "ok", got)
}
```

- [ ] **Step 2: Implement auth interceptor**

Create `apps/api/internal/grpc/auth.go`:

```go
package grpc

import (
	"context"

	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"monorepo-template/libs/go/demoapi"
)

const InternalTokenMetadataKey = demoapi.InternalTokenMetadataKey

func AuthUnaryInterceptor(expectedToken string) grpcpkg.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpcpkg.UnaryServerInfo, handler grpcpkg.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "internal auth token is required")
		}

		values := md.Get(InternalTokenMetadataKey)
		if len(values) != 1 || values[0] != expectedToken {
			return nil, status.Error(codes.Unauthenticated, "internal auth token is invalid")
		}

		return handler(ctx, req)
	}
}
```

- [ ] **Step 3: Add request logging interceptor tests**

Create `apps/api/internal/grpc/logging_test.go`:

```go
package grpc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apigrpc "monorepo-template/apps/api/internal/grpc"
)

func TestLoggingUnaryInterceptorRecordsMethodCodeAndDuration(t *testing.T) {
	core, observed := observer.New(zap.InfoLevel)
	interceptor := apigrpc.LoggingUnaryInterceptor(zap.New(core))

	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/demo.v1.DemoService/GetGreeting"}, func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	})

	require.NoError(t, err)
	entry := observed.FilterMessage("grpc request completed").All()[0]
	require.Equal(t, "/demo.v1.DemoService/GetGreeting", entry.ContextMap()["method"])
	require.Equal(t, codes.OK.String(), entry.ContextMap()["code"])
	require.Contains(t, entry.ContextMap(), "duration_ms")
}

func TestLoggingUnaryInterceptorDoesNotLogTokenValues(t *testing.T) {
	core, observed := observer.New(zap.InfoLevel)
	interceptor := apigrpc.LoggingUnaryInterceptor(zap.New(core))

	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/demo.v1.DemoService/GetGreeting"}, func(ctx context.Context, req any) (any, error) {
		return nil, status.Error(codes.Unauthenticated, "secret-token")
	})

	require.Error(t, err)
	entry := observed.FilterMessage("grpc request completed").All()[0]
	require.NotContains(t, entry.ContextMap(), "token")
	require.NotContains(t, entry.ContextMap(), "metadata")
	require.NotContains(t, fmt.Sprint(entry.ContextMap()), "secret-token")
}
```

- [ ] **Step 4: Implement request logging interceptor**

Create `apps/api/internal/grpc/logging.go`:

```go
package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingUnaryInterceptor(log *zap.Logger) grpcpkg.UnaryServerInterceptor {
	if log == nil {
		log = zap.NewNop()
	}

	return func(ctx context.Context, req any, info *grpcpkg.UnaryServerInfo, handler grpcpkg.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		log.Info("grpc request completed",
			zap.String("method", info.FullMethod),
			zap.String("code", status.Code(err).String()),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
		return resp, err
	}
}
```

- [ ] **Step 5: Add gRPC handler mapping tests**

Create `apps/api/internal/grpc/demo_server_test.go`:

```go
package grpc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apigrpc "monorepo-template/apps/api/internal/grpc"
	"monorepo-template/apps/api/internal/service/demo"
	demov1 "monorepo-template/libs/go/demoapi/gen/demo/v1"
)

type fakeDemoService struct {
	greeting demo.GetGreetingResult
	enqueue  demo.EnqueueEchoResult
	err      error
}

func (f *fakeDemoService) GetGreeting(ctx context.Context, input demo.GetGreetingInput) (demo.GetGreetingResult, error) {
	if f.err != nil {
		return demo.GetGreetingResult{}, f.err
	}
	return f.greeting, nil
}

func (f *fakeDemoService) EnqueueEcho(ctx context.Context, input demo.EnqueueEchoInput) (demo.EnqueueEchoResult, error) {
	if f.err != nil {
		return demo.EnqueueEchoResult{}, f.err
	}
	return f.enqueue, nil
}

func TestDemoServerGetGreetingSuccess(t *testing.T) {
	server := apigrpc.NewDemoServer(&fakeDemoService{greeting: demo.GetGreetingResult{Message: "Hello!", RequestID: "req-1"}})

	response, err := server.GetGreeting(context.Background(), &demov1.GetGreetingRequest{Name: "Template"})

	require.NoError(t, err)
	require.Equal(t, "Hello!", response.GetMessage())
	require.Equal(t, "req-1", response.GetRequestId())
}

func TestDemoServerMapsInvalidInput(t *testing.T) {
	server := apigrpc.NewDemoServer(&fakeDemoService{err: demo.ErrInvalidInput})

	_, err := server.EnqueueEcho(context.Background(), &demov1.EnqueueEchoRequest{Message: ""})

	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestDemoServerMapsUnexpectedFailure(t *testing.T) {
	server := apigrpc.NewDemoServer(&fakeDemoService{err: errors.New("database unavailable")})

	_, err := server.EnqueueEcho(context.Background(), &demov1.EnqueueEchoRequest{Message: "hello"})

	require.Equal(t, codes.Internal, status.Code(err))
	require.NotContains(t, err.Error(), "database unavailable")
}
```

- [ ] **Step 6: Implement demo gRPC handler**

Create `apps/api/internal/grpc/demo_server.go`:

```go
package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"monorepo-template/apps/api/internal/service/demo"
	demov1 "monorepo-template/libs/go/demoapi/gen/demo/v1"
)

type DemoService interface {
	GetGreeting(ctx context.Context, input demo.GetGreetingInput) (demo.GetGreetingResult, error)
	EnqueueEcho(ctx context.Context, input demo.EnqueueEchoInput) (demo.EnqueueEchoResult, error)
}

type DemoServer struct {
	demov1.UnimplementedDemoServiceServer
	service DemoService
}

func NewDemoServer(service DemoService) *DemoServer {
	return &DemoServer{service: service}
}

func (s *DemoServer) GetGreeting(ctx context.Context, req *demov1.GetGreetingRequest) (*demov1.GetGreetingResponse, error) {
	result, err := s.service.GetGreeting(ctx, demo.GetGreetingInput{Name: req.GetName()})
	if err != nil {
		return nil, mapDemoError(err)
	}

	return &demov1.GetGreetingResponse{
		Message:   result.Message,
		RequestId: result.RequestID,
	}, nil
}

func (s *DemoServer) EnqueueEcho(ctx context.Context, req *demov1.EnqueueEchoRequest) (*demov1.EnqueueEchoResponse, error) {
	result, err := s.service.EnqueueEcho(ctx, demo.EnqueueEchoInput{Message: req.GetMessage()})
	if err != nil {
		return nil, mapDemoError(err)
	}

	return &demov1.EnqueueEchoResponse{
		TaskId:    result.TaskID,
		Queue:     result.Queue,
		RequestId: result.RequestID,
	}, nil
}

func mapDemoError(err error) error {
	if errors.Is(err, demo.ErrInvalidInput) {
		return status.Error(codes.InvalidArgument, "invalid demo input")
	}
	return status.Error(codes.Internal, "internal demo service error")
}
```

- [ ] **Step 7: Add server construction tests**

Create `apps/api/internal/grpc/server_test.go`:

```go
package grpc_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	apigrpc "monorepo-template/apps/api/internal/grpc"
	"monorepo-template/apps/api/internal/service/demo"
	demov1 "monorepo-template/libs/go/demoapi/gen/demo/v1"
)

func TestNewServerRegistersDemoServiceAndAuth(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer lis.Close()

	server := apigrpc.NewServer(apigrpc.NewDemoServer(&fakeDemoService{greeting: demo.GetGreetingResult{Message: "Hello!", RequestID: "req-1"}}), "secret", zap.NewNop())
	defer server.Stop()

	go func() {
		_ = server.Serve(lis)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := demov1.NewDemoServiceClient(conn)
	ctx = metadata.AppendToOutgoingContext(ctx, apigrpc.InternalTokenMetadataKey, "secret")
	response, err := client.GetGreeting(ctx, &demov1.GetGreetingRequest{Name: "Template"})

	require.NoError(t, err)
	require.Equal(t, "Hello!", response.GetMessage())
}
```

- [ ] **Step 8: Implement server construction**

Create `apps/api/internal/grpc/server.go`:

```go
package grpc

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/zap"
	grpcpkg "google.golang.org/grpc"

	demov1 "monorepo-template/libs/go/demoapi/gen/demo/v1"
)

func NewServer(demoServer demov1.DemoServiceServer, authToken string, log *zap.Logger) *grpcpkg.Server {
	server := grpcpkg.NewServer(grpcpkg.ChainUnaryInterceptor(
		LoggingUnaryInterceptor(log),
		AuthUnaryInterceptor(authToken),
	))
	demov1.RegisterDemoServiceServer(server, demoServer)
	return server
}

func Run(ctx context.Context, port int, server *grpcpkg.Server) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		server.GracefulStop()
		return nil
	case err := <-errCh:
		return err
	}
}
```

- [ ] **Step 9: Wire gRPC and queue into API main**

Modify `apps/api/cmd/server/main.go`:

```go
import (
	// existing imports
	"sync"

	apigrpc "monorepo-template/apps/api/internal/grpc"
	demoservice "monorepo-template/apps/api/internal/service/demo"
	"monorepo-template/libs/go/queue"
)
```

Update the API config load call so `INTERNAL_GRPC_TOKEN` binds to `grpc.auth_token` before validation:

```go
cfg, err := config.Load[appconfig.Config](config.Options{
	ConfigPath: "config/config.yml",
	EnvFile:    optionalEnvFile(".env"),
	EnvAliases: map[string]string{
		"grpc.auth_token": "INTERNAL_GRPC_TOKEN",
	},
})
```

Add after Redis setup:

```go
var queueClient *queue.Client
if cfg.Queue.Enabled {
	queueClient = queue.NewClient(queue.Config{
		Redis:        cfg.Redis,
		DefaultQueue: cfg.Queue.DefaultQueue,
	})
	defer func() { _ = queueClient.Close() }()
}

demoService := demoservice.NewService(queueClient)
```

Replace the shutdown channel section with a shared application context:

```go
appCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
defer stop()

var wg sync.WaitGroup
if cfg.GRPC.Enabled {
	grpcServer := apigrpc.NewServer(apigrpc.NewDemoServer(demoService), cfg.GRPC.AuthToken, l)
	wg.Add(1)
	go func() {
		defer wg.Done()
		l.Info("starting internal grpc server", zap.Int("port", cfg.GRPC.Port))
		if err := apigrpc.Run(appCtx, cfg.GRPC.Port, grpcServer); err != nil {
			l.Fatal("internal grpc server failed", zap.Error(err))
		}
	}()
}

<-appCtx.Done()
```

Keep existing HTTP shutdown after `<-appCtx.Done()` and call `wg.Wait()` after `httpServer.Shutdown(ctx)`.

- [ ] **Step 10: Verify and commit**

Run: `bunx nx test api`

Expected: PASS including gRPC auth, handler, server, config, and service tests.

```bash
git add apps/api/go.mod apps/api/go.sum apps/api/cmd/server/main.go apps/api/internal/grpc apps/api/internal/appconfig apps/api/internal/service/demo apps/api/config/config.yml
git commit -m "feat: add api internal grpc server"
```

## Task 6: Bot gRPC Client And `/demo` Command

**Files:**

- Modify: `apps/bot/go.mod`
- Modify: `apps/bot/internal/appconfig/config.go`
- Modify: `apps/bot/config/config.yml`
- Create: `apps/bot/internal/appconfig/config_test.go`
- Create: `apps/bot/internal/apiclient/demo.go`
- Create: `apps/bot/internal/apiclient/demo_test.go`
- Create: `apps/bot/internal/handler/demo.go`
- Modify: `apps/bot/internal/handler/help.go`
- Modify: `apps/bot/internal/handler/handler_test.go`
- Modify: `apps/bot/cmd/bot/main.go`

- [ ] **Step 1: Add bot config tests**

Create `apps/bot/internal/appconfig/config_test.go`:

```go
package appconfig_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"monorepo-template/apps/bot/internal/appconfig"
	"monorepo-template/libs/go/config"
)

func TestInternalAPIConfigAllowsDisabledClient(t *testing.T) {
	cfg := appconfig.Config{
		Bot:         appconfig.BotConfig{Token: "token", PollTimeout: time.Second},
		Log:         config.LogConfig{Level: "debug", Format: "text"},
		InternalAPI: appconfig.InternalAPIConfig{Enabled: false},
	}

	require.NoError(t, config.Validate(cfg))
}

func TestInternalAPIConfigRequiresAddressTokenAndTimeoutWhenEnabled(t *testing.T) {
	cfg := appconfig.Config{
		Bot: appconfig.BotConfig{Token: "token", PollTimeout: time.Second},
		Log: config.LogConfig{Level: "debug", Format: "text"},
		InternalAPI: appconfig.InternalAPIConfig{
			Enabled: true,
		},
	}

	err := config.Validate(cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "GRPCAddress")
	require.Contains(t, err.Error(), "AuthToken")
	require.Contains(t, err.Error(), "Timeout")
}
```

- [ ] **Step 2: Extend bot config**

Modify `apps/bot/internal/appconfig/config.go`:

```go
package appconfig

import (
	"time"

	"monorepo-template/libs/go/config"
)

type Config struct {
	Bot         BotConfig         `mapstructure:"bot"`
	Log         config.LogConfig  `mapstructure:"log"`
	InternalAPI InternalAPIConfig `mapstructure:"internal_api"`
}

type BotConfig struct {
	Token       string        `mapstructure:"token" validate:"required"`
	PollTimeout time.Duration `mapstructure:"poll_timeout"`
}

type InternalAPIConfig struct {
	Enabled     bool          `mapstructure:"enabled"`
	GRPCAddress string        `mapstructure:"grpc_address" validate:"required_if=Enabled true"`
	AuthToken   string        `mapstructure:"auth_token" validate:"required_if=Enabled true"`
	Timeout     time.Duration `mapstructure:"timeout" validate:"required_if=Enabled true"`
}
```

Modify `apps/bot/config/config.yml`:

```yaml
internal_api:
  enabled: false
  grpc_address: localhost:9090
  auth_token: ''
  timeout: 3s
```

- [ ] **Step 3: Add bot API client tests**

Create `apps/bot/internal/apiclient/demo_test.go`:

```go
package apiclient_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"monorepo-template/apps/bot/internal/apiclient"
	"monorepo-template/libs/go/demoapi"
	demov1 "monorepo-template/libs/go/demoapi/gen/demo/v1"
)

type testDemoServer struct {
	demov1.UnimplementedDemoServiceServer
	seenToken string
	err       error
}

func (s *testDemoServer) GetGreeting(ctx context.Context, req *demov1.GetGreetingRequest) (*demov1.GetGreetingResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	values := md.Get(demoapi.InternalTokenMetadataKey)
	if len(values) == 1 {
		s.seenToken = values[0]
	}
	if s.err != nil {
		return nil, s.err
	}
	return &demov1.GetGreetingResponse{Message: "Hello, bot!", RequestId: "req-1"}, nil
}

func TestGRPCDemoClientSendsTokenAndReturnsGreeting(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer lis.Close()

	server := grpc.NewServer()
	fake := &testDemoServer{}
	demov1.RegisterDemoServiceServer(server, fake)
	defer server.Stop()
	go func() { _ = server.Serve(lis) }()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := apiclient.NewGRPCDemoClient(conn, "secret", time.Second)
	message, err := client.GetGreeting(context.Background(), "Template")

	require.NoError(t, err)
	require.Equal(t, "Hello, bot!", message)
	require.Equal(t, "secret", fake.seenToken)
}

func TestGRPCDemoClientMapsUnavailable(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer lis.Close()

	server := grpc.NewServer()
	demov1.RegisterDemoServiceServer(server, &testDemoServer{err: status.Error(codes.Unavailable, "down")})
	defer server.Stop()
	go func() { _ = server.Serve(lis) }()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := apiclient.NewGRPCDemoClient(conn, "secret", time.Second)
	_, err = client.GetGreeting(context.Background(), "Template")

	require.ErrorIs(t, err, apiclient.ErrUnavailable)
}
```

- [ ] **Step 4: Implement bot API client**

Create `apps/bot/internal/apiclient/demo.go`:

```go
package apiclient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"monorepo-template/apps/bot/internal/appconfig"
	"monorepo-template/libs/go/demoapi"
	demov1 "monorepo-template/libs/go/demoapi/gen/demo/v1"
)

var ErrUnavailable = errors.New("internal api unavailable")

type DemoClient interface {
	GetGreeting(ctx context.Context, name string) (string, error)
}

type GRPCDemoClient struct {
	client    demov1.DemoServiceClient
	authToken string
	timeout   time.Duration
}

type CloserDemoClient struct {
	DemoClient
	closer io.Closer
}

func DialDemoClient(cfg appconfig.InternalAPIConfig) (*CloserDemoClient, error) {
	conn, err := grpc.NewClient(cfg.GRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &CloserDemoClient{
		DemoClient: NewGRPCDemoClient(conn, cfg.AuthToken, cfg.Timeout),
		closer:    conn,
	}, nil
}

func (c *CloserDemoClient) Close() error {
	return c.closer.Close()
}

func NewGRPCDemoClient(conn grpc.ClientConnInterface, authToken string, timeout time.Duration) *GRPCDemoClient {
	return &GRPCDemoClient{
		client:    demov1.NewDemoServiceClient(conn),
		authToken: authToken,
		timeout:   timeout,
	}
}

func (c *GRPCDemoClient) GetGreeting(ctx context.Context, name string) (string, error) {
	callCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	callCtx = metadata.AppendToOutgoingContext(callCtx, demoapi.InternalTokenMetadataKey, c.authToken)

	response, err := c.client.GetGreeting(callCtx, &demov1.GetGreetingRequest{Name: name})
	if err != nil {
		return "", mapTransportError(err)
	}
	if response.GetMessage() == "" {
		return "", ErrUnavailable
	}
	return response.GetMessage(), nil
}

func mapTransportError(err error) error {
	code := status.Code(err)
	if code == codes.Unavailable || code == codes.DeadlineExceeded || errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%w: %v", ErrUnavailable, err)
	}
	return err
}

type DisabledDemoClient struct{}

func (DisabledDemoClient) GetGreeting(ctx context.Context, name string) (string, error) {
	return "", ErrUnavailable
}
```

Update `apps/bot/go.mod` direct requirements:

```go
google.golang.org/grpc v1.81.1
google.golang.org/protobuf v1.36.11
monorepo-template/libs/go/demoapi v0.0.0
```

Update `apps/bot/go.mod` replace block:

```go
replace (
	monorepo-template/libs/go/config => ../../libs/go/config
	monorepo-template/libs/go/demoapi => ../../libs/go/demoapi
	monorepo-template/libs/go/logger => ../../libs/go/logger
)
```

- [ ] **Step 5: Add bot `/demo` handler tests**

Append to `apps/bot/internal/handler/handler_test.go`:

```go
type fakeDemoClient struct {
	message string
	err     error
}

func (f fakeDemoClient) GetGreeting(ctx context.Context, name string) (string, error) {
	return f.message, f.err
}

func TestDemo_SendsGreeting(t *testing.T) {
	s := &mockSender{}
	h := handler.Demo(s, fakeDemoClient{message: "Hello, bot!"})

	h(context.Background(), nil, &models.Update{Message: &models.Message{Chat: models.Chat{ID: 789}}})

	require.NotNil(t, s.lastParams)
	assert.Equal(t, int64(789), s.lastParams.ChatID)
	assert.Contains(t, s.lastParams.Text, "Hello, bot!")
}

func TestDemo_SendsUnavailableFallback(t *testing.T) {
	s := &mockSender{}
	h := handler.Demo(s, fakeDemoClient{err: errors.New("down")})

	h(context.Background(), nil, &models.Update{Message: &models.Message{Chat: models.Chat{ID: 789}}})

	require.NotNil(t, s.lastParams)
	assert.Contains(t, s.lastParams.Text, "Demo service is unavailable")
}
```

- [ ] **Step 6: Implement `/demo` handler and register it**

Create `apps/bot/internal/handler/demo.go`:

```go
package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"monorepo-template/apps/bot/internal/apiclient"
	"monorepo-template/apps/bot/internal/botapi"
	"monorepo-template/libs/go/logger"
)

func Demo(s botapi.Sender, client apiclient.DemoClient) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		const op = "handler.Demo"
		log := logger.FromContext(ctx).With(zap.String("op", op))
		log.Debug("handling /demo")

		text, err := client.GetGreeting(ctx, "Telegram")
		if err != nil {
			log.Warn("internal demo api unavailable", zap.Error(err))
			text = "Demo service is unavailable. Please try again later."
		}

		if _, err := s.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text,
		}); err != nil {
			log.Error("failed to send message", zap.Error(err))
		}
	}
}
```

Modify help text in `apps/bot/internal/handler/help.go`:

```go
Text: "Available commands:\n/start - Start the bot\n/help - Show this message\n/demo - Call the internal API demo service",
```

Modify `apps/bot/cmd/bot/main.go`:

```go
cfg, err := config.Load[appconfig.Config](config.Options{
	ConfigPath: "config/config.yml",
	EnvAliases: map[string]string{
		"internal_api.auth_token": "INTERNAL_GRPC_TOKEN",
	},
})
```

Add client initialization after logger setup:

```go
demoClient := apiclient.DemoClient(apiclient.DisabledDemoClient{})
if cfg.InternalAPI.Enabled {
	client, err := apiclient.DialDemoClient(cfg.InternalAPI)
	if err != nil {
		log.Fatal("failed to create internal api client", zap.Error(err))
	}
	defer func() { _ = client.Close() }()
	demoClient = client
}
```

Keep the existing `bot.New(...)` block after the client initialization. Add `/demo` beside the existing `/start` and `/help` registrations after `b` has been created:

```go
b.RegisterHandler(bot.HandlerTypeMessageText, "/demo", bot.MatchTypeExact, handler.Demo(b, demoClient))
```

Add imports:

```go
"monorepo-template/apps/bot/internal/apiclient"
```

- [ ] **Step 7: Verify and commit**

Run: `bunx nx test bot`

Expected: PASS including appconfig, apiclient, and handler tests.

```bash
git add apps/bot/go.mod apps/bot/go.sum apps/bot/config/config.yml apps/bot/cmd/bot/main.go apps/bot/internal/appconfig apps/bot/internal/apiclient apps/bot/internal/handler
git commit -m "feat: add bot internal api demo command"
```

## Task 7: Worker App And Demo Echo Handler

**Files:**

- Create: `apps/worker/go.mod`
- Create: `apps/worker/config/config.yml`
- Create: `apps/worker/internal/appconfig/config.go`
- Create: `apps/worker/internal/appconfig/config_test.go`
- Create: `apps/worker/internal/handler/demo_echo.go`
- Create: `apps/worker/internal/handler/demo_echo_test.go`
- Create: `apps/worker/cmd/worker/main.go`
- Create: `apps/worker/project.json`
- Create: `apps/worker/air.toml`
- Modify: `go.work`

- [ ] **Step 1: Add worker config tests**

Create `apps/worker/internal/appconfig/config_test.go`:

```go
package appconfig_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"monorepo-template/apps/worker/internal/appconfig"
	"monorepo-template/libs/go/config"
)

func TestWorkerConfigRequiresQueueWeights(t *testing.T) {
	cfg := appconfig.Config{
		Redis: config.RedisConfig{Host: "localhost", Port: 6379},
		Queue: appconfig.QueueConfig{
			Concurrency: 5,
			Queues:      map[string]int{},
		},
		Log: config.LogConfig{Level: "info", Format: "json"},
	}

	err := config.Validate(cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Queues")
}

func TestWorkerConfigAcceptsRedisQueueAndLog(t *testing.T) {
	cfg := appconfig.Config{
		Redis: config.RedisConfig{Host: "localhost", Port: 6379},
		Queue: appconfig.QueueConfig{
			Concurrency: 5,
			Queues:      map[string]int{"default": 1},
		},
		Log: config.LogConfig{Level: "info", Format: "json"},
	}

	require.NoError(t, config.Validate(cfg))
}
```

- [ ] **Step 2: Add worker module and config**

Create `apps/worker/go.mod`:

```go
module monorepo-template/apps/worker

go 1.25.0

require (
	github.com/hibiken/asynq v0.26.0
	github.com/stretchr/testify v1.11.1
	go.uber.org/zap v1.27.1
	google.golang.org/grpc v1.81.1
	monorepo-template/libs/go/config v0.0.0
	monorepo-template/libs/go/demoapi v0.0.0
	monorepo-template/libs/go/logger v0.0.0
	monorepo-template/libs/go/queue v0.0.0
	monorepo-template/libs/go/tasks v0.0.0
)

replace (
	monorepo-template/libs/go/config => ../../libs/go/config
	monorepo-template/libs/go/demoapi => ../../libs/go/demoapi
	monorepo-template/libs/go/logger => ../../libs/go/logger
	monorepo-template/libs/go/queue => ../../libs/go/queue
	monorepo-template/libs/go/tasks => ../../libs/go/tasks
)
```

Create `apps/worker/config/config.yml`:

```yaml
redis:
  host: localhost
  port: 7502
  password: ''
  db: 0

queue:
  concurrency: 5
  queues:
    default: 1

log:
  level: info
  format: json
```

Create `apps/worker/internal/appconfig/config.go`:

```go
package appconfig

import "monorepo-template/libs/go/config"

type Config struct {
	Redis config.RedisConfig `mapstructure:"redis"`
	Queue QueueConfig        `mapstructure:"queue"`
	Log   config.LogConfig   `mapstructure:"log"`
}

type QueueConfig struct {
	Concurrency int            `mapstructure:"concurrency" validate:"gt=0"`
	Queues      map[string]int `mapstructure:"queues"      validate:"required,min=1,dive,keys,required,endkeys,gt=0"`
}
```

- [ ] **Step 3: Add worker handler tests**

Create `apps/worker/internal/handler/demo_echo_test.go`:

```go
package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"

	workerhandler "monorepo-template/apps/worker/internal/handler"
	"monorepo-template/libs/go/tasks"
)

type fakeProcessor struct {
	payload tasks.DemoEchoPayload
	err     error
}

func (f *fakeProcessor) ProcessDemoEcho(ctx context.Context, payload tasks.DemoEchoPayload) error {
	f.payload = payload
	return f.err
}

func TestDemoEchoHandlerProcessesValidTask(t *testing.T) {
	task, err := tasks.NewDemoEchoTask(tasks.DemoEchoPayload{Message: "hello", RequestID: "req-1"})
	require.NoError(t, err)
	processor := &fakeProcessor{}
	handler := workerhandler.NewDemoEchoHandler(processor)

	err = handler.Handle(context.Background(), task)

	require.NoError(t, err)
	require.Equal(t, "hello", processor.payload.Message)
	require.Equal(t, "req-1", processor.payload.RequestID)
}

func TestDemoEchoHandlerSkipsRetryForInvalidPayload(t *testing.T) {
	handler := workerhandler.NewDemoEchoHandler(&fakeProcessor{})

	err := handler.Handle(context.Background(), asynq.NewTask(tasks.TypeDemoEcho, []byte(`{`)))

	require.ErrorIs(t, err, asynq.SkipRetry)
	require.True(t, tasks.IsInvalidPayload(err))
}

func TestDemoEchoHandlerReturnsTransientErrorsForRetry(t *testing.T) {
	task, err := tasks.NewDemoEchoTask(tasks.DemoEchoPayload{Message: "hello"})
	require.NoError(t, err)
	handler := workerhandler.NewDemoEchoHandler(&fakeProcessor{err: errors.New("temporary store outage")})

	err = handler.Handle(context.Background(), task)

	require.ErrorContains(t, err, "temporary store outage")
	require.NotErrorIs(t, err, asynq.SkipRetry)
}
```

- [ ] **Step 4: Implement worker handler**

Create `apps/worker/internal/handler/demo_echo.go`:

```go
package handler

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"

	"monorepo-template/libs/go/tasks"
)

type DemoEchoProcessor interface {
	ProcessDemoEcho(ctx context.Context, payload tasks.DemoEchoPayload) error
}

type DemoEchoHandler struct {
	processor DemoEchoProcessor
}

type LoggingDemoEchoProcessor struct {
	Log *zap.Logger
}

func NewDemoEchoHandler(processor DemoEchoProcessor) *DemoEchoHandler {
	return &DemoEchoHandler{processor: processor}
}

func (h *DemoEchoHandler) Handle(ctx context.Context, task *asynq.Task) error {
	payload, err := tasks.ParseDemoEchoTask(task)
	if err != nil {
		return tasks.SkipRetryInvalidPayload(err)
	}
	if err := h.processor.ProcessDemoEcho(ctx, payload); err != nil {
		return fmt.Errorf("process demo echo: %w", err)
	}
	return nil
}

func (p LoggingDemoEchoProcessor) ProcessDemoEcho(ctx context.Context, payload tasks.DemoEchoPayload) error {
	p.Log.Info("processed demo echo task",
		zap.String("request_id", payload.RequestID),
		zap.String("message", payload.Message),
	)
	return nil
}
```

- [ ] **Step 5: Add worker entrypoint and Nx project**

Create `apps/worker/cmd/worker/main.go`:

```go
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"

	"monorepo-template/apps/worker/internal/appconfig"
	workerhandler "monorepo-template/apps/worker/internal/handler"
	"monorepo-template/libs/go/config"
	"monorepo-template/libs/go/logger"
	"monorepo-template/libs/go/queue"
	"monorepo-template/libs/go/tasks"
)

func main() {
	cfg, err := config.Load[appconfig.Config](config.Options{
		ConfigPath: "config/config.yml",
		EnvFile:    optionalEnvFile(".env"),
		EnvAliases: map[string]string{
			"queue.concurrency": "WORKER_CONCURRENCY",
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	log, err := logger.New(logger.Config{Level: cfg.Log.Level, Format: cfg.Log.Format})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = log.Sync() }()

	server := queue.NewServer(queue.Config{
		Redis:       cfg.Redis,
		Concurrency: cfg.Queue.Concurrency,
		Queues:      cfg.Queue.Queues,
	})

	mux := asynq.NewServeMux()
	processor := workerhandler.LoggingDemoEchoProcessor{Log: log}
	mux.HandleFunc(tasks.TypeDemoEcho, workerhandler.NewDemoEchoHandler(processor).Handle)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		log.Info("worker started", zap.Int("concurrency", cfg.Queue.Concurrency))
		errCh <- server.Run(mux)
	}()

	select {
	case <-ctx.Done():
		log.Info("worker stopping")
		server.Shutdown()
	case err := <-errCh:
		log.Fatal("worker failed", zap.Error(err))
	}
}

func optionalEnvFile(path string) string {
	if _, err := os.Stat(path); err != nil {
		return ""
	}
	return path
}
```

Create `apps/worker/project.json`:

```json
{
  "name": "worker",
  "$schema": "../../node_modules/nx/schemas/project-schema.json",
  "sourceRoot": "apps/worker",
  "projectType": "application",
  "tags": ["scope:worker", "lang:go"],
  "targets": {
    "build": {
      "executor": "nx:run-commands",
      "options": {
        "command": "cd apps/worker && go build -o ../../dist/apps/worker ./cmd/worker"
      }
    },
    "serve": {
      "executor": "nx:run-commands",
      "options": {
        "command": "cd apps/worker && air -c air.toml"
      }
    },
    "test": {
      "executor": "nx:run-commands",
      "options": {
        "command": "cd apps/worker && mkdir -p ../../dist/coverage/go/worker && go test -coverprofile=../../dist/coverage/go/worker/coverage.out ./..."
      }
    },
    "lint": {
      "executor": "nx:run-commands",
      "options": {
        "command": "cd apps/worker && golangci-lint run"
      }
    },
    "e2e": {
      "executor": "nx:run-commands",
      "options": {
        "command": "bash tools/e2e/worker-smoke.sh"
      },
      "dependsOn": ["build"]
    }
  }
}
```

Create `apps/worker/air.toml`:

```toml
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/worker ./cmd/worker"
bin = "./tmp/worker"
include_ext = ["go", "yml"]
exclude_dir = ["tmp"]
```

Modify `go.work`:

```go
use (
	./apps/api
	./apps/bot
	./apps/worker
	./libs/go/config
	./libs/go/demoapi
	./libs/go/logger
	./libs/go/queue
	./libs/go/tasks
)
```

- [ ] **Step 6: Verify and commit**

Run: `bunx nx test worker`

Expected: PASS including config and handler tests.

Run: `bunx nx build worker`

Expected: PASS.

```bash
git add go.work apps/worker
git commit -m "feat: add asynq worker app"
```

## Task 8: Docker, Compose, And E2E Harness

**Files:**

- Create: `docker/worker.Dockerfile`
- Modify: `docker/api.Dockerfile`
- Modify: `docker/bot.Dockerfile`
- Modify: `docker/docker-compose.dev.yml`
- Modify: `docker/docker-compose.yml`
- Modify: `docker/docker-compose.test.yml`
- Create: `tools/e2e/worker-smoke.sh`
- Create: `apps/worker/internal/integration/transport_smoke_test.go`

- [ ] **Step 1: Add the integration smoke test**

Create `apps/worker/internal/integration/transport_smoke_test.go`:

```go
//go:build integration

package integration_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"monorepo-template/libs/go/demoapi"
	demov1 "monorepo-template/libs/go/demoapi/gen/demo/v1"
)

func TestTransportSmoke(t *testing.T) {
	address := getenv("INTERNAL_API_GRPC_ADDRESS", "127.0.0.1:19090")
	authToken := getenv("INTERNAL_GRPC_TOKEN", "integration-token")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := demov1.NewDemoServiceClient(conn)
	ctx = metadata.AppendToOutgoingContext(ctx, demoapi.InternalTokenMetadataKey, authToken)

	greeting, err := client.GetGreeting(ctx, &demov1.GetGreetingRequest{Name: "Smoke"})
	require.NoError(t, err)
	require.Equal(t, "Hello, Smoke!", greeting.GetMessage())

	enqueued, err := client.EnqueueEcho(ctx, &demov1.EnqueueEchoRequest{Message: "hello worker"})
	require.NoError(t, err)
	require.NotEmpty(t, enqueued.GetTaskId())
	require.Equal(t, "default", enqueued.GetQueue())
	require.NotEmpty(t, enqueued.GetRequestId())
}

func getenv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
```

- [ ] **Step 2: Add the worker smoke runner**

Create `tools/e2e/worker-smoke.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

TEST_POSTGRES_PORT=17501 TEST_REDIS_PORT=17502 docker compose -f docker/docker-compose.test.yml up -d --wait postgres redis

tmpdir="$(mktemp -d)"
api_log="$tmpdir/api.log"
worker_log="$tmpdir/worker.log"
token="integration-token"

cleanup() {
  if [[ -n "${api_pid:-}" ]]; then kill "$api_pid" 2>/dev/null || true; fi
  if [[ -n "${worker_pid:-}" ]]; then kill "$worker_pid" 2>/dev/null || true; fi
  TEST_POSTGRES_PORT=17501 TEST_REDIS_PORT=17502 docker compose -f docker/docker-compose.test.yml down -v
  rm -rf "$tmpdir"
}
trap cleanup EXIT

wait_for_port() {
  local port="$1"
  for _ in {1..100}; do
    if (echo >"/dev/tcp/127.0.0.1/$port") >/dev/null 2>&1; then
      return 0
    fi
    sleep 0.1
  done
  echo "port $port did not open" >&2
  return 1
}

(
  cd apps/api
  SERVER_PORT=18090 \
  GRPC_ENABLED=true \
  GRPC_PORT=19090 \
  INTERNAL_GRPC_TOKEN="$token" \
  QUEUE_ENABLED=true \
  QUEUE_DEFAULT_QUEUE=default \
  POSTGRES_HOST=localhost \
  POSTGRES_PORT=17501 \
  POSTGRES_USER=app \
  POSTGRES_PASSWORD=secret \
  POSTGRES_DB=monorepo_test \
  POSTGRES_SSLMODE=disable \
  REDIS_HOST=localhost \
  REDIS_PORT=17502 \
  go run ./cmd/server
) >"$api_log" 2>&1 &
api_pid="$!"
wait_for_port 19090

(
  cd apps/worker
  REDIS_HOST=localhost \
  REDIS_PORT=17502 \
  WORKER_CONCURRENCY=1 \
  go run ./cmd/worker
) >"$worker_log" 2>&1 &
worker_pid="$!"

cd apps/worker
INTERNAL_API_GRPC_ADDRESS=127.0.0.1:19090 \
INTERNAL_GRPC_TOKEN="$token" \
go test -tags=integration -run TestTransportSmoke -count=1 -v ./internal/integration

for _ in {1..100}; do
  if grep -q "processed demo echo task" "$worker_log" && grep -q "hello worker" "$worker_log"; then
    exit 0
  fi
  sleep 0.1
done

echo "worker did not process demo.echo task" >&2
cat "$api_log" >&2
cat "$worker_log" >&2
exit 1
```

Run: `chmod +x tools/e2e/worker-smoke.sh`

- [ ] **Step 3: Add worker Dockerfile**

Create `docker/worker.Dockerfile`:

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /src

COPY go.work ./
COPY apps/api/go.mod apps/api/go.sum ./apps/api/
COPY apps/bot/go.mod apps/bot/go.sum ./apps/bot/
COPY apps/worker/go.mod apps/worker/go.sum ./apps/worker/
COPY libs/go/config/go.mod libs/go/config/go.sum ./libs/go/config/
COPY libs/go/demoapi/go.mod libs/go/demoapi/go.sum ./libs/go/demoapi/
COPY libs/go/logger/go.mod libs/go/logger/go.sum ./libs/go/logger/
COPY libs/go/queue/go.mod libs/go/queue/go.sum ./libs/go/queue/
COPY libs/go/tasks/go.mod libs/go/tasks/go.sum ./libs/go/tasks/

RUN cd apps/worker && go mod download

COPY apps/worker ./apps/worker
COPY libs/go/config ./libs/go/config
COPY libs/go/demoapi ./libs/go/demoapi
COPY libs/go/logger ./libs/go/logger
COPY libs/go/queue ./libs/go/queue
COPY libs/go/tasks ./libs/go/tasks

RUN cd apps/worker && go build -o /out/worker ./cmd/worker

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /out/worker /app/worker
COPY apps/worker/config/config.yml /app/config/config.yml
ENTRYPOINT ["/app/worker"]
```

- [ ] **Step 4: Update API and bot Dockerfiles for new Go modules**

In both dev and build stages of `docker/api.Dockerfile`, add module copy lines before `go mod download` so every `go.work` module path exists:

```dockerfile
COPY apps/worker/go.mod apps/worker/go.sum ./apps/worker/
COPY libs/go/demoapi/go.mod libs/go/demoapi/go.sum ./libs/go/demoapi/
COPY libs/go/queue/go.mod libs/go/queue/go.sum ./libs/go/queue/
COPY libs/go/tasks/go.mod libs/go/tasks/go.sum ./libs/go/tasks/
```

In both dev and build stages of `docker/api.Dockerfile`, add source copy lines before `WORKDIR /app/apps/api` or the API build:

```dockerfile
COPY libs/go/demoapi ./libs/go/demoapi
COPY libs/go/queue ./libs/go/queue
COPY libs/go/tasks ./libs/go/tasks
```

In both dev and build stages of `docker/bot.Dockerfile`, add module copy lines before `go mod download` so every `go.work` module path exists:

```dockerfile
COPY apps/worker/go.mod apps/worker/go.sum ./apps/worker/
COPY libs/go/demoapi/go.mod libs/go/demoapi/go.sum ./libs/go/demoapi/
COPY libs/go/queue/go.mod libs/go/queue/go.sum ./libs/go/queue/
COPY libs/go/tasks/go.mod libs/go/tasks/go.sum ./libs/go/tasks/
```

In both dev and build stages of `docker/bot.Dockerfile`, add source copy lines before `WORKDIR /app/apps/bot` or the bot build:

```dockerfile
COPY libs/go/demoapi ./libs/go/demoapi
```

- [ ] **Step 5: Update compose files**

In `docker/docker-compose.dev.yml`, ensure these service settings exist:

```yaml
api:
  expose:
    - '9090'
  environment:
    GRPC_ENABLED: 'true'
    GRPC_PORT: '9090'
    INTERNAL_GRPC_TOKEN: '${INTERNAL_GRPC_TOKEN:-dev-internal-token}'
    QUEUE_ENABLED: 'true'
    QUEUE_DEFAULT_QUEUE: default

bot:
  environment:
    INTERNAL_API_ENABLED: 'true'
    INTERNAL_API_GRPC_ADDRESS: api:9090
    INTERNAL_GRPC_TOKEN: '${INTERNAL_GRPC_TOKEN:-dev-internal-token}'
    INTERNAL_API_TIMEOUT: 3s

worker:
  build:
    context: ..
    dockerfile: docker/worker.Dockerfile
  depends_on:
    - redis
  environment:
    REDIS_HOST: redis
    REDIS_PORT: '6379'
    WORKER_CONCURRENCY: '5'
```

In `docker/docker-compose.yml`, expose API gRPC inside the compose network and add worker without host gRPC ports:

```yaml
api:
  expose:
    - '9090'
  environment:
    GRPC_ENABLED: 'true'
    GRPC_PORT: '9090'
    INTERNAL_GRPC_TOKEN: '${INTERNAL_GRPC_TOKEN}'
    QUEUE_ENABLED: 'true'
    QUEUE_DEFAULT_QUEUE: default

worker:
  build:
    context: ..
    dockerfile: docker/worker.Dockerfile
  depends_on:
    - redis
  environment:
    REDIS_HOST: redis
    REDIS_PORT: '6379'
    WORKER_CONCURRENCY: '5'
```

In `docker/docker-compose.test.yml`, ensure Redis binds to the existing test port:

```yaml
redis:
  image: redis:7-alpine
  ports:
    - '${TEST_REDIS_PORT:-17502}:6379'
```

- [ ] **Step 6: Verify and commit**

Run: `bunx nx run worker:e2e`

Expected: PASS with `TestTransportSmoke`.

```bash
git add docker tools/e2e/worker-smoke.sh apps/worker/internal/integration
git commit -m "feat: add worker transport smoke"
```

## Task 9: Coverage And Root Gates

**Files:**

- Modify: `tools/coverage/coverage.config.json`
- Modify: `tools/coverage/preflight.mjs`
- Create: `tools/codegen/check-drift.sh`
- Modify: `package.json`

- [ ] **Step 1: Add coverage config entries**

Modify `tools/coverage/coverage.config.json` by adding these `goProjects` entries:

```json
{
  "name": "go-demoapi",
  "cwd": "libs/go/demoapi",
  "profile": "dist/coverage/go/go-demoapi/coverage.out",
  "packages": "./..."
},
{
  "name": "go-tasks",
  "cwd": "libs/go/tasks",
  "profile": "dist/coverage/go/go-tasks/coverage.out",
  "packages": "./..."
},
{
  "name": "go-queue",
  "cwd": "libs/go/queue",
  "profile": "dist/coverage/go/go-queue/coverage.out",
  "packages": "./..."
},
{
  "name": "worker",
  "cwd": "apps/worker",
  "profile": "dist/coverage/go/worker/coverage.out",
  "packages": "./..."
}
```

Add exact allowlist entries:

```json
{
  "path": "libs/go/demoapi/gen/demo/v1/demo.pb.go",
  "reason": "protobuf generated Go message code",
  "gate": "bunx nx run go-demoapi:codegen && bunx nx build go-demoapi"
},
{
  "path": "libs/go/demoapi/gen/demo/v1/demo_grpc.pb.go",
  "reason": "protobuf generated Go gRPC transport code",
  "gate": "bunx nx run go-demoapi:codegen && bunx nx build go-demoapi"
},
{
  "path": "apps/worker/cmd/worker/main.go",
  "reason": "worker process bootstrap",
  "gate": "bunx nx build worker && bunx nx run worker:e2e"
}
```

- [ ] **Step 2: Add a committed-codegen drift gate**

Create `tools/codegen/check-drift.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

bun run codegen
git diff --exit-code -- \
  apps/api/internal/graph \
  apps/web-admin/src/shared/api/generated \
  libs/go/demoapi/gen/demo/v1
```

Run: `chmod +x tools/codegen/check-drift.sh`

Modify `package.json` scripts:

```json
"codegen:check": "bash tools/codegen/check-drift.sh"
```

Run: `bun run codegen:check`

Expected: PASS when generated files are current; FAIL with a non-zero `git diff --exit-code` when proto, GraphQL schema, or generated client files drift.

- [ ] **Step 3: Include worker smoke and codegen drift in root verification**

Modify `package.json`:

```json
"verify:coverage": "bun run lint && bun run codegen:check && bunx nx run web:typecheck && bunx nx run web-admin:typecheck && bun run build && bun run test:coverage && bun run test:e2e && bunx nx run worker:e2e && xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml && grace lint --path ."
```

Modify `tools/coverage/preflight.mjs` by adding these entries to `requiredFiles`:

```js
'apps/worker/project.json',
'apps/worker/config/config.yml',
'docker/worker.Dockerfile',
'tools/e2e/worker-smoke.sh',
```

The expected assertion is that preflight fails when the worker project, config, Dockerfile, or smoke runner is missing.

- [ ] **Step 4: Verify and commit**

Run: `bun run test:coverage`

Expected: PASS with new Go project profiles generated or validated.

```bash
git add package.json tools/coverage tools/codegen/check-drift.sh
git commit -m "chore: include transport slice in coverage gates"
```

## Task 10: CI, Dokploy, And Operator Docs

**Files:**

- Modify: `tools/ci/src/core.ts`
- Modify: `tools/ci/src/core.test.ts`
- Modify: `tools/ci/src/cli.ts`
- Modify: `tools/ci/src/cli.test.ts`
- Modify: `.gitlab-ci.yml`
- Modify: `deploy/dokploy/docker-compose.template.yml`
- Modify: `docs/infrastructure/ci-cd.md`

- [ ] **Step 1: Add CI helper tests for worker service metadata**

Modify `tools/ci/src/core.test.ts` so service image tests expect worker:

```ts
expect(buildImageRefs('registry.example.com/group/app', 'v1.2.3')).toEqual({
  api: 'registry.example.com/group/app/api:v1.2.3',
  web: 'registry.example.com/group/app/web:v1.2.3',
  bot: 'registry.example.com/group/app/bot:v1.2.3',
  worker: 'registry.example.com/group/app/worker:v1.2.3',
});
```

Modify metadata tests so `digests` includes:

```ts
worker: 'sha256:worker',
```

Add an assertion:

```ts
expect(metadata.services).toContainEqual(
  expect.objectContaining({
    service: 'worker',
    dockerfile: 'docker/worker.Dockerfile',
    image: 'registry.example.com/group/app/worker:v1.2.3',
    digest: 'sha256:worker',
  }),
);
```

- [ ] **Step 2: Update CI helper service list**

Modify `tools/ci/src/core.ts`:

```ts
export const services = [
  { name: 'api', dockerfile: 'docker/api.Dockerfile' },
  { name: 'web', dockerfile: 'docker/web.Dockerfile' },
  { name: 'bot', dockerfile: 'docker/bot.Dockerfile' },
  { name: 'worker', dockerfile: 'docker/worker.Dockerfile' },
] as const;
```

Modify Dokploy image env rendering so the returned object includes:

```ts
WORKER_IMAGE: refs.worker,
```

- [ ] **Step 3: Update CLI metadata tests and digest handling**

Modify `tools/ci/src/cli.test.ts` so `write-image-metadata` test env includes:

```ts
WORKER_IMAGE_DIGEST: 'sha256:worker',
```

Assert the written JSON contains a worker service entry:

```ts
expect(JSON.parse(writeFile.mock.calls[0][1]).services).toContainEqual(
  expect.objectContaining({
    service: 'worker',
    digest: 'sha256:worker',
  }),
);
```

Modify `tools/ci/src/cli.ts` inside the `write-image-metadata` command:

```ts
digests: {
  api: requireEnv(env, 'API_IMAGE_DIGEST'),
  web: requireEnv(env, 'WEB_IMAGE_DIGEST'),
  bot: requireEnv(env, 'BOT_IMAGE_DIGEST'),
  worker: requireEnv(env, 'WORKER_IMAGE_DIGEST'),
} satisfies Record<ServiceName, string>,
```

- [ ] **Step 4: Update GitLab pipeline**

Modify `build:images:dev` script in `.gitlab-ci.yml` by adding:

```bash
docker buildx build --pull --target prod -f docker/worker.Dockerfile -t "$CI_REGISTRY_IMAGE/worker:$IMAGE_TAG" -t "$CI_REGISTRY_IMAGE/worker:dev-latest" --metadata-file dist/ci/worker-image.json --push .
```

Modify release loops that currently use `api web bot`:

```bash
for service in api web bot worker; do
  # existing digest and metadata handling
done
```

Ensure the worker image build uses:

```bash
docker buildx build --pull --target prod -f docker/worker.Dockerfile -t "$CI_REGISTRY_IMAGE/worker:$IMAGE_TAG" --metadata-file dist/ci/worker-image.json --push .
export WORKER_IMAGE_DIGEST="$(jq -r '.["containerimage.digest"]' dist/ci/worker-image.json)"
```

- [ ] **Step 5: Update Dokploy compose template**

Modify `deploy/dokploy/docker-compose.template.yml`:

```yaml
services:
  api:
    expose:
      - '9090'
    environment:
      GRPC_ENABLED: 'true'
      GRPC_PORT: '9090'
      INTERNAL_GRPC_TOKEN: '${INTERNAL_GRPC_TOKEN}'
      QUEUE_ENABLED: 'true'
      QUEUE_DEFAULT_QUEUE: default

  bot:
    depends_on:
      - api
    environment:
      INTERNAL_API_ENABLED: 'true'
      INTERNAL_API_GRPC_ADDRESS: api:9090
      INTERNAL_GRPC_TOKEN: '${INTERNAL_GRPC_TOKEN}'
      INTERNAL_API_TIMEOUT: 3s

  worker:
    image: '${WORKER_IMAGE}'
    depends_on:
      - redis
    environment:
      REDIS_HOST: redis
      REDIS_PORT: '6379'
      WORKER_CONCURRENCY: '5'
```

- [ ] **Step 6: Update operator docs**

Modify `docs/infrastructure/ci-cd.md` with these facts:

```markdown
## Worker Image

The release pipeline builds and publishes `WORKER_IMAGE` beside `API_IMAGE`, `WEB_IMAGE`, and `BOT_IMAGE`. Release metadata records the worker image ref, digest, commit SHA, and pipeline ID.

## Internal API Transport

Dokploy compose exposes the API gRPC listener only inside the compose network through `expose: ["9090"]`. The bot uses `INTERNAL_API_GRPC_ADDRESS=api:9090` and both API and bot receive the shared `INTERNAL_GRPC_TOKEN`.

## Queue Worker

The worker process consumes Asynq tasks from Redis. Queue weights are configured in `apps/worker/config/config.yml`; `WORKER_CONCURRENCY` is the only worker queue environment override in this template.
```

- [ ] **Step 7: Verify and commit**

Run: `bunx nx test ci-tools`

Expected: PASS for CI helper tests.

Run: `bun run lint`

Expected: PASS.

```bash
git add .gitlab-ci.yml deploy/dokploy/docker-compose.template.yml docs/infrastructure/ci-cd.md tools/ci
git commit -m "feat: add worker image to ci deploy flow"
```

## Task 11: GRACE Contract Synchronization

**Files:**

- Modify: `docs/requirements.xml`
- Modify: `docs/technology.xml`
- Modify: `docs/development-plan.xml`
- Modify: `docs/knowledge-graph.xml`
- Modify: `docs/verification-plan.xml`
- Modify: `docs/operational-packets.xml`

- [ ] **Step 1: Update requirements**

In `docs/requirements.xml`, add this use case under `<UseCases>`:

```xml
<UC-008>
  <Actor>Developer</Actor>
  <Action>Uses the template API bot transport baseline for internal service calls and background jobs.</Action>
  <Goal>Provide a neutral reusable example of bot to API gRPC calls and API produced worker jobs.</Goal>
  <Preconditions>Redis is reachable, API internal gRPC is enabled, and the bot has a matching internal token.</Preconditions>
  <AcceptanceCriteria>Bot command `/demo` calls `DemoService.GetGreeting` over internal gRPC; `DemoService.EnqueueEcho` enqueues exactly one `demo.echo` task through Asynq over Redis; `GetGreeting` does not enqueue background work; users REST and GraphQL flows do not depend on gRPC or Asynq; API gRPC is internal-only in Docker and Dokploy compose files.</AcceptanceCriteria>
  <Priority>high</Priority>
  <RelatedFlows>DF-API-BOT-TRANSPORT</RelatedFlows>
</UC-008>
```

Add constraints under `<Constraints>` using the next available numeric ids:

```xml
<constraint-13>The API bot transport baseline uses Asynq over Redis and does not use Redis Streams, RabbitMQ, NATS, or Kafka.</constraint-13>
<constraint-14>The bot process does not run queue workers.</constraint-14>
<constraint-15>The internal gRPC listener is not exposed through public host ports.</constraint-15>
```

- [ ] **Step 2: Update technology**

In `docs/technology.xml`, add approved dependencies under `<Dependencies>`:

```xml
<dep name="google.golang.org/grpc" version="v1.81.1" purpose="Internal API to bot gRPC service contract" />
<dep name="google.golang.org/protobuf" version="v1.36.11" purpose="Generated demo service request and response types" />
<dep name="github.com/hibiken/asynq" version="v0.26.0" purpose="Redis-backed task queue for background jobs" />
<dep name="github.com/bufbuild/buf" version="v1.69.0" purpose="Proto code generation through go generate in libs/go/demoapi" />
```

Add tools under `<Tooling>`:

```xml
<tool name="proto-codegen" value="Buf" version="1.69.0" />
<tool name="queue-worker" value="Asynq worker process" version="0.26.0" />
```

- [ ] **Step 3: Update development plan modules**

In `docs/development-plan.xml`, add semantic module blocks under `<Modules>` using this repository's existing `M-*` tag style. The blocks can keep contracts compact, but they must use unique module tags, `NAME`, `TYPE`, `LAYER`, `ORDER`, `STATUS`, `<depends>`, `<target>`, and `<verification-ref>`:

```xml
<M-DEMOAPI NAME="DemoAPIContract" TYPE="UTILITY" LAYER="0" ORDER="3.5" STATUS="implemented">
  <contract><purpose>Own proto source, generated Go gRPC bindings, and shared internal metadata helpers for the neutral demo API.</purpose></contract>
  <interface><export-DemoService PURPOSE="Expose GetGreeting and EnqueueEcho generated gRPC contract." /></interface>
  <depends>none</depends>
  <target><source>libs/go/demoapi</source><tests>libs/go/demoapi/*_test.go</tests></target>
  <observability><log-prefix>[DemoAPI]</log-prefix><critical-block>BLOCK_GENERATE_DEMO_PROTO</critical-block></observability>
  <verification-ref>V-M-DEMOAPI</verification-ref>
</M-DEMOAPI>
<M-GO-TASKS NAME="GoTasks" TYPE="UTILITY" LAYER="0" ORDER="3.6" STATUS="implemented">
  <contract><purpose>Define Asynq task types, payloads, factories, parsers, and skip-retry helpers.</purpose></contract>
  <interface><export-TypeDemoEcho PURPOSE="Stable demo.echo task type." /><export-DemoEchoPayload PURPOSE="JSON task payload contract." /></interface>
  <depends>none</depends>
  <target><source>libs/go/tasks</source><tests>libs/go/tasks/*_test.go</tests></target>
  <observability><log-prefix>[GoTasks]</log-prefix><critical-block>BLOCK_VALIDATE_TASK_PAYLOAD</critical-block></observability>
  <verification-ref>V-M-GO-TASKS</verification-ref>
</M-GO-TASKS>
<M-GO-QUEUE NAME="GoQueue" TYPE="UTILITY" LAYER="0" ORDER="3.7" STATUS="implemented">
  <contract><purpose>Provide a thin Asynq adapter for Redis options, enqueue defaults, and worker server creation.</purpose></contract>
  <interface><export-NewClient PURPOSE="Construct Asynq client wrapper." /><export-NewServer PURPOSE="Construct Asynq server." /></interface>
  <depends>M-GO-CONFIG, M-GO-TASKS</depends>
  <target><source>libs/go/queue</source><tests>libs/go/queue/*_test.go</tests></target>
  <observability><log-prefix>[GoQueue]</log-prefix><critical-block>BLOCK_ENQUEUE_TASK</critical-block></observability>
  <verification-ref>V-M-GO-QUEUE</verification-ref>
</M-GO-QUEUE>
<M-API-DEMO-SERVICE NAME="APIDemoService" TYPE="CORE_LOGIC" LAYER="1" ORDER="5.1" STATUS="implemented">
  <contract><purpose>Own transport-neutral demo greeting and explicit echo enqueue behavior.</purpose></contract>
  <interface><export-GetGreeting PURPOSE="Return greeting without queue side effects." /><export-EnqueueEcho PURPOSE="Enqueue one demo.echo task." /></interface>
  <depends>M-GO-QUEUE, M-GO-TASKS</depends>
  <target><source>apps/api/internal/service/demo</source><tests>apps/api/internal/service/demo/*_test.go</tests></target>
  <observability><log-prefix>[APIDemoService]</log-prefix><critical-block>BLOCK_HANDLE_DEMO_SERVICE</critical-block></observability>
  <verification-ref>V-M-API-DEMO-SERVICE</verification-ref>
</M-API-DEMO-SERVICE>
<M-API-GRPC NAME="APIInternalGRPC" TYPE="ENTRY_POINT" LAYER="1" ORDER="5.2" STATUS="implemented">
  <contract><purpose>Serve the generated DemoService internally with auth and request logging interceptors.</purpose></contract>
  <interface><export-NewServer PURPOSE="Construct internal gRPC server." /><export-DemoServer PURPOSE="Map generated requests to demo service." /></interface>
  <depends>M-DEMOAPI, M-API-DEMO-SERVICE, M-GO-LOGGER</depends>
  <target><source>apps/api/internal/grpc</source><tests>apps/api/internal/grpc/*_test.go</tests></target>
  <observability><log-prefix>[API-GRPC]</log-prefix><critical-block>BLOCK_HANDLE_INTERNAL_GRPC</critical-block></observability>
  <verification-ref>V-M-API-GRPC</verification-ref>
</M-API-GRPC>
<M-BOT-API-CLIENT NAME="BotAPIClient" TYPE="UTILITY" LAYER="1" ORDER="7.1" STATUS="implemented">
  <contract><purpose>Call the generated DemoService from bot handlers with timeout, metadata auth, and friendly unavailable mapping.</purpose></contract>
  <interface><export-DemoClient PURPOSE="Local bot-facing demo client interface." /><export-GRPCDemoClient PURPOSE="Generated gRPC adapter." /></interface>
  <depends>M-DEMOAPI</depends>
  <target><source>apps/bot/internal/apiclient</source><tests>apps/bot/internal/apiclient/*_test.go</tests></target>
  <observability><log-prefix>[BotAPIClient]</log-prefix><critical-block>BLOCK_CALL_INTERNAL_API</critical-block></observability>
  <verification-ref>V-M-BOT-API-CLIENT</verification-ref>
</M-BOT-API-CLIENT>
<M-WORKER NAME="AsynqWorker" TYPE="ENTRY_POINT" LAYER="1" ORDER="8" STATUS="implemented">
  <contract><purpose>Run background Asynq handlers in a process separate from API and bot.</purpose></contract>
  <interface><export-main PURPOSE="Worker startup and graceful shutdown." /><export-DemoEchoHandler PURPOSE="Process demo.echo tasks." /></interface>
  <depends>M-GO-CONFIG, M-GO-LOGGER, M-GO-QUEUE, M-GO-TASKS</depends>
  <target><source>apps/worker</source><source>docker/worker.Dockerfile</source><tests>apps/worker/**/*_test.go</tests></target>
  <observability><log-prefix>[Worker]</log-prefix><critical-block>BLOCK_PROCESS_QUEUE_TASK</critical-block></observability>
  <verification-ref>V-M-WORKER</verification-ref>
</M-WORKER>
```

Update existing module dependency text exactly:

```text
M-COVERAGE-GATE depends -> M-WORKSPACE, M-GRAPHQL-SCHEMA, M-API, M-WEB-ADMIN, M-WEB, M-BOT, M-DEMOAPI, M-API-GRPC, M-API-DEMO-SERVICE, M-BOT-API-CLIENT, M-GO-TASKS, M-GO-QUEUE, M-WORKER
M-CI-CD depends -> M-WORKSPACE, M-COVERAGE-GATE, M-API, M-WEB-ADMIN, M-WEB, M-BOT, M-WORKER
M-API depends -> M-GO-CONFIG, M-GO-LOGGER, M-GRAPHQL-SCHEMA, M-API-GRPC, M-API-DEMO-SERVICE, M-GO-QUEUE, M-GO-TASKS
M-BOT depends -> M-GO-CONFIG, M-GO-LOGGER, M-BOT-API-CLIENT
```

Add a data flow under `<DataFlow>`:

```xml
<DF-API-BOT-TRANSPORT NAME="APIBotTransportFlow" TRIGGER="Bot demo command or internal enqueue call">
  <step-1>Bot `/demo` handler calls `M-BOT-API-CLIENT`.</step-1>
  <step-2>`M-BOT-API-CLIENT` sends generated `DemoService.GetGreeting` request with internal auth metadata from `M-DEMOAPI`.</step-2>
  <step-3>`M-API-GRPC` validates auth, logs the request without token values, and calls `M-API-DEMO-SERVICE`.</step-3>
  <step-4>`M-API-DEMO-SERVICE.GetGreeting` returns a greeting without queue side effects.</step-4>
  <step-5>`M-API-DEMO-SERVICE.EnqueueEcho` enqueues one `demo.echo` task through `M-GO-QUEUE` and `M-GO-TASKS`.</step-5>
  <step-6>`M-WORKER` consumes and processes the `demo.echo` task from Redis through Asynq.</step-6>
  <evidence>`bunx nx run worker:e2e`, focused module tests, and coverage gates pass.</evidence>
</DF-API-BOT-TRANSPORT>
```

- [ ] **Step 4: Update knowledge graph**

In `docs/knowledge-graph.xml`, add module entries under `<Project>` using the same tag style as existing graph modules:

```xml
<M-DEMOAPI NAME="DemoAPIContract" TYPE="UTILITY" STATUS="implemented">
  <purpose>Generated neutral demo gRPC contract and shared internal metadata helpers.</purpose>
  <path>libs/go/demoapi</path>
  <depends>none</depends>
  <verification-ref>V-M-DEMOAPI</verification-ref>
  <annotations><export-DemoService PURPOSE="Generated GetGreeting and EnqueueEcho contract." /></annotations>
</M-DEMOAPI>
<M-API-GRPC NAME="APIInternalGRPC" TYPE="ENTRY_POINT" STATUS="implemented">
  <purpose>Internal API gRPC server, auth interceptor, request logging interceptor, and demo handler.</purpose>
  <path>apps/api/internal/grpc</path>
  <depends>M-DEMOAPI, M-API-DEMO-SERVICE, M-GO-LOGGER</depends>
  <verification-ref>V-M-API-GRPC</verification-ref>
  <annotations><export-NewServer PURPOSE="Construct authenticated logging gRPC server." /></annotations>
</M-API-GRPC>
<M-API-DEMO-SERVICE NAME="APIDemoService" TYPE="CORE_LOGIC" STATUS="implemented">
  <purpose>Transport-neutral greeting and explicit demo echo enqueue service.</purpose>
  <path>apps/api/internal/service/demo</path>
  <depends>M-GO-QUEUE, M-GO-TASKS</depends>
  <verification-ref>V-M-API-DEMO-SERVICE</verification-ref>
  <annotations><export-EnqueueEcho PURPOSE="Create exactly one demo.echo task." /></annotations>
</M-API-DEMO-SERVICE>
<M-BOT-API-CLIENT NAME="BotAPIClient" TYPE="UTILITY" STATUS="implemented">
  <purpose>Bot-side internal API gRPC client adapter.</purpose>
  <path>apps/bot/internal/apiclient</path>
  <depends>M-DEMOAPI</depends>
  <verification-ref>V-M-BOT-API-CLIENT</verification-ref>
  <annotations><export-DemoClient PURPOSE="Bot-local demo client interface." /></annotations>
</M-BOT-API-CLIENT>
<M-GO-TASKS NAME="GoTasks" TYPE="UTILITY" STATUS="implemented">
  <purpose>Asynq task constants, JSON payloads, factories, and parse helpers.</purpose>
  <path>libs/go/tasks</path>
  <depends>none</depends>
  <verification-ref>V-M-GO-TASKS</verification-ref>
  <annotations><export-TypeDemoEcho PURPOSE="Stable demo.echo task type." /></annotations>
</M-GO-TASKS>
<M-GO-QUEUE NAME="GoQueue" TYPE="UTILITY" STATUS="implemented">
  <purpose>Thin Asynq client and server adapter over shared Redis config.</purpose>
  <path>libs/go/queue</path>
  <depends>M-GO-CONFIG, M-GO-TASKS</depends>
  <verification-ref>V-M-GO-QUEUE</verification-ref>
  <annotations><export-EnqueueDemoEcho PURPOSE="Enqueue demo echo task with default options." /></annotations>
</M-GO-QUEUE>
<M-WORKER NAME="AsynqWorker" TYPE="ENTRY_POINT" STATUS="implemented">
  <purpose>Separate background worker process for Asynq jobs.</purpose>
  <path>apps/worker</path>
  <path>docker/worker.Dockerfile</path>
  <depends>M-GO-CONFIG, M-GO-LOGGER, M-GO-QUEUE, M-GO-TASKS</depends>
  <verification-ref>V-M-WORKER</verification-ref>
  <annotations><export-DemoEchoHandler PURPOSE="Process demo.echo jobs." /></annotations>
</M-WORKER>
```

Add cross-links for all new dependencies and coverage links:

```xml
<CrossLink from="M-API" to="M-API-GRPC" relation="starts internal gRPC listener" />
<CrossLink from="M-API" to="M-API-DEMO-SERVICE" relation="uses neutral demo service" />
<CrossLink from="M-API" to="M-GO-QUEUE" relation="creates queue producer" />
<CrossLink from="M-API" to="M-GO-TASKS" relation="uses demo task contract" />
<CrossLink from="M-BOT" to="M-BOT-API-CLIENT" relation="calls internal API demo client" />
<CrossLink from="M-BOT-API-CLIENT" to="M-DEMOAPI" relation="imports generated demo contract" />
<CrossLink from="M-API-GRPC" to="M-DEMOAPI" relation="implements generated demo contract" />
<CrossLink from="M-API-GRPC" to="M-API-DEMO-SERVICE" relation="delegates demo behavior" />
<CrossLink from="M-API-DEMO-SERVICE" to="M-GO-QUEUE" relation="enqueues demo echo tasks" />
<CrossLink from="M-GO-QUEUE" to="M-GO-TASKS" relation="creates demo tasks" />
<CrossLink from="M-WORKER" to="M-GO-TASKS" relation="handles demo tasks" />
<CrossLink from="M-WORKER" to="M-GO-QUEUE" relation="uses Asynq server adapter" />
<CrossLink from="M-COVERAGE-GATE" to="M-DEMOAPI" relation="runs generated contract coverage and codegen drift gate" />
<CrossLink from="M-COVERAGE-GATE" to="M-API-GRPC" relation="runs API gRPC tests" />
<CrossLink from="M-COVERAGE-GATE" to="M-API-DEMO-SERVICE" relation="runs API demo service tests" />
<CrossLink from="M-COVERAGE-GATE" to="M-BOT-API-CLIENT" relation="runs bot internal API client tests" />
<CrossLink from="M-COVERAGE-GATE" to="M-GO-TASKS" relation="runs task contract tests" />
<CrossLink from="M-COVERAGE-GATE" to="M-GO-QUEUE" relation="runs queue adapter tests" />
<CrossLink from="M-COVERAGE-GATE" to="M-WORKER" relation="runs worker tests and worker:e2e" />
<CrossLink from="M-CI-CD" to="M-WORKER" relation="builds and deploys worker image" />
```

- [ ] **Step 5: Update verification plan**

In `docs/verification-plan.xml`, add verification refs under `<ModuleVerification>` using the current `V-M-*` tag style:

```xml
<V-M-DEMOAPI MODULE="M-DEMOAPI" PRIORITY="high">
  <test-files><file>libs/go/demoapi/contract_test.go</file><file>libs/go/demoapi/proto/demo/v1/demo.proto</file></test-files>
  <module-checks><check-1>bunx nx run go-demoapi:codegen</check-1><check-2>bunx nx test go-demoapi</check-2><check-3>bunx nx build go-demoapi</check-3><check-4>bun run codegen:check</check-4></module-checks>
  <scenarios><scenario-1 kind="success">Generated demo request and response types round-trip through protobuf.</scenario-1><scenario-2 kind="failure">Generated files drift after proto changes fails the codegen drift gate.</scenario-2></scenarios>
</V-M-DEMOAPI>
<V-M-API-GRPC MODULE="M-API-GRPC" PRIORITY="high">
  <test-files><file>apps/api/internal/grpc/*_test.go</file></test-files>
  <module-checks><check-1>bunx nx test api</check-1></module-checks>
  <scenarios><scenario-1 kind="success">Authenticated DemoService calls reach the demo service.</scenario-1><scenario-2 kind="failure">Missing or invalid token returns Unauthenticated without logging token values.</scenario-2></scenarios>
</V-M-API-GRPC>
<V-M-API-DEMO-SERVICE MODULE="M-API-DEMO-SERVICE" PRIORITY="high">
  <test-files><file>apps/api/internal/service/demo/service_test.go</file></test-files>
  <module-checks><check-1>bunx nx test api</check-1></module-checks>
  <scenarios><scenario-1 kind="success">GetGreeting returns a greeting without enqueuing.</scenario-1><scenario-2 kind="success">EnqueueEcho enqueues exactly one demo.echo task.</scenario-2></scenarios>
</V-M-API-DEMO-SERVICE>
<V-M-BOT-API-CLIENT MODULE="M-BOT-API-CLIENT" PRIORITY="high">
  <test-files><file>apps/bot/internal/apiclient/demo_test.go</file><file>apps/bot/internal/handler/handler_test.go</file></test-files>
  <module-checks><check-1>bunx nx test bot</check-1></module-checks>
  <scenarios><scenario-1 kind="success">/demo sends a greeting from internal API.</scenario-1><scenario-2 kind="failure">Unavailable or timed out internal API produces a friendly fallback.</scenario-2></scenarios>
</V-M-BOT-API-CLIENT>
<V-M-GO-TASKS MODULE="M-GO-TASKS" PRIORITY="high">
  <test-files><file>libs/go/tasks/tasks_test.go</file></test-files>
  <module-checks><check-1>bunx nx test go-tasks</check-1></module-checks>
  <scenarios><scenario-1 kind="success">Task factory and parser round-trip demo.echo JSON payloads.</scenario-1><scenario-2 kind="failure">Invalid payload wraps asynq.SkipRetry.</scenario-2></scenarios>
</V-M-GO-TASKS>
<V-M-GO-QUEUE MODULE="M-GO-QUEUE" PRIORITY="high">
  <test-files><file>libs/go/queue/queue_test.go</file></test-files>
  <module-checks><check-1>bunx nx test go-queue</check-1></module-checks>
  <scenarios><scenario-1 kind="success">Redis config and default enqueue options map to Asynq.</scenario-1></scenarios>
</V-M-GO-QUEUE>
<V-M-WORKER MODULE="M-WORKER" PRIORITY="high">
  <test-files><file>apps/worker/internal/handler/demo_echo_test.go</file><file>apps/worker/internal/integration/transport_smoke_test.go</file></test-files>
  <module-checks><check-1>bunx nx test worker</check-1><check-2>bunx nx build worker</check-2><check-3>bunx nx run worker:e2e</check-3></module-checks>
  <scenarios><scenario-1 kind="success">Worker processes one demo.echo task from Redis.</scenario-1><scenario-2 kind="failure">Invalid payload skips retry.</scenario-2></scenarios>
</V-M-WORKER>
```

Add `bunx nx run worker:e2e` to the broader release or coverage handoff gate that already includes `bun run verify:coverage`.

Update existing `V-M-COVERAGE-GATE` by adding these entries to the current sections:

```xml
<file>tools/codegen/check-drift.sh</file>
<file>apps/worker/internal/integration/transport_smoke_test.go</file>
<check-5>bunx nx run worker:e2e</check-5>
<entry path="libs/go/demoapi/gen/demo/v1/demo.pb.go" reason="protobuf generated Go message code" gate="bunx nx run go-demoapi:codegen &amp;&amp; bunx nx build go-demoapi" />
<entry path="libs/go/demoapi/gen/demo/v1/demo_grpc.pb.go" reason="protobuf generated Go gRPC transport code" gate="bunx nx run go-demoapi:codegen &amp;&amp; bunx nx build go-demoapi" />
<entry path="apps/worker/cmd/worker/main.go" reason="worker process bootstrap" gate="bunx nx build worker &amp;&amp; bunx nx run worker:e2e" />
```

Update existing `V-M-CI-CD` by adding these entries to the current sections:

```xml
<file>docker/worker.Dockerfile</file>
<scenario-6 kind="success">SemVer tag pipelines validate tag ancestry, publish api, web, bot, and worker release images, create release metadata, and expose manual production deploy.</scenario-6>
```

- [ ] **Step 6: Update operational packets with transport write-scope guidance**

In `docs/operational-packets.xml`, add this note under `<ExecutionPacketTemplate><ExecutionPacket><notes>`:

```xml
<note-4>API bot transport packets use these write scopes: M-DEMOAPI=libs/go/demoapi/**; M-API-GRPC=apps/api/internal/grpc/**; M-API-DEMO-SERVICE=apps/api/internal/service/demo/**; M-BOT-API-CLIENT=apps/bot/internal/apiclient/**; M-GO-TASKS=libs/go/tasks/**; M-GO-QUEUE=libs/go/queue/**; M-WORKER=apps/worker/** and docker/worker.Dockerfile.</note-4>
```

- [ ] **Step 7: Verify and commit**

Run:

```bash
xmllint --noout docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/verification-plan.xml docs/knowledge-graph.xml docs/operational-packets.xml
grace lint --path .
```

Expected: both commands PASS.

```bash
git add docs/requirements.xml docs/technology.xml docs/development-plan.xml docs/knowledge-graph.xml docs/verification-plan.xml docs/operational-packets.xml
git commit -m "docs: model api bot transport in grace"
```

## Task 12: End-to-End Verification And Drift Checks

**Files:**

- Read: entire changed set
- Modify only when verification exposes a concrete mismatch

- [ ] **Step 1: Run focused transport gates**

Run:

```bash
bunx nx run go-demoapi:codegen
bunx nx test go-demoapi
bunx nx build go-demoapi
bunx nx test api
bunx nx test bot
bunx nx test go-tasks
bunx nx test go-queue
bunx nx test worker
bunx nx build worker
bun run codegen:check
```

Expected: every command PASS.

- [ ] **Step 2: Check codegen drift**

Run: `bun run codegen:check`

Expected: PASS and no generated diff after codegen.

- [ ] **Step 3: Run integration smoke**

Run: `bunx nx run worker:e2e`

Expected: PASS with `TestTransportSmoke`.

- [ ] **Step 4: Run broad gates**

Run:

```bash
bun run lint
bun run test
bun run build
bun run test:coverage
bun run verify:coverage
```

Expected: every command PASS.

- [ ] **Step 5: Inspect security-sensitive logs and compose exposure**

Run:

```bash
rg -n "INTERNAL_GRPC_TOKEN|auth_token|x-internal-token" apps docker deploy docs
rg -n "9090:9090|ports:.*9090" docker deploy
```

Expected:

- `INTERNAL_GRPC_TOKEN` appears only in config, env binding, docs, tests, and compose env references.
- No log statement prints the token value.
- No production-ish compose file maps `9090:9090` through host `ports`.

- [ ] **Step 6: Commit verification fixes or close with clean status**

Run:

```bash
git status --short
```

Expected: clean except user-owned files that were dirty before execution and remain unrelated.

When Step 1 through Step 5 produced transport-owned fixes, stage them with the exact `git add ...` command from the task that owns the changed file, excluding pre-existing user-owned dirty files. Then commit:

```bash
git commit -m "fix: align api bot transport verification"
```

## Self-Review Result

Spec coverage:

- gRPC sync path is covered by Tasks 1, 4, 5, 6, 8, 11, and 12.
- API gRPC auth and request logging interceptors are covered by Task 5 tests.
- Explicit async queue path is covered by Tasks 2, 3, 4, 5, 7, 8, 9, 11, and 12.
- `GetGreeting` no-side-effect behavior is covered by Task 4 service tests.
- Asynq over Redis, not Redis Streams, is covered by Tasks 2, 3, 7, 8, and GRACE constraints in Task 11.
- Separate `apps/worker` ownership is covered by Tasks 7, 8, 10, and 11.
- Docker, Dokploy, CI metadata, worker image, coverage allowlist, and `worker:e2e` are covered by Tasks 8, 9, 10, and 12.
- `M-CI-CD`, existing `M-API`, existing `M-BOT`, and all `M-COVERAGE-GATE` cross-links are covered by Task 11.
- Generated code exact allowlist paths are covered by Task 9.

Placeholder scan:

- The plan does not use floating dependency versions.
- The plan does not use wildcard coverage allowlist entries for generated Go code.
- The plan names concrete files, commands, and expected outcomes.

Type consistency:

- `demoapi.InternalTokenMetadataKey`, `tasks.DemoEchoPayload`, `queue.EnqueueResult`, `demo.EchoEnqueuer`, `apigrpc.InternalTokenMetadataKey`, `apiclient.DemoClient`, and `workerhandler.DemoEchoProcessor` are introduced before later tasks use them.
- gRPC request and response field names match `demo.proto` and generated Go naming.
- Queue names use `default_queue` in API config, `queues` weights in worker config, and `WORKER_CONCURRENCY` as the only worker queue env override.
