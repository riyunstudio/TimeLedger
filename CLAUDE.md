# CLAUDE.md

This file provides comprehensive guidance to Claude Code (claude.ai/code) when working with code in this repository. All PDR documentation has been consolidated here to avoid redundant file reads.

---

## 1. 專案核心定位 (Project Core)

**TimeLedger** - 教師中心化多據點排課平台

- **目標市場**：台灣（LINE-First 生態系）
- **策略**：「SaaS + 人才市場」（高毛利訂閱制）
- **角色**：Lead Developer，優先重視 **教師端行動體驗** 與 **中心端治理功能**

---

## 2. 技術堆疊 (Tech Stack)

| 層面 | 技術 |
|:---:|:---|
| **後端** | Go (Gin) + MySQL 8.0 + Redis（單體架構） |
| **前端** | Nuxt 3 (SSR) + Tailwind CSS + LINE LIFF |
| **部署** | Docker Compose（單一 VPS 容器化部署） |
| **通訊** | HTTP REST API (Gin)、gRPC、WebSocket |

---

## 3. 認證策略 (Authentication)

### 教師端：LINE 單一登入
- **無密碼策略**：完全廢除「設定帳密」功能
- **LIFF Silent Login**：點開即登入，利用 LIFF SDK 取得 `id_token`
- **換手機處理**：安裝 LINE 登入即可自動恢復
- **帳號遺失處理**：聯繫中心管理員，由 Admin 後台重新綁定新的 `line_user_id`

### 管理員端
- **Email/Password + JWT**（24 小時效期）
- 支援角色分級：OWNER、ADMIN、STAFF

---

## 4. 分層架構規範 (Layered Architecture)

### 4.1 架構圖

```
HTTP Request → Middleware → Controller → Request (validation) → Service → Repository → Model
                                                                              ↓
gRPC Request → Interceptors → gRPC Service ─────────────────────────────────→┘
```

### 4.2 分層職責（嚴格遵守）

| 層級 | 職責 | 禁止事項 |
|:---:|:---|:---|
| **Controller** | 僅負責 Request 解析 → 呼叫 Service → 回傳 JSON | 寫入任何業務邏輯 |
| **Request** | 參數驗證（binding 標籤、CheckParam、CheckEnum） | 業務判斷 |
| **Service** | 核心業務邏輯、狀態機流轉、交易控制、依賴注入 | 直接操作資料庫 |
| **Repository** | 純粹 DB 操作（Find/Create/Update） | 任何業務判斷 |
| **Model** | 數據結構定義（GORM Tags） | 商業邏輯 |
| **Resource** | Model 轉換為輸出格式（DTO） | 修改資料狀態 |

### 4.3 Repository 隔離防護
- 所有查詢（除跨租戶的教師私人行程外）**必須**在 WHERE 子句中包含 `center_id`
- **禁止** `SELECT * FROM ... WHERE id = ?`
- **必須** `SELECT * FROM ... WHERE id = ? AND center_id = ?`

### 4.4 Service 層標竿代碼

```go
type UserService struct {
    BaseService
    app            *app.App
    userRepository *repositories.UserRepository
    userResource   *resources.UserResource
}

func NewUserService(app *app.App) *UserService {
    return &UserService{
        app:            app,
        userRepository: repositories.NewUserRepository(app),
        userResource:   resources.NewUserResource(app),
    }
}

func (s *UserService) Get(ctx context.Context, req *requests.UserGetRequest) (datas any, errInfo *errInfos.Res, err error) {
    user, err := s.userRepository.Get(ctx, models.User{ID: uint(req.ID)})
    if err != nil {
        return nil, s.app.Err.New(errInfos.SQL_ERROR), err
    }
    // ... 業務邏輯 ...
    response, _ := s.userResource.Get(ctx, user)
    return response, nil, nil
}
```

---

## 5. 命名慣例 (Naming Conventions)

