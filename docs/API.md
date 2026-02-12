# TimeLedger API 參考文檔

## 1. API 基礎資訊

### 1.1 基礎 URL
```
Production: https://api.timeledger.app
Development: http://localhost:8888
```

### 1.2 認證方式
- **管理員**：Header `Authorization: Bearer <JWT_TOKEN>`
- **教師**：Header `Authorization: Bearer <JWT_TOKEN>` (LINE id_token)
- **公開 API**：無需認證

### 1.3 通用 Headers
```
Content-Type: application/json
Accept-Language: zh-TW | en
```

### 1.4 響應格式
```json
{
  "code": 0,
  "message": "OK",
  "data": { ... }
}
```

### 1.5 錯誤響應
```json
{
  "code": 40001,
  "message": "找不到資源",
  "data": null
}
```

---

## 2. API 端點總覽

### 2.1 認證 (Auth)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| POST | `/api/v1/auth/admin/login` | ❌ | 管理員 Email 登入 |
| POST | `/api/v1/auth/teacher/line/login` | ❌ | 教師 LINE 登入 |
| POST | `/api/v1/auth/refresh` | ✅ | 刷新 Token |
| POST | `/api/v1/auth/logout` | ✅ | 登出 |

### 2.2 地理 (Geo)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/geo/cities` | ❌ | 取得縣市列表 |

### 2.3 教師 - 個人檔案 (Teacher Profile)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/teacher/me/profile` | ✅ | 取得個人檔案 |
| PUT | `/api/v1/teacher/me/profile` | ✅ | 更新個人檔案 |
| GET | `/api/v1/teacher/me/centers` | ✅ | 取得已加入中心列表 |
| GET | `/api/v1/teacher/me/skills` | ✅ | 取得技能列表 |
| POST | `/api/v1/teacher/me/skills` | ✅ | 新增技能 |
| PUT | `/api/v1/teacher/me/skills/:id` | ✅ | 更新技能 |
| DELETE | `/api/v1/teacher/me/skills/:id` | ✅ | 刪除技能 |
| GET | `/api/v1/teacher/me/certificates` | ✅ | 取得證照列表 |
| POST | `/api/v1/teacher/me/certificates` | ✅ | 新增證照 |
| PUT | `/api/v1/teacher/me/certificates/:id` | ✅ | 更新證照 |
| POST | `/api/v1/teacher/me/certificates/upload` | ✅ | 上傳證照圖片 |
| DELETE | `/api/v1/teacher/me/certificates/:id` | ✅ | 刪除證照 |

### 2.4 標籤 (Hashtag)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/hashtags/search` | ✅ | 搜尋標籤 |
| POST | `/api/v1/hashtags` | ✅ | 建立標籤 |

### 2.5 教師 - 私人行程 (Teacher Personal Events)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/teacher/me/personal-events` | ✅ | 取得私人行程列表 |
| POST | `/api/v1/teacher/me/personal-events` | ✅ | 建立私人行程 |
| PATCH | `/api/v1/teacher/me/personal-events/:id` | ✅ | 更新私人行程 |
| DELETE | `/api/v1/teacher/me/personal-events/:id` | ✅ | 刪除私人行程 |
| GET | `/api/v1/teacher/me/personal-events/:id/note` | ✅ | 取得行程備註 |
| PUT | `/api/v1/teacher/me/personal-events/:id/note` | ✅ | 更新行程備註 |

### 2.6 教師 - 排課 (Teacher Schedule)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/teacher/me/schedule` | ✅ | 取得個人課表 |
| GET | `/api/v1/teacher/schedules` | ✅ | 取得所有排課 |
| GET | `/api/v1/teacher/me/centers/:center_id/schedule-rules` | ✅ | 取得中心排課規則 |

### 2.7 教師 - 課程筆記 (Teacher Sessions)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/teacher/sessions/note` | ✅ | 取得課程筆記 |
| PUT | `/api/v1/teacher/sessions/note` | ✅ | 新增/更新課程筆記 |

### 2.8 教師 - 例外 (Teacher Exceptions)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/teacher/exceptions` | ✅ | 取得例外列表 |
| POST | `/api/v1/teacher/exceptions` | ✅ | 提交例外申請 |
| POST | `/api/v1/teacher/exceptions/:id/revoke` | ✅ | 撤回例外 |

### 2.9 教師 - 排課編輯 (Teacher Scheduling)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| POST | `/api/v1/teacher/scheduling/check-rule-lock` | ✅ | 檢查規則鎖定狀態 |
| POST | `/api/v1/teacher/scheduling/preview-recurrence-edit` | ✅ | 預覽循環編輯影響範圍 |
| POST | `/api/v1/teacher/scheduling/edit-recurring` | ✅ | 編輯循環排課 |
| POST | `/api/v1/teacher/scheduling/delete-recurring` | ✅ | 刪除循環排課 |

