<template>
  <div class="p-4 md:p-6 max-w-7xl mx-auto">
    <div class="mb-6 md:mb-8">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl md:text-3xl font-bold text-slate-100 mb-2">
            審核中心
          </h1>
          <p class="text-slate-400 text-sm md:text-base">
            處理課程變更申請（停課/改期/找代課）
          </p>
        </div>
        <!-- 即時更新指示器 -->
        <div class="flex items-center gap-3">
          <div class="flex items-center gap-2 text-xs text-slate-500">
            <span
              class="w-2 h-2 rounded-full animate-pulse"
              :class="isPolling ? 'bg-success-500' : 'bg-slate-500'"
            ></span>
            <span v-if="lastUpdated">更新於 {{ formatLastUpdated }}</span>
            <span v-else>尚未更新</span>
          </div>
          <button
            @click="refreshData"
            :disabled="isRefreshing"
            class="p-2 rounded-lg hover:bg-white/10 transition-colors"
            title="手動重新整理"
          >
            <svg
              class="w-5 h-5 text-slate-400"
              :class="{ 'animate-spin': isRefreshing }"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- 篩選器區域 -->
    <div class="mb-6 flex flex-col gap-4">
      <!-- 第一行：狀態篩選 -->
      <div class="flex flex-wrap gap-2 overflow-x-auto pb-2">
        <button
          @click="activeFilter = 'all'"
          class="glass-btn px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap"
          :class="activeFilter === 'all' ? 'bg-primary-500/30 border-primary-500' : ''"
        >
          全部
        </button>
        <button
          @click="activeFilter = 'pending'"
          class="glass-btn px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap"
          :class="activeFilter === 'pending' ? 'bg-warning-500/30 border-warning-500' : ''"
        >
          待審核 ({{ pendingCount }})
        </button>
        <button
          @click="activeFilter = 'approved'"
          class="glass-btn px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap"
          :class="activeFilter === 'approved' ? 'bg-success-500/30 border-success-500' : ''"
        >
          已核准
        </button>
        <button
          @click="activeFilter = 'rejected'"
          class="glass-btn px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap"
          :class="activeFilter === 'rejected' ? 'bg-critical-500/30 border-critical-500' : ''"
        >
          已拒絕
        </button>
        <button
          @click="activeFilter = 'revoked'"
          class="glass-btn px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap"
          :class="activeFilter === 'revoked' ? 'bg-slate-500/30 border-slate-500' : ''"
        >
          已撤回
        </button>

        <div class="w-px bg-white/10 mx-2"></div>

        <!-- 視角篩選 -->
        <select
          v-model="viewModeFilter"
          class="px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap bg-slate-800/80 border border-white/10 text-slate-300 focus:outline-none focus:border-primary-500"
        >
          <option value="">全部視角</option>
          <optgroup label="老師" class="bg-slate-800">
            <option v-for="teacher in teachers" :key="'t-' + teacher.id" :value="'teacher:' + teacher.id" class="bg-slate-800">
              {{ teacher.name }}
            </option>
          </optgroup>
          <optgroup label="教室" class="bg-slate-800">
            <option v-for="room in rooms" :key="'r-' + room.id" :value="'room:' + room.id" class="bg-slate-800">
              {{ room.name }}
            </option>
          </optgroup>
        </select>

        <button
          v-if="viewModeFilter"
          @click="viewModeFilter = ''"
          class="glass-btn px-3 py-2 rounded-xl text-sm font-medium whitespace-nowrap text-slate-400 hover:text-white"
        >
          清除篩選
        </button>
      </div>

      <!-- 第二行：日期範圍篩選 -->
      <div class="flex flex-wrap items-center gap-3">
        <div class="flex items-center gap-2">
          <label class="text-sm text-slate-400">日期範圍：</label>
          <input
            type="date"
            v-model="dateFrom"
            class="px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white text-sm focus:outline-none focus:border-primary-500"
          />
          <span class="text-slate-500">至</span>
          <input
            type="date"
            v-model="dateTo"
            class="px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white text-sm focus:outline-none focus:border-primary-500"
          />
        </div>
        <button
          @click="applyDateFilter"
          class="px-4 py-2 rounded-lg bg-primary-500 text-white text-sm font-medium hover:bg-primary-600 transition-colors"
        >
          套用
        </button>
        <button
          v-if="dateFrom || dateTo"
          @click="clearDateFilter"
          class="px-3 py-2 rounded-lg bg-white/5 text-slate-400 text-sm hover:bg-white/10 transition-colors"
        >
          清除
        </button>
      </div>
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
        class="glass-card p-4 md:p-5"
      >
        <div class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4 mb-4">
          <div class="flex-1">
            <h3 class="text-lg font-semibold text-slate-100 mb-1">
              {{ exception.offering_name || '課程變更' }}
            </h3>
            <div class="flex flex-wrap items-center gap-2 text-sm text-slate-400">
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
              @click="openDetail(exception)"
              aria-label="查看詳情"
              class="p-2 rounded-lg hover:bg-white/10 transition-colors"
            >
              <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </button>
          </div>
        </div>

        <div v-if="exception.type === 'RESCHEDULE'" class="space-y-2">
          <div class="flex flex-col sm:flex-row sm:items-center gap-2">
            <span class="text-slate-400 text-sm">原時間：</span>
            <span class="text-critical-500 text-sm line-through">{{ getOriginalTimeText(exception) }}</span>
          </div>
          <div v-if="exception.new_start_at" class="flex flex-col sm:flex-row sm:items-center gap-2">
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
  </div>

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
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth-admin',
  layout: 'admin',
})

