<template>
  <div class="h-full flex flex-col glass-card overflow-hidden">
    <div class="p-4 border-b border-white/10 shrink-0">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div class="flex items-center gap-4">
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
        </div>

        <div class="flex items-center gap-4">
          <!-- 資源類型切換 -->
          <div class="flex items-center gap-1 bg-slate-800/80 rounded-lg p-1">
            <button
              @click="$emit('update:resourceType', 'teacher')"
              class="px-3 py-1.5 rounded-md text-sm font-medium transition-all"
              :class="props.resourceType === 'teacher' ? 'bg-primary-500 text-white' : 'text-slate-400 hover:text-white'"
            >
              老師
            </button>
            <button
              @click="$emit('update:resourceType', 'room')"
              class="px-3 py-1.5 rounded-md text-sm font-medium transition-all"
              :class="props.resourceType === 'room' ? 'bg-primary-500 text-white' : 'text-slate-400 hover:text-white'"
            >
              教室
            </button>
          </div>

          <button
            @click="showCreateModal = true"
            class="btn-primary px-4 py-2 text-sm font-medium"
          >
            + 新增排課規則
          </button>
        </div>
      </div>
    </div>

    <div
      class="flex-1 overflow-auto relative"
      @dragover.prevent="handleDragOver"
      @drop="handleDrop"
    >
      <!-- 調試資訊 -->
      <div class="fixed bottom-4 right-4 bg-slate-800 p-2 rounded text-xs text-slate-300 z-50">
        <div>slotWidth: {{ slotWidth }}px</div>
        <div>container: {{ tableContainerRef?.offsetWidth || 'N/A' }}px</div>
      </div>

      <div class="min-w-[800px]" ref="tableContainerRef">
        <!-- 表頭：時段 -->
        <div class="grid" :style="gridStyle">
          <div class="p-3 border-b border-white/10 text-center bg-slate-800/50 sticky top-0 z-20">
            <span class="text-sm font-medium text-slate-300">資源 / 時段</span>
          </div>
          <div
            v-for="time in timeSlots"
            :key="time"
            class="p-3 border-b border-white/10 text-center bg-slate-800/50 sticky top-0 z-20"
          >
            <span class="text-sm font-medium text-slate-300">{{ formatTime(time) }}</span>
          </div>
        </div>

        <!-- 資源列 -->
        <div
          v-for="resource in resourceList"
          :key="resource.id"
          class="grid relative"
          :style="gridStyle"
        >
          <!-- 資源名稱 -->
          <div class="p-3 border-r border-b border-white/10 flex items-center bg-slate-900/50 sticky left-0 z-10">
            <div class="w-8 h-8 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center shrink-0 mr-2">
              <span class="text-white text-sm font-medium">{{ resource.name?.charAt(0) || '?' }}</span>
            </div>
            <span class="text-sm text-slate-300 truncate">{{ resource.name }}</span>
          </div>

          <!-- 時段網格背景 -->
          <div
            v-for="time in timeSlots"
            :key="`${resource.id}-bg-${time}`"
            class="p-1 min-h-[80px] border-b border-white/5 border-r relative"
            :class="getCellClass(resource.id, time)"
            @dragenter="handleDragEnter(resource.id, time)"
            @dragleave="handleDragLeave"
            @dragover.prevent
          />

          <!-- 課程卡片層 - 相對於資源列定位 -->
          <div
            v-for="schedule in getSchedulesForResource(resource.id)"
            :key="`${schedule.rule_id}-${schedule.weekday}-${schedule.start_hour}-${schedule.start_minute}`"
            class="absolute rounded-lg p-2 text-xs cursor-pointer hover:opacity-90 transition-opacity pointer-events-auto"
            :class="getScheduleCardClass(schedule)"
            :style="getScheduleStyle(schedule)"
            @click="selectSchedule(resource, schedule)"
          >
            <!-- 簡短資訊 -->
            <div class="font-medium truncate text-slate-100">
              {{ schedule.offering_name }}
            </div>
            <div class="text-slate-400 truncate">
              {{ schedule.teacher_name || schedule.room_name }}
            </div>
            <div class="text-slate-500 text-[10px] mt-0.5">
              {{ schedule.start_time }} - {{ schedule.end_time }}
            </div>
          </div>
        </div>

        <!-- 無資源時顯示 -->
        <div v-if="resourceList.length === 0" class="text-center py-12 text-slate-500">
          沒有{{ props.resourceType === 'teacher' ? '老師' : '教室' }}資料
        </div>
      </div>
    </div>

    <Teleport to="body">
      <ScheduleDetailPanel
        v-if="selectedCell"
        :time="selectedCell.time"
        :weekday="selectedCell.weekday"
        :schedule="selectedSchedule"
        @close="selectedCell = null"
        @edit="handleEdit"
        @delete="handleDelete"
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
        @close="showEditModal = false; editingRule = null"
        @updated="handleRuleUpdated"
      />
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { formatDateToString, getTodayString } from '~/composables/useTaiwanTime'
import { nextTick } from 'vue'

