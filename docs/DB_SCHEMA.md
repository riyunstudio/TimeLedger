# TimeLedger 資料庫 Schema

## 1. 資料表總覽

| 資料表 | 說明 |
|:---|:---|
| `centers` | 中心資訊 |
| `admin_users` | 管理員帳號 |
| `teachers` | 教師帳號 |
| `center_memberships` | 教師-中心關聯 |
| `center_invitations` | 邀請記錄 |
| `geo_cities` | 縣市資料 |
| `geo_districts` | 區域資料 |
| `courses` | 課程模板 |
| `offerings` | 班別 |
| `rooms` | 教室 |
| `timetable_templates` | 課表範本 |
| `timetable_cells` | 課表範本細胞 |
| `schedule_rules` | 排課規則 |
| `schedule_exceptions` | 例外記錄 |
| `personal_events` | 私人行程 |
| `teacher_skills` | 教師技能 |
| `hashtags` | 標籤字典 |
| `teacher_skill_hashtags` | 技能-標籤關聯 |
| `teacher_personal_hashtags` | 個人標籤 |
| `teacher_certificates` | 教師證照 |
| `teacher_backgrounds` | 教師背景 |
| `center_teacher_notes` | 中心對教師的備註 |
| `center_holidays` | 中心假期 |
| `center_terms` | 學期 |
| `session_notes` | 課程筆記 |
| `audit_logs` | 審計日誌 |
| `notifications` | 通知 |
| `notification_queues` | 通知佇列 |

---

## 2. 資料表詳解

### 2.1 centers (中心)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `name` | VARCHAR(255) NOT NULL | 中心名稱 |
| `plan_level` | VARCHAR(20) DEFAULT 'FREE' | 方案等級 |
| `settings` | JSON | 設定 (allow_public_register, exception_lead_days, etc.) |
| `created_at` | DATETIME NOT NULL | 建立時間 |

### 2.2 admin_users (管理員)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `center_id` | BIGINT UNSIGNED NOT NULL | 所屬中心 |
| `email` | VARCHAR(255) NOT NULL | Email |
| `password_hash` | VARCHAR(255) NOT NULL | 密碼雜湊 |
| `name` | VARCHAR(255) NOT NULL | 姓名 |
| `role` | VARCHAR(20) DEFAULT 'STAFF' | 角色 (OWNER/ADMIN/STAFF) |
| `status` | VARCHAR(20) DEFAULT 'ACTIVE' | 狀態 |
| `line_user_id` | VARCHAR(64) | LINE 用戶 ID |
| `line_binding_code` | VARCHAR(8) | LINE 綁定驗證碼 |
| `line_binding_expires` | DATETIME | 綁定驗證碼過期時間 |
| `line_notify_enabled` | BOOLEAN DEFAULT TRUE | 接收 LINE 通知 |
| `line_bound_at` | DATETIME | LINE 綁定時間 |
| `created_at` | DATETIME NOT NULL | 建立時間 |
| `updated_at` | DATETIME NOT NULL | 更新時間 |
| `deleted_at` | DATETIME (index) | 軟刪除 |

### 2.3 teachers (教師)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `line_user_id` | VARCHAR(255) index | LINE 用戶 ID |
| `line_notify_token` | VARCHAR(255) | LINE 推播 Token |
| `name` | VARCHAR(255) NOT NULL | 姓名 |
| `email` | VARCHAR(255) | Email |
| `avatar_url` | VARCHAR(512) | 頭貼 URL |
| `bio` | TEXT | 個人簡介 |
| `is_open_to_hiring` | BOOLEAN DEFAULT FALSE | 開放職涯搜尋 |
| `is_placeholder` | BOOLEAN DEFAULT FALSE | 是否為暫存教師 |
| `city` | VARCHAR(100) index | 縣市 |
| `district` | VARCHAR(100) | 區域 |
| `public_contact_info` | TEXT | 公開聯絡資訊 |
| `created_at` | DATETIME NOT NULL | 建立時間 |
| `updated_at` | DATETIME NOT NULL | 更新時間 |
| `deleted_at` | DATETIME (index) | 軟刪除 |

