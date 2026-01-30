# TimeLedger 前端實作任務清單

> 生成日期：2026-01-31
> 資料來源：API 對齊分析報告
> 總 API 數量：147 個後端端點
> 前端已實作：約 85 個（57.8%）
> 待實作數量：62 個（42.2%）

---

## 總覽

本文件追蹤 TimeLedger 前端實作進度，確保所有後端 API 都有對應的前端整合。實作優先順序依據功能緊急性與業務價值排列，目標是在三個月內達到 85% API 覆蓋率（125 個 API）。

---

## 第一部分：緊急修復（高優先級）

此部分包含影響現有功能運作的問題，需要立即處理。

### 1.1 Template Cells Reorder API 後端實作

**狀態**：待處理  
**優先順序**：🔴 最高  
**影響範圍**：範本管理功能

**問題描述**：前端 `frontend/pages/admin/templates.vue` 第 775 行嘗試呼叫不存在的後端路由，導致範本儲存格排序功能無法使用。

**前端呼叫程式碼**：

```typescript
await api.put(`/admin/templates/${selectedTemplate.value.id}/cells/reorder`, {
  cells: orderData
})
```

**後端實作需求**：

| 項目 | 內容 |
|:---|:---|
| HTTP 方法 | PUT |
| 路由路徑 | `/api/v1/admin/templates/:templateId/cells/reorder` |
| 檔案位置 | `app/servers/route.go`（約第 142 行後） |
| 中介層 | `authMiddleware.Authenticate()`, `authMiddleware.RequireCenterAdmin()` |
| 控制器方法 | `TimetableTemplateController.ReorderCells` |

**新增路由定義**（route.go）：

```go
{http.MethodPut, "/api/v1/admin/templates/:templateId/cells/reorder", s.action.timetableTemplate.ReorderCells, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
```

**驗證標準**：
- [ ] 後端路由正確註冊
- [ ] 控制器方法回傳 200 狀態碼
- [ ] 前端可成功儲存排序結果
- [ ] 重新整理頁面後排序保持不變

---

## 第二部分：核心功能補齊（中優先級）

此部分包含業務核心流程所需的功能，需要在兩週內完成。

### 2.1 Teacher Invitations 功能完整實作

**狀態**：未實作  
**優先順序**：🔴 高  
**關聯頁面**：`frontend/pages/teacher/invitations.vue`

**後端 API 清單**：

| 方法 | 端點 | 功能說明 |
|:---:|:---|:---|
| GET | `/api/v1/teacher/me/invitations` | 取得邀請列表 |
| POST | `/api/v1/teacher/me/invitations/respond` | 回應邀請（接受/拒絕） |
| GET | `/api/v1/teacher/me/invitations/pending-count` | 待處理邀請數量 |

**前端 Store 實作需求**（`frontend/stores/teacher.ts`）：

```typescript
// 邀請相關狀態
const invitations = ref<Invitation[]>([])
const pendingInvitationsCount = ref(0)

// 取得邀請列表
const fetchInvitations = async () => {
  const api = useApi()
  const response = await api.get<ApiResponse<Invitation[]>>('/teacher/me/invitations')
  invitations.value = response.datas
}

// 回應邀請
const respondToInvitation = async (invitationId: number, action: 'ACCEPT' | 'REJECT') => {
  const api = useApi()
  await api.post('/teacher/me/invitations/respond', { invitation_id: invitationId, action })
  await fetchInvitations()
}

// 取得待處理數量
const fetchPendingCount = async () => {
  const api = useApi()
  const response = await api.get<ApiResponse<{ count: number }>>('/teacher/me/invitations/pending-count')
  pendingInvitationsCount.value = response.datas.count
}
```

**驗證標準**：
- [ ] 邀請列表正確顯示
- [ ] 接受/拒絕功能正常運作
- [ ] 待處理數量即時更新
- [ ] 拒絕邀請後邀請從列表移除

---

### 2.2 Hashtags 功能確認與實作

**狀態**：待確認  
**優先順序**：🟡 中  
**後端 API 清單**：

