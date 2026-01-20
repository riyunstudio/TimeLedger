# PDR｜多中心課表與例外變動管理系統（TimeLedger）

> [!IMPORTANT]
> **Trademark & IP Protection Notice**
> **(c) 2026 TimeLedger Team. All rights reserved.**
版本：v3 FINAL（正式版）  
日期：2026-01-20  
狀態：可進入正式開發 / 商用上線準備  

> 本文件為「最新版且正式」PDR。
> 核心前提：**老師主動加入中心即授權共享時間**（因此不再處理隱私暴露模型）。

---

## 0. 產品定位
TimeLedger 是一套 **Web-only（RWD）** 的「多中心 × 老師」排課協作系統：
- **官方官網 (Landing Page)**：以 **WOW Demo** (免登入沙盒) 吸引用戶，展示排課治理的核心價值。
- **老師端 (免費)**：管理個人行程與可分享的可用時段，打造專業履歷 (Skills/Hashtags/Region)。
- **中心端 (付費)**：使用 **Grid Template** 高效排課，透過 **Smart Match (智慧媒合)** 與 **Talent Search (人才搜尋)** 快速開發師資。
- **治理核心**：以 **Buffer/Overlap/State Machine** 確保營運安全與邏輯一致性。

> Out of Scope：學生報名/點名/金流/LMS、自動生成完整課表（深度排課演算法）、App。

---

## 1. 目標使用者與角色
### 1.1 角色
- **Center Admin**：中心管理員/行政（排課、審核、治理、報表）
- **Teacher**：老師（查看課表、提出停課/改期）
- （選配）**Super Admin**：總控（僅營運支援，不影響核心流程）

### 1.2 核心場景
- Admin：建立模板 → 設定 cells → Grid 排課（bulk）→ validate → 成功/待審 → 審核
- Teacher：查看課表（週）→ 對某堂課提出取消/改期 → validate → 自動生效或待審

---

## 2. 系統邊界與多中心隔離
- **Center（中心/租戶）** 是資料隔離單位
- 所有寫入/查詢必須強制帶入 `center_id` scope（middleware + DB where）
- Teacher 可以加入多中心（membership）
- Teacher 的「共享」僅對 **已加入的中心** 生效

---

## 3. 名詞定義
- **Course**：課程模板（名稱/屬性）
- **Offering**：開課單位（實際排進課表的班別/教室/時長/狀態）
- **Template**：課表模板（rows/cols 的 grid）
- **Cell**：模板格子（row/col + 時間 +（可選）教室）
- **Rule**：週期排課規則（weekday + start/end + room + teacher + offering）
- **Session**：Rule 投影到指定日期後的實際課（含例外後狀態）
- **Exception**：例外事件（停課/改期/補課；可能需審核）
- **Policy**：中心治理政策（buffer / override / approval）

---

## 4. v3 核心能力（相較 MVP）
### 4.1 排課治理與人才動能
- **Hard Overlap（硬衝突）**：任何重疊直接拒絕。
- **Buffer（緩衝）**：Teacher Buffer (轉場) 與 Room Buffer (清潔)。
- **Teacher Talent (人才動能)**：支援服務地區 (Region)、雙層能力分類 (Skills) 與 動態標籤 (Hashtags)。
- **Approval 狀態機**：嚴格定義 PENDING -> APPROVED/REJECTED 流轉，核准時必執行 Re-validate。
- **內評與備註 (Internal Feedback)**：中心私有星等與評價，作為智慧媒合的權重因子。

### 4.2 視覺與體驗 (Premium UX)
- **設計風格**：採用 **Midnight Premium (深藍高質感)** 與 **Glassmorphism (毛玻璃)**。
- **雙模式適配**：支援深色 (20-35歲) 與 淺色 (35-45歲) 雙主題。
- **Grid 排課體驗**：Template + Cells 支援多教室 Grid，提供即時紅/綠/橘視覺回饋。
- 交易語意：
  - 任一筆 Hard Overlap → 全部 rollback
  - Buffer conflict → 可部分成功（created + pending）

### 4.3 商業化（Plan/Gating）
- Free（老師個人）→ Starter/Pro/Team（中心）
- Trial：新中心自動 Pro 試用（例：30 天）
- 到期降級 read-only（可看不可改）
- 限制項：teacher 數、template 數、rules 量等

