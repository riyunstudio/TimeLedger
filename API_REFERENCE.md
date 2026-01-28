# API 端點完整對照表

**最後更新**：2026年1月28日  
**來源**：swagger.json (74 endpoints, 106 models)

---

## 回應格式標準

### 成功響應

```json
{
  "code": 0,
  "message": "Success",
  "datas": { ... }
}
```

### 分頁響應

```json
{
  "code": 0,
  "message": "Success",
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "total_pages": 8,
    "has_next": true,
    "has_prev": false
  }
}
```

### 錯誤響應

```json
{
  "code": 10001,
  "message": "System error",
  "datas": null
}
```

### Rate Limit 響應

```json
{
  "code": 10009,
  "message": "請求頻率過高，請稍後再試",
  "datas": {
    "retry_after": 30
  }
}
```

**Rate Limit Headers**：
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 99
X-RateLimit-Reset: 2026-01-28T10:00:00+08:00
```

---

## 一、管理員相關 API

### 1.1 Admin - LINE 綁定

| Method | Endpoint | 說明 |
|:---|:---|:---|
| GET | `/admin/me/line-binding` | 取得 LINE 綁定狀態 |
| POST | `/admin/me/line/bind` | 初始化 LINE 綁定（產生驗證碼） |
| DELETE | `/admin/me/line/unbind` | 解除 LINE 綁定 |
| PATCH | `/admin/me/line/notify-settings` | 更新 LINE 通知設定 |

**LINEBindingResponse**：
```json
{
  "is_bound": false,
  "bound_at": null,
  "notify_enabled": true
}
```

**UpdateNotifySettingsRequest**：
```json
{
  "notify_enabled": true
}
```

### 1.2 Admin - 中心管理

| Method | Endpoint | 說明 |
|:---|:---|:---|
| GET | `/api/v1/admin/centers` | 取得中心列表 |
| POST | `/api/v1/admin/centers` | 建立新中心 |

**CreateCenterRequest**：
```json
{
  "name": "中心名稱",
  "plan_level": "STARTER",
  "settings": {
    "allow_public_register": true,
    "default_language": "zh-TW",
    "exception_lead_days": 14
  }
}
```

### 1.3 Admin - 教室管理

| Method | Endpoint | 說明 |
|:---|:---|:---|
| GET | `/api/v1/admin/centers/:id/rooms` | 取得教室列表 |
| POST | `/api/v1/admin/centers/:id/rooms` | 建立教室 |
| PUT | `/api/v1/admin/rooms/:id` | 更新教室 |
| DELETE | `/api/v1/admin/rooms/:id` | 刪除教室 |

**CreateRoomRequest / UpdateRoomRequest**：
```json
{
  "name": "瑜伽教室A",
  "capacity": 20,
  "is_active": true
}
```

### 1.4 Admin - 課程管理

| Method | Endpoint | 說明 |
|:---|:---|:---|
| GET | `/api/v1/admin/courses` | 取得課程列表 |
| POST | `/api/v1/admin/courses` | 建立課程 |
| PUT | `/api/v1/admin/courses/{course_id}` | 更新課程 |

**CreateCourseRequest**：
```json
{
  "name": "哈達瑜伽",
  "default_duration": 60,
  "color_hex": "#FF6B6B",
  "room_buffer_min": 10,
  "teacher_buffer_min": 10,
  "is_active": true
}
```

### 1.5 Admin - 開課管理 (Offerings)

| Method | Endpoint | 說明 |
|:---|:---|:---|
| GET | `/api/v1/admin/centers/:id/offerings` | 取得開課列表 |
| POST | `/api/v1/admin/centers/:id/offerings` | 建立開課 |
| PUT | `/api/v1/admin/offerings/:id` | 更新開課 |
| POST | `/api/v1/admin/centers/:id/offerings/:offering_id/copy` | 複製班別 |

**CreateOfferingRequest**：
```json
{
  "course_id": 1,
  "name": "週一早班哈達瑜伽",
  "default_room_id": 1,
  "default_teacher_id": null,
  "allow_buffer_override": false,
  "is_active": true
}
```

**CopyOfferingRequest**：
```json
{
  "target_weekday": 2,
  "start_date": "2026-02-01",
  "end_date": "2026-02-28"
}
```

### 1.6 Admin - 假日管理

| Method | Endpoint | 說明 |
|:---|:---|:---|
| GET | `/api/v1/admin/centers/:id/holidays` | 取得假日列表 |
| POST | `/api/v1/admin/centers/:id/holidays/bulk` | 批次建立假日 |

**BulkCreateHolidaysRequest**：
```json
{
  "holidays": [
    {
      "name": "農曆春節",
      "date": "2026-01-29"
    }
  ]
}
```

### 1.7 Admin - 範本管理 (Timetable Templates)

| Method | Endpoint | 說明 |
|:---|:---|:---|
| GET | `/api/v1/admin/centers/:id/templates` | 取得範本列表 |
| POST | `/api/v1/admin/centers/:id/templates` | 建立範本 |
| PUT | `/api/v1/admin/templates/:id` | 更新範本 |
| DELETE | `/api/v1/admin/templates/:id` | 刪除範本 |

---

## 二、排課與例外 API

### 2.1 Admin - 排課驗證

| Method | Endpoint | 說明 |
|:---|:---|:---|
| POST | `/api/v1/admin/centers/:id/validate` | 驗證排課衝突 |

**ValidationRequest**：
```json
{
  "rule_id": 1,
  "start_date": "2026-02-01",
  "end_date": "2026-02-28",
  "override": false
}
```

**ValidationResponse**：
```json
{
  "valid": true,
  "conflicts": []
}
```

### 2.2 Admin - 排課展開

| Method | Endpoint | 說明 |
|:---|:---|:---|
| POST | `/api/v1/admin/centers/:id/expand` | 展開排課規則 |

**ExpandScheduleRequest**：
```json
{
  "rule_ids": [1, 2, 3],
  "start_date": "2026-02-01",
  "end_date": "2026-02-28"
}
```

### 2.3 Admin - 例外管理

| Method | Endpoint | 說明 |
|:---|:---|:---|
| GET | `/api/v1/admin/centers/:id/exceptions` | 取得例外列表 |
| POST | `/api/v1/admin/centers/:id/exceptions` | 建立例外申請 |
| POST | `/api/v1/admin/centers/:id/exceptions/:exceptionId/review` | 審核例外 |

**ReviewExceptionRequest**：
```json
{
  "status": "APPROVED",
  "review_note": "已核准"
}
```

### 2.4 Teacher - 例外申請

| Method | Endpoint | 說明 |
|:---|:---|:---|
| GET | `/api/v1/teacher/exceptions` | 取得例外列表 |
| POST | `/api/v1/teacher/exceptions` | 申請停課/改期 |
| POST | `/api/v1/teacher/exceptions/:id/revoke` | 撤回申請 |

**TeacherCreateExceptionRequest**：
```json
{
  "center_id": 1,
  "rule_id": 1,
  "original_date": "2026-02-03",
  "exception_type": "CANCEL",
  "reason": "身體不適需要休息"
}
```

---

## 三、智慧媒合 API

### 3.1 Smart Matching

| Method | Endpoint | 說明 |
|:---|:---|:---|
| POST | `/api/v1/admin/smart-matching/matches` | 智慧媒合搜尋 |
| GET | `/api/v1/admin/smart-matching/suggestions` | 搜尋建議 |
| POST | `/api/v1/admin/smart-matching/alternatives` | 替代時段建議 |
| GET | `/api/v1/admin/teachers/:id/sessions` | 教師課表查詢 |

**MatchingRequest**：
```json
{
  "center_id": 1,
  "offering_id": 1,
  "preferred_teacher_ids": [1, 2],
  "exclude_teacher_ids": [],
  "skills": ["瑜伽"],
  "city": "臺北市",
  "min_rating": 4.0
}
```

### 3.2 Talent Pool

| Method | Endpoint | 說明 |
|:---|:---|:---|
| GET | `/api/v1/admin/smart-matching/talent/search` | 人才庫搜尋 |
| GET | `/api/v1/admin/smart-matching/talent/stats` | 人才庫統計 |
| POST | `/api/v1/admin/smart-matching/talent/invite` | 邀請人才 |

### 3.3 Admin - 邀請管理

| Method | Endpoint | 說明 |
|:---|:---|:---|
| POST | `/api/v1/admin/centers/:id/invitations` | 邀請老師加入中心 |

**InviteTeacherRequest**：
```json
{
  "teacher_ids": [1, 2, 3],
  "message": "誠摯邀請您加入我們的教學團隊"
}
```

---

## 四、教師相關 API

### 4.1 Teacher - 個人資料

| Method | Endpoint | 說明 |
|:---|:---|:---|
| GET | `/api/v1/teacher/me/profile` | 取得個人資料 |
| PUT | `/api/v1/teacher/me/profile` | 更新個人資料 |

**UpdateTeacherProfileRequest**：
```json
{
  "bio": "資深瑜伽老師",
  "city": "臺北市",
  "district": "大安區",
  "public_contact_info": null,
  "is_open_to_hiring": true
}
```

### 4.2 Teacher - 所屬中心

| Method | Endpoint | 說明 |
|:---|:---|:---|
| GET | `/api/v1/teacher/me/centers` | 取得所屬中心列表 |

### 4.3 Teacher - 課表

| Method | Endpoint | 說明 |
|:---|:---|:---|
| GET | `/api/v1/teacher/me/schedule` | 取得課表 |

**查詢參數**：
| 參數 | 類型 | 說明 |
|:---|:---|:---|
| `week_start` | string | 週開始日期 (YYYY-MM-DD) |
| `week_end` | string | 週結束日期 (YYYY-MM-DD) |

**TeacherScheduleItem 結構**：
```json
{
  "id": "1_2026-01-06",
  "type": "CENTER_SESSION",
  "title": "週一上午哈達瑜伽",
  "date": "2026-01-06",
  "start_time": "09:00",
  "end_time": "10:00",
  "duration": 60,
  "room_id": 1,
  "room_name": "瑜伽教室A",
  "teacher_id": 1,
  "center_id": 1,
  "center_name": "TimeLedger 旗艦館",
  "color": "#FF6B6B",
  "status": "NORMAL"
}
```

### 4.4 Teacher - 今日摘要

| Method | Endpoint | 說明 |
|:---|:---|:---|
| GET | `/api/v1/teacher/today/summary` | 取得今日摘要 |

**TodaySummaryResponse**：
```json
{
  "date": "2026-01-28",
  "weekday": "週三",
  "total_sessions": 3,
  "total_hours": 4.5,
  "sessions": [...],
  "offerings": [...],
  "teachers": [...],
  "rooms": [...]
}
```

---

## 五、課堂筆記 API

### 5.1 Teacher - 課堂筆記

| Method | Endpoint | 說明 |
|:---|:---|:---|
| GET | `/api/v1/teacher/sessions/note` | 取得課堂筆記 |
| PUT | `/api/v1/teacher/sessions/note` | 新增/更新筆記 |

**UpsertSessionNoteRequest**：
```json
{
  "rule_id": 1,
  "session_date": "2026-01-25",
  "content": "本課程適合初學者",
  "prep_note": "請提前 15 分鐘到場"
}
```

### 5.2 Admin - 教師筆記

| Method | Endpoint | 說明 |
|:---|:---|:---|
| GET | `/api/v1/admin/teachers/:id/notes` | 取得教師筆記 |
| POST | `/api/v1/admin/teachers/:id/notes` | 新增筆記 |
| PUT | `/api/v1/admin/notes/:id` | 更新筆記 |

**UpsertTeacherNoteRequest**：
```json
{
  "note_type": "GENERAL",
  "content": "教學經驗豐富",
  "rating": 5
}
```

---

## 六、認證 API

### 6.1 Admin 登入

| Method | Endpoint | 說明 |
|:---|:---|:---|
| POST | `/api/v1/auth/admin/login` | 管理員登入 |

**AdminLoginRequest**：
```json
{
  "email": "admin@timeledger.com",
  "password": "admin123"
}
```

**LoginResponse**：
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "admin@timeledger.com",
    "name": "超級管理員",
    "role": "OWNER",
    "center_id": 1
  }
}
```

