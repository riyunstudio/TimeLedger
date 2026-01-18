
# MASTER_PROMPT_TIMELEDGER_MYSQL.md
## TimeLedger｜一鍵式 AI 開發主控文件（MySQL 版）

> 本文件等同於「專案憲法」。
> 任何 AI（OpenCode / Claude Code / Cursor）或工程師 **只能依此文件與既有 specs 開發，不得自行發明需求**。

---

## 0. 使用方式（給你 & AI）

### 對你（專案擁有者）
- 把本檔案放在 repo 根目錄
- 指令給 AI：
  > 「**請嚴格依照 MASTER_PROMPT_TIMELEDGER_MYSQL_FULL.md 開發，不要自行假設需求**」

### 對 AI（必讀）
- 本文件優先權 **高於你自己的最佳實務**
- 若發現規格衝突，**停止並回報，不可自行決定**

---

## 1. 專案基本資訊

### 專案名稱
**TimeLedger**

### 一句話定義
> Teacher–Center Schedule Governance System

### 解決的唯一問題
> 在尊重老師時間主權的前提下，
> 讓中心能「安全、可追蹤、可審核」地完成排課。

---

## 2. 明確不做（Hard No）

以下功能 **在任何情況下都不得出現於程式碼中**：

- 學生報名 / 點名 / 學生帳號
- 金流 / 訂閱支付（僅做 trial / plan gating）
- LMS / 教材 / 作業
- App（僅 RWD Web）
- 自動排課演算法
- 教室 / 房間 buffer（MVP 僅做時間 overlap）

---

## 3. 角色定義（RBAC 簡化版）

### Teacher（免費）
- 管理個人課表（包含非任何中心的行程）
- 僅能看到自己行程 + 被授權中心的排課結果

### Center Admin（付費 / 試用）
- 管理 Center
- 建立 Offering / Rule
- 邀請老師加入中心
- 建立與審核例外（Exception）

---

## 4. 技術堆疊（強制）

### Backend
- 語言：Go
- Framework：沿用既有專案框架
- 架構分層（不可更動）：
  ```
  /app/controllers
  /app/requests
  /app/services
  /app/repositories
  /app/resources
  /app/models
  /global/errInfos
  ```

### Database
- **MySQL 8.x**
- InnoDB
- utf8mb4
- 所有欄位必須有 COMMENT
- 所有關聯查詢必須有 INDEX

### Frontend
- Web RWD
- 不做拖曳、不做九宮格
- List + Form + Modal 為主

---

## 5. 規格來源（Single Source of Truth）

AI **只能** 依照以下文件實作：

- `backend_specs/PRD.md`
- `backend_specs/API.md`
- `backend_specs/UIUX.md`
- `backend_specs/DB_SCHEMA.md`（MySQL）
- `docs/validate_pseudocode.txt`
- `teacher_schedule_dev_docs/*`
- `PRICING_AND_GO_TO_MARKET.md`

若文件間衝突：
1. MASTER_PROMPT
2. PRD
3. API
4. DB_SCHEMA
5. UIUX

---

## 6. 資料庫設計原則（MySQL）

### 必要表
- teachers
- teacher_personal_sessions
- centers
- teacher_center_memberships
- offerings
- schedule_rules
- schedule_exceptions
- center_plans
- audit_logs
- event_logs

### 時間欄位
- 一律使用 `DATETIME(3)`
- 不使用 timestamp timezone 轉換

---

## 7. API 實作強制規則

### 7.1 統一回傳格式（不可更動）
```json
{
  "code": 0,
  "message": "OK",
  "datas": {}
}
```

### 7.2 錯誤處理
- 所有業務錯誤需有明確 error code
- validate 衝突必須回傳可顯示的 conflict list

### 7.3 Center 隔離（Hard Rule）
- 所有 Admin API 必須帶入 center_id
- Repository 層必須強制 center scope

---

## 8. 核心業務邏輯（必須完整）

### 8.1 Schedule View（查詢時計算）
1. Rule 投影到日期區間
2. 套用 APPROVED Exception
3. 回傳最終 Session list

### 8.2 Validate（排課前檢查）
必須檢查：
- Teacher 個人行程 overlap
- Center rule overlap
- 同日是否已有有效 Exception

Overlap 定義：
```
(aStart < bEnd) && (bStart < aEnd)
```

### 8.3 Exceptions
- CANCEL：取消當日 session
- RESCHEDULE：取消原日 + 新增一堂
- approve 時必須重新 validate

---

## 9. Trial / Plan Gating（必做）

### Trial
- 建立 Center → 自動給 30 天 Pro 試用

### 到期
- 到期後：
  - Read OK
  - Write 全拒（Rule / Exception / Membership）

### 方案限制
| Plan | Max Teachers |
|-----|--------------|
| Starter | 5 |
| Pro | 15 |
| Team | 無上限 |

---

## 10. Audit & Event（必做）

### audit_logs
- who / what / when
- entity + payload_json

### event_logs（至少）
- trial_started
- teacher_invited
- membership_activated
- rule_created
- validate_conflict
- exception_created
- exception_approved
- trial_expired

---

## 11. 開發順序（不可跳關）

### Stage 1
- DB migration（MySQL）
- Teacher APIs

### Stage 2
- Center / Membership / Offering

### Stage 3
- Rule + Validate

### Stage 4
- Exception + Schedule View

### Stage 5
- Trial gating + Frontend RWD

---

## 12. 最終交付要求

每個 Stage 必須提供：
1. 修改檔案清單
2. Migration 指令
3. 測試方式（curl）
4. 已知限制

---

## 13. 結語（給 AI）

> 你不是在寫一個 CRUD 系統，  
> 你是在實作「時間治理」。

任何讓時間失去可信度的捷徑，都是錯的。
