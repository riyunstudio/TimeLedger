# TimeLedger 資料庫設計 (MySQL)

## 1. ER 模型概觀
系統以 `centers` 為核心租戶隔離單位。老師 (`teachers`) 為跨租戶實體，可加入多個中心。
個人行程 (`personal_events`) 獨立於中心之外，歸屬於老師個人。

---

## 2. 核心資料表 (Tables)

### 2.1 租戶與用戶 (Auth & Tenant)

#### `centers` (中心)
| Column | Type | Description |
|:--- |:--- |:--- |
| id | BIGINT PK | 中心 ID |
| name | VARCHAR | 中心名稱 |
| plan_level | ENUM | 方案等級 (FREE, STARTER, PRO, TEAM) |
| settings | JSON | 政策設定 (如：開放註冊、預設語系等，不再包含課程緩衝) |
| created_at | DATETIME | |

#### `admin_users` (後台管理員)
| Column | Type | Description |
|:--- |:--- |:--- |
| id | BIGINT PK | |
| center_id | BIGINT FK | 所屬中心 (Index) |
| email | VARCHAR | 登入帳號 (Unique under center) |
| password_hash | VARCHAR | |
| **name** | VARCHAR | 管理員姓名 (顯示用) |
| role | ENUM | OWNER, ADMIN, STAFF |
| **status** | ENUM | ACTIVE, INACTIVE |

#### `teachers` (老師/個人用戶)
| Column | Type | Description |
|:--- |:--- |:--- |
| id | BIGINT PK | |
| **line_user_id** | VARCHAR | LINE OpenID (Index, Unique) |
| name | VARCHAR | 显示名稱 |
| email | VARCHAR | 備用/通知 Email |
| avatar_url | VARCHAR | |
| **bio** | TEXT | 老師自我介紹 (公開) |
| **is_open_to_hiring** | BOOLEAN | 是否開放媒合 (開啟後，未加入的中心也能搜尋到) |
| **region** | VARCHAR | 服務地區 (e.g. '台北市大安區', '台中市西屯區') |
| **public_contact_info** | TEXT | 自訂公開聯繫資訊 (Line/Phone/Email 等) |

#### `teacher_skills` (老師專業能力 - 雙層結構)
| Column | Type | Description |
|:--- |:--- |:--- |
| id | BIGINT PK | |
| teacher_id | BIGINT FK | Index. 老師可擁有「多項」專業能力 (如：皮拉提斯 + 重訓)。 |
| **category** | VARCHAR | 大類 (e.g. '有氧', '瑜伽', '舞蹈'...) |
| **skill_name** | VARCHAR | 子類/具體描述 |
| level | ENUM | BASIC, INTERMEDIATE, ADVANCED |

#### `hashtags` (Hashtag 字典表)
存放系統中所有被使用過的標籤。
| Column | Type | Description |
|:--- |:--- |:--- |
| id | BIGINT PK | |
| **name** | VARCHAR | 標籤名稱 (Unique, e.g. '#街舞') |
| **usage_count** | INT | 使用次數 (效能優化，用於清理判斷) |

#### `teacher_skill_hashtags` (老師技能-標籤 關聯表)
| Column | Type | Description |
|:--- |:--- |:--- |
| teacher_skill_id | BIGINT FK | Index. 關聯至特定的老師技能。 |
| hashtag_id | BIGINT FK | Index |

#### `teacher_personal_hashtags` (老師個人品牌標籤)
用於個人課表匯出展示，具備品牌辨識度。
| Column | Type | Description |
|:--- |:--- |:--- |
| teacher_id | BIGINT FK | Index. 關聯老師個人。 |
| hashtag_id | BIGINT FK | Index |
| **sort_order** | TINYINT | 排序 (1-5) |

#### `teacher_certificates` (老師證照/證明文件)
| Column | Type | Description |
|:--- |:--- |:--- |
| id | BIGINT PK | |
| teacher_id | BIGINT FK | Index |
| name | VARCHAR | 證照名稱 (e.g. 'Yamaha 鋼琴檢定 5 級') |
| file_url | VARCHAR | S3/Storage 連結 |
| issued_at | DATE | 發證日期 |
| created_at | DATETIME | |

#### `center_teacher_notes` (中心對老師的私密備註/評價)
| Column | Type | Description |
|:--- |:--- |:--- |
| id | BIGINT PK | |
| center_id | BIGINT FK | Index |
| teacher_id | BIGINT FK | |
| internal_note | TEXT | 中心內部的備註 (如：配合度高、適合教幼兒) |
| rating | TINYINT | 中心內部評分 (1-5) |

