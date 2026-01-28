<template>
  <div class="h-full flex flex-col glass-card overflow-hidden">
    <div class="p-4 border-b border-white/10 shrink-0">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div class="flex items-center gap-4">
          <!-- 週導航區域 -->
          <div class="flex items-center gap-2">
            <button
              @click="changeWeek(-1)"
              class="p-2 rounded-lg hover:bg-white/10 transition-colors"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </button>

            <h2 class="text-lg font-semibold text-slate-100">
              {{ weekLabel }}
            </h2>

            <button
              @click="changeWeek(1)"
              class="p-2 rounded-lg hover:bg-white/10 transition-colors"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </button>

            <!-- 週導航說明 -->
            <HelpTooltip
              placement="bottom"
              title="週期導航"
              description="查看不同週期的排課狀況，預設顯示本週。"
              :usage="['點擊左右箭頭切換上週/下週', '可跨月、跨年查看', '所有視角共用同一週期']"
            />
          </div>

          <!-- 視角切換器 -->
          <div class="flex items-center gap-1 bg-slate-800/80 rounded-lg p-1">
            <button
              @click="viewModeModel = 'calendar'"
              class="px-3 py-1.5 rounded-md text-sm font-medium transition-all"
              :class="viewMode === 'calendar' ? 'bg-primary-500 text-white' : 'text-slate-400 hover:text-white'"
            >
              週曆
            </button>
            <button
              @click="viewModeModel = 'teacher_matrix'"
              class="px-3 py-1.5 rounded-md text-sm font-medium transition-all"
              :class="viewMode === 'teacher_matrix' ? 'bg-primary-500 text-white' : 'text-slate-400 hover:text-white'"
            >
              老師矩陣
            </button>
            <button
              @click="viewModeModel = 'room_matrix'"
              class="px-3 py-1.5 rounded-md text-sm font-medium transition-all"
              :class="viewMode === 'room_matrix' ? 'bg-primary-500 text-white' : 'text-slate-400 hover:text-white'"
            >
              教室矩陣
            </button>
          </div>

          <!-- 矩陣視角選擇器 -->
          <div v-if="viewMode !== 'calendar'" class="flex items-center gap-2">
            <select
              v-model="selectedResourceIdModel"
              class="px-3 py-1.5 rounded-lg text-sm bg-slate-800/80 border border-white/10 text-slate-300 focus:outline-none focus:border-primary-500"
            >
              <option :value="null">選擇{{ viewMode === 'teacher_matrix' ? '老師' : '教室' }}...</option>
              <option v-for="resource in resourceList" :key="resource.id" :value="resource.id">
                {{ resource.name }}
              </option>
            </select>
            <HelpTooltip
              :title="viewMode === 'teacher_matrix' ? '選擇老師' : '選擇教室'"
              :description="`從已加入中心的${viewMode === 'teacher_matrix' ? '老師' : '教室'}中選擇，查看該${viewMode === 'teacher_matrix' ? '老師' : '教室'}的專屬排課表。`"
              :usage="['從下拉選單選擇特定人員', '選擇後畫面會顯示該人員的排課', '可點擊右上角 X 清除篩選']"
            />
          </div>

          <!-- 新增排課按鈕 -->
          <div class="flex items-center gap-2 ml-auto">
            <button
              @click="showCreateModal = true"
              class="btn-primary px-4 py-2 text-sm font-medium"
            >
              + 新增排課規則
            </button>
            <HelpTooltip
              title="新增排課規則"
              description="建立新的課程排課規則，設定課程、老師、教室、時間等資訊。"
              :usage="['點擊按鈕開啟新增表單', '選擇課程、老師、教室', '設定每週固定上課日與時段', '設定有效期限後儲存']"
              shortcut="Ctrl + N"
            />
          </div>
        </div>
      </div>

      <!-- 選中資源提示 -->
      <div
        v-if="viewMode !== 'calendar' && selectedResourceName"
        class="mt-3 flex items-center gap-2 px-3 py-2 bg-primary-500/10 border border-primary-500/30 rounded-lg"
      >
        <span class="text-sm text-primary-400">
          {{ viewMode === 'teacher_matrix' ? '老師' : '教室' }}矩陣：
        </span>
        <span class="text-sm font-medium text-white">{{ selectedResourceName }}</span>
        <button
          @click="clearViewMode"
          class="ml-auto p-1 hover:bg-white/10 rounded transition-colors"
        >
          <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <div
      class="flex-1 overflow-auto p-4"
      @dragover.prevent="handleDragOver"
      @drop="handleDrop"
    >
      <!-- 週曆視圖 -->
      <div v-if="viewMode === 'calendar'" class="min-w-[600px] relative" ref="calendarContainerRef">
        <!-- 表頭 -->
        <div class="grid sticky top-0 z-10 bg-slate-800/90 backdrop-blur-sm" style="grid-template-columns: 80px repeat(7, 1fr);">
          <div class="p-2 border-b border-white/10 text-center">
            <span class="text-xs text-slate-400">時段</span>
          </div>
          <div
            v-for="day in weekDays"
            :key="day.value"
            class="p-2 border-b border-white/10 text-center"
          >
            <span class="text-sm font-medium text-slate-100">{{ day.name }}</span>
          </div>
        </div>

        <!-- 時間列和網格 -->
        <div
          v-for="time in timeSlots"
          :key="time"
          class="grid relative"
          style="grid-template-columns: 80px repeat(7, 1fr);"
        >
          <!-- 時間標籤 -->
          <div class="p-2 border-r border-b border-white/5 text-right text-xs text-slate-400">
            {{ formatTime(time) }}
          </div>

          <!-- 每日網格 -->
          <div
            v-for="day in weekDays"
            :key="`${time}-${day.value}`"
            class="p-0 min-h-[60px] border-b border-white/5 border-r relative"
            :class="getCellClass(time, day.value)"
            @dragenter="handleDragEnter(time, day.value)"
            @dragleave="handleDragLeave"
            @dragover.prevent
          />
        </div>

        <!-- 課程卡片層 - 絕對定位 -->
        <div class="absolute top-0 left-[80px] right-0 bottom-0 pointer-events-none">
          <div
            v-for="schedule in displaySchedules"
            :key="schedule.key"
            class="absolute rounded-lg p-2 text-xs cursor-pointer hover:opacity-90 transition-opacity pointer-events-auto"
            :class="getScheduleCardClass(schedule)"
            :style="getScheduleStyle(schedule)"
            @click="selectSchedule(schedule)"
          >
            <div class="font-medium truncate">
              {{ schedule.offering_name }}
            </div>
            <div class="text-slate-400 truncate">
              {{ schedule.teacher_name }}
            </div>
            <div class="text-slate-500 text-[10px] mt-0.5">
              {{ schedule.start_time }} - {{ schedule.end_time }}
            </div>
          </div>
        </div>
      </div>

      <!-- 矩陣視圖（老師/教室） -->
      <div v-else class="min-w-[800px] relative" ref="matrixContainerRef">
        <div class="grid sticky top-0 z-10 bg-slate-800/90 backdrop-blur-sm" style="grid-template-columns: 200px repeat(7, 1fr);">
          <div class="p-3 border-b border-white/10">
            <span class="text-sm font-medium text-slate-400">
              {{ viewMode === 'teacher_matrix' ? '老師' : '教室' }}
            </span>
          </div>
          <div
            v-for="day in weekDays"
            :key="day.value"
            class="p-3 border-b border-white/10 text-center"
          >
            <span class="text-sm font-medium text-slate-100">{{ day.name }}</span>
          </div>
        </div>

        <!-- 資源列表 -->
        <div
          v-for="resource in matrixResources"
          :key="resource.id"
          class="grid relative"
          style="grid-template-columns: 200px repeat(7, 1fr);"
        >
          <!-- 資源名稱 -->
          <div class="p-3 border-b border-white/5 border-r flex items-center bg-slate-900/50 sticky left-0 z-10">
            <div class="flex items-center gap-2">
              <div
                class="w-8 h-8 rounded-full flex items-center justify-center text-xs font-medium"
                :class="viewMode === 'teacher_matrix' ? 'bg-primary-500/20 text-primary-400' : 'bg-amber-500/20 text-amber-400'"
              >
                {{ resource.name?.charAt(0) || '?' }}
              </div>
              <span class="text-sm text-slate-300">{{ resource.name }}</span>
            </div>
          </div>

          <!-- 每週的網格 -->
          <div
            v-for="day in weekDays"
            :key="`${resource.id}-${day.value}`"
            class="p-0 min-h-[80px] border-b border-white/5 border-r relative"
          />

          <!-- 矩陣課程卡片 -->
          <div class="absolute top-0 left-[200px] right-0 bottom-0 pointer-events-none">
            <div
              v-for="schedule in getMatrixSchedulesForResource(resource.id)"
              :key="schedule.key"
              class="absolute rounded-lg p-2 text-xs cursor-pointer hover:opacity-90 transition-opacity pointer-events-auto"
              :class="getScheduleCardClass(schedule)"
              :style="getMatrixScheduleStyle(schedule, resource.id)"
              @click="selectSchedule(schedule)"
            >
              <div class="font-medium truncate text-white">
                {{ schedule.offering_name }}
              </div>
              <div class="text-slate-400 truncate text-[10px]">
                {{ schedule.start_time }} - {{ schedule.end_time }}
              </div>
            </div>
          </div>
        </div>

        <!-- 空狀態 -->
        <div v-if="matrixResources.length === 0" class="text-center py-12">
          <div class="text-slate-500 mb-2">暫無{{ viewMode === 'teacher_matrix' ? '老師' : '教室' }}資料</div>
          <div class="text-xs text-slate-600">請先{{ viewMode === 'teacher_matrix' ? '新增老師' : '新增教室' }}</div>
        </div>
      </div>
    </div>

    <Teleport to="body">
      <ScheduleDetailPanel
        v-if="selectedCell"
        :time="selectedCell.time"
        :weekday="selectedCell.day"
        :schedule="selectedSchedule"
        @close="selectedCell = null"
        @edit="handleEdit"
        @delete="handleDelete"
      />
    </Teleport>

    <Teleport to="body">
      <UpdateModeModal
        v-if="showUpdateModeModal"
        :show="showUpdateModeModal"
        :rule-name="editingRule?.offering_name"
        :rule-date="editingRule?.date ? new Date(editingRule.date).toLocaleDateString('zh-TW', { year: 'numeric', month: 'long', day: 'numeric' }) : ''"
        @close="handleUpdateModeClose"
        @confirm="handleUpdateModeConfirm"
      />
    </Teleport>

    <Teleport to="body">
      <ScheduleRuleModal
        v-if="showCreateModal"
        @close="showCreateModal = false"
        @created="handleRuleCreated"
      />
      <ScheduleRuleModal
        v-if="showEditModal"
        :editing-rule="editingRule"
        :update-mode="pendingUpdateMode"
        @close="handleEditModalClose"
        @submit="handleRuleUpdated"
      />
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { formatDateToString } from '~/composables/useTaiwanTime'
import { nextTick, ref, computed, watch, onMounted, onUnmounted } from 'vue'

