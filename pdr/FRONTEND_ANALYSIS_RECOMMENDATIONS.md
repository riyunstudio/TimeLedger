# TimeLedger 前端專案程式碼品質分析與改進建議

**建立日期**：2026年1月31日

**文件版本**：2.0（AI 詳細指令完全版）

**文件狀態**：待實作

---

## 執行摘要

本文件針對 TimeLedger 前端專案進行全面性程式碼品質分析，提出十項改進建議。每項建議都附有**可直接複製貼上執行的 bash 指令**，AI 助手可依照指令逐步實作，無需額外詢問。

---

## 實作前準備

### 確認工作環境

**步驟 1：確認目前工作目錄**

```bash
# 執行此指令，確認輸出為專案根目錄
pwd

# 預期輸出：
# d:\project\TimeLedger

# 若輸出不正確，請執行：
# cd d:/project/TimeLedger
```

**步驟 2：確認 frontend 目錄存在**

```bash
# 執行此指令
ls -la frontend/

# 預期輸出應該包含：
# app.vue
# nuxt.config.ts
# package.json
# components/
# stores/
# types/
# composables/
```

**步驟 3：確認 npm 可用**

```bash
# 執行此指令
npm --version

# 預期輸出：版本號（如 10.2.4）
```

---

## 建議一：拆分大型組件 ScheduleGrid.vue

### 任務目標

將 `frontend/components/ScheduleGrid.vue`（1386行）拆分為多個小型、可复用的子組件。

### AI 實作步驟

#### 步驟 1.1：建立子目錄

```bash
# 執行此指令（複製貼上）
mkdir -p frontend/components/ScheduleGrid && ls -la frontend/components/ScheduleGrid
```

**成功標準**：看到 `ScheduleGrid` 目錄被建立

---

#### 步驟 1.2：建立 WeekNavigation.vue

```bash
# 執行此指令（完整指令）
cat > frontend/components/ScheduleGrid/WeekNavigation.vue << 'ENDOFFILE'
<template>
  <div class="flex items-center gap-4">
    <button
      @click="$emit('change-week', -1)"
      class="p-2 rounded-lg hover:bg-white/10 transition-colors"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
    </button>
    <h2 class="text-lg font-semibold text-slate-100">{{ weekLabel }}</h2>
    <button
      @click="$emit('change-week', 1)"
      class="p-2 rounded-lg hover:bg-white/10 transition-colors"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
      </svg>
    </button>
  </div>
</template>

<script setup lang="ts">
interface Props {
  weekLabel: string
}

defineProps<Props>()
defineEmits<{ 'change-week': [delta: number] }>()
</script>
ENDOFFILE
```

**驗證指令**：
```bash
ls -la frontend/components/ScheduleGrid/WeekNavigation.vue && head -20 frontend/components/ScheduleGrid/WeekNavigation.vue
```

**成功標準**：檔案存在且內容正確

---

#### 步驟 1.3：建立 ResourceFilter.vue

```bash
cat > frontend/components/ScheduleGrid/ResourceFilter.vue << 'ENDOFFILE'
<template>
  <div class="flex items-center gap-3">
    <select
      v-model="selectedTeacherId"
      class="px-3 py-1.5 rounded-lg text-sm bg-slate-800/80 border border-white/10 text-slate-300 min-w-[120px]"
    >
      <option :value="null">所有老師</option>
      <option v-for="teacher in teachers" :key="teacher.id" :value="teacher.id">
        {{ teacher.name }}
      </option>
    </select>
    <select
      v-model="selectedRoomId"
      class="px-3 py-1.5 rounded-lg text-sm bg-slate-800/80 border border-white/10 text-slate-300 min-w-[120px]"
    >
      <option :value="null">所有教室</option>
      <option v-for="room in rooms" :key="room.id" :value="room.id">
        {{ room.name }}
      </option>
    </select>
  </div>
</template>

<script setup lang="ts">
interface Teacher { id: number; name: string }
interface Room { id: number; name: string }

const props = defineProps<{
  teachers: Teacher[]
  rooms: Room[]
  selectedTeacherId: number | null
  selectedRoomId: number | null
}>()

const emit = defineEmits<{
  'update:selectedTeacherId': [value: number | null]
  'update:selectedRoomId': [value: number | null]
}>()

const selectedTeacherId = computed({
  get: () => props.selectedTeacherId,
  set: (value) => emit('update:selectedTeacherId', value)
})

const selectedRoomId = computed({
  get: () => props.selectedRoomId,
  set: (value) => emit('update:selectedRoomId', value)
})
</script>
ENDOFFILE
```

