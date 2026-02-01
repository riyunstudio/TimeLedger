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
                  <div
                    class="absolute left-1 right-1 rounded bg-primary-500/20 border border-primary-500/30 p-1 cursor-pointer hover:bg-primary-500/30 transition-colors overflow-hidden"
                    :style="getScheduleCardStyle(schedule, time)"
                    @click="$emit('select-schedule', schedule)"
                  >
                    <div class="text-xs font-medium text-primary-400 truncate">
                      {{ schedule.offering_name }}
                    </div>
                    <div class="text-xs text-slate-400 truncate">
                      {{ schedule.start_time }}-{{ schedule.end_time }}
                    </div>
                    <div v-if="schedule.teacher_name" class="text-xs text-slate-500 truncate">
                      {{ schedule.teacher_name }}
                    </div>
                  </div>
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
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { formatDateToString } from '~/composables/useTaiwanTime'

// ============================================
// Props 定義
// ============================================

const props = defineProps<{
  schedules: any[]
  weekLabel: string
  weekStart: Date
  cardInfoType: 'teacher' | 'center'
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

const updateWidth = () => {
  if (calendarContainerRef.value) {
    containerWidth.value = calendarContainerRef.value.offsetWidth
  }
}

onMounted(() => {
  updateWidth()
  window.addEventListener('resize', updateWidth)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateWidth)
})

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

// 取得某個時間和星期幾的課程
const getSchedulesForCell = (time: number, weekday: number) => {
  return props.schedules.filter(schedule => {
    // 檢查星期幾
    if (schedule.weekday !== weekday) return false

    // 檢查時間（課程開始時間等於或早於當前時段，且結束時間晚於當前時段）
    const scheduleStartHour = schedule.start_hour

    // 解析結束時間，判斷是否跨日
    const [endHourStr, endMinuteStr] = (schedule.end_time || '23:59').split(':').map(Number)
    let scheduleEndHour = endHourStr

    // 跨日課程：結束時間早於開始時間，視為跨越到隔天
    const isOvernight = endHourStr < scheduleStartHour ||
                       (endHourStr === scheduleStartHour && endMinuteStr < schedule.start_minute)

    if (isOvernight) {
      // 跨日課程：顯示在當日格子的 00:00-時段
      // 檢查 00:00 到 24:00 這個範圍內的時段
      return time >= 0 && time < 24
    }

    // 當日課程
    const scheduleEndMinute = endMinuteStr
    return scheduleStartHour < time + 1 &&
           (scheduleEndHour > time || (scheduleEndHour === time && scheduleEndMinute > 0))
  })
}

// 取得格子樣式
const getCellClass = (time: number, weekday: number) => {
  const schedules = getSchedulesForCell(time, weekday)
  if (schedules.length > 0) {
    return 'bg-primary-500/5'
  }
  return 'hover:bg-white/5'
}

// 取得課程卡片樣式
const getScheduleCardStyle = (schedule: any, time: number) => {
  const slotHeight = 64 // 每個時段高度
  const minuteOffset = (schedule.start_minute / 60) * slotHeight
  const durationHeight = (schedule.duration_minutes / 60) * slotHeight

  // 解析結束時間，判斷是否跨日
  const [endHourStr, endMinuteStr] = (schedule.end_time || '23:59').split(':').map(Number)
  const isOvernight = endHourStr < schedule.start_hour ||
                     (endHourStr === schedule.start_hour && endMinuteStr < schedule.start_minute)

  let top: number

  if (isOvernight) {
    // 跨日課程：從 00:00 開始
    top = 0
  } else {
    // 當日課程：從開始時間計算
    top = minuteOffset
  }

  return {
    top: `${top}px`,
    height: `${Math.max(durationHeight - 2, 20)}px`,
  }
}
</script>