| 類型 | 慣例 | 範例 |
|:---|:---|:---|
| **Module** | snake_case | `payment_rule` |
| **Table** | snake_case + plural | `payment_rules` |
| **Go Struct** | PascalCase | `PaymentRule` |
| **Interface** | PascalCase + type suffix | `AuthService`, `Job` |
| **JSON Field** | snake_case | `payment_rule_id` |
| **Method** | PascalCase (exported), camelCase (private) | `CreateUser()`, `validate()` |
| **Variable** | camelCase | `userRepository`, `errInfo` |
| **Constant** | UPPER_SNAKE_CASE | `SQL_ERROR`, `USER_NOT_FOUND` |
| **Context** | `ctx` | - |
| **Error** | `err` | - |
| **Error Info** | `eInfo` or `errInfo` | - |
| **Package** | lowercase, single word | `controllers`, `services` |

---

## 6. 錯誤處理 (Error Handling)

### Triple Return Pattern
```go
func GetUser(ctx context.Context, id uint) (*models.User, *errInfos.Res, error) {
    user, err := repo.GetByID(ctx, id)
    if err != nil {
        return nil, s.app.Err.New(errInfos.USER_NOT_FOUND), err
    }
    return user, nil, nil
}
```

### 錯誤碼規範
- 格式：`FunctionType(1) + Serial(4)`（例如：`10001` = System Error）
  - Type 1: System (10001-10999)
  - Type 2: DB/Cache (20001-20999)
  - Type 3: Other (30001-30999)
  - Type 4: User (40001-40999)
- 定義於 `global/errInfos/code.go`，訊息定義於 `message.go`

---

## 7. 排課驗證引擎 (Validation Engine)

### 7.1 驗證層級

1. **Scope Check**：確保操作都在指定 `center_id` 下
2. **Hard Overlap Check（硬衝突）**：
   - 查詢時段 `[Start, End]` 內，該 `Teacher` 或 `Room` 是否已有其他 `Active Session`
   - 規則：`Existing.Start < New.End AND Existing.End > New.Start`
   - 若 True，直接報錯 `E_OVERLAP`（不可覆寫）
   - 若 `teacher_id` 為空（NULL），系統跳過 Teacher Overlap 檢查，僅驗證 Room
3. **Buffer Check（緩衝）**：
   - Room Buffer：`New.Start - Prev.End < current_course.room_buffer_min`
   - Teacher Buffer：`New.Start - Prev.End < current_course.teacher_buffer_min`
   - 若衝突且 `offering.allow_buffer_override = true` → 允許帶 `override=true` 送出

### 7.2 緩衝時間計算策略

採用「取最大值」而非「相加」：

```go
// 教室緩衝時間
RoomBuffer = max(
    CourseA.room_buffer_min,
    CourseB.room_buffer_min,
    Room.cleaning_time
)

// 老師緩衝時間
TeacherBuffer = max(
    CourseA.teacher_buffer_min,
    CourseB.teacher_buffer_min,
    Teacher.default_buffer_min
)
```

### 7.3 緩衝衝突回應格式

```json
{
  "valid": false,
  "conflicts": [{
    "type": "TEACHER_BUFFER",
    "message": "老師上一堂課（13:00結束）與本堂課（13:05開始）間隔不足",
    "current_gap_minutes": 5,
    "required_buffer_minutes": 15,
    "previous_session": { "id": 123, "course_name": "瑜伽基礎", "end_at": "2026-01-20T13:00:00" },
    "can_override": true
  }]
}
```

---

## 8. 異動審核狀態機 (Exception State Machine)

### 8.1 狀態流轉

```
[PENDING] ──教師撤回──→ [REVOKED]
    │
    ├── 管理員同意 ──→ [APPROVED] ──→ 發送 LINE 通知
    │
    └── 管理員拒絕 ──→ [REJECTED] ──→ 發送 LINE 通知
```

### 8.2 狀態轉換定義

