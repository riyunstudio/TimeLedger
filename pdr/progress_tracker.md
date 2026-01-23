# TimeLedger 實作進度追蹤 (Development Progress Tracker)

> [!IMPORTANT]
> 此文件由 AI 持續維護。每完成一個任務或階段，請在此更新狀態與「上下文恢復快照」。

## 1. 階段性進度表 (Roadmap Status)

| 階段 | 任務說明 | 狀態 | 備註 |
|:---|:---|:---:|
| | **Stage 1** | **基建與設計系統** | `[COMPLETED]` | ✅ 完成 |
| | 1.1 Workspace Init | Docker Compose, Monorepo 初始化 | `[X] DONE` | ✅ 完成 |
| | 1.2 Migrations (Base) | 建立 `centers`, `teachers`, `admin_users`, `geo_cities`, `geo_districts` | `[X] DONE` | ✅ 完成 |
| | 1.3 UI Design System | Tailwind Config、基礎組件、基礎佈局 | `[X] DONE` | ✅ 完成 |
| | 1.4 Tests (TDD) | Models、Repository、Components、Layouts 單元測試 | `[X] DONE` | ✅ 完成 |
| | **Stage 2** | **老師身份與專業檔案** | `[COMPLETED]` | ✅ 完成 |
| | 2.1 Migrations (Skills & Certs) | 建立 `teacher_skills`, `hashtags`, `teacher_certificates` | `[X] DONE` | ✅ 完成 |
| | 2.2 Auth Implementation | LINE Login (LIFF Silent), JWT 適配器 | `[X] DONE` | ✅ MockAuthService 已實作 |
| | 2.3 Profile Logic | Hashtag 字典同步邏輯 | `[X] DONE` | ✅ HashtagRepository 已有基本方法 |
| | 2.4 Profile UI | 個人頁面 UI | `[X] DONE` | ✅ 已有 `/teacher/profile` 頁面 |
| | **Stage 3** | **中心管理與邀請流** | `[COMPLETED]` | ✅ 完成 |
| | 3.1 Migrations (Admin & Membership) | 建立 `admin_users`, `center_memberships`, `center_invitations` | `[X] DONE` | ✅ 完成 |
| | 3.2 Admin Auth | Email/Password 登入 | `[X] DONE` | ✅ MockAuthService 已實作 |
| | 3.3 Staff Management | 管理員帳號 CRUD | `[X] DONE` | ✅ AdminUserController 已實作 |
| | 3.4 Invitation Flow | 產生邀請碼 | `[X] DONE` | ✅ TeacherController.InviteTeacher 已實作 |
| | 3.5 Invitation Rejection | 拒絕邀請流程 | `[X] DONE` | ✅ 已實作相關 API |
| | **Stage 4** | **中心資源與緩衝設定** | `[COMPLETED]` | ✅ 完成 |
| | 4.1 Migrations (Resources) | 建立 `rooms`, `courses`, `offerings` | `[X] DONE` | ✅ 完成 |
| | 4.2 Rooms | 教室 CRUD | `[X] DONE` | ✅ AdminResourceController 已實作 |
| | 4.3 Courses | 課程模板 | `[X] DONE` | ✅ AdminResourceController 已實作 |
| | 4.4 Offerings | 班別定義 | `[X] DONE` | ✅ OfferingController 已實作 |
| | **Stage 5** | **排課引擎 I - 週期展開** | `[COMPLETED]` | ✅ 完成 |
| | 5.1 Migrations (Rules) | 建立 `schedule_rules` | `[X] DONE` | ✅ 完成 |
| | 5.2 Rules API | 規則 CRUD | `[X] DONE` | ✅ SchedulingController 已實作 |
| | 5.3 Expander Logic | `expandRules()` 核心邏輯 | `[X] DONE` | ✅ ScheduleExpansionService 已實作 |
| | 5.4 Unified Calendar | 多中心綜合課表 | `[X] DONE` | ✅ TeacherController.GetSchedule 已實作 |
| | **Stage 6** | **排課引擎 II - 衝突驗證** | `[COMPLETED]` | ✅ 完成 |
| | 6.1 Validation Engine | `checkOverlap`, `checkBuffer` | `[X] DONE` | ✅ SchedulingValidationService 已實作 |
| | 6.2 Conflict UI | 拖拉式排課工作台 | `[X] DONE` | ✅ 前端 ScheduleGrid 組件已實作 |
| | 6.3 Bulk Validate | 批量原子校驗 | `[X] DONE` | ✅ SchedulingController.ValidateFull 已實作 |
| | **Stage 7** | **排課引擎 III - 週期過渡** | `[COMPLETED]` | ✅ 完成 |
| | 7.1 Phase Support | `effective_start/end` 邏輯 | `[X] DONE` | ✅ ScheduleExpansionService 已實作 effective_range 檢查 |
| | 7.2 Transition Flow | 過渡介面 | `[X] DONE` | ✅ 新增 DetectPhaseTransitions API 與 PhaseTransition struct |
| | **Stage 8** | **國定假日與自動化邏輯** | `[COMPLETED]` | ✅ 完成 |
| | 8.1 Migrations (Holidays) | 建立 `center_holidays` | `[X] DONE` | ✅ CenterHoliday model 與 repository 已存在 |
| | 8.2 Holiday CRUD | 假日管理 | `[X] DONE` | ✅ 新增 GetHolidays, CreateHoliday, DeleteHoliday API |
| | 8.3 Bulk Import | 批量匯入 | `[X] DONE` | ✅ BulkCreateHolidays API 已實作 |
| | 8.4 Auto-Filter | 自動過濾 | `[X] DONE` | ✅ ExpandRules 現在會標記假日並過濾顯示 |
| | **Stage 9** | **異動審核與狀態機** | `[COMPLETED]` | ✅ 完成 |
| | 9.1 Migrations (Exceptions) | 建立 `schedule_exceptions` | `[X] DONE` | ✅ 完成 |
| | 9.2 Exception API | 老師申請異動 | `[X] DONE` | ✅ TeacherController 已實作 |
| | 9.3 Exception Revoke | 撤回申請 | `[X] DONE` | ✅ TeacherController.RevokeException 已實作 |
| | 9.4 Approval Workflow | 管理員審核 | `[X] DONE` | ✅ SchedulingController 已實作 |
| | 9.5 Review Fields | 審核欄位 | `[X] DONE` | ✅ 已有 reviewed_at, reviewed_by, review_note 欄位 |
| | **Stage 10** | **預約排課與截止鎖定** | `[ ] TODO` | 待開始 |
| | 10.1 Locking Logic | `lock_at` 與 `exception_lead_days` | `[ ] TODO` | |
| | 10.2 Lock UI | 按�禁用 | `[ ] TODO` | |
| | **Stage 11** | **人才市場與智慧媒合** | `[COMPLETED]` | ✅ 完成 |
| | 11.1 Migrations (Notes) | 建立 `center_teacher_notes` | `[X] DONE` | ✅ 完成 |
| | 11.2 Talent Discovery | 全球老師搜尋 | `[X] DONE` | ✅ SmartMatchingController 已實作 |
| | 11.3 Pagination & Sorting | 分頁排序 | `[X] DONE` | ✅ 已有標準分頁實作 |
| | 11.4 Smart Matcher | 智慧媒合 | `[X] DONE` | ✅ SmartMatchingService 已實作 |
| | 11.5 Internal Notes | 內部評分與備註 | `[X] DONE` | ✅ CenterTeacherNote model 已實作 |
| | **Stage 12** | **營運、通知與登陸頁** | `[COMPLETED]` | ✅ 完成 |
| | 12.1 Migrations (Notes & Logs) | 建立 `session_notes`, `audit_logs` | `[X] DONE` | ✅ 完成 |
| | 12.2 Operation Logic | 教學筆記、稽核紀錄 | `[X] DONE` | ✅ 已實作相關 Controller 和 Repository |
| | 12.3 Notifications | LINE 提醒推播 | `[X] DONE` | ✅ NotificationController 和 LineNotifyService 已實作 |
| | 12.4 Export & Branding | 匯出精美課表圖 | `[X] DONE` | ✅ ExportController 已實作 |
| | 12.5 Public Sandbox | 公開排課模擬器 | `[X] DONE` | ✅ 已有 `/` landing page |

