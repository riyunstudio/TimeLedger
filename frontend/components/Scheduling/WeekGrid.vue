<template>
  <!-- 週曆視圖 - 桌面版 -->
  <div id="weekgrid-container" class="hidden lg:block min-w-[800px] relative" ref="calendarContainerRef">
    <!-- 表頭 -->
    <div class="grid sticky top-0 z-10 bg-slate-900/95 backdrop-blur-sm" style="grid-template-columns: 80px repeat(7, 1fr);">
      <div class="p-2 border-b border-white/10 text-center">
        <span class="text-xs text-slate-400">{{ $t('schedule.timeSlot') }}</span>
      </div>
      <div
        v-for="day in weekDays"
        :key="day.value"
        class="p-2 border-b border-white/10 text-center"
      >
        <span class="text-sm font-medium text-slate-100">{{ day.name }}</span>
        <span class="block text-xs text-slate-400 mt-0.5">{{ day.date }}</span>
      </div>
    </div>

    <!-- 時間列和網格區域 -->
    <div class="relative">
      <!-- 時間格子 -->
      <div
        v-for="time in timeSlots"
        :key="time"
        class="grid relative z-0"
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
          class="p-0 min-h-[60px] border-b border-white/5 border-r relative z-0"
          :class="cellClassMap[`${time}-${day.value}`]"
          @dragenter="$emit('drag-enter', time, day.value)"
          @dragleave="$emit('drag-leave')"
          @dragover.prevent
        />
      </div>

      <!-- 課程卡片層 -->
      <div class="absolute top-0 left-0 right-0 bottom-0 pointer-events-none z-10" style="height: 1440px;">
        <!-- 1440px = 24小時 * 60px 每小時 -->
        <!-- 直接渲染卡片，移除 DynamicScroller 以避免定位問題 -->
        <template v-for="item in schedulesWithOverlapData" :key="item.key">
          <template v-if="item.is_personal_event">
            <!-- 個人行程 -->
            <ScheduleCard
              :schedule="item"
              :position-style="item._scheduleStyle"
              :card-info-type="cardInfoType"
              @click="$emit('select-schedule', item)"
            />
          </template>
          <template v-else>
            <!-- 中心課程 -->
            <template v-if="item._overlapCount === 1">
              <ScheduleCard
                :schedule="item"
                :position-style="item._scheduleStyle"
                :card-info-type="cardInfoType"
                @click="$emit('select-schedule', item)"
              />
            </template>
            <!-- 重疊指示器 -->
            <template v-else-if="item._overlapCount > 1 && item._isFirstInOverlap">
              <div
                class="absolute rounded-lg bg-warning-500/20 border border-warning-500/50 p-2 text-xs cursor-pointer hover:bg-warning-500/30 transition-opacity pointer-events-auto"
                :style="item._scheduleStyle"
                @click="$emit('overlap-click', item)"
              >
                <div class="flex items-center justify-center h-full">
                  <span class="text-warning-400 font-bold text-lg">
                    {{ item._overlapCount }}
                  </span>
                  <span class="text-warning-300 ml-1 text-xs">{{ $t('schedule.courseSession') }}</span>
                </div>
              </div>
            </template>
          </template>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, shallowRef, markRaw, watch, onMounted, onUnmounted, nextTick } from 'vue'
import ScheduleCard from './ScheduleCard.vue'

// ============================================
// Props 定義
// ============================================

const props = defineProps<{
  // 課程資料
  schedules: any[]
  // 週標籤
  weekLabel: string
  // 週起始日期
  weekStart: Date
  // 卡片顯示類型
  cardInfoType: 'teacher' | 'center'
  // 驗證結果
  validationResults: Record<string, any>
  // 槽寬度（由父組件傳入，但自己也會計算）
  slotWidth: number
}>()

// ============================================
// Emits 定義
// ============================================