| From | To | Trigger | Action |
|:---|:---|:---|:---|
| (None) | PENDING | Teacher Submit | Create Record, Validate(Soft) |
| PENDING | REVOKED | Teacher Cancel | Mark Resolved |
| PENDING | APPROVED | Admin Approve | **Re-Validate(Hard)**, Apply to Schedule |
| PENDING | REJECTED | Admin Deny | Mark Resolved, Notify Teacher |
| APPROVED | CANCELLED | Admin Undo | Revert Schedule (if date not past) |

### 8.3 Re-validation 規則
- 管理員按下 Approve 瞬間，系統必須執行 `validate(new_time)`
- **Soft Conflict (Buffer)**：彈出警告，允許 Admin Override
- **Hard Conflict (Overlap)**：**直接報錯**，禁止核准

---

## 9. 權限管控矩陣 (RBAC Matrix)

### 9.1 角色定義
- **Visitor**：未登入訪客（無權限）
- **Teacher (Self)**：登入的老師（僅能操作自己的資料）
- **Center Admin**：中心管理員（僅能操作所屬中心的資料）
- **Super Admin**：系統總管（維運與除錯用）

### 9.2 資源存取控制

| 資源 | 動作 | Teacher | Center Admin | 備註 |
|:---|:---|:---:|:---:|:---|
| **Center Schedule** | View (Read) | ✅ (僅已加入中心) | ✅ (僅所屬中心) | |
| | Create/Edit (Write) | ❌ | ✅ | 老師不可直接改課表，需走 Exception |
| **Personal Event** | View (Read) | ✅ (Own) | ⚠️ (僅 Busy/隱私模式) | 中心僅看到 "Busy" |
| | Create/Edit (Write) | ✅ (Own) | ❌ | 中心不可修改老師私人行程 |
| **Exception (請假單)** | Create (Submit) | ✅ | ✅ (代申請) | |
| | Approve/Reject | ❌ | ✅ | 僅 Admin 有準駁權 |
| **Teacher Profile** | Edit (Bio/Skills/City/District) | ✅ | ❌ | 老師擁有自己的專業履歷 |
| **Teacher Certs** | Upload/Delete | ✅ | ❌ | |
| | View (Read) | ✅ | ✅ (僅已加入中心) | 嚴格限制：未加入之中心不可見 |
| **Talent Pool** | Search/View Profile | ❌ | ✅ | 僅限開啟 `is_open_to_hiring` 的老師 |
| **Settings** | Update Center Policy | ❌ | ✅ | |
| **Room Management** | CRUD (Add/Remove Rooms) | ❌ | ✅ | |
| **Admin Users** | CRUD (Add/Remove) | ❌ | ⚠️ (僅限 OWNER) | 僅擁有者可增刪管理員 |
| **Audit Logs** | View (Read) | ❌ | ✅ | |

---

## 10. 資料隔離防護 (Data Isolation)

### 核心原則：後端隔離，前端透明

**資料隔離是後端的責任**，前端不應在 URL 中暴露 `center_id`。

#### 後端職責
1. JWT Token 包含 `center_id`（Admin 登入時由後端設定）
2. 所有資料查詢必須根據 JWT Token 中的 `center_id` 自動過濾
3. **嚴禁**依賴前端傳遞的 `center_id` 參數

#### 前端職責
1. **禁止**在 URL 中顯示 `center_id`
2. **禁止**在 API 請求中傳遞 `center_id`
3. 完全信任後端的資料隔離機制

### Hard Scope Check（強制 Scope 檢查）

**Admin Request**：
- JWT Claim 必須包含 `role: ADMIN`
- URL Path 若包含 `/centers/:center_id`，必須驗證 `JWT.center_id == Path.center_id`
- **禁止** `SELECT * FROM ... WHERE id = ?`
- **必須** `SELECT * FROM ... WHERE id = ? AND center_id = ?`

**Teacher Request**：
- JWT Claim 必須包含 `role: TEACHER`
- 若存取 `schedule_sessions`，Query 必須內建 `WHERE center_id IN (teacher.joined_centers)`
- 若存取 `personal_events`，Query 必須內建 `WHERE teacher_id = JWT.uid`

