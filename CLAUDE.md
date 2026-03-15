# CLAUDE.md

TimeLedger 開發指引 - 精華版

> 詳細規格請參考 `pdr/` 目錄下的專門文件

---

## 1. 專案定位

- **名稱**：TimeLedger - 教師中心化多據點排課平台
- **市場**：台灣（LINE-First 生態系）
- **策略**：SaaS + 人才市場
- **後端**：Go (Gin) + MySQL 8.0 + Redis
- **前端**：Nuxt 3 (SSR) + Tailwind CSS + LINE LIFF
- **部署**：Docker Compose

---

## 2. 認證策略

### 教師端：LINE 單一登入
- **無密碼**：完全廢除帳密設定，LIFF Silent Login
- **LINE 綁定**：首次登入自動綁定，**不可解除**

### 管理員端
- **Email/Password + JWT**（24小時效期）
- 角色：OWNER、ADMIN、STAFF
- LINE 綁定：用於 Exception 通知（可綁定/解除）

---

## 3. 分層架構

```
HTTP Request → Middleware → Controller → Request → Service → Repository → Model
```

### 分層職責

| 層級 | 職責 | 禁止事項 |
|:---|:---|:---|
| Controller | Request 解析、呼叫 Service、回傳 JSON | 業務邏輯 |
| Service | 核心業務邏輯、狀態機流轉、交易控制 | 直接操作 DB |
| Repository | 純粹 DB 操作（Find/Create/Update） | 業務判斷 |
| Model | 數據結構定義（GORM Tags） | 商業邏輯 |

### Repository 隔離
- **禁止**：`SELECT * FROM ... WHERE id = ?`
- **必須**：`SELECT * FROM ... WHERE id = ? AND center_id = ?`

---

## 4. ContextHelper 模式

控制器使用 `ContextHelper` 統一取值：

```go
func (c *Controller) Handle(ctx *gin.Context) {
    helper := NewContextHelper(ctx)

    id := helper.MustParamUint("id")
    var req Request
    if !helper.MustBindJSON(&req) { return }

    result, errInfo, err := c.service.DoSomething(ctx, id, &req)
    if err != nil {
        helper.ErrorWithInfo(errInfo)
        return
    }
    helper.Success(result)
}
```

### 常用方法
- **取值**：`MustUserID()`, `MustCenterID()`, `MustParamUint(name)`, `MustQueryDate(key)`
- **Binding**：`MustBindJSON(&req)`
- **Query**：`QueryStringOrDefault(key, default)`, `QueryIntOrDefault(key, default)`
- **響應**：`Success(data)`, `Created(data)`, `BadRequest(msg)`, `NotFound(msg)`, `Forbidden(msg)`, `Conflict(msg)`, `Unauthorized(msg)`, `InternalError(msg)`
- **錯誤**：`ErrorWithInfo(errInfo)`, `ErrorWithCode(status, code, msg)`

---

## 5. Service 層標準

```go
type UserService struct {
    BaseService
    repo *repositories.UserRepository
}

func NewUserService(app *app.App) *UserService {
    return &UserService{
        BaseService: *NewBaseService(app, "UserService"),
        repo:        repositories.NewUserRepository(app),
    }
}

// Triple Return Pattern
func (s *UserService) Get(ctx context.Context, id uint) (*models.User, *errInfos.Res, error) {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, s.App.Err.New(errInfos.NOT_FOUND), err
    }
    return user, nil, nil
}
```

### BaseService 提供的功能

**分頁**
- `PaginationParams` / `PaginationResult` - 分頁參數與結果
- `Validate()` - 驗證並修正參數
- `GetOffset()` / `BuildOrderClause()` - 取得偏移量與排序

**FilterBuilder**
- `AddEq()`, `AddNe()`, `AddGt()`, `AddGte()`, `AddLt()`, `AddLte()` - 條件過濾
- `AddLike()`, `AddIn()`, `AddNotIn()` - 模糊/範圍過濾
- `AddBetween()`, `AddDateBetween()` - 區間過濾
- `AddCenterScope(centerID)` - 中心範圍隔離

**ServiceLogger**
- `Debug()`, `Info()`, `Warn()`, `Error()` - 結構化日誌

---

## 6. 命名慣例

