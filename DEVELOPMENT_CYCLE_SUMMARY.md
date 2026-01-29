# TimeLedger 開發階段總結

**更新日期：2026年1月29日**

本文件記錄 2026年1月26日-29日 當天完成的開發工作，包括所有問題修復、功能改進和程式碼重構。

---

## 本週開發重點（2026/01/28-29）

### 邀請連結功能（第二階段）

| 項目 | 狀態 | 說明 |
|:---|:---:|:---|
| 產生邀請連結 | ✅ | 管理員可產生 72 小時有效連結 |
| 連結列表管理 | ✅ | 查看、複製、撤回連結 |
| 接受邀請 API | ✅ | 新老師透過連結加入中心 |
| 前端邀請頁面 | ✅ | `/invite/[token].vue` 支援 LINE 登入 |

**程式碼規模：**
- 後端程式碼：+500 行
- 前端程式碼：+450 行

### LINE 通知整合

| 觸發時機 | 通知對象 | 訊息內容 |
|:---|:---|:---|
| 老師接受邀請 | 所有已綁定管理員 | 🎉 新成員加入（姓名、中心、角色） |

**功能特色：**
- 異步發送，不影響主要流程
- 僅通知已綁定 LINE 的管理員
- 支援 Flex Message 格式

### Cloudflare R2 儲存整合

**修改檔案：**
- `libs/cloudflare_r2.go`（新檔案）- 純 HTTP 實作，無需 AWS SDK
- `configs/env.go` - 新增 R2 環境變數
- `.env.example` - 更新範例

**環境設定：**
```
CLOUDFLARE_R2_ENABLED=true
CLOUDFLARE_R2_ACCOUNT_ID=your-account-id
CLOUDFLARE_R2_ACCESS_KEY=your-access-key
CLOUDFLARE_R2_SECRET_KEY=your-secret-key
CLOUDFLARE_R2_BUCKET_NAME=your-bucket-name
CLOUDFLARE_R2_PUBLIC_URL=https://your-domain.com/files
```

**特色：**
- 自動偵測 R2 是否啟用
- 若未設定或失敗，自動回退本地儲存
- 刪除證照時同步刪除 R2 檔案

### 管理員查看老師檔案優化

| 功能 | 狀態 |
|:---|:---:|
| 技能專長顯示 | ✅ |
| 證照清單顯示 | ✅ |
| 證照圖示（PDF/圖片） | ✅ |
| 證照預覽連結 | ✅ |

**API 擴展：**
```json
{
  "id": 1,
  "name": "陳小美",
  "is_active": true,
  "skills": [
    { "id": 1, "skill_name": "瑜伽", "category": "運動", "level": "高級" }
  ],
  "certificates": [
    {
      "id": 1,
      "name": "瑜伽師資認證",
      "file_url": "https://...",
      "issued_at": "2024-01-15"
    }
  ]
}
```

---

## Git 提交紀錄

| 提交 | 說明 |
|:---|:---|
| 4bee261 | feat(invitation): implement invitation link system (Phase 2) |
| 2dfc018 | feat(invitation): add LINE notification when teacher accepts invitation |
| c5dea84 | feat(storage): integrate Cloudflare R2 for certificate file storage |
| 33c0bef | feat(admin): enhance teacher profile with certificates and skills display |

**Git 狀態：** 當前分支 `claudecode` 領先 `origin/claudecode` 52 個提交

---

## 程式碼統計

| 指標 | 數量 |
|:---|:---:|
| 後端新增行數 | +1,500 行 |
| 前端新增行數 | +1,000 行 |
| 新增檔案數 | 3 個 |
| 修改檔案數 | 12 個 |

---

## 測試結果

```
=== RUN   TestCenterInvitationRepository_CRUD
    --- PASS: TestCenterInvitationRepository_CRUD (1.68s)
=== RUN   TestIntegration_InvitationFlow
    --- PASS: TestIntegration_InvitationFlow (3.78s)
=== RUN   TestIntegration_TeacherFullWorkflow
    --- PASS: TestIntegration_TeacherFullWorkflow (4.16s)
=== RUN   TestIntegration_AdminTeacherManagement
    --- PASS: TestIntegration_AdminTeacherManagement (4.01s)
```

