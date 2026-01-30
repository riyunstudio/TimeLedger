# Cursor Agent 優化指令集 (Command Sheet)

複製以下指令貼入 Cursor (Cmd+K 或 Chat)，即可驅動 AI 自動執行優化。

---

### 指令 1：拆分 Scheduling 模組 (當前首選)
> **CONTEXT: @pdr/task.md @pdr/ARCHITECTURAL_OPTIMIZATION_GUIDE.md @app/controllers/scheduling.go**
> 
> 請參考 Teacher 模組的成功經驗，將 `app/controllers/scheduling.go` 進行深度重構。
> 1. 建立 `app/services/scheduling.go` (ScheduleService)。
> 2. 將控制器中的「衝突檢查」、「排課邏輯」、「例外處理」搬移至 Service。
> 3. 將 `SchedulingController` 改為 Thin Controller，並使用 `@controllers/context_helper.go` 管理上下文。
> 4. 注意：這是一個大檔案，請先從「建立 Service 基礎結構」開始。

---

### 指令 2：效能優化 - 解決 N+1 查詢
> **CONTEXT: @app/services/scheduling_expansion.go @app/repositories/schedule_exception.go**
> 
> 請優化 `scheduling_expansion.go` 中的 `ExpandRules` 方法，解決 N+1 查詢問題。
> 1. 識別迴圈中的 DB 查詢動作。
> 2. 在進入日期迴圈前，利用 ID 列表進行 `In` 查詢批次抓取。
> 3. 使用 `map` 在內存中進行數據組裝。
> 4. 提供重構後的效能比較預測。

---

### 指令 3：Repository 泛型「極簡化」遷移
> **CONTEXT: @app/repositories/generic.go @app/repositories/[target_repo].go**
> 
> 請執行 Repository 泛型遷移與「極簡化」清理：
> 1. 使 `[Target]Repository` 繼承 `GenericRepository[models.Target]`。
> 2. **關鍵精簡**：刪除檔案中所有 `GenericRepository` 已內建的標準方法（如 `GetByID`, `ListByCenterID`, `DeleteByID` 等帶有 CenterScope 的版本）。
> 3. 對於特定查詢（如 `ListActive`），利用基類的 `FindWithCenterScope` 加過濾條件重構。
> 4. 確保 `New` 構造函數正確初始化。

---

### 指令 4：領域資源拆分 (Domain Decomposition)
> **CONTEXT: @pdr/task.md @app/controllers/admin_resource.go**
> 
> 請執行 AdminResource 的「[領域名稱，如 Room/Course/Holiday]」領域拆分：
> 1. 建立 `app/services/[domain].go` 與 `app/controllers/admin_[domain].go`。
> 2. 將 `AdminResourceController` 中所有關於該領域的方法移至新 Service/Controller。
> 3. 確保 Service 層實作時包裝了 `AuditLog` 記錄。
> 4. 在 `apis/base.go` 或 `route.go` 更新路由，並刪除舊代碼。

---

### 指令 5：Offering 模組現代化 (Service + Generic Repo)
> **CONTEXT: @app/controllers/offering.go @app/repositories/offering.go**
> 
> 請對 Offering 模組進行架構升級：
> 1. 使 `OfferingRepository` 繼承 `GenericRepository[models.Offering]` 並刪除冗餘 CRUD。
> 2. 建立 `app/services/offering.go` 並搬移 Controller 中的複雜邏輯 (如 CopyOffering)。
> 3. 重構 `OfferingController` 為 Thin Controller 模式。

---

### 指令 6：排課模板服務化 (Logic Extraction)
> **CONTEXT: @app/controllers/timetable_template.go**
> 
> 請執行 `TimetableTemplate` 邏輯抽離：
> 1. 建立 `app/services/timetable_template.go`。
> 2. 將 `ApplyTemplate` 的核心流程（日期計算、規則生成）移至 Service。
> 3. 使用 `ContextHelper` 標準化控制器中的參數提取。

---

### 指令 7：全局控制器標準化 (Standardization Sweep)
> **CONTEXT: @app/controllers/admin_user.go @app/controllers/smart_matching.go @app/controllers/notification.go**
> 
> 請對這些控制器執行「Teacher 模式標準化」：
> 1. 將所有手動提取 `centerID`, `userID` 的 `if exists` 塊替換為 `ctl.getCenterID(ctx)` (來自 ContextHelper)。
> 2. 統一錯誤回傳格式為 `ctl.respondError(ctx, code, msg)` 或 `helper.ErrorWithInfo`。
> 3. 確保所有 DB 操作都透過對應的 Service 進行。