const emit = defineEmits<{
  selectCell: { time: number, weekday: number }
  'update:viewMode': [value: 'calendar' | 'teacher_matrix' | 'room_matrix']
  'update:selectedResourceId': [value: number | null]
}>()

// Alert composable
const { confirm: confirmDialog, error: alertError } = useAlert()

// Props
const props = defineProps<{
  viewMode: 'calendar' | 'teacher_matrix' | 'room_matrix'
  selectedResourceId: number | null
}>()

// Computed with setter for v-model support
const viewModeModel = computed({
  get: () => props.viewMode,
  set: (value) => emit('update:viewMode', value)
})

const selectedResourceIdModel = computed({
  get: () => props.selectedResourceId,
  set: (value) => emit('update:selectedResourceId', value)
})

// 使用共享的資源緩存
const { resourceCache, fetchAllResources } = useResourceCache()

// DOM 引用
const calendarContainerRef = ref<HTMLElement | null>(null)
const matrixContainerRef = ref<HTMLElement | null>(null)
const slotWidth = ref(100)

// 格子尺寸常量
const TIME_SLOT_HEIGHT = 60 // 每個時段格子的高度 (px)
const RESOURCE_COLUMN_WIDTH = 200 // 資源列寬度 (px)
const TIME_COLUMN_WIDTH = 80 // 時間列寬度 (px)