### 2.10 教師 - 邀請 (Teacher Invitations)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/teacher/me/invitations` | ✅ | 取得邀請列表 |
| POST | `/api/v1/teacher/me/invitations/respond` | ✅ | 回應邀請 |
| GET | `/api/v1/teacher/me/invitations/pending-count` | ✅ | 取得待處理邀請數 |

### 2.11 教師 - 公開註冊

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| POST | `/api/v1/teacher/public/register` | ❌ | LINE Bot 自主註冊 |

### 2.12 管理員 - 教師管理 (Admin Teacher)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/teachers` | ✅ | 取得教師列表 |
| POST | `/api/v1/admin/teachers/placeholder` | ✅ | 建立暫存教師 |
| POST | `/api/v1/admin/teachers/merge` | ✅ | 合併教師 |
| DELETE | `/api/v1/teachers/:id` | ✅ | 刪除教師 |
| DELETE | `/api/v1/admin/centers/:id/teachers/:teacher_id` | ✅ | 從中心移除教師 |
| POST | `/api/v1/admin/centers/:id/invitations` | ✅ | 邀請教師 |

### 2.13 管理員 - 中心管理 (Admin Center)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/admin/centers` | ✅ | 取得中心列表 |
| GET | `/api/v1/admin/centers/:id/settings` | ✅ | 取得中心設定 |
| POST | `/api/v1/admin/centers` | ✅ | 建立中心 |
| PATCH | `/api/v1/admin/centers/:center_id/settings` | ✅ | 更新中心設定 |

### 2.14 管理員 - 教師資源

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/admin/teachers` | ✅ | 取得教師資源列表 |

### 2.15 管理員 - 班別 (Offering)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/admin/offerings` | ✅ | 取得班別列表 |
| POST | `/api/v1/admin/offerings` | ✅ | 建立班別 |
| PUT | `/api/v1/admin/offerings/:offering_id` | ✅ | 更新班別 |
| DELETE | `/api/v1/admin/offerings/:offering_id` | ✅ | 刪除班別 |
| POST | `/api/v1/admin/centers/:id/offerings/:offering_id/copy` | ✅ | 複製班別 |

### 2.16 管理員 - 課表範本 (Timetable Template)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/admin/templates` | ✅ | 取得範本列表 |
| POST | `/api/v1/admin/templates` | ✅ | 建立範本 |
| PUT | `/api/v1/admin/templates/:templateId` | ✅ | 更新範本 |
| DELETE | `/api/v1/admin/templates/:templateId` | ✅ | 刪除範本 |
| GET | `/api/v1/admin/templates/:templateId/cells` | ✅ | 取得範本細胞 |
| POST | `/api/v1/admin/templates/:templateId/cells` | ✅ | 建立範本細胞 |
| PUT | `/api/v1/admin/templates/:templateId/cells/reorder` | ✅ | 重新排序細胞 |
| DELETE | `/api/v1/admin/templates/cells/:cellId` | ✅ | 刪除細胞 |
| POST | `/api/v1/admin/templates/:templateId/apply` | ✅ | 套用範本 |
| POST | `/api/v1/admin/templates/:templateId/validate-apply` | ✅ | 驗證套用 |

### 2.17 管理員 - LINE 綁定

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/admin/me/line-binding` | ✅ | 取得 LINE 綁定狀態 |
| POST | `/api/v1/admin/me/line/bind` | ✅ | 產生綁定驗證碼 |
| DELETE | `/api/v1/admin/me/line/unbind` | ✅ | 解除 LINE 綁定 |
| GET | `/api/v1/admin/me/line/notify-settings` | ✅ | 取得通知設定 |
| PATCH | `/api/v1/admin/me/line/notify-settings` | ✅ | 更新通知設定 |
| GET | `/api/v1/admin/me/line/qrcode` | ✅ | 取得綁定 QR Code |
| GET | `/api/v1/admin/me/line/qrcode-with-code` | ✅ | 取得驗證碼 QR Code |

### 2.18 管理員 - Profile

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/admin/me/profile` | ✅ | 取得管理員 Profile |
| POST | `/api/v1/admin/me/change-password` | ✅ | 變更密碼 |

### 2.19 管理員 - 管理 (Admin Management)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/admin/admins` | ✅ | 取得管理員列表 |
| POST | `/api/v1/admin/admins` | ✅ | 建立管理員 |
| POST | `/api/v1/admin/admins/toggle-status` | ✅ | 切換管理員狀態 |
| POST | `/api/v1/admin/admins/reset-password` | ✅ | 重設管理員密碼 |
| POST | `/api/v1/admin/admins/change-role` | ✅ | 變更管理員角色 |