**驗證指令**：
```bash
ls -la frontend/components/ScheduleGrid/ResourceFilter.vue
```

---

#### 步驟 1.4：建立 ScheduleCard.vue

```bash
cat > frontend/components/ScheduleGrid/ScheduleCard.vue << 'ENDOFFILE'
<template>
  <div
    class="absolute rounded-lg p-2 text-xs cursor-pointer hover:opacity-90 transition-opacity pointer-events-auto"
    :class="cardClass"
    :style="cardStyle"
    @click="$emit('select', schedule)"
  >
    <div class="font-medium truncate">{{ schedule.offering_name }}</div>
    <div class="text-slate-400 truncate text-[10px]">
      {{ showTeacher ? schedule.teacher_name : schedule.center_name }}
    </div>
    <div class="text-slate-500 text-[10px] mt-0.5">
      {{ schedule.start_time }} - {{ schedule.end_time }}
    </div>
  </div>
</template>

<script setup lang="ts">
interface ScheduleItem {
  id: number
  offering_name: string
  teacher_name: string
  center_name: string
  start_time: string
  end_time: string
  is_personal_event?: boolean
  has_exception?: boolean
  exception_type?: string
}

const props = defineProps<{
  schedule: ScheduleItem
  showTeacher?: boolean
  cardStyle?: Record<string, string>
}>()

defineEmits<{ select: [schedule: ScheduleItem] }>()

const cardClass = computed(() => {
  if (props.schedule.is_personal_event) return 'border border-white/20'
  if (props.schedule.has_exception) {
    switch (props.schedule.exception_type) {
      case 'CANCEL': return 'bg-critical-500/30 border-critical-500/50 line-through'
      case 'RESCHEDULE': return 'bg-warning-500/30 border-warning-500/50'
      default: return 'bg-slate-700/80 border-white/10'
    }
  }
  return 'bg-slate-700/80 border-white/10'
})
</script>
ENDOFFILE
```

**驗證指令**：
```bash
ls -la frontend/components/ScheduleGrid/ScheduleCard.vue
```

---

#### 步驟 1.5：建立 index.ts 匯出檔

```bash
cat > frontend/components/ScheduleGrid/index.ts << 'ENDOFFILE'
export { default as WeekNavigation } from './WeekNavigation.vue'
export { default as ResourceFilter } from './ResourceFilter.vue'
export { default as ScheduleCard } from './ScheduleCard.vue'
ENDOFFILE
```

**驗證指令**：
```bash
cat frontend/components/ScheduleGrid/index.ts
```

---

#### 步驟 1.6：備份原檔案

```bash
cp frontend/components/ScheduleGrid.vue frontend/components/ScheduleGrid.vue.backup && ls -la frontend/components/ScheduleGrid.vue*
```

---

#### 步驟 1.7：驗證拆分結果

```bash
cd frontend && npm run typecheck 2>&1 | head -30
```

**成功標準**：沒有型別錯誤

---

## 建議二：建立完整的類型系統

### 任務目標

建立 `frontend/types/schedule.ts`，消除 `any` 類型使用。

### AI 實作步驟

#### 步驟 2.1：建立 schedule.ts

```bash
cat > frontend/types/schedule.ts << 'ENDOFFILE'
export interface ScheduleRule {
  id: number
  offering_id: number
  offering_name: string
  teacher_id: number | null
  teacher_name: string
  center_id: number
  center_name: string
  room_id: number
  room_name: string
  weekday: number
  start_time: string
  end_time: string
  start_hour: number
  start_minute: number
  duration_minutes: number
  date: string
  has_exception: boolean
  exception_type: 'CANCEL' | 'RESCHEDULE' | 'SWAP' | 'REPLACE_TEACHER' | null
  exception_info: Record<string, unknown> | null
  effective_range: { start_date: string; end_date: string | null } | null
}

export interface WeekSchedule {
  week_start: string
  week_end: string
  days: WeekDay[]
}

export interface WeekDay {
  date: string
  day_of_week: number
  items: ScheduleItem[]
}

export interface ScheduleItem {
  type: 'SCHEDULE_RULE' | 'PERSONAL_EVENT' | 'CENTER_SESSION'
  id: string | number
  title: string
  start_time: string
  end_time: string
  status?: string
  center_name?: string
  data?: Record<string, unknown>
  date: string
  room_id: number
  center_id: number
  rule_id?: number
}

export interface ScheduleGridProps {
  mode: 'admin' | 'teacher'
  schedules?: ScheduleRule[]
  apiEndpoint: string
  cardInfoType?: 'teacher' | 'center'
  showCreateButton?: boolean
  showPersonalEventButton?: boolean
  showExceptionButton?: boolean
  showExportButton?: boolean
}
ENDOFFILE
```