### 6.2 Teacher LINE 登入

| Method | Endpoint | 說明 |
|:---|:---|:---|
| POST | `/api/v1/auth/teacher/line/login` | 教師 LINE 登入 |

**TeacherLineLoginRequest**：
```json
{
  "line_user_id": "U000000000000000000000001",
  "id_token": "eyJraWQiOiJ..."
}
```

### 6.3 Token 刷新

| Method | Endpoint | 說明 |
|:---|:---|:---|
| POST | `/api/v1/auth/refresh` | 刷新 Token |
| POST | `/api/v1/auth/logout` | 登出 |

---

## 七、通知 API

### 7.1 通知管理

| Method | Endpoint | 說明 |
|:---|:---|:---|
| GET | `/api/v1/notifications` | 取得通知列表 |
| GET | `/api/v1/notifications/unread-count` | 取得未讀數 |
| POST | `/api/v1/notifications/:id/read` | 標記已讀 |
| POST | `/api/v1/notifications/read-all` | 全部標記已讀 |

### 7.2 佇列監控

| Method | Endpoint | 說明 |
|:---|:---|:---|
| GET | `/api/v1/admin/notifications/queue-stats` | 通知佇列統計 |

**QueueStatsResponse**：
```json
{
  "pending_count": 15,
  "retry_count": 3,
  "completed_count": 1250,
  "failed_count": 12,
  "failure_rate": 0.95,
  "redis_connected": true,
  "worker_running": true
}
```

