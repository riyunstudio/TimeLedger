# TimeLedger 前端重構指令清單（技術優先順序）

> 生成日期：2026-01-31
> 修訂日期：2026-01-31
> 用途：提供給 AI 助手的標準化重構指令
> 核心理念：**先建基建，後補功能** — 好的架構讓後續開發效率倍增
> 數量：12 項重構任務

---

## 重構哲學

在開始任何功能開發之前，必須先建立穩固的基礎設施。否則：
- 每個功能都要重新處理型別、錯誤處理、loading 狀態
- 程式碼不一致，維護困難
- 新人接手需要更多時間理解
- 技術債務會越積越多

**正確順序**：基礎設施 → 緊急修復 → 核心功能 → 管理功能

---

## 第一部分：基礎設施（最優先）

這部分是所有後續開發的基礎，請在開始任何功能前完成。

### 指令 1：統一 API Response 型別定義

**指令內容**：

```
請在 TimeLedger 前端建立統一的 API Response 型別定義系統：

目標：建立完整的 TypeScript 型別系統，作為所有 API 呼叫的基礎

理由：沒有統一的型別，後續每個功能都要重新定義，無法複用

需求：

1. 建立 `frontend/types/api.ts` 通用類型：
   ```typescript
   interface ApiResponse<T> {
     code: string      // 錯誤碼，如 "SUCCESS", "SQL_ERROR"
     message: string   // 訊息
     data?: T          // 單筆資料
     datas?: T         // 多筆資料（部分 API 使用）
   }

   interface PaginationParams {
     page?: number
     limit?: number
     sort_by?: string
     sort_order?: 'ASC' | 'DESC'
   }

   interface PaginationResult {
     page: number
     limit: number
     total: number
     total_pages: number
     has_next: boolean
     has_prev: boolean
   }
   ```

2. 建立模組化型別檔案結構：
   - `frontend/types/api.ts` - API 通用類型
   - `frontend/types/admin.ts` - 管理員相關
   - `frontend/types/teacher.ts` - 教師相關
   - `frontend/types/center.ts` - 中心相關
   - `frontend/types/scheduling.ts` - 排課相關
   - `frontend/types/matching.ts` - 智慧媒合相關
   - `frontend/types/notification.ts` - 通知相關

3. 掃描現有程式碼，找出所有 API 相關的 any 型別並替換

4. 定義常用基礎類型：
   ```typescript
   type ID = number
   type Timestamp = string // ISO 8601 格式
   type DateString = string // YYYY-MM-DD 格式

   interface BaseModel {
     id: ID
     created_at: Timestamp
     updated_at: Timestamp
   }

   interface PaginatedResponse<T> {
     data: T[]
     pagination: PaginationResult
   }
   ```

5. 建立型別匯出 index：
   ```typescript
   // frontend/types/index.ts
   export * from './api'
   export * from './admin'
   export * from './teacher'
   // ... 其他匯出
   ```

驗收標準」：
- [x] `frontend/types/api.ts` 已建立
- [x] 模組化型別檔案已建立
- [x] 所有 Store 中的 API 回傳有正確的型別
- [x] any 型別數量顯著減少
- [x] 新增功能時可直接匯入現有型別

**完成時間**：2026-01-31

**建立檔案**：
| 檔案 | 說明 |
|:---|:---|
| `frontend/types/api.ts` | API 通用類型、分頁、驗證結果等基礎類型 |
| `frontend/types/admin.ts` | 管理員用戶、認證、LINE 綁定等相關類型 |
| `frontend/types/teacher.ts` | 教師、技能、證照、標籤、邀請等相關類型 |
| `frontend/types/center.ts` | 中心、課程、教室、方案、假日等相關類型 |
| `frontend/types/scheduling.ts` | 課表規則、例外、個人行程、課堂筆記等相關類型 |
| `frontend/types/matching.ts` | 智慧媒合、人才庫、替代時段等相關類型 |
| `frontend/types/notification.ts` | 通知、通知佇列、LINE 通知等相關類型 |
| `frontend/types/index.ts` | 統一匯出中心，向後相容性類型 |

---

### 指令 2：統一錯誤處理機制

**指令內容**：

```
請在 TimeLedger 前端建立統一的錯誤處理機制：

目標：建立一致的錯誤處理模式，讓所有 API 錯誤都能被妥善處理

理由：沒有統一的錯誤處理，會導致錯誤訊息不一致、使用者體驗差

需求：

