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
- **LINE 綁定**：首次登入自動綁定，**不可解除綁定**（LINE 即為帳號本身）

### 管理員端
- **Email/Password + JWT**（24 小時效期）
- 支援角色分級：OWNER、ADMIN、STAFF
- **LINE 綁定**：用於接收 Exception 即時通知（可綁定/解除綁定）

---

## 3.1 LINE 通知系統 (LINE Notification System)

### 3.1.1 通知策略

**多管理員通知：全員通知 + 已讀機制**
- 老師提交 Exception 時，發送 LINE 通知給中心所有管理員（OWNER + ADMIN）
- 每位管理員都會收到通知，可依情況處理
- 第一位處理者完成後，狀態更新，其他管理員可看到「已由他人處理」

**通知觸發時機**
| 事件 | 通知對象 | 訊息類型 |
|:---|:---|:---|
| 老師提交 Exception | 所有已綁定管理員 | 新例外申請 |
| 管理員核准 Exception | 申請老師 | 核准結果 |
| 管理員拒絕 Exception | 申請老師 | 拒絕結果 |

### 3.1.2 LINE Bot 訊息格式

**Exception 申請通知（Flex Message）**
```json
{
  "type": "flex",
  "altText": "新的例外申請通知",
  "contents": {
    "type": "bubble",
    "body": {
      "type": "box",
      "layout": "vertical",
      "contents": [
        { "type": "text", "text": "🔔 新的例外申請", "weight": "bold" },
        { "type": "text", "text": "━━━━━━━━━━━━━━" },
        { "type": "text", "text": "👤 申請人：陳小美 老師" },
        { "type": "text", "text": "📋 類型：調課申請" },
        { "type": "text", "text": "📅 日期：2026/01/28 (三)" },
        { "type": "text", "text": "🕐 時間：14:00 → 16:00" },
        { "type": "text", "text": "📝 原因：與客戶會議衝突" },
        { "type": "text", "text": "━━━━━━━━━━━━━━" },
        { "type": "text", "text": "⚠️ 此時段已有其他課程，請確認是否有衝突" }
      ]
    },
    "footer": {
      "type": "box",
      "layout": "horizontal",
      "contents": [
        {
          "type": "button",
          "style": "primary",
          "action": { "type": "uri", "label": "前往處理", "uri": "https://timeledger.app/admin/exceptions/456" }
        }
      ]
    }
  }
}
```

### 3.1.3 管理員 LINE 綁定功能

**綁定流程**
1. 管理員登入後台 → 設定 → 通知設定
2. 點擊「開始綁定」
3. 後端產生 6 位數驗證碼 + 顯示 QR Code
4. 管理員開啟 LINE，搜尋官方帳號並傳送驗證碼
5. LINE Bot 驗證成功，回覆「綁定成功」

**解除綁定流程**
1. 管理員點擊「解除綁定」
2. 彈出確認對話框
3. 點擊「確定解除」
4. 後端清除 `line_user_id`
5. 發送 LINE 通知告知已解除綁定

**通知開關**
- 可選擇性關閉特定類型的通知（不解除綁定）
- 選項：接收新例外通知、接收審核結果通知

### 3.1.4 官方帳號歡迎訊息

**老師歡迎訊息（首次登入/受邀請）**
- 標題：👋 歡迎加入 TimeLedger！
- 內容：中心名稱、功能說明
- 按鈕：立即綁定（開啟 LIFF 頁面）

**管理員歡迎訊息（首次登入且未綁定）**
- 標題：🎉 歡迎使用 TimeLedger！
- 內容：中心名稱、角色、即時通知功能說明
- 按鈕：立即綁定、稍後再說

### 3.1.5 資料庫擴展

```go
// AdminUser - 管理員
type AdminUser struct {
    // ... 現有欄位
    LineUserID         string     `gorm:"type:varchar(64);index" json:"-"`                    // LINE 用戶 ID
    LineBindingCode    string     `gorm:"type:varchar(8)" json:"-"`                          // 綁定驗證碼
    LineBindingExpires *time.Time `json:"-"`                                                 // 驗證碼過期時間
    LineNotifyEnabled  bool       `gorm:"default:true" json:"line_notify_enabled"`           // 是否接收通知
    LineBoundAt        *time.Time `json:"line_bound_at"`                                     // 綁定時間
}

// Teacher - 老師
type Teacher struct {
    // ... 現有欄位
    LineUserID   string     `gorm:"type:varchar(64);index" json:"line_user_id"` // 帳號 ID，不可解除綁定
    IsActive     bool       `gorm:"default:false" json:"is_active"`             // 是否已激活
    InvitedAt    *time.Time `json:"invited_at"`                                 // 邀請時間
    ActivatedAt  *time.Time `json:"activated_at"`                               // 激活時間
}
```

### 3.1.6 API 設計

| Method | Endpoint | 說明 |
|:---:|:---|:---|
| **管理員 LINE 綁定** |
| GET | `/admin/me/line-binding` | 取得綁定狀態 |
| POST | `/admin/me/line/bind` | 產生綁定驗證碼 |
| POST | `/admin/me/line/verify` | 驗證綁定（輸入驗證碼） |
| DELETE | `/admin/me/line/unbind` | 解除綁定 |
| PATCH | `/admin/me/line/notify-settings` | 更新通知開關 |
| **老師邀請** |
| POST | `/admin/teachers/:id/invite` | 發送邀請 Email + LINE 歡迎訊息 |

### 3.1.7 LINE Bot 回覆關鍵字

| 關鍵字 | 回覆 |
|:---|:---|
| `綁定` | 顯示綁定連結 |
| `幫助` | 使用說明 |
| `狀態` | 查詢綁定狀態 |
| `解除綁定` | 顯示解除綁定連結 |

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

