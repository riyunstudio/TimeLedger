# 專案分析與優化建議報告

## 1. 專案完整度分析

### 概述
**TimeLedger** 是一個成熟的系統，前後端分離架構清晰。

- **前端**: 使用 **Nuxt 3** 建構，採用 **Pinia** 進行狀態管理，**TailwindCSS** 進行樣式設計，並有 **Playwright/Vitest** 測試框架。結構標準且組織良好。
- **後端**: 使用 **Go (Gin)** 建構，遵循分層架構 (Controller-Service-Repository)。整合了外部服務 (Line, Redis, RabbitMQ) 並擁有完整的 API 結構。

### 完整度評分: 高
- **功能面**: 程式碼涵蓋了排課、教師管理和 Line Bot 整合等複雜領域。
- **基礎設施**: 包含 Docker、Makefiles 和 CI/CD 配置。
- **測試**: 前端測試設置完善。後端有測試目錄，但覆蓋率需進一步驗證。

### 關鍵觀察
- **後端技術債**:雖然 *結構* 存在，但 *實作* 上有職責洩漏的問題。Controllers 過於複雜（"Fat Controllers"），處理了過多的業務邏輯和直接的資料庫操作。
- **前端機會**: 排課介面複雜。在效能和「感知速度」上有顯著的提升空間。

---

## 2. 前端使用者體驗 (UX) 優化方案

### A. 排課網格的虛擬滾動 (Virtual Scrolling)
**問題**: `ScheduleGrid.vue` 和 `ScheduleTimelineView.vue` 需要處理大量數據（日期 * 教室 * 場次）。渲染數千個 DOM 節點會導致卡頓和初始載入緩慢。
**方案**: 實作 **虛擬化 (Virtualization)** (例如使用 `vue-virtual-scroller`)。
- 只渲染目前視窗可見的項目。
- **效益**: 無論載入多少個月/年的數據，捲動都能保持流暢且即時渲染。

### B. 樂觀 UI 更新 (Optimistic UI)
**問題**: 像「預約時段」或「更新個人資料」等操作，目前可能顯示全域讀取轉圈並等待伺服器回應。這讓應用程式感覺「有距離感」且遲鈍。
**方案**: 在使用者互動後 *立即* 更新 UI，然後在背景與伺服器同步。
- 如果伺服器請求失敗，自動還原變更並顯示錯誤提示。
- **效益**: 應用程式的操作手感會像原生 App 一樣即時。

### C. PWA 與離線支援
**問題**: 作為「帳本」或管理工具，使用者（老師/管理員）可能需要在網路訊號不佳的地方查看課表。目前似乎只是標準的 SPA。
**方案**: 使用 `@vite-pwa/nuxt` 轉換為 **漸進式網頁應用 (PWA)**。
- 快取「今日課表」的 API 回應。
- 允許在離線狀態下以唯讀模式查看網格。
- **效益**: 在任何環境下都能保持可靠性和可存取性。

---

## 3. 後端重構與可維護性優化方案

### A. 重構 "Fat Controllers" (臃腫的控制器)
**問題**: `TeacherController` (及其他) 包含大量的邏輯。例如 `GetSchedule` 手動迭代規則、處理例外和格式化字串。`UpdateProfile` 直接呼叫 `ctl.app.MySQL.WDB`。
**方案**: **嚴格的服務層委派 (Strict Service Layer Delegation)**。
- 將所有排課邏輯（展開規則、套用例外）移動到 `SchedulingService`。
- **重構目標**: Controller 應該只負責解析請求 (`ctx.Bind`)，呼叫 Service 方法 (`service.GetTeacherSchedule(...)`)，然後返回 JSON 回應。

### B. 強制儲存庫模式 (移除直接資料庫存取)
**問題**: Controller 中存在像 `ctl.app.MySQL.WDB.WithContext(ctx).Delete(...)` 這樣的程式碼。這將 HTTP 層與特定的資料庫實作 (GORM/MySQL) 耦合在一起。
**方案**: 確保 **所有** 修改都通過 Repository 進行。
- 如果 Repo 中沒有該方法，則新增它。
- **效益**: 更容易進行單元測試（mock 介面）且程式碼更乾淨。更換 DB driver 時無需重寫 Controller。

### C. 集中式錯誤處理與回應 Middleware
**問題**: 錯誤處理重複性高 (`if err != nil { ctx.JSON(500, ...) }`)。這導致程式碼膨脹且錯誤訊息不一致。
**方案**: 建立 `Gin Middleware` 或 `ResponseHelper`。
- Service 返回的錯誤應該是「領域錯誤 (Domain Errors)」(例如 `ErrUserNotFound`)。
- Middleware 捕捉這些錯誤並嚴格對應到 HTTP 狀態碼 (404, 400, 500)。
- **效益**: 從 Controller 中移除數百行樣板錯誤檢查程式碼。

---

## 4. 進階優化建議 (Phase 2) - 深入分析後的新發現

### A. 資料庫效能: 解決 JSON_EXTRACT 導致的全表掃描 (嚴重)
**問題**: 在 `scheduling_expansion.go` 中，查詢使用了 `JSON_EXTRACT(effective_range, '$.start_date')`。
- **後果**: MySQL 無法對 JSON 內部的欄位使用一般索引，這會導致**每次查詢都進行全表掃描 (Full Table Scan)**。隨著資料量增加，效能會呈指數級下降。
**方案**:
1.  **結構變更**: 將 `start_date` 和 `end_date` 獨立為實體欄位 (Columns)。
2.  **虛擬欄位索引 (Generated Columns)** (如果不更改程式碼結構): 在 MySQL 中建立虛擬欄位並對其建立索引。
    ```sql
    ALTER TABLE schedule_rules ADD COLUMN start_date DATE GENERATED ALWAYS AS (effective_range->>'$.start_date') STORED;
    CREATE INDEX idx_rules_start_date ON schedule_rules (start_date);
    ```