const activeFilter = ref('all')
const viewModeFilter = ref('')
const showReviewModal = ref<any>(null)
const showDetailModal = ref<any>(null)
const notificationUI = useNotification()
const loading = ref(false)
const { getCenterId } = useCenterId()
const toast = useToast()

// 即時更新相關
const isPolling = ref(true)
const isRefreshing = ref(false)
const lastUpdated = ref<Date | null>(null)
let pollTimer: ReturnType<typeof setInterval> | null = null

// 日期範圍篩選
const dateFrom = ref('')
const dateTo = ref('')

const exceptions = ref<any[]>([])
const teachers = ref<any[]>([])
const rooms = ref<any[]>([])

const filteredExceptions = computed(() => {
  let result = exceptions.value

  // 狀態過濾
  if (activeFilter.value !== 'all') {
    const filterStatus = activeFilter.value.toUpperCase()
    const normalizedFilter = filterStatus.replace(/ED$/, '')
    result = result.filter(exc => {
      const normalizedStatus = exc.status.replace(/ED$/, '')
      return normalizedStatus === normalizedFilter || exc.status === filterStatus
    })
  }

  // 視角過濾
  if (viewModeFilter.value) {
    const [type, id] = viewModeFilter.value.split(':')
    const targetId = parseInt(id)
    if (type === 'teacher') {
      result = result.filter(exc => {
        // 原本的老師 (透過 Rule.Teacher.ID)
        const originalTeacherId = exc.rule?.teacher?.id
        // 代課老師 (直接欄位)
        const newTeacherId = exc.new_teacher_id
        return originalTeacherId === targetId || newTeacherId === targetId
      })
    } else if (type === 'room') {
      result = result.filter(exc => {
        // 教室 ID (透過 Rule.Room.ID)
        const roomId = exc.rule?.room?.id || exc.new_room_id
        return roomId === targetId
      })
    }
  }

  // 日期範圍過濾
  if (dateFrom.value) {
    const fromDate = new Date(dateFrom.value)
    result = result.filter(exc => new Date(exc.original_date) >= fromDate)
  }
  if (dateTo.value) {
    const toDate = new Date(dateTo.value)
    toDate.setHours(23, 59, 59, 999) // 設置為當天最後一刻
    result = result.filter(exc => new Date(exc.original_date) <= toDate)
  }

  return result
})

const pendingCount = computed(() => {
  return exceptions.value.filter(exc => exc.status === 'PENDING').length
})

// 格式化最後更新時間
const formatLastUpdated = computed(() => {
  if (!lastUpdated.value) return ''
  const now = new Date()
  const diff = Math.floor((now.getTime() - lastUpdated.value.getTime()) / 1000)

  if (diff < 60) return '剛剛'
  if (diff < 3600) return `${Math.floor(diff / 60)} 分鐘前`
  return lastUpdated.value.toLocaleTimeString('zh-TW', { hour: '2-digit', minute: '2-digit' })
})

const fetchExceptions = async () => {
  try {
    const api = useApi()
    // 查詢所有例外申請（不再只查待審核）
    const response = await api.get<{ code: number; datas: any[] }>('/admin/exceptions/all')
    exceptions.value = response.datas || []
    // 更新最後刷新時間
    lastUpdated.value = new Date()
  } catch (error) {
    console.error('Failed to fetch exceptions:', error)
    exceptions.value = []
  } finally {
    loading.value = false
  }
}

