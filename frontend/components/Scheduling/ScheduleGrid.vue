<template>
  <div class="h-full flex flex-col glass-card overflow-hidden">
    <!-- 日曆標頭 -->
    <CalendarHeader
      :mode="mode"
      :week-label="weekLabel"
      :teacher-list="teacherList"
      :room-list="roomList"
      :selected-teacher-id="selectedTeacherId"
      :selected-room-id="selectedRoomId"
      :selected-teacher-name="selectedTeacherName"
      :selected-room-name="selectedRoomName"
      :show-create-button="effectiveShowCreateButton"
      :show-personal-event-button="effectiveShowPersonalEventButton"
      :show-exception-button="effectiveShowExceptionButton"
      :show-export-button="effectiveShowExportButton"
      :show-help-tooltip="effectiveShowHelpTooltip"
      @change-week="changeWeek"
      @create-schedule="showCreateModal = true"
      @add-personal-event="$emit('add-personal-event')"
      @add-exception="$emit('add-exception')"
      @export="$emit('export')"
      @update:selected-teacher-id="selectedTeacherId = $event"
      @update:selected-room-id="selectedRoomId = $event"
      @clear-teacher-filter="selectedTeacherId = null"
      @clear-room-filter="selectedRoomId = null"
      @clear-all-filters="clearAllFilters"
    />

    <!-- 內容區域 -->
    <div
      class="flex-1 overflow-auto p-4"
      @dragover.prevent="handleDragOver"
      @drop="handleDrop"
    >
      <!-- 桌面版週曆視圖 -->
      <WeekGrid
        :schedules="filteredSchedules"
        :week-label="weekLabel"
        :card-info-type="effectiveCardInfoType"
        :validation-results="validationResults"
        :slot-width="slotWidth"
        @drag-enter="handleDragEnter"
        @drag-leave="handleDragLeave"
        @select-schedule="selectSchedule"
        @overlap-click="handleOverlapClick"
      />

      <!-- 手機版日曆列表視圖 -->
      <MobileWeekView
        :schedules="filteredSchedules"
        :week-label="weekLabel"
        :week-start="weekStart"
        @change-week="changeWeek"
        @select-schedule="selectSchedule"
      />
    </div>

    <!-- 管理員專屬彈窗 -->
    <Teleport to="body">
      <ScheduleDetailPanel
        v-if="selectedSchedule && mode === 'admin'"
        :time="selectedSchedule.start_hour"
        :weekday="selectedSchedule.weekday"
        :schedule="selectedSchedule"
        @close="closeSchedulePanel"
        @edit="handleEdit"
        @delete="handleDelete"
      />
    </Teleport>

    <Teleport to="body">
      <UpdateModeModal
        v-if="showUpdateModeModal && mode === 'admin'"
        :show="showUpdateModeModal"
        :rule-name="editingRule?.offering_name"
        :rule-date="editingRule?.date ? formatDate(editingRule.date) : ''"
        @close="handleUpdateModeClose"
        @confirm="handleUpdateModeConfirm"
      />
    </Teleport>

    <Teleport to="body">
      <ScheduleRuleModal
        v-if="showCreateModal && mode === 'admin'"
        @close="showCreateModal = false"
        @created="handleRuleCreated"
      />
      <ScheduleRuleModal
        v-if="showEditModal && mode === 'admin'"
        :editing-rule="editingRule"
        :update-mode="pendingUpdateMode"
        @close="handleEditModalClose"
        @submit="handleRuleUpdated"
      />
    </Teleport>

    <!-- 老師端專屬彈窗 - 動作選擇對話框 -->
    <ActionDialog
      :visible="showActionDialog"
      :item="actionDialogItem"
      @close="closeActionDialog"
      @action="handleActionSelect"
    />

    <!-- 管理員端專屬彈窗 - 重疊課程選擇對話框 -->
    <OverlapDialog
      :visible="showOverlapDialog"
      :schedules="overlapSchedules"
      :time-slot="overlapTimeSlot"
      :card-info-type="effectiveCardInfoType"
      @close="closeOverlapDialog"
      @select="selectFromOverlap"
    />
  </div>
</template>