const showCreateModal = ref(false)
const showEditModal = ref(false)
const showUpdateModeModal = ref(false)
const editingRule = ref<any>(null)
const pendingUpdateMode = ref<string>('')
const selectedCell = ref<{ time: number, day: number } | null>(null)
const selectedSchedule = ref<any>(null)
const dragTarget = ref<{ time: number, day: number } | null>(null)
const validationResults = ref<Record<string, any>>({})

const handleEdit = () => {
  if (selectedSchedule.value) {
    editingRule.value = selectedSchedule.value
    showUpdateModeModal.value = true
  }
}

const handleUpdateModeClose = () => {
  showUpdateModeModal.value = false
  editingRule.value = null
}

const handleUpdateModeConfirm = (mode: string) => {
  pendingUpdateMode.value = mode
  showUpdateModeModal.value = false
  showEditModal.value = true
}

const handleEditModalClose = () => {
  showEditModal.value = false
  editingRule.value = null
  pendingUpdateMode.value = ''
}

const handleDelete = async () => {
  const confirmed = await confirmDialog('確定要刪除此排課規則？')
  if (!confirmed || !selectedSchedule.value) return

  try {
    const api = useApi()
    await api.delete(`/admin/rules/${selectedSchedule.value.id}`)
    selectedCell.value = null
    selectedSchedule.value = null
    await fetchSchedules()
  } catch (err) {
    console.error('Failed to delete rule:', err)
    await alertError('刪除失敗，請稍後再試')
  }
}