### 2.20 管理員 - 排課驗證

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| POST | `/api/v1/admin/scheduling/check-overlap` | ✅ | 檢查時段重疊 |
| POST | `/api/v1/admin/scheduling/check-teacher-buffer` | ✅ | 檢查老師緩衝 |
| POST | `/api/v1/admin/scheduling/check-room-buffer` | ✅ | 檢查教室緩衝 |
| POST | `/api/v1/admin/scheduling/validate` | ✅ | 完整驗證 |

### 2.21 管理員 - 排課 Dashboard

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/admin/dashboard/today-summary` | ✅ | 今日摘要 |
| GET | `/api/v1/admin/rules` | ✅ | 取得排課規則 |
| POST | `/api/v1/admin/rules` | ✅ | 建立排課規則 |
| PUT | `/api/v1/admin/rules/:ruleId` | ✅ | 更新排課規則 |
| DELETE | `/api/v1/admin/rules/:ruleId` | ✅ | 刪除排課規則 |

### 2.22 管理員 - 例外管理

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| POST | `/api/v1/admin/scheduling/exceptions` | ✅ | 建立例外 |
| POST | `/api/v1/admin/scheduling/exceptions/:exceptionId/review` | ✅ | 審核例外 |
| GET | `/api/v1/admin/rules/:ruleId/exceptions` | ✅ | 取得規則的例外 |
| GET | `/api/v1/admin/exceptions` | ✅ | 取得例外列表（日期範圍） |
| GET | `/api/v1/admin/exceptions/pending` | ✅ | 取得待審核例外 |
| GET | `/api/v1/admin/exceptions/all` | ✅ | 取得所有例外 |
| POST | `/api/v1/admin/expand-rules` | ✅ | 展開排課規則 |
| POST | `/api/v1/admin/detect-phase-transitions` | ✅ | 偵測階段轉換 |
| POST | `/api/v1/admin/scheduling/check-rule-lock` | ✅ | 檢查規則鎖定 |

### 2.23 管理員 - Matrix View

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/admin/scheduling/matrix-view` | ✅ | 取得 Matrix 視圖 |

### 2.24 管理員 - 教室 (Admin Room)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/admin/rooms` | ✅ | 取得教室列表 |
| POST | `/api/v1/admin/rooms` | ✅ | 建立教室 |
| PUT | `/api/v1/admin/rooms/:room_id` | ✅ | 更新教室 |
| GET | `/api/v1/admin/rooms/active` | ✅ | 取得有效教室 |
| PATCH | `/api/v1/admin/rooms/:room_id/toggle-active` | ✅ | 切換教室狀態 |

### 2.25 管理員 - 課程 (Admin Course)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/admin/courses` | ✅ | 取得課程列表 |
| POST | `/api/v1/admin/courses` | ✅ | 建立課程 |
| PUT | `/api/v1/admin/courses/:course_id` | ✅ | 更新課程 |
| DELETE | `/api/v1/admin/courses/:course_id` | ✅ | 刪除課程 |
| GET | `/api/v1/admin/courses/active` | ✅ | 取得有效課程 |
| PATCH | `/api/v1/admin/courses/:course_id/toggle-active` | ✅ | 切換課程狀態 |

### 2.26 管理員 - 假期 (Admin Holiday)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/admin/centers/:id/holidays` | ✅ | 取得假期列表 |
| POST | `/api/v1/admin/centers/:id/holidays` | ✅ | 建立假期 |
| DELETE | `/api/v1/admin/centers/:id/holidays/:holiday_id` | ✅ | 刪除假期 |
| POST | `/api/v1/admin/centers/:id/holidays/bulk` | ✅ | 批量建立假期 |

### 2.27 管理員 - 學期 (Admin Term)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/admin/terms` | ✅ | 取得學期列表 |
| GET | `/api/v1/admin/terms/active` | ✅ | 取得有效學期 |
| POST | `/api/v1/admin/terms` | ✅ | 建立學期 |
| PUT | `/api/v1/admin/terms/:term_id` | ✅ | 更新學期 |
| DELETE | `/api/v1/admin/terms/:term_id` | ✅ | 刪除學期 |
| GET | `/api/v1/admin/occupancy/rules` | ✅ | 取得佔用規則 |
| POST | `/api/v1/admin/terms/copy-rules` | ✅ | 複製規則 |

