# LINE Webhook 2.0 升級開發任務清單

## 🟢 Phase 1: 後端基礎與身份識別
- [ ] **[Models]** 定義 `CombinedIdentity` 與 `AgendaItem` 數據結構。
- [ ] **[Service]** 在 `LineBotService` 實作 `IdentifyUser`：
    - [ ] 同步/併行查詢 Admin 與 Teacher。
    - [ ] 預載入 CenterMembership 與 Center 名稱。
- [ ] **[Testing]** 撰寫 `IdentifyUser` 的單元測試。

## 📅 Phase 2: 行程聚合與排序
- [ ] **[Service]** 在 `PersonalEventService` 實作 `GetTodayOccurrences`：
    - [ ] 處理重複性規則（Weekly, Monthly）展開。
- [ ] **[Service]** 在 `LineBotService` 實作 `GetAggregatedAgenda`：
    - [ ] 獲取所有相關中心的課表。
    - [ ] 獲取個人行程。
    - [ ] 實作時間排序 (Sort by StartTime)。
- [ ] **[Template]** 更新 `LineBotTemplateService`：
    - [ ] 支援循環渲染多筆行程。
    - [ ] 樣式區分：中心課(藍) vs 個人(紫)。

## 📢 Phase 3: 公告廣播系統
- [ ] **[Controller]** 建立 `AdminNotificationController`：
    - [ ] 實作 `Broadcast` 端點。
    - [ ] 加入 `CenterID` 權限校驗。
- [ ] **[Frontend]** 建立 `LineFlexPreview.vue` 組件：
    - [ ] 模擬手機外框與 Flex Message。
- [ ] **[Frontend]** 建立 `Broadcast.vue` 頁面：
    - [ ] 訊息輸入、即時預覽、發送按鈕。

## 🧭 Phase 4: UI 優化與連結
- [ ] **[Layout]** 更新 `admin.vue`：
    - [ ] Logo 點擊跳轉 `/admin/dashboard`。
    - [ ] 側邊欄加入「一鍵公告」入口。
- [ ] **[Template]** 在 Flex Message 底部加入「連結網站」按鈕。

## 🛡️ Phase 5: 觀察與驗證
- [ ] **[Log]** 加入詳細的 Webhook 事件日誌。
- [ ] **[Log]** 加入廣播發送紀錄日誌。
- [ ] **[Safety]** 加入 API 限流防護。
