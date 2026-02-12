# TimeLedger API Documentation
**Version**: 1.1.2
**Base URL**: `http://your-server:8080`
**Generated**: 2026-02-13
**Total API Endpoints**: 171

---

## 📊 API Statistics

| Category | Endpoints | Description |
|:---|---:|:---|
| **Authentication** | 4 | 登入、登出、Token 刷新 |
| **Teacher APIs** | 38 | 教師個人資料、課表、技能、證照 |
| **Admin APIs** | 60+ | 管理員功能、中心管理、資源管理 |
| **Public APIs** | 4 | 邀請連結、日曆訂閱 |
| **LINE Bot** | 2 | Webhook、健康檢查 |
| **Smart Matching** | 7 | 人才搜尋、智慧媒合 |
| **Notifications** | 8 | 通知列表、已讀管理、佇列統計 |
| **Export** | 8 | CSV、PDF、ICS 匯出 |
| **Geo** | 1 | 城市列表 |
| **Hashtags** | 2 | 標籤管理 |
| **Test APIs** | 3 | 測試端點 |
| **Health** | 1 | 健康檢查 |
| **Total** | **171** | - |

> **提示**: 完整的 API 詳細文檔請查看 Swagger UI: `/swagger/index.html`

---

## 🔐 Authentication

所有需要認證的端點都需要在 Header 中帶入 JWT Token：

```http
Authorization: Bearer <your-jwt-token>
```

### Roles

- **OWNER**: 系統擁有者，完整權限
- **ADMIN**: 管理員，可管理所有中心
- **STAFF**: 一般管理人員，可管理指定中心
- **TEACHER**: 教師，只能存取自己的資料

---

## 📚 API Categories

### Authentication (4 APIs)
- `POST /api/v1/auth/admin/login` - 管理員登入
- `POST /api/v1/auth/teacher/line/login` - 教師 LINE 登入
- `POST /api/v1/auth/refresh` - 刷新 Token
- `POST /api/v1/auth/logout` - 登出

### Teacher Profile (8 APIs)
- `GET /api/v1/teacher/me/profile` - 取得個人資料
- `PUT /api/v1/teacher/me/profile` - 更新個人資料
- `GET /api/v1/teacher/me/centers` - 取得已加入的中心
- `GET /api/v1/teacher/me/skills` - 取得技能列表
- `POST /api/v1/teacher/me/skills` - 新增技能
- `PUT /api/v1/teacher/me/skills/:id` - 更新技能
- `DELETE /api/v1/teacher/me/skills/:id` - 刪除技能
- `POST /api/v1/teacher/me/certificates/upload` - 上傳證照

### Teacher Schedule (11 APIs)
- `GET /api/v1/teacher/me/schedule` - 取得綜合課表
- `GET /api/v1/teacher/schedules` - 取得課表（替代參數）
- `GET /api/v1/teacher/me/centers/:id/schedule-rules` - 取得中心排課規則
- `GET /api/v1/teacher/sessions/note` - 取得課堂筆記
- `PUT /api/v1/teacher/sessions/note` - 更新課堂筆記
- `POST /api/v1/teacher/scheduling/check-rule-lock` - 檢查規則鎖定
- `POST /api/v1/teacher/scheduling/preview-recurrence-edit` - 預覽循環編輯
- `POST /api/v1/teacher/scheduling/edit-recurring` - 編輯循環排課
- `POST /api/v1/teacher/scheduling/delete-recurring` - 刪除循環排課
- `GET /api/v1/teacher/me/schedule.ics` - 匯出 ICS
- `GET /api/v1/teacher/me/schedule/image` - 匯出圖片

### Teacher Personal Events (5 APIs)
- `GET /api/v1/teacher/me/personal-events` - 取得個人行程
- `POST /api/v1/teacher/me/personal-events` - 新增個人行程
- `PATCH /api/v1/teacher/me/personal-events/:id` - 更新個人行程
- `DELETE /api/v1/teacher/me/personal-events/:id` - 刪除個人行程
- `GET /api/v1/teacher/me/personal-events/:id/note` - 取得行程備註

### Teacher Exceptions (3 APIs)
- `GET /api/v1/teacher/exceptions` - 取得例外列表
- `POST /api/v1/teacher/exceptions` - 提出停課/改期申請
- `POST /api/v1/teacher/exceptions/:id/revoke` - 撤回申請

### Teacher Invitations (3 APIs)
- `GET /api/v1/teacher/me/invitations` - 取得邀請列表
- `POST /api/v1/teacher/me/invitations/respond` - 回應邀請
- `GET /api/v1/teacher/me/invitations/pending-count` - 取得待處理數量

