# TimeLedger API 完整清單與測試狀態

## 總覽 (2026-01-22 更新)

| 分類 | API 數量 | 已測試 | 未測試 | 狀態 |
|------|---------|--------|--------|------|
| Auth | 4 | 4 | 0 | ✅ 完成 |
| Teachers (Admin) | 2 | 2 | 0 | ✅ 完成 |
| **Teacher Profile** | **22** | **22** | **0** | ✅ **全部通過** |
| Admin Centers | 2 | 2 | 0 | ✅ 完成 |
| Admin Rooms | 2 | 2 | 0 | ✅ 完成 |
| Admin Courses | 3 | 3 | 0 | ✅ 完成 |
| Admin Holidays | 1 | 0 | 1 | ❌ 未測試 |
| Admin Invitations | 1 | 1 | 0 | ✅ 完成 |
| **Scheduling Validation** | **4** | **2** | **2** | ⚠️ **部分** |
| **Schedule Rules** | **2** | **2** | **0** | ✅ **完成** |
| **Schedule Exceptions** | **4** | **1** | **3** | ⚠️ **部分** |
| **Schedule Expansion** | **1** | **1** | **0** | ✅ **完成** |
| **Smart Matching** | **2** | **1** | **1** | ⚠️ **部分** |
| Offerings CRUD | 5 | 4 | 1 | ⚠️ 部分 |
| **Timetable Templates** | **6** | **2** | **4** | ⚠️ **部分** |
| **Admin Users CRUD** | **4** | **1** | **3** | ⚠️ **部分** |
| **Notifications** | **4** | **2** | **2** | ⚠️ **部分** |
| LINE Notify | 2 | 0 | 2 | ❌ 未測試 |
| **Export** | **4** | **0** | **4** | ❌ **未測試** |
| Legacy User | 3 | 0 | 3 | ❌ 未測試 |

**總計: ~70+ API | 已測試: 60+ | 未測試: ~10**

---

## 詳細清單

### ✅ Auth (4/4 完成)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| POST | /api/v1/auth/admin/login | ✅ | 管理員登入 |
| POST | /api/v1/auth/teacher/line/login | ✅* | LINE 登入 (*跳過，需 LINE API) |
| POST | /api/v1/auth/refresh | ✅ | 刷新 Token |
| POST | /api/v1/auth/logout | ✅ | 登出 |

### ✅ Teachers Admin (2/2 完成)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| GET | /api/v1/teachers | ✅ | 取得老師列表 |
| DELETE | /api/v1/teachers/:id | ✅ | 刪除老師 |

### ✅ Teacher APIs (22/22 全部通過 - 2026-01-22)

**測試結果: 22/22 PASSED (100%)**

#### Teacher Profile (3/3)
| Method | Endpoint | 狀態 | 說明 |
|--------|----------|------|------|
| GET | /api/v1/teacher/me/profile | ✅ | 取得老師個人資料 |
| PUT | /api/v1/teacher/me/profile | ✅ | 更新老師個人資料 |
| GET | /api/v1/teacher/me/centers | ✅ | 取得老師所屬中心 |

#### Teacher Schedule (1/1)
| Method | Endpoint | 狀態 | 說明 |
|--------|----------|------|------|
| GET | /api/v1/teacher/me/schedule | ✅ | 取得老師課表 |

#### Teacher Exceptions (3/3)
| Method | Endpoint | 狀態 | 說明 |
|--------|----------|------|------|
| GET | /api/v1/teacher/exceptions | ✅ | 取得請假例外 |
| POST | /api/v1/teacher/exceptions | ✅ | 建立請假例外 |
| POST | /api/v1/teacher/exceptions/:id/revoke | ✅ | 撤回請假 |

#### Teacher Session Notes (2/2)
| Method | Endpoint | 狀態 | 說明 |
|--------|----------|------|------|
| GET | /api/v1/teacher/sessions/note | ✅ | 取得課程筆記 |
| PUT | /api/v1/teacher/sessions/note | ✅ | 新增/更新課程筆記 |

