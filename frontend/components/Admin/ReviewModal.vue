<template>
  <div class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="emit('close')">
    <div class="glass-card w-full max-w-md sm:max-w-lg max-h-[90vh] overflow-y-auto animate-spring" @click.stop>
      <div class="flex items-center justify-between p-4 border-b border-white/10 sticky top-0 bg-slate-900/95 backdrop-blur-sm z-10">
        <h3 class="text-lg font-semibold text-slate-100">
          審核申請
        </h3>
        <button @click="emit('close')" class="p-2 rounded-lg hover:bg-white/10 transition-colors">
          <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div class="p-4 space-y-4">
        <div class="space-y-3">
          <div>
            <h4 class="text-sm font-medium text-slate-300 mb-1">課程</h4>
            <p class="text-slate-100 text-sm sm:text-base">{{ props.exception?.rule?.name || '-' }}</p>
          </div>

          <div>
            <h4 class="text-sm font-medium text-slate-300 mb-1">老師</h4>
            <p class="text-slate-100 text-sm sm:text-base">{{ props.exception?.rule?.teacher?.name || '-' }}</p>
          </div>

          <div>
            <h4 class="text-sm font-medium text-slate-300 mb-1">申請類型</h4>
            <span
              class="px-3 py-1 rounded-full text-sm font-medium"
              :class="getTypeClass(props.exception?.exception_type)"
            >
              {{ getTypeText(props.exception?.exception_type) }}
            </span>
          </div>

          <div v-if="props.exception?.exception_type === 'RESCHEDULE'" class="space-y-2">
            <div class="flex items-center gap-2">
              <span class="text-slate-400 text-sm">原時間：</span>
              <span class="text-critical-500 text-sm line-through">{{ getOriginalTimeText(props.exception) }}</span>
            </div>
            <div class="flex items-center gap-2">
              <span class="text-slate-400 text-sm">新時間：</span>
              <span class="text-success-500 text-sm">{{ formatNewTime(props.exception) }}</span>
            </div>
            <div v-if="props.exception?.new_teacher_name" class="flex items-center gap-2">
              <span class="text-slate-400 text-sm">老師：</span>
              <span class="text-primary-500 text-sm">{{ props.exception?.new_teacher_name }}</span>
            </div>
          </div>

          <div>
            <h4 class="text-sm font-medium text-slate-300 mb-1">申請原因</h4>
            <p class="text-slate-300 bg-slate-700/50 rounded-xl p-3 text-sm sm:text-base">
              {{ props.exception?.reason || '-' }}
            </p>
          </div>
        </div>

        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">審核備註（選填）</label>
          <textarea
            v-model="reviewNote"
            placeholder="輸入審核備註..."
            rows="3"
            class="input-field resize-none text-sm sm:text-base"
          />
        </div>
      </div>

      <div class="flex gap-3 p-4 border-t border-white/10">
        <button
          @click="handleReject"
          class="flex-1 btn-critical py-2.5 sm:py-3 rounded-xl font-medium text-sm sm:text-base"
        >
          拒絕
        </button>
        <button
          @click="handleApprove"
          class="flex-1 btn-success py-2.5 sm:py-3 rounded-xl font-medium text-sm sm:text-base"
        >
          核准
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  exception: any
}>()

const emit = defineEmits<{
  close: []
  approved: [id: number, note: string]
  rejected: [id: number, note: string]
}>()

const reviewNote = ref('')

const getTypeClass = (type: string): string => {
  switch (type) {
    case 'CANCEL':
      return 'bg-critical-500/20 text-critical-500'
    case 'RESCHEDULE':
      return 'bg-warning-500/20 text-warning-500'
    default:
      return 'bg-slate-500/20 text-slate-400'
  }
}

const getTypeText = (type: string): string => {
  switch (type) {
    case 'CANCEL':
      return '停課'
    case 'RESCHEDULE':
      return '改期'
    default:
      return type
  }
}

const getOriginalTimeText = (exception: any): string => {
  if (exception.rule && exception.original_date) {
    const startTime = exception.rule.start_time || ''
    const endTime = exception.rule.end_time || ''
    const date = exception.original_date.split('T')[0]
    if (startTime && endTime) {
      return `${date} ${startTime} - ${endTime}`
    }
  }
  if (exception.rule) {
    const startTime = exception.rule.start_time || ''
    const endTime = exception.rule.end_time || ''
    if (startTime && endTime) {
      return `${startTime} - ${endTime}`
    }
  }
  return '-'
}

// 格式化新時間（調課申請）
const formatNewTime = (exception: any): string => {
  if (exception.new_start_at && exception.new_end_at) {
    const date = exception.new_start_at.split('T')[0]
    const time = exception.new_start_at.split('T')[1]?.substring(0, 5)
    const endTime = exception.new_end_at.split('T')[1]?.substring(0, 5)
    if (time && endTime) {
      return `${date} ${time} - ${endTime}`
    }
  }
  return exception.new_time || '-'
}

const handleApprove = () => {
  emit('approved', props.exception.id, reviewNote.value)
}

const handleReject = () => {
  emit('rejected', props.exception.id, reviewNote.value)
}
</script>