const emit = defineEmits<{
  'drag-enter': [time: number, weekday: number]
  'drag-leave': []
  'select-schedule': [schedule: any]
  'overlap-click': [schedule: any]
}>()

// ============================================
// 常量定義
// ============================================

const TIME_SLOT_HEIGHT = 60 // 每個時段格子的高度 (px)
const TIME_COLUMN_WIDTH = 80 // 時間列寬度 (px)

// 時間段（連續顯示所有時段）
const timeSlots = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23]

// 星期幾（需要傳入週起始日期來計算日期）
const weekDays = computed(() => {
  try {
    if (!props.weekStart) {
      return [
        { value: 1, name: '週一', date: '' },
        { value: 2, name: '週二', date: '' },
        { value: 3, name: '週三', date: '' },
        { value: 4, name: '週四', date: '' },
        { value: 5, name: '週五', date: '' },
        { value: 6, name: '週六', date: '' },
        { value: 7, name: '週日', date: '' },
      ]
    }

    const days = ['週日', '週一', '週二', '週三', '週四', '週五', '週六']
    const result = []

    for (let i = 0; i < 7; i++) {
      const date = new Date(props.weekStart)
      date.setDate(date.getDate() + i)
      const month = date.getMonth() + 1
      const day = date.getDate()

      // 週一到週日對應 value 1-7，直接使用迴圈索引 i
      // i=0 對應週一 (Monday)，i=6 對應週日 (Sunday)
      const value = i + 1

      result.push({
        value,
        name: days[date.getDay()],
        date: `${month}/${day}`
      })
    }

    return result
  } catch (error) {
    console.error('[WeekGrid] Error computing weekDays:', error)
    return []
  }
})

// ============================================
// DOM 引用
// ============================================

const calendarContainerRef = ref<HTMLElement | null>(null)

// ============================================
// 本地 slotWidth 計算（不再依賴父組件傳入）
// ============================================

// 使用 ref 來存儲本地計算的 slotWidth
const localSlotWidth = ref<number>(120)

// 計算 slotWidth 的函數
const calculateSlotWidth = () => {
  if (!calendarContainerRef.value) {
    console.warn('[WeekGrid] calendarContainerRef is null')
    return
  }

  const container = calendarContainerRef.value
  const containerWidth = container.offsetWidth

  // 時間列寬度固定為 80px
  const timeColumnWidth = 80
  const calculatedSlotWidth = Math.max(80, (containerWidth - timeColumnWidth) / 7)

  if (calculatedSlotWidth !== localSlotWidth.value) {
    localSlotWidth.value = calculatedSlotWidth
    console.log('[WeekGrid] calculateSlotWidth:', {
      containerWidth,
      calculatedSlotWidth,
      timeColumnWidth
    })
    // 清除樣式快取，讓卡片重新計算位置
    styleCache.clear()
    overlapDataCache.clear()
  }
}

// ResizeObserver 引用
let resizeObserver: ResizeObserver | null = null

// ============================================
// 快取儲存（使用淺層引用避免深層響應式開銷）
// ============================================

// 課程樣式快取
const styleCache = new Map<string, Record<string, string>>()

// 重疊數據快取
const overlapDataCache = new Map<string, { count: number; firstId: number }>()

// ============================================
// 工具函數
// ============================================

const formatTime = (hour: number): string => {
  return `${hour.toString().padStart(2, '0')}:00`
}

