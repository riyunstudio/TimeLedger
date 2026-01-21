# 技術實作藍圖 (Implementation Blueprint)

本文件依據 `pdr/docs/DEV_GUIDELINES.md` 定義之專案架構，將功能需求精確對映至代碼導航路徑。

---

## 1. 核心功能模組路徑 (Feature Path Mapping)

### 1.1 排課驗證引擎 (Scheduling & Validation)
*   **Decoupled Logic (核心算法)**: 
    *   `app/scheduling/validator/`: `overlap.go`, `buffer.go` (純邏輯，不依賴 DB)。
    *   `app/scheduling/expander/`: `recurrence.go` (週期展開處理)。
*   **Application Layer (應用層)**:
    *   `app/services/schedule_service.go`: 業務流程編排。
    *   `app/repositories/schedule_repo.go`: SQL 查詢與併發鎖 (`FOR UPDATE`)。
    *   `app/models/schedule.go`: 資料結構定義。
    *   `app/servers/route.go`: 路由定義。
    *   `app/controllers/schedule_controller.go`: API 入口。
    *   `app/requests/schedule_request.go`: 參數校驗。
    *   `app/resources/schedule_resource.go`: Response 轉換。

### 1.2 身份與權限 (Auth & RBAC)
*   **Decoupled Logic**: `pkg/auth/line/` (LINE 協議封裝)。
*   **Application Layer**:
    *   `app/services/auth_service.go`: 處理登入/註冊流程。
    *   `app/servers/auth_middleware.go`: 強執權限校驗 (Middleware)。
    *   `app/repositories/user_repo.go`: 用戶數據存取。

### 1.3 媒合與人才搜尋 (Matching & Talent)
*   **Decoupled Logic**: `app/services/logic/matcher/` (計算評分與過濾)。
*   **Application Layer**:
    *   `app/services/talent_service.go`: 實作搜尋與媒合業務。
    *   `app/controllers/talent_controller.go`: 外部接口。

### 1.4 非同步任務與通知 (Workers & Notifications)
*   **Worker Engine**: `app/jobs/` or `app/workers/` (基於 Redis/RabbitMQ)。
    *   `daily_reminder.go`: 每日晚上發送明日課程預告。
    *   `emergency_notifier.go`: 課程異動即時推播。
    *   `invitation_worker.go`: 處理邀請碼發送。
*   **Logic Layer**: `pkg/lineutil/` (Flex Message 樣版與 LINE API 封裝)。

### 1.5 檔案儲存與資源 (Storage & Assets)
*   **Service Layer**: `libs/storage/` (封裝 S3/OSS/本地 儲存介面)。
*   **Controller Layer**: `PATCH /teacher/me/avatar` 與 `POST /teacher/me/certificates` 直接調用 Storage Service。

---

## 2. 標準實作序列 (Implementation Sequence) - **Feature-by-Feature TDD**

為了確保開發不歪掉與中斷續航，必須按「功能單位」進行開發，嚴禁跨功能大規模編碼：

### Step 0: 準備與契約稽核
*   核對 `pdr/API.md` 與 `pdr/Integration_Playbook.md`。
*   **Gap Detection**: 先更新 PDR，再寫代碼。

### Step 1: 垂直切片實作 (Backend Cycle)
1.  **DB Migration**: 核對 `Mysql.md` 並執行目標 Stage 的資料表遷移。
2.  **TDD Baseline**: 撰寫 Service/Logic 層的單元測試 (Mock 驅動)。
3.  **Service Logic**: 實作業務邏輯，確保 TDD 100% 通過。
4.  **API Handler**: 依據 `API.md` 定義路由與 Controller。
5.  **Audit & Commit**: 執行 `POSTMAN/CURL` 初步驗證後 Git Commit。

### Step 2: 垂直切片實作 (Frontend Cycle)
1.  **Contract Sync**: 建立 TS Interface 與 API Composable。
2.  **UI Atomic**: 根據 `UiUX.md` 實作組件 (優先使用 `Base` 共用組件)。
3.  **State Logic**: 串接 API 並處理加載、錯誤與樂觀更新狀態。
4.  **Integration Test**: 執行 `Integration_Playbook.md` 中的目標測試案例。
5.  **Continuity Check**: 更新 `progress_tracker.md` 並 Git Commit。

(完成後，重複上述 Step 1-2 於下一個功能)

---

## 3. 共用工具類 (Shared Utilities)

*   **時間處理**: `pkg/utils/timeutil/` (Asia/Taipei)。
*   **標籤處理**: `pkg/utils/tagutil/` (Hashtag 清理)。
*   **資料處理**: `pkg/utils/jsonutil/` (JSON Merge/Validate)。
*   **通知封裝**: `pkg/utils/lineutil/` (Flex Message)。
