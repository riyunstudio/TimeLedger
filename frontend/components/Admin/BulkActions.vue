<template>
  <transition
    enter-active-class="transition-all duration-300 ease-out"
    enter-from-class="translate-y-full opacity-0"
    enter-to-class="translate-y-0 opacity-100"
    leave-active-class="transition-all duration-300 ease-in"
    leave-from-class="translate-y-0 opacity-100"
    leave-to-class="translate-y-full opacity-0"
  >
    <div
      v-if="selectedCount > 0"
      class="bulk-actions fixed bottom-0 left-0 right-0 z-50 bg-slate-900/95 backdrop-blur-lg border-t border-white/10 shadow-2xl"
    >
      <div class="max-w-7xl mx-auto px-4 py-4">
        <div class="flex flex-col sm:flex-row items-center justify-between gap-4">
          <!-- 選取資訊 -->
          <div class="flex items-center gap-3">
            <BaseBadge variant="info" size="md">
              已選取 {{ selectedCount }} 位
            </BaseBadge>
            <button
              @click="$emit('clear')"
              class="text-sm text-slate-400 hover:text-slate-300 transition-colors"
            >
              清除選取
            </button>
          </div>

          <!-- 操作按鈕 -->
          <div class="flex items-center gap-3">
            <button
              @click="$emit('compare')"
              class="px-4 py-2 rounded-lg bg-white/5 text-slate-300 hover:bg-white/10 hover:text-white transition-colors flex items-center gap-2"
              :disabled="selectedCount < 2"
              :class="{ 'opacity-50 cursor-not-allowed': selectedCount < 2 }"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
              比較所選
            </button>

            <button
              @click="$emit('export')"
              class="px-4 py-2 rounded-lg bg-white/5 text-slate-300 hover:bg-white/10 hover:text-white transition-colors flex items-center gap-2"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
              </svg>
              匯出聯絡資訊
            </button>

            <button
              @click="$emit('bulkInvite')"
              :disabled="bulkLoading"
              class="px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors flex items-center gap-2 disabled:opacity-50"
            >
              <svg v-if="bulkLoading" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
              </svg>
              <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
              </svg>
              批量邀請 ({{ selectedCount }})
            </button>
          </div>
        </div>

        <!-- 進度提示 -->
        <div
          v-if="bulkLoading"
          class="mt-3 flex items-center gap-2 text-sm text-slate-400"
        >
          <div class="flex-1 bg-white/10 rounded-full h-1.5 overflow-hidden">
            <div
              class="h-full bg-primary-500 transition-all duration-300"
              :style="{ width: `${bulkProgress}%` }"
            />
          </div>
          <span>{{ bulkProgress }}%</span>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
defineProps<{
  selectedCount: number
  bulkLoading?: boolean
  bulkProgress?: number
}>()

defineEmits<{
  clear: []
  compare: []
  export: []
  bulkInvite: []
}>()
</script>