func (s *UserService) Get(ctx context.Context, req *requests.UserGetRequest) (datas any, *errInfos.Res, error) {
    user, err := s.userRepository.Get(ctx, models.User{ID: uint(req.ID)})
    if err != nil {
        return nil, s.app.Err.New(errInfos.SQL_ERROR), err
    }
    // ... 業務邏輯 ...
    response, _ := s.userResource.Get(ctx, user)
    return response, nil, nil
}
```

### 4.5 Service 層職責定義

Service 層是業務邏輯的核心樞紐，負責封裝所有業務規則與流程。

**職責範圍：**

| 職責類型 | 說明 | 範例 |
|:---|:---|:---|
| **業務邏輯** | 核心業務規則的實作 | `TeacherProfileService.UpdateProfile()` 更新個人資料並同步標籤 |
| **狀態機流轉** | 管理物件狀態的轉換與驗證 | Exception 從 PENDING → APPROVED |
| **交易控制** | 跨 Repository 的資料一致性 | 更新 Profile 時同時記錄 Audit Log |
| **依賴注入** | 組合多個 Repository 與 Resource | `TeacherProfileService` 注入 teacherRepo、skillRepo、certificateRepo |
| **邊界驗證** | 業務層面的資料驗證 | 技能歸屬權限檢查（TeacherID 匹配） |

**禁止事項：**
- 直接操作資料庫（應透過 Repository）
- 處理 HTTP 請求/響應（應由 Controller 處理）
- 包含與業務無關的通用工具函數

### 4.6 ContextHelper 應用

`ContextHelper` 是控制器重構的核心工具，統一封裝 HTTP 上下文取值與響應格式。

**建立方式：**
```go
func (c *TeacherProfileController) GetProfile(ctx *gin.Context) {
    helper := NewContextHelper(ctx)
    // ...
}
```

**取值方法對照：**

| 操作類型 | 重構前 | 重構後（ContextHelper） |
|:---|:---|:---|
| 使用者 ID | `ctx.GetUint(global.UserIDKey)` | `helper.MustUserID()` |
| URL 參數 | `fmt.Sscanf(ctx.Param("id"), "%d", &id)` | `helper.MustParamUint("id")` |
| Query 參數 | `ctx.Query("status")` | `helper.QueryStringOrDefault("status", "")` |
| Query 日期範圍 | 自行解析 `from`, `to` 參數 | `helper.QueryDateRange("from", "to")` |
| JSON 綁定 | `ctx.ShouldBindJSON(&req)` | `helper.MustBindJSON(&req)` |

**響應方法對照：**

| 響應類型 | 重構前 | 重構後（ContextHelper） |
|:---|:---|:---|
| 成功響應 | `ctx.JSON(http.StatusOK, global.ApiResponse{...})` | `helper.Success(data)` |
| 錯誤響應 | `ctx.JSON(http.StatusBadRequest, global.ApiResponse{...})` | `helper.BadRequest(message)` |
| 未找到 | `ctx.JSON(http.StatusNotFound, ...)` | `helper.NotFound("message")` |
| 禁止存取 | `ctx.JSON(http.StatusForbidden, ...)` | `helper.Forbidden("message")` |
| 衝突錯誤 | `ctx.JSON(http.StatusConflict, ...)` | `helper.Conflict("message")` |
| 內部錯誤 | `ctx.JSON(http.StatusInternalServerError, ...)` | `helper.InternalError("message")` |
| 帶錯誤資訊 | `ctx.JSON(http.StatusBadRequest, global.ApiResponse{Code: ..., Message: ...})` | `helper.ErrorWithInfo(errInfo)` |

**完整範例（重構後）：**
```go
func (c *TeacherProfileController) UpdateProfile(ctx *gin.Context) {
    helper := NewContextHelper(ctx)

    teacherID := helper.MustUserID()
    if teacherID == 0 {
        return
    }

    id := helper.MustParamUint("id")
    if id == 0 {
        return
    }

    var req requests.UpdateProfileRequest
    if !helper.MustBindJSON(&req) {
        return
    }

    // 呼叫 Service 層
    profile, errInfo, err := c.profileService.UpdateProfile(ctx.Request.Context(), id, &req)
    if err != nil {
        helper.ErrorWithInfo(errInfo)
        return
    }

    helper.Success(profile)
}
```

### 4.7 TeacherProfileService 範例

`TeacherProfileService` 是重構階段的示範案例，集中管理教師個人檔案相關業務邏輯。

**服務結構：**
```go
// TeacherProfileService 教師個人檔案相關業務邏輯
type TeacherProfileService struct {
    BaseService
    app             *app.App
    teacherRepo     *repositories.TeacherRepository
    membershipRepo  *repositories.CenterMembershipRepository
    centerRepo      *repositories.CenterRepository
    skillRepo       *repositories.TeacherSkillRepository
    certificateRepo *repositories.TeacherCertificateRepository
    hashtagRepo     *repositories.HashtagRepository
    auditLogRepo    *repositories.AuditLogRepository
}

