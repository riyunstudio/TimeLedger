# 修復老師註冊與邀請流程 - 詳細任務清單

此文件為實作修復老師註冊流程的具體步驟清單，旨在解決新老師無法透過邀請連結註冊的問題。

## 1. 後端修改 (Backend)

### Task 1.1: 更新 Request 結構
- **檔案**: [teacher_invitation.go](file:///d:/project/TimeLedger/app/controllers/teacher_invitation.go)
- **修改點**: 修改 `AcceptInvitationByLinkRequest`。
- **目標**: 將 `LineUserID` 設為必要欄位，並確保描述符合用途。

### Task 1.2: 實作自動建立老師邏輯
- **檔案**: [teacher.go](file:///d:/project/TimeLedger/app/services/teacher.go)
- **修改點**: 修改 `AcceptInvitationByLink` 函式。
- **邏輯**:
    - 在取得老師資料失敗時（`err != nil`），加入邏輯：
    - 從 `CenterInvitation` 取得 Email。
    - 建立 `models.Teacher` 實例，設定 `LineUserID`、`Email` 及預設 `Name`。
    - 呼叫 `s.teacherRepo.Create` 寫入資料庫。
    - 確保新建立的 `teacher` 物件被用於後續的 `CenterMembership` 建立。

### Task 1.3: 編譯驗證
- **指令**:
  ```powershell
  go build ./...
  ```
- **目標**: 確保沒有語法錯誤或引用錯誤。

## 2. 前端修改 (Frontend)

### Task 2.1: 更新邀請接受傳參
- **檔案**: [invite/[token].vue](file:///d:/project/TimeLedger/frontend/pages/invite/%5Btoken%5D.vue)
- **修改點**: `acceptInvitation` 函式。
- **目標**: 將 POST body 中的欄位從 `id_token` 改為 `line_user_id`。
- **資料來源**: `authStore.user.line_user_id`。

## 3. 驗證與測試 (Verification)

### Task 3.1: 手動功能測試
1. **邀請生成**: 以管理員身分登入，生成一個新的老師邀請連結。
2. **註冊流程**: 使用無痕視窗或模擬環境，開啟該連結並點擊「LINE 登入」。
3. **資料確認**: 接受邀請後，確認資料庫中已建立新的 `Teacher` 記錄，且 `CenterMembership` 狀態正確。

## 4. 完工紀錄
- [ ] 修改 `app/controllers/teacher_invitation.go`
- [ ] 修改 `app/services/teacher.go`
- [ ] 修改 `frontend/pages/invite/[token].vue`
- [ ] 執行 `go build` 驗證成功
- [ ] 更新 `task.md` 與產出 `walkthrough.md`
## 5. LINE Bot 自主註冊 (Public Registration)

### Task 5.1: 後端 API 實作
- **檔案**: `app/controllers/teacher.go`
- **修改點**: 新增 `PublicRegister` 方法。
- **邏輯**: 接收 `LineUserID`, `Name`, `Email`。檢查是否已存在，若不存在則建立老師記錄並回傳 JWT Token。
- **檔案**: `app/services/teacher.go`
- **修改點**: 新增 `RegisterPublic` 方法處理寫入與 Token 產生。
- **檔案**: `app/servers/route.go`
- **修改點**: 註冊 `POST /api/v1/teacher/public/register`。

### Task 5.2: 前端註冊頁面
- **檔案**: `frontend/pages/teacher/register.vue` (新檔案)
- **邏輯**:
    - 使用 LIFF 取得 `line_user_id`。
    - 表單收集姓名與 Email。
    - 提交後導向 `/teacher/dashboard`。

### Task 5.3: LINE Rich Menu 配置
- **動作**: 提供 Rich Menu 設定 JSON 範例，連結至 `${liffUrl}/teacher/register`。

## 6. 官網開放註冊頁面 (/register)

### Task 6.1: 前端註冊頁面實作
- **檔案**: [NEW] `frontend/pages/register.vue`
- **邏輯**:
    - 提供明顯的「使用 LINE 註冊」按鈕。
    - 點擊按鈕後導向 LINE OAuth2 授權頁面。
    - 伺服器回傳後，從 URL 或 Callback 取得 `line_user_id`。
    - 顯示姓名與 Email 表單（可複用 `teacher/register.vue` 的邏輯）。
    - 提交後呼叫後端 API：`POST /api/v1/teacher/public/register`。

### Task 6.2: LINE Login 環境變數配置
- **動作**: 確保前端配置了正確的 `LINE_CLIENT_ID` 與 `LINE_REDIRECT_URI`。

## 7. 中心通用邀請連結 (General Center Invitation Link)

### Task 7.1: 修改後端模型與邏輯
- **檔案**: `app/models/center_invitation.go`
- **修改點**: 新增 `InvitationTypeGeneral = "GENERAL"`。
- **檔案**: `app/services/teacher.go`
- **修改點**: 
    - 修改 `AcceptInvitationByLink`：判斷若為 `GENERAL`，不更新邀請狀態為 `ACCEPTED`，且跳過 Email 檢查。
    - 新增 `GenerateGeneralInvitationLink`：支援產生不指定 Email 的通用連結。
- **檔案**: `app/controllers/teacher_invitation.go`
- **修改點**: 新增 API 端點支援產生通用連結。

### Task 7.2: 前端管理介面實作
- **檔案**: `frontend/pages/admin/invitations.vue`
- **修改點**: 
    - 新增「通用連結管理」區塊。
    - 實作「啟用/停用」及「重新產生」按鈕。
