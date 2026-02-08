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
      @select-date="handleDateSelect"
      @create-schedule="showCreateModal = true"
      @add-personal-event="$emit('add-personal-event')"
      @add-exception="$emit('add-exception')"
      @export="$emit('export')"
      @update:selected-teacher-id="selectedTeacherId = $event"
      @update:selected-room-id="selectedRoomId = $event"
      @clear-teacher-filter="selectedTeacherId = -1"
      @clear-room-filter="selectedRoomId = -1"
      @clear-all-filters="clearAllFilters"
    />

    <!-- 內容區域 -->
    <div
      class="flex-1 overflow-auto p-4"
      @dragover.prevent="handleDragOver"
      @drop="handleDrop"
    >
      <!-- 調試資訊（臨時，開發用） -->
      <div v-if="false" class="p-4 bg-slate-900 text-xs text-green-400 font-mono">
        <p>isLoading: {{ isLoading }}</p>
        <p>hasError: {{ hasError }}</p>
        <p>schedules.length: {{ schedules.length }}</p>
        <p>filteredSchedules.length: {{ filteredSchedules.length }}</p>
        <p>weekStart: {{ weekStart }}</p>
        <p>weekEnd: {{ weekEnd }}</p>
        <p v-if="filteredSchedules.length > 0">
          First schedule: {{ JSON.stringify(filteredSchedules[0], null, 2) }}
        </p>
      </div>

      <!-- 載入中狀態 -->
      <div v-if="isLoading" class="flex items-center justify-center h-64">
        <div class="flex flex-col items-center gap-3">
          <div class="w-8 h-8 border-2 border-primary-500 border-t-transparent rounded-full animate-spin"></div>
          <p class="text-sm text-slate-400">載入課表中...</p>
        </div>
      </div>

      <!-- 錯誤狀態 -->
      <div v-else-if="hasError" class="flex items-center justify-center h-64">
        <div class="flex flex-col items-center gap-3">
          <div class="w-12 h-12 rounded-full bg-red-500/20 flex items-center justify-center">
            <svg class="w-6 h-6 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
          <p class="text-white font-medium">載入失敗</p>
          <p class="text-sm text-slate-400">{{ errorMessage }}</p>
          <button
            @click="retryFetch"
            class="mt-2 px-4 py-2 bg-primary-500 text-white rounded-lg hover:bg-primary-600 transition-colors"
          >
            重試
          </button>
        </div>
      </div>

      <!-- 空狀態 -->
      <div v-else-if="filteredSchedules.length === 0" class="flex items-center justify-center h-64">
        <div class="flex flex-col items-center gap-3">
          <div class="w-12 h-12 rounded-full bg-slate-800 flex items-center justify-center">
            <svg class="w-6 h-6 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
          </div>
          <p class="text-white font-medium">暫無課表資料</p>
          <p class="text-sm text-slate-400">請確認日期範圍內是否有排課規則</p>
        </div>
      </div>

      <!-- 正常課表視圖 -->
      <template v-else>
        <!-- 桌面版週曆視圖 (lg 以上) -->
        <WeekGrid
          v-if="isDesktop && weekStart"
          :schedules="filteredSchedules"
          :week-label="weekLabel"
          :week-start="weekStart"
          :card-info-type="effectiveCardInfoType"
          :validation-results="validationResults"
          :slot-width="slotWidth"
          @drag-enter="handleDragEnter"
          @drag-leave="handleDragLeave"
          @select-schedule="selectSchedule"
          @overlap-click="handleOverlapClick"
        />

        <!-- 平板版緊湊視圖 (md 到 lg) -->
        <CompactWeekView
          v-else-if="isTablet"
          :schedules="filteredSchedules"
          :week-label="weekLabel"
          :week-start="weekStart"
          :card-info-type="effectiveCardInfoType"
          @change-week="changeWeek"
          @select-schedule="selectSchedule"
        />

        <!-- 手機版日曆列表視圖 (md 以下) -->
        <MobileWeekView
          v-else
          :schedules="filteredSchedules"
          :week-label="weekLabel"
          :week-start="weekStart"
          @change-week="changeWeek"
          @select-schedule="selectSchedule"
        />
      </template>
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
        @suspend="handleSuspend"
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

    <!-- 停課模式選擇彈窗 -->
    <Teleport to="body">
      <UpdateModeModal
        v-if="showSuspendModal && mode === 'admin'"
        :show="showSuspendModal"
        :rule-name="selectedSchedule?.offering_name || selectedSchedule?.title || ''"
        :rule-date="selectedSchedule?.date ? formatDate(selectedSchedule.date) : ''"
        :is-suspend-mode="true"
        @close="showSuspendModal = false"
        @confirm="processSuspend"
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
import CalendarHeader from './CalendarHeader.vue'
import WeekGrid from './WeekGrid.vue'
import CompactWeekView from './CompactWeekView.vue'
import MobileWeekView from './MobileWeekView.vue'
import ActionDialog from './ActionDialog.vue'
import OverlapDialog from './OverlapDialog.vue'
import ScheduleDetailPanel from './ScheduleDetailPanel.vue'
import UpdateModeModal from './UpdateModeModal.vue'
import ScheduleRuleModal from './ScheduleRuleModal.vue'