const handleRuleUpdated = async (formData: any, updateMode: string) => {
  try {
    const api = useApi()
    await api.put(`/admin/rules/${editingRule.value.id}`, {
      ...formData,
      update_mode: updateMode,
    })
    await fetchSchedules()
    selectedCell.value = null
    selectedSchedule.value = null
    editingRule.value = null
    pendingUpdateMode.value = ''
  } catch (err) {
    console.error('Failed to update rule:', err)
    await alertError('更新失敗，請稍後再試')
  }
}

// 資源列表
const resourceList = computed(() => {
  if (props.viewMode === 'teacher_matrix') {
    return Array.from(resourceCache.value.teachers.values())
  } else if (props.viewMode === 'room_matrix') {
    return Array.from(resourceCache.value.rooms.values())
  }
  return []
})

// 矩陣視圖的資源列表
const matrixResources = computed(() => {
  if (props.viewMode === 'teacher_matrix') {
    return Array.from(resourceCache.value.teachers.values())
  } else if (props.viewMode === 'room_matrix') {
    return Array.from(resourceCache.value.rooms.values())
  }
  return []
})

const getWeekStart = (date: Date): Date => {
  const d = new Date(date)
  const day = d.getDay()
  const diff = d.getDate() - day + (day === 0 ? -6 : 1)
  return new Date(d.setDate(diff))
}

const weekStart = ref(getWeekStart(new Date()))
const weekEnd = computed(() => {
  const end = new Date(weekStart.value)
  end.setDate(end.getDate() + 6)
  return end
})

const weekLabel = computed(() => {
  const start = weekStart.value.toLocaleDateString('zh-TW', { month: 'long', day: 'numeric' })
  const end = weekEnd.value.toLocaleDateString('zh-TW', { month: 'long', day: 'numeric', year: 'numeric' })
  return `${start} - ${end}`
})