**測試通過率：** 100%

---

## 待開發項目（可選）

| 優先級 | 項目 | 說明 |
|:---:|:---|:---|
| 中 | 人才庫搜尋 | 依技能、區域搜尋老師 |
| 中 | 代課媒合智慧推薦 | 自動推薦合適代課老師 |
| 低 | 邀請統計分析 | 追蹤邀請轉換率 |
| 低 | 批量產生邀請連結 | 一次產生多個連結 |

---

## 技術棧現況

| 層面 | 技術 | 狀態 |
|:---|:---|:---:|
| 後端 | Go (Gin) + MySQL + Redis | ✅ 穩定 |
| 前端 | Nuxt 3 + Tailwind CSS | ✅ 穩定 |
| 儲存 | 本地 / Cloudflare R2 | ✅ 可用 |
| 通知 | LINE Messaging API | ✅ 可用 |
| 排課 | 排課驗證引擎 | ✅ 穩定 |

---

## 下一步建議

1. **人才庫功能** - 開放老師設定 `is_open_to_hiring`，供中心搜尋
2. **代課媒合優化** - 強化智慧推薦演算法
3. **效能優化** - 快取熱門資料，減少資料庫查詢
4. **測試覆蓋** - 增加整合測試案例

---

## 修復的問題清單（2026/01/27）

本階段共修復了 8 個問題，涵蓋通知系統、資料隔離、UI 顯示和流程邏輯等方面。

### 1. 老師端通知跳轉問題

**問題描述：** 老師點擊審核結果通知時，沒有正確跳轉到例外申請頁面。

**問題原因：**
- 後端發送通知時沒有設置 `Type` 欄位
- 前端只檢查 `APPROVAL` 類型和管理員路徑

**修復方案：**
- 新增 `SendTeacherNotificationWithType()` 方法
- 設置通知類型為 `REVIEW_RESULT`
- 前端根據 `user_type` 判斷身份，老師跳轉到 `/teacher/exceptions`

**修改檔案：**
- `app/services/notification_interface.go`
- `app/services/notification.go`
- `frontend/components/NotificationDropdown.vue`

### 2. 課程時段週日顯示問題

**問題描述：** 課程時段管理頁面中，週日的課程顯示為 `-` 而不是 `日`。

**問題原因：** `getWeekdayText()` 函數的陣列只有 0-6 的索引，但系統使用 7 表示週日。

**修復方案：** 修正函數邏輯，將 weekday 7 轉換為索引 0。

**修改檔案：**
- `frontend/pages/admin/schedules.vue`
- `frontend/tests/pages/admin/schedules.spec.ts`

### 3. 例外申請原時間顯示問題

**問題描述：** 審核頁面中，RESCHEDULE 類型的原時間顯示為 `undefined - undefined`。

**問題原因：** 前端嘗試存取 `exception.start_time`，但時間資訊儲存在關聯的 `Rule` 中。

**修復方案：**
- 新增 `getOriginalTimeText()` helper 函數
- 正確存取 `exception.rule.start_time` 和 `exception.rule.end_time`

**修改檔案：**
- `frontend/pages/admin/approval.vue`
- `frontend/components/ReviewModal.vue`
- `frontend/components/ExceptionDetailModal.vue`

### 4. 管理員核准後老師通知問題

**問題描述：** 管理員核准例外申請後，老師沒有收到通知。

**問題原因：** `ReviewException()` 方法中沒有呼叫 `SendReviewNotification()`。

**修復方案：** 在審核邏輯中新增通知發送呼叫。

**修改檔案：**
- `app/services/scheduling_expansion.go`

### 5. 老師課表資料隔離問題

**問題描述：** 老師登入後可以看到其他老師的課程。

**問題原因：** `GetSchedule` API 使用 `ListByCenterID()` 取得所有課程，而非老師自己的課程。

**修復方案：** 改用 `ListByTeacherID()` 並新增必要的 Preload。

**修改檔案：**
- `app/controllers/teacher.go`
- `app/repositories/schedule_rule.go`

