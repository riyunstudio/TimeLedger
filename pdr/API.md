# TimeLedger API 契約 (API Specification)

## 1. 共通規範 (Conventions)
- **Base URL**: `/api/v1`
- **Response Format**:
  ```json
  {
    "code": "SUCCESS",   // or ERR_OVERLAP, ERR_AUTH...
    "message": "Operation successful",
    "data": { ... }      // Payload
  }
  ```
- **Authentication**:
  - `Authorization: Bearer <JWT>`
  - Header 區分角色權限 (Admin vs Teacher)。

---

## 2. 認證模組 (Auth)

### `GET /auth/line/login`
- **描述**: 取得 LINE Login 授權跳轉網址。
- **Response**: `{ "url": "https://access.line.me/..." }`

### `POST /auth/line/callback`
- **描述**: 用 Authorization Code 交換 Token 並登入/註冊。
- **Body**: `{ "code": "...", "redirect_uri": "..." }`
- **Response**: `{ "token": "JWT...", "user": { "id": 1, "name": "Teacher A", "line_id": "..." } }`

### `POST /auth/admin/login`
- **描述**: 中心管理員登入。
- **Body**: `{ "email": "...", "password": "..." }`
- **Response**: `{ "token": "JWT...", "user": { "role": "ADMIN", "center_id": 10 } }`

---

## 3. 老師 / 個人端 API (Teacher & Personal)

### `GET /teacher/me/profile`
- **描述**: 取得個人資料 (Bio, Skills, Region, Hashtags)。
- **Response**:
  ```json
  {
      "name": "Teacher A",
      "bio": "...",
      "region": "台北市大安區",
      "skills": [{ "category": "舞蹈", "name": "街舞" }],
      "hashtags": ["#MV舞蹈"],
      "is_open_to_hiring": true
  }
  ```

### `POST /teacher/me/profile`
- **描述**: 更新個人專業檔案。
- **Body**:
  ```json
  {
      "bio": "...",
      "skills": [
          { "category": "舞蹈", "name": "街舞", "level": "ADVANCED" },
          { "category": "瑜伽", "name": "空中瑜伽", "level": "INTERMEDIATE" }
      ],
      "hashtags": ["#小班制", "#MV舞蹈", "#全英文教學"],
      "is_open_to_hiring": true,
      "region": "台北市大安區",
      "public_contact_info": "LineID: happy_teacher"
  }
  ```

### `GET /teacher/me/centers`
- **描述**: 取得已加入的中心列表。
- **Response**: `[ { "id": 1, "name": "Mozart Center", "status": "ACTIVE" } ]`

### `POST /teacher/join`
- **描述**: 輸入邀請碼加入中心。
- **Body**: `{ "invite_code": "XYZ-123" }`

### `POST /teacher/me/certificates`
- **描述**: 上傳專業證照/獎狀。
- **Body** (Multipart): `file` (Image/PDF), `name`, `issued_at`
- **Response**: `{ "id": 1, "url": "..." }`

### `DELETE /teacher/me/certificates/{id}`
- **描述**: 刪除證照。

### `GET /teacher/me/schedule`
- **描述**: 取得指定範圍內的所有行程 (含個人行程 + 各中心課程)。
- **Query**: `from=2026-01-01&to=2026-01-31`
- **Response**:
  ```json
  [
    {
      "id": "evt_123",
      "type": "CENTER_SESSION", // or PERSONAL_EVENT
      "title": "鋼琴課 - 小明",
      "start": "2026-01-19T10:00:00",
      "end": "2026-01-19T11:00:00",
      "center_name": "莫札特音樂教室",
      "color": "#FF5733",
      "status": "NORMAL" // or PENDING_CANCEL, PENDING_RESCHEDULE
    }
  ]
  ```

### `GET /teacher/sessions/note`
- **描述**: 取得特定課程場次的筆記。
- **Query**: `rule_id=55&date=2026-01-20`
- **Response**: `{ "content": "上次教到...", "prep_note": "今天要驗收..." }`

### `PUT /teacher/sessions/note`
- **描述**: 儲存/更新課程筆記。
- **Body**:
  ```json
  {
      "rule_id": 55,
      "date": "2026-01-20",
      "content": "教完第一章",
      "prep_note": "記得帶講義"
  }
  ```

### `POST /teacher/personal-events`
- **描述**: 新增個人行程。
- **Body**:
  ```json
  {
    "title": "私人練琴",
    "start_at": "2026-01-20T14:00:00",
    "end_at": "2026-01-20T16:00:00",
    "recurrence": {
      "type": "WEEKLY",
      "interval": 1,
      "weekdays": [1, 3],
      "until": "2026-06-30" // or "count": 10
    },
    "color": "#AABBCC",  // UX優化：讓用戶自訂顏色
    "reminders": [15, 60] // UX優化：App推播提醒 (分)
  }
  ```

