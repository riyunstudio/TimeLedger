# LINE Webhook 2.0 UX Upgrade - 終極開發藍圖 (Ultimate Technical Spec)

這份文件旨在提供給 **Cursor (AI Coding Assistant)** 進行自動化開發。它包含了精確的檔案修改路徑、邏輯演算法、單元測試規劃及預期結果。

---

## 🚀 專案模組 1：身份感知與分流 (Identity Awareness)

### 1.1 修改範圍與邏輯
*   **檔案**：`app/services/line_bot.go`
*   **功能描述**：識別 LINE 使用者的多重身份（管理員、各中心老師、訪客）。
*   **技術邏輯**：
    1.  建立 `CombinedIdentity` 結構：
        ```go
        type CombinedIdentity struct {
            AdminProfile *models.AdminUser
            TeacherProfile *models.Teacher
            Memberships []models.CenterMembership
            PrimaryRole string // "ADMIN" or "TEACHER" or "GUEST"
        }
        ```
    2.  實作 `IdentifyUser(lineUID string)`：
        - 併行呼叫 `AdminRepo.GetByLineID` 與 `TeacherRepo.GetByLineID`。
        - 若是教師，需同步加載 `CenterMembership` 與關聯的 `Center` 名稱。
*   **單元測試**：
    - 輸入綁定管理員的 UID，預期返回 `PrimaryRole: ADMIN`。
    - 輸入綁定老師的 UID，預期返回包含多個中心的 `Memberships`。

### 1.2 Cursor 實作指令
> 「請修改 `app/services/line_bot.go`，新增 `CombinedIdentity` 結構與 `IdentifyUser` 方法。需同時檢查 `AdminUsers` 與 `Teachers` 表。若是老師，請確保預載入其所有參與中心的名稱。請為此邏輯編寫對應的單元測試。」

---

## 📅 專案模組 2：跨來源行程聚合 (Agenda Aggregator)

### 2.1 修改範圍與邏輯
*   **檔案**：`app/services/line_bot.go`, `app/services/personal_event.go`, `app/services/line_bot_template.go`
*   **功能描述**：合併跨中心課表與個人私事，按時間排序。
*   **技術邏輯**：
    1.  **資料收集**：
        - Loop `Memberships` -> 調用 `ScheduleExpansionService.Expand(today)`。
        - 調用 `PersonalEventService.GetTodayOccurrences(teacherID)`。
    2.  **標準化物件**：統一轉為 `AgendaItem{Time, Title, CenterName, Type}`。
    3.  **排序演算法**：使用 `sort.Slice` 對 `Time` 進行升序排序。
    4.  **視覺範本**：在 `line_bot_template.go` 中，中心課用主題藍色，個人行程用對比紫色。
*   **預期結果**：老師在 LINE 收到一個 Flex Message，依序顯示「09:00 中心 A 數學」、「14:00 個人 牙醫門診」。

### 2.2 Cursor 實作指令
> 「實作 `LineBotService.GetTodayAgenda(lineUID)`。它必須合併來自多個中心的排課規則與老師的個人行程 `PersonalEvent`。請確保所有行程按起始時間正確排序。接著，請更新 `line_bot_template.go` 以支持這種動態列表的渲染。」

---

## 📢 專案模組 3：一鍵廣播系統 (Admin Broadcast)

### 3.1 修改範圍與邏輯
*   **後端檔案**：`app/controllers/admin_notification.go` (NEW), `app/servers/route.go`
*   **前端檔案**：`frontend/pages/admin/broadcast.vue` (NEW), `frontend/components/Notification/LineFlexPreview.vue` (NEW)
*   **功能描述**：管理員後台輸入文字，即時預覽 LINE 效果並一鍵廣播。
*   **技術邏輯**：
    - **前端預覽**：在 Vue 中實作一個模擬手機外殼的組件，根據輸入的 `title` 與 `body` 即時生成 Flex Message 預覽圖。
    - **後端廣播**：API 權限校驗 -> 撈取該中心所有綁定老師的 UID -> 呼叫 `LineBotSvc.Multicast`。
    - **防呆**：發送前需彈出二次確認視窗。
*   **單元測試**：
    - 準備測試資料：中心 A 有 3 位老師，2 位有 LINE 綁定。
    - 呼叫廣播 API，預期 Multicast 被呼叫 2 次。

### 3.2 Cursor 實作指令
> 「1. 在後端建立 `AdminNotificationController.Broadcast` API，僅允許管理員發送給該中心成員。
> 2. 在前端建立 `admin/broadcast.vue` 頁面，左側為輸入框，右側為 `LineFlexPreview` 組件。
> 3. `LineFlexPreview` 必須能根據輸入內容，即時渲染出模擬的 LINE 氣泡訊息樣式。」

---

## 🧭 專案模組 4：導航與連結閉環 (Navigation & Home Link)

### 4.1 修改範圍與邏輯
*   **網頁端**：`frontend/layouts/admin.vue`
*   **LINE 端**：`app/services/line_bot_template.go`
*   **功能描述**：Logo 點擊回首頁，LINE 訊息按鈕連回網站。
*   **技術邏輯**：
    - 將側邊欄 Logo 用 `NuxtLink` 包裝，目標為 `/admin/dashboard`。
    - 在所有今日摘要的 Flex Message 底部增加一個 `UriAction` 按鈕，標籤為「進入系統」，連結為前端 Dashboard。
*   **預期結果**：Logo 點擊必跳轉；LINE 訊息底部必有進入系統之按鈕。

### 4.2 Cursor 實作指令
> 「請將 `frontend/layouts/admin.vue` 中的 Logo 修改為可點擊跳轉至 Dashboard。同時，在 `line_bot_template.go` 生成的所有摘要訊息底部，增加一個導向 Web 版首頁的連結按鈕。」

---

## 🛡️ 驗證與上線 (Testing & Verification)

### ✅ 測試鏈條與回歸測試 (Regression Testing)
1.  **現有功能檢查**：在實作後，必須先執行既存的 `LineBotController` 測試，確保基本的「綁定指令」與「文字回覆」功能未受影響。
2.  **身分模擬**：手動在 DB 建立一個具有雙重身份的測試帳號。
3.  **廣播攔截**：檢查日誌 `global.Log.Info` 是否正確產出廣播紀錄。
4.  **邊界測試**：若當天沒有行程，預期回覆「今日尚無規劃」之友善訊息。

### 🛡️ 邏輯隔離原則 (Safety Guards)
*   **非破壞性修改**：新增功能應以「擴充」為主。例如，修改 `handleMessageEvent` 時，應保留原有的 `default` 處理邏輯，確保未定義的指令仍能正常回傳預設訊息。
*   **例外安全**：身分識別邏輯 (`IdentifyUser`) 應包含 `defer recover` 或強大的 `Error Handling`，即使 SQL 報錯，也要能回傳 `GUEST` 身份，而非導致 Webhook 崩潰。
*   **版本並行**：若更動幅度大，建議先建立 `v2` 版本的 Handler，待測試無誤後再進行切換。

### 📋 最終檢查表 (Definition of Done)
- [ ] 老師能看到跨中心與個人行程的合併排序列表。
- [ ] 管理員有專屬概況畫面且能發送廣播（含預覽）。
- [ ] 導航列 Logo 與 LINE 按鈕均能正確導引回首頁。
- [ ] 系統詳細記錄 Webhook 入口日誌。