| 類型 | 慣例 | 範例 |
|:---|:---|:---|
| 檔案名稱 | snake_case | `admin_user.go`, `center_term.go` |
| Table 名稱 | snake_case + plural | `admin_users`, `center_terms` |
| Go Struct | PascalCase | `AdminUser`, `CenterTerm` |
| JSON Field | snake_case | `line_user_id` |
| Method | PascalCase | `CreateUser()` |
| Variable | camelCase | `userRepository` |
| Constant | UPPER_SNAKE_CASE | `SQL_ERROR` |

### 檔案 vs Struct 對照
| 檔案 | Struct |
|:---|:---|
| `admin_user.go` | `AdminUser` |
| `schedule_rule.go` | `ScheduleRule` |
| `center_term.go` | `CenterTerm` |

---

## 7. 錯誤處理

### Triple Return Pattern
```go
func GetUser(ctx context.Context, id uint) (*models.User, *errInfos.Res, error)
```

### 錯誤碼規範
- 格式：`appID(1) + module(2) + serial(4)` = 7 位數字
- 例如：`1-1-0001` → `110001` → `SQL_ERROR`
- 模組：1=System, 2=DB/Cache, 3=Auth, 4=Resource, 5=Schedule, 6=Exception, 7=Teacher

### 常見錯誤碼範例
```go
const (
    SQL_ERROR      ErrCode = 20001  // 2=DB, 0001
    NOT_FOUND     ErrCode = 40001  // 4=Resource, 0001
    FORBIDDEN     ErrCode = 30002  // 3=Auth, 0002
    SYSTEM_ERROR  ErrCode = 10001  // 1=System, 0001
)
```

---

## 8. 排課驗證引擎

### 三層驗證

1. **Scope Check**：確保操作在指定 `center_id` 下
2. **Hard Overlap Check**：
   - `Existing.Start < New.End AND Existing.End > New.Start`
   - 若衝突，直接報錯 `E_OVERLAP`
3. **Buffer Check**：
   - Room Buffer / Teacher Buffer
   - 若衝突且 `offering.allow_buffer_override = true`，可 override

### Update Mode
- `SINGLE`：僅修改單一場次
- `FUTURE`：修改此場次及之後所有場次
- `ALL`：修改整串循環規則

---

## 9. Exception 狀態機

```
[PENDING] ──教師撤回──→ [REVOKED]
    ├── 管理員同意 ──→ [APPROVED] ──→ 發送 LINE 通知
    └── 管理員拒絕 ──→ [REJECTED] ──→ 發送 LINE 通知
```

### Re-validation 規則
- Approve 瞬間執行 `validate(new_time)`
- Soft Conflict → 彈出警告，允許 Admin Override
- Hard Conflict → **直接報錯**，禁止核准

---

## 10. 權限矩陣

| 資源 | Teacher | Center Admin |
|:---|:---:|:---:|
| Center Schedule | Read | CRUD |
| Personal Event | Own CRUD | ❌ |
| Exception | Submit | Approve/Reject |
| Teacher Profile | Own Edit | ❌ |
| Teacher Certs | Own Upload | Read (已加入中心) |
| Talent Pool | ❌ | Search |
| Settings | ❌ | Update Policy |
| Room | ❌ | CRUD |
| Admin Users | ❌ | OWNER only |

---

## 11. 資料隔離

### 核心原則
- **後端隔離，前端透明**
- JWT Token 包含 `center_id`
- 所有查詢根據 JWT Token 自動過濾

### 禁止做法
- 依賴前端傳遞的 `center_id` 參數
- URL 中顯示 `center_id`

---

## 12. 併發控制

使用 **DB Row Lock**：
```go
tx := db.Begin()
tx.Select("id").From("teachers").Where("id = ?", id).ForUpdate()
 // 執行 overlap 檢查
 // Insert/Update
tx.Commit()
```

---

## 13. 通知佇列系統

### Redis Keys
| Key | 用途 |
|:---|:---|
| `notification:pending` | 待處理的通知 |
| `notification:retry` | 需要重試的通知 |
| `notification:stats` | 統計計數器 |

### 運作機制
- **最大重試次數**：3 次
- **重試延遲**：5 秒（指數遞增）
- **通知類型**：EXCEPTION_SUBMIT, EXCEPTION_RESULT, WELCOME
- **RecipientType**：ADMIN, TEACHER