func NewTeacherProfileService(app *app.App) *TeacherProfileService {
    return &TeacherProfileService{
        app:             app,
        teacherRepo:     repositories.NewTeacherRepository(app),
        membershipRepo:  repositories.NewCenterMembershipRepository(app),
        centerRepo:      repositories.NewCenterRepository(app),
        skillRepo:       repositories.NewTeacherSkillRepository(app),
        certificateRepo: repositories.NewTeacherCertificateRepository(app),
        hashtagRepo:     repositories.NewHashtagRepository(app),
        auditLogRepo:    repositories.NewAuditLogRepository(app),
    }
}
```

**主要方法：**

| 方法 | 功能 | 錯誤碼 |
|:---|:---|:---|
| `GetProfile(ctx, teacherID)` | 取得老師個人資料 | SQL_ERROR |
| `UpdateProfile(ctx, teacherID, req)` | 更新老師個人資料 | SQL_ERROR |
| `GetCenters(ctx, teacherID)` | 取得老師已加入的中心列表 | SQL_ERROR |
| `GetSkills(ctx, teacherID)` | 取得老師技能列表 | SQL_ERROR |
| `CreateSkill(ctx, teacherID, req)` | 新增老師技能 | SQL_ERROR |
| `UpdateSkill(ctx, skillID, teacherID, req)` | 更新老師技能 | NOT_FOUND, FORBIDDEN |
| `DeleteSkill(ctx, skillID, teacherID)` | 刪除老師技能 | NOT_FOUND, FORBIDDEN |
| `GetCertificates(ctx, teacherID)` | 取得老師證照列表 | SQL_ERROR |
| `CreateCertificate(ctx, teacherID, req)` | 新增老師證照 | SQL_ERROR |
| `DeleteCertificate(ctx, certID, teacherID)` | 刪除老師證照 | NOT_FOUND, FORBIDDEN |

**權限檢查範例（UpdateSkill）：**
```go
func (s *TeacherProfileService) UpdateSkill(ctx context.Context, skillID, teacherID uint, req *UpdateSkillRequest) (*models.TeacherSkill, *errInfos.Res, error) {
    skill, err := s.skillRepo.GetByID(ctx, skillID)
    if err != nil {
        return nil, s.app.Err.New(errInfos.NOT_FOUND), err
    }

    // 權限檢查：確保技能歸屬於該老師
    if skill.TeacherID != teacherID {
        return nil, s.app.Err.New(errInfos.FORBIDDEN), nil
    }

    // ... 更新邏輯 ...
    return skill, nil, nil
}
```

**單元測試覆蓋：**
- `app/services/teacher_profile_test.go` - 包含 Profile CRUD、Skill CRUD、Certificate CRUD 等測試案例

### 4.8 BaseService 基礎服務結構（指令 10）

`BaseService` 是服務層的基礎設施，提供分頁、過濾、日誌等通用功能。

**結構定義：**
```go
type BaseService struct {
    App    *app.App
    Logger *ServiceLogger
}

func NewBaseService(app *app.App, component string) *BaseService {
    return &BaseService{
        App:    app,
        Logger: NewServiceLogger(app, component),
    }
}
```

#### 4.8.1 PaginationParams 分頁參數

```go
type PaginationParams struct {
    Page      int    `json:"page"`
    Limit     int    `json:"limit"`
    SortBy    string `json:"sort_by"`
    SortOrder string `json:"sort_order"`
}

// 主要方法
func (p *PaginationParams) Validate()                    // 驗證並修正參數
func (p *PaginationParams) GetOffset() int              // 取得偏移量
func (p *PaginationParams) BuildOrderClause() string    // 建立排序子句
func DefaultPagination() *PaginationParams              // 取得預設分頁參數
```

**驗證邏輯：**
- `page` 為負數時預設為 1
- `limit` 為 0 時預設為 20
- `limit` 超過 100 時上限為 100
- `sort_order` 無效時預設為 DESC

#### 4.8.2 PaginationResult 分頁結果

```go
type PaginationResult struct {
    Data       interface{} `json:"data"`
    Total      int64       `json:"total"`
    Page       int         `json:"page"`
    TotalPages int         `json:"total_pages"`
    HasNext    bool        `json:"has_next"`
    HasPrev    bool        `json:"has_prev"`
}

func NewPaginationResult(data interface{}, total int64, params *PaginationParams) *PaginationResult
```

**計算邏輯：**
- `TotalPages` = `(total + limit - 1) / limit`
- `HasNext` = `page < TotalPages`
- `HasPrev` = `page > 1`

#### 4.8.3 FilterBuilder 過濾建構器

```go
type FilterBuilder struct {
    conditions []string
    args       []interface{}
}

// 支援的方法
func (fb *FilterBuilder) AddEq(column string, value interface{}) *FilterBuilder
func (fb *FilterBuilder) AddNe(column string, value interface{}) *FilterBuilder
func (fb *FilterBuilder) AddIn(column string, values []interface{}) *FilterBuilder
func (fb *FilterBuilder) AddBetween(column string, min, max interface{}) *FilterBuilder
func (fb *FilterBuilder) AddCenterScope(centerID uint) *FilterBuilder
func (fb *FilterBuilder) IsEmpty() bool
func (fb *FilterBuilder) Build() (string, []interface{})
```

**使用範例：**
```go
fb := NewFilterBuilder()
conditions, args := fb.
    AddEq("status", "active").
    AddIn("category", []interface{}{"A", "B"}).
    AddBetween("created_at", "2026-01-01", "2026-12-31").
    AddCenterScope(centerID).
    Build()
// 輸出："status = ? AND category IN (?, ?) AND created_at BETWEEN ? AND ? AND center_id = ?"
```

#### 4.8.4 ServiceLogger 結構化日誌

```go
type ServiceLogger struct {
    logger    *logger.Logger
    component string
    enabled   bool  // 測試環境自動禁用
}

// 支援的方法
func (sl *ServiceLogger) Debug(message string, keysAndValues ...interface{})
func (sl *ServiceLogger) Info(message string, keysAndValues ...interface{})
func (sl *ServiceLogger) Warn(message string, keysAndValues ...interface{})
func (sl *ServiceLogger) Error(message string, keysAndValues ...interface{})
func (sl *ServiceLogger) ErrorWithErr(message string, err error, keysAndValues ...interface{})
```

**日誌格式範例：**
```
[2026/01/30 23:37:47] [Debug] [ScheduleValidation] message=checking overlap center_id=1
[2026/01/30 23:37:47] [Warn] [ScheduleValidation] slow_query_duration=413ms
```

#### 4.8.5 服務層標準範本

```go
type MyService struct {
    BaseService
    repo *MyRepository
}

func NewMyService(app *app.App) *MyService {
    baseSvc := NewBaseService(app, "MyService")
    return &MyService{
        BaseService: *baseSvc,
        repo:        NewMyRepository(app),
    }
}

