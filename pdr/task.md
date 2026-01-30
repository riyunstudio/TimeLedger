# TimeLedger 優化任務清單 (Task List)

## Phase 0: 已完成優化 (Completed)
- [x] **Teacher 控制器拆分**: 已拆分為 Event, Profile, Invitation, Session, Schedule, Exception 等子控制器。 <!-- id: 20 -->
- [x] **Teacher Service 層建立**: 已建立 `TeacherService`, `TeacherProfileService`, `PersonalEventService`。 <!-- id: 101 -->
- [x] **響應標準化**：統一使用 `helper.Success`, `helper.ErrorWithInfo`, `helper.Forbidden` 等方法，確保 API 回傳格式 100% 一致。
- [x] **ContextHelper 導入**: 全面取代手動 Context 提取。 <!-- id: 9 -->

## Phase 1: 架構加固與並發安全 (Architectural & Concurrency)
- [x] **GenericRepository 交易增強**: 優化 `Transaction` 與分層 Repo 繼承方式，修復並發設計隱憂。 <!-- id: 401 -->
- [x] **全項目交易審查**: 已全面優化 `ScheduleService` (Create/Update/Delete) 交易模式，確保 100% 執行緒安全。 <!-- id: 402 -->

## Phase 2: 領域完全解耦 (New Mode Transition)
- [x] **效能擴展與異步任務**: 已導入 Redis 快取 (Course, Room) 與 Asynq 異步通知系統。 <!-- id: 301 -->
    - [x] Center 領域抽出為 `CenterService` + `AdminCenterController`。
    - [x] Extraction of **Room** domain from `AdminResourceController`
    - [x] Extraction of **Course** domain from `AdminResourceController`
    - [x] Extraction of **Holiday** domain from `AdminResourceController`
    - [x] Extraction of **Offering** domain from `AdminResourceController`

## Phase 3: 性能規模化 (Performance & Scalability)
- [x] **Redis 分層緩存**: 實施中心資料與老師課表的 Redis 快取。 <!-- id: 18 -->
- [x] **異步任務隊列 (Asynq)**: 導入 `Asynq` 處理 LINE 通知與耗時 Job。 <!-- id: 6 -->
    - [x] 建立 `NotificationQueueService` 處理任務派發。
    - [x] 實作 LINE 通知異步 Enqueue 與 Payload 管理。
    - [x] 補齊 `main.go` 中的 Asynq Worker 啟動邏輯。
- [x] **N+1 解決方案**: 在 `ExpandRules` 中實施批次加載打破 N+1 查詢。 <!-- id: 17 -->

## Phase 4: 工程化與穩定性 (Stability & Engineering)
- [x] **Service 單元測試**: 為 `ScheduleExpansion` 核心邏輯編寫穩定性測試。 <!-- id: 102 -->
    - [x] 建立 `scheduling_expansion_test.go` 並測試批次查詢與例外處理。
- [x] **結構化日誌**: 已導入 `Zap` 結構化日誌並封裝於 `BaseService`。 <!-- id: 5 -->
- [x] **基類賦能 (BaseService)**: 封裝通用分頁與過濾邏輯。 <!-- id: 403 -->
    - [x] 實作動態分頁器 (Pagination Helper)。
    - [x] 實作動態過濾條件生成器 (Filter Helper)。

## Phase 5: 進階效能與架構美學 (Advanced Performance & Architecture)
- [x] **課表展開深度快取**: 實施 `ExpandRules` 的 Redis 分層快取與主動失效機制。 <!-- id: 501 -->
- [x] **DTO/Resource 層建立**: 導入 `Resource` 套件，完全隔離 DB Model 與 API 響應。 <!-- id: 502 -->
- [x] **AdminResource 終極解耦**: 已將 Invitation 移至專屬控制器，並引入 AdminTeacherResource。 <!-- id: 503 -->
- [x] **Repository 代碼大掃除**: 移除子 Repo 中重複的標準 CRUD 代碼。 <!-- id: 504 -->

## Phase 6: 上線準備與生產穩定性 (Go-Live & Production)
- [ ] **環境變數與安全檢查**: 確保生產環境 Secret 隔離。 <!-- id: 601 -->
- [ ] **遷移腳本安全化**: 將 `AutoMigrate` 改為手動觸發或生產環境禁用。 <!-- id: 602 -->
- [ ] **錯誤監控整合**: 考慮導入 Sentry 或類似工具追蹤 Panic。 <!-- id: 603 -->
- [ ] **API 效能基準測試**: 執行 Load Test 驗證快取與並發表現。 <!-- id: 604 -->