| 方法 | 端點 | 功能說明 |
|:---:|:---|:---|
| GET | `/api/v1/hashtags/search` | 搜尋標籤 |
| POST | `/api/v1/hashtags` | 建立標籤 |

**待確認事項**：

- [ ] Hashtags 功能是否為預留功能
- [ ] 是否需要在前端實作標籤搜尋與建立
- [ ] 標籤使用場景（教師技能、個人品牌等）

**若需要實作**（`frontend/stores/teacher.ts`）：

```typescript
// 搜尋標籤
const searchHashtags = async (query: string) => {
  const api = useApi()
  const response = await api.get<ApiResponse<Hashtag[]>>(`/hashtags/search?q=${query}`)
  return response.datas
}

// 建立標籤
const createHashtag = async (tag: string) => {
  const api = useApi()
  const response = await api.post<ApiResponse<Hashtag>>('/hashtags', { tag })
  return response.datas
}
```

**若不需要**：

- [ ] 移除後端路由定義
- [ ] 清理相關程式碼

---

### 2.3 Smart Matching 前端整合

**狀態**：未實作  
**優先順序**：🟡 中  
**關聯頁面**：`frontend/pages/admin/matching.vue`

**後端 API 清單**：

| 方法 | 端點 | 功能說明 |
|:---:|:---|:---|
| POST | `/admin/smart-matching/matches` | 智慧媒合搜尋 |
| GET | `/admin/smart-matching/suggestions` | 搜尋建議 |
| POST | `/admin/smart-matching/alternatives` | 替代時段建議 |
| GET | `/admin/teachers/:teacher_id/sessions` | 教師課表查詢 |
| GET | `/admin/smart-matching/talent/search` | 人才庫搜尋 |
| GET | `/admin/smart-matching/talent/stats` | 人才庫統計 |
| POST | `/admin/smart-matching/talent/invite` | 邀請人才 |

**實作規劃**：

**第一步：建立 Smart Matching Store**（`frontend/stores/smartMatching.ts`）

```typescript
export const useSmartMatchingStore = defineStore('smartMatching', () => {
  const searchResults = ref<MatchingResult[]>([])
  const talentStats = ref<TalentStats | null>(null)
  const loading = ref(false)

  // 智慧媒合搜尋
  const searchMatches = async (params: MatchingParams) => {
    const api = useApi()
    loading.value = true
    try {
      const response = await api.post<ApiResponse<MatchingResult[]>>('/admin/smart-matching/matches', params)
      searchResults.value = response.datas
    } finally {
      loading.value = false
    }
  }

  // 人才庫統計
  const fetchTalentStats = async () => {
    const api = useApi()
    const response = await api.get<ApiResponse<TalentStats>>('/admin/smart-matching/talent/stats')
    talentStats.value = response.datas
  }

  // 人才庫搜尋
  const searchTalent = async (params: TalentSearchParams) => {
    const api = useApi()
    const response = await api.get<ApiResponse<TalentSearchResult>>('/admin/smart-matching/talent/search', { params })
    return response.datas
  }

  // 邀請人才
  const inviteTalent = async (teacherIds: number[], message?: string) => {
    const api = useApi()
    return await api.post('/admin/smart-matching/talent/invite', { teacher_ids: teacherIds, message })
  }

  return {
    searchResults,
    talentStats,
    loading,
    searchMatches,
    fetchTalentStats,
    searchTalent,
    inviteTalent
  }
})
```

**第二步：整合到管理員介面**

- [ ] 確認 `frontend/pages/admin/matching.vue` 頁面存在
- [ ] 整合搜尋功能與結果顯示
- [ ] 實作人才庫統計卡片
- [ ] 實作人才邀請功能

**驗證標準**：
- [ ] 智慧媒合可找到替代教師
- [ ] 人才庫統計正確顯示
- [ ] 可搜尋與篩選人才
- [ ] 可發送人才邀請

---

## 第三部分：管理功能實作（中優先級）

此部分包含管理後台的各項管理功能。

### 3.1 Admin User Management 管理功能（私下 API 使用，無需前端串接）

