# API 端點完整對照表

## 回應格式標準

所有 API 均採用統一格式：
```json
{
  "code": 0,
  "message": "Success",
  "datas": { ... }
}
```

---

## Teacher API (教師端)

### 個人資料
| Method | Endpoint | 參數 | 回應 | 說明 |
|:---|:---|:---|:---|:---|
| GET | `/api/v1/teacher/me/profile` | - | `{datas: TeacherProfileResource}` | 取得個人資料 |
| PUT | `/api/v1/teacher/me/profile` | `{bio, city, district, public_contact_info, is_open_to_hiring}` | `{datas: TeacherProfileResource}` | 更新個人資料 |

### 所屬中心
| Method | Endpoint | 參數 | 回應 | 說明 |
|:---|:---|:---|:---|:---|
| GET | `/api/v1/teacher/me/centers` | - | `{datas: CenterMembership[]}` | 取得所屬中心列表 |

### 課表
| Method | Endpoint | 參數 | 回應 | 說明 |
|:---|:---|:---|:---|:---|
| GET | `/api/v1/teacher/me/schedule` | `?from=YYYY-MM-DD&to=YYYY-MM-DD` | `{datas: TeacherScheduleItem[]}` | 取得課表 |

**TeacherScheduleItem 結構：**
```go
type TeacherScheduleItem struct {
  ID         string `json:"id"`
  Type       string `json:"type"`       // "CENTER_SESSION"
  Title      string `json:"title"`      // "Center Session"
  Date       string `json:"date"`       // "2006-01-02"
  StartTime  string `json:"start_time"` // "10:00"
  EndTime    string `json:"end_time"`   // "11:00"
  RoomID     uint   `json:"room_id"`
  TeacherID  *uint  `json:"teacher_id"`
  CenterID   uint   `json:"center_id"`
  CenterName string `json:"center_name"`
  Status     string `json:"status"`     // "NORMAL", "PENDING_CANCEL", etc.
}
```

### 例外申請
| Method | Endpoint | 參數 | 回應 | 說明 |
|:---|:---|:---|:---|:---|
| GET | `/api/v1/teacher/exceptions` | `?status=PENDING` (可選) | `{datas: ScheduleException[]}` | 取得例外申請列表 |
| POST | `/api/v1/teacher/exceptions` | `{center_id, rule_id, original_date, type, reason}` | `{datas: ScheduleException}` | 申請停課/改期 |
| POST | `/api/v1/teacher/exceptions/:id/revoke` | - | - | 撤回申請 |

### 課堂筆記 (新增)
| Method | Endpoint | 參數 | 回應 | 說明 |
|:---|:---|:---|:---|:---|
| GET | `/api/v1/teacher/sessions/note` | `?rule_id=1&session_date=2026-01-25` | `{datas: {note: SessionNote, is_new: boolean}}` | 取得課堂筆記 |
| PUT | `/api/v1/teacher/sessions/note` | `{rule_id, session_date, content, prep_note}` | `{datas: SessionNote}` | 新增/更新筆記 |

**SessionNote 結構：**
```go
type SessionNote struct {
  ID          uint      `json:"id"`
  RuleID      *uint     `json:"rule_id"`
  SessionDate string    `json:"session_date"` // "2006-01-02"
  Content     string    `json:"content"`
  PrepNote    string    `json:"prep_note"`
  UpdatedAt   time.Time `json:"updated_at"`
}
```

---

## Notification API (通知)

| Method | Endpoint | 參數 | 回應 | 說明 |
|:---|:---|:---|:---|:---|
| GET | `/api/v1/notifications` | `?limit=50` | `{datas: {notifications: Notification[], unread_count: number}}` | 取得通知列表 |
| GET | `/api/v1/notifications/unread-count` | - | `{datas: {unread_count: number}}` | 取得未讀數 |
| POST | `/api/v1/notifications/:id/read` | - | - | 標記已讀 |
| POST | `/api/v1/notifications/read-all` | - | - | 全部標記已讀 |

---

## Auth API (認證)

| Method | Endpoint | 參數 | 回應 | 說明 |
|:---|:---|:---|:---|:---|
| POST | `/api/v1/auth/admin/login` | `{email, password}` | `{datas: {token, refresh_token, user}}` | 管理員登入 |
| POST | `/api/v1/auth/teacher/line/login` | `{line_user_id, access_token}` | `{datas: {token, refresh_token, teacher}}` | 教師 LINE 登入 |
| POST | `/api/v1/auth/refresh` | `{refresh_token}` | `{datas: {token, refresh_token}}` | 刷新 Token |
| POST | `/api/v1/auth/logout` | - | - | 登出 |

---

## Admin API (管理端)

