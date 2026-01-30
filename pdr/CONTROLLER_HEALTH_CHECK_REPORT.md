# TimeLedger 控制器健康檢查報告 (Project-wide)

本報告根據「Teacher 模組優化模式」對全專案控制器進行深度檢測。

## 1. 核心評價標準
*   **Thin Controller**: 控制器是否僅負責網頁請求與回應？
*   **Service Extraction**: 業務邏輯是否已抽離至 Service 層？
*   **ContextHelper**: 是否使用統一工具提取 User/Center ID？
*   **Domain Focus**: 控制器是否專注於單一領域？

---

## 2. 詳細清單與優化方向

| 控制器檔案 | 健康狀態 | 重點問題 | 優化方向 (根據 Teacher 模式) |
| :--- | :---: | :--- | :--- |
| `admin_resource.go` | 🔴 極差 (God) | 44KB 巨無霸，混合 5 個領域 | **拆分**：分為 Center, Room, Course, Holiday 控制器與 Service。 |
| `timetable_template.go` | 🟠 待改進 | 19KB，含複雜校驗與套用邏輯 | **服務化**：建立 `TimetableTemplateService`，搬移核心算法。 |
| `export.go` | 🟠 待改進 | 19KB，含大量日期、格式解析邏輯 | **清理**：將解析邏輯與 DB 查詢移至 `ExportService`。 |
| `line_bot.go` | 🟠 待改進 | 事件處理遍布控制器 | **解耦**：將 Webhook 事件分發後的邏輯移至 `LineBotService`。 |
| `offering.go` | 🔴 舊模式 | 完全無 Service，直接使用 Repo | **重構**：建立 `OfferingService`，繼承泛型 Repo。 |
| `admin_user.go` | 🟡 尚可 | 手動提取身分邏輯重複 | **標準化**：全面導入 `ContextHelper` 與統一錯誤碼。 |
| `smart_matching.go` | 🟡 尚可 | Context 提取邏輯複雜且非統一 | **標準化**：導入 `ContextHelper` 簡化代碼。 |
| `notification.go` | 🟡 尚可 | 部分內部方法仍直呼 Repo | **清理**：統一所有 DB 操作進入 `NotificationService`。 |
| `auth.go` | 🟢 優良 | 體積小，邏輯明確 | **微調**：檢查是否可進一步簡化 Response 結構。 |
| `user.go` | 🟢 優良 | 體積小 | 無需重大變動。 |

---

## 3. 優先順序建議 (Priority)

### 第一優先 (P0) - 結構性崩壞
*   **`admin_resource.go`**: 這是系統擴張的最大阻礙，應最先進行領域拆分。
*   **`offering.go`**: 雖然體積一般，但模式完全過時，是引入 `GenericRepository` 的最佳實驗場。

### 第二優先 (P1) - 邏輯混亂
*   **`timetable_template.go`**: 排課模組的最後一片拼圖，需完成 Service 化。
*   **`export.go`**: 提高 Export 模組的可測試性。

### 第三優先 (P2) - 代碼清理 (Standardization)
*   全面將 `admin_user`, `smart_matching`, `notification` 等模組的 ID 提取邏輯替換為 `ContextHelper`。

---
**總結**：目前專案已完成 Teacher 與 Scheduling 兩大核心重構，約占系統複雜度的 40%。若能完成 `AdminResource` 的拆解，系統健康度將達到 80% 以上。
