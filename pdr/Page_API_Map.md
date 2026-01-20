# 頁面與 API 對照表 (Page-API Mapping & Gap Analysis)

本文件將前端頁面 (Routes) 逐一列出，並對照所需 API，藉此盤點後端介面缺失。

---

## 1. 老師端頁面 (Teacher - /teacher/*)
*風格：Mobile First, 極簡白底, 底部導航 (Bottom Nav)*

| 頁面 (Route) | 畫面功能 element | 需呼叫 API | 狀態 (Gap Check) |
|:---|:---|:---|:---|
| `/auth/login` | LINE 登入按鈕 | `GET /auth/line/login`<br>`POST /auth/line/callback` | ✅ OK |
| `/teacher/dashboard` | **週/3日課表** | `GET /teacher/me/schedule` (含 range) | ✅ OK |
| | 切換中心 (Filter) | `GET /teacher/me/centers` | ✅ OK |
| | **新增個人行程** (Fab) | `POST /teacher/personal-events` | ✅ OK |
| | 課表細節 (Detail) | `GET /teacher/sessions/note` (筆記) | ✅ OK |
| | 匯出圖片 (Export) | (前端 html2canvas 處理) | N/A |
| `/teacher/profile` | 編輯個人檔案 | `GET /teacher/me/profile` (讀)<br>`POST /teacher/me/profile` (寫) | ✅ OK |
| | 證照列表與上傳 | `GET /teacher/me/certificates`<br>`POST ...` | ✅ OK (GET/POST/DELETE) |
| `/teacher/settings` | 綁定新中心 (Input Code) | `POST /teacher/join` | ✅ OK |
| | 登出 | (前端清除 Token) | N/A |

---

## 3. 公開頁面 (Public - /)

| 頁面 (Route) | 畫面功能 element | 需呼叫 API | 狀態 (Gap Check) |
|:---|:---|:---|:---|
| `/` | 登陸頁 (Branding) | N/A | ✅ OK |
| | **互動 Demo (Sandbox)** | `POST /public/validate` (匿名) | ✅ OK |

---

## 2. 中心後台 (Admin - /admin/*)
*風格：Desktop Dashboard, 資訊密度高, 側邊欄 (Sidebar)*

| 頁面 (Route) | 畫面功能 element | 需呼叫 API | 狀態 (Gap Check) |
|:---|:---|:---|:---|
| `/admin/login` | 帳密登入表單 | `POST /auth/admin/login` | ✅ OK |
| `/admin/dashboard` | **排課網格 (Grid)** | `GET /admin/centers/{id}/schedule` | ✅ OK |
| | 拖曳排課 (Save/Validate) | `POST /admin/centers/{id}/validate`<br>`POST .../rules/bulk` | ✅ OK |
| | **智慧代課 (Drawer)** | `GET .../matching/teachers` | ✅ OK |
| `/admin/talent` | **人才庫搜尋 (Search)** | `GET /admin/talent/search` | ✅ OK |
| `/admin/teachers` | 老師列表 (CRUD) | `GET /admin/centers/{id}/teachers`<br>`POST .../invite` | ✅ OK |
| | 查看老師證照 | `GET .../teachers/{tid}/certificates` | ✅ OK |
| `/admin/rooms` | 教室列表 | `GET /admin/centers/{id}/rooms`<br>`POST/PUT/DELETE` | ✅ OK |
| `/admin/courses` | 課程/班別管理 | `GET .../courses` (課程)<br>`GET .../offerings` (班別) | ✅ OK |
| `/admin/exceptions` | **異動審核列表** | `GET /admin/centers/{id}/exceptions`<br>`POST .../review` | ✅ OK |
| `/admin/settings` | 政策與中心資訊 | `GET /admin/centers/{id}`<br>`PATCH ...` | ✅ OK |
| `/admin/users` | **管理員帳號管理** | `GET /admin/centers/{id}/users`<br>`POST ...` | ✅ OK |
| `/admin/audit` | **操作稽核紀錄** | `GET /admin/centers/{id}/audit-logs` | ✅ OK |

---

### 3.1 已補完清單 (Status Update)
*   **老師端**: `GET /teacher/me/profile`, `GET /teacher/me/centers` 已加入 `API.md`。
*   **後台 CRUD**: 所有課程、教室、老師的 CRUD API 已加入 `API.md`。
*   **排課範本**: `timetable_templates` 與 `cells` API 已加入 `API.md`。