### 範例：正確與錯誤的 API 設計

| 類型 | 錯誤做法 | 正確做法 |
|:---|:---|:---|
| **前端呼叫** | `GET /admin/centers/1/teachers` | `GET /teachers` |
| **後端實作** | 從 URL 取得 center_id | 從 JWT Token 取得 center_id |
| **URL 顯示 center_id** | 是 | 否 |
| **資料隔離依賴** | 前端傳遞參數 | JWT Token 自動過濾 |

### 敏感個資遮蔽
- **Line User ID**：僅供系統綁定，不可回傳給前端
- **Certificates**：圖片 URL 需使用 Signed URL（由 S3/Storage 產生，時效性）

---

## 11. 併發控制策略 (Concurrency Control)

針對「多中心同時排同一位老師」的 Race Condition，採用 **DB Row Lock**：

1. Transaction Start
2. 執行 `SELECT id FROM teachers WHERE id = ? FOR UPDATE`
3. 執行 Overlap 檢查
4. Insert/Update
5. Commit（釋放鎖定）

---

## 12. 智慧媒合評分因子 (Smart Matching)

| 因子 | 權重 | 評分邏輯 |
|:---|:---:|:---|
| **Availability** | 40% | 完全空閒 +40分，Buffer 衝突 +15分，Hard Overlap 0分 |
| **Internal Evaluation** | 40% | 星等評分正規化 0~30分，內部備註關鍵字額外 +10分 |
| **Skill & Region Match** | 20% | 技能命中 +10分，標籤命中 +8分，地區命中 +10分 |

---

## 13. Hashtag 標籤管理

### 儲存時同步
- 教師儲存檔案時，後端同步更新 `hashtags` 字典表與 `usage_count`
- 若標籤不存在：新增至 `hashtags` 表
- 重新計算該標籤的全域使用次數

### 個人品牌標籤限制
- 長度為 3-5 個
- 後端嚴格校驗，不符合回傳 `E_VALIDATION` 錯誤
- 確保匯出圖片的版面美觀

### 自動清理
- 更新 `usage_count` 後，若偵測到某標籤 `usage_count == 0`，立即刪除
- 每日凌晨可選掃描一次殘留資料

### 輸入規範
- 若老師輸入標籤漏掉 `#`，後端自動補上
- 前端輸入時需延遲 300~500ms 才發送搜尋請求
- 若標籤已存在（大小寫不同亦然），自動忽略

---

## 14. 循環行程與例外處理 (Recurrence & Exceptions)

### 循環類型
- `NONE`（單次）、`DAILY`（日）、`WEEKLY`（週）、`MONTHLY`（月）、`CUSTOM`（自訂）

### 展開邏輯
1. 取得 `start_at` 與循環參數
2. 計算候選日期
3. 過濾 Exceptions（`CANCEL` 不顯示，`RESCHEDULE` 顯示新時段）
4. 輸出 Sessions

### 國定假日處理
- 中心定義的「假日」優先權高於所有週期的「規則」
- **無感停課**：系統不需要為每個假日自動生成 `schedule_exceptions`，而是查詢時動態過濾

### Update Mode（更新模式）
- `SINGLE`：僅修改此單一場次（原規則產生 CANCEL 例外，新規則產生 ADD 例外）
- `FUTURE`：修改此場次及之後所有場次（原規則截斷，新規則從此場次開始）
- `ALL`：修改整串循環規則（更新 recurrence 欄位）

---

## 15. 開發鐵律 (Development Rules)

### 15.1 遵循計劃
- 嚴格按照 `pdr/Stages.md` 的檢查清單執行
- **禁止跳階段** 或 超前部署

### 15.2 TDD 強制執行
- 每個 Service 或 Logic 模組 **必須先寫測試**
- **開發階段**：使用現有開發資料庫（MySQL port 3306）進行測試，建立測試資料後驗證功能
- **測試資料**：建立測試資料 → 執行測試 → 驗證結果 → 清理測試資料（或標記便於識別）
- 後端功能未通過測試視為 **未完成**

