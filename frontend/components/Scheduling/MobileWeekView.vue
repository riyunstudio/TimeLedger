<template>
  <div class="lg:hidden">
    <!-- 週選擇器（手機版） -->
    <div class="flex items-center justify-between mb-4">
      <button
        @click="$emit('change-week', -1)"
        class="p-2 rounded-lg hover:bg-white/10 transition-colors"
      >
        <svg class="w-5 h-5 text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>
      <h2 class="text-base font-semibold text-white">{{ weekLabel }}</h2>
      <button
        @click="$emit('change-week', 1)"
        class="p-2 rounded-lg hover:bg-white/10 transition-colors"
      >
        <svg class="w-5 h-5 text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>
    </div>

    <!-- 日曆卡片列表 -->
    <div class="space-y-3">
      <div
        v-for="day in weekDays"
        :key="day.value"
        class="bg-slate-800/50 rounded-lg border border-slate-700 overflow-hidden"
      >
        <!-- 日期標題 -->
        <div class="px-4 py-2 bg-slate-700/30 border-b border-slate-700 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <span class="text-sm font-medium text-white">{{ day.name }}</span>
            <span class="text-xs text-slate-400">{{ getDayDate(day.value) }}</span>
          </div>
          <span v-if="getDayScheduleCount(day.value) > 0" class="text-xs px-2 py-0.5 bg-primary-500/20 text-primary-400 rounded-full">
            {{ getDayScheduleCount(day.value) }} 堂課
          </span>
        </div>

        <!-- 當日課程列表 -->
        <div v-if="getDaySchedules(day.value).length > 0" class="divide-y divide-slate-700/50">
          <div
            v-for="schedule in getDaySchedules(day.value)"
            :key="schedule.id"
            class="p-3 hover:bg-white/5 transition-colors cursor-pointer"
            @click="$emit('select-schedule', schedule)"
          >
            <div class="flex items-start gap-3">
              <!-- 時間 -->
              <div class="w-14 flex-shrink-0">
                <div class="text-sm font-medium text-primary-400">{{ schedule.start_time }}</div>
                <div class="text-xs text-slate-500">{{ schedule.end_time }}</div>
              </div>

              <!-- 課程資訊 -->
              <div class="flex-1 min-w-0">
                <div class="font-medium text-white truncate">{{ schedule.offering_name }}</div>
                <div class="flex items-center gap-2 mt-1">
                  <span v-if="schedule.teacher_name" class="text-xs text-slate-400">{{ schedule.teacher_name }}</span>
                  <span v-if="schedule.room_name" class="text-xs text-slate-500">· {{ schedule.room_name }}</span>
                </div>
              </div>

              <!-- 狀態指示 -->
              <div v-if="schedule.has_exception" class="flex-shrink-0">
                <span
                  v-if="schedule.exception_type === 'CANCEL'"
                  class="text-xs px-2 py-1 bg-critical-500/20 text-critical-400 rounded"
                >
                  已取消
                </span>
                <span
                  v-else-if="schedule.exception_type === 'RESCHEDULE'"
                  class="text-xs px-2 py-1 bg-warning-500/20 text-warning-400 rounded"
                >
                  調課
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- 無課程提示 -->
        <div v-else class="p-6 text-center">
          <p class="text-sm text-slate-500">無課程安排</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
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
}>()

// ============================================
// Emits 定義
// ============================================

defineEmits<{
  'change-week': [delta: number]
  'select-schedule': [schedule: any]
}>()

// ============================================
// 工具函數
// ============================================

import { formatDate } from '~/composables/useTaiwanTime'

const weekDays = [
  { value: 1, name: '週一' },
  { value: 2, name: '週二' },
  { value: 3, name: '週三' },
  { value: 4, name: '週四' },
  { value: 5, name: '週五' },
  { value: 6, name: '週六' },
  { value: 7, name: '週日' },
]

// 取得某個星期幾的日期（格式：YYYY/MM/DD）
const getDayDate = (weekday: number): string => {
  const date = new Date(props.weekStart)
  date.setDate(date.getDate() + (weekday - 1))
  return formatDate(date)
}

// 取得某個星期幾的所有課程
const getDaySchedules = (weekday: number) => {
  return props.schedules.filter((schedule: any) => schedule.weekday === weekday)
}

// 取得某個星期幾的課程數量
const getDayScheduleCount = (weekday: number): number => {
  return getDaySchedules(weekday).length
}
</script>
