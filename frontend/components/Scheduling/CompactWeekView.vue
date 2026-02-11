<template>
  <!-- 三天/五天視圖 - 平板和小型桌面版 -->
  <div class="hidden md:block lg:hidden" ref="calendarContainerRef">
    <!-- 視圖切換和控制列 -->
    <div class="flex items-center justify-between mb-4">
      <div class="flex items-center gap-2">
        <!-- 視圖模式選擇 -->
        <div class="flex bg-slate-800 rounded-lg p-1">
          <button
            v-for="mode in viewModes"
            :key="mode.value"
            @click="viewMode = mode.value"
            class="px-3 py-1 rounded-md text-sm transition-colors"
            :class="viewMode === mode.value 
              ? 'bg-primary-500 text-white' 
              : 'text-slate-400 hover:text-white'"
          >
            {{ mode.label }}
          </button>
        </div>
      </div>
      
      <!-- 日期範圍標籤 -->
      <span class="text-sm text-slate-400">{{ viewRangeLabel }}</span>
    </div>

    <!-- 簡化網格視圖 -->
    <div class="bg-slate-900/50 rounded-lg overflow-hidden">
      <!-- 表頭 -->
      <div 
        class="grid sticky top-0 z-10 bg-slate-900/95 backdrop-blur-sm"
        :style="gridStyle"
      >
        <div 
          v-for="day in visibleDays" 
          :key="day.value"
          class="p-2 border-b border-white/10 text-center"
        >
          <span class="text-xs text-slate-400">{{ day.shortName }}</span>
          <div class="text-sm font-medium text-slate-100 mt-1">{{ day.dateLabel }}</div>
        </div>
      </div>

      <!-- 時間網格 -->
      <div class="relative">
        <!-- 簡化時間軸 -->
        <div class="absolute left-0 top-0 bottom-0 w-16 bg-slate-900/50 border-r border-white/5">
          <div 
            v-for="time in visibleTimeSlots" 
            :key="time"
            class="h-16 border-b border-white/5 flex items-center justify-center"
          >
            <span class="text-xs text-slate-400">{{ formatTime(time) }}</span>
          </div>
        </div>

        <!-- 課程區域 -->
        <div class="ml-16">
          <div 
            class="grid relative"
            :style="gridStyle"
          >
            <!-- 每個時間槽 -->
            <div 
              v-for="time in visibleTimeSlots" 
              :key="time"
              class="contents"
            >
              <!-- 每個星期幾的格子 -->
              <div
                v-for="day in visibleDays"
                :key="`${time}-${day.value}`"
                class="h-16 border-b border-white/5 border-r relative"
                :class="getCellClass(time, day.value)"
              >
                <!-- 當日課程 -->
                <template v-for="schedule in getSchedulesForCell(time, day.value)" :key="schedule.key">
                  <template v-if="schedule._overlapCount === 1">
                    <!-- 單一課程：正常顯示 -->
                    <div
                      class="absolute left-1 right-1 rounded p-1 cursor-pointer hover:opacity-90 transition-colors overflow-hidden"
                      :class="getScheduleCardClass(schedule)"
                      :style="getScheduleCardStyle(schedule, time)"
                      @click="$emit('select-schedule', schedule)"
                    >
                      <div class="text-xs font-medium truncate" :class="getScheduleTitleClass(schedule)">
                        {{ schedule.offering_name }}
                      </div>
                      <div class="text-xs truncate" :class="getScheduleTimeClass(schedule)">
                        {{ schedule.start_time }}-{{ schedule.end_time }}
                      </div>
                      <div v-if="schedule.teacher_name" class="text-xs truncate opacity-70" :class="getScheduleTimeClass(schedule)">
                        {{ schedule.teacher_name }}
                      </div>
                    </div>
                  </template>
                  <!-- 重疊指示器：僅第一個課程顯示 -->
                  <template v-else-if="schedule._overlapCount > 1 && schedule._isFirstInOverlap">
                    <div
                      class="absolute left-1 right-1 rounded p-1 cursor-pointer hover:opacity-90 transition-colors overflow-hidden"
                      :class="getScheduleCardClass(schedule)"
                      :style="getScheduleCardStyle(schedule, time)"
                      @click="$emit('select-schedule', schedule)"
                    >
                      <div class="flex items-center justify-center h-full">
                        <span class="text-warning-400 font-bold text-sm">
                          {{ schedule._overlapCount }}
                        </span>
                        <span class="text-warning-300 ml-1 text-xs">堂課</span>
                      </div>
                    </div>
                  </template>
                </template>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 週導航 -->
    <div class="flex items-center justify-center gap-4 mt-4">
      <button
        @click="$emit('change-week', -1)"
        class="p-2 rounded-lg hover:bg-white/10 transition-colors"
      >
        <svg class="w-5 h-5 text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>
      <span class="text-sm text-slate-400">{{ weekLabel }}</span>
      <button
        @click="$emit('change-week', 1)"
        class="p-2 rounded-lg hover:bg-white/10 transition-colors"
      >
        <svg class="w-5 h-5 text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, shallowRef } from 'vue'