**狀態**：後端 API 已完成，供私下打 API 使用  
**優先順序**：⚪ 無需前端實作  
**說明**：此功能不需要前端介面，管理員可透過 Postman 或其他工具直接呼叫 API

**後端 API 清單**：

| 方法 | 端點 | 功能說明 |
|:---:|:---|:---|
| GET | `/api/v1/admin/me/profile` | 取得管理員個人資料 |
| POST | `/api/v1/admin/me/change-password` | 修改密碼 |
| GET | `/api/v1/admin/admins` | 管理員列表 |
| POST | `/api/v1/admin/admins` | 新增管理員 |
| POST | `/api/v1/admin/admins/toggle-status` | 切換狀態 |
| POST | `/api/v1/admin/admins/reset-password` | 重設密碼 |
| POST | `/api/v1/admin/admins/change-role` | 變更角色 |

**待辦事項**：
- [x] ~~建立 Admin User Store~~
- [x] ~~建立管理員設定頁面~~

**驗證標準**：不適用（無需前端實作）

---

### 3.2 LINE Integration 整合

**狀態**：未實作  
**優先順序**：🟡 中  
**後端 API 清單**：

| 方法 | 端點 | 功能說明 |
|:---:|:---|:---|
| GET | `/api/v1/admin/me/line-binding` | LINE 綁定狀態 |
| POST | `/api/v1/admin/me/line/bind` | 初始化綁定 |
| DELETE | `/api/v1/admin/me/line/unbind` | 解除綁定 |
| PATCH | `/api/v1/admin/me/line/notify-settings` | 更新通知設定 |

**實作規劃**：

**整合到個人設定頁面**（`frontend/pages/admin/settings.vue`）

```typescript
// LINE 綁定狀態
const lineBindingStatus = ref<LineBindingStatus | null>(null)

// 取得綁定狀態
const fetchLineBindingStatus = async () => {
  const api = useApi()
  const response = await api.get<ApiResponse<LineBindingStatus>>('/admin/me/line-binding')
  lineBindingStatus.value = response.datas
}

// 初始化綁定（產生驗證碼）
const initiateLineBind = async () => {
  const api = useApi()
  const response = await api.post<ApiResponse<{ code: string; qrCode: string }>>('/admin/me/line/bind')
  return response.datas
}

// 解除綁定
const unbindLine = async () => {
  const api = useApi()
  await api.delete('/admin/me/line/unbind')
  await fetchLineBindingStatus()
}

// 更新通知設定
const updateLineNotifySettings = async (settings: LineNotifySettings) => {
  const api = useApi()
  await api.patch('/admin/me/line/notify-settings', settings)
}
```

**UI 元件需求**：

- [ ] LINE 綁定狀態顯示
- [ ] 開始綁定按鈕（產生驗證碼與 QR Code）
- [ ] 解除綁定按鈕（含確認對話框）
- [ ] 通知開關（可選擇性關閉特定通知）

**驗證標準**：
- [ ] 綁定狀態正確顯示
- [ ] 可產生驗證碼並完成綁定
- [ ] 可解除綁定
- [ ] 通知設定可儲存

---

### 3.3 Centers CRUD 功能完整實作

**狀態**：部分實作  
**優先順序**：🟡 中  
**後端 API 清單**：

| 方法 | 端點 | 狀態 |
|:---:|:---|::|
| GET | `/api/v1/admin/centers` | ✅ 已使用 |
| POST | `/api/v1/admin/centers` | ✅ 已使用 |
| GET | `/api/v1/admin/centers/:id` | ❌ 未使用 |
| PUT | `/api/v1/admin/centers/:id` | ❌ 未使用 |
| DELETE | `/api/v1/admin/centers/:id` | ❌ 未使用 |

**待辦事項**：

- [ ] 確認是否需要中心編輯功能
- [ ] 確認是否需要中心刪除功能
- [ ] 如需實作，建立 `frontend/pages/admin/centers/[id].vue` 編輯頁面
- [ ] 實作刪除確認對話框與邏輯

---