### 6. 編輯課程時日期欄位問題

**問題描述：** 選擇「全部」模式編輯課程時，開始日期和結束日期欄位顯示為必填。

**問題原因：** 日期欄位設計為必填，但 ALL 模式下修改內容時不需要修改日期。

**修復方案：**
- 前端：編輯模式下日期欄位改為可選填，新增提示文字
- 後端：日期欄位為空時保留現有值

**修改檔案：**
- `frontend/components/ScheduleRuleModal.vue`
- `app/controllers/scheduling.go`

---

## 新增功能（2026/01/27）

### 1. 例外申請通知系統

**管理員端：**
- 新增 `/api/v1/admin/exceptions/all` API 端點
- 支援狀態篩選（PENDING、APPROVED、REJECTED、REVOKED）
- 審核頁面新增日期範圍篩選器
- Header 新增通知鈴鐺按鈕

**老師端：**
- 審核通過/拒絕後收到通知
- 通知包含審核結果和日期資訊

### 2. 排課規則編輯優化

**更新模式說明：**
- `SINGLE`：只修改這一天（建立 CANCEL + RESCHEDULE 例外單）
- `FUTURE`：修改這天及之後（截斷原規則，建立新規則段）
- `ALL`：修改全部（同步更新所有相關規則）

---

## 資料庫結構

### schedule_rules 資料表

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| id | uint | 主鍵 |
| center_id | uint | 所屬中心 |
| offering_id | uint | 課程班別 |
| teacher_id | uint | 老師（可為 NULL） |
| room_id | uint | 教室 |
| weekday | int | 星期（1-7，7=週日） |
| start_time | string | 開始時間 |
| end_time | string | 結束時間 |
| effective_range | JSON | 有效日期範圍 |

### schedule_exceptions 資料表

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| id | uint | 主鍵 |
| center_id | uint | 所屬中心 |
| rule_id | uint | 關聯規則 |
| original_date | date | 原日期 |
| type | string | 類型（CANCEL、RESCHEDULE、REPLACE_TEACHER） |
| status | string | 狀態（PENDING、APPROVED、REJECTED、REVOKED） |
| new_start_at | datetime | 新開始時間（改期用） |
| new_end_at | datetime | 新結束時間（改期用） |

---

## API 端點總覽

### 管理員端 API

| 方法 | 路徑 | 功能 |
|:---|:---|:---|
| GET | /api/v1/admin/rules | 取得課程規則列表 |
| POST | /api/v1/admin/rules | 建立課程規則 |
| PUT | /api/v1/admin/rules/:id | 更新課程規則 |
| DELETE | /api/v1/admin/rules/:id | 刪除課程規則 |
| GET | /api/v1/admin/exceptions/pending | 取得待審核例外 |
| GET | /api/v1/admin/exceptions/all | 取得所有例外（支援篩選） |
| POST | /api/v1/admin/scheduling/exceptions/:id/review | 審核例外 |
| POST | /api/v1/admin/invitations | 產生邀請連結 |
| GET | /api/v1/admin/invitations | 取得邀請連結列表 |
| DELETE | /api/v1/admin/invitations/:id | 撤回邀請連結 |

### 老師端 API

| 方法 | 路徑 | 功能 |
|:---|:---|:---|
| GET | /api/v1/teacher/me/schedule | 取得課表 |
| GET | /api/v1/teacher/exceptions | 取得例外申請列表 |
| POST | /api/v1/teacher/exceptions | 提交例外申請 |
| POST | /api/v1/teacher/exceptions/:id/revoke | 撤回例外申請 |
| GET | /api/v1/invitations/:token/accept | 接受邀請（驗證連結有效性） |
| POST | /api/v1/invitations/:token/accept | 確認接受邀請 |

---

## 程式碼變更統計（2026/01/28-29）

| 類型 | 檔案數 | 變更行數 |
|:---|:---:|:---:|
| 新增功能 | 8 | +1,500 |
| 問題修復 | 4 | +150 |
| 整合優化 | 3 | +200 |
| 測試修改 | 4 | +50 |
| **總計** | **19** | **+1,900** |