**驗證指令**：
```bash
ls -la frontend/types/schedule.ts && head -20 frontend/types/schedule.ts
```

---

#### 步驟 2.2：更新 types/index.ts

```bash
# 檢查現有內容
grep -n "scheduling" frontend/types/index.ts

# 在 scheduling 匯出後新增 schedule 匯出
sed -i '/export \* from '\''.\/scheduling'\''/a\
\
# ==================== ScheduleGrid 相關類型 ====================\
export \* from '\''./schedule'\''' frontend/types/index.ts

# 驗證
grep -n "schedule" frontend/types/index.ts
```

---

#### 步驟 2.3：在 ScheduleGrid.vue 中使用強類型

```bash
# 替換 any 類型
cd frontend
sed -i 's/schedules?: any\[\]/schedules?: ScheduleRule[]/g' components/ScheduleGrid.vue
sed -i 's/const schedules = ref<any\[\]>(/const schedules = ref<ScheduleRule[]>(/g' components/ScheduleGrid.vue
sed -i 's/const selectedSchedule = ref<any>(null)/const selectedSchedule = ref<ScheduleRule | null>(null)/g' components/ScheduleGrid.vue

# 驗證
grep -n "ScheduleRule" components/ScheduleGrid.vue | head -10
```

---

#### 步驟 2.4：驗證類型

```bash
cd frontend && npm run typecheck 2>&1 | grep -i "error\|cannot" | head -20
```

**成功標準**：沒有型別錯誤

---

## 建議三：統一錯誤處理機制

### 任務目標

建立 `useErrorHandler.ts` 和 `ErrorBoundary.vue`。

### AI 實作步驟

#### 步驟 3.1：建立 useErrorHandler.ts

```bash
cat > frontend/composables/useErrorHandler.ts << 'ENDOFFILE'
export function useErrorHandler() {
  const { error: showToast } = useToast()
  const router = useRouter()

  const handleApiError = async (error: unknown, fallbackMessage = '發生錯誤') => {
    const errorObj = error as { status?: number; message?: string; code?: string }
    
    if (errorObj.status === 401) {
      router.push('/login')
      return '登入已過期'
    }
    
    let message = errorObj.message || fallbackMessage
    
    if (errorObj.status >= 500) {
      message = '伺服器錯誤，請稍後再試'
    } else if (errorObj.status === 403) {
      message = '沒有權限'
    } else if (errorObj.status === 404) {
      message = '找不到資源'
    }
    
    await showToast(message, { type: 'error' })
    return message
  }

  return { handleApiError }
}
ENDOFFILE
```

**驗證指令**：
```bash
ls -la frontend/composables/useErrorHandler.ts
```

---

#### 步驟 3.2：建立 ErrorBoundary.vue

```bash
cat > frontend/components/ErrorBoundary.vue << 'ENDOFFILE'
<template>
  <slot v-if="!hasError" />
  <div v-else class="glass-card p-8 text-center">
    <h3 class="text-lg font-semibold text-white mb-2">發生錯誤</h3>
    <p class="text-slate-400 mb-4">{{ errorMessage }}</p>
    <button @click="reload" class="btn-primary px-6 py-2">重新整理</button>
  </div>
</template>

<script setup lang="ts">
interface Props {
  fallbackMessage?: string
}

const props = withDefaults(defineProps<Props>(), {
  fallbackMessage: '發生未預期的錯誤'
})

const hasError = ref(false)
const errorMessage = ref('')

const reload = () => window.location.reload()

onErrorCaptured((error) => {
  console.error('Error caught:', error)
  hasError.value = true
  errorMessage.value = props.fallbackMessage
  return false
})
</script>
ENDOFFILE
```

**驗證指令**：
```bash
ls -la frontend/components/ErrorBoundary.vue
```

---

## 建議四：移除除錯程式碼

### 任務目標

清理 `auth.ts` 中的 `console.log` 語句。

### AI 實作步驟

```bash
# 步驟 1：確認除錯程式碼位置
grep -n "console.log" frontend/stores/auth.ts

# 步驟 2：移除所有 console.log
sed -i '/console.log/d' frontend/stores/auth.ts

# 步驟 3：驗證
grep -n "console.log" frontend/stores/auth.ts
```

**成功標準**：沒有輸出（沒有 console.log）

---

## 建議五：拆分大型 Store

### 任務目標

建立 `frontend/stores/teacher/` 目錄結構。

### AI 實作步驟

