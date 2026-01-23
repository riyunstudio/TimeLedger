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

# Run all tests (uses SQLite mock DB + MinRedis)
go test ./testing/test/... -v

# Run a single test
go test ./testing/test/... -v -run TestUserService_CreateAndGet

# Lint
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
- `app/controllers/` - HTTP handlers (Gin) with Swagger annotations
- `app/services/` - Business logic
- `app/repositories/` - Data access (GORM)
- `app/models/` - Database models with GORM/JSON tags
- `app/requests/` - Request validation with `binding` tags
- `app/resources/` - Response DTOs
- `app/servers/` - HTTP server, routes, middleware
- `app/console/` - Cron jobs (Job interface)
- `grpc/proto/` - Protocol Buffer definitions
- `grpc/services/` - gRPC implementations
- `global/errInfos/` - Centralized error codes
- `testing/test/` - Unit/integration tests

## Code Style Guidelines

### Import Organization
```go
import (
    "context"
    "encoding/json"
    "time"

    "timeLedger/app"
    "timeLedger/app/models"
    "timeLedger/app/services"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)
```

### Naming Conventions
| Type | Convention | Example |
|------|------------|---------|
| Package | lowercase, single word | `controllers`, `services` |
| Struct | PascalCase | `UserService`, `AdminUserRepository` |
| Interface | PascalCase + type suffix | `AuthService`, `Job` |
| Method | PascalCase (exported), lowercase (private) | `CreateUser()`, `validate()` |
| Variable | camelCase | `userRepository`, `errInfo` |
| Constant | UPPER_SNAKE_CASE | `SQL_ERROR`, `USER_NOT_FOUND` |
| Context | `ctx` | - |
| Error | `err` | - |
| Error Info | `eInfo` or `errInfo` | - |

### Formatting
- Use tabs for indentation
- Struct tags with backticks and proper spacing
- No trailing whitespace
- Max line length: keep readable

## Key Patterns

### Database Operations
- **Read:** `app.Mysql.RDB.WithContext(ctx)` (slave)
- **Write:** `app.Mysql.WDB.WithContext(ctx)` (master)
- Always pass `context.Context` as first parameter

### Error Handling
```go
// Triple return pattern
func GetUser(ctx context.Context, id uint) (*models.User, *errInfos.Res, error) {
    user, err := repo.GetByID(ctx, id)
    if err != nil {
        return nil, s.app.Err.New(errInfos.USER_NOT_FOUND), err
    }
    return user, nil, nil
}
```
- Error codes: `FunctionType(1) + Serial(4)` (e.g., `10001` = System Error)
  - Type 1: System (10001-10999)
  - Type 2: DB/Cache (20001-20999)
  - Type 3: Other (30001-30999)
  - Type 4: User (40001-40999)
- Define codes in `global/errInfos/code.go`, messages in `message.go`

### Request Validation
```go
func Validate[T any](ctx *gin.Context) (*T, *errInfos.Res, error) {
    var req T
    if err := ctx.ShouldBindJSON(&req); err != nil {
        return nil, nil, err
    }
    return &req, nil, nil
}
```
- Use `binding:"required"` for required fields
- Request structs in `app/requests/<entity>.go`

### Service Layer
- Business logic in `app/services/`
- Use repositories for data access
- Pass `ctx` to all repository calls

### Controller Layer
- Extract Gin context: `ctl.makeCtx(ctx)`
- Return responses: `ctl.JSON(ctx, global.Ret{...})`

### gRPC Services
1. Define proto in `grpc/proto/` with `go_package`
2. Compile with `protoc`
3. Implement in `grpc/services/` embedding `Unimplemented<Name>ServiceServer`
4. Register in `grpc/server.go`

### General Patterns
- Time fields: Unix timestamps (`int64`)
- JSON fields: stored as strings in DB, unmarshaled in resources
- Use `defer` for cleanup
- Recover panics in goroutines
- Use `app.Tools` (timezone, IP, JSON, trace ID)
- Use `app.Api` for external HTTP calls
- Use `app.Rpc` for RPC calls

## Adding New Endpoints

1. Model → `app/models/<entity>.go`
2. Request → `app/requests/<entity>.go`
3. Repository → `app/repositories/<entity>.go`
4. Resource → `app/resources/<entity>.go`
5. Service → `app/services/<entity>.go`
6. Controller → `app/controllers/<entity>.go`
7. Register route → `app/servers/route.go`

## Agent Skills (`.agent/skills/`)

- **auth-adapter-guard**: Mock Login vs LINE Login abstraction; use `AuthService` interface, never call `liff.*` directly
- **contract-sync**: Keep API specs in sync with Go structs and TypeScript interfaces; update models when changing `pdr/API.md` or `pdr/Mysql.md`
- **scheduling-validator**: TDD for scheduling engine; write tests first for overlap, buffer, and cross-day logic

## Testing

Tests use SQLite mock DB + MinRedis:
```go
sqliteDB, _ := sqlite.Initialize()
rdb, mr, _ := mockRedis.Initialize()
defer mr.Close()

app.Mysql = &mysql.DB{WDB: sqliteDB, RDB: sqliteDB}
app.Redis = &redis.Redis{DB0: rdb}
```
- Use table-driven tests with subtests
- Test naming: `Test<Feature>_<Action>` (e.g., `TestUserService_CreateAndGet`)
- Verify both success and error cases

## Documentation

Specs in `pdr/`:
- `MASTER_PROMPT.md` - Development directives
- `API.md` - API specifications
- `Mysql.md` - Database schema
- `功能業務邏輯.md` - Business logic (Chinese)
- `Stages.md` - Execution roadmap

## Internal Packages

- `gitlab.en.mcbwvx.com/frame/teemo` - Tools (timezone, JSON utilities)
- `gitlab.en.mcbwvx.com/frame/zilean` - Logging
- `gitlab.en.mcbwvx.com/frame/ezreal` - HTTP client wrapper
