# TimeLedger 優化任務清單 (Task List)

## Phase 0: 已完成優化 (Completed)
- [x] **Teacher 控制器拆分**: 已拆分為 Event, Profile, Invitation, Session, Schedule, Exception 等子控制器。 <!-- id: 20 -->
- [x] **Teacher Service 層建立**: 已建立 `TeacherService`, `TeacherProfileService`, `PersonalEventService`。 <!-- id: 101 -->
- [x] **響應標準化**：統一使用 `helper.Success`, `helper.ErrorWithInfo`, `helper.Forbidden` 等方法，確保 API 回傳格式 100% 一致。
- [x] **ContextHelper 導入**: 全面取代手動 Context 提取。 <!-- id: 9 -->

## Phase 1: 架構加固與並發安全 (Architectural & Concurrency)
- [ ] **GenericRepository 交易增強**: 優化 `TransactionWithRepo` 與繼承連接方式。 <!-- id: 401 -->
- [ ] **全項目交易審查**: 確保所有多表寫入 Service 已實施 Transaction。 <!-- id: 402 -->

## Phase 2: 領域完全解耦 (New Mode Transition)
- [/] **AdminResource 控制器拆分**: 正進行中。 <!-- id: 201 -->
    - [x] Center 領域抽出為 `CenterService` + `AdminCenterController`。
    - [ ] Room 領域抽出為 `RoomService` + `AdminRoomController`。
    - [ ] Course 領域抽出為 `CourseService` + `AdminCourseController`。
    - [ ] Holiday & Invitation 領域抽出。
- [x] **全局控制器標準化**: 已完成 AdminUser, SmartMatching, Notification 控制器的 ContextHelper 整合。 <!-- id: 302 -->

## Phase 3: 性能規模化 (Performance & Scalability)
- [ ] **Redis 分層緩存**: 實施中心資料與老師課表的 Redis 快取。 <!-- id: 18 -->
- [ ] **異步任務隊列**: 導入 `Asynq` 處理 LINE 通知與耗時 Job。 <!-- id: 6 -->
- [x] **N+1 解決方案**: 在 `ExpandRules` 中實施批次加載打破 N+1 查詢。 <!-- id: 17 -->

## Phase 4: 工程化與穩定性 (Stability & Engineering)
- [ ] **Service 單元測試**: 為 `ApplyTemplate` 等核心邏輯編寫穩定性測試。 <!-- id: 102 -->
- [ ] **結構化日誌**: 導入 `Zap` 或 `Slog` 取代 `fmt` 打印。 <!-- id: 5 -->
- [ ] **基類賦能 (BaseService)**: 封裝通用分頁與過濾邏輯。 <!-- id: 403 -->
