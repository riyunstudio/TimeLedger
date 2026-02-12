<template>
  <div
    class="absolute rounded-lg p-2 text-xs cursor-pointer hover:opacity-90 transition-opacity pointer-events-auto"
    :class="cardClass"
    :style="cardStyle"
    @click="$emit('click')"
  >
    <div class="font-medium truncate" :class="titleClass">
      <!-- PLANNED 狀態標籤 -->
      <span v-if="schedule.status === 'PLANNED'" class="text-[10px] bg-amber-500/30 px-1.5 py-0.5 rounded mr-1 text-amber-300">預計</span>
      <!-- SUSPENDED 狀態標籤 -->
      <span v-else-if="schedule.status === 'SUSPENDED'" class="text-[10px] bg-slate-500/30 px-1.5 py-0.5 rounded mr-1 text-slate-300">停課</span>
      <!-- ARCHIVED 狀態標籤 -->
      <span v-else-if="schedule.status === 'ARCHIVED'" class="text-[10px] bg-slate-600/30 px-1.5 py-0.5 rounded mr-1 text-slate-400">歸檔</span>
      {{ schedule.offering_name }}
    </div>
    <div v-if="cardInfoType === 'teacher'" class="text-slate-400 truncate">
      {{ schedule.teacher_name }}
    </div>
    <div v-else class="truncate" :class="timeClass">
      {{ schedule.center_name }}
    </div>
    <div class="text-xs mt-0.5" :class="timeClass">
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
  positionStyle?: Record<string, string>
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
// 計算屬性（Vue computed 自動追蹤依賴）
// ============================================

const cardClass = computed(() => {
  const schedule = props.schedule
  if (!schedule) return ''

  // 個人行程使用 color_hex 設定背景顏色
  if (schedule.is_personal_event) {
    const colorHex = schedule.color_hex || '#6366F1'
    // 將 hex 顏色轉換為 RGBA 半透明背景
    const bgColor = hexToRgba(colorHex, 0.3)
    const borderColor = hexToRgba(colorHex, 0.6)
    return `border ${borderColor} text-white`
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

  // 根據狀態返回對應樣式
  const status = schedule.status || ''
  switch (status) {
    case 'PLANNED':
      // 預計課程：虛線邊框 + 斜紋背景
      return 'schedule-card-planned border border-dashed border-amber-500/50 text-amber-300'
    case 'CONFIRMED':
      // 已開課：實線邊框（預設樣式）
      return 'bg-primary-500/30 border border-primary-500/50 text-white'
    case 'SUSPENDED':
      // 停課：灰色樣式
      return 'bg-slate-500/30 border border-slate-500/50 text-slate-300'
    case 'ARCHIVED':
      // 已歸檔：淡化樣式
      return 'bg-slate-600/30 border border-slate-600/30 text-slate-400 opacity-60'
    default:
      // 預設樣式（包含空字串情況）
      return 'bg-slate-700/80 border border-white/10 text-white'
  }
})

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

const titleClass = computed(() => {
  const schedule = props.schedule

  // 個人行程使用白色文字
  if (schedule.is_personal_event) {
    return 'text-white'
  }

  return ''
})

// 時間文字樣式
const timeClass = computed(() => {
  const schedule = props.schedule

  // 個人行程使用白色半透明文字
  if (schedule.is_personal_event) {
    return 'text-white/80'
  }

  return 'text-slate-400'
})

// 組合樣式：位置樣式 + 動態背景顏色
const cardStyle = computed(() => {
  const style = props.positionStyle || {}

  // 個人行程使用 color_hex 設定背景顏色
  if (props.schedule?.is_personal_event) {
    const colorHex = props.schedule.color_hex || '#6366F1'
    const bgColor = hexToRgba(colorHex, 0.4)
    return {
      ...style,
      backgroundColor: bgColor,
      borderColor: hexToRgba(colorHex, 0.8),
    }
  }

  return style
})
</script>

<style scoped>
.schedule-card-planned {
  background-color: rgba(51, 65, 85, 0.6); /* slate-700 with opacity */
  background-image: repeating-linear-gradient(
    45deg,
    transparent,
    transparent 10px,
    rgba(255, 255, 255, 0.05) 10px,
    rgba(255, 255, 255, 0.05) 20px
  );
}
</style>