// 手動重新整理
const refreshData = async () => {
  isRefreshing.value = true
  await fetchExceptions()
  await fetchFilters()
  setTimeout(() => {
    isRefreshing.value = false
  }, 500)
}

// 開始輪詢
const startPolling = () => {
  // 每 30 秒自動刷新
  pollTimer = setInterval(() => {
    if (!document.hidden) {
      fetchExceptions()
    }
  }, 30000)
}

// 停止輪詢
const stopPolling = () => {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
}

// 套用日期篩選
const applyDateFilter = () => {
  // 日期已經綁定到 dateFrom 和 dateTo，computed 會自動過濾
}

// 清除日期篩選
const clearDateFilter = () => {
  dateFrom.value = ''
  dateTo.value = ''
}

// 監聽篩選條件變化，重新獲取數據
watch([activeFilter], () => {
  // 狀態改變時需要重新獲取數據
  fetchExceptions()
})

const fetchFilters = async () => {
  try {
    const api = useApi()
    const [teachersRes, roomsRes] = await Promise.all([
      api.get<{ code: number; datas: any[] }>('/teachers'),
      api.get<{ code: number; datas: any[] }>(`/admin/rooms`)
    ])
    teachers.value = teachersRes.datas || []
    rooms.value = roomsRes.datas || []
  } catch (error) {
    console.error('Failed to fetch filters:', error)
  }
}

const getStatusClass = (status: string): string => {
  switch (status) {
    case 'PENDING':
      return 'bg-warning-500/20 text-warning-500'
    case 'APPROVED':
    case 'APPROVE': // 向后兼容旧数据
      return 'bg-success-500/20 text-success-500'
    case 'REJECTED':
    case 'REJECT': // 向后兼容旧数据
      return 'bg-critical-500/20 text-critical-500'
    case 'REVOKED':
      return 'bg-slate-500/20 text-slate-400'
    default:
      return 'bg-slate-500/20 text-slate-400'
  }
}

const getStatusText = (status: string): string => {
  switch (status) {
    case 'PENDING':
      return '待審核'
    case 'APPROVED':
    case 'APPROVE': // 向后兼容旧数据
      return '已核准'
    case 'REJECTED':
    case 'REJECT': // 向后兼容旧数据
      return '已拒絕'
    case 'REVOKED':
      return '已撤回'
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
    case 'revoked':
      return '目前沒有已撤回的申請'
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

const getOriginalTimeText = (exception: any): string => {
  if (exception.original_time) {
    return exception.original_time
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

const handleApproved = async (_id: number, note: string) => {
  try {
    const api = useApi()
    await api.post(`/admin/scheduling/exceptions/${_id}/review`, {
      action: 'APPROVED',
      reason: note,
    })
    // 重新取得例外列表以確保資料同步（後端可能已修改規則）
    await fetchExceptions()
    // 審核通過後，切換到「已核准」標籤
    activeFilter.value = 'approved'
    toast.success('已成功核准該申請', '核准成功')
  } catch (error) {
    console.error('Failed to approve exception:', error)
    toast.error('核准失敗，請稍後再試', '操作失敗')
    return
  }
  showReviewModal.value = null
}

const handleRejected = async (_id: number, note: string) => {
  try {
    const api = useApi()
    await api.post(`/admin/scheduling/exceptions/${_id}/review`, {
      action: 'REJECTED',
      reason: note,
    })
    // 重新取得例外列表以確保資料同步
    await fetchExceptions()
    // 審核拒絕後，切換到「已拒絕」標籤
    activeFilter.value = 'rejected'
    toast.success('已成功拒絕該申請', '拒絕成功')
  } catch (error) {
    console.error('Failed to reject exception:', error)
    toast.error('拒絕失敗，請稍後再試', '操作失敗')
    return
  }
  showReviewModal.value = null
}

// 打開詳情Modal
const openDetail = (exception: any) => {
  console.log('Opening detail for exception:', exception)
  showDetailModal.value = exception
}

onMounted(async () => {
  loading.value = true
  await fetchExceptions()
  await fetchFilters()
  // 開始輪詢
  startPolling()

  // 頁面可見性變化時刷新
  document.addEventListener('visibilitychange', handleVisibilityChange)
})

onUnmounted(() => {
  stopPolling()
  document.removeEventListener('visibilitychange', handleVisibilityChange)
})

// 處理頁面可見性變化
const handleVisibilityChange = () => {
  if (!document.hidden) {
    // 頁面重新可見時刷新數據
    fetchExceptions()
  }
}
</script>
