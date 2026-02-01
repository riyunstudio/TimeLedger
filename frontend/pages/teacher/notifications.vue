<template>
  <div class="p-4 md:p-6 max-w-7xl mx-auto">
    <div class="mb-6 md:mb-8">
      <h1 class="text-2xl md:text-3xl font-bold text-slate-100 mb-2">
        通知歷史
      </h1>
      <p class="text-slate-400 text-sm md:text-base">
        查看您的所有通知紀錄
      </p>
    </div>

    <!-- 篩選器 -->
    <div class="flex flex-wrap gap-3 mb-6">
      <select
        v-model="filters.readStatus"
        class="px-4 py-2 bg-white/5 border border-white/10 rounded-xl text-slate-300 focus:outline-none focus:border-primary-500"
      >
        <option value="">全部</option>
        <option value="unread">未讀</option>
        <option value="read">已讀</option>
      </select>

      <select
        v-model="filters.type"
        class="px-4 py-2 bg-white/5 border border-white/10 rounded-xl text-slate-300 focus:outline-none focus:border-primary-500"
      >
        <option value="">全部類型</option>
        <option value="APPROVAL">審核通過</option>
        <option value="REVIEW_RESULT">審核結果</option>
        <option value="EXCEPTION">例外申請</option>
        <option value="SCHEDULE">課表通知</option>
        <option value="CENTER_INVITE">中心邀請</option>
      </select>

      <button
        @click="refreshData"
        class="px-4 py-2 bg-primary-500/20 border border-primary-500 text-primary-400 rounded-xl hover:bg-primary-500/30 transition-colors"
      >
        重新整理
      </button>
    </div>

    <!-- 通知列表 -->
    <div class="bg-white/5 rounded-xl border border-white/10 overflow-hidden">
      <div v-if="loading" class="p-8 text-center">
        <div class="inline-block w-8 h-8 border-2 border-primary-500 border-t-transparent rounded-full animate-spin"></div>
        <p class="text-slate-400 mt-2">載入中...</p>
      </div>

      <div v-else-if="filteredNotifications.length === 0" class="p-8 text-center text-slate-400">
        暫無通知記錄
      </div>

      <div v-else class="divide-y divide-white/5">
        <div
          v-for="notification in filteredNotifications"
          :key="notification.id"
          class="p-4 hover:bg-white/5 transition-colors cursor-pointer"
          :class="{ 'bg-primary-500/10': !notification.is_read }"
          @click="handleNotificationClick(notification)"
        >
          <div class="flex items-start gap-3">
            <div class="flex-1">
              <div class="flex items-center gap-2 mb-1">
                <h4 class="font-medium text-slate-100 text-sm">
                  {{ notification.title }}
                </h4>
                <span
                  v-if="!notification.is_read"
                  class="w-2 h-2 rounded-full bg-primary-500 shrink-0"
                />
              </div>
              <p class="text-sm text-slate-400">
                {{ notification.message }}
              </p>
              <div class="flex items-center gap-3 mt-2">
                <p class="text-xs text-slate-500">
                  {{ formatTime(notification.created_at) }}
                </p>
                <span
                  class="px-2 py-0.5 rounded text-xs"
                  :class="{
                    'bg-primary-500/20 text-primary-400': notification.type === 'APPROVAL',
                    'bg-success-500/20 text-success-400': notification.type === 'REVIEW_RESULT',
                    'bg-warning-500/20 text-warning-400': notification.type === 'EXCEPTION',
                    'bg-info-500/20 text-info-400': notification.type === 'SCHEDULE',
                    'bg-secondary-500/20 text-secondary-400': notification.type === 'CENTER_INVITE',
                  }"
                >
                  {{ typeText(notification.type) }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 分頁 -->
      <div v-if="pagination.totalPages > 1" class="px-4 py-3 bg-white/5 border-t border-white/10 flex items-center justify-between">
        <p class="text-slate-400 text-sm">
          第 {{ pagination.page }} 頁，共 {{ pagination.totalPages }} 頁
        </p>
        <div class="flex gap-2">
          <button
            @click="changePage(pagination.page - 1)"
            :disabled="pagination.page === 1"
            class="px-3 py-1 rounded-lg bg-white/5 text-slate-300 hover:bg-white/10 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            上一頁
          </button>
          <button
            @click="changePage(pagination.page + 1)"
            :disabled="pagination.page === pagination.totalPages"
            class="px-3 py-1 rounded-lg bg-white/5 text-slate-300 hover:bg-white/10 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            下一頁
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  auth: 'TEACHER',
  layout: 'default',
})

