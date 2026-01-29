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
- **描述**: 取得個人資料 (Bio, Skills, City, District, Hashtags)。
- **Response**:
  ```json
  {
      "name": "Teacher A",
      "bio": "...",
      "city": "台北市",
      "district": "大安區",
      "skills": [
          { 
            "category": "舞蹈", 
            "name": "街舞", 
            "hashtags": ["#MV舞蹈", "#街舞基礎"] 
          }
      ],
      "personal_hashtags": ["#專業準時", "#親和力強", "#資深認證"],
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
          { "category": "舞蹈", "name": "街舞", "level": "ADVANCED", "hashtags": ["#MV舞蹈", "#小班制"] },
          { "category": "瑜伽", "name": "空中瑜伽", "level": "INTERMEDIATE", "hashtags": ["#安靜", "#全英文教學"] }
      ],
      "personal_hashtags": ["#專業準時", "#親和力強"], 
      "is_open_to_hiring": true,
      "city": "台北市",
      "district": "大安區",
      "public_contact_info": "LineID: happy_teacher"
  }
  ```

### `PATCH /teacher/me/avatar`
- **描述**: 更新個人頭像。
- **Body** (Multipart): `avatar` (Image)
- **Response**: `{ "avatar_url": "..." }`

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

### `POST /teacher/personal-events/{id}/copy`
- **描述**: 複製個人行程。
- **Body**:
```json
{
  "title": "私人練琴 (複製)",
  "start_at": "2026-01-25T14:00:00",
  "end_at": "2026-01-25T16:00:00"
}
```
- **Response**: `{ "id": 456, "title": "私人練琴 (複製)", ... }`

### `PATCH /teacher/personal-events/{id}`
- **描述**: 修改行程（支援拖曳改時、內容修訂）。
- **Body**:
```json
{
  "title": "私人練琴 (改)",
  "start_at": "2026-01-21T14:00:00",
  "end_at": "2026-01-21T16:00:00",
  "update_mode": "SINGLE", // 必填（若為循環行程）
  "recurrence": {
    // 僅當 update_mode = ALL 時可更新循環規則
    "type": "WEEKLY",
    "interval": 1,
    "weekdays": [1, 3],
    "until": "2026-06-30"
  }
}
```
- **update_mode 說明**:
  - `SINGLE`: 僅修改此單一場次（原規則產生一個 CANCEL 例外，新規則產生一個 ADD 例外）
  - `FUTURE`: 修改此場次及之後所有場次（原規則截斷，新規則從此場次開始）
  - `ALL`: 修改整串循環規則（更新 recurrence 欄位）
- **Error Responses**:
  - `E_INVALID_UPDATE_MODE`: update_mode 與行程類型不符（非循環行程不應傳此參數）

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

### `POST /teacher/exceptions/{exception_id}/revoke`
- **描述**: 撤回待審核的異動申請（僅限 PENDING 狀態）。
- **Path Parameters**:
  - `exception_id` (BIGINT): 例外申請 ID
- **Response**:
```json
{
  "code": "SUCCESS",
  "message": "Exception revoked successfully",
  "data": {
    "id": 123,
    "status": "REVOKED",
    "revoked_at": "2026-01-22T10:00:00Z"
  }
}
```
- **Error Responses**:
  - `E_FORBIDDEN`: 僅能撤回自己提交的申請
  - `E_INVALID_STATUS`: 僅 PENDING 狀態可撤回

### `GET /teacher/me/exceptions`
- **描述**: 查看已提交的變更申請 (停課/改期) 狀態歷程。
- **Response**: `[ { "id": 1, "type": "CANCEL", "status": "PENDING", "original_date": "...", "reason": "..." } ]`

### `POST /teacher/exceptions/{exception_id}/reject`
- **描述**: **中心端** 拒絕老師的異動申請（管理員專用）。
- **Path Parameters**:
  - `exception_id` (BIGINT): 例外申請 ID