### Admin Centers (2 APIs)
- `GET /api/v1/admin/centers` - 取得中心列表
- `POST /api/v1/admin/centers` - 新增中心

### Admin Rooms (5 APIs)
- `GET /api/v1/admin/rooms` - 取得教室列表
- `GET /api/v1/admin/rooms/active` - 取得啟用教室
- `POST /api/v1/admin/rooms` - 新增教室
- `PUT /api/v1/admin/rooms/:room_id` - 更新教室
- `PATCH /api/v1/admin/rooms/:room_id/toggle-active` - 切換啟用

### Admin Courses (6 APIs)
- `GET /api/v1/admin/courses` - 取得課程列表
- `GET /api/v1/admin/courses/active` - 取得啟用課程
- `POST /api/v1/admin/courses` - 新增課程
- `PUT /api/v1/admin/courses/:course_id` - 更新課程
- `DELETE /api/v1/admin/courses/:course_id` - 刪除課程
- `PATCH /api/v1/admin/courses/:course_id/toggle-active` - 切換啟用

### Admin Offerings (6 APIs)
- `GET /api/v1/admin/offerings` - 取得開課列表
- `GET /api/v1/admin/offerings/active` - 取得啟用開課
- `POST /api/v1/admin/offerings` - 新增開課
- `PUT /api/v1/admin/offerings/:offering_id` - 更新開課
- `DELETE /api/v1/admin/offerings/:offering_id` - 刪除開課
- `POST /api/v1/admin/centers/:id/offerings/:offering_id/copy` - 複製開課

### Admin Holidays (4 APIs)
- `GET /api/v1/admin/centers/:id/holidays` - 取得假日列表
- `POST /api/v1/admin/centers/:id/holidays` - 新增假日
- `POST /api/v1/admin/centers/:id/holidays/bulk` - 批次新增假日
- `DELETE /api/v1/admin/centers/:id/holidays/:holiday_id` - 刪除假日

### Admin Teachers (4 APIs)
- `GET /api/v1/admin/teachers` - 取得教師列表
- `GET /api/v1/admin/teachers/:teacher_id/note` - 取得教師備註
- `PUT /api/v1/admin/teachers/:teacher_id/note` - 更新教師備註
- `DELETE /api/v1/admin/teachers/:teacher_id/note` - 刪除教師備註

### Admin Invitations (5 APIs)
- `GET /api/v1/admin/centers/:id/invitations` - 取得邀請列表
- `GET /api/v1/admin/centers/:id/invitations/stats` - 取得邀請統計
- `POST /api/v1/admin/centers/:id/invitations` - 邀請教師
- `POST /api/v1/admin/centers/:id/invitations/generate-link` - 產生邀請連結
- `GET /api/v1/admin/centers/:id/invitations/links` - 取得邀請連結列表

### Admin Templates (8 APIs)
- `GET /api/v1/admin/templates` - 取得模板列表
- `POST /api/v1/admin/templates` - 新增模板
- `PUT /api/v1/admin/templates/:template_id` - 更新模板
- `DELETE /api/v1/admin/templates/:template_id` - 刪除模板
- `GET /api/v1/admin/templates/:template_id/cells` - 取得模板儲存格
- `POST /api/v1/admin/templates/:template_id/cells` - 新增模板儲存格
- `DELETE /api/v1/admin/templates/:template_id/cells/:cell_id` - 刪除儲存格
- `POST /api/v1/admin/templates/:template_id/apply` - 套用模板

### Scheduling Validation (4 APIs)
- `POST /api/v1/admin/scheduling/check-overlap` - 檢查重疊
- `POST /api/v1/admin/scheduling/check-teacher-buffer` - 檢查教師緩衝
- `POST /api/v1/admin/scheduling/check-room-buffer` - 檢查教室緩衝
- `POST /api/v1/admin/scheduling/validate` - 完整驗證

### Scheduling Rules (6 APIs)
- `GET /api/v1/admin/rules` - 取得規則列表
- `POST /api/v1/admin/rules` - 新增規則
- `PUT /api/v1/admin/rules/:ruleId` - 更新規則
- `DELETE /api/v1/admin/rules/:ruleId` - 刪除規則
- `POST /api/v1/admin/scheduling/check-rule-lock` - 檢查規則鎖定
- `GET /api/v1/admin/rules/:ruleId/exceptions` - 取得規則例外

### Scheduling Exceptions (6 APIs)
- `POST /api/v1/admin/scheduling/exceptions` - 新增例外
- `GET /api/v1/admin/exceptions` - 取得例外列表（日期範圍）
- `GET /api/v1/admin/exceptions/all` - 取得所有例外
- `GET /api/v1/admin/exceptions/pending` - 取得待審核例外
- `POST /api/v1/admin/scheduling/exceptions/:id/review` - 審核例外
- `POST /api/v1/admin/expand-rules` - 展開規則