- **效益**: 查詢速度將提升 100倍甚至 1000倍。

### B. 基礎設施: 導入專業任務隊列 (Job Queue)
**問題**: 目前的 `redis_queue_service.go` 是手刻的 `LPUSH`/`BRPOP` 迴圈。
- 缺乏視覺化儀表板 (Dashboard) 來監控失敗任務。
- 缺乏「指數退避 (Exponential Backoff)」重試機制 (目前是固定 5 秒)。
- 難以維護死信隊列 (Dead Letter Queue)。
**方案**: 導入 **Hibiken/Asynq** 或 **Machinery**。
- 這些是 Go 語言中成熟的非同步任務庫。
- **效益**: 提供完整的 Web UI 監控、優先級隊列、定時任務和穩定的重試機制。

### C. 可觀測性: 結構化日誌 (Structured Logging)
**問題**: 程式碼中大量使用 `fmt.Println` 或 `fmt.Printf` (如 `main.go` 和 Queue Service)。
- 在生產環境中，這些日誌難以被 ELK (Elasticsearch/Logstash/Kibana) 或 Datadog 解析。
- 無法輕易過濾 `ERROR` vs `INFO`。
**方案**: 導入 **Uber-go/Zap** 或 **Slog** (Go 1.21+ 內建)。
- 輸出為 JSON 格式：`{"level":"info", "ts":16800000, "msg":"User logged in", "user_id": 123}`。
- **效益**: 能夠在日誌管理系統中快速搜尋、過濾和設定警報。

---

## 5. 重構待辦清單 (Refactoring TODO List)

為了確保改動的安全性與可控性，建議按照以下優先順序進行：

### Phase 1: 止血與效能 (Immediate Fixes)
- [ ] **DB 索引優化 (Critical)**:
    - [ ] 修改 `schedule_rules` 資料表，新增 `start_date` 和 `end_date` 虛擬欄位 (Generated Columns)。
    - [ ] 建立索引 `idx_rules_effective_date`。
- [ ] **移除測試後門**:
    - [ ] 在 `middleware/auth.go` 中，移除 `mock-` token 的硬編碼後門，或僅允許在 `GO_ENV=test` 時啟用。

### Phase 2: 架構重整 (Structural Refactoring)
- [ ] **建立 Service 層邊界**:
    - [ ] 建立 `SchedulingService` 介面。
    - [ ] 將 `TeacherController.GetSchedule` 中的迴圈邏輯取出，放入 `ExpandRules` 方法中。
- [ ] **引入 Repo 模式**:
    - [ ] 掃描所有 Controller，找出 `ctl.app.MySQL.WDB` 的直接呼叫。
    - [ ] 將這些呼叫封裝進對應的 Repository (例如 `TeacherRepo.UpdateHashtags`)。

### Phase 3: 基礎設施升級 (Infrastructure)
- [ ] **結構化日誌**:
    - [ ] 引入 `Zap` logger。
    - [ ] 替換 `main.go` 與 Middleware 中的 `fmt.Println`。
- [ ] **非同步任務**:
    - [ ] 引入 `Asynq`。
    - [ ] 重寫 `notification_queue` 邏輯以使用新的 Job Queue。

### Phase 4: 前端優化 (Frontend Polish)
- [ ] **排課表虛擬化**:
    - [ ] 安裝 `vue-virtual-scroller`。
    - [ ] 重構 `ScheduleGrid.vue` 僅渲染可見區域。
- [ ] **PWA 支援**:
    - [ ] 配置 `@vite-pwa/nuxt`。
    - [ ] 設定 Service Worker 快取策略 (Stale-while-revalidate)。

### Phase 5: 使用者體驗與開發體驗優化 (Deep Dive)
- [ ] **Middleware Context 封裝 (DX)**:
    - [ ] 建立 `ContextHelper` (或 `UserContext` struct)。
    - [ ] 實作 `UserID()`, `UserRole()`, `MustGetUserID()` 等方法。
    - [ ] **目的**: 消除 Controller 中重複的 `ctx.GetUint(UserIDKey)` 與型別轉換程式碼。
- [ ] **匯出功能增強 (UX)**:
    - [ ] **ICS 行事曆訂閱**: 不要只給圖片，提供 `.ics` 連結讓老師同步到 Google/Apple Calendar。
    - [ ] **Server-Side Image Generation**: 改用後端 (如 `chromedp` 或 `canvas`) 生成高品質圖片，解決 `html2canvas` 跑版問題。
    - [ ] **客製化背景**: 允許老師上傳自己的照片當背景，增加分享意願。

### Phase 6: DRY 原則重構 (Code Consistency)
- [ ] **後端 Generic Repository (Go 1.18+)**:
    - [ ] 發現 `TeacherRepository` 與 `CenterRepository` 重複了大量的 CRUD 程式碼。
    - [ ] **方案**: 建立 `UsingGenericRepo[T any]`，實作統一的 `GetByID`, `Create`, `Update`, `Delete`。
    - [ ] **效益**: 減少 40% 以上的 Repository 樣板程式碼。
- [ ] **前端 API Client 重構**:
    - [ ] 發現 `useApi.ts` 中 `get/post/put` 方法重複了 Header 建構邏輯。
    - [ ] **方案**: 抽離 `request()` 私有方法，統一處理 Header 與 Error Parsing。
    - [ ] **Token 管理**: 將 `localStorage` key (`admin_token`, `teacher_token`) 統一管理，避免字串散落在各處。
