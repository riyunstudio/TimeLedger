# 排課系統修復總結

本文件記錄排課系統的修復內容，包含問題描述、修復方案、修改檔案及待測試項目。

---

## 第一階段修復：週曆顯示與例外處理

### 修復的問題

本階段修復了四個核心問題，這些問題嚴重影響了排課功能的正確性和使用者體驗。

第一個問題是「新增課程後週曆不更新」。當管理員透過 ScheduleRuleModal 建立新课程规则时，模組僅呼叫 API 儲存資料，但未通知父元件重新整理畫面。這導致使用者必須手動重新整理頁面才能看到新增的課程，造成使用上的極大不便。修復方案是改用 Vue 的 `emit('created')` 機制，在成功建立規則後主動通知父元件刷新資料。

第二個問題是「同一時段多個課程覆蓋」。前端原本使用 `hour-day` 作為資料結構的 Key（例如 "14-2026-01-20"），當多個課程安排在同一時段時，後建立的課程會直接覆蓋先建立的課程，導致課程遺失。修復方案是將資料結構從 `Map`（物件）改為 `Array`（陣列），保留所有課程記錄。

第三個問題是「例外處理未實作」。前端原本沒有處理 `CANCEL`（取消）和 `RESCHEDULE`（改期）兩種例外類型，當課程被取消或改期時，週曆無法正確顯示異動資訊。修復方案是改用後端的 `/admin/expand-rules` API，由後端統一處理例外邏輯並返回完整的展开展現資料。

第四個問題是「SINGLE 模式無效」。當管理員使用 SINGLE 模式修改單一課程時間時，系統沒有正確處理例外單的建立和規則的分割，導致修改後的課程無法正確顯示。修復方案是實作完整的 SINGLE 模式邏輯：建立 CANCEL 例外單（取消原場次）、建立新規則（只針對該日期）、建立 RESCHEDULE 例外單（關聯新舊規則）。

### 修改的檔案

前端部分共修改了三個檔案。`frontend/components/ScheduleRuleModal.vue` 是課程規則建立與編輯的表單元件，這次修改改用 `emit('created')` 機制，在成功建立規則後通知父元件刷新資料。`frontend/components/ScheduleGrid.vue` 是週曆顯示元件，這次修改改用後端的 `/admin/expand-rules` API 取得課程資料，支援例外處理邏輯。

後端部分共修改了三個檔案。`app/controllers/scheduling.go` 是排課相關的 API 控制器，這次修正了 ExpandRules API 的日期解析邏輯，並修正了 SINGLE 模式的處理邏輯。`app/services/scheduling_interface.go` 是排課服務的介面定義，這次在 `ExpandedSchedule` 結構中新增了關聯欄位，用於儲存例外規則與原始規則的關聯資訊。`app/services/scheduling_expansion.go` 是排課展開服務的實作，這次在 ExpandRules 方法中新增了填充關聯資料的邏輯。

### 核心改動

本次修復有三項核心改動，這些改動奠定了後續功能的基礎。

第一項核心改動是前端資料結構從 Map 改為 Array。修復前的程式碼使用 `Record<string, any>` 作為資料結構，這種設計導致同一時段的多個課程會相互覆蓋。修復後改用 `any[]` 陣列結構，所有課程記錄都會被保留，週曆可以正確顯示同一時段的多個課程。

第二項核心改動是使用後端 ExpandRules API。修復前前端直接操作本地資料，無法正確處理例外邏輯。修復後統一由後端處理，呼叫 `/admin/expand-rules` API 並傳入規則 ID 列表、起始日期和結束日期，後端會自動展開循環規則並處理所有例外狀況。

第三項核心改動是 SINGLE 模式正確處理。當管理員選擇 SINGLE 模式修改單一課程時，系統會自動建立三筆記錄：CANCEL 例外單（標記原場次被取消）、新規則（只針對該日期的規則）、RESCHEDULE 例外單（關聯新舊規則）。這種設計確保了修改的歷史軌跡完整保留，同時週曆顯示正確。