### 2.4 center_memberships (教師-中心關聯)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `center_id` | BIGINT UNSIGNED NOT NULL | 中心 ID |
| `teacher_id` | BIGINT UNSIGNED NOT NULL | 教師 ID |
| `role` | VARCHAR(20) DEFAULT 'TEACHER' | 角色 |
| `status` | VARCHAR(20) DEFAULT 'INVITED' | 狀態 (INVITED/ACTIVE) |
| `created_at` | DATETIME NOT NULL | 建立時間 |
| `updated_at` | DATETIME NOT NULL | 更新時間 |
| `deleted_at` | DATETIME (index) | 軟刪除 |

### 2.5 center_invitations (邀請記錄)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `center_id` | BIGINT UNSIGNED NOT NULL | 中心 ID |
| `teacher_id` | BIGINT UNSIGNED | 教師 ID |
| `invited_by` | BIGINT UNSIGNED NOT NULL | 邀請人 |
| `token` | VARCHAR(64) NOT NULL | 邀請 Token |
| `email` | VARCHAR(255) | 邀請 Email |
| `status` | VARCHAR(20) DEFAULT 'PENDING' | 狀態 |
| `expires_at` | DATETIME NOT NULL | 過期時間 |
| `accepted_at` | DATETIME | 接受時間 |
| `created_at` | DATETIME NOT NULL | 建立時間 |
| `updated_at` | DATETIME NOT NULL | 更新時間 |

### 2.6 geo_cities (縣市)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `name` | VARCHAR(100) NOT NULL | 名稱 |
| `code` | VARCHAR(10) NOT NULL | 代碼 |

### 2.7 geo_districts (區域)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `city_id` | BIGINT UNSIGNED NOT NULL | 縣市 ID |
| `name` | VARCHAR(100) NOT NULL | 名稱 |
| `code` | VARCHAR(10) NOT NULL | 代碼 |

### 2.8 courses (課程模板)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `center_id` | BIGINT UNSIGNED NOT NULL | 中心 ID |
| `name` | VARCHAR(255) NOT NULL | 課程名稱 |
| `description` | TEXT | 說明 |
| `default_duration` | INT DEFAULT 60 | 預設時長(分鐘) |
| `room_buffer_min` | INT DEFAULT 0 | 教室緩衝時間 |
| `teacher_buffer_min` | INT DEFAULT 0 | 老師緩衝時間 |
| `is_active` | BOOLEAN DEFAULT TRUE | 是否啟用 |
| `created_at` | DATETIME NOT NULL | 建立時間 |
| `updated_at` | DATETIME NOT NULL | 更新時間 |

### 2.9 offerings (班別)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `center_id` | BIGINT UNSIGNED NOT NULL | 中心 ID |
| `course_id` | BIGINT UNSIGNED NOT NULL | 課程 ID |
| `name` | VARCHAR(255) NOT NULL | 班別名稱 |
| `teacher_id` | BIGINT UNSIGNED | 帶課老師 |
| `room_id` | BIGINT UNSIGNED | 教室 |
| `start_date` | DATE | 開始日期 |
| `end_date` | DATE | 結束日期 |
| `day_of_week` | INT | 星期 (0-6) |
| `start_time` | VARCHAR(10) | 開始時間 |
| `end_time` | VARCHAR(10) | 結束時間 |
| `price` | DECIMAL(10,2) | 收費 |
| `max_students` | INT | 最大人數 |
| `is_active` | BOOLEAN DEFAULT TRUE | 是否啟用 |
| `allow_buffer_override` | BOOLEAN DEFAULT FALSE | 允許覆寫緩衝 |
| `created_at` | DATETIME NOT NULL | 建立時間 |
| `updated_at` | DATETIME NOT NULL | 更新時間 |

