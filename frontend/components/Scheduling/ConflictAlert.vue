<template>
  <div class="p-4">
    <div
      class="rounded-xl p-4 mb-4"
      :class="alertClass"
    >
      <div class="flex items-start gap-3">
        <Icon
          :icon="iconName"
          size="xl"
          class="flex-shrink-0 mt-0.5"
          :class="iconClass"
        />
        <div class="flex-1">
          <h4 class="font-medium mb-2" :class="titleClass">
            {{ title }}
          </h4>
          <p v-if="description" class="text-sm text-slate-400 mb-3">
            {{ description }}
          </p>
          <ul v-if="conflictMessages.length > 0" class="space-y-1">
            <li
              v-for="(message, index) in conflictMessages"
              :key="index"
              class="text-sm text-slate-300 flex items-center gap-2"
            >
              <span class="w-1.5 h-1.5 rounded-full flex-shrink-0" :class="dotClass"></span>
              {{ message }}
            </li>
          </ul>
        </div>
      </div>
    </div>
    <button
      @click="handleDismiss"
      class="btn-secondary w-full py-3 rounded-xl font-medium"
    >
      我知道了
    </button>
  </div>
</template>

<script setup lang="ts">
interface Props {
  /** 衝突訊息列表 */
  conflictMessages: string[]
  /** 衝突類型：硬衝突（不可覆寫）或緩衝衝突（可覆寫） */
  conflictType?: 'hard' | 'buffer'
  /** 自訂標題，預設根據 conflictType 自動設定 */
  customTitle?: string
  /** 自訂描述文字 */
  customDescription?: string
}

const props = withDefaults(defineProps<Props>(), {
  conflictType: 'hard',
  conflictMessages: () => [],
})

const emit = defineEmits<{
  dismiss: []
}>()

/** 衝突標題 */
const title = computed(() => {
  if (props.customTitle) {
    return props.customTitle
  }
  return props.conflictType === 'hard' ? '排課時間衝突' : '緩衝時間不足'
})

/** 衝突描述 */
const description = computed(() => {
  if (props.customDescription) {
    return props.customDescription
  }
  return props.conflictType === 'hard'
    ? '以下時間已有排課，請選擇其他時間或教室：'
    : '與前一堂課的間隔時間不足，是否仍要儲存？'
})

/** 警示區域樣式 */
const alertClass = computed(() => {
  return props.conflictType === 'hard'
    ? 'bg-critical-500/10 border border-critical-500/30'
    : 'bg-warning-500/10 border border-warning-500/30'
})

/** 標題樣式 */
const titleClass = computed(() => {
  return props.conflictType === 'hard' ? 'text-critical-500' : 'text-warning-500'
})

/** 圖示名稱 */
const iconName = computed(() => {
  return props.conflictType === 'hard' ? 'warning' : 'warning'
})

/** 圖示樣式 */
const iconClass = computed(() => {
  return props.conflictType === 'hard' ? 'text-critical-500' : 'text-warning-500'
})

/** 項目圓點樣式 */
const dotClass = computed(() => {
  return props.conflictType === 'hard' ? 'bg-critical-500' : 'bg-warning-500'
})

/** 處理關閉 */
const handleDismiss = () => {
  emit('dismiss')
}
</script>