### Center (中心)
| Method | Endpoint | 參數 | 回應 | 說明 |
|:---|:---|:---|:---|:---|
| GET | `/api/v1/admin/centers` | - | `{datas: Center[]}` | 取得中心列表 |
| POST | `/api/v1/admin/centers` | `{name, ...}` | `{datas: Center}` | 建立中心 |

### Rooms (教室)
| Method | Endpoint | 參數 | 回應 | 說明 |
|:---|:---|:---|:---|:---|
| GET | `/api/v1/admin/centers/:id/rooms` | - | `{datas: Room[]}` | 取得教室列表 |
| POST | `/api/v1/admin/centers/:id/rooms` | `{name, capacity}` | `{datas: Room}` | 建立教室 |

### Courses (課程)
| Method | Endpoint | 參數 | 回應 | 說明 |
|:---|:---|:---|:---|:---|
| GET | `/api/v1/admin/centers/:id/courses` | - | `{datas: Course[]}` | 取得課程列表 |
| POST | `/api/v1/admin/centers/:id/courses` | `{name, type, duration}` | `{datas: Course}` | 建立課程 |

### Offerings (開課)
| Method | Endpoint | 參數 | 回應 | 說明 |
|:---|:---|:---|:---|:---|
| GET | `/api/v1/admin/centers/:id/offerings` | - | `{datas: Offering[]}` | 取得開課列表 |
| POST | `/api/v1/admin/centers/:id/offerings` | `{course_id, room_id, weekday, start_time, end_time}` | `{datas: Offering}` | 建立開課 |

### Templates (範本)
| Method | Endpoint | 參數 | 回應 | 說明 |
|:---|:---|:---|:---|:---|
| GET | `/api/v1/admin/centers/:id/templates` | - | `{datas: TimetableTemplate[]}` | 取得範本列表 |
| POST | `/api/v1/admin/centers/:id/templates` | `{name, cells}` | `{datas: TimetableTemplate}` | 建立範本 |

### Scheduling (排課)
| Method | Endpoint | 參數 | 回應 | 說明 |
|:---|:---|:---|:---|:---|
| POST | `/api/v1/admin/centers/:id/validate` | `{rule_id, start_date, end_date}` | `{datas: ValidationResult}` | 驗證排課衝突 |
| POST | `/api/v1/admin/centers/:id/expand` | `{rule_ids, start_date, end_date}` | `{datas: ExpandedSchedule[]}` | 展開排課規則 |
| POST | `/api/v1/admin/centers/:id/exceptions` | `{rule_id, type, ...}` | `{datas: ScheduleException}` | 建立例外 |
| POST | `/api/v1/admin/centers/:id/exceptions/:exceptionId/review` | `{status, note}` | - | 審核例外 |
| GET | `/api/v1/admin/centers/:id/exceptions` | `?start_date=&end_date=` | `{datas: ScheduleException[]}` | 取得例外列表 |

### Smart Matching (智慧媒合)
| Method | Endpoint | 參數 | 回應 | 說明 |
|:---|:---|:---|:---|:---|
| POST | `/api/v1/admin/centers/:id/matching/teachers` | `{offering_id, ...}` | `{datas: TeacherMatch[]}` | 尋找匹配教師 |
| GET | `/api/v1/admin/centers/:id/matching/teachers/search` | `?keyword=` | `{datas: Teacher[]}` | 搜尋教師 |

### Admin Users (管理員)
| Method | Endpoint | 參數 | 回應 | 說明 |
|:---|:---|:---|:---|:---|
| GET | `/api/v1/admin/centers/:id/users` | - | `{datas: AdminUser[]}` | 取得管理員列表 |
| POST | `/api/v1/admin/centers/:id/users` | `{email, role}` | `{datas: AdminUser}` | 建立管理員 |

### Export (匯出)
| Method | Endpoint | 參數 | 回應 | 說明 |
|:---|:---|:---|:---|:---|
| POST | `/api/v1/admin/export/schedule/csv` | `{center_id, start_date, end_date}` | CSV 檔案 | 匯出課表 CSV |
| POST | `/api/v1/admin/export/schedule/pdf` | `{center_id, start_date, end_date}` | PDF 檔案 | 匯出課表 PDF |
| GET | `/api/v1/admin/export/centers/:centerId/teachers/csv` | - | CSV 檔案 | 匯出教師 CSV |
| GET | `/api/v1/admin/export/centers/:centerId/exceptions/csv` | - | CSV 檔案 | 匯出例外 CSV |

---

## 測試檔案

### Repository Tests
- `testing/test/session_note_test.go` - SessionNote Repository CRUD 測試

### Controller Tests  
- `testing/test/teacher_session_note_test.go` - TeacherController SessionNote API 測試

### 執行測試
```bash
# 需要 CGO 環境
go test ./testing/test/... -v

# 執行 SessionNote 測試
go test ./testing/test/... -v -run TestSessionNote
go test ./testing/test/... -v -run TestTeacherController_SessionNote
```