### 2.10 rooms (教室)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `center_id` | BIGINT UNSIGNED NOT NULL | 中心 ID |
| `name` | VARCHAR(255) NOT NULL | 教室名稱 |
| `capacity` | INT DEFAULT 10 | 容納人數 |
| `is_active` | BOOLEAN DEFAULT TRUE | 是否啟用 |
| `created_at` | DATETIME NOT NULL | 建立時間 |
| `updated_at` | DATETIME NOT NULL | 更新時間 |

### 2.11 timetable_templates (課表範本)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `center_id` | BIGINT UNSIGNED NOT NULL | 中心 ID |
| `name` | VARCHAR(255) NOT NULL | 範本名稱 |
| `created_at` | DATETIME NOT NULL | 建立時間 |
| `updated_at` | DATETIME NOT NULL | 更新時間 |

### 2.12 timetable_cells (課表範本細胞)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `template_id` | BIGINT UNSIGNED NOT NULL | 範本 ID |
| `row` | INT NOT NULL | 列 (時間) |
| `col` | INT NOT NULL | 行 (星期) |
| `course_id` | BIGINT UNSIGNED | 課程 ID |
| `color` | VARCHAR(7) | 顏色 |
| `created_at` | DATETIME NOT NULL | 建立時間 |

### 2.13 schedule_rules (排課規則)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `center_id` | BIGINT UNSIGNED NOT NULL | 中心 ID |
| `offering_id` | BIGINT UNSIGNED NOT NULL | 班別 ID |
| `teacher_id` | BIGINT UNSIGNED | 老師 ID |
| `room_id` | BIGINT UNSIGNED NOT NULL | 教室 ID |
| `name` | VARCHAR(100) | 規則名稱 |
| `code` | VARCHAR(50) | 規則代碼 |
| `weekday` | TINYINT NOT NULL | 星期 (0-6) |
| `start_time` | VARCHAR(10) NOT NULL | 開始時間 |
| `end_time` | VARCHAR(10) NOT NULL | 結束時間 |
| `duration` | INT DEFAULT 60 | 時長(分鐘) |
| `is_cross_day` | BOOLEAN DEFAULT FALSE | 是否跨日 |
| `skip_holiday` | BOOLEAN DEFAULT TRUE | 跳過假日 |
| `effective_range` | JSON NOT NULL | 有效日期範圍 |
| `suspended_dates` | JSON | 暫停日期 |
| `status` | VARCHAR(20) DEFAULT 'CONFIRMED' | 狀態 |
| `lock_at` | DATETIME | 鎖定時間 |
| `created_at` | DATETIME NOT NULL | 建立時間 |
| `updated_at` | DATETIME NOT NULL | 更新時間 |
| `deleted_at` | DATETIME (index) | 軟刪除 |

**effective_range 結構：**
```json
{
  "start_date": "2026-01-01 00:00:00",
  "end_date": "2026-12-31 23:59:59"
}
```

### 2.14 schedule_exceptions (例外記錄)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `center_id` | BIGINT UNSIGNED NOT NULL | 中心 ID |
| `rule_id` | BIGINT UNSIGNED NOT NULL | 規則 ID |
| `original_date` | DATE NOT NULL | 原日期 |
| `exception_type` | VARCHAR(20) NOT NULL | 類型 (LEAVE/RESCHEDULE/SWAP/CANCEL) |
| `status` | VARCHAR(20) DEFAULT 'PENDING' | 狀態 |
| `new_start_at` | DATETIME | 新開始時間 |
| `new_end_at` | DATETIME | 新結束時間 |
| `new_teacher_id` | BIGINT UNSIGNED | 新老師 ID |
| `new_room_id` | BIGINT UNSIGNED | 新教室 ID |
| `reason` | TEXT | 原因 |
| `reviewed_by` | BIGINT UNSIGNED | 審核人 |
| `reviewed_at` | DATETIME | 審核時間 |
| `review_note` | TEXT | 審核備註 |
| `created_at` | DATETIME NOT NULL | 建立時間 |
| `updated_at` | DATETIME NOT NULL | 更新時間 |

