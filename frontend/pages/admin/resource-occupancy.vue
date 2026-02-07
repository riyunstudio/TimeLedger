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

    <!-- 篩選欄 -->
    <div class="glass-card p-4 mb-6">
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
        <div class="flex-1">
          <SearchableSelect
            v-model="selectedResourceId"
            :options="resourceOptions"
            :placeholder="`選擇或搜尋${resourceType === 'teacher' ? '老師' : '教室'}...`"
            :label="`搜尋${resourceType === 'teacher' ? '老師' : '教室'}`"
          />
        </div>
      </div>
    </div>

    <!-- 每週網格視圖 -->
    <div v-if="selectedResource && occupancyRules.length > 0" class="glass-card p-4">
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
        <!-- 時間標籤 -->
        <div class="space-y-2">
          <div
            v-for="hour in timeSlots"
            :key="hour"
            class="h-16 flex items-center justify-center text-xs md:text-sm text-slate-500"
          >
            {{ formatTimeLabel(hour) }}
          </div>
        </div>

        <!-- 每天的時段 -->
        <div
          v-for="day in weekDays"
          :key="day.value"
          class="space-y-2"
        >
          <div
            v-for="hour in timeSlots"
            :key="`${day.value}-${hour}`"
            class="h-16 rounded-lg border relative overflow-hidden transition-all"
            :class="getSlotClass(day.value, hour)"
          >
            <!-- 已佔用時段 -->
            <template v-if="getRulesForSlot(day.value, hour).length > 0">
              <div
                v-for="rule in getRulesForSlot(day.value, hour)"
                :key="rule.id"
                class="absolute left-0 right-0 px-2 py-1 text-xs cursor-move hover:opacity-80 transition-opacity overflow-hidden"
                :class="getRuleClass(rule)"
                :style="getRuleStyle(rule, hour)"
                draggable="true"
                @click="handleSlotClick(day.value, hour)"
                @dragstart="handleDragStart($event, rule)"
                @dragend="handleDragEnd($event)"
                @dragover.prevent="handleDragOver($event, day.value, hour)"
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
                <div class="font-medium truncate">
                  {{ rule.offering_name || rule.course_name }}
                </div>
                <div class="text-[10px] opacity-75 truncate">
                  {{ resourceType === 'teacher' ? rule.room_name : rule.teacher_name }}
                </div>
                <div class="text-[10px] opacity-75">
                  {{ rule.start_time }} - {{ rule.end_time }}
                </div>
              </div>
            </template>

            <!-- 空閒時段提示 -->
            <template v-else>
              <div
                v-if="isWithinTerm(day.value, hour)"
                class="w-full h-full flex items-center justify-center transition-colors"
                :class="isValidDropTarget(day.value, hour)
                  ? 'text-primary-400 bg-primary-500/10 cursor-pointer'
                  : 'text-slate-600 hover:text-slate-400 cursor-pointer'"
                @click="handleSlotClick(day.value, hour)"
                @dragover.prevent="handleEmptySlotDragOver($event, day.value, hour)"
                @dragleave="handleDragLeave"
                @drop="handleDropToEmpty($event, day.value, hour)"
              >
                <!-- 拖拽目標指示器 -->
                <div
                  v-if="isDragging && isValidDropTarget(day.value, hour)"
                  class="absolute inset-0 border-2 border-dashed border-primary-400 rounded bg-primary-500/20 flex items-center justify-center"
                >
                  <svg class="w-6 h-6 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                  </svg>
                </div>
                <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                </svg>
              </div>
            </template>
          </div>
        </div>
      </div>
    </div>

    <!-- 佔用率統計 -->
    <div v-if="selectedResource && occupancyRules.length > 0" class="mt-6 grid grid-cols-2 md:grid-cols-4 gap-4">
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

    <!-- 無數據提示 -->
    <div v-else-if="selectedResource && !loadingRules" class="glass-card p-12 text-center">
      <div class="w-16 h-16 bg-slate-800 rounded-full flex items-center justify-center mx-auto mb-4 border border-white/5">
        <svg class="w-8 h-8 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
      <h3 class="text-lg font-medium text-slate-300 mb-2">尚無排課記錄</h3>
      <p class="text-slate-500 max-w-sm mx-auto mb-6">
        目前 {{ resourceType === 'teacher' ? '這位老師' : '這間教室' }} 在此週段內沒有任何排課安排。
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
const { success, error: toastError } = useToast()
const { warning: alertWarning, confirm: alertConfirm } = useAlert()

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

// 選中的資源 ID
const selectedResourceId = ref<number | null>(null)
const selectedResource = ref<Resource | null>(null)

// 佔用規則
const loadingRules = ref(false)

// 定義後端返回的分組結構
interface GroupedOccupancy {
  day_of_week: number
  day_name: string
  rules: OccupancyRule[]
}