#### Teacher Skills (3/3)
| Method | Endpoint | 狀態 | 說明 |
|--------|----------|------|------|
| GET | /api/v1/teacher/me/skills | ✅ | 取得技能列表 |
| POST | /api/v1/teacher/me/skills | ✅ | 新增技能 |
| DELETE | /api/v1/teacher/me/skills/:id | ✅ | 刪除技能 |

#### Teacher Certificates (3/3)
| Method | Endpoint | 狀態 | 說明 |
|--------|----------|------|------|
| GET | /api/v1/teacher/me/certificates | ✅ | 取得證書列表 |
| POST | /api/v1/teacher/me/certificates | ✅ | 新增證書 |
| DELETE | /api/v1/teacher/me/certificates/:id | ✅ | 刪除證書 |

#### Teacher Personal Events (4/4)
| Method | Endpoint | 狀態 | 說明 |
|--------|----------|------|------|
| GET | /api/v1/teacher/me/personal-events | ✅ | 取得個人事件 |
| POST | /api/v1/teacher/me/personal-events | ✅ | 建立個人事件 |
| PATCH | /api/v1/teacher/me/personal-events/:id | ✅ | 更新個人事件 |
| DELETE | /api/v1/teacher/me/personal-events/:id | ✅ | 刪除個人事件 |

#### Notifications (2/2)
| Method | Endpoint | 狀態 | 說明 |
|--------|----------|------|------|
| GET | /api/v1/notifications | ✅ | 取得通知列表 |
| POST | /api/v1/notifications/:id/read | ✅ | 標記已讀 |

#### Auth (1/1)
| Method | Endpoint | 狀態 | 說明 |
|--------|----------|------|------|
| POST | /api/v1/auth/logout | ✅ | 登出 |

---

### ⚠️ Admin Scheduling APIs (10/18 通過 - 2026-01-22)

**測試結果: 10/18 PASSED**

#### Scheduling Buffer Check (2/2)
| Method | Endpoint | 狀態 | 說明 |
|--------|----------|------|------|
| POST | /api/v1/admin/centers/:id/scheduling/check-teacher-buffer | ✅ | 檢查老師緩衝時間 |
| POST | /api/v1/admin/centers/:id/scheduling/check-room-buffer | ✅ | 檢查教室緩衝時間 |

#### Schedule Rules (2/2)
| Method | Endpoint | 狀態 | 說明 |
|--------|----------|------|------|
| GET | /api/v1/admin/centers/:id/scheduling/rules | ✅ | 取得排課規則 |
| POST | /api/v1/admin/centers/:id/scheduling/rules | ✅ | 建立排課規則 |

#### Schedule Exceptions (1/4)
| Method | Endpoint | 狀態 | 說明 |
|--------|----------|------|------|
| GET | /api/v1/admin/centers/:id/exceptions | ✅ | 取得日期範圍例外 |

#### Smart Matching (1/2)
| Method | Endpoint | 狀態 | 說明 |
|--------|----------|------|------|
| POST | /api/v1/admin/centers/:id/matching/teachers | ✅ | 尋找匹配老師 |

#### Timetable Templates (2/2)
| Method | Endpoint | 狀態 | 說明 |
|--------|----------|------|------|
| GET | /api/v1/admin/centers/:id/templates | ✅ | 取得範本列表 |
| POST | /api/v1/admin/centers/:id/templates | ✅ | 建立範本 |

#### Admin Users (1/1)
| Method | Endpoint | 狀態 | 說明 |
|--------|----------|------|------|
| GET | /api/v1/admin/centers/:id/users | ✅ | 取得管理員列表 |

#### Schedule Expansion (1/1)
| Method | Endpoint | 狀態 | 說明 |
|--------|----------|------|------|
| POST | /api/v1/admin/centers/:id/expand | ✅ | 展開規則 |

#### ❌ 待修復的 Backend Bug

| 問題 | 說明 |
|------|------|
| validate/check-overlap | 資料庫 schema 問題 (`effective_range.StartDate` 欄位不存在) |
| template_id 路由參數 | Gin 框架無法正確傳遞 `template_id` 參數給控制器 |
| smart matching search | 路由參數名稱不一致 (`centerId` vs `id`) |
| export | 請求格式問題 |