const emit = defineEmits<{
  selectCell: { resource: any, time: number, weekday: number }
  'update:resourceType': [value: 'teacher' | 'room']
}>()

// Alert composable
const { confirm: confirmDialog, error: alertError } = useAlert()

// Props
const props = defineProps<{
  resourceType: 'teacher' | 'room'
}>()

const showCreateModal = ref(false)
const showEditModal = ref(false)
const editingRule = ref<any>(null)
const selectedCell = ref<{ resource: any, time: number, weekday: number } | null>(null)
const selectedSchedule = ref<any>(null)
const dragTarget = ref<{ resourceId: number, time: number } | null>(null)
const validationResults = ref<Record<string, any>>({})
const tableContainerRef = ref<HTMLElement | null>(null)
const slotWidth = ref(100) // 每個時段槽的像素寬度

const handleEdit = () => {
  if (selectedSchedule.value) {
    editingRule.value = selectedSchedule.value
    showEditModal.value = true
  }
}

const handleDelete = async () => {
  if (!selectedSchedule.value || !await confirmDialog('確定要刪除此排課規則？')) return

  try {
    const api = useApi()
    await api.delete(`/admin/rules/${selectedSchedule.value.id}`)
    selectedCell.value = null
    selectedSchedule.value = null
    await fetchData()
  } catch (err) {
    console.error('Failed to delete rule:', err)
    await alertError('刪除失敗，請稍後再試')
  }
}

const handleRuleUpdated = async () => {
  await fetchData()
  selectedCell.value = null
  selectedSchedule.value = null
}

const { getCenterId } = useCenterId()

// 擴展時間段：包含 00:00-03:00 和 22:00-23:00
const timeSlots = [0, 1, 2, 3, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23]
const RESOURCE_COLUMN_WIDTH = 120 // 資源名稱列的寬度 (px)
const SLOT_COUNT = 19 // 時段數量

const weekDays = [
  { value: 1, name: '週一' },
  { value: 2, name: '週二' },
  { value: 3, name: '週三' },
  { value: 4, name: '週四' },
  { value: 5, name: '週五' },
  { value: 6, name: '週六' },
  { value: 7, name: '週日' },
]

// 課程資料結構
interface ScheduleData {
  id: number
  rule_id: number
  offering_name: string
  teacher_name: string
  teacher_id: number
  room_id: number
  room_name: string
  weekday: number
  start_time: string
  end_time: string
  start_hour: number
  start_minute: number
  duration: number
  is_cross_day?: boolean
  is_split_entry?: boolean // 標記是否為跨日分割後的條目
}

const schedules = ref<{ teachers: ScheduleData[], rooms: ScheduleData[] }>({
  teachers: [],
  rooms: []
})
const teachers = ref<any[]>([])
const rooms = ref<any[]>([])

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

const resourceList = computed(() => {
  return props.resourceType === 'teacher' ? teachers.value : rooms.value
})

const formatTime = (hour: number): string => {
  return `${hour.toString().padStart(2, '0')}:00`
}

// Grid 樣式
const gridStyle = computed(() => ({
  gridTemplateColumns: `${RESOURCE_COLUMN_WIDTH}px repeat(${SLOT_COUNT}, 1fr)`
}))