### 相關檔案
- `app/services/redis_queue_service.go` - Redis 佇列服務實現

---

## 14. 開發鐵律

### 最小改動原則
- 只修改必要的地方
- 修改前搜尋相關使用處，評估影響範圍

### TDD 強制
- 先寫測試，再開發
- 開發階段使用現有資料庫測試

### 提交規範
- **英文訊息**
- 每次改動立即 commit
- 格式：`feat(scope): description (Ref: PDR章節)`

---

## 15. API 設計

### Response 格式
```json
{ "code": 0, "message": "Operation successful", "datas": { ... } }
```

### 分頁 Response
```json
{
  "code": 0,
  "message": "success",
  "datas": {
    "data": [...],
    "total": 150,
    "page": 1,
    "limit": 20,
    "total_pages": 8,
    "has_next": true,
    "has_prev": false
  }
}
```

### 通用參數
- `page`：頁碼（預設 1）
- `limit`：每頁筆數（預設 20，最大 100）
- `sort_by` / `sort_order`：排序

---

## 16. 專案結構

```
/
├── app/
│   ├── controllers/   # API 入口
│   ├── requests/     # 參數驗證
│   ├── services/     # 業務邏輯
│   ├── repositories/ # DB 存取
│   ├── resources/    # Response 轉換
│   ├── models/       # 數據模型
│   └── scheduling/   # 排課引擎 Domain
├── frontend/         # Nuxt 3
├── pdr/              # 規劃文件
└── database/         # Migrations
```

### 新增端點流程
1. Model → `app/models/`
2. Request → `app/requests/`
3. Repository → `app/repositories/`
4. Resource → `app/resources/`
5. Service → `app/services/`
6. Controller → `app/controllers/`
7. Register → `app/servers/route.go`

---

## 17. 時區中央化

- 統一使用 **Asia/Taipei**
- 前端使用 `composables/useTaiwanTime.ts`
- **禁止**使用 `toISOString()` 處理日期
- ⚠️ **現有違規**：部分程式碼仍使用 `toISOString()`（需逐步修正）

### 正確做法
```typescript
import { formatDateToString, getTodayString } from '~/composables/useTaiwanTime'

// 錯誤
const dateStr = new Date().toISOString().split('T')[0]

// 正確
const dateStr = formatDateToString(new Date())
const todayStr = getTodayString()
```

### 允許使用情境
- iCal 標準格式匯出（需要 UTC）
- 測試檔案
- 僅用於檔案名稱產生（無業務邏輯）

---

## 18. 前端規範

### Composable 使用
```typescript
// 正確
const { getCenterId } = useCenterId()

// 錯誤
import { getCenterId } from '~/composables/useCenterId'
```

### 禁止原生 Alert
```typescript
// 錯誤
alert('錯誤')

// 正確
import { alertError, alertConfirm } from '~/composables/useAlert'
await alertError('錯誤')
```

---

## 19. 常用命令

```bash
# Build
go build -mod=vendor -o main .

# Run (需要 MySQL + Redis)
go run main.go

# Test
go test ./testing/test/...

# Swagger
swag init
```

---

## 20. 環境設定

| 服務 | 預設 Port |
|:---|:---|
| HTTP API | localhost:8888 |
| gRPC | localhost:50051 |
| WebSocket | localhost:8889 |
| Swagger | /swagger/index.html |

---

## 21. Agent 技能

| 技能 | 用途 |
|:---|:---|
| `scheduling-validator` | 排課引擎 TDD |
| `contract-sync` | API 與 Model 同步 |
| `auth-adapter-guard` | LINE Login abstraction |

---

## 📚 詳細參考文件

| 文件 | 內容 |
|:---|:---|
| `pdr/API.md` | API 詳細規格 |
| `pdr/Mysql.md` | 資料庫 Schema |
| `pdr/Stages.md` | 開發階段詳細記錄 |
| `pdr/功能業務邏輯.md` | LINE 通知系統詳解 |
| `app/services/scheduling_*.go` | 排課邏輯源碼 |
| `app/models/` | 數據模型定義 |

---

## ⚠️ 重要提醒

1. **最小改動** - 避免改 A 壞 B
2. **驗證優先** - 完成後執行相關測試
3. **資料隔離** - 後端負責，嚴禁依賴前端傳遞 center_id
4. **時區一致** - 統一使用 Asia/Taipei