func (s *MyService) DoSomething(ctx context.Context, id uint) error {
    s.Logger.Info("starting operation", "id", id)
    // 業務邏輯
    s.Logger.Debug("operation completed", "id", id)
    return nil
}
```

#### 4.8.6 已整合 BaseService 的服務

| 服務名稱 | 檔案 | 狀態 |
|:---|:---|:---:|
| ScheduleService | app/services/scheduling.go | ✅ 已整合 |
| ScheduleValidationServiceImpl | app/services/scheduling_validation.go | ✅ 已整合 |
| ScheduleExpansionServiceImpl | app/services/scheduling_expansion.go | ✅ 已整合 |
| ScheduleExceptionServiceImpl | app/services/scheduling_expansion.go | ✅ 已整合 |

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

## 12. 智慧媒合與人才庫 (Smart Matching & Talent Pool)

### 12.1 API 端點總覽

#### 智慧媒合 API

| Method | Endpoint | 說明 |
|:---:|:---|:---|
| POST | /admin/smart-matching/matches | 智慧媒合搜尋 |
| GET | /admin/smart-matching/suggestions | 搜尋建議 |
| POST | /admin/smart-matching/alternatives | 替代時段建議 |
| GET | /admin/teachers/:id/sessions | 教師課表查詢 |

#### 人才庫 API

| Method | Endpoint | 說明 |
|:---:|:---|:---|
| GET | /admin/smart-matching/talent/search | 人才庫搜尋 |
| GET | /admin/smart-matching/talent/stats | 人才庫統計 |
| POST | /admin/smart-matching/talent/invite | 邀請人才合作 |

#### 系統監控 API

| Method | Endpoint | 說明 |
|:---:|:---|:---|
| GET | /admin/notifications/queue-stats | 通知佇列統計 |

### 12.2 人才庫統計 Response 格式

```json
{
  "total_count": 156,
  "open_hiring_count": 89,
  "member_count": 45,
  "average_rating": 4.2,
  "monthly_change": 12,
  "monthly_trend": [65, 72, 78, 85, 92, 88, 95],
  "pending_invites": 23,
  "accepted_invites": 45,
  "declined_invites": 8,
  "city_distribution": [
    {"name": "台北市", "count": 52},
    {"name": "新北市", "count": 38}
  ],
  "top_skills": [
    {"name": "瑜珈", "count": 45},
    {"name": "鋼琴", "count": 38}
  ]
}
```

### 12.3 邀請人才功能

**API**: `POST /admin/smart-matching/talent/invite`

**Request Body**:
```json
{
  "teacher_ids": [1, 2, 3],
  "message": "誠摯邀請您加入我們的人才庫..."
}
```

**Response**:
```json
{
  "success_count": 2,
  "failed_count": 1,
  "failed_ids": [2],
  "invitations": [
    {"teacher_id": 1, "token": "INV-1-abc123", "status": "PENDING"},
    {"teacher_id": 3, "token": "INV-1-def456", "status": "PENDING"}
  ],
  "message": "1 位老師已有待處理邀請，無法重複邀請"
}
```

**邀請邏輯規則**：
- 同一個老師對同一個中心只能有一筆待處理邀請
- 如果已有待處理邀請，再次邀請會被拒絕並回傳 failed_ids
- 邀請有效期為 7 天
- 發送 LINE 通知（非同步處理）

### 12.4 評分因子

| 因子 | 權重 | 評分邏輯 |
|:---|:---:|:---|
| **Availability** | 40% | 完全空閒 +40分，Buffer 衝突 +15分，Hard Overlap 0分 |
| **Internal Evaluation** | 40% | 星等評分正規化 0~30分，內部備註關鍵字額外 +10分 |
| **Skill & Region Match** | 20% | 技能命中 +10分，標籤命中 +8分，地區命中 +10分 |

---

## 13. 通知佇列系統 (Notification Queue System)

### 13.1 架構

```
前端監控頁面 (/admin/queue-monitor)
         ↓
通知佇列統計 API (/admin/notifications/queue-stats)
         ↓
Redis Queue (notification:pending, notification:retry)
         ↓
Background Worker (非同步處理)
```

### 13.2 Redis Queue 結構

| Queue Key | 說明 |
|:---|:---|
| `notification:pending` | 待處理的通知 |
| `notification:retry` | 需要重試的通知 |
| `notification:completed` | 已完成的通知 |
| `notification:failed` | 失敗的通知（超過最大重試次數） |

### 13.3 佇列統計 API

**Response 格式**：
```json
{
  "pending_count": 15,
  "retry_count": 3,
  "completed_count": 1250,
  "failed_count": 12,
  "failure_rate": 0.95,
  "redis_connected": true,
  "worker_running": true
}
```

### 13.4 Notification Worker 配置

**環境變數**：
```bash
# Notification Worker（預設關閉）
NOTIFICATION_WORKER_ENABLED=true
```

**啟動方式**：
```bash
# 僅啟動 Worker
NOTIFICATION_WORKER_ENABLED=true go run main.go