// ============================================
// Props 定義
// ============================================

const props = defineProps<{
  // 模式：'admin' 或 'teacher'
  mode: 'admin' | 'teacher'
  // 排課資料（可選，如果提供則使用此資料，否則自動獲取）
  schedules?: any[]
  // 週起始日期（可選，如果提供則使用此日期）
  weekStart?: Date
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
  if (!date || !(date instanceof Date) || isNaN(date.getTime())) {
    // 如果日期無效，使用今天的日期
    date = new Date()
  }
  const d = new Date(date)
  d.setHours(0, 0, 0, 0)
  const day = d.getDay()
  const diff = d.getDate() - day + (day === 0 ? -6 : 1)
  return new Date(d.setDate(diff))
}

// 週相關狀態
// 優先使用外部傳入的 weekStart，否則使用今天日期
const weekStart = ref(props.weekStart ? getWeekStart(props.weekStart) : getWeekStart(new Date()))

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

// ============================================
// 篩選狀態 - 修復：確保刷新頁面時預設為 -1（顯示全部）
// ============================================

const selectedTeacherId = ref<number>(-1)
const selectedRoomId = ref<number>(-1)

// 確保刷新頁面時下拉選單顯示"全部"選項
// 使用 watch 監控初始渲染，確保 selectedTeacherId 和 selectedRoomId 正確初始化為 -1
onMounted(() => {
  // 強制確保初始值為 -1（顯示全部）
  if (selectedTeacherId.value !== -1) {
    selectedTeacherId.value = -1
  }
  if (selectedRoomId.value !== -1) {
    selectedRoomId.value = -1
  }
})

// 彈窗狀態
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showUpdateModeModal = ref(false)
const showSuspendModal = ref(false)
const editingRule = ref<any>(null)
const pendingUpdateMode = ref<string>('')
const selectedSchedule = ref<any>(null)
const validationResults = ref<Record<string, any>>({})
const slotWidth = ref(120) // 預設值，WeekGrid 自己會計算實際值

// 老師端動作選擇對話框狀態
const showActionDialog = ref(false)
const actionDialogItem = ref<any>(null)

// 管理員端重疊課程選擇對話框狀態
const showOverlapDialog = ref(false)
const overlapSchedules = ref<any[]>([])
const overlapTimeSlot = ref<{ weekday: number; start_hour: number; start_minute: number } | null>(null)

// 排課資料
const schedules = ref<any[]>([])

// 載入狀態
const isLoading = ref(false)
const hasError = ref(false)
const errorMessage = ref('')

// 響應式視圖狀態
const isDesktop = ref(false)
const isTablet = ref(false)

// 更新視圖狀態
const updateViewState = () => {
  const width = window.innerWidth
  isDesktop.value = width >= 1024 // lg breakpoint
  isTablet.value = width >= 768 && width < 1024 // md to lg
}

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
  selectedTeacherId.value = -1
  selectedRoomId.value = -1
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
  const displayResult = displaySchedules.value

  const result = selectedTeacherId.value <= 0 && selectedRoomId.value <= 0
    ? displayResult
    : displayResult.filter(schedule => {
      const teacherMatch = selectedTeacherId.value <= 0 || schedule.teacher_id === selectedTeacherId.value
      const roomMatch = selectedRoomId.value <= 0 || schedule.room_id === selectedRoomId.value
      return teacherMatch && roomMatch
    })

  return result
})

// ============================================
// 週切換
// ============================================

// 週切換 - 固定加減 7 天，不再強制對齊週一
const changeWeek = (delta: number) => {
  isLoading.value = true
  hasError.value = false

  // 直接加減 7 天，保持選中的日期為該週第一天
  const currentDate = new Date(weekStart.value)
  currentDate.setDate(currentDate.getDate() + delta * 7)
  weekStart.value = currentDate

  emit('update:weekStart', weekStart.value)
}

// 處理日期選擇 - 直接跳轉到指定日期
const handleDateSelect = (date: Date) => {
  isLoading.value = true
  hasError.value = false

  // 設定選中的日期為該週第一天
  const newDate = new Date(date)
  newDate.setHours(0, 0, 0, 0)
  weekStart.value = newDate

  emit('update:weekStart', weekStart.value)
}

