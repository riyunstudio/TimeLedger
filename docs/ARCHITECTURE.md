# TimeLedger 架構概述

## 1. 系統總覽

**TimeLedger** 是一個教師中心化多據點排課平台，採用單體架構（Monolithic Architecture），支援教師端（LINE LIFF）和管理員端分離的應用。

### 1.1 技術堆疊

| 層面 | 技術 |
|:---|:---|
| **後端** | Go (Gin Framework) + MySQL 8.0 + Redis |
| **前端** | Nuxt 3 (SSR) + Tailwind CSS + LINE LIFF |
| **部署** | Docker Compose（單一 VPS 容器化部署） |
| **通訊** | HTTP REST API (Gin)、gRPC、WebSocket |

### 1.2 架構圖

```
┌─────────────────────────────────────────────────────────────┐
│                        Frontend (Nuxt 3)                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐  │
│  │ Teacher App │  │  Admin App  │  │   Public Pages  │  │
│  │   (LIFF)    │  │             │  │                 │  │
│  └─────────────┘  └─────────────┘  └─────────────────┘  │
└────────────────────────────┬────────────────────────────────┘
                            │ HTTP/REST
┌────────────────────────────┴────────────────────────────────┐
│                      Backend (Go/Gin)                       │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                   Middleware Layer                   │   │
│  │  • Auth Middleware  • Rate Limiter  • Recovery     │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                  Controller Layer                     │   │
│  │  • ContextHelper  • Request Validation            │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                    Service Layer                     │   │
│  │  • Business Logic  • State Machine  • Validation  │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                  Repository Layer                    │   │
│  │  • GenericRepository  • Custom Queries              │   │
│  └─────────────────────────────────────────────────────┘   │
└────────────────────────────┬────────────────────────────────┘
                            │
         ┌──────────────────┼──────────────────┐
         │                  │                  │
    ┌────▼────┐       ┌─────▼─────┐     ┌─────▼─────┐
    │  MySQL  │       │   Redis   │     │  External │
    │  (RDB)  │       │           │     │   APIs    │
    │ (Slave) │       │ • Cache   │     │  • LINE   │
    └─────────┘       │ • Queue   │     │  • R2     │
    ┌─────────┐       └───────────┘     └───────────┘
    │  MySQL  │
    │  (WDB)  │
    │(Master) │
    └─────────┘
```

## 2. 領域/模組劃分

### 2.1 核心模組

| 模組 | 命名空間 | 說明 |
|:---|:---|:---|
| **認證模組** | `auth` | 教師 LINE 登入、管理員 Email/Password 登入、JWT 驗證 |
| **教師模組** | `teacher` | 教師個人資料、技能、證照、標籤管理 |
| **中心模組** | `center` | 中心管理、課程、教室、假期設定 |
| **排課模組** | `scheduling` | 規則管理、例外申請、衝突驗證、循環編輯 |
| **通知模組** | `notification` | 站內通知、 LINE 通知、通知佇列 |
| **智慧媒合模組** | `smart_matching` | 人才搜尋、媒合建議、替代時段 |
| **匯出模組** | `export` | CSV/PDF 匯出、ICS 日曆匯出、圖片匯出 |
| **LINE Bot 模組** | `line_bot` | LINE Webhook、帳號綁定、歡迎訊息 |

### 2.2 控制器清單

| 控制器 | 檔案 | 職責 |
|:---|:---|:---|
| AuthController | `app/controllers/auth.go` | 登入、登出、Token 刷新 |
| TeacherController | `app/controllers/teacher.go` | 教師公開註冊 |
| TeacherProfileController | `app/controllers/teacher_profile.go` | 教師個人檔案 CRUD |
| TeacherScheduleController | `app/controllers/teacher_schedule.go` | 教師排課查詢、循環編輯 |
| TeacherSessionController | `app/controllers/teacher_session.go` | 教師課程筆記 |
| TeacherEventController | `app/controllers/teacher_event.go` | 私人行程管理 |
| TeacherExceptionController | `app/controllers/teacher_exception.go` | 例外申請（教師端） |
| TeacherInvitationController | `app/controllers/teacher_invitation.go` | 邀請回應 |
| AdminTeacherController | `app/controllers/admin_teacher.go` | 教師管理、合併、移除 |
| AdminCenterController | `app/controllers/admin_center.go` | 中心 CRUD |
| AdminRoomController | `app/controllers/admin_room.go` | 教室管理 |
| AdminCourseController | `app/controllers/admin_course.go` | 課程模板管理 |
| AdminHolidayController | `app/controllers/admin_holiday.go` | 假期管理 |
| AdminTermController | `app/controllers/admin_term.go` | 學期管理 |
| AdminUserController | `app/controllers/admin_user.go` | 管理員 CRUD、 LINE 綁定 |
| SchedulingController | `app/controllers/scheduling.go` | 排課規則管理、例外審核 |
| SmartMatchingController | `app/controllers/smart_matching.go` | 智慧媒合搜尋 |
| NotificationController | `app/controllers/notification.go` | 通知列表管理 |
| AdminNotificationController | `app/controllers/admin_notification.go` | 廣播通知 |
| ExportController | `app/controllers/export.go` | 資料匯出 |
| LineBotController | `app/controllers/line_bot.go` | LINE Bot Webhook |
| GeoController | `app/controllers/geo.go` | 縣市/區域資料 |
| OfferingController | `app/controllers/offering.go` | 班別管理 |
| TimetableTemplateController | `app/controllers/timetable_template.go` | 課表範本 |
| AdminResourceController | `app/controllers/admin_resource.go` | 教師資源管理 |