1. 建立錯誤碼對照表 `frontend/constants/errorCodes.ts`：
   ```typescript
   export const ERROR_MESSAGES: Record<string, string> = {
     SUCCESS: '操作成功',
     SQL_ERROR: '系統錯誤，請稍後再試',
     NOT_FOUND: '找不到請求的資源',
     FORBIDDEN: '您沒有權限執行此操作',
     UNAUTHORIZED: '請先登入',
     VALIDATION_ERROR: '輸入資料驗證失敗',
     // ... 其他錯誤碼
   }
   ```

2. 建立錯誤處理工具 `frontend/utils/errorHandler.ts`：
   ```typescript
   interface ApiError {
     code: string
     message: string
   }

   class ErrorHandler {
     // 根據錯誤碼取得訊息
     static getMessage(error: ApiError): string {
       return ERROR_MESSAGES[error.code] || error.message || '發生未知錯誤'
     }

     // 處理錯誤並顯示
     static handle(error: ApiError): void {
       const message = this.getMessage(error)
       // 使用 toast 或 alert 顯示
       console.error('API Error:', error)
     }

     // 記錄錯誤到監控服務
     static log(error: ApiError, context?: object): void {
       // 送到監控服務
     }

     // 取得 HTTP 狀態碼對應的處理方式
     static getHandler(statusCode: number): (error: ApiError) => void {
       const handlers: Record<number, (error: ApiError) => void> = {
         401: (error) => {
           // 未授權，導向登入頁
           window.location.href = '/login'
         },
         403: (error) => {
           // 禁止存取，顯示權限錯誤
           alert('您沒有權限執行此操作')
         },
         404: (error) => {
           // 找不到，顯示提示
           alert('找不到請求的資源')
         },
         500: (error) => {
           // 伺服器錯誤
           alert('系統錯誤，請稍後再試')
         }
       }
       return handlers[statusCode] || this.handle
     }
   }
   ```

3. 設定 Axios Interceptor `frontend/plugins/api.ts`：
   ```typescript
   const api = axios.create({
     baseURL: '/api/v1',
     timeout: 30000
   })

   // 請求攔截器
   api.interceptors.request.use(
     (config) => {
       // 加入 JWT token
       const token = useCookie('token')
       if (token.value) {
         config.headers.Authorization = `Bearer ${token.value}`
       }
       return config
     },
     (error) => Promise.reject(error)
   )

   // 響應攔截器
   api.interceptors.response.use(
     (response) => {
       const { code, message, data, datas } = response.data

       if (code !== 'SUCCESS') {
         return Promise.reject({ code, message })
       }

       return { data, datas, pagination: response.data.pagination }
     },
     (error) => {
       const { status } = error.response || {}
       const handler = ErrorHandler.getHandler(status || 500)
       handler(error.response?.data || { code: 'UNKNOWN', message: '網路錯誤' })
       return Promise.reject(error)
     }
   )
   ```

4. 建立通用錯誤提示組件 `frontend/components/ErrorToast.vue`：
   - 支援顯示錯誤碼和訊息
   - 自動倒數關閉
   - 支援重試按鈕

驗收標準」：
- [x] 錯誤碼對照表已建立 (`frontend/constants/errorCodes.ts`)
- [x] 錯誤處理工具已實作 (`frontend/utils/errorHandler.ts`)
- [x] API 統一處理機制已建立 (`frontend/composables/useApi.ts`)
- [x] 所有 API 錯誤都有適當的錯誤訊息
- [x] 401 錯誤自動導向登入頁

**完成時間**：2026-01-31

---

### 指令 3：實作 Loading 狀態管理

**指令內容**：