### 2.15 personal_events (私人行程)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `teacher_id` | BIGINT UNSIGNED NOT NULL | 教師 ID |
| `title` | VARCHAR(255) NOT NULL | 標題 |
| `start_at` | DATETIME NOT NULL | 開始時間 |
| `end_at` | DATETIME NOT NULL | 結束時間 |
| `recurrence_rule` | JSON | 循環規則 |
| `is_all_day` | BOOLEAN DEFAULT FALSE | 全天 |
| `color_hex` | VARCHAR(7) | 顏色 |
| `note` | TEXT | 備註 |
| `created_at` | DATETIME NOT NULL | 建立時間 |
| `updated_at` | DATETIME NOT NULL | 更新時間 |
| `deleted_at` | DATETIME (index) | 軟刪除 |

**recurrence_rule 結構：**
```json
{
  "type": "WEEKLY",
  "interval": 1,
  "weekdays": [1, 3, 5],
  "until": "2026-06-30"
}
```

### 2.16 teacher_skills (教師技能)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `teacher_id` | BIGINT UNSIGNED NOT NULL | 教師 ID |
| `name` | VARCHAR(255) NOT NULL | 技能名稱 |
| `level` | VARCHAR(20) | 程度 |
| `created_at` | DATETIME NOT NULL | 建立時間 |
| `updated_at` | DATETIME NOT NULL | 更新時間 |

### 2.17 hashtags (標籤字典)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `name` | VARCHAR(50) NOT NULL UNIQUE | 標籤名稱 |
| `usage_count` | INT DEFAULT 0 | 使用次數 |
| `created_at` | DATETIME NOT NULL | 建立時間 |

### 2.18 teacher_certificates (教師證照)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `teacher_id` | BIGINT UNSIGNED NOT NULL | 教師 ID |
| `name` | VARCHAR(255) NOT NULL | 證照名稱 |
| `issuer` | VARCHAR(255) | 發照單位 |
| `issued_date` | DATE | 發照日期 |
| `expiry_date` | DATE | 到期日期 |
| `file_url` | VARCHAR(512) | 檔案 URL |
| `created_at` | DATETIME NOT NULL | 建立時間 |
| `updated_at` | DATETIME NOT NULL | 更新時間 |

### 2.19 center_teacher_notes (教師備註)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `center_id` | BIGINT UNSIGNED NOT NULL | 中心 ID |
| `teacher_id` | BIGINT UNSIGNED NOT NULL | 教師 ID |
| `rating` | INT | 評分 (1-5) |
| `note` | TEXT | 備註內容 |
| `created_at` | DATETIME NOT NULL | 建立時間 |
| `updated_at` | DATETIME NOT NULL | 更新時間 |

### 2.20 center_holidays (中心假期)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `center_id` | BIGINT UNSIGNED NOT NULL | 中心 ID |
| `name` | VARCHAR(255) NOT NULL | 假日名稱 |
| `date` | DATE NOT NULL | 日期 |
| `is_recurring` | BOOLEAN DEFAULT FALSE | 是否循環 |
| `created_at` | DATETIME NOT NULL | 建立時間 |
| `updated_at` | DATETIME NOT NULL | 更新時間 |

### 2.21 center_terms (學期)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `center_id` | BIGINT UNSIGNED NOT NULL | 中心 ID |
| `name` | VARCHAR(255) NOT NULL | 學期名稱 |
| `start_date` | DATE NOT NULL | 開始日期 |
| `end_date` | DATE NOT NULL | 結束日期 |
| `is_active` | BOOLEAN DEFAULT FALSE | 是否啟用 |
| `created_at` | DATETIME NOT NULL | 建立時間 |
| `updated_at` | DATETIME NOT NULL | 更新時間 |

### 2.22 session_notes (課程筆記)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `teacher_id` | BIGINT UNSIGNED NOT NULL | 教師 ID |
| `session_date` | DATE NOT NULL | 課程日期 |
| `course_name` | VARCHAR(255) | 課程名稱 |
| `content` | TEXT NOT NULL | 筆記內容 |
| `created_at` | DATETIME NOT NULL | 建立時間 |
| `updated_at` | DATETIME NOT NULL | 更新時間 |