// ============================================
// 取得排課資料
// ============================================

const fetchSchedules = async () => {
  isLoading.value = true
  hasError.value = false
  errorMessage.value = ''

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

    // 處理 API 響應格式
    let expandedSchedules: any[] = []
    if (Array.isArray(response)) {
      // 後端直接返回陣列
      expandedSchedules = response
    } else if (response?.datas) {
      // 後端返回 { code, datas } 格式
      expandedSchedules = response.datas
    } else if (response?.data) {
      // 另一種格式 { code, data }
      expandedSchedules = response.data
    } else {
      expandedSchedules = []
    }

    // 確保 expandedSchedules 是陣列
    if (!Array.isArray(expandedSchedules)) {
      expandedSchedules = []
    }

    const scheduleList = expandedSchedules.map((schedule: any, index: number) => {
      try {
        // 檢測數據格式：後端可能返回排課規則或展開後的課表場次
        const isExpandedFormat = schedule.date !== undefined
        const isRuleFormat = schedule.weekday !== undefined

        let date: Date
        let weekday: number

        if (isExpandedFormat) {
          // 展開後的課表場次格式：已有 date 欄位
          // 注意：後端返回的 date 可能是 "2026-02-01" 格式
          date = new Date(schedule.date + 'T00:00:00+08:00')
          weekday = date.getDay() === 0 ? 7 : date.getDay()
        } else if (isRuleFormat) {
          // 排課規則格式：需要根據 weekday 和 effective_range 計算
          // 對於週視圖，我們需要計算每個規則在該週的具體日期
          const effectiveStartDate = schedule.effective_range?.start_date
          const effectiveEndDate = schedule.effective_range?.end_date

          // 計算規則適用的第一個日期（與週視圖的開始日期最接近的規則日期）
          const ruleWeekday = schedule.weekday
          const weekStartDate = new Date(weekStart.value)

          // 找到週開始日期之後第一個符合 weekday 的日期
          let targetDate = new Date(weekStartDate)
          const targetDay = ruleWeekday === 7 ? 0 : ruleWeekday
          while (targetDate.getDay() !== targetDay) {
            targetDate.setDate(targetDate.getDate() + 1)
          }

          // 檢查是否在有效範圍內
          if (effectiveStartDate) {
            const start = new Date(effectiveStartDate)
            // 比較日期部分（忽略時間）
            const targetTime = new Date(targetDate).setHours(0, 0, 0, 0)
            const startTime = new Date(start).setHours(0, 0, 0, 0)
            
            if (targetTime < startTime) {
              return null
            }
          }

          if (effectiveEndDate) {
            const end = new Date(effectiveEndDate)
            if (targetDate > end) {
              return null
            }
          }

          date = targetDate
          // 確保 weekday 從計算出的 date 重新計算，而不是依賴後端可能錯誤的 weekday
          weekday = date.getDay() === 0 ? 7 : date.getDay()
        } else {
          // 無法識別的格式，跳過
          return null
        }

        const startTime = schedule.start_time || '09:00'
        const endTime = schedule.end_time || '10:00'
        const [startHour, startMinute] = startTime.split(':').map(Number)
        const durationMinutes = calculateDurationMinutes(startTime, endTime)

        // 從關聯資料中取得老師和教室名稱
        const teacherName = schedule.teacher?.name || schedule.teacher_name || '-'
        const roomName = schedule.room?.name || schedule.room_name || '-'
        const offeringName = schedule.offering?.name || schedule.offering_name || schedule.title || '-'

        const dateString = date.toISOString().split('T')[0]

        return {
          id: schedule.rule_id || schedule.id,
          key: `${schedule.rule_id || schedule.id}-${weekday}-${startTime}-${dateString}`,
          offering_name: offeringName,
          teacher_name: teacherName,
          teacher_id: schedule.teacher_id || schedule.teacher?.id,
          center_name: schedule.center_name || '-',
          center_id: schedule.center_id,
          room_id: schedule.room_id || schedule.room?.id,
          room_name: roomName,
          weekday: weekday,
          start_time: startTime,
          end_time: endTime,
          start_hour: startHour,
          start_minute: startMinute,
          duration_minutes: durationMinutes,
          date: dateString,
          has_exception: schedule.has_exception || false,
          exception_type: schedule.exception_type || null,
          exception_info: schedule.exception_info || null,
          rule: schedule.rule || null,
          offering_id: schedule.offering_id || schedule.offering?.id,
          effective_range: schedule.effective_range || null,
          rule_id: schedule.rule_id || schedule.id, // 添加 rule_id 欄位
        }
      } catch (err) {
        console.error('[ScheduleGrid] Error processing schedule at index', index, err)
        return null
      }
    }).filter((item): item is NonNullable<typeof item> => item !== null)

    schedules.value = scheduleList

    // 注意：slotWidth 現在由 WeekGrid 組件自己計算，不需要在這裡呼叫
  } catch (error: any) {
    hasError.value = true
    errorMessage.value = error?.message || '無法載入課表資料，請稍後重試'
    schedules.value = []
  } finally {
    isLoading.value = false
  }
}