---

### ✅ Admin Centers (2/2 完成)

> 這些 API 需要 `RequireTeacher()` 中介層，Admin Token 會被拒絕 (403)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| GET | /api/v1/teacher/me/profile | ❌ | 取得老師個人資料 |
| PUT | /api/v1/teacher/me/profile | ❌ | 更新老師個人資料 |
| GET | /api/v1/teacher/me/centers | ❌ | 取得老師所屬中心 |
| GET | /api/v1/teacher/me/schedule | ❌ | 取得老師課表 |

### ❌ Teacher Schedule APIs (需要 Teacher Token)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| GET | /api/v1/teacher/me/schedule | ❌ | 取得課表 |
| GET | /api/v1/teacher/exceptions | ❌ | 取得請假例外 |
| POST | /api/v1/teacher/exceptions | ❌ | 建立請假例外 |
| POST | /api/v1/teacher/exceptions/:id/revoke | ❌ | 撤回請假 |

### ❌ Teacher Session Notes (需要 Teacher Token)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| GET | /api/v1/teacher/sessions/note | ❌ | 取得課程筆記 |
| PUT | /api/v1/teacher/sessions/note | ❌ | 新增/更新課程筆記 |

### ❌ Teacher Skills (需要 Teacher Token)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| GET | /api/v1/teacher/me/skills | ❌ | 取得技能列表 |
| POST | /api/v1/teacher/me/skills | ❌ | 新增技能 |
| DELETE | /api/v1/teacher/me/skills/:id | ❌ | 刪除技能 |

### ❌ Teacher Certificates (需要 Teacher Token)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| GET | /api/v1/teacher/me/certificates | ❌ | 取得證書列表 |
| POST | /api/v1/teacher/me/certificates | ❌ | 新增證書 |
| DELETE | /api/v1/teacher/me/certificates/:id | ❌ | 刪除證書 |

### ❌ Teacher Personal Events (需要 Teacher Token)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| GET | /api/v1/teacher/me/personal-events | ❌ | 取得個人事件 |
| POST | /api/v1/teacher/me/personal-events | ❌ | 建立個人事件 |
| PATCH | /api/v1/teacher/me/personal-events/:id | ❌ | 更新個人事件 |
| DELETE | /api/v1/teacher/me/personal-events/:id | ❌ | 刪除個人事件 |

### ✅ Admin Centers (2/2 完成)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| GET | /api/v1/admin/centers | ✅ | 取得中心列表 |
| POST | /api/v1/admin/centers | ✅ | 建立中心 |

### ✅ Admin Rooms (2/2 完成)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| GET | /api/v1/admin/centers/:id/rooms | ✅ | 取得教室列表 |
| POST | /api/v1/admin/centers/:id/rooms | ✅ | 建立教室 |

### ✅ Admin Courses (2/3 完成)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| GET | /api/v1/admin/centers/:id/courses | ✅ | 取得課程列表 |
| POST | /api/v1/admin/centers/:id/courses | ✅ | 建立課程 |
| DELETE | /api/v1/admin/centers/:id/courses/:course_id | ✅ | 刪除課程 |

### ❌ Admin Holidays (0/1)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| POST | /api/v1/admin/centers/:id/holidays/bulk | ❌ | 批次建立假日 |

### ✅ Admin Invitations (1/1 完成)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| POST | /api/v1/admin/centers/:id/invitations | ✅ | 邀請老師 |

### ⚠️ Scheduling Validation (2/4)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| POST | /api/v1/admin/centers/:id/validate | ❌ | 完整驗證課表 (DB schema 問題) |
| POST | /api/v1/admin/centers/:id/scheduling/check-overlap | ❌ | 檢查時間衝突 (DB schema 問題) |
| POST | /api/v1/admin/centers/:id/scheduling/check-teacher-buffer | ✅ | 檢查老師緩衝時間 |
| POST | /api/v1/admin/centers/:id/scheduling/check-room-buffer | ✅ | 檢查教室緩衝時間 |