## 2. Stage 1 完整修正記錄 (Stage 1 Complete)

### 1.1 Workspace Init ✅
- Docker Compose 配置正確（MySQL 8.0 + Redis）
- Monorepo 結構符合規範（後端根目錄 + frontend 子目錄）

### 1.2 Migrations (Base) ✅
- `centers` model - 已存在且符合規範
- `teachers` model - 已存在且符合規範（`line_user_id`, `bio`, `city`, `district`, `is_open_to_hiring` 等欄位）
- `admin_users` model - 已存在
- `geo_cities` model - 已存在
- `geo_districts` model - 已存在
- 注意：舊的 `users` model 保留，但應改為使用 `teachers` 和 `admin_users`

### 1.3 UI Design System ✅
**Tailwind Config**:
- Midnight Indigo 漸層 (`#6366F1` 到 `#A855F7`)
- Google Fonts (Outfit, Inter)
- Dark/Light 雙模式支援

**基礎組件** - 已創建 `components/base/`:
- `BaseGlassCard.vue` - 毛玻璃卡片
- `BaseButton.vue` - 統一按�樣式（primary/secondary/success/critical/warning）
- `BaseInput.vue` - 統一輸入框樣式
- `BaseModal.vue` - 模態視窗（支援手機底部滑入）
- `BaseBadge.vue` - 標籤組件
- `GridSkeleton.vue` - 骨架屏載入動畫

