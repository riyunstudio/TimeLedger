# API.md — TimeLedger API Contract（v2.0）
日期：2026-01-18

> Auth（MVP stub）：
- Teacher API：Header `X-Teacher-Id: <int64>`
- Admin API：Header `X-Admin-Token: <string>`（伺服器端對應 admin_users）
- Admin API 必須同時解析 `center_id`（path param）並強制 center scoping

> 統一回傳（所有 API）：
```json
{"code":0,"message":"OK","datas":{}}
```

錯誤碼請參考：`docs/ERR_CODES.md`

---

## 1. Teacher APIs

### 1.1 GET /api/teacher/me
取得老師資訊

Response datas：
```json
{"id":1,"display_name":"王小明","avatar_url":"...","created_at":"...","updated_at":"..."}
```

---

### 1.2 POST /api/teacher/sessions
建立個人行程（Personal Session）

Request:
```json
{
  "start_at":"2026-01-20 10:00:00",
  "end_at":"2026-01-20 12:00:00",
  "title":"私人教學",
  "note":"其他中心"
}
```
Validation:
- start_at < end_at
- 不可與自己既有 personal sessions overlap（否則 TIME_CONFLICT）

Success datas：
```json
{"id":101}
```

TIME_CONFLICT datas：見 ERR_CODES.md 的 conflict list

---

### 1.3 GET /api/teacher/sessions?from=YYYY-MM-DD&to=YYYY-MM-DD
查詢個人行程列表

Response datas：
```json
{
  "items":[
    {"id":101,"start_at":"...","end_at":"...","title":"...","note":"..."}
  ]
}
```

---

### 1.4 PUT /api/teacher/sessions/{id}
更新個人行程（同 1.2 驗證）

---

### 1.5 DELETE /api/teacher/sessions/{id}
刪除個人行程

---

### 1.6 GET /api/teacher/schedule?from=YYYY-MM-DD&to=YYYY-MM-DD
老師視角的 Schedule View（中心排課 + 個人行程）

Response datas：
```json
{
  "items":[
    {
      "type":"CENTER_SESSION|PERSONAL_SESSION",
      "center_id":1,
      "offering_id":10,
      "rule_id":55,
      "session_id":9001,
      "date":"2026-01-20",
      "start_at":"2026-01-20 09:00:00",
      "end_at":"2026-01-20 11:00:00",
      "title":"英文會話-A班",
      "status":"NORMAL|CANCELLED|RESCHEDULED",
      "meta":{"exception_id":123}
    }
  ]
}
```

---

## 2. Admin APIs（Center-scoped）

### 2.1 POST /api/admin/centers
建立中心

Request:
```json
{"name":"台北運動中心","code":"TP001"}
```
Success datas：
```json
{"id":1}
```

---

### 2.2 GET /api/admin/centers
列出 admin 可管理中心

Response datas：
```json
{"items":[{"id":1,"name":"...","code":"...","plan":"pro","trial_end_at":"..."}]}
```

---

### 2.3 POST /api/admin/centers/{center_id}/offerings
建立 Offering

Request:
```json
{"name":"英文會話-A班","duration_min":120,"status":"ACTIVE"}
```
Success datas：
```json
{"id":10}
```

---

### 2.4 GET /api/admin/centers/{center_id}/offerings
Response datas：
```json
{"items":[{"id":10,"name":"...","duration_min":120,"status":"ACTIVE"}]}
```

---

### 2.5 POST /api/admin/centers/{center_id}/memberships/invite
邀請老師加入中心（membership=INVITED）

Request:
```json
{"teacher_id":1}
```
Rules:
- Plan gating：若已達 max_teachers → PLAN_LIMIT_EXCEEDED

Success datas：
```json
{"id":501,"status":"INVITED"}
```

---

### 2.6 POST /api/admin/centers/{center_id}/memberships/{id}/activate
啟用 membership（ACTIVE）

Success datas：
```json
{"id":501,"status":"ACTIVE"}
```

---

### 2.7 POST /api/admin/centers/{center_id}/rules
建立 Rule（固定排課）

Request:
```json
{
  "teacher_id":1,
  "offering_id":10,
  "weekday":2,
  "start_time":"09:00",
  "end_time":"11:00"
}
```
Rules:
- membership 必須 ACTIVE（否則 MEMBERSHIP_NOT_ACTIVE）
- 衝突則拒絕 TIME_CONFLICT（建議直接拒絕）

Success datas：
```json
{"id":55}
```

---

### 2.8 GET /api/admin/centers/{center_id}/rules
Response datas：
```json
{"items":[{"id":55,"teacher_id":1,"offering_id":10,"weekday":2,"start_time":"09:00","end_time":"11:00"}]}
```

---

### 2.9 POST /api/admin/centers/{center_id}/validate
排課前檢查（候選時段）

Request:
```json
{
  "teacher_id":1,
  "date":"2026-01-20",
  "start_at":"2026-01-20 10:00:00",
  "end_at":"2026-01-20 12:00:00",
  "ignore_rule_id":55
}
```

Success datas：
```json
{"ok":true}
```

Conflict：HTTP 409 + TIME_CONFLICT + conflict list（見 ERR_CODES.md）

---

### 2.10 POST /api/admin/centers/{center_id}/exceptions
建立例外（停課/改期）

Request（CANCEL）：
```json
{"type":"CANCEL","rule_id":55,"original_date":"2026-01-20","reason":"颱風"}
```

Request（RESCHEDULE）：
```json
{
  "type":"RESCHEDULE",
  "rule_id":55,
  "original_date":"2026-01-20",
  "new_start_at":"2026-01-21 10:00:00",
  "new_end_at":"2026-01-21 12:00:00",
  "reason":"老師行程"
}
```

Rules:
- 同日有效例外唯一性 → EXCEPTION_ALREADY_EXISTS
- RESCHEDULE 建立/核准都必須 validate 新時段（衝突→TIME_CONFLICT）

Success datas：
```json
{"id":7001,"status":"PENDING"}
```

---

### 2.11 GET /api/admin/centers/{center_id}/exceptions?status=PENDING|APPROVED|REJECTED
Response datas：
```json
{"items":[{"id":7001,"type":"RESCHEDULE","status":"PENDING","rule_id":55,"original_date":"2026-01-20","new_start_at":"...","reason":"..."}]}
```

---

### 2.12 POST /api/admin/centers/{center_id}/exceptions/{id}/approve
Approve 例外（RESCHEDULE 需再次 validate）
Success datas：
```json
{"id":7001,"status":"APPROVED"}
```

---

### 2.13 POST /api/admin/centers/{center_id}/exceptions/{id}/reject
Reject 例外
Success datas：
```json
{"id":7001,"status":"REJECTED"}
```

---

### 2.14 GET /api/admin/centers/{center_id}/schedule?from=YYYY-MM-DD&to=YYYY-MM-DD&teacher_id=1
中心視角 Schedule View（可選 teacher filter）
Response datas：同 teacher schedule items，但可加 teacher_name 等 meta