### 15.3 原子化開發（Vertical Slices）
- 一次僅開發一個獨立子功能
- **嚴禁** 同時改動多個不相關模組
- 開發順序：`Migration → Unit Test → Backend Service → API → Frontend UI → Integration Test`

### 15.4 提交規範
- 後端完成且測試通過 → **Commit**
- 前端完成 → **再次 Commit**
- **每次修改（包含小修正）都必須立即 commit**，避免累積大量未提交的變更
- 每次 Commit 前必須更新 `pdr/progress_tracker.md`
- Commit Message 格式：`feat(scope): 描述 (Ref: PDR章節)`

### 15.5 文件回饋循環（Gap Handling）
發現 API、欄位或邏輯缺失於 PDR 文件時：
1. **暫停**開發
2. **更新**相關 PDR 文件
3. **通知**用戶確認後再繼續

---

## 16. API 設計規範 (API Standards)

### Response 格式
```json
{
  "code": "SUCCESS",
  "message": "Operation successful",
  "data": { ... }
}
```

### 分頁 Response
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "total_pages": 8,
    "has_next": true,
    "has_prev": false
  }
}
```

### 通用查詢參數
| 參數 | 類型 | 必填 | 預設值 | 說明 |
|:---|:---:|:---:|:---:|:---|
| `page` | INT | 否 | 1 | 頁碼 |
| `limit` | INT | 否 | 20 | 每頁筆數（最大 100） |
| `sort_by` | STRING | 否 | 依各 API 定義 | 排序欄位 |
| `sort_order` | STRING | 否 | ASC | 排序方向（ASC/DESC） |

---

## 17. 當前開發階段 (Current Stage)

**Stage 1：基建與設計系統（Core & Design Tokens）**
- [ ] 1.1 Workspace Init：Docker Compose（MySQL 8、Redis）、Monorepo 初始化
- [ ] 1.2 Migrations (Base)：建立 `centers`、`users`、`geo_cities`、`geo_districts`
- [ ] 1.3 UI Design System：
  - [ ] Tailwind Config（Midnight Indigo 漸層）、Google Fonts 引入
  - [ ] 基礎組件：`BaseGlassCard`、`BaseButton`、`BaseInput`
  - [ ] 基礎佈局：Admin Sidebar 與 Mobile Bottom Nav

---

## 18. 專案結構 (Project Structure)

```
/
├── apis/                   # 外部 API 接口 (Interface Layer)
├── app/                    # Go 後端核心 (Monolithic)
│   ├── controllers/        # API 入口
│   ├── requests/           # 參數驗證
│   ├── services/           # 業務邏輯
│   ├── repositories/       # DB 存取
│   ├── resources/          # Response 轉換
│   ├── models/             # 數據模型
│   ├── servers/            # Server, Route & Middleware
│   ├── scheduling/         # 排課引擎專屬邏輯 (Domain)
│   └── pagination/         # 分頁 Helper
├── global/                 # 全域共用 (ErrInfos, Config)
├── libs/                   # 第三方或內部 Library 封裝 (JWT, MQ, WS)
├── database/               # SQL Migrations & Seeders
├── configs/                # 系統環境變數配置
├── grpc/                   # gRPC 定義與實作
├── rpc/                    # RPC 通訊組件
├── testing/                # 自動化測試與 Mocks
├── frontend/               # Nuxt 3 前端應用 (User + Admin)
│   ├── pages/
│   ├── components/
│   │   ├── AdminTeacherProfileModal.vue  # 管理員查看老師檔案 Modal
│   │   ├── GlobalAlert.vue               # 全局 Alert 組件
│   │   └── ...
│   └── nuxt.config.ts
├── pdr/                    # 規劃文件 (Reference Only)
├── main.go                 # Backend Entry Point
└── docker-compose.yml      # 本地開發環境
```

### AdminTeacherProfileModal 組件

管理員查看老師個人檔案的彈窗組件。

**功能特色：**
- 顯示老師頭像、姓名、狀態
- 聯繫資訊（Email、電話、縣市區域）
- 技能標籤（包含程度）
- 證照數量統計
- 玻璃擬態 UI 設計

**Props：**
| 屬性 | 類型 | 說明 |
|:---|:---|:---|
| `teacher` | `TeacherProfile \| null` | 老師資料物件 |

**使用方式：**
```vue
<AdminTeacherProfileModal
  v-if="selectedTeacher"
  :teacher="selectedTeacher"
  @close="selectedTeacher = null"