// 解析時間字串為小時和分鐘
const parseTime = (timeStr: string): { hour: number, minute: number } => {
  // 處理 "HH:mm:ss" 或 "HH:mm" 格式
  const parts = timeStr.split(':')
  const hour = parseInt(parts[0], 10)
  const minute = parseInt(parts[1], 10)
  return { hour, minute }
}

// 計算課程持續小時數
const calculateDuration = (startTime: string, endTime: string): number => {
  const start = parseTime(startTime)
  const end = parseTime(endTime)

  // 計算分鐘差
  const startMinutes = start.hour * 60 + start.minute
  const endMinutes = end.hour * 60 + end.minute

  let durationMinutes = endMinutes - startMinutes

  // 跨日處理：如果結束時間早於開始時間，表示跨日
  if (durationMinutes <= 0) {
    durationMinutes += 24 * 60
  }

  // 返回小時數（向上取整到最接近的整數，確保顯示完整）
  return Math.max(1, Math.ceil(durationMinutes / 60))
}

const fetchData = async () => {
  try {
    const api = useApi()

    const [rulesRes, teachersRes, roomsRes] = await Promise.all([
      api.get<{ code: number; datas: any[] }>('/admin/rules'),
      api.get<{ code: number; datas: any[] }>('/teachers'),
      api.get<{ code: number; datas: any[] }>(`/admin/rooms`)
    ])

    // 處理老師資料
    teachers.value = teachersRes.datas || []

    // 處理教室資料
    rooms.value = roomsRes.datas || []

    // 將規則轉換為 schedule 陣列
    const teacherSchedules: ScheduleData[] = []
    const roomSchedules: ScheduleData[] = []
    const rules = rulesRes.datas || []

    rules.forEach((rule: any) => {
      const day = rule.weekday
      if (!day) return

      const { hour: startHour, minute: startMinute } = parseTime(rule.start_time)
      const duration = calculateDuration(rule.start_time, rule.end_time)

      // 判斷是否跨日
      const { hour: endHour } = parseTime(rule.end_time)
      const isCrossDay = endHour < startHour || (endHour === startHour && rule.end_time > rule.start_time)

      // 判斷是否為跨日分割條目
      const originalRuleId = rule.rule_id || rule.id
      const isSplitEntry = rule.rule_id !== undefined && rule.rule_id !== rule.id

      const scheduleData: ScheduleData = {
        id: rule.id,
        rule_id: originalRuleId,
        offering_name: rule.offering?.name || '-',
        teacher_name: rule.teacher?.name || '-',
        teacher_id: rule.teacher_id,
        room_id: rule.room_id,
        room_name: rule.room?.name || '-',
        weekday: day,
        start_time: rule.start_time,
        end_time: rule.end_time,
        start_hour: startHour,
        start_minute: startMinute,
        duration: duration,
        is_cross_day: isCrossDay,
        is_split_entry: isSplitEntry,
      }

      // 根據資源類型分別存儲
      if (rule.teacher_id) {
        teacherSchedules.push(scheduleData)
      }
      if (rule.room_id) {
        roomSchedules.push(scheduleData)
      }
    })

    schedules.value = {
      teachers: teacherSchedules,
      rooms: roomSchedules
    }
  } catch (error) {
    console.error('Failed to fetch data:', error)
  }
}

// 獲取資源的課程列表（排除跨日分割的條目，只顯示原始條目）
const getSchedulesForResource = (resourceId: number): ScheduleData[] => {
  const scheduleList = props.resourceType === 'teacher' ? schedules.value.teachers : schedules.value.rooms

  // 過濾出非分割條目，或者雖然是分割條目但沒有找到原始條目的
  const filtered = scheduleList.filter(s => {
    // 檢查是否屬於該資源
    const isResourceMatch = s.teacher_id === resourceId || s.room_id === resourceId
    if (!isResourceMatch) return false

    // 如果是分割條目，檢查是否有對應的原始條目
    if (s.is_split_entry) {
      const hasOriginal = scheduleList.some(other =>
        other.rule_id === s.rule_id &&
        !other.is_split_entry &&
        other.id !== s.id
      )
      // 如果有原始條目，不顯示這個分割條目
      if (hasOriginal) return false
    }

    return true
  })

  return filtered
}

