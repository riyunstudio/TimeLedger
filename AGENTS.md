# AGENTS.md

This file provides guidance to agentic coding assistants working in this repository.

## Build/Lint/Test Commands

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

## Code Style Guidelines

### Import Organization
- Third-party imports grouped first, internal packages second
- No blank line between groups
- Use `timeLedger/<path>` format for internal imports
- Example:
  ```go
  import (
      "context"
      "encoding/json"
      "time"
      "timeLedger/app"
      "timeLedger/app/models"
      "timeLedger/app/services"

      "github.com/gin-gonic/gin"
  )
  ```

### Formatting
- Use standard `go fmt` conventions
- Struct tags use backticks with proper spacing
- No trailing whitespace
- Use tabs for indentation
- Max line length not strictly enforced but keep readable

### Types and Naming
- Package names: lowercase, single word (e.g., `controllers`, `services`)
- Struct names: PascalCase (e.g., `UserService`, `UserRepository`)
- Interface names: PascalCase, typically ending with type (e.g., `AuthService`)
- Method names: PascalCase for exported, lowercase for private
- Variable names: camelCase (e.g., `userRepository`, `errInfo`)
- Constants: UPPER_SNAKE_CASE (e.g., `SQL_ERROR`, `USER_NOT_FOUND`)
- Context variables: `ctx`
- Error variables: `err`
- Error info variables: `eInfo` or `errInfo`

### Error Handling
- Always use triple return: `(datas any, errInfo *errInfos.Res, err error)`
- Error codes defined in `global/errInfos/code.go` (format: Type(1) + Serial(4))
- Error messages in `global/errInfos/message.go`
- Use `s.app.Err.New(errInfos.ERROR_CODE)` to create error info
- Check errors immediately: `if err != nil { return nil, eInfo, err }`
- Never ignore errors unless intentionally

### Repository Pattern
- Read operations: `app.Mysql.RDB.WithContext(ctx)` (slave)
- Write operations: `app.Mysql.WDB.WithContext(ctx)` (master)
- Always pass `context.Context` as first parameter
- Repository methods return `(data Model, err error)` or `(Model, error)`
- Model structs in `app/models/` with GORM tags and JSON tags

### Request Validation
- Use `requests.Validate[T](ctx)` for automatic binding
- Add `binding:"required"` tag for required fields
- Request structs in `app/requests/` per entity
- Methods return `(req *T, errInfo *errInfos.Res, err error)`

### Service Layer
- Business logic in `app/services/`
- Services use repositories for data access
- Use resources (DTOs) in `app/resources/` for response formatting
- Context propagation: pass `ctx` to repository calls

### Controller Layer
- HTTP handlers in `app/controllers/` with Swagger annotations
- Use BaseController for common functionality
- Extract Gin context via `ctl.makeCtx(ctx)`
- Return responses via `ctl.JSON(ctx, global.Ret{...})`

### Testing
- Tests in `testing/test/` use SQLite mock DB + MinRedis
- Initialize: `sqliteDB, _ := sqlite.Initialize()`, `rdb, mr, _ := mockRedis.Initialize()`
- Inject: `app.Mysql = &mysql.DB{WDB: sqliteDB, RDB: sqliteDB}`
- Use table-driven tests with subtests
- Verify both success and error cases
- Test naming: `Test<Feature>_<Action>` (e.g., `TestUserService_CreateAndGet`)

### gRPC Services
- Proto files in `grpc/proto/` with `go_package = "./proto/<service>"`
- Implement in `grpc/services/` embedding `Unimplemented<Name>ServiceServer`
- Register in `grpc/server.go`
- Use app.App context, not global app

### Console/Cron Jobs
- Jobs in `app/console/` implement `Job` interface
- Use scheduler with cron expression: "秒 分 時 日 月 星期"
- Initialize TraceLog for logging
- Handle Telegram alerts for errors

### General Patterns
- Use `app.Tools` for utilities (timezone, IP validation, JSON merge, trace ID)
- Use `app.Api` for external HTTP calls
- Use `app.Rpc` for RPC calls
- Time fields stored as Unix timestamps (int64)
- JSON fields stored as strings in DB, unmarshaled in resources
- Use `defer` for cleanup (DB close, etc.)
- Recover panics in goroutines
- Always run lint and typecheck after changes

### File Organization
- Models: `app/models/<entity>.go`
- Requests: `app/requests/<entity>.go`
- Repositories: `app/repositories/<entity>.go`
- Resources: `app/resources/<entity>.go`
- Services: `app/services/<entity>.go`
- Controllers: `app/controllers/<entity>.go`
- gRPC proto: `grpc/proto/<service>.proto`
- gRPC impl: `grpc/services/<service>.go`
