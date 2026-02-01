<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="emit('close')">
    <div class="glass-card w-full max-w-lg animate-spring" @click.stop>
      <div class="flex items-center justify-between p-4 border-b border-white/10">
        <h3 class="text-lg font-semibold text-slate-100">
          {{ $t('exception.requestDetails') }}
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
            <h4 class="text-sm font-medium text-slate-300 mb-1">{{ $t('schedule.course') }}</h4>
            <p class="text-slate-100 font-medium">{{ props.exception?.offering_name }}</p>
          </div>

          <div>
            <h4 class="text-sm font-medium text-slate-300 mb-1">{{ $t('schedule.teacher') }}</h4>
            <p class="text-slate-100">{{ props.exception?.teacher_name }}</p>
          </div>

          <div>
            <h4 class="text-sm font-medium text-slate-300 mb-1">{{ $t('exception.type') }}</h4>
            <span
              class="px-3 py-1 rounded-full text-sm font-medium"
              :class="getTypeClass(props.exception?.type)"
            >
              {{ getTypeText(props.exception?.type) }}
            </span>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <h4 class="text-sm font-medium text-slate-300 mb-1">{{ $t('exception.originalTime') }}</h4>
              <p class="text-slate-100">{{ getOriginalTimeText(props.exception) }}</p>
            </div>
            <div>
              <h4 class="text-sm font-medium text-slate-300 mb-1">{{ $t('exception.originalDateLabel') }}</h4>
              <p class="text-slate-100">{{ formatDate(props.exception?.original_date) }}</p>
            </div>
          </div>

          <div v-if="props.exception?.type === 'RESCHEDULE'" class="grid grid-cols-2 gap-4">
            <div>
              <h4 class="text-sm font-medium text-slate-300 mb-1">{{ $t('exception.newTime') }}</h4>
              <p class="text-success-500 font-medium">{{ props.exception?.new_time }}</p>
            </div>
            <div v-if="props.exception?.new_teacher_name">
              <h4 class="text-sm font-medium text-slate-300 mb-1">{{ $t('exception.newTeacher') }}</h4>
              <p class="text-primary-500 font-medium">{{ props.exception?.new_teacher_name }}</p>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <h4 class="text-sm font-medium text-slate-300 mb-1">{{ $t('exception.status') }}</h4>
              <span
                class="px-3 py-1 rounded-full text-sm font-medium"
                :class="getStatusClass(props.exception?.status)"
              >
                {{ getStatusText(props.exception?.status) }}
              </span>
            </div>
            <div>
              <h4 class="text-sm font-medium text-slate-300 mb-1">{{ $t('exception.appliedDate') }}</h4>
              <p class="text-slate-400">{{ formatDateTime(props.exception?.created_at) }}</p>
            </div>
          </div>

          <div>
            <h4 class="text-sm font-medium text-slate-300 mb-1">{{ $t('exception.reason') }}</h4>
            <p class="text-slate-300 bg-slate-700/50 rounded-xl p-3">
              {{ props.exception?.reason }}
            </p>
          </div>

          <div v-if="props.exception?.review_note">
            <h4 class="text-sm font-medium text-slate-300 mb-1">{{ $t('exception.reviewNote') }}</h4>
            <p class="text-slate-300 bg-primary-500/10 border border-primary-500/30 rounded-xl p-3">
              {{ props.exception?.review_note }}
            </p>
          </div>
        </div>
      </div>

      <div class="flex p-4 border-t border-white/10">
        <button
          @click="emit('close')"
          class="w-full glass-btn py-3 rounded-xl font-medium"
        >
          {{ $t('common.close') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { formatDate, formatDateTime } from '~/composables/useTaiwanTime'

const props = defineProps<{
  exception: any
}>()

const emit = defineEmits<{
  close: []
}>()

// 使用 i18n
const { t: $t } = useI18n()

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
      return '調課'
    default:
      return type
  }
}

const getStatusClass = (status: string): string => {
  switch (status) {
    case 'PENDING':
      return 'bg-warning-500/20 text-warning-500'
    case 'APPROVED':
      return 'bg-success-500/20 text-success-500'
    case 'REJECTED':
      return 'bg-critical-500/20 text-critical-500'
    default:
      return 'bg-slate-500/20 text-slate-400'
  }
}

const getStatusText = (status: string): string => {
  switch (status) {
    case 'PENDING':
      return '待審核'
    case 'APPROVED':
      return '已核准'
    case 'REJECTED':
      return '已拒絕'
    default:
      return status
  }
}

const getOriginalTimeText = (exception: any): string => {
  if (exception.original_date && exception.rule) {
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
</script>