const occupancyRules = ref<OccupancyRule[]>([])

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

// 時間槽 (7:00 - 22:00)
const timeSlots = Array.from({ length: 16 }, (_, i) => i + 7)

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
  // 簡易計算：已佔用小時數 / 總可用小時數 (5天 * 12小時)
  const totalHours = 5 * 12
  let occupiedHours = 0

  occupancyRules.value.forEach(r => {
    const start = parseInt(r.start_time.split(':')[0])
    const end = parseInt(r.end_time.split(':')[0])
    occupiedHours += (end - start)
  })

  return Math.min(100, Math.round((occupiedHours / totalHours) * 100))
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
watch(selectedResourceId, (newId) => {
  if (newId !== null) {
    // 從快取中取得對應的物件
    if (resourceType.value === 'teacher') {
      const teacher = resourceCache.value.teachers.get(newId)
      selectedResource.value = teacher ? {
        id: teacher.id,
        name: teacher.name || `老師 ${teacher.id}`,
        count: teacher.count
      } : null
    } else {
      const room = resourceCache.value.rooms.get(newId)
      selectedResource.value = room ? {
        id: room.id,
        name: room.name || `教室 ${room.id}`,
        capacity: room.capacity
      } : null
    }
    // 立即取得佔用規則
    fetchOccupancyRules()
  } else {
    selectedResource.value = null
    occupancyRules.value = []
  }
})

/**
 * 監聽資源類型切換，清除選中的資源
 */
watch(resourceType, () => {
  selectedResourceId.value = null
  selectedResource.value = null
  occupancyRules.value = []
})

/**
 * 取得佔用規則
 */
const fetchOccupancyRules = async () => {
  if (!selectedResource.value) return

  loadingRules.value = true
  try {
    const response = await api.get<GroupedOccupancy[]>('/admin/occupancy/rules', {
      start_date: queryStartDate.value,
      end_date: queryEndDate.value,
      [resourceType.value === 'teacher' ? 'teacher_id' : 'room_id']: selectedResource.value.id,
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
 * 獲取指定時段的規則
 */
const getRulesForSlot = (day: number, hour: number): OccupancyRule[] => {
  return occupancyRules.value.filter(rule => {
    if (rule.weekday !== day) return false
    const startHour = parseInt(rule.start_time.split(':')[0])
    return startHour === hour
  })
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
 * 獲取規則卡片精確樣式 (處理分鐘位移和高度)
 */
const getRuleStyle = (rule: OccupancyRule, hour: number) => {
  const startParts = rule.start_time.split(':').map(Number)
  const endParts = rule.end_time.split(':').map(Number)
  
  const startMinute = startParts[1]
  const duration = (endParts[0] * 60 + endParts[1]) - (startParts[0] * 60 + startParts[1])
  
  // 每一小時高度 64px (h-16), 間距 8px (space-y-2)
  // 簡化計算：只針對當前小時內的起始進行位移
  const topOffset = (startMinute / 60) * 64
  const height = (duration / 60) * 64
  
  return {
    top: `${topOffset}px`,
    height: `${height}px`,
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
const handleDragOver = (event: DragEvent, weekday: number, hour: number) => {
  if (!draggedRule.value) return

  const rules = getRulesForSlot(weekday, hour)
  const targetRule = rules.find(r => r.id !== draggedRule.value?.id)

  if (targetRule) {
    dragOverSlot.value = {
      weekday,
      hour,
      ruleId: targetRule.id,
    }
  }
}

/**
 * 拖拽經過空閒時段
 */
const handleEmptySlotDragOver = (event: DragEvent, weekday: number, hour: number) => {
  if (!draggedRule.value) return

  if (isValidDropTarget(weekday, hour)) {
    dragOverSlot.value = { weekday, hour, ruleId: 0 }
  }
}

/**
 * 離開拖拽區域
 */
const handleDragLeave = () => {
  // 可以在這裡處理清除狀態
}

/**
 * 檢查是否是有效的放置目標
 */
const isValidDropTarget = (weekday: number, hour: number): boolean => {
  if (!draggedRule.value) return false
  return true
}

/**
 * 點擊佔用時段（查看/編輯規則）
 */
const handleSlotClick = (day: number, hour: number) => {
  const rules = getRulesForSlot(day, hour)
  if (rules.length > 0) {
    navigateToEditRule(rules[0].id)
  } else {
    navigateToCreateRule(day, hour)
  }
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
  const query: any = {
    action: 'create',
    [resourceType.value === 'teacher' ? 'teacher' : 'room']: selectedResource.value?.id,
  }
  if (day) query.weekday = day
  if (hour) query.start_time = `${hour.toString().padStart(2, '0')}:00`
  router.push({ path: '/admin/schedules', query })
}

// ========== Lifecycle ==========

onMounted(() => {
  fetchTerms()
  fetchAllResources()
})
</script>
