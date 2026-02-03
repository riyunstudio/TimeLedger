# 老師註冊與邀請流程修復 - 完工紀錄

## 修復日期
2026年2月3日

## 問題描述
新老師透過邀請連結註冊時，系統會回傳 `teacher not found` 錯誤，導致無法完成註冊流程。這是因為系統預期老師必須已存在於資料庫中，但實際上新老師尚未有任何帳號記錄。

## 修復方案

### 1. 後端修改

#### 1.1 自動建立老師帳號邏輯 (`app/services/teacher.go`)

**修改位置**：`AcceptInvitationByLink` 函式

**修改內容**：
當透過 `LineUserID` 找不到老師時，系統會自動建立新的老師帳號：

```go
// 取得老師資料
teacher, err := s.teacherRepo.GetByLineUserID(ctx, req.LineUserID)
if err != nil {
    // 新老師：自動建立老師帳號
    if invitation.Email != "" {
        newTeacher := models.Teacher{
            LineUserID: req.LineUserID,
            Email:      invitation.Email,
            Name:       "新老師", // 預設名稱
        }

        createdTeacher, err := s.teacherRepo.Create(ctx, newTeacher)
        if err != nil {
            return nil, s.app.Err.New(errInfos.SQL_ERROR), fmt.Errorf("failed to create teacher: %w", err)
        }
        teacher = createdTeacher
    } else {
        return nil, s.app.Err.New(errInfos.NOT_FOUND), fmt.Errorf("teacher not found")
    }
}
```

**邏輯說明**：
- 嘗試透過 `LineUserID` 查找已存在的老師
- 若找不到老師且邀請含有 Email，自動建立新老師
- 新老師使用邀請中的 Email 作為帳號 Email
- 預設名稱設為「新老師」（老師可後續自行修改）
- 建立成功後繼續完成 CenterMembership 建立流程

### 2. 前端修改

#### 2.1 邀請接受 API 傳/invite/[token參 (`frontend/pages].vue`)

**修改位置**：`acceptInvitation` 函式

**修改內容**：
將 API 請求中的 `id_token` 改為 `line_user_id`：

```typescript
// 修改前
body: JSON.stringify({
  id_token: authStore.user.id_token,
})

// 修改後
body: JSON.stringify({
  line_user_id: authStore.user.line_user_id,
})
```

**對應的驗證檢查也一併修改**：

```typescript
// 修改前
if (!authStore.user?.id_token) {
  error.value = '無法取得登入資訊，請重新登入'
  return
}

// 修改後
if (!authStore.user?.line_user_id) {
  error.value = '無法取得登入資訊，請重新登入'
  return
}
```

## 流程說明

### 修復後的新老師註冊流程

```
1. 管理員產生邀請連結（包含 Email）
2. 新老師開啟邀請頁面
3. 點擊「LINE 快速登入」
4. LINE 登入成功後取得 line_user_id
5. 點擊「接受邀請」
6. 後端收到 line_user_id
7. 查找老師（找不到）
8. 自動建立新老師帳號（使用邀請中的 Email）
9. 建立 CenterMembership
10. 回傳成功訊息
```

## 檔案變更清單

| 檔案 | 變更類型 | 說明 |
|:---|:---:|:---|
| `app/services/teacher.go` | 修改 | 新增自動建立老師邏輯 |
| `frontend/pages/invite/[token].vue` | 修改 | 變更 API 傳參欄位 |

## 驗證結果

- **編譯狀態**：`go build ./...` 通過
- **API 相容性**：
  - `AcceptInvitationByLinkRequest.LineUserID` 欄位已設定 `binding:"required"`
  - 前端傳送 `line_user_id` 符合後端預期

## 待測試項目

1. **管理員產生邀請連結**
   - [ ] 確認邀請連結可以正確產生
   - [ ] 確認邀請 Email 正確儲存

2. **新老師透過連結註冊**
   - [ ] 新老師開啟邀請頁面
   - [ ] LINE 登入成功
   - [ ] 點擊接受邀請
   - [ ] 確認資料庫中已建立 Teacher 記錄
   - [ ] 確認 CenterMembership 狀態正確

3. **已存在的老師接受邀請**
   - [ ] 系統正常識別已存在的老師
   - [ ] 正確更新為 ALREADY_MEMBER 狀態

## 風險評估

| 風險項目 | 風險等級 | 緩解措施 |
|:---|:---:|:---|
| Email 衝突 | 低 | GORM 唯一索引會阻止重複 Email |
| 邀請無 Email | 低 | 回傳 error，不自動建立 |
| 並發建立 | 低 | 資料庫交易確保資料一致性 |

## 相關文件

- API 文件：`app/controllers/teacher_invitation.go`
- 模型定義：`app/models/teacher.go`
- 前端頁面：`frontend/pages/invite/[token].vue`
