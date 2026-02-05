# LINE Webhook 2.0 詳細實作藍圖

## 🧠 邏輯實作指南 (Logic Deep-Dive)

### 1. 跨中心課表聚合演算法
當 Webhook 觸發「今日行程」時，後端應執行以下步驟：

```go
func (s *LineBotService) GetAggregatedAgenda(lineUID string) ([]AgendaItem, error) {
    // 1. 識別使用者
    identity := s.IdentifyUser(lineUID)
    
    var agenda []AgendaItem
    
    // 2. 抓取中心排課
    for _, membership := range identity.Memberships {
        rules := s.expansionSvc.Expand(membership.CenterID, today)
        for _, r := range rules {
            agenda = append(agenda, AgendaItem{
                Time: r.StartTime,
                Title: r.OfferingName,
                Center: membership.CenterName,
                Type: "CENTER",
            })
        }
    }
    
    // 3. 抓取個人行程
    if identity.TeacherProfile != nil {
        personal := s.personalEventSvc.GetTodayOccurrences(identity.TeacherProfile.ID, today)
        for _, p := range personal {
            agenda = append(agenda, AgendaItem{
                Time: p.StartAt.Format("15:04"),
                Title: p.Title,
                Center: "個人",
                Type: "PERSONAL",
            })
        }
    }
    
    // 4. 關鍵：排序
    sort.Slice(agenda, func(i, j int) {
        return agenda[i].Time < agenda[j].Time
    })
    
    return agenda, nil
}
```

### 2. 重複性行程展開 (Personal Events)
需在 `app/services/personal_event.go` 處理 RRule：

```go
func (s *PersonalEventService) GetTodayOccurrences(teacherID uint, date time.Time) []models.PersonalEvent {
    // 撈取該老師的所有行程 (含 RecurrenceRule != nil)
    // 遍歷行程，若 (IsSameDay) OR (RecurrenceRule 匹配當日週幾)
    // 則回傳該行程實例
}
```

### 3. 前端即時預覽組件 (Vue)
在 `LineFlexPreview.vue` 中監聽 Props：

```javascript
const previewData = computed(() => {
  return {
    type: "bubble",
    body: {
      type: "box",
      contents: [
        { type: "text", text: props.title, weight: "bold", size: "xl" },
        { type: "text", text: props.content, wrap: true }
      ]
    }
  }
})
```

## 🛡️ 安全防護與回歸測試清單 (Safety & Regression Checklist)

### 1. 防斷點檢查 (Breakpoint Prevention)
- [ ] **身分識別隔離**：確保 `GetCombinedIdentity` 失敗時，會回傳一個預設的 `GUEST` 結構，而不是回傳 `nil` 或噴錯，避免後續邏輯崩潰。
- [ ] **並行安全**：若使用 Goroutine 查詢身份，須確保使用 `sync.WaitGroup` 並有超時控制，防止 Webhook 因資料庫回應過慢而超時。

### 2. 回歸測試清單 (Regression List)
- [ ] **原有指令驗證**：測試輸入「綁定 [驗證碼]」，確保原有的綁定流程依然正常運作。
- [ ] **匿名訊息驗證**：測試未綁定使用者發送訊息，系統應正常回覆「預設導引」，而非報錯。
- [ ] **現有課程影響**：在前端查看原本的課表頁面，確認後端的 `PersonalEvent` 修改並未影響到管理員端的排課網格顯示。

## 🛠️ 下一步開發建議 (Developer Commands)
對 Cursor 下達：
1. `Update app/services/line_bot.go to implement CombinedIdentity and aggregation logic. DON'T modify any existing command handlers (like binding logic) - only extend them.`
2. `Create AdminNotificationController with a broadcast endpoint. Ensure it uses standard admin middleware for safety.`
3. `New page at frontend/pages/admin/broadcast.vue with flex message preview.`