<script setup lang="ts">
import { formatDateToString, formatDate } from '~/composables/useTaiwanTime'
import { nextTick, ref, computed, watch, onMounted, onUnmounted } from 'vue'

// 引入子組件
import CalendarHeader from './Schedule/CalendarHeader.vue'
import WeekGrid from './Schedule/WeekGrid.vue'
import MobileWeekView from './Schedule/MobileWeekView.vue'
import ActionDialog from './Schedule/ActionDialog.vue'
import OverlapDialog from './Schedule/OverlapDialog.vue'

// ============================================
// Props 定義
// ============================================

const props = defineProps<{
  // 模式：'admin' 或 'teacher'
  mode: 'admin' | 'teacher'
  // 排課資料（可選，如果提供則使用此資料，否則自動獲取）
  schedules?: any[]
  // API 端點
  apiEndpoint: string
  // 卡片顯示類型：'teacher' 顯示老師名稱，'center' 顯示中心名稱
  cardInfoType?: 'teacher' | 'center'
  // 是否顯示新增排課按鈕
  showCreateButton?: boolean
  // 是否顯示個人行程按鈕（老師端）
  showPersonalEventButton?: boolean
  // 是否顯示請假/調課按鈕（老師端）
  showExceptionButton?: boolean
  // 是否顯示匯出按鈕（老師端）
  showExportButton?: boolean
  // 是否顯示說明提示
  showHelpTooltip?: boolean
}>()

// ============================================
// Emits 定義
// ============================================

const emit = defineEmits<{
  selectCell: [{ time: number; weekday: number }]
  'update:weekStart': [value: Date]
  'select-schedule': [schedule: any]
  'add-personal-event': []
  'add-exception': []
  'export': []
  'edit-personal-event': [event: any]
  'delete-personal-event': [event: any]
  'personal-event-note': [event: any]
}>()

// ============================================
// Composables
// ============================================

const { confirm: confirmDialog, error: alertError } = useAlert()
const { resourceCache, fetchAllResources } = useResourceCache()

// ============================================
// 常量定義
// ============================================

const TIME_SLOT_HEIGHT = 60 // 每個時段格子的高度 (px)
const TIME_COLUMN_WIDTH = 80 // 時間列寬度 (px)

// ============================================
// 計算屬性（預設值）
// ============================================

const effectiveApiEndpoint = computed(() => props.apiEndpoint || '/admin/expand-rules')
const effectiveCardInfoType = computed(() => props.cardInfoType || 'teacher')
const effectiveShowCreateButton = computed(() => props.showCreateButton ?? true)
const effectiveShowPersonalEventButton = computed(() => props.showPersonalEventButton ?? false)
const effectiveShowExceptionButton = computed(() => props.showExceptionButton ?? false)
const effectiveShowExportButton = computed(() => props.showExportButton ?? false)
const effectiveShowHelpTooltip = computed(() => props.showHelpTooltip ?? true)

// ============================================
// 狀態管理
// ============================================

// 週起始日期計算
const getWeekStart = (date: Date): Date => {
  const d = new Date(date)
  const day = d.getDay()
  const diff = d.getDate() - day + (day === 0 ? -6 : 1)
  return new Date(d.setDate(diff))
}

// 週相關狀態
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

// 篩選狀態
const selectedTeacherId = ref<number | null>(null)
const selectedRoomId = ref<number | null>(null)

// 彈窗狀態
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showUpdateModeModal = ref(false)
const editingRule = ref<any>(null)
const pendingUpdateMode = ref<string>('')
const selectedSchedule = ref<any>(null)
const validationResults = ref<Record<string, any>>({})
const slotWidth = ref(100)

// 老師端動作選擇對話框狀態
const showActionDialog = ref(false)
const actionDialogItem = ref<any>(null)

// 管理員端重疊課程選擇對話框狀態
const showOverlapDialog = ref(false)
const overlapSchedules = ref<any[]>([])
const overlapTimeSlot = ref<{ weekday: number; start_hour: number; start_minute: number } | null>(null)

// 排課資料
const schedules = ref<any[]>([])

// ResizeObserver 引用
const resizeObserver = ref<ResizeObserver | null>(null)