---

### 指令 8：Service 原子性與事務性優化
> **CONTEXT: @app/services/[target_service].go @app/repositories/generic.go**
> 
> 請優化 Service 的原子性與錯誤處理：
> 1. 對於涉及多個 Repository 寫入的操作，導入 `db.Transaction`。
> 2. 定義精細化錯誤常量（如 `errInfos.SCHED_OVERLAP`）。
> 3. 重構 Controller 響應，根據錯誤碼回傳 400 或 409 狀態碼（使用 `helper.ErrorWithInfo`）。

---

### 指令 9：性能規模化 (Performance Scaling)
> **CONTEXT: @app/services/[target_service].go @global/redis**
> 
> 請進行性能優化：
> 1. 為高頻讀取資源 (如 Center 設置、今日排課) 導入 Redis 快取。
> 2. 實施「寫入時清除 (Cache Aside)」策略。
> 3. 將耗時通知 (LINE/Email) 或檔案生成動作改為 `Asynq` 異步任務。

---

### 指令 10：工程化與穩定性 (Stability & Engineering)
> **CONTEXT: @app/services/base.go @app/services/[target_service].go**
> 
> 請提升代碼工程化品質：
> 1. 在 `BaseService` 封裝通用的分頁 (Pagination) 與動態過濾邏輯。
> 2. 使用 `Testify` 為核心業務函數 (如排課衝突檢查) 撰寫單元測試。
> 3. 將 `fmt.Println` 替換為結構化日誌記錄。

---

### 指令 11：Repository 交易安全性加固
> **CONTEXT: @app/repositories/generic.go @app/repositories/[target_repo].go**
> 
> 請修正 Repository 的交易設計：
> 1. 修改 `Transaction` 方法，確保在交易閉包中回傳一個「全新的 Repo 實例」而非修改現有實例。
> 2. 為領域 Repository (如 `ScheduleRuleRepository`) 實作專屬交易入口，確保在交易中仍能調用自定義方法。
> 3. 參考 `offering.go` 的交易模式進行重構。


### 指令 12：課表展開深度快取 (Schedule Expansion Caching)
> **CONTEXT: @app/services/scheduling_expansion.go @app/services/cache.go**
> 
> 請完成課表展開效能優化：
> 1. 在 `ScheduleService` (或相關擴展服務) 導入課表快取。
> 2. 快取鍵設計需包含 `center_id` 或 `teacher_id` 以及日期範圍。
> 3. 實施衝突觸發 (In-place Invalidation)：當 `ScheduleRule` 或 `ScheduleException` 異動時，精準清除受影響的快取。

---

### 指令 13：DTO 與 Resource 映射層 (DTO & Resource Separation)
> **CONTEXT: @app/models/ [target_model].go @app/controllers/ [target_controller].go**
> 
> 請實施數據傳輸對象 (DTO) 模式：
> 1. 在 `app/resources` 建立對應領域的 Resource 結構 (如 `RoomResource`)。
> 2. 移除 Controller 直接返回 DB Model 的行為。
> 3. 在 Service 或映射層處理欄位過濾與格式轉換 (如時間格式化)。

---

### 指令 14：AdminResource 控制器終極解耦 (Final Controller Cleanup)
> **CONTEXT: @app/controllers/admin_resource.go**
> 
> 請完成最後的關注點分離：
> 1. 將 `Invitation` 統計與列表 API 遷移至 `AdminInvitationController`。
> 2. 將 `TeacherNote` 相關 API 遷移至 `AdminTeacherController`。
> 3. 確保所有路由配置同步更新，維持 API 不變。

---

### 指令 15：Repository 代碼大掃除 (Repository Code Refinement)
> **CONTEXT: @app/repositories/ [target_repo].go**
> 
> 根據「極簡化原則」優化 Repository：
> 1. 凡是 `GenericRepository` 已有的標準 CRUD (GetByID, List, Delete) 手寫實作，必須 **100% 移除**。
> 2. 子類 Repo 僅保留具備特殊業務邏輯或複雜查詢的方法。
> 3. 驗證所有 Service 對 Repo 的調用是否依然正確。

---


### 💡 如何高效協作？
1. **先 Service 後 Controller**：要求 AI 先產生 Service，再修改 Controller。
2. **小步快跑**：一次只執行一個「指令」，完成後驗證通過再進行下一個。
3. **提及上下文**：務必使用 `@` 提及對應的 `task.md` 與目標代碼檔案。