### `PATCH /teacher/personal-events/{id}`
- **描述**: 修改行程 (支援拖曳改時、內容修訂)。
- **Body**:
  ```json
  {
    "title": "私人練琴 (改)",
    "start_at": "2026-01-21T14:00:00",
    "end_at": "2026-01-21T16:00:00",
    "update_mode": "SINGLE" // 重要：針對循環行程的更新模式
                            // SINGLE: 僅此單一場次 (變成例外)
                            // FUTURE: 此場次及之後所有
                            // ALL: 修改整串循環規則
  }
  ```

### `DELETE /teacher/personal-events/{id}`
- **描述**: 刪除行程。
- **Query**: `mode=SINGLE|FUTURE|ALL` (若為循環行程必填)

### `GET /teacher/me/schedule/conflicts`
- **描述**: (UX優化) 在新增/拖曳行程時，主動檢查是否與「中心課程」或其他「私人行程」重疊。
- **Query**: `start_at=...&end_at=...&exclude_id=...`
- **Response**: `{ "has_conflict": true, "conflicts": [...] }`

### `POST /teacher/exceptions`
- **描述**: 針對 **中心課程** 提出異動申請 (請假/改期)。
- **Body**:
  ```json
  {
    "rule_id": 55,
    "original_date": "2026-01-22",
    "type": "RESCHEDULE", // CANCEL
    "new_start_at": "2026-01-23T10:00:00", // Required if RESCHEDULE
    "new_end_at": "2026-01-23T11:00:00",
    "reason": "家裡有事"
  }
  ```
- **Response**: `200 OK` (若 Policy 允許自動通過) 或 `202 Accepted` (若需審核, status=PENDING)。

---

## 4. 中心後台 API (Center Admin)

### `POST /admin/centers/{id}/validate`
- **描述**: 預檢排課衝突 (Validate Engine)。
- **Body**: `{ "teacher_id": 1, "room_id": 2, "start_at": "...", "end_at": "...", "override_buffer": false }`
- **Response**:
  ```json
  {
    "valid": false,
    "conflicts": [
      { "type": "TEACHER_BUFFER", "message": "老師上一堂課間隔不足 10 分鐘", "can_override": true }
    ]
  }
  ```

### `GET /admin/centers/{id}/exceptions`
- **描述**: 取得待審核列表。
- **Query**: `status=PENDING`

### `POST /admin/centers/{id}/exceptions/{exception_id}/review`
- **描述**: 審核異動申請。
- **Body**: `{ "action": "APPROVE", "override_buffer": true }` (或 REJECT)
- **Note**: Approve 時後端會執行 Force Validate。

### `POST /admin/centers/{id}/rules/bulk`
- **描述**: 批量新增/更新排課規則 (Grid 儲存)。
- **Body**: `[ { "teacher_id":..., "weekday":..., "time":... } ]`

### `GET /admin/centers/{id}/matching/teachers`
- **描述**: **(亮點功能)** 智慧媒合：根據時段與技能尋找可用老師。
- **Query**:
  - `start_at`: 2026-01-20T10:00:00
  - `end_at`: 2026-01-20T12:00:00
  - `skills`: Piano,Jazz (Comma separated)