### 待測試項目

以下項目需要在正式環境中進行測試驗證，確保修復後的功能正常運作。

第一個測試項目是「新增課程後週曆是否立即刷新」。操作步驟為：打開排課頁面、點擊新增課程、按鈕填寫課程資訊並儲存。預期結果為：儲存成功後週曆立即顯示新课程，無需手動重新整理頁面。

第二個測試項目是「同一時段顯示多個不同課程」。操作步驟為：選擇同一個時段（例如週一 14:00）、建立第一門課程、建立第二門課程（同一時段）。預期結果為：週曆正確顯示兩門課程，滑鼠懸停可看到各自的課程資訊。

第三個測試項目是「SINGLE 模式調整時間後正確顯示」。操作步驟為：選擇一門已存在的循環課程、進入編輯模式、選擇 SINGLE 模式、修改時間並儲存。預期結果為：原規則的該日期場次被標記為取消，週曆顯示新的時間場次，其他日期的場次不受影響。

第四個測試項目是「例外申請核准後週曆正確反應」。操作步驟為：老師提交請假申請、管理員核准申請。預期結果為：核准後週曆立即更新，顯示被取消的課程或改期後的課程。

---

## 第二階段修復：例外審核與通知

### 修復的問題

本階段修復了三個與例外審核和管理員通知相關的問題。

第一個問題是「看不到過往審核過的資料」。管理員只能看到目前待審核的例外申請，無法查看歷史記錄（包括已核准、已拒絕、已撤回的申請）。這導致管理員無法追蹤歷史審核情況，也不便於統計和對帳。修復方案是新增 `/api/v1/admin/exceptions/all` API 端點，支援取得所有狀態的例外申請，並在前方增加「已撤回」篩選按鈕。

第二個問題是「應該可以指定查詢審核資料」。管理員需要能夠依日期範圍篩選例外申請，以便針對特定期間進行審核或統計。修復方案是在審核頁面新增日期範圍篩選器，支援起始日期和結束日期的輸入，讓管理員可以快速找到特定期間的申請。

第三個問題是「老師申請時中心管理要有通知」。當老師提交例外申請時，中心管理員沒有收到即時通知，導致審核延遲。修復方案有兩個部分：後端在 CreateException 服務中新增發送通知給管理員的邏輯；前端在管理員 Header 新增通知鈴鐺按鈕，顯示未讀通知數量紅點，點擊可直接跳轉到審核頁面。

### 修改的檔案

後端部分共修改了七個檔案。`app/services/scheduling_interface.go` 新增了 `GetAllExceptions` 介面定義。`app/services/scheduling_expansion.go` 實作了 `GetAllExceptions` 方法和例外通知發送邏輯。`app/controllers/scheduling.go` 新增了取得所有例外申請的 API 端點。`app/servers/route.go` 註冊了 `/admin/exceptions/all` 路由。`app/repositories/notification.go` 修復了管理員通知查詢邏輯，改用 `center_id` 過濾而非 `user_id`。`app/controllers/notification.go` 修復了通知未讀數量的計算邏輯。`app/services/notification_interface.go` 和 `app/services/notification.go` 修改了 `SendAdminNotification` 方法，增加 `notificationType` 參數以支援不同類型的通知。

前端部分共修改了四個檔案。`frontend/pages/admin/approval.vue` 新增了篩選功能和日期範圍選擇器。`frontend/components/AdminHeader.vue` 新增了通知鈴鐺按鈕。`frontend/layouts/admin.vue` 整合了 NotificationDropdown 通知下拉元件。`frontend/components/NotificationDropdown.vue` 實作了通知列表顯示和點擊跳轉功能。

### 核心改動

本階段的核心改動聚焦於通知系統和管理員審核體驗的提升。

第一項核心改動是例外通知的即時推送。當老師提交例外申請時，系統會自動建立一筆記錄並發送通知給該中心的所有管理員。通知類型使用 `notificationType` 參數區分，方便前端針對不同類型進行不同的處理邏輯。