import { formatDateToString } from '~/composables/useTaiwanTime'

// ============================================
// Props 定義
// ============================================

const props = defineProps<{
  schedules: any[]
  weekLabel: string
  weekStart: Date
  cardInfoType: 'teacher' | 'center'
  // 動態時間段（由矩陣視圖 API 回傳，可選）
  timeSlots?: number[]
}>()

// ============================================
// Emits 定義
// ============================================

const emit = defineEmits<{
  'change-week': [delta: number]
  'select-schedule': [schedule: any]
}>()

// ============================================
// 視圖模式
// ============================================

const viewModes = [
  { value: '3day', label: '3天' },
  { value: '5day', label: '5天' },
  { value: 'week', label: '一週' },
]

const viewMode = ref<'3day' | '5day' | 'week'>('5day')

// ============================================
// 響應式寬度
// ============================================

const calendarContainerRef = ref<HTMLElement | null>(null)
const containerWidth = ref(800)
const cellWidth = ref(120)

const updateWidth = () => {
  if (calendarContainerRef.value) {
    containerWidth.value = calendarContainerRef.value.offsetWidth
    // 計算格子寬度（扣除時間列）
    const timeColumnWidth = 64
    const visibleDaysCount = visibleDays.value.length
    // 確保格子寬度至少 60px
    cellWidth.value = Math.max(60, (containerWidth.value - timeColumnWidth) / visibleDaysCount)
  }
}

// ============================================
// 重疊偵測快取
// ============================================

// 重疊數據快取
const overlapDataCache = new Map<string, { count: number; firstId: number }>()

// 課程快照，用於檢測變化
const schedulesSnapshot = shallowRef<string>('')

// ============================================
// 計算重疊數據
// ============================================

const computeOverlapData = (schedules: any[]) => {
  const countMap: Record<string, number> = {}
  const firstIdMap: Record<string, number> = {}

  for (const schedule of schedules) {
    const key = `${schedule.weekday}-${schedule.start_hour}-${schedule.start_minute}`
    countMap[key] = (countMap[key] || 0) + 1

    // 記錄最小的 ID 作為第一個
    if (!firstIdMap[key] || schedule.id < firstIdMap[key]) {
      firstIdMap[key] = schedule.id
    }
  }

  // 儲存到快取
  overlapDataCache.clear()
  for (const key in countMap) {
    overlapDataCache.set(key, {
      count: countMap[key],
      firstId: firstIdMap[key]
    })
  }
}

// 取得重疊數據
const getOverlapData = (weekday: number, startHour: number, startMinute: number) => {
  const key = `${weekday}-${startHour}-${startMinute}`
  return overlapDataCache.get(key)
}

// ============================================
// 生命週期
// ============================================

onMounted(() => {
  updateWidth()
  window.addEventListener('resize', updateWidth)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateWidth)
})