---

## 5. 排課約束（Scheduling Constraints）
### 5.1 必須防呆
1) **OVERLAP**：同教室/同老師同時間不可重疊（硬衝突）  
2) **BUFFER**：
   - Room Buffer：A 結束到 B 開始需 ≥ X 分
   - Teacher Buffer：A 結束到 B 開始需 ≥ Y 分

### 5.2 Course Buffer Policy (課程緩衝政策)
- **課程級設定 (at Course Level)**:
  - `room_buffer_min`: 該課程專屬的教室清潔時間。
  - `teacher_buffer_min`: 該課程專屬的老師轉場時間。
- **班別級權限 (at Offering Level)**:
  - `allow_buffer_override`: 針對特定班別，決定是否允許在 Buffer 衝突時強制排入。

### 5.3 覆寫規則 (Override Logic)
- **Overlap (硬衝突)**：永遠不可覆寫。
- **Buffer (軟衝突)**：
  - 系統檢查該 Offering 是否開啟 `allow_buffer_override`。
  - 若開啟：管理員可進行「強制排入」，需填寫 `reason` 並記錄於 Audit Log。
  - 若關閉：即便是管理員，在出現 Buffer 衝突時也無法強制排入，確保特定高品質課程的專業間隔。

---

## 6. 核心引擎：Validate（v3 真相來源）
### 6.1 Validate 輸入
- center_id, teacher_id, start_at, end_at
- room_id（可選）
- ignore_rule_id（可選）
- override_buffer（bool，v3 必填）

### 6.2 Validate 輸出（給前端展示）
- conflicts[] 每筆包含：
  - type：OVERLAP / TEACHER_BUFFER / ROOM_BUFFER
  - required_minutes / diff_minutes
  - can_override / require_approval
  - conflict_source（RULE / SESSION / PERSONAL）

### 6.3 決策
- 有 overlap → 直接拒絕
- buffer 衝突：
  - 不允許 override → 拒絕
  - 允許 override 且 require approval → 回傳 approval_required（供產生 PENDING）
  - 允許直接 override → OK，但必須 reason + audit

---

## 7. Exception 規格（停課/改期）
### 7.1 類型
- CANCEL：停課（覆蓋該日 session）
- RESCHEDULE：改期（原日取消 + 新日新增）
- ADD_SESSION（選配）：補課（獨立新增）

### 7.2 狀態
- PENDING（待審）
- APPROVED（已核准/生效）
- REJECTED（已拒絕）
- CANCELLED（撤回，選配）

### 7.3 唯一性
- 同一 rule_id + date：最多只能有一筆 ACTIVE（PENDING/APPROVED）

### 7.4 審核約束（必備）
- approve 時 **必須 re-validate**
- re-validate 失敗 → 不可核准

---

## 8. 功能需求（FR）
### 8.1 Admin 後台（中心）
- **排課 Grid**：三欄式工作台，支持拖拉指派與即時 Validate。
- **人才庫管理**：主動搜尋 (Region/Hashtag) 已開啟媒合之老師，並發送邀請。
- **智慧代課**：針對異動時段自動依 Match Score 推薦老師。
- **管理詳情**：
  - 課程、班別、教室、老師、排課模板 (Templates/Cells)。
  - **管理員管理**：新增/維護管理員帳號，支援 Owner/Admin/Staff 分級。
- **報表與稽核**：
  - 操作紀錄 (Audit Log)：詳細記錄「哪個管理員在何時做了什麼」。

### 8.2 Teacher 端（RWD / LINE）
- **個人專業檔案**：編輯介面支援階層式技能與動態標籤輸入。
- **綜合課表**：整合多中心行程，支援週/3日視圖切換。
- **匯出課表**：生成 Mesh Gradient 背景的高級課表圖片。
- **教學紀錄**：單堂課程的教學與備課筆記讀寫。

---

## 9. API 契約（摘要版；正式以 API.md 為準）
### 9.1 共通
- 統一 envelope：{code,message,datas}
- auth：admin token / teacher token 分離
- 所有寫入都要寫 audit/event

### 9.2 重點 Endpoints
- `POST /public/validate`：官網 WOW Demo 專用匿名驗證。
- `GET /admin/talent/search`：跨租戶老師人才搜尋。
- `POST /admin/centers/{id}/rules/bulk`：網格批量儲存。
- `GET /common/hashtags/search`：標籤模糊搜尋建議。