### Smart Matching (7 APIs)
- `POST /api/v1/admin/smart-matching/matches` - 尋找替代教師
- `GET /api/v1/admin/smart-matching/talent/search` - 搜尋人才
- `GET /api/v1/admin/smart-matching/talent/stats` - 取得人才統計
- `POST /api/v1/admin/smart-matching/talent/invite` - 邀請人才
- `GET /api/v1/admin/smart-matching/suggestions` - 取得搜尋建議
- `POST /api/v1/admin/smart-matching/alternatives` - 取得替代時段
- `GET /api/v1/admin/teachers/:teacher_id/sessions` - 取得教師課程

### Notifications (6 APIs)
- `GET /api/v1/notifications` - 取得通知列表
- `GET /api/v1/notifications/unread-count` - 取得未讀數量
- `POST /api/v1/notifications/:id/read` - 標記已讀
- `POST /api/v1/notifications/read-all` - 全部標記已讀
- `POST /api/v1/notifications/token` - 設定通知 Token
- `POST /api/v1/notifications/test` - 發送測試通知

### Export (8 APIs)
- `POST /api/v1/admin/export/schedule/csv` - 匯出課表 CSV
- `POST /api/v1/admin/export/schedule/pdf` - 匯出課表 PDF
- `GET /api/v1/admin/centers/:id/export/teachers/csv` - 匯出教師 CSV
- `GET /api/v1/admin/centers/:id/export/exceptions/csv` - 匯出例外 CSV
- `GET /api/v1/teacher/me/schedule.ics` - 匯出 ICS
- `POST /api/v1/teacher/me/schedule/subscription` - 建立日曆訂閱
- `DELETE /api/v1/teacher/me/schedule/subscription` - 取消訂閱
- `GET /api/v1/teacher/me/schedule/image` - 匯出圖片

### Public APIs (2 APIs)
- `GET /api/v1/invitations/:token` - 取得公開邀請
- `POST /api/v1/invitations/:token/accept` - 接受邀請

### LINE Bot (2 APIs)
- `POST /api/v1/line/webhook` - LINE Webhook
- `GET /api/v1/line/health` - 健康檢查

### Geo - Cities (1 API)
- `GET /api/v1/geo/cities` - 取得城市列表（公開 API，無需認證）

### Hashtags (2 APIs)
- `GET /api/v1/hashtags/search` - 搜尋標籤建議
- `POST /api/v1/hashtags` - 建立新標籤

### Teacher Certificates (4 APIs)
- `GET /api/v1/teacher/me/certificates` - 取得證照列表
- `POST /api/v1/teacher/me/certificates` - 新增證照
- `PUT /api/v1/teacher/me/certificates/:id` - 更新證照
- `DELETE /api/v1/teacher/me/certificates/:id` - 刪除證照

### Teacher Personal Events Notes (2 APIs)
- `GET /api/v1/teacher/me/personal-events/:id/note` - 取得行程備註
- `PUT /api/v1/teacher/me/personal-events/:id/note` - 更新行程備註

### Teacher Background Images (3 APIs)
- `POST /api/v1/teacher/me/backgrounds` - 上傳背景圖片
- `GET /api/v1/teacher/me/backgrounds` - 取得背景圖片列表
- `DELETE /api/v1/teacher/me/backgrounds/:id` - 刪除背景圖片

### Teacher Public Register (1 API)
- `POST /api/v1/teacher/public/register` - 公開註冊（LINE Bot 使用）

### Admin Teachers Management (6 APIs)
- `GET /api/v1/teachers` - 取得教師列表
- `POST /api/v1/admin/teachers/placeholder` - 建立佔位老師
- `POST /api/v1/admin/teachers/merge` - 合併老師
- `DELETE /api/v1/teachers/:id` - 刪除老師
- `DELETE /api/v1/admin/centers/:id/teachers/:teacher_id` - 從中心移除老師

### Admin Terms (7 APIs)
- `GET /api/v1/admin/terms` - 取得學期列表
- `GET /api/v1/admin/terms/active` - 取得啟用學期
- `POST /api/v1/admin/terms` - 新增學期
- `PUT /api/v1/admin/terms/:term_id` - 更新學期
- `DELETE /api/v1/admin/terms/:term_id` - 刪除學期
- `GET /api/v1/admin/occupancy/rules` - 取得佔位規則
- `POST /api/v1/admin/terms/copy-rules` - 複製規則