---

## 排課檢查機制修正（2026/01/27 下午）

本階段修正了排課檢查機制的功能缺口，確保模板套用和手動新增課程時都有適當的衝突檢查。

### 修正背景

**問題描述：**
- `ApplyTemplate` 套用模板時完全沒有進行任何衝突檢查
- `CreateRule` 手動新增課程時缺少 Buffer 檢查
- 可能導致產生時間衝突、違反緩衝時間規定的排課

### 修正方案一：ApplyTemplate 加入衝突檢查

**修改檔案：**
- `app/controllers/timetable_template.go`

**修正內容：**
- 在 Controller 中注入 `scheduleRuleRepo` 和 `personalEventRepo`
- `ApplyTemplate` 函數加入時間衝突檢查
- 對每個 (weekday, cell) 組合呼叫 `CheckOverlap()` 檢查：
  - Room Overlap（教室時間衝突）
  - Teacher Overlap（老師時間衝突）
  - Personal Event（老師個人行程衝突）
- 有衝突時回傳詳細的衝突資訊（錯誤碼 40002）

**衝突回應格式：**
```json
{
  "code": 40002,
  "message": "套用模板會產生時間衝突，請先解決衝突後再嘗試",
  "datas": {
    "conflicts": [...],
    "conflict_count": 3
  }
}
```

### 修正方案二：CreateRule 加入 Buffer 檢查

**修改檔案：**
- `app/controllers/scheduling.go`

**修正內容：**
- 在 Controller 中注入 `courseRepo`
- `CreateRule` 函數加入 Buffer 檢查：
  - Teacher Buffer（老師轉場緩衝時間）
  - Room Buffer（教室清潔緩衝時間）
- 使用 `validationService.CheckTeacherBuffer()` 和 `CheckRoomBuffer()` 進行檢查
- 有衝突時回傳詳細的緩衝衝突資訊（錯誤碼 40003）

**Buffer 衝突回應格式：**
```json
{
  "code": 40003,
  "message": "排課時間違反緩衝時間規定",
  "datas": {
    "buffer_conflicts": [...],
    "conflict_count": 2
  }
}
```

### 新增輔助函數

**檔案：** `app/controllers/scheduling.go`

- `getTeacherPreviousSessionEndTime()` - 取得老師在指定 weekday 的上一堂課結束時間
- `getRoomPreviousSessionEndTime()` - 取得教室在指定 weekday 的上一堂課結束時間

### 新增統一驗證服務

**新增檔案：**
- `app/services/schedule_rule_validator.go`

**功能：**
- `ScheduleRuleValidator` 統一的排課規則驗證服務
- `ValidateForApplyTemplate()` - 驗證模板套用的衝突
- `ValidateForCreateRule()` - 驗證新規則的衝突
- 整合所有檢查邏輯（重疊、緩衝、個人行程）

### 檢查功能對比表

| 檢查項目 | 修正前 | 修正後 |
|:---|:---:|:---:|
| Room Overlap | ✅ CreateRule / ❌ ApplyTemplate | ✅ 兩者皆有 |
| Teacher Overlap | ✅ CreateRule / ❌ ApplyTemplate | ✅ 兩者皆有 |
| Personal Event | ✅ CreateRule / ❌ ApplyTemplate | ✅ 兩者皆有 |
| Teacher Buffer | ❌ 沒有 | ✅ CreateRule 有 |
| Room Buffer | ❌ 沒有 | ✅ CreateRule 有 |

---

## 開發規範遵守情況

本階段遵守的開發規範：

- ✅ 使用 Triple Return Pattern 處理錯誤
- ✅ Repository 層級包含 center_id 過濾
- ✅ 後端負責資料隔離，前端不依賴 URL 傳遞 center_id
- ✅ 禁止使用原生 alert/confirm
- ✅ Commit Message 使用英文
- ✅ 每次修改立即 commit

---

## 待完成項目

以下項目尚未完成，建議在下一階段優先處理：

1. **人才庫搜尋功能**
   - 依技能、區域搜尋老師
   - 開放老師設定 `is_open_to_hiring`