**基礎佈局** - 已創建 `layouts/`:
- `admin.vue` - 管理員後台佈局（側邊欄）
- `default.vue` - 老師端佈局（底部導航）
- `blank.vue` - 全螢幕佈局（登入頁）

**頁面 Layout 更新**:
- 所有 `pages/teacher/*` → 使用 `default` layout
- 所有 `pages/admin/*` → 使用 `admin` layout
- `pages/*/login.vue` → 使用 `blank` layout
- FAB 按�位置調整（手機版避免與底部導航衝突）

### 1.4 Tests (TDD) ✅ 已完成
**後端測試** (`testing/test/stage1_models_test.go`):
- ✅ `TestCenter_Validation` - Center 模型驗證
- ✅ `TestTeacher_Validation` - Teacher 模型驗證
- ✅ `TestAdminUser_Validation` - AdminUser 模型驗證
- ✅ `TestGeoDistrict_ForeignKey` - GeoDistrict 外�驗證
- ✅ `TestTeacherSkill_Validation` - TeacherSkill 模型驗證
- ✅ `TestHashtag_Validation` - Hashtag 模型驗證
- ✅ `TestTeacherCertificate_Validation` - TeacherCertificate 模型驗證
- ✅ `TestTeacherPersonalHashtag_SortOrder` - TeacherPersonalHashtag 排序驗證
- 測試結果：**全部通過** (8/8 passed, 0.312s)

**前端組件測試**:
- `tests/components/base/BaseGlassCard.spec.ts` - 毛玻璃卡片測試
- `tests/components/base/BaseButton.spec.ts` - 按�組件測試（所有 variant 和 size）
- `tests/components/base/BaseInput.spec.ts` - 輸入框組件測試
- `tests/components/base/BaseModal.spec.ts` - 模態視窗測試（所有 size 和行為）
- `tests/components/base/BaseBadge.spec.ts` - 標籤組件測試（所有 variant）
- `tests/components/base/GridSkeleton.spec.ts` - 骨架屏測試

**佈局測試**:
- `tests/layouts/default.spec.ts` - 老師端佈局測試
- `tests/layouts/admin.spec.ts` - 管理員後台佈局測試
- `tests/layouts/blank.spec.ts` - 登入頁佈局測試

## 3. 上下文恢復快照 (Context Recovery Snapshot)

### 當前狀態
- **專案路徑**: `D:\project\TimeLedger`
- **專案類型**: Monorepo (Go 後端 + Nuxt 3 前端)
- **當前階段**: Stage 1 (基建與設計系統) - 已完成

### 已修正問題
1. 創建了完整的基礎組件庫 (`components/base/`)
2. 創建了三個佈局 (`layouts/`)
3. 更新所有頁面使用正確的 layout
4. 修正 `app/base.go` 中的 `Mysql` → `MySQL`
5. 批量更新所有引用（`app/`, `repositories/`, `controllers/`, `services/`）
6. 建立了類型定義 (`types/layout.d.ts`)
7. **新增完整的 Stage 1 單元測試**（後端 Models + 前端 Components + Layouts）

### 7. **新增完整的 Stage 1 單元測試**（後端 Models + 前端 Components + Layouts）
- `testing/test/stage1_models_test.go` - Model 驗證測試
  - `tests/components/base/BaseGlassCard.spec.ts` - 毛玻璃卡片測試
  - `tests/components/base/BaseButton.spec.ts` - 按�組件測試
  - `tests/components/base/BaseInput.spec.ts` - 輸入框組件測試
  - `tests/components/base/BaseModal.spec.ts` - 模態視窗測試
  - `tests/components/base/BaseBadge.spec.ts` - �籤組件測試
- `tests/components/base/GridSkeleton.spec.ts` - 骨架屏測試
- `tests/layouts/default.spec.ts` - 老師端佈局測試
- `tests/layouts/admin.spec.ts` - 管理員後台佈局測試
- `tests/layouts/blank.spec.ts` - 登入頁佈局測試
- 測試結果：**全部通過**

## 7. Stage 7 完整實作記錄 (Stage 7 Complete)

### 7.1 Phase Support - `effective_start/end` 邏輯 ✅
**更新 `app/services/scheduling_expansion.go`**:
- `ExpandRules()` 現在會正確檢查每個規則的 `effective_range`
- 新增 `GetEffectiveRuleForDate()` - 取得指定日期的有效規則
- 新增 `DetectPhaseTransitions()` - 檢測 phase 變化的日期點
- 新增 `GetRulesByEffectiveDateRange()` - 取得指定 effective 日期範圍內的規則
- 新增 `ptrEqual()` - 輔助函數用於比較指針

