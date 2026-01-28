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

            <HelpTooltip
              placement="bottom"
              title="週期導航"
              description="查看不同週期的排課狀況，預設顯示本週。"
              :usage="['點擊左右箭頭切換上週/下週', '可跨月、跨年查看']"
            />
          </div>
        </div>

        <!-- 快捷操作按鈕 -->
        <div class="flex items-center gap-2 ml-auto">
          <button
            @click="$emit('add-personal-event')"
            class="px-4 py-2 rounded-lg bg-purple-500/20 text-purple-400 hover:bg-purple-500/30 transition-colors flex items-center gap-2 text-sm"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            個人行程
          </button>
          <button
            @click="$emit('add-exception')"
            class="px-4 py-2 rounded-lg bg-warning-500/20 text-warning-400 hover:bg-warning-500/30 transition-colors flex items-center gap-2 text-sm"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01" />
            </svg>
            請假/調課
          </button>
        </div>
      </div>
    </div>

    <div
      class="flex-1 overflow-auto p-4"
    >
      <!-- 週曆視圖 -->
      <div class="min-w-[600px] relative" ref="calendarContainerRef">
        <!-- 表頭 -->
        <div class="grid sticky top-0 z-10 bg-slate-900/95 backdrop-blur-sm" style="grid-template-columns: 80px repeat(7, 1fr);">
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
          />
        </div>

        <!-- 課程卡片層 - 絕對定位 -->
        <div class="absolute top-0 left-0 right-0 bottom-0 pointer-events-none">
          <div
            v-for="schedule in displaySchedules"
            :key="schedule.key"
            class="absolute rounded-lg p-2 text-xs cursor-pointer hover:opacity-90 transition-opacity pointer-events-auto"
            :class="getScheduleCardClass(schedule)"
            :style="getScheduleStyle(schedule)"
            @click="$emit('select-schedule', schedule)"
          >
            <div class="font-medium truncate">
              {{ schedule.offering_name }}
            </div>
            <div class="text-slate-400 truncate">
              {{ schedule.center_name }}
            </div>
            <div class="text-slate-500 text-[10px] mt-0.5">
              {{ schedule.start_time }} - {{ schedule.end_time }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { formatDateToString } from '~/composables/useTaiwanTime'
import { nextTick, ref, computed, onMounted, onUnmounted } from 'vue'

const emit = defineEmits<{
  'update:weekStart': [value: Date]
  'select-schedule': [schedule: any]
  'add-personal-event': []
  'add-exception': []
}>()

// DOM 引用
const calendarContainerRef = ref<HTMLElement | null>(null)
const slotWidth = ref(100)

// 格子尺寸常量
const TIME_SLOT_HEIGHT = 60 // 每個時段格子的高度 (px)
const TIME_COLUMN_WIDTH = 80 // 時間列寬度 (px)

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

const getWeekStart = (date: Date): Date => {
  const d = new Date(date)
  const day = d.getDay()
  const diff = d.getDate() - day + (day === 0 ? -6 : 1)
  return new Date(d.setDate(diff))
}

const changeWeek = (delta: number) => {
  weekStart.value = getWeekStart(new Date(weekStart.value.getTime() + delta * 7 * 24 * 60 * 60 * 1000))
  emit('update:weekStart', weekStart.value)
}

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

const props = defineProps<{
  schedules: any[]
}>()

// 計算格子寬度
const calculateSlotWidth = () => {
  if (calendarContainerRef.value) {
    const containerWidth = calendarContainerRef.value.offsetWidth
    slotWidth.value = Math.max(80, (containerWidth - TIME_COLUMN_WIDTH) / 7)
  }
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

// 去重後的排課
const displaySchedules = computed(() => {
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

const formatTime = (hour: number): string => {
  return `${hour.toString().padStart(2, '0')}:00`
}

// 計算課程卡片樣式
const getScheduleStyle = (schedule: any) => {
  const { weekday, start_hour, start_minute, duration_minutes } = schedule

  // 計算水平位置 - 對齊到星期網格
  const dayIndex = weekday - 1 // 0-6
  const left = TIME_COLUMN_WIDTH + (dayIndex * slotWidth.value)

  // 計算垂直位置
  let topSlotIndex = 0
  for (let t = 0; t < start_hour; t++) {
    if (t >= 0 && t <= 3) {
      topSlotIndex++ // 0-3 時段每個都算
    } else if (t >= 9) {
      topSlotIndex++ // 9 以後的時段每個都算
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

onMounted(async () => {
  await nextTick()
  calculateSlotWidth()

  if (calendarContainerRef.value) {
    const resizeObserver = new ResizeObserver(() => {
      calculateSlotWidth()
    })

    resizeObserver.observe(calendarContainerRef.value)

    onUnmounted(() => {
      resizeObserver.disconnect()
    })
  }
})
</script>