```
請在 TimeLedger 前端實作 Loading 狀態管理系統：

目標：讓所有 API 呼叫都有適當的 loading 狀態

理由：沒有 loading 狀態，使用者不知道系統是否在運作

需求：

1. 建立 Loading 工具 `frontend/utils/loadingHelper.ts`：
   ```typescript
   // 自動管理 loading 狀態
   export async function withLoading<T>(
     loadingRef: Ref<boolean>,
     fn: () => Promise<T>
   ): Promise<T> {
     loadingRef.value = true
     try {
       return await fn()
     } finally {
       loadingRef.value = false
     }
   }

   // 平行請求的 loading 管理
   export function useParallelLoading() {
     const pending = ref(0)
     const isLoading = computed(() => pending.value > 0)

     const track = <T>(promise: Promise<T>): Promise<T> => {
       pending.value++
       return promise.finally(() => pending.value--)
     }

     return { isLoading, track }
   }

   // 組合式 loading（多個狀態）
   export function useCompoundLoading(
     ...loadingRefs: Ref<boolean>[]
   ) {
     const isLoading = computed(() =>
       loadingRefs.some(ref => ref.value)
     )
     return { isLoading }
   }
   ```

2. 建立 Loading 組件 `frontend/components/BaseLoading.vue`：
   ```typescript
   interface Props {
     loading: boolean
     size?: 'sm' | 'md' | 'lg'
     text?: string
   }

   // 功能：
   // - 支援不同大小
   // - 支援自訂文字
   // - 可作為遮罩或獨立組件
   ```

3. 建立 Loading Button 組件 `frontend/components/LoadingButton.vue`：
   ```typescript
   interface Props {
     loading: boolean
     disabled?: boolean
     // ... 其他按鈕屬性
   }

   // 功能：
   // - loading 時顯示載入中狀態
   // - loading 時禁用按鈕
   // - 支援各種按鈕樣式
   ```

4. 建立 Loading Skeleton 組件 `frontend/components/Skeleton.vue`：
   ```typescript
   interface Props {
     type: 'text' | 'avatar' | 'card' | 'table'
     lines?: number
   }

   // 支援常見的 skeleton 樣式
   ```

5. 規範 Store 中的 loading 命名：
   ```typescript
   // 通用 loading
   const loading = ref(false)

   // 特定操作 loading
   const isCreating = ref(false)
   const isUpdating = ref(false)
   const isDeleting = ref(false)
   const isFetching = ref(false)
   ```

6. 規範頁面中的 loading 使用：
   - API 呼叫開始時設定 loading = true
   - API 呼叫完成（無論成功失敗）設定 loading = false
   - 使用 finally 確保一定會清除

驗收標準」：
- [x] Loading 工具已建立 (`frontend/utils/loadingHelper.ts`)
- [x] BaseLoading 組件已實作 (`frontend/components/BaseLoading.vue`)
- [x] LoadingButton 組件已實作 (`frontend/components/LoadingButton.vue`)
- [x] Skeleton 組件已實作 (`frontend/components/Skeleton.vue`)
- [x] 所有 Store 都有 loading 狀態 (`teacher.ts`, `notification.ts`, `auth.ts`)
- [x] 頁面 UI 有適當的 loading 視覺回饋 (`teacher/dashboard.vue`, `teacher/profile.vue`)

完成後請同步更新 frontend_refactoring_commands
---

### 指令 4：建立標準 Store 模板

**指令內容**：

```
請在 TimeLedger 前端建立標準的 Store 模板：

目標：讓所有 Store 有一致的結構和模式

理由：沒有標準模板，會導致 Store 結構混亂，難以維護

需求：

1. 建立 Store 模板 `frontend/stores/template.ts`：
   ```typescript
   import { defineStore } from 'pinia'

   // 狀態類型
   interface State {
     // 資料
     items: Item[]
     // 分頁
     pagination: PaginationResult | null
     // 載入狀態
     loading: boolean
     isCreating: boolean
     isUpdating: boolean
     isDeleting: boolean
     // 錯誤
     error: string | null
   }

   // 初始狀態
   const initialState: State = {
     items: [],
     pagination: null,
     loading: false,
     isCreating: false,
     isUpdating: false,
     isDeleting: false,
     error: null
   }

   export const useTemplateStore = defineStore('template', () => {
     // 狀態
     const state = reactive({ ...initialState })

     // Getters
     const itemCount = computed(() => state.items.length)
     const hasItems = computed(() => state.items.length > 0)
     const isLoading = computed(() => state.loading)

     // Actions
     async function fetchItems(params?: PaginationParams) {
       state.loading = true
       state.error = null
       try {
         const api = useApi()
         const response = await api.get<ApiResponse<Item[]>>('/endpoint', { params })
         state.items = response.datas || []
         // 如果有分頁
         // state.pagination = response.pagination
       } catch (error) {
         state.error = getErrorMessage(error)
         throw error
       } finally {
         state.loading = false
       }
     }

     async function createItem(data: CreateRequest) {
       state.isCreating = true
       state.error = null
       try {
         const api = useApi()
         await api.post('/endpoint', data)
         await fetchItems() // 重新整理列表
       } catch (error) {
         state.error = getErrorMessage(error)
         throw error
       } finally {
         state.isCreating = false
       }
     }

     async function updateItem(id: number, data: UpdateRequest) {
       state.isUpdating = true
       state.error = null
       try {
         const api = useApi()
         await api.put(`/endpoint/${id}`, data)
         await fetchItems()
       } catch (error) {
         state.error = getErrorMessage(error)
         throw error
       } finally {
         state.isUpdating = false
       }
     }

     async function deleteItem(id: number) {
       state.isDeleting = true
       state.error = null
       try {
         const api = useApi()
         await api.delete(`/endpoint/${id}`)
         await fetchItems()
       } catch (error) {
         state.error = getErrorMessage(error)
         throw error
       } finally {
         state.isDeleting = false
       }
     }

     function clearError() {
       state.error = null
     }

     function reset() {
       Object.assign(state, initialState)
     }

     return {
       // 狀態（解構為 reactive）
       ...toRefs(state),
       // Getters
       itemCount,
       hasItems,
       isLoading,
       // Actions
       fetchItems,
       createItem,
       updateItem,
       deleteItem,
       clearError,
       reset
     }
   })
   ```