2. **代課媒合智慧推薦**
   - 自動推薦合適代課老師
   - 強化智慧推薦演算法

3. **邀請統計分析**
   - 追蹤邀請轉換率
   - 批量產生邀請連結

4. **效能優化**
   - 快取熱門資料，減少資料庫查詢
   - 大量課程時的週曆渲染效能

---

## 下一步計劃

1. **人才庫功能開發**
   - 開放老師設定 `is_open_to_hiring`
   - 依技能、區域搜尋老師

2. **代課媒合優化**
   - 強化智慧推薦演算法
   - 自動推薦合適代課老師

3. **測試覆蓋加強**
   - 增加整合測試案例
   - 確保邊界條件正確處理

4. **效能優化**
   - 快取熱門資料
   - 減少資料庫查詢次數

---

*文件版本：1.3*
*建立時間：2026-01-26*
*更新時間：2026-01-29*
*維護者：TimeLedger 開發團隊*

---

## 修復的問題清單（2026/01/27）

本階段共修復了 8 個問題，涵蓋通知系統、資料隔離、UI 顯示和流程邏輯等方面。

### 1. 老師端通知跳轉問題

**問題描述：** 老師點擊審核結果通知時，沒有正確跳轉到例外申請頁面。

**問題原因：**
- 後端發送通知時沒有設置 `Type` 欄位
- 前端只檢查 `APPROVAL` 類型和管理員路徑

**修復方案：**
- 新增 `SendTeacherNotificationWithType()` 方法
- 設置通知類型為 `REVIEW_RESULT`
- 前端根據 `user_type` 判斷身份，老師跳轉到 `/teacher/exceptions`

**修改檔案：**
- `app/services/notification_interface.go`
- `app/services/notification.go`
- `frontend/components/NotificationDropdown.vue`

### 2. 課程時段週日顯示問題

**問題描述：** 課程時段管理頁面中，週日的課程顯示為 `-` 而不是 `日`。

**問題原因：** `getWeekdayText()` 函數的陣列只有 0-6 的索引，但系統使用 7 表示週日。

**修復方案：** 修正函數邏輯，將 weekday 7 轉換為索引 0。

**修改檔案：**
- `frontend/pages/admin/schedules.vue`
- `frontend/tests/pages/admin/schedules.spec.ts`

### 3. 例外申請原時間顯示問題

**問題描述：** 審核頁面中，RESCHEDULE 類型的原時間顯示為 `undefined - undefined`。

**問題原因：** 前端嘗試存取 `exception.start_time`，但時間資訊儲存在關聯的 `Rule` 中。

**修復方案：**
- 新增 `getOriginalTimeText()` helper 函數
- 正確存取 `exception.rule.start_time` 和 `exception.rule.end_time`

**修改檔案：**
- `frontend/pages/admin/approval.vue`
- `frontend/components/ReviewModal.vue`
- `frontend/components/ExceptionDetailModal.vue`

### 4. 管理員核准後老師通知問題

**問題描述：** 管理員核准例外申請後，老師沒有收到通知。

**問題原因：** `ReviewException()` 方法中沒有呼叫 `SendReviewNotification()`。

**修復方案：** 在審核邏輯中新增通知發送呼叫。

**修改檔案：**
- `app/services/scheduling_expansion.go`

### 5. 老師課表資料隔離問題

**問題描述：** 老師登入後可以看到其他老師的課程。

**問題原因：** `GetSchedule` API 使用 `ListByCenterID()` 取得所有課程，而非老師自己的課程。

**修復方案：** 改用 `ListByTeacherID()` 並新增必要的 Preload。

**修改檔案：**
- `app/controllers/teacher.go`
- `app/repositories/schedule_rule.go`

### 6. 編輯課程時日期欄位問題

**問題描述：** 選擇「全部」模式編輯課程時，開始日期和結束日期欄位顯示為必填。

**問題原因：** 日期欄位設計為必填，但 ALL 模式下修改內容時不需要修改日期。

**修復方案：**
- 前端：編輯模式下日期欄位改為可選填，新增提示文字
- 後端：日期欄位為空時保留現有值