// 監控 schedules 變化，更新重疊數據
watch(() => props.schedules, (newSchedules) => {
  if (!newSchedules || newSchedules.length === 0) {
    overlapDataCache.clear()
    schedulesSnapshot.value = ''
    return
  }

  // 建立快照
  const snapshot = newSchedules
    .map(s => `${s.id}-${s.weekday}-${s.start_hour}-${s.start_minute}-${s.duration_minutes}`)
    .sort()
    .join(',')

  // 只有當課程實際變化時才重新計算
  if (snapshot !== schedulesSnapshot.value) {
    schedulesSnapshot.value = snapshot
    computeOverlapData(newSchedules)
  }
}, { deep: false })

// ============================================
// 計算屬性
// ============================================

// 可見的星期幾
const visibleDays = computed(() => {
  const days = [
    { value: 1, name: '週一', shortName: '一', dateLabel: '' },
    { value: 2, name: '週二', shortName: '二', dateLabel: '' },
    { value: 3, name: '週三', shortName: '三', dateLabel: '' },
    { value: 4, name: '週四', shortName: '四', dateLabel: '' },
    { value: 5, name: '週五', shortName: '五', dateLabel: '' },
    { value: 6, name: '週六', shortName: '六', dateLabel: '' },
    { value: 7, name: '週日', shortName: '日', dateLabel: '' },
  ]

  // 計算每個星期幾的日期
  const startDate = new Date(props.weekStart)
  
  days.forEach(day => {
    const date = new Date(startDate)
    date.setDate(date.getDate() + (day.value - 1))
    day.dateLabel = `${date.getMonth() + 1}/${date.getDate()}`
  })

  // 根據視圖模式返回不同的天數
  const mode = viewMode.value
  if (mode === '3day') {
    return days.slice(0, 3)
  } else if (mode === '5day') {
    return days.slice(0, 5)
  } else {
    return days
  }
})

// 可見的時間段
const visibleTimeSlots = computed(() => {
  // 如果父組件提供了動態時間段，使用它
  if (props.timeSlots && props.timeSlots.length > 0) {
    return props.timeSlots
  }
  // 否則使用預設時間段（9-22點）
  return [9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22]
})

// 網格樣式
const gridStyle = computed(() => {
  const columnCount = visibleDays.value.length
  return {
    gridTemplateColumns: `repeat(${columnCount}, 1fr)`,
  }
})

// 視圖範圍標籤
const viewRangeLabel = computed(() => {
  const days = visibleDays.value
  if (days.length === 3) return '週一 至 週三'
  if (days.length === 5) return '週一 至 週五'
  return '完整一週'
})

// ============================================
// 工具函數
// ============================================

const formatTime = (hour: number): string => {
  return `${hour.toString().padStart(2, '0')}:00`
}

// 取得某個時間和星期幾的課程（加入重疊數據）
const getSchedulesForCell = (time: number, weekday: number) => {
  const cellSchedules = props.schedules.filter(schedule => {
    // 檢查星期幾
    if (schedule.weekday !== weekday) return false

    // 檢查時間（課程開始時間等於或早於當前時段，且結束時間晚於當前時段）
    const scheduleStartHour = schedule.start_hour

    // 解析結束時間
    const [endHourStr, endMinuteStr] = (schedule.end_time || '23:59').split(':').map(Number)
    const scheduleEndHour = endHourStr

    // 當日課程
    const scheduleEndMinute = endMinuteStr
    return scheduleStartHour < time + 1 &&
           (scheduleEndHour > time || (scheduleEndHour === time && scheduleEndMinute > 0))
  })

  // 為每個課程添加重疊數據
  return cellSchedules.map(schedule => {
    const overlapData = getOverlapData(weekday, schedule.start_hour, schedule.start_minute)
    return {
      ...schedule,
      _overlapCount: overlapData?.count || 1,
      _isFirstInOverlap: overlapData?.firstId === schedule.id
    }
  })
}

// 取得格子樣式
const getCellClass = (time: number, weekday: number) => {
  const schedules = getSchedulesForCell(time, weekday)
  if (schedules.length > 0) {
    // 檢查是否有重疊
    const hasOverlap = schedules.some(s => s._overlapCount > 1)
    if (hasOverlap) {
      return 'bg-warning-500/5'
    }
    return 'bg-primary-500/5'
  }
  return 'hover:bg-white/5'
}