// 時間段 - 包含 00:00-03:00 和 22:00-23:00
const timeSlots = [0, 1, 2, 3, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23]

const weekDays = [
  { value: 1, name: '週一' },
  { value: 2, name: '週二' },
  { value: 3, name: '週三' },
  { value: 4, name: '週四' },
  { value: 5, name: '週五' },
  { value: 6, name: '週六' },
  { value: 7, name: '週日' },
]

const schedules = ref<any[]>([])
const { getCenterId } = useCenterId()

const selectedResourceName = computed(() => {
  if (props.viewMode === 'teacher_matrix') {
    return resourceCache.value.teachers.get(props.selectedResourceId)?.name || '未知老師'
  } else if (props.viewMode === 'room_matrix') {
    return resourceCache.value.rooms.get(props.selectedResourceId)?.name || '未知教室'
  }
  return ''
})

const clearViewMode = () => {
  viewModeModel.value = 'calendar'
  selectedResourceIdModel.value = null
}

// 計算課程持續分鐘數
const calculateDurationMinutes = (startTime: string, endTime: string): number => {
  const [startHour, startMinute] = startTime.split(':').map(Number)
  const [endHour, endMinute] = endTime.split(':').map(Number)

  const startMinutes = startHour * 60 + startMinute
  let endMinutes = endHour * 60 + endMinute

  // 跨日處理
  if (endMinutes <= startMinutes) {
    endMinutes += 24 * 60
  }

  return endMinutes - startMinutes
}

// 計算格子寬度
const calculateSlotWidth = () => {
  if (props.viewMode === 'calendar' && calendarContainerRef.value) {
    const containerWidth = calendarContainerRef.value.offsetWidth
    slotWidth.value = Math.max(80, (containerWidth - TIME_COLUMN_WIDTH) / 7)
  } else if (props.viewMode !== 'calendar' && matrixContainerRef.value) {
    const containerWidth = matrixContainerRef.value.offsetWidth
    slotWidth.value = Math.max(80, (containerWidth - RESOURCE_COLUMN_WIDTH) / 7)
  }
}

const fetchSchedules = async () => {
  try {
    const api = useApi()

    // 取得當前週的日期範圍
    const startDate = formatDateToString(weekStart.value)
    const endDate = formatDateToString(weekEnd.value)

    // 使用 ExpandRules API，取得已展開並處理例外的排課
    const response = await api.post<{ code: number; datas: any[] }>('/admin/expand-rules', {
      rule_ids: [],
      start_date: startDate,
      end_date: endDate,
    })

    const expandedSchedules = response.datas || []

    // 將展開後的排課轉換為前端格式
    const scheduleList = expandedSchedules.map((schedule: any) => {
      const date = new Date(schedule.date)
      const weekday = date.getDay() === 0 ? 7 : date.getDay()
      const startTime = schedule.start_time || '09:00'
      const endTime = schedule.end_time || '10:00'
      const [startHour, startMinute] = startTime.split(':').map(Number)
      const durationMinutes = calculateDurationMinutes(startTime, endTime)

      return {
        id: schedule.rule_id,
        key: `${schedule.rule_id}-${weekday}-${startTime}-${schedule.date}`,
        offering_name: schedule.offering_name || '-',
        teacher_name: schedule.teacher_name || '-',
        teacher_id: schedule.teacher_id,
        room_id: schedule.room_id,
        room_name: schedule.room_name || '-',
        weekday: weekday,
        start_time: startTime,
        end_time: endTime,
        start_hour: startHour,
        start_minute: startMinute,
        duration_minutes: durationMinutes,
        date: schedule.date,
        has_exception: schedule.has_exception || false,
        exception_type: schedule.exception_type || null,
        exception_info: schedule.exception_info || null,
        rule: schedule.rule || null,
        offering_id: schedule.offering_id,
        effective_range: schedule.effective_range || null,
      }
    })

    schedules.value = scheduleList

    // 等待 DOM 更新後計算槽寬度
    await nextTick()
    calculateSlotWidth()
  } catch (error) {
    console.error('Failed to fetch schedules:', error)
    schedules.value = []
  }
}

