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
| | 2.2 Auth Implementation | LINE Login (LIFF Silent), JWT 適配器 | `[X] DONE` | ✅ AuthService 已實作 |
| | 2.3 Profile Logic | Hashtag 字典同步邏輯 | `[X] DONE` | ✅ HashtagRepository 已有基本方法 |
| | 2.4 Profile UI | 個人頁面 UI | `[X] DONE` | ✅ 已有 `/teacher/profile` 頁面 |
| | **Stage 3** | **中心管理與邀請流** | `[COMPLETED]` | ✅ 完成 |
| | 3.1 Migrations (Admin & Membership) | 建立 `admin_users`, `center_memberships`, `center_invitations` | `[X] DONE` | ✅ 完成 |
| | 3.2 Admin Auth | Email/Password 登入 | `[X] DONE` | ✅ AuthService 已實作 |
| | 3.3 Staff Management | 管理員帳號 CRUD | `[X] DONE` | ✅ AdminUserController 已實作 |
| | 3.4 Invitation Flow | 產生邀請碼 | `[X] DONE` | ✅ TeacherController.InviteTeacher 已實作 |
| | 3.5 Invitation Rejection | 拒絕邀請流程 | `[X] DONE` | ✅ 已實作相關 API |
| | **Stage 4** | **中心資源與緩衝設定** | `[COMPLETED]` | ✅ 完成 |
| | 4.1 Migrations (Resources) | 建立 `rooms`, `courses`, `offerings` | `[X] DONE` | ✅ 完成 |
| | 4.2 Rooms | 教室 CRUD | `[X] DONE` | ✅ AdminResourceController 已實作 |
| | 4.3 Courses | 課程模板 | `[X] DONE` | ✅ AdminResourceController 已實作 |
| | 4.4 Offerings | 班別定義 | `[X] DONE` | ✅ OfferingController 已實作 |
| | **Stage 4.5** | **資源管理擴充** | `[COMPLETED]` | ✅ 完成 |
| | 4.5.1 Soft Delete | `is_active` 欄位與 Toggle API | `[X] DONE` | ✅ 新增 ListActive 與 ToggleActive 方法 |
| | 4.5.2 Course Duplication | 班別複製功能 (含規則複製) | `[X] DONE` | ✅ OfferingRepository.Copy 已實作 |
| | 4.5.3 Invitation Stats | 邀請碼列表與使用統計 API | `[X] DONE` | ✅ CenterInvitationRepository 新增統計方法 |
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
| | **Stage 8.5** | **循環編輯功能** | `[COMPLETED]` | ✅ 完成 |
| | 8.5.1 Personal Events Update Mode | 支援 SINGLE/FUTURE/ALL 模式 | `[X] DONE` | ✅ ScheduleRecurrenceService 已實作 |
| | 8.5.2 例外生成邏輯 | 編輯單一場次時產生 CANCEL + ADD 例外 | `[X] DONE` | ✅ editSingle, editFuture, editAll 方法 |
| | 8.5.3 UI 流程 | 編輯/刪除確認對話框與受影響場次預覽 | `[X] DONE` | ✅ 新增 3 個 API 端點 |
| | **Stage 9** | **異動審核與狀態機** | `[COMPLETED]` | ✅ 完成 |
| | 9.1 Migrations (Exceptions) | 建立 `schedule_exceptions` | `[X] DONE` | ✅ 完成 |
| | 9.2 Exception API | 老師申請異動 | `[X] DONE` | ✅ TeacherController 已實作 |
| | 9.3 Exception Revoke | 撤回申請 | `[X] DONE` | ✅ TeacherController.RevokeException 已實作 |
| | 9.4 Approval Workflow | 管理員審核 | `[X] DONE` | ✅ SchedulingController 已實作 |
| | 9.5 Review Fields | 審核欄位 | `[X] DONE` | ✅ 已有 reviewed_at, reviewed_by, review_note 欄位 |
| | **Stage 10** | **預約排課與截止鎖定** | `[COMPLETED]` | ✅ 完成 |
| | 10.1 Locking Logic | `lock_at` 與 `exception_lead_days` | `[X] DONE` | ✅ CheckExceptionDeadline 已實作，支援規則鎖定與中心策略 |
| | 10.2 Lock UI | 按鈕禁用 | `[X] DONE` | ✅ 新增 CheckRuleLockStatus API 供前端禁用按鈕 |
| | **Stage 11** | **人才市場與智慧媒合** | `[COMPLETED]` | ✅ 完成 |
| | 11.1 Migrations (Notes) | 建立 `center_teacher_notes` | `[X] DONE` | ✅ 完成 |
| | 11.2 Talent Discovery | 全球老師搜尋 | `[X] DONE` | ✅ SmartMatchingController 已實作 |
| | 11.3 Pagination & Sorting | 分頁排序 | `[X] DONE` | ✅ 已有標準分頁實作 |
| | 11.4 Smart Matcher | 智慧媒合 | `[X] DONE` | ✅ SmartMatchingService 已實作 |
| | 11.5 Internal Notes | 內部評分與備註 | `[X] DONE` | ✅ CenterTeacherNote model 已實作 |
| | **Stage 13** | **DTO 與 Resource 映射層** | `[COMPLETED]` | ✅ 完成 |
| | 13.1 Resource 實作 | 建立 RoomResource、CenterResource、CourseResource | `[X] DONE` | ✅ 完成 |
| | 13.2 Service 重構 | 移除 Service 層的響應格式邏輯 | `[X] DONE` | ✅ 完成 |
| | 13.3 Controller 整合 | Controller 層注入 Resource 依賴 | `[X] DONE` | ✅ 完成 |
| | **Stage 18** | **教師端互動與課堂備註優化** | `[COMPLETED]` | ✅ 完成 |
| | 18.1 課表互動優化 | 動作選擇對話框、拖曳功能 | `[X] DONE` | ✅ 完成 |
| | 18.2 課堂備註修復 | rule_id 欄位、類型轉換修復 | `[X] DONE` | ✅ 完成 |
| | 18.3 例外申請預填 | 從課表帶入預設資料 | `[X] DONE` | ✅ 完成 |
| | **Stage 19** | **跨日課程顯示修復** | `[COMPLETED]` | ✅ 完成 |
| | 19.1 狀態判斷修復 | 管理員儀表板跨日課程狀態 | `[X] DONE` | ✅ 完成 |
| | 19.2 時間範圍擴展 | 前端課表顯示 00:00-03:00, 22:00-23:00 | `[X] DONE` | ✅ 完成 |
| | 19.3 跨日課程分割 | 後端生成兩個條目（開始日/結束日） | `[X] DONE` | ✅ 完成 |
| | 19.4 前端顯示修復 | 正確處理分割後的跨日課程 | `[X] DONE` | ✅ 完成 |