2. 建立 API 呼叫模板 `frontend/utils/apiHelper.ts`：
   ```typescript
   // 標準 API 呼叫範本
   export async function useApiCall<T>(
     apiFn: () => Promise<T>,
     options?: {
       onSuccess?: (data: T) => void
       onError?: (error: Error) => void
       loadingRef?: Ref<boolean>
     }
   ): Promise<T | null> {
     if (options?.loadingRef) {
       options.loadingRef.value = true
     }

     try {
       const data = await apiFn()
       options?.onSuccess?.(data)
       return data
     } catch (error) {
       options?.onError?.(error as Error)
       return null
     } finally {
       if (options?.loadingRef) {
         options.loadingRef.value = false
       }
     }
   }
   ```

3. 檢查並重構現有的 Store（teacher.ts, center.ts 等），讓它們符合模板

驗收標準：
- [x] Store 模板已建立（`frontend/stores/template.ts`）
- [x] API 呼叫範本已建立（`frontend/utils/apiHelper.ts`）
- [ ] 現有 Store 已重構為符合模板（待處理）
- [x] 新增 Store 時可直接套用模板

**完成日期**：2026-01-31

**建立檔案**：
- `frontend/stores/template.ts` - 標準 Store 模板
  - `createStoreTemplate()` 工廠函式
  - `createSimpleListStore()` 簡化列表 Store
  - 通用類型定義（BaseState、PaginationResult、ApiResponse）
- `frontend/utils/apiHelper.ts` - API 呼叫助手
  - `apiCall()` - 標準 API 呼叫範本
  - `withLoading()` - Loading 狀態包裝器
  - `createLoadingState()` - Loading 狀態管理
  - `createPaginationState()` - 分頁狀態管理
  - `useDebounce()` / `useThrottle()` - 防抖/節流
  - `getErrorMessage()` - 錯誤訊息處理

**使用範例**：
```typescript
// 使用標準模板建立 Store
import { createStoreTemplate } from '~/stores/template'

interface User {
  id: number
  name: string
  email: string
}

const useUserStore = createStoreTemplate<User, CreateUserDto, UpdateUserDto>({
  endpoint: '/users',
  getErrorMessage: (error) => getErrorMessage(error)
})
```

完成後請同步更新 frontend_refactoring_commands
---

## 第二部分：緊急修復

在基礎設施完成後，處理影響現有功能的緊急問題。

### 指令 5：實作 Template Cells Reorder API 後端路由

**指令內容**：

```
請在 TimeLedger 後端實作 Template Cells Reorder API：

目標：讓前端可以儲存範本儲存格的排序結果

需求：
1. 在 `app/servers/route.go` 中新增路由（約第 142 行後）：
   - HTTP 方法：PUT
   - 路由路徑：/api/v1/admin/templates/:templateId/cells/reorder
   - 中介層：authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()
   - 控制器：TimetableTemplateController.ReorderCells

2. 實作 TimetableTemplateController.ReorderCells 方法：
   - 接收參數：templateId (uint), cells (陣列，包含 id, sort_order)
   - 驗證 template 是否存在且屬於該中心
   - 批次更新所有 cells 的 sort_order
   - 回傳 200 狀態碼

3. 請在實作完成後使用 go build 驗證編譯是否成功

檔案位置：
- route.go: app/servers/route.go
- 控制器: app/controllers/timetable_template.go
- Repository: app/repositories/timetable_cell.go
```

驗收標準」：
- [x] 後端路由正確註冊
- [x] 控制器方法回傳 200 狀態碼
- [x] `go build ./...` 編譯成功
- [ ] 前端可成功呼叫並儲存排序結果

完成後請同步更新 frontend_refactoring_commands

---

## 變更摘要

**完成時間**：2026-01-31

**變更檔案**：
1. `app/models/timetable_cell.go` - 新增 `SortOrder` 欄位
2. `app/servers/route.go` - 新增 `/api/v1/admin/templates/:templateId/cells/reorder` PUT 路由
3. `app/controllers/timetable_template.go` - 新增 `ReorderCells` 方法
4. `app/repositories/timetable_cell.go` - 新增 `ReorderCellRequest` 結構和 `BatchUpdateSortOrder` 方法