/>
```

**TeacherProfile 結構：**
```typescript
interface TeacherProfile {
  id: number
  name: string
  email: string
  phone?: string
  city?: string
  district?: string
  is_active: boolean
  skills?: TeacherSkill[]
  certificates?: any[]
}
```

---

## 18.5 Alert/Confirm UI 規範

### 禁止使用原生 alert/confirm

**嚴格禁止**在前端程式碼中使用原生的 `alert()` 或 `confirm()`。必須使用自定義的美化彈窗組件。

**錯誤做法：**
```javascript
alert('操作失敗')
confirm('確定要刪除嗎？')
```

**正確做法：**
```typescript
import { alertError, alertConfirm, alertSuccess, alertWarning } from '~/composables/useAlert'

// 錯誤提示
await alertError('操作失敗，請稍後再試')

// 確認對話框
if (await alertConfirm('確定要刪除嗎？')) {
  // 執行刪除
}

// 成功提示
await alertSuccess('儲存成功')

// 警告提示
await alertWarning('請填寫完整資訊')
```

### GlobalAlert 組件

`frontend/components/GlobalAlert.vue` 提供美化的 Alert/Confirm 彈窗功能：

| 類型 | 用途 | 預設標題 |
|:---|:---|:---|
| `info` | 一般提示 | 提示 |
| `warning` | 警告提醒 | 提醒 |
| `error` | 錯誤訊息 | 操作失敗 |
| `success` | 成功訊息 | 操作成功 |

### useAlert Composable

**Vue 元件內使用：**
```typescript
const { error: alertError, success: alertSuccess, confirm: alertConfirm } = useAlert()
```

**非 Vue 上下文使用：**
```typescript
import { alertError, alertSuccess } from '~/composables/useAlert'
await alertError('錯誤訊息')
await alertSuccess('成功訊息')
```

### useToast Composition

用於簡短的即時提示（而非阻斷式彈窗）：
```typescript
const { success, error, warning, info } = useToast()
success('操作成功')
error('操作失敗')
```

---

## 19. 通用命令 (Common Commands)

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

---

## 20. 環境設定 (Environment Setup)

Copy `.env.example` to `.env`. Key services：
- HTTP API：`localhost:8888`（Swagger at `/swagger/index.html`）
- gRPC：`localhost:50051`
- WebSocket：`localhost:8889`
- Health check：`/healthy`

MySQL master-slave replication：RDB（read/slave）、WDB（write/master）

---

## 21. 導入組織 (Import Organization)

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

---

## 22. 資料庫操作 (Database Operations)

- **Read：** `app.Mysql.RDB.WithContext(ctx)`（slave）
- **Write：** `app.Mysql.WDB.WithContext(ctx)`（master）
- Always pass `context.Context` as first parameter

---

## 23. 請求驗證 (Request Validation)

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

---

## 24. 通用模式 (General Patterns)

- Time fields：Unix timestamps (`int64`)
- JSON fields：stored as strings in DB, unmarshaled in resources
- Use `defer` for cleanup
- Recover panics in goroutines
- Use `app.Tools` (timezone, IP, JSON, trace ID)
- Use `app.Api` for external HTTP calls
- Use `app.Rpc` for RPC calls

---

## 25. 新增端點流程 (Adding New Endpoints)

1. Model → `app/models/<entity>.go`
2. Request → `app/requests/<entity>.go`
3. Repository → `app/repositories/<entity>.go`
4. Resource → `app/resources/<entity>.go`
5. Service → `app/services/<entity>.go`
6. Controller → `app/controllers/<entity>.go`
7. Register route → `app/servers/route.go`

---

## 26. gRPC 服務 (gRPC Services)

1. Define proto in `grpc/proto/` with `go_package`
2. Compile with `protoc`
3. Implement in `grpc/services/` embedding `Unimplemented<Name>ServiceServer`
4. Register in `grpc/server.go`

---

## 27. 測試規範 (Testing)

### 開發階段測試策略
開發期間使用實際開發資料庫進行測試，簡化測試環境維護：

```go
// 使用實際開發資料庫 (port 3306)
dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
mysqlDB, _ := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})