// 去重後的排課
const displaySchedules = computed(() => {
  const seen = new Set<string>()
  const result: any[] = []

  for (const schedule of schedules.value) {
    const key = `${schedule.id}-${schedule.weekday}-${schedule.start_time}`
    if (!seen.has(key)) {
      seen.add(key)
      result.push(schedule)
    }
  }

  return result
})

// 根據視角模式過濾排課
const filteredSchedules = computed(() => {
  if (props.viewMode === 'calendar' || !props.selectedResourceId) {
    return displaySchedules.value
  }

  return displaySchedules.value.filter(schedule => {
    if (props.viewMode === 'teacher_matrix') {
      return schedule.teacher_id === props.selectedResourceId
    } else if (props.viewMode === 'room_matrix') {
      return schedule.room_id === props.selectedResourceId
    }
    return false
  })
})

// 矩陣視圖：取得特定資源的排課
const getMatrixSchedulesForResource = (resourceId: number) => {
  return filteredSchedules.value.filter(s => {
    if (props.viewMode === 'teacher_matrix') {
      return s.teacher_id === resourceId
    } else if (props.viewMode === 'room_matrix') {
      return s.room_id === resourceId
    }
    return false
  })
}

const changeWeek = (delta: number) => {
  weekStart.value = getWeekStart(new Date(weekStart.value.getTime() + delta * 7 * 24 * 60 * 60 * 1000))
}

// 監聽週變化
watch(weekStart, async () => {
  await fetchSchedules()
})

const formatTime = (hour: number): string => {
  return `${hour.toString().padStart(2, '0')}:00`
}

// 週曆視圖：計算課程卡片樣式
const getScheduleStyle = (schedule: any) => {
  const { weekday, start_hour, start_minute, duration_minutes } = schedule

  // 計算水平位置
  const dayIndex = weekday - 1 // 0-6
  const left = dayIndex * slotWidth.value

  // 計算垂直位置
  // 計算 start_hour 前面有多少個時段格子
  let topSlotIndex = 0
  for (let t = 0; t < start_hour; t++) {
    if (t >= 0 && t <= 3) {
      topSlotIndex++ // 0-3 時段每個都算
    } else if (t >= 9) {
      topSlotIndex++ // 9 以後的時段每個都算
    }
  }

  // 每個格子高度 60px
  const slotHeight = TIME_SLOT_HEIGHT
  const baseTop = topSlotIndex * slotHeight

  // 加上分鐘偏移
  const minuteOffset = (start_minute / 60) * slotHeight
  const top = baseTop + minuteOffset

  // 計算高度（持續分鐘數轉像素）
  const height = (duration_minutes / 60) * slotHeight

  // 計算寬度（略小於格子寬度以留邊距）
  const width = slotWidth.value - 4

  return {
    left: `${left}px`,
    top: `${top}px`,
    width: `${width}px`,
    height: `${height}px`,
  }
}

// 矩陣視圖：計算課程卡片樣式
const getMatrixScheduleStyle = (schedule: any, resourceId: number) => {
  const { weekday, start_hour, start_minute, duration_minutes } = schedule

  // 計算水平位置
  const dayIndex = weekday - 1 // 0-6
  const left = dayIndex * slotWidth.value

  // 計算垂直位置
  let topSlotIndex = 0
  for (let t = 0; t < start_hour; t++) {
    if (t >= 0 && t <= 3) {
      topSlotIndex++
    } else if (t >= 9) {
      topSlotIndex++
    }
  }

  const slotHeight = TIME_SLOT_HEIGHT
  const baseTop = topSlotIndex * slotHeight
  const minuteOffset = (start_minute / 60) * slotHeight
  const top = baseTop + minuteOffset

  const height = (duration_minutes / 60) * slotHeight
  const width = slotWidth.value - 4

  return {
    left: `${left}px`,
    top: `${top}px`,
    width: `${width}px`,
    height: `${height}px`,
  }
}

