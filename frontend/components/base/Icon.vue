<template>
  <!-- user - 用戶頭像 -->
  <svg v-if="icon === 'user'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" :class="computedClass" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <path d="M16 7a4 4 0 1 1-8 0 4 4 0 0 1 8 0z" />
    <path d="M12 14a7 7 0 0 0-7 7h14a7 7 0 0 0-7-7z" />
  </svg>

  <!-- lock - 鎖頭 -->
  <svg v-else-if="icon === 'lock'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" :class="computedClass" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
    <path d="M7 11V7a5 5 0 0 1 10 0v4" />
  </svg>

  <!-- chat - 對話訊息 -->
  <svg v-else-if="icon === 'chat'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" :class="computedClass" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
  </svg>

  <!-- clock - 時鐘 -->
  <svg v-else-if="icon === 'clock'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" :class="computedClass" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <circle cx="12" cy="12" r="10" />
    <path d="M12 6v6l4 2" />
  </svg>

  <!-- spinner - 載入中 -->
  <svg v-else-if="icon === 'spinner'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" :class="computedClass" fill="none" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round">
    <line x1="12" y1="2" x2="12" y2="6" />
    <line x1="12" y1="18" x2="12" y2="22" />
    <line x1="4.93" y1="4.93" x2="7.76" y2="7.76" />
    <line x1="16.24" y1="16.24" x2="19.07" y2="19.07" />
    <line x1="2" y1="12" x2="6" y2="12" />
    <line x1="18" y1="12" x2="22" y2="12" />
    <line x1="4.93" y1="19.07" x2="7.76" y2="16.24" />
    <line x1="16.24" y1="7.76" x2="19.07" y2="4.93" />
  </svg>

  <!-- chevron-left - 左箭頭 -->
  <svg v-else-if="icon === 'chevron-left'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" :class="computedClass" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <path d="M15 18l-6-6 6-6" />
  </svg>

  <!-- chevron-right - 右箭頭 -->
  <svg v-else-if="icon === 'chevron-right'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" :class="computedClass" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <path d="M9 18l6-6-6-6" />
  </svg>

  <!-- calendar - 日曆 -->
  <svg v-else-if="icon === 'calendar'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" :class="computedClass" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
    <line x1="16" y1="2" x2="16" y2="6" />
    <line x1="8" y1="2" x2="8" y2="6" />
    <line x1="3" y1="10" x2="21" y2="10" />
  </svg>

  <!-- close - 關閉 -->
  <svg v-else-if="icon === 'close'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" :class="computedClass" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <line x1="18" y1="6" x2="6" y2="18" />
    <line x1="6" y1="6" x2="18" y2="18" />
  </svg>

  <!-- plus - 加號 -->
  <svg v-else-if="icon === 'plus'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" :class="computedClass" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <line x1="12" y1="5" x2="12" y2="19" />
    <line x1="5" y1="12" x2="19" y2="12" />
  </svg>

  <!-- warning - 警告/請假 -->
  <svg v-else-if="icon === 'warning'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" :class="computedClass" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z" />
    <line x1="12" y1="9" x2="12" y2="13" />
    <line x1="12" y1="17" x2="12.01" y2="17" />
  </svg>

  <!-- download - 下載/匯出 -->
  <svg v-else-if="icon === 'download'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" :class="computedClass" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
    <polyline points="7 10 12 15 17 10" />
    <line x1="12" y1="15" x2="12" y2="3" />
  </svg>
</template>

<script setup lang="ts">
interface Props {
  icon:
    | 'user'
    | 'lock'
    | 'chat'
    | 'clock'
    | 'spinner'
    | 'chevron-left'
    | 'chevron-right'
    | 'calendar'
    | 'close'
    | 'plus'
    | 'warning'
    | 'download'
  size?: string
  class?: string
}

const props = withDefaults(defineProps<Props>(), {
  size: 'w-5 h-5',
  class: '',
})

const computedClass = computed(() => `${props.size} ${props.class}`)
</script>
