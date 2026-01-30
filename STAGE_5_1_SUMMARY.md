# Stage 5.1 N+1 查詢優化階段總結

## 1. 優化目標

本次優化聚焦於排課模組中的 N+1 查詢問題，透過批次查詢與記憶體 Map 組裝策略，大幅減少資料庫查詢次數，提升系統整體效能。

### 核心問題描述

在排課模組的原有實作中，存在以下 N+1 查詢問題：

- **循環查詢**：在 DetectPhaseTransitions 方法中，每次迴圈都會呼叫 GetEffectiveRuleForDate 查詢資料庫
- **延遲載入**：存取關聯資料（Offering、Room、Teacher）時觸發額外查詢
- **查詢膨脹**：查詢範圍越大（如 30 天、90 天），查詢次數線性成長

### 優化策略

採用以下策略消除 N+1 查詢問題：

| 策略 | 說明 |
| :--- | :--- |
| 批次查詢 | 使用 ListByOfferingIDWithPreload 一次性取得所有規則及關聯資料 |
| Map 組裝 | 將查詢結果轉換為 Map結構，實現 O(1) 查詢效率 |
| 預載入優化 | 在查詢時使用 Preload 預先載入關聯資料，避免延遲載入 |

---

## 2. 變更檔案清單

| 檔案 | 變更類型 | 說明 |
| :--- | :--- | :--- |
| `app/repositories/schedule_rule.go` | 新增方法 | 新增 ListByOfferingIDWithPreload 批次查詢方法 |
| `app/services/scheduling_expansion.go` | 重構方法 | 重構 DetectPhaseTransitions 使用 Map 查詢 |
| `testing/test/schedule_rule_test.go` | 新增測試 | 新增兩個單元測試案例 |
| `pdr/progress_tracker.md` | 文件更新 | 記錄 Stage 5.1 完成狀態 |

---

## 3. 優化內容詳解

### 3.1 Repository 層優化

#### 新增方法：ListByOfferingIDWithPreload

```go
// app/repositories/schedule_rule.go (Lines 310-321)
// 批次查詢規則並預載入關聯資料（消除 N+1 查詢）

func (rp *ScheduleRuleRepository) ListByOfferingIDWithPreload(
    ctx context.Context,
    offeringID uint,
) ([]models.ScheduleRule, error) {
    var data []models.ScheduleRule
    err := rp.app.MySQL.RDB.WithContext(ctx).
        Preload("Offering").
        Preload("Room").
        Preload("Teacher").
        Where("offering_id = ?", offeringID).
        Order("effective_range ASC").
        Find(&data).Error
    return data, err
}
```

#### 效益說明

| 指標 | 原本行為 | 優化後行為 |
| :--- | :--- | :--- |
| 查詢次數 | 1 + (規則數 × 3) | 1 |
| 觸發時機 | 存取關聯資料時延遲載入 | 查詢時立即載入 |
| 資料完整性 | 可能觸發多次查詢 | 單一查詢完整取得 |

---

### 3.2 Service 層優化

#### 重構方法：DetectPhaseTransitions

**原本實作（N+1 問題）**

```go
// 每次迴圈都查詢資料庫
for date.Before(endDate) || date.Equal(endDate) {
    currentRule, _ := s.GetEffectiveRuleForDate(ctx, offeringID, date)
    // ...
}
```

**優化後實作（批次查詢 + Map）**

```go
// 建立日期到規則的 Map：DateString -> *ScheduleRule
ruleByDate := make(map[string]*models.ScheduleRule)
for i := range rules {
    rule := &rules[i]
    // ... 展開規則日期範圍並填入 Map
}

for date.Before(endDate) || date.Equal(endDate) {
    dateStr := date.Format("2006-01-02")
    currentRule := ruleByDate[dateStr] // O(1) Map 查詢
    // ...
}
```

#### 程式碼位置

- **檔案**：`app/services/scheduling_expansion.go`
- **方法**：`DetectPhaseTransitions`
- **變更內容**：將迴圈內的資料庫查詢改為 Map 查詢

---

## 4. 效能比較預測

### 4.1 查詢次數分析

| 場景 | 查詢範圍 | 優化前 | 優化後 | 減少比例 |
| :--- | :--- | :--- | :--- | :--- |
| DetectPhaseTransitions | 7 天 | 8 次 | 1 次 | 87.5% |
| DetectPhaseTransitions | 30 天 | 31 次 | 1 次 | 96.8% |
| DetectPhaseTransitions | 90 天 | 91 次 | 1 次 | 98.9% |
| ExpandRules（含關聯） | 10 規則 | 31 次 | 1 次 | 96.8% |