---

## 八、匯出 API

### 8.1 Admin 匯出

| Method | Endpoint | 說明 |
|:---|:---|:---|
| POST | `/api/v1/admin/export/schedule/csv` | 匯出課表 CSV |
| POST | `/api/v1/admin/export/schedule/pdf` | 匯出課表 PDF |
| GET | `/api/v1/admin/export/centers/:centerId/teachers/csv` | 匯出教師 CSV |
| GET | `/api/v1/admin/export/centers/:centerId/exceptions/csv` | 匯出例外 CSV |

**ExportScheduleRequest**：
```json
{
  "center_id": 1,
  "start_date": "2026-02-01",
  "end_date": "2026-02-28",
  "include_notes": true
}
```

---

## 九、技能與證書 API

### 9.1 Teacher Skills

| Method | Endpoint | 說明 |
|:---|:---|:---|
| GET | `/api/v1/teacher/me/skills` | 取得技能列表 |
| POST | `/api/v1/teacher/me/skills` | 新增技能 |
| PUT | `/api/v1/teacher/me/skills/:id` | 更新技能 |
| DELETE | `/api/v1/teacher/me/skills/:id` | 刪除技能 |

**CreateSkillRequest / UpdateSkillRequest**：
```json
{
  "category": "瑜伽",
  "skill_name": "哈達瑜伽",
  "level": "EXPERT"
}
```