// 計算課程卡片的樣式
const getScheduleStyle = (schedule: ScheduleData): Record<string, string> => {
  const slotIndex = timeSlots.indexOf(schedule.start_hour)
  if (slotIndex === -1) return { display: 'none' }

  const { start_minute, duration, start_hour } = schedule

  // 使用計算好的 slotWidth
  const currentSlotWidth = slotWidth.value || 100

  // 左邊位置：資源列寬度 + slotIndex * 每格寬度 + 分鐘偏移
  const left = RESOURCE_COLUMN_WIDTH + slotIndex * currentSlotWidth + (start_minute / 60) * currentSlotWidth

  // 寬度：持續時間 * 每格寬度
  const width = duration * currentSlotWidth

  // 頂部位置（基於分鐘偏移，相對於儲存格高度）
  const top = (start_minute / 60) * 100

  // 高度：持續時間對應的時段數量百分比
  const height = duration * 100

  return {
    left: `${left}px`,
    top: `${top}%`,
    width: `${width}px`,
    height: `${height}%`,
  }
}

const getCellClass = (resourceId: number, time: number): string => {
  const key = `${resourceId}-${time}`
  const validation = validationResults.value[key]

  if (validation?.valid === false) {
    return 'bg-critical-500/10 border-critical-500/50'
  } else if (validation?.warning) {
    return 'bg-warning-500/10 border-warning-500/50'
  } else if (validation?.valid === true) {
    return 'bg-success-500/10 border-success-500/50'
  }

  return ''
}

const getScheduleCardClass = (schedule: ScheduleData): string => {
  if (!schedule) return ''
  // 跨日課程使用不同的樣式
  if (schedule.is_cross_day) {
    return 'bg-gradient-to-r from-amber-500/30 to-orange-500/30 border border-amber-500/50'
  }
  return 'bg-slate-700/80 border border-white/10'
}

const selectSchedule = (resource: any, schedule: ScheduleData) => {
  selectedCell.value = {
    resource,
    time: schedule.start_hour,
    weekday: schedule.weekday
  }
  selectedSchedule.value = schedule
  emit('selectCell', {
    resource,
    time: schedule.start_hour,
    weekday: schedule.weekday
  })
}

const changeWeek = (delta: number) => {
  weekStart.value = getWeekStart(new Date(weekStart.value.getTime() + delta * 7 * 24 * 60 * 60 * 1000))
}

const handleDragEnter = (resourceId: number, time: number) => {
  dragTarget.value = { resourceId, time }
}

const handleDragLeave = () => {
  if (dragTarget.value) {
    const key = `${dragTarget.value.resourceId}-${dragTarget.value.time}`
    delete validationResults.value[key]
  }
  dragTarget.value = null
}

const handleDragOver = (event: DragEvent) => {
  if (dragTarget.value) {
    const key = `${dragTarget.value.resourceId}-${dragTarget.value.time}`
    validationResults.value[key] = { valid: true }
  }
}

const handleDrop = async (event: DragEvent) => {
  event.preventDefault()

  if (!dragTarget.value) return

  const data = event.dataTransfer?.getData('application/json')
  if (!data) return

  const parsed = JSON.parse(data)
  const { type, item } = parsed

  const key = `${dragTarget.value.resourceId}-${dragTarget.value.time}`
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
  fetchData()
}

// 計算時段槽寬度
const calculateSlotWidth = () => {
  if (!tableContainerRef.value) return

  // 找到表格容器的寬度
  const containerWidth = tableContainerRef.value.offsetWidth
  // 減去資源列寬度，再除以時段數量
  slotWidth.value = Math.max(50, (containerWidth - RESOURCE_COLUMN_WIDTH) / SLOT_COUNT)
}

onMounted(async () => {
  fetchData()

  // 等待 DOM 更新後計算槽寬度
  await nextTick()
  calculateSlotWidth()

  // 強制更新一次
  await nextTick()

  // 使用 ResizeObserver 監控表格容器的大小變化
  if (tableContainerRef.value) {
    const resizeObserver = new ResizeObserver(() => {
      calculateSlotWidth()
    })
    resizeObserver.observe(tableContainerRef.value)

    // 頁面卸載時取消監控
    onUnmounted(() => {
      resizeObserver.disconnect()
    })
  }
})
</script>
