<template>
  <div class="min-h-screen bg-slate-900">
    <AdminHeader />

    <main class="p-6 max-w-7xl mx-auto">
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-slate-100 mb-2">
          審核中心
        </h1>
        <p class="text-slate-400">
          處理課程變更申請（停課/改期/找代課）
        </p>
      </div>

      <div class="mb-6 flex gap-3">
        <button
          @click="activeFilter = 'all'"
          class="glass-btn px-4 py-2 rounded-xl text-sm font-medium"
          :class="activeFilter === 'all' ? 'bg-primary-500/30 border-primary-500' : ''"
        >
          全部
        </button>
        <button
          @click="activeFilter = 'pending'"
          class="glass-btn px-4 py-2 rounded-xl text-sm font-medium"
          :class="activeFilter === 'pending' ? 'bg-warning-500/30 border-warning-500' : ''"
        >
          待審核 ({{ pendingCount }})
        </button>
        <button
          @click="activeFilter = 'approved'"
          class="glass-btn px-4 py-2 rounded-xl text-sm font-medium"
          :class="activeFilter === 'approved' ? 'bg-success-500/30 border-success-500' : ''"
        >
          已核准
        </button>
        <button
          @click="activeFilter = 'rejected'"
          class="glass-btn px-4 py-2 rounded-xl text-sm font-medium"
          :class="activeFilter === 'rejected' ? 'bg-critical-500/30 border-critical-500' : ''"
        >
          已拒絕
        </button>
      </div>

      <div
        v-if="filteredExceptions.length === 0"
        class="text-center py-16 text-slate-500"
      >
        {{ getEmptyMessage() }}
      </div>

      <div
        v-else
        class="space-y-4"
      >
        <div
          v-for="exception in filteredExceptions"
          :key="exception.id"
          class="glass-card p-5"
        >
          <div class="flex items-start justify-between mb-4">
            <div class="flex-1">
              <h3 class="text-lg font-semibold text-slate-100 mb-1">
                {{ exception.offering_name || '課程變更' }}
              </h3>
              <div class="flex items-center gap-3 text-sm text-slate-400">
                <span>{{ formatDate(exception.original_date) }}</span>
                <span class="px-2 py-1 rounded-full text-xs font-medium"
                  :class="getStatusClass(exception.status)"
                >
                  {{ getStatusText(exception.status) }}
                </span>
              </div>
            </div>

            <div class="flex items-center gap-2">
              <button
                v-if="exception.status === 'PENDING'"
                @click="showReviewModal = exception"
                class="px-4 py-2 rounded-xl bg-primary-500 text-white text-sm font-medium hover:bg-primary-600 transition-colors"
              >
                審核
              </button>
              <button
                @click="showDetailModal = exception"
                class="p-2 rounded-lg hover:bg-white/10 transition-colors"
              >
                <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9 9h.01" />
                </svg>
              </button>
            </div>
          </div>

          <div v-if="exception.type === 'RESCHEDULE'" class="space-y-2">
            <div class="flex items-center gap-2">
              <span class="text-slate-400 text-sm">原時間：</span>
              <span class="text-critical-500 text-sm line-through">{{ exception.original_time || exception.start_time + ' - ' + exception.end_time }}</span>
            </div>
            <div v-if="exception.new_start_at" class="flex items-center gap-2">
              <span class="text-slate-400 text-sm">新時間：</span>
              <span class="text-success-500 text-sm">{{ formatDateTime(exception.new_start_at) }} - {{ formatDateTime(exception.new_end_at) }}</span>
            </div>
            <div v-if="exception.new_teacher_name" class="flex items-center gap-2">
              <span class="text-slate-400 text-sm">代課老師：</span>
              <span class="text-primary-500 text-sm">{{ exception.new_teacher_name }}</span>
            </div>
          </div>

          <div v-if="exception.reason" class="p-3 rounded-xl bg-slate-700/50 border border-white/10">
            <p class="text-sm text-slate-300">
              <span class="font-medium text-slate-400">原因：</span>
              {{ exception.reason }}
            </p>
          </div>
        </div>
      </div>
    </main>

    <ReviewModal
      v-if="showReviewModal"
      :exception="showReviewModal"
      @close="showReviewModal = null"
      @approved="handleApproved"
      @rejected="handleRejected"
    />

    <ExceptionDetailModal
      v-if="showDetailModal"
      :exception="showDetailModal"
      @close="showDetailModal = null"
    />

    <NotificationDropdown
      v-if="notificationUI.show.value"
      @close="notificationUI.close()"
    />
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth-admin',
})

const activeFilter = ref('all')
const showReviewModal = ref<any>(null)
const showDetailModal = ref<any>(null)
const notificationUI = useNotification()
const loading = ref(false)
const { getCenterId } = useCenterId()

const exceptions = ref<any[]>([])

const filteredExceptions = computed(() => {
  if (activeFilter.value === 'all') return exceptions.value

  return exceptions.value.filter(exc => {
    return exc.status === activeFilter.value.toUpperCase()
  })
})

const pendingCount = computed(() => {
  return exceptions.value.filter(exc => exc.status === 'PENDING').length
})

const fetchExceptions = async () => {
  loading.value = true
  try {
    const api = useApi()
    const centerId = getCenterId()
    const today = new Date()
    const firstDayOfMonth = new Date(today.getFullYear(), today.getMonth(), 1)
    const lastDayOfMonth = new Date(today.getFullYear(), today.getMonth() + 1, 0)
    const startDate = firstDayOfMonth.toISOString().split('T')[0]
    const endDate = lastDayOfMonth.toISOString().split('T')[0]
    const response = await api.get<{ code: number; datas: any[] }>(`/admin/exceptions?start_date=${startDate}&end_date=${endDate}`)
    exceptions.value = response.datas || []
  } catch (error) {
    console.error('Failed to fetch exceptions:', error)
    exceptions.value = []
  } finally {
    loading.value = false
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

const getEmptyMessage = (): string => {
  switch (activeFilter.value) {
    case 'pending':
      return '目前沒有待審核的申請'
    case 'approved':
      return '目前沒有已核准的申請'
    case 'rejected':
      return '目前沒有被拒絕的申請'
    default:
      return '目前沒有任何申請'
  }
}

const formatDate = (dateStr: string): string => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'short',
  })
}

const formatDateTime = (dateStr: string): string => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-TW')
}

const handleApproved = () => {
  const exception = exceptions.value.find(e => e.id === showReviewModal.value.id)
  if (exception) {
    exception.status = 'APPROVED'
  }
  showReviewModal.value = null
}

const handleRejected = () => {
  const exception = exceptions.value.find(e => e.id === showReviewModal.value.id)
  if (exception) {
    exception.status = 'REJECTED'
  }
  showReviewModal.value = null
}

onMounted(() => {
  fetchExceptions()
})
</script>
