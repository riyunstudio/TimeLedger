# 專案優化建議書 (後端專注版)

根據對專案代碼的深入分析，以下列出八個針對後端架構、效能與可維護性的核心優化方向。

## 後端優化建議

### 1. 解決 N+1 查詢以提升 API 效能 (Performance)

**問題描述**：
在 `app/services/scheduling_expansion.go` 的 `ExpandRules` 方法中，存在典型的 N+1 查詢問題。程式碼會在一個日期迴圈內，針對每一天、每一個規則呼叫 `s.exceptionRepo.GetByRuleIDAndDateStr(...)` 查詢資料庫。當排課規則和查詢日期範圍變大時，資料庫查詢次數會呈現指數級增長，嚴重拖慢 `GetSchedule` 等關鍵 API 的回應速度。

**優化方案**：
- **批次查詢 (Batch Fetching)**：在進入迴圈前，根據所有的 `RuleIDs` 和查詢的 `DateRange`，一次性從資料庫撈取所有相關的 `ScheduleException`。
- **記憶體映射 (In-Memory Mapping)**：將撈取出的例外資料轉換為 Map 結構（例如 `map[RuleID]map[DateString][]Exception`），在迴圈中直接以 O(1) 的時間複雜度讀取，消除迴圈內的 I/O 操作。

### 2. 重構 Fat Controller 提升可維護性 (Maintainability)

**問題描述**：
`app/controllers/teacher.go` 檔案過於龐大（超過 3500 行），承載了過多職責，包括：個人資料管理、課表查詢、請假申請、標籤搜尋等。這違反了單一職責原則 (SRP)，導致代碼閱讀困難，多人協作時極易產生 Merge Conflict，且測試難以編寫。

**優化方案**：
- **拆分 Controller**：依據業務領域將 `TeacherController` 拆分為多個獨立的 Controller：
    - `TeacherProfileController`: 專注處理個人資料、證照、技能標籤。
    - `TeacherScheduleController`: 專注處理課表查詢、排課規則展開。
    - `TeacherExceptionController`: 專注處理請假、調課申請流程。
- **統一路由管理**：在路由層統一組合這些子 Controller，保持 API 路徑不變但內部結構清晰。

### 3. 解耦合與依賴注入 (Decoupling & DI)

**問題描述**：
目前的 Controller 直接依賴具體的 Service 和 Repository 結構體（如 `*repositories.TeacherRepository`）。這種強耦合設計使得單元測試變得非常困難，因為無法輕易地用 Mock 物件替換真實的資料庫操作。此外，`NewTeacherController` 中充斥著大量手動初始化的代碼，難以管理。

**優化方案**：
- **定義介面 (Interface)**：為所有的 Service 和 Repository 定義明確的 Interface，宣告其公開方法。
- **依賴介面**：Controller 結構體改為依賴這些 Interface。
- **引入與標準化 DI**：考慮引入 `google/wire` 或 `uber-go/dig` 等依賴注入工具，自動管理依賴關係圖，減少樣板代碼並提升模組的可測試性。

### 4. 引入 Redis 快取策略 (Caching Strategy)

**問題描述**：
課表計算 (`ExpandRules`) 邏輯複雜，涉及規則展開、假日排除、例外處理與跨日計算。目前每次 API 請求都會重新執行這些繁重的計算，即使資料沒有變更。隨著使用者增加，這將成為系統的效能瓶頸。

**優化方案**：
- **實作快取層**：使用 Redis 快取計算後的課表結果。Key 可以設計為 `schedule:teacher:{id}:{start_date}:{end_date}`。
- **智慧失效 (Smart Invalidation)**：
    - **TTL**：設定合理的過期時間 (如 1 小時)。
    - **事件觸發**：當有新的 `ScheduleRule` 或 `ScheduleException` 被建立/修改時，主動清除相關聯老師或中心的快取，確保資料一致性。

### 5. 採用事件驅動架構 (Event-Driven Architecture)

**問題描述**：
目前的業務邏輯耦合度高。例如 `CreateException` 方法中，除了建立資料，還直接呼叫了 `notificationQueue` 發送 LINE 通知，並寫入 Audit Log。若未來需新增「寄送 Email」或「推播 App 通知」，必須直接修改核心 Service 代碼，違反開放封閉原則 (OCP)。