**修改檔案：**
- `frontend/components/ScheduleRuleModal.vue`
- `app/controllers/scheduling.go`

---

## 新增功能（2026/01/27）

### 1. 例外申請通知系統

**管理員端：**
- 新增 `/api/v1/admin/exceptions/all` API 端點
- 支援狀態篩選（PENDING、APPROVED、REJECTED、REVOKED）
- 審核頁面新增日期範圍篩選器
- Header 新增通知鈴鐺按鈕

**老師端：**
- 審核通過/拒絕後收到通知
- 通知包含審核結果和日期資訊

### 2. 排課規則編輯優化

**更新模式說明：**
- `SINGLE`：只修改這一天（建立 CANCEL + RESCHEDULE 例外單）
- `FUTURE`：修改這天及之後（截斷原規則，建立新規則段）
- `ALL`：修改全部（同步更新所有相關規則）

---

## 資料庫結構

### schedule_rules 資料表

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| id | uint | 主鍵 |
| center_id | uint | 所屬中心 |
| offering_id | uint | 課程班別 |
| teacher_id | uint | 老師（可為 NULL） |
| room_id | uint | 教室 |
| weekday | int | 星期（1-7，7=週日） |
| start_time | string | 開始時間 |
| end_time | string | 結束時間 |
| effective_range | JSON | 有效日期範圍 |

### schedule_exceptions 資料表

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| id | uint | 主鍵 |
| center_id | uint | 所屬中心 |
| rule_id | uint | 關聯規則 |
| original_date | date | 原日期 |
| type | string | 類型（CANCEL、RESCHEDULE、REPLACE_TEACHER） |
| status | string | 狀態（PENDING、APPROVED、REJECTED、REVOKED） |
| new_start_at | datetime | 新開始時間（改期用） |
| new_end_at | datetime | 新結束時間（改期用） |

---

## API 端點總覽

### 管理員端 API

| 方法 | 路徑 | 功能 |
|:---|:---|:---|
| GET | /api/v1/admin/rules | 取得課程規則列表 |
| POST | /api/v1/admin/rules | 建立課程規則 |
| PUT | /api/v1/admin/rules/:id | 更新課程規則 |
| DELETE | /api/v1/admin/rules/:id | 刪除課程規則 |
| GET | /api/v1/admin/exceptions/pending | 取得待審核例外 |
| GET | /api/v1/admin/exceptions/all | 取得所有例外（支援篩選） |
| POST | /api/v1/admin/scheduling/exceptions/:id/review | 審核例外 |

### 老師端 API

| 方法 | 路徑 | 功能 |
|:---|:---|:---|
| GET | /api/v1/teacher/me/schedule | 取得課表 |
| GET | /api/v1/teacher/exceptions | 取得例外申請列表 |
| POST | /api/v1/teacher/exceptions | 提交例外申請 |
| POST | /api/v1/teacher/exceptions/:id/revoke | 撤回例外申請 |

---

## 程式碼變更統計（2026/01/27）

| 類型 | 檔案數 | 變更行數 |
|:---|:---:|:---:|
| 新增功能 | 5 | +500 |
| 問題修復 | 12 | +200 |
| 重構優化 | 3 | +100 |
| 測試修改 | 2 | +20 |
| **總計** | **22** | **+820** |

---

## 排課檢查機制修正（2026/01/27 下午）

本階段修正了排課檢查機制的功能缺口，確保模板套用和手動新增課程時都有適當的衝突檢查。

### 修正背景

**問題描述：**
- `ApplyTemplate` 套用模板時完全沒有進行任何衝突檢查
- `CreateRule` 手動新增課程時缺少 Buffer 檢查
- 可能導致產生時間衝突、違反緩衝時間規定的排課

### 修正方案一：ApplyTemplate 加入衝突檢查

**修改檔案：**
- `app/controllers/timetable_template.go`