### 9.2 Teacher Certificates

| Method | Endpoint | 說明 |
|:---|:---|:---|
| GET | `/api/v1/teacher/me/certificates` | 取得證書列表 |
| POST | `/api/v1/teacher/me/certificates` | 上傳證書 |
| DELETE | `/api/v1/teacher/me/certificates/:id` | 刪除證書 |

**CreateCertificateRequest**：
```json
{
  "name": "RYT-500 瑜伽教練證書",
  "issuing_organization": "Yoga Alliance",
  "issue_date": "2020-01-15",
  "expiry_date": null,
  "credential_url": "https://..."
}
```

---

## 十、地理資料 API

### 10.1 縣市與區域

| Method | Endpoint | 說明 |
|:---|:---|:---|
| GET | `/api/v1/geo/cities` | 取得縣市列表 |
| GET | `/api/v1/geo/cities/:id/districts` | 取得區域列表 |

---

## 十一、錯誤碼參照

### 系統錯誤 (1xxxx)

| 錯誤碼 | 訊息 | 說明 |
|:---|:---|:---|
| 10001 | System error | 系統錯誤 |
| 10002 | Invalid request | 無效請求 |
| 10003 | Unauthorized | 未授權 |
| 10004 | Forbidden | 禁止訪問 |
| 10005 | Not found | 資源不存在 |
| 10006 | Conflict | 資源衝突 |
| 10007 | Validation error | 驗證錯誤 |
| 10008 | Rate limit exceeded | 超過速率限制 |
| 10009 | Too many requests | 請求過於頻繁 |