const router = useRouter()

interface Notification {
  id: number
  user_id: number
  user_type: string
  title: string
  message: string
  is_read: boolean
  created_at: string
  type: string
}

const loading = ref(true)
const notifications = ref<Notification[]>([])

const filters = ref({
  readStatus: '',
  type: '',
})

const pagination = ref({
  page: 1,
  limit: 20,
  total: 0,
  totalPages: 0,
})

const typeText = (type: string) => {
  const texts: Record<string, string> = {
    APPROVAL: '審核通過',
    REVIEW_RESULT: '審核結果',
    EXCEPTION: '例外申請',
    SCHEDULE: '課表通知',
    CENTER_INVITE: '中心邀請',
  }
  return texts[type] || type
}

const formatTime = (dateStr: string): string => {
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / (1000 * 60))
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60))
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

  if (diffMins < 1) return '剛剛'
  if (diffMins < 60) return `${diffMins} 分鐘前`
  if (diffHours < 24) return `${diffHours} 小時前`
  if (diffDays < 7) return `${diffDays} 天前`

  return date.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

const fetchNotifications = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams({
      page: pagination.value.page.toString(),
      limit: pagination.value.limit.toString(),
    })

    const api = useApi()
    const response = await api.get<{ code: number; message: string; datas: { notifications: Notification[]; total: number; unread_count: number } }>(
      `/notifications?${params}`
    )
    notifications.value = response.datas?.notifications || []
    pagination.value.total = response.datas?.total || 0
    pagination.value.totalPages = Math.ceil(pagination.value.total / pagination.value.limit)
  } catch (error) {
    console.error('Failed to fetch notifications:', error)
  } finally {
    loading.value = false
  }
}

const refreshData = () => {
  pagination.value.page = 1
  fetchNotifications()
}

const changePage = (page: number) => {
  pagination.value.page = page
  fetchNotifications()
}

const markAsRead = async (id: number) => {
  try {
    const api = useApi()
    await api.post(`/notifications/${id}/read`, {})
    const notification = notifications.value.find((n) => n.id === id)
    if (notification) {
      notification.is_read = true
    }
  } catch (error) {
    console.error('Failed to mark as read:', error)
  }
}

const handleNotificationClick = async (notification: Notification) => {
  if (!notification.is_read) {
    await markAsRead(notification.id)
  }

  // 根據通知類型跳轉
  if (notification.type === 'REVIEW_RESULT' || notification.title === '例外單審核結果') {
    router.push('/teacher/exceptions')
  } else if (notification.type === 'APPROVAL') {
    // 審核通過通知，可能需要跳轉到課表或其他頁面
  }
}

const filteredNotifications = computed(() => {
  let result = notifications.value

  if (filters.value.readStatus === 'unread') {
    result = result.filter((n) => !n.is_read)
  } else if (filters.value.readStatus === 'read') {
    result = result.filter((n) => n.is_read)
  }

  if (filters.value.type) {
    result = result.filter((n) => n.type === filters.value.type)
  }

  return result
})

// 監聽篩選器變化
watch(filters, () => {
  pagination.value.page = 1
  fetchNotifications()
}, { deep: true })

onMounted(() => {
  fetchNotifications()
})
</script>