### Admin LINE Binding (6 APIs)
- `GET /api/v1/admin/me/line-binding` - 取得 LINE 綁定狀態
- `POST /api/v1/admin/me/line/bind` - 初始化 LINE 綁定
- `DELETE /api/v1/admin/me/line/unbind` - 解除 LINE 綁定
- `GET /api/v1/admin/me/line/notify-settings` - 取得通知設定
- `PATCH /api/v1/admin/me/line/notify-settings` - 更新通知設定
- `GET /api/v1/admin/me/line/qrcode` - 取得 QR Code

### Admin LINE QR Code (1 API)
- `GET /api/v1/admin/me/line/qrcode-with-code` - 取得帶驗證碼的 QR Code

### Admin Profile (2 APIs)
- `GET /api/v1/admin/me/profile` - 取得管理員個人資料
- `POST /api/v1/admin/me/change-password` - 變更密碼

### Admin Management (5 APIs)
- `GET /api/v1/admin/admins` - 取得管理員列表
- `POST /api/v1/admin/admins` - 新增管理員
- `POST /api/v1/admin/admins/toggle-status` - 切換管理員狀態
- `POST /api/v1/admin/admins/reset-password` - 重設管理員密碼
- `POST /api/v1/admin/admins/change-role` - 變更管理員角色

### Admin Center Settings (2 APIs)
- `GET /api/v1/admin/centers/:id/settings` - 取得中心設定
- `PATCH /api/v1/admin/centers/:id/settings` - 更新中心設定

### Admin Invitations Links (4 APIs)
- `GET /api/v1/admin/centers/:id/invitations/general-link` - 取得一般邀請連結
- `POST /api/v1/admin/centers/:id/invitations/general-link` - 建立一般邀請連結
- `POST /api/v1/admin/centers/:id/invitations/general-link/toggle` - 切換一般邀請狀態
- `POST /api/v1/admin/centers/:id/invitations/general-link/regenerate` - 重新產生一般邀請連結
- `DELETE /api/v1/admin/invitations/links/:id` - 撤銷邀請連結

### Scheduling Dashboard (2 APIs)
- `GET /api/v1/admin/dashboard/today-summary` - 取得今日摘要
- `GET /api/v1/admin/scheduling/matrix-view` - 取得矩陣視圖

### Scheduling Phase Detection (1 API)
- `POST /api/v1/admin/detect-phase-transitions` - 偵測階段轉換

### Admin Notifications Queue (2 APIs)
- `GET /api/v1/admin/notifications/queue-stats` - 取得通知佇列統計
- `POST /api/v1/admin/notifications/broadcast` - 廣播通知

### Admin Resource (1 API)
- `GET /api/v1/admin/teachers` - 取得教師資源列表

### Templates Additional (2 APIs)
- `PUT /api/v1/admin/templates/:templateId/cells/reorder` - 重新排序儲存格
- `POST /api/v1/admin/templates/:templateId/validate-apply` - 驗證並套用模板
- `DELETE /api/v1/admin/templates/cells/:cellId` - 刪除儲存格

### Health Check (1 API)
- `GET /api/v1/admin/health/redis` - Redis 健康檢查

### Calendar Subscription (1 API)
- `GET /api/v1/calendar/subscribe/:token.ics` - 訂閱日曆（公開 API）

### Test APIs (3 APIs)
- `GET /api/test/r2-status` - R2 儲存狀態
- `POST /api/test/upload` - 測試上傳
- `POST /api/test/upload-batch` - 測試批次上傳

---

## 📄 Response Format

所有 API 統一使用以下回傳格式：

```json
{
  "code": 0,
  "message": "success",
  "datas": <actual_data>
}
```

- `code`: 0 表示成功，非 0 表示錯誤
- `message`: 訊息描述
- `datas`: 實際資料（可能是物件、陣列或 null）

---

## ❌ Error Codes

| Code | Description |
|:---|:---|
| 0 | Success |
| 40001 | Bad Request |
| 40101 | Unauthorized |
| 40301 | Forbidden |
| 40401 | Not Found |
| 40901 | Conflict |
| 50001 | Internal Server Error |

---

## 📝 Version History

| Version | Date | Changes |
|:---|:---|:---|
| 1.1.2 | 2026-02-13 | 修正 API 數量為 171（原為 186）|
| 1.1.1 | 2026-02-13 | 修正 API 數量為 176 |
| 1.1.0 | 2026-02-13 | 新增 Admin Terms、Admin LINE Binding、Admin Profile、Admin Management、Scheduling Dashboard、Geo、Hashtags、Test APIs 等端點 |
| 1.0.0 | 2026-01-31 | 初始版本 |

---

**Last Updated**: 2026-02-13
**Version**: 1.1.2
**Swagger UI**: `/swagger/index.html`
**Swagger JSON**: `/swagger/doc.json`
