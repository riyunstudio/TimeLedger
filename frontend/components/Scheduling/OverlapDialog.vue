<template>
  <Teleport to="body">
    <div
      v-if="visible"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
      @click.self="$emit('close')"
    >
      <div class="glass-card w-full max-w-sm">
        <div class="p-4 border-b border-white/10">
          <h3 class="text-lg font-semibold text-white">選擇課程</h3>
          <p v-if="timeSlot" class="text-sm text-slate-400 mt-1">
            {{ getWeekdayText(timeSlot.weekday) }}
            {{ timeSlot.start_hour.toString().padStart(2, '0') }}:{{ timeSlot.start_minute.toString().padStart(2, '0') }}
          </p>
        </div>
        <div class="max-h-96 overflow-y-auto">
          <div
            v-for="schedule in schedules"
            :key="schedule.id"
            class="p-4 border-b border-white/5 hover:bg-white/5 cursor-pointer transition-colors"
            @click="$emit('select', schedule)"
          >
            <div class="flex items-center justify-between">
              <div>
                <div class="font-medium text-white">{{ schedule.offering_name }}</div>
                <div class="text-sm text-slate-400 mt-1">
                  {{ schedule.start_time }} - {{ schedule.end_time }}
                </div>
                <div v-if="cardInfoType === 'teacher'" class="text-xs text-slate-500 mt-1">
                  {{ schedule.teacher_name }}
                </div>
                <div v-else class="text-xs text-slate-500 mt-1">
                  {{ schedule.center_name }}
                </div>
              </div>
              <div v-if="schedule.room_name" class="text-xs text-slate-500">
                {{ schedule.room_name }}
              </div>
            </div>
          </div>
        </div>
        <div class="p-4 border-t border-white/10">
          <button
            @click="$emit('close')"
            class="w-full px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors"
          >
            取消
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
// ============================================
// Props 定義
// ============================================

const props = defineProps<{
  // 是否顯示
  visible: boolean
  // 重疊的課程列表
  schedules: any[]
  // 時間槽資訊
  timeSlot: { weekday: number; start_hour: number; start_minute: number } | null
  // 卡片顯示類型
  cardInfoType: 'teacher' | 'center'
}>()

// ============================================
// Emits 定義
// ============================================

defineEmits<{
  close: []
  select: [schedule: any]
}>()

// ============================================
// 方法
// ============================================

const getWeekdayText = (weekday: number): string => {
  const days = ['週日', '週一', '週二', '週三', '週四', '週五', '週六']
  return days[weekday - 1] || ''
}
</script>