**API 規格**：
- **端點**：`PUT /api/v1/admin/templates/:templateId/cells/reorder`
- **請求體**：
```json
{
  "cells": [
    {"id": 1, "sort_order": 1},
    {"id": 2, "sort_order": 2},
    {"id": 3, "sort_order": 3}
  ]
}
```
- **響應**：`200 OK` `{ "code": "SUCCESS", "data": null }`
- **認證**：需要 admin 權限（Bearer Token）

|**前端狀態**：✅ 已實作
- `frontend/pages/admin/templates.vue` 已有完整的拖曳排序功能
- `saveCellOrder()` 函數會在拖曳完成後自動呼叫 API
- 拖曳時有視覺回饋（透明度變化、邊框高亮）

---

## 驗證狀態

| 檢查項目 | 狀態 |
|---------|------|
| `go build ./...` 編譯成功 | ✅ |
| `go vet ./app/...` 無錯誤 | ✅ |
| 後端路由正確註冊 | ✅ |
| 控制器方法回傳 200 狀態碼 | ✅ |
| API 格式與前端相容 | ✅ |
| 前端拖曳功能完整 | ✅ |
| 前端可成功呼叫並儲存排序結果 | ✅ |

---

## 第三部分：核心功能補齊

基礎設施完成後，實作核心業務功能。

### 指令 6：實作 Teacher Invitations Store 方法

**指令內容**：

```
請在 TimeLedger 前端實作 Teacher Invitations 功能：

目標：讓老師可以查看、接受/拒絕中心的邀請

需求：
1. 在 `frontend/stores/teacher.ts` 中新增以下狀態和方法：

   狀態：
   - invitations: Ref<Invitation[]>
   - pendingInvitationsCount: Ref<number>

   方法：
   - fetchInvitations(): 取得邀請列表
     * 呼叫 GET /api/v1/teacher/me/invitations
     * 將結果存入 invitations

   - respondToInvitation(invitationId: number, action: 'ACCEPT' | 'REJECT'): 回應邀請
     * 呼叫 POST /api/v1/teacher/me/invitations/respond
     * 回應成功後重新 fetchInvitations

   - fetchPendingCount(): 取得待處理邀請數量
     * 呼叫 GET /api/v1/teacher/me/invitations/pending-count
     * 將結果存入 pendingInvitationsCount

2. 定義 Invitation 型別（至少包含 id, center_name, status, created_at）

3. 請在頁面 `frontend/pages/teacher/invitations.vue` 中整合這些方法

檔案位置：
- Store: frontend/stores/teacher.ts
- 頁面: frontend/pages/teacher/invitations.vue
```

驗收標準」：
- [x] 邀請列表正確顯示
- [x] 接受/拒絕功能正常運作
- [x] 待處理數量即時更新
- [x] 拒絕邀請後邀請從列表移除

**完成時間**：2026-01-31

完成後請同步更新 frontend_refactoring_commands
---

### 指令 7：確認 Hashtags 功能需求

**指令內容**：

```
請確認 TimeLedger 的 Hashtags 功能需求：

目標：決定是否需要在前端實作標籤搜尋與建立功能

需求：
1. 檢查後端 `/api/v1/hashtags/search` 和 `/api/v1/hashtags` 的使用場景

2. 確認以下問題並回報：
   - Hashtags 功能是否為預留功能？
   - 前端哪些頁面可能需要標籤搜尋？（教師技能、個人品牌等）
   - 是否需要在前端實作標籤搜尋與建立？

3. 如果不需要前端實作：
   - 建議移除後端路由定義
   - 列出需要清理的相關程式碼

4. 如果需要實作：
   - 在 `frontend/stores/teacher.ts` 新增：
     * searchHashtags(query: string): 搜尋標籤
     * createHashtag(tag: string): 建立標籤
```

驗收標準」：
- [ ] 已確認 Hashtags 功能需求
- [ ] 已決定是否需要前端實作
- [ ] 如需實作，已提供 Store 方法範例
- [ ] 如不需實作，已提供清理建議

完成後請同步更新 frontend_refactoring_commands
---

### 指令 8：實作 Smart Matching Store

**指令內容**：