> 錯誤碼：見 ERR_CODES_v3。

---

## 10. 資料庫（MySQL）
### 10.1 必備表（邏輯）
- centers / admin_users (含 name / role / status)
- teachers (含 Region / Hiring 開關)
- teacher_skills / hashtags / teacher_hashtags
- timetable_templates / timetable_cells
- schedule_rules / schedule_exceptions / session_notes
- audit_logs / center_invitations

### 10.2 索引與唯一性（重點）
- schedule_exceptions：rule_id + date 只限制 PENDING/APPROVED
- schedule_rules：center_id + weekday + start/end 索引
- timetable_cells：template_id + row_no + col_no 唯一

---

## 11. UI/UX（RWD）
### 11.1 Admin（Desktop）
- Grid 排課：左（Offering/Teacher）中（Grid）右（結果摘要）
- 指派即 validate
- 顏色：紅=不可排、橘=需審、綠=可排
- Bulk 儲存後：created/pending/errors 摘要 + 跳轉待審

### 11.2 Teacher（Mobile-first）
- 週視圖
- 點 session → bottom sheet（cancel/reschedule）
- validate 後提示：不可排/需審/可直接送出

---

## 12. 非功能需求（NFR）
### 12.1 安全
- token 過期（admin/teacher 分離）
- 強制 center scope（不可被繞過）
- audit/event 必須可追溯

### 12.2 可用性/性能
- 50 人規模：單機部署足夠
- 同時在線 <10：validate 需在可接受延遲內完成
- 每日 DB 備份（保留 7–14 天）

### 12.3 可維運
- 統一錯誤碼（ERR_CODES_v3）
- validate / approve 需 log + audit
- plan gating 需可觀測（事件）

---

## 13. Plan / 收費與限制（正式版）
（示例，可依市場調整）
- Free：老師個人（不含中心功能）
- Starter：小中心（teacher 上限、template 上限）
- Pro：中型中心（更高限制）
- Team：客製（無上限/多分店）

Trial：新中心自動 Pro 試用（例：30 天），到期降級 read-only。

---

## 14. 驗收標準（AC）
### 14.1 多中心隔離
- **Audit Log**：詳細記錄 Actor (管理員)、Action 與對象，確保責任可追溯性。
- **後台權限**：僅限擁有者 (Owner) 可增刪其他管理員帳號。

### 14.2 排課治理
- overlap 必須拒絕
- buffer 依該課程 (Course) 設定之時長進行驗證。
- 是否可覆寫 (Override) 依該班別 (Offering) 權限決定。

### 14.3 審核一致性
- approve 必須 re-validate；失敗不可核准

### 14.4 Bulk 行為
- HARD overlap → 全 rollback
- BUFFER → 可部分成功，pending 正確進入審核

### 14.5 商業 gating
- 試用到期後不可新增/修改排課（read-only）
- 超限提示與限制正確生效

### 14.6 彈性排課邏輯 (Drafting)
- **非特定老師排課**：系統必須允許 `teacher_id` 為空之排課紀錄。
- **視覺化警示**：無老師課程在 Grid 上必須以「斜線背景」＋「⚠️ 警示標籤」呈現。
- **關聯鏈路**：點擊無老師課程，必須能跳轉至「人才搜尋」或「智慧代課」面板完成指派。


## 15. 里程碑（工程化）
- S1：Data Foundation (Migrations/Seeds/Hashtag GC) - **基礎建設**
- S2：Auth & Resource CRUD (LINE Login/Profile/Templates) - **支撐層**
- S3：The Brain (Validate Engine/State Machine/Notifications) - **邏輯層**
- S4：Teacher App (Calendar Grid/Export Image/Notes) - **體驗層**
- S5：Admin Workspace (Scheduling Grid/Talent Search/Demo) - **交付層**
- S6：QA & Deploy (Production / Load Test) - **上線層**

---

## 16. 附錄：交付物建議（給 AI / OpenCode）
- 以 `MASTER_PROMPT_v3_REAL.md` 作為最高指令
- 每個 Stage 需提供：變更檔清單、如何測試、如何驗收 AC
- 所有寫入皆需 audit/event，且不能在 Controller 寫業務邏輯
