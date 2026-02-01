# 全局錯誤處理系統

## 概述

提供完整的錯誤處理解決方案，整合錯誤碼對應、使用者友善訊息和 UI 提示（Alert / Toast）。

## 核心功能

- **錯誤碼對應**：將後端 API 錯誤碼映射為使用者友善訊息
- **多種顯示方式**：支援 Alert（阻斷式）和 Toast（非阻斷式）
- **情境處理**：針對權限、驗證、排課衝突等常見錯誤提供專門處理
- **表單驗證**：輕鬆處理表單欄位錯誤
- **佇列處理**：支援錯誤佇列，避免同時顯示多個錯誤

## 檔案結構

```
frontend/
├── composables/
│   ├── useErrorHandler.ts    # 核心錯誤處理（完整功能）
│   ├── useQuickError.ts      # 簡化版錯誤處理（推薦使用）
│   └── index.ts              # 統一匯出
├── plugins/
│   └── error-handler.client.ts  # Nuxt 插件（自動註冊）
└── constants/
    └── errorCodes.ts         # 錯誤碼對照表
```

## 快速開始

### 方式一：使用簡化版 Composable（推薦）

```typescript
// 在任何 Vue 元件中使用
const {
  showError,
  showSuccess,
  showWarning,
  request,
  handle,
} = useQuickError()
```

### 方式二：使用 Nuxt Plugin（自動註冊）

```typescript
// 在任何元件中都可以使用
const { $error, $quickError } = useNuxtApp()

// 處理 API 回應
const result = await $quickError.request(() => api.getUser(id))
```

## 使用範例

### 1. 顯示成功訊息

```typescript
await showSuccess('資料已成功儲存', '儲存成功')
```

### 2. 顯示錯誤訊息

```typescript
await showError('發生錯誤，請稍後再試', '操作失敗')
```

### 3. 顯示確認對話框

```typescript
const shouldDelete = await confirm('確定要刪除嗎？此操作無法復原。', '確認刪除')

if (shouldDelete) {
  // 使用者點擊「確認」
  await showSuccess('刪除成功')
}
```

### 4. 處理 API 回應（自動顯示錯誤）

```typescript
const data = await request(() => api.getUser(id))

if (!data) {
  // 錯誤已自動處理並顯示
  return
}

// 成功處理
console.log(data)
```

### 5. 處理 API 回應（手動控制）

```typescript
const response = await api.getUser(id)
const result = await handle(response)

if (!result.success) {
  // 自訂錯誤處理邏輯
  await showError(result.error)
  return
}

// 成功處理
console.log(result.data)
```

### 6. 處理權限錯誤

```typescript
await handlePermission('FORBIDDEN', '您沒有權限執行此操作')
// 會自動提供「前往登入」按鈕（如果是未登入狀態）
```

### 7. 處理驗證錯誤

```typescript
handleValidation({
  name: ['姓名必填'],
  email: ['電子郵件格式錯誤'],
})
// 會顯示警告訊息，列出所有錯誤
```

### 8. 處理排課衝突

```typescript
await handleConflicts([
  {
    type: 'TEACHER',
    message: '老師張三在 14:00-15:00 已有課程',
  },
  {
    type: 'ROOM',
    message: '教室 A 在 14:00-15:00 已被預約',
  },
])
```

### 9. 表單驗證

```typescript
const {
  errors,
  setFieldError,
  setErrors,
  clearErrors,
  getFieldError,
} = useFormErrors()

// 設置欄位錯誤
setFieldError('email', '電子郵件格式錯誤')

// 取得欄位錯誤
const error = getFieldError('email')

// 清除所有錯誤
clearErrors()
```

### 10. 異步操作（自動錯誤處理）

```typescript
const { loading, execute } = useAsyncOperation()

const result = await execute(async () => {
  await api.saveData(data)
  return '儲存成功'
})

if (result.success) {
  await showSuccess(result.data)
}
```

