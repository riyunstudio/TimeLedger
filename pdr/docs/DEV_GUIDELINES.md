# TimeLedger 開發準則 (Development Guidelines)

本文件依據 `backend_module_generator.md` 與 專案架構規劃 制定，所有開發者必須嚴格遵守。

---

## 1. 專案結構 (Monorepo - Backend Centric)

本專案以 **後端框架為根目錄**，前端應用作為子目錄共存。

```text
/
├── app/                    # Go 後端核心 (Monolithic)
│   ├── controllers/        # API 入口
│   ├── requests/           # 參數驗證
│   ├── services/           # 業務邏輯
│   ├── repositories/       # DB 存取
│   ├── resources/          # Response 轉換
│   └── models/             # 數據模型
├── global/                 # 全域共用 (ErrInfos, Config)
├── backend_specs/          # 模組產生器設定檔
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
4.  **Repository**: 純粹的 DB 操作 (`Find`, `Create`, `Update`). 錯誤直接回傳 `gorm` error。

### 2.2 錯誤處理 (Error Handling)
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
*   **Pages**:
    *   `/admin/*`: 使用 `layouts/admin.vue` (側邊欄佈局).
    *   `/teacher/*`: 使用 `layouts/default.vue` (底部導航佈局).
    *   `/auth/*`: 使用 `layouts/blank.vue` (全螢幕).
*   **Components**:
    *   `base/`: 原子元件 (Button, Input, Card).
    *   `business/`: 業務元件 (ScheduleGrid, MemberList).

### 3.3 互動與狀態
*   **API Call**: 統一使用 `useFetch` 或封裝過的 `$api` composable。
*   **Error Handling**: API 錯誤需攔截 `error.value` 並顯示 Toast。
*   **Optimistic UI**: 排課操作需先更新 UI，失敗再 Rollback。

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
