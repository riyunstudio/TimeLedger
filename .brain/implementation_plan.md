# Implementation Plan - Academic Terms & Resource Occupancy View

Introduce "Terms" as named date ranges for center management and provide a specialized weekly calendar view to analyze teacher and room utilization based on scheduling rules.

## UX Optimizations & Features

### 1. 視覺化輔助 (Visual Aids)
- **衝突警示燈**：在週曆視圖中，若時段重疊，背景轉為紅色並加上驚嘆號。
- **空網格點選 (Ghost Slot)**：點擊週曆上的空白處，直接彈出「新增規則」視窗並預填該時間與星期。
- **拖拽與縮放 (Drag & Drop)**：在週曆上拖拽課程塊可直接變更星期/時間；拉動底部邊緣可調整課程長度。

### 2. 智慧化工具 (Smart Tools)
- **閒置時段搜尋 (Empty Slot Finder)**：針對特定教室，一鍵反白「無課時段」，協助管理員安排新學期的補課或空檔填補。
- **批量編輯模式**：進入編輯模式後，可勾選多個課程區塊，統一修改「任課老師」或「地點」。
- **學期對照檢視 (Split View)**：複製學期時，提供左/右對照視圖，一邊顯示「來源學期」，一邊預覽「複製後的標的學期」，避免漏看課程。

### 3. 操作效率 (Efficiency)
- **快速篩選器預設**：提供「僅顯示有衝突的老師」、「僅顯示空堂較多的教室」等智慧篩選條件。
- **響應式視圖切換**：在手機查看佔用表時，自動從「網格視圖」切換為「垂直列表視圖」，解決行動版螢幕過窄的問題。

## Proposed Changes

### 1. Database & Backend [Component]

#### [NEW] [center_term.go](file:///d:/project/TimeLedger/app/models/center_term.go)
- **Schema**:
    - `ID` (uint): 主鍵
    - `CenterID` (uint): 中心 ID (外鍵)
    - `Name` (string): 期間名稱 (例如 "2026-Q1")
    - `StartDate` (time.Time): 開始日期
    - `EndDate` (time.Time): 結束日期
    - `CreatedAt`, `UpdatedAt`, `DeletedAt`

#### [NEW] [AdminTermController](file:///d:/project/TimeLedger/app/controllers/admin_term_controller.go)
- **CRUD**: 標準 GORM API 實作。
- **Aggregation Endpoint**: `GET /admin/occupancy/rules`
    - **邏輯**:
        1. 根據 `teacher_id` 或 `room_id` 查詢。
        2. 篩選 `ScheduleRule`，其中 `(rule.StartDate <= term.EndDate) AND (rule.EndDate >= term.StartDate)`。
        3. 按 `DayOfWeek` 分組返回，前端負責渲染到週曆格點中。
- **Batch Copy Endpoint**: `POST /admin/terms/copy-rules`
    - **邏輯**:
        1. 取得來源 `rule_ids` 列表。
        2. 針對每條規則進行 `Deep Copy`（排除 ID）。
        3. **日期重對齊 (Date Re-alignment)**：將新規則的 `StartDate` 設為標的學期的 `StartDate`，`EndDate` 設為標的學期的 `EndDate`。
        4. 保存新規則，並觸發 `ScheduleExpansionService` 重新展開課程。

---

### 2. Resource Management UI [Component]

#### [MODIFY] [resources.vue](file:///d:/project/TimeLedger/frontend/pages/admin/resources.vue)
- Add "學期期間 (Terms)" tab.
- Integrate `TermsTab` component.

#### [NEW] [TermsTab.vue](file:///d:/project/TimeLedger/frontend/components/Admin/TermsTab.vue)
- CRUD interface for terms.

---

### 3. Occupancy Visualization [Component]

#### [NEW] [resource-occupancy.vue](file:///d:/project/TimeLedger/frontend/pages/admin/resource-occupancy.vue)
- Filter Bar: Term selection, Teacher/Room search.
- **Rule Weekly Grid**: A specialized grid that displays rules on a 7-day layout.
- **Conflict Detection**: Highlight overlapping rules in the same slot.
- **Batch Copy Wizard**:
    - Select source term -> Filter rules (all/by course) -> Select target term -> Confirm copy.

---

### 4. Integration [Component]

#### [MODIFY] [ScheduleRuleForm.vue](file:///d:/project/TimeLedger/frontend/components/Scheduling/ScheduleRuleForm.vue)
- Add "快速填寫學期" dropdown to auto-fill Start/End dates based on defined Terms.

## Verification Plan

### 1. 自動化測試 (Automated Tests)
- **Backend Unit Tests**:
    - 在 `center_term_test.go` 中測試 CRUD。
    - 在 `term_copy_service_test.go` 中測試複製邏輯：驗證複製後的日期是否正確對齊標的學期。
    - 驗證重複複製時的錯誤處理（或覆蓋邏輯）。

### 2. 手動驗證流程 (Manual Verification)
#### [階段 A] 學期管理
1. 建立兩個學期：「2026 第一學期」(1/1~3/31) 與「2026 第二學期」(4/1~6/30)。
2. 確認列表顯示正確，且起訖日期邏輯正常。

#### [階段 B] 資源週曆 (Occupancy View)
1. 為老師 A 在第一學期排一堂「週一 09:00 - 10:00」的課。
2. 開啟 `resource-occupancy` 頁面，選擇第一學期並搜尋老師 A。
3. 驗證該課程正確顯示在週一網格中，且不帶具體日期標籤。
4. 加入一堂「週一 09:30」的教室重疊課程，驗證系統是否顯示「衝突警告」。

#### [階段 C] 批量複製測試
1. 在週曆視圖使用「複製工具」。
2. 來源：第一學期；標的：第二學期。
3. 選取剛才建立的課程，執行複製。
4. 切換到第二學期視圖，驗證課程已成功複製，且日期已自動調整為 4/1~6/30。
5. 檢查 `schedules` 列表，確認新產生的規則數量正確。