### 2.28 管理員 - 教師備註

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/admin/teachers/:teacher_id/note` | ✅ | 取得教師備註 |
| PUT | `/api/v1/admin/teachers/:teacher_id/note` | ✅ | 新增/更新備註 |
| DELETE | `/api/v1/admin/teachers/:teacher_id/note` | ✅ | 刪除備註 |

### 2.29 智慧媒合 (Smart Matching)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| POST | `/api/v1/admin/smart-matching/matches` | ✅ | 智慧媒合搜尋 |
| GET | `/api/v1/admin/smart-matching/talent/search` | ✅ | 人才庫搜尋 |
| GET | `/api/v1/admin/smart-matching/talent/stats` | ✅ | 人才庫統計 |
| POST | `/api/v1/admin/smart-matching/talent/invite` | ✅ | 邀請人才 |
| GET | `/api/v1/admin/smart-matching/suggestions` | ✅ | 搜尋建議 |
| POST | `/api/v1/admin/smart-matching/alternatives` | ✅ | 替代時段建議 |
| GET | `/api/v1/admin/teachers/:teacher_id/sessions` | ✅ | 教師課表查詢 |

### 2.30 通知 (Notification)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/notifications` | ✅ | 取得通知列表 |
| GET | `/api/v1/notifications/unread-count` | ✅ | 取得未讀數 |
| POST | `/api/v1/notifications/:id/read` | ✅ | 標記已讀 |
| POST | `/api/v1/notifications/read-all` | ✅ | 全部標記已讀 |
| POST | `/api/v1/notifications/token` | ✅ | 設定推播 Token |
| POST | `/api/v1/notifications/test` | ✅ | 測試推播 |
| GET | `/api/v1/admin/notifications/queue-stats` | ✅ | 通知佇列統計 |

### 2.31 管理員 - 廣播

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| POST | `/api/v1/admin/notifications/broadcast` | ✅ | 廣播通知 |

### 2.32 匯出 (Export)

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| POST | `/api/v1/admin/export/schedule/csv` | ✅ | 匯出課表 CSV |
| POST | `/api/v1/admin/export/schedule/pdf` | ✅ | 匯出課表 PDF |
| GET | `/api/v1/admin/centers/:id/export/teachers/csv` | ✅ | 匯出教師 CSV |
| GET | `/api/v1/admin/centers/:id/export/exceptions/csv` | ✅ | 匯出例外 CSV |
| GET | `/api/v1/teacher/me/schedule.ics` | ✅ | 匯出 ICS |
| POST | `/api/v1/teacher/me/schedule/subscription` | ✅ | 建立日曆訂閱 |
| DELETE | `/api/v1/teacher/me/schedule/subscription` | ✅ | 取消日曆訂閱 |
| GET | `/api/v1/teacher/me/schedule/image` | ✅ | 匯出課表圖片 |
| POST | `/api/v1/teacher/me/backgrounds` | ✅ | 上傳背景圖片 |
| GET | `/api/v1/teacher/me/backgrounds` | ✅ | 取得背景圖片列表 |
| DELETE | `/api/v1/teacher/me/backgrounds/:id` | ✅ | 刪除背景圖片 |

### 2.33 公開 - 日曆訂閱

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/calendar/subscribe/:token.ics` | ❌ | 公開日曆訂閱 |

### 2.34 公開 - 邀請

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/v1/invitations/:token` | ❌ | 取得公開邀請資訊 |
| POST | `/api/v1/invitations/:token/accept` | ❌ | 接受邀請連結 |

### 2.35 LINE Bot

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| POST | `/api/v1/line/webhook` | ❌ | LINE Webhook |
| GET | `/api/v1/line/health` | ❌ | LINE Bot 健康檢查 |

### 2.36 測試

| Method | Endpoint | 認證 | 說明 |
|:---:|:---|:---|:---|
| GET | `/api/test/r2-status` | ❌ | R2 狀態測試 |
| POST | `/api/test/upload` | ❌ | R2 上傳測試 |
| POST | `/api/test/upload-batch` | ❌ | R2 批量上傳測試 |

---

## 3. 查詢參數

### 3.1 分頁參數

| 參數 | 類型 | 必填 | 預設值 | 說明 |
|:---|:---:|:---:|:---:|:---|
| `page` | INT | 否 | 1 | 頁碼 |
| `limit` | INT | 否 | 20 | 每頁筆數（最大 100） |
| `sort_by` | STRING | 否 | 依各 API 定義 | 排序欄位 |
| `sort_order` | STRING | 否 | ASC | 排序方向 |

### 3.2 日期範圍參數

| 參數 | 類型 | 必填 | 說明 |
|:---|:---:|:---:|:---|
| `from` | DATE | 是 | 開始日期 (YYYY-MM-DD) |
| `to` | DATE | 是 | 結束日期 (YYYY-MM-DD) |

---

*本文件基於實際程式碼生成，最後更新：2026-02-12*
