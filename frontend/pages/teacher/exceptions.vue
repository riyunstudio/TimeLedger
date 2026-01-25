<template>
  <div class="flex items-center justify-between mb-6">
    <h2 class="text-xl font-semibold text-white">例外申請</h2>
    <button @click="showModal = true" class="px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors">
      新增申請
    </button>
  </div>

  <div class="flex gap-2 mb-4">
    <button
      v-for="status in statusFilters"
      :key="status.value"
      @click="currentFilter = status.value"
      class="px-3 py-1.5 rounded-lg text-sm font-medium transition-all"
      :class="currentFilter === status.value ? 'bg-primary-500 text-white' : 'bg-white/5 text-slate-400 hover:text-white'"
    >
      {{ status.label }}
    </button>
  </div>

  <div class="space-y-3">
    <div
      v-for="exception in filteredExceptions"
      :key="exception.id"
      class="glass-card p-4"
    >
      <div class="flex items-start justify-between">
        <div class="flex-1">
          <div class="flex items-center gap-2 mb-2">
            <span
              class="px-2 py-0.5 rounded-full text-xs font-medium"
              :class="getStatusClass(exception.status)"
            >
              {{ getStatusText(exception.status) }}
            </span>
            <span class="text-xs text-slate-500">
              {{ formatDateTime(exception.created_at) }}
            </span>
          </div>
          <div class="text-white font-medium mb-1">
            {{ exception.type === 'CANCEL' ? '停課' : '改期' }} - {{ exception.original_date }}
          </div>
          <p class="text-sm text-slate-400">{{ exception.reason }}</p>
          <div v-if="exception.type === 'RESCHEDULE'" class="mt-2 text-sm text-slate-300">
            新時間: {{ formatDateTime(exception.new_start_at || '') }} - {{ formatDateTime(exception.new_end_at || '') }}
          </div>
        </div>
        <div class="flex gap-2">
          <button
            v-if="exception.status === 'PENDING'"
            @click="handleRevoke(exception.id)"
            class="px-3 py-1 rounded-lg bg-critical-500/20 text-critical-500 text-sm hover:bg-critical-500/30 transition-colors"
          >
            撤回
          </button>
        </div>
      </div>
    </div>

    <div v-if="filteredExceptions.length === 0" class="text-center py-12 text-slate-500">
      暫無例外申請紀錄
    </div>
  </div>

  <ExceptionModal
    v-if="showModal"
    :centers="centers"
    :schedule-rules="scheduleRules"
    @close="showModal = false"
    @submit="fetchExceptions"
  />

  <NotificationDropdown v-if="notificationUI.show.value" @close="notificationUI.close()" />
  <TeacherSidebar v-if="sidebarStore.isOpen.value" @close="sidebarStore.close()" />
</template>

<script setup lang="ts">
import type { ScheduleException } from '~/types'

 definePageMeta({
   middleware: 'auth-teacher',
   layout: 'default',
 })

 const teacherStore = useTeacherStore()
const sidebarStore = useSidebar()
const notificationUI = useNotification()
const { confirm: alertConfirm } = useAlert()
const showModal = ref(false)
const currentFilter = ref('')

const statusFilters = [
  { value: '', label: '全部' },
  { value: 'PENDING', label: '待審核' },
  { value: 'APPROVED', label: '已核准' },
  { value: 'REJECTED', label: '已拒絕' },
  { value: 'REVOKED', label: '已撤回' },
]

const centers = computed(() => {
  return teacherStore.centers.map(m => ({
    center_id: m.center_id,
    center_name: m.center_name || '',
  }))
})

const scheduleRules = ref<Array<{
  id: number
  title: string
  original_date: string
  start_time: string
  end_time: string
}>>([])

const filteredExceptions = computed(() => {
  if (!currentFilter.value) return teacherStore.exceptions
  return teacherStore.exceptions.filter(e => e.status === currentFilter.value)
})

const getStatusClass = (status: string) => {
  switch (status) {
    case 'PENDING':
      return 'bg-warning-500/20 text-warning-500'
    case 'APPROVED':
      return 'bg-success-500/20 text-success-500'
    case 'REJECTED':
      return 'bg-critical-500/20 text-critical-500'
    case 'REVOKED':
      return 'bg-slate-500/20 text-slate-400'
    default:
      return 'bg-slate-500/20 text-slate-400'
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case 'PENDING':
      return '待審核'
    case 'APPROVED':
      return '已核准'
    case 'REJECTED':
      return '已拒絕'
    case 'REVOKED':
      return '已撤回'
    default:
      return status
  }
}

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-TW')
}

const fetchExceptions = async () => {
  await teacherStore.fetchExceptions(currentFilter.value || undefined)
}

const handleRevoke = async (id: number) => {
  if (await alertConfirm('確定要撤回此申請嗎？')) {
    await teacherStore.revokeException(id)
  }
}

onMounted(async () => {
  await Promise.all([
    teacherStore.fetchCenters(),
    fetchExceptions(),
  ])
})
</script>
