# Cursor 逐步執行指南 (Micro-Task Instruction Guide)

為了確保穩定性，請按照以下順序，每次複製 **一個步驟** 的指令給 Cursor 執行。完成一個步驟並確認編譯通過後，再進行下一個。

---

### 第一階段：數據與後端邏輯層 (Backend & Logic)

#### **步驟 1：定義基礎資料結構**
- **檔案**：`app/services/line_bot.go`
- **指令**：
> 「請在 `app/services/line_bot.go` 中定義 `CombinedIdentity` 與 `AgendaItem` 結構。
> `CombinedIdentity` 應包含 AdminProfiles, TeacherProfile, Memberships 與 PrimaryRole。
> `AgendaItem` 應包含時間、標題、來源名稱（中心名或個人）與來源類型。
> **注意：此步驟僅新增定義，不要修改任何現有函數。**」

#### **步驟 2：實現身分識別邏輯 (Identity Fetching)**
- **檔案**：`app/services/line_bot.go`
- **指令**：
> 「在 `LineBotService` 中實作 `GetCombinedIdentity(lineUserID string)` 方法。
> 它應該查詢管理員表與老師表。如果是老師，請預載入其所有參與中心 (CenterMembership) 的名稱。
> 如果找不到任何身分，回傳一個標記為訪客 (GUEST) 的預設結構。請確保邏輯包含 Error Handling。」

#### **步驟 3：撰寫身分識別測試**
- **檔案**：`app/services/line_bot_test.go`
- **指令**：
> 「請為剛才新增的 `GetCombinedIdentity` 撰寫單元測試。測試案例需包含：1. 僅管理員 2. 僅老師 3. 訪客（未綁定）。」

#### **步驟 4：個人重複行程展開邏輯**
- **檔案**：`app/services/personal_event.go`
- **指令**：
> 「在 `PersonalEventService` 中新增 `GetTodayOccurrences(teacherID uint, targetDate time.Time)` 方法。
> 此方法需撈取該老師的所有行程，並根據 `RecurrenceRule` (RRule) 判斷該行程是否發生在 `targetDate` 當天。若是，則回傳該實例。」

#### **步驟 5：實作行程聚合引擎 (The Aggregator)**
- **檔案**：`app/services/line_bot.go`
- **指令**：
> 「在 `LineBotService` 中實作 `GetAggregatedAgenda(lineUserID string)` 核心方法。
> 1. 先呼叫 `GetCombinedIdentity`。
> 2. 循環各中心調用排課擴展服務獲取當日課表。
> 3. 調用 `PersonalEventService.GetTodayOccurrences` 獲取個人事。
> 4. 將兩者轉換為 `AgendaItem` 並按 `StartTime` 排序。
> **請編寫此聚合邏輯的單元測試，驗證排序是否正確。**」

---

### 第二階段：視覺範本與整合 (Visuals & Integration)

#### **步驟 6：升級 LINE Flex Message 範本**
- **檔案**：`app/services/line_bot_template.go`
- **指令**：
> 「修改 `GenerateAgendaFlex` 範本。
> 1. 支援顯示多筆行程列表。
> 2. 顏色區分：中心課程使用藍色系，個人行程使用紫色系。
> 3. 在訊息底部新增一個 [進入系統首頁] 的按鈕連結。」

#### **步驟 7：無損接入 Webhook 處理器**
- **檔案**：`app/controllers/line_bot.go`
- **指令**：
> 「在 `LineBotController` 的 `handleMessageEvent` 中，當使用者輸入『課表』或相關指令時，改用新開發的 `GetAggregatedAgenda` 來獲取資料。
> **絕對禁止異動到原有的『綁定指令』處理邏輯，確保舊功能（改A不壞B）完全正常。**」

---

### 第三階段：管理端廣播功能 (Admin Broadcast UI)

#### **步驟 8：建立廣播 API (Backend)**
- **檔案**：`app/controllers/admin_notification.go` (NEW), `app/servers/route.go`
- **指令**：
> 「1. 新增 `AdminNotificationController` 並實作 `Broadcast` 方法。要求權限校驗，僅限管理員對其所屬中心的老師發送。
> 2. 在 `route.go` 的 admin 群組註冊此 POST 路由。」

#### **步驟 9：開發 LINE 訊息即時預覽組件 (Frontend)**
- **檔案**：`frontend/components/Notification/LineFlexPreview.vue` (NEW)
- **指令**：
> 「新增 `LineFlexPreview.vue` 組件。它需要模擬 LINE 氣泡訊息的樣式，並根據傳入的標題與內文 props，即時顯示預覽效果。請確保 CSS 樣式美觀且具備手機框質感。」

#### **步驟 10：開發廣播管理頁面 (Frontend)**
- **檔案**：`frontend/pages/admin/broadcast.vue` (NEW)
- **指令**：
> 「新增 `admin/broadcast.vue` 頁面。左側為公告內容輸入框，右側引入 `LineFlexPreview` 組件。點擊發送時應有二次確認視窗，並呼叫後端廣播 API。」

#### **步驟 11：更新導航列入口 (Navigation)**
- **檔案**：`frontend/layouts/admin.vue`
- **指令**：
> 「修改導航列：
> 1. 將 Logo 包裝成 `NuxtLink` 到首頁。
> 2. 在選單中新增『一鍵公告』入口。
> 3. 更新同步到行動版 (Mobile) 選單中。」

---

### 第四階段：維護與觀察 (Observability)

#### **步驟 12：日誌與安全加固**
- **指令**：
> 「請在 LINE Webhook 的入口處加入詳細的監控日誌，記錄 UID 與識別出的 PrimaryRole。
> 同時為廣播功能加入簡單的 Rate Limiting，防止管理員連點導致頻繁發送。」