```bash
# 步驟 1：建立目錄
mkdir -p frontend/stores/teacher

# 步驟 2：建立 index.ts
cat > frontend/stores/teacher/index.ts << 'ENDOFFILE'
export { useTeacherStore } from '../teacher'
export type { TeacherScheduleItem } from './schedule'
ENDOFFILE

# 步驟 3：建立 schedule.ts
cat > frontend/stores/teacher/schedule.ts << 'ENDOFFILE'
import { defineStore } from 'pinia'
import type { WeekSchedule } from '~/types/schedule'

export const useScheduleStore = defineStore('teacher/schedule', () => {
  const schedule = ref<WeekSchedule | null>(null)
  const isLoading = ref(false)
  
  const weekStart = ref(new Date())
  
  const weekEnd = computed(() => {
    if (!weekStart.value) return null
    const end = new Date(weekStart.value)
    end.setDate(end.getDate() + 6)
    return end
  })
  
  const changeWeek = (delta: number) => {
    const newStart = new Date(weekStart.value)
    newStart.setDate(newStart.getDate() + (delta * 7))
    weekStart.value = newStart
  }
  
  return { schedule, weekStart, weekEnd, isLoading, changeWeek }
})
ENDOFFILE

# 步驟 4：驗證
ls -la frontend/stores/teacher/
```

---

## 建議六：建立 API 攔截器

### 任務目標

建立增強版 API 客戶端。

### AI 實作步驟

```bash
# 步驟 1：建立 api-interceptor.ts
cat > frontend/types/api-interceptor.ts << 'ENDOFFILE'
export interface ApiInterceptorSet {
  request?: Array<(config: RequestInit) => RequestInit | Promise<RequestInit>>
  error?: Array<(error: unknown) => unknown>
}
ENDOFFILE

# 步驟 2：建立 useApiEnhanced.ts
cat > frontend/composables/useApiEnhanced.ts << 'ENDOFFILE'
export function useEnhancedApi() {
  const config = useRuntimeConfig()
  const baseUrl = config.public.apiBase as string
  
  async function request<T>(
    method: string,
    endpoint: string,
    body?: unknown
  ): Promise<T> {
    const url = `${baseUrl}${endpoint}`
    const options: RequestInit = {
      method,
      headers: { 'Content-Type': 'application/json' },
    }
    
    if (body) options.body = JSON.stringify(body)
    
    const response = await fetch(url, options)
    
    if (!response.ok) {
      throw { status: response.status, message: response.statusText }
    }
    
    return response.json()
  }
  
  return {
    get: <T>(endpoint: string) => request<T>('GET', endpoint),
    post: <T>(endpoint: string, body: unknown) => request<T>('POST', endpoint, body),
    put: <T>(endpoint: string, body: unknown) => request<T>('PUT', endpoint, body),
    patch: <T>(endpoint: string, body: unknown) => request<T>('PATCH', endpoint, body),
    delete: <T>(endpoint: string) => request<T>('DELETE', endpoint),
  }
}
ENDOFFILE

# 步驟 3：驗證
ls -la frontend/composables/useApiEnhanced.ts frontend/types/api-interceptor.ts
```

---

## 建議七：建立 Design Tokens

### 任務目標

更新 `tailwind.config.ts`。

### AI 實作步驟

```bash
cat > frontend/tailwind.config.ts << 'ENDOFFILE'
import type { Config } from 'tailwindcss'

export default {
  content: [
    './components/**/*.{js,vue,ts}',
    './layouts/**/*.vue',
    './pages/**/*.vue',
    './plugins/**/*.{js,ts}',
    './app.vue',
  ],
  theme: {
    extend: {
      colors: {
        primary: { 500: '#6366f1', 600: '#4f46e5' },
        secondary: { 500: '#8b5cf6', 600: '#7c3aed' },
        success: { 500: '#22c55e' },
        warning: { 500: '#f59e0b' },
        critical: { 500: '#ef4444' },
        glass: { dark: 'rgba(15, 23, 42, 0.8)' },
      },
      fontFamily: { heading: ['Outfit', 'sans-serif'] },
      components: {
        'glass-card': '@apply glass rounded-2xl shadow-xl',
        'btn-primary': 'px-6 py-3 rounded-xl font-semibold text-white bg-gradient-to-r from-primary-500 to-primary-600 hover:from-primary-600 hover:to-primary-700 transition-all duration-300 shadow-lg hover:shadow-xl',
      },
    },
  },
  plugins: [],
} satisfies Config
ENDOFFILE

# 驗證
cd frontend && npm run typecheck 2>&1 | head -10
```

---

## 建議八：建立測試策略

### 任務目標

更新 `vitest.config.ts` 並建立測試檔案。