## API 參考

### useQuickError()

| 方法 | 說明 |
|:---|:---|
| `showError(message, title?)` | 顯示錯誤訊息（Alert） |
| `showSuccess(message, title?)` | 顯示成功訊息（Alert） |
| `showWarning(message, title?)` | 顯示警告訊息（Alert） |
| `showInfo(message, title?)` | 顯示資訊訊息（Alert） |
| `confirm(message, title?)` | 顯示確認對話框 |
| `request(fn, options?)` | 執行 API 請求並自動處理錯誤 |
| `handle(response, options?)` | 處理 API 回應 |
| `handlePermission(code, message, options?)` | 處理權限錯誤 |
| `handleValidation(errors, options?)` | 處理驗證錯誤 |
| `handleConflicts(conflicts, options?)` | 處理排課衝突 |

### useFormErrors()

| 方法 | 說明 |
|:---|:---|
| `errors` | 所有錯誤的響應式物件 |
| `hasErrors` | 是否有任何錯誤 |
| `setFieldError(field, message)` | 設置單一欄位錯誤 |
| `setErrors(errors)` | 設置多個欄位錯誤 |
| `clearFieldError(field)` | 清除單一欄位錯誤 |
| `clearErrors()` | 清除所有錯誤 |
| `getFieldError(field)` | 取得欄位錯誤訊息 |

### useAsyncOperation()

| 屬性/方法 | 說明 |
|:---|:---|
| `loading` | 是否正在載入（唯讀） |
| `execute(fn)` | 執行異步操作，自動捕獲錯誤 |

## 錯誤碼對照表

系統支援以下錯誤碼分類：

| 分類 | 前綴 | 範例 |
|:---|:---|:---|
| 系統錯誤 | 1xxxx | `SYSTEM_ERROR` |
| 資料庫錯誤 | 12xxxx | `SQL_ERROR`, `NOT_FOUND` |
| 權限錯誤 | 4xxxx | `UNAUTHORIZED`, `FORBIDDEN` |
| 驗證錯誤 | 4xxxx | `VALIDATION_ERROR` |
| 業務邏輯錯誤 | 4xxxx | `SCHEDULE_CONFLICT` |
| 例外審核錯誤 | 4xxxx | `EXCEPTION_NOT_FOUND` |
| 教師相關錯誤 | 4xxxx | `TEACHER_NOT_FOUND` |
| 課程相關錯誤 | 4xxxx | `COURSE_NOT_FOUND` |
| 教室相關錯誤 | 4xxxx | `ROOM_NOT_FOUND` |
| 通知相關錯誤 | 4xxxx | `NOTIFICATION_SEND_FAILED` |

## 在 Axios 中使用

建立錯誤處理攔截器：

```typescript
// plugins/axios.ts
export default defineNuxtPlugin((nuxtApp) => {
  const { handle } = useQuickError()

  $axios.interceptors.response.use(
    (response) => response,
    async (error) => {
      const apiResponse = {
        code: error.response?.data?.code || error.response?.status?.toString(),
        message: error.response?.data?.message || '發生錯誤',
        data: error.response?.data?.data,
      }

      await handle(apiResponse)
      return Promise.reject(error)
    }
  )
})
```

## 自訂錯誤碼

如果需要新增自訂錯誤碼，在 `frontend/constants/errorCodes.ts` 中擴展：

```typescript
export const BUSINESS_ERROR_CODES = {
  // ... 現有錯誤碼
  CUSTOM_ERROR: '您的自訂錯誤訊息',
} as const
```

## 注意事項

1. **使用簡化版**：推薦使用 `useQuickError()`，介面更簡潔
2. **避免原生 alert/confirm**：一律使用本系統的錯誤處理
3. **錯誤碼完整性**：新增功能時請同步更新錯誤碼對照表
4. **國際化支援**：如需多語言，可在錯誤碼對照表中擴展
