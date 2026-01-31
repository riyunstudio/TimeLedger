<template>
  <!-- 週曆視圖 - 桌面版 -->
  <div class="hidden lg:block min-w-[800px] relative" ref="calendarContainerRef">
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
          :class="getCellClass(time, day.value)"
          @dragenter="$emit('drag-enter', time, day.value)"
          @dragleave="$emit('drag-leave')"
          @dragover.prevent
        />
      </div>

      <!-- 課程卡片層 -->
      <div class="absolute top-0 left-0 right-0 bottom-0 pointer-events-none z-10">
        <DynamicScroller
          :items="virtualizedSchedules"
          :min-item-size="60"
          class="h-full"
          key-field="key"
          v-if="schedules.length > 0"
        >
          <template #default="{ item, index, active }">
            <DynamicScrollerItem
              :item="item"
              :active="active"
              :size-dependencies="[
                item.offering_name,
                item.start_time,
                item.end_time,
                item.teacher_name,
                item.has_exception
              ]"
              :data-index="index"
            >
              <template v-if="item.is_personal_event">
                <!-- 個人行程 -->
                <ScheduleCard
                  :schedule="item"
                  :style="getScheduleStyle(item)"
                  :card-info-type="cardInfoType"
                  @click="$emit('select-schedule', item)"
                />
              </template>
              <template v-else>
                <!-- 中心課程 -->
                <template v-if="getOverlapCount(item) === 1">
                  <ScheduleCard
                    :schedule="item"
                    :style="getScheduleStyle(item)"
                    :card-info-type="cardInfoType"
                    @click="$emit('select-schedule', item)"
                  />
                </template>
                <!-- 重疊指示器 -->
                <template v-else-if="getOverlapCount(item) > 1 && isFirstInOverlap(item)">
                  <div
                    class="absolute rounded-lg bg-warning-500/20 border border-warning-500/50 p-2 text-xs cursor-pointer hover:bg-warning-500/30 transition-opacity pointer-events-auto"
                    :style="getScheduleStyle(item)"
                    @click="$emit('overlap-click', item)"
                  >
                    <div class="flex items-center justify-center h-full">
                      <span class="text-warning-400 font-bold text-lg">
                        {{ getOverlapCount(item) }}
                      </span>
                      <span class="text-warning-300 ml-1 text-xs">堂課程</span>
                    </div>
                  </div>
                </template>
              </template>
            </DynamicScrollerItem>
          </template>
        </DynamicScroller>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { DynamicScroller, DynamicScrollerItem } from 'vue-virtual-scroller'
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css'
import ScheduleCard from './ScheduleCard.vue'

// ============================================
// Props 定義
// ============================================

const props = defineProps<{
  // 課程資料
  schedules: any[]
  // 週標籤
  weekLabel: string
  // 卡片顯示類型
  cardInfoType: 'teacher' | 'center'
  // 驗證結果
  validationResults: Record<string, any>
  // 槽寬度
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

// 時間段
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

// ============================================
// DOM 引用
// ============================================

const calendarContainerRef = ref<HTMLElement | null>(null)

// ============================================
// 計算屬性
// ============================================

// 去重後的課程
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

// 虛擬滾動用的課程列表
const virtualizedSchedules = computed(() => {
  const schedules = uniqueSchedules.value
  // 如果課程數量少於 50 筆，不需要虛擬滾動，直接返回原陣列
  if (schedules.length < 50) {
    return schedules
  }
  return schedules
})

// 重疊數量映射
const overlapCountMap = computed(() => {
  const countMap: Record<string, number> = {}

  for (const schedule of uniqueSchedules.value) {
    const key = `${schedule.weekday}-${schedule.start_hour}-${schedule.start_minute}`
    countMap[key] = (countMap[key] || 0) + 1
  }

  return countMap
})

// ============================================
// 方法
// ============================================

const formatTime = (hour: number): string => {
  return `${hour.toString().padStart(2, '0')}:00`
}

const getCellClass = (time: number, weekday: number): string => {
  const key = `${time}-${weekday}`
  const validation = props.validationResults[key]

  if (validation?.valid === false) {
    return 'bg-critical-500/10 border-critical-500/50'
  } else if (validation?.warning) {
    return 'bg-warning-500/10 border-warning-500/50'
  } else if (validation?.valid === true) {
    return 'bg-success-500/10 border-success-500/50'
  }

  return 'hover:bg-white/5'
}

const getOverlapCount = (schedule: any) => {
  const key = `${schedule.weekday}-${schedule.start_hour}-${schedule.start_minute}`
  return overlapCountMap.value[key] || 1
}

const isFirstInOverlap = (schedule: any) => {
  const key = `${schedule.weekday}-${schedule.start_hour}-${schedule.start_minute}`
  const allAtSameTime = uniqueSchedules.value.filter(s =>
    `${s.weekday}-${s.start_hour}-${s.start_minute}` === key
  )
  // 按 id 排序，返回第一個
  allAtSameTime.sort((a, b) => a.id - b.id)
  return allAtSameTime[0]?.id === schedule.id
}

// 計算課程卡片樣式
const getScheduleStyle = (schedule: any) => {
  const { weekday, start_hour, start_minute, duration_minutes } = schedule

  // 計算水平位置 - 對齊到星期網格
  const dayIndex = weekday - 1 // 0-6
  const left = TIME_COLUMN_WIDTH + (dayIndex * props.slotWidth)

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

  // 計算高度（持續分鐘數轉像素）
  const height = (duration_minutes / 60) * slotHeight

  // 計算寬度（略小於格子寬度以留邊距）
  const width = props.slotWidth - 4

  return {
    left: `${left}px`,
    top: `${top}px`,
    width: `${width}px`,
    height: `${height}px`,
  }
}
</script>