// 重試載入
const retryFetch = async () => {
  await fetchSchedules()
}

// ============================================
// 計算工具函數
// ============================================

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

// ============================================
// 停課處理
// ============================================

const handleSuspend = () => {
  if (selectedSchedule.value) {
    showSuspendModal.value = true
  }
}

const processSuspend = async (mode: 'SINGLE' | 'FUTURE') => {
  if (!selectedSchedule.value) return

  try {
    const api = useApi()

    if (mode === 'SINGLE') {
      // 單次停課：建立 Exception 取消該場次
      await api.post('/admin/scheduling/exceptions', {
        center_id: selectedSchedule.value.center_id,
        rule_id: selectedSchedule.value.id,
        original_date: selectedSchedule.value.date,
        type: 'CANCEL',
        reason: '由管理員於首頁停課(單次)',
      })
    } else if (mode === 'FUTURE') {
      // 從此以後停課：更新規則結束日期
      const currentDate = new Date(selectedSchedule.value.date)
      currentDate.setDate(currentDate.getDate() - 1)
      const endDate = formatDateToString(currentDate)

      await api.put(`/admin/rules/${selectedSchedule.value.id}`, {
        end_date: endDate,
        update_mode: 'ALL',
        name: selectedSchedule.value.title,
        reason: '由管理員於首頁選擇從此以後停課',
      })
    }

    // 刷新課表
    await fetchSchedules()
  } catch (err) {
    console.error('Failed to suspend schedule:', err)
    await alertError('停課失敗，請稍後再試')
    return
  }

  // 關閉所有彈窗
  showSuspendModal.value = false
  selectedSchedule.value = null
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

// 修復：添加標記來區分 schedules 是來自 props 還是內部獲取
const schedulesFromProps = ref(false)

watch(weekStart, async (newWeekStart, oldWeekStart) => {
  // 如果 schedules 來自 props，不應該調用 fetchSchedules()
  // 因為外部會負責提供數據
  if (schedulesFromProps.value) {
    return
  }

  // 如果 schedules 來自內部獲取 (schedules.value)，則重新獲取
  if (!props.schedules || props.schedules.length === 0) {
    isLoading.value = true
    hasError.value = false
    await fetchSchedules()
  }
})

// 監聽 props.schedules 變化，更新標記
watch(() => props.schedules, (newSchedules) => {
  schedulesFromProps.value = !!(newSchedules && newSchedules.length > 0)

  // 如果有新的 schedules，立即設置 isLoading 為 false
  if (newSchedules && newSchedules.length > 0) {
    isLoading.value = false
    hasError.value = false
    errorMessage.value = ''
  }
}, { immediate: true })

// 監聽外部傳入的 weekStart 變化，同步到內部狀態
watch(() => props.weekStart, (newWeekStart) => {
  if (newWeekStart) {
    weekStart.value = getWeekStart(newWeekStart)
  }
}, { immediate: true })

// ============================================
// 生命週期
// ============================================

onMounted(async () => {
  // 設置 schedules 來源標記
  schedulesFromProps.value = !!(props.schedules && props.schedules.length > 0)

  // 如果已經有 props.schedules，立即設置 isLoading 為 false
  if (props.schedules && props.schedules.length > 0) {
    isLoading.value = false
    hasError.value = false
  }

  if (!props.schedules || props.schedules.length === 0) {
    await fetchSchedules()
  }

  // 修復：無論是管理員還是老師模式，都嘗試獲取資源
  // 這樣下拉選單才能正確渲染
  fetchAllResources()

  // 初始化響應式狀態
  updateViewState()
  window.addEventListener('resize', updateViewState)
})

// 監控 isDesktop 變化
watch([isDesktop, weekStart], async ([newIsDesktop, newWeekStart]) => {
  // 當切換到桌面模式且有週起始日期時，確保週曆視圖能正確渲染
  // WeekGrid 自己會處理 slotWidth 的計算
}, { immediate: true })

onUnmounted(() => {
  // 移除 resize 監聽
  window.removeEventListener('resize', updateViewState)
})
</script>