### 2.23 audit_logs (審計日誌)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `center_id` | BIGINT UNSIGNED NOT NULL | 中心 ID |
| `user_id` | BIGINT UNSIGNED NOT NULL | 操作者 ID |
| `user_type` | VARCHAR(20) NOT NULL | 操作者類型 |
| `action` | VARCHAR(50) NOT NULL | 動作 |
| `target_type` | VARCHAR(50) NOT NULL | 目標類型 |
| `target_id` | BIGINT UNSIGNED NOT NULL | 目標 ID |
| `details` | JSON | 詳情 |
| `created_at` | DATETIME NOT NULL | 建立時間 |

### 2.24 notifications (通知)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `user_id` | BIGINT UNSIGNED NOT NULL | 用戶 ID |
| `user_type` | VARCHAR(20) NOT NULL | 用戶類型 |
| `center_id` | BIGINT UNSIGNED | 中心 ID |
| `title` | VARCHAR(255) NOT NULL | 標題 |
| `message` | TEXT NOT NULL | 訊息內容 |
| `type` | VARCHAR(20) NOT NULL | 類型 |
| `is_read` | BOOLEAN DEFAULT FALSE | 已讀狀態 |
| `read_at` | DATETIME | 閱讀時間 |
| `created_at` | DATETIME NOT NULL | 建立時間 |
| `deleted_at` | DATETIME (index) | 軟刪除 |

### 2.25 notification_queues (通知佇列)

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `id` | BIGINT UNSIGNED PRIMARY KEY | ID |
| `notification_type` | VARCHAR(50) NOT NULL | 通知類型 |
| `recipient_id` | BIGINT UNSIGNED NOT NULL | 收件人 ID |
| `recipient_type` | VARCHAR(20) NOT NULL | 收件人類型 |
| `channel` | VARCHAR(20) NOT NULL | 管道 (LINE/PUSH) |
| `payload` | JSON NOT NULL | 內容 |
| `status` | VARCHAR(20) DEFAULT 'PENDING' | 狀態 |
| `retry_count` | INT DEFAULT 0 | 重試次數 |
| `error_message` | TEXT | 錯誤訊息 |
| `sent_at` | DATETIME | 發送時間 |
| `created_at` | DATETIME NOT NULL | 建立時間 |
| `updated_at` | DATETIME NOT NULL | 更新時間 |

---

## 3. 索引設計

### 3.1 重點索引

| 資料表 | 索引欄位 | 類型 | 說明 |
|:---|:---|:---|:---|
| `schedule_rules` | center_id, weekday, start_time | INDEX | 排課查詢 |
| `schedule_rules` | teacher_id | INDEX | 老師時段查詢 |
| `schedule_rules` | room_id | INDEX | 教室時段查詢 |
| `schedule_exceptions` | rule_id, original_date | INDEX | 例外查詢 |
| `personal_events` | teacher_id, start_at, end_at | INDEX | 私人行程查詢 |
| `center_memberships` | center_id, teacher_id | INDEX | 中心教師查詢 |
| `teachers` | line_user_id | INDEX | LINE 登入 |

---

## 4. 關聯圖

```
centers (1) ──────< (N) admin_users
centers (1) ──────< (N) teachers (via center_memberships)
centers (1) ──────< (N) courses
centers (1) ──────< (N) rooms
centers (1) ──────< (N) offerings
centers (1) ──────< (N) schedule_rules
centers (1) ──────< (N) schedule_exceptions

offerings (1) ───< (N) schedule_rules
teachers (1) ────< (N) teacher_skills
teachers (1) ────< (N) teacher_certificates
teachers (1) ────< (N) personal_events

schedule_rules (1) ──< (N) schedule_exceptions
```

---

*本文件基於實際程式碼生成，最後更新：2026-02-12*
