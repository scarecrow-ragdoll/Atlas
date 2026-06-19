# ---- Dev stage ----
FROM golang:1.25-alpine AS dev

RUN go install github.com/air-verse/air@latest
WORKDIR /app
COPY go.work ./
COPY apps/bot/go.mod apps/bot/go.sum ./apps/bot/
COPY apps/api/go.mod apps/api/go.sum ./apps/api/
COPY libs/go/config/go.mod libs/go/config/go.sum ./libs/go/config/
COPY libs/go/logger/go.mod libs/go/logger/go.sum ./libs/go/logger/
RUN cd apps/bot && go mod download
COPY apps/bot/ ./apps/bot/
COPY libs/go/config/ ./libs/go/config/
COPY libs/go/logger/ ./libs/go/logger/
WORKDIR /app/apps/bot
CMD ["air", "-c", "air.toml"]

# ---- Build stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.work ./
COPY apps/bot/go.mod apps/bot/go.sum ./apps/bot/
COPY apps/api/go.mod apps/api/go.sum ./apps/api/
COPY libs/go/config/go.mod libs/go/config/go.sum ./libs/go/config/
COPY libs/go/logger/go.mod libs/go/logger/go.sum ./libs/go/logger/
RUN cd apps/bot && go mod download
COPY apps/bot/ ./apps/bot/
COPY libs/go/config/ ./libs/go/config/
COPY libs/go/logger/ ./libs/go/logger/
RUN cd apps/bot && CGO_ENABLED=0 GOOS=linux go build -o /bot ./cmd/bot

# ---- Production stage ----
FROM alpine:3.20 AS prod

RUN apk --no-cache add ca-certificates
COPY --from=builder /bot /bot
COPY apps/bot/config/config.yml /config/config.yml
CMD ["/bot"]
