# 開發階段總結

## 修復的問題

### 1. API URL 重複問題

**問題描述**：`nuxt.config.ts` 中 `apiBase` 已設為 `/api/v1`，但前端代碼又再次添加前綴，導致 URL 重複。

**修復檔案**：
- `frontend/pages/admin/settings.vue`
- `frontend/pages/admin/admin-list.vue`
- `frontend/pages/admin/line-bind.vue`

---

### 2. JWT Token 角色問題

**問題描述**：`auth.go` 中的 `AdminLogin` 函數將 `UserType` 硬編碼為 `"ADMIN"`，忽略了資料庫中儲存的实际角色。

**修復檔案**：
- `app/middleware/auth.go`

**修復內容**：
- `AdminLogin` 函數：修正為使用資料庫中的實際角色 (`admin.Role`)
- `RefreshToken` 函數：修正相同的問題

---

### 3. Middleware 權限問題

**問題描述**：`RequireAdmin()` 和 `RequireCenterAdmin()` 只檢查 `userType != "ADMIN"`，导致 OWNER 角色被錯誤拒絕。

**修復檔案**：
- `app/middleware/auth.go`

**修正後**：
```go
if userType != "ADMIN" && userType != "OWNER" {
  // 拒絕存取
}
```

---

### 4. Context Key 類型問題

**問題描述**：使用 `ctx.Get(string(global.UserIDKey))` 多餘的 string 轉換，導致 context key 匹配問題。

**修復檔案**：
- `app/controllers/admin_user.go`
- `app/controllers/line_bot.go`

**修正後**：
```go
ctx.Get(global.UserIDKey)  // 不需要 string() 轉換
```

---

### 5. 錯誤碼定義

**修復檔案**：
- `global/errInfos/code.go`

**新增錯誤碼**：
| 錯誤碼 | 常數名稱 | 說明 |
|:---:|:---|:---|
| 100002 | `ADMIN_EMAIL_EXISTS` | 管理員電子郵件已存在 |

---

### 6. 前端類型定義

**問題描述**：後端 `IdentInfo` 使用 `user_type` 欄位，但前端 `AdminUser` 類型使用 `role` 欄位。

**修復檔案**：
- `frontend/types/index.ts`
- `frontend/stores/auth.ts`

**修正內容**：
- `AdminUser` 接口添加 `user_type` 欄位
- `isAdmin` computed 同時檢查 `user_type` 和 `role` 欄位

---

## 新增功能

### 1. OWNER 角色快速登入按鈕

**修復檔案**：
- `frontend/pages/admin/login.vue`

**新增OWNER、管理員、工作人員三種角色的快速登入按鈕**：

| 按鈕名稱 | 角色 | Email | 密碼 |
|:---|:---|:---|:---|
| 擁有者 | OWNER | admin@timeledger.com | admin123 |
| 管理員 | ADMIN | manager@timeledger.com | admin123 |
| 工作人員 | STAFF | staff@timeledger.com | admin123 |

---

### 2. 直接新增管理員功能

**問題描述**：原本需要產生邀請碼，現在改為直接新增管理員並設置密碼。

**新增 API**：
```
POST /api/v1/admin/admins
{
  "email": "admin@example.com",
  "name": "新管理員",
  "role": "ADMIN",
  "password": "admin123"
}
```

**新增檔案**：
- `app/controllers/admin_user.go`：`CreateAdmin` 函數
- `app/services/admin_user.go`：`CreateAdmin` Service 函數
- `app/servers/route.go`：新增路由
- `frontend/pages/admin/admin-list.vue`：表單新增密碼欄位

**資料庫 SQL（OWNER 帳號）**：

```sql
INSERT INTO admin_users (
    center_id,
    email,
    password,
    name,
    role,
    status,
    created_at,
    updated_at
) VALUES (
    1,
    'admin@timeledger.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.rsS1C1.E/6surBGHqW',  -- 密碼 "admin123"
    '系統擁有者',
    'OWNER',
    'ACTIVE',
    NOW(),
    NOW()
);
```

---

## 修改檔案清單

| 檔案 | 修改類型 | 說明 |
|:---|:---|:---|
| `frontend/pages/admin/settings.vue` | 修改 | 移除重複的 API URL 前綴 |
| `frontend/pages/admin/admin-list.vue` | 修改 | 移除重複的 API URL 前綴、新增管理員功能 |
| `frontend/pages/admin/line-bind.vue` | 修改 | 移除重複的 API URL 前綴 |
| `frontend/pages/admin/login.vue` | 修改 | 新增 OWNER 快速登入按鈕 |
| `app/middleware/auth.go` | 修改 | JWT Token 角色修正、Middleware 權限修正 |
| `app/controllers/admin_user.go` | 修改 | Context Key 修正、新增 CreateAdmin |
| `app/controllers/line_bot.go` | 修改 | Context Key 修正 |
| `app/services/admin_user.go` | 修改 | 新增 CreateAdmin、bcrypt 加密 |
| `app/servers/route.go` | 修改 | 新增 POST /admin/admins 路由 |
| `global/errInfos/code.go` | 修改 | 新增 ADMIN_EMAIL_EXISTS 錯誤碼 |
| `frontend/types/index.ts` | 修改 | AdminUser 添加 user_type 欄位 |
| `frontend/stores/auth.ts` | 修改 | isAdmin 同時檢查 user_type 和 role |

---

## 驗證項目

- [x] Go 程式可正常編譯
- [x] API URL 請求正常
- [x] OWNER 角色可正常登入並存取管理員功能
- [x] ADMIN 角色權限正確
- [x] STAFF 角色權限正確
- [x] 錯誤碼可正常回傳
- [x] 新增管理員功能正常運作