- **Body**:
```json
{
  "reason": "該時段已無法調整，請見諒"
}
```
- **Response**:
```json
{
  "code": "SUCCESS",
  "message": "Exception rejected successfully",
  "data": {
    "id": 123,
    "status": "REJECTED",
    "rejected_at": "2026-01-22T10:00:00Z",
    "reject_reason": "該時段已無法調整，請見諒"
  }
}
```

---

## 4. 中心後台 API (Center Admin)

### `POST /admin/centers/{id}/validate`
- **描述**: 預檢排課衝突 (Validate Engine)。
- **Body**: `{ "offering_id": 10, "teacher_id": 1, "room_id": 2, "start_at": "...", "end_at": "...", "override_buffer": false }`
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
- **Body**: `[ { "offering_id":..., "teacher_id":..., "room_id":..., "weekday":..., "start_time":..., "end_time":..., "effective_start":..., "effective_end":... } ]`

### `GET /admin/centers/{id}/rooms/{room_id}/schedule`
- **描述**: 取得特定教室的課表 (用於檢查教室佔用與空檔)。
- **Query**: `start=2026-01-01&end=2026-01-07`
- **Response**:
  ```json
  [
    {
      "session_id": "rule_55_20260120",
      "course_name": "瑜伽基礎",
      "teacher_name": "Teacher A",
      "start": "2026-01-20T10:00:00",
      "end": "2026-01-20T11:00:00",
      "status": "NORMAL"
    }
  ]
  ```

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
      "city": "台北市",
      "district": "大安區",
      "skills": [{ "category": "有氧", "name": "Zumba" }],
      "hashtags": ["#熱情", "#週末可"],
      "availability": "AVAILABLE", // or CONFLICT, BUFFER_CONFLICT
      "internal_note": "適合教初學者",
      "internal_rating": 5,
      "personal_hashtags": ["#親和力", "#耐心"],
      "is_member": true, 
      "public_contact_info": "LineID: ..." // 僅當 (is_member=false AND open_hiring=true) 時有值
    }
  ]
  ```


### `GET /admin/centers/{id}/matching/teachers/details`
- **描述**: 查看媒合老師的完整履歷 (Bio/Certs)。

### `GET /admin/talent/search`
- **描述**: **(人才庫搜尋)** 主動搜尋已開啟「求職媒合」的老師。
- **Query**:
  - `category`: 有氧 (Filter)
  - `skill`: Zumba (Filter)
  - `tag`: "#熱情" (Filter)
  - `city`: 台北市 (Filter)
  - `district`: 大安區 (Filter)
  - `keyword`: "王小明" (Search Name/Bio)
- **Response**:
  ```json
  [
    {
      "teacher_id": 105,
      "name": "Teacher Bloom",
      "bio": "Expert Yoga Instructor...",
      "city": "台北市",
      "district": "大安區",
      "skills": [{ "category": "瑜伽", "name": "空中瑜伽", "hashtags": ["#安靜", "#平日可"] }],
      "personal_hashtags": ["#專業認證", "#不遲到"],
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

### `POST /admin/centers/{id}/rules`
- **描述**: 新增一條循環排課規則。
- **Body**:
  ```json
  {
    "offering_id": 10,
    "teacher_id": 1,
    "room_id": 2,
    "weekday": 2,
    "start_time": "10:00:00",
    "end_time": "11:00:00",
    "effective_start": "2026-01-01",
    "effective_end": "2026-02-28",
    "lock_at": "2026-01-15T23:59:59"
  }
  ```

---

## 12. 資源管理庫 (Resource CRUD) - Admin Only

### `GET /admin/centers/{id}/teachers`
- **描述**: 列表顯示中心老師。
- **Query**: `page=1&limit=20&sort_by=name&sort_order=ASC&status=ACTIVE`
- **Response**:
```json
{
  "code": "SUCCESS",
  "data": [
    { "id": 1, "name": "...", "status": "ACTIVE", "rating": 5 }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 50,
    "total_pages": 3,
    "has_next": true,
    "has_prev": false
  }
}
```

### `POST /admin/centers/{id}/teachers/invite`
- **描述**: 產生邀請連結或邀請碼，並發送 LINE 通知。
- **Body**: `{ "line_user_id": "...", "contact_info": "請加 Line: abc", "note": "誠徵鋼琴老師" }`
- **Response**: `{ "invite_code": "INV-888", "invite_url": "https://..." }`

### `GET /admin/centers/{id}/invitations`
- **描述**: 取得中心所有邀請碼列表（含狀態統計）。
- **Query**: `page=1&limit=20&status=PENDING&sort_by=created_at&sort_order=DESC`
- **Response**:
```json
{
  "code": "SUCCESS",
  "data": [
    {
      "id": 1,
      "invite_code": "INV-888",
      "status": "PENDING",
      "teacher_name": null,
      "created_at": "2026-01-20T10:00:00",
      "expires_at": "2026-01-27T10:00:00",
      "used_at": null
    }
  ],
  "pagination": { ... },
  "statistics": {
    "total": 100,
    "pending": 50,
    "accepted": 45,
    "expired": 5
  }
}
```

### `GET /admin/centers/{id}/invitations/stats`
- **描述**: 取得中心邀請統計數據（用於人才庫管理頁面）。
- **Response**:
```json
{
  "code": "SUCCESS",
  "data": {
    "total": 156,
    "pending": 23,
    "accepted": 45,
    "declined": 8,
    "expired": 12,
    "recent_pending": 5
  }
}
```

---

## 13. 老師邀請功能 (Teacher Invitations)

### `GET /teacher/me/invitations`
- **描述**: 取得老師收到的所有人才庫邀請列表。
- **Query**: `status=PENDING|ACCEPTED|DECLINED|EXPIRED`（選填，不傳則回傳全部）
- **Response**:
```json
{
  "code": "SUCCESS",
  "data": [
    {
      "id": 1,
      "center_id": 10,
      "center_name": "莫札特音樂教室",
      "invite_type": "TALENT_POOL",
      "message": "誠摯邀請您加入我們的人才庫！",
      "status": "PENDING",
      "created_at": "2026-01-20T10:00:00",
      "expires_at": "2026-01-27T10:00:00",
      "responded_at": null
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 5,
    "total_pages": 1,
    "has_next": false,
    "has_prev": false
  }
}
```

### `POST /teacher/me/invitations/respond`
- **描述**: 老師回應人才庫邀請（接受或婉拒）。
- **Body**:
```json
{
  "invitation_id": 1,
  "response": "ACCEPT" // 或 "DECLINE"
}
```
- **Response（接受）**:
```json
{
  "code": "SUCCESS",
  "message": "Invitation accepted successfully",
  "data": {
    "id": 1,
    "status": "ACCEPTED",
    "responded_at": "2026-01-22T10:00:00Z",
    "center_name": "莫札特音樂教室"
  }
}
```
- **Response（婉拒）**:
```json
{
  "code": "SUCCESS",
  "message": "Invitation declined successfully",
  "data": {
    "id": 1,
    "status": "DECLINED",
    "responded_at": "2026-01-22T10:00:00Z"
  }
}
```
- **Error Responses**:
  - `E_FORBIDDEN`: 無權限回應此邀請
  - `E_INVALID_STATUS`: 邀請已過期或已處理

### `GET /teacher/me/invitations/pending-count`
- **描述**: 取得待處理邀請數量（用於側邊欄 Badge 顯示）。
- **Response**:
```json
{
  "code": "SUCCESS",
  "data": {
    "count": 3
  }
}
```

### `GET /admin/centers/{id}/rooms`
- **描述**: 教室列表。
- **Response**: `[{ "id": 1, "name": "Room A", "capacity": 10, "is_active": true }]`

### `POST /admin/centers/{id}/rooms`
- **Body**: `{ "name": "Room B", "capacity": 5 }`

### `PATCH /admin/centers/{id}/rooms/{room_id}`
- **Body**: `{ "name": "...", "capacity": 10, "is_active": true }`

### `DELETE /admin/centers/{id}/rooms/{room_id}`
- **描述**: 軟刪除教室（將 is_active 設為 false）。
- **Note**: 若教室仍有排課規則，回傳 `E_ROOM_IN_USE` 錯誤。
- **Response**: `{ "code": "SUCCESS", "message": "Room deactivated successfully" }`

### `GET /admin/centers/{id}/courses`
- **描述**: 課程模板列表。
- **Response**: `[{ "id": 10, "name": "...", "duration": 60, "color": "#FF9900", "room_buffer_min": 10, "teacher_buffer_min": 5, "is_active": true }]`

### `POST /admin/centers/{id}/courses`
- **描述**: 新增課程模板。
- **Body**: `{ "name": "...", "duration": 60, "color": "#...", "room_buffer_min": 10, "teacher_buffer_min": 5 }`

### `PATCH /admin/centers/{id}/courses/{course_id}`
- **描述**: 更新課程模板。
- **Body**: `{ "name": "...", "duration": 90, "color": "#...", "room_buffer_min": 15, "teacher_buffer_min": 10 }`

### `DELETE /admin/centers/{id}/courses/{course_id}`
- **描述**: 軟刪除課程模板（將 is_active 設為 false）。
- **Note**: 若課程仍有關聯的班別 (offerings)，回傳 `E_COURSE_IN_USE` 錯誤。
- **Response**: `{ "code": "SUCCESS", "message": "Course deactivated successfully" }`

### `GET /admin/centers/{id}/offerings`
- **描述**: 班別列表。
- **Response**: `[{ "id": 55, "course_id": 10, "course_name": "鋼琴課", "allow_buffer_override": true, "is_active": true }]`

### `POST /admin/centers/{id}/offerings`
- **描述**: 新增班別。
- **Body**:
```json
{
  "course_id": 10,
  "default_room_id": 1,
  "default_teacher_id": null,
  "allow_buffer_override": false
}
```

### `PATCH /admin/centers/{id}/offerings/{offering_id}`
- **描述**: 更新班別。
- **Body**: `{ "default_room_id": 2, "allow_buffer_override": true }`

### `DELETE /admin/centers/{id}/offerings/{offering_id}`
- **描述**: 軟刪除班別（將 is_active 設為 false）。
- **Note**: 若班別仍有排課規則，回傳 `E_OFFERING_HAS_RULES` 錯誤。
- **Response**: `{ "code": "SUCCESS", "message": "Offering deactivated successfully" }`

### `POST /admin/centers/{id}/offerings/{offering_id}/copy`
- **描述**: 複製班別（用於快速建立相似班別）。
- **Body**:
```json
{
  "new_name": "鋼琴課 - 進階班",
  "effective_start": "2026-03-01",
  "effective_end": "2026-06-30",
  "copy_rules": true
}
```
- **Response**:
```json
{
  "code": "SUCCESS",
  "data": {
    "id": 56,
    "name": "鋼琴課 - 進階班",
    "rules_copied": 5
  }
}
```

### `GET /admin/centers/{id}/reports/plan-usage`
- **描述**: 取得方案用量報表 (老師數、排課總時數、方案上限)。
- **Response**: `{ "teacher_count": 15, "teacher_limit": 20, "session_hours": 120, "plan_name": "PRO" }`

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
- **Response**: `{ "name": "Mozart", "plan": "PRO", "policy": { "allow_public_register": true, "exception_lead_days": 14 } }`

### `PATCH /admin/centers/{id}`
- **描述**: 修改中心資訊與預設政策。
- **Body**: `{ "name": "...", "settings": { "exception_lead_days": 7 } }`

### `GET /admin/centers/{id}/holidays`
- **描述**: 取得中心自定義假日列表。
- **Response**: `[ { "id": 1, "date": "2026-02-14", "name": "情人節停課" } ]`

### `POST /admin/centers/{id}/holidays`
- **Body**: `{ "date": "2026-02-14", "name": "情人節停課" }`

### `DELETE /admin/centers/{id}/holidays/{holiday_id}`
- **描述**: 刪除假日設定。

### `POST /admin/centers/{id}/holidays/bulk`
- **描述**: 批量匯入假日（支援常見國定假日）。
- **Body**:
```json
{
  "holidays": [
    { "date": "2026-01-01", "name": "元旦" },
    { "date": "2026-02-14", "name": "春節" },
    { "date": "2026-04-04", "name": "兒童節" },
    { "date": "2026-10-10", "name": "雙十節" }
  ],
  "is_recurring": false
}
```
- **Response**:
```json
{
  "code": "SUCCESS",
  "data": {
    "imported": 4,
    "skipped": 0,
    "holidays": [
      { "id": 1, "date": "2026-01-01", "name": "元旦" },
      { "id": 2, "date": "2026-02-14", "name": "春節" },
      { "id": 3, "date": "2026-04-04", "name": "兒童節" },
      { "id": 4, "date": "2026-10-10", "name": "雙十節" }
    ]
  }
}
```
- **Note**: 若某日期已有假日設定，該筆會被標記為 `skipped`，不回報錯誤。

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
- **Query**: `page=1&limit=50&actor_id=1&start=...&end=...&sort_by=created_at&sort_order=DESC`
- **Response**:
```json
{
  "code": "SUCCESS",
  "data": [
    {
      "id": 1,
      "actor_name": "Admin A",
      "actor_email": "admin@example.com",
      "action": "APPROVE_RESCHEDULE",
      "target_type": "schedule_exception",
      "target_id": 123,
      "target_name": "鋼琴課 - 1/22",
      "payload": {
        "before": { "status": "PENDING" },
        "after": { "status": "APPROVED" }
      },
      "note": "Reason: Approved by request",
      "created_at": "2026-01-20T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 1000,
    "total_pages": 20,
    "has_next": true,
    "has_prev": false
  }
}
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

## 11. 通用分頁與排序規範 (Pagination & Sorting Standards)
所有列表 API 支援以下標準參數：

### 11.1 通用查詢參數
| 參數 | 類型 | 必填 | 預設值 | 說明 |
|:---|:---:|:---:|:---:|:---|
| `page` | INT | 否 | 1 | 頁碼 |
| `limit` | INT | 否 | 20 | 每頁筆數 (最大 100) |
| `sort_by` | STRING | 否 | 依各 API 定義 | 排序欄位 |
| `sort_order` | STRING | 否 | ASC | 排序方向 (ASC/DESC) |

### 11.2 通用分頁 Response 格式
```json
{
  "code": "SUCCESS",
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

### 11.3 支援排序欄位對照表
| API | 支援的 sort_by 值 | 預設排序 |
|:---|:---|:---|
| `/admin/centers/{id}/teachers` | `name`, `rating`, `created_at`, `last_active_at` | `name:ASC` |
| `/admin/talent/search` | `match_score`, `rating`, `name`, `created_at` | `match_score:DESC` |
| `/admin/centers/{id}/teachers/{tid}/certificates` | `issued_at`, `name` | `issued_at:DESC` |
| `/admin/centers/{id}/audit-logs` | `created_at`, `action` | `created_at:DESC` |
| `/admin/centers/{id}/exceptions` | `created_at`, `original_date` | `created_at:DESC` |
| `/teacher/personal-events` | `start_at`, `title` | `start_at:ASC` |
| `/teacher/me/exceptions` | `created_at`, `original_date` | `created_at:DESC` |
| `/admin/centers/{id}/rooms` | `name`, `capacity` | `name:ASC` |
| `/admin/centers/{id}/courses` | `name`, `created_at` | `name:ASC` |
| `/admin/centers/{id}/offerings` | `name`, `created_at` | `created_at:DESC` |
| `/admin/centers/{id}/invitations` | `created_at`, `expires_at` | `created_at:DESC` |
| `/admin/centers/{id}/holidays` | `date`, `name` | `date:ASC` |

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