### 3.4 Rooms 功能完整實作

**狀態**：部分實作  
**優先順序**：🟡 中  
**後端 API 清單**：5 個（全部已定義）

**待辦事項**：

- [ ] 檢查教室管理頁面是否完整
- [ ] 確認啟用/停用功能是否正常運作
- [ ] 修復任何缺失的功能

---

### 3.5 Courses 功能完整實作

**狀態**：部分實作  
**優先順序**：🟡 中  
**後端 API 清單**：6 個

**待辦事項**：

- [ ] 檢查課程管理頁面功能完整性
- [ ] 確認課程切換啟用功能是否正常
- [ ] 修復任何缺失的功能

---

### 3.6 Offerings 功能完整實作

**狀態**：部分實作  
**優先順序**：🟡 中  
**後端 API 清單**：6 個

**待辦事項**：

- [ ] 實作開課複製功能
- [ ] 實作啟用/停用切換
- [ ] 整合到開課管理頁面

---

### 3.7 Templates 功能完整實作

**狀態**：部分實作  
**優先順序**：🟡 中  
**後端 API 清單**：多個

**待辦事項**：

- [ ] 實作範本編輯功能
- [ ] 實作範本更新功能
- [ ] 實作套用前驗證功能

---

## 第四部分：排課驗證與例外管理

### 4.1 Scheduling Validation APIs 整合

**狀態**：部分使用  
**優先順序**：🟡 中  
**後端 API 清單**：

| 方法 | 端點 | 功能說明 |
|:---:|:---|:---|
| POST | `/admin/scheduling/check-overlap` | 檢查重疊 |
| POST | `/admin/scheduling/check-teacher-buffer` | 檢查老師緩衝 |
| POST | `/admin/scheduling/check-room-buffer` | 檢查教室緩衝 |
| POST | `/admin/scheduling/validate` | 完整驗證 |

**待辦事項**：

- [ ] 確認排課表單是否有即時驗證
- [ ] 整合驗證 API 到排課表單
- [ ] 顯示驗證錯誤與警告資訊
- [ ] 支援 Buffer Override 參數

---

### 4.2 Exceptions Management 功能確認

**狀態**：部分實作  
**優先順序**：🟡 中  
**後端 API 清單**：6 個

**待辦事項**：

- [ ] 檢查例外管理頁面功能完整性
- [ ] 確認日期範圍查詢功能
- [ ] 確認所有例外列表篩選功能
- [ ] 確認待審核列表功能

---

## 第五部分：擴充功能（低優先級）

此部分包含較不緊急的功能，可在主要功能完成後實作。

### 5.1 Calendar & Export 功能實作

**狀態**：未實作  
**優先順序**：🟢 低  
**後端 API 清單**：

| 方法 | 端點 | 功能說明 |
|:---:|:---|:---|
| POST | `/api/v1/teacher/me/schedule/subscription` | 建立日曆訂閱 |
| DELETE | `/api/v1/teacher/me/schedule/subscription` | 取消訂閱 |
| POST | `/api/v1/teacher/me/backgrounds` | 上傳背景圖 |
| GET | `/api/v1/teacher/me/backgrounds` | 取得背景圖列表 |
| DELETE | `/api/v1/teacher/me/backgrounds` | 刪除背景圖 |

**實作規劃**：

**第一步：課表匯出功能頁面**

- [ ] 實作 CSV 匯出按鈕與功能
- [ ] 實作 PDF 匯出功能
- [ ] 實作 ICS 日曆匯出功能
- [ ] 實作圖片匯出功能

**第二步：日曆訂閱功能**

- [ ] 實作訂閱 URL 產生與複製
- [ ] 實作取消訂閱功能
- [ ] 在課表頁面顯示訂閱選項

**第三步：背景圖管理**

- [ ] 實作背景圖上傳介面
- [ ] 實作背景圖列表顯示
- [ ] 實作背景圖刪除功能
- [ ] 整合到課表個人化設定

---

### 5.2 Public Invitations 頁面實作

**狀態**：未實作  
**優先順序**：🟢 低  
**後端 API 清單**：