### 資料庫錯誤 (2xxxx)

| 錯誤碼 | 訊息 | 說明 |
|:---|:---|:---|
| 20001 | Database error | 資料庫錯誤 |
| 20002 | Duplicate entry | 重複條目 |
| 20003 | Foreign key constraint | 外鍵約束錯誤 |

### 業務錯誤 (4xxxx)

| 錯誤碼 | 訊息 | 說明 |
|:---|:---|:---|
| 40001 | Center not found | 中心不存在 |
| 40002 | Teacher not found | 教師不存在 |
| 40003 | Schedule conflict | 排課衝突 |
| 40004 | Invalid time range | 無效時間範圍 |
| 40005 | Exception not found | 例外申請不存在 |
| 40006 | Permission denied | 權限不足 |
| 40007 | Token expired | Token 過期 |
| 40008 | Invalid token | 無效 Token |

---

## 十二、通用查詢參數

| 參數 | 類型 | 預設值 | 說明 |
|:---|:---|:---:|:---|
| `page` | int | 1 | 頁碼 |
| `limit` | int | 20 | 每頁筆數 (最大 100) |
| `sort_by` | string | 依各 API 定義 | 排序欄位 |
| `sort_order` | string | ASC | 排序方向 (ASC/DESC) |
| `keyword` | string | - | 關鍵字搜尋 |
| `status` | string | - | 狀態篩選 |
| `start_date` | string | - | 開始日期 |
| `end_date` | string | - | 結束日期 |

---

## 十三、認證與授權

### JWT Token

所有 API（除登入相關）都需要在 Header 中攜帶 JWT Token：

```
Authorization: Bearer <token>
```

### Token 結構

```json
{
  "user_id": 1,
  "user_type": "ADMIN" | "TEACHER",
  "center_id": 1,
  "role": "OWNER" | "ADMIN" | "STAFF",
  "exp": 1738000000,
  "iat": 1737913600
}
```

### Mock Token（測試用）

開發環境支援 `mock-` 前綴的 Token：

| Header | 用途 |
|:---|:---|
| `Authorization: Bearer mock-admin-token` | 管理員測試 |
| `Authorization: Bearer mock-teacher-token` | 教師測試 |

---

## 十四、測試覆蓋

### 測試檔案位置

- `testing/test/` - 後端單元/整合測試
- `frontend/tests/` - 前端 E2E 測試

### 執行測試

```bash
# 後端測試
go test ./testing/test/... -v

# 單元測試
go test ./testing/test/... -run TestCrossDay

# E2E 測試
cd frontend && npm test
```

---

## 十五、Swagger 文件

本文件同步自 `docs/swagger.json`

```bash
# 重新產生 Swagger 文件
swag init -g main.go -o docs --parseDependency
```

---

**文件版本**：2026-01-28  
**維護者**：TimeLedger Dev Team