## 3. 分層職責邊界

### 3.1 層級職責表

| 層級 | 職責 | 禁止事項 |
|:---|:---|:---|
| **Controller** | Request 解析 → 呼叫 Service → 回傳 JSON | 寫入任何業務邏輯 |
| **Request** | 參數驗證（binding 標籤、CheckParam、CheckEnum） | 業務判斷 |
| **Service** | 核心業務邏輯、狀態機流轉、交易控制、依賴注入 | 直接操作資料庫 |
| **Repository** | 純粹 DB 操作（Find/Create/Update） | 任何業務判斷 |
| **Model** | 數據結構定義（GORM Tags） | 商業邏輯 |
| **Resource** | Model 轉換為輸出格式（DTO） | 修改資料狀態 |

### 3.2 資料流向

```
HTTP Request
     │
     ▼
┌────────────────────────────────────────────────────────────┐
│ Middleware                                                 │
│ • AuthMiddleware (JWT 驗證)                                │
│ • RateLimiter (頻率限制)                                   │
│ • RecoverMiddleware (panic 復原)                           │
│ • ResponseSanitizer (響應淨化)                             │
└────────────────────────────────────────────────────────────┘
     │
     ▼
┌────────────────────────────────────────────────────────────┐
│ Controller (Thin Controller)                               │
│ • ContextHelper 統一取值                                  │
│ • MustBindJSON() 參數綁定                                  │
│ • MustParamUint() URL 參數                                │
│ • Success()/ErrorWithInfo() 響應                           │
└────────────────────────────────────────────────────────────┘
     │
     ▼
┌────────────────────────────────────────────────────────────┐
│ Request (Validation)                                       │
│ • binding:"required" 必填欄位                            │
│ • 自定義驗證器 (CheckParam, CheckEnum)                    │
└────────────────────────────────────────────────────────────┘
     │
     ▼
┌────────────────────────────────────────────────────────────┐
│ Service (Business Logic)                                  │
│ • Triple Return Pattern: (data, *errInfo, error)          │
│ • Transaction Control                                     │
│ • Validation Engine                                       │
│ • State Machine Transitions                               │
└────────────────────────────────────────────────────────────┘
     │
     ▼
┌────────────────────────────────────────────────────────────┐
│ Repository (Data Access)                                   │
│ • GenericRepository[T] 通用 CRUD                          │
│ • Custom Queries                                          │
│ • RDB (Read) / WDB (Write) 分離                         │
└────────────────────────────────────────────────────────────┘
     │
     ▼
┌────────────────────────────────────────────────────────────┐
│ Model (Data Structure)                                    │
│ • GORM Tags 定義                                         │
│ • JSON 序列化/反序列化                                   │
│ • 自訂類型 (DateRange, RecurrenceRule, etc.)            │
└────────────────────────────────────────────────────────────┘
     │
     ▼
┌────────────────────────────────────────────────────────────┐
│ Database (MySQL)                                          │
│ • AutoMigrate 自動遷移                                    │
│ • Soft Deletes (gorm.DeletedAt)                           │
│ • Foreign Keys                                            │
└────────────────────────────────────────────────────────────┘
     │
     ▼
┌────────────────────────────────────────────────────────────┐
│ Resource (Response Transform)                             │
│ • Model → DTO 轉換                                       │
│ • 敏感資料遮蔽                                            │
└────────────────────────────────────────────────────────────┘
     │
     ▼
HTTP Response (JSON)
```

## 4. 錯誤處理流程

### 4.1 錯誤碼體系

錯誤碼格式：`AppID(1位) + 類型(2位) + 流水號(4位)`

| 類型 | 範圍 | 說明 |
|:---|:---|:---|
| **系統 (1)** | 10001-10999 | 系統級錯誤 |
| **資料庫 (2)** | 20001-20999 | SQL、交易錯誤 |
| **權限 (3)** | 30001-30999 | 認證、授權錯誤 |
| **業務 (4)** | 40001-40999 | 資源不存在、重複等 |
| **排課 (5)** | 50001-50999 | 重疊、緩衝、鎖定等 |
| **例外 (6)** | 60001-60999 | 例外申請、審核錯誤 |
| **檔案 (7)** | 70001-70999 | 上傳、類型錯誤 |
| **媒合 (8)** | 80001-80999 | 人才庫錯誤 |
| **LINE (9)** | 90001-90999 | LINE Bot 錯誤 |
| **管理員 (10)** | 100001-100099 | 管理員錯誤 |
| **資源鎖定 (11)** | 110001-110099 | 並發、鎖定錯誤 |
| **交易 (12)** | 120001-120099 | 交易執行錯誤 |

