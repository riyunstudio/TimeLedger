# TimeLedger 開發準則 (Development Guidelines)

本文件依據 `backend_module_generator.md` 與 專案架構規劃 制定，所有開發者必須嚴格遵守。

---

## 1. 專案結構 (Monorepo - Backend Centric)

本專案以 **後端框架為根目錄**，前端應用作為子目錄共存。

```text
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
├── configs/                # 系統環境變數配置 (env.go)
├── grpc/                   # gRPC 定義與實作
├── rpc/                    # RPC 通訊組件
├── testing/                # 自動化測試與 Mocks
├── frontend/               # Nuxt 3 前端應用 (User + Admin)
│   ├── pages/
│   ├── components/
│   └── nuxt.config.ts
├── pdr/                    # 規劃文件 (PDR, Specs, Analysis)
├── main.go                 # Backend Entry Point
└── docker-compose.yml      # 本地開發環境
```

---

## 2. 後端開發規範 (Backend Guidelines)

**核心精神**：嚴格分層，禁止邏輯洩漏。

### 2.1 分層職責 (Strict Layering)
1.  **Controller**: 僅負責 `Request Parsing` -> `Call Service` -> `Response JSON`。**禁止寫業務邏輯**。
2.  **Request**: 負責所有參數驗證 (`CheckParam`, `CheckEnum`)。錯誤必須回傳 `errInfo`。
3.  **Service**: 核心業務邏輯、狀態機流轉 (`Guard`)、交易控制 (`Transaction`)。
    - **依賴注入**：Service 應持有 `app.App`、`Repository` 與 `Resource` 的指針。
    - **方法簽名**：`func (s *S) Method(ctx, req) (datas any, errInfo *errInfos.Res, err error)`。
4.  **Repository**: 純粹的 DB 操作 (`Find`, `Create`, `Update`)。
    - **Isolation Guard**: 所有查詢（除了跨租戶的 Teacher 私人行情外）必須強制在 `WHERE` 子句中包含 `center_id`。
    - **No Logic**: 禁止在 Repository 寫入任何業務判斷。
5.  **Model**:
    - **GORM Tags**: 必須定義正確的 `gorm:"primaryKey;autoIncrement"` 與關聯標籤。
    - **Time Standard**: `created_at` 與 `updated_at` 必須由 GORM 自動管理或統一使用 Unix Timestamp。
6.  **Resource**: 負責將 Model 轉換為輸出格式 (DTA)。

### 2.2 服務層標竿代碼 (Service Template Reference)
開發新的 Service 時，必須嚴格遵守以下 `app/services/user.go` 的實做模式：

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

### 2.3 錯誤處理 (Error Handling)
*   所有對外錯誤代碼必須定義於 `global/errInfos`。
*   格式：`MODULE_REASON` (例: `SCHEDULE_OVERLAP`)。
*   Controller 統一使用 `ctl.JSON` 回傳標準格式。

### 2.3 命名慣例 (Naming)
*   **Module**: `snake_case` (例: `payment_rule`).
*   **Table**: `snake_case` + plural (例: `payment_rules`).
*   **Go Struct**: `PascalCase` (例: `PaymentRule`).
*   **JSON Field**: `snake_case` (例: `payment_rule_id`).

---

## 3. 前端開發規範 (Frontend Guidelines)

**核心精神**：Mobile First，模組化開發。

### 3.1 技術堆疊
*   **Framework**: Nuxt 3 (Vue 3 + TS).
*   **Styling**: Tailwind CSS (Utility-first).
*   **State**: Pinia (Store).
*   **Icons**: Phosphor Icons / Heroicons.

### 3.2 目錄結構規範
*   **Pages**:## 5. 原子化開發與續航規範 (Atomic Workflow & Continuity)

為了確保 AI 在開發中斷後能無縫接軌，且避免大規模開發導致邏輯偏移，必須遵守以下鐵律：

### 5.1 段落式開發 (Segmented Dev)
1.  **限縮範圍 (Small Scope)**：一次僅開發一個獨立功能 (如：僅處理「老師邀請碼生成」)。禁止同時改動多個不相關模組。
2.  **TDD 順序**：
    *   **Step A**: 撰寫後端對應 Logic 的單元測試 (Unit Test with Mocks)。
    *   **Step B**: 實作後端 Service / logic。
    *   **Step C**: 實作後端 Handler (API)。
    *   **Step D**: 實作前端組件與排版。
    *   **Step E**: 實作前端 API 對接與狀態機。

### 5.2 原子化 Commit (Atomic Commits)
*   **頻率**：完成上述 Step A~C (後端完成且測試通過) 必須進行一次 Commit。完成 Step D~E (前端完成) 必須再次 Commit。
*   **紀錄格式**：
    ```bash
    git commit -m "feat(scope): 描述 (Ref: PDR章節)"
    # 範例
    git commit -m "feat(auth): 實作老師 LINE 註冊邏輯 (Ref: P3.1)"
    ```

### 5.3 續航日誌 (Continuity Logs)
*   **規則**：在每次功能段落結束時，必須更新 `pdr/progress_tracker.md`，寫下當前「已完成的細節」與「下一個具體步驟」。這能讓重新啟動的 AI 秒速恢復語境。
    *   `/admin/*`: 使用 `layouts/admin.vue` (側邊欄佈局).
    *   `/teacher/*`: 使用 `layouts/default.vue` (底部導航佈局).
    *   `/auth/*`: 使用 `layouts/blank.vue` (全螢幕).
*   **Components**:
    *   `base/`: 原子元件 (Button, Input, Card).
    *   `business/`: 業務元件 (ScheduleGrid, MemberList).

### 3.3 互動與狀態
- **API Call**: 統一使用 `useApi` composable。
- **Type Safety (Crucial)**: 
  - 所有 API 請求與回傳必須定義 **TypeScript Interfaces**。
  - **禁止** 在頁面中直接使用 `any` 或未定型的 JSON。
- **Error Handling (API Error Standard)**:
  - **禁止** 在 API 解析層攔截非 2xx 狀態碼而拋出通用錯誤；必須嘗試解析 Body 以取得後端詳細訊息。
  - **優先權**：後端回傳的 `message` > `NUMERIC_ERROR_CODE_MAP` 對應訊息 > `ERROR_MESSAGES` 常量 > 通用系統錯誤。
  - **映射要求**：所有後端數字錯誤碼（如 160006）必須在 `frontend/constants/errorCodes.ts` 的 `NUMERIC_ERROR_CODE_MAP` 中註冊。
- **狀態顯示**：組件必須處理 `isLoading` 與 `hasError` 狀態，並將 `error_handler` 處理後的 `errorMessage` 綁定至 UI。

---

## 4. Git 工作流 (Git Workflow)

*   **Main Branch**: `main` (Production Ready).
*   **Dev Branch**: `dev` (Integration).
*   **Feature Branch**: `feat/<module-name>` (例: `feat/schedule-engine`).
*   **Commit Message**:
    *   `feat: add new login page`
    *   `fix: resolve buffer calculation bug`
    *   `docs: update API spec`
    *   `chore: update dependencies`

---

## 5. 環境變數 (Environment Variables)

請複製 `.env.example` 為 `.env`。

```bash
# Backend
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=password
REDIS_HOST=localhost

# Frontend
NUXT_PUBLIC_API_BASE=http://localhost:8080/timeledger/api
```