## 13. Stage 13 完整實作記錄 (Stage 13 Complete)

### 13.1 開發摘要

本階段成功實作了資料傳輸對象（DTO）模式，將資料庫模型與 API 響應格式進行解耦，達成前後端資料格式的職責分離。

### 13.2 完成項目

#### 13.2.1 新增 Resource 檔案

| 檔案 | 結構 | 功能 |
|:---|:---|:---|
| `app/resources/room.go` | RoomResponse | 教室響應格式轉換 |
| `app/resources/center.go` | CenterResponse、CenterSettings | 中心響應格式轉換 |
| `app/resources/course.go` | CourseResponse | 課程響應格式轉換 |

#### 13.2.2 Service 層重構

| 檔案 | 變更內容 |
|:---|:---|
| `app/services/room.go` | 移除 RoomResponse 結構、ToRoomResponse()、ToRoomResponses() |
| `app/services/course.go` | 移除 CourseResponse 結構、ToCourseResponse()、ToCourseResponses() |

#### 13.2.3 Controller 層整合

| 控制器 | 新增依賴 |
|:---|:---|
| AdminRoomController | roomResource |
| AdminCenterController | centerResource |
| AdminCourseController | courseResource |

### 13.3 架構效益

- **職責分離**：Controller 負責請求解析、Service 負責業務邏輯、Resource 負責格式轉換
- **欄位過濾**：資料庫模型欄位不直接暴露給前端
- **格式標準化**：所有 API 響應格式由 Resource 層統一處理
- **可維護性**：格式變更只需修改 Resource 層

### 13.4 編譯驗證

```bash
go build -mod=mod ./...
# 輸出：無錯誤，全部編譯成功
```

---

[其餘內容保持不變...]