# 同時啟動 API Server 和 Worker
go run main.go
```

### 13.5 監控頁面

**路徑**：管理員選單 → 系統監控 `/admin/queue-monitor`

**功能特色**：
- 通知佇列統計卡片（待處理/重試/已完成/失敗）
- 失敗率警示（超過 10% 顯示警告）
- Redis 連線狀態
- 人才庫邀請統計
- 自動重新整理（每 30 秒）

---

## 14. Hashtag 標籤管理

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

## 15. 循環行程與例外處理 (Recurrence & Exceptions)

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

## 16. 開發鐵律 (Development Rules)

### 16.1 遵循計劃
- 嚴格按照 `pdr/Stages.md` 的檢查清單執行
- **禁止跳階段** 或 超前部署

### 16.2 TDD 強制執行
- 每個 Service 或 Logic 模組 **必須先寫測試**
- **開發階段**：使用現有開發資料庫（MySQL port 3306）進行測試，建立測試資料後驗證功能
- **測試資料**：建立測試資料 → 執行測試 → 驗證結果 → 清理測試資料（或標記便於識別）
- 後端功能未通過測試視為 **未完成**

### 16.3 原子化開發（Vertical Slices）
- 一次僅開發一個獨立子功能
- **嚴禁** 同時改動多個不相關模組
- 開發順序：`Migration → Unit Test → Backend Service → API → Frontend UI → Integration Test`

### 16.4 提交規範 (Commit Standards)
- 後端完成且測試通過 → **Commit**
- 前端完成 → **再次 Commit**
- **每次修改（包含小修正）都必須立即 commit**，避免累積大量未提交的變更
- 每次 Commit 前必須更新 `pdr/progress_tracker.md`
- **Commit Message 必須使用英文**，避免中文編碼問題
- Commit Message 格式：`feat(scope): description (Ref: PDR章節)`

**正確的 Commit Message 範例：**
```
feat(auth): add quick login buttons for admin and teacher pages
fix(frontend): remove undefined mock function calls
docs: update progress tracker with test coverage results
```

**錯誤的 Commit Message 範例（禁止使用）：**
```
新增快速登入功能 <-- 使用中文
修正登入問題 <-- 使用中文
```

### 16.5 文件回饋循環（Gap Handling）
發現 API、欄位或邏輯缺失於 PDR 文件時：
1. **暫停**開發
2. **更新**相關 PDR 文件
3. **通知**用戶確認後再繼續

---

## 17. API 設計規範 (API Standards)

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

## 18. 當前開發階段 (Current Stage)

**Stage 1：基建與設計系統（Core & Design Tokens）** ✅ 已完成
- [x] 1.1 Workspace Init：Docker Compose（MySQL 8、Redis）、Monorepo 初始化
- [x] 1.2 Migrations (Base)：建立 `centers`、`users`、`geo_cities`、`geo_districts`
- [x] 1.3 UI Design System：
  - [x] Tailwind Config（Midnight Indigo 漸層）、Google Fonts 引入
  - [x] 基礎組件：`BaseGlassCard`、`BaseButton`、`BaseInput`
  - [x] 基礎佈局：Admin Sidebar 與 Mobile Bottom Nav

**Stage 2：領域層與分層架構（Domain & Layered Architecture）** ✅ 已完成
- [x] Repository 層實作（GenericRepository + 自定義 Repository）
- [x] Service 層實作（業務邏輯封裝）
- [x] Controller 層實作（API 入口）
- [x] Request/Resource 層實作（驗證與響應格式）
- [x] 分層職責分離（禁止跨層直接操作）

**Stage 3：單元測試與錯誤處理（Unit Tests & Error Handling）** ✅ 已完成
- [x] TeacherProfileService 單元測試（19/19 案例通過）
- [x] Repository API 差異處理（Generic vs 自定義）
- [x] 錯誤碼前綴處理（appID=1 前綴）
- [x] ContextHelper 工具類（統一取值與響應格式）

**Stage 4：控制器精簡（Controller Simplification）** ✅ 已完成

**Stage 5：Scheduling 模組拆分（Scheduling Module Refactoring）** ✅ 已完成

### Stage 1：基建與設計系統 ✅
| 項目 | 完成時間 | 狀態 |
|:---|:---:|:---:|
| Docker Compose 環境 | 2025-01-30 | ✅ |
| MySQL 8.0 + Redis 配置 | 2025-01-30 | ✅ |
| Monorepo 初始化 | 2025-01-30 | ✅ |
| 基礎 Migration 資料表 | 2025-01-30 | ✅ |
| UI Design System | 2025-01-30 | ✅ |

### Stage 2：領域層與分層架構 ✅
| 層級 | 完成項目 |
|:---|:---|
| Repository 層 | GenericRepository + 自定義 Repository |
| Service 層 | TeacherProfileService、TeacherService、PersonalEventService 等 |
| Controller 層 | TeacherProfileController、TeacherEventController、TeacherInvitationController 等 |
| Request 層 | 參數驗證與 binding 標籤 |
| Resource 層 | Response 格式轉換（UserResource、InvitationResource） |
| 分層隔離 | 禁止跨層直接操作資料庫 |

### Stage 3：單元測試與錯誤處理 ✅
| 測試套件 | 案例數 | 狀態 |
|:---|:---:|:---:|
| GetProfile | 2 | ✅ |
| UpdateProfile | 3 | ✅ |
| SkillCRUD | 7 | ✅ |
| CertificateCRUD | 5 | ✅ |
| GetCenters | 2 | ✅ |
| **總計** | **19** | **✅ 100%** |

**錯誤碼前綴處理：**
```go
// 修復前
const SQL_ERROR = 20002