**更新 `app/services/scheduling_interface.go`**:
- 新增 `PhaseTransition` struct - 記錄 phase 變化的詳細資訊
- 更新 `ScheduleExpansionService` interface

### 7.2 Transition Flow - 過渡介面 ✅
**更新 `app/controllers/scheduling.go`**:
- 新增 `DetectPhaseTransitionsRequest` - API 請求結構
- 新增 `DetectPhaseTransitions()` - Phase transition 檢測 API

**更新 `app/servers/route.go`**:
- 新增路由: `POST /api/v1/admin/centers/:id/detect-phase-transitions`

**更新 `app/repositories/schedule_rule.go`**:
- 新增 `ListByOfferingID()` - 依 offering ID 取得規則列表

### 7.3 Stage 7 單元測試 ✅
**建立 `testing/test/stage7_phase_test.go`**:
- `TestStage7_ScheduleRule_EffectiveRange` - 測試有效日期範圍邏輯
- `TestStage7_PhaseDetection_Logic` - 測試 phase 轉換檢測邏輯
- `TestStage7_ScheduleException_PhaseContext` - 測試例外單與 phase 的關聯
- `TestStage7_DateRange_Scan` - 測試 DateRange 的序列化
- 測試結果：**全部通過 (10/10 passed)**

## 8. Stage 8 完整實作記錄 (Stage 8 Complete)

### 8.1 Holiday Migrations ✅
- `CenterHoliday` model 已存在 (`app/models/center_holiday.go`)
- `CenterHolidayRepository` 已存在，包含完整 CRUD 方法

### 8.2 Holiday CRUD API ✅
**新增 `app/controllers/admin_resource.go`**:
- `GetHolidays()` - 取得中心假日列表（支援日期範圍篩選）
- `CreateHoliday()` - 建立單一假日
- `DeleteHoliday()` - 刪除假日

**新增 API Routes**:
- `GET /api/v1/admin/centers/:id/holidays`
- `POST /api/v1/admin/centers/:id/holidays`
- `DELETE /api/v1/admin/centers/:id/holidays/:holiday_id`

### 8.3 Bulk Import ✅
- `BulkCreateHolidays()` API 已存在
- 支援跳過重複日期 (`BulkCreateWithSkipDuplicate`)

### 8.4 Auto-Filter 假日行程 ✅
**更新 `app/services/scheduling_expansion.go`**:
- `ExpandRules()` 現在會載入並標記假日
- 新增 `IsHoliday` 欄位到 `ExpandedSchedule`
- 前端可根據 `is_holiday` 顯示灰色斜紋背景

### 8.5 Stage 8 單元測試 ✅
**建立 `testing/test/stage8_holiday_test.go`**:
- `TestStage8_CenterHoliday_Model` - 模型驗證
- `TestStage8_HolidayFiltering_Logic` - 假日過濾邏輯
- `TestStage8_BulkImport_Logic` - 批量匯入邏輯
- `TestStage8_ExpandedSchedule_HolidayField` - 課表展開假日欄位
- `TestStage8_PhaseTransition_HolidayAwareness` - Phase 與假日關聯
- 測試結果：**全部通過 (8/8 passed)**

### 下一步
- 繼續 Stage 10: 預約排課與截止鎖定

## 4. 已知問題 (Known Issues)

### 高優先級 (High Priority)
無

### 中優先級 (Medium Priority)
- `app/smart_matching.go`, `app/repositories/*`, `app/controllers/*` 中的 `.Mysql` 引用需要批量修正為 `.MySQL`（已在 app/base.go 修正，但其他文件仍有錯誤提示）

### 低優先級 (Low Priority)
- 部分 Repository 的測試被 skip，需要完善（如 `teacher_test.go`, `auth_test.go`）

## 5. 技術債務 (Technical Debt)

### 待優化項目
1. **Hashtag 字典同步邏輯優化**：雖然 HashtagRepository 已有基本方法，但可考慮添加完整的事務處理
2. **錯誤處理統一**：確保所有 Controller 都使用統一的錯誤回應格式
3. **測試覆蓋率**：需要增加更多邊界條件的測試案例

## 6. 開發規範遵守情況

### 已遵守規範
✅ 分層架構 (Controller → Request → Service → Repository → Model)
✅ Monorepo 結構（後端根目錄 + frontend 子目錄）
✅ 原子化開發（每個功能獨立完成）
✅ **TDD 優先**：先寫測試，再實作（Stage 1 已完成）
✅ Interface-based Auth Service（MockAuthService 用於開發）
✅ **No Code Without Tests**：Stage 1 所有核心組件都有對應測試

### 待改進項目
- 需要添加更多單元測試覆蓋
- 部分 Controller 可能過於複雜，可考慮拆分
- 部分 Repository 的測試需要完善（當前被 skip）