#### `center_memberships` (老師-中心關聯)
| Column | Type | Description |
|:--- |:--- |:--- |
| id | BIGINT PK | |
| center_id | BIGINT FK | |
| teacher_id | BIGINT FK | |
| status | ENUM | ACTIVE, INACTIVE, INVITED |

#### `center_invitations` (邀請函暫存)
追蹤邀請狀態，避免重複發送或過期。
| Column | Type | Description |
|:--- |:--- |:--- |
| id | BIGINT PK | |
| center_id | BIGINT FK | |
| email | VARCHAR | 受邀者 Email (Nullable if via LINE) |
| token | VARCHAR | 唯一邀請碼 (UUID) |
| status | ENUM | PENDING, ACCEPTED, EXPIRED |
| created_at | DATETIME | |
| expires_at | DATETIME | 預設 7 天過期 |

---

### 2.2 排課核心 (Scheduling Core)

#### `courses` (課程模板)
| Column | Type | Description |
|:--- |:--- |:--- |
| id | BIGINT PK | |
| center_id | BIGINT FK | |
| name | VARCHAR | |
| default_duration | INT | 分鐘數 |
| color_hex | VARCHAR | 顯示顏色 |
| **room_buffer_min** | INT | 該課程所需教室緩衝 (清潔時間) |
| **teacher_buffer_min** | INT | 該課程所需老師緩衝 (轉場時間) |

#### `offerings` (開課班別)
| Column | Type | Description |
|:--- |:--- |:--- |
| id | BIGINT PK | |
| center_id | BIGINT FK | |
| course_id | BIGINT FK | |
| default_room_id | BIGINT FK | 預設教室 (可 Null) |
| default_teacher_id| BIGINT FK | 預設老師 (可 Null) |
| **allow_buffer_override** | BOOLEAN | 是否允許無視緩衝強制排入 (特定班別權限) |

#### `timetable_templates` (課表模板)
定義 Grid 的行列結構。
| Column | Type | Description |
|:--- |:--- |:--- |
| id | BIGINT PK | |
| center_id | BIGINT FK | |
| name | VARCHAR | 模板名稱 (例：'音樂教室 A 週末板') |
| row_type | ENUM | ROOM, TEACHER (定義 Row 代表的是教室還是老師) |

#### `timetable_cells` (模板格子)
預先定義好的排課位置。
| Column | Type | Description |
|:--- |:--- |:--- |
| id | BIGINT PK | |
| template_id | BIGINT FK | |
| row_no | INT | 橫向編號 |
| col_no | INT | 縱向編號 (通常對應時間段) |
| start_time | TIME | |
| end_time | TIME | |
| room_id | BIGINT | 預設教室 (若 row_type=TEACHER) |
| teacher_id | BIGINT | **Nullable**. 預設老師 (若 row_type=ROOM) |

#### `schedule_rules` (排課規則 - 週期性)
定義中心正規課程的週期。
| Column | Type | Description |
|:--- |:--- |:--- |
| id | BIGINT PK | |
| center_id | BIGINT FK | Index |
| offering_id | BIGINT FK | |
| teacher_id | BIGINT FK | **Nullable**.所屬老師 (Index)。支援暫不指派老師之排課。 |
| room_id | BIGINT FK | Index |
| **weekday** | TINYINT | 1=Mon, 7=Sun |
| start_time | TIME | '09:00:00' |
| end_time | TIME | '10:30:00' |
| effective_range | DATERANGE | [start_date, end_date] |

#### `schedule_exceptions` (例外管理)
處理停課、改期。
| Column | Type | Description |
|:--- |:--- |:--- |
| id | BIGINT PK | |
| center_id | BIGINT FK | |
| **rule_id** | BIGINT FK | 關聯原本的 Rule |
| **original_date** | DATE | 原本發生的日期 (Index) |
| type | ENUM | **CANCEL**, **RESCHEDULE**, **ADD** |
| status | ENUM | **PENDING**, **APPROVED**, **REJECTED** |
| new_start_at | DATETIME | 改期後的新開始時間 (Nullable) |
| new_end_at | DATETIME | 改期後的新結束時間 (Nullable) |
| **new_teacher_id** | BIGINT | 改期後或新指派的老師 (Nullable) |
| new_room_id | BIGINT | 改期後教室 |
| reason | VARCHAR | 申請/審核原因 |

#### `personal_events` (個人行程 - 免費版功能)
老師個人的私有行程，不歸屬任何中心，但需在「空檔查找」時被納入 Teacher Buffer 計算。
| Column | Type | Description |
|:--- |:--- |:--- |
| id | BIGINT PK | |
| **teacher_id** | BIGINT FK | Index |
| title | VARCHAR | |
| start_at | DATETIME | |
| end_at | DATETIME | |
| **recurrence_rule** | JSON | 循環規則 `{type: 'WEEKLY', interval: 2, days: [1,3]}` |
| is_all_day | BOOLEAN | |
| **recurrence_rule** | JSON | 循環規則 `{type: 'WEEKLY', interval: 2, days: [1,3]}` |