// 修復後（appID=1 前綴）
const SQL_ERROR = 120002
```

### Stage 4：控制器精簡 ✅

#### 程式碼行數變化
| 控制器 | 原始行數 | 精簡後 | 減少比例 |
|:---|:---:|:---:|:---:|
| TeacherEventController | 245 | 184 | 25% |
| TeacherInvitationController | 460 | 388 | 16% |

#### 提取的通用方法

**TeacherEventController：**
```go
func (ctl *TeacherEventController) requireTeacherID(helper *ContextHelper) uint
func (ctl *TeacherEventController) requireEventID(helper *ContextHelper) uint
func (ctl *TeacherEventController) requireTeacherAndEventID(helper *ContextHelper) (uint, uint)
```

**TeacherInvitationController：**
```go
func (ctl *TeacherInvitationController) requireTeacherID(helper *ContextHelper) uint
func (ctl *TeacherInvitationController) requireAdminID(helper *ContextHelper) uint
func (ctl *TeacherInvitationController) requireCenterID(helper *ContextHelper) uint
func (ctl *TeacherInvitationController) requireAdminAndCenterID(helper *ContextHelper) (uint, uint)
func (ctl *TeacherInvitationController) getCenterName(ctx context.Context, centerID uint) string
func (ctl *TeacherInvitationController) buildInvitationLinks(ctx context.Context, invitations []models.CenterInvitation, centerID uint) []services.InvitationLinkResponse
func (ctl *TeacherInvitationController) validateInvitationToken(helper *ContextHelper) (models.CenterInvitation, bool)
```

### Stage 5：Scheduling 模組拆分 ✅

Scheduling 模組是本專案中最複雜的領域之一，將原本單一的大型控制器拆分為符合分層架構的模組化結構，參考 Teacher 模組的成功經驗，實作 Thin Controller 模式。

#### 控制器精簡成果
| 指標 | 數值 |
|:---|:---:|
| 使用 ContextHelper | ✅ 全面採用 |
| 通用方法提取 | ✅ requireCenterID, requireRuleID, requireExceptionID, requireAdminID |
| 業務邏輯下放 | ✅ 所有驗證邏輯移至 Service 層 |
| 控制器行數 | ~700 行（更多端點） |

#### Service 層架構

| Service | 職責 | 檔案 |
|:---|:---|:---|
| ScheduleService | 主排課服務，組合子服務 | scheduling.go |
| ScheduleValidationService | 衝突檢查、緩衝驗證 | scheduling_validation.go |
| ScheduleExpansionService | 規則展開、階段偵測 | scheduling_expansion.go |
| ScheduleExceptionService | 例外申請、審核流程 | scheduling_expansion.go (內嵌) |
| ScheduleRecurrenceService | 循環編輯、影響預覽 | scheduling_expansion.go (內嵌) |

#### 類型定義 (app/services/scheduling_interface.go)

| 類型 | 說明 |
|:---|:---|
| ValidationResult | 驗證結果（含衝突列表） |
| ValidationConflict | 衝突詳細資訊 |
| ExpandedSchedule | 展開後的課表項目 |
| PhaseTransition | 階段轉換偵測結果 |
| RuleLockStatus | 規則鎖定狀態 |
| TodaySummary | 今日課表摘要 |
| RecurrenceEditMode | 循環編輯模式（SINGLE/FUTURE/ALL） |

#### 架構分層圖

```
┌─────────────────────────────────────────────────────────────┐
│                   HTTP Request                              │
└──────────────────────────┬──────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│              SchedulingController (Thin)                     │
│  • ContextHelper 統一取值                                    │
│  • 參數解析與響應格式化                                       │
│  • 呼叫 Service 層                                          │
└──────────────────────────┬──────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│              ScheduleServiceInterface                        │
│  • 衝突檢查 (CheckOverlap/Buffer)                            │
│  • 規則管理 (CRUD)                                          │
│  • 例外管理 (Create/Review)                                  │
│  • 展開與摘要 (Expand/Summary)                               │
└──────────────────────────┬──────────────────────────────────┘
                           ↓
┌───────────────┬───────────────┬───────────────┬──────────────┐
│       ↓       │       ↓       │       ↓       │      ↓       │
│  Validation  │  Expansion  │  Exception  │  Recurrence  │
│   Service    │   Service   │   Service   │   Service    │
│  • 重疊檢查  │  • 規則展開  │  • 申請創建  │  • 循環編輯  │
│  • 緩衝驗證  │  • 例外處理  │  • 審核流程  │  • 影響預覽  │
└───────────────┴───────────────┴───────────────┴──────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│                   Repository Layer                           │
│  • ScheduleRuleRepository                                    │
│  • ScheduleExceptionRepository                               │
└─────────────────────────────────────────────────────────────┘
```

#### 檔案變動清單

**新增檔案：**
- `app/services/scheduling_interface.go` - 介面與類型定義

**修改檔案：**
- `app/controllers/scheduling.go` - 精簡為 Thin Controller
- `app/services/scheduling.go` - 主要 Service 實現
- `app/services/scheduling_validation.go` - 驗證邏輯
- `app/services/scheduling_expansion.go` - 展開與例外服務

#### 驗證結果
- `go build ./app/services/ → Build successful!`

#### 與 Teacher 模組對比

| 指標 | Teacher 模組 | Scheduling 模組 |
|:---|:---:|:---:|
| 控制器行數 | ~200 行 | ~700 行（更多端點） |
| Service 數量 | 3 個 | 5 個（含子服務） |
| ContextHelper 使用 | ✅ | ✅ |
| Repository 注入 | ✅ | ✅ |
| Triple Return 模式 | ✅ | ✅ |

#### 總體統計
| 指標 | 數值 |
|:---|:---:|
| 完成 Stage | 5 個 |
| 新增 Service | 5 個（Validation/Expansion/Exception/Recurrence/Schedule） |
| 新增/精簡 Controller | 6 個 |
| 新增 Resource | 3 個 |
| 單元測試案例 | 19 個（100% 通過）+ Scheduling 測試案例 |
| 提取通用方法 | 15 個 |
| 程式碼減少 | 約 200 行（去重） |

---

### 下一步建議

| 優先順序 | 工作項目 | 預估效益 |
|:---|:---|:---|
| 高 | 單元測試覆蓋 | 為 ScheduleService 編寫單元測試 |
| 高 | ContextHelper 增強 | 減少更多重複程式碼 |
| 中 | 其他控制器精簡 | 持續去重 |
| 中 | 測試覆蓋率提升 | 目標 80% |
| 低 | API 文件更新 | 與程式碼同步 |

**建議優先處理：** 繼續精簡剩餘控制器（如 TeacherExceptionController、TeacherScheduleController 等），將通用模式複製到其他控制器中。

---

## 19. 專案結構 (Project Structure)

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

## 19.5 Alert/Confirm UI 規範

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

## 20. 通用命令 (Common Commands)

```bash
# Build
go build -mod=vendor -o main .

# Run locally (requires MySQL + Redis running)
go run main.go