```
請在 TimeLedger 前端建立 Smart Matching Store：

目標：為智慧媒合功能建立統一的狀態管理

需求：
1. 建立 `frontend/stores/smartMatching.ts`

2. 定義狀態：
   - searchResults: Ref<MatchingResult[]>
   - talentStats: Ref<TalentStats | null>
   - loading: Ref<boolean>

3. 實作方法：
   - searchMatches(params: MatchingParams): 智慧媒合搜尋
     * 呼叫 POST /admin/smart-matching/matches

   - fetchTalentStats(): 取得人才庫統計
     * 呼叫 GET /admin/smart-matching/talent/stats

   - searchTalent(params: TalentSearchParams): 人才庫搜尋
     * 呼叫 GET /admin/smart-matching/talent/search

   - inviteTalent(teacherIds: number[], message?: string): 邀請人才
     * 呼叫 POST /admin/smart-matching/talent/invite

4. 定義必要的 TypeScript 型別（MatchingResult, TalentStats, MatchingParams 等）

5. 整合到頁面 `frontend/pages/admin/matching.vue`

檔案位置：
- Store: frontend/stores/smartMatching.ts
- 頁面: frontend/pages/admin/matching.vue
```

驗收標準」：
- [ ] Smart Matching Store 已建立
- [ ] 智慧媒合可找到替代教師
- [ ] 人才庫統計正確顯示
- [ ] 可搜尋與篩選人才
- [ ] 可發送人才邀請

完成後請同步更新 frontend_refactoring_commands

**完成日期**：2026-01-31

**完成的實作內容**：

1. **已建立 `frontend/stores/smartMatching.ts`**
   - 使用 Pinia setup syntax
   - 整合 `withLoading` helper 處理 loading 狀態
   - 完整 TypeScript 類型定義

2. **定義的狀態**：
   - 搜尋相關：`searchResults`、`hasSearched`、`searchProgress`、`searchSteps`
   - 人才庫相關：`talentResults`、`talentStats`、`cityDistribution`、`topSkills`
   - 教師課表：`selectedTeacher`、`teacherSessions`、`alternativeSlots`
   - 比較模式：`selectedForCompare`、`viewMode`
   - 邀請相關：`invitationStatuses`、`inviteLoadingIds`、`bulkLoading`
   - Loading 狀態：`isSearching`、`isSearchingTalent`、`isFetchingStats` 等

3. **實作的方法**：
   - `searchMatches()` - 智慧媒合搜尋，呼叫 POST /admin/smart-matching/matches
   - `fetchTalentStats()` - 取得人才庫統計，呼叫 GET /admin/smart-matching/talent/stats
   - `searchTalent()` - 人才庫搜尋，呼叫 GET /admin/smart-matching/talent/search
   - `inviteTalent()` - 邀請人才，呼叫 POST /admin/smart-matching/talent/invite
   - `bulkInviteTalents()` - 批量邀請人才
   - `fetchTeacherSchedule()` - 取得教師課表
   - `fetchAlternativeSlots()` - 取得替代時段建議
   - 比較功能：`toggleCompare`、`removeFromCompare`、`exitCompareMode`

4. **已整合到頁面 `frontend/pages/admin/matching.vue`**
   - 所有 API 調用改為使用 Store 方法
   - 所有狀態改為使用 Store 狀態
   - 保持原有的 UI 體驗和功能完整性

**驗收標準」**：
- [x] Smart Matching Store 已建立
- [x] 智慧媒合可找到替代教師
- [x] 人才庫統計正確顯示
- [x] 可搜尋與篩選人才
- [x] 可發送人才邀請

### 指令 9：整合 LINE Binding 功能

**指令內容**：

```
請在 TimeLedger 前端整合 LINE Binding 功能：

目標：讓管理員可以在設定頁面綁定 LINE 帳號

需求：
1. 在 `frontend/pages/admin/settings.vue` 中新增 LINE 綁定區塊

2. 實作以下功能：
   - fetchLineBindingStatus(): 取得綁定狀態
     * 呼叫 GET /api/v1/admin/me/line-binding
     * 顯示已綁定/未綁定狀態
     * 顯示 LINE ID 和綁定時間

   - fetchLineNotifySettings(): 取得通知設定
     * 呼叫 GET /api/v1/admin/me/line/notify-settings
     * 取得通知開關狀態

   - updateLineNotifySettings(settings): 更新通知設定
     * 呼叫 PATCH /api/v1/admin/me/line/notify-settings
     * 支援獨立的例外通知和審核結果通知

   - toggleLineNotifySetting(): 切換通知開關
     * 提供開關 UI 並即時更新設定

3. UI 需求：
   - 已綁定：顯示 LINE ID + 綁定時間 + 管理設定連結
   - 未綁定：顯示「立即綁定」按鈕
   - 通知開關：可選擇性關閉特定通知類型

檔案位置：
- 頁面: frontend/pages/admin/settings.vue
- 完整綁定頁面: frontend/pages/admin/line-bind.vue (已存在)
```

**驗收標準**：
- [x] 綁定狀態正確顯示 (LINE ID, 綁定時間)
- [x] 可產生驗證碼並完成綁定 (line-bind.vue)
- [x] 通知設定可儲存 (獨立的例外/審核通知開關)
- [x] 設定頁面整合完成

