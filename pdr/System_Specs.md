# 系統規格細節 (System Specifications)

本文件補充 API 與 業務邏輯 中未詳列的技術細節，包含 **錯誤碼標準**、**通知文案** 以及 **方案限制邏輯**。

---

## 1. 標準錯誤碼 (Error Codes)
本專案採用 **集中式錯誤管理**。所有錯誤碼需在後端 `global/errInfos` 套件中統一宣告，禁止散落在各個 Controller。

### 1.1 通用類 (Common)
| Code | Message | 說明 |
|:---|:---|:---|
| `SUCCESS` | "Operation successful" | 成功 (HTTP 200) |
| `E_BAD_REQUEST` | "Invalid parameters" | 參數格式錯誤 (HTTP 400) |
| `E_INTERNAL_ERROR`| "Internal server error" | 系統內部錯誤 (HTTP 500) |
| `E_DB_ERROR` | "Database operation failed" | 資料庫存取失敗 |

### 1.2 權限與認證 (Auth)
| Code | Message | 說明 |
|:---|:---|:---|
| `E_UNAUTHORIZED` | "Please login first" | 未登入/Token 無效 (HTTP 401) |
| `E_FORBIDDEN` | "Permission denied" | 權限不足 (HTTP 403) |
| `E_TOKEN_EXPIRED` | "Token expired" | Token 過期，需 Refresh |

### 1.3 業務資源類 (Resource)
| Code | Message | 說明 |
|:---|:---|:---|
| `E_NOT_FOUND` | "Resource not found" | 找不到中心/老師/課程 |
| `E_DUPLICATE` | "Resource already exists" | Email/LINE ID 已註冊 |
| `E_INVITE_INVALID`| "Invalid or expired invite"| 邀請碼無效或已過期 |
| `E_TAG_INVALID` | "Invalid hashtag format" | 標籤格式錯誤 (如長度超限) |
| `E_LIMIT_EXCEEDED`| "Plan limit reached" | 超過方案配額 (老師數/規則數) |

### 1.4 排課核心類 (Scheduling)
| Code | Message | 前端 UI 行為 |
|:---|:---|:---|
| `E_SCHED_OVERLAP` | "Time slot occupied" | 🔴 紅框：時段被佔用 |
| `E_SCHED_BUFFER` | "Insufficient buffer time" | 🟠 橘框：教室清潔/老師轉場不足 |
| `E_SCHED_PAST` | "Cannot book past time" | 🚫 阻擋：不能排過去的時間 |
| `E_SCHED_LOCKED` | "Slot is locked by another" | ⏳ 提示：併發鎖定中 |
| `E_SCHED_CLOSED` | "Center is closed" | 🚫 阻擋：非營業時間 |

---

## 2. LINE 通知文案模板 (Notification Templates)
需使用 LINE Messaging API (Push/Reply) 發送。支援 Flex Message 樣式為佳。

### 2.1 給老師 (To Teacher)
*   **排課成功通知 (New Schedule)**
    *   **Title**: 📅 新課程通知
    *   **Body**: "中心 [CenterName] 已為您安排新的課程。"
    *   **Fields**: 課程名稱、時間、地點。
    *   **Action**: [查看課表]

*   **異動審核結果 (Exception Result)**
    *   **Title**: ✅ 請假核准 (或 ❌ 被拒絕)
    *   **Body**: "您申請的 [Date] [Task] 異動已有結果。"
    *   **Reason**: (若是拒絕，顯示拒絕原因)
    *   **Reason**: (若是拒絕，顯示拒絕原因)
    *   **Action**: [查看詳情]

*   **明日課程提醒 (Daily Reminder)** - *New*
    *   **Trigger**: 系統每日 20:00 自動推播
    *   **Title**: ⏰ 明日課表提醒
    *   **Body**: "Hi [Teacher], 您明天 (1/21) 總共有 3 堂課，第一堂於 10:00 在 [CenterA] 開始。"
    *   **Action**: [查看完整課表]

*   **中心邀請通知 (Center Invitation)** - *New*
    *   **Trigger**: 管理員發送邀請時
    *   **Title**: 🤝 [CenterName] 邀請您加入
    *   **Body**: "您好，我們想邀請您成為本中心老師。聯絡人：[ContactInfo]。"
    *   **Action**: [接受邀請 並 查看中心資訊]

*   **課表緊急異動 (Urgent Change)** - *New/Supplement*
    *   **Trigger**: 管理員強制 Override 修改課程時
    *   **Title**: ⚠️ 課程異動通知
    *   **Body**: "您的 [Date] [Time] 課程已被管理員調整，請確認。"
    *   **Action**: [確認異動]

### 2.2 給中心管理員 (To Admin)
*   **異動申請通知 (New Exception)**
    *   **Title**: 📩 收到請假申請
    *   **Body**: "[TeacherName] 申請 [Date] 的課程改期。"
    *   **Conflict**: (若有 Buffer 衝突，直接顯示 "⚠️ 有輕微衝突")
    *   **Action**: [前往審核]

---

## 3. 方案限制執行邏輯 (Plan Enforcement)
在執行「寫入」操作前的 Middleware 檢查層。

### 3.1 限制項目矩陣
| 項目 | Free | Starter | Pro | Team |
|:---|:---:|:---:|:---:|:---:|
| **老師數 (Members)** | N/A | 5 | 20 | ∞ |
| **Grid 規則數 (Rules)** | N/A | 50 | 300 | ∞ |
| **智慧媒合 (Smart Match)**| ❌ | ❌ | ✅ | ✅ |
| **匯出隱私模式** | ❌ | ✅ | ✅ | ✅ |

### 3.2 檢查時機 (Check Points)
1.  **邀請老師時 (Invite Member)**:
    *   `Count(center.members) >= Plan.limit` ? Throw `E_LIMIT_TEACHERS` : Pass
2.  **儲存排課規則時 (Save Rules)**:
    *   `Count(center.rules) + NewRules.length >= Plan.limit` ? Throw `E_LIMIT_RULES` : Pass
3.  **使用進階功能時**:
    *   API 檢查 `center.plan_level`，若等級不足直接 Reject。

---

## 4. 系統環境與維運設定 (System Config)
*   **Audit Log Retention**: 預設保留 90 天，過期自動封存或刪除 (以符合 GDPA/個資法精神)。
*   **LINE Token Refresh**: 需排程 Job 定期更新 Channel Access Token (若非使用 Long-lived)。
*   **DB Backup**: 每日凌晨 04:00 執行全備份，保留 7 天。