# Run all tests (uses SQLite mock DB + MinRedis)
go test ./testing/test/... -v

# Run a single test
go test ./testing/test -run TestUser/... -vService_CreateAndGet

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

## 21. 環境設定 (Environment Setup)

Copy `.env.example` to `.env`. Key services：
- HTTP API：`localhost:8888`（Swagger at `/swagger/index.html`）
- gRPC：`localhost:50051`
- WebSocket：`localhost:8889`
- Health check：`/healthy`

MySQL master-slave replication：RDB（read/slave）、WDB（write/master）

---

## 22. 導入組織 (Import Organization)

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

## 23. 資料庫操作 (Database Operations)

- **Read：** `app.Mysql.RDB.WithContext(ctx)`（slave）
- **Write：** `app.Mysql.WDB.WithContext(ctx)`（master）
- Always pass `context.Context` as first parameter

---

## 24. 請求驗證 (Request Validation)

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

## 25. 通用模式 (General Patterns)

- Time fields：Unix timestamps (`int64`)
- JSON fields：stored as strings in DB, unmarshaled in resources
- Use `defer` for cleanup
- Recover panics in goroutines
- Use `app.Tools` (timezone, IP, JSON, trace ID)
- Use `app.Api` for external HTTP calls
- Use `app.Rpc` for RPC calls

---

## 26. 新增端點流程 (Adding New Endpoints)

1. Model → `app/models/<entity>.go`
2. Request → `app/requests/<entity>.go`
3. Repository → `app/repositories/<entity>.go`
4. Resource → `app/resources/<entity>.go`
5. Service → `app/services/<entity>.go`
6. Controller → `app/controllers/<entity>.go`
7. Register route → `app/servers/route.go`

---

## 27. gRPC 服務 (gRPC Services)

1. Define proto in `grpc/proto/` with `go_package`
2. Compile with `protoc`
3. Implement in `grpc/services/` embedding `Unimplemented<Name>ServiceServer`
4. Register in `grpc/server.go`

---

## 28. 測試規範 (Testing)

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

## 29. 程式碼格式化 (Formatting)

- Use tabs for indentation
- Struct tags with backticks and proper spacing
- No trailing whitespace
- Max line length：keep readable

---

## 30. 內部套件 (Internal Packages)

- `gitlab.en.mcbwvx.com/frame/teemo` - Tools (timezone, JSON utilities)
- `gitlab.en.mcbwvx.com/frame/zilean` - Logging
- `gitlab.en.mcbwvx.com/frame/ezreal` - HTTP client wrapper

---

## 31. 語言與溝通 (Language)

- **開發文件**：繁體中文與英文混用（代碼、API 為英文）
- **與用戶溝通**：**繁體中文**
- **代碼註解**：視情況使用繁體中文說明業務邏輯

---

## 32. 時區中央化 (Timezone Centralization)

### 32.1 架構設計

整個系統（後端 + 前端）統一使用台灣時區（Asia/Taipei）：

```
┌─────────────────────────────────────────────────────────────┐
│                    TimeLedger 系統                          │
├─────────────────────────────────────────────────────────────┤
│  後端 (Go)                                                 │
│  ├── APP_TIMEZONE=Asia/Taipei (預設)                        │
│  ├── MySQL: loc=Asia/Taipei                               │
│  └── app/timezone.go: 中央化時區管理                        │
├─────────────────────────────────────────────────────────────┤
│  通訊 (API)                                               │
│  └── 日期格式: YYYY-MM-DD (字串)                           │
├─────────────────────────────────────────────────────────────┤
│  前端 (Nuxt 3)                                            │
│  ├── useTaiwanTime.ts: 本地時區工具                        │
│  └── 瀏覽器本地顯示                                       │
└─────────────────────────────────────────────────────────────┘
```

### 32.2 後端時區管理

**環境設定 (`configs/env.go`)**
```go
// 環境變數
APP_TIMEZONE=Asia/Taipei  // 預設值
```

**中央時區工具 (`app/timezone.go`)**
```go
// 使用 sync.Once 確保執行緒安全
var loadTaiwanLocationOnce sync.Once
var taiwanLocation *time.Location

// 載入台灣時區（只執行一次）
func LoadTaiwanLocation() (*time.Location, error)

// 取得台灣時區
func GetTaiwanLocation() *time.Location

// 取得台灣現在時間
func NowInTaiwan() time.Time

// 取得台灣今日日期
func TodayInTaiwan() time.Time
```

**MySQL 連線 (`database/mysql/conn.go`)**
```go
// DSN 增加 loc 參數，確保資料庫時間與應用程式時區一致
dsn := "...&loc=Asia/Taipei"
```

### 32.3 前端時區工具

**`frontend/composables/useTaiwanTime.ts`**
```typescript
// 格式化日期為 YYYY-MM-DD 字串
export function formatDateToString(date: Date): string

// 取得今日日期字串
export function getTodayString(): Date

// 取得週開始/結束日期
export function getWeekStart(date?: Date): Date
export function getWeekEnd(date?: Date): Date
```

**重要：避免使用 toISOString()**
- `toISOString()` 會轉換為 UTC，導致凌晨日期偏移
- 使用本地時間運算避免問題

### 32.4 移除非重複時區載入

以下檔案已移除重複的時區載入邏輯，改用中央時區工具：
- `app/services/scheduling_validation.go`
- `app/services/schedule_rule_validator.go`
- `app/controllers/scheduling.go`
- `app/repositories/schedule_rule.go`

### 32.5 前端時區修正

以下前端檔案已更新使用中央時區工具（`useTaiwanTime.ts`）：

