# TimeLedger 一致性與差距分析

## 1. 實現與文檔不一致之處

### 1.1 API 實現 vs 現有文檔

| 項目 | 發現 |
|:---|:---|
| **Route 數量** | 實際有 150+ 端點，部分現有 API.md 未涵蓋 |
| **Admin Teacher List** | 存在 `/api/v1/teachers` (GET) 和 `/api/v1/admin/teachers` (GET) 兩個端點，功能重疊需確認 |
| **邀請連結管理** | 包含一般邀請連結和一般邀請連結，功能類似但分離 |
| **課表範本** | 功能完整，但 API 文件未詳細描述 |

### 1.2 代碼 vs 設計

| 項目 | 發現 |
|:---|:---|
| **ContextHelper** | 大量使用，但部分控制器仍有重複程式碼 |
| **GenericRepository** | 廣泛使用，但部分自定義 Repository 未完全遵循模式 |
| **Service Layer** | 大部分已實作 Triple Return，但錯誤碼映射偶有不一致 |

---

## 2. 潛在驗證缺失

### 2.1 輸入驗證

| 區域 | 項目 | 狀態 |
|:---|:---|:---|
| ScheduleRule | StartTime/EndTime 格式驗證 | ✅ 已有 |
| ScheduleRule | Weekday 範圍 (0-6) 驗證 | ✅ 已有 |
| ScheduleRule | Effective Range 結束日期 > 開始日期 | ✅ 已有 |
| Exception | 原日期不能是過去 | ⚠️ 部分缺失 |
| Exception | 調課必須提供新時間 | ✅ 已有 |
| Teacher | Email 格式驗證 | ⚠️ 依賴 binding |
| Teacher | 標籤數量限制 (3-5個) | ⚠️ 未在所有 API 強制 |

### 2.2 權限驗證

| 區域 | 項目 | 狀態 |
|:---|:---|:---|
| ScheduleRule | 更新時確認 center_id 匹配 | ✅ 已有 |
| ScheduleRule | 刪除時確認 ownership | ✅ 已有 |
| Exception | 教師不能審核自己的例外 | ✅ 已有 |
| CenterInvitation | 只能接受自己收到的邀請 | ⚠️ 需加強 |

---

## 3. 潛在架構漂移

### 3.1 設計 vs 實現

| 設計意圖 | 實際實現 | 差距 |
|:---|:---|:---|
| Thin Controller | 大部分控制器已精簡 | ✅ 已實現 |
| RDB/WDB 分離 | Repository 層已實作 | ✅ 已實現 |
| Center Scope Check | 大部分查詢已有 | ⚠️ 部分漏網 |
| 錯誤碼前綴 | errInfos 使用 appID 前綴 | ✅ 已實現 |
| LINE 綁定不可解除 | 老師 LINE 綁定不可解除 | ✅ 已實現 |
| 管理員可解除 | 管理員 LINE 可解除 | ⚠️ 設計不明確 |

---

## 4. 模式違規

### 4.1 发现的模式违规

| 項目 | 說明 | 嚴重程度 |
|:---|:---|:---|
| **直接 DB 操作** | 部分 Service 層有直接操作 DB 的情況 | 低 |
| **業務邏輯在 Controller** | 個別端點仍有簡單業務邏輯 | 低 |
| **重複驗證** | Request 層和 Service 層有重複驗證 | 低 |
| **Hard Delete** | 特定資料表有硬刪除需求但未統一處理 | 低 |

---

## 5. 財務/決策邏輯風險區域

### 5.1 需要關注的區域

| 區域 | 風險 | 說明 |
|:---|:---|:---|
| **排課驗證** | 低 | 衝突檢查完整，但緩衝時間計算可能有邊界情況 |
| **例外審核** | 低 | 有狀態機保護，但重疊檢查失敗時的回滾需確認 |
| **併發控制** | 中 | 使用 DB Row Lock，但需確認所有時段操作都有鎖定 |
| **緩衝覆寫** | 中 | allow_buffer_override 開關存在，需確保管理員知悉風險 |

### 5.2 計費相關

| 項目 | 狀態 |
|:---|:---|
| Plan Level 欄位存在 | ✅ |
| 配額檢查 (LIMIT_EXCEEDED) | ✅ |
| 配額超限後的處理邏輯 | ⚠️ 部分缺失 |

---

## 6. 測試覆蓋評估

### 6.1 單元測試

| 模組 | 測試檔案 | 覆蓋率 |
|:---|:---|:---|
| TeacherProfile | teacher_profile_test.go | 19/19 通過 |
| Scheduling Validation | scheduling_validation_test.go | 部分覆蓋 |
| Smart Matching | smart_matching_test.go | 部分覆蓋 |
| Integration | integration_full_workflow_test.go | 端到端測試 |

### 6.2 建議加強

- ScheduleService 完整單元測試
- Exception State Machine 測試
- 併發控制測試
- Buffer Override 邊界測試

---

*本文件基於實際程式碼生成，最後更新：2026-02-12*
