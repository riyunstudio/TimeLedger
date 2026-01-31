<template>
  <div
    class="absolute rounded-lg p-2 text-xs cursor-pointer hover:opacity-90 transition-opacity pointer-events-auto"
    :class="cardClass"
    :style="style"
    @click="$emit('click')"
  >
    <div class="font-medium truncate" :class="titleClass">
      {{ schedule.offering_name }}
    </div>
    <div v-if="cardInfoType === 'teacher'" class="text-slate-400 truncate">
      {{ schedule.teacher_name }}
    </div>
    <div v-else class="text-slate-400 truncate">
      {{ schedule.center_name }}
    </div>
    <div class="text-slate-500 text-[10px] mt-0.5">
      {{ schedule.start_time }} - {{ schedule.end_time }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

// ============================================
// Props 定義
// ============================================

const props = defineProps<{
  // 課程資料
  schedule: any
  // 卡片樣式（從外部傳入的位置樣式）
  style?: Record<string, string>
  // 卡片顯示類型
  cardInfoType: 'teacher' | 'center'
}>()

// ============================================
// Emits 定義
// ============================================

defineEmits<{
  click: []
}>()

// ============================================
// 計算屬性
// ============================================

const cardClass = computed(() => {
  const schedule = props.schedule
  if (!schedule) return ''

  // 個人行程使用不同的樣式
  if (schedule.is_personal_event) {
    const baseColor = schedule.color_hex || '#a855f7' // 紫色預設
    return 'border border-white/20'
  }

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
})

const titleClass = computed(() => {
  const schedule = props.schedule

  // 個人行程使用白色文字
  if (schedule.is_personal_event) {
    return 'text-white'
  }

  return ''
})
</script>