| 方法 | 端點 | 功能說明 |
|:---:|:---|:---|
| GET | `/api/v1/invitations/:token` | 取得公開邀請資訊 |
| POST | `/api/v1/invitations/:token/accept` | 接受邀請 |

**待辦事項**：

- [ ] 確認邀請頁面是否存在
- [ ] 建立公開邀請接受頁面（`frontend/pages/invitations/[token].vue`）
- [ ] 實作驗證碼輸入與接受流程
- [ ] 整合 LINE 登入（如需要）

---

### 5.3 Geo API 整合

**狀態**：未使用  
**優先順序**：🟢 低  
**後端 API 清單**：

| 方法 | 端點 | 功能說明 |
|:---:|:---|:---|
| GET | `/api/v1/geo/cities` | 取得城市列表 |

**待辦事項**：

- [ ] 確認是否用於下拉選單
- [ ] 如需要，在相關表單中整合城市選擇器

---

### 5.4 Notification Queue 監控整合

**狀態**：未實作  
**優先順序**：🟢 低  
**後端 API 清單**：

| 方法 | 端點 | 功能說明 |
|:---:|:---|:---|
| GET | `/api/v1/admin/notifications/queue-stats` | 通知佇列統計 |

**實作規劃**：

- [ ] 整合到管理員儀表板或系統監控頁面
- [ ] 實作佇列統計卡片
- [ ] 實作失敗率警示
- [ ] 實作自動重新整理（每 30 秒）

---

### 5.5 Scheduling Expansion APIs 整合

**狀態**：未使用  
**優先順序**：🟢 低  
**後端 API 清單**：

| 方法 | 端點 | 功能說明 |
|:---:|:---|:---|
| POST | `/api/v1/admin/expand-rules` | 展開規則 |
| POST | `/api/v1/admin/detect-phase-transitions` | 檢測階段轉換 |

**待辦事項**：

- [ ] 確認是否用於課表預覽
- [ ] 整合到排課管理頁面
- [ ] 實作階段轉換偵測顯示

---

## 第六部分：技術債務

### 6.1 Type Definitions 統一管理

**狀態**：需要改善  
**優先順序**：🟡 中  
**待辦事項**：

- [ ] 建立 `frontend/types/api.ts` 統一管理所有 API Response 型別
- [ ] 為所有 API Response 定義 TypeScript 介面
- [ ] 移除 any 型別，改用明確的介面
- [ ] 確保前後端型別一致

**建議的型別檔案結構**：

```
frontend/types/
├── api.ts              // API 通用回應格式
├── admin.ts            // 管理員相關型別
├── teacher.ts          // 教師相關型別
├── center.ts           // 中心相關型別
├── scheduling.ts       // 排課相關型別
├── matching.ts         // 智慧媒合相關型別
└── notification.ts     // 通知相關型別
```

---

### 6.2 Error Handling 統一處理

**狀態**：需要改善  
**優先順序**：🟡 中  
**待辦事項**：

- [ ] 檢查所有 API 呼叫的錯誤處理
- [ ] 實作統一的錯誤提示機制
- [ ] 建立錯誤攔截器（Axios Interceptor）
- [ ] 根據錯誤碼顯示對應的錯誤訊息

---

### 6.3 Loading States 狀態管理

**狀態**：需要改善  
**優先順序**：🟡 中  
**待辦事項**：

- [ ] 在 Stores 中加入 loading 狀態
- [ ] 實作統一的 Loading 組件
- [ ] 在 API 呼叫時顯示 loading 狀態
- [ ] 確保 loading 狀態正確清除

---

## 實作優先順序總結

### Phase 1：緊急修復（已完成/進行中）

| 工作項目 | 狀態 | 預估時間 |
|:---|:---:|:---:|
| Personal Event PATCH 方法修復 | ✅ 已完成 | - |
| Template Cells Reorder API 後端實作 | ⚠️ 待處理 | 2 小時 |

### Phase 2：核心功能補齊（兩週內）