// ============================================
// 資源列表計算屬性
// ============================================

const teacherList = computed(() => Array.from(resourceCache.value.teachers.values()))
const roomList = computed(() => Array.from(resourceCache.value.rooms.values()))

const selectedTeacherName = computed(() => {
  const teacher = resourceCache.value.teachers.get(selectedTeacherId.value)
  return teacher?.name || ''
})

const selectedRoomName = computed(() => {
  const room = resourceCache.value.rooms.get(selectedRoomId.value)
  return room?.name || ''
})

const clearAllFilters = () => {
  selectedTeacherId.value = null
  selectedRoomId.value = null
}

// ============================================
// 篩選後的課程
// ============================================

const displaySchedules = computed(() => {
  const sourceSchedules = props.schedules && props.schedules.length > 0
    ? props.schedules
    : schedules.value

  const seen = new Set<string>()
  const result: any[] = []

  for (const schedule of sourceSchedules) {
    const key = `${schedule.id}-${schedule.weekday}-${schedule.start_time}`
    if (!seen.has(key)) {
      seen.add(key)
      result.push(schedule)
    }
  }

  return result
})

const filteredSchedules = computed(() => {
  if (!selectedTeacherId.value && !selectedRoomId.value) {
    return displaySchedules.value
  }

  return displaySchedules.value.filter(schedule => {
    const teacherMatch = !selectedTeacherId.value || schedule.teacher_id === selectedTeacherId.value
    const roomMatch = !selectedRoomId.value || schedule.room_id === selectedRoomId.value
    return teacherMatch && roomMatch
  })
})

// ============================================
// 週切換
// ============================================

const changeWeek = (delta: number) => {
  weekStart.value = getWeekStart(new Date(weekStart.value.getTime() + delta * 7 * 24 * 60 * 60 * 1000))
  emit('update:weekStart', weekStart.value)
}

// ============================================
// 取得排課資料
// ============================================