第二項核心改動是通知查詢邏輯的修正。原有的通知查詢使用 `user_id` 作為過濾條件，但管理員通知應該以 `center_id` 為依據。修復後，管理員登入時會看到所屬中心的所有通知，確保不會遺漏任何重要資訊。

第三項核心改動是審核頁面的功能增強。新增的篩選功能讓管理員可以依狀態（待審核、已核准、已拒絕、已撤回）和日期範圍快速找到目標申請，提升審核效率。

---

## 檔案異動清單

以下表格列出所有修改的檔案及其修改內容，便於追蹤變更歷史。

| 檔案 | 修改內容 |
|:---|:---|
| `app/services/scheduling_interface.go` | 新增 `GetAllExceptions` 介面 |
| `app/services/scheduling_expansion.go` | 實作 `GetAllExceptions`、例外通知發送 |
| `app/controllers/scheduling.go` | 新增 `GetAllExceptions` API、修正日期解析 |
| `app/servers/route.go` | 註冊 `/admin/exceptions/all` 路由 |
| `app/repositories/notification.go` | 管理員通知查詢邏輯修復 |
| `app/controllers/notification.go` | 修復未讀數量計算 |
| `app/services/notification_interface.go` | `SendAdminNotification` 增加類型參數 |
| `app/services/notification.go` | 實作通知類型參數 |
| `frontend/pages/admin/approval.vue` | 篩選功能、日期範圍選擇器 |
| `frontend/components/AdminHeader.vue` | 通知鈴鐺按鈕 |
| `frontend/layouts/admin.vue` | NotificationDropdown 整合 |
| `frontend/components/NotificationDropdown.vue` | 點擊通知跳轉功能 |
| `frontend/components/ScheduleRuleModal.vue` | 改用 `emit('created')` 通知父元件 |
| `frontend/components/ScheduleGrid.vue` | 改用 `/admin/expand-rules` API |

---

## 測試驗證清單

以下為各功能的測試驗證項目，建議依序執行測試以確保修復完整性。

### 週曆顯示功能測試

週曆顯示功能的測試項目包含：新增課程後週曆即時更新（無需手動重新整理）、同一時段顯示多門課程（不覆蓋）、例外課程正確顯示（CANCEL 顯示刪除線、RESCHEDULE 顯示新時間）、循環規則正確展開（每週重複的課程正確顯示）。

### 例外申請流程測試

例外申請流程的測試項目包含：老師提交請假申請（選擇日期、填寫原因）、管理員收到通知（鈴鐺顯示紅點）、管理員審核頁面看到新申請、核准後週曆正確更新。

### SINGLE 模式測試

SINGLE 模式的測試項目包含：選擇循環課程的特定場次、修改該場次的時間和日期、檢查原場次是否標記為取消、檢查新場次是否正確顯示、檢查其他場次是否不受影響。

### 管理員通知測試

管理員通知的測試項目包含：通知鈴鐺顯示未讀數量、點擊鈴鐺展開通知列表、點擊通知跳轉到正確頁面、已讀通知不再顯示紅點、例外申請通知包含正確資訊。

---

## 後續優化建議

基於本次修復的經驗，以下為未來可以考慮的優化方向。

在週曆效能方面，當課程數量較多時，可以考慮實作虛擬捲動（Virtual Scrolling）優化渲染效能，或者增加課程衝突的視覺提示，讓管理者一目了然。

在通知體驗方面，可以考慮支援通知的批次處理（例如一次核准或拒絕多筆申請）、新增通知設定功能讓管理員選擇接收哪些類型的通知。

在例外審核方面，可以考慮新增例外申請的詳細資訊頁面（目前只有列表）、增加核准時的審核意見記錄功能、統計各老師的例外申請次數並提供圖表。

---

*文件更新日期：2026年1月26日*
*維護者：TimeLedger 開發團隊*
