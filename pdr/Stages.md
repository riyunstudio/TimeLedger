# TimeLedger 開發階段規劃 (Development Stages)

本專案採 **「先後端核心，後前端介面」** 的開發順序。以下為詳細的執行項目 (Todo List)。

---

## Stage 0: 專案基建 (已完成)
- [x] **Monorepo Init**: `backend/` (Go) + `frontend/` (Nuxt).
- [x] **Infra Setup**: `docker-compose.yml` (MySQL + Redis).
- [x] **Docs**: DEV_GUIDELINES, API Docs, Tech Specs.

---

## Stage 1: 資料庫與核心模型 (Data Foundation)
**目標**: 建立穩固的資料層，並產生足以支持開發的測試數據。
- [ ] **Migrations (Core)**:
    - [ ] `centers` (含 Plan/Policy JSON)
    - [ ] `users` (Teachers & Admins)
    - [ ] `memberships` (關聯表)
- [ ] **Migrations (Scheduling)**:
    - [ ] `courses` & `offerings` & `rooms`
    - [ ] `schedule_rules` (RB 週期規則)
    - [ ] `schedule_exceptions` (例外單)
    - [ ] `personal_events` (含 Recurrence JSON)
- [ ] **Migrations (Features)**:
    - [ ] `teacher_skills` (Hierarchical) & `teacher_certificates`
    - [ ] `hashtags` (Master) & `teacher_hashtags` (Mapping)
    - [ ] `timetable_templates` & `timetable_cells` (Grid Layouts)
    - [ ] `session_notes` (Teaching Records)
    - [ ] `center_invitations` (Traceable)
- [ ] **Data Seeding (Fake Data)**:
    - [ ] Seed 3 Centers (Starter, Growth, Pro 方案)
    - [ ] Seed 20 Teachers (混合: 有證照/無證照, 開啟媒合/關閉)
    - [ ] Seed 複雜排課場景 (故意製造衝突以供測試)
- [ ] **UI/UX Foundation (Frontend Base)**:
    - [ ] **Theme Config**: Setup Tailwind colors & Glassmorphism tokens in `frontend/`.
    - [ ] **Component Atomic Structure**: Initialize basic Layouts (Mobile Nav, Admin Sidebar).
- [ ] **Testing (Stage 1)**:
    - [ ] **Repository Tests**: 撰寫對應 DB CRUD 的單元測試，確保 SQL 邏輯正確。

---

## Stage 2: 身份與資源 API (Auth & Resources)
**目標**: 完成登入、權限、以及基礎的 CRUD。
- [ ] **Auth System**:
    - [ ] Admin Login (Email/Pass) + JWT
    - [ ] Teacher Login (LINE OIDC) + Auto Register
    - [ ] RBAC Middleware (`requireAdmin`, `requireTeacher`)
- [ ] **Teacher Profile API**:
    - [ ] GET/PUT `/teacher/me/profile` (含 Hiring Toggle, Public Contact)
    - [ ] Upload `/teacher/me/certificates` (S3/Local)
- [ ] **Admin Resource API**:
    - [ ] CRUD Rooms / Courses / Offerings
    - [ ] CRUD **Timetable Templates & Cells**
    - [ ] **Admin Account CRUD** (List, Create, Delete with Role)
    - [ ] CRUD Teachers (Invite Flow with Contact Info)
    - [ ] GET/PATCH Center Settings
- [ ] **Testing (Stage 2)**:
    - [ ] **API Logic Tests**: 使用 Go 客戶端模擬 API 呼叫，驗證 JWT 與 RBAC 攔截是否正確。

---

## Stage 3: 排課引擎與通知 (The Brain)
**目標**: 本專案最困難的後端邏輯。
- [ ] **Validation Engine**:
    - [ ] `checkOverlap()`: 嚴格時段重疊檢查
    - [ ] `checkBuffer()`: 教室清潔/老師轉場緩衝檢查
- [ ] **Scheduling Logic**:
    - [ ] `expandRules()`: 將週規則展開為實際日期
    - [ ] `createException()`: 停課/改期 狀態機 (Pending -> Approved)
    - [ ] **DB Lock**: 實作 `SELECT ... FOR UPDATE` 防止併發寫入
- [ ] **Smart Features**:
    - [ ] Matching Algo: 實作老師推薦評分 (`findMatches`)
    - [ ] Talent Search: 實作 `searchTalent(skills, keyword)`
- [ ] **Notification System**:
    - [ ] 實作 LINE Notify Client
    - [ ] Cron Job: 每日 20:00 發送「明日課程提醒」
    - [ ] Event Hooks: 邀請通知、緊急異動通知
- [ ] **Testing (Stage 3 - Critical)**:
    - [ ] **Logic Unit Tests**: 針對 `validator` 與 `expander` 撰寫極端案例測試 (邊界值、跨日)。
    - [ ] **Concurrency Tests**: 模擬多執行緒競爭排課，驗證 DB Lock 是否有效。

---

## Stage 4: 老師端 App (Teacher Frontend)
**目標**: 手機版極致體驗。
- [ ] **Dashboard**:
    - [ ] 統一視圖 (Unified View) 呈現多中心課程
    - [ ] 週次切換 (Prev/Next) 與 3日視圖適配
- [ ] **Actions**:
    - [ ] 課程詳情 Modal (顯示教室/人數)
    - [ ] **Session Notes**: 讀寫 上課紀錄/備課筆記
    - [ ] **Personal Event**: 新增/編輯 私人行程
- [ ] **Profile**:
    - [ ] 編輯個人檔案 & 求職設定開關
    - [ ] 上傳證照 UI
- [ ] **Features**:
    - [ ] 匯出精美課表 (html2canvas)
- [ ] **Testing (Stage 4)**:
    - [ ] **Frontend Unit Tests**: 驗證週次切換邏輯、時區轉換是否正確。

---

## Stage 5: 中心後台 (Center Frontend)
**目標**: 高效率排課治理。
- [ ] **Scheduling Grid**:
    - [ ] 大型週曆元件 (支援 Drag & Drop)
    - [ ] 即時衝突檢測顯示 (紅/綠/橘)
- [ ] **Management**:
    - [ ] 審核中心 (Approval Queue) 介面
    - [ ] 資源管理頁 (老師/教室/課程)
- [ ] **Talent Integration**:
    - [ ] 智慧代課 Drawer (推薦名單)
    - [ ] **Talent Search Page**: 搜尋開放老師、查看履歷、發送邀請
- [ ] **Public Landing Page**:
    - [ ] **WOW Hero Section** (Branding)
    - [ ] **Interactive Demo Sandbox** (Validate API integration)
- [ ] **Testing (Stage 5)**:
    - [ ] **Integration Tests**: 驗證拖曳後是否成功觸發後端 API 與錯誤處理。

---

## Stage 6: 驗收與部署 (QA & Deploy)
- [ ] **Integration Test**: 測試「邀請 -> 排課 -> 請假 -> 找代課 -> 審核」全流程
- [ ] **Load Test**: 模擬 50 人同時操作排課
- [ ] **Production Deploy**: Docker Compose 上線 VPS
- [ ] **Documentation**: API 文件與操作手冊