const fetchSchedules = async () => {
  try {
    const api = useApi()

    const startDate = formatDateToString(weekStart.value)
    const endDate = formatDateToString(weekEnd.value)

    let response
    if (props.mode === 'teacher') {
      response = await api.get<{ code: number; datas: any[] }>('/teacher/schedules', {
        start_date: startDate,
        end_date: endDate,
      })
    } else {
      response = await api.post<{ code: number; datas: any[] }>(effectiveApiEndpoint.value, {
        rule_ids: [],
        start_date: startDate,
        end_date: endDate,
      })
    }

    const expandedSchedules = response.datas || []

    const scheduleList = expandedSchedules.map((schedule: any) => {
      const date = new Date(schedule.date)
      const weekday = date.getDay() === 0 ? 7 : date.getDay()
      const startTime = schedule.start_time || '09:00'
      const endTime = schedule.end_time || '10:00'
      const [startHour, startMinute] = startTime.split(':').map(Number)
      const durationMinutes = calculateDurationMinutes(startTime, endTime)

      return {
        id: schedule.rule_id || schedule.id,
        key: `${schedule.rule_id || schedule.id}-${weekday}-${startTime}-${schedule.date}`,
        offering_name: schedule.title || schedule.offering_name || '-',
        teacher_name: schedule.teacher_name || '-',
        teacher_id: schedule.teacher_id,
        center_name: schedule.center_name || '-',
        center_id: schedule.center_id,
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

    await nextTick()
    calculateSlotWidth()
  } catch (error) {
    console.error('Failed to fetch schedules:', error)
    schedules.value = []
  }
}

// ============================================
// 計算工具函數
// ============================================

const calculateSlotWidth = () => {
  const container = document.querySelector('.min-w-\\[800px\\]') as HTMLElement
  if (container) {
    const containerWidth = container.offsetWidth
    slotWidth.value = Math.max(80, (containerWidth - TIME_COLUMN_WIDTH) / 7)
  }
}

const calculateDurationMinutes = (startTime: string, endTime: string): number => {
  const [startHour, startMinute] = startTime.split(':').map(Number)
  const [endHour, endMinute] = endTime.split(':').map(Number)

  const startMinutes = startHour * 60 + startMinute
  let endMinutes = endHour * 60 + endMinute

  if (endMinutes <= startMinutes) {
    endMinutes += 24 * 60
  }

  return endMinutes - startMinutes
}

const formatTime = (hour: number): string => {
  return `${hour.toString().padStart(2, '0')}:00`
}

// ============================================
// 互動處理
// ============================================

const selectSchedule = (schedule: any) => {
  if (props.mode === 'teacher') {
    actionDialogItem.value = schedule
    showActionDialog.value = true
    emit('select-schedule', schedule)
  } else {
    selectedSchedule.value = schedule
    emit('selectCell', { time: schedule.start_hour, weekday: schedule.weekday })
  }
}

const closeActionDialog = () => {
  showActionDialog.value = false
  actionDialogItem.value = null
}

const closeSchedulePanel = () => {
  selectedSchedule.value = null
}

const handleOverlapClick = (schedule: any) => {
  const schedulesAtSameTime = filteredSchedules.value.filter(s =>
    s.weekday === schedule.weekday &&
    s.start_hour === schedule.start_hour &&
    s.start_minute === schedule.start_minute
  )

  if (schedulesAtSameTime.length === 1) {
    selectSchedule(schedule)
  } else {
    overlapSchedules.value = schedulesAtSameTime
    overlapTimeSlot.value = {
      weekday: schedule.weekday,
      start_hour: schedule.start_hour,
      start_minute: schedule.start_minute
    }
    showOverlapDialog.value = true
  }
}

const closeOverlapDialog = () => {
  showOverlapDialog.value = false
  overlapSchedules.value = []
  overlapTimeSlot.value = null
}

const selectFromOverlap = (schedule: any) => {
  selectSchedule(schedule)
  closeOverlapDialog()
}

const handleActionSelect = (action: 'exception' | 'note' | 'edit' | 'delete') => {
  const item = actionDialogItem.value
  if (!item) return

  if (item.is_personal_event) {
    if (action === 'edit') {
      emit('edit-personal-event', item)
    } else if (action === 'delete') {
      emit('delete-personal-event', item)
    } else if (action === 'note') {
      setTimeout(() => {
        emit('personal-event-note', { ...item, action: 'note' })
      }, 100)
    }
  } else {
    if (action === 'exception') {
      emit('add-exception', item)
    } else if (action === 'note') {
      setTimeout(() => {
        emit('select-schedule', { ...item, action: 'note' })
      }, 100)
    }
  }

  closeActionDialog()
}

// ============================================
// 管理員端編輯/刪除
// ============================================

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
    selectedSchedule.value = null
    editingRule.value = null
    pendingUpdateMode.value = ''
  } catch (err) {
    console.error('Failed to update rule:', err)
    await alertError('更新失敗，請稍後再試')
  }
}

const handleRuleCreated = () => {
  fetchSchedules()
}

// ============================================
// 拖曳處理（僅管理員）
// ============================================

const dragTarget = ref<{ time: number; day: number } | null>(null)

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
      start_time: `${formatDateToString(weekStart.value)}T${formatTime(dragTarget.value.time)}:00`,
      end_time: `${formatDateToString(weekStart.value)}T${formatTime(dragTarget.value.time + 1)}:00`,
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

// ============================================
// 監聽週變化
// ============================================

watch(weekStart, async () => {
  if (!props.schedules || props.schedules.length === 0) {
    await fetchSchedules()
  }
})

// ============================================
// 生命週期
// ============================================

onMounted(async () => {
  if (!props.schedules || props.schedules.length === 0) {
    await fetchSchedules()
  }

  if (props.mode === 'admin') {
    fetchAllResources()
  }

  await nextTick()
  calculateSlotWidth()

  // 監控容器大小變化
  const container = document.querySelector('.min-w-\\[800px\\]') as HTMLElement
  if (container) {
    resizeObserver.value = new ResizeObserver(() => {
      calculateSlotWidth()
    })
    resizeObserver.value.observe(container)
  }
})

onUnmounted(() => {
  if (resizeObserver.value) {
    resizeObserver.value.disconnect()
  }
})
</script>