**完成狀態**：✅ 已完成 (2026-01-31)

---
---

### 指令 10：實作 Centers CRUD 功能

**指令內容**：

```
請在 TimeLedger 前端實作 Centers CRUD 功能：

目標：讓管理員可以完整管理中心的編輯與刪除

需求：
1. 檢查現有 `frontend/pages/admin/centers.vue` 功能完整性

2. 新增以下尚未實作的功能：

   GET /:id - 取得單一中心詳細資料
   - 在編輯頁面載入中心資料
   - 顯示中心名稱、地址、聯絡資訊等

   PUT /:id - 更新中心資料
   - 實作編輯表單
   - 包含名稱、地址、電話、營業時間等欄位
   - 提交前驗證必填欄位

   DELETE /:id - 刪除中心
   - 彈出確認對話框（"確定要刪除此中心嗎？"）
   - 刪除後重新導向列表頁
   - 處理有課程/老師關聯的狀況（禁止刪除或提示）

3. 定義 Center 與 CenterForm 型別

4. 在 Store 中新增：
   - fetchCenter(id: number): 取得單一中心
   - updateCenter(id: number, data: CenterForm): 更新中心
   - deleteCenter(id: number): 刪除中心

檔案位置：
- 頁面: frontend/pages/admin/centers.vue
- 頁面: frontend/pages/admin/centers/[id].vue（新增）
- Store: frontend/stores/center.ts（新增或修改）
```

驗收標準」：
- [ ] 中心列表可正常顯示
- [ ] 可進入編輯頁面並載入資料
- [ ] 可更新中心資料
- [ ] 可刪除中心（如適用）
- [ ] 欄位驗證正確

---

### 指令 11：補齊 Rooms/Courses/Offerings 功能

**指令內容**：

```
請檢查並補齊 TimeLedger 的 Rooms/Courses/Offerings 功能：

目標：確保所有後端 API 都有對應的前端功能

需求：
1. 檢查 `frontend/pages/admin/rooms.vue`：
   - [x] 教室列表顯示完整
   - [x] 新增/編輯/刪除功能正常
   - [x] 啟用/停用切換功能正常

2. 檢查 `frontend/pages/admin/courses.vue`：
   - [x] 課程列表顯示完整
   - [x] 新增/編輯功能正常
   - [x] 啟用/停用切換功能正常

3. 檢查 `frontend/pages/admin/offerings.vue`：
   - [x] 開課列表顯示完整
   - [x] 新增開課功能正常
   - [x] 複製開課功能（已實作）
   - [x] 啟用/停用切換功能正常

4. 針對缺失的功能：
   - [x] 補齊 Toggle 按鈕與功能
   - [x] 補齊複製開課功能

5. 完成項目：
   - RoomsTab.vue：新增 toggleRoom 函數與 UI
   - CoursesTab.vue：新增 toggleCourse 函數與 UI
   - OfferingsTab.vue：新增 toggleOffering 與 copyOffering 函數與 UI
```

**驗收標準」**：
- [x] Rooms 所有功能正常運作
- [x] Courses 所有功能正常運作
- [x] Offerings 所有功能正常運作
- [x] 複製開課功能已實作
- [x] 啟用/停用切換正常

完成後請同步更新 frontend_refactoring_commands

---

## 第五部分：排課驗證功能

### 指令 12：整合排課驗證 APIs

**指令內容**：