// 計算課程卡片樣式（帶快取）
const getScheduleStyle = (schedule: any): Record<string, string> => {
  // 使用本地計算的 slotWidth
  const currentSlotWidth = localSlotWidth.value

  // 確保 weekday 從 date 正確計算（而不是依賴後端可能錯誤的 weekday 欄位）
  let weekday: number
  if (schedule.date) {
    const itemDate = new Date(schedule.date + 'T00:00:00+08:00')
    weekday = itemDate.getDay() === 0 ? 7 : itemDate.getDay()
  } else {
    weekday = schedule.weekday || 1 // 預設為週一
  }

  // 使用計算出的 weekday 來建立 cacheKey
  const cacheKey = `${schedule.id}-${weekday}-${schedule.start_time}-${currentSlotWidth}`
  const cached = styleCache.get(cacheKey)
  if (cached) {
    console.log('[WeekGrid] getScheduleStyle - 使用快取:', {
      id: schedule.id,
      cacheKey,
      style: cached
    })
    return cached
  }

  const { start_hour, start_minute, duration_minutes, end_time } = schedule

  // 計算垂直位置 - 時間格子是連續的 0-23 小時
  const slotHeight = TIME_SLOT_HEIGHT

  // 計算 top 和 height
  let top: number
  let height: number

  if (schedule.is_cross_day_part && schedule.start_time === '00:00') {
    // 跨日課程的結束部分（00:00 開始）
    top = 0
    height = (duration_minutes / 60) * slotHeight
  } else {
    // 一般課程
    const baseTop = start_hour * slotHeight
    const minuteOffset = (start_minute / 60) * slotHeight
    top = baseTop + minuteOffset
    height = (duration_minutes / 60) * slotHeight
  }

  // 計算水平位置 - 對齊到星期網格
  // weekday: 1-7 (週一到週日)
  const dayIndex = weekday - 1 // 0-6
  const left = TIME_COLUMN_WIDTH + (dayIndex * currentSlotWidth)

  // 計算寬度（略小於格子寬度以留邊距）
  const width = currentSlotWidth - 4

  const style = {
    left: `${left}px`,
    top: `${top}px`,
    width: `${width}px`,
    height: `${height}px`,
  }

  // 使用 markRaw 避免 Vue 將其轉為響應式Proxy
  styleCache.set(cacheKey, markRaw(style))
  return style
}

// 計算並快取重疊數據
const computeOverlapData = (schedules: any[]) => {
  const countMap: Record<string, number> = {}
  const firstIdMap: Record<string, number> = {}

  // 第一遍：計算每個時段的課程數量和第一個 ID
  for (const schedule of schedules) {
    // 確保 weekday 從 date 正確計算
    let weekday: number
    if (schedule.date) {
      const itemDate = new Date(schedule.date + 'T00:00:00+08:00')
      weekday = itemDate.getDay() === 0 ? 7 : itemDate.getDay()
    } else {
      weekday = schedule.weekday || 1
    }

    const key = `${weekday}-${schedule.start_hour}-${schedule.start_minute}`
    countMap[key] = (countMap[key] || 0) + 1

    // 記錄最小的 ID 作為第一個
    if (!firstIdMap[key] || schedule.id < firstIdMap[key]) {
      firstIdMap[key] = schedule.id
    }
  }

  // 儲存到快取
  for (const key in countMap) {
    overlapDataCache.set(key, {
      count: countMap[key],
      firstId: firstIdMap[key]
    })
  }
}

// ============================================
// 計算屬性
// ============================================

// 去重後的課程（保持引用穩定性）
const uniqueSchedules = computed(() => {
  const seen = new Set<string>()
  const result: any[] = []

  for (const schedule of props.schedules) {
    const key = `${schedule.id}-${schedule.weekday}-${schedule.start_time}`
    if (!seen.has(key)) {
      seen.add(key)
      result.push(schedule)
    }
  }

  return result
})

// 預計算每個課程的重疊數據和樣式
const schedulesWithOverlapData = computed(() => {
  const schedules = uniqueSchedules.value

  // 計算重疊數據
  computeOverlapData(schedules)

  // 為每個課程添加計算後的數據
  return schedules.map(schedule => {
    const overlapKey = `${schedule.weekday}-${schedule.start_hour}-${schedule.start_minute}`
    const overlapData = overlapDataCache.get(overlapKey)

    const style = getScheduleStyle(schedule)

    console.log('[WeekGrid] Style calculation:', {
      id: schedule.id,
      weekday: schedule.weekday,
      start_hour: schedule.start_hour,
      start_minute: schedule.start_minute,
      duration_minutes: schedule.duration_minutes,
      date: schedule.date,
      offering_name: schedule.offering_name,
      style: style,
      localSlotWidth: localSlotWidth.value,
      weekDays: weekDays.value.map(d => ({ value: d.value, name: d.name, date: d.date }))
    })

    // 使用淺層複製，避免修改原始資料
    const enhanced = {
      ...schedule,
      _overlapCount: overlapData?.count || 1,
      _isFirstInOverlap: overlapData?.firstId === schedule.id,
      _scheduleStyle: style
    }

    return enhanced
  })
})

