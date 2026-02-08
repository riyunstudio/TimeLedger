<template>
  <div
    class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm isolate"
    @click.self="$emit('close')"
  >
    <div class="glass-card w-full max-w-md max-h-[90vh] overflow-y-auto animate-spring" @click.stop>
      <div class="flex items-center justify-between p-4 border-b border-white/10 sticky top-0 bg-slate-900/95 backdrop-blur-sm z-10">
        <h3 class="text-lg font-semibold text-slate-100">
          {{ schedule?.offering_name || '選擇時段' }}
        </h3>
        <button
          @click="$emit('close')"
          class="p-2 rounded-lg hover:bg-white/10 transition-colors"
        >
          <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div v-if="schedule" class="space-y-4 p-4">
        <!-- 課程資訊 -->
        <div class="glass p-3 rounded-xl">
          <h4 class="text-sm font-medium text-slate-300 mb-2">課程資訊</h4>
          <div class="space-y-2 text-sm">
            <div class="flex justify-between">
              <span class="text-slate-400">課程</span>
              <span class="text-slate-100">{{ schedule.offering_name }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-slate-400">老師</span>
              <span class="text-slate-100">{{ schedule.teacher_name }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-slate-400">教室</span>
              <span class="text-slate-100">{{ schedule.room_name }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-slate-400">時間</span>
              <span class="text-slate-100">
                {{ schedule.start_time }} - {{ schedule.end_time }}
              </span>
            </div>
          </div>
        </div>

        <!-- 操作按鈕 -->
        <div class="space-y-3">
          <div class="flex items-center gap-2">
            <button
              @click="handleEdit"
              class="flex-1 glass-btn py-3 rounded-xl font-medium flex items-center justify-center gap-2"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 0L21.828 3.172a2 2 0 010-2.828l-7-7a2 2 0 00-2.828 0L2.172 20.828a2 2 0 010 2.828l7 7a2 2 0 0012.828 0l7.172-7.172z" />
              </svg>
              編輯
            </button>
          </div>

          <div class="flex items-center gap-2">
            <button
              @click="handleSuspend"
              class="flex-1 bg-warning-500/20 border border-warning-500/30 text-warning-400 py-3 rounded-xl font-medium flex items-center justify-center gap-2 hover:bg-warning-500/30 transition-colors"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              停課
            </button>
          </div>

          <div class="flex items-center gap-2">
            <button
              @click="handleDelete"
              class="flex-1 btn-critical py-3 rounded-xl font-medium flex items-center justify-center gap-2"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
              刪除
            </button>
          </div>
        </div>
      </div>

      <div
        v-else
        class="text-center py-12 text-slate-500"
      >
        點擊左側網格查看詳情
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { alertConfirm } from '~/composables/useAlert'

const props = defineProps<{
  time?: number
  weekday?: number
  schedule?: any
  validation?: any
}>()

const emit = defineEmits<{
  close: []
  edit: []
  delete: []
  create: []
  findSubstitute: []
  suspend: []
}>()

const formatTime = (hour: number): string => {
  return `${hour}:00`
}

const handleEdit = () => {
  emit('edit')
}

const handleDelete = async () => {
  // 第一次確認
  const firstConfirm = await alertConfirm(
    `確定要刪除「${props.schedule?.offering_name}」嗎？`,
    '刪除確認'
  )
  if (!firstConfirm) return

  // 第二次確認（double check）
  const secondConfirm = await alertConfirm(
    '此操作將無法復原，確定要刪除嗎？',
    '再次確認'
  )
  if (secondConfirm) {
    emit('delete')
  }
}

const handleCreate = () => {
  emit('create')
}

const handleFindSubstitute = () => {
  emit('findSubstitute')
}

const handleSuspend = () => {
  emit('suspend')
}
</script>