---

### 2.4 教學紀錄 (Teaching Records)

#### `session_notes` (課程筆記)
用於老師針對「特定日期」的課程撰寫備課或教學紀錄。
| Column | Type | Description |
|:--- |:--- |:--- |
| id | BIGINT PK | |
| **teacher_id** | BIGINT FK | 撰寫者 |
| **rule_id** | BIGINT FK | 關聯的中心排課規則 (Nullable) |
| **personal_event_id** | BIGINT FK | 關聯的個人行程 (Nullable) |
| **session_date** | DATE | 課程日期 (Index) |
| **content** | TEXT | 上課內容紀錄 (Past) |
| **prep_note** | TEXT | 備課筆記 (Future) |
| updated_at | DATETIME | |

---

### 2.3 治理與紀錄 (Governance)

#### `audit_logs` (操作審計)
| Column | Type | Description |
|:--- |:--- |:--- |
| id | BIGINT PK | |
| center_id | BIGINT FK | |
| actor_type | ENUM | ADMIN, TEACHER |
| actor_id | BIGINT | |
| action | VARCHAR | 'CREATE_RULE', 'APPROVE_EXCEPTION' |
| target_type | VARCHAR | |
| target_id | BIGINT | |
| payload | JSON | 修改前內容/修改後內容 |
| timestamp | DATETIME | |

---

## 3. 關鍵索引優化 (Indexes)

1. **衝突檢測 (Overlap Check)**
   - `schedule_rules`: `(center_id, weekday, start_time)` - 快速查找某中心某天的規則。
   - `schedule_exceptions`: `(rule_id, original_date)` - 快速確認某 Rule 在某天是否有例外。
   - `personal_events`: `(teacher_id, start_at, end_at)` - 快速查找老師私人行程。

2. **多租戶查詢**
   - 所有查詢 `WHERE` 子句第一條件必然是 `center_id` (除個人行程外)。
    - 因此大部分 Index 應以 `center_id` 為 Prefix。

    - `session_notes`: `UNIQUE(rule_id, session_date)` - 確保單一課堂當日只有一份筆記。
    - `center_invitations`: `INDEX(token)` - 加速邀請碼驗證。
    - `teachers`: `INDEX(region)` - 支持地區篩選。

## 4. 併發與資料一致性 (Concurrency Control)
針對「多中心同時排同一位老師」的競爭條件 (Race Condition)：
- **風險**: Center A 與 Center B 同時讀取到 Teacher T 週一 10:00 空閒，並同時送出排課請求，導致雙重預約 (Double Booking)。
- **解決方案 (Locking Strategy)**:
  - **Option A (Redis Lock)**: 在 Controller 層進行 Validate/Write 前，取得 `LOCK:TEACHER:{id}` 分散式鎖 (TTL 5秒)。
  - **Option B (DB Row Lock)**: 在 Transaction 開始時，執行 `SELECT id FROM teachers WHERE id = ? FOR UPDATE`。這會強制序列化針對該老師的所有寫入交易。
  - **建議**: 採用 **Option B (DB Row Lock)**，實作最穩健且不需額外引入 Redis 元件。

## 5. SQL 範例：衝突檢查 (Pseudo)
```sql
-- 檢查某老師在特定時段是否有 '中心規則' 衝突
SELECT * FROM schedule_rules
WHERE teacher_id = ? 
  AND weekday = WEEKDAY(?) 
  AND start_time < ? AND end_time > ?
  AND id NOT IN (SELECT rule_id FROM schedule_exceptions WHERE original_date = ? AND type = 'CANCEL')
-- 需再加上 Exception 的 RESCHEDULE 檢查

---

## 6. 快取層設計 (Cache Layer - Redis)
為了提升系統反應速度與安全性，系統將使用 Redis 處理以下場景：

1. **Session & Auth**:
   - 存放 JWT Blacklist (用於登出或強制失效)。
   - 暫存老師登入時的 LINE Profile (TTL 30分)，減少對 LINE API 的請求。
2. **Hashtag Autocomplete**:
   - 預熱常用標籤 (Top 100) 至 Redis ZSET，支援快速模糊搜尋。
3. **Validate Engine 暫存**:
   - 針對高頻率的 `validate` 請求，快取短暫的 Query 結果 (TTL 1-2秒)，避免重複讀取大型 Rules 表。
```