| 工作項目 | 優先順序 | 預估時間 |
|:---|:---:|:---:|
| Teacher Invitations 功能完整實作 | 🔴 高 | 4 小時 |
| Hashtags 功能確認與實作 | 🟡 中 | 2-4 小時 |
| Smart Matching Store 建立 | 🟡 中 | 4 小時 |
| Smart Matching 頁面整合 | 🟡 中 | 6 小時 |

### Phase 3：管理功能實作（一個月內）

| 工作項目 | 優先順序 | 預估時間 | 狀態 |
|:---|:---:|:---:|:---:|
| ~~Admin User Store 建立~~ | 🟡 中 | 3 小時 | ⚪ 無需實作 |
| ~~Admin Settings 頁面~~ | 🟡 中 | 6 小時 | ⚪ 無需實作 |
| LINE Integration 整合 | 🟡 中 | 4 小時 | 待實作 |
| Centers CRUD 完整實作 | 🟡 中 | 8 小時 | 待實作 |
| Rooms/Courses/Offerings 功能補齊 | 🟡 中 | 8 小時 | 待實作 |

### Phase 4：排課驗證與例外（一個月內）

| 工作項目 | 優先順序 | 預估時間 |
|:---|:---:|:---:|
| Scheduling Validation 整合 | 🟡 中 | 6 小時 |
| Exceptions Management 補齊 | 🟡 中 | 4 小時 |

### Phase 5：擴充功能（兩個月內）

| 工作項目 | 優先順序 | 預估時間 |
|:---|:---:|:---:|
| Calendar & Export 功能 | 🟢 低 | 8 小時 |
| Public Invitations 頁面 | 🟢 低 | 4 小時 |
| Notification Queue 監控 | 🟢 低 | 2 小時 |

### Phase 6：技術債務（二個月內）

| 工作項目 | 優先順序 | 預估時間 |
|:---|:---:|:---:|
| TypeScript 型別定義統一 | 🟡 中 | 16 小時 |
| 錯誤處理機制統一 | 🟡 中 | 8 小時 |
| Loading 狀態管理 | 🟡 中 | 4 小時 |

---

## 進度追蹤

### 統計摘要

| 指標 | 數值 |
|:---|:---:|
| 總 API 數量 | 147 |
| 已實作 | 85 (57.8%) |
| 待實作 | 62 (42.2%) |
| 緊急修復完成率 | 1/2 (50%) |
| 功能確認進度 | 0/3 (0%) |

### 目標

- **短期目標**（2 週）：達到 70% API 覆蓋率（103 個）
- **中期目標**（1 個月）：達到 80% API 覆蓋率（118 個）
- **長期目標**（3 個月）：達到 85% API 覆蓋率（125 個）

---

## 更新記錄

| 日期 | 版本 | 變更內容 |
|:---|:---:|:---|
| 2026-01-31 | 1.0 | 初版建立 |
| 2026-02-07 | - | 下次進度檢視 |

---

## 檢查清單

### 緊急修復檢查清單

- [ ] Template Cells Reorder API 後端路由已註冊
- [ ] TimetableTemplateController.ReorderCells 方法已實作
- [ ] 前端可成功呼叫並儲存排序結果

### 核心功能檢查清單

- [ ] Teacher Invitations Store 方法已補齊
- [ ] 邀請列表正確顯示
- [ ] 接受/拒絕功能正常
- [ ] Hashtags 功能已確認需求
- [ ] Smart Matching Store 已建立
- [ ] 智慧媒合功能可正常運作

### 管理功能檢查清單

- [] ~~Admin User Store 已建立~~
- [] ~~Admin Settings 頁面已實作~~
- [ ] LINE 綁定功能已整合
- [ ] Centers CRUD 功能完整
- [ ] Rooms/Courses/Offerings 功能完整

### 排課功能檢查清單

- [ ] Scheduling Validation 已整合到表單
- [ ] Exceptions Management 功能完整
- [ ] 即時驗證功能正常運作

### 技術債務檢查清單

- [ ] API Response 型別已統一定義
- [ ] 錯誤處理機制已統一
- [ ] Loading 狀態已正確管理