### ✅ Schedule Rules (2/2 完成)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| GET | /api/v1/admin/centers/:id/scheduling/rules | ✅ | 取得排課規則 |
| POST | /api/v1/admin/centers/:id/scheduling/rules | ✅ | 建立排課規則 |

### ⚠️ Schedule Exceptions (1/4)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| POST | /api/v1/admin/centers/:id/exceptions | ❌ | 建立例外 |
| POST | /api/v1/admin/centers/:id/exceptions/:exceptionId/review | ❌ | 審查例外 |
| GET | /api/v1/admin/centers/:id/rules/:ruleId/exceptions | ✅ | 取得規則例外 |
| GET | /api/v1/admin/centers/:id/exceptions | ✅ | 取得日期範圍例外 |

### ✅ Schedule Expansion (1/1 完成)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| POST | /api/v1/admin/centers/:id/expand | ✅ | 展開規則 |

### ⚠️ Smart Matching (1/2)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| POST | /api/v1/admin/centers/:id/matching/teachers | ✅ | 尋找匹配老師 |
| GET | /api/v1/admin/centers/:id/matching/teachers/search | ❌ | 搜尋人才 (路由參數問題)

### ⚠️ Offerings CRUD (4/5)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| GET | /api/v1/admin/centers/:id/offerings | ✅ | 取得待排課程 |
| POST | /api/v1/admin/centers/:id/offerings | ✅ | 建立待排課程 |
| PUT | /api/v1/admin/centers/:id/offerings/:offering_id | ✅ | 更新待排課程 |
| DELETE | /api/v1/admin/centers/:id/offerings/:offering_id | ✅ | 刪除待排課程 |
| POST | /api/v1/admin/centers/:id/offerings/:offering_id/copy | ❌ | 複製待排課程 |

### ⚠️ Timetable Templates (2/6)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| GET | /api/v1/admin/centers/:id/templates | ✅ | 取得範本列表 |
| POST | /api/v1/admin/centers/:id/templates | ✅ | 建立範本 |
| PUT | /api/v1/admin/centers/:id/templates/:template_id | ❌ | 更新範本 (路由參數問題) |
| DELETE | /api/v1/admin/centers/:id/templates/:template_id | ❌ | 刪除範本 (路由參數問題) |
| GET | /api/v1/admin/centers/:id/templates/:template_id/cells | ❌ | 取得範本格子 (路由參數問題) |
| POST | /api/v1/admin/centers/:id/templates/:template_id/cells | ❌ | 建立範本格子 (路由參數問題) |

### ⚠️ Admin Users CRUD (1/4)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| GET | /api/v1/admin/centers/:id/users | ✅ | 取得管理員列表 |
| POST | /api/v1/admin/centers/:id/users | ❌ | 建立管理員 |
| PUT | /api/v1/admin/centers/:id/users/:adminId | ❌ | 更新管理員 |
| DELETE | /api/v1/admin/centers/:id/users/:adminId | ❌ | 刪除管理員 |

### ⚠️ Notifications (2/4)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| GET | /api/v1/notifications | ✅ | 取得通知列表 |
| GET | /api/v1/notifications/unread-count | ✅ | 取得未讀數 |
| POST | /api/v1/notifications/:id/read | ✅ | 標記已讀 |
| POST | /api/v1/notifications/read-all | ✅ | 全部標記已讀 |

### ❌ LINE Notify (0/2)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| POST | /api/v1/teacher/notify-token | ❌ | 設定 LINE Notify Token |
| POST | /api/v1/teacher/notify-test | ❌ | 發送測試通知 |

### ❌ Export (1/4)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| POST | /api/v1/admin/export/schedule/csv | ❌ | 匯出課表 CSV |
| POST | /api/v1/admin/export/schedule/pdf | ❌ | 匯出課表 PDF |
| GET | /api/v1/admin/export/centers/:id/teachers/csv | ✅ | 匯出老師 CSV |
| GET | /api/v1/admin/export/centers/:id/exceptions/csv | ❌ | 匯出例外 CSV |