| 檔案 | 修正內容 |
|:---|:---|
| `stores/teacher.ts` | 新增 `formatDateTimeForApi()` 函數，API 資料傳送改用台灣時區 |
| `components/ExceptionModal.vue` | `today` computed 改用 `getTodayString()` |
| `components/ScheduleMatrixView.vue` | `formatDate()` 改用 `formatDateToString()` |
| `components/PersonalEventModal.vue` | 表單初始值與 `formatDateTimeForApi()` 改用台灣時區 |
| `components/ScheduleTimelineView.vue` | `date` 格式化改用 `formatDateToString()` |
| `pages/admin/matching.vue` | CSV 匯出、API 查詢參數、請求資料改用台灣時區 |

### 32.6 禁止使用 toISOString() 處理日期

**嚴格禁止**在前端程式碼中使用 `toISOString()` 處理日期相關邏輯：

```typescript
// ❌ 錯誤做法
const dateStr = new Date().toISOString().split('T')[0]

// ✅ 正確做法
import { formatDateToString, getTodayString } from '~/composables/useTaiwanTime'

const dateStr = formatDateToString(new Date())
const todayStr = getTodayString()
```

**例外**：以下情境可繼續使用 `toISOString()`：
- iCal 標準格式匯出（需要 UTC）
- 測試檔案（模擬資料）
- 僅用於檔案名稱產生（無業務邏輯）

### 32.7 效益

| 項目 | 改善內容 |
|:---|:---|
| 時區一致性 | 後端、資料庫、前端統一使用台灣時區 |
| 日期正確性 | 避免 toISOString() 導致的凌晨日期偏移問題 |
| 程式碼維護 | 中央化時區工具，減少重複程式碼 |
| 執行緒安全 | 使用 sync.Once 確保時區只載入一次 |
| 可設定性 | 可透過環境變數調整時區 |

---

## 33. Agent 技能 (Agent Skills)

- **auth-adapter-guard**：Mock Login vs LINE Login abstraction；使用 `AuthService` interface，永遠不要直接呼叫 `liff.*`
- **contract-sync**：保持 API 規格與 Go struct 和 TypeScript interface 同步；修改 `pdr/API.md` 或 `pdr/Mysql.md` 時更新 model
- **scheduling-validator**：排課引擎 TDD；為 overlap、buffer、cross-day 邏輯先寫測試

---

> **注意**：所有 PDR 文件已整併至此，開發時請直接參考本文件。
> 原始 PDR 文件位於 `pdr/` 目錄，僅供查閱參考用。
---

## 13. AI 開發工程規範 (AI Development Principles)

為確保開發品質，避免「改 A 壞 B」，所有 AI 助理在執行任務時必須嚴格遵守以下規範：

### 13.1 預防回歸錯誤 (Prevention of Regressions)
- **修改前搜尋**：在修改任何函數或欄位前，**必須**先搜尋該項目的所有使用處 (Usage Search)，評估修改對其他功能模組、測試案例或 DTO 的影響。
- **最小改動原則**：優先採用不會破壞現有合約 (Contract) 的修改方式。

### 13.2 驗證導向開發 (Verification-Driven Development)
- **TDD 優先**：
    - 修復 Bug 前，應先撰寫能重現該 Bug 的測試案例 (Red Test)。
    - 完成修復後，確保測試案例通過 (Green Test)。
- **自動化驗證**：
    - 完成開發後，**必須**自動尋找並執行相關的測試組 (如 `go test ./testing/test/...`)。
    - 若涉及 UI 變動，應使用 `browser` 工具進行關鍵流程的視覺化確認。

### 13.3 透明化開發流程 (Transparency & Reporting)
- **實作計畫 (Implementation Plan)**：在開始大規模修改前，先產出 `implementation_plan.md` 供使用者審核：
    - 包含：Bug 根因分析、變動檔案清單、潛在影響評估、驗證策略。
- **成果報告 (Walkthrough)**：完成後產出 `walkthrough.md`：
    - 包含：具體修改點、測試執行截圖/日誌、驗證成功的證據。

### 13.4 程式碼品質審查
- **自我審核**：在完成代碼後，主動檢查：
    - 是否漏掉 `center_id` 隔離檢查？
    - 是否有無效的 JSON 欄位提取 (如 `JSON_EXTRACT` 漏掉 `JSON_UNQUOTE`)？
    - 是否符合 `CLAUDE.md` 中的分層架構規範？

### 13.5 環境適應性 (Environment Adaptability)
- **指令執行前檢查**：在執行任何終端機指令 (Terminal Commands) 前，**必須**先確認當下的作業系統環境 (如：Windows vs. macOS/Linux)。
- **跨平台指令相容性**：
    - Windows 環境下優先使用 PowerShell 或相容指令。
    - 避開僅限 Unix-like 的指令 (如 `ls`, `grep`, `export`) 或確保有對應的替代方案 (如 `dir`, `findstr`, `$env:VAR`)。
    - 嚴禁盲目嘗試指令，應先從 `user_information` 或 `run_command` 的系統提示中確認 OS 版本，以減少無謂的錯誤重試與資源浪費。

---

## 34. 前端開發規範 (Nuxt 3 Frontend Standards)

為確保開發品質並預防導覽列或頁面整合問題，AI 助手必須遵循以下規範：

### 34.1 Composable (Hooks) 使用
- **解構執行**：所有 Composable（如 `useCenterId`）不可直接 import 內部函數。必須先執行 Hook 並進行解構。
- **正確範例**：`const { getCenterId } = useCenterId()`

### 34.2 側邊欄與佈局安全
- **層級管理**：桌面版側邊欄 `<aside>` 必須具備 `relative` 或 `fixed` 定位，且 `z-index` 至少為 `20`。
- **自動狀態管理**：在 `layouts/admin.vue` 中必須實作 `watch` 監聽 `route.path`，根據當前路徑自動將 `expandedMenus` 對應項設為 `true`。

### 34.3 導覽安全性
- **跳轉前檢查**：在修改側邊欄或 `NuxtLink` 連結後，必須確認目標頁面的 script 區塊無語法錯誤或致命 Import 錯誤。前端模組崩潰會導致路由系統拒絕執行跳轉。