**優化方案**：
- **引入 Event Bus**：建立一個簡單的內部事件匯流排 (基於 Go Channel) 或引入外部 Message Queue。
- **發布/訂閱模式**：核心 Service 僅負責「發生了什麼事」（如 `ExceptionCreated`），並發布事件。
- **非同步處理**：獨立的 Handler (如 `NotificationHandler`, `AuditHandler`)訂閱這些事件並執行副作用。這樣不僅降低了耦合，還能避免通知服務的延遲影響主流程的回應速度。

### 6. 領域邏輯封裝 (Domain Logic Encapsulation)

**問題描述**：
目前的專案採用「貧血模型 (Anemic Domain Model)」，即 `models` 只包含資料結構，而所有的業務規則（如「檢查是否可請假」、「計算課程費用」）都散落在 `Services` 甚至是 `Controllers` 中。這導致邏輯重複且難以復用。

**優化方案**：
- **充血模型**：將屬於該實體的業務邏輯移回 Model struct 中。
    - 例如：`func (r *ScheduleRule) IsActive(date time.Time) bool`
    - 例如：`func (t *Teacher) CanTeach(skillID uint) bool`
- **好處**：提高代碼的內聚性 (Cohesion)，讓業務規則更集中，Service 層只需負責協調工作流程。

### 7. DTO 模式分離 (Data Transfer Object Pattern)

**問題描述**：
API 的輸入輸出經常直接使用資料庫 Model (如 `binding:"json"` 直接掛在 Model struct 上)。這會導致：
1.  **安全性風險**：若前端惡意傳入 `is_admin=true`，可能會被意外寫入資料庫 (Mass Assignment)。
2.  **API 契約不穩**：資料庫欄位修改會直接破壞前端 API 格式。

**優化方案**：
- **引入 DTOs**：在 `app/requests` 和 `app/resources` 中明確定義 API 專用的 Request/Response struct。
- **映射層**：使用 `automapper` 或手寫轉換函式在 DTO 與 Model 之間進行轉換。確保 API 合約 (Contract) 與資料庫結構 (Schema) 解耦。

### 8. 統一錯誤處理機制 (Unified Error Handling)

**問題描述**：
目前的錯誤處理大多是 `if err != nil { return err }`，缺乏上下文資訊。且 Controller 層直接回傳 HTTP Status Code，導致業務邏輯中的錯誤（如「餘額不足」）與技術錯誤（如「DB 連線失敗」）難以區分與統一格式化。

**優化方案**：
- **自定義錯誤類型**：定義 `AppError` 結構，包含 `Code` (業務錯誤碼), `Message`, `HTTPStatus`。
- **Centralized Error Handler**：在 Middleware 層統一捕捉錯誤並格式化為標準 JSON 回應。
- **Wrap Errors**：使用 `fmt.Errorf("failed to fetch user: %w", err)` 包裹錯誤，保留完整的 Stack Trace 與上下文，便於 Debug。

---

## 🚀 深度優化專題 (Deep Refinement Topics)

### 9. Repository 泛型模式「再強化」(Strengthened Pattern)
**目標**：消除繼承泛型基類後仍存在的冗餘。
*   **極簡化原則**：凡是 `GenericRepository` 已有的方法（如 `GetByIDWithCenterScope`, `FindWithCenterScope`, `DeleteByIDWithCenterScope`），子類 Repo 應 **100% 移除** 手寫實現。
*   **成果預期**：`course.go` 等檔案可縮減 70% 以上的代碼量，減少維護成本。

### 10. Service 層事務與安全性 (Transactional Safety & Error Codes)
**目標**：解決排課、資源分配等涉及多表操作的原子性問題。
*   **事務包裹 (Transactions)**：在涉及「建立規則 + 生成場次」等連鎖操作時，必須在 Service 層使用 `db.Transaction`，避免執行一半失敗導致的髒數據（Dangling Data）。
*   **精細化錯誤碼**：Service 不應只回傳錯誤訊息，應定義具體的業務錯誤常量（如 `ErrOverlapConflict`），讓 Controller 能更準確地調用 `ContextHelper` 的響應工具。
*   **Controller 瘦身**：全面移除 Controller 內定義的 `requireCenterID` 等重複工具，改用 `helper.MustCenterID()` 等標準內建方法。