### 4.2 錯誤處理流程

```
Service Layer
     │
     │ return nil, app.Err.New(errInfos.NOT_FOUND), err
     ▼
Controller Layer
     │
     │ helper.ErrorWithInfo(errInfo)
     ▼
ContextHelper.ErrorWithInfo()
     │
     ├─ 檢查 errInfo 是否為 nil
     │
     ├─ 根據錯誤碼映射 HTTP Status
     │   • UNAUTHORIZED → 401
     │   • FORBIDDEN → 403
     │   • NOT_FOUND → 404
     │   • BAD_REQUEST → 400
     │   • SCHED_OVERLAP → 409 (Conflict)
     │   • SCHED_BUFFER → 409 (Conflict)
     │
     ▼
HTTP Response
{
  "code": 40001,
  "message": "找不到資源",
  "data": null
}
```

## 5. RDB/WDB 使用說明

### 5.1 讀寫分離架構

| 連線 | 用途 | 實作 |
|:---|:---|:---|
| **RDB** (Read Database) | 查詢操作，使用 Slave | `app.MySQL.RDB.WithContext(ctx)` |
| **WDB** (Write Database) | 寫入操作，使用 Master | `app.MySQL.WDB.WithContext(ctx)` |

### 5.2 Repository 使用模式

```go
// GenericRepository 自動處理 RDB/WDB
type GenericRepository[T models.IModel] struct {
    dbRead  *gorm.DB  // RDB (Slave)
    dbWrite *gorm.DB  // WDB (Master)
}

// 讀取操作 → RDB
func (rp *GenericRepository[T]) GetByID(ctx context.Context, id uint) (T, error) {
    return rp.dbRead.Table(rp.table).Where("id = ?", id).First(&data)
}

// 寫入操作 → WDB
func (rp *GenericRepository[T]) Create(ctx context.Context, data T) (T, error) {
    return rp.dbWrite.Table(rp.table).Create(&data)
}

// 交易操作
func (rp *GenericRepository[T]) Transaction(ctx context.Context, fn func(txRepo *GenericRepository[T]) error) error {
    return rp.dbWrite.Transaction(func(tx *gorm.DB) error {
        txRepo := &GenericRepository[T]{dbRead: tx, dbWrite: tx}
        return fn(txRepo)
    })
}
```

## 6. 驗證引擎（排課核心）

### 6.1 驗證層級

1. **Scope Check**：確保操作都在指定 `center_id` 下
2. **Hard Overlap Check（硬衝突）**：
   - 查詢時段 `[Start, End]` 內，該 `Teacher` 或 `Room` 是否已有其他 `Active Session`
   - 規則：`Existing.Start < New.End AND Existing.End > New.Start`
   - 若 True，直接報錯 `E_OVERLAP`（不可覆寫）
3. **Buffer Check（緩衝）**：
   - Room Buffer：`New.Start - Prev.End < current_course.room_buffer_min`
   - Teacher Buffer：`New.Start - Prev.End < current_course.teacher_buffer_min`
   - 若衝突且 `offering.allow_buffer_override = true` → 允許帶 `override=true` 送出

### 6.2 例外狀態機

```
[PENDING] ──教師撤回──→ [REVOKED]
 │
 ├── 管理員同意 ──→ [APPROVED] ──→ 發送 LINE 通知
 │
 └── 管理員拒絕 ──→ [REJECTED] ──→ 發送 LINE 通知
```

## 7. 環境變數

| 變數 | 說明 | 預設值 |
|:---|:---|:---|
| `APP_TIMEZONE` | 系統時區 | `Asia/Taipei` |
| `APP_ID` | 應用 ID（錯誤碼前綴） | `1` |
| `LOG_LEVEL` | 日誌級別 | `info` |
| `MYSQL_RDB_DSN` | MySQL 讀取連線 | - |
| `MYSQL_WDB_DSN` | MySQL 寫入連線 | - |
| `REDIS_ADDR` | Redis 位址 | `localhost:6379` |
| `NOTIFICATION_WORKER_ENABLED` | 通知 Worker | `false` |
| `ASYNQ_WORKER_ENABLED` | Asynq Worker | `false` |

## 8. 專案結構

```
/
├── main.go                    # Backend Entry Point
├── app/                      # Go 後端核心
│   ├── controllers/           # API 入口
│   ├── requests/              # 參數驗證
│   ├── services/              # 業務邏輯
│   ├── repositories/          # DB 存取
│   ├── resources/             # Response 轉換
│   ├── models/                # 數據模型
│   ├── servers/               # Server, Route & Middleware
│   ├── middleware/            # 中間件
│   └── pagination/            # 分頁 Helper
├── global/                    # 全域共用
│   ├── errInfos/              # 錯誤碼與訊息
│   ├── logger/                # 日誌
│   └── define.go              # 全域常數
├── database/                  # 資料庫
│   └── mysql/                 # MySQL 連線與遷移
├── configs/                   # 環境變數
├── frontend/                  # Nuxt 3 前端
├── pdr/                      # 規劃文件
└── testing/                   # 測試
```

---

*本文件基於實際程式碼生成，最後更新：2026-02-12*