- **Response**:
  ```json
  [
    {
      "teacher_id": 101,
      "name": "Teacher A",
      "match_score": 90,
      "region": "台北市大安區",
      "skills": [{ "category": "有氧", "name": "Zumba" }],
      "hashtags": ["#熱情", "#週末可"],
      "availability": "AVAILABLE", // or CONFLICT, BUFFER_CONFLICT
      "internal_note": "適合教初學者",
      "internal_rating": 5,
      "is_member": true, // false = 公開媒合到的非本中心老師
      "public_contact_info": "LineID: ..." // 僅當 (is_member=false AND open_hiring=true) 時有值
    }
  ]
  ```

  ```

### `GET /admin/talent/search`
- **描述**: **(人才庫搜尋)** 主動搜尋已開啟「求職媒合」的老師。
- **Query**:
  - `category`: 有氧 (Filter)
  - `skill`: Zumba (Filter)
  - `tag`: "#熱情" (Filter)
  - `region`: 台北市 (Filter)
  - `keyword`: "王小明" (Search Name/Bio)
- **Response**:
  ```json
  [
    {
      "teacher_id": 105,
      "name": "Teacher Bloom",
      "bio": "Expert Yoga Instructor...",
      "region": "台北市大安區",
      "skills": [{ "category": "瑜伽", "name": "空中瑜伽" }],
      "hashtags": ["#安靜", "#平日可"],
      "certificates": [{ "name": "RYT-200", "url": "..." }],
      "public_contact_info": "Line: bloom_yoga"
    }
  ]
  ```

### `GET /admin/centers/{id}/teachers/{teacher_id}/certificates`
- **描述**: 查看特定老師的證照列表 (僅限已加入該中心的老師)。
- **Response**: `[ { "name": "Yamaha Rank 5", "url": "...", "issued_at": "..." } ]`

### `GET /admin/centers/{id}/teachers/{teacher_id}/internal-note`
- **描述**: 取得中心對老師的私密備註與評分。
- **Response**: `{ "rating": 5, "internal_note": "..." }`

### `PUT /admin/centers/{id}/teachers/{teacher_id}/internal-note`
- **描述**: 更新中心對老師的私密備註與評分。
- **Body**: `{ "rating": 5, "internal_note": "配合度極高" }`

### `GET /admin/centers/{id}/schedule`
- **描述**: 取得範圍內的完整排班 (Grid View)。含 Rules (週期) 與 Excpetions (例外)。
- **Query**: `start=2026-01-01&end=2026-01-07`

---

## 6. 資源管理庫 (Resource CRUD) - Admin Only

### `GET /admin/centers/{id}/teachers`
- **描述**: 列表顯示中心老師。
- **Response**: `[{ "id": 1, "name": "...", "status": "ACTIVE", "rating": 5 }]`

### `POST /admin/centers/{id}/teachers/invite`
- **描述**: 產生邀請連結或邀請碼，並發送 LINE 通知。
- **Body**: `{ "line_user_id": "...", "contact_info": "請加 Line: abc", "note": "誠徵鋼琴老師" }`
- **Response**: `{ "invite_code": "INV-888", "invite_url": "https://..." }`

### `GET /admin/centers/{id}/rooms`
- **描述**: 教室列表。
- **Response**: `[{ "id": 1, "name": "Room A", "capacity": 10 }]`

### `POST /admin/centers/{id}/rooms`
- **Body**: `{ "name": "Room B", "capacity": 5 }`

### `GET /admin/centers/{id}/courses`
- **描述**: 課程模板列表。
- **Response**: `[{ "id": 10, "name": "Piano Basic", "duration": 60, "color": "#FF9900" }]`

---

## 7. 課表模板管理 (Grid Templates) - Admin Only

### `GET /admin/centers/{id}/templates`
- **描述**: 取得中心所有排課模板。

### `POST /admin/centers/{id}/templates`
- **Body**: `{ "name": "...", "row_type": "ROOM" }`

### `GET /admin/centers/{id}/templates/{id}/cells`
- **描述**: 取得模板內的格子。

### `POST /admin/centers/{id}/templates/{id}/cells/bulk`
- **描述**: 批量更新格子的定義 (row/col/time/refs)。

---

## 7. 中心設定 (Center Config) - Admin Only

### `GET /admin/centers/{id}`
- **描述**: 取得中心基本資訊與 Policy。
- **Response**: `{ "name": "Mozart", "plan": "PRO", "policy": { "teacher_buffer": 10 } }`

### `PATCH /admin/centers/{id}`
- **描述**: 修改中心資訊。
- **Body**: `{ "name": "new name", "settings": {...} }`

---

## 8. 管理員帳號管理 (Admin User Management) - Owner/Admin Only

### `GET /admin/centers/{id}/users`
- **描述**: 取得中心所有管理員列表。
- **Response**: `[{ "id": 1, "email": "...", "role": "OWNER", "name": "Admin A" }]`

### `POST /admin/centers/{id}/users`
- **描述**: 新增管理員。
- **Body**: `{ "email": "...", "password": "...", "role": "STAFF", "name": "New Staff" }`

### `DELETE /admin/centers/{id}/users/{uid}`
- **描述**: 移除或禁用管理員。

---

## 9. 稽核紀錄 (Audit Logs)

### `GET /admin/centers/{id}/audit-logs`
- **描述**: 取得中心的操作紀錄。
- **Query**: `actor_id=1&start=...&end=...`
- **Response**: 
  ```json
  [
    {
      "id": 1,
      "actor_name": "Admin A",
      "action": "APPROVE_RESCHEDULE",
      "target": "Course X",
      "created_at": "2026-01-20T10:00:00",
      "note": "Reason: Approved by request"
    }
  ]
  ```

---

## 8. 匯出相關 (Export)

> 註：圖片匯出建議採前端 `html2canvas` 方案，故無專屬 Backend API。
> 若需匯出資料 (CSV/Excel)，可呼叫以下 API：

### `GET /teacher/me/schedule/export`
- **Query**: `format=ics` (iCalendar 格式，可匯入 Google Calendar)
- **Response**: `.ics` 檔案下載。

---

## 9. 公用與輔助 (Common & Utils)

### `GET /common/hashtags/search`
- **描述**: (UX優化) 老師輸入標籤時的模糊搜尋建議。
- **Query**: `q=舞蹈`
- **Response**: `{ "results": ["#街舞", "#國標舞", "#MV舞蹈"] }`

---

## 10. 公開與展示 (Public & Demo)

### `POST /public/validate`
- **描述**: 供官網 Demo 沙盒使用的匿名驗證接口。
- **Note**: 僅支援簡單的 Overlap/Buffer 邏輯，不檢查中心 Policy 或權限。
- **Body**: `{ "start_at": "...", "end_at": "...", "teacher_id": 1, "room_id": 2 }`
- **Response**: `{ "valid": true, "conflicts": [] }`