const getCellClass = (time: number, weekday: number): string => {
  const key = `${time}-${weekday}`
  const validation = validationResults.value[key]

  if (validation?.valid === false) {
    return 'bg-critical-500/10 border-critical-500/50'
  } else if (validation?.warning) {
    return 'bg-warning-500/10 border-warning-500/50'
  } else if (validation?.valid === true) {
    return 'bg-success-500/10 border-success-500/50'
  }

  return 'hover:bg-white/5'
}

const getScheduleCardClass = (schedule: any): string => {
  if (!schedule) return ''

  if (schedule.has_exception) {
    switch (schedule.exception_type) {
      case 'CANCEL':
        return 'bg-critical-500/30 border border-critical-500/50 line-through'
      case 'RESCHEDULE':
        return 'bg-warning-500/30 border border-warning-500/50'
      case 'SWAP':
        return 'bg-primary-500/30 border border-primary-500/50'
      default:
        return 'bg-slate-700/80 border border-white/10'
    }
  }

  return 'bg-slate-700/80 border border-white/10'
}

const selectSchedule = (schedule: any) => {
  selectedCell.value = { time: schedule.start_hour, day: schedule.weekday }
  selectedSchedule.value = schedule
  emit('selectCell', { time: schedule.start_hour, weekday: schedule.weekday })
}

const handleDragOver = (event: DragEvent) => {
  if (dragTarget.value) {
    const key = `${dragTarget.value.time}-${dragTarget.value.day}`
    validationResults.value[key] = { valid: true }
  }
}

const handleDragEnter = (time: number, day: number) => {
  dragTarget.value = { time, day }
}

const handleDragLeave = () => {
  if (dragTarget.value) {
    const key = `${dragTarget.value.time}-${dragTarget.value.day}`
    delete validationResults.value[key]
  }
  dragTarget.value = null
}

const handleDrop = async (event: DragEvent) => {
  event.preventDefault()

  if (!dragTarget.value) return

  const data = event.dataTransfer?.getData('application/json')
  if (!data) return

  const parsed = JSON.parse(data)
  const { type, item } = parsed

  const key = `${dragTarget.value.time}-${dragTarget.value.day}`
  validationResults.value[key] = { valid: 'checking' }

  try {
    const api = useApi()
    const teacherId = type === 'teacher' ? item.id : (item.teacher_id || null)
    const roomId = type === 'room' ? item.id : (item.room_id || null)

    const response = await api.post<any>('/admin/scheduling/check-overlap', {
      teacher_id: teacherId,
      room_id: roomId,
      start_time: `${formatDate(weekStart.value)}T${formatTime(dragTarget.value.time)}:00`,
      end_time: `${formatDate(weekStart.value)}T${formatTime(dragTarget.value.time + 1)}:00`,
    })

    if (response.data.valid) {
      validationResults.value[key] = { valid: true }
    } else {
      validationResults.value[key] = { valid: false, conflicts: response.data.conflicts }
    }
  } catch (error) {
    console.error('Validation failed:', error)
    validationResults.value[key] = { valid: false, error: true }
  }

  dragTarget.value = null
}

const formatDate = (date: Date): string => {
  return formatDateToString(date)
}

const handleRuleCreated = () => {
  fetchSchedules()
}

onMounted(async () => {
  await fetchSchedules()
  fetchAllResources()

  // 等待 DOM 更新後計算槽寬度
  await nextTick()
  calculateSlotWidth()

  // 監控容器大小變化
  if (calendarContainerRef.value || matrixContainerRef.value) {
    const resizeObserver = new ResizeObserver(() => {
      calculateSlotWidth()
    })

    if (calendarContainerRef.value) {
      resizeObserver.observe(calendarContainerRef.value)
    }
    if (matrixContainerRef.value) {
      resizeObserver.observe(matrixContainerRef.value)
    }

    onUnmounted(() => {
      resizeObserver.disconnect()
    })
  }
})
</script>