// 取得課程卡片樣式
const getScheduleCardStyle = (schedule: any, time: number) => {
  const slotHeight = 64 // 每個時段高度
  const minuteOffset = (schedule.start_minute / 60) * slotHeight
  const durationHeight = (schedule.duration_minutes / 60) * slotHeight

  // 取得第一個時段的小時數
  const firstSlotHour = visibleTimeSlots.value.length > 0 ? visibleTimeSlots.value[0] : 9

  // 從開始時間計算位置（相對於第一個時段）
  let relativeStartHour = schedule.start_hour - firstSlotHour

  // 保護：如果課程開始時間早於營業時間，將其夾緊到 0
  if (relativeStartHour < 0) {
    relativeStartHour = 0
  }

  const top = (relativeStartHour * slotHeight) + minuteOffset

  // 確保 top 不會是負值
  const clampedTop = Math.max(0, top)

  // 重疊課程調整寬度和位置
  let left = '0.25rem'
  let right = '0.25rem'

  if (schedule._overlapCount > 1) {
    // 重疊課程使用警告色
    const colorHex = '#F59E0B' // warning color
    const bgColor = hexToRgba(colorHex, 0.3)
    return {
      top: `${clampedTop}px`,
      height: `${Math.max(durationHeight - 2, 20)}px`,
      left,
      right,
      backgroundColor: bgColor,
      borderColor: hexToRgba(colorHex, 0.6),
    }
  }

  // 個人行程使用 color_hex 設定背景顏色
  if (schedule.is_personal_event) {
    const colorHex = schedule.color_hex || '#6366F1'
    const bgColor = hexToRgba(colorHex, 0.4)
    return {
      top: `${clampedTop}px`,
      height: `${Math.max(durationHeight - 2, 20)}px`,
      left,
      right,
      backgroundColor: bgColor,
      borderColor: hexToRgba(colorHex, 0.8),
    }
  }

  return {
    top: `${clampedTop}px`,
    height: `${Math.max(durationHeight - 2, 20)}px`,
    left,
    right,
  }
}

// 輔助函數：將 hex 顏色轉換為 RGBA
const hexToRgba = (hex: string, alpha: number): string => {
  if (!hex) return `rgba(99, 102, 241, ${alpha})` // 預設 indigo

  // 移除 # 前綴
  const hexStr = hex.replace('#', '')

  // 解析 RGB
  let r, g, b
  if (hexStr.length === 6) {
    r = parseInt(hexStr.substring(0, 2), 16)
    g = parseInt(hexStr.substring(2, 4), 16)
    b = parseInt(hexStr.substring(4, 6), 16)
  } else if (hexStr.length === 3) {
    r = parseInt(hexStr.substring(0, 1) + hexStr.substring(0, 1), 16)
    g = parseInt(hexStr.substring(1, 2) + hexStr.substring(1, 2), 16)
    b = parseInt(hexStr.substring(2, 3) + hexStr.substring(2, 3), 16)
  } else {
    return `rgba(99, 102, 241, ${alpha})` // 預設 indigo
  }

  return `rgba(${r}, ${g}, ${b}, ${alpha})`
}

// 取得卡片樣式類別
const getScheduleCardClass = (schedule: any): string => {
  if (schedule._overlapCount > 1) {
    // 重疊課程使用警告樣式
    return 'border'
  }
  if (schedule.is_personal_event) {
    const colorHex = schedule.color_hex || '#6366F1'
    const borderColor = hexToRgba(colorHex, 0.8)
    return `border ${borderColor}`
  }
  return 'bg-primary-500/20 border border-primary-500/30'
}

// 取得標題樣式類別
const getScheduleTitleClass = (schedule: any): string => {
  if (schedule._overlapCount > 1) {
    return 'text-warning-400'
  }
  if (schedule.is_personal_event) {
    return 'text-white'
  }
  return 'text-primary-400'
}

// 取得時間樣式類別
const getScheduleTimeClass = (schedule: any): string => {
  if (schedule._overlapCount > 1) {
    return 'text-warning-300'
  }
  if (schedule.is_personal_event) {
    return 'text-white/80'
  }
  return 'text-slate-400'
}
</script>
