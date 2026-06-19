# ---- Dev stage ----
FROM golang:1.25-alpine AS dev

RUN go install github.com/air-verse/air@latest
WORKDIR /app
COPY go.work ./
COPY apps/api/go.mod apps/api/go.sum ./apps/api/
COPY apps/bot/go.mod apps/bot/go.sum ./apps/bot/
COPY libs/go/config/go.mod libs/go/config/go.sum ./libs/go/config/
COPY libs/go/logger/go.mod libs/go/logger/go.sum ./libs/go/logger/
RUN cd apps/api && go mod download
COPY apps/api/ ./apps/api/
COPY libs/go/config/ ./libs/go/config/
COPY libs/go/logger/ ./libs/go/logger/
WORKDIR /app/apps/api
CMD ["air", "-c", "air.toml"]

# ---- Build stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.work ./
COPY apps/api/go.mod apps/api/go.sum ./apps/api/
COPY apps/bot/go.mod apps/bot/go.sum ./apps/bot/
COPY libs/go/config/go.mod libs/go/config/go.sum ./libs/go/config/
COPY libs/go/logger/go.mod libs/go/logger/go.sum ./libs/go/logger/
RUN cd apps/api && go mod download
COPY apps/api/ ./apps/api/
COPY libs/go/config/ ./libs/go/config/
COPY libs/go/logger/ ./libs/go/logger/
RUN cd apps/api && CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# ---- Production stage ----
FROM alpine:3.20 AS prod

RUN apk --no-cache add ca-certificates
COPY --from=builder /server /server
COPY apps/api/config/config.yml /config/config.yml
EXPOSE 8080
CMD ["/server"]
