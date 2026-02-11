<template>
  <div class="p-4 md:p-6 max-w-7xl mx-auto">
    <!-- 頁面標題 -->
    <div class="mb-6 md:mb-8">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 class="text-2xl md:text-3xl font-bold text-slate-100 mb-2">
            資源佔用表
          </h1>
          <p class="text-slate-400 text-sm md:text-base">
            檢視老師與教室的每週排課佔用情況
          </p>
        </div>
        <button
          @click="openCopyRulesWizard"
          class="px-4 py-2 bg-primary-500/30 border border-primary-500 text-primary-400 rounded-xl hover:bg-primary-500/40 transition-colors flex items-center gap-2"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
          </svg>
          <span>複製規則</span>
        </button>
      </div>
    </div>

    <!-- 篩選欄 (使用 fixed 定位確保不被週曆表遮擋) -->
    <div class="glass-card p-4 mb-6 relative z-50">
      <div class="flex flex-col lg:flex-row gap-4">
        <!-- 學期選擇 (可選) -->
        <div class="lg:w-64">
          <label class="block text-sm font-medium text-slate-300 mb-2">
            學期 (快速跳轉)
          </label>
          <select
            v-model="selectedTermId"
            @change="handleTermChange"
            class="w-full px-4 py-2.5 bg-white/5 border border-white/10 rounded-xl text-slate-100 focus:outline-none focus:border-primary-500 transition-colors appearance-none"
          >
            <option :value="0">選擇學期...</option>
            <option v-for="term in terms" :key="term.id" :value="term.id">
              {{ term.name }}
            </option>
          </select>
        </div>

        <!-- 日期區間選擇 -->
        <div class="flex-1">
          <label class="block text-sm font-medium text-slate-300 mb-2">
            日期區間
          </label>
          <div class="flex items-center gap-2">
            <input
              v-model="queryStartDate"
              type="date"
              class="flex-1 px-4 py-2.5 bg-white/5 border border-white/10 rounded-xl text-slate-100 focus:outline-none focus:border-primary-500 transition-colors"
              @change="fetchOccupancyRules"
            />
            <span class="text-slate-500">至</span>
            <input
              v-model="queryEndDate"
              type="date"
              class="flex-1 px-4 py-2.5 bg-white/5 border border-white/10 rounded-xl text-slate-100 focus:outline-none focus:border-primary-500 transition-colors"
              @change="fetchOccupancyRules"
            />
          </div>
        </div>

        <!-- 資源類型 -->
        <div class="lg:w-40">
          <label class="block text-sm font-medium text-slate-300 mb-2">
            資源類型
          </label>
          <div class="flex p-1 bg-white/5 border border-white/10 rounded-xl">
            <button
              @click="resourceType = 'teacher'"
              class="flex-1 py-1.5 text-sm font-medium rounded-lg transition-all"
              :class="resourceType === 'teacher' ? 'bg-primary-500 text-white shadow-lg' : 'text-slate-400 hover:text-slate-200'"
            >
              老師
            </button>
            <button
              @click="resourceType = 'room'"
              class="flex-1 py-1.5 text-sm font-medium rounded-lg transition-all"
              :class="resourceType === 'room' ? 'bg-primary-500 text-white shadow-lg' : 'text-slate-400 hover:text-slate-200'"
            >
              教室
            </button>
          </div>
        </div>

        <!-- 資源搜尋 -->
        <div class="flex-1 relative">
          <SearchableSelect
            v-model="selectedResourceIds"
            :options="resourceOptions"
            :placeholder="`選擇或搜尋${resourceType === 'teacher' ? '老師' : '教室'}...`"
            :label="`搜尋${resourceType === 'teacher' ? '老師' : '教室'}`"
            :multiple="true"
          />
        </div>
      </div>
    </div>

    <!-- 佔用率統計 -->
    <div v-if="selectedResource.length > 0" class="mt-6 grid grid-cols-2 md:grid-cols-4 gap-4">
      <div class="glass-card p-4">
        <div class="text-sm text-slate-400 mb-1">總課程數</div>
        <div class="text-2xl font-bold text-slate-100">{{ occupancyRules.length }}</div>
      </div>
      <div class="glass-card p-4">
        <div class="text-sm text-slate-400 mb-1">佔用天數</div>
        <div class="text-2xl font-bold text-slate-100">{{ occupiedDays }}</div>
      </div>
      <div class="glass-card p-4">
        <div class="text-sm text-slate-400 mb-1">衝突次數</div>
        <div class="text-2xl font-bold text-red-400">{{ conflictCount }}</div>
      </div>
      <div class="glass-card p-4">
        <div class="text-sm text-slate-400 mb-1">佔用率</div>
        <div class="text-2xl font-bold text-primary-400">{{ occupancyRate }}%</div>
      </div>
    </div>

    <!-- 每週網格視圖 - 借鑑首頁 WeekGrid 風格 -->
    <div v-if="selectedResource.length > 0 && occupancyRules.length > 0" class="glass-card p-4 mt-6">
      <!-- 週網格標題 -->
      <div class="grid grid-cols-8 gap-2 mb-2">
        <div class="text-center text-sm text-slate-400 font-medium">
          時間
        </div>
        <div
          v-for="day in weekDays"
          :key="day.value"
          class="text-center text-sm text-slate-200 font-medium py-2"
        >
          {{ day.label }}
        </div>
      </div>

      <!-- 時段網格 -->
      <div class="grid grid-cols-8 gap-2">
        <!-- 時間標籤 - 使用與網格相同的 h-16 高度，移除 space-y-2 改用直接排列 -->
        <div>
          <div
            v-for="hour in visibleTimeSlots"
            :key="hour"
            class="h-16 flex items-center justify-center text-xs md:text-sm text-slate-500 border-t border-white/5"
            :class="{ 'border-t-0': hour === visibleTimeSlots[0] }"
          >
            {{ formatTimeLabel(hour) }}
          </div>
        </div>

        <!-- 每天的時段 - 真正的週曆風格 -->
        <div
          v-for="day in weekDays"
          :key="day.value"
          class="relative"
        >
          <!-- 每小時的網格線 -->
          <div
            v-for="hour in visibleTimeSlots"
            :key="`${day.value}-${hour}-grid`"
            class="h-16 border-t border-white/5 w-full"
            :class="{ 'border-t-0': hour === visibleTimeSlots[0] }"
          />

          <!-- 課程卡片 - 根據實際時間絕對定位 -->
          <div
            v-for="rule in getRulesForDay(day.value)"
            :key="rule.id"
            class="absolute left-0 right-0 px-2 py-1 text-xs rounded-lg cursor-pointer hover:opacity-80 transition-opacity overflow-hidden"
            :class="getRuleClass(rule)"
            :style="getRuleStyle(rule)"
            draggable="true"
            @click="handleSlotClick(rule)"
            @dragstart="handleDragStart($event, rule)"
            @dragend="handleDragEnd($event)"
            @dragover.prevent="handleDragOver($event, rule)"
            @drop="handleDrop($event, rule)"
          >
            <!-- 衝突警示 -->
            <div
              v-if="hasConflict(rule)"
              class="absolute -top-1 -right-1 w-4 h-4 bg-red-500 rounded-full flex items-center justify-center z-10"
            >
              <svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            </div>
            <!-- 拖拽時的視覺提示 -->
            <div
              v-if="isDragging && draggedRule?.id === rule.id"
              class="absolute inset-0 bg-primary-500/50 border-2 border-primary-400 border-dashed rounded"
            >
              <svg class="w-4 h-4 mx-auto mt-1 text-primary-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4" />
              </svg>
            </div>
            <!-- 拖拽目標指示器 -->
            <div
              v-if="dragOverSlot && dragOverSlot.ruleId === rule.id"
              class="absolute bottom-0 left-0 right-0 h-1 bg-primary-400 animate-pulse"
            />
            <!-- 課程名稱 -->
            <div class="font-medium truncate">
              {{ rule.offering_name || rule.course_name }}
            </div>
            <!-- 關聯資源名稱 -->
            <div class="text-[10px] opacity-75 truncate">
              {{ resourceType === 'teacher' ? rule.room_name : rule.teacher_name }}
            </div>
            <!-- 時間 -->
            <div class="text-[10px] opacity-75">
              {{ rule.start_time }} - {{ rule.end_time }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 無數據提示 -->
    <div v-else-if="selectedResource.length > 0 && !loadingRules" class="glass-card p-12 text-center">
      <div class="w-16 h-16 bg-slate-800 rounded-full flex items-center justify-center mx-auto mb-4 border border-white/5">
        <svg class="w-8 h-8 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
      <h3 class="text-lg font-medium text-slate-300 mb-2">尚無排課記錄</h3>
      <p class="text-slate-500 max-w-sm mx-auto mb-6">
        目前 {{ resourceType === 'teacher' ? '這位老師' : '這些教室' }} 在此週段內沒有任何排課安排。
      </p>
      <button
        @click="navigateToCreateRule()"
        class="px-6 py-2 bg-primary-500 text-white rounded-xl hover:bg-primary-600 transition-colors"
      >
        新增排課規則
      </button>
    </div>

    <!-- 複製規則精靈 -->
    <CopyRulesWizard
      v-if="showCopyWizard"
      :terms="termsList"
      @close="showCopyWizard = false"
      @copied="handleCopyComplete"
    />
  </div>
</template>

<script setup lang="ts">
import type { Term, OccupancyRule } from '~/types/scheduling'
import CopyRulesWizard from '~/components/Admin/CopyRulesWizard.vue'
import SearchableSelect, { type SelectOption } from '~/components/Common/SearchableSelect.vue'
import { useResourceCache } from '~/composables/useResourceCache'

/**
 * 資源（老師或教室）
 */
interface Resource {
  id: number
  name: string
  count?: number
  capacity?: number
}

definePageMeta({
  auth: 'ADMIN',
  layout: 'admin',
})

// ========== State ==========

const api = useApi()
const router = useRouter()
const config = useRuntimeConfig()
const { success, error: toastError } = useToast()
const { warning: alertWarning, confirm: alertConfirm } = useAlert()

// API 基礎 URL
const API_BASE = config.public.apiBase

// Center ID
const { getCenterId } = useCenterId()

// 學期列表
const terms = ref<Term[]>([])
const selectedTermId = ref<number>(0)

// 查詢日期區間
const queryStartDate = ref('')
const queryEndDate = ref('')

// Initialize dats to current week
const today = new Date()
const day = today.getDay()
const diff = today.getDate() - day + (day === 0 ? -6 : 1) // 調整到週一
const monday = new Date(today.setDate(diff))
const sunday = new Date(monday)
sunday.setDate(monday.getDate() + 6)
queryStartDate.value = monday.toISOString().split('T')[0]
queryEndDate.value = sunday.toISOString().split('T')[0]

// 資源類型
const resourceType = ref<'teacher' | 'room'>('teacher')

// 資源快取
const { resourceCache, fetchAllResources } = useResourceCache()

// 選中的資源 ID（單選或多選）
const selectedResourceIds = ref<(number | string)[]>([])
const selectedResource = ref<Resource[]>([])

// 佔用規則
const loadingRules = ref(false)

// 定義後端返回的分組結構
interface GroupedOccupancy {
  day_of_week: number
  day_name: string
  rules: OccupancyRule[]
}

const occupancyRules = ref<OccupancyRule[]>([])

// 中心營業時間設定
const operatingStartTime = ref('00:00')
const operatingEndTime = ref('23:00')

// 複製規則精靈
const showCopyWizard = ref(false)

// ========== 拖拽狀態 ==========
const isDragging = ref(false)
const draggedRule = ref<OccupancyRule | null>(null)
const dragOverSlot = ref<{ weekday: number; hour: number; ruleId: number } | null>(null)

// ========== 常量 ==========

// 週日期
const weekDays = [
  { label: '週一', value: 1 },
  { label: '週二', value: 2 },
  { label: '週三', value: 3 },
  { label: '週四', value: 4 },
  { label: '週五', value: 5 },
  { label: '週六', value: 6 },
  { label: '週日', value: 7 },
]

// 時間槽 (基礎陣列)
const timeSlots = Array.from({ length: 24 }, (_, i) => i)

// 顯示用的時間槽（根據中心營業時間動態生成）
const visibleTimeSlots = computed(() => {
  const startStr = operatingStartTime.value
  const endStr = operatingEndTime.value

  // 解析開始時間，若無效則使用預設值 0
  let start = 0
  if (startStr && startStr.includes(':')) {
    const parsed = parseInt(startStr.split(':')[0], 10)
    start = isNaN(parsed) ? 0 : parsed
  }

  // 解析結束時間，若無效則使用預設值 23
  let end = 23
  if (endStr && endStr.includes(':')) {
    const parsed = parseInt(endStr.split(':')[0], 10)
    end = isNaN(parsed) ? 23 : parsed
  }

  // 確保開始時間在有效範圍內 (0-23)
  start = Math.max(0, Math.min(23, start))
  end = Math.max(0, Math.min(23, end))

  // 確保結束時間 >= 開始時間
  const actualEnd = end < start ? 23 : end

  const slots = []
  for (let i = start; i <= actualEnd; i++) {
    slots.push(i)
  }

  // 如果沒有有效的時間段，回退到預設值
  if (slots.length === 0) {
    return Array.from({ length: 24 }, (_, i) => i)
  }

  return slots
})

// ========== Computed ==========

/**
 * 老師選項列表
 */
const teacherOptions = computed<SelectOption[]>(() => {
  return Array.from(resourceCache.value.teachers.values()).map((teacher: any) => ({
    id: teacher.id,
    name: teacher.name || `老師 ${teacher.id}`
  }))
})

/**
 * 教室選項列表
 */
const roomOptions = computed<SelectOption[]>(() => {
  return Array.from(resourceCache.value.rooms.values()).map((room: any) => ({
    id: room.id,
    name: room.name || `教室 ${room.id}`
  }))
})

/**
 * 根據資源類型取得選項列表
 */
const resourceOptions = computed<SelectOption[]>(() => {
  return resourceType.value === 'teacher' ? teacherOptions.value : roomOptions.value
})

/**
 * 學期列表（確保是陣列格式）
 */
const termsList = computed(() => {
  return Array.isArray(terms.value) ? terms.value : []
})

/**
 * 佔用天數
 */
const occupiedDays = computed(() => {
  const days = new Set(occupancyRules.value.map(r => r.weekday))
  return days.size
})

/**
 * 衝突次數
 */
const conflictCount = computed(() => {
  let count = 0
  const rules = occupancyRules.value

  for (let i = 0; i < rules.length; i++) {
    for (let j = i + 1; j < rules.length; j++) {
      // 只檢查同一天，且時段有重疊
      if (rules[i].weekday === rules[j].weekday) {
        if (checkOverlap(rules[i], rules[j])) {
          count++
        }
      }
    }
  }
  return count
})

/**
 * 佔用率
 */
const occupancyRate = computed(() => {
  if (occupancyRules.value.length === 0) return 0

  // 計算實際涵蓋的天數（週一到週日都算）
  const days = new Set(occupancyRules.value.map(r => r.weekday))
  const activeDays = days.size > 0 ? days.size : 1

  // 計算這段時間內每個rule佔用的小時數
  let occupiedMinutes = 0
  occupancyRules.value.forEach(r => {
    const toMin = (s: string) => {
      const parts = s.split(':').map(Number)
      return parts[0] * 60 + parts[1]
    }
    const startMin = toMin(r.start_time)
    const endMin = toMin(r.end_time)
    if (endMin > startMin) {
      occupiedMinutes += (endMin - startMin)
    }
  })

  // 總可用分鐘數：涵蓋的天數 * 每小時60分鐘 * 12小時 (7:00-19:00)
  const totalMinutes = activeDays * 60 * 12
  const rate = Math.round((occupiedMinutes / totalMinutes) * 100)

  return Math.max(0, Math.min(100, rate))
})

// ========== Methods ==========

/**
 * 取得學期列表
 */
const fetchTerms = async () => {
  try {
    const response = await api.get<Term[]>('/admin/terms')
    if (response) {
      terms.value = response
    }
  } catch (err) {
    console.error('Failed to fetch terms:', err)
    toastError('載入學期失敗')
  }
}

/**
 * 取得中心設定（營業時間）
 */
const fetchCenterSettings = async () => {
  try {
    const centerId = getCenterId()
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${API_BASE}/admin/centers/${centerId}/settings`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })

    if (!response.ok) {
      console.error('取得中心設定失敗: HTTP', response.status)
      return
    }

    const text = await response.text()
    if (!text.trim()) {
      console.error('取得中心設定失敗: 空响应')
      return
    }

    const data = JSON.parse(text)
    if (data.datas) {
      // 取得營業時間設定，若無效則使用預設值
      const apiStartTime = data.datas.operating_start_time
      const apiEndTime = data.datas.operating_end_time

      // 驗證並設置開始時間
      if (apiStartTime && apiStartTime.includes(':')) {
        const hour = parseInt(apiStartTime.split(':')[0], 10)
        if (!isNaN(hour) && hour >= 0 && hour <= 23) {
          operatingStartTime.value = apiStartTime
        }
      }

      // 驗證並設置結束時間
      if (apiEndTime && apiEndTime.includes(':')) {
        const hour = parseInt(apiEndTime.split(':')[0], 10)
        if (!isNaN(hour) && hour >= 0 && hour <= 23) {
          operatingEndTime.value = apiEndTime
        }
      }
    }
  } catch (err) {
    console.error('取得中心設定失敗:', err)
  }
}

/**
 * 處理學期切換
 */
const handleTermChange = () => {
  const term = terms.value.find(t => t.id === selectedTermId.value)
  if (term) {
    // 自動帶入學期起訖日
    queryStartDate.value = term.start_date.split('T')[0]
    queryEndDate.value = term.end_date.split('T')[0]
    fetchOccupancyRules()
  }
}

/**
 * 監聽資源選擇變化
 */
watch(selectedResourceIds, (newIds) => {
  if (!newIds || newIds.length === 0) {
    selectedResource.value = []
    occupancyRules.value = []
    return
  }
  const ids = newIds as number[]
  // 從快取中取得對應的物件
  if (resourceType.value === 'teacher') {
    selectedResource.value = ids.map((id: number | string) => {
      const teacher = resourceCache.value.teachers.get(Number(id))
      if (!teacher) return null
      const resource: Resource = {
        id: teacher.id as number,
        name: (teacher.name || `老師 ${teacher.id}`) as string,
      }
      if (teacher.count !== undefined) {
        resource.count = teacher.count as number
      }
      return resource
    }).filter((r: Resource | null): r is Resource => r !== null)
  } else {
    selectedResource.value = ids.map((id: number | string) => {
      const room = resourceCache.value.rooms.get(Number(id))
      if (!room) return null
      const resource: Resource = {
        id: room.id as number,
        name: (room.name || `教室 ${room.id}`) as string,
      }
      if (room.capacity !== undefined) {
        resource.capacity = room.capacity as number
      }
      return resource
    }).filter((r: Resource | null): r is Resource => r !== null)
  }
  // 立即取得佔用規則
  fetchOccupancyRules()
})

/**
 * 監聽資源類型切換，清除選中的資源並載入資料
 */
watch(resourceType, async () => {
  selectedResourceIds.value = []
  selectedResource.value = []
  occupancyRules.value = []

  // 確保資源已載入
  await fetchAllResources()
})

/**
 * 取得佔用規則
 */
const fetchOccupancyRules = async () => {
  if (selectedResource.value.length === 0) return

  loadingRules.value = true
  try {
    const roomIds = resourceType.value === 'room'
      ? (selectedResource.value.map(r => r.id) as number[])
      : undefined
    const teacherIds = resourceType.value === 'teacher'
      ? (selectedResource.value.map(r => r.id) as number[])
      : undefined

    const response = await api.get<GroupedOccupancy[]>('/admin/occupancy/rules', {
      start_date: queryStartDate.value,
      end_date: queryEndDate.value,
      ...(teacherIds && { teacher_ids: teacherIds.join(',') }),
      ...(roomIds && { room_ids: roomIds.join(',') }),
    })

    if (response) {
      // 將分組數據平舖
      occupancyRules.value = response.flatMap(group => group.rules)
    }
  } catch (err) {
    console.error('Failed to fetch occupancy rules:', err)
    toastError('載入佔用規則失敗')
  } finally {
    loadingRules.value = false
  }
}

/**
 * 獲取指定時段的規則（按小時格子）- 保留用於向後相容
 */
const getRulesForSlot = (day: number, hour: number): OccupancyRule[] => {
  return occupancyRules.value.filter(rule => {
    if (rule.weekday !== day) return false
    const startHour = parseInt(rule.start_time.split(':')[0])
    return startHour === hour
  })
}

/**
 * 獲取指定日期的所有規則（用於連續時間軸）
 */
const getRulesForDay = (day: number): OccupancyRule[] => {
  return occupancyRules.value.filter(rule => rule.weekday === day)
}

/**
 * 檢查指定時間格子是否為空（用於連續時間軸）
 * 只有當整個格子都沒有課程時才視為空
 */
const isEmptySlot = (day: number, hour: number): boolean => {
  const slotStart = hour * 60
  const slotEnd = (hour + 1) * 60

  // 檢查是否有課程佔用這個時間格子的任何部分
  return !occupancyRules.value.some(rule => {
    if (rule.weekday !== day) return false
    const ruleStart = parseInt(rule.start_time.split(':')[0]) * 60 + parseInt(rule.start_time.split(':')[1])
    const ruleEnd = parseInt(rule.end_time.split(':')[0]) * 60 + parseInt(rule.end_time.split(':')[1])
    // 課程與格子有任何重疊就視為不空
    return ruleStart < slotEnd && ruleEnd > slotStart
  })
}

/**
 * 檢查指定時間是否為有效的放置目標
 */
const isValidDropTargetForTime = (day: number, hour: number): boolean => {
  return true
}

/**
 * 獲取時段樣式
 */
const getSlotClass = (day: number, hour: number): string => {
  const hasRules = getRulesForSlot(day, hour).length > 0
  if (hasRules) return 'bg-slate-800/50 border-white/5'
  return 'bg-white/5 border-white/5'
}

/**
 * 獲取規則卡片樣式
 */
const getRuleClass = (rule: OccupancyRule): string => {
  if (hasConflict(rule)) return 'bg-red-500/20 text-red-100 border-red-500/50'
  return 'bg-primary-500/20 text-primary-100 border-primary-500/50'
}

/**
 * 獲取規則卡片精確樣式 (根據可見時間槽定位)
 */
const getRuleStyle = (rule: OccupancyRule) => {
  const startParts = rule.start_time.split(':').map(Number)
  const endParts = rule.end_time.split(':').map(Number)

  // 確保時間部分是有效的數字
  const startHour = isNaN(startParts[0]) ? 0 : startParts[0]
  const startMinute = isNaN(startParts[1]) ? 0 : startParts[1]
  const endHour = isNaN(endParts[0]) ? startHour + 1 : endParts[0]
  const endMinute = isNaN(endParts[1]) ? 0 : endParts[1]

  // 將時間轉換為分鐘
  const startMinutes = startHour * 60 + startMinute
  let endMinutes = endHour * 60 + endMinute

  // 處理跨夜課程：如果結束時間早於開始時間，表示跨越午夜
  let duration = endMinutes - startMinutes
  if (duration <= 0) {
    // 跨越午夜：加上一整天（24小時 = 1440分鐘）
    endMinutes = endMinutes + 24 * 60
    duration = endMinutes - startMinutes
  }

  // 每小時格子的高度是 64px
  const hourHeight = 64

  // 計算相對於 visibleTimeSlots 起始時間的位置
  const minVisibleHour = visibleTimeSlots.value.length > 0 ? visibleTimeSlots.value[0] : 0

  // 確保 minVisibleHour 是有效的數字
  const safeMinHour = isNaN(minVisibleHour) ? 0 : minVisibleHour

  // 計算卡片頂部位置
  let topOffset = (startHour - safeMinHour) * hourHeight + startMinute * (hourHeight / 60)
  const cardHeight = duration * (hourHeight / 60)

  // 保護：如果課程開始時間早於營業時間，將其夾緊到 0
  const clampedTop = Math.max(0, topOffset)

  // 確保高度至少 20px，且不會超出範圍
  const clampedHeight = Math.max(Math.min(cardHeight, 24 * hourHeight - clampedTop), 20)

  return {
    top: `${clampedTop}px`,
    height: `${clampedHeight}px`,
    zIndex: 10
  }
}

/**
 * 格式化時間標籤
 */
const formatTimeLabel = (hour: number) => {
  return `${hour.toString().padStart(2, '0')}:00`
}

/**
 * 檢查是否在學期範圍內
 */
const isWithinTerm = (weekday: number, hour: number): boolean => {
  return true
}

// ========== 拖拽相關方法 ==========

/**
 * 開始拖拽規則
 */
const handleDragStart = (event: DragEvent, rule: OccupancyRule) => {
  isDragging.value = true
  draggedRule.value = rule

  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'move'
    event.dataTransfer.setData('application/json', JSON.stringify({ id: rule.id }))
  }

  const target = event.target as HTMLElement
  if (target) {
    target.style.opacity = '0.5'
  }
}

/**
 * 結束拖拽
 */
const handleDragEnd = (event: DragEvent) => {
  isDragging.value = false
  draggedRule.value = null
  dragOverSlot.value = null

  const target = event.target as HTMLElement
  if (target) {
    target.style.opacity = ''
  }
}

/**
 * 拖拽經過規則上方
 */
const handleDragOver = (event: DragEvent, rule: OccupancyRule) => {
  if (!draggedRule.value || draggedRule.value.id === rule.id) return
  dragOverSlot.value = { weekday: rule.weekday, hour: parseInt(rule.start_time.split(':')[0]), ruleId: rule.id }
}

/**
 * 拖拽經過空閒時段
 */
const handleEmptySlotDragOver = (event: DragEvent, weekday: number, hour: number) => {
  if (!draggedRule.value) return
  dragOverSlot.value = { weekday, hour, ruleId: 0 }
}

/**
 * 離開拖拽區域
 */
const handleDragLeave = () => {
  dragOverSlot.value = null
}

/**
 * 檢查是否是有效的放置目標
 */
const isValidDropTarget = (weekday: number, hour: number): boolean => {
  if (!draggedRule.value) return false
  return true
}

/**
 * 點擊佔用時段（查看/編輯規則）- 直接傳入課程物件
 */
const handleSlotClick = (rule: OccupancyRule) => {
  navigateToEditRule(rule.id)
}

/**
 * 導航至編輯規則
 */
const navigateToEditRule = (ruleId: number) => {
  router.push({
    path: '/admin/schedules',
    query: {
      action: 'edit',
      rule_id: ruleId
    }
  })
}

/**
 * 放置到規則上
 */
const handleDrop = async (event: DragEvent, targetRule: OccupancyRule) => {
  if (!draggedRule.value) return
  const rule = draggedRule.value
  if (rule.id === targetRule.id) return

  await updateRuleTime(rule.id, targetRule.weekday, targetRule.start_time, targetRule.end_time)
  dragOverSlot.value = null
}

/**
 * 放置到空閒時段
 */
const handleDropToEmpty = async (event: DragEvent, weekday: number, hour: number) => {
  if (!draggedRule.value) return
  const rule = draggedRule.value

  const durationParts = getDuration(rule.start_time, rule.end_time)
  const startTime = `${hour.toString().padStart(2, '0')}:00`
  const endTime = addDuration(startTime, durationParts.h, durationParts.m)

  await updateRuleTime(rule.id, weekday, startTime, endTime)
  dragOverSlot.value = null
}

/**
 * 更新規則時間
 */
const updateRuleTime = async (ruleId: number, weekday: number, startTime: string, endTime: string) => {
  try {
    const hasConflictNow = checkConflictForRule(ruleId, weekday, startTime, endTime)
    if (hasConflictNow) {
      const confirmed = await alertConfirm('注意：此時段與其他課程衝突，是否確定要進行調整？')
      if (!confirmed) return
    }

    await api.put(`/admin/scheduling/rules/${ruleId}`, {
      weekday,
      start_time: startTime,
      end_time: endTime,
      update_mode: 'SINGLE'
    })

    await fetchOccupancyRules()
    success('規則已更新')
  } catch (err) {
    console.error('Failed to update rule time:', err)
    toastError('更新失敗')
  }
}

/**
 * 檢查衝突
 */
const checkConflictForRule = (ruleId: number, weekday: number, startTime: string, endTime: string): boolean => {
  const toMinutes = (timeStr: string) => {
    const [h, m] = timeStr.split(':').map(Number)
    return h * 60 + m
  }
  const startMin = toMinutes(startTime)
  const endMin = toMinutes(endTime)

  for (const other of occupancyRules.value) {
    if (other.id === ruleId) continue
    if (other.weekday !== weekday) continue
    const otherS = toMinutes(other.start_time)
    const otherE = toMinutes(other.end_time)
    if (startMin < otherE && endMin > otherS) return true
  }
  return false
}

const checkOverlap = (r1: OccupancyRule, r2: OccupancyRule) => {
  const toMin = (s: string) => s.split(':').map(Number)[0] * 60 + s.split(':').map(Number)[1]
  return toMin(r1.start_time) < toMin(r2.end_time) && toMin(r1.end_time) > toMin(r2.start_time)
}

const hasConflict = (rule: OccupancyRule) => {
  return occupancyRules.value.some(other => {
    if (other.id === rule.id || other.weekday !== rule.weekday) return false
    return checkOverlap(rule, other)
  })
}

const getDuration = (start: string, end: string) => {
  const s = start.split(':').map(Number)
  const e = end.split(':').map(Number)
  const total = (e[0] * 60 + e[1]) - (s[0] * 60 + s[1])
  return { h: Math.floor(total / 60), m: total % 60 }
}

const addDuration = (start: string, h: number, m: number) => {
  const s = start.split(':').map(Number)
  let eh = s[0] + h
  let em = s[1] + m
  if (em >= 60) {
    eh += Math.floor(em / 60)
    em = em % 60
  }
  return `${eh.toString().padStart(2, '0')}:${em.toString().padStart(2, '0')}`
}

/**
 * 打開複製精靈
 */
const openCopyRulesWizard = () => {
  showCopyWizard.value = true
}

/**
 * 處理複製完成
 */
const handleCopyComplete = () => {
  success('規則複製成功')
  showCopyWizard.value = false
  fetchOccupancyRules()
}

/**
 * 導航到新增規則
 */
const navigateToCreateRule = (day?: number, hour?: number) => {
  // 多選時使用第一個選中的資源
  const firstResource = selectedResource.value[0]
  if (!firstResource) return

  const query: any = {
    action: 'create',
    [resourceType.value === 'teacher' ? 'teacher' : 'room']: firstResource.id,
  }
  if (day) query.weekday = day
  if (hour) query.start_time = `${hour.toString().padStart(2, '0')}:00`
  router.push({ path: '/admin/schedules', query })
}

// ========== Lifecycle ==========

onMounted(async () => {
  await Promise.all([
    fetchTerms(),
    fetchAllResources(),
    fetchCenterSettings(),
  ])
})
</script>