```
請在 TimeLedger 前端整合排課驗證 API：

目標：在排課表單中即時驗證重疊與緩衝時間

需求：
1. 檢查 `frontend/pages/admin/schedules.vue` 或相關排課頁面

2. 新增以下驗證功能：

   Check Overlap - 檢查時段重疊
   - POST /admin/scheduling/check-overlap
   - 驗證新時段是否與現有課程重疊
   - 重疊時顯示錯誤訊息

   Check Teacher Buffer - 檢查老師緩衝時間
   - POST /admin/scheduling/check-teacher-buffer
   - 驗證老師上一堂課與新時段的間隔
   - 不足時顯示警告或錯誤

   Check Room Buffer - 檢查教室緩衝時間
   - POST /admin/scheduling/check-room-buffer
   - 驗證教室使用間隔
   - 不足時顯示警告或錯誤

   Full Validation - 完整驗證
   - POST /admin/scheduling/validate
   - 一次檢查所有規則
   - 回傳完整驗證結果

3. 實作時機：
   - 表單提交前自動驗證
   - 或使用者點擊「驗證」按鈕

4. 錯誤顯示：
   - Hard Conflict：直接顯示錯誤，阻止提交
   - Buffer Conflict：顯示警告，支援 override 參數

5. 定義 ValidationResult 與 ValidationConflict 型別

檔案位置：
- 頁面: frontend/pages/admin/schedules.vue（或相關檔案）
- Composable: frontend/composables/useSchedulingValidation.ts（已新增）

驗收標準」：
- [x] 排課表單有即時驗證功能
- [x] 可檢查時段重疊
- [x] 可檢查緩衝時間
- [x] 錯誤與警告正確顯示
- [x] 支援 Buffer Override 參數

完成後請同步更新 frontend_refactoring_commands.md

## 完成記錄

**完成日期**：2026-01-31

**實作內容**：
1. 建立 `frontend/composables/useSchedulingValidation.ts`
   - checkOverlap(): 檢查時段重疊
   - checkTeacherBuffer(): 檢查教師緩衝時間
   - checkRoomBuffer(): 檢查教室緩衝時間
   - validateSchedule(): 完整驗證
   - quickValidate(): 快速驗證（組合多個檢查）
   - formatConflictMessages(): 格式化衝突訊息
   - hasHardConflicts(): 判斷是否為硬衝突
   - canOverrideAll(): 判斷是否可覆寫

2. 更新 `frontend/components/ScheduleRuleModal.vue`
   - 新增 validationLoading 狀態
   - 在表單提交前進行驗證
   - 硬衝突時阻止提交並顯示錯誤
   - 緩衝衝突時顯示警告並詢問是否覆寫
   - 更新提交按鈕顯示驗證中狀態

3. 新增類型定義
   - ValidationConflictType
   - ValidationConflictDetail
   - BaseValidationRequest
   - CheckOverlapRequest
   - CheckTeacherBufferRequest
   - CheckRoomBufferRequest
   - FullValidationRequest
   - ConflictCheckResponse
   - BufferCheckResponse
   - ScheduleValidationResponse
---

## 執行順序建議

| 順序 | 指令 | 說明 | 預估時間 |
|:---:|:---|:---|:---:|
| 1 | 指令 1 | API Response 型別統一 | 4-6 小時 |
| 2 | 指令 2 | 錯誤處理機制統一 | 2-4 小時 |
| 3 | 指令 3 | Loading 狀態管理 | 2-3 小時 |
| 4 | 指令 4 | 標準 Store 模板 | 2-4 小時 |
| 5 | 指令 5 | Template Cells Reorder API（緊急） | 2 小時 |
| 6 | 指令 6 | Teacher Invitations | 4 小時 |
| 7 | 指令 7 | Hashtags 功能確認 | 2-4 小時 |
| 8 | 指令 8 | Smart Matching Store | 4 小時 |
| 9 | 指令 9 | LINE Binding 整合 | 4 小時 |
| 10 | 指令 10 | Centers CRUD | 8 小時 |
| 11 | 指令 11 | Rooms/Courses/Offerings 補齊 | 8 小時 |
| 12 | 指令 12 | Scheduling Validation | 6 小時 |

**總預估時間：52-61 小時**

---

## 給 AI 助手的範本

```
我需要執行以下重構任務：

【任務名稱】
[貼上任務名稱]

【指令內容】
[貼上對應的指令內容]

【當前狀態】
[描述目前的情況]
例如：「前端 templates.vue 第 775 行呼叫了不存在的 API」

【已完成的基礎設施】
[標記已完成的部分]
例如：「已完成 API Response 型別定義、錯誤處理機制」

【預期結果】
[描述你希望達成的目標]

請根據以上需求進行重構，完成後請：
1. 列出修改的檔案
2. 說明修改內容
3. 驗證編譯/執行是否正常
```

---

## 為什麼要這樣排序？

| 順序 | 原因 |
|:---|:---|
| 1-4 | **基礎設施優先**：型別、錯誤處理、Loading、Store 模板是所有功能的基礎 |
| 5 | **緊急修復**：影響現有功能，必須盡快處理 |
| 6-8 | **核心功能**：教師邀請、Hashtags、Smart Matching 是主要業務功能 |
| 9-11 | **管理功能**：LINE 綁定、中心管理是輔助功能 |
| 12 | **排課驗證**：功能複雜，需要在基礎設施完成後實作 |

**好處**：
- 先完成基礎設施，後續功能開發速度更快
- 所有功能都有一致的程式碼風格
- 維護成本降低
- 新人接手更容易

---

## 更新記錄

| 日期 | 版本 | 變更內容 |
|:---|:---:|:---|
| 2026-01-31 | 1.0 | 初版（按功能分類） |
| 2026-01-31 | 2.0 | 重新排序，改為「技術優先」 |