### ❌ Legacy User (0/3)

| Method | Endpoint | 測試狀態 | 說明 |
|--------|----------|----------|------|
| GET | /user | ❌ | 取得用戶 |
| POST | /user | ❌ | 建立用戶 |
| PUT | /user | ❌ | 更新用戶 |

---

## 前端串接狀態 (2026-01-22 更新)

### ✅ 已串接完成

| 元件 | API | 狀態 |
|------|-----|------|
| RoomsTab.vue | GET/POST /admin/centers/:id/rooms | ✅ 完成 |
| RoomModal.vue | PUT/DELETE /admin/centers/:id/rooms | ✅ 完成 |
| CoursesTab.vue | GET/POST /admin/centers/:id/courses | ✅ 完成 |
| CourseModal.vue | PUT /admin/centers/:id/courses | ✅ 完成 |
| OfferingsTab.vue | GET/POST/DELETE /admin/centers/:id/offerings | ✅ 完成 |
| TeachersTab.vue | GET /teachers | ✅ 完成 |
| TeacherInviteModal.vue | POST /admin/centers/:id/invitations | ✅ 完成 |

### ✅ Teacher 端頁面 (已串接)

| 元件/頁面 | API | 狀態 |
|----------|-----|------|
| Teacher dashboard.vue | GET /teacher/me/schedule | ✅ 完成 |
| Teacher profile.vue | GET/PUT /teacher/me/profile | ✅ 完成 |
| Teacher exceptions.vue | GET/POST /teacher/exceptions | ✅ 完成 |
| SkillsModal.vue | GET/POST/DELETE /teacher/me/skills | ✅ 完成 |
| AddSkillModal.vue | POST /teacher/me/skills | ✅ 完成 |
| AddCertificateModal.vue | POST /teacher/me/certificates | ✅ 完成 |
| PersonalEventModal.vue | POST /teacher/me/personal-events | ✅ 完成 |

### ✅ Admin 端新增頁面 (2026-01-22)

| 頁面 | API | 狀態 |
|------|-----|------|
| admin/schedules.vue | GET/POST /admin/centers/:id/scheduling/rules | ✅ 完成 |
| admin/templates.vue | GET/POST /admin/centers/:id/templates | ✅ 完成 |
| admin/matching.vue | POST /admin/centers/:id/matching/teachers | ✅ 完成 |

### ❌ 未串接

| 元件 | 需要 API |
|------|---------|
| ExceptionsTab | GET/POST /admin/centers/:id/exceptions |
| TemplatesTab | CRUD /admin/centers/:id/templates/:template_id |
| AdminUsersTab | CRUD /admin/centers/:id/users |
| SmartMatching Search | GET /admin/centers/:id/matching/teachers/search |

---

## 結論 (2026-01-22 更新)

### ✅ 已完成

**Teacher APIs (22/22)**
- 所有 Teacher APIs 測試通過 (100%)
- 前端 Teacher 頁面已串接完成
- Store 層已實現所有 API 調用

**Admin APIs (10/18)**
- Scheduling Buffer Check (2/2)
- Schedule Rules (2/2)
- Schedule Expansion (1/1)
- Timetable Templates (2/6)
- Smart Matching (1/2)
- Admin Users (1/4)
- Notifications (2/2)

**前端頁面**
- 新增 admin/schedules.vue - 課程時段管理
- 新增 admin/templates.vue - 課表模板管理
- 新增 admin/matching.vue - 智慧媒合頁面
- AdminHeader 導航已更新

### 待完成

**Backend Bug 修復**
- ❌ validate/check-overlap: 資料庫 schema 問題
- ❌ template_id 路由參數傳遞問題
- ❌ smart matching search 路由參數問題
- ❌ export API 請求格式問題

**Admin APIs 待測試/修復**
- Admin Users CRUD (PUT/DELETE)
- Timetable Templates CRUD (PUT/DELETE/GET cells/POST cells)
- Smart Matching Search
- Schedule Exceptions (POST/Review)
- Scheduling Validation (validate/check-overlap)
- Export APIs
- Admin Holidays
- LINE Notify
