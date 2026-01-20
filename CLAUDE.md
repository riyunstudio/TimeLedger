# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

TimeLedger is a Go backend microservice for a teacher-centric multi-center scheduling platform targeting Taiwan's LINE-first ecosystem. It provides HTTP REST API (Gin), gRPC services, WebSocket, and scheduled jobs. The frontend is a separate Nuxt 3 project.

## Common Commands

```bash
# Build
go build -mod=vendor -o main .

# Run locally (requires MySQL + Redis running)
go run main.go

# Run tests (uses SQLite mock DB + MinRedis)
go test ./testing/test/... -v

# Lint (via Docker build or local golangci-lint)
golangci-lint run --timeout 10m

# Generate Swagger docs
swag init

# Compile Protocol Buffers
protoc --go_out=./grpc --go-grpc_out=./grpc grpc/proto/<service>.proto

# Configure private Go modules
go env -w GOPRIVATE=gitlab.en.mcbwvx.com
```

## Environment Setup

Copy `.env.example` to `.env`. Key services:
- HTTP API: `localhost:8888` (Swagger at `/swagger/index.html`)
- gRPC: `localhost:50051`
- WebSocket: `localhost:8889`
- Health check: `/healthy`

MySQL master-slave replication: RDB (read/slave), WDB (write/master).

## Architecture

**Layered Architecture with Repository Pattern:**
```
HTTP Request → Middleware → Controller → Request (validation) → Service → Repository → Model
                                                                              ↓
gRPC Request → Interceptors → gRPC Service ─────────────────────────────────→┘
```

Key directories:
- `app/controllers/` - HTTP handlers (Gin)
- `app/services/` - Business logic
- `app/repositories/` - Data access (GORM)
- `app/models/` - Database models
- `app/requests/` - Request validation with `binding` tags
- `app/resources/` - Response DTOs
- `app/servers/` - HTTP server, routes, middleware
- `app/console/` - Cron jobs
- `grpc/proto/` - Protocol Buffer definitions
- `grpc/services/` - gRPC implementations
- `global/errInfos/` - Centralized error codes

## Key Patterns

**Database Operations:**
- Read: `app.Mysql.RDB.WithContext(ctx)` (slave)
- Write: `app.Mysql.WDB.WithContext(ctx)` (master)
- Always pass `context.Context` to repository methods

**Error Codes:** Format `ProjectID(2) + FunctionType(2) + Serial(4)`
- Type 1: System, Type 2: DB/Cache, Type 3: Other, Type 4: User

**Request Validation:** Use `requests.Validate[T](ctx)` generic function

**Adding New Endpoints:**
1. Model in `app/models/`
2. Request in `app/requests/`
3. Repository in `app/repositories/`
4. Resource in `app/resources/`
5. Service in `app/services/`
6. Controller in `app/controllers/` with Swagger annotations
7. Register route in `app/servers/route.go`

**gRPC Services:**
1. Define proto in `grpc/proto/`
2. Compile with protoc
3. Implement in `grpc/services/` (embed `UnimplementedXxxServiceServer`)
4. Register in `grpc/server.go` `registerServices()`

## Agent Skills (`.agent/skills/`)

- **auth-adapter-guard**: Mock Login vs LINE Login abstraction; use `AuthService` interface, never call `liff.*` directly
- **contract-sync**: Keep API specs in sync with Go structs and TypeScript interfaces; update models when changing `pdr/API.md` or `pdr/Mysql.md`
- **scheduling-validator**: TDD for scheduling engine; write tests first for overlap, buffer, and cross-day logic

## Documentation

Product specs in `pdr/`:
- `MASTER_PROMPT.md` - Development directives
- `API.md` - API specifications
- `Mysql.md` - Database schema
- `功能業務邏輯.md` - Business logic (Chinese)
- `Stages.md` - Execution roadmap

## Internal Packages

- `gitlab.en.mcbwvx.com/frame/teemo` - Tools (timezone, JSON utilities)
- `gitlab.en.mcbwvx.com/frame/zilean` - Logging
- `gitlab.en.mcbwvx.com/frame/ezreal` - HTTP client wrapper