// 預計算網格樣式（避免每個格子渲染時重複計算）
const cellClassMap = computed(() => {
  try {
    const weekDaysValue = weekDays.value
    // 防護：確保 weekDays 是可迭代的陣列
    if (!Array.isArray(weekDaysValue) || weekDaysValue.length === 0) {
      return {}
    }

    const classMap: Record<string, string> = {}

    for (const time of timeSlots) {
      for (const day of weekDaysValue) {
        const key = `${time}-${day.value}`
        const validation = props.validationResults[key]

        if (validation?.valid === false) {
          classMap[key] = 'bg-critical-500/10 border-critical-500/50'
        } else if (validation?.warning) {
          classMap[key] = 'bg-warning-500/10 border-warning-500/50'
        } else if (validation?.valid === true) {
          classMap[key] = 'bg-success-500/10 border-success-500/50'
        } else {
          classMap[key] = 'hover:bg-white/5'
        }
      }
    }

    return classMap
  } catch (error) {
    console.error('[WeekGrid] Error computing cellClassMap:', error)
    return {}
  }
})

// ============================================
// 監聽 props 變化，清除快取
// ============================================

// 當 schedules 變化時，清除樣式快取
// 注意：slotWidth 現在由本地計算，不再監聽 props.slotWidth
const lastSchedulesSnapshot = shallowRef<any[]>([])

// 調試：監控傳入的 props
watch([() => props.schedules, () => props.weekStart], ([newSchedules, newWeekStart], [oldSchedules, oldWeekStart]) => {
  console.log('[WeekGrid] Props changed:', {
    schedulesCount: newSchedules?.length || 0,
    weekStart: newWeekStart,
    firstSchedule: newSchedules?.[0] ? {
      id: newSchedules[0].id,
      weekday: newSchedules[0].weekday,
      start_hour: newSchedules[0].start_hour,
      start_minute: newSchedules[0].start_minute,
      start_time: newSchedules[0].start_time
    } : null
  })
}, { immediate: true, deep: false })

watch(
  () => props.schedules,
  (newSchedules, oldSchedules) => {
    const snapshot = newSchedules.map(s => `${s.id}-${s.weekday}-${s.start_time}-${s.duration_minutes}`).join(',')
    const lastSnapshot = lastSchedulesSnapshot.value.join(',')

    if (snapshot !== lastSnapshot) {
      styleCache.clear()
      overlapDataCache.clear()
      console.log('[WeekGrid] Cache cleared due to schedules change')
      lastSchedulesSnapshot.value = newSchedules.map(s => `${s.id}-${s.weekday}-${s.start_time}-${s.duration_minutes}`)
    }
  },
  { deep: false }
)

// ============================================
// 生命週期
// ============================================

onMounted(async () => {
  // 等待 DOM 渲染完成
  await nextTick()
  await nextTick()

  // 計算初始 slotWidth
  calculateSlotWidth()

  // 設置 ResizeObserver 監控容器大小變化
  if (calendarContainerRef.value) {
    resizeObserver = new ResizeObserver(() => {
      calculateSlotWidth()
    })
    resizeObserver.observe(calendarContainerRef.value)
  }
})

onUnmounted(() => {
  if (resizeObserver) {
    resizeObserver.disconnect()
    resizeObserver = null
  }
})
</script>