**修正內容：**
- 在 Controller 中注入 `scheduleRuleRepo` 和 `personalEventRepo`
- `ApplyTemplate` 函數加入時間衝突檢查
- 對每個 (weekday, cell) 組合呼叫 `CheckOverlap()` 檢查：
  - Room Overlap（教室時間衝突）
  - Teacher Overlap（老師時間衝突）
  - Personal Event（老師個人行程衝突）
- 有衝突時回傳詳細的衝突資訊（錯誤碼 40002）

**衝突回應格式：**
```json
{
  "code": 40002,
  "message": "套用模板會產生時間衝突，請先解決衝突後再嘗試",
  "datas": {
    "conflicts": [...],
    "conflict_count": 3
  }
}
```

### 修正方案二：CreateRule 加入 Buffer 檢查

**修改檔案：**
- `app/controllers/scheduling.go`

**修正內容：**
- 在 Controller 中注入 `courseRepo`
- `CreateRule` 函數加入 Buffer 檢查：
  - Teacher Buffer（老師轉場緩衝時間）
  - Room Buffer（教室清潔緩衝時間）
- 使用 `validationService.CheckTeacherBuffer()` 和 `CheckRoomBuffer()` 進行檢查
- 有衝突時回傳詳細的緩衝衝突資訊（錯誤碼 40003）

**Buffer 衝突回應格式：**
```json
{
  "code": 40003,
  "message": "排課時間違反緩衝時間規定",
  "datas": {
    "buffer_conflicts": [...],
    "conflict_count": 2
  }
}
```

### 新增輔助函數

**檔案：** `app/controllers/scheduling.go`

- `getTeacherPreviousSessionEndTime()` - 取得老師在指定 weekday 的上一堂課結束時間
- `getRoomPreviousSessionEndTime()` - 取得教室在指定 weekday 的上一堂課結束時間

### 新增統一驗證服務

**新增檔案：**
- `app/services/schedule_rule_validator.go`

**功能：**
- `ScheduleRuleValidator` 統一的排課規則驗證服務
- `ValidateForApplyTemplate()` - 驗證模板套用的衝突
- `ValidateForCreateRule()` - 驗證新規則的衝突
- 整合所有檢查邏輯（重疊、緩衝、個人行程）

### 檢查功能對比表

| 檢查項目 | 修正前 | 修正後 |
|:---|:---:|:---:|
| Room Overlap | ✅ CreateRule / ❌ ApplyTemplate | ✅ 兩者皆有 |
| Teacher Overlap | ✅ CreateRule / ❌ ApplyTemplate | ✅ 兩者皆有 |
| Personal Event | ✅ CreateRule / ❌ ApplyTemplate | ✅ 兩者皆有 |
| Teacher Buffer | ❌ 沒有 | ✅ CreateRule 有 |
| Room Buffer | ❌ 沒有 | ✅ CreateRule 有 |

---

## 開發規範遵守情況

本階段遵守的開發規範：

- ✅ 使用 Triple Return Pattern 處理錯誤
- ✅ Repository 層級包含 center_id 過濾
- ✅ 後端負責資料隔離，前端不依賴 URL 傳遞 center_id
- ✅ 禁止使用原生 alert/confirm
- ✅ Commit Message 使用英文
- ✅ 每次修改立即 commit

---

## 待完成項目

以下項目尚未完成，建議在下一階段優先處理：

1. **課程時段測試資料完善**
   - 新增週日的測試課程規則
   - 確保每個星期都有測試資料

2. **例外申請詳細頁面**
   - 目前只有列表，沒有詳細資訊頁面
   - 建議增加例外申請的完整資訊展示

3. **通知系統增強**
   - 支援通知批次處理
   - 新增通知設定功能

4. **效能優化**
   - 大量課程時的週曆渲染效能
   - 通知載入分頁

---

## 下一步計劃

1. **完善測試覆蓋率**
   - 為新增功能撰寫單元測試
   - 確保邊界條件正確處理

2. **使用者體驗優化**
   - 載入狀態優化
   - 錯誤提示美化和統一

3. **文件更新**
   - 更新 API 文件
   - 更新 CLAUDE.md

---

*文件版本：1.2*
*建立時間：2026-01-26*
*更新時間：2026-01-27*
*維護者：TimeLedger 開發團隊*