rdb, mr, _ := mockRedis.Initialize()

appInstance := &app.App{
    MySQL: &mysql.DB{WDB: mysqlDB, RDB: mysqlDB},
    Redis: &redis.Redis{DB0: rdb},
}
```

### 測試資料策略

**Mock JWT Token 驗證**
- 後端支援 `mock-` 前綴的 JWT token 進行測試
- 格式：`Authorization: Bearer mock-teacher-token` 或 `mock-admin-token`
- 使用 mock token 時會跳過 JWT 簽名驗證，自動設定：
  - `user_id`: 1
  - `user_type`: ADMIN 或 TEACHER
  - `center_id`: 1
- 適用場景：API 端點測試、功能驗證

**使用現有資料庫資料**
- 開發階段測試直接連接 MySQL port 3306（開發資料庫）
- **無需建立測試資料**：可直接查詢現有資料進行測試
- 若資料不足，使用 `t.Skip()` 跳過測試而非建立新資料
- 查詢現有資料範例：
  ```go
  var center models.Center
  if err := appInstance.MySQL.RDB.WithContext(ctx).Order("id DESC").First(&center).Error; err != nil {
      t.Skipf("No center data available, skipping test: %v", err)
      return
  }
  ```

### 測試檔案位置
- `testing/test/`

### 測試撰寫規範
- Use table-driven tests with subtests
- Test naming：`Test<Feature>_<Action>` (e.g., `TestScheduleRuleUpdateMode_Single`)
- Verify both success and error cases
- 使用現有資料驗證功能，不強求資料完整性

### CI/CD 測試資料庫
未來建立正式 CI/CD 時，可再配置獨立的測試資料庫（port 3307）。

---

## 28. 程式碼格式化 (Formatting)

- Use tabs for indentation
- Struct tags with backticks and proper spacing
- No trailing whitespace
- Max line length：keep readable

---

## 29. 內部套件 (Internal Packages)

- `gitlab.en.mcbwvx.com/frame/teemo` - Tools (timezone, JSON utilities)
- `gitlab.en.mcbwvx.com/frame/zilean` - Logging
- `gitlab.en.mcbwvx.com/frame/ezreal` - HTTP client wrapper

---

## 30. 語言與溝通 (Language)

- **開發文件**：繁體中文與英文混用（代碼、API 為英文）
- **與用戶溝通**：**繁體中文**
- **代碼註解**：視情況使用繁體中文說明業務邏輯

---

## 31. Agent 技能 (Agent Skills)

- **auth-adapter-guard**：Mock Login vs LINE Login abstraction；使用 `AuthService` interface，永遠不要直接呼叫 `liff.*`
- **contract-sync**：保持 API 規格與 Go struct 和 TypeScript interface 同步；修改 `pdr/API.md` 或 `pdr/Mysql.md` 時更新 model
- **scheduling-validator**：排課引擎 TDD；為 overlap、buffer、cross-day 邏輯先寫測試

---

> **注意**：所有 PDR 文件已整併至此，開發時請直接參考本文件。
> 原始 PDR 文件位於 `pdr/` 目錄，僅供查閱參考用。
