# TimeLedger 開發階段規劃 (Development Stages)

> [!IMPORTANT]
> **開發鐵律 (Atomic Vertical Slices)**
> 1. **功能隔離**：一次僅開發一個 Stage 中的單一子項。
> 2. **開發順序**：`Migration -> Unit Test (Mock) -> Backend Service -> API -> Frontend UI -> Integration Test`。
> 3. **原子提交**：每完成一個子項的「後端」或「前端」皆須 Commit 並更新 `progress_tracker.md`。

---

## Stage 1: 基建與設計系統 (Core & Design Tokens)
- [ ] **1.1 Workspace Init**: Docker Compose (MySQL 8, Redis), Monorepo 初始化。
- [ ] **1.2 Migrations (Base)**: 建立 `centers`, `users`, `geo_cities`, `geo_districts`。
- [ ] **1.3 UI Design System**: 
    - [ ] Tailwind Config (Midnight Indigo 漸層)、Google Fonts 引入。
    - [ ] 基礎組件：`BaseGlassCard`, `BaseButton`, `BaseInput`。
    - [ ] 基礎佈局：Admin Sidebar 與 Mobile Bottom Nav。

## Stage 2: 老師身份與專業檔案 (Teacher Identity & Skills)
- [ ] **2.1 Migrations (Skills & Certs)**: 建立 `teacher_skills`, `hashtags`, `teacher_ certificates` 等表。
- [ ] **2.2 Auth Implementation**: LINE Login (LIFF Silent), JWT 適配器。
- [ ] **2.3 Profile Logic**: 實作 **Hashtag 字典同步邏輯** (Update usage_count / GC)。
- [ ] **2.4 Profile UI**: 具有技能標籤掛載、證照上傳與求職開關的個人頁面。

## Stage 3: 中心管理與邀請流 (Center Admin & Invitation)
- [ ] **2.1 Migrations (Admin & Membership)**: 建立 `admin_users`, `center_memberships`, `center_invitations`。
- [ ] **3.2 Admin Auth**: Email/Password 登入與 Role-based 權限牆。
- [ ] **3.3 Staff Management**: 管理員帳號 CRUD 與停權功能。
- [ ] **3.4 Invitation Flow**: 產生邀請碼 Token 至老師點擊連結加入中心的整合流程。

## Stage 4: 中心資源與緩衝設定 (Resources & Buffers)
- [ ] **4.1 Migrations (Resources)**: 建立 `rooms`, `courses`, `offerings`。
- [ ] **4.2 Rooms**: 教室 CRUD (含 **Capacity** 校驗)。
- [ ] **4.3 Courses**: 課程模板 (含 **Buffer Min** 分鐘數設定)。
- [ ] **4.4 Offerings**: 開課班別定義 (含 **Buffer Override** 覆寫旗標)。

## Stage 5: 排課引擎 I - 週期展開 (Scheduling Engine: Rules)
- [ ] **5.1 Migrations (Rules)**: 建立 `schedule_rules`。
- [ ] **5.2 Rules API**: 循環規則基本 CRUD (TDD)。
- [ ] **5.3 Expander Logic**: `expandRules()` 核心邏輯，支援每週/隔週規律。
- [ ] **5.4 Unified Calendar**: 老師端 **多中心綜合課表**。

## Stage 6: 排課引擎 II - 衝突驗證 (Scheduling Engine: Validation)
- [ ] **6.1 Validation Engine**: `checkOverlap`, `checkBuffer` 核心檢測 (TDD)。
- [ ] **6.2 Conflict UI**: 拖拉式排課工作台，整合紅/綠/橘視覺回饋。
- [ ] **6.3 Bulk Validate**: 儲存時的批量原子校驗 logic。

## Stage 7: 排課引擎 III - 週期過渡 (Dynamic Phases)
- [ ] **7.1 Phase Support**: 實作 `effective_start/end` 邏輯控管。
- [ ] **7.2 Transition Flow**: 處理同一班別不同月份更換教室或時段的過渡介面。

## Stage 8: 國定假日與自動化邏輯 (Holiday Automation)
- [ ] **8.1 Migrations (Holidays)**: 建立 `center_holidays`。
- [ ] **8.2 Holiday CRUD**: 各中心自定義假日管理 UI。
- [ ] **8.3 Auto-Filter**: 展開課表時「無感」隱藏假日行程，且不產生髒資料。

## Stage 9: 異動審核與狀態機 (Exceptions & Approvals)
- [ ] **9.1 Migrations (Exceptions)**: 建立 `schedule_exceptions`。
- [ ] **9.2 Exception API**: 老師端發起停課 / 改期請求。
- [ ] **9.3 Approval Workflow**: 管理員審核、反對與 **Re-validation (前置防撞)**。

## Stage 10: 預約排課與截止鎖定 (Deadlines & Locking)
- [ ] **10.1 Locking Logic**: `lock_at` 與 `exception_lead_days` 動態鎖定檢測。
- [ ] **10.2 Lock UI**: 老師端功能按鈕的自動禁用與提示。

## Stage 11: 人才市場與智慧媒合 (Talent Search & Match)
- [ ] **11.1 Migrations (Notes)**: 建立 `center_teacher_notes`。
- [ ] **11.2 Talent Discovery**: 全球老師搜尋介面 (Skill/Region 篩選)。
- [ ] **11.3 Smart Matcher**: 評分權限因子排序推薦列表，代課建議 Drawer。
- [ ] **11.4 Internal Notes**: 中心內的評分與私密備註系統。

## Stage 12: 營運、通知與登陸頁 (Operational Polish)
- [ ] **12.1 Migrations (Notes & Logs)**: 建立 `session_notes`, `audit_logs`。
- [ ] **12.2 Operation Logic**: 教學筆記、管理員稽核歷史紀錄。
- [ ] **12.3 Notifications**: LINE 每日提醒、緊急異動與邀請推播。
- [ ] **12.4 Export & Branding**: 匯出具備 **個人品牌標籤** 的精美週課表圖。
- [ ] **12.5 Public Sandbox**: 首頁免登入的高質感排課模擬器。
- [ ] **Final QA**: 通過 `Integration_Playbook.md` 全功能驗收。