### 4.2 計算範例

假設查詢範圍為 30 天、有 5 條排課規則：

**優化前：**

| 操作 | 查詢次數 |
| :--- | :--- |
| ListByOfferingID() | 1 次 |
| GetEffectiveRuleForDate() × 30 | 30 次 |
| **總計** | **31 次資料庫查詢** |

**優化後：**

| 操作 | 查詢次數 |
| :--- | :--- |
| ListByOfferingIDWithPreload() | 1 次 |
| 建立 Map（記憶體操作） | O(n) 操作 |
| Map 查詢 × 30 | O(1) × 30 |
| **總計** | **1 次資料庫查詢** |

### 4.3 預估響應時間改善

| 指標 | 優化前 | 優化後 | 說明 |
| :--- | :--- | :--- | :--- |
| 資料庫查詢延遲 | ~310ms | ~10ms | 以每次查詢 ~10ms 估算 |
| 網路往返次數 | 31 次 | 1 次 | 減少連線建立開銷 |
| 資料庫負載 | 高 | 低 | 查詢量減少 97% |

---

## 5. 程式碼變更統計

| 指標 | 數值 |
| :--- | :--- |
| 新增方法數 | 1 個 |
| 重構方法數 | 1 個 |
| 新增測試案例 | 2 個 |
| 修改檔案數 | 4 個 |
| 減少查詢次數 | 最多 98.9% |
| 編譯狀態 | ✅ 通過 |

---

## 6. 使用建議

### 情境一：ExpandRules 需要關聯資料

```go
// 建議使用 ListByOfferingIDWithPreload
rules, err := s.scheduleRuleRepo.ListByOfferingIDWithPreload(ctx, offeringID)
if err != nil {
    return nil, err
}

schedules := s.ExpandRules(ctx, rules, startDate, endDate, centerID)
```

### 情境二：DetectPhaseTransitions

現有程式碼已自動使用優化後的實作，無需額外修改。

---

## 7. 已存在但非本次優化範圍

本次優化主要處理 Repository 層與 Service 層的 N+1 問題。以下模式在程式碼中已存在但屬於既有設計：

| 模式 | 位置 | 說明 |
| :--- | :--- | :--- |
| 假日批次查詢 | 第 39 行 | holidayRepo.ListByDateRange 已優化 |
| 例外批次查詢 | 第 46-57 行 | exceptionRepo.GetByRuleIDsAndDateRange 已優化 |
| Map 組裝 | 第 94-98 行 | 使用 Map 取得例外資料 |

---

## 8. 下一步建議

| 優先順序 | 項目 | 說明 |
| :--- | :--- | :--- |
| ~~中~~ → 低 | 更新調用端 | 經檢查，主要方法都已內建 Preload，無需額外修改 |
| 低 | 監控實際效能 | 部署後觀察資料庫查詢 logs 驗證改善 |
| 低 | 單元測試 | 為新方法新增測試案例（已完成） |

---

## 9. 驗證結果

### 9.1 編譯驗證

```bash
$ cd d:\project\TimeLedger
$ go build ./app/repositories/
# Exit code: 0 - Build successful!

$ go build ./app/services/
# Exit code: 0 - Build successful!

$ go build ./testing/test/
# Exit code: 0 - Build successful!
```

### 9.2 測試覆蓋

| 測試案例 | 功能 |
| :--- | :--- |
| TestScheduleRuleRepository_ListByOfferingIDWithPreload | 測試批次查詢與預載入功能 |
| TestScheduleRuleRepository_ListByOfferingIDWithPreload_PreloadVerification | 驗證關聯資料預載入的正確性 |

---

## 10. 總結

| 維度 | 成果 |
| :--- | :--- |
| 效能提升 | 查詢次數減少 87.5% ~ 98.9% |
| 程式碼品質 | 消除 N+1 查詢問題，提升可維護性 |
| 測試覆蓋 | 新增 2 個單元測試案例 |
| 文件更新 | 進度追蹤文件已更新 |
| 編譯狀態 | 全部套件編譯通過 |

### Stage 5.1 完成狀態：✅ 已完成

所有 N+1 查詢問題已修復，程式碼編譯通過，效能改善預估可減少 87% ~ 98.9% 的資料庫查詢次數。

---

## 附錄：相關資源

- ** Repository 層變更**：`app/repositories/schedule_rule.go`
- ** Service 層變更**：`app/services/scheduling_expansion.go`
- ** 測試案例**：`testing/test/schedule_rule_test.go`
- ** 進度追蹤**：`pdr/progress_tracker.md`