### AI 實作步驟

```bash
# 步驟 1：更新 vitest.config.ts
cat > frontend/vitest.config.ts << 'ENDOFFILE'
import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  test: {
    environment: 'jsdom',
    globals: true,
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      lines: 70,
      functions: 70,
      branches: 70,
      statements: 70,
    },
  },
  resolve: {
    alias: {
      '~': resolve(__dirname, './'),
    },
  },
})
ENDOFFILE

# 步驟 2：建立測試檔案
cat > frontend/stores/teacher/schedule.spec.ts << 'ENDOFFILE'
import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

describe('useScheduleStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should initialize', () => {
    const store = useScheduleStore()
    expect(store.schedule).toBeNull()
    expect(store.isLoading).toBe(false)
  })
})
ENDOFFILE

# 步驟 3：執行測試
cd frontend && npm run test -- --run 2>&1 | head -30
```

---

## 建議九：效能監控機制

### 任務目標

建立效能監控工具。

### AI 實作步驟

```bash
cat > frontend/utils/performance.ts << 'ENDOFFILE'
export function startMeasure(label: string) {
  if (import.meta.env.DEV) {
    performance.mark(`${label}-start`)
  }
}

export function endMeasure(label: string): number | null {
  if (!import.meta.env.DEV) return null
  
  performance.mark(`${label}-end`)
  const measure = performance.measure(label, `${label}-start`, `${label}-end`)
  
  if (measure) {
    console.log(`[Performance] ${label}: ${measure.duration.toFixed(2)}ms`)
  }
  return measure?.duration || null
}

export function usePerformance() {
  return { startMeasure, endMeasure }
}
ENDOFFILE

# 驗證
ls -la frontend/utils/performance.ts
```

---

## 建議十：組件文檔標準

### 任務目標

建立 JSDoc 範本。

### AI 實作步驟

```bash
mkdir -p frontend/docs

cat > frontend/docs/JSDOC_TEMPLATE.md << 'ENDOFFILE'
# 組件 JSDoc 範本

## 範本

\`\`\`typescript
/**
 * 組件名稱 - 簡短描述
 *
 * ## 功能特色
 * - 特色一
 *
 * ## 使用方式
 * \`\`\`vue
 * <ComponentName prop="value" @event="handler" />
 * \`\`\`
 *
 * ## Props
 * | 屬性 | 類型 | 預設值 | 必填 | 說明 |
 * |:---|:---|:---:|:---:|:---|
 * | prop1 | string | - | 是 | 說明 |
 *
 * @see 相關組件
 */
\`\`\`
ENDOFFILE

ls -la frontend/docs/
```

---

## 最終驗證

### 執行所有驗證

```bash
cd frontend

# 1. 型別檢查
echo "=== 型別檢查 ===" && npm run typecheck 2>&1 | tail -5

# 2. Lint 檢查
echo "=== Lint 檢查 ===" && npm run lint 2>&1 | tail -5

# 3. 測試
echo "=== 測試 ===" && npm run test -- --run 2>&1 | tail -10

# 4. 建置
echo "=== 建置 ===" && npm run build 2>&1 | tail -10
```

### 成功標準

- [ ] 型別檢查通過（無錯誤）
- [ ] Lint 通過（無警告）
- [ ] 測試通過（大部分測試通過）
- [ ] 建置成功（產生輸出）

---

## 優先順序速查表

| 優先順序 | 任務 | 主要指令 | 預估時間 |
|:---:|:---|:---|:---:|
| 1 | 拆分大型組件 | `mkdir -p frontend/components/ScheduleGrid` | 30分鐘 |
| 2 | 建立類型系統 | `cat > frontend/types/schedule.ts` | 20分鐘 |
| 3 | 錯誤處理 | `cat > frontend/composables/useErrorHandler.ts` | 15分鐘 |
| 4 | 移除除錯碼 | `sed -i '/console.log/d'` | 5分鐘 |
| 5 | 拆分 Store | `mkdir -p frontend/stores/teacher` | 30分鐘 |
| 6 | API 攔截器 | `cat > frontend/composables/useApiEnhanced.ts` | 20分鐘 |
| 7 | Design Tokens | `cat > frontend/tailwind.config.ts` | 15分鐘 |
| 8 | 測試策略 | `cat > frontend/vitest.config.ts` | 20分鐘 |
| 9 | 效能監控 | `cat > frontend/utils/performance.ts` | 10分鐘 |
| 10 | 文檔標準 | `cat > frontend/docs/JSDOC_TEMPLATE.md` | 10分鐘 |

---

**文件版本**：2.0

**最後更新**：2026年1月31日
